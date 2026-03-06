package repositories

import (
	"fmt"
	"strings"
	"time"

	helpers "gitlab.com/mims-api-service/helpers"
	models "gitlab.com/mims-api-service/models"
	requests "gitlab.com/mims-api-service/requests"
	responses "gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadCondition/handlers"
	"gitlab.com/mims-api-service/src/roadCondition/usecases"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type roadConditionRepository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewRoadConditionRepositoryHandler(conn *gorm.DB) *handlers.RoadConditionHandler {
	useCase := usecases.NewRoadConditionUseCase(&roadConditionRepository{conn})
	handler := handlers.NewRoadConditionHandler(useCase)
	return handler
}

// /////////////////////////
func (t *roadConditionRepository) GetRole(userId uint) ([]models.UserRole, error) {
	var userRole []models.UserRole
	if err := t.conn.Where("user_id = ?", userId).Find(&userRole).Error; err != nil {
		return userRole, err
	}
	return userRole, nil
}

func (t *roadConditionRepository) GetAccessControl(roles []int) ([]models.AccessControl, error) {
	var accessControl []models.AccessControl
	query := t.conn
	query = query.Joins("JOIN role_access_control on access_control.id = role_access_control.access_control_id")
	err := query.Find(&accessControl).Error
	if err != nil {
		fmt.Println(err)
	}
	return accessControl, nil
}

func (t *roadConditionRepository) GetUserById(userId uint) (models.Users, error) {
	var user models.Users
	if err := t.conn.Where("id = ?", userId).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (t *roadConditionRepository) GetUserDepartmentById(userId int) (models.UserDepartment, error) {
	var user models.UserDepartment
	if err := t.conn.Where("id = ?", userId).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (t *roadConditionRepository) GetRoadDetailMenu(accessKeys []string, departmentId int) ([]responses.RoadMenuData, error) {
	var data []responses.RoadMenuData
	// if helpers.HasPermission([]string{"edit_asset_in data", "approve_asset_in_data", "edit_asset_out_data", "approve_asset_out_data"}, accessKeys) {
	query := t.conn
	query = query.Select("DISTINCT ref_asset.id as group_id, ref_asset.name as group_name, ref_asset_table.id as asset_id, ref_asset_table.table_label as asset_name, ref_asset_table.geom_type as geom_type, ref_asset_table.is_in_road as is_in_road, ref_asset_table.icon_filepath as icon_filepath,	ref_asset.seq AS a_seq, ref_asset_table.seq AS as_seq")
	query = query.Table("ref_asset")
	query = query.Joins("JOIN ref_asset_table on ref_asset.id = ref_asset_table.ref_asset_id")

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

func (t *roadConditionRepository) GetRoadGroupList(prams requests.RoadPrams) []models.RoadList {
	var roadList []models.RoadList
	query := t.conn
	// ChildRoad.RoadInfo
	query = query.Preload("ChildRoad.RoadInfo", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,road_id, year, name, km_start, km_end,   ref_direction_id, st_astext(the_geom) as the_geom, revision, status, st_astext(old_the_geom) as old_the_geom, old_km_start, old_km_end, ramp_id,road_color_code, created_by,created_date,updated_by,updated_date")
	})

	//RoadInfo
	query = query.Preload("RoadInfo", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,road_id, year, name, km_start, km_end,   ref_direction_id, st_astext(the_geom) as the_geom, revision, status, st_astext(old_the_geom) as old_the_geom, old_km_start, old_km_end, ramp_id,road_color_code, created_by,created_date,updated_by,updated_date")
	})

	query = query.Preload("ChildRoad", func(db *gorm.DB) *gorm.DB {
		return db.Order("seq")
	})

	query = query.Preload("Direction")
	query = query.Preload("ChildRoad.Direction")
	query = query.Preload("ChildRoad.RoadType")
	query = query.Preload("RoadType")

	query = query.Preload("RoadSurfaceData.RoadSurfaceLane")
	query = query.Preload("ChildRoad.RoadSurfaceData.RoadSurfaceLane")
	query = query.Preload("RoadDamage", func(db *gorm.DB) *gorm.DB {
		var wheres []string
		wheres = append(wheres, "status = 'W'")
		wheres = append(wheres, "status = 'R'")
		queryWhereString := strings.Join(wheres, " or ")
		return db.Where(queryWhereString)
	})
	query = query.Preload("ChildRoad.RoadDamage", func(db *gorm.DB) *gorm.DB {
		var wheres []string
		wheres = append(wheres, "status = 'W'")
		wheres = append(wheres, "status = 'R'")
		queryWhereString := strings.Join(wheres, " or ")
		return db.Where(queryWhereString)
	})
	query = query.Preload("RoadCondition", func(db *gorm.DB) *gorm.DB {
		var wheres []string
		wheres = append(wheres, "status = 'W'")
		wheres = append(wheres, "status = 'R'")
		queryWhereString := strings.Join(wheres, " or ")
		return db.Where(queryWhereString)
	})
	query = query.Preload("ChildRoad.RoadCondition", func(db *gorm.DB) *gorm.DB {
		var wheres []string
		wheres = append(wheres, "status = 'W'")
		wheres = append(wheres, "status = 'R'")
		queryWhereString := strings.Join(wheres, " or ")
		return db.Where(queryWhereString)
	})
	query = query.Preload("RoadSurface", func(db *gorm.DB) *gorm.DB {
		var wheres []string
		wheres = append(wheres, "status = 'W'")
		wheres = append(wheres, "status = 'R'")
		queryWhereString := strings.Join(wheres, " or ")
		return db.Where(queryWhereString)
	})
	query = query.Preload("ChildRoad.RoadSurface", func(db *gorm.DB) *gorm.DB {
		var wheres []string
		wheres = append(wheres, "status = 'W'")
		wheres = append(wheres, "status = 'R'")
		queryWhereString := strings.Join(wheres, " or ")
		return db.Where(queryWhereString)
	})

	query = query.Preload("RoadAsset", func(db *gorm.DB) *gorm.DB {
		var wheres []string
		wheres = append(wheres, "status = 'W'")
		wheres = append(wheres, "status = 'R'")
		queryWhereString := strings.Join(wheres, " or ")
		return db.Where(queryWhereString)
	})
	query = query.Preload("ChildRoad.RoadAsset", func(db *gorm.DB) *gorm.DB {
		var wheres []string
		wheres = append(wheres, "status = 'W'")
		wheres = append(wheres, "status = 'R'")
		queryWhereString := strings.Join(wheres, " or ")
		return db.Where(queryWhereString)
	})

	query = query.Preload("RoadGeom", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,road_id, lane_no, km_start,km_end,ST_ASTEXT(ST_FORCE2D(the_geom)) as the_geom,revision,status,old_km_start,old_km_end,old_the_geom,remark,created_by,created_date,updated_by,updated_date").Where("status = ?", "A").Order("revision DESC").Order("revision")
	})
	query = query.Preload("ChildRoad.RoadGeom", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,road_id, lane_no, km_start,km_end,ST_ASTEXT(ST_FORCE2D(the_geom)) as the_geom,revision,status,old_km_start,old_km_end,old_the_geom,remark,created_by,created_date,updated_by,updated_date").Where("status = ?", "A").Order("revision DESC").Order("revision")
	})

	if err := query.Order("id").Find(&roadList).Error; err != nil {
		// fmt.Println(err)
		// return user, err

	}
	return roadList
}

