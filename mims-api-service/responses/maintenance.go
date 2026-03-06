package responses

import (
	"time"

	"gitlab.com/mims-api-service/models"
)

type MaintenanceID struct {
	ID       int `json:"id"`
	IDParent int `json:"id_parent"`
}

type MaintenanceRoadID struct {
	ID       int `json:"id"`
	IDParent int `json:"id_parent"`
}

type MaintenanceRoadHisID struct {
	ID       int `json:"id"`
	IDParent int `json:"id_parent"`
}
type MaintenancePlanRes struct {
	ID    int                       `json:"id"`
	Name  string                    `json:"name"`
	Steps []MaintenancePlanStepsRes `json:"steps"`
}

type MaintenancePlanStepsRes struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	Date             *string  `json:"date"`
	ProcessPlan      *float64 `json:"process_plan"`
	DisbursementPlan *float64 `json:"disbursement_plan"`
}

type MaintainPreloadGetAll struct {
	Maintenance
	RemainingTime string `json:"remaining_time"`
	Color         string `json:"color"`
	//RemainingTimeHeader string                         `json:"remaining_time_header"`
	KmTotal        float64                        `json:"km_total"`
	OwnerCode      string                         `json:"owner_code"`
	OwnerName      string                         `json:"owner_name"`
	RoadGroupNames []string                       `json:"road_group_names"`
	CreatedBy      UserBy                         `json:"created_by" gorm:"ForeignKey:CreatedBy;AssociationForeignKey:Id"`
	UpdateBy       UserBy                         `json:"update_by" gorm:"ForeignKey:UpdatedBy;AssociationForeignKey:Id"`
	Budget         models.SettingBudgetData       `json:"budget" gorm:"ForeignKey:BudgetID;AssociationForeignKey:Id"`
	BudgetMethod   models.SettingBudgetMethodData `json:"budget_method" gorm:"ForeignKey:BudgetMethodID;AssociationForeignKey:Id"`
	RefDivision    models.RefDivisionRes          `json:"ref_division"  gorm:"foreignKey:RefDivisionCode; references:DivisionCode"`
	RefDistrict    models.RefDistrictRes          `json:"ref_district"  gorm:"foreignKey:RefDistrictCode; references:DistrictCode"`
	RefDepot       models.RefDepotRes             `json:"ref_depot"  gorm:"foreignKey:RefDepotCode; references:DepotCode"`
	IsShowMethod   bool                           `json:"is_show_method"`
}

type MaintenanceRoadPreload struct {
	ID               int      `json:"id"` // รหัส
	IDParent         int      `json:"id_parent"`
	RoadID           int      `json:"road_id"` // รหัสสายทาง
	Color            string   `json:"color"`
	RoadGroupName    string   `json:"road_group_name"`
	RoadName         string   `json:"road_name"` // ชื่อ จาก ถึง
	LaneNo           *int     `json:"lane_no"`   // เลนส์
	GridNo           *int     `json:"grid_no"`   // กริด
	KmStart          float64  `json:"km_start"`  // ช่วง กม. เริ่มต้น
	KmEnd            float64  `json:"km_end"`    // ช่วง กม. สิ้นสุด
	MaintenanceType  int      `json:"maintenance_type"`
	Distance         float64  `json:"distance"`
	RefDirectionId   int      `json:"ref_direction_id"`
	RefDirectionName string   `json:"ref_direction_name"`
	LaneTotal        int      `json:"lane_total"`
	TheGeom          GeomJSON `json:"the_geom"`
	//IsShowMethod         bool                            `json:"is_show_method"`
	InterventionCriteria models.InterventionCriteriaData `json:"intervention_criteria" gorm:"ForeignKey:InterventionCriteriaID;AssociationForeignKey:ID"`
}

