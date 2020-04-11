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
				var s string
				switch x := n.(type) {
				case *ast.BasicLit:
					s = x.Value
				case *ast.Ident:
					s = x.Name
				}
				if s != "" {
					fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
				}
				return true
			})
		}
	}
}
