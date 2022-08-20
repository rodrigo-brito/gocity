package analyzer

import (
	"go/ast"
	"go/token"

	"github.com/rodrigo-brito/gocity/pkg/lib"
)

type NodeInfo struct {
	File             string
	ObjectName       string
	NumberLines      int
	NumberMethods    int
	NumberAttributes int
	Line             int
}

type Visitor struct {
	FileSet     *token.FileSet
	StructInfo  map[string]*NodeInfo
	PackageName string
	Path        string
	TmpFolder   string
}

func (v Visitor) getNumberOfLines(start, end token.Pos) int {
	return v.FileSet.Position(end).Line - v.FileSet.Position(start).Line + 1
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch definition := node.(type) {
	case *ast.ValueSpec: // Atributes
		identifier := lib.GetIdentifier(v.TmpFolder, v.Path, v.PackageName, "")

		if _, ok := v.StructInfo[identifier]; !ok {
			v.StructInfo[identifier] = new(NodeInfo)
		}

		v.StructInfo[identifier].NumberAttributes++
		v.StructInfo[identifier].Line = v.FileSet.Position(definition.Pos()).Line

	case *ast.TypeSpec: // Structs
		if structObj, ok := definition.Type.(*ast.StructType); ok {
			identifier := lib.GetIdentifier(v.TmpFolder, v.Path, v.PackageName, definition.Name.Name)

			if _, ok := v.StructInfo[identifier]; !ok {
				v.StructInfo[identifier] = new(NodeInfo)
			}

			v.StructInfo[identifier].ObjectName = definition.Name.Name
			v.StructInfo[identifier].NumberAttributes = len(structObj.Fields.List)
			v.StructInfo[identifier].NumberLines += v.getNumberOfLines(structObj.Pos(), structObj.End())
			v.StructInfo[identifier].Line = v.FileSet.Position(structObj.Pos()).Line
		}
	case *ast.FuncDecl: // Methods
		var structName = ""
		if definition.Recv != nil && len(definition.Recv.List) > 0 {
			typeObj := definition.Recv.List[0].Type
			if ident, ok := typeObj.(*ast.Ident); ok {
				structName = ident.Name
			} else {
				if ident, ok := typeObj.(*ast.StarExpr).X.(*ast.Ident); ok {
					structName = ident.Name
				} else if ident, ok := typeObj.(*ast.StarExpr).X.(*ast.SelectorExpr); ok {
					structName = ident.Sel.Name
				}
			}
		}

		identifier := lib.GetIdentifier(v.TmpFolder, v.Path, v.PackageName, structName)

		if _, ok := v.StructInfo[identifier]; !ok {
			v.StructInfo[identifier] = new(NodeInfo)
			v.StructInfo[identifier].ObjectName = structName
		}

		v.StructInfo[identifier].NumberMethods++
		if definition.Body != nil {
			v.StructInfo[identifier].NumberLines += v.getNumberOfLines(definition.Body.Pos(), definition.Body.End())
		}
	}

	return v
}
