package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/maintenance/domains"
	"gopkg.in/validator.v2"
)

type Handler struct {
	Usecase domains.Usecase
}

func NewHandler(usecase domains.Usecase) *Handler {
	return &Handler{Usecase: usecase}
}

// @Summary สถานะของแผนดำเนินการ
// @Description
// @tags Maintenance Plan
// @id GetMaintenancePlanStatus
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object}  responses.Success{data=[]models.RefMaintenancePlanStatus}  "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/plan_stauts [get]
// func (h *Handler) GetMaintenancePlanStatus(c *gin.Context) {
// 	data, err := h.Usecase.GetMaintenancePlanStatus()
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusUnprocessableEntity, errResponse)
// 		return
// 	}
// 	c.JSON(200, responses.SuccessResponse(data, 200))
// }

// @Summary ประเภทงบประมาณ
// @Description
// @tags Maintenance
// @id GetMaintenanceBudget
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object}  responses.Success{data=[]models.RefMaintenancePlanStatus}  "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/budgets [get]
func (h *Handler) GetMaintenanceBudget(c *gin.Context) {
	data, err := h.Usecase.GetMaintenanceBudget()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @Summary วิธีการซ่อมบำรุง
// @Description
// @tags Maintenance
// @id GetInterventionCriteria
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object}  responses.Success{data=[]models.RefCriteriaMethod}  "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/intervention_criteria [get]
func (h *Handler) GetInterventionCriteria(c *gin.Context) {
	data, err := h.Usecase.GetRefCriteriaMethod()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary Road Dropdown List
// @description
// @tags Maintenance
// @id GetRoadDropdownList
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object}  responses.Success{data=[]responses.RoadListInitData}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/maintenance/road_dropdown_list [get]
func (h *Handler) GetRoadDropdownList(c *gin.Context) {
	userID := helpers.GetUserID(c)
	data, err := h.Usecase.GetRoadDropdownList(userID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary Road Dropdown List
// @description
// @tags Maintenance
// @id GetRoadDropdownListAnalyze
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object}  responses.Success{data=[]responses.RoadListInitData}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/maintenance/road_dropdown_list_analyze [get]
func (h *Handler) GetRoadDropdownListAnalyze(c *gin.Context) {
	userID := helpers.GetUserID(c)
	data, err := h.Usecase.GetRoadDropdownListAnalyze(userID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary Road Dropdown List
// @description
// @tags Maintenance
// @id GetRoadDropdownListDashboard
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object}  responses.Success{data=[]responses.RoadListInitData}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/maintenance/road_dropdown_list_dashboard [get]
func (h *Handler) GetRoadDropdownListDashboard(c *gin.Context) {
	userID := helpers.GetUserID(c)
	typeData := c.Query("type")
	ownerCode := c.Query("owner_code")
	helpers.PrintlnJson(typeData)
	data, err := h.Usecase.GetRoadDropdownListDashboard(userID, typeData, ownerCode)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(200, responses.SuccessResponse(data, 200))
}

// // @summary Road Dropdown List
// // @description
// // @tags Maintenance
// // @id GetRoadDivisionFilter
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @Success 200 {object}  responses.Success{data=[]responses.RoadListInitData}  "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/maintenance/road_division_filter [get]
// func (h *Handler) GetRoadDivisionFilter(c *gin.Context) {
// 	data, err := h.Usecase.GetRoadDivisionFilter()
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(errResponse.Code, errResponse)
// 		return
// 	}

// 	c.JSON(200, responses.SuccessResponse(data, 200))
// }

// @Summary วิธีการซ่อมบำรุง
// @Description
// @tags Maintenance
// @id GetDivisionList
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object}  responses.Success{data=[]models.RefDivisionList}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/division [get]
func (h *Handler) GetDivisionList(c *gin.Context) {
	userID := helpers.GetUserID(c)
	data, err := h.Usecase.GetDivisionList(userID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
}

// @Summary วิธีการซ่อมบำรุง หน้าสายทาง
// @Description
// @tags Maintenance
// @id GetMaintenanceByRoadID
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @Param year query string false "ปีงบประมาณ"
// @Success 200 {object}  responses.Success{data=[]responses.MaintenancHistoryeList}  "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/roads/{id}/maintenance [get]
func (m *Handler) GetMaintenanceByRoadID(c *gin.Context) {
	roadID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	year := 0
	yearParam, ok := c.Request.URL.Query()["year"] // ปีงบประมาณ
	if ok {
		year, err = strconv.Atoi(yearParam[0])
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(http.StatusBadRequest, errResponse)
			return
		}
	}

	responds, err := m.Usecase.GetMaintenanceByRoadID(roadID, year)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(responds, 200))
}

// @Summary ปี วิธีการซ่อมบำรุง หน้าสายทาง
// @Description
// @tags Maintenance
// @id GetMaintenanceByRoadYear
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @Success 200 {object}  responses.Success{data=[]responses.IntRes}  "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/roads/{id}/maintenance_year [get]
func (m *Handler) GetMaintenanceByRoadYear(c *gin.Context) {
	roadID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	responds, err := m.Usecase.GetRoadMaintenanceYear(roadID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(responds, 200))
}

