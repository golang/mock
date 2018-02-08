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

// This file contains the model construction by parsing source files.

import (
	"flag"
	"strings"

	"github.com/golang/mock/mockgen/model"
	"github.com/golang/mock/mockgen/parser"
)

var (
	imports  = flag.String("imports", "", "(source mode) Comma-separated name=path pairs of explicit imports to use.")
	auxFiles = flag.String("aux_files", "", "(source mode) Comma-separated pkg=path pairs of auxiliary Go source files.")
)

// TODO: simplify error reporting

func ParseFile(source string) (*model.Package, error) {
	return parser.ParseFile(source, commaSeparatedArgs(imports), commaSeparatedArgs(auxFiles))
}

func commaSeparatedArgs(in *string) []string {
	s := strings.TrimSpace(*in)
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}
