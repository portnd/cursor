package repositories

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/role/handlers"
	"gitlab.com/mims-api-service/src/role/usecases"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type roleRepository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewRoleRepositoryHandler(conn *gorm.DB) *handlers.RoleHandler {
	useCase := usecases.NewRoleUseCase(&roleRepository{conn})
	handler := handlers.NewRoleHandler(useCase)
	return handler
}

// //===================== query =====================
func (t *roleRepository) GetRole() ([]models.Role, error) {
	var roles []models.Role
	if err := t.conn.Where("is_active = ?", true).Find(&roles).Error; err != nil {
		return roles, err
	}
	return roles, nil
}

func (t *roleRepository) GetRoleById(roleId int) (models.Role, error) {
	var roles models.Role
	if err := t.conn.Where("id = ?", roleId).Where("is_active = ?", true).First(&roles).Error; err != nil {
		return roles, err
	}
	return roles, nil
}

func (t *roleRepository) GetRoleData(limit, offset int64, filter string) ([]models.Role, error) {
	var role []models.Role
	query := t.conn
	if filter != "" {
		query = query.Where(filter)
	}

	query = query.Where("is_active = ?", true).Order("id").Limit(int(limit)).Offset(int(offset))
	err := query.Find(&role).Error
	if err != nil {
		return role, err
	}
	return role, nil
}

func (t *roleRepository) CountRoleAll() (int64, error) {
	var role []models.Role
	var count int64
	if err := t.conn.Model(&role).Where("is_active = ?", true).Count(&count).Error; err != nil {
		return int64(0), err
	}
	return int64(count), nil
}

func (t *roleRepository) GetMenuAccessControl(accCtrl []string) ([]models.AccessControl, error) {
	var accessControl []models.AccessControl
	if err := t.conn.Where("id IN (?)", accCtrl).Find(&accessControl).Error; err != nil {
		return accessControl, err
	}
	return accessControl, nil
}

func (t *roleRepository) DeleteRoleAccessControl(rouleId int) error {
	var roleAccessControl []models.RoleAccessControl
	if err := t.conn.Where("role_id = ?", rouleId).Delete(roleAccessControl).Error; err != nil {
		return err
	}
	return nil
}

func (t *roleRepository) CreateRoleAccessControl(valueStrings []string, valueArgs []interface{}) error {
	sqlStr := "INSERT INTO role_access_control(role_id,access_control_id) VALUES  %s "
	sqlStr = fmt.Sprintf(sqlStr, strings.Join(valueStrings, ","))
	tx := t.conn.Begin()
	if err := tx.Exec(sqlStr, valueArgs...).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (t *roleRepository) GetAccessGroup() ([]models.AccessGroup, error) {
	var accessGroup []models.AccessGroup
	if err := t.conn.Where("status = ?", 1).Where("parent_id = ?", 0).Find(&accessGroup).Error; err != nil {
		return accessGroup, err
	}
	return accessGroup, nil
}

func (t *roleRepository) GetAccessGroupAll() ([]models.AccessGroup, error) {
	var accessGroup []models.AccessGroup
	if err := t.conn.Where("status = ?", 1).Order("id asc").Find(&accessGroup).Error; err != nil {
		return accessGroup, err
	}
	return accessGroup, nil
}

func (t *roleRepository) GetAccessControlByGrpId(grpId int) ([]models.AccessControl, error) {
	var accCtrl []models.AccessControl
	if err := t.conn.Where("access_grp_id = ?", grpId).Order("seq asc").Find(&accCtrl).Error; err != nil {
		return accCtrl, err
	}
	return accCtrl, nil
}

func (t *roleRepository) GetRoleAccessControlAll() ([]models.RoleAccessControl, error) {
	var roleAccCtrl []models.RoleAccessControl
	if err := t.conn.Find(&roleAccCtrl).Error; err != nil {
		return roleAccCtrl, err
	}
	return roleAccCtrl, nil
}

func (t *roleRepository) GetRoleAccessControlByRoleID(ID int) ([]models.RoleAccessControl, error) {
	var roleAccCtrl []models.RoleAccessControl
	if err := t.conn.Where("role_id = ?", ID).Find(&roleAccCtrl).Error; err != nil {
		return roleAccCtrl, err
	}
	return roleAccCtrl, nil
}

func (t *roleRepository) CheckRoleAccessControlByRoleIDAccUD(ID, accID int) (bool, error) {
	var roleAccCtrl models.RoleAccessControl
	if err := t.conn.Where("role_id = ?", ID).Where("access_control_id = ?", accID).First(&roleAccCtrl).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (t *roleRepository) UpdateRoleName(roleId int, name string) error {
	var role models.Role
	role.Name = name
	if err := t.conn.Model(role).Where("id = ?", roleId).Updates(role).Error; err != nil {
		return err
	}
	return nil
}

func (t *roleRepository) CreateRole(name string) (int, error) {
	var role models.Role
	role.Name = name
	role.IsActive = true
	if err := t.conn.Create(&role).Error; err != nil {
		return 0, err
	}
	return role.Id, nil
}

func (t *roleRepository) CheckRoleDuplicate(roleId int, name string) error {
	var role []models.Role
	query := t.conn
	query = query.Where("is_active = ?", true)
	if roleId != 0 {
		query = query.Where("id <> ?", roleId)
	}
	if err := query.Where("name = ?", name).Find(&role).Error; err != nil {
		return err
	}
	if len(role) > 0 {
		return errors.New("duplicate name '" + name + "'")
	}
	return nil
}

func (t *roleRepository) DeleteRole(roleId int) error {
	if err := t.conn.Exec(fmt.Sprintf("UPDATE role SET is_active = 'f' WHERE id = %d", roleId)).Error; err != nil {
		return err
	}

	var roleAccessControl []models.RoleAccessControl
	if err := t.conn.Where("role_id = ?", roleId).Delete(roleAccessControl).Error; err != nil {
		return err
	}
	return nil
}
