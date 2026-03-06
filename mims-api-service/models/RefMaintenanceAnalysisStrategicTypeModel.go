package models

type MaintenanceAnalysisStrategicType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type MaintenanceAnalysisStrategicTypePreload struct {
	MaintenanceAnalysisStrategicType
	Budget []MaintenanceAnalysisStrategicBudgetTypePreload `json:"budget" gorm:"ForeignKey:MaintenanceAnalysisStrategicTypeId;AssociationForeignKey:Id"`
}

func (rg *MaintenanceAnalysisStrategicType) TableName() string {
	return "ref_maintenance_analysis_strategic_type"
}
