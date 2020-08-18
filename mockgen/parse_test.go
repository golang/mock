package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileParser_ParseFile(t *testing.T) {
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "internal/tests/custom_package_name/greeter/greeter.go", nil, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	p := fileParser{
		fileSet:            fs,
		imports:            make(map[string]importedPackage),
		importedInterfaces: make(map[string]map[string]*ast.InterfaceType),
	}

	pkg, err := p.parseFile("", file)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checkGreeterImports(t, p.imports)

	expectedName := "greeter"
	if pkg.Name != expectedName {
		t.Fatalf("Expected name to be %v but got %v", expectedName, pkg.Name)
	}

	expectedInterfaceName := "InputMaker"
	if pkg.Interfaces[0].Name != expectedInterfaceName {
		t.Fatalf("Expected interface name to be %v but got %v", expectedInterfaceName, pkg.Interfaces[0].Name)
	}
}

func TestFileParser_ParsePackage(t *testing.T) {
	fs := token.NewFileSet()
	_, err := parser.ParseFile(fs, "internal/tests/custom_package_name/greeter/greeter.go", nil, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	p := fileParser{
		fileSet:            fs,
		imports:            make(map[string]importedPackage),
		importedInterfaces: make(map[string]map[string]*ast.InterfaceType),
	}

	newP, err := p.parsePackage("github.com/golang/mock/mockgen/internal/tests/custom_package_name/greeter")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	checkGreeterImports(t, newP.imports)
}

func TestImportsOfFile(t *testing.T) {
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "internal/tests/custom_package_name/greeter/greeter.go", nil, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	imports, _ := importsOfFile(file)
	checkGreeterImports(t, imports)
}

func checkGreeterImports(t *testing.T, imports map[string]importedPackage) {
	// check that imports have stdlib package "fmt"
	if fmtPackage, ok := imports["fmt"]; !ok {
		t.Errorf("Expected imports to have key \"fmt\"")
	} else {
		expectedFmtPackage := "fmt"
		if fmtPackage.Path() != expectedFmtPackage {
			t.Errorf("Expected fmt key to have value %s but got %s", expectedFmtPackage, fmtPackage.Path())
		}
	}

	// check that imports have package named "validator"
	if validatorPackage, ok := imports["validator"]; !ok {
		t.Errorf("Expected imports to have key \"fmt\"")
	} else {
		expectedValidatorPackage := "github.com/golang/mock/mockgen/internal/tests/custom_package_name/validator"
		if validatorPackage.Path() != expectedValidatorPackage {
			t.Errorf("Expected validator key to have value %s but got %s", expectedValidatorPackage, validatorPackage.Path())
		}
	}

	// check that imports have package named "client"
	if clientPackage, ok := imports["client"]; !ok {
		t.Errorf("Expected imports to have key \"client\"")
	} else {
		expectedClientPackage := "github.com/golang/mock/mockgen/internal/tests/custom_package_name/client/v1"
		if clientPackage.Path() != expectedClientPackage {
			t.Errorf("Expected client key to have value %s but got %s", expectedClientPackage, clientPackage.Path())
		}
	}

	// check that imports don't have package named "v1"
	if _, ok := imports["v1"]; ok {
		t.Errorf("Expected import not to have key \"v1\"")
	}
}

func Benchmark_parseFile(b *testing.B) {
	source := "internal/tests/performance/big_interface/big_interface.go"
	for n := 0; n < b.N; n++ {
		sourceMode(source)
	}
}

func TestParsePackageImport(t *testing.T) {
	testRoot, err := ioutil.TempDir("", "test_root")
	if err != nil {
		t.Fatal("cannot create tempdir")
	}
	defer func() {
		if err = os.RemoveAll(testRoot); err != nil {
			t.Errorf("cannot clean up tempdir at %s: %v", testRoot, err)
		}
	}()
	barDir := filepath.Join(testRoot, "gomod/bar")
	if err = os.MkdirAll(barDir, 0755); err != nil {
		t.Fatalf("error creating %s: %v", barDir, err)
	}
	if err = ioutil.WriteFile(filepath.Join(barDir, "bar.go"), []byte("package bar"), 0644); err != nil {
		t.Fatalf("error creating gomod/bar/bar.go: %v", err)
	}
	if err = ioutil.WriteFile(filepath.Join(testRoot, "gomod/go.mod"), []byte("module github.com/golang/foo"), 0644); err != nil {
		t.Fatalf("error creating gomod/go.mod: %v", err)
	}
	goPath := filepath.Join(testRoot, "gopath")
	for _, testCase := range []struct {
		name    string
		envs    map[string]string
		dir     string
		pkgPath string
		err     error
	}{
		{
			name:    "go mod default",
			envs:    map[string]string{"GO111MODULE": ""},
			dir:     barDir,
			pkgPath: "github.com/golang/foo/bar",
		},
		{
			name:    "go mod off",
			envs:    map[string]string{"GO111MODULE": "off", "GOPATH": goPath},
			dir:     filepath.Join(testRoot, "gopath/src/example.com/foo"),
			pkgPath: "example.com/foo",
		},
		{
			name: "outside GOPATH",
			envs: map[string]string{"GO111MODULE": "off", "GOPATH": goPath},
			dir:  "testdata",
			err:  errOutsideGoPath,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			for key, value := range testCase.envs {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("unable to set environment variable %q to %q: %v", key, value, err)
				}
			}
			pkgPath, err := parsePackageImport(filepath.Clean(testCase.dir))
			if err != testCase.err {
				t.Errorf("expect %v, got %v", testCase.err, err)
			}
			if pkgPath != testCase.pkgPath {
				t.Errorf("expect %s, got %s", testCase.pkgPath, pkgPath)
			}
		})
	}
}

