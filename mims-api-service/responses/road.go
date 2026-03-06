package responses

import (
	"time"

	"gitlab.com/mims-api-service/models"
)

type RoadMenuData struct {
	GroupId      int    `json:"group_id"`
	GroupName    string `json:"group_name"`
	AssetId      int    `json:"asset_id"`
	AssetName    string `json:"asset_name"`
	GeomType     int    `json:"geom_type"`
	IsInRoad     bool   `json:"is_in_road"`
	IconFilepath string `json:"icon_filepath"`
	IsActive     bool   `json:"is_active"`
}

type RoadGroup struct {
	Id            int        `json:"id"`
	Name          string     `json:"name"`
	Code          string     `json:"code"`
	CountWaiting  int        `json:"count_waiting"`
	CountRejected int        `json:"count_rejected"`
	Road          []RoadMain `json:"roads"`
}

type RoadMain struct {
	RoadID       int     `json:"road_id"`
	Code         string  `json:"code"`
	Seq          int     `json:"seq"`
	Name         string  `json:"name"`
	KmStart      float32 `json:"km_start"`
	KmEnd        float32 `json:"km_end"`
	RoadLevel    int     `json:"road_level"`
	ParentRoadID int     `json:"parent_road_id"`
	StatusLatest string  `json:"status_latest"`
	LaneCount    int     `json:"lane_count"`
	Direction    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"direction"`
	RoadType struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"road_type"`
	GeomCl         string      `json:"geom_cl"`
	CountWaiting   int         `json:"count_waiting"`
	CountRejected  int         `json:"count_rejected"`
	RoadColorCode  string      `json:"road_color_code"`
	RoadTypeIconID int         `json:"road_type_icon_id"`
	ChildRoad      []ChildRoad `json:"child_road"`
}

type RoadMainById struct {
	RoadID       int     `json:"road_id"`
	Code         string  `json:"code"`
	Seq          int     `json:"seq"`
	Name         string  `json:"name"`
	KmStart      float32 `json:"km_start"`
	KmEnd        float32 `json:"km_end"`
	RoadLevel    int     `json:"road_level"`
	ParentRoadID int     `json:"parent_road_id"`
	// StatusLatest string  `json:"status_latest"`
	LaneCount int `json:"lane_count"`
	Direction struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"direction"`
	RoadType struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"road_type"`
	GeomCl string `json:"geom_cl"`
	// CountWaiting   int         `json:"count_waiting"`
	// CountRejected  int         `json:"count_rejected"`
	RoadColorCode  string `json:"road_color_code"`
	RoadTypeIconID int    `json:"road_type_icon_id"`
}

type ChildRoad struct {
	RoadID       int     `json:"road_id"`
	Code         string  `json:"code"`
	Seq          int     `json:"seq"`
	Name         string  `json:"name"`
	KmStart      float32 `json:"km_start"`
	KmEnd        float32 `json:"km_end"`
	RoadLevel    int     `json:"road_level"`
	ParentRoadID int     `json:"parent_road_id"`
	StatusLatest string  `json:"status_latest"`
	LaneCount    int     `json:"lane_count"`
	Direction    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"direction"`
	RoadType struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"road_type"`
	GeomCl         string `json:"geom_cl"`
	CountWaiting   int    `json:"count_waiting"`
	CountRejected  int    `json:"count_rejected"`
	RoadColorCode  string `json:"road_color_code"`
	RoadTypeIconID int    `json:"road_type_icon_id"`
}

type AccessGroupMenu struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId int    `json:"parent_id"`
	Route    string `json:"route"`
	Icon     string `json:"icon"`
}

type StatusCount struct {
	Category string `json:"category"`
	ID       int    `json:"id"`
	GroupId  int    `json:"group_id"`
	// Status        string `json:"status"`
	CountTemp     int `json:"count_temp"`
	CountWaiting  int `json:"count_waiting"`
	CountRejected int `json:"Count_rejected"`
}

