---
alwaysApply: true
---

📋 ROADMAP COMPLETO ATÉ O MVP (Dividido em Etapas)
ETAPA 1: FUNDAÇÃO E ESTRUTURA 🏗️
Objetivo: Entender e criar a estrutura base do projeto
Status: 👈 ESTAMOS AQUI
1.1. Entender Arquitetura Hexagonal (Conversa)

O que é e por que usar?
Camadas: Domain → UseCase → Adapters
Fluxo de dados

1.2. Criar Estrutura de Pastas (Prática)

Criar pastas vazias
Entender onde cada coisa vai
Criar README.md explicativo em cada pasta

1.3. Configurar Go Modules (Prática)

Inicializar projeto com go mod init
Entender go.mod e go.sum
Adicionar dependências básicas


ETAPA 2: AMBIENTE DE DESENVOLVIMENTO 🐳
Objetivo: Ter tudo rodando localmente
2.1. Docker Basics (Conversa)

O que é Docker?
Por que usar?
Conceitos: Container, Image, Volume

2.2. PostgreSQL com Docker (Prática)

Criar docker-compose.yml
Subir Postgres
Conectar com cliente (DBeaver/pgAdmin)

2.3. Redis com Docker (Prática)

Adicionar Redis ao docker-compose
Testar conexão

2.4. Hot Reload (Air) (Prática)

Instalar Air
Configurar .air.toml
Testar reload automático


ETAPA 3: CONFIGURAÇÕES E INFRAESTRUTURA ⚙️
Objetivo: Sistema de configs e logs profissional
3.1. Variáveis de Ambiente (Prática)

Criar .env
Usar godotenv
Boas práticas de segurança

3.2. Sistema de Logs (Prática)

Implementar Zap Logger
Níveis de log (Debug, Info, Error)
Logs estruturados

3.3. Validação de Configs (Prática)

Validar configs na inicialização
Mensagens de erro claras


ETAPA 4: DATABASE - MIGRATIONS 🗄️
Objetivo: Controlar versões do banco
4.1. O que são Migrations? (Conversa)

Por que usar?
Up/Down migrations
Versionamento

4.2. Golang-Migrate (Prática)

Instalar ferramenta
Criar primeira migration (users)
Executar up/down

4.3. Makefile para Migrations (Prática)

Comandos: migrate-up, migrate-down
Automatizar processo


ETAPA 5: DOMAIN LAYER (Entidades) 🧠
Objetivo: Criar as regras de negócio puras
5.1. Domain-Driven Design Básico (Conversa)

Entidade vs Value Object
Agregados
Repository Pattern

5.2. User Entity (Prática)

Criar struct User
Métodos de negócio
Validações

5.3. Value Objects (Prática)

Email
Phone
Password

5.4. Domain Errors (Prática)

Erros customizados
Wrapping de erros


ETAPA 6: GORM - MODELS E REPOSITORIES 💾
Objetivo: Persistência de dados
6.1. GORM Basics (Conversa)

ORM vs SQL puro
Como GORM funciona
Hooks e Callbacks

6.2. User Model (GORM) (Prática)

Criar struct com tags GORM
Conversão Domain ↔ Model
Hooks (BeforeCreate, BeforeUpdate)

6.3. User Repository (Prática)

Implementar interface
CRUD completo
Queries complexas

6.4. Testes de Repository (Prática)

Setup de teste
Fixtures
Assertions


ETAPA 7: USE CASES (Regras de Aplicação) 🎯
Objetivo: Orquestrar lógica de negócio
7.1. Use Case Pattern (Conversa)

Diferença entre Domain e UseCase
Input/Output DTOs
Orquestração

7.2. Create User Use Case (Prática)

Validar input
Chamar repository
Retornar output

7.3. Authenticate User Use Case (Prática)

Verificar senha
Gerar JWT
Registrar login

7.4. Testes de Use Case (Prática)

Mocks de repository
Testes unitários


