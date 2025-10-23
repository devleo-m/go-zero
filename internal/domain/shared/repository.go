// internal/domain/shared/repository.go
package shared

import (
	"context"

	"github.com/google/uuid"
)

// Repository é a interface genérica para todas as entidades
// T representa o tipo da entidade (User, Patient, etc)
type Repository[T any] interface {
	// ==========================================
	// CRUD BÁSICO
	// ==========================================

	// Create cria uma nova entidade
	Create(ctx context.Context, entity T) error

	// FindOne busca UMA entidade baseado em filtros
	// Retorna erro se não encontrar
	FindOne(ctx context.Context, filter QueryFilter) (T, error)

	// FindMany busca MÚLTIPLAS entidades baseado em filtros
	// Retorna slice vazio se não encontrar nenhuma
	FindMany(ctx context.Context, filter QueryFilter) ([]T, error)

	// FindByID busca uma entidade por ID (atalho comum)
	FindByID(ctx context.Context, id uuid.UUID) (T, error)

	// Update atualiza uma entidade existente
	Update(ctx context.Context, id uuid.UUID, entity T) error

	// Delete remove uma entidade (soft delete por padrão)
	Delete(ctx context.Context, id uuid.UUID) error

	// HardDelete remove permanentemente uma entidade
	HardDelete(ctx context.Context, id uuid.UUID) error

	// Restore restaura uma entidade soft deleted
	Restore(ctx context.Context, id uuid.UUID) error

	// ==========================================
	// PAGINAÇÃO
	// ==========================================

	// Paginate retorna resultados paginados com metadados
	Paginate(ctx context.Context, filter QueryFilter) (*PaginatedResult[T], error)

	// ==========================================
	// AGREGAÇÕES E CONTAGENS
	// ==========================================

	// Count conta entidades baseado em filtros
	Count(ctx context.Context, filter QueryFilter) (int64, error)

	// Exists verifica se existe alguma entidade que corresponde aos filtros
	Exists(ctx context.Context, filter QueryFilter) (bool, error)

	// Sum soma valores de um campo específico
	Sum(ctx context.Context, field string, filter QueryFilter) (float64, error)

	// Avg calcula a média de um campo
	Avg(ctx context.Context, field string, filter QueryFilter) (float64, error)

	// Min retorna o valor mínimo de um campo
	Min(ctx context.Context, field string, filter QueryFilter) (interface{}, error)

	// Max retorna o valor máximo de um campo
	Max(ctx context.Context, field string, filter QueryFilter) (interface{}, error)

	// ==========================================
	// OPERAÇÕES EM LOTE
	// ==========================================

	// CreateMany cria múltiplas entidades de uma vez
	CreateMany(ctx context.Context, entities []T) error

	// UpdateMany atualiza múltiplas entidades
	UpdateMany(ctx context.Context, filter QueryFilter, updates map[string]interface{}) (int64, error)

	// DeleteMany remove múltiplas entidades
	DeleteMany(ctx context.Context, filter QueryFilter) (int64, error)

	// ==========================================
	// TRANSAÇÕES
	// ==========================================

	// WithTransaction executa uma função dentro de uma transação
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error

	// ==========================================
	// QUERIES AVANÇADAS
	// ==========================================

	// FindFirst busca a primeira entidade que corresponde aos filtros
	FindFirst(ctx context.Context, filter QueryFilter) (T, error)

	// FindLast busca a última entidade que corresponde aos filtros
	FindLast(ctx context.Context, filter QueryFilter) (T, error)

	// Distinct retorna valores únicos de um campo
	Distinct(ctx context.Context, field string, filter QueryFilter) ([]interface{}, error)

	// GroupBy agrupa resultados por um campo
	GroupBy(ctx context.Context, field string, filter QueryFilter) (map[string][]T, error)
}
