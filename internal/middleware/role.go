package middleware

import (
	"post-backend/internal/custom"
	"post-backend/internal/helper"
	"slices"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleString, err := helper.GetUserRole(c)
		if err != nil {
			helper.HandleErrorResponde(c, err)
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
