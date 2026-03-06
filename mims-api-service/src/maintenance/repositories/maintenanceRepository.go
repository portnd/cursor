package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	servicesDB "gitlab.com/mims-api-service/services/database"
	"gitlab.com/mims-api-service/src/maintenance/handlers"
	"gitlab.com/mims-api-service/src/maintenance/usecases"
	"gorm.io/gorm"
)

type Repository struct {
	conn *gorm.DB
}

func NewMaintainRepoHandler(conn *gorm.DB) *handlers.Handler {
	servicesDB := servicesDB.NewServicesDatabase(conn)
	usecase := usecases.NewUsecase(&Repository{conn}, servicesDB)
	handler := handlers.NewHandler(usecase)
	return handler
}

func (r *Repository) GetRoadGroupData() ([]models.RoadGroup, error) {
	var result []models.RoadGroup
	err := r.conn.Order("id").Find(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *Repository) GetRoadGroupInfoByRoadID(roadID int, year string) ([]models.MaintainPreloadGetAll, error) {
	var maintenance []models.MaintainPreloadGetAll
	query := r.conn
	query = query.Select("DISTINCT maintenance.* ")
	query = query.Preload("Budget")
	query = query.Preload("BudgetMethod")

	query = query.Joins("join maintenance_road on  maintenance_road.maintenance_id = maintenance.id ")
	if year != "" {
		query = query.Where("budget_year = ?", year)
	}
	query = query.Where("maintenance.status = ?", "A")
	query = query.Where("maintenance_road.road_id = ?", roadID)
	err := query.Order("id desc").Find(&maintenance).Error
	if err != nil {
		return maintenance, err
	}
	return maintenance, nil
}

func (r *Repository) GetRoadMaintenanceYear(roadID int) ([]models.MaintainPreloadGetAll, error) {
	var maintenance []models.MaintainPreloadGetAll
	query := r.conn
	query = query.Select("DISTINCT maintenance.budget_year")
	query = query.Joins("join maintenance_road on  maintenance_road.maintenance_id = maintenance.id ")
	query = query.Where("maintenance.status = ?", "A")
	query = query.Where("maintenance_road.road_id = ?", roadID)
	query = query.Where("maintenance_road.status = ?", "A")
	err := query.Order("maintenance.budget_year desc").Find(&maintenance).Error
	if err != nil {
		return maintenance, err
	}
	return maintenance, nil
}

func (r *Repository) GetMaintenanceListCount(queryFilter string, prams requests.MaintenancePrams, isAllData, isOwnerData bool, depotCode string) (int64, error) {
	var count int64
	query := r.conn.Model(&models.Maintenance{}).
		Joins("LEFT join maintenance_road on maintenance_road.maintenance_id = maintenance.id and maintenance_road.status = 'A'").
		Joins("LEFT JOIN road_group on road_group.id = maintenance_road.road_group_id")
	if isOwnerData && !isAllData {
		query = query.Where("maintenance.ref_depot_code = ?", depotCode)
	}
	if !isOwnerData && !isAllData {
		query = query.Where("maintenance.ref_depot_code = ?", "no_data")
	}
	if queryFilter != "" {
		query = query.Where(queryFilter)
	}
	if len(prams.BudgetMethodId) != 0 {
		query = query.Where("maintenance.budget_method_id in ?", prams.BudgetMethodId)
	}
	if len(prams.RoadGroupID) != 0 {
		query = query.Where("EXISTS (SELECT 1 FROM maintenance_road mr2 WHERE mr2.maintenance_id = maintenance.id AND mr2.status = 'A' AND mr2.road_group_id IN ?)", prams.RoadGroupID)
	} else if len(prams.RoadGroupIDDashboard) != 0 {
		queryString := []string{}
		for _, v := range prams.RoadGroupIDDashboard {
			queryString = append(queryString, fmt.Sprintf("EXISTS (SELECT 1 FROM maintenance_road mr2 WHERE mr2.maintenance_id = maintenance.id AND mr2.status = 'A' AND mr2.road_group_id = %v)", v))
		}
		query = query.Where(strings.Join(queryString, " OR "))
	}
	query = query.Where("maintenance.status = ?", "A")
	err := query.Select("COUNT(DISTINCT maintenance.id)").Scan(&count).Error
	return count, err
}

func (r *Repository) GetMaintenanceList(queryFilter string, prams requests.MaintenancePrams, isAllData, isOwnerData bool, depotCode string, limit, offset int64) ([]models.MaintenanceList, error) {
	var results []models.MaintenanceList

	// Subquery: ดึง maintenance.id ที่ผ่าน filter แล้ว limit/offset เพื่อแก้ปัญหา GORM Limit กับ Preload/Group
	idSubQuery := r.conn.Model(&models.Maintenance{}).
		Select("maintenance.id").
		Joins("LEFT join maintenance_road on maintenance_road.maintenance_id = maintenance.id and maintenance_road.status = 'A'").
		Joins("LEFT JOIN road_group on road_group.id = maintenance_road.road_group_id").
		Group("maintenance.id")
	if isOwnerData && !isAllData {
		idSubQuery = idSubQuery.Where("maintenance.ref_depot_code = ?", depotCode)
	}
	if !isOwnerData && !isAllData {
		idSubQuery = idSubQuery.Where("maintenance.ref_depot_code = ?", "no_data")
	}
	if queryFilter != "" {
		idSubQuery = idSubQuery.Where(queryFilter)
	}
	if len(prams.BudgetMethodId) != 0 {
		idSubQuery = idSubQuery.Where("maintenance.budget_method_id in ?", prams.BudgetMethodId)
	}
	if len(prams.RoadGroupID) != 0 {
		idSubQuery = idSubQuery.Where("EXISTS (SELECT 1 FROM maintenance_road mr2 WHERE mr2.maintenance_id = maintenance.id AND mr2.status = 'A' AND mr2.road_group_id IN ?)", prams.RoadGroupID)
	} else if len(prams.RoadGroupIDDashboard) != 0 {
		queryString := []string{}
		for _, v := range prams.RoadGroupIDDashboard {
			queryString = append(queryString, fmt.Sprintf("EXISTS (SELECT 1 FROM maintenance_road mr2 WHERE mr2.maintenance_id = maintenance.id AND mr2.status = 'A' AND mr2.road_group_id = %v)", v))
		}
		idSubQuery = idSubQuery.Where(strings.Join(queryString, " OR "))
	}
	idSubQuery = idSubQuery.Where("maintenance.status = ?", "A").Order("maintenance.id ASC")
	if limit > 0 {
		idSubQuery = idSubQuery.Limit(int(limit))
	}
	if offset > 0 {
		idSubQuery = idSubQuery.Offset(int(offset))
	}

	query := r.conn.Model(&models.Maintenance{}).
		Where("maintenance.id IN (?)", idSubQuery).
		Select(`maintenance.*, SUM(ABS(maintenance_road.km_start - maintenance_road.km_end))/1000 as km_total,
							ARRAY_AGG(DISTINCT ltrim(road_group.number, '0')) as road_group_names, ARRAY_AGG(DISTINCT maintenance_road.road_group_id) as road_group_id`).
		Joins("LEFT join maintenance_road on maintenance_road.maintenance_id = maintenance.id and maintenance_road.status = 'A'").
		Joins("LEFT JOIN road_group on road_group.id = maintenance_road.road_group_id").
		Group("maintenance.id").
		Order("maintenance.id ASC")

	query = query.Preload("Budget")
	query = query.Preload("BudgetMethod")
	query = query.Preload("RefDivision")
	query = query.Preload("RefDistrict")
	query = query.Preload("RefDepot")

	query = query.Preload("CreatedByUser", func(db *gorm.DB) *gorm.DB {
		return db.Select("* , CONCAT(firstname, ' ', lastname) AS name")
	})
	query = query.Preload("CreatedByUser.Department")

	query = query.Preload("UpdateByUser", func(db *gorm.DB) *gorm.DB {
		return db.Select("* , CONCAT(firstname, ' ', lastname) AS name")
	})
	query = query.Preload("UpdateByUser.Department")
	query = query.Preload("Roads", func(db *gorm.DB) *gorm.DB {
		db = db.Select(`maintenance_road.*,ST_ASTEXT(ST_FORCE2D(maintenance_road.the_geom)) as the_geom, 
						CONCAT('ทางหลวงพิเศษหมายเลข ', ltrim(road_group.number, '0')) as road_group_name ,
						road_section.name_origin_th as road_sec_name_or , road_section.name_destination_th as road_sec_name_des, 
						road.road_level ,road_info.name as road_name, road_info.ref_direction_id ,ABS(maintenance_road.km_start - maintenance_road.km_end)/1000 as distance ,
						COUNT(DISTINCT road_geom.lane_no) as lane_total, ref_direction.name as ref_direction_name`).
			Joins("LEFT JOIN road_group on road_group.id = maintenance_road.road_group_id").
			Joins("LEFT JOIN road on road.id = maintenance_road.road_id and is_active = true").
			Joins("LEFT JOIN road_section on road_section.id = road.road_section_id").
			Joins("LEFT JOIN road_info on road_info.road_id = maintenance_road.road_id and road_info.status = 'A'").
			Joins("LEFT JOIN road_geom on road_geom.road_id = maintenance_road.road_id and road_geom.status = 'A'").
			Joins("LEFT JOIN ref_direction on ref_direction.id = road_info.ref_direction_id").
			Where("maintenance_road.status = ? ", "A").
			Group("maintenance_road.id , road_group.number ,road_section.name_origin_th, road_section.name_destination_th, road.road_level,road_info.name ,road_info.ref_direction_id ,ref_direction.name")

		if len(prams.RoadGroupID) != 0 {
			db = db.Where("maintenance_road.road_group_id in ?", prams.RoadGroupID)
		} else if len(prams.RoadGroupIDDashboard) != 0 {
			queryString := []string{}
			for _, v := range prams.RoadGroupIDDashboard {
				queryString = append(queryString, fmt.Sprintf(`maintenance_road.road_group_id = %v`, v))
			}
			db = db.Where(strings.Join(queryString, " or "))
		}

		return db
	})
	query = query.Preload("Roads.InterventionCriteria")

	err := query.Find(&results).Error
	if err != nil {
		return results, err
	}
	return results, nil
}

func (r *Repository) GetMaintenanceListHistory(queryFilter string) ([]models.MaintainPreloadGetAll, error) {
	var results []models.MaintainPreloadGetAll

	query := r.conn
	query = query.Preload("Budget")
	query = query.Preload("BudgetMethod")
	if queryFilter != "" {
		query = query.Where(queryFilter)
	}

	query = query.Where("is_deleted = ?", false)
	query = query.Where("is_complete = ?", true)
	query = query.Order("id ASC")
	err := query.Find(&results).Error
	if err != nil {
		return results, err
	}
	return results, nil
}

func (r *Repository) GetTotalKmByMaintenanceID(maintenanceID int) (float64, error) {
	type resultTotalKm struct {
		TotalKm float64 `gorm:"column:total_km"`
	}
	var result resultTotalKm
	err := r.conn.Model(&models.MaintenanceRoad{}).
		Select("SUM(ABS(km_start-km_end)) as total_km").
		Where("is_deleted = ? AND maintenance_id = ?", false, maintenanceID).
		Group("maintenance_id").
		Scan(&result).
		Error

	if err != nil {
		logs.Error(err)
		return 0, err
	}
	return result.TotalKm, nil
}

// func (r *Repository) GetTotalProgressByID(maintenanceID int) (interface{}, error) {
// 	// var plan models.MaintenancePlanDetail
// 	disbursementPlan := 0
// 	var maintenanceplan responses.MaintenanceProgress
// 	query := r.conn
// 	query = query.Table("maintenance_plan")
// 	query = query.Select("SUM(maintenance_plan_detail.disbursement_plan) as disbursement_plan")
// 	query = query.Joins("join maintenance_plan_detail on maintenance_plan.id = maintenance_plan_detail.maintenance_plan_id")
// 	query = query.Where("maintenance_id = ?", maintenanceID).Where("is_current = ?", true)
// 	if err := query.Scan(&maintenanceplan).Error; err != nil {
// 		disbursementPlan = 0
// 	} else {
// 		disbursementPlan = int(maintenanceplan.DisbursementPlan)
// 	}
// 	var maintenanceProgress responses.MaintenanceProgress
// 	err := r.conn.Model(&models.MaintenancePlanDetailProgress{}).
// 		Select("SUM(progress) as progress, SUM(disbursement) as disbursement").
// 		Where("maintenance_id = ?", maintenanceID).
// 		Group("maintenance_id").
// 		Scan(&maintenanceProgress).Error
// 	disbursement := 0.0
// 	if err != nil {
// 		disbursement = 0.0
// 	} else {
// 		disbursement = maintenanceProgress.Disbursement
// 	}
// 	disbursementPercen := 0.0
// 	// helpers.PrintlnJson("disbursementPlandisbursementPlandisbursementPlan", disbursementPlan)
// 	if disbursementPlan != 0 {
// 		disbursementPercen = disbursement / float64(disbursementPlan) * 100
// 	} else {
// 		disbursementPercen = 0.0
// 	}

// 	DIGIT, _ := strconv.Atoi(os.Getenv("DIGIT"))
// 	maintenanceProgress.Disbursement = helpers.RoundFloat(disbursementPercen, DIGIT)
// 	return maintenanceProgress, nil
// }

func (r *Repository) GetRoadGroupInfoByMaintenanceID(maintenanceID int) (models.RoadGroup, error) {
	type resultTotalKm struct {
		models.Maintenance
		TotalKm float64 `gorm:"column:total_km"`
	}
	var result models.RoadGroup
	query := r.conn
	err := query.Raw(`
	SELECT maintenance_road.road_group_id, road_group.*
	FROM maintenance_road 
	JOIN road_group ON maintenance_road.road_group_id = road_group.id  
	WHERE maintenance_road.maintenance_id = ?
`, maintenanceID).Find(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *Repository) GetMaintenanceListByID(IDParent int) (models.MaintainPreload, error) {
	var maintenance models.MaintainPreload
	query := r.conn
	query = query.Select(`maintenance.*, SUM(ABS(maintenance_road.km_start - maintenance_road.km_end))/1000 as km_total,
						ARRAY_AGG(DISTINCT ltrim(road_group.number, '0')) as road_group_names`).
		Joins("LEFT join maintenance_road on maintenance_road.maintenance_id = maintenance.id and maintenance_road.status = 'A'").
		Joins("LEFT JOIN road_group on road_group.id = maintenance_road.road_group_id").
		Group("maintenance.id")

	query = query.Preload("Budget")
	query = query.Preload("Attachments")
	query = query.Preload("BudgetMethod")
	query = query.Preload("RefDivision")
	query = query.Preload("RefDistrict")
	query = query.Preload("RefDepot")

	query = query.Preload("CreatedByUser", func(db *gorm.DB) *gorm.DB {
		return db.Select("* , CONCAT(firstname, ' ', lastname) AS name")
	})
	query = query.Preload("CreatedByUser.Department")

	query = query.Preload("UpdateByUser", func(db *gorm.DB) *gorm.DB {
		return db.Select("* , CONCAT(firstname, ' ', lastname) AS name")
	})
	query = query.Preload("UpdateByUser.Department")
	query = query.Preload("Roads", func(db *gorm.DB) *gorm.DB {
		db = db.Select(`maintenance_road.*,ST_ASTEXT(ST_FORCE2D(maintenance_road.the_geom)) as the_geom, 
						CONCAT('ทางหลวงพิเศษหมายเลข ', ltrim(road_group.number, '0')) as road_group_name ,
						road_section.name_origin_th as road_sec_name_or , road_section.name_destination_th as road_sec_name_des, 
						road.road_level ,road_info.name as road_name, road_info.ref_direction_id ,ABS(maintenance_road.km_start - maintenance_road.km_end)/1000 as distance , 
						COUNT(DISTINCT road_geom.lane_no) as lane_total , ref_direction.name as ref_direction_name`).
			Joins("LEFT JOIN road_group on road_group.id = maintenance_road.road_group_id").
			Joins("LEFT JOIN road on road.id = maintenance_road.road_id and is_active = true").
			Joins("LEFT JOIN road_section on road_section.id = road.road_section_id").
			Joins("LEFT JOIN road_info on road_info.road_id = maintenance_road.road_id and road_info.status = 'A'").
			Joins("LEFT JOIN road_geom on road_geom.road_id = maintenance_road.road_id and road_geom.status = 'A'").
			Joins("LEFT JOIN ref_direction on ref_direction.id = road_info.ref_direction_id").
			Where("maintenance_road.status = ? ", "A").
			Group("maintenance_road.id , road_group.number ,road_section.name_origin_th, road_section.name_destination_th, road.road_level,road_info.name ,road_info.ref_direction_id,ref_direction.name")
		return db
	})
	query = query.Preload("Roads.InterventionCriteria")
	query = query.Preload("RoadHistories", func(db *gorm.DB) *gorm.DB {
		db = db.Select(`maintenance_road_history.*,ST_ASTEXT(ST_FORCE2D(maintenance_road_history.the_geom)) as the_geom, 
						CONCAT('ทางหลวงพิเศษหมายเลข ', ltrim(road_group.number, '0')) as road_group_name ,
						road_section.name_origin_th as road_sec_name_or , road_section.name_destination_th as road_sec_name_des, 
						road.road_level ,road_info.name as road_name, road_info.ref_direction_id ,ABS(maintenance_road_history.km_start - maintenance_road_history.km_end)/1000 as distance , 
						COUNT(DISTINCT road_geom.lane_no) as lane_total , ref_direction.name as ref_direction_name`).
			Joins("LEFT JOIN road_group on road_group.id = maintenance_road_history.road_group_id").
			Joins("LEFT JOIN road on road.id = maintenance_road_history.road_id and is_active = true").
			Joins("LEFT JOIN road_section on road_section.id = road.road_section_id").
			Joins("LEFT JOIN road_info on road_info.road_id = maintenance_road_history.road_id and road_info.status = 'A'").
			Joins("LEFT JOIN road_geom on road_geom.road_id = maintenance_road_history.road_id and road_geom.status = 'A'").
			Joins("LEFT JOIN ref_direction on ref_direction.id = road_info.ref_direction_id").
			Where("maintenance_road_history.status = ? ", "A").
			Group("maintenance_road_history.id , road_group.number ,road_section.name_origin_th, road_section.name_destination_th, road.road_level,road_info.name ,road_info.ref_direction_id,ref_direction.name")
		return db
	})

	query = query.Preload("RoadHistories.InterventionCriteria")

	// query = query.Preload("SettingInterventionCriteria")
	query = query.Where("maintenance.id_parent = ?", IDParent)
	query = query.Where("maintenance.status = ?", "A")

	err := query.First(&maintenance).Error
	if err != nil {
		return maintenance, err
	}
	return maintenance, nil
}

func (r *Repository) GetMaintenanceRoadID(maintenanceID int, mRoadID int) (*models.MaintenanceRoadPreloadById, error) {
	var mRoad models.MaintenanceRoadPreloadById
	query := r.conn
	query = query.Select(`maintenance_road.*,ST_ASTEXT(ST_FORCE2D(maintenance_road.the_geom)) as the_geom, 
	CONCAT('ทางหลวงพิเศษหมายเลข ', ltrim(road_group.number, '0')) as road_group_name ,
	road_section.name_origin_th as road_sec_name_or , road_section.name_destination_th as road_sec_name_des, 
	road.road_level ,road_info.name as road_name, road_info.ref_direction_id ,ABS(maintenance_road.km_start - maintenance_road.km_end)/1000 as distance , 
	COUNT(DISTINCT road_geom.lane_no) as lane_total , ref_direction.name as ref_direction_name`).
		Joins("LEFT JOIN road_group on road_group.id = maintenance_road.road_group_id").
		Joins("LEFT JOIN road on road.id = maintenance_road.road_id and is_active = true").
		Joins("LEFT JOIN road_section on road_section.id = road.road_section_id").
		Joins("LEFT JOIN road_info on road_info.road_id = maintenance_road.road_id and road_info.status = 'A'").
		Joins("LEFT JOIN road_geom on road_geom.road_id = maintenance_road.road_id and road_geom.status = 'A'").
		Joins("LEFT JOIN ref_direction on ref_direction.id = road_info.ref_direction_id").
		Where("maintenance_road.status = ? ", "A").
		Group("maintenance_road.id , road_group.number ,road_section.name_origin_th, road_section.name_destination_th, road.road_level,road_info.name ,road_info.ref_direction_id,ref_direction.name")
	query = query.Preload("InterventionCriteria")
	query = query.Preload("Maintenance", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", "A")
	})
	err := query.Where("maintenance_road.id = ? and maintenance_road.status = ? and maintenance_road.maintenance_id = ?", mRoadID, "A", maintenanceID).First(&mRoad).Error
	if err != nil {
		return nil, err
	}
	return &mRoad, nil
}

func (r *Repository) GetMaintenanceListByIDWithNotFilterIsComplete(maintenanceID int) (models.MaintainPreload, error) {
	var maintenance models.MaintainPreload
	query := r.conn
	query = query.Preload("Budget")
	query = query.Preload("BudgetMethod")
	query = query.Preload("MaintenanceRoads", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, maintenance_id, road_group_id, road_id, lane, km_start, km_end, maintenance_method_id, ref_surface_id, ref_surface_params_id, intervention_criteria_id, intervention_criteria_id_params, ST_ASTEXT(the_geom) as the_geom")
	})
	query = query.Preload("MaintenanceRoads.InterventionCriteria")
	query = query.Preload("MaintenanceRoads.RoadGroup")
	query = query.Preload("MaintenanceRoads.RoadInfo")
	query = query.Preload("MaintenanceRoads.RoadInfo.Direction")
	query = query.Preload("UserDepartment")
	// query = query.Preload("UserDepartment.Department")

	// query = query.Preload("SettingInterventionCriteria")
	query = query.Where("id = ?", maintenanceID)
	query = query.Where("is_deleted = ?", false)
	err := query.First(&maintenance).Error
	if err != nil {
		return maintenance, err
	}
	return maintenance, nil
}

