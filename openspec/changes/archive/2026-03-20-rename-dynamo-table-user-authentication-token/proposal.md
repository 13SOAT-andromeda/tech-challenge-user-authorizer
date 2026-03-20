## Why

The DynamoDB table name `user-auth-tokens` is inconsistent with the rest of the project, which already partially migrated to `user-authentication-token` (config default and Makefile). The `.env` and `.env.example` files still reference the old name, creating a mismatch that can cause runtime failures in local development.

## What Changes

- **Update `.env`**: change `SESSION_TABLE_NAME` from `user-auth-tokens` to `user-authentication-token`
- **Update `.env.example`**: change `SESSION_TABLE_NAME` from `user-auth-tokens` to `user-authentication-token`
- **Update openspec artifacts** in `validate-jti-from-dynamodb`: correct all remaining references to the old table name in proposal.md, design.md, specs, and tasks.md

## Capabilities

### New Capabilities
<!-- none -->

### Modified Capabilities
- `jti-session-lookup`: Table name referenced in the spec changes from `user-auth-tokens` to `user-authentication-token`

## Impact

- **`.env`** — local dev environment variable corrected
- **`.env.example`** — template for new developers corrected
- **`openspec/changes/validate-jti-from-dynamodb/`** — documentation artifacts updated for consistency
- No code changes required; `internal/config/config.go` and `Makefile` already use the correct name
