package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	helpers "gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	requests "gitlab.com/mims-api-service/requests"
	responses "gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/road/domains"
	"gopkg.in/validator.v2"

	"github.com/gin-gonic/gin"
)

// init handler
type RoadHandler struct {
	roadUseCase domains.RoadUseCase
}

// init handler
func NewRoadHandler(usecase domains.RoadUseCase) *RoadHandler {
	return &RoadHandler{
		roadUseCase: usecase,
	}
}

// ================================== start function  ==================================
// @summary
// @description
// @tags Roads
// @id Roads
// @Param keyword query string false "keyword" example(ฉิมพลี)
// @Param road_id query string false "road_id" example(1,2)
// @Param road_group_id query string false "road_group_id" example(1,2)
// @Param road_section_id query string false "road_section_id" example(1,2,3)
// @Param km_start query string false "km_start" example(0)
// @Param km_end query string false "km_end" example(8000)
// @Param depot_code query string false "depot_code" example(26104,26105)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=[]responses.RoadList} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads [get]
func (t *RoadHandler) GetRoadList(c *gin.Context) {

	var params requests.RoadPrams

	keyword, ok := c.Request.URL.Query()["keyword"]
	if ok {
		params.Keyword = keyword[0]
	}

	roadID, ok := c.Request.URL.Query()["road_id"]
	if ok {
		idSplit := strings.Split(roadID[0], ",")
		arrayInt, err := helpers.ConvertToArrayInt(idSplit)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		params.RoadId = arrayInt
	}

	roadGroupId, ok := c.Request.URL.Query()["road_group_id"]
	if ok {
		idSplit := strings.Split(roadGroupId[0], ",")
		arrayInt, err := helpers.ConvertToArrayInt(idSplit)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		params.RoadGroupId = arrayInt
	}
	roadSectionId, ok := c.Request.URL.Query()["road_section_id"]
	if ok {
		idSplit := strings.Split(roadSectionId[0], ",")
		arrayInt, err := helpers.ConvertToArrayInt(idSplit)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		params.RoadSectionId = arrayInt
	}

	refSurfaceId, ok := c.Request.URL.Query()["ref_surface_id"]
	if ok {
		idSplit := strings.Split(refSurfaceId[0], ",")
		arrayInt, err := helpers.ConvertToArrayInt(idSplit)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		params.RefSurfaceId = arrayInt
	}

	depotCode, ok := c.Request.URL.Query()["depot_code"]
	if ok {
		params.DepotCode = strings.Split(depotCode[0], ",")
	}

	kmStart, ok := c.Request.URL.Query()["km_start"]
	if ok {
		dataFloat, err := strconv.ParseFloat(kmStart[0], 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		value := float32(dataFloat)
		params.KmStart = &value
	}

	kmEnd, ok := c.Request.URL.Query()["km_end"]
	if ok {
		dataFloat, err := strconv.ParseFloat(kmEnd[0], 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		value := float32(dataFloat)
		params.KmEnd = &value
	}

	isIri1000, ok := c.Request.URL.Query()["is_iri_1000"]
	if ok {
		if isIri1000[0] != "" {
			dataBool, err := strconv.ParseBool(isIri1000[0])
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}
			value := dataBool
			params.IsIri1000 = &value
		}
	}

	isIri100, ok := c.Request.URL.Query()["is_iri_100"]
	if ok {
		if isIri100[0] != "" {
			dataBool, err := strconv.ParseBool(isIri100[0])
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}
			value := dataBool
			params.IsIri100 = &value
		}
	}

	isRut100, ok := c.Request.URL.Query()["is_rut_100"]
	if ok {
		if isRut100[0] != "" {
			dataBool, err := strconv.ParseBool(isRut100[0])
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}
			value := dataBool
			params.IsRut100 = &value
		}
	}

	isIfi100, ok := c.Request.URL.Query()["is_ifi_100"]
	if ok {
		if isIfi100[0] != "" {
			dataBool, err := strconv.ParseBool(isIfi100[0])
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}
			value := dataBool
			params.IsIfi100 = &value
		}
	}

	isG7100, ok := c.Request.URL.Query()["is_g7_100"]
	if ok {
		if isG7100[0] != "" {
			dataBool, err := strconv.ParseBool(isG7100[0])
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}
			value := dataBool
			params.IsG7100 = &value
		}
	}
	// c.JSON(200, params)
	// return
	// if ok {
	// 	dataFloat, err := strconv.ParseFloat(kmEnd[0], 32)
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, err)
	// 		return
	// 	}
	// 	value := float32(dataFloat)
	// 	params.KmEnd = &value
	// }
	userID := helpers.GetUserID(c)
	result, err := t.roadUseCase.GetRoadGroupList(userID, params)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))
}

