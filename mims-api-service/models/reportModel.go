package models

import (
	"time"

	"github.com/lib/pq"
)

type Year struct {
	Year int
}

type ReportYear struct {
	Year []int `json:"year"`
}

// road group data
type RoadGroupData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// road group infomation
type RoadGroupReportInfo struct {
	Name string
	Code string
}

// road infomation
type RoadReportInfo struct {
	RoadGroupName string
	RoadName      string
	RoadCode      string
	RoadColorCode string
	KmStart       float64
	KmEnd         float64
	RoadLengthStr string

	RoadID                       int
	RoadGroupNumber              string
	RoadSectionNumber            string
	RoadSectionDistance          float64
	RoadSectionNameOriginTh      string
	RoadSectionNameDestinationTh string
	RoadMainName                 string
}

// multi-road infomation
type MultiRoadInfo struct {
	RoadID   int    `gorm:"column:id"`
	RoadName string `gorm:"column:name"`
	KmStart  int
	KmEnd    int
}

// road damage report

// type DataReportDamage struct {
// 	RoadGroupName string
// 	RoadName      string
// 	RoadCode      string
// 	KmStart       int
// 	KmEnd         int
// 	StrKmStart    string
// 	StrKmEnd      string
// 	RoadLengthStr string
// 	IsNull        bool

// 	Detail []DataRoadDamage
// }
// type DataRoadDamage struct {
// 	LaneNo int

// 	ACIcrack        float64
// 	ACUcrack        float64
// 	ACRavelling     float64
// 	ACPatching      float64
// 	ACPothole       float64
// 	ACSurfaceDeform float64
// 	ACBleeding      float64

// 	CCTransverseCrack    float64
// 	CCNonTransverseCrack float64
// 	CCFaulting           float64
// 	CCSpalling           float64
// 	CCCornerbreaks       float64
// 	CCJointSealDamage    float64
// 	CCPatching           float64

// 	Position []CovertPositionDamage
// }
// type CovertPositionDamage struct {
// 	No            int
// 	Km            int
// 	StrKm         string
// 	Surface       string
// 	DamageType    string
// 	DamageTypeENG string
// 	Value         float64
// 	Unit          string
// 	Image         string
// }
// type PositionDamage struct {
// 	LaneNO      int
// 	Km          int
// 	StrKm       string
// 	ImgFilepath string

// 	ACIcrack        float64
// 	ACUcrack        float64
// 	ACRavelling     float64
// 	ACPatching      float64
// 	ACPothole       float64
// 	ACSurfaceDeform float64
// 	ACBleeding      float64

// 	CCTransverseCrack    float64
// 	CCNonTransverseCrack float64
// 	CCFaulting           float64
// 	CCSpalling           float64
// 	CCCornerbreaks       float64
// 	CCJointSealDamage    float64
// 	CCPatching           float64
// }

// road condition report

type DataReportCondition struct {
	RoadGroupName string
	RoadName      string
	RoadCode      string
	KmStart       int
	KmEnd         int
	StrKmStart    string
	StrKmEnd      string
	RoadLengthStr string
	IsNull        bool

	Lane []DataRoadCondition
}

type DataReportConditionDetail struct {
	RoadGroupName string
	RoadName      string
	RoadCode      string
	KmStart       int
	KmEnd         int
	StrKmStart    string
	StrKmEnd      string
	RoadLengthStr string
	IsNull        bool
}
type DataRoadCondition struct {
	LaneNo  int
	KmStart int
	KmEnd   int
	Iri     float64
	Rut     float64
	Mpd     float64
	Ifi     float64

	RoadLengthStr string

	Detail []DataRoadConditionM
}

type DataRoadConditionDetail struct {
	LaneNo  int
	KmStart int
	KmEnd   int
	Iri     float64
	Rut     float64
	Mpd     float64
	Ifi     float64

	RoadLengthStr string

	// Detail []DataRoadConditionM
}

type DataRoadConditionM struct {
	LaneNo int

	KmStart    int
	KmEnd      int
	StrKmStart string
	StrKmEnd   string
	Iri        float64
	Rut        float64
	Mpd        float64
	Ifi        float64
}

// road maintenance history report

type DataMaintenance struct {
	ID         int
	Name       string
	BudgetYear int
}

type DataResponseMaintenance struct {
	RoadGroupID int    `json:"road_group_id"`
	Name        string `json:"road_group_name"`
	BudgetYear  []int  `json:"year"`
}

