package domains

import (
	"time"

	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
)

type UseCase interface {
	GetYears(typeName string) (interface{}, error)
	GetAsset(roadIDs []int, depotCodes []string, filter requests.Asset) (interface{}, int64, error)
	GetAssetDetail(roadIDs []int, depotCodes []string, filter requests.Asset) (interface{}, int64, error)
	GetAssetMap(roadIDs []string, assetIDs []string, depotCodes []string, filter requests.AssetMap) (interface{}, error)
	GetAssetMapDetailByID(ID, assetTableID int) (interface{}, error)
	GetRoadDashboard() (interface{}, error)

	GetDashboardCondition(roadIDs []string, depotCodes []string, filter requests.Condition) (interface{}, error)
	GetDashboardConditionMap(roadIDs []string, depotCodes []string, page, limit string, filter requests.ConditionMap) (interface{}, error)

	//AssetDashboard(roadIDs []int, assetType string) (interface{}, error)

	GetDataMart(userID int) (interface{}, error)
	GetDataMartCheck() (interface{}, error)

	GetSurfaceDashboard(roadIDs []int, depotCodes []string, filter requests.Asset) (interface{}, error)
	GetSurfaceDashboardMap(roadIDs []int, depotCodes []string, filter requests.Asset, display int) (interface{}, error)

	GetRoadConditionList(roadId int) (interface{}, error)
	GetRoadConditionDetails(conditionRangeType string, idParent int) (interface{}, error)
	GetRoadRetroReflectivityList(roadID int) (interface{}, error)
	GetRoadRetroReflectivityDetails(rangeType string, refStripeTypeIDs string, idParent int) (interface{}, error)
	ConditionGradeAnalysisMap(roadConditions []models.RoadConditionDashboard, conditionType string, conditionGrades []models.ParamsConditionPreload) ([]responses.DashboardConditionMap, error)
	RefactiveGradeAnalysisMap(roadRetroReflectivitys []models.RoadRetroReflectivityDashboard, refactiveGrades []models.ParamsRoadLinePreload) ([]responses.DashboardConditionMap, error)
}
type Repository interface {
	GetAssetMap() ([]models.TableResult, error)
	GetTableResult(roadIDs []string, assetIDs []string, depotCodes []string, filter requests.AssetMap) ([]models.TableResult, error)
	FindImgFilepath(assetID int) ([]string, error)
	CheckSignImgId(assetID int) (bool, error)
	GetRoadAssetSignDataByRoadIDs(t models.TableResult, roadIDs []string, filter requests.AssetMap) ([]models.RoadAssetSign, error)
	GetRoadAssetTheGeomCuster(buildQuery string, roadIDs []string, assetIDs []string, depotCodes []string, filter requests.AssetMap) ([]models.RoadAssetGeomCuster, error)

	GetRoadAssetDetailColumn(roadAssetID int, refAssetTableID int) ([]responses.RoadAssetDetailColumn, error)
	GetRefDataFromSelect(tableName string, ID int64) (string, error)
	GetPavementSurfaceByRoadGroupID(roadGroupID int) ([]models.PavementSurface, error)
	GetRoadDashboard() ([]responses.RoadGroupDashboard, error)
	GetVolumeAADTByRoadGroupID(roadGroupID int) ([]responses.VolumeAADTDashboard, error)
	GetRoadGroups() ([]models.RoadGroup, error)
	//GetAssetRawData(roadIDs []int, filter requests.Asset) ([]models.RoadAssetSummary, error)
	GetAsset(roadIDs []int, depotCodes []string, filter requests.Asset) ([]models.RefAssetDashboard, int64, error)
	GetDashboardMaxMinYear() (responses.DashboardYearMaxMin, error)
	GetDashboardYear(tableName string) ([]int, error)
	GetRefAssetTableColumns(refAssetTableID int) ([]responses.RoadAssetDetailColumn, error)

	GetRoadAssetDetails(ID int, data []responses.RoadAssetDetailColumn) ([]map[string]interface{}, error)

	GetRoadByIDs(roadIDs []int) ([]models.RoadById, error)
	GetRoadSurfaceByGroupRoadID(roadID int, laneNo int) ([]int, error)
	GetRoadSurfaceByRoadID(roadID int, laneNo int, grp int) ([]models.RoadSurfacePrepareData, error)
	GetMaintenanceHistoryByRoadID(roadID int, year int) (models.MaintainPreloadGetAll, error)
	GetRefSurfaceParam(id int) (models.RefSurfaceParam, error)
	GetRoadSurfaceFirstGrpByRoadID(roadId int, laneNo int) ([]models.RoadSurfacePrepareData, error)
	GetMaintenanceData(roadId int, laneNo int) ([]models.MaintenanceData2, error)
	GetInterventionCriteriaParamsById(ID int) (models.SettingInterventionCriteriaParams, error)
	GetRefCriteriaMethodByID(ID int) (models.RefCriteriaMethod, error)
	GetRoadInfoByID(roadID int) (models.RoadInfo, error)
	GetRoadDatebegin(roadID int) ([]models.RoadDatebegin, error)
	GetRoadByRoadID(roadID int) (models.Road, error)
	GetRoadGroupByID(ID int) (models.RoadGroup, error)
	GetRoadGeomLaneByID(roadID int, laneNo int) (models.RoadGeom, error)

	CreateDataMart(data []models.DataMart) ([]models.DataMart, error)
	GetDataMartInfo(roadIDs []int, depotCodes []string, filter requests.Asset) ([]models.SurfaceInfo, error)
	GetRoadSurfaceAll() (map[int][]models.RoadSurfaceData2, error)
	GetRoadIDAll() ([]int, error)
	GetRoadGeomByRoadIDLaneNo(roadID int, laneNo int) (models.RoadGeom, error)
	GetRoadConditionData(roadID int, laneNo int) ([]models.RoadConditionSurveyM, time.Time, error)
	UpdatePercentage(percent float64, userID int) error
	GetDataMartCheck() (interface{}, error)

	GetRoadSurfaceInfo(roadID []int) ([]models.SurfaceInfo, error)
	GetInitialSurfaceArray() ([]models.Surface, error)
	GetRoadConditionGradesByID(ownerID int, conditionType string) ([]models.ParamsConditionPreload, error)
	GetRoadLineGradesByID(ownerID int) ([]models.ParamsRoadLinePreload, error)
	GetRoadCondition100M(where string) ([]models.RoadConditionSurvey100M2, error)
	GetRoadCondition(where string) ([]models.RoadConditionSurveyM2, error)

	GetRoadConditionDashboard(roadIDs []string, depotCodes []string, filter requests.Condition) ([]models.RoadConditionDashboard, error)
	GetRefGrade() ([]models.RefGrade, error)

	GetRoadByID(roadID int) (models.RoadInfoGeomDirection, error)
	GetAllRoadConditionByIdParent(idParent int) (models.RoadConditionAll, error)
	GetUserDepartmentById(userId int) (models.UserDepartment, error)
	GetRoadConditionList(rcStatus string, roadId int) ([]models.RoadConditionList, error)
	GetRoadRetroReflectivityList(roadID int) ([]models.RetroReflectivityList, error)
	GetRoadRetroReflectivityDetailsByIdParent(idParent int, refStripeTypeIDs []string) (models.RoadRetroReflectivityPreload, error)
	GetRoadRetroReflectivityDashboard(roadIDs []string, depotCodes []string, filter requests.Condition) ([]models.RoadRetroReflectivityDashboard, error)
	GetRoadConditionDashboardMap(roadIDs []string, depotCodes []string, filter requests.ConditionMap) ([]models.RoadConditionDashboard, error)
	GetRoadRetroReflectivityDashboardMap(roadIDs []string, depotCodes []string, filter requests.ConditionMap) ([]models.RoadRetroReflectivityDashboard, error)

	GetDataList(model interface{}, where string) error
}
