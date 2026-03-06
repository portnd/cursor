package usecases

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/dashboard/domains"
	"gorm.io/gorm"
)

type useCase struct {
	repo domains.Repository
}

func NewUsecase(repo domains.Repository) domains.UseCase {
	return &useCase{repo: repo}
}

func (u *useCase) GetRoadDashboard() (interface{}, error) {

	var resp responses.RoadDashboard

	roadDashboards, err := u.repo.GetRoadDashboard()
	if err != nil {
		return nil, responses.NewAppErr(500, err.Error())
	}

	for i, roadDashboard := range roadDashboards {

		resp.RoadLabel.Name = "จำนวนถนน (สาย)"
		resp.RoadLabel.Data = append(resp.RoadLabel.Data, roadDashboard.TotalRoad)
		resp.RoadLabel.Label = append(resp.RoadLabel.Label, "หมายเลข "+roadDashboard.RoadGroupName)
		resp.RoadLabel.RoadGroupID = append(resp.RoadLabel.RoadGroupID, roadDashboard.ID)

		if i < len(Colors) {
			resp.RoadLabel.Color = append(resp.RoadLabel.Color, Colors[i])
		}

		if resp.RoadLabel.Label == nil {
			resp.RoadLabel.Label = []string{}
		}

		if resp.RoadLabel.Data == nil {
			resp.RoadLabel.Data = []int{}
		}
		if resp.RoadLabel.Color == nil {
			resp.RoadLabel.Color = []string{}
		}

		pavementSurfaceRoads, err := u.repo.GetPavementSurfaceByRoadGroupID(roadDashboard.ID)
		if err != nil {
			return nil, responses.NewAppErr(500, err.Error())
		}

		var respLength responses.LengthRoad
		respLength.Name = "ความยาวทางหลวงพิเศษหมายเลข " + roadDashboard.RoadGroupName + " ทั้งหมด"

		if i < len(Colors) {
			respLength.Color = Colors[i]
		}

		var lengthRoad float64
		var lengthSurface float64
		var lengthSurfaceCC float64
		var lengthSurfaceAC float64
		for _, v := range pavementSurfaceRoads {
			lengthRoad = lengthRoad + v.Length
			for _, sur := range v.SurfaceData {
				for _, surL := range sur.SurfaceLane {
					lengthSurface = lengthSurface + sur.Length
					if surL.SurfaceTyepe == "Concrete" {
						lengthSurfaceCC = lengthSurfaceCC + sur.Length
					} else if surL.SurfaceTyepe == "Asphalt" {
						lengthSurfaceAC = lengthSurfaceAC + sur.Length
					}
				}
			}
		}
		respLength.Total = helpers.RoundFloat(lengthRoad, 2)
		respLength.Asphalt = helpers.RoundFloat(((lengthSurfaceAC / lengthSurface) * 100), 2)
		respLength.Concrete = helpers.RoundFloat(((lengthSurfaceCC / lengthSurface) * 100), 2)

		resp.LengthRoad = append(resp.LengthRoad, respLength)

		volumeAADTRoads, err := u.repo.GetVolumeAADTByRoadGroupID(roadDashboard.ID)
		if err != nil {
			return nil, responses.NewAppErr(500, err.Error())
		}

		time, err := helpers.GetBangkokTimeNow()
		if err != nil {
			return nil, responses.NewAppErr(500, err.Error())
		}

		var currentYearVolume, lastYearVolume float64
		var currentYear, lastYear int

		var respAadtRoad responses.AadtRoad
		respAadtRoad.Name = "ปริมาณการจราจรทางหลวงพิเศษหมายเลข " + roadDashboard.RoadGroupName
		if len(volumeAADTRoads) > 1 {

			for _, volumeAADTRoad := range volumeAADTRoads {

				if volumeAADTRoad.Year == time.Year() {
					currentYear = volumeAADTRoad.Year //ปีปัจจุบัน
					currentYearVolume = float64(volumeAADTRoad.Total)

				} else {
					lastYear = volumeAADTRoad.Year //ปีก่อนหน้า
					lastYearVolume = float64(volumeAADTRoad.Total)
				}

				respAadtRoad.Aadt = int(currentYearVolume)

				// Calculate the percentage difference
				percentDiff := ((currentYearVolume - lastYearVolume) / lastYearVolume) * 100
				respAadtRoad.Percent = helpers.RoundFloat(percentDiff, 2)
				if percentDiff > 0 {
					respAadtRoad.GrowthRate = "up"
				} else {
					respAadtRoad.GrowthRate = "down"
				}

			}

		}
		if currentYear != 0 {
			respAadtRoad.Year2 = currentYear
		} else {
			respAadtRoad.Year2 = time.Year()
		}
		if lastYear != 0 {
			respAadtRoad.Year1 = lastYear
		} else {
			respAadtRoad.Year1 = time.Year() - 1
		}

		resp.AadtRoad = append(resp.AadtRoad, respAadtRoad)
	}

	return resp, nil
}

// func (u *useCase) GetAsset(roadIDs []int, depotCodes []string, filter requests.Asset) (interface{}, int64, error) {

// 	var resps []responses.DashboardAsset
// 	results, total, err := u.repo.GetAsset(roadIDs, depotCodes, filter)
// 	if err != nil {
// 		return results, 0, responses.NewAppErr(500, err.Error())
// 	}

// 	for _, result := range results {
// 		var resp responses.DashboardAsset

// 		if len(result.RefAssetTables) == 0 {
// 			continue
// 		}

// 		countByRoadgroup := make(map[string]int)

// 		resp.Label = []string{}
// 		resp.Data = []int{}

// 		for _, refAssetTable := range result.RefAssetTables {

// 			fmt.Println("refAssetTableID", refAssetTable.ID)
// 			fmt.Println("RefAssetID", refAssetTable.RefAssetID)
// 			resp.RefAssetID = refAssetTable.RefAssetID
// 			resp.ID = refAssetTable.RefAssetTable.ID
// 			resp.Name = refAssetTable.RefAssetTable.TableLabel

// 			for _, roadAsset := range refAssetTable.RoadAssets {

// 				if len(depotCodes) > 0 && roadAsset.Road.RoadSection.Id == 0 {
// 					continue
// 				}

// 				if (filter.Year != "" || filter.KmStart >= 0.0 || filter.KmEnd >= 0.0) && roadAsset.RoadInfo.Id == 0 {
// 					continue
// 				}

// 				countByRoadgroup[roadAsset.Road.RoadGroup.Number] += roadAsset.AssetCount

// 			}

// 			// for groupNumber, count := range countByRoadgroup {
// 			// 	if count == 0 {
// 			// 		continue
// 			// 	}

// 			// 	resp.Lable = append(resp.Lable, groupNumber)
// 			// 	resp.Data = append(resp.Data, count)
// 			// 	//resps = append(resps, resp)
// 			// }

// 			addedLabels := make(map[string]bool)

// 			for groupNumber, count := range countByRoadgroup {
// 				if count == 0 {
// 					continue
// 				}

// 				if !addedLabels[groupNumber] {
// 					resp.Label = append(resp.Label, groupNumber)
// 					resp.Data = append(resp.Data, count)
// 					addedLabels[groupNumber] = true
// 				}
// 			}

// 			resps = append(resps, resp)

// 		}

// 	}

// 	sort.Slice(resps, func(i, j int) bool {
// 		return resps[i].ID < resps[j].ID
// 	})

// 	sort.Slice(resps, func(i, j int) bool {
// 		return resps[i].RefAssetID < resps[j].RefAssetID
// 	})

// 	return resps, total, nil
// }

func (u *useCase) GetAsset(roadIDs []int, depotCodes []string, filter requests.Asset) (interface{}, int64, error) {

	var resps []responses.DashboardAsset
	results, total, err := u.repo.GetAsset(roadIDs, depotCodes, filter)
	if err != nil {
		return results, 0, responses.NewAppErr(500, err.Error())
	}

	for _, result := range results {
		var resp responses.DashboardAsset

		if len(result.RefAssetTables) == 0 {
			continue
		}

		for _, refAssetTable := range result.RefAssetTables {
			resp.Label = []string{}
			resp.Data = []int{}

			countByRoadgroup := make(map[int]map[string]int)
			addedLabels := make(map[string]bool)

			resp.RefAssetID = refAssetTable.RefAssetID
			resp.ID = refAssetTable.RefAssetTable.ID
			resp.Name = refAssetTable.RefAssetTable.TableLabel

			for _, roadAsset := range refAssetTable.RoadAssets {

				if len(depotCodes) > 0 && roadAsset.Road.RoadSection.Id == 0 {
					continue
				}

				if (filter.Year != "" || filter.KmStart >= 0.0 || filter.KmEnd >= 0.0) && roadAsset.RoadInfo.Id == 0 {
					continue
				}

				if countByRoadgroup[roadAsset.RefAssetTableId] == nil {
					countByRoadgroup[roadAsset.RefAssetTableId] = make(map[string]int)
				}

				countByRoadgroup[roadAsset.RefAssetTableId][roadAsset.Road.RoadGroup.Number] += roadAsset.AssetCount
			}

			type groupData struct {
				GroupNumber string
				Count       int
			}

			var groupDataList []groupData

			for groupNumber, count := range countByRoadgroup[refAssetTable.ID] {
				if count > 0 {
					groupDataList = append(groupDataList, groupData{GroupNumber: groupNumber, Count: count})
				}
			}

			// for groupNumber, count := range countByRoadgroup {
			// 	if count == 0 {
			// 		continue
			// 	}

			// 	if !addedLabels[groupNumber] {
			// 		resp.Label = append(resp.Label, groupNumber)
			// 		resp.Data = append(resp.Data, count)
			// 		addedLabels[groupNumber] = true
			// 	}
			// }
			// Sort the collected data
			sort.Slice(groupDataList, func(i, j int) bool {
				return groupDataList[i].GroupNumber < groupDataList[j].GroupNumber
			})

			// Append the sorted data to resp
			for _, data := range groupDataList {
				if !addedLabels[data.GroupNumber] {
					resp.Label = append(resp.Label, data.GroupNumber)
					resp.Data = append(resp.Data, data.Count)
					addedLabels[data.GroupNumber] = true
				}
			}

			resps = append(resps, resp)
		}

	}

	sort.Slice(resps, func(i, j int) bool {
		return resps[i].ID < resps[j].ID
	})

	sort.Slice(resps, func(i, j int) bool {
		return resps[i].RefAssetID < resps[j].RefAssetID
	})

	return resps, total, nil
}

func genAssetList(datas []models.RefAssetTableDashboard, roadIDs []int) []responses.AssetList {
	assetLists := make([]responses.AssetList, 0)
	STORAGE_IP := os.Getenv("STORAGE_IP") + "/"

	for _, data := range datas {
		iconFilepath := ""
		if data.IconFilepath != "" {
			iconFilepath = STORAGE_IP + data.IconFilepath
		}
		var assetList responses.AssetList
		assetList.Asset.ID = data.ID
		assetList.Asset.Name = data.TableLabel
		assetList.Asset.DefaultIconURL = iconFilepath
		assetList.Asset.ThumbnailIconURL = iconFilepath
		assetList.Asset.DefaultColor = data.LineColorCode
		if data.GeomType == 2 {
			assetList.IsRange = true
		} else {
			assetList.IsRange = false
		}
		if len(data.RoadAssets) == 0 {
			assetList.Value = 0
		} else {
			assetList.Value = CountValue(data.RoadAssets, roadIDs)
		}
		assetLists = append(assetLists, assetList)
	}

	return assetLists
}

