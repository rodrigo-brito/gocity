package model

import (
	"fmt"
	"strings"

	"github.com/rodrigo-brito/gocity/utils"

	"github.com/rodrigo-brito/gocity/analyzer"
)

type NodeType string

const (
	StructType  NodeType = "STRUCT"
	FileType    NodeType = "FILE"
	PackageType NodeType = "PACKAGE"
)

type Tree struct {
	Repository string  `json:"repository"`
	Children   []*Node `json:"children"`
}

type Node struct {
	Name               string   `json:"name"`
	Type               NodeType `json:"type"`
	Size               int      `json:"size"`
	Position           Position `json:"position"`
	NumberOfLine       int      `json:"numberOfLine"`
	NumberOfMethods    int      `json:"numberOfMethods"`
	NumberOfAttributes int      `json:"numberOfAttributes"`
	Children           []*Node  `json:"children"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func getPathAndFile(fullPath string) (paths []string, fileName, structName string) {
	pathlist := strings.Split(fullPath, "/")
	paths = pathlist[:len(pathlist)-1]
	fileName, structName = utils.GetFileAndStruct(pathlist[len(pathlist)-1])
	return
}

func New(tree map[string]*analyzer.NodeInfo, repositoryName string) Tree {
	value := Tree{
		Repository: repositoryName,
	}

	for key, _ := range tree {
		paths, fileName, structName := getPathAndFile(key)
		println(fileName, structName)
		fmt.Printf("%#v\n", paths)
	}

	return value
}
