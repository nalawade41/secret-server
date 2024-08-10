package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// deriveKeyFromHash derives a secret key from the hash
func deriveKeyFromHash(hash string) []byte {
	// Use the first 32 characters of the hash as the key
	key, _ := hex.DecodeString(hash[:32])
	return key
}

// EncryptMessage encrypts the plaintext message using AES encryption
func EncryptMessage(plaintext string, hash string) (string, error) {
	key := deriveKeyFromHash(hash)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return hex.EncodeToString(ciphertext), nil
}
