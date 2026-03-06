package middlewares

import (
	"github.com/gin-gonic/gin"
)

func AnalyzeAccessPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"maintenance_analysis_access"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func AnalyzeManagePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"maintenance_analysis_manage_data", "maintenance_analysis_access"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}
