package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	helpers "gitlab.com/mims-api-service/helpers"
	requests "gitlab.com/mims-api-service/requests"
	responses "gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadCondition/domains"
	"gopkg.in/validator.v2"

	"github.com/gin-gonic/gin"
)

// init handler
type RoadConditionHandler struct {
	rcUseCase domains.RoadConditionUseCase
}

// init handler
func NewRoadConditionHandler(usecase domains.RoadConditionUseCase) *RoadConditionHandler {
	return &RoadConditionHandler{
		rcUseCase: usecase,
	}
}

// @summary
// @description
// @tags RoadConditions
// @id GetRoadConditionList
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @response 200 {object} responses.Success{data=responses.RoadConditionList}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/condition_list [get]
func (t *RoadConditionHandler) GetRoadConditionList(c *gin.Context) {

	userID, _ := c.Get("userID")
	roadID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	uid := uint(userID.(float64))
	data, err := t.rcUseCase.GetMenu(uid)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	var accessKeys []string
	for _, item := range data {
		accessKeys = append(accessKeys, item.AccessKey)
	}

	// permissions := []string{
	// 	"road_condition_manage_data", "road_condition_access",
	// }
	// if helpers.HasPermission(permissions, accessKeys) {
	roadConditionLists, err := t.rcUseCase.GetRoadConditionList(accessKeys, roadID)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(roadConditionLists, 200))
	// } else {
	// 	errResponse := responses.FailRespone(errors.New("access denied"))
	// 	c.JSON(http.StatusBadRequest, errResponse)
	// 	return
	// }

}

// @summary
// @description
// @tags RoadConditions
// @id GetRoadConditionDetails
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id_parent path string true "Insert your id_parent"
// @param condition_rang_type query string true "Insert your condition_rang_type"
// @response 200 {object} responses.Success{data=[]responses.RoadDamageImport}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/condition_details/{id_parent} [get]
func (t *RoadConditionHandler) GetRoadConditionDetails(c *gin.Context) {

	userID, _ := c.Get("userID")
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

	uid := uint(userID.(float64))
	data, err := t.rcUseCase.GetMenu(uid)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	var accessKeys []string
	for _, item := range data {
		accessKeys = append(accessKeys, item.AccessKey)
	}

	// permissions := []string{
	// 	"road_condition_manage_data", "road_condition_access",
	// }

	// if helpers.HasPermission(permissions, accessKeys) {

	conditionLists, err := t.rcUseCase.GetRoadConditionDetails(uid, req.ConditionRangeType, IDParent)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(conditionLists, 200))
	// } else {
	// 	errResponse := responses.FailRespone(errors.New("access denied"))
	// 	c.JSON(http.StatusBadRequest, errResponse)
	// 	return
	// }

}

