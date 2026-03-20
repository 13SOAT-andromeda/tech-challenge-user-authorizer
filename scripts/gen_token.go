package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "5b9b178c235820c6e69fbf54876bc4df3ffb4f3ab5ec87305b8b42d2481358c3"
	}
	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "tech-challenge-s1"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": issuer,
		"sub": "1",
		"jti": "test-jti-1",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	fmt.Print(tokenString)
}
