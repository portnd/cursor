package repositories

import (
	"fmt"
	"strings"
	"time"

	models "gitlab.com/mims-api-service/models"
	requests "gitlab.com/mims-api-service/requests"
	responses "gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/road/handlers"
	"gitlab.com/mims-api-service/src/road/usecases"

	_ "github.com/go-sql-driver/mysql"
	helperFile "gitlab.com/mims-api-service/helpers/file"
	servicesDB "gitlab.com/mims-api-service/services/database"
	"gorm.io/gorm"
)

type roadRepository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewRoadRepositoryHandler(conn *gorm.DB) *handlers.RoadHandler {
	servicesDB := servicesDB.NewServicesDatabase(conn)
	helperFile := helperFile.NewHelpersFile()
	useCase := usecases.NewRoadUseCase(&roadRepository{conn}, helperFile, servicesDB)
	handler := handlers.NewRoadHandler(useCase)
	return handler
}

// /////////////////////////
func (t *roadRepository) GetRole(userId uint) ([]models.UserRole, error) {
	var userRole []models.UserRole
	if err := t.conn.Where("user_id = ?", userId).Find(&userRole).Error; err != nil {
		return userRole, err
	}
	return userRole, nil
}

func (t *roadRepository) GetAccessControl(roles []int) ([]models.AccessControl, error) {
	var accessControl []models.AccessControl
	query := t.conn
	query = query.Joins("JOIN role_access_control on access_control.id = role_access_control.access_control_id")
	err := query.Find(&accessControl).Error
	if err != nil {
		fmt.Println(err)
	}
	return accessControl, nil
}

