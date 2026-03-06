package responses

import "time"

type RoadDashboard struct {
	RoadLabel  RoadLabel    `json:"road"`
	LengthRoad []LengthRoad `json:"length_roads"`
	AadtRoad   []AadtRoad   `json:"aadt_roads"`
}

type RoadLabel struct {
	Name        string   `json:"name"`
	Label       []string `json:"label"`
	Data        []int    `json:"data"`
	Color       []string `json:"color"`
	RoadGroupID []int    `json:"road_group_id"`
}

type LengthRoad struct {
	Name     string  `json:"name"`
	Total    float64 `json:"total"`
	Asphalt  float64 `json:"asphalt"`
	Concrete float64 `json:"concrete"`
	Color    string  `json:"color"`
}

type AadtRoad struct {
	Name       string  `json:"name"`
	Aadt       int     `json:"aadt"`
	Year1      int     `json:"year1"`
	Year2      int     `json:"year2"`
	Percent    float64 `json:"percent"`
	GrowthRate string  `json:"growth_rate"`
}

type PavementSurface struct {
	Length       float64 `json:"length"`
	SurfaceGroup string  `json:"surface_group"`
	Percentage   float64 `json:"percentage"`
}

type RoadGroupDashboard struct {
	ID            int     `json:"id"`
	TotalRoad     int     `json:"total_road"`
	RoadGroupName string  `json:"road_group_name"`
	Distance      float64 `json:"distance"`
}

type VolumeAADTDashboard struct {
	Year          int    `json:"year"`
	RoadGroupName string `json:"road_group_name"`
	Total         int    `json:"total"`
}

type DashboardAsset struct {
	ID         int      `json:"id"`
	RefAssetID int      `json:"ref_asset_id"`
	Name       string   `json:"name"`
	Label      []string `json:"label"`
	Data       []int    `json:"data"`
}

type DashboardYearMaxMin struct {
	MaxYear int `gorm:"column:overall_max_year"`
	MinYear int `gorm:"column:overall_min_year"`
}

type DashboardYear struct {
	Year int `json:"year"`
}

type DataMartCheck struct {
	Stauts    bool      `json:"stauts"`
	Percent   float64   `json:"percent"`
	UpdatedBy string    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ConditionDashboard struct {
	HasMutipleRoad bool             `json:"has_mutiple_road"`
	Chart          ConditionChart   `json:"chart"`
	Table          []ConditionTable `json:"table"`
}

type ConditionDashboardStr struct {
	HasMutipleRoad bool                `json:"has_mutiple_road"`
	Chart          ConditionChart      `json:"chart"`
	Table          []ConditionTableStr `json:"table"`
}

type ConditionPieChart struct {
	Name  string   `json:"name"`
	Lable []string `json:"lable"`
	Data  []int    `json:"data"`
	Km    []string `json:"km"`
}

type ConditionChart struct {
	Name   string    `json:"name"`
	Lable  []string  `json:"lable"`
	Data   []float64 `json:"data"`
	Color  []string  `json:"color"`
	RoadID []string  `json:"road_id"`
}

type NumberMaintenanceChart struct {
	Name        string   `json:"name"`
	Lable       []string `json:"lable"`
	Data        []int    `json:"data"`
	Color       []string `json:"color"`
	RoadGroupID []string `json:"road_group_id"`
}

type MaintenanceBudgetChart struct {
	Name        string    `json:"name"`
	Lable       []string  `json:"lable"`
	Data        []float64 `json:"data"`
	Color       []string  `json:"color"`
	RoadGroupID []string  `json:"road_group_id"`
}

type ConditionTable struct {
	LaneNo   int                 `json:"lane_no"`
	TotalKm  float64             `json:"total_km"`
	AvgValue float64             `json:"avg_value"`
	DetailKm []ConditionDetailKm `json:"detail_km"`
}

type ConditionTableStr struct {
	LaneNo   int                    `json:"lane_no"`
	TotalKm  float64                `json:"total_km"`
	AvgValue float64                `json:"avg_value"`
	DetailKm []ConditionDetailKmStr `json:"detail_km"`
}

type ConditionDetailKm struct {
	RefGradeID   int     `json:"ref_grade_id"`
	RefGradeName string  `json:"ref_grade_name"`
	Value        float64 `json:"value"`
	ValuePercent float64 `json:"value_percent"`
}

type ConditionDetailKmStr struct {
	RefGradeID   int    `json:"ref_grade_id"`
	RefGradeName string `json:"ref_grade_name"`
	Value        string `json:"value"`
	ValuePercent string `json:"value_percent"`
}

type DashboardConditionMap struct {
	Color   string   `json:"color"`
	TheGeom GeomJSON `json:"the_geom"`
}

type TopTenMaintenanceBudgetChart struct {
	Name          string    `json:"name"`
	Lable         []string  `json:"lable"`
	Data          []float64 `json:"data"`
	Color         []string  `json:"color"`
	MaintenanceID []int     `json:"maintenance_id"`
}

type MaintenanceTable struct {
	ID                      int      `json:"id"`
	ContractNumber          string   `json:"contract_number"`
	RoadName                []string `json:"road_name"`
	SectionName             []string `json:"section_name"`
	RefDepotName            []string `json:"ref_depot_name"`
	Budget                  float64  `json:"budget"`
	GuaranteeExpirationDate string   `json:"guarantee_expiration_date"`
	RemainDate              int      `json:"remain_date"`
}

type MaintenanceDashboard struct {
	UpdatedAt                    time.Time                    `json:"updated_at"`
	NumberMaintenanceChart       NumberMaintenanceChart       `json:"number_maintenance_chart"`
	MaintenanceBudgetChart       MaintenanceBudgetChart       `json:"maintenance_budget_chart"`
	TopTenMaintenanceBudgetChart TopTenMaintenanceBudgetChart `json:"top_ten_maintenance_budget_chart"`
}

type MaintenanceMapDashboard struct {
	IDParent       int         `json:"id_parent"`
	Title          string      `json:"title"`
	RoadName       string      `json:"road_name"`
	SectionName    string      `json:"section_name"`
	RefDepotName   string      `json:"ref_depot_name"`
	ContractNumber string      `json:"contract_number"`
	LaneNo         int         `json:"lane_no"`
	Name           string      `json:"name"`
	KmStart        string      `json:"km_start"`
	KmEnd          string      `json:"km_end"`
	KmTotal        float64     `json:"km_total"`
	Color          string      `json:"color"`
	TheGeom        TheGeomJson `json:"the_geom"`
}

type TheGeomJson struct {
	Coordinates [][]float64 `json:"coordinates"`
	Type        string      `json:"type"`
}

type TheGeomJsonPoint struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}
