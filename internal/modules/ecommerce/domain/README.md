# 🧠 DOMAIN LAYER - ECOMMERCE

## O que é este módulo?

O **Domain Layer** é o **coração** da aplicação GO ZERO. Aqui ficam todas as **regras de negócio puras** - sem dependências de frameworks, banco de dados ou HTTP.

## 🏗️ Estrutura

```
domain/
├── entities/           # Entidades de negócio
│   ├── user.go        # Entidade User
│   └── user_test.go   # Testes da entidade User
├── valueobjects/       # Objetos de valor
│   ├── email.go       # Email com validação
│   ├── password.go    # Senha com hash seguro
│   └── money.go       # Valor monetário com precisão
├── services/          # Serviços de domínio
│   └── password_service.go  # Serviços de senha
├── repositories/      # Interfaces de repositórios
│   └── user_repository.go   # Contrato para acesso a dados
├── errors/           # Erros específicos do domínio
│   └── domain_errors.go     # Erros customizados
├── domain.go         # Ponto de entrada do módulo
└── README.md         # Este arquivo
```

## 🎯 Entidades

### User (Usuário)
**O que é:** Representa um usuário do sistema com todas as regras de negócio.

**Características:**
- ✅ Validação de email e senha
- ✅ Controle de status (ativo, inativo, bloqueado)
- ✅ Verificação de email
- ✅ Controle de papéis (admin, client)
- ✅ Soft delete
- ✅ Rastreamento de login

**Exemplo de uso:**
```go
// Import direto (mais claro)
import "github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/entities"
import "github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/valueobjects"

// Criar usuário
user, err := entities.NewUser("user@example.com", "Password123!", "João Silva", entities.UserRoleClient)
if err != nil {
    log.Fatal(err)
}

// Verificar senha
if user.VerifyPassword("Password123!") {
    fmt.Println("Senha correta!")
}

// Ativar usuário
user.Activate()
user.VerifyEmail()

// Verificar se pode fazer login
if user.CanLogin() {
    fmt.Println("Usuário pode fazer login!")
}
```

## 💎 Value Objects

### Email
**O que é:** Email válido com validação rigorosa.

**Características:**
- ✅ Validação de formato
- ✅ Normalização (lowercase, trim)
- ✅ Extração de domínio
- ✅ Comparação segura

**Exemplo:**
```go
import "github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/valueobjects"

email, err := valueobjects.NewEmail("USER@EXAMPLE.COM")
if err != nil {
    log.Fatal(err)
}

fmt.Println(email.String())        // "user@example.com"
fmt.Println(email.GetDomain())     // "example.com"
fmt.Println(email.GetLocalPart())  // "user"
```

### Password
**O que é:** Senha segura com hash Argon2.

**Características:**
- ✅ Validação de força
- ✅ Hash seguro (Argon2)
- ✅ Verificação de senha
- ✅ Resistente a ataques

**Exemplo:**
```go
import "github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/valueobjects"

password, err := valueobjects.NewPassword("MySecure123!")
if err != nil {
    log.Fatal(err)
}

// Verificar senha
if password.Verify("MySecure123!") {
    fmt.Println("Senha correta!")
}
```

### Money
**O que é:** Valor monetário com precisão decimal.

**Características:**
- ✅ Precisão decimal (evita problemas de float)
- ✅ Suporte a múltiplas moedas
- ✅ Operações matemáticas seguras
- ✅ Formatação para exibição

**Exemplo:**
```go
import "github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/valueobjects"

price, err := valueobjects.NewMoney(99.99, "BRL")
if err != nil {
    log.Fatal(err)
}

fmt.Println(price.String())     // "99.99 BRL"
fmt.Println(price.Format())     // "99,99 BRL"

// Operações (agora mais simples!)
total := price.Multiply(2)      // R$ 199,98
sum := price.Add(price)         // R$ 199,98
```

## 🔧 Serviços de Domínio

### PasswordService
**O que é:** Serviços relacionados a senhas.

**Funcionalidades:**
- ✅ Geração de senhas aleatórias
- ✅ Validação de força
- ✅ Geração de hash
- ✅ Verificação de senha

**Exemplo:**
```go
import "github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/services"

service := services.NewPasswordService()

// Gerar senha aleatória
password, err := service.GenerateRandomPassword(12)
if err != nil {
    log.Fatal(err)
}

// Validar força (agora usa o value object internamente!)
isValid, errors := service.ValidatePasswordStrength("MyPass123!")
if !isValid {
    fmt.Println("Erros:", errors)
}
```

