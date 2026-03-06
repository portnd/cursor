package models

import (
	"time"
)

type RoadRetroReflectivity struct {
	ID            int       `gorm:"column:id"` // assuming `db` is the struct tag used by your ORM
	RoadID        int       `gorm:"column:road_id"`
	LineNo        int       `gorm:"column:line_no"`
	Year          int       `gorm:"column:year"`
	KmStart       float64   `gorm:"column:km_start"`
	KmEnd         float64   `gorm:"column:km_end"`
	SurveyedDate  time.Time `gorm:"column:surveyed_date"`
	CreatedBy     int       `gorm:"column:created_by"`
	CreatedDate   time.Time `gorm:"column:created_date"`
	UpdatedBy     int       `gorm:"column:updated_by"`
	UpdatedDate   time.Time `gorm:"column:updated_date"`
	RetroMin      *float64  `gorm:"column:retro_min"`      // Nullable field
	RetroMax      *float64  `gorm:"column:retro_max"`      // Nullable field
	RetroAvg      *float64  `gorm:"column:retro_avg"`      // Nullable field
	Remarks       string    `gorm:"column:remarks"`        // Nullable field
	InputFilePath string    `gorm:"column:input_filepath"` // Nullable field
	Revision      int       `gorm:"column:revision"`
	Status        string    `gorm:"column:status"`    // 'bpchar' is usually a fixed length string
	IDParent      int       `gorm:"column:id_parent"` // Nullable field
}

type RoadRetroReflectivityRange struct {
	ID                      int      `gorm:"column:id"`
	RoadRetroReflectivityID int      `gorm:"column:road_retro_reflectivity_id"`
	KmStart                 float64  `gorm:"column:km_start"`
	KmEnd                   float64  `gorm:"column:km_end"`
	RetroMin                *float64 `gorm:"column:retro_min"`
	RetroMax                *float64 `gorm:"column:retro_max"`
	RetroAvg                *float64 `gorm:"column:retro_avg"`
	TheGeom                 string   `gorm:"column:the_geom"`
	RefStripeColorID        int      `gorm:"column:ref_stripe_color_id"`
	RefStripeTypeID         int      `gorm:"column:ref_stripe_type_id"`
}

type RoadRetroReflectivityM struct {
	ID                           int      `gorm:"column:id"`
	RoadRetroReflectivityRangeID int      `gorm:"column:road_retro_reflectivity_range_id"`
	KmStart                      float64  `gorm:"column:km_start"`
	KmEnd                        float64  `gorm:"column:km_end"`
	RetroMin                     *float64 `gorm:"column:retro_min"`
	RetroMax                     *float64 `gorm:"column:retro_max"`
	RetroAvg                     *float64 `gorm:"column:retro_avg"`
	TheGeom                      string   `gorm:"column:the_geom"`
	RefStripeColorID             int      `gorm:"column:ref_stripe_color_id"`
	RefStripeTypeID              int      `gorm:"column:ref_stripe_type_id"`
}

func (b *RoadRetroReflectivity) TableName() string {
	return "road_retro_reflectivity"
}

func (b *RoadRetroReflectivityRange) TableName() string {
	return "road_retro_reflectivity_range"
}

func (b *RoadRetroReflectivityM) TableName() string {
	return "road_retro_reflectivity_m"
}

type RoadRetroReflectivityCSV struct {
	RoadID     int      `json:"road_id"`
	RoadCode   string   `json:"road_code"`
	Name       string   `json:"name"`
	KMStart    float64  `json:"km_start"`
	KMEnd      float64  `json:"km_end"`
	RetroMin   *float64 `json:"retro_min"`
	RetroMax   *float64 `json:"retro_max"`
	RetroAvg   *float64 `json:"retro_avg"`
	Color      string   `json:"color"`
	StripeType string   `json:"stripe_type"`
}

type RoadRetroReflectivityData struct {
	RetroMinAverage      *float64
	RetroMaxAverage      *float64
	RetroAvgAverage      *float64
	DividerCountRetroMin *float64
	DividerCountRetroMax *float64
	DividerCountRetroAvg *float64
	RetroMin             *float64
	RetroMax             *float64
	RetroAvg             *float64
	Rut100m              *float64
	TotalM               float64
}

