package models

import "time"

// Todo ...
type UserRole struct {
	Id        uint      `gorm:"primary_key" json:"id"`
	UserId    int       `json:"user_id"`
	RoleID    int       `json:"role_id"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName use to specific table
func (b *UserRole) TableName() string {
	return "user_role"
}
