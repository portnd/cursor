package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	servicesDB "gitlab.com/mims-api-service/services/database"
	"gitlab.com/mims-api-service/src/maintenanceAnalysis/handlers"
	"gitlab.com/mims-api-service/src/maintenanceAnalysis/usecases"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type Repository struct {
	conn    *gorm.DB
	connMon *mongo.Client
}

func NewRepositoryHandler(conn *gorm.DB, connMon *mongo.Client) *handlers.Handler {
	servicesDB := servicesDB.NewServicesDatabase(conn)
	useCase := usecases.NewUseCase(&Repository{conn, connMon}, servicesDB)
	handler := handlers.NewHandler(useCase)
	return handler
}

func (r *Repository) GetMaintenanceAnalysisCount(filter requests.AnalysisFilter, isAllData, isOwnerData bool, depotCode string) (int64, error) {
	var count int64
	query := r.conn.Model(&models.MaintenanceAnalysis{})
	if !isOwnerData && !isAllData {
		query = query.Where("1 = 0")
	} else if isOwnerData && !isAllData {
		query = query.Where("EXISTS (SELECT 1 FROM maintenance_analysis_road mar JOIN road r ON r.id = mar.road_id AND r.is_active = true JOIN road_section rs ON rs.id = r.road_section_id WHERE mar.maintenance_analysis_id = maintenance_analysis.id AND rs.ref_depot_code = ?)", depotCode)
	}
	if filter.TypeAnalysis != nil && *filter.TypeAnalysis != 0 {
		query = query.Where("maintenance_analysis_type_id = ?", filter.TypeAnalysis)
	}
	if filter.Condition != nil && *filter.Condition != "" {
		query = query.Where("name like ?", "%"+*filter.Condition+"%")
	}
	err := query.Where("status = 'A' OR status = 'I' OR status = 'E'").Count(&count).Error
	return count, err
}

func (r *Repository) GetMaintenanceAnalysis(filter requests.AnalysisFilter, isAllData, isOwnerData bool, depotCode string, limit, offset int64) ([]models.MaintenanceAnalysisPreload, error) {
	var maintenanceAnalysisPreload []models.MaintenanceAnalysisPreload
	query := r.conn

	if !isOwnerData && !isAllData {
		query = query.Where("1 = 0")
	} else if isOwnerData && !isAllData {
		query = query.Where("EXISTS (SELECT 1 FROM maintenance_analysis_road mar JOIN road r ON r.id = mar.road_id AND r.is_active = true JOIN road_section rs ON rs.id = r.road_section_id WHERE mar.maintenance_analysis_id = maintenance_analysis.id AND rs.ref_depot_code = ?)", depotCode)
		query = query.Where("EXISTS (SELECT 1 FROM prepare_data pd WHERE pd.maintenance_analysis_id = maintenance_analysis.id)")
	}

	if filter.TypeAnalysis != nil && *filter.TypeAnalysis != 0 {
		query = query.Where("maintenance_analysis_type_id = ?", filter.TypeAnalysis)
	}

	if filter.Condition != nil && *filter.Condition != "" {
		query = query.Where("name like ?", "%"+*filter.Condition+"%")
	}
	query = query.Table("maintenance_analysis").Where("status = 'A' OR status = 'I' OR status = 'E'").Order("is_favorite desc").Order("updated_at desc")
	if limit > 0 {
		query = query.Limit(int(limit))
	}
	if offset > 0 {
		query = query.Offset(int(offset))
	}
	err := query.Find(&maintenanceAnalysisPreload).Error
	if err != nil {
		return []models.MaintenanceAnalysisPreload{}, err
	}
	return maintenanceAnalysisPreload, nil
}

func (r *Repository) GetMaintenanceAnalysisById(id int) (models.MaintenanceAnalysis, error) {
	var maintenanceAnalysis models.MaintenanceAnalysis
	err := r.conn.Where("id = ?", id).Order("id ASC").First(&maintenanceAnalysis).Error
	if err != nil {
		return models.MaintenanceAnalysis{}, err
	}
	return maintenanceAnalysis, nil
}

func (r *Repository) CreateMaintenanceAnalysis(maintenanceAnalysis models.MaintenanceAnalysis) (models.MaintenanceAnalysis, error) {
	if err := r.conn.Table("maintenance_analysis").Save(&maintenanceAnalysis).Error; err != nil {
		return models.MaintenanceAnalysis{}, err
	}
	return maintenanceAnalysis, nil
}

