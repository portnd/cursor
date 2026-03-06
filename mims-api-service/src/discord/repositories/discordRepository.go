package repositories

import (
	"context"
	"os"

	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/discord/handlers"
	"gitlab.com/mims-api-service/src/discord/usecases"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"gorm.io/gorm"
)

type Repository struct {
	conn    *gorm.DB
	connMon *mongo.Client
}

func NewRepository(conn *gorm.DB, connMon *mongo.Client) *handlers.Handler {
	useCase := usecases.NewUsecase(&Repository{conn, connMon})
	handler := handlers.NewHandler(useCase)
	return handler
}

///////////////// MongoDB /////////////////

func (r *Repository) GetHrisAndHsmsLogWithFilter(filter primitive.D, collectionDb string) ([]models.MongoDbLog, error) {
	mongoDb := r.connMon
	collection := mongoDb.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv(collectionDb))

	var mongoDbLog []models.MongoDbLog
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return mongoDbLog, err
	}

	if err = cursor.All(context.TODO(), &mongoDbLog); err != nil {
		return mongoDbLog, err
	}

	return mongoDbLog, nil
}
