package repository

import (
	"github.com/nalawade41/secret-server/db"
)

type BaseRepository struct {
	DBConnection db.DynamoDBAPI
}

// If we need, we can add the common methods here like transaction handlers and other common methods
