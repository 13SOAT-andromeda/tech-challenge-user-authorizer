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
	"os"

	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var ErrSessionNotFound = errors.New("session not found")

type Session struct {
	UserID string
}

type Store interface {
	GetSessionByJTI(ctx context.Context, jti string) (*Session, error)
}

type DynamoStore struct {
	client    *dynamodb.DynamoDB
	tableName string
}

func NewDynamoStore(tableName string) (*DynamoStore, error) {
	cfg := aws.NewConfig().WithRegion(getEnv("AWS_REGION", "us-east-1"))
	if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
		cfg = cfg.WithEndpoint(endpoint)
	}

	sess, err := awsSession.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create aws session: %w", err)
	}

	return &DynamoStore{
		client:    dynamodb.New(sess),
		tableName: tableName,
	}, nil
}

func (s *DynamoStore) GetSessionByJTI(ctx context.Context, jti string) (*Session, error) {
	out, err := s.client.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"token_id": {S: aws.String(jti)},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(out.Item) == 0 {
		return nil, ErrSessionNotFound
	}

	userIDAttr, ok := out.Item["user_id"]
	if !ok || userIDAttr.S == nil || *userIDAttr.S == "" {
		return nil, fmt.Errorf("invalid user_id in session record")
	}

	return &Session{
		UserID: *userIDAttr.S,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
