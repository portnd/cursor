package responses

import (
	"time"

	"gitlab.com/mims-api-service/models"
)

type ChangeStatusResponse struct {
	Status        models.RefDataStatus `json:"status" extensions:"x-order=0"`
	ModifiedLabel string               `json:"modified_label" extensions:"x-order=1"`
	RoadInfo      RoadInfoData         `json:"road_info" extensions:"x-order=2"`
	IdParent      int                  `json:"id_parent" extensions:"x-order=3"`
	Category      string               `json:"category" extensions:"x-order=4"`
	AssetId       int                  `json:"asset_id" extensions:"x-order=5"`
	UpdatedDate   time.Time            `json:"updated_date" extensions:"x-order=6"`
}

type RoadInfoData struct {
	RoadId    int                 `json:"road_id" extensions:"x-order=0"`
	Code      string              `json:"code" extensions:"x-order=1"`
	Name      string              `json:"name" extensions:"x-order=2"`
	Direction models.RefDirection `json:"direction" extensions:"x-order=3"`
	RoadType  models.RefRoadType  `json:"road_type" extensions:"x-order=4"`
	KmStart   float32             `json:"km_start" extensions:"x-order=5"`
	KmEnd     float32             `json:"km_end" extensions:"x-order=6"`
}

type ChangeAssetDetailResponse struct {
	IconFilePath  string                    `json:"icon_filepath" extensions:"x-order=0"`
	LineColorCode string                    `json:"line_color_code" extensions:"x-order=1"`
	DataColumns   []map[string]ColumnDetail `json:"data_columns" extensions:"x-order=2"`
	DataChanged   DataDetail                `json:"data_changed" extensions:"x-order=3"`
	DataCurrent   DataDetail                `json:"data_current" extensions:"x-order=4"`
}

// type DataDetail struct {
// 	UpdatedBy   UpdatedBy   `json:"updated_by" extensions:"x-order=0"`
// 	UpdatedDate string      `json:"updated_date" extensions:"x-order=1"`
// 	Status      string      `json:"status" extensions:"x-order=2"`
// 	Items       interface{} `json:"items" extensions:"x-order=3"`
// }

type DataDetail struct {
	UpdatedBy    UpdatedBy   `json:"updated_by" extensions:"x-order=0"`
	UpdatedDate  string      `json:"updated_date" extensions:"x-order=1"`
	Status       string      `json:"status" extensions:"x-order=2"`
	StatusCode   string      `json:"status_code" extensions:"x-order=3"`
	RejectReason string      `json:"reject_reason"`
	Items        interface{} `json:"items" extensions:"x-order=4"`
}

type DataDetailVolume struct {
	UpdatedBy    UpdatedBy   `json:"updated_by" extensions:"x-order=0"`
	UpdatedDate  string      `json:"updated_date" extensions:"x-order=1"`
	Status       string      `json:"status" extensions:"x-order=2"`
	StatusCode   string      `json:"status_code" extensions:"x-order=3"`
	RejectReason string      `json:"reject_reason"`
	Item         interface{} `json:"item" extensions:"x-order=4"`
}

type UpdatedBy struct {
	Id             uint              `json:"id" extensions:"x-order=0"`
	Email          string            `json:"email" extensions:"x-order=1"`
	FullName       string            `json:"full_name" extensions:"x-order=2"`
	Department     models.Department `json:"department" extensions:"x-order=3"`
	ProfilePicture string            `json:"profile_picture" extensions:"x-order=4"`
}

type ColumnDetail struct {
	ComponentTitle string `json:"component_title" extensions:"x-order=0"`
	ComponentType  string `json:"component_type" extensions:"x-order=1"`
	ColumnDataType string `json:"column_data_type" extensions:"x-order=2"`
}

type ChangeSurfaceDetailResponse struct {
	DataChanged DataDetail `json:"data_changed" extensions:"x-order=0"`
	DataCurrent DataDetail `json:"data_current" extensions:"x-order=1"`
}

