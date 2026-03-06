package models

import "time"

// Todo ...
type Role struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedBy int       `json:"-"`
	UpdatedBy int       `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	IsActive  bool      `json:"-"`
}

// TableName use to specific table
func (b *Role) TableName() string {
	return "role"
}
