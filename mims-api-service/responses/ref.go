package responses

import (
	"gitlab.com/mims-api-service/models"
)

// these response have been only created for using in swagger api spec.
type InitData struct {
	RefAsset                  []models.RefAsset                  `json:"ref_asset"`
	RefAssetPosition          []models.RefAssetPosition          `json:"ref_asset_position"`
	RefAssetArea              []models.RefAssetArea              `json:"ref_asset_area"`
	RefAssetGuardrail         []models.RefAssetGuardrail         `json:"ref_asset_guardrail"`
	RefAssetReflecType        []models.RefAssetReflecType        `json:"ref_asset_reflec_type"`
	RefAssetOwner             []models.RefAssetOwner             `json:"ref_asset_owner"`
	RefAssetLightType         []models.RefAssetLightType         `json:"ref_asset_light_type"`
	RefAssetLightWatt         []models.RefAssetLightWatt         `json:"ref_asset_light_watt"`
	RefAssetCleranceType      []models.RefAssetCleranceType      `json:"ref_asset_clerance_type"`
	RefAssetCrashcushionType  []models.RefAssetCrashcushionType  `json:"ref_asset_crashcushion_type"`
	RefAssetFenceType         []models.RefAssetFenceType         `json:"ref_asset_fence_type"`
	RefAssetKmstoneType       []models.RefAssetKmstoneType       `json:"ref_asset_kmstone_type"`
	RefAssetQutterType        []models.RefAssetQutterType        `json:"ref_asset_qutter_type"`
	RefAssetTrafficCameraType []models.RefAssetTrafficCameraType `json:"ref_asset_traffic_camera_type"`
	RefAssetWeightStationType []models.RefAssetWeightStationType `json:"ref_asset_weight_station_type"`
	RefAssetBuildingType      []models.RefAssetBuildingType      `json:"ref_asset_building_type"`
	RefAssetNoiseBarrier      []models.RefAssetNoiseBarrier      `json:"ref_asset_noise_barrier"`
	RefAssetSignImage         []models.RefAssetSignImage         `json:"ref_asset_sign_image"`
	RefAssetTable             []models.RefAssetTable             `json:"ref_asset_table"`
	RefAssetTableColumns      []models.RefAssetTableColumns      `json:"ref_asset_table_columns"`
	RefAssetTableStaff        []models.RefAssetTableStaff        `json:"ref_asset_table_staff"`

	/////////
	RefDataStaus                    []models.RefDataStatus        `json:"ref_data_status"`
	RefDepartment                   []models.RefDepartment        `json:"ref_department"`
	RefDirection                    []models.RefDirection         `json:"ref_direction"`
	RefGrade                        []models.RefGrade             `json:"ref_grade"`
	RefMaterialBase                 []models.RefMaterialBase      `json:"ref_material_base"`
	RefMaterialSubbase              []models.RefMaterialSubbase   `json:"ref_material_subbase"`
	RefMaterialSubgrade             []models.RefMaterialSubgrade  `json:"ref_material_subgrade"`
	RefOwner                        []models.RefOwner             `json:"ref_owner"`
	RefRoadType                     []models.RefRoadType          `json:"ref_road_type"`
	RefSurface                      []models.RefSurface           `json:"ref_surface"`
	RefSurfaceType                  []models.RefSurfaceType       `json:"ref_surface_type"`
	RefStructureSurface             []models.RefStructureSurface  `json:"ref_structure_surface"`
	RefSurfaceGroup                 []RefSurfaceGroup             `json:"ref_surface_group"`
	RefColorList                    []models.RefColorList         `json:"ref_color_list"`
	RefRoadTypeIcon                 []models.RefRoadTypeIcon      `json:"ref_road_type_icon"`
	RefTableList                    []models.RefTableList         `json:"ref_table_list"`
	ConditionGrade                  []ConditionRespondInit        `json:"condition_grade"`
	RoadLineList                    []RoadLineListInit            `json:"reflectivity_grade"`
	RefConditionRange               []models.RefConditionRange    `json:"ref_condition_range"`
	RefReflectivityRange            []models.RefReflectivityRange `json:"ref_reflectivity_range"`
	RoadGroup                       []RoadGroupInitData           `json:"road_group"`
	RoadSection                     []models.RoadSectionInitData  `json:"road_section"`
	RefDistrict                     []models.RefDistrictInitData  `json:"ref_district"`
	RefRoadTypeLevelOne             []models.RefRoadTypeInit      `json:"ref_road_type_level_first"`
	RefRoadTypeLevelTwo             []models.RefRoadTypeInit      `json:"ref_road_type_level_second"`
	RefStripeType                   []models.RefStripeType        `json:"ref_stripe_type"`
	RefStripeColor                  []models.RefStripeColor       `json:"ref_stripe_color"`
	RefDivision                     interface{}                   `json:"ref_division"`
	RefDivisionDashboardAsset       interface{}                   `json:"ref_division_dashboard_asset"`
	RefDivisionDashboardCondition   interface{}                   `json:"ref_division_dashboard_condition"`
	RefDivisionDashboardSurface     interface{}                   `json:"ref_division_dashboard_surface"`
	RefDivisionDashboardMaintenance interface{}                   `json:"ref_division_dashboard_maintenance"`

	RefCriteriaType   []models.RefCriteriaType   `json:"ref_criteria_type"`
	RefCriteriaMethod []models.RefCriteriaMethod `json:"ref_criteria_method"`
	RefUserOwner      []RefUserOwner             `json:"ref_user_owner"`
}

