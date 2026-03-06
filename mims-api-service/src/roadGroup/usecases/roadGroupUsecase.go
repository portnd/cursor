package usecases

import (
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadGroup/domains"
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

func (u *UseCase) GetRoadGroup() (interface{}, error) {
	data, err := u.Repo.GetRoadGroup()
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	return data, nil
}

func (u *UseCase) GetRoadGroupByID(id int) (models.RoadGroupByID, error) {
	result, err := u.Repo.GetRoadGroupByID(id)
	if err != nil {
		logs.Error(err)
		return result, responses.NewAppErr(404, err.Error())
	}
	roads, err := u.Repo.GetRoadByRoadGroupID(id)
	if err != nil {
		logs.Error(err)
		return result, responses.NewAppErr(404, err.Error())
	}
	for i, v := range roads {
		lanes, err := u.Repo.GetLaneByRoadID(v.ID)
		if err != nil {
			logs.Error(err)
			return result, responses.NewAppErr(404, err.Error())
		}
		roads[i].Lanes = lanes
	}
	result.Roads = roads
	return result, nil
}
