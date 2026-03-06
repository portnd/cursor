package handlers

import (
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadDamage/domains"
	"gopkg.in/validator.v2"

	"github.com/gin-gonic/gin"
)

// init handler
type RoadDamageHandler struct {
	roadDamageUseCase domains.RoadDamageUseCase
}

// init handler
func NewRoadDamageHandler(usecase domains.RoadDamageUseCase) *RoadDamageHandler {
	return &RoadDamageHandler{
		roadDamageUseCase: usecase,
	}
}

// @summary
// @description
// @tags Roads Damage
// @id GetRoadDamageList
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @response 200 {object} []responses.Success{data=responses.RoadDamageList}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/damage_list [get]
func (t *RoadDamageHandler) GetRoadDamageList(c *gin.Context) {
	//-------- Permission --------
	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	permission := helpers.GetAccessControl(c)
	roadDamage, err := t.roadDamageUseCase.GetRoadDamageList(roadID, permission)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(roadDamage, http.StatusOK))
}

// @summary
// @description
// @tags Roads Damage
// @id GetRoadDamageForImport
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @param id_parent path string true "Insert your id_parent"
// @response 200 {object} []responses.Success{data=responses.RoadDamageImport}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/damage_import/{id_parent} [get]
func (t *RoadDamageHandler) GetRoadDamageForImport(c *gin.Context) {
	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	IDParent, _ := strconv.Atoi(c.Params.ByName("id_parent"))

	direction, err := t.roadDamageUseCase.GetDirectionById(roadID)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	data, err := t.roadDamageUseCase.GetRoadDamageForImport(roadID, IDParent)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	var resp responses.RoadDamageImport
	resp.ID = data.Id
	resp.IDParent = data.IdParent
	resp.LaneNo = data.LaneNo
	resp.SurveyedDate = (data.SurveyedDate)
	resp.DamageFilenameint = os.Getenv("STORAGE_IP") + "/" + data.DamageInputFilepath
	if data.ImgFilepath != "" {
		resp.ImgFilepath = os.Getenv("STORAGE_IP") + "/" + data.ImgFilepath
	} else {
		resp.ImgFilepath = ""
	}

	resp.Direction.ID = direction.ID
	resp.Direction.Name = direction.Name
	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags Roads Damage
// @id GetRoadDamageDetail
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @param id_parent path string true "Insert your id_parent"
// @response 200 {object} []responses.Success{data=responses.RoadDamageListDetail}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/damage_detail/{id_parent} [get]
func (t *RoadDamageHandler) GetRoadDamageDetail(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Params.ByName("id"))
	idParent, _ := strconv.Atoi(c.Params.ByName("id_parent"))
	permission := helpers.GetAccessControl(c)
	dart, err := t.roadDamageUseCase.GetRoadDamageDetail(roleID, idParent, permission)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(dart, http.StatusOK))
}

// @summary
// @description
// @tags Roads Damage
// @id GetRoadDamageTemplate
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} []responses.Success{data}  "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/damage_template [get]
func (t *RoadDamageHandler) GetRoadDamageTemplate(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Params.ByName("id"))
	dart, err := t.roadDamageUseCase.GetRoadDamageTemplate(roleID)
	// helpers.DownloadHandler(c, dart.(string))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(dart, http.StatusOK))
}

// @summary
// @description
// @tags Roads Damage
// @id SetRoadDamageFromImport
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @param lane_no formData string false "Insert your lane_no"
// @param surveyed_date formData string false "Insert your surveyed_date"
// @param damage_filename formData file false "upload your damage_filename"
// @param image_filename formData file false "upload your image_filename"
// @response 201 {object} []responses.Success{data=responses.RoadDamageSetData} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/damage_import [post]
func (t *RoadDamageHandler) SetRoadDamageFromImport(c *gin.Context) {
	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	dstPath := ""
	imgPath := ""
	userId, _ := c.Get("userID")
	var req requests.RoadDamageImport
	err := c.ShouldBind(&req)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	validateMsgArr := map[string]string{}
	if validateErr := validator.Validate(req); validateErr != nil {
		validateMsgArr = helpers.ConverstError(validateErr)
	}

	damageFilename, _ := c.FormFile("damage_filename")
	if damageFilename == nil {
		validateMsgArr["damage_filename"] = "โปรดระบุ"
	}

	if len(validateMsgArr) > 0 {
		errResponse := responses.ValidateResponse(validateMsgArr)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	imageFilename, _ := c.FormFile("image_filename")

	if damageFilename != nil {
		dstPath, err = helpers.SaveFile(c, damageFilename, os.Getenv("ROAD_DAMAGE_CSV_DIR"))
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(400, errResponse)
			return
		}
	}

	if imageFilename != nil {
		imgPath, err = helpers.SaveFile(c, imageFilename, os.Getenv("ROAD_DAMAGE_IMG_DIR"))
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(400, errResponse)
			return
		}
	}

	// check csv data
	_, err = t.roadDamageUseCase.RoadDamageReadCSVFile(roadID, dstPath)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	var rcData requests.RcData
	surveyedDate, err := time.Parse("2006-01-02 15:04:05", req.SurveyedDate)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	rcData.SurveyedDate = surveyedDate
	rcData.Year = surveyedDate.Year()
	rcData.CreatedDate = time.Now()
	rcData.CreatedBy = int(userId.(float64))
	rcData.UpdatedDate = time.Now()
	rcData.UpdatedBy = int(userId.(float64))
	rcData.LaneNo = req.LaneNo
	rcData.IDParent = 0
	rcData.Status = "A"
	rcData.Revision = 0
	idParent := 0
	data, err := t.roadDamageUseCase.SetRoadDamageFromImport(roadID, idParent, int(userId.(float64)), dstPath, imgPath, req, rcData, damageFilename, imageFilename)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(data, http.StatusCreated))
}

