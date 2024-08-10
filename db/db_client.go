package db

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	lConfig "github.com/nalawade41/secret-server/config"
	"github.com/nalawade41/secret-server/internal/common/logger"
)

var (
	dynamoDBClient DynamoDBAPI
	dynamoDBOnce   = new(sync.Once)
)

// InitDynamoDB initializes the DynamoDB connection
func InitDynamoDB(cfg *lConfig.Config) (DynamoDBAPI, error) {
	var initErr error

	dynamoDBOnce.Do(func() {
		awsConfig, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(cfg.AWS.Region),
		)
		if err != nil {
			logger.Errorf("failed to load AWS config: %v", err)
			initErr = err
			return
		}

		if cfg.Environment == lConfig.EnvLocal {
			awsConfig, err = config.LoadDefaultConfig(context.TODO(),
				config.WithRegion(cfg.AWS.Region),
				config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
					func(service, region string, options ...interface{}) (aws.Endpoint, error) {
						if service == dynamodb.ServiceID && region == "us-west-2" {
							return aws.Endpoint{URL: fmt.Sprintf("http://%s:%s", cfg.Database.Host, cfg.Database.Port)}, nil
						}
						return aws.Endpoint{}, &aws.EndpointNotFoundError{}
					})),
				config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
					Value: aws.Credentials{
						AccessKeyID: "a1b2c3", SecretAccessKey: "a1b2c3", SessionToken: "a1b2c3",
						Source: "Mock credentials used above for local instance",
					},
				}),
			)
			if err != nil {
				logger.Errorf("failed to load AWS config: %v", err)
				initErr = err
				return
			}
		}

		// Create DynamoDB client
		dynamoDBClient = dynamodb.NewFromConfig(awsConfig)

		// Check if the table exists
		exists, err := doesTableExist(context.TODO(), dynamoDBClient, cfg.Database.TableName)
		if err != nil {
			logger.Errorf("failed to check table existence: %v", err)
			initErr = err
			return
		}

		// Create the table if it doesn't exist
		if !exists {
			err = createTable(context.TODO(), dynamoDBClient, cfg.Database.TableName)
			if err != nil {
				logger.Errorf("failed to create table: %v", err)
				initErr = err
				return
			} else {
				fmt.Printf("Table %s created successfully\n", cfg.Database.TableName)
			}
		} else {
			logger.Infof("Table %s already exists", cfg.Database.TableName)
		}
	})

	return dynamoDBClient, initErr
}

// doesTableExist checks if a DynamoDB table exists
func doesTableExist(ctx context.Context, svc DynamoDBAPI, tableName string) (bool, error) {
	// Use DescribeTable to check if the table exists
	_, err := svc.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	// Check for errors
	if err != nil {
		var notFoundErr *types.ResourceNotFoundException
		if ok := errors.As(err, &notFoundErr); ok {
			return false, nil // Table does not exist
		}
		return false, fmt.Errorf("failed to describe table: %w", err)
	}

	return true, nil // Table exists
}

// createTable creates a new DynamoDB table
func createTable(ctx context.Context, svc DynamoDBAPI, tableName string) error {
	_, err := svc.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("hash"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("hash"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Wait for the table to become active
	return waitForTableToBeActive(ctx, svc, tableName)
}

// waitForTableToBeActive waits for a DynamoDB table to become active
func waitForTableToBeActive(ctx context.Context, svc DynamoDBAPI, tableName string) error {
	waitTime := 5 * time.Second
	maxRetries := 12 // Retry for up to 1 minute

	for i := 0; i < maxRetries; i++ {
		desc, err := svc.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return fmt.Errorf("failed to describe table: %w", err)
		}

		if desc.Table.TableStatus == types.TableStatusActive {
			return nil
		}

		fmt.Printf("Waiting for table %s to become active...\n", tableName)
		time.Sleep(waitTime)
	}

	return fmt.Errorf("table %s did not become active in time", tableName)
}
