package models

import "time"

// Todo ...
type MaintenanceRoad struct {
	ID                           int       `json:"id"`             // รหัส
	MaintenanceID                int       `json:"maintenance_id"` // รหัสโครงการ
	IDParent                     int       `json:"id_parent"`
	Revision                     int       `json:"revision"`
	Status                       string    `json:"status"`                          // สถานะ
	RoadGroupId                  int       `json:"road_group_Id"`                   //
	RoadID                       int       `json:"road_id"`                         // รหัสสายทาง
	MaintenanceMethodID          int       `json:"maintenance_method_id"`           // ประเภทการซ่อมบำรุง (fk setting_intervention_criteria_method)
	InterventionCriteriaID       int       `json:"intervention_criteria_id"`        // รหัสวิธีการซ่อมบำรุง
	InterventionCriteriaIDParams int       `json:"intervention_criteria_id_params"` // รหัสวิธีการซ่อมบำรุง(params)
	RefSurfaceID                 int       `json:"ref_surface_id"`                  // รหัสประเภทผิว
	RefSurfaceParamsID           int       `json:"ref_surface_params_id"`           // รหัสประเภทผิว(params)
	KmStart                      float64   `json:"km_start"`                        // ช่วง กม. เริ่มต้น
	KmEnd                        float64   `json:"km_end"`
	TheGeom                      string    `json:"the_geom"`
	MaintenanceType              int       `json:"maintenance_type"` // ประเภทการซ่อมบำรุง
	LaneNo                       int       `json:"lane_no"`          // เลนส์
	GridNo                       int       `json:"grid_no"`          // กริด
	CreatedBy                    int       `json:"created_by"`       // รหัสผู้ใช้งานที่สร้่างข้อมูล
	UpdatedBy                    int       `json:"updated_by"`       // รหัสผู้ใช้งานที่อัพเดตข้อมูล
	CreatedAt                    time.Time `json:"created_at"`       // วันที่ที่สร้่างข้อมูล
	UpdatedAt                    time.Time `json:"updated_at"`       // วันที่ที่อัพเดตข้อมูล
}

type MaintenanceRoadData struct {
	ID            int `json:"id" note:"รหัสข้อมูลการซ่อมบำรุง"`
	MaintenanceID int `json:"maintenance_id" note:"รหัสโครงการ"`
	RoadGroupID   int `json:"road_group_id" note:"รหัสสายทาง (fk road_group)"`
	RoadID        int `json:"road_id" note:"ตอนควบคุม (fk road)"`
	Lane          int `json:"lane" note:"ช่องจราจร"`
	// InterventionCriteriaID int          `json:"intervention_criteria_id" note:"วิธีการซ่อมบำรุง (fk setting_intervention_criteria_params)"`
	KmStart float64 `json:"km_start" note:"ช่วง กม. เริ่มต้น"`
	KmEnd   float64 `json:"km_end" note:"ช่วง กม. สิ้นสุด"`
	// RoadSurfaceID          int          `json:"road_surface_id" note:"ตัวแปล (fk road_surface)"`
	// RoadSurfaceParamsID    int          `json:"road_surface_params_id" note:"ตัวแปล (fk road_surface_params)"`
	MaintenanceMethodID          int                      `json:"maintenance_method_id"` // ประเภทการซ่อมบำรุง (fk setting_intervention_criteria_method)
	RefSurfaceID                 int                      `json:"ref_surface_id"`        // ชนิดผิวทาง
	RefSurfaceParamsID           int                      `json:"ref_surface_params_id"`
	InterventionCriteriaID       int                      `json:"intervention_criteria_id"`
	InterventionCriteriaIDParams int                      `json:"intervention_criteria_id_params"`
	TheGeom                      string                   `json:"the_geom"` // geometry
	RoadGroup                    RoadGroup                `json:"road_group" gorm:"ForeignKey:RoadGroupID;AssociationForeignKey:Id"`
	RoadInfo                     RoadInfoData             `json:"road_info" gorm:"ForeignKey:RoadID;AssociationForeignKey:RoadId"`
	InterventionCriteria         InterventionCriteriaData `json:"intervention_criteria" gorm:"ForeignKey:InterventionCriteriaID;AssociationForeignKey:Id"`
}

