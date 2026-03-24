// Package session provides DynamoDB-backed session lookup for the authorizer.
//
// Environment (see README.md and .env.example):
//   - SESSION_TABLE_NAME — table name (default in config: user-authentication-token)
//   - AWS_REGION — AWS SDK region (default us-east-1)
//   - DYNAMODB_ENDPOINT — optional override for LocalStack (e.g. http://host.docker.internal:4566)
package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var ErrSessionNotFound = errors.New("session not found")

type Session struct {
	UserID string
}

type Store interface {
	GetSessionByJTI(ctx context.Context, jti string) (*Session, error)
}

type DynamoStore struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoStore(tableName, region, endpoint string) (*DynamoStore, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}
	if endpoint != "" {
		opts = append(opts, config.WithBaseEndpoint(endpoint))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	return &DynamoStore{
		client:    dynamodb.NewFromConfig(cfg),
		tableName: tableName,
	}, nil
}

func (s *DynamoStore) GetSessionByJTI(ctx context.Context, jti string) (*Session, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]types.AttributeValue{
			"token_id": &types.AttributeValueMemberS{Value: jti},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(out.Item) == 0 {
		return nil, ErrSessionNotFound
	}

	userIDAttr, ok := out.Item["user_id"]
	if !ok {
		return nil, fmt.Errorf("invalid user_id in session record")
	}
	userIDVal, ok := userIDAttr.(*types.AttributeValueMemberS)
	if !ok || userIDVal.Value == "" {
		return nil, fmt.Errorf("invalid user_id in session record")
	}

	return &Session{
		UserID: userIDVal.Value,
	}, nil
}
