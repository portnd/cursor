package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/dashboard/domains"
)

type HandlerMaintenance struct {
	useCase domains.UseCaseMaintenance
}

func NewHandlerMaintenance(usecase domains.UseCaseMaintenance) *HandlerMaintenance {
	return &HandlerMaintenance{
		useCase: usecase,
	}
}

// @summary
// @description
// @tags Dashboard
// @id GetMaintenanceDashboard
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param road_id query string false "สายทาง 1,2"
// @Param depot_code query string false "หน่วยงานที่รับผิดชอบ 26103,26104"
// @Param km_start query string false "กม. เริ่มต้น 0"
// @Param km_end query string false " กม. สิ้นสุด 1000"
// @Param year query string false "ปี 2020"
// @response 200 {object} responses.Success{data=responses.MaintenanceDashboard} "OK"
// @failure 400 {object} responses.BadRequestErrorResponse
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/maintenance [get]
func (h *HandlerMaintenance) GetMaintenanceDashboard(c *gin.Context) {
	var roadIds []int
	var strDepotCode []string
	var filter requests.MaintenanceDashboard
	if err := c.ShouldBind(&filter); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	reqDepotCode := c.Query("depot_code")
	if reqDepotCode != "" {
		strDepotCode = strings.Split(reqDepotCode, ",")
	}

	reqRoadId := c.Query("road_id")
	strRoadId := strings.Split(reqRoadId, ",")
	if reqRoadId != "" {
		for _, id := range strRoadId {
			newId, err := helpers.ConvertStringToInt(id)
			if err != nil {
				errResponse := responses.FailRespone(err)
				c.JSON(http.StatusBadRequest, errResponse)
				return
			}
			roadIds = append(roadIds, newId)
		}
	}

	result, err := h.useCase.GetMaintenanceDashboard(roadIds, strDepotCode, filter)
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))
}

// @summary
// @description
// @tags Dashboard
// @id GetMaintenanceMapDashboard
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param road_id query string false "สายทาง 1,2"
// @Param depot_code query string false "หน่วยงานที่รับผิดชอบ 26103,26104"
// @Param km_start query string false "กม. เริ่มต้น 0"
// @Param km_end query string false " กม. สิ้นสุด 1000"
// @Param year query string false "ปี 2020"
// @response 200 {object} responses.Success{data=[]responses.MaintenanceMapDashboard} "OK"
// @failure 400 {object} responses.BadRequestErrorResponse
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/maintenance_map [get]
func (h *HandlerMaintenance) GetMaintenanceMapDashboard(c *gin.Context) {
	var roadIds []int
	var strDepotCode []string
	var filter requests.MaintenanceDashboard
	if err := c.ShouldBind(&filter); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	reqDepotCode := c.Query("depot_code")
	if reqDepotCode != "" {
		strDepotCode = strings.Split(reqDepotCode, ",")
	}

	reqRoadId := c.Query("road_id")
	strRoadId := strings.Split(reqRoadId, ",")
	if reqRoadId != "" {
		for _, id := range strRoadId {
			newId, err := helpers.ConvertStringToInt(id)
			if err != nil {
				errResponse := responses.FailRespone(err)
				c.JSON(http.StatusBadRequest, errResponse)
				return
			}
			roadIds = append(roadIds, newId)
		}
	}

	result, err := h.useCase.GetMaintenanceMapDashboard(roadIds, strDepotCode, filter)
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))
}

// @summary
// @description
// @tags Dashboard
// @id GetMaintenanceTableDashboard
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param page query string true "Insert your page number"
// @param limit query string true "Insert your limit number"
// @Param road_id query string false "สายทาง 1,2"
// @Param depot_code query string false "หน่วยงานที่รับผิดชอบ 26103,26104"
// @Param km_start query string false "กม. เริ่มต้น 0"
// @Param km_end query string false " กม. สิ้นสุด 1000"
// @Param year query string false "ปี 2020"
// @response 200 {object} responses.Success{data=responses.Pagination{item=[]responses.MaintenanceTable}} "OK"
// @failure 400 {object} responses.BadRequestErrorResponse
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/maintenance_table [get]
func (h *HandlerMaintenance) GetMaintenanceTableDashboard(c *gin.Context) {
	var roadIds []int
	var strDepotCode []string
	var filter requests.MaintenanceDashboard
	if err := c.ShouldBind(&filter); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	reqDepotCode := c.Query("depot_code")
	if reqDepotCode != "" {
		strDepotCode = strings.Split(reqDepotCode, ",")
	}

	reqRoadId := c.Query("road_id")
	strRoadId := strings.Split(reqRoadId, ",")
	if reqRoadId != "" {
		for _, id := range strRoadId {
			newId, err := helpers.ConvertStringToInt(id)
			if err != nil {
				errResponse := responses.FailRespone(err)
				c.JSON(http.StatusBadRequest, errResponse)
				return
			}
			roadIds = append(roadIds, newId)
		}
	}

	result, err := h.useCase.GetMaintenanceTableDashboard(roadIds, strDepotCode, filter)
	if err != nil {
		c.JSON(400, err)
		return
	}

	limitParam := strconv.Itoa(filter.Limit)
	pageParam := strconv.Itoa(filter.Page)

	respTable := result.([]responses.MaintenanceTable)
	totalItems := int64(len(respTable))
	limit, offset, page := helpers.GetlimitOffsetPage(limitParam, pageParam, totalItems)

	if totalItems == 0 {
		respTable = []responses.MaintenanceTable{}
	} else if limit+offset > totalItems {
		respTable = respTable[offset:totalItems]
	} else {
		respTable = respTable[offset : limit+offset]
	}

	pagination := helpers.Pagination(respTable, limit, page, totalItems)

	c.JSON(http.StatusOK, responses.SuccessResponse(pagination, http.StatusOK))
}