// @Summary รายการติดตามซ่อมบำรุง
// @Description
// @tags Maintenance
// @id GetMaintenance
// @Accept  json
// @Produce  json
// @Param road_group_id query string false "สายทาง"
// @Param road_group_id_dashboard query string false "สายทาง action dashboard"
// @Param owner_code query string false "หน่วยงานที่รับผิดชอบ"
// @Param budget_year query string false "ปีงบประมาณ"
// @Param budget_method_id query string false "ประเภทงบประมาณ 1,2,3,4"
// @Param name query string false "ชื่อโครงการ"
// @param page query string true "Insert your page number"
// @param limit query string true "Insert your limit number"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.Pagination{items=[]responses.MaintenanceList}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance [get]
func (m *Handler) GetMaintenance(c *gin.Context) {
	userID := helpers.GetUserID(c)
	var params requests.MaintenancePrams
	owerCode, ok := c.Request.URL.Query()["owner_code"]
	if ok {
		owerCode := owerCode[0]
		params.OwnerCode = &owerCode
	}

	budgetYear, ok := c.Request.URL.Query()["budget_year"] // ปีงบประมาณ
	if ok {
		budgetYear := helpers.StrToInt(budgetYear[0])
		params.BudgetYear = &budgetYear
	}

	budgetTypeID, ok := c.Request.URL.Query()["budget_method_id"] //ประเภทงบประมาณ
	if ok {
		idSplit := strings.Split(budgetTypeID[0], ",")
		arrayInt, err := helpers.ConvertToArrayInt(idSplit)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		params.BudgetMethodId = arrayInt
	}

	roadGroupId, ok := c.Request.URL.Query()["road_group_id"]
	if ok {
		idSplit := strings.Split(roadGroupId[0], ",")
		arrayInt, err := helpers.ConvertToArrayInt(idSplit)
		if err != nil {
			arrayInt = []int{}
		}
		params.RoadGroupID = arrayInt
	}

	roadGroupIDDashboard, ok := c.Request.URL.Query()["road_group_id_dashboard"]
	if ok {
		idSplit := strings.Split(roadGroupIDDashboard[0], ",")
		arrayInt, err := helpers.ConvertToArrayInt(idSplit)
		if err != nil {
			arrayInt = []int{}
		}
		params.RoadGroupIDDashboard = arrayInt
	}

	name, ok := c.Request.URL.Query()["name"] //ชื่อโครการ
	if ok {
		name := name[0]
		params.Name = &name
	}
	limitParam := c.Query("limit")
	pageParam := c.Query("page")
	limit, offset, page := helpers.GetlimitOffsetPage(limitParam, pageParam, 0)
	respData, totalItems, err := m.Usecase.GetMaintenanceList(userID, params, limit, offset)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	pagination := helpers.Pagination(respData, limit, page, totalItems)
	c.JSON(http.StatusOK, responses.SuccessResponse(pagination, http.StatusOK))
}

// @Summary รายละเอียดการติดตามซ่อมบำรุง
// @description
// @tags Maintenance
// @id GetMaintenanceByID
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your maintenance ID Parent"
// @response 200 {object} []responses.Success{data=responses.MaintenanceList}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/{id_parent} [get]
func (m *Handler) GetMaintenanceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id_parent"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}
	responds, err := m.Usecase.GetMaintenanceListByID(id)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusNotFound, errResponse)
		return
	}
	// if len(responds) == 0 {
	// 	c.JSON(http.StatusOK, responses.SuccessResponse([]string{}, http.StatusOK))
	// 	return
	// }

	c.JSON(http.StatusOK, responses.SuccessResponse(responds, http.StatusOK))

}

