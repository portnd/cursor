package mockdata

import (
	"time"

	"gitlab.com/mims-api-service/models"
)

var RoadMock = models.Road{
	Id:            1,
	Seq:           1,
	ParentRoadId:  nil,
	RoadLevel:     1,
	RoadCode:      "ABC123",
	IsActive:      true,
	RoadGroupId:   100,
	IsInit:        true,
	RoadSectionId: 1,
	CreatedBy:     1,
	CreatedAt:     time.Now(),
}

var parentRoadId = 1
var RoadChildMock = models.Road{
	Id:            2,
	Seq:           1,
	ParentRoadId:  &parentRoadId,
	RoadLevel:     2,
	RoadCode:      "ABC456",
	IsActive:      true,
	RoadGroupId:   100,
	IsInit:        true,
	RoadSectionId: 1,
	CreatedBy:     1,
	CreatedAt:     time.Now(),
}

var RoadByIdMock = models.RoadById{
	Road:            RoadMock,
	RoadInfo:        RoadInfoAddDataMock,
	RoadSurfaceIcon: []models.RoadSurfaceIcon{RoadSurfaceIconMock},
	RoadGeom:        []models.RoadGeom{RoadGeomMock},
	RoadSection:     RoadSectionByIdMock,
}

var RoadDataMock = models.RoadData{
	Road:            RoadMock,
	RoadInfo:        RoadInfoAddDataMock,
	RefSurface:      RefSurfaceRoadMock,
	RoadSurfaceIcon: []models.RoadSurfaceIcon{RoadSurfaceIconMock},
	RoadGeom:        []models.RoadGeom{RoadGeomMock},
	ChildRoads:      []models.ChildRoadData{ChildRoadDataMaock},
}

var ChildRoadDataMaock = models.ChildRoadData{
	Road:            RoadChildMock,
	RoadInfo:        RoadInfoAddDataMock,
	RefSurface:      RefSurfaceRoadMock,
	RoadSurfaceIcon: []models.RoadSurfaceIcon{RoadSurfaceIconMock},
	RoadGeom:        []models.RoadGeom{RoadGeomMock},
	ChildRoads:      []models.ChildRoadData{},
}

var ChildRoadDataFoundRefSurfaceMaock = models.ChildRoadData{
	Road:            RoadChildMock,
	RoadInfo:        RoadInfoAddDataMock,
	RefSurface:      RefSurfaceRoadForChildMock,
	RoadSurfaceIcon: []models.RoadSurfaceIcon{RoadSurfaceIconMock},
	RoadGeom:        []models.RoadGeom{RoadGeomMock},
	ChildRoads:      []models.ChildRoadData{},
}

var RoadListMock = models.RoadList{
	RoadGroup: RoadGroupMock,
	Sections:  []models.RoadSectionData{RoadSectionDataMock},
}
