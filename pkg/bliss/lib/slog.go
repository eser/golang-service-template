package lib

import (
	"log/slog"
	"strings"
)

func SerializeSlogAttrs(attrs ...slog.Attr) string {
	length := len(attrs)

	if length == 0 {
		return ""
	}

	result := make([]string, len(attrs))

	for _, attr := range attrs {
		result = append(result, attr.String())
	}

	return strings.Join(result, ", ")
}
