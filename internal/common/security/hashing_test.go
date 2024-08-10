package security_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nalawade41/secret-server/internal/common/security"
	"github.com/nalawade41/secret-server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSHA256Hash_SingleInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncryptor := mocks.NewMockEncryptor(ctrl)

	input := "Hello, World!"
	expectedHash := sha256.Sum256([]byte(input))
	expectedHashString := hex.EncodeToString(expectedHash[:])

	// Set up mock expectations
	mockEncryptor.EXPECT().GenerateSHA256Hash(input).Return(expectedHashString)

	hash := mockEncryptor.GenerateSHA256Hash(input)

	assert.Equal(t, expectedHashString, hash, "The hash should match the expected SHA-256 hash of the input")
}

func TestGenerateSHA256Hash_MultipleInputs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncryptor := mocks.NewMockEncryptor(ctrl)

	inputs := []string{"Hello, ", "World!", "123"}
	concatenatedInputs := "Hello, World!123"
	expectedHash := sha256.Sum256([]byte(concatenatedInputs))
	expectedHashString := hex.EncodeToString(expectedHash[:])

	// Set up mock expectations
	mockEncryptor.EXPECT().GenerateSHA256Hash("Hello, ", "World!", "123").Return(expectedHashString)

	// Use ... to pass the slice as variadic arguments
	hash := mockEncryptor.GenerateSHA256Hash(inputs...)

	assert.Equal(t, expectedHashString, hash, "The hash should match the expected SHA-256 hash of the concatenated inputs")
}

func TestGenerateSHA256Hash_EmptyInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncryptor := mocks.NewMockEncryptor(ctrl)

	expectedHash := sha256.Sum256([]byte(""))
	expectedHashString := hex.EncodeToString(expectedHash[:])

	// Set up mock expectations
	mockEncryptor.EXPECT().GenerateSHA256Hash("").Return(expectedHashString)

	hash := mockEncryptor.GenerateSHA256Hash("")

	assert.Equal(t, expectedHashString, hash, "The hash of an empty input should match the expected SHA-256 hash for an empty string")
}

func TestGenerateSHA256Hash_IdenticalConcatenation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncryptor := mocks.NewMockEncryptor(ctrl)

	inputs1 := []string{"Hello", "World"}
	inputs2 := []string{"HelloWorld"}

	expectedHash := sha256.Sum256([]byte("HelloWorld"))
	expectedHashString := hex.EncodeToString(expectedHash[:])

	// Set up mock expectations
	mockEncryptor.EXPECT().GenerateSHA256Hash("Hello", "World").Return(expectedHashString)
	mockEncryptor.EXPECT().GenerateSHA256Hash("HelloWorld").Return(expectedHashString)

	// Use ... to pass the slice as variadic arguments
	hash1 := mockEncryptor.GenerateSHA256Hash(inputs1...)
	hash2 := mockEncryptor.GenerateSHA256Hash(inputs2...)

	assert.Equal(t, hash1, hash2, "The hash should be the same for identical concatenated inputs")
}

func TestGenerateSHA256Hash_SingleInput_Actual(t *testing.T) {
	encryptor := security.RealEncryptor{}

	input := "Hello, World!"
	expectedHash := sha256.Sum256([]byte(input))
	expectedHashString := hex.EncodeToString(expectedHash[:])

	hash := encryptor.GenerateSHA256Hash(input)

	assert.Equal(t, expectedHashString, hash, "The hash should match the expected SHA-256 hash of the input")
}

func TestGenerateSHA256Hash_MultipleInputs_Actual(t *testing.T) {
	encryptor := security.RealEncryptor{}

	inputs := []string{"Hello, ", "World!", "123"}
	concatenatedInputs := "Hello, World!123"
	expectedHash := sha256.Sum256([]byte(concatenatedInputs))
	expectedHashString := hex.EncodeToString(expectedHash[:])

	hash := encryptor.GenerateSHA256Hash(inputs...)

	assert.Equal(t, expectedHashString, hash, "The hash should match the expected SHA-256 hash of the concatenated inputs")
}

func TestGenerateSHA256Hash_EmptyInput_Actual(t *testing.T) {
	encryptor := security.RealEncryptor{}

	expectedHash := sha256.Sum256([]byte(""))
	expectedHashString := hex.EncodeToString(expectedHash[:])

	hash := encryptor.GenerateSHA256Hash("")

	assert.Equal(t, expectedHashString, hash, "The hash of an empty input should match the expected SHA-256 hash for an empty string")
}

func TestGenerateSHA256Hash_IdenticalConcatenation_Actual(t *testing.T) {
	encryptor := security.RealEncryptor{}

	inputs1 := []string{"Hello", "World"}
	inputs2 := []string{"HelloWorld"}

	hash1 := encryptor.GenerateSHA256Hash(inputs1...)
	hash2 := encryptor.GenerateSHA256Hash(inputs2...)

	assert.Equal(t, hash1, hash2, "The hash should be the same for identical concatenated inputs")
}
