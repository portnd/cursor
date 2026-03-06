package models

import "time"

type MaintenanceRoadHistoryAttachment struct {
	ID                       *int      `json:"id"`
	MaintenanceID            int       `json:"maintenance_id"`
	MaintenanceRoadHistoryID int       `json:"maintenance_road_history_id"`
	Path                     string    `json:"path"`
	FileName                 string    `json:"file_name"`
	FileType                 string    `json:"file_type"`
	CreatedBy                int       `json:"created_by"`
	UpdatedBy                int       `json:"updated_by"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type MaintenanceRoadHistoryAttachmentData struct {
	ID                       int    `json:"id"`        // รหัส
	MaintenanceRoadHistoryID int    `json:"-"`         // รหัสโครงการ
	Path                     string `json:"path"`      // path folder
	FileName                 string `json:"file_name"` // ชื่อไฟล์
}

func (ral *MaintenanceRoadHistoryAttachment) TableName() string {
	return "maintenance_road_history_attachment"
}

func (ral *MaintenanceRoadHistoryAttachmentData) TableName() string {
	return "maintenance_road_history_attachment"
}
