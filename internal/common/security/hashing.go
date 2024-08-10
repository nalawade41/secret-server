package security

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// GenerateSHA256Hash generates a SHA-256 hash from the given inputs
func GenerateSHA256Hash(inputs ...string) string {
	// Concatenate all input strings
	concatenatedInputs := strings.Join(inputs, "")

	// Create a new SHA-256 hasher
	hasher := sha256.New()

	// Write the concatenated string to the hasher
	hasher.Write([]byte(concatenatedInputs))

	// Compute the hash
	hashBytes := hasher.Sum(nil)

	// Convert the hash to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
