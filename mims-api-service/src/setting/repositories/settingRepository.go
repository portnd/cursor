package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/setting/handlers"
	"gitlab.com/mims-api-service/src/setting/usecases"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"gorm.io/gorm"
)

type settingRepository struct {
	conn    *gorm.DB
	connMon *mongo.Client
}

func NewSettingRepositoryHandler(conn *gorm.DB, connMon *mongo.Client) *handlers.SettingHandler {
	useCase := usecases.NewSettingUseCase(&settingRepository{conn, connMon})
	handler := handlers.NewSettingHandler(useCase)
	return handler
}

// this repo below can be used in any setting usecase
func (sr *settingRepository) GetAll(records interface{}) error {
	err := sr.conn.Where("status = ?", 1).Order("id ASC").Find(records).Error
	if err != nil {
		return err
	}

	return nil
}

func (sr *settingRepository) GetByID(record interface{}, id int) error {
	err := sr.conn.Where("status = ? AND id = ?", 1, id).First(record).Error
	if err != nil {
		return err
	}

	return nil
}

func (sr *settingRepository) UpdateByID(tableName, updateName string, id int) error {
	err := sr.conn.Table(tableName).Where("id = ?", id).Update("name", updateName).Error
	if err != nil {
		return err
	}

	return nil
}

func (sr *settingRepository) DeleteByID(tableName string, id int) error {
	err := sr.conn.Table(tableName).Where("id = ?", id).Update("status", 0).Error
	if err != nil {
		return err
	}
	return nil
}

// be used in settings/asset_groups route
func (sr *settingRepository) CountAssetGroups(searchName string) (int64, error) {
	var count int64
	if err := sr.conn.Model(&models.RefAsset{}).Where("status = ?", 1).Where("name ILIKE ?", "%"+searchName+"%").Count(&count).Error; err != nil {
		return int64(0), err
	}
	return int64(count), nil
}

func (sr *settingRepository) GetAssetGroups(limit, offset int64, searchName string) ([]models.RefAsset, error) {
	assetGroups := []models.RefAsset{}
	query := sr.conn
	query = query.Where("status = ?", 1).Where("name ILIKE ?", "%"+searchName+"%").Order("id").Limit(int(limit)).Offset(int(offset))
	err := query.Find(&assetGroups).Error
	if err != nil {
		log.Println(err)
		return assetGroups, err
	}
	return assetGroups, nil
}

func (sr *settingRepository) CreateAssetGroup(assetGroup models.RefAsset) error {
	var maxSeq int
	query := sr.conn

	if err := query.Model(&models.RefAsset{}).Select("MAX(seq)").Row().Scan(&maxSeq); err != nil {
		return err
	}

	if err := query.Model(&models.RefAsset{}).Where("name = ?", "อื่น ๆ").Update("seq", maxSeq+1).Error; err != nil {
		return err
	}

	assetGroup.Seq = maxSeq

	if err := query.Create(&assetGroup).Error; err != nil {
		return err
	}

	return nil
}

// settings/departments
func (sr *settingRepository) CountDepartments(searchName string) (int64, error) {
	var count int64
	if err := sr.conn.Model(&models.RefDepartment{}).Where("status = ?", 1).Where("name ILIKE ?", "%"+searchName+"%").Count(&count).Error; err != nil {
		return int64(0), err
	}
	return int64(count), nil
}

func (sr *settingRepository) GetDepartments(limit, offset int64, searchName string) ([]models.RefDepartment, error) {
	departments := []models.RefDepartment{}
	query := sr.conn
	query = query.Where("status = ?", 1).Where("name ILIKE ?", "%"+searchName+"%").Order("id").Limit(int(limit)).Offset(int(offset))
	err := query.Find(&departments).Error
	if err != nil {
		log.Println(err)
		return departments, err
	}
	return departments, nil
}

func (sr *settingRepository) CreateDepartment(department models.RefDepartment) error {
	err := sr.conn.Create(&department).Error
	if err != nil {
		return err
	}
	return nil
}

// be used in settings/condition_lists and settings/owner_lists route
func (sr *settingRepository) GetParamsCondition(ownerId int) ([]models.ParamsConditionPreload, error) {
	records := []models.ParamsConditionPreload{}

	err := sr.conn.Where("ref_owner_id = ?", ownerId).Preload("RefOwner").Preload("RefOwner.RefConditionRange").Preload("RefGrade").Order("id ASC").Find(&records).Error
	if err != nil {
		return []models.ParamsConditionPreload{}, err
	}

	return records, nil
}

func (sr *settingRepository) GetParamsRoadLine(ownerId int) ([]models.ParamsRoadLinePreload, error) {
	records := []models.ParamsRoadLinePreload{}

	err := sr.conn.Where("ref_owner_road_line_id = ?", ownerId).Preload("RefOwnerRoadLine").Preload("RefGrade").Order("id ASC").Find(&records).Error
	//.Preload("RefConditionRange")
	if err != nil {
		return []models.ParamsRoadLinePreload{}, err
	}

	return records, nil
}

