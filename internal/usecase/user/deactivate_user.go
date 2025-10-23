package user

import (
	"context"
	"fmt"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/domain/user"
	"github.com/google/uuid"
)

// DeactivateUserUseCase implementa o caso de uso de desativação de usuário
type DeactivateUserUseCase struct {
	userRepo     user.Repository
	emailService EmailService
	logger       Logger
}

// NewDeactivateUserUseCase cria uma nova instância do caso de uso
func NewDeactivateUserUseCase(
	userRepo user.Repository,
	emailService EmailService,
	logger Logger,
) *DeactivateUserUseCase {
	return &DeactivateUserUseCase{
		userRepo:     userRepo,
		emailService: emailService,
		logger:       logger,
	}
}

// Execute executa o caso de uso de desativação de usuário
func (uc *DeactivateUserUseCase) Execute(ctx context.Context, input DeactivateUserInput) (*DeactivateUserOutput, error) {
	// 1. VALIDAR INPUT
	if err := uc.validateInput(input); err != nil {
		uc.logger.Warn("Invalid input for deactivate user",
			"error", err,
			"user_id", input.UserID,
		)
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	uc.logger.Debug("Starting deactivate user process",
		"user_id", input.UserID,
	)

	// 2. BUSCAR USUÁRIO
	domainUser, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		uc.logger.Error("Failed to find user for deactivation",
			"error", err,
			"user_id", input.UserID,
		)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if domainUser == nil {
		uc.logger.Warn("User not found for deactivation",
			"user_id", input.UserID,
		)
		return nil, user.ErrUserNotFound
	}

	uc.logger.Debug("User found for deactivation",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
		"current_status", domainUser.Status.String(),
	)

	// 3. VERIFICAR SE O USUÁRIO PODE SER DESATIVADO
	if domainUser.Status == user.StatusInactive {
		uc.logger.Warn("User is already inactive",
			"user_id", domainUser.ID,
			"current_status", domainUser.Status.String(),
		)
		return nil, fmt.Errorf("user is already inactive")
	}

	// 4. DESATIVAR USUÁRIO
	if err := domainUser.Deactivate(); err != nil {
		uc.logger.Error("Failed to deactivate user",
			"error", err,
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("failed to deactivate user: %w", err)
	}

	// 5. SALVAR NO BANCO DE DADOS
	if err := uc.userRepo.Update(ctx, domainUser.ID, domainUser); err != nil {
		uc.logger.Error("Failed to save deactivated user",
			"error", err,
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("failed to save deactivated user: %w", err)
	}

	uc.logger.Info("User deactivated successfully",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
	)

	// 6. ENVIAR EMAIL DE NOTIFICAÇÃO DE DESATIVAÇÃO (ASSÍNCRONO)
	go func() {
		if err := uc.emailService.SendDeactivationNotificationEmail(ctx, domainUser.Email.String(), domainUser.Name); err != nil {
			uc.logger.Error("Failed to send deactivation notification email",
				"error", err,
				"user_id", domainUser.ID,
				"email", domainUser.Email.String(),
			)
		} else {
			uc.logger.Info("Deactivation notification email sent successfully",
				"user_id", domainUser.ID,
				"email", domainUser.Email.String(),
			)
		}
	}()

	// 7. RETORNAR RESULTADO
	output := &DeactivateUserOutput{
		User:    ToUserOutput(domainUser),
		Message: "User deactivated successfully",
	}

	uc.logger.Info("Deactivate user use case completed successfully",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
	)

	return output, nil
}

// validateInput valida os dados de entrada
func (uc *DeactivateUserUseCase) validateInput(input DeactivateUserInput) error {
	// Verificar se o ID é válido
	if input.UserID == uuid.Nil {
		return fmt.Errorf("user ID cannot be empty")
	}

	return nil
}

// ==========================================
// MÉTODOS AUXILIARES
// ==========================================

// DeactivateUserByEmail desativa usuário por email
func (uc *DeactivateUserUseCase) DeactivateUserByEmail(ctx context.Context, email string) (*DeactivateUserOutput, error) {
	uc.logger.Debug("Deactivating user by email",
		"email", email,
	)

	// Buscar usuário por email
	domainUser, err := uc.findUserByEmail(ctx, email)
	if err != nil {
		uc.logger.Error("Failed to find user by email for deactivation",
			"error", err,
			"email", email,
		)
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if domainUser == nil {
		uc.logger.Warn("User not found by email for deactivation",
			"email", email,
		)
		return nil, user.ErrUserNotFound
	}

	// Usar o método principal de desativação
	input := DeactivateUserInput{
		UserID: domainUser.ID,
	}

	return uc.Execute(ctx, input)
}

// findUserByEmail busca usuário por email usando filtros
func (uc *DeactivateUserUseCase) findUserByEmail(ctx context.Context, email string) (*user.User, error) {
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

	users, err := uc.userRepo.FindMany(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	return users[0], nil
}
