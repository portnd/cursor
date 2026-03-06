package requests

import "mime/multipart"

type RefRequest struct {
	Name string `json:"name" validate:"nonzero" example:"input string"`
}

type UpdateConditionListRequets struct {
	OwnerId             int             `json:"owner_id" validate:"nonzero" extensions:"x-order=0"`
	RefConditionRangeID int             `json:"ref_condition_range_id" validate:"nonzero" extensions:"x-order=1"`
	ConditionList       []ConditionList `json:"condition_list" extensions:"x-order=2"`
}

type RoadLineList struct {
	GradeID        int     `json:"grade_id" example:"1" extensions:"x-order=0"`
	LeftValueWhite float64 `json:"left_value_white" example:"10" extensions:"x-order=1"`
	// LeftCondition  string  `json:"left_condition" example:"<="  extensions:"x-order=2"`
	RightValueWhite float64 `json:"right_value_white" example:"100" extensions:"x-order=3"`
	// RightCondition string  `json:"right_condition" example:"<" extensions:"x-order=4"`
	LeftValueYellow  float64 `json:"left_value_yellow" example:"10" extensions:"x-order=4"`
	RightValueYellow float64 `json:"right_value_yellow" example:"100" extensions:"x-order=5"`

	// ConditionType string `json:"condition_type" example:"IRI" extensions:"x-order=6"`
}

type ConditionList struct {
	GradeID     int     `json:"grade_id" example:"1" extensions:"x-order=0"`
	LeftValueAC float64 `json:"left_value_ac" example:"10" extensions:"x-order=1"`
	// LeftCondition  string  `json:"left_condition" example:"<="  extensions:"x-order=2"`
	RightValueAC float64 `json:"right_value_ac" example:"100" extensions:"x-order=3"`
	// RightCondition string  `json:"right_condition" example:"<" extensions:"x-order=4"`
	LeftValueCC  float64 `json:"left_value_cc" example:"10" extensions:"x-order=4"`
	RightValueCC float64 `json:"right_value_cc" example:"100" extensions:"x-order=5"`

	ConditionType string `json:"condition_type" example:"IRI" extensions:"x-order=6"`
}

type SignImageRequest struct {
	Name        string                `form:"name" extensions:"x-order=0"`
	Abbr        string                `form:"abbr" extensions:"x-order=1"`
	ImageStatus string                `form:"image_status" extensions:"x-order=3"`
	Image       *multipart.FileHeader `form:"image" swaggerignore:"true" extensions:"x-order=2"`
}

type SwaggerAssetTable struct {
	Data         AssetTableData `json:"data" extension:"x-order=0"`
	IconFilePath string         `json:"icon_filepath" extensions:"x-order=1" example:"the actual type of icon_filepath is file type, not string. the type string use for example only"`
}

type AssetTableRequest struct {
	Data               AssetTableData        `form:"data"`
	IconFilePathStatus string                `form:"icon_filepath_status"`
	IconFilePath       *multipart.FileHeader `form:"icon_filepath"`
}

type AssetTableData struct {
	TableName     string `json:"table_name" extensions:"x-order=0"`
	TableLabel    string `json:"table_label" extensions:"x-order=1"`
	AssetType     string `json:"asset_type" extensions:"x-order=2"`
	AssetGroup    int    `json:"asset_group" extensions:"x-order=3"`
	GeomType      string `json:"geom_type" extensions:"x-order=4"`
	LineColorCode string `json:"line_color_code" extensions:"x-order=5"`
	// ApproverID    []int     `json:"approver_id" extensions:"x-order=6"`
	// ViewerID      []int     `json:"viewer_id" extensions:"x-order=7"`
	Columns       []Columns `json:"columns" extensions:"x-order=9"`
	DeleteColumns []int     `json:"delete_columns" extensions:"x-order=10"`
}

type Columns struct {
	ColumnID        int    `json:"column_id" extensions:"x-order=1"`
	ColumnName      string `json:"column_name" extensions:"x-order=2"`
	TableNameRef    string `json:"table_name_ref" extensions:"x-order=3"`
	ComponentTitle  string `json:"component_title" extensions:"x-order=4"`
	ComponentType   string `json:"component_type" extensions:"x-order=5"`
	IsRequired      bool   `json:"is_required" extensions:"x-order=6"`
	IsVisibleView   bool   `json:"is_visible_view" extensions:"x-order=7"`
	IsVisibleEdit   bool   `json:"is_visible_edit" extensions:"x-order=8"`
	IsVisibleReport bool   `json:"is_visible_report" extensions:"x-order=9"`
}