func (sr *settingRepository) CountOwnersRoadLine(params requests.QueryParamsReflectivityRange) (int64, error) {
	var count int64
	query := sr.conn
	query = query.Model(&models.RefOwnerRoadLine{})
	query = query.Where("is_active = ?", true)
	query = query.Where("name ILIKE ?", "%"+params.Name+"%")
	if params.RefReflectivityRangeID != 0 {
		query = query.Where("ref_reflectivity_range_id = ?", params.RefReflectivityRangeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return int64(0), err
	}
	return int64(count), nil
}

func (sr *settingRepository) GetOwnersRoadLine(limit, offset int64, params requests.QueryParamsReflectivityRange) ([]models.RefOwnerRoadLinePreload, error) {
	owners := []models.RefOwnerRoadLinePreload{}
	query := sr.conn
	query = query.Preload("RefReflectivityRange")
	query = query.Where("is_active = ?", true).Where("name ILIKE ?", "%"+params.Name+"%")
	if params.RefReflectivityRangeID != 0 {
		query = query.Where("ref_reflectivity_range_id = ?", params.RefReflectivityRangeID)
	}
	query = query.Order("id").Limit(int(limit)).Offset(int(offset))
	err := query.Find(&owners).Error
	if err != nil {
		log.Println(err)
		return owners, err
	}
	return owners, nil
}

func (sr *settingRepository) GetOwnerRoadLineByID(id int) (models.RefOwnerRoadLine, error) {
	owner := models.RefOwnerRoadLine{}

	err := sr.conn.Where("id = ? AND is_active = ?", id, true).First(&owner).Error
	if err != nil {
		return owner, err
	}

	return owner, nil
}

func (sr *settingRepository) GetOwnersRoadLineAll() ([]models.RefOwnerRoadLine, error) {
	owners := []models.RefOwnerRoadLine{}
	err := sr.conn.Where("is_active = ?", true).Find(&owners).Error
	if err != nil {
		return owners, err
	}

	return owners, nil
}

func (sr *settingRepository) CreateOwnerRoadLine(owner *models.RefOwnerRoadLine) error {
	err := sr.conn.Create(&owner).Error
	if err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) UpdateOwnerRoadLineByID(id int, owner models.RefOwnerRoadLine) error {
	err := sr.conn.Model(&models.RefOwnerRoadLine{}).Where("id = ?", id).Updates(owner).Error
	if err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) DeleteOwnerRoadLineByID(ownerId int) error {
	err := sr.conn.Model(&models.RefOwnerRoadLine{}).Where("id = ?", ownerId).Update("is_active", false).Error
	if err != nil {
		return err
	}

	var paramsCondition models.ParamsRoadLine
	if err := sr.conn.Where("ref_owner_road_line_id = ?", ownerId).Delete(&paramsCondition).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) CountOwners(params requests.QueryParams) (int64, error) {
	var count int64
	query := sr.conn
	query = query.Where("is_active = ?", true)
	query = query.Where("name ILIKE ?", "%"+params.Name+"%")
	helpers.PrintlnJson("params.RefConditionRangeID", params.RefConditionRangeID)
	if params.RefConditionRangeID != 0 {
		query = query.Where("ref_condition_range_id = ?", params.RefConditionRangeID)
	}
	if err := query.Model(&models.RefOwner{}).Count(&count).Error; err != nil {
		return int64(0), err
	}
	return int64(count), nil
}

func (sr *settingRepository) GetOwners(limit, offset int64, params requests.QueryParams) ([]models.RefOwnerPreload, error) {
	owners := []models.RefOwnerPreload{}
	query := sr.conn
	query = query.Preload("RefConditionRange")
	if params.RefConditionRangeID != 0 {
		query = query.Where("ref_condition_range_id = ?", params.RefConditionRangeID)
	}
	query = query.Where("is_active = ?", true)
	query = query.Where("name ILIKE ?", "%"+params.Name+"%")
	query = query.Order("id").Limit(int(limit)).Offset(int(offset))
	err := query.Find(&owners).Error
	if err != nil {
		log.Println(err)
		return owners, err
	}
	return owners, nil
}

func (sr *settingRepository) GetOwnerByID(id int) (models.RefOwner, error) {
	owner := models.RefOwner{}

	err := sr.conn.Where("id = ? AND is_active = ?", id, true).First(&owner).Error
	if err != nil {
		return owner, err
	}

	return owner, nil
}

func (sr *settingRepository) GetOwnersAll() ([]models.RefOwner, error) {
	owners := []models.RefOwner{}
	err := sr.conn.Where("is_active = ?", true).Find(&owners).Error
	if err != nil {
		return owners, err
	}

	return owners, nil
}

func (sr *settingRepository) CreateOwner(owner *models.RefOwner) error {
	err := sr.conn.Create(&owner).Error
	if err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) UpdateOwnerByID(id int, owner models.RefOwner) error {
	err := sr.conn.Model(&models.RefOwner{}).Where("id = ?", id).Updates(owner).Error
	if err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) DeleteOwnerByID(ownerId int) error {
	err := sr.conn.Model(&models.RefOwner{}).Where("id = ?", ownerId).Update("is_active", false).Error
	if err != nil {
		return err
	}

	var paramsCondition models.ParamsCondition
	if err := sr.conn.Where("ref_owner_id = ?", ownerId).Delete(&paramsCondition).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetOwnerByRoadID(roadId int) (models.RoadOwner, error) {
	roadOwner := models.RoadOwner{}

	err := sr.conn.Where("road_id = ?", roadId).First(&roadOwner).Error
	if err != nil {
		return models.RoadOwner{}, err
	}

	return roadOwner, nil
}

// func (sr *settingRepository) CreateCondition(ownerId int, newRecords []models.ParamsCondition) error {
// 	valueStrings := []string{}
// 	valueArgs := []interface{}{}
// 	for _, t := range newRecords {
// 		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")
// 		valueArgs = append(valueArgs, ownerId)
// 		valueArgs = append(valueArgs, t.RefGradeID)
// 		valueArgs = append(valueArgs, t.LeftValue)
// 		valueArgs = append(valueArgs, t.LeftCondition)
// 		valueArgs = append(valueArgs, t.RightValue)
// 		valueArgs = append(valueArgs, t.RightCondition)
// 		valueArgs = append(valueArgs, t.ConditionType)
// 	}
// 	query := sr.conn
// 	smt := `INSERT INTO params_condition(ref_owner_id, ref_grade_id, left_value, left_condition, right_value, right_condition, condition_type) VALUES %s `
// 	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

// 	tx := query.Begin()
// 	if err := tx.Exec(smt, valueArgs...).Error; err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	tx.Commit()
// 	return nil
// }

func (sr *settingRepository) DeleteCondition(ownerId int) error {
	var paramsCondition models.ParamsCondition
	if err := sr.conn.Where("ref_owner_id = ?", ownerId).Delete(&paramsCondition).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) DeleteRoadLine(ownerId int) error {
	var paramsRoadLine models.ParamsRoadLine
	if err := sr.conn.Where("ref_owner_road_line_id = ?", ownerId).Delete(&paramsRoadLine).Error; err != nil {
		return err
	}
	return nil
}

// be used in settings/signs route
func (sr *settingRepository) CountSigns(searchName string) (int64, error) {
	var count int64
	if err := sr.conn.Model(&models.RefAssetSignImage{}).Where("status = ?", 1).Where("name ILIKE ?", "%"+searchName+"%").Count(&count).Error; err != nil {
		return int64(0), err
	}
	return int64(count), nil
}

func (sr *settingRepository) GetSigns(limit, offset int64, searchName string) ([]models.RefAssetSignImage, error) {
	signs := []models.RefAssetSignImage{}
	query := sr.conn
	query = query.Where("status = ?", 1).Where("name ILIKE ?", "%"+searchName+"%").Order("id").Limit(int(limit)).Offset(int(offset))
	err := query.Find(&signs).Error
	if err != nil {
		log.Println(err)
		return signs, err
	}
	return signs, nil
}

func (sr *settingRepository) CreateSignImage(record *models.RefAssetSignImage) error {
	err := sr.conn.Create(&record).Error
	if err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) UpdateSignImage(id int, record models.RefAssetSignImage) error {
	err := sr.conn.Model(&models.RefAssetSignImage{}).Where("id = ?", id).Updates(record).Error
	if err != nil {
		return err
	}
	return nil
}

// be used in settings/asset_tables route
func (sr *settingRepository) CountAssetTables(assetType, name string, groupID int) (int64, error) {
	var count int64
	query := sr.conn
	query = query.Model(&models.RefAssetTable{})
	if name != "" {
		query = query.Where("table_label ILIKE ?", "%"+name+"%")
	}
	if groupID != 0 {
		query = query.Where("ref_asset_id = ?", groupID)
	}
	if assetType != "" {
		query = helpers.GetAssetTypeCondition(query, assetType)
	}

	err := query.Where("is_active = ?", true).Count(&count).Error
	if err != nil {
		return int64(0), err
	}
	return int64(count), nil
}

func (sr *settingRepository) GetAssetTables(limit, offset int64, assetType, name string, groupID int) ([]models.AssetTable, error) {
	assetTable := []models.AssetTable{}
	query := sr.conn
	query = query.Model(&models.RefAssetTable{})
	if name != "" {
		query = query.Where("table_label ILIKE ?", "%"+name+"%")
	}
	if groupID != 0 {
		query = query.Where("ref_asset_id = ?", groupID)
	}
	if assetType != "" {
		query = helpers.GetAssetTypeCondition(query, assetType)
	}

	query = query.Where("is_active = ?", true)
	// query := sr.conn.Where("is_active = ?", true).Where("table_label LIKE ?", "%"+name+"%")
	// if assetType != "" {
	// 	query = helpers.GetAssetTypeCondition(query, assetType)
	// }

	err := query.Preload("RefAsset").Order("id ASC").Limit(int(limit)).Offset(int(offset)).Find(&assetTable).Error
	if err != nil {
		return assetTable, err
	}

	return assetTable, nil
}

func (sr *settingRepository) GetAssetTableStaffs() ([]models.AssetTableStaff, error) {
	assetTableStaffs := []models.AssetTableStaff{}

	err := sr.conn.Preload("RefDepartment").Select("ref_asset_table_id, ref_department_id").
		Group("1, 2").Order("1, 2").
		Find(&assetTableStaffs).Error
	if err != nil {
		return assetTableStaffs, err
	}

	return assetTableStaffs, nil
}

func (sr *settingRepository) GetAssetTableByID(id int) (models.AssetTable, error) {
	assetTable := models.AssetTable{}

	err := sr.conn.Where("id = ?", id).Preload("RefAsset").Preload("AssetTableColumns").First(&assetTable).Error
	if err != nil {
		return assetTable, err
	}

	return assetTable, nil
}

func (sr *settingRepository) GetAssetTableStaffByID(id int) ([]models.AssetTableStaff, error) {
	assetTableStaffs := []models.AssetTableStaff{}

	err := sr.conn.Where("ref_asset_table_id = ?", id).Preload("RefDepartment").Select("ref_asset_table_id, ref_department_id, is_approver").
		Group("1, 2, 3").Order("1, 2").
		Find(&assetTableStaffs).Error
	if err != nil {
		return assetTableStaffs, err
	}

	return assetTableStaffs, nil
}

func (r *settingRepository) CreateData(model interface{}) error {

	var err error
	err = r.conn.CreateInBatches(model, 100).Error
	if err == nil {
		return nil // Success, return nil error
	}
	// Return the last error after all retries
	return err
}

func (sr *settingRepository) CreateAssetTable(data helpers.CreateAssetTable) error {
	tx := sr.conn.Begin()

	if err := tx.Exec(data.QueryString).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := sr.InsertRefAssetTable(tx, &data.AssetTable); err != nil {
		tx.Rollback()
		return err
	}

	for _, column := range data.InsertColumns {
		column.RefAssetTableID = data.AssetTable.ID
		if err := sr.InsertRefAssetTableColumns(tx, column); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, staff := range data.Staffs {
		staff.RefAssetTableID = data.AssetTable.ID
		if err := sr.InsertRefAssetTableStaffs(tx, staff); err != nil {
			tx.Rollback()
			return err
		}

	}

	tx.Commit()

	return nil
}

func (sr *settingRepository) InsertRefAssetTable(tx *gorm.DB, record *models.RefAssetTable) error {
	err := tx.Create(&record).Error
	if err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) InsertRefAssetTableColumns(tx *gorm.DB, record models.RefAssetTableColumns) error {
	err := tx.Create(&record).Error
	if err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) InsertRefAssetTableStaffs(tx *gorm.DB, record models.RefAssetTableStaff) error {
	err := tx.Create(&record).Error
	if err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) UpdateAssetTable(data helpers.UpdateAssetTable) error {
	tx := sr.conn.Begin()
	// update asset teble data
	if err := tx.Save(&data.AssetTable).Error; err != nil {
		tx.Rollback()
		return err
	}

	// delete column
	for _, columnId := range data.DeleteColumns {
		err := tx.Where("id = ?", columnId).Delete(&models.RefAssetTableColumns{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if data.RenameColumnsQuery != "" {
		if err := tx.Exec(data.RenameColumnsQuery).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// add new column
	for _, column := range data.InsertNewColumns {
		column.RefAssetTableID = data.AssetTable.ID
		if err := tx.Create(&column).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(data.UpdateColumns) > 0 {
		for _, column := range data.UpdateColumns {
			if err := tx.Model(models.RefAssetTableColumns{}).Where("id = ?", column.ID).Save(&column).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if data.AddColumnsQuery != "" {
		if err := tx.Exec(data.AddColumnsQuery).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// remove old staff
	if err := tx.Where("ref_asset_table_id = ?", data.AssetTable.ID).Delete(&models.RefAssetTableStaff{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// add new staff
	for _, staff := range data.Staffs {
		staff.RefAssetTableID = data.AssetTable.ID
		if err := tx.Create(&staff).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// update comment for asset table and asset table columns
	if err := tx.Exec(data.CommentQuery).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (sr *settingRepository) GetColumnsByID(id int) (models.RefAssetTableColumns, error) {
	column := models.RefAssetTableColumns{}
	err := sr.conn.Where("id = ?", id).First(&column).Error
	if err != nil {
		return column, err
	}

	return column, nil
}

func (sr *settingRepository) GetColumnMaxSeqByAssetTableID(id int) (int, error) {
	var maxSeq int
	err := sr.conn.Model(&models.RefAssetTableColumns{}).Select("MAX(column_seq)").Where("ref_asset_table_id = ?", id).Row().Scan(&maxSeq)
	if err != nil {
		return 0, err
	}

	return maxSeq, nil
}

func (sr *settingRepository) GetAllAssetTables() ([]models.AssetTable, error) {
	assetTables := []models.AssetTable{}

	err := sr.conn.Find(&assetTables).Error
	if err != nil {
		return assetTables, err

	}

	return assetTables, nil
}

func (sr *settingRepository) CountAssetTableColumns(tableName string) (int, error) {
	var count int

	queryString := fmt.Sprintf("SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '%s';", tableName)

	err := sr.conn.Raw(queryString).Row().Scan(&count)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (sr *settingRepository) DeleteAssetTable(id int) error {
	err := sr.conn.Model(&models.RefAssetTable{}).Where("id = ?", id).Update("is_active", false).Error
	if err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetGrade() ([]models.RefGrade, error) {
	var grades []models.RefGrade
	if err := sr.conn.Find(&grades).Error; err != nil {
		return grades, err
	}
	return grades, nil
}

func (sr *settingRepository) InsertRefSurfaceParam(tx *gorm.DB, dataToInsert models.RefSurfaceParam) error {
	err := tx.Table("ref_surface_param").Where("ref_surface_id = ?", dataToInsert.RefSurfaceID).Update("is_latest", false).Error
	if err != nil {
		return err
	}
	err = tx.Table("ref_surface_param").Create(&dataToInsert).Error
	if err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetLatestSurfaceParamsById(id int) (models.RefSurfaceParam, error) {
	var latestSurfaceParams models.RefSurfaceParam
	err := sr.conn.Table("ref_surface_param").Where("is_latest = ? AND ref_surface_id = ?", true, id).Find(&latestSurfaceParams).Error
	if err != nil {
		return latestSurfaceParams, err
	}
	return latestSurfaceParams, nil
}

func (sr *settingRepository) CreateInterventionCriteriaParams(settingInterventionCriteria *models.SettingInterventionCriteria) error {
	if err := sr.conn.Table("setting_intervention_criteria_params").Save(&settingInterventionCriteria).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) UpdateInterventionCriteriaParamsByIsLatestIsFalse() error {
	if err := sr.conn.Table("setting_intervention_criteria_params").Where("is_latest = true").Update("is_latest", false).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) UpdateInterventionCriteria(interventionCriteria *models.InterventionCriteria) error {
	if err := sr.conn.Table("setting_intervention_criteria").Save(&interventionCriteria).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetInterventionCriteriaById(id int) (models.InterventionCriteria, error) {
	var interventionCriteria models.InterventionCriteria
	if err := sr.conn.Table("setting_intervention_criteria").Where("is_deleted = ? and id = ?", false, id).Find(&interventionCriteria).Error; err != nil {
		return models.InterventionCriteria{}, err
	}
	return interventionCriteria, nil
}

func (sr *settingRepository) GetInterventionCriteriaByNotId(id int) ([]models.InterventionCriteria, error) {
	var interventionCriteria []models.InterventionCriteria
	if err := sr.conn.Table("setting_intervention_criteria").Where("is_deleted = ? and id != ?", false, id).Find(&interventionCriteria).Error; err != nil {
		return []models.InterventionCriteria{}, err
	}
	return interventionCriteria, nil
}

func (sr *settingRepository) GetRefCriteriaMethod() ([]models.RefCriteriaMethod, error) {
	var refCriteriaMethod []models.RefCriteriaMethod
	if err := sr.conn.Table("ref_criteria_method").Where("is_deleted = ?", false).Order("id ASC").Find(&refCriteriaMethod).Error; err != nil {
		return []models.RefCriteriaMethod{}, err
	}
	return refCriteriaMethod, nil
}

func (sr *settingRepository) GetInterventionCriteria() ([]models.InterventionCriteria, error) {
	var interventionCriteria []models.InterventionCriteria
	if err := sr.conn.Table("setting_intervention_criteria").Where("is_deleted = ?", false).Order("maintenance_sequence ASC").Where("is_show = ?", true).Find(&interventionCriteria).Error; err != nil {
		return []models.InterventionCriteria{}, err
	}
	return interventionCriteria, nil
}

func (sr *settingRepository) CountInterventionCriteriaByMaintenanceMethod(surface string) (models.InterventionCriteriaCount, error) {
	var interventionCriteriaCount models.InterventionCriteriaCount
	query := sr.conn
	query = query.Table("setting_intervention_criteria").
		Select("COUNT(setting_intervention_criteria.id) as count").
		Joins("LEFT JOIN ref_criteria_method ON ref_criteria_method.id = setting_intervention_criteria.maintenance_method").
		Where("setting_intervention_criteria.is_deleted = ? AND ref_criteria_method.surface = ? AND setting_intervention_criteria.is_show = ? ", false, surface, true)

	query = query.Find(&interventionCriteriaCount)
	err := query.Error
	if err != nil {
		return models.InterventionCriteriaCount{}, err
	}

	return interventionCriteriaCount, nil
}

func (sr *settingRepository) GetInterventionCriteriaConditionSequence() ([]models.InterventionCriteriaCondition, error) {
	var interventionCriteriaCondition []models.InterventionCriteriaCondition
	if err := sr.conn.Table("setting_intervention_criteria_condition").Where("is_deleted = ?", false).Order("condition_sequence ASC").Find(&interventionCriteriaCondition).Error; err != nil {
		return []models.InterventionCriteriaCondition{}, err
	}
	return interventionCriteriaCondition, nil
}

func (sr *settingRepository) UpdateInterventionCriteriaCondition(interventionCriteriaCondition *models.InterventionCriteriaCondition) error {
	if err := sr.conn.Table("setting_intervention_criteria_condition").Save(&interventionCriteriaCondition).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) DeleteInterventionCriteriaConditionById(id int) error {
	if err := sr.conn.Table("setting_intervention_criteria_condition").Where("intervention_criteria_id = ?", id).Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetInterventionCriteriaConditionById(id int) (models.InterventionCriteriaCondition, error) {
	var interventionCriteriaCondition models.InterventionCriteriaCondition
	if err := sr.conn.Table("setting_intervention_criteria_condition").Where("is_deleted = ? and id = ?", false, id).Order("condition_sequence ASC").Find(&interventionCriteriaCondition).Error; err != nil {
		return models.InterventionCriteriaCondition{}, err
	}
	return interventionCriteriaCondition, nil
}

func (sr *settingRepository) GetInterventionCriteriaConditionListByInterventionCriteriaId(id int) ([]models.InterventionCriteriaCondition, error) {
	var interventionCriteriaCondition []models.InterventionCriteriaCondition
	if err := sr.conn.Table("setting_intervention_criteria_condition").Where("is_deleted = ? and intervention_criteria_id = ?", false, id).Order("condition_sequence ASC").Find(&interventionCriteriaCondition).Error; err != nil {
		return []models.InterventionCriteriaCondition{}, err
	}
	return interventionCriteriaCondition, nil
}

func (sr *settingRepository) GetCriteriaMethodById(id int) (models.RefCriteriaMethod, error) {
	var refCriteriaMethod models.RefCriteriaMethod
	if err := sr.conn.Table("ref_criteria_method").Where("is_deleted = ? and id = ?", false, id).Find(&refCriteriaMethod).Error; err != nil {
		return models.RefCriteriaMethod{}, err
	}
	return refCriteriaMethod, nil
}

func (sr *settingRepository) GetCriteriaMethodByName(name string) (models.RefCriteriaMethod, error) {
	var refCriteriaMethod models.RefCriteriaMethod
	if err := sr.conn.Table("ref_criteria_method").Where("is_deleted = ? and name = ?", false, name).Find(&refCriteriaMethod).Error; err != nil {
		return models.RefCriteriaMethod{}, err
	}
	return refCriteriaMethod, nil
}

func (sr *settingRepository) GetCriteriaMethod() ([]models.RefCriteriaMethod, error) {
	var refCriteriaMethod []models.RefCriteriaMethod
	if err := sr.conn.Table("ref_criteria_method").Find(&refCriteriaMethod).Error; err != nil {
		return []models.RefCriteriaMethod{}, err
	}
	return refCriteriaMethod, nil
}

func (sr *settingRepository) InsertRefSurface(tx *gorm.DB, data models.NewRefSurface) (int, error) {
	var refSurface models.NewRefSurface
	data.CanDelete = true
	err := tx.Table("ref_surface").Where("name = ?", data.Name).First(&refSurface).Error
	if err == gorm.ErrRecordNotFound {
		err = tx.Table("ref_surface").Create(&data).Last(&refSurface).Error
		if err != nil {
			return 0, err
		}
		return refSurface.ID, nil
	} else if err != nil {
		return 0, err
	}
	return 0, errors.New("This name already use.")

}

func (sr *settingRepository) UpdateRefSurface(tx *gorm.DB, data models.NewRefSurface) (bool, error) {
	var refSurface models.NewRefSurface
	var checkName models.NewRefSurface
	if !(data.ID > 0) {
		return false, errors.New("Invalid ID")
	}
	err := sr.conn.Table("ref_surface").Where("id = ?", data.ID).First(&refSurface).Error
	if err != nil {
		return false, err
	}
	err = tx.Table("ref_surface").Where("name = ?", data.Name).First(&checkName).Error
	if err == gorm.ErrRecordNotFound {
		data.CanDelete = refSurface.CanDelete
		err = tx.Table("ref_surface").Save(&data).Error
		if err != nil {
			return false, err
		}
		return data.CanDelete, nil
	} else if data.Name == checkName.Name && data.ID == checkName.ID {
		data.CanDelete = checkName.CanDelete
		err = tx.Table("ref_surface").Save(&data).Error
		if err != nil {
			return false, err
		}
		return data.CanDelete, nil
	} else if err != nil {
		return false, err
	}
	return false, errors.New("This name already use.")

}

func (r *settingRepository) GetRefSurface(condition string) ([]models.NewRefSurface, error) {
	var results []models.NewRefSurface
	query := r.conn.Table("ref_surface")
	if condition != "" {
		query = query.Where(condition)
	}
	err := query.Find(&results).Error
	if err != nil {
		return results, err
	}
	return results, nil
}

func (r *settingRepository) GetRefSurfaceByID(id int) (models.NewRefSurface, error) {
	var result models.NewRefSurface
	err := r.conn.Table("ref_surface").Where("id = ?", id).First(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *settingRepository) GetNewRefSurfaceByID(id int) (models.RefSurfaceNew, error) {
	var result models.RefSurfaceNew
	err := r.conn.Select("*").Where("id = ?", id).First(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *settingRepository) GetParamRefSurfaceByID(id int) ([]models.RefSurfaceParam, error) {
	var results []models.RefSurfaceParam
	var err error
	if id != 0 {
		err = r.conn.Table("ref_surface_param").Where("ref_surface_id = ? and is_latest = ?", id, true).Find(&results).Error
		if len(results) == 0 {
			err = gorm.ErrRecordNotFound
		}
	} else {
		err = r.conn.Table("ref_surface_param").Where("is_latest = ?", true).Find(&results).Error

	}
	if err != nil {
		return results, err
	}
	return results, nil
}

func (r *settingRepository) StartTransSection() *gorm.DB {
	tx := r.conn.Begin()
	return tx
}

func (r *settingRepository) RollBack(tx *gorm.DB) error {
	tx.Rollback()
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func (r *settingRepository) Commit(tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (sr *settingRepository) GetRoadWorkEffectByIsLatest() (models.SettingRoadWorkEffect, error) {
	var settingRoadWorkEffect models.SettingRoadWorkEffect
	if err := sr.conn.Table("setting_road_work_effect_params").Where("is_latest = ?", true).Scan(&settingRoadWorkEffect).Error; err != nil {
		return settingRoadWorkEffect, err
	}
	return settingRoadWorkEffect, nil
}

func (sr *settingRepository) UpdateRoadWorkEffectByIsLatestIsFalse() error {
	if err := sr.conn.Table("setting_road_work_effect_params").Update("is_latest", false).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) CreateRoadWorkEffect(settingRoadWorkEffect *models.SettingRoadWorkEffect) error {
	if err := sr.conn.Table("setting_road_work_effect_params").Save(&settingRoadWorkEffect).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetBudget() ([]models.SettingBudget, error) {
	var budgetMethods []models.SettingBudget
	if err := sr.conn.Table("setting_budget").Where("is_deleted = ?", false).Find(&budgetMethods).Error; err != nil {
		return []models.SettingBudget{}, err
	}
	return budgetMethods, nil
}

func (sr *settingRepository) CountGetBudgetByName(name string) (int64, error) {
	var count int64
	if err := sr.conn.Table("setting_budget").Where("is_deleted = ?", false).Where("name ILIKE ?", "%"+name+"%").Count(&count).Error; err != nil {
		return int64(0), err
	}
	return int64(count), nil
}

func (sr *settingRepository) GetBudgetByName(limit, offset int64, name string) ([]responses.BudgetList, error) {
	var budgetList []responses.BudgetList
	query := sr.conn

	query = query.Table("setting_budget").Where("is_deleted = ?", false).Where("name ILIKE ?", "%"+name+"%").Order("id").Limit(int(limit)).Offset(int(offset))
	err := query.Find(&budgetList).Error
	if err != nil {
		log.Println(err)
		return []responses.BudgetList{}, err
	}
	return budgetList, nil
}

func (sr *settingRepository) UpdateBudget(budget *models.SettingBudget) error {
	if err := sr.conn.Table("setting_budget").Save(&budget).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetBudgetById(budgetId int) (models.SettingBudget, error) {
	var settingBudgetdget models.SettingBudget
	if err := sr.conn.Table("setting_budget").Where("is_deleted = ? AND id = ?", false, budgetId).Find(&settingBudgetdget).Error; err != nil {
		return models.SettingBudget{}, err
	}
	return settingBudgetdget, nil
}

func (sr *settingRepository) GetBudgetMethodByBudgetId(budgetId int) ([]models.SettingBudgetMethod, error) {
	var budgetMethod []models.SettingBudgetMethod
	if err := sr.conn.Table("setting_budget_method").Where("is_deleted = ? AND budget_id = ?", false, budgetId).Order("is_show_method DESC").Order("id ASC").Find(&budgetMethod).Error; err != nil {
		return []models.SettingBudgetMethod{}, err
	}
	return budgetMethod, nil
}

func (sr *settingRepository) GetBudgetMethodById(budgetId int) (models.SettingBudgetMethod, error) {
	var budgetMethod models.SettingBudgetMethod
	if err := sr.conn.Table("setting_budget_method").Where("is_deleted = ? AND id = ?", false, budgetId).Find(&budgetMethod).Error; err != nil {
		return models.SettingBudgetMethod{}, err
	}
	return budgetMethod, nil
}

func (sr *settingRepository) GetBudgetMethodListByBudgetId(budgetId int) ([]models.SettingBudgetMethod, error) {
	var budgetMethod []models.SettingBudgetMethod
	if err := sr.conn.Table("setting_budget_method").Where("is_deleted = ? AND budget_id = ?", false, budgetId).Find(&budgetMethod).Error; err != nil {
		return []models.SettingBudgetMethod{}, err
	}
	return budgetMethod, nil
}

func (sr *settingRepository) UpdateBudgetMethod(budgetMethod *models.SettingBudgetMethod) error {
	if err := sr.conn.Table("setting_budget_method").Save(&budgetMethod).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) UpdateAadtGrowthRateToFalse() error {
	if err := sr.conn.Table("setting_aadt_growth_rate").Update("is_latest", false).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetAadtGrowthRateByRoadGroupId(roadGroupId int) (models.AadtGrowthRate, error) {
	var aadtGrowthRate models.AadtGrowthRate
	if err := sr.conn.Table("setting_aadt_growth_rate").Where("is_deleted = ? AND road_group_id = ?", false, roadGroupId).Scan(&aadtGrowthRate).Error; err != nil {
		return aadtGrowthRate, err
	}
	return aadtGrowthRate, nil
}

func (sr *settingRepository) UpdateAadtGrowthRate(aadtGrowthRate *models.AadtGrowthRate) error {
	if err := sr.conn.Table("setting_aadt_growth_rate").Save(&aadtGrowthRate).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetAadtGrowthRate() ([]models.GetAadtGrowthRate, error) {

	var getAadtGrowthRate []models.GetAadtGrowthRate
	query := sr.conn
	query = query.Table("road_group").
		Select("road_group.id as road_group_id, road_group.number, road_group.name as road_group_name, setting_aadt_growth_rate.r").
		Joins("LEFT JOIN setting_aadt_growth_rate on setting_aadt_growth_rate.road_group_id = road_group.id").
		//Joins("FULL OUTER JOIN setting_aadt_growth_rate on setting_aadt_growth_rate.road_group_id = road_group.id").
		Order("road_group.id ASC")

	query = query.Find(&getAadtGrowthRate)
	err := query.Error
	if err != nil {
		return nil, err
	}

	return getAadtGrowthRate, err

}

func (sr *settingRepository) GetAadtPercentageVehicleTypeByRoadGroupId(roadGroupId int) (models.AadtPercentageVehicleType, error) {
	var aadtPercentageVehicleType models.AadtPercentageVehicleType
	if err := sr.conn.Table("setting_aadt_percentage_vehicle_type").Where("is_deleted = ? AND road_group_id = ?", false, roadGroupId).Scan(&aadtPercentageVehicleType).Error; err != nil {
		return aadtPercentageVehicleType, err
	}
	return aadtPercentageVehicleType, nil
}

func (sr *settingRepository) GetRoadGroup() ([]models.RoadGroup, error) {
	var roadGroup []models.RoadGroup
	if err := sr.conn.Table("road_group").Order("id").Scan(&roadGroup).Error; err != nil {
		return roadGroup, err
	}
	return roadGroup, nil
}

func (sr *settingRepository) UpdateAadtPercentageVehicleType(addtGrowthRate *models.AadtPercentageVehicleTypeParams) error {
	if err := sr.conn.Table("setting_aadt_percentage_vehicle_type").Save(&addtGrowthRate).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) UpdateAadtPercentageVehicleTypeByIsLatestIsFalseAndRoadGroupId(id int) error {
	if err := sr.conn.Table("setting_aadt_percentage_vehicle_type").Where("road_group_id = ?", id).Update("is_latest", false).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetAllAadtPercentageVehicleType() ([]models.AadtPercentageVehicleTypeParams, error) {
	var aadtPercentageVehicleTypeParams []models.AadtPercentageVehicleTypeParams
	query := sr.conn
	query = query.Where("is_latest = ?", true).Find(&aadtPercentageVehicleTypeParams)
	err := query.Error
	if err != nil {
		return []models.AadtPercentageVehicleTypeParams{}, err
	}

	return aadtPercentageVehicleTypeParams, err
}

func (sr *settingRepository) GetAadtPercentageVehicleTypeWithRoadGroupByRoadGroupId(roadGroupId int) (models.AadtPercentageVehicleTypeParams, error) {
	var aadtPercentageVehicleTypeParams models.AadtPercentageVehicleTypeParams
	fmt.Println(roadGroupId)
	query := sr.conn
	query = query.Where("is_latest = ? And road_group_id = ?", true, roadGroupId).Find(&aadtPercentageVehicleTypeParams)
	err := query.Error
	if err != nil {
		return models.AadtPercentageVehicleTypeParams{}, err
	}

	return aadtPercentageVehicleTypeParams, err
}

func (sr *settingRepository) UpdateAadtParameter(aadtParameter *models.AadtParameter) error {
	if err := sr.conn.Table("setting_aadt_parameter").Save(&aadtParameter).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetAadtParameterByRoadGroupId(roadGroupId int) (models.AadtParameter, error) {
	var aadtParameter models.AadtParameter
	query := sr.conn
	query = query.Select("setting_aadt_parameter.id as id, road_group.id as road_group_id, setting_aadt_parameter.elane, setting_aadt_parameter.four_wheel_axle_number, setting_aadt_parameter.four_wheel_vehicle_volume, setting_aadt_parameter.six_wheel_axle_number_id, setting_aadt_parameter.six_wheel_vehicle_volume, setting_aadt_parameter.six_wheel_percentage_truck, setting_aadt_parameter.six_wheel_factor_result, setting_aadt_parameter.ten_wheel_axle_number_id, setting_aadt_parameter.ten_wheel_vehicle_volume, setting_aadt_parameter.ten_wheel_percentage_truck, setting_aadt_parameter.ten_wheel_factor_result, setting_aadt_parameter.is_truck_factor, setting_aadt_parameter.speed_average, setting_aadt_parameter.speed_heavy_truck, setting_aadt_parameter.lane_distribution_factor , setting_aadt_parameter.directional_distribution_factor").
		Table("road_group").
		Joins("FULL OUTER JOIN setting_aadt_parameter on setting_aadt_parameter.road_group_id = road_group.id").
		Where("road_group.id = ?", roadGroupId)

	query = query.Find(&aadtParameter)
	err := query.Error
	if err != nil {
		return models.AadtParameter{}, err
	}

	return aadtParameter, err
}

func (r *settingRepository) GetVolumeByRoadGroupId(roadId int, status string) (models.VolumeAadt, error) {
	var volumeAadt models.VolumeAadt
	query := r.conn
	if err := query.Where("status = ?", status).Where("road_id = ?", roadId).Order("revision DESC").Order("year DESC").Order("id_parent DESC").First(&volumeAadt).Error; err != nil {
		fmt.Println(err)
		return volumeAadt, err
	}
	return volumeAadt, nil
}

func (r *settingRepository) GetParameterVehicleTypeById(id int) (models.RefAadtParameterVehicleType, error) {
	var refAadtParameterVehicleTypeList models.RefAadtParameterVehicleType
	query := r.conn
	query = query.Where("id != ?", id)
	err := query.Find(&refAadtParameterVehicleTypeList).Error
	if err != nil {
		return models.RefAadtParameterVehicleType{}, err
	}
	return refAadtParameterVehicleTypeList, nil
}

func (sr *settingRepository) GetAllAadtGrowthRate() ([]models.GetAadtGrowthRate, error) {

	var getAadtGrowthRate []models.GetAadtGrowthRate
	query := sr.conn
	query = query.Table("road_group").
		Select("road_group.id as road_group_id, road_group.number, road_group.name as road_group_name, setting_aadt_growth_rate.r").
		Joins("FULL OUTER JOIN setting_aadt_growth_rate on setting_aadt_growth_rate.road_group_id = road_group.id").
		Order("road_group.number DESC")

	query = query.Find(&getAadtGrowthRate)
	err := query.Error
	if err != nil {
		return nil, err
	}
	return getAadtGrowthRate, err
}

func (sr *settingRepository) GetAllAadtParameter() ([]models.GetAadtParameter, error) {
	var getAadtParameter []models.GetAadtParameter
	query := sr.conn
	query = query.Select("road_group.id as road_group_id, setting_aadt_parameter.elane, setting_aadt_parameter.four_wheel_axle_number, setting_aadt_parameter.four_wheel_vehicle_volume, setting_aadt_parameter.six_wheel_axle_number_id, setting_aadt_parameter.six_wheel_vehicle_volume, setting_aadt_parameter.six_wheel_percentage_truck, setting_aadt_parameter.six_wheel_factor_result, setting_aadt_parameter.ten_wheel_axle_number_id, setting_aadt_parameter.ten_wheel_vehicle_volume, setting_aadt_parameter.ten_wheel_percentage_truck, setting_aadt_parameter.ten_wheel_factor_result, setting_aadt_parameter.is_truck_factor, setting_aadt_parameter.speed_average, setting_aadt_parameter.speed_heavy_truck, setting_aadt_parameter.lane_distribution_factor , setting_aadt_parameter.directional_distribution_factor").
		Table("road_group").
		Joins("FULL OUTER JOIN setting_aadt_parameter on setting_aadt_parameter.road_group_id = road_group.id")
	query = query.Find(&getAadtParameter)
	err := query.Error
	if err != nil {
		return []models.GetAadtParameter{}, err
	}

	return getAadtParameter, err
}

func (sr *settingRepository) UpdateAadtParams(getAadtParams models.AadtParams) error {
	if err := sr.conn.Table("setting_aadt_params").Save(&getAadtParams).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) UpdateAadtParamsByIsLatestIsFalse() error {
	if err := sr.conn.Table("setting_aadt_params").Where("is_latest = is_latest").Update("is_latest", false).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetRoadWorkEffect() (models.SettingRoadWorkEffect, error) {
	var settingRoadWorkEffect models.SettingRoadWorkEffect
	if err := sr.conn.Table("setting_road_work_effect").First(&settingRoadWorkEffect).Error; err != nil {
		return models.SettingRoadWorkEffect{}, err
	}
	return settingRoadWorkEffect, nil
}

func (sr *settingRepository) UpdateRoadWorkEffect(settingRoadWorkEffect *models.SettingRoadWorkEffect) error {
	if err := sr.conn.Table("setting_road_work_effect").Save(&settingRoadWorkEffect).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) UpdateRoadWorkEffectParamsByIsLatestIsFalse() error {
	if err := sr.conn.Table("setting_road_work_effect_params").Where("is_latest = true").Update("is_latest", false).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) CreateRoadWorkEffectParams(settingRoadWorkEffectParams *models.SettingRoadWorkEffectParams) error {
	if err := sr.conn.Table("setting_road_work_effect_params").Save(&settingRoadWorkEffectParams).Error; err != nil {
		return err
	}
	return nil
}

func (r *settingRepository) DeleteSettingRefSurfaceByID(ID int) error {
	query := r.conn.Begin()
	err := query.Delete(models.RefSurface{}, ID).Error
	if err != nil {
		query.Rollback()
		return err
	}
	query.Commit()
	return nil
}

func (sr *settingRepository) CreateRoadUserCostLossValueParams(costLossValue *models.SettingRoadUserCostLossValue) error {
	if err := sr.conn.Table("setting_road_user_cost_acc_loss_value_params").Save(&costLossValue).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetRoadUserCostLossValueParamsByLatest() (models.SettingRoadUserCostLossValue, error) {
	var costLossValue models.SettingRoadUserCostLossValue
	if err := sr.conn.Table("setting_road_user_cost_acc_loss_value_params").Where("is_latest = ?", true).First(&costLossValue).Error; err != nil {
		return models.SettingRoadUserCostLossValue{}, err
	}
	return costLossValue, nil
}

func (sr *settingRepository) UpdateRoadUserCostLossValueParamsByIsLatestIsFalse() error {
	if err := sr.conn.Table("setting_road_user_cost_acc_loss_value_params").Where("is_latest = ?", true).Update("is_latest", false).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetRoadUserCostAccChanceOfAccidentParamsByLatest() ([]models.SettingRoadUserCostChanceOfAccident, error) {
	var settingRoadUserCostChanceOfAccident []models.SettingRoadUserCostChanceOfAccident
	if err := sr.conn.Table("setting_road_user_cost_acc_chance_of_accident_params").Where("is_latest = ?", true).Find(&settingRoadUserCostChanceOfAccident).Error; err != nil {
		return []models.SettingRoadUserCostChanceOfAccident{}, err
	}
	return settingRoadUserCostChanceOfAccident, nil
}

func (sr *settingRepository) GetRoadUserCostAccChanceOfAccidentParamsByLatestAndRoadGroupId(id int) (models.SettingRoadUserCostChanceOfAccident, error) {
	var chanceOfAccident models.SettingRoadUserCostChanceOfAccident
	if err := sr.conn.Table("setting_road_user_cost_acc_chance_of_accident_params").Where("is_latest = ? AND road_group_id = ?", true, id).First(&chanceOfAccident).Error; err != nil {
		return models.SettingRoadUserCostChanceOfAccident{}, err
	}
	return chanceOfAccident, nil
}

func (sr *settingRepository) CreateChanceOfAccidentParams(chanceOfAccident *models.SettingRoadUserCostChanceOfAccident) error {
	if err := sr.conn.Table("setting_road_user_cost_acc_chance_of_accident_params").Save(&chanceOfAccident).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) UpdateRoadUserCostLossValueParamsByIsLatestIsFalseAndRoadGroupId(id int) error {
	if err := sr.conn.Table("setting_road_user_cost_acc_chance_of_accident_params").Where("road_group_id = ?", id).Update("is_latest", false).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) CreateRoadUserCostRuc(params models.SettingRoadUserCost) error {
	if err := sr.conn.Table("setting_road_user_cost_ruc").Save(&params).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetRoadUserCostRuc() (models.SettingRoadUserCost, error) {
	var roadUserCost models.SettingRoadUserCost
	if err := sr.conn.Table("setting_road_user_cost_ruc").Where("is_latest = ?", true).First(&roadUserCost).Error; err != nil {
		return models.SettingRoadUserCost{}, err
	}
	return roadUserCost, nil
}

func (sr *settingRepository) UpdateRoadUserCostParamsByIsLatestIsFalse() error {
	if err := sr.conn.Table("setting_road_user_cost_ruc_params").Where("is_latest =  true").Update("is_latest", false).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) CreateRoadUserCostParams(roadUserCostParams *models.SettingRoadUserCostParams) error {
	if err := sr.conn.Table("setting_road_user_cost_ruc_params").Save(&roadUserCostParams).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) UpdateOptimizationByIsLatestIsFalse() error {
	if err := sr.conn.Table("setting_optimization_params").Where("is_latest =  true").Update("is_latest", false).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) CreateOptimizationParams(settingOptimization *models.SettingOptimization) error {
	if err := sr.conn.Table("setting_optimization_params").Save(&settingOptimization).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) GetOptimizationParams() (models.SettingOptimization, error) {
	var settingOptimization models.SettingOptimization
	if err := sr.conn.Table("setting_optimization_params").Where("is_latest = ?", true).First(&settingOptimization).Error; err != nil {
		return models.SettingOptimization{}, err
	}
	return settingOptimization, nil
}

func (sr *settingRepository) GetDeterioration() (models.SettingDeterioration, error) {
	var settingDeterioration models.SettingDeterioration
	if err := sr.conn.Table("setting_deterioration").First(&settingDeterioration).Error; err != nil {
		return models.SettingDeterioration{}, err
	}
	return settingDeterioration, nil
}

func (sr *settingRepository) GetDeteriorationList() ([]models.SettingDeterioration, error) {
	var settingDeteriorationList []models.SettingDeterioration
	if err := sr.conn.Table("setting_deterioration").Find(&settingDeteriorationList).Error; err != nil {
		return []models.SettingDeterioration{}, err
	}
	return settingDeteriorationList, nil
}

func (sr *settingRepository) GetDeteriorationByRoadGroupId(roadGroupId int) (models.SettingDeterioration, error) {
	var settingDeterioration models.SettingDeterioration
	if err := sr.conn.Table("setting_deterioration").Where("road_group_id = ?", roadGroupId).First(&settingDeterioration).Error; err != nil {
		return models.SettingDeterioration{}, err
	}
	return settingDeterioration, nil
}

func (sr *settingRepository) UpdateDeterioration(settingDeterioration *models.SettingDeterioration) error {
	if err := sr.conn.Table("setting_deterioration").Save(&settingDeterioration).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) UpdateDeteriorationParamsByIsLatestIsFalse() error {
	if err := sr.conn.Table("setting_deterioration_params").Where("is_latest = true").Update("is_latest", false).Error; err != nil {
		return err
	}
	return nil
}

func (sr *settingRepository) CreateDeteriorationParams(settingDeteriorationParams *models.SettingDeteriorationParams) error {
	if err := sr.conn.Table("setting_deterioration_params").Save(&settingDeteriorationParams).Error; err != nil {
		return err
	}
	return nil
}

func (r *settingRepository) GetDataList(model interface{}, where string) error {
	query := r.conn
	if where != "" {
		query = query.Where(where)
	}
	err := query.Find(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *settingRepository) GetHsms() (interface{}, error) {
	//  select * from hsms_01_bridge where road_id == null
	//  select * from hsms_01_guard where road_id == null
	//  select * from hsms_01_interchange where road_id == null
	//  select * from hsms_01_intersection where road_id == null
	//  select * from hsms_01_light where road_id == null
	//  select * from hsms_01_railwaycrossing where road_id == null
	//  select * from hsms_01_signal where road_id == null
	//  select * from hsms_01_uturnbridge where road_id == null
	return "", nil
}

func (r *settingRepository) GetHris() ([]models.RefHris, error) {
	var refHris []models.RefHris
	query := r.conn
	query = query.Where("is_deleted = false")
	query = query.Order("id DESC")
	err := query.Find(&refHris).Error
	if err != nil {
		return refHris, err
	}
	return refHris, nil
}

func (r *settingRepository) GetHrisByStatus() ([]models.RefHris, error) {
	var refHris []models.RefHris
	query := r.conn
	query = query.Where("status = true")
	query = query.Where("is_deleted = false")
	query = query.Order("id DESC")
	err := query.Find(&refHris).Error
	if err != nil {
		return refHris, err
	}
	return refHris, nil
}

func (r *settingRepository) GetHrisByStatusBylatest() (models.RefHris, error) {
	var refHris models.RefHris
	query := r.conn
	query = query.Where("is_deleted = false")
	query = query.Order("id DESC")
	err := query.First(&refHris).Error
	if err != nil {
		return refHris, err
	}
	return refHris, nil
}

func (r *settingRepository) CreateHris(refHris models.RefHris) error {
	query := r.conn
	err := query.Create(&refHris).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *settingRepository) UpdateHris(refHris models.RefHris) error {
	query := r.conn
	err := query.Save(&refHris).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *settingRepository) GetHrisById(id int) (models.RefHris, error) {
	var refHris models.RefHris
	query := r.conn
	query = query.Where("id = ?", id)
	query = query.Where("is_deleted = false")
	query = query.Order("id DESC")
	err := query.Find(&refHris).Error
	if err != nil {
		return refHris, err
	}
	return refHris, nil
}

func (r *settingRepository) GetSectionGeomWithFilter(filter primitive.D) ([]models.Item, error) {
	mongoDb := r.connMon
	collection := mongoDb.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_HRIS2_SECTION_GEOM_COLLECTION"))

	var sectionGeom []models.Item
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return sectionGeom, err
	}

	if err = cursor.All(context.TODO(), &sectionGeom); err != nil {
		return sectionGeom, err
	}

	return sectionGeom, nil
}

func (r *settingRepository) GetRoadLatest(filter primitive.D) ([]models.RoadLatest, error) {
	mongoDb := r.connMon
	collection := mongoDb.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_HRIS2_ROAD_LATEST_COLLECTION"))

	var roadLatest []models.RoadLatest
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return roadLatest, err
	}

	if err = cursor.All(context.TODO(), &roadLatest); err != nil {
		return roadLatest, err
	}

	return roadLatest, nil
}

func (r *settingRepository) InsertSectionGeom(data []interface{}) error {
	mongoDb := r.connMon
	collection := mongoDb.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_HRIS2_SECTION_GEOM_COLLECTION"))

	filter := bson.D{{"is_latested", true}}
	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}

	_, err = collection.InsertMany(context.TODO(), data)
	if err != nil {
		return err
	}

	return nil
}

func (r *settingRepository) InsertRoadLatest(data []interface{}) error {
	mongoDb := r.connMon
	collection := mongoDb.Database(os.Getenv("MONGODB_DB")).Collection(os.Getenv("MONGODB_HRIS2_ROAD_LATEST_COLLECTION"))

	filter := bson.D{{"is_latested", true}}
	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}

	_, err = collection.InsertMany(context.TODO(), data)
	if err != nil {
		return err
	}

	return nil
}

func (r *settingRepository) InsertRoadSection(roadSection models.InsertRoadSection, conn *gorm.DB) error {
	if err := conn.Save(&roadSection).Error; err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *settingRepository) InsertRoadGroup(roadGroup models.InsertRoadGroup, conn *gorm.DB) error {
	if err := conn.Save(&roadGroup).Error; err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *settingRepository) UpdateRoadGroupNameByNumber(number, name, shortName string, conn *gorm.DB) error {
	err := conn.Model(&models.InsertRoadGroup{}).Where("number = ?", number).Update("name", name).Update("short_name", shortName).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *settingRepository) GetRefHirsWithXt(conn *gorm.DB) ([]models.RefHris, error) {
	var refHris []models.RefHris
	err := conn.Where("status = true AND is_deleted = false").Find(&refHris).Error
	if err != nil {
		conn.Rollback()
		return refHris, err
	}
	return refHris, nil
}

func (r *settingRepository) GetRoadGroupWithXt(conn *gorm.DB) ([]models.RoadGroup, error) {
	var roadGroup []models.RoadGroup
	err := conn.Find(&roadGroup).Error
	if err != nil {
		conn.Rollback()
		return roadGroup, err
	}
	return roadGroup, nil
}

func (r *settingRepository) GetRoadSectionWithXt(conn *gorm.DB) ([]models.RoadSection, error) {
	var roadSection []models.RoadSection
	err := conn.Find(&roadSection).Error
	if err != nil {
		conn.Rollback()
		return roadSection, err
	}
	return roadSection, nil
}

func (r *settingRepository) InsertRoadGroupWithXt(roadGroup models.InsertRoadGroup, conn *gorm.DB) error {
	if err := conn.Save(&roadGroup).Error; err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *settingRepository) InsertRoadSectionWithXt(roadSection models.InsertRoadSection, conn *gorm.DB) error {
	if err := conn.Save(&roadSection).Error; err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *settingRepository) UpdateRoadGroupNameByNumberWithXt(number, name, shortName string, conn *gorm.DB) error {
	err := conn.Model(&models.InsertRoadGroup{}).Where("number = ?", number).Update("name", name).Update("short_name", shortName).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	return nil
}

func (r *settingRepository) GetHsmsBridge() ([]models.Hsms01Bridge, error) {
	query := r.conn
	var hsms01Bridge []models.Hsms01Bridge

	if err := query.Where("road_id is NULL").Find(&hsms01Bridge).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return hsms01Bridge, nil
		}
		return hsms01Bridge, err
	}

	return hsms01Bridge, nil
}

func (r *settingRepository) GetHsmsGuard() ([]models.Hsms01Guard, error) {
	query := r.conn
	var hsms01Guard []models.Hsms01Guard

	if err := query.Where("road_id is NULL").Find(&hsms01Guard).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return hsms01Guard, nil
		}
		return hsms01Guard, err
	}

	return hsms01Guard, nil
}

func (r *settingRepository) GetHsmsInterchange() ([]models.Hsms01Interchange, error) {
	query := r.conn
	var hsms01Interchange []models.Hsms01Interchange

	if err := query.Where("road_id is NULL").Find(&hsms01Interchange).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return hsms01Interchange, nil
		}
		return hsms01Interchange, err
	}

	return hsms01Interchange, nil
}

func (r *settingRepository) GetHsmsIntersection() ([]models.Hsms01Intersection, error) {
	query := r.conn
	var hsms01Intersection []models.Hsms01Intersection

	if err := query.Where("road_id is NULL").Find(&hsms01Intersection).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return hsms01Intersection, nil
		}
		return hsms01Intersection, err
	}

	return hsms01Intersection, nil
}

func (r *settingRepository) GetHsmsStreetlight() ([]models.Hsms01Light, error) {
	query := r.conn
	var hsms01Light []models.Hsms01Light

	if err := query.Where("road_id is NULL").Find(&hsms01Light).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return hsms01Light, nil
		}
		return hsms01Light, err
	}

	return hsms01Light, nil
}

func (r *settingRepository) GetHsmsRailwaycrossing() ([]models.Hsms01Railwaycrossing, error) {
	query := r.conn
	var hsms01Railwaycrossing []models.Hsms01Railwaycrossing

	if err := query.Where("road_id is NULL").Find(&hsms01Railwaycrossing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return hsms01Railwaycrossing, nil
		}
		return hsms01Railwaycrossing, err
	}

	return hsms01Railwaycrossing, nil
}

func (r *settingRepository) GetHsmsTrafficlight() ([]models.Hsms01Signal, error) {
	query := r.conn
	var hsms01Signal []models.Hsms01Signal

	if err := query.Where("road_id is NULL").Find(&hsms01Signal).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return hsms01Signal, nil
		}
		return hsms01Signal, err
	}

	return hsms01Signal, nil
}

func (r *settingRepository) GetHsmsUturnbridge() ([]models.Hsms01Uturnbridge, error) {
	query := r.conn
	var hsms01Uturnbridge []models.Hsms01Uturnbridge

	if err := query.Where("road_id is NULL").Find(&hsms01Uturnbridge).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return hsms01Uturnbridge, nil
		}
		return hsms01Uturnbridge, err
	}

	return hsms01Uturnbridge, nil
}

func (r *settingRepository) GetRoadSection() ([]models.RoadSection, error) {
	query := r.conn
	var roadSection []models.RoadSection
	err := query.Find(&roadSection).Error
	if err != nil {
		return roadSection, err
	}

	return roadSection, nil
}

func (r *settingRepository) GetRefAssetTable() ([]models.RefAssetTable, error) {
	query := r.conn
	var assetTable []models.RefAssetTable
	query = query.Where("is_active = true")
	err := query.Find(&assetTable).Error

	if err != nil {
		return assetTable, err
	}

	return assetTable, nil
}

func (r *settingRepository) DeleteHsmsByTypeAndId(typeData, idInteger string) error {
	query := r.conn
	sql := "DELETE FROM " + typeData + " WHERE id = " + idInteger
	err := query.Exec(sql).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *settingRepository) GetTableFromStruct(value interface{}) string {
	query := r.conn
	stmt := &gorm.Statement{DB: query}
	stmt.Parse(&value)

	return stmt.Schema.Table
}
