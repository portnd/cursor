package repositories

import (
	"fmt"
	"log"
	"strings"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/user/handlers"
	"gitlab.com/mims-api-service/src/user/usecases"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type userRepository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewUserRepositoryHandler(conn *gorm.DB) *handlers.UserHandler {
	useCase := usecases.NewUserUseCase(&userRepository{conn})
	handler := handlers.NewUserHandler(useCase)
	return handler
}

// //===================== query =====================
func (t *userRepository) GetUserAll(user *[]models.Users, offset, limit int) (err error) {
	query := t.conn
	query = query.Offset(helpers.QueryOffset(offset))
	query = query.Limit(helpers.QueryLimit(limit))
	if err := query.Order("id asc").Find(&user).Error; err != nil {
		return err
	}
	return nil
}

func (t *userRepository) GetRole() ([]models.Role, error) {
	var role []models.Role
	if err := t.conn.Where("is_active = ?", true).Order("id").Find(&role).Error; err != nil {
		return role, err

	}
	return role, nil
}

func (t *userRepository) GetUserRole(userId int) ([]models.UserRole, error) {
	var userRole []models.UserRole
	if err := t.conn.Where("user_id = ?", userId).Find(&userRole).Group("role_id").Error; err != nil {
		return userRole, err

	}
	return userRole, nil
}

func (t *userRepository) GetRoleAccessControl(roleIds []int) ([]models.AccessControl, error) {
	var accessControl []models.AccessControl
	if err := t.conn.Where("role_access_control.role_id in (?)", roleIds).Joins("JOIN role_access_control on access_control.id = role_access_control.access_control_id").Find(&accessControl).Group("role_access_control.access_control_id").Error; err != nil {
		fmt.Println(err)
		return accessControl, err

	}
	return accessControl, nil //
}

func (t *userRepository) CheckPassword(userId int) (models.UsersPassword, error) {
	var user models.UsersPassword
	if err := t.conn.Where("id = ?", userId).First(&user).Error; err != nil {
		fmt.Println(err)
		return user, err

	}
	return user, nil
}

func (t *userRepository) GetUserById(userId int) (models.UserDepartment, error) {
	var user models.UserDepartment
	if err := t.conn.Where("id = ?", userId).Preload("RefUserOwner").Preload("RefDepot").Find(&user).Error; err != nil {
		fmt.Println(err)
		return user, err

	}
	return user, nil
}

func (t *userRepository) CreateUserRole(userID int, roleID []int) error {
	var userRole models.UserRole
	if err := t.conn.Where("user_id = ?", userID).Delete(userRole).Error; err != nil {
		return err
	}

	var role models.UserRole
	for _, roleID := range roleID {
		role.RoleID = roleID
		role.UserId = userID
		if err := t.conn.Create(&role).Error; err != nil {
			return err
		}
	}
	return nil
}

func (t *userRepository) GetUserByUsername(username string) ([]models.Users, error) {
	var user []models.Users
	if err := t.conn.Where("username = ?", username).Find(&user).Error; err != nil {
		return user, err

	}
	return user, nil
}

func (t *userRepository) DeleteUserById(userId int) error {
	var userRole models.UserRole
	if err := t.conn.Where("user_id = ?", userId).Delete(userRole).Error; err != nil {
		return err
	}

	var user models.Users
	if err := t.conn.Where("id = ?", userId).Delete(user).Error; err != nil {
		return err
	}
	return nil
}

func (t *userRepository) CreateUser(user models.Users) (models.Users, error) {
	// helpers.PrintlnJson(user)

	if err := t.conn.Create(&user).Error; err != nil {
		fmt.Println(err)
		return models.Users{}, err
	}
	return user, nil
}

func (t *userRepository) UpdateUser(userId int, user models.Users) error {
	// 	userData.ProfileImgPath = ""
	// userData.Username = req.Username
	// userData.Email = req.Email
	// userData.DepartmentId = req.DepartmentId
	// userData.Firstname = req.Firstname
	// userData.Lastname = req.Lastname
	// userData.Status = req.Status
	// userData.CreatedBy = req.CreatedBy
	// userData.Tel = req.Tel
	if err := t.conn.Debug().Table("users").Where("id = ?", userId).Updates(map[string]interface{}{"profile_img_path": user.ProfileImgPath, "email": user.Email, "ref_user_owner_id": user.RefUserOwnerID, "ref_depot_id": user.RefDepotID, "firstname": user.Firstname, "lastname": user.Lastname, "status": user.Status, "tel": user.Tel}).Error; err != nil {
		fmt.Println(err)
		return err
	}
	if !user.Status {
		if err := t.conn.Raw("UPDATE users SET status = ? WHERE id = ? RETURNING id", false, userId).Scan(&user).Error; err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func (t *userRepository) UpdatePassword(userId int, user models.Users) error {
	// 	userData.ProfileImgPath = ""
	// userData.Username = req.Username
	// userData.Email = req.Email
	// userData.DepartmentId = req.DepartmentId
	// userData.Firstname = req.Firstname
	// userData.Lastname = req.Lastname
	// userData.Status = req.Status
	// userData.CreatedBy = req.CreatedBy
	// userData.Tel = req.Tel
	if err := t.conn.Table("users").Select("*").Where("id = ?", userId).Updates(map[string]interface{}{"password": user.Password}).Error; err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (t *userRepository) CountUserAll() (int64, error) {
	var users []models.Users
	var count int64
	if err := t.conn.Model(&users).Count(&count).Error; err != nil {
		return int64(0), err
	}
	return int64(count), nil
}

func (t *userRepository) CountUserFilter(filter, roleID string) (int64, error) {
	var users []models.Users
	var count int64
	if roleID != "" {
		if err := t.conn.Model(&users).
			Joins("JOIN user_role ON user_role.user_id = users.id").
			Where(filter).
			Count(&count).Error; err != nil {
			return int64(0), err
		}
		return int64(count), nil

	} else {
		if err := t.conn.Model(&users).Where(filter).Count(&count).Error; err != nil {
			return int64(0), err
		}
		return int64(count), nil
	}

}

func (t *userRepository) GetUserData(limit, offset int64, filter string) ([]models.UserDepartment, error) {
	var users []models.UserDepartment
	query := t.conn
	if filter != "" {
		if strings.Contains(filter, "user_role.role_id") {
			query = query.Joins("JOIN user_role ON user_role.user_id = users.id")
		}
		fmt.Println(filter)
		query = query.Where(filter)

	}

	query = query.Order("id").Limit(int(limit)).Offset(int(offset))
	query = query.Preload("RefUserOwner")
	query = query.Preload("RefDepot")
	err := query.Find(&users).Error
	if err != nil {
		log.Println(err.Error())
		return users, err
	}

	return users, nil
}

func (t *userRepository) GetRoleByIds(ids []int) ([]models.Role, error) {
	var roles []models.Role
	query := t.conn
	if err := query.Where("id IN (?)", ids).Find(&roles).Error; err != nil {
		log.Println(err.Error())
		return roles, err
	}
	return roles, nil
}

func (t *userRepository) CheckEmail(userID int, email string) ([]models.Users, error) {
	var users []models.Users
	query := t.conn
	if userID == 0 {
		if err := query.Where("email = ?", email).Find(&users).Error; err != nil {
			log.Println(err.Error())
			return users, err
		}
	} else {
		if err := query.Where("id != ?", userID).Where("email = ?", email).Find(&users).Error; err != nil {
			log.Println(err.Error())
			return users, err
		}
	}

	return users, nil
}
