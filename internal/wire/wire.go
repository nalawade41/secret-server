//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/nalawade41/secret-server/db"
	"github.com/nalawade41/secret-server/internal/secret"
	"github.com/nalawade41/secret-server/internal/secret/handler"
)

func InitializeRouteProvider(dbConnection db.DynamoDBAPI, tableName string) *handler.SecretManagerHandler {
	panic(wire.Build(secret.ManagerProviderSet))
}
