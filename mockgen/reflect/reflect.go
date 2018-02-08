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

package reflect

// This file contains the model construction by reflection.

import (
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/golang/mock/mockgen/model"
)

func Generate(program io.Writer, importPath string, symbols []string) (err error) {
	// Generate program.
	data := reflectData{
		ImportPath: importPath,
		Symbols:    symbols,
	}
	return reflectProgram.Execute(program, &data)
}

func Write(program []byte, flags string) (progPath string, remover func(), err error) {
	const progSource = "prog.go"
	var progBinary = "prog.bin"
	if runtime.GOOS == "windows" {
		// Windows won't execute a program unless it has a ".exe" suffix.
		progBinary += ".exe"
	}
	pwd, _ := os.Getwd()
	// We use TempDir instead of TempFile so we can control the filename.
	// Try to place the TempDir under pwd, so that if there is some package in
	// vendor directory, 'go build' can also load/mock it.
	tmpDir, err := ioutil.TempDir(pwd, "gomock_reflect_")
	if err != nil {
		return
	}
	remover = func() {
		os.RemoveAll(tmpDir)
	}
	defer func() {
		if err == nil {
			return
		}
		remover()
		remover = nil
	}()
	if err = ioutil.WriteFile(filepath.Join(tmpDir, progSource), program, 0600); err != nil {
		return
	}

	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, "build")
	if flags != "" {
		cmdArgs = append(cmdArgs, flags)
	}
	cmdArgs = append(cmdArgs, "-o", progBinary, progSource)

	// Build the program.
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return
	}
	progPath = filepath.Join(tmpDir, progBinary)
	return
}

func Run(progPath string) (*model.Package, error) {
	// Run it.
	cmd := exec.Command(progPath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Process output.
	var pkg model.Package
	if err := gob.NewDecoder(&stdout).Decode(&pkg); err != nil {
		return nil, err
	}
	return &pkg, nil
}

type reflectData struct {
	ImportPath string
	Symbols    []string
}

// This program reflects on an interface value, and prints the
// gob encoding of a model.Package to standard output.
// JSON doesn't work because of the model.Type interface.
var reflectProgram = template.Must(template.New("program").Parse(`
package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/golang/mock/mockgen/model"

	pkg_ {{printf "%q" .ImportPath}}
)

func main() {
	its := []struct{
		sym string
		typ reflect.Type
	}{
		{{range .Symbols}}
		{ {{printf "%q" .}}, reflect.TypeOf((*pkg_.{{.}})(nil)).Elem()},
		{{end}}
	}
	pkg := &model.Package{
		// NOTE: This behaves contrary to documented behaviour if the
		// package name is not the final component of the import path.
		// The reflect package doesn't expose the package name, though.
		Name: path.Base({{printf "%q" .ImportPath}}),
	}

	for _, it := range its {
		intf, err := model.InterfaceFromInterfaceType(it.typ)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Reflection: %v\n", err)
			os.Exit(1)
		}
		intf.Name = it.sym
		pkg.Interfaces = append(pkg.Interfaces, intf)
	}
	if err := gob.NewEncoder(os.Stdout).Encode(pkg); err != nil {
		fmt.Fprintf(os.Stderr, "gob encode: %v\n", err)
		os.Exit(1)
	}
}
`))
