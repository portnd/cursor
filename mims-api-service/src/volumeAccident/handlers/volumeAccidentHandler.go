package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/volumeAccident/domains"
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

// ================================== start function  ==================================

// @summary
// @description
// @tags volume accident
// @id GetVolumeRevisionAccident
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road group ID"
// @response 200 {object} responses.Success{data=responses.VolumeAadtRevision} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/road_group/{id}/volume_accident/revision [get]
func (h *Handler) GetVolumeRevision(c *gin.Context) {
	roadGrpID, _ := strconv.Atoi(c.Params.ByName("id"))
	permission := helpers.GetAccessControl(c)
	data, err := h.UseCase.GetVolumeRevision(roadGrpID, permission)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
}

// @summary
// @description
// @tags volume accident
// @id GetVolumeAccident
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road group ID"
// @param accident_id path string true "Insert your accident_id"
// @response 200 {object} responses.Success{data=[]models.VolumeAadt} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/road_group/{id}/volume_accident/{accident_id} [get]
func (h *Handler) GetVolume(c *gin.Context) {
	roadGrpID, _ := strconv.Atoi(c.Params.ByName("id"))
	ID, _ := strconv.Atoi(c.Params.ByName("accident_id"))
	// permission := helpers.GetAccessControl(c)
	data, err := h.UseCase.GetVolume(roadGrpID, ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
}

// @summary
// @description
// @tags volume accident
// @id CreateVolumeAccident
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road group ID"
// @param Asset body requests.VolumeAccidentReq true "Update your data"
// @response 200 {object} []responses.Success{data=responses.Volume} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/road_group/{id}/volume_accident [post]
func (h *Handler) CreateVolume(c *gin.Context) {
	roadGrpID, _ := strconv.Atoi(c.Params.ByName("id"))
	userID, _ := c.Get("userID")
	var req requests.VolumeAccidentReq
	err := c.ShouldBind(&req)
	if err != nil {
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

	// permission := helpers.GetAccessControl(c)
	// data, err := h.UseCase.UpdateVolume(roadGrpID, int(userID.(float64)), req)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, err.Error())
	// 	return
	// }
	// c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
	data, err := h.UseCase.CreateVolume(roadGrpID, 0, 0, int(userID.(float64)), req)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
}

// @summary
// @description
// @tags volume accident
// @id UpdateVolumeAccident
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road group ID"
// @param id_parent path string true "Insert your id_parent"
// @param Asset body requests.VolumeAccidentReq true "Update your data"
// @response 200 {object} []responses.Success{data=responses.Volume} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/road_group/{id}/volume_accident/{id_parent}/accident/{accident_id} [put]
func (h *Handler) UpdateVolume(c *gin.Context) {
	roadGrpID, _ := strconv.Atoi(c.Params.ByName("id"))
	IDParent, _ := strconv.Atoi(c.Params.ByName("id_parent"))
	accidentID, _ := strconv.Atoi(c.Params.ByName("accident_id"))
	userID, _ := c.Get("userID")
	var req requests.VolumeAccidentReq
	err := c.ShouldBind(&req)
	if err != nil {
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

	// permission := helpers.GetAccessControl(c)
	data, err := h.UseCase.CreateVolume(roadGrpID, IDParent, accidentID, int(userID.(float64)), req)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
}

// @summary
// @description
// @tags volume accident
// @id DeleteVolumeAccident
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your road group ID"
// @param accident_id path string true "Insert your 	accident_id"
// @response 204 "No Content"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/road_group/{id}/volume_accident/{accident_id} [delete]
func (h *Handler) DeleteVolume(c *gin.Context) {
	userId, _ := c.Get("userID")
	roadGrpId, _ := strconv.Atoi(c.Params.ByName("id"))
	accidentID, _ := strconv.Atoi(c.Params.ByName("accident_id"))
	_, err := h.UseCase.DeleteVolume(roadGrpId, accidentID, int(userId.(float64)))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	c.Status(204)
}
