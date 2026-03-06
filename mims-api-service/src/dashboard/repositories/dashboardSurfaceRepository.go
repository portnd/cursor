package repositories

import (
	"fmt"
	"time"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gorm.io/gorm"
)

func (r *repository) GetRoadByIDs(roadIDs []int) ([]models.RoadById, error) {
	var roadList []models.RoadById
	query := r.conn

	//RoadInfo
	query = query.Preload("RoadInfo", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,road_id, year, name, km_start, km_end,   ref_direction_id, st_astext(the_geom) as the_geom, revision, status, st_astext(old_the_geom) as old_the_geom, old_km_start, old_km_end, ramp_id,road_color_code, created_by,created_date,updated_by,updated_date").Where("road_id in (?)", roadIDs)
	})

	query = query.Preload("Direction")
	query = query.Preload("RoadType")

	query = query.Preload("RoadGeom", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,road_id, lane_no, km_start,km_end,ST_ASTEXT(ST_FORCE2D(the_geom)) as the_geom,revision,status,old_km_start,old_km_end,old_the_geom,remark,created_by,created_date,updated_by,updated_date").Where("status = ?", "A").Where("road_id in (?)", roadIDs).Order("revision DESC").Order("revision")
	})

	if err := query.Where("id in (?)", roadIDs).Order("id").Find(&roadList).Error; err != nil {
		return roadList, err
	}
	return roadList, nil
}

func (r *repository) GetRoadSurfaceByGroupRoadID(roadID int, laneNo int) ([]int, error) {
	var roadSurface []models.RoadSurface
	if err := r.conn.Select("surface_group").Where("status = ? or status = ?", "A", "I").Where("road_id = ?", roadID).Order("surface_group DESC").Group("surface_group").Find(&roadSurface).Error; err != nil {
		return []int{}, err
	}
	surfaceCnt := []int{}
	for _, item := range roadSurface {
		surfaceCnt = append(surfaceCnt, item.SurfaceGroup)
	}
	return surfaceCnt, nil
}

func (r *repository) GetRoadSurfaceByRoadID(roadID, laneNo, grp int) ([]models.RoadSurfacePrepareData, error) {

	var roadSurfacePrepareData []models.RoadSurfacePrepareData
	query := r.conn
	query = query.Preload("RefStructureSurface")
	query = query.Preload("RefSurfaceShoulderLeft")
	query = query.Preload("RoadSurfaceLane", func(db *gorm.DB) *gorm.DB {
		return db.Where("lane_no = ?", laneNo)
	})
	query = query.Preload("RoadSurfaceLane.RefSurface")
	query = query.Preload("RefMaterialBase")
	query = query.Preload("RefMaterialSubbase")
	query = query.Preload("RefMaterialSubgrade")
	if err := query.Where("status = ? or status = ?", "A", "I").Where("road_id = ?", roadID).Where("surface_group = ?", grp).Find(&roadSurfacePrepareData).Error; err != nil {
		return []models.RoadSurfacePrepareData{}, nil
	}

	if len(roadSurfacePrepareData) > 0 {
		return roadSurfacePrepareData, nil
	}

	return []models.RoadSurfacePrepareData{}, nil
}

func (r *repository) GetMaintenanceHistoryByRoadID(roadID, year int) (models.MaintainPreloadGetAll, error) {
	// var responds []responses.MaintenancHistoryeList
	result, err := r.GetRoadGroupInfoByRoadID(roadID, year)
	if err != nil {
		logs.Error(err)
		return models.MaintainPreloadGetAll{}, err
	}
	return result, nil
}

