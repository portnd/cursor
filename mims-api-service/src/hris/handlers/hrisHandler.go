package handlers

import (
	"net/http"

	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/hris/domains"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Usecase domains.UseCase
}

func NewHandler(usecase domains.UseCase) *Handler {
	return &Handler{
		Usecase: usecase,
	}
}

// @summary
// @description
// @tags hris
// @id get_section_geom
// @response 200 {object} responses.Success{data=[]string{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/hris/section_geom [get]
func (h *Handler) GetSectionGeom(c *gin.Context) {
	err := h.Usecase.GetSectionGeom()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse([]string{}, http.StatusOK))
}

// @summary
// @description
// @tags hris
// @id get_road_latest
// @response 200 {object} responses.Success{data=[]string{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/hris/road_latest [get]
func (h *Handler) GetRoadLatest(c *gin.Context) {
	err := h.Usecase.GetRoadLatest()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse([]string{}, http.StatusOK))
}

// @summary
// @description
// @tags hris
// @id get_match_data
// @response 200 {object} responses.Success{data=[]string{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/hris/match_data [get]
func (h *Handler) MatchData(c *gin.Context) {
	response, err := h.Usecase.MatchData()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(response, http.StatusOK))
}
