package dynamo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang/mock/gomock"
	"github.com/nalawade41/secret-server/internal/common/repository"
	"github.com/nalawade41/secret-server/internal/domain"
	"github.com/nalawade41/secret-server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSecret_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDynamoDBAPI(ctrl)
	repo := SecretManagerRepository{BaseRepository: repository.BaseRepository{DBConnection: mockDB}, TableName: "secrets"}

	hash := "testhash"

	// Set expectations for DeleteItem
	mockDB.EXPECT().DeleteItem(gomock.Any(), &dynamodb.DeleteItemInput{
		TableName: aws.String("secrets"),
		Key: map[string]types.AttributeValue{
			"hash": &types.AttributeValueMemberS{Value: hash},
		},
	}).Return(&dynamodb.DeleteItemOutput{}, nil)

	err := repo.DeleteSecret(context.Background(), hash)

	assert.NoError(t, err)
}

func TestDeleteSecret_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDynamoDBAPI(ctrl)
	repo := SecretManagerRepository{BaseRepository: repository.BaseRepository{DBConnection: mockDB}, TableName: "secrets"}

	hash := "testhash"

	// Set expectations for DeleteItem to return an error
	mockDB.EXPECT().DeleteItem(gomock.Any(), &dynamodb.DeleteItemInput{
		TableName: aws.String("secrets"),
		Key: map[string]types.AttributeValue{
			"hash": &types.AttributeValueMemberS{Value: hash},
		},
	}).Return(nil, errors.New("delete error"))

	err := repo.DeleteSecret(context.Background(), hash)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete secret")
}

func TestUpdateSecretViews_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDynamoDBAPI(ctrl)
	repo := SecretManagerRepository{BaseRepository: repository.BaseRepository{DBConnection: mockDB}, TableName: "secrets"}

	hash := "testhash"
	remainingViews := 4

	// Set expectations for UpdateItem
	mockDB.EXPECT().UpdateItem(gomock.Any(), &dynamodb.UpdateItemInput{
		TableName: aws.String("secrets"),
		Key: map[string]types.AttributeValue{
			"hash": &types.AttributeValueMemberS{Value: hash},
		},
		UpdateExpression: aws.String("SET remainingViews = :remainingViews"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":remainingViews": &types.AttributeValueMemberN{Value: "4"},
			":zero":           &types.AttributeValueMemberN{Value: "0"},
		},
		ConditionExpression: aws.String("remainingViews > :zero"),
	}).Return(&dynamodb.UpdateItemOutput{}, nil)

	err := repo.UpdateSecretViews(context.Background(), hash, remainingViews)

	assert.NoError(t, err)
}

func TestUpdateSecretViews_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDynamoDBAPI(ctrl)
	repo := SecretManagerRepository{BaseRepository: repository.BaseRepository{DBConnection: mockDB}, TableName: "secrets"}

	hash := "testhash"
	remainingViews := 4

	// Set expectations for UpdateItem to return an error
	mockDB.EXPECT().UpdateItem(gomock.Any(), &dynamodb.UpdateItemInput{
		TableName: aws.String("secrets"),
		Key: map[string]types.AttributeValue{
			"hash": &types.AttributeValueMemberS{Value: hash},
		},
		UpdateExpression: aws.String("SET remainingViews = :remainingViews"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":remainingViews": &types.AttributeValueMemberN{Value: "4"},
			":zero":           &types.AttributeValueMemberN{Value: "0"},
		},
		ConditionExpression: aws.String("remainingViews > :zero"),
	}).Return(nil, errors.New("update error"))

	err := repo.UpdateSecretViews(context.Background(), hash, remainingViews)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update remaining views")
}