type RoadConditionList struct {
	Year  int                     `json:"year" extensions:"x-order=0"`
	Items []RoadConditionListItem `json:"items" extensions:"x-order=1"`
}
type RoadConditionListItem struct {
	ID        int `json:"id" extensions:"x-order=0"`
	IDParent  int `json:"id_parent" extensions:"x-order=1"`
	Direction struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"direction" extensions:"x-order=2"`
	LaneNo       int       `json:"lane_no" extensions:"x-order=3"`
	SurveyedDate time.Time `json:"surveyed_date" extensions:"x-order=4"`
	Revision     int       `json:"-"`
}

type RoadRetroReflectivityList struct {
	Year  int                             `json:"year" extensions:"x-order=0"`
	Items []RoadRetroReflectivityListItem `json:"items" extensions:"x-order=1"`
}

type RoadRetroReflectivityListItem struct {
	ID        int `json:"id" extensions:"x-order=0"`
	IDParent  int `json:"id_parent" extensions:"x-order=1"`
	Direction struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"direction" extensions:"x-order=2"`
	LineNo       int       `json:"line_no" extensions:"x-order=3"`
	SurveyedDate time.Time `json:"surveyed_date" extensions:"x-order=4"`
	Revision     int       `json:"-"`
}

type Permissions struct {
	CanEdit    bool `json:"can_edit"`
	CanDelete  bool `json:"can_delete"`
	CanApprove bool `json:"can_approve"`
	CanSend    bool `json:"can_send"`
	CanReject  bool `json:"can_reject"`
}

type Direction struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type RoadConditionDetails struct {
	ID             int                         `json:"id"`
	IDParent       int                         `json:"id_parent"`
	UpdatedDate    string                      `json:"updated_date"`
	UpdatedBy      RoadConditionDetailUserInfo `json:"updated_by"`
	Status         string                      `json:"status"`
	Direction      models.RefDirection         `json:"direction"`
	ConditionTypes []RoadConditionDetailData   `json:"condition_datas"`
}

type RetroReflectivityDetails struct {
	ID            int                               `json:"id"`
	IDParent      int                               `json:"id_parent"`
	UpdatedDate   string                            `json:"updated_date"`
	UpdatedBy     RoadUserInfo                      `json:"updated_by"`
	Status        string                            `json:"status"`
	Direction     models.RefDirection               `json:"direction"`
	HasWhiteLine  bool                              `json:"has_white_line"`
	HasYellowLine bool                              `json:"has_yellow_line"`
	Datas         []RoadRetroReflectivityDetailData `json:"datas"`
}

type RoadConditionDetailData struct {
	ConditionType                 string                          `json:"condition_type"`
	RoadConditionDetailHeaderItem []RoadConditionDetailHeaderItem `json:"items"`
}

type RoadRetroReflectivityDetailData struct {
	Color                                  string                                  `json:"color"`
	RoadRetroReflectivityDetailHeaderItems []RoadRetroReflectivityDetailHeaderItem `json:"items"`
}

type RoadConditionDetailUserInfo struct {
	UID        string `json:"uid"`
	UserName   string `json:"user_name"`
	FullName   string `json:"full_name"`
	Department struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"department"`
	ProfilePicture string `json:"profile_picture"`
}

type RoadUserInfo struct {
	UID        string `json:"uid"`
	UserName   string `json:"user_name"`
	FullName   string `json:"full_name"`
	Department struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"department"`
	ProfilePicture string `json:"profile_picture"`
}
type RoadConditionDetailHeaderItem struct {
	KmStart    int                           `json:"km_start"`
	KmEnd      int                           `json:"km_end"`
	Value      *float64                      `json:"value"`
	GeomCl     string                        `json:"geom_cl"`
	SurveyType string                        `json:"survey_type"`
	Items      []RoadConditionDetailBodyItem `json:"items"`
}

type RoadConditionDetailBodyItem struct {
	KmStart int      `json:"km_start"`
	KmEnd   int      `json:"km_end"`
	Value   *float64 `json:"value"`
	//	Grade       int     `json:"grade"`
	GeomCl      string `json:"geom_cl"`
	ImgFilepath string `json:"img_filepath"`
	SurveyType  string `json:"survey_type"`
}

