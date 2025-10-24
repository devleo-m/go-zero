package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSConfig representa a configuração do CORS
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	MaxAge           int
	AllowCredentials bool
}

// CORS cria um middleware de CORS
func CORS(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Verificar se a origem é permitida
		if isOriginAllowed(origin, config.AllowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
		} else if len(config.AllowedOrigins) > 0 {
			// Se não for permitida e não for wildcard, usar a primeira origem permitida
			c.Header("Access-Control-Allow-Origin", config.AllowedOrigins[0])
		} else {
			// Se não houver origens específicas, permitir todas
			c.Header("Access-Control-Allow-Origin", "*")
		}

		// Configurar outros headers
		if len(config.AllowedMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", joinStrings(config.AllowedMethods, ", "))
		} else {
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}

		if len(config.AllowedHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", joinStrings(config.AllowedHeaders, ", "))
		} else {
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		}

		if len(config.ExposedHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", joinStrings(config.ExposedHeaders, ", "))
		}

		if config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Responder a requisições OPTIONS
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
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}

	return false
}

// joinStrings une strings com um separador
func joinStrings(strs []string, separator string) string {
	if len(strs) == 0 {
		return ""
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += separator + strs[i]
	}

	return result
}
