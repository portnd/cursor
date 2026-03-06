package domains

import (
	"mime/multipart"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SettingUseCase interface {
	GetAssetGroups(params requests.QueryParams) (interface{}, error)
	GetAssetGroupByID(id string) (models.RefAsset, error)
	DeleteAssetGroupByID(id string) error
	CreateAssetGroup(assetGroupName string) error
	UpdateAssetGroupByID(id, assetGroupName string) error

	GetDepartments(params requests.QueryParams) (interface{}, error)
	GetDepartmentByID(id string) (models.RefDepartment, error)
	DeleteDepartmentByID(id string) error
	CreateDepartment(departmentName string) error
	UpdateDepartmentByID(id, departmentName string) error

	GetOwners(params requests.QueryParams) (interface{}, error)
	GetOwnerByID(id int) (models.RefOwner, error)
	CreateOwner(request requests.OwnerRequest) (int, error)
	UpdateOwnerByID(int, requests.OwnerRequest) error
	DeleteOwnerByID(id string) error

	GetConditionList(int) ([]models.ParamsConditionPreload, error)
	CreateConditionList(int, requests.OwnerRequest) (interface{}, error)
	UpdateConditionList(int, requests.OwnerRequest) error

	GetSigns(c *gin.Context, params requests.QueryParams) (interface{}, error)
	GetSignByID(c *gin.Context, id string) (models.RefAssetSignImage, error)
	DeleteSignByID(id string) error
	CreateSign(c *gin.Context, requet requests.SignImageRequest) error
	UpdateSignByID(c *gin.Context, request requests.SignImageRequest) error

	GetAssetTables(params requests.AssetTableQueryParams) (helpers.ReturnValueOfGetAssetTables, error)
	GetAssetTableByID(c *gin.Context, id string) (helpers.ReturnValueOfGetAssetTableById, error)
	CreateAssetTable(request requests.AssetTableData, icon *multipart.FileHeader, c *gin.Context) error
	UpdateAssetTableByID(request requests.AssetTableData, icon *multipart.FileHeader, iconFilePathStatus string, c *gin.Context) error
	DeleteAssetTableByID(id string) error

	CreateInterventionCriteria(params requests.InterventionCriteria, c *gin.Context) (requests.InterventionCriteria, error)
	GetInterventionCriteria(c *gin.Context) (responses.InterventionCriteria, error)
	GetInterventionCriteriaMethod(c *gin.Context) (responses.InterventionCriteriaSequenceCriteriaMethod, error)
	GetInterventionCriteriaSequence(c *gin.Context) (responses.InterventionCriteriaSequenceCriteriaMethod, error)
	CreateInterventionCriteriaSequence(params requests.InterventionCriteriaSequenceCriteriaMethod, c *gin.Context) (requests.InterventionCriteriaSequenceCriteriaMethod, error)
	DeleteInterventionCriteriaById(id string, c *gin.Context) error

	PostRefSurface(request requests.RefSurface, uID int, isC2Nil bool) error
	PutRefSurface(request requests.RefSurface, uID int, id int, isC2Nil bool) error
	GetRefSurface(condition string) ([]models.NewRefSurface, error)
	GetRefSurfaceByID(id int) (interface{}, error)
	GetParamRefSurface(id int) (interface{}, error)
	DeleteSettingRefSurfaceByID(ID int) error

	CreateRoadWorkEffectAsphalt(params requests.SettingRoadWorkEffectAsphalt, c *gin.Context) (responses.SettingRoadWorkEffectAsphalt, error)
	GetRoadWorkEffectAsphalt(c *gin.Context) (responses.SettingRoadWorkEffectAsphalt, error)
	CreateRoadWorkEffectConcrete(params requests.SettingRoadWorkEffectConcrete, c *gin.Context) (responses.SettingRoadWorkEffectConcrete, error)
	GetRoadWorkEffectConcrete(c *gin.Context) (responses.SettingRoadWorkEffectConcrete, error)
	MergeRoadWorkEffectToParams(userID int) error

	CreateBudget(params requests.Budget, c *gin.Context) (requests.Budget, error)
	GetBudget(params requests.QueryParams, c *gin.Context) (interface{}, error)
	GetBudgetById(id string, c *gin.Context) (responses.Budget, error)
	UpdateBudget(params requests.UpdateBudget, c *gin.Context) (requests.UpdateBudget, error)
	DeleteBudgetById(id string, c *gin.Context) (interface{}, error)

	CreateAadtGrowthRate(params []requests.CreateAadtGrowthRate, c *gin.Context) ([]requests.CreateAadtGrowthRate, error)
	GetAadtGrowthRate(c *gin.Context) ([]responses.AadtGrowthRate, error)

	CreateAadtPercentageVehicleType(params requests.AadtPercentageVehicleType, c *gin.Context) (requests.AadtPercentageVehicleType, error)
	GetAadtPercentageVehicleType(c *gin.Context, id string) (responses.AadtPercentageVehicleType, error)

	CreateAadtParameter(params requests.CreateAadtParameter, c *gin.Context) (models.AadtParameter, error)
	GetAadtParameter(c *gin.Context, id string) (responses.AadtParameter, error)
	GetAadtParameterRoadGroupWithVolumeAadt(c *gin.Context) ([]responses.RoadGroupWithVolumeAadt, error)

	MergeSettingAadtToParams(userID int) (interface{}, error)

	CreateRoadUserCostAccLossValue(params requests.RoadUserCostAccLossValue, c *gin.Context) (responses.RoadUserCostAccLossValue, error)
	GetRoadUserCostAccLossValue(c *gin.Context) (responses.RoadUserCostAccLossValue, error)

	CreateRoadUserCostAccChanceOfAccident(params requests.RoadUserCostAccChanceOfAccident, c *gin.Context) (requests.RoadUserCostAccChanceOfAccident, error)
	GetRoadUserCostAccChanceOfAccident(id string, c *gin.Context) (responses.RoadUserCostAccChanceOfAccident, error)

	CreateRoadUserCostRucDefaultData(params requests.RoadUserCostRusDefaultData, c *gin.Context) (requests.RoadUserCostRusDefaultData, error)
	GetRoadUserCostRucDefaultData(c *gin.Context) (responses.RoadUserCostRusDefaultData, error)
	CreateRoadUserCostRucDriving(params requests.RoadUserCostRusDriving, c *gin.Context) (requests.RoadUserCostRusDriving, error)
	GetRoadUserCostRucDriving(c *gin.Context) (responses.RoadUserCostRusDriving, error)
	CreateRoadUserCostRucEngineSpeed(params requests.RoadUserCostRusEngineSpeed, c *gin.Context) (requests.RoadUserCostRusEngineSpeed, error)
	GetRoadUserCostRucEngineSpeed(c *gin.Context) (responses.RoadUserCostRusEngineSpeed, error)
	CreateRoadUserCostRucFuelConsumption(params requests.RoadUserCostRusFuelConsumption, c *gin.Context) (requests.RoadUserCostRusFuelConsumption, error)
	GetRoadUserCostRucFuelConsumption(c *gin.Context) (responses.RoadUserCostRusFuelConsumption, error)
	CreateRoadUserCostRucLubricantConsumption(params requests.RoadUserCostRusLubricantConsumption, c *gin.Context) (requests.RoadUserCostRusLubricantConsumption, error)
	GetRoadUserCostRucLubricantConsumption(c *gin.Context) (responses.RoadUserCostRusLubricantConsumption, error)
	CreateRoadUserCostRucWasteOfConsumption(params requests.RoadUserCostRusWasteOfConsumption, c *gin.Context) (requests.RoadUserCostRusWasteOfConsumption, error)
	GetRoadUserCostRucWasteOfConsumption(c *gin.Context) (responses.RoadUserCostRusWasteOfConsumption, error)
	CreateRoadUserCostRucMaintenance(params requests.RoadUserCostRusMaintenance, c *gin.Context) (requests.RoadUserCostRusMaintenance, error)
	GetRoadUserCostRucMaintenance(c *gin.Context) (responses.RoadUserCostRusMaintenance, error)
	CreateRoadUserCostRucTravelTime(params requests.RoadUserCostRusTravelTime, c *gin.Context) (requests.RoadUserCostRusTravelTime, error)
	GetRoadUserCostRucTravelTime(c *gin.Context) (responses.RoadUserCostRusTravelTime, error)
	CreateRoadUserCostRucVehicleSpeedCalculation(params requests.RoadUserCostRusVehicleSpeedCalculation, c *gin.Context) (requests.RoadUserCostRusVehicleSpeedCalculation, error)
	GetRoadUserCostRucVehicleSpeedCalculation(c *gin.Context) (responses.RoadUserCostRusVehicleSpeedCalculation, error)
	CreateRoadUserCostRucTrafficData(params requests.RoadUserCostRusTrafficData, c *gin.Context) (requests.RoadUserCostRusTrafficData, error)
	GetRoadUserCostRucTrafficData(c *gin.Context) (responses.RoadUserCostRusTrafficData, error)

	CreateOptimization(params requests.Optimization, c *gin.Context) (requests.Optimization, error)
	GetOptimization(c *gin.Context) (responses.Optimization, error)

	CreateDeteriorationAsphalt(params requests.DeteriorationAsphalt, c *gin.Context) (requests.DeteriorationAsphalt, error)
	GetDeteriorationAsphalt(roadGroupId string, c *gin.Context) (responses.DeteriorationAsphalt, error)
	CreateDeteriorationConcrete(params requests.DeteriorationConcrete, c *gin.Context) (requests.DeteriorationConcrete, error)
	GetDeteriorationConcrete(roadGroupId string, c *gin.Context) (responses.DeteriorationConcrete, error)

	GetRoadLineList(ownerId int) (responses.RoadLineList, error)
	CreateRoadLineList(ownerId int, requests requests.OwnerRoadLineRequest) (interface{}, error)
	UpdateRoadLineList(ownerId int, requests requests.OwnerRoadLineRequest) error

	GetOwnersRoadLine(params requests.QueryParamsReflectivityRange) (interface{}, error)
	GetOwnerRoadLineByID(ownerId int) (models.RefOwnerRoadLine, error)
	CreateOwnerRoadLine(request requests.OwnerRoadLineRequest) (int, error)
	UpdateOwnerRoadLineByID(ownerId int, request requests.OwnerRoadLineRequest) error
	DeleteOwnerRoadLineByID(id string) error

	GetHris() (interface{}, error)
	GetHrisById(id string) (interface{}, error)
	GetHrisPreview() (interface{}, error)
	CreateHris(request requests.CreateRefHris, userId int) (interface{}, error)
	UpdateHris(request requests.UpdateRefHris, id string, userId int) (interface{}, error)
	DeleteHris(id string, userId int) (interface{}, error)
	ImportHris() (interface{}, error)

	GetAllHsms(requests.FilterHsms) ([]responses.HsmsAll, error)
	DeleteHsmsByTypeAndId(typeData, id string) error
}

type SettingRepository interface {
	GetAll(records interface{}) error
	GetByID(record interface{}, id int) error
	UpdateByID(tableName, updateName string, id int) error
	DeleteByID(tableName string, id int) error

	CountAssetGroups(searchName string) (int64, error)
	GetAssetGroups(limit, offset int64, searchName string) ([]models.RefAsset, error)
	CreateAssetGroup(assetGroup models.RefAsset) error

	CountDepartments(searchName string) (int64, error)
	GetDepartments(limit, offset int64, searchName string) ([]models.RefDepartment, error)
	CreateDepartment(department models.RefDepartment) error
	GetOwnersRoadLine(limit, offset int64, params requests.QueryParamsReflectivityRange) ([]models.RefOwnerRoadLinePreload, error)
	CountOwnersRoadLine(params requests.QueryParamsReflectivityRange) (int64, error)
	GetOwnerRoadLineByID(id int) (models.RefOwnerRoadLine, error)
	GetOwnersRoadLineAll() ([]models.RefOwnerRoadLine, error)
	CreateOwnerRoadLine(owner *models.RefOwnerRoadLine) error
	UpdateOwnerRoadLineByID(id int, owner models.RefOwnerRoadLine) error
	DeleteOwnerRoadLineByID(ownerId int) error

	CountOwners(params requests.QueryParams) (int64, error)
	GetOwners(limit, offset int64, params requests.QueryParams) ([]models.RefOwnerPreload, error)
	GetOwnerByID(id int) (models.RefOwner, error)
	GetOwnersAll() ([]models.RefOwner, error)
	CreateOwner(*models.RefOwner) error
	UpdateOwnerByID(id int, owner models.RefOwner) error
	DeleteOwnerByID(id int) error

	GetParamsCondition(ownerId int) ([]models.ParamsConditionPreload, error)
	GetOwnerByRoadID(roadId int) (models.RoadOwner, error)
	// CreateCondition(int, []models.ParamsCondition) error
	DeleteCondition(ownerId int) error

	CountSigns(searchName string) (int64, error)
	GetSigns(limit, offset int64, searchName string) ([]models.RefAssetSignImage, error)
	CreateSignImage(record *models.RefAssetSignImage) error
	UpdateSignImage(id int, record models.RefAssetSignImage) error

	CountAssetTables(assetType, name string, grpId int) (int64, error)
	GetAssetTables(limit, offset int64, assetType, name string, grpId int) ([]models.AssetTable, error)
	GetAssetTableStaffs() ([]models.AssetTableStaff, error)
	GetAssetTableByID(id int) (models.AssetTable, error)
	GetAssetTableStaffByID(id int) ([]models.AssetTableStaff, error)
	CreateAssetTable(data helpers.CreateAssetTable) error
	InsertRefAssetTable(tx *gorm.DB, record *models.RefAssetTable) error
	InsertRefAssetTableColumns(tx *gorm.DB, record models.RefAssetTableColumns) error
	InsertRefAssetTableStaffs(tx *gorm.DB, record models.RefAssetTableStaff) error
	GetColumnsByID(id int) (models.RefAssetTableColumns, error)
	GetColumnMaxSeqByAssetTableID(id int) (int, error)
	UpdateAssetTable(data helpers.UpdateAssetTable) error
	GetAllAssetTables() ([]models.AssetTable, error)
	CountAssetTableColumns(tableName string) (int, error)
	DeleteAssetTable(id int) error

	GetGrade() ([]models.RefGrade, error)

	GetCriteriaMethodByName(name string) (models.RefCriteriaMethod, error)
	GetCriteriaMethodById(int int) (models.RefCriteriaMethod, error)
	GetCriteriaMethod() ([]models.RefCriteriaMethod, error)

	UpdateInterventionCriteria(interventionCriteria *models.InterventionCriteria) error
	DeleteInterventionCriteriaConditionById(id int) error
	UpdateInterventionCriteriaCondition(interventionCriteriaCondition *models.InterventionCriteriaCondition) error
	GetRefCriteriaMethod() ([]models.RefCriteriaMethod, error)
	GetInterventionCriteriaById(id int) (models.InterventionCriteria, error)
	GetInterventionCriteriaByNotId(id int) ([]models.InterventionCriteria, error)
	GetInterventionCriteriaConditionListByInterventionCriteriaId(id int) ([]models.InterventionCriteriaCondition, error)
	GetInterventionCriteriaConditionById(id int) (models.InterventionCriteriaCondition, error)
	GetInterventionCriteria() ([]models.InterventionCriteria, error)
	GetInterventionCriteriaConditionSequence() ([]models.InterventionCriteriaCondition, error)
	CreateInterventionCriteriaParams(settingInterventionCriteria *models.SettingInterventionCriteria) error
	UpdateInterventionCriteriaParamsByIsLatestIsFalse() error
	CountInterventionCriteriaByMaintenanceMethod(surface string) (models.InterventionCriteriaCount, error)

	GetLatestSurfaceParamsById(id int) (models.RefSurfaceParam, error)

	InsertRefSurface(tx *gorm.DB, data models.NewRefSurface) (int, error)
	UpdateRefSurface(tx *gorm.DB, data models.NewRefSurface) (bool, error)
	InsertRefSurfaceParam(tx *gorm.DB, dataToInsert models.RefSurfaceParam) error
	GetRefSurface(condition string) ([]models.NewRefSurface, error)
	GetRefSurfaceByID(id int) (models.NewRefSurface, error)
	GetNewRefSurfaceByID(id int) (models.RefSurfaceNew, error)
	GetParamRefSurfaceByID(id int) ([]models.RefSurfaceParam, error)
	DeleteSettingRefSurfaceByID(ID int) error

	CreateRoadWorkEffectParams(*models.SettingRoadWorkEffectParams) error
	UpdateRoadWorkEffectParamsByIsLatestIsFalse() error
	GetRoadWorkEffect() (models.SettingRoadWorkEffect, error)
	UpdateRoadWorkEffect(settingRoadWorkEffect *models.SettingRoadWorkEffect) error

	GetBudget() ([]models.SettingBudget, error)
	GetBudgetById(budgetId int) (models.SettingBudget, error)
	GetBudgetMethodByBudgetId(budgetId int) ([]models.SettingBudgetMethod, error)
	UpdateBudget(budget *models.SettingBudget) error
	UpdateBudgetMethod(budgetMethod *models.SettingBudgetMethod) error
	GetBudgetMethodById(budgetId int) (models.SettingBudgetMethod, error)
	GetBudgetMethodListByBudgetId(budgetId int) ([]models.SettingBudgetMethod, error)
	CountGetBudgetByName(name string) (int64, error)
	GetBudgetByName(limit, offset int64, name string) ([]responses.BudgetList, error)

	GetAadtGrowthRateByRoadGroupId(roadGroupId int) (models.AadtGrowthRate, error)
	UpdateAadtGrowthRate(*models.AadtGrowthRate) error
	GetAadtGrowthRate() ([]models.GetAadtGrowthRate, error)
	GetAllAadtGrowthRate() ([]models.GetAadtGrowthRate, error)

	GetAadtPercentageVehicleTypeByRoadGroupId(roadGroupId int) (models.AadtPercentageVehicleType, error)
	UpdateAadtPercentageVehicleType(*models.AadtPercentageVehicleTypeParams) error
	UpdateAadtPercentageVehicleTypeByIsLatestIsFalseAndRoadGroupId(id int) error
	GetAadtPercentageVehicleTypeWithRoadGroupByRoadGroupId(roadGroupId int) (models.AadtPercentageVehicleTypeParams, error)
	GetAllAadtPercentageVehicleType() ([]models.AadtPercentageVehicleTypeParams, error)

	UpdateAadtParameter(aadtParameter *models.AadtParameter) error
	GetAadtParameterByRoadGroupId(roadGroupId int) (models.AadtParameter, error)
	GetVolumeByRoadGroupId(roadGroupId int, status string) (models.VolumeAadt, error)
	GetParameterVehicleTypeById(id int) (models.RefAadtParameterVehicleType, error)
	GetAllAadtParameter() ([]models.GetAadtParameter, error)
	GetRoadGroup() ([]models.RoadGroup, error)

	UpdateAadtParams(models.AadtParams) error
	UpdateAadtParamsByIsLatestIsFalse() error

	CreateRoadUserCostLossValueParams(costLossValue *models.SettingRoadUserCostLossValue) error
	UpdateRoadUserCostLossValueParamsByIsLatestIsFalse() error
	GetRoadUserCostLossValueParamsByLatest() (models.SettingRoadUserCostLossValue, error)

	GetRoadUserCostAccChanceOfAccidentParamsByLatest() ([]models.SettingRoadUserCostChanceOfAccident, error)
	GetRoadUserCostAccChanceOfAccidentParamsByLatestAndRoadGroupId(id int) (models.SettingRoadUserCostChanceOfAccident, error)
	CreateChanceOfAccidentParams(chanceOfAccident *models.SettingRoadUserCostChanceOfAccident) error
	UpdateRoadUserCostLossValueParamsByIsLatestIsFalseAndRoadGroupId(id int) error

	CreateRoadUserCostRuc(params models.SettingRoadUserCost) error
	GetRoadUserCostRuc() (models.SettingRoadUserCost, error)
	CreateRoadUserCostParams(roadUserCostParams *models.SettingRoadUserCostParams) error
	UpdateRoadUserCostParamsByIsLatestIsFalse() error

	GetOptimizationParams() (models.SettingOptimization, error)
	CreateOptimizationParams(settingOptimization *models.SettingOptimization) error
	UpdateOptimizationByIsLatestIsFalse() error

	GetDeterioration() (models.SettingDeterioration, error)
	GetDeteriorationList() ([]models.SettingDeterioration, error)
	GetDeteriorationByRoadGroupId(roadGroupId int) (models.SettingDeterioration, error)
	CreateDeteriorationParams(settingDeteriorationParams *models.SettingDeteriorationParams) error
	UpdateDeteriorationParamsByIsLatestIsFalse() error
	UpdateDeterioration(settingDeterioration *models.SettingDeterioration) error

	GetDataList(model interface{}, where string) error
	DeleteRoadLine(ownerId int) error
	GetParamsRoadLine(ownerId int) ([]models.ParamsRoadLinePreload, error)
	CreateData(model interface{}) error
	StartTransSection() *gorm.DB
	RollBack(tx *gorm.DB) error
	Commit(tx *gorm.DB) error

	GetHris() ([]models.RefHris, error)
	GetHrisByStatus() ([]models.RefHris, error)
	GetHrisById(id int) (models.RefHris, error)
	GetSectionGeomWithFilter(filter primitive.D) ([]models.Item, error)
	GetRoadLatest(filter primitive.D) ([]models.RoadLatest, error)
	CreateHris(refHris models.RefHris) error
	GetHrisByStatusBylatest() (models.RefHris, error)
	UpdateHris(refHris models.RefHris) error
	InsertSectionGeom(data []interface{}) error
	InsertRoadLatest(data []interface{}) error
	InsertRoadGroup(roadGroup models.InsertRoadGroup, conn *gorm.DB) error
	UpdateRoadGroupNameByNumber(number, name, shortName string, conn *gorm.DB) error

	GetRefHirsWithXt(conn *gorm.DB) ([]models.RefHris, error)
	GetRoadSectionWithXt(conn *gorm.DB) ([]models.RoadSection, error)
	GetRoadGroupWithXt(conn *gorm.DB) ([]models.RoadGroup, error)
	InsertRoadGroupWithXt(roadGroup models.InsertRoadGroup, conn *gorm.DB) error
	InsertRoadSectionWithXt(roadSection models.InsertRoadSection, conn *gorm.DB) error
	UpdateRoadGroupNameByNumberWithXt(number, name, shortName string, conn *gorm.DB) error

	GetTableFromStruct(value interface{}) string
	GetRefAssetTable() ([]models.RefAssetTable, error)
	GetRoadSection() ([]models.RoadSection, error)
	GetHsmsBridge() ([]models.Hsms01Bridge, error)
	GetHsmsGuard() ([]models.Hsms01Guard, error)
	GetHsmsInterchange() ([]models.Hsms01Interchange, error)
	GetHsmsIntersection() ([]models.Hsms01Intersection, error)
	GetHsmsStreetlight() ([]models.Hsms01Light, error)
	GetHsmsRailwaycrossing() ([]models.Hsms01Railwaycrossing, error)
	GetHsmsTrafficlight() ([]models.Hsms01Signal, error)
	GetHsmsUturnbridge() ([]models.Hsms01Uturnbridge, error)

	DeleteHsmsByTypeAndId(typeData, idInteger string) error
}
