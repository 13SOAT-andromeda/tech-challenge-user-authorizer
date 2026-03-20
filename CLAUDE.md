# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

An AWS Lambda authorizer written in Go. It validates JWT tokens and cross-checks the token's `jti` and `userId` against an active session record in DynamoDB before allowing or denying requests.

## Commands

```bash
# Build Lambda binary (linux/amd64, no CGO)
make build

# Run all tests
go test ./...

# Run tests in a specific package
go test ./internal/auth/...
go test ./cmd/authorizer/...

# Run a single test
go test ./cmd/authorizer/... -run TestHandler/Valid_Token_and_Session

# Format code
go fmt ./...

# Local development: start LocalStack
docker compose up -d

# Build + deploy to LocalStack (creates Lambda if missing, updates if exists)
make deploy

# Seed DynamoDB with a test session (USER_ID and JTI must match the JWT you generate)
make dynamodb-bootstrap USER_ID=1 JTI=<jti-from-your-jwt>

# Generate a test token and invoke the Lambda in LocalStack
./scripts/test_local.sh

# Tail Lambda logs via CloudWatch Logs (LocalStack)
make logs
```

## Architecture

The authorizer follows a two-step validation flow:

1. **JWT validation** (`internal/auth/jwt.go`): Parses the `Authorization: Bearer <token>` header, verifies the HMAC-SHA256 signature using `JWT_SECRET`, checks expiry, and validates the `iss` claim against `JWT_ISSUER`.

2. **Active session check** (`internal/session/store.go`): After JWT validation, looks up the DynamoDB record keyed by `userId` (numeric PK). Compares the stored `jti` with the token's `jti` claim — mismatches or missing sessions result in a 401.

### Key design points

- `session.Store` is an interface — in tests, `sessionStore` is swapped out for a `mockSessionStore`. The production implementation is `DynamoStore`.
- `appConfig` and `sessionStore` are package-level vars in `cmd/authorizer/main.go`, initialized lazily on first invocation (Lambda warm start optimization). Tests reset them to `nil` between runs.
- `getUserIDFromClaims` checks `user_id` first (string or float64), then falls back to `sub`. This handles tokens where the user identity is in either field.
- `DynamoStore.GetSessionByUserID` tries a numeric key (`AttributeType=N`) first, then retries with a string key (`AttributeType=S`) on `ValidationException` — supporting both DynamoDB table schemas.

### Environment variables

| Variable | Default | Purpose |
|---|---|---|
| `JWT_SECRET` | — | HMAC signing secret |
| `JWT_ISSUER` | — | Expected `iss` claim value |
| `SESSION_TABLE_NAME` | `user-sessions` | DynamoDB table name |
| `AWS_REGION` | `us-east-1` | AWS SDK region |
| `DYNAMODB_ENDPOINT` | _(empty = real AWS)_ | Override for LocalStack |

Copy `.env.example` to `.env` and set values before running locally.

### Production deployment

CI (`.github/workflows/deploy.yml`) builds a Docker image, pushes to ECR, then runs `terraform apply` from the `terraform/` directory. Triggers on push to `main`.

## Code Style

Follow [Effective Go](https://go.dev/doc/effective_go): `gofmt`-formatted, `MixedCaps` naming, small interfaces, explicit error handling (never discard errors with `_`), `panic` only for unrecoverable situations.
