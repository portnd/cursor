package models

type RefMaterialSubbase struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	IsInitial        bool    `json:"is_initial"`
	LayerCoefficient float64 `json:"layer_coefficient"`
	Drainage         float64 `json:"drainage"`
}

func (rms *RefMaterialSubbase) TableName() string {
	return "ref_material_subbase"
}