type MaintenanceRoadPrepareData struct {
	ID            int `json:"id" note:"รหัสข้อมูลการซ่อมบำรุง"`
	MaintenanceID int `json:"maintenance_id" note:"รหัสโครงการ"`
	RoadGroupID   int `json:"road_group_id" note:"รหัสสายทาง (fk road_group)"`
	RoadID        int `json:"road_id" note:"ตอนควบคุม (fk road)"`
	Lane          int `json:"lane" note:"ช่องจราจร"`
	// InterventionCriteriaID int          `json:"intervention_criteria_id" note:"วิธีการซ่อมบำรุง (fk setting_intervention_criteria_params)"`
	KmStart float64 `json:"km_start" note:"ช่วง กม. เริ่มต้น"`
	KmEnd   float64 `json:"km_end" note:"ช่วง กม. สิ้นสุด"`
	// RoadSurfaceID          int          `json:"road_surface_id" note:"ตัวแปล (fk road_surface)"`
	// RoadSurfaceParamsID    int          `json:"road_surface_params_id" note:"ตัวแปล (fk road_surface_params)"`
	MaintenanceMethodID          int          `json:"maintenance_method_id"` // ประเภทการซ่อมบำรุง (fk setting_intervention_criteria_method)
	RefSurfaceID                 int          `json:"ref_surface_id"`        // ชนิดผิวทาง
	RefSurfaceParamsID           int          `json:"ref_surface_params_id"`
	InterventionCriteriaID       int          `json:"intervention_criteria_id"`
	InterventionCriteriaIDParams int          `json:"intervention_criteria_id_params"`
	TheGeom                      string       `json:"the_geom"` // geometry
	RoadGroup                    RoadGroup    `json:"road_group" gorm:"ForeignKey:RoadGroupID;AssociationForeignKey:Id"`
	RoadInfo                     RoadInfoData `json:"road_info" gorm:"ForeignKey:RoadID;AssociationForeignKey:RoadId"`
	// InterventionCriteria         InterventionCriteriaData `json:"intervention_criteria" gorm:"ForeignKey:InterventionCriteriaID;AssociationForeignKey:Id"`
	InterventionCriteria MaintenanceInterventionCriteria `json:"intervention_criteria"`
	RefSurface           RefSurfaceParam                 `json:"ref_surface" gorm:"ForeignKey:RefSurfaceParamsID;AssociationForeignKey:ID"`
	LastInspectionDate   time.Time                       `json:"last_inspection_date"`
}

type MaintenanceRoadPrepareData2 struct {
	ID            int `json:"id" note:"รหัสข้อมูลการซ่อมบำรุง"`
	MaintenanceID int `json:"maintenance_id" note:"รหัสโครงการ"`
	RoadGroupID   int `json:"road_group_id" note:"รหัสสายทาง (fk road_group)"`
	RoadID        int `json:"road_id" note:"ตอนควบคุม (fk road)"`
	Lane          int `json:"lane" note:"ช่องจราจร"`
	// InterventionCriteriaID int          `json:"intervention_criteria_id" note:"วิธีการซ่อมบำรุง (fk setting_intervention_criteria_params)"`
	KmStart float64 `json:"km_start" note:"ช่วง กม. เริ่มต้น"`
	KmEnd   float64 `json:"km_end" note:"ช่วง กม. สิ้นสุด"`
	// RoadSurfaceID          int          `json:"road_surface_id" note:"ตัวแปล (fk road_surface)"`
	// RoadSurfaceParamsID    int          `json:"road_surface_params_id" note:"ตัวแปล (fk road_surface_params)"`
	MaintenanceMethodID          int                      `json:"maintenance_method_id"` // ประเภทการซ่อมบำรุง (fk setting_intervention_criteria_method)
	RefSurfaceID                 int                      `json:"ref_surface_id"`        // ชนิดผิวทาง
	RefSurfaceParamsID           int                      `json:"ref_surface_params_id"`
	InterventionCriteriaID       int                      `json:"intervention_criteria_id"`
	InterventionCriteriaIDParams int                      `json:"intervention_criteria_id_params"`
	TheGeom                      string                   `json:"the_geom"` // geometry
	RoadGroup                    RoadGroup                `json:"road_group" gorm:"ForeignKey:RoadGroupID;AssociationForeignKey:Id"`
	RoadInfo                     RoadInfoData             `json:"road_info" gorm:"ForeignKey:RoadID;AssociationForeignKey:RoadId"`
	InterventionCriteria         InterventionCriteriaData `json:"intervention_criteria" gorm:"ForeignKey:InterventionCriteriaID;AssociationForeignKey:Id"`
	// InterventionCriteria MaintenanceInterventionCriteria `json:"intervention_criteria"`
	RefSurface         RefSurfaceParam `json:"ref_surface" gorm:"ForeignKey:RefSurfaceParamsID;AssociationForeignKey:ID"`
	LastInspectionDate time.Time       `json:"last_inspection_date"`
}

