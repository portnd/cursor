package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/responses"
)

// @Summary รายงานข้อมูลความเสียหาย
// @Description Get a  report of road damage
// @Tags Report
// @ID GetReportDamage
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param year query string true "Insert a year" default(<Add a year here>)
// @Param road_section_id query string true "Insert the road ID" default(<Add a road ID here>)
// @Param type query string true "Insert the report type" default(<html, pdf, excel>)
// @Produce json
// @Success 200 {object} string "OK"
// @Failure 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @Failure 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/report/type10 [get]
func (h *Handler) Report10(c *gin.Context) {
	budgetYear := c.Query("year")
	roadID := c.Query("road_section_id")
	typ := c.Query("type")
	resp, err := h.UseCase.Report10(budgetYear, roadID, typ)
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

// @Summary รายงานสรุปข้อมูลความเสียหาย
// @Description Get a summary report of road damage
// @Tags Report
// @ID GetReportSummaryDamage
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param year query string true "Insert a year" default(<Add a year here>)
// @Param road_section_id query string true "Insert the road ID" default(<Add a road ID here>)
// @Param type query string true "Insert the report type" default(<html, pdf, excel>)
// @Produce json
// @Success 200 {object} string "OK"
// @Failure 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @Failure 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/report/type11 [get]
func (h *Handler) Report11(c *gin.Context) {
	budgetYear := c.Query("year")
	roadID := c.Query("road_section_id")
	typ := c.Query("type")
	resp, err := h.UseCase.Report11(budgetYear, roadID, typ)
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
