package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/entities"
	domainErrors "github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/errors"
	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/repositories"
	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/valueobjects"
	"gorm.io/gorm"
)

// UserRepository implements the domain UserRepository interface using GORM
// This is the Infrastructure layer - it knows about GORM and PostgreSQL
// It translates between Domain entities and Database models using the Converter
type UserRepository struct {
	db        *gorm.DB
	converter *UserConverter
}

// NewUserRepository creates a new instance of UserRepository
// This is a factory function that implements Dependency Injection
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserRepository{
		db:        db,
		converter: NewUserConverter(),
	}
}

// Create inserts a new user into the database
// Flow: Domain Entity → Converter → GORM Model → PostgreSQL
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	// 1. Convert Domain entity to GORM model
	model := r.converter.ToModel(user)

	// 2. Insert into database using GORM
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return r.handleError(result.Error, "create user")
	}

	// 3. Update Domain entity with database-generated values
	// This ensures the entity has the correct ID and timestamps
	// Note: In our case, ID is generated in Domain, but this pattern
	// is useful when DB generates IDs (auto-increment, sequences)

	return nil
}

// GetByID retrieves a user by their ID
// Flow: PostgreSQL → GORM Model → Converter → Domain Entity
func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	var model UserModel

	// 1. Query database using GORM
	result := r.db.WithContext(ctx).First(&model, "id = ?", id)

	// 2. Handle not found error
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrUserNotFound
		}
		return nil, r.handleError(result.Error, "find user by ID")
	}

	// 3. Convert GORM model to Domain entity
	user, err := r.converter.ToDomain(&model)
	if err != nil {
		return nil, fmt.Errorf("failed to convert model to domain: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by their email
// This is a common query for authentication
func (r *UserRepository) GetByEmail(ctx context.Context, email valueobjects.Email) (*entities.User, error) {
	var model UserModel

	// 1. Query database using email string
	result := r.db.WithContext(ctx).First(&model, "email = ?", email.String())

	// 2. Handle not found error
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrUserNotFound
		}
		return nil, r.handleError(result.Error, "find user by email")
	}

	// 3. Convert to domain
	user, err := r.converter.ToDomain(&model)
	if err != nil {
		return nil, fmt.Errorf("failed to convert model to domain: %w", err)
	}

	return user, nil
}

// Update modifies an existing user in the database
// Only updates mutable fields, preserving ID and CreatedAt
func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	// 1. First, check if user exists
	var existingModel UserModel
	result := r.db.WithContext(ctx).First(&existingModel, "id = ?", user.ID())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domainErrors.ErrUserNotFound
		}
		return r.handleError(result.Error, "find user for update")
	}

	// 2. Update model with domain data (preserving immutable fields)
	r.converter.UpdateModelFromDomain(&existingModel, user)

	// 3. Save to database
	result = r.db.WithContext(ctx).Save(&existingModel)
	if result.Error != nil {
		return r.handleError(result.Error, "update user")
	}

	return nil
}

// Delete performs a soft delete on a user
// The record remains in database but is marked as deleted
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	// GORM's Delete automatically does soft delete when DeletedAt field exists
	result := r.db.WithContext(ctx).Delete(&UserModel{}, "id = ?", id)

	// Check if record was found
	if result.RowsAffected == 0 {
		return domainErrors.ErrUserNotFound
	}

	if result.Error != nil {
		return r.handleError(result.Error, "delete user")
	}

	return nil
}

// List retrieves users with pagination
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	var userModels []*UserModel

	query := r.db.WithContext(ctx).Model(&UserModel{})

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	// Default ordering
	query = query.Order("created_at DESC")

	// Execute query
	result := query.Find(&userModels)
	if result.Error != nil {
		return nil, r.handleError(result.Error, "list users")
	}

	// Convert to domain entities
	users, err := r.converter.ToDomains(userModels)
	if err != nil {
		return nil, fmt.Errorf("failed to convert models to domain: %w", err)
	}

	return users, nil
}

// Count returns the total number of users
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	result := r.db.WithContext(ctx).Model(&UserModel{}).Count(&count)
	if result.Error != nil {
		return 0, r.handleError(result.Error, "count users")
	}

	return count, nil
}

// ExistsByEmail checks if a user with the given email already exists
// Useful for validation before creating new users
func (r *UserRepository) ExistsByEmail(ctx context.Context, email valueobjects.Email) (bool, error) {
	var count int64

	result := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("email = ?", email.String()).
		Count(&count)

	if result.Error != nil {
		return false, r.handleError(result.Error, "check email exists")
	}

	return count > 0, nil
}

