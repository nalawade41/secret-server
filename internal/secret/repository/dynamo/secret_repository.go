package dynamo

import (
	"github.com/nalawade41/secret-server/internal/common/repository"
	"github.com/nalawade41/secret-server/internal/domain"
)

type SecretManagerRepository struct {
	repository.BaseRepository
}

var _ domain.SecretRepository = (*SecretManagerRepository)(nil)
