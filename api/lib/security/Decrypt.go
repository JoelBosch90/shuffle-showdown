package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func Decrypt(encrypted string) (string, error) {
	ciphertext, base64Error := base64.URLEncoding.DecodeString(encrypted)
	if base64Error != nil {
		return "", base64Error
	}

	if len(ciphertext) <= aes.BlockSize {
		return "", nil
	}

	block, blockError := CreateCipherKey()
	if blockError != nil {
		return "", blockError
	}

	// Extract the IV from the start of the cipher.
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	decrypter := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	decrypter.XORKeyStream(plaintext, ciphertext)

	return string(plaintext), nil
}
