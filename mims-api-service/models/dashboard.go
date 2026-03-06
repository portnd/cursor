package models

import (
	"time"

	"github.com/lib/pq"
)

type RefAssetDashboard struct {
	RefAsset
	RefAssetTables []RefAssetTableDashboard `json:"ref_asset_table" gorm:"ForeignKey:RefAssetID;AssociationForeignKey:ID"`
	// RoadAsset	[]RoadAsset	`json:"road_asset" gorm:"ForeignKey:RefAssetTableId;AssociationForeignKey:ID"`
}

func (b *RefAssetDashboard) TableName() string {
	return "ref_asset"
}

type RefAssetTableDashboard struct {
	RefAssetTable
	RoadAssets []RoadAssetForDashboard `json:"road_asset" gorm:"ForeignKey:RefAssetTableId;AssociationForeignKey:ID"`
}

func (b *RefAssetTableDashboard) TableName() string {
	return "ref_asset_table"
}

type RoadAssetForDashboard struct {
	Id              int              `json:"id" gorm:"column:id"`
	RoadId          int              `json:"road_id"`
	Road            RoadForDashboard `json:"road" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadInfo        RoadInfo         `json:"road_info" gorm:"ForeignKey:RoadId;references:RoadId"`
	RefAssetTableId int              `json:"ref_asset_table_id"`
	Revision        int              `json:"revision"`
	Status          string           `json:"status"`
	IdParent        int              `json:"id_parent"`
	IsExclusiveLock bool             `json:"is_exclusive_lock"`
	RejectReason    string           `json:"reject_reason"`
	IsActiveData    bool             `json:"is_active_data"`
	IdTemp          int              `json:"id_temp"`
	CreatedBy       int              `json:"created_by"`
	CreatedDate     time.Time        `json:"created_date"`
	UpdatedBy       int              `json:"updated_by"`
	UpdatedDate     time.Time        `json:"updated_date"`
	AssetCount      int              `json:"asset_count"`
}

func (b *RoadAssetForDashboard) TableName() string {
	return "road_asset"
}

type RoadForDashboard struct {
	Id           int    `json:"id"`
	Seq          int    `json:"seq"`
	ParentRoadId *int   `json:"parent_road_id"`
	RoadLevel    int    `json:"road_level"`
	RoadCode     string `json:"road_code"`
	IsActive     bool   `json:"is_active"`
	// RefDirectionId int       `json:"ref_direction_id"`
	RoadGroupId   int                  `json:"road_group_id"`
	RoadGroup     RoadGroup            `json:"road_group" gorm:"ForeignKey:RoadGroupId;AssociationForeignKey:Id"`
	IsInit        bool                 `json:"is_init"`
	RoadSectionId int                  `json:"road_section_id"`
	RoadSection   RoadSectionDashboard `json:"road_section" gorm:"ForeignKey:RoadSectionId;AssociationForeignKey:Id"`
	CreatedBy     int                  `json:"created_by"`
	CreatedAt     time.Time            `json:"created_at"`
}

type RoadSectionDashboard struct {
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
	RefDivisionCode   string         `json:"ref_division_code"`
	RefDivision       RefDivision    `json:"ref_division"  gorm:"foreignKey:RefDivisionCode; references:DivisionCode"`
	RefDistrictCode   string         `json:"ref_district_code"`
	RefDistrict       RefDistrict    `json:"ref_district"  gorm:"foreignKey:RefDistrictCode; references:DistrictCode"`
	RefDepotCode      string         `json:"ref_depot_code"`
	RefDepot          RefDepot       `json:"ref_depot"  gorm:"foreignKey:RefDepotCode; references:DepotCode"`
}

func (b *RoadSectionDashboard) TableName() string {
	return "road_section"
}

func (b *RoadForDashboard) TableName() string {
	return "road"
}

type RoadConditionDashboard struct {
	RoadCondition
	Road                 RoadForDashboard               `json:"road" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadInfo             RoadInfo                       `json:"road_info" gorm:"ForeignKey:RoadId;references:RoadId"`
	RoadConditionSurveys []RoadConditionSurveyDashboard ` gorm:"ForeignKey:RoadConditionID;references:ID"`
}

func (rc *RoadConditionDashboard) TableName() string {
	return "road_condition"
}

type RoadRetroReflectivityDashboard struct {
	RoadRetroReflectivity
	Road                        RoadForDashboard                      `json:"road" gorm:"ForeignKey:RoadID;AssociationForeignKey:Id"`
	RoadInfo                    RoadInfo                              `json:"road_info" gorm:"ForeignKey:RoadID;references:RoadId"`
	RoadRetroReflectivityRanges []RoadRetroReflectivityRangeDashboard ` gorm:"ForeignKey:RoadRetroReflectivityID;references:ID"`
}

func (rc *RoadRetroReflectivityDashboard) TableName() string {
	return "road_retro_reflectivity"
}

type PavementSurface struct {
	ID          int           `json:"id"`
	Length      float64       `json:"length"`
	SurfaceData []SurfaceData `json:"surface_data" gorm:"ForeignKey:RoadID;references:ID"`
}

type SurfaceData struct {
	ID          int               `json:"id"`
	RoadID      int               `json:"road_id"`
	Length      float64           `json:"length"`
	SurfaceLane []SurfaceLaneData `json:"surface_lane" gorm:"ForeignKey:RoadSurfaceID;references:ID"`
}

func (rc *SurfaceData) TableName() string {
	return "road_surface"
}

type SurfaceLaneData struct {
	ID            int    `json:"id"`
	RoadSurfaceID int    `json:"road_surface_id"`
	LaneNo        int    `json:"lane_no"`
	SurfaceTyepe  string `json:"surface_tyepe"`
}

func (rc *SurfaceLaneData) TableName() string {
	return "road_surface_lane"
}
