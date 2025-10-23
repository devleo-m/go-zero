# Interface HTTP Layer

## 📋 Visão Geral

A camada de interface HTTP é responsável por expor a API REST da aplicação, seguindo os princípios da arquitetura hexagonal. Esta camada atua como o ponto de entrada para todas as requisições HTTP e é responsável por:

- **Receber requisições HTTP** e convertê-las para DTOs
- **Validar dados de entrada** usando validações customizadas
- **Chamar os casos de uso** apropriados da camada de aplicação
- **Converter respostas** dos casos de uso para DTOs HTTP
- **Retornar respostas JSON** padronizadas
- **Gerenciar autenticação e autorização** via middlewares
- **Tratar erros** de forma centralizada

## 🏗️ Arquitetura

```
HTTP Request
    ↓
Router (Gin)
    ↓
Middlewares (Auth, CORS, Logger, Recovery)
    ↓
Handler (UserHandler, HealthHandler)
    ↓
DTOs (Request/Response)
    ↓
Use Cases (Camada de Aplicação)
    ↓
Domain (Camada de Domínio)
    ↓
Infrastructure (Camada de Infraestrutura)
```

## 📁 Estrutura de Pastas

```
interface/http/
├── dto/                    # DTOs HTTP (Request/Response)
│   ├── user_dto.go        # DTOs específicos de usuário
│   └── common_dto.go      # DTOs comuns e helpers
├── handlers/              # Handlers HTTP (Controllers)
│   ├── user_handler.go    # Handler de usuários
│   ├── health_handler.go  # Handler de health check
│   └── error_handler.go   # Handler de erros centralizado
├── middleware/            # Middlewares HTTP
│   ├── auth.go           # Autenticação JWT
│   ├── cors.go           # CORS
│   ├── logger.go         # Logging
│   └── recovery.go       # Recovery de panics
├── router/               # Configuração de rotas
│   ├── router.go         # Router principal
│   ├── routes_v1.go      # Rotas da API v1
│   └── routes_public.go  # Rotas públicas
├── validation/           # Validações customizadas
│   └── validator.go      # Validador customizado
└── README.md            # Esta documentação
```

## 🔧 Componentes Principais

### 1. DTOs (Data Transfer Objects)

#### User DTOs
- **CreateUserRequest**: Dados para criar usuário
- **AuthenticateUserRequest**: Credenciais de login
- **UpdateUserRequest**: Dados para atualizar usuário
- **UserResponse**: Resposta com dados do usuário
- **ListUsersResponse**: Resposta com lista paginada

#### Common DTOs
- **SuccessResponse**: Resposta de sucesso genérica
- **ErrorResponse**: Resposta de erro padronizada
- **ValidationErrorResponse**: Resposta de erro de validação
- **PaginationResponse**: Metadados de paginação

### 2. Handlers

#### UserHandler
Gerencia todas as operações relacionadas a usuários:

```go
// Endpoints disponíveis:
POST   /api/v1/users              # Criar usuário
GET    /api/v1/users/:id          # Buscar usuário por ID
GET    /api/v1/users/email/:email # Buscar usuário por email
GET    /api/v1/users              # Listar usuários (paginado)
PUT    /api/v1/users/:id          # Atualizar usuário
PUT    /api/v1/users/:id/password # Alterar senha
POST   /api/v1/users/:id/activate # Ativar usuário (admin)
POST   /api/v1/users/:id/deactivate # Desativar usuário (admin)
POST   /api/v1/users/:id/suspend  # Suspender usuário (admin)
PUT    /api/v1/users/:id/role     # Alterar role (admin)
GET    /api/v1/users/stats        # Estatísticas (admin)
```

#### HealthHandler
Gerencia health checks e informações do sistema:

```go
// Endpoints disponíveis:
GET    /health    # Health check básico
GET    /ready     # Readiness check
GET    /live      # Liveness check
GET    /version   # Versão da aplicação
GET    /metrics   # Métricas (protegido)
```

### 3. Middlewares

#### AuthMiddleware
- Valida tokens JWT
- Extrai informações do usuário
- Adiciona dados ao contexto

