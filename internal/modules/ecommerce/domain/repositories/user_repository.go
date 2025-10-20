package repositories

import (
	"context"

	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/entities"
	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/valueobjects"
)

// UserRepository define o contrato para operações de usuário
type UserRepository interface {
	// Create cria um novo usuário
	Create(ctx context.Context, user *entities.User) error

	// GetByID busca um usuário pelo ID
	GetByID(ctx context.Context, id string) (*entities.User, error)

	// GetByEmail busca um usuário pelo email
	GetByEmail(ctx context.Context, email valueobjects.Email) (*entities.User, error)

	// Update atualiza um usuário existente
	Update(ctx context.Context, user *entities.User) error

	// Delete remove um usuário (soft delete)
	Delete(ctx context.Context, id string) error

	// List retorna uma lista de usuários com paginação
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)

	// Count retorna o total de usuários
	Count(ctx context.Context) (int64, error)

	// Exists verifica se um usuário existe
	Exists(ctx context.Context, id string) (bool, error)

	// ExistsByEmail verifica se um usuário existe pelo email
	ExistsByEmail(ctx context.Context, email valueobjects.Email) (bool, error)

	// ListByRole retorna usuários por papel
	ListByRole(ctx context.Context, role entities.UserRole, limit, offset int) ([]*entities.User, error)

	// ListByStatus retorna usuários por status
	ListByStatus(ctx context.Context, status entities.UserStatus, limit, offset int) ([]*entities.User, error)

	// Search busca usuários por critérios
	Search(ctx context.Context, query string, limit, offset int) ([]*entities.User, error)

	// GetActiveUsers retorna usuários ativos
	GetActiveUsers(ctx context.Context, limit, offset int) ([]*entities.User, error)

	// GetVerifiedUsers retorna usuários verificados
	GetVerifiedUsers(ctx context.Context, limit, offset int) ([]*entities.User, error)

	// GetUsersCreatedBetween retorna usuários criados em um período
	GetUsersCreatedBetween(ctx context.Context, start, end string, limit, offset int) ([]*entities.User, error)

	// GetUsersByLastLogin retorna usuários por último login
	GetUsersByLastLogin(ctx context.Context, limit, offset int) ([]*entities.User, error)

	// BulkUpdate atualiza múltiplos usuários
	BulkUpdate(ctx context.Context, users []*entities.User) error

	// BulkDelete remove múltiplos usuários
	BulkDelete(ctx context.Context, ids []string) error

	// Restore restaura um usuário deletado
	Restore(ctx context.Context, id string) error

	// GetDeletedUsers retorna usuários deletados
	GetDeletedUsers(ctx context.Context, limit, offset int) ([]*entities.User, error)

	// HardDelete remove permanentemente um usuário
	HardDelete(ctx context.Context, id string) error

	// GetUserStats retorna estatísticas dos usuários
	GetUserStats(ctx context.Context) (*UserStats, error)
}

// UserStats representa estatísticas dos usuários
type UserStats struct {
	TotalUsers      int64 `json:"total_users"`
	ActiveUsers     int64 `json:"active_users"`
	InactiveUsers   int64 `json:"inactive_users"`
	BlockedUsers    int64 `json:"blocked_users"`
	VerifiedUsers   int64 `json:"verified_users"`
	UnverifiedUsers int64 `json:"unverified_users"`
	AdminUsers      int64 `json:"admin_users"`
	ClientUsers     int64 `json:"client_users"`
	DeletedUsers    int64 `json:"deleted_users"`
}