type QueryParams struct {
	Page                string `form:"page"`
	Limit               string `form:"limit"`
	Name                string `form:"name"`
	RefConditionRangeID int    `form:"ref_condition_range_id" `
}

type QueryParamsReflectivityRange struct {
	Page                   string `form:"page"`
	Limit                  string `form:"limit"`
	Name                   string `form:"name"`
	RefReflectivityRangeID int    `form:"ref_reflectivity_range_id" `
}

type AssetTableQueryParams struct {
	QueryParams
	AssetType string `form:"asset_type"`
	GroupID   string `form:"group_id"`
}

type OwnerRequest struct {
	Name                string          `json:"name" example:"EXAT" validate:"nonzero" extensions:"x-order=1"`
	RefConditionRangeID int             `json:"ref_condition_range_id" validate:"nonzero" extensions:"x-order=2"`
	ConditionList       []ConditionList `json:"condition_list" extensions:"x-order=3"`
}

type OwnerRoadLineRequest struct {
	Name                   string         `json:"name" example:"EXAT" validate:"nonzero" extensions:"x-order=1"`
	RefReflectivityRangeID int            `json:"ref_reflectivity_range_id" validate:"nonzero" extensions:"x-order=2"`
	ConditionList          []RoadLineList `json:"road_line_list" extensions:"x-order=3"`
}

type Budget struct {
	Name   string         `json:"name" validate:"nonzero"`
	Budget []BudgetMethod `json:"budget" `
}

type BudgetMethod struct {
	CostPerUnit *float64 `json:"cost_per_unit"`
	MethodName  string   `json:"method_name" validate:"nonzero"`
}

type UpdateBudget struct {
	Id     int                  `json:"id" validate:"nonzero"`
	Name   string               `json:"name" validate:"nonzero"`
	Budget []UpdateBudgetMethod `json:"budget" `
}

type UpdateBudgetMethod struct {
	Id          int      `json:"id"`
	CostPerUnit *float64 `json:"cost_per_unit"`
	MethodName  string   `json:"method_name" validate:"nonzero"`
}

type InterventionCriteria struct {
	Asphalt  []InterventionCriteriaGroup `json:"asphalt"`
	Concrete []InterventionCriteriaGroup `json:"concrete"`
}

type InterventionCriteriaGroup struct {
	ID                        int                        `json:"id"`
	Name                      string                     `json:"name"`
	InterventionCriteriaDatas []InterventionCriteriaData `json:"intervention_criterias"`
}

type InterventionCriteriaData struct {
	Id    int  `json:"id"`
	IsNew bool `json:"is_new"`
	// MaintenanceMethodID      int                             `json:"maintenance_method_id"`
	// MaintenanceMethod        string                          `json:"maintenance_method"`
	MaintenanceCostPerUnit   float64                         `json:"maintenance_cost_per_unit"`
	MaintenanceDescription   string                          `json:"maintenance_description"`
	MaintenanceScraping      float64                         `json:"maintenance_scraping"`
	MaintenanceStandardName  string                          `json:"maintenance_standard_name"`
	MaintenanceSurfaceTypeId int                             `json:"maintenance_surface_type_id"`
	MaintenanceThickness     float64                         `json:"maintenance_thickness"`
	MaintenanceType          string                          `json:"maintenance_type"`
	MaintenanceCondition     []InterventionCriteriaCindition `json:"maintenance_condition"`
}

type InterventionCriteriaCindition struct {
	Id                  int     `json:"id"`
	IsNew               bool    `json:"is_new"`
	ConditionCriterion  string  `json:"condition_criterion"`
	ConditionLink       string  `json:"condition_link"`
	ConditionOperation1 string  `json:"condition_operation_1"`
	ConditionOperation2 string  `json:"condition_operation_2"`
	ConditionValue1     float64 `json:"condition_value_1"`
	ConditionValue2     float64 `json:"condition_value_2"`
}