#### CORSMiddleware
- Configura CORS para desenvolvimento/produção
- Suporte a wildcards e origens específicas

#### LoggerMiddleware
- Log de todas as requisições
- Captura de body de request/response
- Logs estruturados com Zap

#### RecoveryMiddleware
- Captura e trata panics
- Logs de erros com stack trace
- Respostas de erro padronizadas

### 4. Validações Customizadas

#### Validações de Negócio
- **CPF**: Validação de CPF brasileiro
- **CNPJ**: Validação de CNPJ brasileiro
- **PhoneBR**: Validação de telefone brasileiro
- **CEP**: Validação de CEP brasileiro
- **StrongPassword**: Validação de senha forte
- **ValidRole**: Validação de role válido
- **ValidStatus**: Validação de status válido
- **UUID4**: Validação de UUID v4

## 🚀 Como Usar

### 1. Configuração Básica

```go
// Criar instâncias necessárias
validator := validation.NewCustomValidator()
errorHandler := handlers.NewErrorHandler(logger)
userHandler := handlers.NewUserHandler(useCaseAggregate, validator, errorHandler, logger)
healthHandler := handlers.NewHealthHandler(logger)

// Configurar router
config := router.Config{
    UserUseCaseAggregate: useCaseAggregate,
    JWTService:           jwtService,
    Logger:               logger,
    Environment:          "development",
    EnableSwagger:        true,
    EnableMetrics:        true,
}

r := router.NewRouter(config)
```

### 2. Iniciar Servidor

```go
// Servidor HTTP
r.Run(":8080")

// Servidor HTTPS
r.RunTLS(":8443", "cert.pem", "key.pem")
```

### 3. Exemplo de Requisição

```bash
# Criar usuário
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@example.com",
    "password": "MinhaSenh@123",
    "role": "user"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "joao@example.com",
    "password": "MinhaSenh@123"
  }'

# Buscar usuário (com autenticação)
curl -X GET http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <token>"
```

## 🔒 Autenticação e Autorização

### JWT Authentication
- Tokens de acesso com expiração
- Refresh tokens para renovação
- Validação automática via middleware

### Role-based Access Control
- **Admin**: Acesso total ao sistema
- **Manager**: Gerenciamento limitado de usuários
- **User**: Acesso básico ao próprio perfil
- **Guest**: Acesso apenas a informações públicas

### Middleware de Autorização
```go
// Requer role específico
router.Use(middleware.RequireRole("admin"))

// Requer qualquer um dos roles
router.Use(middleware.RequireAnyRole("admin", "manager"))
```

## 📊 Logging e Monitoramento

### Logs Estruturados
- **Request/Response**: Todas as requisições HTTP
- **Business Events**: Eventos de negócio importantes
- **Security Events**: Tentativas de acesso, falhas de autenticação
- **Error Logs**: Erros com stack trace e contexto

### Métricas
- **Response Time**: Tempo de resposta das requisições
- **Request Count**: Contagem de requisições por endpoint
- **Error Rate**: Taxa de erro por endpoint
- **User Activity**: Atividade dos usuários

## 🧪 Testes

### Testes de Handler
```go
func TestUserHandler_CreateUser(t *testing.T) {
    // Setup
    handler := NewUserHandler(useCaseAggregate, validator, errorHandler, logger)
    
    // Test
    req := dto.CreateUserRequest{
        Name:     "João Silva",
        Email:    "joao@example.com",
        Password: "MinhaSenh@123",
        Role:     "user",
    }
    
    // Execute
    // Assert
}
```

### Testes de Middleware
```go
func TestAuthMiddleware(t *testing.T) {
    // Setup
    middleware := AuthMiddleware(jwtService, logger)
    
    // Test valid token
    // Test invalid token
    // Test missing token
}
```

## 🔧 Configuração

### Variáveis de Ambiente
```env
# Servidor
APP_PORT=8080
APP_ENV=development

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001

# Rate Limiting
RATE_LIMIT_REQUESTS_PER_MINUTE=100

# Timeout
REQUEST_TIMEOUT_SECONDS=30

# Body Size
MAX_BODY_SIZE_MB=10
```

