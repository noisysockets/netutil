// SPDX-License-Identifier: MPL-2.0
/*
 * Copyright (C) 2024 The Noisy Sockets Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package defaults_test

import (
	"testing"

	"github.com/noisysockets/netutil/defaults"
	"github.com/noisysockets/netutil/ptr"
	"github.com/stretchr/testify/require"
)

func TestWithDefaults(t *testing.T) {
	type config struct {
		A string
		B []int
		C *bool
	}

	defaultConf := config{
		A: "default",
		B: []int{1, 2, 3},
		C: ptr.To(true),
	}

	t.Run("Nil", func(t *testing.T) {
		conf, err := defaults.WithDefaults(nil, &defaultConf)
		require.NoError(t, err)

		require.Equal(t, defaultConf.A, conf.A)
		require.Equal(t, defaultConf.B, conf.B)
		require.Equal(t, *defaultConf.C, *conf.C)
	})

	t.Run("Empty", func(t *testing.T) {
		conf, err := defaults.WithDefaults(&config{}, &defaultConf)
		require.NoError(t, err)

		require.Equal(t, defaultConf.A, conf.A)
		require.Equal(t, defaultConf.B, conf.B)
		require.Equal(t, *defaultConf.C, *conf.C)
	})

	t.Run("Partial", func(t *testing.T) {
		conf, err := defaults.WithDefaults(&config{A: "partial", C: ptr.To(false)}, &defaultConf)
		require.NoError(t, err)

		require.Equal(t, "partial", conf.A)
		require.Equal(t, defaultConf.B, conf.B)
		require.False(t, *conf.C)
	})
}
