package domains

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/role/requests"
)

// business logic
type RoleUseCase interface {
	GetRole(requests.RoleQueryParams) (interface{}, error)
	GetRoleById(int) (models.Role, error)
	CreateRole(requests.ReqRoleAccCtrl) (int, error)
	UpdateRole(int, requests.ReqRoleAccCtrl) error
	DeleteRole(roleId int) error
	GetAccessGroup() ([]models.AccessGroup, error)
	GetAccessGroupAll() ([]models.AccessGroup, error)
	GetAccessControlByGrpId(grpId int) ([]models.AccessControl, error)
	GetRoleAccessControlAll() (map[string]int, error)
	GetRoleAccessControlByRoleID(ID int) (map[int]int, error)
	CheckRoleAccessControlByRoleIDAccUD(ID, accID int) (bool, error)
}

// อะไรเชื่อมต่อกับ DB
type RoleRepository interface {
	CountRoleAll() (int64, error)
	GetRole() ([]models.Role, error)
	GetRoleById(int) (models.Role, error)
	GetRoleData(int64, int64, string) ([]models.Role, error)
	// GetMenuAll() ([]models.Menu, error)
	DeleteRoleAccessControl(roleId int) error
	CreateRoleAccessControl([]string, []interface{}) error
	GetAccessGroup() ([]models.AccessGroup, error)
	GetAccessGroupAll() ([]models.AccessGroup, error)
	GetAccessControlByGrpId(grpId int) ([]models.AccessControl, error)
	GetRoleAccessControlAll() ([]models.RoleAccessControl, error)
	GetRoleAccessControlByRoleID(ID int) ([]models.RoleAccessControl, error)
	CheckRoleAccessControlByRoleIDAccUD(ID int, accID int) (bool, error)

	UpdateRoleName(int, string) error
	CreateRole(string) (int, error)
	CheckRoleDuplicate(int, string) error
	DeleteRole(roleId int) error
}
