// This module is taken from the Go standard library and
// modified to work with the Bliss framework.

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the go stdlib's LICENSE
// file.

package uris

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyPattern                     = errors.New("empty pattern")
	ErrPatternParsing                   = errors.New("unable to parse pattern")
	ErrNonConnectPatternWithUncleanPath = errors.New("non-CONNECT pattern with unclean path can never match")
	ErrInvalidWildcard                  = errors.New("invalid wildcard")
	ErrEmptyWildcard                    = errors.New("empty wildcard")
	ErrInvalidMethod                    = errors.New("invalid method")
)

// A pattern is something that can be matched against an HTTP request.
// It has an optional method, an optional host, and a path.
type Pattern struct {
	Str    string // original string
	Method string
	Host   string
	Path   string
	// The representation of a path differs from the surface syntax, which
	// simplifies most algorithms.
	//
	// Paths ending in '/' are represented with an anonymous "..." wildcard.
	// For example, the path "a/" is represented as a literal segment "a" followed
	// by a segment with multi==true.
	//
	// Paths ending in "{$}" are represented with the literal segment "/".
	// For example, the path "a/{$}" is represented as a literal segment "a" followed
	// by a literal segment "/".
	Segments []Segment
	Loc      string // source location of registering call, for helpful messages
}

func (p *Pattern) String() string { return p.Str }

// A segment is a pattern piece that matches one or more path segments, or
// a trailing slash.
//
// If wild is false, it matches a literal segment, or, if s == "/", a trailing slash.
// Examples:
//
//	"a" => segment{s: "a"}
//	"/{$}" => segment{s: "/"}
//
// If wild is true and multi is false, it matches a single path segment.
// Example:
//
//	"{x}" => segment{s: "x", wild: true}
//
// If both wild and multi are true, it matches all remaining path segments.
// Example:
//
//	"{rest...}" => segment{s: "rest", wild: true, multi: true}
type Segment struct {
	Str   string // literal or wildcard name or "/" for "/{$}".
	Wild  bool
	Multi bool // "..." wildcard
}

// parsePattern parses a string into a Pattern.
// The string's syntax is
//
//	[METHOD] [HOST]/[PATH]
//
// where:
//   - METHOD is an HTTP method
//   - HOST is a hostname
//   - PATH consists of slash-separated segments, where each segment is either
//     a literal or a wildcard of the form "{name}", "{name...}", or "{$}".
//
// METHOD, HOST and PATH are all optional; that is, the string can be "/".
// If METHOD is present, it must be followed by a single space.
// Wildcard names must be valid Go identifiers.
// The "{$}" and "{name...}" wildcard must occur at the end of PATH.
// PATH may end with a '/'.
// Wildcard names in a path must be distinct.
func ParsePattern(s string) (_ *Pattern, err error) { //nolint:funlen,gocognit,cyclop
	if len(s) == 0 {
		return nil, ErrEmptyPattern
	}

	off := 0 // offset into string

	defer func() {
		if err != nil {
			err = fmt.Errorf("at offset %d: %w", off, err)
		}
	}()

	method, rest, found := strings.Cut(s, " ")
	if !found {
		rest = method
		method = ""
	}

	if method != "" && !IsValidMethod(method) {
		return nil, fmt.Errorf("invalid method %q: %w", method, ErrInvalidMethod)
	}

	p := &Pattern{Str: s, Method: method}

	if found {
		off = len(method) + 1
	}

	i := strings.IndexByte(rest, '/')

	if i < 0 {
		return nil, fmt.Errorf("host/path missing /: %w", ErrPatternParsing)
	}

	p.Host = rest[:i]
	rest = rest[i:]

	if j := strings.IndexByte(p.Host, '{'); j >= 0 {
		off += j

		return nil, fmt.Errorf("host contains '{' (missing initial '/'?): %w", ErrPatternParsing)
	}
	// At this point, rest is the path.
	off += i

	// An unclean path with a method that is not CONNECT can never match,
	// because paths are cleaned before matching.
	if method != "" && method != "CONNECT" && rest != CleanPath(rest) {
		return nil, ErrNonConnectPatternWithUncleanPath
	}

	p.Path = rest
	seenNames := map[string]bool{} // remember wildcard names to catch dups

	for len(rest) > 0 {
		// Invariant: rest[0] == '/'.
		rest = rest[1:]
		off = len(s) - len(rest)

		if len(rest) == 0 {
			// Trailing slash.
			p.Segments = append(p.Segments, Segment{Wild: true, Multi: true})

			break
		}

		i := strings.IndexByte(rest, '/')

		if i < 0 {
			i = len(rest)
		}

		var seg string
		seg, rest = rest[:i], rest[i:]

		if i := strings.IndexByte(seg, '{'); i < 0 { //nolint:nestif
			// Literal.
			seg = tryPathUnescape(seg)
			p.Segments = append(p.Segments, Segment{Str: seg})
		} else {
			// Wildcard.
			if i != 0 {
				return nil, fmt.Errorf("bad wildcard segment (must start with '{'): %w", ErrInvalidWildcard)
			}

			if seg[len(seg)-1] != '}' {
				return nil, fmt.Errorf("bad wildcard segment (must end with '}'): %w", ErrInvalidWildcard)
			}

			name := seg[1 : len(seg)-1]

			if name == "$" {
				if len(rest) != 0 {
					return nil, fmt.Errorf("{$} wildcard not at end: %w", ErrInvalidWildcard)
				}

				p.Segments = append(p.Segments, Segment{Str: "/"})

				break
			}

			name, multi := strings.CutSuffix(name, "...")

			if multi && len(rest) != 0 {
				return nil, fmt.Errorf("{...} wildcard not at end: %w", ErrInvalidWildcard)
			}

			if name == "" {
				return nil, ErrEmptyWildcard
			}

			if !isValidWildcardName(name) {
				return nil, fmt.Errorf("bad wildcard name %q: %w", name, ErrInvalidWildcard)
			}

			if seenNames[name] {
				return nil, fmt.Errorf("duplicate wildcard name %q: %w", name, ErrInvalidWildcard)
			}

			seenNames[name] = true
			p.Segments = append(p.Segments, Segment{Str: name, Wild: true, Multi: multi})
		}
	}

	return p, nil
}