func (r *Repository) GetMaintenanceHistoryByID(maintenanceID int, mRoadHisId int) (models.MaintenanceRoadHistoryPreloadByID, error) {
	var maintenanceRoadHis models.MaintenanceRoadHistoryPreloadByID
	query := r.conn
	query = query.Select(`maintenance_road_history.*,ST_ASTEXT(ST_FORCE2D(maintenance_road_history.the_geom)) as the_geom, 
							CONCAT('ทางหลวงพิเศษหมายเลข ', ltrim(road_group.number, '0')) as road_group_name ,
							road_section.name_origin_th as road_sec_name_or , road_section.name_destination_th as road_sec_name_des, 
							road.road_level ,road_info.name as road_name, road_info.ref_direction_id ,ABS(maintenance_road_history.km_start - maintenance_road_history.km_end)/1000 as distance , 
							COUNT(DISTINCT road_geom.lane_no) as lane_total , ref_direction.name as ref_direction_name`).
		Joins("LEFT JOIN road_group on road_group.id = maintenance_road_history.road_group_id").
		Joins("LEFT JOIN road on road.id = maintenance_road_history.road_id and is_active = true").
		Joins("LEFT JOIN road_section on road_section.id = road.road_section_id").
		Joins("LEFT JOIN road_info on road_info.road_id = maintenance_road_history.road_id and road_info.status = 'A'").
		Joins("LEFT JOIN road_geom on road_geom.road_id = maintenance_road_history.road_id and road_geom.status = 'A'").
		Joins("LEFT JOIN ref_direction on ref_direction.id = road_info.ref_direction_id").
		Where("maintenance_road_history.status = ? ", "A").
		Group("maintenance_road_history.id , road_group.number ,road_section.name_origin_th, road_section.name_destination_th, road.road_level,road_info.name ,road_info.ref_direction_id,ref_direction.name")

	query = query.Preload("InterventionCriteria")
	query = query.Preload("Maintenance")
	query = query.Preload("Attachments")
	err := query.Where("maintenance_road_history.id = ? and maintenance_road_history.maintenance_id = ?", mRoadHisId, maintenanceID).First(&maintenanceRoadHis).Error
	if err != nil {
		return maintenanceRoadHis, err
	}
	return maintenanceRoadHis, nil
}

