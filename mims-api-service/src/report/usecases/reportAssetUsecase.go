package usecases

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	"gorm.io/gorm"
)

func (u *UseCase) GetReport1(roadID, assetID, reportTyp string) (interface{}, error) {
	id, err := strconv.Atoi(roadID)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}

	roads, err := u.Repo.GetRoadInfo(id)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		logs.Error(err)
		return responses.RoadConditionDetails{}, err
	}

	data := []models.DataReportMap{}

	for _, road := range roads {
		road.RoadLengthStr = helpers.FormatNumberFloat(road.RoadSectionDistance)

		tableName, err := u.Repo.GetTableName(fmt.Sprint(road.RoadID), assetID)
		if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
			logs.Error(err)
			return responses.RoadConditionDetails{}, err
		}

		column, err := u.Repo.GetColumn(assetID, reportTyp)
		if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
			logs.Error(err)
			return responses.RoadConditionDetails{}, err
		}

		if tableName == nil || column == nil {
			assetName, err := u.Repo.GetAssetName(assetID)
			if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
				logs.Error(err)
				return responses.RoadConditionDetails{}, err
			}

			data = append(data, models.DataReportMap{
				AssetName:     assetName.AssetName,
				RoadGroupName: road.RoadGroupName,
				RoadName:      road.RoadName,
				RoadCode:      road.RoadCode,
				RoadColorCode: road.RoadColorCode,
				KmStart:       helpers.FormatKM(int64(road.KmStart)),
				KmEnd:         helpers.FormatKM(int64(road.KmEnd)),
				StrRoadLength: road.RoadLengthStr,

				IsNull: true,

				RoadID:                       road.RoadID,
				RoadGroupNumber:              road.RoadGroupNumber,
				RoadSectionNumber:            road.RoadSectionNumber,
				RoadSectionDistance:          road.RoadSectionDistance,
				RoadSectionNameOriginTh:      road.RoadSectionNameOriginTh,
				RoadSectionNameDestinationTh: road.RoadSectionNameDestinationTh,
			})
		} else {

			var columnName []string
			var columnTitle []string
			var key []string
			typ := make(map[string]string)
			var join string

			for _, i := range column {

				if i.ColumnName != "the_geom" {

					if strings.HasSuffix(i.ColumnName, "_image_id") {
						columnName = append(columnName, strings.TrimSuffix(i.ColumnName, "_id")+".sign_image_filepath AS "+i.ColumnName)
						join += " LEFT JOIN " + strings.TrimSuffix(i.ColumnName, "_id") + " on " + strings.TrimSuffix(i.ColumnName, "_id") + ".id = " + tableName.TableName + "." + i.ColumnName

					} else if i.ColumnName == "road_id" {
						columnName = append(columnName, strings.TrimSuffix(i.ColumnName, "_id")+".id AS "+i.ColumnName)
						join += " LEFT JOIN " + strings.TrimSuffix(i.ColumnName, "_id") + " on " + strings.TrimSuffix(i.ColumnName, "_id") + ".id = " + tableName.TableName + "." + i.ColumnName

					} else if strings.HasSuffix(i.ColumnName, "_id") {
						columnName = append(columnName, strings.TrimSuffix(i.ColumnName, "_id")+".name AS "+i.ColumnName)
						join += " LEFT JOIN " + strings.TrimSuffix(i.ColumnName, "_id") + " on " + strings.TrimSuffix(i.ColumnName, "_id") + ".id = " + tableName.TableName + "." + i.ColumnName

					} else {
						columnName = append(columnName, tableName.TableName+"."+i.ColumnName)
					}

					key = append(key, i.ColumnName)
					columnTitle = append(columnTitle, i.ComponentTitle)
					typ[i.ColumnName] = i.ComponentType
				}
			}

			// columnName = append(columnName, "road_asset.updated_date")
			// key = append(key, "updated_date")
			// columnTitle = append(columnTitle, "วันที่สำรวจ")

			rowMap, err := u.Repo.GetRow(columnName, fmt.Sprint(road.RoadID), assetID, tableName.TableName, join, reportTyp)
			if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
				logs.Error(err)
				return responses.RoadConditionDetails{}, err
			}

			if rowMap == nil {
				newRowMap := make(map[string]interface{})
				for _, i := range key {
					newRowMap[i] = nil
				}
				rowMap = append(rowMap, newRowMap)
			}

			var rows [][]interface{}

			for _, i := range rowMap {
				var row []interface{}

				for _, j := range key {

					if i[j] == nil {
						row = append(row, "")

					} else if typ[j] == "datepicker" {
						row = append(row, helpers.TimeThai(i[j].(time.Time)))

					} else if typ[j] == "text-km" {
						row = append(row, helpers.FormatKM(int64(i[j].(float64))))

					} else if j == "ref_asset_sign_image_id" {
						row = append(row, os.Getenv("STORAGE_IP")+"/"+fmt.Sprint(i[j]))

					} else if typ[j] == "text-number" {
						switch i[j].(type) {
						case float64:
							row = append(row, helpers.FormatNumberFloat(i[j].(float64)))
						case int:
							row = append(row, helpers.FormatNumber(i[j].(int)))
						}
					} else if typ[j] == "text-year" {
						// แปลงจาก string เป็น int
						year, err := strconv.Atoi(i[j].(string))
						if err != nil {
							fmt.Println("Error:", err)
							return nil, responses.NewAppErr(400, err.Error())
						}
						row = append(row, year+543)
					} else {
						ext := strings.ToLower(filepath.Ext(fmt.Sprint(i[j])))
						if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
							row = append(row, os.Getenv("STORAGE_IP")+"/"+fmt.Sprint(i[j]))

						} else {
							row = append(row, i[j])

						}
					}
				}

				rows = append(rows, row)
			}

			var mapGeom []string

			geomDB, err := u.Repo.GetMapGeom(fmt.Sprint(road.RoadID), assetID, tableName.TableName, reportTyp)
			if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
				logs.Error(err)
				return responses.RoadConditionDetails{}, err
			}

			for _, i := range geomDB {
				mapGeom = append(mapGeom, i.StringGeom)
			}

			newData := models.DataReportMap{
				AssetName:     tableName.TableLabel,
				RoadGroupName: road.RoadGroupName,
				RoadName:      road.RoadName,
				RoadCode:      road.RoadCode,
				RoadColorCode: road.RoadColorCode,
				KmStart:       helpers.FormatKM(int64(road.KmStart)),
				KmEnd:         helpers.FormatKM(int64(road.KmEnd)),
				StrRoadLength: road.RoadLengthStr,

				Column: columnTitle,
				Key:    key,

				Row: rows,

				RoadID:                       road.RoadID,
				RoadGroupNumber:              road.RoadGroupNumber,
				RoadSectionNumber:            road.RoadSectionNumber,
				RoadSectionDistance:          road.RoadSectionDistance,
				RoadSectionNameOriginTh:      road.RoadSectionNameOriginTh,
				RoadSectionNameDestinationTh: road.RoadSectionNameDestinationTh,
			}
			if newData.Row[0][0] == "" {
				newData.IsNull = true
			}

			if key[0] == "km" {
				newData.PointGeom = mapGeom
				newData.PinGeom = os.Getenv("STORAGE_IP") + "/public/assets/images/pin.png"
			} else {
				newData.LineGeom = mapGeom
			}

			data = append(data, newData)
		}
	}

	// generate report
	var pathResult interface{}

	if reportTyp == "excel" {
		pathResult, err = helpers.ExportExcelType1(data)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	} else {
		pathResult, err = helpers.RequestExport(data, "TEMPLATE_GENARAL_TYPE1", reportTyp)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	}
	//_ = pathResult
	return pathResult, nil
}
