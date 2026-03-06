package models

import "time"

type SettingBudget struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CanDelete bool      `json:"can_delete"`
	IsDeleted bool      `json:"is_deleted"`
	UpdatedBy int       `json:"updated_by"`
	CreatedBy int       `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type SettingBudgetProload struct {
	SettingBudget
	BudgetMethods []SettingBudgetMethod `json:"budget_methods" gorm:"ForeignKey:BudgetId;AssociationForeignKey:Id"`
}

type SettingBudgetData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (ract *SettingBudget) TableName() string {
	return "setting_budget"
}

func (ract *SettingBudgetData) TableName() string {
	return "setting_budget"
}
