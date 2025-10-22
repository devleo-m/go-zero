# Módulo User - Domain Layer

## 🎯 O que é este módulo?

Este módulo contém toda a **lógica de negócio** relacionada aos usuários do sistema. É o **coração** do domínio User, seguindo os princípios da **Arquitetura Hexagonal** e **Domain-Driven Design (DDD)**.

## 🏗️ Estrutura do Módulo

```
internal/domain/user/
├── entity.go           # User Entity (regras de negócio)
├── value_objects.go    # Email, Password, Phone (Value Objects)
├── password_hash.go     # Lógica de hash de senha (Argon2)
├── types.go            # Role, Status (Enums/Types)
├── errors.go           # Erros específicos do domínio
└── repository.go       # Interface do Repository
```

## 🧠 Conceitos Fundamentais

### 1. **Entity (User)**
- Representa um usuário com **identidade única**
- Contém **regras de negócio** (autenticação, autorização, validações)
- É **mutável** (pode ser alterada ao longo do tempo)

### 2. **Value Objects (Email, Password, Phone)**
- Representam **conceitos** que têm valor, não identidade
- São **imutáveis** (nunca mudam)
- Sempre **válidos** quando criados

### 3. **Repository Interface**
- Define **O QUE** o domínio precisa para persistir dados
- **NÃO** define **COMO** persistir (isso fica na Infrastructure)
- Garante que o domínio não depende de tecnologias específicas

## 🔐 Segurança Implementada

### Senhas
- **Argon2** para hash de senhas (mais seguro que bcrypt)
- Validação de força: min 8 chars, maiúscula, minúscula, número, especial
- **Constant time comparison** para evitar timing attacks

### Autenticação
- Verificação de status do usuário antes do login
- Controle de tentativas de login
- Tokens JWT (implementado na Infrastructure)

### Autorização
- Sistema de **roles** hierárquico
- **Permissões** granulares
- Verificação de **capacidade de gerenciar** outros usuários

## 📋 Funcionalidades Principais

### Criação de Usuário
```go
user, err := NewUser(
    "João Silva",
    email,      // Value Object Email
    password,   // Value Object Password
    RoleUser,   // Enum Role
)
```

### Autenticação
```go
err := user.Authenticate("senha123")
if err != nil {
    // Tratar erro de autenticação
}
```

### Mudança de Senha
```go
err := user.ChangePassword("senhaAntiga", "novaSenha123!")
```

### Controle de Acesso
```go
if user.CanAccess("read_users") {
    // Usuário pode ler outros usuários
}

if user.CanManage(otherUser) {
    // Usuário pode gerenciar outro usuário
}
```

## 🎭 Roles e Permissões

### Hierarquia de Roles
1. **Admin** - Acesso total
2. **Manager** - Gerenciamento limitado
3. **User** - Acesso básico
4. **Guest** - Acesso mínimo

### Sistema de Permissões
- **Granular**: Cada ação tem uma permissão específica
- **Hierárquico**: Roles superiores herdam permissões dos inferiores
- **Flexível**: Fácil adicionar novas permissões

## 📊 Status do Usuário

### Estados Possíveis
- **Pending** - Aguardando ativação
- **Active** - Usuário ativo
- **Inactive** - Usuário desativado
- **Suspended** - Usuário suspenso

### Transições Válidas
- Pending → Active (ativação)
- Active → Inactive (desativação)
- Active → Suspended (suspensão)
- Suspended → Active (remoção de suspensão)

## 🔍 Repository Interface

### Operações Principais
```go
type Repository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
    FindByEmail(ctx context.Context, email Email) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, filters ListFilters) ([]*User, int64, error)
    // ... outras operações
}
```

### Filtros de Listagem
- **Paginação**: Page, PageSize
- **Filtros**: Role, Status, Search
- **Ordenação**: SortBy, SortOrder

## ⚠️ Erros de Domínio

### Erros Específicos
- `ErrUserNotFound` - Usuário não encontrado
- `ErrUserAlreadyExists` - Usuário já existe
- `ErrInvalidCredentials` - Credenciais inválidas
- `ErrUserInactive` - Usuário inativo
- `ErrEmailAlreadyInUse` - Email já em uso
- E muitos outros...

### Tratamento de Erros
```go
if err != nil {
    switch {
    case errors.Is(err, ErrUserNotFound):
        // Usuário não encontrado
    case errors.Is(err, ErrInvalidCredentials):
        // Credenciais inválidas
    default:
        // Erro genérico
    }
}
```

## 🧪 Testabilidade

### Princípios Seguidos
- **Zero dependências** externas
- **Interfaces** para todos os contratos
- **Métodos pequenos** e focados
- **Validações** explícitas

### Como Testar
```go
func TestUser_Authenticate(t *testing.T) {
    // Arrange
    user := createTestUser()
    
    // Act
    err := user.Authenticate("wrong_password")
    
    // Assert
    assert.Error(t, err)
    assert.Equal(t, ErrInvalidCredentials, err)
}
```

## 🔄 Fluxo de Dados

```
HTTP Request
    ↓
Use Case Layer
    ↓
Domain Layer (User Entity)
    ↓
Repository Interface
    ↓
Infrastructure Layer (GORM Implementation)
    ↓
Database
```

## 🎯 Regras de Negócio Implementadas

### Validações
- ✅ Nome obrigatório (2-100 caracteres)
- ✅ Email válido e único
- ✅ Senha forte (8+ chars, maiúscula, minúscula, número, especial)
- ✅ Telefone opcional mas válido
- ✅ Role válido
- ✅ Status válido

### Regras de Segurança
- ✅ Senhas nunca serializadas em JSON
- ✅ Hash seguro com Argon2
- ✅ Verificação de permissões antes de operações
- ✅ Prevenção de auto-desativação/suspensão
- ✅ Hierarquia de roles respeitada

### Regras de Negócio
- ✅ Usuários começam como "Pending"
- ✅ Apenas usuários ativos podem fazer login
- ✅ Controle de tentativas de login
- ✅ Auditoria de alterações (timestamps)

## 🚀 Próximos Passos

1. **Use Case Layer** - Orquestração das operações
2. **Infrastructure Layer** - Implementação do Repository com GORM
3. **Interface Layer** - Handlers HTTP com Gin
4. **Testes** - Unit tests, integration tests, E2E tests

## 📚 Conceitos Aprendidos

- ✅ **Domain-Driven Design** (DDD)
- ✅ **Arquitetura Hexagonal**
- ✅ **Value Objects** vs **Entities**
- ✅ **Repository Pattern**
- ✅ **Domain Errors**
- ✅ **Segurança** (Argon2, validações)
- ✅ **Sistema de Permissões**
- ✅ **Imutabilidade** e **Encapsulamento**

---

**🎉 Parabéns!** Você criou um Domain Layer **100% profissional** seguindo todas as melhores práticas da indústria!
