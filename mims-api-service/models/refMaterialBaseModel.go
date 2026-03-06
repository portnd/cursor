package models

type RefMaterialBase struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	IsInitial        bool    `json:"is_initial"`
	LayerCoefficient float64 `json:"layer_coefficient"`
	Drainage         float64 `json:"drainage"`
	Type             string  `json:"type"`
}

func (rmb *RefMaterialBase) TableName() string {
	return "ref_material_base"
}
