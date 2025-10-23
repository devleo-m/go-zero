// internal/infrastructure/persistence/postgres/user/model.go
package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserModel representa o modelo GORM para usuários
// Mapeia diretamente para a tabela 'users' no PostgreSQL
type UserModel struct {
	// ==========================================
	// CAMPOS BASE (BaseEntity)
	// ==========================================

	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `gorm:"not null;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;index" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// ==========================================
	// CAMPOS ESPECÍFICOS DO USER
	// ==========================================

	// Dados pessoais
	Name  string `gorm:"type:varchar(255);not null;index" json:"name"`
	Email string `gorm:"type:varchar(255);not null;uniqueIndex;index" json:"email"`

	// Senha (hash)
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`

	// Telefone (opcional)
	Phone *string `gorm:"type:varchar(20);index" json:"phone,omitempty"`

	// Status e role
	Status string `gorm:"type:varchar(50);not null;default:'pending';index" json:"status"`
	Role   string `gorm:"type:varchar(50);not null;default:'user';index" json:"role"`

	// Dados de autenticação
	EmailVerifiedAt *time.Time `gorm:"index" json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time `gorm:"index" json:"last_login_at,omitempty"`
	LoginCount      int        `gorm:"default:0" json:"login_count"`

	// Dados de recuperação
	PasswordResetToken   *string    `gorm:"type:varchar(255);index" json:"-"`
	PasswordResetExpires *time.Time `json:"password_reset_expires,omitempty"`

	// Dados de ativação
	ActivationToken   *string    `gorm:"type:varchar(255);index" json:"-"`
	ActivationExpires *time.Time `json:"activation_expires,omitempty"`

	// Dados de 2FA
	TwoFactorSecret  *string `gorm:"type:varchar(255)" json:"-"`
	TwoFactorEnabled bool    `gorm:"default:false" json:"two_factor_enabled"`

	// Metadados
	IPAddress *string `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent *string `gorm:"type:text" json:"user_agent,omitempty"`
	Timezone  *string `gorm:"type:varchar(50);default:'UTC'" json:"timezone,omitempty"`
	Language  *string `gorm:"type:varchar(10);default:'pt-BR'" json:"language,omitempty"`

	// Campos de auditoria
	CreatedBy *uuid.UUID `gorm:"type:uuid;index" json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid;index" json:"updated_by,omitempty"`
	DeletedBy *uuid.UUID `gorm:"type:uuid;index" json:"deleted_by,omitempty"`
}

// TableName define o nome da tabela no banco
func (UserModel) TableName() string {
	return "users"
}

// ==========================================
// HOOKS GORM
// ==========================================

// BeforeCreate hook executado antes de criar
func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	// Gerar ID se não existir
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	// Definir timestamps se não existirem
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = now
	}

	// Validar campos obrigatórios
	if u.Name == "" {
		return gorm.ErrInvalidData
	}
	if u.Email == "" {
		return gorm.ErrInvalidData
	}
	if u.PasswordHash == "" {
		return gorm.ErrInvalidData
	}

	return nil
}

// BeforeUpdate hook executado antes de atualizar
func (u *UserModel) BeforeUpdate(tx *gorm.DB) error {
	// Atualizar timestamp
	u.UpdatedAt = time.Now()

	// Validar campos obrigatórios
	if u.Name == "" {
		return gorm.ErrInvalidData
	}
	if u.Email == "" {
		return gorm.ErrInvalidData
	}

	return nil
}

// BeforeDelete hook executado antes de deletar
func (u *UserModel) BeforeDelete(tx *gorm.DB) error {
	// Atualizar timestamp
	u.UpdatedAt = time.Now()

	return nil
}

// ==========================================
// MÉTODOS DE CONVENIÊNCIA
// ==========================================

// IsActive verifica se o usuário está ativo
func (u *UserModel) IsActive() bool {
	return u.Status == "active"
}

// IsAdmin verifica se o usuário é admin
func (u *UserModel) IsAdmin() bool {
	return u.Role == "admin"
}

