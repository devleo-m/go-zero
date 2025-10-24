package http

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configura as rotas para usuários
func SetupRoutes(router *gin.Engine, handler *Handler) {
	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", handler.CreateUser)       // POST /api/v1/users
			users.GET("", handler.ListUsers)         // GET /api/v1/users
			users.GET("/:id", handler.GetUser)       // GET /api/v1/users/:id
			users.PUT("/:id", handler.UpdateUser)    // PUT /api/v1/users/:id
			users.DELETE("/:id", handler.DeleteUser) // DELETE /api/v1/users/:id
		}
	}
}