type RoadRetroReflectivityDetailHeaderItem struct {
	KmStart          int                                   `json:"km_start"`
	KmEnd            int                                   `json:"km_end"`
	RetroAvg         *float64                              `json:"retro_avg"`
	GeomCl           string                                `json:"geom_cl"`
	RefStripeColorID int                                   `json:"ref_stripe_color_id"`
	RefStripeColor   models.RefStripeColor                 `json:"ref_stripe_color"`
	RefStripeTypeID  int                                   `json:"ref_stripe_type_id"`
	RefStripeType    models.RefStripeType                  `json:"ref_stripe_type"`
	Items            []RoadRetroReflectivityDetailBodyItem `json:"items"`
}
type RoadRetroReflectivityDetailBodyItem struct {
	KmStart          int                   `json:"km_start"`
	KmEnd            int                   `json:"km_end"`
	RetroAvg         *float64              `json:"retro_avg"`
	GeomCl           string                `json:"geom_cl"`
	RefStripeColorID int                   `json:"ref_stripe_color_id"`
	RefStripeColor   models.RefStripeColor `json:"ref_stripe_color"`
	RefStripeTypeID  int                   `json:"ref_stripe_type_id"`
	RefStripeType    models.RefStripeType  `json:"ref_stripe_type"`
}

type RoadConditionCreate struct {
	Id int `json:"id"`
}

type RoadRetroReflectivityCreate struct {
	ID int `json:"id"`
}

type RoadConditionUpdate struct {
	ID       int `json:"id"`
	IDParent int `json:"id_parent"`
}

type RoadCondition struct {
	Id           int                 `json:"id" extensions:"x-order=0"`
	IDParent     int                 `json:"id_parent" extensions:"x-order=1"`
	LaneNo       int                 `json:"lane_no" extensions:"x-order=2"`
	SurveyedDate time.Time           `json:"surveyed_date" extensions:"x-order=3"`
	Remarks      string              `json:"remarks" extensions:"x-order=4"`
	IriFilename  string              `json:"iri_filename" extensions:"x-order=5"`
	ImgFilepath  string              `json:"img_filepath" extensions:"x-order=6"`
	Direction    models.RefDirection `json:"direction" extensions:"x-order=7"`
}

type RoadRetroReflectivity struct {
	ID           int                 `json:"id" extensions:"x-order=0"`
	IDParent     int                 `json:"id_parent" extensions:"x-order=1"`
	LineNo       int                 `json:"line_no" extensions:"x-order=2"`
	SurveyedDate time.Time           `json:"surveyed_date" extensions:"x-order=3"`
	Remarks      string              `json:"remarks" extensions:"x-order=4"`
	CsvFile      string              `json:"csv_file" extensions:"x-order=5"`
	Direction    models.RefDirection `json:"direction" extensions:"x-order=7"`
}

type RoadDamageSetData struct {
	ID       int `json:"id"`
	IDParent int `json:"id_parent"`
}

type RoadDamageList struct {
	Year int                             `json:"year" extensions:"x-order=0"`
	Item []models.RoadDamageListResponse `json:"items" extensions:"x-order=1"`
}

// type RoadDamageListResponses struct {
// 	Year int                      `json:"year"`
// 	Item []RoadDamageListResponse `json:"item"`
// }

type RoadDamageImport struct {
	ID                int       `json:"id" extensions:"x-order=0"`
	IDParent          int       `json:"id_parent" extensions:"x-order=1"`
	LaneNo            int       `json:"lane_no" extensions:"x-order=2"`
	SurveyedDate      time.Time `json:"surveyed_date" extensions:"x-order=3"`
	DamageFilenameint string    `json:"damage_filename" extensions:"x-order=4"`
	ImgFilepath       string    `json:"img_filepath" extensions:"x-order=5"`
	Direction         struct {
		ID   int    `json:"id" extensions:"x-order=0"`
		Name string `json:"name" extensions:"x-order=1"`
	} `json:"direction" extensions:"x-order=6"`
}

type RoadDamageListDetail struct {
	Id          int                   `json:"id" extensions:"x-order=0"`
	IdParent    int                   `json:"id_parent" extensions:"x-order=1"`
	UpdatedDate string                `json:"updated_date" extensions:"x-order=2"`
	UpdatedBy   models.UserDepartment `json:"updated_by" extensions:"x-order=3"`
	Status      string                `json:"status" extensions:"x-order=4"`
	Permission  struct {
		CanEdit    bool `json:"can_edit"`
		CanDelete  bool `json:"can_delete"`
		CanApprove bool `json:"can_approve"`
		CanSend    bool `json:"can_send"`
		CanReject  bool `json:"can_reject"`
	} `json:"permissions" extensions:"x-order=5"`
	TheGeom    string                  `json:"the_geom" extensions:"x-order=6"`
	RoadDamage RoadDamageDetailRespond `json:"road_damage" extensions:"x-order=7"`
}