func (r *repository) GetRoadGroupInfoByRoadID(roadID int, year int) (models.MaintainPreloadGetAll, error) {
	var maintenance models.MaintainPreloadGetAll
	query := r.conn
	query = query.Select("DISTINCT maintenance.* ")
	query = query.Preload("Budget")
	query = query.Preload("BudgetMethod")

	query = query.Joins("join maintenance_road on  maintenance_road.maintenance_id = maintenance.id ")
	query = query.Where("maintenance.budget_year > ?", year)
	query = query.Where("maintenance.status > ?", "A")
	query = query.Where("maintenance_road.road_id = ?", roadID)
	err := query.Order("id asc").First(&maintenance).Error
	if err != nil {
		return maintenance, err
	}
	return maintenance, nil
}

func (r *repository) GetRefSurfaceParam(ID int) (models.RefSurfaceParam, error) {
	var refSurfaceParam models.RefSurfaceParam
	query := r.conn
	if err := query.Where("id = ?", ID).First(&refSurfaceParam).Error; err != nil {
		return refSurfaceParam, err
	}
	return refSurfaceParam, nil
}

func (r *repository) GetRoadSurfaceFirstGrpByRoadID(roadID int, laneNo int) ([]models.RoadSurfacePrepareData, error) {
	var roadSurface models.RoadSurface
	if err := r.conn.Select("max(surface_group) as surface_group").Where("status = ? or status = ?", "A", "I").Where("road_id = ?", roadID).Find(&roadSurface).Error; err != nil {
		return []models.RoadSurfacePrepareData{}, err
	}
	maxSurfaceGrp := roadSurface.SurfaceGroup

	var roadSurfacePrepareData []models.RoadSurfacePrepareData
	query := r.conn
	query = query.Preload("RefStructureSurface")
	query = query.Preload("RefSurfaceShoulderLeft")
	query = query.Preload("RoadSurfaceLane", func(db *gorm.DB) *gorm.DB {
		return db.Where("lane_no = ?", laneNo)
	})
	// query = query.Preload("RoadSurfaceLane.RefSurface")
	query = query.Preload("RefMaterialBase")
	query = query.Preload("RefMaterialSubbase")
	query = query.Preload("RefMaterialSubgrade")
	if err := query.Where("status = ? or status = ?", "A", "I").Where("road_id = ?", roadID).Where("surface_group = ?", maxSurfaceGrp).Find(&roadSurfacePrepareData).Error; err != nil {
		return []models.RoadSurfacePrepareData{}, nil
	}

	if len(roadSurfacePrepareData) > 0 {
		return roadSurfacePrepareData, nil
	}

	return []models.RoadSurfacePrepareData{}, nil

}

func (r *repository) GetMaintenanceData(roadId, laneNo int) ([]models.MaintenanceData2, error) {
	var maintenanceData []models.MaintenanceData2
	query := r.conn
	query = query.Preload("Budget")
	query = query.Preload("BudgetMethod")

	query = query.Preload("MaintenanceRoads", func(db *gorm.DB) *gorm.DB {
		return db.Where("lane_no = ?", laneNo).Where("road_id = ?", roadId).Where("intervention_criteria_id != ?", 0)
	})
	query = query.Preload("MaintenanceRoads.RefSurface")
	query = query.Preload("MaintenanceRoads.RoadGroup")
	query = query.Preload("MaintenanceRoads.RoadInfo")
	query = query.Preload("MaintenanceRoads.RoadInfo.Direction")
	if err := query.Where("status = ?", "A").Find(&maintenanceData).Error; err != nil {
		return maintenanceData, err
	}
	return maintenanceData, nil
}

func (r *repository) GetRefCriteriaMethodByID(ID int) (models.RefCriteriaMethod, error) {
	var refCriteriaMethod models.RefCriteriaMethod
	query := r.conn
	err := query.Where("id = ?", ID).First(&refCriteriaMethod).Error
	if err != nil {
		return refCriteriaMethod, err
	}
	return refCriteriaMethod, nil
}

