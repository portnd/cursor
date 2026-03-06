package models

import "time"

type MaintenancePlanDetail struct {
	ID                         int      `json:"id"`                             // รหัสข้อมูลแผลน
	MaintenancePlanID          int      `json:"maintenance_plan_id"`            // รหัสแผลน
	Schedule                   string   `json:"schedule"`                       // เดือน/ปี
	RefMaintenancePlanStatusID int      `json:"ref_maintenance_plan_status_id"` // สถานะโครงการ
	PlanDate                   string   `json:"plan_date"`                      // วัน / เดือน / ปี
	ProgressPlan               *float64 `json:"progress_plan"`                  // แผนความก้าวหน้า
	DisbursementPlan           *float64 `json:"disbursement_plan"`              // แผนการเบิกจ่าย
	// IsDeleted                  bool      `json:"is_deleted,omitempty"`           // ถูกลบหรือไม่
	CreatedBy int       `json:"created_by"` // รหัสผู้ใช้งานที่สร้างข้อมูล
	UpdatedBy int       `json:"updated_by"` // รหัสผู้ใช้งานที่อัพเดตข้อมูล
	CreatedAt time.Time `json:"created_at"` // วันที่ที่สร้างข้อมูล
	UpdatedAt time.Time `json:"updated_at"` // วันที่ที่อัพเดตข้อมูล
}

type MaintenancePlanDetailData struct {
	MaintenancePlanDetail
	Status RefMaintenancePlanStatus `json:"status" gorm:"ForeignKey:RefMaintenancePlanStatusID;AssociationForeignKey:Id"`
}

func (ral *MaintenancePlanDetail) TableName() string {
	return "maintenance_plan_detail"
}