type RoadDamage struct {
	RoadId               int       `json:"road_id"`
	LaneNo               int       `json:"lane_no"`
	Year                 int       `json:"year"`
	KmStart              float64   `json:"km_start"`
	KmEnd                float64   `json:"km_end"`
	SurveyedDate         time.Time `json:"surveyed_date"`
	Revision             int       `json:"revision"`
	Status               string    `json:"status"`
	IdParent             int       `json:"id_parent"`
	ImgFilepath          string    `json:"img_filepath"`
	DamageInputFilepath  string    `json:"damage_input_filepath"`
	RejectReason         string    `json:"reject_reason"`
	AcIcrack             float64   `json:"ac_icrack"`
	AcUcrack             float64   `json:"ac_ucrack"`
	AcRavelling          float64   `json:"ac_ravelling"`
	AcPatching           float64   `json:"ac_patching"`
	AcPothole            float64   `json:"ac_pothole"`
	AcSurfaceDeform      float64   `json:"ac_surface_deform"`
	AcBleeding           float64   `json:"ac_bleeding"`
	CcTransverseCrack    float64   `json:"cc_transverse_crack"`
	CcNonTransverseCrack float64   `json:"cc_non_transverse_crack"`
	CcFaulting           float64   `json:"cc_faulting"`
	CcSpalling           float64   `json:"cc_spalling"`
	CcCornerbreaks       float64   `json:"cc_cornerbreaks"`
	CcJointSealDamage    float64   `json:"cc_joint_seal_damage"`
	CcPatching           float64   `json:"cc_patching"`
}

type RoadConditionLane struct {
	Year  int                     `json:"year" extensions:"x-order=0"`
	Items []RoadConditionLaneItem `json:"items"`
}

type RoadRetroReflectivityLine struct {
	Year  int                             `json:"year" extensions:"x-order=0"`
	Items []RoadRetroReflectivityLineItem `json:"items"`
}

type RoadConditionLaneItem struct {
	LaneNo  int      `json:"lane_no" extensions:"x-order=1"`
	KmStart int      `json:"km_start" extensions:"x-order=2"`
	KmEnd   int      `json:"km_end" extensions:"x-order=3"`
	Value   *float64 `json:"value" extensions:"x-order=4"`
}

type RoadRetroReflectivityLineItem struct {
	LineNo  int      `json:"line_no" extensions:"x-order=1"`
	KmStart int      `json:"km_start" extensions:"x-order=2"`
	KmEnd   int      `json:"km_end" extensions:"x-order=3"`
	Value   *float64 `json:"value" extensions:"x-order=4"`
}

type RoadConditionYear struct {
	Lane  int                     `json:"lane" extensions:"x-order=0"`
	Items []RoadConditionYearItem `json:"items"`
}

type RoadRetroReflectivityYearItem struct {
	Line    int     `json:"-"`
	Year    int     `json:"year"`
	KmStart int     `json:"km_start"`
	KmEnd   int     `json:"km_end"`
	Value   float64 `json:"value"`
}

type RoadRetroReflectivityYear struct {
	Line  int                             `json:"line" extensions:"x-order=0"`
	Items []RoadRetroReflectivityYearItem `json:"items"`
}

type RoadConditionYearItem struct {
	Lane    int      `json:"-"`
	Year    int      `json:"year"`
	KmStart int      `json:"km_start"`
	KmEnd   int      `json:"km_end"`
	Value   *float64 `json:"value"`
}

type RoadConditionAverage struct {
	Lane  int                               `json:"lane" extensions:"x-order=0"`
	Items []models.RoadConditionAverageItem `json:"items"`
}

type RoadRetroReflectivityAverage struct {
	Lane  int                                  `json:"lane" extensions:"x-order=0"`
	Items []models.RoadReflectivityAverageItem `json:"items"`
}

type RoadLaneList struct {
	LaneNo   int    `json:"lane_no"`
	LaneName string `json:"lane_name"`
}

type RoadRetroReflectivityLineList struct {
	LineNo int `json:"line_no"`
}

