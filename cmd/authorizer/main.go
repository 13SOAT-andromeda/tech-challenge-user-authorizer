package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"tech-challenge-user-authorizer/internal/auth"
	"tech-challenge-user-authorizer/internal/config"
	"tech-challenge-user-authorizer/pkg/utils"
)

var appConfig *config.Config

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if appConfig == nil {
		appConfig = config.LoadConfig()
	}

	utils.InfoLogger.Printf("Processing request: %s", request.Path)

	authHeader, ok := request.Headers["Authorization"]
	if !ok {
		// Fallback to lowercase authorization (API Gateway sometimes lowercases headers)
		authHeader, ok = request.Headers["authorization"]
		if !ok {
			utils.ErrorLogger.Printf("Missing Authorization header")
			return unauthorizedResponse("Missing Authorization header"), nil
		}
	}

	tokenString, err := auth.ExtractBearerToken(authHeader)
	if err != nil {
		utils.ErrorLogger.Printf("Invalid Authorization header format: %v", err)
		return unauthorizedResponse("Invalid Authorization header format"), nil
	}

	_, err = auth.ValidateToken(tokenString, appConfig.JWTSecret, appConfig.JWTIssuer)
	if err != nil {
		utils.ErrorLogger.Printf("Invalid token: %v", err)
		return unauthorizedResponse("Invalid or expired token"), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "{\"message\": \"Authorized\"}",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func unauthorizedResponse(message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 401,
		Body:       "{\"error\": \"" + message + "\"}",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func main() {
	lambda.Start(handler)
}
