package responses

import "gitlab.com/mims-api-service/models"

type FilterRoadType1 struct {
	FilterRoad []FilterRoadGroupNoSection `json:"filter_road"`
}

type FilterRoadType2 struct {
	FilterRoad []FilterRoadYear `json:"filter_road"`
}

type FilterRoadType3 struct {
	FilterRoad  []FilterRoadYear `json:"filter_road"`
	FilterRange []int            `json:"filter_range"`
}

type FilterRoadType4 struct {
	FilterRoad      []FilterRoadYear `json:"filter_road"`
	FilterCondition []string         `json:"filter_condition"`
	FilterCriteria  []FilterRefOwner `json:"filter_criteria"`
}

type FilterRoadType5 struct {
	FilterRoad     []FilterRoadYear     `json:"filter_road"`
	FilterCriteria []FilterRefOwnerLine `json:"filter_criteria"`
}

type FilterRoadType6 struct {
	FilterRoad     []FilterRoadYear     `json:"filter_road"`
	FilterCriteria []FilterRefOwnerLine `json:"filter_criteria"`
}

type FilterRoadDamageType1 struct {
	FilterRoad []FilterRoadYear `json:"filter_road"`
}

type FilterRoadDamageType2 struct {
	FilterRoad []FilterRoadYear `json:"filter_road"`
}

type FilterAadtType1 struct {
	FilterRoad []FilterAadtYear `json:"filter_road"`
}

type FilterMaintenanceKpiType1 struct {
	FilterCondition []FilterMaintenanceConditionKpi `json:"filter_condition"`
}

type FilterMaintenanceConditionKpi struct {
	Name string                     `json:"name"`
	Year []FilterMaintenanceYearKpi `json:"year"`
}

type FilterMaintenanceYearKpi struct {
	Year      int               `json:"year"`
	RoadGroup []FilterRoadGroup `json:"road_group"`
}

type FilterCondition struct {
	Name       string           `json:"name"`
	FilterYear []FilterRoadYear `json:"filter_year"`
}

type FilterYear struct {
	Id              int               `json:"id"`
	Name            string            `json:"name"`
	FilterRoadGroup []FilterRoadGroup `json:"filter_road_group"`
}

type FilterMaintenanceType1 struct {
	FilterRoad []FilterRoadGroup     `json:"filter_road"`
	FilterYear FilterMaintenanceYear `json:"filter_Year"`
}

type FilterMaintenanceYear struct {
	StartYear []int `json:"start_year"`
	EndYear   []int `json:"end_year"`
}

type FilterRefOwner struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type FilterRefOwnerLine struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type FilterRoadYear struct {
	Year  int           `json:"year"`
	Depot []FilterDepot `json:"depot"`
}

type FilterAssetType1 struct {
	FilterRoad  []FilterDepot   `json:"filter_road"`
	FilterAsset []AssetRefAsset `json:"filter_asset"`
}

type FilterAssetType2 struct {
	FilterRoad  []FilterDepot   `json:"filter_road"`
	FilterAsset []AssetRefAsset `json:"filter_asset"`
}

type FilterAssetType3 struct {
	FilterRoad []FilterDepot `json:"filter_road"`
}

type AssetRefAsset struct {
	Id    int                  `json:"id"`
	Name  string               `json:"name"`
	Asset []AssetRefAssetTable `json:"asset"`
}

type AssetRefAssetTable struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type FilterDepot struct {
	Id        int               `json:"id"`
	Name      string            `json:"name"`
	RoadGroup []FilterRoadGroup `json:"road_group"`
}

type FilterRoadGroup struct {
	Id          int                 `json:"id"`
	Name        string              `json:"name"`
	RoadSection []FilterRoadSection `json:"road_section"`
}

type FilterRoadSection struct {
	Id     int    `json:"id"`
	Number int    `json:"-"`
	Name   string `json:"name"`
}

type FilterRoadGroupNoSection struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type FilterAadtYear struct {
	Year      int               `json:"year"`
	RoadGroup []FilterRoadGroup `json:"road_group"`
}

type ReportStatusId struct {
	Id int `json:"id"`
}

type ReportStatus struct {
	IsFinish bool   `json:"is_finish"`
	Path     string `json:"path"`
}

type Report8 struct {
	RoadGroupNumber   string        `json:"road_group_number"`
	RoadSectionNumber string        `json:"road_section_number"`
	RoadSectionName   string        `json:"road_section_Name"`
	RoadGroupName     string        `json:"road_group_name"`
	KmStart           string        `json:"km_start"`
	KmEnd             string        `json:"km_end"`
	Distance          string        `json:"distance"`
	Road              []Report8Road `json:"road"`
}

type Report8Road struct {
	RoadId int                  `json:"road_id"`
	Name   string               `json:"name"`
	Table  []Report8TableLevel1 `json:"table"`
}

type Report8TableLevel1 struct {
	Line                 int                  `json:"line"`
	Color                string               `json:"color"`
	Type                 string               `json:"type"`
	Distance             float64              `json:"distance"`
	Sub                  []Report8TableLevel2 `json:"sub"`
	RetroReflectivityAvg string               `json:"retro_reflectivity_avg"`
}

