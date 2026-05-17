include .env
export

# ===== CONFIG =====
FUNCTION_NAME?=user-authorizer
RUNTIME=provided.al2023
API_ID?=
STAGE?=dev
AUTHORIZE_PATH?=authorize
LOG_GROUP?=/aws/lambda/$(FUNCTION_NAME)
COMPOSE_SERVICE?=localstack

LAMBDA_ENV_VARS := JWT_SECRET=$(JWT_SECRET),JWT_ISSUER=$(JWT_ISSUER),DYNAMODB_TABLE_NAME=$(DYNAMODB_TABLE_NAME),AWS_REGION=$(AWS_REGION)
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
	dynamodb-delete-session dynamodb-delete-table dynamodb-bootstrap \
	sam-env sam

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
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(FUNCTION_BINARY_NAME) ./cmd/authorizer

# ===== ZIP =====
zip: build ## Generate deployment zip
	@echo "Zipping artifact..."
	zip -q $(FUNCTION_ZIP_NAME) $(FUNCTION_BINARY_NAME)

# ===== CREATE / UPDATE LAMBDA =====
create: zip ## Create Lambda function in LocalStack (primeira vez)
	@$(MAKE) lambda-create-only

lambda-create-only:
	@echo "Creating Lambda $(FUNCTION_NAME)..."
	aws --endpoint-url=$(AWS_ENDPOINT_URL) --region=$(AWS_REGION) lambda create-function \
		--function-name $(FUNCTION_NAME) \
		--runtime $(RUNTIME) \
		--handler $(FUNCTION_BINARY_NAME) \
		--role $(AWS_ROLE_ARN) \
		--zip-file fileb://$(FUNCTION_ZIP_NAME) \
		--environment "Variables={$(LAMBDA_ENV_VARS)}"

lambda-update-code-only:
	@echo "Updating Lambda code..."
	aws --endpoint-url=$(AWS_ENDPOINT_URL) --region=$(AWS_REGION) lambda update-function-code \
		--function-name $(FUNCTION_NAME) \
		--zip-file fileb://$(FUNCTION_ZIP_NAME)

update-code: zip ## Atualiza o zip na Lambda (cria a função se ainda não existir)
	@aws lambda get-function \
		--endpoint-url=$(AWS_ENDPOINT_URL) \
		--region=$(AWS_REGION) \
		--function-name $(FUNCTION_NAME) \
		>/dev/null 2>&1 \
		&& $(MAKE) lambda-update-code-only \
		|| $(MAKE) lambda-create-only

# ===== UPDATE ENV =====
update-env: ## Update Lambda environment variables (função precisa existir)
	@echo "Updating environment variables..."
	aws --endpoint-url=$(AWS_ENDPOINT_URL) --region=$(AWS_REGION) lambda update-function-configuration \
		--function-name $(FUNCTION_NAME) \
		--environment "Variables={$(LAMBDA_ENV_VARS)}"

# ===== DEPLOY (FULL) =====
deploy: update-code update-env ## Build zip, update code and env vars
	@echo "Deploy complete."

# ===== INVOKE =====
invoke: ## Invoke Lambda using event.json
	@echo "Invoking Lambda..."
	@test -f event.json || (echo "event.json not found"; exit 1)
	aws --endpoint-url=$(AWS_ENDPOINT_URL) --region=$(AWS_REGION) lambda invoke \
		--function-name $(FUNCTION_NAME) \
		--payload file://event.json \
		--cli-binary-format raw-in-base64-out \
		response.json
	@echo "Response saved to response.json"

# ===== LOGS (LocalStack / CloudWatch Logs) =====
logs: ## Acompanha logs da função em tempo real (Ctrl+C para sair)
	@aws logs tail "$(LOG_GROUP)" \
		--endpoint-url=$(AWS_ENDPOINT_URL) \
		--region=$(AWS_REGION) \
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

cloudwatch-tail: ## CloudWatch Logs: tail do log group da Lambda (aws logs tail --follow)
	@$(MAKE) logs

container-logs: ## Monitora logs do container em tempo real (docker compose logs -f)
	docker compose logs -f --tail=200 $(COMPOSE_SERVICE)

# ===== DELETE =====
delete: ## Delete Lambda function
	@echo "Deleting Lambda..."
	aws --endpoint-url=$(AWS_ENDPOINT_URL) --region=$(AWS_REGION) lambda delete-function \
		--function-name $(FUNCTION_NAME) || true

# ===== RECREATE =====
recreate: delete create ## Recreate Lambda from scratch

