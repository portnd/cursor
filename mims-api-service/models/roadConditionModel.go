package models

import (
	"time"
)

type RoadCondition struct {
	ID               int       `gorm:"primaryKey;autoIncrement"`
	RoadId           int       `json:"road_id" gorm:"column:road_id"`
	LaneNo           int       `json:"lane_no" gorm:"column:lane_no"`
	Year             int       `json:"year" gorm:"column:year"`
	KmStart          float64   `json:"km_start" gorm:"column:km_start"`
	KmEnd            float64   `json:"km_end" gorm:"column:km_end"`
	SurveyedDate     time.Time `json:"surveyed_date" gorm:"column:surveyed_date"`
	CreatedBy        int       `json:"created_by" gorm:"column:created_by"`
	CreatedDate      time.Time `json:"created_date" gorm:"column:created_date"`
	UpdatedBy        int       `json:"updated_by" gorm:"column:updated_by"`
	UpdatedDate      time.Time `json:"updated_date" gorm:"column:updated_date"`
	IRI              *float64  `json:"iri" gorm:"column:iri"`
	MPD              *float64  `json:"mpd" gorm:"column:mpd"`
	RUT              *float64  `json:"rut" gorm:"column:rut"`
	IFI              *float64  `json:"ifi" gorm:"column:ifi"`
	Remarks          string    `json:"remarks" gorm:"column:remarks"`
	IRIInputFilePath string    `json:"iri_input_filepath" gorm:"column:iri_input_filepath"`
	Revision         int       `json:"revision" gorm:"column:revision"`
	Status           string    `json:"status" gorm:"column:status"`
	IDParent         int       `json:"id_parent" gorm:"column:id_parent"`
	ImgFilePath      string    `json:"img_filepath" gorm:"column:img_filepath" default:""`
	RejectReason     string    `json:"reject_reason" gorm:"column:reject_reason"`
}

func (rc *RoadCondition) TableName() string {
	return "road_condition"
}

type RoadConditionList struct {
	ID            int       `json:"id" `
	RoadId        int       `json:"roadId" `
	Year          int       `json:"year" `
	IDParent      int       `json:"id_parent" `
	Revision      int       `json:"revision" `
	DirectionId   int       `json:"direction_id" `
	DirectionName string    `json:"direaction_name" `
	LaneNo        int       `json:"lane_no" `
	SurveyedDate  time.Time `json:"surveyed_date" `
	RoadInfo      RoadInfo  `json:"road_info" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
}

func (rc *RoadConditionList) TableName() string {
	return "road"
}

type RetroReflectivityList struct {
	ID            int       `json:"id" `
	RoadId        int       `json:"roadId" `
	Year          int       `json:"year" `
	IDParent      int       `json:"id_parent" `
	Revision      int       `json:"revision" `
	DirectionId   int       `json:"direction_id" `
	DirectionName string    `json:"direaction_name" `
	LineNo        int       `json:"line_no" `
	SurveyedDate  time.Time `json:"surveyed_date" `
	RoadInfo      RoadInfo  `json:"road_info" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
}

func (rc *RetroReflectivityList) TableName() string {
	return "road"
}

type RoadConditionForCount struct {
	ID     int    `json:"id"`
	RoadId int    `json:"road_id"`
	Status string `json:"status"`
}

func (rc *RoadConditionForCount) TableName() string {
	return "road_condition"
}

type RoadConditionDetails struct {
	RoadID        int       `gorm:"column:road_id"`
	ID            int       `gorm:"column:id"`
	IDParent      int       `gorm:"column:id_parent"`
	UpdatedDate   time.Time `gorm:"column:updated_date"`
	UpdatedBy     int       `gorm:"column:updated_by"`
	Status        string    `gorm:"column:status"`
	StatusText    string    `gorm:"column:status_text"`
	Revision      int       `gorm:"column:revision"`
	DirectionID   int       `gorm:"column:direction_id"`
	DirectionName string    `gorm:"column:direction_name"`
	RoadTypeID    int       `gorm:"column:road_type_id"`
	ImgFilepath   string    `gorm:"column:img_filepath"`
	KmStartKM     int       `gorm:"column:km_start_km"`
	KmEndKM       int       `gorm:"column:km_end_km"`
	KmStartM      int       `gorm:"column:km_start_m"`
	KmEndM        int       `gorm:"column:km_end_m"`
	Value         *float64  `gorm:"column:value"`
	// GradeID       int       `gorm:"column:grade_id"`
	// GradeName     string    `gorm:"column:grade_name"`
	GeomCL       string `gorm:"column:geom_cl"`
	RejectReason string `gorm:"column:reject_reason"`
}

func (rc *RoadConditionDetails) TableName() string {
	return "road"
}

type SqlCondition struct {
	Select string `json:"select" `
	Join   string `json:"join" `
	Where  string `json:"where" `
}

