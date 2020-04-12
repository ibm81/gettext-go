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
	"reflect"
	"sort"
)

func main() {
	fset := token.NewFileSet()
	pkgs, firstErr := parser.ParseDir(fset, "../../examples/hi", nil, 0)
	if firstErr != nil {
		log.Fatal(firstErr)
	}

	var files []*ast.File
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			files = append(files, f)
		}
	}

	// https://github.com/golang/go/issues/26504
	config := types.Config{
		Importer:    importer.For("source", nil),
		FakeImportC: true,
	}
	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	pkg, err := config.Check("", fset, files, info)
	if err != nil {
		log.Fatal(err)
	}
	for _, node := range getSortedKeys(info.Uses) {
		fmt.Println("use:", node)
	}

	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			return true
		})
	}

	_ = pkg
}

func getSortedKeys(m interface{}) []ast.Node {
	mValue := reflect.ValueOf(m)
	nodes := make([]ast.Node, mValue.Len())

	keys := mValue.MapKeys()
	for i := range keys {
		nodes[i] = keys[i].Interface().(ast.Node)
	}

	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Pos() == nodes[j].Pos() {
			return nodes[i].End() < nodes[j].End()
		}
		return nodes[i].Pos() < nodes[j].Pos()
	})

	return nodes
}
