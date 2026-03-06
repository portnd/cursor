package repositories

import (
	"sort"

	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// ////////////// NEW MIMS ////////////////
func (r *Repository) GetInitialSurfaceArrayForRoadCondition() ([]models.Surface, error) {
	var surface []models.Surface
	query := r.conn
	err := query.Table("ref_surface AS ref_s").Select("id , name, color ").Order("id").Scan(&surface).Error
	if err != nil {
		return surface, err
	}
	return surface, nil
}

func (r *Repository) GetDataMartInfoForRoadCondition(roadID int) ([]models.SurfaceInfo, error) {
	var surfaceInfo []models.SurfaceInfo
	query := r.conn

	query = query.Where("data_mart.road_id = ?", roadID)
	err := query.Table("data_mart").
		Select("road_group.name as road_group_name, road_info.name as road_name, data_mart.km_start as km_start,  data_mart.km_end as km_end,  data_mart.road_id, road.road_code, lane_no, ref_surface.id as surface_id, ref_surface.name as surface_name, ref_surface.color as color_code, st_asgeojson(ST_FORCE2D(data_mart.the_geom)) AS geom_cl, road_surface_id, data_mart.lane_count, data_mart.age, data_mart.contract_number, data_mart.year, data_mart.last_inspection_date, ref_surface.surface_group").
		Joins("JOIN ref_surface ON data_mart.ref_surface_id = ref_surface.id").
		Joins("JOIN road ON road.id = data_mart.road_id").
		Joins("JOIN road_info ON road_info.road_id = data_mart.road_id").
		Joins("JOIN road_group ON road.road_group_id = road_group.id").
		Scan(&surfaceInfo).Error
	if err != nil {
		return surfaceInfo, err
	}

	return surfaceInfo, nil
}

func (r *Repository) GetRoadInfoForRoadCondition(roadID int) (*models.RoadReportInfo, error) {
	var data models.RoadReportInfo
	err := r.conn.
		Table("road").
		Select(`
        road_group.name AS road_group_name,
        road_info.name AS road_name,
        road.road_code,
        road_info.km_start,
        road_info.km_end,
		road_info.road_color_code
    	`).
		Joins("JOIN road_group ON road_group.id = road.road_group_id").
		Joins("JOIN road_info ON road_info.road_id = road.id").
		Where("road_id = ?", roadID).
		Find(&data).Error

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *Repository) GetSurfaceForRoadCondition() ([]models.Surface, error) {
	result := []models.Surface{}
	err := r.conn.Order("id ASC").Find(&result).Error
	if err != nil {
		return []models.Surface{}, err
	}
	return result, nil
}

func (r *Repository) GetRoadSectionByIDForRoadCondition(id int) (models.RoadSection, error) {
	result := models.RoadSection{}
	err := r.conn.Where("id = ?", id).First(&result).Error
	if err != nil {
		return models.RoadSection{}, err
	}
	return result, nil
}

func (r *Repository) GetRoadGroupByIDForRoadCondition(id int) (models.RoadGroup, error) {
	result := models.RoadGroup{}
	err := r.conn.Where("id = ?", id).First(&result).Error
	if err != nil {
		return models.RoadGroup{}, err
	}
	return result, nil
}

func (r *Repository) GetRoadFromSectionIDForRoadCondition(SectionID int) ([]int, error) {
	var roadIDs []int
	err := r.conn.Model(&models.Road{}).Where("road_section_id = ?", SectionID).Pluck("id", &roadIDs).Error
	if err != nil {
		return nil, err
	}
	return roadIDs, nil
}

func (r *Repository) GetReportConditionForRoadCondition(year, roadID int, dis string) (*models.DataReportCondition, []models.DataRoadCondition, []models.DataRoadConditionM, error) {
	var data models.DataReportCondition
	var header models.DataReportConditionDetail
	err := r.conn.Debug().
		Table("road").
		Select(`
        road_group.name AS road_group_name,
        road_info.name AS road_name,
        road.road_code,
        road_info.km_start,
        road_info.km_end
    	`).
		Joins("JOIN road_group ON road_group.id = road.road_group_id").
		Joins("JOIN road_info ON road_info.road_id = road.id").
		Where("road.id  = ?", roadID).
		Scan(&header).Error
	if err != nil {
		return nil, nil, nil, err
	}
	copier.Copy(&data, &header)
	var lane []models.DataRoadCondition
	var laneData []models.DataRoadConditionDetail
	err = r.conn.
		Table("road_condition").
		Select("lane_no, km_start, km_end, iri, rut, mpd, ifi").
		Where("status = ? AND year = ? AND road_id = ?", "A", year, roadID).
		Order("lane_no ASC").
		Find(&laneData).Error
	if err != nil {
		return nil, nil, nil, err
	}
	copier.Copy(&lane, &laneData)
	var detail []models.DataRoadConditionM

	if dis == "25" || dis == "100" {
		err = r.conn.
			Table("road_condition").
			Select(`
			road_condition.lane_no, 
			road_condition_survey_m.km_start, 
			road_condition_survey_m.km_end, 
			road_condition_survey_m.iri, 
			road_condition_survey_m.rut, 
			road_condition_survey_m.mpd, 
			road_condition_survey_m.ifi`).
			Joins("JOIN road_condition_survey ON road_condition.id = road_condition_survey.road_condition_id").
			Joins("JOIN road_condition_survey_100m ON road_condition_survey.id = road_condition_survey_100m.road_condition_survey_id").
			Joins("JOIN road_condition_survey_m on road_condition_survey_100m.id  = road_condition_survey_m.road_condition_survey_100m_id").
			Where("status = ? AND year = ? AND road_id = ?", "A", year, roadID).
			Order("road_condition.lane_no ASC, road_condition_survey_m.km_start ASC").
			Find(&detail).Error
		if err != nil {
			return nil, nil, nil, err
		}
	} else if dis == "1000" {
		err = r.conn.
			Table("road_condition").
			Select(`road_condition.lane_no, 
			road_condition_survey.km_start, 
			road_condition_survey.km_end, 
			road_condition_survey.iri, 
			road_condition_survey.rut, 
			road_condition_survey.mpd, 
			road_condition_survey.ifi`).
			Joins("JOIN road_condition_survey ON road_condition.id = road_condition_survey.road_condition_id").
			Where("status = ? AND year = ? AND road_id = ?", "A", year, roadID).
			Order("road_condition.lane_no ASC, road_condition_survey.km_start ASC").
			Find(&detail).Error
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return &data, lane, detail, err
}

func (r *Repository) GetMeasureValueSummaryConditionForRoadCondition(factor, measureID string) ([]models.DataGradeSummaryCondition, error) {
	var resp []models.DataGradeSummaryCondition
	err := r.conn.
		Table("params_condition").
		Select(`
		ref_grade_id,
		left_value,
		right_value
		`).
		Where("params_condition.condition_type = ?", factor).
		Where("ref_owner_id = ?", measureID).
		Order("ref_grade_id ASC").
		Find(&resp).Error
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Repository) GetRoadConditionGradesByIDForRoadCondition(ownerID int, conditionType string) ([]models.ParamsConditionPreload, error) {
	var paramsCondition []models.ParamsConditionPreload
	query := r.conn
	query = query.Preload("RefOwner")
	query = query.Preload("RefGrade")
	query = query.Where("ref_owner_id = ? AND LOWER(condition_type) = ? ", ownerID, conditionType)
	err := query.Find(&paramsCondition).Error
	if err != nil {
		return paramsCondition, err
	}
	sort.Slice(paramsCondition, func(i, j int) bool {
		return paramsCondition[i].RefGradeID < paramsCondition[j].RefGradeID
	})
	return paramsCondition, nil
}

func (r *Repository) GetRoadLineGradesByIDForRoadCondition(ownerID int) ([]models.ParamsRoadLinePreload, error) {
	var paramsRoadLine []models.ParamsRoadLinePreload
	query := r.conn
	query = query.Preload("RefOwnerRoadLine")
	query = query.Preload("RefGrade")
	query = query.Where("ref_owner_road_line_id = ?", ownerID)
	err := query.Find(&paramsRoadLine).Error
	if err != nil {
		return paramsRoadLine, err
	}
	sort.Slice(paramsRoadLine, func(i, j int) bool {
		return paramsRoadLine[i].RefGradeID < paramsRoadLine[j].RefGradeID
	})
	return paramsRoadLine, nil
}

func (r *Repository) GetRoadConditionDashboardForRoadCondition(roadIDs []int, depotCodes []string, filter requests.Condition) ([]models.RoadConditionReport, error) {
	var roadCondition []models.RoadConditionReport

	query := r.conn.Model(&models.RoadConditionReport{})

	query = query.Preload("Road.RoadSection", func(db *gorm.DB) *gorm.DB {
		if len(depotCodes) > 0 {
			return db.Where("ref_depot_code IN (?)", depotCodes)
		}
		return db
	})

	query = query.Preload("RoadConditionSurveys")
	query = query.Preload("RoadConditionSurveys.RoadConditionSurvey100Ms")
	query = query.Preload("RoadConditionSurveys.RoadConditionSurvey100Ms.RoadConditionSurveyMs")

	// Applying filters
	if len(roadIDs) > 0 {
		query = query.Where("road_condition.road_id IN ?", roadIDs)
	}
	if filter.Year != 0 {
		query = query.Where("road_condition.year = ?", filter.Year)
	}
	if filter.KmStart != 0.0 && filter.KmEnd != 0.0 {
		query = query.Where("road_condition.km_start >= ? AND road_condition.km_end <= ?", filter.KmStart, filter.KmEnd)
	} else if filter.KmStart != 0.0 {
		query = query.Where("road_condition.km_start >= ?", filter.KmStart)
	} else if filter.KmEnd != 0.0 {
		query = query.Where("road_condition.km_end <= ?", filter.KmEnd)
	}
	query = query.Where("road_condition.status = ?", "A")

	// Subquery to find the latest road condition per road ID based on created_date
	subquery := r.conn.Table("road_condition").
		Select("road_id, MAX(created_date) AS latest_date").
		Where("status = ?", "A").
		Group("road_id,lane_no")
	if filter.Year != 0 {
		subquery = subquery.Where("year = ?", filter.Year)
	}

	// Main query with subquery to get the latest records
	query = query.Joins("JOIN (?) AS latest ON road_condition.road_id = latest.road_id AND road_condition.created_date = latest.latest_date", subquery)

	// Execute the query
	err := query.Find(&roadCondition).Error
	if err != nil {
		return nil, err
	}

	return roadCondition, nil
}

func (r *Repository) GetRoadRetroReflectivityDashboardForRoadCondition(roadIDs []int, depotCodes []string, filter requests.Condition) ([]models.RoadRetroReflectivityDashboard, error) {
	var retroReflectivity []models.RoadRetroReflectivityDashboard

	query := r.conn.Model(&models.RoadRetroReflectivityDashboard{})

	query = query.Preload("Road.RoadSection", func(db *gorm.DB) *gorm.DB {
		if len(depotCodes) > 0 {
			return db.Where("ref_depot_code IN (?)", depotCodes)
		}
		return db
	})

	query = query.Preload("RoadRetroReflectivityRanges")
	query = query.Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs")

	// Applying filters
	if len(roadIDs) > 0 {
		query = query.Where("road_retro_reflectivity.road_id IN ?", roadIDs)
	}
	if filter.Year != 0 {
		query = query.Where("road_retro_reflectivity.year = ?", filter.Year)
	}
	if filter.KmStart != 0.0 && filter.KmEnd != 0.0 {
		query = query.Where("road_retro_reflectivity.km_start >= ? AND road_condition.km_end <= ?", filter.KmStart, filter.KmEnd)
	} else if filter.KmStart != 0.0 {
		query = query.Where("road_retro_reflectivity.km_start >= ?", filter.KmStart)
	} else if filter.KmEnd != 0.0 {
		query = query.Where("road_retro_reflectivity.km_end <= ?", filter.KmEnd)
	}
	query = query.Where("road_retro_reflectivity.status = ?", "A")

	// Subquery to find the latest road condition per road ID based on created_date
	subquery := r.conn.Table("road_retro_reflectivity").
		Select("road_id, MAX(created_date) AS latest_date").
		Where("status = ?", "A").
		Group("road_id,line_no")
	if filter.Year != 0 {
		subquery = subquery.Where("year = ?", filter.Year)
	}

	// Main query with subquery to get the latest records
	query = query.Joins("JOIN (?) AS latest ON road_retro_reflectivity.road_id = latest.road_id AND road_retro_reflectivity.created_date = latest.latest_date", subquery)

	// Execute the query
	err := query.Find(&retroReflectivity).Error
	if err != nil {
		return nil, err
	}

	return retroReflectivity, nil
}

func (r *Repository) GetDataListForRoadCondition(model interface{}, where string) error {
	query := r.conn
	if where != "" {
		query = query.Where(where)
	}

	return query.Find(model).Error
}

func (r *Repository) GetRoadConditionForRoadCondition(where string) ([]models.RoadConditionSurveyM2, error) {
	var model []models.RoadConditionSurveyM2
	query := r.conn

	if err := query.Debug().Select("*, st_asgeojson(the_geom) as geojson").Find(&model).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (r *Repository) GetRoadCondition100MForRoadCondition(where string) ([]models.RoadConditionSurvey100M2, error) {
	var model []models.RoadConditionSurvey100M2
	query := r.conn

	if err := query.Debug().Select("*, st_asgeojson(the_geom) as geojson").Find(&model).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (r *Repository) GetCenterFromRoadSectionIDForRoadCondition(roadSectionID int) string {
	type LineStr struct {
		Point string `gorm:"column:center_point"`
	}
	result := LineStr{}
	err := r.conn.Debug().Raw(`
	WITH road_ids AS (
		SELECT id
		FROM road
		WHERE road_section_id = ? 
	),
	line_strings AS (
		SELECT ri.the_geom
		FROM road_info ri
		INNER JOIN road r ON ri.road_id = r.id
		WHERE r.id IN (SELECT id FROM road_ids)
	),
	merged_line AS (
		SELECT ST_LineMerge(ST_Collect(the_geom)) AS merged_geom
		FROM line_strings
	)
	SELECT ST_AsText(ST_Centroid(merged_geom)) AS center_point
	FROM merged_line;
	
	`, roadSectionID).Scan(&result).Error
	if err != nil {
		logs.Error(err)
	}

	return result.Point
}
