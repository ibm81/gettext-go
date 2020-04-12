// Copyright 2020 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The xgettext-go program extracts translatable strings from Go packages.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	fset := token.NewFileSet()
	pkgs, firstErr := parser.ParseDir(fset, "../../examples/hi", nil, 0)
	if firstErr != nil {
		log.Fatal(firstErr)
	}

	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.CallExpr:
					switch sel := x.Fun.(type) {
					case *ast.SelectorExpr:
						if sel.X.(*ast.Ident).Name == "gettext" && sel.Sel.Name == "Gettext" {
							fmt.Println("msgstr:", evalStringValue(x.Args[0]))
						}
					}
				}
				return true
			})
		}
	}
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