type SurfaceDetail struct {
	KmStart              float32                `json:"km_start" extensions:"x-order=0"`
	KmEnd                float32                `json:"km_end" extensions:"x-order=1"`
	WidthSurface         float32                `json:"width_surface" extensions:"x-order=2"`
	ThicknessSurface     float32                `json:"thickness_surface" extensions:"x-order=3"`
	WidthShoulderLeft    float32                `json:"width_shoulder_left" extensions:"x-order=4"`
	SurfaceShoulderLeft  SurfaceDetailIDAndName `json:"surface_shoulder_left" extensions:"x-order=5"`
	WidthShoulderRight   float32                `json:"width_shoulder_right" extensions:"x-order=6"`
	SurfaceShoulderRight SurfaceDetailIDAndName `json:"surface_shoulder_right" extensions:"x-order=7"`
	ThicknessBase        float32                `json:"thickness_base" extensions:"x-order=8"`
	MaterialBase         SurfaceDetailIDAndName `json:"material_base" extensions:"x-order=9"`
	ThicknessSubbase     float32                `json:"thickness_subbase" extensions:"x-order=10"`
	MaterialSubbase      SurfaceDetailIDAndName `json:"material_subbase" extensions:"x-order=11"`
	ThicknessSubgrade    float32                `json:"thickness_subgrade" extensions:"x-order=12"`
	MaterialSubgrade     SurfaceDetailIDAndName `json:"material_subgrade" extensions:"x-order=13"`
	Lanes                []SurfaceLane          `json:"lanes" extensions:"x-order=14"`
	GeomCl               string                 `json:"geom_cl" extensions:"x-order=15"`
	CompareStatus        string                 `json:"compare_status" extensions:"x-order=16"`
}

type SurfaceLane struct {
	Surface   SurfaceDetailIDAndName `json:"surface" extensions:"x-order=0"`
	Direction string                 `json:"direction" extensions:"x-order=1"`
	LaneNo    int                    `json:"lane_no" extensions:"x-order=2"`
}

type SurfaceDetailIDAndName struct {
	ID   int    `json:"id" extensions:"x-order=0"`
	Name string `json:"name" extensions:"x-order=1"`
}

type ChangeConditionDetailResponse struct {
	DataChanged DataDetail `json:"data_changed" extensions:"x-order=0"`
	DataCurrent DataDetail `json:"data_current" extensions:"x-order=1"`
}

type RoadConditionDetail struct {
	KmStart                 float32                    `json:"km_start" extensions:"x-order=0"`
	KmEnd                   float32                    `json:"km_end" extensions:"x-order=1"`
	Value                   float32                    `json:"value" extensions:"x-order=2"`
	RoadRangeConditionItems []RoadRangeConditionDetail `json:"items" extensions:"x-order=3"`
}

type RoadRangeConditionDetail struct {
	KmStart float32 `json:"km_start" extensions:"x-order=0"`
	KmEnd   float32 `json:"km_end" extensions:"x-order=1"`
	Value   float32 `json:"value" extensions:"x-order=2"`
	// Grade         int     `json:"grade" extensions:"x-order=3"`
	GeomCl        string `json:"geom_cl" extensions:"x-order=4"`
	ImgFilepath   string `json:"img_filepath" extensions:"x-order=5"`
	CompareStatus string `json:"compare_status" extensions:"x-order=6"`
}

type ChangeDamageDetailResponse struct {
	DataChanged DataDetail `json:"data_changed" extensions:"x-order=0"`
	DataCurrent DataDetail `json:"data_current" extensions:"x-order=1"`
}

type ChangeVolumeAADTDetail struct {
	RoadGroup   models.RoadGroup `json:"road_group_info" extensions:"x-order=0"`
	GeomCl      interface{}      `json:"geom_cl" extensions:"x-order=1"`
	DataChanged DataDetail       `json:"data_changed" extensions:"x-order=2"`
	DataCurrent DataDetail       `json:"data_current" extensions:"x-order=3"`
}

