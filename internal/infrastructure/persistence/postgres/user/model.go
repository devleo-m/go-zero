package user

import (
	"time"

	"gorm.io/gorm"
)

// UserModel represents the database table structure for users
// This is the Infrastructure layer - it's tightly coupled to GORM and PostgreSQL
// It's separate from the Domain layer to maintain clean architecture
type UserModel struct {
	// Primary Key and Timestamps
	ID        string         `gorm:"type:varchar(36);primaryKey;not null"`
	CreatedAt time.Time      `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"not null;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete support

	// User Basic Information
	Name     string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(255);not null;uniqueIndex:idx_email"`
	Password string `gorm:"type:varchar(255);not null"` // Stored as hash

	// User Status and Role
	Status     string `gorm:"type:varchar(50);not null;index:idx_status;default:'active'"`
	Role       string `gorm:"type:varchar(50);not null;index:idx_role;default:'client'"`
	IsVerified bool   `gorm:"not null;default:false"`

	// Security and Audit
	LastLoginAt *time.Time `gorm:"type:timestamp"`
}

// TableName specifies the table name for GORM
// This ensures consistency and allows for custom naming conventions
func (UserModel) TableName() string {
	return "users"
}

// BeforeCreate is a GORM hook that runs before inserting a new record
// Used for: auto-generating IDs, setting timestamps, validation
func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	// If ID is not set, generate a new UUID
	if u.ID == "" {
		u.ID = generateID()
	}

	// Set default timestamps if not already set
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = now
	}

	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a record
// Ensures UpdatedAt is always current
func (u *UserModel) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// generateID creates a unique identifier for the user
// Using a simple UUID v4 generation
func generateID() string {
	// Simple UUID generation (in production, use google/uuid or similar)
	return generateUUID()
}

// generateUUID is a helper function to generate UUID v4
// This is a simplified version - in production use a proper UUID library
func generateUUID() string {
	// This is a placeholder - will use proper UUID library
	// For now, using timestamp + random to ensure uniqueness
	return time.Now().Format("20060102150405") + randomString(10)
}

// randomString generates a random string of specified length
// Used for ID generation (temporary solution)
func randomString(length int) string {
	// Simple implementation - in production use crypto/rand
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// IsDeleted checks if the record has been soft deleted
func (u *UserModel) IsDeleted() bool {
	return u.DeletedAt.Valid
}

// UpdateLoginInfo updates login-related fields
// This is a convenience method for tracking user authentication
func (u *UserModel) UpdateLoginInfo() {
	now := time.Now()
	u.LastLoginAt = &now
}

// Scopes are reusable query modifiers
// They follow the GORM Scopes pattern for composable queries

// ActiveUsers returns a scope that filters only active users
func ActiveUsers(db *gorm.DB) *gorm.DB {
	return db.Where("status = ? AND deleted_at IS NULL", "active")
}

// VerifiedUsers returns a scope that filters only email-verified users
func VerifiedUsers(db *gorm.DB) *gorm.DB {
	return db.Where("is_verified = ?", true)
}

// ByRole returns a scope that filters users by role
func ByRole(role string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("role = ?", role)
	}
}

// ByStatus returns a scope that filters users by status
func ByStatus(status string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

// WithEmailLike returns a scope that filters users by email pattern
func WithEmailLike(email string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email ILIKE ?", "%"+email+"%")
	}
}

// WithNameLike returns a scope that filters users by name pattern
func WithNameLike(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name ILIKE ?", "%"+name+"%")
	}
}

// OrderByCreatedAt returns a scope that orders by creation date
func OrderByCreatedAt(desc bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		order := "created_at ASC"
		if desc {
			order = "created_at DESC"
		}
		return db.Order(order)
	}
}
