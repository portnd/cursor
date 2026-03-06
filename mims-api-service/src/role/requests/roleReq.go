package requests

type AccessControl struct {
	AccessControlId int `json:"access_control_id"`
}

type ReqRoleAccCtrl struct {
	Name          string          `json:"name" validate:"min=1" extensions:"x-order=0"`
	AccessControl []AccessControl `json:"access_control" extensions:"x-order=1"`
}

type ReqRoleCreate struct {
	Name string `json:"name" validate:"min=1"`
}

type RoleQueryParams struct {
	Page  string `form:"page"`
	Limit string `form:"limit"`
	Name  string `form:"name"`
}
