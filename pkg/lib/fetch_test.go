package lib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	tmpFolder, _ := os.MkdirTemp("", "")

	f := NewFetcher(tmpFolder)
	assert.Implements(t, new(Fetcher), f)

	t.Run("invalid package", func(t *testing.T) {
		path, err := f.Fetch("invalid", "master")
		assert.Error(t, err)
		require.Empty(t, path)
	})

	t.Run("valid package", func(t *testing.T) {
		path, err := f.Fetch("github.com/rodrigo-brito/gocity", "master")
		assert.NoError(t, err)
		require.NotEmpty(t, path)
	})
}
