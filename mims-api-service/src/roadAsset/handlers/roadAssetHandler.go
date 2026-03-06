package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadAsset/domains"
	"gopkg.in/validator.v2"

	"github.com/gin-gonic/gin"
)

// init handler
type RoadAssetHandler struct {
	roadAssetUseCase domains.RoadAssetUseCase
}

// init handler
func NewRoadAssetHandler(usecase domains.RoadAssetUseCase) *RoadAssetHandler {
	return &RoadAssetHandler{
		roadAssetUseCase: usecase,
	}
}

// ================================== start function  ==================================
// @summary
// @description
// @tags Roads Asset
// @id GetRoadAssetDetail
// @param id path string true "Insert your road ID"
// @param road_asset_id path string true "Insert your road asset ID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param ref_asset_table_id query string true "ref_asset_table_id" example(2)
// @param page query string false "Insert your page"
// @param limit query string false "Insert your limit"
// @response 200 {object} responses.Success{data=responses.Pagination{items=[]responses.RoadAssetData}}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/asset_details/{road_asset_id} [get]
func (t *RoadAssetHandler) GetRoadAssetDetail(c *gin.Context) {
	// roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	roadAssetId, _ := strconv.Atoi(c.Params.ByName("road_asset_id"))
	var queryParams requests.AssetDetailsQueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	limit := queryParams.Limit
	page := queryParams.Page
	permissions := helpers.GetAccessControl(c)
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))

	data, total, err := t.roadAssetUseCase.GetRoadAssetDetail(queryParams, permissions, roadAssetId, userID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	pagination := helpers.Pagination(data, int64(limit), int64(page), int64(total))
	c.JSON(200, responses.SuccessResponse(pagination, 200))
}

// @summary
// @description
// @tags Roads Asset
// @id GetRoadAssetPermission
// @Param ref_asset_table_id query string false "ref_asset_table_id" example(2)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.Permission}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/asset_permission [get]
func (t *RoadAssetHandler) GetRoadAssetPermission(c *gin.Context) {
	var queryParams requests.AssetPermissionQueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	permissions := helpers.GetAccessControl(c)
	data, err := t.roadAssetUseCase.GetRoadAssetPermission(queryParams, permissions)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags Roads Asset
// @id GetRoadAssetRevisions
// @param id path string true "Insert your road ID"
// @Param ref_asset_table_id query string false "ref_asset_table_id" example(2)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} []responses.Success{data=responses.RoadAssetRevision}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/asset_revision_list [get]
func (t *RoadAssetHandler) GetRoadAssetRevisions(c *gin.Context) {
	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	var queryParams requests.AssetRevisionsQueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	permissions := helpers.GetAccessControl(c)
	data, err := t.roadAssetUseCase.GetRoadAssetRevisions(queryParams, permissions, roadID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags Roads Asset
// @id GetRoadAssetTemplate
// @param id path string true "Insert your road ID"
// @Param ref_asset_table_id query string false "ref_asset_table_id" example(2)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} []responses.Success{data=responses.RoadAssetTemplateColumn}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/asset_edit_template [get]
func (t *RoadAssetHandler) GetRoadAssetTemplate(c *gin.Context) {
	var queryParams requests.AssetTemplateQueryParams
	asset_object_id, ok := c.Request.URL.Query()["asset_object_id"] // ปีงบประมาณ
	if ok {
		assetObjectID, _ := strconv.Atoi(asset_object_id[0])
		queryParams.AssetObjectID = assetObjectID
	}

	action, ok := c.Request.URL.Query()["action"] // ปีงบประมาณ
	if ok {
		queryParams.Action = action[0]
	}

	ref_asset_table_id, ok := c.Request.URL.Query()["ref_asset_table_id"] // ปีงบประมาณ
	if ok {
		refAssetTableID, _ := strconv.Atoi(ref_asset_table_id[0])
		queryParams.RefAssetTableID = refAssetTableID
	}

	permissions := helpers.GetAccessControl(c)
	data, err := t.roadAssetUseCase.GetRoadAssetTemplate(queryParams, permissions)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags Roads Asset
// @id CreateRoadAsset
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @param Asset body map[string]interface{} true "Update your data"
// @response 201 {object} []responses.Success{data=responses.RoadDamageSetData} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/asset [post]
func (t *RoadAssetHandler) CreateRoadAsset(c *gin.Context) {
	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, err)
		return
	}

	var reqs map[string]interface{}
	err = json.Unmarshal(body, &reqs)
	if err != nil {
		c.JSON(400, err)
		return
	}
	data, err := t.roadAssetUseCase.CreateRoadAsset(reqs, roadID, 0, userID)
	if err != nil {
		helpers.PrintlnJson(err)
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(data, http.StatusCreated))
}

// @summary
// @description
// @tags Roads Asset
// @id UpdateRoadAsset
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @param id_parent_asset path string true "Insert your road id_parent_asset"
// @param Asset body map[string]interface{} true "Update your data"
// @response 200 {object} []responses.Success{data=responses.RoadDamageSetData} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/asset/{id_parent_asset} [PUT]
func (t *RoadAssetHandler) UpdateRoadAsset(c *gin.Context) {
	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	IDParentAsset, _ := strconv.Atoi(c.Params.ByName("id_parent_asset"))
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	var reqs map[string]interface{}
	err = json.Unmarshal(body, &reqs)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	roadInfo, err := t.roadAssetUseCase.GetRoadInfoByRoadID(roadID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	if reqs["km"] != nil { // asset in
		validateMsgArr := map[string]string{}
		switch roadInfo.RefDirectionId {
		case 1:
			if reqs["km"].(float64) < float64(roadInfo.KmStart) || reqs["km"].(float64) > float64(roadInfo.KmEnd) {
				validateMsgArr["km"] = constants.INVALID_GEOM_RANGE
				errResponse := responses.ValidateResponse(validateMsgArr)
				c.JSON(http.StatusUnprocessableEntity, errResponse)
				return
			}
		case 2:
			if reqs["km"].(float64) > float64(roadInfo.KmStart) || reqs["km"].(float64) < float64(roadInfo.KmEnd) {
				validateMsgArr["km"] = constants.INVALID_GEOM_RANGE
				errResponse := responses.ValidateResponse(validateMsgArr)
				c.JSON(http.StatusUnprocessableEntity, errResponse)
				return
			}
		}
	}

	data, err := t.roadAssetUseCase.CreateRoadAsset(reqs, roadID, IDParentAsset, userID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(http.StatusAccepted, responses.SuccessResponse(data, http.StatusAccepted))
}

// @summary
// @description
// @tags Roads Asset
// @id ConfirmRoadAsset
// @param id_parent path string true "Insert your road asset id parent"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.AssetID}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/asset_confirm/{id_parent} [put]
func (t *RoadAssetHandler) ConfirmRoadAsset(c *gin.Context) {
	idParent, _ := strconv.Atoi(c.Params.ByName("id_parent"))
	err := c.ShouldBind(&idParent)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))
	resp, err := t.roadAssetUseCase.ConfirmRoadAsset(idParent, uid)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(responses.AssetID{ID: resp}, 200))
}

// @summary
// @description
// @tags Roads Asset
// @id CancelRoadAsset
// @param id_parent path string true "Insert your road asset id parent"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 204 "No Content"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/asset_cancel/{id_parent} [put]
func (t *RoadAssetHandler) CancelRoadAsset(c *gin.Context) {
	idParent, _ := strconv.Atoi(c.Params.ByName("id_parent"))
	err := c.ShouldBind(&idParent)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))
	resp, err := t.roadAssetUseCase.CancelRoadAsset(idParent, uid)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	c.JSON(http.StatusNoContent, responses.SuccessResponse(resp, 204))
}

