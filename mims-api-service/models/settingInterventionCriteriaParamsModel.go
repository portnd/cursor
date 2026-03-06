package models

import "time"

// Todo ...
type SettingInterventionCriteriaParams struct {
	Id        int       `json:"id"`
	Params    string    `json:"params"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	IsLatest  bool      `json:"is_latest"`
}

func (b *SettingInterventionCriteriaParams) TableName() string {
	return "setting_intervention_criteria_params"
}
