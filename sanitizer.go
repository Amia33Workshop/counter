package main

import (
	"strings"
	"unicode"
)

func SanitizeLogString(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, ch := range s {
		if unicode.IsPrint(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}
