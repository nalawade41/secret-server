package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAWSConfig_ValidEnvVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("AWS_REGION", "us-west-2")
	os.Setenv("AWS_PROFILE", "default")

	defer func() {
		// Unset environment variables after the test
		os.Unsetenv("AWS_REGION")
		os.Unsetenv("AWS_PROFILE")
	}()

	// Load AWS config
	awsConfig := LoadAWSConfig()

	// Assertions
	assert.Equal(t, "us-west-2", awsConfig.Region)
	assert.Equal(t, "default", awsConfig.Profile)
}

func TestLoadAWSConfig_MissingEnvVariables(t *testing.T) {
	// Ensure environment variables are not set
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_PROFILE")

	// Load AWS config
	awsConfig := LoadAWSConfig()

	// Assertions
	assert.Empty(t, awsConfig.Region)
	assert.Empty(t, awsConfig.Profile)
}