type InterventionCriteriaSequenceCriteriaMethod struct {
	Concrete []int `json:"concrete"`
	Asphalt  []int `json:"asphalt"`
}

type InterventionCriteriaSequence struct {
	Id int `json:"id"`
}

type CreateAadtGrowthRate struct {
	RoadGroupId int     `json:"road_group_id" validate:"nonzero"`
	R           float64 `json:"r"`
}

type AadtPercentageVehicleType struct {
	RoadGroupId   int                           `json:"road_group_id" validate:"nonzero"`
	FourWheel     AadtPercentageVehicleTypeFour `json:"four_wheel"`
	SixToTenWheel AadtPercentageVehicleTypeSix  `json:"six_to_ten_wheel"`
	OverTenWheel  AadtPercentageVehicleTypeTen  `json:"over_ten_wheel"`
}

type AadtPercentageVehicleTypeFour struct {
	CarLessThanEqualSeven float64 `json:"car_less_than_equal_seven"`
	CarOverThanSeven      float64 `json:"car_over_than_seven"`
	LightBus              float64 `json:"light_bus"`
	LightTruck            float64 `json:"light_truck"`
}

type AadtPercentageVehicleTypeSix struct {
	MediumBus   float64 `json:"medium_bus"`
	MediumTruck float64 `json:"medium_truck"`
}

type AadtPercentageVehicleTypeTen struct {
	HeavyBus    float64 `json:"heavy_bus"`
	HeavyTruck  float64 `json:"heavy_truck"`
	FullTrailor float64 `json:"full_trailor"`
	SemiTrailor float64 `json:"semi_trailor"`
}

type CreateAadtParameter struct {
	RoadGroupId                   int     `json:"road_group_id" validate:"nonzero"`
	Elane                         float64 `json:"elane"`
	FourWheelAxleNumber           int     `json:"four_wheel_axle_number"`
	FourWheelVehicleVolume        float64 `json:"four_wheel_vehicle_volume"`
	SixWheelAxleNumberId          int     `json:"six_wheel_axle_number_id"`
	SixWheelVehicleVolume         float64 `json:"six_wheel_vehicle_volume"`
	SixWheelPercentageTruck       float64 `json:"six_wheel_percentage_truck"`
	SixWheelFactorResult          float64 `json:"six_wheel_factor_result"`
	TenWheelAxleNumberId          int     `json:"ten_wheel_axle_number_id"`
	TenWheelVehicleVolume         float64 `json:"ten_wheel_vehicle_volume"`
	TenWheelPercentageTruck       float64 `json:"ten_wheel_percentage_truck"`
	TenWheelFactorResult          float64 `json:"ten_wheel_factor_result"`
	IsTruckFactor                 bool    `json:"is_truck_factor"`
	SpeedAverage                  float64 `json:"speed_average"`
	SpeedHeavyTruck               float64 `json:"speed_heavy_truck"`
	LaneDistributionFactor        float64 `json:"lane_distribution_factor"`
	DirectionalDistributionFactor float64 `json:"directional_distribution_factor"`
}

type CalculateAadtParameter struct {
	RoadGroupId            int `json:"road_group_id" validate:"nonzero"`
	ParameterVehicleTypeId int `json:"parameter_vehicle_type_id" validate:"nonzero"`
}

