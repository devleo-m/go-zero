# Repository Genérico Profissional 🎯

## O que é?

Sistema de Repository genérico que substitui interfaces específicas por uma interface única, poderosa e flexível usando Go Generics.

## Por que usar?

### ❌ Repository Específico (Problema):
```go
type UserRepository interface {
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByPhone(ctx context.Context, phone string) (*User, error)
    FindByRole(ctx context.Context, role Role) ([]*User, error)
    FindByStatus(ctx context.Context, status Status) ([]*User, error)
    FindActiveUsers(ctx context.Context) ([]*User, error)
    // ... 50+ métodos específicos
}
```

**Problemas:**
- Interface gigante
- Difícil de manter
- Muito acoplado
- Duplicação de código

### ✅ Repository Genérico (Solução):
```go
type Repository[T any] interface {
    Create(ctx context.Context, entity T) error
    FindOne(ctx context.Context, filter QueryFilter) (T, error)
    FindMany(ctx context.Context, filter QueryFilter) ([]T, error)
    Update(ctx context.Context, id uuid.UUID, entity T) error
    Delete(ctx context.Context, id uuid.UUID) error
    Paginate(ctx context.Context, filter QueryFilter) (*PaginatedResult[T], error)
}
```

**Vantagens:**
- Interface enxuta
- Flexível (busca qualquer campo)
- Fácil de manter
- Reutilizável
- MUITO mais profissional

## Estrutura do Sistema

### 1. Interface Base Genérica
```go
// internal/domain/shared/repository.go
type Repository[T any] interface {
    // CRUD básico
    Create(ctx context.Context, entity T) error
    FindOne(ctx context.Context, filter QueryFilter) (T, error)
    FindMany(ctx context.Context, filter QueryFilter) ([]T, error)
    FindByID(ctx context.Context, id uuid.UUID) (T, error)
    Update(ctx context.Context, id uuid.UUID, entity T) error
    Delete(ctx context.Context, id uuid.UUID) error
    
    // Paginação
    Paginate(ctx context.Context, filter QueryFilter) (*PaginatedResult[T], error)
    
    // Agregações
    Count(ctx context.Context, filter QueryFilter) (int64, error)
    Exists(ctx context.Context, filter QueryFilter) (bool, error)
    Sum(ctx context.Context, field string, filter QueryFilter) (float64, error)
    Avg(ctx context.Context, field string, filter QueryFilter) (float64, error)
    
    // Operações em lote
    CreateMany(ctx context.Context, entities []T) error
    UpdateMany(ctx context.Context, filter QueryFilter, updates map[string]interface{}) (int64, error)
    DeleteMany(ctx context.Context, filter QueryFilter) (int64, error)
    
    // Transações
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
```

### 2. QueryFilter - O Coração do Sistema
```go
// internal/domain/shared/query_filter.go
type QueryFilter struct {
    Where         []Condition `json:"where,omitempty"`
    Or            [][]Condition `json:"or,omitempty"`
    OrderBy       []OrderBy   `json:"order_by,omitempty"`
    Page          int         `json:"page"`
    PageSize      int         `json:"page_size"`
    Include       []string    `json:"include,omitempty"`
    Select        []string    `json:"select,omitempty"`
    Omit          []string    `json:"omit,omitempty"`
    IncludeDeleted bool       `json:"include_deleted"`
    OnlyDeleted   bool        `json:"only_deleted"`
    GroupBy       []string    `json:"group_by,omitempty"`
    Having        []Condition `json:"having,omitempty"`
    Limit         int         `json:"limit,omitempty"`
    Offset        int         `json:"offset,omitempty"`
}

type Condition struct {
    Field         string      `json:"field"`
    Operator      Operator    `json:"operator"`
    Value         interface{} `json:"value"`
    CaseSensitive bool        `json:"case_sensitive"`
}
```

