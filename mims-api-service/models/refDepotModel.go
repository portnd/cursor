package models

type RefDepot struct {
	Id        int    `json:"id"`
	DepotCode string `json:"depot_code"`
	Name      string `json:"name"`
	TheGeom   string `json:"-"`
}
type RefDepotRes struct {
	Id        int    `json:"id"`
	DepotCode string `json:"depot_code"`
	Name      string `json:"name"`
}

type RefDepotData struct {
	Id           int    `json:"id"`
	DepotCode    string `json:"depot_code"`
	Name         string `json:"name"`
	DistrictCode string `json:"-"`
	TheGeom      string `json:"the_geom"`
}

type RefDepotInit struct {
	Id           int    `json:"id"`
	DepotCode    string `json:"depot_code"`
	Name         string `json:"name"`
	OwnerCodeKey string `json:"owner_code_key"`
	DistrictCode string `json:"-"`
}
type RefDepotInitData struct {
	Id            int          `json:"id"`
	DepotCode     string       `json:"depot_code"`
	Name          string       `json:"name"`
	DistrictCode  string       `json:"-"`
	RoadSectionID IntDataArray `json:"road_section_id" gorm:"type:integer[]"`
}

func (b *RefDepot) TableName() string {
	return "ref_depot"
}

func (b *RefDepotInitData) TableName() string {
	return "ref_depot"
}

func (b *RefDepotInit) TableName() string {
	return "ref_depot"
}

func (b *RefDepotData) TableName() string {
	return "ref_depot"
}

func (b *RefDepotRes) TableName() string {
	return "ref_depot"
}
