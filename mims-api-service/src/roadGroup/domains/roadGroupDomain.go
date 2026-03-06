package domains

import (
	"gitlab.com/mims-api-service/models"
)

// business logic
type UseCase interface {
	GetRoadGroup() (interface{}, error)
	GetRoadGroupByID(id int) (models.RoadGroupByID, error)
}

// อะไรเชื่อมต่อกับ DB
type Repository interface {
	GetRoadGroup() ([]models.RoadGroup, error)
	GetRoadByRoadGroupID(id int) ([]models.RoadInRoadGroup, error)
	GetRoadGroupByID(id int) (models.RoadGroupByID, error)
	GetLaneByRoadID(id int) ([]models.RoadLane, error)
}
