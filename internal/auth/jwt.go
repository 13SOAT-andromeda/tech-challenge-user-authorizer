package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// ValidateToken parses a JWT token string, validates its signature using a symmetric key,
// checks its expiration, and verifies the issuer claim.
func ValidateToken(tokenString, secret, expectedIssuer string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Verify issuer
		iss, err := claims.GetIssuer()
		if err != nil {
			return nil, err
		}
		if iss != expectedIssuer {
			return nil, errors.New("invalid issuer")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ExtractBearerToken extracts the token from the Authorization header.
func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) || len(authHeader) <= len(prefix) {
		return "", errors.New("invalid authorization header format")
	}

	return authHeader[len(prefix):], nil
}
