// Copyright 2010 Google Inc.
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

// MockGen generates mock implementations of Go interfaces.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	gomockImportPath = "gomock.googlecode.com/hg/gomock"
)

var (
	source      = flag.String("source", "", "Input Go source file.")
	destination = flag.String("destination", "", "Output file; defaults to stdout.")
	packageOut  = flag.String("package", "", "Package of the generated code; defaults to the package of the input file with a 'mock_' prefix.")
	imports     = flag.String("imports", "", "Comma-separated name=path pairs of explicit imports to use.")
)

func main() {
	flag.Parse()

	if *source == "" {
		log.Fatalf("No source passed in --source flag")
	}

	dst := os.Stdout
	if len(*destination) > 0 {
		f, err := os.Create(*destination)
		if err != nil {
			log.Fatalf("Failed opening destination file: %v", err)
		}
		defer f.Close()
		dst = f
	}

	file, err := parser.ParseFile(token.NewFileSet(), *source, nil, 0)
	if err != nil {
		log.Fatalf("Failed parsing source file: %v", err)
	}

	pkg := *packageOut
	if pkg == "" {
		pkg = "mock_" + file.Name.Name
	}

	g := generator{
		w:                    dst,
		filename:             *source,
		imports:              make(map[string]string),
		explicitNamedImports: make(map[string]string),
	}
	if *imports != "" {
		g.SetImports(*imports)
	}
	g.ScanImports(file)
	if err := g.Generate(file, pkg); err != nil {
		log.Fatalf("Failed generating mock: %v", err)
	}
}

func printAst(node interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, nil, node); err != nil {
		log.Fatalf("Unexpected error in printAst: %v", err)
	}
	return buf.String()
}

type generator struct {
	w        io.Writer
	filename string
	indent   string

	imports map[string]string // map from package name to import path

	// Named imports passed via --imports; map from package name to import path
	explicitNamedImports map[string]string
	// "Dot" imports passed via --imports (i.e. 'import . "..."')
	explicitDotImports []string
}

func (g *generator) p(format string, args ...interface{}) {
	fmt.Fprintf(g.w, g.indent+format+"\n", args...)
}

func (g *generator) in() {
	g.indent += "\t"
}

func (g *generator) out() {
	if len(g.indent) > 0 {
		g.indent = g.indent[0 : len(g.indent)-1]
	}
}

type namedInterface struct {
	name *ast.Ident
	it   *ast.InterfaceType
}

// Create an iterator over all interfaces in file.
func iterInterfaces(file *ast.File) <-chan namedInterface {
	ch := make(chan namedInterface)
	go func() {
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				it, ok := ts.Type.(*ast.InterfaceType)
				if !ok {
					continue
				}
				ch <- namedInterface{ts.Name, it}
			}
		}
		close(ch)
	}()
	return ch
}

func removeDot(s string) string {
	if len(s) > 0 && s[len(s)-1] == '.' {
		return s[0 : len(s)-1]
	}
	return s
}

// Return the list of packages needed by this type. Not guaranteed to be unique.
func packagesOfType(t ast.Expr) []string {
	switch v := t.(type) {
	case *ast.ArrayType:
		// slice or array
		return packagesOfType(v.Elt)
	case *ast.Ellipsis:
		// a "..." type
		return packagesOfType(v.Elt)
	case *ast.Ident:
		// raw identifier
		return []string{}
	case *ast.InterfaceType:
		// TODO: Handle more than just interface{}
		return []string{}
	case *ast.SelectorExpr:
		return []string{v.X.(*ast.Ident).Name}
	case *ast.StarExpr:
		return packagesOfType(v.X)
	}
	log.Fatalf("Can't deduce package for a %T", t)
	return nil // never reached
}

func (g *generator) SetImports(imps string) {
	for _, kv := range strings.Split(imps, ",") {
		eq := strings.Index(kv, "=")
		k, v := kv[:eq], kv[eq+1:]
		if k == "." {
			// TODO: Catch dupes?
			g.explicitDotImports = append(g.explicitDotImports, v)
		} else {
			// TODO: Catch dupes?
			g.explicitNamedImports[k] = v
		}
	}
}

