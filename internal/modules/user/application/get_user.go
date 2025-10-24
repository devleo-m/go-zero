package application

import (
	"context"
	"fmt"

	"github.com/devleo-m/go-zero/internal/modules/user/domain"
	"github.com/google/uuid"
)

// GetUserUseCase implementa o caso de uso de buscar usuário
type GetUserUseCase struct {
	userRepo domain.Repository
}

// NewGetUserUseCase cria uma nova instância do caso de uso
func NewGetUserUseCase(userRepo domain.Repository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepo: userRepo,
	}
}

// GetUserInput representa os dados de entrada
type GetUserInput struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// GetUserOutput representa os dados de saída
type GetUserOutput struct {
	User *domain.User `json:"user"`
}

// Execute executa o caso de uso
func (uc *GetUserUseCase) Execute(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	user, err := uc.userRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &GetUserOutput{
		User: user,
	}, nil
}
