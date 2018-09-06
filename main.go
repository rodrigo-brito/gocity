package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

func getIdentifier(pkg, name string) string {
	return fmt.Sprintf("%s.%s", pkg, name)
}

func main() {
	root := "example"

	summary := make(map[string]*info)

	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		fileSet := token.NewFileSet()
		if !f.IsDir() {
			return nil
		}

		packages, err := parser.ParseDir(fileSet, path, nil, parser.AllErrors)
		if err != nil {
			log.Fatalf("invalid input %s: %s", path, err)
		}

		for _, pkg := range packages {
			v := &Visitor{
				FileSet:     fileSet,
				PackageName: pkg.Name,
				Path:        path,
				StructInfo:  summary,
			}

			ast.Walk(v, pkg)
			if err != nil {
				log.Fatalf("error on walk: %s", err)
				return err
			}

			v.Print()
		}

		return nil
	})

	if err != nil {
		log.Fatalf("error on read directory %s", root)
	}
}
