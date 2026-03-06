package services

import (
	"fmt"
	"log"
	"os"

	// helpers "gitlab.com/mims-api-service/helpers/file"

	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	"gorm.io/gorm"
)

type servicesDatabase struct {
	conn *gorm.DB
}

func NewServicesDatabase(conn *gorm.DB) *servicesDatabase {
	return &servicesDatabase{conn: conn}
}

func (s *servicesDatabase) UserInfo(userID int) (*responses.UserInfoRespond, error) {
	var userInfo responses.UserInfoRespond
	user, err := s.GetUserById(userID)
	var newRole []responses.RoleRespond
	if err != nil {
		return &userInfo, responses.NewAppErr(400, err.Error())
	}

	roles, err := s.GetUserRole(userID)
	if err != nil {
		return &userInfo, responses.NewAppErr(400, err.Error())
	}

	var ids []int
	for _, item := range roles {
		ids = append(ids, item.RoleID)
	}
	roleAccessControl, err := s.GetRoleAccessControl(ids)
	if err != nil {
		return &userInfo, responses.NewAppErr(400, err.Error())
	}

	if user.ProfileImgPath == "" {
		user.ProfileImgPath = ""
	} else {
		user.ProfileImgPath = os.Getenv("STORAGE_IP") + "/" + user.ProfileImgPath
	}
	userInfo.RefUserOwnerID = user.RefUserOwnerID
	userInfo.RefDepot = user.RefDepot
	userInfo.Roles = newRole
	userInfo.AccessControl = roleAccessControl
	return &userInfo, nil
}

func (s *servicesDatabase) GetUserById(userID int) (models.UserDepartment, error) {
	var user models.UserDepartment
	if err := s.conn.Where("id = ?", userID).Preload("RefDepot").Find(&user).Error; err != nil {
		fmt.Println(err)
		return user, err

	}
	return user, nil
}

func (s *servicesDatabase) GetUserRole(userID int) ([]models.UserRole, error) {
	var userRole []models.UserRole
	if err := s.conn.Where("user_id = ?", userID).Find(&userRole).Group("role_id").Error; err != nil {
		return userRole, err

	}
	return userRole, nil
}
func (s *servicesDatabase) GetRoleByIds(ids []int) ([]models.Role, error) {
	var roles []models.Role
	query := s.conn
	if err := query.Where("id IN (?)", ids).Find(&roles).Error; err != nil {
		log.Println(err.Error())
		return roles, err
	}
	return roles, nil
}

func (s *servicesDatabase) GetRoleAccessControl(roleIds []int) ([]models.AccessControl, error) {
	var accessControl []models.AccessControl
	if err := s.conn.Where("role_access_control.role_id in (?)", roleIds).Joins("JOIN role_access_control on access_control.id = role_access_control.access_control_id").Find(&accessControl).Group("role_access_control.access_control_id").Error; err != nil {
		fmt.Println(err)
		return accessControl, err

	}
	return accessControl, nil //
}
