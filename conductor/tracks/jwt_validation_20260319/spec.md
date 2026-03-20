# Track Specification: JWT Validation

## Overview
This track involves implementing a robust JWT (JSON Web Token) validation mechanism for the User Authorizer Lambda function written in Go. The validation will focus on verifying tokens received via the `Authorization` header, specifically focusing on symmetric key validation and issuer verification.

## Functional Requirements
- **Bearer Token Extraction:** Extract the JWT from the `Authorization: Bearer <token>` header in the incoming Lambda event.
- **Symmetric Key Validation:** Validate the JWT using a symmetric key (HMAC) stored in an environment variable (e.g., `JWT_SECRET`).
- **Issuer Verification:** Verify that the `iss` (issuer) claim in the JWT matches the expected value defined in the environment (e.g., `JWT_ISSUER`).
- **Error Handling:** Return appropriate error responses (e.g., 401 Unauthorized) if the token is missing, invalid, or expired.
- **Environment Configuration:** Load the secret key and expected issuer from a `.env` file or environment variables.

## Non-Functional Requirements
- **Security:** Ensure the JWT secret is handled securely and not exposed in logs or code.
- **Performance:** Validation should be efficient to minimize Lambda execution time.
- **Reliability:** Handle various edge cases like malformed tokens or missing headers gracefully.

## Acceptance Criteria
- [ ] Successfully extract JWT from the `Authorization` header.
- [ ] Validate JWT signature using the symmetric key from the `.env` file.
- [ ] Verify that the `iss` claim matches the expected issuer.
- [ ] Return a 200 OK or appropriate success response for valid tokens.
- [ ] Return a 401 Unauthorized response for invalid or missing tokens.
- [ ] Unit tests cover both successful and failed validation scenarios.

## Out of Scope
- Integration with DynamoDB (to be handled in a separate track or phase).
- Asymmetric key (RSA/ECDSA) validation.
- User permission/role-based access control (RBAC).
