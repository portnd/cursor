package repositories

import (
	"gitlab.com/mims-api-service/models"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

func (r *Repository) GetReportRoad(roadGroupIDs []int) ([]models.RoadListReport, error) {
	var roadList []models.RoadListReport
	query := r.conn
	query = query.Where("id IN ?", roadGroupIDs)
	query = query.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		db = db.Select("road_section.id,road_section.road_group_id,road_section.number,road_section.name_origin_th,road_section.name_destination_th,road_section.name_origin_en,road_section.name_destination_en,road_section.km_start, road_section.km_end, road_section.distance,road_section.ref_division_code,road_section.ref_district_code,road_section.ref_depot_code, ARRAY_AGG(ref_province.name) AS province").
			Joins("JOIN ref_province ON ref_province.province_code = ANY(road_section.province_code)").Group("road_section.id,road_section.road_group_id,road_section.number,road_section.name_origin_th,road_section.name_destination_th,road_section.name_origin_en,road_section.name_destination_en,road_section.km_start, road_section.km_end, road_section.distance,road_section.ref_division_code,road_section.ref_district_code,road_section.ref_depot_code")
		return db.Order("road_section.number ASC")
	})

	// query = query.Preload("Sections.RefDivision", func(db *gorm.DB) *gorm.DB {
	// 	return db.Select("id, division_code, name, name_en")
	// })

	// query = query.Preload("Sections.RefDistrict", func(db *gorm.DB) *gorm.DB {
	// 	return db.Select("id, district_code, name, name_en")
	// })

	query = query.Preload("Sections.RefDepot", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, depot_code, name")
	})

	// query = query.Preload("Sections.Roads", func(db *gorm.DB) *gorm.DB {
	// 	db = db.Where("road.is_active = true").Order("road.id ASC")
	// 	return db
	// })

	// query = query.Preload("Sections.Roads.RoadInfo", func(db *gorm.DB) *gorm.DB {
	// 	db = db.Select("id,road_id, year, name, km_start, km_end,   ref_direction_id, revision, status, ramp_id,road_color_code").Where("status = ?", "A")

	// 	return db
	// })

	// query = query.Preload("Sections.Roads.RefSurface", func(db *gorm.DB) *gorm.DB {
	// 	db = db.Select("road_surface.road_id, ARRAY_AGG(DISTINCT road_surface_lane.ref_surface_id) as ref_surface_id").
	// 		Joins("LEFT JOIN road_surface_lane ON road_surface_lane.road_surface_id = road_surface.id").
	// 		Where("road_surface.status = ? and road_surface_lane.ref_surface_id <> ?", "A", 0).
	// 		Group("road_surface.road_id")
	// 	return db
	// })

	// query = query.Preload("Sections.Roads.RoadSurfaceIcon", func(db *gorm.DB) *gorm.DB {
	// 	db = db.Select("DISTINCT ON (ref_surface.surface_group,road_surface.road_id) ref_surface.surface_group, ref_surface.id , road_surface.road_id, ref_surface.surface_group as name, CASE WHEN ref_surface.surface_group = 'Concrete' THEN '#398BF7' ELSE '#7460EE' END AS color_code").
	// 		Joins("RIGHT JOIN road_surface_lane ON road_surface.id = road_surface_lane.road_surface_id").
	// 		Joins("LEFT JOIN ref_surface ON ref_surface.id = road_surface_lane.ref_surface_id").
	// 		Where("road_surface.status = ? and road_surface_lane.ref_surface_id <> ?", "A", 0)

	// 	return db
	// })

	// query = query.Preload("Sections.Roads.RoadInfo.RefRoadType")

	// query = query.Preload("Sections.Roads.RoadGeom", func(db *gorm.DB) *gorm.DB {
	// 	db = db.Select("id,road_id, lane_no, km_start,km_end,ST_ASTEXT(ST_FORCE2D(the_geom)) as the_geom,revision,status,remark,created_by,created_at,updated_by,updated_at").Where("status = ?", "A").Order("revision DESC").Order("revision")
	// 	return db
	// })

	// query = query.Preload("Sections.Roads.ChildRoads", func(db *gorm.DB) *gorm.DB {
	// 	db = db.Where("is_active = true")
	// 	return db
	// })

	// query = query.Preload("Sections.Roads.ChildRoads.RefSurface", func(db *gorm.DB) *gorm.DB {
	// 	db = db.Select("road_surface.road_id, ARRAY_AGG(DISTINCT road_surface_lane.ref_surface_id) as ref_surface_id").
	// 		Joins("LEFT JOIN road_surface_lane ON road_surface_lane.road_surface_id = road_surface.id").
	// 		Where("road_surface.status = ? and road_surface_lane.ref_surface_id <> ?", "A", 0).
	// 		Group("road_surface.road_id")
	// 	return db
	// })

	// query = query.Preload("Sections.Roads.ChildRoads.RoadSurfaceIcon", func(db *gorm.DB) *gorm.DB {
	// 	db = db.Select("DISTINCT ON (ref_surface.surface_group,road_surface.road_id) ref_surface.surface_group, ref_surface.id , road_surface.road_id, ref_surface.surface_group as name, CASE WHEN ref_surface.surface_group = 'Concrete' THEN '#398BF7' ELSE '#7460EE' END AS color_code").
	// 		Joins("RIGHT JOIN road_surface_lane ON road_surface.id = road_surface_lane.road_surface_id").
	// 		Joins("LEFT JOIN ref_surface ON ref_surface.id = road_surface_lane.ref_surface_id").
	// 		Where("road_surface.status = ? and road_surface_lane.ref_surface_id <> ?", "A", 0)
	// 	return db
	// })

	// query = query.Preload("Sections.Roads.ChildRoads.RoadGeom", func(db *gorm.DB) *gorm.DB {
	// 	db = db.Select("id,road_id, lane_no, km_start,km_end,ST_ASTEXT(ST_FORCE2D(the_geom)) as the_geom,revision,status,remark,created_by,created_at,updated_by,updated_at").Where("status = ?", "A").Order("revision DESC").Order("revision")
	// 	return db
	// })

	// query = query.Preload("Sections.Roads.ChildRoads.RoadInfo", func(db *gorm.DB) *gorm.DB {
	// 	db = db.Select("id,road_id, year, name, km_start, km_end,   ref_direction_id, revision, status, ramp_id,road_color_code, created_by,created_at,updated_by,updated_at,ref_road_type_id").Where("status = ?", "A")
	// 	return db
	// })

	// query = query.Preload("Sections.Roads.ChildRoads.RoadInfo.RefRoadType")

	// query = query.Preload("Sections.Roads.ChildRoads.ChildRoads")

	if err := query.Find(&roadList).Error; err != nil {
		return roadList, err
	}

	return roadList, nil
}
