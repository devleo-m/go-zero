# 🗄️ PERSISTENCE LAYER - PostgreSQL Adapter

## O que é esta camada?

A **Persistence Layer** é a camada de **Infrastructure** que implementa a comunicação com o banco de dados PostgreSQL usando GORM. Ela traduz entre as entidades do Domain e as tabelas do banco de dados.

## 🎯 Conceitos Fundamentais

### Separação Domain vs Infrastructure

```
┌─────────────────────────────────────────────────────────┐
│ DOMAIN LAYER (Regras de Negócio)                       │
│                                                         │
│  type User struct {                                     │
│      Email    Email     // Value Object com validação  │
│      Password Password  // Hash seguro                 │
│  }                                                      │
│                                                         │
│  func (u *User) Authenticate(pass string) error {      │
│      // Lógica de autenticação                         │
│  }                                                      │
└─────────────────────────────────────────────────────────┘
                          ↕️ CONVERTER
┌─────────────────────────────────────────────────────────┐
│ INFRASTRUCTURE LAYER (Persistência)                     │
│                                                         │
│  type UserModel struct {                                │
│      Email    string `gorm:"uniqueIndex"`  // String   │
│      Password string `gorm:"not null"`     // String   │
│  }                                                      │
│                                                         │
│  // Sem métodos de negócio! Só dados.                  │
└─────────────────────────────────────────────────────────┘
                          ↕️
                   PostgreSQL Database
```

### Por que separar?

✅ **Trocar de banco sem dor**: PostgreSQL → MongoDB (só muda Infrastructure)  
✅ **Testar sem banco**: Domain testa sozinho  
✅ **Código limpo**: Cada camada tem UMA responsabilidade  
✅ **Manutenibilidade**: Bug no banco? Olhe Infrastructure. Bug na regra? Domain.

## 📁 Estrutura

```
postgres/
├── models/                 # GORM Models (mapeamento de tabelas)
│   └── user_model.go      # UserModel com tags GORM
│
├── converters/            # Tradutores Domain ↔ Model
│   └── user_converter.go  # UserConverter
│
└── repositories/          # Implementações de Repository
    └── user_repository.go # UserRepository (implementa interface do Domain)
```

## 🗄️ 1. Models (GORM)

### O que é?

Um **Model** é uma struct Go com tags GORM que mapeia para uma tabela PostgreSQL.

### Exemplo: UserModel

```go
type UserModel struct {
    // Primary Key and Timestamps
    ID        string         `gorm:"type:varchar(36);primaryKey;not null"`
    CreatedAt time.Time      `gorm:"not null;autoCreateTime"`
    UpdatedAt time.Time      `gorm:"not null;autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete
    
    // User Information
    Name     string `gorm:"type:varchar(255);not null"`
    Email    string `gorm:"type:varchar(255);not null;uniqueIndex:idx_email"`
    Password string `gorm:"type:varchar(255);not null"`
    
    // Status and Role
    Status     string `gorm:"type:varchar(50);not null;index:idx_status"`
    Role       string `gorm:"type:varchar(50);not null;index:idx_role"`
    IsVerified bool   `gorm:"not null;default:false"`
    
    // Audit
    LastLoginAt *time.Time `gorm:"type:timestamp"`
}
```

### Tags GORM Explicadas:

| Tag | Significado | Exemplo |
|-----|-------------|---------|
| `type:varchar(255)` | Tipo da coluna SQL | VARCHAR no PostgreSQL |
| `primaryKey` | Chave primária | ID único |
| `not null` | Campo obrigatório | Não aceita NULL |
| `uniqueIndex` | Índice único | Email não duplica |
| `index` | Índice normal | Busca rápida |
| `default:false` | Valor padrão | Valor inicial |
| `autoCreateTime` | Timestamp automático | CreatedAt gerado |
| `autoUpdateTime` | Timestamp automático | UpdatedAt atualizado |

### Hooks do GORM

```go
// BeforeCreate: executa ANTES de inserir
func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
    if u.ID == "" {
        u.ID = generateID()
    }
    return nil
}

// BeforeUpdate: executa ANTES de atualizar
func (u *UserModel) BeforeUpdate(tx *gorm.DB) error {
    u.UpdatedAt = time.Now()
    return nil
}
```

### Scopes (Queries Reutilizáveis)

```go
// ActiveUsers filtra apenas usuários ativos
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("status = ? AND deleted_at IS NULL", "active")
}

// VerifiedUsers filtra apenas verificados
func VerifiedUsers(db *gorm.DB) *gorm.DB {
    return db.Where("is_verified = ?", true)
}

