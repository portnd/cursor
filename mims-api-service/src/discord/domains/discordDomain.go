package domains

import (
	"gitlab.com/mims-api-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UseCase interface {
	GetHrisAndHsmsLog() error
}

type Repository interface {
	GetHrisAndHsmsLogWithFilter(filter primitive.D, collectionDb string) ([]models.MongoDbLog, error)
}
