# ===========================================
# GO ZERO - MAKEFILE PARA DESENVOLVIMENTO
# ===========================================

# Variáveis
APP_NAME=go-zero
DOCKER_COMPOSE=docker-compose
GO=go

# Cores para output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help
help: ## Mostra esta ajuda
	@echo "$(GREEN)GO ZERO - Comandos Disponíveis$(NC)"
	@echo "=================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

# ===========================================
# COMANDOS DE DOCKER
# ===========================================

.PHONY: docker-up
docker-up: ## Sobe todos os serviços (PostgreSQL + Redis)
	@echo "$(GREEN)🚀 Subindo serviços Docker...$(NC)"
	$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)✅ Serviços subiram com sucesso!$(NC)"
	@echo "$(YELLOW)📊 Serviços disponíveis:$(NC)"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379"

.PHONY: docker-down
docker-down: ## Para todos os serviços
	@echo "$(YELLOW)🛑 Parando serviços Docker...$(NC)"
	$(DOCKER_COMPOSE) down
	@echo "$(GREEN)✅ Serviços parados!$(NC)"

.PHONY: docker-logs
docker-logs: ## Mostra logs dos serviços
	$(DOCKER_COMPOSE) logs -f

.PHONY: docker-restart
docker-restart: ## Reinicia todos os serviços
	@echo "$(YELLOW)🔄 Reiniciando serviços...$(NC)"
	$(DOCKER_COMPOSE) restart
	@echo "$(GREEN)✅ Serviços reiniciados!$(NC)"

.PHONY: docker-status
docker-status: ## Mostra status dos containers
	@echo "$(GREEN)📊 Status dos Containers:$(NC)"
	$(DOCKER_COMPOSE) ps

# ===========================================
# COMANDOS DE DESENVOLVIMENTO
# ===========================================

.PHONY: dev
dev: ## Roda a aplicação em modo desenvolvimento (com Air)
	@echo "$(GREEN)🚀 Iniciando aplicação em modo desenvolvimento...$(NC)"
	@echo "$(YELLOW)💡 Certifique-se de que o arquivo .env existe!$(NC)"
	air -c .air.toml

.PHONY: run
run: ## Roda a aplicação sem hot-reload
	@echo "$(GREEN)🚀 Iniciando aplicação...$(NC)"
	$(GO) run cmd/server/main.go

.PHONY: build
build: ## Compila a aplicação
	@echo "$(GREEN)🔨 Compilando aplicação...$(NC)"
	$(GO) build -o bin/$(APP_NAME) cmd/server/main.go
	@echo "$(GREEN)✅ Aplicação compilada em bin/$(APP_NAME)$(NC)"

# ===========================================
# COMANDOS DE DEPENDÊNCIAS
# ===========================================

.PHONY: deps
deps: ## Baixa dependências Go
	@echo "$(GREEN)📦 Baixando dependências...$(NC)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)✅ Dependências atualizadas!$(NC)"

.PHONY: deps-update
deps-update: ## Atualiza dependências
	@echo "$(GREEN)🔄 Atualizando dependências...$(NC)"
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "$(GREEN)✅ Dependências atualizadas!$(NC)"

# ===========================================
# COMANDOS DE TESTE
# ===========================================

.PHONY: test
test: ## Roda todos os testes
	@echo "$(GREEN)🧪 Executando testes...$(NC)"
	$(GO) test -v ./...

.PHONY: test-coverage
test-coverage: ## Roda testes com coverage
	@echo "$(GREEN)📊 Executando testes com coverage...$(NC)"
	$(GO) test -v -cover ./...

# ===========================================
# COMANDOS DE LIMPEZA
# ===========================================

.PHONY: clean
clean: ## Limpa arquivos temporários
	@echo "$(YELLOW)🧹 Limpando arquivos temporários...$(NC)"
	rm -rf tmp/
	rm -rf bin/
	$(GO) clean
	@echo "$(GREEN)✅ Limpeza concluída!$(NC)"

.PHONY: clean-docker
clean-docker: ## Remove containers e volumes Docker
	@echo "$(YELLOW)🧹 Limpando Docker...$(NC)"
	$(DOCKER_COMPOSE) down -v --remove-orphans
	docker system prune -f
	@echo "$(GREEN)✅ Docker limpo!$(NC)"

# ===========================================
# COMANDOS DE SETUP
# ===========================================