type RoadAssetData struct {
	// Page            int                   `json:"page"`
	ID           int          `json:"id"`
	IDParent     int          `json:"id_parent"`
	UpdatedDate  string       `json:"updated_date"`
	Revision     int          `json:"revision"`
	Status       string       `json:"status"`
	CanEdit      bool         `json:"can_edit"`
	StatusCode   string       `json:"status_code"`
	RejectReason string       `json:"reject_reason"`
	UpdatedBy    RoadUserInfo `json:"updated_by"`
	Permissions  struct {
		CanEdit    bool `json:"can_edit"`
		CanDelete  bool `json:"can_delete"`
		CanApprove bool `json:"can_approve"`
		CanSend    bool `json:"can_send"`
		CanReject  bool `json:"can_reject"`
	} `json:"permissions"`
	IsExclusiveLock       bool        `json:"is_exclusive_lock"`
	RoadAssets            interface{} `json:"road_assets"`
	DataColumns           interface{} `json:"data_columns"`
	IconFilepath          string      `json:"icon_filepath"`
	ThumbnailIconFilepath string      `json:"thumbnail_icon_filepath"`
	LineColorCode         string      `json:"line_color_code"`
}

type RoadAssetDetailColumn struct {
	// IDParent       int    `json:"id_parent"`
	RefAssetID     int    `json:"-"`
	IconFilepath   string `json:"icon_filepath"`
	LineColorCode  string `json:"line_color_code"`
	TableName      string `json:"table_name"`
	ColumnName     string `json:"column_name"`
	TableNameRef   string `json:"table_name_ref"`
	ColumnDataType string `json:"column_data_type"`
	ComponentTitle string `json:"component_title"`
	ComponentType  string `json:"component_type"`
	Seq            int    `json:"seq"`
}

type RoadAssetRevision struct {
	ID              int    `json:"id"`
	IDParent        int    `json:"id_parent"`
	UpdatedDate     string `json:"updated_date"`
	RevisionNo      int    `json:"revision_no"`
	IsExclusiveLock bool   `json:"is_exclusive_lock"`
	Status          string `json:"status"`
	StatusText      string `json:"-"`
}

type RoadAssetTemplateColumn struct {
	Seq            int    `json:"-"`
	TableName      string `json:"table_name"`
	TableNameRef   string `json:"table_name_ref"`
	ColumnName     string `json:"column_name"`
	ComponentTitle string `json:"component_title"`
	ComponentType  string `json:"component_type"`
	ColumnDataType string `json:"column_data_type"`
	IsRequired     bool   `json:"is_required"`
	Value          string `json:"value"`
}

type RoadKmByGeom struct {
	RoadID  int     `json:"road_id"`
	LaneNo  int     `json:"lane_no"`
	TheGeom string  `json:"the_geom"`
	KmStart float64 `json:"km_start"`
	KmEnd   float64 `json:"km_end"`
	Km      float64 `json:"km"`
	Type    string  `json:"type"`
}

type AssetTableType struct {
	GeomTypeID       int    `json:"geom_type_id"`
	GeomTypeName     string `json:"geom_type_name"`
	Color            string `json:"color"`
	DepartmentManage []int  `json:"department_manage"`
}

type AssetID struct {
	ID int `json:"id" `
}

type RoadTree struct {
	ID       int            `json:"id"`
	Label    string         `json:"label"`
	Children []RoadChildren `json:"children"`
}

type RoadChildren struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}
type RoadSection struct {
	ID                int                `json:"id"`
	RoadGroupId       int                `json:"road_group_id"`
	Number            string             `json:"number"`
	NameOriginTH      string             `json:"name_origin_th"`
	NameDestinationTH string             `json:"name_destination_th"`
	NameOriginEn      string             `json:"name_origin_en"`
	NameDestinationEn string             `json:"name_destination_en"`
	KmStart           float32            `json:"km_start"`
	KmEnd             float32            `json:"km_end"`
	Distance          float32            `json:"distance"`
	ProvinceCode      []string           `json:"-" gorm:"type:character[]"`
	Province          []string           `json:"province" gorm:"type:character[]"`
	RefDivisionCode   string             `json:"-"`
	RefDivision       models.RefDivision `json:"ref_division"`
	RefDistrictCode   string             `json:"-"`
	RefDistrict       models.RefDistrict `json:"ref_district"`
	RefDepotCode      string             `json:"-"`
	RefDepot          models.RefDepot    `json:"ref_depot"`
}

type RoadSectionData struct {
	RoadSection
	Roads []models.RoadData `json:"roads"`
}
type RoadList struct {
	models.RoadGroup
	Sections []RoadSectionData `json:"sections"`
}

