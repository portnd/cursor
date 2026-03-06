package repositories

import (
	"gitlab.com/mims-api-service/src/department/handlers"
	"gitlab.com/mims-api-service/src/department/usecases"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type departmentRepository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewDepartmentRepositoryHandler(conn *gorm.DB) *handlers.DepartmentHandler {
	useCase := usecases.NewDepartmentUseCase(&departmentRepository{conn})
	handler := handlers.NewDepartmentHandler(useCase)
	return handler
}

// //===================== query =====================
