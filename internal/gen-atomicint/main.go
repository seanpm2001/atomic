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

// gen-atomicint generates an atomic wrapper around an integer type.
//
//  gen-atomicint -name Int32 -wrapped int32 -file out.go
//
// The generated wrapper will use the functions in the sync/atomic package
// named after the generated type.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"text/template"
)

func main() {
	log.SetFlags(0)
	if err := run(os.Args[1:]); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(args []string) error {
	var opts struct {
		Name     string
		Wrapped  string
		File     string
		Unsigned bool
	}

	flag := flag.NewFlagSet("gen-atomicint", flag.ContinueOnError)

	flag.StringVar(&opts.Name, "name", "", "name of the generated type (e.g. Int32)")
	flag.StringVar(&opts.Wrapped, "wrapped", "", "name of the wrapped type (e.g. int32)")
	flag.StringVar(&opts.File, "file", "", "output file path (default: stdout)")
	flag.BoolVar(&opts.Unsigned, "unsigned", false, "whether the type is unsigned")

	if err := flag.Parse(args); err != nil {
		return err
	}

	if len(opts.Name) == 0 || len(opts.Wrapped) == 0 {
		return errors.New("flags -name and -wrapped are required")
	}

	var w io.Writer = os.Stdout
	if file := opts.File; len(file) > 0 {
		f, err := os.Create(file)
		if err != nil {
			return fmt.Errorf("create %q: %v", file, err)
		}
		defer f.Close()

		w = f
	}

	data := struct {
		Name     string
		Wrapped  string
		Unsigned bool
	}{
		Name:     opts.Name,
		Wrapped:  opts.Wrapped,
		Unsigned: opts.Unsigned,
	}

	var buff bytes.Buffer
	if err := _tmpl.Execute(&buff, data); err != nil {
		return fmt.Errorf("render template: %v", err)
	}

	bs, err := format.Source(buff.Bytes())
	if err != nil {
		return fmt.Errorf("reformat source: %v", err)
	}

	_, err = w.Write(bs)
	return err
}

var _tmpl = template.Must(template.New("int.go").Parse(`// Copyright (c) 2020 Uber Technologies, Inc.
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

// {{ .Name }} is an atomic wrapper around {{ .Wrapped }}.
type {{ .Name }} struct{ v {{ .Wrapped }} }

// New{{ .Name }} creates a new {{ .Name }}.
func New{{ .Name }}(i {{ .Wrapped }}) *{{ .Name }} {
	return &{{ .Name }}{i}
}

// Load atomically loads the wrapped value.
func (i *{{ .Name }}) Load() {{ .Wrapped }} {
	return atomic.Load{{ .Name }}(&i.v)
}

// Add atomically adds to the wrapped {{ .Wrapped }} and returns the new value.
func (i *{{ .Name }}) Add(n {{ .Wrapped }}) {{ .Wrapped }} {
	return atomic.Add{{ .Name }}(&i.v, n)
}

// Sub atomically subtracts from the wrapped {{ .Wrapped }} and returns the new value.
func (i *{{ .Name }}) Sub(n {{ .Wrapped }}) {{ .Wrapped }} {
	return atomic.Add{{ .Name }}(&i.v,
		{{- if .Unsigned -}}
			^(n - 1)
		{{- else -}}
			-n
		{{- end -}}
	)
}

// Inc atomically increments the wrapped {{ .Wrapped }} and returns the new value.
func (i *{{ .Name }}) Inc() {{ .Wrapped }} {
	return i.Add(1)
}

// Dec atomically decrements the wrapped {{ .Wrapped }} and returns the new value.
func (i *{{ .Name }}) Dec() {{ .Wrapped }} {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *{{ .Name }}) CAS(old, new {{ .Wrapped }}) bool {
	return atomic.CompareAndSwap{{ .Name }}(&i.v, old, new)
}

// Store atomically stores the passed value.
func (i *{{ .Name }}) Store(n {{ .Wrapped }}) {
	atomic.Store{{ .Name }}(&i.v, n)
}

// Swap atomically swaps the wrapped {{ .Wrapped }} and returns the old value.
func (i *{{ .Name }}) Swap(n {{ .Wrapped }}) {{ .Wrapped }} {
	return atomic.Swap{{ .Name }}(&i.v, n)
}

// MarshalJSON encodes the wrapped {{ .Wrapped }} into JSON.
func (i *{{ .Name }}) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Load())
}

// UnmarshalJSON decodes JSON into the wrapped {{ .Wrapped }}.
func (i *{{ .Name }}) UnmarshalJSON(b []byte) error {
	var v {{ .Wrapped }}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	i.Store(v)
	return nil
}
`))
