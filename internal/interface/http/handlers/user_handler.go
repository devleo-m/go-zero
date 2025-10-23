package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/devleo-m/go-zero/internal/interface/http/dto"
	"github.com/devleo-m/go-zero/internal/interface/http/middleware"
	"github.com/devleo-m/go-zero/internal/interface/http/validation"
	"github.com/devleo-m/go-zero/internal/usecase/user"
)

// UserHandler handler para operações de usuário
type UserHandler struct {
	useCaseAggregate user.UserUseCaseAggregate
	validator        *validation.CustomValidator
	errorHandler     *ErrorHandler
	logger           *zap.Logger
}

// NewUserHandler cria uma nova instância do UserHandler
func NewUserHandler(
	useCaseAggregate user.UserUseCaseAggregate,
	validator *validation.CustomValidator,
	errorHandler *ErrorHandler,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		useCaseAggregate: useCaseAggregate,
		validator:        validator,
		errorHandler:     errorHandler,
		logger:           logger,
	}
}

// GetUseCaseAggregate retorna o use case aggregate
func (h *UserHandler) GetUseCaseAggregate() user.UserUseCaseAggregate {
	return h.useCaseAggregate
}

// ==========================================
// USER CREATION ENDPOINTS
// ==========================================

// CreateUser cria um novo usuário
// @Summary Create a new user
// @Description Create a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User data"
// @Success 201 {object} dto.CreateUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	// Bind e validar request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Validar com validador customizado
	if err := h.validator.ValidateStruct(req); err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para input do use case
	input := req.ToCreateUserInput()

	// Executar use case
	output, err := h.useCaseAggregate.CreateUser(c.Request.Context(), input)
	if err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.CreateUserResponse{
		Success: true,
		Message: output.Message,
		Data:    dto.ToUserResponseFromOutput(output.User),
	}

	// Log do evento
	h.logger.Info("User created successfully",
		zap.String("user_id", output.User.ID.String()),
		zap.String("email", output.User.Email),
		zap.String("role", output.User.Role),
	)

	c.JSON(http.StatusCreated, response)
}

// ==========================================
// AUTHENTICATION ENDPOINTS
// ==========================================

// AuthenticateUser autentica um usuário
// @Summary Authenticate user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.AuthenticateUserRequest true "Login credentials"
// @Success 200 {object} dto.AuthenticateUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *UserHandler) AuthenticateUser(c *gin.Context) {
	var req dto.AuthenticateUserRequest

	// Bind e validar request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Validar com validador customizado
	if err := h.validator.ValidateStruct(req); err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para input do use case
	input := req.ToAuthenticateUserInput()

	// Executar use case
	output, err := h.useCaseAggregate.AuthenticateUser(c.Request.Context(), input)
	if err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.AuthenticateUserResponse{
		Success:      true,
		Message:      "Authentication successful",
		Data:         dto.ToUserResponseFromOutput(output.User),
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		ExpiresIn:    output.ExpiresIn,
	}

	// Log do evento
	h.logger.Info("User authenticated successfully",
		zap.String("user_id", output.User.ID.String()),
		zap.String("email", output.User.Email),
		zap.String("ip", c.ClientIP()),
	)

	c.JSON(http.StatusOK, response)
}

// ==========================================
// USER RETRIEVAL ENDPOINTS
// ==========================================

// GetUser busca um usuário por ID
// @Summary Get user by ID
// @Description Get user information by user ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.GetUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	// Extrair ID da URL
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.errorHandler.HandleValidationError(c, "id", "Invalid user ID format")
		return
	}

	// Criar input do use case
	input := user.GetUserInput{
		ID: userID,
	}

	// Executar use case
	output, err := h.useCaseAggregate.GetUser(c.Request.Context(), input)
	if err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.GetUserResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    dto.ToUserResponseFromOutput(output.User),
	}

	c.JSON(http.StatusOK, response)
}

// GetUserByEmail busca um usuário por email
// @Summary Get user by email
// @Description Get user information by email address
// @Tags users
// @Produce json
// @Param email query string true "User email"
// @Success 200 {object} dto.GetUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/email/{email} [get]
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	// Extrair email da URL
	email := c.Param("email")
	if email == "" {
		h.errorHandler.HandleValidationError(c, "email", "Email is required")
		return
	}

	// Executar use case
	output, err := h.useCaseAggregate.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.GetUserResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    dto.ToUserResponseFromOutput(output.User),
	}

	c.JSON(http.StatusOK, response)
}

