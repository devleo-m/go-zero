package user

import (
	"context"
	"fmt"

	"github.com/devleo-m/go-zero/internal/domain/user"
)

// CreateUserUseCase implementa o caso de uso de criação de usuário
type CreateUserUseCase struct {
	userRepo          user.Repository
	queryHelper       *UserQueryHelper
	passwordValidator *PasswordValidator
	nameValidator     *NameValidator
	emailValidator    *EmailValidator
	phoneValidator    *PhoneValidator
	emailService      EmailService
	logger            Logger
}

// NewCreateUserUseCase cria uma nova instância do caso de uso
func NewCreateUserUseCase(
	userRepo user.Repository,
	queryHelper *UserQueryHelper,
	passwordValidator *PasswordValidator,
	nameValidator *NameValidator,
	emailValidator *EmailValidator,
	phoneValidator *PhoneValidator,
	emailService EmailService,
	logger Logger,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo:          userRepo,
		queryHelper:       queryHelper,
		passwordValidator: passwordValidator,
		nameValidator:     nameValidator,
		emailValidator:    emailValidator,
		phoneValidator:    phoneValidator,
		emailService:      emailService,
		logger:            logger,
	}
}

// Execute executa o caso de uso de criação de usuário
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	// 1. VALIDAR INPUT
	if err := uc.validateInput(input); err != nil {
		uc.logger.Warn("Invalid input for create user",
			"error", err,
			"email", input.Email,
		)
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	uc.logger.Debug("Starting create user process",
		"email", input.Email,
		"name", input.Name,
		"role", input.Role,
	)

	// 2. VERIFICAR SE EMAIL JÁ EXISTE
	emailExists, err := uc.queryHelper.CheckEmailExists(ctx, input.Email)
	if err != nil {
		uc.logger.Error("Failed to check email existence",
			"error", err,
			"email", input.Email,
		)
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}

	if emailExists {
		uc.logger.Warn("Email already exists",
			"email", input.Email,
		)
		return nil, user.ErrEmailAlreadyInUse
	}

	// 3. VERIFICAR SE TELEFONE JÁ EXISTE (se fornecido)
	if input.Phone != "" {
		phoneExists, err := uc.queryHelper.CheckPhoneExists(ctx, input.Phone)
		if err != nil {
			uc.logger.Error("Failed to check phone existence",
				"error", err,
				"phone", input.Phone,
			)
			return nil, fmt.Errorf("failed to check phone existence: %w", err)
		}

		if phoneExists {
			uc.logger.Warn("Phone already exists",
				"phone", input.Phone,
			)
			return nil, user.ErrPhoneAlreadyInUse
		}
	}

	// 4. CRIAR ENTIDADE DE DOMÍNIO
	domainUser, err := uc.createDomainUser(input)
	if err != nil {
		uc.logger.Error("Failed to create domain user",
			"error", err,
			"email", input.Email,
		)
		return nil, fmt.Errorf("failed to create domain user: %w", err)
	}

	// 5. SALVAR NO BANCO DE DADOS
	if err := uc.userRepo.Create(ctx, domainUser); err != nil {
		uc.logger.Error("Failed to save user to database",
			"error", err,
			"user_id", domainUser.ID,
			"email", input.Email,
		)
		return nil, fmt.Errorf("failed to save user to database: %w", err)
	}

	uc.logger.Info("User created successfully",
		"user_id", domainUser.ID,
		"email", input.Email,
		"name", input.Name,
		"role", input.Role,
	)

	// 6. ENVIAR EMAIL DE BOAS-VINDAS (ASSÍNCRONO)
	go func() {
		if err := uc.emailService.SendWelcomeEmail(ctx, input.Email, input.Name); err != nil {
			uc.logger.Error("Failed to send welcome email",
				"error", err,
				"user_id", domainUser.ID,
				"email", input.Email,
			)
		} else {
			uc.logger.Info("Welcome email sent successfully",
				"user_id", domainUser.ID,
				"email", input.Email,
			)
		}
	}()

	// 7. RETORNAR RESULTADO
	output := &CreateUserOutput{
		User:    ToUserOutput(domainUser),
		Message: "User created successfully",
	}

	uc.logger.Info("Create user use case completed successfully",
		"user_id", domainUser.ID,
		"email", input.Email,
	)

	return output, nil
}

// validateInput valida os dados de entrada
func (uc *CreateUserUseCase) validateInput(input CreateUserInput) error {
	// Validar nome
	if err := uc.nameValidator.Validate(input.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	// Validar email
	if err := uc.emailValidator.Validate(input.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	// Validar senha
	if err := uc.passwordValidator.Validate(input.Password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	// Validar telefone (se fornecido)
	if input.Phone != "" {
		if err := uc.phoneValidator.Validate(input.Phone); err != nil {
			return fmt.Errorf("invalid phone: %w", err)
		}
	}

	// Validar role
	if !user.Role(input.Role).IsValid() {
		return fmt.Errorf("invalid role: %s", input.Role)
	}

	return nil
}

// createDomainUser cria a entidade de domínio
func (uc *CreateUserUseCase) createDomainUser(input CreateUserInput) (*user.User, error) {
	// Criar value objects
	emailVO, err := user.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to create email value object: %w", err)
	}

	passwordVO, err := user.NewPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create password value object: %w", err)
	}

	// Criar telefone (se fornecido)
	var phoneVO *user.Phone
	if input.Phone != "" {
		phone, err := user.NewPhone(input.Phone)
		if err != nil {
			return nil, fmt.Errorf("failed to create phone value object: %w", err)
		}
		phoneVO = &phone
	}

	// Criar entidade de domínio
	domainUser, err := user.NewUser(
		input.Name,
		emailVO,
		passwordVO,
		user.Role(input.Role),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	// Definir telefone se fornecido
	if phoneVO != nil {
		domainUser.Phone = phoneVO
	}

	return domainUser, nil
}

// ==========================================
// MÉTODOS AUXILIARES
// ==========================================

// CreateUserWithRole cria usuário com role específico
func (uc *CreateUserUseCase) CreateUserWithRole(ctx context.Context, input CreateUserInput, role user.Role) (*CreateUserOutput, error) {
	// Sobrescrever role do input
	input.Role = role.String()

	return uc.Execute(ctx, input)
}

// CreateAdminUser cria usuário administrador
func (uc *CreateUserUseCase) CreateAdminUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	return uc.CreateUserWithRole(ctx, input, user.RoleAdmin)
}

// CreateManagerUser cria usuário gerente
func (uc *CreateUserUseCase) CreateManagerUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	return uc.CreateUserWithRole(ctx, input, user.RoleManager)
}

// CreateRegularUser cria usuário comum
func (uc *CreateUserUseCase) CreateRegularUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	return uc.CreateUserWithRole(ctx, input, user.RoleUser)
}

// CreateGuestUser cria usuário convidado
func (uc *CreateUserUseCase) CreateGuestUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	return uc.CreateUserWithRole(ctx, input, user.RoleGuest)
}
