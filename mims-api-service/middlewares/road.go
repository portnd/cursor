package middlewares

import (
	"github.com/gin-gonic/gin"
)

func RoadPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"dashboard_road_condition_access", "road_summary_access", "road_summary_manage", "road_in_asset_access", "road_in_asset_manage_data", "road_out_asset_access", "road_out_asset_manage_data", "road_condition_access", "road_condition_manage_data", "road_damage_access", "road_damage_manage_data", "road_aadt_access", "road_aadt_manage_data", "road_accident_access", "road_accident_manage_data"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func RoadConditionAccessPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"road_condition_access"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func RoadConditionManagePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"road_condition_manage_data"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func RoadDamageAccessPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"road_damage_access"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func RoadDamageManagePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"road_damage_manage_data"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func RoadAssetAccessPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// // if !helpers.HasPermission([]string{"road_in_asset_access", "road_out_asset_access"}, permissions) {
		// // 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// // 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// // }
	}
}

func RoadAssetManagePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"road_in_asset_manage_data", "road_out_asset_manage_data"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func RoadAadtAccessPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"road_aadt_access"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func RoadAadtManagePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"road_aadt_manage_data"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func RoadAccidentAccessPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"road_accident_access"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}

func RoadAccidentManagePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission, _ := c.Get("accessControl")
		// permissions := []string{}
		// for _, item := range userPermission.([]interface{}) {
		// 	permissions = append(permissions, item.(string))
		// }
		// if !helpers.HasPermission([]string{"road_accident_manage_data"}, permissions) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
	}
}
