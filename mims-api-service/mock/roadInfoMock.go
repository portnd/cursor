package mockdata

import (
	"time"

	"gitlab.com/mims-api-service/models"
)

var RoadInfoMock = models.RoadInfo{
	Id:                      1,
	RoadId:                  1,
	Year:                    nil,
	RefDirectionId:          1,
	Name:                    "Sample Road",
	KmStart:                 0.0,
	KmEnd:                   10.5,
	TheGeom:                 "LINESTRING(100.521147 13.625611)",
	Revision:                1,
	Status:                  "A",
	RampId:                  "R001",
	RoadColorCode:           "#FF0000",
	CreatedBy:               123,
	CreatedAt:               time.Now(),
	UpdatedBy:               123,
	UpdatedAt:               time.Now(),
	Remark:                  "Sample Remark",
	RefRoadTypeID:           1,
	CenterLaneShapeFilePath: "/path/to/center_lane_shape_file",
	CenterLineShapeFilePath: "/path/to/center_line_shape_file",
}

var RoadInfoAddDataMock = models.RoadInfoAddData{
	RoadInfo:            RoadInfoMock,
	OriginToDestination: "City A - City B",
	RoadCode:            "00071212",
	ResponsibleCode:     "RESP456",
	KmRange:             "0+000 - 2+000",
	User:                User,
	RefRoadType:         RefRoadType,
}
