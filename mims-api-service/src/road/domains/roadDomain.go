package domains

import (
	"github.com/gin-gonic/gin"
	models "gitlab.com/mims-api-service/models"
	requests "gitlab.com/mims-api-service/requests"
	responses "gitlab.com/mims-api-service/responses"
	"gorm.io/gorm"
)

// business logic
type RoadUseCase interface {
	GetMenu(uint) ([]models.AccessControl, error)
	GetRoadDetailMenu(uint, []string, string) (interface{}, error)
	GetRoadGroupList(userID int, req requests.RoadPrams) (interface{}, error)
	GetRoadByID(roadID, userID int) (*responses.RoadById, error)
	GetRoadGroup() ([]models.RoadGroup, error)
	GetRoadTypeIcon() (map[string]int, error)
	// GetRoadDetailStatus(int, []string)
	GetRoadDirectionLaneList(int) (interface{}, error)
	GetRoadTree() (interface{}, error)
	CntVolumeApproved(roadGrpID int) (int, int, error)
	GetRoadSectionByID(int) (*models.RoadSectionById, error)
	GetDataById(interface{}, int) error
	GetRoadInit(id, level, refRoadTypeId int) (*responses.RoadInit, error)

	CreateRoad(c *gin.Context, userID int, req requests.Road) (interface{}, error)
	UpdateRoad(c *gin.Context, roadId int, userID int, req requests.RoadUpdate) (interface{}, error)
	UpdateRoadInit(c *gin.Context, roadId int, userID int, req requests.RoadUpdateInit) (interface{}, error)
	DeleteRoad(roadID int, userID int) (interface{}, error)
	GetLastRoadInfoByID(roadID int) (*models.RoadInfoGeomData, error)

	GetRoadLanes(roadID int) ([]responses.RoadLanes, error)

	GetRfp(params requests.RoadPrams) (interface{}, error)
}

// อะไรเชื่อมต่อกับ DB
type RoadRepository interface {
	GetRole(uint) ([]models.UserRole, error)
	GetAccessControl([]int) ([]models.AccessControl, error)
	GetRoadDetailMenu([]string, string) ([]responses.RoadMenuData, error)
	GetUserById(uint) (models.Users, error)
	GetUserDepartmentById(int) (models.UserDepartment, error)
	GetRoadGroupList(req requests.RoadPrams, roadIDs []int, isAllData, isOwnerData bool, depotCode string) ([]models.RoadList, error)
	GetRoadByID(roadID int) (*models.RoadById, error)
	GetRoadByRoadGrpID(roadGrpID int) ([]models.RoadInfo, error)
	GetRoadGroup() ([]models.RoadGroup, error)
	GetRoadTypeIcon() ([]models.RefRoadTypeIcon, error)
	GetRoadStatusSurface(int) (responses.StatusCount, error)
	GetRoadDirectionLaneList(int) (models.RoadInfoGeomDirection, error)
	GetVolumeAadtApproved(roadGrpID int) (responses.StatusCount, error)
	GetVolumeAccidentApproved(roadGrpID int) (responses.StatusCount, error)

	GetRoadInfoByRoadID(roadID int) (models.RoadInfo, error)

	GetRoadMaxSeq() (int, error)
	CreateData(tx *gorm.DB, model interface{}) error
	CreateRoadInfo(tx *gorm.DB, roadInfo *models.RoadInfo) error
	CreateRoadGeom(tx *gorm.DB, roadGeoms []models.RoadGeom) error
	UpdateRoadInfo(tx *gorm.DB, roadInfo *models.RoadInfo) error
	UpdateRoadGeom(tx *gorm.DB, roadGeoms []models.RoadGeom) error

	GetRoadSectionByID(int) (*models.RoadSectionById, error)
	GetDataById(interface{}, int) error
	// GetRoadStatusAsset(int) ([]responses.StatusCount, error)
	// GetRoadStatusCout(int)
	// GetRoadDetailStatus(int) error
	GetLastRoadInfoByID(roadID int) (models.RoadInfoGeomData, error)
	UpdateDirectionRoad(tx *gorm.DB, roadID, refDirectionId int) error
	CountParent(tx *gorm.DB, roadId int) (*int64, error)
	DeleteRoad(tx *gorm.DB, roadID int, userID int) error

	DeleteRoadInfo(tx *gorm.DB, roadID int, userID int) error

	DeleteRoadGeom(tx *gorm.DB, userID int, roadID int) error

	GetRoadLanes(roadID int) ([]models.RoadLanes, error)

	StartTransSection() *gorm.DB
	RollBack(tx *gorm.DB) error
	Commit(tx *gorm.DB) error

	GetParamsConditionByID(ID int, conditionType string) (models.ParamsCondition, error)
	GetParamsRoadLine(ID int) (models.ParamsRoadLine, error)
	GetroadConditionSurvey(roadID int, laneNo int) ([]models.RoadConditionSurvey, error)
	GetRoadSurfaceLaneAll(roadIds []int) (map[int][]models.RoadSurfaceLane, error)
	GetroadConditionSurvey100M(roadID, laneNo int) ([]models.RoadConditionSurvey100M, error)
	GetRoadRetroReflectivity100M(roadID, laneNo int) ([]models.RoadRetroReflectivityRange, error)
	GetRoadID(params requests.RoadPrams) ([]int, error)
}
