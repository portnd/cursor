package models

import "time"

type MaintenancePlanDetailProgress struct {
	ID                      *int       `json:"id"`                         // รหัสข้อมูลความก้าวหน้า
	MaintenanceID           int        `json:"maintenance_id"`             // รหัสข้อมูลแผลน
	MaintenancePlanDetailID int        `json:"maintenance_plan_detail_id"` // รหัสข้อมูลแผลน
	MaintenancePlanID       int        `json:"maintenance_plan_id"`        // รหัสแผลน
	IsSelect                bool       `json:"is_select"`                  // เลือกไปแสดงหรือไม่
	Schedule                time.Time  `json:"schedule"`
	DisbursementDate        *time.Time `json:"disbursement_date"`    // วันที่เบิกจ่าย
	Progress                *float64   `json:"progress"`             // ความก้าวหน้า
	Disbursement            *float64   `json:"disbursement"`         // ผลเบิกจ่าย
	Problems                string     `json:"problems"`             // ปัญหาและอุปสรรค
	IsDeleted               bool       `json:"is_deleted,omitempty"` // ถูกลบหรือไม่
	CreatedBy               int        `json:"created_by"`           // รหัสผู้ใช้งานที่สร้างข้อมูล
	UpdatedBy               int        `json:"updated_by"`           // รหัสผู้ใช้งานที่อัพเดตข้อมูล
	CreatedAt               time.Time  `json:"created_at"`           // วันที่ที่สร้างข้อมูล
	UpdatedAt               time.Time  `json:"updated_at"`           // วันที่ที่อัพเดตข้อมูล
}

func (ral *MaintenancePlanDetailProgress) TableName() string {
	return "maintenance_plan_detail_progress"
}
