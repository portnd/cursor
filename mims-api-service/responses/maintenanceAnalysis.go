package responses

import (
	"time"

	"gitlab.com/mims-api-service/models"
)

type ResponseMaintenanceAnalysis struct {
	MaintenanceAnalysis models.MaintenanceAnalysis `json:"maintenance_analysis"`
	PrepareData         interface{}                `json:"prepare_data"`
}

type AnalysisByIDRes struct {
	ID                        int                               `json:"id"`
	Name                      string                            `json:"name"`
	MaintenanceAnalysisTypeId int                               `json:"maintenance_analysis_type_id"`
	Roads                     []models.MaintenanceAnalysisRoads `json:"roads"`
	SurfaceTypeId             int                               `json:"surface_type_id"`
	LaneTypeId                int                               `json:"lane_type_id"`
	Iri1                      *float64                          `json:"iri1"`
	Iri2                      *float64                          `json:"iri2"`
	Aadt1                     *float64                          `json:"aadt1"`
	Aadt2                     *float64                          `json:"aadt2"`
	Ifi1                      *float64                          `json:"ifi1"`
	Ifi2                      *float64                          `json:"ifi2"`
	GroupKm                   float64                           `json:"group_km"`
	Condition                 int                               `json:"condition"`
	NumberPlan                int                               `json:"number_plan"`
	CreatedBy                 int                               `json:"created_by"`
	UpdatedBy                 int                               `json:"updated_by"`
	CreatedAt                 time.Time                         `json:"created_at"`
	UpdatedAt                 time.Time                         `json:"updated_at"`
}

type PrepareDataRes struct {
	PrepareData interface{} `json:"prepare_data"`
}

type AnalysisRes struct {
	ID                         int       `json:"id"`
	Name                       string    `json:"name"`
	TypeAnalysis               string    `json:"type_analysis"`
	MaintenanceAnalysisTypeId  int       `json:"maintenance_analysis_type_id"`
	MaintenanceConditionTypeId int       `json:"maintenance_condition_type_id"`
	Comment                    string    `json:"comment"`
	AnalysisDate               time.Time `json:"analysis_date"`
	Percentage                 int       `json:"percentage"`
	Status                     string    `json:"status"`
	IsFavorite                 bool      `json:"is_favorite"`
}

type AnalysisStep2Res struct {
	ConditionID *int        `json:"condition_id"`
	SurfaceType string      `json:"surface_type"`
	TotalKm     float64     `json:"total_km"`
	IRIAvg      float64     `json:"iri_avg"`
	IFIAvg      float64     `json:"ifi_avg"`
	Discount    *float64    `json:"discount"`
	Year        *int        `json:"year"`
	Target      *int        `json:"target"`
	NumberPlan  *int        `json:"number_plan"`
	Comment     *string     `json:"comment"`
	Budget      *float64    `json:"budget"`
	IRI         *float64    `json:"iri"`
	IFI         *float64    `json:"ifi"`
	Plans       interface{} `json:"plans"`
}

type Report4 struct {
	RoadID               int     `json:"road_id"`
	Year                 int     `json:"year"`
	RoadCode             string  `json:"road_coad"`
	RoadName             string  `json:"road_name"`
	RoadInfoName         string  `json:"road_info_name"`
	KmStart              int     `json:"km_start"`
	KmEnd                int     `json:"km_end"`
	KmTotal              float64 `json:"km_total"`
	LaneNo               int     `json:"lane_no"`
	InterventionCriteria string  `json:"intervention_criteria"`
	Area                 float64 `json:"area"`
	Budget               float64 `json:"budget"`
	BC                   float64 `json:"b_c"`
	Aadt                 float64 `json:"aadt"`
	IriBefore            float64 `json:"iri_before"`
	IriAfter             float64 `json:"iri_after"`
	Acc                  float64 `json:"acc"`
	Voc                  float64 `json:"voc"`
	Vot                  float64 `json:"vot"`
	Ruc                  float64 `json:"ruc"`
	Benefit              float64 `json:"benifit"`
	AccRm                float64 `json:"acc_rm"`
	VocRm                float64 `json:"voc_rm"`
	VotRm                float64 `json:"vot_rm"`
	RucRm                float64 `json:"ruc_rm"`
}

type Report4Res struct {
	RoadID               int     `json:"road_id"`
	Year                 int     `json:"year"`
	RoadCode             string  `json:"road_coad"`
	RoadName             string  `json:"road_name"`
	RoadInfoName         string  `json:"road_info_name"`
	KmStart              int     `json:"km_start"`
	KmEnd                int     `json:"km_end"`
	KmTotal              float64 `json:"km_total"`
	LaneNo               int     `json:"lane_no"`
	InterventionCriteria string  `json:"intervention_criteria"`
	Area                 string  `json:"area"`
	Budget               string  `json:"budget"`
	BC                   string  `json:"b_c"`
	Aadt                 string  `json:"aadt"`
	IriBefore            string  `json:"iri_before"`
	IriAfter             string  `json:"iri_after"`
	Acc                  string  `json:"acc"`
	Voc                  string  `json:"voc"`
	Vot                  string  `json:"vot"`
	Ruc                  string  `json:"ruc"`
	Benefit              string  `json:"benifit"`
	AccRm                string  `json:"acc_rm"`
	VocRm                string  `json:"voc_rm"`
	VotRm                string  `json:"vot_rm"`
	RucRm                string  `json:"ruc_rm"`
}

type Report4Data struct {
	Report4   []Report4Res
	Title     string `json:"title"`
	Date      string `json:"date"`
	User      string `json:"user"`
	Condition string `json:"condition"`
	Target    string `json:"target"`
	PlanName  string `json:"plan_name"`
}

type Report3 struct {
	Seq   float64     `json:"seq"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type Report3Data struct {
	Report3   []Report3
	Year      interface{}
	Title     string `json:"title"`
	Date      string `json:"date"`
	User      string `json:"user"`
	Condition string `json:"condition"`
	Target    string `json:"target"`
}

type Report3Year struct {
	Year interface{}
}

type CheckPrepareDataStatus struct {
	Status bool `json:"status"`
}

type PrepareDataWithPagination struct {
	ID                    int     `json:"id"`
	MaintenanceAnalysisID int     `json:"maintenance_analysis_id"`
	IsSelected            bool    `json:"is_selected"`
	GroupName             string  `json:"group_name"`
	RoadName              string  `json:"road_name"`
	LaneNo                int     `json:"lane_no"`
	KmStart               float64 `json:"km_start"`
	KmEnd                 float64 `json:"km_end"`
	Length                float64 `json:"length"`
	Iri                   float64 `json:"iri"`
	Ifi                   float64 `json:"ifi"`
	AADT                  float64 `json:"aadt"`
}

type AnalysisIsFavorite struct {
	IsFavorite bool `json:"is_favorite"`
}
