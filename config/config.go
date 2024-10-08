package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

const (
	EnvLocal = "local"
	Prod     = "prod"
)

var (
	// Get current file full path from runtime
	_, b, _, _ = runtime.Caller(0)

	// ProjectRootPath Root folder of this project
	ProjectRootPath = filepath.Join(filepath.Dir(b), "../")
)

type (
	Config struct {
		Environment string
		HTTP        *HttpConfig
		Database    *DynamoConfig
		AWS         *AWSConfig
	}
)

// Init populates Config struct with values from a config file
// located at filepath and environment variables.
func Init() (*Config, error) {
	// Get the environment
	env := os.Getenv("APP_ENV")

	// Load .env only if APP_ENV is "local"
	if env == EnvLocal {
		// Load environment-specific .env file if it exists
		envFilePath := ProjectRootPath + "/.env." + env
		if _, err := os.Stat(envFilePath); err == nil {
			if err := godotenv.Load(envFilePath); err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("Error loading enviornment %s file", envFilePath))
			}
		}
	}

	http := LoadHttpConfig()
	db := LoadDynamoConfig()
	aws := LoadAWSConfig()

	config := &Config{
		Environment: env,
		HTTP:        http,
		Database:    db,
		AWS:         aws,
	}
	return config, nil
}