func (r *Repository) CreateMaintenanceAnalysisRoad(ID, userID int, roadIDs []int) error {
	query := r.conn
	tx := query.Begin()

	if err := tx.Where("maintenance_analysis_id = ?", ID).Delete(&models.MaintenanceAnalysisRoad{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, roadID := range roadIDs {
		var road models.Road
		if err := r.conn.Where("id = ?", roadID).First(&road).Error; err != nil {
			// tx.Rollback()
			// return err
			continue
		}
		var analysisRoad models.MaintenanceAnalysisRoad
		analysisRoad.MaintenanceAnalysisID = ID
		analysisRoad.RoadGroupID = road.RoadGroupId
		analysisRoad.RoadID = roadID
		analysisRoad.CreatedBy = userID
		analysisRoad.CreatedAt = time.Now()
		if err := tx.Create(&analysisRoad).Error; err != nil {
			tx.Rollback()
		}
	}
	tx.Commit()
	return nil
}

func (r *Repository) SelectdPrepareData(ID int, prepareDataID []int) error {
	tx := r.conn.Begin()
	if err := tx.Model(&models.PrepareData{}).Where("maintenance_analysis_id = ?", ID).Update("is_selected", false).Error; err != nil {
		tx.Rollback()
		fmt.Println(err)
		return err
	}

	if err := tx.Model(&models.PrepareData{}).Where("id IN (?)", prepareDataID).Where("maintenance_analysis_id = ?", ID).Update("is_selected", true).Error; err != nil {
		tx.Rollback()
		fmt.Println(err)
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) DeleteMaintenanceAnalysisPlandById(id int) error {
	err := r.conn.Model(&models.MaintenanceAnalysisPlan{}).Where("maintenance_analysis_strategic_id = ?", id).Update("is_deleted", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateMaintenanceAnalysisStrategic(maintenanceAnalysisStrategic models.MaintenanceAnalysisStrategic) (models.MaintenanceAnalysisStrategic, error) {
	if err := r.conn.Table("maintenance_analysis_strategic").Save(&maintenanceAnalysisStrategic).Error; err != nil {
		return models.MaintenanceAnalysisStrategic{}, err
	}
	return maintenanceAnalysisStrategic, nil
}

func (r *Repository) CreateMaintenanceAnalysisPlan(ID int, maintenanceAnalysisPlan []models.MaintenanceAnalysisPlan) error {
	query := r.conn
	tx := query.Begin()

	ids := []int{}
	for _, item := range maintenanceAnalysisPlan {
		if item.ID != nil {
			ids = append(ids, *item.ID)
		}
	}

	if len(ids) > 0 {
		err := query.Where("id NOT IN (?)", ids).Where("maintenance_analysis_id = ?", ID).Delete(models.MaintenanceAnalysisPlan{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, item := range maintenanceAnalysisPlan {
		if err := tx.Save(&item).Error; err != nil {
			tx.Rollback()
		}

	}

	tx.Commit()
	return nil
}

func (r *Repository) GetMaintenanceAnalysisConditionById(id int) (models.MaintenanceAnalysisStrategic, error) {
	var maintenanceAnalysisStrategic models.MaintenanceAnalysisStrategic
	err := r.conn.Where("is_deleted = ? AND id = ?", false, id).Order("id ASC").First(&maintenanceAnalysisStrategic).Error
	if err != nil {
		return models.MaintenanceAnalysisStrategic{}, err
	}
	return maintenanceAnalysisStrategic, nil
}

func (r *Repository) GetMaintenanceAnalysisPlanById(id int) (models.MaintenanceAnalysisPlan, error) {
	var maintenanceAnalysisPlan models.MaintenanceAnalysisPlan
	err := r.conn.Where("id = ?", id).Order("id ASC").First(&maintenanceAnalysisPlan).Error
	if err != nil {
		return models.MaintenanceAnalysisPlan{}, err
	}
	return maintenanceAnalysisPlan, nil
}

func (r *Repository) GetMaintenanceAnalysisPlanByAnalysisID(ID int) ([]models.MaintenanceAnalysisPlan, error) {
	var maintenanceAnalysisPlan []models.MaintenanceAnalysisPlan
	err := r.conn.Where("maintenance_analysis_id = ?", ID).Order("id ASC").Find(&maintenanceAnalysisPlan).Error
	if err != nil {
		return maintenanceAnalysisPlan, err
	}
	return maintenanceAnalysisPlan, nil
}

func (r *Repository) UpdateMaintenanceAnalysisStep2(ID int, req models.MaintenanceAnalysis) error {
	var interventionCriteriaParams models.SettingInterventionCriteriaParams
	if err := r.conn.Where("is_latest = ?", true).First(&interventionCriteriaParams).Error; err != nil {
		return err
	}
	var roadWorkEffectParams models.SettingRoadWorkEffectParams
	if err := r.conn.Where("is_latest = ?", true).First(&roadWorkEffectParams).Error; err != nil {
		return err
	}
	var roadUserCostParams models.SettingRoadUserCostParams
	if err := r.conn.Where("is_latest = ?", true).First(&roadUserCostParams).Error; err != nil {
		return err
	}
	var deteriorationParams models.SettingDeteriorationParams
	if err := r.conn.Where("is_latest = ?", true).First(&deteriorationParams).Error; err != nil {
		return err
	}
	var optimization models.SettingOptimization
	if err := r.conn.Where("is_latest = ?", true).First(&optimization).Error; err != nil {
		return err
	}
	var aadtParams models.SettingAadtParams
	if err := r.conn.Where("is_latest = ?", true).First(&aadtParams).Error; err != nil {
		return err
	}

	req.InterventionCriteriaParmasID = interventionCriteriaParams.Id
	req.RoadWorkEffectParmasID = roadWorkEffectParams.Id
	req.RoadUserCostParmasID = roadUserCostParams.Id
	req.DeterationParmasID = deteriorationParams.Id
	req.OptimizationParmasID = optimization.Id

	req.AadtParmasID = aadtParams.ID

	if err := r.conn.Table("maintenance_analysis").Where("id = ?", ID).Save(&req).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateMaintenanceAnalysisPlan(maintenanceAnalysisPlan models.MaintenanceAnalysisPlan) error {
	if err := r.conn.Table("maintenance_analysis_plan").Save(&maintenanceAnalysisPlan).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) ClearPrepareData(ID int) error {
	query := r.conn
	if err := query.Table("prepare_data").Where("maintenance_analysis_id = ?", ID).Delete(models.PrepareData{}).Error; err != nil {
		return err
	}

	if err := query.Table("maintenance_analysis_road").Where("maintenance_analysis_id = ?", ID).Delete(models.PrepareData{}).Error; err != nil {
		return err
	}

	if err := query.Table("maintenance_analysis").Where("id = ?", ID).Delete(models.PrepareData{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetPrepareDataByAnalysisID(analysisID int) (interface{}, error) {
	data, err := r.GetMaintenanceAnalysisById(analysisID)
	if err != nil {
		return "", err
	}
	query := r.conn
	aadt1 := data.Aadt1
	aadt2 := data.Aadt2
	iri1 := data.Iri1
	iri2 := data.Iri2
	ifi1 := data.Ifi1
	ifi2 := data.Ifi2
	if aadt1 != nil && aadt2 != nil {
		query = query.Where("aadt > ? AND aadt < ?", aadt1, aadt2)
	} else if aadt1 != nil {
		query = query.Where("aadt > ?", aadt1)
	} else if aadt2 != nil {
		query = query.Where("aadt < ?", aadt2)
	}

	if iri1 != nil && iri2 != nil {
		query = query.Where("iri > ? AND iri < ?", iri1, iri2)
	} else if iri1 != nil {
		query = query.Where("iri > ?", iri1)
	} else if iri2 != nil {
		query = query.Where("iri < ?", iri2)
	}

	if ifi1 != nil && ifi2 != nil {
		query = query.Where("ifi > ? AND ifi < ?", ifi1, ifi2)
	} else if ifi1 != nil {
		query = query.Where("ifi > ?", ifi1)
	} else if ifi2 != nil {
		query = query.Where("ifi < ?", ifi2)
	}

	var prepareData []models.PrepareData
	if err := query.Where("maintenance_analysis_id = ?", analysisID).Order("id asc").Find(&prepareData).Error; err != nil {
		return "", err
	}
	return prepareData, nil
}

func (r *Repository) GetPrepareDataByAnalysis(analysisID int, prepareDataID []int) ([]models.PrepareData, error) {
	query := r.conn
	var prepareData []models.PrepareData
	if err := query.Where("maintenance_analysis_id = ?", analysisID).Where("id IN (?)", prepareDataID).Find(&prepareData).Error; err != nil {
		return prepareData, err
	}
	return prepareData, nil
}

func (r *Repository) DeleteMaintenanceAnalysis(ID int) error {
	query := r.conn
	tx := query.Begin()
	if err := tx.Table("prepare_data").Where("maintenance_analysis_id = ?", ID).Delete(models.PrepareData{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("maintenance_analysis_road").Where("maintenance_analysis_id = ?", ID).Delete(models.MaintenanceAnalysisRoad{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("maintenance_analysis_plan").Where("maintenance_analysis_id = ?", ID).Delete(models.MaintenanceAnalysisPlan{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("maintenance_analysis").Where("id = ?", ID).Delete(models.MaintenanceAnalysis{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *Repository) CopyMaintenanceAnalysis(ID int) (models.MaintenanceAnalysis, error) {
	query := r.conn
	tx := query.Begin()
	var analysis models.MaintenanceAnalysis
	if err := tx.Where("id = ?", ID).First(&analysis).Error; err != nil {
		tx.Rollback()
		return analysis, err
	}
	analysis.ID = 0
	analysis.IsFavorite = false
	analysis.CreatedAt = time.Now()
	analysis.UpdatedAt = time.Now()
	analysis.Status = ""
	if err := tx.Create(&analysis).Error; err != nil {
		tx.Rollback()
		return analysis, err
	}

	var prepareData []models.PrepareDataForCopy
	if err := tx.Table("prepare_data").Where("maintenance_analysis_id = ?", ID).Find(&prepareData).Error; err != nil {
		tx.Rollback()
		return analysis, err
	}

	for _, item := range prepareData {
		item.ID = 0
		item.MaintenanceAnalysisID = analysis.ID
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return analysis, err
		}
	}

	var roads []models.MaintenanceAnalysisRoad
	if err := tx.Where("maintenance_analysis_id = ?", ID).Find(&roads).Error; err != nil {
		tx.Rollback()
		return analysis, err
	}

	for _, item := range roads {
		item.ID = 0
		item.MaintenanceAnalysisID = analysis.ID
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return analysis, err
		}
	}

	var plans []models.MaintenanceAnalysisPlan
	if err := tx.Where("maintenance_analysis_id = ?", ID).Find(&plans).Error; err != nil {
		tx.Rollback()
		return analysis, err
	}
	for _, item := range plans {
		item.ID = nil
		item.MaintenanceAnalysisID = analysis.ID
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return analysis, err
		}
	}

	tx.Commit()
	return analysis, nil
}

func (r *Repository) FavoriteMaintenanceAnalysis(ID int) error {
	var maintenanceAnalysis models.MaintenanceAnalysis
	query := r.conn
	if err := query.Where("id = ?", ID).First(&maintenanceAnalysis).Error; err != nil {
		return err
	}

	if maintenanceAnalysis.IsFavorite {
		if err := query.Model(&maintenanceAnalysis).Where("id = ?", ID).Updates(map[string]interface{}{"is_favorite": false}).Error; err != nil {
			return err
		}
	} else {
		if err := query.Model(&models.MaintenanceAnalysis{}).Where("id = ?", ID).Update("is_favorite", true).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) GetMaintenanceAnalysisRoadByID(ID int) ([]models.MaintenanceAnalysisRoads, error) {
	var maintenanceAnalysisRoad []models.MaintenanceAnalysisRoads
	query := r.conn
	if err := query.Where("maintenance_analysis_id = ?", ID).Order("id desc").Find(&maintenanceAnalysisRoad).Error; err != nil {
		return maintenanceAnalysisRoad, err
	}
	return maintenanceAnalysisRoad, nil
}

func (r *Repository) GetMaintenanceAnalysisRoadDataByID(ID int) (models.MaintenanceAnalysisPreload, error) {
	// var maintenanceAnalysisRoad []models.MaintenanceAnalysisRoadData
	// query := r.conn
	// query = query.Preload("RoadInfo")
	// if err := query.Select("").Where("maintenance_analysis_id = ?", ID).Group("road_id").Find(&maintenanceAnalysisRoad).Error; err != nil {
	// 	return maintenanceAnalysisRoad, err
	// }
	// return maintenanceAnalysisRoad, nil
	var maintenanceAnalysisPreload models.MaintenanceAnalysisPreload
	query := r.conn
	// query = query.Preload("List.RoadGroup")
	query = query.Preload("PrepareData")
	query = query.Preload("TargetData")
	query = query.Preload("ConditionData")

	// query = query.Preload("Strategic")
	// query = query.Where("is_deleted = ? AND status = ?", false, 1)

	query = query.Where("id = ?", ID).Order("is_favorite , id desc").First(&maintenanceAnalysisPreload)
	err := query.Error
	if err != nil {
		return models.MaintenanceAnalysisPreload{}, err
	}
	return maintenanceAnalysisPreload, nil
}

func (r *Repository) GetModelResultDataById(id int) ([]models.ModelResult, error) {

	query := r.connMon.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_MAINTENANCE_ANALYSIS_RESULT_COLLECTION"))
	filter := bson.D{{Key: "maintenance_analysis_id", Value: id}}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cursor, err := query.Find(ctx, filter)
	if err != nil {
		fmt.Println(err.Error())
		return []models.ModelResult{}, nil
	}

	var modelResult []models.ModelResult
	if err = cursor.All(context.TODO(), &modelResult); err != nil {
		fmt.Println(err.Error())
		return []models.ModelResult{}, nil
	}

	return modelResult, nil
}

func (r *Repository) GetLatestPrepareDataByMaintenanceAnalysisId(id int) (models.PrepareData, error) {
	var prepareData models.PrepareData
	query := r.conn
	err := query.Table("prepare_data").Where("maintenance_analysis_id = ?", id).Order("id DESC").Limit(1).Find(&prepareData).Error
	if err != nil {
		return models.PrepareData{}, err
	}
	return prepareData, nil
}

func (r *Repository) GetMaintenanceAnalysisStrategicBudgetTypeById(id int) (models.RefMaintenanceAnalysisCondition, error) {
	var maintenanceAnalysisStrategicBudgetType models.RefMaintenanceAnalysisCondition
	query := r.conn
	err := query.Where("id = ?", id).Find(&maintenanceAnalysisStrategicBudgetType).Error
	if err != nil {
		return models.RefMaintenanceAnalysisCondition{}, err
	}
	return maintenanceAnalysisStrategicBudgetType, nil
}

func (r *Repository) GetMaintenanceAnalysisStrategiTargetById(id int) (models.RefMaintenanceAnalysisTarget, error) {
	var refMaintenanceAnalysisTarget models.RefMaintenanceAnalysisTarget
	query := r.conn
	err := query.Where("id = ?", id).Find(&refMaintenanceAnalysisTarget).Error
	if err != nil {
		return models.RefMaintenanceAnalysisTarget{}, err
	}
	return refMaintenanceAnalysisTarget, nil
}

func (r *Repository) GetMaintenanceAnalysisRoadNameById(id int) ([]models.RoadInfo, error) {
	var roadName []models.RoadInfo
	query := r.conn
	query = query.Table("maintenance_analysis_road")
	query = query.Select(`road_info.name`)
	query = query.Joins("LEFT JOIN road on road.id = maintenance_analysis_road.road_id")
	query = query.Joins("LEFT JOIN road_info on road_info.road_id = maintenance_analysis_road.road_id")
	query = query.Where("maintenance_analysis_road.maintenance_analysis_id  = ?", id)
	query = query.Order("road.id ASC")

	err := query.Find(&roadName).Error
	if err != nil {
		return []models.RoadInfo{}, err
	}
	return roadName, nil

}

func (r *Repository) GetMaintenanceAnalysisRoadIdById(id int) ([]models.MaintenanceAnalysisRoad, error) {
	var maintenanceAnalysisRoad []models.MaintenanceAnalysisRoad
	query := r.conn
	err := query.Select("road_id").Where("maintenance_analysis_id = ?", id).Group("road_id").Find(&maintenanceAnalysisRoad).Error
	if err != nil {
		return []models.MaintenanceAnalysisRoad{}, err
	}
	return maintenanceAnalysisRoad, nil
}

func (r *Repository) GetRoad() ([]models.RoadInfo, error) {
	var roadInfo []models.RoadInfo
	query := r.conn
	query = query
	if err := query.Group("id ,name").Find(&roadInfo).Error; err != nil {
		return roadInfo, err
	}
	return roadInfo, nil
}

//////////////////////// report ////////////////////////////

func (r *Repository) GetInterventionCriteriaParamsByID() (models.SettingInterventionCriteriaParams, error) {
	var data models.SettingInterventionCriteriaParams
	query := r.conn
	if err := query.Where("is_latest = ?", true).First(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (r *Repository) GetRefCriteriaMethodBySurface(surface string) ([]models.RefCriteriaMethod, error) {
	var data []models.RefCriteriaMethod
	query := r.conn
	if err := query.Where("surface = ?", surface).Find(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (r *Repository) GetMaintenanceAnalysisResult(ID, roadID, plan, analysisTypeID int) ([]models.ModelResult, error) {
	query := r.connMon.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_MAINTENANCE_ANALYSIS_RESULT_COLLECTION"))
	// filter := bson.D{{"maintenance_analysis_id", 2}, {Key: "plan_id", Value: 2}}
	var filter primitive.D
	if analysisTypeID == 1 {
		filter = bson.D{
			{Key: "maintenance_analysis_id", Value: ID},
			{Key: "plan", Value: plan},
			{Key: "road_id", Value: roadID},
			{Key: "data.repair", Value: true},
			{Key: "maintenance_analysis_type_id", Value: analysisTypeID},
			// {Key: "analyst_year", Value: 2025},
			// {Key: "repair_budget_type", Value: "Limited Budget"},
		}
	} else {
		filter = bson.D{
			{Key: "maintenance_analysis_id", Value: ID},
			{Key: "road_id", Value: roadID},
			{Key: "data.repair", Value: true},
			{Key: "maintenance_analysis_type_id", Value: analysisTypeID},
			// {Key: "analyst_year", Value: 2027},
			// {Key: "repair_budget_type", Value: "Limited Budget"},
		}
	}

	cursor, err := query.Find(context.TODO(), filter)
	if err != nil {
		return []models.ModelResult{}, err
	}

	var modelResult []models.ModelResult
	if err = cursor.All(context.TODO(), &modelResult); err != nil {
		fmt.Println(err.Error())
		return []models.ModelResult{}, err
	}

	return modelResult, nil
}

// GetMaintenanceAnalysisResultByRoadIDs fetches results for multiple roads in one MongoDB query (faster for report type 5).
func (r *Repository) GetMaintenanceAnalysisResultByRoadIDs(ID int, roadIDs []int, plan, analysisTypeID int) (map[int][]models.ModelResult, error) {
	if len(roadIDs) == 0 {
		return map[int][]models.ModelResult{}, nil
	}
	query := r.connMon.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_MAINTENANCE_ANALYSIS_RESULT_COLLECTION"))
	var filter primitive.D
	if analysisTypeID == 1 {
		filter = bson.D{
			{Key: "maintenance_analysis_id", Value: ID},
			{Key: "plan", Value: plan},
			{Key: "road_id", Value: bson.D{{Key: "$in", Value: roadIDs}}},
			{Key: "data.repair", Value: true},
			{Key: "maintenance_analysis_type_id", Value: analysisTypeID},
		}
	} else {
		filter = bson.D{
			{Key: "maintenance_analysis_id", Value: ID},
			{Key: "road_id", Value: bson.D{{Key: "$in", Value: roadIDs}}},
			{Key: "data.repair", Value: true},
			{Key: "maintenance_analysis_type_id", Value: analysisTypeID},
		}
	}
	cursor, err := query.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var all []models.ModelResult
	if err = cursor.All(context.TODO(), &all); err != nil {
		return nil, err
	}
	byRoad := make(map[int][]models.ModelResult)
	for _, r := range roadIDs {
		byRoad[r] = nil
	}
	for i := range all {
		roadID := all[i].RoadID
		byRoad[roadID] = append(byRoad[roadID], all[i])
	}
	return byRoad, nil
}

func (r *Repository) GetMaintenanceAnalysisResultGroupByYear(ID, plan, analysisTypeID int) ([]int, error) {
	query := r.connMon.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_MAINTENANCE_ANALYSIS_RESULT_COLLECTION"))
	var pipeline []primitive.M
	if analysisTypeID == 1 {
		pipeline = []bson.M{
			{
				"$match": bson.M{
					"maintenance_analysis_id":      ID,
					"plan":                         plan,
					"maintenance_analysis_type_id": analysisTypeID,
					// "repair_budget_type":      "Limited Budget",
				},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"analyst_year": "$analyst_year",
					},
				},
			},
		}
	} else {
		pipeline = []bson.M{
			{
				"$match": bson.M{
					"maintenance_analysis_id":      ID,
					"maintenance_analysis_type_id": analysisTypeID,
					// "repair_budget_type":      "Limited Budget",
				},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"analyst_year": "$analyst_year",
					},
				},
			},
		}
	}

	cursor, err := query.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return []int{}, err
	}

	defer cursor.Close(context.Background())
	var years []int
	for cursor.Next(context.Background()) {
		var result bson.M
		var result2 responses.BsonM
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)

		err = mapstructure.Decode(result, &result2)
		if err != nil {
			log.Fatal(err)
		}
		years = append(years, int(result2.ID["analyst_year"].(int32)))

	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	var modelResult []models.ModelResultReport
	if err = cursor.All(context.TODO(), &modelResult); err != nil {
		fmt.Println(err.Error())
		return []int{}, err
	}

	return years, nil
}

func (r *Repository) GetUserByID(userId uint) (models.Users, error) {
	user := models.Users{}
	err := r.conn.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return models.Users{}, err
	}
	return user, nil
}

func (r *Repository) GetMaintenanceHistoryByRoadID(roadID, year int) (models.MaintainPreloadGetAll, error) {
	// var responds []responses.MaintenancHistoryeList
	result, err := r.GetRoadGroupInfoByRoadID(roadID, year)
	if err != nil {
		logs.Error(err)
		return models.MaintainPreloadGetAll{}, err
	}

	// for _, item := range result {
	// 	maintenanceRoad, err := r.GetRoadMaintenance(item.ID)
	// 	if err != nil {
	// 		if errors.Is(err, gorm.ErrRecordNotFound) {
	// 			continue
	// 		}
	// 		logs.Error(err)
	// 		return responds, err
	// 	}
	// 	var respond responses.MaintenancHistoryeList
	// 	copier.Copy(&respond, &item)
	// 	var roads []models.MaintenanceRoadData
	// 	for _, item2 := range maintenanceRoad {
	// 		var road models.MaintenanceRoadData
	// 		copier.Copy(&road, &item2)
	// 		roads = append(roads, road)
	// 	}
	// 	respond.MaintenanceRoads = roads
	// 	////////////////////////////////////////////////////
	// 	maintenanceRoadHis, err := r.GetRoadMaintenanceHistory(item.ID)
	// 	if err != nil {
	// 		if errors.Is(err, gorm.ErrRecordNotFound) {
	// 			continue
	// 		}
	// 		logs.Error(err)
	// 		return responds, err
	// 	}
	// 	var roadHis []models.MaintenanceRoadHistoryData
	// 	for _, item3 := range maintenanceRoadHis {
	// 		var road models.MaintenanceRoadHistoryData
	// 		copier.Copy(&road, &item3)

	// 		roadHis = append(roadHis, road)
	// 	}
	// 	if len(roadHis) <= 0 {
	// 		respond.MaintenanceRoadHistories = []string{}
	// 	} else {
	// 		respond.MaintenanceRoadHistories = roadHis
	// 	}

	// 	responds = append(responds, respond)
	// }

	return result, nil
}

func (r *Repository) GetRoadGroupInfoByRoadID(roadID int, year int) (models.MaintainPreloadGetAll, error) {
	var maintenance models.MaintainPreloadGetAll
	query := r.conn
	query = query.Select("DISTINCT maintenance.* ")
	query = query.Preload("Budget")
	query = query.Preload("BudgetMethod")

	query = query.Joins("join maintenance_road on  maintenance_road.maintenance_id = maintenance.id ")
	query = query.Where("maintenance.budget_year > ?", year)
	query = query.Where("maintenance.status = 'A'")
	query = query.Where("maintenance_road.road_id = ?", roadID)
	err := query.Order("id asc").First(&maintenance).Error
	if err != nil {
		return maintenance, err
	}
	return maintenance, nil
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

func (r *Repository) GetReport4(id, plan int) ([]models.ModelResult, error) {
	query := r.connMon.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_MAINTENANCE_ANALYSIS_RESULT_COLLECTION"))
	filter := bson.D{{Key: "maintenance_analysis_id", Value: id}, {Key: "plan", Value: plan}}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "road_id", Value: 1}, {Key: "analyst_year", Value: 1}, {Key: "data.prepare_data.road_geom.lane_no", Value: 1}})

	cursor, err := query.Find(context.TODO(), filter, findOptions)
	if err != nil {
		fmt.Println(err.Error())
		return []models.ModelResult{}, nil
	}

	var modelResult []models.ModelResult
	if err = cursor.All(context.TODO(), &modelResult); err != nil {
		fmt.Println(err.Error())
		return []models.ModelResult{}, nil
	}
	return modelResult, nil
}

func (r *Repository) GetRoadAll() ([]models.Road, error) {
	var roads []models.Road
	query := r.conn
	err := query.Find(&roads).Error
	if err != nil {
		return roads, err
	}
	return roads, nil
}

func (r *Repository) GetMaintenanceAnalysisPreviousID(ID int) (int, error) {
	var maintenanceAnalysis models.MaintenanceAnalysis
	query := r.conn
	if err := query.Where("previous_id = ?", ID).First(&maintenanceAnalysis).Error; err != nil {
		return 0, nil
	}
	return maintenanceAnalysis.ID, nil
}

func (r *Repository) GetAnalysisModel(id int, req requests.MaintenanceAnalysisModel) ([]models.ModelResult, error) {

	query := r.connMon.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_MAINTENANCE_ANALYSIS_RESULT_COLLECTION"))
	// filterConditions := bson.D{{"maintenance_analysis_id", id}}
	// Initialize filter with a mandatory condition
	filter := bson.D{{Key: "maintenance_analysis_id", Value: id}}
	if req.InterventionCriteriaID != nil {
		filter = append(filter, bson.E{Key: "data.intervention_criteria_id", Value: *req.InterventionCriteriaID})
	}

	// Create an array to hold $and conditions
	var andConditions []bson.E

	// // Check and add conditions for 'aadt'
	// if req.Aadt1 != nil && req.Aadt2 != nil {
	// 	andConditions = append(andConditions,
	// 		bson.E{Key: "aadt", Value: bson.D{{Key: "$gt", Value: req.Aadt1}}},
	// 		bson.E{Key: "aadt", Value: bson.D{{Key: "$lt", Value: req.Aadt2}}},
	// 	)
	// } else if req.Aadt1 != nil {
	// 	andConditions = append(andConditions, bson.E{Key: "aadt", Value: bson.D{{Key: "$gt", Value: req.Aadt1}}})
	// } else if req.Aadt2 != nil {
	// 	andConditions = append(andConditions, bson.E{Key: "aadt", Value: bson.D{{Key: "$lt", Value: req.Aadt2}}})
	// }

	// // Check and add conditions for 'iri'
	// if req.Iri1 != nil && req.Iri2 != nil {
	// 	andConditions = append(andConditions,
	// 		bson.E{Key: "data.rwe_result.iri", Value: bson.D{{Key: "$gt", Value: req.Iri1}}},
	// 		// bson.E{Key: "data.rwe_result.iri", Value: bson.D{{Key: "$lt", Value: req.Iri2}}},
	// 	)
	// } else if req.Iri1 != nil {
	// 	andConditions = append(andConditions, bson.E{Key: "data.rwe_result.iri", Value: bson.D{{Key: "$gt", Value: req.Iri1}}})
	// } else if req.Iri2 != nil {
	// 	andConditions = append(andConditions, bson.E{Key: "data.rwe_result.iri", Value: bson.D{{Key: "$lt", Value: req.Iri2}}})
	// }

	// // Check and add conditions for 'gn'
	// if req.Gn1 != nil && req.Gn2 != nil {
	// 	andConditions = append(andConditions,
	// 		bson.E{Key: "gn", Value: bson.D{{Key: "$gt", Value: req.Gn1}}},
	// 		bson.E{Key: "gn", Value: bson.D{{Key: "$lt", Value: req.Gn2}}},
	// 	)
	// } else if req.Gn1 != nil {
	// 	andConditions = append(andConditions, bson.E{Key: "gn", Value: bson.D{{Key: "$gt", Value: req.Gn1}}})
	// } else if req.Gn2 != nil {
	// 	andConditions = append(andConditions, bson.E{Key: "gn", Value: bson.D{{Key: "$lt", Value: req.Gn2}}})
	// }

	// Merge andConditions with the initial filter only if there are additional conditions
	if len(andConditions) > 0 {
		filter = append(filter, bson.E{Key: "$and", Value: andConditions})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "road_id", Value: 1}, {Key: "analyst_year", Value: 1}})

	cursor, err := query.Find(ctx, filter, findOptions)
	// cursor, err := query.Find(ctx, filter)
	if err != nil {
		fmt.Println(err.Error())
		return []models.ModelResult{}, nil
	}

	var modelResult []models.ModelResult
	if err = cursor.All(context.TODO(), &modelResult); err != nil {
		fmt.Println(err.Error())
		return []models.ModelResult{}, nil
	}

	return modelResult, nil
}

func (r *Repository) GetRoadInfo() (map[int]models.RoadInfo, error) {
	var roadInfos []models.RoadInfo
	query := r.conn
	if err := query.Where("status = ?", "A").Find(&roadInfos).Error; err != nil {
		return map[int]models.RoadInfo{}, nil
	}
	roadInfoData := make(map[int]models.RoadInfo)
	for _, item := range roadInfos {
		roadInfoData[item.RoadId] = item
	}
	return roadInfoData, nil
}

func (r *Repository) GetRoadGroup() (map[int]models.RoadGroup, error) {
	var roadGrp []models.RoadGroup
	query := r.conn
	if err := query.Find(&roadGrp).Error; err != nil {
		return map[int]models.RoadGroup{}, nil
	}
	roadData := make(map[int]models.RoadGroup)
	for _, item := range roadGrp {
		roadData[item.Id] = item
	}
	return roadData, nil
}

func (r *Repository) GetRoads() (map[int]models.Road, error) {
	var roads []models.Road
	query := r.conn
	if err := query.Find(&roads).Error; err != nil {
		return map[int]models.Road{}, nil
	}
	roadData := make(map[int]models.Road)
	for _, item := range roads {
		roadData[item.Id] = item
	}
	return roadData, nil
}

func (r *Repository) GetRefCriteriaMethod() ([]models.RefCriteriaMethod, error) {
	var refCriteriaMethod []models.RefCriteriaMethod
	query := r.conn
	if err := query.Find(&refCriteriaMethod).Error; err != nil {
		return refCriteriaMethod, err
	}
	return refCriteriaMethod, nil
}

func (r *Repository) GetInterventionCriteriaByID(ID int) ([]models.InterventionCriteria, error) {
	var interventionCriteria []models.InterventionCriteria
	query := r.conn
	if err := query.Where("maintenance_method = ?", ID).Find(&interventionCriteria).Error; err != nil {
		return interventionCriteria, err
	}
	return interventionCriteria, nil
}

func (r *Repository) UpdateAnalysisModel(ID int, dataID string, interventionCriteriaID int) (*mongo.SingleResult, error) {
	// query := r.connMon.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_MAINTENANCE_ANALYSIS_RESULT_COLLECTION"))
	_id, _ := primitive.ObjectIDFromHex(dataID)

	// 1) Create the context
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()

	// 2) Create the connection
	database := r.connMon.Database(os.Getenv("MONGODB_DB"))
	// 3) Select the database
	collection := database.Collection(os.Getenv("MONGODB_MAINTENANCE_ANALYSIS_RESULT_COLLECTION"))

	// 4) Select the collection
	update := bson.M{
		"$set": bson.M{"data.intervention_criteria_change_id": interventionCriteriaID},
	}

	// 5) Create the search filter
	filter := bson.M{"_id": _id}
	// 7) Create an instance of an options and set the desired options
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := collection.FindOneAndUpdate(ctx, filter, update, &opt)
	if result.Err() != nil {
		return result, result.Err()
	}
	return result, nil
}

func (r *Repository) GetCriteriaMethodColor() (map[string]string, error) {
	var data []models.RefCriteriaMethod
	query := r.conn
	if err := query.Find(&data).Error; err != nil {
		return map[string]string{}, err
	}
	color := make(map[string]string)
	for _, item := range data {
		color[item.Name] = item.Color
	}
	return color, nil
}

func (r *Repository) GetRefOwner() ([]models.RefOwner, error) {
	var data []models.RefOwner
	query := r.conn
	if err := query.Where("is_active = ?", true).Find(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (r *Repository) DashboardMap(ID, analysisTypeID, year int, plan *int) ([]models.ModelResultReportMap, error) {
	query := r.connMon.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_MAINTENANCE_ANALYSIS_RESULT_COLLECTION"))
	var pipeline []primitive.M
	if analysisTypeID == 1 {
		if plan != nil {
			pipeline = []bson.M{
				{
					"$match": bson.M{
						"maintenance_analysis_id":      ID,
						"plan":                         plan,
						"maintenance_analysis_type_id": analysisTypeID,
						"year":                         year,
					},
				},
			}
		} else {
			pipeline = []bson.M{
				{
					"$match": bson.M{
						"maintenance_analysis_id":      ID,
						"maintenance_analysis_type_id": analysisTypeID,
						"year":                         year,
					},
				},
			}
		}

	} else {
		pipeline = []bson.M{
			{
				"$match": bson.M{
					"maintenance_analysis_id":      ID,
					"maintenance_analysis_type_id": analysisTypeID,
					"year":                         year,
				},
			},
		}
	}

	cursor, err := query.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return []models.ModelResultReportMap{}, err
	}

	defer cursor.Close(context.Background())
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	var modelResult []models.ModelResultReportMap
	if err = cursor.All(context.TODO(), &modelResult); err != nil {
		fmt.Println(err.Error())
		return []models.ModelResultReportMap{}, err
	}

	return modelResult, nil
}

func (r *Repository) GetGrade(ids []int) ([]models.RefGrade, error) {
	var grade []models.RefGrade
	query := r.conn
	query = query.Where("id IN (?)", ids)
	err := query.Find(&grade).Error
	if err != nil {
		return grade, err
	}

	return grade, nil
}

func (r *Repository) GetRoadConditionGradesByID(ID int) ([]models.ParamsConditionPreload, error) {
	var paramsCondition []models.ParamsConditionPreload
	query := r.conn
	query = query.Preload("RefOwner")
	query = query.Preload("RefGrade")
	query = query.Where("ref_owner_id = ?", ID)
	query = query.Where("condition_type = ?", "IRI")
	err := query.Find(&paramsCondition).Error
	if err != nil {
		return paramsCondition, err
	}

	return paramsCondition, nil
}

func (r *Repository) GetRefOwnerCondition(ID int) (map[string][]responses.Grade, error) {
	gradeData := make(map[string][]responses.Grade)
	var data []models.ParamsCondition //[]responses.Grade
	query := r.conn
	if err := query.Table("params_condition").Select("name,color,left_value_cc,left_condition_cc,right_value_cc,right_condition_cc,left_value_ac,left_condition_ac,right_value_ac,right_condition_ac, ref_grade_id").Joins("INNER JOIN ref_grade on params_condition.ref_grade_id = ref_grade.id").Where("ref_owner_id = ?", ID).Where("condition_type = ?", "IRI").Find(&data).Error; err != nil {
		return gradeData, err
	}

	var refGrades []models.RefGrade
	refGradeData := make(map[int]models.RefGrade)
	if err := query.Find(&refGrades).Error; err != nil {
		return gradeData, err
	}

	for _, item := range refGrades {
		refGradeData[item.ID] = item
	}

	for _, item := range data {
		refGrade := refGradeData[item.RefGradeID]
		gradeAC := responses.Grade{LeftValue: item.LeftValueAC, RightValue: item.RightValueAC, Color: refGrade.Color, Name: fmt.Sprintf("%.1f <= %s < %.1f", item.LeftValueAC, refGrade.Name, item.RightValueAC)}
		gradeCC := responses.Grade{LeftValue: item.LeftValueCC, RightValue: item.RightValueCC, Color: refGrade.Color, Name: fmt.Sprintf("%.1f <= %s < %.1f", item.LeftValueCC, refGrade.Name, item.RightValueCC)}
		gradeData["AC"] = append(gradeData["AC"], gradeAC)
		gradeData["CC"] = append(gradeData["CC"], gradeCC)
	}

	return gradeData, nil
}

func (r *Repository) GetPrepareDataIdByMaintenanceAnalysisId(id int) ([]int, error) {
	var prepareDataId []int
	query := r.conn
	query = query.Table("prepare_data")
	query = query.Select("prepare_data.id")
	query = query.Joins("LEFT JOIN maintenance_analysis on prepare_data.maintenance_analysis_id = maintenance_analysis.id")
	query = query.Where("prepare_data.maintenance_analysis_id = ?", id)
	query = query.Where("maintenance_analysis.prepare_data_status = ?", true)
	query = query.Order("prepare_data.id ASC")
	err := query.Find(&prepareDataId).Error
	if err != nil {
		return []int{}, err
	}
	return prepareDataId, nil
}

func (r *Repository) GetPrepareDataById(id int) ([]responses.PrepareDataWithPagination, error) {

	data, err := r.GetMaintenanceAnalysisById(id)
	if err != nil {
		return []responses.PrepareDataWithPagination{}, err
	}
	query := r.conn
	aadt1 := data.Aadt1
	aadt2 := data.Aadt2
	iri1 := data.Iri1
	iri2 := data.Iri2
	ifi1 := data.Ifi1
	ifi2 := data.Ifi2
	if aadt1 != nil && aadt2 != nil {
		query = query.Where("aadt > ? AND aadt < ?", aadt1, aadt2)
	} else if aadt1 != nil {
		query = query.Where("aadt > ?", aadt1)
	} else if aadt2 != nil {
		query = query.Where("aadt < ?", aadt2)
	}

	if iri1 != nil && iri2 != nil {
		query = query.Where("prepare_data.iri > ? AND prepare_data.iri < ?", iri1, iri2)
	} else if iri1 != nil {
		query = query.Where("prepare_data.iri > ?", iri1)
	} else if iri2 != nil {
		query = query.Where("prepare_data.iri < ?", iri2)
	}

	if ifi1 != nil && ifi2 != nil {
		query = query.Where("prepare_data.ifi > ? AND prepare_data.ifi < ?", ifi1, ifi2)
	} else if ifi1 != nil {
		query = query.Where("prepare_data.ifi > ?", ifi1)
	} else if ifi2 != nil {
		query = query.Where("prepare_data.ifi < ?", ifi2)
	}
	helpers.PrintlnJson("SurfaceTypeId", data.SurfaceTypeId)
	if data.SurfaceTypeId == 1 {
		surfaceType := []int{1, 2, 3, 4}
		query = query.Where("last_type in (?)", surfaceType)
	} else {
		surfaceType := []int{5}
		query = query.Where("last_type in (?)", surfaceType)
	}
	var prepareData []responses.PrepareDataWithPagination

	query = query.Table("prepare_data")
	query = query.Select("prepare_data.id")
	query = query.Joins("LEFT JOIN maintenance_analysis on prepare_data.maintenance_analysis_id = maintenance_analysis.id")
	query = query.Where("prepare_data.maintenance_analysis_id = ?", id)
	query = query.Where("maintenance_analysis.prepare_data_status = ?", true)
	query = query.Order("prepare_data.id ASC")
	err = query.Debug().Find(&prepareData).Error
	if err != nil {
		return prepareData, err
	}
	IDs := []int{}
	for _, item := range prepareData {
		IDs = append(IDs, item.ID)
	}

	query2 := r.conn
	if err := query2.Where("maintenance_analysis_id = ?", id).Where("id NOT IN (?)", IDs).Delete(&models.PrepareData{}).Error; err != nil {
		return prepareData, err
	}

	var prepareData2 []responses.PrepareDataWithPagination
	query3 := r.conn
	query3 = query3.Table("prepare_data")
	query3 = query3.Select("prepare_data.*")
	query3 = query3.Joins("LEFT JOIN maintenance_analysis on prepare_data.maintenance_analysis_id = maintenance_analysis.id")
	query3 = query3.Where("prepare_data.maintenance_analysis_id = ?", id)
	query3 = query3.Where("maintenance_analysis.prepare_data_status = ?", true)
	query3 = query3.Order("prepare_data.id ASC")
	err = query3.Debug().Find(&prepareData2).Error
	if err != nil {
		return prepareData2, err
	}

	return prepareData2, nil
}

func (r *Repository) UpdatePrepareDataStatusById(id int) error {
	query := r.conn
	if err := query.Model(&models.MaintenanceAnalysis{}).Where("id = ?", id).Update("prepare_data_status", true).Error; err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *Repository) GetPrepareDataAllByAnalysisSelected(ID int) ([]models.PrepareData, error) {
	query := r.conn
	var prepareData []models.PrepareData
	if err := query.Table("prepare_data").Where("maintenance_analysis_id = ?", ID).Where("is_selected = true").Find(&prepareData).Error; err != nil {
		return prepareData, err
	}
	return prepareData, nil
}

func (r *Repository) GetRefDepotByRoad() (map[int]string, error) {
	query := r.conn
	var data []responses.RefDepotByRoad
	if err := query.Select("ref_depot_code, road.id as road_id").Table("road").Joins("JOIN road_section on road.road_section_id = road_section.id").Where("road.is_active = true").Find(&data).Error; err != nil {
		return map[int]string{}, err
	}

	dataRes := make(map[int]string)
	for _, item := range data {
		dataRes[item.RoadID] = item.RefDepotCode
	}
	return dataRes, nil
}

func (r *Repository) DashboardMapAllYear(ID, analysisTypeID int, plan *int) ([]models.ModelResultReportMap, error) {
	query := r.connMon.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_MAINTENANCE_ANALYSIS_RESULT_COLLECTION"))
	var pipeline []primitive.M
	if analysisTypeID == 1 {
		if plan != nil {
			helpers.PrintlnJson("dfdfdfdfd")
			pipeline = []bson.M{
				{
					"$match": bson.M{
						"maintenance_analysis_id": ID,
						// "plan":                    plan,
						"maintenance_analysis_type_id": analysisTypeID,
					},
				},
			}
		} else {
			pipeline = []bson.M{
				{
					"$match": bson.M{
						"maintenance_analysis_id":      ID,
						"maintenance_analysis_type_id": analysisTypeID,
					},
				},
			}
		}

	} else {
		pipeline = []bson.M{
			{
				"$match": bson.M{
					"maintenance_analysis_id":      ID,
					"maintenance_analysis_type_id": analysisTypeID,
				},
			},
		}
	}

	cursor, err := query.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return []models.ModelResultReportMap{}, err
	}

	defer cursor.Close(context.Background())
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	var modelResult []models.ModelResultReportMap
	if err = cursor.All(context.TODO(), &modelResult); err != nil {
		fmt.Println(err.Error())
		return []models.ModelResultReportMap{}, err
	}

	return modelResult, nil
}
