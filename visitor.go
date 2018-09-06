package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

type info struct {
	Name             string
	NumberLines      int
	NumberFunctions  int
	NumberAttributes int
}

type Visitor struct {
	FileSet     *token.FileSet
	StructInfo  map[string]*info
	PackageName string
	Path        string
}

func (v Visitor) Print() {
	fmt.Println("Package: ", v.PackageName)
	for ID, v := range v.StructInfo {
		fmt.Println("ID: ", ID)
		fmt.Println("Source: ", v.Name)
		fmt.Println("Lines: ", v.NumberLines)
		fmt.Println("Attrs: ", v.NumberAttributes)
		fmt.Println("Funcs: ", v.NumberFunctions)
	}
}

func (v Visitor) getNumberOfLines(start, end token.Pos) int {
	return v.FileSet.Position(end).Line - v.FileSet.Position(start).Line + 1
}

func (v *Visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch d := n.(type) {
	case *ast.TypeSpec:
		if structObj, ok := d.Type.(*ast.StructType); ok {
			identifier := getIdentifier(v.PackageName, d.Name.Name)

			if _, ok := v.StructInfo[identifier]; !ok {
				v.StructInfo[identifier] = new(info)
			}

			v.StructInfo[identifier].Name = d.Name.Name
			v.StructInfo[identifier].NumberAttributes = len(structObj.Fields.List)
			v.StructInfo[identifier].NumberLines += v.getNumberOfLines(structObj.Pos(), structObj.End())
		}
	case *ast.FuncDecl:
		var structName = "(Orphan)"
		if d.Recv != nil && len(d.Recv.List) > 0 {
			typeObj := d.Recv.List[0].Type
			structName = typeObj.(*ast.StarExpr).X.(*ast.Ident).Name
		}

		identifier := getIdentifier(v.PackageName, structName)

		if _, ok := v.StructInfo[identifier]; !ok {
			v.StructInfo[identifier] = new(info)
			v.StructInfo[identifier].Name = structName
		}

		v.StructInfo[identifier].NumberFunctions += 1
		v.StructInfo[identifier].NumberLines += v.getNumberOfLines(d.Body.Pos(), d.Body.End())
	}

	return v
}
