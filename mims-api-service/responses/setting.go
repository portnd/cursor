package responses

import (
	"gitlab.com/mims-api-service/models"
)

type ConditionListResponse struct {
	Status bool          `json:"status" example:"true"`
	Data   ConditionList `json:"data"`
}

type ConditionList struct {
	// Owner     models.RefOwner `json:"owner" extensions:"x-order=0"`
	ID                  int         `json:"id" example:"1" extensions:"x-order=0"`
	Name                string      `json:"name" example:"BEM" extensions:"x-order=1"`
	RefConditionRangeID int         `json:"ref_condition_range_id" example:"BEM" extensions:"x-order=2"`
	ConditionList       []Condition `json:"condition" extensions:"x-order=3"`
}

type OwnerList struct {
	Owner models.RefOwner
}

type RoadLineList struct {
	RefReflectivityRangeID int                 `json:"ref_reflectivity_range_id"`
	OwnerName              string              `json:"owner_name"`
	RoadLine               SurfaceTypeRoadLine `json:"road_line" extensions:"x-order=2"`
}

type RoadLineListInit struct {
	ID                     int                     `json:"id" example:"1" extensions:"x-order=0"`
	RefReflectivityRangeID int                     `json:"ref_reflectivity_range_id"`
	OwnerName              string                  `json:"owner_name"`
	RoadLine               SurfaceTypeRoadLineInit `json:"road_line" extensions:"x-order=2"`
}

type RoadLine struct {
	Grade                models.RefGrade `json:"grade" extensions:"x-order=0"`
	ConditionRange       models.RefConditionRange
	LeftValueWhite       float64 `gorm:"column:left_value_white" json:"left_value_white"`
	LeftConditionWhite   string  `gorm:"column:left_condition_white" json:"left_condition_white"`
	RightValueWhite      float64 `gorm:"column:right_value_white" json:"right_value_white"`
	RightConditionWhite  string  `gorm:"column:right_condition_white" json:"right_condition_white"`
	LeftValueYellow      float64 `gorm:"column:left_value_yellow" json:"left_value_yellow"`
	LeftConditionYellow  string  `gorm:"column:left_condition_yellow" json:"left_condition_yellow"`
	RightValueYellow     float64 `gorm:"column:right_value_yellow" json:"right_value_yellow"`
	RightConditionYellow string  `gorm:"column:right_condition_yellow" json:"right_condition_yellow"`
}

type ConditionValue struct {
	Grade                models.RefGrade `json:"grade" extensions:"x-order=0"`
	LeftValueWhite       float64         `gorm:"column:left_value_white" json:"left_value_white"`
	LeftConditionWhite   string          `gorm:"column:left_condition_white" json:"left_condition_white"`
	RightValueWhite      float64         `gorm:"column:right_value_white" json:"right_value_white"`
	RightConditionWhite  string          `gorm:"column:right_condition_white" json:"right_condition_white"`
	LeftValueYellow      float64         `gorm:"column:left_value_yellow" json:"left_value_yellow"`
	LeftConditionYellow  string          `gorm:"column:left_condition_yellow" json:"left_condition_yellow"`
	RightValueYellow     float64         `gorm:"column:right_value_yellow" json:"right_value_yellow"`
	RightConditionYellow string          `gorm:"column:right_condition_yellow" json:"right_condition_yellow"`
}

type ConditionTypeYellow struct {
	Grade                models.RefGrade `json:"grade" extensions:"x-order=0"`
	LeftValueYellow      float64         `gorm:"column:left_value_yellow" json:"left_value_yellow"`
	LeftConditionYellow  string          `gorm:"column:left_condition_yellow" json:"left_condition_yellow"`
	RightValueYellow     float64         `gorm:"column:right_value_yellow" json:"right_value_yellow"`
	RightConditionYellow string          `gorm:"column:right_condition_yellow" json:"right_condition_yellow"`
}

type ConditionTypeWhite struct {
	Grade               models.RefGrade `json:"grade" extensions:"x-order=0"`
	LeftValueWhite      float64         `gorm:"column:left_value_white" json:"left_value_white"`
	LeftConditionWhite  string          `gorm:"column:left_condition_white" json:"left_condition_white"`
	RightValueWhite     float64         `gorm:"column:right_value_white" json:"right_value_white"`
	RightConditionWhite string          `gorm:"column:right_condition_white" json:"right_condition_white"`
}