// @summary
// @description
// @tags Roads
// @id GetRoadByID
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @response 200 {object} []responses.Success{data=responses.RoadById}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id} [get]
func (t *RoadHandler) GetRoadByID(c *gin.Context) {
	userID := helpers.GetUserID(c)
	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	result, err := t.roadUseCase.GetRoadByID(roadID, userID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))
}

// @summary
// @description
// @tags Roads
// @id GetRoadDetailMenu
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param assetType query string false "assetin or assetout" example(assetin or assetout)
// @response 200 {object} []responses.RoadGroup "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/menu [get]
func (t *RoadHandler) GetRoadDetailMenu(c *gin.Context) {
	userId, _ := c.Get("userID")
	uid := uint(userId.(float64))
	// data, err := t.roadUseCase.GetMenu(uid)
	// if err != nil {
	// 	appErr, _ := err.(*responses.AppErr)
	// 	errResponse := responses.FailRespone(appErr)
	// 	c.JSON(appErr.StatusCode, errResponse)
	// 	return
	// }

	assetType, Ok := c.Request.URL.Query()["assetType"]
	if !Ok {
		errResponse := responses.FailRespone(errors.New("assetType is required"))
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	// var accessKeys []string
	// for _, item := range data {
	// 	accessKeys = append(accessKeys, item.AccessKey)
	// }
	accessKeys := helpers.GetAccessControl(c)

	// permissions := []string{"road_in_asset_access", "view_all_road_in_asset", "view_department_road_in_asset"}
	// if helpers.HasPermission(permissions, accessKeys) {
	data, err := t.roadUseCase.GetRoadDetailMenu(uid, accessKeys, assetType[0])
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
	// } else {
	// 	errResponse := responses.FailRespone(errors.New("access denied"))
	// 	c.JSON(http.StatusBadRequest, errResponse)
	// 	return
	// }
}

// @summary
// @description
// @tags Roads
// @id GetRoadDirectionLaneList
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @response 200 {object} []responses.RoadLaneList "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/lane_list [get]
func (t RoadHandler) GetRoadDirectionLaneList(c *gin.Context) {
	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	data, _ := t.roadUseCase.GetRoadDirectionLaneList(roadID)
	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
}

// func (t *RoadHandler) GetRoadStatusCout(c *gin.Context) {
// 	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
// 	permissions := helpers.GetAccessControl(c)
// 	t.roadUseCase.GetRoadDetailStatus(roadID, permissions)
// 	/// test

// }

// @summary road tree
// @description
// @tags Roads
// @id GetRoadTree
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object}  responses.Success{data=[]responses.RoadTree}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/tree [get]
func (t *RoadHandler) GetRoadTree(c *gin.Context) {
	data, err := t.roadUseCase.GetRoadTree()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags Roads
// @id GetRoadInit
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road section ID or road Id"
// @param level path string true "Insert your road level"
// @Param ref_road_type_id query string false "ref_road_type_id" example(1)
// @response 200 {object} []responses.Success{data=responses.RoadInit}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/init/{id}/{level} [get]
func (t *RoadHandler) GetRoadInit(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusNotFound, errResponse)
		return
	}

	level, err := strconv.Atoi(c.Params.ByName("level"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusNotFound, errResponse)
		return
	}

	refRoadTypeIdStr, ok := c.Request.URL.Query()["ref_road_type_id"]
	var refRoadTypeId int
	if ok {
		refRoadTypeId, err = strconv.Atoi(refRoadTypeIdStr[0])
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(http.StatusNotFound, errResponse)
			return
		}
	}

	result, err := t.roadUseCase.GetRoadInit(id, level, refRoadTypeId)

	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusNotFound, errResponse)
		return
	}

	c.JSON(200, responses.SuccessResponse(result, 200))
}

