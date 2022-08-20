package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPositionGenerator(t *testing.T) {
	generator := NewGenerator(4)
	position := generator.NextPosition(5, 10)
	assert.Equal(t, float64(3), position.X)
	assert.Equal(t, float64(5.5), position.Y)

	position = generator.NextPosition(2, 2)
	assert.Equal(t, float64(1.5), position.X)
	assert.Equal(t, float64(12.5), position.Y)

	position = generator.NextPosition(10, 10)
	assert.Equal(t, float64(11.5), position.X)
	assert.Equal(t, float64(5.5), position.Y)

	position = generator.NextPosition(10, 10)
	assert.Equal(t, float64(11.5), position.X)
	assert.Equal(t, float64(16.5), position.Y)

	bounds := generator.GetBounds()
	assert.Equal(t, float64(17), bounds.X)
	assert.Equal(t, float64(22), bounds.Y)
}
