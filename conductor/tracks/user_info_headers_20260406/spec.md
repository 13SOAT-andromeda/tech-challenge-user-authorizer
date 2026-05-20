# Specification: User Information Headers

## Overview
This track implements the extraction of user information (ID, Role, and Email) from the JWT token and returns this data in the Lambda Authorizer's response. This allows the API Gateway to forward these details as headers to the backend services.

## Functional Requirements
- **Claim Extraction:** Extract the following claims from the validated JWT token:
    - `sub` (Subject) mapping to `userId`.
    - `role` mapping to `userRole`.
    - `email` mapping to `userEmail`.
- **Missing Claims:** If any of these required claims are missing from a valid token, the authorizer must deny the request (`isAuthorized: false`).
- **Context Mapping:** The extracted information must be included in the `Context` map of the `APIGatewayV2CustomAuthorizerSimpleResponse`.
- **Header Convention:** The mapped headers in the API Gateway should be configured as:
    - `X-User-Id` from `context.userId`
    - `X-User-Role` from `context.userRole`
    - `X-User-Email` from `context.userEmail`

## Non-Functional Requirements
- **Performance:** Maintain low latency during token extraction and validation.
- **Security:** Ensure that only valid tokens can produce these headers.
- **Type Safety:** Use robust Go type assertions when extracting claims.

## Acceptance Criteria
- [ ] Authorizer returns `isAuthorized: true` with `userId`, `userRole`, and `userEmail` in the `context` when all claims are present in a valid JWT.
- [ ] Authorizer returns `isAuthorized: false` if `sub`, `role`, or `email` claims are missing, even if the token is otherwise valid.
- [ ] Unit tests cover success cases with all claims and failure cases with missing claims.
- [ ] Existing functionality (JWT validation, session lookup) remains unaffected.

## Out of Scope
- **API Gateway Configuration:** The actual mapping of `context` to HTTP headers in AWS API Gateway is outside the scope of this authorizer code change.
- **New Auth Logic:** No changes to how tokens are validated or how sessions are stored in DynamoDB.