// @summary
// @description
// @tags Roads
// @id CreateRoad
// @Security Bearer
// @Param name formData string false "Name"
// @Param road_section_id formData int false "Road Section ID"
// @Param road_id formData int false "Road ID"
// @Param road_level formData int true "Road level"
// @Param ramp_id formData int false "ramp id"
// @Param km_start formData float32 false "KM Start"
// @Param km_end formData float32 false "KM End"
// @Param road_color_code formData string true "Road Color Code"
// @Param ref_road_type_id formData int true "Reference Road Type ID"
// @Param register_date formData string false "Register date"
// @Param center_line_shape_file formData file true "Center Line Shape File"
// @Param center_lane_shape_file formData file true "Center Lane Shape File"
// @Param remark formData string false "remark"
// @response 201 {object} responses.Success{data=responses.DataId}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads [post]
func (t *RoadHandler) CreateRoad(c *gin.Context) {

	var req requests.Road
	_ = c.ShouldBind(&req)
	var checkType int
	if req.RoadLevel == 2 {
		roadInfo, err := t.roadUseCase.GetLastRoadInfoByID(req.RoadID)
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(http.StatusNotFound, errResponse)
			return
		}
		checkType = roadInfo.RefDirectionId
	} else {
		checkType = req.RefRoadTypeID
	}

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

	if len(allErrors) > 0 {
		errResponse := responses.ValidateResponse(allErrors)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	userID := helpers.GetUserID(c)
	id, err := t.roadUseCase.CreateRoad(c, userID, req)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(responses.DataId{Id: id.(int)}, http.StatusCreated))
}

func validateCreate(req requests.Road, checkType int) error {
	errMsg := []string{}

	if req.KmStart == 0 && req.KmEnd == 0 {
		errMsg = append(errMsg, "km_start: zero value")
		errMsg = append(errMsg, "km_end: zero value")
	}

	if req.RoadLevel == 1 && req.RoadSectionID == 0 {
		errMsg = append(errMsg, "road_section_id: zero value")
	} else if req.RoadLevel == 2 {
		if req.RoadID == 0 {
			errMsg = append(errMsg, "road_id: zero value")
		}

		if req.Name == "" {
			errMsg = append(errMsg, "name: zero value")
		}

		if req.RampId == "" {
			errMsg = append(errMsg, "ramp_id: zero value")
		}
	}

	if (checkType == 1 || checkType == 3) && req.KmStart > req.KmEnd {
		errMsg = append(errMsg, "km_start: incorrect")
		errMsg = append(errMsg, "km_end: incorrect")
	} else if (checkType == 2 || checkType == 4) && req.KmStart < req.KmEnd {
		errMsg = append(errMsg, "km_start: incorrect")
		errMsg = append(errMsg, "km_end: incorrect")
	}

	if req.RegisterDate != "" {
		if !helpers.IsValidDate(req.RegisterDate) {
			errMsg = append(errMsg, "register_date: invalid date format")
		}
	}

	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, ","))
	}

	return nil
}

