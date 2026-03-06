package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/responses"
)

// @summary รายงานสรุปรายละเอียดชนิดผิวทาง
// @description Get a summay report of summary of road surface detail report
// @tags Report
// @id Report5
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param year query string true "Insert the year" default(<Add a year>)
// @Param road_section_id query string true "Insert the road section ID" default(<Add a road section ID>)
// @Param type query string true "Insert the report type" default(<html, pdf, excel>)
// @Success 200 {object} string "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/type5 [get]
func (h *Handler) Report5(c *gin.Context) {
	year := c.Query("year")
	roadSectionID := c.Query("road_section_id")
	typ := c.Query("type")
	resp, err := h.UseCase.Report5(year, roadSectionID, typ)
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

// @summary รายงานสรุปข้อมูลสภาพทาง
// @description Get a summay report of summary of road surface detail report
// @tags Report
// @id Report7
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param year query string true "Insert the year" default(<Add a year>)
// @Param road_section_id query string true "Insert the road section ID" default(<Add a road section ID>)
// @Param type query string true "Insert the report type" default(<html, pdf, excel>)
// @Param factor query string true "Insert the filter_condition" default(<IRI,MPD,RUT,IFI>)
// @Param filter_criteria_id query string true "Insert the filter_criteria id"
// @Success 200 {object} string "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/type7 [get]
func (h *Handler) Report7(c *gin.Context) {
	year := c.Query("year")
	roadSectionID := c.Query("road_section_id")
	typ := c.Query("type")
	factor := c.Query("factor")
	measureID := c.Query("filter_criteria_id")
	resp, err := h.UseCase.Report7(factor, year, roadSectionID, measureID, typ)
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

// @Summary รายงานข้อมูลสภาพทาง
// @Description Get a summary report of road condition
// @Tags Report
// @ID GetReportCondition
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param year query string true "Insert a year" default(<Add a year here>)
// @Param road_section_id query string true "Insert the road ID" default(<Add a road ID here>)
// @Param type query string true "Insert the report type" default(<html, pdf, excel>)
// @Param distance query string true "Insert the report distance" default(<25, 100, 1000>)
// @Produce json
// @Success 200 {object} string "OK"
// @Failure 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @Failure 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/report/type6 [get]
func (h *Handler) Report6(c *gin.Context) {
	year := c.Query("year")
	roadID := c.Query("road_section_id")
	dis := c.Query("distance")
	typ := c.Query("type")
	resp, err := h.UseCase.Report6(year, roadID, typ, dis)
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