// @Summary เพิ่มโครงการติดตามซ่อมบำรุง
// @description
// @tags Maintenance
// @id CreateMaintenance
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param maintenance body requests.MaintenanceReq true "JSON object that represents the maintenance project to be created"
// @Success 201 {object} responses.Success{data=responses.MaintenanceID} "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance [post]
func (h *Handler) CreateMaintenance(c *gin.Context) {
	var req requests.MaintenanceReq
	if err := c.ShouldBind(&req); err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	validateMsgArr := map[string]string{}
	if validateErr := validator.Validate(req); validateErr != nil {
		validateMsgArr = helpers.ConverstError(validateErr)
	}

	if req.BudgetMethodID == 0 {
		validateMsgArr["budget_method_id"] = "โปรดระบุ"
	}

	if req.BudgetProcurement == nil {
		validateMsgArr["budget_procurement"] = "โปรดระบุ"
	}

	isDup, err := h.Usecase.CheckMaintenanceDuplicate(0, req.Name)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	if isDup {
		validateMsgArr["name"] = "ชื่อโครงการมีอยุ่แล้ว"
	}

	if len(validateMsgArr) > 0 {
		errResponse := responses.ValidateResponse(validateMsgArr)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	userId := helpers.GetUserInfo(c).UserId

	attachments := req.Attachments
	var maintenanceAttachmentsReq []requests.MaintenanceAttachmentsReq
	for i, item := range attachments {
		switch item.Status {
		case "delete":
			maintenanceAttachmentsReq = append(maintenanceAttachmentsReq, requests.MaintenanceAttachmentsReq{ID: item.ID, FileName: item.FileName, Path: "", FileType: "", Status: item.Status})
		case "upload":
			imgBase64 := item.File
			base64 := ""
			if strings.Contains(imgBase64, ",") {
				dataImgBase64 := helpers.Explode(",", imgBase64)
				base64 = dataImgBase64[1]
			} else {
				base64 = imgBase64
			}
			pathOutput := os.Getenv("MAINTENANCE_ROAD_HISTORY")

			err := helpers.EnsureDir(pathOutput)
			if err != nil {
				appErr, ok := err.(*responses.AppErr)
				if !ok {
					// Handle the case where err is not *responses.AppErr
					c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
					return
				}
				errResponse := responses.FailRespone(appErr)
				c.JSON(http.StatusBadRequest, errResponse)
				return
			}

			fileName := "maintenance_" + time.Now().Format("20060102150405") + "_" + strconv.Itoa(i)
			fileNames := helpers.Explode(".", item.FileName)
			fileNameType := ""
			if len(fileNames) > 1 {
				fileNameType = fileNames[len(fileNames)-1]
			} else {
				errResponse := responses.FailRespone(fmt.Errorf("ชื่อไฟล์ไม่ถูกต้อง"))
				c.JSON(http.StatusBadRequest, errResponse)
				return
			}

			filePath, fileType, err := helpers.DecodeFileBase64FileType(base64, pathOutput, fileName, fileNameType, "file")
			if err != nil {
				errResponse := responses.FailRespone(err)
				c.JSON(http.StatusBadRequest, errResponse)
				return
			}
			maintenanceAttachmentsReq = append(maintenanceAttachmentsReq, requests.MaintenanceAttachmentsReq{ID: item.ID, FileName: item.FileName, Path: filePath, FileType: fileType, Status: item.Status})
		}
	}

	id, err := h.Usecase.CreateMaintenance(req, userId, maintenanceAttachmentsReq)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(201, responses.SuccessResponse(responses.MaintenanceID{ID: id.(int), IDParent: id.(int)}, 201))
}

// @Summary แก้ไขโครงการติดตามซ่อมบำรุง
// @description
// @tags Maintenance
// @id UpdateMaintenance
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your maintenance ID Parent"
// @Param maintenance body requests.MaintenanceReq true "JSON object that represents the maintenance project to be created"
// @Success 202 {object} responses.Success{data=responses.MaintenanceID} "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/{id_parent} [put]
func (h *Handler) UpdateMaintenance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id_parent"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}

	var req requests.MaintenanceReq
	if err := c.ShouldBind(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	validateMsgArr := map[string]string{}
	if validateErr := validator.Validate(req); validateErr != nil {
		validateMsgArr = helpers.ConverstError(validateErr)
	}

	if req.BudgetMethodID == 0 {
		validateMsgArr["budget_method_id"] = "โปรดระบุ"
	}

	if req.BudgetProcurement == nil {
		validateMsgArr["budget_procurement"] = "โปรดระบุ"
	}

	isDup, err := h.Usecase.CheckMaintenanceDuplicate(id, req.Name)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	if isDup {
		validateMsgArr["name"] = "ชื่อโครงการมีอยุ่แล้ว"
	}

	if len(validateMsgArr) > 0 {
		errResponse := responses.ValidateResponse(validateMsgArr)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	attachments := req.Attachments
	var maintenanceAttachmentsReq []requests.MaintenanceAttachmentsReq
	for i, item := range attachments {
		switch item.Status {
		case "delete":

			if item.ID == nil {

				err := responses.NewAppErr(400, "โปรดระบุ ID ของไฟล์ที่ต้องการลบ")
				errResponse := responses.FailRespone(err)
				c.JSON(http.StatusBadRequest, errResponse)
				return
			}

			att, err := h.Usecase.GetMaintenanceaAttachmentByID(*item.ID)
			if err != nil {
				errResponse := responses.FailRespone(err)
				c.JSON(400, errResponse)
				return
			}
			filePath := att.Path
			if filePath != "" {
				err = helpers.DeleteFile(filePath)
				if err != nil {
					c.JSON(400, err)
					return
				}
			}
			maintenanceAttachmentsReq = append(maintenanceAttachmentsReq, requests.MaintenanceAttachmentsReq{ID: item.ID, FileName: item.FileName, Path: "", FileType: "", Status: item.Status})
		case "upload":
			imgBase64 := item.File
			base64 := ""
			if strings.Contains(imgBase64, ",") {
				dataImgBase64 := helpers.Explode(",", imgBase64)
				base64 = dataImgBase64[1]
			} else {
				base64 = imgBase64
			}

			pathOutput := os.Getenv("MAINTENANCE_ROAD_HISTORY")
			err := helpers.EnsureDir(pathOutput)
			if err != nil {
				appErr, ok := err.(*responses.AppErr)
				if !ok {
					// Handle the case where err is not *responses.AppErr
					c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
					return
				}
				errResponse := responses.FailRespone(appErr)
				c.JSON(http.StatusBadRequest, errResponse)
				return
			}

			fileName := "maintenance_" + time.Now().Format("20060102150405") + "_" + strconv.Itoa(i)
			fileNames := helpers.Explode(".", item.FileName)
			fileNameType := ""
			if len(fileNames) > 1 {
				fileNameType = fileNames[len(fileNames)-1]
			} else {
				errResponse := responses.FailRespone(fmt.Errorf("ชื่อไฟล์ไม่ถูกต้อง"))
				c.JSON(http.StatusBadRequest, errResponse)
				return
			}

			filePath, fileType, err := helpers.DecodeFileBase64FileType(base64, pathOutput, fileName, fileNameType, "")
			if err != nil {
				appErr, ok := err.(*responses.AppErr)
				if !ok {
					// Handle the case where err is not *responses.AppErr
					c.JSON(http.StatusBadRequest, responses.NewAppErr(http.StatusBadRequest, err.Error()))
					return
				}
				errResponse := responses.FailRespone(appErr)
				c.JSON(http.StatusBadRequest, errResponse)
				return
			}
			maintenanceAttachmentsReq = append(maintenanceAttachmentsReq, requests.MaintenanceAttachmentsReq{ID: item.ID, FileName: item.FileName, Path: filePath, FileType: fileType, Status: item.Status})
		}
	}

	userId := helpers.GetUserInfo(c).UserId
	newId, idParent, err := h.Usecase.UpdateMaintenance(id, req, userId, maintenanceAttachmentsReq)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(202, responses.SuccessResponse(responses.MaintenanceID{ID: *newId, IDParent: *idParent}, 202))

}

// @Summary ลบโครงการติดตามซ่อมบำรุง
// @description
// @tags Maintenance
// @id DeleteMaintenance
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your maintenance ID Parent"
// @response 204 {object} responses.NoDataResponse "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/{id_parent} [delete]
func (m *Handler) DeleteMaintenance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id_parent"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}
	err = m.Usecase.DeleteMaintenance(id)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.Status(http.StatusNoContent)
}

// // @Summary สิ้นสุดโครงการ
// // @description
// // @tags Maintenance
// // @id MaintenanceFinished
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @param id path string true "Insert your maintenance ID Parent"
// // @Param maintenance body requests.MaintenanceFinished true "JSON object that represents the maintenance project to be created"
// // @Success 200 {object} responses.Success{} "OK"
// // @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// // @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @Router /api/v1/maintenance/{id_parent}/finished [post]
// func (h *Handler) MaintenanceFinished(c *gin.Context) {
// 	ID, err := strconv.Atoi(c.Param("id_parent"))
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusBadRequest, errResponse)
// 		return
// 	}
// 	var req requests.MaintenanceFinished
// 	if err := c.ShouldBind(&req); err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusBadRequest, errResponse)
// 		return
// 	}

// 	validateMsgArr := map[string]string{}
// 	if validateErr := validator.Validate(req); validateErr != nil {
// 		validateMsgArr = helpers.ConverstError(validateErr)
// 	}
// 	if len(validateMsgArr) > 0 {
// 		errResponse := responses.ValidateResponse(validateMsgArr)
// 		c.JSON(http.StatusUnprocessableEntity, errResponse)
// 		return
// 	}

// 	err = h.Usecase.MaintenanceFinished(ID, req)
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusBadRequest, errResponse)
// 		return
// 	}
// 	c.JSON(200, responses.SuccessResponse("", 200))
// }

// @Summary สถานะโครงการ โดยเลือกตามแผน หน้ารายละเอียด
// @Description
// @tags Maintenance Plan
// @id GetMaintenanceStatus
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your maintenance ID"
// @Success 200 {object}  responses.Success{data=[]responses.MaintenanceStatus}  "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/{id}/plan_stauts [get]
// func (h *Handler) GetMaintenanceStatus(c *gin.Context) {
// 	ID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusBadRequest, errResponse)
// 		return
// 	}
// 	data, err := h.Usecase.GetMaintenanceStatus(ID)
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusUnprocessableEntity, errResponse)
// 		return
// 	}
// 	c.JSON(200, responses.SuccessResponse(data, 200))
// }

