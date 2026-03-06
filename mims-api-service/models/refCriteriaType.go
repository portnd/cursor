package models

type RefCriteriaType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (ral *RefCriteriaType) TableName() string {
	return "ref_criteria_type"
}
