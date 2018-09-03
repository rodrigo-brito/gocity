package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch d := n.(type) {
	case *ast.Package:
		fmt.Printf("PACK %s (%d file(s))\n", d.Name, len(d.Files))
	case *ast.File:
		fmt.Printf("FILE %s = %d-%d\n", d.Name.Name, v.FileSet.Position(d.Pos()).Line, v.FileSet.Position(d.End()).Line)
	case *ast.StructType:
		fmt.Printf("STRUCT %v = %d-%d\n", d.Struct, v.FileSet.Position(d.Pos()).Line, v.FileSet.Position(d.End()).Line)
	case *ast.FuncDecl:
		var structName string
		if d.Recv != nil && len(d.Recv.List) > 0 {
			typeObj := d.Recv.List[0].Type
			structName = typeObj.(*ast.StarExpr).X.(*ast.Ident).Name
		}

		fmt.Printf("(%s) FUNC %s = %d-%d\n", structName, d.Name.Name, v.FileSet.Position(d.Body.Pos()).Line, v.FileSet.Position(d.Body.End()).Line)
	}

	return v
}

func main() {
	packageName := "example"
	fileSet := token.NewFileSet()
	files, err := parser.ParseDir(fileSet, packageName, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("invalid input %s: %s", packageName, err)
	}

	for _, file := range files {
		visitor := &visitor{FileSet: fileSet}
		ast.Walk(visitor, file)
		if err != nil {
			log.Fatalf("error on walk: %s", err)
		}
	}
}

type visitor struct {
	FileSet         *token.FileSet
	NumberStructs   int
	NumberFunctions int
}
