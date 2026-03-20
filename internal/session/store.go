// Package session provides DynamoDB-backed session lookup for the authorizer.
//
// Environment (see README.md and .env.example):
//   - SESSION_TABLE_NAME — table name (default in config: user-sessions)
//   - AWS_REGION — AWS SDK region (default us-east-1)
//   - DYNAMODB_ENDPOINT — optional override for LocalStack (e.g. http://host.docker.internal:4566)
package session

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var ErrSessionNotFound = errors.New("session not found")

type Session struct {
	JTI    string
	UserID string
}

type Store interface {
	GetSessionByUserID(ctx context.Context, userID string) (*Session, error)
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

func (s *DynamoStore) GetSessionByUserID(ctx context.Context, userID string) (*Session, error) {
	item, err := s.getItem(ctx, userID, "N")
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "ValidationException" {
			item, err = s.getItem(ctx, userID, "S")
		}
	}
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrSessionNotFound
	}

	jti, err := firstStringAttr(item, []string{"jti", "Jti", "JTI"})
	if err != nil || jti == "" {
		return nil, fmt.Errorf("invalid jti in session record")
	}
	sessionUserID, err := firstStringAttr(item, []string{"userId", "UserId", "user_id"})
	if err != nil || sessionUserID == "" {
		return nil, fmt.Errorf("invalid userId in session record")
	}

	return &Session{
		JTI:    jti,
		UserID: sessionUserID,
	}, nil
}

func (s *DynamoStore) getItem(ctx context.Context, userID string, attrType string) (map[string]*dynamodb.AttributeValue, error) {
	key := map[string]*dynamodb.AttributeValue{}
	if attrType == "N" {
		key["userId"] = &dynamodb.AttributeValue{N: aws.String(userID)}
	} else {
		key["userId"] = &dynamodb.AttributeValue{S: aws.String(userID)}
	}

	out, err := s.client.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	if len(out.Item) == 0 {
		return nil, ErrSessionNotFound
	}

	return out.Item, nil
}

func firstStringAttr(item map[string]*dynamodb.AttributeValue, keys []string) (string, error) {
	for _, k := range keys {
		if v, ok := item[k]; ok {
			s, err := attributeValueToString(v)
			if err == nil && s != "" {
				return s, nil
			}
		}
	}
	return "", errors.New("missing attribute")
}

func attributeValueToString(attr *dynamodb.AttributeValue) (string, error) {
	if attr == nil {
		return "", errors.New("missing attribute")
	}
	if attr.S != nil {
		return *attr.S, nil
	}
	if attr.N != nil {
		return *attr.N, nil
	}

	return "", errors.New("unsupported attribute type")
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
