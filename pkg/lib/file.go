package lib

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	regexpFile   = regexp.MustCompile(`([^/]+\.go)(?:\.\((\w+)\))?$`)
	regexpGithub = regexp.MustCompile(`github\.com\/([^\/]+)\/([^\/]+)`)
)

func TrimTemporaryPath(tmpFolder, path, repository string) string {
	return strings.TrimPrefix(path, fmt.Sprintf("%s/%s", tmpFolder, repository))
}

func GetFileAndStruct(identifier string) (fileName, structName string) {
	result := regexpFile.FindStringSubmatch(identifier)
	if len(result) > 1 {
		fileName = result[1]
	}

	if len(result) > 2 {
		structName = result[2]
	}

	return
}

func GetIdentifier(tmpFolder, path, pkg, name string) string {
	path = TrimTemporaryPath(tmpFolder, path, pkg)
	if len(name) > 0 {
		return fmt.Sprintf("%s.(%s)", path, name)
	}
	return path
}

func IsGoFile(name string) bool {
	return strings.HasSuffix(name, ".go")
}

func GetGithubBaseURL(path string) (string, bool) {
	result := regexpGithub.FindStringSubmatch(path)
	if len(result) > 2 {
		return fmt.Sprintf("github.com/%s/%s", result[1], result[2]), true
	}
	return "", false
}
