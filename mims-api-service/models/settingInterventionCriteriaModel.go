package models

import "time"

type SettingInterventionCriteria struct {
	Id        int       `json:"id"`
	Params    string    `json:"params"`
	CreatedBy int       `json:"-"`
	CreatedAt time.Time `json:"-"`
	IsLatest  bool      `json:"is_latest"`
}

type InterventionCriteriaCount struct {
	Count int `json:"count"`
}

type InterventionCriteria struct {
	Id                             int       `json:"id"`
	MaintenanceMethod              int       `json:"maintenance_method"`
	MaintenanceCostPerUnit         float64   `json:"maintenance_cost_per_unit"`
	MaintenanceDescription         string    `json:"maintenance_description"`
	MaintenanceScraping            float64   `json:"maintenance_scraping"`
	MaintenanceSequence            int       `json:"maintenance_sequence"`
	MaintenanceStandardName        string    `json:"maintenance_standard_name"`
	MaintenanceSurfaceTypeId       int       `json:"maintenance_surface_type_id"`
	MaintenanceSurfaceTypeIdParams int       `json:"maintenance_surface_type_id_params"`
	MaintenanceThickness           float64   `json:"maintenance_thickness"`
	IsDeleted                      bool      `json:"is_deleted"`
	IsShow                         bool      `json:"is_show"`
	UpdatedBy                      int       `json:"updated_by"`
	CreatedBy                      int       `json:"created_by"`
	UpdatedAt                      time.Time `json:"updated_at"`
	CreatedAt                      time.Time `json:"created_at"`
}

type InterventionCriteriaData struct {
	Id                      int    `json:"id"`
	MaintenanceStandardName string `json:"maintenance_standard_name"`
}

type InterventionCriteriaCondition struct {
	Id                     int       `json:"id"`
	InterventionCriteriaId int       `json:"intervention_criteria_id"`
	ConditionSequence      int       `json:"condition_sequence"`
	ConditionCriterion     string    `json:"condition_criterion"`
	ConditionLink          string    `json:"condition_link"`
	ConditionOperation_1   string    `json:"condition_operation_1"`
	ConditionOperation_2   string    `json:"condition_operation_2"`
	ConditionValue_1       float64   `json:"condition_value_1"`
	ConditionValue_2       float64   `json:"condition_value_2"`
	IsDeleted              bool      `json:"is_deleted"`
	UpdatedBy              int       `json:"updated_by"`
	CreatedBy              int       `json:"created_by"`
	UpdatedAt              time.Time `json:"updated_at"`
	CreatedAt              time.Time `json:"created_at"`
}

type InterventionCriteriaChildren struct {
	ID                int    `json:"id"`
	Label             string `json:"label"`
	MaintenanceMethod int    `json:"-"`
}

func (b *InterventionCriteria) TableName() string {
	return "setting_intervention_criteria"
}

func (b *InterventionCriteriaData) TableName() string {
	return "setting_intervention_criteria"
}

func (b *InterventionCriteriaCondition) TableName() string {
	return "setting_intervention_criteria_condition"
}

func (b *SettingInterventionCriteria) TableName() string {
	return "setting_intervention_criteria"
}

func (b *InterventionCriteriaChildren) TableName() string {
	return "setting_intervention_criteria"
}
