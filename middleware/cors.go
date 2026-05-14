package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS returns a Gin middleware that injects CORS headers when enabled.
// allowedOrigins is the list of permitted origins from configuration.
func CORS(enabled bool, allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enabled {
			c.Next()
			return
		}

		origin := c.Request.Header.Get("Origin")
		allowed := matchOrigin(origin, allowedOrigins)

		c.Header("Access-Control-Allow-Origin", allowed)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// matchOrigin returns the matching allowed origin, or "*" if the list contains "*".
func matchOrigin(origin string, allowedOrigins []string) string {
	for _, o := range allowedOrigins {
		if o == "*" {
			return "*"
		}
		if strings.EqualFold(o, origin) {
			return origin
		}
	}
	if len(allowedOrigins) > 0 {
		return allowedOrigins[0]
	}
	return "*"
}