func (t *roadConditionRepository) GetRoadGroup() ([]models.RoadGroup, error) {
	var roadGroup []models.RoadGroup
	if err := t.conn.Find(&roadGroup).Error; err != nil {
		return roadGroup, err
	}
	return roadGroup, nil
}

func (t *roadConditionRepository) GetRoadStatusSurface(roadId int) (responses.StatusCount, error) {

	var statusCount responses.StatusCount
	statusCount.Category = "SURFACE"
	query := t.conn
	query = query.Select("SUM(CASE WHEN status = 'T' THEN 1 ELSE 0 END) as count_temp, SUM(CASE WHEN status = 'W' THEN 1 ELSE 0 END) as count_waiting, SUM(CASE WHEN status = 'R' THEN 1 ELSE 0 END) as count_rejected")
	query = query.Table("road_surface")
	if err := query.Where("road_id = ?", roadId).Find(&statusCount).Error; err != nil {
		return statusCount, err
	}
	return statusCount, nil
}

// func (t *roadConditionRepository) GetRoadStatusAsset(roadId int) ([]responses.StatusCount, error) {
// 	// helpers.PrintlnJson(data)
// 	var statusCount []responses.StatusCount
// 	query2 := t.conn
// 	if err := query2.Raw("SELECT  ref_asset_table.id as id, ref_asset_table.ref_asset_id as group_id, road_asset.status as status, COUNT (CASE WHEN road_asset.status = 'T' THEN 1 ELSE NULL END) as count_temp, COUNT (CASE WHEN road_asset.status = 'W' THEN 1 ELSE NULL END ) as count_waiting, COUNT (CASE WHEN road_asset.status = 'R' THEN 1 ELSE NULL END ) as count_rejected FROM road_asset INNER JOIN ref_asset_table ON road_asset.ref_asset_table_id = ref_asset_table.id INNER JOIN ref_asset_table_staff ON ref_asset_table_staff.ref_asset_table_id = ref_asset_table.id WHERE (road_asset.road_id = 1) AND (ref_asset_table.is_in_road) AND (ref_asset_table.is_active) AND (road_asset.status = 'T' or road_asset.status = 'W' or road_asset.status = 'R') AND (ref_asset_table_staff.is_approver = true) AND (ref_asset_table_staff.ref_department_id = 1) GROUP BY ref_asset_table.id, ref_asset_table.ref_asset_id,road_asset.status").Find(&statusCount).Error; err != nil {
// 		return statusCount, err
// 	}
// 	return statusCount, nil
// }

