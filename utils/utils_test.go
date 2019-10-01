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
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestTrimGoPath(t *testing.T) {
	type args struct {
		path       string
		repository string
	}
	tests := []struct {
		args args
		got  string
		want string
	}{
		{args{"/home/dude/go/src/github.com/repo", "github.com/dude/repo.git"}, "somerepo", "home/bla/src/github/somerp"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Trim GOPATH"), func(t *testing.T) {
			got := TrimGoPath(tt.args.path, tt.args.repository)
			fmt.Println(got)
			assert.Equal(t, got, tt.want)
		})
	}
}
