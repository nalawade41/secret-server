package secret

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/wire"
	"github.com/nalawade41/secret-server/internal/common/repository"
	"github.com/nalawade41/secret-server/internal/domain"
	"github.com/nalawade41/secret-server/internal/secret/handler"
	"github.com/nalawade41/secret-server/internal/secret/repository/dynamo"
	"github.com/nalawade41/secret-server/internal/secret/usecase"
)

var (
	secretHandler *handler.SecretManagerHandler
	hdlOnce       sync.Once

	secretUseCase *usecase.SecretManagerUseCase
	ucOnce        sync.Once

	repo     *dynamo.SecretManagerRepository
	repoOnce sync.Once

	ManagerProviderSet wire.ProviderSet = wire.NewSet(
		NewSecretManagerHandler,
		NewSecretManagerUseCase,
		NewSecretManagerRepository,

		wire.Bind(new(domain.SecretUseCase), new(*usecase.SecretManagerUseCase)),
		wire.Bind(new(domain.SecretRepository), new(*dynamo.SecretManagerRepository)),
	)
)

func NewSecretManagerUseCase(repo domain.SecretRepository) *usecase.SecretManagerUseCase {
	ucOnce.Do(func() {
		secretUseCase = &usecase.SecretManagerUseCase{
			SecretRepo: repo,
		}
	})
	return secretUseCase
}

func NewSecretManagerHandler(rs domain.SecretUseCase) *handler.SecretManagerHandler {
	hdlOnce.Do(func() {
		secretHandler = &handler.SecretManagerHandler{
			SecretManager: rs,
		}
	})
	return secretHandler
}

// NewSecretManagerRepository creates new secret repository
func NewSecretManagerRepository(db *dynamodb.Client, tableName string) *dynamo.SecretManagerRepository {
	repoOnce.Do(func() {
		repo = &dynamo.SecretManagerRepository{
			BaseRepository: repository.BaseRepository{
				DBConnection: db,
			},
			TableName: tableName,
		}
	})
	return repo
}
