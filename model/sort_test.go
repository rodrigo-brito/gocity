package model

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeSortByWidth(t *testing.T) {
	nodes := []*Node{
		{
			Name:  "second",
			Width: 10.1,
		},
		{
			Name:  "third",
			Width: 10.2,
		},
		{
			Name:  "first",
			Width: 1,
		},
	}

	sort.Sort(byWidth(nodes))
	assert.Equal(t, nodes[0].Name, "first")
	assert.Equal(t, nodes[1].Name, "second")
	assert.Equal(t, nodes[2].Name, "third")
}
