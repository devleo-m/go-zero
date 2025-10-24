# Makefile para GO ZERO - Repository Genérico
.PHONY: help run-examples test lint clean

# Cores para output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
PURPLE=\033[0;35m
CYAN=\033[0;36m
WHITE=\033[1;37m
NC=\033[0m # No Color

# Configurações
GO_VERSION := $(shell go version | cut -d' ' -f3)
PROJECT_NAME := go-zero

help: ## Mostra esta ajuda
	@echo "$(CYAN)🚀 GO ZERO - Repository Genérico Profissional$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@echo ""
	@echo "$(WHITE)Comandos disponíveis:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(YELLOW)Exemplos de uso:$(NC)"
	@echo "  make test           # Executa todos os testes"
	@echo "  make lint           # Executa linter"
	@echo "  make clean          # Limpa arquivos temporários"


test: ## Executa todos os testes
	@echo "$(BLUE)🧪 Executando testes...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go test -v -cover ./...

test-unit: ## Executa apenas testes unitários
	@echo "$(BLUE)🧪 Executando testes unitários...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go test -v -cover ./internal/domain/...

test-integration: ## Executa apenas testes de integração
	@echo "$(BLUE)🧪 Executando testes de integração...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go test -v -cover ./tests/integration/...

test-e2e: ## Executa apenas testes end-to-end
	@echo "$(BLUE)🧪 Executando testes end-to-end...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go test -v -cover ./tests/e2e/...

lint: ## Executa linter
	@echo "$(BLUE)🔍 Executando linter...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@golangci-lint run

lint-fix: ## Executa linter e corrige problemas automaticamente
	@echo "$(BLUE)🔧 Executando linter e corrigindo problemas...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@golangci-lint run --fix

format: ## Formata código Go
	@echo "$(BLUE)🎨 Formatando código...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go fmt ./...

vet: ## Executa go vet
	@echo "$(BLUE)🔍 Executando go vet...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go vet ./...

clean: ## Limpa arquivos temporários
	@echo "$(BLUE)🧹 Limpando arquivos temporários...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go clean
	@rm -rf coverage.out
	@rm -rf tmp/
	@rm -rf .air/

build: ## Compila o projeto
	@echo "$(BLUE)🔨 Compilando projeto...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go build -o bin/$(PROJECT_NAME) ./cmd/api

run: ## Executa o projeto
	@echo "$(BLUE)🚀 Executando projeto...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go run ./cmd/api

dev: ## Executa em modo desenvolvimento com hot reload
	@echo "$(BLUE)🔥 Executando em modo desenvolvimento...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@air

docker-up: ## Sobe containers Docker
	@echo "$(BLUE)🐳 Subindo containers Docker...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@docker-compose up -d

docker-down: ## Para containers Docker
	@echo "$(BLUE)🐳 Parando containers Docker...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@docker-compose down

docker-logs: ## Mostra logs dos containers
	@echo "$(BLUE)📋 Mostrando logs dos containers...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@docker-compose logs -f

migrate-up: ## Executa migrations para cima
	@echo "$(BLUE)📊 Executando migrations...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@migrate -path database/migrations -database "postgres://postgres:postgres@localhost:5432/go_zero?sslmode=disable" up

migrate-down: ## Executa migrations para baixo
	@echo "$(BLUE)📊 Revertendo migrations...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@migrate -path database/migrations -database "postgres://postgres:postgres@localhost:5432/go_zero?sslmode=disable" down

migrate-create: ## Cria nova migration (use: make migrate-create NAME=nome_da_migration)
	@echo "$(BLUE)📊 Criando nova migration...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@migrate create -ext sql -dir database/migrations -seq $(NAME)

# Seeds
seed: ## Executa todos os seeds
	@echo "$(BLUE)🌱 Executando seeds...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go run cmd/seed/main.go

seed-users: ## Executa apenas seed de usuários
	@echo "$(BLUE)👥 Executando seed de usuários...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go run cmd/seed/main.go -action=users

seed-clear: ## Limpa todos os dados seedados
	@echo "$(BLUE)🧹 Limpando dados seedados...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go run cmd/seed/main.go -action=clear

seed-clear-users: ## Limpa apenas dados de usuários seedados
	@echo "$(BLUE)🧹 Limpando dados de usuários seedados...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go run cmd/seed/main.go -action=clear-users

dev-setup: ## Configura ambiente de desenvolvimento (migrations + seeds)
	@echo "$(BLUE)🚀 Configurando ambiente de desenvolvimento...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@echo "$(GREEN)1. Executando migrations...$(NC)"
	@make migrate-up
	@echo "$(GREEN)2. Executando seeds...$(NC)"
	@make seed
	@echo "$(GREEN)✅ Ambiente configurado e populado!$(NC)"

coverage: ## Gera relatório de cobertura
	@echo "$(BLUE)📊 Gerando relatório de cobertura...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Relatório de cobertura gerado: coverage.html$(NC)"

benchmark: ## Executa benchmarks
	@echo "$(BLUE)⚡ Executando benchmarks...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go test -bench=. -benchmem ./...

deps: ## Instala dependências
	@echo "$(BLUE)📦 Instalando dependências...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go mod download
	@go mod tidy

deps-update: ## Atualiza dependências
	@echo "$(BLUE)📦 Atualizando dependências...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go get -u ./...
	@go mod tidy

install-tools: ## Instala ferramentas de desenvolvimento
	@echo "$(BLUE)🛠️ Instalando ferramentas de desenvolvimento...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

docs: ## Gera documentação
	@echo "$(BLUE)📚 Gerando documentação...$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@swag init -g cmd/api/main.go -o docs/swagger

check: lint vet test ## Executa todas as verificações
	@echo "$(GREEN)✅ Todas as verificações passaram!$(NC)"

ci: deps check ## Executa pipeline de CI
	@echo "$(GREEN)✅ Pipeline de CI executado com sucesso!$(NC)"


# Informações do projeto
info: ## Mostra informações do projeto
	@echo "$(CYAN)📋 Informações do Projeto$(NC)"
	@echo "$(YELLOW)================================================$(NC)"
	@echo "$(WHITE)Projeto:$(NC) $(PROJECT_NAME)"
	@echo "$(WHITE)Go Version:$(NC) $(GO_VERSION)"
	@echo ""
	@echo "$(WHITE)Estrutura do Repository Genérico:$(NC)"
	@echo "  📁 internal/domain/shared/"
	@echo "    ├── repository.go          # Interface genérica"
	@echo "    ├── query_filter.go        # Sistema de filtros"
	@echo "    ├── paginated_result.go    # Paginação profissional"
	@echo "    ├── aggregation.go         # Agregações"
	@echo "    ├── transaction.go         # Transações"
	@echo "    └── README.md              # Documentação"
	@echo ""
	@echo "$(WHITE)Comandos de Database:$(NC)"
	@echo "  make migrate-up              # Executa migrations"
	@echo "  make migrate-down            # Reverte migrations"
	@echo "  make migrate-create NAME=... # Cria nova migration"
	@echo "  make seed                    # Executa todos os seeds"
	@echo "  make seed-users              # Executa seed de usuários"
	@echo "  make dev-setup               # Setup completo (migrate + seed)"

# Comandos de desenvolvimento
dev-setup-full: ## Setup completo do ambiente (deps + migrate + seed)
	@echo "$(BLUE)🚀 Setup completo do ambiente...$(NC)"
	@make install-deps
	@make docker-up
	@sleep 10
	@make migrate-up
	@make seed
	@echo "$(GREEN)✅ Ambiente completo configurado!$(NC)"

install-deps: ## Instala todas as dependências e ferramentas
	@echo "$(BLUE)📦 Instalando dependências...$(NC)"
	@go mod download
	@go mod tidy
	@make install-tools

check-deps: ## Verifica se todas as dependências estão instaladas
	@echo "$(BLUE)🔍 Verificando dependências...$(NC)"
	@which air > /dev/null || (echo "$(RED)❌ Air não instalado$(NC)" && exit 1)
	@which golangci-lint > /dev/null || (echo "$(RED)❌ golangci-lint não instalado$(NC)" && exit 1)
	@which migrate > /dev/null || (echo "$(RED)❌ migrate não instalado$(NC)" && exit 1)
	@echo "$(GREEN)✅ Todas as dependências estão instaladas$(NC)"

# Comandos de desenvolvimento
dev-logs: ## Mostra logs da aplicação em desenvolvimento
	@echo "$(BLUE)📋 Logs da aplicação...$(NC)"
	@docker-compose logs -f app

dev-restart: ## Reinicia apenas a aplicação
	@echo "$(BLUE)🔄 Reiniciando aplicação...$(NC)"
	@docker-compose restart app

# Comandos de banco
db-reset: ## Reseta completamente o banco (CUIDADO!)
	@echo "$(RED)⚠️  Resetando banco de dados...$(NC)"
	@make docker-down
	@docker volume rm go-zero_postgres_data 2>/dev/null || true
	@make docker-up
	@sleep 10
	@make migrate-up
	@make seed
	@echo "$(GREEN)✅ Banco resetado e populado$(NC)"

# Comandos de teste
test-watch: ## Executa testes em modo watch
	@echo "$(BLUE)👀 Executando testes em modo watch...$(NC)"
	@air -c .air.test.toml

# Comandos de documentação
docs-serve: ## Serve documentação localmente
	@echo "$(BLUE)📚 Servindo documentação...$(NC)"
	@swag init -g cmd/api/main.go -o docs/swagger
	@echo "$(GREEN)✅ Documentação disponível em /swagger$(NC)"

# Comandos de monitoramento
logs-all: ## Mostra logs de todos os serviços
	@echo "$(BLUE)📋 Logs de todos os serviços...$(NC)"
	@docker-compose logs -f

logs-db: ## Mostra logs do banco de dados
	@echo "$(BLUE)📋 Logs do banco de dados...$(NC)"
	@docker-compose logs -f db

logs-redis: ## Mostra logs do Redis
	@echo "$(BLUE)📋 Logs do Redis...$(NC)"
	@docker-compose logs -f redis

# Comandos de limpeza
clean-all: ## Limpa tudo (containers, volumes, imagens)
	@echo "$(BLUE)🧹 Limpando tudo...$(NC)"
	@make docker-down
	@docker system prune -f
	@docker volume prune -f
	@make clean
	@echo "$(GREEN)✅ Limpeza completa realizada$(NC)"

# Comandos de status
status: ## Mostra status dos serviços
	@echo "$(BLUE)📊 Status dos serviços...$(NC)"
	@docker-compose ps

health: ## Verifica saúde dos serviços
	@echo "$(BLUE)🏥 Verificando saúde dos serviços...$(NC)"
	@curl -s http://localhost:8080/health || echo "$(RED)❌ API não está respondendo$(NC)"
	@docker-compose ps | grep -q "Up" && echo "$(GREEN)✅ Serviços estão rodando$(NC)" || echo "$(RED)❌ Alguns serviços estão parados$(NC)"

# Comandos de backup
backup-db: ## Faz backup do banco de dados
	@echo "$(BLUE)💾 Fazendo backup do banco...$(NC)"
	@docker-compose exec -T db pg_dump -U postgres go_zero > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)✅ Backup criado$(NC)"

# Comandos de desenvolvimento avançado
dev-shell: ## Abre shell no container da aplicação
	@echo "$(BLUE)🐚 Abrindo shell no container...$(NC)"
	@docker-compose exec app sh

dev-db-shell: ## Abre shell no banco de dados
	@echo "$(BLUE)🐚 Abrindo shell no banco...$(NC)"
	@docker-compose exec db psql -U postgres -d go_zero

# Comando padrão
.DEFAULT_GOAL := help