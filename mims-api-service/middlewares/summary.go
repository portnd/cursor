package middlewares

import (
	"github.com/gin-gonic/gin"
)

func SummaryPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// 'edit_asset_in_data', 'view_asset_in_data', 'edit_asset_out_data', 'view_asset_out_data'
		// if !helpers.HasPermission([]string{"dashboard_road_condition_access", "dashboard_road_surface_access", "dashboard_road_asset_access"}, permissions) {
		// errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }

	}
}