// @summary
// @description
// @tags Roads Damage
// @id UpdateRoadDamageFromImport
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @param id_parent path string true "Insert your id_parent"
// @param lane_no formData string false "Insert your lane_no"
// @param surveyed_date formData string false "Insert your surveyed_date"
// @param damage_filename formData file false "upload your damage_filename"
// @param image_filename formData file false "upload your image_filename"
// @response 200 {object} []responses.Success{data=responses.RoadDamageSetData} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/damage_import/{id_parent} [put]
func (t *RoadDamageHandler) UpdateRoadDamageFromImport(c *gin.Context) {
	roadID, _ := strconv.Atoi(c.Params.ByName("id"))
	idParent, _ := strconv.Atoi(c.Params.ByName("id_parent"))
	csvPath := ""
	imgPath := ""
	userId, _ := c.Get("userID")
	var req requests.RoadDamageImport
	err := c.ShouldBind(&req)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	validateMsgArr := map[string]string{}
	if validateErr := validator.Validate(req); validateErr != nil {
		validateMsgArr = helpers.ConverstError(validateErr)
	}

	// var damageFilename multipart.FileHeader

	damageFilename, _ := c.FormFile("damage_filename")
	if req.DamageFilenameStatus == "upload" {
		if len(validateMsgArr) > 0 {
			errResponse := responses.ValidateResponse(validateMsgArr)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}

		if damageFilename != nil {
			csvPath, err = helpers.SaveFile(c, damageFilename, os.Getenv("ROAD_DAMAGE_CSV_DIR"))
			if err != nil {
				errResponse := responses.FailRespone(err)
				c.JSON(400, errResponse)
				return
			}
		}
		// check csv data
		_, err = t.roadDamageUseCase.RoadDamageReadCSVFile(roadID, csvPath)
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(400, errResponse)
			return
		}
	} else if req.DamageFilenameStatus == "no_file" || req.DamageFilenameStatus == "delete" {
		validateMsgArr["damage_filename"] = "โปรดระบุ"
		if len(validateMsgArr) > 0 {
			errResponse := responses.ValidateResponse(validateMsgArr)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}
	}

	imageFilename, _ := c.FormFile("image_filename")
	if imageFilename != nil {
		imgPath, err = helpers.SaveFile(c, imageFilename, os.Getenv("ROAD_DAMAGE_IMG_DIR"))
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(400, errResponse)
			return
		}
	}

	var rcData requests.RcData
	surveyedDate, err := time.Parse("2006-01-02 15:04:05", req.SurveyedDate)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}

	roadDamage, err := t.roadDamageUseCase.GetRoadDamageByIDParent(idParent)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	rcData.SurveyedDate = surveyedDate
	rcData.Year = surveyedDate.Year()
	rcData.CreatedDate = roadDamage.CreatedDate
	rcData.CreatedBy = roadDamage.CreatedBy
	rcData.UpdatedDate = time.Now()
	rcData.UpdatedBy = int(userId.(float64))
	rcData.LaneNo = roadDamage.LaneNo
	rcData.IDParent = idParent
	rcData.Status = "A"
	rcData.Revision = roadDamage.Revision + 1
	if roadDamage.Status == "A" {
		// check have csv file
		if req.DamageFilenameStatus == "not_edit" {
			if roadDamage.Status == "A" {
				if err := t.roadDamageUseCase.UpdateRoadDamage(roadDamage.Id, int(userId.(float64)), req); err != nil {
					errResponse := responses.FailRespone(err)
					c.JSON(400, errResponse)
					return
				}
			}
		}
	}

	// int(userId.(float64))
	if damageFilename == nil {
		var damageFilename *multipart.FileHeader
		data, err := t.roadDamageUseCase.SetRoadDamageFromImport(roadID, idParent, int(userId.(float64)), csvPath, imgPath, req, rcData, damageFilename, imageFilename)
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(400, errResponse)
			return
		}

		if req.DamageFilenameStatus == "upload" {
			err = t.roadDamageUseCase.UpdateRoadDamageStatusI(roadDamage.Id, req)
			if err != nil {
				errResponse := responses.FailRespone(err)
				c.JSON(400, errResponse)
				return
			}
		}

		c.JSON(http.StatusAccepted, responses.SuccessResponse(data, http.StatusAccepted))
	} else {
		data, err := t.roadDamageUseCase.SetRoadDamageFromImport(roadID, idParent, int(userId.(float64)), csvPath, imgPath, req, rcData, damageFilename, imageFilename)
		if err != nil {
			errResponse := responses.FailRespone(err)
			c.JSON(400, errResponse)
			return
		}

		// Updat roadDamage.Id
		if req.DamageFilenameStatus == "upload" {
			err = t.roadDamageUseCase.UpdateRoadDamageStatusI(roadDamage.Id, req)
			if err != nil {
				errResponse := responses.FailRespone(err)
				c.JSON(400, errResponse)
				return
			}
		}

		c.JSON(http.StatusAccepted, responses.SuccessResponse(data, http.StatusAccepted))
	}
}

// @summary
// @description
// @tags Roads Damage
// @id DeleteRoadDamageForImport
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road ID"
// @param id_parent path string true "Insert your id_parent"
// @response 204 "No Content"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roads/{id}/damage_import/{id_parent} [delete]
func (t *RoadDamageHandler) DeleteRoadDamageForImport(c *gin.Context) {
	userId, _ := c.Get("userID")
	idParent, _ := strconv.Atoi(c.Params.ByName("id_parent"))
	_, err := t.roadDamageUseCase.DeleteRoadDamageForImport(idParent, int(userId.(float64)))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	c.Status(204)
}
