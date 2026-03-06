package domains

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gorm.io/gorm"
)

type Usecase interface {
	// GetMaintenance(interface{}, error)
	GetRoadMaintenanceYear(roadID int) (interface{}, error)
	GetMaintenanceByRoadID(roadID int, year int) (interface{}, error)

	//GetMaintenanceListHistory(maintenanceID int, prams requests.MaintenancePrams) (interface{}, error)
	//GetMaintenanceYearHistory(maintenanceID int) (interface{}, error)
	CalculateDistance(request requests.CalDistance) (float64, error)

	GetMaintenanceList(userID int, req requests.MaintenancePrams, limit, offset int64) ([]responses.MaintenanceList, int64, error)
	CreateMaintenance(request requests.MaintenanceReq, userID int, attReqs []requests.MaintenanceAttachmentsReq) (interface{}, error)
	UpdateMaintenance(idParent int, request requests.MaintenanceReq, userID int, attReqs []requests.MaintenanceAttachmentsReq) (*int, *int, error)
	DeleteMaintenance(idParent int) error

	GetMaintenanceListByID(maintenanceID int) (interface{}, error)
	GetMaintenanceHistoryByID(maintenanceID int, mRoadHisId int) (interface{}, error)

	//MaintenanceFinished(maintenanceID int, req requests.MaintenanceFinished) error
	CheckMaintenanceDuplicate(maintenanceID int, name string) (bool, error)
	CheckValidateIsMethod(IDParent int) (bool, error)
	GetMaintenanceYear(roadId int) ([]int, error)

	GetRefCriteriaMethod() (interface{}, error)

	CreateMaintenancePlanReport(dataCharts interface{}, dataTable interface{}, maintenanceID int) (string, error)
	CreateMaintenanceHistoryPlanReport(dataCharts interface{}, dataTable interface{}, history interface{}, maintenanceID int) (string, error)

	// GetMaintenanceStatus(maintenanceID int) (interface{}, error)

	GetMaintenanceBudget() (interface{}, error)

	GetMaintenanceHistory(idParent int, hisID int) (interface{}, error)
	CreateMaintenanceHistory(idParent int, userID int, reqs requests.MaintenanceRoadHistoryReq) (interface{}, error)
	UpdateMaintenanceHistory(idParent int, hisID int, userID int, reqs requests.MaintenanceRoadHistoryReq) (*int, *int, error)
	DeleteMaintenanceHistory(idParent int, hisID int) (interface{}, error)
	GetMaintenanceaAttachmentByID(id int) (*models.MaintenanceAttachment, error)

	GetMaintenanceRoadByID(IDParent int, mRoadId int) (interface{}, error)
	CreateMaintenanceRoad(idParent int, userID int, req requests.MaintenanceRoadsReq) (interface{}, error)
	UpdateMaintenanceRoad(idParent int, mRoadId int, userID int, req requests.MaintenanceRoadsReq) (*int, *int, error)
	DeleteMaintenanceRoad(idParent int, mRoadId int) error

	GetLastRoadInfoByID(roadID int) (*models.RoadInfoGeomData, error)

	GetDivisionList(userID int) (interface{}, error)
	GetRoadDropdownList(userID int) (interface{}, error)
	GetRoadDropdownListAnalyze(userID int) (interface{}, error)
	GetRoadDropdownListDashboard(userID int, dataType string, ownerCode string) (interface{}, error)
	//GetRoadDivisionFilter() (interface{}, error)
}