type SettingRoadWorkEffectAsphalt struct {
	AsOlOverlayA0                           float64 `json:"as_ol_overlay_a0"`
	AsOlArVb                                float64 `json:"as_ol_ar_vb"`
	AsOlPoTb                                float64 `json:"as_ol_po_tb"`
	AsOlAcAb                                float64 `json:"as_ol_ac_ab"`
	AsOlRdMb                                float64 `json:"as_ol_rd_mb"`
	AsSsRweSsModelA0                        float64 `json:"as_ss_rwe_ss_model_a0"`
	AsSsDefaultLowerBoundIriAfterSlurrySeal float64 `json:"as_ss_default_lower_bound_iri_after_slurry_seal"`
	AsSsArVb                                float64 `json:"as_ss_ar_vb"`
	AsSsApoTb                               float64 `json:"as_ss_apo_tb"`
	AsSsAcAb                                float64 `json:"as_ss_ac_ab"`
	AsSsRdMb                                float64 `json:"as_ss_rd_mb"`
	AsMolIriAfterMillOverlay                float64 `json:"as_mol_iri_after_mill_overlay"`
	AsMolArVb                               float64 `json:"as_mol_ar_vb"`
	AsMolApoTb                              float64 `json:"as_mol_apo_tb"`
	AsMolAcAb                               float64 `json:"as_mol_ac_ab"`
	AsMolRdMb                               float64 `json:"as_mol_rd_mb"`
	AsRclSnc                                float64 `json:"as_rcl_snc"`
	AsRclIriAfterRecycling                  float64 `json:"as_rcl_iri_after_recycling"`
	AsRclArVb                               float64 `json:"as_rcl_ar_vb"`
	AsRclApoTb                              float64 `json:"as_rcl_apo_tb"`
	AsRclAcAb                               float64 `json:"as_rcl_ac_ab"`
	AsRclRdMb                               float64 `json:"as_rcl_rd_mb"`
	AsRclDefaultHsOld                       float64 `json:"as_rcl_default_hs_old"`
	AsRcSnc                                 float64 `json:"as_rc_snc"`
	AsRcIriAfterReconstruction              float64 `json:"as_rc_iri_after_reconstruction"`
	AsRcArVb                                float64 `json:"as_rc_ar_vb"`
	AsRcApoTb                               float64 `json:"as_rc_apo_tb"`
	AsRcAcAb                                float64 `json:"as_rc_ac_ab"`
	AsRcRdMb                                float64 `json:"as_rc_rd_mb"`
}

type SettingRoadWorkEffectConcrete struct {
	CcFdrIriAfterFdr     float64 `json:"cc_fdr_iri_after_fdr"`
	CcFdrFaulting        float64 `json:"cc_fdr_faulting"`
	CcFdrCracking        float64 `json:"cc_fdr_cracking"`
	CcFdrSpalling        float64 `json:"cc_fdr_spalling"`
	CcBcoIriAfterBco     float64 `json:"cc_bco_iri_after_bco"`
	CcBcoFaulting        float64 `json:"cc_bco_faulting"`
	CcBcoCracking        float64 `json:"cc_bco_cracking"`
	CcBcoSpalling        float64 `json:"cc_bco_spalling"`
	CcMolIriAfterMol     float64 `json:"cc_mol_iri_after_mol"`
	CcMolFaulting        float64 `json:"cc_mol_faulting"`
	CcMolCracking        float64 `json:"cc_mol_cracking"`
	CcMolSpalling        float64 `json:"cc_mol_spalling"`
	CcSealIriAfterSeal   float64 `json:"cc_seal_iri_after_seal"`
	CcSealFaulting       float64 `json:"cc_seal_faulting"`
	CcSealCracking       float64 `json:"cc_seal_cracking"`
	CcSealSpalling       float64 `json:"cc_seal_spalling"`
	CcRbcIri             float64 `json:"cc_rbc_iri"`
	CcRbcSlabthk         float64 `json:"cc_rbc_slabthk"`
	CcRbcPercentFaulting float64 `json:"cc_rbc_percent_faulting"`
	CcRbcPercentSpalling float64 `json:"cc_rbc_percent_spalling"`
	CcRbcPercentCracking float64 `json:"cc_rbc_percent_cracking"`
}

type RoadUserCostAccLossValue struct {
	ValueOfFatalAccidents               *float64 `json:"value_of_fatal_accidents"`
	ValueOfAccidentsWithSeriousInjuries *float64 `json:"value_of_accidents_with_serious_injuries"`
	ValueOfAccidentsWithMinorInjuries   *float64 `json:"value_of_accidents_with_minor_injuries"`
	ValueOfAccidentsWithPropertyDamaged *float64 `json:"value_of_accidents_with_property_damaged"`
}

type RoadUserCostAccChanceOfAccident struct {
	RoadGroupId                          int      `json:"road_group_id" validate:"nonzero"`
	NumberOfFatalAccidents               *float64 `json:"number_of_fatal_accidents"`
	NumberOfAccidentsWithSeriousInjuries *float64 `json:"number_of_accidents_with_serious_injuries"`
	NumberOfAccidentsWithMinorInjuries   *float64 `json:"number_of_accidents_with_minor_injuries"`
	NumberOfAccidentsWithPropertyDamage  *float64 `json:"number_of_accidents_with_property_damaged"`
}

