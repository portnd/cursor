package usecases

import (
	"context"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/xuri/excelize/v2"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gorm.io/gorm"
)

// ////////////// NEW MIMS ////////////////

func (u *UseCase) Report5(year, SectionIDstr, typ string) (interface{}, error) {
	summaryData := []models.Summary{}
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	SectionID, err := strconv.Atoi(SectionIDstr)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	// var pathResult interface{}
	Surface, err := u.Repo.GetSurfaceForRoadCondition()
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	SurfaceReportResponds := []responses.SurfaceReportRespond{}
	roadIDs, err := u.Repo.GetRoadFromSectionIDForRoadCondition(SectionID)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	for _, roadID := range roadIDs {
		data, summary, err := u.Report5ByID(yearInt, roadID)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
		SurfaceReportResponds = append(SurfaceReportResponds, data)
		summaryData = append(summaryData, summary...)
	}

	tableData := []responses.SurfaceReportTable{}
	for _, v := range SurfaceReportResponds {
		tableData = append(tableData, v.Table...)
	}

	tableDataReport, tableSum := combineTable(tableData, Surface)
	dataForPieChart := combinePieChart(summaryData)
	var series []float64
	var colors []string
	var labels []string
	for _, item := range dataForPieChart {
		series = append(series, item.Value)
		colors = append(colors, item.Summary.ColorCode)
		labels = append(labels, item.Summary.Name)
	}
	var pieChart responses.PieChart
	pieChart.Colors = colors
	pieChart.Labels = labels
	pieChart.Series = series
	var data responses.SurfaceReportRespond
	data.PieChart = pieChart
	data.Table = tableDataReport
	data.TableSum = tableSum
	//////////
	roadSectionInfo, err := u.Repo.GetRoadSectionByIDForRoadCondition(SectionID)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	roadGroup, err := u.Repo.GetRoadGroupByIDForRoadCondition(roadSectionInfo.RoadGroupId)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	data.RoadGroupName = roadGroup.Number
	data.RoadSectionNumber = roadSectionInfo.Number
	data.RoadSectionName = roadSectionInfo.NameOriginTH + " - " + roadSectionInfo.NameDestinationTH
	data.KmStart = helpers.FormatKM(int64(roadSectionInfo.KmStart))
	data.KmEnd = helpers.FormatKM(int64(roadSectionInfo.KmEnd))
	data.StrRoadLength = fmt.Sprintf(`%.3f`, roadSectionInfo.Distance)
	data.Year = fmt.Sprint(yearInt + 543)
	if typ == "excel" {
		pathResult, err := ExportExcelType5(&data)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
		return pathResult, nil
	}
	pathResult, err := helpers.RequestExport(data, "TEMPLATE_GENARAL_TYPE5_AEK", typ)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}

	return pathResult, nil
}

func combineTable(data []responses.SurfaceReportTable, surface []models.Surface) ([]responses.SurfaceReportTable, responses.SurfaceLaneTypeReport) {
	sumOneLane := 0.0
	sumTwoLane := 0.0
	sumThreeLane := 0.0
	sumFourLane := 0.0
	sumMoreLane := 0.0
	sum := 0.0
	mapTable := make(map[string]responses.SurfaceReportTable)
	for _, v := range data {
		if table, ok := mapTable[v.SurfaceName]; ok {
			// var item
			newTable := table
			item := table.SurfaceLaneType
			oneLane, _ := strconv.ParseFloat(item.OneLane, 64)
			twoLane, _ := strconv.ParseFloat(item.TwoLane, 64)
			threeLane, _ := strconv.ParseFloat(item.ThreeLane, 64)
			fourLane, _ := strconv.ParseFloat(item.FourLane, 64)
			moreLane, _ := strconv.ParseFloat(item.MoreThanFour, 64)

			newData := v.SurfaceLaneType
			newOneLane, _ := strconv.ParseFloat(newData.OneLane, 64)
			newTwoLane, _ := strconv.ParseFloat(newData.TwoLane, 64)
			newThreeLane, _ := strconv.ParseFloat(newData.ThreeLane, 64)
			newFourLane, _ := strconv.ParseFloat(newData.FourLane, 64)
			newMoreLane, _ := strconv.ParseFloat(newData.MoreThanFour, 64)

			item.OneLane = fmt.Sprintf("%.2f", oneLane+newOneLane)
			item.TwoLane = fmt.Sprintf("%.2f", twoLane+newTwoLane)
			item.ThreeLane = fmt.Sprintf("%.2f", threeLane+newThreeLane)
			item.FourLane = fmt.Sprintf("%.2f", fourLane+newFourLane)
			item.MoreThanFour = fmt.Sprintf("%.2f", moreLane+newMoreLane)
			item.Sum = fmt.Sprintf("%.2f", oneLane+newOneLane+twoLane+newTwoLane+threeLane+newThreeLane+fourLane+newFourLane+moreLane+newMoreLane)
			newTable.SurfaceLaneType = item
			mapTable[v.SurfaceName] = newTable

			sumOneLane += newOneLane
			sumTwoLane += newTwoLane
			sumThreeLane += newThreeLane
			sumFourLane += newFourLane
			sumMoreLane += newMoreLane
			sum += newOneLane + newTwoLane + newThreeLane + newFourLane + newMoreLane

		} else {
			mapTable[v.SurfaceName] = v
			newData := v.SurfaceLaneType
			newOneLane, _ := strconv.ParseFloat(newData.OneLane, 64)
			newTwoLane, _ := strconv.ParseFloat(newData.TwoLane, 64)
			newThreeLane, _ := strconv.ParseFloat(newData.ThreeLane, 64)
			newFourLane, _ := strconv.ParseFloat(newData.FourLane, 64)
			newMoreLane, _ := strconv.ParseFloat(newData.MoreThanFour, 64)

			sumOneLane += newOneLane
			sumTwoLane += newTwoLane
			sumThreeLane += newThreeLane
			sumFourLane += newFourLane
			sumMoreLane += newMoreLane
			sum += newOneLane + newTwoLane + newThreeLane + newFourLane + newMoreLane

		}
	}

	result := []responses.SurfaceReportTable{}
	for _, v := range surface {
		if _, ok := mapTable[v.Name]; ok {
			result = append(result, mapTable[v.Name])
		}
	}
	var tablesSum responses.SurfaceLaneTypeReport
	tablesSum.OneLane = fmt.Sprintf("%.2f", sumOneLane)
	tablesSum.TwoLane = fmt.Sprintf("%.2f", sumTwoLane)
	tablesSum.ThreeLane = fmt.Sprintf("%.2f", sumThreeLane)
	tablesSum.FourLane = fmt.Sprintf("%.2f", sumFourLane)
	tablesSum.MoreThanFour = fmt.Sprintf("%.2f", sumMoreLane)
	tablesSum.Sum = fmt.Sprintf("%.2f", sum)
	return result, tablesSum
}

func combinePieChart(Summary []models.Summary) []models.Summary {
	mapSummary := make(map[string]models.Summary)
	result := []models.Summary{}
	for _, v := range Summary {
		if data, ok := mapSummary[v.Summary.Name]; ok {
			newData := data
			newData.Value += v.Value
			mapSummary[v.Summary.Name] = newData
		} else {
			mapSummary[v.Summary.Name] = v
		}
	}

	for _, v := range mapSummary {
		result = append(result, v)
	}
	return result
}

