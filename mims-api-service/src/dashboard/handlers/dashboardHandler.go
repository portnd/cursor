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
	"gopkg.in/validator.v2"
)

type Handler struct {
	useCase domains.UseCase
}

func NewHandler(usecase domains.UseCase) *Handler {
	return &Handler{
		useCase: usecase,
	}
}

// @summary
// @description
// @tags Dashboard
// @id GetRoadDashboard
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RoadDashboard} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/road [get]
func (h *Handler) GetRoadDashboard(c *gin.Context) {

	data, err := h.useCase.GetRoadDashboard()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
}

// @summary
// @description
// @tags Dashboard
// @id GetAsset
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param page query string true "Insert your page number"
// @param limit query string true "Insert your limit number"
// @Param road_id query string false "สายทาง 1,2"
// @Param depot_code query string false "หน่วยงานที่รับผิดชอบ 26103,26104"
// @Param km_start query string false "กม. เริ่มต้น 0"
// @Param km_end query string false " กม. สิ้นสุด 1000"
// @Param year query string false "ปี 2020"
// @Success 200 {object} responses.Pagination{items=[]responses.DashboardAsset}  "OK"
// @failure 400 {object} responses.BadRequestErrorResponse
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/asset [get]
func (h *Handler) GetAsset(c *gin.Context) {
	var roadIds []int
	var strDepotCode []string
	var filter requests.Asset
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

	result, total, err := h.useCase.GetAsset(roadIds, strDepotCode, filter)
	if err != nil {
		c.JSON(400, err)
		return
	}

	responds := result.([]responses.DashboardAsset)
	if len(responds) == 0 {
		c.JSON(http.StatusOK, responses.SuccessResponse([]int{}, http.StatusOK))
		return
	}

	resp := helpers.Pagination(result, int64(filter.Limit), int64(filter.Page), total)

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags Dashboard
// @id GetAssetDetail
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param page query string true "Insert your page number"
// @param limit query string true "Insert your limit number"
// @Param road_id query string false "สายทาง 1,2"
// @Param depot_code query string false "หน่วยงานที่รับผิดชอบ 26103,26104"
// @Param km_start query string false "กม. เริ่มต้น 0"
// @Param km_end query string false " กม. สิ้นสุด 1000"
// @Param year query string false "ปี 2020"
// @success 200 {object} responses.Pagination{items=[]responses.AssetRespond} "OK"
// @failure 400 {object} responses.BadRequestErrorResponse
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/asset_detail [get]
func (h *Handler) GetAssetDetail(c *gin.Context) {
	var roadIds []int
	var strDepotCode []string
	reqRoadId := c.Query("road_id")

	var filter requests.Asset
	if err := c.ShouldBind(&filter); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	reqDepotCode := c.Query("depot_code")
	if reqDepotCode != "" {
		strDepotCode = strings.Split(reqDepotCode, ",")
	}

	strRoadId := strings.Split(reqRoadId, ",")
	if reqRoadId != "" {
		for _, id := range strRoadId {
			newId, err := helpers.ConvertStringToInt(id)
			if err != nil {
				errResponse := responses.FailRespone(err)
				c.JSON(errResponse.Code, errResponse)
				return
			}
			roadIds = append(roadIds, newId)
		}
	}
	result, total, err := h.useCase.GetAssetDetail(roadIds, strDepotCode, filter)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	resp := helpers.Pagination(result, int64(filter.Limit), int64(filter.Page), total)

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags Dashboard
// @id GetAssetMap
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param depot_code query string false "หน่วยงานที่รับผิดชอบ 26103,26104"
// @Param road_id query string false "สายทาง 1,2"
// @Param km_start query string false "กม. เริ่มต้น 0"
// @Param km_end query string false " กม. สิ้นสุด 1000"
// @Param year query string false "ปี 2020"
// @Param ref_asset_id query string false "รหัสกลุ่มสินทรัพย์ 1,2,3,4,5 "
// @Param left query string false "พิกัดซ้าย 100.47408466074495"
// @Param bottom query string false "พิกัดล่าง 13.449818134362182"
// @Param right query string false "พิกัดขวา 101.00609456630438"
// @Param top query string false "พิกัดบน ex 13.911723542477098"
// @Param zoom query string false "ซูม ex 12"
// @response 200 {object} interface{} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/asset_map [get]
func (h *Handler) GetAssetMap(c *gin.Context) {

	// depotCode := c.Query("depot_code")
	// roadID := c.Query("road_id")
	// kmStart := c.Query("km_start")
	// kmEnd := c.Query("km_end")
	// year := c.Query("year")
	// left := c.Query("left")
	// right := c.Query("right")
	// bottom := c.Query("bottom")
	// top := c.Query("top")
	// zoom := c.Query("zoom")

	var filter requests.AssetMap
	if err := c.ShouldBind(&filter); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	var refAssetIDs []string
	refAssetID := c.Query("ref_asset_id")
	if refAssetID != "" {
		refAssetIDs = strings.Split(refAssetID, ",")
	}

	var depotCodes []string
	depotCode := c.Query("depot_code")
	if depotCode != "" {
		depotCodes = strings.Split(depotCode, ",")
	}

	var roadIds []string
	roadId := c.Query("road_id")
	if roadId != "" {
		roadIds = strings.Split(roadId, ",")
	}

	data, err := h.useCase.GetAssetMap(roadIds, refAssetIDs, depotCodes, filter)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}
	if err := helpers.SortStructByID(data, false); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}
	responds := data.([]responses.AssetMap)
	if len(responds) == 0 {
		c.JSON(http.StatusOK, responses.SuccessResponse([]int{}, http.StatusOK))
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
}

// @summary
// @description
// @tags Dashboard
// @id GetAssetMapDetailByID
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @success 200 {object} responses.AssetRespond
// @failure 400 {object} responses.BadRequestErrorResponse
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/asset_map/:id/detail/:asset_table_id [get]
func (h *Handler) GetAssetMapDetailByID(c *gin.Context) {

	assetTableID, err := strconv.Atoi(c.Params.ByName("asset_table_id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	ID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	result, err := h.useCase.GetAssetMapDetailByID(ID, assetTableID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))
}

// @summary
// @description
// @tags Dashboard
// @id GetDashboardYears
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param type_name query string false "Ex (road | condition | surface | maintenance | "" )"
// @success 200 {object} responses.AssetRespond
// @failure 400 {object} responses.BadRequestErrorResponse
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/years [get]
func (h *Handler) GetDashboardYears(c *gin.Context) {
	typeName := c.Query("type_name")

	result, err := h.useCase.GetYears(typeName)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))
}

// @summary
// @description
// @tags Dashboard
// @id GetDashboardCondition
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param depot_code query string false "หน่วยงานที่รับผิดชอบ 26103,26104"
// @param road_id query string false "Comma-separated list of road IDs (e.g., 1,2,3)"
// @Param km_start query string false "กม. เริ่มต้น 0"
// @Param km_end query string false " กม. สิ้นสุด 1000"
// @param condition_type query string true "Condition type (e.g., IFI=1,IRI=2,RUT=3,MPD=4,แถบสะท้อนแสง=5)"
// @param condition_owner_id query string true "Owner ID "
// @param year query string false "Year or 'latest' for the latest year available"
// @success 200 {object} responses.ConditionDashboard "Successful operation"
// @failure 400 {object} responses.BadRequestErrorResponse "Bad request"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/condition [get]
func (h *Handler) GetDashboardCondition(c *gin.Context) {

	var filter requests.Condition
	if err := c.ShouldBind(&filter); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	var roadIds []string
	roadId := c.Query("road_id")
	if roadId != "" {
		roadIds = strings.Split(roadId, ",")
	}

	var depotCodes []string
	depotCode := c.Query("depot_code")
	if depotCode != "" {
		depotCodes = strings.Split(depotCode, ",")
	}

	result, err := h.useCase.GetDashboardCondition(roadIds, depotCodes, filter)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))
}

