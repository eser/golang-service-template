package uris_test

import (
	"testing"

	"github.com/eser/go-service/pkg/bliss/httpfx/uris"
	"github.com/stretchr/testify/assert"
)

func TestDifferencePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		p1   *uris.Pattern
		p2   *uris.Pattern
		want string
	}{
		{
			name: "Both patterns have multi segments",
			p1:   &uris.Pattern{Segments: []uris.Segment{{Multi: true}, {Multi: true}}},
			p2:   &uris.Pattern{Segments: []uris.Segment{{Multi: true}, {Multi: true}}},
			want: "/",
		},
		{
			name: "p1 has multi, p2 doesn't, s2 ends in /",
			p1:   &uris.Pattern{Segments: []uris.Segment{{Multi: true, Str: "foo"}, {Multi: true}}},
			p2:   &uris.Pattern{Segments: []uris.Segment{{Str: "bar"}, {Str: "/"}}},
			want: "/",
		},
		{
			name: "p1 has multi, p2 doesn't, s2 doesn't end in /",
			p1:   &uris.Pattern{Segments: []uris.Segment{{Multi: true}, {Multi: true}}},
			p2:   &uris.Pattern{Segments: []uris.Segment{{Str: "bar"}, {Str: "baz"}}},
			want: "/",
		},
		{
			name: "p2 has multi, p1 doesn't",
			p1:   &uris.Pattern{Segments: []uris.Segment{{Str: "foo"}, {Str: "bar"}}},
			p2:   &uris.Pattern{Segments: []uris.Segment{{Multi: true}, {Multi: true}}},
			want: "/foo/bar",
		},
		{
			name: "Both patterns have wild segments, same name",
			p1:   &uris.Pattern{Segments: []uris.Segment{{Wild: true, Str: "foo"}, {Wild: true, Str: "bar"}}},
			p2:   &uris.Pattern{Segments: []uris.Segment{{Wild: true, Str: "foo"}, {Wild: true, Str: "bar"}}},
			want: "/foo/bar",
		},
		{
			name: "p1 has wild, p2 doesn't, different names",
			p1:   &uris.Pattern{Segments: []uris.Segment{{Wild: true, Str: "foo"}, {Wild: true, Str: "bar"}}},
			p2:   &uris.Pattern{Segments: []uris.Segment{{Str: "baz"}, {Str: "qux"}}},
			want: "/foo/bar",
		},
		{
			name: "p1 has wild, p2 doesn't, same names",
			p1:   &uris.Pattern{Segments: []uris.Segment{{Wild: true, Str: "foo"}, {Wild: true, Str: "bar"}}},
			p2:   &uris.Pattern{Segments: []uris.Segment{{Str: "foo"}, {Str: "bar"}}},
			want: "/foox/barx",
		},
		{
			name: "p2 has wild, p1 doesn't",
			p1:   &uris.Pattern{Segments: []uris.Segment{{Str: "foo"}, {Str: "bar"}}},
			p2:   &uris.Pattern{Segments: []uris.Segment{{Wild: true, Str: "baz"}, {Wild: true, Str: "qux"}}},
			want: "/foo/bar",
		},
		{
			name: "Both are literals, same",
			p1:   &uris.Pattern{Segments: []uris.Segment{{Str: "foo"}, {Str: "bar"}}},
			p2:   &uris.Pattern{Segments: []uris.Segment{{Str: "foo"}, {Str: "bar"}}},
			want: "/foo/bar",
		},
		{
			name: "Both are literals, different (should panic)",
			p1:   &uris.Pattern{Segments: []uris.Segment{{Str: "foo"}, {Str: "bar"}}},
			p2:   &uris.Pattern{Segments: []uris.Segment{{Str: "baz"}, {Str: "qux"}}},
			want: "", // This test should panic
		},
		{
			name: "p1 is longer than p2",
			p1:   &uris.Pattern{Segments: []uris.Segment{{Str: "foo"}, {Str: "bar"}, {Str: "baz"}}},
			p2:   &uris.Pattern{Segments: []uris.Segment{{Str: "foo"}, {Str: "bar"}}},
			want: "/foo/bar/baz",
		},
		{
			name: "p2 is longer than p1",
			p1:   &uris.Pattern{Segments: []uris.Segment{{Str: "foo"}, {Str: "bar"}}},
			p2:   &uris.Pattern{Segments: []uris.Segment{{Str: "foo"}, {Str: "bar"}, {Str: "baz"}}},
			want: "/foo/bar/baz",
		},
		{
			name: "p1 has multi and it is empty, and p2 is only slash",
			p1:   &uris.Pattern{Segments: []uris.Segment{{Multi: true}}},
			p2:   &uris.Pattern{Segments: []uris.Segment{{Str: "/"}}},
			want: "/x",
		},
		{
			name: "p1 has multi and it is only slah, and p2 is only slash",
			p1:   &uris.Pattern{Segments: []uris.Segment{{Str: "/", Multi: true}}},
			p2:   &uris.Pattern{Segments: []uris.Segment{{Str: "/"}}},
			want: "//",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				r := recover()

				if tt.want == "" {
					assert.NotNil(t, r, "DifferencePath() did not panic")

					return
				}

				assert.Nil(t, r, "DifferencePath() panicked: %v", r)
			}()

			got := uris.DifferencePath(tt.p1, tt.p2)
			assert.Equal(t, tt.want, got)
		})
	}
}
