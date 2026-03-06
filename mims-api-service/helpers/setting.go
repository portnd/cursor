package helpers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/models"
	"gorm.io/gorm"
)

func GetAssetTypeCondition(query *gorm.DB, assetType string) *gorm.DB {
	if assetType == "in" {
		return query.Where("is_in_road = ?", true)
	}

	return query.Where("is_in_road = ?", false)
}

type CreateAssetTable struct {
	QueryString   string
	AssetTable    models.RefAssetTable
	InsertColumns []models.RefAssetTableColumns
	Staffs        []models.RefAssetTableStaff
}

type UpdateAssetTable struct {
	CommentQuery       string
	RenameColumnsQuery string
	AddColumnsQuery    string
	AssetTable         models.RefAssetTable
	InsertNewColumns   []models.RefAssetTableColumns
	UpdateColumns      []models.RefAssetTableColumns
	Staffs             []models.RefAssetTableStaff
	DeleteColumns      []int
}

func GetLimitNumberOfColumns() int {
	limit, _ := strconv.Atoi(os.Getenv("LIMIT_COLUMNS"))
	return limit
}

func IsImageTypeJPEGOrPNG(fileName string) bool {
	isPNG := CheckImageType(fileName, "png")
	isJPEG := CheckImageType(fileName, "jpeg")
	isJPG := CheckImageType(fileName, "jpg")

	return isPNG || isJPEG || isJPG
}

func CheckImageType(fileName, fileType string) bool {
	return strings.ToLower(fileName[len(fileName)-len(fileType):]) == fileType
}

func GetHostAndScheme(c *gin.Context) string {
	if c.Request.URL.Scheme == "" {
		return fmt.Sprintf("http://%s/", c.Request.Host)
	}
	return fmt.Sprintf("https://%s/", c.Request.Host)
}

func MatchInteger(s []int, str int) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

type ReturnValueOfGetAssetTables struct {
	AssetTables      []models.AssetTable
	AssetTableStaffs []models.AssetTableStaff
	Limit            int64
	Page             int64
	Total            int64
}

type ReturnValueOfGetAssetTableById struct {
	AssetTable       models.AssetTable
	AssetTableStaffs []models.AssetTableStaff
}