// @summary
// @description
// @tags Dashboard
// @id GetDashboardConditionMap
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param depot_code query string false "หน่วยงานที่รับผิดชอบ 26103,26104"
// @param road_id query string false "Comma-separated list of road IDs (e.g., 1,2,3)"
// @Param km_start query string false "กม. เริ่มต้น 0"
// @Param km_end query string false " กม. สิ้นสุด 1000"
// @param condition_type query string true "Condition type (e.g., IFI=1,IRI=2,RUT=3,MPD=4,แถบสะท้อนแสง=5)"
// @param condition_owner_id query string true "Owner ID "
// @param year query string false "Year or 'latest' for the latest year available"
// @param left query string true "Bounding Map Left Value"
// @param right query string true "Bounding Map Right Value"
// @param bottom query string true "Bounding Map Bottom Value"
// @param top query string true "Bounding Map Top Value"
// @param page query string false "Year or 'page"
// @param limit query string false "Limit"
// @success 200 {object} responses.Success{data=responses.Pagination{items=[]responses.DashboardConditionMap}}
// @failure 400 {object} responses.BadRequestErrorResponse
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/condition_map [get]
func (h *Handler) GetDashboardConditionMap(c *gin.Context) {

	var filter requests.ConditionMap
	if err := c.ShouldBind(&filter); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	var roadIDs []string
	roadId := c.Query("road_id")
	if roadId != "" {
		roadIDs = strings.Split(roadId, ",")
	}

	var depotCodes []string
	depotCode := c.Query("depot_code")
	if depotCode != "" {
		depotCodes = strings.Split(depotCode, ",")
	}

	page := c.Query("page")
	limit := c.Query("limit")

	result, err := h.useCase.GetDashboardConditionMap(roadIDs, depotCodes, page, limit, filter)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))
}