func CountValue(data []models.RoadAssetForDashboard, roadIDs []int) int {
	var value int
	for _, v := range data {
		if helpers.ContainsInt(v.RoadId, roadIDs) || len(roadIDs) == 0 {
			value += v.AssetCount
		}

	}
	return value
}

func (u *useCase) GetAssetDetail(roadIDs []int, depotCodes []string, filter requests.Asset) (interface{}, int64, error) {

	var responds []responses.AssetRespond

	result, total, err := u.repo.GetAsset(roadIDs, depotCodes, filter)
	if err != nil {
		return result, 0, responses.NewAppErr(500, err.Error())
	}
	for _, data := range result {
		if len(data.RefAssetTables) == 0 {
			continue
		}
		var respond responses.AssetRespond
		respond.AssetGroup.ID = data.ID
		respond.AssetGroup.Name = data.Name
		refTables := genAssetList(data.RefAssetTables, roadIDs)
		respond.AssetList = append(respond.AssetList, refTables...)
		if len(refTables) == 0 {
			respond.AssetList = []responses.AssetList{}
		}
		responds = append(responds, respond)
	}

	for i := range responds {
		sort.Slice(responds[i].AssetList, func(a, b int) bool {
			return responds[i].AssetList[a].Asset.ID < responds[i].AssetList[b].Asset.ID
		})
	}

	sort.Slice(responds, func(i, j int) bool {
		return responds[i].AssetGroup.ID < responds[j].AssetGroup.ID
	})

	return responds, total, nil
}

func (u *useCase) GetAssetMap(roadIDs []string, assetIDs []string, depotCodes []string, filter requests.AssetMap) (interface{}, error) {
	var resps []responses.AssetMap

	if len(assetIDs) == 0 {
		return []responses.AssetMap{}, nil
	}

	listTable, err := u.repo.GetTableResult(roadIDs, assetIDs, depotCodes, filter)
	if err != nil {
		return nil, responses.NewAppErr(500, err.Error())
	}

	zoom := 1
	if filter.Zoom != "" {
		zoom, err = strconv.Atoi(filter.Zoom)
		if err != nil {
			return nil, responses.NewAppErr(500, err.Error())
		}
	}

	if zoom <= 0 {
		return "", responses.NewAppErr(400, "Zoom level must be greater than 0")
	}

	tables := []string{}
	var isZoom12 bool
	for _, t := range listTable {

		if zoom <= 12 {
			isZoom12 = true
			tables = append(tables, t.TableName)
		} else {

			roadSigns, err := u.repo.GetRoadAssetSignDataByRoadIDs(t, roadIDs, filter)
			if err != nil {
				return nil, responses.NewAppErr(500, err.Error())
			}

			createAssetLocation := CreateAssetLocation(t, roadSigns)

			for _, v := range createAssetLocation {
				var resp responses.AssetMap

				resp.ID = v.ID
				resp.RoadID = v.RoadID
				resp.AssetTableID = t.ID
				resp.Name = v.Name
				resp.IconFilepath = v.IconFilepath
				resp.ThumbnailIconFilepath = v.ThumbnailIconFilepath
				resp.LineColorCode = v.LineColorCode

				geomJson, err := helpers.ConvertThegeomToGeomJSON(v.Wkt)
				if err != nil {
					return nil, responses.NewAppErr(500, err.Error())
				}
				resp.TheGeom = geomJson
				resps = append(resps, resp)

			}

		}

	}

	if isZoom12 {

		buildQuery := buildUnionQuery(tables, filter)
		theGeomCusters, err := u.repo.GetRoadAssetTheGeomCuster(buildQuery, roadIDs, assetIDs, depotCodes, filter)
		if err != nil {
			return nil, responses.NewAppErr(500, err.Error())
		}

		for _, theGeoms := range theGeomCusters {
			var resp responses.AssetMap

			if len(theGeoms.TheGeomCluster) == 0 {
				continue
			}

			theGeomCuster, err := helpers.ConvertThegeomToGeomJSON(theGeoms.TheGeomCluster)
			if err != nil {
				return nil, responses.NewAppErr(500, err.Error())
			}

			resp.IsCluster = true
			resp.Cluster = theGeoms.TotalTheGeomCluster
			resp.TheGeom = theGeomCuster

			resps = append(resps, resp)

		}

	}

	return resps, nil

}

func buildUnionQuery(tables []string, filter requests.AssetMap) string {
	var parts []string

	for _, table := range tables {
		parts = append(parts, generateSQLForTable(table, filter))
	}
	return strings.Join(parts, " UNION ALL ")
}

func generateSQLForTable(tableName string, filter requests.AssetMap) string {
	// Assume tableName is validated or fetched from a controlled source to prevent SQL Injection

	return fmt.Sprintf(`
        SELECT t.the_geom FROM %s t
        JOIN road_asset ra ON ra.id = t.road_asset_id
		Left Join road r ON ra.road_id = r.id 
		Left Join road_info ri ON ra.road_id = ri.road_id 
		Left Join road_section rs ON r.road_section_id = rs.id 
        WHERE ra.status = 'A' AND t.is_deleted = false AND ST_Transform(t.the_geom, 4326) && ST_MakeEnvelope(%s, %s, %s, %s, 4326) group by t.id 
    `, tableName, filter.Left, filter.Bottom, filter.Right, filter.Top)
}

func CreateAssetLocation(t models.TableResult, roadSign []models.RoadAssetSign) []models.AssetLocation {
	var results []models.AssetLocation
	STORAGE_IP := os.Getenv("STORAGE_IP") + "/"
	for _, g := range roadSign {
		var assLo models.AssetLocation
		ThumbnailIconFilepath := t.IconFilepath
		if g.SignImageFilePath != "" {
			ThumbnailIconFilepath = g.SignImageFilePath
		} else if g.ImgFilePath != "" {
			ThumbnailIconFilepath = g.ImgFilePath
		}
		//type_var := strings.Split(g.GeomCl, "(")[0]
		assLo.ID = g.ID
		assLo.Name = t.TableLabel
		//assLo.Type = type_var
		assLo.RoadID = g.RoadID
		if t.IconFilepath != "" {
			assLo.IconFilepath = STORAGE_IP + t.IconFilepath
		}

		if ThumbnailIconFilepath != "" {
			assLo.ThumbnailIconFilepath = STORAGE_IP + ThumbnailIconFilepath
		}

		assLo.LineColorCode = t.LineColorCode

		assLo.Wkt = g.TheGeom
		results = append(results, assLo)
	}
	return results
}

func (u *useCase) GetAssetMapDetailByID(ID, assetTableID int) (interface{}, error) {

	colList, err := u.repo.GetRefAssetTableColumns(assetTableID)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewAppErr(500, err.Error())
	}

	result, err := u.repo.GetRoadAssetDetails(ID, colList)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewAppErr(500, err.Error())
	}

	htmlDetils, err := u.ConvertToHtml(colList, result)
	if err != nil {
		return nil, responses.NewAppErr(500, err.Error())
	}

	return htmlDetils, nil
}

func (u *useCase) ConvertToHtml(columns []responses.RoadAssetDetailColumn, values []map[string]interface{}) (string, error) {

	lines := []string{}

	for _, val := range values {
		for _, column := range columns {
			line := ""

			value := val[column.ColumnName]

			if value == nil {
				line = fmt.Sprintf("<p>%s: -</p>", column.ComponentTitle)
			} else if column.ComponentType == "select" {
				// Handle 'select' type component
				if id, ok := value.(int64); ok {
					name, err := u.repo.GetRefDataFromSelect(column.TableNameRef, id)
					if err != nil {
						return "", err
					}
					line = fmt.Sprintf("<p>%s: %s</p>", column.ComponentTitle, name)
				}
			} else if column.ComponentType == "text-km" {

				if val, ok := value.(float64); ok {

					kmString := helpers.FormatKM(int64(val))

					line = fmt.Sprintf("<p>%s: %s</p>", column.ComponentTitle, kmString)
				}

			} else {

				switch v := value.(type) {
				case float64:
					if column.ColumnName == "road_code" {
						line = fmt.Sprintf("<p>%s: %v</p>", column.ComponentTitle, v)
					} else if column.ColumnName == "section_code" {
						line = fmt.Sprintf("<p>%s: %v</p>", column.ComponentTitle, v)
					} else {
						// Format floating-point numbers with '%f' or '%g'
						line = fmt.Sprintf("<p>%s: %.2f</p>", column.ComponentTitle, v)
					}
				case int, int64:
					// Integers can be formatted with '%d'
					line = fmt.Sprintf("<p>%s: %d</p>", column.ComponentTitle, v)
				case string:
					// Strings can continue using '%s'
					line = fmt.Sprintf("<p>%s: %s</p>", column.ComponentTitle, v)
				default:
					// Fallback for other data types
					if column.ComponentType == "datepicker" {
						// พยายามแปลงเป็น time.Time
						if timeValue, ok := v.(time.Time); ok {
							// เช็คว่าเป็นวันที่ว่างไหม
							if timeValue.IsZero() {
								line = fmt.Sprintf("<p>%s: -</p>", column.ComponentTitle)
							} else {
								// ถ้าไม่ว่าง แปลงเป็นวันที่ภาษาไทย
								date := helpers.ConvertToThaiFullCalendarNoTime(timeValue)
								line = fmt.Sprintf("<p>%s: %v</p>", column.ComponentTitle, date)
							}
						} else {
							// กรณีแปลงเป็น time.Time ไม่ได้
							line = fmt.Sprintf("<p>%s: -</p>", column.ComponentTitle)
						}
					} else {
						// กรณีไม่ใช่ datepicker
						line = fmt.Sprintf("<p>%s: %v</p>", column.ComponentTitle, v)
					}
				}

			}
			lines = append(lines, line)

		}

	}

	html := strings.Join(lines, "")
	return html, nil

}

var Colors = map[int]string{
	0:  constants.YELLOW_1,
	1:  constants.YELLOW_2,
	2:  constants.YELLOW_3,
	3:  constants.ORANGE_1,
	4:  constants.ORANGE_2,
	5:  constants.ORANGE_3,
	6:  constants.BRIGHT_YELLOW_1,
	7:  constants.BRIGHT_YELLOW_2,
	8:  constants.BRIGHT_YELLOW_3,
	9:  constants.RED_ORANGE_1,
	10: constants.RED_ORANGE_2,
	11: constants.RED_ORANGE_3,
}

func (u *useCase) GetYears(typeName string) (interface{}, error) {

	var table string
	var getAll bool

	switch typeName {
	case "road":
		table = "road_info"
	case "condition":
		table = "road_condition"
	case "surface":
		table = "road_surface"
	case "maintenance":
		table = "maintenance"
	case "":
		getAll = true
	default:
		return nil, responses.NewAppErr(400, "Invalid type name")

	}

	var resp []int
	resp = []int{}

	if strings.TrimSpace(typeName) != "" {

		years, err := u.repo.GetDashboardYear(table)
		if err != nil {
			return nil, responses.NewAppErr(500, err.Error())
		}

		resp = years

	}

	if getAll {

		maxminYear, err := u.repo.GetDashboardMaxMinYear()
		if err != nil {
			return nil, responses.NewAppErr(500, err.Error())
		}

		if maxminYear.MinYear == 0 && maxminYear.MaxYear == 0 {
			return resp, nil
		}

		for i := maxminYear.MinYear; i <= maxminYear.MaxYear; i++ {

			resp = append(resp, i)
		}

	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i] > resp[j]
	})

	return resp, nil
}

