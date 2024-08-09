package repository

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type BaseRepository struct {
	DBConnection *dynamodb.Client
}

// If we need, we can add the common methods here like transaction handlers and other common methods
