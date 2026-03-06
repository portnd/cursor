package helpers

import (
	"context"
	"os"
	"time"

	"gitlab.com/mims-api-service/databases"
)

type MongoDbInsert struct {
	IsSuccess bool      `json:"is_success"`
	DateTime  time.Time `json:"date_time"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
}

func MongoDbLog(title, message, collectionMongoDb string, isSuccess bool) error {

	mongoDb := databases.MongoDb
	collection := mongoDb.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv(collectionMongoDb))

	var mongoDbInsert MongoDbInsert
	mongoDbInsert.IsSuccess = isSuccess
	mongoDbInsert.DateTime = time.Now()
	mongoDbInsert.Title = title
	mongoDbInsert.Message = message

	_, err := collection.InsertOne(context.TODO(), mongoDbInsert)
	if err != nil {
		return err
	}

	return nil

}
