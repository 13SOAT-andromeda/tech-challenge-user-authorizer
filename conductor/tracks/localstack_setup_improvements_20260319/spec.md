# Track Specification: LocalStack Setup & Authorizer Improvements

## Overview
This track aims to streamline the local development experience for the User Authorizer and improve the application's production-readiness by adding structured logging and enhanced error handling.

## Functional Requirements
- **Automated Local Setup:**
    - Add a `make setup-local` target to:
        - Start LocalStack containers.
        - Create required DynamoDB tables (if still needed for future use).
        - Build, zip, and deploy the Lambda function to LocalStack.
        - Set up the environment variables.
- **Enhanced JWT Validation:**
    - Ensure robust validation of standard claims (exp, iss).
    - Provide clear, actionable error messages for invalid tokens.
- **Improved Token Testing:**
    - Update `scripts/gen_token.go` to be more flexible (accepting issuer, secret, and claims as input).
    - Add a `make test-auth` target that generates a token and invokes the authorizer.

## Non-Functional Requirements
- **Structured Logging:**
    - Replace the basic `log` package with JSON-formatted structured logging (e.g., using `log/slog` from the Go standard library).
    - Log essential request metadata (path, request ID if available) and detailed error contexts.
- **Error Handling:**
    - Standardize error responses from the authorizer (JSON format).
    - Differentiate between client-side errors (401 Unauthorized) and server-side errors (500 Internal Server Error) with appropriate logs.

## Acceptance Criteria
- [ ] `make setup-local` successfully provisions the entire environment in one command.
- [ ] The authorizer returns JSON-formatted error responses for all failure cases.
- [ ] Application logs are emitted in JSON format and include contextual information.
- [ ] `make test-auth` successfully triggers an end-to-end authorization flow on LocalStack.
- [ ] The authorizer includes unit tests for the new logging and error handling logic.

## Out of Scope
- Re-enabling DynamoDB session validation (at this stage).
- Integrating with a real AWS environment (focused on LocalStack).