func (u *UseCase) Report5ByID(year, roadID int) (responses.SurfaceReportRespond, []models.Summary, error) {
	// yearInt, err := strconv.Atoi(year)
	// if err != nil {
	//  return nil, responses.NewAppErr(400, err.Error())
	// }

	// roadID, err := strconv.Atoi(roadIDStr)
	// if err != nil {
	//  return nil, responses.NewAppErr(400, err.Error())
	// }

	road, err := u.Repo.GetRoadInfoForRoadCondition(roadID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		logs.Error(err)
		return responses.SurfaceReportRespond{}, []models.Summary{}, err
	}
	road.RoadLengthStr = fmt.Sprint(math.Abs(float64(road.KmEnd)-float64(road.KmStart)) / 1000)

	var result responses.SurfaceRespond
	// surfaceInfo, err := s.summaryRepo.GetRoadSurfaceInfo(roadID)
	// if err != nil {
	//  return result, err
	// }
	surfaceInfo, err := u.Repo.GetDataMartInfoForRoadCondition(roadID)
	if err != nil {
		return responses.SurfaceReportRespond{}, []models.Summary{}, err
	}
	// return surfaceInfo, nil
	if len(surfaceInfo) == 0 {
		return responses.SurfaceReportRespond{}, []models.Summary{}, nil
	}
	summary, err := u.Repo.GetInitialSurfaceArrayForRoadCondition()
	if err != nil {
		return responses.SurfaceReportRespond{}, []models.Summary{}, err
	}
	summaryData := summaryPart(summary, surfaceInfo)
	result.Summary = append(result.Summary, summaryData.Summary...)
	detailAndGeom := combineData(surfaceInfo)
	surfaceTable := detailAndGeom.SurfaceDashboardTable
	sort.Slice(surfaceTable, func(i, j int) bool {
		return surfaceTable[i].ID < surfaceTable[j].ID
	})
	result.SurfaceDashboardTable = surfaceTable
	// result.GeomList = append(result.GeomList, detailAndGeom.GeomList...)

	sort.Slice(result.Detail.DetailKm, func(i, j int) bool {
		return result.Detail.DetailKm[i].LaneNo < result.Detail.DetailKm[j].LaneNo
	})

	// return result, nil
	var tables []responses.SurfaceReportTable
	oneLane := 0.0
	twoLane := 0.0
	threeLane := 0.0
	fourLane := 0.0
	moreThanFour := 0.0
	sum := 0.0
	for _, item := range result.SurfaceDashboardTable {
		var table responses.SurfaceReportTable
		table.SurfaceName = item.SurfaceName
		table.SurfaceLaneType.OneLane = fmt.Sprintf("%.2f", item.SurfaceLaneType.OneLane)
		table.SurfaceLaneType.TwoLane = fmt.Sprintf("%.2f", item.SurfaceLaneType.TwoLane)
		table.SurfaceLaneType.ThreeLane = fmt.Sprintf("%.2f", item.SurfaceLaneType.ThreeLane)
		table.SurfaceLaneType.FourLane = fmt.Sprintf("%.2f", item.SurfaceLaneType.FourLane)
		table.SurfaceLaneType.MoreThanFour = fmt.Sprintf("%.2f", item.SurfaceLaneType.MoreThanFour)
		table.SurfaceLaneType.Sum = fmt.Sprintf("%.2f", item.SurfaceLaneType.OneLane+item.SurfaceLaneType.TwoLane+item.SurfaceLaneType.ThreeLane+item.SurfaceLaneType.FourLane+item.SurfaceLaneType.MoreThanFour)
		tables = append(tables, table)
		oneLane += item.SurfaceLaneType.OneLane
		twoLane += item.SurfaceLaneType.TwoLane
		threeLane += item.SurfaceLaneType.ThreeLane
		fourLane += item.SurfaceLaneType.FourLane
		moreThanFour += item.SurfaceLaneType.MoreThanFour
		sum += item.SurfaceLaneType.OneLane + item.SurfaceLaneType.TwoLane + item.SurfaceLaneType.ThreeLane + item.SurfaceLaneType.FourLane + item.SurfaceLaneType.MoreThanFour
	}
	var series []float64
	var colors []string
	var labels []string
	for _, item := range result.Summary {
		series = append(series, item.Value)
		colors = append(colors, item.Summary.ColorCode)
		labels = append(labels, item.Summary.Name)
	}
	var pieChart responses.PieChart
	pieChart.Colors = colors
	pieChart.Labels = labels
	pieChart.Series = series

	var tablesSum responses.SurfaceLaneTypeReport
	tablesSum.OneLane = fmt.Sprintf("%.2f", oneLane)
	tablesSum.TwoLane = fmt.Sprintf("%.2f", twoLane)
	tablesSum.ThreeLane = fmt.Sprintf("%.2f", threeLane)
	tablesSum.FourLane = fmt.Sprintf("%.2f", fourLane)
	tablesSum.MoreThanFour = fmt.Sprintf("%.2f", moreThanFour)
	tablesSum.Sum = fmt.Sprintf("%.2f", sum)
	var data responses.SurfaceReportRespond
	data.PieChart = pieChart
	data.Table = tables
	data.TableSum = tablesSum
	//////////
	data.RoadGroupName = road.RoadGroupName
	// data.RoadName = road.RoadName
	// data.RoadCode = road.RoadCode
	data.KmStart = helpers.FormatKM(int64(road.KmStart))
	data.KmEnd = helpers.FormatKM(int64(road.KmEnd))
	data.StrRoadLength = road.RoadLengthStr
	data.Year = fmt.Sprint(year + 543)

	// return data, nil
	// if typ == "excel" {
	// pathResult, err = helpers.ExportExcelType8(&data)
	// if err != nil {
	//  return nil, responses.NewAppErr(400, err.Error())
	// }
	// } else {
	// pathResult, err = helpers.RequestExport(data, "TEMPLATE_GENARAL_TYPE5_AEK", typ)
	// if err != nil {
	//  return responses.SurfaceReportRespond{}, responses.NewAppErr(400, err.Error())
	// }
	// }

	return data, result.Summary, nil
}

func summaryPart(surfaces []models.Surface, surfaceInfo []models.SurfaceInfo) responses.SurfaceRespond {
	var arrData responses.SurfaceRespond
	dict := make(map[int]float64)
	for _, v := range surfaceInfo {
		rangeM := math.Abs(v.KmEnd - v.KmStart)
		rangeKm := rangeM / 1000
		dict[int(v.SurfaceID)] += rangeKm
	}
	sort.Slice(surfaces, func(i, j int) bool {
		return surfaces[i].ID < surfaces[j].ID
	})
	for _, s := range surfaces {
		_, isval := dict[s.ID]
		if isval {
			value := dict[s.ID]
			arrData.Summary = append(arrData.Summary, models.Summary{Summary: s, Value: value})
		}

	}
	return arrData
}

func combineData(surfaceInfo []models.SurfaceInfo) responses.SurfaceRespond {
	var result responses.SurfaceRespond
	laneTypeGroup := findLaneType(surfaceInfo)
	result.Summary = make([]models.Summary, len(surfaceInfo))

	// mapGeomID := make(map[string]models.GeomList)
	for key, rsID := range laneTypeGroup {
		mapDetailIdxs := make(map[string]models.SubDetailKm)
		for _, s := range surfaceInfo {

			if !(helpers.ContainsInt(s.RoadSurfaceID, rsID)) {
				continue
			}
			rangeM := math.Abs(s.KmEnd - s.KmStart)
			rangeKm := rangeM / 1000
			// surface_idx := s.SurfaceID - 1
			// result.Summary[surface_idx].Value += float64(rangeKm)
			// if !(helpers.ContainsInt(s.LaneNo, result.Detail.LaneCountList)) {
			//  result.Detail.LaneCountList = append(result.Detail.LaneCountList, s.LaneNo)
			// }
			// detailID := strconv.Itoa(s.LaneNo)
			idStr := strconv.Itoa(int(s.SurfaceID))
			// idStr := detailID + "_" + strSurIdx
			if _, ok := mapDetailIdxs[idStr]; !ok {
				var SubDetailKm models.SubDetailKm
				SubDetailKm.Surface.ID = int(s.SurfaceID)
				SubDetailKm.Surface.Name = s.SurfaceName
				SubDetailKm.LaneNo = s.LaneNo
				SubDetailKm.Value = 0
				SubDetailKm.LaneType = key
				mapDetailIdxs[idStr] = SubDetailKm
			}
			newRange := mapDetailIdxs[idStr]
			newRange.Value += float64(rangeKm)
			mapDetailIdxs[idStr] = newRange
			// geomID := strconv.Itoa(int(s.RoadID)) + "_" + strconv.Itoa(int(s.SurfaceID))

			// if _, ok := mapGeomID[geomID]; !ok {
			//  var addGeom models.GeomList
			//  addGeom.ID = int(s.RoadID)
			//  addGeom.Surface.ID = int(s.SurfaceID)
			//  addGeom.Surface.Name = s.SurfaceName
			//  addGeom.Surface.ColorCode = s.ColorCode
			//  addGeom.Code = s.RoadCode
			//  mapGeomID[geomID] = addGeom

			// }
			// newGeom := mapGeomID[geomID]
			// newGeom.GeomCl = append(newGeom.GeomCl, s.Geometry)
			// mapGeomID[geomID] = newGeom

		}
		for _, detail := range mapDetailIdxs {
			result.Detail.DetailKm = append(result.Detail.DetailKm, detail)
		}
	}
	laneTypeInfo := make(map[string]responses.SurfaceLaneType)
	surfaceKey := make(map[string]int)
	for _, v := range result.Detail.DetailKm {
		helpers.PrintlnJson(v.LaneType)
		switch v.LaneType {
		case "1เลน":
			surfaceKey[v.Surface.Name] = v.Surface.ID
			oldValue := laneTypeInfo[v.Surface.Name]
			oldValue.OneLane += v.Value
			laneTypeInfo[v.Surface.Name] = oldValue
		case "2เลน":
			surfaceKey[v.Surface.Name] = v.Surface.ID
			oldValue := laneTypeInfo[v.Surface.Name]
			oldValue.TwoLane += v.Value
			laneTypeInfo[v.Surface.Name] = oldValue
		case "3เลน":
			surfaceKey[v.Surface.Name] = v.Surface.ID
			oldValue := laneTypeInfo[v.Surface.Name]
			oldValue.ThreeLane += v.Value
			laneTypeInfo[v.Surface.Name] = oldValue
		case "4เลน":
			surfaceKey[v.Surface.Name] = v.Surface.ID
			oldValue := laneTypeInfo[v.Surface.Name]
			oldValue.FourLane += v.Value
			laneTypeInfo[v.Surface.Name] = oldValue
		default:
			surfaceKey[v.Surface.Name] = v.Surface.ID
			oldValue := laneTypeInfo[v.Surface.Name]
			oldValue.MoreThanFour += v.Value
			laneTypeInfo[v.Surface.Name] = oldValue
		}
	}
	for key, value := range laneTypeInfo {
		var surfaceLane responses.SurfaceDashboardTable

		_, isval := surfaceKey[key]
		if isval {
			surfaceLane.ID = surfaceKey[key]
		}

		surfaceLane.SurfaceName = key
		surfaceLane.SurfaceLaneType = value
		result.SurfaceDashboardTable = append(result.SurfaceDashboardTable, surfaceLane)
	}
	// for _, geom := range mapGeomID {
	//  result.GeomList = append(result.GeomList, geom)
	// }
	return result
}