// @summary
// @description
// @tags RoadConditions
// @id CreateRoadCondition
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your Road ID"
// @param lane_no formData string false "Insert your lane_no"
// @param surveyed_date formData string false "Insert your surveyed_date"
// @param remarks formData string false "Insert your remarks"
// @param iri_filename formData file false "upload your iri_filename"
// @param image_filename formData file false "upload your image_filename"
// @response 201 {object} responses.Success{data=responses.RoadConditionCreate} "Created"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/condition [post]
func (t *RoadConditionHandler) CreateRoadCondition(c *gin.Context) {

	var req requests.RoadCondition

	roadID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	} else {
		req.RoadID = roadID
	}

	// if err := c.ShouldBind(&req); err != nil {
	// 	if validateErr := validator.Validate(req); validateErr != nil {
	// 		err := helpers.ConverstError(validateErr)
	// 		errResponse := responses.ValidateResponse(err)
	// 		c.JSON(http.StatusUnprocessableEntity, errResponse)
	// 		return
	// 	}
	// }

	// if validateErr := validator.Validate(req); validateErr != nil {
	// 	err := helpers.ConverstError(validateErr)
	// 	errResponse := responses.ValidateResponse(err)
	// 	c.JSON(http.StatusUnprocessableEntity, errResponse)
	// 	return
	// }

	var reqFile requests.RoadConditionFiles
	imgFile, _ := c.FormFile("image_filename")
	if imgFile != nil {
		reqFile.ImageFilename = imgFile
	}

	remarks, has := c.GetPostForm("remarks")
	if has {
		req.Remarks = remarks
	}

	surveyedDate, hasSurveyedDate := c.GetPostForm("surveyed_date")
	if hasSurveyedDate {
		req.SurveyedDate = surveyedDate
	}

	LaneNo, hasLaneNo := c.GetPostForm("lane_no")
	if hasLaneNo {
		laneNo, _ := strconv.Atoi(LaneNo)
		req.LaneNo = laneNo
	}

	iriFile, err := c.FormFile("iri_filename")
	if iriFile != nil {
		reqFile.IriFilename = iriFile

	}

	if hasSurveyedDate && surveyedDate == "" || hasLaneNo && LaneNo == "" || err != nil {
		errMsg := ""

		if hasSurveyedDate && surveyedDate == "" {
			errMsg += "surveyed_date : zero value, "
		}
		if hasLaneNo && LaneNo == "" {
			errMsg += "lane_no : zero value, "
		}
		if err != nil {
			errMsg += "iri_filename : zero value"
		}

		err := helpers.ConverstError(errors.New(errMsg))
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))
	data, err := t.rcUseCase.GetMenu(uid)
	if err != nil {
		err := fmt.Errorf("userID", err.Error())
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return

	}

	var accessKeys []string
	for _, item := range data {
		accessKeys = append(accessKeys, item.AccessKey)
	}

	// permissions := []string{
	// 	"road_condition_manage_data",
	// }

	// if helpers.HasPermission(permissions, accessKeys) {

	var resp responses.RoadConditionCreate
	rcId, err := t.rcUseCase.CreateRoadCondition(c, uid, req, reqFile)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}
	resp.Id = rcId

	c.JSON(http.StatusCreated, responses.SuccessResponse(resp, 201))
	// } else {
	// 	errResponse := responses.FailRespone(errors.New("access denied"))
	// 	c.JSON(http.StatusBadRequest, errResponse)
	// 	return
	// }
}

