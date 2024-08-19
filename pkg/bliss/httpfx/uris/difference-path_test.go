package uris

import "testing"

func TestDifferencePath(t *testing.T) {
	tests := []struct {
		name string
		p1   *Pattern
		p2   *Pattern
		want string
	}{
		{
			name: "Both patterns have multi segments",
			p1:   &Pattern{Segments: []Segment{{multi: true}, {multi: true}}},
			p2:   &Pattern{Segments: []Segment{{multi: true}, {multi: true}}},
			want: "/",
		},
		{
			name: "p1 has multi, p2 doesn't, s2 ends in /",
			p1:   &Pattern{Segments: []Segment{{multi: true, s: "foo"}, {multi: true}}},
			p2:   &Pattern{Segments: []Segment{{s: "bar"}, {s: "/"}}},
			want: "/",
		},
		{
			name: "p1 has multi, p2 doesn't, s2 doesn't end in /",
			p1:   &Pattern{Segments: []Segment{{multi: true}, {multi: true}}},
			p2:   &Pattern{Segments: []Segment{{s: "bar"}, {s: "baz"}}},
			want: "/",
		},
		{
			name: "p2 has multi, p1 doesn't",
			p1:   &Pattern{Segments: []Segment{{s: "foo"}, {s: "bar"}}},
			p2:   &Pattern{Segments: []Segment{{multi: true}, {multi: true}}},
			want: "/foo/bar",
		},
		{
			name: "Both patterns have wild segments, same name",
			p1:   &Pattern{Segments: []Segment{{wild: true, s: "foo"}, {wild: true, s: "bar"}}},
			p2:   &Pattern{Segments: []Segment{{wild: true, s: "foo"}, {wild: true, s: "bar"}}},
			want: "/foo/bar",
		},
		{
			name: "p1 has wild, p2 doesn't, different names",
			p1:   &Pattern{Segments: []Segment{{wild: true, s: "foo"}, {wild: true, s: "bar"}}},
			p2:   &Pattern{Segments: []Segment{{s: "baz"}, {s: "qux"}}},
			want: "/foo/bar",
		},
		{
			name: "p1 has wild, p2 doesn't, same names",
			p1:   &Pattern{Segments: []Segment{{wild: true, s: "foo"}, {wild: true, s: "bar"}}},
			p2:   &Pattern{Segments: []Segment{{s: "foo"}, {s: "bar"}}},
			want: "/foox/barx",
		},
		{
			name: "p2 has wild, p1 doesn't",
			p1:   &Pattern{Segments: []Segment{{s: "foo"}, {s: "bar"}}},
			p2:   &Pattern{Segments: []Segment{{wild: true, s: "baz"}, {wild: true, s: "qux"}}},
			want: "/foo/bar",
		},
		{
			name: "Both are literals, same",
			p1:   &Pattern{Segments: []Segment{{s: "foo"}, {s: "bar"}}},
			p2:   &Pattern{Segments: []Segment{{s: "foo"}, {s: "bar"}}},
			want: "/foo/bar",
		},
		{
			name: "Both are literals, different (should panic)",
			p1:   &Pattern{Segments: []Segment{{s: "foo"}, {s: "bar"}}},
			p2:   &Pattern{Segments: []Segment{{s: "baz"}, {s: "qux"}}},
			want: "", // This test should panic
		},
		{
			name: "p1 is longer than p2",
			p1:   &Pattern{Segments: []Segment{{s: "foo"}, {s: "bar"}, {s: "baz"}}},
			p2:   &Pattern{Segments: []Segment{{s: "foo"}, {s: "bar"}}},
			want: "/foo/bar/baz",
		},
		{
			name: "p2 is longer than p1",
			p1:   &Pattern{Segments: []Segment{{s: "foo"}, {s: "bar"}}},
			p2:   &Pattern{Segments: []Segment{{s: "foo"}, {s: "bar"}, {s: "baz"}}},
			want: "/foo/bar/baz",
		},
		{
			name: "p1 has multi and it is empty, and p2 is only slash",
			p1:   &Pattern{Segments: []Segment{{multi: true}}},
			p2:   &Pattern{Segments: []Segment{{s: "/"}}},
			want: "/x",
		},
		{
			name: "p1 has multi and it is only slah, and p2 is only slash",
			p1:   &Pattern{Segments: []Segment{{s: "/", multi: true}}},
			p2:   &Pattern{Segments: []Segment{{s: "/"}}},
			want: "//",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tt.want != "" {
						t.Errorf("DifferencePath() panicked: %v, want %q", r, tt.want)
					}
				}
			}()

			got := DifferencePath(tt.p1, tt.p2)
			if got != tt.want {
				t.Errorf("DifferencePath() = %q, want %q", got, tt.want)
			}
		})
	}
}
