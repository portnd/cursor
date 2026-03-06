package models

import "time"

type AadtGrowthRate struct {
	Id          int       `json:"id"`
	RoadGroupId int       `json:"road_group_id"`
	R           float64   `json:"r"`
	IsLatest    bool      `json:"is_latest"`
	IsDeleted   bool      `json:"is_deleted"`
	UpdatedBy   int       `json:"updated_by"`
	CreatedBy   int       `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetAadtGrowthRate struct {
	RoadGroupId   int     `json:"road_group_id"`
	R             float64 `json:"r"`
	Number        string  `json:"number"`
	RoadGroupName string  `json:"road_group_name"`
}

type AadtPercentageVehicleTypeParams struct {
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

type AadtParameter struct {
	Id                            int       `json:"id"`
	RoadGroupId                   int       `json:"road_group_id" validate:"nonzero"`
	Elane                         float64   `json:"elane"`
	FourWheelAxleNumber           int       `json:"four_wheel_axle_number"`
	FourWheelVehicleVolume        float64   `json:"four_wheel_vehicle_volume"`
	SixWheelAxleNumberId          int       `json:"six_wheel_axle_number_id"`
	SixWheelVehicleVolume         float64   `json:"six_wheel_vehicle_volume"`
	SixWheelPercentageTruck       float64   `json:"six_wheel_percentage_truck"`
	SixWheelFactorResult          float64   `json:"six_wheel_factor_result"`
	TenWheelAxleNumberId          int       `json:"ten_wheel_axle_number_id"`
	TenWheelVehicleVolume         float64   `json:"ten_wheel_vehicle_volume"`
	TenWheelPercentageTruck       float64   `json:"ten_wheel_percentage_truck"`
	TenWheelFactorResult          float64   `json:"ten_wheel_factor_result"`
	IsTruckFactor                 bool      `json:"is_truck_factor"`
	SpeedAverage                  float64   `json:"speed_average"`
	SpeedHeavyTruck               float64   `json:"speed_heavy_truck"`
	LaneDistributionFactor        float64   `json:"lane_distribution_factor"`
	DirectionalDistributionFactor float64   `json:"directional_distribution_factor"`
	IsLatest                      bool      `json:"is_latest"`
	IsDeleted                     bool      `json:"is_deleted"`
	UpdatedBy                     int       `json:"updated_by"`
	CreatedBy                     int       `json:"created_by"`
	UpdatedAt                     time.Time `json:"updated_at"`
	CreatedAt                     time.Time `json:"created_at"`
}

type GetAadtParameter struct {
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

type GetAadtParams struct {
	GetAadtGrowthRate            []GetAadtGrowthRate         `json:"aadt_growth_rate"`
	GetAadtPercentageVehicleType []AadtPercentageVehicleType `json:"aadt_percentage_vehicleType"`
	GetAadtParameter             []GetAadtParameter          `json:"aadt_parameter"`
}

type AadtParams struct {
	Id        int       `json:"id"`
	Params    string    `json:"params"`
	CreatedBy int       `json:"-"`
	CreatedAt time.Time `json:"-"`
	IsLatest  bool      `json:"is_latest"`
}

func (b *AadtGrowthRate) TableName() string {
	return "setting_aadt_growth_rate"
}

func (b *AadtPercentageVehicleTypeParams) TableName() string {
	return "setting_aadt_percentage_vehicle_type"
}
