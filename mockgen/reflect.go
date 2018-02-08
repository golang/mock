// Copyright 2012 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// This file contains the model construction by reflection.

import (
	"bytes"
	"flag"
	"io"
	"os"

	"github.com/golang/mock/mockgen/model"
	"github.com/golang/mock/mockgen/reflect"
)

var (
	progOnly   = flag.Bool("prog_only", false, "(reflect mode) Only generate the reflection program; write it to stdout.")
	execOnly   = flag.String("exec_only", "", "(reflect mode) If set, execute this reflection program.")
	buildFlags = flag.String("build_flags", "", "(reflect mode) Additional flags for go build.")
)

func Reflect(importPath string, symbols []string) (pkg *model.Package, err error) {
	// TODO: sanity check arguments

	progPath := *execOnly
	if progPath == "" {
		buf := &bytes.Buffer{}
		if err = reflect.Generate(buf, importPath, symbols); err != nil {
			return
		}
		if *progOnly {
			io.Copy(os.Stdout, buf)
			os.Exit(0)
		}
		var remover func()
		if progPath, remover, err = reflect.Write(buf.Bytes(), *buildFlags); err != nil {
			return
		}
		defer remover()
	}
	return reflect.Run(progPath)
}
