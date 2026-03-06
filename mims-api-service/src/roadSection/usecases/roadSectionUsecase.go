package usecases

import (
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadSection/domains"
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

func (u *UseCase) GetRoadSection(roadGroupId *int) ([]models.RoadSection, error) {
	data, err := u.Repo.GetRoadSection(roadGroupId)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewAppErr(400, err.Error())
	}
	return data, nil
}

func (u *UseCase) GetRoadSectionByID(id int) (*models.RoadSection, error) {
	result, err := u.Repo.GetRoadSectionByID(id)
	if err != nil {
		logs.Error(err)
		return result, responses.NewAppErr(404, err.Error())
	}
	return result, nil
}
