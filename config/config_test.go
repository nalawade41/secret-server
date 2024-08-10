package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestInit_LocalEnvironment(t *testing.T) {
	// Set up the environment variable for local testing
	os.Setenv("APP_ENV", EnvLocal)
	defer os.Unsetenv("APP_ENV")

	// Load the local .env file if it exists
	envFilePath := filepath.Join(ProjectRootPath, ".env.local")
	_ = godotenv.Load(envFilePath)

	config, err := Init()
	assert.NoError(t, err)
	assert.Equal(t, EnvLocal, config.Environment)

	// Check that the HTTP config was loaded
	assert.NotNil(t, config.HTTP)
	assert.NotEmpty(t, config.HTTP.Port)

	// Check that the Database config was loaded
	assert.NotNil(t, config.Database)
	assert.NotEmpty(t, config.Database.Host)

	// Check that the AWS config was loaded
	assert.NotNil(t, config.AWS)
	assert.NotEmpty(t, config.AWS.Region)
}

func TestInit_ProdEnvironment(t *testing.T) {
	// Set up the environment variable for production testing
	os.Setenv("APP_ENV", Prod)
	defer os.Unsetenv("APP_ENV")

	config, err := Init()
	assert.NoError(t, err)
	assert.Equal(t, Prod, config.Environment)

	// Check that the HTTP config was loaded
	assert.NotNil(t, config.HTTP)

	// Check that the Database config was loaded
	assert.NotNil(t, config.Database)

	// Check that the AWS config was loaded
	assert.NotNil(t, config.AWS)
}

func TestApplicationStartup(t *testing.T) {
	// Set up environment variables
	os.Setenv("APP_ENV", "local")
	os.Setenv("HTTP_HOST", "localhost")
	os.Setenv("HTTP_PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "8000")
	os.Setenv("DB_TABLE_NAME", "Secrets")
	os.Setenv("AWS_REGION", "us-west-2")
	os.Setenv("AWS_PROFILE", "default")

	defer func() {
		// Unset environment variables after the test
		os.Unsetenv("APP_ENV")
		os.Unsetenv("HTTP_HOST")
		os.Unsetenv("HTTP_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_TABLE_NAME")
		os.Unsetenv("AWS_REGION")
		os.Unsetenv("AWS_PROFILE")
	}()

	// Initialize configuration
	config, err := Init()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Verify that configuration values are correctly loaded
	assert.Equal(t, "localhost", config.HTTP.Host)
	assert.Equal(t, "8080", config.HTTP.Port)
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, "8000", config.Database.Port)
	assert.Equal(t, "Secrets", config.Database.TableName)
	assert.Equal(t, "us-west-2", config.AWS.Region)
	assert.Equal(t, "default", config.AWS.Profile)
}