.PHONY: setup
setup: ## Setup inicial do projeto
	@echo "$(GREEN)🛠️  Configurando projeto GO ZERO...$(NC)"
	@if [ ! -f .env ]; then \
		echo "$(YELLOW)📝 Criando arquivo .env...$(NC)"; \
		cp .env.example .env; \
		echo "$(GREEN)✅ Arquivo .env criado! Edite as configurações se necessário.$(NC)"; \
	else \
		echo "$(GREEN)✅ Arquivo .env já existe!$(NC)"; \
	fi
	@echo "$(GREEN)📦 Baixando dependências...$(NC)"
	$(GO) mod download
	@echo "$(GREEN)🚀 Subindo serviços Docker...$(NC)"
	$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)✅ Setup concluído!$(NC)"
	@echo "$(YELLOW)💡 Próximos passos:$(NC)"
	@echo "  1. Edite o arquivo .env se necessário"
	@echo "  2. Execute 'make dev' para iniciar o desenvolvimento"

# ===========================================
# COMANDOS DE VERIFICAÇÃO
# ===========================================

.PHONY: check-db
check-db: ## Verifica conexão com PostgreSQL
	@echo "$(GREEN)🔍 Verificando conexão com PostgreSQL...$(NC)"
	@docker exec $(APP_NAME)-db psql -U postgres -d go_zero_dev -c "SELECT version();" || echo "$(RED)❌ Erro ao conectar com PostgreSQL$(NC)"

.PHONY: check-redis
check-redis: ## Verifica conexão com Redis
	@echo "$(GREEN)🔍 Verificando conexão com Redis...$(NC)"
	@docker exec $(APP_NAME)-redis redis-cli ping || echo "$(RED)❌ Erro ao conectar com Redis$(NC)"

.PHONY: check-all
check-all: ## Verifica todos os serviços
	@echo "$(GREEN)🔍 Verificando todos os serviços...$(NC)"
	@make check-db
	@make check-redis
	@echo "$(GREEN)✅ Verificação concluída!$(NC)"

# ===========================================
# COMANDOS DE DESENVOLVIMENTO AVANÇADO
# ===========================================

.PHONY: shell-db
shell-db: ## Abre shell do PostgreSQL
	@echo "$(GREEN)🐘 Abrindo shell do PostgreSQL...$(NC)"
	docker exec -it $(APP_NAME)-db psql -U postgres -d go_zero_dev

.PHONY: shell-redis
shell-redis: ## Abre shell do Redis
	@echo "$(GREEN)🔴 Abrindo shell do Redis...$(NC)"
	docker exec -it $(APP_NAME)-redis redis-cli

# ===========================================
# COMANDOS DE MIGRATIONS
# ===========================================

.PHONY: migrate-up
migrate-up: ## Aplica todas as migrations
	@echo "$(GREEN)🚀 Aplicando migrations...$(NC)"
	$(GO) run cmd/migrate/main.go -direction=up
	@echo "$(GREEN)✅ Migrations aplicadas!$(NC)"

.PHONY: migrate-down
migrate-down: ## Reverte última migration
	@echo "$(YELLOW)🔄 Revertendo última migration...$(NC)"
	$(GO) run cmd/migrate/main.go -direction=down -steps=1
	@echo "$(GREEN)✅ Migration revertida!$(NC)"

.PHONY: migrate-down-all
migrate-down-all: ## Reverte todas as migrations
	@echo "$(RED)⚠️  Revertendo TODAS as migrations...$(NC)"
	$(GO) run cmd/migrate/main.go -direction=down
	@echo "$(GREEN)✅ Todas as migrations revertidas!$(NC)"

.PHONY: migrate-force
migrate-force: ## Força versão específica (uso: make migrate-force version=1)
	@echo "$(YELLOW)🔧 Forçando versão da migration...$(NC)"
	$(GO) run cmd/migrate/main.go -direction=force -steps=$(version)
	@echo "$(GREEN)✅ Versão forçada!$(NC)"

.PHONY: migrate-version
migrate-version: ## Mostra versão atual das migrations
	@echo "$(GREEN)📊 Verificando versão das migrations...$(NC)"
	$(GO) run cmd/migrate/main.go -direction=up -steps=0
	@echo "$(GREEN)✅ Verificação concluída!$(NC)"

.PHONY: migrate-create
migrate-create: ## Cria nova migration (uso: make migrate-create name=add_phone_to_users)
	@echo "$(GREEN)📝 Criando nova migration...$(NC)"
	@cd internal/infra/database/migrations && migrate create -ext sql -dir . -seq $(name)
	@echo "$(GREEN)✅ Migration criada!$(NC)"
	@echo "$(YELLOW)💡 Edite os arquivos .up.sql e .down.sql$(NC)"

# ===========================================
# COMANDO PADRÃO
# ===========================================

.DEFAULT_GOAL := help
