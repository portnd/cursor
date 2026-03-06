package usecases

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/xuri/excelize/v2"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

func (u *UseCase) GetReportType5(id, userID int, plan int, typeFile string) (interface{}, error) {
	filePath := os.Getenv("MAINTENANCE_ANALYSIS_PDF") + fmt.Sprintf("%d/", id)
	templateName := os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE5")

	roads, err := u.Repo.GetMaintenanceAnalysisRoadIdById(id)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	analysis, err := u.Repo.GetMaintenanceAnalysisById(id)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	analysisTypeID := analysis.MaintenanceAnalysisTypeId
	if analysisTypeID == 2 {
		plan = 1
	}

	var refCriteriaMethods []models.RefCriteriaMethod
	colorCriteria := make(map[string]string)
	dataAllYear, err := u.Repo.DashboardMapAllYear(id, analysis.MaintenanceAnalysisTypeId, &plan)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return dataAllYear, responses.NewAppErr(400, err.Error())
	}

	colorIndex := 0
	for _, item := range dataAllYear {
		if item.Data.Repair {
			//set color\
			if colorCriteria[item.Data.IcResult.Name] == "" {
				colorCriteria[item.Data.IcResult.Name] = colors[colorIndex]
				refCriteriaMethods = append(refCriteriaMethods, models.RefCriteriaMethod{Name: item.Data.IcResult.Name, Color: colors[colorIndex]})
				colorIndex++
			}
		}
	}

	helpers.PrintlnJson(refCriteriaMethods)

	years, err := u.Repo.GetMaintenanceAnalysisResultGroupByYear(id, plan, analysisTypeID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	var roadGroup []models.RoadGroup
	var roadSection []models.RoadSection
	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		var err error
		roadGroup, err = u.Repo.GetRoadGroups()
		return err
	})
	g.Go(func() error {
		var err error
		roadSection, err = u.Repo.GetRoadSections()
		return err
	})
	if err := g.Wait(); err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	roadGroupMap := map[int]models.RoadGroup{}
	for _, item := range roadGroup {
		roadGroupMap[item.Id] = item
	}
	roadSectionMap := map[int]models.RoadSection{}
	for _, item := range roadSection {
		roadSectionMap[item.Id] = item
	}

	sort.Ints(years)
	round := int(math.Ceil(float64(len(years)) / 3))
	var res []responses.Report5Road

	// Fetch all roads' result data in one MongoDB query instead of one per road (much faster for HTML/PDF)
	roadIDs := make([]int, 0, len(roads))
	for _, r := range roads {
		roadIDs = append(roadIDs, r.RoadID)
	}
	resultByRoad, err := u.Repo.GetMaintenanceAnalysisResultByRoadIDs(id, roadIDs, plan, analysisTypeID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	for _, road := range roads {
		roadID := road.RoadID
		roadInfo, err := u.Repo.GetRoadInfoByID(roadID)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}

		road, err := u.Repo.GetRoadByID(roadID)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}

		// road, err := u.Repo.GetRoadByID(road.RoadID)
		// if err != nil {
		// 	return "", responses.NewAppErr(400, err.Error())
		// }

		// if road.RoadId != 3 {
		// 	continue
		// }
		yearMapData := make(map[int]int)
		for index, item := range years {
			yearMapData[index] = item
		}

		start := 0
		end := 0
		// j := 0
		rowData := []interface{}{}
		batchSize := 3
		if round > 1 {
			for i := 0; i < round; i++ {
				start = end //+= i + end - j
				// j = 1
				end = start + batchSize
				fmt.Println("endend", start, end)
				if end > len(years) {
					end = len(years)
				}
				var yearData []int
				yearData = years[start:end]
				roadResult := resultByRoad[roadID]
				if roadResult == nil {
					roadResult = []models.ModelResult{}
				}
				rows, err := u.GetRows(id, plan, analysisTypeID, roadInfo, yearData, colorCriteria, roadResult)
				if err != nil {
					logs.Error(err)
					return "", responses.NewAppErr(400, err.Error())
				}
				rowData = append(rowData, rows)
			}
		} else {
			if len(years) == 1 {
				start = 0
				end = 1
			} else if len(years) == 2 {
				start = 0
				end = 2
			} else if len(years) == 3 {
				start = 0
				end = 3
				helpers.PrintlnJson(years[start:end])
			}
			roadResult := resultByRoad[roadID]
			if roadResult == nil {
				roadResult = []models.ModelResult{}
			}
			rows, err := u.GetRows(id, plan, analysisTypeID, roadInfo, years[start:end], colorCriteria, roadResult)
			if err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, err.Error())
			}

			// return rows, nil
			rowData = append(rowData, rows)
		}

		var res22 responses.Report5Road
		res22.RoadID = road.RoadCode
		res22.RoadName = roadInfo.Name
		res22.Years = years
		res22.RoadGroupCode = roadGroupMap[road.RoadGroupId].Number
		res22.RoadSectionCode = roadSectionMap[road.RoadSectionId].Number
		res22.Rows = rowData
		res = append(res, res22)
	}
	yearStart := fmt.Sprintf("%d", years[0]+543)
	yearEnd := fmt.Sprintf("%d", years[len(years)-1]+543)
	condition, err := u.Repo.GetMaintenanceAnalysisStrategicBudgetTypeById(analysis.Condition)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	target, err := u.Repo.GetMaintenanceAnalysisStrategiTargetById(*analysis.Target)
	if err != nil {
		logs.Error(err)
		// return "", responses.NewAppErr(400, err.Error())
	}

	user, err := u.Repo.GetUserByID(uint(userID))
	// analysis.Condition = condition.Name
	// analysis.Target = 1
	var report5Res responses.Report5Res
	report5Res.ReportRoad = res
	if len(years) > 1 {
		report5Res.YearLength = yearStart + " - " + yearEnd
	} else {
		report5Res.YearLength = yearStart
	}
	layout := "2006-01-02 15:04:05.999999 -0700 -0700"

	// Parse the input time with the given layout
	t, err := time.Parse(layout, analysis.CreatedAt.String())
	if err != nil {
		fmt.Println("Error:", err)
		// return
	}

	// Set the time zone to Asia/Bangkok (Thailand)
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println("Error:", err)
		// return
	}
	t = t.In(loc)

	report5Res.Plan = plan
	report5Res.Condition = condition.Name
	report5Res.Target = target.Name
	report5Res.CreatedAt = helpers.ConvertToThaiFullCalendar(analysis.CreatedAt.Add(time.Hour * 7))
	report5Res.User = user.Firstname + " " + user.Lastname

	// report5Res.Target =
	// var rows responses.Res
	// rows.Data = res
	// rows.Rows = rowData
	// return report5Res, nil
	// return report5Res, nil
	// if typeFile == "html" {
	// 	filename := "type5_1.html"
	// 	result, err := helpers.InitDataToHtmlFunc(os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE5_1"), report5Res, filePath, filename)
	// 	if err != nil {
	// 		// helpers.PrintlnJson(err.Error())
	// 		logs.Error(err)
	// 		return "", err
	// 	}
	// 	return os.Getenv("STORAGE_IP") + "/" + result, nil
	// }

	filename := "type5.html"
	result, err := helpers.InitDataToHtmlFunc(templateName, report5Res, filePath, filename)
	if err != nil {
		// helpers.PrintlnJson(err.Error())
		logs.Error(err)
		return "", err
	}

	if typeFile == "html" {
		texts := strings.Split(result, "/")
		text := texts[len(texts)-1]

		new := strings.ReplaceAll(result, text, "แผนการดำเนินงานการปรับปรุงผิวทาง.html")

		os.Rename(result, new)

		return os.Getenv("STORAGE_IP") + "/" + new, nil
	}

	html, err := os.ReadFile(result)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	e := os.Remove(result)
	if e != nil {
		logs.Error(err)
		return "", err
	}

	var buf []byte
	ctx, cancel := helpers.NewChromedpContext(context.Background(), log.Printf)
	err = chromedp.Run(ctx, helpers.PrintToPDFWithDelay(string(html), &buf, 3, "div#success-pagejs"))
	cancel()
	if err != nil {
		logs.Error(err)
		return "", err
	}

	if err = os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return "", err
	}

	name := "แผนการดำเนินงานการปรับปรุงผิวทาง"

	fullFilePath := filePath + name + ".pdf"
	if err := ioutil.WriteFile(fullFilePath, buf, 0777); err != nil {
		logs.Error(err)
		return "", err
	}

	return os.Getenv("STORAGE_IP") + "/" + fullFilePath, nil
}

