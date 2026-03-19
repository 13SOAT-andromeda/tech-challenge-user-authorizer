package auth

import (
	"testing"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

func TestValidateToken(t *testing.T) {
	secret := "test_secret"
	issuer := "test_issuer"
	
	t.Run("Valid Token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": issuer,
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(secret))

		var claims jwt.MapClaims
		var err error
		claims, err = ValidateToken(tokenString, secret, issuer)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		iss, _ := claims.GetIssuer()
		if iss != issuer {
			t.Errorf("expected issuer %s, got %v", issuer, iss)
		}
	})

	t.Run("Invalid Secret", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": issuer,
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte("wrong_secret"))

		_, err := ValidateToken(tokenString, secret, issuer)
		if err == nil {
			t.Error("expected error for invalid secret, got nil")
		}
	})

	t.Run("Expired Token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": issuer,
			"exp": time.Now().Add(-time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(secret))

		_, err := ValidateToken(tokenString, secret, issuer)
		if err == nil {
			t.Error("expected error for expired token, got nil")
		}
	})

	t.Run("Invalid Issuer", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "wrong_issuer",
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(secret))

		_, err := ValidateToken(tokenString, secret, issuer)
		if err == nil {
			t.Error("expected error for invalid issuer, got nil")
		}
	})

	t.Run("Invalid Signing Method", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
			"iss": issuer,
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString := token.Raw

		_, err := ValidateToken(tokenString, secret, issuer)
		if err == nil {
			t.Error("expected error for invalid signing method, got nil")
		}
	})

	t.Run("Missing Issuer", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(secret))

		_, err := ValidateToken(tokenString, secret, issuer)
		if err == nil {
			t.Error("expected error for missing issuer, got nil")
		}
	})
}
