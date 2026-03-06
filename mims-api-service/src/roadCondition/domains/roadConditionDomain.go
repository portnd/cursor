package domains

import (
	"github.com/gin-gonic/gin"
	models "gitlab.com/mims-api-service/models"
	requests "gitlab.com/mims-api-service/requests"
	responses "gitlab.com/mims-api-service/responses"
	"gorm.io/gorm"
)

// business logic
type RoadConditionUseCase interface {
	GetMenu(uint) ([]models.AccessControl, error)
	GetRoadConditionList([]string, int) (interface{}, error)
	GetRoadConditionDetails(uid uint, conditionRangType string, idParent int) (interface{}, error)
	CreateRoadCondition(*gin.Context, uint, requests.RoadCondition, requests.RoadConditionFiles) (int, error)
	UpdateRoadCondition(*gin.Context, uint, int, requests.RoadConditionUpdate, requests.RoadConditionFiles, string, string) (interface{}, error)
	GetRoadCondition(*gin.Context, uint, int) (responses.RoadCondition, error)
	DeleteRoadCondition(*gin.Context, uint, int) (bool, error)
	GetRoadConditionTemplate(uint, int) (interface{}, error)
	GetRoadConditionCompareLane(int, requests.RoadConditionCompare) (interface{}, error)
	GetRoadConditionCompareYear(int, requests.RoadConditionCompare) (interface{}, error)
	GetRoadConditionCompareAverage(int, int) (interface{}, error)
}

// อะไรเชื่อมต่อกับ DB
type RoadConditionRepository interface {
	GetRole(uint) ([]models.UserRole, error)
	GetAccessControl([]int) ([]models.AccessControl, error)
	GetUserById(uint) (models.Users, error)
	GetUserDepartmentById(int) (models.UserDepartment, error)
	GetRoadConditionList(string, int) ([]models.RoadConditionList, error)
	//GetRoadConditionDetails(int, models.SqlCondition, models.SqlCondition) ([]models.RoadConditionDetails, error)
	//GetRoadConditionGrades(int) ([]models.RoadConditionGrade, error)
	GetRoadById(int) (models.Road, error)
	GetRoadPreloadById(int) (models.RoadPreload, error)
	GetFullGeom(int, int) (models.FullGeom, error)
	CreateRoadCondition(*gorm.DB, models.RoadCondition) (int, error)
	CreateRoadConditionSurvey(*gorm.DB, models.RoadConditionSurvey) (int, error)
	CreateRoadConditionSurvey100M(*gorm.DB, models.RoadConditionSurvey100M) (int, error)
	CreateRoadConditionSurveyM(*gorm.DB, models.RoadConditionSurveyM) error
	GetRoadConditionByIdParent(int) (models.RoadCondition, error)
	GetRoadConditionBeforeLastRevitionByIdParent(int, int) (models.RoadCondition, error)
	GetAllRoadConditionByIdParent(int) (models.RoadConditionAll, error)
	UpdateRoadCondition(*gorm.DB, models.RoadCondition) (models.RoadCondition, error)
	UpdateRoadConditionNoIriFile(*gorm.DB, models.RoadCondition) (models.RoadCondition, error) // hariphan

	DeleteRoadConditionByID(ID int, UserID int) error

	UpdateStatusIByID(tx *gorm.DB, ID int, UserID uint) error
	UpdateRoadConditionServeyM(*gorm.DB, models.RoadConditionSurveyM) error
	DeleteRoadConditionByIDParent(idParent int) error
	DeleteRoadCondition(models.RoadCondition) error
	GetRoadConditionTemplate(int) (models.RoadConditionTemplate, error)
	GetRoadConditionCompare(int, models.RoadConditionCompareLane) ([]models.RoadConditionAll, error)
	GetRoadConditionCompareAverage(int, int) ([]models.RoadConditionAll, error)
	UpdateRoadConditionFilepath(tx *gorm.DB, rcID int, csvFilepath string, imgFilepath string) error
	UpdateRoadConditionImgPath(rcID int, imgPath string) error
	UpdateRoadConditionSurveyMImgPath(rcSurveyID int, imgPath string) error
	GetRoadKmRange(roadId int, kmStart int, kmEnd int, direction int) ([]models.RoadKmRage, error)
	StartTransSection() *gorm.DB
	RollBack(tx *gorm.DB) error
	Commit(tx *gorm.DB) error

	GetRoadByID(roadID int) (models.RoadInfoGeomDirection, error)
	// hariphan
	GetRoadSurfaceByRoadIDByLaneNo(roadID int, laneNo int) ([]models.RoadSurfaceData, error)
	GetRefSurface() (map[int]string, error)

	GetRoadDirectionByRoadID(roadID int) (models.RoadPreloadConditionAll, error)
}
