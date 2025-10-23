package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/devleo-m/go-zero/internal/domain/user"
	"github.com/google/uuid"
)

// ChangePasswordUseCase implementa o caso de uso de alteração de senha
type ChangePasswordUseCase struct {
	userRepo user.Repository
	logger   Logger
}

// NewChangePasswordUseCase cria uma nova instância do caso de uso
func NewChangePasswordUseCase(
	userRepo user.Repository,
	logger Logger,
) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

// Execute executa o caso de uso de alteração de senha
func (uc *ChangePasswordUseCase) Execute(ctx context.Context, input ChangePasswordInput) (*ChangePasswordOutput, error) {
	// 1. VALIDAR INPUT
	if err := uc.validateInput(input); err != nil {
		uc.logger.Warn("Invalid input for change password",
			"error", err,
			"user_id", input.UserID,
		)
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	uc.logger.Debug("Starting change password process",
		"user_id", input.UserID,
	)

	// 2. BUSCAR USUÁRIO
	domainUser, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		uc.logger.Error("Failed to find user for password change",
			"error", err,
			"user_id", input.UserID,
		)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if domainUser == nil {
		uc.logger.Warn("User not found for password change",
			"user_id", input.UserID,
		)
		return nil, user.ErrUserNotFound
	}

	uc.logger.Debug("User found for password change",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
	)

	// 3. VERIFICAR SE A NOVA SENHA É DIFERENTE DA ATUAL
	if domainUser.Password.Verify(input.NewPassword) {
		uc.logger.Warn("New password is the same as current password",
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("new password must be different from current password")
	}

	// 4. VERIFICAR SE A NOVA SENHA NÃO É MUITO COMUM
	if isCommonPassword(input.NewPassword) {
		uc.logger.Warn("New password is too common",
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("password is too common, please choose a stronger password")
	}

	// 5. VERIFICAR SE A NOVA SENHA NÃO É MUITO SIMILAR À ATUAL
	if uc.isPasswordTooSimilar(input.OldPassword, input.NewPassword) {
		uc.logger.Warn("New password is too similar to current password",
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("new password is too similar to current password")
	}

	// 6. ALTERAR SENHA (isso já valida a senha atual)
	if err := domainUser.ChangePassword(input.OldPassword, input.NewPassword); err != nil {
		uc.logger.Warn("Failed to change password",
			"error", err,
			"user_id", domainUser.ID,
		)
		return nil, err
	}

	// 7. SALVAR NO BANCO DE DADOS
	if err := uc.userRepo.Update(ctx, domainUser.ID, domainUser); err != nil {
		uc.logger.Error("Failed to save user with new password",
			"error", err,
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("failed to save user with new password: %w", err)
	}

	uc.logger.Info("Password changed successfully",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
	)

	// 8. RETORNAR RESULTADO
	output := &ChangePasswordOutput{
		Message: "Password changed successfully",
	}

	uc.logger.Info("Change password use case completed successfully",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
	)

	return output, nil
}

// validateInput valida os dados de entrada
func (uc *ChangePasswordUseCase) validateInput(input ChangePasswordInput) error {
	// Verificar se o ID é válido
	if input.UserID == uuid.Nil {
		return fmt.Errorf("user ID cannot be empty")
	}

	// Verificar se a senha atual não está vazia
	if input.OldPassword == "" {
		return fmt.Errorf("current password cannot be empty")
	}

	// Verificar se a nova senha não está vazia
	if input.NewPassword == "" {
		return fmt.Errorf("new password cannot be empty")
	}

	// Verificar se as senhas são diferentes
	if input.OldPassword == input.NewPassword {
		return fmt.Errorf("new password must be different from current password")
	}

	// Verificar se a nova senha não é muito longa (possível ataque)
	if len(input.NewPassword) > 128 {
		return fmt.Errorf("new password too long")
	}

	return nil
}

// isPasswordTooSimilar verifica se a nova senha é muito similar à atual
func (uc *ChangePasswordUseCase) isPasswordTooSimilar(oldPassword, newPassword string) bool {
	// Verificar se a nova senha contém a senha atual
	if len(newPassword) >= len(oldPassword) &&
		(len(newPassword) == len(oldPassword) ||
			strings.Contains(newPassword, oldPassword)) {
		return true
	}

	// Verificar se a senha atual contém a nova senha
	if len(oldPassword) >= len(newPassword) &&
		(len(oldPassword) == len(newPassword) ||
			strings.Contains(oldPassword, newPassword)) {
		return true
	}

	// Verificar similaridade por caracteres comuns (mais de 80% de similaridade)
	similarity := uc.calculateSimilarity(oldPassword, newPassword)
	return similarity > 0.8
}

// calculateSimilarity calcula a similaridade entre duas strings
func (uc *ChangePasswordUseCase) calculateSimilarity(s1, s2 string) float64 {
	if len(s1) == 0 && len(s2) == 0 {
		return 1.0
	}
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	// Usar algoritmo de Levenshtein para calcular similaridade
	distance := uc.levenshteinDistance(s1, s2)
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}

	return 1.0 - float64(distance)/float64(maxLen)
}

// levenshteinDistance calcula a distância de Levenshtein entre duas strings
func (uc *ChangePasswordUseCase) levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

// min retorna o menor valor entre três inteiros
func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

// ==========================================
// MÉTODOS AUXILIARES
// ==========================================

// isCommonPassword verifica se a senha é muito comum
func isCommonPassword(password string) bool {
	commonPasswords := []string{
		"password", "123456", "123456789", "qwerty", "abc123",
		"password123", "admin", "letmein", "welcome", "monkey",
		"1234567890", "password1", "qwerty123", "dragon", "master",
		"12345678", "1234567", "1234566", "1234565", "1234564",
		"1234563", "1234562", "1234561", "1234560", "123455",
		"123454", "123453", "123452", "123451", "123450",
	}

	for _, common := range commonPasswords {
		if password == common {
			return true
		}
	}
	return false
}

// ChangePasswordWithConfirmation executa alteração de senha com confirmação
func (uc *ChangePasswordUseCase) ChangePasswordWithConfirmation(ctx context.Context, input ChangePasswordWithConfirmationInput) (*ChangePasswordWithConfirmationOutput, error) {
	// Verificar se as senhas coincidem
	if input.NewPassword != input.ConfirmNewPassword {
		uc.logger.Warn("Password confirmation does not match",
			"user_id", input.UserID,
		)
		return nil, fmt.Errorf("password confirmation does not match")
	}

	// Converter para input genérico
	genericInput := ChangePasswordInput{
		UserID:      input.UserID,
		OldPassword: input.OldPassword,
		NewPassword: input.NewPassword,
	}

	// Executar caso de uso genérico
	result, err := uc.Execute(ctx, genericInput)
	if err != nil {
		return nil, err
	}

	// Converter para output específico
	return &ChangePasswordWithConfirmationOutput{
		Message: result.Message,
	}, nil
}
