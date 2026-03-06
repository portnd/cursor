package models

type RefMaterialSubgrade struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	IsInitial        bool    `json:"is_initial"`
	LayerCoefficient float64 `json:"layer_coefficient"`
	Drainage         float64 `json:"drainage"`
}

func (rmg *RefMaterialSubgrade) TableName() string {
	return "ref_material_subgrade"
}
