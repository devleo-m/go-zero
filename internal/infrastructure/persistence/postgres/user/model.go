package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserModel representa a tabela principal de usuários
// Contém APENAS os campos essenciais para operações básicas
type UserModel struct {
	// Campos base (auditoria)
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time      `gorm:"not null;index"`
	UpdatedAt time.Time      `gorm:"not null;index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Campos essenciais do usuário
	Name         string  `gorm:"type:varchar(255);not null;index"`
	Email        string  `gorm:"type:varchar(255);not null;uniqueIndex"`
	PasswordHash string  `gorm:"type:varchar(255);not null"`
	Phone        *string `gorm:"type:varchar(20);index"`
	Status       string  `gorm:"type:varchar(50);not null;index;default:'pending'"`
	Role         string  `gorm:"type:varchar(50);not null;index;default:'user'"`

	// Relacionamentos (lazy load)
	Profile     *UserProfileModel     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	AuthData    *UserAuthDataModel    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Preferences *UserPreferencesModel `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName define o nome da tabela
func (UserModel) TableName() string {
	return "users"
}

// BeforeCreate hook executado antes de criar o usuário
func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
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

	// Validar Status
	if !isValidStatus(u.Status) {
		return gorm.ErrInvalidData
	}

	// Validar Role
	if !isValidRole(u.Role) {
		return gorm.ErrInvalidData
	}

	// Gerar ID se não existir
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	// Definir timestamps
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = now
	}

	return nil
}

// BeforeUpdate hook executado antes de atualizar o usuário
func (u *UserModel) BeforeUpdate(tx *gorm.DB) error {
	// Validar Status se foi alterado
	if tx.Statement.Changed("Status") && !isValidStatus(u.Status) {
		return gorm.ErrInvalidData
	}

	// Validar Role se foi alterado
	if tx.Statement.Changed("Role") && !isValidRole(u.Role) {
		return gorm.ErrInvalidData
	}

	// Atualizar timestamp
	u.UpdatedAt = time.Now()

	return nil
}

// IsActive verifica se o usuário está ativo
func (u *UserModel) IsActive() bool {
	return u.Status == "active"
}

// IsAdmin verifica se o usuário é admin
func (u *UserModel) IsAdmin() bool {
	return u.Role == "admin"
}

// GetDisplayName retorna o nome para exibição
func (u *UserModel) GetDisplayName() string {
	if u.Name != "" {
		return u.Name
	}
	return u.Email
}

// Validação de enums
func isValidStatus(status string) bool {
	valid := []string{"pending", "active", "inactive", "suspended"}
	for _, s := range valid {
		if s == status {
			return true
		}
	}
	return false
}

func isValidRole(role string) bool {
	valid := []string{"admin", "manager", "user", "guest"}
	for _, r := range valid {
		if r == role {
			return true
		}
	}
	return false
}
