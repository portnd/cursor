package models

import (
	"database/sql/driver"
	"strings"
	"time"
)

// Todo ...
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
	BudgetID                int       `json:"budget_id"`                 // ประเภทงบประมาณ (fk setting_budget)
	BudgetMethodID          int       `json:"budget_method_id"`          // ประเภทการซ่อมบำรุง (fk setting_intervention_criteria)
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
	CreatedBy               int       `json:"created_by"`                // รหัสผู้ใช้งานที่สร้่างข้อมูล
	UpdatedBy               int       `json:"updated_by"`                // รหัสผู้ใช้งานที่อัพเดตข้อมูล
	CreatedAt               time.Time `json:"created_at"`                // วันที่ที่สร้่างข้อมูล
	UpdatedAt               time.Time `json:"updated_at"`                // วันที่ที่อัพเดตข้อมูล
}

type MaintenancePrepareData struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	BudgetYear        int     `json:"budget_year"`
	BudgetMaintenance float64 `json:"budget_maintenance"`
	BudgetID          int     `json:"budget_id"`
	BudgetMethodID    int     `json:"budget_method_id"` // ประเภทการซ่อมบำรุง
	ContractorName    string  `json:"contractor_name"`
	//BudgetProcurement       float64   `json:"budget_procurement"`
	AdviserName             string    `json:"adviser_name"`
	ProjectEndDate          time.Time `json:"project_end_date"`
	GuaranteeExpirationDate time.Time `json:"guarantee_expiration_date"`
	// IsComplete              bool      `json:"is_complete"`
	LastInspectionDate time.Time `json:"last_inspection_date"`
}

type MaintenanceData struct {
	MaintenancePrepareData
	Budget           SettingBudgetData            `json:"budget" gorm:"ForeignKey:BudgetID;AssociationForeignKey:Id"`
	BudgetMethod     SettingBudgetMethodData      `json:"budget_method" gorm:"ForeignKey:BudgetMethodID;AssociationForeignKey:Id"`
	MaintenanceRoads []MaintenanceRoadPrepareData `json:"maintenance_roads" gorm:"ForeignKey:MaintenanceID;AssociationForeignKey:ID"`
}

type MaintenanceData2 struct {
	MaintenancePrepareData
	Budget           SettingBudgetData             `json:"budget" gorm:"foreignKey:BudgetID;references:Id"`
	BudgetMethod     SettingBudgetMethodData       `json:"budget_method" gorm:"foreignKey:BudgetMethodID;references:Id"`
	MaintenanceRoads []MaintenanceRoadPrepareData2 `json:"maintenance_roads" gorm:"foreignKey:MaintenanceID;references:ID"`
}

// type MaintenanceData2 struct {
// 	MaintenancePrepareData
// 	BudgetID         int                           `json:"budget_id"`
// 	Budget           SettingBudgetData             `json:"budget" gorm:"foreignKey:BudgetID"`
// 	BudgetMethodID   int                           `json:"budget_method_id"`
// 	BudgetMethod     SettingBudgetMethodData       `json:"budget_method" gorm:"foreignKey:BudgetMethodID"`
// 	MaintenanceRoads []MaintenanceRoadPrepareData2 `json:"maintenance_roads" gorm:"foreignKey:MaintenanceID"`
// }

type MaintenanceInterventionCriteriaCondition struct {
	ConditionCriterion  string  `json:"condition_criterion"`
	ConditionLink       string  `json:"condition_link"`
	ConditionOperation1 string  `json:"condition_operation_1"`
	ConditionOperation2 string  `json:"condition_operation_2"`
	ConditionValue1     float64 `json:"condition_value_1"`
	ConditionValue2     int     `json:"condition_value_2"`
	ID                  int     `json:"id"`
}

type MaintenanceInterventionCriteria struct {
	ID int `json:"id"`
	// MaintenanceCondition     []MaintenanceInterventionCriteriaCondition `json:"maintenance_condition"`
	MaintenanceCostPerUnit   int    `json:"maintenance_cost_per_unit"`
	MaintenanceDescription   string `json:"maintenance_description"`
	MaintenanceMethod        string `json:"maintenance_method"`
	MaintenanceScraping      int    `json:"maintenance_scraping"`
	MaintenanceSequence      int    `json:"maintenance_sequence"`
	MaintenanceStandardName  string `json:"maintenance_standard_name"`
	MaintenanceSurfaceTypeID int    `json:"maintenance_surface_type_id"`
	MaintenanceThickness     int    `json:"maintenance_thickness"`
	MaintenanceType          string `json:"maintenance_type"`
}

type MaintainPreload struct {
	MaintainPreloadGetAll
	Roads         []MaintenanceRoadPreload        `json:"roads" gorm:"ForeignKey:MaintenanceID;AssociationForeignKey:ID"`
	RoadHistories []MaintenanceRoadHistoryPreload `json:"roads_histories" gorm:"ForeignKey:MaintenanceID;AssociationForeignKey:ID"`
}

