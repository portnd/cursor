package mockdata

import (
	"time"

	"gitlab.com/mims-api-service/models"
)

var RoadGeomMock = models.RoadGeom{
	Id:        1,
	RoadId:    1001,
	LaneNo:    1,
	KmStart:   0.0,
	KmEnd:     10.5,
	TheGeom:   "LINESTRING(100.521147 13.625611)",
	Revision:  1,
	Status:    "A",
	Remark:    "Sample Remark",
	CreatedBy: 123,
	CreatedAt: time.Now(),
	UpdatedBy: 123,
	UpdatedAt: time.Now(),
}
