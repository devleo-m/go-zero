package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/domain/user"
)

// Logger interface para logging
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// Repository implementa a interface user.Repository usando GORM
type Repository struct {
	db     *gorm.DB
	logger Logger
}

// NewRepository cria uma nova instância do repository
func NewRepository(db *gorm.DB, logger Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ==========================================
// MÉTODOS CRUD BÁSICOS
// ==========================================

// Create cria um novo usuário
func (r *Repository) Create(ctx context.Context, entity *user.User) error {
	model := ToModel(entity)

	r.logger.Debug("Creating user",
		"email", model.Email,
		"role", model.Role,
		"status", model.Status,
	)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Error("Failed to create user",
			"error", err,
			"email", model.Email,
		)
		return fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.Info("User created successfully",
		"id", model.ID,
		"email", model.Email,
	)

	return nil
}

// FindByID busca um usuário por ID
func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var model UserModel

	r.logger.Debug("Finding user by ID", "id", id)

	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("User not found", "id", id)
			return nil, nil
		}

		r.logger.Error("Failed to find user by ID",
			"error", err,
			"id", id,
		)
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	domainUser, err := ToDomain(&model)
	if err != nil {
		r.logger.Error("Failed to convert user to domain",
			"error", err,
			"id", id,
		)
		return nil, fmt.Errorf("failed to convert user to domain: %w", err)
	}

	r.logger.Debug("User found successfully", "id", id)
	return domainUser, nil
}

// FindByEmail busca um usuário por email
func (r *Repository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var model UserModel

	r.logger.Debug("Finding user by email", "email", email)

	if err := r.db.WithContext(ctx).First(&model, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("User not found by email", "email", email)
			return nil, nil
		}

		r.logger.Error("Failed to find user by email",
			"error", err,
			"email", email,
		)
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	domainUser, err := ToDomain(&model)
	if err != nil {
		r.logger.Error("Failed to convert user to domain",
			"error", err,
			"email", email,
		)
		return nil, fmt.Errorf("failed to convert user to domain: %w", err)
	}

	r.logger.Debug("User found by email", "email", email)
	return domainUser, nil
}

// Update atualiza um usuário
func (r *Repository) Update(ctx context.Context, entity *user.User) error {
	model := ToModel(entity)

	r.logger.Debug("Updating user",
		"id", model.ID,
		"email", model.Email,
	)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		r.logger.Error("Failed to update user",
			"error", err,
			"id", model.ID,
		)
		return fmt.Errorf("failed to update user: %w", err)
	}

	r.logger.Info("User updated successfully", "id", model.ID)
	return nil
}

// Delete remove um usuário (soft delete)
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("Deleting user", "id", id)

	if err := r.db.WithContext(ctx).Delete(&UserModel{}, "id = ?", id).Error; err != nil {
		r.logger.Error("Failed to delete user",
			"error", err,
			"id", id,
		)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	r.logger.Info("User deleted successfully", "id", id)
	return nil
}

// ==========================================
// MÉTODOS DE BUSCA AVANÇADA
// ==========================================

// FindMany busca múltiplos usuários com filtros
func (r *Repository) FindMany(ctx context.Context, filter shared.QueryFilter) ([]*user.User, error) {
	var models []*UserModel

	r.logger.Debug("Finding users with filter",
		"where_count", len(filter.Where),
		"limit", filter.Limit,
	)

	db, err := r.applyFilter(r.db.WithContext(ctx), filter)
	if err != nil {
		r.logger.Error("Failed to apply filter",
			"error", err,
		)
		return nil, err
	}

	if err := db.Find(&models).Error; err != nil {
		r.logger.Error("Failed to find users",
			"error", err,
		)
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	users, err := ToDomainSlice(models)
	if err != nil {
		r.logger.Error("Failed to convert users to domain",
			"error", err,
		)
		return nil, fmt.Errorf("failed to convert users to domain: %w", err)
	}

	r.logger.Debug("Users found", "count", len(users))
	return users, nil
}

// Count conta usuários com filtros
func (r *Repository) Count(ctx context.Context, filter shared.QueryFilter) (int64, error) {
	var count int64

	r.logger.Debug("Counting users with filter",
		"where_count", len(filter.Where),
	)

	db, err := r.applyFilter(r.db.WithContext(ctx).Model(&UserModel{}), filter)
	if err != nil {
		r.logger.Error("Failed to apply filter for count",
			"error", err,
		)
		return 0, err
	}

	if err := db.Count(&count).Error; err != nil {
		r.logger.Error("Failed to count users",
			"error", err,
		)
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	r.logger.Debug("Users counted", "count", count)
	return count, nil
}

// Paginate busca usuários com paginação
func (r *Repository) Paginate(ctx context.Context, filter shared.QueryFilter) (*shared.PaginatedResult[*user.User], error) {
	// Contar total de registros
	total, err := r.Count(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count users for pagination: %w", err)
	}

	// Buscar registros
	users, err := r.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find users for pagination: %w", err)
	}

	// Calcular metadados de paginação
	page := filter.Page
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20 // default
	}
	if page <= 0 {
		page = 1
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	return &shared.PaginatedResult[*user.User]{
		Data: users,
		Pagination: shared.PaginationMeta{
			CurrentPage: page,
			TotalPages:  totalPages,
			PageSize:    pageSize,
			TotalItems:  total,
			ItemsInPage: len(users),
			HasNext:     hasNext,
			HasPrevious: hasPrev,
		},
	}, nil
}

// ==========================================
// MÉTODOS DE TRANSAÇÃO
// ==========================================

// WithTransaction executa uma função dentro de uma transação
func (r *Repository) WithTransaction(ctx context.Context, fn func(*Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &Repository{
			db:     tx,
			logger: r.logger,
		}
		return fn(txRepo)
	})
}