type RoadUserCostRusDefaultData struct {
	CarLessThanEqualSeven RoadUserCostRusDefaultDataParameters `json:"car_less_than_equal_seven"`
	CarOverThanSeven      RoadUserCostRusDefaultDataParameters `json:"car_over_than_seven"`
	LightBus              RoadUserCostRusDefaultDataParameters `json:"light_bus"`
	MediumBus             RoadUserCostRusDefaultDataParameters `json:"medium_bus"`
	HeavyBus              RoadUserCostRusDefaultDataParameters `json:"heavy_bus"`
	LightTruck            RoadUserCostRusDefaultDataParameters `json:"light_truck"`
	MediumTruck           RoadUserCostRusDefaultDataParameters `json:"medium_truck"`
	HeavyTruck            RoadUserCostRusDefaultDataParameters `json:"heavy_truck"`
	FullTrailor           RoadUserCostRusDefaultDataParameters `json:"full_trailor"`
	SemiTrailor           RoadUserCostRusDefaultDataParameters `json:"semi_trailor"`
}

type RoadUserCostRusDefaultDataParameters struct {
	VehicleName string   `json:"vehicle_name"`
	FuelUCost   *float64 `json:"fuel_u_cost"`
	OilUCost    *float64 `json:"oil_u_cost"`
	TypeUCost   *float64 `json:"type_u_cost"`
	VehUCost    *float64 `json:"veh_u_cost"`
	MUpper      *float64 `json:"m_upper"`
	MLower      *float64 `json:"m_lower"`
	Wheels      *float64 `json:"wheels"`
	NumPass     *float64 `json:"num_pass"`
	TUCost      *float64 `json:"t_u_cost"`
}

type RoadUserCostRusDriving struct {
	CarLessThanEqualSeven RoadUserCostRusDrivingParameters `json:"car_less_than_equal_seven"`
	CarOverThanSeven      RoadUserCostRusDrivingParameters `json:"car_over_than_seven"`
	LightBus              RoadUserCostRusDrivingParameters `json:"light_bus"`
	MediumBus             RoadUserCostRusDrivingParameters `json:"medium_bus"`
	HeavyBus              RoadUserCostRusDrivingParameters `json:"heavy_bus"`
	LightTruck            RoadUserCostRusDrivingParameters `json:"light_truck"`
	MediumTruck           RoadUserCostRusDrivingParameters `json:"medium_truck"`
	HeavyTruck            RoadUserCostRusDrivingParameters `json:"heavy_truck"`
	FullTrailor           RoadUserCostRusDrivingParameters `json:"full_trailor"`
	SemiTrailor           RoadUserCostRusDrivingParameters `json:"semi_trailor"`
}

type RoadUserCostRusDrivingParameters struct {
	CrA1  *float64 `json:"cr_a1"`
	CrA2  *float64 `json:"cr_a2"`
	P     *float64 `json:"p"`
	Cd    *float64 `json:"cd"`
	Cduml *float64 `json:"cdmul"`
	Ad    *float64 `json:"af"`
}

type RoadUserCostRusEngineSpeed struct {
	CarLessThanEqualSeven RoadUserCostRusEngineSpeedParameters `json:"car_less_than_equal_seven"`
	CarOverThanSeven      RoadUserCostRusEngineSpeedParameters `json:"car_over_than_seven"`
	LightBus              RoadUserCostRusEngineSpeedParameters `json:"light_bus"`
	MediumBus             RoadUserCostRusEngineSpeedParameters `json:"medium_bus"`
	HeavyBus              RoadUserCostRusEngineSpeedParameters `json:"heavy_bus"`
	LightTruck            RoadUserCostRusEngineSpeedParameters `json:"light_truck"`
	MediumTruck           RoadUserCostRusEngineSpeedParameters `json:"medium_truck"`
	HeavyTruck            RoadUserCostRusEngineSpeedParameters `json:"heavy_truck"`
	FullTrailor           RoadUserCostRusEngineSpeedParameters `json:"full_trailor"`
	SemiTrailor           RoadUserCostRusEngineSpeedParameters `json:"semi_trailor"`
}

