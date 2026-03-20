## Why

The authorizer currently validates sessions by looking up DynamoDB using `userId` as the partition key (table `user-sessions`). The authentication lambda (`tech-challenge-user-authentication`) stores sessions in a different table (`user-auth-tokens`) with the **JTI (UUID) as the partition key** and `user_id` as a regular attribute. Because the two services use incompatible schemas and table names, every authorization request fails with "session not found" — JTI-based revocation is broken end-to-end.

## What Changes

- **Replace `GetSessionByUserID`** with `GetSessionByJTI`: look up the session directly by JTI (`token_id`) from the `user-auth-tokens` table instead of looking up by `userId` from `user-sessions`.
- **Update `Store` interface**: replace `GetSessionByUserID(ctx, userID)` with `GetSessionByJTI(ctx, jti)`.
- **Update `DynamoStore`**: `GetItem` by `token_id` (String) instead of `userId` (Number/String fallback), reading `user_id` and `expires_at` attributes from the result.
- **Update `handler`**: call `GetSessionByJTI(ctx, tokenJTI)` instead of `GetSessionByUserID(ctx, userID)`, then validate `userID` from the session record matches the token claim.
- **Update config/env**: `SESSION_TABLE_NAME` default changes from `user-sessions` to `user-auth-tokens`.
- **Update Makefile & deploy scripts**: align DynamoDB table name and seed commands with `user-auth-tokens` schema.

## Capabilities

### New Capabilities

- `jti-session-lookup`: Retrieve and validate user sessions from DynamoDB using the JTI as the partition key, aligned with the authentication lambda's storage schema.

### Modified Capabilities

<!-- No existing openspec specs — this is net-new spec coverage. -->

## Impact

- **`internal/session/store.go`** — replace `GetSessionByUserID` with `GetSessionByJTI`; key type changes from `userId (N/S)` to `token_id (S)`
- **`cmd/authorizer/main.go`** — call `GetSessionByJTI` instead of `GetSessionByUserID`; remove `userId`-based lookup
- **`internal/config/config.go`** — update default table name from `user-sessions` to `user-auth-tokens`
- **`Makefile`** — update DynamoDB bootstrap targets to use `user-auth-tokens` schema
- **`.env` / `.env.example`** — update `SESSION_TABLE_NAME`
- No API contract changes; no new external dependencies
- **BREAKING**: `Store` interface method renamed — any existing mock implementations must be updated
