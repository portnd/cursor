package usecases

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/user/domains"

	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo domains.UserRepository
}

// init usecase
func NewUserUseCase(repo domains.UserRepository) domains.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

type UserInfoRespond struct {
	models.UserDepartment
	Roles         []RoleRespond          `json:"roles"`
	AccessControl []models.AccessControl `json:"access_control"`
}

type RoleRespond struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CreatedBy int    `json:"-"`
	UpdatedBy int    `json:"-"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	IsActive  bool   `json:"-"`
}

// =========================================================
func (t *userUseCase) GetUser(param requests.UserQueryParams) (interface{}, error) {
	// count all users
	total, err := t.userRepo.CountUserAll()
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}
	// get limit, offset, page
	limit, offset, page := helpers.GetlimitOffsetPage(param.Limit, param.Page, total)
	// create filter
	// return param, nil
	queryFilter, err := createFilter(param)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}
	// query
	users, err := t.userRepo.GetUserData(limit, offset, queryFilter)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}
	// return users, nil

	// roles, err := t.userRepo.GetRole()
	// if err != nil {
	// 	return "", responses.NewAppErr(400, err.Error())
	// }
	var usersDatas []responses.UsersDataRes
	for _, item := range users {
		userRole, _ := t.userRepo.GetUserRole(int(item.Id))
		ids := []int{}
		for _, item := range userRole {
			ids = append(ids, item.RoleID)
		}

		roleData, _ := t.userRepo.GetRoleByIds(ids)

		var usersData responses.UsersDataRes
		usersData.CreatedBy = item.CreatedBy
		usersData.RefUserOwnerID = item.RefUserOwnerID
		usersData.RefUserOwner = item.RefUserOwner
		usersData.RefDepotID = item.RefDepotID
		usersData.RefDepot = item.RefDepot
		usersData.Email = item.Email
		usersData.Firstname = item.Firstname
		usersData.Lastname = item.Lastname
		usersData.Id = item.Id
		// usersData.Password = item.Password
		if item.ProfileImgPath == "" {
			usersData.ProfileImgPath = ""
		} else {
			usersData.ProfileImgPath = os.Getenv("STORAGE_IP") + "/" + item.ProfileImgPath
		}

		usersData.Status = item.Status
		usersData.Tel = item.Tel
		usersData.UpdatedBy = item.UpdatedBy
		usersData.Username = item.Username
		// usersData.Roles = role
		rolesRes := []responses.RoleRes{}
		isRole := true
		role, _ := strconv.Atoi(param.Permission)
		for _, item := range roleData {
			if param.Permission != "" {
				if item.Id == role {
					isRole = true
				} else {
					isRole = false
				}

			}
			role := responses.RoleRes{}
			role.Id = item.Id
			role.Name = item.Name
			role.IsChecked = true
			rolesRes = append(rolesRes, role)
		}

		// usersData.Department = item.Department
		usersData.Roles = rolesRes

		// old logic to make some data lose
		/*
			if isRole {
			usersDatas = append(usersDatas, usersData)
			}
		*/

		// new logic
		_ = isRole
		usersDatas = append(usersDatas, usersData)
	}

	// new logic to count the users
	total, err = t.userRepo.CountUserFilter(queryFilter, param.Permission)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	// old logic to count all users (this code make the result has no next page)
	/*
		if queryFilter != "" {
			total = int64(len(usersDatas))
		}
	*/

	var data interface{}
	if len(usersDatas) == 0 {
		data = []string{}
	} else {
		data = usersDatas
	}

	// pagination
	pagination := helpers.Pagination(data, limit, page, total)
	return pagination, nil
}

func (t *userUseCase) UserInfo(userId int) (interface{}, error) {
	user, err := t.userRepo.GetUserById(userId)
	var newRole []RoleRespond
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	roles, err := t.userRepo.GetUserRole(userId)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	var ids []int
	for _, item := range roles {
		ids = append(ids, item.RoleID)
	}

	roleData, err := t.userRepo.GetRoleByIds(ids)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}
	copier.Copy(&newRole, &roleData)
	for i := range roleData {
		newRole[i].CreatedAt = helpers.SetTimeToString(roleData[i].CreatedAt)
		newRole[i].UpdatedAt = helpers.SetTimeToString(roleData[i].UpdatedAt)
	}
	roleAccessControl, err := t.userRepo.GetRoleAccessControl(ids)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	if user.ProfileImgPath == "" {
		user.ProfileImgPath = ""
	} else {
		user.ProfileImgPath = os.Getenv("STORAGE_IP") + "/" + user.ProfileImgPath
	}

	var userInfo UserInfoRespond
	userInfo.UserDepartment = user
	userInfo.Roles = newRole
	userInfo.AccessControl = roleAccessControl
	return userInfo, nil
}

func (t *userUseCase) UpdateUserInfo(userId int, req requests.UserInfoUpdateReq) error {
	if strings.Contains(req.Firstname, " ") && strings.Contains(req.Lastname, " ") {
		return responses.NewAppErr(http.StatusUnprocessableEntity, "Firstname:have space, Lastname:have space")
	}

	if strings.Contains(req.Firstname, " ") {
		return responses.NewAppErr(http.StatusUnprocessableEntity, "Firstname:have space")
	}

	if strings.Contains(req.Lastname, " ") {
		return responses.NewAppErr(http.StatusUnprocessableEntity, "Lastname:have space")
	}
	var userData models.Users
	imgBase64 := req.ProfileImgPath
	if imgBase64 != "" {
		base64 := ""
		if strings.Contains(imgBase64, ",") {
			dataImgBase64 := helpers.Explode(",", imgBase64)
			base64 = dataImgBase64[1]
		} else {
			base64 = imgBase64
		}
		pathOutput := os.Getenv("PROFILE_PIC")
		fileName := "profile" + time.Now().Format("20060102150405")
		image_data := strings.Split(imgBase64, ",")[0]
		typeFile := strings.Split(strings.Split(image_data, "/")[1], ";")[0]
		filePath, _, err := helpers.DecodeFileBase64FileType(base64, pathOutput, fileName, typeFile, "img")
		if err != nil {
			logs.Error(err)
			return responses.NewAppErr(400, err.Error())
		}
		userData.ProfileImgPath = filePath
	} else {

		userData.ProfileImgPath = ""
	}

	if req.Email != "" {
		users, err := t.userRepo.CheckEmail(userId, req.Email)
		if err != nil {
			logs.Error(err)
			return responses.NewAppErr(400, err.Error())
		}
		if len(users) > 0 {
			return responses.NewAppErr(400, "Email นี้ถูกใช้ในระบบแล้ว")
		}
	}

	userData.Firstname = req.Firstname
	userData.Lastname = req.Lastname
	userData.UpdatedBy = userId
	userData.Tel = req.Tel
	userData.Email = req.Email
	err := t.userRepo.UpdateUser(userId, userData)
	if err != nil {
		logs.Error(err)
		return responses.NewInternalServerError()
	}
	return nil
}

func (t *userUseCase) UserInfoChangePassword(userId int, req requests.UpdatePasswordUserInfoReq) error {

	if req.NewPassword != req.ConfirmNewPassword {
		return responses.NewAppErr(http.StatusBadRequest, "รหัสผ่านไม่ตรงกัน")
	}
	newPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}

	userData := models.Users{
		Password:  string(newPassword),
		UpdatedBy: userId,
	}

	err = t.userRepo.UpdatePassword(userId, userData)
	if err != nil {
		return responses.NewInternalServerError()
	}
	return nil
}

func (t *userUseCase) CheckPassword(userId int, password string) (bool, error) {
	user, err := t.userRepo.CheckPassword(userId)
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println(err)
		return false, responses.NewAppErr(400, err.Error())
	}
	return true, nil
}

func (t *userUseCase) GetUserById(userId int) (interface{}, error) {
	user, err := t.userRepo.GetUserById(userId)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}
	// return user, nil
	// count all users

	// query

	roles, err := t.userRepo.GetRole()
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}
	// var usersDatas []responses.UsersDataRes
	// for _, item := range users {
	userRole, _ := t.userRepo.GetUserRole(int(user.Id))
	ids := []int{}
	for _, item := range userRole {
		ids = append(ids, item.RoleID)
	}

	roleData, _ := t.userRepo.GetRoleByIds(ids)
	temp := make(map[int]string)
	for _, item := range roleData {
		temp[item.Id] = item.Name
	}

	var usersData responses.UsersDataRes
	usersData.CreatedBy = user.CreatedBy
	// usersData.DepartmentId = user.DepartmentId
	usersData.Email = user.Email
	usersData.Firstname = user.Firstname
	usersData.Lastname = user.Lastname
	usersData.Id = user.Id
	usersData.RefUserOwnerID = user.RefUserOwnerID
	usersData.RefUserOwner = user.RefUserOwner
	usersData.RefDepotID = user.RefDepotID
	usersData.RefDepot = user.RefDepot
	if user.ProfileImgPath == "" {
		usersData.ProfileImgPath = ""
	} else {
		usersData.ProfileImgPath = os.Getenv("STORAGE_IP") + "/" + user.ProfileImgPath
	}
	usersData.Status = user.Status
	usersData.Tel = user.Tel
	usersData.UpdatedBy = user.UpdatedBy
	usersData.Username = user.Username
	// usersData.Roles = role
	rolesRes := []responses.RoleRes{}
	for _, item := range roles {
		role := responses.RoleRes{}
		role.Id = item.Id
		role.Name = item.Name
		_, isval := temp[item.Id]
		if isval {
			role.IsChecked = true
		} else {
			role.IsChecked = false
		}
		rolesRes = append(rolesRes, role)
	}
	// usersData.Department = user.Department
	usersData.Roles = rolesRes
	return usersData, nil
	// }

}

func (t *userUseCase) GetUserByUsername(username string) (interface{}, error) {
	user, err := t.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}
	return user, nil
}

func (t *userUseCase) DeleteUserById(userId int) error {
	err := t.userRepo.DeleteUserById(userId)
	if err != nil {
		return responses.NewAppErr(400, err.Error())
	}
	return nil
}

func (t *userUseCase) CreateUser(req requests.UserReq) (interface{}, error) {
	var userData models.Users
	username := req.Username
	user, err := t.userRepo.GetUserByUsername(username)
	if err != nil {
		return userData, err
	}
	// check have username
	if len(user) > 0 {
		return userData, responses.NewAppErr(400, "ชื่อผู้ใช้งานนี้ถูกใช้ในระบบแล้ว") //
	}

	if req.Email != "" {
		users, err := t.userRepo.CheckEmail(0, req.Email)
		if err != nil {
			return userData, responses.NewAppErr(400, err.Error())
		}
		if len(users) > 0 {
			return userData, responses.NewAppErr(400, "Email นี้ถูกใช้ในระบบแล้ว")
		}
	}

	password, bcryptErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if bcryptErr != nil {
		return userData, responses.NewAppErr(400, bcryptErr.Error()) //errors.New("password error")
	}
	if req.ProfileImgPath != "" {
		imgBase64 := req.ProfileImgPath
		base64 := ""
		if strings.Contains(imgBase64, ",") {
			dataImgBase64 := helpers.Explode(",", imgBase64)
			base64 = dataImgBase64[1]
		} else {
			base64 = imgBase64
		}
		pathOutput := os.Getenv("PROFILE_PIC")
		fileName := "profile" + time.Now().Format("20060102150405")
		image_data := strings.Split(imgBase64, ",")[0]
		typeFile := strings.Split(strings.Split(image_data, "/")[1], ";")[0]
		filePath, _, err := helpers.DecodeFileBase64FileType(base64, pathOutput, fileName, typeFile, "img")
		if err != nil {
			return "", responses.NewAppErr(400, err.Error())
		}
		userData.ProfileImgPath = filePath
	} else {
		userData.ProfileImgPath = ""
	}

	userData.Username = req.Username
	userData.Email = req.Email
	userData.Password = string(password)
	userData.RefUserOwnerID = req.RefUserOwnerID
	userData.RefDepotID = req.RefDepotID
	userData.Firstname = req.Firstname
	userData.Lastname = req.Lastname
	userData.Status = req.Status
	userData.CreatedBy = req.CreatedBy
	userData.Tel = req.Tel

	user2, _ := t.userRepo.CreateUser(userData)
	if err != nil {
		return userData, responses.NewAppErr(400, err.Error())
	}

	t.userRepo.CreateUserRole(int(user2.Id), req.Roles)
	return userData, nil
}

func (t *userUseCase) UpdateUser(userId int, req requests.UserReq) (interface{}, error) {
	user1, err := t.userRepo.GetUserById(userId)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	if req.Email != "" {
		users, err := t.userRepo.CheckEmail(userId, req.Email)
		if err != nil {
			return "", responses.NewAppErr(400, err.Error())
		}
		if len(users) > 0 {
			return "", responses.NewAppErr(400, "Email นี้ถูกใช้ในระบบแล้ว")
		}
	}

	var userData models.Users
	if req.Password != "" {
		password, bcryptErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if bcryptErr != nil {
			return userData, responses.NewAppErr(400, bcryptErr.Error()) //errors.New("password error")
		}
		userData.Password = string(password)
	} else {
		userData.Password = user1.Password
	}

	imgBase64 := req.ProfileImgPath
	if imgBase64 != "" {
		base64 := ""
		if strings.Contains(imgBase64, ",") {
			dataImgBase64 := helpers.Explode(",", imgBase64)
			base64 = dataImgBase64[1]
		} else {
			base64 = imgBase64
		}
		pathOutput := os.Getenv("PROFILE_PIC")
		fileName := "profile" + time.Now().Format("20060102150405")
		image_data := strings.Split(imgBase64, ",")[0]
		typeFile := strings.Split(strings.Split(image_data, "/")[1], ";")[0]
		filePath, _, err := helpers.DecodeFileBase64FileType(base64, pathOutput, fileName, typeFile, "img")
		if err != nil {
			return "", responses.NewAppErr(400, err.Error())
		}
		userData.ProfileImgPath = filePath
	} else {

		userData.ProfileImgPath = ""
	}

	userData.Username = req.Username
	userData.Email = req.Email
	userData.RefUserOwnerID = req.RefUserOwnerID
	userData.RefDepotID = req.RefDepotID
	userData.Firstname = req.Firstname
	userData.Lastname = req.Lastname
	userData.Status = req.Status
	userData.CreatedBy = req.CreatedBy
	userData.Tel = req.Tel
	err = t.userRepo.UpdateUser(userId, userData)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	err = t.userRepo.CreateUserRole(userId, req.Roles)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}
	data, err := t.GetUserById(userId)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}
	return data, nil
}

func createFilter(params requests.UserQueryParams) (string, error) {
	var queryWhere []string
	if params.Username != "" {
		queryWhere = append(queryWhere, "username like '%"+params.Username+"%'")
	}

	if params.Fullname != "" {
		queryWhere = append(queryWhere, "CONCAT(firstname, lastname) LIKE '%"+params.Fullname+"%'")
	}

	if params.Status != "" {
		status := fmt.Sprintf("status = %s", params.Status)
		queryWhere = append(queryWhere, status)
	}

	if params.RefUserOwnerID != "" {
		department := fmt.Sprintf("ref_user_owner_id = %s", params.RefUserOwnerID)
		queryWhere = append(queryWhere, department)

		if params.RefUserOwnerID == "3" {
			department := fmt.Sprintf("ref_depot_id = %s", params.RefDepotID)
			queryWhere = append(queryWhere, department)
		}

	}

	if params.Permission != "" {
		permission := fmt.Sprintf("user_role.role_id = %s", params.Permission)
		queryWhere = append(queryWhere, permission)
	}

	query := strings.Join(queryWhere, " and ")
	helpers.PrintlnJson(query)
	return query, nil
}
