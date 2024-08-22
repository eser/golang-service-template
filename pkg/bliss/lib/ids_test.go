package lib_test

import (
	"testing"

	"github.com/eser/go-service/pkg/bliss/lib"
	"github.com/stretchr/testify/assert"
)

func TestIdsGenerateUnique(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "Test 1",
		},
		{
			name: "Test 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := lib.IdsGenerateUnique()

			assert.Len(t, got, 26)
		})
	}
}