// @summary
// @description
// @tags RoadConditions
// @id UpdateRoadCondition
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your Road ID"
// @param id_parent path string true "Insert your id_parent"
// @param lane_no formData string false "Insert your lane_no"
// @param surveyed_date formData string false "Insert your surveyed_date"
// @param remarks formData string false "Insert your remarks"
// @param iri_filename formData file false "upload your iri_filename"
// @param image_filename formData file false "upload your image_filename"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @response 200 {object} responses.Success{data=responses.RoadConditionUpdate}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/condition/{id_parent} [put]
func (t *RoadConditionHandler) UpdateRoadCondition(c *gin.Context) {

	var req requests.RoadConditionUpdate

	// if err := c.ShouldBind(&req); err != nil {
	// 	if validateErr := validator.Validate(req); validateErr != nil {
	// 		err := helpers.ConverstError(validateErr)
	// 		errResponse := responses.ValidateResponse(err)
	// 		c.JSON(http.StatusUnprocessableEntity, errResponse)
	// 		return
	// 	}
	// }

	// if validateErr := validator.Validate(req); validateErr != nil {
	// 	err := helpers.ConverstError(validateErr)
	// 	errResponse := responses.ValidateResponse(err)
	// 	c.JSON(http.StatusUnprocessableEntity, errResponse)
	// 	return
	// }
	validateMsgArr := map[string]string{}
	roadID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	} else {
		req.RoadID = roadID
	}

	var reqFile requests.RoadConditionFiles
	imgFile, _ := c.FormFile("image_filename")
	if imgFile != nil {
		reqFile.ImageFilename = imgFile
	}

	iriFile, err := c.FormFile("iri_filename")
	if iriFile != nil {
		reqFile.IriFilename = iriFile
	}

	remarks, has := c.GetPostForm("remarks")
	if has {
		req.Remarks = remarks
	}

	// hariphan แก้เรื่อง validate
	// hariphan เพิ่ม status iri_filename_status, image_filename_status
	iriFilenameStatus, _ := c.GetPostForm("iri_filename_status")
	imageFilenameStatus, _ := c.GetPostForm("image_filename_status")

	switch iriFilenameStatus {
	case "no_file":
		validateMsgArr["iri_filename"] = "โปรดระบุ"
	case "delete":
		validateMsgArr["iri_filename"] = "โปรดระบุ"
	case "":
		validateMsgArr["iri_filename"] = "โปรดระบุ"
	}

	lane, hasLane := c.GetPostForm("lane_no")
	if hasLane {
		req.LaneNo, _ = strconv.Atoi(lane)
	}

	surveyedDate, hasSurveyedDate := c.GetPostForm("surveyed_date")
	if hasSurveyedDate {
		req.SurveyedDate = surveyedDate
	}

	// hariphan แก้เรื่อง validate
	if hasSurveyedDate && surveyedDate == "" || err != nil {
		if hasSurveyedDate && surveyedDate == "" {
			validateMsgArr["surveyed_date"] = "โปรดระบุ"
		}
		if hasLane && lane == "" {
			validateMsgArr["lane"] = "โปรดระบุ"
		}
	}
	if len(validateMsgArr) > 0 {
		errResponse := responses.ValidateResponse(validateMsgArr)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	idParent, _ := strconv.Atoi(c.Params.ByName("id_parent"))
	err = c.ShouldBind(&idParent)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))
	data, err := t.rcUseCase.GetMenu(uid)
	if err != nil {
		err := fmt.Errorf("userID", err.Error())
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	var accessKeys []string
	for _, item := range data {
		accessKeys = append(accessKeys, item.AccessKey)
	}

	// permissions := []string{
	// 	"road_condition_manage_data",
	// }

	// if helpers.HasPermission(permissions, accessKeys) {
	// hariphan เพิ่ม iriFilenameStatus, imageFilenameStatus
	resp, err := t.rcUseCase.UpdateRoadCondition(c, uid, idParent, req, reqFile, iriFilenameStatus, imageFilenameStatus)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusAccepted, responses.SuccessResponse(resp, 202))
	// } else {
	// 	errResponse := responses.FailRespone(errors.New("access denied"))
	// 	c.JSON(http.StatusBadRequest, errResponse)
	// 	return
	// }
}

// @summary
// @description
// @tags RoadConditions
// @id GetRoadCondition
// @param id_parent path string true "Insert your id_parent"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RoadCondition}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/condition/{id_parent} [get]
func (t *RoadConditionHandler) GetRoadCondition(c *gin.Context) {

	idParent, _ := strconv.Atoi(c.Params.ByName("id_parent"))
	err := c.ShouldBind(&idParent)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))
	data, err := t.rcUseCase.GetMenu(uid)
	if err != nil {
		err := fmt.Errorf("userID", err.Error())
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	var accessKeys []string
	for _, item := range data {
		accessKeys = append(accessKeys, item.AccessKey)
	}

	// permissions := []string{
	// 	"road_condition_manage_data",
	// }

	// if helpers.HasPermission(permissions, accessKeys) {

	resp, err := t.rcUseCase.GetRoadCondition(c, uid, idParent)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, 200))
	// } else {
	// 	errResponse := responses.FailRespone(errors.New("access denied"))
	// 	c.JSON(http.StatusBadRequest, errResponse)
	// 	return
	// }
}

// @summary
// @description
// @tags RoadConditions
// @id DeleteRoadCondition
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id_parent path string true "Insert your id_parent"
// @response 204 "No Content"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/condition/{id_parent} [delete]
func (t *RoadConditionHandler) DeleteRoadCondition(c *gin.Context) {

	idParent, _ := strconv.Atoi(c.Params.ByName("id_parent"))
	err := c.ShouldBind(&idParent)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))
	data, err := t.rcUseCase.GetMenu(uid)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	var accessKeys []string
	for _, item := range data {
		accessKeys = append(accessKeys, item.AccessKey)
	}

	// permissions := []string{
	// 	"road_condition_manage_data",
	// }

	// if helpers.HasPermission(permissions, accessKeys) {

	_, err = t.rcUseCase.DeleteRoadCondition(c, uid, idParent)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(http.StatusNoContent, responses.SuccessResponse(nil, 204))
	// } else {
	// 	errResponse := responses.FailRespone(errors.New("access denied"))
	// 	c.JSON(http.StatusBadRequest, errResponse)
	// 	return
	// }
}

