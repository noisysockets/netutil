//go:build windows
// +build windows

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
	"testing"
)

func TestReadlineWin(t *testing.T) {
	testCases := []readlineTestCase{
		{"foo\r\nbar\r\nbaz\r\n", []string{"foo", "bar", "baz"}},
		{"foo\r\nbar\r\nbaz", []string{"foo", "bar", "baz"}},
		{"foo\r\nbar\r\nbz\r\n\r\n", []string{"foo", "bar", "bz", ""}},
		{"foo\nbar\rbz\n\r", []string{"foo", "bar\rbz", ""}},
		{"foo\nbar\r\nbz\r\n\n", []string{"foo", "bar", "bz", ""}},
	}

	testReadline(t, testCases)
}
