package main

import (
	"context"
	"testing"
	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	ctx := context.Background()
	request := events.APIGatewayProxyRequest{
		Path: "/test",
	}

	response, err := handler(ctx, request)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("expected status code 200, got %d", response.StatusCode)
	}

	expectedBody := "{\"message\": \"Hello from User Authorizer!\"}"
	if response.Body != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, response.Body)
	}
}