// IsEmailVerified verifica se o email foi verificado
func (u *UserModel) IsEmailVerified() bool {
	return u.EmailVerifiedAt != nil
}

// IsTwoFactorEnabled verifica se 2FA está habilitado
func (u *UserModel) IsTwoFactorEnabled() bool {
	return u.TwoFactorEnabled && u.TwoFactorSecret != nil
}

// MarkEmailVerified marca o email como verificado
func (u *UserModel) MarkEmailVerified() {
	now := time.Now()
	u.EmailVerifiedAt = &now
	u.UpdatedAt = now
}

// MarkLastLogin registra o último login
func (u *UserModel) MarkLastLogin(ip, userAgent string) {
	now := time.Now()
	u.LastLoginAt = &now
	u.LoginCount++
	u.IPAddress = &ip
	u.UserAgent = &userAgent
	u.UpdatedAt = now
}

// SetPasswordResetToken define token de reset de senha
func (u *UserModel) SetPasswordResetToken(token string, expiresAt time.Time) {
	u.PasswordResetToken = &token
	u.PasswordResetExpires = &expiresAt
	u.UpdatedAt = time.Now()
}

// ClearPasswordResetToken limpa token de reset de senha
func (u *UserModel) ClearPasswordResetToken() {
	u.PasswordResetToken = nil
	u.PasswordResetExpires = nil
	u.UpdatedAt = time.Now()
}

// SetActivationToken define token de ativação
func (u *UserModel) SetActivationToken(token string, expiresAt time.Time) {
	u.ActivationToken = &token
	u.ActivationExpires = &expiresAt
	u.UpdatedAt = time.Now()
}

// ClearActivationToken limpa token de ativação
func (u *UserModel) ClearActivationToken() {
	u.ActivationToken = nil
	u.ActivationExpires = nil
	u.UpdatedAt = time.Now()
}

// EnableTwoFactor habilita 2FA
func (u *UserModel) EnableTwoFactor(secret string) {
	u.TwoFactorSecret = &secret
	u.TwoFactorEnabled = true
	u.UpdatedAt = time.Now()
}

// DisableTwoFactor desabilita 2FA
func (u *UserModel) DisableTwoFactor() {
	u.TwoFactorSecret = nil
	u.TwoFactorEnabled = false
	u.UpdatedAt = time.Now()
}

// ==========================================
// ÍNDICES E CONSTRAINTS
// ==========================================

// GetIndexes retorna os índices da tabela
func (UserModel) GetIndexes() []string {
	return []string{
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);",
		"CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);",
		"CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);",
		"CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);",
		"CREATE INDEX IF NOT EXISTS idx_users_updated_at ON users(updated_at);",
		"CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);",
		"CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);",
		"CREATE INDEX IF NOT EXISTS idx_users_email_verified_at ON users(email_verified_at);",
		"CREATE INDEX IF NOT EXISTS idx_users_last_login_at ON users(last_login_at);",
		"CREATE INDEX IF NOT EXISTS idx_users_created_by ON users(created_by);",
		"CREATE INDEX IF NOT EXISTS idx_users_updated_by ON users(updated_by);",
		"CREATE INDEX IF NOT EXISTS idx_users_deleted_by ON users(deleted_by);",
		"CREATE INDEX IF NOT EXISTS idx_users_password_reset_token ON users(password_reset_token);",
		"CREATE INDEX IF NOT EXISTS idx_users_activation_token ON users(activation_token);",
	}
}

// GetConstraints retorna as constraints da tabela
func (UserModel) GetConstraints() []string {
	return []string{
		"ALTER TABLE users ADD CONSTRAINT chk_users_status CHECK (status IN ('pending', 'active', 'inactive', 'suspended'));",
		"ALTER TABLE users ADD CONSTRAINT chk_users_role CHECK (role IN ('admin', 'manager', 'user', 'guest'));",
		"ALTER TABLE users ADD CONSTRAINT chk_users_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$');",
		"ALTER TABLE users ADD CONSTRAINT chk_users_login_count CHECK (login_count >= 0);",
	}
}
