package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/databases"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
)

func CheckPermission(c *gin.Context, userID int, permissions []string) {
	// return func(c *gin.Context) {

	userInfo, err := UserInfo(userID)
	if err != nil {
		errResponse := responses.FailRespone(responses.NewAppErr(http.StatusBadRequest, err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, errResponse)
	}

	myPermission := []string{}
	for _, item := range userInfo.AccessControl {
		myPermission = append(myPermission, item.AccessKey)
	}
	if !helpers.HasPermission(permissions, myPermission) {
		errResponse := responses.FailRespone(responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION))
		c.AbortWithStatusJSON(http.StatusForbidden, errResponse)
	}
	// }

}

func UserInfo(userId int) (*responses.UserInfoRespond, error) {
	var userInfo responses.UserInfoRespond
	user, err := GetUserById(userId)
	var newRole []responses.RoleRespond
	if err != nil {
		return &userInfo, responses.NewAppErr(400, err.Error())
	}

	roles, err := GetUserRole(userId)
	if err != nil {
		return &userInfo, responses.NewAppErr(400, err.Error())
	}

	var ids []int
	for _, item := range roles {
		ids = append(ids, item.RoleID)
	}

	roleData, err := GetRoleByIds(ids)
	if err != nil {
		return &userInfo, responses.NewAppErr(400, err.Error())
	}
	copier.Copy(&newRole, &roleData)
	for i := range roleData {
		newRole[i].CreatedAt = helpers.SetTimeToString(roleData[i].CreatedAt)
		newRole[i].UpdatedAt = helpers.SetTimeToString(roleData[i].UpdatedAt)
	}
	roleAccessControl, err := GetRoleAccessControl(ids)
	if err != nil {
		return &userInfo, responses.NewAppErr(400, err.Error())
	}

	if user.ProfileImgPath == "" {
		user.ProfileImgPath = ""
	} else {
		user.ProfileImgPath = os.Getenv("STORAGE_IP") + "/" + user.ProfileImgPath
	}

	userInfo.Roles = newRole
	userInfo.AccessControl = roleAccessControl
	return &userInfo, nil
}

func GetUserById(userId int) (models.UserDepartment, error) {
	var user models.UserDepartment
	if err := databases.DB.Where("id = ?", userId).Preload("RefUserOwner").Preload("RefDepot").Find(&user).Error; err != nil {
		fmt.Println(err)
		return user, err

	}
	return user, nil
}

func GetUserRole(userId int) ([]models.UserRole, error) {
	var userRole []models.UserRole
	if err := databases.DB.Where("user_id = ?", userId).Find(&userRole).Group("role_id").Error; err != nil {
		return userRole, err

	}
	return userRole, nil
}
func GetRoleByIds(ids []int) ([]models.Role, error) {
	var roles []models.Role
	query := databases.DB
	if err := query.Where("id IN (?)", ids).Find(&roles).Error; err != nil {
		log.Println(err.Error())
		return roles, err
	}
	return roles, nil
}

func GetRoleAccessControl(roleIds []int) ([]models.AccessControl, error) {
	var accessControl []models.AccessControl
	if err := databases.DB.Where("role_access_control.role_id in (?)", roleIds).Joins("JOIN role_access_control on access_control.id = role_access_control.access_control_id").Find(&accessControl).Group("role_access_control.access_control_id").Error; err != nil {
		fmt.Println(err)
		return accessControl, err

	}
	return accessControl, nil //
}
