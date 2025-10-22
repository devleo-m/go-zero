package user

import (
	"time"

	"github.com/devleo-m/go-zero/internal/domain/shared"
)

// User é a entidade principal do domínio
// Representa um usuário do sistema com todas as regras de negócio
type User struct {
	shared.BaseEntity

	// Informações básicas
	Name     string   `json:"name"`
	Email    Email    `json:"email"`
	Password Password `json:"-"` // Nunca serializar a senha

	// Informações opcionais
	Phone *Phone `json:"phone,omitempty"`

	// Controle de acesso
	Role   Role   `json:"role"`
	Status Status `json:"status"`

	// Metadados
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	LoginCount  int        `json:"login_count"`
}

// NewUser cria um novo usuário com validações de negócio
func NewUser(name string, email Email, password Password, role Role) (*User, error) {
	// Validar role
	if !role.IsValid() {
		return nil, NewInvalidRoleError(role.String())
	}

	// Criar usuário
	user := &User{
		BaseEntity: shared.NewBaseEntity(),
		Name:       name,
		Email:      email,
		Password:   password,
		Role:       role,
		Status:     StatusPending, // Usuário sempre começa como pendente
		LoginCount: 0,
	}

	// Validar usuário
	if err := user.validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// validate valida as regras de negócio do usuário
func (u *User) validate() error {
	if u.Name == "" {
		return shared.NewDomainError("INVALID_USER", "name cannot be empty")
	}

	if len(u.Name) < 2 {
		return shared.NewDomainError("INVALID_USER", "name must be at least 2 characters long")
	}

	if len(u.Name) > 100 {
		return shared.NewDomainError("INVALID_USER", "name too long")
	}

	return nil
}

// Authenticate autentica o usuário verificando a senha
func (u *User) Authenticate(plainPassword string) error {
	// Verificar se usuário está ativo
	if !u.Status.CanLogin() {
		switch u.Status {
		case StatusInactive:
			return ErrUserInactive
		case StatusSuspended:
			return ErrUserSuspended
		case StatusPending:
			return ErrUserPendingActivation
		}
	}

	// Verificar senha
	if !u.Password.Verify(plainPassword) {
		return ErrInvalidCredentials
	}

	// Atualizar informações de login
	now := time.Now()
	u.LastLoginAt = &now
	u.LoginCount++
	u.Touch()

	return nil
}

// ChangePassword altera a senha do usuário
func (u *User) ChangePassword(oldPassword, newPassword string) error {
	// Verificar senha atual
	if !u.Password.Verify(oldPassword) {
		return ErrInvalidCredentials
	}

	// Criar nova senha
	newPasswordVO, err := NewPassword(newPassword)
	if err != nil {
		return err
	}

	// Atualizar senha
	u.Password = newPasswordVO
	u.Touch()

	return nil
}

// UpdateProfile atualiza informações do perfil
func (u *User) UpdateProfile(name string, phone *Phone) error {
	// Validar nome
	if name == "" {
		return shared.NewDomainError("INVALID_USER", "name cannot be empty")
	}

	if len(name) < 2 {
		return shared.NewDomainError("INVALID_USER", "name must be at least 2 characters long")
	}

	if len(name) > 100 {
		return shared.NewDomainError("INVALID_USER", "name too long")
	}

	// Atualizar campos
	u.Name = name
	u.Phone = phone
	u.Touch()

	return nil
}

// Activate ativa o usuário
func (u *User) Activate() error {
	if u.Status != StatusPending {
		return shared.NewDomainError("INVALID_OPERATION", "only pending users can be activated")
	}

	u.Status = StatusActive
	u.Touch()

	return nil
}

// Deactivate desativa o usuário
func (u *User) Deactivate() error {
	if u.Status == StatusInactive {
		return shared.NewDomainError("INVALID_OPERATION", "user is already inactive")
	}

	u.Status = StatusInactive
	u.Touch()

	return nil
}

// Suspend suspende o usuário
func (u *User) Suspend() error {
	if u.Status == StatusSuspended {
		return shared.NewDomainError("INVALID_OPERATION", "user is already suspended")
	}

	u.Status = StatusSuspended
	u.Touch()

	return nil
}

// Unsuspend remove a suspensão do usuário
func (u *User) Unsuspend() error {
	if u.Status != StatusSuspended {
		return shared.NewDomainError("INVALID_OPERATION", "only suspended users can be unsuspended")
	}

	u.Status = StatusActive
	u.Touch()

	return nil
}

// ChangeRole altera o role do usuário
func (u *User) ChangeRole(newRole Role, requesterRole Role) error {
	// Validar novo role
	if !newRole.IsValid() {
		return NewInvalidRoleError(newRole.String())
	}

	// Verificar permissões do solicitante
	if !requesterRole.CanManage(u.Role) {
		return NewCannotChangeRoleError(requesterRole.String(), u.Role.String())
	}

	// Verificar se não está tentando alterar para um role superior
	if !requesterRole.CanManage(newRole) {
		return NewCannotChangeRoleError(requesterRole.String(), newRole.String())
	}

	u.Role = newRole
	u.Touch()

	return nil
}

// CanAccess verifica se o usuário pode acessar um recurso
func (u *User) CanAccess(permission string) bool {
	return u.Role.HasPermission(permission)
}

// CanManage verifica se o usuário pode gerenciar outro usuário
func (u *User) CanManage(targetUser *User) bool {
	return u.Role.CanManage(targetUser.Role)
}

// IsActive verifica se o usuário está ativo
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// IsPending verifica se o usuário está pendente de ativação
func (u *User) IsPending() bool {
	return u.Status == StatusPending
}

// IsSuspended verifica se o usuário está suspenso
func (u *User) IsSuspended() bool {
	return u.Status == StatusSuspended
}

// GetDisplayName retorna o nome para exibição
func (u *User) GetDisplayName() string {
	return u.Name
}

// GetContactInfo retorna informações de contato formatadas
func (u *User) GetContactInfo() map[string]string {
	info := map[string]string{
		"email": u.Email.String(),
	}

	if u.Phone != nil {
		info["phone"] = u.Phone.Format()
	}

	return info
}

// GetRolePermissions retorna as permissões do role do usuário
func (u *User) GetRolePermissions() []string {
	permissions := []string{}

	// Lista de todas as permissões possíveis
	allPermissions := []string{
		"read_users", "create_users", "update_users", "delete_users",
		"read_reports", "manage_appointments", "read_own_profile",
		"update_own_profile", "create_appointments", "read_own_appointments",
		"read_public_info",
	}

	// Filtrar permissões baseadas no role
	for _, perm := range allPermissions {
		if u.Role.HasPermission(perm) {
			permissions = append(permissions, perm)
		}
	}

	return permissions
}