type MaintenanceRoadHistoryPreload struct {
	ID                   int                             `json:"id"`      // รหัส
	RoadID               int                             `json:"road_id"` // รหัสสายทาง
	Color                string                          `json:"color"`
	RoadGroupName        string                          `json:"road_group_name"`
	RoadName             string                          `json:"road_name"` // ชื่อ จาก ถึง
	LaneNo               *int                            `json:"lane_no"`   // เลนส์
	GridNo               *int                            `json:"grid_no"`   // กริด
	KmStart              float64                         `json:"km_start"`  // ช่วง กม. เริ่มต้น
	KmEnd                float64                         `json:"km_end"`    // ช่วง กม. สิ้นสุด
	MaintenanceType      int                             `json:"maintenance_type"`
	RefDirectionId       int                             `json:"ref_direction_id"`
	RefDirectionName     string                          `json:"ref_direction_name"`
	LaneTotal            int                             `json:"lane_total"`
	TheGeom              GeomJSON                        `json:"the_geom"`
	InterventionCriteria models.InterventionCriteriaData `json:"intervention_criteria" gorm:"ForeignKey:InterventionCriteriaID;AssociationForeignKey:ID"`
	//IsShowMethod         bool                                          `json:"is_show_method"`
	Attachments []models.MaintenanceRoadHistoryAttachmentData `json:"attachments"`
}

type MaintenanceList struct {
	MaintainPreloadGetAll
	Roads []MaintenanceRoadPreload `json:"roads"`
}

type MaintenanceById struct {
	MaintainPreloadGetAll
	Attachments   []models.MaintenanceAttachmentData `json:"attachments" gorm:"ForeignKey:MaintenanceID;AssociationForeignKey:ID"`
	Roads         []MaintenanceRoadPreload           `json:"roads"`
	RoadHistories []MaintenanceRoadPreload           `json:"road_histories"`
}

type MaintenanceByRoadId struct {
	MaintainPreloadGetAll
	Attachments   []models.MaintenanceAttachmentData `json:"attachments" gorm:"ForeignKey:MaintenanceID;AssociationForeignKey:ID"`
	Roads         []MaintenanceRoadPreload           `json:"roads"`
	RoadHistories []MaintenanceRoadPreload           `json:"road_histories"`
}
type MaintenanceProgress struct {
	Progress         float64 `json:"progress"`
	Disbursement     float64 `json:"disbursement"`
	DisbursementPlan float64 `json:"disbursement_plan"`
}

type MaintenancHistoryeList struct {
	models.MaintainPreloadGetAll
	MaintenanceRoads         []models.MaintenanceRoadData `json:"maintenance_roads"`
	MaintenanceRoadHistories interface{}                  `json:"maintenance_road_histories"`
	PercentProgress          float64                      `json:"percent_progress"`
	PercentPay               float64                      `json:"percent_pay"`
}

type MaintenanceListByID struct {
	Maintenance
	MaintenanceRoadRes               []MaintenanceRoad                  `json:"maintenance_road" gorm:"ForeignKey:MaintenanceID;AssociationForeignKey:ID"`
	MaintenanceProjectStepProcessRes []MaintenanceProjectStepProcess    `json:"maintenance_project_step_process" gorm:"ForeignKey:MaintenanceID;AssociationForeignKey:ID"`
	SettingBudget                    SettingBudget                      `json:"setting_budget"`
	SettingInterventionCriteria      models.SettingInterventionCriteria `json:"setting_intervention_criteria" `
}

