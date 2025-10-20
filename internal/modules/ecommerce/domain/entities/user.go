package entities

import (
	"time"

	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/errors"
	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/valueobjects"
	"github.com/google/uuid"
)

// UserRole representa o papel do usuário no sistema
type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleClient UserRole = "client"
)

// UserStatus representa o status do usuário
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBlocked  UserStatus = "blocked"
)

// User representa um usuário do sistema
type User struct {
	id          string
	email       valueobjects.Email
	password    valueobjects.Password
	fullName    string
	role        UserRole
	status      UserStatus
	isVerified  bool
	lastLoginAt *time.Time
	createdAt   time.Time
	updatedAt   time.Time
	deletedAt   *time.Time
}

// NewUser cria um novo usuário
func NewUser(email, password, fullName string, role UserRole) (*User, error) {
	// Validar email
	emailVO, err := valueobjects.NewEmail(email)
	if err != nil {
		return nil, err
	}

	// Validar senha
	passwordVO, err := valueobjects.NewPassword(password)
	if err != nil {
		return nil, err
	}

	// Validar nome
	if fullName == "" {
		return nil, errors.NewDomainError("INVALID_FULL_NAME", "nome completo é obrigatório", nil)
	}

	// Validar role
	if role != UserRoleAdmin && role != UserRoleClient {
		return nil, errors.NewDomainError("INVALID_ROLE", "papel do usuário inválido", nil)
	}

	now := time.Now()
	return &User{
		id:         uuid.New().String(),
		email:      *emailVO,
		password:   *passwordVO,
		fullName:   fullName,
		role:       role,
		status:     UserStatusActive,
		isVerified: false,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

// NewUserFromData cria um usuário a partir de dados existentes (para reconstruir do banco)
func NewUserFromData(
	id, email, passwordHash, fullName string,
	role UserRole,
	status UserStatus,
	isVerified bool,
	lastLoginAt *time.Time,
	createdAt, updatedAt time.Time,
	deletedAt *time.Time,
) (*User, error) {
	// Validar email
	emailVO, err := valueobjects.NewEmail(email)
	if err != nil {
		return nil, err
	}

	// Validar senha (hash existente)
	passwordVO, err := valueobjects.NewPasswordFromHash(passwordHash)
	if err != nil {
		return nil, err
	}

	// Validar ID
	if id == "" {
		return nil, errors.NewDomainError("INVALID_ID", "ID é obrigatório", nil)
	}

	// Validar nome
	if fullName == "" {
		return nil, errors.NewDomainError("INVALID_FULL_NAME", "nome completo é obrigatório", nil)
	}

	return &User{
		id:          id,
		email:       *emailVO,
		password:    *passwordVO,
		fullName:    fullName,
		role:        role,
		status:      status,
		isVerified:  isVerified,
		lastLoginAt: lastLoginAt,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		deletedAt:   deletedAt,
	}, nil
}

// Getters
func (u *User) ID() string {
	return u.id
}

func (u *User) Email() valueobjects.Email {
	return u.email
}

func (u *User) Password() valueobjects.Password {
	return u.password
}

func (u *User) FullName() string {
	return u.fullName
}

func (u *User) Role() UserRole {
	return u.role
}

func (u *User) Status() UserStatus {
	return u.status
}

func (u *User) IsVerified() bool {
	return u.isVerified
}

func (u *User) LastLoginAt() *time.Time {
	return u.lastLoginAt
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) DeletedAt() *time.Time {
	return u.deletedAt
}

// Business Methods

// IsActive verifica se o usuário está ativo
func (u *User) IsActive() bool {
	return u.status == UserStatusActive && u.deletedAt == nil
}

// IsAdmin verifica se o usuário é administrador
func (u *User) IsAdmin() bool {
	return u.role == UserRoleAdmin
}

// IsClient verifica se o usuário é cliente
func (u *User) IsClient() bool {
	return u.role == UserRoleClient
}

// CanLogin verifica se o usuário pode fazer login
func (u *User) CanLogin() bool {
	return u.IsActive() && u.isVerified
}

// VerifyPassword verifica se a senha está correta
func (u *User) VerifyPassword(password string) bool {
	return u.password.Verify(password)
}

// ChangePassword altera a senha do usuário
func (u *User) ChangePassword(newPassword string) error {
	// Verificar se não é a mesma senha (usar verificação, não comparação de hash)
	if u.password.Verify(newPassword) {
		return errors.NewDomainError("SAME_PASSWORD", "nova senha deve ser diferente da atual", nil)
	}

	// Validar nova senha
	passwordVO, err := valueobjects.NewPassword(newPassword)
	if err != nil {
		return err
	}

	u.password = *passwordVO
	u.updatedAt = time.Now()
	return nil
}

// UpdateFullName atualiza o nome completo
func (u *User) UpdateFullName(fullName string) error {
	if fullName == "" {
		return errors.NewDomainError("INVALID_FULL_NAME", "nome completo é obrigatório", nil)
	}

	u.fullName = fullName
	u.updatedAt = time.Now()
	return nil
}

// VerifyEmail marca o email como verificado
func (u *User) VerifyEmail() {
	u.isVerified = true
	u.updatedAt = time.Now()
}

// Activate ativa o usuário
func (u *User) Activate() {
	u.status = UserStatusActive
	u.updatedAt = time.Now()
}

// Deactivate desativa o usuário
func (u *User) Deactivate() {
	u.status = UserStatusInactive
	u.updatedAt = time.Now()
}

// Block bloqueia o usuário
func (u *User) Block() {
	u.status = UserStatusBlocked
	u.updatedAt = time.Now()
}

// Unblock desbloqueia o usuário
func (u *User) Unblock() {
	u.status = UserStatusActive
	u.updatedAt = time.Now()
}

// RecordLogin registra o último login
func (u *User) RecordLogin() {
	now := time.Now()
	u.lastLoginAt = &now
	u.updatedAt = now
}

// SoftDelete marca o usuário como deletado (soft delete)
func (u *User) SoftDelete() {
	now := time.Now()
	u.deletedAt = &now
	u.updatedAt = now
}

// Restore restaura um usuário deletado
func (u *User) Restore() {
	u.deletedAt = nil
	u.updatedAt = time.Now()
}

// IsDeleted verifica se o usuário foi deletado
func (u *User) IsDeleted() bool {
	return u.deletedAt != nil
}

// ChangeRole altera o papel do usuário
func (u *User) ChangeRole(newRole UserRole) error {
	if newRole != UserRoleAdmin && newRole != UserRoleClient {
		return errors.NewDomainError("INVALID_ROLE", "papel do usuário inválido", nil)
	}

	u.role = newRole
	u.updatedAt = time.Now()
	return nil
}

// Equals verifica se dois usuários são iguais
func (u *User) Equals(other *User) bool {
	return u.id == other.id
}

// ToMap converte o usuário para um mapa (para serialização)
func (u *User) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"id":          u.id,
		"email":       u.email.String(),
		"full_name":   u.fullName,
		"role":        string(u.role),
		"status":      string(u.status),
		"is_verified": u.isVerified,
		"created_at":  u.createdAt,
		"updated_at":  u.updatedAt,
	}

	if u.lastLoginAt != nil {
		result["last_login_at"] = *u.lastLoginAt
	}

	if u.deletedAt != nil {
		result["deleted_at"] = *u.deletedAt
	}

	return result
}

// String retorna representação string do usuário
func (u *User) String() string {
	return u.email.String()
}
