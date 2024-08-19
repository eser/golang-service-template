package lib_test

import (
	"testing"

	"github.com/eser/go-service/pkg/bliss/lib"
)

func TestIdsGenerateUnique(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test 1",
		},
		{
			name: "Test 2",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lib.IdsGenerateUnique()
			if len(got) != 26 {
				t.Errorf("unexpected length of generated ID, got: %d, want: 26", len(got))
			}
			// Add more assertions as needed
		})
	}
}
