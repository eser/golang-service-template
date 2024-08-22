package uris_test

import (
	"testing"

	"github.com/eser/go-service/pkg/bliss/httpfx/uris"
	"github.com/stretchr/testify/assert"
)

func TestCommonPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		p1       uris.Pattern
		p2       uris.Pattern
		expected string
	}{
		{
			name:     "SameSegments",
			p1:       uris.Pattern{Segments: []uris.Segment{{Wild: false, Str: "foo"}}},
			p2:       uris.Pattern{Segments: []uris.Segment{{Wild: false, Str: "foo"}}},
			expected: "/foo",
		},
		{
			name:     "DifferentSegments",
			p1:       uris.Pattern{Segments: []uris.Segment{{Wild: false, Str: "foo"}}},
			p2:       uris.Pattern{Segments: []uris.Segment{{Wild: false, Str: "bar"}}},
			expected: "/foo",
		},
		{
			name:     "DifferentSegmentsWithWildcard",
			p1:       uris.Pattern{Segments: []uris.Segment{{Wild: true, Str: "foo"}}},
			p2:       uris.Pattern{Segments: []uris.Segment{{Wild: false, Str: "bar"}}},
			expected: "/bar",
		},
		{
			name:     "WildcardAndMultipleSegments",
			p1:       uris.Pattern{Segments: []uris.Segment{{Wild: true, Str: "foo", Multi: true}}},
			p2:       uris.Pattern{Segments: []uris.Segment{{Wild: false, Str: "bar"}}},
			expected: "/bar",
		},
		{
			name: "FourSegments",
			p1: uris.Pattern{Segments: []uris.Segment{
				{Wild: true, Str: "foo"},
				{Wild: true, Str: "bar"},
				{Wild: true, Str: "baz"},
				{Wild: true, Str: "qux"},
			}},
			p2: uris.Pattern{Segments: []uris.Segment{
				{Wild: true, Str: "foo"},
				{Wild: true, Str: "bar"},
				{Wild: true, Str: "baz"},
				{Wild: true, Str: "qux"},
			}},
			expected: "/foo/bar/baz/qux",
		},
		{
			name: "FourP1SegmentWithThreeP2Segment",
			p1: uris.Pattern{Segments: []uris.Segment{
				{Wild: true, Str: "foo"},
				{Wild: true, Str: "bar"},
				{Wild: true, Str: "baz"},
				{Wild: true, Str: "qux"},
			}},
			p2: uris.Pattern{Segments: []uris.Segment{
				{Wild: true, Str: "foo"},
				{Wild: true, Str: "bar"},
				{Wild: true, Str: "test"},
			}},
			expected: "/foo/bar/test/qux",
		},
		{
			name: "ThreeP1SegmentWithFourP2Segment",
			p1: uris.Pattern{Segments: []uris.Segment{
				{Wild: true, Str: "foo"},
				{Wild: true, Str: "bar"},
				{Wild: true, Str: "baz"},
			}},
			p2: uris.Pattern{Segments: []uris.Segment{
				{Wild: true, Str: "foo"},
				{Wild: true, Str: "bar"},
				{Wild: true, Str: "baz"},
				{Wild: true, Str: "qux"},
			}},
			expected: "/foo/bar/baz/qux",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := uris.CommonPath(&tt.p1, &tt.p2)

			assert.Equal(t, tt.expected, result)
		})
	}
}
