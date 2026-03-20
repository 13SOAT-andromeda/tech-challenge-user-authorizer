-include .env
export

# ===== CONFIG =====
FUNCTION_NAME=user-authorizer
ENDPOINT=http://localhost:4566
REGION=us-east-1
RUNTIME=provided.al2023
API_ID?=
STAGE?=dev
AUTHORIZE_PATH?=authorize
LOG_GROUP?=/aws/lambda/$(FUNCTION_NAME)
COMPOSE_SERVICE?=localstack

JWT_SECRET=5b9b178c235820c6e69fbf54876bc4df3ffb4f3ab5ec87305b8b42d2481358c3
JWT_ISSUER=tech-challenge-s1
SESSION_TABLE_NAME=user-sessions
# Endpoint do DynamoDB dentro da Lambda (LocalStack). Vazio = AWS real. Ver README / .env.example
DYNAMODB_ENDPOINT?=

# AWS CLI não aceita valor vazio em Variables={...,KEY=} — omitir DYNAMODB_ENDPOINT se vazio
LAMBDA_ENV_VARS := JWT_SECRET=$(JWT_SECRET),JWT_ISSUER=$(JWT_ISSUER),SESSION_TABLE_NAME=$(SESSION_TABLE_NAME),AWS_REGION=$(REGION)
ifneq ($(strip $(DYNAMODB_ENDPOINT)),)
LAMBDA_ENV_VARS := $(LAMBDA_ENV_VARS),DYNAMODB_ENDPOINT=$(DYNAMODB_ENDPOINT)
endif

# DynamoDB seed (override: make dynamodb-put-session USER_ID=1 JTI=your-jwt-jti)
USER_ID?=1
JTI?=test-jti-1

.DEFAULT_GOAL := help
.PHONY: help makehelp build zip create lambda-create-only lambda-update-code-only \
	update-code update-env deploy invoke delete recreate curl logs \
	cloudwatch-tail container-logs localstack-restart \
	dynamodb-create-table dynamodb-wait-table dynamodb-put-session dynamodb-get-session \
	dynamodb-delete-session dynamodb-delete-table dynamodb-bootstrap

help: ## Show available Make targets
	@echo ""
	@echo "Available commands:"
	@echo ""
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z0-9_.-]+:.*## / {printf "  %-20s %s\n", $$1, $$2}' Makefile
	@echo ""

# ===== BUILD =====
makehelp: help ## Alias for help

build: ## Build Lambda bootstrap binary
	@echo "Building Lambda..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap ./cmd/authorizer

# ===== ZIP =====
zip: build ## Generate deployment zip
	@echo "Zipping artifact..."
	zip -q function.zip bootstrap

# ===== CREATE / UPDATE LAMBDA =====
create: zip ## Create Lambda function in LocalStack (primeira vez)
	@$(MAKE) lambda-create-only

lambda-create-only:
	@echo "Creating Lambda $(FUNCTION_NAME)..."
	aws --endpoint-url=$(ENDPOINT) --region=$(REGION) lambda create-function \
		--function-name $(FUNCTION_NAME) \
		--runtime $(RUNTIME) \
		--handler bootstrap \
		--role arn:aws:iam::000000000000:role/lambda-role \
		--zip-file fileb://function.zip \
		--environment "Variables={$(LAMBDA_ENV_VARS)}"

lambda-update-code-only:
	@echo "Updating Lambda code..."
	aws --endpoint-url=$(ENDPOINT) --region=$(REGION) lambda update-function-code \
		--function-name $(FUNCTION_NAME) \
		--zip-file fileb://function.zip

# Se a função não existir no LocalStack, cria em vez de falhar com ResourceNotFoundException
update-code: zip ## Atualiza o zip na Lambda (cria a função se ainda não existir)
	@aws lambda get-function \
		--endpoint-url=$(ENDPOINT) \
		--region=$(REGION) \
		--function-name $(FUNCTION_NAME) \
		>/dev/null 2>&1 \
		&& $(MAKE) lambda-update-code-only \
		|| $(MAKE) lambda-create-only

# ===== UPDATE ENV =====
update-env: ## Update Lambda environment variables (função precisa existir)
	@echo "Updating environment variables..."
	aws --endpoint-url=$(ENDPOINT) --region=$(REGION) lambda update-function-configuration \
		--function-name $(FUNCTION_NAME) \
		--environment "Variables={$(LAMBDA_ENV_VARS)}"

# ===== DEPLOY (FULL) =====
deploy: update-code update-env ## Build zip, update code and env vars
	@echo "Deploy complete."

# ===== INVOKE =====
invoke: ## Invoke Lambda using event.json
	@echo "Invoking Lambda..."
	@test -f event.json || (echo "event.json not found"; exit 1)
	aws --endpoint-url=$(ENDPOINT) --region=$(REGION) lambda invoke \
		--function-name $(FUNCTION_NAME) \
		--payload file://event.json \
		--cli-binary-format raw-in-base64-out \
		response.json
	@echo "Response saved to response.json"