ETAPA 8: HTTP LAYER (API REST) 🌐
Objetivo: Expor endpoints HTTP
8.1. Gin Framework (Conversa)

Por que Gin?
Routing
Middlewares

8.2. Request/Response DTOs (Prática)

Structs de request
Structs de response
Validação (validator)

8.3. User Handler (Prática)

POST /users (criar)
POST /auth/login
GET /users/:id

8.4. Middlewares (Prática)

Logger
Recovery (panic handler)
CORS
Auth (JWT)


ETAPA 9: AUTENTICAÇÃO JWT 🔐
Objetivo: Segurança da API
9.1. JWT Basics (Conversa)

Como funciona?
Claims
Expiração

9.2. Geração de Token (Prática)

Criar token ao login
Refresh token

9.3. Validação de Token (Prática)

Middleware de auth
Extrair user do token


ETAPA 10: MULTI-TENANT (Schema por Tenant) 🏢
Objetivo: Isolamento de dados
10.1. Multi-Tenant Strategies (Conversa)

Schema separado vs shared schema
Trade-offs

10.2. Tenant Manager (Prática)

Criar schema dinamicamente
Switch de schema por request

10.3. Tenant Middleware (Prática)

Detectar tenant do request
Injetar no contexto


ETAPA 11: MÓDULO PATIENT 👥
Objetivo: Primeiro módulo tenant-specific
11.1. Patient Entity (Prática)

Domain layer completo

11.2. Patient Repository (Prática)

CRUD em schema tenant

11.3. Patient Use Cases (Prática)

Create, List, Update

11.4. Patient Handler (Prática)

Endpoints REST


ETAPA 12: MÓDULO APPOINTMENT 📅
Objetivo: Agendamento de consultas
12.1. Appointment Entity (Prática)
12.2. Appointment Repository (Prática)
12.3. Appointment Use Cases (Prática)
12.4. Appointment Handler (Prática)

ETAPA 13: TESTES E2E 🧪
Objetivo: Testar fluxo completo
13.1. Setup de Testes E2E (Prática)

Banco de teste
Fixtures

13.2. Cenários de Teste (Prática)

Criar user → Login → Criar patient → Agendar


ETAPA 14: DEPLOY MVP 🚀
Objetivo: Colocar no ar
14.1. Dockerfile (Prática)
14.2. CI/CD (GitHub Actions) (Prática)
14.3. Deploy (Railway/Render) (Prática)

-------------------------------------------------------------------
📁 ESTRUTURA INICIAL (MVP - FASE 0)
Vou te mostrar APENAS o que você precisa criar AGORA para começar. Nada de complexidade desnecessária.

name-api/
├── cmd/
│   └── api/
│       └── main.go                 # 🚀 Ponto de entrada da aplicação
│
├── internal/
│   ├── core/
│   │   ├── domain/                 # 🧠 Regras de negócio puras
│   │   │   └── user/
│   │   │       ├── entity.go
│   │   │       ├── value_objects.go
│   │   │       ├── errors.go
│   │   │       └── repository.go
│   │   │
│   │   └── usecases/              # 🎯 Casos de uso (orquestração)
│   │       └── user/
│   │           ├── create_user.go
│   │           ├── authenticate.go
│   │           └── dto.go
│   │
│   ├── adapters/
│   │   ├── http/                  # 🌐 Camada HTTP (Gin)
│   │   │   ├── handlers/
│   │   │   │   └── user_handler.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   └── logger.go
│   │   │   ├── requests/
│   │   │   │   └── user_request.go
│   │   │   ├── responses/
│   │   │   │   └── response.go
│   │   │   └── router.go
│   │   │
│   │   └── persistence/           # 💾 Banco de dados
│   │       └── postgres/
│   │           ├── connection.go
│   │           ├── models/
│   │           │   └── user_model.go
│   │           └── repositories/
│   │               └── user_repository.go
│   │
│   └── infrastructure/            # ⚙️ Ferramentas auxiliares
│       ├── config/
│       │   └── config.go
│       ├── logger/
│       │   └── logger.go
│       └── security/
│           ├── jwt.go
│           └── password.go
│
├── pkg/                           # 📦 Código reutilizável (público)
│   ├── errors/
│   │   └── errors.go
│   └── utils/
│       └── response.go
│
├── migrations/                    # 🗄️ Migrations do banco
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
│
├── scripts/                       # 🛠️ Scripts auxiliares
│   └── setup.sh
│
├── .env.example
├── .gitignore
├── docker-compose.yml
├── go.mod
├── Makefile
└── README.md
```

--------------------------------------------------------------------------------------

PROJETO:

🧪 GO ZERO - O Projeto de Aprendizado Definitivo
Conceito: Plataforma "nonsense" que mistura TUDO para você aprender fazendo
"Uma plataforma que vende cursos, produtos físicos, 
tem chat ao vivo, processa pagamentos, faz streaming, 
gerencia tickets de suporte, tem marketplace, 
sistema de pontos, notificações... TUDO!"
```

