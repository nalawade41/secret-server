
# Secret Server

This project is a Secret Server API that allows storing and sharing secrets using a randomly generated URL. Each secret can only be read a limited number of times and may have a TTL (time-to-live) after which it expires.

## Features

- **Store and Share Secrets**: Save secrets with a unique URL for sharing.
- **Access Control**: Limit the number of views for each secret.
- **Expiration**: Set a TTL for secrets after which they are no longer accessible.
- **JSON/XML Response**: Supports JSON and XML responses based on the `Accept` header.
- **Swagger Documentation**: Provides API documentation and testing via Swagger UI.

## Prerequisites

- Go 1.18 or later
- AWS SDK for Go v2
- Docker (for running DynamoDB locally)
- AWS CLI (for deploying with CDK)

## Running the Server Locally

### Step 1: Install Dependencies

Ensure you have Go installed and set up your Go workspace.

```bash
# Install dependencies
go mod tidy
```

### Step 2: Run DynamoDB Locally

You can run DynamoDB locally using Docker Compose. A `docker-compose.yml` file is provided in the repository.

```bash
# Start DynamoDB using Docker Compose
docker-compose up -d
```

This will start DynamoDB on `http://localhost:8000`.

### Step 3: Configure Environment Variables

Create a `.env.local` file in the root of your project and set the following environment variables:

```plaintext
ENVIRONMENT=local
HTTP_HOST=localhost
HTTP_PORT=3000
HTTP_PROTOCOL=http
READ_TIMEOUT=5
WRITE_TIMEOUT=5
MAX_HEADER_BYTES=1048576
DB_HOST=localhost
DB_PORT=8000
DB_TABLE_NAME=secrets
```
Please check `.env_example` for more details and other environment variables.

### Step 4: Run the Server

Navigate to the `cmd/local` directory and run the `main.go` file:

```bash
cd cmd/local
go run main.go
```

The server will start on `http://localhost:3000`.

### Step 5: Access Swagger Documentation

Once the server is running, you can access the Swagger UI for API documentation and testing at:

```
http://localhost:8080/swagger/index.html
```

## Project Structure

```
.
├── cmd
│   ├── local
│   │   └── main.go      # Entry point for running the server locally
│   └── app             
│       └── main.go      # Entry point for running the server in production (lambda)
├── config               # Configuration management
├── db                   # Database client setup
├── docs                 # Swagger documentation files
│   └── swagger.yaml
├── infra               # Infrastructure setup
│   └── app
│       └── deploy      # CDK deployment scripts
├── server               # HTTP server setup
│   └── server.go
├── internal
│   ├── common           # Common utilities and helpers
│   │   ├── security     # Encryption and hashing utilities
│   │   │   ├── encrypt.go
│   │   │   └── hash.go
│   │   ├── logger       # Logging utilities
│   │   │   └── logger.go
│   │   ├── repository   # Base Repository
│   │   │   └── repository.go
│   │   ├── response     # Response utilities
│   │   │   └── response.go
│   │   └── constants    # Constants
│   ├── domain           # Internal domain models and interfaces
│   └── secret           # Contains the business logic
│       ├── handler
│           └── handler.go
│       ├── repository
│           └── repository.go
│       ├── requests
│           └── request.go
│       ├── response
│           └── response.go
│       ├── usecase
│           └── usecase.go
│       └── provider.go  # Dependency injection setup
├── wire                 # Dependency injection setup
├── router               # HTTP routing setup
│   └── router.go
├── trace                # Request tracing setup
│   └── trace.go
├── .env_example         # Example environment variables
├── .gitignore           # Git ignore file
├── .swagignore          # Swagger ignore file
├── go.mod               # Go module file
├── Makefile             # Makefile for common tasks
├── docker-compose.yml   # Docker Compose configuration for local DynamoDB
└── README.md            # Project README

```

## Testing

Run unit tests with:

```bash
go test ./...
```

## API Endpoints

### Create a Secret

- **Endpoint**: `/api/v1/secrets`
- **Method**: `POST`
- **Description**: Create a new secret with specified TTL and view count.
- **Request Body**:
  ```json
  {
    "secretText": "your-secret",
    "ttl": 60,  // TTL in minutes
    "remainingViews": 5
  }
  ```
- **Response**: Returns the created secret's hash and other details.

### Get a Secret

- **Endpoint**: `/api/v1/secrets/{hash}`
- **Method**: `GET`
- **Description**: Retrieve a secret by its hash. The response format is based on the `Accept` header (JSON/XML).
- **Response**: Returns the secret text if it is not expired or exceeded its view count.

## Configuration

The server can be configured using environment variables defined in the `.env` file or directly set in the environment:

- `AWS_REGION`: The AWS region where DynamoDB is hosted.
- `DYNAMODB_ENDPOINT`: The endpoint for connecting to DynamoDB (useful for local testing).
- `TABLE_NAME`: The name of the DynamoDB table.
- `ENVIRONMENT`: The environment mode (e.g., `local`, `dev`, `prod`).

## CDK Deployment

The project is configured to deploy to AWS using AWS Cloud Development Kit (CDK). The deployment is triggered automatically on push to the `dev` branch.

### Deployment Process

1. **GitHub Actions**: The deployment workflow is configured in `.github/workflows/deploy.yml`.
2. **Automatic Deployment**: On push to the `dev` branch, the workflow runs the CDK deployment script.
3. **CDK Stack**: The CDK stack defines the necessary AWS resources, such as Lambda functions, API Gateway, and DynamoDB tables.

### Setting Up CDK

Ensure you have the AWS CLI and CDK installed:

```bash
npm install -g aws-cdk
```

Configure your AWS credentials:

```bash
aws configure
```

Deploy the stack manually (if needed):

```bash
cd cdk
cdk deploy
```

## License

This project is licensed under the MIT License.
