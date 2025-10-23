package user

import (
	"context"
	"time"

	"github.com/devleo-m/go-zero/internal/domain/shared"
	"github.com/devleo-m/go-zero/internal/domain/user"
	"github.com/google/uuid"
)

// userUseCaseAggregate implementa UserUseCaseAggregate interface
type userUseCaseAggregate struct {
	// Casos de uso individuais
	createUserUC       *CreateUserUseCase
	authenticateUserUC *AuthenticateUserUseCase
	getUserUC          *GetUserUseCase
	listUsersUC        *ListUsersUseCase
	updateUserUC       *UpdateUserUseCase
	changePasswordUC   *ChangePasswordUseCase
	activateUserUC     *ActivateUserUseCase
	deactivateUserUC   *DeactivateUserUseCase
	suspendUserUC      *SuspendUserUseCase
	changeRoleUC       *ChangeRoleUseCase

	// Helpers e validators
	queryHelper       *UserQueryHelper
	passwordValidator *PasswordValidator
	nameValidator     *NameValidator
	emailValidator    *EmailValidator
	phoneValidator    *PhoneValidator

	// Dependências compartilhadas
	userRepo     user.Repository
	emailService EmailService
	jwtService   JWTService
	tokenService TokenService
	logger       Logger
}

// NewUserUseCaseAggregate cria uma nova instância do agregado
func NewUserUseCaseAggregate(
	userRepo user.Repository,
	emailService EmailService,
	jwtService JWTService,
	tokenService TokenService,
	logger Logger,
) UserUseCaseAggregate {
	// Criar helpers e validators
	queryHelper := NewUserQueryHelper(userRepo)
	passwordValidator := NewPasswordValidator()
	nameValidator := NewNameValidator()
	emailValidator := NewEmailValidator()
	phoneValidator := NewPhoneValidator()

	// Criar instâncias dos casos de uso individuais
	createUserUC := NewCreateUserUseCase(
		userRepo,
		queryHelper,
		passwordValidator,
		nameValidator,
		emailValidator,
		phoneValidator,
		emailService,
		logger,
	)

	authenticateUserUC := NewAuthenticateUserUseCase(
		userRepo,
		queryHelper,
		jwtService,
		tokenService,
		logger,
	)

	getUserUC := NewGetUserUseCase(
		userRepo,
		queryHelper,
		logger,
	)

	listUsersUC := NewListUsersUseCase(
		userRepo,
		queryHelper,
		logger,
	)

	updateUserUC := NewUpdateUserUseCase(
		userRepo,
		logger,
	)

	changePasswordUC := NewChangePasswordUseCase(
		userRepo,
		logger,
	)

	activateUserUC := NewActivateUserUseCase(
		userRepo,
		emailService,
		logger,
	)

	deactivateUserUC := NewDeactivateUserUseCase(
		userRepo,
		emailService,
		logger,
	)

	suspendUserUC := NewSuspendUserUseCase(
		userRepo,
		emailService,
		logger,
	)

	changeRoleUC := NewChangeRoleUseCase(
		userRepo,
		emailService,
		logger,
	)

	return &userUseCaseAggregate{
		createUserUC:       createUserUC,
		authenticateUserUC: authenticateUserUC,
		getUserUC:          getUserUC,
		listUsersUC:        listUsersUC,
		updateUserUC:       updateUserUC,
		changePasswordUC:   changePasswordUC,
		activateUserUC:     activateUserUC,
		deactivateUserUC:   deactivateUserUC,
		suspendUserUC:      suspendUserUC,
		changeRoleUC:       changeRoleUC,
		queryHelper:        queryHelper,
		passwordValidator:  passwordValidator,
		nameValidator:      nameValidator,
		emailValidator:     emailValidator,
		phoneValidator:     phoneValidator,
		userRepo:           userRepo,
		emailService:       emailService,
		jwtService:         jwtService,
		tokenService:       tokenService,
		logger:             logger,
	}
}

// ==========================================
// MÉTODOS DE DELEGAÇÃO PARA CASOS DE USO
// ==========================================

// CreateUser delega para CreateUserUseCase
func (uc *userUseCaseAggregate) CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	return uc.createUserUC.Execute(ctx, input)
}

// AuthenticateUser delega para AuthenticateUserUseCase
func (uc *userUseCaseAggregate) AuthenticateUser(ctx context.Context, input AuthenticateUserInput) (*AuthenticateUserOutput, error) {
	return uc.authenticateUserUC.Execute(ctx, input)
}