type ChangeVolumeAADTDetailResponse struct {
	RoadGroup   models.RoadGroup `json:"road_group_info" extensions:"x-order=0"`
	GeomCl      interface{}      `json:"geom_cl" extensions:"x-order=1"`
	DataChanged DataDetailVolume `json:"data_changed" extensions:"x-order=1"`
	DataCurrent DataDetailVolume `json:"data_current" extensions:"x-order=2"`
}

type ChangeVolumeAccidentDetail struct {
	RoadGroup   models.RoadGroup `json:"road_group_info" extensions:"x-order=0"`
	GeomCl      interface{}      `json:"geom_cl" extensions:"x-order=1"`
	DataChanged DataDetail       `json:"data_changed" extensions:"x-order=1"`
	DataCurrent DataDetail       `json:"data_current" extensions:"x-order=2"`
}

type ChangeVolumeAccidentDetailResponse struct {
	RoadGroup   models.RoadGroup `json:"road_group_info" extensions:"x-order=0"`
	GeomCl      interface{}      `json:"geom_cl" extensions:"x-order=1"`
	DataChanged DataDetailVolume `json:"data_changed" extensions:"x-order=1"`
	DataCurrent DataDetailVolume `json:"data_current" extensions:"x-order=2"`
}

type RoadDamageDetail struct {
	KmStart float32 `json:"km_start" extensions:"x-order=0"`
	KmEnd   float32 `json:"km_end" extensions:"x-order=1"`
	// AcCracks             float32             `json:"ac_cracks" extensions:"x-order=2"`
	AcIcrack             float32             `json:"ac_icrack" extensions:"x-order=3"`
	AcUcrack             float32             `json:"ac_ucrack" extensions:"x-order=4"`
	AcRavelling          float32             `json:"ac_ravelling" extensions:"x-order=5"`
	AcPatching           float32             `json:"ac_patching" extensions:"x-order=6"`
	AcPothole            float32             `json:"ac_pothole" extensions:"x-order=7"`
	AcSurfaceDeform      float32             `json:"ac_surface_deform" extensions:"x-order=8"`
	AcBleeding           float32             `json:"ac_bleeding" extensions:"x-order=9"`
	CcTransverseCrack    float32             `json:"cc_transverse_crack" extensions:"x-order=10"`
	CcNonTransverseCrack float32             `json:"cc_non_transverse_crack" extensions:"x-order=11"`
	CcFaulting           float32             `json:"cc_faulting" extensions:"x-order=12"`
	CcSpalling           float32             `json:"cc_spalling" extensions:"x-order=13"`
	CcCornerBreaks       float32             `json:"cc_cornerbreaks" extensions:"x-order=14"`
	CcJointSealDamage    float32             `json:"cc_joint_seal_damage" extensions:"x-order=15"`
	CcPatching           float32             `json:"cc_patching"`
	GeomCl               string              `json:"geom_cl" extensions:"x-order=16"`
	RangeDamageItems     []RangeDamageDetail `json:"items" extensions:"x-order=17"`
}