type DataGetMaintenance struct {
	RoadID            int     `json:"road_id"`
	RoadGroupNumber   string  `json:"road_group_number"`
	RoadSectionNumber string  `json:"road_section_number"`
	RoadMainName      string  `json:"road_main_name"`
	SecKmStrat        int     `json:"sec_km_start"`
	SecKmEnd          int     `json:"sec_km_end"`
	StrKmStart        string  `json:"str_km_start"`
	StrKmEnd          string  `json:"str_km_end"`
	Distance          float32 `json:"distance"`
	DistanceStr       string  `json:"distance_str"`
	RoadName          string  `gorm:"column:name" json:"road_name"`
	KmStart           int
	KmEnd             int
}

type DataReportMaintenance struct {
	RoadID            int     `json:"road_id"`
	RoadGroupNumber   string  `json:"road_group_number"`
	RoadSectionNumber string  `json:"road_section_number"`
	RoadMainName      string  `json:"road_main_name"`
	SecKmStrat        int     `json:"sec_km_start"`
	SecKmEnd          int     `json:"sec_km_end"`
	StrKmStart        string  `json:"str_km_start"`
	StrKmEnd          string  `json:"str_km_end"`
	Distance          float32 `json:"distance"`
	DistanceStr       string  `json:"distance_str"`
	RoadName          string  `gorm:"column:name" json:"road_name"`
	KmStart           int
	KmEnd             int
	Rows              []interface{}        `json:"rows"`
	RowExcel          []DataRowMaintenance `json:"row_excel"`
	Years             []int                `json:"years"`
	YearStart         int                  `json:"year_start"`
	YearEnd           int                  `json:"year_end"`
	IsNull            bool
}

type DataRowMaintenance struct {
	KmStart    int                   `json:"km_start"`
	KmEnd      int                   `json:"km_end"`
	StrKmStart string                `json:"str_km_start"`
	StrKmEnd   string                `json:"str_km_end"`
	Data       []DataYearMaintenance `json:"data"`
}

type DataYearMaintenance struct {
	Year int                   `json:"year"`
	Data []DataLaneMaintenance `json:"data"`
}

type DataLaneMaintenance struct {
	LaneNo int                   `json:"lane_no"`
	Data   DataDetailMaintenance `json:"data"`
}

type DataDetailMaintenance struct {
	Km      int               `json:"km"`
	KmStart string            `json:"km_start"`
	KmEnd   string            `json:"km_end"`
	LaneNo  int               `json:"lane_no"`
	Span    int               `json:"span"`
	IsWrite int               `json:"is_write"`
	Year    int               `json:"year"`
	Data    MethodMaintenance `json:"data"`
}

type MethodMaintenance struct {
	RoadID     int
	Year       int
	Lane       int
	KmStart    int    `json:"km_start"`
	KmEnd      int    `json:"km_end"`
	StrKmStart string `json:"str_km_start"`
	StrKmEnd   string `json:"str_km_end"`
	Range      int
	Method     string `json:"method"`
	Color      string `json:"color"`
	Unique     string
}
type RangeMaintenance struct {
	Start                    int
	End                      int
	Value                    int
	Writed                   int
	Year                     int
	Lane                     int
	Range                    int
	InterventionCriteriaName string
	Method                   string
}

// road maintenance tracking report
type DataMaintenanceTracking struct {
	Year          int
	RoadGroupID   int
	RoadGroupName string
}

type DataResponseMaintenanceTracking struct {
	Year      int                 `json:"year"`
	RoadGroup []RoadGroupTracking `json:"road_group"`
}

type RoadGroupTracking struct {
	RoadGroupID   int    `json:"road_group_id"`
	RoadGroupName string `json:"road_group_name"`
}

type DataResponseReportMaintenanceTracking struct {
	Year    int `json:"year"`
	IsNull  bool
	Project []DataReportMaintenanceTracking `json:"project"`
}

type DataReportMaintenanceTracking struct {
	ProjectName          string
	RoadGroupName        string
	ContractNumber       string
	BudgetYear           int
	BudgetType           string
	MaintenanceType      string
	ContractorName       string
	AdviserName          string
	ProjectSecretaryName string
	BudgetMaintenance    float64
	MiddlePrice          float64
	ContractWorkValue    float64
	BudgetProcurement    float64
	StrBudgetMaintenance string
	StrMiddlePrice       string
	StrContractWorkValue string
	StrBudgetProcurement string
	MaintenanceDetail    []MaintenanceDetailTracking
	ProgressDetail       []ProgressDetailTracking
}

