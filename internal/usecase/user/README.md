# User Use Cases Module

## 📋 Visão Geral

Este módulo contém todos os casos de uso (Use Cases) relacionados ao domínio de usuários, implementando a camada de aplicação da arquitetura hexagonal. Cada caso de uso representa uma operação de negócio específica e orquestra as interações entre o domínio e a infraestrutura.

## 🏗️ Arquitetura

```
internal/usecase/user/
├── dtos.go                      # DTOs de entrada e saída
├── interfaces.go                # Interfaces dos serviços externos
├── create_user.go               # Caso de uso: Criar usuário
├── authenticate_user.go         # Caso de uso: Autenticar usuário
├── get_user.go                  # Caso de uso: Buscar usuário
├── list_users.go                # Caso de uso: Listar usuários
├── update_user.go               # Caso de uso: Atualizar usuário
├── change_password.go           # Caso de uso: Alterar senha
├── activate_user.go             # Caso de uso: Ativar usuário
├── deactivate_user.go           # Caso de uso: Desativar usuário
├── suspend_user.go              # Caso de uso: Suspender usuário
├── change_role.go               # Caso de uso: Alterar role
├── user_usecase_aggregate.go    # Agregado principal
└── README.md                    # Esta documentação
```

## 🎯 Casos de Uso Implementados

### 1. **CreateUserUseCase**
- **Propósito**: Criar um novo usuário no sistema
- **Validações**: Email único, senha forte, dados obrigatórios
- **Ações**: Cria entidade, salva no banco, envia email de boas-vindas
- **Input**: `CreateUserInput`
- **Output**: `CreateUserOutput`

### 2. **AuthenticateUserUseCase**
- **Propósito**: Autenticar usuário com email e senha
- **Validações**: Credenciais válidas, usuário ativo
- **Ações**: Verifica senha, gera tokens JWT, atualiza login
- **Input**: `AuthenticateUserInput`
- **Output**: `AuthenticateUserOutput`

### 3. **GetUserUseCase**
- **Propósito**: Buscar usuário por ID
- **Validações**: ID válido, usuário existe
- **Ações**: Busca no repositório, retorna dados seguros
- **Input**: `GetUserInput`
- **Output**: `GetUserOutput`

### 4. **ListUsersUseCase**
- **Propósito**: Listar usuários com paginação e filtros
- **Validações**: Parâmetros de paginação válidos
- **Ações**: Aplica filtros, busca paginada, retorna metadados
- **Input**: `ListUsersInput`
- **Output**: `ListUsersOutput`

### 5. **UpdateUserUseCase**
- **Propósito**: Atualizar perfil do usuário
- **Validações**: Dados válidos, telefone único
- **Ações**: Atualiza entidade, salva no banco
- **Input**: `UpdateUserInput`
- **Output**: `UpdateUserOutput`

### 6. **ChangePasswordUseCase**
- **Propósito**: Alterar senha do usuário
- **Validações**: Senha atual correta, nova senha forte
- **Ações**: Verifica senha atual, define nova senha
- **Input**: `ChangePasswordInput`
- **Output**: `ChangePasswordOutput`

### 7. **ActivateUserUseCase**
- **Propósito**: Ativar usuário pendente
- **Validações**: Usuário está pendente
- **Ações**: Ativa usuário, envia confirmação
- **Input**: `ActivateUserInput`
- **Output**: `ActivateUserOutput`

### 8. **DeactivateUserUseCase**
- **Propósito**: Desativar usuário
- **Validações**: Usuário não está inativo
- **Ações**: Desativa usuário, envia notificação
- **Input**: `DeactivateUserInput`
- **Output**: `DeactivateUserOutput`

### 9. **SuspendUserUseCase**
- **Propósito**: Suspender usuário
- **Validações**: Usuário não está suspenso
- **Ações**: Suspende usuário, envia notificação
- **Input**: `SuspendUserInput`
- **Output**: `SuspendUserOutput`

### 10. **ChangeRoleUseCase**
- **Propósito**: Alterar role do usuário
- **Validações**: Permissões do solicitante, role válido
- **Ações**: Altera role, envia notificação
- **Input**: `ChangeRoleInput`
- **Output**: `ChangeRoleOutput`

## 🔧 Como Usar

### Instanciação

```go
// Criar dependências
userRepo := postgres.NewUserRepository(db, logger)
emailService := email.NewService(config)
jwtService := auth.NewJWTService(config)

// Criar agregado de casos de uso
userUC := user.NewUserUseCaseAggregate(
    userRepo,
    emailService,
    jwtService,
    logger,
)
```

### Exemplos de Uso

#### Criar Usuário
```go
input := user.CreateUserInput{
    Name:     "João Silva",
    Email:    "joao@example.com",
    Password: "MinhaSenh@123",
    Phone:    "+5511999999999",
    Role:     "user",
}

result, err := userUC.CreateUser(ctx, input)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Usuário criado: %+v\n", result.User)
```

#### Autenticar Usuário
```go
input := user.AuthenticateUserInput{
    Email:    "joao@example.com",
    Password: "MinhaSenh@123",
}

result, err := userUC.AuthenticateUser(ctx, input)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Token: %s\n", result.AccessToken)
```

