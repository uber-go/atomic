// Copyright (c) 2026 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

//go:build 386 || arm || mips || mipsle

package atomic

import "unsafe"

// On these platforms int64 and uint64 are only 32-bit aligned, but the 64-bit
// sync/atomic primitives panic when given an address that is not 64-bit
// aligned. The wrappers must therefore stay 64-bit aligned even when a user
// embeds them in a larger struct.
//
// These assertions fail to compile if that guarantee is lost.
var (
	_ [unsafe.Alignof(Int64{}) - 8]struct{}
	_ [8 - unsafe.Alignof(Int64{})]struct{}

	_ [unsafe.Alignof(Uint64{}) - 8]struct{}
	_ [8 - unsafe.Alignof(Uint64{})]struct{}

	_ [unsafe.Alignof(Duration{}) - 8]struct{}
	_ [8 - unsafe.Alignof(Duration{})]struct{}

	_ [unsafe.Alignof(Float64{}) - 8]struct{}
	_ [8 - unsafe.Alignof(Float64{})]struct{}
)
