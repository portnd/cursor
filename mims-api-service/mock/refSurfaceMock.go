package mockdata

import "gitlab.com/mims-api-service/models"

var RefSurfaceRoadMock = models.RefSurfaceRoad{
	Id:           1,
	RoadId:       1,
	RefSurfaceId: []int{1, 2, 3},
}

var RefSurfaceRoadForChildMock = models.RefSurfaceRoad{
	Id:           1,
	RoadId:       1,
	RefSurfaceId: []int{5, 6, 7},
}
