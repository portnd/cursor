package domains

import (
	"gitlab.com/mims-api-service/models"
)

// business logic
type UseCase interface {
	GetRoadSection(*int) ([]models.RoadSection, error)
	GetRoadSectionByID(id int) (*models.RoadSection, error)
}

// อะไรเชื่อมต่อกับ DB
type Repository interface {
	GetRoadSection(*int) ([]models.RoadSection, error)
	GetRoadSectionByID(id int) (*models.RoadSection, error)
}
