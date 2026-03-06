package domains

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
)

// business logic
type RoadAssetUseCase interface {
	GetRoadAssetDetail(params requests.AssetDetailsQueryParams, permissions []string, roadAssetId int, userID int) (interface{}, int, error)
	GetRoadAssetPermission(params requests.AssetPermissionQueryParams, permissions []string) (interface{}, error)
	GetRoadAssetRevisions(params requests.AssetRevisionsQueryParams, permissions []string, roadID int) (interface{}, error)
	GetRoadAssetTemplate(params requests.AssetTemplateQueryParams, permissions []string) (interface{}, error)
	CreateRoadAsset(reqs map[string]interface{}, roadID int, IDParentAsset, userID int) (interface{}, error)
	ConfirmRoadAsset(idParent int, userID uint) (int, error)
	CancelRoadAsset(idParent int, userID uint) (interface{}, error)
	DeleteRoadAsset(c *gin.Context, idParent int, userID uint) (interface{}, error)
	DeleteRoadAssetObject(assetID int, refAssetTableID int, parentAssetID int, userID uint) (interface{}, error)
	GetRoadKmByGeom(geom string, roadID int) (responses.RoadKmByGeom, error)
	GetAssetTableByID(roadAssetId int) (interface{}, error)
	GetRoadInfoByRoadID(roadID int) (models.RoadInfo, error)
}

// อะไรเชื่อมต่อกับ DB
type RoadAssetRepository interface {
	GetRoadAssetDetailColumn(permissions []string, roadAssetID int, refAssetTableID int) ([]responses.RoadAssetDetailColumn, error)
	GetRoadAssetDetailData(int, []responses.RoadAssetDetailColumn, []string, int, int, bool) ([]map[string]interface{}, error)
	GetRoadAssetDetailInfo(int, []string) (models.RoadAssetRefDataStatus, error)
	IsApproverAssetTableStaff(int, int) (bool, error)
	GetRoadAssetRevisions(params requests.AssetRevisionsQueryParams, permissions []string, roadID int) ([]responses.RoadAssetRevision, error)
	GetRoadAssetTemplateColumn(params requests.AssetTemplateQueryParams, permissions []string) ([]responses.RoadAssetTemplateColumn, error)
	GetRoadAssetTemplateData(params requests.AssetTemplateQueryParams, colList []responses.RoadAssetTemplateColumn) (interface{}, error)

	GetMaxRevision(roadID int, idParent int, refAssetTableID int) (models.RaData, error)
	UpdateRoadAssetStauts(RoadAssetId int, status string) error
	UpdateRoadAssetUpdatedDate(RoadAssetId int, updatedDate time.Time) error
	UpdateRoadAssetIDParent(RoadAssetId int, idParent int) error
	UpdateRoadAssetData(RoadAssetId int, status string) error
	GetRefAssetTableById(ID int) (models.RefAssetTable, error)
	InsertRoadAsset(models.RoadAsset) (int, error)
	LoadData(tableName string, RoadAssetId int) ([]map[string]interface{}, error)
	GetRoadAssetStatusTByIdParent(idParent int) (models.RoadAsset, error)
	GetLastRoadAssetByIdParent(idParent int, status string) (models.RoadAsset, error)
	GetRoadAssetExclusiveLockByIdParent(idParent int, exclusiveLock bool, status string) (models.RoadAsset, error)
	GetRoadAssetByID(assetId int, status string) (models.RoadAsset, error)
	CreateRoadAssetRevision(roadAsset models.RoadAsset) (int, error)
	UpdateConfirmRoadAsset(roadAssetID int, userID int) (models.RoadAsset, error)
	UpdateCancelRoadAsset(roadAssetID int, userID int) (models.RoadAsset, error)
	UpdateRoadAssetStatus(roadAssetID int, status string) (bool, error)
	DeleteRoadAssetTable(idParent int, roadAssetID int, tableName string) (bool, error)
	DeleteRoadAssetTableByRoadAssetID(roadAssetID int, tableName string) (bool, error)
	DeleteRoadAssetTableByID(assetObjectID int, tableName string) (bool, error)
	UndeleteRoadAssetTableByID(assetObjectID int, tableName string) (bool, error)
	GetAllRoadAssetTableByAssetId(roadAssetID int, tableName string) ([]interface{}, error)
	CreateAssetTableFromOldAsset(oldAssetData []map[string]interface{}, tableName string) error
	CreateAssetTableFromOldAssetByNewRoadAssetId(roadAssetID int, oldAssetData []map[string]interface{}, tableName string) error
	RawQuery(sql string) error
	RawQueryInsert(sql string) (int, error)
	FileColumnFilepath(refAssetTableID int) ([]models.RefAssetTableColumns, error)
	GetOldAsset(idParentAsset int, maxRevisionID int, tableName string, selects []string) ([]map[string]interface{}, error)
	UpdateTableIsDeletedByRoadAssetID(roadAssetID int, tableName string) error
	UpdateTableIsDeletedByID(ID int, tableName string) error
	UpdateTableIsDeletedByRoadAssetIDAndIDParent(assetID int, IDParent int, tableName string) error
	UpdateIDParentByID(ID int, IDParent int, tableName string) error
	UpdateFilePathByID(ID int, column string, filePath string, tableName string) error
	GetRoadKmByGeomLine(geomString string, roadID int) (float64, float64, int, error)
	GetRoadKmByGeomPoint(geomString string, roadID int) (float64, int, error)
	UpdatetTriggerHashData(ID int, tableName string, hash string) error
	GetStatus(statusCode string) (string, error)
	GetClosestRoad(roadID int, point string) (string, error)
	GetAssetTableByID(assetTableId int) (models.RefAssetTable, error)
	GetAssetTableByIDStaff(assetTableID int) ([]models.RefAssetTableStaff, error)
	GetRoadInfoByRoadID(roadID int) (models.RoadInfo, error)
	GetUserDepartmentById(userID int) (models.UserDepartment, error)
	UpdateRoadAssetStatuByIdParent(idParent int, status string) (models.RoadAsset, error)
}
