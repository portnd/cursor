package models

import (
	"time"
)

type MaintenanceAttachment struct {
	ID            int       `json:"id"`             // รหัส
	MaintenanceID int       `json:"maintenance_id"` // รหัสโครงการ
	Path          string    `json:"path"`           // path folder
	FileName      string    `json:"file_name"`      // ชื่อไฟล์
	FileType      string    `json:"file_type"`      // นามสกุลไฟล์
	CreatedBy     int       `json:"created_by"`     // รหัสผู้ใช้งานที่สร้่างข้อมูล
	UpdatedBy     int       `json:"updated_by"`     // รหัสผู้ใช้งานที่อัพเดตข้อมูล
	CreatedAt     time.Time `json:"created_at"`     // วันที่ที่สร้่างข้อมูล
	UpdatedAt     time.Time `json:"updated_at"`     // วันที่ที่อัพเดตข้อมูล
}

type MaintenanceAttachmentData struct {
	ID            int    `json:"id"`        // รหัส
	MaintenanceID int    `json:"-"`         // รหัสโครงการ
	Path          string `json:"path"`      // path folder
	FileName      string `json:"file_name"` // ชื่อไฟล์
}

func (b *MaintenanceAttachment) TableName() string {
	return "maintenance_attachment"
}

func (b *MaintenanceAttachmentData) TableName() string {
	return "maintenance_attachment"
}
