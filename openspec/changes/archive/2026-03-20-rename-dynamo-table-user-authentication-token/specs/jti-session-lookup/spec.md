## MODIFIED Requirements

### Requirement: Default session table is user-auth-tokens
The default value of `SESSION_TABLE_NAME` SHALL be `user-authentication-token` to align with the authentication lambda's table.

#### Scenario: No SESSION_TABLE_NAME env var uses correct default
- **WHEN** `SESSION_TABLE_NAME` is not set in the environment
- **THEN** the authorizer SHALL connect to the `user-authentication-token` DynamoDB table

### Requirement: Session is looked up by JTI from DynamoDB
The authorizer SHALL retrieve the session record from DynamoDB using the JTI extracted from the JWT (`jti` claim) as the lookup key (`token_id` partition key), against the `user-authentication-token` table.

#### Scenario: Valid JTI with existing session is authorized
- **WHEN** a request arrives with a valid JWT whose `jti` exists as `token_id` in DynamoDB, without remove the actual jwt validation
- **THEN** the authorizer SHALL return HTTP 200 `{"message": "Authorized"}`

#### Scenario: JTI not found in DynamoDB is rejected
- **WHEN** a request arrives with a valid JWT whose `jti` does NOT exist in DynamoDB
- **THEN** the authorizer SHALL return HTTP 401 `{"error": "Invalid or expired token"}`

#### Scenario: Missing JTI claim in token is rejected
- **WHEN** a JWT is presented without a `jti` claim
- **THEN** the authorizer SHALL return HTTP 401 `{"error": "Invalid or expired token"}`
