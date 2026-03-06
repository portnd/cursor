package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/maintenanceAnalysis/domains"
	"gopkg.in/validator.v2"
)

type Handler struct {
	Usecase domains.Usecase
}

func NewHandler(usecase domains.Usecase) *Handler {
	return &Handler{Usecase: usecase}
}

// @summary หน้ารายการวิเคราะห์ซ่อมบำรุง
// @description
// @tags Analysis
// @id get_maintenance_analysis
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param type query string false "ประเภท"
// @Param condition query string false "เงื่อนไข"
// @param page query string false "Insert your page number"
// @param limit query string false "Insert your limit number"
// @response 200 {object} responses.Success{data=responses.Pagination{items=[]responses.AnalysisRes}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze [get]
func (h *Handler) GetMaintenanceAnalysis(c *gin.Context) {
	userID := helpers.GetUserID(c)
	var analysisFilter requests.AnalysisFilter

	typeAnalysis, ok := c.Request.URL.Query()["type"] // ปีงบประมาณ
	if ok {
		typeAnalysis, err := strconv.Atoi(typeAnalysis[0])
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(http.StatusBadRequest, errResponse)
			return
		}
		analysisFilter.TypeAnalysis = &typeAnalysis
	}

	condition, ok := c.Request.URL.Query()["condition"] // ปีงบประมาณ
	if ok {
		analysisFilter.Condition = &condition[0]
	}

	limitParam := c.Query("limit")
	pageParam := c.Query("page")
	limit, offset, page := helpers.GetlimitOffsetPage(limitParam, pageParam, 0)

	respData, totalItems, err := h.Usecase.GetMaintenanceAnalysis(userID, analysisFilter, limit, offset)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	_, _, page = helpers.GetlimitOffsetPage(limitParam, pageParam, totalItems)
	pagination := helpers.Pagination(respData, limit, page, totalItems)
	c.JSON(http.StatusOK, responses.SuccessResponse(pagination, http.StatusOK))
}

// @summary หน้ารายละเอียด วิเคราะห์ซ่อมบำรุง (condition + prepare data)
// @description
// @tags Analysis
// @id get_maintenance_analysis_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @response 200 {object} responses.Success{data=responses.AnalysisByIDRes} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/{id} [get]
func (h *Handler) GetMaintenanceAnalysisById(c *gin.Context) {
	userID := helpers.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	resp, err := h.Usecase.GetMaintenanceAnalysisById(id, userID)
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

// @summary ค้นหา prepare data ตามเงือนไข
// @description
// @tags Analysis
// @id CreateMaintenanceAnalysis
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param maintenance_analysis body requests.MaintenanceAnalysis true "Insert your maintenance analysis"
// @response 200 {object} responses.Success{data=responses.AnalysisByIDRes} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/prepare_data [post]
func (h *Handler) CreateMaintenanceAnalysis(c *gin.Context) {
	var req requests.MaintenanceAnalysis
	if err := c.ShouldBind(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if validateErr := validator.Validate(req); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))

	res, err := h.Usecase.CreateMaintenanceAnalysis(userID, req)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(200, responses.SuccessResponse(res, 200))
}

// @summary ค้นหา prepare data ตามเงือนไข กรณีกลับมาแก้ไข
// @description
// @tags Analysis
// @id UpdateMaintenanceAnalysis
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @param maintenance_analysis body requests.MaintenanceAnalysis true "Insert your maintenance analysis"
// @response 200 {object} responses.Success{data=responses.AnalysisByIDRes} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/{id}/prepare_data [put]
func (h *Handler) UpdateMaintenanceAnalysis(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	var req requests.MaintenanceAnalysis
	if err := c.ShouldBind(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))

	if validateErr := validator.Validate(req); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	res, err := h.Usecase.UpdateMaintenanceAnalysis(ID, userID, req)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(200, responses.SuccessResponse(res, 200))
}

