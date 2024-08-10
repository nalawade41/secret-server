package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nalawade41/secret-server/internal/common/security"
	"github.com/nalawade41/secret-server/internal/domain"
)

type SecretManagerUseCase struct {
	SecretRepo domain.SecretRepository
}

// CreateSecretMessage creates a secret message and stores it in the repository
func (s SecretManagerUseCase) CreateSecretMessage(ctx context.Context, message domain.Secret) (domain.Secret, error) {
	// Generate a unique hash for the secret
	hash := security.GenerateSHA256Hash(message.SecretText, message.CreatedAt.String())

	// Set the hash in the secret
	message.Hash = hash

	// Encrypt the message
	encryptedText, err := security.EncryptMessage(message.SecretText, hash)
	if err != nil {
		return domain.Secret{}, fmt.Errorf("failed to encrypt secret: %w", err)
	}
	message.SecretText = encryptedText

	// Store the secret in the repository
	if err := s.SecretRepo.Save(ctx, message); err != nil {
		return domain.Secret{}, fmt.Errorf("failed to store secret: %w", err)
	}

	return message, nil
}

// GetSecretMessage retrieves a secret from the repository and decrements the remaining views
func (s SecretManagerUseCase) GetSecretMessage(ctx context.Context, hash string) (domain.Secret, error) {
	// Retrieve the secret from the repository
	secret, err := s.SecretRepo.GetByHash(ctx, hash)
	if err != nil {
		return domain.Secret{}, fmt.Errorf("failed to retrieve secret: %w", err)
	}

	// Check if the secret has expired or if there are no remaining views
	if secret.ExpiresAt.Before(time.Now().UTC()) || secret.RemainingViews <= 0 {
		// TODO:This part we can do asynchronously using queue services like SQS, RabbitMQ, etc.
		// Delete the secret from the repository
		if err := s.SecretRepo.DeleteSecret(ctx, hash); err != nil {
			return domain.Secret{}, fmt.Errorf("failed to delete expired or fully viewed secret: %w", err)
		}
		return domain.Secret{}, errors.New("secret expired or no remaining views")
	}

	// Decrement the remaining views
	secret.RemainingViews -= 1

	// If the views reach 0 after decrementing, delete the secret
	if secret.RemainingViews == 0 {
		// TODO:This part we can do asynchronously using queue services like SQS, RabbitMQ, etc.
		if err := s.SecretRepo.DeleteSecret(ctx, hash); err != nil {
			return domain.Secret{}, fmt.Errorf("failed to delete secret after reaching 0 views: %w", err)
		}
		return secret, nil
	}

	// TODO:This part we can do asynchronously using queue services like SQS, RabbitMQ, etc.
	// Update the remaining views in the repository
	err = s.SecretRepo.UpdateSecretViews(ctx, hash, secret.RemainingViews)
	if err != nil {
		return domain.Secret{}, fmt.Errorf("failed to update remaining views: %w", err)
	}

	return secret, nil
}

var _ domain.SecretUseCase = (*SecretManagerUseCase)(nil)