## 📋 Repositórios

### UserRepository
**O que é:** Interface que define como acessar dados de usuários.

**Operações:**
- ✅ CRUD básico
- ✅ Busca por email
- ✅ Listagem com paginação
- ✅ Busca por critérios
- ✅ Estatísticas

**Exemplo:**
```go
import "github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/repositories"

type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    GetByID(ctx context.Context, id string) (*entities.User, error)
    GetByEmail(ctx context.Context, email valueobjects.Email) (*entities.User, error)
    Update(ctx context.Context, user *entities.User) error
    Delete(ctx context.Context, id string) error
    // ... mais operações
}
```

## ❌ Erros do Domínio

### DomainError
**O que é:** Erros específicos do domínio com códigos e contexto.

**Características:**
- ✅ Códigos padronizados
- ✅ Mensagens claras
- ✅ Contexto adicional
- ✅ Fácil identificação

**Exemplo:**
```go
import "github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/errors"

// Erro específico
err := errors.ErrUserNotFound

// Erro com contexto
err := errors.NewDomainError("INVALID_EMAIL", "email inválido", map[string]interface{}{
    "email": "invalid-email",
    "field": "email",
})
```

## 🧪 Testes

### Como executar
```bash
# Todos os testes
go test ./internal/modules/ecommerce/domain/...

# Com coverage
go test -cover ./internal/modules/ecommerce/domain/...

# Verbose
go test -v ./internal/modules/ecommerce/domain/...
```

### Cobertura de testes
- ✅ **User Entity:** 100% dos métodos testados
- ✅ **Value Objects:** Validação e operações
- ✅ **Domain Services:** Funcionalidades principais
- ✅ **Error Handling:** Cenários de erro

## 🎯 Princípios Seguidos

### 1. **Pure Business Logic**
- ❌ Sem dependências externas
- ❌ Sem frameworks
- ❌ Sem banco de dados
- ✅ Só regras de negócio

### 2. **Rich Domain Model**
- ✅ Entidades com comportamento
- ✅ Validações no domínio
- ✅ Regras de negócio encapsuladas

### 3. **Testabilidade**
- ✅ Fácil de testar
- ✅ Sem dependências externas
- ✅ Testes isolados

### 4. **Reutilização**
- ✅ Usado em qualquer camada
- ✅ Independente de framework
- ✅ Portável

## 🚀 Como Usar

### Import
```go
import "github.com/devleo-m/go-zero/internal/modules/ecommerce/domain"
```

### Exemplo Completo
```go
import (
    "github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/entities"
    "github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/valueobjects"
)

// Criar usuário
user, err := entities.NewUser("user@example.com", "Password123!", "João Silva", entities.UserRoleClient)
if err != nil {
    log.Fatal(err)
}

// Verificar se pode fazer login
if user.CanLogin() {
    fmt.Println("Usuário pode fazer login!")
}

// Alterar senha
err = user.ChangePassword("NewPassword123!")
if err != nil {
    log.Fatal(err)
}

// Ativar usuário
user.Activate()
user.VerifyEmail()

// Exemplo com Money (operações simplificadas!)
price, _ := valueobjects.NewMoney(99.99, "BRL")
total := price.Multiply(2)  // Sem erro, mais simples!
```

## 🔄 Próximos Passos

1. **Use Cases Layer** - Orquestrar lógica de negócio
2. **Adapters Layer** - Implementar repositórios
3. **HTTP Layer** - Expor via API REST

## 📚 Conceitos Importantes

- **Entity:** Objeto com identidade única
- **Value Object:** Objeto identificado pelo valor
- **Domain Service:** Lógica que não pertence a uma entidade
- **Repository Interface:** Contrato para acesso a dados
- **Domain Error:** Erro específico do domínio

## ⚠️ Common Mistakes

❌ **NÃO FAÇA:**
```go
// Dependência de framework no domain
type User struct {
    gorm.Model  // ❌ GORM no domain
}

// Lógica de infraestrutura no domain
func (u *User) ToJSON() string {  // ❌ HTTP no domain
    return gin.H{"id": u.ID}.String()
}
```

✅ **FAÇA:**
```go
// Só tipos Go nativos
type User struct {
    id    string
    email Email
}

// Só regras de negócio
func (u *User) IsActive() bool {
    return u.status == UserStatusActive
}
```

---

**Este é o coração da sua aplicação! Mantenha-o puro e focado nas regras de negócio.** 🧠