type RoadUserCostRusEngineSpeedParameters struct {
	RpmA0   *float64 `json:"rpm_a0"`
	RpmA1   *float64 `json:"rpm_a1"`
	RpmA2   *float64 `json:"rpm_a2"`
	RpmIdle *float64 `json:"rpm_idle"`
	Rpm100  *float64 `json:"rpm100"`
}

type RoadUserCostRusFuelConsumption struct {
	CarLessThanEqualSeven RoadUserCostRusFuelConsumptionParameters `json:"car_less_than_equal_seven"`
	CarOverThanSeven      RoadUserCostRusFuelConsumptionParameters `json:"car_over_than_seven"`
	LightBus              RoadUserCostRusFuelConsumptionParameters `json:"light_bus"`
	MediumBus             RoadUserCostRusFuelConsumptionParameters `json:"medium_bus"`
	HeavyBus              RoadUserCostRusFuelConsumptionParameters `json:"heavy_bus"`
	LightTruck            RoadUserCostRusFuelConsumptionParameters `json:"light_truck"`
	MediumTruck           RoadUserCostRusFuelConsumptionParameters `json:"medium_truck"`
	HeavyTruck            RoadUserCostRusFuelConsumptionParameters `json:"heavy_truck"`
	FullTrailor           RoadUserCostRusFuelConsumptionParameters `json:"full_trailor"`
	SemiTrailor           RoadUserCostRusFuelConsumptionParameters `json:"semi_trailor"`
}

type RoadUserCostRusFuelConsumptionParameters struct {
	IdleFuel *float64 `json:"idle_fuel"`
	DfFuel   *float64 `json:"df_fuel"`
	ZeTab    *float64 `json:"ze_tab"`
	Ehp      *float64 `json:"ehp"`
	Edt      *float64 `json:"edt"`
	Prat     *float64 `json:"prat"`
	Kpea     *float64 `json:"kpea"`
	PaccsA0  *float64 `json:"paccs_a0"`
	PctPeng  *float64 `json:"pct_peng"`
}

type RoadUserCostRusLubricantConsumption struct {
	CarLessThanEqualSeven RoadUserCostRusLubricantConsumptionParameters `json:"car_less_than_equal_seven"`
	CarOverThanSeven      RoadUserCostRusLubricantConsumptionParameters `json:"car_over_than_seven"`
	LightBus              RoadUserCostRusLubricantConsumptionParameters `json:"light_bus"`
	MediumBus             RoadUserCostRusLubricantConsumptionParameters `json:"medium_bus"`
	HeavyBus              RoadUserCostRusLubricantConsumptionParameters `json:"heavy_bus"`
	LightTruck            RoadUserCostRusLubricantConsumptionParameters `json:"light_truck"`
	MediumTruck           RoadUserCostRusLubricantConsumptionParameters `json:"medium_truck"`
	HeavyTruck            RoadUserCostRusLubricantConsumptionParameters `json:"heavy_truck"`
	FullTrailor           RoadUserCostRusLubricantConsumptionParameters `json:"full_trailor"`
	SemiTrailor           RoadUserCostRusLubricantConsumptionParameters `json:"semi_trailor"`
}

type RoadUserCostRusLubricantConsumptionParameters struct {
	OilCont *float64 `json:"oil_cont"`
	OilOper *float64 `json:"oil_oper"`
}

type RoadUserCostRusWasteOfConsumption struct {
	CarLessThanEqualSeven RoadUserCostRusWasteOfConsumptionParameters `json:"car_less_than_equal_seven"`
	CarOverThanSeven      RoadUserCostRusWasteOfConsumptionParameters `json:"car_over_than_seven"`
	LightBus              RoadUserCostRusWasteOfConsumptionParameters `json:"light_bus"`
	MediumBus             RoadUserCostRusWasteOfConsumptionParameters `json:"medium_bus"`
	HeavyBus              RoadUserCostRusWasteOfConsumptionParameters `json:"heavy_bus"`
	LightTruck            RoadUserCostRusWasteOfConsumptionParameters `json:"light_truck"`
	MediumTruck           RoadUserCostRusWasteOfConsumptionParameters `json:"medium_truck"`
	HeavyTruck            RoadUserCostRusWasteOfConsumptionParameters `json:"heavy_truck"`
	FullTrailor           RoadUserCostRusWasteOfConsumptionParameters `json:"full_trailor"`
	SemiTrailor           RoadUserCostRusWasteOfConsumptionParameters `json:"semi_trailor"`
}