// ListUsers lista usuários com paginação e filtros
// @Summary List users
// @Description Get a paginated list of users with optional filters
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param role query string false "Filter by role" Enums(admin,manager,user,guest)
// @Param status query string false "Filter by status" Enums(active,inactive,pending,suspended)
// @Param search query string false "Search term"
// @Success 200 {object} dto.ListUsersResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req dto.ListUsersRequest

	// Bind query parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Validar com validador customizado
	if err := h.validator.ValidateStruct(req); err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para input do use case
	input := req.ToListUsersInput()

	// Executar use case
	output, err := h.useCaseAggregate.ListUsers(c.Request.Context(), input)
	if err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.ListUsersResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    dto.ToUserResponseSliceFromOutput(output.Users),
		Pagination: dto.PaginationResponse{
			CurrentPage: output.Pagination.CurrentPage,
			TotalPages:  output.Pagination.TotalPages,
			PageSize:    output.Pagination.PageSize,
			TotalItems:  output.Pagination.TotalItems,
			ItemsInPage: output.Pagination.ItemsInPage,
			HasNext:     output.Pagination.HasNext,
			HasPrevious: output.Pagination.HasPrevious,
		},
	}

	c.JSON(http.StatusOK, response)
}

// ==========================================
// USER UPDATE ENDPOINTS
// ==========================================

// UpdateUser atualiza um usuário
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dto.UpdateUserRequest true "User data"
// @Success 200 {object} dto.UpdateUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Extrair ID da URL
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.errorHandler.HandleValidationError(c, "id", "Invalid user ID format")
		return
	}

	var req dto.UpdateUserRequest

	// Bind e validar request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Validar com validador customizado
	if err := h.validator.ValidateStruct(req); err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para input do use case
	input := req.ToUpdateUserInput(userID)

	// Executar use case
	output, err := h.useCaseAggregate.UpdateUser(c.Request.Context(), input)
	if err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.UpdateUserResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    dto.ToUserResponseFromOutput(output.User),
	}

	// Log do evento
	updatedBy, _ := middleware.GetUserIDFromContext(c)
	h.logger.Info("User updated successfully",
		zap.String("user_id", userID.String()),
		zap.String("updated_by", updatedBy),
	)

	c.JSON(http.StatusOK, response)
}

// ==========================================
// PASSWORD MANAGEMENT ENDPOINTS
// ==========================================

// ChangePassword altera a senha do usuário
// @Summary Change user password
// @Description Change user password with old password verification
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param password body dto.ChangePasswordRequest true "Password data"
// @Success 200 {object} dto.ChangePasswordResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/{id}/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// Extrair ID da URL
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.errorHandler.HandleValidationError(c, "id", "Invalid user ID format")
		return
	}

	var req dto.ChangePasswordRequest

	// Bind e validar request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Validar com validador customizado
	if err := h.validator.ValidateStruct(req); err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Verificar se as senhas coincidem
	if req.NewPassword != req.ConfirmNewPassword {
		h.errorHandler.HandleValidationError(c, "confirm_new_password", "Passwords do not match")
		return
	}

	// Converter para input do use case
	input := req.ToChangePasswordInput(userID)

	// Executar use case
	_, err = h.useCaseAggregate.ChangePassword(c.Request.Context(), input)
	if err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.ChangePasswordResponse{
		Success: true,
		Message: "Password changed successfully",
	}

	// Log do evento
	changedBy, _ := middleware.GetUserIDFromContext(c)
	h.logger.Info("User password changed successfully",
		zap.String("user_id", userID.String()),
		zap.String("changed_by", changedBy),
	)

	c.JSON(http.StatusOK, response)
}

// ==========================================
// USER STATUS MANAGEMENT ENDPOINTS
// ==========================================

// ActivateUser ativa um usuário
// @Summary Activate user
// @Description Activate a pending user
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.ActivateUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/{id}/activate [post]
func (h *UserHandler) ActivateUser(c *gin.Context) {
	// Extrair ID da URL
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.errorHandler.HandleValidationError(c, "id", "Invalid user ID format")
		return
	}

	// Criar input do use case
	input := user.ActivateUserInput{
		UserID: userID,
	}

	// Executar use case
	output, err := h.useCaseAggregate.ActivateUser(c.Request.Context(), input)
	if err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.ActivateUserResponse{
		Success: true,
		Message: "User activated successfully",
		Data:    dto.ToUserResponseFromOutput(output.User),
	}

	// Log do evento
	activatedBy, _ := middleware.GetUserIDFromContext(c)
	h.logger.Info("User activated successfully",
		zap.String("user_id", userID.String()),
		zap.String("activated_by", activatedBy),
	)

	c.JSON(http.StatusOK, response)
}

