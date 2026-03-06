package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MaintenanceAnalysis struct {
	ID                        int      `json:"id"`
	Name                      string   `json:"name"`
	MaintenanceAnalysisTypeId int      `json:"maintenance_analysis_type_id"`
	SurfaceTypeId             int      `json:"surface_type_id"`
	LaneTypeId                int      `json:"lane_type_id"`
	RoadGroupId               int      `json:"road_group_id"`
	Iri1                      *float64 `json:"iri1"`
	Iri2                      *float64 `json:"iri2"`
	Aadt1                     *float64 `json:"aadt1"`
	Aadt2                     *float64 `json:"aadt2"`
	Ifi1                      *float64 `json:"ifi1"`
	Ifi2                      *float64 `json:"ifi2"`
	GroupKm                   float64  `json:"group_km"`
	Percentage                float64  `json:"percentage"`
	Status                    string   `json:"status"`
	IsFavorite                bool     `json:"is_favorite"`
	Condition                 int      `json:"condition"`
	Discount                  *float64 `json:"discount"`
	Year                      *int     `json:"year"`
	Target                    *int     `json:"target"`
	NumberPlan                *int     `json:"number_plan"`
	Comment                   string   `json:"comment"`
	Budget                    *float64 `json:"budget"`
	Iri                       *float64 `json:"iri"`
	Ifi                       *float64 `json:"ifi"`

	InterventionCriteriaParmasID int `json:"intervention_criteria_parmas_id"`
	RoadWorkEffectParmasID       int `json:"road_work_effect_parmas_id"`
	RoadUserCostParmasID         int `json:"road_user_cost_parmas_id"`
	DeterationParmasID           int `json:"deteration_parmas_id"`
	OptimizationParmasID         int `json:"optimization_parmas_id"`
	AadtParmasID                 int `json:"aadt_parmas_id"`

	FilterData string `json:"filter_data"`
	PreviousID int    `json:"previous_id"`

	IsLatest  bool      `json:"is_latest"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	PrepareDataStatus bool `json:"prepare_data_status"`
}

// type MaintenanceAnalysisRoad struct {
// 	Id                    int       `json:"id"`
// 	MaintenanceAnalysisId int       `json:"maintenance_analysis_id"`
// 	RoadGroupId           int       `json:"road_group_id"`
// 	RoadId                int       `json:"road_id"`
// 	KmStart               float64   `json:"km_start"`
// 	KmEnd                 float64   `json:"km_end"`
// 	LaneTypeId            int       `json:"lane_type_id"`
// 	Iri                   float64   `json:"iri"`
// 	Aadt                  float64   `json:"aadt"`
// 	Gn                    float64   `json:"gn"`
// 	IsSelected            bool      `json:"is_selected"`
// 	IsLatest              bool      `json:"is_latest"`
// 	IsDeleted             bool      `json:"is_deleted"`
// 	CreatedBy             int       `json:"created_by"`
// 	UpdatedBy             int       `json:"updated_by"`
// 	CreatedAt             time.Time `json:"created_at"`
// 	UpdatedAt             time.Time `json:"updated_at"`
// }

type MaintenanceAnalysisPreload struct {
	MaintenanceAnalysis
	PrepareData   []PrepareData                   `json:"prepare_data" gorm:"ForeignKey:MaintenanceAnalysisID;AssociationForeignKey:ID"`
	TargetData    RefMaintenanceAnalysisTarget    `json:"target_data" gorm:"ForeignKey:Target;AssociationForeignKey:ID"`
	ConditionData RefMaintenanceAnalysisCondition `json:"condition_data" gorm:"ForeignKey:Condition;AssociationForeignKey:ID"`
	// Strategic   MaintenanceAnalysisStrategicPreload `json:"strategic" gorm:"ForeignKey:Id;AssociationForeignKey:MaintenanceAnalysisId"`
}
type MaintenanceAnalysisRoadPreload struct {
	// MaintenanceAnalysisRoad
	PrepareData
	RoadGroup RoadGroup `json:"road_group" gorm:"ForeignKey:RoadGroupId;AssociationForeignKey:Id"`
	Road      Road      `json:"road" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
}

type MaintenanceAnalysisStrategicPreload struct {
	MaintenanceAnalysisStrategic
	PlanList []MaintenanceAnalysisPlan `json:"plan_list" gorm:"ForeignKey:MaintenanceAnalysisStrategicId;AssociationForeignKey:Id"`
}

type MaintenanceAnalysisStrategic struct {
	Id                                 int       `json:"id"`
	MaintenanceAnalysisId              int       `json:"maintenance_analysis_id"`
	RefMaintenanceAnalysisBudgetTypeId int       `json:"ref_maintenance_analysis_budget_type_id"`
	RefMaintenanceAnalysisTargetTypeId int       `json:"ref_maintenance_analysis_target_type_id"`
	Discount                           float64   `json:"discount"`
	RangeYear                          int       `json:"range_year"`
	NumPlan                            int       `json:"num_plan"`
	Comment                            string    `json:"comment"`
	Budget                             float64   `json:"budget"`
	Iri                                float64   `json:"iri"`
	Gn                                 float64   `json:"gn"`
	AnalysisStrategicDate              time.Time `json:"analysis_strategic_date"`
	IsLatest                           bool      `json:"is_latest"`
	IsDeleted                          bool      `json:"is_deleted"`
	CreatedBy                          int       `json:"created_by"`
	UpdatedBy                          int       `json:"updated_by"`
	CreatedAt                          time.Time `json:"created_at"`
	UpdatedAt                          time.Time `json:"updated_at"`
}

