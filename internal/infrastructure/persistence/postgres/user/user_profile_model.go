package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserProfileModel representa dados de perfil do usuário
// Separado para lazy loading e melhor performance
type UserProfileModel struct {
	// Campos base
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time      `gorm:"not null;index"`
	UpdatedAt time.Time      `gorm:"not null;index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Chave estrangeira
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;constraint:OnDelete:CASCADE"`

	// Dados de perfil
	EmailVerifiedAt *time.Time `gorm:"index"`
	LastLoginAt     *time.Time `gorm:"index"`
	LoginCount      int        `gorm:"default:0;not null"`
	IPAddress       *string    `gorm:"type:varchar(45)"` // IPv6 suporta até 45 chars
	UserAgent       *string    `gorm:"type:text"`
	AvatarURL       *string    `gorm:"type:varchar(500)"`
	Bio             *string    `gorm:"type:text"`
	Location        *string    `gorm:"type:varchar(255)"`
	Website         *string    `gorm:"type:varchar(500)"`
}

// TableName define o nome da tabela
func (UserProfileModel) TableName() string {
	return "user_profiles"
}

// BeforeCreate hook executado antes de criar o perfil
func (p *UserProfileModel) BeforeCreate(tx *gorm.DB) error {
	// Validar UserID
	if p.UserID == uuid.Nil {
		return gorm.ErrInvalidData
	}

	// Gerar ID se não existir
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	// Definir timestamps
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}

	return nil
}

// BeforeUpdate hook executado antes de atualizar o perfil
func (p *UserProfileModel) BeforeUpdate(tx *gorm.DB) error {
	// Atualizar timestamp
	p.UpdatedAt = time.Now()
	return nil
}

// IsEmailVerified verifica se o email foi verificado
func (p *UserProfileModel) IsEmailVerified() bool {
	return p.EmailVerifiedAt != nil
}

// MarkEmailAsVerified marca o email como verificado
func (p *UserProfileModel) MarkEmailAsVerified() {
	now := time.Now()
	p.EmailVerifiedAt = &now
}

// RecordLogin registra um novo login
func (p *UserProfileModel) RecordLogin(ipAddress, userAgent string) {
	now := time.Now()
	p.LastLoginAt = &now
	p.LoginCount++

	if ipAddress != "" {
		p.IPAddress = &ipAddress
	}
	if userAgent != "" {
		p.UserAgent = &userAgent
	}
}