func (r *Repository) InsertMaintenance(data models.Maintenance, conn *gorm.DB) (int, error) {
	err := conn.Table("maintenance").Create(&data).Error
	if err != nil {
		conn.Rollback()
		return 0, err
	}
	return data.ID, nil
}

func (r *Repository) UpdateMaintenance(maintenanceID int, data models.Maintenance, conn *gorm.DB) error {
	err := conn.Table("maintenance").Where("id = ?", maintenanceID).Updates(&data).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *Repository) UpdateIDParentMaintenance(maintenanceID int, idParent int, conn *gorm.DB) error {
	err := conn.Table("maintenance").Where("id = ?", maintenanceID).Updates(&models.Maintenance{IDParent: &idParent}).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *Repository) GetMaxRevisionMaintenanceByIDParent(IDParent int) (*int, *int, error) {
	type Result struct {
		Revision int
		ID       int `gorm:"column:id"`
	}

	result := Result{}

	err := r.conn.Table("maintenance").Select("revision , id").Where("id_parent = ? and status = ?", IDParent, "A").First(&result).Error
	if err != nil {
		return nil, nil, err
	}

	return &result.Revision, &result.ID, nil
}

func (r *Repository) UpdateStatusMaintenance(maintenanceID int, conn *gorm.DB) error {
	err := conn.Table("maintenance").Where("id = ? and status = ?", maintenanceID, "A").Updates(&models.Maintenance{Status: "I"}).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *Repository) GetMaintenanceRoadByMaintenanceID(maintenanceID int) ([]models.MaintenanceRoad, error) {
	var result []models.MaintenanceRoad
	err := r.conn.Where("maintenance_id = ? and status = ?", maintenanceID, "A").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) GetMaintenanceRoadHistoryByMaintenanceID(maintenanceID int) ([]models.MaintenanceRoadHistory, error) {
	var result []models.MaintenanceRoadHistory
	err := r.conn.Where("maintenance_id = ? and status = ?", maintenanceID, "A").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) InsertMaintenanceRoad(data []models.MaintenanceRoad, conn *gorm.DB) error {
	err := conn.Table("maintenance_road").Create(&data).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *Repository) InsertMaintenanceRoadHistory(data []models.MaintenanceRoadHistory, conn *gorm.DB) error {
	helpers.PrintlnJson(data)
	err := conn.Table("maintenance_road_history").Create(&data).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *Repository) UpdateStatusMaintenanceRoad(ids []int, conn *gorm.DB) error {
	err := conn.Table("maintenance_road").Where("id in (?) and status = ?", ids, "A").Updates(&models.MaintenanceRoad{Status: "I"}).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *Repository) UpdateStatusMaintenanceRoadHistory(ids []int, conn *gorm.DB) error {
	err := conn.Table("maintenance_road_history").Where("id in (?) and status = ?", ids, "A").Updates(&models.MaintenanceRoadHistory{Status: "I"}).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *Repository) GetInterventionCriteriaByID(interventionCriteriaID int) (models.InterventionCriteria, error) {
	var interventionCriteria models.InterventionCriteria
	query := r.conn
	err := query.Where("id = ?", interventionCriteriaID).First(&interventionCriteria).Error
	if err != nil {
		return interventionCriteria, err
	}
	return interventionCriteria, nil
}

