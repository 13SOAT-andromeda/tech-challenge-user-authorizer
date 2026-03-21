## Context

Two Lambda functions share a DynamoDB session contract:

| Service | Table | Partition Key | Attributes |
|---|---|---|---|
| `tech-challenge-user-authentication` | `user-authentication-token` | `token_id` (String = JTI UUID) | `user_id` (String), `expires_at` (String) |
| `tech-challenge-user-authorizer` (current) | `user-sessions` | `userId` (Number/String) | `jti` (String) |

The authorizer was built against a different schema than what the auth lambda produces. As a result, `GetSessionByUserID` will never find a record in `user-authentication-token` — the partition key name and type both differ. The fix is to invert the lookup: use the JTI (which the authorizer already extracts from the JWT claims) as the key, and read `user_id` from the stored session to cross-validate.

## Goals / Non-Goals

**Goals:**
- Authorizer reads from `user-authentication-token` using `token_id` (JTI) as the GetItem key
- Authorizer validates that `user_id` stored in DynamoDB matches the `user_id` claim in the JWT
- `Store` interface expresses the correct abstraction (`GetSessionByJTI`)
- Local dev (LocalStack, Makefile) aligned with the new table/schema

**Non-Goals:**
- Session expiry enforcement (DynamoDB TTL or `expires_at` check in code — can be added separately)
- Migrating `user-sessions` data to `user-authentication-token`
- Changing the auth lambda

## Decisions

### D1: Look up by JTI, not by userID
**Decision**: `GetItem` using `token_id = jti` (String partition key).

**Why**: The auth lambda writes sessions with JTI as the PK. The authorizer already has the JTI from the JWT claims before making any DynamoDB call. A direct GetItem by JTI is O(1) and removes the fragile N/S type-fallback logic.

**Alternative considered**: Query by `user_id` as a GSI — rejected because it requires a new GSI on `user-authentication-token`, adds latency, and is unnecessary when the JTI is already available.

---

### D2: Validate userID from the stored session record
**Decision**: After fetching the session by JTI, compare `session.UserID` (from DynamoDB `user_id` attribute) against the `user_id` claim in the JWT.

**Why**: Prevents a token with a valid JTI but wrong `user_id` claim from being authorized. Defense-in-depth.

---

### D3: Keep AWS SDK v1 (`aws-sdk-go`)
**Decision**: Do not migrate to AWS SDK v2 as part of this change.

**Why**: The authorizer uses SDK v1 throughout. Migrating would be a separate, larger refactor unrelated to the session lookup bug.

---

### D4: Update default `SESSION_TABLE_NAME` to `user-authentication-token`
**Decision**: Change the default in `config.go` and `.env`.

**Why**: Zero-config alignment — developers can run locally without setting an extra env var. The old default (`user-sessions`) points to a table that no longer exists in the system.

## Risks / Trade-offs

- **Risk**: Existing tests mock `GetSessionByUserID` — they will break at compile time → **Mitigation**: Update all mocks in `main_test.go` and any store tests as part of this change.
- **Risk**: `user-authentication-token` doesn't exist in the authorizer's LocalStack environment → **Mitigation**: Update `make dynamodb-bootstrap` to create the table with the correct schema.
- **Risk**: `expires_at` is stored as a String (ISO 8601) — not enforced at read time → **Mitigation**: Acceptable for now; TTL enforcement is a separate task.

## Migration Plan

1. Update `internal/session/store.go` — new interface method + GetItem by `token_id`
2. Update `cmd/authorizer/main.go` — call `GetSessionByJTI`
3. Update `internal/config/config.go` — default table name
4. Update `.env`, `.env.example`, `Makefile`
5. Fix tests
6. Build, deploy to LocalStack, run end-to-end smoke test using a JWT from the auth lambda
