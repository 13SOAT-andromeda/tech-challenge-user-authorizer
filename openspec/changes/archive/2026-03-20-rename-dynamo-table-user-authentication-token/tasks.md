## 1. Fix environment files

- [x] 1.1 In `.env`, update `SESSION_TABLE_NAME` from `user-auth-tokens` to `user-authentication-token`
- [x] 1.2 In `.env.example`, update `SESSION_TABLE_NAME` from `user-auth-tokens` to `user-authentication-token`

## 2. Fix openspec documentation artifacts

- [x] 2.1 In `openspec/changes/validate-jti-from-dynamodb/proposal.md`, replace all occurrences of `user-auth-tokens` with `user-authentication-token`
- [x] 2.2 In `openspec/changes/validate-jti-from-dynamodb/design.md`, replace all occurrences of `user-auth-tokens` with `user-authentication-token`
- [x] 2.3 In `openspec/changes/validate-jti-from-dynamodb/specs/jti-session-lookup/spec.md`, replace all occurrences of `user-auth-tokens` with `user-authentication-token`
- [x] 2.4 In `openspec/changes/validate-jti-from-dynamodb/tasks.md`, replace all occurrences of `user-auth-tokens` with `user-authentication-token`

## 3. Verify correctness

- [x] 3.1 Run `grep -r "user-auth-tokens" .` and confirm zero matches remain (excluding this change's own docs)
- [x] 3.2 Run `go test ./...` to confirm no tests are broken
