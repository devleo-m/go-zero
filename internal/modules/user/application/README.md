# 🎯 Camada de Aplicação (Application Layer)

> **🎯 Objetivo:** Orquestrar casos de uso do sistema, coordenando entre a camada de domínio e infraestrutura sem conhecer detalhes de implementação.

## 📚 O que é a Camada de Aplicação?

A **Application Layer** é o **maestro** da sua aplicação! É aqui que:

- 🎭 **Orquestramos** operações complexas
- 🔄 **Coordenamos** entre domínio e infraestrutura  
- 📋 **Definimos** casos de uso específicos
- 🧪 **Facilitamos** testes unitários

## 🎓 Por que separar Use Cases?

### ❌ **Problema sem Use Cases:**
```go
// Lógica espalhada no controller
func CreateUser(c *gin.Context) {
    // Validação
    // Verificar email único
    // Criar usuário
    // Salvar no banco
    // Enviar email
    // Log de auditoria
    // ... tudo misturado!
}
```

### ✅ **Solução com Use Cases:**
```go
// Use Case específico e testável
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
    // 1. Verificar regras de negócio
    // 2. Orquestrar operações
    // 3. Retornar resultado
}
```

## 🏗️ Estrutura da Camada

```
application/
├── 📄 README.md           # Este arquivo - conceitos de aplicação
├── create_user.go        # 🎯 Caso de uso: criar usuário
├── get_user.go           # 🎯 Caso de uso: buscar usuário
├── list_users.go         # 🎯 Caso de uso: listar usuários
├── update_user.go        # 🎯 Caso de uso: atualizar usuário
└── delete_user.go        # 🎯 Caso de uso: deletar usuário
```

## 🎯 Use Cases Implementados

### **1. CreateUserUseCase - Criar Usuário**

#### **📋 Estrutura do Use Case:**
```go
type CreateUserUseCase struct {
    userRepo domain.Repository  // Dependência injetada
}

type CreateUserInput struct {
    Name     string  `json:"name" validate:"required,min=2,max=100"`
    Email    string  `json:"email" validate:"required,email"`
    Password string  `json:"password" validate:"required,min=8"`
    Phone    *string `json:"phone,omitempty"`
}

type CreateUserOutput struct {
    User    *domain.User `json:"user"`
    Message string       `json:"message"`
}
```

#### **🔄 Fluxo de Execução:**
```go
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
    // 1️⃣ Verificar se email já existe
    existingUser, err := uc.userRepo.GetByEmail(ctx, input.Email)
    if err != nil && err != domain.ErrUserNotFound {
        return nil, fmt.Errorf("failed to check email: %w", err)
    }
    if existingUser != nil {
        return nil, domain.ErrEmailAlreadyInUse
    }

    // 2️⃣ Criar usuário usando regras de domínio
    user, err := domain.NewUser(input.Name, input.Email, input.Password)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    // 3️⃣ Definir telefone se fornecido
    if input.Phone != nil {
        user.Phone = input.Phone
    }

    // 4️⃣ Salvar no banco
    if err := uc.userRepo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to save user: %w", err)
    }

    // 5️⃣ Retornar resultado
    return &CreateUserOutput{
        User:    user,
        Message: "User created successfully",
    }, nil
}
```

### **2. GetUserUseCase - Buscar Usuário**

#### **🔍 Lógica Simples e Direta:**
```go
func (uc *GetUserUseCase) Execute(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
    // 1️⃣ Buscar usuário por ID
    user, err := uc.userRepo.GetByID(ctx, input.ID)
    if err != nil {
        if err == domain.ErrUserNotFound {
            return nil, domain.ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    // 2️⃣ Retornar usuário encontrado
    return &GetUserOutput{
        User: user,
    }, nil
}
```

### **3. ListUsersUseCase - Listar Usuários**

