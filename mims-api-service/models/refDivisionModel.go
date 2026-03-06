package models

type RefDivision struct {
	Id           int    `json:"id"`
	DivisionCode string `json:"division_code"`
	Name         string `json:"name"`
	NameEn       string `json:"name_en"`
	TheGeom      string `json:"the_geom"`
}
type RefDivisionRes struct {
	Id           int    `json:"id"`
	DivisionCode string `json:"division_code"`
	Name         string `json:"name"`
	NameEn       string `json:"name_en"`
}

type RefDivisionInit struct {
	Id           int    `json:"id"`
	DivisionCode string `json:"division_code"`
	OwnerCodeKey string `json:"owner_code_key"`
	Name         string `json:"name"`
}

type RefDistrictData struct {
	Id           int            `json:"id"`
	DistrictCode string         `json:"district_code"`
	Name         string         `json:"name"`
	DivisionCode string         `json:"-"`
	NameEn       string         `json:"name_en"`
	TheGeom      string         `json:"the_geom"`
	Depots       []RefDepotData `json:"depots" gorm:"foreignKey:DistrictCode;references:DistrictCode"`
}

type RefDivisionList struct {
	RefDivisionInit
	Districts []RefDistrictInit `json:"districts" gorm:"foreignKey:DivisionCode;references:DivisionCode"`
}

func (b *RefDivision) TableName() string {
	return "ref_division"
}

func (b *RefDistrictData) TableName() string {
	return "ref_district"
}

func (b *RefDivisionInit) TableName() string {
	return "ref_division"
}

func (b *RefDivisionRes) TableName() string {
	return "ref_division"
}
