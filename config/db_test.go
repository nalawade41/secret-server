package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadDynamoConfig_ValidEnvVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "8000")
	os.Setenv("DB_TABLE_NAME", "Secrets")

	defer func() {
		// Unset environment variables after the test
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_TABLE_NAME")
	}()

	// Load DynamoDB config
	dynamoConfig := LoadDynamoConfig()

	// Assertions
	assert.Equal(t, "localhost", dynamoConfig.Host)
	assert.Equal(t, "8000", dynamoConfig.Port)
	assert.Equal(t, "Secrets", dynamoConfig.TableName)
}

func TestLoadDynamoConfig_MissingEnvVariables(t *testing.T) {
	// Ensure environment variables are not set
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_TABLE_NAME")

	// Load DynamoDB config
	dynamoConfig := LoadDynamoConfig()

	// Assertions
	assert.Empty(t, dynamoConfig.Host)
	assert.Empty(t, dynamoConfig.Port)
	assert.Empty(t, dynamoConfig.TableName)
}
