package dto

import (
	"time"

	"github.com/devleo-m/go-zero/internal/domain/user"
	usecaseuser "github.com/devleo-m/go-zero/internal/usecase/user"
	"github.com/google/uuid"
)

// ==========================================
// REQUEST DTOs (HTTP INPUT)
// ==========================================

// CreateUserRequest representa a requisição HTTP para criar usuário
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100" example:"João Silva"`
	Email    string `json:"email" binding:"required,email" example:"joao@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"MinhaSenh@123"`
	Phone    string `json:"phone,omitempty" binding:"omitempty,min=10,max=15" example:"(11) 99999-9999"`
	Role     string `json:"role" binding:"required,oneof=admin manager user guest" example:"user"`
}

// AuthenticateUserRequest representa a requisição HTTP para autenticação
type AuthenticateUserRequest struct {
	Email    string `json:"email" binding:"required,email" example:"joao@example.com"`
	Password string `json:"password" binding:"required" example:"MinhaSenh@123"`
}

// UpdateUserRequest representa a requisição HTTP para atualizar usuário
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100" example:"João Silva Santos"`
	Phone string `json:"phone,omitempty" binding:"omitempty,min=10,max=15" example:"(11) 88888-8888"`
}

// ChangePasswordRequest representa a requisição HTTP para alterar senha
type ChangePasswordRequest struct {
	OldPassword        string `json:"old_password" binding:"required" example:"MinhaSenh@123"`
	NewPassword        string `json:"new_password" binding:"required,min=8" example:"NovaSenh@456"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required" example:"NovaSenh@456"`
}

// ChangeRoleRequest representa a requisição HTTP para alterar role
type ChangeRoleRequest struct {
	NewRole string `json:"new_role" binding:"required,oneof=admin manager user guest" example:"manager"`
}

// ListUsersRequest representa a requisição HTTP para listar usuários
type ListUsersRequest struct {
	Page   int    `form:"page" binding:"min=1" example:"1"`
	Limit  int    `form:"limit" binding:"min=1,max=100" example:"10"`
	Role   string `form:"role,omitempty" binding:"omitempty,oneof=admin manager user guest" example:"user"`
	Status string `form:"status,omitempty" binding:"omitempty,oneof=active inactive pending suspended" example:"active"`
	Search string `form:"search,omitempty" binding:"omitempty,max=100" example:"joão"`
}

// ==========================================
// RESPONSE DTOs (HTTP OUTPUT)
// ==========================================

// UserResponse representa a resposta HTTP de um usuário
type UserResponse struct {
	ID          uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string     `json:"name" example:"João Silva"`
	Email       string     `json:"email" example:"joao@example.com"`
	Phone       *string    `json:"phone,omitempty" example:"(11) 99999-9999"`
	Role        string     `json:"role" example:"user"`
	Status      string     `json:"status" example:"active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" example:"2024-01-15T10:30:00Z"`
	LoginCount  int        `json:"login_count" example:"5"`
	CreatedAt   time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// CreateUserResponse representa a resposta da criação de usuário
type CreateUserResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"User created successfully"`
	Data    UserResponse `json:"data"`
}

// AuthenticateUserResponse representa a resposta da autenticação
type AuthenticateUserResponse struct {
	Success      bool         `json:"success" example:"true"`
	Message      string       `json:"message" example:"Authentication successful"`
	Data         UserResponse `json:"data"`
	AccessToken  string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int          `json:"expires_in" example:"3600"`
}

// UpdateUserResponse representa a resposta da atualização de usuário
type UpdateUserResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"User updated successfully"`
	Data    UserResponse `json:"data"`
}

// ChangePasswordResponse representa a resposta da alteração de senha
type ChangePasswordResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Password changed successfully"`
}

// ChangeRoleResponse representa a resposta da alteração de role
type ChangeRoleResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"Role changed successfully"`
	Data    UserResponse `json:"data"`
}

// ActivateUserResponse representa a resposta da ativação de usuário
type ActivateUserResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"User activated successfully"`
	Data    UserResponse `json:"data"`
}

// DeactivateUserResponse representa a resposta da desativação de usuário
type DeactivateUserResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"User deactivated successfully"`
	Data    UserResponse `json:"data"`
}

// SuspendUserResponse representa a resposta da suspensão de usuário
type SuspendUserResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"User suspended successfully"`
	Data    UserResponse `json:"data"`
}