type MaintenanceDetailTracking struct {
	ContractNumber          string
	No                      int
	RoadName                string
	Lane                    int
	MaintenanceStandardName string
	KmStart                 int
	KmEnd                   int
	StrKmStart              string
	StrKmEnd                string
	Distance                float64
	StrDistance             string
	SumDistance             float64
	StrSumDistance          string
	LastRow                 bool
}

type ProgressDetailTracking struct {
	ContractNumber string
	Schedule       time.Time
	StrSchedule    string

	ProgressPlan       float64
	SumProgressPlan    float64
	DisProgressPlan    interface{}
	DisSumProgressPlan interface{}

	Progress       float64
	SumProgress    float64
	DisProgress    interface{}
	DisSumProgress interface{}

	DisbursementPlan       float64
	SumDisbursementPlan    float64
	DisDisbursementPlan    interface{}
	DisSumDisbursementPlan interface{}

	Disbursement       float64
	SumDisbursement    float64
	DisDisbursement    interface{}
	DisSumDisbursement interface{}

	StrDisbursementPlan    string
	StrSumDisbursementPlan string

	StrDisbursement    string
	StrSumDisbursement string

	LastRow bool
}

// traffic volume report
type DataVolume struct {
	Year      []int           `json:"year"`
	RoadGroup []RoadGroupData `json:"road_group"`
}

type DataReportTrafficVolume struct {
	Code            string
	RoadGroupName   string
	RoadSectionName string
	RoadName        string
	KmStart         float64
	KmEnd           float64
	TotalKm         float64
	Veh1            int
	Veh2            int
	Veh3            int
	Total           int
	SurveyedDate    string
	Year            int
}

// accident volume report
type DataReportAccidentVolume struct {
	Code            string
	Name            string
	Acc1            int
	Acc2            int
	Acc3            int
	Acc4            int
	SurveyedDate    time.Time
	Sum             int
	StrAcc1         string
	StrAcc2         string
	StrAcc3         string
	StrAcc4         string
	StrSurveyedDate string
	StrSum          string
	IsNull          bool

	Year int
}

// summary of road condition report
type DataSummaryCondition struct {
	Type    []string        `json:"type"`
	Year    []int           `json:"year"`
	Road    []RoadGroupData `json:"road"`
	Measure []RoadGroupData `json:"measure"`
}

type DataReportSummaryCondition struct {
	RoadGroupName     string
	RoadSectionNumber string
	RoadSectionName   string
	KmStart           string
	KmEnd             string
	RoadLength        float64
	Year              int
	IsNull            bool

	Grade DataMeasureSummaryCondition

	Lane []DataLaneSummaryCondition

	SumRoadLengthStr string
	SumAvg           string
	SumCountA        string
	SumCountB        string
	SumCountC        string
	SumCountD        string
}

type DataMeasureSummaryCondition struct {
	A      DataGradeSummaryCondition
	B      DataGradeSummaryCondition
	C      DataGradeSummaryCondition
	D      DataGradeSummaryCondition
	IsNull bool
}
type DataGradeSummaryCondition struct {
	RefGradeID int
	LeftValue  float64
	RightValue float64
}
type DataLaneSummaryCondition struct {
	No        int
	Length    float64
	StrLength string
	Avg       float64
	StrAvg    string
	CountA    float64
	CountB    float64
	CountC    float64
	CountD    float64
	StrCountA string
	StrCountB string
	StrCountC string
	StrCountD string
}

// map of asset report
type DataAsset struct {
	AssetGroupID int
	Name         string
	AssetID      int
	TableLabel   string
}

type DataMap struct {
	Road       []RoadGroupData  `json:"road"`
	AssetGroup []DataAssetGroup `json:"group"`
}

type DataAssetGroup struct {
	ID    int             `json:"id"`
	Name  string          `json:"name"`
	Asset []RoadGroupData `json:"asset"`
}

type TableName struct {
	TableLabel string
	TableName  string
}

type Column struct {
	ColumnSeq      int
	ColumnName     string
	ComponentTitle string
	ComponentType  string
}

