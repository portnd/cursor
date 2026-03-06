package domains

import (
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
)

type UseCaseMaintenance interface {
	GetMaintenanceDashboard(roadIDs []int, depotCodes []string, filter requests.MaintenanceDashboard) (interface{}, error)
	GetMaintenanceTableDashboard(roadIDs []int, depotCodes []string, filter requests.MaintenanceDashboard) (interface{}, error)
	GetMaintenanceMapDashboard(roadIDs []int, depotCodes []string, filter requests.MaintenanceDashboard) (interface{}, error)
}
type RepositoryMaintenance interface {
	GetRoadGroup() ([]models.RoadGroup, error)
	GetMaintenance(limit int, roadIDs []int, depotCodes []string, filter requests.MaintenanceDashboard) ([]models.MaintenancePreloadForDashboard, error)
}
