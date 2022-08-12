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
	"fmt"
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
		t, err := p.parseType(pkg, v.Index, tps)
		if err != nil {
			return nil, err
		}
		nm.TypeParams = &model.TypeParametersType{TypeParameters: []model.Type{t}}
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
		var ts []model.Type
		for _, expr := range v.Indices {
			t, err := p.parseType(pkg, expr, tps)
			if err != nil {
				return nil, err
			}
			ts = append(ts, t)
		}
		nm.TypeParams = &model.TypeParametersType{TypeParameters: ts}
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

func (p *fileParser) parseEmbeddedGenericIface(iface *model.Interface, field *ast.Field, pkg string, tps map[string]bool) (wasGeneric bool, err error) {
	switch v := field.Type.(type) {
	case *ast.IndexExpr, *ast.IndexListExpr:
		wasGeneric = true
		// generic embedded interface
		// may or may not be external pkg
		// *ast.IndexExpr for embedded generic iface with single index e.g. DoSomething[T]
		// *ast.IndexListExpr for embedded generic iface with multiple indexes e.g. DoSomething[T, K]
		var (
			ident    *ast.Ident
			selIdent *ast.Ident // selector identity only used in external import
			// path       string
			typeParams []model.Type // normalize to slice whether IndexExpr or IndexListExpr to make it consistent to work with
		)
		if ie, ok := v.(*ast.IndexExpr); ok {
			if se, ok := ie.X.(*ast.SelectorExpr); ok {
				ident, selIdent = se.X.(*ast.Ident), se.Sel
			} else {
				ident = ie.X.(*ast.Ident)
			}
			var typParam model.Type
			if typParam, err = p.parseType(pkg, ie.Index, tps); err != nil {
				return
			}
			typeParams = append(typeParams, typParam)
		} else {
			ile := v.(*ast.IndexListExpr)
			if se, ok := ile.X.(*ast.SelectorExpr); ok {
				ident, selIdent = se.X.(*ast.Ident), se.Sel
			} else {
				ident = ile.X.(*ast.Ident)
			}
			var typParam model.Type
			for i := range ile.Indices {
				if typParam, err = p.parseType(pkg, ile.Indices[i], tps); err != nil {
					return
				}
				typeParams = append(typeParams, typParam)
			}
		}

		var (
			embeddedIface *model.Interface
		)

		if selIdent == nil {
			if embeddedIface, err = p.retrieveEmbeddedIfaceModel(pkg, ident.Name, ident.Pos(), false); err != nil {
				return
			}
		} else {
			filePkg, sel := ident.String(), selIdent.String()
			if embeddedIface, err = p.retrieveEmbeddedIfaceModel(filePkg, sel, ident.Pos(), true); err != nil {
				return
			}
		}

		// Copy the methods.
		// TODO: apply shadowing rules.
		for _, m := range embeddedIface.Methods {
			// non-trivial part - we have to match up the as-used type params with the as-defined
			//    defined as DoSomething[T any, K any]
			//    used as    DoSomething[somPkg.SomeType, int64]
			// meaning methods may be like in definition:
			//    Do(T) (K, error)
			// but need to be like this in implementation:
			//    Do(somePkg.SomeType) (int64, error)
			gm := m.Clone() // clone so we can change without changing source def

			// overwrite all typed params for incoming/outgoing params
			// to get the implementor-specified typing over the definition-specified typing

			for pinIdx, pin := range gm.In {
				switch t := pin.Type.(type) {
				case *model.NamedType:
					p.populateTypedParamsFromNamedType(t, embeddedIface, typeParams)
				case model.PredeclaredType:
					p.populateTypedParamsFromPredeclaredType(t, pinIdx, gm.In, embeddedIface, typeParams)
				case *model.PointerType:
					p.populateTypedParamsFromPointerType(t, embeddedIface, typeParams)
				}
			}
			for outIdx, out := range gm.Out {
				switch t := out.Type.(type) {
				case *model.NamedType:
					p.populateTypedParamsFromNamedType(t, embeddedIface, typeParams)
				case model.PredeclaredType:
					p.populateTypedParamsFromPredeclaredType(t, outIdx, gm.Out, embeddedIface, typeParams)
				case *model.PointerType:
					p.populateTypedParamsFromPointerType(t, embeddedIface, typeParams)
				}
			}

			iface.AddMethod(gm)
		}
	}

	return
}

func (p *fileParser) populateTypedParamsFromNamedType(nt *model.NamedType, iface *model.Interface, knownTypeParams []model.Type) {
	if nt.TypeParams == nil {
		return
	}

	for i, tp := range nt.TypeParams.TypeParameters {
		switch tpt := tp.(type) {
		case *model.PointerType:
			p.populateTypedParamsFromPointerType(tpt, iface, knownTypeParams)
		default:
			if srcParamIdx := iface.TypeParamIndexByName(tp.String(nil, "")); srcParamIdx > -1 && srcParamIdx < len(knownTypeParams) {
				dstParamTyp := knownTypeParams[srcParamIdx]
				nt.TypeParams.TypeParameters[i] = dstParamTyp
			}
		}
	}
}

func (p *fileParser) populateTypedParamsFromPredeclaredType(pt model.PredeclaredType, paramIdx int, inOrOutParams []*model.Parameter, iface *model.Interface, knownTypParams []model.Type) {
	if srcParamIdx := iface.TypeParamIndexByName(pt.String(nil, "")); srcParamIdx > -1 {
		dstParamTyp := knownTypParams[srcParamIdx]
		inOrOutParams[paramIdx] = &model.Parameter{
			Name: "",
			Type: dstParamTyp,
		}
	}
}

func (p *fileParser) populateTypedParamsFromPointerType(pt *model.PointerType, iface *model.Interface, knownTypeParams []model.Type) {
	switch t := pt.Type.(type) {
	case model.PredeclaredType:
		parms := make([]*model.Parameter, 1)
		p.populateTypedParamsFromPredeclaredType(t, 0, parms, iface, knownTypeParams)
		if parms[0] != nil {
			pt.Type = parms[0].Type
		}
	case *model.NamedType:
		p.populateTypedParamsFromNamedType(t, iface, knownTypeParams)
	default:
		fmt.Println("unhandled model PointerType")
	}
}
