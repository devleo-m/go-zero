// internal/domain/shared/usage_examples.go
package shared

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ==========================================
// EXEMPLOS DE USO DO REPOSITORY GENÉRICO
// ==========================================

// ExemploUserRepository demonstra como usar o Repository genérico
type ExemploUserRepository struct{}

// ==========================================
// 1. QUERY BUILDER - FACILITA MUITO!
// ==========================================

func (r *ExemploUserRepository) BuscarUsuariosAtivos() QueryFilter {
	// Antes (manual):
	// filter := QueryFilter{
	//     Where: []Condition{
	//         {Field: "status", Operator: OpEqual, Value: "active"},
	//     },
	//     OrderBy: []OrderBy{
	//         {Field: "created_at", Order: SortDesc},
	//     },
	//     Page: 1,
	//     PageSize: 20,
	// }

	// Agora (QueryBuilder):
	return NewQueryBuilder().
		WhereEqual("status", "active").
		OrderByDesc("created_at").
		Page(1).
		PageSize(20).
		Build()
}

func (r *ExemploUserRepository) BuscarUsuariosPorNome(nome string) QueryFilter {
	return NewQueryBuilder().
		WhereILike("name", nome).
		WhereEqual("status", "active").
		OrderByAsc("name").
		Build()
}

func (r *ExemploUserRepository) BuscarUsuariosCriadosHoje() QueryFilter {
	return NewQueryBuilder().
		CreatedToday().
		OrderByDesc("created_at").
		Build()
}

func (r *ExemploUserRepository) BuscarUsuariosComRoles(roles []string) QueryFilter {
	return NewQueryBuilder().
		WhereIn("role", convertToStringSlice(roles)).
		WhereEqual("status", "active").
		OrderByDesc("created_at").
		Build()
}

// ==========================================
// 2. SPECIFICATION PATTERN - REUTILIZAÇÃO!
// ==========================================

func (r *ExemploUserRepository) BuscarAdminsAtivos() QueryFilter {
	// Usando especificações reutilizáveis
	activeAdmins := ActiveAdminsSpecification[User]()
	return activeAdmins.ToQueryFilter()
}

func (r *ExemploUserRepository) BuscarUsuariosPendentes() QueryFilter {
	// Usando especificação específica
	pending := PendingActivationSpecification[User]()
	return pending.ToQueryFilter()
}

func (r *ExemploUserRepository) BuscarUsuariosCriadosEstaSemana() QueryFilter {
	// Usando especificação de tempo
	thisWeek := CreatedThisWeekSpecification[User]()
	return thisWeek.ToQueryFilter()
}

func (r *ExemploUserRepository) BuscarUsuariosComplexos() QueryFilter {
	// Combinando especificações
	activeUsers := ActiveSpecification[User]()
	adminRole := RoleSpecification[User]("admin")
	userRole := RoleSpecification[User]("user")

	// Admins OU usuários ativos
	activeAdminsOrUsers := activeUsers.And(adminRole.Or(userRole))

	return activeAdminsOrUsers.ToQueryFilter()
}

// ==========================================
// 3. DOMAIN EVENTS - ESCALABILIDADE!
// ==========================================

// UserCreatedEvent evento de usuário criado
type UserCreatedEvent struct {
	BaseDomainEvent
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	UserEmail string    `json:"user_email"`
}

func NewUserCreatedEvent(userID uuid.UUID, name, email string) DomainEvent {
	return &UserCreatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			eventType:   "user.created",
			occurredAt:  time.Now(),
			aggregateID: userID,
			eventData:   map[string]interface{}{"name": name, "email": email},
		},
		UserID:    userID,
		UserName:  name,
		UserEmail: email,
	}
}

// UserUpdatedEvent evento de usuário atualizado
type UserUpdatedEvent struct {
	BaseDomainEvent
	UserID    uuid.UUID              `json:"user_id"`
	Changes   map[string]interface{} `json:"changes"`
	UpdatedBy uuid.UUID              `json:"updated_by"`
}

