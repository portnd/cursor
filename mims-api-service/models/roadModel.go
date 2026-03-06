package models

import "time"

// Todo ...
type Road struct {
	Id           int    `json:"id"`
	Seq          int    `json:"seq"`
	ParentRoadId *int   `json:"parent_road_id"`
	RoadLevel    int    `json:"road_level"`
	RoadCode     string `json:"road_code"`
	IsActive     bool   `json:"is_active"`
	// RefDirectionId int       `json:"ref_direction_id"`
	RoadGroupId   int       `json:"road_group_id"`
	IsInit        bool      `json:"is_init"`
	RoadSectionId int       `json:"road_section_id"`
	CreatedBy     int       `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}

type RoadData struct {
	Road
	SurveyStatus         bool                            `json:"survey_status"`
	ConditionStatus      bool                            `json:"condition_status"`
	ConditionStatusColor string                          `json:"condition_status_color"`
	RetroStatus          bool                            `json:"retro_status"`
	RetroStatusColor     string                          `json:"retro_status_color"`
	DamageStatus         bool                            `json:"damage_status"`
	DamageStatusColor    string                          `json:"damage_status_color"`
	RoadCondition        []RoadConditionSurveyDate       `json:"-" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	MaintenanceRoad      []MaintenanceRoadProjectEndDate `json:"-" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadInfo             RoadInfoAddData                 `json:"road_info" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RefSurface           RefSurfaceRoad                  `json:"RefSurface" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadSurfaceIcon      []RoadSurfaceIcon               `json:"road_surface_icon" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadGeom             []RoadGeom                      `json:"road_geom" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	ChildRoads           []ChildRoadData                 `json:"child_roads" gorm:"ForeignKey:ParentRoadId;AssociationForeignKey:Id"`
}

type ChildRoadData struct {
	Road
	SurveyStatus         bool                            `json:"survey_status"`
	ConditionStatus      bool                            `json:"condition_status"`
	ConditionStatusColor string                          `json:"condition_status_color"`
	RetroStatus          bool                            `json:"retro_status"`
	RetroStatusColor     string                          `json:"retro_status_color"`
	DamageStatus         bool                            `json:"damage_status"`
	DamageStatusColor    string                          `json:"damage_status_color"`
	RoadCondition        []RoadConditionSurveyDate       `json:"-" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	MaintenanceRoad      []MaintenanceRoadProjectEndDate `json:"-" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadInfo             RoadInfoAddData                 `json:"road_info" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RefSurface           RefSurfaceRoad                  `json:"RefSurface" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadSurfaceIcon      []RoadSurfaceIcon               `json:"road_surface_icon" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadGeom             []RoadGeom                      `json:"road_geom" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	ChildRoads           []ChildRoadData                 `json:"child_roads" gorm:"ForeignKey:ParentRoadId;AssociationForeignKey:Id"`
}

type RoadPreload struct {
	Id            int    `json:"id"`
	Seq           int    `json:"seq"`
	ParentRoadId  int    `json:"parent_road_id"`
	RoadLevel     int    `json:"road_level"`
	RoadCode      string `json:"road_code"`
	IsActive      bool   `json:"is_active"`
	RoadGroupId   int    `json:"road_group_id"`
	RefRoadTypeId int    `json:"ref_road_type_id"`
	CreatedBy     int    `json:"created_by"`
	CreatedDate   string `json:"created_date"`
}

func (b *RoadPreload) TableName() string {
	return "road"
}

type RoadSectionData struct {
	RoadSection
	Roads []RoadData `json:"roads" gorm:"ForeignKey:RoadSectionId;AssociationForeignKey:Id"`
}

type RoadListNew struct {
	RoadGroup
	RefDepot []RefDepotNew `json:"ref_depot"`
}

type RefDepotNew struct {
	RefDepot
	Section []RoadSectionData `json:"section"`
}

type RoadList struct {
	RoadGroup
	Sections []RoadSectionData `json:"sections" gorm:"foreignKey:RoadGroupId;AssociationForeignKey:Id"`
}

type RoadListReport struct {
	RoadGroupReport
	Sections []RoadSectionDataReport `json:"sections" gorm:"foreignKey:RoadGroupId;AssociationForeignKey:Id"`
}

type RoadSectionDataReport struct {
	No int `json:"no"`
	RoadSectionReport
	// Roads []RoadDataReport `json:"roads" gorm:"ForeignKey:RoadSectionId;AssociationForeignKey:Id"`
}

type RoadForDashboardMaintenance struct {
	Id            int                                `json:"id"`
	Seq           int                                `json:"seq"`
	ParentRoadId  *int                               `json:"parent_road_id"`
	RoadLevel     int                                `json:"road_level"`
	RoadCode      string                             `json:"road_code"`
	IsActive      bool                               `json:"is_active"`
	RoadGroupId   int                                `json:"road_group_id"`
	IsInit        bool                               `json:"is_init"`
	RoadSectionId int                                `json:"road_section_id"`
	CreatedBy     int                                `json:"created_by"`
	CreatedAt     time.Time                          `json:"created_at"`
	RoadInfo      RoadInfoForDashboard               `json:"road_info" gorm:"ForeignKey:Id;AssociationForeignKey:RoadID"`
	RoadSection   RoadSectionForDashboardMaintenance `json:"road_section" gorm:"ForeignKey:RoadSectionId;AssociationForeignKey:Id"`
}

func (b *RoadForDashboardMaintenance) TableName() string {
	return "road"
}

// type RoadList struct {

//	type RoadList struct {
//		Road
//		Direction       RefDirection            `json:"direction" gorm:"ForeignKey:RefDirectionId;AssociationForeignKey:Id"`
//		RoadType        RefRoadType             `json:"road_type" gorm:"ForeignKey:RefRoadTypeId;AssociationForeignKey:Id"`
//		RoadInfo        RoadInfo                `json:"road_info" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
//		RoadSurfaceData []RoadSurfaceLaneData   `json:"road_surface_data" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
//		RoadDamage      []RoadDamageForCount    `json:"road_damage" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
//		RoadCondition   []RoadConditionForCount `json:"road_condition" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
//		RoadSurface     []RoadSurfaceForCount   `json:"road_surface" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
//		RoadAsset       []RoadAssetForCount     `json:"road_asset" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
//		RoadGeom        []RoadGeom              `json:"road_geom" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
//		ChildRoad       []ChildRoad             `json:"child_road" gorm:"ForeignKey:ParentRoadId;AssociationForeignKey:Id"`
//	}
type RoadSectionById struct {
	RoadSection
	RoadGroup RoadGroup `json:"road_group" gorm:"foreignKey:Id;references:RoadGroupId"`
}

type RoadById struct {
	Road
	RoadInfo        RoadInfoAddData   `json:"road_info" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadSurfaceIcon []RoadSurfaceIcon `json:"road_surface_icon" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadGeom        []RoadGeom        `json:"road_geom" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadSection     RoadSectionById   `json:"road_section" gorm:"foreignKey:Id;references:RoadSectionId"`
}

// type RoadById struct {
// 	Road
// 	Direction RefDirection `json:"direction" gorm:"ForeignKey:RefDirectionId;AssociationForeignKey:Id"`
// 	RoadType  RefRoadType  `json:"road_type" gorm:"ForeignKey:RefRoadTypeId;AssociationForeignKey:Id"`
// 	RoadInfo  RoadInfo     `json:"road_info" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
// 	RoadGeom  []RoadGeom   `json:"road_geom" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
// }

type ChildRoad struct {
	Road
	Direction       RefDirection            `json:"direction" gorm:"ForeignKey:RefDirectionId;AssociationForeignKey:Id"`
	RoadType        RefRoadType             `json:"road_type" gorm:"ForeignKey:RefRoadTypeId;AssociationForeignKey:Id"`
	RoadInfo        RoadInfo                `json:"road_info" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadSurfaceData []RoadSurfaceLaneData   `json:"road_surface_data" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadDamage      []RoadDamageForCount    `json:"road_damage" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadCondition   []RoadConditionForCount `json:"road_condition" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadSurface     []RoadSurfaceForCount   `json:"road_surface" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadAsset       []RoadAssetForCount     `json:"road_asset" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadGeom        []RoadGeom              `json:"road_geom" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
}

type RoadInfoGeom struct {
	Road
	RoadInfo    RoadInfo        `json:"road_info" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadSection RoadSectionById `json:"road_section" gorm:"foreignKey:Id;references:RoadSectionId"`
	RoadGeom    []RoadGeom      `json:"road_geom" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
}

type RoadLastInfoGeom struct {
	Road
	RoadInfo    RoadInfo    `json:"road_info" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RefRoadType RefRoadType `json:"ref_road_type" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadGeom    []RoadGeom  `json:"road_geom" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
}

type RoadInfoGeomDirection struct {
	Road
	RoadInfo RoadInfoDataDirection `json:"road_info" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	//Direction RefDirection `json:"direction" gorm:"ForeignKey:RefDirectionId;AssociationForeignKey:Id"`
	RoadGeom []RoadGeom `json:"road_geom" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
}

type RoadSurfaceLaneData struct {
	RoadSurfaceForPreload
	RoadSurfaceLane []RoadSurfaceLaneCountLane `json:"road_surface_lane" gorm:"ForeignKey:RoadSurfaceId;AssociationForeignKey:Id"`
}

type RoadInit struct {
	Id             int     `json:"id"`
	Name           string  `json:"name"`
	RoadLevel      int     `json:"road_level"`
	RefRoadTypeId  int     `json:"ref_road_type_id"`
	RoadSectionId  int     `json:"-"`
	RefDirectionId int     `json:"ref_direction_id"`
	LaneTotal      int     `json:"lane_total"`
	KMStart        float64 `json:"km_start"`
	KMEnd          float64 `json:"km_end"`
}

type RoadListInit struct {
	RoadGroupInit
	RoadSection []RoadSectionInit `json:"road_sections" gorm:"ForeignKey:RoadGroupId;AssociationForeignKey:Id"`
}

// type Direction struct {
// 	// Road
// 	RefDirection RefDirection `json:"direction" gorm:"ForeignKey:RefDirectionId;AssociationForeignKey:Id"`
// }

// TableName use to specific table
func (b *Road) TableName() string {
	return "road"
}

func (b *RoadInit) TableName() string {
	return "road"
}

type RoadGroupFilter struct {
	Id         int    `json:"id"`
	RoadNumber string `json:"road_number"`
	ShortName  string `json:"short_name"`
}

func (b *RoadGroupFilter) TableName() string {
	return "road_group"
}

type RoadDivisionFilter struct {
	RoadGroupFilter
	RoadSections []RoadSectionFilter `json:"road_sections" gorm:"ForeignKey:RoadGroupId;AssociationForeignKey:Id"`
}

type RoadSectionFilter struct {
	Id                int                 `json:"id"`
	RoadGroupId       int                 `json:"road_group_id"`
	Number            string              `json:"number"`
	NameOriginTH      string              `json:"name_origin"`
	NameDestinationTH string              `json:"name_destination"`
	RefDistrictCode   string              `json:"ref_district_code"`
	Districts         []RefDistrictFilter `json:"districts" gorm:"foreignKey:DistrictCode;references:RefDistrictCode"`
}

func (b *RoadSectionFilter) TableName() string {
	return "road_section"
}

type RefDivisionFilter struct {
	Id           int    `json:"id"`
	DivisionCode string `json:"division_code"`
	OwnerCodeKey string `json:"owner_code_key"`
	Name         string `json:"name"`
}

func (b *RefDivisionFilter) TableName() string {
	return "ref_division"
}

func (b *RefDistrictFilter) TableName() string {
	return "ref_district"
}

type RefDistrictFilter struct {
	Id           int              `json:"id"`
	DistrictCode string           `json:"district_code"`
	Name         string           `json:"name"`
	OwnerCodeKey string           `json:"owner_code_key"`
	DivisionCode string           `json:"-"`
	Depots       []RefDepotFilter `json:"depots" gorm:"foreignKey:DistrictCode;references:DistrictCode"`
}

func (b *RefDepotFilter) TableName() string {
	return "ref_depot"
}

type RefDepotFilter struct {
	Id           int    `json:"id"`
	DepotCode    string `json:"depot_code"`
	Name         string `json:"name"`
	OwnerCodeKey string `json:"owner_code_key"`
	DistrictCode string `json:"-"`
}
