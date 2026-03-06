package models

type RefDirection struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rd *RefDirection) TableName() string {
	return "ref_direction"
}
