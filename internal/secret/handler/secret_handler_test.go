package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/nalawade41/secret-server/internal/domain"
	"github.com/nalawade41/secret-server/internal/secret/response"
	"github.com/nalawade41/secret-server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAddSecret_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockSecretUseCase(ctrl)

	handler := SecretManagerHandler{SecretManager: mockUseCase}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/secret", bytes.NewBufferString(`{"secret":"This is a test secret","expireAfter":10,"expireAfterViews":5}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON) // Set to application/json
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedSecret := domain.Secret{
		Hash:           "testhash",
		SecretText:     "This is a test secret",
		ExpiresAt:      time.Now().Add(10 * time.Minute),
		RemainingViews: 5,
		CreatedAt:      time.Now().UTC(),
	}

	// Set up the expectation for CreateSecretMessage call
	mockUseCase.EXPECT().CreateSecretMessage(gomock.Any(), gomock.Any()).Return(expectedSecret, nil)

	if assert.NoError(t, handler.AddSecret(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var secretResponse response.SecretResponse
		err := json.Unmarshal(rec.Body.Bytes(), &secretResponse)
		assert.NoError(t, err)
		assert.Equal(t, expectedSecret.Hash, secretResponse.Hash)
		assert.Equal(t, expectedSecret.SecretText, secretResponse.SecretText)
	}
}

func TestAddSecret_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockSecretUseCase(ctrl)

	handler := SecretManagerHandler{SecretManager: mockUseCase}

	e := echo.New()
	// Set content type to application/json to trigger a bind error with malformed JSON
	req := httptest.NewRequest(http.MethodPost, "/api/v1/secret", bytes.NewBufferString(`invalid json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.AddSecret(c)

	// Check the response code and error message
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Error parsing data")
	}
}

func TestAddSecret_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockSecretUseCase(ctrl)

	handler := SecretManagerHandler{SecretManager: mockUseCase}

	e := echo.New()
	// Use application/json for JSON payload
	req := httptest.NewRequest(http.MethodPost, "/api/v1/secret", bytes.NewBufferString(`{"secret":"","expireAfter":-1,"expireAfterViews":-1}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.AddSecret(c)

	// Echo does not return errors from the handler directly, it writes them to the response
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid input")
	}
}

func TestAddSecret_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockSecretUseCase(ctrl)

	handler := SecretManagerHandler{SecretManager: mockUseCase}

	e := echo.New()
	// Use application/json for JSON payload
	req := httptest.NewRequest(http.MethodPost, "/api/v1/secret", bytes.NewBufferString(`{"secret":"This is a test secret","expireAfter":10,"expireAfterViews":5}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up the mock expectation for a repository error
	mockUseCase.EXPECT().CreateSecretMessage(gomock.Any(), gomock.Any()).Return(domain.Secret{}, errors.New("repository error"))

	err := handler.AddSecret(c)

	// Echo does not return errors from the handler directly, it writes them to the response
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		assert.Contains(t, rec.Body.String(), "Error creating secret message")
	}
}

func TestGetSecretByHash_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockSecretUseCase(ctrl)

	handler := SecretManagerHandler{SecretManager: mockUseCase}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/secret/testhash", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("hash")
	c.SetParamValues("testhash")

	expectedSecret := domain.Secret{
		Hash:           "testhash",
		SecretText:     "This is a test secret",
		ExpiresAt:      time.Now().Add(10 * time.Minute),
		RemainingViews: 5,
		CreatedAt:      time.Now().UTC(),
	}

	mockUseCase.EXPECT().GetSecretMessage(gomock.Any(), "testhash").Return(expectedSecret, nil)

	if assert.NoError(t, handler.GetSecretByHash(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var secretResponse response.SecretResponse
		err := json.Unmarshal(rec.Body.Bytes(), &secretResponse)
		assert.NoError(t, err)
		assert.Equal(t, expectedSecret.Hash, secretResponse.Hash)
		assert.Equal(t, expectedSecret.SecretText, secretResponse.SecretText)
	}
}

func TestGetSecretByHash_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockSecretUseCase(ctrl)

	handler := SecretManagerHandler{SecretManager: mockUseCase}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/secret/nonexistenthash", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("hash")
	c.SetParamValues("nonexistenthash")

	// Set up the expectation for GetSecretMessage to return an error indicating the secret was not found
	mockUseCase.EXPECT().GetSecretMessage(gomock.Any(), "nonexistenthash").Return(domain.Secret{}, errors.New("secret not found"))

	// Call the handler
	if assert.NoError(t, handler.GetSecretByHash(c)) {
		// Check that the status code and response body contain the expected error message
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "Error getting secret message")
	}
}

func TestGetSecretByHash_MissingHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockSecretUseCase(ctrl)

	handler := SecretManagerHandler{SecretManager: mockUseCase}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/secret/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// We need to explicitly unset the parameter to simulate a missing hash
	c.SetParamNames("hash")
	c.SetParamValues("") // Simulating missing hash by setting an empty string

	// Call the handler
	if assert.NoError(t, handler.GetSecretByHash(c)) {
		// Verify that the status code is 400 and the response contains the expected error message
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Hash is required")
	}
}
