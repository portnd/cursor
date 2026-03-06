package domains

import (
	"mime/multipart"

	models "gitlab.com/mims-api-service/models"
	requests "gitlab.com/mims-api-service/requests"
)

// business logic
type RoadDamageUseCase interface {
	GetRoadDamageList(roadID int, permission []string) (interface{}, error)
	GetRoadDamageForImport(int, int) (models.RoadDamage, error)
	GetRoadById(int) (models.Road, error)
	GetRoadGeom(int, int) (models.RoadGeom, error)
	UpdateRoadDamage(int, int, requests.RoadDamageImport) error
	UpdateRoadDamageStatusT(int, requests.RoadDamageImport) error
	UpdateRoadDamageStatusI(int, requests.RoadDamageImport) error
	SetRoadDamageFromImport(int, int, int, string, string, requests.RoadDamageImport, requests.RcData, *multipart.FileHeader, *multipart.FileHeader) (interface{}, error)
	GetRoadDamageByIDParent(int) (models.RoadDamage, error)
	GetDirectionById(int) (models.RefDirection, error)
	GetRoadDamageDetail(int, int, []string) (interface{}, error)
	GetRoadDamageTemplate(int) (interface{}, error)
	DeleteRoadDamageForImport(int, int) (interface{}, error)
	RoadDamageReadCSVFile(roadID int, filePath string) ([]requests.RoadDamageCsv, error)
}

// อะไรเชื่อมต่อกับ DB
type RoadDamageRepository interface {
	GetRoadDamageByIDParent(int) (models.RoadDamage, error)
	GetRoadDamageList(RoadID int, permissions []string) ([]models.RoadDamageList, error)
	GetRoadDamageForImport(int, int) (models.RoadDamage, error)
	GetRoadById(int) (models.Road, error)
	GetRoadGeom(int, int) (models.RoadGeom, error)
	CreateRoadDamage(requests.SumValues, requests.RcData) (models.RoadDamage, error)
	CreateRoadDamageRange(requests.RaodRangeItem, requests.RcData, int) (int, string, error)
	CreateRoadDamageM(requests.RoadDamageMItem, requests.RcData) error
	UpdateRoadDamage(int, int, requests.RoadDamageImport) error
	UpdateRoadDamageStatusT(int, requests.RoadDamageImport) error
	UpdateRoadDamageStatusI(int, requests.RoadDamageImport) error
	UpdateRoadDamageIDParent(int) error
	UpdateRoadDamageCsvFilepath(int, string) error
	UpdateRoadDamageImgFilepath(int, string) error
	GetRoadDamageById(int) (models.RoadDamage, error)
	GetDirectionById(int) (models.RefDirection, error)
	GetRoadDamageDirection(int, int, []string) (int, error)
	GetRoadDamageDetail(int, int, int, []string) (models.RoadDamageDetail, error)
	GetUserById(int) (models.UserDepartment, error)
	GetRoadDamageTemplate(int) (models.RoadInfoGeom, error)
	DeleteRoadDamageForImport(int, int) error
	GetRoadInfo(roadID int) (models.RoadInfo, error)
	GetRoadGeomByRoadIDByLaneNo(roadId int, laneNo int) (models.RoadGeom, error)
	UpdateRoadDamageImgPath(rdID int, imgPath string) error
	GetRoadSurfaceLane(roadID int, laneNo int, km float64) (string, error)
	GetLastRoadInfoByID(roadID int) (models.RoadInfoGeomData, error)
}
