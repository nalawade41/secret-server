package domain

import (
	"context"
	"time"
)

type Secret struct {
	Hash           string    `dynamodbav:"hash"`
	SecretText     string    `dynamodbav:"secretText"`
	CreatedAt      time.Time `dynamodbav:"createdAt"`
	ExpiresAt      time.Time `dynamodbav:"expiresAt"`
	RemainingViews int       `dynamodbav:"remainingViews"`
}

// SecretRepository represents interface providers for secret repository
type SecretRepository interface {
	Save(ctx context.Context, secret Secret) error
	GetByHash(ctx context.Context, hash string) (Secret, error)
	DeleteSecret(ctx context.Context, hash string) error
	UpdateSecretViews(ctx context.Context, hash string, remainingViews int) error
}

// SecretUseCase represents interface for secret use cases
type SecretUseCase interface {
	CreateSecretMessage(ctx context.Context, message Secret) (Secret, error)
	GetSecretMessage(ctx context.Context, hash string) (Secret, error)
}