func (u *useCase) GetRoadConditionList(roadId int) (interface{}, error) {

	rcStatus := "rc.status = 'A'"

	result, err := u.repo.GetRoadConditionList(rcStatus, roadId)
	if err != nil {
		return nil, err
	}

	data := make(map[int][]responses.RoadConditionListItem)
	for _, r := range result {
		year := r.Year
		item := responses.RoadConditionListItem{
			ID:       r.ID,
			IDParent: r.IDParent,
			Revision: r.Revision,
			Direction: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{ID: r.DirectionId, Name: r.DirectionName},
			LaneNo:       r.LaneNo,
			SurveyedDate: (r.SurveyedDate),
		}

		if data[year] == nil {
			data[year] = []responses.RoadConditionListItem{item}
		} else {
			found := false
			for i, existingItem := range data[year] {
				if existingItem.IDParent == item.IDParent {
					found = true
					if existingItem.Revision < item.Revision {
						data[year][i] = item
					}
					break
				}
			}
			if !found {
				data[year] = append(data[year], item)
			}
		}
	}

	var results []responses.RoadConditionList
	for year, items := range data {

		sort.Slice(items, func(i, j int) bool {
			return items[i].LaneNo < items[j].LaneNo
		})

		results = append(results, responses.RoadConditionList{Year: year, Items: items})
	}

	if len(results) == 0 {
		return []responses.RoadConditionList{}, nil
	}

	return results, nil
}

func (u *useCase) GetRoadConditionDetails(conditionRangeType string, idParent int) (interface{}, error) {

	var resp responses.RoadConditionDetails

	surveys, err := u.repo.GetAllRoadConditionByIdParent(idParent)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil // or return an appropriate empty response
		}

		return nil, responses.NewAppErr(500, err.Error())
	}

	road, err := u.repo.GetRoadByID(surveys.RoadId)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil // or return an appropriate empty response
		}

		return nil, responses.NewAppErr(500, err.Error())
	}

	conditionTypesMap := make(map[string][]responses.RoadConditionDetailHeaderItem)

	for _, survey := range surveys.RoadConditionSurveys {

		// Process conditions like IRI, MPD, RUT, IFI, etc.
		for _, conditionType := range []string{"IRI", "MPD", "RUT", "IFI"} {

			conditionValue := getConditionValue(conditionType, &survey)
			if conditionValue == nil {
				continue // Skip if the condition value is nil
			}

			geom := ""
			if survey.TheGeom != "" {
				geom, err = helpers.StringToGeom(survey.TheGeom)
				if err != nil {
					logs.Error(err)
					return nil, err
				}
			}

			// Generate header item for this condition type
			headerItem := responses.RoadConditionDetailHeaderItem{
				KmStart:    int(survey.KmStart),
				KmEnd:      int(survey.KmEnd),
				Value:      conditionValue, // Example, you might want to average or summarize
				GeomCl:     geom,
				SurveyType: survey.SurveyType, // Example, adjust based on your data
				Items:      []responses.RoadConditionDetailBodyItem{},
			}

			if conditionRangeType == "2" {
				// Generate body items for associated RoadConditionSurvey100Ms, if any
				for _, survey100M := range survey.RoadConditionSurvey100Ms {

					geom := ""
					if survey100M.TheGeom != "" {
						geom, err = helpers.StringToGeom(survey100M.TheGeom)
						if err != nil {
							logs.Error(err)
							return nil, err
						}
					}

					detailValue := getConditionValue(conditionType, &survey100M)
					if detailValue != nil {
						bodyItem := responses.RoadConditionDetailBodyItem{
							KmStart:     int(survey100M.KmStart),
							KmEnd:       int(survey100M.KmEnd),
							Value:       detailValue,
							GeomCl:      geom,                  // Assuming similar fields exist in detail
							ImgFilepath: "",                    // Assume you have this data
							SurveyType:  survey100M.SurveyType, // Example
						}
						headerItem.Items = append(headerItem.Items, bodyItem)
					}
				}
			}

			if conditionRangeType == "1" {
				// Generate body items for associated RoadConditionSurveyMs, if any
				for _, survey100M := range survey.RoadConditionSurvey100Ms {

					for _, surveyM := range survey100M.RoadConditionSurveyMs {

						geom := ""
						if surveyM.TheGeom != "" {
							geom, err = helpers.StringToGeom(surveyM.TheGeom)
							if err != nil {
								logs.Error(err)
								return nil, err
							}
						}

						detailValue := getConditionValue(conditionType, &surveyM)

						img := ""
						if surveyM.ImgFilepath != "" {
							img = os.Getenv("STORAGE_IP") + "/" + surveyM.ImgFilepath
						}
						if detailValue != nil {
							bodyItem := responses.RoadConditionDetailBodyItem{
								KmStart:     int(surveyM.KmStart),
								KmEnd:       int(surveyM.KmEnd),
								Value:       detailValue,
								GeomCl:      geom,
								ImgFilepath: img,
								SurveyType:  surveyM.SurveyType,
							}

							headerItem.Items = append(headerItem.Items, bodyItem)
						}
					}
				}
			}

			switch road.RoadInfo.RefDirectionId {
			case 1:
				sort.Slice(headerItem.Items, func(i, j int) bool {

					return headerItem.Items[i].KmStart < headerItem.Items[j].KmStart
				})

			case 2:
				sort.Slice(headerItem.Items, func(i, j int) bool {
					return headerItem.Items[i].KmStart > headerItem.Items[j].KmStart
				})
			}

			// Append the header item with all its body items to the response
			//headerItems = append(headerItems, headerItem)

			conditionTypesMap[conditionType] = append(conditionTypesMap[conditionType], headerItem)

		}
	}

	for conditionType, items := range conditionTypesMap {

		switch road.RoadInfo.RefDirectionId {
		case 1:
			sort.Slice(items, func(i, j int) bool {
				return items[i].KmStart < items[j].KmStart
			})

		case 2:
			sort.Slice(items, func(i, j int) bool {
				return items[i].KmStart > items[j].KmStart
			})
		}

		rcType := responses.RoadConditionDetailData{
			ConditionType:                 conditionType,
			RoadConditionDetailHeaderItem: items,
		}
		resp.ConditionTypes = append(resp.ConditionTypes, rcType)

		sort.Slice(resp.ConditionTypes, func(i, j int) bool {
			return resp.ConditionTypes[i].ConditionType < resp.ConditionTypes[j].ConditionType
		})
	}

	resp.ID = surveys.ID
	resp.IDParent = surveys.IDParent
	resp.Status = surveys.Status
	resp.UpdatedDate = helpers.SetTimeToString(surveys.UpdatedDate)
	resp.Direction = road.RoadInfo.Direction

	userInfo, err := u.repo.GetUserDepartmentById(surveys.UpdatedBy)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logs.Error(err)
			return responses.RoadConditionDetails{}, err
		}
	}

	resp.UpdatedBy.UID = fmt.Sprintf("%d", userInfo.Id)
	resp.UpdatedBy.UserName = userInfo.Username
	resp.UpdatedBy.FullName = userInfo.Firstname + " " + userInfo.Lastname
	// resp.UpdatedBy.Department.ID = userInfo.Department.ID
	// resp.UpdatedBy.Department.Name = userInfo.Department.Name
	if userInfo.ProfileImgPath != "" {
		resp.UpdatedBy.ProfilePicture = os.Getenv("STORAGE_IP") + "/" + userInfo.ProfileImgPath
	}

	return resp, nil
}

