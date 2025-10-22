package shared

import (
	"time"

	"github.com/google/uuid"
)

// BaseEntity contém campos comuns a todas as entidades do domínio
// É uma abstração que evita repetição de código
type BaseEntity struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// NewBaseEntity cria uma nova BaseEntity com ID único e timestamps
func NewBaseEntity() BaseEntity {
	now := time.Now()
	return BaseEntity{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Touch atualiza o UpdatedAt para o momento atual
// Usado quando a entidade é modificada
func (e *BaseEntity) Touch() {
	e.UpdatedAt = time.Now()
}

// SoftDelete marca a entidade como deletada (soft delete)
// Não remove fisicamente do banco, apenas marca como deletada
func (e *BaseEntity) SoftDelete() {
	now := time.Now()
	e.DeletedAt = &now
	e.Touch()
}

// IsDeleted verifica se a entidade foi deletada
func (e *BaseEntity) IsDeleted() bool {
	return e.DeletedAt != nil
}

// Restore restaura uma entidade que foi soft deleted
func (e *BaseEntity) Restore() {
	e.DeletedAt = nil
	e.Touch()
}
