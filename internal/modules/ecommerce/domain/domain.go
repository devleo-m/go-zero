package domain

// Este arquivo serve como ponto de entrada para o domain layer
// Aqui podemos exportar os tipos principais para facilitar imports

import (
	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/entities"
	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/errors"
	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/repositories"
	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/services"
	"github.com/devleo-m/go-zero/internal/modules/ecommerce/domain/valueobjects"
)

// Re-exportar apenas os tipos principais (não funções)
type (
	// Entities
	User       = entities.User
	UserRole   = entities.UserRole
	UserStatus = entities.UserStatus

	// Value Objects
	Email    = valueobjects.Email
	Password = valueobjects.Password
	Money    = valueobjects.Money

	// Services
	PasswordService = services.PasswordService

	// Repositories
	UserRepository = repositories.UserRepository
	UserStats      = repositories.UserStats

	// Errors
	DomainError = errors.DomainError
)

// Re-exportar apenas constantes importantes
const (
	// User Roles
	UserRoleAdmin  = entities.UserRoleAdmin
	UserRoleClient = entities.UserRoleClient

	// User Status
	UserStatusActive   = entities.UserStatusActive
	UserStatusInactive = entities.UserStatusInactive
	UserStatusBlocked  = entities.UserStatusBlocked
)

// Re-exportar apenas erros mais comuns
var (
	ErrInvalidEmail      = errors.ErrInvalidEmail
	ErrInvalidPassword   = errors.ErrInvalidPassword
	ErrInvalidMoney      = errors.ErrInvalidMoney
	ErrUserNotFound      = errors.ErrUserNotFound
	ErrUserAlreadyExists = errors.ErrUserAlreadyExists
	ErrUserInactive      = errors.ErrUserInactive
	ErrUserNotVerified   = errors.ErrUserNotVerified
)