func (u *useCase) GetRoadRetroReflectivityList(roadID int) (interface{}, error) {

	result, err := u.repo.GetRoadRetroReflectivityList(roadID)
	if err != nil {
		return nil, err
	}

	data := make(map[int][]responses.RoadRetroReflectivityListItem)
	for _, r := range result {
		year := r.Year
		item := responses.RoadRetroReflectivityListItem{
			ID:       r.ID,
			IDParent: r.IDParent,
			Revision: r.Revision,
			Direction: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{ID: r.DirectionId, Name: r.DirectionName},
			LineNo:       r.LineNo,
			SurveyedDate: (r.SurveyedDate),
		}

		if data[year] == nil {
			data[year] = []responses.RoadRetroReflectivityListItem{item}
		} else {
			found := false
			for i, existingItem := range data[year] {
				if existingItem.IDParent == item.IDParent {
					found = true
					if existingItem.Revision < item.Revision {
						data[year][i] = item
					}
					break
				}
			}
			if !found {
				data[year] = append(data[year], item)
			}
		}
	}

	var results []responses.RoadRetroReflectivityList
	for year, items := range data {

		sort.Slice(items, func(i, j int) bool {
			return items[i].LineNo < items[j].LineNo
		})

		results = append(results, responses.RoadRetroReflectivityList{Year: year, Items: items})
	}

	if len(results) == 0 {
		return []responses.RoadRetroReflectivityList{}, nil
	}

	return results, nil
}
func (u *useCase) GetRoadRetroReflectivityDetails(rangeType string, refStripeTypeIDs string, idParent int) (interface{}, error) {

	var resp responses.RetroReflectivityDetails
	var stripeTypes []string
	if refStripeTypeIDs != "" {
		stripeTypes = strings.Split(refStripeTypeIDs, ",")
		fmt.Println(stripeTypes)
	}

	retroReflectivitys, err := u.repo.GetRoadRetroReflectivityDetailsByIdParent(idParent, stripeTypes)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil // or return an appropriate empty response
		}

		return nil, responses.NewAppErr(500, err.Error())
	}

	road, err := u.repo.GetRoadByID(retroReflectivitys.RoadID)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil // or return an appropriate empty response
		}

		return nil, responses.NewAppErr(500, err.Error())
	}

	RetroReflectivityTypesMap := make(map[string][]responses.RoadRetroReflectivityDetailHeaderItem)

	for _, retroReflectivityRange := range retroReflectivitys.RoadRetroReflectivityRanges {

		geom := ""
		if retroReflectivityRange.TheGeom != "" {
			geom, err = helpers.StringToGeom(retroReflectivityRange.TheGeom)
			if err != nil {
				logs.Error(err)
				return nil, err
			}
		}

		// Generate header item for this condition type
		headerItem := responses.RoadRetroReflectivityDetailHeaderItem{
			KmStart:          int(retroReflectivityRange.KmStart),
			KmEnd:            int(retroReflectivityRange.KmEnd),
			RetroAvg:         retroReflectivityRange.RetroAvg, // Example, you might want to average or summarize
			GeomCl:           geom,
			RefStripeTypeID:  retroReflectivityRange.RefStripeTypeID,
			RefStripeType:    retroReflectivityRange.RefStripeType,
			RefStripeColorID: retroReflectivityRange.RefStripeColorID,
			RefStripeColor:   retroReflectivityRange.RefStripeColor,
			Items:            []responses.RoadRetroReflectivityDetailBodyItem{},
		}

		if rangeType == "1" {
			// Generate body items for associated RoadConditionSurvey100Ms, if any
			for _, retroReflectivityM := range retroReflectivityRange.RoadRetroReflectivityMs {

				geom := ""
				if retroReflectivityM.TheGeom != "" {
					geom, err = helpers.StringToGeom(retroReflectivityM.TheGeom)
					if err != nil {
						logs.Error(err)
						return nil, err
					}
				}

				bodyItem := responses.RoadRetroReflectivityDetailBodyItem{
					KmStart:          int(retroReflectivityM.KmStart),
					KmEnd:            int(retroReflectivityM.KmEnd),
					RetroAvg:         retroReflectivityM.RetroAvg,
					GeomCl:           geom,
					RefStripeTypeID:  retroReflectivityRange.RefStripeTypeID,
					RefStripeType:    retroReflectivityRange.RefStripeType,
					RefStripeColorID: retroReflectivityRange.RefStripeColorID,
					RefStripeColor:   retroReflectivityRange.RefStripeColor,
				}
				headerItem.Items = append(headerItem.Items, bodyItem)

			}
		}

		switch road.RoadInfo.RefDirectionId {
		case 1:
			sort.Slice(headerItem.Items, func(i, j int) bool {

				return headerItem.Items[i].KmStart < headerItem.Items[j].KmStart
			})

		case 2:
			sort.Slice(headerItem.Items, func(i, j int) bool {
				return headerItem.Items[i].KmStart > headerItem.Items[j].KmStart
			})
		}

		// Append the header item with all its body items to the response
		//headerItems = append(headerItems, headerItem)

		RetroReflectivityTypesMap[retroReflectivityRange.RefStripeColor.Name] = append(RetroReflectivityTypesMap[retroReflectivityRange.RefStripeColor.Name], headerItem)

	}

	for retroReflectivityType, items := range RetroReflectivityTypesMap {

		switch retroReflectivityType {
		case "White":
			resp.HasWhiteLine = true
		case "Yellow":
			resp.HasYellowLine = true
		}

		switch road.RoadInfo.RefDirectionId {
		case 1:
			sort.Slice(items, func(i, j int) bool {
				return items[i].KmStart < items[j].KmStart
			})

		case 2:
			sort.Slice(items, func(i, j int) bool {
				return items[i].KmStart > items[j].KmStart
			})
		}

		retroReflectivity := responses.RoadRetroReflectivityDetailData{
			Color:                                  retroReflectivityType,
			RoadRetroReflectivityDetailHeaderItems: items,
		}
		resp.Datas = append(resp.Datas, retroReflectivity)

		sort.Slice(resp.Datas, func(i, j int) bool {
			return resp.Datas[i].Color < resp.Datas[j].Color
		})
	}

	resp.ID = retroReflectivitys.ID
	resp.IDParent = retroReflectivitys.IDParent
	resp.Status = retroReflectivitys.Status
	resp.UpdatedDate = helpers.SetTimeToString(retroReflectivitys.UpdatedDate)
	resp.Direction = road.RoadInfo.Direction

	userInfo, err := u.repo.GetUserDepartmentById(retroReflectivitys.UpdatedBy)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logs.Error(err)
			return responses.RoadConditionDetails{}, err
		}
	}

	resp.UpdatedBy.UID = fmt.Sprintf("%d", userInfo.Id)
	resp.UpdatedBy.UserName = userInfo.Username
	resp.UpdatedBy.FullName = userInfo.Firstname + " " + userInfo.Lastname
	if userInfo.ProfileImgPath != "" {
		resp.UpdatedBy.ProfilePicture = os.Getenv("STORAGE_IP") + "/" + userInfo.ProfileImgPath
	}

	if len(resp.Datas) == 0 {
		return responses.NoData{}, nil
	}

	return resp, nil
}

func (u *useCase) GetDashboardCondition(roadIDs []string, depotCodes []string, filter requests.Condition) (interface{}, error) {

	var resp responses.ConditionDashboard
	var roadConditions []models.RoadConditionDashboard

	var roadRetroReflectivitys []models.RoadRetroReflectivityDashboard
	var err error

	if filter.ConditionType == 5 {

		roadRetroReflectivitys, err = u.repo.GetRoadRetroReflectivityDashboard(roadIDs, depotCodes, filter)
		if err != nil {
			logs.Error(err)
			return nil, responses.NewAppErr(500, err.Error())
		}

	} else {

		roadConditions, err = u.repo.GetRoadConditionDashboard(roadIDs, depotCodes, filter)
		if err != nil {
			logs.Error(err)
			return nil, responses.NewAppErr(500, err.Error())
		}

	}

	conditionType := ""

	switch filter.ConditionType {
	case 1:
		conditionType = "iri"
	case 2:
		conditionType = "mpd"
	case 3:
		conditionType = "rut"
	case 4:
		conditionType = "ifi"
	case 5:

		roadLineGrades, err := u.repo.GetRoadLineGradesByID(filter.ConditionOwnerID)
		if err != nil {
			logs.Error(err)
			return resp, err
		}

		resp, err = u.RefactiveGradeAnalysis(roadRetroReflectivitys, roadLineGrades)
		if err != nil {
			logs.Error(err)
			return resp, err
		}

	default:
		logs.Error("ไม่พบประเภทสภาพทางที่เลือก")
		return resp, responses.NewAppErr(400, "ไม่พบประเภทสภาพทางที่เลือก")
	}

	if conditionType != "" {
		conditionGrades, err := u.repo.GetRoadConditionGradesByID(filter.ConditionOwnerID, conditionType)
		if err != nil {
			logs.Error(err)
			return resp, err
		}

		resp, err = u.ConditionGradeAnalysis(roadConditions, conditionType, conditionGrades)
		if err != nil {
			logs.Error(err)
			return resp, err
		}
	}

	if len(roadIDs) > 1 || len(depotCodes) > 1 {
		resp.HasMutipleRoad = true
	}

	return resp, nil
}

type GradeValue struct {
	Grade models.RefGrade
	Value float64
}

func (u *useCase) RefactiveGradeAnalysis(roadRetroReflectivitys []models.RoadRetroReflectivityDashboard, refactiveGrades []models.ParamsRoadLinePreload) (responses.ConditionDashboard, error) {
	var resp responses.ConditionDashboard

	resp.Chart = responses.ConditionChart{}

	summaryByConditionGrade := make(map[models.RefGrade]float64)

	summaryConditionGradeByLane := make(map[int]map[models.RefGrade]float64)

	totalConditionValueByLane := make(map[int]float64)

	for _, roadRetroReflectivity := range roadRetroReflectivitys {

		for _, roadRetroReflectivityRange := range roadRetroReflectivity.RoadRetroReflectivityRanges {

			//100 m = 2  | RoadRetroReflectivityRanges
			//25 m = 1  | roadRetroReflectivityM

			if len(refactiveGrades) == 0 {
				return responses.ConditionDashboard{}, responses.NewAppErr(400, "ไม่พบเกณฑ์การจำแนกสภาพทางที่เลือก")
			}
			if refactiveGrades[0].RefOwnerRoadLine.RefReflectivityRangeID == 2 {

				conditionValue := roadRetroReflectivityRange.RetroAvg

				if conditionValue == nil {
					continue
				}

				totalRange := math.Abs(roadRetroReflectivityRange.KmStart - roadRetroReflectivityRange.KmEnd)
				totalConditionValueByLane[roadRetroReflectivity.LineNo] += *conditionValue * totalRange

				processRoadRetroReflectivityGrades(roadRetroReflectivity.LineNo, totalRange, refactiveGrades, conditionValue, roadRetroReflectivityRange.RefStripeColorID, summaryByConditionGrade, summaryConditionGradeByLane)

			}

			for _, roadRetroReflectivityM := range roadRetroReflectivityRange.RoadRetroReflectivityMs {
				if refactiveGrades[0].RefOwnerRoadLine.RefReflectivityRangeID == 1 {
					conditionValue := roadRetroReflectivityM.RetroAvg

					totalRange := math.Abs(roadRetroReflectivityM.KmStart - roadRetroReflectivityM.KmEnd)

					totalConditionValueByLane[roadRetroReflectivity.LineNo] += *conditionValue * totalRange

					processRoadRetroReflectivityGrades(roadRetroReflectivity.LineNo, totalRange, refactiveGrades, conditionValue, roadRetroReflectivityM.RefStripeColorID, summaryByConditionGrade, summaryConditionGradeByLane)

				}

			}

		}

	}

	km := 1000.00

	var grades []GradeValue
	for grade, value := range summaryByConditionGrade {
		grades = append(grades, GradeValue{Grade: grade, Value: value})
	}
	// Sort slice by Grade.ID
	sort.Slice(grades, func(i, j int) bool {
		return grades[i].Grade.ID < grades[j].Grade.ID
	})

	var labels []string
	var data []float64
	var colors []string

	resp.Chart.Lable = []string{}
	resp.Chart.Data = []float64{}
	resp.Chart.Color = []string{}

	for _, gv := range grades {
		labels = append(labels, gv.Grade.Name)
		data = append(data, gv.Value/km)
		colors = append(colors, gv.Grade.Color)
	}

	resp.Chart.Name = ""
	resp.Chart.Lable = labels
	resp.Chart.Data = data
	resp.Chart.Color = colors

	for laneNo, gradeMap := range summaryConditionGradeByLane {

		var totalValue float64
		var details []responses.ConditionDetailKm

		for grade, value := range gradeMap {

			totalValue += value
			details = append(details, responses.ConditionDetailKm{
				RefGradeID:   grade.ID,
				RefGradeName: grade.Name,
				Value:        value / km,
			})
		}

		totalKm := totalValue / km

		avgValue := 0.0
		if len(details) > 0 {

			avgValue = totalConditionValueByLane[laneNo] / totalValue
		}

		for i := range details {
			if totalValue != 0 {
				details[i].ValuePercent = (details[i].Value / totalKm) * 100
			} else {
				details[i].ValuePercent = 0
			}
		}

		sort.Slice(details, func(i, j int) bool {
			return details[i].RefGradeID < details[j].RefGradeID
		})

		resp.Table = append(resp.Table, responses.ConditionTable{
			LaneNo:   laneNo,
			TotalKm:  totalKm,
			AvgValue: avgValue,
			DetailKm: details,
		})

	}
	sort.Slice(resp.Table, func(i, j int) bool {
		return resp.Table[i].LaneNo < resp.Table[j].LaneNo
	})

	return resp, nil
}