#### Listar Usuários
```go
input := user.ListUsersInput{
    Page:     1,
    PageSize: 20,
    Role:     "user",
    Status:   "active",
    Search:   "joão",
}

result, err := userUC.ListUsers(ctx, input)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Usuários encontrados: %d\n", len(result.Users))
```

## 📊 DTOs (Data Transfer Objects)

### Input DTOs
- `CreateUserInput`: Dados para criar usuário
- `AuthenticateUserInput`: Credenciais de login
- `GetUserInput`: ID do usuário
- `ListUsersInput`: Parâmetros de listagem
- `UpdateUserInput`: Dados para atualizar
- `ChangePasswordInput`: Dados para alterar senha
- `ActivateUserInput`: ID do usuário para ativar
- `DeactivateUserInput`: ID do usuário para desativar
- `SuspendUserInput`: ID do usuário para suspender
- `ChangeRoleInput`: Dados para alterar role

### Output DTOs
- `UserOutput`: Dados seguros do usuário
- `CreateUserOutput`: Resultado da criação
- `AuthenticateUserOutput`: Tokens de autenticação
- `ListUsersOutput`: Lista paginada de usuários
- `PaginationOutput`: Metadados de paginação
- E outros outputs específicos...

## 🔌 Interfaces de Serviços

### Logger
```go
type Logger interface {
    Debug(msg string, fields ...interface{})
    Info(msg string, fields ...interface{})
    Warn(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
}
```

### EmailService
```go
type EmailService interface {
    SendWelcomeEmail(ctx context.Context, email, name string) error
    SendActivationConfirmationEmail(ctx context.Context, email, name string) error
    // ... outros métodos
}
```

### JWTService
```go
type JWTService interface {
    GenerateAccessToken(userID string, email string, role string) (string, int, error)
    GenerateRefreshToken(userID string) (string, error)
    ValidateToken(token string) (map[string]interface{}, error)
    // ... outros métodos
}
```

## 🛡️ Validações e Segurança

### Validações de Input
- **Email**: Formato válido, único no sistema
- **Senha**: Mínimo 8 caracteres, maiúscula, minúscula, número, especial
- **Nome**: 2-100 caracteres, sem caracteres perigosos
- **Telefone**: Formato válido, único no sistema
- **Role**: Valores permitidos (admin, manager, user, guest)
- **Status**: Valores permitidos (active, inactive, pending, suspended)

### Regras de Negócio
- Usuários começam como "pending" e precisam ser ativados
- Apenas usuários ativos podem fazer login
- Admins podem gerenciar todos os usuários
- Managers podem gerenciar users e guests
- Users podem gerenciar apenas guests
- Não é possível alterar o próprio role
- Senhas não podem ser muito similares à atual

### Logs de Segurança
- Tentativas de login falhadas
- Alterações de role
- Ativações/desativações
- Alterações de senha
- Acessos a dados sensíveis

## 🧪 Testes

### Estrutura de Testes
```
tests/
├── unit/
│   └── usecase/
│       └── user/
│           ├── create_user_test.go
│           ├── authenticate_user_test.go
│           └── ...
├── integration/
│   └── usecase/
│       └── user/
│           └── user_usecase_test.go
└── e2e/
    └── user_flows_test.go
```

### Exemplo de Teste Unitário
```go
func TestCreateUserUseCase_Execute(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    mockEmailService := &MockEmailService{}
    mockLogger := &MockLogger{}
    
    useCase := NewCreateUserUseCase(mockRepo, mockEmailService, mockLogger)
    
    input := CreateUserInput{
        Name:     "João Silva",
        Email:    "joao@example.com",
        Password: "MinhaSenh@123",
        Role:     "user",
    }
    
    // Act
    result, err := useCase.Execute(context.Background(), input)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "João Silva", result.User.Name)
}
```

## 📈 Métricas e Monitoramento

### Métricas Coletadas
- Total de usuários
- Usuários por status
- Usuários por role
- Tentativas de login
- Operações de CRUD
- Tempo de resposta dos casos de uso

### Logs Estruturados
- Operações de negócio
- Eventos de segurança
- Erros e exceções
- Performance e latência

## 🔄 Fluxo de Dados

```
HTTP Request
    ↓
Handler (Gin)
    ↓
Use Case (Orquestração)
    ↓
Domain Entity (Regras de Negócio)
    ↓
Repository (Persistência)
    ↓
Database (PostgreSQL)
    ↓
Response DTO
    ↓
JSON Response
```

## 🚀 Próximos Passos

1. **Implementar testes unitários** para todos os casos de uso
2. **Adicionar cache** para consultas frequentes
3. **Implementar rate limiting** para operações sensíveis
4. **Adicionar auditoria** completa de operações
5. **Implementar notificações** push e in-app
6. **Adicionar métricas** de Prometheus
7. **Implementar circuit breakers** para serviços externos

## 📚 Referências

- [Arquitetura Hexagonal](https://alistair.cockburn.us/hexagonal-architecture/)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Best Practices](https://golang.org/doc/effective_go.html)