type SurfaceTypeCondition struct {
	AC []ConditionTypeAC `json:"ac"`
	CC []ConditionTypeCC `json:"cc"`
}

type SurfaceTypeRoadLine struct {
	Yellow []ConditionTypeYellow `json:"yellow"`
	White  []ConditionTypeWhite  `json:"white"`
}

type SurfaceTypeRoadLineInit struct {
	Yellow []ConditionTypeInit `json:"yellow"`
	White  []ConditionTypeInit `json:"white"`
}

// type RoadLineListNew struct {
// 	ConditionType string              `json:"condition_type"`
// 	SurfaceType   SurfaceTypeRoadLine `json:"surface_type"`
// }

type ConditionListNew struct {
	ConditionType string               `json:"condition_type"`
	SurfaceType   SurfaceTypeCondition `json:"surface_type"`
}

type ConditionListNewInit struct {
	ConditionType string                   `json:"condition_type"`
	SurfaceType   SurfaceTypeConditionInit `json:"surface_type"`
}

type ConditionRespond struct {
	ID                  int                `json:"id"`
	RefConditionRangeID int                `json:"ref_condition_range_id"`
	OwnerName           string             `json:"owner_name"`
	ConditionList       []ConditionListNew `json:"condition_list"`
}

type ConditionRespondInit struct {
	ID                  int                    `json:"id"`
	RefConditionRangeID int                    `json:"ref_condition_range_id"`
	OwnerName           string                 `json:"owner_name"`
	ConditionList       []ConditionListNewInit `json:"condition_list"`
}

type ConditionLisInit struct {
	ConditionType string               `json:"condition_type"`
	SurfaceType   SurfaceTypeCondition `json:"surface_type"`
}

type SurfaceTypeConditionInit struct {
	AC []ConditionTypeInit `json:"ac"`
	CC []ConditionTypeInit `json:"cc"`
}

type ConditionTypeInit struct {
	Grade          models.RefGrade `json:"grade" extensions:"x-order=0"`
	LeftValue      float64         ` json:"left_value"`
	LeftCondition  string          ` json:"left_condition"`
	RightValue     float64         ` json:"right_value"`
	RightCondition string          ` json:"right_condition"`
}

type Condition struct {
	IFI []ConditionType `json:"ifi" extensions:"x-order=0"`
	IRI []ConditionType `json:"iri" extensions:"x-order=1"`
	MPD []ConditionType `json:"mpd" extensions:"x-order=2"`
	RUT []ConditionType `json:"rut" extensions:"x-order=3"`
}

type ConditionType struct {
	Grade            models.RefGrade `json:"grade" extensions:"x-order=0"`
	LeftValueAC      float64         `gorm:"column:left_value_ac" json:"left_value_ac"`
	LeftConditionAC  string          `gorm:"column:left_condition_ac" json:"left_condition_ac"`
	RightValueAC     float64         `gorm:"column:right_value_ac" json:"right_value_ac"`
	RightConditionAC string          `gorm:"column:right_condition_ac" json:"right_condition_ac"`
	LeftValueCC      float64         `gorm:"column:left_value_cc" json:"left_value_cc"`
	LeftConditionCC  string          `gorm:"column:left_condition_cc" json:"left_condition_cc"`
	RightValueCC     float64         `gorm:"column:right_value_cc" json:"right_value_cc"`
	RightConditionCC string          `gorm:"column:right_condition_cc" json:"right_condition_cc"`
}

type InterventionCriteriaSurface struct {
	ID                               int                               `json:"id"`
	Name                             string                            `json:"name"`
	InterventionCriteriaMaintenances []InterventionCriteriaMaintenance `json:"intervention_criterias"`
}

type ConditionTypeAC struct {
	Grade            models.RefGrade `json:"grade" extensions:"x-order=0"`
	LeftValueAC      float64         `gorm:"column:left_value_ac" json:"left_value_ac"`
	LeftConditionAC  string          `gorm:"column:left_condition_ac" json:"left_condition_ac"`
	RightValueAC     float64         `gorm:"column:right_value_ac" json:"right_value_ac"`
	RightConditionAC string          `gorm:"column:right_condition_ac" json:"right_condition_ac"`
}

