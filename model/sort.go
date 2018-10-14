package model

import (
	"strings"
)

type byWidth []*Node

func (s byWidth) Len() int {
	return len(s)
}
func (s byWidth) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byWidth) Less(i, j int) bool {
	if s[i].Width < s[j].Width {
		return true
	}

	if s[i].Width > s[j].Width {
		return false
	}

	return strings.Compare(s[i].Name, s[j].Name) == -1
}
