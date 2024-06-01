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
	"bufio"
	"io"
	"strings"
	"testing"
)

type readlineTestCase struct {
	in  string
	out []string
}

func testReadline(t *testing.T, testCases []readlineTestCase) {
	for _, tc := range testCases {
		var e error
		var l string

		r := bufio.NewReader(strings.NewReader(tc.in))
		i := 0

		for l, e = readline(r); e == nil; l, e = readline(r) {
			if i >= len(tc.out) {
				t.Fatalf("Too many lines received")
			}

			if tc.out[i] != l {
				t.Fatalf("Line does not match.\n"+
					"\tExpected: %q\n"+
					"\tActual  : %q\n",
					tc.out[i], l)
			}
			i += 1
		}

		if e != io.EOF {
			t.Fatalf("Expected EOF, but exception is %T.", e)
		}

		if i != len(tc.out) {
			t.Fatalf("Not enough lines received")
		}
	}
}

func TestReadline(t *testing.T) {
	testCases := []readlineTestCase{
		{"foo\nbar\nbaz\n", []string{"foo", "bar", "baz"}},
		{"foo\nbar\nbaz", []string{"foo", "bar", "baz"}},
		{"foo\nbar\nbaz\n\n", []string{"foo", "bar", "baz", ""}},
		{"foo\nbar\nbaz\n\nx", []string{"foo", "bar", "baz", "", "x"}},
	}

	testReadline(t, testCases)
}
