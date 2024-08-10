package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nalawade41/secret-server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeriveKeyFromHash(t *testing.T) {
	// Test that the derived key is of the correct length and format
	hash := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	expectedKey, _ := hex.DecodeString(hash[:32])
	key, err := deriveKeyFromHash(hash)
	assert.NoError(t, err, "deriveKeyFromHash should not return an error for valid hash")
	assert.Equal(t, expectedKey, key, "The derived key should match the expected key")
	assert.Equal(t, 16, len(key), "The derived key should be 16 bytes long for AES-128")
}

// Mock of the security package using Encryptor interface

func TestEncryptMessage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncryptor := mocks.NewMockEncryptor(ctrl)

	plaintext := "Hello, World!"
	hash := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	// Set up mock expectations for encryption
	expectedCiphertext := "expectedciphertext"
	mockEncryptor.EXPECT().EncryptMessage(plaintext, hash).Return(expectedCiphertext, nil)

	// Encrypt the plaintext using the mock
	ciphertext, err := mockEncryptor.EncryptMessage(plaintext, hash)
	assert.NoError(t, err, "EncryptMessage should not return an error")
	assert.Equal(t, expectedCiphertext, ciphertext, "Ciphertext should match expected value")
}

func TestEncryptMessage_InvalidKeyLength(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncryptor := mocks.NewMockEncryptor(ctrl)

	plaintext := "Hello, World!"
	hash := "short"

	// Simulate encryption error for short hash
	mockEncryptor.EXPECT().EncryptMessage(plaintext, hash).Return("", errors.New("hash must be at least 32 characters long"))

	_, err := mockEncryptor.EncryptMessage(plaintext, hash)
	assert.Error(t, err, "EncryptMessage should return an error for invalid key length")
	assert.Contains(t, err.Error(), "hash must be at least 32 characters long", "Error should indicate hash length requirement")
}

func TestEncryptMessage_EmptyPlaintext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncryptor := mocks.NewMockEncryptor(ctrl)

	plaintext := ""
	hash := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	// Set up mock expectations for encryption with empty plaintext
	expectedCiphertext := "emptyciphertext"
	mockEncryptor.EXPECT().EncryptMessage(plaintext, hash).Return(expectedCiphertext, nil)

	ciphertext, err := mockEncryptor.EncryptMessage(plaintext, hash)
	assert.NoError(t, err, "EncryptMessage should not return an error for empty plaintext")
	assert.Equal(t, expectedCiphertext, ciphertext, "Ciphertext should match expected value")
}

func TestDeriveKeyFromHash_ValidHash(t *testing.T) {
	hash := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" // 64 hex characters = 32 bytes

	key, err := deriveKeyFromHash(hash)
	expectedKey, _ := hex.DecodeString(hash[:32])

	assert.NoError(t, err, "deriveKeyFromHash should not return an error for a valid hash")
	assert.Equal(t, expectedKey, key, "The derived key should match the first 32 characters of the hash")
}

func TestDeriveKeyFromHash_InvalidHash(t *testing.T) {
	hash := "short"

	key, err := deriveKeyFromHash(hash)

	assert.Error(t, err, "deriveKeyFromHash should return an error for an invalid hash")
	assert.Nil(t, key, "The derived key should be nil for an invalid hash")
	assert.Contains(t, err.Error(), "hash must be at least 32 characters long", "The error message should indicate the hash length requirement")
}

func TestEncryptMessage_Success_Actual(t *testing.T) {
	encryptor := RealEncryptor{}

	plaintext := "Hello, World!"
	hash := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" // 64 hex characters = 32 bytes

	// Encrypt the plaintext
	ciphertext, err := encryptor.EncryptMessage(plaintext, hash)
	assert.NoError(t, err, "EncryptMessage should not return an error")

	// Decode ciphertext from hex
	ciphertextBytes, err := hex.DecodeString(ciphertext)
	assert.NoError(t, err, "Decoding ciphertext should not return an error")

	// Derive key and check the decryption
	key, err := deriveKeyFromHash(hash)
	assert.NoError(t, err, "deriveKeyFromHash should not return an error for a valid hash")

	block, err := aes.NewCipher(key)
	assert.NoError(t, err, "Creating cipher block should not return an error")

	iv := ciphertextBytes[:aes.BlockSize]
	encryptedMessage := ciphertextBytes[aes.BlockSize:]
	decryptedMessage := make([]byte, len(encryptedMessage))

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decryptedMessage, encryptedMessage)

	assert.Equal(t, plaintext, string(decryptedMessage), "Decrypted message should match the original plaintext")
}

func TestEncryptMessage_InvalidKeyLength_Actual(t *testing.T) {
	encryptor := RealEncryptor{}

	plaintext := "Hello, World!"
	hash := "short"

	_, err := encryptor.EncryptMessage(plaintext, hash)

	assert.Error(t, err, "EncryptMessage should return an error for invalid key length")
	assert.Contains(t, err.Error(), "hash must be at least 32 characters long", "The error message should indicate the hash length requirement")
}

func TestEncryptMessage_EmptyPlaintext_Actual(t *testing.T) {
	encryptor := RealEncryptor{}

	plaintext := ""
	hash := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" // 64 hex characters = 32 bytes

	ciphertext, err := encryptor.EncryptMessage(plaintext, hash)
	assert.NoError(t, err, "EncryptMessage should not return an error for empty plaintext")
	assert.NotEmpty(t, ciphertext, "Ciphertext should not be empty even if plaintext is empty")

	// Decode ciphertext from hex
	ciphertextBytes, err := hex.DecodeString(ciphertext)
	assert.NoError(t, err, "Decoding ciphertext should not return an error")

	// Derive key and check the decryption
	key, err := deriveKeyFromHash(hash)
	assert.NoError(t, err, "deriveKeyFromHash should not return an error for a valid hash")

	block, err := aes.NewCipher(key)
	assert.NoError(t, err, "Creating cipher block should not return an error")

	iv := ciphertextBytes[:aes.BlockSize]
	encryptedMessage := ciphertextBytes[aes.BlockSize:]
	decryptedMessage := make([]byte, len(encryptedMessage))

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decryptedMessage, encryptedMessage)

	assert.Equal(t, plaintext, string(decryptedMessage), "Decrypted message should match the original plaintext")
}