// ==========================================
// MÉTODOS DE AGREGAÇÃO
// ==========================================

// GetStats retorna estatísticas dos usuários
func (r *Repository) GetStats(ctx context.Context) (*shared.AggregationResult, error) {
	var stats struct {
		Total     int64 `json:"total"`
		Active    int64 `json:"active"`
		Inactive  int64 `json:"inactive"`
		Pending   int64 `json:"pending"`
		Suspended int64 `json:"suspended"`
	}

	r.logger.Debug("Getting user statistics")

	// Total de usuários
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Count(&stats.Total).Error; err != nil {
		r.logger.Error("Failed to count total users", "error", err)
		return nil, fmt.Errorf("failed to count total users: %w", err)
	}

	// Usuários por status
	statusCounts := make(map[string]int64)
	var results []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	if err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error; err != nil {
		r.logger.Error("Failed to count users by status", "error", err)
		return nil, fmt.Errorf("failed to count users by status: %w", err)
	}

	for _, result := range results {
		statusCounts[result.Status] = result.Count
	}

	stats.Active = statusCounts["active"]
	stats.Inactive = statusCounts["inactive"]
	stats.Pending = statusCounts["pending"]
	stats.Suspended = statusCounts["suspended"]

	r.logger.Debug("User statistics retrieved",
		"total", stats.Total,
		"active", stats.Active,
	)

	return &shared.AggregationResult{
		Count: stats.Total,
	}, nil
}

// ==========================================
// MÉTODOS AUXILIARES
// ==========================================

// applyFilter é um helper que aplica QueryFilter ao GORM
func (r *Repository) applyFilter(db *gorm.DB, filter shared.QueryFilter) (*gorm.DB, error) {
	query, err := QueryFilterToGORM(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to convert filter: %w", err)
	}

	return db.Scopes(query), nil
}

// ==========================================
// MÉTODOS DE MANUTENÇÃO
// ==========================================

// CleanupExpiredTokens remove tokens expirados
func (r *Repository) CleanupExpiredTokens(ctx context.Context) error {
	r.logger.Debug("Cleaning up expired tokens")

	now := time.Now()

	// Limpar tokens de reset de senha expirados
	if err := r.db.WithContext(ctx).
		Model(&UserAuthDataModel{}).
		Where("password_reset_expires < ?", now).
		Updates(map[string]interface{}{
			"password_reset_token":   nil,
			"password_reset_expires": nil,
		}).Error; err != nil {
		r.logger.Error("Failed to cleanup password reset tokens", "error", err)
		return fmt.Errorf("failed to cleanup password reset tokens: %w", err)
	}

	// Limpar tokens de ativação expirados
	if err := r.db.WithContext(ctx).
		Model(&UserAuthDataModel{}).
		Where("activation_expires < ?", now).
		Updates(map[string]interface{}{
			"activation_token":   nil,
			"activation_expires": nil,
		}).Error; err != nil {
		r.logger.Error("Failed to cleanup activation tokens", "error", err)
		return fmt.Errorf("failed to cleanup activation tokens: %w", err)
	}

	// Limpar refresh tokens expirados
	if err := r.db.WithContext(ctx).
		Model(&UserAuthDataModel{}).
		Where("refresh_expires < ?", now).
		Updates(map[string]interface{}{
			"refresh_token":      nil,
			"refresh_expires":    nil,
			"refresh_token_hash": nil,
		}).Error; err != nil {
		r.logger.Error("Failed to cleanup refresh tokens", "error", err)
		return fmt.Errorf("failed to cleanup refresh tokens: %w", err)
	}

	r.logger.Info("Expired tokens cleaned up successfully")
	return nil
}