type MaintainPreloadGetAll struct {
	Maintenance
	RoadGroupID    IntDataArray                `json:"road_group_id" gorm:"type:integer[]"`
	RoadGroupNames StringDataArray             `json:"road_group_names" gorm:"type:character[]"`
	KmTotal        float64                     `json:"km_total"`
	RefDivision    RefDivisionRes              `json:"ref_division"  gorm:"foreignKey:RefDivisionCode; references:DivisionCode"`
	RefDistrict    RefDistrictRes              `json:"ref_district"  gorm:"foreignKey:RefDistrictCode; references:DistrictCode"`
	RefDepot       RefDepotRes                 `json:"ref_depot"  gorm:"foreignKey:RefDepotCode; references:DepotCode"`
	CreatedByUser  UserBy                      `json:"created_by_user" gorm:"ForeignKey:CreatedBy;AssociationForeignKey:Id"`
	UpdateByUser   UserBy                      `json:"update_by_user" gorm:"ForeignKey:UpdatedBy;AssociationForeignKey:Id"`
	Budget         SettingBudgetData           `json:"budget" gorm:"ForeignKey:BudgetID;AssociationForeignKey:Id"`
	BudgetMethod   SettingBudgetMethodData     `json:"budget_method" gorm:"ForeignKey:BudgetMethodID;AssociationForeignKey:Id"`
	Attachments    []MaintenanceAttachmentData `json:"attachments" gorm:"ForeignKey:MaintenanceID;AssociationForeignKey:ID"`
}

type MaintenanceList struct {
	MaintainPreloadGetAll
	Roads []MaintenanceRoadPreload `json:"roads" gorm:"ForeignKey:MaintenanceID;AssociationForeignKey:ID"`
}

type MaintainGetAll struct {
	Maintenance
	RoadGroupId   int    `json:"road_group_id" gorm:"column:road_group_id"`
	RoadGroupCode string `json:"road_group_code" gorm:"column:road_group_code"`
	RoadGroupName string `json:"road_group_name" gorm:"column:road_group_name"`

	KmTotal       float64     `json:"km_total" gorm:"column:total_km"`
	StatusProject interface{} `json:"status_project"`
}

type MaintainReport struct {
	MaintainHistoryTableReport []MaintenanceRoadHistoryData `json:"maintain_history_table_report"`
	MaintainChartReport        string                       `json:"maintain_chart_report"`
	MaintainTableReport        MaintainTableReport          `json:"maintain_table_report"`
	MaintainRoadPath           string                       `json:"maintain_road_PATH"`
	IsProblems                 bool                         `json:"is_problems"`
	IsHistory                  bool                         `json:"is_history"`
	Date                       string                       `json:"date"`
	Maintain                   MaintainPreload              `json:"maintain"`
	ApiKey                     string                       `json:"api_key"`
}

type MaintainChartReport struct {
	Name     string                    `json:"name"`
	Data     []MaintainChartDataReport `json:"data"`
	Schedule []string                  `json:"schedule"`
}

type MaintainChartDataReport struct {
	Name  string    `json:"name"`
	Data  []float64 `json:"data"`
	Color string    `json:"color"`
}

type MaintainTableReport struct {
	MaintenancePlan []MaintainTableReportValue `json:"maintenance_plan"`
	Problems        []string                   `json:"problems"`
}

type MaintainTableReportValue struct {
	PlanName string                     `json:"plan_name"`
	Value    []MaintainTableValueReport `json:"value"`
}

type MaintainTableValueReport struct {
	Schedule                  string  `json:"schedule"`
	Plan                      float64 `json:"plan"`
	PlanTotal                 float64 `json:"plan_total"`
	ProgressPlan              float64 `json:"progress_plan"`
	ProgressPlanTotal         float64 `json:"progress_plan_total"`
	DisbursementPlan          float64 `json:"disbursement_plan"`
	DisbursementPlanTotal     float64 `json:"disbursement_plan_total"`
	DisbursementProgress      float64 `json:"disbursement_progress"`
	DisbursementProgressTotal float64 `json:"disbursement_progress_total"`
}

type MaintenancePreloadForDashboard struct {
	Maintenance
	MaintenanceRoads []MaintenanceRoadPreloadForDashboard `json:"maintenance_roads" gorm:"ForeignKey:MaintenanceID;AssociationForeignKey:ID"`
}

func (b *Maintenance) TableName() string {
	return "maintenance"
}

func (b *MaintenancePrepareData) TableName() string {
	return "maintenance"
}

type StringDataArray []string

func (a *StringDataArray) Scan(src interface{}) error {
	nilArray := []string{}
	if src == "{NULL}" {
		*a = nilArray
	} else {
		trimLeft := strings.TrimLeft(src.(string), "{")
		trimRight := strings.TrimRight(trimLeft, "}")
		arrayStr := strings.Split(trimRight, ",")
		*a = arrayStr
	}
	return nil
}

func (a StringDataArray) Value() (driver.Value, error) {
	return a, nil
}