// @Summary รายละเอียดข้อมูลซ่อมบำรุง
// @Description
// @tags Maintenance Road
// @id GetMaintenanceRoadByID
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id_parent path string true "Insert your maintenance ID Parent"
// @param m_road_id path string true "Insert your maintenance road ID"
// @response 200 {object} []responses.Success{data=responses.MaintenanceRoadPreload}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/{id_parent}/road/{m_road_id} [get]
func (m *Handler) GetMaintenanceRoadByID(c *gin.Context) {
	idParent, err := strconv.Atoi(c.Param("id_parent"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}

	mRoadID, err := strconv.Atoi(c.Param("m_road_id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return

	}
	responds, err := m.Usecase.GetMaintenanceRoadByID(idParent, mRoadID)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(responds, http.StatusOK))

}

// @Summary เพิ่มข้อมูลประวัติซ่อมบำรุง
// @description
// @tags Maintenance Road
// @id CreateMaintenanceRoad
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your maintenance Parent"
// @Param maintenance body requests.MaintenanceRoadsReq true "JSON object that represents the maintenance project to be created"
// @Success 201 {object} responses.Success{data=responses.MaintenanceID} "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/{id_parent}/road [post]
func (h *Handler) CreateMaintenanceRoad(c *gin.Context) {
	idParent, err := strconv.Atoi(c.Param("id_parent"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}
	userID := helpers.GetUserInfo(c).UserId
	var req requests.MaintenanceRoadsReq
	if err := c.ShouldBind(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	var checkType int

	roadInfo, err := h.Usecase.GetLastRoadInfoByID(*req.RoadID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusNotFound, errResponse)
		return
	}
	checkType = roadInfo.RefDirectionId

	var allErrors map[string]string = make(map[string]string)
	if validateErr := validator.Validate(req); validateErr != nil {
		errs := helpers.ConverstError(validateErr)
		for k, v := range errs {
			allErrors[k] = v
		}
	}

	validateErr := validateCreate(req, checkType)
	if validateErr != nil {
		errs := helpers.ConverstError(validateErr)
		for k, v := range errs {
			allErrors[k] = v
		}
	}

	checkValidateIsMethod, err := h.Usecase.CheckValidateIsMethod(idParent)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if (checkValidateIsMethod && req.InterventionCriteriaID == nil) || (checkValidateIsMethod && *req.InterventionCriteriaID == 0) {
		allErrors["intervention_criteria_id"] = "โปรดระบุ"
	}

	if len(allErrors) > 0 {
		errResponse := responses.ValidateResponse(allErrors)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	data, err := h.Usecase.CreateMaintenanceRoad(idParent, userID, req)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	id, _ := data.(*int)

	c.JSON(201, responses.SuccessResponse(responses.MaintenanceRoadID{ID: *id, IDParent: *id}, 201))
}

func validateCreate(req requests.MaintenanceRoadsReq, checkType int) error {
	errMsg := []string{}

	if *req.KmStart == 0 && *req.KmEnd == 0 {
		errMsg = append(errMsg, "km_start: zero value")
		errMsg = append(errMsg, "km_end: zero value")
	}

	switch req.MaintenanceType {
	case 1, 2: // Handles both case 1 and 2 as they both need to check LaneNo
		if req.LaneNo == nil || *req.LaneNo == 0 {
			errMsg = append(errMsg, "lane_no: zero value")
		}

		if req.MaintenanceType == 2 { // Additional check for GridNo in case 2 only
			if req.GridNo == nil || *req.GridNo == 0 {
				errMsg = append(errMsg, "grid_no: zero value")
			}
		}
	default:
		if req.LaneNo == nil || *req.LaneNo == 0 {
			errMsg = append(errMsg, "lane_no: zero value")
		}

		if req.GridNo == nil || *req.GridNo == 0 {
			errMsg = append(errMsg, "grid_no: zero value")
		}
	}

	if (checkType == 1) && *req.KmStart > *req.KmEnd {
		errMsg = append(errMsg, "km_start: incorrect")
		errMsg = append(errMsg, "km_end: incorrect")
	} else if (checkType == 2) && *req.KmStart < *req.KmEnd {
		errMsg = append(errMsg, "km_start: incorrect")
		errMsg = append(errMsg, "km_end: incorrect")
	}

	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, ","))
	}

	return nil
}

