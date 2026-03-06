package middlewares

import (
	"github.com/gin-gonic/gin"
)

func MaintenancePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"maintenance_view_data", "maintenance_manage_data"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func MaintenanceAccessPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"maintenance_view_data"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func MaintenanceManagePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"maintenance_manage_data"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func MaintenanceHisPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"maintenance_history_view_data", "maintenance_history_manage_data"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func MaintenanceHisAccessPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"maintenance_history_view_data"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func MaintenanceHisManagePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"maintenance_history_manage_data"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}
