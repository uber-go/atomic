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
	"strconv"
	"sync/atomic"
)

// Bool is an atomic Boolean.
type Bool struct {
	nocmp // disallow non-atomic comparison

	v uint32
}

// NewBool creates a Bool.
func NewBool(initial bool) *Bool {
	return &Bool{v: boolToInt(initial)}
}

// Load atomically loads the Boolean.
func (b *Bool) Load() bool {
	return truthy(atomic.LoadUint32(&b.v))
}

// CAS is an atomic compare-and-swap.
func (b *Bool) CAS(old, new bool) bool {
	return atomic.CompareAndSwapUint32(&b.v, boolToInt(old), boolToInt(new))
}

// Store atomically stores the passed value.
func (b *Bool) Store(new bool) {
	atomic.StoreUint32(&b.v, boolToInt(new))
}

// Swap sets the given value and returns the previous value.
func (b *Bool) Swap(new bool) bool {
	return truthy(atomic.SwapUint32(&b.v, boolToInt(new)))
}

// Toggle atomically negates the Boolean and returns the previous value.
func (b *Bool) Toggle() bool {
	for {
		old := b.Load()
		if b.CAS(old, !old) {
			return old
		}
	}
}

func truthy(n uint32) bool {
	return n == 1
}

func boolToInt(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

// MarshalJSON encodes the wrapped bool into JSON.
func (b *Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Load())
}

// UnmarshalJSON decodes JSON into the wrapped bool.
func (b *Bool) UnmarshalJSON(t []byte) error {
	var v bool
	if err := json.Unmarshal(t, &v); err != nil {
		return err
	}
	b.Store(v)
	return nil
}

// String encodes the wrapped value as a string.
func (b *Bool) String() string {
	return strconv.FormatBool(b.Load())
}
