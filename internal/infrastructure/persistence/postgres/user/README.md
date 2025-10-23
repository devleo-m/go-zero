# User Repository - PostgreSQL Implementation 🚀

## Estrutura do Padrão

```
internal/infrastructure/persistence/postgres/user/
├── model.go          # GORM Model (tabela users)
├── converter.go      # Domain ↔ Model conversion
├── repository.go     # Repository implementation
├── queries.go        # Queries específicas otimizadas
└── README.md         # Esta documentação
```

## 1. Model (model.go)

**Propósito:** Mapeia diretamente para a tabela `users` no PostgreSQL

**Características:**
- ✅ Campos base (ID, timestamps, soft delete)
- ✅ Campos específicos do User
- ✅ Hooks GORM (BeforeCreate, BeforeUpdate, BeforeDelete)
- ✅ Métodos de conveniência (IsActive, IsAdmin, etc.)
- ✅ Índices e constraints definidos
- ✅ Tags GORM para otimização

**Exemplo:**
```go
type UserModel struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    CreatedAt time.Time      `gorm:"not null;index"`
    UpdatedAt time.Time      `gorm:"not null;index"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
    
    Name         string `gorm:"type:varchar(255);not null;index"`
    Email        string `gorm:"type:varchar(255);not null;uniqueIndex"`
    PasswordHash string `gorm:"type:varchar(255);not null"`
    // ... outros campos
}
```

## 2. Converter (converter.go)

**Propósito:** Conversão bidirecional entre Domain e Model

**Funções principais:**
- `ToModel(domainUser *user.User) *UserModel` - Domain → Model
- `ToDomain(model *UserModel) (*user.User, error)` - Model → Domain
- `ToDomainSlice(models []*UserModel) ([]*user.User, error)` - Lote
- `QueryFilterToGORM(filter shared.QueryFilter)` - Filtros

**Exemplo:**
```go
// Domain → Model
func ToModel(domainUser *user.User) *UserModel {
    return &UserModel{
        ID:        domainUser.ID,
        Name:      domainUser.Name,
        Email:     domainUser.Email.String(),
        Status:    domainUser.Status.String(),
        Role:      domainUser.Role.String(),
        // ...
    }
}