func (t *roadRepository) GetUserById(userId uint) (models.Users, error) {
	var user models.Users
	if err := t.conn.Where("id = ?", userId).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (t *roadRepository) GetUserDepartmentById(userId int) (models.UserDepartment, error) {
	var user models.UserDepartment
	if err := t.conn.Where("id = ?", userId).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (t *roadRepository) GetRoadDetailMenu(accessKeys []string, assetType string) ([]responses.RoadMenuData, error) {
	var data []responses.RoadMenuData
	// if helpers.HasPermission([]string{"road_in_asset_manage_data", "approve_road_in_asset_access", "road_out_asset_manage_data", "approve_road_out_asset_access"}, accessKeys) {
	query := t.conn
	query = query.Select("DISTINCT ref_asset.id as group_id, ref_asset.name as group_name, ref_asset_table.id as asset_id, ref_asset_table.table_label as asset_name, ref_asset_table.geom_type as geom_type, ref_asset_table.is_in_road as is_in_road, ref_asset_table.icon_filepath as icon_filepath,	ref_asset.seq AS a_seq, ref_asset_table.seq AS as_seq, ref_asset_table.is_active as is_active")
	query = query.Table("ref_asset")
	query = query.Joins("JOIN ref_asset_table on ref_asset.id = ref_asset_table.ref_asset_id")
	// query = query.Joins("JOIN ref_asset_table_staff on ref_asset_table.id = ref_asset_table_staff.ref_asset_table_id")

	// if !helpers.HasPermission([]string{"view_all_road_in_asset"}, accessKeys) {
	// 	query = query.Where("ref_asset_table_staff.ref_department_id = ?", departmentId)
	// }

	// if assetType == "assetin" {
	// query = query.Where("ref_asset_table.is_in_road = ?", true)
	// } else {
	// 	query = query.Where("ref_asset_table.is_in_road = ?", false)
	// }

	// helpers.PrintlnJson(accessKeys)
	// permissions = []string{"road_out_asset_manage_data", "road_out_asset_access"}
	// if helpers.HasPermission(permissions, accessKeys) {
	// 	query = query.Or("ref_asset_table.is_in_road = ?", false)
	// }

	// permissions = []string{"road_in_asset_manage_data", "road_out_asset_manage_data"}
	// if helpers.HasPermission(permissions, accessKeys) {
	// 	query = query.Where("ref_asset_table_staff.is_approver = ?", true)
	// 	query = query.Where("ref_asset_table_staff.ref_department_id = ?", departmentId)
	// }

	// permissions = []string{"road_in_asset_access", "road_out_asset_access"}
	// if helpers.HasPermission(permissions, accessKeys) {
	// 	query = query.Where("ref_asset_table_staff.is_approver = ?", false)
	// 	query = query.Where("ref_asset_table_staff.ref_department_id = ?", departmentId)
	// }

	query = query.Where("ref_asset_table.is_active = ?", true)
	query = query.Order("ref_asset.seq")
	query = query.Order("ref_asset_table.seq")
	err := query.Find(&data).Error
	if err != nil {
		return data, err
	}
	// }
	return data, nil
}

func (t *roadRepository) GetRoadGroupList(prams requests.RoadPrams, roadIDs []int, isAllData, isOwnerData bool, depotCode string) ([]models.RoadList, error) {
	var roadList []models.RoadList
	query := t.conn
	if len(prams.RoadGroupId) != 0 {
		query = query.Where("id IN ?", prams.RoadGroupId)
	}

	query = query.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		db = db.Select("road_section.id,road_section.road_group_id,road_section.number,road_section.name_origin_th,road_section.name_destination_th,road_section.name_origin_en,road_section.name_destination_en,road_section.km_start, road_section.km_end, road_section.distance,road_section.ref_division_code,road_section.ref_district_code,road_section.ref_depot_code, ARRAY_AGG(ref_province.name) AS province").
			Joins("JOIN ref_province ON ref_province.province_code = ANY(road_section.province_code)").Group("road_section.id,road_section.road_group_id,road_section.number,road_section.name_origin_th,road_section.name_destination_th,road_section.name_origin_en,road_section.name_destination_en,road_section.km_start, road_section.km_end, road_section.distance,road_section.ref_division_code,road_section.ref_district_code,road_section.ref_depot_code")

		if len(prams.RoadSectionId) != 0 {
			db = db.Where("road_section.id IN ?", prams.RoadSectionId)
		}

		if len(prams.DepotCode) != 0 {
			db = db.Where("road_section.ref_depot_code IN ?", prams.DepotCode)
		}
		if isOwnerData && !isAllData {
			db = db.Where("road_section.ref_depot_code = ?", depotCode)
		}
		if !isOwnerData && !isAllData {
			db = db.Where("road_section.ref_depot_code = ?", "no_data")
		}
		// if !isAllData && !isOwnerData {
		// 	continue
		// }
		// if prams.KmStart != nil && prams.KmEnd != nil {
		// 	db = db.Where("road_section.km_start >= ? and road_section.km_end <= ?", *prams.KmStart, *prams.KmEnd)
		// } else if prams.KmStart != nil {
		// 	db = db.Where("road_section.km_start >= ?", *prams.KmStart)
		// } else if prams.KmEnd != nil {
		// 	db = db.Where("road_section.km_end <= ?", *prams.KmEnd)
		// }

		return db.Order("road_section.number ASC")
	})

	// query = query.Preload("Sections.RefProvince")

	// query = query.Preload("Sections.RefProvince.Data", func(db *gorm.DB) *gorm.DB {
	// 	return db.Select("id, province_code,node ,name, name_en, region, region_name ,st_astext(the_geom) as the_geom, status")
	// })

	query = query.Preload("Sections.RefDivision", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, division_code, name, name_en, st_astext(the_geom) as the_geom")
	})

	query = query.Preload("Sections.RefDistrict", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, district_code, name, name_en, st_astext(the_geom) as the_geom")
	})

	query = query.Preload("Sections.RefDepot", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, depot_code, name, st_astext(the_geom) as the_geom")
	})

	//and road_condition."year" = date_part('year', now())
	//and road_retro_reflectivity."year" = date_part('year', now())
	//and road_damage."year" = date_part('year', now())
	query = query.Preload("Sections.Roads", func(db *gorm.DB) *gorm.DB {
		db = db.Select(`road.*, 
		CASE 
			WHEN count(road_condition.id) > 0 THEN true
			ELSE false 
		END AS condition_status,
		CASE 
			WHEN count(road_retro_reflectivity.id) > 0 THEN true 
			ELSE false 
		END AS retro_status,
		CASE 
			WHEN count(road_damage.id) > 0 THEN true 
			ELSE false 
		END AS damage_status,
		-- Color Logic based on current year or fallback
		CASE 
			WHEN count(road_condition.id) > 0 AND MAX(road_condition."year") = date_part('year', now()) THEN '#F27A0D' 
			ELSE '#C8C8C8' 
		END AS condition_status_color,
		CASE 
			WHEN count(road_retro_reflectivity.id) > 0 AND MAX(road_retro_reflectivity."year") = date_part('year', now()) THEN '#F2A30B' 
			ELSE '#C8C8C8' 
		END AS retro_status_color,
		CASE 
			WHEN count(road_damage.id) > 0 AND MAX(road_damage."year") = date_part('year', now()) THEN '#F30B0D' 
			ELSE '#C8C8C8' 
		END AS damage_status_color`).
			Joins("LEFT JOIN road_info ON road_info.road_id = road.id AND road_info.status = 'A'").
			Joins(`LEFT JOIN road_condition ON road_condition.road_id = road.id 
			AND road_condition.status = 'A' 
			AND road_condition."year" <= date_part('year', now())`).
			Joins(`LEFT JOIN road_retro_reflectivity ON road_retro_reflectivity.road_id  = road.id 
			AND road_retro_reflectivity.status = 'A' 
			AND road_retro_reflectivity."year" <= date_part('year', now())`).
			Joins(`LEFT JOIN road_damage ON road_damage.road_id  = road.id 
			AND road_damage.status = 'A' 
			AND road_damage."year" <= date_part('year', now())`)

			// query = query.Preload("Sections.Roads", func(db *gorm.DB) *gorm.DB {
			// 	db = db.Select(`road.*, CASE
			// 						WHEN count(road_condition.id)  > 0 THEN true
			// 							ELSE false
			// 						END AS condition_status,
			// 						CASE
			// 							WHEN count(road_retro_reflectivity.id)  > 0 THEN true
			// 							ELSE false
			// 						END AS retro_status,
			// 						CASE
			// 							WHEN count(road_damage.id)  > 0 THEN true
			// 							ELSE false
			// 						END AS damage_status`).
			// 		Joins("LEFT JOIN road_info ON road_info.road_id = road.id AND road_info.status = 'A'").
			// 		Joins(`LEFT JOIN road_condition ON road_condition.road_id = road.id AND road_condition.status = 'A' and road_condition."year" = date_part('year', now())`).
			// 		Joins(`LEFT JOIN road_retro_reflectivity ON road_retro_reflectivity.road_id  = road.id AND road_retro_reflectivity.status = 'A' and road_retro_reflectivity."year" = date_part('year', now())`).
			// 		Joins(`LEFT JOIN road_damage ON road_damage.road_id  = road.id AND road_damage.status = 'A' and road_damage."year" = date_part('year', now())`)

		if prams.KmStart != nil || prams.KmEnd != nil {
			if prams.KmStart != nil && prams.KmEnd != nil {
				db = db.Where("(road_info.ref_road_type_id IN (1, 3) AND ((? BETWEEN road_info.km_start AND road_info.km_end) OR (? BETWEEN road_info.km_start AND road_info.km_end))) OR (road_info.ref_road_type_id IN (2, 4) AND ((? BETWEEN road_info.km_end AND road_info.km_start) OR (? BETWEEN road_info.km_end AND road_info.km_start)))", prams.KmStart, prams.KmEnd, prams.KmStart, prams.KmEnd)
			} else if prams.KmStart != nil {
				db = db.Where("((road_info.ref_road_type_id = 1 OR road_info.ref_road_type_id = 3) AND (? >= road_info.km_start and ? < road_info.km_end) ) OR ((road_info.ref_road_type_id = 2 OR road_info.ref_road_type_id = 4) AND (? <= road_info.km_start and ? > road_info.km_end))", prams.KmStart, prams.KmStart, prams.KmStart, prams.KmStart)
			} else if prams.KmEnd != nil {
				db = db.Where("((road_info.ref_road_type_id = 1 OR road_info.ref_road_type_id = 3) AND (? <= road_info.km_end and ? > road_info.km_start) ) OR ((road_info.ref_road_type_id = 2 OR road_info.ref_road_type_id = 4) AND  (? >= road_info.km_end and ? < road_info.km_start))", prams.KmEnd, prams.KmEnd, prams.KmEnd, prams.KmEnd)
			}
		}
		db = db.Where("road.road_level = 1 and road.is_active = true").Group("road.id").Order("road.id ASC")
		if len(roadIDs) > 0 {
			db = db.Where("road.id in (?)", roadIDs)
		}
		return db
	})

	query = query.Preload("Sections.Roads.RoadCondition", func(db *gorm.DB) *gorm.DB {
		if len(roadIDs) > 0 {
			db = db.Where("road_id in (?)", roadIDs)
		}
		db = db.Select("road_id, lane_no, MAX(surveyed_date) as surveyed_date").Where("status = ?", "A").Group("road_id, lane_no")

		return db
	})

	query = query.Preload("Sections.Roads.MaintenanceRoad", func(db *gorm.DB) *gorm.DB {
		if len(roadIDs) > 0 {
			db = db.Where("road_id in (?)", roadIDs)
		}
		db = db.Select("road_id, lane_no, MAX(maintenance.project_end_date) as project_end_date").Where("maintenance_road.status = ?", "A").
			Joins("left join maintenance on maintenance_road.maintenance_id  = maintenance.id and maintenance.status = 'A'").Group("road_id,lane_no")

		return db
	})

	query = query.Preload("Sections.Roads.RoadInfo", func(db *gorm.DB) *gorm.DB {
		if len(roadIDs) > 0 {
			db = db.Where("road_id in (?)", roadIDs)
		}
		db = db.Select("id,road_id, year, name, km_start, km_end,   ref_direction_id, st_astext(the_geom) as the_geom, revision, status, ramp_id,road_color_code, created_by,created_at,updated_by,updated_at,ref_road_type_id").Where("status = ?", "A")

		return db
	})

	query = query.Preload("Sections.Roads.RefSurface", func(db *gorm.DB) *gorm.DB {
		if len(roadIDs) > 0 {
			db = db.Where("road_surface.road_id in (?)", roadIDs)
		}
		db = db.Select("road_surface.road_id, ARRAY_AGG(DISTINCT road_surface_lane.ref_surface_id) as ref_surface_id").
			Joins("LEFT JOIN road_surface_lane ON road_surface_lane.road_surface_id = road_surface.id").
			Where("road_surface.status = ? and road_surface_lane.ref_surface_id <> ?", "A", 0).
			Group("road_surface.road_id")
		return db
	})

	query = query.Preload("Sections.Roads.RoadSurfaceIcon", func(db *gorm.DB) *gorm.DB {
		if len(roadIDs) > 0 {
			db = db.Where("road_surface.road_id in (?)", roadIDs)
		}
		db = db.Select("DISTINCT ON (ref_surface.surface_group,road_surface.road_id) ref_surface.surface_group, ref_surface.id , road_surface.road_id, ref_surface.surface_group as name, CASE WHEN ref_surface.surface_group = 'Concrete' THEN '#398BF7' ELSE '#7460EE' END AS color_code").
			Joins("RIGHT JOIN road_surface_lane ON road_surface.id = road_surface_lane.road_surface_id").
			Joins("LEFT JOIN ref_surface ON ref_surface.id = road_surface_lane.ref_surface_id").
			Where("road_surface.status = ? and road_surface_lane.ref_surface_id <> ?", "A", 0)

		return db
	})

	query = query.Preload("Sections.Roads.RoadInfo.RefRoadType")
	query = query.Preload("Sections.Roads.RoadInfo.RefDirection")

	query = query.Preload("Sections.Roads.RoadGeom", func(db *gorm.DB) *gorm.DB {
		if len(roadIDs) > 0 {
			db = db.Where("road_id in (?)", roadIDs)
		}
		db = db.Select("id,road_id, lane_no, km_start,km_end,ST_ASTEXT(ST_FORCE2D(the_geom)) as the_geom,revision,status,remark,created_by,created_at,updated_by,updated_at").Where("status = ?", "A").Order("revision DESC").Order("revision")

		return db
	})

	query = query.Preload("Sections.Roads.ChildRoads", func(db *gorm.DB) *gorm.DB {
		db = db.Select(`road.*, 
		CASE 
			WHEN count(road_condition.id) > 0 THEN true
			ELSE false 
		END AS condition_status,
		CASE 
			WHEN count(road_retro_reflectivity.id) > 0 THEN true 
			ELSE false 
		END AS retro_status,
		CASE 
			WHEN count(road_damage.id) > 0 THEN true 
			ELSE false 
		END AS damage_status,
		-- Color Logic based on current year or fallback
		CASE 
			WHEN count(road_condition.id) > 0 AND MAX(road_condition."year") = date_part('year', now()) THEN '#F27A0D' 
			ELSE '#C8C8C8' 
		END AS condition_status_color,
		CASE 
			WHEN count(road_retro_reflectivity.id) > 0 AND MAX(road_retro_reflectivity."year") = date_part('year', now()) THEN '#F2A30B' 
			ELSE '#C8C8C8' 
		END AS retro_status_color,
		CASE 
			WHEN count(road_damage.id) > 0 AND MAX(road_damage."year") = date_part('year', now()) THEN '#F30B0D' 
			ELSE '#C8C8C8' 
		END AS damage_status_color`).
			Joins("LEFT JOIN road_info ON road_info.road_id = road.id AND road_info.status = 'A'").
			Joins(`LEFT JOIN road_condition ON road_condition.road_id = road.id 
			AND road_condition.status = 'A' 
			AND road_condition."year" <= date_part('year', now())`).
			Joins(`LEFT JOIN road_retro_reflectivity ON road_retro_reflectivity.road_id  = road.id 
			AND road_retro_reflectivity.status = 'A' 
			AND road_retro_reflectivity."year" <= date_part('year', now())`).
			Joins(`LEFT JOIN road_damage ON road_damage.road_id  = road.id 
			AND road_damage.status = 'A' 
			AND road_damage."year" <= date_part('year', now())`)

		// db = db.Select(`road.*, CASE
		// 				WHEN count(road_condition.id)  > 0 THEN true
		// 					ELSE false
		// 				END AS condition_status,
		// 				CASE
		// 					WHEN count(road_retro_reflectivity.id)  > 0 THEN true
		// 					ELSE false
		// 				END AS retro_status,
		// 				CASE
		// 					WHEN count(road_damage.id)  > 0 THEN true
		// 					ELSE false
		// 				END AS damage_status`).
		// 	Joins(`LEFT JOIN road_condition ON road_condition.road_id = road.id AND road_condition.status = 'A' and road_condition."year" = date_part('year', now())`).
		// 	Joins(`LEFT JOIN road_retro_reflectivity ON road_retro_reflectivity.road_id  = road.id AND road_retro_reflectivity.status = 'A' and road_retro_reflectivity."year" = date_part('year', now())`).
		// 	Joins(`LEFT JOIN road_damage ON road_damage.road_id  = road.id AND road_damage.status = 'A' and road_damage."year" = date_part('year', now())`)
		if len(roadIDs) > 0 {
			db = db.Where("road.id in (?)", roadIDs)
		}
		db = db.Where("road.is_active = true").Group("road.id").Order("road.id")
		return db
	})

	query = query.Preload("Sections.Roads.ChildRoads.RoadCondition", func(db *gorm.DB) *gorm.DB {
		if len(roadIDs) > 0 {
			db = db.Where("road_id in (?)", roadIDs)
		}
		db = db.Select("road_id, lane_no, MAX(surveyed_date) as surveyed_date").Where("status = ?", "A").Group("road_id, lane_no")

		return db
	})

	query = query.Preload("Sections.Roads.ChildRoads.MaintenanceRoad", func(db *gorm.DB) *gorm.DB {
		if len(roadIDs) > 0 {
			db = db.Where("road_id in (?)", roadIDs)
		}
		db = db.Select("road_id, lane_no, MAX(maintenance.project_end_date) as project_end_date").Where("maintenance_road.status = ?", "A").
			Joins("left join maintenance on maintenance_road.maintenance_id  = maintenance.id and maintenance.status = 'A'").Group("road_id,lane_no")

		return db
	})

	query = query.Preload("Sections.Roads.ChildRoads.RefSurface", func(db *gorm.DB) *gorm.DB {
		if len(roadIDs) > 0 {
			db = db.Where("road_surface.road_id in (?)", roadIDs)
		}
		db = db.Select("road_surface.road_id, ARRAY_AGG(DISTINCT road_surface_lane.ref_surface_id) as ref_surface_id").
			Joins("LEFT JOIN road_surface_lane ON road_surface_lane.road_surface_id = road_surface.id").
			Where("road_surface.status = ? and road_surface_lane.ref_surface_id <> ?", "A", 0).
			Group("road_surface.road_id")
		return db
	})

	query = query.Preload("Sections.Roads.ChildRoads.RoadSurfaceIcon", func(db *gorm.DB) *gorm.DB {
		if len(roadIDs) > 0 {
			db = db.Where("road_surface.road_id in (?)", roadIDs)
		}
		db = db.Select("DISTINCT ON (ref_surface.surface_group,road_surface.road_id) ref_surface.surface_group, ref_surface.id , road_surface.road_id, ref_surface.surface_group as name, CASE WHEN ref_surface.surface_group = 'Concrete' THEN '#398BF7' ELSE '#7460EE' END AS color_code").
			Joins("RIGHT JOIN road_surface_lane ON road_surface.id = road_surface_lane.road_surface_id").
			Joins("LEFT JOIN ref_surface ON ref_surface.id = road_surface_lane.ref_surface_id").
			Where("road_surface.status = ? and road_surface_lane.ref_surface_id <> ?", "A", 0)
		return db
	})

	query = query.Preload("Sections.Roads.ChildRoads.RoadGeom", func(db *gorm.DB) *gorm.DB {
		if len(roadIDs) > 0 {
			db = db.Where("road_id in (?)", roadIDs)
		}
		db = db.Select("id,road_id, lane_no, km_start,km_end,ST_ASTEXT(ST_FORCE2D(the_geom)) as the_geom,revision,status,remark,created_by,created_at,updated_by,updated_at").Where("status = ?", "A").Order("revision DESC").Order("revision")
		return db
	})

	query = query.Preload("Sections.Roads.ChildRoads.RoadInfo", func(db *gorm.DB) *gorm.DB {
		if len(roadIDs) > 0 {
			db = db.Where("road_id in (?)", roadIDs)
		}
		db = db.Select("id,road_id, year, name, km_start, km_end,   ref_direction_id, st_astext(the_geom) as the_geom, revision, status, ramp_id,road_color_code, created_by,created_at,updated_by,updated_at,ref_road_type_id").Where("status = ?", "A")
		return db
	})

	query = query.Preload("Sections.Roads.ChildRoads.RoadInfo.RefRoadType")

	query = query.Preload("Sections.Roads.ChildRoads.RoadInfo.RefDirection")

	query = query.Preload("Sections.Roads.ChildRoads.ChildRoads")

	if err := query.Find(&roadList).Error; err != nil {
		return roadList, err
	}

	return roadList, nil
}

