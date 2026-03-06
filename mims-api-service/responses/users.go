package responses

import (
	"gitlab.com/mims-api-service/models"
)

type UserInfo struct {
	models.UserDepartment
	AccessControl []models.AccessControl `json:"access_control"`
}

type Users struct {
	Items        []models.Users `json:"items"`
	CurrentPage  int64          `json:"current_page"`
	NextPage     int64          `json:"next_page"`
	PreviousPage int64          `json:"previous_page"`
	SizePerPage  int64          `json:"size_per_page"`
	TotalPages   int64          `json:"total_pages"`
	TotalItems   int64          `json:"total_items"`
}

type User struct {
	Data models.Users `json:"data"`
}

type UsersDataRes struct {
	Id             uint        `json:"id"`
	Email          string      `json:"email"`
	Username       string      `json:"username"`
	Password       string      `json:"-"`
	RefUserOwnerID int         `json:"ref_user_owner_id"`
	RefUserOwner   interface{} `json:"ref_user_owner"`
	RefDepotID     int         `json:"ref_depot_id"`
	RefDepot       interface{} `json:"ref_depot"`
	Firstname      string      `json:"firstname"`
	Lastname       string      `json:"lastname"`
	ProfileImgPath string      `json:"profile_img_path"`
	Status         bool        `json:"status"`
	Tel            string      `json:"tel"`
	CreatedBy      int         `json:"created_by"`
	UpdatedBy      int         `json:"updated_by"`
	Roles          []RoleRes   `json:"roles"`
}

type RoleRes struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	IsChecked bool   `json:"is_checked"`
}

type UserBy struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	ProfileImgPath string `json:"profile_pic" gorm:"column:profile_img_path"`
	DepartName     string `json:"depart_name"`
}

type RefUserOwner struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	RefDepot interface{} `json:"ref_depot"`
}

type RefDepot struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