### 3. Operadores Disponíveis
```go
const (
    // Comparação
    OpEqual              Operator = "="
    OpNotEqual           Operator = "!="
    OpGreaterThan        Operator = ">"
    OpGreaterThanOrEqual Operator = ">="
    OpLessThan           Operator = "<"
    OpLessThanOrEqual    Operator = "<="
    
    // String
    OpLike        Operator = "LIKE"
    OpNotLike     Operator = "NOT LIKE"
    OpILike       Operator = "ILIKE"
    OpStartsWith  Operator = "STARTS_WITH"
    OpEndsWith    Operator = "ENDS_WITH"
    OpContains    Operator = "CONTAINS"
    
    // Array
    OpIn    Operator = "IN"
    OpNotIn Operator = "NOT IN"
    
    // Null
    OpIsNull    Operator = "IS NULL"
    OpIsNotNull Operator = "IS NOT NULL"
    
    // Range
    OpBetween    Operator = "BETWEEN"
    OpNotBetween Operator = "NOT BETWEEN"
)
```

### 4. Paginação Profissional
```go
// internal/domain/shared/paginated_result.go
type PaginatedResult[T any] struct {
    Data       []T            `json:"data"`
    Pagination PaginationMeta `json:"pagination"`
    Aggregations map[string]interface{} `json:"aggregations,omitempty"`
    AppliedFilters *QueryFilter `json:"applied_filters,omitempty"`
}

type PaginationMeta struct {
    CurrentPage    int  `json:"current_page"`
    TotalPages     int  `json:"total_pages"`
    PageSize       int  `json:"page_size"`
    TotalItems     int64 `json:"total_items"`
    ItemsInPage    int  `json:"items_in_page"`
    HasPrevious    bool `json:"has_previous"`
    HasNext        bool `json:"has_next"`
    PreviousPage   *int `json:"previous_page,omitempty"`
    NextPage       *int `json:"next_page,omitempty"`
    FirstItemIndex int  `json:"first_item_index"`
    LastItemIndex  int  `json:"last_item_index"`
}
```

## Como Usar

### 1. Definir Repository Específico
```go
// internal/domain/user/repository.go
type Repository interface {
    shared.Repository[*User]
    
    // Métodos específicos APENAS se realmente necessários
    // Tente sempre usar o genérico primeiro!
}
```

### 2. Exemplos de Uso

#### Buscar por Email
```go
// Antes (específico):
user, err := repo.FindByEmail(ctx, "user@example.com")

// Agora (genérico):
user, err := repo.FindOne(ctx, shared.QueryFilter{
    Where: []shared.Condition{
        {
            Field:    "email",
            Operator: shared.OpEqual,
            Value:    "user@example.com",
        },
    },
})
```

#### Buscar Usuários Ativos com Role Admin
```go
users, err := repo.FindMany(ctx, shared.QueryFilter{
    Where: []shared.Condition{
        {
            Field:    "status",
            Operator: shared.OpEqual,
            Value:    "active",
        },
        {
            Field:    "role",
            Operator: shared.OpEqual,
            Value:    "admin",
        },
    },
    OrderBy: []shared.OrderBy{
        {Field: "created_at", Order: shared.SortDesc},
    },
})
```

#### Busca com LIKE (Busca Parcial)
```go
users, err := repo.FindMany(ctx, shared.QueryFilter{
    Where: []shared.Condition{
        {
            Field:    "name",
            Operator: shared.OpILike,
            Value:    "%john%",
        },
    },
})
```

#### Busca com IN (Múltiplos Valores)
```go
users, err := repo.FindMany(ctx, shared.QueryFilter{
    Where: []shared.Condition{
        {
            Field:    "role",
            Operator: shared.OpIn,
            Value:    []string{"admin", "manager"},
        },
    },
})
```

#### Busca com Range (BETWEEN)
```go
users, err := repo.FindMany(ctx, shared.QueryFilter{
    Where: []shared.Condition{
        {
            Field:    "created_at",
            Operator: shared.OpBetween,
            Value:    []time.Time{startDate, endDate},
        },
    },
})
```

