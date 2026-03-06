package models

import "time"

type SettingAADTParameter struct {
	ID                            int64     `json:"id"`
	RoadGroupID                   int32     `json:"road_group_id"`
	Elane                         float64   `json:"elane"`
	FourWheelAxleNumber           int32     `json:"four_wheel_axle_number"`
	FourWheelVehicleVolume        float64   `json:"four_wheel_vehicle_volume"`
	SixWheelAxleNumberID          int32     `json:"six_wheel_axle_number_id"`
	SixWheelVehicleVolume         float64   `json:"six_wheel_vehicle_volume"`
	SixWheelPercentageTruck       float64   `json:"six_wheel_percentage_truck"`
	SixWheelFactorResult          float64   `json:"six_wheel_factor_result"`
	TenWheelAxleNumberID          int32     `json:"ten_wheel_axle_number_id"`
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
	CreatedBy                     int32     `json:"created_by"`
	UpdatedBy                     int32     `json:"updated_by"`
	CreatedAt                     time.Time `json:"created_at"`
	UpdatedAt                     time.Time `json:"updated_at"`
}

type AadtParameterData struct {
	Elane                         float64 `json:"elane" gorm:"column:elane"`
	TrackFactorOpenInput          bool    `json:"track_factor_open_input" gorm:"column:is_truck_factor"`
	Car4Axle                      float64 `json:"car_4_axle" gorm:"column:four_wheel_axle_number"`
	Truck6Axle                    int     `json:"truck_6_axle" gorm:"column:six_wheel_axle_number_id"`
	Truck6LoadEquivalent          float64 `json:"truck_6_load_equivalent" gorm:"column:load_equivalent"`
	Truck6TruckFactorInput        float64 `json:"truck_6_truck_factor_input" gorm:"column:six_wheel_factor_result"`
	Truck10Axle                   int     `json:"truck_10_axle" gorm:"column:ten_wheel_axle_number_id"`
	Truck10LoadEquivalent         float64 `json:"truck_10_load_equivalent" gorm:"column:load_equivalent"`
	Truck10TruckFactorInput       float64 `json:"truck_10_truck_factor_input" gorm:"column:ten_wheel_factor_result"`
	SpeedAverage                  float64 `json:"speed_average" gorm:"column:speed_average"`
	SpeedHeavyTruck               float64 `json:"speed_heavy_truck" gorm:"column:speed_heavy_truck"`
	LaneDistributionFactor        float64 `json:"lane_distribution_factor" gorm:"column:lane_distribution_factor"`
	DirectionalDistributionFactor float64 `json:"directional_distribution_factor" gorm:"column:directional_distribution_factor"`
}

func (ract *SettingAADTParameter) TableName() string {
	return "setting_aadt_parameter"
}
