package responses

import "gitlab.com/mims-api-service/models"

type SurfaceReportRespond struct {
	RoadGroupName     string
	RoadSectionNumber string
	RoadSectionName   string
	KmStart           string
	KmEnd             string
	StrRoadLength     string
	Year              string
	IsNull            bool
	PieChart          PieChart              `json:"pie_chart"`
	Table             []SurfaceReportTable  `json:"table"`
	TableSum          SurfaceLaneTypeReport `json:"table_sum"`
}

type PieChart struct {
	Series []float64 `json:"series"`
	Colors []string  `json:"colors"`
	Labels []string  `json:"labels"`
}

type SurfaceReportTable struct {
	SurfaceName     string                `json:"surface_name"`
	SurfaceLaneType SurfaceLaneTypeReport `json:"surface_lane_type"`
	ColorCode       string                `json:"color_code"`
}

type SurfaceLaneTypeReport struct {
	OneLane      string `json:"one_lane"`
	TwoLane      string `json:"two_lane"`
	ThreeLane    string `json:"three_lane"`
	FourLane     string `json:"four_lane"`
	MoreThanFour string `json:"more_than_four"`
	Sum          string `json:"sum"`
}

type Report7 struct {
	ConditionDashboardStr
	Header            string
	Year              int
	Type              string
	RoadGroupName     string
	RoadSectionNumber string
	RoadSectionName   string
	KmStart           string
	KmEnd             string
	RoadLength        string
	Summary           []string
	Grade             []models.ParamsConditionPreload
	MapData           interface{}
}

type Report6 struct {
	Data              []models.DataReportCondition
	Header            string
	Year              int
	Type              string
	RoadGroupName     string
	RoadSectionNumber string
	RoadSectionName   string
	KmStart           string
	KmEnd             string
	RoadLength        string
	IsNull            bool
}
