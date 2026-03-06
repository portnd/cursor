package models

import "time"

type SettingOptimization struct {
	Id        int       `json:"id"`
	Params    string    `json:"params"`
	IsLatest  bool      `json:"is_latest"`
	IsDeleted bool      `json:"is_deleted"`
	UpdatedBy int       `json:"updated_by"`
	CreatedBy int       `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Optimization struct {
	BcRatioConstraint *float64 `json:"bc_ratio_constraint"`
	DefaultDesignLife *float64 `json:"default_design_life"`
}

func (b *SettingOptimization) TableName() string {
	return "setting_optimization_params"
}
