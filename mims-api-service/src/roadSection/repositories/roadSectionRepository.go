package repositories

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/roadSection/handlers"
	"gitlab.com/mims-api-service/src/roadSection/usecases"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type Repository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewRepositoryHandler(conn *gorm.DB) *handlers.Handler {
	useCase := usecases.NewUseCase(&Repository{conn})
	handler := handlers.NewHandler(useCase)
	return handler
}

func (r *Repository) GetRoadSection(roadGroupId *int) ([]models.RoadSection, error) {
	query := r.conn
	var roadSection []models.RoadSection

	query = query.Preload("RefDivision", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, division_code, name, name_eng, st_astext(the_geom) as the_geom")
	})

	query = query.Preload("RefDistrict", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, district_code, name, name_eng, st_astext(the_geom) as the_geom")
	})

	query = query.Preload("RefDepot", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, depot_code, name, st_astext(the_geom) as the_geom")
	})

	if roadGroupId != nil {
		query.Where("road_group_id = ?", *roadGroupId)
	}

	if err := query.Order("id").Find(&roadSection).Error; err != nil {
		return nil, err
	}
	return roadSection, nil
}

func (r *Repository) GetRoadSectionByID(id int) (*models.RoadSection, error) {
	var result models.RoadSection
	query := r.conn
	query = query.Table("road_group").Select("id,name").Where("id = ?", id)
	err := query.Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
