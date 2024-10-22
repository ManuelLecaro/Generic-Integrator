package middleware

import (
	"generic-integration-platform/internal/infra/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyMiddleware validates the presence and correctness of the x-api-key header.
func APIKeyMiddleware(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-api-key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key is missing"})
			c.Abort()
			return
		}

		if apiKey != config.APIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
