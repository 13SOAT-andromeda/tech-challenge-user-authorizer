package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

func TestHandler(t *testing.T) {
	// Setup environment variables for config
	os.Setenv("JWT_SECRET", "test_secret")
	os.Setenv("JWT_ISSUER", "test_issuer")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("JWT_ISSUER")

	ctx := context.Background()

	t.Run("Missing Authorization Header", func(t *testing.T) {
		request := events.APIGatewayProxyRequest{
			Path:    "/test",
			Headers: map[string]string{},
		}
		response, err := handler(ctx, request)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if response.StatusCode != 401 {
			t.Errorf("expected status code 401, got %d", response.StatusCode)
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		request := events.APIGatewayProxyRequest{
			Path: "/test",
			Headers: map[string]string{
				"Authorization": "Bearer invalid_token_string",
			},
		}
		response, err := handler(ctx, request)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if response.StatusCode != 401 {
			t.Errorf("expected status code 401, got %d", response.StatusCode)
		}
	})

	t.Run("Valid Token", func(t *testing.T) {
		// Generate a valid token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "test_issuer",
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte("test_secret"))

		request := events.APIGatewayProxyRequest{
			Path: "/test",
			Headers: map[string]string{
				"Authorization": "Bearer " + tokenString,
			},
		}
		response, err := handler(ctx, request)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if response.StatusCode != 200 {
			t.Errorf("expected status code 200, got %d", response.StatusCode)
		}
	})
}
