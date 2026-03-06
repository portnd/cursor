package models

import "time"

type ChangedStatus struct {
	RoadId        int
	Code          string
	Name          string
	DirectionId   int
	DirectionName string
	RoadTypeId    int
	RoadTypeName  string
	KmStart       float32
	KmEnd         float32
	StatusId      int
	StatusCode    string
	StatusName    string
	IdParent      int
	AssetId       int
	AssetName     string
	UpdatedDate   time.Time
}

type ChangedVolumeAADTStatus struct {
	VolumeAadt
	Code       string `json:"code" gorm:"column:code"`
	Name       string `json:"name" gorm:"column:name"`
	StatusId   int    `json:"status_id" gorm:"column:status_id"`
	StatusName string `json:"status_name" gorm:"column:status_name"`
}

// type VolumeAadt struct {
// 	ID           int       `json:"id"`
// 	RoadGroupID  int       `json:"road_group_id" gorm:"column:road_group_id"`
// 	Year         int       `json:"years"`
// 	CreatedBy    int       `json:"created_by" gorm:"column:created_by"`
// 	CreatedDate  time.Time `gorm:"column:created_date"`
// 	UpdatedBy    int       `json:"updated_by"`
// 	UpdatedDate  time.Time `json:"updated_date"`
// 	TheGeom      string    `json:"the_geom"`
// 	Revision     int       `json:"revision"`
// 	Status       string    `json:"status" gorm:"column:status"`
// 	IdParent     int       `json:"id_parent" gorm:"column:id_parent"`
// 	RejectReason string    `json:"reject_reason" gorm:"column:reject_reason"`
// 	Veh1         int       `json:"veh1" gorm:"column:veh1"`
// 	Veh2         int       `json:"veh2" gorm:"column:veh2"`
// 	Veh3         int       `json:"veh3" gorm:"column:veh3"`
// 	Veh4         int       `json:"veh4" gorm:"column:veh4"`
// 	Aadt         int       `json:"aadt" gorm:"column:aadt"`
// 	Esal         int       `json:"esal" gorm:"column:esal"`
// 	Yax          int       `json:"yax" gorm:"column:yax"`
// 	SurveyedDate time.Time `json:"surveyed_date" gorm:"column:surveyed_date"`
// 	HashData     string    `json:"hash_data" gorm:"column:hash_data"`
// }

type ChangedVolumeAccidentStatus struct {
	VolumeAccident
	Code       string `json:"code" gorm:"column:code"`
	Name       string `json:"name" gorm:"column:name"`
	StatusId   int    `json:"status_id" gorm:"column:status_id"`
	StatusName string `json:"status_name" gorm:"column:status_name"`
}

//	type VolumeAccident struct {
//		ID           int       `json:"id"`
//		RoadGroupID  int       `json:"road_group_id" gorm:"column:road_group_id"`
//		Year         int       `json:"years"`
//		CreatedBy    int       `json:"created_by" gorm:"column:created_by"`
//		CreatedDate  time.Time `json:"created_date" gorm:"column:created_date"`
//		UpdatedBy    int       `json:"updated_by"`
//		UpdatedDate  time.Time `json:"updated_date"`
//		TheGeom      string    `json:"the_geom"`
//		Revision     int       `json:"revision"`
//		Status       string    `json:"status" gorm:"column:status"`
//		IdParent     int       `json:"id_parent" gorm:"column:id_parent"`
//		RejectReason string    `json:"reject_reason" gorm:"column:reject_reason"`
//		Acc1         int       `json:"acc1" gorm:"column:acc1"`
//		Acc2         int       `json:"acc2" gorm:"column:acc2"`
//		Acc3         int       `json:"acc3" gorm:"column:acc3"`
//		Acc4         int       `json:"acc4" gorm:"column:acc4"`
//		Total        int       `json:"total" gorm:"column:total"`
//		SurveyedDate time.Time `json:"surveyed_date" gorm:"column:surveyed_date"`
//		HashData     string    `json:"hash_data" gorm:"column:hash_data"`
//	}
type ColumnsList struct {
	TableName      string
	GeomType       int
	IconFilepath   string
	LineColorCode  string
	ColumnName     string
	TableNameRef   string
	ColumnDataType string
	ComponentTitle string
	ComponentType  string
	Seq            int
}