// DeactivateUser desativa um usuário
// @Summary Deactivate user
// @Description Deactivate an active user
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.DeactivateUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/{id}/deactivate [post]
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	// Extrair ID da URL
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.errorHandler.HandleValidationError(c, "id", "Invalid user ID format")
		return
	}

	// Criar input do use case
	input := user.DeactivateUserInput{
		UserID: userID,
	}

	// Executar use case
	output, err := h.useCaseAggregate.DeactivateUser(c.Request.Context(), input)
	if err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.DeactivateUserResponse{
		Success: true,
		Message: "User deactivated successfully",
		Data:    dto.ToUserResponseFromOutput(output.User),
	}

	// Log do evento
	deactivatedBy, _ := middleware.GetUserIDFromContext(c)
	h.logger.Info("User deactivated successfully",
		zap.String("user_id", userID.String()),
		zap.String("deactivated_by", deactivatedBy),
	)

	c.JSON(http.StatusOK, response)
}

// SuspendUser suspende um usuário
// @Summary Suspend user
// @Description Suspend an active user
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.SuspendUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/{id}/suspend [post]
func (h *UserHandler) SuspendUser(c *gin.Context) {
	// Extrair ID da URL
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.errorHandler.HandleValidationError(c, "id", "Invalid user ID format")
		return
	}

	// Criar input do use case
	input := user.SuspendUserInput{
		UserID: userID,
	}

	// Executar use case
	output, err := h.useCaseAggregate.SuspendUser(c.Request.Context(), input)
	if err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.SuspendUserResponse{
		Success: true,
		Message: "User suspended successfully",
		Data:    dto.ToUserResponseFromOutput(output.User),
	}

	// Log do evento
	suspendedBy, _ := middleware.GetUserIDFromContext(c)
	h.logger.Info("User suspended successfully",
		zap.String("user_id", userID.String()),
		zap.String("suspended_by", suspendedBy),
	)

	c.JSON(http.StatusOK, response)
}

// ==========================================
// ROLE MANAGEMENT ENDPOINTS
// ==========================================

// ChangeRole altera o role de um usuário
// @Summary Change user role
// @Description Change user role (requires appropriate permissions)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param role body dto.ChangeRoleRequest true "New role"
// @Success 200 {object} dto.ChangeRoleResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/{id}/role [put]
func (h *UserHandler) ChangeRole(c *gin.Context) {
	// Extrair ID da URL
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.errorHandler.HandleValidationError(c, "id", "Invalid user ID format")
		return
	}

	var req dto.ChangeRoleRequest

	// Bind e validar request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Validar com validador customizado
	if err := h.validator.ValidateStruct(req); err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Obter ID do usuário que está fazendo a requisição
	requesterID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		h.errorHandler.HandleCustomError(c, "UNAUTHORIZED", "User not authenticated", http.StatusUnauthorized)
		return
	}

	requesterUUID, err := uuid.Parse(requesterID)
	if err != nil {
		h.errorHandler.HandleCustomError(c, "INVALID_REQUESTER", "Invalid requester ID", http.StatusBadRequest)
		return
	}

	// Converter para input do use case
	input := req.ToChangeRoleInput(userID, requesterUUID)

	// Executar use case
	output, err := h.useCaseAggregate.ChangeRole(c.Request.Context(), input)
	if err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.ChangeRoleResponse{
		Success: true,
		Message: "Role changed successfully",
		Data:    dto.ToUserResponseFromOutput(output.User),
	}

	// Log do evento
	h.logger.Info("User role changed successfully",
		zap.String("user_id", userID.String()),
		zap.String("new_role", req.NewRole),
		zap.String("changed_by", requesterID),
	)

	c.JSON(http.StatusOK, response)
}

// ==========================================
// UTILITY ENDPOINTS
// ==========================================

// GetUserStats obtém estatísticas dos usuários
// @Summary Get user statistics
// @Description Get statistics about users in the system
// @Tags users
// @Produce json
// @Success 200 {object} dto.UserStatsOutput
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/stats [get]
func (h *UserHandler) GetUserStats(c *gin.Context) {
	// Executar use case
	output, err := h.useCaseAggregate.GetUserStats(c.Request.Context())
	if err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.SuccessResponse{
		Success: true,
		Message: "User statistics retrieved successfully",
		Data:    output,
	}

	c.JSON(http.StatusOK, response)
}

// CheckUserExists verifica se um usuário existe
// @Summary Check if user exists
// @Description Check if a user exists by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/{id}/exists [get]
func (h *UserHandler) CheckUserExists(c *gin.Context) {
	// Extrair ID da URL
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.errorHandler.HandleValidationError(c, "id", "Invalid user ID format")
		return
	}

	// Executar use case
	exists, err := h.useCaseAggregate.CheckUserExists(c.Request.Context(), userID.String())
	if err != nil {
		h.errorHandler.HandleError(c, err)
		return
	}

	// Converter para response
	response := dto.SuccessResponse{
		Success: true,
		Message: "User existence checked successfully",
		Data: map[string]interface{}{
			"exists": exists,
		},
	}

	c.JSON(http.StatusOK, response)
}