func (r *Repository) GetInterventionCriteriaParamsLatest() (models.SettingInterventionCriteriaParams, error) {
	var interventionCriteria models.SettingInterventionCriteriaParams
	query := r.conn
	err := query.Where("is_latest = ?", true).First(&interventionCriteria).Error
	if err != nil {
		return interventionCriteria, err
	}
	return interventionCriteria, nil
}

func (r *Repository) CreateMaintenanceRoad(mRoad models.MaintenanceRoad, tx *gorm.DB) (*int, error) {
	var id int

	sql := fmt.Sprintf(`
        INSERT INTO maintenance_road (
            maintenance_id,
            id_parent,
            revision,
            status,
            road_id,
			road_group_id,
			maintenance_method_id,
            intervention_criteria_id,
            intervention_criteria_id_params,
            ref_surface_id,
            ref_surface_params_id,
            km_start,
            km_end,
            the_geom,
            maintenance_type,
            lane_no,
            grid_no,
            created_by,
            updated_by,
            created_at,
            updated_at
        ) VALUES (
            %d, %d, %d, '%s', %d, %d,%d, %d, %d, %d, %d, %f, %f, %s, %d, %d, %d, %d, %d, '%s', '%s'
        ) RETURNING id`,
		mRoad.MaintenanceID,
		mRoad.IDParent,
		mRoad.Revision,
		mRoad.Status,
		mRoad.RoadID,
		mRoad.RoadGroupId,
		mRoad.MaintenanceMethodID,
		mRoad.InterventionCriteriaID,
		mRoad.InterventionCriteriaIDParams,
		mRoad.RefSurfaceID,
		mRoad.RefSurfaceParamsID,
		mRoad.KmStart,
		mRoad.KmEnd,
		mRoad.TheGeom,
		mRoad.MaintenanceType,
		mRoad.LaneNo,
		mRoad.GridNo,
		mRoad.CreatedBy,
		mRoad.UpdatedBy,
		mRoad.CreatedAt.Format("2006-01-02 15:04:05"),
		mRoad.UpdatedAt.Format("2006-01-02 15:04:05"),
	)

	// Execute the SQL statement and scan the result into id
	err := tx.Raw(sql).Row().Scan(&id)
	if err != nil {
		return nil, err
	}

	// Return the id
	return &id, nil
}

func (r *Repository) DeleteMaintenanceRoad(maintenanceID int, mRoadId int) error {
	query := r.StartTransSection()
	if err := query.Debug().Where("id = ?", mRoadId).Updates(models.MaintenanceRoad{Status: "D"}).Error; err != nil {
		return err
	}

	if err := query.Commit().Error; err != nil {
		query.Rollback()
		return err
	}
	return nil
}

func (r *Repository) UpdateIDParentMaintenanceRoad(maintenanceID int, idParent int, conn *gorm.DB) error {
	err := conn.Table("maintenance_road").Where("id = ?", maintenanceID).Updates(&models.MaintenanceRoad{IDParent: idParent}).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *Repository) GetMaxRevisionMaintenanceRoad(maintenanceRoadID int) (*int, *int, error) {
	type Result struct {
		Revision int
		IDparent int `gorm:"column:id_parent"`
	}

	result := Result{}

	err := r.conn.Table("maintenance_road").Select("revision , id_parent").Where("id = ? and status = ?", maintenanceRoadID, "A").Order("revision  DESC").Find(&result).Error
	if err != nil {
		return nil, nil, err
	}

	return &result.Revision, &result.IDparent, nil
}

func (r *Repository) CreateMaintenanceRoadHistory(data models.MaintenanceRoadHistory, tx *gorm.DB) (*int, error) {
	var id int

	sql := fmt.Sprintf(`
        INSERT INTO maintenance_road_history (
            maintenance_id,
            id_parent,
            revision,
            status,
            road_id,
			road_group_id,
			maintenance_method_id,
            intervention_criteria_id,
            intervention_criteria_id_params,
            ref_surface_id,
            ref_surface_params_id,
            km_start,
            km_end,
            the_geom,
            maintenance_type,
            lane_no,
            grid_no,
            created_by,
            updated_by,
            created_at,
            updated_at
        ) VALUES (
            %d, %d, %d, '%s', %d, %d, %d,%d, %d, %d, %d, %f, %f, %s, %d, %d, %d, %d, %d, '%s', '%s'
        ) RETURNING id`,
		data.MaintenanceID,
		data.IDParent,
		data.Revision,
		data.Status,
		data.RoadID,
		data.RoadGroupId,
		data.MaintenanceMethodID,
		data.InterventionCriteriaID,
		data.InterventionCriteriaIDParams,
		data.RefSurfaceID,
		data.RefSurfaceParamsID,
		data.KmStart,
		data.KmEnd,
		data.TheGeom,
		data.MaintenanceType,
		data.LaneNo,
		data.GridNo,
		data.CreatedBy,
		data.UpdatedBy,
		data.CreatedAt.Format("2006-01-02 15:04:05"),
		data.UpdatedAt.Format("2006-01-02 15:04:05"),
	)

	// Execute the SQL statement and scan the result into id
	err := tx.Raw(sql).Row().Scan(&id)
	if err != nil {
		return nil, err
	}

	// Return the id
	return &id, nil
}

func (r *Repository) UpdateMaintenanceRoadHistory(data models.MaintenanceRoadHistory, attReqs []requests.MaintenanceAttachmentsReq, historyID int, ids []int) error {
	query := r.StartTransSection()

	// sql := fmt.Sprintf(`INSERT INTO maintenance_road_history (maintenance_id, road_group_id, road_id, lane, intervention_criteria_id, km_start, km_end, the_geom, is_deleted, created_by, updated_by, maintenance_method_id, intervention_criteria_id_params, ref_surface_id, ref_surface_params_id) VALUES (%d, %d, %d, %d, %d, %v, %v, %v, %v, %d, %d, %d, %d, %d, %d)`, data.MaintenanceID, data.RoadGroupID, data.RoadID, data.Lane, data.InterventionCriteriaID, data.KmStart, data.KmEnd, data.TheGeom, data.IsDeleted, data.CreatedBy, data.UpdatedBy, data.MaintenanceMethodID, data.InterventionCriteriaIDParams, data.RefSurfaceID, data.RefSurfaceParamsID)
	// var maintenanceRoadHistory models.MaintenanceRoadHistory
	// if err := query.Exec(sql).Find(&maintenanceRoadHistory).Error; err != nil {
	// 	query.Rollback()
	// 	return err
	// }
	// err := query.Where("id IN (?)", ids).Delete(models.MaintenanceRoadHistory{}).Error
	// if err != nil {
	// 	query.Rollback()
	// 	return err
	// }

	// query = query.Where(s).Updates(&data)
	var maintenanceRoadHistory models.MaintenanceRoadHistory
	sql := fmt.Sprintf(`
    UPDATE maintenance_road 
    SET 
        maintenance_id = %d,
        id_parent = %d,
        revision = %d,
        status = '%s',
        road_id = %d,
		maintenance_method_id = %d,
        intervention_criteria_id = %d,
        intervention_criteria_id_params = %d,
        ref_surface_id = %d,
        ref_surface_params_id = %d,
        km_start = %f,
        km_end = %f,
        the_geom = '%s',
        maintenance_type = %d,
        lane_no = %d,
        grid_no = %d,
        created_by = %d,
        updated_by = %d,
        created_at = '%s',
        updated_at = '%s'
    WHERE id = %d`,
		data.MaintenanceID,
		data.IDParent,
		data.Revision,
		data.Status,
		data.RoadID,
		data.MaintenanceMethodID,
		data.InterventionCriteriaID,
		data.InterventionCriteriaIDParams,
		data.RefSurfaceID,
		data.RefSurfaceParamsID,
		data.KmStart,
		data.KmEnd,
		data.TheGeom,
		data.MaintenanceType,
		data.LaneNo,
		data.GridNo,
		data.CreatedBy,
		data.UpdatedBy,
		data.CreatedAt.Format("2006-01-02 15:04:05"),
		data.UpdatedAt.Format("2006-01-02 15:04:05"),
		data.ID, // ID of the row to update
	)
	if err := query.Exec(sql).Find(&maintenanceRoadHistory).Error; err != nil {
		query.Rollback()
		return err
	}

	for _, item := range attReqs {
		var maintenanceRoadAttachHistory models.MaintenanceRoadHistoryAttachment
		copier.Copy(&maintenanceRoadAttachHistory, &item)
		maintenanceRoadAttachHistory.MaintenanceID = data.MaintenanceID
		maintenanceRoadAttachHistory.MaintenanceRoadHistoryID = maintenanceRoadHistory.ID
		maintenanceRoadAttachHistory.CreatedBy = data.CreatedBy
		maintenanceRoadAttachHistory.UpdatedBy = data.UpdatedBy
		maintenanceRoadAttachHistory.CreatedAt = time.Now().UTC()
		maintenanceRoadAttachHistory.UpdatedAt = time.Now().UTC()
		if err := query.Save(&maintenanceRoadAttachHistory).Error; err != nil {
			query.Rollback()
			return err
		}
	}

	err := query.Where("id IN (?)", ids).Delete(models.MaintenanceRoadHistoryAttachment{}).Error
	if err != nil {
		query.Rollback()
		return err
	}

	query.Commit()
	return nil
}

