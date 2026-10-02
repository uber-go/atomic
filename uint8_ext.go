// Copyright (c) 2020-2023 Uber Technologies, Inc.
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

package atomic

import (
	"strconv"
)

//go:generate bin/gen-atomicwrapper -name=Uint8 -type=uint8 -wrapped=Uint32 -pack=packUint8 -unpack=unpackUint8 -cas -swap -json -file=uint8.go

func packUint8(n uint8) uint32 {
	return uint32(n)
}

func unpackUint8(n uint32) uint8 {
	return uint8(n)
}

// Add atomically adds to the wrapped uint8 and returns the new value.
func (i *Uint8) Add(delta uint8) uint8 {
	for {
		old := i.Load()
		new := old + delta
		if i.CompareAndSwap(old, new) {
			return new
		}
	}
}

// Sub atomically subtracts from the wrapped uint8 and returns the new value.
func (i *Uint8) Sub(delta uint8) uint8 {
	return i.Add(-delta)
}

// Inc atomically increments the wrapped uint8 and returns the new value.
func (i *Uint8) Inc() uint8 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped uint8 and returns the new value.
func (i *Uint8) Dec() uint8 {
	return i.Sub(1)
}

// String encodes the wrapped value as a string.
func (i *Uint8) String() string {
	return strconv.FormatUint(uint64(i.Load()), 10)
}

// MarshalText encodes the wrapped uint8 into a textual form.
func (i *Uint8) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(i.Load()), 10)), nil
}

// UnmarshalText decodes text into the wrapped uint8.
func (i *Uint8) UnmarshalText(b []byte) error {
	v, err := strconv.ParseUint(string(b), 10, 8)
	if err != nil {
		return err
	}
	i.Store(uint8(v))
	return nil
}