// these structs below use for building example responses at swagger (API document)
// TODO: Need to remove these structs below because we use data{...} in swagger instead of create struct for example.
// first group is defined ref_asset tables which have only two columns including id and name.
type FirstGroup struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"ex: ป้ายจราจร"`
}

// first group is defined ref_asset tables which have 3 columns including id, namd, seq.
type SecondGroup struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"ex: ป้ายจราจร"`
	Seq  int    `json:"seq" example:"1"`
}

type Material struct {
	ID        int    `json:"id" example:"1"`
	Name      string `json:"name" example:"ดิน"`
	IsInitial bool   `json:"true" example:"true"`
}

type MaterialResponse struct {
	Status bool       `json:"status" example:"true"`
	Data   []Material `json:"data"`
}

type FirstGroupResponse struct {
	Status bool         `json:"status" example:"true"`
	Data   []FirstGroup `json:"data"`
}

type SecondGroupResponse struct {
	Status bool          `json:"status" example:"true"`
	Data   []SecondGroup `json:"data"`
}

type RefAsestSignImageReponse struct {
	Status bool                       `json:"status" example:"true"`
	Data   []models.RefAssetSignImage `json:"data"`
}

type RefAssetTableResponse struct {
	Status bool                   `json:"status" example:"true"`
	Data   []models.RefAssetTable `json:"data"`
}

type RefAssetTableColumnsResponse struct {
	Status bool                          `json:"status" example:"true"`
	Data   []models.RefAssetTableColumns `json:"data"`
}

type RefAssetTableStaffResponse struct {
	Status bool                        `json:"status" example:"true"`
	Data   []models.RefAssetTableStaff `json:"data"`
}

type RefDataStatusResponse struct {
	Status bool                   `json:"status" example:"true"`
	Data   []models.RefDataStatus `json:"data"`
}

type RefOwnerResponse struct {
	Status bool              `json:"status" example:"true"`
	Data   []models.RefOwner `json:"data"`
}

type RefAsset struct {
	ID        int    `json:"id" example:"1"`
	Name      string `json:"name" example:"ป้ายจราจร"`
	CanDelete bool   `json:"can_delete" example:"false"`
}

type RefAssetResponse struct {
	Status bool       `json:"status" example:"true"`
	Data   []RefAsset `json:"data"`
}

type RefTableList struct {
	RefName string `json:"ref_name"`
	RefDesc string `json:"ref_desc"`
}

type RefTableListResponse struct {
	Status bool           `json:"status" example:"true"`
	Data   []RefTableList `json:"data"`
}

type InitDataResponse struct {
	Status bool     `json:"status" example:"true"`
	Data   InitData `json:"data"`
}

type RefSurfaceGroup struct {
	Name      string `json:"name"`
	ColorCode string `json:"color_code"`
}

type ConditionGradeResponse struct {
	ID              int                `json:"id"`
	Name            string             `json:"name"`
	ConditionGroups []ConditionListNew `json:"condition_groups"`
}
type ConditionGroup struct {
	ConditionType string           `json:"condition_type"`
	Conditions    []ConditionGrade `json:"conditions"`
}

type ConditionGrade struct {
	Grade            models.RefGrade `json:"grade"`
	LeftValueAC      float64         `gorm:"column:left_value_ac" json:"left_value_ac"`
	LeftConditionAC  string          `gorm:"column:left_condition_ac" json:"left_condition_ac"`
	RightValueAC     float64         `gorm:"column:right_value_ac" json:"right_value_ac"`
	RightConditionAC string          `gorm:"column:right_condition_ac" json:"right_condition_ac"`
	LeftValueCC      float64         `gorm:"column:left_value_cc" json:"left_value_cc"`
	LeftConditionCC  string          `gorm:"column:left_condition_cc" json:"left_condition_cc"`
	RightValueCC     float64         `gorm:"column:right_value_cc" json:"right_value_cc"`
	RightConditionCC string          `gorm:"column:right_condition_cc" json:"right_condition_cc"`
}

type RefAadtParameterVehicleType struct {
	Id             int     `json:"id"`
	NumWheel       int     `json:"num_wheel"`
	Name           string  `json:"name"`
	NumAxle        int     `json:"num_axle"`
	LoadEquivalent float64 `json:"load_equivalent"`
	ImagePath      string  `json:"image_path"`
}

type RoadGroupInitData struct {
	Id               int      `json:"id"`
	Number           string   `json:"number"`
	ShortName        string   `json:"short_name"`
	RefDivisionCodes []string `json:"ref_division_codes"`
	RefDistrictCodes []string `json:"ref_district_codes"`
}
