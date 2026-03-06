package requests

import (
	"mime/multipart"
	"time"
)

type RoadPrams struct {
	Keyword       string   `json:"keyword" form:"keyword"`
	RoadId        []int    `json:"road_id" form:"road_id"`
	RoadGroupId   []int    `json:"road_group_id" form:"road_group_id"`
	RoadSectionId []int    `json:"road_section_id" form:"road_section_id"`
	KmStart       *float32 `json:"km_start" form:"km_start"`
	KmEnd         *float32 `json:"km_end" form:"km_end"`
	RefSurfaceId  []int    `json:"ref_road_surface_id" form:"ref_road_surface_id"`
	DepotCode     []string `json:"depot_code" form:"depot_code"`
	IsIri1000     *bool    `json:"is_iri_1000" form:"is_iri_1000"`
	IsIri100      *bool    `json:"is_iri_100" form:"is_iri_100"`
	IsRut100      *bool    `json:"is_rut_100" form:"is_rut_100"`
	IsIfi100      *bool    `json:"is_ifi_100" form:"is_ifi_100"`
	IsG7100       *bool    `json:"is_g7_100" form:"is_g7_100"`
}

type RoadDamageListPrams struct {
	Uuid   string `json:"uuid"`
	RoadId string `json:"road_id"`
}

// type RoadDamageImport struct {
// 	LaneNo         int                   `json:"lane_no" form:"lane_no" validate:"nonzero" extensions:"x-order=0"`
// 	SurveyedDate   string                `json:"surveyed_date" form:"surveyed_date" validate:"nonzero" extensions:"x-order=1"`
// 	DamageFilename *multipart.FileHeader `json:"damage_filename" form:"damage_filename" validate:"nonzero" extensions:"x-order=2"`
// 	ImageFilename  *multipart.FileHeader `json:"image_filename" form:"image_filename"  extensions:"x-order=3"`
// }

type RoadDamageImport struct {
	LaneNo               int    `form:"lane_no" validate:"nonzero" extensions:"x-order=0"`
	SurveyedDate         string `form:"surveyed_date" validate:"nonzero" extensions:"x-order=1"`
	DamageFilenameStatus string `form:"damage_filename_status"`
	ImageFilenameStatus  string `form:"image_filename_status"`
}

// type RoadDamageImportFile struct {
// 	DamageFilename *multipart.FileHeader `json:"damage_filename" form:"damage_filename" validate:"nonzero" extensions:"x-order=2"`
// 	ImageFilename  *multipart.FileHeader `json:"image_filename" form:"image_filename"  extensions:"x-order=3"`
// }

type RoadCondition struct {
	RoadID       int
	LaneNo       int    `form:"lane_no"   extensions:"x-order=2"`
	SurveyedDate string `form:"surveyed_date"  extensions:"x-order=1" `
	Remarks      string `form:"remarks" extensions:"x-order=3"`
}

type RoadRetroReflectivity struct {
	RoadID       int
	LineNo       int    `form:"line_no"   extensions:"x-order=2"`
	SurveyedDate string `form:"surveyed_date"  extensions:"x-order=1" `
	Remarks      string `form:"remarks" extensions:"x-order=3"`
}

type RoadRetroReflectivityFiles struct {
	CsvFilename *multipart.FileHeader `form:"csv_filename" swaggerignore:"true"  extensions:"x-order=5" `
}

type RoadConditionFiles struct {
	IriFilename   *multipart.FileHeader `form:"iri_filename" swaggerignore:"true"  extensions:"x-order=5" `
	ImageFilename *multipart.FileHeader `form:"image_filename" swaggerignore:"true"    extensions:"x-order=7" `
}

type RoadConditionUpdate struct {
	RoadID       int
	LaneNo       int    `form:"lane_no"  extensions:"x-order=2" `
	SurveyedDate string `form:"surveyed_date"  extensions:"x-order=1"`
	Remarks      string `form:"remarks" extensions:"x-order=3"`
}

