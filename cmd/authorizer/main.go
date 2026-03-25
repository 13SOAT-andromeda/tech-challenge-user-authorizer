package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"tech-challenge-user-authorizer/internal/auth"
	"tech-challenge-user-authorizer/internal/config"
	"tech-challenge-user-authorizer/internal/session"
	"tech-challenge-user-authorizer/pkg/utils"

	ddlambda "github.com/DataDog/dd-trace-go/contrib/aws/datadog-lambda-go/v2"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var appConfig *config.Config

var sessionStore session.Store
var newSessionStore = func(tableName, region, endpoint string) (session.Store, error) {
	return session.NewDynamoStore(tableName, region, endpoint)
}

func handler(ctx context.Context, request events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {

	if appConfig == nil {
		appConfig = config.LoadConfig()
	}

	utils.InfoLogger.Printf("Configuring session Store: %s", appConfig.DynamoDBTableName)

	if sessionStore == nil {
		store, err := newSessionStore(appConfig.DynamoDBTableName, appConfig.AWSRegion, appConfig.DynamoDBEndpoint)
		if err != nil {
			utils.ErrorLogger.Printf("Failed to initialize session store: %v", err)
			return denyResponse(), nil
		}
		sessionStore = store
	}

	utils.InfoLogger.Printf("Processing request: %s", request.RawPath)

	// API Gateway v2 always lowercases headers
	authHeader, ok := request.Headers["authorization"]
	if !ok {
		utils.ErrorLogger.Printf("Missing Authorization header")
		return denyResponse(), nil
	}

	tokenString, err := auth.ExtractBearerToken(authHeader)
	if err != nil {
		utils.ErrorLogger.Printf("Invalid Authorization header format: %v", err)
		return denyResponse(), nil
	}

	claims, err := auth.ValidateToken(tokenString, appConfig.JWTSecret, appConfig.JWTIssuer)
	if err != nil {
		utils.ErrorLogger.Printf("Invalid token: %v", err)
		return denyResponse(), nil
	}

	tokenJTI, err := getRequiredStringClaim(claims, "jti")
	if err != nil {
		utils.ErrorLogger.Printf("Invalid token claims: %v", err)
		return denyResponse(), nil
	}
	userID, err := getUserIDFromClaims(claims)
	if err != nil {
		utils.ErrorLogger.Printf("Invalid user identifier in token: %v", err)
		return denyResponse(), nil
	}

	activeSession, err := sessionStore.GetSessionByJTI(ctx, tokenJTI)
	if err != nil {
		if errors.Is(err, session.ErrSessionNotFound) {
			utils.ErrorLogger.Printf("No active session found for jti %s", tokenJTI)
			return denyResponse(), nil
		}
		utils.ErrorLogger.Printf("Failed to validate active session: %v", err)
		return denyResponse(), nil
	}
	if strings.TrimSpace(activeSession.UserID) != strings.TrimSpace(userID) {
		utils.ErrorLogger.Printf("Session userId mismatch for token user %s (stored=%s)", userID, activeSession.UserID)
		return denyResponse(), nil
	}

	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
		Context: map[string]interface{}{
			"userId": userID,
		},
	}, nil
}

func denyResponse() events.APIGatewayV2CustomAuthorizerSimpleResponse {
	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: false,
	}
}

func getRequiredStringClaim(claims map[string]interface{}, key string) (string, error) {
	value, exists := claims[key]
	if !exists {
		return "", fmt.Errorf("missing %s claim", key)
	}

	strValue, ok := value.(string)
	if !ok || strValue == "" {
		return "", fmt.Errorf("invalid %s claim", key)
	}

	return strValue, nil
}

func getUserIDFromClaims(claims map[string]interface{}) (string, error) {
	if value, exists := claims["user_id"]; exists {
		switch typed := value.(type) {
		case string:
			if typed != "" {
				return typed, nil
			}
		case float64:
			return strconv.FormatInt(int64(typed), 10), nil
		}
	}

	subject, err := getRequiredStringClaim(claims, "sub")
	if err != nil {
		return "", err
	}
	return subject, nil
}

func main() {
	lambda.Start(ddlambda.WrapFunction(handler, nil))
}
