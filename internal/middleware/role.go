package middleware

import (
	"post-backend/internal/custom"
	"post-backend/internal/helper"
	"slices"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			helper.HandleErrorResponde(c, custom.ErrUnauthorized)
			c.Abort()
		}

		userRoleString, ok := userRole.(string)
		if !ok {
			helper.HandleErrorResponde(c, custom.ErrInternal)
			c.Abort()
		}

		exists = slices.Contains(roles, userRoleString)
		if !exists {
			helper.HandleErrorResponde(c, custom.ErrUnauthorized)
			c.Abort()
		}

		c.Next()
	}
}