type RoadConditionCompare struct {
	Years         string `form:"years"`
	Lanes         string `form:"lanes"`
	ConditionType string `form:"condition_type"`
}

type RoadRetroReflectivityCompare struct {
	Years string `form:"years"`
	Lines string `form:"lines"`
}

// type RoadConditionCompareYear struct {
// 	Year  string `form:"years"`
// 	Lanes string `form:"lanes"`
// 	Input string `form:"input"`
// }

type RoadConditionDetails struct {
	ConditionRangeType string `form:"condition_range_type" `
}

type RoadRetroReflectivityDetails struct {
	RangeType        string `form:"range_type" `
	RefStripeTypeIDs string `form:"ref_stripe_type_ids" `
}

type RoadConditionCompareAverage struct {
	Lanes string `form:"lanes"  validate:"nonzero"`
}

type RcData struct {
	SurveyedDate        time.Time `json:"surveyed_date"`
	Year                int       `json:"year"`
	RoadId              int       `json:"road_id"`
	LaneNo              int       `json:"lane_no"`
	IDParent            int       `json:"id_parent"`
	Status              string    `json:"status"`
	Revision            int       `json:"revision"`
	CreatedBy           int       `json:"created_by"`
	CreatedDate         time.Time `json:"created_date"`
	UpdatedBy           int       `json:"updated_by"`
	UpdatedDate         time.Time `json:"updated_date"`
	KmStart             float64   `json:"km_start"`
	KmEnd               float64   `json:"km_end"`
	DamageInputFilepath string    `json:"damage_input_filepath"`
	DamageImgFilepath   string    `json:"damage_img_filepath"`
}

type MItem struct {
	KmStart float64 `json:"km_start"`
	KmEnd   float64 `json:"km_end"`
	TheGeom string  `json:"the_geom"`
}
type RoadDamageItem struct {
	KM                   float64 `json:"km"`
	SurveyType           string  `json:"survey_type"`
	AcIcrack             float64 `json:"ac_icrack"`
	AcUcrack             float64 `json:"ac_ucrack"`
	AcRavelling          float64 `json:"ac_ravelling"`
	AcPatching           float64 `json:"ac_patching"`
	AcPotholeArea        float64 `json:"ac_pothole_area"`
	AcBleeding           float64 `json:"ac_bleeding"`
	AcPotholeCount       float64 `json:"ac_pothole_count"`
	CcTransverseCrack    float64 `json:"cc_transverse_crack"`
	CcNonTransverseCrack float64 `json:"cc_non_transverse_crack"`
	CcCornerBreak        float64 `json:"cc_corner_break"`
	CcJointSealDamage    float64 `json:"cc_joint_seal_damage"`
	CcPatching           float64 `json:"cc_patching"`
	CcSpalling           float64 `json:"cc_spalling"`
	CcScaling            float64 `json:"cc_scaling"`
	ImgFilepath          string  `json:"img_filepath"`
	TheGeomPoint         string  `json:"the_geom"`
}

type RaodRangeItem struct {
	RoaDamageId          int     `json:"road_damage_id"`
	KMStart              float64 `json:"km_start"`
	KMEnd                float64 `json:"km_end"`
	SurveyType           string  `json:"survey_type"`
	AcIcrack             float64 `json:"ac_icrack"`
	AcUcrack             float64 `json:"ac_ucrack"`
	AcRavelling          float64 `json:"ac_ravelling"`
	AcPatching           float64 `json:"ac_patching"`
	AcPotholeArea        float64 `json:"ac_pothole_area"`
	AcBleeding           float64 `json:"ac_bleeding"`
	AcPotholeCount       float64 `json:"ac_pothole_count"`
	CcTransverseCrack    float64 `json:"cc_transverse_crack"`
	CcNonTransverseCrack float64 `json:"cc_non_transverse_crack"`
	CcCornerBreak        float64 `json:"cc_corner_break"`
	CcJointSealDamage    float64 `json:"cc_joint_seal_damage"`
	CcPatching           float64 `json:"cc_patching"`
	CcSpalling           float64 `json:"cc_spalling"`
	CcScaling            float64 `json:"cc_scaling"`
	TheGeom              string  `json:"the_geom"`
}