func TestSave_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDynamoDBAPI(ctrl)
	repo := SecretManagerRepository{BaseRepository: repository.BaseRepository{DBConnection: mockDB}, TableName: "secrets"}

	secret := domain.Secret{
		Hash:           "testhash",
		SecretText:     "This is a test secret",
		ExpiresAt:      time.Now().Add(10 * time.Minute),
		RemainingViews: 5,
		CreatedAt:      time.Now().UTC(),
	}

	item, _ := attributevalue.MarshalMap(secret)

	// Set expectations for PutItem
	mockDB.EXPECT().PutItem(gomock.Any(), &dynamodb.PutItemInput{
		TableName: aws.String("secrets"),
		Item:      item,
	}).Return(&dynamodb.PutItemOutput{}, nil)

	err := repo.Save(context.Background(), secret)

	assert.NoError(t, err)
}

func TestSave_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDynamoDBAPI(ctrl)
	repo := SecretManagerRepository{BaseRepository: repository.BaseRepository{DBConnection: mockDB}, TableName: "secrets"}

	secret := domain.Secret{
		Hash:           "testhash",
		SecretText:     "This is a test secret",
		ExpiresAt:      time.Now().Add(10 * time.Minute),
		RemainingViews: 5,
		CreatedAt:      time.Now().UTC(),
	}

	item, _ := attributevalue.MarshalMap(secret)

	// Set expectations for PutItem to return an error
	mockDB.EXPECT().PutItem(gomock.Any(), &dynamodb.PutItemInput{
		TableName: aws.String("secrets"),
		Item:      item,
	}).Return(nil, errors.New("put item error"))

	err := repo.Save(context.Background(), secret)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to put item")
}

func TestGetByHash_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDynamoDBAPI(ctrl)
	repo := SecretManagerRepository{BaseRepository: repository.BaseRepository{DBConnection: mockDB}, TableName: "secrets"}

	hash := "testhash"
	secret := domain.Secret{
		Hash:           hash,
		SecretText:     "This is a test secret",
		ExpiresAt:      time.Now().UTC().Add(10 * time.Minute),
		RemainingViews: 5,
		CreatedAt:      time.Now().UTC(),
	}

	item, _ := attributevalue.MarshalMap(secret)

	// Set expectations for GetItem
	mockDB.EXPECT().GetItem(gomock.Any(), &dynamodb.GetItemInput{
		TableName: aws.String("secrets"),
		Key: map[string]types.AttributeValue{
			"hash": &types.AttributeValueMemberS{Value: hash},
		},
	}).Return(&dynamodb.GetItemOutput{Item: item}, nil)

	result, err := repo.GetByHash(context.Background(), hash)

	assert.NoError(t, err)
	assert.Equal(t, secret, result)
}

func TestGetByHash_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDynamoDBAPI(ctrl)
	repo := SecretManagerRepository{BaseRepository: repository.BaseRepository{DBConnection: mockDB}, TableName: "secrets"}

	hash := "nonexistenthash"

	// Set expectations for GetItem to return no item
	mockDB.EXPECT().GetItem(gomock.Any(), &dynamodb.GetItemInput{
		TableName: aws.String("secrets"),
		Key: map[string]types.AttributeValue{
			"hash": &types.AttributeValueMemberS{Value: hash},
		},
	}).Return(&dynamodb.GetItemOutput{Item: nil}, nil)

	_, err := repo.GetByHash(context.Background(), hash)

	assert.Error(t, err)
	assert.Equal(t, "secret not found", err.Error())
}

func TestGetByHash_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDynamoDBAPI(ctrl)
	repo := SecretManagerRepository{BaseRepository: repository.BaseRepository{DBConnection: mockDB}, TableName: "secrets"}

	hash := "testhash"

	// Set expectations for GetItem to return an error
	mockDB.EXPECT().GetItem(gomock.Any(), &dynamodb.GetItemInput{
		TableName: aws.String("secrets"),
		Key: map[string]types.AttributeValue{
			"hash": &types.AttributeValueMemberS{Value: hash},
		},
	}).Return(nil, errors.New("get item error"))

	_, err := repo.GetByHash(context.Background(), hash)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to retrieve item")
}