func (u *UseCase) Report5Excel(id, userID int, plan int, typeFile string) (interface{}, error) {

	// filePath := os.Getenv("MAINTENANCE_ANALYSIS_PDF")
	// templateName := os.Getenv("TEMPLATE_MAINTENANCE_ANALYSIS_TYPE5")

	roads, err := u.Repo.GetMaintenanceAnalysisRoadIdById(id)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	analysis, err := u.Repo.GetMaintenanceAnalysisById(id)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}
	analysisTypeID := analysis.MaintenanceAnalysisTypeId
	if analysisTypeID == 2 {
		plan = 1
	}

	var refCriteriaMethods []models.RefCriteriaMethod
	colorCriteria := make(map[string]string)
	dataAllYear, err := u.Repo.DashboardMapAllYear(id, analysis.MaintenanceAnalysisTypeId, &plan)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return dataAllYear, responses.NewAppErr(400, err.Error())
	}

	colorIndex := 0
	for _, item := range dataAllYear {
		if item.Data.Repair {
			//set color\
			if colorCriteria[item.Data.IcResult.Name] == "" {
				colorCriteria[item.Data.IcResult.Name] = colors[colorIndex]
				refCriteriaMethods = append(refCriteriaMethods, models.RefCriteriaMethod{Name: item.Data.IcResult.Name, Color: colors[colorIndex]})
				colorIndex++
			}
		}
	}

	years, err := u.Repo.GetMaintenanceAnalysisResultGroupByYear(id, plan, analysisTypeID)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	var roadGroup []models.RoadGroup
	var roadSection []models.RoadSection
	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		var err error
		roadGroup, err = u.Repo.GetRoadGroups()
		return err
	})
	g.Go(func() error {
		var err error
		roadSection, err = u.Repo.GetRoadSections()
		return err
	})
	if err := g.Wait(); err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	roadGroupMap := map[int]models.RoadGroup{}
	for _, item := range roadGroup {
		roadGroupMap[item.Id] = item
	}
	roadSectionMap := map[int]models.RoadSection{}
	for _, item := range roadSection {
		roadSectionMap[item.Id] = item
	}

	sort.Ints(years)
	var res []responses.Report5Road
	// var rowData []responses.KMStr
	rowData := []interface{}{}
	for _, road := range roads {

		roadInfo, err := u.Repo.GetRoadInfoByID(road.RoadID)
		if err != nil {
			return "", responses.NewAppErr(400, err.Error())
		}

		road, err := u.Repo.GetRoadByID(road.RoadID)
		if err != nil {
			return "", responses.NewAppErr(400, err.Error())
		}

		rows, err := u.GetRows(id, plan, analysisTypeID, roadInfo, years, colorCriteria, nil)
		if err != nil {
			return "", responses.NewAppErr(400, err.Error())
		}
		rowData = append(rowData, rows)

		var res22 responses.Report5Road
		res22.RoadID = road.RoadCode
		res22.RoadName = roadInfo.Name
		res22.RoadGroupCode = roadGroupMap[road.RoadGroupId].Number
		res22.RoadSectionCode = roadSectionMap[road.RoadSectionId].Number
		res22.Years = years
		res22.Rows = rowData[0]
		res = append(res, res22)
	}

	yearStart := fmt.Sprintf("%d", years[0]+543)
	yearEnd := fmt.Sprintf("%d", years[len(years)-1]+543)
	condition, err := u.Repo.GetMaintenanceAnalysisStrategicBudgetTypeById(analysis.Condition)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	target, err := u.Repo.GetMaintenanceAnalysisStrategiTargetById(*analysis.Target)
	if err != nil {
		// return "", responses.NewAppErr(400, err.Error())
	}

	user, err := u.Repo.GetUserByID(uint(userID))
	// analysis.Condition = condition.Name
	// analysis.Target = 1
	var report5Res responses.Report5Res
	report5Res.ReportRoad = res
	if len(years) > 1 {
		report5Res.YearLength = yearStart + " - " + yearEnd
	} else {
		report5Res.YearLength = yearStart
	}
	layout := "2006-01-02 15:04:05.999999 -0700 -0700"

	// Parse the input time with the given layout
	t, err := time.Parse(layout, analysis.CreatedAt.String())
	if err != nil {
		fmt.Println("xxxxxxxxxxx", analysis.CreatedAt)
		fmt.Println("Error:", err)
		// return
	}

	// Set the time zone to Asia/Bangkok (Thailand)
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println("yyyyyyyyy")
		fmt.Println("Error:", err)
		// return
	}
	t = t.In(loc)

	// Format the time as desired
	formattedTime := t.Format("2006-01-02 15:04:05")

	report5Res.Plan = plan
	report5Res.Condition = condition.Name
	report5Res.Target = target.Name
	report5Res.CreatedAt = formattedTime
	report5Res.User = user.Firstname + " " + user.Lastname
	filePath := os.Getenv("ANALYSIS_EXCEL")
	template := "TEMPLATE_MAINTENANCE_ANALYSIS_TYPE5_" + fmt.Sprint(len(years)) + "_EXCEL"
	f, err := excelize.OpenFile(os.Getenv(template))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for i := 1; i <= 512; i++ {
		if i > len(report5Res.ReportRoad) {
			f.DeleteSheet("Sheet1 " + "(" + fmt.Sprint(i) + ")")
		}
	}

	textLeft, err := f.NewStyle(
		&excelize.Style{
			Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true, ShrinkToFit: true},
			Font:      &excelize.Font{Bold: true, Size: 12, Family: "TH SarabunPSK"},
		},
	)
	if err != nil {
		return nil, err
	}

	for index, item := range report5Res.ReportRoad {
		plan := fmt.Sprint(report5Res.Plan)
		if plan == "0" {
			plan = "(แผนไม่จำกัดงบประมาณ)"
		} else {
			plan = "(แผนที่ " + plan + ")"
		}
		f.SetCellValue("Sheet1 "+"("+fmt.Sprint(index+1)+")", "D4", "แผนการดําเนินงานการปรับปรุงผิวทาง"+plan+" พ.ศ."+report5Res.YearLength)
		f.SetCellValue("Sheet1 "+"("+fmt.Sprint(index+1)+")", "D5", "")
		f.UnmergeCell("Sheet1 "+"("+fmt.Sprint(index+1)+")", "B5", "I7")
		f.SetCellValue("Sheet1 "+"("+fmt.Sprint(index+1)+")", "B5", " ")
		f.SetCellValue("Sheet1 "+"("+fmt.Sprint(index+1)+")", "B7", " ")

		f.SetColWidth("Sheet1 "+"("+fmt.Sprint(index+1)+")", "A1", "A1", 20.0)

		f.MergeCell("Sheet1 "+"("+fmt.Sprint(index+1)+")", "B6", "E6")
		f.SetRowHeight("Sheet1 "+"("+fmt.Sprint(index+1)+")", 6, 120)
		f.SetCellStyle("Sheet1 "+"("+fmt.Sprint(index+1)+")", "B6", "B6", textLeft)

		f.SetCellValue("Sheet1 "+"("+fmt.Sprint(index+1)+")", "B6", "เงื่อนไข : "+report5Res.Condition+"\nเป้าหมาย : "+report5Res.Target+"\nผู้ประมวลผล : "+user.Firstname+" "+user.Lastname+"\nวันที่วิเคราะห์ : "+helpers.ConvertToThaiFullCalendar(analysis.CreatedAt.Add(time.Hour*7))+" น."+"\nหมายเลขทางหลวง : "+item.RoadGroupCode+"\nตอนควบคุม : "+item.RoadSectionCode)

		f.SetCellValue("Sheet1 "+"("+fmt.Sprint(index+1)+")", "B7", "ชื่อสายทาง : "+item.RoadName)

		yearColumn := []string{"D", "J", "P", "V", "AB", "AH", "AN", "AT", "AZ", "BF"}
		laneColumn := [][]string{
			{"D", "F", "H"},
			{"J", "L", "N"},
			{"P", "R", "T"},
			{"V", "X", "Z"},
			{"AB", "AD", "AF"},
			{"AH", "AJ", "AL"},
			{"AN", "AP", "AR"},
			{"AT", "AV", "AX"},
			{"AZ", "BB", "BD"},
			{"BF", "BH", "BJ"}}

		for yIndex, y := range item.Years {
			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(index+1)+")", yearColumn[yIndex]+"8", y+543)
		}

		for jIndex, j := range item.Rows.([]responses.KMStr) {

			if jIndex < len(item.Rows.([]responses.KMStr))-1 {
				f.DuplicateRow("Sheet1 "+"("+fmt.Sprint(index+1)+")", 12+jIndex)
			}

			f.SetCellValue("Sheet1 "+"("+fmt.Sprint(index+1)+")", "B"+fmt.Sprint(12+jIndex), fmt.Sprint(j.KmStart)+"-"+fmt.Sprint(j.KmEnd))
			for yearIndex, lane := range j.Data.([]responses.YearData) {

				for laneIndex, laneData := range lane.Data.([]responses.LaneData) {
					// report := laneData.Data //data.Data.(responses.Report5)
					if laneData.Data.KMStart != "" {
						f.SetCellRichText("Sheet1 "+"("+fmt.Sprint(index+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex), []excelize.RichTextRun{
							{
								Text: fmt.Sprint(helpers.FormatKM(int64(laneData.Data.Data.KmStart))) + " - " + fmt.Sprint(helpers.FormatKM(int64(laneData.Data.Data.KmEnd))) + "\n",
							},
							{
								Text: laneData.Data.Data.InterventionCriteriaName + "\n",
							},
							{
								Text: laneData.Data.Data.Method,
							},
						})
					}
					// f.SetCellStyle("Sheet1 "+"("+fmt.Sprint(index+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex-1), laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex), wrapTextStyle1)
					preCell, _ := f.GetCellValue("Sheet1 "+"("+fmt.Sprint(index+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex-1))
					currentCell, _ := f.GetCellValue("Sheet1 "+"("+fmt.Sprint(index+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex))

					if preCell != "" && currentCell != "" && preCell == currentCell {

						f.MergeCell("Sheet1 "+"("+fmt.Sprint(index+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex-1), laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex))

						// style and wrap text
						wrapTextStyle, _ := f.NewStyle(&excelize.Style{
							Font: &excelize.Font{
								Size:   10,
								Family: "TH SarabunPSK",
							},
							Alignment: &excelize.Alignment{
								WrapText:   true,
								Horizontal: "center",
								Vertical:   "center",
							},
							Fill: excelize.Fill{
								Type:    "pattern",
								Color:   []string{laneData.Data.Data.Color},
								Pattern: 1,
							},
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

						f.SetCellStyle("Sheet1 "+"("+fmt.Sprint(index+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex-1), laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex), wrapTextStyle)
					} else if currentCell != "" {
						// style and wrap text
						wrapTextStyle, _ := f.NewStyle(&excelize.Style{
							Font: &excelize.Font{
								Size:   10,
								Family: "TH SarabunPSK",
							},
							Alignment: &excelize.Alignment{
								WrapText:   true,
								Horizontal: "center",
								Vertical:   "center",
							},
							Fill: excelize.Fill{
								Type:    "pattern",
								Color:   []string{laneData.Data.Data.Color},
								Pattern: 1,
							},
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

						f.SetCellStyle("Sheet1 "+"("+fmt.Sprint(index+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex), laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex), wrapTextStyle)

						// if laneData.Data.Span == 1 {
						// 	f.SetRowHeight("Sheet1 "+"("+fmt.Sprint(index+1)+")", 12+jIndex, 60)
						// }

					} else {

						// f.SetCellStyle("Sheet1 "+"("+fmt.Sprint(index+1)+")", laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex), laneColumn[yearIndex][laneIndex]+fmt.Sprint(12+jIndex), wrapTextStyle)

					}
				}
			}
		}

		f.SetSheetName("Sheet1 "+"("+fmt.Sprint(index+1)+")", item.RoadName)
	}

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	name := "แผนการดำเนินงานการปรับปรุงผิวทาง"

	f.SaveAs(os.Getenv("ANALYSIS_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}

func (u *UseCase) GetRows(id, plan, analysisTypeID int, roadData models.RoadInfo, years []int, colorCriteria map[string]string, preFetchedResult []models.ModelResult) (interface{}, error) {
	roadID := roadData.RoadId

	kmStart := roadData.KmStart
	kmEnd := roadData.KmEnd
	directionID := roadData.RefDirectionId
	var kms []responses.KMStr

	var resultData []models.ModelResult
	if preFetchedResult != nil {
		resultData = preFetchedResult
	} else {
		resultData, _ = u.Repo.GetMaintenanceAnalysisResult(id, roadID, plan, analysisTypeID)
	}
	// return resultData, nil
	// dataTemp := map[string]int{}
	cnt := 0
	// ffff := map[string]string{}
	spanData := map[string]int{}
	stack := Stack{}
	var test []responses.Report5
	if directionID == 1 {
		for i := kmStart; i < kmEnd; {
			// ////////////////////////////////////
			result, span := CalKm(roadID, directionID, int(kmStart), int(kmEnd), resultData)
			// helpers.PrintlnJson(result)
			// return result, nil
			// return result, nil
			if cnt == 0 {
				spanData = span
				stack.Push(span)
			} else {
				// helpers.PrintlnJson("stack.itemsstack.items", stack.items)
				// ddddd := stack.items[0].(map[string]int)

				// spanData = ddddd
			}
			cnt++
			////////////////////////////////
			var yearData []responses.YearData
			for _, year := range years {
				if year == 0 {
					continue
				}
				var laneDatas []responses.LaneData
				//////// lane 1 ////////
				data1, temp, err := u.GetReportTypeData(result, roadID, directionID, int(kmStart), int(kmEnd), int(i), int(i+200), 1, year, spanData, colorCriteria)
				if err != nil {
					return "", responses.NewAppErr(400, err.Error())
				}

				if len(data1.([]responses.Report5)) > 0 {
					for _, item := range data1.([]responses.Report5) {
						test = append(test, item)
					}
					// return data1.([]responses.Report5), nil
				}

				spanData = temp
				// helpers.PrintlnJson("===", spanData)
				tempLane1 := make(map[string]interface{})
				for _, item := range data1.([]responses.Report5) {
					if item.Data.InterventionCriteriaName == "" {
						continue
					}
					str1 := fmt.Sprintf("%s-%s-%d-%d", item.KMStart, item.KMEnd, item.LaneNo, item.Year)
					tempLane1[str1] = item
				}
				lanes1 := []responses.Report5{}
				for key := range tempLane1 {
					value, ok := tempLane1[key].(responses.Report5)
					if ok {
						lanes1 = append(lanes1, value)
					}
				}
				if len(lanes1) == 0 {
					var laneData1 responses.LaneData
					laneData1.LaneNo = 1
					laneData1.Data = responses.Report5{}
					laneDatas = append(laneDatas, laneData1)
				} else {
					var laneData1 responses.LaneData
					laneData1.LaneNo = 1
					laneData1.Data = lanes1[0]
					laneDatas = append(laneDatas, laneData1)
				}

				//////// lane 2 ////////
				data2, _, err := u.GetReportTypeData(result, roadID, directionID, int(kmStart), int(kmEnd), int(i), int(i+200), 2, year, spanData, colorCriteria)
				if err != nil {
					return "", responses.NewAppErr(400, err.Error())
				}
				tempLane2 := make(map[string]interface{})
				for _, item := range data2.([]responses.Report5) {
					if item.Data.InterventionCriteriaName == "" {
						continue
					}
					str := fmt.Sprintf("%s-%s-%d-%d", item.KMStart, item.KMEnd, item.LaneNo, item.Year)
					tempLane2[str] = item
				}
				lanes2 := []responses.Report5{}
				for key := range tempLane2 {
					value, ok := tempLane2[key].(responses.Report5)
					if ok {
						lanes2 = append(lanes2, value)
					}
				}
				if len(lanes2) == 0 {
					var laneData2 responses.LaneData
					laneData2.LaneNo = 2
					laneData2.Data = responses.Report5{}
					laneDatas = append(laneDatas, laneData2)
				} else {
					var laneData2 responses.LaneData
					laneData2.LaneNo = 2
					laneData2.Data = lanes2[0]
					laneDatas = append(laneDatas, laneData2)
				}

				//////// lane 3 ////////
				data3, _, err := u.GetReportTypeData(result, roadID, directionID, int(kmStart), int(kmEnd), int(i), int(i+200), 3, year, spanData, colorCriteria)
				if err != nil {
					return "", responses.NewAppErr(400, err.Error())
				}
				tempLane3 := make(map[string]interface{})
				for _, item := range data3.([]responses.Report5) {
					if item.Data.InterventionCriteriaName == "" {
						continue
					}
					str := fmt.Sprintf("%s-%s-%d-%d", item.KMStart, item.KMEnd, item.LaneNo, item.Year)
					tempLane3[str] = item
				}
				lanes3 := []responses.Report5{}
				for key := range tempLane3 {
					value, ok := tempLane3[key].(responses.Report5)
					if ok {
						lanes3 = append(lanes3, value)
					}
				}
				if len(lanes3) == 0 {
					var laneData3 responses.LaneData
					laneData3.LaneNo = 3
					laneData3.Data = responses.Report5{}
					laneDatas = append(laneDatas, laneData3)
				} else {
					var laneData3 responses.LaneData
					laneData3.LaneNo = 3
					laneData3.Data = lanes3[0]
					laneDatas = append(laneDatas, laneData3)
				}

				yearData = append(yearData, responses.YearData{Year: year, Data: laneDatas})

			}
			kms = append(kms, responses.KMStr{Data: yearData, KmStart: helpers.FormatKM(int64(i)), KmEnd: helpers.FormatKM(int64(i + 200))})
			i += 200
		}

	} else {

		for i := kmStart; i > kmEnd; {
			// ////////////////////////////////////
			result, span := CalKm(roadID, directionID, int(kmStart), int(kmEnd), resultData)

			if cnt == 0 {
				spanData = span
				stack.Push(span)
			} else {
				// spanData = span
				// helpers.PrintlnJson(spanData)
			}
			cnt++
			// helpers.PrintlnJson(dataTemp)
			// return result, nil
			////////////////////////////////
			var yearData []responses.YearData
			for _, year := range years {
				if year == 0 {
					continue
				}
				var laneDatas []responses.LaneData
				//////// lane 1 ////////
				data1, temp, err := u.GetReportTypeData(result, roadID, directionID, int(kmStart), int(kmEnd), int(i), int(i-200), 1, year, spanData, colorCriteria)
				if err != nil {
					return "", responses.NewAppErr(400, err.Error())
				}
				spanData = temp
				tempLane1 := make(map[string]interface{})
				for _, item := range data1.([]responses.Report5) {
					if item.Data.InterventionCriteriaName == "" {
						continue
					}
					str1 := fmt.Sprintf("%s-%s-%d-%d", item.KMStart, item.KMEnd, item.LaneNo, item.Year)
					tempLane1[str1] = item
				}
				lanes1 := []responses.Report5{}
				for key := range tempLane1 {
					value, ok := tempLane1[key].(responses.Report5)
					if ok {
						lanes1 = append(lanes1, value)
					}
				}
				if len(lanes1) == 0 {
					var laneData1 responses.LaneData
					laneData1.LaneNo = 1
					laneData1.Data = responses.Report5{}
					laneDatas = append(laneDatas, laneData1)
				} else {
					var laneData1 responses.LaneData
					laneData1.LaneNo = 1
					laneData1.Data = lanes1[0]
					laneDatas = append(laneDatas, laneData1)
				}

				//////// lane 2 ////////
				data2, _, err := u.GetReportTypeData(result, roadID, directionID, int(kmStart), int(kmEnd), int(i), int(i+200), 2, year, spanData, colorCriteria)
				if err != nil {
					return "", responses.NewAppErr(400, err.Error())
				}
				tempLane2 := make(map[string]interface{})
				for _, item := range data2.([]responses.Report5) {
					if item.Data.InterventionCriteriaName == "" {
						continue
					}
					str := fmt.Sprintf("%s-%s-%d-%d", item.KMStart, item.KMEnd, item.LaneNo, item.Year)
					tempLane2[str] = item
				}
				lanes2 := []responses.Report5{}
				for key := range tempLane2 {
					value, ok := tempLane2[key].(responses.Report5)
					if ok {
						lanes2 = append(lanes2, value)
					}
				}
				if len(lanes2) == 0 {
					var laneData2 responses.LaneData
					laneData2.LaneNo = 2
					laneData2.Data = responses.Report5{}
					laneDatas = append(laneDatas, laneData2)
				} else {
					var laneData2 responses.LaneData
					laneData2.LaneNo = 2
					laneData2.Data = lanes2[0]
					laneDatas = append(laneDatas, laneData2)
				}

				//////// lane 3 ////////
				data3, _, err := u.GetReportTypeData(result, roadID, directionID, int(kmStart), int(kmEnd), int(i), int(i+200), 3, year, spanData, colorCriteria)
				if err != nil {
					return "", responses.NewAppErr(400, err.Error())
				}
				tempLane3 := make(map[string]interface{})
				for _, item := range data3.([]responses.Report5) {
					if item.Data.InterventionCriteriaName == "" {
						continue
					}
					str := fmt.Sprintf("%s-%s-%d-%d", item.KMStart, item.KMEnd, item.LaneNo, item.Year)
					tempLane3[str] = item
				}
				lanes3 := []responses.Report5{}
				for key := range tempLane3 {
					value, ok := tempLane3[key].(responses.Report5)
					if ok {
						lanes3 = append(lanes3, value)
					}
				}
				if len(lanes3) == 0 {
					var laneData3 responses.LaneData
					laneData3.LaneNo = 3
					laneData3.Data = responses.Report5{}
					laneDatas = append(laneDatas, laneData3)
				} else {
					var laneData3 responses.LaneData
					laneData3.LaneNo = 3
					laneData3.Data = lanes3[0]
					laneDatas = append(laneDatas, laneData3)
				}

				yearData = append(yearData, responses.YearData{Year: year, Data: laneDatas})

			}
			kms = append(kms, responses.KMStr{Data: yearData, KmStart: helpers.FormatKM(int64(i)), KmEnd: helpers.FormatKM(int64(i - 200))})
			i -= 200
		}

	}
	// return test, nil
	return kms, nil
}

func CalKm(roadID, directionID, kmStart, kmEnd int, resultData []models.ModelResult) ([]responses.Report5, map[string]int) {
	kms := []int{}
	kms2 := make(map[string]int)
	// helpers.PrintlnJson(kmStart, kmEnd)
	if directionID == 1 {
		for i := kmStart; i < kmEnd; {
			kms2[fmt.Sprintf("%d - %d", int(float64(i)), int(float64(i))+200)] = int(float64(i))
			kms = append(kms, int(float64(i)))
			i += 200
			// helpers.PrintlnJson(i, kmEnd)

		}
	} else {
		for i := kmStart; i > kmEnd; {
			kms2[fmt.Sprintf("%d - %d", int(float64(i)), int(float64(i))+200)] = int(float64(i))
			kms = append(kms, int(float64(i)))
			i -= 200
		}
	}

	// helpers.PrintlnJson(kms)
	dataRes := make(map[int]interface{})
	// helpers.PrintlnJson("====", kms)
	var kms3 []responses.KMReport
	// isHave := make(map[string]string)
	for _, data := range resultData {

		kmStart := data.KmStart
		kmEnd := data.KmEnd

		for _, km := range kms {
			start := 0
			end := 0
			if directionID == 1 {
				start = km
				end = km + 200
			} else {
				start = km
				end = km - 200
			}

			if dataRes[start] != nil {
				// continue
			}
			// fmt.Println("directionID", directionID)
			if directionID == 1 {

				// if kmStart < end+1 && kmEnd >= start+1 {
				//1570 < 1000&& 2000 >= 1200
				//1570 < 1600 && 2000 >= 1400
				if kmStart < end+1 && kmEnd > start+1 {
					var km responses.KMReport
					km.KmStart = kmStart
					km.KmEnd = kmEnd
					// helpers.PrintlnJson(kmStart, kmEnd)
					// helpers.PrintlnJson("lane", lane)
					km.LaneNo = data.Data.PrepareDataBefore.RoadGeom.LaneNo
					km.Year = data.AnalystYear
					km.InterventionCriteriaName = data.Data.IcResult.Name
					km.Method = data.Data.IcResult.Method
					km.KM = start
					dataRes[start] = km
					kms3 = append(kms3, km)

				} else {
					var km responses.KMReport
					km.KM = start
					km.Year = data.AnalystYear
					dataRes[start] = km
				}
			} else {
				// fmt.Println(kmStart, end-1, kmEnd, start-1)
				if kmStart > end-1 && kmEnd < start-1 {
					var km responses.KMReport
					km.KmStart = kmStart
					km.KmEnd = kmEnd
					// helpers.PrintlnJson("lane", lane)
					km.LaneNo = data.Data.PrepareDataBefore.RoadGeom.LaneNo
					km.Year = data.AnalystYear
					km.InterventionCriteriaName = data.Data.IcResult.Name
					km.Method = data.Data.IcResult.Method
					km.KM = start
					dataRes[start] = km
					kms3 = append(kms3, km)

				} else {
					var km responses.KMReport
					km.KM = start
					km.Year = data.AnalystYear
					dataRes[start] = km
				}
			}
		}
	}
	// helpers.PrintlnJson("kms3", kms3)
	// return kms3, nil

	var report5s []responses.Report5
	for index, item := range kms3 {
		KMStart := 0
		KMEnd := 0
		if directionID == 1 {
			KMStart = item.KM
			KMEnd = item.KM + 200
			if KMStart < item.KmStart {
				// continue
			}
		} else {
			KMStart = item.KM
			KMEnd = item.KM - 200
		}
		// if start == item.KM && end == item.KM+200 {
		var report5 responses.Report5
		report5.KM = index
		// if item != nil {
		// report5.LaneNo = item
		report5.KM = item.KM
		report5.KMStart = helpers.FormatKM(int64(KMStart))
		report5.KMEnd = helpers.FormatKM(int64(KMEnd))
		report5.Data.KmStart = item.KmStart
		report5.Data.KmEnd = item.KmEnd
		report5.Data.InterventionCriteriaName = item.InterventionCriteriaName
		report5.Data.Method = item.Method
		report5.Data.Unique = fmt.Sprintf("%d%d%d%d%s", item.KmStart, item.KmEnd, item.LaneNo, item.Year, item.InterventionCriteriaName)
		report5.LaneNo = item.LaneNo
		report5.Year = item.Year

		// }
		report5s = append(report5s, report5)
		// }
	}

	sort.Slice(report5s, func(i, j int) bool {
		return report5s[i].Year < report5s[j].Year
	})

	// cntSpan := map[string]int{}
	// for _, report := range report5s {
	// 	temp := fmt.Sprintf("%d-%v-%v-%v-%v", roadID, report.Data.KmStart, report.Data.KmEnd, report.Year, report.LaneNo)
	// 	cntSpan[temp]++
	// 	// dataTemp = append(dataTemp, temp)
	// }
	// helpers.PrintlnJson(cntSpan)

	return report5s, map[string]int{}
	// return report5s, nil
	// return dataRes, nil

	// fff := dataTemp[0]
	// helpers.PrintlnJson(fff)
	// return dataTemp, dataTemp, nil
}
func (u *UseCase) GetReportTypeData(report5s []responses.Report5, roadID, directionID, kmStart, kmEnd int, start, end, laneNo, year int, dataTemp map[string]int, colorCriteria map[string]string) (interface{}, map[string]int, error) {
	var reportData []responses.Report5
	data := map[string]int{}
	for _, data := range report5s {
		dataKmStart := 0
		dataKmEnd := 0
		if directionID == 1 {
			dataKmStart = data.KM
			dataKmEnd = data.KM + 200
		} else {
			dataKmStart = data.KM
			dataKmEnd = data.KM - 200
		}

		// ff :=
		if start == dataKmStart && end == dataKmEnd && data.LaneNo == laneNo && data.Year == year {
			var report responses.Report5
			report.KM = data.KM
			report.KMStart = helpers.FormatKM(int64(data.KM))
			report.KMEnd = helpers.FormatKM(int64(report.KM + 200))
			report.LaneNo = data.LaneNo
			report.Year = data.Year
			report.Data.KmStart = data.Data.KmStart
			report.Data.KmEnd = data.Data.KmEnd
			report.Data.InterventionCriteriaName = data.Data.InterventionCriteriaName
			report.Data.Method = data.Data.Method
			report.Data.Unique = fmt.Sprintf("%d%d%d%d%s", data.Data.KmStart, data.Data.KmEnd, data.LaneNo, data.Year, data.Data.InterventionCriteriaName)
			report.Data.Color = colorCriteria[data.Data.Method]
			reportData = append(reportData, report)

			//
		}
	}

	var report2Data []responses.Report5
	// ffff := map[string]string{}

	for _, item := range reportData {
		itemKmStart := 0
		itemKmEnd := 0
		if directionID == 1 {
			itemKmStart = item.KM
			itemKmEnd = item.KM + 200
		} else {
			itemKmStart = item.KM
			itemKmEnd = item.KM - 200
		}

		if start == itemKmStart && end == itemKmEnd {
			var report2 responses.Report5
			report2.KM = item.KM
			report2.LaneNo = item.LaneNo
			report2.KMStart = item.KMStart
			report2.KMEnd = item.KMEnd
			if directionID == 1 {
				if item.Data.KmStart <= item.KM && item.Data.KmEnd > item.KM {
					report2.Data.KmStart = item.Data.KmStart
					report2.Data.KmEnd = item.Data.KmEnd
					report2.Data.InterventionCriteriaName = item.Data.InterventionCriteriaName
					report2.Data.Method = item.Data.Method
					report2.Year = item.Year
					report2.Data.Unique = fmt.Sprintf("%d%d%d%d%s", item.Data.KmStart, item.Data.KmEnd, item.LaneNo, item.Year, item.Data.InterventionCriteriaName)
					report2.Data.Color = item.Data.Color

				} else {
					if (item.Data.KmStart >= item.KM && item.KM+200 <= item.Data.KmEnd) && item.Data.KmStart != item.KM+200 {
						report2.Data.KmStart = item.Data.KmStart
						report2.Data.KmEnd = item.Data.KmEnd
						report2.Data.InterventionCriteriaName = item.Data.InterventionCriteriaName
						report2.Data.Method = item.Data.Method
						report2.Year = item.Year
						report2.Data.Unique = fmt.Sprintf("%d%d%d%d%s", item.Data.KmStart, item.Data.KmEnd, item.LaneNo, item.Year, item.Data.InterventionCriteriaName)
						report2.Data.Color = item.Data.Color
					} else {
						if item.Data.KmStart == item.KM+200 {
							if item.Data.KmStart <= item.KM && item.Data.KmEnd >= item.KM {
								report2.Data.KmStart = item.Data.KmStart
								report2.Data.KmEnd = item.Data.KmEnd
								report2.Data.InterventionCriteriaName = item.Data.InterventionCriteriaName
								report2.Data.Method = item.Data.Method
								report2.Year = item.Year
								report2.Data.Unique = fmt.Sprintf("%d%d%d%d%s", item.Data.KmStart, item.Data.KmEnd, item.LaneNo, item.Year, item.Data.InterventionCriteriaName)
								report2.Data.Color = item.Data.Color
							} else {
								// helpers.PrintlnJson("11111", item.Data.KmStart, item.Data.KmEnd)
							}
						} else {
							if item.KM <= item.Data.KmStart && item.KM+200 >= item.Data.KmEnd {
								report2.Data.KmStart = item.Data.KmStart
								report2.Data.KmEnd = item.Data.KmEnd
								report2.Data.InterventionCriteriaName = item.Data.InterventionCriteriaName
								report2.Data.Method = item.Data.Method
								report2.Year = item.Year
								report2.Data.Unique = fmt.Sprintf("%d%d%d%d%s", item.Data.KmStart, item.Data.KmEnd, item.LaneNo, item.Year, item.Data.InterventionCriteriaName)
								report2.Data.Color = item.Data.Color
							}
						}
					}
				}
			} else {
				// if item.Data.KmStart > item.KM-1 && item.Data.KmEnd < item.KM-201 {
				if item.Data.KmStart >= item.KM && item.Data.KmEnd < item.KM {
					report2.Data.KmStart = item.Data.KmStart
					report2.Data.KmEnd = item.Data.KmEnd
					report2.Data.InterventionCriteriaName = item.Data.InterventionCriteriaName
					report2.Data.Method = item.Data.Method
					report2.Year = item.Year
					report2.Data.Unique = fmt.Sprintf("%d%d%d%d%s", item.Data.KmStart, item.Data.KmEnd, item.LaneNo, item.Year, item.Data.InterventionCriteriaName)
					report2.Data.Color = item.Data.Color
				} else {
					if (item.Data.KmStart <= item.KM && item.KM-200 >= item.Data.KmEnd) && item.Data.KmStart != item.KM-200 {
						report2.Data.KmStart = item.Data.KmStart
						report2.Data.KmEnd = item.Data.KmEnd
						report2.Data.InterventionCriteriaName = item.Data.InterventionCriteriaName
						report2.Data.Method = item.Data.Method
						report2.Year = item.Year
						report2.Data.Unique = fmt.Sprintf("%d%d%d%d%s", item.Data.KmStart, item.Data.KmEnd, item.LaneNo, item.Year, item.Data.InterventionCriteriaName)
						report2.Data.Color = item.Data.Color
					} else {
						if item.Data.KmStart == item.KM-200 {
							if item.Data.KmStart >= item.KM && item.Data.KmEnd <= item.KM {
								report2.Data.KmStart = item.Data.KmStart
								report2.Data.KmEnd = item.Data.KmEnd
								report2.Data.InterventionCriteriaName = item.Data.InterventionCriteriaName
								report2.Data.Method = item.Data.Method
								report2.Year = item.Year
								report2.Data.Unique = fmt.Sprintf("%d%d%d%d%s", item.Data.KmStart, item.Data.KmEnd, item.LaneNo, item.Year, item.Data.InterventionCriteriaName)
								report2.Data.Color = item.Data.Color
							}
						} else {
							if item.KM >= item.Data.KmStart && item.KM-200 <= item.Data.KmEnd {
								report2.Data.KmStart = item.Data.KmStart
								report2.Data.KmEnd = item.Data.KmEnd
								report2.Data.InterventionCriteriaName = item.Data.InterventionCriteriaName
								report2.Data.Method = item.Data.Method
								report2.Year = item.Year
								report2.Data.Unique = fmt.Sprintf("%d%d%d%d%s", item.Data.KmStart, item.Data.KmEnd, item.LaneNo, item.Year, item.Data.InterventionCriteriaName)
								report2.Data.Color = item.Data.Color
							}
						}
					}
				}
			}

			// }

			report2Data = append(report2Data, report2)
		}
	}
	sort.Slice(report2Data, func(i, j int) bool {
		return report2Data[i].Year < report2Data[j].Year
	})
	// helpers.PrintlnJson(data)
	return report2Data, data, nil
}