type RefTableData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ChangeDetailDataOnlyUpdateDateUpdateByAndStatus struct {
	NewUpdatedDate  time.Time
	OldUpdatedDate  time.Time
	NewUpdatedBy    int
	OldUpdatedBy    int
	NewStatus       string
	OldStatus       string
	NewRejectReason string
	OldRejectReason string
}

type ChangeSurfaceDetailData struct {
	NewUpdatedDate              time.Time
	OldUpdatedDate              time.Time
	NewUpdatedBy                int
	OldUpdatedBy                int
	NewStatus                   string
	OldStatus                   string
	NewGeomCl                   string
	OldGeomCl                   string
	NewKmStart                  float32
	OldKmStart                  float32
	NewKmEnd                    float32
	OldKmEnd                    float32
	NewHashData                 string
	OldHashData                 string
	NewWidthSurface             float32
	OldWidthSurface             float32
	NewThicknessSurface         float32
	OldThicknessSurface         float32
	NewWidthShoulderLeft        float32
	OldWidthShoulderLeft        float32
	NewSurfaceShoulderLeftID    int
	OldSurfaceShoulderLeftID    int
	NewSurfaceShoulderLeftName  string
	OldSurfaceShoulderLeftName  string
	NewWidthShoulderRight       float32
	OldWidthShoulderRight       float32
	NewSurfaceShoulderRightID   int
	OldSurfaceShoulderRightID   int
	NewSurfaceShoulderRightName string
	OldSurfaceShoulderRightName string
	NewThicknessBase            float32
	OldThicknessBase            float32
	NewMaterialBaseID           int
	OldMaterialBaseID           int
	NewMaterialBaseName         string
	OldMaterialBaseName         string
	NewThicknessSubbase         float32
	OldThicknessSubbase         float32
	NewMaterialSubbaseID        int
	OldMaterialSubbaseID        int
	NewMaterialSubbaseName      string
	OldMaterialSubbaseName      string
	NewThicknessSubgrade        float32
	OldThicknessSubgrade        float32
	NewMaterialSubgradeID       int
	OldMaterialSubgradeID       int
	NewMaterialSubgradeName     string
	OldMaterialSubgradeName     string
	NewLaneSurfaceID            int
	OldLaneSurfaceID            int
	NewLaneSurfaceName          string
	OldLaneSurfaceName          string
	NewDirectionName            string
	OldDirectionName            string
	NewLaneNo                   int
	OldLaneNo                   int
	NewRejectReason             string
	CompareStatus               string
}

type ChangeAccidentDetailNew struct {
	ID             int       `json:"new_id"`
	RoadGroupID    int       `json:"new_road_group_id" gorm:"column:road_group_id"`
	Year           int       `json:"new_years"`
	CreatedBy      int       `json:"new_created_by" gorm:"column:created_by"`
	CreatedDate    time.Time `json:"new_created_date" gorm:"column:created_date"`
	UpdatedBy      int       `json:"new_updated_by"`
	UpdatedDate    time.Time `json:"new_updated_date"`
	TheGeom        string    `json:"new_the_geom"`
	Revision       int       `json:"new_revision"`
	Status         string    `json:"new_status" gorm:"column:status"`
	IdParent       int       `json:"new_id_parent" gorm:"column:id_parent"`
	RejectReason   string    `json:"new_reject_reason" gorm:"column:reject_reason"`
	Acc1           int       `json:"new_acc1" gorm:"column:acc1"`
	Acc2           int       `json:"new_acc2" gorm:"column:acc2"`
	Acc3           int       `json:"new_acc3" gorm:"column:acc3"`
	Acc4           int       `json:"new_acc4" gorm:"column:acc4"`
	Total          int       `json:"new_total" gorm:"column:total"`
	SurveyedDate   time.Time `json:"new_surveyed_date" gorm:"column:surveyed_date"`
	HashData       string    `json:"new_hash_data" gorm:"column:hash_data"`
	CompareStatus  string    `json:"compare_status"`
	NewUpdatedBy   int
	NewUpdatedDate time.Time
	NewStatus      string
}

