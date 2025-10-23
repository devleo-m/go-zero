// internal/infrastructure/persistence/postgres/user/repository.go
package user

import (
	"context"
	"fmt"
	"time"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository implementa a interface do repositório de usuários usando GORM
type Repository struct {
	db *gorm.DB
}

// NewRepository cria uma nova instância do repositório
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// ==========================================
// CRUD BÁSICO
// ==========================================

// Create cria um novo usuário
func (r *Repository) Create(ctx context.Context, entity *user.User) error {
	model := ToModel(entity)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Atualizar entidade com dados do banco
	updatedUser, err := ToDomain(model)
	if err != nil {
		return fmt.Errorf("failed to convert model to domain: %w", err)
	}
	*entity = *updatedUser

	return nil
}

// FindOne busca uma entidade baseado em filtros
func (r *Repository) FindOne(ctx context.Context, filter shared.QueryFilter) (*user.User, error) {
	var model UserModel

	query, err := QueryFilterToGORM(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to convert filter: %w", err)
	}

	if err := r.db.WithContext(ctx).Scopes(query).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.NewDomainError("USER_NOT_FOUND", "user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return ToDomain(&model)
}

// FindMany busca múltiplas entidades baseado em filtros
func (r *Repository) FindMany(ctx context.Context, filter shared.QueryFilter) ([]*user.User, error) {
	var models []*UserModel

	query, err := QueryFilterToGORM(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to convert filter: %w", err)
	}

	if err := r.db.WithContext(ctx).Scopes(query).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	return ToDomainSlice(models)
}

// FindByID busca uma entidade por ID
func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var model UserModel

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.NewDomainError("USER_NOT_FOUND", "user not found")
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return ToDomain(&model)
}

// Update atualiza uma entidade existente
func (r *Repository) Update(ctx context.Context, id uuid.UUID, entity *user.User) error {
	model := ToModel(entity)

	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return shared.NewDomainError("USER_NOT_FOUND", "user not found")
	}

	return nil
}

// Delete remove uma entidade (soft delete)
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&UserModel{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return shared.NewDomainError("USER_NOT_FOUND", "user not found")
	}

	return nil
}

// HardDelete remove permanentemente uma entidade
func (r *Repository) HardDelete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&UserModel{})
	if result.Error != nil {
		return fmt.Errorf("failed to hard delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return shared.NewDomainError("USER_NOT_FOUND", "user not found")
	}

	return nil
}

// Restore restaura uma entidade soft deleted
func (r *Repository) Restore(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Unscoped().Model(&UserModel{}).Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return fmt.Errorf("failed to restore user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return shared.NewDomainError("USER_NOT_FOUND", "user not found")
	}

	return nil
}

// ==========================================
// PAGINAÇÃO
// ==========================================

// Paginate retorna resultados paginados com metadados
func (r *Repository) Paginate(ctx context.Context, filter shared.QueryFilter) (*shared.PaginatedResult[*user.User], error) {
	// Validar filtro
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	// Contar total de registros
	var total int64
	countQuery, err := QueryFilterToGORM(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to convert filter: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&UserModel{}).Scopes(countQuery).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Buscar registros
	users, err := r.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	// Criar resultado paginado
	result := shared.NewPaginatedResult(
		users,
		total,
		filter.Page,
		filter.PageSize,
	)

	return result, nil
}

// ==========================================
// AGREGAÇÕES E CONTAGENS
// ==========================================

