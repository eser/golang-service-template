package lib

import (
	"bytes"
	"strings"
	"unicode"
)

func StringsTrimLeadingSpaceFromBytes(src []byte) []byte {
	if len(src) == 0 {
		return src
	}

	return bytes.TrimLeftFunc(src, unicode.IsSpace)
}

func StringsTrimTrailingSpaceFromBytes(src []byte) []byte {
	if len(src) == 0 {
		return src
	}

	return bytes.TrimRightFunc(src, unicode.IsSpace)
}

func StringsTrimLeadingSpace(src string) string {
	if len(src) == 0 {
		return src
	}

	return strings.TrimLeftFunc(src, unicode.IsSpace)
}

func StringsTrimTrailingSpace(src string) string {
	if len(src) == 0 {
		return src
	}

	return strings.TrimRightFunc(src, unicode.IsSpace)
}