func (t *roadRepository) GetRoadByID(roadID int) (*models.RoadById, error) {
	var road models.RoadById

	query := t.conn

	query = query.Where("id = ?", roadID)

	query = query.Preload("RoadInfo", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,road_id, year, name, km_start, km_end,   ref_direction_id, st_astext(the_geom) as the_geom, revision, status, ramp_id,road_color_code, created_by,created_at,updated_by,updated_at,ref_road_type_id", "center_lane_shape_file_path", "center_line_shape_file_path", "remark", "year_construction_completed").Where("status = ?", "A")
	})

	query = query.Preload("RoadSurfaceIcon", func(db *gorm.DB) *gorm.DB {
		return db.Select("DISTINCT ON (ref_surface.surface_group) ref_surface.surface_group, ref_surface.id , road_surface.road_id, ref_surface.surface_group as name, CASE WHEN ref_surface.surface_group = 'Concrete' THEN '#398BF7' ELSE '#7460EE' END AS color_code").
			Joins("RIGHT JOIN road_surface_lane ON road_surface.id = road_surface_lane.road_surface_id").
			Joins("LEFT JOIN ref_surface ON ref_surface.id = road_surface_lane.ref_surface_id").
			Where("road_surface.status = ? and road_surface_lane.ref_surface_id <> ?", "A", 0)
	})

	query = query.Preload("RoadInfo.User")

	query = query.Preload("RoadGeom", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,road_id, lane_no, km_start,km_end,ST_ASTEXT(ST_FORCE2D(the_geom)) as the_geom,revision,status,remark,created_by,created_at,updated_by,updated_at").Where("status = ?", "A").Order("revision DESC").Order("revision")
	})

	query = query.Preload("RoadInfo.RefRoadType")

	query = query.Preload("RoadSection", func(db *gorm.DB) *gorm.DB {
		db = db.Select("road_section.id,road_section.road_group_id,road_section.number,road_section.name_origin_th,road_section.name_destination_th,road_section.name_origin_en,road_section.name_destination_en,road_section.km_start, road_section.km_end, road_section.distance,road_section.ref_division_code,road_section.ref_district_code,road_section.ref_depot_code, ARRAY_AGG(ref_province.name) AS province").
			Joins("JOIN ref_province ON ref_province.province_code = ANY(road_section.province_code)").Group("road_section.id,road_section.road_group_id,road_section.number,road_section.name_origin_th,road_section.name_destination_th,road_section.name_origin_en,road_section.name_destination_en,road_section.km_start, road_section.km_end, road_section.distance,road_section.ref_division_code,road_section.ref_district_code,road_section.ref_depot_code")

		return db.Order("road_section.id ASC")
	})

	query = query.Preload("RoadSection.RefDivision", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, division_code, name, name_en, st_astext(the_geom) as the_geom")
	})

	query = query.Preload("RoadSection.RefDistrict", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, district_code, name, name_en, st_astext(the_geom) as the_geom")
	})

	query = query.Preload("RoadSection.RefDepot", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, depot_code, name, st_astext(the_geom) as the_geom")
	})

	query = query.Preload("RoadSection.RoadGroup")

	if err := query.Find(&road).Error; err != nil {
		return nil, err
	}

	return &road, nil
}