type MaintenanceRoadPreload struct {
	ID                           int                      `json:"id"`        // รหัส
	IDParent                     int                      `json:"id_parent"` //
	MaintenanceID                int                      `json:"-"`         // รหัสโครงการ
	Status                       string                   `json:"status"`    // สถานะ
	Color                        string                   `json:"color"`
	RoadGroupId                  int                      `json:"road_group_Id"` //
	RoadGroupName                string                   `json:"road_group_name"`
	RoadID                       int                      `json:"road_id"`   // รหัสสายทาง
	RoadName                     string                   `json:"road_name"` // ชื่อ จาก ถึง
	RoadLevel                    int                      `json:"road_level"`
	LaneTotal                    int                      `json:"lane_total"`
	RefDirectionId               int                      `json:"ref_direction_id"`
	RefDirectionName             string                   `json:"ref_direction_name"`
	RoadSecNameOr                string                   `json:"road_sec_name_or"`
	RoadSecNameDes               string                   `json:"road_sec_name_des"`
	Distance                     float64                  `json:"distance"`
	InterventionCriteriaID       int                      `json:"intervention_criteria_id"`        // รหัสวิธีการซ่อมบำรุง
	InterventionCriteriaIDParams int                      `json:"intervention_criteria_id_params"` // รหัสวิธีการซ่อมบำรุง(params)
	KmStart                      float64                  `json:"km_start"`                        // ช่วง กม. เริ่มต้น
	KmEnd                        float64                  `json:"km_end"`                          // ช่วง กม. สิ้นสุด
	TheGeom                      string                   `json:"the_geom"`
	TheGeomString                string                   `json:"the_geom_string"`
	MaintenanceType              int                      `json:"maintenance_type"` // ประเภทการซ่อมบำรุง
	LaneNo                       int                      `json:"lane_no"`          // เลนส์
	GridNo                       int                      `json:"grid_no"`          // กริด
	InterventionCriteria         InterventionCriteriaData `json:"intervention_criteria" gorm:"ForeignKey:InterventionCriteriaID;AssociationForeignKey:ID"`
}

type MaintenanceRoadPreloadById struct {
	ID                           int                      `json:"id"`        // รหัส
	IDParent                     int                      `json:"id_parent"` //
	MaintenanceID                int                      `json:"-"`         // รหัสโครงการ
	Status                       string                   `json:"status"`    // สถานะ
	Color                        string                   `json:"color"`
	RoadGroupId                  int                      `json:"road_group_Id"` //
	RoadGroupName                string                   `json:"road_group_name"`
	RoadID                       int                      `json:"road_id"`   // รหัสสายทาง
	RoadName                     string                   `json:"road_name"` // ชื่อ จาก ถึง
	RoadLevel                    int                      `json:"road_level"`
	RefDirectionId               int                      `json:"ref_direction_id"`
	RefDirectionName             string                   `json:"ref_direction_name"`
	RoadSecNameOr                string                   `json:"road_sec_name_or"`
	RoadSecNameDes               string                   `json:"road_sec_name_des"`
	Distance                     float64                  `json:"distance"`
	LaneTotal                    int                      `json:"lane_total"`
	InterventionCriteriaID       int                      `json:"intervention_criteria_id"`        // รหัสวิธีการซ่อมบำรุง
	InterventionCriteriaIDParams int                      `json:"intervention_criteria_id_params"` // รหัสวิธีการซ่อมบำรุง(params)
	KmStart                      float64                  `json:"km_start"`                        // ช่วง กม. เริ่มต้น
	KmEnd                        float64                  `json:"km_end"`                          // ช่วง กม. สิ้นสุด
	TheGeom                      string                   `json:"the_geom"`
	MaintenanceType              int                      `json:"maintenance_type"` // ประเภทการซ่อมบำรุง
	LaneNo                       int                      `json:"lane_no"`          // เลนส์
	GridNo                       int                      `json:"grid_no"`          // กริด
	InterventionCriteria         InterventionCriteriaData `json:"intervention_criteria" gorm:"ForeignKey:InterventionCriteriaID;AssociationForeignKey:ID"`
	Maintenance                  Maintenance              `json:"maintenance" gorm:"ForeignKey:ID; references:MaintenanceID"`
}
type MaintenanceRoadPreloadForDashboard struct {
	MaintenanceRoad
	TheGeomJson string                      `json:"the_geom_json"`
	Road        RoadForDashboardMaintenance `json:"road" gorm:"ForeignKey:RoadID;AssociationForeignKey:Id"`
}

func (b *MaintenanceRoad) TableName() string {
	return "maintenance_road"
}

func (b *MaintenanceRoadData) TableName() string {
	return "maintenance_road"
}
func (b *MaintenanceRoadPrepareData) TableName() string {
	return "maintenance_road"
}

func (b *MaintenanceRoadPrepareData2) TableName() string {
	return "maintenance_road"
}

func (b *MaintenanceRoadPreload) TableName() string {
	return "maintenance_road"
}

func (b *MaintenanceRoadPreloadById) TableName() string {
	return "maintenance_road"
}

type MaintenanceRoadProjectEndDate struct {
	RoadId         int       `json:"road_id"` // รหัสสายทาง
	LaneNo         int       `json:"lane_no"`
	ProjectEndDate time.Time `json:"project_end_date"` // วันที่สิ้นสุดโครงการ
}

func (b *MaintenanceRoadProjectEndDate) TableName() string {
	return "maintenance_road"
}
