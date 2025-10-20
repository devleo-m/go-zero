package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		fullName    string
		role        UserRole
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid user",
			email:       "test@example.com",
			password:    "Password123!",
			fullName:    "John Doe",
			role:        UserRoleClient,
			expectError: false,
		},
		{
			name:        "invalid email",
			email:       "invalid-email",
			password:    "Password123!",
			fullName:    "John Doe",
			role:        UserRoleClient,
			expectError: true,
			errorMsg:    "email inválido",
		},
		{
			name:        "invalid password",
			email:       "test@example.com",
			password:    "weak",
			fullName:    "John Doe",
			role:        UserRoleClient,
			expectError: true,
			errorMsg:    "senha inválida",
		},
		{
			name:        "empty full name",
			email:       "test@example.com",
			password:    "Password123!",
			fullName:    "",
			role:        UserRoleClient,
			expectError: true,
			errorMsg:    "nome completo é obrigatório",
		},
		{
			name:        "invalid role",
			email:       "test@example.com",
			password:    "Password123!",
			fullName:    "John Doe",
			role:        UserRole("invalid"),
			expectError: true,
			errorMsg:    "papel do usuário inválido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.password, tt.fullName, tt.role)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email().String())
				assert.Equal(t, tt.fullName, user.FullName())
				assert.Equal(t, tt.role, user.Role())
				assert.Equal(t, UserStatusActive, user.Status())
				assert.False(t, user.IsVerified())
				assert.True(t, user.IsActive())
			}
		})
	}
}

func TestUser_ChangePassword(t *testing.T) {
	tests := []struct {
		name        string
		newPassword string
		expectError bool
	}{
		{
			name:        "valid new password",
			newPassword: "NewPassword123!",
			expectError: false,
		},
		{
			name:        "same password",
			newPassword: "OldPassword123!",
			expectError: true,
		},
		{
			name:        "invalid password",
			newPassword: "weak",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar novo usuário para cada teste (evita estado compartilhado)
			user, err := NewUser("test@example.com", "OldPassword123!", "John Doe", UserRoleClient)
			require.NoError(t, err)

			err = user.ChangePassword(tt.newPassword)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, user.VerifyPassword(tt.newPassword))
			}
		})
	}
}

func TestUser_StatusChanges(t *testing.T) {
	user, err := NewUser("test@example.com", "Password123!", "John Doe", UserRoleClient)
	require.NoError(t, err)

	// Test activation
	assert.True(t, user.IsActive())
	user.Deactivate()
	assert.False(t, user.IsActive())
	assert.Equal(t, UserStatusInactive, user.Status())

	// Test activation
	user.Activate()
	assert.True(t, user.IsActive())
	assert.Equal(t, UserStatusActive, user.Status())

	// Test blocking
	user.Block()
	assert.False(t, user.IsActive())
	assert.Equal(t, UserStatusBlocked, user.Status())

	// Test unblocking
	user.Unblock()
	assert.True(t, user.IsActive())
	assert.Equal(t, UserStatusActive, user.Status())
}

func TestUser_EmailVerification(t *testing.T) {
	user, err := NewUser("test@example.com", "Password123!", "John Doe", UserRoleClient)
	require.NoError(t, err)

	assert.False(t, user.IsVerified())
	assert.False(t, user.CanLogin())

	user.VerifyEmail()
	assert.True(t, user.IsVerified())
	assert.True(t, user.CanLogin())
}

func TestUser_SoftDelete(t *testing.T) {
	user, err := NewUser("test@example.com", "Password123!", "John Doe", UserRoleClient)
	require.NoError(t, err)

	assert.False(t, user.IsDeleted())
	assert.True(t, user.IsActive())

	user.SoftDelete()
	assert.True(t, user.IsDeleted())
	assert.False(t, user.IsActive())

	user.Restore()
	assert.False(t, user.IsDeleted())
	assert.True(t, user.IsActive())
}

func TestUser_RoleChanges(t *testing.T) {
	user, err := NewUser("test@example.com", "Password123!", "John Doe", UserRoleClient)
	require.NoError(t, err)

	assert.True(t, user.IsClient())
	assert.False(t, user.IsAdmin())

	err = user.ChangeRole(UserRoleAdmin)
	assert.NoError(t, err)
	assert.True(t, user.IsAdmin())
	assert.False(t, user.IsClient())

	err = user.ChangeRole(UserRole("invalid"))
	assert.Error(t, err)
}

func TestUser_LoginTracking(t *testing.T) {
	user, err := NewUser("test@example.com", "Password123!", "John Doe", UserRoleClient)
	require.NoError(t, err)

	assert.Nil(t, user.LastLoginAt())

	user.RecordLogin()
	assert.NotNil(t, user.LastLoginAt())
	assert.True(t, time.Since(*user.LastLoginAt()) < time.Second)
}

func TestUser_Equals(t *testing.T) {
	user1, err := NewUser("test1@example.com", "Password123!", "John Doe", UserRoleClient)
	require.NoError(t, err)

	user2, err := NewUser("test2@example.com", "Password123!", "Jane Doe", UserRoleClient)
	require.NoError(t, err)

	assert.True(t, user1.Equals(user1))
	assert.False(t, user1.Equals(user2))
}

func TestUser_ToMap(t *testing.T) {
	user, err := NewUser("test@example.com", "Password123!", "John Doe", UserRoleClient)
	require.NoError(t, err)

	userMap := user.ToMap()

	assert.Equal(t, user.ID(), userMap["id"])
	assert.Equal(t, user.Email().String(), userMap["email"])
	assert.Equal(t, user.FullName(), userMap["full_name"])
	assert.Equal(t, string(user.Role()), userMap["role"])
	assert.Equal(t, string(user.Status()), userMap["status"])
	assert.Equal(t, user.IsVerified(), userMap["is_verified"])
	assert.Equal(t, user.CreatedAt(), userMap["created_at"])
	assert.Equal(t, user.UpdatedAt(), userMap["updated_at"])
}

func TestUser_String(t *testing.T) {
	user, err := NewUser("test@example.com", "Password123!", "John Doe", UserRoleClient)
	require.NoError(t, err)

	assert.Equal(t, "test@example.com", user.String())
}