type ConditionTypeCC struct {
	Grade            models.RefGrade `json:"grade" extensions:"x-order=0"`
	LeftValueCC      float64         `gorm:"column:left_value_cc" json:"left_value_cc"`
	LeftConditionCC  string          `gorm:"column:left_condition_cc" json:"left_condition_cc"`
	RightValueCC     float64         `gorm:"column:right_value_cc" json:"right_value_cc"`
	RightConditionCC string          `gorm:"column:right_condition_cc" json:"right_condition_cc"`
}
type AssetTableResponse struct {
	ID              int    `json:"id" extensions:"x-order=0"`
	TableName       string `json:"table_name" extensions:"x-order=1"`
	TableLabel      string `json:"table_label" extensions:"x-order=2"`
	AssetGroup      string `json:"asset_group" extensions:"x-order=3"`
	ResponsibleDept string `json:"responsible_dept" extensions:"x-order=4"`
	CanDelete       bool   `json:"can_delete" extensions:"x-order=5"`
}

type AssetTableDetailResponse struct {
	AssetID       int             `json:"asset_id" extensions:"x-order=0"`
	TableName     string          `json:"table_name" extensions:"x-order=1"`
	TableLabel    string          `json:"table_label" extensions:"x-order=2"`
	GeomType      int             `json:"geom_type" extensions:"x-order=3"`
	IconFilePath  string          `json:"icon_filepath" extensions:"x-order=4"`
	LineColorCode string          `json:"line_color_code" extensions:"x-order=5"`
	AssetGroup    models.RefAsset `json:"asset_group" extensions:"x-order=6"`
	// Approver      []models.RefDepartment `json:"approver" extensions:"x-order=7"`
	// Viewer        []models.RefDepartment `json:"viewer" extensions:"x-order=8"`
	Columns []Columns `json:"columns" extensions:"x-order=9"`
}

type Columns struct {
	ColumnID        int    `json:"column_id" extensions:"x-order=0"`
	ColumnName      string `json:"column_name" extensions:"x-order=1"`
	TableNameRef    string `json:"table_name_ref" extensions:"x-order=2"`
	ColumnDataType  string `json:"column_data_type" extensions:"x-order=3"`
	ComponentTitle  string `json:"component_title" extensions:"x-order=4"`
	ComponentType   string `json:"component_type" extensions:"x-order=5"`
	IsRequired      bool   `json:"is_required" extensions:"x-order=6"`
	IsVisibleView   bool   `json:"is_visible_view" extensions:"x-order=7"`
	IsVisibleEdit   bool   `json:"is_visible_edit" extensions:"x-order=8"`
	IsMandatory     bool   `json:"is_mandatory" extensions:"x-order=9"`
	IsVisibleReport bool   `json:"is_visible_report" extensions:"x-order=10"`
}

type Budget struct {
	Id        int            `json:"id"`
	Name      string         `json:"name"`
	CanDelete bool           `json:"can_delete"`
	Budget    []BudgetMethod `json:"budget" `
}

type BudgetMethod struct {
	Id           int      `json:"id"`
	CostPerUnit  *float64 `json:"cost_per_unit"`
	MethodName   string   `json:"method_name"`
	IsShowMethod bool     `json:"is_show_method"`
}

type BudgetList struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CanDelete bool   `json:"can_delete"`
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

type CreateInterventionCriteria struct {
	MaintenanceCostPerUnit   float64                         `json:"maintenance_cost_per_unit"`
	MaintenanceDescription   string                          `json:"maintenance_description"`
	MaintenanceScraping      float64                         `json:"maintenance_scraping"`
	MaintenanceSequence      int                             `json:"maintenance_sequence"`
	MaintenanceStandardName  string                          `json:"maintenance_standard_name"`
	MaintenanceSurfaceTypeId int                             `json:"maintenance_surface_type_id"`
	MaintenanceThickness     float64                         `json:"maintenance_thickness"`
	MaintenanceType          string                          `json:"maintenance_type"`
	MaintenanceCondition     []InterventionCriteriaCindition `json:"maintenance_condition"`
}

type InterventionCriteriaSequenceCriteriaMethod struct {
	Concrete []InterventionCriteriaSequence `json:"concrete"`
	Asphalt  []InterventionCriteriaSequence `json:"asphalt"`
}

type InterventionCriteriaSequence struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type InterventionCriteriaParams struct {
	Asphalt  []InterventionCriteriaSurfaceParams `json:"asphalt"`
	Concrete []InterventionCriteriaSurfaceParams `json:"concrete"`
}