type ModelResult struct {
	ID                             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	MaintenanceAnalysisId          int                `json:"maintenance_analysis_id" bson:"maintenance_analysis_id"`
	MaintenanceAnalysisTypeID      int                `json:"maintenance_analysis_type_id" bson:"maintenance_analysis_type_id"`
	MaintenanceAnalysisConditionID int                `json:"maintenance_analysis_condition_id" bson:"maintenance_analysis_condition_id"`
	RoadID                         int                `json:"road_id" bson:"road_id"`
	Plan                           int                `json:"plan" bson:"plan"`
	PlanSequence                   int                `json:"plan_sequence" bson:"plan_sequence"`
	RepairBudgetType               string             `json:"repair_budget_type" bson:"repair_budget_type"`
	Relation                       string             `json:"relation" bson:"relation"`
	RoadNumber                     int                `json:"road_number" bson:"road_number"`
	Year                           int                `json:"year" bson:"year"`
	KmStart                        int                `json:"km_start" bson:"km_start"`
	AnalystYear                    int                `json:"analyst_year" bson:"analyst_year"`
	KmEnd                          int                `json:"km_end" bson:"km_end"`
	Data                           struct {
		BcAfter                      float64 `json:"bc_after" bson:"bc_after"`
		IriAfter                     float64 `json:"iri_after" bson:"iri_after"`
		IfiAfter                     float64 `json:"ifi_after" bson:"ifi_after"`
		InterventionCriteriaID       int     `json:"intervention_criteria_id" bson:"intervention_criteria_id"`
		InterventionCriteriaChangeID int     `json:"intervention_criteria_change_id" bson:"intervention_criteria_change_id"`
		BC                           bool    `json:"BC" bson:"BC"`
		Budget                       bool    `json:"budget" bson:"budget"`
		Repair                       bool    `json:"repair" bson:"repair"`
		Ic                           bool    `json:"ic" bson:"ic"`
		DeteriorationResult          struct {
			Result struct {
				Bc  float64 `json:"bc" bson:"bc"`
				Iri float64 `json:"iri" bson:"iri"`
				Ifi float64 `json:"ifi" bson:"ifi"`
			} `json:"result" bson:"result"`
		} `json:"deterioration_result" bson:"deterioration_result"`
		IcResult struct {
			IcParamsID      int     `json:"ic_params_id" bson:"ic_params_id"`
			IcId            int     `json:"ic_id" bson:"ic_id"`
			Name            string  `json:"name" bson:"name"`
			Description     string  `json:"description" bson:"description"`
			Type            string  `json:"type" bson:"type"`
			Method          string  `json:"method" bson:"method"`
			ThicknessScrape float64 `json:"thickness_scrape" bson:"thickness_scrape"`
			ThicknessRepair float64 `json:"thickness_repair" bson:"thickness_repair"`
			PricePerUnit    int     `json:"price_per_unit" bson:"price_per_unit"`
			RefSurface      struct {
				ID               int     `json:"id" bson:"id"`
				Name             string  `json:"name" bson:"name"`
				Type             string  `json:"type" bson:"type"`
				SurfaceGroup     string  `json:"surface_group" bson:"surface_group"`
				LayerCoefficient float64 `json:"layer_coefficient" bson:"layer_coefficient"`
				Drainage         float64 `json:"drainage" bson:"drainage"`
				A                float64 `json:"a" bson:"a"`
				B                float64 `json:"b" bson:"b"`
				CBase            float64 `json:"c_base" bson:"c_base"`
				CExp             float64 `json:"c_exp" bson:"c_exp"`
				Crt              float64 `json:"crt" bson:"crt"`
				Rrf              float64 `json:"rrf" bson:"rrf"`
				Raveling         struct {
					Initial struct {
						A0 float64 `json:"a0" bson:"a0"`
						A1 float64 `json:"a1" bson:"a1"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						A0 float64 `json:"a0" bson:"a1"`
						A1 float64 `json:"a1" bson:"a2"`
						A2 float64 `json:"a2" bson:"a3"`
					} `json:"progression" bson:"progression"`
				} `json:"raveling" bson:"raveling"`
				AllStructuralCrack struct {
					Initial struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
							A3 float64 `json:"a3" bson:"a3"`
							A4 float64 `json:"a4" bson:"a4"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
							A3 float64 `json:"a3" bson:"a3"`
							A4 float64 `json:"a4" bson:"a4"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"progression" bson:"progression"`
				} `json:"all_structural_crack" bson:"all_structural_crack"`
				WideStructuralCrack struct {
					Initial struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"10"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"progression" bson:"progression"`
				} `json:"wide_structural_crack" bson:"wide_structural_crack"`
				RuttingPlasticDeformation struct {
					A0 float64 `json:"a0" bson:"a0"`
					A1 float64 `json:"a1" bson:"a1"`
					A2 float64 `json:"a2" bson:"a2"`
				} `json:"rutting_plastic_deformation" bson:"rutting_plastic_deformation"`
			} `json:"ref_surface" bson:"ref_surface"`
		} `json:"ic_result" bson:"ic_result"`
		OptimizationResult struct {
			Benefit float64 `json:"benefit" bson:"benefit"`
			DIRI    float64 `json:"dIRI" bson:"dIRI"`
			DGN     float64 `json:"dGN" bson:"dGN"`
			Cost    float64 `json:"cost" bson:"cost"`
			BC      float64 `json:"bc" bson:"bc"`
			FV      float64 `json:"FV" bson:"FV"`
			NPV     float64 `json:"NPV" bson:"NPV"`
		} `json:"optimization_result" bson:"optimization_result"`
		RucResult struct {
			RucBeforeResult struct {
				Summary struct {
					RoadUserCostParamsID int     `json:"road_user_cost_params_id" bson:"road_user_cost_params_id"`
					AadtParamsID         int     `json:"aadt_params_id" bson:"aadt_params_id"`
					TotalPcu             float64 `json:"total_pcu" bson:"total_pcu"`
					Vkt                  float64 `json:"vkt" bson:"vkt"`
					ACCCostR             float64 `json:"ACC_Cost_r" bson:"ACC_Cost_r"`
					CMF                  float64 `json:"CMF" bson:"CMF"`
					ACC                  float64 `json:"ACC" bson:"ACC"`
					ACCFINAL             float64 `json:"ACC_FINAL" bson:"ACC_FINAL"`
					VOCFINAL             float64 `json:"VOC_FINAL" bson:"VOC_FINAL"`
					VOTFINAL             float64 `json:"VOT_FINAL" bson:"VOT_FINAL"`
					TOTALRUCNOACC        float64 `json:"TOTAL_RUC_NO_ACC" bson:"TOTAL_RUC_NO_ACC"`
					TOTALRUCWITHACC      float64 `json:"TOTAL_RUC_WITH_ACC" bson:"TOTAL_RUC_WITH_ACC"`
				} `json:"summary" bson:"summary"`
			} `json:"ruc_before_result" bson:"ruc_before_result"`
			RucAfterResult struct {
				Summary struct {
					RoadUserCostParamsID int     `json:"road_user_cost_params_id" bson:"road_user_cost_params_id"`
					AadtParamsID         int     `json:"aadt_params_id" bson:"aadt_params_id"`
					TotalPcu             float64 `json:"total_pcu" bson:"total_pcu"`
					Vkt                  float64 `json:"vkt" bson:"vkt"`
					ACCCostR             float64 `json:"ACC_Cost_r" bson:"ACC_Cost_r"`
					CMF                  float64 `json:"CMF" bson:"CMF"`
					ACC                  float64 `json:"ACC" bson:"ACC"`
					ACCFINAL             float64 `json:"ACC_FINAL" bson:"ACC_FINAL"`
					VOCFINAL             float64 `json:"VOC_FINAL" bson:"VOC_FINAL"`
					VOTFINAL             float64 `json:"VOT_FINAL" bson:"VOT_FINAL"`
					TOTALRUCNOACC        float64 `json:"TOTAL_RUC_NO_ACC" bson:"TOTAL_RUC_NO_ACC"`
					TOTALRUCWITHACC      float64 `json:"TOTAL_RUC_WITH_ACC" bson:"TOTAL_RUC_WITH_ACC"`
				} `json:"summary" bson:"summary"`
			} `json:"ruc_after_result" bson:"ruc_after_result"`
		} `json:"ruc_result" bson:"ruc_result"`
		RweResult struct {
			Age struct {
				YearLastOverlay        int `json:"year_last_overlay" bson:"year_last_overlay"`
				YearLastSeal           int `json:"year_last_seal" bson:"year_last_seal"`
				YearLastMolRcl         int `json:"year_last_mol_rcl" bson:"year_last_mol_rcl"`
				YearLastReconstruction int `json:"year_last_reconstruction" bson:"year_last_reconstruction"`
				Age                    int `json:"age" bson:"age"`
			} `json:"age" bson:"age"`
			AreaAcIcrack   float64 `json:"area_ac_icrack" bson:"area_ac_icrack"`
			AreaAcUcrack   float64 `json:"area_ac_ucrack" bson:"area_ac_ucrack"`
			CurrentSurface struct {
				ID               int     `json:"id" bson:"id"`
				Name             string  `json:"name" bson:"name"`
				Type             string  `json:"type" bson:"type"`
				SurfaceGroup     string  `json:"surface_group" bson:"surface_group"`
				LayerCoefficient float64 `json:"layer_coefficient" bson:"layer_coefficient"`
				Drainage         float64 `json:"drainage" bson:"drainage"`
				A                float64 `json:"a" bson:"a"`
				B                float64 `json:"b" bson:"b"`
				CBase            float64 `json:"c_base" bson:"c_base"`
				CExp             float64 `json:"c_exp" bson:"c_exp"`
				Raveling         struct {
					Initial struct {
						A0 float64 `json:"a0" bson:"a0"`
						A1 float64 `json:"a1" bson:"a1"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						A0 float64 `json:"a0" bson:"a0"`
						A1 float64 `json:"a1" bson:"a1"`
						A2 float64 `json:"a2" bson:"a2"`
					} `json:"progression" bson:"progression"`
				} `json:"raveling" bson:"raveling"`
				AllStructuralCrack struct {
					Initial struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
							A3 float64 `json:"a3" bson:"a3"`
							A4 float64 `json:"a4" bson:"a4"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
							A3 float64 `json:"a3" bson:"a3"`
							A4 float64 `json:"a4" bson:"a4"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"progression" bson:"progression"`
				} `json:"all_structural_crack" bson:"all_structural_crack"`
				WideStructuralCrack struct {
					Initial struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"progression" bson:"progression"`
				} `json:"wide_structural_crack" bson:"wide_structural_crack"`
				RuttingPlasticDeformation struct {
					A0 float64 `json:"a0" bson:"a0"`
					A1 float64 `json:"a1" bson:"a1"`
					A2 float64 `json:"a2" bson:"a2"`
				} `json:"rutting_plastic_deformation" bson:"rutting_plastic_deformation"`
			} `json:"current_surface" bson:"current_surface"`
			Gn         float64 `json:"gn" bson:"gn"`
			HsoldHsnew struct {
				Hsold float64 `json:"hsold" bson:"hsold"`
				Hsnew float64 `json:"hsnew" bson:"hsnew"`
			} `json:"hsold_hsnew" bson:"hsold_hsnew"`
			Iri                float64 `json:"iri" bson:"iri"`
			NumberOfPothole    float64 `json:"number_of_pothole" bson:"number_of_pothole"`
			PercentAcIcrack    float64 `json:"percent_ac_icrack" bson:"percent_ac_icrack"`
			PercentAcRavelling float64 `json:"percent_ac_ravelling" bson:"percent_ac_ravelling"`
			PercentAcUcrack    float64 `json:"percent_ac_ucrack" bson:"percent_ac_ucrack"`
			Rut                float64 `json:"rut" bson:"rut"`
			Snp                struct {
			} `json:"snp" bson:"snp"`
		} `json:"rwe_result" bson:"rwe_result"`
		PrepareData struct {
			Area    float64 `json:"area" bson:"area"`
			Aadt    float64 `json:"aadt" bson:"aadt"`
			KmStart int     `json:"km_start" bson:"km_start"`
			KmEnd   int     `json:"km_end" bson:"km_end"`
			Length  int     `json:"length" bson:"length"`
			Road    struct {
				RoadID         int    `json:"road_id" bson:"road_id"`
				RefDirectionID int    `json:"ref_direction_id" bson:"ref_direction_id"`
				RoadName       string `json:"road_name" bson:"road_name"`
				RoadGroupName  string `json:"road_group_name" bson:"road_group_name"`
			} `json:"road" bson:"road"`
			RoadGeom struct {
				RoadID int `json:"road_id" bson:"road_id"`
				LaneNo int `json:"lane_no" bson:"lane_no"`
			} `json:"road_geom" bson:"road_geom"`
			RoadCondition struct {
				Rut float64 `json:"rut" bson:"rut"`
				Iri float64 `json:"iri" bson:"iri"`
				Gn  float64 `json:"gn" bson:"gn"`
			} `json:"road_condition" bson:"road_condition"`
			RoadInfo struct {
				RoadID int    `json:"road_id" bson:"road_id"`
				Name   string `json:"name" bson:"name"`
			} `json:"road_info" bson:"road_info"`
		} `json:"prepare_data" bson:"prepare_data"`
		PrepareDataBefore struct {
			Area    float64 `json:"area" bson:"area"`
			Aadt    float64 `json:"aadt" bson:"aadt"`
			KmStart int     `json:"km_start" bson:"km_start"`
			KmEnd   int     `json:"km_end" bson:"km_end"`
			Length  int     `json:"length" bson:"length"`
			Road    struct {
				RoadID         int    `json:"road_id" bson:"road_id"`
				RefDirectionID int    `json:"ref_direction_id" bson:"ref_direction_id"`
				RoadName       string `json:"road_name" bson:"road_name"`
				RoadGroupName  string `json:"road_group_name" bson:"road_group_name"`
			} `json:"road" bson:"road"`
			RoadGeom struct {
				RoadID int `json:"road_id" bson:"road_id"`
				LaneNo int `json:"lane_no" bson:"lane_no"`
			} `json:"road_geom" bson:"road_geom"`
			RoadCondition struct {
				Rut float64 `json:"rut" bson:"rut"`
				Iri float64 `json:"iri" bson:"iri"`
				Gn  float64 `json:"gn" bson:"gn"`
			} `json:"road_condition" bson:"road_condition"`
			RoadInfo struct {
				RoadID int    `json:"road_id" bson:"road_id"`
				Name   string `json:"name" bson:"name"`
			} `json:"road_info" bson:"road_info"`
		} `json:"prepare_data_before" bson:"prepare_data_before"`
		UsedBudget   float64 `json:"used_budget" bson:"used_budget"`
		AllBudget    float64 `json:"all_budget" bson:"all_budget"`
		LeftBudget   float64 `json:"left_budget" bson:"left_budget"`
		RepireBudget float64 `json:"repire_budget" bson:"repire_budget"`
		Year         string  `json:"year" bson:"year"`
	} `json:"data" bson:"data"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type DashboardAnnualMaintenance struct {
	Road        []string    `json:"road"`
	Filter      Filter      `json:"filter"`
	Condition   Condition   `json:"condition"`
	Comment     string      `json:"comment"`
	Bar1        AnnualBar   `json:"bar1"`
	Bar2        AnnualBar2  `json:"bar2"`
	AnalystYear int         `json:"analyst_year"`
	Table       AnnualTable `json:"table"`
}

type AnnualBar struct {
	Name  string    `json:"name"`
	Lable []string  `json:"lable"`
	Data  []float64 `json:"data"`
	Area  []float64 `json:"area"`
	Color []string  `json:"color"`
}

type AnnualBar2 struct {
	Name   string    `json:"name"`
	Lable  []string  `json:"lable"`
	Data   []float64 `json:"data"`
	Budget []float64 `json:"budget"`
	Color  []string  `json:"color"`
}

type AnnualTable struct {
	Table1 AnnualTableData1   `json:"table1"`
	Table2 []AnnualTableData2 `json:"table2"`
}

type AnnualTableData1 struct {
	Budget    float64 `json:"budget"`
	IriAfter  float64 `json:"iri_after"`
	IriBefore float64 `json:"iri_before"`
}

type AnnualTableData2 struct {
	Name   string  `json:"name"`
	Area   float64 `json:"aera"`
	Budget float64 `json:"iri_after"`
	Range  float64 `json:"range"`
}

type DashboardStrategicMaintenance struct {
	Road       []string       `json:"road"`
	NumberPlan int            `json:"number_plan"`
	Filter     Filter         `json:"filter"`
	Condition  Condition      `json:"condition"`
	Comment    string         `json:"comment"`
	Graph1     StrategicGraph `json:"graph1"`
	Bar1       StrategicBar   `json:"bar1"`
	Bar2       StrategicBar2  `json:"bar2"`
	Table      StrategicTable `json:"table"`
}

type StrategicTable struct {
	Summary       []StrategicSummaryPlanTable `json:"summary"`
	Plan1         []StrategicDataTable        `json:"plan_1"`
	Plan2         []StrategicDataTable        `json:"plan_2"`
	Plan3         []StrategicDataTable        `json:"plan_3"`
	UnlimitedPlan []StrategicDataTable        `json:"unlimited_plan"`
}

type StrategicDataTable struct {
	MethodName string                    `json:"method_name"`
	Data       []StrategicDataTableOther `json:"data"`
}

type StrategicDataTableOther struct {
	Year   int     `json:"year"`
	Km     float64 `json:"km"`
	Budget float64 `json:"budget"`
}
type StrategicSummaryPlanTable struct {
	Name string                      `json:"name"`
	Data []StrategicSummaryDataTable `json:"data"`
}

type StrategicSummaryDataTable struct {
	Name  string    `json:"name"`
	Value []float64 `json:"value"`
}

type StrategicGraph struct {
	Name  string      `json:"name"`
	Lable []int       `json:"lable"`
	Value [][]float64 `json:"value"`
	Line  []string    `json:"line"`
	Color []string    `json:"color"`
}

type StrategicBar struct {
	Name     string        `json:"name"`
	Lable    []int         `json:"lable"`
	Datasets []DatasetBar1 `json:"datasets"`
	Color    []string      `json:"color"`
}

type DatasetBar1 struct {
	Lable string    `json:"lable"`
	Value []float64 `json:"value"`
}

type Filter struct {
	SurfaceType string   `json:"surface_type"`
	Lane        string   `json:"lane"`
	Km          float64  `json:"km"`
	Filter      []string `json:"filter"`
}

type Condition struct {
	Condition string   `json:"condition"`
	Target    string   `json:"target"`
	Discount  *float64 `json:"discount"`
}

type StrategicBar2 struct {
	Name     string                 `json:"name"`
	Lable    []int                  `json:"lable"`
	Datasets []StrategicDatasetBar2 `json:"datasets"`
}

type StrategicDatasetBar2 struct {
	Plan string              `json:"plan"`
	Data []StrategicDataBar2 `json:"data"`
}

type StrategicDataBar2 struct {
	Lable  []string  `json:"lable"`
	Value  []float64 `json:"value"`
	Budget []float64 `json:"budget"`
	Color  []string  `json:"color"`
}

func (b *MaintenanceAnalysisStrategic) TableName() string {
	return "maintenance_analysis_strategic"
}

func (b *MaintenanceAnalysis) TableName() string {
	return "maintenance_analysis"
}

type ModelResultReport struct {
	MaintenanceAnalysisId int    `json:"maintenance_analysis_id" bson:"maintenance_analysis_id"`
	RoadID                int    `json:"road_id" bson:"road_id"`
	Plan                  int    `json:"plan" bson:"plan"`
	PlanSequence          int    `json:"plan_sequence" bson:"plan_sequence"`
	RepairBudgetType      string `json:"repair_budget_type" bson:"repair_budget_type"`
	Relation              string `json:"relation" bson:"relation"`
	RoadNumber            int    `json:"road_number" bson:"road_number"`
	Year                  int    `json:"year" bson:"year"`
	KmStart               int    `json:"km_start" bson:"km_start"`
	AnalystYear           int    `json:"analyst_year" bson:"analyst_year"`
	KmEnd                 int    `json:"km_end" bson:"km_end"`
	Data                  struct {
		BC       bool `json:"BC" bson:"BC"`
		Budget   bool `json:"budget" bson:"budget"`
		Repair   bool `json:"repair" bson:"repair"`
		Ic       bool `json:"ic" bson:"ic"`
		IcResult struct {
			IcParamsID      int     `json:"ic_params_id" bson:"ic_params_id"`
			IcId            int     `json:"ic_id" bson:"ic_id"`
			Name            string  `json:"name" bson:"name"`
			Description     string  `json:"description" bson:"description"`
			Type            string  `json:"type" bson:"type"`
			Method          string  `json:"method" bson:"method"`
			ThicknessScrape float64 `json:"thickness_scrape" bson:"thickness_scrape"`
			ThicknessRepair float64 `json:"thickness_repair" bson:"thickness_repair"`
			PricePerUnit    int     `json:"price_per_unit" bson:"price_per_unit"`
			RefSurface      struct {
				ID               int     `json:"id" bson:"id"`
				Name             string  `json:"name" bson:"name"`
				Type             string  `json:"type" bson:"type"`
				SurfaceGroup     string  `json:"surface_group" bson:"surface_group"`
				LayerCoefficient float64 `json:"layer_coefficient" bson:"layer_coefficient"`
				Drainage         float64 `json:"drainage" bson:"drainage"`
				A                float64 `json:"a" bson:"a"`
				B                float64 `json:"b" bson:"b"`
				CBase            float64 `json:"c_base" bson:"c_base"`
				CExp             float64 `json:"c_exp" bson:"c_exp"`
				Crt              float64 `json:"crt" bson:"crt"`
				Rrf              float64 `json:"rrf" bson:"rrf"`
				Raveling         struct {
					Initial struct {
						A0 float64 `json:"a0" bson:"a0"`
						A1 float64 `json:"a1" bson:"a1"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						A0 float64 `json:"a0" bson:"a1"`
						A1 float64 `json:"a1" bson:"a2"`
						A2 float64 `json:"a2" bson:"a3"`
					} `json:"progression" bson:"progression"`
				} `json:"raveling" bson:"raveling"`
				AllStructuralCrack struct {
					Initial struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
							A3 float64 `json:"a3" bson:"a3"`
							A4 float64 `json:"a4" bson:"a4"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
							A3 float64 `json:"a3" bson:"a3"`
							A4 float64 `json:"a4" bson:"a4"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"progression" bson:"progression"`
				} `json:"all_structural_crack" bson:"all_structural_crack"`
				WideStructuralCrack struct {
					Initial struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"10"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"progression" bson:"progression"`
				} `json:"wide_structural_crack" bson:"wide_structural_crack"`
				RuttingPlasticDeformation struct {
					A0 float64 `json:"a0" bson:"a0"`
					A1 float64 `json:"a1" bson:"a1"`
					A2 float64 `json:"a2" bson:"a2"`
				} `json:"rutting_plastic_deformation" bson:"rutting_plastic_deformation"`
			} `json:"ref_surface" bson:"ref_surface"`
		} `json:"ic_result" bson:"ic_result"`
		OptimizationResult struct {
			DIRI float64 `json:"dIRI" bson:"dIRI"`
			DGN  float64 `json:"dGN" bson:"dGN"`
			Cost float64 `json:"Cost" bson:"Cost"`
			BC   float64 `json:"BC" bson:"BC"`
			FV   float64 `json:"FV" bson:"FV"`
			NPV  float64 `json:"NPV" bson:"NPV"`
		} `json:"optimization_result" bson:"optimization_result"`
		RucResult struct {
			RucBeforeResult struct {
				Summary struct {
					RoadUserCostParamsID int     `json:"road_user_cost_params_id" bson:"road_user_cost_params_id"`
					AadtParamsID         int     `json:"aadt_params_id" bson:"aadt_params_id"`
					TotalPcu             float64 `json:"total_pcu" bson:"total_pcu"`
					Vkt                  float64 `json:"vkt" bson:"vkt"`
					ACCCostR             float64 `json:"ACC_Cost_r" bson:"ACC_Cost_r"`
					CMF                  float64 `json:"CMF" bson:"CMF"`
					ACC                  float64 `json:"ACC" bson:"ACC"`
					ACCFINAL             float64 `json:"ACC_FINAL" bson:"ACC_FINAL"`
					VOCFINAL             float64 `json:"VOC_FINAL" bson:"VOC_FINAL"`
					VOTFINAL             float64 `json:"VOT_FINAL" bson:"VOT_FINAL"`
					TOTALRUCNOACC        float64 `json:"TOTAL_RUC_NO_ACC" bson:"TOTAL_RUC_NO_ACC"`
					TOTALRUCWITHACC      float64 `json:"TOTAL_RUC_WITH_ACC" bson:"TOTAL_RUC_WITH_ACC"`
				} `json:"summary" bson:"summary"`
			} `json:"ruc_before_result" bson:"ruc_before_result"`
			RucAfterResult struct {
				Summary struct {
					RoadUserCostParamsID int     `json:"road_user_cost_params_id" bson:"road_user_cost_params_id"`
					AadtParamsID         int     `json:"aadt_params_id" bson:"aadt_params_id"`
					TotalPcu             float64 `json:"total_pcu" bson:"total_pcu"`
					Vkt                  float64 `json:"vkt" bson:"vkt"`
					ACCCostR             float64 `json:"ACC_Cost_r" bson:"ACC_Cost_r"`
					CMF                  float64 `json:"CMF" bson:"CMF"`
					ACC                  float64 `json:"ACC" bson:"ACC"`
					ACCFINAL             float64 `json:"ACC_FINAL" bson:"ACC_FINAL"`
					VOCFINAL             float64 `json:"VOC_FINAL" bson:"VOC_FINAL"`
					VOTFINAL             float64 `json:"VOT_FINAL" bson:"VOT_FINAL"`
					TOTALRUCNOACC        float64 `json:"TOTAL_RUC_NO_ACC" bson:"TOTAL_RUC_NO_ACC"`
					TOTALRUCWITHACC      float64 `json:"TOTAL_RUC_WITH_ACC" bson:"TOTAL_RUC_WITH_ACC"`
				} `json:"summary" bson:"summary"`
			} `json:"ruc_after_result" bson:"ruc_after_result"`
		} `json:"ruc_result" bson:"ruc_result"`
		RweResult struct {
			Age struct {
				YearLastOverlay        int `json:"year_last_overlay" bson:"year_last_overlay"`
				YearLastSeal           int `json:"year_last_seal" bson:"year_last_seal"`
				YearLastMolRcl         int `json:"year_last_mol_rcl" bson:"year_last_mol_rcl"`
				YearLastReconstruction int `json:"year_last_reconstruction" bson:"year_last_reconstruction"`
				Age                    int `json:"age" bson:"age"`
			} `json:"age" bson:"age"`
			AreaAcIcrack   float64 `json:"area_ac_icrack" bson:"area_ac_icrack"`
			AreaAcUcrack   float64 `json:"area_ac_ucrack" bson:"area_ac_ucrack"`
			CurrentSurface struct {
				ID               int     `json:"id" bson:"id"`
				Name             string  `json:"name" bson:"name"`
				Type             string  `json:"type" bson:"type"`
				SurfaceGroup     string  `json:"surface_group" bson:"surface_group"`
				LayerCoefficient float64 `json:"layer_coefficient" bson:"layer_coefficient"`
				Drainage         float64 `json:"drainage" bson:"drainage"`
				A                float64 `json:"a" bson:"a"`
				B                float64 `json:"b" bson:"b"`
				CBase            float64 `json:"c_base" bson:"c_base"`
				CExp             float64 `json:"c_exp" bson:"c_exp"`
				Crt              float64 `json:"crt" bson:"crt"`
				Rrf              float64 `json:"rrf" bson:"rrf"`
				Raveling         struct {
					Initial struct {
						A0 float64 `json:"a0" bson:"a0"`
						A1 float64 `json:"a1" bson:"a1"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						A0 float64 `json:"a0" bson:"a0"`
						A1 float64 `json:"a1" bson:"a1"`
						A2 float64 `json:"a2" bson:"a2"`
					} `json:"progression" bson:"progression"`
				} `json:"raveling" bson:"raveling"`
				AllStructuralCrack struct {
					Initial struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
							A3 float64 `json:"a3" bson:"a3"`
							A4 float64 `json:"a4" bson:"a4"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
							A3 float64 `json:"a3" bson:"a3"`
							A4 float64 `json:"a4" bson:"a4"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"progression" bson:"progression"`
				} `json:"all_structural_crack" bson:"all_structural_crack"`
				WideStructuralCrack struct {
					Initial struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"progression" bson:"progression"`
				} `json:"wide_structural_crack" bson:"wide_structural_crack"`
				RuttingPlasticDeformation struct {
					A0 float64 `json:"a0" bson:"a0"`
					A1 float64 `json:"a1" bson:"a1"`
					A2 float64 `json:"a2" bson:"a2"`
				} `json:"rutting_plastic_deformation" bson:"rutting_plastic_deformation"`
			} `json:"current_surface" bson:"current_surface"`
			Gn         float64 `json:"gn" bson:"gn"`
			HsoldHsnew struct {
				Hsold float64 `json:"hsold" bson:"hsold"`
				Hsnew float64 `json:"hsnew" bson:"hsnew"`
			} `json:"hsold_hsnew" bson:"hsold_hsnew"`
			Iri                float64 `json:"iri" bson:"iri"`
			NumberOfPothole    float64 `json:"number_of_pothole" bson:"number_of_pothole"`
			PercentAcIcrack    float64 `json:"percent_ac_icrack" bson:"percent_ac_icrack"`
			PercentAcRavelling float64 `json:"percent_ac_ravelling" bson:"percent_ac_ravelling"`
			PercentAcUcrack    float64 `json:"percent_ac_ucrack" bson:"percent_ac_ucrack"`
			Rut                float64 `json:"rut" bson:"rut"`
			Snp                struct {
			} `json:"snp" bson:"snp"`
		} `json:"rwe_result" bson:"rwe_result"`
		PrepareData struct {
			Area    float64 `json:"area" bson:"area"`
			Aadt    float64 `json:"aadt" bson:"aadt"`
			KmStart int     `json:"km_start" bson:"km_start"`
			KmEnd   int     `json:"km_end" bson:"km_end"`
			Length  int     `json:"length" bson:"length"`
			Road    struct {
				RoadID         int    `json:"road_id" bson:"road_id"`
				RefDirectionID int    `json:"ref_direction_id" bson:"ref_direction_id"`
				RoadName       string `json:"road_name" bson:"road_name"`
				RoadGroupName  string `json:"road_group_name" bson:"road_group_name"`
			} `json:"road" bson:"road"`
			RoadGeom struct {
				RoadID int `json:"road_id" bson:"road_id"`
				LaneNo int `json:"lane_no" bson:"lane_no"`
			} `json:"road_geom" bson:"road_geom"`
			RoadCondition struct {
				Rut float64 `json:"rut" bson:"rut"`
				Iri float64 `json:"iri" bson:"iri"`
				Gn  float64 `json:"gn" bson:"gn"`
			} `json:"road_condition" bson:"road_condition"`
			RoadInfo struct {
				RoadID int    `json:"road_id" bson:"road_id"`
				Name   string `json:"name" bson:"name"`
			} `json:"road_info" bson:"road_info"`
		} `json:"prepare_data" bson:"prepare_data"`
		UsedBudget   float64 `json:"used_budget" bson:"used_budget"`
		AllBudget    float64 `json:"all_budget" bson:"all_budget"`
		LeftBudget   float64 `json:"left_budget" bson:"left_budget"`
		RepireBudget float64 `json:"repire_budget" bson:"repire_budget"`
		Year         string  `json:"year" bson:"year"`
	} `json:"data" bson:"data"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type ModelResultReportMap struct {
	MaintenanceAnalysisId int    `json:"maintenance_analysis_id" bson:"maintenance_analysis_id"`
	RoadID                int    `json:"road_id" bson:"road_id"`
	Plan                  int    `json:"plan" bson:"plan"`
	PlanSequence          int    `json:"plan_sequence" bson:"plan_sequence"`
	RepairBudgetType      string `json:"repair_budget_type" bson:"repair_budget_type"`
	Relation              string `json:"relation" bson:"relation"`
	RoadNumber            int    `json:"road_number" bson:"road_number"`
	Year                  int    `json:"year" bson:"year"`
	KmStart               int    `json:"km_start" bson:"km_start"`
	AnalystYear           int    `json:"analyst_year" bson:"analyst_year"`
	KmEnd                 int    `json:"km_end" bson:"km_end"`
	Data                  struct {
		Geometry Geometry `json:"geometry" bson:"geometry"`
		// BC       bool     `json:"BC" bson:"BC"`
		// Budget   bool     `json:"budget" bson:"budget"`
		Repair bool `json:"repair" bson:"repair"`
		// Ic       bool     `json:"ic" bson:"ic"`
		IcResult struct {
			IcParamsID      int     `json:"ic_params_id" bson:"ic_params_id"`
			IcId            int     `json:"ic_id" bson:"ic_id"`
			Name            string  `json:"name" bson:"name"`
			Description     string  `json:"description" bson:"description"`
			Type            string  `json:"type" bson:"type"`
			Method          string  `json:"method" bson:"method"`
			Color           string  `json:"color" bson:"color"`
			ThicknessScrape float64 `json:"thickness_scrape" bson:"thickness_scrape"`
			ThicknessRepair float64 `json:"thickness_repair" bson:"thickness_repair"`
			PricePerUnit    int     `json:"price_per_unit" bson:"price_per_unit"`
			RefSurface      struct {
				ID               int     `json:"id" bson:"id"`
				Name             string  `json:"name" bson:"name"`
				Type             string  `json:"type" bson:"type"`
				SurfaceGroup     string  `json:"surface_group" bson:"surface_group"`
				LayerCoefficient float64 `json:"layer_coefficient" bson:"layer_coefficient"`
				Drainage         float64 `json:"drainage" bson:"drainage"`
				A                float64 `json:"a" bson:"a"`
				B                float64 `json:"b" bson:"b"`
				Crt              float64 `json:"crt" bson:"crt"`
				Raveling         struct {
					Initial struct {
						A0 float64 `json:"a0" bson:"a0"`
						A1 float64 `json:"a1" bson:"a1"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						A0 float64 `json:"a0" bson:"a1"`
						A1 float64 `json:"a1" bson:"a2"`
						A2 float64 `json:"a2" bson:"a3"`
					} `json:"progression" bson:"progression"`
				} `json:"raveling" bson:"raveling"`
				AllStructuralCrack struct {
					Initial struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
							A3 float64 `json:"a3" bson:"a3"`
							A4 float64 `json:"a4" bson:"a4"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
							A3 float64 `json:"a3" bson:"a3"`
							A4 float64 `json:"a4" bson:"a4"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"progression" bson:"progression"`
				} `json:"all_structural_crack" bson:"all_structural_crack"`
				WideStructuralCrack struct {
					Initial struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"10"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"progression" bson:"progression"`
				} `json:"wide_structural_crack" bson:"wide_structural_crack"`
				RuttingPlasticDeformation struct {
					A0 float64 `json:"a0" bson:"a0"`
					A1 float64 `json:"a1" bson:"a1"`
					A2 float64 `json:"a2" bson:"a2"`
				} `json:"rutting_plastic_deformation" bson:"rutting_plastic_deformation"`
			} `json:"ref_surface" bson:"ref_surface"`
		} `json:"ic_result" bson:"ic_result"`
		OptimizationResult struct {
			DIRI float64 `json:"dIRI" bson:"dIRI"`
			DGN  float64 `json:"dGN" bson:"dGN"`
			Cost float64 `json:"Cost" bson:"Cost"`
			BC   float64 `json:"BC" bson:"BC"`
			FV   float64 `json:"FV" bson:"FV"`
			NPV  float64 `json:"NPV" bson:"NPV"`
		} `json:"optimization_result" bson:"optimization_result"`
		RucResult struct {
			RucBeforeResult struct {
				Summary struct {
					RoadUserCostParamsID int     `json:"road_user_cost_params_id" bson:"road_user_cost_params_id"`
					AadtParamsID         int     `json:"aadt_params_id" bson:"aadt_params_id"`
					TotalPcu             float64 `json:"total_pcu" bson:"total_pcu"`
					Vkt                  float64 `json:"vkt" bson:"vkt"`
					ACCCostR             float64 `json:"ACC_Cost_r" bson:"ACC_Cost_r"`
					CMF                  float64 `json:"CMF" bson:"CMF"`
					ACC                  float64 `json:"ACC" bson:"ACC"`
					ACCFINAL             float64 `json:"ACC_FINAL" bson:"ACC_FINAL"`
					VOCFINAL             float64 `json:"VOC_FINAL" bson:"VOC_FINAL"`
					VOTFINAL             float64 `json:"VOT_FINAL" bson:"VOT_FINAL"`
					TOTALRUCNOACC        float64 `json:"TOTAL_RUC_NO_ACC" bson:"TOTAL_RUC_NO_ACC"`
					TOTALRUCWITHACC      float64 `json:"TOTAL_RUC_WITH_ACC" bson:"TOTAL_RUC_WITH_ACC"`
				} `json:"summary" bson:"summary"`
			} `json:"ruc_before_result" bson:"ruc_before_result"`
			RucAfterResult struct {
				Summary struct {
					RoadUserCostParamsID int     `json:"road_user_cost_params_id" bson:"road_user_cost_params_id"`
					AadtParamsID         int     `json:"aadt_params_id" bson:"aadt_params_id"`
					TotalPcu             float64 `json:"total_pcu" bson:"total_pcu"`
					Vkt                  float64 `json:"vkt" bson:"vkt"`
					ACCCostR             float64 `json:"ACC_Cost_r" bson:"ACC_Cost_r"`
					CMF                  float64 `json:"CMF" bson:"CMF"`
					ACC                  float64 `json:"ACC" bson:"ACC"`
					ACCFINAL             float64 `json:"ACC_FINAL" bson:"ACC_FINAL"`
					VOCFINAL             float64 `json:"VOC_FINAL" bson:"VOC_FINAL"`
					VOTFINAL             float64 `json:"VOT_FINAL" bson:"VOT_FINAL"`
					TOTALRUCNOACC        float64 `json:"TOTAL_RUC_NO_ACC" bson:"TOTAL_RUC_NO_ACC"`
					TOTALRUCWITHACC      float64 `json:"TOTAL_RUC_WITH_ACC" bson:"TOTAL_RUC_WITH_ACC"`
				} `json:"summary" bson:"summary"`
			} `json:"ruc_after_result" bson:"ruc_after_result"`
		} `json:"ruc_result" bson:"ruc_result"`
		RweResult struct {
			Age struct {
				YearLastOverlay        int `json:"year_last_overlay" bson:"year_last_overlay"`
				YearLastSeal           int `json:"year_last_seal" bson:"year_last_seal"`
				YearLastMolRcl         int `json:"year_last_mol_rcl" bson:"year_last_mol_rcl"`
				YearLastReconstruction int `json:"year_last_reconstruction" bson:"year_last_reconstruction"`
				Age                    int `json:"age" bson:"age"`
			} `json:"age" bson:"age"`
			AreaAcIcrack   float64 `json:"area_ac_icrack" bson:"area_ac_icrack"`
			AreaAcUcrack   float64 `json:"area_ac_ucrack" bson:"area_ac_ucrack"`
			CurrentSurface struct {
				ID               int     `json:"id" bson:"id"`
				Name             string  `json:"name" bson:"name"`
				Type             string  `json:"type" bson:"type"`
				SurfaceGroup     string  `json:"surface_group" bson:"surface_group"`
				LayerCoefficient float64 `json:"layer_coefficient" bson:"layer_coefficient"`
				Drainage         float64 `json:"drainage" bson:"drainage"`
				A                float64 `json:"a" bson:"a"`
				B                float64 `json:"b" bson:"b"`
				Crt              float64 `json:"crt" bson:"crt"`
				Rrf              float64 `json:"rrf" bson:"rrf"`
				Raveling         struct {
					Initial struct {
						A0 float64 `json:"a0" bson:"a0"`
						A1 float64 `json:"a1" bson:"a1"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						A0 float64 `json:"a0" bson:"a0"`
						A1 float64 `json:"a1" bson:"a1"`
						A2 float64 `json:"a2" bson:"a2"`
					} `json:"progression" bson:"progression"`
				} `json:"raveling" bson:"raveling"`
				AllStructuralCrack struct {
					Initial struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
							A3 float64 `json:"a3" bson:"a3"`
							A4 float64 `json:"a4" bson:"a4"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
							A3 float64 `json:"a3" bson:"a3"`
							A4 float64 `json:"a4" bson:"a4"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"progression" bson:"progression"`
				} `json:"all_structural_crack" bson:"all_structural_crack"`
				WideStructuralCrack struct {
					Initial struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
							A2 float64 `json:"a2" bson:"a2"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"initial" bson:"initial"`
					Progression struct {
						HSOLDO struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD_O" bson:"HSOLD_O"`
						HSOLD struct {
							A0 float64 `json:"a0" bson:"a0"`
							A1 float64 `json:"a1" bson:"a1"`
						} `json:"HSOLD" bson:"HSOLD"`
					} `json:"progression" bson:"progression"`
				} `json:"wide_structural_crack" bson:"wide_structural_crack"`
				RuttingPlasticDeformation struct {
					A0 float64 `json:"a0" bson:"a0"`
					A1 float64 `json:"a1" bson:"a1"`
					A2 float64 `json:"a2" bson:"a2"`
				} `json:"rutting_plastic_deformation" bson:"rutting_plastic_deformation"`
			} `json:"current_surface" bson:"current_surface"`
			Gn         float64 `json:"gn" bson:"gn"`
			HsoldHsnew struct {
				Hsold float64 `json:"hsold" bson:"hsold"`
				Hsnew float64 `json:"hsnew" bson:"hsnew"`
			} `json:"hsold_hsnew" bson:"hsold_hsnew"`
			Iri                float64 `json:"iri" bson:"iri"`
			NumberOfPothole    float64 `json:"number_of_pothole" bson:"number_of_pothole"`
			PercentAcIcrack    float64 `json:"percent_ac_icrack" bson:"percent_ac_icrack"`
			PercentAcRavelling float64 `json:"percent_ac_ravelling" bson:"percent_ac_ravelling"`
			PercentAcUcrack    float64 `json:"percent_ac_ucrack" bson:"percent_ac_ucrack"`
			Rut                float64 `json:"rut" bson:"rut"`
			Snp                struct {
			} `json:"snp" bson:"snp"`
		} `json:"rwe_result" bson:"rwe_result"`
		PrepareData struct {
			Area    float64 `json:"area" bson:"area"`
			Aadt    float64 `json:"aadt" bson:"aadt"`
			KmStart int     `json:"km_start" bson:"km_start"`
			KmEnd   int     `json:"km_end" bson:"km_end"`
			Length  int     `json:"length" bson:"length"`
			Road    struct {
				RoadID         int    `json:"road_id" bson:"road_id"`
				RefDirectionID int    `json:"ref_direction_id" bson:"ref_direction_id"`
				RoadName       string `json:"road_name" bson:"road_name"`
				RoadGroupName  string `json:"road_group_name" bson:"road_group_name"`
			} `json:"road" bson:"road"`
			RoadGeom struct {
				RoadID int `json:"road_id" bson:"road_id"`
				LaneNo int `json:"lane_no" bson:"lane_no"`
			} `json:"road_geom" bson:"road_geom"`
			RoadCondition struct {
				Rut float64 `json:"rut" bson:"rut"`
				Iri float64 `json:"iri" bson:"iri"`
				Gn  float64 `json:"gn" bson:"gn"`
			} `json:"road_condition" bson:"road_condition"`
			RoadInfo struct {
				RoadID int    `json:"road_id" bson:"road_id"`
				Name   string `json:"name" bson:"name"`
			} `json:"road_info" bson:"road_info"`
		} `json:"prepare_data" bson:"prepare_data"`
		UsedBudget          float64 `json:"used_budget" bson:"used_budget"`
		AllBudget           float64 `json:"all_budget" bson:"all_budget"`
		LeftBudget          float64 `json:"left_budget" bson:"left_budget"`
		RepireBudget        float64 `json:"repire_budget" bson:"repire_budget"`
		Year                string  `json:"year" bson:"year"`
		IriAfter            float64 `json:"iri_after" bson:"iri_after"`
		DeteriorationResult struct {
			IsSurfaceAc bool `json:"is_surface_ac" bson:"is_surface_ac"`
			Result      struct {
				Iri float64 `json:"iri" bson:"iri"`
			} `json:"result" bson:"result"`
		} `json:"deterioration_result" bson:"deterioration_result"`
	} `json:"data" bson:"data"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}
