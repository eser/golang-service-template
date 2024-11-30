package uris_test

import (
	"testing"

	"github.com/eser/go-service/pkg/bliss/httpfx/uris"
	"github.com/stretchr/testify/assert"
)

func TestCommonPath(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := []struct {
		name     string
		p1       uris.Pattern
		p2       uris.Pattern
		expected string
	}{
		{
			name:     "SameSegments",
			p1:       uris.Pattern{Segments: []uris.Segment{{Wild: false, Str: "foo"}}}, //nolint:exhaustruct
			p2:       uris.Pattern{Segments: []uris.Segment{{Wild: false, Str: "foo"}}}, //nolint:exhaustruct
			expected: "/foo",
		},
		{
			name:     "DifferentSegments",
			p1:       uris.Pattern{Segments: []uris.Segment{{Wild: false, Str: "foo"}}}, //nolint:exhaustruct
			p2:       uris.Pattern{Segments: []uris.Segment{{Wild: false, Str: "bar"}}}, //nolint:exhaustruct
			expected: "/foo",
		},
		{
			name:     "DifferentSegmentsWithWildcard",
			p1:       uris.Pattern{Segments: []uris.Segment{{Wild: true, Str: "foo"}}},  //nolint:exhaustruct
			p2:       uris.Pattern{Segments: []uris.Segment{{Wild: false, Str: "bar"}}}, //nolint:exhaustruct
			expected: "/bar",
		},
		{
			name:     "WildcardAndMultipleSegments",
			p1:       uris.Pattern{Segments: []uris.Segment{{Wild: true, Str: "foo", Multi: true}}}, //nolint:exhaustruct
			p2:       uris.Pattern{Segments: []uris.Segment{{Wild: false, Str: "bar"}}},             //nolint:exhaustruct
			expected: "/bar",
		},
		{
			name: "FourSegments",
			p1: uris.Pattern{Segments: []uris.Segment{ //nolint:exhaustruct
				{Wild: true, Str: "foo"}, //nolint:exhaustruct
				{Wild: true, Str: "bar"}, //nolint:exhaustruct
				{Wild: true, Str: "baz"}, //nolint:exhaustruct
				{Wild: true, Str: "qux"}, //nolint:exhaustruct
			}},
			p2: uris.Pattern{Segments: []uris.Segment{ //nolint:exhaustruct
				{Wild: true, Str: "foo"}, //nolint:exhaustruct
				{Wild: true, Str: "bar"}, //nolint:exhaustruct
				{Wild: true, Str: "baz"}, //nolint:exhaustruct
				{Wild: true, Str: "qux"}, //nolint:exhaustruct
			}},
			expected: "/foo/bar/baz/qux",
		},
		{
			name: "FourP1SegmentWithThreeP2Segment",
			p1: uris.Pattern{Segments: []uris.Segment{ //nolint:exhaustruct
				{Wild: true, Str: "foo"}, //nolint:exhaustruct
				{Wild: true, Str: "bar"}, //nolint:exhaustruct
				{Wild: true, Str: "baz"}, //nolint:exhaustruct
				{Wild: true, Str: "qux"}, //nolint:exhaustruct
			}},
			p2: uris.Pattern{Segments: []uris.Segment{ //nolint:exhaustruct
				{Wild: true, Str: "foo"},  //nolint:exhaustruct
				{Wild: true, Str: "bar"},  //nolint:exhaustruct
				{Wild: true, Str: "test"}, //nolint:exhaustruct
			}},
			expected: "/foo/bar/test/qux",
		},
		{
			name: "ThreeP1SegmentWithFourP2Segment",
			p1: uris.Pattern{Segments: []uris.Segment{ //nolint:exhaustruct
				{Wild: true, Str: "foo"}, //nolint:exhaustruct
				{Wild: true, Str: "bar"}, //nolint:exhaustruct
				{Wild: true, Str: "baz"}, //nolint:exhaustruct
			}},
			p2: uris.Pattern{Segments: []uris.Segment{ //nolint:exhaustruct
				{Wild: true, Str: "foo"}, //nolint:exhaustruct
				{Wild: true, Str: "bar"}, //nolint:exhaustruct
				{Wild: true, Str: "baz"}, //nolint:exhaustruct
				{Wild: true, Str: "qux"}, //nolint:exhaustruct
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
