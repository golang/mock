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
	"encoding/gob"
	"errors"
	"flag"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/golang/mock/mockgen/model"
)

const (
	pkgOutputFileName = "gomock_pkg_content_output"
)

var (
	progOnly   = flag.Bool("prog_only", false, "(reflect mode) Only generate the reflection program; write it to stdout and exit.")
	execOnly   = flag.String("exec_only", "", "(reflect mode) If set, execute this reflection program.")
	buildFlags = flag.String("build_flags", "", "(reflect mode) Additional flags for go build.")
)

func writeProgram(importPath string, symbols []string) ([]byte, error) {
	var program bytes.Buffer
	data := reflectData{
		ImportPath:               importPath,
		Symbols:                  symbols,
		PkgContentOutputFileName: pkgOutputFileName,
	}
	if err := reflectProgram.Execute(&program, &data); err != nil {
		return nil, err
	}
	return program.Bytes(), nil
}

// run the given command and parse the output as a model.Package.
func run(command string) (*model.Package, error) {
	// Run the program.
	cmd := exec.Command(command)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Process output.
	progDir, _ := path.Split(command)
	outputFile := filepath.Join(progDir, pkgOutputFileName)

	pkgBytes, err := ioutil.ReadFile(outputFile)
	// remove output file explicitly in case temporary output file is not in a tmp dir.
	defer func() { os.Remove(outputFile) }()
	if err != nil {
		return nil, errors.New("read pkg temporary output file failed")
	}

	reader := bytes.NewReader(pkgBytes)
	var pkg model.Package
	if err := gob.NewDecoder(reader).Decode(&pkg); err != nil {
		return nil, err
	}
	return &pkg, nil
}

// runInDir writes the given program into the given dir, runs it there, and
// parses the output as a model.Package.
func runInDir(program []byte, dir string) (*model.Package, error) {
	// We use TempDir instead of TempFile so we can control the filename.
	tmpDir, err := ioutil.TempDir(dir, "gomock_reflect_")
	if err != nil {
		return nil, err
	}
	defer func() { os.RemoveAll(tmpDir) }()
	const progSource = "prog.go"
	var progBinary = "prog.bin"
	if runtime.GOOS == "windows" {
		// Windows won't execute a program unless it has a ".exe" suffix.
		progBinary += ".exe"
	}

	if err := ioutil.WriteFile(filepath.Join(tmpDir, progSource), program, 0600); err != nil {
		return nil, err
	}

	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, "build")
	if *buildFlags != "" {
		cmdArgs = append(cmdArgs, *buildFlags)
	}
	cmdArgs = append(cmdArgs, "-o", progBinary, progSource)

	// Build the program.
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return run(filepath.Join(tmpDir, progBinary))
}

func Reflect(importPath string, symbols []string) (*model.Package, error) {
	// TODO: sanity check arguments

	if *execOnly != "" {
		return run(*execOnly)
	}

	program, err := writeProgram(importPath, symbols)
	if err != nil {
		return nil, err
	}

	if *progOnly {
		os.Stdout.Write(program)
		os.Exit(0)
	}

	wd, _ := os.Getwd()

	// Try to run the program in the same directory as the input package.
	if p, err := build.Import(importPath, wd, build.FindOnly); err == nil {
		dir := p.Dir
		if p, err := runInDir(program, dir); err == nil {
			return p, nil
		}
	}

	// Since that didn't work, try to run it in the current working directory.
	if p, err := runInDir(program, wd); err == nil {
		return p, nil
	}
	// Since that didn't work, try to run it in a standard temp directory.
	return runInDir(program, "")
}

type reflectData struct {
	ImportPath               string
	Symbols                  []string
	PkgContentOutputFileName string
}

// This program reflects on an interface value, and prints the
// gob encoding of a model.Package to standard output.
// JSON doesn't work because of the model.Type interface.
var reflectProgram = template.Must(template.New("program").Parse(`
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"

	"github.com/golang/mock/mockgen/model"

	pkg_ {{printf "%q" .ImportPath}}
)

func getProgDir() (string, error) {
	progBinary, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}

	progDir, _ := path.Split(progBinary)
	return progDir, nil
}

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

	var pkgContent bytes.Buffer
	if err := gob.NewEncoder(&pkgContent).Encode(pkg); err != nil {
		fmt.Fprintf(os.Stderr, "gob encode: %v\n", err)
		os.Exit(1)
	}

	progDir, err := getProgDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "get program path: %v\n", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(filepath.Join(progDir, {{printf "%q" .PkgContentOutputFileName}}), pkgContent.Bytes(), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "write pkg content to temporary output file: %v\n", err)
		os.Exit(1)
	}
}
`))
