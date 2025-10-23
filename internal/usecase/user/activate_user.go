package user

import (
	"context"
	"fmt"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/domain/user"
	"github.com/google/uuid"
)

// ActivateUserUseCase implementa o caso de uso de ativação de usuário
type ActivateUserUseCase struct {
	userRepo     user.Repository
	emailService EmailService
	logger       Logger
}

// NewActivateUserUseCase cria uma nova instância do caso de uso
func NewActivateUserUseCase(
	userRepo user.Repository,
	emailService EmailService,
	logger Logger,
) *ActivateUserUseCase {
	return &ActivateUserUseCase{
		userRepo:     userRepo,
		emailService: emailService,
		logger:       logger,
	}
}

// Execute executa o caso de uso de ativação de usuário
func (uc *ActivateUserUseCase) Execute(ctx context.Context, input ActivateUserInput) (*ActivateUserOutput, error) {
	// 1. VALIDAR INPUT
	if err := uc.validateInput(input); err != nil {
		uc.logger.Warn("Invalid input for activate user",
			"error", err,
			"user_id", input.UserID,
		)
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	uc.logger.Debug("Starting activate user process",
		"user_id", input.UserID,
	)

	// 2. BUSCAR USUÁRIO
	domainUser, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		uc.logger.Error("Failed to find user for activation",
			"error", err,
			"user_id", input.UserID,
		)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if domainUser == nil {
		uc.logger.Warn("User not found for activation",
			"user_id", input.UserID,
		)
		return nil, user.ErrUserNotFound
	}

	uc.logger.Debug("User found for activation",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
		"current_status", domainUser.Status.String(),
	)

	// 3. VERIFICAR SE O USUÁRIO PODE SER ATIVADO
	if !domainUser.IsPending() {
		uc.logger.Warn("User cannot be activated",
			"user_id", domainUser.ID,
			"current_status", domainUser.Status.String(),
		)
		return nil, fmt.Errorf("user cannot be activated, current status: %s", domainUser.Status.String())
	}

	// 4. ATIVAR USUÁRIO
	if err := domainUser.Activate(); err != nil {
		uc.logger.Error("Failed to activate user",
			"error", err,
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("failed to activate user: %w", err)
	}

	// 5. SALVAR NO BANCO DE DADOS
	if err := uc.userRepo.Update(ctx, domainUser.ID, domainUser); err != nil {
		uc.logger.Error("Failed to save activated user",
			"error", err,
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("failed to save activated user: %w", err)
	}

	uc.logger.Info("User activated successfully",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
	)

	// 6. ENVIAR EMAIL DE CONFIRMAÇÃO DE ATIVAÇÃO (ASSÍNCRONO)
	go func() {
		if err := uc.emailService.SendActivationConfirmationEmail(ctx, domainUser.Email.String(), domainUser.Name); err != nil {
			uc.logger.Error("Failed to send activation confirmation email",
				"error", err,
				"user_id", domainUser.ID,
				"email", domainUser.Email.String(),
			)
		} else {
			uc.logger.Info("Activation confirmation email sent successfully",
				"user_id", domainUser.ID,
				"email", domainUser.Email.String(),
			)
		}
	}()

	// 7. RETORNAR RESULTADO
	output := &ActivateUserOutput{
		User:    ToUserOutput(domainUser),
		Message: "User activated successfully",
	}

	uc.logger.Info("Activate user use case completed successfully",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
	)

	return output, nil
}

// validateInput valida os dados de entrada
func (uc *ActivateUserUseCase) validateInput(input ActivateUserInput) error {
	// Verificar se o ID é válido
	if input.UserID == uuid.Nil {
		return fmt.Errorf("user ID cannot be empty")
	}

	return nil
}

// ==========================================
// MÉTODOS AUXILIARES
// ==========================================

// ActivateUserByEmail ativa usuário por email
func (uc *ActivateUserUseCase) ActivateUserByEmail(ctx context.Context, email string) (*ActivateUserOutput, error) {
	uc.logger.Debug("Activating user by email",
		"email", email,
	)

	// Buscar usuário por email
	domainUser, err := uc.findUserByEmail(ctx, email)
	if err != nil {
		uc.logger.Error("Failed to find user by email for activation",
			"error", err,
			"email", email,
		)
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if domainUser == nil {
		uc.logger.Warn("User not found by email for activation",
			"email", email,
		)
		return nil, user.ErrUserNotFound
	}

	// Usar o método principal de ativação
	input := ActivateUserInput{
		UserID: domainUser.ID,
	}

	return uc.Execute(ctx, input)
}

// ActivateUserWithToken ativa usuário com token de ativação
func (uc *ActivateUserUseCase) ActivateUserWithToken(ctx context.Context, token string) (*ActivateUserOutput, error) {
	uc.logger.Debug("Activating user with token",
		"token", token[:10]+"...", // Log apenas parte do token por segurança
	)

	// TODO: Implementar validação de token de ativação
	// Por enquanto, vamos assumir que o token é válido e contém o user_id

	// Em uma implementação real, você validaria o token e extrairia o user_id
	// Por exemplo:
	// userID, err := uc.tokenService.ValidateActivationToken(token)
	// if err != nil {
	//     return nil, fmt.Errorf("invalid activation token: %w", err)
	// }

	// Por enquanto, vamos retornar erro para indicar que não está implementado
	return nil, fmt.Errorf("token-based activation not implemented yet")
}

// findUserByEmail busca usuário por email usando filtros
func (uc *ActivateUserUseCase) findUserByEmail(ctx context.Context, email string) (*user.User, error) {
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
