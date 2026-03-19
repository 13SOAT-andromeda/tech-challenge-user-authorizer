# tech-challenge-user-authentication

Authentication and authorization API for the Tech Challenge platform.

## Local Development and Testing

This project uses LocalStack to simulate the AWS environment locally.

### Prerequisites

- Docker and Docker Compose
- Go 1.22+
- AWS CLI

### Getting Started

1. **Start LocalStack:**
   ```bash
   docker compose up -d
   ```

2. **Deploy to LocalStack:**
   This script builds the Go binary, packages it, and deploys it to LocalStack.
   ```bash
   ./scripts/deploy_local.sh
   ```

3. **Test the Lambda:**
   This script generates a valid JWT and invokes the Lambda in LocalStack.
   ```bash
   ./scripts/test_local.sh
   ```

### Project Structure

- `cmd/authorizer`: Lambda entry point.
- `internal/auth`: JWT validation logic.
- `internal/config`: Configuration management.
- `pkg/utils`: Logging and utilities.
- `scripts/`: Helper scripts for local development.
- `terraform/`: Infrastructure as Code (Production).
