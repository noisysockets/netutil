// SPDX-License-Identifier: MPL-2.0
/*
 * Copyright (C) 2024 The Noisy Sockets Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package cidr

import (
	"errors"
	"net/netip"

	"github.com/noisysockets/netutil/uint128"
)

var (
	// ErrHostNumberOutOfRange is returned when the host number is out of range.
	ErrHostNumberOutOfRange = errors.New("host number out of range")
)

// Host returns the n-th host address in the given prefix.
func Host(prefix netip.Prefix, num int) (netip.Addr, error) {
	prefixAddrBytes := prefix.Addr().AsSlice()

	var b [16]byte
	copy(b[16-len(prefixAddrBytes):], prefixAddrBytes)

	intVal := uint128.FromBytesBE(b[:]).Add(uint128.From64(uint64(num)))
	intValBytes := intVal.BytesBE()

	var addr netip.Addr
	if intVal.Hi == 0 {
		addr = netip.AddrFrom4([4]byte(intValBytes[12:]))
	} else {
		addr = netip.AddrFrom16(intValBytes)
	}

	// Check if the address is within the prefix.
	if !prefix.Contains(addr) {
		return netip.Addr{}, ErrHostNumberOutOfRange
	}

	return addr, nil
}
