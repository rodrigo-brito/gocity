package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rodrigo-brito/gocity/pkg/analyzer"
	"github.com/rodrigo-brito/gocity/pkg/lib"
)

type NodeType string

const (
	StructType  NodeType = "STRUCT"
	FileType    NodeType = "FILE"
	PackageType NodeType = "PACKAGE"
)

type Node struct {
	Name               string   `json:"name"`
	URL                string   `json:"url"`
	Branch             string   `json:"branch"`
	Type               NodeType `json:"type"`
	Width              float64  `json:"width"`
	Depth              float64  `json:"depth"`
	Position           Position `json:"position"`
	NumberOfLines      int      `json:"numberOfLines"`
	NumberOfMethods    int      `json:"numberOfMethods"`
	NumberOfAttributes int      `json:"numberOfAttributes"`
	Children           []*Node  `json:"children"`
	Line               int      `json:"-"`
	childrenMap        map[string]*Node
}

const (
	BaseTypeFlag       = "{{TYPE}}"
	PackageBaseTypeURL = "tree"
	FileBaseTypeURL    = "blob"
)

func getNodeURL(node *Node, parentPath string) (raw string, formatted string) {
	if node.Type == StructType {
		formatted = fmt.Sprintf("%s#L%d", strings.Replace(parentPath, BaseTypeFlag, FileBaseTypeURL, -1), node.Line)
		return formatted, formatted
	}

	if len(node.Name) > 0 {
		raw = fmt.Sprintf("%s/%s", parentPath, node.Name)
	} else {
		raw = parentPath
	}

	if node.Type == PackageType {
		formatted = strings.Replace(raw, BaseTypeFlag, PackageBaseTypeURL, -1)
		return
	}
	formatted = strings.Replace(raw, BaseTypeFlag, FileBaseTypeURL, -1)
	return
}

func (n *Node) GenerateChildList(parentPath string) {
	for _, child := range n.childrenMap {
		baseName, nodeURL := getNodeURL(child, parentPath)
		child.URL = nodeURL
		n.Children = append(n.Children, child)
		if len(child.childrenMap) > 0 {
			child.GenerateChildList(baseName)
		}
	}

	// Sort by width
	sort.Sort(sort.Reverse(byWidth(n.Children)))
}

func (n *Node) GenerateChildrenPosition() {
	if len(n.Children) == 0 {
		n.Width = float64(n.NumberOfAttributes) + 1
		n.Depth = float64(n.NumberOfAttributes) + 1
		return
	}

	positionGenerator := NewGenerator(len(n.Children))
	for _, child := range n.Children {
		child.GenerateChildrenPosition()
		child.Position = positionGenerator.NextPosition(child.Width, child.Depth)
	}

	bounds := positionGenerator.GetBounds()
	n.Width, n.Depth = bounds.X, bounds.Y

	for _, child := range n.Children {
		child.Position.X -= n.Width / 2.0
		child.Position.Y -= n.Depth / 2.0
	}

	if n.Type == FileType {
		n.Width += float64(n.NumberOfAttributes)
		n.Depth += float64(n.NumberOfAttributes)
	}
}

func getPathAndFile(fullPath string) (paths []string, fileName, structName string) {
	pathlist := strings.Split(fullPath, "/")
	paths = pathlist[:len(pathlist)-1]
	fileName, structName = lib.GetFileAndStruct(pathlist[len(pathlist)-1])
	return
}

func New(items map[string]*analyzer.NodeInfo, repositoryName string, repositoryBranch string) *Node {
	tree := &Node{
		Name:        repositoryName,
		Branch:      repositoryBranch,
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
					Line:               value.Line,
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

	tree.GenerateChildList(fmt.Sprintf("https://%s/%s/%s", repositoryName, BaseTypeFlag, repositoryBranch))
	tree.GenerateChildrenPosition()

	return tree
}
