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

// TODO: This does not support recursive embedded interfaces.
// TODO: This does not support embedding package-local interfaces in a separate file.

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
	gomockImportPath = "github.com/dsymonds/gomock/gomock"
)

var (
	source      = flag.String("source", "", "Input Go source file.")
	destination = flag.String("destination", "", "Output file; defaults to stdout.")
	packageOut  = flag.String("package", "", "Package of the generated code; defaults to the package of the input file with a 'mock_' prefix.")
	imports     = flag.String("imports", "", "Comma-separated name=path pairs of explicit imports to use.")
	auxFiles    = flag.String("aux_files", "", "Comma-separated pkg=path pairs of auxilliary Go source files.")
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

	if err := parseAuxFiles(*auxFiles); err != nil {
		log.Fatalf("Failed parsing auxilliary files: %v", err)
	}

	file, err := parser.ParseFile(token.NewFileSet(), *source, nil, 0)
	if err != nil {
		log.Fatalf("Failed parsing source file: %v", err)
	}
	addAuxInterfacesFromFile("", file)

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

var (
	auxFileList   []*ast.File
	auxInterfaces = make(map[string]map[string]*ast.InterfaceType)
)

func parseAuxFiles(auxFiles string) error {
	auxFiles = strings.TrimSpace(auxFiles)
	if auxFiles == "" {
		return nil
	}
	for _, kv := range strings.Split(auxFiles, ",") {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("bad aux file spec: %v", kv)
		}
		file, err := parser.ParseFile(token.NewFileSet(), parts[1], nil, 0)
		if err != nil {
			return err
		}
		auxFileList = append(auxFileList, file)
		addAuxInterfacesFromFile(parts[0], file)
	}
	return nil
}

func addAuxInterfacesFromFile(pkg string, file *ast.File) {
	if _, ok := auxInterfaces[pkg]; !ok {
		auxInterfaces[pkg] = make(map[string]*ast.InterfaceType)
	}
	for ni := range iterInterfaces(file) {
		auxInterfaces[pkg][ni.name.Name] = ni.it
	}
}

func auxInterface(pkg, name string) *ast.InterfaceType {
	m, ok := auxInterfaces[pkg]
	if !ok {
		log.Fatalf("Don't have any auxilliary interfaces for package %q", pkg)
	}
	it, ok := m[name]
	if !ok {
		log.Fatalf("Don't have an auxilliary interface %q in package %q", name, pkg)
	}
	return it
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
	case *ast.ChanType:
		return packagesOfType(v.Value)
	case *ast.Ellipsis:
		// a "..." type
		return packagesOfType(v.Elt)
	case *ast.FuncType:
		var pkgs []string
		if v.Params != nil {
			for _, f := range v.Params.List {
				pkgs = append(pkgs, packagesOfType(f.Type)...)
			}
		}
		if v.Results != nil {
			for _, f := range v.Results.List {
				pkgs = append(pkgs, packagesOfType(f.Type)...)
			}
		}
		return pkgs
	case *ast.Ident:
		// raw identifier
		return []string{}
	case *ast.InterfaceType:
		// TODO: Handle more than just interface{}
		return []string{}
	case *ast.MapType:
		return append(packagesOfType(v.Key), packagesOfType(v.Value)...)
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

// importsOfFile returns a map of package name to import path
// of the imports in file.
func importsOfFile(file *ast.File) map[string]string {
	/* We have to make guesses about some imports, because imports are not required
	 * to have names. Named imports are always certain. Unnamed imports are guessed
	 * to have a name of the last path component; if the last path component has dots,
	 * the first dot-delimited field is used as the name.
	 */

	m := make(map[string]string)
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
				if is.Name.Name == "_" {
					continue
				}
				pkg = removeDot(is.Name.Name)
			} else {
				_, last := path.Split(importPath)
				pkg = strings.SplitN(last, ".", 2)[0]
			}
			if _, ok := m[pkg]; ok {
				log.Fatalf("imported package collision: %q imported twice", pkg)
			}
			m[pkg] = importPath
		}
	}
	return m
}

// packagesUsedByInterface returns the package names used by an interface.
func packagesUsedByInterface(it *ast.InterfaceType) map[string]int {
	m := make(map[string]int)
	for _, method := range it.Methods.List {
		switch t := method.Type.(type) {
		case *ast.FuncType:
			for _, pkg := range packagesOfType(t) {
				m[pkg] = 1
			}
		case *ast.Ident:
			// Embedded interface in this package.
			for pkg := range packagesUsedByInterface(auxInterface("", t.Name)) {
				m[pkg] = 1
			}
		case *ast.SelectorExpr:
			// Embedded interface in another package.
			for pkg := range packagesUsedByInterface(auxInterface(t.X.(*ast.Ident).Name, t.Sel.Name)) {
				m[pkg] = 1
			}
		default:
			log.Fatalf("Don't know how to find packages used by a %T inside an interface", t)
		}
	}
	return m
}

