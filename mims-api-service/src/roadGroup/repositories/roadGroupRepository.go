package repositories

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/roadGroup/handlers"
	"gitlab.com/mims-api-service/src/roadGroup/usecases"

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

func (r *Repository) GetRoadGroup() ([]models.RoadGroup, error) {
	query := r.conn
	var roadGroup []models.RoadGroup
	if err := query.Order("id").Find(&roadGroup).Error; err != nil {
		return roadGroup, err
	}
	return roadGroup, nil
}

func (r *Repository) GetRoadByRoadGroupID(id int) ([]models.RoadInRoadGroup, error) {
	var result []models.RoadInRoadGroup
	query := r.conn
	query = query.Table("road").Select("road.id,parent_road_id as parent_id,road.ref_direction_id as direction_id,road_info.name as name").Joins("JOIN road_info on road.id = road_info.road_id")
	query = query.Where("road_group_id = ?", id)
	query = query.Where("road.is_active = ?", true).Order("id")
	err := query.Scan(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *Repository) GetRoadGroupByID(id int) (models.RoadGroupByID, error) {
	var result models.RoadGroupByID
	query := r.conn
	query = query.Table("road_group").Select("id,name").Where("id = ?", id)
	err := query.Scan(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *Repository) GetLaneByRoadID(id int) ([]models.RoadLane, error) {
	var results []models.RoadLane
	query := r.conn
	query = query.Table("road_geom").Select("lane_no ,km_start,km_end").Where("road_id = ?", id).Where("status = ?", "A").Order("lane_no")
	err := query.Scan(&results).Error
	if err != nil {
		return results, err
	}
	return results, nil
}
