package domains

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
)

type UseCase interface {

	//////////////// NEW MIMS ////////////////
	ReportStatus() (interface{}, error)
	CheckReportStatus(id int) (interface{}, error)

	FilterAssetType1(userID int) (interface{}, error)
	FilterAssetType2(userID int) (interface{}, error)
	FilterAssetType3(userID int) (interface{}, error)

	FilterRoadType1(userID int) (interface{}, error)
	FilterRoadType2(userID int) (interface{}, error)
	FilterRoadType3(userID int) (interface{}, error)
	FilterRoadType4(userID int) (interface{}, error)
	FilterRoadType5(userID int) (interface{}, error)
	FilterRoadType6(userID int) (interface{}, error)

	FilterAadtType1(userID int) (interface{}, error)

	FiltertRoadDamageType1(userID int) (interface{}, error)
	FiltertRoadDamageType2(userID int) (interface{}, error)

	FilterMaintenanceKpiType1(userID int) (interface{}, error)

	FilterMaintenanceFilterType1(userID int) (interface{}, error)

	Report8(roadSectionId, filterCriteriaId, year int, typ string) (interface{}, error)
	Report9(roadSectionId, filterCriteriaId, year int, typ string) (interface{}, error)
	Report12(roadSectionId, year int, filterConditionName, typ string) (interface{}, error)

	//////////////// NEW MIMS ////////////////

	// road report
	GetReportRoad(roadGroupIDs []int, typ string) (interface{}, error)

	//รานงานปริมาณจราจร
	GetReportTrafficVolume(roadSectionIDs []string, typ string, year int) (interface{}, error)
	GetReport1(string, string, string) (interface{}, error)
	GetReport2(string, string, string) (interface{}, error)
	GetReport3(roadID, typ string) (interface{}, error)
	Report13(req requests.Report13) (interface{}, error)

	Report5(year, roadIDStr, typ string) (interface{}, error)
	Report7(factor, year, SectionIDstr, measureID, typ string) (interface{}, error)
	Report6(year, roadSectionID, typ, dis string) (interface{}, error)
	Report11(year, roadSectionID, typ string) (interface{}, error)
	Report10(year, roadSectionID, typ string) (interface{}, error)
}

