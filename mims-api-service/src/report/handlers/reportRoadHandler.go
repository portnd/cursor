package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/responses"
)

// @summary รายงานสรุปข้อมูลสายทาง
// @description Get a summay report of summary of road surface detail report
// @tags Report
// @id type4
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param road_group_ids query string true "Insert the road group ID" default(<1,2,3,4>)
// @Param type query string true "Insert the report type" default(<html, pdf, excel>)
// @response 200 {object} []models.DataMaintenance "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/type4 [get]
func (h *Handler) GetReportRoad(c *gin.Context) {
	roadGroupIDStr := c.Query("road_group_ids")
	typ := c.Query("type")
	var roadGroupIDs []int
	if roadGroupIDStr != "" {
		for _, str := range strings.Split(roadGroupIDStr, ",") {
			i, err := strconv.Atoi(str)
			if err != nil {
				fmt.Printf("Failed to convert '%s' to int: %s\n", str, err)
				continue
			}
			roadGroupIDs = append(roadGroupIDs, i)
		}

	}

	resp, err := h.UseCase.GetReportRoad(roadGroupIDs, typ)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}
