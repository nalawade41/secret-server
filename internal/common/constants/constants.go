package constants

import "os"

const (
	AwsRegion      = "AWS_REGION"
	LambdaTaskRoot = "LAMBDA_TASK_ROOT" // this is used to check if it is local environment or lambda environment
)

var IsLambdaDeployed = os.Getenv(LambdaTaskRoot) != ""