type Repository interface {
	//////////////// NEW MIMS ////////////////
	CreateReportStatus() (int, error)
	UpdateReportStatus(id int, path string) error
	CheckReportStatusById(id int) (models.ReportStatus, error)

	FilterAssetRoad(isAllData, isOwnerData bool, depotCode string) ([]models.FilterAssetRoad, error)
	FilterAsset() ([]models.FilterAsset, error)
	FilterRoadGroup(isAllData, isOwnerData bool, depotCode string) ([]models.RoadGroup, error)
	FilterRoadSurface(isAllData, isOwnerData bool, depotCode string) ([]models.FilterRoadSurface, error)
	FilterRoadCondition(isAllData, isOwnerData bool, depotCode string) ([]models.FilterRoadCondition, error)
	FilterRefOwner() ([]models.RefOwner, error)
	FilterRoadRetroReflectivity(isAllData, isOwnerData bool, depotCode string) ([]models.FilterRoadRetroReflectivity, error)
	FilterRefOwnerLine() ([]models.RefOwnerRoadLine, error)

	FilterRoadDamage(isAllData, isOwnerData bool, depotCode string) ([]models.FilterRoadDamage, error)

	FilterMaintenance(isAllData, isOwnerData bool, depotCode string) ([]models.FilterMaintenance, error)

	FilterAadt(isAllData, isOwnerData bool, depotCode string) ([]models.FilterAadt, error)

	ReportRetroReflectivity(roadSectionId, filterCriteriaId, year int) (models.ReportRetroReflectivityRoadSection, error)
	ParamsRoadLine(filterCriteriaId int) ([]models.ParamsRoadLinePreload, error)
	RefStripeColor() ([]models.RefStripeColor, error)
	RefGrade() ([]models.RefGrade, error)
	RefStripeType() ([]models.RefStripeType, error)

	Report12RoadGroup(roadSectionId int) (models.ReportKpiRoadGroup, error)

	Report12Iri(roadSectionId, year int) (models.ResultIri, error)
	Report12Ifi(roadSectionId, year int) ([]models.ResultIfi100, error)
	Report12Rut(roadSectionId, year int) ([]models.ResultRut100, error)
	Report12G7(roadSectionId, year int) ([]models.ResultG7100, error)

	//////////////// NEW MIMS ////////////////

	// road report
	GetReportRoad(roadGroupIDs []int) ([]models.RoadListReport, error)

	//รายงานปริมาณจราจร

	GetRoadDetailsByRoadSectionID(roadSectionIDs []string) ([]models.ReportTrafficVolumeHeader, error)
	GetReportTrafficVolume(roadIDs []int, year int) ([]models.VolumeAadt, error)

	GetRoadInfo(roadID int) ([]models.RoadReportInfo, error)
	GetTableName(roadID, assetID string) (*models.TableName, error)
	GetColumn(assetID, typ string) ([]models.Column, error)
	GetAssetName(assetID string) (*models.AssetName, error)
	GetRow(columnName []string, roadID, assetID, tableName, join, typ string) ([]map[string]interface{}, error)
	GetRoadGeom(roadID string) ([]models.MapGeom, error)
	GetMapGeom(roadID, assetID, tableName, typ string) ([]models.MapGeom, error)

	GetRoadInfoForAssetMap(roadID int) ([]models.RoadReportInfo, error)
	GetTableNameForAssetMap(roadID, assetID string) (*models.TableName, error)
	GetColumnForAssetMap(assetID, typ string) ([]models.Column, error)
	GetAssetNameForAssetMap(assetID string) (*models.AssetName, error)
	GetRowForAssetMap(columnName []string, roadID, assetID, tableName, join, typ string) ([]map[string]interface{}, error)
	GetRoadGeomForAssetMap(roadID string) ([]models.MapGeom, error)
	GetMapGeomForAssetMap(roadID, assetID, tableName, typ string) ([]models.MapGeom, error)

	GetRoadInfoForAssetSummary(roadSectionID int) ([]models.RoadReportInfo, error)
	GetReportSummayAssetForAssetSummary(roadID []int) ([]models.RefSummaryAsset, error)
	GetNoTypeCountForAssetSummary(roadAssetID []int32, tableName string) (*models.CountSummaryAsset, error)
	GetTypeCountForAssetSummary(roadAssetID []int32, tableName string, joinName string) ([]models.CountAndTypeSummaryAsset, error)
	GetLightCountForAssetSummary(roadAssetID []int32) ([]models.CountLightSummaryAsset, error)

	GetReportMaintenance(yearStart, yearEnd int, roadSectionId int) ([]models.DataGetMaintenance, []models.MethodMaintenance, error)
	GetMultiRoadInfo(roadSectionId int) ([]models.MultiRoadInfo, error)

	GetInitialSurfaceArrayForRoadCondition() ([]models.Surface, error)
	GetDataMartInfoForRoadCondition(roadID int) ([]models.SurfaceInfo, error)
	GetRoadInfoForRoadCondition(roadID int) (*models.RoadReportInfo, error)
	GetSurfaceForRoadCondition() ([]models.Surface, error)
	GetRoadSectionByIDForRoadCondition(id int) (models.RoadSection, error)
	GetRoadGroupByIDForRoadCondition(id int) (models.RoadGroup, error)
	GetRoadFromSectionIDForRoadCondition(SectionID int) ([]int, error)
	GetReportConditionForRoadCondition(year, roadID int, dis string) (*models.DataReportCondition, []models.DataRoadCondition, []models.DataRoadConditionM, error)
	GetMeasureValueSummaryConditionForRoadCondition(factor, measureID string) ([]models.DataGradeSummaryCondition, error)
	GetRoadConditionGradesByIDForRoadCondition(ownerID int, conditionType string) ([]models.ParamsConditionPreload, error)
	GetRoadLineGradesByIDForRoadCondition(ownerID int) ([]models.ParamsRoadLinePreload, error)
	GetRoadConditionDashboardForRoadCondition(roadIDs []int, depotCodes []string, filter requests.Condition) ([]models.RoadConditionReport, error)
	GetRoadRetroReflectivityDashboardForRoadCondition(roadIDs []int, depotCodes []string, filter requests.Condition) ([]models.RoadRetroReflectivityDashboard, error)
	GetDataListForRoadCondition(model interface{}, where string) error
	GetRoadConditionForRoadCondition(where string) ([]models.RoadConditionSurveyM2, error)
	GetRoadCondition100MForRoadCondition(where string) ([]models.RoadConditionSurvey100M2, error)
	GetCenterFromRoadSectionIDForRoadCondition(roadSectionID int) string

	GetRoadInfoForRoadDamage(roadID int) (*models.RoadReportInfo, error)
	GetRoadFromSectionIDForRoadDamage(SectionID int) ([]int, error)
	GetRoadSectionByIDForRoadDamage(id int) (models.RoadSection, error)
	GetRoadGroupByIDForRoadDamage(id int) (models.RoadGroup, error)
	GetReportDamageForRoadDamage(year, roadID int) (*models.DataReportDamage, []models.DataRoadDamage, []models.PositionDamage, error)

	// ForRoadCondition

	// ForRoadDamage
}
