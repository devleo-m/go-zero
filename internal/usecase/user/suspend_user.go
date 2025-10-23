package user

import (
	"context"
	"fmt"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/domain/user"
	"github.com/google/uuid"
)

// SuspendUserUseCase implementa o caso de uso de suspensão de usuário
type SuspendUserUseCase struct {
	userRepo     user.Repository
	emailService EmailService
	logger       Logger
}

// NewSuspendUserUseCase cria uma nova instância do caso de uso
func NewSuspendUserUseCase(
	userRepo user.Repository,
	emailService EmailService,
	logger Logger,
) *SuspendUserUseCase {
	return &SuspendUserUseCase{
		userRepo:     userRepo,
		emailService: emailService,
		logger:       logger,
	}
}

// Execute executa o caso de uso de suspensão de usuário
func (uc *SuspendUserUseCase) Execute(ctx context.Context, input SuspendUserInput) (*SuspendUserOutput, error) {
	// 1. VALIDAR INPUT
	if err := uc.validateInput(input); err != nil {
		uc.logger.Warn("Invalid input for suspend user",
			"error", err,
			"user_id", input.UserID,
		)
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	uc.logger.Debug("Starting suspend user process",
		"user_id", input.UserID,
	)

	// 2. BUSCAR USUÁRIO
	domainUser, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		uc.logger.Error("Failed to find user for suspension",
			"error", err,
			"user_id", input.UserID,
		)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if domainUser == nil {
		uc.logger.Warn("User not found for suspension",
			"user_id", input.UserID,
		)
		return nil, user.ErrUserNotFound
	}

	uc.logger.Debug("User found for suspension",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
		"current_status", domainUser.Status.String(),
	)

	// 3. VERIFICAR SE O USUÁRIO PODE SER SUSPENSO
	if domainUser.Status == user.StatusSuspended {
		uc.logger.Warn("User is already suspended",
			"user_id", domainUser.ID,
			"current_status", domainUser.Status.String(),
		)
		return nil, fmt.Errorf("user is already suspended")
	}

	// 4. SUSPENDER USUÁRIO
	if err := domainUser.Suspend(); err != nil {
		uc.logger.Error("Failed to suspend user",
			"error", err,
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("failed to suspend user: %w", err)
	}

	// 5. SALVAR NO BANCO DE DADOS
	if err := uc.userRepo.Update(ctx, domainUser.ID, domainUser); err != nil {
		uc.logger.Error("Failed to save suspended user",
			"error", err,
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("failed to save suspended user: %w", err)
	}

	uc.logger.Info("User suspended successfully",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
	)

	// 6. ENVIAR EMAIL DE NOTIFICAÇÃO DE SUSPENSÃO (ASSÍNCRONO)
	go func() {
		if err := uc.emailService.SendSuspensionNotificationEmail(ctx, domainUser.Email.String(), domainUser.Name); err != nil {
			uc.logger.Error("Failed to send suspension notification email",
				"error", err,
				"user_id", domainUser.ID,
				"email", domainUser.Email.String(),
			)
		} else {
			uc.logger.Info("Suspension notification email sent successfully",
				"user_id", domainUser.ID,
				"email", domainUser.Email.String(),
			)
		}
	}()

	// 7. RETORNAR RESULTADO
	output := &SuspendUserOutput{
		User:    ToUserOutput(domainUser),
		Message: "User suspended successfully",
	}

	uc.logger.Info("Suspend user use case completed successfully",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
	)

	return output, nil
}

// validateInput valida os dados de entrada
func (uc *SuspendUserUseCase) validateInput(input SuspendUserInput) error {
	// Verificar se o ID é válido
	if input.UserID == uuid.Nil {
		return fmt.Errorf("user ID cannot be empty")
	}

	return nil
}

// ==========================================
// MÉTODOS AUXILIARES
// ==========================================

// SuspendUserByEmail suspende usuário por email
func (uc *SuspendUserUseCase) SuspendUserByEmail(ctx context.Context, email string) (*SuspendUserOutput, error) {
	uc.logger.Debug("Suspending user by email",
		"email", email,
	)

	// Buscar usuário por email
	domainUser, err := uc.findUserByEmail(ctx, email)
	if err != nil {
		uc.logger.Error("Failed to find user by email for suspension",
			"error", err,
			"email", email,
		)
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if domainUser == nil {
		uc.logger.Warn("User not found by email for suspension",
			"email", email,
		)
		return nil, user.ErrUserNotFound
	}

	// Usar o método principal de suspensão
	input := SuspendUserInput{
		UserID: domainUser.ID,
	}

	return uc.Execute(ctx, input)
}

// findUserByEmail busca usuário por email usando filtros
func (uc *SuspendUserUseCase) findUserByEmail(ctx context.Context, email string) (*user.User, error) {
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
