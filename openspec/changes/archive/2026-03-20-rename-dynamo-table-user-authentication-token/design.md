## Context

The `validate-jti-from-dynamodb` change introduced DynamoDB session lookups using a new table named `user-auth-tokens`. During implementation the table name was updated to `user-authentication-token` in `config.go` and `Makefile`, but the `.env` / `.env.example` files and the openspec documentation artifacts were not updated, leaving an inconsistency.

Current state:
- `internal/config/config.go` default: `user-authentication-token` ✓
- `Makefile` `SESSION_TABLE_NAME`: `user-authentication-token` ✓
- `.env` `SESSION_TABLE_NAME`: `user-auth-tokens` ✗
- `.env.example` `SESSION_TABLE_NAME`: `user-auth-tokens` ✗
- `openspec/changes/validate-jti-from-dynamodb/` docs: `user-auth-tokens` throughout ✗

## Goals / Non-Goals

**Goals:**
- Align `.env` and `.env.example` with the table name already used in code
- Correct all documentation references in the `validate-jti-from-dynamodb` openspec change

**Non-Goals:**
- Changing any Go source code (already correct)
- Migrating or recreating DynamoDB tables (table name change applies to new/local environments only)
- Updating Terraform or CI configuration (not affected)

## Decisions

**D1: Update env files directly, no migration script needed**

Since `SESSION_TABLE_NAME` is an env var override, the default in `config.go` already uses `user-authentication-token`. Changing `.env` and `.env.example` is a simple string replacement with no ordering requirements.

**D2: Update openspec docs in-place**

The `validate-jti-from-dynamodb` change is archived/complete. Updating its artifacts is a documentation correction, not a new change. No re-apply is needed.

## Risks / Trade-offs

- **[Risk]** Developers with existing LocalStack environments may have a table named `user-auth-tokens` — **Mitigation**: `make dynamodb-create-table` uses `SESSION_TABLE_NAME` from environment; running it with the updated `.env` will create the correctly named table.
- **[Risk]** `.env` is gitignored; this change only affects `.env.example` in the repo — **Mitigation**: Document in tasks to update both files; developers must re-copy or manually edit their local `.env`.