func (u *useCase) ConditionGradeAnalysis(roadConditions []models.RoadConditionDashboard, conditionType string, conditionGrades []models.ParamsConditionPreload) (responses.ConditionDashboard, error) {

	//surveyType := make(map[string][]models.RoadConditionSurveyM)

	var resp responses.ConditionDashboard
	resp.Chart = responses.ConditionChart{}

	summaryByConditionGrade := make(map[models.RefGrade]float64)

	summaryConditionGradeByLane := make(map[int]map[models.RefGrade]float64)

	totalConditionValueByLane := make(map[int]float64)

	summaryConditionGradeByRoadID := make(map[int]map[models.RefGrade]float64)

	for _, roadCondition := range roadConditions {
		summaryByConditionGradeRoad := make(map[models.RefGrade]float64)
		if summaryConditionGradeByRoadID[roadCondition.RoadId] != nil {
			summaryByConditionGradeRoad = summaryConditionGradeByRoadID[roadCondition.RoadId]
		}

		for _, roadConditionSurvey := range roadCondition.RoadConditionSurveys {

			//1 km = 3 | roadConditionSurvey
			//100 m = 2  | roadConditionSurvey100m
			//25 m = 1  | roadConditionSurveyM

			if len(conditionGrades) == 0 {
				return responses.ConditionDashboard{}, responses.NewAppErr(400, "ไม่พบเกณฑ์การจำแนกสภาพทางที่เลือก")
			}
			if conditionGrades[0].RefOwner.RefConditionRangeID == 3 {

				conditionValue := getConditionValueDashboard(conditionType, roadConditionSurvey)
				if conditionValue == nil {
					continue
				}

				totalRange := math.Abs(roadConditionSurvey.KmStart - roadConditionSurvey.KmEnd)
				totalConditionValueByLane[roadCondition.LaneNo] += *conditionValue * totalRange

				processConditionGrades(roadCondition.LaneNo, totalRange, conditionGrades, conditionValue, roadConditionSurvey.SurveyType, summaryByConditionGrade, summaryConditionGradeByLane, summaryByConditionGradeRoad)

			}

			for _, roadConditionSurvey100M := range roadConditionSurvey.RoadConditionSurvey100Ms {
				if conditionGrades[0].RefOwner.RefConditionRangeID == 2 {
					conditionValue := getConditionValueDashboard(conditionType, roadConditionSurvey100M)
					if conditionValue == nil {
						continue
					}
					totalRange := math.Abs(roadConditionSurvey100M.KmStart - roadConditionSurvey100M.KmEnd)

					totalConditionValueByLane[roadCondition.LaneNo] += *conditionValue * totalRange

					processConditionGrades(roadCondition.LaneNo, totalRange, conditionGrades, conditionValue, roadConditionSurvey100M.SurveyType, summaryByConditionGrade, summaryConditionGradeByLane, summaryByConditionGradeRoad)

				}

				for _, roadConditionSurveyM := range roadConditionSurvey100M.RoadConditionSurveyMs {
					if conditionGrades[0].RefOwner.RefConditionRangeID == 1 {

						conditionValue := getConditionValueDashboard(conditionType, roadConditionSurveyM)
						if conditionValue == nil {
							continue
						}

						totalRange := math.Abs(roadConditionSurveyM.KmStart - roadConditionSurveyM.KmEnd)

						totalConditionValueByLane[roadCondition.LaneNo] += *conditionValue * totalRange

						processConditionGrades(roadCondition.LaneNo, totalRange, conditionGrades, conditionValue, roadConditionSurveyM.SurveyType, summaryByConditionGrade, summaryConditionGradeByLane, summaryByConditionGradeRoad)

					}

				}
			}

		}

		summaryConditionGradeByRoadID[roadCondition.RoadId] = summaryByConditionGradeRoad
	}

	gradeConditionRoadIDs := make(map[int][]string)

	for roadID, grades := range summaryConditionGradeByRoadID {
		for grade := range grades {
			gradeConditionRoadIDs[grade.ID] = append(gradeConditionRoadIDs[grade.ID], fmt.Sprintf(`%v`, roadID))
		}
	}

	km := 1000.00

	var grades []GradeValue
	for grade, value := range summaryByConditionGrade {
		grades = append(grades, GradeValue{Grade: grade, Value: value})
	}
	// Sort slice by Grade.ID
	sort.Slice(grades, func(i, j int) bool {
		return grades[i].Grade.ID < grades[j].Grade.ID
	})

	sumAllGrade := map[string]float64{}

	for laneNo, gradeMap := range summaryConditionGradeByLane {

		var totalValue float64
		var details []responses.ConditionDetailKm

		for grade, value := range gradeMap {

			totalValue += value
			details = append(details, responses.ConditionDetailKm{
				RefGradeID:   grade.ID,
				RefGradeName: grade.Name,
				Value:        value / km,
			})

			sumAllGrade[grade.Name] += value / km

		}

		totalKm := totalValue / km

		avgValue := 0.0
		if len(details) > 0 {

			avgValue = totalConditionValueByLane[laneNo] / totalValue
		}

		for i := range details {
			if totalValue != 0 {
				details[i].ValuePercent = (details[i].Value / totalKm) * 100
			} else {
				details[i].ValuePercent = 0
			}
		}

		sort.Slice(details, func(i, j int) bool {
			return details[i].RefGradeID < details[j].RefGradeID
		})

		resp.Table = append(resp.Table, responses.ConditionTable{
			LaneNo:   laneNo,
			TotalKm:  totalKm,
			AvgValue: avgValue,
			DetailKm: details,
		})

	}

	var labels []string
	var data []float64
	var colors []string
	var roadID []string
	DIGIT, _ := strconv.Atoi(os.Getenv("DIGIT"))
	for _, gv := range grades {
		labels = append(labels, gv.Grade.Name)
		data = append(data, helpers.RoundFloat(sumAllGrade[gv.Grade.Name], DIGIT))
		colors = append(colors, gv.Grade.Color)
		roadID = append(roadID, strings.Join(gradeConditionRoadIDs[gv.Grade.ID], ","))
	}

	resp.Chart.Name = ""

	resp.Chart.Lable = labels
	if resp.Chart.Lable == nil {
		resp.Chart.Lable = []string{}

	}
	resp.Chart.Data = data
	if resp.Chart.Data == nil {
		resp.Chart.Data = []float64{}
	}
	resp.Chart.Color = colors
	if resp.Chart.Color == nil {
		resp.Chart.Color = []string{}
	}

	resp.Chart.RoadID = roadID
	if resp.Chart.RoadID == nil {
		resp.Chart.RoadID = []string{}
	}

	if resp.Table == nil {
		resp.Table = []responses.ConditionTable{}
	}

	sort.Slice(resp.Table, func(i, j int) bool {
		return resp.Table[i].LaneNo < resp.Table[j].LaneNo
	})

	return resp, nil
}

func (u *useCase) ConditionGradeAnalysisMap(roadConditions []models.RoadConditionDashboard, conditionType string, conditionGrades []models.ParamsConditionPreload) ([]responses.DashboardConditionMap, error) {

	//surveyType := make(map[string][]models.RoadConditionSurveyM)

	var resps []responses.DashboardConditionMap

	for _, roadCondition := range roadConditions {
		var resp responses.DashboardConditionMap

		for _, roadConditionSurvey := range roadCondition.RoadConditionSurveys {

			//1 km = 3 | roadConditionSurvey
			//100 m = 2  | roadConditionSurvey100m
			//25 m = 1  | roadConditionSurveyM

			if len(conditionGrades) == 0 {
				return []responses.DashboardConditionMap{}, responses.NewAppErr(400, "ไม่พบเกณฑ์การจำแนกสภาพทางที่เลือก")
			}
			if conditionGrades[0].RefOwner.RefConditionRangeID == 3 {

				conditionValue := getConditionValueDashboard(conditionType, roadConditionSurvey)
				if conditionValue == nil {
					continue
				}
				grade := processConditionGradesMapColor(conditionGrades, conditionValue, roadConditionSurvey.SurveyType)
				resp.Color = grade.Color

				surveyGeomJson, err := helpers.ConvertThegeomToGeomJSON(roadConditionSurvey.SurveyTheGeom)
				if err != nil {
					logs.Error(err)
					return resps, err
				}

				resp.TheGeom = surveyGeomJson
				resps = append(resps, resp)
			}

			for _, roadConditionSurvey100M := range roadConditionSurvey.RoadConditionSurvey100Ms {
				if conditionGrades[0].RefOwner.RefConditionRangeID == 2 {
					conditionValue := getConditionValueDashboard(conditionType, roadConditionSurvey100M)
					if conditionValue == nil {
						continue
					}

					grade := processConditionGradesMapColor(conditionGrades, conditionValue, roadConditionSurvey100M.SurveyType)
					resp.Color = grade.Color

					surveyGeomJson, err := helpers.ConvertThegeomToGeomJSON(roadConditionSurvey.SurveyTheGeom)
					if err != nil {
						logs.Error(err)
						return resps, err
					}

					resp.TheGeom = surveyGeomJson

					resps = append(resps, resp)

				}

				for _, roadConditionSurveyM := range roadConditionSurvey100M.RoadConditionSurveyMs {
					if conditionGrades[0].RefOwner.RefConditionRangeID == 1 {

						conditionValue := getConditionValueDashboard(conditionType, roadConditionSurveyM)
						if conditionValue == nil {
							continue
						}

						grade := processConditionGradesMapColor(conditionGrades, conditionValue, roadConditionSurveyM.SurveyType)
						resp.Color = grade.Color

						surveyGeomJson, err := helpers.ConvertThegeomToGeomJSON(roadConditionSurvey.SurveyTheGeom)
						if err != nil {
							logs.Error(err)
							return resps, err
						}

						resp.TheGeom = surveyGeomJson

						resps = append(resps, resp)

					}

				}
			}

		}

	}

	return resps, nil
}

func processConditionGradesMapColor(conditionGrades []models.ParamsConditionPreload, conditionValue *float64, surveyType string) models.RefGrade {

	var resp models.RefGrade
	for _, conditionGrade := range conditionGrades {

		if conditionValue != nil && checkCondition(*conditionValue, surveyType, conditionGrade) {
			return conditionGrade.RefGrade
		}
	}
	return resp
}

func processRoadRetroReflectivityGradesMapColor(retroReflectivityGrades []models.ParamsRoadLinePreload, conditionValue *float64, stripeColor int) models.RefGrade {
	var resp models.RefGrade
	for _, retroReflectivityGrade := range retroReflectivityGrades {

		if conditionValue != nil && checkRetroReflectivity(*conditionValue, stripeColor, retroReflectivityGrade) {
			return retroReflectivityGrade.RefGrade

		}
	}
	return resp
}

func processRoadRetroReflectivityGrades(laneNo int, totalKm float64, retroReflectivityGrades []models.ParamsRoadLinePreload, conditionValue *float64, stripeColor int, summaryByConditionGrade map[models.RefGrade]float64, summaryConditionGradeByLane map[int]map[models.RefGrade]float64) {

	// Ensure the lane map exists
	if _, exists := summaryConditionGradeByLane[laneNo]; !exists {
		summaryConditionGradeByLane[laneNo] = make(map[models.RefGrade]float64)
	}

	for _, retroReflectivityGrade := range retroReflectivityGrades {

		// Check for nil conditionValue to avoid panic
		if conditionValue != nil && checkRetroReflectivity(*conditionValue, stripeColor, retroReflectivityGrade) {
			// Add totalKm to the accumulated kilometers for each condition met
			summaryByConditionGrade[retroReflectivityGrade.RefGrade] += totalKm
			summaryConditionGradeByLane[laneNo][retroReflectivityGrade.RefGrade] += totalKm

		} else {
			// Ensure that every grade is initialized in the map
			if _, exists := summaryByConditionGrade[retroReflectivityGrade.RefGrade]; !exists {
				summaryByConditionGrade[retroReflectivityGrade.RefGrade] = 0
			}
			if _, exists := summaryConditionGradeByLane[laneNo][retroReflectivityGrade.RefGrade]; !exists {
				summaryConditionGradeByLane[laneNo][retroReflectivityGrade.RefGrade] = 0
			}
		}
	}
}

