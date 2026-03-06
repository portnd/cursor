package middlewares

import (
	"github.com/gin-gonic/gin"
)

func ReportPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"report_access", "report_download_asset", "report_download_asset_map", "report_download_asset_summary", "report_download_improve_asset", "report_download_road_condition", "report_download_road_condition_summary", "report_download_road_damage_summary", "report_download_road_surface", "report_download_project_follow_up", "report_download_project_maintenance", "report_download_volume_aadt", "report_download_volume_accident"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}