type RoadRetroReflectivityPreload struct {
	RoadRetroReflectivity
	RoadRetroReflectivityRanges []RoadRetroReflectivityRangePreload ` gorm:"ForeignKey:RoadRetroReflectivityID;references:ID"`
}

type RoadRetroReflectivityRangePreload struct {
	RoadRetroReflectivityRange

	RefStripeColor          RefStripeColor                  `gorm:"ForeignKey:RefStripeColorID;references:ID"`
	RefStripeType           RefStripeType                   `gorm:"ForeignKey:RefStripeTypeID;references:ID"`
	RoadRetroReflectivityMs []RoadRetroReflectivityMPreload ` gorm:"ForeignKey:RoadRetroReflectivityRangeID;references:ID"`
}

type RoadRetroReflectivityMPreload struct {
	ID                           int            `gorm:"column:id"`
	RoadRetroReflectivityRangeID int            `gorm:"column:road_retro_reflectivity_range_id"`
	KmStart                      float64        `gorm:"column:km_start"`
	KmEnd                        float64        `gorm:"column:km_end"`
	RetroMin                     *float64       `gorm:"column:retro_min"`
	RetroMax                     *float64       `gorm:"column:retro_max"`
	RetroAvg                     *float64       `gorm:"column:retro_avg"`
	TheGeom                      string         `gorm:"column:the_geom"`
	RefStripeColorID             int            `gorm:"column:ref_stripe_color_id"`
	RefStripeColor               RefStripeColor `gorm:"ForeignKey:RefStripeColorID;references:ID"`
	RefStripeTypeID              int            `gorm:"column:ref_stripe_type_id"`
	RefStripeType                RefStripeType  `gorm:"ForeignKey:RefStripeTypeID;references:ID"`
}

type RoadRetroReflectivityRangeDashboard struct {
	RoadRetroReflectivityRange
	RetroRangeTheGeom       []byte                            `gorm:"column:retro_range_the_geom"`
	RefStripeColor          RefStripeColor                    `gorm:"ForeignKey:RefStripeColorID;references:ID"`
	RefStripeType           RefStripeType                     `gorm:"ForeignKey:RefStripeTypeID;references:ID"`
	RoadRetroReflectivityMs []RoadRetroReflectivityMDashboard ` gorm:"ForeignKey:RoadRetroReflectivityRangeID;references:ID"`
}

type RoadRetroReflectivityMDashboard struct {
	ID                           int            `gorm:"column:id"`
	RoadRetroReflectivityRangeID int            `gorm:"column:road_retro_reflectivity_range_id"`
	KmStart                      float64        `gorm:"column:km_start"`
	KmEnd                        float64        `gorm:"column:km_end"`
	RetroMin                     *float64       `gorm:"column:retro_min"`
	RetroMax                     *float64       `gorm:"column:retro_max"`
	RetroAvg                     *float64       `gorm:"column:retro_avg"`
	RetroMTheGeom                []byte         `gorm:"column:retro_m_the_geom"`
	RefStripeColorID             int            `gorm:"column:ref_stripe_color_id"`
	RefStripeColor               RefStripeColor `gorm:"ForeignKey:RefStripeColorID;references:ID"`
	RefStripeTypeID              int            `gorm:"column:ref_stripe_type_id"`
	RefStripeType                RefStripeType  `gorm:"ForeignKey:RefStripeTypeID;references:ID"`
}

func (b *RoadRetroReflectivityPreload) TableName() string {
	return "road_retro_reflectivity"
}

func (b *RoadRetroReflectivityRangePreload) TableName() string {
	return "road_retro_reflectivity_range"
}

func (b *RoadRetroReflectivityMPreload) TableName() string {
	return "road_retro_reflectivity_m"
}

func (b *RoadRetroReflectivityRangeDashboard) TableName() string {
	return "road_retro_reflectivity_range"
}

func (b *RoadRetroReflectivityMDashboard) TableName() string {
	return "road_retro_reflectivity_m"
}
