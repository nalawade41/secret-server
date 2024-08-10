package response

import (
	"testing"
	"time"

	"github.com/nalawade41/secret-server/internal/domain"
	"github.com/stretchr/testify/assert"
)

// TestNewSecretResponse tests the conversion from domain.Secret to SecretResponse
func TestNewSecretResponse(t *testing.T) {
	createdAt := time.Now().UTC()
	expiresAt := createdAt.Add(10 * time.Minute)

	domainSecret := domain.Secret{
		Hash:           "testhash",
		SecretText:     "This is a test secret",
		CreatedAt:      createdAt,
		ExpiresAt:      expiresAt,
		RemainingViews: 5,
	}

	expectedResponse := SecretResponse{
		Hash:           domainSecret.Hash,
		SecretText:     domainSecret.SecretText,
		CreatedAt:      domainSecret.CreatedAt,
		ExpiresAt:      domainSecret.ExpiresAt,
		RemainingViews: domainSecret.RemainingViews,
	}

	response := NewSecretResponse(domainSecret)

	assert.Equal(t, expectedResponse.Hash, response.Hash, "Hash should match")
	assert.Equal(t, expectedResponse.SecretText, response.SecretText, "SecretText should match")
	assert.Equal(t, expectedResponse.CreatedAt, response.CreatedAt, "CreatedAt should match")
	assert.Equal(t, expectedResponse.ExpiresAt, response.ExpiresAt, "ExpiresAt should match")
	assert.Equal(t, expectedResponse.RemainingViews, response.RemainingViews, "RemainingViews should match")
}
