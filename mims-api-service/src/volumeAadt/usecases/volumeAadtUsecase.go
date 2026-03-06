package usecases

import (
	"fmt"
	"time"

	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/volumeAadt/domains"
	"gorm.io/gorm"
)

type UseCase struct {
	Repo domains.Repository
}

// init usecase
func NewUseCase(repo domains.Repository) domains.UseCase {
	return &UseCase{
		Repo: repo,
	}
}

func (u *UseCase) GetVolumeRevision(roadGrpID int, permissions []string) (interface{}, error) {
	data, err := u.Repo.GetVolumeRevision(roadGrpID, permissions)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	// Group unique IDParent entries by year, keeping the latest revision per IDParent.
	// Data is pre-sorted by year DESC, surveyed_date DESC, revision DESC.
	temp := make(map[int]bool)
	dataYear := make(map[int][]models.VolumeAadtRevision)
	for _, item := range data {
		if !temp[item.IDParent] {
			temp[item.IDParent] = true
			dataYear[item.Year] = append(dataYear[item.Year], item)
		}
	}

	var response []responses.VolumeAadtRevision
	for year, items := range dataYear {
		var entry responses.VolumeAadtRevision
		entry.Year = year
		entry.Item = items
		response = append(response, entry)
	}
	if len(response) == 0 {
		return []string{}, nil
	}
	return response, nil
}

func (u *UseCase) GetVolume(roadGrpID, aadtID int) (interface{}, error) {
	// var result responses.VolumeAadtRespond
	data, err := u.Repo.GetVolume(roadGrpID, aadtID)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return data, responses.NewAppErr(400, err.Error())
	}
	// 	result.CreatedDate = helpers.SetTimeToString(data.CreatedDate)
	// 	result.UpdatedDate = helpers.SetTimeToString(data.UpdatedDate)

	// get geom
	// roadInfo, err := u.Repo.GetTheGeomByRoadGrpID(roadGrpID)
	// if err != nil {
	// 	logs.Error(err)
	// 	if err == gorm.ErrRecordNotFound {
	// 		return responses.NoData{}, nil
	// 	}
	// 	return data, responses.NewAppErr(400, err.Error())
	// }
	status := ""
	if data.Status != "I" {
		status, err = u.Repo.GetStatus(data.Status)
		if err != nil {
			logs.Error(err)
			if err == gorm.ErrRecordNotFound {
				return responses.NoData{}, nil
			}
			return data, responses.NewAppErr(400, err.Error())
		}
	} else {
		revision := data.Revision
		if revision == 0 {
			status = "ข้อมูลเริ่มต้น"
		} else {
			status = fmt.Sprintf("ครั้งที่ %d", revision)
		}
	}
	data.StatusCode = data.Status
	data.Status = status
	return data, nil
}

func (u *UseCase) CreateVolume(roadId, IDParent, aadtID, userID int, req requests.VolumeAadtReq) (interface{}, error) {
	var newData models.VolumeAadt
	if IDParent != 0 {
		maxRevision, err := u.Repo.GetMaxRevision(roadId, IDParent)
		if err != nil {
			return "", responses.NewAppErr(400, err.Error())
		}

		aadt, err := u.Repo.GetVolumeByID(roadId, aadtID)
		if err != nil {
			return "", responses.NewAppErr(400, err.Error())
		}

		layout := "2006-01-02"
		surveyedDate, err := time.Parse(layout, req.SurveyedDate)
		if err != nil {
			return "", responses.NewAppErr(400, err.Error())
		}

		newData.Year = surveyedDate.Year()
		newData.SurveyedDate = surveyedDate
		newData.Veh1 = req.Veh1
		newData.Veh2 = req.Veh2
		newData.Veh3 = req.Veh3
		newData.CreatedBy = aadt.CreatedBy
		newData.UpdatedBy = userID
		newData.CreatedDate = aadt.CreatedDate
		newData.UpdatedDate = time.Now().UTC()
		newData.RoadId = roadId
		newData.IDParent = maxRevision.IDParent
		newData.Revision = maxRevision.Revision + 1
		newData.Total = req.Veh1 + req.Veh2 + req.Veh3

		newData.Status = "A"
		if aadt.Status == "A" {
			err := u.Repo.UpdateStatusI(aadt.IDParent)
			if err != nil {
				return "", responses.NewAppErr(400, err.Error())
			}
		}
	} else {
		// if maxRevision.Status == "T" || maxRevision.Status == "R" {
		// 	err := u.Repo.UpdateStatusD(maxRevision.ID)
		// 	if err != nil {
		// 		return "", responses.NewAppErr(400, err.Error())
		// 	}
		// } else if maxRevision.Status == "W" {
		// 	return "", errors.New(constants.DATA_WAITING_APPROVAL)
		// }
		layout := "2006-01-02"
		surveyedDate, err := time.Parse(layout, req.SurveyedDate)
		if err != nil {
			return "", responses.NewAppErr(400, err.Error())
		}
		newData.Year = surveyedDate.Year()
		newData.SurveyedDate = surveyedDate
		newData.Veh1 = req.Veh1
		newData.Veh2 = req.Veh2
		newData.Veh3 = req.Veh3
		newData.CreatedBy = userID
		newData.UpdatedBy = userID
		newData.CreatedDate = time.Now().UTC()
		newData.UpdatedDate = time.Now().UTC()
		newData.RoadId = roadId
		newData.Total = req.Veh1 + req.Veh2 + req.Veh3
		newData.Revision = 0
		newData.Status = "A"
		newData.IDParent = 0
	}

	// err := u.Repo.UpdateVolumeStatusT_To_DByGrpID(roadGrpID)
	// if err != nil {
	// 	return "", responses.NewAppErr(400, err.Error())
	// }

	ID, IDParent, err := u.Repo.CreateVolume(newData)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	// update id_parent new volume
	if IDParent == 0 {
		err = u.Repo.UpdateVolumUpdateIdParent(ID)
		if err != nil {
			return "", responses.NewAppErr(400, err.Error())
		}
	}

	var volume responses.Volume
	volume.ID = ID
	if IDParent != 0 {
		volume.IDParent = IDParent
	} else {
		volume.IDParent = ID
	}
	return volume, nil
}

func (u *UseCase) DeleteVolume(roadId, aadtID, userID int) (interface{}, error) {
	err := u.Repo.DeleteVolume(roadId, aadtID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		} else {
			return "", responses.NewAppErr(400, err.Error())
		}
	}

	err = u.Repo.UpdateStatusIToA(roadId, aadtID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		} else {
			return "", responses.NewAppErr(400, err.Error())
		}

	}
	return "", nil
}