// Exists checks if a user exists by ID
func (r *UserRepository) Exists(ctx context.Context, id string) (bool, error) {
	var count int64

	result := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("id = ?", id).
		Count(&count)

	if result.Error != nil {
		return false, r.handleError(result.Error, "check user exists")
	}

	return count > 0, nil
}

// ListByRole retrieves users by role with pagination
func (r *UserRepository) ListByRole(ctx context.Context, role entities.UserRole, limit, offset int) ([]*entities.User, error) {
	var userModels []*UserModel

	query := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("role = ?", string(role))

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	query = query.Order("created_at DESC")

	result := query.Find(&userModels)
	if result.Error != nil {
		return nil, r.handleError(result.Error, "list users by role")
	}

	users, err := r.converter.ToDomains(userModels)
	if err != nil {
		return nil, fmt.Errorf("failed to convert models to domain: %w", err)
	}

	return users, nil
}

// ListByStatus retrieves users by status with pagination
func (r *UserRepository) ListByStatus(ctx context.Context, status entities.UserStatus, limit, offset int) ([]*entities.User, error) {
	var userModels []*UserModel

	query := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("status = ?", string(status))

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	query = query.Order("created_at DESC")

	result := query.Find(&userModels)
	if result.Error != nil {
		return nil, r.handleError(result.Error, "list users by status")
	}

	users, err := r.converter.ToDomains(userModels)
	if err != nil {
		return nil, fmt.Errorf("failed to convert models to domain: %w", err)
	}

	return users, nil
}

// Search searches for users by name or email
func (r *UserRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entities.User, error) {
	var userModels []*UserModel

	searchPattern := "%" + query + "%"

	q := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	q = q.Order("created_at DESC")

	result := q.Find(&userModels)
	if result.Error != nil {
		return nil, r.handleError(result.Error, "search users")
	}

	users, err := r.converter.ToDomains(userModels)
	if err != nil {
		return nil, fmt.Errorf("failed to convert models to domain: %w", err)
	}

	return users, nil
}

// GetActiveUsers retrieves only active users
func (r *UserRepository) GetActiveUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	var userModels []*UserModel

	query := r.db.WithContext(ctx).
		Scopes(ActiveUsers)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	query = query.Order("created_at DESC")

	result := query.Find(&userModels)
	if result.Error != nil {
		return nil, r.handleError(result.Error, "get active users")
	}

	users, err := r.converter.ToDomains(userModels)
	if err != nil {
		return nil, fmt.Errorf("failed to convert models to domain: %w", err)
	}

	return users, nil
}

// GetVerifiedUsers retrieves only verified users
func (r *UserRepository) GetVerifiedUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	var userModels []*UserModel

	query := r.db.WithContext(ctx).
		Scopes(VerifiedUsers)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	query = query.Order("created_at DESC")

	result := query.Find(&userModels)
	if result.Error != nil {
		return nil, r.handleError(result.Error, "get verified users")
	}

	users, err := r.converter.ToDomains(userModels)
	if err != nil {
		return nil, fmt.Errorf("failed to convert models to domain: %w", err)
	}

	return users, nil
}

// GetUsersCreatedBetween retrieves users created within a date range
func (r *UserRepository) GetUsersCreatedBetween(ctx context.Context, start, end string, limit, offset int) ([]*entities.User, error) {
	var userModels []*UserModel

	query := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("created_at BETWEEN ? AND ?", start, end)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	query = query.Order("created_at DESC")

	result := query.Find(&userModels)
	if result.Error != nil {
		return nil, r.handleError(result.Error, "get users created between dates")
	}

	users, err := r.converter.ToDomains(userModels)
	if err != nil {
		return nil, fmt.Errorf("failed to convert models to domain: %w", err)
	}

	return users, nil
}

// GetUsersByLastLogin retrieves users ordered by last login
func (r *UserRepository) GetUsersByLastLogin(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	var userModels []*UserModel

	query := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("last_login_at IS NOT NULL")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	query = query.Order("last_login_at DESC")

	result := query.Find(&userModels)
	if result.Error != nil {
		return nil, r.handleError(result.Error, "get users by last login")
	}

	users, err := r.converter.ToDomains(userModels)
	if err != nil {
		return nil, fmt.Errorf("failed to convert models to domain: %w", err)
	}

	return users, nil
}

