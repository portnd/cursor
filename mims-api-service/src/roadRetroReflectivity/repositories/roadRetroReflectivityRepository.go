package repositories

import (
	"fmt"
	"time"

	helpers "gitlab.com/mims-api-service/helpers"
	models "gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/roadRetroReflectivity/handlers"
	"gitlab.com/mims-api-service/src/roadRetroReflectivity/usecases"

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

// /////////////////////////
func (t *Repository) GetRole(userId uint) ([]models.UserRole, error) {
	var userRole []models.UserRole
	if err := t.conn.Where("user_id = ?", userId).Find(&userRole).Error; err != nil {
		return userRole, err
	}
	return userRole, nil
}

func (t *Repository) GetAccessControl(roles []int) ([]models.AccessControl, error) {
	var accessControl []models.AccessControl
	query := t.conn
	query = query.Joins("JOIN role_access_control on access_control.id = role_access_control.access_control_id")
	err := query.Find(&accessControl).Error
	if err != nil {
		fmt.Println(err)
	}
	return accessControl, nil
}

func (t *Repository) GetUserById(userId uint) (models.Users, error) {
	var user models.Users
	if err := t.conn.Where("id = ?", userId).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (t *Repository) GetUserDepartmentById(userId int) (models.UserDepartment, error) {
	var user models.UserDepartment
	if err := t.conn.Where("id = ?", userId).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (t *Repository) GetRoadRetroReflectivityList(roadID int) ([]models.RetroReflectivityList, error) {
	var rcList []models.RetroReflectivityList
	err := t.conn.Table("road_info AS ri").
		Select("DISTINCT ri.road_id AS road_id, rrt.year AS year, rrt.id AS id, rrt.id_parent AS id_parent, rrt.revision AS revision, ref_d.id AS direction_id, ref_d.name AS direction_name, rrt.line_no AS line_no, rrt.surveyed_date AS surveyed_date,rrt.updated_date AS updated_date").
		Joins("LEFT JOIN ref_direction ref_d ON ri.ref_direction_id = ref_d.id").
		Joins("LEFT JOIN road_retro_reflectivity rrt ON ri.road_id= rrt.road_id  ").
		Where("rrt.status = 'A'").
		Where("ri.road_id = ?", roadID).
		Order("direction_id, line_no ASC, surveyed_date DESC, updated_date DESC").
		Find(&rcList).Error
	if err != nil {
		return nil, err
	}
	return rcList, nil

}

// func (t *Repository) GetRoadConditionDetails(idParent int, rcStatus models.SqlCondition, sqlCondition models.SqlCondition) ([]models.RoadConditionDetails, error) {

// 	var roadConditionDetails []models.RoadConditionDetails
// 	var roads struct {
// 		Revision    int `gorm:"column:revision"`
// 		DirectionID int `gorm:"column:direction_id"`
// 	}
// 	err := t.conn.
// 		Table("road r").
// 		Select("CASE WHEN r.ref_direction_id = 2 AND r.ref_road_type_id = 1 THEN 2 ELSE 1 END as direction_id, rc.revision").
// 		Joins("LEFT JOIN road_condition rc ON r.id = rc.road_id").
// 		Where(rcStatus.Where).
// 		Where("rc.id_parent = ? AND r.is_active = ? ", idParent, true).
// 		Order("rc.revision DESC").
// 		Find(&roads).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	OrderBy := ""
// 	switch roads.DirectionID {
// 	case 1:
// 		OrderBy = "revision DESC, road_id ASC , updated_date DESC , km_start_m ASC"
// 	case 2:
// 		OrderBy = "revision DESC, road_id ASC , updated_date DESC , km_start_m DESC"
// 	}

// 	query := t.conn.Table("road r")
// 	query = query.Select(`r.id AS road_id,
// 	rc.id AS id,
// 	rc.id_parent AS id_parent,
// 	rc.updated_date AS updated_date,
// 	rc.updated_by AS updated_by,
// 	rc.status AS status,
// 	ref_ds.name AS status_text,
// 	rc.revision AS revision,
// 	ref_d.id AS direction_id,
// 	ref_d.name AS direction_name,
// 	r.ref_road_type_id AS road_type_id,
// 	rc_m.img_filepath AS img_filepath,
// 	rc_km.km_start AS km_start_km,
// 	rc_km.km_end AS km_end_km,
// 	rc_m.km_start AS km_start_m,
// 	rc_m.km_end AS km_end_m,
// 	` + sqlCondition.Select + `
// 	ST_ASTEXT(ST_FORCE2D(rc_m.the_geom)) AS geom_cl,reject_reason`)
// 	query = query.Joins("LEFT JOIN ref_direction ref_d ON r.ref_direction_id = ref_d.id")
// 	query = query.Joins("LEFT JOIN road_condition  rc ON r.id = rc.road_id")
// 	query = query.Joins("LEFT  JOIN road_info  ri ON r.id = ri.road_id")
// 	query = query.Joins("LEFT JOIN ref_data_status ref_ds ON ref_ds.status_code  = rc.status")
// 	query = query.Joins("LEFT JOIN road_condition_survey rc_km ON rc_km.road_condition_id = rc.id")
// 	query = query.Joins("LEFT JOIN road_condition_survey_m rc_m ON rc_m.road_condition_survey_id = rc_km.id")
// 	query = query.Where(rcStatus.Where)
// 	query = query.Where("rc.id_parent = ? AND r.is_active = ?", idParent, true)
// 	query = query.Order(OrderBy)

// 	query = query.Find(&roadConditionDetails)
// 	err = query.Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return roadConditionDetails, nil

// }

// func (t *Repository) GetRoadConditionGrades(roadId int) ([]models.RoadConditionGrade, error) {
// 	var grades []models.RoadConditionGrade

// 	err := t.conn.Table("road r").
// 		Select("pc.left_value, pc.left_condition, pc.right_value, pc.right_condition,pc.condition_type AS condition_type, pc.ref_grade_id AS grade_id").
// 		Joins("LEFT JOIN road_owner ro ON r.id = ro.road_id").
// 		Joins("LEFT JOIN params_condition pc ON pc.ref_owner_id = ro.ref_owner_id").
// 		Where("r.id = ?", roadId).
// 		Order("condition_type").
// 		Order("grade_id").
// 		Find(&grades).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return grades, nil
// }

func (t *Repository) GetFullGeom(roadId int, laneNo int) (models.FullGeom, error) {
	var geom models.FullGeom

	err := t.conn.Where("road_geom.road_id = ? AND road_geom.lane_no = ? AND road_geom.status = ?", roadId, laneNo, "A").
		Order("road_geom.revision DESC").
		First(&geom).Error
	if err != nil {
		return models.FullGeom{}, err
	}

	return geom, nil
}

func (t *Repository) GetRoadById(roadId int) (models.Road, error) {
	var road models.Road
	if err := t.conn.Where("id = ?", roadId).Find(&road).Error; err != nil {
		return models.Road{}, err
	}
	return road, nil
}

func (t *Repository) GetRoadPreloadById(roadId int) (models.RoadPreload, error) {
	var road models.RoadPreload
	if err := t.conn.Preload("RefDirection").Where("id = ?", roadId).Find(&road).Error; err != nil {
		return models.RoadPreload{}, err
	}
	return road, nil
}

func (t *Repository) GetRoadRefactiveStripByIdParent(idParent int) (models.RoadRetroReflectivity, error) {
	var rrs models.RoadRetroReflectivity
	if err := t.conn.Where("id_parent = ? and status != ?", idParent, "D").Last(&rrs).Order("revision ASC").Error; err != nil {
		return rrs, err
	}
	return rrs, nil
}

func (t *Repository) GetRoadRefactiveStripBeforeLastRevitionByIdParent(idParent int, revision int) (models.RoadRetroReflectivity, error) {
	var rrs models.RoadRetroReflectivity
	if err := t.conn.Where("id_parent = ? and revision = ? and status != ?", idParent, revision, "D").First(&rrs).Order("revision ASC").Error; err != nil {
		return rrs, err
	}
	return rrs, nil
}
func (t *Repository) GetAllRoadRetroReflectivityByIdParent(idParent int) (models.RoadRetroReflectivityPreload, error) {
	var rrs models.RoadRetroReflectivityPreload

	err := t.conn.
		Preload("RoadRetroReflectivityRanges").
		Preload("RoadRetroReflectivityRanges.RefStripeColor").
		Preload("RoadRetroReflectivityRanges.RefStripeType").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs.RefStripeColor").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs.RefStripeType").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs").
		Where("id_parent = ? AND status = ?", idParent, "A").
		Find(&rrs).Error // Use Find to retrieve the data and store it in 'roadCondition'.

	// Check for errors in querying the database and return them if present.
	if err != nil {
		return rrs, err // Return the zero value of roadCondition and the error.
	}

	// If no error, return the retrieved roadCondition data and nil for the error.
	return rrs, nil
}

func (t *Repository) GetRoadRetroReflectivityDetailsByIdParent(idParent int, refStripeTypeIDs []string) (models.RoadRetroReflectivityPreload, error) {
	var rrs models.RoadRetroReflectivityPreload
	test := len(refStripeTypeIDs)
	fmt.Println("refStripeTypeIDs", refStripeTypeIDs)
	fmt.Println("test", test)

	query := t.conn.
		Preload("RoadRetroReflectivityRanges", func(db *gorm.DB) *gorm.DB {
			// Apply condition only if refStripeTypeIDs is not empty.
			if len(refStripeTypeIDs) > 0 {
				db = db.Where("ref_stripe_type_id IN (?)", refStripeTypeIDs)
			}
			return db
		}).
		Preload("RoadRetroReflectivityRanges.RefStripeColor").
		Preload("RoadRetroReflectivityRanges.RefStripeType").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs.RefStripeColor").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs.RefStripeType").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs").
		Where("id_parent = ? AND status = ?", idParent, "A")

	// Execute the query.
	err := query.Find(&rrs).Error
	if err != nil {
		// Return early if there's an error.
		return models.RoadRetroReflectivityPreload{}, err
	}

	// Return the result and nil error if the operation was successful.
	return rrs, nil
}

func (t *Repository) GetRoadRetroReflectivityByIdParent(idParent int) (models.RoadRetroReflectivity, error) {
	var rrs models.RoadRetroReflectivity
	if err := t.conn.Where("id_parent = ? and status != ?", idParent, "D").Last(&rrs).Order("revision ASC").Error; err != nil {
		return rrs, err
	}
	return rrs, nil
}

func (t *Repository) GetLastRevisionByIdParent(idParent int) (models.RoadRetroReflectivity, error) {
	var rrs models.RoadRetroReflectivity
	if err := t.conn.Where("id_parent = ? ", idParent).Last(&rrs).Error; err != nil {
		return rrs, err
	}
	return rrs, nil
}

func (t *Repository) CreateRoadRetroReflectivity(tx *gorm.DB, rc models.RoadRetroReflectivity) (int, error) {

	err := tx.Create(&rc).Error
	if err != nil {
		return 0, err
	}

	if err := tx.Model(&rc).Update("IDParent", rc.ID).Error; err != nil {
		return 0, err
	}

	return rc.ID, nil

}

func (t *Repository) CreateRoadRetroReflectivityRange(tx *gorm.DB, rsrRange models.RoadRetroReflectivityRange) (int, error) {

	// err := tx.Create(&rcServay).Error
	// if err != nil {
	// 	return 0, err
	// }
	var id int
	smt := fmt.Sprintf(`
	INSERT INTO road_retro_reflectivity_range
	(road_retro_reflectivity_id, km_start, km_end, retro_min, retro_max, retro_avg,ref_stripe_color_id,ref_stripe_type_id,the_geom)
	VALUES (%d,%f,%f, %s, %s, %s,%d,%d, %s) RETURNING id`,

		rsrRange.RoadRetroReflectivityID,
		rsrRange.KmStart,
		rsrRange.KmEnd,
		helpers.ConvertNullableFloat64(rsrRange.RetroMin), // Same note as above
		helpers.ConvertNullableFloat64(rsrRange.RetroMax),
		helpers.ConvertNullableFloat64(rsrRange.RetroAvg),
		rsrRange.RefStripeColorID,
		rsrRange.RefStripeTypeID,
		rsrRange.TheGeom, // Make sure this is correctly formatted or converted
	)
	err := tx.Raw(smt).Scan(&id).Error
	if err != nil {
		return 0, err
	}

	return id, nil

}

func (t *Repository) CreateRoadRetroReflectivityM(tx *gorm.DB, rsrM models.RoadRetroReflectivityM) error {

	// err := tx.Create(&rcServay100M).Error
	// if err != nil {
	// 	return 0, err
	// }

	smt := fmt.Sprintf(`
	INSERT INTO road_retro_reflectivity_m
	(road_retro_reflectivity_range_id, km_start, km_end, retro_min, retro_max, retro_avg, ref_stripe_color_id,ref_stripe_type_id,the_geom)
	VALUES (%d,%f,%f, %s, %s, %s, %d,%d, %s)`,

		rsrM.RoadRetroReflectivityRangeID,
		rsrM.KmStart,
		rsrM.KmEnd,
		helpers.ConvertNullableFloat64(rsrM.RetroMin), // Same note as above
		helpers.ConvertNullableFloat64(rsrM.RetroMax),
		helpers.ConvertNullableFloat64(rsrM.RetroAvg),
		rsrM.RefStripeColorID,
		rsrM.RefStripeTypeID,
		rsrM.TheGeom, // Make sure this is correctly formatted or convertedconverted
	)
	if err := tx.Exec(smt).Error; err != nil {
		return err
	}

	return nil

}

// ใช้ insert แบบนี้เพื่อที่จะสามารถ insert the_geom ได้
func (t *Repository) UpdateRoadRetroReflectivityM(tx *gorm.DB, rrtM models.RoadRetroReflectivityM) error {

	smt := `
		INSERT INTO road_condition_survey_m
		(    road_condition_survey_100m_id,
			 km_start,
			 km_end,
			 iri,
			 rut,
			 mpd,
			 ifi,
			 the_geom,
			 img_filepath,
			 survey_type
		 ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,?)`

	if err := tx.Exec(smt,
		rrtM.RoadRetroReflectivityRangeID,
		rrtM.KmStart,
		rrtM.KmEnd,
		helpers.ConvertNullableFloat64(rrtM.RetroMin),
		helpers.ConvertNullableFloat64(rrtM.RetroMax),
		helpers.ConvertNullableFloat64(rrtM.RetroAvg),
		rrtM.TheGeom,
		rrtM.RefStripeTypeID,
		rrtM.RefStripeColorID).Error; err != nil {
		return err
	}
	return nil
}

func (t *Repository) UpdateRoadRetroReflectivity(tx *gorm.DB, rrs models.RoadRetroReflectivity) (models.RoadRetroReflectivity, error) {
	if err := tx.Create(&rrs).Error; err != nil {
		return rrs, err
	}
	return rrs, nil
}

func (t *Repository) UpdateRoadRetroReflectivityNoIriFile(tx *gorm.DB, rrs models.RoadRetroReflectivity) (models.RoadRetroReflectivity, error) {

	// Update only the specified fields. Zero-value fields will be ignored.
	if err := tx.Model(&models.RoadRetroReflectivity{}).Where("id = ?", rrs.ID).Updates(&rrs).Error; err != nil {
		return models.RoadRetroReflectivity{}, err
	}

	return models.RoadRetroReflectivity{}, nil

}

func (t *Repository) DeleteRoadRetroReflectivity(rrs models.RoadRetroReflectivity) error {
	if err := t.conn.Save(&rrs).Error; err != nil {
		return err
	}
	return nil
}

func (t *Repository) DeleteRoadRetroReflectivityByID(ID int, UserID int) error {

	rc := models.RoadRetroReflectivity{
		ID:          ID,
		Status:      "D",
		UpdatedDate: time.Now(),
		UpdatedBy:   UserID,
	}

	if err := t.conn.Model(&models.RoadRetroReflectivity{}).Updates(rc).Error; err != nil {
		return err
	}
	return nil

}

func (t *Repository) DeleteRoadConditionByIDParent(idParent int) error {

	rc := models.RoadCondition{
		Status: "D",
	}

	if err := t.conn.Model(&models.RoadCondition{}).Where("id_parent = ?", idParent).Updates(rc).Error; err != nil {
		return err
	}
	return nil

}
func (t *Repository) GetRoadConditionTemplate(roadId int) (models.RoadConditionTemplate, error) {
	var roadCondition models.RoadConditionTemplate
	if err := t.conn.Preload("RoadInfo").Preload("RoadGeom").Where("id = ? AND is_active = ? ", roadId, true).Find(&roadCondition).Order("road_code ASC,revision DESC ").Error; err != nil {
		return roadCondition, err
	}

	return roadCondition, nil
}

func (t *Repository) GetRoadKmRange(roadID int, kmStart int, kmEnd int, direction int) ([]models.RoadKmRage, error) {
	var kmRanges []models.RoadKmRage

	sql := "SELECT * FROM get_km_range(?, ?, ?, ?, ?)"
	err := t.conn.Raw(sql, roadID, kmStart, kmEnd, 25, direction).Scan(&kmRanges).Error
	if err != nil {
		return nil, err
	}

	return kmRanges, nil
}

func (t *Repository) GetRoadRetroReflectivityCompare(roadId int, req models.RoadRetroReflectivityCompareLine) ([]models.RoadRetroReflectivityPreload, error) {
	var rrt []models.RoadRetroReflectivityPreload

	if err := t.conn.
		Preload("RoadRetroReflectivityRanges").
		Preload("RoadRetroReflectivityRanges.RefStripeColor").
		Preload("RoadRetroReflectivityRanges.RefStripeType").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs.RefStripeColor").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs.RefStripeType").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs").
		Where("status = ? AND road_id = ? AND year IN (?) AND line_no IN (?) ", "A", roadId, req.Years, req.Lines).Order("year DESC, surveyed_date DESC, line_no asc,revision DESC,km_start ASC").Find(&rrt).Error; err != nil {
		return rrt, err
	}
	return rrt, nil
}

