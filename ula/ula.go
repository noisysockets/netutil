// SPDX-License-Identifier: MPL-2.0
/*
 * Copyright (C) 2024 The Noisy Sockets Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package ula

import (
	"crypto/rand"
	"net/netip"
)

// Generate a new Unique Local Address (ULA) IPv6 prefix.
func Generate() (netip.Prefix, error) {
	var ipv6Addr [16]byte
	ipv6Addr[0] = 0xfd
	_, err := rand.Read(ipv6Addr[1:6])
	if err != nil {
		return netip.Prefix{}, err
	}

	return netip.PrefixFrom(netip.AddrFrom16(ipv6Addr), 48), nil
}