type Report8TableLevel2 struct {
	KmStart           string  `json:"km_start"`
	KmEnd             string  `json:"km_end"`
	Distance          float64 `json:"distance"`
	RetroReflectivity string  `json:"retro_reflectivity"`
}

type Report9 struct {
	RoadGroupNumber   string              `json:"road_group_number"`
	RoadSectionNumber string              `json:"road_section_number"`
	RoadSectionName   string              `json:"road_section_Name"`
	RoadGroupName     string              `json:"road_group_name"`
	KmStart           string              `json:"km_start"`
	KmEnd             string              `json:"km_end"`
	Distance          string              `json:"distance"`
	Criteria          []Report9Criteria   `json:"criteria"`
	Graph             Report9Graph        `json:"graph"`
	Table             []Report9Table      `json:"table"`
	TableSummary      Report9TableSummary `json:"table_summary"`
	MapCenter         string              `json:"map_center"`
	Map               []Report9Map        `json:"map"`
	Year              string              `json:"year"`
}

type Report9Map struct {
	Color   string   `json:"color"`
	TheGeom []string `json:"the_grom"`
}

type Report9Criteria struct {
	ColorName     string                 `json:"color_name"`
	CriteriaValue []Report9CriteriaValue `json:"criteria_value"`
}

type Report9CriteriaValue struct {
	Criteria       string  `json:"criteria"`
	OperatorLeft   string  `json:"operator_left"`
	OperationLeft  float64 `json:"operation_left"`
	OperatorRight  string  `json:"operator_right"`
	OperationRight float64 `json:"operation_right"`
}

type Report9TableSummary struct {
	G7Avg   float64 `json:"g7_avg"`
	Pass    float64 `json:"pass"`
	NotPass float64 `json:"not_pass"`
}

type Report9Table struct {
	Line    int     `json:"line"`
	Color   string  `json:"color"`
	G7Avg   float64 `json:"g7_avg"`
	Pass    float64 `json:"pass"`
	NotPass float64 `json:"not_pass"`
}

type Report9Graph struct {
	Lable []string  `json:"lable"`
	Color []string  `json:"color"`
	Value []float64 `json:"value"`
}

type ReportTrafficVolume struct {
	RoadID          int    `json:"-"`
	RoadGroupName   string `json:"road_group_name"`
	RoadSectionName string `json:"road_section_name"`
	RoadName        string `json:"road_name"`
	KmStart         string `json:"km_start"`
	KmEnd           string `json:"km_end"`
	TotalKm         string `json:"total_km"`
	Veh1            string `json:"veh1"`
	Veh2            string `json:"veh2"`
	Veh3            string `json:"veh3"`
	Total           string `json:"total"`
	SurveyedDate    string `json:"surveyed_date"`
	Year            string `json:"year"`
}

type Report12Iri struct {
	RoadGroupNumber    string           `json:"road_group_number"`
	RoadSectionNumber  string           `json:"road_section_number"`
	RoadSectionName    string           `json:"road_section_Name"`
	RoadGroupName      string           `json:"road_group_name"`
	KmStart            string           `json:"km_start"`
	KmEnd              string           `json:"km_end"`
	Distance           string           `json:"distance"`
	Table              models.ResultIri `json:"table"`
	LastRowOnPage1000M int              `json:"last_row_on_page_1000m"`
	LastRowOnPage100M  int              `json:"last_row_on_page_100m"`
}

type Report12Ifi struct {
	RoadGroupNumber   string                `json:"road_group_number"`
	RoadSectionNumber string                `json:"road_section_number"`
	RoadSectionName   string                `json:"road_section_Name"`
	RoadGroupName     string                `json:"road_group_name"`
	KmStart           string                `json:"km_start"`
	KmEnd             string                `json:"km_end"`
	Distance          string                `json:"distance"`
	Table             []models.ResultIfi100 `json:"table"`
	LastRowOnPage     int                   `json:"last_row_on_page"`
}

type Report12Rut struct {
	RoadGroupNumber   string                `json:"road_group_number"`
	RoadSectionNumber string                `json:"road_section_number"`
	RoadSectionName   string                `json:"road_section_Name"`
	RoadGroupName     string                `json:"road_group_name"`
	KmStart           string                `json:"km_start"`
	KmEnd             string                `json:"km_end"`
	Distance          string                `json:"distance"`
	Table             []models.ResultRut100 `json:"table"`
	LastRowOnPage     int                   `json:"last_row_on_page"`
}

type Report12G7 struct {
	RoadGroupNumber   string               `json:"road_group_number"`
	RoadSectionNumber string               `json:"road_section_number"`
	RoadSectionName   string               `json:"road_section_Name"`
	RoadGroupName     string               `json:"road_group_name"`
	KmStart           string               `json:"km_start"`
	KmEnd             string               `json:"km_end"`
	Distance          string               `json:"distance"`
	Table             []models.ResultG7100 `json:"table"`
	LastRowOnPage     int                  `json:"last_row_on_page"`
}
