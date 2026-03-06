package models

// Todo ...
type RoleAccessControl struct {
	Id              int `json:"id"`
	RoleId          int `json:"role_id"`
	AccessControlId int `json:"access_control_id"`
}

type RoleAccessJoinControl struct {
	Id              int    `json:"id"`
	RoleId          int    `json:"role_id"`
	AccessControlId int    `json:"access_control_id"`
	AccessKey       string `json:"access_key"`
}

// TableName use to specific table
func (b *RoleAccessControl) TableName() string {
	return "role_access_control"
}

func (b *RoleAccessJoinControl) TableName() string {
	return "role_access_control"
}