#### Paginação Completa
```go
result, err := repo.Paginate(ctx, shared.QueryFilter{
    Where: []shared.Condition{
        {
            Field:    "status",
            Operator: shared.OpEqual,
            Value:    "active",
        },
    },
    OrderBy: []shared.OrderBy{
        {Field: "name", Order: shared.SortAsc},
    },
    Page:     2,
    PageSize: 20,
})

// Resposta:
// {
//   "data": [...],
//   "pagination": {
//     "current_page": 2,
//     "total_pages": 5,
//     "page_size": 20,
//     "total_items": 95,
//     "items_in_page": 20,
//     "has_previous": true,
//     "has_next": true,
//     "previous_page": 1,
//     "next_page": 3,
//     "first_item_index": 21,
//     "last_item_index": 40
//   }
// }
```

#### Contar Usuários
```go
count, err := repo.Count(ctx, shared.QueryFilter{
    Where: []shared.Condition{
        {
            Field:    "status",
            Operator: shared.OpEqual,
            Value:    "active",
        },
    },
})
```

#### Verificar se Existe
```go
exists, err := repo.Exists(ctx, shared.QueryFilter{
    Where: []shared.Condition{
        {
            Field:    "email",
            Operator: shared.OpEqual,
            Value:    "test@example.com",
        },
    },
})
```

#### Update em Lote
```go
affected, err := repo.UpdateMany(ctx, 
    shared.QueryFilter{
        Where: []shared.Condition{
            {
                Field:    "status",
                Operator: shared.OpEqual,
                Value:    "pending",
            },
        },
    },
    map[string]interface{}{
        "status": "active",
    },
)
```

#### Delete em Lote
```go
affected, err := repo.DeleteMany(ctx, shared.QueryFilter{
    Where: []shared.Condition{
        {
            Field:    "created_at",
            Operator: shared.OpLessThan,
            Value:    time.Now().AddDate(0, -6, 0), // Mais de 6 meses
        },
        {
            Field:    "status",
            Operator: shared.OpEqual,
            Value:    "inactive",
        },
    },
})
```

#### Transação
```go
err := repo.WithTransaction(ctx, func(ctx context.Context) error {
    // Criar usuário
    if err := repo.Create(ctx, user1); err != nil {
        return err // Rollback automático
    }
    
    // Criar perfil
    if err := profileRepo.Create(ctx, profile); err != nil {
        return err // Rollback automático
    }
    
    return nil // Commit automático
})
```

#### Agregações
```go
// Somar campo
total, err := repo.Sum(ctx, "login_count", shared.QueryFilter{
    Where: []shared.Condition{
        {
            Field:    "status",
            Operator: shared.OpEqual,
            Value:    "active",
        },
    },
})

// Média
avg, err := repo.Avg(ctx, "login_count", shared.QueryFilter{})

// Mínimo
min, err := repo.Min(ctx, "created_at", shared.QueryFilter{})

// Máximo
max, err := repo.Max(ctx, "created_at", shared.QueryFilter{})
```

## Vantagens do Repository Genérico

### ✅ Interface Enxuta
- **Antes**: 50+ métodos específicos
- **Agora**: 10 métodos poderosos

### ✅ Flexibilidade Total
- Qualquer combinação de filtros
- Qualquer campo, qualquer operador
- Paginação profissional inclusa

### ✅ Fácil Manutenção
- Um lugar para mexer
- Padrão consistente
- Menos bugs

### ✅ Reutilizável
- Serve para TODAS entidades
- User, Patient, Product, etc.
- Mesmo padrão sempre

### ✅ Profissional
- Padrão usado por Prisma, TypeORM, etc.
- Type-safe com Go Generics
- Testável (mockar interface pequena)

## Quando Adicionar Métodos Específicos?

**APENAS** quando:

1. **Query é MUITO complexa** e usada constantemente
2. **Envolve múltiplas tabelas** (joins complexos)
3. **Tem lógica de negócio** dentro da query
4. **Performance crítica** (query otimizada manualmente)

