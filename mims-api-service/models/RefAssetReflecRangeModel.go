package models

type RefAssetReflecRange struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (ra *RefAssetReflecRange) TableName() string {
	return "ref_reflectivity_range"
}
