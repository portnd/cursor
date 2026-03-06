package handlers

import (
	"net/http"
	"strconv"

	responses "gitlab.com/mims-api-service/responses"

	"github.com/gin-gonic/gin"
)

// @summary
// @description
// @tags RoadConditions
// @id GetRoadConditionTemplate
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @response 200 {object} responses.Success{data}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/condition_template [get]
func (t *RoadConditionHandler) GetRoadConditionTemplate(c *gin.Context) {

	roadId, _ := strconv.Atoi(c.Params.ByName("id"))
	err := c.ShouldBind(&roadId)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))
	data, err := t.rcUseCase.GetMenu(uid)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	var accessKeys []string
	for _, item := range data {
		accessKeys = append(accessKeys, item.AccessKey)
	}

	// permissions := []string{
	// 	"road_condition_manage_data",
	// }

	// if helpers.HasPermission(permissions, accessKeys) {

	resp, err := t.rcUseCase.GetRoadConditionTemplate(uid, roadId)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(int(http.StatusOK), responses.SuccessResponse(resp, 200))
	// } else {
	// 	errResponse := responses.FailRespone(errors.New("access denied"))
	// 	c.JSON(http.StatusBadRequest, errResponse)
	// 	return
	// }
}