func findLaneType(surfaceInfo []models.SurfaceInfo) map[string][]int {
	result := make(map[string][]int)
	result["1เลน"] = []int{}
	result["2เลน"] = []int{}
	result["3เลน"] = []int{}
	result["4เลน"] = []int{}
	result["มากกว่า4เลน"] = []int{}

	count := make(map[int]int)
	for _, v := range surfaceInfo {
		count[v.RoadSurfaceID] = v.LaneCount
	}
	// helpers.PrintlnJson(count)
	for key, value := range count {
		switch value {
		case 1:
			result["1เลน"] = append(result["1เลน"], key)
		case 2:
			result["2เลน"] = append(result["2เลน"], key)
		case 3:
			result["3เลน"] = append(result["3เลน"], key)
		case 4:
			result["4เลน"] = append(result["4เลน"], key)
		default:
			result["มากกว่า4เลน"] = append(result["มากกว่า4เลน"], key)
		}
	}
	return result
}

func ExportExcelType7(data responses.Report7, factor string) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	template := "TEMPLATE_GENARAL_TYPE7_AEK_EXCEL"
	path, err := ExportPNGChartForRoadCondition(data, "TEMPLATE_GENARAL_TYPE7_AEK_MAP", true)
	if err != nil {
		return nil, err
	}

	// img size checker
	file, err := os.Open(fmt.Sprint(path))
	if err != nil {
		fmt.Println("can not open file", err)
	}
	defer file.Close()

	imageConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Println("can not check image size", err)
	}

	defer os.Remove(fmt.Sprint(path))

	f, err := excelize.OpenFile(os.Getenv(template))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	//start

	sheetName := ""
	switch len(data.Grade) {
	case 2:
		sheetName = "Sheet1"
		f.DeleteSheet("Sheet2")
		f.DeleteSheet("Sheet3")
	case 3:
		sheetName = "Sheet2"
		f.DeleteSheet("Sheet1")
		f.DeleteSheet("Sheet3")
	case 4:
		sheetName = "Sheet3"
		f.DeleteSheet("Sheet2")
		f.DeleteSheet("Sheet1")
	}

	columns := []string{"B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N"}
	totalWidth := 0.0

	// Iterate over the columns and sum their widths
	for _, col := range columns {
		width, err := f.GetColWidth(sheetName, col)
		if err != nil {
			log.Fatal(err)
		}
		totalWidth += width
	}

	f.SetCellValue(sheetName, "C2", data.Header)

	f.SetCellValue(sheetName, "C3", fmt.Sprintf("หมายเลขทางหลวง : %s ตอนควบคุม : %s ", data.RoadGroupName, data.RoadSectionNumber))
	f.SetCellValue(sheetName, "C4", "ชื่อสายทาง : "+data.RoadSectionName)
	f.SetCellValue(sheetName, "C5", "กม.เริ่มต้น "+data.KmStart+" กม.สิ้นสุด "+data.KmEnd+" ระยะทาง "+data.RoadLength+" กม.")

	f.SetCellValue(sheetName, "C14", fmt.Sprintf("ข้อมูลสภาพทาง %s ปี %d", data.Type, data.Year))
	f.SetCellValue(sheetName, "C30", fmt.Sprintf("ข้อมูลสภาพทาง %s ปี %d", data.Type, data.Year))

	f.SetCellValue(sheetName, "F45", fmt.Sprintf("%s เฉลี่ย", data.Type))
	f.SetCellValue(sheetName, "I45", fmt.Sprintf("ระยะทางในแต่ละช่วง %s (กม.)", data.Type))

	startCell := 47
	for i, v := range data.Grade {
		column := ""
		columnChart := ""
		switch i {
		case 0:
			f.SetCellValue(sheetName, "C"+fmt.Sprintf("%d", 8), v.LeftValueAC)
			f.SetCellValue(sheetName, "E"+fmt.Sprintf("%d", 8), v.RefGrade.Name)
			f.SetCellValue(sheetName, "G"+fmt.Sprintf("%d", 8), v.RightValueAC)
			f.SetCellValue(sheetName, "J"+fmt.Sprintf("%d", 8), v.LeftValueCC)
			f.SetCellValue(sheetName, "L"+fmt.Sprintf("%d", 8), v.RefGrade.Name)
			f.SetCellValue(sheetName, "N"+fmt.Sprintf("%d", 8), v.RightValueCC)

			column = "I"
			columnChart = "I"
		case 1:
			f.SetCellValue(sheetName, "C"+fmt.Sprintf("%d", 9), v.LeftValueAC)
			f.SetCellValue(sheetName, "E"+fmt.Sprintf("%d", 9), v.RefGrade.Name)
			f.SetCellValue(sheetName, "G"+fmt.Sprintf("%d", 9), v.RightValueAC)
			f.SetCellValue(sheetName, "J"+fmt.Sprintf("%d", 9), v.LeftValueCC)
			f.SetCellValue(sheetName, "L"+fmt.Sprintf("%d", 9), v.RefGrade.Name)
			f.SetCellValue(sheetName, "N"+fmt.Sprintf("%d", 9), v.RightValueCC)

			column = "K"
			columnChart = "J"
		case 2:
			f.SetCellValue(sheetName, "C"+fmt.Sprintf("%d", 10), v.LeftValueAC)
			f.SetCellValue(sheetName, "E"+fmt.Sprintf("%d", 10), v.RefGrade.Name)
			f.SetCellValue(sheetName, "G"+fmt.Sprintf("%d", 10), v.RightValueAC)
			f.SetCellValue(sheetName, "J"+fmt.Sprintf("%d", 10), v.LeftValueCC)
			f.SetCellValue(sheetName, "L"+fmt.Sprintf("%d", 10), v.RefGrade.Name)
			f.SetCellValue(sheetName, "N"+fmt.Sprintf("%d", 10), v.RightValueCC)
			column = "M"
			columnChart = "K"
		case 3:
			f.SetCellValue(sheetName, "C"+fmt.Sprintf("%d", 11), v.LeftValueAC)
			f.SetCellValue(sheetName, "E"+fmt.Sprintf("%d", 11), v.RefGrade.Name)
			f.SetCellValue(sheetName, "G"+fmt.Sprintf("%d", 11), v.RightValueAC)
			f.SetCellValue(sheetName, "J"+fmt.Sprintf("%d", 11), v.LeftValueCC)
			f.SetCellValue(sheetName, "L"+fmt.Sprintf("%d", 11), v.RefGrade.Name)
			f.SetCellValue(sheetName, "N"+fmt.Sprintf("%d", 11), v.RightValueCC)
			column = "O"
			columnChart = "L"

		}
		f.SetCellValue(sheetName, column+fmt.Sprintf("%d", 46), v.RefGrade.Name)
		f.SetCellValue(sheetName, columnChart+fmt.Sprintf("%d", 17), v.RefGrade.Name)

	}
	if len(data.Table) > 0 {

		for _, item := range data.Table {
			f.MergeCell(sheetName, "B"+fmt.Sprintf("%d", startCell), "C"+fmt.Sprintf("%d", startCell))
			f.MergeCell(sheetName, "D"+fmt.Sprintf("%d", startCell), "E"+fmt.Sprintf("%d", startCell))
			f.MergeCell(sheetName, "F"+fmt.Sprintf("%d", startCell), "H"+fmt.Sprintf("%d", startCell))

			f.SetCellValue(sheetName, "B"+fmt.Sprintf("%d", startCell), item.LaneNo)

			f.SetCellValue(sheetName, "D"+fmt.Sprintf("%d", startCell), fmt.Sprintf("%.3f", item.TotalKm))

			f.SetCellValue(sheetName, "F"+fmt.Sprintf("%d", startCell), fmt.Sprintf("%.3f", item.AvgValue))
			for i, v := range item.DetailKm {
				column := ""
				switch i {
				case 0:
					column = "I"
					f.MergeCell(sheetName, "I"+fmt.Sprintf("%d", startCell), "J"+fmt.Sprintf("%d", startCell))
				case 1:
					column = "K"
					f.MergeCell(sheetName, "K"+fmt.Sprintf("%d", startCell), "L"+fmt.Sprintf("%d", startCell))
				case 2:
					column = "M"
					f.MergeCell(sheetName, "M"+fmt.Sprintf("%d", startCell), "N"+fmt.Sprintf("%d", startCell))

				case 3:
					column = "O"
					f.MergeCell(sheetName, "O"+fmt.Sprintf("%d", startCell), "P"+fmt.Sprintf("%d", startCell))

				}
				f.SetCellValue(sheetName, column+fmt.Sprintf("%d", startCell), v.Value)

			}
			startCell++
		}
		f.MergeCell(sheetName, "B"+fmt.Sprintf("%d", startCell), "C"+fmt.Sprintf("%d", startCell))
		f.MergeCell(sheetName, "D"+fmt.Sprintf("%d", startCell), "E"+fmt.Sprintf("%d", startCell))
		f.MergeCell(sheetName, "F"+fmt.Sprintf("%d", startCell), "H"+fmt.Sprintf("%d", startCell))

		f.SetCellValue(sheetName, "B"+fmt.Sprintf("%d", startCell), "รวม")
		totalKm, _ := strconv.ParseFloat(data.Summary[0], 64)
		avg, _ := strconv.ParseFloat(data.Summary[1], 64)

		f.SetCellValue(sheetName, "D"+fmt.Sprintf("%d", startCell), fmt.Sprintf("%.3f", totalKm))
		f.SetCellValue(sheetName, "F"+fmt.Sprintf("%d", startCell), fmt.Sprintf("%.3f", avg))
		switch len(data.Grade) {
		case 2:
			f.MergeCell(sheetName, "I"+fmt.Sprintf("%d", startCell), "J"+fmt.Sprintf("%d", startCell))
			f.MergeCell(sheetName, "K"+fmt.Sprintf("%d", startCell), "L"+fmt.Sprintf("%d", startCell))

			one, _ := strconv.ParseFloat(data.Summary[2], 64)
			two, _ := strconv.ParseFloat(data.Summary[3], 64)

			//Table
			f.SetCellValue(sheetName, "I"+fmt.Sprintf("%d", startCell), one)
			f.SetCellValue(sheetName, "K"+fmt.Sprintf("%d", startCell), two)

			//Graph
			f.SetCellValue(sheetName, "I"+fmt.Sprintf("%d", 18), fmt.Sprintf("%.3f", one))
			f.SetCellValue(sheetName, "J"+fmt.Sprintf("%d", 18), fmt.Sprintf("%.3f", two))

		case 3:
			f.MergeCell(sheetName, "I"+fmt.Sprintf("%d", startCell), "J"+fmt.Sprintf("%d", startCell))
			f.MergeCell(sheetName, "K"+fmt.Sprintf("%d", startCell), "L"+fmt.Sprintf("%d", startCell))
			f.MergeCell(sheetName, "M"+fmt.Sprintf("%d", startCell), "N"+fmt.Sprintf("%d", startCell))

			one, _ := strconv.ParseFloat(data.Summary[2], 64)
			two, _ := strconv.ParseFloat(data.Summary[3], 64)
			three, _ := strconv.ParseFloat(data.Summary[4], 64)

			//Table
			f.SetCellValue(sheetName, "I"+fmt.Sprintf("%d", startCell), fmt.Sprintf("%.3f", one))
			f.SetCellValue(sheetName, "K"+fmt.Sprintf("%d", startCell), fmt.Sprintf("%.3f", two))
			f.SetCellValue(sheetName, "M"+fmt.Sprintf("%d", startCell), fmt.Sprintf("%.3f", three))

			//Graph
			f.SetCellValue(sheetName, "I"+fmt.Sprintf("%d", 18), one)
			f.SetCellValue(sheetName, "J"+fmt.Sprintf("%d", 18), two)
			f.SetCellValue(sheetName, "K"+fmt.Sprintf("%d", 18), three)

		case 4:
			f.MergeCell(sheetName, "I"+fmt.Sprintf("%d", startCell), "J"+fmt.Sprintf("%d", startCell))
			f.MergeCell(sheetName, "K"+fmt.Sprintf("%d", startCell), "L"+fmt.Sprintf("%d", startCell))
			f.MergeCell(sheetName, "M"+fmt.Sprintf("%d", startCell), "N"+fmt.Sprintf("%d", startCell))
			f.MergeCell(sheetName, "O"+fmt.Sprintf("%d", startCell), "P"+fmt.Sprintf("%d", startCell))

			one, _ := strconv.ParseFloat(data.Summary[2], 64)
			two, _ := strconv.ParseFloat(data.Summary[3], 64)
			three, _ := strconv.ParseFloat(data.Summary[4], 64)
			four, _ := strconv.ParseFloat(data.Summary[5], 64)

			//Table
			f.SetCellValue(sheetName, "I"+fmt.Sprintf("%d", startCell), fmt.Sprintf("%.3f", one))
			f.SetCellValue(sheetName, "K"+fmt.Sprintf("%d", startCell), fmt.Sprintf("%.3f", two))
			f.SetCellValue(sheetName, "M"+fmt.Sprintf("%d", startCell), fmt.Sprintf("%.3f", three))
			f.SetCellValue(sheetName, "O"+fmt.Sprintf("%d", startCell), fmt.Sprintf("%.3f", four))

			//Graph
			f.SetCellValue(sheetName, "I"+fmt.Sprintf("%d", 18), one)
			f.SetCellValue(sheetName, "J"+fmt.Sprintf("%d", 18), two)
			f.SetCellValue(sheetName, "K"+fmt.Sprintf("%d", 18), three)
			f.SetCellValue(sheetName, "L"+fmt.Sprintf("%d", 18), four)

		}
		f.AddPicture(sheetName, "B"+fmt.Sprintf("%d", startCell+2), fmt.Sprint(path), &excelize.GraphicOptions{ScaleX: 700 / float64(imageConfig.Width), ScaleY: 600 / float64(imageConfig.Height)})

	}
	wrapTextStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{{
			Type:  "left",
			Color: "000000",
			Style: 1,
		}, {
			Type:  "top",
			Color: "000000",
			Style: 1,
		}, {
			Type:  "right",
			Color: "000000",
			Style: 1,
		}, {
			Type:  "bottom",
			Color: "000000",
			Style: 1,
		}},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		}, Font: &excelize.Font{
			Family: "TH SarabunPSK",
		},
	})
	switch len(data.Grade) {
	case 2:
		f.SetCellStyle(sheetName, "B"+fmt.Sprintf("%d", 47), "L"+fmt.Sprintf("%d", startCell), wrapTextStyle)

	case 3:
		f.SetCellStyle(sheetName, "B"+fmt.Sprintf("%d", 47), "N"+fmt.Sprintf("%d", startCell), wrapTextStyle)

	case 4:
		f.SetCellStyle(sheetName, "B"+fmt.Sprintf("%d", 47), "P"+fmt.Sprintf("%d", startCell), wrapTextStyle)

	}
	helpers.AddFooter(f, sheetName)

	reportName, err := helpers.ReportName("TEMPLATE_GENARAL_TYPE5_AEK_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}

func ExportExcelType5(data *responses.SurfaceReportRespond) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	template := "TEMPLATE_GENARAL_TYPE5_AEK_EXCEL"
	path, err := ExportPNGChartForRoadCondition(data, "TEMPLATE_GENARAL_TYPE5_AEK_PIE", false)
	if err != nil {
		return nil, err
	}

	// img size checker
	file, err := os.Open(fmt.Sprint(path))
	if err != nil {
		fmt.Println("can not open file", err)
	}
	defer file.Close()

	imageConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Println("can not check image size", err)
	}

	defer os.Remove(fmt.Sprint(path))

	f, err := excelize.OpenFile(os.Getenv(template))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	f.AddPicture("Sheet1", "E8", fmt.Sprint(path), &excelize.GraphicOptions{ScaleX: 320 / float64(imageConfig.Width), ScaleY: 330 / float64(imageConfig.Height)})
	//start

	f.SetCellValue("Sheet1", "C3", fmt.Sprintf("หมายเลขทางหลวง : %s ตอนควบคุม : %s ", data.RoadGroupName, data.RoadSectionNumber))
	f.SetCellValue("Sheet1", "C4", "ชื่อสายทาง : "+data.RoadSectionName)
	f.SetCellValue("Sheet1", "C5", "กม.เริ่มต้น "+data.KmStart+" กม.สิ้นสุด "+data.KmEnd+" ระยะทาง "+data.StrRoadLength+" กม.")

	startCell := 29
	for _, item := range data.Table {
		f.SetCellValue("Sheet1", "C"+fmt.Sprintf("%d", startCell), item.SurfaceName)
		f.SetCellValue("Sheet1", "D"+fmt.Sprintf("%d", startCell), item.SurfaceLaneType.MoreThanFour)
		f.SetCellValue("Sheet1", "E"+fmt.Sprintf("%d", startCell), item.SurfaceLaneType.FourLane)
		f.SetCellValue("Sheet1", "F"+fmt.Sprintf("%d", startCell), item.SurfaceLaneType.ThreeLane)
		f.SetCellValue("Sheet1", "G"+fmt.Sprintf("%d", startCell), item.SurfaceLaneType.TwoLane)
		f.SetCellValue("Sheet1", "H"+fmt.Sprintf("%d", startCell), item.SurfaceLaneType.OneLane)
		f.SetCellValue("Sheet1", "I"+fmt.Sprintf("%d", startCell), item.SurfaceLaneType.Sum)
		startCell++
	}
	f.SetCellValue("Sheet1", "C"+fmt.Sprintf("%d", startCell), "รวม")
	f.SetCellValue("Sheet1", "D"+fmt.Sprintf("%d", startCell), data.TableSum.MoreThanFour)
	f.SetCellValue("Sheet1", "E"+fmt.Sprintf("%d", startCell), data.TableSum.FourLane)
	f.SetCellValue("Sheet1", "F"+fmt.Sprintf("%d", startCell), data.TableSum.ThreeLane)
	f.SetCellValue("Sheet1", "G"+fmt.Sprintf("%d", startCell), data.TableSum.TwoLane)
	f.SetCellValue("Sheet1", "H"+fmt.Sprintf("%d", startCell), data.TableSum.OneLane)
	f.SetCellValue("Sheet1", "I"+fmt.Sprintf("%d", startCell), data.TableSum.Sum)
	wrapTextStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{{
			Type:  "left",
			Color: "000000",
			Style: 1,
		}, {
			Type:  "top",
			Color: "000000",
			Style: 1,
		}, {
			Type:  "right",
			Color: "000000",
			Style: 1,
		}, {
			Type:  "bottom",
			Color: "000000",
			Style: 1,
		}},
	})
	f.SetCellStyle("Sheet1", "B"+fmt.Sprintf("%d", 29), "I"+fmt.Sprintf("%d", startCell), wrapTextStyle)

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	helpers.AddFooter(f, "Sheet1")

	reportName, err := helpers.ReportName("TEMPLATE_GENARAL_TYPE5_AEK_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}

func ExportPNGChartForRoadCondition(data interface{}, HTMLTemplate string, isRp7 bool) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	templateName := os.Getenv(HTMLTemplate)
	result, err := helpers.InitDataToHtml(templateName, data, filePath)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	html, err := os.ReadFile(result)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	e := os.Remove(result)
	if e != nil {
		logs.Error(err)
		return nil, err
	}

	ctx, cancel := helpers.NewChromedpContext(context.Background(), log.Printf)
	defer cancel()

	var buf []byte
	if isRp7 {
		if err := chromedp.Run(ctx, PrintToPNGChartForRoadConditionRp7(string(html), &buf, false)); err != nil {
			logs.Error(err)
			return nil, err
		}
	} else {
		if err := chromedp.Run(ctx, PrintToPNGChartForRoadCondition(string(html), &buf, false)); err != nil {
			logs.Error(err)
			return nil, err
		}
	}

	unix := strconv.Itoa(int(time.Now().Unix()))

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	fullFilePath := filePath + unix + ".png"
	if err := ioutil.WriteFile(fullFilePath, buf, 0777); err != nil {
		logs.Error(err)
		return nil, err
	}

	return fullFilePath, nil
}

