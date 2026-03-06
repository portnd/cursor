package models

type RefGrade struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (rg *RefGrade) TableName() string {
	return "ref_grade"
}
