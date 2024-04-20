package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

func generateIV() []byte {
	iv := make([]byte, aes.BlockSize)
	_, ioError := io.ReadFull(rand.Reader, iv)

	if ioError != nil {
		return []byte{}
	}

	return iv
}

func Encrypt(secret string) (string, error) {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	plaintext := []byte(secret)

	block, blockError := aes.NewCipher(key)
	if blockError != nil {
		return "", blockError
	}

	iv := generateIV()
	if len(iv) == 0 {
		return "", errors.New("unable to generate an IV byte array")
	}

	encrypter := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	encrypter.XORKeyStream(ciphertext, plaintext)

	// Prepend the IV to the ciphertext.
	result := append(iv, ciphertext...)

	return base64.URLEncoding.EncodeToString(result), nil
}
