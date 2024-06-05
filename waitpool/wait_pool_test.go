// SPDX-License-Identifier: MPL-2.0
/*
 * Copyright (C) 2024 The Noisy Sockets Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package waitpool_test

import (
	"testing"
	"time"

	"github.com/noisysockets/netutil/waitpool"
	"github.com/stretchr/testify/require"
)

func TestWaitPool(t *testing.T) {
	p := waitpool.New(10, func() []byte { return make([]byte, 512) })

	var bufs [10][]byte
	for i := 0; i < 10; i++ {
		bufs[i] = p.Get()
	}

	count := p.Count()
	require.Equal(t, 10, count)

	done := make(chan struct{})
	go func() {
		defer close(done)

		p.Get()
	}()

	// Should block.
	select {
	case <-done:
		t.Fatal("Get returned before Put")
	case <-time.After(10 * time.Millisecond):
	}

	// Return all the buffers to the pool.
	for i := 0; i < 10; i++ {
		p.Put(bufs[i])
	}

	// Should no longer block.
	select {
	case <-done:
	case <-time.After(10 * time.Millisecond):
		t.Fatal("Get did not return after Put")
	}

	// Get the buffer that was put back.
	buf := p.Get()
	require.Len(t, buf, 512)
}
