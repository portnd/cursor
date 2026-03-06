package domains

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
)

// business logic
type MenuUseCase interface {
	GetMenu(map[string]string) ([]responses.AccessGroupMenu, error)
}

// อะไรเชื่อมต่อกับ DB
type MenuRepository interface {
	GetMenu() ([]models.AccessGroup, error)
}
