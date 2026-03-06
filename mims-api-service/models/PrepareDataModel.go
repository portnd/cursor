package models

import "time"

// Todo ...
type PrepareData struct {
	ID                             int         `json:"id"`
	MaintenanceAnalysisID          int         `json:"maintenance_analysis_id"`
	IsSelected                     bool        `json:"is_selected"`
	GroupName                      string      `json:"group_name"`
	RoadName                       string      `json:"road_name"`
	RoadGroupID                    int         `json:"road_group_id"`
	RoadID                         int         `json:"road_id"`
	LaneNo                         int         `json:"lane_no"`
	LaneKmStart                    float64     `json:"lane_km_start"`
	LaneKmEnd                      float64     `json:"lane_km_end"`
	LaneLength                     float64     `json:"lane_length"`
	LaneWidth                      float64     `json:"lane_width"`
	KmStart                        float64     `json:"km_start"`
	KmEnd                          float64     `json:"km_end"`
	Length                         float64     `json:"length"`
	Area                           float64     `json:"area"`
	LastType                       int         `json:"last_type"`
	Type                           int         `json:"type"`
	AnalystYear                    int         `json:"analyst_year"`
	YearRoadBegin                  int         `json:"year_road_begin"`
	YearLastOverlay                int         `json:"year_last_overlay"`
	YearLastSeal                   int         `json:"year_last_seal"`
	YearLastMolRcl                 int         `json:"year_last_mol_rcl"`
	YearLastReconstruction         int         `json:"year_last_reconstruction"`
	Age                            int         `json:"age"`
	Rut                            float64     `json:"rut"`
	Iri                            float64     `json:"iri"`
	Ifi                            float64     `json:"ifi"`
	NumberOfPothole                float64     `json:"number_of_pothole"`
	AreaAcIcrack                   float64     `json:"area_ac_icrack"`
	PercentAcIcrack                float64     `json:"percent_ac_icrack"`
	AreaAcUcrack                   float64     `json:"area_ac_ucrack"`
	PercentAcUcrack                float64     `json:"percent_ac_ucrack"`
	PercentAcRavelling             float64     `json:"percent_ac_ravelling"`
	CcTransverseCrack              float64     `json:"cc_transverse_crack"`
	CcFaulting                     float64     `json:"cc_faulting"`
	CcSpalling                     float64     `json:"cc_spalling"`
	CurrentSurfaceID               int         `json:"current_surface_id"`
	CurrentSurfaceName             string      `json:"current_surface_name"`
	CurrentSurfaceType             string      `json:"current_surface_type"`
	CurrentSurfaceSurfaceGroup     string      `json:"current_surface_surface_group"`
	CurrentSurfaceLayerCoefficient float64     `json:"current_surface_layer_coefficient"`
	CurrentSurfaceDrainage         float64     `json:"current_surface_drainage"`
	CurrentSurfaceA                float64     `json:"current_surface_a"`
	CurrentSurfaceB                float64     `json:"current_surface_b"`
	CurrentSurfaceCBase            float64     `json:"current_surface_c_base"`
	CurrentSurfaceCExp             float64     `json:"current_surface_c_exp"`
	CurrentSurfaceCRT              float64     `json:"current_surface_crt"`
	CurrentSurfaceRRF              float64     `json:"current_surface_rrf"`
	Hsold                          float64     `json:"hsold"`
	Hsnew                          float64     `json:"hsnew"`
	SNPSurface                     float64     `json:"snp_surface"`
	SNPBase                        float64     `json:"snp_base"`
	SNPSubbase                     float64     `json:"snp_subbase"`
	SNP                            float64     `json:"snp"`
	AADT                           float64     `json:"aadt"`
	TruckFactor                    float64     `json:"truck_factor"`
	ESAL                           float64     `json:"esal"`
	YAX                            float64     `json:"yax"`
	Data                           interface{} `json:"data" gorm:"type:text"`
	TheGeom                        string      `json:"the_geom"`
	CreatedBy                      int         `json:"created_by"` // รหัสผู้ใช้งานที่สร้างข้อมูล
	UpdatedBy                      int         `json:"updated_by"` // รหัสผู้ใช้งานที่อัพเดตข้อมูล
	CreatedAt                      time.Time   `json:"created_at"` // วันที่ที่สร้างข้อมูล
	UpdatedAt                      time.Time   `json:"updated_at"` // วันที่ที่อัพเดตข้อมูล
}