func (t *roadConditionRepository) GetRoadTypeIcon() ([]models.RefRoadTypeIcon, error) {
	var roadTypeIcon []models.RefRoadTypeIcon

	if err := t.conn.Find(&roadTypeIcon).Error; err != nil {
		return roadTypeIcon, err
	}
	return roadTypeIcon, nil
}

func (t *roadConditionRepository) GetRoadConditionList(rcStatus string, roadId int) ([]models.RoadConditionList, error) {
	var rcList []models.RoadConditionList
	err := t.conn.Table("road_info AS ri").
		Select("DISTINCT ri.road_id AS road_id, rc.year AS year, rc.id AS id, rc.id_parent AS id_parent, rc.revision AS revision, ref_d.id AS direction_id, ref_d.name AS direction_name, rc.lane_no AS lane_no, rc.surveyed_date AS surveyed_date").
		Joins("LEFT JOIN ref_direction ref_d ON ri.ref_direction_id = ref_d.id").
		Joins("LEFT JOIN road_condition rc ON ri.road_id = rc.road_id").
		Where(rcStatus).
		Where("ri.road_id = ?", roadId).
		Order("year DESC, lane_no ASC, surveyed_date DESC, revision DESC").
		Find(&rcList).Error
	if err != nil {
		return nil, err
	}
	return rcList, nil

}

// func (t *roadConditionRepository) GetRoadConditionDetails(idParent int, rcStatus models.SqlCondition, sqlCondition models.SqlCondition) ([]models.RoadConditionDetails, error) {

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

// func (t *roadConditionRepository) GetRoadConditionGrades(roadId int) ([]models.RoadConditionGrade, error) {
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

func (t *roadConditionRepository) GetFullGeom(roadId int, laneNo int) (models.FullGeom, error) {
	var geom models.FullGeom

	err := t.conn.Where("road_geom.road_id = ? AND road_geom.lane_no = ? AND road_geom.status = ?", roadId, laneNo, "A").
		Order("road_geom.revision DESC").
		First(&geom).Error
	if err != nil {
		return models.FullGeom{}, err
	}

	return geom, nil
}

func (t *roadConditionRepository) GetRoadById(roadId int) (models.Road, error) {
	var road models.Road
	if err := t.conn.Where("id = ?", roadId).Find(&road).Error; err != nil {
		return models.Road{}, err
	}
	return road, nil
}

func (t *roadConditionRepository) GetRoadPreloadById(roadId int) (models.RoadPreload, error) {
	var road models.RoadPreload
	if err := t.conn.Preload("RefDirection").Where("id = ?", roadId).Find(&road).Error; err != nil {
		return models.RoadPreload{}, err
	}
	return road, nil
}

