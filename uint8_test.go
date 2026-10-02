// Copyright (c) 2020-2026 Uber Technologies, Inc.
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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUint8(t *testing.T) {
	atom := NewUint8(42)

	require.Equal(t, uint8(42), atom.Load(), "Load didn't work.")
	require.Equal(t, uint8(46), atom.Add(4), "Add didn't work.")
	require.Equal(t, uint8(44), atom.Sub(2), "Sub didn't work.")
	require.Equal(t, uint8(45), atom.Inc(), "Inc didn't work.")
	require.Equal(t, uint8(44), atom.Dec(), "Dec didn't work.")

	require.True(t, atom.CAS(44, 0), "CAS didn't report a swap.")
	require.Equal(t, uint8(0), atom.Load(), "CAS didn't set the correct value.")

	require.True(t, atom.CompareAndSwap(0, 1), "CompareAndSwap didn't report a swap.")
	require.Equal(t, uint8(1), atom.Load(), "CompareAndSwap didn't set the correct value.")

	require.Equal(t, uint8(1), atom.Swap(2), "Swap didn't return the old value.")
	require.Equal(t, uint8(2), atom.Load(), "Swap didn't set the correct value.")

	atom.Store(42)
	require.Equal(t, uint8(42), atom.Load(), "Store didn't set the correct value.")

	t.Run("JSON/Marshal", func(t *testing.T) {
		bytes, err := json.Marshal(atom)
		require.NoError(t, err, "json.Marshal errored unexpectedly.")
		require.Equal(t, []byte("42"), bytes, "json.Marshal encoded the wrong bytes.")
	})

	t.Run("JSON/Unmarshal", func(t *testing.T) {
		err := json.Unmarshal([]byte("40"), &atom)
		require.NoError(t, err, "json.Unmarshal errored unexpectedly.")
		require.Equal(t, uint8(40), atom.Load(), "json.Unmarshal didn't set the correct value.")
	})

	t.Run("JSON/Unmarshal/Error", func(t *testing.T) {
		err := json.Unmarshal([]byte(`"40"`), &atom)
		require.Error(t, err, "json.Unmarshal didn't error as expected.")
		assertErrorJSONUnmarshalType(t, err,
			"json.Unmarshal failed with unexpected error %v, want UnmarshalTypeError.", err)
	})

	t.Run("Text/MarshalUnmarshal", func(t *testing.T) {
		atom := NewUint8(42)
		bytes, err := atom.MarshalText()
		require.NoError(t, err)
		require.Equal(t, []byte("42"), bytes)

		var atom2 Uint8
		err = atom2.UnmarshalText([]byte("40"))
		require.NoError(t, err)
		require.Equal(t, uint8(40), atom2.Load())

		err = atom2.UnmarshalText([]byte("invalid"))
		require.Error(t, err)

		err = atom2.UnmarshalText([]byte("256"))
		require.Error(t, err)
	})

	t.Run("String", func(t *testing.T) {
		atom := NewUint8(math.MaxUint8)
		assert.Equal(t, "255", atom.String(),
			"String() returned an unexpected value.")
	})

	t.Run("Wrapping", func(t *testing.T) {
		atom := NewUint8(math.MaxUint8)
		require.Equal(t, uint8(0), atom.Add(1), "Add didn't wrap around.")
		require.Equal(t, uint8(math.MaxUint8), atom.Sub(1), "Sub didn't wrap around.")
		require.Equal(t, uint8(0), atom.Inc(), "Inc didn't wrap around.")
		require.Equal(t, uint8(math.MaxUint8), atom.Dec(), "Dec didn't wrap around.")
	})

	t.Run("Uninitialized", func(t *testing.T) {
		var zeroAtom Uint8
		require.Equal(t, uint8(0), zeroAtom.Load())
		require.True(t, zeroAtom.CompareAndSwap(0, 10))
		require.Equal(t, uint8(10), zeroAtom.Load())
	})
}