// GetUser delega para GetUserUseCase
func (uc *userUseCaseAggregate) GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	return uc.getUserUC.Execute(ctx, input)
}

// ListUsers delega para ListUsersUseCase
func (uc *userUseCaseAggregate) ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	return uc.listUsersUC.Execute(ctx, input)
}

// UpdateUser delega para UpdateUserUseCase
func (uc *userUseCaseAggregate) UpdateUser(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error) {
	return uc.updateUserUC.Execute(ctx, input)
}

// ChangePassword delega para ChangePasswordUseCase
func (uc *userUseCaseAggregate) ChangePassword(ctx context.Context, input ChangePasswordInput) (*ChangePasswordOutput, error) {
	return uc.changePasswordUC.Execute(ctx, input)
}

// ChangePasswordWithConfirmation delega para ChangePasswordUseCase
func (uc *userUseCaseAggregate) ChangePasswordWithConfirmation(ctx context.Context, input ChangePasswordWithConfirmationInput) (*ChangePasswordWithConfirmationOutput, error) {
	return uc.changePasswordUC.ChangePasswordWithConfirmation(ctx, input)
}

// ActivateUser delega para ActivateUserUseCase
func (uc *userUseCaseAggregate) ActivateUser(ctx context.Context, input ActivateUserInput) (*ActivateUserOutput, error) {
	return uc.activateUserUC.Execute(ctx, input)
}

// DeactivateUser delega para DeactivateUserUseCase
func (uc *userUseCaseAggregate) DeactivateUser(ctx context.Context, input DeactivateUserInput) (*DeactivateUserOutput, error) {
	return uc.deactivateUserUC.Execute(ctx, input)
}

// SuspendUser delega para SuspendUserUseCase
func (uc *userUseCaseAggregate) SuspendUser(ctx context.Context, input SuspendUserInput) (*SuspendUserOutput, error) {
	return uc.suspendUserUC.Execute(ctx, input)
}

// ChangeRole delega para ChangeRoleUseCase
func (uc *userUseCaseAggregate) ChangeRole(ctx context.Context, input ChangeRoleInput) (*ChangeRoleOutput, error) {
	return uc.changeRoleUC.Execute(ctx, input)
}

// ==========================================
// MÉTODOS CONVÊNIENTES E AUXILIARES
// ==========================================

// GetUserByEmail busca usuário por email
func (uc *userUseCaseAggregate) GetUserByEmail(ctx context.Context, email string) (*GetUserOutput, error) {
	input := GetUserInput{
		Email: email,
	}
	return uc.GetUser(ctx, input)
}

// GetUserByID busca usuário por ID
func (uc *userUseCaseAggregate) GetUserByID(ctx context.Context, userID uuid.UUID) (*GetUserOutput, error) {
	input := GetUserInput{
		ID: userID,
	}
	return uc.GetUser(ctx, input)
}

// CheckUserExists verifica se usuário existe
func (uc *userUseCaseAggregate) CheckUserExists(ctx context.Context, userID string) (bool, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return false, err
	}

	_, err = uc.GetUserByID(ctx, parsedID)
	if err != nil {
		if err == user.ErrUserNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetUserBasicInfo retorna informações básicas do usuário
func (uc *userUseCaseAggregate) GetUserBasicInfo(ctx context.Context, userID string) (*UserBasicInfo, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	userOutput, err := uc.GetUserByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	return &UserBasicInfo{
		ID:     userOutput.User.ID,
		Name:   userOutput.User.Name,
		Email:  userOutput.User.Email,
		Role:   userOutput.User.Role,
		Status: userOutput.User.Status,
	}, nil
}

// ListUsersByRole lista usuários por role
func (uc *userUseCaseAggregate) ListUsersByRole(ctx context.Context, role string, page, limit int) (*ListUsersOutput, error) {
	input := ListUsersInput{
		Role:  role,
		Page:  page,
		Limit: limit,
	}
	return uc.ListUsers(ctx, input)
}

// ListUsersByStatus lista usuários por status
func (uc *userUseCaseAggregate) ListUsersByStatus(ctx context.Context, status string, page, limit int) (*ListUsersOutput, error) {
	input := ListUsersInput{
		Status: status,
		Page:   page,
		Limit:  limit,
	}
	return uc.ListUsers(ctx, input)
}

// SearchUsers busca usuários por termo
func (uc *userUseCaseAggregate) SearchUsers(ctx context.Context, searchTerm string, page, limit int) (*ListUsersOutput, error) {
	input := ListUsersInput{
		Search: searchTerm,
		Page:   page,
		Limit:  limit,
	}
	return uc.ListUsers(ctx, input)
}

// GetUserStats retorna estatísticas dos usuários
func (uc *userUseCaseAggregate) GetUserStats(ctx context.Context) (*UserStatsOutput, error) {
	return uc.queryHelper.GetUserStats(ctx)
}

// ==========================================
// MÉTODOS DE VALIDAÇÃO E PERMISSÕES
// ==========================================

// CanUserAccess verifica se usuário pode acessar recurso
func (uc *userUseCaseAggregate) CanUserAccess(ctx context.Context, userID string, resource string) (bool, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return false, err
	}

	userOutput, err := uc.GetUserByID(ctx, parsedID)
	if err != nil {
		return false, err
	}

	// Converter para entidade de domínio para usar métodos de negócio
	emailVO, _ := user.NewEmail(userOutput.User.Email)
	domainUser := &user.User{
		BaseEntity: shared.BaseEntity{ID: userOutput.User.ID},
		Name:       userOutput.User.Name,
		Email:      emailVO,
		Role:       user.Role(userOutput.User.Role),
		Status:     user.Status(userOutput.User.Status),
	}

	return domainUser.CanAccess(resource), nil
}

