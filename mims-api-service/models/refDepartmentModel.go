package models

type RefDepartment struct {
	ID        int    `json:"id" extensions:"x-order=0"`
	Name      string `json:"name" extensions:"x-order=1"`
	Status    int    `json:"-"`
	CanDelete bool   `json:"can_delete" extensions:"x-order=2"`
}

func (rd *RefDepartment) TableName() string {
	return "ref_department"
}
