package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/nalawade41/secret-server/internal/common/logger"
	"github.com/nalawade41/secret-server/internal/domain"
	"github.com/pkg/errors"
)

type SecretManagerUseCase struct {
	SecretRepo domain.SecretRepository
	Encryptor  domain.Encryptor
}

// CreateSecretMessage creates a secret message and stores it in the repository
func (s SecretManagerUseCase) CreateSecretMessage(ctx context.Context, message domain.Secret) (domain.Secret, error) {
	// Generate a unique hash for the secret
	hash := s.Encryptor.GenerateSHA256Hash(message.SecretText, message.CreatedAt.String())

	// Set the hash in the secret
	message.Hash = hash

	// Encrypt the message
	encryptedText, err := s.Encryptor.EncryptMessage(message.SecretText, hash)
	if err != nil {
		return domain.Secret{}, errors.Wrap(err, fmt.Sprintf("failed to encrypt secret: %v", err))
	}

	message.SecretText = encryptedText

	// Store the secret in the repository
	if err := s.SecretRepo.Save(ctx, message); err != nil {
		return domain.Secret{}, errors.Wrap(err, fmt.Sprintf("failed to store secret: %v", err))
	}

	return message, nil
}

// GetSecretMessage retrieves a secret from the repository and decrements the remaining views
func (s SecretManagerUseCase) GetSecretMessage(ctx context.Context, hash string) (domain.Secret, error) {
	// Retrieve the secret from the repository
	secret, err := s.SecretRepo.GetByHash(ctx, hash)
	if err != nil {
		return domain.Secret{}, errors.Wrap(err, fmt.Sprintf("failed to retrieve secret: %v", err))
	}

	// Check if the secret has expired or if there are no remaining views
	if secret.ExpiresAt.Before(time.Now().UTC()) || secret.RemainingViews <= 0 {
		// TODO:This part we can do asynchronously using queue services like SQS, RabbitMQ, etc.
		// Delete the secret from the repository
		if err := s.SecretRepo.DeleteSecret(ctx, hash); err != nil {
			return domain.Secret{}, errors.Wrap(err, fmt.Sprintf("failed to delete expired or fully viewed secret: %v", err))
		}
		return domain.Secret{}, errors.New("secret expired or no remaining views")
	}

	// Decrement the remaining views
	secret.RemainingViews -= 1

	// If the views reach 0 after decrementing, delete the secret
	if secret.RemainingViews == 0 {
		// TODO:This part we can do asynchronously using queue services like SQS, RabbitMQ, etc.
		if err := s.SecretRepo.DeleteSecret(ctx, hash); err != nil {
			logger.Error("failed to delete secret: %v", err)
		}
		return secret, nil
	}

	// TODO:This part we can do asynchronously using queue services like SQS, RabbitMQ, etc.
	// Update the remaining views in the repository
	err = s.SecretRepo.UpdateSecretViews(ctx, hash, secret.RemainingViews)
	if err != nil {
		logger.Errorf("failed to update remaining views: %w", err)
	}

	return secret, nil
}

var _ domain.SecretUseCase = (*SecretManagerUseCase)(nil)
