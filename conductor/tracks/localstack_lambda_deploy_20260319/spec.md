# Specification: LocalStack Lambda & API Gateway Authorizer

## Overview
This track focuses on deploying the existing User Authorizer (Go-based) as an AWS Lambda function within a LocalStack environment. The Lambda will be integrated with API Gateway as a "Lambda Authorizer" to secure simulated platform endpoints.

## Objectives
- Deploy the `cmd/authorizer/` binary to LocalStack as a Lambda function.
- Create a `terraform/localstack.tf` file to manage local AWS resources.
- Update `scripts/deploy_local.sh` to automate the build and deployment process.
- Configure API Gateway to use the Lambda as a Custom Authorizer.

## Functional Requirements
- **Go Build:** Compile the Go source in `cmd/authorizer/` for the `provided.al2023` runtime.
- **Terraform Configuration:**
  - Define an AWS Lambda function resource.
  - Define an API Gateway (REST or HTTP as appropriate for an authorizer).
  - Configure the API Gateway Authorizer to point to the Lambda function.
  - Create a mock/test endpoint in API Gateway to verify authorization.
- **Local Automation:** Update `scripts/deploy_local.sh` to handle the `go build`, `terraform init`, and `terraform apply` sequence targeting LocalStack.

## Non-Functional Requirements
- **Environment:** Must run entirely within the Docker/LocalStack environment defined in `docker-compose.yml`.
- **Consistency:** Use the same Go version and dependencies as the project's standard.
- **Maintainability:** Ensure the `localstack.tf` is clean and well-commented.

## Acceptance Criteria
- [ ] Running `scripts/deploy_local.sh` successfully builds the Go binary.
- [ ] Terraform successfully provisions the Lambda and API Gateway in LocalStack.
- [ ] A `curl` request to the API Gateway test endpoint (with a valid/invalid token) correctly triggers the Lambda Authorizer.
- [ ] The Lambda Authorizer returns an IAM Policy (Allow/Deny) that the API Gateway respects.

## Out of Scope
- Production AWS deployment (this track is strictly for local simulation).
- Extensive JWT validation logic (unless already implemented).
- Advanced IAM permission sets (local development roles only).
