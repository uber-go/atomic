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
	"time"
)

// Duration is an atomic wrapper around time.Duration
// https://godoc.org/time#Duration
type Duration struct {
	nocmp // disallow non-atomic comparison

	v Int64
}

// NewDuration creates a Duration.
func NewDuration(d time.Duration) *Duration {
	return &Duration{v: *NewInt64(int64(d))}
}

// Load atomically loads the wrapped value.
func (d *Duration) Load() time.Duration {
	return time.Duration(d.v.Load())
}

// Store atomically stores the passed value.
func (d *Duration) Store(n time.Duration) {
	d.v.Store(int64(n))
}

// Add atomically adds to the wrapped time.Duration and returns the new value.
func (d *Duration) Add(n time.Duration) time.Duration {
	return time.Duration(d.v.Add(int64(n)))
}

// Sub atomically subtracts from the wrapped time.Duration and returns the new value.
func (d *Duration) Sub(n time.Duration) time.Duration {
	return time.Duration(d.v.Sub(int64(n)))
}

// Swap atomically swaps the wrapped time.Duration and returns the old value.
func (d *Duration) Swap(n time.Duration) time.Duration {
	return time.Duration(d.v.Swap(int64(n)))
}

// CAS is an atomic compare-and-swap.
func (d *Duration) CAS(old, new time.Duration) bool {
	return d.v.CAS(int64(old), int64(new))
}

// MarshalJSON encodes the wrapped time.Duration into JSON.
func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Load())
}

// UnmarshalJSON decodes JSON into the wrapped time.Duration.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v time.Duration
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	d.Store(v)
	return nil
}

// String encodes the wrapped value as a string.
func (d *Duration) String() string {
	return d.Load().String()
}