func (g *generator) ScanImports(file *ast.File) {
	allImports := importsOfFile(file)

	// Include imports from auxilliary files, because we might need them for embedded interfaces.
	// Give the current file precedence over any conflicts.
	for _, f := range auxFileList {
		for pkg, path := range importsOfFile(f) {
			if _, ok := allImports[pkg]; !ok {
				allImports[pkg] = path
			}
		}
	}

	/* Now we scan the interfaces to look for uses of the packages. */

	usedPackages := make(map[string]int)
	for ni := range iterInterfaces(file) {
		for pkg := range packagesUsedByInterface(ni.it) {
			usedPackages[pkg] = 1
		}
	}

	/* Finally, gather the imports that are actually used */
	for pkg := range usedPackages {
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

func (g *generator) Generate(file *ast.File, pkg string) error {
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

func (g *generator) GenerateMockInterface(typeName *ast.Ident, it *ast.InterfaceType) error {
	mockType := mockName(typeName)

	g.p("")
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

	g.GenerateMockMethods(mockType, it, "")

	return nil
}

func (g *generator) GenerateMockMethods(mockType string, it *ast.InterfaceType, pkgOverride string) {
	for _, field := range it.Methods.List {
		switch ft := field.Type.(type) {
		case *ast.FuncType:
			if len(field.Names) != 1 {
				log.Fatal("unexpected case: there should be exactly one Ident for a method in an interface")
			}
			g.p("")
			g.GenerateMockMethod(mockType, field.Names[0].String(), field.Type.(*ast.FuncType), pkgOverride)
			g.p("")
			g.GenerateMockRecorderMethod(mockType, field.Names[0].String(), field.Type.(*ast.FuncType))
		case *ast.Ident:
			// Embedded interface in this package.
			g.GenerateMockMethodsForEmbedded(mockType, "", ft.Name)
		case *ast.SelectorExpr:
			// Embedded interface in another package.
			g.GenerateMockMethodsForEmbedded(mockType, ft.X.(*ast.Ident).Name, ft.Sel.Name)
		default:
			log.Fatalf("Don't know how to mock method of type %T", field.Type)
		}
	}
}

func (g *generator) GenerateMockMethodsForEmbedded(mockType, pkg, name string) {
	it := auxInterface(pkg, name)

	nicePkg := ""
	if pkg != "" {
		nicePkg = pkg + "."
	}
	g.p("")
	g.p("// Methods for embedded interface %s%s", nicePkg, name)
	g.GenerateMockMethods(mockType, it, pkg)
}

func typeString(f ast.Expr, pkgOverride string) string {
	switch v := f.(type) {
	case *ast.ArrayType:
		// slice or array
		if v.Len == nil {
			// slice
			return "[]" + typeString(v.Elt, pkgOverride)
		}
		if bl, ok := v.Len.(*ast.BasicLit); ok && bl.Kind == token.INT {
			// array
			return fmt.Sprintf("[%v]%s", bl.Value, typeString(v.Elt, pkgOverride))
		}
		log.Printf("WARNING: odd *ast.ArrayType: %v", v)
	case *ast.ChanType:
		var s string
		switch v.Dir {
		case ast.SEND:
			s = "chan<-"
		case ast.RECV:
			s = "<-chan"
		default:
			s = "chan"
		}
		return s + " " + typeString(v.Value, pkgOverride)
	case *ast.Ellipsis:
		return "..." + typeString(v.Elt, pkgOverride)
	case *ast.FuncType:
		inStr, outStr := flattenFieldList(v.Params, pkgOverride).typeString(), ""
		out := flattenFieldList(v.Results, pkgOverride)
		switch nOut := len(out.t); {
		case nOut == 1:
			outStr = " " + out.typeString()
		case nOut > 1:
			outStr = " (" + out.typeString() + ")"
		}
		return "func (" + inStr + ")" + outStr
	case *ast.Ident:
		// NOTE(dsymonds): This is an approximation, but a reasonable one.
		// It breaks if the foreign interface being embedded refers to a type T
		// that it knows about through a dot import. Dot imports are discouraged
		// anyway, so this is a reasonable heuristic.
		if pkgOverride != "" && ast.IsExported(v.Name) {
			return pkgOverride + "." + v.Name
		}
		return v.Name
	case *ast.InterfaceType:
		// TODO: Support more than just interface{}
		return "interface{}"
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", typeString(v.Key, pkgOverride), typeString(v.Value, pkgOverride))
	case *ast.SelectorExpr:
		// a foreign type.
		return fmt.Sprintf("%v.%v", v.X.(*ast.Ident), v.Sel)
	case *ast.StarExpr:
		// pointer
		return "*" + typeString(v.X, pkgOverride)
	}
	log.Printf("WARNING: failed to generate type string for %T", f)
	return fmt.Sprintf("<%T>", f)
}

// Represents a flattened list of parameters (either arguments or return values).
type parameterList struct {
	name, t []string // name and type of each arg (name may be empty)
}

func flattenFieldList(fields *ast.FieldList, pkgOverride string) *parameterList {
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
		ts := typeString(f.Type, pkgOverride)
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

// GenerateMockMethod generates a mock method implementation.
// If non-empty, pkgOverride is the package in which unqualified types reside.
func (g *generator) GenerateMockMethod(mockType, methodName string, f *ast.FuncType, pkgOverride string) error {
	args := flattenFieldList(f.Params, pkgOverride)
	rets := flattenFieldList(f.Results, pkgOverride)

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

func (g *generator) GenerateMockRecorderMethod(mockType, methodName string, f *ast.FuncType) error {
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
