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

2. **Deploy to LocalStack (Makefile):**
   Na primeira vez o alvo cria a função; depois só atualiza o código e as variáveis de ambiente.
   ```bash
   make deploy
   ```
   Equivalente manual: `make create` (primeira vez) e depois `make update-code` / `make update-env`.

   Ou use o script:
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
- `internal/config`: Configuration management (`JWT_*`, `SESSION_TABLE_NAME`; see `.env.example`).
- `internal/session`: **DynamoDB** access for active session (`jti` + `userId`). Implementation: `internal/session/store.go`.
- `pkg/utils`: Logging and utilities.
- `scripts/`: Helper scripts for local development.
- `terraform/`: Infrastructure as Code (Production).
- `Makefile`: alvos `dynamodb-*` para criar tabela e sessão no LocalStack.

### DynamoDB (sessão ativa)

Depois de validar o JWT, o authorizer consulta uma tabela DynamoDB e confere se o **`jti`** e o **`userId`** do item batem com os claims do token (`jti`, `user_id` / `sub`).

| Onde | O quê |
|------|--------|
| Código | `internal/session/store.go` — `NewDynamoStore`, `GetSessionByUserID` |
| Config da tabela | variável de ambiente `SESSION_TABLE_NAME` (padrão `user-sessions`) |
| Região AWS | `AWS_REGION` (padrão `us-east-1` no client) |
| Endpoint customizado (LocalStack / dev) | `DYNAMODB_ENDPOINT` — se **vazio**, o SDK usa o endpoint público da AWS (Dynamo real) |

**Item esperado na tabela** (PK `userId` numérica no Makefile local):

| Atributo | Valor (deve bater com o JWT) |
|----------|------------------------------|
| `userId` | mesmo identificador que `user_id` ou `sub` no token |
| `jti` | mesmo valor do claim `jti` |

Aliases de nome de atributo aceitos na leitura: `jti` / `Jti` / `JTI`, `userId` / `UserId` / `user_id`.

**LocalStack**

1. Suba os serviços: `docker compose up -d` (DynamoDB habilitado no `docker-compose.yml`).
2. Crie a tabela e uma sessão de teste:
   ```bash
   make dynamodb-bootstrap USER_ID=1 JTI=<jti-do-seu-jwt>
   ```
3. Para a **Lambda** enxergar o Dynamo do LocalStack, defina na variável de ambiente da função o endpoint acessível **de dentro do runtime** (ex.: `sam local` / container):
   - Windows / macOS / WSL com Docker Desktop: muitas vezes `http://host.docker.internal:4566`
   - Linux puro: use o IP do host na rede Docker (ex.: `http://172.17.0.1:4566`) ou o hostname do serviço LocalStack na mesma rede.

Copie `.env.example` para `.env` e ajuste `JWT_*`, `SESSION_TABLE_NAME`, `DYNAMODB_ENDPOINT` e `AWS_REGION` conforme o ambiente.

### Monitoramento em tempo real (CloudWatch Logs + container)

O `docker-compose.yml` habilita **`cloudwatch`** (métricas/alarmes via API) e **`logs`** (CloudWatch Logs — grupos/streams onde a Lambda publica). Assim você monitora:

| Comando | O que mostra |
|---------|----------------|
| `make logs` ou `make cloudwatch-tail` | **CloudWatch Logs** no LocalStack: eventos do log group da Lambda (`aws logs tail --follow`). |
| `make container-logs` | **Container** LocalStack: `stdout`/`stderr` em tempo real (`docker compose logs -f`). Ajuste verbosidade com `LS_LOG` no compose (ex.: `LS_LOG=info`). |

**CloudWatch Logs (Lambda):**

```bash
make cloudwatch-tail
# equivalente: make logs
```

Grupo padrão: `/aws/lambda/user-authorizer` (sobrescreva com `make logs LOG_GROUP=/caminho/outro`).

**Logs do container (Docker):**

```bash
make container-logs
```

Requer **AWS CLI v2** para `make logs` / `make cloudwatch-tail`.

Se aparecer **`Service 'logs' is not enabled`**:

1. O `docker-compose.yml` já inclui `logs` em `SERVICES` — é preciso **recriar** o container:
   ```bash
   make localstack-restart
   ```
2. Se ainda falhar, o diretório persistido `./localstack-data` pode estar com config antiga:
   ```bash
   rm -rf localstack-data
   make localstack-restart
   ```
3. Confirme que o container em execução usa o compose atual: `docker compose ps`.
