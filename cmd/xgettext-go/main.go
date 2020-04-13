// Copyright 2020 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The xgettext-go program extracts translatable strings from Go packages.
package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
)

func main() {
	fset := token.NewFileSet()
	astPackages, firstErr := parser.ParseDir(fset, "../../examples/hi", nil, 0)
	if firstErr != nil {
		log.Fatal(firstErr)
	}

	var astFiles []*ast.File
	for _, pkg := range astPackages {
		for _, f := range pkg.Files {
			astFiles = append(astFiles, f)
		}
	}

	// https://github.com/golang/go/issues/26504
	typesConfig := &types.Config{
		Importer:    importer.For("source", nil),
		FakeImportC: true,
	}
	typesInfo := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	typesPackage, err := typesConfig.Check("", fset, astFiles, typesInfo)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range astPackages {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.CallExpr:
					switch sel := x.Fun.(type) {
					case *ast.SelectorExpr:
						if isGettextPackage(fset, typesPackage, sel.X) {
							fmt.Println("msgstr:", evalStringValue(x.Args[0]))
						}
					}
				}
				return true
			})
		}
	}
}

func isGettextPackage(fset *token.FileSet, pkg *types.Package, node ast.Node) bool {
	inner := pkg.Scope().Innermost(node.Pos())
	if ident, ok := node.(*ast.Ident); ok {
		if _, obj := inner.LookupParent(ident.Name, node.Pos()); obj != nil {
			if pkgName, ok := obj.(*types.PkgName); ok {
				if pkg := pkgName.Imported(); pkg != nil {
					if pkg.Path() == "github.com/chai2010/gettext-go" {
						return true
					}
				}
			}
		}
	}
	return false
}

func evalStringValue(val interface{}) string {
	switch val.(type) {
	case *ast.BasicLit:
		return val.(*ast.BasicLit).Value
	case *ast.BinaryExpr:
		if val.(*ast.BinaryExpr).Op != token.ADD {
			return ""
		}
		left := evalStringValue(val.(*ast.BinaryExpr).X)
		right := evalStringValue(val.(*ast.BinaryExpr).Y)
		return left[0:len(left)-1] + right[1:len(right)]
	default:
		panic(fmt.Sprintf("unknown type: %v", val))
	}
}
