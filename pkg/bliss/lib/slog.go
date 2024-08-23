package lib

import (
	"log/slog"
	"strings"
)

func SerializeSlogAttrs(attrs []slog.Attr) string {
	length := len(attrs)

	if length == 0 {
		return ""
	}

	result := make([]string, length)

	for i, attr := range attrs {
		result[i] = attr.String()
	}

	return strings.Join(result, ", ")
}
