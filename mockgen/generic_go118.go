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

	"github.com/golang/mock/mockgen/model"
)

func getTypeSpecTypeParams(ts *ast.TypeSpec) []*ast.Field {
	if ts == nil || ts.TypeParams == nil {
		return nil
	}
	return ts.TypeParams.List
}

func parseGenericType(typ ast.Expr) (model.Type, bool) {
	switch v := typ.(type) {
	case *ast.IndexExpr:
		_ = v
	case *ast.IndexListExpr:
	}
	return nil, false
}
