package repositories

import (
	_ "github.com/go-sql-driver/mysql"
	models "gitlab.com/mims-api-service/models"
	"gorm.io/gorm"
)

func (t *roadConditionRepository) GetRoadSurfaceByRoadIDByLaneNo(roadID, laneNo int) ([]models.RoadSurfaceData, error) {
	var roadSurface []models.RoadSurfaceData
	query := t.conn
	// query = query.Preload("RoadSurface")
	query = query.Preload("RoadSurfaceLane")
	query = query.Preload("RoadSurfaceLane", func(db *gorm.DB) *gorm.DB {
		return db.Where("lane_no = ?", laneNo)
	})
	query = query.Where("road_id = ?", roadID)
	query = query.Where("status = ?", "A")
	if err := query.Order("id asc").Find(&roadSurface).Error; err != nil {
		return roadSurface, err
	}
	return roadSurface, nil
}

func (t *roadConditionRepository) GetRefSurface() (map[int]string, error) {
	var refSurfaces []models.RefSurface
	res := make(map[int]string)
	query := t.conn
	if err := query.Order("id asc").Find(&refSurfaces).Error; err != nil {
		return res, err
	}

	for _, item := range refSurfaces {
		if item.SurfaceGroup == "Asphalt" {
			res[item.ID] = "AC"
		} else {
			res[item.ID] = "CC"
		}

	}
	return res, nil
}