func (t *roadConditionRepository) GetRoadConditionByIdParent(idParent int) (models.RoadCondition, error) {
	var roadCondition models.RoadCondition
	if err := t.conn.Where("id_parent = ? and status != ?", idParent, "D").Last(&roadCondition).Order("revision ASC").Error; err != nil {
		return roadCondition, err
	}
	return roadCondition, nil
}

func (t *roadConditionRepository) GetRoadConditionBeforeLastRevitionByIdParent(idParent int, revision int) (models.RoadCondition, error) {
	var roadCondition models.RoadCondition
	if err := t.conn.Where("id_parent = ? and revision = ? and status != ?", idParent, revision, "D").First(&roadCondition).Order("revision ASC").Error; err != nil {
		return roadCondition, err
	}
	return roadCondition, nil
}
func (t *roadConditionRepository) GetAllRoadConditionByIdParent(idParent int) (models.RoadConditionAll, error) {
	var roadCondition models.RoadConditionAll
	err := t.conn.
		Preload("RoadConditionSurveys").
		Preload("RoadConditionSurveys.RoadConditionSurvey100Ms").
		Preload("RoadConditionSurveys.RoadConditionSurvey100Ms.RoadConditionSurveyMs").
		Where("id_parent = ? AND status = ?", idParent, "A").
		Find(&roadCondition).Error // Use Find to retrieve the data and store it in 'roadCondition'.

	// Check for errors in querying the database and return them if present.
	if err != nil {
		return roadCondition, err // Return the zero value of roadCondition and the error.
	}

	// If no error, return the retrieved roadCondition data and nil for the error.
	return roadCondition, nil
}

func (t *roadConditionRepository) CreateRoadCondition(tx *gorm.DB, rc models.RoadCondition) (int, error) {

	err := tx.Create(&rc).Error
	if err != nil {
		return 0, err
	}

	if err := tx.Model(&rc).Update("IDParent", rc.ID).Error; err != nil {
		return 0, err
	}

	return rc.ID, nil

}

func (t *roadConditionRepository) CreateRoadConditionSurvey(tx *gorm.DB, rcServay models.RoadConditionSurvey) (int, error) {

	// err := tx.Create(&rcServay).Error
	// if err != nil {
	// 	return 0, err
	// }
	var id int
	smt := fmt.Sprintf(`
	INSERT INTO road_condition_survey
	(road_condition_id, km_start, km_end, iri, rut, mpd, ifi, survey_type,the_geom)
	VALUES (%d,%f,%f, %s, %s, %s, %s,'%s', %s) RETURNING id`,

		rcServay.RoadConditionId,
		rcServay.KmStart,
		rcServay.KmEnd,
		helpers.ConvertNullableFloat64(rcServay.IRI), // Same note as above
		helpers.ConvertNullableFloat64(rcServay.RUT),
		helpers.ConvertNullableFloat64(rcServay.MPD),
		helpers.ConvertNullableFloat64(rcServay.IFI),
		rcServay.SurveyType,
		rcServay.TheGeom, // Make sure this is correctly formatted or converted
	)
	err := tx.Raw(smt).Scan(&id).Error
	if err != nil {
		return 0, err
	}

	return id, nil

}

func (t *roadConditionRepository) CreateRoadConditionSurvey100M(tx *gorm.DB, rcSurvey100M models.RoadConditionSurvey100M) (int, error) {

	// err := tx.Create(&rcServay100M).Error
	// if err != nil {
	// 	return 0, err
	// }

	var id int
	smt := fmt.Sprintf(`
	INSERT INTO road_condition_survey_100m
	(road_condition_survey_id, km_start, km_end, iri, rut, mpd, ifi, survey_type,the_geom)
	VALUES (%d,%f,%f, %s, %s, %s, %s,'%s', %s) RETURNING id`,

		rcSurvey100M.RoadConditionSurveyId,
		rcSurvey100M.KmStart,
		rcSurvey100M.KmEnd,
		helpers.ConvertNullableFloat64(rcSurvey100M.IRI), // Same note as above
		helpers.ConvertNullableFloat64(rcSurvey100M.RUT),
		helpers.ConvertNullableFloat64(rcSurvey100M.MPD),
		helpers.ConvertNullableFloat64(rcSurvey100M.IFI),
		rcSurvey100M.SurveyType,
		rcSurvey100M.TheGeom, // Make sure this is correctly formatted or converted
	)
	err := tx.Raw(smt).Scan(&id).Error
	if err != nil {
		return 0, err
	}

	return id, nil

}

