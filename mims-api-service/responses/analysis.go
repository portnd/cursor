package responses

import (
	"time"

	"gitlab.com/mims-api-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KMReport struct {
	KM                       int    `json:"km"`
	KmStart                  int    `json:"km_start"`
	KmEnd                    int    `json:"km_end"`
	LaneNo                   int    `json:"lane_no"`
	Year                     int    `json:"year"`
	InterventionCriteriaName string `json:"intervention_criteria_name"`
	Method                   string `json:"method"`
}
type KM struct {
	KmStart int `json:"km_start"`
	KmEnd   int `json:"km_end"`
}

type KMData struct {
	KmStart                  int    `json:"km_start"`
	KmEnd                    int    `json:"km_end"`
	InterventionCriteriaName string `json:"intervention_criteria_name"`
	Method                   string `json:"method"`
	Unique                   string `json:"unique"`
	Color                    string `json:"color"`
}

type KMStr struct {
	KmStart string      `json:"km_start"`
	KmEnd   string      `json:"km_end"`
	Data    interface{} `json:"data"`
}

type KMFolat64 struct {
	KmStart int `json:"km_start"`
	KmEnd   int `json:"km_end"`
}

type mainternanceKM struct {
	KM
	Maintenances []models.MaintenanceData `json:"maintenances"`
	Surfaces     []models.RoadSurface     `json:"surfaces"`
}

type PrepareData struct {
	KM
	RoadGroupID         int                        `json:"road_group_id"`
	RoadGroupName       string                     `json:"road_group_name"`
	RoadID              int                        `json:"road_id"`
	RoadName            string                     `json:"road_name"`
	RoadLength          float64                    `json:"road_length"`
	LaneNo              int                        `json:"lane_no"`
	RefStructureSurface models.RefStructureSurface `json:"ref_structure_surface"`
	RefSurface          models.RefSurface          `json:"ref_surface"`
	RefSurfaceParam     models.RefSurface          `json:"ref_surface_param"`

	RefMaterialBase     models.RefMaterialBase     `json:"ref_material_base"`
	RefMaterialSubbase  models.RefMaterialSubbase  `json:"ref_material_subbase"`
	RefMaterialSubgrade models.RefMaterialSubgrade `json:"ref_material_subgrade"`

	WidthSurface     float64 `json:"width_surface"`
	ThicknessSurface float64 `json:"thickness_surface"`

	ThicknessBase     *float64 `json:"thickness_base"`
	ThicknessSubbase  *float64 `json:"thickness_subbase"`
	ThicknessSubgrade *float64 `json:"thickness_subgrade"`

	Year int `json:"year"`

	YearLastOverlay        int                      `json:"year_last_overlay"`
	YearLastSeal           int                      `json:"year_last_seal"`
	YearLastMol            int                      `json:"year_last_mol"`
	YearLastRcl            int                      `json:"year_last_rcl"`
	YearLastReconstruction int                      `json:"year_last_reconstruction"`
	AcHsold                float64                  `json:"ac_hsold"`
	AcHsnew                float64                  `json:"ac_hsnew"`
	Maintenances           []models.MaintenanceData `json:"-"`
}

type Data struct {
	Maintenances []mainternanceKM `json:"maintenances"`
}

type Road struct {
	RoadID                    int    `json:"road_id"`
	RoadName                  string `json:"road_name"`
	RoadGroupID               int    `json:"road_group_id"`
	RoadGroupName             string `json:"road_group_name"`
	RefDirectionID            int    `json:"ref_direction_id"`
	YearConstructionCompleted int    `json:"year_construction_completed"`
}

type RoadDateBegins struct {
	RoadID  int     `json:"road_id"`
	KmStart float64 `json:"km_start"`
	KmEnd   float64 `json:"km_end"`
	Year    int     `json:"year"`
}

type RoadGeoms struct {
	RoadID  int     `json:"road_id"`
	KmStart float64 `json:"km_start"`
	KmEnd   float64 `json:"km_end"`
	LaneNo  int     `json:"lane_no"`
}

