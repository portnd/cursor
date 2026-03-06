package middlewares

import (
	"github.com/gin-gonic/gin"
)

func RolePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, _ := c.Get("userID")
		userID := int(user_id.(float64))
		CheckPermission(c, userID, []string{"setting_permission_acess"})
		c.Next()
	}
}
