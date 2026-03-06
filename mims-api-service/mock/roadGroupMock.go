package mockdata

import "gitlab.com/mims-api-service/models"

var RoadGroupMock = models.RoadGroup{
	Id:        1,
	Number:    "0007",
	Code:      "",
	Name:      "Road Group 1",
	ShortName: "Road G",
	KmStart:   0.0,
	KmEnd:     100.0,
	Distance:  100.0,
}
