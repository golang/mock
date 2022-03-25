package main

import (
	"fmt"
	"go/token"
	"go/types"
	"log"
	"os"

	"github.com/golang/mock/mockgen/model"

	"golang.org/x/tools/go/gcexportdata"
)

func archiveMode(importpath, archive string) (*model.Package, error) {
	f, err := os.Open(archive)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r, err := gcexportdata.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("read export data %q: %v", archive, err)
	}

	fset := token.NewFileSet()
	imports := make(map[string]*types.Package)
	tp, err := gcexportdata.Read(r, fset, imports, importpath)
	if err != nil {
		return nil, err
	}

	pkg := &model.Package{
		Name:    tp.Name(),
		PkgPath: tp.Path(),
	}
	for _, name := range tp.Scope().Names() {
		m := tp.Scope().Lookup(name)
		tn, ok := m.(*types.TypeName)
		if !ok {
			continue
		}
		ti, ok := tn.Type().Underlying().(*types.Interface)
		if !ok {
			continue
		}
		it, err := model.InterfaceFromGoTypesType(ti)
		if err != nil {
			log.Fatal(err)
		}
		it.Name = m.Name()
		pkg.Interfaces = append(pkg.Interfaces, it)
	}
	return pkg, nil
}