func (t *roadRepository) GetRoadGroup() ([]models.RoadGroup, error) {
	var roadGroup []models.RoadGroup
	if err := t.conn.Order("id").Find(&roadGroup).Error; err != nil {
		return roadGroup, err
	}
	return roadGroup, nil
}

func (t *roadRepository) GetRoadStatusSurface(roadID int) (responses.StatusCount, error) {

	var statusCount responses.StatusCount
	statusCount.Category = "SURFACE"
	query := t.conn
	query = query.Select("SUM(CASE WHEN status = 'T' THEN 1 ELSE 0 END) as count_temp, SUM(CASE WHEN status = 'W' THEN 1 ELSE 0 END) as count_waiting, SUM(CASE WHEN status = 'R' THEN 1 ELSE 0 END) as count_rejected")
	query = query.Table("road_surface")
	if err := query.Where("road_id = ?", roadID).Find(&statusCount).Error; err != nil {
		return statusCount, err
	}
	return statusCount, nil
}

func (t *roadRepository) GetRoadTypeIcon() ([]models.RefRoadTypeIcon, error) {
	var roadTypeIcon []models.RefRoadTypeIcon

	if err := t.conn.Find(&roadTypeIcon).Error; err != nil {
		return roadTypeIcon, err
	}
	return roadTypeIcon, nil
}