func (r *repository) GetInterventionCriteriaParamsById(ID int) (models.SettingInterventionCriteriaParams, error) {
	var interventionCriteria models.SettingInterventionCriteriaParams
	query := r.conn
	err := query.Where("id = ?", ID).First(&interventionCriteria).Error
	if err != nil {
		return interventionCriteria, err
	}
	return interventionCriteria, nil
}

func (r *repository) GetRoadInfoByID(roadID int) (models.RoadInfo, error) {
	var roadInfo models.RoadInfo
	query := r.conn
	if err := query.Where("road_id = ?", roadID).First(&roadInfo).Error; err != nil {
		return roadInfo, err
	}
	return roadInfo, nil
}

func (r *repository) GetRoadDatebegin(roadID int) ([]models.RoadDatebegin, error) {
	var roadDatebegin []models.RoadDatebegin
	query := r.conn
	if err := query.Where("road_id = ?", roadID).Find(&roadDatebegin).Error; err != nil {
		return roadDatebegin, err
	}
	return roadDatebegin, nil
}

func (r *repository) GetRoadByRoadID(roadID int) (models.Road, error) {
	var road models.Road
	query := r.conn
	if err := query.Where("id = ?", roadID).First(&road).Error; err != nil {
		return road, err
	}
	return road, nil
}

func (r *repository) GetRoadGroupByID(ID int) (models.RoadGroup, error) {
	var roadGroup models.RoadGroup
	query := r.conn
	if err := query.Where("id = ?", ID).First(&roadGroup).Error; err != nil {
		return roadGroup, err
	}
	return roadGroup, nil
}

func (r *repository) GetRoadGeomLaneByID(roadID, laneNo int) (models.RoadGeom, error) {
	var roadInfo models.RoadGeom
	query := r.conn
	if err := query.Where("road_id = ?", roadID).Where("lane_no = ?", laneNo).Where("status = ?", "A").First(&roadInfo).Error; err != nil {
		return roadInfo, err
	}
	return roadInfo, nil
}

func (r *repository) CreateDataMart(data []models.DataMart) ([]models.DataMart, error) {
	var dataMart models.DataMart
	tx := r.conn.Begin()
	if tx.Where("id > ?", 0).Delete(&dataMart).Error != nil {
		tx.Rollback()
	}
	// helpers.PrintlnJson(data)
	for _, item := range data {
		// if item.RefSurfaceID == 0 {
		// 	continue
		// }
		var lastInspectionDate string
		if item.LastInspectionDate == nil {
			lastInspectionDate = "NULL"
		} else {
			str := "'" + *item.LastInspectionDate + "'"
			lastInspectionDate = str
		}
		sql := fmt.Sprintf("INSERT INTO data_mart(road_id, road_surface_id, lane_count, year, km_start, km_end, lane_no, ref_surface_id, the_geom, age, surface_year, contract_number, last_inspection_date) VALUES(%d, %d, %d, %d, %f, %f, %d, %d, %s, %d, %d, '%s', %v)", item.RoadID, item.RoadSurfaceID, item.LaneCount, item.Year, item.KmStart, item.KmEnd, item.LaneNo, item.RefSurfaceID, item.TheGeom, item.Age, item.SurfaceYear, item.ContractNumber, lastInspectionDate)
		// fmt.Println(sql)
		err := tx.Exec(sql).Error
		if err != nil {
			logs.Error(err)
			tx.Rollback()
			// return responses.NewAppErr(400, err.Error())
		}
		// query := r.conn
		// if err := query.Create(&item).Error; err != nil {
		// 	return data, err
		// }
	}
	tx.Commit()

	return data, nil
}

