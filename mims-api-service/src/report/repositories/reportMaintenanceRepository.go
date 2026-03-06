package repositories

import (
	"gitlab.com/mims-api-service/models"

	_ "github.com/go-sql-driver/mysql"
)

func (r *Repository) GetReportMaintenance(yearStart, yearEnd int, roadSectionId int) ([]models.DataGetMaintenance, []models.MethodMaintenance, error) {
	var resp []models.DataGetMaintenance
	err := r.conn.
		Table("maintenance").
		Select(`DISTINCT 
		road.id as road_id, 
		road_info.name, 
		road_info.km_start, 
		road_info.km_end,
		road_group.number as road_group_number,
		road_section.number as road_section_number, 
		CASE
        WHEN road_section.name_destination_th IS NOT NULL AND road_section.name_destination_th <> '' THEN
			CONCAT(road_section.name_origin_th , ' - ' ,road_section.name_destination_th)
        ELSE
			road_section.name_destination_th
    	END AS road_main_name,
		road_section.distance as distance,
		road_section.km_start as sec_km_start,
		road_section.km_end as sec_km_end
		`).
		Joins("JOIN maintenance_road ON maintenance.ID = maintenance_road.maintenance_id").
		Joins("JOIN road ON maintenance_road.road_id = road.id").
		Joins("JOIN road_section ON road.road_section_id = road_section.id").
		Joins("JOIN road_group ON road.road_group_id = road_group.id").
		Joins("JOIN road_info ON road.id = road_info.road_id").
		Where("budget_year BETWEEN ? AND ?", yearStart, yearEnd).
		Where("road.road_section_id = ?", roadSectionId).
		Order("road.id ASC").
		Find(&resp).Error
	if err != nil {
		return nil, nil, err
	}

	var detail []models.MethodMaintenance
	err = r.conn.
		Table("maintenance").
		Select(`
        maintenance_road.road_id,
        maintenance.budget_year as year,
        maintenance_road.lane_no as lane, 
        maintenance_road.km_start, 
        maintenance_road.km_end,
        ref_criteria_method."name" as method
		`).
		Joins("JOIN maintenance_road ON maintenance.ID = maintenance_road.maintenance_id").
		Joins("JOIN road ON maintenance_road.road_id = road.id").
		Joins("JOIN ref_criteria_method ON ref_criteria_method.id = maintenance_road.maintenance_method_id").
		Where("road.road_section_id = ?", roadSectionId).
		Where("maintenance.status = 'A'").
		Order("maintenance_road.road_id ASC, maintenance.budget_year ASC, maintenance_road.lane_no ASC").
		Find(&detail).Error
	if err != nil {
		return nil, nil, err
	}

	return resp, detail, nil
}

func (r *Repository) GetMultiRoadInfo(roadSectionId int) ([]models.MultiRoadInfo, error) {
	var resp []models.MultiRoadInfo
	err := r.conn.
		Table("maintenance").
		Select(`DISTINCT
		road.id,  
		road.road_code, 
		road_info.name, 
		road_info.km_start, 
		road_info.km_end
		`).
		Joins("JOIN maintenance_road ON maintenance.ID = maintenance_road.maintenance_id").
		Joins("JOIN road ON maintenance_road.road_id = road.id").
		Joins("JOIN road_info ON road.id = road_info.road_id").
		Where("road.road_section_id = ?", roadSectionId).
		Order("road.id ASC").
		Find(&resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}