# ===== CURL TEST =====
curl: ## Call API Gateway authorize route (requires API_ID)
	@test -n "$(API_ID)" || (echo "Set API_ID, e.g. make curl API_ID=abc123"; exit 1)
	curl -i -X POST "$(AWS_ENDPOINT_URL)/restapis/$(API_ID)/$(STAGE)/_user_request_/$(AUTHORIZE_PATH)" \
	-H "Authorization: Bearer TEST"

# ===== DynamoDB (LocalStack) =====

dynamodb-create-table: ## Create DYNAMODB_TABLE_NAME with hash key token_id (S)
	@echo "Creating DynamoDB table $(DYNAMODB_TABLE_NAME)..."
	aws dynamodb create-table \
		--endpoint-url=$(AWS_ENDPOINT_URL) \
		--region=$(AWS_REGION) \
		--table-name $(DYNAMODB_TABLE_NAME) \
		--attribute-definitions AttributeName=token_id,AttributeType=S \
		--key-schema AttributeName=token_id,KeyType=HASH \
		--billing-mode PAY_PER_REQUEST \
		|| true

dynamodb-wait-table: ## Wait until table exists and is ACTIVE
	aws dynamodb wait table-exists \
		--endpoint-url=$(AWS_ENDPOINT_URL) \
		--region=$(AWS_REGION) \
		--table-name $(DYNAMODB_TABLE_NAME)

dynamodb-put-session: ## Upsert session: token_id (JTI) + user_id + expires_at
	@echo "Putting session token_id=$(JTI) user_id=$(USER_ID) (must match JWT jti and user_id/sub)..."
	aws dynamodb put-item \
		--endpoint-url=$(AWS_ENDPOINT_URL) \
		--region=$(AWS_REGION) \
		--table-name $(DYNAMODB_TABLE_NAME) \
		--item "{\"token_id\":{\"S\":\"$(JTI)\"},\"user_id\":{\"S\":\"$(USER_ID)\"},\"expires_at\":{\"S\":\"2099-12-31T23:59:59Z\"}}"

dynamodb-get-session: ## Get session by JTI
	aws dynamodb get-item \
		--endpoint-url=$(AWS_ENDPOINT_URL) \
		--region=$(AWS_REGION) \
		--table-name $(DYNAMODB_TABLE_NAME) \
		--key "{\"token_id\":{\"S\":\"$(JTI)\"}}" \
		--consistent-read

dynamodb-delete-session: ## Remove session row for JTI
	aws dynamodb delete-item \
		--endpoint-url=$(AWS_ENDPOINT_URL) \
		--region=$(AWS_REGION) \
		--table-name $(DYNAMODB_TABLE_NAME) \
		--key "{\"token_id\":{\"S\":\"$(JTI)\"}}"

dynamodb-delete-table: ## Delete entire table (LocalStack)
	aws dynamodb delete-table \
		--endpoint-url=$(AWS_ENDPOINT_URL) \
		--region=$(AWS_REGION) \
		--table-name $(DYNAMODB_TABLE_NAME) \
		|| true

dynamodb-bootstrap: dynamodb-create-table dynamodb-wait-table dynamodb-put-session ## Create table + seed session
	@echo "DynamoDB bootstrap done."

# ===== SAM Local =====

sam-env: ## Gera sam-env.json com variáveis do .env
	@echo '{"UserAuthorizerFunction":{' \
		'"JWT_SECRET":"$(JWT_SECRET)",' \
		'"JWT_ISSUER":"$(JWT_ISSUER)",' \
		'"DYNAMODB_TABLE_NAME":"$(DYNAMODB_TABLE_NAME)",' \
		'"DYNAMODB_ENDPOINT":"$(DYNAMODB_ENDPOINT)",' \
		'"AWS_ENDPOINT_URL":"$(AWS_ENDPOINT_URL)",' \
		'"AWS_ACCESS_KEY_ID":"$(AWS_ACCESS_KEY_ID)",' \
		'"AWS_SECRET_ACCESS_KEY":"$(AWS_SECRET_ACCESS_KEY)",' \
		'"AWS_REGION":"$(AWS_REGION)"' \
		'}}' > sam-env.json

sam: sam-env ## Start SAM local API (port $(SAM_PORT))
	sam build
	sam local start-api \
		--port $(SAM_PORT) \
		--docker-network administrative-api \
		--env-vars sam-env.json \
		--debug \
		--warm-containers EAGER

# ===== CLEAN =====
clean: ## Remove build artifacts
	rm -f $(FUNCTION_BINARY_NAME) $(FUNCTION_ZIP_NAME) sam-env.json
