package models

type SettingAadtParams struct {
	ID       int    `json:"id"`
	Params   string `json:"params"`
	IsLatest string `json:"is_latest"`
}

func (ract *SettingAadtParams) TableName() string {
	return "setting_aadt_params"
}
