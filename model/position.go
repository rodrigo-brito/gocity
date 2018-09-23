package model

import (
	"math"
)

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type generator struct {
	maxSize     float64
	numberNodes int
	margin      int
	positions   []Position
}

func NewGenerator(numberNode int, maxSize float64) *generator {
	generator := &generator{
		maxSize:     maxSize,
		numberNodes: numberNode,
	}
	generator.Init()
	return generator
}

func (g *generator) Init() {
	dimension := math.Ceil(math.Sqrt(float64(g.numberNodes)))
	halfWidth := float64(g.maxSize / 2.0)
	shiftCorner := dimension*halfWidth - halfWidth
	for i := 0; i < int(dimension); i++ {
		for j := 0; j < int(dimension); j++ {
			g.positions = append(g.positions, Position{
				X: g.maxSize*float64(j) - shiftCorner,
				Y: g.maxSize*float64(i) - shiftCorner,
			})
		}
	}
}

func (g *generator) NextPosition() Position {
	position := g.positions[0]
	g.positions = append(g.positions[:0], g.positions[1:]...)
	return position
}