func (r *Repository) UpdateIDParentMaintenanceRoadHistory(maintenanceRoadHisID int, idParent int, conn *gorm.DB) error {
	err := conn.Table("maintenance_road_history").Where("id = ?", maintenanceRoadHisID).Updates(&models.MaintenanceRoadHistory{IDParent: idParent}).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *Repository) GetMaxRevisionMaintenanceRoadHistory(maintenanceRoadHisID int) (*int, *int, error) {
	type Result struct {
		Revision int
		IDparent int `gorm:"column:id_parent"`
	}

	result := Result{}

	err := r.conn.Table("maintenance_road_history").Select("revision , id_parent").Where("id = ? and status = ?", maintenanceRoadHisID, "A").Order("revision  DESC").Find(&result).Error
	if err != nil {
		return nil, nil, err
	}

	return &result.Revision, &result.IDparent, nil
}

func (r *Repository) UpdateMaintenanceRoad(datas []models.MaintenanceRoad, conn *gorm.DB) error {
	for _, data := range datas {
		//sql := fmt.Sprintf("UPDATE maintenance_road  SET maintenance_id = %d, road_group_id = %d, road_id = %d, lane = %d, maintenance_method_id = %d, ref_surface_params_id = %d, km_start = %v, km_end = %v, the_geom = %v, is_deleted = %t, created_by = %d, updated_by = %d , intervention_criteria_id = %d   WHERE id = %d", data.MaintenanceID, data.RoadGroupID, data.RoadID, data.Lane, data.MaintenanceMethodID, data.RefSurfaceParamsID, data.KmStart, data.KmEnd, data.TheGeom, data.IsDeleted, data.CreatedBy, data.UpdatedBy, data.InterventionCriteriaID, data.ID)
		sql := fmt.Sprintf(`
    UPDATE maintenance_road 
    SET 
        maintenance_id = %d,
        id_parent = %d,
        revision = %d,
        status = '%s',
        road_id = %d,
		maintenance_method_id = %d,
        intervention_criteria_id = %d,
        intervention_criteria_id_params = %d,
        ref_surface_id = %d,
        ref_surface_params_id = %d,
        km_start = %f,
        km_end = %f,
        the_geom = '%s',
        maintenance_type = %d,
        lane_no = %d,
        grid_no = %d,
        created_by = %d,
        updated_by = %d,
        created_at = '%s',
        updated_at = '%s'
    WHERE id = %d`,
			data.MaintenanceID,
			data.IDParent,
			data.Revision,
			data.Status,
			data.RoadID,
			data.MaintenanceMethodID,
			data.InterventionCriteriaID,
			data.InterventionCriteriaIDParams,
			data.RefSurfaceID,
			data.RefSurfaceParamsID,
			data.KmStart,
			data.KmEnd,
			data.TheGeom,
			data.MaintenanceType,
			data.LaneNo,
			data.GridNo,
			data.CreatedBy,
			data.UpdatedBy,
			data.CreatedAt.Format("2006-01-02 15:04:05"),
			data.UpdatedAt.Format("2006-01-02 15:04:05"),
			data.ID, // ID of the row to update
		)
		if err := conn.Exec(sql).Error; err != nil {
			conn.Rollback()
			return err
		}
	}
	return nil
}

// func (r *Repository) DeleteMaintenanceRoad(maintenanceID int, ids []int, conn *gorm.DB) error {
// 	err := conn.Not(ids).Where("maintenance_id = ?", maintenanceID).Delete(models.MaintenanceRoad{}).Error
// 	if err != nil {
// 		conn.Rollback()
// 		return err
// 	}
// 	return nil
// }

func (r *Repository) GetMaintenanceRoad(maintenanceID, roadID int) ([]models.MaintenanceRoadData, error) {
	var maintenanceRoad []models.MaintenanceRoadData
	query := r.conn
	query = query.Preload("RoadInfo")
	query = query.Preload("RoadInfo.Direction")
	query = query.Preload("RoadGroup")
	query = query.Preload("InterventionCriteria")
	err := query.Where("maintenance_id = ?", maintenanceID).Where("road_id = ?", roadID).Find(&maintenanceRoad).Error
	if err != nil {
		return maintenanceRoad, err
	}
	return maintenanceRoad, nil
}

func (r *Repository) GetRoadMaintenance(maintenanceID int) ([]models.MaintenanceRoadData, error) {
	var maintenanceRoad []models.MaintenanceRoadData
	query := r.conn
	query = query.Preload("RoadInfo")
	query = query.Preload("RoadInfo.Direction")
	query = query.Preload("RoadGroup")
	query = query.Preload("InterventionCriteria")
	err := query.Where("maintenance_id = ?", maintenanceID).Find(&maintenanceRoad).Error
	if err != nil {
		return maintenanceRoad, err
	}
	return maintenanceRoad, nil
}

func (r *Repository) GetRoadMaintenanceHistory(maintenanceID int) ([]models.MaintenanceRoadHistoryData, error) {
	var maintenanceRoadHistory []models.MaintenanceRoadHistoryData
	query := r.conn
	query = query.Preload("RoadInfo")
	query = query.Preload("RoadInfo.Direction")
	query = query.Preload("RoadGroup")
	query = query.Preload("InterventionCriteria")
	query = query.Preload("Attacchments")
	err := query.Where("maintenance_id = ?", maintenanceID).Find(&maintenanceRoadHistory).Error
	if err != nil {
		return maintenanceRoadHistory, err
	}
	return maintenanceRoadHistory, nil
}

func (r *Repository) GetMaintenanceRoadHistory(maintenanceID, roadID int) ([]models.MaintenanceRoadHistoryData, error) {
	var maintenanceRoadHistory []models.MaintenanceRoadHistoryData
	query := r.conn
	query = query.Preload("RoadInfo")
	query = query.Preload("RoadInfo.Direction")
	query = query.Preload("RoadGroup")
	query = query.Preload("InterventionCriteria")
	err := query.Where("maintenance_id = ?", maintenanceID).Where("road_id = ?", roadID).Find(&maintenanceRoadHistory).Error
	if err != nil {
		return maintenanceRoadHistory, err
	}
	return maintenanceRoadHistory, nil
}

func (r *Repository) DeleteMaintenance(maintenanceID int) error {
	query := r.StartTransSection()
	err := query.Where("id = ? and status = ?", maintenanceID, "A").Updates(models.Maintenance{Status: "D"}).Error
	if err != nil {
		query.Rollback()
		return err
	}

	err = query.Where("maintenance_id = ?", maintenanceID).Delete(models.MaintenanceAttachment{}).Error
	if err != nil {
		query.Rollback()
		return err
	}

	err = query.Where("maintenance_id = ? and status = ?", maintenanceID, "A").Updates(models.MaintenanceRoad{Status: "D"}).Error
	if err != nil {
		query.Rollback()
		return err
	}

	err = query.Where("maintenance_id = ? and status = ?", maintenanceID, "A").Updates(models.MaintenanceRoadHistory{Status: "D"}).Error
	if err != nil {
		query.Rollback()
		return err
	}

	err = query.Where("maintenance_id = ?", maintenanceID).Delete(models.MaintenanceRoadHistoryAttachment{}).Error
	if err != nil {
		query.Rollback()
		return err
	}

	query.Commit()

	return nil
}

func (r *Repository) FindRefDirectionIDByRoadID(id int) (int, error) {
	var result models.RoadInfo
	err := r.conn.Where("road_id = ? and status = ?", id, "A").First(&result).Error
	if err != nil {
		return 0, err
	}
	return result.RefDirectionId, nil
}

func (r *Repository) GetGeomByRoadID(roadID int) (models.RoadInfo, error) {
	var result models.RoadInfo
	err := r.conn.Select("the_geom,km_start,km_end").Where("road_id = ? AND status = 'A'", roadID).Order("revision DESC").First(&result).Error

	return result, err
}

