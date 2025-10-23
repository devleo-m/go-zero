package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/devleo-m/go-zero/internal/interface/http/dto"
)

// JWTService interface para serviços de JWT
type JWTService interface {
	ValidateToken(token string) (map[string]interface{}, error)
	ExtractUserID(claims map[string]interface{}) (string, error)
	ExtractUserRole(claims map[string]interface{}) (string, error)
	ExtractUserEmail(claims map[string]interface{}) (string, error)
}

// AuthMiddleware middleware de autenticação JWT
func AuthMiddleware(jwtService JWTService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Extrair token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing authorization header",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("ip", c.ClientIP()),
			)

			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED",
				"Authorization header is required",
			))
			c.Abort()
			return
		}

		// Verificar se o header tem o formato correto "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.Warn("Invalid authorization header format",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("ip", c.ClientIP()),
			)

			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED",
				"Invalid authorization header format. Expected: Bearer <token>",
			))
			c.Abort()
			return
		}

		token := parts[1]

		// Validar token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			logger.Warn("Invalid JWT token",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("ip", c.ClientIP()),
				zap.Error(err),
			)

			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED",
				"Invalid or expired token",
			))
			c.Abort()
			return
		}

		// Extrair informações do usuário dos claims
		userID, err := jwtService.ExtractUserID(claims)
		if err != nil {
			logger.Error("Failed to extract user ID from token",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("ip", c.ClientIP()),
				zap.Error(err),
			)

			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED",
				"Invalid token claims",
			))
			c.Abort()
			return
		}

		userRole, err := jwtService.ExtractUserRole(claims)
		if err != nil {
			logger.Error("Failed to extract user role from token",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("ip", c.ClientIP()),
				zap.Error(err),
			)

			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED",
				"Invalid token claims",
			))
			c.Abort()
			return
		}

		userEmail, err := jwtService.ExtractUserEmail(claims)
		if err != nil {
			logger.Error("Failed to extract user email from token",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("ip", c.ClientIP()),
				zap.Error(err),
			)

			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED",
				"Invalid token claims",
			))
			c.Abort()
			return
		}

		// Adicionar informações do usuário ao contexto
		c.Set("user_id", userID)
		c.Set("user_role", userRole)
		c.Set("user_email", userEmail)
		c.Set("jwt_claims", claims)

		// Log de sucesso
		logger.Info("User authenticated successfully",
			zap.String("user_id", userID),
			zap.String("user_role", userRole),
			zap.String("user_email", userEmail),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("ip", c.ClientIP()),
			zap.Duration("duration", time.Since(start)),
		)

		c.Next()
	}
}

// OptionalAuthMiddleware middleware de autenticação opcional
// Se o token estiver presente, valida e adiciona ao contexto
// Se não estiver presente, continua sem erro
func OptionalAuthMiddleware(jwtService JWTService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Verificar formato do header
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		token := parts[1]

		// Validar token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			// Se o token for inválido, continua sem adicionar ao contexto
			c.Next()
			return
		}

		// Extrair informações do usuário
		userID, err := jwtService.ExtractUserID(claims)
		if err != nil {
			c.Next()
			return
		}

		userRole, err := jwtService.ExtractUserRole(claims)
		if err != nil {
			c.Next()
			return
		}

		userEmail, err := jwtService.ExtractUserEmail(claims)
		if err != nil {
			c.Next()
			return
		}

		// Adicionar informações do usuário ao contexto
		c.Set("user_id", userID)
		c.Set("user_role", userRole)
		c.Set("user_email", userEmail)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// RequireRole middleware que verifica se o usuário tem um role específico
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED",
				"User role not found in context",
			))
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
				"INTERNAL_ERROR",
				"Invalid user role format",
			))
			c.Abort()
			return
		}

		// Verificar se o usuário tem o role necessário
		if !hasRequiredRole(role, requiredRole) {
			c.JSON(http.StatusForbidden, dto.NewErrorResponse(
				"FORBIDDEN",
				"Insufficient permissions. Required role: "+requiredRole,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole middleware que verifica se o usuário tem qualquer um dos roles especificados
func RequireAnyRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED",
				"User role not found in context",
			))
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
				"INTERNAL_ERROR",
				"Invalid user role format",
			))
			c.Abort()
			return
		}

		// Verificar se o usuário tem algum dos roles necessários
		hasRole := false
		for _, requiredRole := range requiredRoles {
			if hasRequiredRole(role, requiredRole) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, dto.NewErrorResponse(
				"FORBIDDEN",
				"Insufficient permissions. Required roles: "+strings.Join(requiredRoles, ", "),
			))
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasRequiredRole verifica se o usuário tem o role necessário
// Implementa uma hierarquia de roles: admin > manager > user > guest
func hasRequiredRole(userRole, requiredRole string) bool {
	// Admin tem acesso a tudo
	if userRole == "admin" {
		return true
	}

	// Se o role necessário for admin, apenas admin pode acessar
	if requiredRole == "admin" {
		return userRole == "admin"
	}

	// Manager pode acessar manager, user e guest
	if userRole == "manager" {
		return requiredRole == "manager" || requiredRole == "user" || requiredRole == "guest"
	}

	// User pode acessar user e guest
	if userRole == "user" {
		return requiredRole == "user" || requiredRole == "guest"
	}

	// Guest só pode acessar guest
	if userRole == "guest" {
		return requiredRole == "guest"
	}

	return false
}

// GetUserIDFromContext extrai o user ID do contexto
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	id, ok := userID.(string)
	return id, ok
}

// GetUserRoleFromContext extrai o user role do contexto
func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", false
	}

	role, ok := userRole.(string)
	return role, ok
}

// GetUserEmailFromContext extrai o user email do contexto
func GetUserEmailFromContext(c *gin.Context) (string, bool) {
	userEmail, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	email, ok := userEmail.(string)
	return email, ok
}

// GetJWTClaimsFromContext extrai os claims JWT do contexto
func GetJWTClaimsFromContext(c *gin.Context) (map[string]interface{}, bool) {
	claims, exists := c.Get("jwt_claims")
	if !exists {
		return nil, false
	}

	claimsMap, ok := claims.(map[string]interface{})
	return claimsMap, ok
}
