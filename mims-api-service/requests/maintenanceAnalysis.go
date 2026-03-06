package requests

type MaintenanceAnalysisQueryParams struct {
	Page              string `form:"page"`
	Limit             string `form:"limit"`
	AnalysisType      int    `json:"analysis_type"`
	AnalysisCondition string `json:"analysis_condition"`
}

type MaintenanceAnalysisStrategic struct {
	MaintenanceAnalysisId              int                       `json:"maintenance_analysis_id"`
	RefMaintenanceAnalysisBudgetTypeId int                       `json:"ref_maintenance_analysis_budget_type_id"` // เงื่อนไข
	RefMaintenanceAnalysisTargetTypeId int                       `json:"ref_maintenance_analysis_target_type_id"`
	Discount                           float64                   `json:"discount"`
	NumPlan                            int                       `json:"num_plan"`
	RangeYear                          int                       `json:"range_year"`
	Comment                            string                    `json:"comment"`
	Budget                             float64                   `json:"budget"`
	Iri                                float64                   `json:"iri"`
	Ifi                                float64                   `json:"ifi"`
	PlanList                           []MaintenanceAnalysisPlan `json:"plan_list"`
}

type MaintenanceAnalysisPlan struct {
	ID       *int     `json:"id"`
	Plan1    *float64 `json:"plan_1"`
	Plan2    *float64 `json:"plan_2"`
	Plan3    *float64 `json:"plan_3"`
	PlanYear float64  `json:"plan_year"`
}

type UpdateMaintenanceAnalysisPlan struct {
	Id         int     `json:"id"`
	PlanNumber float64 `json:"plan_number"`
	PlanYear   float64 `json:"plan_year"`
	PlanValue  float64 `json:"plan_value"`
}

type MaintenanceAnalysis struct {
	MaintenanceAnalysisTypeId int      `json:"maintenance_analysis_type_id"`
	SurfaceTypeId             int      `json:"surface_type_id"`
	LaneTypeId                int      `json:"lane_type_id"`
	Roads                     []int    `json:"roads"`
	Iri1                      *float64 `json:"iri1"`
	Iri2                      *float64 `json:"iri2"`
	Aadt1                     *float64 `json:"aadt1"`
	Aadt2                     *float64 `json:"aadt2"`
	Age1                      *float64 `json:"age1"`
	Age2                      *float64 `json:"age2"`
	Ifi1                      *float64 `json:"ifi1"`
	Ifi2                      *float64 `json:"ifi2"`
	GroupKm                   int      `json:"group_km"`
	Name                      string   `json:"name"`
	// List                      []MaintenanceAnalysisRoad `json:"lists"`
}

type MaintenanceAnalysisRoad struct {
	Id                    int     `json:"id"`
	MaintenanceAnalysisId int     `json:"maintenance_analysis_id"`
	RoadGroupId           int     `json:"road_group_id"`
	RoadId                int     `json:"road_id"`
	KmStart               float64 `json:"km_start"`
	KmEnd                 float64 `json:"km_end"`
	LaneTypeId            int     `json:"lane_type_id"`
	Iri                   float64 `json:"iri"`
	Aadt                  float64 `json:"aadt"`
	Gn                    float64 `json:"gn"`
	IsSelected            bool    `json:"is_selected"`
}

type UpdateMaintenanceAnalysisRoad struct {
	Id                    int     `json:"id"`
	MaintenanceAnalysisId int     `json:"maintenance_analysis_id"`
	RoadGroupId           int     `json:"road_group_id"`
	RoadId                int     `json:"road_id"`
	KmStart               float64 `json:"km_start"`
	KmEnd                 float64 `json:"km_end"`
	LaneTypeId            int     `json:"lane_type_id"`
	Iri                   float64 `json:"iri"`
	Aadt                  float64 `json:"aadt"`
	Gn                    float64 `json:"gn"`
	IsSelected            bool    `json:"is_selected"`
}

type AnalyzingReq struct {
	ConditionID *int   `json:"condition_id"`
	SurfaceType string `json:"surface_type"`
	// TotalKm     float64                            `json:"total_km"`
	// IRI        float64                            `json:"iri"`
	// GN         float64                            `json:"gn"`
	Name          string                    `json:"name"`
	Discount      *float64                  `json:"discount"`
	Year          *int                      `json:"year"`
	Target        *int                      `json:"target"`
	NumberPlan    *int                      `json:"number_plan"`
	Comment       *string                   `json:"comment"`
	Budget        *float64                  `json:"budget"`
	Iri           *float64                  `json:"iri"`
	Ifi           *float64                  `json:"ifi"`
	Plans         []MaintenanceAnalysisPlan `json:"plans"`
	PrepareDataID []int                     `json:"prepare_data_id"`
}

type ChkPrepareDataReq struct {
	IDs []int `json:"id"`
}
type PrepareDataIDReq struct {
	PrepareDataID []int `json:"prepare_data_id"`
}

type AnalysisFilter struct {
	TypeAnalysis *int    `json:"type_analysis"`
	Condition    *string `json:"condition"`
}

type MapFilter struct {
	Year     *int `json:"year"`
	Plan     *int `json:"plan"`
	Display  int  `json:"display"`
	Criteria *int `json:"criteria"`
}