func (r *repository) GetDataMartInfo(roadIDs []int, depotCodes []string, filter requests.Asset) ([]models.SurfaceInfo, error) {
	var surfaceInfo []models.SurfaceInfo
	query := r.conn.Debug()
	helpers.PrintlnJson(roadIDs)
	// if len(roadIDs) > 0 {
	// 	query = query.Where("data_mart.road_id in (?)", roadIDs)
	// }

	if filter.Year != "" {
		query = query.Where("data_mart.surface_year = ?", filter.Year)
	}

	if filter.KmStart != 0.0 && filter.KmEnd != 0.0 {
		query = query.Where("data_mart.km_start >= ? AND data_mart.km_end <= ?", filter.KmStart, filter.KmEnd)
	}

	if filter.KmStart != 0.0 {
		query = query.Where("data_mart.km_start >= ?", filter.KmStart)
	} else {
		query = query.Where("data_mart.km_start >= ?", 0)
	}

	if filter.KmEnd != 0.0 {
		query = query.Where("data_mart.km_end <= ?", filter.KmEnd)
	} else if len(depotCodes) != 0 {
		query = query.Where("road_section.ref_depot_code IN ? ", depotCodes)
	}

	if len(roadIDs) != 0 {
		query = query.Where("road.id IN (?)", roadIDs)
	}

	err := query.Table("data_mart").
		Select("road_group.name as road_group_name, road_info.name as road_name, data_mart.km_start as km_start,  data_mart.km_end as km_end,  data_mart.road_id, road.road_code, lane_no, ref_surface.id as surface_id, ref_surface.name as surface_name, ref_surface.color as color_code, st_asgeojson(ST_FORCE2D(data_mart.the_geom)) AS geom_cl, road_surface_id, data_mart.lane_count, data_mart.age, data_mart.contract_number, data_mart.year, data_mart.last_inspection_date, data_mart.surface_year, ref_surface.surface_group").
		Joins("JOIN ref_surface ON data_mart.ref_surface_id = ref_surface.id").
		Joins("JOIN road ON road.id = data_mart.road_id").
		Joins("JOIN road_info ON road_info.road_id = data_mart.road_id").
		Joins("JOIN road_section ON road.road_section_id = road_section.id").
		Joins("JOIN road_group ON road.road_group_id = road_group.id").
		Where("road_info.status = ?", "A").
		Order("data_mart.road_id, data_mart.lane_no, data_mart.km_start ASC").
		Scan(&surfaceInfo).Error
	if err != nil {
		return surfaceInfo, err
	}

	return surfaceInfo, nil
}

func (r *repository) GetRoadSurfaceAll() (map[int][]models.RoadSurfaceData2, error) {
	var surfaceInfo []models.RoadSurfaceData2
	query := r.conn
	data := make(map[int][]models.RoadSurfaceData2)
	query = query.Preload("RoadSurfaceLane")
	if err := query.Where("status != ?", "D").Order("surface_group DESC").Find(&surfaceInfo).Error; err != nil {
		return data, err
	}

	for _, item := range surfaceInfo {
		data[item.RoadId] = append(data[item.RoadId], item)
	}

	for key, surfaceInfo := range data {
		surface := data[key]
		maxSurface := findMaxSurface(surface)
		for _, item := range surfaceInfo {
			if maxSurface == item.SurfaceGroup {
				data[item.RoadId] = append(data[item.RoadId], item)
			}
		}
	}
	return data, nil
}

func findMaxSurface(surfaces []models.RoadSurfaceData2) int {
	if len(surfaces) == 0 {
		// return nil // Return nil if the slice is empty
	}
	var maxSurface int // Start with the first surface as the max
	for _, surface := range surfaces {
		if surface.SurfaceGroup > maxSurface {
			maxSurface = surface.SurfaceGroup // We found a new maximum
		}
	}
	return maxSurface // Return a pointer to the RoadSurface with the largest Size
}

func (r *repository) GetRoadIDAll() ([]int, error) {
	var roadInfo []models.RoadInfo
	query := r.conn
	if err := query.Where("status = ?", "A").Find(&roadInfo).Error; err != nil {
		return []int{}, err
	}

	roadIDs := []int{}
	for _, item := range roadInfo {
		roadIDs = append(roadIDs, item.RoadId)
	}
	return roadIDs, nil
}

