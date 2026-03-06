package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadSection/domains"
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

// @Summary กลุ่มตอนควบคุม
// @Description
// @tags Road Section
// @id GetRoadSection
// @Param road_group_id query string false "road_group_id" example(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object}  responses.Success{data=[]responses.RoadSection}  "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/road_section [get]
func (h *Handler) GetRoadSection(c *gin.Context) {
	roadGroupIdQuery, ok := c.Request.URL.Query()["road_group_number"]
	var roadGroupId *int
	if ok {
		RoadGroupIdInt, err := strconv.Atoi(roadGroupIdQuery[0])
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(http.StatusBadRequest, errResponse)
			return
		}
		roadGroupId = &RoadGroupIdInt
	}
	data, err := h.UseCase.GetRoadSection(roadGroupId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
}

// @Summary ตอนควบคุม
// @Description
// @tags Road Section
// @id GetRoadSectionByID
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road section ID"
// @Success 200 {object}  responses.Success{data=responses.RoadSection}  "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/road_section/{id} [get]
func (h *Handler) GetRoadSectionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	result, err := h.UseCase.GetRoadSectionByID(id)
	if err != nil {
		// if err == gorm.ErrRecordNotFound {
		// 	c.JSON(http.StatusOK, responses.SuccessResponse(responses.NoData{}, http.StatusOK))
		// 	return
		// }
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusNotFound, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(*result, http.StatusOK))

}
