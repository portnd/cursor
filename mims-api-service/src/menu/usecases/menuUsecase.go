package usecases

import (
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/menu/domains"
)

type menuUseCase struct {
	menuRepo domains.MenuRepository
}

// init usecase
func NewMenuUseCase(repo domains.MenuRepository) domains.MenuUseCase {
	return &menuUseCase{
		menuRepo: repo,
	}
}

// =========================================================
func (t *menuUseCase) GetMenu(accCtrls map[string]string) ([]responses.AccessGroupMenu, error) {
	var accessGroups []responses.AccessGroupMenu
	data, err := t.menuRepo.GetMenu()
	if err != nil {
		return accessGroups, err
	}
	helpers.PrintlnJson(accCtrls)
	for _, item := range data {
		var accessGroup responses.AccessGroupMenu
		keys := []string{}
		keys = helpers.Explode(item.AccessKey, ",")
		i := 0
		for _, key := range keys {
			// fmt.Println(accCtrls[key], "===", key)
			if accCtrls[key] == key {
				if i == item.Id {
					continue
				}
				i = item.Id
				accessGroup.Id = item.Id
				accessGroup.Name = item.Name
				accessGroup.ParentId = item.ParentId
				accessGroup.Route = item.Route
				accessGroup.Icon = item.Icon
				accessGroups = append(accessGroups, accessGroup)
			}
		}

	}
	// helpers.PrintlnJson(accessGroups)
	return accessGroups, nil
}
