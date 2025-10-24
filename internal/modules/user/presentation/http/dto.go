package http

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse representa a resposta de um usuário
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone,omitempty"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest representa a requisição de criação de usuário
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Phone    string `json:"phone,omitempty"`
}

// UpdateUserRequest representa a requisição de atualização de usuário
type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Phone string `json:"phone,omitempty"`
}

// ListUsersRequest representa a requisição de listagem de usuários
type ListUsersRequest struct {
	Limit  int `json:"limit" form:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" form:"offset" validate:"min=0"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse representa uma resposta de sucesso
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
