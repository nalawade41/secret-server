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
	@make run APP_ENV=production

# Target to run the application in development mode
.PHONY: run-development
run-development:
	@make run APP_ENV=development


.PHONY: run-wire
run-wire:
	cd internal/wire && go run -mod=mod github.com/google/wire/cmd/wire