// BulkUpdate updates multiple users in a transaction
func (r *UserRepository) BulkUpdate(ctx context.Context, users []*entities.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, user := range users {
			var existingModel UserModel
			if err := tx.First(&existingModel, "id = ?", user.ID()).Error; err != nil {
				return err
			}

			r.converter.UpdateModelFromDomain(&existingModel, user)

			if err := tx.Save(&existingModel).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BulkDelete soft deletes multiple users
func (r *UserRepository) BulkDelete(ctx context.Context, ids []string) error {
	result := r.db.WithContext(ctx).Delete(&UserModel{}, ids)
	if result.Error != nil {
		return r.handleError(result.Error, "bulk delete users")
	}
	return nil
}

// Restore restores a soft-deleted user
func (r *UserRepository) Restore(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil)

	if result.Error != nil {
		return r.handleError(result.Error, "restore user")
	}

	if result.RowsAffected == 0 {
		return domainErrors.ErrUserNotFound
	}

	return nil
}

// GetDeletedUsers retrieves soft-deleted users
func (r *UserRepository) GetDeletedUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	var userModels []*UserModel

	query := r.db.WithContext(ctx).
		Unscoped().
		Model(&UserModel{}).
		Where("deleted_at IS NOT NULL")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	query = query.Order("deleted_at DESC")

	result := query.Find(&userModels)
	if result.Error != nil {
		return nil, r.handleError(result.Error, "get deleted users")
	}

	users, err := r.converter.ToDomains(userModels)
	if err != nil {
		return nil, fmt.Errorf("failed to convert models to domain: %w", err)
	}

	return users, nil
}

// HardDelete permanently deletes a user from database
func (r *UserRepository) HardDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Unscoped().
		Delete(&UserModel{}, "id = ?", id)

	if result.Error != nil {
		return r.handleError(result.Error, "hard delete user")
	}

	if result.RowsAffected == 0 {
		return domainErrors.ErrUserNotFound
	}

	return nil
}

// GetUserStats returns statistics about users
func (r *UserRepository) GetUserStats(ctx context.Context) (*repositories.UserStats, error) {
	stats := &repositories.UserStats{}

	// Total users
	r.db.WithContext(ctx).Model(&UserModel{}).Count(&stats.TotalUsers)

	// Active users
	r.db.WithContext(ctx).Model(&UserModel{}).
		Where("status = ?", string(entities.UserStatusActive)).
		Count(&stats.ActiveUsers)

	// Inactive users
	r.db.WithContext(ctx).Model(&UserModel{}).
		Where("status = ?", string(entities.UserStatusInactive)).
		Count(&stats.InactiveUsers)

	// Blocked users
	r.db.WithContext(ctx).Model(&UserModel{}).
		Where("status = ?", string(entities.UserStatusBlocked)).
		Count(&stats.BlockedUsers)

	// Verified users
	r.db.WithContext(ctx).Model(&UserModel{}).
		Where("is_verified = ?", true).
		Count(&stats.VerifiedUsers)

	// Unverified users
	r.db.WithContext(ctx).Model(&UserModel{}).
		Where("is_verified = ?", false).
		Count(&stats.UnverifiedUsers)

	// Admin users
	r.db.WithContext(ctx).Model(&UserModel{}).
		Where("role = ?", string(entities.UserRoleAdmin)).
		Count(&stats.AdminUsers)

	// Client users
	r.db.WithContext(ctx).Model(&UserModel{}).
		Where("role = ?", string(entities.UserRoleClient)).
		Count(&stats.ClientUsers)

	// Deleted users
	r.db.WithContext(ctx).Unscoped().Model(&UserModel{}).
		Where("deleted_at IS NOT NULL").
		Count(&stats.DeletedUsers)

	return stats, nil
}

// handleError converts GORM/database errors to Domain errors
// This maintains clean architecture by not exposing infrastructure errors to domain
func (r *UserRepository) handleError(err error, operation string) error {
	// Handle not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domainErrors.ErrUserNotFound
	}

	// Handle unique constraint violations
	if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
		if strings.Contains(err.Error(), "email") {
			return domainErrors.ErrEmailAlreadyInUse
		}
		// Could add more specific constraint handling here
		return domainErrors.NewDomainError(
			"DUPLICATE_ENTRY",
			"registro duplicado no banco de dados",
			map[string]interface{}{"operation": operation, "error": err.Error()},
		)
	}

	// Handle foreign key violations
	if strings.Contains(err.Error(), "foreign key") {
		return domainErrors.NewDomainError(
			"FOREIGN_KEY_VIOLATION",
			"violação de chave estrangeira",
			map[string]interface{}{"operation": operation, "error": err.Error()},
		)
	}

	// Generic database error (don't expose internal details to domain)
	return domainErrors.NewDomainError(
		"DATABASE_ERROR",
		fmt.Sprintf("erro ao %s", operation),
		map[string]interface{}{"operation": operation},
	)
}
