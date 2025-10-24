package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserModel representa o modelo GORM para User
type UserModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string         `gorm:"size:100;not null"`
	Email     string         `gorm:"size:254;uniqueIndex;not null"`
	Password  string         `gorm:"size:255;not null"`
	Phone     *string        `gorm:"size:20"`
	Role      string         `gorm:"size:20;not null;default:'user'"`
	Status    string         `gorm:"size:20;not null;default:'active'"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName define o nome da tabela
func (UserModel) TableName() string {
	return "users"
}