// @Summary ปรับปรุงข้อมูลประวัติซ่อมบำรุง
// @description
// @tags Maintenance Road
// @id UpdateMaintenanceRoad
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your maintenance ID"
// @param m_road_id path string true "Insert your maintenance road ID Parent"
// @Param maintenance body requests.MaintenanceRoadsReq true "JSON object that represents the maintenance project to be created"
// @Success 202 {object} responses.Success{data=responses.MaintenanceRoadID} "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/{id_parent}/road/{m_road_id} [put]
func (h *Handler) UpdateMaintenanceRoad(c *gin.Context) {
	idParent, err := strconv.Atoi(c.Param("id_parent"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}
	mRoadID, err := strconv.Atoi(c.Param("m_road_id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}
	userID := helpers.GetUserInfo(c).UserId
	var req requests.MaintenanceRoadsReq
	if err := c.ShouldBind(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	var checkType int

	roadInfo, err := h.Usecase.GetLastRoadInfoByID(*req.RoadID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusNotFound, errResponse)
		return
	}
	checkType = roadInfo.RefDirectionId

	var allErrors map[string]string = make(map[string]string)
	if validateErr := validator.Validate(req); validateErr != nil {
		errs := helpers.ConverstError(validateErr)
		for k, v := range errs {
			allErrors[k] = v
		}
	}

	validateErr := validateUpdate(req, checkType)
	if validateErr != nil {
		errs := helpers.ConverstError(validateErr)
		for k, v := range errs {
			allErrors[k] = v
		}
	}

	checkValidateIsMethod, err := h.Usecase.CheckValidateIsMethod(idParent)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if (checkValidateIsMethod && req.InterventionCriteriaID == nil) || (checkValidateIsMethod && *req.InterventionCriteriaID == 0) {
		allErrors["intervention_criteria_id"] = "โปรดระบุ"
	}

	if len(allErrors) > 0 {
		errResponse := responses.ValidateResponse(allErrors)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	id, _, err := h.Usecase.UpdateMaintenanceRoad(idParent, mRoadID, userID, req)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(202, responses.SuccessResponse(responses.MaintenanceRoadID{ID: *id, IDParent: idParent}, 202))
}

func validateUpdate(req requests.MaintenanceRoadsReq, checkType int) error {
	errMsg := []string{}

	if *req.KmStart == 0 && *req.KmEnd == 0 {
		errMsg = append(errMsg, "km_start: zero value")
		errMsg = append(errMsg, "km_end: zero value")
	}

	switch req.MaintenanceType {
	case 1, 2: // Handles both case 1 and 2 as they both need to check LaneNo
		if req.LaneNo == nil || *req.LaneNo == 0 {
			errMsg = append(errMsg, "lane_no: zero value")
		}

		if req.MaintenanceType == 2 { // Additional check for GridNo in case 2 only
			if req.GridNo == nil || *req.GridNo == 0 {
				errMsg = append(errMsg, "grid_no: zero value")
			}
		}
	default:
		if req.LaneNo == nil || *req.LaneNo == 0 {
			errMsg = append(errMsg, "lane_no: zero value")
		}

		if req.GridNo == nil || *req.GridNo == 0 {
			errMsg = append(errMsg, "grid_no: zero value")
		}
	}

	if (checkType == 1) && *req.KmStart > *req.KmEnd {
		errMsg = append(errMsg, "km_start: incorrect")
		errMsg = append(errMsg, "km_end: incorrect")
	} else if (checkType == 2) && *req.KmStart < *req.KmEnd {
		errMsg = append(errMsg, "km_start: incorrect")
		errMsg = append(errMsg, "km_end: incorrect")
	}

	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, ","))
	}

	return nil
}

// @Summary ลบข้อมูลการประวัติซ่อมบำรุง
// @description
// @tags Maintenance Road
// @id DeletedMaintenanceRoad
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your maintenance ID Parent"
// @param m_road_id path string true "Insert your maintenance road ID"
// @Success 204 {object} responses.Success{data=responses.NoData{}} "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/{id_parent}/road/{m_road_id} [delete]
func (h *Handler) DeleteMaintenanceRoad(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id_parent"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}
	hisID, err := strconv.Atoi(c.Param("m_road_id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}

	err = h.Usecase.DeleteMaintenanceRoad(ID, hisID)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(errResponse.Code, errResponse)
		return
	}
	c.JSON(204, responses.SuccessResponse(responses.NoData{}, 204))
}

// /////////////////////////////// maintenance_history  /////////////////////////////////
// func (h *Handler) GetMaintenanceHistory(c *gin.Context) {
// 	maintenanceID, _ := strconv.Atoi(c.Param("id"))
// 	data, err := h.Usecase.GetMaintenanceHistory(maintenanceID)
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusBadRequest, errResponse)
// 		return
// 	}
// 	c.JSON(200, responses.SuccessResponse(data, 200))
// }

// // @Summary รายการประวัติติดตามซ่อมบำรุง
// // @Description
// // @tags Maintenance History
// // @id GetMaintenanceHistory
// // @Accept  json
// // @Produce  json
// // @Param road_group_id query string false "สายทาง"
// // @Param budget_year query string false "ปีงบประมาณ"
// // @Param budget_type_id query string false "ประเภทงบประมาณ"
// // @Param budget_maintenance query string false "งบประมาณการซ่อมบำรุง"
// // @Param name query string false "ชื่อโครงการ"
// // @param page query string true "Insert your page number"
// // @param limit query string true "Insert your limit number"
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.Success{data=responses.Pagination{items=[]responses.MaintenanceList}} "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// // @response 422 {object} responses.Validate "Unprocessable Entity"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @Router /api/v1/maintenance/history [get]
// func (m *Handler) GetMaintenanceHistory(c *gin.Context) {
// 	// ID, err := strconv.Atoi(c.Param("id_parent"))
// 	// if err != nil {
// 	// 	c.JSON(http.StatusBadRequest, err)
// 	// 	return
// 	// }
// 	var params requests.MaintenancePrams
// 	budgetYear, ok := c.Request.URL.Query()["budget_year"] // ปีงบประมาณ
// 	if ok {
// 		budgetYear := helpers.StrToInt(budgetYear[0])
// 		params.BudgetYear = &budgetYear
// 	}

// 	budgetTypeID, ok := c.Request.URL.Query()["budget_id"] //ประเภทงบประมาณ
// 	if ok {
// 		idSplit := strings.Split(budgetTypeID[0], ",")
// 		arrayInt, err := helpers.ConvertToArrayInt(idSplit)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, err)
// 			return
// 		}
// 		params.BudgetMethodId = arrayInt
// 	}

// 	roadGroupId, ok := c.Request.URL.Query()["road_group_id"]
// 	if ok {
// 		idSplit := strings.Split(roadGroupId[0], ",")
// 		arrayInt, err := helpers.ConvertToArrayInt(idSplit)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, err)
// 			return
// 		}
// 		params.RoadGroupID = arrayInt
// 	}

// 	name, ok := c.Request.URL.Query()["name"] //ชื่อโครการ
// 	if ok {
// 		name := name[0]
// 		params.Name = &name
// 	}
// 	// c.JSON(200, params)
// 	// return
// 	// maintainName := c.Query("name")
// 	limitParam := c.Query("limit")
// 	pageParam := c.Query("page")
// 	responds, err := m.Usecase.GetMaintenanceListHistory(params)
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusInternalServerError, errResponse)
// 		return
// 	}

// 	// if len(responds) == 0 {
// 	// 	c.JSON(http.StatusOK, responses.SuccessResponse([]string{}, http.StatusOK))
// 	// 	return
// 	// }
// 	respData := responds.([]responses.MaintenanceList)
// 	totalItems := int64(len(respData))
// 	limit, offset, page := helpers.GetlimitOffsetPage(limitParam, pageParam, totalItems)

// 	if totalItems == 0 {
// 		respData = []responses.MaintenanceList{}
// 	} else if limit+offset > totalItems {
// 		respData = respData[offset:totalItems]
// 	} else {
// 		respData = respData[offset : limit+offset]
// 	}

// 	pagination := helpers.Pagination(respData, limit, page, totalItems)

