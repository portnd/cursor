package helpers

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/models"
)

type HelpersFileDomain interface {
	SaveFile(*gin.Context, *multipart.FileHeader, string) (string, error)
	EnsureDir(dirName string) error
	ReadCenterLineShepesFile(centerLinePath string) (*string, *string, *string, error)
	ReadCenterLaneShepesFile(centerLanePath string) ([]models.RoadGeom, error)
}
