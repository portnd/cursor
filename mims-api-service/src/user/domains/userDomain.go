package domains

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
)

// business logic
type UserUseCase interface {
	UserInfo(int) (interface{}, error)
	UpdateUserInfo(int, requests.UserInfoUpdateReq) error
	UserInfoChangePassword(userId int, req requests.UpdatePasswordUserInfoReq) error
	GetUser(requests.UserQueryParams) (interface{}, error)
	GetUserById(int) (interface{}, error)
	GetUserByUsername(username string) (interface{}, error)
	CreateUser(requests.UserReq) (interface{}, error)
	UpdateUser(int, requests.UserReq) (interface{}, error)
	DeleteUserById(userId int) error
	CheckPassword(userId int, password string) (bool, error)
}

// อะไรเชื่อมต่อกับ DB
type UserRepository interface {
	GetUserData(int64, int64, string) ([]models.UserDepartment, error)
	CountUserAll() (int64, error)
	CountUserFilter(string, string) (int64, error)
	GetUserAll(user *[]models.Users, offset, limit int) error
	CheckPassword(userID int) (models.UsersPassword, error)
	GetUserById(UserID int) (models.UserDepartment, error)
	GetUserByUsername(string) ([]models.Users, error)
	GetRole() ([]models.Role, error)
	GetUserRole(int) ([]models.UserRole, error)
	CreateUserRole(userID int, roleID []int) error
	GetRoleAccessControl([]int) ([]models.AccessControl, error)
	CreateUser(models.Users) (models.Users, error)
	UpdateUser(int, models.Users) error
	DeleteUserById(int) error
	GetRoleByIds(ids []int) ([]models.Role, error)
	UpdatePassword(userId int, user models.Users) error
	CheckEmail(userID int, email string) ([]models.Users, error)
}