func (g *generator) ScanImports(file *ast.File) {
	/* We have to make guesses about some imports, because imports are not required
	 * to have names. Named imports are always certain. Unnamed imports are guessed
	 * to have a name of the last path component; if the last path component has dots,
	 * the first dot-delimited field is used as the name.
	 */

	allImports := make(map[string]string)
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.IMPORT {
			continue
		}
		for _, spec := range gd.Specs {
			is, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}
			pkg, importPath := "", string(is.Path.Value)
			importPath = importPath[1 : len(importPath)-1] // remove quotes
			if is.Name != nil {
				pkg = removeDot(is.Name.Name)
			} else {
				_, last := path.Split(importPath)
				pkg = strings.SplitN(last, ".", 2)[0]
			}
			if _, ok := allImports[pkg]; ok {
				log.Fatalf("imported package collision: %q imported twice", pkg)
			}
			allImports[pkg] = importPath
		}
	}

	/* Now we scan the interfaces to look for uses of the packages. */

	usedPackages := make(map[string]int)
	for ni := range iterInterfaces(file) {
		for _, method := range ni.it.Methods.List {
			ft := method.Type.(*ast.FuncType)
			if ft.Params != nil {
				for _, f := range ft.Params.List {
					for _, pkg := range packagesOfType(f.Type) {
						usedPackages[pkg] = 1
					}
				}
			}
			if ft.Results != nil {
				for _, f := range ft.Results.List {
					for _, pkg := range packagesOfType(f.Type) {
						usedPackages[pkg] = 1
					}
				}
			}
		}
	}

	/* Finally, gather the imports that are actually used */
	for pkg, _ := range usedPackages {
		// Explicit named imports take precedence.
		importPath, ok := g.explicitNamedImports[pkg]
		if !ok {
			importPath, ok = allImports[pkg]
		}
		if !ok {
			log.Fatalf("Package %q used but not imported!", pkg)
		}
		g.imports[pkg] = importPath
	}
}

func (g *generator) Generate(file *ast.File, pkg string) os.Error {
	g.p("// Automatically generated by MockGen. DO NOT EDIT!")
	g.p("// Source: %v", g.filename)
	g.p("")

	if _, ok := g.imports["gomock"]; ok {
		log.Fatalf("interface uses gomock package, so is not mockable.")
	}
	g.imports["gomock"] = gomockImportPath

	g.p("package %v", pkg)
	g.p("")
	g.p("import (")
	g.in()
	for pkg, path := range g.imports {
		g.p("%v %q", pkg, path)
	}
	for _, path := range g.explicitDotImports {
		g.p(". %q", path)
	}
	g.out()
	g.p(")")
	g.p("")

	for ni := range iterInterfaces(file) {
		if err := g.GenerateMockInterface(ni.name, ni.it); err != nil {
			return err
		}
	}

	return nil
}

// The name of the mock type to use for the given interface identifier.
func mockName(typeName *ast.Ident) string {
	return fmt.Sprintf("Mock%v", typeName)
}

func (g *generator) GenerateMockInterface(typeName *ast.Ident, it *ast.InterfaceType) os.Error {
	mockType := mockName(typeName)

	g.p("// Mock of %v interface", typeName)
	g.p("type %v struct {", mockType)
	g.in()
	g.p("ctrl     *gomock.Controller")
	g.p("recorder *_%vRecorder", mockType)
	g.out()
	g.p("}")
	g.p("")

	g.p("// Recorder for %v (not exported)", mockType)
	g.p("type _%vRecorder struct {", mockType)
	g.in()
	g.p("mock *%v", mockType)
	g.out()
	g.p("}")
	g.p("")

	// TODO: Re-enable this if we can import the interface reliably.
	//g.p("// Verify that the mock satisfies the interface at compile time.")
	//g.p("var _ %v = (*%v)(nil)", typeName, mockType)
	//g.p("")

	g.p("func New%v(ctrl *gomock.Controller) *%v {", mockType, mockType)
	g.in()
	g.p("mock := &%v{ctrl: ctrl}", mockType)
	g.p("mock.recorder = &_%vRecorder{mock}", mockType)
	g.p("return mock")
	g.out()
	g.p("}")
	g.p("")

	// XXX: possible name collision here if someone has EXPECT in their interface.
	g.p("func (m *%v) EXPECT() *_%vRecorder {", mockType, mockType)
	g.in()
	g.p("return m.recorder")
	g.out()
	g.p("}")

	for _, field := range it.Methods.List {
		if len(field.Names) != 1 {
			log.Fatal("unexpected case: there should be exactly one Ident for a method in an interface")
		}
		g.p("")
		g.GenerateMockMethod(mockType, field.Names[0].String(), field.Type.(*ast.FuncType))
		g.p("")
		g.GenerateMockRecorderMethod(mockType, field.Names[0].String(), field.Type.(*ast.FuncType))
	}

	return nil
}

