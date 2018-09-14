package main

import (
	"fmt"
	"strings"
)

func getIdentifier(path, name string) string {
	if len(name) > 0 {
		return fmt.Sprintf("%s.(%s)", path, name)
	}
	return path
}

func isGoFile(name string) bool {
	return strings.HasSuffix(name, ".go")
}
