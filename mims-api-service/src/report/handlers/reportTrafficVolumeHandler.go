package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/responses"
)

// @summary รายงานสรุปข้อมูลปริมาณจราจร
// @description
// @tags Report
// @id type14
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param year query string true "Insert the year" default(2024)
// @Param road_section_ids query string true "Insert the Road Section ID" default(11566) Ex. 1,2,3,4
// @Param type query string true "Insert the report type" default(pdf) Ex. tml,pdf,excel
// @response 200 {object} responses.Success{data} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/type14 [get]
func (h *Handler) GetReportTrafficVolume(c *gin.Context) {
	roadSectionIDStr := c.Query("road_section_ids")
	typ := c.Query("type")
	year := c.Query("year")

	var roadSectionIDs []string
	if roadSectionIDStr != "" {
		roadSectionIDs = strings.Split(roadSectionIDStr, ",")

	}

	yearInt, _ := strconv.Atoi(year)

	resp, err := h.UseCase.GetReportTrafficVolume(roadSectionIDs, typ, yearInt)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}