func typeString(f ast.Expr) string {
	switch v := f.(type) {
	case *ast.ArrayType:
		// slice or array
		if v.Len == nil {
			// slice
			return "[]" + typeString(v.Elt)
		}
		if bl, ok := v.Len.(*ast.BasicLit); ok && bl.Kind == token.INT {
			// array
			return fmt.Sprintf("[%v]%s", bl.Value, typeString(v.Elt))
		}
		log.Printf("WARNING: odd *ast.ArrayType: %v", v)
	case *ast.Ellipsis:
		return "..." + typeString(v.Elt)
	case *ast.Ident:
		return v.Name
	case *ast.InterfaceType:
		// TODO: Support more than just interface{}
		return "interface{}"
	case *ast.SelectorExpr:
		// a foreign type.
		return fmt.Sprintf("%v.%v", v.X.(*ast.Ident), v.Sel)
	case *ast.StarExpr:
		// pointer
		return "*" + typeString(v.X)
	}
	log.Printf("WARNING: failed to generate type string for %T", f)
	return fmt.Sprintf("<%T>", f)
}

// Represents a flattened list of parameters (either arguments or return values).
type parameterList struct {
	name, t []string // name and type of each arg (name may be empty)
}

func flattenFieldList(fields *ast.FieldList) *parameterList {
	n := 0
	if fields != nil && fields.List != nil {
		n = fields.NumFields()
	}
	params := &parameterList{
		name: make([]string, n),
		t:    make([]string, n),
	}
	if n == 0 {
		return params
	}
	i := 0 // destination index
	for _, f := range fields.List {
		ts := typeString(f.Type)
		switch len(f.Names) {
		case 0:
			// anonymous arg; allocate a fake name
			params.name[i] = fmt.Sprintf("_param%d", i)
			params.t[i] = ts
			i++
		default:
			for _, name := range f.Names {
				params.name[i] = name.Name
				params.t[i] = ts
				i++
			}
		}
	}
	return params
}

// A string suitable for use as a method argument.
func (p *parameterList) argumentString() string {
	strs := make([]string, len(p.t))
	for i, t := range p.t {
		// XXX: doesn't handle anonymous params.
		strs[i] = fmt.Sprintf("%v %v", p.name[i], t)
	}
	return strings.Join(strs, ", ")
}

// A string just consisting of the types.
func (p *parameterList) typeString() string {
	return strings.Join(p.t, ", ")
}

func (g *generator) GenerateMockMethod(mockType, methodName string, f *ast.FuncType) os.Error {
	args := flattenFieldList(f.Params)
	rets := flattenFieldList(f.Results)

	retString := strings.Join(rets.t, ", ")
	if len(rets.t) > 1 {
		retString = "(" + retString + ")"
	}
	if retString != "" {
		retString = " " + retString
	}

	g.p("func (m *%v) %v(%v)%v {", mockType, methodName, args.argumentString(), retString)
	g.in()

	callArgs := strings.Join(args.name, ", ")
	if callArgs != "" {
		callArgs = ", " + callArgs
	}
	if f.Results == nil || len(f.Results.List) == 0 {
		g.p(`m.ctrl.Call(m, "%v"%v)`, methodName, callArgs)
	} else {
		g.p(`ret := m.ctrl.Call(m, "%v"%v)`, methodName, callArgs)

		// Go does not allow "naked" type assertions on nil values, so we use the two-value form here.
		// The value of that is either (x.(T), true) or (Z, false), where Z is the zero value for T.
		// Happily, this coincides with the semantics we want here.
		for i, t := range rets.t {
			g.p("ret%d, _ := ret[%d].(%s)", i, i, t)
		}

		retAsserts := make([]string, len(rets.t))
		for i, _ := range rets.t {
			retAsserts[i] = fmt.Sprintf("ret%d", i)
		}
		g.p("return " + strings.Join(retAsserts, ", "))
	}

	g.out()
	g.p("}")
	return nil
}

// isVariadic returns whether the function is variadic.
func isVariadic(f *ast.FuncType) bool {
	nargs := len(f.Params.List)
	if nargs == 0 {
		return false
	}
	_, ok := f.Params.List[nargs-1].Type.(*ast.Ellipsis)
	return ok
}

func (g *generator) GenerateMockRecorderMethod(mockType, methodName string, f *ast.FuncType) os.Error {
	nargs, variadic := f.Params.NumFields(), isVariadic(f)
	if variadic {
		nargs--
	}
	args := make([]string, nargs)
	for i := 0; i < nargs; i++ {
		args[i] = "arg" + strconv.Itoa(i)
	}
	argString := strings.Join(args, ", ")
	if nargs > 0 {
		argString += " interface{}"
	}
	if variadic {
		argString += fmt.Sprintf(", arg%d ...interface{}", nargs)
	}

	g.p("func (mr *_%vRecorder) %v(%v) *gomock.Call {", mockType, methodName, argString)
	g.in()

	callArgs := strings.Join(args, ", ")
	if nargs > 0 {
		callArgs = ", " + callArgs
	}
	if variadic {
		callArgs += fmt.Sprintf(", arg%d", nargs)
	}
	g.p(`return mr.mock.ctrl.RecordCall(mr.mock, "%v"%v)`, methodName, callArgs)

	g.out()
	g.p("}")
	return nil
}