type MaintenanceRoad struct {
	ID                     int     `json:"id" note:"รหัสข้อมูลการซ่อมบำรุง"`
	MaintenanceID          int     `json:"maintenance_id" note:"รหัสโครงการ"`
	RoadGroupID            int     `json:"road_group_id" note:"รหัสสายทาง (fk road_group)"`
	RoadID                 int     `json:"road_id" note:"ตอนควบคุม (fk road)"`
	Lane                   int     `json:"lane" note:"ช่องจราจร"`
	InterventionCriteriaID int     `json:"intervention_criteria_id" note:"วิธีการซ่อมบำรุง (fk setting_intervention_criteria_params)"`
	KmStart                float64 `json:"km_start" note:"ช่วง กม. เริ่มต้น"`
	KmEnd                  float64 `json:"km_end" note:"ช่วง กม. สิ้นสุด"`
	KmTotal                float64 `json:"km_total" note:"ระยะทาง (กม.)"`
	TheGeom                string  `json:"the_geom" sql:"default:null"`
	RoadSurfaceID          int     `json:"road_surface_id" note:"ตัวแปล (fk road_surface)"`
	RoadSurfaceParamsID    int     `json:"road_surface_params_id" note:"ตัวแปล (fk road_surface_params)"`
	IsDeleted              bool    `json:"is_deleted" default:"false" note:"ถูกลบหรือไม่"`
	CreatedBy              int     `json:"created_by" note:"รหัสผู้ใช้งานที่สร้างข้อมูล"`
	UpdatedBy              int     `json:"updated_by" note:"รหัสผู้ใช้งานที่อัปเดตข้อมูล"`
	CreatedAt              string  `json:"created_at" note:"วันที่ที่สร้างข้อมูล"`
	UpdatedAt              string  `json:"updated_at" note:"วันที่ที่อัปเดตข้อมูล"`
}

type MaintenanceProjectStepProcess struct {
	ID            int    `json:"id" note:"รหัสดำเนินโครงการ"`
	MaintenanceID int    `json:"maintenance_id" note:"รหัสโครงการ"`
	StepTypeID    int    `json:"step_type_id" note:"รหัสประเภทขั้นตอนของโครงการ"`
	Name          string `json:"name" note:"ชื่อของแต่ละงวด"`
	IsCurrent     bool   `json:"is_current" note:"ขั้นตอนนี้สำเร็จแล้วหรือไม่"`
	IsDeleted     bool   `json:"is_deleted" default:"false" note:"ถูกลบหรือไม่"`
	CreatedBy     int    `json:"created_by" note:"รหัสผู้ใช้งานที่สร้างข้อมูล"`
	UpdatedBy     int    `json:"updated_by" note:"รหัสผู้ใช้งานที่อัปเดตข้อมูล"`
	CreatedAt     string `json:"created_at" note:"วันที่ที่สร้างข้อมูล"`
	UpdatedAt     string `json:"updated_at" note:"วันที่ที่อัปเดตข้อมูล"`
}
type RoadGroupMaintain struct {
	Id   int    `json:"id" gorm:"column:road_group_id"`
	Code string `json:"code" gorm:"column:code"`
	Name string `json:"name" gorm:"column:road_group_name"`
}

type SettingBudget struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	BudgetTypeId int     `json:"budget_type_id"`
	CostPerUnit  float64 `json:"cost_per_unit"`
	IsShowMethod bool    `json:"is_show_method"`
	IsDeleted    bool    `json:"is_deleted"`
	UpdatedBy    int     `json:"updated_by"`
	CreatedBy    int     `json:"created_by"`
	CreatedAt    string  `json:"created_at" `
	UpdatedAt    string  `json:"updated_at" `
}

type Maintenance struct {
	ID                      int       `json:"id"` // รหัสโครงการ
	IDParent                *int      `json:"id_parent"`
	Revision                int       `json:"revision"`
	Status                  string    `json:"status"`                    // สถานะ
	Name                    string    `json:"name"`                      // ชื่อโครงการ
	RefDivisionCode         *string   `json:"ref_division_code"`         // แผนก
	RefDistrictCode         *string   `json:"ref_district_code"`         // เขตทาง
	RefDepotCode            *string   `json:"ref_depot_code"`            // หมวด
	ContractNumber          string    `json:"contract_number"`           // เลขที่สัญญา
	BudgetYear              int       `json:"budget_year"`               // ปีงบประมาณ (ค.ศ.)
	BudgetMaintenance       float64   `json:"budget_maintenance"`        // วงเงินงบประมาณ (ไม่รวม VAT)
	MiddlePrice             float64   `json:"middle_price"`              // ราคากลาง (รวม vat)
	ContractWorkValue       float64   `json:"contract_work_value"`       // มูลค่างานตามสัญญา (รวม vat)
	BudgetProcurement       float64   `json:"budget_procurement"`        // ราคาที่จัดซื้อจัดจ้าง
	AdvisorName             string    `json:"advisor_name"`              // บริษัทที่ปรึกษาโครงการ
	ContractorName          string    `json:"contractor_name"`           // บริษัทผู้รับจ้าง
	ProjectSecretaryName    string    `json:"project_secretary_name"`    // ชื่อ-นามสกุลเลขาโครงการ
	ProjectEndDate          time.Time `json:"project_end_date"`          // วันที่สิ้นสุดโครงการ
	GuaranteeExpirationDate time.Time `json:"guarantee_expiration_date"` // วันที่หมดการค้ำประกัน
	ProjectDetails          string    `json:"project_details"`           // รายละเอียดโครงการ
	CreatedAt               time.Time `json:"created_at"`                // วันที่ที่สร้่างข้อมูล
	UpdatedAt               time.Time `json:"updated_at"`                // วันที่ที่อัพเดตข้อมูล
}