type DataReportMap struct {
	AssetName     string
	RoadGroupName string
	RoadName      string
	RoadCode      string
	RoadColorCode string
	KmStart       string
	KmEnd         string
	StrRoadLength string
	IsNull        bool

	Column []string
	Key    []string

	Row [][]interface{}

	RoadGeom  string
	PointGeom []string
	PinGeom   string
	LineGeom  []string

	Zoom int

	RoadID                       int
	RoadGroupNumber              string
	RoadSectionNumber            string
	RoadSectionDistance          float64
	RoadSectionNameOriginTh      string
	RoadSectionNameDestinationTh string
}

type AssetName struct {
	AssetName string `gorm:"column:table_label"`
}

type MapGeom struct {
	StringGeom string
}

// summary asset
type RefSummaryAsset struct {
	ID           int
	Name         string
	SummaryAsset []SummaryAsset `gorm:"ForeignKey:RefAssetID;AssociationForeignKey:ID"`
}

// summary asset
type SummaryAsset struct {
	ID               int
	TableLabel       string
	RefAssetID       int
	RoadAssetID      pq.Int32Array `json:"road_asset_id" gorm:"type:integer[]"`
	TableName        string
	FirstUpdatedDate time.Time
	LastUpdatedDate  time.Time
}
type SummaryAssetIDs struct {
	IDs         []int
	TableLabel  string
	RoadAssetID int
	TableName   string
	UpdatedDate time.Time
}

type TitleSummaryAsset struct {
	Title string
	Table []ListSummaryAsset
}

type ListSummaryAsset struct {
	Topic  string
	Header []string
	Row    [][]string
}

type CountSummaryAsset struct {
	Count int
}

type CountAndTypeSummaryAsset struct {
	Name  string
	Count int
}

type CountLightSummaryAsset struct {
	Name  string
	Type  string
	Count int
}

type DataReportSummaryAsset struct {
	RoadID                       int
	RoadGroupNumber              string
	RoadSectionNumber            string
	RoadSectionDistance          float64
	RoadSectionNameOriginTh      string
	RoadSectionNameDestinationTh string
	RoadMainName                 string

	RoadGroupName string
	RoadName      string
	RoadCode      string
	KmStart       string
	KmEnd         string
	StrRoadLength string
	IsNull        bool

	Table []TitleSummaryAsset
}

// asset adjustment report
type DataAssetAdjustment struct {
	Year  []int           `json:"year"`
	Month []string        `json:"month"`
	Road  []RoadGroupData `json:"road"`
}

type AssetYear struct {
	Year int
}

type Topic struct {
	ID            int
	TableLabel    string
	IDParent      int
	IDParentSlice []string
	TableName     string
}
type ListAssetAdjustment struct {
	Topic  string
	Header []string
	Row    [][]interface{}
}

type DataReportAssetAdjustment struct {
	Month         string
	Year          string
	RoadGroupName string
	RoadName      string
	RoadCode      string
	KmStart       string
	KmEnd         string
	StrRoadLength string
	IsNull        bool

	Table []ListAssetAdjustment
}

// summary of road surface report
type DataSurface struct {
	Year []int           `json:"year"`
	Road []RoadGroupData `json:"road"`
}

type DatabaseSurface struct {
	Typ     string
	Lane    int
	KmStart float64
	KmEnd   float64
}

type DataSurfaceDetail struct {
	Lane1 float64
	Lane2 float64
	Lane3 float64
	Lane4 float64
	Lane5 float64
	Sum   float64

	StrLane1 string
	StrLane2 string
	StrLane3 string
	StrLane4 string
	StrLane5 string
	StrSum   string
}

type SumSurface struct {
	StrLane1 string
	StrLane2 string
	StrLane3 string
	StrLane4 string
	StrLane5 string
	StrSum   string
}

type DataSurfaceType struct {
	PMA      DataSurfaceDetail
	AC       DataSurfaceDetail
	Slurry   DataSurfaceDetail
	Porous   DataSurfaceDetail
	Concrete DataSurfaceDetail
	Sum      SumSurface
}

type DataReportSurface struct {
	RoadGroupName string
	RoadName      string
	RoadCode      string
	KmStart       string
	KmEnd         string
	StrRoadLength string
	Year          string

	IsNull bool

	Surface DataSurfaceType
}

// ////////////// NEW MIMS ////////////////

type FilterAssetRoad struct {
	RoadSection
	RoadGroup RoadGroup `json:"road_group" gorm:"ForeignKey:RoadGroupId;AssociationForeignKey:Id"`
}

type FilterAsset struct {
	RefAsset
	RefAssetTable []RefAssetTable `json:"ref_asset_table" gorm:"ForeignKey:RefAssetID;AssociationForeignKey:ID"`
}

