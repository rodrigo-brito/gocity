package utils

import (
	"fmt"
	"regexp"
	"strings"
)

var pattern = regexp.MustCompile(`(\w+\.go)(?:\.\((\w+)\))?$`)

func GetFileAndStruct(identifier string) (fileName, structName string) {
	result := pattern.FindStringSubmatch(identifier)
	if len(result) > 1 {
		fileName = result[1]
	}

	if len(result) > 2 {
		structName = result[2]
	}

	return
}

func GetIdentifier(path, name string) string {
	if len(name) > 0 {
		return fmt.Sprintf("%s.(%s)", path, name)
	}
	return path
}

func IsGoFile(name string) bool {
	return strings.HasSuffix(name, ".go")
}
