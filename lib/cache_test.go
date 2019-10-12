package lib

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	c := NewCache()
	assert.Implements(t, new(Cache), c)

	ok, res := c.Get("test")
	assert.False(t, ok)
	assert.Empty(t, res)

	c.Set("test", []byte("data"), 10*time.Second)

	ok, res = c.Get("test")
	assert.True(t, ok)
	assert.NotEmpty(t, res)

	res, err := c.GetSet("testing", func() ([]byte, error) {
		return []byte{}, errors.New("some error")
	}, 3*time.Second)

	assert.Empty(t, res)
	assert.Error(t, err)

	res, err = c.GetSet("testing", func() ([]byte, error) {
		return []byte("data"), nil
	}, 3*time.Second)

	assert.NotEmpty(t, res)
	assert.NoError(t, err)

	res, err = c.GetSet("test", func() ([]byte, error) {
		return []byte{}, nil
	}, 3*time.Second)

	assert.NotEmpty(t, res)
	assert.NoError(t, err)
}
