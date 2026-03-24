package main

import (
	"context"
	"os"
	"tech-challenge-user-authorizer/internal/session"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

type mockSessionStore struct {
	sessionByJTI map[string]*session.Session
	err          error
}

func (m *mockSessionStore) GetSessionByJTI(_ context.Context, jti string) (*session.Session, error) {
	if m.err != nil {
		return nil, m.err
	}
	sessionData, ok := m.sessionByJTI[jti]
	if !ok {
		return nil, session.ErrSessionNotFound
	}
	return sessionData, nil
}

func TestHandler(t *testing.T) {
	appConfig = nil
	sessionStore = nil

	// Setup environment variables for config
	os.Setenv("JWT_SECRET", "test_secret")
	os.Setenv("JWT_ISSUER", "test_issuer")
	os.Setenv("SESSION_TABLE_NAME", "sessions-test")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("JWT_ISSUER")
	defer os.Unsetenv("SESSION_TABLE_NAME")

	ctx := context.Background()

	t.Run("Missing Authorization Header", func(t *testing.T) {
		request := events.APIGatewayV2HTTPRequest{
			RawPath: "/authorize",
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

	t.Run("Invalid Token Format", func(t *testing.T) {
		request := events.APIGatewayV2HTTPRequest{
			RawPath: "/authorize",
			Headers: map[string]string{
				"authorization": "Bearer invalid_token_string",
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

	t.Run("Valid Token and Session", func(t *testing.T) {
		tokenJTI := "test-jti"
		userID := "1"
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "test_issuer",
			"sub": userID,
			"jti": tokenJTI,
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte("test_secret"))

		sessionStore = &mockSessionStore{
			sessionByJTI: map[string]*session.Session{
				tokenJTI: {
					UserID: userID,
				},
			},
		}

		request := events.APIGatewayV2HTTPRequest{
			RawPath: "/authorize",
			Headers: map[string]string{
				"authorization": "Bearer " + tokenString,
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

	t.Run("Token Missing JTI", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "test_issuer",
			"sub": "1",
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte("test_secret"))

		request := events.APIGatewayV2HTTPRequest{
			RawPath: "/authorize",
			Headers: map[string]string{
				"authorization": "Bearer " + tokenString,
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

	t.Run("Session Not Found", func(t *testing.T) {
		userID := "1"
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "test_issuer",
			"sub": userID,
			"jti": "some-jti",
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte("test_secret"))

		sessionStore = &mockSessionStore{
			sessionByJTI: map[string]*session.Session{},
		}

		request := events.APIGatewayV2HTTPRequest{
			RawPath: "/authorize",
			Headers: map[string]string{
				"authorization": "Bearer " + tokenString,
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

	t.Run("Session UserID Mismatch", func(t *testing.T) {
		tokenJTI := "token-jti"
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "test_issuer",
			"sub": "1",
			"jti": tokenJTI,
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte("test_secret"))

		sessionStore = &mockSessionStore{
			sessionByJTI: map[string]*session.Session{
				tokenJTI: {
					UserID: "999",
				},
			},
		}

		request := events.APIGatewayV2HTTPRequest{
			RawPath: "/authorize",
			Headers: map[string]string{
				"authorization": "Bearer " + tokenString,
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
}
