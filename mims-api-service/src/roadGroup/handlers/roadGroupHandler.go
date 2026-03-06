package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadGroup/domains"
)

// init handler
type Handler struct {
	UseCase domains.UseCase
}

// init handler
func NewHandler(usecase domains.UseCase) *Handler {
	return &Handler{
		UseCase: usecase,
	}
}

// ================================== start function  ==================================

// @Summary กลุ่มสายทาง
// @Description
// @tags Road Group
// @id GetRoadGroup
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object}  responses.Success{data=[]models.RoadGroup}  "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/road_group [get]
func (h *Handler) GetRoadGroup(c *gin.Context) {
	data, err := h.UseCase.GetRoadGroup()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
}

// @Summary สายทาง
// @Description
// @tags Road Group
// @id GetRoadGroupByID
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object}  responses.Success{data=[]models.RoadGroupByID}  "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/road_group/{id} [get]
func (h *Handler) GetRoadGroupByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	result, err := h.UseCase.GetRoadGroupByID(id)
	if err != nil {
		// if err == gorm.ErrRecordNotFound {
		// 	c.JSON(http.StatusOK, responses.SuccessResponse(responses.NoData{}, http.StatusOK))
		// 	return
		// }
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusNotFound, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))

}
