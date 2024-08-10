package db

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang/mock/gomock"
	"github.com/nalawade41/secret-server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDoesTableExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoClient := mocks.NewMockDynamoDBAPI(ctrl)

	tableName := "secrets"

	// Test case: Table exists
	mockDynamoClient.EXPECT().
		DescribeTable(gomock.Any(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		}).
		Return(&dynamodb.DescribeTableOutput{}, nil)

	exists, err := doesTableExist(context.TODO(), mockDynamoClient, tableName)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test case: Table does not exist
	mockDynamoClient.EXPECT().
		DescribeTable(gomock.Any(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		}).
		Return(nil, &types.ResourceNotFoundException{})

	exists, err = doesTableExist(context.TODO(), mockDynamoClient, tableName)
	assert.NoError(t, err)
	assert.False(t, exists)

	// Test case: Other error
	mockDynamoClient.EXPECT().
		DescribeTable(gomock.Any(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		}).
		Return(nil, errors.New("unexpected error"))

	exists, err = doesTableExist(context.TODO(), mockDynamoClient, tableName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to describe table")
	assert.False(t, exists)
}

func TestCreateTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoClient := mocks.NewMockDynamoDBAPI(ctrl)

	tableName := "secrets"

	// Mock CreateTable and wait for table to be active
	mockDynamoClient.EXPECT().
		CreateTable(gomock.Any(), gomock.Any()).
		Return(&dynamodb.CreateTableOutput{}, nil)

	mockDynamoClient.EXPECT().
		DescribeTable(gomock.Any(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		}).
		Return(&dynamodb.DescribeTableOutput{
			Table: &types.TableDescription{
				TableStatus: types.TableStatusActive,
			},
		}, nil).AnyTimes()

	err := createTable(context.TODO(), mockDynamoClient, tableName)
	assert.NoError(t, err)

	// Test case: CreateTable error
	mockDynamoClient.EXPECT().
		CreateTable(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("create table error"))

	err = createTable(context.TODO(), mockDynamoClient, tableName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create table")
}