func PrintToPNGChartForRoadConditionRp7(html string, res *[]byte, isDelay bool) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {

			lctx, cancel := context.WithCancel(ctx)
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)

			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					cancel()
					wg.Done()
				}
			})
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			if err := page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx); err != nil {
				return err
			}
			delay := 5
			if isDelay {
				delay = 20
			}

			defer chromedp.Run(
				ctx,
				helpers.RunWithTimeOut(&ctx, time.Duration(delay), chromedp.Tasks{
					chromedp.WaitVisible("div#success"),
				}),
			)

			wg.Wait()
			return nil
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			width, height := 800, 600

			if err := emulation.SetDeviceMetricsOverride(int64(width), int64(height), 1.0, false).Do(ctx); err != nil {
				return err
			}

			buf, err := page.CaptureScreenshot().Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

func PrintToPNGChartForRoadCondition(html string, res *[]byte, isDelay bool) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {

			lctx, cancel := context.WithCancel(ctx)
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)

			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					cancel()
					wg.Done()
				}
			})
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			if err := page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx); err != nil {
				return err
			}
			delay := 5
			if isDelay {
				delay = 20
			}

			defer chromedp.Run(
				ctx,
				helpers.RunWithTimeOut(&ctx, time.Duration(delay), chromedp.Tasks{
					chromedp.WaitVisible("div#success"),
				}),
			)

			wg.Wait()
			return nil
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			width, height := 400, 450
			if err := emulation.SetDeviceMetricsOverride(int64(width), int64(height), 1.0, false).Do(ctx); err != nil {
				return err
			}

			buf, err := page.CaptureScreenshot().Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

