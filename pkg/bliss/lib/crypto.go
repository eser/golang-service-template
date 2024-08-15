package lib

import (
	"crypto/rand"
	"fmt"
)

func CryptoGetRandomBytes(size int) ([]byte, error) {
	key := make([]byte, size)

	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("error on CryptoGetRandomBytes: %w", err)
	}

	return key, nil
}