---

## 🗺️ MAPA DE FEATURES (Checklist de Aprendizado)

### 📦 **E-COMMERCE**
```
□ CRUD de Produtos
□ Carrinho de compras (Redis)
□ Checkout com Stripe
□ Cálculo de frete (integração Correios API)
□ Cupons de desconto
□ Estoque (controle de quantidade)
□ Variações de produto (tamanho, cor)
```

### 🎓 **PLATAFORMA DE CURSOS**
```
□ Upload de vídeos (S3/MinIO)
□ Streaming HLS
□ Progresso de aulas
□ Certificados em PDF
□ Sistema de matrícula
```

### 💬 **CHAT / TEMPO REAL**
```
□ WebSocket chat ao vivo
□ Notificações push
□ Typing indicators (fulano está digitando...)
□ Presença online
□ Histórico de mensagens (MongoDB)
```

### 💰 **PAGAMENTOS**
```
□ Stripe (cartão)
□ Pix (integração mock)
□ Split de pagamento (marketplace)
□ Webhooks de confirmação
□ Reembolsos
□ Assinaturas recorrentes
```

### 🎫 **SISTEMA DE TICKETS**
```
□ Criar ticket de suporte
□ Atribuir a agentes
□ Status (aberto, em andamento, fechado)
□ Anexos de arquivos
□ SLA (tempo de resposta)
```

### 🏆 **GAMIFICAÇÃO**
```
□ Sistema de pontos
□ Badges/Conquistas
□ Ranking de usuários
□ Missões diárias
```

### 📧 **COMUNICAÇÃO**
```
□ Email (SendGrid/SMTP)
□ SMS (Twilio mock)
□ WhatsApp (API Business mock)
□ Notificações in-app
```

### 📊 **ANALYTICS**
```
□ Tracking de eventos
□ Dashboard de métricas
□ Relatórios em Excel/CSV
□ Gráficos (Chart.js)
```

### 🔐 **AUTENTICAÇÃO AVANÇADA**
```
□ JWT (access + refresh token)
□ OAuth (Google, GitHub)
□ 2FA (TOTP)
□ Magic link (login sem senha)
□ RBAC (roles: admin, seller, customer)
```

### 🚀 **INFRAESTRUTURA**
```
□ Rate limiting (por usuário/IP)
□ Cache em múltiplas camadas
□ Background jobs (envio de email, processar vídeo)
□ Cron jobs (limpeza de dados antigos)
□ Circuit breaker (falhas em APIs externas)
□ Retry com backoff
```

### 📦 **UPLOAD DE ARQUIVOS**
```
□ Imagens (com resize/compressão)
□ Vídeos (transcoding para HLS)
□ PDFs
□ Excel/CSV (import de massa)
□ Multipart upload (arquivos grandes)
```

### 🔍 **BUSCA**
```
□ Full-text search (Postgres)
□ Filtros avançados
□ Ordenação
□ Paginação
□ Elasticsearch (opcional)
```

