package user

import (
	"context"
)

// ==========================================
// INTERFACES DE SERVIÇOS EXTERNOS
// ==========================================

// Logger interface para logging
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// EmailService interface para serviços de email
type EmailService interface {
	// Enviar email de boas-vindas
	SendWelcomeEmail(ctx context.Context, email, name string) error

	// Enviar email de confirmação de ativação
	SendActivationConfirmationEmail(ctx context.Context, email, name string) error

	// Enviar email de notificação de desativação
	SendDeactivationNotificationEmail(ctx context.Context, email, name string) error

	// Enviar email de notificação de suspensão
	SendSuspensionNotificationEmail(ctx context.Context, email, name string) error

	// Enviar email de notificação de mudança de role
	SendRoleChangeNotificationEmail(ctx context.Context, email, name, newRole string) error

	// Enviar email de notificação de mudança de senha
	SendPasswordChangeNotificationEmail(ctx context.Context, email, name string) error

	// Enviar email de reset de senha
	SendPasswordResetEmail(ctx context.Context, email, name, resetToken string) error
}

// JWTService interface para serviços de JWT
type JWTService interface {
	// Gerar token de acesso
	GenerateAccessToken(userID string, email string, role string) (string, int, error)

	// Gerar token de refresh
	GenerateRefreshToken(userID string) (string, error)

	// Validar token
	ValidateToken(token string) (map[string]interface{}, error)

	// Renovar token de acesso usando refresh token
	RefreshAccessToken(refreshToken string) (string, int, error)

	// Revogar token
	RevokeToken(token string) error
}

// TokenService interface para serviços de token (ativação, reset, etc.)
type TokenService interface {
	// Validar token de ativação
	ValidateActivationToken(token string) (string, error) // retorna userID

	// Gerar token de ativação
	GenerateActivationToken(userID string) (string, error)

	// Validar token de reset de senha
	ValidatePasswordResetToken(token string) (string, error) // retorna userID

	// Gerar token de reset de senha
	GeneratePasswordResetToken(userID string) (string, error)

	// Invalidar token
	InvalidateToken(token string) error
}

// CacheService interface para serviços de cache
type CacheService interface {
	// Obter valor do cache
	Get(ctx context.Context, key string) (string, error)

	// Definir valor no cache
	Set(ctx context.Context, key string, value string, ttl int) error

	// Deletar valor do cache
	Delete(ctx context.Context, key string) error

	// Verificar se chave existe
	Exists(ctx context.Context, key string) (bool, error)

	// Incrementar contador
	Increment(ctx context.Context, key string) (int64, error)

	// Definir expiração
	Expire(ctx context.Context, key string, ttl int) error
}

// NotificationService interface para serviços de notificação
type NotificationService interface {
	// Enviar notificação por email
	SendEmailNotification(ctx context.Context, to, subject, body string) error

	// Enviar notificação por SMS
	SendSMSNotification(ctx context.Context, phone, message string) error

	// Enviar notificação push
	SendPushNotification(ctx context.Context, userID, title, body string) error

	// Enviar notificação in-app
	SendInAppNotification(ctx context.Context, userID, message string) error
}

// AuditService interface para serviços de auditoria
type AuditService interface {
	// Registrar evento de auditoria
	LogEvent(ctx context.Context, userID, action, resource string, details map[string]interface{}) error

	// Registrar evento de segurança
	LogSecurityEvent(ctx context.Context, userID, event string, details map[string]interface{}) error

	// Registrar evento de negócio
	LogBusinessEvent(ctx context.Context, userID, event string, details map[string]interface{}) error
}

// ==========================================
// INTERFACES DE CASOS DE USO
// ==========================================

// CreateUserUseCaseInterface interface para o caso de uso de criação
type CreateUserUseCaseInterface interface {
	Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error)
}

// AuthenticateUserUseCaseInterface interface para o caso de uso de autenticação
type AuthenticateUserUseCaseInterface interface {
	Execute(ctx context.Context, input AuthenticateUserInput) (*AuthenticateUserOutput, error)
}

// GetUserUseCaseInterface interface para o caso de uso de busca
type GetUserUseCaseInterface interface {
	Execute(ctx context.Context, input GetUserInput) (*GetUserOutput, error)
	GetUserByEmail(ctx context.Context, email string) (*GetUserOutput, error)
	CheckUserExists(ctx context.Context, userID string) (bool, error)
	GetUserBasicInfo(ctx context.Context, userID string) (*UserBasicInfo, error)
}

