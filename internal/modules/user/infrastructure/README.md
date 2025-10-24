# 🗄️ Camada de Infraestrutura (Infrastructure Layer)

> **🎯 Objetivo:** Implementar detalhes técnicos de persistência, conectando o domínio com tecnologias externas como banco de dados, APIs, etc.

## 📚 O que é a Camada de Infraestrutura?

A **Infrastructure Layer** é o **braço técnico** da sua aplicação! É aqui que:

- 🗄️ **Persistimos** dados no banco
- 🔌 **Conectamos** com serviços externos
- 📡 **Integramos** com APIs
- 🛠️ **Implementamos** detalhes técnicos

## 🎓 Por que separar Infraestrutura?

### ❌ **Problema sem separação:**
```go
// Lógica de negócio misturada com SQL
func CreateUser(user *User) error {
    // Regras de negócio aqui
    if user.Email == "" {
        return errors.New("email required")
    }
    
    // SQL direto aqui
    query := "INSERT INTO users (name, email) VALUES (?, ?)"
    _, err := db.Exec(query, user.Name, user.Email)
    return err
}
```

### ✅ **Solução com Infrastructure Layer:**
```go
// Domain define a interface
type Repository interface {
    Create(ctx context.Context, user *User) error
}

// Infrastructure implementa os detalhes
func (r *PostgresRepository) Create(ctx context.Context, user *User) error {
    // Apenas detalhes de persistência
    model := toModel(user)
    return r.db.WithContext(ctx).Create(model).Error
}
```

## 🏗️ Estrutura da Camada

```
infrastructure/
├── 📄 README.md           # Este arquivo - conceitos de infraestrutura
└── postgres/              # 🗄️ Implementação PostgreSQL
    ├── repository.go      # Repositório PostgreSQL
    └── user.go           # Modelo de dados GORM
```

## 🗄️ Implementação PostgreSQL

### **1. Repository - Implementação da Interface**

```go
type Repository struct {
    db *gorm.DB  // Dependência do GORM
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}
```

### **2. Operações CRUD Implementadas**

#### **📝 Create - Criar Usuário**
```go
func (r *Repository) Create(ctx context.Context, user *domain.User) error {
    // 1️⃣ Converter domain.User para UserModel
    model := toModel(user)
    
    // 2️⃣ Persistir no banco
    if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    // 3️⃣ Atualizar ID gerado pelo banco
    user.ID = model.ID
    return nil
}
```

#### **🔍 GetByID - Buscar por ID**
```go
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
    var model UserModel
    
    // 1️⃣ Query com filtro de soft delete
    if err := r.db.WithContext(ctx).
        Where("id = ? AND deleted_at IS NULL", id).
        First(&model).Error; err != nil {
        
        // 2️⃣ Tratar erro específico do GORM
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to get user by ID: %w", err)
    }
    
    // 3️⃣ Converter para domain.User
    return toDomain(&model), nil
}
```

#### **📧 GetByEmail - Buscar por Email**
```go
func (r *Repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    var model UserModel
    
    // 1️⃣ Query com filtro de soft delete
    if err := r.db.WithContext(ctx).
        Where("email = ? AND deleted_at IS NULL", email).
        First(&model).Error; err != nil {
        
        // 2️⃣ Tratar erro específico do GORM
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }
    
    // 3️⃣ Converter para domain.User
    return toDomain(&model), nil
}
```

#### **📄 List - Listar com Paginação**
```go
func (r *Repository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
    var models []UserModel
    
    // 1️⃣ Query com paginação e filtro de soft delete
    if err := r.db.WithContext(ctx).
        Where("deleted_at IS NULL").
        Limit(limit).
        Offset(offset).
        Find(&models).Error; err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }
    
    // 2️⃣ Converter slice de models para domain.Users
    users := make([]*domain.User, len(models))
    for i, model := range models {
        users[i] = toDomain(&model)
    }
    
    return users, nil
}
```

#### **🔄 Update - Atualizar Usuário**
```go
func (r *Repository) Update(ctx context.Context, user *domain.User) error {
    // 1️⃣ Converter domain.User para UserModel
    model := toModel(user)
    
    // 2️⃣ Salvar alterações
    if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    return nil
}
```

#### **🗑️ Delete - Soft Delete**
```go
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
    now := time.Now()
    
    // 1️⃣ Soft delete - apenas marca deleted_at
    err := r.db.WithContext(ctx).Model(&UserModel{}).
        Where("id = ?", id).
        Update("deleted_at", now).Error
    
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    return nil
}
```

## 🔄 Mapeamento Domain ↔ Infrastructure

### **toModel - Domain → Database**
```go
func toModel(user *domain.User) *UserModel {
    return &UserModel{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Password:  user.Password,
        Phone:     user.Phone,
        Role:      user.Role,
        Status:    user.Status,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
        DeletedAt: gorm.DeletedAt{},  // GORM gerencia soft delete
    }
}
```

### **toDomain - Database → Domain**
```go
func toDomain(model *UserModel) *domain.User {
    var deletedAt *time.Time
    if model.DeletedAt.Valid {
        deletedAt = &model.DeletedAt.Time
    }
    
    return &domain.User{
        ID:        model.ID,
        Name:      model.Name,
        Email:     model.Email,
        Password:  model.Password,
        Phone:     model.Phone,
        Role:      model.Role,
        Status:    model.Status,
        CreatedAt: model.CreatedAt,
        UpdatedAt: model.UpdatedAt,
        DeletedAt: deletedAt,
    }
}
```

## 🎯 Padrões Aplicados

### **1. Repository Pattern**
- ✅ **Abstração** - Domain não conhece detalhes do banco
- ✅ **Testabilidade** - Fácil criar mocks
- ✅ **Flexibilidade** - Trocar banco sem afetar lógica

### **2. Data Mapper Pattern**
- ✅ **Separação** - Domain e Database models separados
- ✅ **Conversão** - Funções específicas para mapear
- ✅ **Evolução** - Mudanças no banco não afetam domain

### **3. Soft Delete**
- ✅ **Auditoria** - Dados não são perdidos
- ✅ **Recuperação** - Possível restaurar dados
- ✅ **Performance** - Queries filtram deletados

## 🧪 Como Testar

### **Teste de Integração:**
```go
func TestUserRepository(t *testing.T) {
    // Arrange
    db := setupTestDB()
    repo := NewRepository(db)
    
    user := &domain.User{
        ID:    uuid.New(),
        Name:  "João",
        Email: "joao@test.com",
    }
    
    // Act
    err := repo.Create(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    
    // Verify
    found, err := repo.GetByID(context.Background(), user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Name, found.Name)
}
```

## 🎓 Benefícios da Separação

### **1. Independência**
- ✅ **Domain** não depende de frameworks
- ✅ **Infrastructure** pode ser trocada facilmente
- ✅ **Testes** isolados por camada

### **2. Manutenibilidade**
- ✅ **Mudanças** no banco não afetam lógica
- ✅ **Evolução** independente das camadas
- ✅ **Debugging** mais fácil

### **3. Reutilização**
- ✅ **Repository** pode ser usado em diferentes contextos
- ✅ **Domain** pode ser usado com diferentes bancos
- ✅ **Use Cases** funcionam independente da persistência

## 🚀 Próximos Passos

1. **Explore o código** do repositório PostgreSQL
2. **Entenda** como o mapeamento funciona
3. **Veja** como os erros são tratados
4. **Pratique** criando outros repositórios (Redis, MongoDB, etc.)

---

> **💡 Dica:** A camada de infraestrutura deve **implementar** as interfaces definidas pelo domínio, mas nunca **definir** regras de negócio!
