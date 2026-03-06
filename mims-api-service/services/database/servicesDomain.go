package services

import (
	"gitlab.com/mims-api-service/responses"
)

//mockery --dir=services/database --name=ServicesDatabaseDomain --filename=database_mock.go --output=services/database/mocks --outpkg=servicemocks

type ServicesDatabaseDomain interface {
	UserInfo(userID int) (*responses.UserInfoRespond, error)
}