type GradeValue struct {
	Grade models.RefGrade
	Value float64
}

func (u *UseCase) Report7(factor, year, SectionIDstr, measureID, typ string) (interface{}, error) {
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	SectionID, err := strconv.Atoi(SectionIDstr)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	OwnerID, err := strconv.Atoi(measureID)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	roadIDs, err := u.Repo.GetRoadFromSectionIDForRoadCondition(SectionID)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	filter := requests.Condition{Year: yearInt, ConditionOwnerID: OwnerID}
	data, grade, err := u.GetDashboardCondition(roadIDs, []string{}, filter, factor)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	mapData, err := u.GetDashboardConditionMap(roadIDs, filter, factor)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	// pathMap, err := helpers.RequestExport(mapData, "TEMPLATE_GENARAL_TYPE7_AEK_MAP", "html")
	// if err != nil {
	// 	return nil, responses.NewAppErr(400, err.Error())
	// }
	// fmt.Println(pathMap)
	mapData.AllLineString = u.Repo.GetCenterFromRoadSectionIDForRoadCondition(SectionID)
	report7 := responses.Report7{
		ConditionDashboardStr: data,
		MapData:               mapData,
	}
	summary := []string{}
	gradeID := []int{}
	var sumLen, SumAvg float64
	mapSum := make(map[int]float64)
	for _, v := range data.Table {
		sumLen += v.TotalKm
		SumAvg += v.AvgValue
		for _, k := range v.DetailKm {
			kVal, _ := strconv.ParseFloat(k.Value, 64)

			mapSum[k.RefGradeID] += kVal
			if !helpers.IsContainsSlice(k.RefGradeID, gradeID) {
				gradeID = append(gradeID, k.RefGradeID)
			}
		}
	}
	summary = append(summary, fmt.Sprintf("%.3f", sumLen), fmt.Sprintf("%.3f", SumAvg))
	helpers.SortSlice(gradeID, false)
	for _, v := range gradeID {
		data := mapSum[v]
		summary = append(summary, fmt.Sprintf("%.3f", data))
	}
	roadSectionInfo, err := u.Repo.GetRoadSectionByIDForRoadCondition(SectionID)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	roadGroup, err := u.Repo.GetRoadGroupByIDForRoadCondition(roadSectionInfo.RoadGroupId)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	switch strings.ToLower(factor) {
	case "iri":
		report7.Header = "รายงานสรุปข้อมูลค่าเฉลี่ยความขรุขระสากล (International Roughness Index : IRI)"
	case "ifi":
		report7.Header = "รายงานสรุปข้อมูลค่าเฉลี่ยดัชนีความเสียดทานสากล (International Friction Index : IFI)"
	case "mpd":
		report7.Header = "รายงานสรุปข้อมูลค่าเฉลี่ยความลึกของผิวทาง (Mean Profile Depth : MPD)"
	case "rut":
		report7.Header = "รายงานสรุปข้อมูลค่าเฉลี่ยความลึกร่องล้อ (Rutting : RUT)"
	}
	report7.Type = strings.ToUpper(factor)
	report7.RoadGroupName = roadGroup.Number
	report7.RoadSectionNumber = roadSectionInfo.Number
	report7.RoadSectionName = roadSectionInfo.NameOriginTH + " - " + roadSectionInfo.NameDestinationTH
	report7.KmStart = helpers.FormatKM(int64(roadSectionInfo.KmStart))
	report7.KmEnd = helpers.FormatKM(int64(roadSectionInfo.KmEnd))
	report7.RoadLength = helpers.FormatNumberFloat(float64(roadSectionInfo.Distance))
	report7.Year = yearInt + 543
	report7.Summary = summary
	report7.Grade = grade

	table := []responses.ConditionTableStr{}
	switch len(grade) {
	case 2:
		if len(summary) != 4 {

			report7.ConditionDashboardStr.Table = table
			// err := errors.New("grade and data value not match")
			// logs.Error(err)
			// return nil, responses.NewAppErr(400, err.Error())
		}
	case 3:
		if len(summary) != 5 {

			report7.ConditionDashboardStr.Table = table
			// err := errors.New("grade and data value not match")
			// logs.Error(err)
			// return nil, responses.NewAppErr(400, err.Error())
		}
	case 4:
		if len(summary) != 6 {

			report7.ConditionDashboardStr.Table = table
			// err := errors.New("grade and data value not match")
			// logs.Error(err)
			// return nil, responses.NewAppErr(400, err.Error())
		}
	default:

		report7.ConditionDashboardStr.Table = table
		// err := errors.New("number of grade data error")
		// logs.Error(err)
		// return nil, responses.NewAppErr(400, err.Error())
	}
	var pathResult interface{}
	if typ == "excel" {
		pathResult, err = ExportExcelType7(report7, factor)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	} else {
		pathResult, err = helpers.RequestExport(report7, "TEMPLATE_GENARAL_TYPE7_AEK", typ)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	}
	return pathResult, nil
}

