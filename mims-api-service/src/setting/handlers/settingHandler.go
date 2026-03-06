package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/setting/domains"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gopkg.in/validator.v2"
)

type SettingHandler struct {
	settingUseCase domains.SettingUseCase
}

func NewSettingHandler(usecase domains.SettingUseCase) *SettingHandler {
	return &SettingHandler{settingUseCase: usecase}
}

// @summary
// @description
// @tags setting
// @id get_asset_groups
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param page query string false "Insert your page number"
// @param limit query string false "Insert your limit number"
// @param name query string false "Insert your asset name"
// @response 200 {object} responses.Success{data=responses.Pagination{items=[]models.RefAsset}} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/asset_groups [get]
func (sh *SettingHandler) GetAssetGroups(c *gin.Context) {
	var params requests.QueryParams
	if err := c.ShouldBind(&params); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	resp, err := sh.settingUseCase.GetAssetGroups(params)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags setting
// @id get_asset_groups_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your asset ID"
// @response 200 {object} responses.Success{data=models.RefAsset} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/asset_groups/{id} [get]
func (sh *SettingHandler) GetAssetGroupByID(c *gin.Context) {
	resp, err := sh.settingUseCase.GetAssetGroupByID(c.Param("id"))
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags setting
// @id create_asset_groups
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param AssetGroup body requests.RefRequest true "Insert your asset group name"
// @response 201 {object} responses.CreateResponse "Create Successfully"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/asset_groups [post]
func (sh *SettingHandler) CreateAssetGroup(c *gin.Context) {
	var request requests.RefRequest
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

	err := sh.settingUseCase.CreateAssetGroup(request.Name)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok && err.StatusCode == 422 {
			err := helpers.ConverstError(err)
			errResponse := responses.ValidateResponse(err)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(responses.NoData{}, http.StatusCreated))
}

// @summary
// @description
// @tags setting
// @id update_asset_groups_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your asset ID"
// @param AssetGroup body requests.RefRequest true "Insert your asset group name"
// @response 202 {object} responses.UpdateResponse "Update Successfully"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/asset_groups/{id} [put]
func (sh *SettingHandler) UpdateAssetGroupByID(c *gin.Context) {
	var request requests.RefRequest
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

	err := sh.settingUseCase.UpdateAssetGroupByID(c.Param("id"), request.Name)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok && err.StatusCode == 422 {
			err := helpers.ConverstError(err)
			errResponse := responses.ValidateResponse(err)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusAccepted, responses.SuccessResponse(responses.NoData{}, http.StatusAccepted))
}

// @summary
// @description
// @tags setting
// @id delete_asset_groups_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your asset ID"
// @response 204 "No Content"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/asset_groups/{id} [delete]
func (sh *SettingHandler) DeleteAssetGroupByID(c *gin.Context) {
	err := sh.settingUseCase.DeleteAssetGroupByID(c.Param("id"))
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// @summary
// @description
// @tags setting
// @id get_departments
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param page query string false "Insert your page number"
// @param limit query string false "Insert your limit number"
// @param name query string false "Insert your department name"
// @response 200 {object} responses.Success{data=responses.Pagination{items=[]models.RefDepartment}} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/departments [get]
func (sh *SettingHandler) GetDepartments(c *gin.Context) {
	var params requests.QueryParams
	if err := c.ShouldBind(&params); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	resp, err := sh.settingUseCase.GetDepartments(params)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags setting
// @id get_department_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your department ID"
// @response 200 {object} responses.Success{data=models.RefDepartment} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/departments/{id} [get]
func (sh *SettingHandler) GetDepartmentByID(c *gin.Context) {
	resp, err := sh.settingUseCase.GetDepartmentByID(c.Param("id"))
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags setting
// @id create_department
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param Department body requests.RefRequest true "Insert your department name"
// @response 201 {object} responses.CreateResponse "Create Successfully"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/departments [post]
func (sh *SettingHandler) CreateDepartment(c *gin.Context) {
	var request requests.RefRequest
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

	err := sh.settingUseCase.CreateDepartment(request.Name)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok && err.StatusCode == 422 {
			err := helpers.ConverstError(err)
			errResponse := responses.ValidateResponse(err)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(responses.NoData{}, http.StatusCreated))
}

// @summary
// @description
// @tags setting
// @id update_department_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your department ID"
// @param Department body requests.RefRequest true "Insert your department name"
// @response 202 {object} responses.UpdateResponse "Update Successfully"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/departments/{id} [put]
func (sh *SettingHandler) UpdateDepartmentByID(c *gin.Context) {
	var request requests.RefRequest
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

	err := sh.settingUseCase.UpdateDepartmentByID(c.Param("id"), request.Name)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok && err.StatusCode == 422 {
			err := helpers.ConverstError(err)
			errResponse := responses.ValidateResponse(err)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusAccepted, responses.SuccessResponse(responses.NoData{}, http.StatusAccepted))
}

// @summary
// @description
// @tags setting
// @id delete_department_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your department ID"
// @response 204 "No Content"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/departments/{id} [delete]
func (sh *SettingHandler) DeleteDepartmentByID(c *gin.Context) {
	err := sh.settingUseCase.DeleteDepartmentByID(c.Param("id"))
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// @summary
// @description
// @tags setting
// @id get_owners
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param page query string false "Insert your page number"
// @param limit query string false "Insert your limit number"
// @param name query string false "Insert your department name"
// @response 200 {object} responses.Success{data=responses.Pagination{items=[]responses.FirstGroup}} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/owners [get]
func (sh *SettingHandler) GetOwners(c *gin.Context) {
	var params requests.QueryParams
	if err := c.ShouldBind(&params); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	resp, err := sh.settingUseCase.GetOwners(params)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags setting
// @id get_owner_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your owner id"
// @response 200 {object} responses.Success{data=models.RefOwner} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/owners/{id} [get]
func (sh *SettingHandler) GetOwnersByID(c *gin.Context) {
	ownerID, _ := strconv.Atoi(c.Param("id"))
	_, err := sh.settingUseCase.GetOwnerByID(ownerID)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	paramsConditions, err := sh.settingUseCase.GetConditionList(ownerID)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	owner := GetOwner(paramsConditions)
	data := responses.ConditionRespond{
		ID:                  owner.ID,
		RefConditionRangeID: owner.RefConditionRangeID,
		OwnerName:           owner.Name,
		ConditionList:       GetCondition(paramsConditions),
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
}

// @summary
// @description
// @tags setting
// @id craete_owner
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param OwnerData body requests.OwnerRequest false "Insert your owner data"
// @response 201 {object} responses.CreateResponse "Create Successfully"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/owners [post]
func (sh *SettingHandler) CreateOwner(c *gin.Context) {
	// request
	var request requests.OwnerRequest
	if err := c.ShouldBind(&request); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	// validate
	if validateErr := validator.Validate(request); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	// create owner
	ownerId, err := sh.settingUseCase.CreateOwner(request)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok && err.StatusCode == 422 {
			err := helpers.ConverstError(err)
			errResponse := responses.ValidateResponse(err)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	// create Condition
	// conditionReq := request.ConditionList
	_, err = sh.settingUseCase.CreateConditionList(ownerId, request)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	paramsConditions, err := sh.settingUseCase.GetConditionList(ownerId)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	owner := GetOwner(paramsConditions)
	data := responses.ConditionRespond{
		ID:                  owner.ID,
		RefConditionRangeID: owner.RefConditionRangeID,
		OwnerName:           owner.Name,
		ConditionList:       GetCondition(paramsConditions),
	}

	c.JSON(201, responses.SuccessResponse(data, http.StatusCreated))
}

// @summary
// @description
// @tags setting
// @id update_owner_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your owner id"
// @param OwnerData body requests.OwnerRequest false "Insert your owner data"
// @response 202 {object} responses.UpdateResponse "Update Successfully"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/owners/{id} [put]
func (sh *SettingHandler) UpdateOwnerByID(c *gin.Context) {
	// request
	ownerID, _ := strconv.Atoi(c.Param("id"))
	var request requests.OwnerRequest
	if err := c.ShouldBind(&request); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	// validate
	if validateErr := validator.Validate(request); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	// create owner
	err := sh.settingUseCase.UpdateOwnerByID(ownerID, request)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok && err.StatusCode == 422 {
			err := helpers.ConverstError(err)
			errResponse := responses.ValidateResponse(err)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	// // create Condition
	// conditionReq := request.ConditionList
	err = sh.settingUseCase.UpdateConditionList(ownerID, request)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	paramsConditions, err := sh.settingUseCase.GetConditionList(ownerID)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	owner := GetOwner(paramsConditions)
	data := responses.ConditionRespond{
		ID:                  owner.ID,
		RefConditionRangeID: owner.RefConditionRangeID,
		OwnerName:           owner.Name,
		ConditionList:       GetCondition(paramsConditions),
	}
	c.JSON(http.StatusAccepted, responses.SuccessResponse(data, http.StatusAccepted))
}

// @summary
// @description
// @tags setting
// @id delete_owner_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your owner id"
// @response 204 "No Content"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/owners/{id} [delete]
func (sh *SettingHandler) DeleteOwnerByID(c *gin.Context) {
	err := sh.settingUseCase.DeleteOwnerByID(c.Param("id"))
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.Status(http.StatusNoContent)
}

// @summary
// @description
// @tags setting
// @id get_owners_road_line
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param page query string false "Insert your page number"
// @param limit query string false "Insert your limit number"
// @param name query string false "Insert your department name"
// @param ref_reflectivity_range_id query string false "Insert your ref reflectivity range id "
// @response 200 {object} responses.Success{data=responses.Pagination{items=[]responses.FirstGroup}} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/owners_road_line [get]
func (sh *SettingHandler) GetOwnersRoadLine(c *gin.Context) {
	var params requests.QueryParamsReflectivityRange
	if err := c.ShouldBind(&params); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	resp, err := sh.settingUseCase.GetOwnersRoadLine(params)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags setting
// @id get_owner_road_line_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your owner id"
// @response 200 {object} responses.Success{data=models.RefOwner} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/owners_road_line/{id} [get]
func (sh *SettingHandler) GetOwnersRoadLineByID(c *gin.Context) {
	ownerID, _ := strconv.Atoi(c.Param("id"))
	_, err := sh.settingUseCase.GetOwnerRoadLineByID(ownerID)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	data, err := sh.settingUseCase.GetRoadLineList(ownerID)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	// owner := GetOwner(paramsConditions)
	// data := responses.ConditionList{
	// 	ID:        owner.ID,
	// 	Name:      owner.Name,
	// 	Condition: GetCondition(paramsConditions),
	// }
	c.JSON(http.StatusOK, responses.SuccessResponse(data, http.StatusOK))
}

// @summary
// @description
// @tags setting
// @id create_owner_road_line
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param OwnerData body requests.OwnerRoadLineRequest false "Insert your owner data"
// @response 201 {object} responses.CreateResponse "Create Successfully"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/owners_road_line [post]
func (sh *SettingHandler) CreateOwnerRoadLine(c *gin.Context) {
	// request
	var request requests.OwnerRoadLineRequest
	if err := c.ShouldBind(&request); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	// validate
	if validateErr := validator.Validate(request); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	// create owner
	ownerId, err := sh.settingUseCase.CreateOwnerRoadLine(request)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok && err.StatusCode == 422 {
			err := helpers.ConverstError(err)
			errResponse := responses.ValidateResponse(err)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	// create Condition
	// conditionReq := request.ConditionList
	_, err = sh.settingUseCase.CreateRoadLineList(ownerId, request)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	data, err := sh.settingUseCase.GetRoadLineList(ownerId)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	// owner := GetOwner(paramsConditions)
	// data := responses.ConditionList{
	// 	ID:        owner.ID,
	// 	Name:      owner.Name,
	// 	Condition: GetCondition(paramsConditions),
	// }
	c.JSON(201, responses.SuccessResponse(data, http.StatusOK))
}

// @summary
// @description
// @tags setting
// @id update_owner_road_line_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your owner id"
// @param OwnerData body requests.OwnerRoadLineRequest false "Insert your owner data"
// @response 202 {object} responses.UpdateResponse "Update Successfully"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/owners_road_line/{id} [put]
func (sh *SettingHandler) UpdateOwnerRoadLineByID(c *gin.Context) {
	// request
	ownerID, _ := strconv.Atoi(c.Param("id"))
	var request requests.OwnerRoadLineRequest
	if err := c.ShouldBind(&request); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	// validate
	if validateErr := validator.Validate(request); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	// create owner
	err := sh.settingUseCase.UpdateOwnerRoadLineByID(ownerID, request)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok && err.StatusCode == 422 {
			err := helpers.ConverstError(err)
			errResponse := responses.ValidateResponse(err)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	// // create Condition
	// conditionReq := request.ConditionList
	err = sh.settingUseCase.UpdateRoadLineList(ownerID, request)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	data, err := sh.settingUseCase.GetRoadLineList(ownerID)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusAccepted, responses.SuccessResponse(data, http.StatusAccepted))
}

// @summary
// @description
// @tags setting
// @id delete_owner_road_line_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your owner id"
// @response 204 "No Content"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/owners_road_line/{id} [delete]
func (sh *SettingHandler) DeleteOwnerRoadLineByID(c *gin.Context) {
	err := sh.settingUseCase.DeleteOwnerRoadLineByID(c.Param("id"))
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.Status(http.StatusNoContent)
}

func GetConditionType(condition models.ParamsConditionPreload) responses.ConditionType {
	return responses.ConditionType{
		Grade:            condition.RefGrade,
		LeftValueAC:      condition.LeftValueAC,
		LeftConditionAC:  condition.LeftConditionAC,
		RightValueAC:     condition.RightValueAC,
		RightConditionAC: condition.RightConditionAC,
		LeftValueCC:      condition.LeftValueCC,
		LeftConditionCC:  condition.LeftConditionCC,
		RightValueCC:     condition.RightValueCC,
		RightConditionCC: condition.RightConditionCC,
	}
}

func GetOwner(paramsCondition []models.ParamsConditionPreload) models.RefOwnerPreload {
	if len(paramsCondition) == 0 {
		return models.RefOwnerPreload{}
	}

	return paramsCondition[0].RefOwner
}

func GetCondition(paramsCondition []models.ParamsConditionPreload) []responses.ConditionListNew {
	ConditionListNews := []responses.ConditionListNew{}
	if len(paramsCondition) == 0 {
		return ConditionListNews
	}
	condition := responses.Condition{}
	for _, v := range paramsCondition {
		switch v.ConditionType {
		case "IFI":
			ifi := GetConditionType(v)
			condition.IFI = append(condition.IFI, ifi)
		case "IRI":
			iri := GetConditionType(v)
			condition.IRI = append(condition.IRI, iri)
		case "MPD":
			mpd := GetConditionType(v)
			condition.MPD = append(condition.MPD, mpd)
		case "RUT":
			rut := GetConditionType(v)
			condition.RUT = append(condition.RUT, rut)
		default:
			log.Println("there is no condition type")
		}
	}
	if len(condition.IFI) != 0 {
		ConditionListNew := GenRespond(condition, "IFI")
		ConditionListNews = append(ConditionListNews, ConditionListNew)
	}
	if len(condition.IRI) != 0 {
		ConditionListNew := GenRespond(condition, "IRI")
		ConditionListNews = append(ConditionListNews, ConditionListNew)
	}
	if len(condition.MPD) != 0 {
		ConditionListNew := GenRespond(condition, "MPD")
		ConditionListNews = append(ConditionListNews, ConditionListNew)
	}
	if len(condition.RUT) != 0 {
		ConditionListNew := GenRespond(condition, "RUT")
		ConditionListNews = append(ConditionListNews, ConditionListNew)
	}
	return ConditionListNews
}

func GenRespond(condition responses.Condition, conditionType string) responses.ConditionListNew {
	ConditionListNew := responses.ConditionListNew{ConditionType: conditionType}
	SurfaceTypeCondition := responses.SurfaceTypeCondition{}
	AC := []responses.ConditionTypeAC{}
	CC := []responses.ConditionTypeCC{}
	switch conditionType {
	case "IFI":
		for _, v := range condition.IFI {
			ac := responses.ConditionTypeAC{Grade: v.Grade,
				LeftValueAC:      v.LeftValueAC,
				LeftConditionAC:  v.LeftConditionAC,
				RightValueAC:     v.RightValueAC,
				RightConditionAC: v.RightConditionAC}
			cc := responses.ConditionTypeCC{Grade: v.Grade,
				LeftValueCC:      v.LeftValueCC,
				LeftConditionCC:  v.LeftConditionCC,
				RightValueCC:     v.RightValueCC,
				RightConditionCC: v.RightConditionCC}
			AC = append(AC, ac)
			CC = append(CC, cc)

		}
	case "IRI":
		for _, v := range condition.IRI {
			ac := responses.ConditionTypeAC{Grade: v.Grade,
				LeftValueAC:      v.LeftValueAC,
				LeftConditionAC:  v.LeftConditionAC,
				RightValueAC:     v.RightValueAC,
				RightConditionAC: v.RightConditionAC}
			cc := responses.ConditionTypeCC{Grade: v.Grade,
				LeftValueCC:      v.LeftValueCC,
				LeftConditionCC:  v.LeftConditionCC,
				RightValueCC:     v.RightValueCC,
				RightConditionCC: v.RightConditionCC}
			AC = append(AC, ac)
			CC = append(CC, cc)

		}
	case "MPD":
		for _, v := range condition.MPD {
			ac := responses.ConditionTypeAC{Grade: v.Grade,
				LeftValueAC:      v.LeftValueAC,
				LeftConditionAC:  v.LeftConditionAC,
				RightValueAC:     v.RightValueAC,
				RightConditionAC: v.RightConditionAC}
			cc := responses.ConditionTypeCC{Grade: v.Grade,
				LeftValueCC:      v.LeftValueCC,
				LeftConditionCC:  v.LeftConditionCC,
				RightValueCC:     v.RightValueCC,
				RightConditionCC: v.RightConditionCC}
			AC = append(AC, ac)
			CC = append(CC, cc)

		}
	case "RUT":
		for _, v := range condition.RUT {
			ac := responses.ConditionTypeAC{Grade: v.Grade,
				LeftValueAC:      v.LeftValueAC,
				LeftConditionAC:  v.LeftConditionAC,
				RightValueAC:     v.RightValueAC,
				RightConditionAC: v.RightConditionAC}
			cc := responses.ConditionTypeCC{Grade: v.Grade,
				LeftValueCC:      v.LeftValueCC,
				LeftConditionCC:  v.LeftConditionCC,
				RightValueCC:     v.RightValueCC,
				RightConditionCC: v.RightConditionCC}
			AC = append(AC, ac)
			CC = append(CC, cc)

		}
	}

	SurfaceTypeCondition.AC = AC
	SurfaceTypeCondition.CC = CC
	ConditionListNew.SurfaceType = SurfaceTypeCondition
	return ConditionListNew
}

// @summary
// @description
// @tags setting
// @id get_signs
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param page query string false "Insert your page number"
// @param limit query string false "Insert your limit number"
// @param name query string false "Insert your sign name"
// @response 200 {object} responses.Success{data=responses.Pagination{items=[]models.RefAssetSignImage}} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/signs [get]
func (sh *SettingHandler) GetSigns(c *gin.Context) {
	var params requests.QueryParams
	if err := c.ShouldBind(&params); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	resp, err := sh.settingUseCase.GetSigns(c, params)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags setting
// @id get_sign_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your asset ID"
// @response 200 {object} responses.Success{data=models.RefAssetSignImage} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/signs/{id} [get]
func (sh *SettingHandler) GetSignByID(c *gin.Context) {
	resp, err := sh.settingUseCase.GetSignByID(c, c.Param("id"))
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags setting
// @id create_sign
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param name formData string false "Insert your sign name"
// @param abbr formData string false "Insert your sign abbreviation"
// @param image formData file false "upload your sign image"
// @param image_status formData string false "sign image status"
// @response 201 {object} responses.CreateResponse "Create Successfully"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/signs [post]
func (sh *SettingHandler) CreateSign(c *gin.Context) {
	var request requests.SignImageRequest

	// no handle error because user can send
	// empty file from multipart/form-data
	c.ShouldBind(&request)

	errMessageSlice := []string{}
	if request.Image == nil {
		errMessageSlice = append(errMessageSlice, "Image:zero value")
	}

	if request.Name == "" {
		errMessageSlice = append(errMessageSlice, "Name:zero value")
	}

	if request.Abbr == "" {
		errMessageSlice = append(errMessageSlice, "Abbr:zero value")
	}

	if err := helpers.ValidateRequest(errMessageSlice); err != nil {
		validateErr := helpers.ConverstError(err)
		errResponse := responses.ValidateResponse(validateErr)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	err := sh.settingUseCase.CreateSign(c, request)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(responses.NoData{}, http.StatusCreated))
}

// @summary
// @description
// @tags setting
// @id update_sign_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your sign image ID"
// @param name formData string false "Insert your sign name"
// @param abbr formData string false "Insert your sign abbreviation"
// @param image formData file false "upload your sign image"
// @param image_status formData string false "sign image status"
// @response 202 {object} responses.UpdateResponse "Update Successfully"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/signs/{id} [put]
func (sh *SettingHandler) UpdateSignByID(c *gin.Context) {
	var request requests.SignImageRequest

	// no handle error because user can send
	// empty file from multipart/form-data
	c.ShouldBind(&request)

	errMessageSlice := []string{}
	if request.Name == "" {
		errMessageSlice = append(errMessageSlice, "Name:zero value")
	}

	if request.Abbr == "" {
		errMessageSlice = append(errMessageSlice, "Abbr:zero value")
	}

	if request.ImageStatus == "delete" || request.ImageStatus == "no_file" {
		errMessageSlice = append(errMessageSlice, "Image:zero value")
	}

	if err := helpers.ValidateRequest(errMessageSlice); err != nil {
		validateErr := helpers.ConverstError(err)
		errResponse := responses.ValidateResponse(validateErr)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	err := sh.settingUseCase.UpdateSignByID(c, request)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusAccepted, responses.SuccessResponse(responses.NoData{}, http.StatusAccepted))
}

// @summary
// @description
// @tags setting
// @id delete_sign_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your asset ID"
// @response 204 "No Content"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/signs/{id} [delete]
func (sh *SettingHandler) DeleteSignByID(c *gin.Context) {
	err := sh.settingUseCase.DeleteSignByID(c.Param("id"))
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// @summary
// @description
// @tags setting
// @id get_asset_tables
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param page query string false "Insert your page number"
// @param limit query string false "Insert your limit number"
// @param asset_type query string false "Insert your asset type, please choose which one between 'in' or 'out'"
// @param name query string false "Insert your asset table name"
// @param group_id query string false "Insert your asset group id"
// @response 200 {object} responses.Success{data=responses.Pagination{items=[]responses.AssetTableResponse}} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/asset_tables [get]
func (sh *SettingHandler) GetAssetTables(c *gin.Context) {
	var params requests.AssetTableQueryParams
	if err := c.ShouldBind(&params); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	assetTablesData, err := sh.settingUseCase.GetAssetTables(params)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	assetTables := assetTablesData.AssetTables
	assetTableStaffs := assetTablesData.AssetTableStaffs

	items := []responses.AssetTableResponse{}
	for _, assetTable := range assetTables {
		respItem := responses.AssetTableResponse{
			ID:         assetTable.ID,
			TableName:  assetTable.TableNameColumn,
			TableLabel: assetTable.TableLabel,
			AssetGroup: assetTable.RefAsset.Name,
			CanDelete:  assetTable.CanDelete,
		}

		items = append(items, respItem)
	}

	responsibleDept := make(map[int]string)

	//cathering responsible department for any asset table
	for _, staff := range assetTableStaffs {
		responsibleDept[staff.RefAssetTableID] += staff.RefDepartment.Name + ", "
	}

	//remove "," and " " on the end of responsible department line
	for idx, admins := range responsibleDept {
		responsibleDept[idx] = admins[:len(admins)-2]
	}

	//set responsible department to items
	for i := 0; i < len(items); i++ {
		items[i].ResponsibleDept = responsibleDept[items[i].ID]
	}

	resp := helpers.Pagination(items, assetTablesData.Limit, assetTablesData.Page, assetTablesData.Total)

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags setting
// @id get_asset_table_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert your asset table ID"
// @response 200 {object} responses.Success{data=responses.AssetTableDetailResponse} "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/asset_tables/{id} [get]
func (sh *SettingHandler) GetAssetTableByID(c *gin.Context) {
	assetTableData, err := sh.settingUseCase.GetAssetTableByID(c, c.Param("id"))
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	assetTable := assetTableData.AssetTable
	//assetTableStaffs := assetTableData.AssetTableStaffs
	iconFilepath := ""
	if assetTable.IconFilepath != "" {
		iconFilepath = os.Getenv("STORAGE_IP") + "/" + assetTable.IconFilepath
	} else {
		iconFilepath = ""
	}
	assetTableDetailResponse := responses.AssetTableDetailResponse{
		AssetID:       assetTable.ID,
		TableName:     assetTable.TableNameColumn,
		TableLabel:    assetTable.TableLabel,
		GeomType:      assetTable.GeomType,
		IconFilePath:  iconFilepath,
		LineColorCode: assetTable.LineColorCode,
		AssetGroup:    assetTable.RefAsset,
	}

	// for _, staff := range assetTableStaffs {
	// 	switch staff.IsApprover {
	// 	case false:
	// 		assetTableDetailResponse.Viewer = append(assetTableDetailResponse.Viewer, staff.RefDepartment)
	// 	case true:
	// 		assetTableDetailResponse.Approver = append(assetTableDetailResponse.Approver, staff.RefDepartment)
	// 	default:
	// 		log.Println("is_approver field type must be a boolean")
	// 	}
	// }

	for _, column := range assetTable.AssetTableColumns {
		column := responses.Columns{
			ColumnID:        column.ID,
			ColumnName:      column.ColumnName,
			TableNameRef:    column.TableNameRef,
			ColumnDataType:  column.ColumnDataType,
			ComponentTitle:  column.ComponentTitle,
			ComponentType:   column.ComponentType,
			IsRequired:      column.IsRequired,
			IsVisibleView:   column.IsVisibleView,
			IsVisibleEdit:   column.IsVisibleEdit,
			IsMandatory:     column.IsMandatory,
			IsVisibleReport: column.IsVisibleReport,
		}

		assetTableDetailResponse.Columns = append(assetTableDetailResponse.Columns, column)
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(assetTableDetailResponse, http.StatusOK))
}

// @param AssetTable body requests.SwaggerAssetTable true "NOTE: please use postman to test this api endpoint and use form-data to insert data"

// @summary
// @description
// @tags setting
// @id create_asset_table
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param icon_filepath formData file false "upload your sign image"
// @param icon_filepath_status formData string false "icon_filepath_status"
// @param data formData string false "data"
// @response 201 {object} responses.CreateResponse "Create Successfully"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/asset_tables [post]
func (sh *SettingHandler) CreateAssetTable(c *gin.Context) {
	var request requests.AssetTableRequest

	// no handle error because user can send
	// empty file from multipart/form-data
	c.ShouldBind(&request)

	data := request.Data

	errMessageSlice := []string{}
	if data.TableName == "" {
		errMessageSlice = append(errMessageSlice, "TableName:zero value")
	}

	if data.TableLabel == "" {
		errMessageSlice = append(errMessageSlice, "TableLabel:zero value")
	}

	if data.AssetType == "" {
		errMessageSlice = append(errMessageSlice, "AssetType:zero value")
	}

	if data.AssetGroup <= 0 {
		errMessageSlice = append(errMessageSlice, "AssetGroup:zero value")
	}

	if data.GeomType == "" {
		errMessageSlice = append(errMessageSlice, "GeomType:zero value")
	}

	if err := helpers.ValidateRequest(errMessageSlice); err != nil {
		validateErr := helpers.ConverstError(err)
		errResponse := responses.ValidateResponse(validateErr)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	validateMsgArr := map[string]string{}
	pattern := "^[a-zA-Z0-9_]+$"
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(data.TableName) {
		validateMsgArr["table_name"] = "กรอกได้เฉพาะ a-z, A-Z, 0-9"
		errResponse := responses.ValidateResponse(validateMsgArr)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	if strings.HasPrefix(data.TableName, "ref_asset") {
		validateMsgArr["table_name"] = "โปรดตั้งชื่อตารางข้อมูลที่ไม่ขึ้นต้นด้วยคำว่า ref_asset"
		errResponse := responses.ValidateResponse(validateMsgArr)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	err := sh.settingUseCase.CreateAssetTable(data, request.IconFilePath, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok && err.StatusCode == 422 {
			err := helpers.ConverstError(err)
			errResponse := responses.ValidateResponse(err)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(responses.NoData{}, http.StatusCreated))
}

// @param AssetTable body requests.SwaggerAssetTable true "NOTE: please use postman to test this api endpoint and use form-data to insert data"

// @summary
// @description
// @tags setting
// @id update_asset_table_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your asset table ID"
// @param icon_filepath formData file false "upload your sign image"
// @param icon_filepath_status formData string false "icon_filepath_status"
// @param data formData string false "data"
// @response 202 {object} responses.UpdateResponse "Update Successfully"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/asset_tables/{id} [put]
func (sh *SettingHandler) UpdateAssetTableByID(c *gin.Context) {
	var request requests.AssetTableRequest

	// no handle error because user can send
	// empty file from multipart/form-data
	c.ShouldBind(&request)

	data := request.Data
	errMessageSlice := []string{}
	if data.TableLabel == "" {
		errMessageSlice = append(errMessageSlice, "TableLabel:zero value")
	}

	if data.AssetType == "" {
		errMessageSlice = append(errMessageSlice, "AssetType:zero value")
	}

	if data.AssetGroup <= 0 {
		errMessageSlice = append(errMessageSlice, "AssetGroup:zero value")
	}

	if err := helpers.ValidateRequest(errMessageSlice); err != nil {
		validateErr := helpers.ConverstError(err)
		errResponse := responses.ValidateResponse(validateErr)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	err := sh.settingUseCase.UpdateAssetTableByID(data, request.IconFilePath, request.IconFilePathStatus, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusAccepted, responses.SuccessResponse(responses.NoData{}, http.StatusAccepted))
}

// @summary
// @description
// @tags setting
// @id delete_asset_table_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your asset table ID"
// @response 204 "No Content"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/asset_tables/{id} [delete]
func (sh *SettingHandler) DeleteAssetTableByID(c *gin.Context) {
	err := sh.settingUseCase.DeleteAssetTableByID(c.Param("id"))
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// @summary
// @description
// @tags setting_models
// @id delete_intervention_criteria
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your ID"
// @response 204 "No Content"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/intervention_criteria/{id} [delete]
func (sh *SettingHandler) DeleteInterventionCriteria(c *gin.Context) {
	err := sh.settingUseCase.DeleteInterventionCriteriaById(c.Param("id"), c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// @summary
// @description
// @tags setting_models
// @id create_intervention_criteria
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param InterventionCriteria body requests.InterventionCriteria true "Insert your json"
// @response 200 {object} responses.Success{data=requests.InterventionCriteria} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/intervention_criteria [post]
func (sh *SettingHandler) CreateInterventionCriteria(c *gin.Context) {
	var request requests.InterventionCriteria
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

	interventionCriteria, err := sh.settingUseCase.CreateInterventionCriteria(request, c)
	if err != nil {
		if err.Error() == "standard_name" {
			err := helpers.ConverstError(errors.New("standard_name : already used"))
			errResponse := responses.ValidateResponse(err)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}

		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(interventionCriteria, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id get_intervention_criteria
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.InterventionCriteria} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/intervention_criteria [get]
func (sh *SettingHandler) GetInterventionCriteria(c *gin.Context) {
	interventionCriteria, err := sh.settingUseCase.GetInterventionCriteria(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(interventionCriteria, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id GetInterventionCriteriaList
// @Security Bearer
// @response 200 {object} responses.Success{data=responses.InterventionCriteriaSequenceCriteriaMethod} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/intervention_criteria/list [get]
func (sh *SettingHandler) GetInterventionCriteriaList(c *gin.Context) {
	interventionCriteria, err := sh.settingUseCase.GetInterventionCriteriaMethod(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(interventionCriteria, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id get_intervention_criteria_sequence
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.InterventionCriteriaSequenceCriteriaMethod} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/intervention_criteria/sequence [get]
func (sh *SettingHandler) GetInterventionCriteriaSequence(c *gin.Context) {
	interventionCriteria, err := sh.settingUseCase.GetInterventionCriteriaSequence(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(interventionCriteria, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_intervention_criteria_sequence
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param InterventionCriteriaSequenceCriteriaMethod body requests.InterventionCriteriaSequenceCriteriaMethod true "Insert your json"
// @response 200 {object} responses.Success{data=requests.InterventionCriteriaSequenceCriteriaMethod} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/intervention_criteria/sequence [post]
func (sh *SettingHandler) CreateInterventionCriteriaSequence(c *gin.Context) {
	var request requests.InterventionCriteriaSequenceCriteriaMethod
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

	interventionCriteria, err := sh.settingUseCase.CreateInterventionCriteriaSequence(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(interventionCriteria, http.StatusCreated))
}

// @summary Create a new ref surface
// @description This endpoint allows you to create a new ref surface with the given payload
// @tags setting
// @id post_ref_surface
// @accept json
// @produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param requestBody body requests.RefSurface true "Request body to create a new ref surface"
// @response 201 {object} responses.CreateResponse "Create Successfully"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/ref/surface [post]
func (sh *SettingHandler) PostRefSurface(c *gin.Context) {
	var requestFirst requests.RefSurfacePointer
	var request requests.RefSurface
	if err := c.ShouldBind(&requestFirst); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	// var isC2nil bool
	// if requestFirst.C2 == nil {
	// 	isC2nil = true
	// 	num := 1
	// 	requestFirst.C2 = &num
	// }
	if strings.ToLower(requestFirst.SurfaceGroup) != "asphalt" {
		if validateErr := helpers.CheckFloatIntNonNil(requestFirst, "CRT", "RRF"); validateErr != nil {
			err := helpers.ConverstError(validateErr)
			errResponse := responses.ValidateResponse(err)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}
	} else {
		if validateErr := helpers.CheckFloatIntNonNil(requestFirst); validateErr != nil {
			err := helpers.ConverstError(validateErr)
			errResponse := responses.ValidateResponse(err)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}
	}

	if validateErr := validator.Validate(requestFirst); validateErr != nil {
		err := helpers.ConverstError(validateErr)
		errResponse := responses.ValidateResponse(err)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	copier.Copy(&request, &requestFirst)
	helpers.CopyPointerToValueFloatInt(&request, &requestFirst)

	if request.C1 == 0 && request.C2 == 0 {
		err := errors.New("Input number(0^0) cannot calculate")
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	if strings.ToLower(request.SurfaceGroup) != "concrete" && strings.ToLower(request.SurfaceGroup) != "asphalt" {
		err := errors.New("Surface Group ไม่ถูกต้อง")
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	userIdValue, exists := c.Get("userID")
	if !exists {
		err := errors.New("userId not found")
		logs.Error(err)
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	userId, ok := userIdValue.(float64)
	if !ok {
		err := errors.New("Error to find UserID")
		logs.Error(err)
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	err := sh.settingUseCase.PostRefSurface(request, int(userId), false)
	if err != nil {
		errRespond := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errRespond)
		return
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(responses.NoData{}, http.StatusCreated))

}

// @summary Update a new ref surface
// @description This endpoint allows you to update ref surface with the given payload
// @tags setting
// @id put_ref_surface
// @accept json
// @produce json
// @Param id path int true "ID of the reference surface"
// @param requestBody body requests.RefSurface true "Request body to update a new ref surface"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 202 {object} responses.UpdateResponse "Update Successfully"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/ref/surface/{id} [put]
func (sh *SettingHandler) PutRefSurface(c *gin.Context) {
	var requestFirst requests.RefSurfacePointer
	var request requests.RefSurface
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		logs.Error(err)
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	if err := c.ShouldBind(&requestFirst); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	validateMsgArr := map[string]string{}

	if strings.ToLower(requestFirst.SurfaceGroup) != "asphalt" {
		if validateErr := helpers.CheckFloatIntNonNil(requestFirst, "CRT", "RRF"); validateErr != nil {
			errs := helpers.ConverstError(validateErr)
			for k, v := range errs {
				validateMsgArr[k] = v
			}
		}
	} else {
		if validateErr := helpers.CheckFloatIntNonNil(requestFirst); validateErr != nil {
			errs := helpers.ConverstError(validateErr)
			for k, v := range errs {
				validateMsgArr[k] = v
			}
		}
	}
	if validateErr := validator.Validate(requestFirst); validateErr != nil {
		errs := helpers.ConverstError(validateErr)
		for k, v := range errs {
			validateMsgArr[k] = v
		}
	}

	if requestFirst.Drainage == nil {
		validateMsgArr["drainage"] = " โปรดระบุ"
	}

	if requestFirst.LayerCoefficient == nil {
		validateMsgArr["layer_coefficient"] = " โปรดระบุ"
	}

	if requestFirst.A == nil {
		validateMsgArr["a"] = " โปรดระบุ"
	}

	if requestFirst.B == nil {
		validateMsgArr["b"] = " โปรดระบุ"
	}

	if requestFirst.C1 == nil {
		validateMsgArr["c1"] = " โปรดระบุ"
	}

	if requestFirst.C2 == nil {
		validateMsgArr["c2"] = " โปรดระบุ"
	}

	if requestFirst.CRT == nil {
		validateMsgArr["crt"] = " โปรดระบุ"
	}

	if requestFirst.RRF == nil {
		validateMsgArr["rrf"] = " โปรดระบุ"
	}

	if len(validateMsgArr) > 0 {
		errResponse := responses.ValidateResponse(validateMsgArr)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}
	// var isC2nil bool
	// if requestFirst.C2 == nil {
	// 	isC2nil = true
	// 	num := 1
	// 	requestFirst.C2 = &num
	// }

	copier.Copy(&request, &requestFirst)
	helpers.CopyPointerToValueFloatInt(&request, &requestFirst)

	if request.C1 == 0 && request.C2 == 0 {
		err = errors.New("Input number(0^0) cannot calculate")
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	if strings.ToLower(request.SurfaceGroup) != "concrete" && strings.ToLower(request.SurfaceGroup) != "asphalt" {
		err := errors.New("Surface Group ไม่ถูกต้อง")
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	userIdValue, exists := c.Get("userID")
	if !exists {
		err := errors.New("userId not found")
		logs.Error(err)
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	userId, ok := userIdValue.(float64)
	if !ok {
		err := errors.New("Error to find UserID")
		logs.Error(err)
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	err = sh.settingUseCase.PutRefSurface(request, int(userId), id, false)
	if err != nil {
		errRespond := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errRespond)
		return
	}
	c.JSON(http.StatusAccepted, responses.SuccessResponse(responses.NoData{}, http.StatusAccepted))
}

// @Summary
// @Description
// @Tags setting
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your ref surface ID"
// @response 204 "No Content"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorize
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @Router /api/v1/settings/ref/surface/{id} [delete]
func (m *SettingHandler) DeleteSettingRefSurfaceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	err = m.settingUseCase.DeleteSettingRefSurfaceByID(id)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

// @summary Retrieve all ref surfaces
// @description This endpoint allows you to retrieve all reference surfaces
// @tags setting
// @id get_ref_surfaces
// @accept json
// @produce json
// @Param name query string false "name of the reference surface"
// @Param type query string false "type of the reference surface"
// @param page query string false "Insert your page number"
// @param limit query string false "Insert your limit number"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} []models.NewRefSurface "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @router /api/v1/settings/ref/surface [get]
func (sh *SettingHandler) GetRefSurface(c *gin.Context) {
	surfaceType := c.Query("type")
	surfaceName := c.Query("name")
	limitParam := c.Query("limit")
	pageParam := c.Query("page")
	whereCondition := ""
	if surfaceName != "" && surfaceType != "" {
		whereCondition = "name ILIKE '%" + surfaceName + "%' AND LOWER(type) = '" + strings.ToLower(surfaceType) + "'"
	} else if surfaceName != "" {
		whereCondition = "name ILIKE '%" + surfaceName + "%'"
	} else if surfaceType != "" {
		whereCondition = fmt.Sprintf("LOWER(type) = '%s' ", strings.ToLower(surfaceType))
	}
	responds, err := sh.settingUseCase.GetRefSurface(whereCondition)
	if err != nil {
		errRespond := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errRespond)
		return
	}

	totalItems := int64(len(responds))
	limit, offset, page := helpers.GetlimitOffsetPage(limitParam, pageParam, totalItems)
	if totalItems == 0 {
		responds = []models.NewRefSurface{}
	} else if limit+offset > totalItems {
		responds = responds[offset:totalItems]
	} else {
		responds = responds[offset : limit+offset]
	}
	pagination := helpers.Pagination(responds, limit, page, totalItems)

	c.JSON(http.StatusOK, responses.SuccessResponse(pagination, http.StatusOK))
}

// @summary Retrieve a ref surface by ID
// @description This endpoint allows you to retrieve a reference surface by its ID
// @tags setting
// @id get_ref_surface_by_id
// @accept json
// @produce json
// @Param id path int true "ID of the reference surface"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} models.NewRefSurface "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @router /api/v1/settings/ref/surface/{id} [get]
func (sh *SettingHandler) GetRefSurfaceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		errRespond := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errRespond)
		return
	}
	respond, err := sh.settingUseCase.GetRefSurfaceByID(id)
	if err != nil {
		errRespond := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errRespond)
		return
	}
	if respond == nil {
		respond = responses.NoData{}
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(respond, http.StatusOK))
}

// @summary Retrieve a ref surface with an optional ID parameter
// @description This endpoint allows you to retrieve a reference surface with an optional ID parameter. If ID parameter is not provided, it retrieves all reference surfaces.
// @tags setting
// @id get_param_ref_surface
// @accept json
// @produce json
// @Param id query int false "ID of the reference surface (optional)"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} []models.RefSurfaceParam "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 422 {object} responses.Validate "Unprocessable Entity"
// @router /api/v1/settings/ref/surface_param [get]
func (sh *SettingHandler) GetParamRefSurface(c *gin.Context) {
	var id int
	var err error
	idParam := c.Query("id")
	if idParam != "" {
		id, err = strconv.Atoi(idParam)
		if err != nil {
			errRespond := responses.FailRespone(err)
			c.JSON(http.StatusBadRequest, errRespond)
			return
		}
	}

	respond, err := sh.settingUseCase.GetParamRefSurface(id)
	if err != nil {
		errRespond := responses.FailRespone(err)
		c.JSON(http.StatusUnprocessableEntity, errRespond)
		return
	}
	if respond == nil {
		c.JSON(http.StatusOK, responses.SuccessResponse([]string{}, http.StatusOK))
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(respond, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_budget
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param SettingBudget body requests.Budget true "Insert your json"
// @response 201 {object} responses.Success{data=requests.Budget} "CREATED"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/budget [post]
func (sh *SettingHandler) CreateBudget(c *gin.Context) {
	var request requests.Budget
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

	resp, err := sh.settingUseCase.CreateBudget(request, c)
	if err != nil {
		if err.Error() == "name" {
			err := helpers.ConverstError(errors.New("name : already used"))
			errResponse := responses.ValidateResponse(err)
			c.JSON(http.StatusUnprocessableEntity, errResponse)
			return
		}

		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(resp, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id get_budget
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param page query string false "Insert your page number"
// @param limit query string false "Insert your limit number"
// @param name query string false "Insert your asset name"
// @response 200 {object} responses.Success{data=responses.Pagination{items=[]responses.BudgetList}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/budget [get]
func (sh *SettingHandler) GetBudget(c *gin.Context) {
	var params requests.QueryParams
	if err := c.ShouldBind(&params); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	resp, err := sh.settingUseCase.GetBudget(params, c)
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
// @description
// @tags setting_models
// @id get_budget_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert budget id"
// @response 200 {object} responses.Success{data=responses.Budget} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/budget/{id} [get]
func (sh *SettingHandler) GetBudgetById(c *gin.Context) {
	resp, err := sh.settingUseCase.GetBudgetById(c.Param("id"), c)
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
// @description
// @tags setting_models
// @id update_budget
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param UpdateBudget body requests.UpdateBudget true "Insert your json"
// @response 202 {object} responses.Success{data=requests.UpdateBudget} "ACCEPTED"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/budget [put]
func (sh *SettingHandler) UpdateBudget(c *gin.Context) {
	var request requests.UpdateBudget
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

	resp, err := sh.settingUseCase.UpdateBudget(request, c)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusAccepted, responses.SuccessResponse(resp, http.StatusAccepted))
}

// @summary
// @description
// @tags setting_models
// @id delete_budget
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Insert budget id"
// @response 204 "No Content"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/budget/{id} [delete]
func (sh *SettingHandler) DeleteBudget(c *gin.Context) {
	_, err := sh.settingUseCase.DeleteBudgetById(c.Param("id"), c)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.Status(204)
}

// @summary
// @description
// @tags setting_models
// @id create_aadt_growth_rate
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param CreateAadtGrowthRate body []requests.CreateAadtGrowthRate true "Insert your growth rate json"
// @response 200 {object} responses.Success{data=[]requests.CreateAadtGrowthRate} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/aadt/growth_rate [post]
func (sh *SettingHandler) CreateAadtGrowthRate(c *gin.Context) {
	var request []requests.CreateAadtGrowthRate
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

	resp, err := sh.settingUseCase.CreateAadtGrowthRate(request, c)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(resp, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id get_setting_aadt_growth_rate
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=[]responses.AadtGrowthRate} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/aadt/growth_rate [get]
func (sh *SettingHandler) GetAadtGrowthRate(c *gin.Context) {

	resp, err := sh.settingUseCase.GetAadtGrowthRate(c)
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
// @description
// @tags setting_models
// @id create_setting_aadt_percentage_vehicle_type
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param AadtPercentageVehicleType body requests.AadtPercentageVehicleType true "Insert your json"
// @response 201 {object} responses.Success{data=requests.AadtPercentageVehicleType} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/aadt/percentage_vehicle_type [post]
func (sh *SettingHandler) CreateAadtPercentageVehicleType(c *gin.Context) {
	var request requests.AadtPercentageVehicleType
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

	resp, err := sh.settingUseCase.CreateAadtPercentageVehicleType(request, c)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(resp, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id get_setting_aadt_percentage_vehicle_type
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param road_group_id path string false "Insert your road group id"
// @response 200 {object} responses.Success{data=responses.AadtPercentageVehicleType} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/aadt/percentage_vehicle_type/{road_group_id} [get]
func (sh *SettingHandler) GetAadtPercentageVehicleType(c *gin.Context) {
	resp, err := sh.settingUseCase.GetAadtPercentageVehicleType(c, c.Param("road_group_id"))
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
// @description
// @tags setting_models
// @id create_setting_aadt_parameter
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param CreateAadtParameter body requests.CreateAadtParameter true "Insert your parameter json"
// @response 201 {object} responses.Success{data=models.AadtParameter} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/aadt/parameter [post]
func (sh *SettingHandler) CreateAadtParameter(c *gin.Context) {
	var request requests.CreateAadtParameter
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

	resp, err := sh.settingUseCase.CreateAadtParameter(request, c)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse(resp, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id get_setting_aadt_parameter
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param road_group_id path string false "Insert your road group id"
// @response 200 {object} responses.Success{data=interface{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/aadt/parameter/{road_group_id} [get]
func (sh *SettingHandler) GetAadtParameter(c *gin.Context) {
	resp, err := sh.settingUseCase.GetAadtParameter(c, c.Param("road_group_id"))
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
// @description
// @tags setting_models
// @id get_aadt_parameter_road_group_with_volume_aadt
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=[]responses.RoadGroupWithVolumeAadt} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/aadt/parameter/road_group_with_volume_aadt [get]
func (sh *SettingHandler) GetAadtParameterRoadGroupWithVolumeAadt(c *gin.Context) {
	resp, err := sh.settingUseCase.GetAadtParameterRoadGroupWithVolumeAadt(c)
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
// @description
// @tags setting_models
// @id create_road_work_effect_asphalt
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param SettingRoadWorkEffectAsphalt body requests.SettingRoadWorkEffectAsphalt true "Insert your json"
// @response 201 {object} responses.Success{data=responses.SettingRoadWorkEffectAsphalt} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_work_effect/asphalt [post]
func (sh *SettingHandler) CreateRoadWorkEffectAsphalt(c *gin.Context) {
	var request requests.SettingRoadWorkEffectAsphalt
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

	roadWorkEffect, err := sh.settingUseCase.CreateRoadWorkEffectAsphalt(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadWorkEffect, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id create_road_work_effect_concrete
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param SettingRoadWorkEffectConcrete body requests.SettingRoadWorkEffectConcrete true "Insert your json"
// @response 201 {object} responses.Success{data=responses.SettingRoadWorkEffectConcrete} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_work_effect/concrete [post]
func (sh *SettingHandler) CreateRoadWorkEffectConcrete(c *gin.Context) {
	var request requests.SettingRoadWorkEffectConcrete
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

	roadWorkEffect, err := sh.settingUseCase.CreateRoadWorkEffectConcrete(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadWorkEffect, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id get_road_work_effect_concrete
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.SettingRoadWorkEffectConcrete} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_work_effect/concrete [get]
func (sh *SettingHandler) GetRoadWorkEffectConcrete(c *gin.Context) {
	resp, err := sh.settingUseCase.GetRoadWorkEffectConcrete(c)
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
// @description
// @tags setting_models
// @id get_road_work_effect_asphalt
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.SettingRoadWorkEffectAsphalt} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_work_effect/asphalt [get]
func (sh *SettingHandler) GetRoadWorkEffectAsphalt(c *gin.Context) {
	resp, err := sh.settingUseCase.GetRoadWorkEffectAsphalt(c)
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
// @description
// @tags setting_models
// @id create_road_user_cost_acc_loss_value
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param RoadUserCostAccLossValue body requests.RoadUserCostAccLossValue true "Insert your json"
// @response 201 {object} responses.Success{data=responses.RoadUserCostAccLossValue} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/acc/loss_value [post]
func (sh *SettingHandler) CreateRoadUserCostAccLossValue(c *gin.Context) {
	var request requests.RoadUserCostAccLossValue
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

	roadUserCost, err := sh.settingUseCase.CreateRoadUserCostAccLossValue(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadUserCost, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id get_road_user_cost_acc_loss_value
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RoadUserCostAccLossValue} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/acc/loss_value [get]
func (sh *SettingHandler) GetRoadUserCostAccLossValue(c *gin.Context) {
	roadUserCost, err := sh.settingUseCase.GetRoadUserCostAccLossValue(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(roadUserCost, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_acc_chance_of_accident
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param RoadUserCostAccChanceOfAccident body requests.RoadUserCostAccChanceOfAccident true "Insert your json"
// @response 201 {object} responses.Success{data=requests.RoadUserCostAccChanceOfAccident} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/acc/chance_of_accident [post]
func (sh *SettingHandler) CreateRoadUserCostAccChanceOfAccident(c *gin.Context) {
	var request requests.RoadUserCostAccChanceOfAccident
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

	roadUserCost, err := sh.settingUseCase.CreateRoadUserCostAccChanceOfAccident(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadUserCost, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id get_road_user_cost_acc_chance_of_accident
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param road_group_id path string false "Insert your road group id"
// @response 200 {object} responses.Success{data=responses.RoadUserCostAccLossValue} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/acc/chance_of_accident/{road_group_id} [get]
func (sh *SettingHandler) GetRoadUserCostAccChanceOfAccident(c *gin.Context) {
	roadUserCost, err := sh.settingUseCase.GetRoadUserCostAccChanceOfAccident(c.Param("road_group_id"), c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(roadUserCost, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_default_data
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param RoadUserCostRusDefaultData body requests.RoadUserCostRusDefaultData true "Insert your json"
// @response 201 {object} responses.Success{data=requests.RoadUserCostRusDefaultData} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/default_data [post]
func (sh *SettingHandler) CreateRoadUserCostRucDefaultData(c *gin.Context) {
	var request requests.RoadUserCostRusDefaultData
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

	roadUserCost, err := sh.settingUseCase.CreateRoadUserCostRucDefaultData(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadUserCost, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_default_data
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RoadUserCostRusDefaultData} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/default_data [get]
func (sh *SettingHandler) GetRoadUserCostRucDefaultData(c *gin.Context) {
	roadUserCost, err := sh.settingUseCase.GetRoadUserCostRucDefaultData(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(roadUserCost, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_driving
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param RoadUserCostRusDriving body requests.RoadUserCostRusDriving true "Insert your json"
// @response 201 {object} responses.Success{data=requests.RoadUserCostRusDriving} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/driving [post]
func (sh *SettingHandler) CreateRoadUserCostRucDriving(c *gin.Context) {
	var request requests.RoadUserCostRusDriving
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

	roadUserCost, err := sh.settingUseCase.CreateRoadUserCostRucDriving(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadUserCost, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_driving
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RoadUserCostRusDriving} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/driving [get]
func (sh *SettingHandler) GetRoadUserCostRucDriving(c *gin.Context) {
	roadUserCost, err := sh.settingUseCase.GetRoadUserCostRucDriving(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(roadUserCost, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_engine_speed
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param RoadUserCostRusEngineSpeed body requests.RoadUserCostRusEngineSpeed true "Insert your json"
// @response 201 {object} responses.Success{data=requests.RoadUserCostRusEngineSpeed} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/engine_speed [post]
func (sh *SettingHandler) CreateRoadUserCostRucEngineSpeed(c *gin.Context) {
	var request requests.RoadUserCostRusEngineSpeed
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

	roadUserCost, err := sh.settingUseCase.CreateRoadUserCostRucEngineSpeed(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadUserCost, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_engine_speed
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RoadUserCostRusEngineSpeed} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/engine_speed [get]
func (sh *SettingHandler) GetRoadUserCostRucEngineSpeed(c *gin.Context) {
	roadUserCost, err := sh.settingUseCase.GetRoadUserCostRucEngineSpeed(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(roadUserCost, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_fuel_consumption
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param RoadUserCostRusFuelConsumption body requests.RoadUserCostRusFuelConsumption true "Insert your json"
// @response 201 {object} responses.Success{data=requests.RoadUserCostRusFuelConsumption} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/fuel_consumption [post]
func (sh *SettingHandler) CreateRoadUserCostRucFuelConsumption(c *gin.Context) {
	var request requests.RoadUserCostRusFuelConsumption
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

	roadUserCost, err := sh.settingUseCase.CreateRoadUserCostRucFuelConsumption(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadUserCost, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_fuel_consumption
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RoadUserCostRusFuelConsumption} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/fuel_consumption [get]
func (sh *SettingHandler) GetRoadUserCostRucFuelConsumption(c *gin.Context) {
	roadUserCost, err := sh.settingUseCase.GetRoadUserCostRucFuelConsumption(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(roadUserCost, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_lubricant_consumption
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param RoadUserCostRusLubricantConsumption body requests.RoadUserCostRusLubricantConsumption true "Insert your json"
// @response 201 {object} responses.Success{data=requests.RoadUserCostRusLubricantConsumption} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/lubricant_consumption [post]
func (sh *SettingHandler) CreateRoadUserCostRucLubricantConsumption(c *gin.Context) {
	var request requests.RoadUserCostRusLubricantConsumption
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

	roadUserCost, err := sh.settingUseCase.CreateRoadUserCostRucLubricantConsumption(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadUserCost, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_lubricant_consumption
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RoadUserCostRusLubricantConsumption} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/lubricant_consumption [get]
func (sh *SettingHandler) GetRoadUserCostRucLubricantConsumption(c *gin.Context) {
	roadUserCost, err := sh.settingUseCase.GetRoadUserCostRucLubricantConsumption(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(roadUserCost, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_waste_of_consumption
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param RoadUserCostRusWasteOfConsumption body requests.RoadUserCostRusWasteOfConsumption true "Insert your json"
// @response 201 {object} responses.Success{data=requests.RoadUserCostRusWasteOfConsumption} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/waste_of_consumption [post]
func (sh *SettingHandler) CreateRoadUserCostRucWasteOfConsumption(c *gin.Context) {
	var request requests.RoadUserCostRusWasteOfConsumption
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

	roadUserCost, err := sh.settingUseCase.CreateRoadUserCostRucWasteOfConsumption(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadUserCost, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_waste_of_consumption
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RoadUserCostRusWasteOfConsumption} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/waste_of_consumption [get]
func (sh *SettingHandler) GetRoadUserCostRucWasteOfConsumption(c *gin.Context) {
	roadUserCost, err := sh.settingUseCase.GetRoadUserCostRucWasteOfConsumption(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(roadUserCost, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_maintenance
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param RoadUserCostRusMaintenance body requests.RoadUserCostRusMaintenance true "Insert your json"
// @response 201 {object} responses.Success{data=requests.RoadUserCostRusMaintenance} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/maintenance [post]
func (sh *SettingHandler) CreateRoadUserCostRucMaintenance(c *gin.Context) {
	var request requests.RoadUserCostRusMaintenance
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

	roadUserCost, err := sh.settingUseCase.CreateRoadUserCostRucMaintenance(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadUserCost, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_maintenance
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RoadUserCostRusMaintenance} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/maintenance [get]
func (sh *SettingHandler) GetRoadUserCostRucMaintenance(c *gin.Context) {
	roadUserCost, err := sh.settingUseCase.GetRoadUserCostRucMaintenance(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(roadUserCost, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_travel_time
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param RoadUserCostRusTravelTime body requests.RoadUserCostRusTravelTime true "Insert your json"
// @response 201 {object} responses.Success{data=requests.RoadUserCostRusTravelTime} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/travel_time [post]
func (sh *SettingHandler) CreateRoadUserCostRucTravelTime(c *gin.Context) {
	var request requests.RoadUserCostRusTravelTime
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

	roadUserCost, err := sh.settingUseCase.CreateRoadUserCostRucTravelTime(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadUserCost, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_travel_time
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RoadUserCostRusTravelTime} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/travel_time [get]
func (sh *SettingHandler) GetRoadUserCostRucTravelTime(c *gin.Context) {
	roadUserCost, err := sh.settingUseCase.GetRoadUserCostRucTravelTime(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(roadUserCost, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_vehicle_speed_calculation
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param RoadUserCostRusVehicleSpeedCalculation body requests.RoadUserCostRusVehicleSpeedCalculation true "Insert your json"
// @response 201 {object} responses.Success{data=requests.RoadUserCostRusVehicleSpeedCalculation} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/vehicle_speed_calculation [post]
func (sh *SettingHandler) CreateRoadUserCostRucVehicleSpeedCalculation(c *gin.Context) {
	var request requests.RoadUserCostRusVehicleSpeedCalculation
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

	roadUserCost, err := sh.settingUseCase.CreateRoadUserCostRucVehicleSpeedCalculation(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadUserCost, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_vehicle_speed_calculation
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RoadUserCostRusVehicleSpeedCalculation} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/vehicle_speed_calculation [get]
func (sh *SettingHandler) GetRoadUserCostRucVehicleSpeedCalculation(c *gin.Context) {
	roadUserCost, err := sh.settingUseCase.GetRoadUserCostRucVehicleSpeedCalculation(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(roadUserCost, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_traffic_data
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param RoadUserCostRusTrafficData body requests.RoadUserCostRusTrafficData true "Insert your json"
// @response 201 {object} responses.Success{data=requests.RoadUserCostRusTrafficData} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/traffic_data [post]
func (sh *SettingHandler) CreateRoadUserCostRucTrafficData(c *gin.Context) {
	var request requests.RoadUserCostRusTrafficData
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

	roadUserCost, err := sh.settingUseCase.CreateRoadUserCostRucTrafficData(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadUserCost, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id create_road_user_cost_ruc_traffic_data
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RoadUserCostRusTrafficData} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/road_user_cost/ruc/traffic_data [get]
func (sh *SettingHandler) GetRoadUserCostRucTrafficData(c *gin.Context) {
	roadUserCost, err := sh.settingUseCase.GetRoadUserCostRucTrafficData(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(roadUserCost, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_optimization
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param Optimization body requests.Optimization true "Insert your json"
// @response 201 {object} responses.Success{data=requests.Optimization} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/optimization [post]
func (sh *SettingHandler) CreateOptimization(c *gin.Context) {
	var request requests.Optimization
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

	roadUserCost, err := sh.settingUseCase.CreateOptimization(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadUserCost, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id get_optimization
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.Optimization} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/optimization [get]
func (sh *SettingHandler) GetOptimization(c *gin.Context) {
	roadUserCost, err := sh.settingUseCase.GetOptimization(c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(roadUserCost, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_deterioration_asphalt
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param DeteriorationAsphalt body requests.DeteriorationAsphalt true "Insert your json"
// @response 201 {object} responses.Success{data=requests.DeteriorationAsphalt} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/deterioration/asphalt [post]
func (sh *SettingHandler) CreateDeteriorationAsphalt(c *gin.Context) {
	var request requests.DeteriorationAsphalt
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

	roadWorkEffect, err := sh.settingUseCase.CreateDeteriorationAsphalt(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadWorkEffect, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id get_deterioration_asphalt
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param road_group_id path string false "Insert your road group id"
// @response 200 {object} responses.Success{data=responses.DeteriorationAsphalt} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/deterioration/asphalt/{road_group_id} [get]
func (sh *SettingHandler) GetDeteriorationAsphalt(c *gin.Context) {
	resp, err := sh.settingUseCase.GetDeteriorationAsphalt(c.Param("road_group_id"), c)
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
// @description
// @tags setting_models
// @id create_deterioration_concrete
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param DeteriorationConcrete body requests.DeteriorationConcrete true "Insert your json"
// @response 201 {object} responses.Success{data=requests.DeteriorationConcrete} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/deterioration/concrete [post]
func (sh *SettingHandler) CreateDeteriorationConcrete(c *gin.Context) {
	var request requests.DeteriorationConcrete
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

	roadWorkEffect, err := sh.settingUseCase.CreateDeteriorationConcrete(request, c)
	if err != nil {
		err, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(err)
			c.JSON(err.StatusCode, errResponse)
			return
		}
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(roadWorkEffect, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id get_deterioration_concrete
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param road_group_id path string false "Insert your road group id"
// @response 200 {object} responses.Success{data=responses.DeteriorationConcrete} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/deterioration/concrete/{road_group_id} [get]
func (sh *SettingHandler) GetDeteriorationConcrete(c *gin.Context) {
	resp, err := sh.settingUseCase.GetDeteriorationConcrete(c.Param("road_group_id"), c)
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
// @description
// @id get_hris
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param page query string true "Insert your page number"
// @param limit query string true "Insert your limit number"
// @response 200 {object} responses.Success{data=responses.Pagination{items=[]responses.RefHris}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/hris [get]
func (sh *SettingHandler) GetHris(c *gin.Context) {
	resp, err := sh.settingUseCase.GetHris()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	limitParam := c.Query("limit")
	pageParam := c.Query("page")

	respData := resp.([]responses.RefHris)
	totalItems := int64(len(respData))
	limit, offset, page := helpers.GetlimitOffsetPage(limitParam, pageParam, totalItems)

	if totalItems == 0 {
		respData = []responses.RefHris{}
	} else if limit+offset > totalItems {
		respData = respData[offset:totalItems]
	} else {
		respData = respData[offset : limit+offset]
	}

	pagination := helpers.Pagination(respData, limit, page, totalItems)

	c.JSON(http.StatusOK, responses.SuccessResponse(pagination, http.StatusOK))

}

// @summary
// @description
// @tags setting_models
// @id get_hris_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Hris id"
// @response 200 {object} responses.Success{data=responses.RefHris{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/hris/{id} [get]
func (sh *SettingHandler) GetHrisById(c *gin.Context) {
	resp, err := sh.settingUseCase.GetHrisById(c.Param("id"))
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id get_hris_preview
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RefHrisPreview{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/hris_preview [get]
func (sh *SettingHandler) GetHrisPreview(c *gin.Context) {
	resp, err := sh.settingUseCase.GetHrisPreview()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags setting_models
// @id create_hris
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param Hris body requests.CreateRefHris true "Insert your json"
// @response 201 {object} responses.Success{data=responses.RefHris{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/hris [post]
func (sh *SettingHandler) CreateHris(c *gin.Context) {
	var request requests.CreateRefHris
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

	userID, _ := c.Get("userID")
	uid := int(userID.(float64))
	response, err := sh.settingUseCase.CreateHris(request, uid)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusAccepted, responses.SuccessResponse(response, http.StatusCreated))
}

// @summary
// @description
// @tags setting_models
// @id update_hris
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Hris id"
// @param Hris body requests.UpdateRefHris true "Insert your json"
// @response 201 {object} responses.Success{data=responses.RefHris{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/hris/{id} [put]
func (sh *SettingHandler) UpdateHris(c *gin.Context) {
	var request requests.UpdateRefHris
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

	userID, _ := c.Get("userID")
	uid := int(userID.(float64))
	response, err := sh.settingUseCase.UpdateHris(request, c.Param("id"), uid)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusAccepted, responses.SuccessResponse(response, http.StatusAccepted))
}

// @summary
// @description
// @tags setting_models
// @id Delete_hris
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string false "Hris id"
// @response 201 {object} responses.Success{data=responses.NoData{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/hris/{id} [delete]
func (sh *SettingHandler) DeleteHris(c *gin.Context) {

	userID, _ := c.Get("userID")
	uid := int(userID.(float64))
	_, err := sh.settingUseCase.DeleteHris(c.Param("id"), uid)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusNoContent, responses.SuccessResponse(responses.NoData{}, 204))
}

// @summary
// @description
// @tags setting_models
// @id Create_hris
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=responses.RefHrisPreview{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/hris_import [post]
func (sh *SettingHandler) ImportHris(c *gin.Context) {
	response, err := sh.settingUseCase.ImportHris()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(response, http.StatusOK))

}

// @summary
// @description
// @tags setting_models
// @id get_all_hsms
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param asset_name query string false "Insert your asset name"
// @param page query string false "Insert your page number"
// @param limit query string false "Insert your limit number"
// @response 200 {object} responses.Success{data=responses.Pagination{items=[]responses.HsmsAll}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/hsms [get]
func (sh *SettingHandler) GetAllHsms(c *gin.Context) {

	var request requests.FilterHsms
	request.AssetName = c.Query("asset_name")
	respTable, err := sh.settingUseCase.GetAllHsms(request)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	limitParam := c.Query("limit")
	pageParam := c.Query("page")

	totalItems := int64(len(respTable))
	limit, offset, page := helpers.GetlimitOffsetPage(limitParam, pageParam, totalItems)

	if totalItems == 0 {
		respTable = []responses.HsmsAll{}
	} else if limit+offset > totalItems {
		respTable = respTable[offset:totalItems]
	} else {
		respTable = respTable[offset : limit+offset]
	}

	pagination := helpers.Pagination(respTable, limit, page, totalItems)

	c.JSON(http.StatusOK, responses.SuccessResponse(pagination, http.StatusOK))

}

// @summary
// @description
// @tags setting_models
// @id delete_hsms_by_type_and_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param type path string true "Insert your type"
// @param id path string true "Insert your road ID"
// @Success 204 {object} responses.Success{data=responses.NoData{}} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/settings/hsms/{type}/table/{id} [delete]
func (sh *SettingHandler) DeleteHsmsByTypeAndId(c *gin.Context) {
	typeData := c.Param("type")
	id := c.Param("id")
	err := sh.settingUseCase.DeleteHsmsByTypeAndId(typeData, id)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse("", http.StatusNoContent))

}