func processConditionGrades(laneNo int, totalKm float64, conditionGrades []models.ParamsConditionPreload, conditionValue *float64, surveyType string, summaryByConditionGrade map[models.RefGrade]float64, summaryConditionGradeByLane map[int]map[models.RefGrade]float64, summaryByConditionGradeRoad map[models.RefGrade]float64) {

	// Ensure the lane map exists
	if _, exists := summaryConditionGradeByLane[laneNo]; !exists {
		summaryConditionGradeByLane[laneNo] = make(map[models.RefGrade]float64)
	}

	for _, conditionGrade := range conditionGrades {

		// Check for nil conditionValue to avoid panic
		if conditionValue != nil && checkCondition(*conditionValue, surveyType, conditionGrade) {
			// Add totalKm to the accumulated kilometers for each condition met
			summaryByConditionGrade[conditionGrade.RefGrade] += totalKm
			summaryConditionGradeByLane[laneNo][conditionGrade.RefGrade] += totalKm
			summaryByConditionGradeRoad[conditionGrade.RefGrade] += totalKm
		} else {
			// Ensure that every grade is initialized in the map
			if _, exists := summaryByConditionGrade[conditionGrade.RefGrade]; !exists {
				summaryByConditionGrade[conditionGrade.RefGrade] = 0
			}
			if _, exists := summaryConditionGradeByLane[laneNo][conditionGrade.RefGrade]; !exists {
				summaryConditionGradeByLane[laneNo][conditionGrade.RefGrade] = 0
			}

			if _, exists := summaryByConditionGradeRoad[conditionGrade.RefGrade]; !exists {
				summaryByConditionGrade[conditionGrade.RefGrade] = 0
			}
		}
	}
}

func checkRetroReflectivity(input float64, stripeColor int, retroReflectivity models.ParamsRoadLinePreload) bool {
	var leftCondition, rightCondition string
	var leftValue, rightValue float64

	switch stripeColor {
	case 1: //White
		leftCondition, rightCondition = retroReflectivity.LeftConditionWhite, retroReflectivity.RightConditionWhite
		leftValue, rightValue = retroReflectivity.LeftValueWhite, retroReflectivity.RightValueWhite
	case 2: //Yellow
		leftCondition, rightCondition = retroReflectivity.LeftConditionYellow, retroReflectivity.RightConditionYellow
		leftValue, rightValue = retroReflectivity.LeftValueYellow, retroReflectivity.RightValueYellow
	default:
		fmt.Printf("Unsupported stripe color: %d\n", stripeColor)
		return false
	}

	// fmt.Printf("%s_LEFT: %f %s %f\n", surfaceType, input, leftCondition, leftValue)
	// fmt.Printf("%s_RIGHT: %f %s %f\n", surfaceType, input, rightCondition, rightValue)

	leftConditionMet := compareConditionLeft(input, leftCondition, leftValue)
	rightConditionMet := compareConditionRight(input, rightCondition, rightValue)

	return leftConditionMet && rightConditionMet
}

func checkCondition(input float64, surfaceType string, condition models.ParamsConditionPreload) bool {
	var leftCondition, rightCondition string
	var leftValue, rightValue float64

	switch surfaceType {
	case "AC":
		leftCondition, rightCondition = condition.LeftConditionAC, condition.RightConditionAC
		leftValue, rightValue = condition.LeftValueAC, condition.RightValueAC
	case "CC":
		leftCondition, rightCondition = condition.LeftConditionCC, condition.RightConditionCC
		leftValue, rightValue = condition.LeftValueCC, condition.RightValueCC
	default:
		fmt.Printf("Unsupported surface type: %s\n", surfaceType)
		return false
	}

	// fmt.Printf("%s_LEFT: %f %s %f\n", surfaceType, input, leftCondition, leftValue)
	// fmt.Printf("%s_RIGHT: %f %s %f\n", surfaceType, input, rightCondition, rightValue)

	leftConditionMet := compareConditionLeft(input, leftCondition, leftValue)
	rightConditionMet := compareConditionRight(input, rightCondition, rightValue)

	return leftConditionMet && rightConditionMet
}

func getConditionValueDashboard(conditionType string, survey interface{}) *float64 {
	switch v := survey.(type) {
	case models.RoadConditionSurveyDashboard:
		switch conditionType {
		case "iri":
			return v.IRI
		case "mpd":
			return v.MPD
		case "rut":
			return v.RUT
		case "ifi":
			return v.IFI
		}
	case models.RoadConditionSurvey100MDashboard:
		switch conditionType {
		case "iri":
			return v.IRI
		case "mpd":
			return v.MPD
		case "rut":
			return v.RUT
		case "ifi":
			return v.IFI
		}
	case models.RoadConditionSurveyMDashboard:
		switch conditionType {
		case "iri":
			return v.IRI
		case "mpd":
			return v.MPD
		case "rut":
			return v.RUT
		case "ifi":
			return v.IFI
		}
	}
	return nil
}

func getConditionValue(conditionType string, survey interface{}) *float64 {
	switch v := survey.(type) {
	case models.RoadConditionSurveyPreload:
		switch conditionType {
		case "iri":
			return v.IRI
		case "mpd":
			return v.MPD
		case "rut":
			return v.RUT
		case "ifi":
			return v.IFI
		}
	case models.RoadConditionSurvey100MPreload:
		switch conditionType {
		case "iri":
			return v.IRI
		case "mpd":
			return v.MPD
		case "rut":
			return v.RUT
		case "ifi":
			return v.IFI
		}
	case models.RoadConditionSurveyM:
		switch conditionType {
		case "iri":
			return v.IRI
		case "mpd":
			return v.MPD
		case "rut":
			return v.RUT
		case "ifi":
			return v.IFI
		}
	}
	return nil
}

func compareConditionLeft(input float64, operator string, value float64) bool {
	switch operator {
	case "<=":
		return value <= input
	case "<":
		return value < input
	case ">=":
		return value >= input
	case ">":
		return value > input
	case "==":
		return value == input
	case "!=":
		return value != input
	default:
		fmt.Printf("Unsupported operator: %s\n", operator)
		return false
	}
}

func compareConditionRight(input float64, operator string, value float64) bool {
	switch operator {
	case "<=":
		return input <= value
	case "<":
		return input < value
	case ">=":
		return input >= value
	case ">":
		return input > value
	case "==":
		return input == value
	case "!=":
		return input != value
	default:
		fmt.Printf("Unsupported operator: %s\n", operator)
		return false
	}
}

func (u *useCase) GetDashboardConditionMap(roadIDs []string, depotCodes []string, pageParam, limitParam string, filter requests.ConditionMap) (interface{}, error) {

	var pagination responses.Pagination
	var resps []responses.DashboardConditionMap
	var roadConditions []models.RoadConditionDashboard

	var roadRetroReflectivitys []models.RoadRetroReflectivityDashboard
	var err error

	if filter.ConditionType == 5 {

		roadRetroReflectivitys, err = u.repo.GetRoadRetroReflectivityDashboardMap(roadIDs, depotCodes, filter)
		if err != nil {
			logs.Error(err)
			return nil, responses.NewAppErr(500, err.Error())
		}

	} else {

		roadConditions, err = u.repo.GetRoadConditionDashboardMap(roadIDs, depotCodes, filter)
		if err != nil {
			logs.Error(err)
			return nil, responses.NewAppErr(500, err.Error())
		}

	}

	conditionType := ""

	switch filter.ConditionType {
	case 1:
		conditionType = "iri"
	case 2:
		conditionType = "mpd"
	case 3:
		conditionType = "rut"
	case 4:
		conditionType = "ifi"
	case 5:

		roadLineGrades, err := u.repo.GetRoadLineGradesByID(filter.ConditionOwnerID)
		if err != nil {
			logs.Error(err)
			return pagination, err
		}

		resps, err = u.RefactiveGradeAnalysisMap(roadRetroReflectivitys, roadLineGrades)
		if err != nil {
			logs.Error(err)
			return pagination, err
		}

	default:
		logs.Error("ไม่พบประเภทสภาพทางที่เลือก")
		return pagination, responses.NewAppErr(400, "ไม่พบประเภทสภาพทางที่เลือก")
	}

	if conditionType != "" {

		conditionGrades, err := u.repo.GetRoadConditionGradesByID(filter.ConditionOwnerID, conditionType)
		if err != nil {
			logs.Error(err)
			return pagination, responses.NewAppErr(500, err.Error())
		}

		resps, err = u.ConditionGradeAnalysisMap(roadConditions, conditionType, conditionGrades)
		if err != nil {
			logs.Error(err)
			return pagination, responses.NewAppErr(500, err.Error())
		}
	}

	colorGroupedGeoJsons := make(map[string][]IGeoJson)
	for _, item := range resps {

		coordinates, err := convertInterfaceToFloat64Slice(item.TheGeom.Coordinates)
		if err != nil {
			log.Println(err)
			return pagination, responses.NewAppErr(500, err.Error())
		}

		geoJson := IGeoJson{Coordinates: coordinates}
		colorGroupedGeoJsons[item.Color] = append(colorGroupedGeoJsons[item.Color], geoJson)
	}

	// Merge lines for each color group
	var result []responses.DashboardConditionMap
	for color, geoJsons := range colorGroupedGeoJsons {
		mergedLines := mergeDashboardLineStrings(geoJsons, color)
		result = append(result, mergedLines...)
	}

	total := int64(len(result))

	limit, offset, page := helpers.GetlimitOffsetPage(limitParam, pageParam, total)

	if total == 0 {
		result = []responses.DashboardConditionMap{}
	} else if limit+offset > total {
		result = result[offset:total]
	} else {
		result = result[offset : limit+offset]
	}

	pagination = helpers.Pagination(result, int64(limit), int64(page), total)

	if result == nil {
		pagination.Items = []responses.DashboardConditionMap{}
	}

	return pagination, nil

}

func convertInterfaceToFloat64Slice(input interface{}) ([][]float64, error) {
	rawSlice, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("input is not a []interface{}")
	}

	if len(rawSlice) == 0 {
		return nil, fmt.Errorf("input slice is empty")
	}

	// Determine the type of the first element
	firstElem := rawSlice[0]
	switch firstElem.(type) {
	case []interface{}:
		// Input is [][]interface{}
		result := make([][]float64, len(rawSlice))
		for i, inner := range rawSlice {
			innerSlice, ok := inner.([]interface{})
			if !ok {
				fmt.Println("input", input)
				return nil, fmt.Errorf("item at index %d is not a []interface{}", i)
			}

			floatSlice := make([]float64, len(innerSlice))
			for j, value := range innerSlice {
				floatValue, ok := value.(float64)
				if !ok {
					return nil, fmt.Errorf("value at [%d][%d] is not of type float64", i, j)
				}
				floatSlice[j] = floatValue
			}
			result[i] = floatSlice
		}
		return result, nil
	case float64:
		// Input is []interface{} of float64 values
		floatSlice := make([]float64, len(rawSlice))
		for i, value := range rawSlice {
			floatValue, ok := value.(float64)
			if !ok {
				return nil, fmt.Errorf("value at index %d is not of type float64", i)
			}
			floatSlice[i] = floatValue
		}
		// Wrap the flat slice in an outer slice to return [][]float64
		return [][]float64{floatSlice}, nil
	default:
		return nil, fmt.Errorf("unsupported element type %T", firstElem)
	}
}