type ChangeAccidentDetailOld struct {
	ID             int       `json:"old_id"`
	RoadGroupID    int       `json:"old_road_group_id" gorm:"column:road_group_id"`
	Year           int       `json:"old_years"`
	CreatedBy      int       `json:"old_created_by" gorm:"column:created_by"`
	CreatedDate    time.Time `json:"old_created_date" gorm:"column:created_date"`
	UpdatedBy      int       `json:"old_updated_by"`
	UpdatedDate    time.Time `json:"old_updated_date"`
	TheGeom        string    `json:"old_the_geom"`
	Revision       int       `json:"old_revision"`
	Status         string    `json:"old_status" gorm:"column:status"`
	IdParent       int       `json:"old_id_parent" gorm:"column:id_parent"`
	RejectReason   string    `json:"old_reject_reason" gorm:"column:reject_reason"`
	Acc1           int       `json:"old_acc1" gorm:"column:acc1"`
	Acc2           int       `json:"old_acc2" gorm:"column:acc2"`
	Acc3           int       `json:"old_acc3" gorm:"column:acc3"`
	Acc4           int       `json:"old_acc4" gorm:"column:acc4"`
	Total          int       `json:"old_total" gorm:"column:total"`
	SurveyedDate   time.Time `json:"old_surveyed_date" gorm:"column:surveyed_date"`
	HashData       string    `json:"old_hash_data" gorm:"column:hash_data"`
	OldUpdatedBy   int
	OldUpdatedDate time.Time
	OldStatus      string
}
type ChangeAADTDetailNew struct {
	ID             int       `json:"new_id" gorm:"column:id"`
	RoadGroupID    int       `json:"new_road_group_id" gorm:"column:road_group_id"`
	Year           int       `json:"new_years" gorm:"column:year"`
	CreatedBy      int       `json:"new_created_by" gorm:"column:created_by"`
	CreatedDate    time.Time `json:"new_created_date" gorm:"column:create_date"`
	UpdatedBy      int       `json:"new_updated_by" gorm:"column:updated_by"`
	UpdatedDate    time.Time `json:"new_updated_date" gorm:"column:updated_date"`
	TheGeom        string    `json:"new_the_geom" gorm:"column:the_geom"`
	Revision       int       `json:"new_revision" gorm:"column:revision"`
	Status         string    `json:"new_status" gorm:"column:status"`
	IdParent       int       `json:"new_id_parent" gorm:"column:id_parent"`
	RejectReason   string    `json:"new_reject_reason" gorm:"column:reject_reason"`
	Veh1           int       `json:"new_veh1" gorm:"column:veh1"`
	Veh2           int       `json:"new_veh2" gorm:"column:veh2"`
	Veh3           int       `json:"new_veh3" gorm:"column:veh3"`
	Veh4           int       `json:"new_veh4" gorm:"column:veh4"`
	Aadt           int       `json:"new_aadt" gorm:"column:aadt"`
	Esal           int       `json:"new_esal" gorm:"column:esal"`
	Yax            int       `json:"new_yax" gorm:"column:yax"`
	SurveyedDate   time.Time `json:"new_surveyed_date" gorm:"column:surveyed_date"`
	HashData       string    `json:"new_hash_data" gorm:"column:hash_data"`
	CompareStatus  string    `json:"compare_status"`
	NewUpdatedBy   int
	NewUpdatedDate time.Time
	NewStatus      string
}
type ChangeAADTDetailOld struct {
	ID             int       `json:"old_id" gorm:"column:id"`
	RoadGroupID    int       `json:"old_road_group_id" gorm:"column:road_group_id"`
	Year           int       `json:"old_years" gorm:"column:year"`
	CreatedBy      int       `json:"old_created_by" gorm:"column:created_by"`
	CreatedDate    time.Time `json:"old_created_date" gorm:"column:create_date"`
	UpdatedBy      int       `json:"old_updated_by" gorm:"column:updated_by"`
	UpdatedDate    time.Time `json:"old_updated_date" gorm:"column:update_date"`
	TheGeom        string    `json:"old_the_geom" gorm:"column:the_geom"`
	Revision       int       `json:"old_revision" gorm:"column:revision"`
	Status         string    `json:"old_status" gorm:"column:status"`
	IdParent       int       `json:"old_id_parent" gorm:"column:id_parent"`
	RejectReason   string    `json:"old_reject_reason" gorm:"column:reject_reason"`
	Veh1           int       `json:"old_veh1" gorm:"column:veh1"`
	Veh2           int       `json:"old_veh2" gorm:"column:veh2"`
	Veh3           int       `json:"old_veh3" gorm:"column:veh3"`
	Veh4           int       `json:"old_veh4" gorm:"column:veh4"`
	Aadt           int       `json:"old_aadt" gorm:"column:aadt"`
	Esal           int       `json:"old_esal" gorm:"column:esal"`
	Yax            int       `json:"old_yax" gorm:"column:yax"`
	SurveyedDate   time.Time `json:"old_surveyed_date" gorm:"column:surveyed_date"`
	HashData       string    `json:"old_hash_data" gorm:"column:hash_data"`
	OldUpdatedBy   int
	OldUpdatedDate time.Time
	OldStatus      string
}

