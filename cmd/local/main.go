package main

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/nalawade41/secret-server/db"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nalawade41/secret-server/config"
	_ "github.com/nalawade41/secret-server/docs"
	"github.com/nalawade41/secret-server/internal/common/logger"
	"github.com/nalawade41/secret-server/router"
	"github.com/nalawade41/secret-server/server"
)

//	@title			My API
//	@version		1.0
//	@description	This is a sample server for a secret API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/
//	@schemes	http
func main() {
	var err error

	// Initialize the configuration and get the configuration object
	var cfg *config.Config
	if cfg, err = config.Init(); err != nil {
		logger.Error(err)
		return
	}

	// Initialize the dynamo client
	var dbConnect *dynamodb.Client
	if dbConnect, err = db.InitDynamoDB(cfg); err != nil {
		logger.Error(err)
		return
	}

	// Initialize the server with the configuration object and the router handler
	srv := server.NewServer(cfg, router.NewHandler(cfg, dbConnect).Init())

	// Start the server in a goroutine
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	// Log the server address to the logger
	logger.Infof("Server started and listening on port %s", srv.GetAddress())

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	logger.Info("Server stopped")

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}
