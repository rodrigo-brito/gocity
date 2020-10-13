package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	f := NewFetcher()
	assert.Implements(t, new(Fetcher), f)

	err := f.Fetch("invalid", "master")
	assert.Error(t, err)

	err = f.Fetch("github.com/rodrigo-brito/gocity", "master")
	assert.NoError(t, err)

}
