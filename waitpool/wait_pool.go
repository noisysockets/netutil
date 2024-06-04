// SPDX-License-Identifier: MPL-2.0
/*
 * Copyright (C) 2024 The Noisy Sockets Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * Portions of this file are based on code originally from wireguard-go,
 *
 * Copyright (C) 2017-2023 WireGuard LLC. All Rights Reserved.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
 * of the Software, and to permit persons to whom the Software is furnished to do
 * so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package waitpool

import (
	"sync"
	"sync/atomic"
)

// WaitPool is a sync.Pool that blocks on Get when the maximum number of items
// are in use.
type WaitPool[T any] struct {
	pool  sync.Pool
	cond  sync.Cond
	lock  sync.Mutex
	count atomic.Int32
	max   uint32
}

// New creates a new WaitPool with a maximum size and a function to create new
// items. If max is 0, the pool has no maximum size and effectively becomes a
// typed sync.Pool.
func New[T any](max uint32, new func() T) *WaitPool[T] {
	p := &WaitPool[T]{pool: sync.Pool{New: func() any { return new() }}, max: max}
	p.cond = sync.Cond{L: &p.lock}
	return p
}

func (p *WaitPool[T]) Get() T {
	if p.max != 0 {
		p.lock.Lock()
		for uint32(p.count.Load()) >= p.max {
			p.cond.Wait()
		}
		p.count.Add(1)
		p.lock.Unlock()
	}
	return p.pool.Get().(T)
}

func (p *WaitPool[T]) Put(x T) {
	p.pool.Put(x)
	if p.max == 0 {
		return
	}
	p.count.Add(-1)
	p.cond.Signal()
}
