package models

type RefDataStatus struct {
	ID             int    `json:"id" example:"1"`
	StatusCode     string `json:"status_code" example:"A"`
	Name           string `json:"name" example:"อนุมัติ"`
	NextActionList string `json:"-"`
	Seq            int    `json:"-"`
}

func (rds *RefDataStatus) TableName() string {
	return "ref_data_status"
}
