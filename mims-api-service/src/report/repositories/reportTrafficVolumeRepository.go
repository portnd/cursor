package repositories

import (
	"gitlab.com/mims-api-service/models"

	_ "github.com/go-sql-driver/mysql"
)

func (r *Repository) GetReportTrafficVolume(roadIDs []int, year int) ([]models.VolumeAadt, error) {
	var data []models.VolumeAadt
	err := r.conn.Debug().
		Select(`
		veh1,
		veh2,
		veh3,
		total,
		surveyed_date,
        year,
		road_id
		`).
		Where("year = ?", year).
		Where("road_id IN(?) ", roadIDs).
		Where("status = ?", "A").
		Order("surveyed_date DESC").
		Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Repository) GetRoadDetailsByRoadSectionID(roadSectionIDs []string) ([]models.ReportTrafficVolumeHeader, error) {
	var resps []models.ReportTrafficVolumeHeader

	// Start constructing the query
	query := r.conn.Debug().
		Table("road").
		Select(`road.id AS road_id,
			road.road_section_id,
			road_group.number AS road_group_name,
			road_section.number AS road_section_name,
			CASE
			WHEN road_section.name_destination_th IS NOT NULL AND road_section.name_destination_th <> '' THEN
				CONCAT(road_section.name_origin_th , ' - ' ,road_section.name_destination_th)
			ELSE
				road_section.name_destination_th
			END AS road_name,
			road_section.km_start,
			road_section.km_end,
			road_section.distance AS total_km`).
		Joins("LEFT JOIN road_info ON road.id = road_info.road_id").
		Joins("LEFT JOIN road_group ON road.road_group_id = road_group.id").
		Joins("LEFT JOIN road_section ON road.road_section_id = road_section.id").
		Where("road_info.status = ?", "A").
		Where("road.is_active = ?", true)

	if len(roadSectionIDs) > 0 {
		query = query.Where("road.road_section_id = (?) ", roadSectionIDs)
	}

	// Execute the query
	err := query.Find(&resps).Error
	if err != nil {
		return nil, err
	}

	return resps, nil
}
