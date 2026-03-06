package responses

import "gitlab.com/mims-api-service/models"

type DetailKm struct {
	RefGrade     models.RefGrade
	Value        float64
	ValuePercent float64
}

type RefGrade struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type SummaryCondition struct {
	TotalKm  float64    `json:"total_km"`
	AvgValue float64    `json:"avg_value"`
	DetailKm []DetailKm `json:"detail_km"`
}

type DetailCondition struct {
	LaneNO int `json:"lane_no"`
	SummaryCondition
}

type ConditionSumRespond struct {
	SummaryCondition `json:"summary"`
	DetailCondition  []DetailCondition `json:"detail"`
}

type ConditionGeomSumRespond struct {
	TotalPage int `json:"total_page"`
	// Offset    int        `json:"offset"`
	Page     int        `json:"page"`
	GeomList []GeomList `json:"geom_list"`
}

type ConditionTotalPageRespond struct {
	TotalPage int `json:"total_page"`
	Page      int `json:"page"`
}

type GeomList struct {
	// ID int `json:"id"`
	// Grade  models.RefGrade `json:"grade"`
	GradeID int         `json:"grade_id"`
	Color   string      `json:"color"`
	GeomCL  interface{} `json:"geom_cl"`
}

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type SurfaceRespond struct {
	Summary               []models.Summary        `json:"summary"`
	SurfaceDashboardTable []SurfaceDashboardTable `json:"surface_dashboard_table"`
	Detail                models.Detail           `json:"-"`
	// GeomList              []models.GeomList       `json:"geom_list"`
}

type SurfaceDashboardTable struct {
	ID              int             `json:"id"`
	SurfaceName     string          `json:"surface_name"`
	SurfaceLaneType SurfaceLaneType `json:"surface_lane_type"`
	// ColorCode       string          `json:"color_code"`
}

type SurfaceLaneType struct {
	OneLane      float64 `json:"one_lane"`
	TwoLane      float64 `json:"two_lane"`
	ThreeLane    float64 `json:"three_lane"`
	FourLane     float64 `json:"four_lane"`
	MoreThanFour float64 `json:"more_than_four"`
}

type AssetLocationRespond struct {
	Code    string `json:"code"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Data    []models.AssetLocation
}

type AssetRespond struct {
	AssetGroup AssetGroup  `json:"asset_group"`
	AssetList  []AssetList `json:"asset_list"`
}

type AssetGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AssetList struct {
	Asset   Asset `json:"asset"`
	Value   int   `json:"value"`
	IsRange bool  `json:"is_range"`
}

type Asset struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	DefaultIconURL   string `json:"default_icon_url"`
	ThumbnailIconURL string `json:"thumbnail_icon_url"`
	DefaultColor     string `json:"default_color"`
}
