// SPDX-License-Identifier: MPL-2.0
/*
 * Copyright (C) 2024 The Noisy Sockets Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package cidr_test

import (
	"net/netip"
	"testing"

	"github.com/noisysockets/netutil/cidr"
	"github.com/stretchr/testify/require"
)

func TestHost(t *testing.T) {
	t.Run("IPv4", func(t *testing.T) {
		prefix := netip.MustParsePrefix("100.64.0.0/10")

		addr, err := cidr.Host(prefix, 1)
		require.NoError(t, err)

		require.Equal(t, "100.64.0.1", addr.String())
	})

	t.Run("IPv6", func(t *testing.T) {
		prefix := netip.MustParsePrefix("fd00::/48")

		addr, err := cidr.Host(prefix, 1)
		require.NoError(t, err)

		require.Equal(t, "fd00::1", addr.String())
	})
}