type RangeDamageDetail struct {
	Km float32 `json:"km" extensions:"x-order=0"`
	// AcCracks             float32 `json:"ac_cracks" extensions:"x-order=1"`
	AcIcrack             float32 `json:"ac_icrack" extensions:"x-order=2"`
	AcUcrack             float32 `json:"ac_ucrack" extensions:"x-order=3"`
	AcRavelling          float32 `json:"ac_ravelling" extensions:"x-order=4"`
	AcPatching           float32 `json:"ac_patching" extensions:"x-order=5"`
	AcPothole            float32 `json:"ac_pothole" extensions:"x-order=6"`
	AcSurfaceDeform      float32 `json:"ac_surface_deform" extensions:"x-order=7"`
	AcBleeding           float32 `json:"ac_bleeding" extensions:"x-order=8"`
	CcTransverseCrack    float32 `json:"cc_transverse_crack" extensions:"x-order=9"`
	CcNonTransverseCrack float32 `json:"cc_non_transverse_crack" extensions:"x-order=10"`
	CcFaulting           float32 `json:"cc_faulting" extensions:"x-order=11"`
	CcSpalling           float32 `json:"cc_spalling" extensions:"x-order=12"`
	CcCornerBreaks       float32 `json:"cc_cornerbreaks" extensions:"x-order=13"`
	CcJointSealDamage    float32 `json:"cc_joint_seal_damage" extensions:"x-order=14"`
	CcPatching           float32 `json:"cc_patching"`
	GeomCl               string  `json:"geom_cl" extensions:"x-order=15"`

	ImgFilepath   string `json:"img_filepath" extensions:"x-order=16"`
	CompareStatus string `json:"compare_status" extensions:"x-order=17"`
}

type ChangeVolumeStatusResponse struct {
	Status          models.RefDataStatus `json:"status" extensions:"x-order=0"`
	ModifiedLabel   string               `json:"modified_label" extensions:"x-order=1"`
	RoadGroupInfo   models.RoadGroup     `json:"road_group_info" extensions:"x-order=2"`
	RoadGroupDetail interface{}          `json:"road_group_detail" extensions:"x-order=3"`
	IdParent        int                  `json:"id_parent" extensions:"x-order=4"`
	Category        string               `json:"category" extensions:"x-order=5"`
	UpdateAt        time.Time            `json:"-"`
}

type ChangeAADTDetailOld struct {
	ID            int       `json:"id" gorm:"column:id"`
	RoadGroupID   int       `json:"road_group_id" gorm:"column:road_group_id"`
	Year          int       `json:"years" gorm:"column:year"`
	CreatedBy     int       `json:"created_by" gorm:"column:created_by"`
	CreatedDate   string    `json:"created_date" gorm:"column:create_date"`
	UpdatedBy     int       `json:"updated_by" gorm:"column:updated_by"`
	UpdatedDate   string    `json:"updated_date" gorm:"column:update_date"`
	TheGeom       string    `json:"the_geom" gorm:"column:the_geom"`
	Revision      int       `json:"revision" gorm:"column:revision"`
	Status        string    `json:"status" gorm:"column:status"`
	IdParent      int       `json:"id_parent" gorm:"column:id_parent"`
	RejectReason  string    `json:"reject_reason" gorm:"column:reject_reason"`
	Veh1          int       `json:"veh1" gorm:"column:veh1"`
	Veh2          int       `json:"veh2" gorm:"column:veh2"`
	Veh3          int       `json:"veh3" gorm:"column:veh3"`
	Veh4          int       `json:"veh4" gorm:"column:veh4"`
	Aadt          int       `json:"aadt" gorm:"column:aadt"`
	Esal          int       `json:"esal" gorm:"column:esal"`
	Yax           int       `json:"yax" gorm:"column:yax"`
	SurveyedDate  time.Time `json:"surveyed_date" gorm:"column:surveyed_date"`
	HashData      string    `json:"hash_data" gorm:"column:hash_data"`
	CompareStatus string    `json:"compare_status"`
}

