//  Copyright (c) 2021 Uber Technologies, Inc.
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

import "time"

//go:generate bin/gen-atomicwrapper -name=Time -type=time.Time -wrapped=Value -pack=packTime -unpack=unpackTime -imports time -file=time.go

type packedTime struct{ Value time.Time }

func packTime(t time.Time) interface{} {
	return packedTime{t}
}

func unpackTime(v interface{}) time.Time {
	if t, ok := v.(packedTime); ok {
		return t.Value
	}
	return time.Time{}
}

// Add atomically adds to the wrapped time.Time and returns the new value.
func (t *Time) Add(d time.Duration) time.Time {
	return t.Load().Add(d)
}

// Sub atomically subtracts from the wrapped time.Duration and returns the new value.
func (t *Time) Sub(t2 time.Time) time.Duration {
	return t.Load().Sub(t2)
}

// Round atomically rounds the wrapped time.Time and returns the new value.
func (t *Time) Round(d time.Duration) time.Time {
	return t.Load().Round(d)
}
