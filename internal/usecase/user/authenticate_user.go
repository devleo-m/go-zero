package user

import (
	"context"
	"fmt"
	"time"

	"github.com/devleo-m/go-zero/internal/domain/user"
	"github.com/google/uuid"
)

// AuthenticateUserUseCase implementa o caso de uso de autenticação
type AuthenticateUserUseCase struct {
	userRepo     user.Repository
	queryHelper  *UserQueryHelper
	jwtService   JWTService
	tokenService TokenService
	logger       Logger
}

// NewAuthenticateUserUseCase cria uma nova instância do caso de uso
func NewAuthenticateUserUseCase(
	userRepo user.Repository,
	queryHelper *UserQueryHelper,
	jwtService JWTService,
	tokenService TokenService,
	logger Logger,
) *AuthenticateUserUseCase {
	return &AuthenticateUserUseCase{
		userRepo:     userRepo,
		queryHelper:  queryHelper,
		jwtService:   jwtService,
		tokenService: tokenService,
		logger:       logger,
	}
}

// Execute executa o caso de uso de autenticação
func (uc *AuthenticateUserUseCase) Execute(ctx context.Context, input AuthenticateUserInput) (*AuthenticateUserOutput, error) {
	// 1. VALIDAR INPUT
	if err := uc.validateInput(input); err != nil {
		uc.logger.Warn("Invalid input for authentication",
			"error", err,
			"email", input.Email,
		)
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	uc.logger.Debug("Starting authentication process",
		"email", input.Email,
	)

	// 2. BUSCAR USUÁRIO POR EMAIL
	domainUser, err := uc.queryHelper.FindUserByEmail(ctx, input.Email)
	if err != nil {
		uc.logger.Error("Failed to find user by email",
			"error", err,
			"email", input.Email,
		)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if domainUser == nil {
		uc.logger.Warn("User not found during authentication",
			"email", input.Email,
		)
		return nil, user.ErrUserNotFound
	}

	uc.logger.Debug("User found for authentication",
		"user_id", domainUser.ID,
		"email", input.Email,
		"status", domainUser.Status.String(),
	)

	// 3. VERIFICAR SE USUÁRIO ESTÁ ATIVO
	if !domainUser.IsActive() {
		uc.logger.Warn("Inactive user attempted authentication",
			"user_id", domainUser.ID,
			"email", input.Email,
			"status", domainUser.Status.String(),
		)
		return nil, user.ErrUserInactive
	}

	// 4. VERIFICAR SENHA
	if err := domainUser.Authenticate(input.Password); err != nil {
		uc.logger.Warn("Invalid password during authentication",
			"user_id", domainUser.ID,
			"email", input.Email,
		)
		return nil, user.ErrInvalidCredentials
	}

	uc.logger.Debug("Password verified successfully",
		"user_id", domainUser.ID,
		"email", input.Email,
	)

	// 5. ATUALIZAR DADOS DE LOGIN
	now := time.Now()
	domainUser.LastLoginAt = &now
	domainUser.LoginCount++

	// 6. SALVAR ALTERAÇÕES NO BANCO
	if err := uc.userRepo.Update(ctx, domainUser.ID, domainUser); err != nil {
		uc.logger.Error("Failed to update user login data",
			"error", err,
			"user_id", domainUser.ID,
		)
		// Não falhar a autenticação por causa disso, apenas logar
		uc.logger.Warn("Authentication succeeded but failed to update login data",
			"user_id", domainUser.ID,
		)
	}

	// 7. GERAR TOKENS JWT
	accessToken, expiresIn, err := uc.jwtService.GenerateAccessToken(
		domainUser.ID.String(),
		domainUser.Email.String(),
		domainUser.Role.String(),
	)
	if err != nil {
		uc.logger.Error("Failed to generate access token",
			"error", err,
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(domainUser.ID.String())
	if err != nil {
		uc.logger.Error("Failed to generate refresh token",
			"error", err,
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	uc.logger.Info("User authenticated successfully",
		"user_id", domainUser.ID,
		"email", input.Email,
		"role", domainUser.Role.String(),
		"login_count", domainUser.LoginCount,
	)

	// 8. RETORNAR RESULTADO
	output := &AuthenticateUserOutput{
		User:         ToUserOutput(domainUser),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}

	uc.logger.Info("Authenticate user use case completed successfully",
		"user_id", domainUser.ID,
		"email", input.Email,
	)

	return output, nil
}

// validateInput valida os dados de entrada
func (uc *AuthenticateUserUseCase) validateInput(input AuthenticateUserInput) error {
	if input.Email == "" {
		return fmt.Errorf("email is required")
	}

	if input.Password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

// ==========================================
// MÉTODOS AUXILIARES
// ==========================================

// AuthenticateWithToken autentica usando token de refresh
func (uc *AuthenticateUserUseCase) AuthenticateWithToken(ctx context.Context, refreshToken string) (*AuthenticateUserOutput, error) {
	uc.logger.Debug("Starting token-based authentication",
		"refresh_token", refreshToken[:10]+"...", // Log apenas parte do token por segurança
	)

	// Validar token de refresh
	userID, err := uc.tokenService.ValidateActivationToken(refreshToken)
	if err != nil {
		uc.logger.Warn("Invalid refresh token",
			"error", err,
		)
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Buscar usuário por ID
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	domainUser, err := uc.userRepo.FindByID(ctx, userIDUUID)
	if err != nil {
		uc.logger.Error("Failed to find user by ID from token",
			"error", err,
			"user_id", userID,
		)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Verificar se usuário está ativo
	if !domainUser.IsActive() {
		uc.logger.Warn("Inactive user attempted token authentication",
			"user_id", domainUser.ID,
			"status", domainUser.Status.String(),
		)
		return nil, user.ErrUserInactive
	}

	// Gerar novos tokens
	accessToken, expiresIn, err := uc.jwtService.GenerateAccessToken(
		domainUser.ID.String(),
		domainUser.Email.String(),
		domainUser.Role.String(),
	)
	if err != nil {
		uc.logger.Error("Failed to generate new access token",
			"error", err,
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := uc.jwtService.GenerateRefreshToken(domainUser.ID.String())
	if err != nil {
		uc.logger.Error("Failed to generate new refresh token",
			"error", err,
			"user_id", domainUser.ID,
		)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	uc.logger.Info("Token authentication successful",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
	)

	return &AuthenticateUserOutput{
		User:         ToUserOutput(domainUser),
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// RefreshAccessToken renova token de acesso usando refresh token
func (uc *AuthenticateUserUseCase) RefreshAccessToken(ctx context.Context, refreshToken string) (*AuthenticateUserOutput, error) {
	uc.logger.Debug("Starting access token refresh",
		"refresh_token", refreshToken[:10]+"...",
	)

	// Renovar token usando o serviço JWT
	newAccessToken, expiresIn, err := uc.jwtService.RefreshAccessToken(refreshToken)
	if err != nil {
		uc.logger.Warn("Failed to refresh access token",
			"error", err,
		)
		return nil, fmt.Errorf("failed to refresh access token: %w", err)
	}

	// Extrair user ID do token para buscar dados do usuário
	claims, err := uc.jwtService.ValidateToken(newAccessToken)
	if err != nil {
		uc.logger.Error("Failed to validate new access token",
			"error", err,
		)
		return nil, fmt.Errorf("failed to validate new access token: %w", err)
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		uc.logger.Error("Invalid user_id in token claims",
			"claims", claims,
		)
		return nil, fmt.Errorf("invalid user_id in token claims")
	}

	// Buscar usuário
	userIDUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	domainUser, err := uc.userRepo.FindByID(ctx, userIDUUID)
	if err != nil {
		uc.logger.Error("Failed to find user during token refresh",
			"error", err,
			"user_id", userIDStr,
		)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	uc.logger.Info("Access token refreshed successfully",
		"user_id", domainUser.ID,
		"email", domainUser.Email.String(),
	)

	return &AuthenticateUserOutput{
		User:         ToUserOutput(domainUser),
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken, // Manter o mesmo refresh token
		ExpiresIn:    expiresIn,
	}, nil
}

// Logout invalida tokens do usuário
func (uc *AuthenticateUserUseCase) Logout(ctx context.Context, accessToken, refreshToken string) error {
	uc.logger.Debug("Starting logout process")

	// Revogar access token
	if err := uc.jwtService.RevokeToken(accessToken); err != nil {
		uc.logger.Warn("Failed to revoke access token",
			"error", err,
		)
		// Não falhar o logout por causa disso
	}

	// Revogar refresh token
	if err := uc.jwtService.RevokeToken(refreshToken); err != nil {
		uc.logger.Warn("Failed to revoke refresh token",
			"error", err,
		)
		// Não falhar o logout por causa disso
	}

	uc.logger.Info("Logout completed successfully")
	return nil
}
