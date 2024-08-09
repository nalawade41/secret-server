package main

import (
	"context"
	"errors"
	"github.com/nalawade41/secret-server/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nalawade41/secret-server/config"
	"github.com/nalawade41/secret-server/internal/util/logger"
	"github.com/nalawade41/secret-server/server"
)

func main() {
	var cfg *config.Config
	var err error

	// Initialize the configuration and get the configuration object
	if cfg, err = config.Init(); err != nil {
		logger.Error(err)
		return
	}

	// Initialize the server with the configuration object and the router handler
	srv := server.NewServer(cfg, router.NewHandler(cfg).Init())

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
