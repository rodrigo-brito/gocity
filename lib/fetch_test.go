package lib

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	tmpFolder, _ := ioutil.TempDir("", "")

	f := NewFetcher(tmpFolder)
	assert.Implements(t, new(Fetcher), f)

	err := f.Fetch("invalid", "master")
	assert.Error(t, err)

	err = f.Fetch("github.com/rodrigo-brito/gocity", "master")
	assert.NoError(t, err)

}
