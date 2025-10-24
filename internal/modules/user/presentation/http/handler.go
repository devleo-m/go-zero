package http

import (
	"strconv"

	"github.com/devleo-m/go-zero/internal/modules/user/application"
	"github.com/devleo-m/go-zero/internal/modules/user/domain"
	"github.com/devleo-m/go-zero/internal/shared/response"
	"github.com/devleo-m/go-zero/internal/shared/validation"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler gerencia as rotas HTTP para usuários
type Handler struct {
	createUserUseCase *application.CreateUserUseCase
	getUserUseCase    *application.GetUserUseCase
	listUsersUseCase  *application.ListUsersUseCase
	updateUserUseCase *application.UpdateUserUseCase
	deleteUserUseCase *application.DeleteUserUseCase
}

// NewHandler cria uma nova instância do handler
func NewHandler(
	createUserUseCase *application.CreateUserUseCase,
	getUserUseCase *application.GetUserUseCase,
	listUsersUseCase *application.ListUsersUseCase,
	updateUserUseCase *application.UpdateUserUseCase,
	deleteUserUseCase *application.DeleteUserUseCase,
) *Handler {
	return &Handler{
		createUserUseCase: createUserUseCase,
		getUserUseCase:    getUserUseCase,
		listUsersUseCase:  listUsersUseCase,
		updateUserUseCase: updateUserUseCase,
		deleteUserUseCase: deleteUserUseCase,
	}
}

// CreateUser cria um novo usuário
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", err.Error())
		return
	}

	// Validar dados de entrada
	if err := validation.ValidateName(req.Name); err != nil {
		response.BadRequest(c, "INVALID_NAME", err.Error())
		return
	}

	if err := validation.ValidateEmail(req.Email); err != nil {
		response.BadRequest(c, "INVALID_EMAIL", err.Error())
		return
	}

	if err := validation.ValidatePassword(req.Password); err != nil {
		response.BadRequest(c, "INVALID_PASSWORD", err.Error())
		return
	}

	var phone *string
	if req.Phone != "" {
		if err := validation.ValidatePhone(req.Phone); err != nil {
			response.BadRequest(c, "INVALID_PHONE", err.Error())
			return
		}
		phone = &req.Phone
	}

	input := application.CreateUserInput{
		Name:     validation.SanitizeString(req.Name),
		Email:    validation.SanitizeString(req.Email),
		Password: req.Password,
		Phone:    phone,
	}

	result, err := h.createUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		response.BadRequest(c, "CREATE_USER_FAILED", err.Error())
		return
	}

	response.Created(c, toUserResponse(result.User), result.Message)
}

// GetUser busca um usuário por ID
func (h *Handler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	if err := validation.ValidateUUID(idStr); err != nil {
		response.BadRequest(c, "INVALID_ID", err.Error())
		return
	}

	id, _ := uuid.Parse(idStr)
	input := application.GetUserInput{ID: id}
	result, err := h.getUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(c, "USER_NOT_FOUND", "User not found")
			return
		}
		response.InternalServerError(c, "GET_USER_FAILED", err.Error())
		return
	}

	response.Success(c, toUserResponse(result.User))
}

// ListUsers lista usuários
func (h *Handler) ListUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	// Validar parâmetros de paginação
	if err := validation.ValidatePagination(offset/limit+1, limit); err != nil {
		response.BadRequest(c, "INVALID_PAGINATION", err.Error())
		return
	}

	input := application.ListUsersInput{
		Limit:  limit,
		Offset: offset,
	}

	result, err := h.listUsersUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		response.InternalServerError(c, "LIST_USERS_FAILED", err.Error())
		return
	}

	users := make([]UserResponse, len(result.Users))
	for i, user := range result.Users {
		users[i] = toUserResponse(user)
	}

	// Usar sistema de paginação
	page := offset/limit + 1
	meta := response.NewMeta(page, limit, int64(result.Total))

	response.Paginated(c, map[string]interface{}{
		"users": users,
	}, meta)
}

// UpdateUser atualiza um usuário
func (h *Handler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	if err := validation.ValidateUUID(idStr); err != nil {
		response.BadRequest(c, "INVALID_ID", err.Error())
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", err.Error())
		return
	}

	// Validar dados de entrada
	if err := validation.ValidateName(req.Name); err != nil {
		response.BadRequest(c, "INVALID_NAME", err.Error())
		return
	}

	var phone *string
	if req.Phone != "" {
		if err := validation.ValidatePhone(req.Phone); err != nil {
			response.BadRequest(c, "INVALID_PHONE", err.Error())
			return
		}
		phone = &req.Phone
	}

	id, _ := uuid.Parse(idStr)
	input := application.UpdateUserInput{
		ID:    id,
		Name:  validation.SanitizeString(req.Name),
		Phone: phone,
	}

	result, err := h.updateUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(c, "USER_NOT_FOUND", "User not found")
			return
		}
		response.BadRequest(c, "UPDATE_USER_FAILED", err.Error())
		return
	}

	response.Success(c, toUserResponse(result.User), result.Message)
}

// DeleteUser deleta um usuário
func (h *Handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	if err := validation.ValidateUUID(idStr); err != nil {
		response.BadRequest(c, "INVALID_ID", err.Error())
		return
	}

	id, _ := uuid.Parse(idStr)
	input := application.DeleteUserInput{ID: id}
	result, err := h.deleteUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(c, "USER_NOT_FOUND", "User not found")
			return
		}
		response.InternalServerError(c, "DELETE_USER_FAILED", err.Error())
		return
	}

	response.Success(c, nil, result.Message)
}

// toUserResponse converte domain.User para UserResponse
func toUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
