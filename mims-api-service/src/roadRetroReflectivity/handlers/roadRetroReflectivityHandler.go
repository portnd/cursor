package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadRetroReflectivity/domains"
	"gopkg.in/validator.v2"
)

// init handler

type Handler struct {
	UseCase domains.UseCase
}

// init handler
func NewHandler(usecase domains.UseCase) *Handler {
	return &Handler{
		UseCase: usecase,
	}
}

// @summary
// @description
// @tags RoadRetroReflectivity
// @id GetRoadRetroReflectivityList
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @response 200 {object} responses.Success{data=[]responses.RoadRetroReflectivityList}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/retro_reflectivity/list [get]
func (t *Handler) GetRoadRetroReflectivityList(c *gin.Context) {

	roadID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	roadConditionLists, err := t.UseCase.GetRoadRetroReflectivityList(roadID)
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
// @tags RoadRetroReflectivity
// @id GetRoadRetroReflectivityDetails
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id_parent path string true "Insert your id_parent"
// @param range_type query string true "Insert your ref_reflectivity_range"
// @response 200 {object} responses.Success{data=[]responses.RetroReflectivityDetails}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/retro_reflectivity/details/{id_parent} [get]
func (t *Handler) GetRoadRetroReflectivityDetails(c *gin.Context) {

	userID, _ := c.Get("userID")
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

	uid := uint(userID.(float64))

	RetroReflectivityDetails, err := t.UseCase.GetRoadRetroReflectivityDetails(uid, req.RangeType, req.RefStripeTypeIDs, IDParent)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(RetroReflectivityDetails, 200))

}

// @summary
// @description
// @tags RoadRetroReflectivity
// @id CreateRoadRetroReflectivity
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your Road ID"
// @param line_no formData string false "Insert your line_no"
// @param surveyed_date formData string false "Insert your surveyed_date"
// @param remarks formData string false "Insert your remarks"
// @param csv_file formData file false "upload your csv_file"
// @response 201 {object} responses.Success{data=responses.RoadRetroReflectivityCreate} "Created"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/retro_reflectivity [post]
func (t *Handler) CreateRoadRetroReflectivity(c *gin.Context) {

	var req requests.RoadRetroReflectivity

	roadID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	} else {
		req.RoadID = roadID
	}

	var reqFile requests.RoadRetroReflectivityFiles

	remarks, has := c.GetPostForm("remarks")
	if has {
		req.Remarks = remarks
	}

	surveyedDate, hasSurveyedDate := c.GetPostForm("surveyed_date")
	if hasSurveyedDate {
		req.SurveyedDate = surveyedDate
	}

	LineNo, hasLineNo := c.GetPostForm("line_no")
	if hasLineNo {
		lineNo, _ := strconv.Atoi(LineNo)
		req.LineNo = lineNo
	}

	csvFile, err := c.FormFile("csv_file")
	if csvFile != nil {
		reqFile.CsvFilename = csvFile

	}

	if hasSurveyedDate && surveyedDate == "" || hasLineNo && LineNo == "null" || err != nil {
		errMsg := ""

		if hasSurveyedDate && surveyedDate == "" {
			errMsg += "surveyed_date : zero value, "
		}
		if hasLineNo && LineNo == "null" {
			errMsg += "line_no : zero value, "
		}
		if err != nil {
			errMsg += "csv_file : zero value"
		}

		err := helpers.ConverstError(errors.New(errMsg))
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))

	resp, err := t.UseCase.CreateRoadRetroReflectivity(c, uid, req, reqFile)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(resp, 201))

}

