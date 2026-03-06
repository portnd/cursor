package repositories

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/menu/handlers"
	"gitlab.com/mims-api-service/src/menu/usecases"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type menuRepository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewMenuRepositoryHandler(conn *gorm.DB) *handlers.MenuHandler {
	useCase := usecases.NewMenuUseCase(&menuRepository{conn})
	handler := handlers.NewMenuHandler(useCase)
	return handler
}

// //===================== query =====================
func (t *menuRepository) GetMenu() ([]models.AccessGroup, error) {
	var acctrolGrp []models.AccessGroup
	if err := t.conn.Where("status = ?", 1).Order("id asc").Find(&acctrolGrp).Error; err != nil {
		return acctrolGrp, err
	}
	return acctrolGrp, nil
}