// ใช้ insert แบบนี้เพื่อที่จะสามารถ insert the_geom ได้
func (t *roadConditionRepository) CreateRoadConditionSurveyM(tx *gorm.DB, rcServayM models.RoadConditionSurveyM) error {

	smt := fmt.Sprintf(`
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
	 ) VALUES (%d,%f,%f, %s, %s, %s, %s, %s,'%s','%s')`,
		rcServayM.RoadConditionSurvay100mID,
		rcServayM.KmStart,
		rcServayM.KmEnd,
		helpers.ConvertNullableFloat64(rcServayM.IRI),
		helpers.ConvertNullableFloat64(rcServayM.RUT),
		helpers.ConvertNullableFloat64(rcServayM.MPD),
		helpers.ConvertNullableFloat64(rcServayM.IFI),
		rcServayM.TheGeom,
		rcServayM.ImgFilepath,
		rcServayM.SurveyType)
	if err := tx.Exec(smt).Error; err != nil {
		return err
	}

	return nil
}

// ใช้ insert แบบนี้เพื่อที่จะสามารถ insert the_geom ได้
func (t *roadConditionRepository) UpdateRoadConditionServeyM(tx *gorm.DB, rcServayM models.RoadConditionSurveyM) error {

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
		rcServayM.RoadConditionSurvay100mID,
		rcServayM.KmStart,
		rcServayM.KmEnd,
		helpers.ConvertNullableFloat64(rcServayM.IRI),
		helpers.ConvertNullableFloat64(rcServayM.RUT),
		helpers.ConvertNullableFloat64(rcServayM.MPD),
		helpers.ConvertNullableFloat64(rcServayM.IFI),
		rcServayM.TheGeom,
		rcServayM.ImgFilepath,
		rcServayM.SurveyType).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadConditionRepository) UpdateRoadCondition(tx *gorm.DB, rc models.RoadCondition) (models.RoadCondition, error) {
	if err := tx.Create(&rc).Error; err != nil {
		return rc, err
	}
	return rc, nil
}

func (t *roadConditionRepository) UpdateRoadConditionNoIriFile(tx *gorm.DB, rc models.RoadCondition) (models.RoadCondition, error) {

	// Update only the specified fields. Zero-value fields will be ignored.
	if err := tx.Model(&models.RoadCondition{}).Where("id = ?", rc.ID).Updates(rc).Error; err != nil {
		return models.RoadCondition{}, err
	}

	return models.RoadCondition{}, tx.Commit().Error

}

func (t *roadConditionRepository) DeleteRoadCondition(rc models.RoadCondition) error {
	if err := t.conn.Save(&rc).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadConditionRepository) DeleteRoadConditionByID(ID int, UserID int) error {

	rc := models.RoadCondition{
		ID:          ID,
		Status:      "D",
		UpdatedDate: time.Now(),
		UpdatedBy:   UserID,
	}

	if err := t.conn.Model(&models.RoadCondition{}).Updates(rc).Error; err != nil {
		return err
	}
	return nil

}

func (t *roadConditionRepository) DeleteRoadConditionByIDParent(idParent int) error {

	rc := models.RoadCondition{
		Status: "D",
	}

	if err := t.conn.Model(&models.RoadCondition{}).Where("id_parent = ?", idParent).Updates(rc).Error; err != nil {
		return err
	}
	return nil

}

func (t *roadConditionRepository) GetRoadConditionTemplate(roadId int) (models.RoadConditionTemplate, error) {
	var roadCondition models.RoadConditionTemplate
	if err := t.conn.Preload("RoadInfo").Preload("RoadSection").Preload("RoadSection.RoadGroup").Preload("RoadGeom").Where("id = ? AND is_active = ? ", roadId, true).Find(&roadCondition).Order("road_code ASC,revision DESC ").Error; err != nil {
		return roadCondition, err
	}

	return roadCondition, nil
}

