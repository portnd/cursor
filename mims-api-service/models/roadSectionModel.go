package models

import "github.com/lib/pq"

type RoadSection struct {
	Id                int            `json:"id"`
	RoadGroupId       int            `json:"road_group_id"`
	Number            string         `json:"number"`
	NameOriginTH      string         `json:"name_origin_th"`
	NameDestinationTH string         `json:"name_destination_th"`
	NameOriginEn      string         `json:"name_origin_en"`
	NameDestinationEn string         `json:"name_destination_en"`
	KmStart           float32        `json:"km_start"`
	KmEnd             float32        `json:"km_end"`
	Distance          float32        `json:"distance"`
	ProvinceCode      pq.StringArray `json:"-" gorm:"type:character[]"`
	Province          pq.StringArray `json:"province" gorm:"type:character[]"`
	RefDivisionCode   string         `json:"-"`
	RefDivision       RefDivision    `json:"ref_division"  gorm:"foreignKey:RefDivisionCode; references:DivisionCode"`
	RefDistrictCode   string         `json:"-"`
	RefDistrict       RefDistrict    `json:"ref_district"  gorm:"foreignKey:RefDistrictCode; references:DistrictCode"`
	RefDepotCode      string         `json:"-"`
	RefDepot          RefDepot       `json:"ref_depot"  gorm:"foreignKey:RefDepotCode; references:DepotCode"`
}

type RoadSectionReport struct {
	Id                int     `json:"id"`
	RoadGroupId       int     `json:"road_group_id"`
	Number            string  `json:"number"`
	NameOriginTH      string  `json:"name_origin_th"`
	NameDestinationTH string  `json:"name_destination_th"`
	NameOriginEn      string  `json:"name_origin_en"`
	NameDestinationEn string  `json:"name_destination_en"`
	KmStart           string  `json:"km_start"`
	KmEnd             string  `json:"km_end"`
	Distance          float32 `json:"distance"`
	StrDistance       string
	ProvinceCode      pq.StringArray `json:"-" gorm:"type:character[]"`
	Province          pq.StringArray `json:"province" gorm:"type:character[]"`
	RefDivisionCode   string         `json:"-"`
	RefDivision       RefDivision    `json:"ref_division"  gorm:"foreignKey:RefDivisionCode; references:DivisionCode"`
	RefDistrictCode   string         `json:"-"`
	RefDistrict       RefDistrict    `json:"ref_district"  gorm:"foreignKey:RefDistrictCode; references:DistrictCode"`
	RefDepotCode      string         `json:"-"`
	RefDepot          RefDepot       `json:"ref_depot"  gorm:"foreignKey:RefDepotCode; references:DepotCode"`
}

// RefProvince       []RefRoadSectionProvince `json:"ref_province"  gorm:"foreignKey:RoadSectionId; references:ID"`
// RefProvince     []RefProvince  `json:"ref_province"  gorm:"foreignKey:ProvinceCode; references:ProvinceCode"`
type RefRoadSectionProvince struct {
	ID            int         `json:"-"`
	RoadSectionId int         `json:"-"`
	ProvinceCode  string      `json:"-"`
	Data          RefProvince `json:"data" gorm:"foreignKey:ProvinceCode; references:ProvinceCode"`
}

type RoadSectionInitData struct {
	Id                int          `json:"id"`
	RoadGroupId       int          `json:"road_group_id"`
	Number            string       `json:"number"`
	NameOriginTH      string       `json:"name_origin"`
	NameDestinationTH string       `json:"name_destination"`
	RefSurfaceId      IntDataArray `json:"ref_surface_id" gorm:"type:integer[]"`
}

type RoadSectionInit struct {
	Id                int        `json:"id"`
	RoadGroupId       int        `json:"road_group_id"`
	Number            string     `json:"number"`
	NameOriginTH      string     `json:"name_origin"`
	NameDestinationTH string     `json:"name_destination"`
	Roads             []RoadInit `json:"roads"  gorm:"ForeignKey:RoadSectionId;AssociationForeignKey:Id"`
}

type RoadSectionPreload struct {
	ID                int            `json:"id"`
	RoadGroupId       int            `json:"road_group_id"`
	Number            string         `json:"number"`
	NameOriginTH      string         `json:"name_origin_th"`
	NameDestinationTH string         `json:"name_destination_th"`
	NameOriginEn      string         `json:"name_origin_en"`
	NameDestinationEn string         `json:"name_destination_en"`
	KmStart           float32        `json:"km_start"`
	KmEnd             float32        `json:"km_end"`
	Distance          float32        `json:"distance"`
	ProvinceCode      pq.StringArray `json:"-" gorm:"type:character[]"`
	Province          pq.StringArray `json:"province" gorm:"type:character[]"`
	RefDivisionCode   string         `json:"-"`
	RefDivision       RefDivision    `json:"ref_division"  gorm:"foreignKey:RefDivisionCode; references:DivisionCode"`
	RefDistrictCode   string         `json:"-"`
	RefDistrict       RefDistrict    `json:"ref_district"  gorm:"foreignKey:RefDistrictCode; references:DistrictCode"`
	RefDepotCode      string         `json:"-"`
	RefDepot          RefDepot       `json:"ref_depot"  gorm:"foreignKey:RefDepotCode; references:DepotCode"`
}

type RoadSectionForDashboardMaintenance struct {
	Id                int            `json:"id"`
	RoadGroupId       int            `json:"road_group_id"`
	Number            string         `json:"number"`
	NameOriginTH      string         `json:"name_origin_th"`
	NameDestinationTH string         `json:"name_destination_th"`
	NameOriginEn      string         `json:"name_origin_en"`
	NameDestinationEn string         `json:"name_destination_en"`
	KmStart           float32        `json:"km_start"`
	KmEnd             float32        `json:"km_end"`
	Distance          float32        `json:"distance"`
	ProvinceCode      pq.StringArray `json:"-" gorm:"type:character[]"`
	Province          pq.StringArray `json:"province" gorm:"type:character[]"`
	RefDivisionCode   string         `json:"-"`
	RefDivision       RefDivision    `json:"ref_division"  gorm:"foreignKey:RefDivisionCode; references:DivisionCode"`
	RefDistrictCode   string         `json:"-"`
	RefDistrict       RefDistrict    `json:"ref_district"  gorm:"foreignKey:RefDistrictCode; references:DistrictCode"`
	RefDepotCode      string         `json:"-"`
	RefDepot          RefDepot       `json:"ref_depot"  gorm:"foreignKey:RefDepotCode; references:DepotCode"`
}

func (b *RoadSectionForDashboardMaintenance) TableName() string {
	return "road_section"
}

func (b *RoadSection) TableName() string {
	return "road_section"
}

func (b *RoadSectionInitData) TableName() string {
	return "road_section"
}

func (b *RefRoadSectionProvince) TableName() string {
	return "ref_road_section_province"
}

func (b *RoadSectionInit) TableName() string {
	return "road_section"
}

func (b *RoadSectionReport) TableName() string {
	return "road_section"
}
