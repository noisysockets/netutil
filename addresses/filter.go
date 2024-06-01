// SPDX-License-Identifier: MPL-2.0
/*
 * Copyright (C) 2024 The Noisy Sockets Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package addresses

import "net/netip"

// FilterByNetwork returns a slice of netip.Addr matching the given network family.
// Known networks are: "ip", "ip4", and "ip6".
func FilterByNetwork(addrs []netip.Addr, network string) []netip.Addr {
	if network == "ip" {
		return addrs
	} else if network == "ip4" {
		return filterByIPv4(addrs)
	} else if network == "ip6" {
		return filterByIPv6(addrs)
	}

	return nil
}

func filterByIPv4(addrs []netip.Addr) []netip.Addr {
	var filteredAddrs []netip.Addr
	for _, addr := range addrs {
		if addr.Is4() {
			filteredAddrs = append(filteredAddrs, addr)
		}
	}
	return filteredAddrs
}

func filterByIPv6(addrs []netip.Addr) []netip.Addr {
	var filtered []netip.Addr
	for _, addr := range addrs {
		if addr.Is6() {
			filtered = append(filtered, addr)
		}
	}
	return filtered
}
