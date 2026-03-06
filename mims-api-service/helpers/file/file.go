package helpers

import (
	"fmt"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jonas-p/go-shp"
	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
)

type helpersFile struct {
}

func NewHelpersFile() *helpersFile {
	return &helpersFile{}
}

func (t *helpersFile) SaveFile(c *gin.Context, file *multipart.FileHeader, dir string) (string, error) {
	//add time to make image file unique
	dstPath := dir + time.Now().Format("20060102150405") + "_" + file.Filename
	err := c.SaveUploadedFile(file, dstPath)
	if err != nil {
		return "", err
	}

	return dstPath, nil
}

func (t *helpersFile) EnsureDir(dirName string) error {
	// Check if the directory exists
	info, err := os.Stat(dirName)

	// If the directory does not exist, create it
	if os.IsNotExist(err) {
		return os.MkdirAll(dirName, 0775) // You can adjust the file permissions as needed
	}

	// If the path exists but is not a directory, return an error
	if info != nil && !info.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", dirName)
	}

	return nil
}

func (s *helpersFile) ReadCenterLineShepesFile(centerLinePath string) (*string, *string, *string, error) {
	centerLineInZip, err := shp.ShapesInZip(centerLinePath)
	if err != nil {
		logs.Error(err)
		return nil, nil, nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LINE_FILE)
	}

	centerLineInZip = helpers.FilterByPrefix(centerLineInZip, "__MACOSX/")

	if len(centerLineInZip) > 1 {
		return nil, nil, nil, responses.NewAppErr(400, constants.HAS_MANY_SHAPE_FILE_IN_ZIP)
	}
	//Open Center Line Shape File From Zip
	shapeLine, err := shp.OpenShapeFromZip(centerLinePath, centerLineInZip[0])
	if err != nil {
		logs.Error(err)
		return nil, nil, nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LINE_FILE)

	}

	lineString := ""
	var kmStartStr, kmEndStr string
	// loop through all features in the shapefile
	fields := shapeLine.Fields()
	for shapeLine.Next() {
		_, p := shapeLine.Shape()

		// Check if the geometry is a PolyLine
		polyline, ok := p.(*shp.PolyLine)
		if !ok {
			return nil, nil, nil, responses.NewAppErr(400, constants.INVALID_POLY_LINE)
		}

		//Convert Center Line Polyline To LineString
		lineString = convertPolylineToLineString(polyline)
		for k, f := range fields {
			val := shapeLine.Attribute(k)
			fieldName := strings.TrimSpace(strings.TrimRight(string(f.Name[:]), "\x00"))

			var handled bool
			switch fieldName {
			case "road_code":
				handled = true
			case "name":
				handled = true
			case "km_start":
				handled = true
				kmStartStr = val

			case "km_end":
				handled = true
				kmEndStr = val

			}
			if !handled {
				return nil, nil, nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LINE_FILE)
			}
		}
	}
	return &lineString, &kmStartStr, &kmEndStr, nil
}

func (s *helpersFile) ReadCenterLaneShepesFile(centerLanePath string) ([]models.RoadGeom, error) {
	centerLaneInZip, err := shp.ShapesInZip(centerLanePath)
	if err != nil {
		return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LANE_FILE)

	}
	centerLaneInZip = helpers.FilterByPrefix(centerLaneInZip, "__MACOSX/")
	if len(centerLaneInZip) > 1 {
		return nil, responses.NewAppErr(400, constants.HAS_MANY_SHAPE_FILE_IN_ZIP)
	}

	shapeLane, err := shp.OpenShapeFromZip(centerLanePath, centerLaneInZip[0])
	if err != nil {
		fmt.Println(err)
		return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LANE_FILE)

	}

	var roadGeoms []models.RoadGeom
	// loop through all features in the shapefile
	for shapeLane.Next() {
		_, p := shapeLane.Shape()

		// Check if the geometry is a PolyLine
		polyline, ok := p.(*shp.PolyLine)
		if !ok {
			return nil, responses.NewAppErr(400, constants.INVALID_POLY_LINE)
		}

		lineString := convertPolylineToLineString(polyline)
		// shapeLaneGeom, err := t.roadRepo.ConvertLineStringToGeom(lineString)
		// if err != nil {
		// 	return 0, responses.NewAppErr(400, constants.INVALID_CONVERT_LINE_STRING_TO_GEOM)
		// }

		var roadGeom models.RoadGeom

		fields := shapeLane.Fields()
		for k, f := range fields {
			val := shapeLane.Attribute(k)
			fieldName := strings.TrimSpace(strings.TrimRight(string(f.Name[:]), "\x00"))
			// Convert [11]byte to string and trim any whitespace

			var handled bool
			switch fieldName {
			case "road_code":
				handled = true
			case "lane_no":
				laneNo, err := strconv.Atoi(val)
				if err == nil {
					roadGeom.LaneNo = laneNo
				}
				handled = true
			case "km_start":

				val := strings.Replace(val, "+", "", -1)
				kmStart, err := strconv.ParseFloat(val, 64)
				if err == nil {
					roadGeom.KmStart = kmStart
				}
				handled = true
			case "km_end":

				val := strings.Replace(val, "+", "", -1)
				kmEnd, err := strconv.ParseFloat(val, 64)
				if err == nil {
					roadGeom.KmEnd = kmEnd
				}
				handled = true

			}

			if !handled {
				return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LANE_FILE)
			}
		}

		roadGeom.TheGeom = lineString
		roadGeoms = append(roadGeoms, roadGeom)

	}
	return roadGeoms, nil
}

func convertPolylineToLineString(polyline *shp.PolyLine) string {
	var points []string
	for _, point := range polyline.Points {
		points = append(points, fmt.Sprintf("%f %f", point.X, point.Y))
	}
	return fmt.Sprintf(`LINESTRING(%s)`, strings.Join(points, ", "))
}
