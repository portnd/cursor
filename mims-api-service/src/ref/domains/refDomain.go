package domains

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
)

type RefUseCase interface {
	GetRefAsset() (interface{}, error)
	GetRefAssetPosition() (interface{}, error)
	GetRefAssetArea() (interface{}, error)
	GetRefAssetGuardrail() (interface{}, error)
	GetRefAssetReflecType() (interface{}, error)
	GetRefAssetOwner() (interface{}, error)
	GetRefAssetLightType() (interface{}, error)
	GetRefAssetLightWatt() (interface{}, error)
	GetRefAssetCleranceType() (interface{}, error)
	GetRefAssetCrashcushionType() (interface{}, error)
	GetRefAssetFenceType() (interface{}, error)
	GetRefAssetKmstoneType() (interface{}, error)
	GetRefAssetQutterType() (interface{}, error)
	GetRefAssetTrafficCameraType() (interface{}, error)
	GetRefAssetWeightStationType() (interface{}, error)
	GetRefAssetBuildingType() (interface{}, error)
	GetRefAssetNoiseBarrier() (interface{}, error)

	GetRefAssetSign() (interface{}, error)
	GetRefAssetSignImage() (interface{}, error)
	GetRefAssetSignType() (interface{}, error)
	GetRefAssetTable() (interface{}, error)
	GetRefAssetTableColumns() (interface{}, error)
	GetRefAssetTableStaff() (interface{}, error)

	///////////////////////////////////////
	GetRefDataStatus() (interface{}, error)
	GetRefDepartment() (interface{}, error)
	GetRefDirection() (interface{}, error)
	GetRefGrade() (interface{}, error)
	GetRefMaterialBase() (interface{}, error)
	GetRefMaterialSubbase() (interface{}, error)
	GetRefMaterialSubgrade() (interface{}, error)
	GetRefOwner() (interface{}, error)
	GetRefRoadType() (interface{}, error)
	GetRefSurface() (interface{}, error)
	GetRefSurfaceType() (interface{}, error)
	GetRefSurfaceGroup() (interface{}, error)
	GetRefStructureSurface() (interface{}, error)
	GetRefTableList() ([]models.RefTableList, error)
	GetRefColorList() (interface{}, error)
	GetRefRoadTypeIcon() (interface{}, error)
	GetRoadConditionGrades() (interface{}, error)
	GetParameterVehicleTypeList(road_group_id string) ([]responses.RefAadtParameterVehicleType, error)
	GetRefCriteriaType() (interface{}, error)
	GetRoadUserCostAcc() (interface{}, error)
	GetRoadUserCostRuc() (interface{}, error)

	GetRefStripeColor() (interface{}, error)
	GetRefStripeType() (interface{}, error)

	GetMaintenanceAnalysisStrategicBudgetType() (interface{}, error)
	GetMaintenanceAnalysisStrategicTargetType() (interface{}, error)
	GetMaintenanceAnalysisStrategic() (interface{}, error)

	GetRoadLineConditionGrades() (interface{}, error)
	GetRefRoadConditionRange() (interface{}, error)

	GetRefDistrictsList() (interface{}, error)
	GetRoadGroupList() ([]responses.RoadGroupInitData, error)
	GetRoadSectionList(userID int) ([]models.RoadSectionInitData, error)
	GetRefDistrictsInitList(userID int) (interface{}, error)
	GetRefRoadTypeLevel(string) ([]models.RefRoadTypeInit, error)
	GetRefDivisionInitList(userID int) (interface{}, error)
	GetRefDivisionInitListDashboardAsset(userID int) (interface{}, error)
	GetRefDivisionInitListDashboardCondition(userID int) (interface{}, error)
	GetRefDivisionInitListDashboardSurface(userID int) (interface{}, error)
	GetRefDivisionInitListDashboardMaintenance(userID int) (interface{}, error)
	GetRefReflectivityRange() (interface{}, error)

	GetRefCriteriaMethod() ([]models.RefCriteriaMethod, error)

	GetRefUserOwner() ([]responses.RefUserOwner, error)
}

type RefRepository interface {
	GetRef(result interface{}) error
	GetRefAssetSignImage(result interface{}) error
	GetRefStatus(result interface{}) error
	GetRefTableList() ([]models.RefTableList, error)
	GetRefColorList() ([]models.RefColorList, error)
	GetRefRoadTypeIcon() ([]models.RefRoadTypeIcon, error)
	GetRoadConditionGrades() ([]models.ParamsConditionPreload, error)
	GetRefCriteriaType() ([]models.RefCriteriaType, error)
	GetParameterVehicleTypeListByRoadGroupId(road_group_id int) ([]models.RefAadtParameterVehicleType, error)
	GetRoadUserCostAcc() ([]models.RefRoadUserCostAcc, error)
	GetRoadUserCostRuc() ([]models.RefRoadUserCostRuc, error)
	GetDataList(model interface{}, where string) error
	GetMaintenanceAnalysisStrategicBudgetType() ([]models.RefMaintenanceAnalysisCondition, error)
	GetMaintenanceAnalysisStrategicTargetType() ([]models.RefMaintenanceAnalysisTarget, error)
	GetMaintenanceAnalysisStrategic() ([]models.MaintenanceAnalysisStrategicTypePreload, error)

	GetParamsCondition(ownerId int) ([]models.ParamsConditionPreload, error)
	GetParamsRoadLine(ownerId int) ([]models.ParamsRoadLinePreload, error)
	GetRefDistrictsList() (interface{}, error)
	GetRefDistrictsInitList(isAllData, isOwnerData bool, depotCode string) (interface{}, error)
	GetRoadGroupList() ([]models.RoadGroupInitData, error)
	GetRoadSectionList(isAllData, isOwnerData bool, depotCode string) ([]models.RoadSectionInitData, error)
	GetRefDivisionInitList(isAllData, isOwnerData bool, depotCode string) ([]models.RefDivisionList, error)

	GetRefCriteriaMethod() ([]models.RefCriteriaMethod, error)

	GetRefUserOwner() ([]responses.RefUserOwner, error)
}
