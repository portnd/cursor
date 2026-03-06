package usecases

import (
	"strconv"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/src/hsms/domains"
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

func (u *UseCase) GetHsmsBridge() (interface{}, error) {
	data, err := u.Repo.GetHsmsBridge()
	if err != nil {
		err = helpers.MongoDbLog("Hsms Bridge", err.Error(), "MONGODB_HSMS_ERROR_LOG", false)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = helpers.MongoDbLog("Hsms Bridge", "บันทึนข้อมูลลง Postgresql จำนวน "+strconv.Itoa(len(data))+" ข้อมูล", "MONGODB_HSMS_ERROR_LOG", true)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *UseCase) GetHsmsGuard() (interface{}, error) {
	data, err := u.Repo.GetHsmsGuard()
	if err != nil {
		err = helpers.MongoDbLog("Hsms Guard", err.Error(), "MONGODB_HSMS_ERROR_LOG", false)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = helpers.MongoDbLog("Hsms Guard", "บันทึนข้อมูลลง Postgresql จำนวน "+strconv.Itoa(len(data))+" ข้อมูล", "MONGODB_HSMS_ERROR_LOG", true)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *UseCase) GetHsmsInterchange() (interface{}, error) {
	data, err := u.Repo.GetHsmsInterchange()
	if err != nil {
		err = helpers.MongoDbLog("Hsms Interchange", err.Error(), "MONGODB_HSMS_ERROR_LOG", false)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = helpers.MongoDbLog("Hsms Interchange", "บันทึนข้อมูลลง Postgresql จำนวน "+strconv.Itoa(len(data))+" ข้อมูล", "MONGODB_HSMS_ERROR_LOG", true)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *UseCase) GetHsmsIntersection() (interface{}, error) {
	data, err := u.Repo.GetHsmsIntersection()
	if err != nil {
		err = helpers.MongoDbLog("Hsms Intersection", err.Error(), "MONGODB_HSMS_ERROR_LOG", false)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = helpers.MongoDbLog("Hsms Intersection", "บันทึนข้อมูลลง Postgresql จำนวน "+strconv.Itoa(len(data))+" ข้อมูล", "MONGODB_HSMS_ERROR_LOG", true)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *UseCase) GetHsmsStreetlight() (interface{}, error) {
	data, err := u.Repo.GetHsmsStreetlight()
	if err != nil {
		err = helpers.MongoDbLog("Hsms Streetlight", err.Error(), "MONGODB_HSMS_ERROR_LOG", false)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = helpers.MongoDbLog("Hsms Streetlight", "บันทึนข้อมูลลง Postgresql จำนวน "+strconv.Itoa(len(data))+" ข้อมูล", "MONGODB_HSMS_ERROR_LOG", true)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *UseCase) GetHsmsRailwaycrossing() (interface{}, error) {
	data, err := u.Repo.GetHsmsRailwaycrossing()
	if err != nil {
		err = helpers.MongoDbLog("Hsms Railwaycrossing", err.Error(), "MONGODB_HSMS_ERROR_LOG", false)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = helpers.MongoDbLog("Hsms Railwaycrossing", "บันทึนข้อมูลลง Postgresql จำนวน "+strconv.Itoa(len(data))+" ข้อมูล", "MONGODB_HSMS_ERROR_LOG", true)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *UseCase) GetHsmsTrafficlight() (interface{}, error) {
	data, err := u.Repo.GetHsmsTrafficlight()
	if err != nil {
		err = helpers.MongoDbLog("Hsms Trafficlight", err.Error(), "MONGODB_HSMS_ERROR_LOG", false)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = helpers.MongoDbLog("Hsms Trafficlight", "บันทึนข้อมูลลง Postgresql จำนวน "+strconv.Itoa(len(data))+" ข้อมูล", "MONGODB_HSMS_ERROR_LOG", true)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *UseCase) GetHsmsUturnbridge() (interface{}, error) {
	data, err := u.Repo.GetHsmsUturnbridge()
	if err != nil {
		err = helpers.MongoDbLog("Hsms Uturnbridge", err.Error(), "MONGODB_HSMS_ERROR_LOG", false)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = helpers.MongoDbLog("Hsms Uturnbridge", "บันทึนข้อมูลลง Postgresql จำนวน "+strconv.Itoa(len(data))+" ข้อมูล", "MONGODB_HSMS_ERROR_LOG", true)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