func (t *roadRepository) GetRoadDirectionLaneList(roadID int) (models.RoadInfoGeomDirection, error) {
	var roadInfoGeomDirection models.RoadInfoGeomDirection
	query := t.conn

	query = query.Preload("RoadInfo")

	query = query.Preload("RoadInfo.Direction", func(db *gorm.DB) *gorm.DB {
		return db.Order("id")
	})
	query = query.Preload("RoadGeom", func(db *gorm.DB) *gorm.DB {
		return db.Order("lane_no")
	})
	if err := query.Where("id = ?", roadID).Where("is_active = ?", true).First(&roadInfoGeomDirection).Error; err != nil {
		return roadInfoGeomDirection, err
	}
	return roadInfoGeomDirection, nil
}

func (t *roadRepository) GetRoadByRoadGrpID(roadGrpID int) ([]models.RoadInfo, error) {
	query := t.conn
	var roadInfo []models.RoadInfo
	if err := query.Table("road").Select("road_info.road_id as road_id, road_info.name as name").Joins("join road_info on road.id = road_info.road_id").Where("status = ?", "A").Where("road_group_id = ?", roadGrpID).Order("road.road_level asc").Order("road.ref_direction_id asc").Order("road_id asc").Scan(&roadInfo).Error; err != nil {
		return roadInfo, err
	}

	return roadInfo, nil
}

