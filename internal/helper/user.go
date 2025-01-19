package helper

import (
	"post-backend/internal/custom"

	"github.com/gin-gonic/gin"
)

func GetUserId(c *gin.Context) (int, error) {
	userId, exists := c.Get("userId")
	if !exists {
		return -1, custom.ErrUnauthorized
	}

	userIdInt, ok := userId.(int)
	if !ok {
		return -1, custom.ErrUnauthorized
	}

	return userIdInt, nil
}

func GetUserRole(c *gin.Context) (string, error) {
	userRole, exists := c.Get("userRole")
	if !exists {
		return "", custom.ErrUnauthorized

	}

	userRoleString, ok := userRole.(string)
	if !ok {
		return "", custom.ErrUnauthorized
	}

	return userRoleString, nil
}
