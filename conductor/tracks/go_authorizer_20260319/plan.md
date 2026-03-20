# Implementation Plan: Go-based Lambda Authorizer

This plan outlines the steps for implementing a serverless authorization function in Go with JWT validation and DynamoDB session check.

## Phase 1: Foundations & Environment [checkpoint: 81920a6]
- [x] Task: Project Scaffolding 81920a6
    - [x] Initialize a new Go module (`go mod init`).
    - [x] Create a project directory structure following best practices (`cmd/`, `internal/`, `pkg/`).
    - [x] Configure basic logging and environment variable handling.
- [x] Task: Go Lambda Setup 81920a6
    - [x] Create a basic Lambda handler in `cmd/authorizer/main.go`.
    - [x] Implement the Lambda entry point and basic response structure (HTTP 200/401).
- [x] Task: DynamoDB Client Configuration 81920a6
    - [x] Initialize the AWS SDK for Go.
    - [x] Configure the DynamoDB client with appropriate timeouts and retries.
- [x] Task: Conductor - User Manual Verification 'Phase 1: Foundations & Environment' (Protocol in workflow.md) 81920a6

## Phase 2: JWT Validation & Session Lookup (TDD) [checkpoint: 5add1f9]
- [x] Task: Implement JWT Parsing & Expiration Check 5add1f9
    - [x] Write failing unit tests for JWT parsing and `exp` claim validation (Red Phase).
    - [x] Implement the logic to extract and validate the JWT from the `Authorization` header (Green Phase).
    - [x] Refactor the parsing logic and ensure all tests pass.
- [x] Task: Implement DynamoDB Session Verification 5add1f9
    - [x] Create an interface for the DynamoDB session lookup to enable mocking.
    - [x] Write failing unit tests for the session lookup, including success and failure cases (Red Phase).
    - [x] Implement the `GetSession` function to query the `sessions` table by `user_id_session` (Green Phase).
    - [x] Refactor the session lookup logic and ensure all tests pass.
- [x] Task: Integrate Authorization Logic into Lambda Handler 5add1f9
    - [x] Write unit tests for the complete authorizer flow (Red Phase).
    - [x] Integrate JWT validation and session check into the Lambda handler (Green Phase).
    - [x] Refactor the handler logic and ensure all tests pass.
- [x] Task: Conductor - User Manual Verification 'Phase 2: JWT Validation & Session Lookup (TDD)' (Protocol in workflow.md) 5add1f9

## Phase 3: Infrastructure & Final Verification
- [ ] Task: Define Infrastructure with Terraform
    - [ ] Create a `terraform/` directory and define the Lambda function resource.
    - [ ] Configure IAM roles and policies for the Lambda function (e.g., CloudWatch Logs, DynamoDB Read).
    - [ ] Define the `sessions` table name as an environment variable for the Lambda function.
- [ ] Task: Local Verification & Cleanup
    - [ ] Perform a local execution of the Lambda handler with sample events (valid/invalid tokens).
    - [ ] Finalize documentation and update the README if necessary.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Infrastructure & Final Verification' (Protocol in workflow.md)
