package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"github.com/pkg/errors"
	"io"
)

type RealEncryptor struct{}

// deriveKeyFromHash derives a secret key from the hash
func deriveKeyFromHash(hash string) ([]byte, error) {
	if len(hash) < 32 {
		return nil, errors.New("hash must be at least 32 characters long")
	}
	// Use the first 32 characters of the hash as the key
	key, _ := hex.DecodeString(hash[:32])
	return key, nil
}

// EncryptMessage encrypts the plaintext message using AES encryption
func (e RealEncryptor) EncryptMessage(plaintext string, hash string) (string, error) {
	key, err := deriveKeyFromHash(hash)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "failed to create new cipher")
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.Wrap(err, "failed to read random bytes")
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return hex.EncodeToString(ciphertext), nil
}
