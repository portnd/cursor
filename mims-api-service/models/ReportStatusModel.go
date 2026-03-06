package models

import "time"

type ReportStatus struct {
	Id        int       `json:"id"`
	Path      string    `json:"path"`
	IsFinish  bool      `json:"is_finish"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName use to specific table
func (b *ReportStatus) TableName() string {
	return "report_status"
}