type RoadDamageMItem struct {
	RoadDamageRangeID    int     `json:"road_damage_range_id"`
	KM                   float64 `json:"km"`
	KMEnd                float64 `json:"km_end"`
	TheGeomPoint         string  `json:"the_geom"`
	ImgFilepath          string  `json:"img_filepath"`
	HashData             string  `json:"hash_data"`
	SurveyType           string  `json:"survey_type"`
	AcIcrack             float64 `json:"ac_icrack"`
	AcUcrack             float64 `json:"ac_ucrack"`
	AcRavelling          float64 `json:"ac_ravelling"`
	AcPatching           float64 `json:"ac_patching"`
	AcPotholeArea        float64 `json:"ac_pothole_area"`
	AcBleeding           float64 `json:"ac_bleeding"`
	AcPotholeCount       float64 `json:"ac_pothole_count"`
	CcTransverseCrack    float64 `json:"cc_transverse_crack"`
	CcNonTransverseCrack float64 `json:"cc_non_transverse_crack"`
	CcCornerBreak        float64 `json:"cc_corner_break"`
	CcJointSealDamage    float64 `json:"cc_joint_seal_damage"`
	CcPatching           float64 `json:"cc_patching"`
	CcSpalling           float64 `json:"cc_spalling"`
	CcScaling            float64 `json:"cc_scaling"`
}

type SumValues struct {
	SurveyType           string  `json:"survey_type"`
	AcIcrack             float64 `json:"ac_icrack"`
	AcUcrack             float64 `json:"ac_ucrack"`
	AcRavelling          float64 `json:"ac_ravelling"`
	AcPatching           float64 `json:"ac_patching"`
	AcPotholeArea        float64 `json:"ac_pothole_area"`
	AcBleeding           float64 `json:"ac_bleeding"`
	AcPotholeCount       float64 `json:"ac_pothole_count"`
	CcTransverseCrack    float64 `json:"cc_transverse_crack"`
	CcNonTransverseCrack float64 `json:"cc_non_transverse_crack"`
	CcCornerBreak        float64 `json:"cc_corner_break"`
	CcJointSealDamage    float64 `json:"cc_joint_seal_damage"`
	CcPatching           float64 `json:"cc_patching"`
	CcSpalling           float64 `json:"cc_spalling"`
	CcScaling            float64 `json:"cc_scaling"`
}

type RoadRangeItems struct {
	MItem
	SurveyType           string  `json:"survey_type"`
	AcIcrack             float64 `json:"ac_icrack"`
	AcUcrack             float64 `json:"ac_ucrack"`
	AcRavelling          float64 `json:"ac_ravelling"`
	AcPatching           float64 `json:"ac_patching"`
	AcPotholeArea        float64 `json:"ac_pothole_area"`
	AcBleeding           float64 `json:"ac_bleeding"`
	AcPotholeCount       float64 `json:"ac_pothole_count"`
	CcTransverseCrack    float64 `json:"cc_transverse_crack"`
	CcNonTransverseCrack float64 `json:"cc_non_transverse_crack"`
	CcCornerBreak        float64 `json:"cc_corner_break"`
	CcJointSealDamage    float64 `json:"cc_joint_seal_damage"`
	CcPatching           float64 `json:"cc_patching"`
	CcSpalling           float64 `json:"cc_spalling"`
	CcScaling            float64 `json:"cc_scaling"`
}

