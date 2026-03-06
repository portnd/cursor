package domains

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gorm.io/gorm"
)

type RoadSurfaceUsecases interface {
	GetRoadSurfaceList(roadId string, permissions []string, uID int) ([]responses.RoadSurfaceResponds, error)
	PostRoadSurface(request requests.RoadSurface, uid int) ([]int, error)
	GetTotalKm(roadID int) (float64, error)
	GetRoadSurfaceIconById(roadId int) ([]models.RoadSurfaceIcon, error)
}

type RoadSurfaceRepositories interface {
	GetRoadSurfaceData(roadId string, hasPermission bool) ([]models.RoadSurfacePreload, error)
	FindRefSurface(id int) (models.SurfaceShoulder, error)
	FindRefMaterial(table string, id int) (interface{}, error)
	FindRefDirectionName(roadID int) (string, error)
	FindRefSurfaceNameAndSurfaceGroup(id int) (string, string, error)
	FindDataStatus(code string) (string, error)
	GetRole(userId uint) ([]models.UserRole, error)
	GetUserByID(userId uint) (models.Users, error)
	GetUserInfoByUpdatedBy(updatedBy int) models.UserDepartment
	GetMaxRevisionByRoadID(roadID int) (int, error)
	FindRefDirectionIDByRoadID(id int) (int, error)
	InsertRoadSurface(data models.RoadSurfacePointer, tx *gorm.DB) (int, error)
	GetGeomByRoadID(roadID int) ([]models.RoadGeom, error)
	InsertRoadSurfaceLane(data models.RoadSurfaceLane, tx *gorm.DB) error
	UpdateIDParent(id, idParent int, tx *gorm.DB) error
	StartTransSection() *gorm.DB
	RollBack(tx *gorm.DB) error
	Commit(tx *gorm.DB) error
	ClearPreviousData(roadID int, uid int, kmStart float64, kmEnd float64, numberLane int, tx *gorm.DB) error
	GetNewIDs(roadID int) ([]int, error)
	GetTotalKm(roadID int) (float64, error)
	GetRoadSurfaceGroupByRoadID(roadID int) (models.RoadSurface, error)
	GetRefSurfaceParam() ([]models.RefSurfaceParam, error)
	UpdateStatusTToDeleteAll(roadID int, tx *gorm.DB) error
	GetDataList(model interface{}, where string) error
	GetLastRoadInfoByID(roadId int) (*models.RoadInfo, error)
	ClearPreviousDataStatus(roadId int, uid int, tx *gorm.DB) error
	GetRoadSurfaceIconById(roadId int) ([]models.RoadSurfaceIcon, error)
}
