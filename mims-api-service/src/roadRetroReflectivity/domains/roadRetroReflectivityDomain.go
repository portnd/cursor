package domains

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gorm.io/gorm"
)

// business logic
type UseCase interface {
	CreateRoadRetroReflectivity(c *gin.Context, uid uint, req requests.RoadRetroReflectivity, files requests.RoadRetroReflectivityFiles) (responses.RoadRetroReflectivityCreate, error)
	UpdateRoadRetroReflectivity(c *gin.Context, uid uint, IDParent int, rrsImport requests.RoadRetroReflectivity, files requests.RoadRetroReflectivityFiles, csvFilenameStatus string) (interface{}, error)
	DeleteRoadRetroReflectivity(c *gin.Context, uid uint, idParent int) (bool, error)
	GetRoadRetroReflectivityTemplate(uid uint, roadID int) (interface{}, error)
	GetRoadRetroReflectivity(c *gin.Context, uid uint, idParent int) (responses.RoadRetroReflectivity, error)
	GetRoadRetroReflectivityList(roadID int) (interface{}, error)
	GetRoadRetroReflectivityLineList(roadID int) (interface{}, error)
	GetRoadRetroReflectivityDetails(uid uint, rangeType string, refStripeTypeIDs string, idParent int) (interface{}, error)
	GetRoadRetroReflectivityCompareLine(int, requests.RoadRetroReflectivityCompare) (interface{}, error)
	GetRoadRetroReflectivityCompareYear(int, requests.RoadRetroReflectivityCompare) (interface{}, error)
	GetRoadRetroReflectivityCompareAverage(int, int) (interface{}, error)
}

// อะไรเชื่อมต่อกับ DB

// อะไรเชื่อมต่อกับ DB
type Repository interface {
	GetRole(uint) ([]models.UserRole, error)
	GetAccessControl([]int) ([]models.AccessControl, error)
	GetUserById(uint) (models.Users, error)
	GetUserDepartmentById(int) (models.UserDepartment, error)

	GetRoadRetroReflectivityDetailsByIdParent(idParent int, refStripeTypeIDs []string) (models.RoadRetroReflectivityPreload, error)
	GetRoadByID(roadID int) (*models.RoadById, error)
	GetRoadPreloadById(int) (models.RoadPreload, error)
	GetFullGeom(int, int) (models.FullGeom, error)
	CreateRoadRetroReflectivity(*gorm.DB, models.RoadRetroReflectivity) (int, error)
	CreateRoadRetroReflectivityRange(*gorm.DB, models.RoadRetroReflectivityRange) (int, error)
	CreateRoadRetroReflectivityM(tx *gorm.DB, rsrRangeM models.RoadRetroReflectivityM) error
	GetLastRevisionByIdParent(idParent int) (models.RoadRetroReflectivity, error)
	GetRoadRetroReflectivityList(roadID int) ([]models.RetroReflectivityList, error)
	GetRoadRetroReflectivityByIdParent(int) (models.RoadRetroReflectivity, error)
	GetRoadRefactiveStripBeforeLastRevitionByIdParent(idParent int, revision int) (models.RoadRetroReflectivity, error)
	GetAllRoadRetroReflectivityByIdParent(idParent int) (models.RoadRetroReflectivityPreload, error)
	UpdateRoadRetroReflectivity(tx *gorm.DB, rrs models.RoadRetroReflectivity) (models.RoadRetroReflectivity, error)
	UpdateRoadRetroReflectivityNoIriFile(tx *gorm.DB, rrs models.RoadRetroReflectivity) (models.RoadRetroReflectivity, error)

	DeleteRoadRetroReflectivity(models.RoadRetroReflectivity) error
	DeleteRoadRetroReflectivityByID(ID int, UserID int) error

	UpdateStatusIByID(tx *gorm.DB, ID int, UserID uint) error
	UpdateRoadRetroReflectivityM(*gorm.DB, models.RoadRetroReflectivityM) error
	DeleteRoadConditionByIDParent(idParent int) error

	GetRoadConditionTemplate(int) (models.RoadConditionTemplate, error)
	GetRoadRetroReflectivityCompare(roadId int, req models.RoadRetroReflectivityCompareLine) ([]models.RoadRetroReflectivityPreload, error)
	GetRoadRetroReflectivityAverage(int, int) ([]models.RoadRetroReflectivityPreload, error)

	UpdateRoadRetroReflectivityFilepath(tx *gorm.DB, rcID int, csvFilepath string) error
	GetRoadKmRange(roadId int, kmStart int, kmEnd int, direction int) ([]models.RoadKmRage, error)
	StartTransSection() *gorm.DB
	RollBack(tx *gorm.DB) error
	Commit(tx *gorm.DB) error

	GetRoadDirectionByRoadID(roadID int) (models.RoadPreloadConditionAll, error)
	GetRefStripeColor() ([]models.RefStripeColor, error)
	GetRefStripeType() ([]models.RefStripeType, error)

	GetRoadDirectionLaneList(roadID int) (models.RoadInfoGeomDirection, error)

	GetTotalLanesByRoadID(roadID int) (int64, error)

	GetLineListByRoadID(roadId int) ([]models.RoadRetroReflectivityPreload, error)
}