// @summary
// @description
// @tags Roads
// @id UpdateRoad
// @Security Bearer
// @param id path string true "Insert your road ID"
// @Param name formData string false "Name"
// @Param ramp_id formData int false "ramp id"
// @Param destination formData string false "Destination"
// @Param km_start formData float32 false "KM Start"
// @Param km_end formData float32 false "KM End"
// @Param road_color_code formData string false "Road Color Code"
// @Param ref_road_type_id formData int false "Reference Road Type ID"
// @Param register_date formData string false "Register Date"
// @Param center_line_shape_file formData file false "Center Line Shape File"
// @Param center_line_shape_file_status formData string false "Center Line Shape File Status"
// @Param center_lane_shape_file formData file true "Center Lane Shape File"
// @Param center_lane_shape_file_status formData string true "Center Lane Shape File Status"
// @Param remark formData string false "remark"
// @response 202 {object} responses.Success{data=responses.NoData{}}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id} [put]
func (t *RoadHandler) UpdateRoad(c *gin.Context) {
	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	var road models.Road
	err := t.roadUseCase.GetDataById(&road, roadID)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	roadInfo, err := t.roadUseCase.GetLastRoadInfoByID(roadID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusNotFound, errResponse)
		return
	}

	if road.IsInit {
		var req requests.RoadUpdateInit
		_ = c.ShouldBind(&req)
		var allErrors map[string]string = make(map[string]string)
		if validateErr := validator.Validate(req); validateErr != nil {
			errs := helpers.ConverstError(validateErr)
			for k, v := range errs {
				allErrors[k] = v
			}
		}

		var checkType int
		if road.RoadLevel == 2 {
			checkType = roadInfo.RefDirectionId
		} else {
			checkType = req.RefRoadTypeID
		}

		validateErr := validateUpdateInit(req, checkType, road.RoadLevel)
		if validateErr != nil {
			errs := helpers.ConverstError(validateErr)
			for k, v := range errs {
				allErrors[k] = v
			}
		}

		if len(allErrors) > 0 {
			errResponse := responses.ValidateResponse(allErrors)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}

		userID := helpers.GetUserID(c)
		_, err := t.roadUseCase.UpdateRoadInit(c, roadID, userID, req)
		if err != nil {
			appErr, _ := err.(*responses.AppErr)
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
			return
		}
		c.JSON(http.StatusAccepted, responses.SuccessResponse(responses.NoData{}, http.StatusAccepted))
	} else {
		var req requests.RoadUpdate

		_ = c.ShouldBind(&req)
		var allErrors map[string]string = make(map[string]string)
		if validateErr := validator.Validate(req); validateErr != nil {
			errs := helpers.ConverstError(validateErr)
			for k, v := range errs {
				allErrors[k] = v
			}
		}

		var checkType int
		if road.RoadLevel == 2 {
			checkType = roadInfo.RefDirectionId
		} else {
			checkType = req.RefRoadTypeID
		}

		validateErr := validateUpdate(req, checkType, road.RoadLevel)
		if validateErr != nil {
			errs := helpers.ConverstError(validateErr)
			for k, v := range errs {
				allErrors[k] = v
			}
		}

		if len(allErrors) > 0 {
			errResponse := responses.ValidateResponse(allErrors)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}

		userID := helpers.GetUserID(c)
		_, err = t.roadUseCase.UpdateRoad(c, roadID, userID, req)
		if err != nil {
			appErr, _ := err.(*responses.AppErr)
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
			return
		}
		c.JSON(http.StatusAccepted, responses.SuccessResponse(responses.NoData{}, http.StatusAccepted))
	}

	// no_file
	//    csv [validate], zip[path รูปต้องเป็น path เดิม]
	// delete
	//     csv [validate], zip[ลบpathรูป, ข้อมูล รูปหน้ากล้อง road_condition  ต้องเป็นค่าว่างทั้งหมด]
	// not_edit
	//    csv[ไม่มีการแก้ไขข้อมูล],zip[ไม่มีการแก้ไขข้อมูล]
	// upload
	//    csv[update ข้อมูล road_condition],zip[อัพเดท path img และ ข้อมูลรูปหน้ากล้อง]

}

