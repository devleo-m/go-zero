package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/devleo-m/go-zero/internal/modules/user/domain"
)

// Repository implementa domain.Repository usando GORM.
type Repository struct {
	db *gorm.DB
}

// NewRepository cria uma nova instância do repositório.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create cria um novo usuário.
func (r *Repository) Create(ctx context.Context, user *domain.User) error {
	model := toModel(user)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Atualizar o ID gerado
	user.ID = model.ID

	return nil
}

// GetByID busca um usuário por ID (excluindo deletados).
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var model UserModel

	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return toDomain(&model), nil
}

// GetByEmail busca um usuário por email (excluindo deletados).
func (r *Repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model UserModel

	if err := r.db.WithContext(ctx).
		Where("email = ? AND deleted_at IS NULL", email).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return toDomain(&model), nil
}

// List lista usuários com paginação (excluindo deletados).
func (r *Repository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var models []UserModel

	if err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Limit(limit).
		Offset(offset).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*domain.User, len(models))
	for i, model := range models {
		users[i] = toDomain(&model)
	}

	return users, nil
}

// Count conta o total de usuários (excluindo deletados).
func (r *Repository) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// Update atualiza um usuário.
func (r *Repository) Update(ctx context.Context, user *domain.User) error {
	model := toModel(user)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deleta um usuário (soft delete).
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()

	err := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("id = ?", id).
		Update("deleted_at", now).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// toModel converte domain.User para UserModel.
func toModel(user *domain.User) *UserModel {
	model := &UserModel{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Phone:     user.Phone,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Converter DeletedAt corretamente
	if user.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{
			Time:  *user.DeletedAt,
			Valid: true,
		}
	}

	return model
}

// toDomain converte UserModel para domain.User.
func toDomain(model *UserModel) *domain.User {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &domain.User{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		Password:  model.Password,
		Phone:     model.Phone,
		Role:      model.Role,
		Status:    model.Status,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: deletedAt,
	}
}
