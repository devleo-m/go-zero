# 🐳 Docker - Guia de Troubleshooting e Setup

> Documentação completa de como configurar e resolver problemas do Docker no projeto GO ZERO

---

## 📋 Índice

1. [Problemas Identificados](#problemas-identificados)
2. [Soluções Aplicadas](#soluções-aplicadas)
3. [Comandos Executados](#comandos-executados)
4. [Referência Rápida](#referência-rápida)
5. [Comandos Úteis do Dia a Dia](#comandos-úteis-do-dia-a-dia)

---

## 🔍 Problemas Identificados

### Problema 1: Container não permanecia rodando
**Sintoma:** Container subia e parava imediatamente

**Causa:** O `Dockerfile` tinha o comando `CMD` comentado na linha 23:
```dockerfile
# CMD ["air", "-c", ".air.toml"]  ❌ COMENTADO
```

**Por quê isso é um problema?**
- Containers Docker precisam de um processo principal para permanecer rodando
- Sem o `CMD`, o container não sabe o que executar e termina imediatamente
- O Air é nossa ferramenta de hot-reload que mantém a aplicação rodando

---

### Problema 2: Aplicação não conectava ao banco de dados
**Sintoma:** Erros de conexão com PostgreSQL e Redis

**Causa:** Arquivo `.env` não existia, ou tinha configurações incorretas:
```env
DB_HOST=localhost      ❌ ERRADO (dentro do Docker)
REDIS_HOST=localhost   ❌ ERRADO (dentro do Docker)
```

**Por quê isso é um problema?**
- Dentro do Docker Compose, os serviços se comunicam por nome, não por `localhost`
- `localhost` dentro do container da aplicação aponta para o próprio container
- O banco está em outro container chamado `db`

**Correto:**
```env
DB_HOST=db           ✅ Nome do serviço no docker-compose.yml
REDIS_HOST=redis     ✅ Nome do serviço no docker-compose.yml
```

---

### Problema 3: Porta 8080 em conflito (RESOLVIDO)
**Sintoma:** Erro ao subir container: "bind: porta já em uso"

**Causa:** Outro processo estava usando a porta 8080

**Como identificar:** No nosso caso, a porta estava livre após verificação

---

### Problema 4: Falta de healthchecks
**Sintoma:** Aplicação tentava conectar ao banco antes dele estar pronto

**Causa:** `docker-compose.yml` não tinha healthchecks configurados

**Por quê isso é um problema?**
- PostgreSQL demora alguns segundos para inicializar
- Se a aplicação tentar conectar antes, vai falhar
- Healthchecks garantem ordem correta de inicialização

---

## 🛠️ Soluções Aplicadas

### ✅ Solução 1: Corrigir Dockerfile

**Arquivo:** `Dockerfile` (linha 23)

**ANTES:**
```dockerfile
# CMD ["air", "-c", ".air.toml"]
```

**DEPOIS:**
```dockerfile
CMD ["air", "-c", ".air.toml"]
```

**O que faz:**
- `CMD` define o comando padrão que o container executa ao iniciar
- `air` é a ferramenta de hot-reload para Go
- `-c .air.toml` especifica o arquivo de configuração do Air
- Com isso, o container inicia o Air, que compila e executa a aplicação Go

---

### ✅ Solução 2: Criar e configurar arquivo .env

**Comando executado:**
```powershell
# Copiar arquivo de exemplo para .env
Copy-Item env.example .env

# Ajustar DB_HOST para apontar ao container do banco
(Get-Content .env) -replace 'DB_HOST=localhost', 'DB_HOST=db' | Set-Content .env

# Ajustar REDIS_HOST para apontar ao container do Redis
(Get-Content .env) -replace 'REDIS_HOST=localhost', 'REDIS_HOST=redis' | Set-Content .env
```

**Por que cada comando:**

1. **`Copy-Item env.example .env`**
   - Cria o arquivo `.env` a partir do template
   - `.env` não é versionado no Git (está no .gitignore)
   - Cada desenvolvedor tem seu próprio `.env` com configurações locais

2. **`(Get-Content .env) -replace 'DB_HOST=localhost', 'DB_HOST=db'`**
   - Lê todo o conteúdo do arquivo `.env`
   - Substitui `DB_HOST=localhost` por `DB_HOST=db`
   - `db` é o nome do serviço PostgreSQL no `docker-compose.yml`
   - Dentro do Docker, containers se comunicam pelo nome do serviço

3. **`(Get-Content .env) -replace 'REDIS_HOST=localhost', 'REDIS_HOST=redis'`**
   - Mesmo conceito, mas para o Redis
   - `redis` é o nome do serviço no `docker-compose.yml`

---

### ✅ Solução 3: Melhorar docker-compose.yml com healthchecks

**Melhorias aplicadas:**

#### 3.1 - Healthcheck no PostgreSQL
```yaml
db:
  # ... outras configurações
  healthcheck:
    test: ["CMD-SHELL", "pg_isready -U postgres"]
    interval: 5s
    timeout: 5s
    retries: 5
```

**O que faz:**
- `test`: Executa `pg_isready` para verificar se o PostgreSQL está aceitando conexões
- `interval: 5s`: Verifica a cada 5 segundos
- `timeout: 5s`: Se não responder em 5s, considera falha
- `retries: 5`: Tenta 5 vezes antes de considerar unhealthy

#### 3.2 - Healthcheck na Aplicação
```yaml
app:
  # ... outras configurações
  healthcheck:
    test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
    interval: 10s
    timeout: 5s
    retries: 3
    start_period: 40s
```

**O que faz:**
- `test`: Faz requisição HTTP para `/health` da aplicação
- `start_period: 40s`: Aguarda 40s antes de começar os testes (tempo para compilar)
- Se `/health` retornar 200 OK, container é considerado healthy

#### 3.3 - Depends_on com condition
```yaml
app:
  depends_on:
    db:
      condition: service_healthy  # ← Espera DB ficar healthy
    redis:
      condition: service_started   # ← Espera Redis iniciar
```

**O que faz:**
- Garante que o banco está **pronto** (não apenas iniciado) antes de subir a aplicação
- Evita erros de conexão durante inicialização

---

### ✅ Solução 4: Rebuild e restart dos containers

**Comandos executados:**
```powershell
# 1. Parar todos os containers
docker-compose down

# 2. Subir containers reconstruindo as imagens
docker-compose up -d --build
```

**Por que cada comando:**

1. **`docker-compose down`**
   - Para todos os containers
   - Remove os containers (mas mantém os volumes com dados)
   - Necessário para aplicar mudanças no `docker-compose.yml`

2. **`docker-compose up -d --build`**
   - `-d`: "detached mode" - roda em background
   - `--build`: Reconstrói a imagem Docker (necessário após mudar o Dockerfile)
   - Sobe todos os serviços definidos no `docker-compose.yml`

---

## 📝 Comandos Executados (Ordem Cronológica)

### Passo 1: Verificar porta 8080
```powershell
netstat -ano | findstr :8080
```
**Por quê:** Verificar se algum processo estava usando a porta 8080  
**Resultado:** Porta estava livre ✅

---

### Passo 2: Corrigir Dockerfile
```powershell
# Arquivo editado manualmente
# Linha 23: Descomentado CMD ["air", "-c", ".air.toml"]
```
**Por quê:** Container precisa de um comando para executar  
**Resultado:** Container agora tem processo principal ✅

---

### Passo 3: Criar arquivo .env
```powershell
Copy-Item env.example .env
```
**Por quê:** Aplicação precisa de variáveis de ambiente  
**Resultado:** Arquivo `.env` criado ✅

---

### Passo 4: Ajustar DB_HOST
```powershell
(Get-Content .env) -replace 'DB_HOST=localhost', 'DB_HOST=db' | Set-Content .env
```
**Por quê:** Aplicação precisa se conectar ao container do banco  
**Resultado:** DB_HOST apontando corretamente ✅

---

### Passo 5: Ajustar REDIS_HOST
```powershell
(Get-Content .env) -replace 'REDIS_HOST=localhost', 'REDIS_HOST=redis' | Set-Content .env
```
**Por quê:** Aplicação precisa se conectar ao container do Redis  
**Resultado:** REDIS_HOST apontando corretamente ✅

---

### Passo 6: Parar containers existentes
```powershell
docker-compose down
```
**Por quê:** Aplicar mudanças no docker-compose.yml  
**Resultado:** Containers parados ✅

---

### Passo 7: Subir containers com rebuild
```powershell
docker-compose up -d --build
```
**Por quê:** Aplicar mudanças no Dockerfile e subir tudo  
**Resultado:** Todos os containers rodando ✅

---

### Passo 8: Verificar status
```powershell
docker-compose ps
```
**Por quê:** Confirmar que todos os containers estão rodando  
**Resultado:**
```
go-zero-app    Up (healthy)
go-zero-db     Up (healthy)
go-zero-redis  Up
```

---

### Passo 9: Testar aplicação
```powershell
curl http://localhost:8080/health
```
**Por quê:** Verificar se aplicação está respondendo  
**Resultado:** Status 200 OK com resposta JSON ✅

---

### Passo 10: Ver logs
```powershell
docker-compose logs --tail=100 app
```
**Por quê:** Confirmar que aplicação iniciou corretamente  
**Resultado:** Logs mostram servidor rodando na porta 8080 ✅

---

## 🚀 Referência Rápida

### Setup Inicial (Primeira vez ou após clonar projeto)

```powershell
# 1. Copiar arquivo de configuração
Copy-Item env.example .env

# 2. Ajustar configurações do Docker no .env
(Get-Content .env) -replace 'DB_HOST=localhost', 'DB_HOST=db' | Set-Content .env
(Get-Content .env) -replace 'REDIS_HOST=localhost', 'REDIS_HOST=redis' | Set-Content .env

# 3. Subir containers
docker-compose up -d --build

# 4. Aguardar 30-40 segundos e verificar status
docker-compose ps

# 5. Testar aplicação
curl http://localhost:8080/health
```

---

### Se container parar de funcionar

```powershell
# Ver logs para identificar problema
docker-compose logs app

# Ver logs em tempo real
docker-compose logs -f app

# Reiniciar apenas a aplicação
docker-compose restart app

# Reiniciar tudo
docker-compose restart
```

---

### Se precisar reconstruir tudo do zero

```powershell
# 1. Parar e remover tudo (inclusive volumes)
docker-compose down -v

# 2. Limpar cache do Docker (opcional)
docker system prune -f

# 3. Subir tudo novamente
docker-compose up -d --build
```

---

## 🔧 Comandos Úteis do Dia a Dia

### Ver logs
```powershell
# Logs de todos os serviços
docker-compose logs

# Logs apenas da aplicação
docker-compose logs app

# Logs em tempo real (seguir)
docker-compose logs -f app

# Ver últimas 50 linhas
docker-compose logs --tail=50 app
```

---

### Gerenciar containers
```powershell
# Ver status de todos os containers
docker-compose ps

# Parar todos os containers
docker-compose down

# Subir todos os containers
docker-compose up -d

# Reiniciar um serviço específico
docker-compose restart app

# Parar um serviço específico
docker-compose stop app

# Iniciar um serviço específico
docker-compose start app
```

---

### Acessar containers
```powershell
# Abrir shell no container da aplicação
docker-compose exec app sh

# Abrir shell no PostgreSQL
docker-compose exec db psql -U postgres -d go_zero_dev

# Executar comando no container
docker-compose exec app go version
```

---

### Verificar saúde
```powershell
# Verificar se aplicação está respondendo
curl http://localhost:8080/health

# Verificar métricas Prometheus
curl http://localhost:8080/metrics

# Ver status detalhado com healthcheck
docker inspect go-zero-app | grep -A 10 Health
```

---

### Limpar recursos
```powershell
# Remover containers parados
docker-compose down

# Remover containers e volumes (CUIDADO: apaga dados do banco)
docker-compose down -v

# Limpar imagens não utilizadas
docker image prune -f

# Limpar tudo (sistema completo)
docker system prune -af --volumes
```

---

### Rebuildar apenas a aplicação
```powershell
# Rebuildar e reiniciar apenas o container app
docker-compose up -d --build app

# Ver se rebuild funcionou
docker-compose logs app
```

---

### Verificar recursos
```powershell
# Ver uso de recursos (CPU, memória)
docker stats

# Ver uso de disco
docker system df

# Ver todos os containers (incluindo parados)
docker ps -a
```

---

## 🎯 Troubleshooting Comum

### Erro: "porta já em uso"
```powershell
# Descobrir o que está usando a porta
netstat -ano | findstr :8080

# Se encontrar um processo, encerrar pelo PID
taskkill /PID <número> /F

# Ou mudar a porta no .env
APP_PORT=8081
```

---

### Erro: "cannot connect to database"
```powershell
# Verificar se DB está healthy
docker-compose ps

# Ver logs do banco
docker-compose logs db

# Verificar se .env tem DB_HOST=db (não localhost)
cat .env | findstr DB_HOST

# Recriar containers
docker-compose down
docker-compose up -d
```

---

### Erro: "module not found" ou dependências
```powershell
# Rebuildar com --no-cache
docker-compose build --no-cache app
docker-compose up -d app

# Ou entrar no container e atualizar deps
docker-compose exec app go mod download
docker-compose exec app go mod tidy
```

---

### Container fica reiniciando infinitamente
```powershell
# Ver logs para identificar o erro
docker-compose logs app

# Verificar se CMD está descomentado no Dockerfile
cat Dockerfile | findstr CMD

# Verificar se arquivo .air.toml existe
ls .air.toml

# Testar localmente sem Docker
go run ./cmd/api/main.go
```

---

### Hot reload (Air) não está funcionando
```powershell
# Verificar se volume está montado corretamente
docker-compose config | findstr volumes

# Deve aparecer: - .:/app

# Reiniciar container
docker-compose restart app

# Ver logs do Air
docker-compose logs -f app
```

---

## 📚 Makefile - Atalhos Disponíveis

O projeto tem um `Makefile` com comandos prontos:

```bash
# Docker
make docker-up        # Sobe containers
make docker-down      # Para containers
make docker-logs      # Ver logs
make status           # Ver status dos containers

# Desenvolvimento
make dev              # Rodar localmente com Air
make run              # Rodar localmente sem Air
make build            # Compilar aplicação

# Banco de dados
make migrate-up       # Executar migrations
make migrate-down     # Reverter migrations
make db-reset         # Resetar banco completamente

# Úteis
make health           # Verificar saúde dos serviços
make dev-restart      # Reiniciar apenas aplicação
make clean-all        # Limpar tudo
```

---

## 🔑 Conceitos Importantes

### Docker Compose Network
- Containers no mesmo `docker-compose.yml` ficam na mesma rede
- Se comunicam pelo **nome do serviço** (não por IP ou localhost)
- Exemplo: serviço `db` é acessível em `db:5432`

### Healthchecks
- Verificam se serviço está realmente pronto (não apenas iniciado)
- `condition: service_healthy` garante ordem correta de inicialização
- Evita erros de "connection refused" durante startup

### Volumes
- `- .:/app` mapeia código local para dentro do container
- Mudanças no código local refletem imediatamente no container
- Permite hot-reload funcionar

### Build Context
- `context: .` define raiz do projeto como contexto
- Dockerfile pode acessar todos os arquivos do projeto
- `.dockerignore` define o que não copiar

---

## ✅ Checklist de Verificação

Antes de reportar problemas, verifique:

- [ ] Arquivo `.env` existe e está configurado corretamente
- [ ] `DB_HOST=db` (não localhost)
- [ ] `REDIS_HOST=redis` (não localhost)
- [ ] `CMD` está descomentado no Dockerfile
- [ ] Porta 8080 está livre
- [ ] Docker Desktop está rodando
- [ ] `docker-compose ps` mostra containers "Up" e "healthy"
- [ ] `curl http://localhost:8080/health` retorna 200 OK

---

## 📞 Serviços e Portas

Após `docker-compose up -d`, você tem acesso a:

| Serviço | URL | Descrição |
|---------|-----|-----------|
| API Go | http://localhost:8080 | Backend principal |
| PostgreSQL | localhost:5432 | Banco de dados |
| Redis | localhost:6379 | Cache |
| MongoDB | localhost:27017 | Chat database |
| MinIO | http://localhost:9000 | Object storage |
| MinIO Console | http://localhost:9001 | Interface do MinIO |
| Mailhog | http://localhost:8025 | Interface de emails |
| Prometheus | http://localhost:9090 | Métricas |
| Grafana | http://localhost:3000 | Dashboards |

---

## 🎓 Aprendizados

### Por que Air?
- Air detecta mudanças nos arquivos Go
- Recompila automaticamente
- Reinicia a aplicação
- Desenvolvimento muito mais rápido (não precisa parar e reiniciar manualmente)

### Por que Healthchecks?
- PostgreSQL demora ~5 segundos para aceitar conexões
- Sem healthcheck, app tenta conectar antes e falha
- Com healthcheck, app só inicia quando banco está pronto

### Por que .env separado?
- Cada desenvolvedor tem configurações diferentes
- Produção tem secrets diferentes
- `.env` não deve ir para o Git (está no .gitignore)

---

**Criado em:** 28 de Outubro de 2025  
**Última atualização:** 28 de Outubro de 2025  
**Versão:** 1.0.0

---

💡 **Dica:** Salve este arquivo e consulte sempre que tiver problemas com Docker!

