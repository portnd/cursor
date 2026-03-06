package usecases

import (
	"strconv"
	"strings"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/role/domains"
	"gitlab.com/mims-api-service/src/role/requests"
)

type roleUseCase struct {
	roleRepo domains.RoleRepository
}

// init usecase
func NewRoleUseCase(repo domains.RoleRepository) domains.RoleUseCase {
	return &roleUseCase{
		roleRepo: repo,
	}
}

// =========================================================
func (t *roleUseCase) GetRole(param requests.RoleQueryParams) (interface{}, error) {
	total, err := t.roleRepo.CountRoleAll()
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	queryFilter, err := createFilter(param)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	// get limit, offset, page
	limit, offset, page := helpers.GetlimitOffsetPage(param.Limit, param.Page, total)
	data, err := t.roleRepo.GetRoleData(limit, offset, queryFilter) //.GetRoleData(limit, offset, queryFilter)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	if queryFilter != "" {
		total = int64(len(data))
	}
	// pagination
	pagination := helpers.Pagination(data, limit, page, total)
	return pagination, nil
}

func (t *roleUseCase) GetRoleById(roleId int) (models.Role, error) {
	roles, err := t.roleRepo.GetRoleById(roleId)
	if err != nil {
		return roles, responses.NewAppErr(400, err.Error())
	}
	return roles, nil
}
func (t *roleUseCase) GetAccessGroupAll() ([]models.AccessGroup, error) {
	roles, err := t.roleRepo.GetAccessGroupAll()
	if err != nil {
		return roles, responses.NewAppErr(400, err.Error())
	}
	return roles, nil
}

func (t *roleUseCase) GetAccessGroup() ([]models.AccessGroup, error) {
	roles, err := t.roleRepo.GetAccessGroup()
	if err != nil {
		return roles, responses.NewAppErr(400, err.Error())
	}
	return roles, nil
}

func (t *roleUseCase) GetAccessControlByGrpId(grpId int) ([]models.AccessControl, error) {
	accCtrl, err := t.roleRepo.GetAccessControlByGrpId(grpId)
	if err != nil {
		return accCtrl, responses.NewAppErr(400, err.Error())
	}
	return accCtrl, nil
}

func (t *roleUseCase) GetRoleAccessControlAll() (map[string]int, error) {
	roleAccCtrlArr := make(map[string]int)
	roleAccCtrl, err := t.roleRepo.GetRoleAccessControlAll()
	if err != nil {
		return roleAccCtrlArr, responses.NewAppErr(400, err.Error())
	}

	for _, item := range roleAccCtrl {
		roleAccCtrlArr[strconv.Itoa(item.RoleId)+"-"+strconv.Itoa(item.AccessControlId)] = item.AccessControlId
	}

	return roleAccCtrlArr, nil
}

func (t *roleUseCase) GetRoleAccessControlByRoleID(ID int) (map[int]int, error) {
	roleAccCtrlArr := make(map[int]int)
	roleAccCtrl, err := t.roleRepo.GetRoleAccessControlByRoleID(ID)
	if err != nil {
		return roleAccCtrlArr, responses.NewAppErr(400, err.Error())
	}

	for _, item := range roleAccCtrl {
		roleAccCtrlArr[item.AccessControlId] = item.AccessControlId
	}

	return roleAccCtrlArr, nil
}

func (t *roleUseCase) CheckRoleAccessControlByRoleIDAccUD(ID, accID int) (bool, error) {
	is, _ := t.roleRepo.CheckRoleAccessControlByRoleIDAccUD(ID, accID)
	// if err != nil {
	// 	return false, responses.NewAppErr(400, err.Error())
	// }

	return is, nil
}

func (t *roleUseCase) CreateRole(req requests.ReqRoleAccCtrl) (int, error) {
	err := t.roleRepo.CheckRoleDuplicate(0, req.Name)
	if err != nil {
		return 0, responses.NewDuplicatedNameError("Name:duplicate")
	}
	roleId, err := t.roleRepo.CreateRole(req.Name)
	if err != nil {
		return 0, responses.NewInternalServerError()
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, item := range req.AccessControl {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, roleId)
		valueArgs = append(valueArgs, item.AccessControlId)
	}
	err = t.roleRepo.CreateRoleAccessControl(valueStrings, valueArgs)
	if err != nil {
		return 0, responses.NewInternalServerError()
	}
	return roleId, nil
}

func (t *roleUseCase) UpdateRole(roleId int, req requests.ReqRoleAccCtrl) error {
	err := t.roleRepo.CheckRoleDuplicate(roleId, req.Name)
	if err != nil {
		return responses.NewDuplicatedNameError("Name:duplicate")
	}

	err = t.roleRepo.UpdateRoleName(roleId, req.Name)
	if err != nil {
		return responses.NewInternalServerError()
	}

	err = t.roleRepo.DeleteRoleAccessControl(roleId)
	if err != nil {
		return responses.NewInternalServerError()
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, item := range req.AccessControl {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, roleId)
		valueArgs = append(valueArgs, item.AccessControlId)
	}
	err = t.roleRepo.CreateRoleAccessControl(valueStrings, valueArgs)
	if err != nil {
		return responses.NewInternalServerError()
	}
	return nil
}

func (t roleUseCase) DeleteRole(roleId int) error {
	err := t.roleRepo.DeleteRole(roleId)
	if err != nil {
		return responses.NewInternalServerError()
	}
	return nil
}

func createFilter(params requests.RoleQueryParams) (string, error) {
	var queryWhere []string
	if params.Name != "" {
		queryWhere = append(queryWhere, "name like '%"+params.Name+"%'")
	}
	query := strings.Join(queryWhere, " and ")
	return query, nil
}
