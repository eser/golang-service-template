package uris_test

import (
	"testing"

	"github.com/eser/go-service/pkg/bliss/httpfx/uris"
)

func TestIsValidMethod(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		expected bool
	}{
		{
			name:     "ValidMethod",
			method:   "GET",
			expected: true,
		},
		{
			name:     "ValidMethodWithLowerCase",
			method:   "post",
			expected: true,
		},
		{
			name:     "InvalidMethodWithSpace",
			method:   "GET ",
			expected: false,
		},
		{
			name:     "InvalidMethodWithSpecialCharacter",
			method:   "GET@",
			expected: false,
		},
		{
			name:     "EmptyMethod",
			method:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uris.IsValidMethod(tt.method)
			if result != tt.expected {
				t.Errorf("IsValidMethod() = %v, want %v", result, tt.expected)
			}
		})
	}
}
