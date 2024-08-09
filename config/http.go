package config

import (
	"github.com/nalawade41/secret-server/internal/util/logger"
	"os"
	"strconv"
	"time"
)

type HttpConfig struct {
	Host               string
	Port               string
	HttpProtocol       string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	MaxHeaderMegabytes int
}

func LoadHttpConfig() *HttpConfig {
	// Get read timeout, write timeout, and max header bytes from environment variables
	// and parse them to a time.Duration and an int respectively
	var readTimeout, writeTimeout time.Duration
	var maxHeaderBytes int

	// Create an HttpConfig struct and populate it with values from environment variables
	http := HttpConfig{
		Host:         os.Getenv("HTTP_HOST"),
		Port:         os.Getenv("HTTP_PORT"),
		HttpProtocol: os.Getenv("HTTP_PROTOCOL"),
	}

	var err error
	if readTimeout, err = time.ParseDuration(os.Getenv("READ_TIMEOUT")); err != nil {
		logger.Warnf("Failed to parse READ_TIMEOUT: %v. Using default value", err)
		readTimeout = 5 * time.Second
	}

	if writeTimeout, err = time.ParseDuration(os.Getenv("WRITE_TIMEOUT")); err != nil {
		logger.Warnf("Failed to parse WRITE_TIMEOUT: %v. Using default value", err)
		writeTimeout = 5 * time.Second
	}

	if maxHeaderBytes, err = strconv.Atoi(os.Getenv("MAX_HEADER_BYTES")); err != nil {
		logger.Warnf("Failed to parse MAX_HEADER_BYTES: %v. Using default value", err)
		maxHeaderBytes = 1048576 // 1MB
	}

	http.ReadTimeout = readTimeout
	http.WriteTimeout = writeTimeout
	http.MaxHeaderMegabytes = maxHeaderBytes
	return &http
}
