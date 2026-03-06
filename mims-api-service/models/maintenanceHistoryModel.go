package models

import "time"

// ประวัติซ่อมบำรุง
type MaintenanceHistory struct {
	ID            int       `json:"id" note:"รหัสไฟล์ของแต่ละแผน"`
	MaintenanceID int       `json:"maintenance_id" note:"รหัสโครงการ"`
	Params        string    `json:"params" note:"json ประวัติการซ่อมบำรุง"`
	IsDeleted     bool      `json:"is_deleted" default:"false" note:"ถูกลบหรือไม่"`
	CreatedBy     int       `json:"created_by" note:"รหัสผู้ใช้งานที่สร้างข้อมูล"`
	UpdatedBy     int       `json:"updated_by" note:"รหัสผู้ใช้งานที่อัปเดตข้อมูล"`
	CreatedAt     time.Time `json:"created_at" note:"วันที่ที่สร้างข้อมูล"`
	UpdatedAt     time.Time `json:"updated_at" note:"วันที่ที่อัปเดตข้อมูล"`
}

func (ral *MaintenanceHistory) TableName() string {
	return "maintenance_history"
}
