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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {
	atom := NewBool(false)
	require.False(t, atom.Toggle(), "Expected Toggle to return previous value.")
	require.True(t, atom.Toggle(), "Expected Toggle to return previous value.")
	require.False(t, atom.Toggle(), "Expected Toggle to return previous value.")
	require.True(t, atom.Load(), "Unexpected state after swap.")

	require.True(t, atom.CAS(true, true), "CAS should swap when old matches")
	require.True(t, atom.Load(), "CAS should have no effect")
	require.True(t, atom.CAS(true, false), "CAS should swap when old matches")
	require.False(t, atom.Load(), "CAS should have modified the value")
	require.False(t, atom.CAS(true, false), "CAS should fail on old mismatch")
	require.False(t, atom.Load(), "CAS should not have modified the value")

	atom.Store(false)
	require.False(t, atom.Load(), "Unexpected state after store.")

	prev := atom.Swap(false)
	require.False(t, prev, "Expected Swap to return previous value.")

	prev = atom.Swap(true)
	require.False(t, prev, "Expected Swap to return previous value.")

	t.Run("JSON/Marshal", func(t *testing.T) {
		atom.Store(true)
		bytes, err := json.Marshal(atom)
		require.NoError(t, err, "json.Marshal errored unexpectedly.")
		require.Equal(t, []byte("true"), bytes, "json.Marshal encoded the wrong bytes.")
	})

	t.Run("JSON/Unmarshal", func(t *testing.T) {
		err := json.Unmarshal([]byte("false"), &atom)
		require.NoError(t, err, "json.Unmarshal errored unexpectedly.")
		require.False(t, atom.Load(), "json.Unmarshal didn't set the correct value.")
	})

	t.Run("JSON/Unmarshal/Error", func(t *testing.T) {
		err := json.Unmarshal([]byte("42"), &atom)
		require.Error(t, err, "json.Unmarshal didn't error as expected.")
		assertErrorJSONUnmarshalType(t, err,
			"json.Unmarshal failed with unexpected error %v, want UnmarshalTypeError.", err)
	})

	t.Run("String", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			assert.Equal(t, "true", NewBool(true).String(),
				"String() returned an unexpected value.")
		})

		t.Run("false", func(t *testing.T) {
			var b Bool
			assert.Equal(t, "false", b.String(),
				"String() returned an unexpected value.")
		})
	})
}

func TestBool_InitializeDefaults(t *testing.T) {
	tests := []struct {
		msg     string
		newBool func() *Bool
	}{
		{
			msg: "Uninitialized",
			newBool: func() *Bool {
				var b Bool
				return &b
			},
		},
		{
			msg: "NewBool with default",
			newBool: func() *Bool {
				return NewBool(false)
			},
		},
		{
			msg: "Bool swapped with default",
			newBool: func() *Bool {
				b := NewBool(true)
				b.Swap(false)
				return b
			},
		},
		{
			msg: "Bool CAS'd with default",
			newBool: func() *Bool {
				b := NewBool(true)
				b.CompareAndSwap(true, false)
				return b
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			t.Run("MarshalJSON", func(t *testing.T) {
				b := tt.newBool()
				marshalled, err := b.MarshalJSON()
				require.NoError(t, err)
				assert.Equal(t, "false", string(marshalled))
			})

			t.Run("String", func(t *testing.T) {
				b := tt.newBool()
				assert.Equal(t, "false", b.String())
			})

			t.Run("CompareAndSwap", func(t *testing.T) {
				b := tt.newBool()
				require.True(t, b.CompareAndSwap(false, true))
				assert.Equal(t, true, b.Load())
			})

			t.Run("Swap", func(t *testing.T) {
				b := tt.newBool()
				assert.Equal(t, false, b.Swap(true))
			})
		})
	}
}
