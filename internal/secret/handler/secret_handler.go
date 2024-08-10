package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nalawade41/secret-server/internal/common/responses"
	"github.com/nalawade41/secret-server/internal/domain"
	"github.com/nalawade41/secret-server/internal/secret/requests"
	"github.com/nalawade41/secret-server/internal/secret/response"
)

type SecretManagerHandler struct {
	SecretManager domain.SecretUseCase
}

func (h *SecretManagerHandler) InitRoutes(e *echo.Group) {
	e.POST("/secret", h.AddSecret)
	e.GET("/secret/:hash", h.GetSecretByHash)
}

// AddSecret godoc
//	@Summary		Add a new secret
//	@Description	Add a new secret with expiration controls
//	@Tags			secret
//	@ID				addSecret
//	@Accept			application/x-www-form-urlencoded
//	@Produce		application/json, application/xml
//	@Param			secret	body		requests.CreateSecretRequest	true	"Create Secret Message"
//	@Success		200		{object}	response.SecretResponse			"successful operation"
//	@Failure		400		{object}	responses.Error					"Bad request"
//	@Failure		405		{object}	responses.Error					"Invalid input"
//	@Router			/api/v1/secret [post]
func (h *SecretManagerHandler) AddSecret(c echo.Context) error {
	ctx := c.Request().Context()
	var err error

	request := new(requests.CreateSecretRequest)
	if err := c.Bind(request); err != nil {
		return responses.ErrorResponseWithMessage(c, http.StatusBadRequest, "Error parsing data")
	}

	if err := request.Validate(); err != nil {
		return responses.ErrorResponseWithMessage(c, http.StatusBadRequest, "Invalid input")
	}

	var res domain.Secret
	if res, err = h.SecretManager.CreateSecretMessage(ctx, request.ToDomain()); err != nil {
		return responses.ErrorResponseWithMessage(c, http.StatusMethodNotAllowed, "Error creating secret message, Try Again!!!")
	}

	return responses.Response(c, http.StatusOK, response.NewSecretResponse(res))
}

// GetSecretByHash godoc
//	@Summary		Find a secret by hash
//	@Description	Returns a single secret
//	@ID				getSecretByHash
//	@Tags			Secret
//	@Produce		application/json, application/xml
//	@Param			hash	path		string					true	"Unique hash to identify the secret"
//	@Success		200		{object}	response.SecretResponse	"successful operation"
//	@Failure		400		{object}	responses.Error			"Bad request, hash missing"
//	@Failure		404		{object}	responses.Error			"Secret not found"
//	@Router			/api/v1/secret/{hash} [get]
func (h *SecretManagerHandler) GetSecretByHash(c echo.Context) error {
	ctx := c.Request().Context()
	var err error

	hash := c.Param("hash")
	if hash == "" {
		return responses.ErrorResponseWithMessage(c, http.StatusBadRequest, "Hash is required")
	}

	var res domain.Secret
	if res, err = h.SecretManager.GetSecretMessage(ctx, hash); err != nil {
		return responses.ErrorResponseWithMessage(c, http.StatusNotFound, "Error getting secret message")
	}

	return responses.Response(c, http.StatusOK, response.NewSecretResponse(res))
}
