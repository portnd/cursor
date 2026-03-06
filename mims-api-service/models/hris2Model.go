package models

import (
	"time"

	"github.com/lib/pq"
)

type SectionGeom struct {
	Item []Item `xml:"Item" json:"item"`
}

type Item struct {
	IsLatested               bool   `json:"is_latested" bson:"is_latested"`
	RoadNumber               string `xml:"RoadNumber" json:"road_number" bson:"road_number"`
	SectionRoadNumber        string `xml:"SectionRoadNumber" json:"section_road_number" bson:"section_road_number"`
	SectionRoadThName        string `xml:"SectionRoadThName" json:"section_road_th_name" bson:"section_road_th_name"`
	SectionRoadEngName       string `xml:"SectionRoadEngName" json:"section_road_eng_name" bson:"section_road_eng_name"`
	SubsectionRoadNumber     string `xml:"SubsectionRoadNumber" json:"subsection_road_number" bson:"subsection_road_number"`
	KmBegin                  string `xml:"KmBegin" json:"km_begin" bson:"km_begin"`
	KmEnd                    string `xml:"KmEnd" json:"km_end" bson:"km_end"`
	Distance                 string `xml:"Distance" json:"distance" bson:"distance"`
	DepotCode                string `xml:"DepotCode" json:"depot_code"  bson:"depot_code"`
	NameDepotCode            string `xml:"NameDepotCode" json:"name_depot_code"  bson:"name_depot_code"`
	HighwayDistrictCode      string `xml:"HighwayDistrictCode" json:"highway_district_code"  bson:"highway_district_code"`
	NameHighwayDistrictCode  string `xml:"NameHighwayDistrictCode" json:"name_highway_district_code"  bson:"name_highway_district_code"`
	OfficeOfHighwaysCode     string `xml:"OfficeOfHighwaysCode" json:"office_of_highways_code"  bson:"office_of_highways_code"`
	NameOfficeOfHighwaysCode string `xml:"NameOfficeOfHighwaysCode" json:"name_office_of_highways_code"  bson:"name_office_of_highways_code"`
	ProvinceCord             string `xml:"ProvinceCord" json:"province_cord"  bson:"province_cord"`
	DistrictCord             string `xml:"DistrictCord" json:"district_cord"  bson:"district_cord"`
	SubDistrictCord          string `xml:"SubDistrictCord" json:"sub_district_cord"  bson:"sub_district_cord"`
	RegionCord               string `xml:"RegionCord" json:"region_cord"  bson:"region_cord"`
	SurveyDate               string `xml:"SurveyDate" json:"survey_date"  bson:"survey_date"`
	StatusRoad               string `xml:"StatusRoad" json:"status_road"  bson:"status_road"`
	SectionPartID            string `xml:"section_part_id" json:"section_part_id" bson:"section_part_id"`
	// TheGeom                  string `xml:"the_geom" json:"the_geom"`
}

type RoadLatest struct {
	IsLatested bool   `json:"is_latested" bson:"is_latested"`
	RoadID     string `xml:"road_id"  json:"road_id" bson:"road_id"`
	RoadCode   string `xml:"road_code"  json:"road_code" bson:"road_code"`
	RoadName   string `xml:"road_name"  json:"road_name" bson:"road_name"`
	KmStart    string `xml:"km_start"  json:"km_start" bson:"km_start"`
	KmEnd      string `xml:"km_end"  json:"km_end" bson:"km_end"`
	Length     string `xml:"length"  json:"length" bson:"length"`
	Revision   string `xml:"revision"  json:"revision" bson:"revision"`
	MostRecent string `xml:"most_recent"  json:"most_recent" bson:"most_recent"`
	Status     string `xml:"status"  json:"status" bson:"status"`
}

type RefHris struct {
	Id                   int       `json:"id"`
	RoadNumber           string    `json:"road_number"`
	OfficeOfHighwaysCode string    `json:"office_of_highways_code"`
	SectionRoadNumber    string    `json:"section_road_number"`
	Status               bool      `json:"status"`
	IsDeleted            bool      `json:"is_deleted"`
	CreatedBy            int       `json:"created_by"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedBy            int       `json:"updated_by"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type InsertRoadGroup struct {
	Id        int     `json:"id"`
	Number    string  `json:"number"`
	Name      string  `json:"name"`
	ShortName string  `json:"short_name"`
	KmStart   float32 `json:"km_start"`
	KmEnd     float32 `json:"km_end"`
	Distance  float32 `json:"distance"`
}

type InsertRoadSection struct {
	Id                int           `json:"id"`
	RoadGroupId       int           `json:"road_group_id"`
	Number            string        `json:"number"`
	NameOriginTH      string        `json:"name_origin_th"`
	NameDestinationTH string        `json:"name_destination_th"`
	NameOriginEn      string        `json:"name_origin_en"`
	NameDestinationEn string        `json:"name_destination_en"`
	KmStart           float32       `json:"km_start"`
	KmEnd             float32       `json:"km_end"`
	Distance          float32       `json:"distance"`
	ProvinceCode      pq.Int64Array `json:"province_code" gorm:"type:int[]"`
	RefDivisionCode   string        `json:"-"`
	RefDivision       RefDivision   `json:"ref_division"  gorm:"foreignKey:RefDivisionCode; references:DivisionCode"`
	RefDistrictCode   string        `json:"-"`
	RefDistrict       RefDistrict   `json:"ref_district"  gorm:"foreignKey:RefDistrictCode; references:DistrictCode"`
	RefDepotCode      string        `json:"-"`
	RefDepot          RefDepot      `json:"ref_depot"  gorm:"foreignKey:RefDepotCode; references:DepotCode"`
}

func (b *InsertRoadGroup) TableName() string {
	return "road_group"
}

func (b *InsertRoadSection) TableName() string {
	return "road_section"
}
