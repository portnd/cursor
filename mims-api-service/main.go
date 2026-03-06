package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"gitlab.com/mims-api-service/databases"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/routes"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var err error

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	logs.InitLogger()
	logs.InitLoggerRequest()
	defer logs.CloseLogReq()
	defer logs.Close()

	DB_HOST := os.Getenv("DB_HOST")
	DB_DATABASE := os.Getenv("DB_DATABASE")
	DB_USERNAME := os.Getenv("DB_USERNAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		DB_HOST, DB_USERNAME, DB_PASSWORD, DB_DATABASE, DB_PORT)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logs.Error(err)
		panic(err)
	}

	databases.DB = db

	sqlDB, err := databases.DB.DB()
	if err != nil {
		log.Fatal("Error getting sqlDB from GORM: ", err)
	}
	defer sqlDB.Close()

	uri := os.Getenv("MONGODB_URI")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	databases.MongoDb, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer databases.MongoDb.Disconnect(context.TODO())

	// run the migrations: todo struct
	// databases.DB.AutoMigrate(&models.User{})

	//setup routes
	r := routes.SetupRouter()
	// running
	r.Run(":8080")
}
