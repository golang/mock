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
				if genType, hasGeneric := p.getTypedParamForGeneric(pin.Type, embeddedIface, typeParams); hasGeneric {
					gm.In[pinIdx].Type = genType
				}
			}
			for outIdx, out := range gm.Out {
				if genType, hasGeneric := p.getTypedParamForGeneric(out.Type, embeddedIface, typeParams); hasGeneric {
					gm.Out[outIdx].Type = genType
				}
			}
			if gm.Variadic != nil {
				if vGenType, hasGeneric := p.getTypedParamForGeneric(gm.Variadic.Type, embeddedIface, typeParams); hasGeneric {
					gm.Variadic.Type = vGenType
				}
			}

			iface.AddMethod(gm)
		}
	}

	return
}

// getTypedParamForGeneric is recursive func to hydrate all generic types within a model.Type
// so they get populated instead with the actual desired target types
func (p *fileParser) getTypedParamForGeneric(t model.Type, iface *model.Interface, knownTypeParams []model.Type) (model.Type, bool) {
	switch typ := t.(type) {
	case *model.ArrayType:
		if gType, wasGeneric := p.getTypedParamForGeneric(typ.Type, iface, knownTypeParams); wasGeneric {
			typ.Type = gType
			return typ, true
		}
	case *model.ChanType:
		if gType, wasGeneric := p.getTypedParamForGeneric(typ.Type, iface, knownTypeParams); wasGeneric {
			typ.Type = gType
			return typ, true
		}
	case *model.FuncType:
		hasGeneric := false
		for inIdx, inParam := range typ.In {
			if genType, ok := p.getTypedParamForGeneric(inParam.Type, iface, knownTypeParams); ok {
				hasGeneric = true
				typ.In[inIdx].Type = genType
			}
		}
		for outIdx, outParam := range typ.Out {
			if genType, ok := p.getTypedParamForGeneric(outParam.Type, iface, knownTypeParams); ok {
				hasGeneric = true
				typ.Out[outIdx].Type = genType
			}
		}
		if typ.Variadic != nil {
			if genType, ok := p.getTypedParamForGeneric(typ.Variadic.Type, iface, knownTypeParams); ok {
				hasGeneric = true
				typ.Variadic.Type = genType
			}
		}
		if hasGeneric {
			return typ, true
		}
	case *model.MapType:
		var (
			keyTyp, valTyp               model.Type
			wasKeyGeneric, wasValGeneric bool
		)
		if keyTyp, wasKeyGeneric = p.getTypedParamForGeneric(typ.Key, iface, knownTypeParams); wasKeyGeneric {
			typ.Key = keyTyp
		}
		if valTyp, wasValGeneric = p.getTypedParamForGeneric(typ.Value, iface, knownTypeParams); wasValGeneric {
			typ.Value = valTyp
		}
		if wasKeyGeneric || wasValGeneric {
			return typ, true
		}
	case *model.NamedType:
		if typ.TypeParams == nil {
			return nil, false
		}
		hasGeneric := false
		for i, tp := range typ.TypeParams.TypeParameters {
			// it will either be a type with name matching a generic parameter
			// or it will be something like ptr or slice etc...
			if srcParamIdx := iface.TypeParamIndexByName(tp.String(nil, "")); srcParamIdx > -1 && srcParamIdx < len(knownTypeParams) {
				hasGeneric = true
				dstParamTyp := knownTypeParams[srcParamIdx]
				typ.TypeParams.TypeParameters[i] = dstParamTyp
			} else if _, ok := p.getTypedParamForGeneric(tp, iface, knownTypeParams); ok {
				hasGeneric = true
			}
		}
		if hasGeneric {
			return typ, true
		}
	case model.PredeclaredType:
		if srcParamIdx := iface.TypeParamIndexByName(typ.String(nil, "")); srcParamIdx > -1 {
			dstParamTyp := knownTypeParams[srcParamIdx]
			return dstParamTyp, true
		}
	case *model.PointerType:
		if gType, hasGeneric := p.getTypedParamForGeneric(typ.Type, iface, knownTypeParams); hasGeneric {
			typ.Type = gType
			return typ, true
		}
	}

	return nil, false
}