// ================================================ STEP 2 ================================================
// @summary get data to step 2
// @description
// @tags Analysis
// @id GetMaintenanceAnalysisCondition
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @param maintenance_analysis body requests.PrepareDataIDReq true "Insert your maintenance analysis"
// @response 200 {object} responses.Success{data=responses.AnalysisStep2Res} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/{id}/condition [post]
func (h *Handler) GetMaintenanceAnalysisCondition(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	var req requests.PrepareDataIDReq
	if err := c.ShouldBind(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	res, err := h.Usecase.GetMaintenanceAnalysisConditionById(ID, req.PrepareDataID)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

// @summary เริ่มวิเคราะห์
// @description ccc
// @tags Analysis
// @id post_maintenance_analysis_strategic
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @param maintenance_analysis_strategic body requests.AnalyzingReq true "Insert your maintenance analysis"
// @response 200 {object} responses.Success{data=responses.AnalysisStep2Res} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/{id}/analyzing [post]
func (h *Handler) CreateMaintenanceAnalysisAnalyzing(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))

	var req requests.AnalyzingReq
	if err := c.ShouldBind(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if validateErr := validator.Validate(req); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	resp, err := h.Usecase.UpdateMaintenanceAnalysisCondition(ID, userID, req)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	err = h.Usecase.SelectdPrepareData(ID, req.PrepareDataID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	h.Usecase.AnalysisModel(ID)

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary ลบข้อมูลวิเคราะห์
// @description
// @tags Analysis
// @id DeleteMaintenanceAnalysis
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @response 204 {object} responses.Success{data=responses.NoData} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/{id} [delete]
func (h *Handler) DeleteMaintenanceAnalysis(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	err = h.Usecase.DeleteMaintenanceAnalysis(ID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	c.JSON(204, responses.SuccessResponse("", 204))
}

// @summary คัดลอก
// @description
// @tags Analysis
// @id CopyMaintenanceAnalysis
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @response 200 {object} responses.Success{data=models.MaintenanceAnalysis} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/{id}/copy [post]
func (h *Handler) CopyMaintenanceAnalysis(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	res, err := h.Usecase.CopyMaintenanceAnalysis(ID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

// @summary บุ๊กมาร์ค
// @description
// @tags Analysis
// @id FavoriteMaintenanceAnalysis
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @response 200 {object} responses.Success{data=responses.AnalysisIsFavorite{}}"OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/{id}/favorite [post]
func (h *Handler) FavoriteMaintenanceAnalysis(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	resp, err := h.Usecase.FavoriteMaintenanceAnalysis(ID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary dashboard
// @description
// @tags Analysis
// @id get_strategic_maintenance_analysis_dashboard
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @response 200 {object} responses.Success{data=models.DashboardStrategicMaintenance} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/dashboard/strategic/{id} [get]
func (h *Handler) DashboardStrategicMaintenanceAnalysis(c *gin.Context) {
	userID := helpers.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	resp, err := h.Usecase.DashboardStrategicMaintenanceAnalysis(strconv.Itoa(id), userID)
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

// @summary dashboard
// @description
// @tags Analysis
// @id get_annual_maintenance_analysis_dashboard
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @response 200 {object} responses.Success{data=models.DashboardAnnualMaintenance} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/dashboard/annual/{id} [get]
func (h *Handler) DashboardAnnualMaintenanceAnalysis(c *gin.Context) {
	userID := helpers.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	resp, err := h.Usecase.DashboardAnnualMaintenanceAnalysis(strconv.Itoa(id), userID)
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

// @summary dashboard
// @description
// @tags Analysis
// @id ExportData
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @Success 200 {object}  responses.Success{data=string{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/{id}/export [get]
func (h *Handler) ExportData(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	resp, err := h.Usecase.ExportData(id)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary หน้ารายละเอียด วิเคราะห์ซ่อมบำรุง
// @description
// @tags Analysis
// @id get_analysis_model
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @response 200 {object} responses.Success{data=responses.Pagination{items=[]responses.AnalysesModel}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/{id}/model [get]
func (h *Handler) GetAnalysisModel(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	var req requests.MaintenanceAnalysisModel
	if err := c.ShouldBind(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	responds, err := h.Usecase.GetAnalysisModel(ID, req)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	limitParam := c.Query("limit")
	pageParam := c.Query("page")

	resp := responds.([]responses.AnalysesModel)
	totalItems := int64(len(resp))
	limit, offset, page := helpers.GetlimitOffsetPage(limitParam, pageParam, totalItems)

	if totalItems == 0 {
		resp = []responses.AnalysesModel{}
	} else if limit+offset > totalItems {
		resp = resp[offset:totalItems]
	} else {
		resp = resp[offset : limit+offset]
	}

	pagination := helpers.Pagination(resp, limit, page, totalItems)

	c.JSON(http.StatusOK, responses.SuccessResponse(pagination, http.StatusOK))

}

// @summary หน้ารายละเอียด วิเคราะห์ซ่อมบำรุง (วิเคราะห์ใหม่)
// @description
// @tags Analysis
// @id update_analysis_model
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param maintenance_analysis body requests.MaintenanceAnalysisModelReq true "Insert your analysis model"
// @param id path string false "Insert your maintenance analysis id"
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/{id}/model [put]
func (h *Handler) UpdateAnalysisModel(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	var req []requests.MaintenanceAnalysisModelReq
	if err := c.ShouldBind(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	resp, err := h.Usecase.UpdateAnalysisModel(ID, req)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	h.Usecase.AnalysisModel(ID)
	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusAccepted))
}

// @summary หน้ารายละเอียด วิเคราะห์ซ่อมบำรุง (วิธีการซ่อมบำรุง)
// @description
// @tags Analysis
// @id get_ref_criteria_method
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/intervention_criterias [get]
func (h *Handler) GetRefCriteriaMethod(c *gin.Context) {
	resp, err := h.Usecase.GetRefCriteriaMethod()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary filter map
// @description
// @tags Analysis Map
// @id DashboardMapFilter
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @Success 200 {object}  responses.Success{data=responses.Filter} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/dashboard-map/{id}/filter [get]
func (h *Handler) DashboardMapFilter(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	resp, err := h.Usecase.DashboardMapFilter(ID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))

}

// @summary  map
// @description
// @tags Analysis Map
// @id DashboardMap
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @Param plan query string false "แผน"
// @Param year query string false "ปี"
// @Param display query string false "การแสดงผล"
// @Param criteria query string false "เกรณฑ์"
// @Success 200 {object}  responses.Success{data=responses.DashboardMap} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/dashboard-map/{id} [get]
func (h *Handler) DashboardMap(c *gin.Context) {
	//
	var mapFilter requests.MapFilter
	year, ok := c.Request.URL.Query()["year"] // ปีงบประมาณ
	if ok {
		year, _ := strconv.Atoi(year[0])
		mapFilter.Year = &year
	}

	plan, ok := c.Request.URL.Query()["plan"] // ปีงบประมาณ
	if ok {
		plan, _ := strconv.Atoi(plan[0])
		mapFilter.Plan = &plan
	}

	display, ok := c.Request.URL.Query()["display"] // ปีงบประมาณ
	if ok {
		display, _ := strconv.Atoi(display[0])
		mapFilter.Display = display
	}

	criteria, ok := c.Request.URL.Query()["criteria"] // ปีงบประมาณ
	if ok {
		criteria, _ := strconv.Atoi(criteria[0])
		mapFilter.Criteria = &criteria
	}

	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	resp, err := h.Usecase.DashboardMap(ID, mapFilter)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))

}

// @summary เช็คสถานะการสร้าง Prepare Data
// @description
// @tags Analysis
// @id check_prepare_data
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @response 200 {object} responses.Success{data=responses.CheckPrepareDataStatus} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/{id}/check_prepare_data [get]
func (h *Handler) CheckPrepareDataById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	resp, err := h.Usecase.CheckPrepareDataById(id)
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

// @summary รายการ Prepare Data Id
// @description
// @tags Analysis
// @id prepare_data_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @response 200 {object} responses.Success{data=[]int} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/{id}/prepare_data_id [get]
func (h *Handler) GetPrepareDataIdById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	resp, err := h.Usecase.GetPrepareDataIdById(id)
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

// @summary รายการ Prepare Data
// @description
// @tags Analysis
// @id prepare_data
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your maintenance analysis id"
// @param page query string false "Insert your page number"
// @param limit query string false "Insert your limit number"
// @response 200 {object} responses.Success{data=responses.Pagination{items=[]responses.PrepareDataWithPagination}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/{id}/prepare_data [get]
func (h *Handler) GetPrepareDataById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	responds, err := h.Usecase.GetPrepareDataById(id)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	limitParam := c.Query("limit")
	pageParam := c.Query("page")

	resp := responds.([]responses.PrepareDataWithPagination)
	totalItems := int64(len(resp))
	limit, offset, page := helpers.GetlimitOffsetPage(limitParam, pageParam, totalItems)

	if totalItems == 0 {
		resp = []responses.PrepareDataWithPagination{}
	} else if limit+offset > totalItems {
		fmt.Println(limit, " + ", offset, " > ", totalItems)

		resp = resp[offset:totalItems]
	} else {
		resp = resp[offset : limit+offset]
	}

	pagination := helpers.Pagination(resp, limit, page, totalItems)

	c.JSON(http.StatusOK, responses.SuccessResponse(pagination, http.StatusOK))
	return
}

// @Summary ID ของ prepare data selected
// @description
// @tags Analysis
// @id GetPrepareDataAllByAnalysisSelected
// @Security Bearer
// @param id path string true "Maintenance Analysis ID"
// @response 200 {object} responses.Success{data=[]int{}}  "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/analyze/{id}/prepare_data_selected [get]
func (h *Handler) GetPrepareDataAllByAnalysisSelected(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	res, err := h.Usecase.GetPrepareDataAllByAnalysisSelected(ID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusNotFound, errResponse)
	}
	// res.PrepareData = prepareData
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}
