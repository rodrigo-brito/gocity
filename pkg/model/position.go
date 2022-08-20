package model

import (
	"math"
)

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

const defaultMargin = 1

type generator struct {
	numberNodes  int
	dimension    int
	xReference   float64
	yReference   float64
	currentIndex int
	maxWidth     float64
	maxHeight    float64
}

func NewGenerator(numberNodes int) *generator {
	generator := &generator{
		numberNodes: numberNodes,
		dimension:   int(math.Ceil(math.Sqrt(float64(numberNodes)))),
	}

	return generator
}

func (g *generator) GetBounds() Position {
	return Position{
		X: g.maxWidth + defaultMargin,
		Y: g.maxHeight + defaultMargin,
	}
}

func (g *generator) NextPosition(width, height float64) Position {
	g.currentIndex++

	if g.currentIndex > g.dimension && g.yReference+height >= g.maxWidth {
		g.currentIndex = 0
		g.yReference = 0
		g.xReference = g.maxWidth + defaultMargin
	}

	position := Position{X: g.xReference + (width+defaultMargin)/2, Y: g.yReference + (height+defaultMargin)/2}

	if g.xReference+width > g.maxWidth {
		g.maxWidth = g.xReference + width
	}

	if g.yReference+height > g.maxHeight {
		g.maxHeight = g.yReference + height
	}

	g.yReference += height + defaultMargin

	return position
}
