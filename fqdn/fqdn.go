// SPDX-License-Identifier: MPL-2.0
/*
 * Copyright (C) 2024 The Noisy Sockets Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * Portions of this file are based on code originally:
 *
 * Copyright since 2015 Showmax s.r.o.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fqdn

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
)

// isalnum(3p) in POSIX locale
func isalnum(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9')
}

const (
	maxHostnameLen = 254
)

// Validate hostname, based on musl-c version of this function.
func isValidHostname(s string) bool {
	if len(s) > maxHostnameLen {
		return false
	}

	for _, c := range s {
		if !(c >= 0x80 || c == '.' || c == '-' || isalnum(c)) {
			return false
		}
	}

	return true
}

func parseHostLine(host string, line string) (string, bool) {
	const (
		StateSkipWhite = iota
		StateIp
		StateCanonFirst
		StateCanon
		StateAliasFirst
		StateAlias
	)

	var (
		canon     string
		state     int
		nextState int

		i     int
		start int
	)

	isWhite := func(b byte) bool {
		return b == ' ' || b == '\t'
	}
	isLast := func() bool {
		return i == len(line)-1 || isWhite(line[i+1])
	}
	partStr := func() string {
		return line[start : i+1]
	}

	state = StateSkipWhite
	nextState = StateIp

	slog.Debug("Looking for host", slog.String("host", host), slog.String("line", line))
	for i = 0; i < len(line); i += 1 {
		if line[i] == '#' {
			slog.Debug("Found comment, terminating")
			break
		}

		switch state {
		case StateSkipWhite:
			if !isWhite(line[i]) {
				state = nextState
				i -= 1
			}
		case StateIp:
			if isLast() {
				state = StateSkipWhite
				nextState = StateCanonFirst
			}
		case StateCanonFirst:
			start = i
			state = StateCanon
			i -= 1
		case StateCanon:
			slog.Debug("Canon so far", slog.String("part", partStr()))
			if isLast() {
				canon = partStr()
				if !isValidHostname(canon) {
					return "", false
				}

				if canon == host {
					slog.Debug("Canon match")
					return canon, true
				}

				state = StateSkipWhite
				nextState = StateAliasFirst
			}
		case StateAliasFirst:
			start = i
			state = StateAlias
			i -= 1
		case StateAlias:
			slog.Debug("Alias so far", slog.String("part", partStr()))
			if isLast() {
				alias := partStr()
				if alias == host {
					slog.Debug("Alias match")
					return canon, true
				}

				state = StateSkipWhite
				nextState = StateAliasFirst
			}
		default:
			panic(fmt.Sprintf("BUG: State not handled: %d", state))
		}
	}

	slog.Debug("No match")
	return "", false
}

// Reads hosts(5) file and tries to get canonical name for host.
func fromHosts(host string) (string, error) {
	var (
		fqdn string
		line string
		err  error
		file *os.File
		r    *bufio.Reader
		ok   bool
	)

	file, err = os.Open(hostsPath)
	if err != nil {
		err = fmt.Errorf("cannot open hosts file: %w", err)
		goto out
	}
	defer file.Close()

	r = bufio.NewReader(file)
	for line, err = readline(r); err == nil; line, err = readline(r) {
		fqdn, ok = parseHostLine(host, line)
		if ok {
			goto out
		}
	}

	if err != io.EOF {
		err = fmt.Errorf("failed to read file: %w", err)
		goto out
	}
	err = errFqdnNotFound{}

out:
	return fqdn, err
}

func fromLookup(host string) (string, error) {
	var (
		fqdn  string
		err   error
		addrs []net.IP
		hosts []string
	)

	fqdn, err = net.LookupCNAME(host)
	if err == nil && len(fqdn) != 0 {
		slog.Debug("LookupCNAME success", slog.String("fqdn", fqdn))
		goto out
	}
	slog.Debug("LookupCNAME failed", slog.Any("error", err))

	slog.Debug("Looking up", slog.String("host", host))
	addrs, err = net.LookupIP(host)
	if err != nil {
		err = errFqdnNotFound{err}
		goto out
	}
	slog.Debug("Resolved addrs", slog.Any("addrs", addrs))

	for _, addr := range addrs {
		slog.Debug("Trying", slog.Any("addr", addr))
		hosts, err = net.LookupAddr(addr.String())
		// On windows it can return err == nil but empty list of hosts
		if err != nil || len(hosts) == 0 {
			continue
		}
		slog.Debug("Resolved hosts", slog.Any("hosts", hosts))

		// First one should be the canonical hostname
		fqdn = hosts[0]

		goto out
	}

	err = errFqdnNotFound{}
out:
	// For some reason we wanted the canonical hostname without
	// trailing dot. So if it is present, strip it.
	if len(fqdn) > 0 && fqdn[len(fqdn)-1] == '.' {
		fqdn = fqdn[:len(fqdn)-1]
	}

	return fqdn, err
}

// Try to get fully qualified hostname for current machine.
//
// It tries to mimic how `hostname -f` works, so except for few edge cases you
// should get the same result from both. One thing that needs to be mentioned is
// that it does not guarantee that you get back fqdn. There is no way to do that
// and `hostname -f` can also return non-fqdn hostname if your /etc/hosts is
// fucked up.
//
// It checks few sources in this order:
//
//  1. hosts file
//     It parses hosts file if present and readable and returns first canonical
//     hostname that also references your hostname. See hosts(5) for more
//     details.
//  2. dns lookup
//     If lookup in hosts file fails, it tries to ask dns.
//
// If none of steps above succeeds, ErrFqdnNotFound is returned as error. You
// will probably want to just use output from os.Hostname() at that point.
func FqdnHostname() (string, error) {
	var (
		fqdn string
		host string
		err  error
	)

	host, err = os.Hostname()
	if err != nil {
		err = errHostnameFailed{err}
		goto out
	}
	slog.Debug("Hostname", slog.String("host", host))

	fqdn, err = fromHosts(host)
	if err == nil {
		slog.Debug("Fqdn fetched from hosts", slog.String("fqdn", fqdn))
		goto out
	}

	fqdn, err = fromLookup(host)
	if err == nil {
		slog.Debug("Fqdn fetched from lookup", slog.String("fqdn", fqdn))
		goto out
	}

	slog.Debug("Fqdn fetch failed", slog.Any("error", err))
out:
	return fqdn, err
}
