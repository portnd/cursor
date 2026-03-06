package middlewares

import (
	"github.com/gin-gonic/gin"
)

func AssetGroupPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, _ := c.Get("userID")
		userID := int(user_id.(float64))
		CheckPermission(c, userID, []string{"setting_asset_place_group_access"})
		c.Next()
	}
}

func SurfacePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, _ := c.Get("userID")
		userID := int(user_id.(float64))
		CheckPermission(c, userID, []string{"setting_surface_type_access"})
		c.Next()
	}
}

func ConditionListPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userPermission := helpers.GetUserInfo(c).UserPermission
		// if !helpers.HasPermission([]string{"setting_survey_rule_access"}, userPermission) {
		// 	errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		// 	c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
		// }
		// c.Next()
	}
}

func SignPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, _ := c.Get("userID")
		userID := int(user_id.(float64))
		CheckPermission(c, userID, []string{"setting_traffic_sign_access"})
		c.Next()
	}
}

func AssetTablePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, _ := c.Get("userID")
		userID := int(user_id.(float64))
		CheckPermission(c, userID, []string{"setting_in_asset_access"})
		c.Next()
	}
}

func BudgetPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, _ := c.Get("userID")
		userID := int(user_id.(float64))
		CheckPermission(c, userID, []string{"setting_ref_budget_access"})
		c.Next()
	}
}

func ModelPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, _ := c.Get("userID")
		userID := int(user_id.(float64))
		CheckPermission(c, userID, []string{"setting_model_access"})
		c.Next()
	}
}

func HrisPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, _ := c.Get("userID")
		userID := int(user_id.(float64))
		CheckPermission(c, userID, []string{"setting_hris_access"})
		c.Next()
	}
}

func HsmsPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, _ := c.Get("userID")
		userID := int(user_id.(float64))
		CheckPermission(c, userID, []string{"setting_hsms_access"})
		c.Next()
	}
}

func RoadConditionCriteriaPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, _ := c.Get("userID")
		userID := int(user_id.(float64))
		CheckPermission(c, userID, []string{"setting_condition_criteria_access"})
		c.Next()
	}
}

func RoadG7CriteriaPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, _ := c.Get("userID")
		userID := int(user_id.(float64))
		CheckPermission(c, userID, []string{"setting_retro_criteria_access"})
		c.Next()
	}
}