func mergeDashboardLineStrings(geojsonData []IGeoJson, color string) []responses.DashboardConditionMap {
	mergedLines := make([]responses.DashboardConditionMap, 0)
	processedIndices := make(map[int]bool)
	processedCoordinates := make(map[string]bool) // Track processed coordinates

	for index, line := range geojsonData {
		if processedIndices[index] {
			continue
		}

		currentLine := line.Coordinates

		// Convert coordinates to string for duplicate checking
		currentLineStr := fmt.Sprintf("%v", currentLine)
		if processedCoordinates[currentLineStr] {
			continue // Skip processing if coordinates are already processed
		}

		// Traverse the remaining lines to find the ones that connect to the current line
		for i := index + 1; i < len(geojsonData); i++ {
			nextLine := geojsonData[i]
			lastPoint := currentLine[len(currentLine)-1]
			if len(nextLine.Coordinates) == 0 {
				continue
			}
			firstPoint := nextLine.Coordinates[0]

			// Check if the last point of the current line is equal to the first point of the next line
			if lastPoint[0] == firstPoint[0] && lastPoint[1] == firstPoint[1] {
				currentLine = append(currentLine, nextLine.Coordinates[1:]...)
				processedIndices[i] = true
			}
		}

		mergedLines = append(mergedLines, responses.DashboardConditionMap{
			Color: color,
			TheGeom: responses.GeomJSON{
				Type:        "LineString",
				Coordinates: currentLine,
			},
		})

		processedIndices[index] = true
		processedCoordinates[currentLineStr] = true // Mark coordinates as processed
	}
	return mergedLines
}

func (u *useCase) RefactiveGradeAnalysisMap(roadRetroReflectivitys []models.RoadRetroReflectivityDashboard, refactiveGrades []models.ParamsRoadLinePreload) ([]responses.DashboardConditionMap, error) {
	var resps []responses.DashboardConditionMap

	for _, roadRetroReflectivity := range roadRetroReflectivitys {
		var resp responses.DashboardConditionMap

		for _, roadRetroReflectivityRange := range roadRetroReflectivity.RoadRetroReflectivityRanges {

			//100 m = 2  | RoadRetroReflectivityRanges
			//25 m = 1  | roadRetroReflectivityM

			if len(refactiveGrades) == 0 {
				return resps, responses.NewAppErr(400, "ไม่พบเกณฑ์การจำแนกสภาพทางที่เลือก")
			}
			if refactiveGrades[0].RefOwnerRoadLine.RefReflectivityRangeID == 2 {

				conditionValue := roadRetroReflectivityRange.RetroAvg

				if conditionValue == nil {
					continue
				}

				grade := processRoadRetroReflectivityGradesMapColor(refactiveGrades, conditionValue, roadRetroReflectivityRange.RefStripeColorID)
				resp.Color = grade.Color

				surveyGeomJson, err := helpers.ConvertThegeomToGeomJSON(roadRetroReflectivityRange.RetroRangeTheGeom)
				if err != nil {
					logs.Error(err)
					return resps, err
				}

				resp.TheGeom = surveyGeomJson
				resps = append(resps, resp)

			}

			for _, roadRetroReflectivityM := range roadRetroReflectivityRange.RoadRetroReflectivityMs {
				if refactiveGrades[0].RefOwnerRoadLine.RefReflectivityRangeID == 1 {
					conditionValue := roadRetroReflectivityM.RetroAvg

					grade := processRoadRetroReflectivityGradesMapColor(refactiveGrades, conditionValue, roadRetroReflectivityM.RefStripeColorID)
					resp.Color = grade.Color

					surveyGeomJson, err := helpers.ConvertThegeomToGeomJSON(roadRetroReflectivityM.RetroMTheGeom)
					if err != nil {
						logs.Error(err)
						return resps, err
					}

					resp.TheGeom = surveyGeomJson
					resps = append(resps, resp)

				}

			}

		}

	}

	return resps, nil
}

// func (u *useCase) GetDashboardConditionMap(roadIDs []int, depotCodes []string, page, limit, offset int, filter requests.Condition) (interface{}, error) {

// 	conditionType := ""

// 	switch filter.ConditionType {
// 	case 1:
// 		conditionType = "iri"
// 	case 2:
// 		conditionType = "mpd"
// 	case 3:
// 		conditionType = "rut"
// 	case 4:
// 		conditionType = "ifi"
// 	case 5:

// 	default:
// 		logs.Error("ไม่พบประเภทสภาพทางที่เลือก")
// 		return nil, responses.NewAppErr(400, "ไม่พบประเภทสภาพทางที่เลือก")

// 	}

// 	var err error
// 	respondConditionsSum := responses.ConditionGeomSumRespond{}
// 	geomList, total, err := u.geomCondition(roadIDs, conditionType, filter.ConditionOwnerID, filter.Year, page, offset, limit)
// 	if err != nil {
// 		logs.Error(err)
// 		return "", responses.NewAppErr(400, err.Error())
// 	}
// 	respondConditionsSum.GeomList = geomList
// 	respondConditionsSum.TotalPage = int(math.Ceil(float64(total) / float64(limit)))
// 	// respondConditionsSum.Offset = offset
// 	respondConditionsSum.Page = page
// 	return respondConditionsSum, nil

// }

func (u *useCase) GetConditionTotalPage(roadIDs []int, conditionType string, year, ownerID, page, offset, limit int) (int, error) {
	var err error
	total, err := u.geomConditionTotalPage(roadIDs, conditionType, ownerID, year, page, offset, limit)
	if err != nil {
		logs.Error(err)
		return 0, responses.NewAppErr(400, err.Error())
	}

	return total, nil
}

func (u *useCase) geomConditionTotalPage(roadIDs []int, conditionType string, ownerID, year, page, offset, limit int) (int, error) {
	RoadConditions := []models.RoadCondition{}
	if len(roadIDs) != 0 {
		whereRoad := helpers.SliceIntToString(roadIDs)
		query := fmt.Sprintf("road_id IN (%s) AND year = %d AND status = 'A'", whereRoad, year)
		err := u.repo.GetDataList(&RoadConditions, query)
		if err != nil {
			logs.Error(err)
			return 0, responses.NewAppErr(400, err.Error())
		}

	} else {
		err := u.repo.GetDataList(&RoadConditions, fmt.Sprintf(" year = %d", year))
		if err != nil {
			logs.Error(err)
			return 0, responses.NewAppErr(400, err.Error())
		}
	}
	return len(RoadConditions), nil
}
func (u *useCase) geomCondition(roadIDs []int, conditionType string, ownerID, year, page, offset, limit int) ([]responses.GeomList, int, error) {
	dict := make(map[int][]responses.Geometry)
	geomCls := []responses.GeomList{}

	refOwner, err := u.repo.GetRoadConditionGradesByID(ownerID, conditionType)
	if err != nil {
		return geomCls, 0, err
	}

	RoadConditions := []models.RoadCondition{}
	if len(roadIDs) != 0 {
		whereRoad := helpers.SliceIntToString(roadIDs)
		query := fmt.Sprintf("road_id IN (%s) AND year = %d AND status = 'A'", whereRoad, year)
		err := u.repo.GetDataList(&RoadConditions, query)
		if err != nil {
			logs.Error(err)
			return geomCls, 0, responses.NewAppErr(400, err.Error())
		}
	} else {
		err = u.repo.GetDataList(&RoadConditions, fmt.Sprintf(" year = %d", year))
		if err != nil {
			logs.Error(err)
			return geomCls, 0, responses.NewAppErr(400, err.Error())
		}
	}
	if len(RoadConditions) == 0 {
		return geomCls, 0, nil
	}

	detailKm, err := u.initDetailKm(ownerID, conditionType)
	if err != nil {
		logs.Error(err)
		return geomCls, 0, err
	}
	detailKms := make(map[int]responses.DetailKm)
	for _, item := range detailKm {
		detailKms[item.RefGrade.ID] = item
	}

	var roadSurveys []models.RoadConditionSurvey
	err = u.repo.GetDataList(&roadSurveys, "")
	if err != nil {
		logs.Error(err)
		return geomCls, 0, responses.NewAppErr(400, err.Error())
	}

	roadSurvey100Ms, err := u.repo.GetRoadCondition100M("")
	if err != nil {
		logs.Error(err)
		return geomCls, 0, responses.NewAppErr(400, err.Error())
	}

	roadSurveyMs, err := u.repo.GetRoadCondition("")
	if err != nil {
		logs.Error(err)
		return geomCls, 0, responses.NewAppErr(400, err.Error())
	}

	// RoadByID
	var RoadByID []models.Road
	err = u.repo.GetDataList(&RoadByID, "")
	if err != nil {
		logs.Error(err)
		return geomCls, 0, responses.NewAppErr(400, err.Error())
	}

	var roadConditionData []models.RoadCondition
	totalItems := int64(len(RoadConditions))
	limitData := fmt.Sprintf("%d", limit)
	pageData := fmt.Sprintf("%d", page)

	limit1, offset1, _ := helpers.GetlimitOffsetPage(limitData, pageData, (totalItems))
	if totalItems == 0 {
		return []responses.GeomList{}, 0, nil
	} else if limit1+offset1 > totalItems {
		roadConditionData = RoadConditions[offset1:totalItems]
	} else {
		roadConditionData = RoadConditions[offset1 : limit1+offset1]
	}

	for _, roadCondition := range roadConditionData {

		var RoadSurveys []models.RoadConditionSurvey
		for _, i := range roadSurveys {
			if i.RoadConditionId == roadCondition.ID {
				RoadSurveys = append(RoadSurveys, i)
			}
		}

		var RoadSurvey100Ms []models.RoadConditionSurvey100M2
		for _, roadSurvey := range RoadSurveys {
			for _, i := range roadSurvey100Ms {
				if i.RoadConditionSurveyId == roadSurvey.ID {
					RoadSurvey100Ms = append(RoadSurvey100Ms, i)
				}
			}
		}

		var RoadSurveyMs []models.RoadConditionSurveyM2
		for _, roadSurvey100M := range RoadSurvey100Ms {
			for _, i := range roadSurveyMs {
				if i.RoadConditionSurvay100mID == roadSurvey100M.ID {
					RoadSurveyMs = append(RoadSurveyMs, i)
				}
			}
		}

		for index, k := range RoadSurveyMs {
			geoString, err := helpers.GeoJson(k.Geojson, index, RoadSurveyMs)
			if err != nil {
				return geomCls, 0, err
			}
			switch conditionType {
			case "ifi":
				gradeID := FindGradeID(*k.IFI, k.SurveyType, refOwner)
				dict[gradeID] = append(dict[gradeID], geoString)
			case "iri":
				gradeID := FindGradeID(*k.IRI, k.SurveyType, refOwner)
				dict[gradeID] = append(dict[gradeID], geoString)
			case "rut":
				gradeID := FindGradeID(*k.RUT, k.SurveyType, refOwner)
				dict[gradeID] = append(dict[gradeID], geoString)
			case "mpd":
				gradeID := FindGradeID(*k.MPD, k.SurveyType, refOwner)
				dict[gradeID] = append(dict[gradeID], geoString)
			}
		}

		// var roadByID models.Road
		// for _, i := range RoadByID {
		// 	if i.Id == roadCondition.RoadId {
		// 		roadByID = i
		// 	}
		// }

		for key, value := range dict {
			var geoJsons []IGeoJson
			for _, item := range value {
				var geoJson IGeoJson
				geoJson.Coordinates = item.Coordinates
				geoJsons = append(geoJsons, geoJson)
			}

			result := mergeLineStrings(geoJsons)

			var geomCl responses.GeomList
			geomCl.GradeID = detailKms[key].RefGrade.ID
			geomCl.Color = detailKms[key].RefGrade.Color
			geomCl.GeomCL = result
			geomCls = append(geomCls, geomCl)
		}
		// }
	}
	return geomCls, len(RoadConditions), nil
}

