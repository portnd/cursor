package requests

type MaintenanceReq struct {
	Name                    string   `json:"name" validate:"nonzero"`               // ชื่อโครงการ
	OwnerCode               string   `json:"owner_code" validate:"nonzero"`         // เลขที่สัญญา
	BudgetYear              int      `json:"budget_year" validate:"nonzero"`        // ปีงบประมาณ (ค.ศ.)
	ContractNumber          string   `json:"contract_number"`                       // เลขที่สัญญา
	BudgetID                int      `json:"budget_id" validate:"nonzero"`          // ประเภทงบประมาณ (fk setting_budget)
	BudgetMethodID          int      `json:"budget_method_id" validate:"nonzero"`   // ประเภทการซ่อมบำรุง
	BudgetMaintenance       *float64 `json:"budget_maintenance" `                   // วงเงินงบประมาณ (ไม่รวม vat)
	MiddlePrice             *float64 `json:"middle_price"`                          // ราคากลาง (รวม vat)
	ContractWorkValue       *float64 `json:"contract_work_value"`                   // มูลค่างานตามสัญญา (รวม vat)
	AdvisorName             string   `json:"advisor_name"`                          // ชื่อ-นามสกุลที่ปรึกษาโครงการ
	ContractorName          string   `json:"contractor_name"`                       // ชื่อ-นามสกุลผู้รับจ้าง
	BudgetProcurement       *float64 `json:"budget_procurement" validate:"nonzero"` // ราคาที่จัดซื้อจัดจ้าง
	ProjectSecretaryName    string   `json:"project_secretary_name"`                // ชื่อ-นามสกุลเลขาโครงการ
	ProjectEndDate          string   `json:"project_end_date"  validate:"nonzero"`
	GuaranteeExpirationDate string   `json:"guarantee_expiration_date"  validate:"nonzero"`
	ProjectDetails          string   `json:"project_details"` // รายละเอียดโครงการ
	Attachments             []struct {
		ID       *int   `json:"id"`
		File     string `json:"file"`
		FileName string `json:"file_name"`
		Status   string `json:"status"`
	} `json:"attachments"`
}

type MaintenanceRoadsReq struct {
	RoadID                 *int     `json:"road_id" validate:"nonzero"` // ตอนควบคุม (fk road)
	LaneNo                 *int     `json:"lane_no"`                    // ช่องจราจร
	GridNo                 *int     `json:"grid_no"`
	MaintenanceType        int      `json:"maintenance_type" validate:"nonzero"`
	InterventionCriteriaID *int     `json:"intervention_criteria_id"` // ประเภทการซ่อมบำรุง (fk setting_intervention_criteria_method)
	KmStart                *float64 `json:"km_start"`                 // ช่วง กม. เริ่มต้น
	KmEnd                  *float64 `json:"km_end"`                   // ช่วง กม. สิ้นสุด
}

type MaintenanceRoadHistoryReq struct {
	RoadID                 *int     `json:"road_id" validate:"nonzero"` // ตอนควบคุม (fk road)
	LaneNo                 *int     `json:"lane_no"`                    // ช่องจราจร
	GridNo                 *int     `json:"grid_no"`
	MaintenanceType        int      `json:"maintenance_type" validate:"nonzero"`
	InterventionCriteriaID *int     `json:"intervention_criteria_id"` // ประเภทการซ่อมบำรุง (fk setting_intervention_criteria_method)
	KmStart                *float64 `json:"km_start"`                 // ช่วง กม. เริ่มต้น
	KmEnd                  *float64 `json:"km_end"`                   // ช่วง กม. สิ้นสุด
	// Attachments            []struct {
	// 	ID       *int   `json:"id"`
	// 	File     string `json:"file"`
	// 	FileName string `json:"file_name"`
	// 	Status   string `json:"status"`
	// } `json:"attachments"`
}