#### **📄 Paginação e Filtros:**
```go
func (uc *ListUsersUseCase) Execute(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
    // 1️⃣ Validar parâmetros de paginação
    if input.Limit <= 0 || input.Limit > 100 {
        input.Limit = 10
    }
    if input.Offset < 0 {
        input.Offset = 0
    }

    // 2️⃣ Buscar usuários com paginação
    users, total, err := uc.userRepo.List(ctx, input.Limit, input.Offset)
    if err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }

    // 3️⃣ Retornar lista paginada
    return &ListUsersOutput{
        Users: users,
        Total: total,
    }, nil
}
```

### **4. UpdateUserUseCase - Atualizar Usuário**

#### **🔄 Atualização Segura:**
```go
func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error) {
    // 1️⃣ Buscar usuário existente
    user, err := uc.userRepo.GetByID(ctx, input.ID)
    if err != nil {
        if err == domain.ErrUserNotFound {
            return nil, domain.ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    // 2️⃣ Atualizar usando regras de domínio
    if err := user.UpdateProfile(input.Name, input.Phone); err != nil {
        return nil, fmt.Errorf("failed to update profile: %w", err)
    }

    // 3️⃣ Salvar alterações
    if err := uc.userRepo.Update(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to save user: %w", err)
    }

    // 4️⃣ Retornar resultado
    return &UpdateUserOutput{
        User:    user,
        Message: "User updated successfully",
    }, nil
}
```

### **5. DeleteUserUseCase - Deletar Usuário**

#### **🗑️ Soft Delete:**
```go
func (uc *DeleteUserUseCase) Execute(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, error) {
    // 1️⃣ Buscar usuário existente
    user, err := uc.userRepo.GetByID(ctx, input.ID)
    if err != nil {
        if err == domain.ErrUserNotFound {
            return nil, domain.ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    // 2️⃣ Soft delete usando regras de domínio
    user.SoftDelete()

    // 3️⃣ Salvar alterações
    if err := uc.userRepo.Update(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to delete user: %w", err)
    }

    // 4️⃣ Retornar resultado
    return &DeleteUserOutput{
        Message: "User deleted successfully",
    }, nil
}
```

## 🎯 Padrões Aplicados

### **1. Use Case Pattern**
- ✅ **Uma responsabilidade** por use case
- ✅ **Input/Output** bem definidos
- ✅ **Orquestração** clara de operações

### **2. Dependency Injection**
- ✅ **Repositório injetado** no construtor
- ✅ **Fácil de testar** com mocks
- ✅ **Inversão de dependência**

### **3. Error Handling**
- ✅ **Erros de domínio** preservados
- ✅ **Contexto** adicionado aos erros
- ✅ **Tratamento** consistente

## 🧪 Como Testar Use Cases

### **Teste com Mock:**
```go
func TestCreateUserUseCase(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    useCase := NewCreateUserUseCase(mockRepo)
    
    input := CreateUserInput{
        Name:     "João",
        Email:    "joao@test.com",
        Password: "12345678",
    }
    
    // Act
    result, err := useCase.Execute(context.Background(), input)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "João", result.User.Name)
    assert.True(t, mockRepo.CreateCalled)
}
```

## 🎓 Benefícios dos Use Cases

### **1. Testabilidade**
- ✅ **Isolados** - Fácil de testar unitariamente
- ✅ **Mocks** - Dependências podem ser mockadas
- ✅ **Cobertura** - Teste de lógica de negócio

### **2. Reutilização**
- ✅ **HTTP** - Pode ser usado em controllers
- ✅ **CLI** - Pode ser usado em comandos
- ✅ **gRPC** - Pode ser usado em serviços

### **3. Manutenibilidade**
- ✅ **Clareza** - Cada use case tem uma responsabilidade
- ✅ **Evolução** - Fácil de modificar sem afetar outros
- ✅ **Debugging** - Fácil de encontrar problemas

## 🚀 Próximos Passos

1. **Explore cada use case** individualmente
2. **Entenda o fluxo** de cada operação
3. **Veja como** os erros são tratados
4. **Pratique** criando novos use cases

---

> **💡 Dica:** Use Cases são a **ponte** entre o que o usuário quer fazer e como o sistema vai fazer. Eles orquestram, mas não implementam detalhes técnicos!
