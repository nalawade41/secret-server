package usecase

import "github.com/nalawade41/secret-server/internal/domain"

type SecretManagerUseCase struct {
	SecretRepo domain.SecretRepository
}

var _ domain.SecretUseCase = (*SecretManagerUseCase)(nil)
