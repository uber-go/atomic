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

package atomic

// String is an atomic type-safe wrapper around Value for strings.
type String struct{ v Value }

// NewString creates a String.
func NewString(str string) *String {
	s := &String{}
	if str != "" {
		s.Store(str)
	}
	return s
}

// Load atomically loads the wrapped string.
func (s *String) Load() string {
	v := s.v.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}

// MarshalText encodes the wrapped string into a textual form.
//
// This makes it encodable as JSON, YAML, XML, and more.
func (s *String) MarshalText() ([]byte, error) {
	return []byte(s.Load()), nil
}

// UnmarshalText decodes text and replaces the wrapped string with it.
//
// This makes it decodable from JSON, YAML, XML, and more.
func (s *String) UnmarshalText(b []byte) error {
	s.Store(string(b))
	return nil
}

// Store atomically stores the passed string.
// Note: Converting the string to an interface{} to store in the Value
// requires an allocation.
func (s *String) Store(str string) {
	s.v.Store(str)
}
