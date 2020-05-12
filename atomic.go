// Copyright (c) 2016-2020 Uber Technologies, Inc.
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

// Package atomic provides simple wrappers around numerics to enforce atomic
// access.
package atomic

import (
	"encoding/json"
	"math"
	"sync/atomic"
	"time"
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

// Value shadows the type of the same name from sync/atomic
// https://godoc.org/sync/atomic#Value
type Value struct {
	nocmp // disallow non-atomic comparison
	atomic.Value
}