// Uso:
db.Scopes(ActiveUsers, VerifiedUsers).Find(&users)
```

## 🔄 2. Converters (Tradutores)

### O que é?

O **Converter** é o **TRADUTOR** entre duas linguagens diferentes:
- **Domain Entity** (com Value Objects, comportamentos)
- **GORM Model** (struct simples para o banco)

### Métodos Principais:

#### ToModel: Domain → GORM (para salvar)

```go
func (c *UserConverter) ToModel(user *entities.User) *models.UserModel {
    return &models.UserModel{
        ID:         user.ID(),
        Name:       user.FullName(),
        Email:      user.Email().String(),    // Email VO → string
        Password:   user.Password().Hash(),   // Password VO → hash
        Status:     string(user.Status()),    // enum → string
        Role:       string(user.Role()),      // enum → string
        IsVerified: user.IsVerified(),
        LastLoginAt: user.LastLoginAt(),
        CreatedAt:  user.CreatedAt(),
        UpdatedAt:  user.UpdatedAt(),
    }
}
```

#### ToDomain: GORM → Domain (depois de buscar)

```go
func (c *UserConverter) ToDomain(model *models.UserModel) (*entities.User, error) {
    // Parse enums
    status, err := c.parseUserStatus(model.Status)
    if err != nil {
        return nil, err
    }
    
    role, err := c.parseUserRole(model.Role)
    if err != nil {
        return nil, err
    }
    
    // Reconstruct User from database data
    user, err := entities.NewUserFromData(
        model.ID,
        model.Email,
        model.Password,
        model.Name,
        role,
        status,
        model.IsVerified,
        model.LastLoginAt,
        model.CreatedAt,
        model.UpdatedAt,
        getDeletedAt(model.DeletedAt),
    )
    
    return user, err
}
```

### Fluxo na Prática:

**CREATE** (Salvar):  
`User (Domain) → ToModel → UserModel (GORM) → PostgreSQL`

**FIND** (Buscar):  
`PostgreSQL → UserModel (GORM) → ToDomain → User (Domain)`

## 💾 3. Repositories (Implementações)

### O que é?

O **Repository** implementa a **interface** definida no Domain, usando GORM para acessar o banco de dados.

### Estrutura:

```go
type UserRepository struct {
    db        *gorm.DB              // Conexão GORM
    converter *converters.UserConverter  // Tradutor
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
    return &UserRepository{
        db:        db,
        converter: converters.NewUserConverter(),
    }
}
```

### CRUD Completo:

#### CREATE

```go
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
    // 1. Converter Domain → Model
    model := r.converter.ToModel(user)
    
    // 2. Salvar no banco
    result := r.db.WithContext(ctx).Create(model)
    if result.Error != nil {
        return r.handleError(result.Error, "create user")
    }
    
    return nil
}
```

#### READ (GetByID)

```go
func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
    var model models.UserModel
    
    // 1. Buscar do banco
    result := r.db.WithContext(ctx).First(&model, "id = ?", id)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, domainErrors.ErrUserNotFound
        }
        return nil, r.handleError(result.Error, "find user by ID")
    }
    
    // 2. Converter Model → Domain
    return r.converter.ToDomain(&model)
}
```

#### UPDATE

```go
func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
    // 1. Buscar model existente
    var existingModel models.UserModel
    if err := r.db.First(&existingModel, "id = ?", user.ID()).Error; err != nil {
        return domainErrors.ErrUserNotFound
    }
    
    // 2. Atualizar model com dados do domain
    r.converter.UpdateModelFromDomain(&existingModel, user)
    
    // 3. Salvar
    return r.db.Save(&existingModel).Error
}
```

#### DELETE (Soft Delete)

```go
func (r *UserRepository) Delete(ctx context.Context, id string) error {
    // GORM faz soft delete automaticamente quando DeletedAt existe
    result := r.db.Delete(&models.UserModel{}, "id = ?", id)
    
    if result.RowsAffected == 0 {
        return domainErrors.ErrUserNotFound
    }
    
    return result.Error
}
```

### Queries Avançadas:

#### List com Filtros

```go
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
    var userModels []*models.UserModel
    
    query := r.db.WithContext(ctx).Model(&models.UserModel{})
    
    // Paginação
    if limit > 0 {
        query = query.Limit(limit)
    }
    if offset > 0 {
        query = query.Offset(offset)
    }
    
    // Ordenação
    query = query.Order("created_at DESC")
    
    // Executar
    result := query.Find(&userModels)
    if result.Error != nil {
        return nil, r.handleError(result.Error, "list users")
    }
    
    // Converter todos
    return r.converter.ToDomains(userModels)
}
```

#### Search (Busca por Nome ou Email)

```go
func (r *UserRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entities.User, error) {
    var userModels []*models.UserModel
    
    searchPattern := "%" + query + "%"
    
    q := r.db.WithContext(ctx).
        Model(&models.UserModel{}).
        Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern).
        Limit(limit).
        Offset(offset).
        Order("created_at DESC")
    
    result := q.Find(&userModels)
    if result.Error != nil {
        return nil, r.handleError(result.Error, "search users")
    }
    
    return r.converter.ToDomains(userModels)
}
```

### Tratamento de Erros:

```go
func (r *UserRepository) handleError(err error, operation string) error {
    // Not found
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return domainErrors.ErrUserNotFound
    }
    
    // Email duplicado
    if strings.Contains(err.Error(), "duplicate") && 
       strings.Contains(err.Error(), "email") {
        return domainErrors.ErrEmailAlreadyInUse
    }
    
    // Foreign key violation
    if strings.Contains(err.Error(), "foreign key") {
        return domainErrors.NewDomainError(
            "FOREIGN_KEY_VIOLATION",
            "violação de chave estrangeira",
            map[string]interface{}{"operation": operation},
        )
    }
    
    // Erro genérico (sem expor detalhes internos)
    return domainErrors.NewDomainError(
        "DATABASE_ERROR",
        fmt.Sprintf("erro ao %s", operation),
        map[string]interface{}{"operation": operation},
    )
}
```

**Por que isso é importante?**

- ✅ Converte erros do GORM para erros do Domain
- ✅ Mantém arquitetura limpa (Domain não conhece GORM)
- ✅ Não expõe detalhes internos do banco
- ✅ Retorna erros específicos e significativos

## 🔥 Métodos Implementados

### CRUD Básico
- ✅ `Create` - Criar usuário
- ✅ `GetByID` - Buscar por ID
- ✅ `GetByEmail` - Buscar por email
- ✅ `Update` - Atualizar usuário
- ✅ `Delete` - Soft delete

### Queries Avançadas
- ✅ `List` - Listar com paginação
- ✅ `Count` - Contar usuários
- ✅ `Exists` - Verificar existência
- ✅ `ExistsByEmail` - Verificar email
- ✅ `ListByRole` - Filtrar por papel
- ✅ `ListByStatus` - Filtrar por status
- ✅ `Search` - Buscar por nome/email

### Queries Específicas
- ✅ `GetActiveUsers` - Usuários ativos
- ✅ `GetVerifiedUsers` - Usuários verificados
- ✅ `GetUsersCreatedBetween` - Por período
- ✅ `GetUsersByLastLogin` - Por último login

### Operações em Lote
- ✅ `BulkUpdate` - Atualizar múltiplos
- ✅ `BulkDelete` - Deletar múltiplos

### Soft Delete Avançado
- ✅ `Restore` - Restaurar deletado
- ✅ `GetDeletedUsers` - Listar deletados
- ✅ `HardDelete` - Deletar permanentemente

### Estatísticas
- ✅ `GetUserStats` - Estatísticas completas

## 🎯 Princípios Seguidos

### 1. **Dependency Injection**
Repository recebe `*gorm.DB` no construtor, não cria conexão própria.

### 2. **Interface Implementation**
Repository implementa **100%** da interface do Domain.

### 3. **Error Handling**
Todos os erros são tratados e convertidos para erros do Domain.

### 4. **Context Usage**
Todas as operações usam `context.Context` para:
- Cancelamento de operações
- Timeouts
- Propagação de valores

### 5. **Separation of Concerns**
- **Model**: Apenas estrutura de dados
- **Converter**: Apenas tradução
- **Repository**: Apenas acesso a dados

## 🚀 Como Usar

### 1. Criar Repository

```go
import (
    "github.com/devleo-m/go-zero/internal/modules/ecommerce/adapters/persistence/postgres/repositories"
    "gorm.io/gorm"
)

