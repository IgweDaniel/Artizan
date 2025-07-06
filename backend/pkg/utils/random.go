package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func Randomize(length int) (string, error) {
	randomBytes := make([]byte, length)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	t := hex.EncodeToString(randomBytes)

	return t, nil
}
