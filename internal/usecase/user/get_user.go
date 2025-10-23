package user

import (
	"context"
	"fmt"

	"github.com/devleo-m/go-zero/internal/domain/user"
	"github.com/google/uuid"
)

// GetUserUseCase implementa o caso de uso de busca de usuário
type GetUserUseCase struct {
	userRepo    user.Repository
	queryHelper *UserQueryHelper
	logger      Logger
}

// NewGetUserUseCase cria uma nova instância do caso de uso
func NewGetUserUseCase(
	userRepo user.Repository,
	queryHelper *UserQueryHelper,
	logger Logger,
) *GetUserUseCase {
	return &GetUserUseCase{
		userRepo:    userRepo,
		queryHelper: queryHelper,
		logger:      logger,
	}
}

// Execute executa o caso de uso de busca de usuário
func (uc *GetUserUseCase) Execute(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	// 1. VALIDAR INPUT
	if err := uc.validateInput(input); err != nil {
		uc.logger.Warn("Invalid input for get user",
			"error", err,
			"input", input,
		)
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	uc.logger.Debug("Starting get user process",
		"input", input,
	)

	var domainUser *user.User
	var err error

	// 2. BUSCAR USUÁRIO POR ID OU EMAIL
	if input.ID != uuid.Nil {
		// Buscar por ID usando o repository genérico
		user, err := uc.userRepo.FindByID(ctx, input.ID)
		if err != nil {
			uc.logger.Error("Failed to find user by ID",
				"error", err,
				"user_id", input.ID,
			)
			return nil, fmt.Errorf("failed to find user by ID: %w", err)
		}
		domainUser = user
	} else if input.Email != "" {
		// Buscar por email usando o helper
		domainUser, err = uc.queryHelper.FindUserByEmail(ctx, input.Email)
		if err != nil {
			uc.logger.Error("Failed to find user by email",
				"error", err,
				"email", input.Email,
			)
			return nil, fmt.Errorf("failed to find user by email: %w", err)
		}
	}

	// 3. VERIFICAR SE USUÁRIO FOI ENCONTRADO
	if domainUser == nil {
		uc.logger.Warn("User not found",
			"input", input,
		)
		return nil, user.ErrUserNotFound
	}

	uc.logger.Info("User found successfully",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
		"name", domainUser.Name,
	)

	// 4. RETORNAR RESULTADO
	output := &GetUserOutput{
		User: ToUserOutput(domainUser),
	}

	uc.logger.Info("Get user use case completed successfully",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
	)

	return output, nil
}

// validateInput valida os dados de entrada
func (uc *GetUserUseCase) validateInput(input GetUserInput) error {
	// Deve ter pelo menos ID ou Email
	if input.ID == uuid.Nil && input.Email == "" {
		return fmt.Errorf("either ID or Email must be provided")
	}

	// Não pode ter ambos
	if input.ID != uuid.Nil && input.Email != "" {
		return fmt.Errorf("cannot provide both ID and Email")
	}

	return nil
}

// ==========================================
// MÉTODOS AUXILIARES
// ==========================================

// GetUserByEmail busca usuário por email
func (uc *GetUserUseCase) GetUserByEmail(ctx context.Context, email string) (*GetUserOutput, error) {
	input := GetUserInput{
		Email: email,
	}
	return uc.Execute(ctx, input)
}

// GetUserByID busca usuário por ID
func (uc *GetUserUseCase) GetUserByID(ctx context.Context, userID uuid.UUID) (*GetUserOutput, error) {
	input := GetUserInput{
		ID: userID,
	}
	return uc.Execute(ctx, input)
}

// CheckUserExists verifica se usuário existe
func (uc *GetUserUseCase) CheckUserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	_, err := uc.GetUserByID(ctx, userID)
	if err != nil {
		if err == user.ErrUserNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetUserBasicInfo retorna informações básicas do usuário
func (uc *GetUserUseCase) GetUserBasicInfo(ctx context.Context, userID uuid.UUID) (*UserBasicInfo, error) {
	userOutput, err := uc.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserBasicInfo{
		ID:     userOutput.User.ID,
		Name:   userOutput.User.Name,
		Email:  userOutput.User.Email,
		Role:   userOutput.User.Role,
		Status: userOutput.User.Status,
	}, nil
}

// GetUserByEmailString busca usuário por email (string)
func (uc *GetUserUseCase) GetUserByEmailString(ctx context.Context, email string) (*GetUserOutput, error) {
	return uc.GetUserByEmail(ctx, email)
}

// GetUserByIDString busca usuário por ID (string)
func (uc *GetUserUseCase) GetUserByIDString(ctx context.Context, userIDStr string) (*GetUserOutput, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	return uc.GetUserByID(ctx, userID)
}

// CheckUserExistsString verifica se usuário existe (string ID)
func (uc *GetUserUseCase) CheckUserExistsString(ctx context.Context, userIDStr string) (bool, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return false, fmt.Errorf("invalid user ID format: %w", err)
	}

	return uc.CheckUserExists(ctx, userID)
}

// GetUserBasicInfoString retorna informações básicas do usuário (string ID)
func (uc *GetUserUseCase) GetUserBasicInfoString(ctx context.Context, userIDStr string) (*UserBasicInfo, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	return uc.GetUserBasicInfo(ctx, userID)
}