func (t *roadRepository) GetVolumeAadtApproved(roadGrpID int) (responses.StatusCount, error) {
	query := t.conn
	// var aadt []models.VolumeAadt
	var statusCount responses.StatusCount
	if err := query.Table("volume_aadt").Select("SUM ( CASE WHEN status = 'T' THEN 1 ELSE 0 END ) AS count_temp, SUM ( CASE WHEN status = 'W' THEN 1 ELSE 0 END ) AS count_waiting, SUM ( CASE WHEN status = 'R' THEN 1 ELSE 0 END ) AS count_rejected").Where("road_group_id = ?", roadGrpID).Find(&statusCount).Error; err != nil {
		return statusCount, err
	}
	return statusCount, nil
}

func (t *roadRepository) GetVolumeAccidentApproved(roadGrpID int) (responses.StatusCount, error) {
	query := t.conn
	var statusCount responses.StatusCount
	if err := query.Table("volume_accident").Select("SUM ( CASE WHEN status = 'T' THEN 1 ELSE 0 END ) AS count_temp, SUM ( CASE WHEN status = 'W' THEN 1 ELSE 0 END ) AS count_waiting, SUM ( CASE WHEN status = 'R' THEN 1 ELSE 0 END ) AS count_rejected").Where("road_group_id = ?", roadGrpID).Find(&statusCount).Error; err != nil {
		return statusCount, err
	}
	return statusCount, nil
}

func (t *roadRepository) GetRoadInfoByRoadID(roadID int) (models.RoadInfo, error) {
	var roadInfo models.RoadInfo
	if err := t.conn.Where("road_id = ? AND status = 'A'", roadID).First(&roadInfo).Error; err != nil {
		return roadInfo, err
	}

	return roadInfo, nil
}

func (t *roadRepository) GetRoadSectionByID(RoadsectionId int) (*models.RoadSectionById, error) {
	var roadSection models.RoadSectionById
	query := t.conn
	query = query.Select("road_section.id,road_section.road_group_id,road_section.number,road_section.name_origin_th,road_section.name_destination_th,road_section.name_origin_en,road_section.name_destination_en,road_section.km_start, road_section.km_end, road_section.distance,road_section.ref_division_code,road_section.ref_district_code,road_section.ref_depot_code, ARRAY_AGG(ref_province.name) AS province").
		Joins("JOIN ref_province ON ref_province.province_code = ANY(road_section.province_code)").Group("road_section.id,road_section.road_group_id,road_section.number,road_section.name_origin_th,road_section.name_destination_th,road_section.name_origin_en,road_section.name_destination_en,road_section.km_start, road_section.km_end, road_section.distance,road_section.ref_division_code,road_section.ref_district_code,road_section.ref_depot_code").Where("road_section.id = ?", RoadsectionId)

	query = query.Preload("RefDivision", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, division_code, name, name_en, st_astext(the_geom) as the_geom")
	})

	query = query.Preload("RefDistrict", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, district_code, name, name_en, st_astext(the_geom) as the_geom")
	})

	query = query.Preload("RefDepot", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, depot_code, name, st_astext(the_geom) as the_geom")
	})

	query = query.Preload("RoadGroup")

	if err := query.Find(&roadSection).Error; err != nil {
		return nil, err
	}

	return &roadSection, nil
}

