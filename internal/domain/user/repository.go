package user

import (
	"context"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/google/uuid"
)

// Repository define o contrato para persistência de usuários
// Esta interface fica no DOMAIN, mas a implementação fica na INFRASTRUCTURE
// Isso garante que o domínio não depende de tecnologias específicas (GORM, MongoDB, etc)
type Repository interface {
	// Create cria um novo usuário no banco de dados
	Create(ctx context.Context, user *User) error

	// FindByID busca um usuário pelo ID
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	// FindByEmail busca um usuário pelo email
	FindByEmail(ctx context.Context, email Email) (*User, error)

	// FindByPhone busca um usuário pelo telefone
	FindByPhone(ctx context.Context, phone Phone) (*User, error)

	// Update atualiza um usuário existente
	Update(ctx context.Context, user *User) error

	// Delete remove um usuário (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// List retorna uma lista paginada de usuários
	List(ctx context.Context, filters ListFilters) ([]*User, int64, error)

	// ExistsByEmail verifica se já existe um usuário com o email fornecido
	ExistsByEmail(ctx context.Context, email Email) (bool, error)

	// ExistsByPhone verifica se já existe um usuário com o telefone fornecido
	ExistsByPhone(ctx context.Context, phone Phone) (bool, error)

	// CountByRole conta quantos usuários existem para cada role
	CountByRole(ctx context.Context) (map[Role]int64, error)

	// FindByStatus busca usuários por status
	FindByStatus(ctx context.Context, status Status) ([]*User, error)

	// FindPendingActivation busca usuários pendentes de ativação
	FindPendingActivation(ctx context.Context) ([]*User, error)
}

// ListFilters define filtros para listagem de usuários
type ListFilters struct {
	// Paginação
	Page     int `json:"page"`
	PageSize int `json:"page_size"`

	// Filtros opcionais
	Role   *Role   `json:"role,omitempty"`
	Status *Status `json:"status,omitempty"`

	// Busca por texto
	Search string `json:"search,omitempty"`

	// Ordenação
	SortBy    string `json:"sort_by"`    // "name", "email", "created_at", "updated_at"
	SortOrder string `json:"sort_order"` // "asc", "desc"
}

// DefaultListFilters retorna filtros padrão para listagem
func DefaultListFilters() ListFilters {
	return ListFilters{
		Page:      1,
		PageSize:  20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// Validate valida os filtros de listagem
func (f *ListFilters) Validate() error {
	if f.Page < 1 {
		f.Page = 1
	}

	if f.PageSize < 1 || f.PageSize > 100 {
		f.PageSize = 20
	}

	if f.Role != nil && !f.Role.IsValid() {
		return NewInvalidRoleError(f.Role.String())
	}

	if f.Status != nil && !f.Status.IsValid() {
		return NewInvalidStatusError(f.Status.String())
	}

	validSortFields := []string{"name", "email", "created_at", "updated_at", "last_login_at"}
	if f.SortBy != "" {
		valid := false
		for _, field := range validSortFields {
			if field == f.SortBy {
				valid = true
				break
			}
		}
		if !valid {
			return shared.NewDomainError("INVALID_SORT_FIELD", "invalid sort field: "+f.SortBy)
		}
	}

	if f.SortOrder != "" && f.SortOrder != "asc" && f.SortOrder != "desc" {
		return shared.NewDomainError("INVALID_SORT_ORDER", "sort order must be 'asc' or 'desc'")
	}

	return nil
}

// GetOffset calcula o offset para paginação
func (f *ListFilters) GetOffset() int {
	return (f.Page - 1) * f.PageSize
}