// Count conta entidades baseado em filtros
func (r *Repository) Count(ctx context.Context, filter shared.QueryFilter) (int64, error) {
	var count int64

	query, err := QueryFilterToGORM(filter)
	if err != nil {
		return 0, fmt.Errorf("failed to convert filter: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&UserModel{}).Scopes(query).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// Exists verifica se existe alguma entidade que corresponde aos filtros
func (r *Repository) Exists(ctx context.Context, filter shared.QueryFilter) (bool, error) {
	count, err := r.Count(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Sum soma valores de um campo específico
func (r *Repository) Sum(ctx context.Context, field string, filter shared.QueryFilter) (float64, error) {
	var result float64

	query, err := QueryFilterToGORM(filter)
	if err != nil {
		return 0, fmt.Errorf("failed to convert filter: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&UserModel{}).Scopes(query).Select(fmt.Sprintf("SUM(%s)", field)).Scan(&result).Error; err != nil {
		return 0, fmt.Errorf("failed to sum field %s: %w", field, err)
	}

	return result, nil
}

// Avg calcula a média de um campo
func (r *Repository) Avg(ctx context.Context, field string, filter shared.QueryFilter) (float64, error) {
	var result float64

	query, err := QueryFilterToGORM(filter)
	if err != nil {
		return 0, fmt.Errorf("failed to convert filter: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&UserModel{}).Scopes(query).Select(fmt.Sprintf("AVG(%s)", field)).Scan(&result).Error; err != nil {
		return 0, fmt.Errorf("failed to avg field %s: %w", field, err)
	}

	return result, nil
}

// Min retorna o valor mínimo de um campo
func (r *Repository) Min(ctx context.Context, field string, filter shared.QueryFilter) (interface{}, error) {
	var result interface{}

	query, err := QueryFilterToGORM(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to convert filter: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&UserModel{}).Scopes(query).Select(fmt.Sprintf("MIN(%s)", field)).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to min field %s: %w", field, err)
	}

	return result, nil
}

// Max retorna o valor máximo de um campo
func (r *Repository) Max(ctx context.Context, field string, filter shared.QueryFilter) (interface{}, error) {
	var result interface{}

	query, err := QueryFilterToGORM(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to convert filter: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&UserModel{}).Scopes(query).Select(fmt.Sprintf("MAX(%s)", field)).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to max field %s: %w", field, err)
	}

	return result, nil
}

// ==========================================
// OPERAÇÕES EM LOTE
// ==========================================

// CreateMany cria múltiplas entidades de uma vez
func (r *Repository) CreateMany(ctx context.Context, entities []*user.User) error {
	if len(entities) == 0 {
		return nil
	}

	models := ToModelSlice(entities)

	if err := r.db.WithContext(ctx).CreateInBatches(models, 100).Error; err != nil {
		return fmt.Errorf("failed to create users: %w", err)
	}

	return nil
}

// UpdateMany atualiza múltiplas entidades
func (r *Repository) UpdateMany(ctx context.Context, filter shared.QueryFilter, updates map[string]interface{}) (int64, error) {
	query, err := QueryFilterToGORM(filter)
	if err != nil {
		return 0, fmt.Errorf("failed to convert filter: %w", err)
	}

	// Adicionar timestamp de atualização
	updates["updated_at"] = time.Now()

	result := r.db.WithContext(ctx).Model(&UserModel{}).Scopes(query).Updates(updates)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to update users: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// DeleteMany remove múltiplas entidades
func (r *Repository) DeleteMany(ctx context.Context, filter shared.QueryFilter) (int64, error) {
	query, err := QueryFilterToGORM(filter)
	if err != nil {
		return 0, fmt.Errorf("failed to convert filter: %w", err)
	}

	result := r.db.WithContext(ctx).Scopes(query).Delete(&UserModel{})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete users: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// ==========================================
// TRANSAÇÕES
// ==========================================

// WithTransaction executa uma função dentro de uma transação
func (r *Repository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Criar novo contexto com transação
		txCtx := context.WithValue(ctx, "tx", tx)
		return fn(txCtx)
	})
}

// ==========================================
// QUERIES AVANÇADAS
// ==========================================

// FindFirst busca a primeira entidade que corresponde aos filtros
func (r *Repository) FindFirst(ctx context.Context, filter shared.QueryFilter) (*user.User, error) {
	var model UserModel

	query, err := QueryFilterToGORM(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to convert filter: %w", err)
	}

	if err := r.db.WithContext(ctx).Scopes(query).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.NewDomainError("USER_NOT_FOUND", "user not found")
		}
		return nil, fmt.Errorf("failed to find first user: %w", err)
	}

	return ToDomain(&model)
}

// FindLast busca a última entidade que corresponde aos filtros
func (r *Repository) FindLast(ctx context.Context, filter shared.QueryFilter) (*user.User, error) {
	var model UserModel

	query, err := QueryFilterToGORM(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to convert filter: %w", err)
	}

	if err := r.db.WithContext(ctx).Scopes(query).Last(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.NewDomainError("USER_NOT_FOUND", "user not found")
		}
		return nil, fmt.Errorf("failed to find last user: %w", err)
	}

	return ToDomain(&model)
}

// Distinct retorna valores únicos de um campo
func (r *Repository) Distinct(ctx context.Context, field string, filter shared.QueryFilter) ([]interface{}, error) {
	var results []interface{}

	query, err := QueryFilterToGORM(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to convert filter: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&UserModel{}).Scopes(query).Distinct(field).Pluck(field, &results).Error; err != nil {
		return nil, fmt.Errorf("failed to get distinct values for field %s: %w", field, err)
	}

	return results, nil
}

// GroupBy agrupa resultados por um campo
func (r *Repository) GroupBy(ctx context.Context, field string, filter shared.QueryFilter) (map[string][]*user.User, error) {
	var results []struct {
		Group string `json:"group"`
		UserModel
	}

	query, err := QueryFilterToGORM(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to convert filter: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&UserModel{}).Scopes(query).Select(field + " as group, *").Group(field).Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to group by field %s: %w", field, err)
	}

	// Agrupar resultados
	groups := make(map[string][]*user.User)
	for _, result := range results {
		user, err := ToDomain(&result.UserModel)
		if err != nil {
			return nil, fmt.Errorf("failed to convert user: %w", err)
		}
		groups[result.Group] = append(groups[result.Group], user)
	}

	return groups, nil
}
