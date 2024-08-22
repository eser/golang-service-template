package lib_test

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/eser/go-service/pkg/bliss/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestCryptoGetRandomBytes(t *testing.T) { //nolint:paralleltest
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
			originalRand := rand.Reader
			defer func() {
				rand.Reader = originalRand
			}() // Restore original rand.Reader

			if tt.mockReader.fail {
				rand.Reader = tt.mockReader
			}

			const size = 16
			result, err := lib.CryptoGetRandomBytes(size)

			if tt.expectedError {
				require.Error(t, err, "CryptoGetRandomBytes() error = nil, expectedError true")

				return
			}

			require.NoError(t, err, "CryptoGetRandomBytes() error = %v, expectedError false", err)
			assert.Equal(t, size, len(result), "CryptoGetRandomBytes() = %v, want length %v", result, size)
		})
	}
}