func TestParsePackageImport_FallbackGoPath(t *testing.T) {
	goPath, err := ioutil.TempDir("", "gopath")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err = os.RemoveAll(goPath); err != nil {
			t.Error(err)
		}
	}()
	srcDir := filepath.Join(goPath, "src/example.com/foo")
	err = os.MkdirAll(srcDir, 0755)
	if err != nil {
		t.Error(err)
	}
	key := "GOPATH"
	value := goPath
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("unable to set environment variable %q to %q: %v", key, value, err)
	}
	key = "GO111MODULE"
	value = "on"
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("unable to set environment variable %q to %q: %v", key, value, err)
	}
	pkgPath, err := parsePackageImport(srcDir)
	expected := "example.com/foo"
	if pkgPath != expected {
		t.Errorf("expect %s, got %s", expected, pkgPath)
	}
}

func TestParsePackageImport_FallbackMultiGoPath(t *testing.T) {
	var goPathList []string

	// first gopath
	goPath, err := ioutil.TempDir("", "gopath1")
	if err != nil {
		t.Error(err)
	}
	goPathList = append(goPathList, goPath)
	defer func() {
		if err = os.RemoveAll(goPath); err != nil {
			t.Error(err)
		}
	}()
	srcDir := filepath.Join(goPath, "src/example.com/foo")
	err = os.MkdirAll(srcDir, 0755)
	if err != nil {
		t.Error(err)
	}

	// second gopath
	goPath, err = ioutil.TempDir("", "gopath2")
	if err != nil {
		t.Error(err)
	}
	goPathList = append(goPathList, goPath)
	defer func() {
		if err = os.RemoveAll(goPath); err != nil {
			t.Error(err)
		}
	}()

	goPaths := strings.Join(goPathList, string(os.PathListSeparator))
	key := "GOPATH"
	value := goPaths
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("unable to set environment variable %q to %q: %v", key, value, err)
	}
	key = "GO111MODULE"
	value = "on"
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("unable to set environment variable %q to %q: %v", key, value, err)
	}
	pkgPath, err := parsePackageImport(srcDir)
	expected := "example.com/foo"
	if pkgPath != expected {
		t.Errorf("expect %s, got %s", expected, pkgPath)
	}
}
