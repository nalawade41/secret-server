package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/nalawade41/secret-server/internal/domain"
)

type SecretManagerHandler struct {
	SecretManager domain.SecretUseCase
}

func (h *SecretManagerHandler) InitRoutes(e *echo.Group) {
	// TODO: Add routes
}

// TODO: Add handler methods
