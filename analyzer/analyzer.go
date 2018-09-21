package analyzer

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/rodrigo-brito/gocity/utils"

	"github.com/rodrigo-brito/gocity/lib"
)

var ErrInvalidPackage = errors.New("invalid package")

type Analyzer interface {
	FetchPackage() error
	Analyze() (map[string]*NodeInfo, error)
}

type analyzer struct {
	PackageName string
}

func NewAnalyzer(packageName string) Analyzer {
	return &analyzer{
		PackageName: packageName,
	}
}

func (p *analyzer) FetchPackage() error {
	fetcher := lib.NewFetcher()
	ok, err := fetcher.Fetch(p.PackageName)
	if err != nil {
		return err
	}

	if !ok {
		return ErrInvalidPackage
	}

	return nil
}

func (a *analyzer) Analyze() (map[string]*NodeInfo, error) {
	summary := make(map[string]*NodeInfo)
	root := fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), a.PackageName)
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("error on file walk: %s", err)
		}

		fileSet := token.NewFileSet()
		if f.IsDir() || !utils.IsGoFile(f.Name()) {
			return nil
		}

		fmt.Printf("processing file %s...\n", path)

		file, err := parser.ParseFile(fileSet, path, nil, parser.AllErrors)
		if err != nil {
			log.Fatalf("invalid input %s: %s", path, err)
		}

		v := &Visitor{
			FileSet:     fileSet,
			PackageName: path,
			Path:        path,
			StructInfo:  summary,
		}

		ast.Walk(v, file)
		if err != nil {
			log.Fatalf("error on walk: %s", err)
			return err
		}

		return nil
	})

	return summary, err
}
