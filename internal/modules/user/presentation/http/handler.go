package http

import (
	"net/http"
	"strconv"

	"github.com/devleo-m/go-zero/internal/modules/user/application"
	"github.com/devleo-m/go-zero/internal/modules/user/domain"
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
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	input := application.CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
	}

	result, err := h.createUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "CREATE_USER_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: result.Message,
		Data:    toUserResponse(result.User),
	})
}

// GetUser busca um usuário por ID
func (h *Handler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_ID",
			Message: "Invalid user ID",
		})
		return
	}

	input := application.GetUserInput{ID: id}
	result, err := h.getUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "USER_NOT_FOUND",
				Message: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "GET_USER_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: toUserResponse(result.User),
	})
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

	input := application.ListUsersInput{
		Limit:  limit,
		Offset: offset,
	}

	result, err := h.listUsersUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "LIST_USERS_FAILED",
			Message: err.Error(),
		})
		return
	}

	users := make([]UserResponse, len(result.Users))
	for i, user := range result.Users {
		users[i] = toUserResponse(user)
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: map[string]interface{}{
			"users": users,
			"total": result.Total,
		},
	})
}

// UpdateUser atualiza um usuário
func (h *Handler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_ID",
			Message: "Invalid user ID",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	input := application.UpdateUserInput{
		ID:    id,
		Name:  req.Name,
		Phone: req.Phone,
	}

	result, err := h.updateUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "USER_NOT_FOUND",
				Message: "User not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "UPDATE_USER_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: result.Message,
		Data:    toUserResponse(result.User),
	})
}

// DeleteUser deleta um usuário
func (h *Handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_ID",
			Message: "Invalid user ID",
		})
		return
	}

	input := application.DeleteUserInput{ID: id}
	result, err := h.deleteUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "USER_NOT_FOUND",
				Message: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "DELETE_USER_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: result.Message,
	})
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