// ListUsersUseCaseInterface interface para o caso de uso de listagem
type ListUsersUseCaseInterface interface {
	Execute(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error)
	GetUsersByRole(ctx context.Context, role string, page, pageSize int) (*ListUsersOutput, error)
	GetUsersByStatus(ctx context.Context, status string, page, pageSize int) (*ListUsersOutput, error)
	SearchUsers(ctx context.Context, searchTerm string, page, pageSize int) (*ListUsersOutput, error)
	GetActiveUsers(ctx context.Context, page, pageSize int) (*ListUsersOutput, error)
	GetPendingUsers(ctx context.Context, page, pageSize int) (*ListUsersOutput, error)
	GetUserStats(ctx context.Context) (*UserStatsOutput, error)
}

// UpdateUserUseCaseInterface interface para o caso de uso de atualização
type UpdateUserUseCaseInterface interface {
	Execute(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error)
	UpdateUserProfile(ctx context.Context, input UpdateUserProfileInput) (*UpdateUserProfileOutput, error)
}

// ChangePasswordUseCaseInterface interface para o caso de uso de alteração de senha
type ChangePasswordUseCaseInterface interface {
	Execute(ctx context.Context, input ChangePasswordInput) (*ChangePasswordOutput, error)
	ChangePasswordWithConfirmation(ctx context.Context, input ChangePasswordWithConfirmationInput) (*ChangePasswordWithConfirmationOutput, error)
}

// ActivateUserUseCaseInterface interface para o caso de uso de ativação
type ActivateUserUseCaseInterface interface {
	Execute(ctx context.Context, input ActivateUserInput) (*ActivateUserOutput, error)
	ActivateUserByEmail(ctx context.Context, email string) (*ActivateUserOutput, error)
	ActivateUserWithToken(ctx context.Context, token string) (*ActivateUserOutput, error)
}

// DeactivateUserUseCaseInterface interface para o caso de uso de desativação
type DeactivateUserUseCaseInterface interface {
	Execute(ctx context.Context, input DeactivateUserInput) (*DeactivateUserOutput, error)
	DeactivateUserByEmail(ctx context.Context, email string) (*DeactivateUserOutput, error)
}

// SuspendUserUseCaseInterface interface para o caso de uso de suspensão
type SuspendUserUseCaseInterface interface {
	Execute(ctx context.Context, input SuspendUserInput) (*SuspendUserOutput, error)
	SuspendUserByEmail(ctx context.Context, email string) (*SuspendUserOutput, error)
}

// ChangeRoleUseCaseInterface interface para o caso de uso de alteração de role
type ChangeRoleUseCaseInterface interface {
	Execute(ctx context.Context, input ChangeRoleInput) (*ChangeRoleOutput, error)
	ChangeRoleByEmail(ctx context.Context, email, newRole string, requesterID string) (*ChangeRoleOutput, error)
}

// ==========================================
// INTERFACE AGGREGATE (FACADE)
// ==========================================

// UserUseCaseAggregate interface que agrega todos os casos de uso
type UserUseCaseAggregate interface {
	// Casos de uso básicos
	CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error)
	AuthenticateUser(ctx context.Context, input AuthenticateUserInput) (*AuthenticateUserOutput, error)
	GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error)
	ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error)
	UpdateUser(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error)

	// Casos de uso de senha
	ChangePassword(ctx context.Context, input ChangePasswordInput) (*ChangePasswordOutput, error)
	ChangePasswordWithConfirmation(ctx context.Context, input ChangePasswordWithConfirmationInput) (*ChangePasswordWithConfirmationOutput, error)

	// Casos de uso de status
	ActivateUser(ctx context.Context, input ActivateUserInput) (*ActivateUserOutput, error)
	DeactivateUser(ctx context.Context, input DeactivateUserInput) (*DeactivateUserOutput, error)
	SuspendUser(ctx context.Context, input SuspendUserInput) (*SuspendUserOutput, error)

	// Casos de uso de role
	ChangeRole(ctx context.Context, input ChangeRoleInput) (*ChangeRoleOutput, error)

	// Métodos auxiliares
	GetUserByEmail(ctx context.Context, email string) (*GetUserOutput, error)
	CheckUserExists(ctx context.Context, userID string) (bool, error)
	GetUserBasicInfo(ctx context.Context, userID string) (*UserBasicInfo, error)
	GetUserStats(ctx context.Context) (*UserStatsOutput, error)
}
