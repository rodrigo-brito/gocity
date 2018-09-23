package model

import (
	"math"
	"strings"

	"github.com/rodrigo-brito/gocity/utils"

	"github.com/rodrigo-brito/gocity/analyzer"
)

type NodeType string

const (
	StructType  NodeType = "STRUCT"
	FileType    NodeType = "FILE"
	PackageType NodeType = "PACKAGE"

	PackageSizeMargin = 1
)

type Node struct {
	Name               string   `json:"name"`
	URL                string   `json:"url"`
	Type               NodeType `json:"type"`
	Size               float64  `json:"size"`
	Position           Position `json:"position"`
	NumberOfLines      int      `json:"numberOfLines"`
	NumberOfMethods    int      `json:"numberOfMethods"`
	NumberOfAttributes int      `json:"numberOfAttributes"`
	Children           []*Node  `json:"children"`
	childrenMap        map[string]*Node
}

func (n *Node) GetSize() float64 {
	numberChildren := len(n.childrenMap)
	if numberChildren > 0 {
		dimension := math.Ceil(math.Sqrt(float64(numberChildren)))

		var maxSize float64
		for _, child := range n.childrenMap {
			size := child.GetSize()
			if size > maxSize {
				maxSize = size
			}
		}

		maxSize += PackageSizeMargin
		n.Size = float64(dimension) * maxSize
		return n.Size
	}

	n.Size = float64(n.NumberOfAttributes) + 1
	return n.Size
}

func (n *Node) GenerateChildList() {
	for _, child := range n.childrenMap {
		n.Children = append(n.Children, child)
		if len(child.childrenMap) > 0 {
			child.GenerateChildList()
		}
	}
}

func (n *Node) GenerateChildrenPosition() {
	if len(n.childrenMap) == 0 {
		return
	}

	var maxSize float64
	for _, child := range n.childrenMap {
		if child.Size > maxSize {
			maxSize = child.Size
		}
	}

	positionGenerator := NewGenerator(len(n.childrenMap), maxSize)
	for _, child := range n.childrenMap {
		child.Position = positionGenerator.NextPosition()
		child.GenerateChildrenPosition()
	}
}

func getPathAndFile(fullPath string) (paths []string, fileName, structName string) {
	pathlist := strings.Split(fullPath, "/")
	paths = pathlist[:len(pathlist)-1]
	fileName, structName = utils.GetFileAndStruct(pathlist[len(pathlist)-1])
	return
}

func New(items map[string]*analyzer.NodeInfo, repositoryName string) *Node {
	tree := &Node{
		Name:        repositoryName,
		childrenMap: make(map[string]*Node),
		Children:    make([]*Node, 0),
	}

	for key, value := range items {
		currentNode := tree
		paths, fileName, structName := getPathAndFile(key)
		for _, path := range paths {
			_, ok := currentNode.childrenMap[path]
			if !ok {
				currentNode.childrenMap[path] = &Node{
					Name:        path,
					Type:        PackageType,
					childrenMap: make(map[string]*Node),
				}
			}
			currentNode = currentNode.childrenMap[path]
		}

		_, ok := currentNode.childrenMap[fileName]
		if !ok {
			currentNode.childrenMap[fileName] = &Node{
				Name:        fileName,
				Type:        FileType,
				childrenMap: make(map[string]*Node),
			}
		}

		fileNode := currentNode.childrenMap[fileName]

		if len(structName) > 0 {
			structNode, ok := fileNode.childrenMap[structName]
			if !ok {
				fileNode.childrenMap[structName] = &Node{
					Name:               structName,
					Type:               StructType,
					NumberOfAttributes: value.NumberAttributes,
					NumberOfMethods:    value.NumberMethods,
					NumberOfLines:      value.NumberLines,
				}
			} else {
				structNode.NumberOfAttributes += value.NumberAttributes
				structNode.NumberOfLines += value.NumberLines
				structNode.NumberOfMethods += value.NumberMethods
			}
		} else {
			fileNode.NumberOfAttributes += value.NumberAttributes
			fileNode.NumberOfLines += value.NumberLines
			fileNode.NumberOfMethods += value.NumberMethods
		}
	}

	tree.GetSize()
	tree.GenerateChildrenPosition()
	tree.GenerateChildList()

	return tree
}