type PrepareDataForCopy struct {
	ID                             int       `json:"id"`
	MaintenanceAnalysisID          int       `json:"maintenance_analysis_id"`
	IsSelected                     bool      `json:"is_selected"`
	GroupName                      string    `json:"group_name"`
	RoadName                       string    `json:"road_name"`
	RoadGroupID                    int       `json:"road_group_id"`
	RoadID                         int       `json:"road_id"`
	LaneNo                         int       `json:"lane_no"`
	LaneKmStart                    float64   `json:"lane_km_start"`
	LaneKmEnd                      float64   `json:"lane_km_end"`
	LaneLength                     float64   `json:"lane_length"`
	LaneWidth                      float64   `json:"lane_width"`
	KmStart                        float64   `json:"km_start"`
	KmEnd                          float64   `json:"km_end"`
	Length                         float64   `json:"length"`
	Area                           float64   `json:"area"`
	LastType                       int       `json:"last_type"`
	Type                           int       `json:"type"`
	AnalystYear                    int       `json:"analyst_year"`
	YearRoadBegin                  int       `json:"year_road_begin"`
	YearLastOverlay                int       `json:"year_last_overlay"`
	YearLastSeal                   int       `json:"year_last_seal"`
	YearLastMolRcl                 int       `json:"year_last_mol_rcl"`
	YearLastReconstruction         int       `json:"year_last_reconstruction"`
	Age                            int       `json:"age"`
	Rut                            float64   `json:"rut"`
	Iri                            float64   `json:"iri"`
	Ifi                            float64   `json:"ifi"`
	NumberOfPothole                float64   `json:"number_of_pothole"`
	AreaAcIcrack                   float64   `json:"area_ac_icrack"`
	PercentAcIcrack                float64   `json:"percent_ac_icrack"`
	AreaAcUcrack                   float64   `json:"area_ac_ucrack"`
	PercentAcUcrack                float64   `json:"percent_ac_ucrack"`
	PercentAcRavelling             float64   `json:"percent_ac_ravelling"`
	CcTransverseCrack              float64   `json:"cc_transverse_crack"`
	CcFaulting                     float64   `json:"cc_faulting"`
	CcSpalling                     float64   `json:"cc_spalling"`
	CurrentSurfaceID               int       `json:"current_surface_id"`
	CurrentSurfaceName             string    `json:"current_surface_name"`
	CurrentSurfaceType             string    `json:"current_surface_type"`
	CurrentSurfaceSurfaceGroup     string    `json:"current_surface_surface_group"`
	CurrentSurfaceLayerCoefficient float64   `json:"current_surface_layer_coefficient"`
	CurrentSurfaceDrainage         float64   `json:"current_surface_drainage"`
	CurrentSurfaceA                float64   `json:"current_surface_a"`
	CurrentSurfaceB                float64   `json:"current_surface_b"`
	CurrentSurfaceCBase            float64   `json:"current_surface_c_base"`
	CurrentSurfaceCExp             float64   `json:"current_surface_c_exp"`
	CurrentSurfaceCRT              float64   `json:"current_surface_crt"`
	CurrentSurfaceRRF              float64   `json:"current_surface_rrf"`
	Hsold                          float64   `json:"hsold"`
	Hsnew                          float64   `json:"hsnew"`
	SNPSurface                     float64   `json:"snp_surface"`
	SNPBase                        float64   `json:"snp_base"`
	SNPSubbase                     float64   `json:"snp_subbase"`
	SNP                            float64   `json:"snp"`
	AADT                           float64   `json:"aadt"`
	TruckFactor                    float64   `json:"truck_factor"`
	ESAL                           float64   `json:"esal"`
	YAX                            float64   `json:"yax"`
	Data                           string    `json:"data" gorm:"type:text"`
	TheGeom                        string    `json:"the_geom"`
	CreatedBy                      int       `json:"created_by"` // รหัสผู้ใช้งานที่สร้างข้อมูล
	UpdatedBy                      int       `json:"updated_by"` // รหัสผู้ใช้งานที่อัพเดตข้อมูล
	CreatedAt                      time.Time `json:"created_at"` // วันที่ที่สร้างข้อมูล
	UpdatedAt                      time.Time `json:"updated_at"` // วันที่ที่อัพเดตข้อมูล
}

func (b *PrepareDataForCopy) TableName() string {
	return "prepare_data"
}

func (b *PrepareData) TableName() string {
	return "prepare_data"
}
