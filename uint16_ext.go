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

//go:generate bin/gen-atomicwrapper -name=Uint16 -type=uint16 -wrapped=Uint32 -pack=packUint16 -unpack=unpackUint16 -cas -swap -json -file=uint16.go

func packUint16(n uint16) uint32 {
	return uint32(n)
}

func unpackUint16(n uint32) uint16 {
	return uint16(n)
}

// Add atomically adds to the wrapped uint16 and returns the new value.
func (i *Uint16) Add(delta uint16) uint16 {
	for {
		old := i.Load()
		new := old + delta
		if i.CompareAndSwap(old, new) {
			return new
		}
	}
}

// Sub atomically subtracts from the wrapped uint16 and returns the new value.
func (i *Uint16) Sub(delta uint16) uint16 {
	return i.Add(-delta)
}

// Inc atomically increments the wrapped uint16 and returns the new value.
func (i *Uint16) Inc() uint16 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped uint16 and returns the new value.
func (i *Uint16) Dec() uint16 {
	return i.Sub(1)
}

// String encodes the wrapped value as a string.
func (i *Uint16) String() string {
	return strconv.FormatUint(uint64(i.Load()), 10)
}

// MarshalText encodes the wrapped uint16 into a textual form.
func (i *Uint16) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(i.Load()), 10)), nil
}

// UnmarshalText decodes text into the wrapped uint16.
func (i *Uint16) UnmarshalText(b []byte) error {
	v, err := strconv.ParseUint(string(b), 10, 16)
	if err != nil {
		return err
	}
	i.Store(uint16(v))
	return nil
}
