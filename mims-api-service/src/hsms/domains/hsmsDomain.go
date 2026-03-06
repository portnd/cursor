package domains

import (
	"gitlab.com/mims-api-service/models"
	"gorm.io/gorm"
)

// business logic
type UseCase interface {
	GetHsmsBridge() (interface{}, error)
	GetHsmsGuard() (interface{}, error)
	GetHsmsInterchange() (interface{}, error)
	GetHsmsIntersection() (interface{}, error)
	GetHsmsStreetlight() (interface{}, error)
	GetHsmsRailwaycrossing() (interface{}, error)
	GetHsmsTrafficlight() (interface{}, error)
	GetHsmsUturnbridge() (interface{}, error)
}

// อะไรเชื่อมต่อกับ DB
type Repository interface {
	GetHsmsBridge() ([]models.Hsms01Bridge, error)
	GetHsmsGuard() ([]models.Hsms01Guard, error)
	GetHsmsInterchange() ([]models.Hsms01Interchange, error)
	GetHsmsIntersection() ([]models.Hsms01Intersection, error)
	GetHsmsStreetlight() ([]models.Hsms01Light, error)
	GetHsmsRailwaycrossing() ([]models.Hsms01Railwaycrossing, error)
	GetHsmsTrafficlight() ([]models.Hsms01Signal, error)
	GetHsmsUturnbridge() ([]models.Hsms01Uturnbridge, error)

	StartTransSection() *gorm.DB
	RollBack(tx *gorm.DB) error
	Commit(tx *gorm.DB) error
}