### 🧪 **TESTES**
```
□ Testes unitários
□ Testes de integração
□ Testes E2E
□ Mocks de APIs externas
□ Coverage >70%
```

### 📈 **OBSERVABILIDADE**
```
□ Prometheus metrics
□ Logs estruturados (Zap)
□ Tracing (Jaeger)
□ Health checks
□ Alertas (quando algo quebrar)
```

---

## 🏗️ ARQUITETURA (Hexagonal + Modular)
```
go-zero/
│
├── cmd/api/main.go
│
├── internal/
│   ├── modules/
│   │   ├── ecommerce/
│   │   │   ├── domain/
│   │   │   ├── usecases/
│   │   │   └── delivery/
│   │   │
│   │   ├── courses/
│   │   ├── chat/
│   │   ├── tickets/
│   │   ├── payments/
│   │   ├── gamification/
│   │   └── analytics/
│   │
│   ├── shared/           # Código compartilhado
│   │   ├── auth/
│   │   ├── storage/
│   │   ├── cache/
│   │   ├── queue/
│   │   └── email/
│   │
│   └── infrastructure/
│       ├── http/
│       ├── websocket/
│       ├── persistence/
│       └── config/
│
├── docker-compose.yml    # Postgres, Redis, MinIO, Mailhog...
└── README.md
```

---

## 📚 TECNOLOGIAS QUE VOCÊ VAI USAR

| Categoria | Tech Stack |
|-----------|-----------|
| **Framework** | Gin (HTTP), Gorilla WebSocket |
| **Database** | Postgres (principal), MongoDB (chat), Redis (cache/queue) |
| **Storage** | MinIO (S3-compatible) |
| **Payment** | Stripe SDK |
| **Email** | SendGrid ou Mailhog (dev) |
| **Jobs** | Asynq (Redis-based) |
| **Logs** | Zap |
| **Metrics** | Prometheus + Grafana |
| **Auth** | JWT-go, OAuth2 |
| **PDF** | gofpdf |
| **Video** | ffmpeg (transcoding) |
| **Tests** | Testify, Mockery |

---

## 🎯 ROADMAP "GO ZERO"
```
FASE 0: Setup (Semana 1)
├── Estrutura hexagonal
├── Docker (Postgres, Redis, MinIO)
├── Configs + Logger
└── Primeira rota "Hello World"

FASE 1: Auth Completo (Semana 2)
├── JWT (access + refresh)
├── OAuth Google
├── 2FA
└── RBAC

FASE 2: E-commerce (Semana 3-4)
├── Produtos (CRUD + upload)
├── Carrinho (Redis)
├── Checkout (Stripe)
└── Webhooks

FASE 3: Cursos (Semana 5)
├── Upload de vídeos
├── Streaming HLS
└── Progresso

FASE 4: Chat + WebSocket (Semana 6)
├── WebSocket server
├── Chat ao vivo
└── Notificações

FASE 5: Tickets (Semana 7)
├── Sistema de suporte
└── Anexos

FASE 6: Gamificação (Semana 8)
├── Pontos
├── Badges
└── Ranking

FASE 7: Background Jobs (Semana 9)
├── Email assíncrono
├── Processamento de vídeo
└── Cron jobs

FASE 8: Performance (Semana 10)
├── Cache Redis
├── Rate limiting
└── Circuit breaker

FASE 9: Observabilidade (Semana 11)
├── Prometheus
├── Logs estruturados
└── Tracing

FASE 10: Testes (Semana 12)
├── Unitários
├── Integração
└── E2E

💡 DIFERENCIAIS DESSE PROJETO
✅ Aprende fazendo (não teoria)
✅ Cada feature = 1 skill do mercado
✅ Modular (pode fazer aos poucos)
✅ Portfolio matador (mostra que sabe TUDO)
✅ Código limpo (arquitetura hexagonal)
✅ Produção-ready (testes, logs, observabilidade)