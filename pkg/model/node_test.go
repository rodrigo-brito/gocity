package model

import (
	"testing"

	"github.com/rodrigo-brito/gocity/pkg/analyzer"
	"github.com/stretchr/testify/assert"
)

func TestNode_GenerateChildList(t *testing.T) {
	n := Node{
		childrenMap: map[string]*Node{
			"1": {
				Type: StructType,
				Line: 1,
			},
			"2": {
				Name: "name",
				Type: PackageType,
			},
			"3": {
				Type: FileType,
				childrenMap: map[string]*Node{
					"4": {
						Type: FileType,
					},
				},
			},
		},
	}

	n.GenerateChildList("https://github.com/rodrigo-brito/gocity/blob/master")
	assert.Equal(t, len(n.Children), 3)
	assert.Equal(t, n.childrenMap["1"].URL, "https://github.com/rodrigo-brito/gocity/blob/master#L1")
	assert.Equal(t, n.childrenMap["2"].URL, "https://github.com/rodrigo-brito/gocity/blob/master/name")
	assert.Equal(t, n.childrenMap["3"].URL, "https://github.com/rodrigo-brito/gocity/blob/master")
	assert.Equal(t, n.childrenMap["3"].childrenMap["4"].URL, "https://github.com/rodrigo-brito/gocity/blob/master")
}

func TestNode_GenerateChildrenPosition(t *testing.T) {
	n := Node{}
	n.GenerateChildrenPosition()
	assert.Equal(t, n.Width, float64(1))
	assert.Equal(t, n.Depth, float64(1))

	n.Children = append(n.Children, []*Node{
		{Width: 2, Depth: 2, Type: FileType, Children: []*Node{{Width: 16, Depth: 12}}},
		{Width: 10, Depth: 4},
	}...)
	n.GenerateChildrenPosition()
	assert.Equal(t, n.Children[0].Position.X, float64(0))
	assert.Equal(t, n.Children[0].Position.Y, float64(-1))
	assert.Equal(t, n.Children[0].Width, float64(2))
	assert.Equal(t, n.Children[0].Depth, float64(2))
	assert.Equal(t, n.Children[0].Children[0].Position.X, float64(0))
	assert.Equal(t, n.Children[0].Children[0].Position.Y, float64(0))
	assert.Equal(t, n.Children[1].Position.X, -0.5)
	assert.Equal(t, n.Children[1].Position.Y, 1.5)
	assert.Equal(t, n.Width, float64(3))
	assert.Equal(t, n.Depth, float64(5))
}

func TestNew(t *testing.T) {
	n := New(map[string]*analyzer.NodeInfo{"github.com/rodrigo-brito/gocity/blob/master/model/node.go.(Test)": {
		File:       "main.go",
		ObjectName: "main",
		Line:       1,
	}}, "gocity", "master")

	assert.Equal(t, n.Name, "gocity")
	assert.Equal(t, n.Branch, "master")
	assert.Equal(t, n.Depth, float64(9))
	assert.Equal(t, n.Width, float64(9))
	assert.Equal(t, n.Position.X, float64(0))
	assert.Equal(t, n.Position.Y, float64(0))
	assert.Equal(t, len(n.Children), 1)

	n = New(map[string]*analyzer.NodeInfo{"github.com/rodrigo-brito/gocity/blob/master/model/node.go": {
		File:       "main.go",
		ObjectName: "main",
		Line:       1,
	}}, "gocity", "master")

	assert.Equal(t, n.Name, "gocity")
	assert.Equal(t, n.Branch, "master")
	assert.Equal(t, n.Depth, float64(8))
	assert.Equal(t, n.Width, float64(8))
	assert.Equal(t, n.Position.X, float64(0))
	assert.Equal(t, n.Position.Y, float64(0))
	assert.Equal(t, len(n.Children), 1)
}
