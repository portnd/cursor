package models

// Todo ...
type Department struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// TableName use to specific table
func (b *Department) TableName() string {
	return "department"
}
