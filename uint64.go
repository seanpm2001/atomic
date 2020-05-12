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
	"sync/atomic"
)

// Uint64 is an atomic wrapper around uint64.
type Uint64 struct{ v uint64 }

// NewUint64 creates a new Uint64.
func NewUint64(i uint64) *Uint64 {
	return &Uint64{i}
}

// Load atomically loads the wrapped value.
func (i *Uint64) Load() uint64 {
	return atomic.LoadUint64(&i.v)
}

// Add atomically adds to the wrapped uint64 and returns the new value.
func (i *Uint64) Add(n uint64) uint64 {
	return atomic.AddUint64(&i.v, n)
}

// Sub atomically subtracts from the wrapped uint64 and returns the new value.
func (i *Uint64) Sub(n uint64) uint64 {
	return atomic.AddUint64(&i.v, ^(n - 1))
}

// Inc atomically increments the wrapped uint64 and returns the new value.
func (i *Uint64) Inc() uint64 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped uint64 and returns the new value.
func (i *Uint64) Dec() uint64 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Uint64) CAS(old, new uint64) bool {
	return atomic.CompareAndSwapUint64(&i.v, old, new)
}

// Store atomically stores the passed value.
func (i *Uint64) Store(n uint64) {
	atomic.StoreUint64(&i.v, n)
}

// Swap atomically swaps the wrapped uint64 and returns the old value.
func (i *Uint64) Swap(n uint64) uint64 {
	return atomic.SwapUint64(&i.v, n)
}

// MarshalJSON encodes the wrapped uint64 into JSON.
func (i *Uint64) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Load())
}

// UnmarshalJSON decodes JSON into the wrapped uint64.
func (i *Uint64) UnmarshalJSON(b []byte) error {
	var v uint64
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	i.Store(v)
	return nil
}