// 	c.JSON(http.StatusOK, responses.SuccessResponse(pagination, http.StatusOK))
// }

// @Summary รายละเอียดประวัติติดตามซ่อมบำรุง
// @Description
// @tags Maintenance History
// @id GetMaintenanceHistoryByID
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your maintenance ID Parent"
// @response 200 {object} []responses.Success{data=responses.MaintenanceList}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/{id_parent}/road_history/{m_road_his_id} [get]
func (m *Handler) GetMaintenanceHistoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id_parent"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}

	mRoadHisId, err := strconv.Atoi(c.Param("m_road_his_id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}

	responds, err := m.Usecase.GetMaintenanceHistoryByID(id, mRoadHisId)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(responds, http.StatusOK))

}

// // @Summary รายละเอียดประวัติติดตามซ่อมบำรุง
// // @Description
// // @tags Maintenance History
// // @id GetMaintenanceYearHistory
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @param id path string true "Insert your maintenance ID Parent"
// // @response 200 {object} []responses.Success{data=[]responses.IntRes}  "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// // @response 422 {object} responses.Validate "Unprocessable Entity"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @Router /api/v1/maintenance/history/{id_parent}/year [get]
// func (m *Handler) GetMaintenanceYearHistory(c *gin.Context) {
// 	ID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusInternalServerError, errResponse)
// 		return
// 	}
// 	// var params requests.MaintenancePrams
// 	data, err := m.Usecase.GetMaintenanceYearHistory(ID)
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusInternalServerError, errResponse)
// 		return
// 	}
// 	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
// }

// @Summary ประวิติการซ่อมในช่วงประกัน
// @description
// @tags Maintenance History
// @id CreateMaintenanceRoadHistory
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your maintenance road history ID Parent"
// @Param maintenance body requests.MaintenanceRoadsReq true "JSON object that represents the maintenance project to be created"
// @Success 201 {object} responses.Success{data=responses.MaintenanceRoadHisID{}} "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/{id_parent}/road_history [post]
func (h *Handler) CreateMaintenanceRoadHistory(c *gin.Context) {
	idParent, err := strconv.Atoi(c.Param("id_parent"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}

	userID := helpers.GetUserInfo(c).UserId
	var req requests.MaintenanceRoadHistoryReq
	if err := c.ShouldBind(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	var checkType int

	roadInfo, err := h.Usecase.GetLastRoadInfoByID(*req.RoadID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusNotFound, errResponse)
		return
	}
	checkType = roadInfo.RefDirectionId

	var allErrors map[string]string = make(map[string]string)
	if validateErr := validator.Validate(req); validateErr != nil {
		errs := helpers.ConverstError(validateErr)
		for k, v := range errs {
			allErrors[k] = v
		}
	}

	validateErr := validateCreateRoadHis(req, checkType)
	if validateErr != nil {
		errs := helpers.ConverstError(validateErr)
		for k, v := range errs {
			allErrors[k] = v
		}
	}

	checkValidateIsMethod, err := h.Usecase.CheckValidateIsMethod(idParent)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if (checkValidateIsMethod && req.InterventionCriteriaID == nil) || (checkValidateIsMethod && *req.InterventionCriteriaID == 0) {
		allErrors["intervention_criteria_id"] = "โปรดระบุ"
	}

	if len(allErrors) > 0 {
		errResponse := responses.ValidateResponse(allErrors)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	// attachments := req.Attachments
	// var maintenanceAttachmentsReq []requests.MaintenanceAttachmentsReq
	// for i, item := range attachments {
	// 	switch item.Status {
	// 	case "delete":
	// 		maintenanceAttachmentsReq = append(maintenanceAttachmentsReq, requests.MaintenanceAttachmentsReq{ID: item.ID, FileName: item.FileName, Path: "", FileType: "", Status: item.Status})
	// 	case "upload":
	// 		imgBase64 := item.File
	// 		base64 := ""
	// 		if strings.Contains(imgBase64, ",") {
	// 			dataImgBase64 := helpers.Explode(",", imgBase64)
	// 			base64 = dataImgBase64[1]
	// 		} else {
	// 			base64 = imgBase64
	// 		}
	// 		pathOutput := os.Getenv("MAINTENANCE_HIS_FILE_DIR")
	// 		fileName := "maintenance_his" + time.Now().Format("20060102150405") + "_" + strconv.Itoa(i)
	// 		fileNames := helpers.Explode(".", item.FileName)
	// 		fileNameType := ""
	// 		if len(fileNames) > 1 {
	// 			fileNameType = fileNames[len(fileNames)-1]
	// 		} else {
	// 			errResponse := responses.FailRespone(fmt.Errorf("ชื่อไฟล์ไม่ถูกต้อง"))
	// 			c.JSON(http.StatusBadRequest, errResponse)
	// 			return
	// 		}

	// 		filePath, fileType, err := helpers.DecodeFileBase64FileType(base64, pathOutput, fileName, fileNameType, "img")
	// 		if err != nil {
	// 			errResponse := responses.FailRespone(err)
	// 			c.JSON(http.StatusBadRequest, errResponse)
	// 			return
	// 		}
	// 		maintenanceAttachmentsReq = append(maintenanceAttachmentsReq, requests.MaintenanceAttachmentsReq{ID: item.ID, FileName: item.FileName, Path: filePath, FileType: fileType, Status: item.Status})
	// 	}
	// }

	data, err := h.Usecase.CreateMaintenanceHistory(idParent, userID, req)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(errResponse.Code, errResponse)
		return
	}
	id, _ := data.(*int)

	c.JSON(201, responses.SuccessResponse(responses.MaintenanceRoadHisID{ID: *id, IDParent: *id}, 201))
}

func validateCreateRoadHis(req requests.MaintenanceRoadHistoryReq, checkType int) error {
	errMsg := []string{}

	if *req.KmStart == 0 && *req.KmEnd == 0 {
		errMsg = append(errMsg, "km_start: zero value")
		errMsg = append(errMsg, "km_end: zero value")
	}

	switch req.MaintenanceType {
	case 1, 2: // Handles both case 1 and 2 as they both need to check LaneNo
		if req.LaneNo == nil || *req.LaneNo == 0 {
			errMsg = append(errMsg, "lane_no: zero value")
		}

		if req.MaintenanceType == 2 { // Additional check for GridNo in case 2 only
			if req.GridNo == nil || *req.GridNo == 0 {
				errMsg = append(errMsg, "grid_no: zero value")
			}
		}
	default:
		if req.LaneNo == nil || *req.LaneNo == 0 {
			errMsg = append(errMsg, "lane_no: zero value")
		}

		if req.GridNo == nil || *req.GridNo == 0 {
			errMsg = append(errMsg, "grid_no: zero value")
		}
	}

	if (checkType == 1) && *req.KmStart > *req.KmEnd {
		errMsg = append(errMsg, "km_start: incorrect")
		errMsg = append(errMsg, "km_end: incorrect")
	} else if (checkType == 2) && *req.KmStart < *req.KmEnd {
		errMsg = append(errMsg, "km_start: incorrect")
		errMsg = append(errMsg, "km_end: incorrect")
	}

	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, ","))
	}

	return nil
}

