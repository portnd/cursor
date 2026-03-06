package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
)

// @summary เช็คว่าดึงข้อมูลเสร็จหรือยัง
// @description Get DataMart Check
// @tags Dashboard
// @accept json
// @produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @success 200 {object} responses.DataMartCheck
// @failure 400 {object} responses.BadRequestErrorResponse
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/surface/data_mart_check [get]
func (h *Handler) GetDataMartCheck(c *gin.Context) {
	result, err := h.useCase.GetDataMartCheck()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))
}

// @summary ดึงข้อมูล
// @description Get DataMart Check
// @tags Dashboard
// @accept json
// @produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @success 200 {object} []models.DataMart
// @failure 400 {object} responses.BadRequestErrorResponse
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/surface/data_mart [get]
func (h *Handler) GetDataMart(c *gin.Context) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	result, err := h.useCase.GetDataMart(userID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))
}

func (h *Handler) GetDataMartSys(c *gin.Context) {
	result, err := h.useCase.GetDataMart(0)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))
}

// @summary ข้อมูลกราฟ + ตาราง
// @description Get surface dashboard data for specified road IDs
// @tags Dashboard
// @accept json
// @produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param road_id query string false "สายทาง 1,2"
// @Param depot_code query string false "หน่วยงานที่รับผิดชอบ 26103,26104"
// @Param km_start query string false "กม. เริ่มต้น 0"
// @Param km_end query string false " กม. สิ้นสุด 1000"
// @Param year query string false "ปี 2020"
// @success 200 {object} responses.SurfaceRespond
// @failure 400 {object} responses.BadRequestErrorResponse
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/surface [get]
func (h *Handler) GetSurfaceDashboard(c *gin.Context) {
	var strDepotCode []string
	var filter requests.Asset
	if err := c.ShouldBind(&filter); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	kmStart1 := 0.0
	kmStart := c.Query("km_start")
	if kmStart != "" {
		floatVal, _ := strconv.ParseFloat(kmStart, 64)
		kmStart1 = floatVal
	}
	kmEnd1 := 0.0
	kmEnd := c.Query("km_end")
	if kmStart != "" {
		floatVal, _ := strconv.ParseFloat(kmEnd, 64)
		kmEnd1 = floatVal
	}
	filter.KmStart = kmStart1
	filter.KmEnd = kmEnd1
	reqDepotCode := c.Query("depot_code")
	if reqDepotCode != "" {
		strDepotCode = strings.Split(reqDepotCode, ",")
	}

	var roadIds []int
	roadId := c.Query("road_id")
	if roadId != "" {
		for _, item := range strings.Split(roadId, ",") {
			roadIdInt, _ := strconv.Atoi(item)
			roadIds = append(roadIds, roadIdInt)
		}
	}
	helpers.PrintlnJson(filter)
	result, err := h.useCase.GetSurfaceDashboard(roadIds, strDepotCode, filter)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))
}

// @summary ข้อมูลแผนท่ี
// @description Get SurfaceDashboard Map
// @tags Dashboard
// @accept json
// @produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param road_id query string false "สายทาง 1,2"
// @Param depot_code query string false "หน่วยงานที่รับผิดชอบ 26103,26104"
// @Param km_start query string false "กม. เริ่มต้น 0"
// @Param km_end query string false " กม. สิ้นสุด 1000"
// @Param year query string false "ปี 2020"
// @param display query string false "1 = อายุผิวทาง 2 = ชนิดผิวทาง"
// @success 200 {object}  []models.GeomList
// @failure 400 {object} responses.BadRequestErrorResponse
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/surface_map [get]
func (h *Handler) GetSurfaceDashboardMap(c *gin.Context) {
	var strDepotCode []string
	var filter requests.Asset
	if err := c.ShouldBind(&filter); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	kmStart1 := 0.0
	kmStart := c.Query("km_start")
	if kmStart != "" {
		floatVal, _ := strconv.ParseFloat(kmStart, 64)
		kmStart1 = floatVal
	}
	kmEnd1 := 0.0
	kmEnd := c.Query("km_end")
	if kmStart != "" {
		floatVal, _ := strconv.ParseFloat(kmEnd, 64)
		kmEnd1 = floatVal
	}
	filter.KmStart = kmStart1
	filter.KmEnd = kmEnd1
	var roadIds []int
	roadId := c.Query("road_id")
	if roadId != "" {
		for _, item := range strings.Split(roadId, ",") {
			roadIdInt, _ := strconv.Atoi(item)
			roadIds = append(roadIds, roadIdInt)
		}
	}
	reqDepotCode := c.Query("depot_code")
	if reqDepotCode != "" {
		strDepotCode = strings.Split(reqDepotCode, ",")
	}
	displayStr := c.Query("display")
	validateMsgArr := map[string]string{}
	if displayStr == "" || (displayStr != "1" && displayStr != "2") {
		validateMsgArr["display"] = "กรุณาเลือกประเภทแสดงผล"

		if len(validateMsgArr) > 0 {
			errResponse := responses.ValidateResponse(validateMsgArr)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}
	}
	display, _ := strconv.Atoi(displayStr)
	result, err := h.useCase.GetSurfaceDashboardMap(roadIds, strDepotCode, filter, display)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))
}
