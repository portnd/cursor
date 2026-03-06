package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
)

// @summary รายงานประวัติการซ่อมบำรุง
// @description Get a summay report of summary of road surface detail report
// @tags Report
// @id Report13
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param year_start query string true "Insert the year" default(<Add a year>)
// @Param year_end query string true "Insert the year" default(<Add a year>)
// @Param road_section_id query string true "Insert the road section ID" default(<Add a road section ID>)
// @Param type query string true "Insert the report type" default(<html, pdf, excel>)
// @Success 200 {object} string "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/type13 [get]
func (h *Handler) Report13(c *gin.Context) {
	var req requests.Report13
	if err := c.ShouldBindQuery(&req); err != nil {
		newError := responses.NewAppErr(500, err.Error())
		errResponse := responses.FailRespone(newError)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	req.Group = 3
	resp, err := h.UseCase.Report13(req)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}