### Exemplo Válido:
```go
type UserRepository interface {
    shared.Repository[*User]
    
    // Query complexa com múltiplos joins e subqueries
    FindUsersWithActiveSubscriptionsAndRecentActivity(
        ctx context.Context,
        days int,
    ) ([]*UserWithSubscription, error)
}
```

## Comparação: Antes vs Agora

### Repository Específico ❌
```go
type UserRepository interface {
    FindByEmail(email string) (*User, error)
    FindByPhone(phone string) (*User, error)
    FindByRole(role Role) ([]*User, error)
    FindActiveUsers() ([]*User, error)
    FindInactiveUsers() ([]*User, error)
    // ... 50+ métodos
}
```

### Repository Genérico ✅
```go
type Repository[T any] interface {
    FindOne(ctx, filter QueryFilter) (T, error)
    FindMany(ctx, filter QueryFilter) ([]T, error)
    Paginate(ctx, filter QueryFilter) (*PaginatedResult[T], error)
    // ... poucos métodos poderosos
}
```

## 4. Query Builder Helpers - Facilita MUITO o uso! 🛠️

```go
// Antes (manual):
filter := QueryFilter{
    Where: []Condition{
        {Field: "status", Operator: OpEqual, Value: "active"},
        {Field: "role", Operator: OpIn, Value: []string{"admin", "user"}},
    },
    OrderBy: []OrderBy{
        {Field: "created_at", Order: SortDesc},
    },
    Page: 1,
    PageSize: 20,
}

// Agora (QueryBuilder):
filter := NewQueryBuilder().
    WhereEqual("status", "active").
    WhereIn("role", []interface{}{"admin", "user"}).
    OrderByDesc("created_at").
    Page(1).
    PageSize(20).
    Build()
```

**Métodos disponíveis:**
- `WhereEqual()`, `WhereNotEqual()`, `WhereIn()`, `WhereLike()`
- `WhereNull()`, `WhereNotNull()`, `WhereBetween()`
- `OrderByAsc()`, `OrderByDesc()`
- `Page()`, `PageSize()`, `Include()`, `Select()`
- `Active()`, `Inactive()`, `CreatedToday()`, `CreatedThisWeek()`

## 5. Specification Pattern - Reutilização de regras! 🎯

```go
// Especificações reutilizáveis
activeUsers := ActiveSpecification[User]()
adminUsers := RoleSpecification[User]("admin")
activeAdmins := activeUsers.And(adminUsers)

// Usar no repository
users, err := repo.FindMany(ctx, activeAdmins.ToQueryFilter())
```

**Especificações disponíveis:**
- `ActiveSpecification[T]()` - Entidades ativas
- `InactiveSpecification[T]()` - Entidades inativas
- `CreatedTodaySpecification[T]()` - Criadas hoje
- `RoleSpecification[T](role)` - Por role específico
- `EmailSpecification[T](email)` - Por email
- `ActiveAdminsSpecification[T]()` - Admins ativos

## 6. Domain Events - Escalabilidade! 🚀

```go
// Criar evento
event := NewUserCreatedEvent(userID, name, email)
user.AddDomainEvent(event)

// Publicar eventos
for _, event := range user.GetDomainEvents() {
    eventBus.Publish(ctx, event)
}
user.ClearDomainEvents()
```

**Eventos disponíveis:**
- `UserCreatedEvent` - Usuário criado
- `UserUpdatedEvent` - Usuário atualizado
- `UserDeletedEvent` - Usuário deletado

## Conclusão

O Repository Genérico + QueryBuilder + Specification + Domain Events é **SUPERIOR** porque:

✅ Interface enxuta (10 métodos vs 50+)  
✅ 100% flexível (qualquer busca possível)  
✅ Fácil de manter  
✅ Padrão da indústria  
✅ Type-safe com Go Generics  
✅ Reutilizável para TODAS entidades  
✅ Paginação profissional inclusa  
✅ Agregações poderosas  
✅ Transações simples  
✅ QueryBuilder facilita uso  
✅ Specification reutiliza regras  
✅ Domain Events para escalabilidade  

**Agora sim está 100% ENTERPRISE-READY!** 🎯🔥