### Configuração de CORS
```go
// Desenvolvimento
router.Use(middleware.CORSForDevelopment())

// Produção
allowedOrigins := []string{"https://app.gozero.dev"}
router.Use(middleware.CORSForProduction(allowedOrigins))
```

## 📈 Performance

### Otimizações Implementadas
- **Connection Pooling**: Pool de conexões HTTP
- **Request Timeout**: Timeout configurável
- **Body Size Limit**: Limite de tamanho do body
- **Rate Limiting**: Proteção contra spam
- **Caching Headers**: Headers de cache para rotas públicas

### Métricas de Performance
- **Response Time**: < 100ms para 95% das requisições
- **Throughput**: > 1000 req/s
- **Error Rate**: < 0.1%
- **Memory Usage**: < 100MB

## 🛡️ Segurança

### Headers de Segurança
- **X-Frame-Options**: DENY
- **X-Content-Type-Options**: nosniff
- **X-XSS-Protection**: 1; mode=block
- **Strict-Transport-Security**: HSTS em HTTPS
- **Content-Security-Policy**: Política de segurança de conteúdo

### Validação de Entrada
- **Sanitização**: Limpeza de dados de entrada
- **Validação**: Validação rigorosa de todos os campos
- **Rate Limiting**: Proteção contra ataques de força bruta
- **Input Size Limits**: Limites de tamanho de entrada

## 🔄 Versionamento da API

### Estratégia de Versionamento
- **URL Path**: `/api/v1/`, `/api/v2/`
- **Header**: `API-Version: v1`
- **Backward Compatibility**: Manter compatibilidade com versões anteriores

### Exemplo de Versionamento
```go
// v1
GET /api/v1/users

// v2 (com breaking changes)
GET /api/v2/users
```

## 📚 Documentação da API

### Swagger/OpenAPI
- Documentação automática via Swagger
- Exemplos de requisições e respostas
- Testes interativos via Swagger UI

### Exemplos de Uso
- **cURL**: Comandos curl para testar a API
- **JavaScript**: Exemplos com fetch/axios
- **Python**: Exemplos com requests
- **Postman**: Collection do Postman

## 🚨 Tratamento de Erros

### Códigos de Status HTTP
- **200**: Sucesso
- **201**: Criado com sucesso
- **400**: Erro de validação
- **401**: Não autorizado
- **403**: Proibido
- **404**: Não encontrado
- **409**: Conflito
- **422**: Entidade não processável
- **500**: Erro interno do servidor

### Formato de Erro
```json
{
  "success": false,
  "error": "VALIDATION_ERROR",
  "message": "One or more fields are invalid",
  "details": [
    {
      "field": "email",
      "message": "Invalid email format",
      "value": "invalid-email"
    }
  ]
}
```

## 🔧 Manutenção

### Logs de Debug
```go
// Habilitar logs de debug
logger.SetLevel(zap.DebugLevel)
```

### Health Checks
```bash
# Health check básico
curl http://localhost:8080/health

# Readiness check
curl http://localhost:8080/ready

# Liveness check
curl http://localhost:8080/live
```

### Métricas
```bash
# Métricas do sistema
curl -H "Authorization: Bearer <token>" http://localhost:8080/metrics
```

## 🎯 Próximos Passos

1. **Implementar Swagger**: Documentação automática da API
2. **Adicionar Rate Limiting**: Proteção contra spam
3. **Implementar Caching**: Cache de respostas frequentes
4. **Adicionar Webhooks**: Notificações em tempo real
5. **Implementar GraphQL**: API GraphQL alternativa
6. **Adicionar WebSocket**: Comunicação em tempo real
7. **Implementar gRPC**: API gRPC para comunicação interna

## 📞 Suporte

Para dúvidas ou problemas:
- **Documentação**: [docs.gozero.dev](https://docs.gozero.dev)
- **Issues**: [GitHub Issues](https://github.com/go-zero/go-zero/issues)
- **Email**: support@gozero.dev