func (r *Repository) GetGeomByRoadGeomID(roadID, lane int) (models.RoadGeom, error) {
	var roadGeom models.RoadGeom
	err := r.conn.Select("the_geom,km_start,km_end").Where("road_id = ? AND lane_no = ? AND status = 'A'", roadID, lane).Order("revision DESC").First(&roadGeom).Error

	return roadGeom, err
}

func (r *Repository) StartTransSection() *gorm.DB {
	tx := r.conn.Begin()
	return tx
}

func (r *Repository) RollBack(tx *gorm.DB) error {
	tx.Rollback()
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) Commit(tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

/////////////////

func (r *Repository) GetInterventionCriteria(ID int) ([]models.InterventionCriteria, error) {
	var interventionCriteria []models.InterventionCriteria
	query := r.conn
	if err := query.Where("maintenance_method = ?", ID).Find(&interventionCriteria).Error; err != nil {
		return interventionCriteria, err
	}
	return interventionCriteria, nil
}

func (r *Repository) GetRefCriteriaMethod() ([]models.RefCriteriaMethodData, error) {
	var refCriteriaMethod []models.RefCriteriaMethodData
	query := r.conn
	query = query.Preload("Children", func(db *gorm.DB) *gorm.DB {
		return db.Select("* ,maintenance_standard_name as label").Order("maintenance_sequence ASC")
	})
	if err := query.Select("* ,name as label").Order("id ASC").Find(&refCriteriaMethod).Error; err != nil {
		return refCriteriaMethod, err
	}
	return refCriteriaMethod, nil
}

////////////////////

func (r *Repository) GetMaintenanceBudget() (interface{}, error) {
	var budget []models.SettingBudgetProload
	query := r.conn
	query = query.Preload("BudgetMethods", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_deleted = ?", false)
	})

	if err := query.Where("is_deleted = ?", false).Find(&budget).Error; err != nil {
		return "", err
	}
	return budget, nil
}