type RoadDamageCsv struct {
	RoadId               int     `json:"road_id"`
	RoadCode             string  `json:"road_code"`
	Name                 string  `json:"name"`
	KM                   float64 `json:"km"`
	SurveyType           string  `json:"survey_type"`
	AcIcrack             float64 `json:"ac_icrack"`
	AcUcrack             float64 `json:"ac_ucrack"`
	AcRavelling          float64 `json:"ac_ravelling"`
	AcPatching           float64 `json:"ac_patching"`
	AcPotholeArea        float64 `json:"ac_pothole_area"`
	AcBleeding           float64 `json:"ac_bleeding"`
	AcPotholeCount       float64 `json:"ac_pothole_count"`
	CcTransverseCrack    float64 `json:"cc_transverse_crack"`
	CcNonTransverseCrack float64 `json:"cc_non_transverse_crack"`
	CcCornerBreak        float64 `json:"cc_corner_break"`
	CcJointSealDamage    float64 `json:"cc_joint_seal_damage"`
	CcPatching           float64 `json:"cc_patching"`
	CcSpalling           float64 `json:"cc_spalling"`
	CcScaling            float64 `json:"cc_scaling"`
	ImgFilepath          string  `json:"img_filepath"`
}
type AssetDetailsQueryParams struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
	// RoadAssetID     int `form:"road_asset_id"`      // 3055
	RefAssetTableID int `form:"ref_asset_table_id"` //8
}

type AssetPermissionQueryParams struct {
	RefAssetTableID int `form:"ref_asset_table_id"` //8
}

type AssetRevisionsQueryParams struct {
	RefAssetTableID int `form:"ref_asset_table_id"` //8
}

type AssetTemplateQueryParams struct {
	RefAssetTableID int    `form:"ref_asset_table_id"` //8
	Action          string `form:"action"`             //8
	AssetObjectID   int    `form:"asset_object_id"`
	// IDParent        int `form:"id_parent"`
}

type AssetCreateRequest struct {
	RefAssetTableID int            `form:"ref_asset_table_id"` //8
	IDParent        int            `form:"id_parent"`
	IDParent_asset  int            `form:"id_parent_asset"`
	SurveyedDate    time.Time      `form:"surveyed_date"`
	Data            AssetTableData `form:"data"`
	Files           []struct {
		Filename   string `form:"filename"`
		FileBase64 string `form:"file_base64"`
	} `form:"files"`
}

type DeleteRoadAssetObject struct {
	IDParentAsset int `form:"id_parent_asset"`
}

type RoadAssetIDParent struct {
	IDParent int `form:"id_parent" validate:"nonzero"`
}
type RoadAssetKm struct {
	Geom string `form:"geom" validate:"nonzero"`
}

// type RoadRangeItem struct {
// 	RoaDamageId          int     `json:"road_damage_id"`
// 	KMStart              float64 `json:"km_start"`
// 	KMEnd                float64 `json:"km_end"`
// 	AcIcrack             float64 `json:"ac_icrack"`
// 	AcUcrack             float64 `json:"ac_ucrack"`
// 	AcRavelling          float64 `json:"ac_ravelling"`
// 	AcPatching           float64 `json:"ac_patching"`
// 	AcPothole            float64 `json:"ac_pothole"`
// 	AcSurfaceDeform      float64 `json:"ac_surface_deform"`
// 	AcBleeding           float64 `json:"ac_bleeding"`
// 	CcTransverseCrack    float64 `json:"cc_transverse_crack"`
// 	CcNonTransverseCrack float64 `json:"cc_non_transverse_crack"`
// 	CcFaulting           float64 `json:"cc_faulting"`
// 	CcSpalling           float64 `json:"cc_spalling"`
// 	CcCornerbreaks       float64 `json:"cc_corner_breaks"`
// 	CcJointSealDamage    float64 `json:"cc_joint_seal_damage"`
// 	CcPatching           float64 `json:"cc_patching"`
// 	CcFaultingMax        float64 `json:"cc_faulting_max"`
// 	CcScaling            float64 `json:"cc_scaling"`
// 	// ImgFilepath          string  `json:"img_filepath"`
// 	TheGeom string `json:"the_geom"`
// }
