package lib

import (
	"bytes"
	"strings"
	"unicode"
)

func TrimLeadingSpaceFromBytes(src []byte) []byte {
	return bytes.TrimLeftFunc(src, unicode.IsSpace)
}

func TrimTrailingSpaceFromBytes(src []byte) []byte {
	return bytes.TrimRightFunc(src, unicode.IsSpace)
}

func TrimLeadingSpace(src string) string {
	return strings.TrimLeftFunc(src, unicode.IsSpace)
}

func TrimTrailingSpace(src string) string {
	return strings.TrimRightFunc(src, unicode.IsSpace)
}
