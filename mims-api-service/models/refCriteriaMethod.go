package models

type RefCriteriaMethod struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surface string `json:"surface"`
	Color   string `json:"color"`
}

type RefCriteriaMethodPreLoad struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

type RefCriteriaMethodData struct {
	RefCriteriaMethodPreLoad
	Children []InterventionCriteriaChildren `json:"children" gorm:"ForeignKey:MaintenanceMethod;AssociationForeignKey:Id"`
}

func (ral *RefCriteriaMethod) TableName() string {
	return "ref_criteria_method"
}

func (ral *RefCriteriaMethodPreLoad) TableName() string {
	return "ref_criteria_method"
}
