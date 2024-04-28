package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"os"
)

func CreateCipherKey() (cipher.Block, error) {
	// Get the ENCRYPTION_KEY from the environment variables
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		return nil, errors.New("ENCRYPTION_KEY is not set")
	}

	// Create a new SHA-256 hash
	hasher := sha256.New()
	hasher.Write([]byte(encryptionKey))

	// Get the hashed key
	key := hasher.Sum(nil)

	return aes.NewCipher(key)
}
