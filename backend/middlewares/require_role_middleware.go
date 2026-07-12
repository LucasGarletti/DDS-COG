package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "authentication required",
			})
			return
		}

		currentRole, ok := roleValue.(string)
		if !ok || currentRole != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "insufficient permissions",
			})
			return
		}

		c.Next()
	}
}
