package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig configurações do CORS
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig retorna configuração padrão de CORS
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:3001",
			"http://127.0.0.1:8080",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"HEAD",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"X-Request-ID",
			"X-User-ID",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-ID",
			"X-Total-Count",
			"X-Page-Count",
		},
		AllowCredentials: true,
		MaxAge:           86400, // 24 horas
	}
}

// CORS middleware para Cross-Origin Resource Sharing
func CORS(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Verificar se a origem é permitida
		if isOriginAllowed(origin, config.AllowOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
		} else if len(config.AllowOrigins) > 0 && config.AllowOrigins[0] == "*" {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		// Configurar outros headers CORS
		c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
		c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))

		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		// Responder a requisições OPTIONS (preflight)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// isOriginAllowed verifica se a origem é permitida
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}

	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" {
			return true
		}
		if allowedOrigin == origin {
			return true
		}
		// Suporte para wildcards simples (ex: *.example.com)
		if strings.Contains(allowedOrigin, "*") {
			if matchWildcard(origin, allowedOrigin) {
				return true
			}
		}
	}

	return false
}

// matchWildcard verifica se uma origem corresponde a um padrão wildcard
func matchWildcard(origin, pattern string) bool {
	// Implementação simples de wildcard
	// Ex: *.example.com deve corresponder a subdomain.example.com
	if strings.HasPrefix(pattern, "*.") {
		suffix := pattern[2:] // Remove "*. "
		return strings.HasSuffix(origin, suffix)
	}

	return false
}

// CORSForDevelopment configuração CORS para desenvolvimento
func CORSForDevelopment() gin.HandlerFunc {
	config := CORSConfig{
		AllowOrigins: []string{
			"http://localhost:*",
			"http://127.0.0.1:*",
			"http://0.0.0.0:*",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"HEAD",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"*",
		},
		ExposeHeaders: []string{
			"*",
		},
		AllowCredentials: true,
		MaxAge:           0, // Sem cache para desenvolvimento
	}

	return CORS(config)
}

// CORSForProduction configuração CORS para produção
func CORSForProduction(allowedOrigins []string) gin.HandlerFunc {
	config := CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 horas
	}

	return CORS(config)
}