func NewUserUpdatedEvent(userID uuid.UUID, changes map[string]interface{}, updatedBy uuid.UUID) DomainEvent {
	return &UserUpdatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			eventType:   "user.updated",
			occurredAt:  time.Now(),
			aggregateID: userID,
			eventData:   changes,
		},
		UserID:    userID,
		Changes:   changes,
		UpdatedBy: updatedBy,
	}
}

// UserDeletedEvent evento de usuário deletado
type UserDeletedEvent struct {
	BaseDomainEvent
	UserID    uuid.UUID `json:"user_id"`
	DeletedBy uuid.UUID `json:"deleted_by"`
}

func NewUserDeletedEvent(userID, deletedBy uuid.UUID) DomainEvent {
	return &UserDeletedEvent{
		BaseDomainEvent: BaseDomainEvent{
			eventType:   "user.deleted",
			occurredAt:  time.Now(),
			aggregateID: userID,
			eventData:   map[string]interface{}{"deleted_by": deletedBy},
		},
		UserID:    userID,
		DeletedBy: deletedBy,
	}
}

// ==========================================
// 4. EXEMPLO DE USO COMPLETO
// ==========================================

func (r *ExemploUserRepository) ExemploCompleto(ctx context.Context) {
	// 1. Buscar usuários ativos com QueryBuilder
	activeUsersFilter := NewQueryBuilder().
		WhereEqual("status", "active").
		WhereIn("role", []interface{}{"admin", "user"}).
		OrderByDesc("created_at").
		Page(1).
		PageSize(20).
		Include("profile", "permissions").
		Build()

	// 2. Buscar admins com Specification
	adminFilter := ActiveAdminsSpecification[User]().ToQueryFilter()

	// 3. Buscar usuários criados hoje
	todayFilter := CreatedTodaySpecification[User]().ToQueryFilter()

	// 4. Combinar especificações
	complexFilter := ActiveSpecification[User]().
		And(RoleSpecification[User]("admin").
			Or(RoleSpecification[User]("manager"))).
		ToQueryFilter()

	// 5. Usar em repository (exemplo)
	// users, err := userRepo.FindMany(ctx, activeUsersFilter)
	// admins, err := userRepo.FindMany(ctx, adminFilter)
	// todayUsers, err := userRepo.FindMany(ctx, todayFilter)
	// complexUsers, err := userRepo.FindMany(ctx, complexFilter)

	_ = activeUsersFilter
	_ = adminFilter
	_ = todayFilter
	_ = complexFilter
}

// ==========================================
// 5. UTILITÁRIOS
// ==========================================

func convertToStringSlice(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

// User representa uma entidade de exemplo
type User struct {
	BaseEntity
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// ==========================================
// 6. COMPARAÇÃO: ANTES vs AGORA
// ==========================================

func ComparacaoAntesDepois() {
	// ❌ ANTES (Repository específico):
	// users, err := repo.FindByStatus(ctx, "active")
	// users, err := repo.FindByRole(ctx, "admin")
	// users, err := repo.FindByCreatedAtBetween(ctx, start, end)
	// users, err := repo.FindByNameLike(ctx, "%john%")
	// users, err := repo.FindByRoleIn(ctx, []string{"admin", "user"})

	// ✅ AGORA (Repository genérico + QueryBuilder + Specification):

	// QueryBuilder - Fácil e legível
	filter1 := NewQueryBuilder().
		WhereEqual("status", "active").
		WhereIn("role", []interface{}{"admin", "user"}).
		OrderByDesc("created_at").
		Page(1).
		PageSize(20).
		Build()

	// Specification - Reutilizável
	filter2 := ActiveAdminsSpecification[User]().ToQueryFilter()

	// Combinado - Poderoso
	filter3 := ActiveSpecification[User]().
		And(RoleSpecification[User]("admin").
			Or(RoleSpecification[User]("manager"))).
		ToQueryFilter()

	_ = filter1
	_ = filter2
	_ = filter3
}
