package user

import (
	"github.com/devleo-m/go-zero/internal/domain/shared"
)

// Erros específicos do domínio User
var (
	// ErrUserNotFound é retornado quando um usuário não é encontrado
	ErrUserNotFound = shared.NewDomainError("USER_NOT_FOUND", "user not found")

	// ErrUserAlreadyExists é retornado quando tentamos criar um usuário que já existe
	ErrUserAlreadyExists = shared.NewDomainError("USER_ALREADY_EXISTS", "user already exists")

	// ErrInvalidCredentials é retornado quando as credenciais são inválidas
	ErrInvalidCredentials = shared.NewDomainError("INVALID_CREDENTIALS", "invalid email or password")

	// ErrUserInactive é retornado quando tentamos fazer login com usuário inativo
	ErrUserInactive = shared.NewDomainError("USER_INACTIVE", "user account is inactive")

	// ErrUserSuspended é retornado quando tentamos fazer login com usuário suspenso
	ErrUserSuspended = shared.NewDomainError("USER_SUSPENDED", "user account is suspended")

	// ErrUserPendingActivation é retornado quando o usuário precisa ser ativado
	ErrUserPendingActivation = shared.NewDomainError("USER_PENDING_ACTIVATION", "user account requires activation")

	// ErrInvalidRole é retornado quando um role inválido é fornecido
	ErrInvalidRole = shared.NewDomainError("INVALID_ROLE", "invalid user role")

	// ErrInvalidStatus é retornado quando um status inválido é fornecido
	ErrInvalidStatus = shared.NewDomainError("INVALID_STATUS", "invalid user status")

	// ErrPasswordTooWeak é retornado quando a senha não atende aos critérios
	ErrPasswordTooWeak = shared.NewDomainError("PASSWORD_TOO_WEAK", "password does not meet security requirements")

	// ErrCannotChangeRole é retornado quando não é possível alterar o role
	ErrCannotChangeRole = shared.NewDomainError("CANNOT_CHANGE_ROLE", "insufficient permissions to change user role")

	// ErrCannotDeactivateSelf é retornado quando um usuário tenta desativar a si mesmo
	ErrCannotDeactivateSelf = shared.NewDomainError("CANNOT_DEACTIVATE_SELF", "cannot deactivate your own account")

	// ErrCannotSuspendSelf é retornado quando um usuário tenta suspender a si mesmo
	ErrCannotSuspendSelf = shared.NewDomainError("CANNOT_SUSPEND_SELF", "cannot suspend your own account")

	// ErrEmailAlreadyInUse é retornado quando o email já está sendo usado por outro usuário
	ErrEmailAlreadyInUse = shared.NewDomainError("EMAIL_ALREADY_IN_USE", "email address is already in use")

	// ErrPhoneAlreadyInUse é retornado quando o telefone já está sendo usado por outro usuário
	ErrPhoneAlreadyInUse = shared.NewDomainError("PHONE_ALREADY_IN_USE", "phone number is already in use")
)

// NewUserNotFoundError cria um erro de usuário não encontrado com detalhes
func NewUserNotFoundError(userID string) shared.DomainError {
	return ErrUserNotFound.WithDetails("user_id: " + userID)
}

// NewUserAlreadyExistsError cria um erro de usuário já existente com detalhes
func NewUserAlreadyExistsError(email string) shared.DomainError {
	return ErrUserAlreadyExists.WithDetails("email: " + email)
}

// NewEmailAlreadyInUseError cria um erro de email já em uso com detalhes
func NewEmailAlreadyInUseError(email string) shared.DomainError {
	return ErrEmailAlreadyInUse.WithDetails("email: " + email)
}

// NewPhoneAlreadyInUseError cria um erro de telefone já em uso com detalhes
func NewPhoneAlreadyInUseError(phone string) shared.DomainError {
	return ErrPhoneAlreadyInUse.WithDetails("phone: " + phone)
}

// NewInvalidRoleError cria um erro de role inválido com detalhes
func NewInvalidRoleError(role string) shared.DomainError {
	return ErrInvalidRole.WithDetails("provided role: " + role)
}

// NewInvalidStatusError cria um erro de status inválido com detalhes
func NewInvalidStatusError(status string) shared.DomainError {
	return ErrInvalidStatus.WithDetails("provided status: " + status)
}

// NewCannotChangeRoleError cria um erro de permissão insuficiente para alterar role
func NewCannotChangeRoleError(currentRole, targetRole string) shared.DomainError {
	return ErrCannotChangeRole.WithDetails("current role: " + currentRole + ", target role: " + targetRole)
}
