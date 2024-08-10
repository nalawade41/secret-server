package config

import (
	"os"
)

// DynamoConfig holds config details for the dynamo db server
type DynamoConfig struct {
	Host      string
	Port      string
	TableName string
}

// LoadDynamoConfig loads the DynamoConfig struct
func LoadDynamoConfig() *DynamoConfig {
	dynamoDb := DynamoConfig{
		Host:      os.Getenv("DB_HOST"),
		Port:      os.Getenv("DB_PORT"),
		TableName: os.Getenv("DB_TABLE_NAME"),
	}
	return &dynamoDb
}
