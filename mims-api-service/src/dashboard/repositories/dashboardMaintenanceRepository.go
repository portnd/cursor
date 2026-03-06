package repositories

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/src/dashboard/handlers"
	"gitlab.com/mims-api-service/src/dashboard/usecases"
	"gorm.io/gorm"
)

type repositoryMaintenance struct {
	conn *gorm.DB
}

// init Repository Handler
func NewRepositoryHandlerMaintenance(conn *gorm.DB) *handlers.HandlerMaintenance {
	usecase := usecases.NewUsecaseMaintenance(&repositoryMaintenance{conn})
	handler := handlers.NewHandlerMaintenance(usecase)
	return handler
}

func (r *repositoryMaintenance) GetRoadGroup() ([]models.RoadGroup, error) {
	query := r.conn
	var roadGroup []models.RoadGroup
	err := query.Find(&roadGroup).Error
	if err != nil {
		return roadGroup, err
	}

	return roadGroup, nil
}

func (r *repositoryMaintenance) GetMaintenance(limit int, roadIDs []int, depotCodes []string, filter requests.MaintenanceDashboard) ([]models.MaintenancePreloadForDashboard, error) {
	query := r.conn
	var maintainaintenancePreloadForDashboard []models.MaintenancePreloadForDashboard

	query = query.Preload("MaintenanceRoads", func(db *gorm.DB) *gorm.DB {
		db = db.Select("ST_AsGeoJSON(maintenance_road.the_geom) as the_geom_json, maintenance_road.*")
		db = db.Where("status = 'A'")
		if len(roadIDs) != 0 {
			db = db.Where("road_id IN (?)", roadIDs)
		}

		if filter.KmStart != 0.0 {
			db = db.Where("km_start >= ?", filter.KmStart)
		}

		if filter.KmEnd != 0.0 {
			db = db.Where("km_end <= ?", filter.KmEnd)
		}
		return db
	})

	query = query.Preload("MaintenanceRoads.Road", func(db *gorm.DB) *gorm.DB {
		db = db.Where("is_active = true")
		return db
	})

	query = query.Preload("MaintenanceRoads.Road.RoadInfo", func(db *gorm.DB) *gorm.DB {
		db = db.Select("ST_AsGeoJSON(road_info.the_geom) as the_geom_json, road_info.*")
		db = db.Where("status = 'A'")
		return db
	})

	query = query.Preload("MaintenanceRoads.Road.RoadSection")

	query = query.Preload("MaintenanceRoads.Road.RoadSection.RefDepot")

	query = query.Where("status = 'A'")
	if filter.Year != "" {
		query = query.Where("budget_year = ?", filter.Year)
	}
	if len(depotCodes) != 0 {
		query = query.Where("ref_depot_code IN (?)", depotCodes)
	}
	query = query.Order("budget_maintenance DESC")

	if limit != 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&maintainaintenancePreloadForDashboard).Error
	if err != nil {
		return maintainaintenancePreloadForDashboard, err
	}

	return maintainaintenancePreloadForDashboard, nil
}
