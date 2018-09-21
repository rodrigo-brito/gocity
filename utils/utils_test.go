package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileAndStruct(t *testing.T) {
	tt := []struct {
		Input      string
		FileName   string
		StructName string
	}{
		{Input: "foo/bar/file.go", FileName: "file.go", StructName: ""},
		{Input: "foo/bar/file.go.(Test)", FileName: "file.go", StructName: "Test"},
		{Input: "foo/bar/test.go.(test)", FileName: "test.go", StructName: "test"},
		{Input: "foo/bar/9999", FileName: "", StructName: ""},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("given the input %s", tc.Input), func(t *testing.T) {
			fileName, structName := GetFileAndStruct(tc.Input)
			assert.Equal(t, tc.FileName, fileName)
			assert.Equal(t, tc.StructName, structName)
		})
	}
}
