package config

import (
	"os"
)

type AWSConfig struct {
	Region  string
	Profile string
}

func LoadAWSConfig() *AWSConfig {
	aws := AWSConfig{
		Region:  os.Getenv("AWS_REGION"),
		Profile: os.Getenv("AWS_PROFILE"),
	}
	return &aws
}