type InterventionCriteriaSurfaceParams struct {
	ID                               int                               `json:"id"`
	Name                             string                            `json:"name"`
	InterventionCriteriaMaintenances []InterventionCriteriaMaintenance `json:"intervenrion_creteria_maintenances"`
}

type InterventionCriteriaSurfaceAspaltParams struct {
	OlOverlay        []InterventionCriteriaMaintenance `json:"OL-Overlay"`
	SsSlurrySeal     []InterventionCriteriaMaintenance `json:"SS-SlurrySeal"`
	MolMillOverlay   []InterventionCriteriaMaintenance `json:"M&OL-Mill&Overlay"`
	RclRecycling     []InterventionCriteriaMaintenance `json:"RCL-Recycling"`
	RcReconstruction []InterventionCriteriaMaintenance `json:"Rc-Reconstruction"`
}

type InterventionCriteriaSurfaceConcreteParams struct {
	Fdr  []InterventionCriteriaMaintenance `json:"FDR"`
	Bco  []InterventionCriteriaMaintenance `json:"BCO"`
	Mol  []InterventionCriteriaMaintenance `json:"M&OL"`
	Seal []InterventionCriteriaMaintenance `json:"Seal"`
}

type InterventionCriteria struct {
	Asphalt  []InterventionCriteriaSurface `json:"asphalt"`
	Concrete []InterventionCriteriaSurface `json:"concrete"`
}

type InterventionCriteriaSurfaceAspalt struct {
	OlOverlay        []InterventionCriteriaMaintenance `json:"ol_overlay"`
	SsSlurrySeal     []InterventionCriteriaMaintenance `json:"ss_slurry_seal"`
	MolMillOverlay   []InterventionCriteriaMaintenance `json:"mol_mill_overlay"`
	RclRecycling     []InterventionCriteriaMaintenance `json:"rcl_recycling"`
	RcReconstruction []InterventionCriteriaMaintenance `json:"rc_reconstruction"`
}

type InterventionCriteriaSurfaceConcrete struct {
	Fdr  []InterventionCriteriaMaintenance `json:"fdr"`
	Bco  []InterventionCriteriaMaintenance `json:"bco"`
	Mol  []InterventionCriteriaMaintenance `json:"mol"`
	Seal []InterventionCriteriaMaintenance `json:"seal"`
}

type InterventionCriteriaMaintenance struct {
	Id                       int                             `json:"id"`
	MaintenanceMethod        string                          `json:"maintenance_method"`
	MaintenanceCostPerUnit   float64                         `json:"maintenance_cost_per_unit"`
	MaintenanceDescription   string                          `json:"maintenance_description"`
	MaintenanceScraping      float64                         `json:"maintenance_scraping"`
	MaintenanceSequence      int                             `json:"maintenance_sequence"`
	MaintenanceStandardName  string                          `json:"maintenance_standard_name"`
	MaintenanceSurfaceTypeId int                             `json:"maintenance_surface_type_id"`
	MaintenanceThickness     float64                         `json:"maintenance_thickness"`
	MaintenanceType          string                          `json:"maintenance_type"`
	MaintenanceCondition     []InterventionCriteriaCindition `json:"maintenance_condition"`
}

type InterventionCriteriaCindition struct {
	Id                  int     `json:"id"`
	ConditionCriterion  string  `json:"condition_criterion"`
	ConditionLink       string  `json:"condition_link"`
	ConditionOperation1 string  `json:"condition_operation_1"`
	ConditionOperation2 string  `json:"condition_operation_2"`
	ConditionValue1     float64 `json:"condition_value_1"`
	ConditionValue2     float64 `json:"condition_value_2"`
}