type RoadUserCostRusWasteOfConsumptionParameters struct {
	C0Tc  *float64 `json:"c0_tc"`
	CtCte *float64 `json:"ct_cte"`
	Vol   *float64 `json:"vol"`
}

type RoadUserCostRusMaintenance struct {
	CarLessThanEqualSeven RoadUserCostRusMaintenanceParameters `json:"car_less_than_equal_seven"`
	CarOverThanSeven      RoadUserCostRusMaintenanceParameters `json:"car_over_than_seven"`
	LightBus              RoadUserCostRusMaintenanceParameters `json:"light_bus"`
	MediumBus             RoadUserCostRusMaintenanceParameters `json:"medium_bus"`
	HeavyBus              RoadUserCostRusMaintenanceParameters `json:"heavy_bus"`
	LightTruck            RoadUserCostRusMaintenanceParameters `json:"light_truck"`
	MediumTruck           RoadUserCostRusMaintenanceParameters `json:"medium_truck"`
	HeavyTruck            RoadUserCostRusMaintenanceParameters `json:"heavy_truck"`
	FullTrailor           RoadUserCostRusMaintenanceParameters `json:"full_trailor"`
	SemiTrailor           RoadUserCostRusMaintenanceParameters `json:"semi_trailor"`
}

type RoadUserCostRusMaintenanceParameters struct {
	Kpc   *float64 `json:"kpc"`
	Akmo  *float64 `json:"akmo"`
	Life0 *float64 `json:"life0"`
	Kp    *float64 `json:"kp"`
	A0    *float64 `json:"a0"`
	A1    *float64 `json:"a1"`
	CpCon *float64 `json:"cp_con"`
}

type RoadUserCostRusTravelTime struct {
	CarLessThanEqualSeven RoadUserCostRusTravelTimeParameters `json:"car_less_than_equal_seven"`
	CarOverThanSeven      RoadUserCostRusTravelTimeParameters `json:"car_over_than_seven"`
	LightBus              RoadUserCostRusTravelTimeParameters `json:"light_bus"`
	MediumBus             RoadUserCostRusTravelTimeParameters `json:"medium_bus"`
	HeavyBus              RoadUserCostRusTravelTimeParameters `json:"heavy_bus"`
	LightTruck            RoadUserCostRusTravelTimeParameters `json:"light_truck"`
	MediumTruck           RoadUserCostRusTravelTimeParameters `json:"medium_truck"`
	HeavyTruck            RoadUserCostRusTravelTimeParameters `json:"heavy_truck"`
	FullTrailor           RoadUserCostRusTravelTimeParameters `json:"full_trailor"`
	SemiTrailor           RoadUserCostRusTravelTimeParameters `json:"semi_trailor"`
}

type RoadUserCostRusTravelTimeParameters struct {
	PcTwk *float64 `json:"pc_twk"`
}

type RoadUserCostRusVehicleSpeedCalculation struct {
	CarLessThanEqualSeven RoadUserCostRusVehicleSpeedCalculationParameters `json:"car_less_than_equal_seven"`
	CarOverThanSeven      RoadUserCostRusVehicleSpeedCalculationParameters `json:"car_over_than_seven"`
	LightBus              RoadUserCostRusVehicleSpeedCalculationParameters `json:"light_bus"`
	MediumBus             RoadUserCostRusVehicleSpeedCalculationParameters `json:"medium_bus"`
	HeavyBus              RoadUserCostRusVehicleSpeedCalculationParameters `json:"heavy_bus"`
	LightTruck            RoadUserCostRusVehicleSpeedCalculationParameters `json:"light_truck"`
	MediumTruck           RoadUserCostRusVehicleSpeedCalculationParameters `json:"medium_truck"`
	HeavyTruck            RoadUserCostRusVehicleSpeedCalculationParameters `json:"heavy_truck"`
	FullTrailor           RoadUserCostRusVehicleSpeedCalculationParameters `json:"full_trailor"`
	SemiTrailor           RoadUserCostRusVehicleSpeedCalculationParameters `json:"semi_trailor"`
}