type FilterRoadCondition struct {
	RoadCondition
	Road FilterRoadSurfaceRoad `json:"road" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
}

type FilterRoadSurface struct {
	RoadSurface
	Road FilterRoadSurfaceRoad `json:"road" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
}

type FilterRoadRetroReflectivity struct {
	RoadRetroReflectivity
	Road FilterRoadSurfaceRoad `json:"road" gorm:"ForeignKey:RoadID;AssociationForeignKey:Id"`
}

type FilterRoadSurfaceRoad struct {
	Road
	RoadSection FilterRoadSurfaceRoadSection `json:"road_section" gorm:"ForeignKey:RoadSectionId;AssociationForeignKey:Id"`
}

type FilterRoadDamage struct {
	RoadDamage
	Road FilterRoadSurfaceRoad `json:"road" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
}

type FilterRoadDamageRoad struct {
	Road
	RoadSection FilterRoadSurfaceRoadSection `json:"road_section" gorm:"ForeignKey:RoadSectionId;AssociationForeignKey:Id"`
}

type FilterRoadSurfaceRoadSection struct {
	RoadSection
	RoadGroup RoadGroup `json:"road_group" gorm:"ForeignKey:RoadGroupId;AssociationForeignKey:Id"`
}

type FilterAadt struct {
	VolumeAadt
	Road FilterAadtRoad `json:"road" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
}

type FilterMaintenance struct {
	Maintenance
	MaintenanceRoad []FilterMaintenanceRoad `json:"maintenance_road" gorm:"ForeignKey:MaintenanceID;AssociationForeignKey:ID"`
}

type FilterMaintenanceRoad struct {
	MaintenanceRoad
	Road FilterAadtRoad `json:"road" gorm:"ForeignKey:RoadID;AssociationForeignKey:Id"`
}

type FilterAadtRoad struct {
	Road
	RoadSection FilterAadtRoadSection `json:"road_section" gorm:"ForeignKey:RoadSectionId;AssociationForeignKey:Id"`
}

type FilterAadtRoadSection struct {
	RoadSection
	RoadGroup RoadGroup `json:"road_group" gorm:"ForeignKey:RoadGroupId;AssociationForeignKey:Id"`
}

type ReportRetroReflectivityRoadSection struct {
	RoadSection           RoadSection
	RoadGroup             RoadGroup `json:"road_group" gorm:"ForeignKey:RoadGroupId;AssociationForeignKey:Id"`
	RoadRetroReflectivity []ReportRetroReflectivity
}

type ReportKpiRoadGroup struct {
	RoadSection RoadSection
	RoadGroup   RoadGroup `json:"road_group" gorm:"ForeignKey:RoadGroupId;AssociationForeignKey:Id"`
}

type ReportRetroReflectivity struct {
	RoadRetroReflectivity
	RoadInfo                   RoadInfo `json:"road_info"`
	RoadRetroReflectivityRange []RoadRetroReflectivityAllMeter
}

type RoadRetroReflectivityAllMeter struct {
	ID                           int      `gorm:"column:id"`
	RoadRetroReflectivityRangeID int      `gorm:"column:road_retro_reflectivity_range_id"`
	KmStart                      float64  `gorm:"column:km_start"`
	KmEnd                        float64  `gorm:"column:km_end"`
	RetroMin                     *float64 `gorm:"column:retro_min"`
	RetroMax                     *float64 `gorm:"column:retro_max"`
	RetroAvg                     *float64 `gorm:"column:retro_avg"`
	TheGeom                      string   `gorm:"column:the_geom"`
	RefStripeColorID             int      `gorm:"column:ref_stripe_color_id"`
	RefStripeTypeID              int      `gorm:"column:ref_stripe_type_id"`
}

type RoadSurfaceLaneMaximumLane struct {
	RoadId  int
	NumLane int
}

type ResultIri struct {
	ResultIri1000 []ResultIri1000 `json:"result_iri_1000"`
	ResultIri100  []ResultIri100  `json:"result_iri_100"`
}

type RoadConditionSurveyForKpi struct {
	Id              int
	RoadConditionId int
	RoadId          int
	LaneNo          int
	KmStart         int
	KmEnd           int
	Iri             float64
	Rut             float64
	Ifi             float64
	SurveyType      string
	SurveyDate      time.Time
	Year            int
}

