package requests

type UserReq struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	RefUserOwnerID  int    `json:"ref_user_owner_id"`
	RefDepotID      int    `json:"ref_depot_id"`
	TitleName       string `json:"title_name"`
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	ProfileImgPath  string `json:"profile_img_path"`
	Email           string `json:"email"`
	Status          bool   `json:"status"`
	Tel             string `json:"tel"`
	CreatedBy       int    `json:"created_by"`
	UpdatedBy       int    `json:"updated_by"`
	Roles           []int  `json:"roles"`
}
type UserUpdateReq struct {
	RefUserOwnerID int    `json:"ref_user_owner_id"`
	RefDepotID     int    `json:"ref_depot_id"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	Tel            string `json:"tel"`
	ProfileImgPath string `json:"profile_img_path"`
	// Status         string `json:"status"`
	CreatedBy int `json:"created_by" swaggerignore:"true"`
	UpdatedBy int `json:"updated_by" swaggerignore:"true"`
}

type UserInfoUpdateReq struct {
	Firstname      string `json:"firstname" validate:"nonzero" extensions:"x-order=0"`
	Lastname       string `json:"lastname" validate:"nonzero" extensions:"x-order=1"`
	Email          string `json:"email" extensions:"x-order=2"`
	Tel            string `json:"tel" extensions:"x-order=3"`
	ProfileImgPath string `json:"profile_img_path" extensions:"x-order=4"`
}

type UpdatePasswordUserInfoReq struct {
	CurrentPassword    string `json:"current_password" validate:"nonzero" extensions:"x-order=0"`
	NewPassword        string `json:"new_password" validate:"nonzero" extensions:"x-order=1"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"nonzero" extensions:"x-order=2"`
}

type UserQueryParams struct {
	Page           string `form:"page"`
	Limit          string `form:"limit"`
	Username       string `form:"username"`
	Fullname       string `form:"fullname"`
	RefUserOwnerID string `form:"ref_user_owner_id"`
	RefDepotID     string `form:"ref_depot_id"`
	Permission     string `form:"permission"`
	Status         string `form:"status"`
}
