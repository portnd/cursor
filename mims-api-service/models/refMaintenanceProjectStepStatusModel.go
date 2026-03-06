package models

// // สถานะโครงการ
type RefMaintenanceProjectStepStatus struct {
	ID                int    `json:"id" note:"รหัสสถานะโครงการ"`
	ProjectStepTypeID int    `json:"project_step_type_id" note:"รหัสประเภทขั้นตอนของโครงการ"`
	Name              string `json:"name" note:"ชื่อสถานะโครงการ"`
	IsDeleted         bool   `json:"is_deleted" default:"false" note:"ถูกลบหรือไม่"`
	CreatedBy         int    `json:"created_by" note:"รหัสผู้ใช้งานที่สร้างข้อมูล"`
	UpdatedBy         int    `json:"updated_by" note:"รหัสผู้ใช้งานที่อัปเดตข้อมูล"`
	CreatedAt         int    `json:"created_at" note:"วันที่ที่สร้างข้อมูล"`
	UpdatedAt         int    `json:"updated_at" note:"วันที่ที่อัปเดตข้อมูล"`
}

func (ral *RefMaintenanceProjectStepStatus) TableName() string {
	return "ref_maintenance_project_step_status"
}
