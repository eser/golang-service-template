// This module is taken from the Go standard library and
// modified to work with the Bliss framework.

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the go stdlib's LICENSE
// file.

package uris

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/eser/go-service/pkg/bliss/results"
)

var (
	ErrPatternParsing  = results.Define("ERRBHUP001", "unable to parse pattern")
	ErrInvalidWildcard = results.Define("ERRBHUP001", "invalid wildcard")
	ErrInvalidMethod   = results.Define("ERRBHUP001", "invalid method")
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
		return nil, ErrPatternParsing.New("empty pattern")
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
		return nil, ErrInvalidMethod.New().
			WithAttribute(
				slog.String("pattern", s),
				slog.String("method", method),
			)
	}

	p := &Pattern{Str: s, Method: method}

	if found {
		off = len(method) + 1
	}

	i := strings.IndexByte(rest, '/')

	if i < 0 {
		return nil, ErrPatternParsing.New("host/path missing /").WithAttribute(slog.String("pattern", p.Str), slog.String("method", p.Method))
	}

	p.Host = rest[:i]
	rest = rest[i:]

	if j := strings.IndexByte(p.Host, '{'); j >= 0 {
		off += j

		return nil, ErrPatternParsing.New("host contains '{' (missing initial '/'?)").WithAttribute(slog.String("pattern", p.Str), slog.String("method", p.Method), slog.String("host", p.Host))
	}
	// At this point, rest is the path.
	off += i

	// An unclean path with a method that is not CONNECT can never match,
	// because paths are cleaned before matching.
	if method != "" && method != "CONNECT" && rest != CleanPath(rest) {
		return nil, ErrPatternParsing.New("non-CONNECT pattern with unclean path can never match").WithAttribute(slog.String("pattern", p.Str), slog.String("method", p.Method), slog.String("host", p.Host), slog.String("path", rest))
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
				return nil, ErrInvalidWildcard.New("bad wildcard segment (must start with '{')").WithAttribute(slog.String("pattern", p.Str), slog.String("method", p.Method), slog.String("segment", seg), slog.String("expected", "{"))
			}

			if seg[len(seg)-1] != '}' {
				return nil, ErrInvalidWildcard.New("bad wildcard segment (must end with '}')").WithAttribute(slog.String("pattern", p.Str), slog.String("method", p.Method), slog.String("segment", seg), slog.String("expected", "}"))
			}

			name := seg[1 : len(seg)-1]

			if name == "$" {
				if len(rest) != 0 {
					return nil, ErrInvalidMethod.New("{$} wildcard not at end").WithAttribute(slog.String("pattern", p.Str), slog.String("method", p.Method), slog.String("segment", seg))
				}

				p.Segments = append(p.Segments, Segment{Str: "/"})

				break
			}

			name, multi := strings.CutSuffix(name, "...")

			if multi && len(rest) != 0 {
				return nil, ErrInvalidWildcard.New("{...} wildcard not at end").WithAttribute(slog.String("pattern", p.Str), slog.String("method", p.Method), slog.String("segment", seg))
			}

			if name == "" {
				return nil, ErrInvalidWildcard.New("empty wildcard").WithAttribute(slog.String("pattern", p.Str), slog.String("method", p.Method))
			}

			if !isValidWildcardName(name) {
				return nil, ErrInvalidWildcard.New("bad wildcard name").WithAttribute(slog.String("pattern", p.Str), slog.String("method", p.Method), slog.String("segment", seg), slog.String("name", name))
			}

			if seenNames[name] {
				return nil, ErrInvalidWildcard.New("duplicate wildcard name").WithAttribute(slog.String("pattern", p.Str), slog.String("method", p.Method), slog.String("segment", seg), slog.String("name", name))
			}

			seenNames[name] = true
			p.Segments = append(p.Segments, Segment{Str: name, Wild: true, Multi: multi})
		}
	}

	return p, nil
}
