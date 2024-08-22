package lib_test

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/eser/go-service/pkg/bliss/lib"
)

// MockReader is a mock implementation of io.Reader that simulates
// a failure in the Read method.
type MockReader struct {
	fail bool
}

// Read implements the io.Reader interface for MockReader.
func (m *MockReader) Read(p []byte) (n int, err error) {
	if m.fail {
		return 0, fmt.Errorf("mock read error")
	}

	// Simulate successful read
	for i := range p {
		p[i] = byte(i)
	}

	return len(p), nil
}

func TestCryptoGetRandomBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockReader    *MockReader
		expectedError bool
	}{
		{
			name:          "Successful read",
			mockReader:    &MockReader{fail: false},
			expectedError: false,
		},
		{
			name:          "Read error",
			mockReader:    &MockReader{fail: true},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			originalRand := rand.Reader
			defer func() {
				rand.Reader = originalRand
			}() // Restore original rand.Reader

			if tt.mockReader.fail {
				rand.Reader = tt.mockReader
			}

			result, err := lib.CryptoGetRandomBytes(16)

			if (err != nil) != tt.expectedError {
				t.Errorf("CryptoGetRandomBytes() error = %v, wantErr %v", err, tt.expectedError)

				return
			}

			if !tt.expectedError && len(result) != 16 {
				t.Errorf("CryptoGetRandomBytes() = %v, want length %v", result, 16)
			}
		})
	}
}
