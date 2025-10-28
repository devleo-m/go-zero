# 🧠 Camada de Domínio (Domain Layer)

> **🎯 Objetivo:** Centralizar as regras de negócio e entidades do sistema, garantindo que a lógica crítica esteja isolada e testável.

## 📚 O que é a Camada de Domínio?

A **Domain Layer** é o **coração** da sua aplicação! É aqui que vivem:

- 🏗️ **Entidades** - Objetos com identidade única
- 📋 **Regras de Negócio** - Lógica que não pode ser violada
- 🚫 **Validações** - Garantias de integridade dos dados
- 🔄 **Comportamentos** - Ações que as entidades podem realizar

## 🎓 Por que usar Domain-Driven Design (DDD)?

### ❌ **Problema sem DDD:**
```go
// Lógica espalhada em controllers
func CreateUser(c *gin.Context) {
    // Validação aqui
    if len(password) < 8 {
        return error
    }
    
    // Hash da senha aqui
    hashedPassword := bcrypt.Hash(password)
    
    // Salvar no banco aqui
    db.Save(user)
}
```

### ✅ **Solução com DDD:**
```go
// Lógica centralizada na entidade
func NewUser(name, email, password string) (*User, error) {
    // Todas as regras de negócio em um lugar
    if err := validateUserData(name, email, password); err != nil {
        return nil, err
    }
    
    hashedPassword, err := hashPassword(password)
    if err != nil {
        return nil, err
    }
    
    return &User{
        ID:       uuid.New(),
        Name:     name,
        Email:    email,
        Password: hashedPassword,
        // ... outros campos
    }, nil
}
```

## 🏗️ Estrutura da Camada

```
domain/
├── 📄 README.md           # Este arquivo - conceitos de domínio
├── user.go               # 🧠 Entidade User com regras de negócio
├── repository.go         # 📋 Interface do repositório
└── errors.go            # 🚫 Erros específicos do domínio
```

## 🧠 Entidade User - Análise Detalhada

### **1. Estrutura da Entidade**
```go
type User struct {
    ID        uuid.UUID  `json:"id"`         // Identidade única
    Name      string     `json:"name"`       // Nome do usuário
    Email     string     `json:"email"`      // Email único
    Password  string     `json:"-"`          // Senha (nunca serializar!)
    Phone     *string    `json:"phone"`      // Telefone opcional
    Role      string     `json:"role"`       // Papel no sistema
    Status    string     `json:"status"`  // Status do usuário
    CreatedAt time.Time  `json:"created_at"` // Data de criação
    UpdatedAt time.Time  `json:"updated_at"` // Data de atualização
    DeletedAt *time.Time `json:"deleted_at"` // Soft delete
}
```

### **2. Regras de Negócio Implementadas**

#### **🔐 Criação Segura de Usuário**
```go
func NewUser(name, email, password string) (*User, error) {
    // ✅ Validação de nome (mínimo 2 caracteres)
    if name == "" || len(name) < 2 {
        return nil, ErrInvalidName
    }
    
    // ✅ Validação de email (não pode ser vazio)
    if email == "" {
        return nil, ErrInvalidEmail
    }
    
    // ✅ Validação de senha (mínimo 8 caracteres)
    if password == "" || len(password) < 8 {
        return nil, ErrInvalidPassword
    }
    
    // ✅ Hash seguro da senha
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, ErrPasswordHash
    }
    
    // ✅ Valores padrão seguros
    return &User{
        ID:        uuid.New(),        // UUID único
        Name:      name,
        Email:     email,
        Password:  string(hashedPassword),
        Role:      "user",            // Papel padrão
        Status:    "active",          // Status ativo
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }, nil
}
```

#### **🔒 Validação de Senha**
```go
func (u *User) ValidatePassword(password string) error {
    return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
```

#### **🔄 Atualização de Perfil**
```go
func (u *User) UpdateProfile(name string, phone *string) error {
    // ✅ Validação de nome
    if name == "" || len(name) < 2 {
        return ErrInvalidName
    }
    
    // ✅ Atualização segura
    u.Name = name
    u.Phone = phone
    u.UpdatedAt = time.Now()  // Timestamp automático
    return nil
}
```

#### **🗑️ Soft Delete**
```go
func (u *User) SoftDelete() {
    now := time.Now()
    u.DeletedAt = &now      // Marca como deletado
    u.UpdatedAt = now       // Atualiza timestamp
}

func (u *User) IsDeleted() bool {
    return u.DeletedAt != nil
}
```

## 📋 Interface do Repositório

### **Por que usar Interface?**
- ✅ **Testabilidade** - Fácil criar mocks
- ✅ **Flexibilidade** - Trocar implementação sem afetar lógica
- ✅ **Inversão de Dependência** - Domain não depende de infraestrutura

```go
type Repository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id uuid.UUID) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    List(ctx context.Context, limit, offset int) ([]*User, int64, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

## 🚫 Erros de Domínio

### **Por que erros específicos?**
- ✅ **Clareza** - Cada erro tem um significado específico
- ✅ **Tratamento** - Handlers podem tratar cada erro adequadamente
- ✅ **Manutenibilidade** - Fácil de entender e corrigir

```go
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrEmailAlreadyInUse = errors.New("email already in use")
    ErrInvalidName       = errors.New("invalid name")
    ErrInvalidEmail      = errors.New("invalid email")
    ErrInvalidPassword   = errors.New("invalid password")
    ErrPasswordHash      = errors.New("failed to hash password")
)
```

## 🎯 Princípios Aplicados

### **1. Encapsulamento**
- Dados privados protegidos
- Métodos públicos para operações
- Validações internas

### **2. Imutabilidade**
- UUID gerado automaticamente
- Timestamps controlados pela entidade
- Senha sempre hasheada

### **3. Responsabilidade Única**
- User gerencia apenas dados de usuário
- Repository gerencia apenas persistência
- Errors definem apenas erros de domínio

## 🧪 Como Testar

### **Teste Unitário da Entidade:**
```go
func TestNewUser(t *testing.T) {
    user, err := NewUser("João", "joao@test.com", "12345678")
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "João", user.Name)
    assert.Equal(t, "joao@test.com", user.Email)
    assert.NotEmpty(t, user.ID)
    assert.True(t, user.ValidatePassword("12345678"))
}
```

## 🚀 Próximos Passos

1. **Explore o código** da entidade `User`
2. **Entenda as regras** de negócio implementadas
3. **Veja como** os erros são tratados
4. **Pratique** criando novas validações

---

> **💡 Dica:** A camada de domínio deve ser **independente** de frameworks, bancos de dados e tecnologias externas. Ela representa o "coração" do seu negócio!