// road_id	road_code	name	km_start	km_end	img_filepath	iri	mpd	rut	ifi	gn	load	speed	flow	sp	fr60
type RoadConditionCSV struct {
	RoadId      int      `json:"road_id"`
	RoadCode    string   `json:"road_code"`
	Name        string   `json:"name"`
	KMStart     float64  `json:"km_start"`
	KMEnd       float64  `json:"km_end"`
	IRI         *float64 `json:"iri"`
	MPD         *float64 `json:"mpd"`
	RUT         *float64 `json:"rut"`
	IFI         *float64 `json:"ifi"`
	SurveyType  string   `json:"survey_type"`
	ImgFilepath string   `json:"img_filepath"`
}

// type RoadConditionGrade struct {
// 	LeftValue      float64 `gorm:"column:left_value"`
// 	LeftCondition  string  `gorm:"column:left_condition"`
// 	RightValue     float64 `gorm:"column:right_value"`
// 	RightCondition string  `gorm:"column:right_condition"`
// 	CoditionType   string  `gorm:"column:condition_type"`
// 	GradeId        int     `gorm:"column:grade_id"`
// }

type FullGeom struct {
	Geom    string  `gorm:"column:the_geom"`
	KmStart float64 `gorm:"column:km_start"`
	KmEnd   float64 `gorm:"column:km_end"`
}

func (b *FullGeom) TableName() string {
	return "road_geom"
}

type RoadConditionAll struct {
	RoadCondition
	RoadConditionSurveys []RoadConditionSurveyPreload ` gorm:"ForeignKey:RoadConditionID;references:ID"`
}

func (rc *RoadConditionAll) TableName() string {
	return "road_condition"
}

type RoadPreloadConditionAll struct {
	RoadInfo
	RefDirection  RefDirection       `json:"direction" gorm:"ForeignKey:RefDirectionId;AssociationForeignKey:Id"`
	RoadCondition []RoadConditionAll ` gorm:"ForeignKey:RoadId;references:RoadId"`
}

func (rc *RoadPreloadConditionAll) TableName() string {
	return "road_info"
}

type RoadConditionTemplate struct {
	Road
	RoadInfo    RoadInfo        ` gorm:"ForeignKey:RoadId;references:Id"`
	RoadSection RoadSectionById ` gorm:"foreignKey:Id;references:RoadSectionId"`
	RoadGeom    []RoadGeom      ` gorm:"ForeignKey:RoadId;references:Id"`
}

func (rc *RoadConditionTemplate) TableName() string {
	return "road"
}

type RoadConditionCompare struct {
	RoadCondition
	RoadConditionSurveys []RoadConditionSurveyPreload ` gorm:"ForeignKey:RoadConditionId;references:Id"`
}

func (rc *RoadConditionCompare) TableName() string {
	return "road_condition"
}

type RoadConditionCompareLane struct {
	Years []int `form:"years"`
	Lanes []int `form:"lanes"`
}

type RoadRetroReflectivityCompareLine struct {
	Years []int `form:"years"`
	Lines []int `form:"lines"`
}

type RoadConditionAverage struct {
	Lane     int                        `json:"lane" extensions:"x-order=0"`
	Items    []RoadConditionAverageItem `json:"items"`
	Revision map[int]int
	IDParent map[int]int
}

type RoadRetroReflectivityAverage struct {
	Line     int                           `json:"line" extensions:"x-order=0"`
	Items    []RoadReflectivityAverageItem `json:"items"`
	Revision map[int]int
	IDParent map[int]int
}

type RoadReflectivityAverageItem struct {
	Year     int      `json:"year" extensions:"x-order=0"`
	KmStart  int      `json:"km_start" extensions:"x-order=1"`
	KmEnd    int      `json:"km_end" extensions:"x-order=2"`
	RetroAvg *float64 `json:"retro_avg" extensions:"x-order=3"`
}

type RoadConditionAverageItem struct {
	Year    int      `json:"year" extensions:"x-order=0"`
	KmStart int      `json:"km_start" extensions:"x-order=1"`
	KmEnd   int      `json:"km_end" extensions:"x-order=2"`
	IRI     *float64 `json:"iri" extensions:"x-order=3"`
	MPD     *float64 `json:"mpd" extensions:"x-order=4"`
	RUT     *float64 `json:"rut" extensions:"x-order=5"`
	IFI     *float64 `json:"ifi" extensions:"x-order=6"`
}

type RoadKmRage struct {
	RoadID   int
	RoadCode string
	RoadName string
	KmStart  int
	KmEnd    int
}

type RoadKmRage25M struct {
	RoadID     int
	RoadCode   string
	RoadName   string
	KmStart    float64
	KmEnd      float64
	SurveyType int
}

type RoadConditionSurveyDate struct {
	RoadId       int       `json:"road_id" gorm:"column:road_id"`
	LaneNo       int       `json:"lane_no" gorm:"column:lane_no"`
	SurveyedDate time.Time `json:"surveyed_date" gorm:"column:surveyed_date"`
}

func (rc *RoadConditionSurveyDate) TableName() string {
	return "road_condition"
}