func (t *roadRepository) GetDataById(model interface{}, id int) error {
	query := t.conn
	query = query.Where("id = ?", id)
	err := query.Find(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *roadRepository) GetRoadMaxSeq() (int, error) {
	var road models.Road
	if err := t.conn.Select("max(seq) as seq").Find(&road).Error; err != nil {
		return 0, err
	}
	return road.Seq, nil
}

func (r *roadRepository) CreateData(tx *gorm.DB, model interface{}) error {
	err := tx.Create(model).Error
	if err != nil {

		return err
	}

	return nil
}

func (r *roadRepository) CreateRoadInfo(tx *gorm.DB, roadInfo *models.RoadInfo) error {
	sqlStatement := `
        INSERT INTO road_info (
            road_id, year, ref_direction_id, name, km_start, km_end, the_geom, revision, status,
            ramp_id, road_color_code, created_by, created_at, updated_by, updated_at, remark,
            ref_road_type_id, center_lane_shape_file_path, center_line_shape_file_path, year_construction_completed
        ) VALUES (?, ?, ?, ?, ?, ?, ST_GeomFromText(?, 4326), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Convert the WKT LINESTRING into a string and pass it twice to the query
	wktLineString := roadInfo.TheGeom

	result := tx.Exec(sqlStatement,
		roadInfo.RoadId, roadInfo.Year, roadInfo.RefDirectionId, roadInfo.Name,
		roadInfo.KmStart, roadInfo.KmEnd, wktLineString, roadInfo.Revision, roadInfo.Status,
		roadInfo.RampId, roadInfo.RoadColorCode, roadInfo.CreatedBy, roadInfo.CreatedAt,
		roadInfo.UpdatedBy, roadInfo.UpdatedAt, roadInfo.Remark, roadInfo.RefRoadTypeID,
		roadInfo.CenterLaneShapeFilePath, roadInfo.CenterLineShapeFilePath, roadInfo.YearConstructionCompleted)

	// Handle errors
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (t *roadRepository) CreateRoadGeom(tx *gorm.DB, roadGeoms []models.RoadGeom) error {

	sqlStatement := `
	INSERT INTO road_geom (
		road_id, lane_no, km_start, km_end, the_geom, revision, status,
		remark, created_by, created_at, updated_by, updated_at
	) VALUES `

	var valueStrings []string
	var valueArgs []interface{}

	for _, roadGeom := range roadGeoms {
		// Convert the WKT LINESTRING into a string and pass it twice to the query
		wktLineString := roadGeom.TheGeom

		valueStrings = append(valueStrings, "(?, ?, ?, ?, ST_GeomFromText(?, 4326), ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs,
			roadGeom.RoadId, roadGeom.LaneNo, roadGeom.KmStart, roadGeom.KmEnd,
			wktLineString, roadGeom.Revision, roadGeom.Status,
			roadGeom.Remark, roadGeom.CreatedBy, roadGeom.CreatedAt,
			roadGeom.UpdatedBy, roadGeom.UpdatedAt)
	}

	sqlStatement += strings.Join(valueStrings, ", ")
	result := tx.Exec(sqlStatement, valueArgs...)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *roadRepository) UpdateRoadInfo(tx *gorm.DB, roadInfo *models.RoadInfo) error {
	sqlStatement := `
        INSERT INTO road_info (
            road_id, year, ref_direction_id, name, km_start, km_end, the_geom, revision, status,
            ramp_id, road_color_code, created_by, created_at, updated_by, updated_at, remark,
            ref_road_type_id, center_lane_shape_file_path, center_line_shape_file_path, year_construction_completed
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Convert the WKT LINESTRING into a string and pass it twice to the query
	wktLineString := roadInfo.TheGeom

	result := tx.Exec(sqlStatement,
		roadInfo.RoadId, roadInfo.Year, roadInfo.RefDirectionId, roadInfo.Name,
		roadInfo.KmStart, roadInfo.KmEnd, wktLineString, roadInfo.Revision, roadInfo.Status,
		roadInfo.RampId, roadInfo.RoadColorCode, roadInfo.CreatedBy, roadInfo.CreatedAt,
		roadInfo.UpdatedBy, roadInfo.UpdatedAt, roadInfo.Remark, roadInfo.RefRoadTypeID,
		roadInfo.CenterLaneShapeFilePath, roadInfo.CenterLineShapeFilePath, roadInfo.YearConstructionCompleted)

	// Handle errors
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (t *roadRepository) UpdateRoadGeom(tx *gorm.DB, roadGeoms []models.RoadGeom) error {
	sqlStatement := `
	INSERT INTO road_geom (
		road_id, lane_no, km_start, km_end, the_geom, revision, status,
		remark, created_by, created_at, updated_by, updated_at
	) VALUES `

	var valueStrings []string
	var valueArgs []interface{}

	for _, roadGeom := range roadGeoms {
		// Convert the WKT LINESTRING into a string and pass it twice to the query
		wktLineString := roadGeom.TheGeom

		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs,
			roadGeom.RoadId, roadGeom.LaneNo, roadGeom.KmStart, roadGeom.KmEnd,
			wktLineString, roadGeom.Revision, roadGeom.Status,
			roadGeom.Remark, roadGeom.CreatedBy, roadGeom.CreatedAt,
			roadGeom.UpdatedBy, roadGeom.UpdatedAt)
	}

	sqlStatement += strings.Join(valueStrings, ", ")
	result := tx.Exec(sqlStatement, valueArgs...)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (t *roadRepository) DeleteRoad(tx *gorm.DB, roadID int, userID int) error {
	if err := tx.Model(&models.Road{}).Where("id = ?", roadID).Update("is_active", "false").Error; err != nil {
		return err
	}
	return nil
}

func (t *roadRepository) DeleteRoadInfo(tx *gorm.DB, roadID int, userID int) error {
	roadInfoUpdate := models.RoadInfo{
		Status:    "D",
		UpdatedAt: time.Now(),
		UpdatedBy: userID,
	}
	if err := tx.Model(&models.RoadInfo{}).Where("road_id = ?", roadID).Updates(roadInfoUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadRepository) GetLastRoadInfoByID(roadID int) (models.RoadInfoGeomData, error) {
	var roadInfo models.RoadInfoGeomData
	query := t.conn
	query = query.Select("* , ST_ASTEXT(the_geom) as line_string")
	query = query.Table("road_info")
	query = query.Preload("RoadGeom", func(db *gorm.DB) *gorm.DB {
		return db.Select("*, ST_ASTEXT(the_geom) as line_string").Where("status = ?", "A")
	})
	if err := query.Where("road_id = ?", roadID).Where("status = ?", "A").Order("road_id asc").Find(&roadInfo).Error; err != nil {
		return roadInfo, err
	}

	return roadInfo, nil
}

func (t *roadRepository) UpdateDirectionRoad(tx *gorm.DB, roadID, refDirectionId int) error {
	roadUpdate := models.RoadInfo{
		RefDirectionId: refDirectionId,
	}
	if err := tx.Model(&models.RoadInfo{}).Where("road_id = ? AND status = ?", roadID, "A").Updates(roadUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadRepository) CountParent(tx *gorm.DB, roadId int) (*int64, error) {
	var count int64
	if err := tx.Model(&models.Road{}).Where("parent_road_id = ? AND is_active = ?", roadId, true).Count(&count).Error; err != nil {
		return nil, err
	}

	return &count, nil
}

func (t *roadRepository) DeleteRoadGeom(tx *gorm.DB, roadID int, userID int) error {
	roadGeomUpdate := models.RoadGeom{
		Status:    "D",
		UpdatedAt: time.Now(),
		UpdatedBy: userID,
	}
	if err := tx.Model(&models.RoadGeom{}).Where("road_id = ?", roadID).Updates(roadGeomUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadRepository) GetRoadLanes(roadID int) ([]models.RoadLanes, error) {
	var result []models.RoadLanes
	query := t.conn.Select(`road_geom.road_id, road_geom.lane_no, road_info.ref_direction_id as ref_direction_id ,  
							ref_direction.name as ref_direction_name`).
		Joins("LEFT JOIN road_info on road_info.road_id = road_geom.road_id and road_info.status = 'A'").
		Joins("LEFT JOIN ref_direction on ref_direction.id = road_info.ref_direction_id")

	err := query.Where("road_geom.road_id = ? and road_geom.status = ?", roadID, "A").
		Group("road_geom.road_id, road_geom.lane_no, road_info.ref_direction_id, ref_direction.name").
		Order(`CASE
        WHEN road_info.ref_direction_id = 1 THEN road_geom.lane_no
        WHEN road_info.ref_direction_id = 2 THEN road_geom.lane_no * -1
        ELSE road_geom.lane_no
    	END ASC`).
		Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *roadRepository) StartTransSection() *gorm.DB {
	tx := t.conn.Begin()
	return tx
}

func (t *roadRepository) RollBack(tx *gorm.DB) error {
	tx.Rollback()
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func (t *roadRepository) Commit(tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *roadRepository) GetParamsConditionByID(ID int, conditionType string) (models.ParamsCondition, error) {
	var paramsCondition models.ParamsCondition
	query := r.conn.Debug()
	query = query.Where("ref_owner_id = ?", ID)
	if conditionType != "" {
		query = query.Where("condition_type = ?", conditionType)
	}
	err := query.First(&paramsCondition).Error
	if err != nil {
		return paramsCondition, err
	}
	return paramsCondition, nil
}

func (r *roadRepository) GetParamsRoadLine(ID int) (models.ParamsRoadLine, error) {
	var paramsRoadLine models.ParamsRoadLine
	query := r.conn
	query = query.Where("ref_owner_road_line_id = ?", ID)
	err := query.First(&paramsRoadLine).Error
	if err != nil {
		return paramsRoadLine, err
	}
	return paramsRoadLine, nil
}

func (r *roadRepository) GetroadConditionSurvey(roadID, laneNo int) ([]models.RoadConditionSurvey, error) {
	var roadConditionSurvey []models.RoadConditionSurvey
	query := r.conn
	query.Table("road_condition_survey")
	query = query.Joins("INNER JOIN road_condition ON road_condition_survey.road_condition_id = road_condition.id")
	query = query.Where("road_condition.road_id = ?", roadID)
	query = query.Where("road_condition.lane_no = ?", laneNo)
	query = query.Order("road_condition.Surveyed_date DESC")
	if err := query.Find(&roadConditionSurvey).Error; err != nil {
		return roadConditionSurvey, err
	}
	return roadConditionSurvey, nil
}

func (r *roadRepository) GetroadConditionSurvey100M(roadID, laneNo int) ([]models.RoadConditionSurvey100M, error) {
	var roadConditionSurvey100M []models.RoadConditionSurvey100M
	query := r.conn
	query = query.Table("road_condition_survey")
	query = query.Select("road_condition_survey.id as id ,road_condition.id as road_condition_id ,road_condition.road_id as road_id ,road_condition.lane_no as lane_no ,road_condition_survey_100m.km_start as km_start,road_condition_survey_100m.km_end as km_end ,road_condition_survey_100m.iri as iri ,road_condition_survey_100m.rut as rut,road_condition_survey_100m.ifi as ifi ,road_condition_survey_100m.survey_type as survey_type, road_condition.surveyed_date as survey_date ,road_condition.year as year")
	query = query.Joins("INNER JOIN road_condition ON road_condition_survey.road_condition_id = road_condition.id")
	query = query.Joins("INNER JOIN road_condition_survey_100m ON road_condition_survey.id = road_condition_survey_100m.road_condition_survey_id")
	query = query.Where("road_condition.road_id = ?", roadID)
	query = query.Where("road_condition.lane_no = ?", laneNo)
	query = query.Order("road_condition.Surveyed_date DESC")
	if err := query.Find(&roadConditionSurvey100M).Error; err != nil {
		return roadConditionSurvey100M, err
	}
	return roadConditionSurvey100M, nil
}

func (r *roadRepository) GetRoadSurfaceLaneAll(roadIds []int) (map[int][]models.RoadSurfaceLane, error) {
	var roadSurfaceLanes []models.RoadSurfaceLane
	query := r.conn
	query.Table("road_surface_lane")
	query = query.Joins("INNER JOIN road_surface ON road_surface_lane.road_surface_id = road_surface.id")
	query = query.Where("status = ?", "A")
	query = query.Where("status = ?", "A")
	if err := query.Find(&roadSurfaceLanes).Error; err != nil {
		return map[int][]models.RoadSurfaceLane{}, err
	}
	roadSurfaceData := make(map[int][]models.RoadSurfaceLane)
	for _, item := range roadSurfaceLanes {
		roadSurfaceData[item.RoadId] = append(roadSurfaceData[item.RoadId], item)
	}
	return roadSurfaceData, nil
}

func (r *roadRepository) GetRoadRetroReflectivity100M(roadID, laneNo int) ([]models.RoadRetroReflectivityRange, error) {
	var roadRetroReflectivityRange []models.RoadRetroReflectivityRange
	query := r.conn
	query = query.Table("road_retro_reflectivity")
	query = query.Joins("INNER JOIN road_retro_reflectivity_range ON road_retro_reflectivity.id = road_retro_reflectivity_range.road_retro_reflectivity_id")
	query = query.Where("road_retro_reflectivity.road_id = ?", roadID)
	query = query.Where("road_retro_reflectivity.lane_no = ?", laneNo)
	if err := query.Find(&roadRetroReflectivityRange).Error; err != nil {
		return roadRetroReflectivityRange, err
	}
	return roadRetroReflectivityRange, nil
}

type D struct {
	RoadID int `json:"road_id"`
}

func (r *roadRepository) GetRoadID(params requests.RoadPrams) ([]int, error) {
	var d []D
	query := r.conn
	query = query.Table("road_info ri")
	query = query.Select("ri.road_id")
	query = query.Joins("INNER JOIN road r on ri.road_id = r.id")
	query = query.Joins("INNER JOIN road_group rp on r.road_group_id = rp.id")
	query = query.Joins("INNER JOIN road_section rs on r.road_section_id = rs.id")
	if len(params.RoadGroupId) != 0 {
		query = query.Where("rp.id IN ?", params.RoadGroupId)
	}
	if len(params.RoadSectionId) != 0 {
		query = query.Where("rs.id IN ?", params.RoadSectionId)
	}
	if len(params.DepotCode) != 0 {
		query = query.Where("rs.ref_depot_code IN ?", params.DepotCode)
	}
	query = query.Where("ri.status = 'A'")
	query = query.Where("r.is_active = true")
	if err := query.Find(&d).Error; err != nil {

	}
	roadIDs := []int{}
	for _, item := range d {
		roadIDs = append(roadIDs, item.RoadID)
	}

	return roadIDs, nil
}