type RefSurface struct {
	ID                        int                       `json:"id"`
	Name                      string                    `json:"name"`
	Type                      string                    `json:"type"`
	SurfaceGroup              string                    `json:"surface_group"`
	LayerCoefficient          float64                   `json:"layer_coefficient"`
	Drainage                  float64                   `json:"drainage"`
	A                         float64                   `json:"a"`
	B                         float64                   `json:"b"`
	Cbase                     float64                   `json:"c_base"`
	Cexp                      float64                   `json:"c_exp"`
	Crt                       float64                   `json:"crt"`
	Rrf                       float64                   `json:"rrf"`
	Raveling                  Raveling                  `json:"raveling"`
	AllStructuralCrack        AllStructuralCrack        `json:"all_structural_crack"`
	WideStructuralCrack       WideStructuralCrack       `json:"wide_structural_crack"`
	RuttingPlasticDeformation RuttingPlasticDeformation `json:"rutting_plastic_deformation"`
	ThicknessRepair           float64                   `json:"thickness_repair"`
	ThicknessScrape           float64                   `json:"thickness_scrape"`
}

type Raveling struct {
	Initial struct {
		A0 float64 `json:"a0"`
		A1 float64 `json:"a1"`
	} `json:"initial"`
	Progression struct {
		A0 float64 `json:"a0"`
		A1 float64 `json:"a1"`
		A2 float64 `json:"a2"`
	} `json:"progression"`
}

// ////////////////////
type AllStructuralCrack struct {
	Initial struct {
		HSOLD_O struct {
			A0 float64 `json:"a0"`
			A1 float64 `json:"a1"`
			A2 float64 `json:"a2"`
			A3 float64 `json:"a3"`
			A4 float64 `json:"a4"`
		} `json:"hsold_o"`
		HSOLD struct {
			A0 float64 `json:"a0"`
			A1 float64 `json:"a1"`
			A2 float64 `json:"a2"`
			A3 float64 `json:"a3"`
			A4 float64 `json:"a4"`
		} `json:"hsold"`
	} `json:"initial"`
	Progression struct {
		HSOLD_O struct {
			A0 float64 `json:"a0"`
			A1 float64 `json:"a1"`
		} `json:"hsold_o"`
		HSOLD struct {
			A0 float64 `json:"a0"`
			A1 float64 `json:"a1"`
		} `json:"hsold"`
	} `json:"progression"`
}

// ////////////////////
type WideStructuralCrack struct {
	Initial struct {
		HSOLD_O struct {
			A0 float64 `json:"a0"`
			A1 float64 `json:"a1"`
			A2 float64 `json:"a2"`
		} `json:"hsold_o"`
		HSOLD struct {
			A0 float64 `json:"a0"`
			A1 float64 `json:"a1"`
			A2 float64 `json:"a2"`
		} `json:"hsold"`
	} `json:"initial"`
	Progression struct {
		HSOLD_O struct {
			A0 float64 `json:"a0"`
			A1 float64 `json:"a1"`
		} `json:"hsold_o"`
		HSOLD struct {
			A0 float64 `json:"a0"`
			A1 float64 `json:"a1"`
		} `json:"hsold"`
	} `json:"progression"`
}

// ////////////////////////////
type RuttingPlasticDeformation struct {
	A0 float64 `json:"a0"`
	A1 float64 `json:"a1"`
	A2 float64 `json:"a2"`
}

type IcResult struct {
	Type            string     `json:"type"`
	Method          string     `json:"method"`
	ThicknessRepair float64    `json:"thickness_repair"`
	ThicknessScrape float64    `json:"thickness_scrape"`
	RefSurface      RefSurface `json:"ref_surface"`

	// AllStructuralCrack        AllStructuralCrack        `json:"all_structural_crack"`
	// WideStructuralCrack       WideStructuralCrack       `json:"wide_structural_crack"`
	// RuttingPlasticDeformation RuttingPlasticDeformation `json:"rutting_plastic_deformation"`
	// ... other fields ...
}