func validateUpdateInit(req requests.RoadUpdateInit, checkType, roadLevel int) error {
	errMsg := []string{}

	if req.KmStart == 0 && req.KmEnd == 0 {
		errMsg = append(errMsg, "km_start: zero value")
		errMsg = append(errMsg, "km_end: zero value")
	}

	if roadLevel == 2 {
		if req.Name == "" {
			errMsg = append(errMsg, "name: zero value")
		}

		if req.RampId == "" {
			errMsg = append(errMsg, "ramp_id: zero value")
		}
	}

	if (checkType == 1 || checkType == 3) && req.KmStart > req.KmEnd {
		errMsg = append(errMsg, "km_start: incorrect")
		errMsg = append(errMsg, "km_end: incorrect")
	} else if (checkType == 2 || checkType == 4) && req.KmStart < req.KmEnd {
		errMsg = append(errMsg, "km_start: incorrect")
		errMsg = append(errMsg, "km_end: incorrect")
	}

	if req.RegisterDate != "" {
		if !helpers.IsValidDate(req.RegisterDate) {
			errMsg = append(errMsg, "register_date: invalid date format")
		}
	}

	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, ","))
	}

	return nil
}

func validateUpdate(req requests.RoadUpdate, checkType, roadLevel int) error {
	errMsg := []string{}

	if req.KmStart == 0 && req.KmEnd == 0 {
		errMsg = append(errMsg, "km_start: zero value")
		errMsg = append(errMsg, "km_end: zero value")
	}

	if req.RegisterDate != "" {
		if !helpers.IsValidDate(req.RegisterDate) {
			errMsg = append(errMsg, "register_date: invalid date format")
		}
	}

	if roadLevel == 2 {
		if req.Name == "" {
			errMsg = append(errMsg, "name: zero value")
		}

		if req.RampId == "" {
			errMsg = append(errMsg, "ramp_id: zero value")
		}
	}

	if (checkType == 1 || checkType == 3) && req.KmStart > req.KmEnd {
		errMsg = append(errMsg, "km_start: incorrect")
		errMsg = append(errMsg, "km_end: incorrect")
	} else if (checkType == 2 || checkType == 4) && req.KmStart < req.KmEnd {
		errMsg = append(errMsg, "km_start: incorrect")
		errMsg = append(errMsg, "km_end: incorrect")
	}

	if req.RegisterDate != "" {
		if !helpers.IsValidDate(req.RegisterDate) {
			errMsg = append(errMsg, "register_date: invalid date format")
		}
	}

	switch req.CenterLaneShapeFileStatus {
	case "no_file":
		errMsg = append(errMsg, "center_lane_shape_file: zero value")
	case "delete":
		errMsg = append(errMsg, "center_lane_shape_file: zero value")
	case "":
		errMsg = append(errMsg, "center_lane_shape_file: zero value")
	}

	switch req.CenterLineShapeFileStatus {
	case "no_file":
		errMsg = append(errMsg, "center_line_shape_file: zero value")
	case "delete":
		errMsg = append(errMsg, "center_line_shape_file: zero value")
	case "":
		errMsg = append(errMsg, "center_line_shape_file: zero value")
	}

	if (req.CenterLaneShapeFile == nil || req.CenterLaneShapeFile.Filename == "") && (req.CenterLaneShapeFileStatus != "not_edit") {
		errMsg = append(errMsg, "center_lane_shape_file: zero value")
	}

	if (req.CenterLineShapeFile == nil || req.CenterLaneShapeFile.Filename == "") && (req.CenterLineShapeFileStatus != "not_edit") {
		errMsg = append(errMsg, "center_line_shape_file: zero value")
	}

	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, ","))
	}

	return nil
}

// @summary
// @description
// @tags Roads
// @id DeleteRoad
// @Security Bearer
// @param id path string true "Insert your road ID"
// @response 204 {object} []responses.Success{data=responses.NoData{}}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id} [delete]
func (t *RoadHandler) DeleteRoad(c *gin.Context) {
	userID := helpers.GetUserID(c)
	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	_, err := t.roadUseCase.DeleteRoad(roadID, userID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(http.StatusNoContent, responses.SuccessResponse(responses.NoData{}, 204))

}

// @summary road tree
// @description
// @tags Roads
// @id GetRoadLanes
// @param id path string true "Insert your road ID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object}  responses.Success{data=[]responses.RoadLanes}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/lanes [get]
func (t *RoadHandler) GetRoadLanes(c *gin.Context) {
	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	data, err := t.roadUseCase.GetRoadLanes(roadID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}
