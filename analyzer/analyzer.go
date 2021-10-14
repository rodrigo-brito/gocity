package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/rodrigo-brito/gocity/lib"

	log "github.com/sirupsen/logrus"
)

type Analyzer interface {
	FetchPackage() error
	Analyze() (map[string]*NodeInfo, error)
}

type analyzer struct {
	PackageName string
	BranchName  string
	IgnoreNodes []string
	fetcher     lib.Fetcher
}

type Option func(a *analyzer)

func NewAnalyzer(packageName string, branchName string, options ...Option) Analyzer {
	analyzer := &analyzer{
		PackageName: packageName,
		BranchName:  branchName,
		fetcher:     lib.NewFetcher(),
	}

	for _, option := range options {
		option(analyzer)
	}

	return analyzer
}

func WithIgnoreList(files ...string) Option {
	return func(a *analyzer) {
		a.IgnoreNodes = files
	}
}

func (p *analyzer) FetchPackage() error {
	return p.fetcher.Fetch(p.PackageName, p.BranchName)
}

func (p *analyzer) IsInvalidPath(path string) bool {
	for _, value := range p.IgnoreNodes {
		return strings.Contains(path, value)
	}

	return false
}

func (a *analyzer) Analyze() (map[string]*NodeInfo, error) {
	summary := make(map[string]*NodeInfo)
	root := fmt.Sprintf("%s/src/%s", os.Getenv("GOCITY_CACHE"), a.PackageName)
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error on file walk: %s", err)
		}

		fileSet := token.NewFileSet()
		if f.IsDir() || !lib.IsGoFile(f.Name()) || a.IsInvalidPath(path) {
			return nil
		}

		file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			log.WithField("file", path).Warn(err)
			return nil
		}

		v := &Visitor{
			FileSet:     fileSet,
			PackageName: a.PackageName,
			Path:        path,
			StructInfo:  summary,
		}

		ast.Walk(v, file)
		if err != nil {
			return fmt.Errorf("error on walk: %s", err)
		}

		return nil
	})

	return summary, err
}
