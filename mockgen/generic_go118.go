// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build go1.18
// +build go1.18

package main

import (
	"go/ast"
	"strings"

	"github.com/golang/mock/mockgen/model"
)

func getTypeSpecTypeParams(ts *ast.TypeSpec) []*ast.Field {
	if ts == nil || ts.TypeParams == nil {
		return nil
	}
	return ts.TypeParams.List
}

func (p *fileParser) parseGenericType(pkg string, typ ast.Expr, tps map[string]bool) (model.Type, error) {
	switch v := typ.(type) {
	case *ast.IndexExpr:
		m, err := p.parseType(pkg, v.X, tps)
		if err != nil {
			return nil, err
		}
		nm, ok := m.(*model.NamedType)
		if !ok {
			return m, nil
		}
		nm.TypeParams = formatTypeArgs(v.Index.(*ast.Ident).Name)
		return m, nil
	case *ast.IndexListExpr:
		m, err := p.parseType(pkg, v.X, tps)
		if err != nil {
			return nil, err
		}
		nm, ok := m.(*model.NamedType)
		if !ok {
			return m, nil
		}
		var ts []string
		for _, expr := range v.Indices {
			ts = append(ts, expr.(*ast.Ident).Name)
		}
		nm.TypeParams = formatTypeArgs(ts...)
		return m, nil
	}
	return nil, nil
}

func getIdentTypeParams(decl interface{}) string {
	if decl == nil {
		return ""
	}
	ts, ok := decl.(*ast.TypeSpec)
	if !ok {
		return ""
	}
	if ts.TypeParams == nil || len(ts.TypeParams.List) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range ts.TypeParams.List {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.Names[0].Name)
	}
	sb.WriteString("]")
	return sb.String()
}

func formatTypeArgs(args ...string) string {
	if len(args) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range args {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v)
	}
	sb.WriteString("]")
	return sb.String()
}