// Criar repository
userRepo := repositories.NewUserRepository(db)
```

### 2. Criar Usuário

```go
import "github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/entities"

user, _ := entities.NewUser(
    "user@example.com",
    "Password123!",
    "João Silva",
    entities.UserRoleClient,
)

err := userRepo.Create(ctx, user)
```

### 3. Buscar Usuário

```go
user, err := userRepo.GetByID(ctx, "user-id-123")
if err != nil {
    if errors.Is(err, domainErrors.ErrUserNotFound) {
        // Usuário não encontrado
    }
}
```

### 4. Listar com Filtros

```go
users, err := userRepo.ListByRole(ctx, entities.UserRoleClient, 10, 0)
```

### 5. Buscar Usuários

```go
users, err := userRepo.Search(ctx, "joão", 20, 0)
```

## ⚠️ Common Mistakes

❌ **NÃO FAÇA:**

```go
// ❌ Lógica de negócio no repository
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
    if user.IsAdmin() {
        // Lógica de negócio aqui é ERRADO!
        user.Activate()
    }
    // ...
}

// ❌ Expor erros do GORM ao Domain
func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
    // ...
    return nil, result.Error  // ❌ ERRADO! Expõe GORM ao Domain
}

// ❌ Criar conexão dentro do repository
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
    db, _ := gorm.Open(...)  // ❌ ERRADO!
    // ...
}
```

✅ **FAÇA:**

```go
// ✅ Repository só acessa dados
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
    model := r.converter.ToModel(user)
    return r.db.Create(model).Error
}

// ✅ Converta erros para Domain errors
func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
    // ...
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, domainErrors.ErrUserNotFound  // ✅ CORRETO!
    }
    // ...
}

// ✅ Receba DB por injeção de dependência
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
    return &UserRepository{db: db}  // ✅ CORRETO!
}
```

## 📊 Próximos Passos

1. ✅ **Models** - Criado
2. ✅ **Converters** - Criado
3. ✅ **Repositories** - Criado
4. ⏳ **Tests** - Próximo passo!
5. ⏳ **Use Cases** - Depois dos testes

---

**Esta é a ponte entre seu Domain puro e o PostgreSQL!** 🗄️

Mantenha-a focada apenas em persistência de dados, sem lógica de negócio.