// @summary
// @description
// @tags Roads Asset
// @id DeleteRoadAsset
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id_parent path string true "Insert your asset id_parent"
// @response 204 "No Content"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/asset_delete/{id_parent}  [delete]
func (t *RoadAssetHandler) DeleteRoadAsset(c *gin.Context) {

	idParent, _ := strconv.Atoi(c.Params.ByName("id_parent"))
	// err := c.ShouldBind(&idParent)
	// if err != nil {
	// 	errResponse := responses.FailRespone(err)
	// 	c.JSON(400, errResponse)
	// 	return
	// }

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))

	_, err := t.roadAssetUseCase.DeleteRoadAsset(c, idParent, uid)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	c.JSON(http.StatusNoContent, responses.SuccessResponse(nil, 204))
}

// @id GetRoadAssetKm
// @param id path string true "Insert your road ID"
// @Param geom query string true "Insert geom point or lineString" example(LINESTRING 100.550018 13.812126339999997,100.55003801492836 13.81226394772273,100.55005475207957 13.81242085851546,100.55007288399342 13.812572887639076,100.55009241066983 13.812733285338302,100.55010635829585 13.812869274692002,100.55011622807295 13.812980309684542 OR POINT 99.76503923535347+13.990306811858371) ))
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} []responses.Success{data=responses.RoadKmByGeom}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Fail "validate error"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/km [get]
func (t *RoadAssetHandler) GetRoadKmByGeom(c *gin.Context) {
	var request requests.RoadAssetKm
	if err := c.ShouldBind(&request); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	if validateErr := validator.Validate(request); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	data, err := t.roadAssetUseCase.GetRoadKmByGeom(request.Geom, roadID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	if data.LaneNo == 0 && data.KmStart == 0 && data.KmEnd == 0 && data.Km == 0 {
		c.JSON(http.StatusOK, responses.SuccessResponse(responses.NoData{}, http.StatusOK))
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
}

// @summary
// @description
// @tags Roads Asset
// @id DeleteRoadAssetObject
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your asset ID"
// @param ref_asset_table_id path string true "Insert your ref_asset_table_id"
// @Param asset_object_id path string true "Insert  your asset_object_id"
// @response 204 "No Content"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/asset/{id}/table/{ref_asset_table_id}/asset_object/{asset_object_id} [delete]
func (t *RoadAssetHandler) DeleteRoadAssetObject(c *gin.Context) {

	assetObjectID, err := strconv.Atoi(c.Params.ByName("asset_object_id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	// return
	// idParent, err := strconv.Atoi(c.Params.ByName("ref_asset_table_id"))
	// if err != nil {
	// 	errResponse := responses.FailRespone(err)
	// 	c.JSON(400, errResponse)
	// 	return
	// }

	assetID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	refAssetTableID, err := strconv.Atoi(c.Params.ByName("ref_asset_table_id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))
	_, err = t.roadAssetUseCase.DeleteRoadAssetObject(assetID, refAssetTableID, assetObjectID, uid)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(http.StatusNoContent, responses.SuccessResponse(nil, 204))
}

// @summary
// @description
// @tags Roads Asset
// @id GetAssetRoadType
// @param id path string true "Insert your asset table ID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.AssetTableType}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/asset_table/{id} [get]
func (t *RoadAssetHandler) GetAssetRoadType(c *gin.Context) {
	assetTableId, _ := strconv.Atoi(c.Params.ByName("id"))
	data, err := t.roadAssetUseCase.GetAssetTableByID(assetTableId)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}