// CanUserManage verifica se usuário pode gerenciar outro usuário
func (uc *userUseCaseAggregate) CanUserManage(ctx context.Context, managerID, targetID string) (bool, error) {
	managerParsedID, err := uuid.Parse(managerID)
	if err != nil {
		return false, err
	}

	targetParsedID, err := uuid.Parse(targetID)
	if err != nil {
		return false, err
	}

	managerOutput, err := uc.GetUserByID(ctx, managerParsedID)
	if err != nil {
		return false, err
	}

	targetOutput, err := uc.GetUserByID(ctx, targetParsedID)
	if err != nil {
		return false, err
	}

	// Converter para entidades de domínio
	managerEmailVO, _ := user.NewEmail(managerOutput.User.Email)
	managerUser := &user.User{
		BaseEntity: shared.BaseEntity{ID: managerOutput.User.ID},
		Name:       managerOutput.User.Name,
		Email:      managerEmailVO,
		Role:       user.Role(managerOutput.User.Role),
		Status:     user.Status(managerOutput.User.Status),
	}

	targetEmailVO, _ := user.NewEmail(targetOutput.User.Email)
	targetUser := &user.User{
		BaseEntity: shared.BaseEntity{ID: targetOutput.User.ID},
		Name:       targetOutput.User.Name,
		Email:      targetEmailVO,
		Role:       user.Role(targetOutput.User.Role),
		Status:     user.Status(targetOutput.User.Status),
	}

	return managerUser.CanManage(targetUser), nil
}

// ==========================================
// MÉTODOS DE AUDITORIA E LOGS
// ==========================================

// GetUserActivityLog retorna log de atividades do usuário
func (uc *userUseCaseAggregate) GetUserActivityLog(ctx context.Context, userID string, limit int) ([]UserActivityLog, error) {
	// Implementar busca de logs de atividade
	// Por enquanto, retornar slice vazio
	return []UserActivityLog{}, nil
}

// LogUserActivity registra atividade do usuário
func (uc *userUseCaseAggregate) LogUserActivity(ctx context.Context, userID string, activity, details string) error {
	uc.logger.Info("User activity logged",
		"user_id", userID,
		"activity", activity,
		"details", details,
		"timestamp", time.Now(),
	)
	return nil
}

// ==========================================
// MÉTODOS DE LIMPEZA E MANUTENÇÃO
// ==========================================

// CleanupInactiveUsers remove usuários inativos há muito tempo
func (uc *userUseCaseAggregate) CleanupInactiveUsers(ctx context.Context, olderThanDays int) (int, error) {
	// Implementar limpeza de usuários inativos
	// Por enquanto, retornar 0
	return 0, nil
}

// ArchiveOldUsers arquiva usuários antigos
func (uc *userUseCaseAggregate) ArchiveOldUsers(ctx context.Context, olderThanDays int) (int, error) {
	// Implementar arquivamento de usuários antigos
	// Por enquanto, retornar 0
	return 0, nil
}