// @Summary ประวิติการซ่อมในช่วงประกัน
// @description
// @tags Maintenance History
// @id UpdateMaintenanceHistory
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your maintenance ID Parent"
// @param m_road_his_id path string true "Insert your maintenance road history ID"
// @Param maintenance body requests.MaintenanceRoadsReq true "JSON object that represents the maintenance project to be created"
// @Success 200 {object} responses.Success{data=responses.MaintenanceRoadHisID{}} "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/{id_parent}/road_history/{m_road_his_id} [put]
func (h *Handler) UpdateMaintenanceRoadHistory(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id_parent"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}

	historyID, err := strconv.Atoi(c.Param("m_road_his_id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}
	userID := helpers.GetUserInfo(c).UserId
	var req requests.MaintenanceRoadHistoryReq

	if err := c.ShouldBind(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	var checkType int

	roadInfo, err := h.Usecase.GetLastRoadInfoByID(*req.RoadID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}
	checkType = roadInfo.RefDirectionId

	var allErrors map[string]string = make(map[string]string)
	if validateErr := validator.Validate(req); validateErr != nil {
		errs := helpers.ConverstError(validateErr)
		for k, v := range errs {
			allErrors[k] = v
		}
	}

	validateErr := validateUpdateRoadHis(req, checkType)
	if validateErr != nil {
		errs := helpers.ConverstError(validateErr)
		for k, v := range errs {
			allErrors[k] = v
		}
	}

	checkValidateIsMethod, err := h.Usecase.CheckValidateIsMethod(ID)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if (checkValidateIsMethod && req.InterventionCriteriaID == nil) || (checkValidateIsMethod && *req.InterventionCriteriaID == 0) {
		allErrors["intervention_criteria_id"] = "โปรดระบุ"
	}

	if len(allErrors) > 0 {
		errResponse := responses.ValidateResponse(allErrors)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	// attachments := req.Attachments
	// var maintenanceAttachmentsReq []requests.MaintenanceAttachmentsReq
	// for i, item := range attachments {
	// 	switch item.Status {
	// 	case "delete":
	// 		att, err := h.Usecase.GetMaintenanceaAttachmentByID(*item.ID)
	// 		if err != nil {
	// 			errResponse := responses.FailRespone(err)
	// 			c.JSON(400, errResponse)
	// 			return
	// 		}
	// 		filePath := att.Path
	// 		if filePath != "" {
	// 			err = helpers.DeleteFile(filePath)
	// 			if err != nil {
	// 				c.JSON(400, err)
	// 				return
	// 			}
	// 		}
	// 		maintenanceAttachmentsReq = append(maintenanceAttachmentsReq, requests.MaintenanceAttachmentsReq{ID: item.ID, FileName: item.FileName, Path: "", FileType: "", Status: item.Status})
	// 	case "upload":
	// 		imgBase64 := item.File
	// 		base64 := ""
	// 		if strings.Contains(imgBase64, ",") {
	// 			dataImgBase64 := helpers.Explode(",", imgBase64)
	// 			base64 = dataImgBase64[1]
	// 		} else {
	// 			base64 = imgBase64
	// 		}
	// 		pathOutput := os.Getenv("MAINTENANCE_HIS_FILE_DIR")
	// 		fileName := "maintenance_his" + time.Now().Format("20060102150405") + "_" + strconv.Itoa(i)
	// 		fileNames := helpers.Explode(".", item.FileName)
	// 		fileNameType := ""
	// 		if len(fileNames) > 1 {
	// 			fileNameType = fileNames[len(fileNames)-1]
	// 		} else {
	// 			errResponse := responses.FailRespone(fmt.Errorf("ชื่อไฟล์ไม่ถูกต้อง"))
	// 			c.JSON(http.StatusBadRequest, errResponse)
	// 			return
	// 		}
	// 		filePath, fileType, err := helpers.DecodeFileBase64FileType(base64, pathOutput, fileName, fileNameType, "img")
	// 		if err != nil {
	// 			errResponse := responses.FailRespone(err)
	// 			c.JSON(http.StatusBadRequest, errResponse)
	// 			return
	// 		}
	// 		maintenanceAttachmentsReq = append(maintenanceAttachmentsReq, requests.MaintenanceAttachmentsReq{ID: item.ID, FileName: item.FileName, Path: filePath, FileType: fileType, Status: item.Status})
	// 	}
	// }

	id, idParent, err := h.Usecase.UpdateMaintenanceHistory(ID, historyID, userID, req)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(errResponse.Code, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(responses.MaintenanceRoadHisID{ID: *id, IDParent: *idParent}, 200))
}

func validateUpdateRoadHis(req requests.MaintenanceRoadHistoryReq, checkType int) error {
	errMsg := []string{}

	if *req.KmStart == 0 && *req.KmEnd == 0 {
		errMsg = append(errMsg, "km_start: zero value")
		errMsg = append(errMsg, "km_end: zero value")
	}
	switch req.MaintenanceType {
	case 1, 2: // Handles both case 1 and 2 as they both need to check LaneNo
		if req.LaneNo == nil || *req.LaneNo == 0 {
			errMsg = append(errMsg, "lane_no: zero value")
		}

		if req.MaintenanceType == 2 { // Additional check for GridNo in case 2 only
			if req.GridNo == nil || *req.GridNo == 0 {
				errMsg = append(errMsg, "grid_no: zero value")
			}
		}
	default:
		if req.LaneNo == nil || *req.LaneNo == 0 {
			errMsg = append(errMsg, "lane_no: zero value")
		}

		if req.GridNo == nil || *req.GridNo == 0 {
			errMsg = append(errMsg, "grid_no: zero value")
		}
	}

	if (checkType == 1) && *req.KmStart > *req.KmEnd {
		errMsg = append(errMsg, "km_start: incorrect")
		errMsg = append(errMsg, "km_end: incorrect")
	} else if (checkType == 2) && *req.KmStart < *req.KmEnd {
		errMsg = append(errMsg, "km_start: incorrect")
		errMsg = append(errMsg, "km_end: incorrect")
	}

	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, ","))
	}

	return nil
}

// @Summary ลบประวิติการซ่อมในช่วงประกัน
// @description
// @tags Maintenance History
// @id UpdateMaintenanceHistory
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your maintenance ID Parent"
// @param m_road_his_id path string true "Insert your history ID"
// @Success 204 {object} responses.Success{} "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/{id_parent}/road_history/{m_road_his_id} [delete]
func (h *Handler) DeleteMaintenanceRoadHistory(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id_parent"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}

	historyID, err := strconv.Atoi(c.Param("m_road_his_id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		errResponse.Code = http.StatusBadRequest
		c.JSON(errResponse.Code, errResponse)
		return
	}

	data, err := h.Usecase.DeleteMaintenanceHistory(ID, historyID)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(errResponse.Code, errResponse)
		return
	}
	c.JSON(204, responses.SuccessResponse(data, 204))
}

// @Summary GetMaintenancePlanReport
// @description
// @tags Maintenance Plan Progress
// @id GetMaintenancePlanReport
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your maintenance ID"
// @Param plan_id query string false "เรียงลำดับ id ปีงบประมาณ"
// @Success 200 {object}  responses.Success{data=string}  "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/{id}/plan_progress_export_report [get]
// func (h *Handler) GetMaintenancePlanReport(c *gin.Context) {
// 	planID, ok := c.Request.URL.Query()["plan_id"] // ปีงบประมาณ
// 	planIDs := []int{}
// 	if ok {
// 		strArr := strings.Split(planID[0], ",")
// 		for _, val := range strArr {
// 			planIDs = append(planIDs, helpers.StrToInt(val))
// 		}
// 	} else {
// 		planIDs = []int{}
// 	}

// 	maintenanceID, _ := strconv.Atoi(c.Param("id"))
// 	dataChart, err := h.Usecase.GetMaintenancePlanProgressReport(maintenanceID, planIDs)
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusUnprocessableEntity, errResponse)
// 		return
// 	}