type MaintenanceData struct {
	RoadGroupID    int       `json:"road_group_id"`
	RoadID         int       `json:"road_id"`
	KmStart        float64   `json:"km_start"`
	KmEnd          float64   `json:"km_end"`
	LaneNo         int       `json:"lane_no"`
	Year           int       `json:"year"`
	BudgetType     float64   `json:"budget_type"`
	IcResult       IcResult  `json:"ic_result"`
	ProjectEndDate time.Time `json:"project_end_date"`
	// LastInspectionDate time.Time `json:"last_inspection_date"`
	ContractNumber string `json:"contract_number"`
	BudgetYear     int    `json:"budget_year"`
}
type RoadConditionPrepareDataRes struct {
	RUT        float64     `json:"rut"`
	IRI        float64     `json:"iri"`
	IFI        float64     `json:"ifi"`
	SurveyDate interface{} `json:"survey_date"`
}

type RoadDamagePrepareDataRes struct {
	NumberOfPothole    float64     `json:"number_of_pothole"`
	AreaAcIcrack       float64     `json:"area_ac_icrack"`
	PercentAcIcrack    float64     `json:"percent_ac_icrack"`
	AreaAcUcrack       float64     `json:"area_ac_ucrack"`
	PercentAcUcrack    float64     `json:"percent_ac_ucrack"`
	PercentAcRavelling float64     `json:"percent_ac_ravelling"`
	CcTransverseCrack  float64     `json:"cc_transverse_crack"`
	CcFaulting         float64     `json:"cc_faulting"`
	CcSpalling         float64     `json:"cc_spalling"`
	SurveyDate         interface{} `json:"survey_date"`
}

type VolumeAadtPrepareDataRes struct {
	Veh1 int `json:"veh1"`
	Veh2 int `json:"veh2"`
	Veh3 int `json:"veh3"`
}

type VolumeRainPrepareDataRes struct {
	RoadGroupID int     `json:"road_group_id"`
	MinRain     float64 `json:"min_rain"`
	MaxRain     float64 `json:"max_rain"`
	AvgRain     float64 `json:"avg_rain"`
}

type AgePrepareDataRes struct {
	YearLastOverlay        int `json:"year_last_overlay"`
	YearLastSeal           int `json:"year_last_seal"`
	YearLastMolRcl         int `json:"year_last_mol_rcl"`
	YearLastReconstruction int `json:"year_last_reconstruction"`
	YearLastFdr            int `json:"year_last_fdr"`
	YearLastOvl            int `json:"year_last_ovl"`
	YearLastMol            int `json:"year_last_mol"`
	Age                    int `json:"age"`
}

type HsoldHsnewPrepareDataRes struct {
	Hsold float64 `json:"hsold"`
	Hsnew float64 `json:"hsnew"`
}

type SnpPrepareDataRes struct {
	Method           string  `json:"-"`
	Thickness        float64 `json:"-"`
	LayerCoefficient float64 `json:"-"`
	Drainage         float64 `json:"-"`
	SnpSurface       float64 `json:"snp_surface"`
	SnpBase          float64 `json:"snp_base"`
	SnpSubbase       float64 `json:"snp_subbase"`
	Snp              float64 `json:"snp"`
}

type ReportListRes struct {
	ReportName string `json:"report_name"`
	Url        string `json:"url"`
}

type Report5 struct {
	KM      int    `json:"km"`
	KMStart string `json:"km_start"`
	KMEnd   string `json:"km_end"`
	LaneNo  int    `json:"lane_no"`
	Span    int    `json:"span"`
	Desc    int    `json:"desc"`
	IsWrite int    `json:"is_write"`
	Year    int    `json:"year"`
	// Years   []int  `json:"years"`
	Data KMData `json:"data"`
}

type Report5Res struct {
	Plan       int           `json:"plan"`
	Condition  string        `json:"condition"`
	Target     string        `json:"target"`
	YearLength string        `json:"year_length"`
	CreatedAt  string        `json:"CreatedAt"`
	User       string        `json:"user"`
	RoadCode   string        `json:"road_code"`
	ReportRoad []Report5Road `json:"report_road"`
}

