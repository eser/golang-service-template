package lib

import (
	"bytes"
	"strings"
	"unicode"
)

func StringsTrimLeadingSpaceFromBytes(src []byte) []byte {
	return bytes.TrimLeftFunc(src, unicode.IsSpace)
}

func StringsTrimTrailingSpaceFromBytes(src []byte) []byte {
	return bytes.TrimRightFunc(src, unicode.IsSpace)
}

func StringsTrimLeadingSpace(src string) string {
	return strings.TrimLeftFunc(src, unicode.IsSpace)
}

func StringsTrimTrailingSpace(src string) string {
	return strings.TrimRightFunc(src, unicode.IsSpace)
}
