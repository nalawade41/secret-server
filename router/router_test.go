package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/golang/mock/gomock"
	"github.com/nalawade41/secret-server/config"
	"github.com/nalawade41/secret-server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	cfg := &config.Config{
		Environment: config.EnvLocal,
	}
	db := &dynamodb.Client{}

	handler := NewHandler(cfg, db)

	assert.NotNil(t, handler)
	assert.Equal(t, cfg, handler.config)
	assert.Equal(t, db, handler.dbConnect)
}

func TestHandler_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoClient := mocks.NewMockDynamoDBAPI(ctrl)

	cfg := &config.Config{
		Environment: config.EnvLocal,
		HTTP:        nil,
		Database: &config.DynamoConfig{
			Host:      "localhost",
			Port:      "8000",
			TableName: "secrets",
		},
		AWS: nil,
	}
	handler := NewHandler(cfg, mockDynamoClient)

	// Initialize the Echo instance
	e := handler.Init()

	// Test the HealthCheck route
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, HealthCheck(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"data":"Server is up and running"}`, rec.Body.String())
	}

	// Test the Swagger route
	req = httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandler_Init_Production(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoClient := mocks.NewMockDynamoDBAPI(ctrl)

	cfg := &config.Config{
		Environment: config.Prod,
		HTTP:        nil,
		Database: &config.DynamoConfig{
			Host:      "localhost",
			Port:      "8000",
			TableName: "secrets",
		},
		AWS: nil,
	}
	handler := NewHandler(cfg, mockDynamoClient)

	// Initialize the Echo instance
	e := handler.Init()

	// Ensure Swagger docs are not available in production
	req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
