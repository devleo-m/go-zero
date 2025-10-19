# 🚀 COMANDOS RÁPIDOS - GO ZERO

## 📋 **COMANDOS DE MIGRATIONS**

### Aplicar Migrations
```bash
go run cmd/migrate/main.go -direction=up
```
**O que faz:** Aplica todas as migrations pendentes no banco de dados

### Reverter Última Migration
```bash
go run cmd/migrate/main.go -direction=down -steps=1
```
**O que faz:** Reverte apenas a última migration aplicada

### Reverter Todas as Migrations
```bash
go run cmd/migrate/main.go -direction=down
```
**O que faz:** Reverte TODAS as migrations (cuidado!)

### Ver Versão Atual
```bash
go run cmd/migrate/main.go -direction=up -steps=0
```
**O que faz:** Mostra a versão atual das migrations sem aplicar nada

### Criar Nova Migration
```bash
cd internal/infra/database/migrations
migrate create -ext sql -dir . -seq nome_da_migration
cd ../../../..
```
**O que faz:** Cria arquivos .up.sql e .down.sql para nova migration

---

## 🐳 **COMANDOS DOCKER**

### Subir Serviços
```bash
docker-compose up -d
```
**O que faz:** Sobe PostgreSQL e Redis em background

### Parar Serviços
```bash
docker-compose down
```
**O que faz:** Para todos os containers

### Ver Logs
```bash
docker-compose logs -f
```
**O que faz:** Mostra logs em tempo real

### Ver Status
```bash
docker-compose ps
```
**O que faz:** Lista containers e status

---

## 🚀 **COMANDOS DE DESENVOLVIMENTO**

### Rodar Aplicação (com hot-reload)
```bash
air -c .air.toml
```
**O que faz:** Inicia servidor com reload automático

### Rodar Aplicação (sem hot-reload)
```bash
go run cmd/server/main.go
```
**O que faz:** Inicia servidor normalmente

### Compilar Aplicação
```bash
go build -o bin/go-zero cmd/server/main.go
```
**O que faz:** Gera executável em bin/go-zero

---

## 🧪 **COMANDOS DE TESTE**

### Rodar Todos os Testes
```bash
go test -v ./...
```
**O que faz:** Executa todos os testes do projeto

### Rodar Testes com Coverage
```bash
go test -v -cover ./...
```
**O que faz:** Executa testes mostrando cobertura

---

## 🗄️ **COMANDOS DE BANCO**

### Conectar ao PostgreSQL
```bash
docker exec -it go-zero-db psql -U postgres -d go_zero_dev
```
**O que faz:** Abre shell do PostgreSQL

### Conectar ao Redis
```bash
docker exec -it go-zero-redis redis-cli
```
**O que faz:** Abre shell do Redis

### Ver Tabelas do Banco
```bash
docker exec go-zero-db psql -U postgres -d go_zero_dev -c "\dt"
```
**O que faz:** Lista todas as tabelas

---

## 🧹 **COMANDOS DE LIMPEZA**

### Limpar Arquivos Temporários
```bash
go clean
rm -rf bin/
rm -rf tmp/
```
**O que faz:** Remove arquivos compilados e temporários

### Limpar Docker
```bash
docker-compose down -v --remove-orphans
docker system prune -f
```
**O que faz:** Remove containers, volumes e imagens não usadas

---

## 📦 **COMANDOS DE DEPENDÊNCIAS**

### Baixar Dependências
```bash
go mod download
go mod tidy
```
**O que faz:** Baixa e organiza dependências do Go

### Atualizar Dependências
```bash
go get -u ./...
go mod tidy
```
**O que faz:** Atualiza todas as dependências

---

## ⚡ **COMANDOS MAIS USADOS (SEQUÊNCIA)**

### Setup Inicial
```bash
# 1. Subir serviços
docker-compose up -d

# 2. Aplicar migrations
go run cmd/migrate/main.go -direction=up

# 3. Rodar aplicação
go run cmd/server/main.go
```

### Desenvolvimento Diário
```bash
# 1. Aplicar migrations
go run cmd/migrate/main.go -direction=up

# 2. Rodar com hot-reload
air -c .air.toml
```

### Deploy
```bash
# 1. Compilar
go build -o bin/go-zero cmd/server/main.go

# 2. Testar
go test -v ./...

# 3. Rodar
./bin/go-zero
```

---

## 🆘 **COMANDOS DE EMERGÊNCIA**

### Reset Completo do Banco
```bash
# 1. Parar aplicação
# 2. Reverter todas migrations
go run cmd/migrate/main.go -direction=down

# 3. Aplicar novamente
go run cmd/migrate/main.go -direction=up
```

### Verificar Status Geral
```bash
# 1. Status Docker
docker-compose ps

# 2. Status Migrations
go run cmd/migrate/main.go -direction=up -steps=0

# 3. Testar Aplicação
curl http://localhost:8080/health
```

---

## 📝 **NOTAS IMPORTANTES**

- **SEMPRE** suba o Docker antes de rodar migrations
- **SEMPRE** teste migrations em desenvolvimento primeiro
- **NUNCA** faça `migrate-down` em produção sem backup
- **SEMPRE** use `go mod tidy` após mudanças no go.mod
- **SEMPRE** teste com `go test` antes de commit

---

## 🎯 **COMANDOS MAKEFILE (quando instalar make)**

```bash
make migrate-up          # Aplicar migrations
make migrate-down        # Reverter migration
make migrate-create name=add_phone  # Criar migration
make docker-up          # Subir Docker
make docker-down        # Parar Docker
make dev               # Rodar com hot-reload
make test              # Rodar testes
make build             # Compilar
make clean             # Limpar arquivos
make help              # Ver todos os comandos
```