type RoadG7100 struct {
	Id         int
	RoadId     int
	LineNo     int
	KmStart    int
	KmEnd      int
	Retro      float64
	NameStrip  string
	SurveyDate time.Time
	Year       int
}
type RoadSurveyIri100 struct {
	RoadConditionSurveyId int
	KmStart               int
	KmEnd                 int
	Iri                   float64
	SurveyType            string
}

type RoadConditionId struct {
	Id int
}

type result2 struct {
	ResultIri1000 []ResultIri1000 `json:"result_iri_1000"`
	ResultIri100  []ResultIri100  `json:"result_iri_100"`
}

type ResultIri1000 struct {
	RoadId          int
	RoadName        string
	LaneNo          int
	KmStart         string
	KmEnd           string
	Iri             string
	Comment         []CommentString
	RoadCode        string
	SectionCode     string
	SurveyType      string
	MaintExpireDate string
	IsExpire        bool
}
type ResultIri100 struct {
	RoadId          int
	RoadName        string
	LaneNo          int
	KmStart         string
	KmEnd           string
	Iri             string
	Comment         []CommentString
	RoadCode        string
	SectionCode     string
	SurveyType      string
	MaintExpireDate string
	IsExpire        bool
}
type ResultIfi100 struct {
	RoadId          int
	RoadName        string
	LaneNo          int
	KmStart         string
	KmEnd           string
	Ifi             string
	Comment         []CommentString
	RoadCode        string
	SectionCode     string
	SurveyType      string
	MaintExpireDate string
	IsExpire        bool
}
type ResultRut100 struct {
	RoadId          int
	RoadName        string
	LaneNo          int
	KmStart         string
	KmEnd           string
	Rut             string
	Comment         []CommentString
	RoadCode        string
	SectionCode     string
	SurveyType      string
	MaintExpireDate string
	IsExpire        bool
}
type ResultG7100 struct {
	RoadId          int
	RoadName        string
	LineNo          int
	KmStart         string
	KmEnd           string
	Retro           string
	Comment         []CommentString
	RoadCode        string
	SectionCode     string
	NameStrip       string
	MaintExpireDate string
	IsExpire        bool
}

type CommentString struct {
	ProjectEndDate    string
	KmStart           string
	KmEnd             string
	MaintenanceMethod string
}

type Comment struct {
	ProjectEndDate    time.Time
	KmStart           int
	KmEnd             int
	MaintenanceMethod string
}

type RoadInfoForKpi struct {
	RoadId       int
	Name         string
	RoadCode     string
	SectionCode  string
	RefDirection int
}

type Condition1000 struct {
	IriCondition1000Ac float64
	IriCondition1000Cc float64
}

type Iri100 struct {
	IriCondition100Ac float64
	IriCondition100Cc float64
}

type Rut100 struct {
	RutCondition float64
}

type Ifi100 struct {
	IfiConditionAc float64
	IfiConditionCc float64
}
type ParamCondition struct {
	Id               int
	RefOwnerId       int
	RefGradeId       int
	LeftValueCc      float64
	LeftConditionCc  string
	RightValueCc     float64
	RightConditionCc string
	ConditionType    string
	LeftValueAc      float64
	LeftConditionAc  string
	RightValueAc     float64
	RightConditionAc string
}
type G7Condition struct {
	Id                   int
	RefOwnerRoadLineId   int
	RefGradeId           int
	LeftValueYellow      float64
	LeftConditionYellow  string
	RightValueYellow     float64
	RightConditionYellow string
	LeftValueWhite       float64
	LeftConditionWhite   string
	RightValueWhite      float64
	RightConditionWhite  string
}

type MaintStruct struct {
	Name              string
	ProjectEndDate    time.Time
	KmStart           int
	KmEnd             int
	MaintenanceMethod string
}

//////////////// NEW MIMS ////////////////

type ReportTrafficVolumeHeader struct {
	RoadID          int     `json:"road_id" gorm:"column:road_id"`
	RoadGroupName   string  `json:"road_group_name" gorm:"column:road_group_name"`
	RoadSectionName string  `json:"road_section_name" gorm:"column:road_section_name"`
	RoadName        string  `json:"road_name" gorm:"column:road_name"`
	KmStart         float64 `json:"km_start" gorm:"column:km_start"`
	KmEnd           float64 `json:"km_end" gorm:"column:km_end"`
	TotalKm         float64 `json:"total_km" gorm:"column:total_km"`
}

// func (ReportTrafficVolumeHeader) TableName() string {
// 	return "road"
// }
