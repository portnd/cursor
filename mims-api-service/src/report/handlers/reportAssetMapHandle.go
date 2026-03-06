package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/responses"
)

// @summary รายงานแผนที่สินทรัพย์
// @description Get a summay report of map of asset report
// @tags Report
// @id GetReportType2
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param road_section_id query string true "Insert the road section ID" default(<Add a road section ID>)
// @Param asset_id query string true "Insert an asset ID" default(<Add an asset ID here>)
// @Param type query string true "Insert the report type" default(<html, pdf, excel>)
// @response 200 {object} []models.DataMaintenance "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/type2 [get]
func (h *Handler) GetReport2(c *gin.Context) {
	roadID := c.Query("road_section_id")
	assetID := c.Query("asset_id")
	typ := c.Query("type")

	resp, err := h.UseCase.GetReport2(roadID, assetID, typ)
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