type MaintenanceAttachmentsReq struct {
	ID       *int   `json:"id"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	Path     string `json:"path"`
	Status   string `json:"status"`
}

type MaintenancePlan struct {
	ID           *int           `json:"id"`
	Name         string         `json:"name" validate:"nonzero"`
	StartDate    string         `json:"start_date" validate:"nonzero"`
	EndDate      string         `json:"end_date" validate:"nonzero"`
	IsCurrent    bool           `json:"is_current" validate:"nonzero"`
	SchedulePlan []SchedulePlan `json:"schedule_plan" validate:"nonzero"`
}

type SchedulePlan struct {
	ID               *int     `json:"id"`
	Schedule         string   `json:"schedule" validate:"nonzero"`
	Status           int      `json:"status" validate:"nonzero"`
	PlanDate         string   `json:"plan_date" validate:"nonzero"`
	ProgressPlan     float64  `json:"progress_plan" validate:"nonzero"`
	DisbursementPlan *float64 `json:"disbursement_plan" validate:"nonzero"`
}

type MaintenancePlanProgres struct {
	ID               *int     `json:"id"`
	IsSelect         bool     `json:"is_select"`
	Schedule         string   `json:"schedule"`
	DisbursementDate *string  `json:"disbursement_date"`
	Progress         *float64 `json:"progress"`
	Disbursement     *float64 `json:"disbursement"`
	Problems         string   `json:"problems"`
	// Attachments      []MaintenancePlanAtt `json:"attachments"`
}

type MaintenancePlanProgresDataReq struct {
	DisbursementDate *string                  `json:"disbursement_date"`
	Progress         *float64                 `json:"progress"`
	Disbursement     *float64                 `json:"disbursement"`
	Problems         string                   `json:"problems"`
	Attachments      []MaintenanceAttachments `json:"attachments"`

	// FileCount        int      `form:"file_count"`
	// Attachments      []MaintenancePlanAtt `json:"attachments"`
}

type MaintenanceAttachments struct {
	ID       *int   `่json:"id"`
	File     string `่json:"file"`
	FileName string `json:"file_name"`
	Status   string `่json:"status"`
}

type MaintenancePlanAttReq struct {
	Path     string `json:"path"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
}

type MaintenancePrams struct {
	OwnerCode            *string `json:"owner_code"`
	BudgetYear           *int    `json:"budget_year"`             // ปีงบประมาณ
	BudgetMethodId       []int   `json:"budget_type"`             //ประเภทงบประมาณ
	RoadGroupID          []int   `json:"road_group_id"`           // สายทาง
	Name                 *string `json:"name"`                    //ชื่อโครการ
	RoadGroupIDDashboard []int   `json:"road_group_id_dashboard"` // สายทาง

}

type MaintenanceAttPrams struct {
	FileType []string `json:"file_type"`
	Order    string   `json:"order"`
}

type MaintenanceFinished struct {
	LastInspectionDate      string `json:"last_inspection_date" validate:"nonzero"`      // วันที่ตรวจรับงานงวดสุดท้าย
	GuaranteeExpirationDate string `json:"guarantee_expiration_date" validate:"nonzero"` // วันที่หมดการค้ำประกันโครงการ
}

type MaintenanceAnalysisModel struct {
	Iri1                   *float64 `json:"iri1" form:"iri1"`
	Iri2                   *float64 `json:"iri2" form:"iri2"`
	Aadt1                  *float64 `json:"aadt1" form:"aadt1"`
	Aadt2                  *float64 `json:"aadt2" form:"aadt2"`
	Age1                   *float64 `json:"age1" form:"age1"`
	Age2                   *float64 `json:"age2" form:"age2"`
	Ifi1                   *float64 `json:"ifi1" form:"ifi1"`
	Ifi2                   *float64 `json:"ifi2" form:"ifi2"`
	LaneNo                 *int     `json:"lane_no" form:"lane_no"`
	InterventionCriteriaID *int     `json:"intervention_criteria_id" form:"intervention_criteria_id"`
	RoadIDs                string   `json:"road_ids" form:"road_ids"`
}

type MaintenanceAnalysisModelReq struct {
	ID                     string `json:"id"`
	InterventionCriteriaID int    `json:"intervention_criteria_id"`
}
