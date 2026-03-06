package models

// Todo ...
type AccessControl struct {
	Id                int    `json:"id"`
	AccessTitle       string `json:"access_title"`
	AccessDesc        string `json:"access_desc"`
	AccessGrpId       int    `json:"access_grp_id"`
	AccessGrpParentId int    `json:"access_grp_parent_id"`
	AccessKey         string `json:"access_key"`
}

// TableName use to specific table
func (b *AccessControl) TableName() string {
	return "access_control"
}