func (r *repository) GetRoadGeomByRoadIDLaneNo(roadID int, laneNo int) (models.RoadGeom, error) {
	var roadGeom models.RoadGeom
	query := r.conn
	err := query.Where("road_id = ?", roadID).Where("lane_no = ?", laneNo).Where("status = ?", "A").First(&roadGeom).Error
	if err != nil {
		return roadGeom, err
	}
	return roadGeom, nil
}

func (r *repository) GetRoadConditionData(roadID, laneNo int) ([]models.RoadConditionSurveyM, time.Time, error) {
	query := r.conn.Debug()
	var roadCondition models.RoadCondition
	var roadConditionSurveyM []models.RoadConditionSurveyM
	err := query.Select("max(year) as year").Where("road_id = ?", roadID).Where("lane_no = ?", laneNo).Where("status = ?", "A").Find(&roadCondition).Error
	if err != nil {
		return []models.RoadConditionSurveyM{}, time.Time{}, err
	}

	err = query.Where("road_id = ?", roadID).Where("year = ?", roadCondition.Year).Where("lane_no = ?", laneNo).Where("status = ?", "A").First(&roadCondition).Error
	if err != nil {
		return []models.RoadConditionSurveyM{}, time.Time{}, err
	}

	// Subquery for MAX Year for Road Conditions
	maxYearSubQuery := query.Model(&models.RoadCondition{}).Select("MAX(year)").
		Where("road_id = ? AND lane_no = ? AND status = ?", roadID, laneNo, "A")

	// Subquery for getting Road Condition ID
	roadConditionIdSubQuery := query.Model(&models.RoadCondition{}).Select("id").
		Where("road_id = ? AND lane_no = ? AND status = ? AND year = (?)", roadID, laneNo, "A", maxYearSubQuery).
		Limit(1)

	// Subquery for getting Road Condition Survey ID
	roadConditionSurveyIdSubQuery := query.Model(&models.RoadConditionSurvey{}).Select("id").
		Where("road_condition_id IN (?)", roadConditionIdSubQuery)

	// Subquery for getting Road Condition Survey 100m ID
	roadConditionSurvey100mIdSubQuery := query.Model(&models.RoadConditionSurvey100M{}).
		Select("id").
		Where("road_condition_survey_id IN (?)", roadConditionSurveyIdSubQuery)

	// Main query to get required records from RoadConditionSurveyM table
	// var roadConditionSurveyM []models.RoadConditionSurveyM
	result := query.Where("road_condition_survey_100m_id IN (?)", roadConditionSurvey100mIdSubQuery).
		Find(&roadConditionSurveyM)

	// Error handling and processing the result
	if result.Error != nil {
		return []models.RoadConditionSurveyM{}, time.Time{}, err
	}
	if err != nil {
		return []models.RoadConditionSurveyM{}, time.Time{}, err
	}
	return roadConditionSurveyM, roadCondition.SurveyedDate, nil
}

func (r *repository) UpdatePercentage(percent float64, userID int) error {
	stauts := false
	if percent == 100 {
		stauts = true
	}
	query := r.conn
	sql := fmt.Sprintf("UPDATE data_mart_check SET stauts = %t, percent = %.2f, updated_by = %d, updated_at = now() - interval '7 hours' WHERE id = (select max(id) from data_mart_check)", stauts, percent, userID)
	// fmt.Println(sql)
	if err := query.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) GetDataMartCheck() (interface{}, error) {
	var dataMartCheck models.DataMartCheck
	query := r.conn
	if err := query.Last(&dataMartCheck).Error; err != nil {
		return dataMartCheck, err
	}

	var res responses.DataMartCheck

	if dataMartCheck.UpdatedBy != 0 {
		var user models.Users
		if err := query.Where("id = ?", dataMartCheck.UpdatedBy).First(&user).Error; err != nil {
			return dataMartCheck, err
		}
		res.UpdatedBy = user.Firstname + " " + user.Lastname
	} else {
		res.UpdatedBy = "System"
	}
	res.UpdatedAt = dataMartCheck.UpdatedAt
	res.Stauts = dataMartCheck.Stauts
	res.Percent = dataMartCheck.Percent

	return res, nil
}

