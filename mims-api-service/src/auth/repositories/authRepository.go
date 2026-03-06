package repositories

import (
	"fmt"

	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/auth/handlers"
	"gitlab.com/mims-api-service/src/auth/usecases"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type authRepository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewAuthRepositoryHandler(conn *gorm.DB) *handlers.AuthHandler {
	useCase := usecases.NewAuthUseCase(&authRepository{conn})
	handler := handlers.NewAuthHandler(useCase)
	return handler
}

// //===================== query =====================
func (t *authRepository) GetUserByUserName(username string) (models.Users, error) {
	user := models.Users{}
	err := t.conn.Where("username = ?", username).First(&user).Error
	if err != nil {
		return models.Users{}, err
	}
	return user, nil
}

func (t *authRepository) GetUserByID(userId uint) (models.Users, error) {
	user := models.Users{}
	err := t.conn.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return models.Users{}, err
	}
	return user, nil
}

func (t *authRepository) GetUserByEmail(email string) (models.Users, error) {
	user := models.Users{}
	err := t.conn.Where("email = ?", email).First(&user).Error
	if err != nil {
		return models.Users{}, err
	}

	return user, nil
}

func (t *authRepository) GetUserByResetPasswordToken(token string) (models.Users, error) {
	user := models.Users{}
	err := t.conn.Where("reset_password_token = ?", token).First(&user).Error
	if err != nil {
		return models.Users{}, err
	}
	return user, nil
}

func (t *authRepository) UpdateUser(user models.Users) error {
	err := t.conn.Model(&models.Users{}).Updates(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *authRepository) CreateAuth(authDetail models.Auth) error {
	err := t.conn.Create(&authDetail).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *authRepository) GetAuthByUserId(userId uint) (models.Auth, error) {
	authDetail := models.Auth{}

	err := t.conn.Where("user_id = ?", userId).First(&authDetail).Error
	if err != nil {
		return models.Auth{}, err
	}

	return authDetail, nil
}

func (t *authRepository) DeleteAuthByUserId(userId uint) error {
	err := t.conn.Where("user_id = ?", userId).Delete(&models.Auth{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *authRepository) UpdateAuth(authDetail models.Auth) error {
	err := t.conn.Model(&models.Auth{}).Updates(&authDetail).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *authRepository) GetRoleByUser(userId uint) ([]models.UserRole, error) {
	var userRole []models.UserRole
	err := t.conn.Where("user_id = ?", userId).Find(&userRole).Error
	if err != nil {
		return userRole, err
	}

	return userRole, nil
}

func (t *authRepository) GetAccessControlByRole(roleId int) ([]models.RoleAccessJoinControl, error) {
	var roleAccessControl []models.RoleAccessJoinControl

	err := t.conn.Select("*").Table("role_access_control").Joins("JOIN access_control on role_access_control.access_control_id =  access_control.id").Where("role_id = ?", roleId).Find(&roleAccessControl).Error
	if err != nil {
		fmt.Println(err)
		return roleAccessControl, err
	}

	return roleAccessControl, nil
}

func (t *authRepository) GetUserByVerifyEmailToken(token string) (models.Users, error) {
	user := models.Users{}
	err := t.conn.Where("verify_email_token = ?", token).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
