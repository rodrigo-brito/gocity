package utils

import (
	"fmt"
	"strings"
)

func GetIdentifier(path, name string) string {
	if len(name) > 0 {
		return fmt.Sprintf("%s.(%s)", path, name)
	}
	return path
}

func IsGoFile(name string) bool {
	return strings.HasSuffix(name, ".go")
}
