package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	_ "gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadSurface/domains"
)

type RoadSurfaceHandler struct {
	roadSurfaceHandler domains.RoadSurfaceUsecases
}

func NewRoadSurfaceHandler(usecase domains.RoadSurfaceUsecases) *RoadSurfaceHandler {
	return &RoadSurfaceHandler{roadSurfaceHandler: usecase}
}

// @summary
// @description
// @tags Road_Surface
// @id get_surface_data
// @param Authorization header string true "JWT token in the format of Bearer {token}"
// @param road_id query string false "Insert your road ID"
// @response 200 {object} responses.Success{data=[]responses.RoadSurfaceResponds} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/road/surface [get]
func (r *RoadSurfaceHandler) GetRoadSurface(c *gin.Context) {
	reqRoadId := c.Query("road_id")
	userId := helpers.GetUserInfo(c).UserId
	userPermission := helpers.GetUserInfo(c).UserPermission
	result, err := r.roadSurfaceHandler.GetRoadSurfaceList(reqRoadId, userPermission, userId)
	if err != nil {
		logs.Error(err)
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	if len(result) == 0 {
		c.JSON(http.StatusOK, responses.SuccessResponse([]string{}, http.StatusOK))
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))

}

// @summary Create a new road surface
// @description This endpoint allows you to create a new road surface with the given payload
// @tags Road_Surface
// @id post_road_surface
// @accept json
// @produce json
// @param Authorization header string true "JWT token in the format of Bearer {token}"
// @param requestBody body requests.RoadSurface true "Request body to create a new road surface"
// @response 200 {object} responses.Success{data=responses.ResPostRoadSurface} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/road/surface [post]
func (r *RoadSurfaceHandler) PostRoadSurface(c *gin.Context) {
	var reqBody requests.RoadSurface
	userIdValue, exists := c.Get("userID")
	if !exists {
		err := errors.New("userId not found")
		logs.Error(err)
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	userId, ok := userIdValue.(float64)
	if !ok {
		err := errors.New("Error to find UserID")
		logs.Error(err)
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	err := c.ShouldBind(&reqBody)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	// reqBody.CreatedBy = int(userId)
	result, err := r.roadSurfaceHandler.PostRoadSurface(reqBody, int(userId))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	var res responses.ResPostRoadSurface
	res.ID = result
	c.JSON(http.StatusAccepted, responses.SuccessResponse(res, http.StatusAccepted))

}

// @summary
// @description
// @tags Road_Surface
// @id get_surface_data_icon
// @param Authorization header string true "JWT token in the format of Bearer {token}"
// @param road_id path string true "Insert your road Id"
// @response 200 {object} responses.Success{data=[]models.RoadSurfaceIcon} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/road/surface/icon/{road_id} [get]
func (r *RoadSurfaceHandler) GetIconRoadSurface(c *gin.Context) {
	reqRoadId, err := strconv.Atoi(c.Params.ByName("road_id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusNotFound, errResponse)
		return
	}
	result, err := r.roadSurfaceHandler.GetRoadSurfaceIconById(reqRoadId)
	if err != nil {
		logs.Error(err)
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))

}