type ChangeAADTDetailNew struct {
	ID            int       `json:"id" gorm:"column:id"`
	RoadGroupID   int       `json:"road_group_id" gorm:"column:road_group_id"`
	Year          int       `json:"years" gorm:"column:year"`
	CreatedBy     int       `json:"created_by" gorm:"column:created_by"`
	CreatedDate   string    `json:"created_date" gorm:"column:create_date"`
	UpdatedBy     int       `json:"updated_by" gorm:"column:updated_by"`
	UpdatedDate   string    `json:"updated_date" gorm:"column:updated_date"`
	TheGeom       string    `json:"the_geom" gorm:"column:the_geom"`
	Revision      int       `json:"revision" gorm:"column:revision"`
	Status        string    `json:"status" gorm:"column:status"`
	IdParent      int       `json:"id_parent" gorm:"column:id_parent"`
	RejectReason  string    `json:"reject_reason" gorm:"column:reject_reason"`
	Veh1          int       `json:"veh1" gorm:"column:veh1"`
	Veh2          int       `json:"veh2" gorm:"column:veh2"`
	Veh3          int       `json:"veh3" gorm:"column:veh3"`
	Veh4          int       `json:"veh4" gorm:"column:veh4"`
	Aadt          int       `json:"aadt" gorm:"column:aadt"`
	Esal          int       `json:"esal" gorm:"column:esal"`
	Yax           int       `json:"yax" gorm:"column:yax"`
	SurveyedDate  time.Time `json:"surveyed_date" gorm:"column:surveyed_date"`
	HashData      string    `json:"hash_data" gorm:"column:hash_data"`
	CompareStatus string    `json:"compare_status"`
}

type ChangeAccidentDetailOld struct {
	ID            int       `json:"id"`
	RoadGroupID   int       `json:"road_group_id" gorm:"column:road_group_id"`
	Year          int       `json:"years"`
	CreatedBy     int       `json:"created_by" gorm:"column:created_by"`
	CreatedDate   string    `json:"created_date" gorm:"column:created_date"`
	UpdatedBy     int       `json:"updated_by"`
	UpdatedDate   string    `json:"updated_date"`
	TheGeom       string    `json:"the_geom"`
	Revision      int       `json:"revision"`
	Status        string    `json:"status" gorm:"column:status"`
	IdParent      int       `json:"id_parent" gorm:"column:id_parent"`
	RejectReason  string    `json:"reject_reason" gorm:"column:reject_reason"`
	Acc1          int       `json:"acc1" gorm:"column:acc1"`
	Acc2          int       `json:"acc2" gorm:"column:acc2"`
	Acc3          int       `json:"acc3" gorm:"column:acc3"`
	Acc4          int       `json:"acc4" gorm:"column:acc4"`
	Total         int       `json:"total" gorm:"column:total"`
	SurveyedDate  time.Time `json:"surveyed_date" gorm:"column:surveyed_date"`
	HashData      string    `json:"old_hash_data" gorm:"column:hash_data"`
	CompareStatus string    `json:"compare_status"`
}

type ChangeAccidentDetailNew struct {
	ID            int       `json:"id"`
	RoadGroupID   int       `json:"road_group_id" gorm:"column:road_group_id"`
	Year          int       `json:"years"`
	CreatedBy     int       `json:"created_by" gorm:"column:created_by"`
	CreatedDate   string    `json:"created_date" gorm:"column:created_date"`
	UpdatedBy     int       `json:"updated_by"`
	UpdatedDate   string    `json:"updated_date"`
	TheGeom       string    `json:"the_geom"`
	Revision      int       `json:"revision"`
	Status        string    `json:"status" gorm:"column:status"`
	IdParent      int       `json:"id_parent" gorm:"column:id_parent"`
	RejectReason  string    `json:"reject_reason" gorm:"column:reject_reason"`
	Acc1          int       `json:"acc1" gorm:"column:acc1"`
	Acc2          int       `json:"acc2" gorm:"column:acc2"`
	Acc3          int       `json:"acc3" gorm:"column:acc3"`
	Acc4          int       `json:"acc4" gorm:"column:acc4"`
	Total         int       `json:"total" gorm:"column:total"`
	SurveyedDate  time.Time `json:"surveyed_date" gorm:"column:surveyed_date"`
	HashData      string    `json:"hash_data" gorm:"column:hash_data"`
	CompareStatus string    `json:"compare_status"`
}

type VolumeGeom struct {
	TheGeom string `json:"the_geom"`
	Color   string `json:"color"`
}
