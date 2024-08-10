package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadHttpConfig_ValidEnvVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("HTTP_HOST", "localhost")
	os.Setenv("HTTP_PORT", "8080")
	os.Setenv("HTTP_PROTOCOL", "http")
	os.Setenv("READ_TIMEOUT", "10s")
	os.Setenv("WRITE_TIMEOUT", "15s")
	os.Setenv("MAX_HEADER_BYTES", "2048")

	defer func() {
		// Unset environment variables after the test
		os.Unsetenv("HTTP_HOST")
		os.Unsetenv("HTTP_PORT")
		os.Unsetenv("HTTP_PROTOCOL")
		os.Unsetenv("READ_TIMEOUT")
		os.Unsetenv("WRITE_TIMEOUT")
		os.Unsetenv("MAX_HEADER_BYTES")
	}()

	// Load HTTP config
	httpConfig := LoadHttpConfig()

	// Assertions
	assert.Equal(t, "localhost", httpConfig.Host)
	assert.Equal(t, "8080", httpConfig.Port)
	assert.Equal(t, "http", httpConfig.HttpProtocol)
	assert.Equal(t, 10*time.Second, httpConfig.ReadTimeout)
	assert.Equal(t, 15*time.Second, httpConfig.WriteTimeout)
	assert.Equal(t, 2048, httpConfig.MaxHeaderMegabytes)
}

func TestLoadHttpConfig_DefaultValues(t *testing.T) {
	// Unset environment variables to simulate defaults
	os.Unsetenv("READ_TIMEOUT")
	os.Unsetenv("WRITE_TIMEOUT")
	os.Unsetenv("MAX_HEADER_BYTES")

	// Set environment variables for other values
	os.Setenv("HTTP_HOST", "localhost")
	os.Setenv("HTTP_PORT", "8080")
	os.Setenv("HTTP_PROTOCOL", "http")

	defer func() {
		// Unset environment variables after the test
		os.Unsetenv("HTTP_HOST")
		os.Unsetenv("HTTP_PORT")
		os.Unsetenv("HTTP_PROTOCOL")
	}()

	// Load HTTP config
	httpConfig := LoadHttpConfig()

	// Assertions
	assert.Equal(t, "localhost", httpConfig.Host)
	assert.Equal(t, "8080", httpConfig.Port)
	assert.Equal(t, "http", httpConfig.HttpProtocol)
	assert.Equal(t, 5*time.Second, httpConfig.ReadTimeout)  // Default value
	assert.Equal(t, 5*time.Second, httpConfig.WriteTimeout) // Default value
	assert.Equal(t, 1048576, httpConfig.MaxHeaderMegabytes) // Default value
}

func TestLoadHttpConfig_InvalidValues(t *testing.T) {
	// Set invalid environment variables
	os.Setenv("READ_TIMEOUT", "invalid")
	os.Setenv("WRITE_TIMEOUT", "invalid")
	os.Setenv("MAX_HEADER_BYTES", "invalid")

	defer func() {
		// Unset environment variables after the test
		os.Unsetenv("READ_TIMEOUT")
		os.Unsetenv("WRITE_TIMEOUT")
		os.Unsetenv("MAX_HEADER_BYTES")
	}()

	// Set environment variables for other values
	os.Setenv("HTTP_HOST", "localhost")
	os.Setenv("HTTP_PORT", "8080")
	os.Setenv("HTTP_PROTOCOL", "http")

	defer func() {
		// Unset environment variables after the test
		os.Unsetenv("HTTP_HOST")
		os.Unsetenv("HTTP_PORT")
		os.Unsetenv("HTTP_PROTOCOL")
	}()

	// Load HTTP config
	httpConfig := LoadHttpConfig()

	// Assertions
	assert.Equal(t, 5*time.Second, httpConfig.ReadTimeout)  // Default value due to invalid input
	assert.Equal(t, 5*time.Second, httpConfig.WriteTimeout) // Default value due to invalid input
	assert.Equal(t, 1048576, httpConfig.MaxHeaderMegabytes) // Default value due to invalid input
}
