//go:build wireinject

package wire

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/wire"
	"github.com/nalawade41/secret-server/internal/secret"
	"github.com/nalawade41/secret-server/internal/secret/handler"
)

func InitializeRouteProvider(dbConnection *dynamodb.Client, tableName string) *handler.SecretManagerHandler {
	panic(wire.Build(secret.ManagerProviderSet))
}
