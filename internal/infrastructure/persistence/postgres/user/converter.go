package user

import (
	"fmt"
	"time"

	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/entities"
	"gorm.io/gorm"
)

// UserConverter handles bidirectional conversion between Domain and Infrastructure layers
// This is the TRANSLATOR between two different languages:
// - Domain Entity (with Value Objects and business logic)
// - GORM Model (simple struct for database persistence)
type UserConverter struct{}

// NewUserConverter creates a new instance of UserConverter
func NewUserConverter() *UserConverter {
	return &UserConverter{}
}

// ToModel converts a Domain User entity to a GORM UserModel
// Direction: Domain → Infrastructure (for saving to database)
//
// Process:
// 1. Extract primitive values from Value Objects
// 2. Map Domain types to Database types
// 3. Handle nullable fields (pointers)
func (c *UserConverter) ToModel(user *entities.User) *UserModel {
	model := &UserModel{
		ID:         user.ID(),
		CreatedAt:  user.CreatedAt(),
		UpdatedAt:  user.UpdatedAt(),
		Name:       user.FullName(),
		Email:      user.Email().String(),  // Email VO → string
		Password:   user.Password().Hash(), // Password VO → hash string
		Status:     string(user.Status()),  // UserStatus → string
		Role:       string(user.Role()),    // UserRole → string
		IsVerified: user.IsVerified(),
	}

	// Handle nullable timestamp fields
	if user.LastLoginAt() != nil {
		model.LastLoginAt = user.LastLoginAt()
	}

	return model
}

// ToDomain converts a GORM UserModel to a Domain User entity
// Direction: Infrastructure → Domain (after fetching from database)
//
// Process:
// 1. Reconstruct Value Objects with validation
// 2. Map Database types to Domain types
// 3. Handle nullable fields safely
//
// Returns error if Value Object validation fails
func (c *UserConverter) ToDomain(model *UserModel) (*entities.User, error) {
	// Parse Status enum
	status, err := c.parseUserStatus(model.Status)
	if err != nil {
		return nil, fmt.Errorf("invalid status in database: %w", err)
	}

	// Parse Role enum
	role, err := c.parseUserRole(model.Role)
	if err != nil {
		return nil, fmt.Errorf("invalid role in database: %w", err)
	}

	// Create User entity using the constructor for reconstruction from database
	// Note: We're using NewUserFromData because data is already validated in DB
	user, err := entities.NewUserFromData(
		model.ID,
		model.Email,
		model.Password,
		model.Name,
		role,
		status,
		model.IsVerified,
		model.LastLoginAt,
		model.CreatedAt,
		model.UpdatedAt,
		c.getDeletedAt(model.DeletedAt),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct user from database: %w", err)
	}

	return user, nil
}

// ToDomains converts multiple GORM UserModels to Domain User entities
// Useful for batch operations and list queries
func (c *UserConverter) ToDomains(models []*UserModel) ([]*entities.User, error) {
	users := make([]*entities.User, 0, len(models))

	for _, model := range models {
		user, err := c.ToDomain(model)
		if err != nil {
			// Log error but continue processing
			// In production, you might want to collect errors
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateModelFromDomain updates an existing GORM model with data from Domain entity
// Used in UPDATE operations to preserve database-generated fields
//
// Only updates mutable fields, preserving:
// - ID (immutable)
// - CreatedAt (never changes)
// - System-managed fields
func (c *UserConverter) UpdateModelFromDomain(model *UserModel, user *entities.User) {
	// Update mutable fields
	model.Name = user.FullName()
	model.Email = user.Email().String()
	model.Password = user.Password().Hash()
	model.Status = string(user.Status())
	model.Role = string(user.Role())
	model.IsVerified = user.IsVerified()
	model.UpdatedAt = time.Now()

	// Update nullable timestamps
	if user.LastLoginAt() != nil {
		model.LastLoginAt = user.LastLoginAt()
	}
}

// Helper Methods

// parseUserStatus converts string to UserStatus enum
func (c *UserConverter) parseUserStatus(status string) (entities.UserStatus, error) {
	switch status {
	case string(entities.UserStatusActive):
		return entities.UserStatusActive, nil
	case string(entities.UserStatusInactive):
		return entities.UserStatusInactive, nil
	case string(entities.UserStatusBlocked):
		return entities.UserStatusBlocked, nil
	default:
		return "", fmt.Errorf("unknown user status: %s", status)
	}
}

// parseUserRole converts string to UserRole enum
func (c *UserConverter) parseUserRole(role string) (entities.UserRole, error) {
	switch role {
	case string(entities.UserRoleAdmin):
		return entities.UserRoleAdmin, nil
	case string(entities.UserRoleClient):
		return entities.UserRoleClient, nil
	default:
		return "", fmt.Errorf("unknown user role: %s", role)
	}
}

// getDeletedAt safely extracts deleted_at from GORM's DeletedAt type
func (c *UserConverter) getDeletedAt(deletedAt gorm.DeletedAt) *time.Time {
	if deletedAt.Valid {
		return &deletedAt.Time
	}
	return nil
}
