package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nalawade41/secret-server/internal/domain"
	"github.com/nalawade41/secret-server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateSecretMessage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSecretRepository(ctrl)
	mockEncryptor := mocks.NewMockEncryptor(ctrl)

	useCase := SecretManagerUseCase{SecretRepo: mockRepo, Encryptor: mockEncryptor}

	message := domain.Secret{
		SecretText:     "This is a test secret",
		ExpiresAt:      time.Now().Add(10 * time.Minute),
		RemainingViews: 5,
		CreatedAt:      time.Now().UTC(),
	}

	// Mocked hash and encrypted text
	expectedHash := "mockedhash"
	expectedEncryptedText := "encryptedText"

	// Set expectations for mock methods
	mockEncryptor.EXPECT().GenerateSHA256Hash(message.SecretText, message.CreatedAt.String()).Return(expectedHash)
	mockEncryptor.EXPECT().EncryptMessage(message.SecretText, expectedHash).Return(expectedEncryptedText, nil)

	// Set expectations for mock repository
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

	result, err := useCase.CreateSecretMessage(context.Background(), message)

	assert.NoError(t, err)
	assert.Equal(t, expectedHash, result.Hash)
	assert.Equal(t, expectedEncryptedText, result.SecretText)
	assert.Equal(t, message.RemainingViews, result.RemainingViews)
}

func TestCreateSecretMessage_EncryptionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSecretRepository(ctrl)
	mockEncryptor := mocks.NewMockEncryptor(ctrl)

	useCase := SecretManagerUseCase{SecretRepo: mockRepo, Encryptor: mockEncryptor}

	message := domain.Secret{
		SecretText:     "This is a test secret",
		ExpiresAt:      time.Now().Add(10 * time.Minute),
		RemainingViews: 5,
		CreatedAt:      time.Now().UTC(),
	}

	expectedHash := "mockedhash"

	// Mock the hash generation
	mockEncryptor.EXPECT().GenerateSHA256Hash(message.SecretText, message.CreatedAt.String()).Return(expectedHash)

	// Simulate encryption error
	mockEncryptor.EXPECT().EncryptMessage(message.SecretText, expectedHash).Return("", errors.New("encryption failed"))

	_, err := useCase.CreateSecretMessage(context.Background(), message)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to encrypt secret")
}

func TestCreateSecretMessage_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSecretRepository(ctrl)
	mockEncryptor := mocks.NewMockEncryptor(ctrl)

	useCase := SecretManagerUseCase{SecretRepo: mockRepo, Encryptor: mockEncryptor}

	message := domain.Secret{
		SecretText:     "This is a test secret",
		ExpiresAt:      time.Now().Add(10 * time.Minute),
		RemainingViews: 5,
		CreatedAt:      time.Now().UTC(),
	}

	expectedHash := "mockedhash"
	expectedEncryptedText := "encryptedText"

	// Mock the hash generation
	mockEncryptor.EXPECT().GenerateSHA256Hash(message.SecretText, message.CreatedAt.String()).Return(expectedHash)

	// Set expectations for mock encryptor
	mockEncryptor.EXPECT().EncryptMessage(message.SecretText, expectedHash).Return(expectedEncryptedText, nil)

	// Set expectations for mock repository to return an error
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("repository save error"))

	_, err := useCase.CreateSecretMessage(context.Background(), message)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to store secret")
}

func TestGetSecretMessage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSecretRepository(ctrl)
	mockEncryptor := mocks.NewMockEncryptor(ctrl)

	useCase := SecretManagerUseCase{SecretRepo: mockRepo, Encryptor: mockEncryptor}

	hash := "testhash"
	secret := domain.Secret{
		Hash:           hash,
		SecretText:     "Encrypted text",
		ExpiresAt:      time.Now().Add(10 * time.Minute),
		RemainingViews: 5,
		CreatedAt:      time.Now().UTC(),
	}

	// Set expectations for mock repository
	mockRepo.EXPECT().GetByHash(gomock.Any(), hash).Return(secret, nil)
	mockRepo.EXPECT().UpdateSecretViews(gomock.Any(), hash, 4).Return(nil)

	result, err := useCase.GetSecretMessage(context.Background(), hash)

	assert.NoError(t, err)
	assert.Equal(t, 4, result.RemainingViews)
}

func TestGetSecretMessage_SecretExpired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSecretRepository(ctrl)
	mockEncryptor := mocks.NewMockEncryptor(ctrl)

	useCase := SecretManagerUseCase{SecretRepo: mockRepo, Encryptor: mockEncryptor}

	hash := "testhash"
	secret := domain.Secret{
		Hash:           hash,
		SecretText:     "Encrypted text",
		ExpiresAt:      time.Now().Add(-10 * time.Minute), // Already expired
		RemainingViews: 5,
		CreatedAt:      time.Now().UTC(),
	}

	// Set expectations for mock repository
	mockRepo.EXPECT().GetByHash(gomock.Any(), hash).Return(secret, nil)
	mockRepo.EXPECT().DeleteSecret(gomock.Any(), hash).Return(nil)

	_, err := useCase.GetSecretMessage(context.Background(), hash)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "secret expired or no remaining views")
}

func TestGetSecretMessage_NoRemainingViews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSecretRepository(ctrl)
	mockEncryptor := mocks.NewMockEncryptor(ctrl)

	useCase := SecretManagerUseCase{SecretRepo: mockRepo, Encryptor: mockEncryptor}

	hash := "testhash"
	secret := domain.Secret{
		Hash:           hash,
		SecretText:     "Encrypted text",
		ExpiresAt:      time.Now().Add(10 * time.Minute),
		RemainingViews: 0, // No remaining views
		CreatedAt:      time.Now().UTC(),
	}

	// Set expectations for mock repository
	mockRepo.EXPECT().GetByHash(gomock.Any(), hash).Return(secret, nil)
	mockRepo.EXPECT().DeleteSecret(gomock.Any(), hash).Return(nil)

	_, err := useCase.GetSecretMessage(context.Background(), hash)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "secret expired or no remaining views")
}

func TestGetSecretMessage_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSecretRepository(ctrl)
	mockEncryptor := mocks.NewMockEncryptor(ctrl)

	useCase := SecretManagerUseCase{SecretRepo: mockRepo, Encryptor: mockEncryptor}

	hash := "testhash"

	// Set expectations for mock repository to return an error
	mockRepo.EXPECT().GetByHash(gomock.Any(), hash).Return(domain.Secret{}, errors.New("repository error"))

	_, err := useCase.GetSecretMessage(context.Background(), hash)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to retrieve secret")
}
