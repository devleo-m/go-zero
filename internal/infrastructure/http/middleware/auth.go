package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims representa as claims do JWT.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware cria um middleware de autenticação JWT.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "AUTHORIZATION_REQUIRED",
				"message": "Authorization header is required",
			})
			c.Abort()

			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "INVALID_TOKEN_FORMAT",
				"message": "Bearer token is required",
			})
			c.Abort()

			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "INVALID_TOKEN",
				"message": "Invalid or expired token",
			})
			c.Abort()

			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "INVALID_TOKEN",
				"message": "Invalid token",
			})
			c.Abort()

			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "INVALID_TOKEN_CLAIMS",
				"message": "Invalid token claims",
			})
			c.Abort()

			return
		}

		// Verificar se o token não expirou
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "TOKEN_EXPIRED",
				"message": "Token has expired",
			})
			c.Abort()

			return
		}

		// Adicionar informações do usuário ao contexto
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("token_claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware cria um middleware de autenticação opcional.
func OptionalAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.Next()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.Next()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.Next()
			return
		}

		// Verificar se o token não expirou
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			c.Next()
			return
		}

		// Adicionar informações do usuário ao contexto se o token for válido
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("token_claims", claims)

		c.Next()
	}
}

// RequireRole cria um middleware que requer um role específico.
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "AUTHENTICATION_REQUIRED",
				"message": "Authentication is required",
			})
			c.Abort()

			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "INVALID_ROLE",
				"message": "Invalid user role",
			})
			c.Abort()

			return
		}

		// Verificar se o usuário tem o role necessário
		if !hasRequiredRole(role, requiredRole) {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "INSUFFICIENT_PERMISSIONS",
				"message": "Insufficient permissions",
			})
			c.Abort()

			return
		}

		c.Next()
	}
}

// RequireAnyRole cria um middleware que requer qualquer um dos roles especificados.
func RequireAnyRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "AUTHENTICATION_REQUIRED",
				"message": "Authentication is required",
			})
			c.Abort()

			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "INVALID_ROLE",
				"message": "Invalid user role",
			})
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
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "INSUFFICIENT_PERMISSIONS",
				"message": "Insufficient permissions",
			})
			c.Abort()

			return
		}

		c.Next()
	}
}

// hasRequiredRole verifica se o usuário tem o role necessário.
func hasRequiredRole(userRole, requiredRole string) bool {
	// Hierarquia de roles (do menor para o maior)
	roleHierarchy := map[string]int{
		"user":        1,
		"moderator":   2,
		"admin":       3,
		"super_admin": 4,
	}

	userLevel, userExists := roleHierarchy[userRole]
	requiredLevel, requiredExists := roleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return false
	}

	// O usuário deve ter pelo menos o nível necessário
	return userLevel >= requiredLevel
}

// GetUserID extrai o ID do usuário do contexto.
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	id, ok := userID.(string)

	return id, ok
}

// GetUserRole extrai o role do usuário do contexto.
func GetUserRole(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", false
	}

	role, ok := userRole.(string)

	return role, ok
}

// GetUserEmail extrai o email do usuário do contexto.
func GetUserEmail(c *gin.Context) (string, bool) {
	userEmail, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	email, ok := userEmail.(string)

	return email, ok
}
