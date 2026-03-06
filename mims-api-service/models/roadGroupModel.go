package models

import "github.com/lib/pq"

// Todo ...
type RoadGroup struct {
	Id        int     `json:"id"`
	Number    string  `json:"number"`
	Code      string  `json:"-"`
	Name      string  `json:"name"`
	ShortName string  `json:"short_name"`
	KmStart   float32 `json:"km_start"`
	KmEnd     float32 `json:"km_end"`
	Distance  float32 `json:"distance"`
}

type RoadGroupReport struct {
	Id          int     `json:"id"`
	Number      string  `json:"number"`
	Code        string  `json:"-"`
	Name        string  `json:"name"`
	ShortName   string  `json:"short_name"`
	KmStart     string  `json:"km_start"`
	KmEnd       string  `json:"km_end"`
	Distance    float64 `json:"distance"`
	StrDistance string
}

type RoadGroupByID struct {
	ID             int               `json:"id" gorm:"column:id"`
	RoadsGroupName string            `json:"roads_group_name" gorm:"column:name"`
	Roads          []RoadInRoadGroup `json:"roads"`
}
type RoadInRoadGroup struct {
	ID       int        `json:"id" gorm:"column:id"`
	RoadName string     `json:"name" gorm:"column:name"`
	Lanes    []RoadLane `json:"lanes"`
}

type RoadLane struct {
	Lane    int     `json:"lane" gorm:"column:lane_no"`
	KmStart float64 `json:"km_start" gorm:"column:km_start"`
	KmEnd   float64 `json:"km_end" gorm:"column:km_end"`
}

type RoadGroupInitData struct {
	Id               int            `json:"id"`
	Number           string         `json:"number"`
	ShortName        string         `json:"short_name"`
	RefDivisionCodes pq.StringArray `json:"ref_division_codes" gorm:"type:character[]"`
	RefDistrictCodes pq.StringArray `json:"ref_district_codes" gorm:"type:character[]"`
}

type RoadGroupInit struct {
	Id         int    `json:"id"`
	RoadNumber string `json:"road_number"`
	Prefix     string `json:"prefix"`
	ShortName  string `json:"short_name"`
}

// TableName use to specific table
func (b *RoadGroup) TableName() string {
	return "road_group"
}

func (b *RoadGroupInitData) TableName() string {
	return "road_group"
}

func (b *RoadGroupInit) TableName() string {
	return "road_group"
}

func (b *RoadGroupReport) TableName() string {
	return "road_group"
}
