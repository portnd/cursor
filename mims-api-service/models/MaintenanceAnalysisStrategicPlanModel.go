package models

import (
	"time"
)

type MaintenanceAnalysisPlan struct {
	ID                    *int      `json:"id"`
	MaintenanceAnalysisID int       `json:"maintenance_analysis_id"`
	Plan1                 *float64  `json:"plan_1" gorm:"column:plan_1"`
	Plan2                 *float64  `json:"plan_2" gorm:"column:plan_2"`
	Plan3                 *float64  `json:"plan_3" gorm:"column:plan_3"`
	PlanYear              float64   `json:"plan_year"`
	CreatedBy             int       `json:"created_by"`
	UpdatedBy             int       `json:"updated_by"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func (b *MaintenanceAnalysisPlan) TableName() string {
	return "maintenance_analysis_plan"
}