# ===== LOGS (LocalStack / CloudWatch Logs) =====
# Requer AWS CLI v2 (`aws logs tail`). LocalStack precisa do serviço `logs` em SERVICES (ver docker-compose.yml).
logs: ## Acompanha logs da função em tempo real (Ctrl+C para sair)
	@aws logs tail "$(LOG_GROUP)" \
		--endpoint-url=$(ENDPOINT) \
		--region=$(REGION) \
		--follow \
	|| { \
		echo ""; \
		echo "Dica: erro 'logs is not enabled' → o container LocalStack foi iniciado sem o serviço logs ou com estado antigo."; \
		echo "      Rode: make localstack-restart"; \
		echo "      Se continuar: rm -rf localstack-data && make localstack-restart"; \
		exit 1; \
	}

localstack-restart: ## Reinicia LocalStack (--force-recreate) para aplicar SERVICES (incl. logs)
	docker compose down
	docker compose up -d --force-recreate
	@echo "Aguarde o LocalStack ficar pronto (~10–40s), depois: make logs"

# Mesma coisa que `logs` — nome explícito para CloudWatch Logs (API emulada no LocalStack)
cloudwatch-tail: ## CloudWatch Logs: tail do log group da Lambda (aws logs tail --follow)
	@$(MAKE) logs

# Stdout/stderr do container Docker (LocalStack) — complementa cloudwatch-tail para ver o runtime local
container-logs: ## Monitora logs do container em tempo real (docker compose logs -f)
	docker compose logs -f --tail=200 $(COMPOSE_SERVICE)

# ===== DELETE =====
delete: ## Delete Lambda function
	@echo "Deleting Lambda..."
	aws --endpoint-url=$(ENDPOINT) --region=$(REGION) lambda delete-function \
		--function-name $(FUNCTION_NAME) || true

# ===== RECREATE =====
recreate: delete create ## Recreate Lambda from scratch

# ===== CURL TEST =====
curl: ## Call API Gateway authorize route (requires API_ID)
	@test -n "$(API_ID)" || (echo "Set API_ID, e.g. make curl API_ID=abc123"; exit 1)
	curl -i -X POST "$(ENDPOINT)/restapis/$(API_ID)/$(STAGE)/_user_request_/$(AUTHORIZE_PATH)" \
	-H "Authorization: Bearer TEST"

# ===== DynamoDB (LocalStack) =====
# Schema: PK userId (N). Each item MUST include:
#   - userId: same value as JWT user_id or sub (authorizer checks item.userId == token user)
#   - jti: same value as JWT jti claim (authorizer checks item.jti == token jti)
# Aliases also accepted on read: Jti/JTI, UserId/user_id (see internal/session/store.go)
# Set AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY if needed (LocalStack: test / test).

dynamodb-create-table: ## Create SESSION_TABLE_NAME with hash key userId (N)
	@echo "Creating DynamoDB table $(SESSION_TABLE_NAME)..."
	aws dynamodb create-table \
		--endpoint-url=$(ENDPOINT) \
		--region=$(REGION) \
		--table-name $(SESSION_TABLE_NAME) \
		--attribute-definitions AttributeName=userId,AttributeType=N \
		--key-schema AttributeName=userId,KeyType=HASH \
		--billing-mode PAY_PER_REQUEST \
		|| true

dynamodb-wait-table: ## Wait until table exists and is ACTIVE
	aws dynamodb wait table-exists \
		--endpoint-url=$(ENDPOINT) \
		--region=$(REGION) \
		--table-name $(SESSION_TABLE_NAME)

dynamodb-put-session: ## Upsert session mirroring JWT: userId (N) + jti (S)
	@echo "Putting session userId=$(USER_ID) jti=$(JTI) (must match JWT user_id/sub and jti)..."
	aws dynamodb put-item \
		--endpoint-url=$(ENDPOINT) \
		--region=$(REGION) \
		--table-name $(SESSION_TABLE_NAME) \
		--item "{\"userId\":{\"N\":\"$(USER_ID)\"},\"jti\":{\"S\":\"$(JTI)\"}}"

dynamodb-get-session: ## Get session by USER_ID
	aws dynamodb get-item \
		--endpoint-url=$(ENDPOINT) \
		--region=$(REGION) \
		--table-name $(SESSION_TABLE_NAME) \
		--key "{\"userId\":{\"N\":\"$(USER_ID)\"}}" \
		--consistent-read

dynamodb-delete-session: ## Remove session row for USER_ID
	aws dynamodb delete-item \
		--endpoint-url=$(ENDPOINT) \
		--region=$(REGION) \
		--table-name $(SESSION_TABLE_NAME) \
		--key "{\"userId\":{\"N\":\"$(USER_ID)\"}}"

dynamodb-delete-table: ## Delete entire table (LocalStack)
	aws dynamodb delete-table \
		--endpoint-url=$(ENDPOINT) \
		--region=$(REGION) \
		--table-name $(SESSION_TABLE_NAME) \
		|| true

dynamodb-bootstrap: dynamodb-create-table dynamodb-wait-table dynamodb-put-session ## Create table + seed session
	@echo "DynamoDB bootstrap done."