type RoadConditionAnalysisData struct {
	RoadConditionSurvey     []models.RoadConditionSurvey
	RoadConditionSurvey100M []models.RoadConditionSurvey100M
	RoadConditionSurveyM    []models.RoadConditionSurveyM
	RoadData                RoadData
}

type RoadRetroReflectivityAnalysisData struct {
	RoadRetroReflectivity      []models.RoadRetroReflectivity
	RoadRetroReflectivityRange []models.RoadRetroReflectivityRange
	RoadRetroReflectivityM     []models.RoadRetroReflectivityM
	RoadData                   models.RoadRetroReflectivityData
}

type RoadData struct {
	IfiAverage          *float64
	IriAverage          *float64
	MpdAverage          *float64
	RutAverage          *float64
	IriKm               *float64
	IfiKm               *float64
	MpdKm               *float64
	RutKm               *float64
	Iri100m             *float64
	Ifi100m             *float64
	Mpd100m             *float64
	Rut100m             *float64
	DividerCount100mIri *float64
	DividerCount100mIfi *float64
	DividerCount100mMpd *float64
	DividerCount100mRut *float64
	DividerCountIri     *float64
	DividerCountIfi     *float64
	DividerCountMpd     *float64
	DividerCountRut     *float64
	TotalM              float64
}

type RoadById struct {
	models.Road
	RoadCode            string                   `json:"road_code"`
	RoadSectionNameTH   string                   `json:"road_section_name_th"`
	RoadSectionNameEN   string                   `json:"road_section_name_en"`
	Province            string                   `json:"province"`
	ResponsibleCode     string                   `json:"responsible_code"`
	OriginToDestination string                   `json:"origin_to_destination"`
	KmRange             string                   `json:"km_range"`
	Distance            float64                  `json:"distance"`
	RoadInfo            models.RoadInfoAddData   `json:"road_info"`
	RoadSurfaceIcon     []models.RoadSurfaceIcon `json:"road_surface_icon"`
	RoadGeom            []models.RoadGeom        `json:"road_geom"`
	RefDepot            models.RefDepot          `json:"ref_depot"`
}

type RoadInit struct {
	RoadCode          string `json:"road_code"`
	RoadSectionNameTH string `json:"road_section_name_th"`
	RoadSectionNameEN string `json:"road_section_name_en"`
	Province          string `json:"province"`
	District          string `json:"district"`
	Depot             string `json:"depot"`
	Origin            string `json:"origin"`
	Destination       string `json:"destination"`
}

type RoadInitData struct {
	Id             int     `json:"id"`
	Name           string  `json:"name"`
	KMStart        float64 `json:"km_start"`
	KMEnd          float64 `json:"km_end"`
	RefDirectionId int     `json:"ref_direction_id"`
	LaneTotal      int     `json:"lane_total"`
	RoadSectionId  int     `json:"-"`
}

type RoadGropInitData struct {
	Id         int    `json:"id"`
	RoadNumber string `json:"road_number"`
	ShortName  string `json:"short_name"`
}

type RoadSectionInitData struct {
	Id                int            `json:"id"`
	RoadGroupId       int            `json:"road_group_id"`
	Number            string         `json:"number"`
	NameOriginTH      string         `json:"name_origin"`
	NameDestinationTH string         `json:"name_destination"`
	Roads             []RoadInitData `json:"roads"  gorm:"ForeignKey:RoadSectionId;AssociationForeignKey:Id"`
}
type RoadListInitData struct {
	RoadGropInitData
	RoadSection []RoadSectionInitData `json:"road_sections" gorm:"ForeignKey:RoadGroupId;AssociationForeignKey:Id"`
}

type RoadLanes struct {
	Id               int    `json:"-"`
	RoadId           int    `json:"road_id"`
	RefDirectionId   int    `json:"ref_direction_id"`
	RefDirectionName string `json:"ref_direction_name"`
	LaneNo           int    `json:"lane_no"`
}

type GeomJSON struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

type DataRes struct {
	Iri1000 map[int][]int
	Iri100  map[int][]int
	Rut100  map[int][]int
	Ifi100  map[int][]int
	G7100   map[int][]int
}

type RoadListNew struct {
	RoadGroup
	RefDepot []RefDepotNew `json:"ref_depot"`
}

type RefDepotNew struct {
	RefDepot
	Section []RoadSectionData `json:"section"`
}
