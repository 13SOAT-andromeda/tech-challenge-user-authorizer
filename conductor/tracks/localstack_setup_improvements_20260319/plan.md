# Implementation Plan: LocalStack Setup & Authorizer Improvements

## Phase 1: Foundation & Application Improvements
This phase focuses on upgrading the application's logging and error handling to be more robust and production-ready.

- [ ] Task: Implement Structured JSON Logging in `pkg/utils/logger.go` using `log/slog`.
    - [ ] Add `Logger` struct or global helper to `pkg/utils/logger.go`.
    - [ ] Implement `Info`, `Error`, and `Warn` methods that output JSON.
    - [ ] Update `pkg/utils/logger_test.go` to verify JSON output.
- [ ] Task: Refactor `cmd/authorizer/main.go` to use the new JSON logger.
    - [ ] Replace `utils.InfoLogger` and `utils.ErrorLogger` calls with the new structured logger.
    - [ ] Add relevant context to log entries (e.g., `request_path`, `error_type`).
- [ ] Task: Enhance Error Handling in `cmd/authorizer/main.go`.
    - [ ] Standardize the JSON structure for error responses (e.g., `{"error": "...", "code": "..."}`).
    - [ ] Ensure proper differentiation between 401 (Unauthorized) and 500 (Internal Server Error) responses.
- [ ] Task: Update `internal/auth/jwt.go` and `internal/auth/jwt_test.go`.
    - [ ] Add more granular error types/messages for JWT validation (e.g., `ErrExpiredToken`, `ErrInvalidIssuer`).
    - [ ] Ensure tests cover these granular error cases.
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Foundation & Application Improvements' (Protocol in workflow.md)

## Phase 2: Automation & Testing
This phase focuses on streamlining the local development workflow with LocalStack and adding end-to-end testing capabilities.

- [ ] Task: Enhance the JWT generator script `scripts/gen_token.go`.
    - [ ] Update the script to accept parameters (or use environment variables) for `user_id`, `jti`, `issuer`, and `secret`.
    - [ ] Ensure it generates tokens compatible with the authorizer's expectations.
- [ ] Task: Update `Makefile` with automation targets.
    - [ ] Implement `make setup-local` to orchestrate `docker compose up`, `dynamodb-create-table`, and `deploy`.
    - [ ] Implement `make test-auth` to generate a token, invoke the Lambda via `aws lambda invoke`, and verify the response.
    - [ ] Add `make clean` to stop containers and remove temporary files (`function.zip`, `response.json`).
- [ ] Task: Final Verification and Documentation.
    - [ ] Run `make setup-local` from a clean state.
    - [ ] Run `make test-auth` and confirm success.
    - [ ] Update `README.md` with instructions for the new Makefile targets.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Automation & Testing' (Protocol in workflow.md)
