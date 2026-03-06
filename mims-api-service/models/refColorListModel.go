package models

type RefColorList struct {
	Name string `json:"ref_color_name"`
	Code string `json:"ref_color_code"`
}

func (rat *RefColorList) TableName() string {
	return "ref_color_list"
}
