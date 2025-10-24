package application

import (
	"context"
	"fmt"

	"github.com/devleo-m/go-zero/internal/modules/user/domain"
	"github.com/google/uuid"
)

// DeleteUserUseCase implementa o caso de uso de deletar usuário
type DeleteUserUseCase struct {
	userRepo domain.Repository
}

// NewDeleteUserUseCase cria uma nova instância do caso de uso
func NewDeleteUserUseCase(userRepo domain.Repository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepo: userRepo,
	}
}

// DeleteUserInput representa os dados de entrada
type DeleteUserInput struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// DeleteUserOutput representa os dados de saída
type DeleteUserOutput struct {
	Message string `json:"message"`
}

// Execute executa o caso de uso
func (uc *DeleteUserUseCase) Execute(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, error) {
	// Verificar se usuário existe
	_, err := uc.userRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Deletar usuário (soft delete)
	if err := uc.userRepo.Delete(ctx, input.ID); err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &DeleteUserOutput{
		Message: "User deleted successfully",
	}, nil
}
