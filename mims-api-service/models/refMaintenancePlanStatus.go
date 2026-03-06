package models

import "time"

type RefMaintenancePlanStatus struct {
	ID        int       `json:"id"`   // รหัสข้อมูลแผลน
	Name      string    `json:"name"` // ชื่อ
	IsDeleted bool      `json:"-"`    // ถูกลบหรือไม่
	CreatedBy int       `json:"-"`    // รหัสผู้ใช้งานที่สร้างข้อมูล
	UpdatedBy int       `json:"-"`    // รหัสผู้ใช้งานที่อัพเดตข้อมูล
	CreatedAt time.Time `json:"-"`    // วันที่ที่สร้างข้อมูล
	UpdatedAt time.Time `json:"-"`    // วันที่ที่อัพเดตข้อมูล
}

func (ral *RefMaintenancePlanStatus) TableName() string {
	return "ref_maintenance_plan_status"
}
