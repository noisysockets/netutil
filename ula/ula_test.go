// SPDX-License-Identifier: MPL-2.0
/*
 * Copyright (C) 2024 The Noisy Sockets Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package ula_test

import (
	"testing"

	"github.com/noisysockets/netutil/ula"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	prefix, err := ula.Generate()
	require.NoError(t, err)

	require.True(t, prefix.Addr().IsGlobalUnicast())
}
