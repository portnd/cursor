package models

type RefStructureSurface struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IsInitial bool   `json:"is_initial"`
}

func (rts *RefStructureSurface) TableName() string {
	return "ref_structure_surface"
}
