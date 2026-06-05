package middlewares

import (
	"net/http"
	"strings"

	"backend/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "authorization token is required",
			})
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "authorization token must use Bearer format",
			})
			return
		}

		claims, err := utils.ValidateJWT(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "invalid or expired token",
			})
			return
		}

		c.Set("user_id", claims.ID)
		c.Set("name", claims.Name)
		c.Set("email", claims.Email)
		c.Next()
	}
}