type MaintenanceInterventionCriteria struct {
	ID       int                                       `json:"id"`
	Label    string                                    `json:"label"`
	Children []MaintenanceInterventionCriteriaChildren `json:"children"`
}

type MaintenanceInterventionCriteriaChildren struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

type MaintenancePlanProgres struct {
	ID               *int       `json:"id"`
	IsSelect         bool       `json:"is_select"`
	Schedule         string     `json:"schedule"`
	DisbursementDate *time.Time `json:"disbursement_date"`
	Progress         *float64   `json:"progress"`
	Disbursement     *float64   `json:"disbursement"`
	Problems         string     `json:"problems"`
	//Attachments      []models.MaintenancePlanProgressAttachmentData `json:"attachments"`
}

type MaintenancePlanAtt struct {
	Filepath string `json:"filepath"`
}

type CheckMaintenancePlanDelete struct {
	Message string `json:"message"`
}

type MaintenanceStatus struct {
	Name      string                       `json:"name"`
	IsCurrent bool                         `json:"is_current"`
	Schedules []MaintenanceStatusSchedules `json:"schedules"`
}

type MaintenanceStatusSchedules struct {
	Schedule         string     `json:"schedule"`
	DisbursementDate *time.Time `json:"disbursement_date"`
	IsChecked        bool       `json:"is_checked"`
	Status           string     `json:"status"`
}

type MaintenanceReportRes struct {
	Name     string              `json:"name"`
	Data     []MaintenanceReport `json:"data"`
	Schedule []string            `json:"schedule"`
}
type MaintenanceReport struct {
	Name  string      `json:"name"`
	Data  interface{} `json:"data"`
	Color string      `json:"color"`
}

type MaintenanceTableReporDatatRes struct {
	Problems        interface{}                 `json:"problems"`
	MaintenancePlan []MaintenanceTableReportRes `json:"maintenance_plan"`
}

type MaintenanceTableReportRes struct {
	PlanName string      `json:"plan_name"`
	Values   interface{} `json:"value"`
	// Problems []string    `json:"problems"`
}

type MaintenanceTableReport struct {
	Schedule                  string      `json:"schedule"`
	Plan                      interface{} `json:"plan"`
	PlanTotal                 interface{} `json:"plan_total"`
	ProgressPlan              interface{} `json:"progress_plan"`
	ProgressPlanTotal         interface{} `json:"progress_plan_total"`
	DisbursementPlan          interface{} `json:"disbursement_plan"`
	DisbursementPlanTotal     interface{} `json:"disbursement_plan_total"`
	DisbursementProgress      interface{} `json:"disbursement_progress"`
	DisbursementProgressTotal interface{} `json:"disbursement_progress_total"`
}

type MaintenanceInterventionCriteriaModel struct {
	ID       int                                            `json:"id"`
	Label    string                                         `json:"label"`
	Children []MaintenanceInterventionCriteriaChildrenModel `json:"children"`
}

type MaintenanceInterventionCriteriaChildrenModel struct {
	ID       int    `json:"id"`
	Label    string `json:"label"`
	Selected bool   `json:"-"`
}
