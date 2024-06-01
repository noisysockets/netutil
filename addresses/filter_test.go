// SPDX-License-Identifier: MPL-2.0
/*
 * Copyright (C) 2024 The Noisy Sockets Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package addresses_test

import (
	"net/netip"
	"testing"

	"github.com/noisysockets/netutil/addresses"
	"github.com/stretchr/testify/require"
)

func TestFilterByNetwork(t *testing.T) {
	addrs := []netip.Addr{
		netip.MustParseAddr("192.168.1.10"),
		netip.MustParseAddr("172.16.5.4"),
		netip.MustParseAddr("10.0.0.1"),
		netip.MustParseAddr("203.0.113.45"),
		netip.MustParseAddr("8.8.8.8"),
		netip.MustParseAddr("2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
		netip.MustParseAddr("fd00::1"),
		netip.MustParseAddr("2607:f8b0:4005:805::200e"),
		netip.MustParseAddr("fe80::1ff:fe23:4567:890a"),
		netip.MustParseAddr("2400:cb00:2049:1::a29f:1801"),
	}

	t.Run("ip", func(t *testing.T) {
		filteredAddrs := addresses.FilterByNetwork(addrs, "ip")
		require.Equal(t, addrs, filteredAddrs)
	})

	t.Run("ip4", func(t *testing.T) {
		filteredAddrs := addresses.FilterByNetwork(addrs, "ip4")
		require.Len(t, filteredAddrs, 5)
		require.Equal(t, netip.MustParseAddr("192.168.1.10"), filteredAddrs[0])
	})

	t.Run("ip6", func(t *testing.T) {
		filteredAddrs := addresses.FilterByNetwork(addrs, "ip6")
		require.Len(t, filteredAddrs, 5)
		require.Equal(t, netip.MustParseAddr("2001:0db8:85a3::8a2e:0370:7334"), filteredAddrs[0])
	})
}
