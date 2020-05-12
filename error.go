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

// Error is an atomic type-safe wrapper for error values.
type Error struct{ v Value }

type storedError struct{ Value error }

func wrapError(v error) storedError {
	return storedError{v}
}

func unwrapError(v storedError) error {
	return v.Value
}

// NewError creates a new Error.
func NewError(v error) *Error {
	x := &Error{}
	if v != nil {
		x.Store(v)
	}
	return x
}

// Load atomically loads the wrapped error.
func (x *Error) Load() error {
	v := x.v.Load()
	if v == nil {
		return nil
	}
	return unwrapError(v.(storedError))
}

// Store atomically stores the passed error.
//
// NOTE: This will cause an allocation.
func (x *Error) Store(v error) {
	x.v.Store(wrapError(v))
}
