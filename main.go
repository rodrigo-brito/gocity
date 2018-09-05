package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func getIdentifier(pkg, name string) string {
	return fmt.Sprintf("%s.%s", pkg, name)
}

func (v visitor) getNumberOfLines(start, end token.Pos) int {
	return v.FileSet.Position(end).Line - v.FileSet.Position(start).Line + 1
}

func (v *visitor) Visit(n ast.Node) ast.Visitor {
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

func main() {
	packageName := "example"
	fileSet := token.NewFileSet()
	packages, err := parser.ParseDir(fileSet, packageName, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("invalid input %s: %s", packageName, err)
	}

	for name, pkg := range packages {
		visitor := &visitor{FileSet: fileSet, PackageName: pkg.Name, StructInfo: make(map[string]*info)}
		fmt.Println("reading pack ", name)
		ast.Walk(visitor, pkg)
		if err != nil {
			log.Fatalf("error on walk: %s", err)
		}

		visitor.Print()
	}
}

type info struct {
	Name             string
	NumberLines      int
	NumberFunctions  int
	NumberAttributes int
}

type visitor struct {
	FileSet     *token.FileSet
	PackageName string
	StructInfo  map[string]*info
}

func (v visitor) Print() {
	fmt.Println("Package: ", v.PackageName)
	for _, v := range v.StructInfo {
		fmt.Println("Source: ", v.Name)
		fmt.Println("Lines: ", v.NumberLines)
		fmt.Println("Attrs: ", v.NumberAttributes)
		fmt.Println("Funcs: ", v.NumberFunctions)
	}
}
