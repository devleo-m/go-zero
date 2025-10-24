package application

import (
	"context"
	"fmt"

	"github.com/devleo-m/go-zero/internal/modules/user/domain"
	"github.com/google/uuid"
)

// UpdateUserUseCase implementa o caso de uso de atualizar usuário
type UpdateUserUseCase struct {
	userRepo domain.Repository
}

// NewUpdateUserUseCase cria uma nova instância do caso de uso
func NewUpdateUserUseCase(userRepo domain.Repository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepo: userRepo,
	}
}

// UpdateUserInput representa os dados de entrada
type UpdateUserInput struct {
	ID    uuid.UUID `json:"id" validate:"required"`
	Name  string    `json:"name" validate:"required,min=2,max=100"`
	Phone string    `json:"phone,omitempty"`
}

// UpdateUserOutput representa os dados de saída
type UpdateUserOutput struct {
	User    *domain.User `json:"user"`
	Message string       `json:"message"`
}

// Execute executa o caso de uso
func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error) {
	// Buscar usuário existente
	user, err := uc.userRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Atualizar perfil
	var phone *string
	if input.Phone != "" {
		phone = &input.Phone
	}

	if err := user.UpdateProfile(input.Name, phone); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Salvar no banco
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return &UpdateUserOutput{
		User:    user,
		Message: "User updated successfully",
	}, nil
}
