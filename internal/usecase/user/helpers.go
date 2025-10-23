package user

import (
	"context"
	"fmt"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/domain/user"
)

// UserQueryHelper ajuda com queries comuns usando o repository genérico
type UserQueryHelper struct {
	userRepo user.Repository
}

// NewUserQueryHelper cria uma nova instância do helper
func NewUserQueryHelper(userRepo user.Repository) *UserQueryHelper {
	return &UserQueryHelper{
		userRepo: userRepo,
	}
}

// FindUserByEmail busca usuário por email usando o repository genérico
func (h *UserQueryHelper) FindUserByEmail(ctx context.Context, email string) (*user.User, error) {
	filter := shared.QueryFilter{
		Where: []shared.Condition{
			{
				Field:    "email",
				Operator: shared.OpEqual,
				Value:    email,
			},
		},
		Limit: 1,
	}

	// Usar FindFirst do repository genérico
	domainUser, err := h.userRepo.FindFirst(ctx, filter)
	if err != nil {
		// Se não encontrou, retornar nil (não é erro)
		return nil, nil
	}

	return domainUser, nil
}

// FindUserByPhone busca usuário por telefone usando o repository genérico
func (h *UserQueryHelper) FindUserByPhone(ctx context.Context, phone string) (*user.User, error) {
	filter := shared.QueryFilter{
		Where: []shared.Condition{
			{
				Field:    "phone",
				Operator: shared.OpEqual,
				Value:    phone,
			},
		},
		Limit: 1,
	}

	// Usar FindFirst do repository genérico
	domainUser, err := h.userRepo.FindFirst(ctx, filter)
	if err != nil {
		// Se não encontrou, retornar nil (não é erro)
		return nil, nil
	}

	return domainUser, nil
}

// CheckEmailExists verifica se email já está em uso usando o repository genérico
func (h *UserQueryHelper) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	filter := shared.QueryFilter{
		Where: []shared.Condition{
			{
				Field:    "email",
				Operator: shared.OpEqual,
				Value:    email,
			},
		},
	}

	// Usar Exists do repository genérico
	exists, err := h.userRepo.Exists(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check email exists: %w", err)
	}

	return exists, nil
}

// CheckPhoneExists verifica se telefone já está em uso usando o repository genérico
func (h *UserQueryHelper) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	filter := shared.QueryFilter{
		Where: []shared.Condition{
			{
				Field:    "phone",
				Operator: shared.OpEqual,
				Value:    phone,
			},
		},
	}

	// Usar Exists do repository genérico
	exists, err := h.userRepo.Exists(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check phone exists: %w", err)
	}

	return exists, nil
}

// FindUsersByRole busca usuários por role usando o repository genérico
func (h *UserQueryHelper) FindUsersByRole(ctx context.Context, role string, limit int) ([]*user.User, error) {
	filter := shared.QueryFilter{
		Where: []shared.Condition{
			{
				Field:    "role",
				Operator: shared.OpEqual,
				Value:    role,
			},
		},
		Limit: limit,
	}

	// Usar FindMany do repository genérico
	users, err := h.userRepo.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find users by role: %w", err)
	}

	// Converter para slice de ponteiros
	result := make([]*user.User, len(users))
	copy(result, users)

	return result, nil
}

// FindUsersByStatus busca usuários por status usando o repository genérico
func (h *UserQueryHelper) FindUsersByStatus(ctx context.Context, status string, limit int) ([]*user.User, error) {
	filter := shared.QueryFilter{
		Where: []shared.Condition{
			{
				Field:    "status",
				Operator: shared.OpEqual,
				Value:    status,
			},
		},
		Limit: limit,
	}

	// Usar FindMany do repository genérico
	users, err := h.userRepo.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find users by status: %w", err)
	}

	// Converter para slice de ponteiros
	result := make([]*user.User, len(users))
	copy(result, users)

	return result, nil
}

// SearchUsers busca usuários por termo de busca usando o repository genérico
func (h *UserQueryHelper) SearchUsers(ctx context.Context, searchTerm string, limit int) ([]*user.User, error) {
	searchPattern := "%" + searchTerm + "%"

	filter := shared.QueryFilter{
		Where: []shared.Condition{
			{
				Field:    "name",
				Operator: shared.OpILike,
				Value:    searchPattern,
			},
			{
				Field:    "email",
				Operator: shared.OpILike,
				Value:    searchPattern,
			},
		},
		Limit: limit,
	}

	// Usar FindMany do repository genérico
	users, err := h.userRepo.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	// Converter para slice de ponteiros
	result := make([]*user.User, len(users))
	copy(result, users)

	return result, nil
}

// CountUsersByRole conta usuários por role usando o repository genérico
func (h *UserQueryHelper) CountUsersByRole(ctx context.Context, role string) (int64, error) {
	filter := shared.QueryFilter{
		Where: []shared.Condition{
			{
				Field:    "role",
				Operator: shared.OpEqual,
				Value:    role,
			},
		},
	}

	// Usar Count do repository genérico
	count, err := h.userRepo.Count(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by role: %w", err)
	}

	return count, nil
}

// CountUsersByStatus conta usuários por status usando o repository genérico
func (h *UserQueryHelper) CountUsersByStatus(ctx context.Context, status string) (int64, error) {
	filter := shared.QueryFilter{
		Where: []shared.Condition{
			{
				Field:    "status",
				Operator: shared.OpEqual,
				Value:    status,
			},
		},
	}

	// Usar Count do repository genérico
	count, err := h.userRepo.Count(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by status: %w", err)
	}

	return count, nil
}

// GetUserStats obtém estatísticas dos usuários usando o repository genérico
func (h *UserQueryHelper) GetUserStats(ctx context.Context) (*UserStatsOutput, error) {
	// Contar total de usuários
	totalFilter := shared.QueryFilter{}
	totalCount, err := h.userRepo.Count(ctx, totalFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to count total users: %w", err)
	}

	// Contar usuários ativos
	activeCount, err := h.CountUsersByStatus(ctx, user.StatusActive.String())
	if err != nil {
		return nil, fmt.Errorf("failed to count active users: %w", err)
	}

	// Contar usuários pendentes
	pendingCount, err := h.CountUsersByStatus(ctx, user.StatusPending.String())
	if err != nil {
		return nil, fmt.Errorf("failed to count pending users: %w", err)
	}

	// Contar usuários suspensos
	suspendedCount, err := h.CountUsersByStatus(ctx, user.StatusSuspended.String())
	if err != nil {
		return nil, fmt.Errorf("failed to count suspended users: %w", err)
	}

	// Contar usuários inativos
	inactiveCount, err := h.CountUsersByStatus(ctx, user.StatusInactive.String())
	if err != nil {
		return nil, fmt.Errorf("failed to count inactive users: %w", err)
	}

	return &UserStatsOutput{
		Total:     totalCount,
		Active:    activeCount,
		Pending:   pendingCount,
		Suspended: suspendedCount,
		Inactive:  inactiveCount,
	}, nil
}
