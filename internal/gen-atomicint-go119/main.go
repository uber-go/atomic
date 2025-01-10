// Copyright (c) 2020-2022 Uber Technologies, Inc.
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
//	gen-atomicint -name Int32 -wrapped int32 -file out.go
//
// The generated wrapper will use the functions in the sync/atomic package
// named after the generated type.
package main

import (
	"bytes"
	"embed"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"text/template"
	"time"
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

	flag := flag.NewFlagSet("gen-atomicint-go119", flag.ContinueOnError)

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
		ToYear   int
	}{
		Name:     opts.Name,
		Wrapped:  opts.Wrapped,
		Unsigned: opts.Unsigned,
		ToYear:   time.Now().Year(),
	}

	var buff bytes.Buffer
	if err := _tmpl.ExecuteTemplate(&buff, "wrapper.tmpl", data); err != nil {
		return fmt.Errorf("render template: %v", err)
	}

	bs, err := format.Source(buff.Bytes())
	if err != nil {
		return fmt.Errorf("reformat source: %v", err)
	}

	io.WriteString(w, "// @generated Code generated by gen-atomicint.\n\n")
	_, err = w.Write(bs)
	return err
}

var (
	//go:embed *.tmpl
	_tmplFS embed.FS

	_tmpl = template.Must(template.New("atomicint").ParseFS(_tmplFS, "*.tmpl"))
)