func (t *Repository) GetRoadRetroReflectivityAverage(roadID int, Line int) ([]models.RoadRetroReflectivityPreload, error) {

	var rrt []models.RoadRetroReflectivityPreload

	if err := t.conn.
		Preload("RoadRetroReflectivityRanges").
		Preload("RoadRetroReflectivityRanges.RefStripeColor").
		Preload("RoadRetroReflectivityRanges.RefStripeType").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs.RefStripeColor").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs.RefStripeType").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs").
		Where("status = ? AND road_id = ? AND line_no = ? ", "A", roadID, Line).Order("year DESC, surveyed_date DESC").Find(&rrt).Error; err != nil {
		return rrt, err
	}

	return rrt, nil
}

func (t *Repository) UpdateRoadRetroReflectivityFilepath(tx *gorm.DB, rcID int, csvFilepath string) error {
	var rrs models.RoadRetroReflectivity
	if err := tx.Model(rrs).Where("id = ?", rcID).Updates(models.RoadRetroReflectivity{InputFilePath: csvFilepath}).Error; err != nil {
		return err
	}
	return nil
}

func (t *Repository) StartTransSection() *gorm.DB {
	tx := t.conn.Begin()
	return tx
}

func (t *Repository) RollBack(tx *gorm.DB) error {
	tx.Rollback()
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func (t *Repository) Commit(tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (t *Repository) UpdateStatusIByID(tx *gorm.DB, ID int, UserID uint) error {

	if err := tx.Model(&models.RoadRetroReflectivity{}).Where("id = ?", ID).Updates(map[string]interface{}{
		"status":       "I",
		"updated_by":   UserID,
		"updated_date": time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil

}

func (t *Repository) GetRoadSectionDataID(roadID int) (models.RoadInfoGeom, error) {
	var roadInfoGeom models.RoadInfoGeom
	query := t.conn
	query = query.Preload("RoadInfo")
	query = query.Preload("Direction", func(db *gorm.DB) *gorm.DB {
		return db.Order("id")
	})
	query = query.Preload("RoadGeom", func(db *gorm.DB) *gorm.DB {
		return db.Order("lane_no")
	})
	if err := query.Where("id = ?", roadID).Where("is_active = ?", true).First(&roadInfoGeom).Error; err != nil {
		return roadInfoGeom, err
	}
	return roadInfoGeom, nil
}

func (t *Repository) GetRoadByID(roadID int) (*models.RoadById, error) {
	var road models.RoadById
	query := t.conn
	query = query.Where("id = ?", roadID)
	query = query.Preload("RoadInfo", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,road_id, year, name, km_start, km_end,   ref_direction_id, st_astext(the_geom) as the_geom, revision, status, ramp_id,road_color_code, created_by,created_at,updated_by,updated_at,ref_road_type_id", "center_lane_shape_file_path", "center_line_shape_file_path", "remark").Where("status = ?", "A")
	})
	query = query.Preload("RoadSurfaceIcon", func(db *gorm.DB) *gorm.DB {
		return db.Select("DISTINCT ON (ref_surface.surface_group) ref_surface.surface_group, ref_surface.id , road_surface.road_id, ref_surface.surface_group as name, CASE WHEN ref_surface.surface_group = 'Concrete' THEN '#398BF7' ELSE '#7460EE' END AS color_code").
			Joins("RIGHT JOIN road_surface_lane ON road_surface.id = road_surface_lane.road_surface_id").
			Joins("LEFT JOIN ref_surface ON ref_surface.id = road_surface_lane.ref_surface_id").
			Where("road_surface.status = ? and road_surface_lane.ref_surface_id <> ?", "A", 0)
	})
	query = query.Preload("RoadInfo.User")
	query = query.Preload("RoadInfo.RefDirection")
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

func (t *Repository) GetRoadDirectionByRoadID(roadID int) (models.RoadPreloadConditionAll, error) {

	var roadConditions models.RoadPreloadConditionAll
	query := t.conn
	query = query.Preload("RefDirection")
	err := query.Where("id = ? AND is_active = ?", roadID, true).First(&roadConditions).Error
	if err != nil {
		return roadConditions, err
	}
	return roadConditions, nil

}

func (r *Repository) GetRefStripeColor() ([]models.RefStripeColor, error) {
	var refStripeColor []models.RefStripeColor

	query := r.conn
	err := query.Find(&refStripeColor).Error
	if err != nil {
		return nil, err
	}

	return refStripeColor, nil

}

func (r *Repository) GetRefStripeType() ([]models.RefStripeType, error) {
	var refStripeType []models.RefStripeType

	query := r.conn
	err := query.Find(&refStripeType).Error
	if err != nil {
		return nil, err
	}

	return refStripeType, nil

}

func (t *Repository) GetRoadDirectionLaneList(roadID int) (models.RoadInfoGeomDirection, error) {
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
	helpers.PrintlnJson(roadInfoGeomDirection)
	return roadInfoGeomDirection, nil
}

func (t *Repository) GetTotalLanesByRoadID(roadID int) (int64, error) {
	var roadLaneCount int64
	query := t.conn.Model(&models.RoadGeom{})
	if err := query.Where("road_id = ? AND status = ?", roadID, "A").Count(&roadLaneCount).Error; err != nil {
		return 0, err
	}
	return roadLaneCount, nil
}

func (t *Repository) GetLineListByRoadID(roadId int) ([]models.RoadRetroReflectivityPreload, error) {
	var rrt []models.RoadRetroReflectivityPreload

	if err := t.conn.Select("DISTINCT line_no").
		Where("status = ? AND road_id = ?", "A", roadId).Order("line_no asc").Find(&rrt).Error; err != nil {
		return rrt, err
	}
	return rrt, nil
}
