package models

import "time"

type SettingRoadUserCostLossValue struct {
	Id        int       `json:"id"`
	Params    string    `json:"params"`
	IsLatest  bool      `json:"is_latest"`
	IsDeleted bool      `json:"is_deleted"`
	UpdatedBy int       `json:"updated_by"`
	CreatedBy int       `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type SettingRoadUserCostChanceOfAccident struct {
	Id          int       `json:"id"`
	RoadGroupId int       `json:"road_group_id"`
	Params      string    `json:"params"`
	IsLatest    bool      `json:"is_latest"`
	IsDeleted   bool      `json:"is_deleted"`
	UpdatedBy   int       `json:"updated_by"`
	CreatedBy   int       `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type SettingRoadUserCost struct {
	Id                      int       `json:"id"`
	DefaultData             string    `json:"default_data"`
	Driving                 string    `json:"driving"`
	EngineSpeed             string    `json:"engine_speed"`
	FuelConsumption         string    `json:"fuel_consumption"`
	LubricantConsumption    string    `json:"lubricant_consumption"`
	WasteOfConsumption      string    `json:"waste_of_consumption"`
	Maintenance             string    `json:"maintenance"`
	TravelTime              string    `json:"travel_time"`
	VehicleSpeedCalculation string    `json:"vehicle_speed_calculation"`
	TrafficData             string    `json:"traffic_data"`
	IsLatest                bool      `json:"is_latest"`
	IsDeleted               bool      `json:"is_deleted"`
	UpdatedBy               int       `json:"updated_by"`
	CreatedBy               int       `json:"created_by"`
	UpdatedAt               time.Time `json:"updated_at"`
	CreatedAt               time.Time `json:"created_at"`
}

type SettingRoadUserCostParams struct {
	Id        int       `json:"id"`
	Params    string    `json:"params"`
	IsLatest  bool      `json:"is_latest"`
	IsDeleted bool      `json:"is_deleted"`
	UpdatedBy int       `json:"updated_by"`
	CreatedBy int       `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type SettingRoadUserCostMergeParams struct {
	RucRoadUserCost     SettingRoadUserCost                   `json:"ruc_road_user_cost"`
	AccLossValue        SettingRoadUserCostLossValue          `json:"acc_loss_value"`
	AccChanceOfAccident []SettingRoadUserCostChanceOfAccident `json:"acc_chance_of_accident"`
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

type MergeRoadUserCost struct {
	Ruc MergeRoadUserCostRuc `json:"ruc"`
	Acc MergeRoadUserCostAcc `json:"acc"`
}

type MergeRoadUserCostRuc struct {
	DefaultData             RoadUserCostRusDefaultData             `json:"default_data"`
	Driving                 RoadUserCostRusDriving                 `json:"driving"`
	EngineSpeed             RoadUserCostRusEngineSpeed             `json:"engine_speed"`
	FuelConsumption         RoadUserCostRusFuelConsumption         `json:"fuel_consumption"`
	LubricantConsumption    RoadUserCostRusLubricantConsumption    `json:"lubricant_consumption"`
	WasteOfConsumption      RoadUserCostRusWasteOfConsumption      `json:"waste_of_consumption"`
	Maintenance             RoadUserCostRusMaintenance             `json:"maintenance"`
	TravelTime              RoadUserCostRusTravelTime              `json:"travel_time"`
	VehicleSpeedCalculation RoadUserCostRusVehicleSpeedCalculation `json:"vehicle_speed_calculation"`
	TrafficData             RoadUserCostRusTrafficData             `json:"traffic_data"`
}
type MergeRoadUserCostAcc struct {
	LossValue        RoadUserCostAccLossValue          `json:"loss_value"`
	ChanceOfAccident []RoadUserCostAccChanceOfAccident `json:"chance_of_accident"`
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

func (b *SettingRoadUserCostLossValue) TableName() string {
	return "setting_road_user_cost_acc_loss_value_params"
}

func (b *SettingRoadUserCostChanceOfAccident) TableName() string {
	return "setting_road_user_cost_acc_chance_of_accident_params"
}

func (b *SettingRoadUserCost) TableName() string {
	return "setting_road_user_cost_ruc"
}

func (b *SettingRoadUserCostParams) TableName() string {
	return "setting_road_user_cost_ruc_params"
}
