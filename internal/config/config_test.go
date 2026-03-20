package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	os.Setenv("JWT_ISSUER", "test_issuer")
	os.Setenv("SESSION_TABLE_NAME", "sessions-test")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("JWT_ISSUER")
	defer os.Unsetenv("SESSION_TABLE_NAME")

	cfg := LoadConfig()

	if cfg.JWTSecret != "test_secret" {
		t.Errorf("expected JWTSecret to be 'test_secret', got '%s'", cfg.JWTSecret)
	}

	if cfg.JWTIssuer != "test_issuer" {
		t.Errorf("expected JWTIssuer to be 'test_issuer', got '%s'", cfg.JWTIssuer)
	}

	if cfg.SessionTableName != "sessions-test" {
		t.Errorf("expected SessionTableName to be 'sessions-test', got '%s'", cfg.SessionTableName)
	}
}

func TestGetEnv_Existing(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	val := getEnv("TEST_KEY", "default_value")

	if val != "test_value" {
		t.Errorf("expected 'test_value', got '%s'", val)
	}
}

func TestGetEnv_Default(t *testing.T) {
	val := getEnv("NON_EXISTENT_KEY", "default_value")

	if val != "default_value" {
		t.Errorf("expected 'default_value', got '%s'", val)
	}
}
