// Copyright (c) 2016 Uber Technologies, Inc.
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
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringNoInitialValue(t *testing.T) {
	atom := &String{}
	require.Equal(t, "", atom.Load(), "Initial value should be blank string")
}

func TestString(t *testing.T) {
	atom := NewString("")
	require.Equal(t, "", atom.Load(), "Expected Load to return initialized value")

	atom.Store("abc")
	require.Equal(t, "abc", atom.Load(), "Unexpected value after Store")

	atom = NewString("bcd")
	require.Equal(t, "bcd", atom.Load(), "Expected Load to return initialized value")

	bytes, err := json.Marshal(atom)
	require.NoError(t, err, "json.Marshal errored unexpectedly.")
	require.Equal(t, []byte("\"bcd\""), bytes, "json.Marshal encoded the wrong bytes.")

	err = json.Unmarshal([]byte("\"abc\""), &atom)
	require.NoError(t, err, "json.Unmarshal errored unexpectedly.")
	require.Equal(t, "abc", atom.Load(), "json.Unmarshal didn't set the correct value.")

	err = json.Unmarshal([]byte("42"), &atom)
	require.Error(t, err, "json.Unmarshal didn't error as expected.")
	require.True(t, errors.As(err, new(*json.UnmarshalTypeError)),
		"json.Unmarshal failed with unexpected error %v, want UnmarshalTypeError.", err)
}
