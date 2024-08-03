package lib

import (
	"crypto/rand"
	"fmt"
)

func GetRandomBytes(size int) ([]byte, error) {
	key := make([]byte, size)

	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("error on GetRandomBytes: %w", err)
	}

	return key, nil
}
