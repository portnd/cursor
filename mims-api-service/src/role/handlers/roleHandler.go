package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/role/domains"
	"gitlab.com/mims-api-service/src/role/requests"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gopkg.in/validator.v2"
)

// init handler
type RoleHandler struct {
	roleUseCase domains.RoleUseCase
}

// init handler
func NewRoleHandler(usecase domains.RoleUseCase) *RoleHandler {
	return &RoleHandler{
		roleUseCase: usecase,
	}
}

// ================================== start function  ==================================

// request form
type LoginCredentials struct {
	Email    string `form:"email" validate:"min=1"`
	Password string `form:"password" validate:"min=1"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type Param struct {
	Offset      int               `json:"offset"`
	Limit       int               `json:"limit"`
	Search      map[string]string `json:"search"`
	Column      string            `json:"column"`
	Dir         string            `json:"dir"`
	ColumnOrder map[string]string `json:"column_order"`
}

type RoleData struct {
	Id          int           `json:"id"`
	Role        string        `json:"role"`
	AccessGroup []AccessGroup `json:"access_group"`
}

type AccessGroup struct {
	// Id   int         `json:"id" extensions:"x-order=0"`
	Name string `json:"name" extensions:"x-order=0"`
	Menu []Menu `json:"menu" extensions:"x-order=1"`
	// Menu []struct {
	// 	Name          string          `json:"name" extensions:"x-order=0"`
	// 	AccessControl []AccessControl `json:"access_control" extensions:"x-order=1"`
	// } `json:"menu" extensions:"x-order=1"`
	AccessControl []AccessControl `json:"-" extensions:"x-order=1"`
}

type MenuData struct {
	Name          string      `json:"name" extensions:"x-order=0"`
	AccessControl interface{} `json:"access_control" extensions:"x-order=2"`
}

type Menu struct {
	Name          string          `json:"name" extensions:"x-order=0"`
	AccGrpId      int             `json:"-" extensions:"x-order=1"`
	AccessControl []AccessControl `json:"access_control" extensions:"x-order=2"`
}

type AccessControl struct {
	Id        int    `json:"id" extensions:"x-order=0"`
	Name      string `json:"name" extensions:"x-order=1"`
	AccessKey string `json:"access_key" extensions:"x-order=2"`
	IsCheck   bool   `json:"is_check" extensions:"x-order=3"`
	ParentId  int    `json:"-" extensions:"x-order=4"`
	AccGrpId  int    `json:"-" extensions:"x-order=5"`
}

// @summary
// @description
// @tags role
// @id get_role
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param name query string false "Anonymous user" example(Anonymous user)
// @response 200 {object} responses.Role "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roles [get]
func (t *RoleHandler) GetRole(c *gin.Context) {
	var queryParams requests.RoleQueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	data, err := t.roleUseCase.GetRole(queryParams)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags role
// @id get_role_by_id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your role ID"
// @response 200 {object} RoleData "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roles/{id} [get]
func (t *RoleHandler) GetRoleById(c *gin.Context) {
	roleId, _ := strconv.Atoi(c.Params.ByName("id"))
	data, err := t.getRoleData(roleId)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags role
// @id create_role
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param Role body requests.ReqRoleAccCtrl true "Insert your data"
// @response 201 {object} RoleData "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roles [post]
func (t *RoleHandler) CreateRole(c *gin.Context) {
	var req requests.ReqRoleAccCtrl //ReqRoleCreate
	err := c.ShouldBind(&req)
	if err != nil {
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
	roleId, err := t.roleUseCase.CreateRole(req)
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

	data, err := t.getRoleData(roleId)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(201, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags role
// @id update_role
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your role ID"
// @param Role body requests.ReqRoleAccCtrl true "Update your data"
// @response 200 {object} RoleData "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roles/{id} [put]
func (t *RoleHandler) UpdateRole(c *gin.Context) {
	roleId, _ := strconv.Atoi(c.Params.ByName("id"))
	var req requests.ReqRoleAccCtrl
	err := c.ShouldBind(&req)
	if err != nil {
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

	err = t.roleUseCase.UpdateRole(roleId, req)
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

	data, err := t.getRoleData(roleId)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.JSON(200, responses.SuccessResponse(data, 200))
}

// @summary
// @description
// @tags role
// @id delete_role
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param id path string true "Insert your role ID"
// @response 204 {object} responses.Empty "OK"
// @response 400 {object} responses.BadRequestErrorResponse "Bad Request"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roles/{id} [delete]
func (t *RoleHandler) DeleteRole(c *gin.Context) {
	roleId, _ := strconv.Atoi(c.Params.ByName("id"))
	err := t.roleUseCase.DeleteRole(roleId)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	c.Status(204)
}

// ///////////////////////////
func (t *RoleHandler) getRoleData(roleId int) (interface{}, error) {
	menuNames := make(map[int]string)
	role, err := t.roleUseCase.GetRoleById(roleId)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	asccessGrps, err := t.roleUseCase.GetAccessGroupAll()
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	var accessGroups []AccessGroup
	for _, item := range asccessGrps {
		if item.ParentId == 0 {
			grpId := item.Id
			accessControlData, err := t.roleUseCase.GetAccessControlByGrpId(grpId)

			if err != nil {
				return "", err
			}
			var accArr []AccessControl
			for _, item := range accessControlData {
				accCtrlId := item.Id
				var acc AccessControl
				acc.Id = accCtrlId
				acc.AccessKey = item.AccessKey
				acc.Name = item.AccessDesc
				acc.ParentId = item.AccessGrpParentId
				acc.AccGrpId = item.AccessGrpId
				accArr = append(accArr, acc)
			}

			var accessGroup AccessGroup
			accessGroup.Name = item.Name
			accessGroup.AccessControl = accArr
			accessGroups = append(accessGroups, accessGroup)
		}
		menuNames[item.Id] = item.Name
	}

	// roleAccCtrlArr, err := t.roleUseCase.GetRoleAccessControlByRoleID(roleId)
	// if err != nil {
	// 	return "", responses.NewAppErr(400, err.Error())
	// }

	var menus []Menu
	parentIndex := make(map[int]int)
	for _, item := range accessGroups {
		for _, item2 := range item.AccessControl {
			if item2.ParentId == 0 {
				continue
			}
			var menu Menu
			menu.Name = menuNames[item2.ParentId]
			menu.AccGrpId = item2.AccGrpId
			var accCtrlDataArr []AccessControl
			for _, item3 := range item.AccessControl {
				var accCtrlData AccessControl
				accCtrlData.Id = item3.Id
				accCtrlData.Name = item3.Name
				accCtrlData.AccessKey = item3.AccessKey
				// if roleAccCtrlArr[strconv.Itoa(roleId)+"-"+strconv.Itoa(item2.Id)] == item3.Id {
				// 	accCtrlData.IsCheck = true
				// } else {
				accCtrlData.IsCheck = false
				// }
				if item2.ParentId == item3.ParentId {
					accCtrlDataArr = append(accCtrlDataArr, accCtrlData)
				}
			}

			if parentIndex[item2.ParentId] == 0 {
				parentIndex[item2.ParentId] = item2.ParentId
				menu.AccessControl = accCtrlDataArr
				menus = append(menus, menu)
			}
		}
	}

	var accessGroupDataArr []AccessGroup
	for _, item := range asccessGrps {
		if item.ParentId == 0 {
			var accessGroupData AccessGroup
			var data []Menu
			for _, menu := range menus {
				if item.Id == menu.AccGrpId {
					data = append(data, menu)
				}
			}
			accessGroupData.Name = item.Name
			accessGroupData.Menu = data
			accessGroupDataArr = append(accessGroupDataArr, accessGroupData)
		}
	}
	// return accessGroupDataArr, nil
	var accessGroupData []AccessGroup
	for _, item := range accessGroupDataArr {
		var accessGroup AccessGroup
		var menus []Menu
		for _, menuItem := range item.Menu {
			var menu Menu
			var accCtrlDatas []AccessControl
			// helpers.PrintlnJson(menu)
			for _, acc := range menuItem.AccessControl {
				var accCtrlData AccessControl
				copier.Copy(&accCtrlData, &acc)
				// if acc.Id =
				isCheck, _ := t.roleUseCase.CheckRoleAccessControlByRoleIDAccUD(roleId, acc.Id)
				if isCheck {
					accCtrlData.IsCheck = true
				} else {
					accCtrlData.IsCheck = false
				}
				accCtrlDatas = append(accCtrlDatas, accCtrlData)
			}
			menu.Name = menuItem.Name
			menu.AccessControl = accCtrlDatas
			menus = append(menus, menu)
		}
		accessGroup.Menu = menus
		accessGroup.Name = item.Name
		accessGroupData = append(accessGroupData, accessGroup)
	}

	var roleData RoleData
	roleData.Id = roleId
	roleData.Role = role.Name
	roleData.AccessGroup = accessGroupData
	return roleData, nil
}

// @summary
// @description
// @tags role
// @id GetRoleAccessControlAll
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} AccessGroup "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/roles/access_control [get]
func (t *RoleHandler) GetRoleAccessControlAll(c *gin.Context) {
	menuNames := make(map[int]string)
	asccessGrps, err := t.roleUseCase.GetAccessGroupAll()
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.JSON(400, errResponse)
		return
	}
	var accessGroups []AccessGroup
	for _, item := range asccessGrps {
		if item.ParentId == 0 {
			grpId := item.Id
			accessControlData, err := t.roleUseCase.GetAccessControlByGrpId(grpId)

			if err != nil {
				errResponse := responses.FailRespone(err)
				c.JSON(400, errResponse)
				return
			}
			var accArr []AccessControl
			for _, item := range accessControlData {
				accCtrlId := item.Id
				var acc AccessControl
				acc.Id = accCtrlId
				acc.AccessKey = item.AccessKey
				acc.Name = item.AccessDesc
				acc.ParentId = item.AccessGrpParentId
				acc.AccGrpId = item.AccessGrpId
				accArr = append(accArr, acc)
			}

			var accessGroup AccessGroup
			accessGroup.Name = item.Name
			accessGroup.AccessControl = accArr
			accessGroups = append(accessGroups, accessGroup)
		}
		menuNames[item.Id] = item.Name
	}

	var menus []Menu
	parentIndex := make(map[int]int)
	for _, item := range accessGroups {
		for _, item2 := range item.AccessControl {
			if item2.ParentId == 0 {
				continue
			}
			var menu Menu
			menu.Name = menuNames[item2.ParentId]
			menu.AccGrpId = item2.AccGrpId
			var accCtrlDataArr []AccessControl
			for _, item3 := range item.AccessControl {
				var accCtrlData AccessControl
				accCtrlData.Id = item3.Id
				accCtrlData.Name = item3.Name
				accCtrlData.AccessKey = item3.AccessKey
				accCtrlData.IsCheck = false
				if item2.ParentId == item3.ParentId {
					accCtrlDataArr = append(accCtrlDataArr, accCtrlData)
				}
			}

			if parentIndex[item2.ParentId] == 0 {
				parentIndex[item2.ParentId] = item2.ParentId
				menu.AccessControl = accCtrlDataArr
				menus = append(menus, menu)
			}

		}

	}
	var accessGroupDataArr []AccessGroup
	for _, item := range asccessGrps {
		if item.ParentId == 0 {
			var accessGroupData AccessGroup
			var data []Menu
			for _, menu := range menus {
				if item.Id == menu.AccGrpId {
					data = append(data, menu)
				}
			}
			accessGroupData.Name = item.Name
			if len(data) == 0 {
				accessGroupData.Menu = []Menu{}
			} else {
				accessGroupData.Menu = data
			}
			accessGroupDataArr = append(accessGroupDataArr, accessGroupData)
		}
	}
	c.JSON(200, responses.SuccessResponse(accessGroupDataArr, 200))
}
