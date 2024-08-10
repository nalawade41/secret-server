package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/nalawade41/secret-server/config"
	"github.com/nalawade41/secret-server/db"
	_ "github.com/nalawade41/secret-server/docs"
	"github.com/nalawade41/secret-server/internal/common/logger"
	"github.com/nalawade41/secret-server/router"
	"github.com/nalawade41/secret-server/trace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
)

var echoLambda *echoadapter.EchoLambda

func init() {
	logger.Info("Initializing the Lambda function")

	cfg, err := config.Init()
	if err != nil {
		logger.Error(err)
		return
	}

	// Initialize the dynamo client
	var dbConnect db.DynamoDBAPI
	if dbConnect, err = db.InitDynamoDB(cfg); err != nil {
		logger.Error(err)
		return
	}

	// Initialize the server with the configuration object and the router handler
	echoLambda = echoadapter.New(router.NewHandler(cfg, dbConnect).Init())
}

// @title My API
// @version 1.0
// @description This is a sample server for a secret API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes	http
func main() {
	ctx, tp, err := trace.SetTracing()
	if err != nil {
		logger.Errorf("error setting up tracing: %v", err)
		return
	}

	defer func(ctx context.Context) {
		err := tp.Shutdown(ctx)
		if err != nil {
			logger.Infof("error shutting down tracer provider: %v", err)
		}
	}(ctx)

	lambda.Start(otellambda.InstrumentHandler(Handler, xrayconfig.WithRecommendedOptions(tp)...))

	logger.Info("Lambda started")
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Create handlers for routing
	// Create Echo Router with routes
	// Wrap the Echo router using the Echo Lambda adapter
	// Forward the Lambda request to the Echo router
	response, err := echoLambda.ProxyWithContext(ctx, req)
	if err != nil {
		logger.Errorf(fmt.Sprintf("Error while processing the Lambda request: %v", err.Error()))
		return events.APIGatewayProxyResponse{}, err
	}
	return response, nil
}