func (r *repository) GetRoadSurfaceInfo(roadID []int) ([]models.SurfaceInfo, error) {
	var surfaceInfo []models.SurfaceInfo

	var subQueryResult []int
	var query string
	var args []interface{}

	if len(roadID) != 0 {
		query = "road_id IN (?) AND status = ?"
		args = append(append(args, roadID), "A")
	} else {
		query = "status = ?"
		args = append(args, "A")
	}

	err := r.conn.
		Table("road_surface").
		Select("id").
		Where(query, args...).
		Pluck("id", &subQueryResult).Error

	if err != nil {
		return nil, err
	}

	err = r.conn.Table("road_surface_lane").
		Select("road_surface.road_id, road.road_code, lane_no, ref_surface.id as surface_id, ref_surface.name as surface_name, ref_surface.color as color_code, km_start, km_end, ST_ASTEXT(ST_FORCE2D(road_surface_lane.the_geom)) AS geom_cl, road_surface_id").
		Joins("JOIN ref_surface ON road_surface_lane.ref_surface_id = ref_surface.id").
		Joins("JOIN road_surface ON road_surface.id = road_surface_lane.road_surface_id").
		Joins("JOIN road ON road.id = road_surface.road_id").
		Where("road_surface_id IN (?)", subQueryResult).
		Find(&surfaceInfo).Error

	// //old query
	// selectCol := "r.id as road_id,road_code,rg.lane_no as lane_no ,ref_s.id as surface_id,ref_s.name  as surface_name ,rs.km_start as km_start,rs.km_end as km_end,ST_ASTEXT(ST_FORCE2D(rg.the_geom)) AS geom_cl,rsl.road_surface_id as road_surface_id"

	// query := s.conn
	// query.LogMode(false)
	// query = query.Table("road as r").
	// 	Select(selectCol).
	// 	Joins("INNER JOIN road_info ri on r.id =ri.road_id").
	// 	Joins("INNER JOIN road_surface rs on r.id = rs.road_id").
	// 	Joins("INNER JOIN road_surface_lane rsl on rs.id  = rsl.road_surface_id").
	// 	Joins("INNER JOIN road_geom rg on r.id = rg.road_id and rsl.lane_no = rg.lane_no and (CASE WHEN rs.km_start < rs.km_end THEN (rg.km_start < rs.km_end) AND (rg.km_end > rs.km_start) ELSE (rg.km_start > rs.km_end) AND (rg.km_end < rs.km_start) END)").
	// 	Joins("INNER JOIN ref_surface ref_s on ref_s.id =rsl.ref_surface_id").
	// 	Where("r.is_active and rs.status = 'A' ")
	// if len(roadID) != 0 {
	// 	query = query.Where("r.id  in (?)", roadID)
	// }
	// query = query.Group("r.id, road_code, rg.lane_no, ref_s.id, ref_s.name, rs.km_start, rs.km_end, ST_ASTEXT(ST_FORCE2D(rg.the_geom)), rsl.road_surface_id").Order("r.id ,lane_no")

	// err := query.Scan(&surfaceInfo).Error

	if err != nil {
		return surfaceInfo, err
	}

	return surfaceInfo, nil
}

func (r *repository) GetInitialSurfaceArray() ([]models.Surface, error) {
	var surface []models.Surface
	query := r.conn
	err := query.Debug().Table("ref_surface AS ref_s").Select("id , name, color").Order("id").Scan(&surface).Error
	if err != nil {
		return surface, err
	}
	return surface, nil
}