// ListUsersResponse representa a resposta da listagem de usuários
type ListUsersResponse struct {
	Success    bool               `json:"success" example:"true"`
	Message    string             `json:"message" example:"Users retrieved successfully"`
	Data       []UserResponse     `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

// GetUserResponse representa a resposta da busca de usuário
type GetUserResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"User retrieved successfully"`
	Data    UserResponse `json:"data"`
}

// ==========================================
// ERROR RESPONSE DTOs
// ==========================================

// ValidationError representa um erro de validação específico
type ValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Invalid email format"`
	Value   string `json:"value,omitempty" example:"invalid-email"`
}

// ==========================================
// CONVERSION FUNCTIONS
// ==========================================

// ToUserResponse converte uma entidade User para UserResponse
func ToUserResponse(u *user.User) UserResponse {
	response := UserResponse{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email.String(),
		Role:        u.Role.String(),
		Status:      u.Status.String(),
		LastLoginAt: u.LastLoginAt,
		LoginCount:  u.LoginCount,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}

	// Adicionar telefone se existir
	if u.Phone != nil {
		phone := u.Phone.String()
		response.Phone = &phone
	}

	return response
}

// ToUserResponseSlice converte uma slice de entidades User para UserResponse
func ToUserResponseSlice(users []*user.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, u := range users {
		responses[i] = ToUserResponse(u)
	}
	return responses
}

// ToUserResponseFromOutput converte UserOutput do use case para UserResponse
func ToUserResponseFromOutput(u usecaseuser.UserOutput) UserResponse {
	response := UserResponse{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Phone:       u.Phone,
		Role:        u.Role,
		Status:      u.Status,
		LastLoginAt: u.LastLoginAt,
		LoginCount:  u.LoginCount,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
	return response
}

// ToUserResponseSliceFromOutput converte uma slice de UserOutput para UserResponse
func ToUserResponseSliceFromOutput(users []usecaseuser.UserOutput) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, u := range users {
		responses[i] = ToUserResponseFromOutput(u)
	}
	return responses
}

// ToCreateUserInput converte CreateUserRequest para CreateUserInput do use case
func (r *CreateUserRequest) ToCreateUserInput() usecaseuser.CreateUserInput {
	return usecaseuser.CreateUserInput{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
		Phone:    r.Phone,
		Role:     r.Role,
	}
}

// ToAuthenticateUserInput converte AuthenticateUserRequest para AuthenticateUserInput do use case
func (r *AuthenticateUserRequest) ToAuthenticateUserInput() usecaseuser.AuthenticateUserInput {
	return usecaseuser.AuthenticateUserInput{
		Email:    r.Email,
		Password: r.Password,
	}
}

// ToUpdateUserInput converte UpdateUserRequest para UpdateUserInput do use case
func (r *UpdateUserRequest) ToUpdateUserInput(userID uuid.UUID) usecaseuser.UpdateUserInput {
	return usecaseuser.UpdateUserInput{
		ID:    userID,
		Name:  r.Name,
		Phone: r.Phone,
	}
}

// ToChangePasswordInput converte ChangePasswordRequest para ChangePasswordInput do use case
func (r *ChangePasswordRequest) ToChangePasswordInput(userID uuid.UUID) usecaseuser.ChangePasswordInput {
	return usecaseuser.ChangePasswordInput{
		UserID:      userID,
		OldPassword: r.OldPassword,
		NewPassword: r.NewPassword,
	}
}

// ToChangeRoleInput converte ChangeRoleRequest para ChangeRoleInput do use case
func (r *ChangeRoleRequest) ToChangeRoleInput(userID, requesterID uuid.UUID) usecaseuser.ChangeRoleInput {
	return usecaseuser.ChangeRoleInput{
		UserID:      userID,
		NewRole:     r.NewRole,
		RequesterID: requesterID,
	}
}

// ToListUsersInput converte ListUsersRequest para ListUsersInput do use case
func (r *ListUsersRequest) ToListUsersInput() usecaseuser.ListUsersInput {
	// Valores padrão se não especificados
	page := r.Page
	if page <= 0 {
		page = 1
	}

	limit := r.Limit
	if limit <= 0 {
		limit = 10
	}

	return usecaseuser.ListUsersInput{
		Page:   page,
		Limit:  limit,
		Role:   r.Role,
		Status: r.Status,
		Search: r.Search,
	}
}