func (u *UseCase) GetDashboardCondition(roadIDs []int, depotCodes []string, filter requests.Condition, factor string) (responses.ConditionDashboardStr, []models.ParamsConditionPreload, error) {

	var resp responses.ConditionDashboardStr
	var roadConditions []models.RoadConditionReport

	// var roadRetroReflectivitys []models.RoadRetroReflectivityDashboard
	var err error

	roadConditions, err = u.Repo.GetRoadConditionDashboardForRoadCondition(roadIDs, depotCodes, filter)
	if err != nil {
		logs.Error(err)
		return responses.ConditionDashboardStr{}, []models.ParamsConditionPreload{}, responses.NewAppErr(500, err.Error())
	}

	conditionType := strings.ToLower(factor)

	// roadLineGrades, err := u.Repo.GetRoadLineGradesByID(filter.ConditionOwnerID)
	// if err != nil {
	// 	logs.Error(err)
	// 	return resp, []models.ParamsConditionPreload{}, err
	// }

	// resp, err = u.RefactiveGradeAnalysis(roadRetroReflectivitys, roadLineGrades)
	// if err != nil {
	// 	logs.Error(err)
	// 	return resp, []models.ParamsConditionPreload{}, err
	// }
	conditionGrades := []models.ParamsConditionPreload{}
	if conditionType != "" {
		conditionGrades, err = u.Repo.GetRoadConditionGradesByIDForRoadCondition(filter.ConditionOwnerID, conditionType)
		if err != nil {
			logs.Error(err)
			return resp, []models.ParamsConditionPreload{}, err
		}

		resp, err = u.ConditionGradeAnalysis(roadConditions, conditionType, conditionGrades)
		if err != nil {
			logs.Error(err)
			return resp, []models.ParamsConditionPreload{}, err
		}
	}

	if len(roadIDs) > 1 || len(depotCodes) > 1 {
		resp.HasMutipleRoad = true
	}

	return resp, conditionGrades, nil
}

