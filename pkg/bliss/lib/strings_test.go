package lib_test

import (
	"testing"

	"github.com/eser/go-service/pkg/bliss/lib"
)

func TestStringsTrimLeadingSpaceFromBytes(t *testing.T) {
	tests := []struct {
		name string
		src  []byte
		want []byte
	}{
		{
			name: "Empty input",
			src:  []byte{},
			want: []byte{},
		},
		{
			name: "No leading spaces",
			src:  []byte("Hello, World!"),
			want: []byte("Hello, World!"),
		},
		{
			name: "Leading spaces",
			src:  []byte("   Hello, World!"),
			want: []byte("Hello, World!"),
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lib.StringsTrimLeadingSpaceFromBytes(tt.src)
			if string(got) != string(tt.want) {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestStringsTrimTrailingSpaceFromBytes(t *testing.T) {
	tests := []struct {
		name string
		src  []byte
		want []byte
	}{
		{
			name: "Empty input",
			src:  []byte{},
			want: []byte{},
		},
		{
			name: "No trailing spaces",
			src:  []byte("Hello, World!"),
			want: []byte("Hello, World!"),
		},
		{
			name: "Trailing spaces",
			src:  []byte("Hello, World!   "),
			want: []byte("Hello, World!"),
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lib.StringsTrimTrailingSpaceFromBytes(tt.src)
			if string(got) != string(tt.want) {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestStringsTrimLeadingSpace(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "Empty input",
			src:  "",
			want: "",
		},
		{
			name: "No leading spaces",
			src:  "Hello, World!",
			want: "Hello, World!",
		},
		{
			name: "Leading spaces",
			src:  "   Hello, World!",
			want: "Hello, World!",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lib.StringsTrimLeadingSpace(tt.src)
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestStringsTrimTrailingSpace(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "Empty input",
			src:  "",
			want: "",
		},
		{
			name: "No trailing spaces",
			src:  "Hello, World!",
			want: "Hello, World!",
		},
		{
			name: "Trailing spaces",
			src:  "Hello, World!   ",
			want: "Hello, World!",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lib.StringsTrimTrailingSpace(tt.src)
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}
