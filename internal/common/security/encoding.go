package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type RealEncryptor struct{}

// deriveKeyFromHash derives a secret key from the hash
func deriveKeyFromHash(hash string) ([]byte, error) {
	if len(hash) < 32 {
		return nil, errors.New("hash must be at least 32 characters long")
	}

	// Use the first 32 characters of the hash as the key
	key, err := hex.DecodeString(hash[:32])
	if err != nil {
		return nil, errors.Wrap(err, "invalid hex in key derivation")
	}
	return key, nil
}

// EncryptMessage encrypts the plaintext message using AES encryption
func (e RealEncryptor) EncryptMessage(plaintext string, hash string) (string, error) {
	key, err := deriveKeyFromHash(hash)
	if err != nil {
		return "", err
	}

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "failed to create AES cipher")
	}
	paddedText, _ := pkcs7Pad([]byte(plaintext), block.BlockSize())

	// The IV needs to be unique, but not secure.
	// Therefore, it's common to include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(paddedText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.Wrap(err, "failed to generate IV")
	}

	bm := cipher.NewCBCEncrypter(block, iv)
	bm.CryptBlocks(ciphertext[aes.BlockSize:], paddedText)

	return fmt.Sprintf("%x", ciphertext), nil
}

func pkcs7Pad(b []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, errors.New("invalid blocksize")
	}

	if b == nil || len(b) == 0 {
		return nil, errors.New("invalid PKCS7 data (empty or not padded)")
	}

	n := blockSize - (len(b) % blockSize)
	pb := make([]byte, len(b)+n)

	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))

	return pb, nil
}
