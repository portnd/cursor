package models

// Todo ...
type AccessGroup struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	ParentId  int    `json:"parent_id"`
	Route     string `json:"route"`
	Icon      string `json:"icon"`
	Status    int    `json:"status"`
	AccessKey string `json:"access_key"`
}

// TableName use to specific table
func (b *AccessGroup) TableName() string {
	return "access_group"
}
