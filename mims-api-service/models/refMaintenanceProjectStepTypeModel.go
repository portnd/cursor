package models

import "time"

// Todo ...
// สถานะขั้นตอนโครงการ
type RefMaintenanceProjectStepType struct {
	ID        int       `json:"id" note:"รหัสประเภทขั้นตอนของโครงการ"`
	Name      string    `json:"name" note:"ชื่อขั้นตอนของโครงการ"`
	Seq       int       `json:"seq" note:"ลำดับขั้นตอนของโครงการ"`
	IsDeleted bool      `json:"is_deleted" default:"false" note:"ถูกลบหรือไม่"`
	CreatedBy int       `json:"created_by" note:"รหัสผู้ใช้งานที่สร้างข้อมูล"`
	UpdatedBy int       `json:"updated_by" note:"รหัสผู้ใช้งานที่อัปเดตข้อมูล"`
	CreatedAt time.Time `json:"created_at" note:"วันที่ที่สร้างข้อมูล"`
	UpdatedAt time.Time `json:"updated_at" note:"วันที่ที่อัปเดตข้อมูล"`
}

func (b *RefMaintenanceProjectStepType) TableName() string {
	return "ref_maintenance_project_step_type"
}
