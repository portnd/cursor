package models

type RefMaintenanceAnalysisCondition struct {
	Id                                 int    `json:"id"`
	MaintenanceAnalysisStrategicTypeId int    `json:"maintenance_analysis_strategic_type_id"`
	Name                               string `json:"name"`
}

type MaintenanceAnalysisStrategicBudgetTypePreload struct {
	RefMaintenanceAnalysisCondition
	Target []RefMaintenanceAnalysisTargetPrelaoad `json:"target" gorm:"ForeignKey:MaintenanceAnalysisStrategicBudgetTypeId;AssociationForeignKey:Id"`
}

func (rg *RefMaintenanceAnalysisCondition) TableName() string {
	return "ref_maintenance_analysis_condition"
}
