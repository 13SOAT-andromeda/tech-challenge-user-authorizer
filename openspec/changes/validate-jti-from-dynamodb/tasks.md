## 1. Update Store Interface and DynamoDB Implementation

- [x] 1.1 In `internal/session/store.go`, rename interface method from `GetSessionByUserID` to `GetSessionByJTI(ctx context.Context, jti string) (*Session, error)`
- [x] 1.2 Replace `getItem` logic: use `token_id` (String) as the GetItem key instead of `userId` (N/S fallback)
- [x] 1.3 Read `user_id` attribute (String) from the DynamoDB item to populate `Session.UserID`
- [x] 1.4 Remove the N/S type-fallback logic for `userId` (no longer needed)
- [x] 1.5 Rename the `DynamoStore` method to `GetSessionByJTI` and wire it to the new `getItem` logic

## 2. Update Handler in cmd/authorizer/main.go

- [x] 2.1 Replace `sessionStore.GetSessionByUserID(ctx, userID)` call with `sessionStore.GetSessionByJTI(ctx, tokenJTI)`
- [x] 2.2 After fetching session by JTI, compare `activeSession.UserID` against `userID` from JWT claims
- [x] 2.3 Remove the now-redundant `activeSession.JTI != tokenJTI` check (JTI was used as the lookup key — a match is implicit)

## 3. Update Config and Environment

- [x] 3.1 In `internal/config/config.go`, change the default value of `SessionTableName` from `user-sessions` to `user-auth-tokens`
- [x] 3.2 Update `.env` — set `SESSION_TABLE_NAME=user-auth-tokens`
- [x] 3.3 Update `.env.example` — set `SESSION_TABLE_NAME=user-auth-tokens`

## 4. Update Makefile and Local Dev Scripts

- [x] 4.1 Update `make dynamodb-create-table` to create `user-auth-tokens` with `token_id` (String) as the partition key
- [x] 4.2 Update `make dynamodb-put-session` to write a `token_id` (JTI) + `user_id` + `expires_at` item
- [x] 4.3 Update `make dynamodb-get-session` to use `token_id` key
- [x] 4.4 Update `make dynamodb-delete-session` to use `token_id` key
- [x] 4.5 Update `DYNAMODB_TABLE_NAME` / `SESSION_TABLE_NAME` variable references in Makefile

## 5. Fix Tests

- [x] 5.1 In `cmd/authorizer/main_test.go`, rename mock method from `GetSessionByUserID` to `GetSessionByJTI` in all mock implementations
- [x] 5.2 Update test cases that previously seeded sessions by `userId` to seed by `jti` / `token_id`
- [x] 5.3 Add or update test case: valid JTI in DynamoDB + matching `user_id` → 200
- [x] 5.4 Add or update test case: valid JTI in DynamoDB + mismatched `user_id` → 401
- [x] 5.5 Add or update test case: JTI not in DynamoDB → 401

## 6. Build and Smoke Test

- [x] 6.1 Run `go build ./...` and confirm zero compilation errors
- [x] 6.2 Run `go test ./...` and confirm all tests pass
- [ ] 6.3 Run `make dynamodb-bootstrap` against LocalStack to create `user-auth-tokens` with correct schema
- [ ] 6.4 Login via `tech-challenge-user-authentication` Lambda to obtain a real JWT (with JTI written to DynamoDB)
- [ ] 6.5 Build and deploy authorizer: `make build && make deploy`
- [ ] 6.6 Invoke authorizer with the JWT from step 6.4 and confirm HTTP 200
- [ ] 6.7 Confirm HTTP 401 when invoking with an expired or revoked token (delete the session item and retry)
