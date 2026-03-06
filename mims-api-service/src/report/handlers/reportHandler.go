package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/report/domains"
)

type Handler struct {
	UseCase domains.UseCase
}

func NewHandler(usecase domains.UseCase) *Handler {
	return &Handler{
		UseCase: usecase,
	}
}

//////////////// NEW MIMS ////////////////

// @summary รายงานสินทรัพย์
// @description Get Asset Filter Type1
// @tags Report Filter
// @id GetAssetFilterType1
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.FilterAssetType1{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/asset/filter/type1 [get]
func (h *Handler) GetAssetFilterType1(c *gin.Context) {
	userID := helpers.GetUserID(c)
	resp, err := h.UseCase.FilterAssetType1(userID)
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

// @summary รายงานแผนที่สินทรัพย์
// @description Get Asset Filter Type2
// @tags Report Filter
// @id GetAssetFilterType2
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.FilterAssetType2{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/asset/filter/type2 [get]
func (h *Handler) GetAssetFilterType2(c *gin.Context) {
	userID := helpers.GetUserID(c)
	resp, err := h.UseCase.FilterAssetType2(userID)
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

// @summary รายงานสรุปสินทรัพย์
// @description Get Asset Filter Type3
// @tags Report Filter
// @id GetAssetFilterType3
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.FilterAssetType3{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/asset/filter/type3 [get]
func (h *Handler) GetAssetFilterType3(c *gin.Context) {
	userID := helpers.GetUserID(c)
	resp, err := h.UseCase.FilterAssetType3(userID)
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

// @summary รายงานสรุปข้อมูลสายทาง
// @description Get Road Filter Type1
// @tags Report Filter
// @id GetRoadFilterType1
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.FilterRoadType1{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/road/filter/type1 [get]
func (h *Handler) GetRoadFilterType1(c *gin.Context) {
	userID := helpers.GetUserID(c)
	resp, err := h.UseCase.FilterRoadType1(userID)
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

// @summary รายงานสรุปรายละเอียดชนิดผิวทาง
// @description Get Road Filter Type2
// @tags Report Filter
// @id GetRoadFilterType2
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.FilterRoadType2{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/road/filter/type2 [get]
func (h *Handler) GetRoadFilterType2(c *gin.Context) {
	userID := helpers.GetUserID(c)
	resp, err := h.UseCase.FilterRoadType2(userID)
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

// @summary รายงานข้อมูลสภาพทาง
// @description Get Road Filter Type3
// @tags Report Filter
// @id GetRoadFilterType3
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.FilterRoadType3{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/road/filter/type3 [get]
func (h *Handler) GetRoadFilterType3(c *gin.Context) {
	userID := helpers.GetUserID(c)
	resp, err := h.UseCase.FilterRoadType3(userID)
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
// @description Get Road Filter Type4
// @tags Report Filter
// @id GetRoadFilterType4
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.FilterRoadType4{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/road/filter/type4 [get]
func (h *Handler) GetRoadFilterType4(c *gin.Context) {
	userID := helpers.GetUserID(c)
	resp, err := h.UseCase.FilterRoadType4(userID)
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

// @summary รายงานค่าการสะท้อนแสงของเส้นจราจร
// @description Get Road Filter Type5
// @tags Report Filter
// @id GetRoadFilterType5
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.FilterRoadType5{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/road/filter/type5 [get]
func (h *Handler) GetRoadFilterType5(c *gin.Context) {
	userID := helpers.GetUserID(c)
	resp, err := h.UseCase.FilterRoadType5(userID)
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

// @summary รายงานสรุปข้อมูลค่าสะท้อนแสงของเส้นจราจร
// @description Get Road Filter Type6
// @tags Report Filter
// @id GetRoadFilterType6
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.FilterRoadType6{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/road/filter/type6 [get]
func (h *Handler) GetRoadFilterType6(c *gin.Context) {
	userID := helpers.GetUserID(c)
	resp, err := h.UseCase.FilterRoadType6(userID)
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

// @summary รายงานข้อมูลความเสียหาย
// @description Get Road Damage Filter Type1
// @tags Report Filter
// @id GetRoadDamageFilterType1
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.FilterRoadDamageType1{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/road_damage/filter/type1 [get]
func (h *Handler) GetRoadDamageFilterType1(c *gin.Context) {
	userID := helpers.GetUserID(c)
	resp, err := h.UseCase.FiltertRoadDamageType1(userID)
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

// @summary รายงานสรุปข้อมูลความเสียหาย
// @description Get Road Damage Filter Type2
// @tags Report Filter
// @id GetRoadDamageFilterType2
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.FilterRoadDamageType2{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/road_damage/filter/type2 [get]
func (h *Handler) GetRoadDamageFilterType2(c *gin.Context) {
	userID := helpers.GetUserID(c)
	resp, err := h.UseCase.FiltertRoadDamageType2(userID)
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

// @summary รายงานประวัติการซ่อมบำรุง KPI
// @description Get Maintenance Kpi Filter Type1
// @tags Report Filter
// @id GetMaintenanceKpiKpiType1
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.FilterMaintenanceKpiType1{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/maintenance_kpi/filter/type1 [get]
func (h *Handler) GetMaintenanceKpiFilterType1(c *gin.Context) {
	userID := helpers.GetUserID(c)
	resp, err := h.UseCase.FilterMaintenanceKpiType1(userID)
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

// @summary รายงานประวัติการซ่อมบำรุง
// @description Get Maintenance Filter Type1
// @tags Report Filter
// @id GetMaintenanceFilterType1
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.FilterMaintenanceType1{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/maintenance/filter/type1 [get]
func (h *Handler) GetMaintenanceFilterType1(c *gin.Context) {
	userID := helpers.GetUserID(c)
	resp, err := h.UseCase.FilterMaintenanceFilterType1(userID)
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

// @summary รายงานปริมาณจราจร
// @description Get Aadt Filter Type1
// @tags Report Filter
// @id GetAadtFilterType1
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.FilterAadtType1{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/aadt/filter/type1 [get]
func (h *Handler) GetAadtFilterType1(c *gin.Context) {
	userID := helpers.GetUserID(c)
	resp, err := h.UseCase.FilterAadtType1(userID)
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

// @summary
// @description Report Status
// @tags Report
// @id ReportStatus
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/status [get]
func (h *Handler) ReportStatus(c *gin.Context) {
	resp, err := h.UseCase.ReportStatus()
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

// @summary
// @description Report Status
// @tags Report
// @id CheckReportStatus
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/check/status/{id} [get]
func (h *Handler) CheckReportStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	resp, err := h.UseCase.CheckReportStatus(id)
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

// @summary รายงานสรุปข้อมูลค่าสะท้อนแสงของเส้นจราจร
// @description Report 9
// @tags Report
// @id Report9
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param year query string true "Insert the year"
// @Param road_section_id query string true "Insert the road section ID"
// @Param filter_criteria_id query string true "Insert the filter criteria  ID"
// @Param type query string true "Insert the report type" default(<html, pdf, excel>)
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/type9 [get]
func (h *Handler) Report9(c *gin.Context) {
	roadSectionId, _ := strconv.Atoi(c.Query("road_section_id"))
	filterCriteriaId, _ := strconv.Atoi(c.Query("filter_criteria_id"))
	year, _ := strconv.Atoi(c.Query("year"))
	typ := c.Query("type")

	resp, err := h.UseCase.Report9(roadSectionId, filterCriteriaId, year, typ)
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

// @summary รายงานค่าสะท้อนแสงของเส้นจราจร
// @description Report 8
// @tags Report
// @id Report8
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param year query string true "Insert the year"
// @Param road_section_id query string true "Insert the road section ID"
// @Param filter_criteria_id query string true "Insert the filter criteria  ID"
// @Param type query string true "Insert the report type" default(<html, pdf, excel>)
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/type8 [get]
func (h *Handler) Report8(c *gin.Context) {
	roadSectionId, _ := strconv.Atoi(c.Query("road_section_id"))
	filterCriteriaId, _ := strconv.Atoi(c.Query("filter_criteria_id"))
	year, _ := strconv.Atoi(c.Query("year"))
	typ := c.Query("type")

	resp, err := h.UseCase.Report8(roadSectionId, filterCriteriaId, year, typ)
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

// @summary รายงานการซ่อมบํารุงตามเกณฑ์ KPI ค่าดัชนีความขรุขระสากล
// @description Report 12
// @tags Report
// @id Report12
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param year query string true "Insert the year"
// @Param road_section_id query string true "Insert the road section ID"
// @Param filter_condition_name query string true "Insert the filter criteria name" default(<IRI, IFI, RUT, G7>)
// @Param type query string true "Insert the report type" default(<html, pdf, excel>)
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/report/type12 [get]
func (h *Handler) Report12(c *gin.Context) {
	roadSectionId, _ := strconv.Atoi(c.Query("road_section_id"))
	filterConditionName := c.Query("filter_condition_name")
	year, _ := strconv.Atoi(c.Query("year"))
	typ := c.Query("type")

	resp, err := h.UseCase.Report12(roadSectionId, year, filterConditionName, typ)
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

//////////////// NEW MIMS ////////////////
