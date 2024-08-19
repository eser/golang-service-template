package lib_test

import (
	"reflect"
	"testing"

	"github.com/eser/go-service/pkg/bliss/lib"
)

func TestArraysCopy(t *testing.T) {
	tests := []struct {
		name  string
		items [][]int
		want  []int
	}{
		{
			name:  "Empty slices",
			items: [][]int{},
			want:  []int{},
		},
		{
			name:  "Single slice",
			items: [][]int{{1, 2, 3}},
			want:  []int{1, 2, 3},
		},
		{
			name:  "Multiple slices",
			items: [][]int{{1, 2}, {3, 4, 5}, {6}},
			want:  []int{1, 2, 3, 4, 5, 6},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lib.ArraysCopy(tt.items...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
