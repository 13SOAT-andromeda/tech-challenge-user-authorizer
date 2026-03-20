## ADDED Requirements

### Requirement: Session is looked up by JTI from DynamoDB
The authorizer SHALL retrieve the session record from DynamoDB using the JTI extracted from the JWT (`jti` claim) as the lookup key (`token_id` partition key), against the `user-auth-tokens` table.

#### Scenario: Valid JTI with existing session is authorized
- **WHEN** a request arrives with a valid JWT whose `jti` exists as `token_id` in DynamoDB, without remove the actual jwt validation
- **THEN** the authorizer SHALL return HTTP 200 `{"message": "Authorized"}`

#### Scenario: JTI not found in DynamoDB is rejected
- **WHEN** a request arrives with a valid JWT whose `jti` does NOT exist in DynamoDB
- **THEN** the authorizer SHALL return HTTP 401 `{"error": "Invalid or expired token"}`

#### Scenario: Missing JTI claim in token is rejected
- **WHEN** a JWT is presented without a `jti` claim
- **THEN** the authorizer SHALL return HTTP 401 `{"error": "Invalid or expired token"}`

---

### Requirement: UserID from session record is validated against JWT claim
After retrieving the session by JTI, the system SHALL compare the `user_id` attribute from the DynamoDB record against the `user_id` (or `sub`) claim in the JWT.

#### Scenario: Matching userID passes validation
- **WHEN** the DynamoDB session's `user_id` matches the JWT `user_id` claim
- **THEN** the authorizer SHALL return HTTP 200

#### Scenario: Mismatched userID is rejected
- **WHEN** the DynamoDB session's `user_id` does NOT match the JWT `user_id` claim
- **THEN** the authorizer SHALL return HTTP 401 `{"error": "Invalid or expired token"}`

---

### Requirement: Store interface exposes GetSessionByJTI
The `Store` interface SHALL declare `GetSessionByJTI(ctx context.Context, jti string) (*Session, error)` as its lookup method. `GetSessionByUserID` SHALL be removed.

#### Scenario: Successful lookup returns Session with JTI and UserID
- **WHEN** `GetSessionByJTI` is called with a JTI that exists in DynamoDB
- **THEN** it SHALL return a `*Session` with `JTI` and `UserID` populated

#### Scenario: Missing item returns ErrSessionNotFound
- **WHEN** `GetSessionByJTI` is called with a JTI that does not exist in DynamoDB
- **THEN** it SHALL return `nil, ErrSessionNotFound`

---

### Requirement: Default session table is user-auth-tokens
The default value of `SESSION_TABLE_NAME` SHALL be `user-auth-tokens` to align with the authentication lambda's table.

#### Scenario: No SESSION_TABLE_NAME env var uses correct default
- **WHEN** `SESSION_TABLE_NAME` is not set in the environment
- **THEN** the authorizer SHALL connect to the `user-auth-tokens` DynamoDB table
