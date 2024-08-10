package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/nalawade41/secret-server/internal/common/repository"
	"github.com/nalawade41/secret-server/internal/domain"
	"github.com/pkg/errors"
)

type SecretManagerRepository struct {
	repository.BaseRepository
	TableName string
}

func (s SecretManagerRepository) DeleteSecret(ctx context.Context, hash string) error {
	_, err := s.DBConnection.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(s.TableName),
		Key: map[string]types.AttributeValue{
			"hash": &types.AttributeValueMemberS{Value: hash},
		},
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to delete secret for hash: %s", hash))
	}

	return nil
}

func (s SecretManagerRepository) UpdateSecretViews(ctx context.Context, hash string, remainingViews int) error {
	_, err := s.DBConnection.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(s.TableName),
		Key: map[string]types.AttributeValue{
			"hash": &types.AttributeValueMemberS{Value: hash},
		},
		UpdateExpression: aws.String("SET remainingViews = :remainingViews"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":remainingViews": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", remainingViews)},
			":zero":           &types.AttributeValueMemberN{Value: "0"},
		},
		ConditionExpression: aws.String("remainingViews > :zero"),
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to update remaining views for hash: %s", hash))
	}

	return nil
}

func (s SecretManagerRepository) Save(ctx context.Context, secret domain.Secret) error {
	// Marshal the secret into a map of DynamoDB attribute values
	item, err := attributevalue.MarshalMap(secret)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to marshal secret: %w", err))
	}

	// Put the item into the DynamoDB table
	_, err = s.DBConnection.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.TableName),
		Item:      item,
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to put item: %w", err))
	}

	return nil
}

func (s SecretManagerRepository) GetByHash(ctx context.Context, hash string) (domain.Secret, error) {
	result, err := s.DBConnection.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.TableName),
		Key: map[string]types.AttributeValue{
			"hash": &types.AttributeValueMemberS{Value: hash},
		},
	})
	if err != nil {
		return domain.Secret{}, errors.Wrap(err, fmt.Sprintf("failed to get item for hash: %s", hash))
	}

	if result.Item == nil {
		return domain.Secret{}, errors.New("secret not found")
	}

	// Unmarshal the result into a domain.Secret struct
	var secret domain.Secret
	if err := attributevalue.UnmarshalMap(result.Item, &secret); err != nil {
		return domain.Secret{}, errors.Wrap(err, fmt.Sprintf("failed to unmarshal item: %w", err))
	}

	return secret, nil
}

var _ domain.SecretRepository = (*SecretManagerRepository)(nil)
