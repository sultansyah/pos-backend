package middleware

import (
	"post-backend/internal/custom"
	"post-backend/internal/helper"
	"slices"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			helper.HandleErrorResponde(c, custom.ErrUnauthorized)
			c.Abort()
			return
		}

		userRoleString, ok := userRole.(string)
		if !ok {
			helper.HandleErrorResponde(c, custom.ErrInternal)
			c.Abort()
			return
		}

		isSame := slices.Contains(roles, userRoleString)
		if !isSame {
			helper.HandleErrorResponde(c, custom.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}
