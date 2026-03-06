package domains

import (
	"gitlab.com/mims-api-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type UseCase interface {
	GetSectionGeom() error
	GetRoadLatest() error
	MatchData() (interface{}, error)
}

type Repository interface {
	InsertSectionGeom(data []interface{}) error
	InsertRoadLatest(data []interface{}) error
	GetSectionGeomWithFilter(filter primitive.D) ([]models.Item, error)
	GetRoadLatest(filter primitive.D) ([]models.RoadLatest, error)

	GetRefHirs(conn *gorm.DB) ([]models.RefHris, error)
	GetRoadSection(conn *gorm.DB) ([]models.RoadSection, error)
	GetRoadGroup(conn *gorm.DB) ([]models.RoadGroup, error)
	InsertRoadGroup(roadGroup models.InsertRoadGroup, conn *gorm.DB) error
	InsertRoadSection(roadSection models.InsertRoadSection, conn *gorm.DB) error
	UpdateRoadGroupNameByNumber(number, name, shortName string, conn *gorm.DB) error

	StartTransSection() *gorm.DB
	RollBack(tx *gorm.DB) error
	Commit(tx *gorm.DB) error
}