type RoadUserCostRusVehicleSpeedCalculationParameters struct {
	Cw1     *float64 `json:"cw1"`
	Cw2     *float64 `json:"cw2"`
	Cw3     *float64 `json:"cw3"`
	A2      *float64 `json:"a2"`
	A3      *float64 `json:"a3"`
	Pd      *float64 `json:"pd"`
	Pb      *float64 `json:"pb"`
	ArvMax  *float64 `json:"arv_max"`
	AUpper0 *float64 `json:"a_upper0"`
	ALower0 *float64 `json:"a_lower0"`
	A1      *float64 `json:"a1"`
}

type RoadUserCostRusTrafficData struct {
	CarLessThanEqualSeven RoadUserCostRusTrafficDataParameters `json:"car_less_than_equal_seven"`
	CarOverThanSeven      RoadUserCostRusTrafficDataParameters `json:"car_over_than_seven"`
	LightBus              RoadUserCostRusTrafficDataParameters `json:"light_bus"`
	MediumBus             RoadUserCostRusTrafficDataParameters `json:"medium_bus"`
	HeavyBus              RoadUserCostRusTrafficDataParameters `json:"heavy_bus"`
	LightTruck            RoadUserCostRusTrafficDataParameters `json:"light_truck"`
	MediumTruck           RoadUserCostRusTrafficDataParameters `json:"medium_truck"`
	HeavyTruck            RoadUserCostRusTrafficDataParameters `json:"heavy_truck"`
	FullTrailor           RoadUserCostRusTrafficDataParameters `json:"full_trailor"`
	SemiTrailor           RoadUserCostRusTrafficDataParameters `json:"semi_trailor"`
}

type RoadUserCostRusTrafficDataParameters struct {
	PcuEquivalent *float64 `json:"pcu_equivalent"`
}

type Optimization struct {
	BcRatioConstraint *float64 `json:"bc_ratio_constraint"`
	DefaultDesignLife *float64 `json:"default_design_life"`
}

type DeteriorationAsphalt struct {
	RoadGroupId int     `json:"road_group_id"  validate:"nonzero"`
	Tlf         float64 `json:"tlf"`
	Cdb         float64 `json:"cdb"`
	Cds         float64 `json:"cds"`
	Comp        float64 `json:"comp"`
	Kvi         float64 `json:"kvi"`
	Kvp         float64 `json:"kvp"`
	Kpi         float64 `json:"kpi"`
	Kpp         float64 `json:"kpp"`
	Krid        float64 `json:"krid"`
	Krst        float64 `json:"krst"`
	Krpd        float64 `json:"krpd"`
	Kgm         float64 `json:"kgm"`
	Kgp         float64 `json:"kgp"`
	Kcia        float64 `json:"kcia"`
	Cmod        float64 `json:"cmod"`
	Kciw        float64 `json:"kciw"`
	Kcpa        float64 `json:"kcpa"`
	Kcpw        float64 `json:"kcpw"`
}

type DeteriorationConcrete struct {
	RoadGroupId int     `json:"road_group_id"  validate:"nonzero"`
	PSteel      float64 `json:"p_steel"`
	Ec          float64 `json:"ec"`
	Mi          float64 `json:"mi"`
	Fi          float64 `json:"fi"`
	Kjrc        float64 `json:"kjrc"`
	BStress     float64 `json:"b_stress"`
	JtSpace     float64 `json:"jt_space"`
	Kjrf        float64 `json:"kjrf"`
	Widened     float64 `json:"widened"`
	PredSeal    float64 `json:"pred_seal"`
	DwlCor      float64 `json:"dwl_cor"`
	Kjrs        float64 `json:"kjrs"`
	Kjrr        float64 `json:"kjrr"`
}

type CreateRefHris struct {
	RoadNumber           string `json:"road_number"`
	OfficeOfHighwaysCode string `json:"office_of_highways_code"`
	SectionRoadNumber    string `json:"section_road_number"`
	Status               bool   `json:"status"`
}

type UpdateRefHris struct {
	RoadNumber           string `json:"road_number"`
	OfficeOfHighwaysCode string `json:"office_of_highways_code"`
	SectionRoadNumber    string `json:"section_road_number"`
	Status               bool   `json:"status"`
}

type FilterHsms struct {
	AssetName string `json:"asset_name"`
}
