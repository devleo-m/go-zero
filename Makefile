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

# Comando padrão
.DEFAULT_GOAL := help