package models

type RefMaintenanceAnalysisTarget struct {
	ID                                       int    `json:"id"`
	MaintenanceAnalysisStrategicBudgetTypeId int    `json:"maintenance_analysis_strategic_budget_type_id"`
	Name                                     string `json:"name"`
}

type RefMaintenanceAnalysisTargetPrelaoad struct {
	RefMaintenanceAnalysisTarget
}

func (rg *RefMaintenanceAnalysisTarget) TableName() string {
	return "ref_maintenance_analysis_target"
}
