package models

import "time"

type SettingBudgetMethod struct {
	Id           int       `json:"id"`
	MethodName   string    `json:"method_name"`
	BudgetId     int       `json:"budget_id"`
	CostPerUnit  *float64  `json:"cost_per_unit"`
	IsShowMethod bool      `json:"is_show_method"`
	IsDeleted    bool      `json:"is_deleted"`
	UpdatedBy    int       `json:"updated_by"`
	CreatedBy    int       `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type SettingBudgetMethodData struct {
	Id         int    `json:"id"`
	MethodName string `json:"method_name"`
}

func (ract *SettingBudgetMethod) TableName() string {
	return "setting_budget_method"
}

func (ract *SettingBudgetMethodData) TableName() string {
	return "setting_budget_method"
}
