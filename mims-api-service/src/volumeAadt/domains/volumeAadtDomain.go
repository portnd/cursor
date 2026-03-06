package domains

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
)

// business logic
type UseCase interface {
	GetVolumeRevision(roadGrpID int, permission []string) (interface{}, error)
	GetVolume(roadID int, aadtID int) (interface{}, error)
	CreateVolume(roadGrpID int, IDParent int, aadtID int, userID int, req requests.VolumeAadtReq) (interface{}, error)
	DeleteVolume(roadGrpID int, aadtID int, userID int) (interface{}, error)
}

// อะไรเชื่อมต่อกับ DB
type Repository interface {
	GetVolumeRevision(roadGrpID int, permissions []string) ([]models.VolumeAadtRevision, error)
	GetVolume(roadGrpID int, aadtID int) (models.VolumeAadtList, error)
	GetVolumeByID(roadGrpID int, aadtID int) (models.VolumeAadt, error)
	CreateVolume(data models.VolumeAadt) (int, int, error)
	UpdateVolumeStatusT_To_DByGrpID(roadGrpID int) error
	UpdateStatusD(volID int) error
	UpdateStatusI(IDParent int) error
	GetMaxRevision(roadGrpID int, IDParent int) (models.VolumeAadt, error)
	DeleteVolume(roadGrpID int, aadtID int, userID int) error
	UpdateStatusIToA(volID int, IDParent int) error
	UpdateVolumUpdateIdParent(ID int) error
	GetTheGeomByRoadGrpID(roadGrpID int) ([]models.RoadInfo, error)
	GetStatus(statusCode string) (string, error)
}
