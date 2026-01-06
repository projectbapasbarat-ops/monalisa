package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequirePermission(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permsRaw, exists := c.Get("permissions")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "no permissions",
			})
			return
		}

		permsSlice, ok := permsRaw.([]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "invalid permissions format",
			})
			return
		}

		for _, p := range permsSlice {
			if perm, ok := p.(string); ok && perm == required {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "forbidden",
		})
	}
}
