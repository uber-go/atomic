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
		Name string
		Type string
		Zero string
		File string
	}

	flag := flag.NewFlagSet("gen-atomicint", flag.ContinueOnError)

	flag.StringVar(&opts.Name, "name", "", "name of the generated type (e.g. Int32)")
	flag.StringVar(&opts.Type, "type", "", "name of the wrapped type (e.g. int32)")
	flag.StringVar(&opts.Zero, "zero", "", "zero value of the wrapped type (e.g. nil)")
	flag.StringVar(&opts.File, "file", "", "output file path (default: stdout)")

	if err := flag.Parse(args); err != nil {
		return err
	}

	if len(opts.Name) == 0 || len(opts.Type) == 0 || len(opts.Zero) == 0 {
		return errors.New("flags -name, -type, and -zero are required")
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
		Type     string
		Zero     string
		Nillable bool
	}{
		Name:     opts.Name,
		Type:     opts.Type,
		Zero:     opts.Zero,
		Nillable: opts.Zero == "nil",
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

// {{ .Name }} is an atomic type-safe wrapper for {{ .Type }} values.
type {{ .Name }} struct{ v Value }

{{/* atomic.Value panics for nil values. Generate a wrapper for these types. */}}
{{ $stored := .Type }}
{{ $wrap := .Type }}
{{ $unwrap := .Type }}
{{ if .Nillable -}}
	{{ $stored = printf "stored%s" .Name }}
	{{ $wrap = printf "wrap%s" .Name }}
	{{ $unwrap = printf "unwrap%s" .Name }}

	type {{ $stored }} struct{ Value {{ .Type }} }

	func {{ $wrap }}(v {{ .Type }}) {{ $stored }} {
		return {{ $stored }}{v}
	}

	func {{ $unwrap }}(v {{ $stored }}) {{ .Type }} {
		return v.Value
	}
{{- end }}

// New{{ .Name }} creates a new {{ .Name }}.
func New{{ .Name }}(v {{ .Type }}) *{{ .Name }} {
	x := &{{ .Name }}{}
	if v != {{ .Zero }} {
		x.Store(v)
	}
	return x
}

// Load atomically loads the wrapped {{ .Type }}.
func (x *{{ .Name }}) Load() {{ .Type }} {
	v := x.v.Load()
	if v == nil {
		return {{ .Zero }}
	}
	return {{ $unwrap }}(v.({{ $stored }}))
}

// Store atomically stores the passed {{ .Type }}.
//
// NOTE: This will cause an allocation.
func (x *{{ .Name }}) Store(v {{ .Type }}) {
	x.v.Store({{ $wrap }}(v))
}
`))