// Model → Domain
func ToDomain(model *UserModel) (*user.User, error) {
    email, err := user.NewEmail(model.Email)
    if err != nil {
        return nil, err
    }
    
    return &user.User{
        BaseEntity: shared.BaseEntity{...},
        Name:       model.Name,
        Email:      email,
        // ...
    }, nil
}
```

## 3. Repository (repository.go)

**Propósito:** Implementa a interface `shared.Repository[*user.User]`

**Métodos implementados:**
- ✅ CRUD básico (Create, FindOne, FindMany, Update, Delete)
- ✅ Paginação (Paginate)
- ✅ Agregações (Count, Exists, Sum, Avg, Min, Max)
- ✅ Operações em lote (CreateMany, UpdateMany, DeleteMany)
- ✅ Transações (WithTransaction)
- ✅ Queries avançadas (FindFirst, FindLast, Distinct, GroupBy)

**Exemplo:**
```go
func (r *Repository) Create(ctx context.Context, entity *user.User) error {
    model := ToModel(entity)
    
    if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    // Atualizar entidade com dados do banco
    updatedUser, err := ToDomain(model)
    if err != nil {
        return fmt.Errorf("failed to convert model to domain: %w", err)
    }
    *entity = *updatedUser
    
    return nil
}
```

## 4. Queries (queries.go)

**Propósito:** Queries específicas otimizadas para performance

**Categorias:**
- 🔍 **Busca específica:** FindByEmail, FindByPhone, FindByStatus
- 📊 **Estatísticas:** GetUserStats, contadores por role/status
- 🔎 **Busca e filtros:** SearchUsers, FindUsersByDateRange
- 🔑 **Tokens:** FindByPasswordResetToken, FindByActivationToken
- 📈 **Performance:** FindUsersWithPagination, queries otimizadas
- 🧹 **Manutenção:** FindExpiredTokens, CleanExpiredTokens

**Exemplo:**
```go
func (r *Repository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
    var model UserModel
    
    if err := r.db.WithContext(ctx).
        Where("email = ?", email).
        First(&model).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, shared.NewDomainError("USER_NOT_FOUND", "user not found")
        }
        return nil, err
    }
    
    return ToDomain(&model)
}
```

## Padrão para Outras Entidades

### 1. Criar estrutura de pastas
```bash
mkdir -p internal/infrastructure/persistence/postgres/{entity_name}
```

### 2. Criar arquivos seguindo o padrão
- `model.go` - GORM Model
- `converter.go` - Domain ↔ Model
- `repository.go` - Repository implementation
- `queries.go` - Queries específicas
- `README.md` - Documentação

### 3. Estrutura do Model
```go
type {Entity}Model struct {
    // Campos base
    ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    CreatedAt time.Time      `gorm:"not null;index"`
    UpdatedAt time.Time      `gorm:"not null;index"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
    
    // Campos específicos da entidade
    // ... campos específicos
    
    // Hooks GORM
    func (m *{Entity}Model) BeforeCreate(tx *gorm.DB) error
    func (m *{Entity}Model) BeforeUpdate(tx *gorm.DB) error
    func (m *{Entity}Model) BeforeDelete(tx *gorm.DB) error
}
```

### 4. Estrutura do Converter
```go
// Domain → Model
func ToModel(domainEntity *domain.{Entity}) *{Entity}Model {
    // Converter campos da entidade de domínio para modelo
}

// Model → Domain
func ToDomain(model *{Entity}Model) (*domain.{Entity}, error) {
    // Converter campos do modelo para entidade de domínio
    // Criar Value Objects quando necessário
}

// Lote
func ToDomainSlice(models []*{Entity}Model) ([]*domain.{Entity}, error)
func ToModelSlice(entities []*domain.{Entity}) []*{Entity}Model
```

### 5. Estrutura do Repository
```go
type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}

// Implementar todos os métodos da interface shared.Repository[T]
```

### 6. Estrutura das Queries
```go
// Queries específicas otimizadas
func (r *Repository) FindBy{Field}(ctx context.Context, value string) (*domain.{Entity}, error)
func (r *Repository) Get{Entity}Stats(ctx context.Context) (*{Entity}Stats, error)
func (r *Repository) Search{Entity}s(ctx context.Context, query string) ([]*domain.{Entity}, error)
// ... outras queries específicas
```

## Boas Práticas

### ✅ FAÇA:
- Use índices para campos frequentemente consultados
- Implemente hooks GORM para validações
- Trate erros com wrapping (`fmt.Errorf`)
- Use transações para operações complexas
- Documente queries complexas
- Teste conversões Domain ↔ Model
- Use Value Objects do domínio
- Implemente soft delete
- Otimize queries com índices

### ❌ NÃO FAÇA:
- Lógica de negócio no Model
- Validações de domínio no Model
- Dependências do domínio no Model
- Queries N+1 (use Preload)
- Ignorar erros de conversão
- Hardcode de strings (use constantes)
- Queries sem índices
- Conversões sem tratamento de erro

## Exemplo de Uso

```go
// Injeção de dependência
userRepo := user.NewRepository(db)

// Usar QueryBuilder
filter := shared.NewQueryBuilder().
    WhereEqual("status", "active").
    WhereIn("role", []interface{}{"admin", "user"}).
    OrderByDesc("created_at").
    Page(1).
    PageSize(20).
    Build()

// Buscar usuários
users, err := userRepo.FindMany(ctx, filter)

// Usar Specification
activeAdmins := shared.ActiveAdminsSpecification[user.User]()
users, err := userRepo.FindMany(ctx, activeAdmins.ToQueryFilter())

// Paginação
result, err := userRepo.Paginate(ctx, filter)
```

## Conclusão

Este padrão garante:
- ✅ **Separação clara** entre Domain e Infrastructure
- ✅ **Conversões seguras** entre camadas
- ✅ **Performance otimizada** com índices
- ✅ **Reutilização** do Repository genérico
- ✅ **Manutenibilidade** com código organizado
- ✅ **Testabilidade** com interfaces claras

**Resultado:** Infrastructure Layer 100% profissional e escalável! 🚀