func (t *roadConditionRepository) GetRoadKmRange(roadID int, kmStart int, kmEnd int, direction int) ([]models.RoadKmRage, error) {
	var kmRanges []models.RoadKmRage

	sql := "SELECT * FROM get_km_range(?, ?, ?, ?, ?)"
	err := t.conn.Raw(sql, roadID, kmStart, kmEnd, 25, direction).Scan(&kmRanges).Error
	if err != nil {
		return nil, err
	}

	return kmRanges, nil
}

func (t *roadConditionRepository) GetRoadConditionCompare(roadId int, req models.RoadConditionCompareLane) ([]models.RoadConditionAll, error) {
	var roadCondition []models.RoadConditionAll
	query := t.conn
	if err := query.Preload("RoadConditionSurveys").
		Preload("RoadConditionSurveys.RoadConditionSurvey100Ms").
		Preload("RoadConditionSurveys.RoadConditionSurvey100Ms.RoadConditionSurveyMs").
		Where("status = ? AND road_id = ? AND year IN (?) AND lane_no IN (?) ", "A", roadId, req.Years, req.Lanes).Order("year DESC, surveyed_date DESC, lane_no asc,revision DESC,km_start ASC").Find(&roadCondition).Error; err != nil {
		return roadCondition, err
	}
	return roadCondition, nil
}

func (t *roadConditionRepository) GetRoadConditionCompareAverage(roadId int, Lane int) ([]models.RoadConditionAll, error) {

	var roadCondition []models.RoadConditionAll

	query := t.conn
	// if err := query.Preload("RoadConditionSurveys.RoadConditionSurvey100Ms").
	// 	Preload("RoadConditionSurveys.RoadConditionSurvey100Ms.RoadConditionSurveyMs").
	// 	Where("status = ? AND road_id = ? AND lane_no = ? ", "A", roadId, Lane).Order("year DESC").Find(&roadCondition).Error; err != nil {
	// 	return roadCondition, err
	// }

	if err := query.Where("status = ? AND road_id = ? AND lane_no = ? ", "A", roadId, Lane).Order("year DESC, surveyed_date DESC").Find(&roadCondition).Error; err != nil {
		return roadCondition, err
	}
	return roadCondition, nil
}

func (t *roadConditionRepository) UpdateRoadConditionFilepath(tx *gorm.DB, rcID int, csvFilepath string, imgFilepath string) error {
	var roadCondition models.RoadCondition
	if err := tx.Model(roadCondition).Where("id = ?", rcID).Updates(models.RoadCondition{IRIInputFilePath: csvFilepath, ImgFilePath: imgFilepath}).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadConditionRepository) UpdateRoadConditionImgPath(ID int, imgPath string) error {

	if err := t.conn.Model(&models.RoadCondition{}).Where("id = ?", ID).Updates(map[string]interface{}{
		"img_filepath": imgPath,
	}).Error; err != nil {
		return err
	}
	return nil

}

func (t *roadConditionRepository) UpdateRoadConditionSurveyMImgPath(rcSurveyID int, imgPath string) error {
	var roadConditionSurveyM models.RoadConditionSurveyM
	if err := t.conn.Model(roadConditionSurveyM).Where("road_condition_survey_id = ?", rcSurveyID).Update("img_filepath", imgPath).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadConditionRepository) StartTransSection() *gorm.DB {
	tx := t.conn.Begin()
	return tx
}

func (t *roadConditionRepository) RollBack(tx *gorm.DB) error {
	tx.Rollback()
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func (t *roadConditionRepository) Commit(tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (t *roadConditionRepository) UpdateStatusIByID(tx *gorm.DB, ID int, UserID uint) error {

	if err := tx.Model(&models.RoadCondition{}).Where("id = ?", ID).Updates(map[string]interface{}{
		"status":       "I",
		"updated_by":   UserID,
		"updated_date": time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil

}

func (t *roadConditionRepository) GetRoadByID(roadID int) (models.RoadInfoGeomDirection, error) {
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

func (t *roadConditionRepository) GetRoadDirectionByRoadID(roadID int) (models.RoadPreloadConditionAll, error) {

	var roadConditions models.RoadPreloadConditionAll
	query := t.conn
	query = query.Preload("RefDirection")
	err := query.Where("road_id = ? AND status = ?", roadID, 'A').First(&roadConditions).Error
	if err != nil {
		return roadConditions, err
	}
	return roadConditions, nil

}
