package application

import (
	"context"
	"fmt"

	"github.com/devleo-m/go-zero/internal/modules/user/domain"
)

// ListUsersUseCase implementa o caso de uso de listar usuários
type ListUsersUseCase struct {
	userRepo domain.Repository
}

// NewListUsersUseCase cria uma nova instância do caso de uso
func NewListUsersUseCase(userRepo domain.Repository) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepo: userRepo,
	}
}

// ListUsersInput representa os dados de entrada
type ListUsersInput struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

// ListUsersOutput representa os dados de saída
type ListUsersOutput struct {
	Users []*domain.User `json:"users"`
	Total int            `json:"total"`
}

// Execute executa o caso de uso
func (uc *ListUsersUseCase) Execute(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	// Definir valores padrão
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Offset < 0 {
		input.Offset = 0
	}

	users, err := uc.userRepo.List(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return &ListUsersOutput{
		Users: users,
		Total: len(users),
	}, nil
}
