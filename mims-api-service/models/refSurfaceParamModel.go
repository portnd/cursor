package models

import "time"

type RefSurfaceParam struct {
	ID           int       `json:"id"`
	RefSurfaceID int       `json:"ref_surface_id" gorm:"column:ref_surface_id"`
	Params       string    `json:"params" gorm:"column:params"`
	CreateBy     int       `json:"create_by" gorm:"column:create_by"`
	CreateDate   time.Time `json:"create_date" gorm:"column:create_date"`
	IsLatest     bool      `json:"is_latest" gorm:"column:is_latest"`
}

func (rs *RefSurfaceParam) TableName() string {
	return "ref_surface_param"
}
