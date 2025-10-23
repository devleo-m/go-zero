package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/domain/user"
	"github.com/google/uuid"
)

// UpdateUserUseCase implementa o caso de uso de atualização de usuário
type UpdateUserUseCase struct {
	userRepo user.Repository
	logger   Logger
}

// NewUpdateUserUseCase cria uma nova instância do caso de uso
func NewUpdateUserUseCase(
	userRepo user.Repository,
	logger Logger,
) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

// Execute executa o caso de uso de atualização de usuário
func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error) {
	// 1. VALIDAR INPUT
	if err := uc.validateInput(input); err != nil {
		uc.logger.Warn("Invalid input for update user",
			"error", err,
			"user_id", input.ID,
		)
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	uc.logger.Debug("Starting update user process",
		"user_id", input.ID,
		"name", input.Name,
		"phone", input.Phone,
	)

	// 2. BUSCAR USUÁRIO EXISTENTE
	domainUser, err := uc.userRepo.FindByID(ctx, input.ID)
	if err != nil {
		uc.logger.Error("Failed to find user for update",
			"error", err,
			"user_id", input.ID,
		)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if domainUser == nil {
		uc.logger.Warn("User not found for update",
			"user_id", input.ID,
		)
		return nil, user.ErrUserNotFound
	}

	uc.logger.Debug("User found for update",
		"user_id", domainUser.ID,
		"current_name", domainUser.Name,
		"current_phone", domainUser.Phone,
	)

	// 3. VERIFICAR SE TELEFONE JÁ ESTÁ EM USO (se fornecido)
	if input.Phone != "" {
		existingUser, err := uc.findUserByPhone(ctx, input.Phone)
		if err != nil {
			uc.logger.Error("Failed to check phone availability",
				"error", err,
				"phone", input.Phone,
			)
			return nil, fmt.Errorf("failed to check phone availability: %w", err)
		}

		// Se encontrou um usuário com o telefone e não é o mesmo usuário
		if existingUser != nil && existingUser.ID != domainUser.ID {
			uc.logger.Warn("Phone already in use by another user",
				"phone", input.Phone,
				"current_user_id", domainUser.ID,
				"existing_user_id", existingUser.ID,
			)
			return nil, user.NewPhoneAlreadyInUseError(input.Phone)
		}
	}

	// 4. CRIAR VALUE OBJECTS
	var phoneVO *user.Phone
	if input.Phone != "" {
		phone, err := user.NewPhone(input.Phone)
		if err != nil {
			uc.logger.Warn("Invalid phone format",
				"error", err,
				"phone", input.Phone,
			)
			return nil, fmt.Errorf("invalid phone: %w", err)
		}
		phoneVO = &phone
	}

	// 5. ATUALIZAR PERFIL DO USUÁRIO
	if err := domainUser.UpdateProfile(input.Name, phoneVO); err != nil {
		uc.logger.Error("Failed to update user profile",
			"error", err,
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	// 6. SALVAR NO BANCO DE DADOS
	if err := uc.userRepo.Update(ctx, domainUser.ID, domainUser); err != nil {
		uc.logger.Error("Failed to save updated user",
			"error", err,
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("failed to save updated user: %w", err)
	}

	uc.logger.Info("User updated successfully",
		"user_id", domainUser.ID,
		"name", domainUser.Name,
		"phone", domainUser.Phone,
	)

	// 7. RETORNAR RESULTADO
	output := &UpdateUserOutput{
		User:    ToUserOutput(domainUser),
		Message: "User profile updated successfully",
	}

	uc.logger.Info("Update user use case completed successfully",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
	)

	return output, nil
}

// validateInput valida os dados de entrada
func (uc *UpdateUserUseCase) validateInput(input UpdateUserInput) error {
	// Verificar se o ID é válido
	if input.ID == uuid.Nil {
		return fmt.Errorf("user ID cannot be empty")
	}

	// Verificar se o nome não está vazio
	if input.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	// Verificar se o nome não contém caracteres especiais perigosos
	if containsSpecialChars(input.Name) {
		return fmt.Errorf("name contains invalid characters")
	}

	// Verificar se o telefone não é muito longo (se fornecido)
	if input.Phone != "" && len(input.Phone) > 20 {
		return fmt.Errorf("phone number too long")
	}

	return nil
}

// findUserByPhone busca usuário por telefone (método auxiliar)
func (uc *UpdateUserUseCase) findUserByPhone(ctx context.Context, phone string) (*user.User, error) {
	// Como não temos um método específico no repository para buscar por telefone,
	// vamos usar o método genérico FindMany com filtro
	filter := shared.QueryFilter{
		Where: []shared.Condition{
			{
				Field:    "phone",
				Operator: shared.OpEqual,
				Value:    phone,
			},
		},
		Limit: 1, // Apenas o primeiro resultado
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

// ==========================================
// MÉTODOS AUXILIARES
// ==========================================

// containsSpecialChars verifica se o nome contém caracteres especiais perigosos
func containsSpecialChars(name string) bool {
	// Lista de caracteres que não devem estar no nome
	dangerousChars := []string{"<", ">", "&", "\"", "'", "\\", "/", "=", "+", "(", ")", "[", "]", "{", "}", "|", "`", "~", "!", "@", "#", "$", "%", "^", "*", "?", ":", ";", ",", ".", " "}

	for _, char := range dangerousChars {
		if strings.Contains(name, char) {
			return true
		}
	}
	return false
}

// UpdateUserProfile executa atualização específica de perfil
func (uc *UpdateUserUseCase) UpdateUserProfile(ctx context.Context, input UpdateUserProfileInput) (*UpdateUserProfileOutput, error) {
	// Converter para input genérico
	genericInput := UpdateUserInput{
		ID:    input.UserID,
		Name:  input.Name,
		Phone: input.Phone,
	}

	// Executar caso de uso genérico
	result, err := uc.Execute(ctx, genericInput)
	if err != nil {
		return nil, err
	}

	// Converter para output específico
	return &UpdateUserProfileOutput{
		User:    result.User,
		Message: result.Message,
	}, nil
}
