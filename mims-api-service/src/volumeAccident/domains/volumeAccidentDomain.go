package domains

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
)

// business logic
type UseCase interface {
	GetVolumeRevision(roadGrpID int, permission []string) (interface{}, error)
	GetVolume(roadID int, IDParent int) (interface{}, error)
	// CreateVolume(roadGrpID int, userID int, req requests.VolumeAccidentReq) (interface{}, error)
	CreateVolume(roadGrpID int, IDParent, accidentID int, userID int, req requests.VolumeAccidentReq) (interface{}, error)
	DeleteVolume(roadGrpID int, accidentID int, userID int) (interface{}, error)
}

// อะไรเชื่อมต่อกับ DB
type Repository interface {
	GetVolumeRevision(roadGrpID int, permissions []string) ([]models.VolumeAccidentRevision, error)
	GetVolume(roadGrpID int, IDParent int) (models.VolumeAccidentList, error)
	GetVolumeByID(roadGrpID int, ID int) (models.VolumeAccident, error)
	UpdateVolumeStatusT_To_DByGrpID(roadGrpID int) error
	// CreateVolume(roadGrpID int, userID int, req requests.VolumeAccidentReq) (int, error)
	CreateVolume(data models.VolumeAccident) (int, int, error)
	// UpdateStatusA(volID int, req requests.VolumeAadtReq) error
	UpdateStatusD(volID int) error
	UpdateStatusI(volID int) error
	GetMaxRevision(roadGrpID int, IDParent int) (models.VolumeAccident, error)
	DeleteVolume(roadGrpID int, accidentID int, userID int) error
	UpdateStatusIToA(volID int, IDParent int) error
	UpdateVolumUpdateIdParent(ID int) error
	GetTheGeomByRoadGrpID(roadGrpID int) ([]models.RoadInfo, error)
	GetStatus(statusCode string) (string, error)
}
