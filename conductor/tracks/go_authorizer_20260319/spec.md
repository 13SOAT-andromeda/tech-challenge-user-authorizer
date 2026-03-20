# Track Specification: Go-based Lambda Authorizer

## Overview
This track involves implementing a serverless authorization function in Go (Golang) for the AWS Lambda runtime. The authorizer will validate JWT tokens received via the `Authorization` header and perform a session check against a DynamoDB table.

## Functional Requirements
- **JWT Parsing:** Extract and parse the JWT from the `Authorization: Bearer <token>` header.
- **Expiration Validation:** Verify the token's `exp` claim to ensure it has not expired.
- **Session Lookup (DynamoDB):**
    - Query the `sessions` table in DynamoDB using the `user_id_session` key (extracted from the JWT's `user_id` claim).
    - Validate that a session exists and is active for the given user.
- **Success Response:**
    - Return an HTTP `200 OK` status.
    - Optional JSON body: `{ "authorized": true }`.
- **Failure Response:**
    - Return an HTTP `401 Unauthorized` status for any of the following:
        - Missing or malformed `Authorization` header.
        - Expired JWT token.
        - No valid session found in the DynamoDB `sessions` table.

## Non-Functional Requirements
- **Performance:** Ensure the authorizer executes efficiently to minimize latency for authorized requests.
- **Security:** Use secure coding practices to handle JWT secrets (if applicable) and DynamoDB credentials.
- **Observability:** Implement structured logging for debugging and audit purposes.
- **Runtime:** AWS Lambda running on Amazon Linux 2023.

## Acceptance Criteria
- [ ] Successfully parses a valid JWT token.
- [ ] Correctly identifies and rejects expired tokens with a `401 Unauthorized` status.
- [ ] Successfully performs a DynamoDB lookup and validates the session.
- [ ] Returns a `200 OK` status for authorized requests.
- [ ] Returns a `401 Unauthorized` status for unauthorized requests.
- [ ] The authorizer is implemented in idiomatic Go and includes unit tests with >80% coverage.

## Out of Scope
- Implementation of the `sessions` table in DynamoDB (assumed to already exist or managed separately).
- Complex role-based access control (RBAC) beyond simple session validation.
- Custom IAM policy generation (unless specifically requested in the future).
