package repositories

import (
	"context"
	"os"

	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/hris/handlers"
	"gitlab.com/mims-api-service/src/hris/usecases"
	"go.mongodb.org/mongo-driver/bson"
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

func (r *Repository) InsertSectionGeom(data []interface{}) error {
	mongoDb := r.connMon
	collection := mongoDb.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_HRIS2_SECTION_GEOM_COLLECTION"))

	filter := bson.D{{"is_latested", true}}
	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}

	_, err = collection.InsertMany(context.TODO(), data)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) InsertRoadLatest(data []interface{}) error {
	mongoDb := r.connMon
	collection := mongoDb.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_HRIS2_ROAD_LATEST_COLLECTION"))

	filter := bson.D{{"is_latested", true}}
	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}

	_, err = collection.InsertMany(context.TODO(), data)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetSectionGeomWithFilter(filter primitive.D) ([]models.Item, error) {
	mongoDb := r.connMon
	collection := mongoDb.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_HRIS2_SECTION_GEOM_COLLECTION"))

	var sectionGeom []models.Item
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return sectionGeom, err
	}

	if err = cursor.All(context.TODO(), &sectionGeom); err != nil {
		return sectionGeom, err
	}

	return sectionGeom, nil
}

func (r *Repository) GetRoadLatest(filter primitive.D) ([]models.RoadLatest, error) {
	mongoDb := r.connMon
	collection := mongoDb.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_HRIS2_ROAD_LATEST_COLLECTION"))

	var roadLatest []models.RoadLatest
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return roadLatest, err
	}

	if err = cursor.All(context.TODO(), &roadLatest); err != nil {
		return roadLatest, err
	}

	return roadLatest, nil
}

///////////////// Postgresql /////////////////

func (r *Repository) GetRefHirs(conn *gorm.DB) ([]models.RefHris, error) {
	var refHris []models.RefHris
	err := conn.Where("status = true AND is_deleted = false").Find(&refHris).Error
	if err != nil {
		conn.Rollback()
		return refHris, err
	}
	return refHris, nil
}

func (r *Repository) GetRoadGroup(conn *gorm.DB) ([]models.RoadGroup, error) {
	var roadGroup []models.RoadGroup
	err := conn.Find(&roadGroup).Error
	if err != nil {
		conn.Rollback()
		return roadGroup, err
	}
	return roadGroup, nil
}

func (r *Repository) GetRoadSection(conn *gorm.DB) ([]models.RoadSection, error) {
	var roadSection []models.RoadSection
	err := conn.Find(&roadSection).Error
	if err != nil {
		conn.Rollback()
		return roadSection, err
	}
	return roadSection, nil
}

func (r *Repository) InsertRoadGroup(roadGroup models.InsertRoadGroup, conn *gorm.DB) error {
	if err := conn.Save(&roadGroup).Error; err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *Repository) InsertRoadSection(roadSection models.InsertRoadSection, conn *gorm.DB) error {
	if err := conn.Save(&roadSection).Error; err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *Repository) UpdateRoadGroupNameByNumber(number, name, shortName string, conn *gorm.DB) error {
	err := conn.Model(&models.InsertRoadGroup{}).Where("number = ?", number).Update("name", name).Update("short_name", shortName).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

///////////////// Postgresql TransSection /////////////////

func (r *Repository) StartTransSection() *gorm.DB {
	tx := r.conn.Begin()
	return tx
}

func (r *Repository) RollBack(tx *gorm.DB) error {
	tx.Rollback()
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) Commit(tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