func (r *Repository) MaintenanceFinished(maintenanceID int, req requests.MaintenanceFinished) error {
	// lastInspectionDate, err := helpers.StringToDate(req.LastInspectionDate)
	// if err != nil {
	// 	return err
	// }

	guaranteeExpirationDate, err := helpers.StringToDate(req.GuaranteeExpirationDate)
	if err != nil {
		return err
	}
	query := r.conn
	//if err := query.Table("maintenance").Where("id = ?", maintenanceID).Updates(models.Maintenance{LastInspectionDate: lastInspectionDate, GuaranteeExpirationDate: guaranteeExpirationDate, IsComplete: true}).Error; err != nil {
	if err := query.Table("maintenance").Where("id = ?", maintenanceID).Updates(models.Maintenance{GuaranteeExpirationDate: guaranteeExpirationDate}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) MaintenanceCheckProgressComplete(maintenanceID int) (float64, error) {
	query := r.conn
	var sum responses.Sum
	if err := query.Table("maintenance_plan_detail_progress").Select("sum(progress) as sum").Where("maintenance_id = ?", maintenanceID).Find(&sum).Error; err != nil {
		return 0, err
	}
	return sum.Sum, nil
}

// ///////////////////////////// maintenance_history  /////////////////////////////////
func (r *Repository) GetMaintenanceHistory(maintenanceID int) ([]models.MaintenanceRoadHistoryData, error) {
	var maintenanceRoadHistory []models.MaintenanceRoadHistoryData
	query := r.conn
	query = query.Preload("InterventionCriteria")
	query = query.Preload("RoadGroup")
	query = query.Preload("RoadInfo")
	query = query.Preload("RoadInfo.Direction")
	query = query.Where("maintenance_id = ?", maintenanceID)
	query = query.Where("is_deleted = ?", false)
	err := query.Find(&maintenanceRoadHistory).Error
	if err != nil {
		return maintenanceRoadHistory, err
	}
	return maintenanceRoadHistory, nil
}

func (r *Repository) CheckMaintenanceDuplicate(idParent int, name string) (bool, error) {
	var count int64
	query := r.conn
	query = query.Model(&models.Maintenance{})
	if idParent != 0 {
		query = query.Where("id_parent <> ?", idParent)
	}

	query = query.Where("name = ? AND status = ?", name, "A")

	query = query.Count(&count)

	if query.Error != nil {
		return true, query.Error
	}
	return count > 0, nil
}

func (r *Repository) GetRoadMaintenanceByID(maintenanceID int) (models.Maintenance, error) {
	var maintenance models.Maintenance
	query := r.conn
	query = query.Where("id = ? ", maintenanceID).Find(&maintenance)
	if query.Error != nil {
		return maintenance, query.Error
	}
	return maintenance, nil
}

func (r *Repository) DeleteMaintenanceHistory(maintenanceID int, hisID int) error {
	query := r.StartTransSection()
	if err := query.Where("maintenance_id = ?", maintenanceID).Where("id = ?", hisID).Updates(models.MaintenanceRoadHistory{Status: "D"}).Error; err != nil {
		return err
	}

	query = query.Where("maintenance_id = ?", maintenanceID).Where("maintenance_road_history_id = ?", hisID).Delete(models.MaintenanceRoadHistoryAttachment{})

	if err := query.Commit().Error; err != nil {
		query.Rollback()
		return err
	}
	return nil
}

func (r *Repository) GetMaintenanceRoadTheGeomByID(maintenanceID int) (models.MaintenanceRoad, error) {
	var maintenance models.MaintenanceRoad
	query := r.conn
	query = query.Select("ST_ASTEXT(the_geom) as the_geom")
	query = query.Where("maintenance_id = ?", maintenanceID)

	err := query.First(&maintenance).Error
	if err != nil {
		return maintenance, err
	}
	return maintenance, nil
}

func (r *Repository) GetMaintenanceYearHistory(maintenanceID int) ([]models.MaintainPreloadGetAll, error) {
	var maintenance []models.MaintainPreloadGetAll
	query := r.conn
	query = query.Select("DISTINCT maintenance.budget_year")
	query = query.Where("is_deleted = ?", false)
	query = query.Where("is_complete = ?", true)
	query = query.Where("id = ?", maintenanceID)
	err := query.Order("maintenance.budget_year desc").Find(&maintenance).Error
	if err != nil {
		return maintenance, err
	}
	return maintenance, nil
}

func (r *Repository) StAstext(theGeom string) (string, error) {
	var res responses.TheGeom
	query := r.conn
	if theGeom == "" {
		return "", nil
	}
	if err := query.Raw("select ST_AsText($1) as the_geom", theGeom).Scan(&res).Error; err != nil {
		return "", err
	}
	return res.TheGeom, nil
}

func (r *Repository) GetMaintenanceYear(roadId int) ([]models.Maintenance, error) {
	var maintenance []models.Maintenance
	query := r.conn
	query = query.Select("budget_year").Joins("LEFT JOIN maintenance_road on maintenance_road.maintenance_id = maintenance.id and maintenance_road.status = 'A'").
		Where("maintenance.status = ?", "A").Group("maintenance.budget_year")
	if roadId != 0 {
		query = query.Where("maintenance_road.road_id = ?", roadId).Order("maintenance.budget_year DESC")
	}
	if err := query.Find(&maintenance).Error; err != nil {
		return maintenance, err
	}
	return maintenance, nil
}

func (r *Repository) InsertMaintenanceaAttachment(data []models.MaintenanceAttachment, conn *gorm.DB) error {
	err := conn.Table("maintenance_attachment").Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteMaintenanceAttachment(ids []int, conn *gorm.DB) error {
	err := conn.Table("maintenance_attachment").Where("id IN (?)", ids).Delete(&models.MaintenanceAttachment{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateMaintenanceAttachment(oldMaintenanceID int, newMaintenanceID int, conn *gorm.DB) error {
	err := conn.Table("maintenance_attachment").Where("maintenance_id = ?", oldMaintenanceID).Updates(&models.MaintenanceAttachment{MaintenanceID: newMaintenanceID}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetMaintenanceaRoadHisAttachmentByMaintenanceID(maintenanceID int) ([]models.MaintenanceRoadHistoryAttachment, error) {
	var result []models.MaintenanceRoadHistoryAttachment
	err := r.conn.Where("maintenance_id = ?", maintenanceID).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) GetMaintenanceaRoadHisAttachmentByMRoadHisID(mRoadHisID int) ([]models.MaintenanceRoadHistoryAttachment, error) {
	var result []models.MaintenanceRoadHistoryAttachment
	err := r.conn.Where("maintenance_road_history_id = ?", mRoadHisID).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

// func (r *Repository) InsertMaintenanceaRoadHisAttachment(data []models.MaintenanceRoadHistoryAttachment, conn *gorm.DB) error {
// 	err := conn.Table("maintenance_road_history_attachment").Create(&data).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *Repository) DeleteMaintenanceRoadHisAttachment(ids []int, conn *gorm.DB) error {
	err := conn.Table("maintenance_road_history_attachment").Where("id IN (?)", ids).Delete(&models.MaintenanceRoadHistoryAttachment{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateMaintenanceRoadHisAttachment(oldMaintenanceRoadHisID int, newMaintenanceRoadHisID int, conn *gorm.DB) error {
	err := conn.Table("maintenance_road_history_attachment").Where("maintenance_road_history_id = ?", oldMaintenanceRoadHisID).Updates(&models.MaintenanceRoadHistoryAttachment{MaintenanceRoadHistoryID: newMaintenanceRoadHisID}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetMaintenanceaAttachmentByMaintenanceID(maintenanceID int) ([]models.MaintenanceAttachment, error) {
	var result []models.MaintenanceAttachment
	err := r.conn.Model(&models.MaintenanceAttachment{}).Where("maintenance_id = ?", maintenanceID).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) GetMaintenanceaAttachmentByID(id int) (*models.MaintenanceAttachment, error) {
	var result models.MaintenanceAttachment
	err := r.conn.Model(&models.MaintenanceAttachment{}).Where("id = ?", id).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (t *Repository) GetLastRoadInfoByID(roadID int) (*models.RoadInfoGeomData, error) {
	var roadInfo models.RoadInfoGeomData
	query := t.conn
	query = query.Select("* , ST_ASTEXT(the_geom) as line_string")
	query = query.Table("road_info")
	query = query.Preload("RoadGeom", func(db *gorm.DB) *gorm.DB {
		return db.Select("*, ST_ASTEXT(the_geom) as line_string").Where("status = ?", "A")
	})
	if err := query.Where("road_id = ?", roadID).Where("status = ?", "A").Order("road_id asc").Find(&roadInfo).Error; err != nil {
		return nil, err
	}

	return &roadInfo, nil
}

func (t *Repository) GetRoadGroupByRoadId(roadID int) (*int, error) {
	var roadGropId int
	query := t.conn
	query = query.Table("road").Select("road_section.road_group_id").
		Joins("LEFT JOIN road_section on road_section.id = road.road_section_id")

	if err := query.Where("road.id = ? and road.is_active = ?", roadID, true).Find(&roadGropId).Error; err != nil {
		return nil, err
	}

	return &roadGropId, nil
}

func (t *Repository) GetMaintenanceByRoadID(roadID int, year int) ([]models.MaintainPreload, error) {
	var maintenances []models.MaintainPreload
	query := t.conn

	query = query.Select(`maintenance.*, SUM(ABS(maintenance_road.km_start - maintenance_road.km_end))/1000 as km_total,
						ARRAY_AGG(DISTINCT ltrim(road_group.number, '0')) as road_group_names`).
		Joins("LEFT join maintenance_road on maintenance_road.maintenance_id = maintenance.id and maintenance_road.status = 'A'").
		Joins("LEFT JOIN road_group on road_group.id = maintenance_road.road_group_id").
		Group("maintenance.id")

	query = query.Preload("Budget")
	query = query.Preload("Attachments")
	query = query.Preload("BudgetMethod")
	query = query.Preload("RefDivision")
	query = query.Preload("RefDistrict")
	query = query.Preload("RefDepot")

	query = query.Preload("CreatedByUser", func(db *gorm.DB) *gorm.DB {
		return db.Select("* , CONCAT(firstname, ' ', lastname) AS name")
	})
	query = query.Preload("CreatedByUser.Department")

	query = query.Preload("UpdateByUser", func(db *gorm.DB) *gorm.DB {
		return db.Select("* , CONCAT(firstname, ' ', lastname) AS name")
	})

	query = query.Preload("UpdateByUser.Department")
	query = query.Preload("Roads", func(db *gorm.DB) *gorm.DB {
		db = db.Select(`maintenance_road.*,ST_ASTEXT(ST_FORCE2D(maintenance_road.the_geom)) as the_geom,
		ST_AsGeoJSON(ST_ASTEXT(ST_FORCE2D(maintenance_road.the_geom))) as the_geom_string,
						CONCAT('ทางหลวงพิเศษหมายเลข ', ltrim(road_group.number, '0')) as road_group_name ,
						road_section.name_origin_th as road_sec_name_or , road_section.name_destination_th as road_sec_name_des, 
						road.road_level ,road_info.name as road_name, road_info.ref_direction_id ,ABS(maintenance_road.km_start - maintenance_road.km_end)/1000 as distance , 
						COUNT(DISTINCT road_geom.lane_no) as lane_total , ref_direction.name as ref_direction_name`).
			Joins("LEFT JOIN road_group on road_group.id = maintenance_road.road_group_id").
			Joins("LEFT JOIN road on road.id = maintenance_road.road_id and is_active = true").
			Joins("LEFT JOIN road_section on road_section.id = road.road_section_id").
			Joins("LEFT JOIN road_info on road_info.road_id = maintenance_road.road_id and road_info.status = 'A'").
			Joins("LEFT JOIN road_geom on road_geom.road_id = maintenance_road.road_id and road_geom.status = 'A'").
			Joins("LEFT JOIN ref_direction on ref_direction.id = road_info.ref_direction_id").
			Where("maintenance_road.status = ? ", "A").
			Group("maintenance_road.id , road_group.number ,road_section.name_origin_th, road_section.name_destination_th, road.road_level,road_info.name ,road_info.ref_direction_id,ref_direction.name")
		return db
	})
	query = query.Preload("Roads.InterventionCriteria")
	query = query.Preload("RoadHistories", func(db *gorm.DB) *gorm.DB {
		db = db.Select(`maintenance_road_history.*,ST_AsGeoJSON(ST_ASTEXT(ST_FORCE2D(maintenance_road_history.the_geom))) as the_geom, 
						CONCAT('ทางหลวงพิเศษหมายเลข ', ltrim(road_group.number, '0')) as road_group_name ,
						road_section.name_origin_th as road_sec_name_or , road_section.name_destination_th as road_sec_name_des, 
						road.road_level ,road_info.name as road_name, road_info.ref_direction_id ,ABS(maintenance_road_history.km_start - maintenance_road_history.km_end)/1000 as distance , 
						COUNT(DISTINCT road_geom.lane_no) as lane_total , ref_direction.name as ref_direction_name`).
			Joins("LEFT JOIN road_group on road_group.id = maintenance_road_history.road_group_id").
			Joins("LEFT JOIN road on road.id = maintenance_road_history.road_id and is_active = true").
			Joins("LEFT JOIN road_section on road_section.id = road.road_section_id").
			Joins("LEFT JOIN road_info on road_info.road_id = maintenance_road_history.road_id and road_info.status = 'A'").
			Joins("LEFT JOIN road_geom on road_geom.road_id = maintenance_road_history.road_id and road_geom.status = 'A'").
			Joins("LEFT JOIN ref_direction on ref_direction.id = road_info.ref_direction_id").
			Where("maintenance_road_history.status = ? ", "A").
			Group("maintenance_road_history.id , road_group.number ,road_section.name_origin_th, road_section.name_destination_th, road.road_level,road_info.name ,road_info.ref_direction_id,ref_direction.name")
		return db
	})

	query = query.Preload("RoadHistories.InterventionCriteria")

	// query = query.Preload("SettingInterventionCriteria")
	query = query.Where("maintenance_road.road_id = ?", roadID)
	query = query.Where("maintenance.status = ?", "A")

	if year != 0 {
		query = query.Where("maintenance.budget_year = ?", year)
	}

	err := query.Find(&maintenances).Error
	if err != nil {
		return nil, err
	}
	return maintenances, nil
}

func (r *Repository) GetMaintenanceByIDParent(IDParent int) (models.Maintenance, error) {
	var maintenance models.Maintenance
	query := r.conn
	err := query.Where("id_parent = ? AND status = ?", IDParent, "A").First(&maintenance).Error
	if err != nil {
		return maintenance, err
	}
	return maintenance, nil
}

func (r *Repository) GetGeomJsonFromMaintenanceRoadID(maintenanceRoadID int) ([]byte, error) {
	type Result struct {
		TheGeom []byte `gorm:"column:the_geom"`
	}

	result := Result{}

	err := r.conn.Table("maintenance_road").Select("ST_AsGeoJSON(the_geom) AS the_geom").Where("id = ? and status = ?", maintenanceRoadID, "A").First(&result).Error
	if err != nil {
		return nil, err
	}

	return result.TheGeom, nil
}

// GetGeomJsonByMaintenanceRoadIDs returns GeoJSON for multiple maintenance_road IDs in one query (avoids N+1).
func (r *Repository) GetGeomJsonByMaintenanceRoadIDs(ids []int) (map[int][]byte, error) {
	out := make(map[int][]byte)
	if len(ids) == 0 {
		return out, nil
	}
	type Row struct {
		ID     int    `gorm:"column:id"`
		TheGeom []byte `gorm:"column:the_geom"`
	}
	var rows []Row
	err := r.conn.Table("maintenance_road").
		Select("id, ST_AsGeoJSON(ST_Force2D(the_geom)) AS the_geom").
		Where("id IN ? AND status = ?", ids, "A").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.ID] = row.TheGeom
	}
	return out, nil
}

func (r *Repository) GetGeomJsonFromMaintenanceRoadHistoryID(maintenanceRoadHisID int) ([]byte, error) {
	type Result struct {
		TheGeom []byte `gorm:"column:the_geom"`
	}

	result := Result{}

	err := r.conn.Table("maintenance_road_history").Select("ST_AsGeoJSON(the_geom) AS the_geom").Where("id = ? and status = ?", maintenanceRoadHisID, "A").First(&result).Error
	if err != nil {
		return nil, err
	}

	return result.TheGeom, nil
}

// GetGeomJsonByMaintenanceRoadHistoryIDs returns GeoJSON for multiple maintenance_road_history IDs in one query (avoids N+1).
func (r *Repository) GetGeomJsonByMaintenanceRoadHistoryIDs(ids []int) (map[int][]byte, error) {
	out := make(map[int][]byte)
	if len(ids) == 0 {
		return out, nil
	}
	type Row struct {
		ID      int    `gorm:"column:id"`
		TheGeom []byte `gorm:"column:the_geom"`
	}
	var rows []Row
	err := r.conn.Table("maintenance_road_history").
		Select("id, ST_AsGeoJSON(ST_Force2D(the_geom)) AS the_geom").
		Where("id IN ? AND status = ?", ids, "A").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		out[row.ID] = row.TheGeom
	}
	return out, nil
}

func (r *Repository) GetDivisionList(isAllData, isOwnerData bool, depotCode string) ([]models.RefDivisionList, error) {
	var division []models.RefDivisionList

	query := r.conn
	query = query.Select("ref_division.id, ref_division.division_code, ref_division.name , concat('ref_division_code-', division_code) as owner_code_key").
		Joins("LEFT JOIN road_section on road_section.ref_division_code = ref_division.division_code").
		Group("ref_division.id, ref_division.division_code, ref_division.name").
		Order("ref_division.id ASC ")

	query = query.Preload("Districts", func(db *gorm.DB) *gorm.DB {
		return db.Select("*, concat('ref_district_code-', district_code) as owner_code_key")
	})
	query = query.Preload("Districts.Depots", func(db *gorm.DB) *gorm.DB {
		if isOwnerData && !isAllData {
			db = db.Where("depot_code = ?", depotCode)
		}
		if !isOwnerData && !isAllData {
			db = db.Where("depot_code = ?", "no_data")
		}
		db = db.Select("*, concat('ref_depot_code-',depot_code) as owner_code_key")
		return db
	})
	if err := query.Order("ref_division.id ASC").Find(&division).Error; err != nil {
		return nil, err
	}
	return division, nil
}

func (r *Repository) GetRoadDropdownList(isAllData, isOwnerData bool, depotCode string) ([]models.RoadListInit, error) {
	var roadInit []models.RoadListInit
	query := r.conn
	query = query.Debug().Preload("RoadSection", func(db *gorm.DB) *gorm.DB {
		if isOwnerData && !isAllData {
			db = db.Where("ref_depot_code = ?", depotCode)
		}
		if !isOwnerData && !isAllData {
			db = db.Where("ref_depot_code = ?", "no_data")
		}
		db = db.Order("id ASC")
		return db
	})
	// query = query.Preload("RoadSection.Roads", func(db *gorm.DB) *gorm.DB {
	// 	return db.Select("road.*, road_info.name as name, road_info.km_start as km_start, road_info.km_end as km_end, road_info.ref_direction_id as ref_direction_id, road_info.ref_road_type_id as ref_road_type_id, COUNT(DISTINCT road_geom.lane_no) as lane_total").
	// 		Joins("LEFT JOIN road_info on road_info.road_id = road.id and road_info.status = 'A'").
	// 		Joins("LEFT JOIN road_geom on road_geom.road_id = road.id and road_geom.status = 'A'").
	// 		Where("road.is_active = true").Order("road_level ASC").Order("id ASC")
	// })
	query = query.Preload("RoadSection.Roads", func(db *gorm.DB) *gorm.DB {
		return db.
			Select("road.*, road_info.name, road_info.km_start, road_info.km_end, road_info.ref_direction_id, road_info.ref_road_type_id, COUNT(DISTINCT road_geom.lane_no) as lane_total").
			Joins("LEFT JOIN road_info on road_info.road_id = road.id AND road_info.status = 'A'").
			Joins("LEFT JOIN road_geom on road_geom.road_id = road.id AND road_geom.status = 'A'").
			Where("road.is_active = true").
			Group("road.id,road_info.name, road_info.km_start, road_info.km_end, road_info.ref_direction_id, road_info.ref_road_type_id").
			Order("road_level ASC, id ASC")
	})

	err := query.Select("*, CONCAT('ทางหลวงพิเศษหมายเลข ', ltrim(number, '0')) as road_number").Find(&roadInit).Error
	if err != nil {
		return nil, err
	}
	return roadInit, nil
}

func (r *Repository) GetRoadDropdownListDashboard(isAllData, isOwnerData bool, depotCode string, queryRelated *string) ([]models.RoadListInit, error) {
	var roadInit []models.RoadListInit
	query := r.conn
	query = query.Debug().Preload("RoadSection", func(db *gorm.DB) *gorm.DB {
		if queryRelated != nil {
			db = db.Where(*queryRelated)
		}
		// if isOwnerData && !isAllData {
		// 	db = db.Where("ref_depot_code = ?", depotCode)
		// }
		// if !isOwnerData && !isAllData {
		// 	db = db.Where("ref_depot_code = ?", "no_data")
		// }
		db = db.Order("id ASC")
		return db
	})
	// query = query.Preload("RoadSection.Roads", func(db *gorm.DB) *gorm.DB {
	// 	return db.Select("road.*, road_info.name as name, road_info.km_start as km_start, road_info.km_end as km_end, road_info.ref_direction_id as ref_direction_id, road_info.ref_road_type_id as ref_road_type_id, COUNT(DISTINCT road_geom.lane_no) as lane_total").
	// 		Joins("LEFT JOIN road_info on road_info.road_id = road.id and road_info.status = 'A'").
	// 		Joins("LEFT JOIN road_geom on road_geom.road_id = road.id and road_geom.status = 'A'").
	// 		Where("road.is_active = true").Order("road_level ASC").Order("id ASC")
	// })
	query = query.Preload("RoadSection.Roads", func(db *gorm.DB) *gorm.DB {
		return db.
			Select("road.*, road_info.name, road_info.km_start, road_info.km_end, road_info.ref_direction_id, road_info.ref_road_type_id, COUNT(DISTINCT road_geom.lane_no) as lane_total").
			Joins("LEFT JOIN road_info on road_info.road_id = road.id AND road_info.status = 'A'").
			Joins("LEFT JOIN road_geom on road_geom.road_id = road.id AND road_geom.status = 'A'").
			Where("road.is_active = true").
			Group("road.id,road_info.name, road_info.km_start, road_info.km_end, road_info.ref_direction_id, road_info.ref_road_type_id").
			Order("road_level ASC, id ASC")
	})

	err := query.Select("*, CONCAT('ทางหลวงหมายเลข ', ltrim(number, '0')) as road_number").Find(&roadInit).Error
	if err != nil {
		return nil, err
	}
	return roadInit, nil
}

// func (r *Repository) GetRoadDivisionFilter() ([]models.RoadDivisionFilter, error) {
// 	var division []models.RoadDivisionFilter

// 	query := r.conn

// 	// query = query.Select("ref_division.id, ref_division.division_code, ref_division.name , concat('ref_division_code-', division_code) as owner_code_key").
// 	// 	Joins("RIGHT JOIN road_section on road_section.ref_division_code = ref_division.division_code").
// 	// 	Group("ref_division.id, ref_division.division_code, ref_division.name").
// 	// 	Order("ref_division.id ASC ")

// 	query = query.Preload("RoadSections", func(db *gorm.DB) *gorm.DB {
// 		return db.Order("id ASC")
// 	})

// 	query = query.Preload("RoadSections.Districts", func(db *gorm.DB) *gorm.DB {
// 		return db.Select("*, concat('ref_district_code-', district_code) as owner_code_key")
// 	})

//		query = query.Preload("RoadSections.Districts.Depots", func(db *gorm.DB) *gorm.DB {
//			return db.Select("*, concat('ref_depot_code-',depot_code) as owner_code_key")
//		})
//		if err := query.Select("*, CONCAT('ทางหลวงพิเศษหมายเลข ', ltrim(number, '0')) as road_number").Find(&division).Error; err != nil {
//			return nil, err
//		}
//		return division, nil
//	}
func (r *Repository) GetSettingBudgetMethodByID(ID int) (models.SettingBudgetMethod, error) {
	var settingBudgetMethod models.SettingBudgetMethod
	query := r.conn
	err := query.Select("is_show_method").Where("id = ? AND is_deleted = ?", ID, false).First(&settingBudgetMethod).Error
	if err != nil {
		return settingBudgetMethod, err
	}
	return settingBudgetMethod, nil
}