// @summary
// @description
// @tags RoadConditions
// @id GetRoadConditionCompareLane
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @param years query string false "Insert your years"
// @param lanes query string false "Insert your lanes"
// @param condition_type query string false "Insert your condition_type"
// @response 200 {object} responses.Success{data=[]responses.RoadConditionLane}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/condition_compare_lane [get]
func (t *RoadConditionHandler) GetRoadConditionCompareLane(c *gin.Context) {

	var req requests.RoadConditionCompare
	if err := c.ShouldBindQuery(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	roadId, _ := strconv.Atoi(c.Params.ByName("id"))
	err := c.ShouldBind(&roadId)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))
	data, err := t.rcUseCase.GetMenu(uid)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	var accessKeys []string
	for _, item := range data {
		accessKeys = append(accessKeys, item.AccessKey)
	}

	// permissions := []string{
	// 	"road_condition_manage_data", "road_condition_access",
	// }

	// if helpers.HasPermission(permissions, accessKeys) {

	resp, err := t.rcUseCase.GetRoadConditionCompareLane(roadId, req)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(int(http.StatusOK), responses.SuccessResponse(resp, 200))
	// } else {
	// 	errResponse := responses.FailRespone(errors.New("access denied"))
	// 	c.JSON(http.StatusBadRequest, errResponse)
	// 	return
	// }
}

// @summary
// @description
// @tags RoadConditions
// @id GetRoadConditionCompareYear
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @param years query string false "Insert your years"
// @param lanes query string false "Insert your lanes"
// @param condition_type query string false "Insert your condition_type"
// @response 200 {object} responses.Success{data=[]responses.RoadConditionYear}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/condition_compare_year [get]
func (t *RoadConditionHandler) GetRoadConditionCompareYear(c *gin.Context) {

	var req requests.RoadConditionCompare
	if err := c.ShouldBindQuery(&req); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	roadId, _ := strconv.Atoi(c.Params.ByName("id"))
	err := c.ShouldBind(&roadId)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))
	data, err := t.rcUseCase.GetMenu(uid)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	var accessKeys []string
	for _, item := range data {
		accessKeys = append(accessKeys, item.AccessKey)
	}

	// permissions := []string{
	// 	"road_condition_manage_data", "road_condition_access",
	// }

	// if helpers.HasPermission(permissions, accessKeys) {

	resp, err := t.rcUseCase.GetRoadConditionCompareYear(roadId, req)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(int(http.StatusOK), responses.SuccessResponse(resp, 200))
	// } else {
	// 	errResponse := responses.FailRespone(errors.New("access denied"))
	// 	c.JSON(http.StatusBadRequest, errResponse)
	// 	return
	// }
}

// @summary
// @description
// @tags RoadConditions
// @id GetRoadConditionCompareAverage
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @param lane path string true "Insert your lane"
// @response 200 {object} responses.Success{data=[]responses.RoadConditionAverage}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/condition_compare_average/{lane} [get]
func (t *RoadConditionHandler) GetRoadConditionCompareAverage(c *gin.Context) {

	lane, _ := strconv.Atoi(c.Params.ByName("lane"))
	err := c.ShouldBind(&lane)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	roadId, _ := strconv.Atoi(c.Params.ByName("id"))
	err = c.ShouldBind(&roadId)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))
	data, err := t.rcUseCase.GetMenu(uid)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	var accessKeys []string
	for _, item := range data {
		accessKeys = append(accessKeys, item.AccessKey)
	}

	// permissions := []string{
	// 	"road_condition_manage_data", "road_condition_access",
	// }

	// if helpers.HasPermission(permissions, accessKeys) {

	resp, err := t.rcUseCase.GetRoadConditionCompareAverage(roadId, lane)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(int(http.StatusOK), responses.SuccessResponse(resp, 200))
	// } else {
	// 	errResponse := responses.FailRespone(errors.New("access denied"))
	// 	c.JSON(http.StatusBadRequest, errResponse)
	// 	return
	// }
}