func (u *useCase) initDetailKm(ownerID int, conditionType string) ([]responses.DetailKm, error) {
	var detailKms []responses.DetailKm
	paramCon := []models.ParamsCondition{}
	err := u.repo.GetDataList(&paramCon, fmt.Sprintf("ref_owner_id = %d AND LOWER(condition_type) = '%s'", ownerID, conditionType))
	if err != nil {
		return detailKms, err
	}
	if len(paramCon) == 0 {
		logs.Info("ParamCondition not found check owner_id or condition_type")
		return detailKms, nil
	}
	refGradeIDs := []int{}
	for _, v := range paramCon {
		if !helpers.ContainsInt(v.RefGradeID, refGradeIDs) {
			refGradeIDs = append(refGradeIDs, v.RefGradeID)
		}
	}
	idString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(refGradeIDs)), ","), "[]")
	refGrades := []models.RefGrade{}
	err = u.repo.GetDataList(&refGrades, fmt.Sprintf("id IN (%s)", idString))
	if err != nil {
		return detailKms, err
	}

	for _, v := range refGrades {
		var detailKm responses.DetailKm
		detailKm.RefGrade.ID = v.ID
		detailKm.RefGrade.Name = v.Name
		detailKm.RefGrade.Color = v.Color
		detailKm.Value = 0.0
		detailKm.ValuePercent = 0.0
		detailKms = append(detailKms, detailKm)
	}
	sort.Slice(detailKms, func(i, j int) bool {
		return detailKms[i].RefGrade.ID < detailKms[j].RefGrade.ID
	})
	return detailKms, nil
}

type IGeoJson struct {
	Coordinates [][]float64
	Connections []int // Indices of the connected lines
}

func mergeLineStrings(geojsonData []IGeoJson) []map[string]interface{} {
	mergedLines := make([]map[string]interface{}, 0)
	processedIndices := make(map[int]bool)

	for index, line := range geojsonData {
		if processedIndices[index] {
			continue
		}

		currentLine := line.Coordinates

		// Traverse the remaining lines to find the ones that connect to the current line
		for i := index + 1; i < len(geojsonData); i++ {

			nextLine := geojsonData[i]
			lastPoint := currentLine[len(currentLine)-1]
			if len(nextLine.Coordinates) == 0 {
				continue
			}
			firstPoint := nextLine.Coordinates[0]

			// Check if the last point of the current line is equal to the first point of the next line
			if lastPoint[0] == firstPoint[0] && lastPoint[1] == firstPoint[1] {
				currentLine = append(currentLine, nextLine.Coordinates[1:]...)
				processedIndices[i] = true
			}
		}

		mergedLines = append(mergedLines, map[string]interface{}{
			"coordinates": currentLine,
		})

		processedIndices[index] = true
	}

	return mergedLines
}

func FindGradeID(value float64, conditionType string, refOwner []models.ParamsConditionPreload) int {
	for _, item := range refOwner {
		if checkCondition(value, conditionType, item) {
			return item.RefGradeID
		}
	}
	return 0

}

// func FindGradeID2(refOwner []responses.ConditionTypeInit, value float64) int {
// 	for index, item := range refOwner {
// 		if checkCondition(value, item) {
// 			return index
// 		}
// 	}
// 	return 0

// }

// func GetConditionInit(paramsCondition []models.ParamsConditionPreload) []responses.ConditionListNewInit {
// 	ConditionListNews := []responses.ConditionListNewInit{}
// 	if len(paramsCondition) == 0 {
// 		return ConditionListNews
// 	}
// 	condition := responses.Condition{}
// 	for _, v := range paramsCondition {

// 		if !v.RefOwner.IsActive {
// 			continue
// 		}

// 		switch v.ConditionType {
// 		case "IFI":
// 			ifi := GetConditionType(v)
// 			condition.IFI = append(condition.IFI, ifi)
// 		case "IRI":
// 			iri := GetConditionType(v)
// 			condition.IRI = append(condition.IRI, iri)
// 		case "MPD":
// 			mpd := GetConditionType(v)
// 			condition.MPD = append(condition.MPD, mpd)
// 		case "RUT":
// 			rut := GetConditionType(v)
// 			condition.RUT = append(condition.RUT, rut)
// 		default:
// 			log.Println("there is no condition type")
// 		}
// 	}
// 	if len(condition.IFI) != 0 {
// 		ConditionListNew := GenRespondInit(condition, "IFI")
// 		ConditionListNews = append(ConditionListNews, ConditionListNew)
// 	}
// 	if len(condition.IRI) != 0 {
// 		ConditionListNew := GenRespondInit(condition, "IRI")
// 		ConditionListNews = append(ConditionListNews, ConditionListNew)
// 	}
// 	if len(condition.MPD) != 0 {
// 		ConditionListNew := GenRespondInit(condition, "MPD")
// 		ConditionListNews = append(ConditionListNews, ConditionListNew)
// 	}
// 	if len(condition.RUT) != 0 {
// 		ConditionListNew := GenRespondInit(condition, "RUT")
// 		ConditionListNews = append(ConditionListNews, ConditionListNew)
// 	}
// 	return ConditionListNews
// }

// func checkCondition(input float64, condition responses.ConditionTypeInit) bool {
// 	leftConditionMet := false
// 	rightConditionMet := false

// 	switch condition.LeftCondition {
// 	case "<=":
// 		if input >= condition.LeftValue {
// 			leftConditionMet = true
// 		}
// 	case "<":
// 		if input > condition.LeftValue {
// 			leftConditionMet = true
// 		}
// 	case ">=":
// 		if input <= condition.LeftValue {
// 			leftConditionMet = true
// 		}
// 	case ">":
// 		if input < condition.LeftValue {
// 			leftConditionMet = true
// 		}
// 	case "==":
// 		if input == condition.LeftValue {
// 			leftConditionMet = true
// 		}
// 	case "!=":
// 		if input != condition.LeftValue {
// 			leftConditionMet = true
// 		}
// 	}

// 	switch condition.RightCondition {
// 	case "<=":
// 		if input <= condition.RightValue {
// 			rightConditionMet = true
// 		}
// 	case "<":
// 		if input < condition.RightValue {
// 			rightConditionMet = true
// 		}
// 	case ">=":
// 		if input >= condition.RightValue {
// 			rightConditionMet = true
// 		}
// 	case ">":
// 		if input > condition.RightValue {
// 			rightConditionMet = true
// 		}
// 	case "==":
// 		if input == condition.RightValue {
// 			rightConditionMet = true
// 		}
// 	case "!=":
// 		if input != condition.RightValue {
// 			rightConditionMet = true
// 		}
// 	}

// 	return leftConditionMet && rightConditionMet
// }

// func GetConditionType(condition models.ParamsConditionPreload) responses.ConditionType {
// 	return responses.ConditionType{
// 		Grade:            condition.RefGrade,
// 		LeftValueAC:      condition.LeftValueAC,
// 		LeftConditionAC:  condition.LeftConditionAC,
// 		RightValueAC:     condition.RightValueAC,
// 		RightConditionAC: condition.RightConditionAC,
// 		LeftValueCC:      condition.LeftValueCC,
// 		LeftConditionCC:  condition.LeftConditionCC,
// 		RightValueCC:     condition.RightValueCC,
// 		RightConditionCC: condition.RightConditionCC,
// 	}
// }

// func GenRespondInit(condition responses.Condition, conditionType string) responses.ConditionListNewInit {
// 	ConditionListNew := responses.ConditionListNewInit{ConditionType: conditionType}
// 	SurfaceTypeCondition := responses.SurfaceTypeConditionInit{}
// 	AC := []responses.ConditionTypeInit{}
// 	CC := []responses.ConditionTypeInit{}
// 	switch conditionType {
// 	case "IFI":
// 		for _, v := range condition.IFI {
// 			ac := responses.ConditionTypeInit{Grade: v.Grade,
// 				LeftValue:      v.LeftValueAC,
// 				LeftCondition:  v.LeftConditionAC,
// 				RightValue:     v.RightValueAC,
// 				RightCondition: v.RightConditionAC}
// 			cc := responses.ConditionTypeInit{Grade: v.Grade,
// 				LeftValue:      v.LeftValueCC,
// 				LeftCondition:  v.LeftConditionCC,
// 				RightValue:     v.RightValueCC,
// 				RightCondition: v.RightConditionCC}
// 			AC = append(AC, ac)
// 			CC = append(CC, cc)

// 		}
// 	case "IRI":
// 		for _, v := range condition.IRI {
// 			ac := responses.ConditionTypeInit{Grade: v.Grade,
// 				LeftValue:      v.LeftValueAC,
// 				LeftCondition:  v.LeftConditionAC,
// 				RightValue:     v.RightValueAC,
// 				RightCondition: v.RightConditionAC}
// 			cc := responses.ConditionTypeInit{Grade: v.Grade,
// 				LeftValue:      v.LeftValueCC,
// 				LeftCondition:  v.LeftConditionCC,
// 				RightValue:     v.RightValueCC,
// 				RightCondition: v.RightConditionCC}
// 			AC = append(AC, ac)
// 			CC = append(CC, cc)

// 		}
// 	case "MPD":
// 		for _, v := range condition.MPD {
// 			ac := responses.ConditionTypeInit{Grade: v.Grade,
// 				LeftValue:      v.LeftValueAC,
// 				LeftCondition:  v.LeftConditionAC,
// 				RightValue:     v.RightValueAC,
// 				RightCondition: v.RightConditionAC}
// 			cc := responses.ConditionTypeInit{Grade: v.Grade,
// 				LeftValue:      v.LeftValueCC,
// 				LeftCondition:  v.LeftConditionCC,
// 				RightValue:     v.RightValueCC,
// 				RightCondition: v.RightConditionCC}
// 			AC = append(AC, ac)
// 			CC = append(CC, cc)

// 		}
// 	case "RUT":
// 		for _, v := range condition.RUT {
// 			ac := responses.ConditionTypeInit{Grade: v.Grade,
// 				LeftValue:      v.LeftValueAC,
// 				LeftCondition:  v.LeftConditionAC,
// 				RightValue:     v.RightValueAC,
// 				RightCondition: v.RightConditionAC}
// 			cc := responses.ConditionTypeInit{Grade: v.Grade,
// 				LeftValue:      v.LeftValueCC,
// 				LeftCondition:  v.LeftConditionCC,
// 				RightValue:     v.RightValueCC,
// 				RightCondition: v.RightConditionCC}
// 			AC = append(AC, ac)
// 			CC = append(CC, cc)

// 		}
// 	}

// 	SurfaceTypeCondition.AC = AC
// 	SurfaceTypeCondition.CC = CC
// 	ConditionListNew.SurfaceType = SurfaceTypeCondition
// 	return ConditionListNew
// }
