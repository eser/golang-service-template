package uris

import (
	"testing"
)

func TestCommonPath(t *testing.T) {
	tests := []struct {
		name     string
		p1       Pattern
		p2       Pattern
		expected string
	}{
		{
			name:     "SameSegments",
			p1:       Pattern{Segments: []Segment{{wild: false, s: "foo"}}},
			p2:       Pattern{Segments: []Segment{{wild: false, s: "foo"}}},
			expected: "/foo",
		},
		{
			name:     "DifferentSegments",
			p1:       Pattern{Segments: []Segment{{wild: false, s: "foo"}}},
			p2:       Pattern{Segments: []Segment{{wild: false, s: "bar"}}},
			expected: "/foo",
		},
		{
			name:     "DifferentSegmentsWithWildcard",
			p1:       Pattern{Segments: []Segment{{wild: true, s: "foo"}}},
			p2:       Pattern{Segments: []Segment{{wild: false, s: "bar"}}},
			expected: "/bar",
		},
		{
			name:     "WilcardAndMultipleSegments",
			p1:       Pattern{Segments: []Segment{{wild: true, s: "foo", multi: true}}},
			p2:       Pattern{Segments: []Segment{{wild: false, s: "bar"}}},
			expected: "/bar",
		},
		{
			name: "FourSegments",
			p1: Pattern{Segments: []Segment{
				{wild: true, s: "foo"},
				{wild: true, s: "bar"},
				{wild: true, s: "baz"},
				{wild: true, s: "qux"},
			}},
			p2: Pattern{Segments: []Segment{
				{wild: true, s: "foo"},
				{wild: true, s: "bar"},
				{wild: true, s: "baz"},
				{wild: true, s: "qux"},
			}},
			expected: "/foo/bar/baz/qux",
		},
		{
			name: "FourP1SegmentWithThreeP2Segment",
			p1: Pattern{Segments: []Segment{
				{wild: true, s: "foo"},
				{wild: true, s: "bar"},
				{wild: true, s: "baz"},
				{wild: true, s: "qux"},
			}},
			p2: Pattern{Segments: []Segment{
				{wild: true, s: "foo"},
				{wild: true, s: "bar"},
				{wild: true, s: "test"},
			}},
			expected: "/foo/bar/test/qux",
		},
		{
			name: "ThreeP1SegmentWithFourP2Segment",
			p1: Pattern{Segments: []Segment{
				{wild: true, s: "foo"},
				{wild: true, s: "bar"},
				{wild: true, s: "baz"},
			}},
			p2: Pattern{Segments: []Segment{
				{wild: true, s: "foo"},
				{wild: true, s: "bar"},
				{wild: true, s: "baz"},
				{wild: true, s: "qux"},
			}},
			expected: "/foo/bar/baz/qux",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CommonPath(&tt.p1, &tt.p2)
			if result != tt.expected {
				t.Errorf("CommonPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}
