package usecases

import (
	"fmt"
	"time"

	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/volumeAccident/domains"
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
	// helpers.PrintlnJson(data)
	year := make(map[int]int)
	dataYear := make(map[int][]models.VolumeAccidentRevision)
	var volumes []models.VolumeAccidentRevision
	temp := make(map[int]string)
	for _, item := range data {
		_, ok := year[item.Year]
		if !ok {
			year[item.Year] = item.Year
		}

		_, isVal := temp[item.IDParent]
		if !isVal {
			temp[item.IDParent] = item.Status
			volumes = append(volumes, item)
			dataYear[item.Year] = volumes
		}

	}

	var response []responses.VolumeAccidentRevision
	for year, items := range dataYear {
		var data responses.VolumeAccidentRevision
		data.Year = year
		for _, item := range items {

			if item.SurveyedDate.Year() == year {
				data.Item = append(data.Item, item)
			}
		}
		response = append(response, data)
	}
	if len(response) == 0 {
		return []string{}, nil
	}
	return response, nil
}

func (u *UseCase) GetVolume(roadGrpID, ID int) (interface{}, error) {
	// var result responses.VolumeAccidentRespond
	data, err := u.Repo.GetVolume(roadGrpID, ID)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return data, responses.NewAppErr(400, err.Error())
	}
	// return data, nil
	// get geom
	roadInfo, err := u.Repo.GetTheGeomByRoadGrpID(roadGrpID)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return data, responses.NewAppErr(400, err.Error())
	}
	// get status
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
	data.RoadInfo = roadInfo
	// //
	// copier.Copy(&result, &data)
	// result.CreatedDate = helpers.SetTimeToString(data.CreatedDate)
	// result.UpdatedDate = helpers.SetTimeToString(data.UpdatedDate)

	return data, nil
}

func (u *UseCase) CreateVolume(roadGrpID, IDParent, accidentID, userID int, req requests.VolumeAccidentReq) (interface{}, error) {
	// return req, nil
	var newData models.VolumeAccident
	if IDParent != 0 {
		maxRevision, err := u.Repo.GetMaxRevision(roadGrpID, IDParent)
		if err != nil {
			return "", responses.NewAppErr(400, err.Error())
		}

		accident, err := u.Repo.GetVolumeByID(roadGrpID, accidentID)
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
		newData.Acc1 = req.Acc1
		newData.Acc2 = req.Acc2
		newData.Acc3 = req.Acc3
		newData.Acc4 = req.Acc4
		newData.CreatedBy = accident.CreatedBy
		newData.UpdatedBy = userID
		newData.CreatedDate = accident.CreatedDate
		newData.UpdatedDate = time.Now().UTC()
		newData.RoadGroupID = roadGrpID
		newData.Status = "T"
		newData.IDParent = maxRevision.IDParent
		newData.Revision = maxRevision.Revision + 1

		if accident.Status == "T" || accident.Status == "R" {
			err := u.Repo.UpdateStatusD(accident.IDParent)
			if err != nil {
				return "", responses.NewAppErr(400, err.Error())
			}
		} else if accident.Status == "W" {
			return "", responses.NewAppErr(400, constants.DATA_WAITING_APPROVAL)
		} else if accident.Status == "A" {
			// err := u.Repo.UpdateStatusI(accident.IDParent)
			// if err != nil {
			// 	return "", responses.NewAppErr(400, err.Error())
			// }
		}
	} else {
		layout := "2006-01-02"
		surveyedDate, err := time.Parse(layout, req.SurveyedDate)
		if err != nil {
			return "", responses.NewAppErr(400, err.Error())
		}
		newData.Year = surveyedDate.Year()
		newData.SurveyedDate = surveyedDate
		newData.Acc1 = req.Acc1
		newData.Acc2 = req.Acc2
		newData.Acc3 = req.Acc3
		newData.Acc4 = req.Acc4
		newData.CreatedBy = userID
		newData.UpdatedBy = userID
		newData.CreatedDate = time.Now().UTC()
		newData.UpdatedDate = time.Now().UTC()
		newData.RoadGroupID = roadGrpID
		newData.Revision = 0
		newData.Status = "T"
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

func (u *UseCase) DeleteVolume(roadGrpID, accidentID, userID int) (interface{}, error) {
	err := u.Repo.DeleteVolume(roadGrpID, accidentID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		} else {
			return "", responses.NewAppErr(400, err.Error())
		}
	}
	err = u.Repo.UpdateStatusIToA(roadGrpID, accidentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		} else {
			return "", responses.NewAppErr(400, err.Error())
		}

	}
	return "", nil
}
