package models

type RefProvince struct {
	ID           int    `json:"id"`
	ProvinceCode string `json:"province_code"`
	Node         int    `json:"node"`
	Name         string `json:"name_th"`
	NameEn       string `json:"name_en"`
	Region       string `json:"region"`
	RegionName   string `json:"region_name"`
	TheGeom      string `json:"the_geom"`
	Status       bool   `json:"status"`
}

func (b *RefProvince) TableName() string {
	return "ref_province"
}