type AadtGrowthRate struct {
	RoadGroupId   int     `json:"road_group_id"`
	R             float64 `json:"r"`
	Code          string  `json:"code"`
	RoadGroupName string  `json:"road_group_name"`
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

type RoadGroupWithVolumeAadt struct {
	RoadGroupId   int                     `json:"road_group_id"`
	RoadGroupName string                  `json:"road_group_name"`
	VolumeAadt    VolumeAadtWithRoadGroup `json:"volume_aadt"`
}

type VolumeAadtWithRoadGroup struct {
	Veh1      int                              `json:"veh1"`
	Veh2      int                              `json:"Veh2"`
	Veh3      int                              `json:"Veh3"`
	Veh4      int                              `json:"Veh4"`
	Calculate VolumeAadtWithRoadGroupCalculate `json:"calculate"`
}

type VolumeAadtWithRoadGroupCalculate struct {
	FourWheelTotal          int     `json:"four_wheel_total"`
	SixToTenWheelTotal      int     `json:"six_to_ten_wheel_total"`
	TenWheelTotal           int     `json:"ten_wheel_total"`
	SixToTenWheelPercentage float64 `json:"six_to_ten_wheel_percentage"`
	TenWheelPercentage      float64 `json:"ten_wheel_percentage"`
}

type AadtParameter struct {
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

type CalculateAadtParameterCar struct {
	VehicleVolume float64 `json:"vehicle_volume"`
}

type CalculateAadtParameterTruck struct {
	VehicleVolume   float64 `json:"vehicle_volume"`
	PercentageTruck float64 `json:"percentage_truck"`
	LoadEquivalent  float64 `json:"load_equivalent"`
	TruckFactor     float64 `json:"truck_factor"`
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
	RoadGroupId int     `json:"road_group_id"`
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
	RoadGroupId int     `json:"road_group_id"`
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

type InterventionCriteriaReportData struct {
	UpdatedBy            string `json:"updated_by"`
	UpdatedAt            string `json:"updated_at"`
	User                 string `json:"user"`
	PrintDate            string `json:"print_date"`
	InterventionCriteria []InterventionCriteriaReportRes
}
type InterventionCriteriaReportRes struct {
	Seq          int                                       `json:"seq"`
	StandardName string                                    `json:"standard_name"`
	Description  string                                    `json:"description"`
	Conditions   []InterventionCriteriaConditionsReportRes `json:"conditions"`
}

type InterventionCriteriaReport2Data struct {
	UpdatedBy            string `json:"updated_by"`
	UpdatedAt            string `json:"updated_at"`
	User                 string `json:"user"`
	PrintDate            string `json:"print_date"`
	InterventionCriteria []InterventionCriteriaReport2Res
}

type InterventionCriteriaReport2Res struct {
	No                     int     `json:"no"`
	Seq                    float64 `json:"seq"`
	StandardName           string  `json:"standard_name"`
	MaintenanceCostPerUnit string  `json:"maintenance_cost_per_unit"`
}
type AnalyzeReport1 struct {
	Data []InterventionCriteriaReportRes `json:"data"`
}

type AnalyzeReport2 struct {
	Data []InterventionCriteriaReport2Res `json:"data"`
}
type InterventionCriteriaConditionsReportRes struct {
	ConditionCriterion  string `json:"condition_criterion"`
	ConditionLink       string `json:"condition_link"`
	ConditionOperation1 string `json:"condition_operation_1"`
	ConditionOperation2 string `json:"condition_operation_2"`
	ConditionValue1     string `json:"condition_value_1"`
	ConditionValue2     string `json:"condition_value_2"`
}

type RefHris struct {
	Id                   int    `json:"id"`
	RoadNumber           string `json:"road_number"`
	OfficeOfHighwaysCode string `json:"office_of_highways_code"`
	SectionRoadNumber    string `json:"section_road_number"`
	Status               bool   `json:"status"`
}

type RefHrisPreview struct {
	RoadGroup   []HrisRoadList    `json:"road_group"`
	RoadSection []HrisSectionGeom `json:"road_section"`
}

type HrisSectionGeom struct {
	RoadGroupNumber    string `json:"road_group_number"`
	SectionRoadNumber  string `json:"section_road_number"`
	SectionRoadThName  string `json:"section_road_th_name"`
	SectionRoadEngName string `json:"section_road_eng_name"`
	KmStart            string `json:"km_start"`
	KmEnd              string `json:"km_end"`
}

type HrisRoadList struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}

type HsmsAll struct {
	Id               int    `json:"id"`
	Type             string `json:"type"`
	AssetName        string `json:"asset_name"`
	RoadGroupName    string `json:"road_group_name"`
	RoadName         string `json:"road_name"`
	Km               string `json:"km"`
	KmRange          string `json:"km_range"`
	LocationName     string `json:"location_name"`
	LocationTypeName string `json:"location_type_name"`
	DepotName        string `json:"depot_name"`
}
