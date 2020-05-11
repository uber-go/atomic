// Copyright (c) 2020 Uber Technologies, Inc.
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
	"encoding/json"
	"sync/atomic"
)

// Int64 is an atomic wrapper around an int64.
type Int64 struct{ v int64 }

// NewInt64 creates an Int64.
func NewInt64(i int64) *Int64 {
	return &Int64{i}
}

// Load atomically loads the wrapped value.
func (i *Int64) Load() int64 {
	return atomic.LoadInt64(&i.v)
}

// Add atomically adds to the wrapped int64 and returns the new value.
func (i *Int64) Add(n int64) int64 {
	return atomic.AddInt64(&i.v, n)
}

// Sub atomically subtracts from the wrapped int64 and returns the new value.
func (i *Int64) Sub(n int64) int64 {
	return atomic.AddInt64(&i.v, -n)
}

// Inc atomically increments the wrapped int64 and returns the new value.
func (i *Int64) Inc() int64 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped int64 and returns the new value.
func (i *Int64) Dec() int64 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Int64) CAS(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&i.v, old, new)
}

// Store atomically stores the passed value.
func (i *Int64) Store(n int64) {
	atomic.StoreInt64(&i.v, n)
}

// Swap atomically swaps the wrapped int64 and returns the old value.
func (i *Int64) Swap(n int64) int64 {
	return atomic.SwapInt64(&i.v, n)
}

// MarshalJSON encodes the wrapped int64 into JSON.
func (i *Int64) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Load())
}

// UnmarshalJSON decodes JSON into the wrapped int64.
func (i *Int64) UnmarshalJSON(b []byte) error {
	var v int64
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	i.Store(v)
	return nil
}
