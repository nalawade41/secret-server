package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nalawade41/secret-server/config"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	cfg := &config.Config{
		HTTP: &config.HttpConfig{
			Port:               "8080",
			ReadTimeout:        10 * time.Second,
			WriteTimeout:       10 * time.Second,
			MaxHeaderMegabytes: 1,
		},
	}

	handler := http.NewServeMux()
	server := NewServer(cfg, handler)

	assert.NotNil(t, server)
	assert.Equal(t, ":8080", server.GetAddress())
	assert.Equal(t, handler, server.httpServer.Handler)
	assert.Equal(t, cfg.HTTP.ReadTimeout, server.httpServer.ReadTimeout)
	assert.Equal(t, cfg.HTTP.WriteTimeout, server.httpServer.WriteTimeout)
	assert.Equal(t, cfg.HTTP.MaxHeaderMegabytes<<20, server.httpServer.MaxHeaderBytes)
}

func TestServer_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock server
	handler := http.NewServeMux()
	server := &Server{
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: handler,
		},
	}

	go func() {
		time.Sleep(100 * time.Millisecond) // Wait a bit to ensure server tries to start
		_ = server.Stop(context.Background())
	}()

	err := server.Run()

	// Since we stop the server immediately, ListenAndServe should return nil or an error indicating it was closed
	if err != nil {
		assert.Contains(t, err.Error(), "Server closed")
	}
}

func TestServer_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := http.NewServeMux()
	server := &Server{
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: handler,
		},
	}

	go func() {
		// Start the server in a goroutine
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			t.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)

	// Stop the server and check for errors
	err := server.Stop(context.Background())
	assert.NoError(t, err)
}
