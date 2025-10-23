package user

import (
	"time"

	"github.com/devleo-m/go-zero/internal/domain/user"
	"github.com/google/uuid"
)

// ==========================================
// DTOs DE INPUT (REQUEST)
// ==========================================

// CreateUserInput representa os dados de entrada para criar um usuário
type CreateUserInput struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,min=10,max=15"`
	Role     string `json:"role" validate:"required,oneof=admin manager user guest"`
}

// AuthenticateUserInput representa os dados de entrada para autenticação
type AuthenticateUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateUserInput representa os dados de entrada para atualizar um usuário
type UpdateUserInput struct {
	ID    uuid.UUID `json:"id" validate:"required"`
	Name  string    `json:"name" validate:"required,min=2,max=100"`
	Phone string    `json:"phone,omitempty" validate:"omitempty,min=10,max=15"`
}

// ChangePasswordInput representa os dados de entrada para alterar senha
type ChangePasswordInput struct {
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	OldPassword string    `json:"old_password" validate:"required"`
	NewPassword string    `json:"new_password" validate:"required,min=8"`
}

// ActivateUserInput representa os dados de entrada para ativar um usuário
type ActivateUserInput struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// DeactivateUserInput representa os dados de entrada para desativar um usuário
type DeactivateUserInput struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// SuspendUserInput representa os dados de entrada para suspender um usuário
type SuspendUserInput struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// ChangeRoleInput representa os dados de entrada para alterar role
type ChangeRoleInput struct {
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	NewRole     string    `json:"new_role" validate:"required,oneof=admin manager user guest"`
	RequesterID uuid.UUID `json:"requester_id" validate:"required"`
}

// ListUsersInput representa os dados de entrada para listar usuários
type ListUsersInput struct {
	Page   int    `json:"page" validate:"min=1"`
	Limit  int    `json:"limit" validate:"min=1,max=100"`
	Role   string `json:"role,omitempty" validate:"omitempty,oneof=admin manager user guest"`
	Status string `json:"status,omitempty" validate:"omitempty,oneof=active inactive pending suspended"`
	Search string `json:"search,omitempty" validate:"omitempty,max=100"`
}

// GetUserInput representa os dados de entrada para buscar um usuário
type GetUserInput struct {
	ID    uuid.UUID `json:"id,omitempty" validate:"omitempty"`
	Email string    `json:"email,omitempty" validate:"omitempty,email"`
}

// ==========================================
// DTOs DE OUTPUT (RESPONSE)
// ==========================================

// UserOutput representa os dados de saída de um usuário (sem informações sensíveis)
type UserOutput struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Phone       *string    `json:"phone,omitempty"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	LoginCount  int        `json:"login_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateUserOutput representa a saída da criação de usuário
type CreateUserOutput struct {
	User    UserOutput `json:"user"`
	Message string     `json:"message"`
}

// AuthenticateUserOutput representa a saída da autenticação
type AuthenticateUserOutput struct {
	User         UserOutput `json:"user"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresIn    int        `json:"expires_in"`
}

// UpdateUserOutput representa a saída da atualização de usuário
type UpdateUserOutput struct {
	User    UserOutput `json:"user"`
	Message string     `json:"message"`
}

// ChangePasswordOutput representa a saída da alteração de senha
type ChangePasswordOutput struct {
	Message string `json:"message"`
}

// ActivateUserOutput representa a saída da ativação de usuário
type ActivateUserOutput struct {
	User    UserOutput `json:"user"`
	Message string     `json:"message"`
}

// DeactivateUserOutput representa a saída da desativação de usuário
type DeactivateUserOutput struct {
	User    UserOutput `json:"user"`
	Message string     `json:"message"`
}

// SuspendUserOutput representa a saída da suspensão de usuário
type SuspendUserOutput struct {
	User    UserOutput `json:"user"`
	Message string     `json:"message"`
}

// ChangeRoleOutput representa a saída da alteração de role
type ChangeRoleOutput struct {
	User    UserOutput `json:"user"`
	Message string     `json:"message"`
}

// ListUsersOutput representa a saída da listagem de usuários
type ListUsersOutput struct {
	Users      []UserOutput     `json:"users"`
	Pagination PaginationOutput `json:"pagination"`
}

// PaginationOutput representa os metadados de paginação
type PaginationOutput struct {
	CurrentPage int  `json:"current_page"`
	TotalPages  int  `json:"total_pages"`
	PageSize    int  `json:"page_size"`
	TotalItems  int  `json:"total_items"`
	ItemsInPage int  `json:"items_in_page"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}

// GetUserOutput representa a saída da busca de usuário
type GetUserOutput struct {
	User UserOutput `json:"user"`
}

// ==========================================
// DTOs ADICIONAIS
// ==========================================

// ChangePasswordWithConfirmationInput representa input com confirmação de senha
type ChangePasswordWithConfirmationInput struct {
	UserID             uuid.UUID `json:"user_id" validate:"required"`
	OldPassword        string    `json:"old_password" validate:"required"`
	NewPassword        string    `json:"new_password" validate:"required,min=8"`
	ConfirmNewPassword string    `json:"confirm_new_password" validate:"required"`
}

// ChangePasswordWithConfirmationOutput representa output com confirmação
type ChangePasswordWithConfirmationOutput struct {
	Message string `json:"message"`
}

// UpdateUserProfileInput representa input específico para atualização de perfil
type UpdateUserProfileInput struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Name   string    `json:"name" validate:"required,min=2,max=100"`
	Phone  string    `json:"phone,omitempty" validate:"omitempty,min=10,max=15"`
}

// UpdateUserProfileOutput representa output específico para atualização de perfil
type UpdateUserProfileOutput struct {
	User    UserOutput `json:"user"`
	Message string     `json:"message"`
}

// UserBasicInfo representa informações básicas de um usuário (otimizado)
type UserBasicInfo struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	Status string    `json:"status"`
}

// UserStatsOutput representa estatísticas dos usuários
type UserStatsOutput struct {
	Total     int64 `json:"total"`
	Active    int64 `json:"active"`
	Pending   int64 `json:"pending"`
	Suspended int64 `json:"suspended"`
	Inactive  int64 `json:"inactive"`
}

// UserActivityLog representa um log de atividade do usuário
type UserActivityLog struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Activity  string    `json:"activity"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

// ==========================================
// FUNÇÕES DE CONVERSÃO
// ==========================================

// ToUserOutput converte uma entidade User para UserOutput
func ToUserOutput(u *user.User) UserOutput {
	output := UserOutput{
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
		output.Phone = &phone
	}

	return output
}

// ToUserOutputSlice converte uma slice de entidades User para UserOutput
func ToUserOutputSlice(users []*user.User) []UserOutput {
	outputs := make([]UserOutput, len(users))
	for i, u := range users {
		outputs[i] = ToUserOutput(u)
	}
	return outputs
}

// ==========================================
// VALIDAÇÕES CUSTOMIZADAS
// ==========================================

// ValidateCreateUserInput valida os dados de entrada para criação
func ValidateCreateUserInput(input CreateUserInput) error {
	// Validações específicas de negócio podem ser adicionadas aqui
	// As validações básicas são feitas pelas tags de struct

	// Exemplo: validar se o email não está em uso
	// Isso seria feito no use case, não aqui

	return nil
}

// ValidateAuthenticateUserInput valida os dados de entrada para autenticação
func ValidateAuthenticateUserInput(input AuthenticateUserInput) error {
	// Validações específicas de negócio para autenticação
	return nil
}
