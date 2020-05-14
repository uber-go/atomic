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
	"math"
	"strconv"
	"sync/atomic"
)

// Float64 is an atomic wrapper around float64.
type Float64 struct {
	nocmp // disallow non-atomic comparison

	v uint64
}

// NewFloat64 creates a Float64.
func NewFloat64(f float64) *Float64 {
	return &Float64{v: math.Float64bits(f)}
}

// Load atomically loads the wrapped value.
func (f *Float64) Load() float64 {
	return math.Float64frombits(atomic.LoadUint64(&f.v))
}

// Store atomically stores the passed value.
func (f *Float64) Store(s float64) {
	atomic.StoreUint64(&f.v, math.Float64bits(s))
}

// Add atomically adds to the wrapped float64 and returns the new value.
func (f *Float64) Add(s float64) float64 {
	for {
		old := f.Load()
		new := old + s
		if f.CAS(old, new) {
			return new
		}
	}
}

// Sub atomically subtracts from the wrapped float64 and returns the new value.
func (f *Float64) Sub(s float64) float64 {
	return f.Add(-s)
}

// CAS is an atomic compare-and-swap.
func (f *Float64) CAS(old, new float64) bool {
	return atomic.CompareAndSwapUint64(&f.v, math.Float64bits(old), math.Float64bits(new))
}

// MarshalJSON encodes the wrapped float64 into JSON.
func (f *Float64) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Load())
}

// UnmarshalJSON decodes JSON into the wrapped float64.
func (f *Float64) UnmarshalJSON(b []byte) error {
	var v float64
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	f.Store(v)
	return nil
}

// String encodes the wrapped value as a string.
func (f *Float64) String() string {
	// 'g' is the behavior for floats with %v.
	return strconv.FormatFloat(f.Load(), 'g', -1, 64)
}
