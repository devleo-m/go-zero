package user

import (
	"context"
	"fmt"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/domain/user"
	"github.com/google/uuid"
)

// ChangeRoleUseCase implementa o caso de uso de alteração de role
type ChangeRoleUseCase struct {
	userRepo     user.Repository
	emailService EmailService
	logger       Logger
}

// NewChangeRoleUseCase cria uma nova instância do caso de uso
func NewChangeRoleUseCase(
	userRepo user.Repository,
	emailService EmailService,
	logger Logger,
) *ChangeRoleUseCase {
	return &ChangeRoleUseCase{
		userRepo:     userRepo,
		emailService: emailService,
		logger:       logger,
	}
}

// Execute executa o caso de uso de alteração de role
func (uc *ChangeRoleUseCase) Execute(ctx context.Context, input ChangeRoleInput) (*ChangeRoleOutput, error) {
	// 1. VALIDAR INPUT
	if err := uc.validateInput(input); err != nil {
		uc.logger.Warn("Invalid input for change role",
			"error", err,
			"user_id", input.UserID,
			"new_role", input.NewRole,
			"requester_id", input.RequesterID,
		)
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	uc.logger.Debug("Starting change role process",
		"user_id", input.UserID,
		"new_role", input.NewRole,
		"requester_id", input.RequesterID,
	)

	// 2. BUSCAR USUÁRIO ALVO
	targetUser, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		uc.logger.Error("Failed to find target user for role change",
			"error", err,
			"user_id", input.UserID,
		)
		return nil, fmt.Errorf("failed to find target user: %w", err)
	}

	if targetUser == nil {
		uc.logger.Warn("Target user not found for role change",
			"user_id", input.UserID,
		)
		return nil, user.ErrUserNotFound
	}

	// 3. BUSCAR USUÁRIO SOLICITANTE
	requesterUser, err := uc.userRepo.FindByID(ctx, input.RequesterID)
	if err != nil {
		uc.logger.Error("Failed to find requester user for role change",
			"error", err,
			"requester_id", input.RequesterID,
		)
		return nil, fmt.Errorf("failed to find requester user: %w", err)
	}

	if requesterUser == nil {
		uc.logger.Warn("Requester user not found for role change",
			"requester_id", input.RequesterID,
		)
		return nil, user.ErrUserNotFound
	}

	uc.logger.Debug("Users found for role change",
		"target_user_id", targetUser.ID,
		"target_current_role", targetUser.Role.String(),
		"requester_id", requesterUser.ID,
		"requester_role", requesterUser.Role.String(),
		"new_role", input.NewRole,
	)

	// 4. VERIFICAR SE O SOLICITANTE PODE ALTERAR O ROLE
	newRole := user.Role(input.NewRole)
	if !newRole.IsValid() {
		uc.logger.Warn("Invalid new role",
			"new_role", input.NewRole,
		)
		return nil, user.NewInvalidRoleError(input.NewRole)
	}

	// 5. ALTERAR ROLE DO USUÁRIO
	if err := targetUser.ChangeRole(newRole, requesterUser.Role); err != nil {
		uc.logger.Warn("Failed to change user role",
			"error", err,
			"target_user_id", targetUser.ID,
			"new_role", input.NewRole,
			"requester_role", requesterUser.Role.String(),
		)
		return nil, err
	}

	// 6. SALVAR NO BANCO DE DADOS
	if err := uc.userRepo.Update(ctx, targetUser.ID, targetUser); err != nil {
		uc.logger.Error("Failed to save user with new role",
			"error", err,
			"user_id", targetUser.ID,
		)
		return nil, fmt.Errorf("failed to save user with new role: %w", err)
	}

	uc.logger.Info("User role changed successfully",
		"user_id", targetUser.ID,
		"email", targetUser.Email.String(),
		"old_role", targetUser.Role.String(),
		"new_role", newRole.String(),
		"requester_id", requesterUser.ID,
	)

	// 7. ENVIAR EMAIL DE NOTIFICAÇÃO DE MUDANÇA DE ROLE (ASSÍNCRONO)
	go func() {
		if err := uc.emailService.SendRoleChangeNotificationEmail(ctx, targetUser.Email.String(), targetUser.Name, newRole.String()); err != nil {
			uc.logger.Error("Failed to send role change notification email",
				"error", err,
				"user_id", targetUser.ID,
				"email", targetUser.Email.String(),
			)
		} else {
			uc.logger.Info("Role change notification email sent successfully",
				"user_id", targetUser.ID,
				"email", targetUser.Email.String(),
			)
		}
	}()

	// 8. RETORNAR RESULTADO
	output := &ChangeRoleOutput{
		User:    ToUserOutput(targetUser),
		Message: fmt.Sprintf("User role changed to %s successfully", newRole.String()),
	}

	uc.logger.Info("Change role use case completed successfully",
		"user_id", targetUser.ID,
		"email", targetUser.Email.String(),
		"new_role", newRole.String(),
	)

	return output, nil
}

// validateInput valida os dados de entrada
func (uc *ChangeRoleUseCase) validateInput(input ChangeRoleInput) error {
	// Verificar se o ID do usuário alvo é válido
	if input.UserID == uuid.Nil {
		return fmt.Errorf("target user ID cannot be empty")
	}

	// Verificar se o ID do solicitante é válido
	if input.RequesterID == uuid.Nil {
		return fmt.Errorf("requester user ID cannot be empty")
	}

	// Verificar se não está tentando alterar o próprio role
	if input.UserID == input.RequesterID {
		uc.logger.Warn("User trying to change their own role",
			"user_id", input.UserID,
		)
		return fmt.Errorf("cannot change your own role")
	}

	// Verificar se o novo role não está vazio
	if input.NewRole == "" {
		return fmt.Errorf("new role cannot be empty")
	}

	return nil
}

// ==========================================
// MÉTODOS AUXILIARES
// ==========================================

// ChangeRoleByEmail altera role por email
func (uc *ChangeRoleUseCase) ChangeRoleByEmail(ctx context.Context, email, newRole string, requesterID uuid.UUID) (*ChangeRoleOutput, error) {
	uc.logger.Debug("Changing role by email",
		"email", email,
		"new_role", newRole,
		"requester_id", requesterID,
	)

	// Buscar usuário por email
	domainUser, err := uc.findUserByEmail(ctx, email)
	if err != nil {
		uc.logger.Error("Failed to find user by email for role change",
			"error", err,
			"email", email,
		)
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if domainUser == nil {
		uc.logger.Warn("User not found by email for role change",
			"email", email,
		)
		return nil, user.ErrUserNotFound
	}

	// Usar o método principal de alteração de role
	input := ChangeRoleInput{
		UserID:      domainUser.ID,
		NewRole:     newRole,
		RequesterID: requesterID,
	}

	return uc.Execute(ctx, input)
}

// findUserByEmail busca usuário por email usando filtros
func (uc *ChangeRoleUseCase) findUserByEmail(ctx context.Context, email string) (*user.User, error) {
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
