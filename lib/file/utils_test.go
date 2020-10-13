package utils

import (
	"fmt"
	"os"
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
		{Input: "foo/bar/file.pb.go", FileName: "file.pb.go", StructName: ""},
		{Input: "foo/bar/file.pb_test.go", FileName: "file.pb_test.go", StructName: ""},
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

func TestGetGithubBaseURL(t *testing.T) {
	tt := []struct {
		Input   string
		Output  string
		IsValid bool
	}{
		{"github.com/foo/bar", "github.com/foo/bar", true},
		{"https://github.com/foo/bar", "github.com/foo/bar", true},
		{"github.com/foo/bar/subpackage", "github.com/foo/bar", true},
		{"www.github.com/foo/bar/subpackage", "github.com/foo/bar", true},
		{"www.gitlab.com/foo/bar/subpackage", "", false},
		{"github.com/foo", "", false},
		{"invalid", "", false},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("given the input %s", tc.Input), func(t *testing.T) {
			output, valid := GetGithubBaseURL(tc.Input)
			assert.Equal(t, tc.IsValid, valid)
			assert.Equal(t, tc.Output, output)
		})
	}
}

func TestIsGoFile(t *testing.T) {
	tests := []struct {
		got  string
		want bool
	}{
		{"foo.go", true},
		{"bar.gol", false},
		{"foobar", false},
		{"fubar.g", false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("given the filename %s", tt.got), func(t *testing.T) {
			got := IsGoFile(tt.got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTrimGoPath(t *testing.T) {
	tests := []struct {
		path       string
		repository string
		want       string
	}{
		{fmt.Sprintf("%s/src/gocity/main.go", os.Getenv("GOPATH")), "gocity", "/main.go"},
		{fmt.Sprintf("%s/src/gocity/foo/bar.go", os.Getenv("GOPATH")), "gocity", "/foo/bar.go"},
		{fmt.Sprintf("%s/src/gocity/vendor", os.Getenv("GOPATH")), "gocity", "/vendor"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("given project %s/%s", tt.path, tt.repository), func(t *testing.T) {
			got := TrimGoPath(tt.path, tt.repository)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetIdentifier(t *testing.T) {
	tests := []struct {
		path string
		pkg  string
		name string
		want string
	}{
		{fmt.Sprintf("%s/src/gocity/main.go", os.Getenv("GOPATH")), "gocity", "/main.go", "/main.go.(/main.go)"},
		{fmt.Sprintf("%s/src/gocity/foo/bar.go", os.Getenv("GOPATH")), "gocity", "/foo/bar.go", "/foo/bar.go.(/foo/bar.go)"},
		{fmt.Sprintf("%s/src/gocity/vendor", os.Getenv("GOPATH")), "gocity", "/vendor", "/vendor.(/vendor)"},
		{fmt.Sprintf("%s/src/gocity/vendor", os.Getenv("GOPATH")), "gocity", "", "/vendor"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("given path %s, pkg %s and name %s", tt.path, tt.pkg, tt.name), func(t *testing.T) {
			got := GetIdentifier(tt.path, tt.pkg, tt.name)
			assert.Equal(t, tt.want, got)
		})
	}
}