type ChangeAADTDetailDelta struct {
	ChangeAADTDetailNew
	ChangeAADTDetailOld
}

type ChangeAccidentDetailDelta struct {
	ChangeAccidentDetailNew
	ChangeAccidentDetailOld
}
type ChangeDamageDetailData struct {
	NewUpdatedDate           time.Time
	OldUpdatedDate           time.Time
	NewUpdatedBy             int
	OldUpdatedBy             int
	NewStatus                string
	OldStatus                string
	NewKmStartM              float32
	OldKmStartM              float32
	NewKmEndM                float32
	OldKmEndM                float32
	NewAcCracksM             float32
	OldAcCracksM             float32
	NewAcIcrackM             float32
	OldAcIcrackM             float32
	NewAcUcrackM             float32
	OldAcUcrackM             float32
	NewAcRavellingM          float32
	OldAcRavellingM          float32
	NewAcPatchingM           float32
	OldAcPatchingM           float32
	NewAcPotholeM            float32
	OldAcPotholeM            float32
	NewAcSurfaceDeformM      float32
	OldAcSurfaceDeformM      float32
	NewAcBleedingM           float32
	OldAcBleedingM           float32
	NewCcTransverseCrackM    float32
	OldCcTransverseCrackM    float32
	NewCcNonTransverseCrackM float32
	OldCcNonTransverseCrackM float32
	NewCcFaultingM           float32
	OldCcFaultingM           float32
	NewCcSpallingM           float32
	OldCcSpallingM           float32
	NewCcCornerbreaksM       float32
	OldCcCornerbreaksM       float32
	NewCcJointSealDamageM    float32
	OldCcJointSealDamageM    float32
	NewCcPatchingM           float32
	OldCcPatchingM           float32
	NewGeomClM               string
	OldGeomClM               string
	NewKmP                   float32
	OldKmP                   float32
	NewAcCracksP             float32
	OldAcCracksP             float32
	NewAcIcrackP             float32
	OldAcIcrackP             float32
	NewAcUcrackP             float32
	OldAcUcrackP             float32
	NewAcRavellingP          float32
	OldAcRavellingP          float32
	NewAcPatchingP           float32
	OldAcPatchingP           float32
	NewAcPotholeP            float32
	OldAcPotholeP            float32
	NewAcSurfaceDeformP      float32
	OldAcSurfaceDeformP      float32
	NewAcBleedingP           float32
	OldAcBleedingP           float32
	NewCcTransverseCrackP    float32
	OldCcTransverseCrackP    float32
	NewCcNonTransverseCrackP float32
	OldCcNonTransverseCrackP float32
	NewCcFaultingP           float32
	OldCcFaultingP           float32
	NewCcSpallingP           float32
	OldCcSpallingP           float32
	NewCcCornerbreaksP       float32
	OldCcCornerbreaksP       float32
	NewCcJointSealDamageP    float32
	OldCcJointSealDamageP    float32
	NewCcPatchingP           float32
	OldCcPatchingP           float32
	NewGeomClP               string
	OldGeomClP               string
	NewImgFilepath           string
	OldImgFilepath           string
	NewRejectReason          string
	CompareStatus            string
}

type ChangeConditionDetailData struct {
	NewUpdatedDate time.Time
	OldUpdatedDate time.Time
	NewUpdatedBy   int
	OldUpdatedBy   int
	NewStatus      string
	OldStatus      string
	NewKmStartKm   float32
	OldKmStartKm   float32
	NewKmEndKm     float32
	OldKmEndKm     float32
	NewValueKm     float32
	OldValueKm     float32
	NewKmStartM    float32
	OldKmStartM    float32
	NewKmEndM      float32
	OldKmEndM      float32
	NewValueM      float32
	OldValueM      float32
	// NewGrade       int
	// OldGrade       int
	NewGeomCl       string
	OldGeomCl       string
	CompareStatus   string
	NewImgFilepath  string
	OldImgFilepath  string
	NewRejectReason string
}
