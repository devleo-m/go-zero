package application

import (
	"context"
	"fmt"

	"github.com/devleo-m/go-zero/internal/modules/user/domain"
)

// CreateUserUseCase implementa o caso de uso de criação de usuário.
type CreateUserUseCase struct {
	userRepo domain.Repository
}

// NewCreateUserUseCase cria uma nova instância do caso de uso.
func NewCreateUserUseCase(userRepo domain.Repository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo: userRepo,
	}
}

// CreateUserInput representa os dados de entrada.
type CreateUserInput struct {
	Phone    *string `json:"phone,omitempty"`
	Name     string  `json:"name" validate:"required,min=2,max=100"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8"`
}

// CreateUserOutput representa os dados de saída.
type CreateUserOutput struct {
	User    *domain.User `json:"user"`
	Message string       `json:"message"`
}

// Execute executa o caso de uso.
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	// Verificar se email já existe
	existingUser, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil && err != domain.ErrUserNotFound {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	if existingUser != nil {
		return nil, domain.ErrEmailAlreadyInUse
	}

	// Criar usuário
	user, err := domain.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Definir telefone se fornecido
	if input.Phone != nil {
		user.Phone = input.Phone
	}

	// Salvar no banco
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return &CreateUserOutput{
		User:    user,
		Message: "User created successfully",
	}, nil
}
