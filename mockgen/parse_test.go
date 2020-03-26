package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
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
			dir:     "testdata/gomod/bar",
			pkgPath: "github.com/golang/foo/bar",
		},
		{
			name:    "go mod off",
			envs:    map[string]string{"GO111MODULE": "off", "GOPATH": "testdata/gopath"},
			dir:     "testdata/gopath/src/example.com/foo",
			pkgPath: "example.com/foo",
		},
		{
			name: "outside GOPATH",
			envs: map[string]string{"GO111MODULE": "off", "GOPATH": "testdata/gopath"},
			dir:  "testdata",
			err:  errOutsideGoPath,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			for key, value := range testCase.envs {
				os.Setenv(key, value)
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
	os.Setenv("GOPATH", goPath)
	os.Setenv("GO111MODULE", "on")
	pkgPath, err := parsePackageImport(srcDir)
	expected := "example.com/foo"
	if pkgPath != expected {
		t.Errorf("expect %s, got %s", expected, pkgPath)
	}
}
