package models

// Todo ...
type Users struct {
	Id       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
	// DepartmentId   int    `json:"department_id"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	ProfileImgPath string `json:"profile_img_path" gorm:"column:profile_img_path"`
	Status         bool   `json:"status"`
	Tel            string `json:"tel"`
	RefUserOwnerID int    `json:"ref_user_owner_id"`
	RefDepotID     int    `json:"ref_depot_id"`
	CreatedBy      int    `json:"created_by"`
	UpdatedBy      int    `json:"updated_by"`
}

type UsersPassword struct {
	Password string `json:"password"`
}

type UserDepartment struct {
	Users
	RefUserOwner RefUserOwner `json:"ref_user_owner" gorm:"ForeignKey:RefUserOwnerID;AssociationForeignKey:ID"`
	RefDepot     RefDepot     `json:"ref_depot" gorm:"ForeignKey:RefDepotID;AssociationForeignKey:ID"`
}

type UserRes struct {
	Id        uint   `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type UserBy struct {
	Id             int           `json:"id"`
	Name           string        `json:"name"`
	ProfileImgPath string        `json:"profile_pic" gorm:"column:profile_img_path"`
	DepartmentId   int           `json:"-"`
	Department     RefDepartment `json:"department" gorm:"ForeignKey:DepartmentId;AssociationForeignKey:Id"`
}

// TableName use to specific table
func (b *Users) TableName() string {
	return "users"
}

// TableName use to specific table
func (b *UsersPassword) TableName() string {
	return "users"
}

func (b *UserRes) TableName() string {
	return "users"
}

func (b *UserBy) TableName() string {
	return "users"
}
