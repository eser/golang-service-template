// This module is taken from the Go standard library and
// modified to work with the Bliss framework.

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the go stdlib's LICENSE
// file.

package uris

import (
	"net/url"
	"strings"
	"unicode"
)

func isValidWildcardName(s string) bool {
	if s == "" {
		return false
	}
	// Valid Go identifier.
	for i, c := range s {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}

	return true
}

func tryPathUnescape(path string) string {
	u, err := url.PathUnescape(path)
	if err != nil {
		// Invalidly escaped path; use the original
		return path
	}

	return u
}

// writeMatchingPath writes to b a path that matches the segments.
func writeMatchingPath(b *strings.Builder, segs []Segment) {
	for _, s := range segs {
		writeSegment(b, s)
	}
}

func writeSegment(b *strings.Builder, s Segment) {
	b.WriteByte('/')

	if !s.Multi && s.Str != "/" {
		b.WriteString(s.Str)
	}
}