// @summary
// @description
// @tags Dashboard
// @id GetDashboardRoadConditionList
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @response 200 {object} responses.Success{data=responses.RoadConditionList}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/condition_list [get]
func (h *Handler) GetRoadConditionList(c *gin.Context) {

	roadID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	roadConditionLists, err := h.useCase.GetRoadConditionList(roadID)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(roadConditionLists, 200))

}

// @summary
// @description
// @tags Dashboard
// @id GetDashboardRoadConditionDetails
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id_parent path string true "Insert your id_parent"
// @param condition_rang_type query string true "Insert your condition_rang_type"
// @response 200 {object} responses.Success{data=[]responses.RoadDamageImport}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/condition_details/{id_parent} [get]
func (h *Handler) GetRoadConditionDetails(c *gin.Context) {

	IDParent, err := strconv.Atoi(c.Params.ByName("id_parent"))
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}
	var req requests.RoadConditionDetails
	if err := c.ShouldBindQuery(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	if validateErr := validator.Validate(req); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	conditionLists, err := h.useCase.GetRoadConditionDetails(req.ConditionRangeType, IDParent)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(conditionLists, 200))

}

// @summary
// @description
// @tags Dashboard
// @id GetDashboardRoadRetroReflectivityList
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @response 200 {object} responses.Success{data=[]responses.RoadRetroReflectivityList}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/retro_reflectivity_list [get]
func (h *Handler) GetRoadRetroReflectivityList(c *gin.Context) {

	roadID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	roadConditionLists, err := h.useCase.GetRoadRetroReflectivityList(roadID)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(roadConditionLists, 200))

}

// @summary
// @description
// @tags Dashboard
// @id GetDashboardRoadRetroReflectivityDetails
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id_parent path string true "Insert your id_parent"
// @param range_type query string true "Insert your ref_reflectivity_range"
// @response 200 {object} responses.Success{data=[]responses.RetroReflectivityDetails}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/dashboard/retro_reflectivity_details/{id_parent} [get]
func (h *Handler) GetRoadRetroReflectivityDetails(c *gin.Context) {

	IDParent, err := strconv.Atoi(c.Params.ByName("id_parent"))
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}
	var req requests.RoadRetroReflectivityDetails
	if err := c.ShouldBindQuery(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	if validateErr := validator.Validate(req); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	RetroReflectivityDetails, err := h.useCase.GetRoadRetroReflectivityDetails(req.RangeType, req.RefStripeTypeIDs, IDParent)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(RetroReflectivityDetails, 200))

}
