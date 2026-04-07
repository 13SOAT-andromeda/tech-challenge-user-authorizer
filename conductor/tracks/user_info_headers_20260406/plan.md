# Implementation Plan: User Information Headers

## Phase 1: Research & Preparation [checkpoint: 8c8ab75]
- [x] Task: Review `internal/auth/jwt.go` to understand current claim extraction logic
- [x] Task: Review `cmd/authorizer/main_test.go` to identify relevant test cases for authorizer response context
- [x] Task: Conductor - User Manual Verification 'Phase 1: Research & Preparation' (Protocol in workflow.md) 8c8ab75

## Phase 2: Implementation (TDD)
- [ ] Task: Write failing unit tests in `cmd/authorizer/main_test.go` to verify:
    - [ ] Response `Context` contains `userId`, `userRole`, and `userEmail` for a valid token with all claims.
    - [ ] Authorizer returns `isAuthorized: false` when a valid token is missing the `role` claim.
    - [ ] Authorizer returns `isAuthorized: false` when a valid token is missing the `email` claim.
- [ ] Task: Implement extraction of `role` and `email` claims from the JWT in `cmd/authorizer/main.go`
- [ ] Task: Update authorizer logic to deny access if any required claim (`sub`, `role`, `email`) is missing
- [ ] Task: Update the response structure to include `userId`, `userRole`, and `userEmail` in the `Context` map
- [ ] Task: Refactor and ensure all tests pass, including existing ones
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Implementation (TDD)' (Protocol in workflow.md)

## Phase 3: Verification & Quality Gate
- [ ] Task: Run full test suite and verify code coverage >80% (`go test ./... -cover`)
- [ ] Task: Verify that existing session-related checks (JTI, session store) still work as expected
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Verification & Quality Gate' (Protocol in workflow.md)
