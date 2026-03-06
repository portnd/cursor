package models

import "time"

type RoadOwner struct {
	ID            int
	RoadID        int
	Year          int
	RefOwnerID    int
	EffectiveDate time.Time
	IsActive      bool
}

func (ro *RoadOwner) TableName() string {
	return "road_owner"
}