// 	dataTable, err := h.Usecase.GetMaintenancePlanProgressTableReport(maintenanceID, planIDs)
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusUnprocessableEntity, errResponse)
// 		return
// 	}

// 	resp, err := h.Usecase.CreateMaintenancePlanReport(dataChart, dataTable, maintenanceID)
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusUnprocessableEntity, errResponse)
// 		return
// 	}

// 	c.JSON(200, responses.SuccessResponse(resp, 200))
// }

// @Summary GetMaintenancePlanHistoryReport
// @description
// @tags Maintenance Plan History Progress
// @id GetMaintenancePlanHistoryReport
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your maintenance ID"
// @Param plan_id query string false "เรียงลำดับ id ปีงบประมาณ"
// @Success 200 {object}  responses.Success{data=string}  "OK"
// @Failure 400 {object} responses.BadRequestErrorResponse "Bad Request - the request was invalid or cannot be served"
// @response 422 {object} responses.Validate "Unprocessable Entity - the request was well-formed but was unable to be followed due to semantic errors"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/{id}/plan_progress_history_export_report [get]
// func (h *Handler) GetMaintenancePlanHistoryReport(c *gin.Context) {
// 	planID, ok := c.Request.URL.Query()["plan_id"] // ปีงบประมาณ
// 	planIDs := []int{}
// 	if ok {
// 		strArr := strings.Split(planID[0], ",")
// 		for _, val := range strArr {
// 			planIDs = append(planIDs, helpers.StrToInt(val))
// 		}
// 	} else {
// 		planIDs = []int{}
// 	}

// 	maintenanceID, _ := strconv.Atoi(c.Param("id"))
// 	dataChart, err := h.Usecase.GetMaintenancePlanProgressReport(maintenanceID, planIDs)
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusUnprocessableEntity, errResponse)
// 		return
// 	}

// 	dataTable, err := h.Usecase.GetMaintenancePlanProgressTableReport(maintenanceID, planIDs)
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusUnprocessableEntity, errResponse)
// 		return
// 	}

// 	history, err := h.Usecase.GetMaintenanceHistoryByID(maintenanceID)
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusInternalServerError, errResponse)
// 		return
// 	}

// 	resp, err := h.Usecase.CreateMaintenanceHistoryPlanReport(dataChart, dataTable, history, maintenanceID)
// 	if err != nil {
// 		errResponse := responses.FailRespone(err)
// 		c.JSON(http.StatusUnprocessableEntity, errResponse)
// 		return
// 	}

// 	c.JSON(200, responses.SuccessResponse(resp, 200))
// }

// @Summary ปี
// @Description
// @tags Maintenance
// @id GetMaintenanceYear
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object}  responses.Success{data=string{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/maintenance/years [get]
func (h *Handler) GetMaintenanceYear(c *gin.Context) {
	var roadID int
	roadIDStr, ok := c.Request.URL.Query()["road_id"] // ปีงบประมาณ
	if ok {
		roadID = helpers.StrToInt(roadIDStr[0])
	}
	years, err := h.Usecase.GetMaintenanceYear(roadID)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if !ok {
			// Handle the case where err is not *responses.AppErr
			c.JSON(http.StatusInternalServerError, responses.NewAppErr(http.StatusInternalServerError, err.Error()))
			return
		}
		errResponse := responses.FailRespone(appErr)
		c.JSON(errResponse.Code, errResponse)
		return
	}
	// currentTime := time.Now()
	// maxTime := currentTime.AddDate(2, 0, 0).Year()
	// minTime := currentTime.AddDate(-13, 0, 0).Year()
	// data := []int{}
	// for i := minTime; i <= maxTime; i++ {
	// 	data = append(data, i)
	// }

	// for _, item := range years {
	// 	isVal := helpers.InArrayInt(item, data)
	// 	if !isVal {
	// 		data = append(data, item)
	// 	}
	// }
	sort.Slice(years, func(i, j int) bool {
		return years[i] > years[j]
	})
	c.JSON(http.StatusOK, responses.SuccessResponse(years, http.StatusOK))
}
