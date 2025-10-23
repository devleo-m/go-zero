package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserAuthDataModel representa dados de autenticação do usuário
// Separado para segurança e lazy loading
type UserAuthDataModel struct {
	// Campos base
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time      `gorm:"not null;index"`
	UpdatedAt time.Time      `gorm:"not null;index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Chave estrangeira
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;constraint:OnDelete:CASCADE"`

	// Tokens de reset de senha
	PasswordResetToken   *string    `gorm:"type:varchar(255);uniqueIndex"`
	PasswordResetExpires *time.Time `gorm:"index"`

	// Tokens de ativação
	ActivationToken   *string    `gorm:"type:varchar(255);uniqueIndex"`
	ActivationExpires *time.Time `gorm:"index"`

	// 2FA (Two-Factor Authentication)
	TwoFactorSecret      *string  `gorm:"type:varchar(255)"`
	TwoFactorEnabled     bool     `gorm:"default:false;not null"`
	TwoFactorBackupCodes []string `gorm:"type:text;serializer:json"`

	// Refresh tokens (para JWT)
	RefreshToken     *string    `gorm:"type:varchar(500);uniqueIndex"`
	RefreshExpires   *time.Time `gorm:"index"`
	RefreshTokenHash *string    `gorm:"type:varchar(255)"` // Hash do token para segurança

	// Dados de sessão
	LastPasswordChange  *time.Time `gorm:"index"`
	FailedLoginAttempts int        `gorm:"default:0;not null"`
	LockedUntil         *time.Time `gorm:"index"`
}

// TableName define o nome da tabela
func (UserAuthDataModel) TableName() string {
	return "user_auth_data"
}

// BeforeCreate hook executado antes de criar os dados de auth
func (a *UserAuthDataModel) BeforeCreate(tx *gorm.DB) error {
	// Validar UserID
	if a.UserID == uuid.Nil {
		return gorm.ErrInvalidData
	}

	// Gerar ID se não existir
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}

	// Definir timestamps
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	if a.UpdatedAt.IsZero() {
		a.UpdatedAt = now
	}

	return nil
}

// BeforeUpdate hook executado antes de atualizar os dados de auth
func (a *UserAuthDataModel) BeforeUpdate(tx *gorm.DB) error {
	// Atualizar timestamp
	a.UpdatedAt = time.Now()
	return nil
}

// IsPasswordResetTokenValid verifica se o token de reset é válido
func (a *UserAuthDataModel) IsPasswordResetTokenValid() bool {
	if a.PasswordResetToken == nil || a.PasswordResetExpires == nil {
		return false
	}
	return time.Now().Before(*a.PasswordResetExpires)
}

// IsActivationTokenValid verifica se o token de ativação é válido
func (a *UserAuthDataModel) IsActivationTokenValid() bool {
	if a.ActivationToken == nil || a.ActivationExpires == nil {
		return false
	}
	return time.Now().Before(*a.ActivationExpires)
}

// IsRefreshTokenValid verifica se o refresh token é válido
func (a *UserAuthDataModel) IsRefreshTokenValid() bool {
	if a.RefreshToken == nil || a.RefreshExpires == nil {
		return false
	}
	return time.Now().Before(*a.RefreshExpires)
}

// IsAccountLocked verifica se a conta está bloqueada
func (a *UserAuthDataModel) IsAccountLocked() bool {
	if a.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*a.LockedUntil)
}

// SetPasswordResetToken define um novo token de reset
func (a *UserAuthDataModel) SetPasswordResetToken(token string, expiresIn time.Duration) {
	expires := time.Now().Add(expiresIn)
	a.PasswordResetToken = &token
	a.PasswordResetExpires = &expires
}

// SetActivationToken define um novo token de ativação
func (a *UserAuthDataModel) SetActivationToken(token string, expiresIn time.Duration) {
	expires := time.Now().Add(expiresIn)
	a.ActivationToken = &token
	a.ActivationExpires = &expires
}

// SetRefreshToken define um novo refresh token
func (a *UserAuthDataModel) SetRefreshToken(token, tokenHash string, expiresIn time.Duration) {
	expires := time.Now().Add(expiresIn)
	a.RefreshToken = &token
	a.RefreshTokenHash = &tokenHash
	a.RefreshExpires = &expires
}

// ClearPasswordResetToken limpa o token de reset
func (a *UserAuthDataModel) ClearPasswordResetToken() {
	a.PasswordResetToken = nil
	a.PasswordResetExpires = nil
}

// ClearActivationToken limpa o token de ativação
func (a *UserAuthDataModel) ClearActivationToken() {
	a.ActivationToken = nil
	a.ActivationExpires = nil
}

// ClearRefreshToken limpa o refresh token
func (a *UserAuthDataModel) ClearRefreshToken() {
	a.RefreshToken = nil
	a.RefreshTokenHash = nil
	a.RefreshExpires = nil
}

// RecordFailedLogin registra uma tentativa de login falhada
func (a *UserAuthDataModel) RecordFailedLogin() {
	a.FailedLoginAttempts++

	// Bloquear conta após 5 tentativas por 30 minutos
	if a.FailedLoginAttempts >= 5 {
		lockUntil := time.Now().Add(30 * time.Minute)
		a.LockedUntil = &lockUntil
	}
}

// ResetFailedLogins reseta as tentativas de login falhadas
func (a *UserAuthDataModel) ResetFailedLogins() {
	a.FailedLoginAttempts = 0
	a.LockedUntil = nil
}
