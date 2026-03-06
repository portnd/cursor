package middlewares

import (
	"github.com/gin-gonic/gin"
)

func UserPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, _ := c.Get("userID")
		userID := int(user_id.(float64))
		CheckPermission(c, userID, []string{"setting_user_acess"})
		c.Next()
	}
}
