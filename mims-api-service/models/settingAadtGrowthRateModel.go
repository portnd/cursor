package models

type SettingAadtGrowthRate struct {
	ID          int64   `json:"id"`
	RoadGroupID int32   `json:"road_group_id"`
	GrowthRate  float64 `json:"growth_rate" gorm:"column:r"`
}

func (ract *SettingAadtGrowthRate) TableName() string {
	return "setting_aadt_growth_rate"
}