func (u *UseCase) RefactiveGradeAnalysis(roadRetroReflectivitys []models.RoadRetroReflectivityDashboard, refactiveGrades []models.ParamsRoadLinePreload) (responses.ConditionDashboard, error) {
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

func (u *UseCase) ConditionGradeAnalysis(roadConditions []models.RoadConditionReport, conditionType string, conditionGrades []models.ParamsConditionPreload) (responses.ConditionDashboardStr, error) {

	//surveyType := make(map[string][]models.RoadConditionSurveyM)

	var resp responses.ConditionDashboardStr
	resp.Chart = responses.ConditionChart{}

	summaryByConditionGrade := make(map[models.RefGrade]float64)

	summaryConditionGradeByLane := make(map[int]map[models.RefGrade]float64)

	totalConditionValueByLane := make(map[int]float64)

	for _, roadCondition := range roadConditions {

		for _, roadConditionSurvey := range roadCondition.RoadConditionSurveys {

			//1 km = 3 | roadConditionSurvey
			//100 m = 2  | roadConditionSurvey100m
			//25 m = 1  | roadConditionSurveyM

			if len(conditionGrades) == 0 {
				return responses.ConditionDashboardStr{}, responses.NewAppErr(400, "ไม่พบเกณฑ์การจำแนกสภาพทางที่เลือก")
			}
			if conditionGrades[0].RefOwner.RefConditionRangeID == 3 {

				conditionValue := getConditionValue(conditionType, roadConditionSurvey)
				if conditionValue == nil {
					continue
				}

				totalRange := math.Abs(roadConditionSurvey.KmStart - roadConditionSurvey.KmEnd)
				totalConditionValueByLane[roadCondition.LaneNo] += *conditionValue * totalRange

				processConditionGrades(roadCondition.LaneNo, totalRange, conditionGrades, conditionValue, roadConditionSurvey.SurveyType, summaryByConditionGrade, summaryConditionGradeByLane)

			}

			for _, roadConditionSurvey100M := range roadConditionSurvey.RoadConditionSurvey100Ms {
				if conditionGrades[0].RefOwner.RefConditionRangeID == 2 {
					conditionValue := getConditionValue(conditionType, roadConditionSurvey100M)
					if conditionValue == nil {
						continue
					}
					totalRange := math.Abs(roadConditionSurvey100M.KmStart - roadConditionSurvey100M.KmEnd)

					totalConditionValueByLane[roadCondition.LaneNo] += *conditionValue * totalRange

					processConditionGrades(roadCondition.LaneNo, totalRange, conditionGrades, conditionValue, roadConditionSurvey.SurveyType, summaryByConditionGrade, summaryConditionGradeByLane)

				}

				for _, roadConditionSurveyM := range roadConditionSurvey100M.RoadConditionSurveyMs {
					if conditionGrades[0].RefOwner.RefConditionRangeID == 1 {

						conditionValue := getConditionValue(conditionType, roadConditionSurveyM)
						if conditionValue == nil {
							continue
						}

						totalRange := math.Abs(roadConditionSurveyM.KmStart - roadConditionSurveyM.KmEnd)

						totalConditionValueByLane[roadCondition.LaneNo] += *conditionValue * totalRange

						processConditionGrades(roadCondition.LaneNo, totalRange, conditionGrades, conditionValue, roadConditionSurvey.SurveyType, summaryByConditionGrade, summaryConditionGradeByLane)

					}

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
		var details []responses.ConditionDetailKmStr

		for grade, value := range gradeMap {

			totalValue += value
			detailValue := value / km
			details = append(details, responses.ConditionDetailKmStr{
				RefGradeID:   grade.ID,
				RefGradeName: grade.Name,
				Value:        fmt.Sprintf("%.3f", detailValue),
			})
		}

		totalKm := totalValue / km

		avgValue := 0.0
		if len(details) > 0 {

			avgValue = totalConditionValueByLane[laneNo] / totalValue
		}

		for i := range details {

			detailsValue, _ := strconv.ParseFloat(details[i].Value, 64)

			valuePercent := (detailsValue / totalKm) * 100
			if totalValue != 0 {
				details[i].ValuePercent = fmt.Sprintf("%.3f", valuePercent)
			} else {
				details[i].ValuePercent = fmt.Sprintf("%.3f", 0.0)
			}
		}

		sort.Slice(details, func(i, j int) bool {
			return details[i].RefGradeID < details[j].RefGradeID
		})

		resp.Table = append(resp.Table, responses.ConditionTableStr{
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

func processConditionGrades(laneNo int, totalKm float64, conditionGrades []models.ParamsConditionPreload, conditionValue *float64, surveyType string, summaryByConditionGrade map[models.RefGrade]float64, summaryConditionGradeByLane map[int]map[models.RefGrade]float64) {

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

		} else {
			// Ensure that every grade is initialized in the map
			if _, exists := summaryByConditionGrade[conditionGrade.RefGrade]; !exists {
				summaryByConditionGrade[conditionGrade.RefGrade] = 0
			}
			if _, exists := summaryConditionGradeByLane[laneNo][conditionGrade.RefGrade]; !exists {
				summaryConditionGradeByLane[laneNo][conditionGrade.RefGrade] = 0
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

func (u *UseCase) GetDashboardConditionMap(roadIDs []int, filter requests.Condition, factor string) (GeomData, error) {

	conditionType := strings.ToLower(factor)

	var err error
	geomList, _, err := u.geomCondition(roadIDs, conditionType, filter.ConditionOwnerID, filter.Year)
	if err != nil {
		logs.Error(err)
		return GeomData{}, responses.NewAppErr(400, err.Error())
	}
	return geomList, nil

}

type GeomList struct {
	GradeID int      `json:"grade_id"`
	Color   string   `json:"color"`
	GeomCL  []string `json:"geom_cl"`
}

type GeomData struct {
	GeomList      []GeomList
	AllLineString string
	NumberOfData  int
}

func (u *UseCase) geomCondition(roadIDsForQuery []int, conditionType string, ownerID, year int) (GeomData, int, error) {
	dict := make(map[int][]responses.Geometry)
	geomCls := []GeomList{}
	allGeoJson := []IGeoJson{}
	result := GeomData{}
	refOwner, err := u.Repo.GetRoadConditionGradesByIDForRoadCondition(ownerID, conditionType)
	if err != nil {
		return result, 0, err
	}

	RoadConditions := []models.RoadCondition{}
	if len(roadIDsForQuery) != 0 {
		whereRoad := helpers.SliceIntToString(roadIDsForQuery)
		query := fmt.Sprintf("road_id IN (%s) AND year = %d AND status = 'A'", whereRoad, year)
		err := u.Repo.GetDataListForRoadCondition(&RoadConditions, query)
		if err != nil {
			logs.Error(err)
			return result, 0, responses.NewAppErr(400, err.Error())
		}
	} else {
		err = u.Repo.GetDataListForRoadCondition(&RoadConditions, fmt.Sprintf(" year = %d", year))
		if err != nil {
			logs.Error(err)
			return result, 0, responses.NewAppErr(400, err.Error())
		}
	}
	if len(RoadConditions) == 0 {
		return result, 0, nil
	}

	detailKm, err := u.initDetailKm(ownerID, conditionType)
	if err != nil {
		logs.Error(err)
		return result, 0, err
	}
	detailKms := make(map[int]responses.DetailKm)
	for _, item := range detailKm {
		detailKms[item.RefGrade.ID] = item
	}

	var roadSurveys []models.RoadConditionSurvey
	err = u.Repo.GetDataListForRoadCondition(&roadSurveys, "")
	if err != nil {
		logs.Error(err)
		return result, 0, responses.NewAppErr(400, err.Error())
	}

	roadSurvey100Ms, err := u.Repo.GetRoadCondition100MForRoadCondition("")
	if err != nil {
		logs.Error(err)
		return result, 0, responses.NewAppErr(400, err.Error())
	}

	roadSurveyMs, err := u.Repo.GetRoadConditionForRoadCondition("")
	if err != nil {
		logs.Error(err)
		return result, 0, responses.NewAppErr(400, err.Error())
	}

	// RoadByID
	var RoadByID []models.Road
	err = u.Repo.GetDataListForRoadCondition(&RoadByID, "")
	if err != nil {
		logs.Error(err)
		return result, 0, responses.NewAppErr(400, err.Error())
	}

	for _, roadCondition := range RoadConditions {
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
				if geoString.Type == "Point" {
					continue
				}
				logs.Error(err)
				return result, 0, err
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
				if len(item.Coordinates) == 0 {
					continue
				}
				var geoJson IGeoJson
				geoJson.Coordinates = item.Coordinates
				geoJsons = append(geoJsons, geoJson)
				allGeoJson = append(allGeoJson, geoJson)
			}
			result := mergeLineStrings(geoJsons)
			var geomCl GeomList
			geomCl.GradeID = detailKms[key].RefGrade.ID
			geomCl.Color = detailKms[key].RefGrade.Color
			geomCl.GeomCL = result
			geomCls = append(geomCls, geomCl)
		}

		// }
	}
	result.GeomList = geomCls
	result.NumberOfData = len(allGeoJson)
	return result, len(RoadConditions), nil
}

func (u *UseCase) initDetailKm(ownerID int, conditionType string) ([]responses.DetailKm, error) {
	var detailKms []responses.DetailKm
	paramCon := []models.ParamsCondition{}
	err := u.Repo.GetDataListForRoadCondition(&paramCon, fmt.Sprintf("ref_owner_id = %d AND LOWER(condition_type) = '%s'", ownerID, conditionType))
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
	err = u.Repo.GetDataListForRoadCondition(&refGrades, fmt.Sprintf("id IN (%s)", idString))
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

func FindGradeID(value float64, conditionType string, refOwner []models.ParamsConditionPreload) int {
	for _, item := range refOwner {
		if checkCondition(value, conditionType, item) {
			return item.RefGradeID
		}
	}
	return 0

}

type IGeoJson struct {
	Coordinates [][]float64
	Connections []int // Indices of the connected lines
}

func mergeLineStrings(geojsonData []IGeoJson) []string {
	mergedLines := []string{}
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

		// Generate the WKT LINESTRING representation
		wkt := "LINESTRING("
		for i, coord := range currentLine {
			if i > 0 {
				wkt += ","
			}
			wkt += fmt.Sprintf("%f %f", coord[0], coord[1])
		}
		wkt += ")"

		mergedLines = append(mergedLines, wkt)

		processedIndices[index] = true
	}

	return mergedLines
}

func (u *UseCase) Report6(year, roadSectionID, typ, dis string) (interface{}, error) {

	y, err := strconv.Atoi(year)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}

	SectionID, err := strconv.Atoi(roadSectionID)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	result := responses.Report6{}

	roadIDs, err := u.Repo.GetRoadFromSectionIDForRoadCondition(SectionID)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	for _, r := range roadIDs {
		data, lane, detail, err := u.Repo.GetReportConditionForRoadCondition(y, r, dis)
		if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
			logs.Error(err)
			return responses.RoadConditionDetails{}, err
		}

		if data == nil || len(lane) == 0 || len(detail) == 0 {

			roadInfo, err := u.Repo.GetRoadInfoForRoadCondition(r)
			if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
				logs.Error(err)
				return responses.RoadConditionDetails{}, err
			}
			// calulate road distance
			length := math.Abs(float64(roadInfo.KmEnd) - float64(roadInfo.KmStart))
			roadInfo.RoadLengthStr = fmt.Sprintf(`%.3f`, length/1000)

			newData := models.DataReportCondition{
				RoadGroupName: roadInfo.RoadGroupName,
				RoadName:      roadInfo.RoadName,
				RoadCode:      roadInfo.RoadCode,
				StrKmStart:    helpers.FormatKM(int64(roadInfo.KmStart)),
				StrKmEnd:      helpers.FormatKM(int64(roadInfo.KmEnd)),
				RoadLengthStr: roadInfo.RoadLengthStr,
				IsNull:        true,

				Lane: []models.DataRoadCondition{
					{
						LaneNo:  0,
						KmStart: 0,
						KmEnd:   0,
						Iri:     0,
						Rut:     0,
						Mpd:     0,
						Ifi:     0,

						RoadLengthStr: "",

						Detail: []models.DataRoadConditionM{
							{
								LaneNo: 0,

								KmStart: 0,
								KmEnd:   0,
								Iri:     0,
								Rut:     0,
								Mpd:     0,
								Ifi:     0,
							},
						},
					},
				},
			}
			data = &newData
		} else {
			// convert dis 25 to dis 100
			if dis == "100" {
				var newDetail []models.DataRoadConditionM
				count := 1
				iri, rut, mpd, ifi := 0.0, 0.0, 0.0, 0.0

				for iIndex, i := range detail {

					if count <= 4 {
						iri += i.Iri
						rut += i.Rut
						mpd += i.Mpd
						ifi += i.Ifi
					}

					if count == 4 {
						var r models.DataRoadConditionM

						if i.KmStart > i.KmEnd {
							r = models.DataRoadConditionM{
								LaneNo:     i.LaneNo,
								KmStart:    i.KmStart,
								KmEnd:      detail[iIndex-3].KmEnd,
								StrKmStart: helpers.FormatKM(int64(i.KmStart)),
								StrKmEnd:   helpers.FormatKM(int64(detail[iIndex-3].KmEnd)),
								Iri:        math.Round(iri/4*100) / 100,
								Rut:        math.Round(rut/4*100) / 100,
								Mpd:        math.Round(mpd/4*100) / 100,
								Ifi:        math.Round(ifi/4*100) / 100,
							}
						} else {
							r = models.DataRoadConditionM{
								LaneNo:     i.LaneNo,
								KmStart:    detail[iIndex-3].KmStart,
								KmEnd:      i.KmEnd,
								StrKmStart: helpers.FormatKM(int64(detail[iIndex-3].KmStart)),
								StrKmEnd:   helpers.FormatKM(int64(i.KmEnd)),
								Iri:        math.Round(iri/4*100) / 100,
								Rut:        math.Round(rut/4*100) / 100,
								Mpd:        math.Round(mpd/4*100) / 100,
								Ifi:        math.Round(ifi/4*100) / 100,
							}
						}

						newDetail = append(newDetail, r)
						count, iri, rut, mpd, ifi = 0, 0.0, 0.0, 0.0, 0.0

					} else if iIndex == len(detail)-1 || i.LaneNo != detail[iIndex+1].LaneNo { //the data can not divide by 4
						var r models.DataRoadConditionM

						if i.KmStart > i.KmEnd {
							r = models.DataRoadConditionM{
								LaneNo:     i.LaneNo,
								KmStart:    i.KmStart,
								KmEnd:      detail[iIndex-count+1].KmEnd,
								StrKmStart: helpers.FormatKM(int64(i.KmStart)),
								StrKmEnd:   helpers.FormatKM(int64(detail[iIndex-count+1].KmEnd)),
								Iri:        math.Round(iri/float64(count)*100) / 100,
								Rut:        math.Round(rut/float64(count)*100) / 100,
								Mpd:        math.Round(mpd/float64(count)*100) / 100,
								Ifi:        math.Round(ifi/float64(count)*100) / 100,
							}
						} else {
							r = models.DataRoadConditionM{
								LaneNo:     i.LaneNo,
								KmStart:    detail[iIndex-count+1].KmStart,
								KmEnd:      i.KmEnd,
								StrKmStart: helpers.FormatKM(int64(detail[iIndex-count+1].KmStart)),
								StrKmEnd:   helpers.FormatKM(int64(i.KmEnd)),
								Iri:        math.Round(iri/float64(count)*100) / 100,
								Rut:        math.Round(rut/float64(count)*100) / 100,
								Mpd:        math.Round(mpd/float64(count)*100) / 100,
								Ifi:        math.Round(ifi/float64(count)*100) / 100,
							}
						}

						newDetail = append(newDetail, r)
						count, iri, rut, mpd, ifi = 0, 0.0, 0.0, 0.0, 0.0
					}

					count++
				}

				detail = newDetail

			} else {

				for iIndex, i := range detail {
					detail[iIndex].StrKmStart = helpers.FormatKM(int64(i.KmStart))
					detail[iIndex].StrKmEnd = helpers.FormatKM(int64(i.KmEnd))
				}
			}

			// calulate road distance
			length := math.Abs(float64(data.KmEnd) - float64(data.KmStart))
			data.RoadLengthStr = helpers.FormatNumberFloat(length / 1000)
			data.StrKmStart = helpers.FormatKM(int64(data.KmStart))
			data.StrKmEnd = helpers.FormatKM(int64(data.KmEnd))

			//insert detail to lane
			for index, i := range lane {
				// calulate lane distance
				lane[index].RoadLengthStr = helpers.FormatNumber(int(math.Abs(float64(i.KmEnd) - float64(i.KmStart))))

				if i.KmStart > i.KmEnd {

					sortDetail := detail
					sort.Slice(sortDetail, func(i, j int) bool {
						return sortDetail[i].KmStart > sortDetail[j].KmStart
					})

					detail = sortDetail
				}

				for _, j := range detail {
					if i.LaneNo == j.LaneNo {
						lane[index].Detail = append(lane[index].Detail, j)
					}
				}
			}

			data.Lane = lane
		}
		result.Data = append(result.Data, *data)
	}
	roadSectionInfo, err := u.Repo.GetRoadSectionByIDForRoadCondition(SectionID)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	roadGroup, err := u.Repo.GetRoadGroupByIDForRoadCondition(roadSectionInfo.RoadGroupId)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}

	result.RoadGroupName = roadGroup.Number
	result.RoadSectionNumber = roadSectionInfo.Number
	result.RoadSectionName = roadSectionInfo.NameOriginTH + " - " + roadSectionInfo.NameDestinationTH
	result.KmStart = helpers.FormatKM(int64(roadSectionInfo.KmStart))
	result.KmEnd = helpers.FormatKM(int64(roadSectionInfo.KmEnd))
	result.RoadLength = fmt.Sprintf(`%.3f`, roadSectionInfo.Distance)
	result.Year = y + 543
	// generate report
	var pathResult interface{}

	if typ == "excel" {
		pathResult, err = ExportExcelType6(result)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	} else {
		pathResult, err = helpers.RequestExport(result, "TEMPLATE_GENARAL_TYPE6_AEK", typ)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	}

	return pathResult, nil
}

func ExportExcelType6(data responses.Report6) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE6_AEK_EXCEL"))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// file, err := os.Open(fmt.Sprint(logoPath))
	// if err != nil {
	// 	fmt.Println("can not open file", err)
	// 	return nil, err
	// }
	// defer file.Close()
	templateIndex, _ := f.GetSheetIndex("Sheet1")
	sheetList := make(map[string]int)

	// Initialize a map to keep track of sheets that already have the picture
	for _, v := range data.Data {
		logoPath := os.Getenv("LOGO")

		fileLogo, err := os.Open(fmt.Sprint(logoPath))
		if err != nil {
			return nil, err
		}
		defer fileLogo.Close()

		imageLogoConfig, _, err := image.DecodeConfig(fileLogo)
		if err != nil {
			return nil, err
		}

		// Create a new sheet by copying the template sheet
		sheetName := v.RoadName
		if utf8.RuneCountInString(sheetName) > 31 {
			runes := []rune(sheetName)
			sheetName = string(runes[:25])
		}
		if count, exists := sheetList[sheetName]; exists {
			sheetList[sheetName]++
			if len(sheetName) > 75 {
				sheetName = sheetName[:75]
			}
			sheetName = fmt.Sprintf("%s (%d)", sheetName, count)
			newSheetIndex, err := f.NewSheet(sheetName)
			if err != nil {
				fmt.Println("Error creating new sheet:", err)
				return nil, err
			}

			err = f.CopySheet(templateIndex, newSheetIndex)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		} else {
			newSheetIndex, err := f.NewSheet(sheetName)
			if err != nil {
				fmt.Println("Error creating new sheet:", err)
				return nil, err
			}

			err = f.CopySheet(templateIndex, newSheetIndex)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			sheetList[sheetName]++
		}

		f.SetCellValue(sheetName, "C3", fmt.Sprintf("หมายเลขทางหลวง : %s ตอนควบคุม : %s ", data.RoadGroupName, data.RoadSectionNumber))
		f.SetCellValue(sheetName, "C4", "ชื่อสายทาง : "+v.RoadName)
		f.SetCellValue(sheetName, "C5", "กม.เริ่มต้น "+v.StrKmStart+" กม.สิ้นสุด "+v.StrKmEnd+" ระยะทาง "+v.RoadLengthStr+" กม.")
		f.SetCellValue(sheetName, "B6", "ชื่อสายทาง : "+v.RoadName)
		if !v.IsNull {
			Row := 9
			StrRow := fmt.Sprint(Row)

			for iIndex, i := range v.Lane {
				f.SetCellValue(sheetName, "B"+StrRow, i.LaneNo)
				f.SetCellValue(sheetName, "C"+StrRow, i.RoadLengthStr)
				f.SetCellValue(sheetName, "D"+StrRow, "")
				f.SetCellValue(sheetName, "E"+StrRow, "")
				f.SetCellValue(sheetName, "F"+StrRow, "")
				f.SetCellValue(sheetName, "G"+StrRow, fmt.Sprintf(`%.3f`, i.Iri))
				f.SetCellValue(sheetName, "H"+StrRow, "")
				f.SetCellValue(sheetName, "I"+StrRow, fmt.Sprintf(`%.3f`, i.Rut))
				f.SetCellValue(sheetName, "J"+StrRow, "")
				f.SetCellValue(sheetName, "K"+StrRow, fmt.Sprintf(`%.3f`, i.Mpd))
				f.SetCellValue(sheetName, "L"+StrRow, "")
				f.SetCellValue(sheetName, "M"+StrRow, fmt.Sprintf(`%.3f`, i.Ifi))

				f.DuplicateRow(sheetName, Row)
				Row++
				StrRow = fmt.Sprint(Row)

				for jIndex, j := range i.Detail {
					f.SetCellValue(sheetName, "B"+StrRow, "")
					f.SetCellValue(sheetName, "C"+StrRow, "")
					f.SetCellValue(sheetName, "D"+StrRow, j.StrKmStart)
					f.SetCellValue(sheetName, "E"+StrRow, j.StrKmEnd)
					f.SetCellValue(sheetName, "F"+StrRow, fmt.Sprintf(`%.3f`, j.Iri))
					f.SetCellValue(sheetName, "G"+StrRow, "")
					f.SetCellValue(sheetName, "H"+StrRow, fmt.Sprintf(`%.3f`, j.Rut))
					f.SetCellValue(sheetName, "I"+StrRow, "")
					f.SetCellValue(sheetName, "J"+StrRow, fmt.Sprintf(`%.3f`, j.Mpd))
					f.SetCellValue(sheetName, "K"+StrRow, "")
					f.SetCellValue(sheetName, "L"+StrRow, fmt.Sprintf(`%.3f`, j.Ifi))
					f.SetCellValue(sheetName, "M"+StrRow, "")

					if jIndex < len(i.Detail)-1 {
						f.DuplicateRow(sheetName, Row)
						Row++
						StrRow = fmt.Sprint(Row)
					}
				}

				if iIndex < len(v.Lane)-1 {
					f.DuplicateRow(sheetName, Row)
					Row++
					StrRow = fmt.Sprint(Row)
				}
			}
		}

		// Define the target dimensions in inches
		targetWidthInches := 1.55
		targetHeightInches := 1.0

		// Convert inches to pixels (assuming 96 DPI)
		dpi := 96.0
		targetWidthPixels := targetWidthInches * dpi
		targetHeightPixels := targetHeightInches * dpi

		// Calculate the scaling factors based on the image's original dimensions
		scaleX := targetWidthPixels / float64(imageLogoConfig.Width)
		scaleY := targetHeightPixels / float64(imageLogoConfig.Height)

		// Add the picture with the calculated scaling factors
		err = f.AddPicture(sheetName, "B2", fmt.Sprint(logoPath), &excelize.GraphicOptions{
			ScaleX: scaleX,
			ScaleY: scaleY,
		})
		if err != nil {
			fmt.Println("Error adding picture:", err)
			return nil, err
		}

		helpers.AddFooter(f, sheetName)
	}
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}
	f.DeleteSheet("Sheet1")
	reportName, err := helpers.ReportName("TEMPLATE_GENARAL_TYPE6_AEK_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	err = f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")
	if err != nil {
		return nil, err
	}

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}
