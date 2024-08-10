# Makefile for running the application

# Set default environment
APP_ENV ?= local

# Path to the Go source file
SRC = main.go

# Default target
.PHONY: run
run:
	@echo "Running application in $(APP_ENV) environment..."
	APP_ENV=$(APP_ENV) go run $(SRC)

# Target to run the application in production mode
.PHONY: run-production
run-production:
	@make run APP_ENV=prod

# Target to run the application in local mode
.PHONY: run-local
run-development:
	@make run APP_ENV=local


.PHONY: run-wire
run-wire:
	cd internal/wire && go run -mod=mod github.com/google/wire/cmd/wire

.PHONY: run-swagger
run-swagger:
	swag init -g cmd/local/main.go


# Define a target to generate mocks
.PHONY: mockgen
mockgen:
	@echo "Generating mocks..."
	@mockgen -source=internal/common/repository/base_repository.go -destination=mocks/base_repository_mock.go -package=mocks
	@mockgen -source=db/db_api.go -destination=mocks/dynamodbapi_mock.go -package=mocks
	@mockgen -destination=mocks/secret_repository_mock.go -package=mocks github.com/nalawade41/secret-server/internal/domain SecretRepository
	@mockgen -destination=mocks/encryptor_mock.go -package=mocks github.com/nalawade41/secret-server/internal/domain Encryptor
	@mockgen -destination=mocks/mock_secret_usecase.go -package=mocks github.com/nalawade41/secret-server/internal/domain SecretUseCase
