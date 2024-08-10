package response

import (
	"github.com/nalawade41/secret-server/internal/domain"
	"time"
)

type SecretResponse struct {
	Hash           string    `xml:"hash" json:"hash"`
	SecretText     string    `xml:"secretText" json:"secretText"`
	CreatedAt      time.Time `xml:"createdAt" json:"createdAt"`
	ExpiresAt      time.Time `xml:"expiresAt" json:"expiresAt"`
	RemainingViews int       `xml:"remainingViews" json:"remainingViews"`
}

// NewSecretResponse converts data to SecretResponse
func NewSecretResponse(data domain.Secret) SecretResponse {
	return SecretResponse{
		Hash:           data.Hash,
		SecretText:     data.SecretText,
		CreatedAt:      data.CreatedAt,
		ExpiresAt:      data.ExpiresAt,
		RemainingViews: data.RemainingViews,
	}
}