// @summary
// @description
// @tags RoadRetroReflectivity
// @id UpdateRoadRetroReflectivity
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your Road ID"
// @param id_parent path string true "Insert your id_parent"
// @param line_no formData string false "Insert your line_no"
// @param surveyed_date formData string false "Insert your surveyed_date"
// @param remarks formData string false "Insert your remarks"
// @param csv_file formData file false "upload your csv_file"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @response 200 {object} responses.Success{data=responses.RoadConditionUpdate}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/retro_reflectivity/{id_parent} [put]
func (t *Handler) UpdateRoadRetroReflectivity(c *gin.Context) {

	var req requests.RoadRetroReflectivity

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

	var reqFile requests.RoadRetroReflectivityFiles

	csv, err := c.FormFile("csv_file")
	if csv != nil {
		reqFile.CsvFilename = csv
	}

	remarks, has := c.GetPostForm("remarks")
	if has {
		req.Remarks = remarks
	}

	// hariphan แก้เรื่อง validate
	// hariphan เพิ่ม status iri_filename_status, image_filename_status
	csvFileStatus, _ := c.GetPostForm("csv_file_status")

	switch csvFileStatus {
	case "no_file":
		validateMsgArr["csv_file"] = "โปรดระบุ"
	case "delete":
		validateMsgArr["csv_file"] = "โปรดระบุ"
	case "":
		validateMsgArr["csv_file"] = "โปรดระบุ"
	}

	line, hasLine := c.GetPostForm("line_no")
	if hasLine {
		req.LineNo, _ = strconv.Atoi(line)
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
		if hasLine && line == "" {
			validateMsgArr["line_no"] = "โปรดระบุ"
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

	// hariphan เพิ่ม iriFilenameStatus, imageFilenameStatus
	resp, err := t.UseCase.UpdateRoadRetroReflectivity(c, uid, idParent, req, reqFile, csvFileStatus)
	if err != nil {
		appErr, _ := err.(*responses.AppErr)
		errResponse := responses.FailRespone(appErr)
		c.JSON(appErr.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusAccepted, responses.SuccessResponse(resp, 202))

}

// @summary
// @description
// @tags RoadRetroReflectivity
// @id GetRoadRetroReflectivity
// @param id_parent path string true "Insert your id_parent"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RoadRetroReflectivity}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/retro_reflectivity/{id_parent} [get]
func (t *Handler) GetRoadRetroReflectivity(c *gin.Context) {

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

	resp, err := t.UseCase.GetRoadRetroReflectivity(c, uid, idParent)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, 200))

}

// @summary
// @description
// @tags RoadRetroReflectivity
// @id DeleteRoadRetroReflectivity
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id_parent path string true "Insert your id_parent"
// @response 204 "No Content"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/retro_reflectivity/{id_parent} [delete]
func (t *Handler) DeleteRoadRetroReflectivity(c *gin.Context) {

	idParent, _ := strconv.Atoi(c.Params.ByName("id_parent"))
	err := c.ShouldBind(&idParent)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))

	_, err = t.UseCase.DeleteRoadRetroReflectivity(c, uid, idParent)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(http.StatusNoContent, responses.SuccessResponse(nil, 204))

}

// @summary
// @description
// @tags RoadRetroReflectivity
// @id GetRoadRetroReflectivityCompareLine
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @param years query string false "Insert your years"
// @param lines query string false "Insert your lines"
// @param condition_type query string false "Insert your condition_type"
// @response 200 {object} responses.Success{data=[]responses.RoadRetroReflectivityLine}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/retro_reflectivity/compare_line [get]
func (t *Handler) GetRoadRetroReflectivityCompareLine(c *gin.Context) {

	var req requests.RoadRetroReflectivityCompare
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

	//userID, _ := c.Get("userID")
	//uid := uint(userID.(float64))

	resp, err := t.UseCase.GetRoadRetroReflectivityCompareLine(roadId, req)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(int(http.StatusOK), responses.SuccessResponse(resp, 200))

}

// @summary
// @description
// @tags RoadRetroReflectivity
// @id GetRoadRetroReflectivityCompareYear
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @param years query string false "Insert your years"
// @param lines query string false "Insert your lines"
// @response 200 {object} responses.Success{data=[]responses.RoadRetroReflectivityYear}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/retro_reflectivity/compare_year [get]
func (t *Handler) GetRoadRetroReflectivityCompareYear(c *gin.Context) {

	var req requests.RoadRetroReflectivityCompare
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

	//userID, _ := c.Get("userID")

	resp, err := t.UseCase.GetRoadRetroReflectivityCompareYear(roadId, req)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(int(http.StatusOK), responses.SuccessResponse(resp, 200))

}

// @summary
// @description
// @tags RoadRetroReflectivity
// @id GetRoadRetroReflectivityCompareAverage
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @param lane path string true "Insert your lane"
// @response 200 {object} responses.Success{data=[]responses.RoadRetroReflectivityAverage}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/retro_reflectivity/compare_average/{lane} [get]
func (t *Handler) GetRoadRetroReflectivityCompareAverage(c *gin.Context) {

	line, _ := strconv.Atoi(c.Params.ByName("line"))
	err := c.ShouldBind(&line)
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

	resp, err := t.UseCase.GetRoadRetroReflectivityCompareAverage(roadId, line)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(int(http.StatusOK), responses.SuccessResponse(resp, 200))

}

// @summary
// @description
// @tags RoadRetroReflectivity
// @id GetRoadRetroReflectivityTemplate
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @response 200 {object} responses.Success{data}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/retro_reflectivity/template [get]
func (t *Handler) GetRoadRetroReflectivityTemplate(c *gin.Context) {

	roadId, _ := strconv.Atoi(c.Params.ByName("id"))
	err := c.ShouldBind(&roadId)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	userID, _ := c.Get("userID")
	uid := uint(userID.(float64))

	resp, err := t.UseCase.GetRoadRetroReflectivityTemplate(uid, roadId)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	c.JSON(int(http.StatusOK), responses.SuccessResponse(resp, 200))

}

// @summary
// @description
// @tags RoadRetroReflectivity
// @id GetRoadRetroReflectivityLineList
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @response 200 {object} []responses.RoadRetroReflectivityLineList "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/retro_reflectivity/line_list [get]
func (t Handler) GetRoadRetroReflectivityLineList(c *gin.Context) {
	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	data, _ := t.UseCase.GetRoadRetroReflectivityLineList(roadID)
	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
}