type Repository interface {
	// CreateMaintenance(request requests.MaintenanceReq, userID int) (interface{}, error)
	UpdateMaintenance(maintenanceID int, data models.Maintenance, tx *gorm.DB) error
	InsertMaintenance(data models.Maintenance, tx *gorm.DB) (int, error)
	InsertMaintenanceRoad(data []models.MaintenanceRoad, tx *gorm.DB) error
	UpdateMaintenanceRoad(data []models.MaintenanceRoad, tx *gorm.DB) error

	DeleteMaintenanceRoad(maintenanceID int, id int) error
	UpdateStatusMaintenanceRoad(ids []int, conn *gorm.DB) error
	DeleteMaintenance(maintenanceID int) error
	MaintenanceFinished(maintenanceID int, req requests.MaintenanceFinished) error
	MaintenanceCheckProgressComplete(maintenanceID int) (float64, error)
	CheckMaintenanceDuplicate(maintenanceID int, name string) (bool, error)
	GetRoadMaintenanceByID(maintenanceID int) (models.Maintenance, error)

	GetRoadGroupInfoByRoadID(roadID int, year string) ([]models.MaintainPreloadGetAll, error)
	GetRoadMaintenanceYear(roadID int) ([]models.MaintainPreloadGetAll, error)

	GetMaintenanceList(whereCondi string, prams requests.MaintenancePrams, isAllData, isOwnerData bool, depotCode string, limit, offset int64) ([]models.MaintenanceList, error)
	GetMaintenanceListCount(whereCondi string, prams requests.MaintenancePrams, isAllData, isOwnerData bool, depotCode string) (int64, error)
	GetMaintenanceListHistory(whereCondi string) ([]models.MaintainPreloadGetAll, error)
	//GetMaintenanceYearHistory(maintenanceID int) ([]models.MaintainPreloadGetAll, error)
	FindRefDirectionIDByRoadID(id int) (int, error)
	GetGeomByRoadID(roadID int) (models.RoadInfo, error)
	GetGeomByRoadGeomID(roadID int, lane int) (models.RoadGeom, error)
	GetMaintenanceListByID(maintenanceID int) (models.MaintainPreload, error)
	GetMaintenanceListByIDWithNotFilterIsComplete(maintenanceID int) (models.MaintainPreload, error)
	GetMaintenanceHistoryByID(maintenanceID int, mRoadHisId int) (models.MaintenanceRoadHistoryPreloadByID, error)
	GetTotalKmByMaintenanceID(maintenanceID int) (float64, error)
	// GetTotalProgressByID(maintenanceID int) (interface{}, error)
	GetRoadGroupInfoByMaintenanceID(maintenanceID int) (models.RoadGroup, error)
	StartTransSection() *gorm.DB
	RollBack(tx *gorm.DB) error
	Commit(tx *gorm.DB) error
	GetInterventionCriteriaByID(interventionCriteriaID int) (models.InterventionCriteria, error)
	GetInterventionCriteriaParamsLatest() (models.SettingInterventionCriteriaParams, error)

	GetMaintenanceRoad(maintenanceID, roadID int) ([]models.MaintenanceRoadData, error)
	GetRoadMaintenance(maintenanceID int) ([]models.MaintenanceRoadData, error)
	GetRoadMaintenanceHistory(maintenanceID int) ([]models.MaintenanceRoadHistoryData, error)
	GetMaintenanceRoadHistory(maintenanceID, roadID int) ([]models.MaintenanceRoadHistoryData, error)

	GetInterventionCriteria(ID int) ([]models.InterventionCriteria, error)
	GetRefCriteriaMethod() ([]models.RefCriteriaMethodData, error)

	GetMaintenanceBudget() (interface{}, error)

	GetMaintenanceRoadTheGeomByID(maintenanceID int) (models.MaintenanceRoad, error)

	// maintenance_history
	GetMaintenanceHistory(maintenanceID int) ([]models.MaintenanceRoadHistoryData, error)
	InsertMaintenanceRoadHistory(data []models.MaintenanceRoadHistory, conn *gorm.DB) error
	CreateMaintenanceRoadHistory(data models.MaintenanceRoadHistory, tx *gorm.DB) (*int, error)
	UpdateMaintenanceRoadHistory(data models.MaintenanceRoadHistory, attReqs []requests.MaintenanceAttachmentsReq, historyID int, ids []int) error
	DeleteMaintenanceHistory(maintenanceID int, hisID int) error
	UpdateStatusMaintenanceRoadHistory(ids []int, conn *gorm.DB) error

	UpdateIDParentMaintenanceRoadHistory(maintenanceRoadHisID int, idParent int, conn *gorm.DB) error

	GetMaintenanceaRoadHisAttachmentByMaintenanceID(maintenanceID int) ([]models.MaintenanceRoadHistoryAttachment, error)
	GetMaintenanceaRoadHisAttachmentByMRoadHisID(mRoadHisID int) ([]models.MaintenanceRoadHistoryAttachment, error)
	//InsertMaintenanceaRoadHisAttachment(data []models.MaintenanceRoadHistoryAttachment, conn *gorm.DB) error
	DeleteMaintenanceRoadHisAttachment(ids []int, conn *gorm.DB) error
	UpdateMaintenanceRoadHisAttachment(oldMaintenanceID int, newMaintenanceID int, conn *gorm.DB) error
	GetMaxRevisionMaintenanceRoadHistory(maintenanceRoadID int) (*int, *int, error)
	StAstext(TheGeom string) (string, error)

	GetMaintenanceYear(roadId int) ([]models.Maintenance, error)
	InsertMaintenanceaAttachment(data []models.MaintenanceAttachment, conn *gorm.DB) error
	UpdateMaintenanceAttachment(oldMaintenanceID int, newMaintenanceID int, conn *gorm.DB) error
	DeleteMaintenanceAttachment(idsAttDelete []int, conn *gorm.DB) error
	GetMaintenanceaAttachmentByID(id int) (*models.MaintenanceAttachment, error)
	GetMaintenanceaAttachmentByMaintenanceID(id int) ([]models.MaintenanceAttachment, error)
	UpdateIDParentMaintenance(maintenanceID int, idParent int, conn *gorm.DB) error
	GetMaxRevisionMaintenanceByIDParent(maintenanceID int) (*int, *int, error)
	UpdateStatusMaintenance(maintenanceID int, conn *gorm.DB) error

	GetMaintenanceRoadByMaintenanceID(maintenanceID int) ([]models.MaintenanceRoad, error)
	GetMaintenanceRoadHistoryByMaintenanceID(maintenanceID int) ([]models.MaintenanceRoadHistory, error)

	GetMaintenanceRoadID(maintenanceID int, mRoadID int) (*models.MaintenanceRoadPreloadById, error)
	CreateMaintenanceRoad(data models.MaintenanceRoad, conn *gorm.DB) (*int, error)
	UpdateIDParentMaintenanceRoad(maintenanceID int, idParent int, conn *gorm.DB) error
	GetMaxRevisionMaintenanceRoad(maintenanceRoadID int) (*int, *int, error)

	GetLastRoadInfoByID(roadID int) (*models.RoadInfoGeomData, error)
	GetRoadGroupByRoadId(roadID int) (*int, error)

	GetMaintenanceByRoadID(roadId int, year int) ([]models.MaintainPreload, error)
	GetMaintenanceByIDParent(IDParent int) (models.Maintenance, error)
	GetGeomJsonFromMaintenanceRoadID(maintenanceRoadID int) ([]byte, error)
	GetGeomJsonByMaintenanceRoadIDs(ids []int) (map[int][]byte, error)
	GetGeomJsonFromMaintenanceRoadHistoryID(maintenanceRoadHisID int) ([]byte, error)
	GetGeomJsonByMaintenanceRoadHistoryIDs(ids []int) (map[int][]byte, error)

	GetDivisionList(isAllData, isOwnerData bool, depotCode string) ([]models.RefDivisionList, error)
	GetRoadDropdownList(isAllData, isOwnerData bool, depotCode string) ([]models.RoadListInit, error)
	GetRoadDropdownListDashboard(isAllData, isOwnerData bool, depotCode string, queryRelated *string) ([]models.RoadListInit, error)
	//GetRoadDivisionFilter() ([]models.RoadDivisionFilter, error)
	GetSettingBudgetMethodByID(ID int) (models.SettingBudgetMethod, error)
}
