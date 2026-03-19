package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"tech-challenge-user-authorizer/pkg/utils"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	utils.InfoLogger.Printf("Processing request: %s", request.Path)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "{\"message\": \"Hello from User Authorizer!\"}",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