type Report5Road struct {
	// Plan       int         `json:"plan"`
	// Condition  string      `json:"condition"`
	// Target     string      `json:"target"`
	// YearLength string      `json:"year_length"`
	RoadSectionCode string      `json:"road_section_code"`
	RoadGroupCode   string      `json:"road_group_code"`
	RoadID          string      `json:"road_id"`
	RoadName        string      `json:"road_name"`
	Rows            interface{} `json:"rows"`
	Years           []int       `json:"years"`
	Data            interface{} `json:"data"` // []Report5ResYear `json:"yearData"`
}

type Report5ResYear struct {
	Year int       `json:"year"`
	Data []Report5 `json:"data"`
}

type Res struct {
	Rows interface{} `json:"rows"`
	Data interface{} `json:"data"`
}

type LaneData struct {
	LaneNo interface{} `json:"lane_no"`
	Data   Report5     `json:"data"`
}

type YearData struct {
	Year int         `json:"year"`
	Data interface{} `json:"data"`
}

type DataData struct {
	Data interface{} `json:"data"`
}

type BsonM struct {
	ID map[string]interface{} `mapstructure:"_id"`
}

type AnalysesReportMaintenance struct {
	models.ModelResult
	RoadID    int    `json:"road_id"`
	RoadCode  string `json:"road_code"`
	RoadName  string `gorm:"column:name" json:"road_name"`
	KmStart   int
	KmEnd     int
	Rows      []interface{} `json:"rows"`
	Years     []int         `json:"years"`
	YearStart int           `json:"year_start"`
	YearEnd   int           `json:"year_end"`
}

type AnalysesModel struct {
	ID                   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Year                 int                `json:"year"`
	Plan                 string             `json:"plan"`
	RoadGroupName        string             `gorm:"column:name" json:"road_group_name"`
	RoadName             string             `gorm:"column:name" json:"road_name"`
	KmStart              float64            `json:"km_start"`
	KmEnd                float64            `json:"km_end"`
	Distance             float64            `json:"distance"`
	Lane                 int                `json:"lane"`
	InterventionCriteria interface{}        `json:"intervention_criteria"`
	Area                 float64            `json:"area"`
	Cost                 float64            `json:"cost"`
	BC                   float64            `json:"bc"`
	VolumeAadt           int                `json:"volume_aadt"`
	IriBefore            float64            `json:"iri_before"`
	IriAfter             float64            `json:"iri_after"`
}

type Plan struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Displays struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Criterias struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Grade   []Grade `json:"grade"`
	GradeCC []Grade `json:"grade_cc"`
}

type Methods struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Filter struct {
	Year     []int       `json:"year"`
	Display  []Displays  `json:"display"`
	Criteria []Criterias `json:"criteria"`
	Plan     []Plan      `json:"plan"`
	Method   []Methods   `json:"method"`
}

type Grade struct {
	Name           string  `json:"name"`
	Color          string  `json:"color"`
	LeftValue      float64 `json:"left_value"`
	LeftCondition  string  `json:"left_condition"`
	RightValue     float64 `json:"right_value"`
	RightCondition string  `json:"right_condition"`
	ConditionType  string  `json:"condition_type"`
}

type DashboardMap struct {
	CriteriaMethod []DashboardMapCriteriaMethod `json:"criteria_method"`
	Items          []DashboardMapData           `json:"items"`
}

type DashboardMapCriteriaMethod struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type DashboardMapData struct {
	Title    string `json:"title"`
	Display  int    `json:"display"`
	RoadName string `json:"road_name"`
	KmStart  int    `json:"km_start"`
	KmEnd    int    `json:"km_end"`
	TheGeom  struct {
		Type        string      `json:"type"`
		Coordinates [][]float64 `json:"coordinates"`
	} `json:"the_geom"`
	IriBefore float64 `json:"iri_before"`
	IriAfter  float64 `json:"iri_after"`
	Year      *int    `json:"year"`
	Color     string  `json:"color"`
}

type RefDepotByRoad struct {
	RefDepotCode string `json:"ref_depot_code"`
	RoadID       int    `json:"road_id"`
}
