# Implementation Plan: Go-based Lambda Authorizer

This plan outlines the steps for implementing a serverless authorization function in Go with JWT validation and DynamoDB session check.

## Phase 1: Foundations & Environment
- [ ] Task: Project Scaffolding
    - [ ] Initialize a new Go module (`go mod init`).
    - [ ] Create a project directory structure following best practices (`cmd/`, `internal/`, `pkg/`).
    - [ ] Configure basic logging and environment variable handling.
- [ ] Task: Go Lambda Setup
    - [ ] Create a basic Lambda handler in `cmd/authorizer/main.go`.
    - [ ] Implement the Lambda entry point and basic response structure (HTTP 200/401).
- [ ] Task: DynamoDB Client Configuration
    - [ ] Initialize the AWS SDK for Go.
    - [ ] Configure the DynamoDB client with appropriate timeouts and retries.
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Foundations & Environment' (Protocol in workflow.md)

## Phase 2: JWT Validation & Session Lookup (TDD)
- [ ] Task: Implement JWT Parsing & Expiration Check
    - [ ] Write failing unit tests for JWT parsing and `exp` claim validation (Red Phase).
    - [ ] Implement the logic to extract and validate the JWT from the `Authorization` header (Green Phase).
    - [ ] Refactor the parsing logic and ensure all tests pass.
- [ ] Task: Implement DynamoDB Session Verification
    - [ ] Create an interface for the DynamoDB session lookup to enable mocking.
    - [ ] Write failing unit tests for the session lookup, including success and failure cases (Red Phase).
    - [ ] Implement the `GetSession` function to query the `sessions` table by `user_id_session` (Green Phase).
    - [ ] Refactor the session lookup logic and ensure all tests pass.
- [ ] Task: Integrate Authorization Logic into Lambda Handler
    - [ ] Write unit tests for the complete authorizer flow (Red Phase).
    - [ ] Integrate JWT validation and session check into the Lambda handler (Green Phase).
    - [ ] Refactor the handler logic and ensure all tests pass.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: JWT Validation & Session Lookup (TDD)' (Protocol in workflow.md)

## Phase 3: Infrastructure & Final Verification
- [ ] Task: Define Infrastructure with Terraform
    - [ ] Create a `terraform/` directory and define the Lambda function resource.
    - [ ] Configure IAM roles and policies for the Lambda function (e.g., CloudWatch Logs, DynamoDB Read).
    - [ ] Define the `sessions` table name as an environment variable for the Lambda function.
- [ ] Task: Local Verification & Cleanup
    - [ ] Perform a local execution of the Lambda handler with sample events (valid/invalid tokens).
    - [ ] Finalize documentation and update the README if necessary.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Infrastructure & Final Verification' (Protocol in workflow.md)
