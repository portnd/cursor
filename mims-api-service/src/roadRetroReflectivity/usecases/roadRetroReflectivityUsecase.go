package usecases

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadRetroReflectivity/domains"
	"gorm.io/gorm"
)

type UseCase struct {
	Repo domains.Repository
}

// init usecase
func NewUseCase(repo domains.Repository) domains.UseCase {
	return &UseCase{
		Repo: repo,
	}
}

func (t *UseCase) CreateRoadRetroReflectivity(c *gin.Context, uid uint, req requests.RoadRetroReflectivity, files requests.RoadRetroReflectivityFiles) (responses.RoadRetroReflectivityCreate, error) {

	// start transaction

	tx := t.Repo.StartTransSection()

	CSVDir := os.Getenv("ROAD_RETRO_REFLECTIVITY_CSV_DIR")

	surveyedDate, err := time.Parse("2006-01-02 15:04:00", req.SurveyedDate)
	if err != nil {
		return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(400, err.Error())
	}
	rrs := models.RoadRetroReflectivity{
		SurveyedDate: surveyedDate,
		Year:         surveyedDate.Year(),
		LineNo:       req.LineNo,
		CreatedBy:    int(uid),
		UpdatedBy:    int(uid),
		IDParent:     0,
		Remarks:      req.Remarks,
		Status:       "A",
		Revision:     0,
	}

	var irifilePath string

	if files.CsvFilename.Filename != "" {

		irifilePath, err = helpers.SaveFile(c, files.CsvFilename, CSVDir)
		if err != nil {
			return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())
		}
		err := helpers.SetPermissionsTo0775(CSVDir)
		if err != nil {
			return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())
		}
	}

	CSVData, err := readRoadRetroReflectivityCSVFile(irifilePath)
	if err != nil {
		return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())

	}

	if len(CSVData) == 0 {
		return responses.RoadRetroReflectivityCreate{}, nil
	}
	rrs.RoadID = CSVData[0].RoadID
	rrs.KmStart = CSVData[0].KMStart
	rrs.KmEnd = CSVData[len(CSVData)-1].KMEnd

	if rrs.RoadID != req.RoadID {
		return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(400, constants.INVALID_ROAD_RETRO_REFLECTIVITY_GEOM_RANGE)
	}

	road, err := t.Repo.GetRoadByID(rrs.RoadID)
	if err != nil {
		return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())

	}

	totalLanes, err := t.Repo.GetTotalLanesByRoadID(rrs.RoadID)
	if err != nil {
		return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())
	}

	line := totalLanes + 1

	//เส้นจราจรจะเท่ากับ lane ทั้งหมดบวก 1 เสมอ
	if int64(req.LineNo) > line || req.LineNo <= 0 {
		return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(400, constants.INVALID_ROAD_RETRO_REFLECTIVITY_LINE+fmt.Sprintf(" %d", line))
	}

	// var getKmType int
	// switch road.RefDirectionId {
	// case 1:
	// 	getKmType = 2

	// case 2:
	// 	getKmType = 1

	// }

	// kmValues, err := t.Repo.GetRoadKmRange(rc.RoadId, int(fullGeom.KmStart), int(fullGeom.KmEnd), getKmType)
	// if err != nil {
	// 	return responses.RoadRetroReflectivityCreate{}, fmt.Errorf("fullGrom:", err.Error())
	// }

	refStripeTypeTemp := make(map[string]models.RefStripeType)
	refStripeTypes, err := t.Repo.GetRefStripeType()
	if err != nil {
		return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())

	}
	for _, refStripeType := range refStripeTypes {
		refStripeTypeTemp[refStripeType.Name] = refStripeType
	}

	refStripeColorTemp := make(map[string]models.RefStripeColor)
	refStripeColors, err := t.Repo.GetRefStripeColor()
	if err != nil {
		return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())

	}
	for _, refStripeColor := range refStripeColors {
		refStripeColorTemp[refStripeColor.Name] = refStripeColor
	}

	err = validateData(CSVData, refStripeTypeTemp, refStripeColorTemp, road.RoadInfo.RoadInfo, road.RoadInfo.RefDirectionId)
	if err != nil {
		return responses.RoadRetroReflectivityCreate{}, err
	}

	kmCountDivider := rrs.KmStart
	analysisData, err := RoadRetroReflectivityAnalysisData(CSVData, refStripeTypeTemp, refStripeColorTemp, kmCountDivider, road.RoadInfo.RoadInfo)
	if err != nil {
		return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())
	}

	if analysisData.RoadData.TotalM >= 0 {
		rrs.RetroMin = helpers.CalculateRoadCondition(analysisData.RoadData.RetroMinAverage, &analysisData.RoadData.TotalM)
		rrs.RetroMax = helpers.CalculateRoadCondition(analysisData.RoadData.RetroMaxAverage, &analysisData.RoadData.TotalM)
		rrs.RetroAvg = helpers.CalculateRoadCondition(analysisData.RoadData.RetroAvgAverage, &analysisData.RoadData.TotalM)

	}

	roundedStruct := helpers.RoundStructPointerFloats(&rrs, 3).(*models.RoadRetroReflectivity)
	rrs = *roundedStruct

	rrs.CreatedDate = time.Now()
	rrs.CreatedBy = int(uid)
	rrs.UpdatedDate = time.Now()
	rrsID, err := t.Repo.CreateRoadRetroReflectivity(tx, rrs)
	if err != nil {
		logs.Error(err.Error())
		t.Repo.RollBack(tx)
		return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())
	}

	filePath := fmt.Sprintf("storages/road/attachments/%d/retro-reflectivity/%d/", rrs.RoadID, rrsID)
	if err := os.MkdirAll(filePath, 0775); err != nil {
		logs.Error(err)
		t.Repo.RollBack(tx)
		return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
	}

	iriFileName := filepath.Base(irifilePath)
	// update road damage csv file path
	fullCsvName := filePath + iriFileName
	// move file csv upload
	if err := helpers.MoveFile(irifilePath, fullCsvName); err != nil {
		logs.Error(err)
		t.Repo.RollBack(tx)
		return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())
	}

	err = t.Repo.UpdateRoadRetroReflectivityFilepath(tx, rrsID, fullCsvName)
	if err != nil {
		t.Repo.RollBack(tx)
		logs.Error(err.Error())
		return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())
	}

	m100List := make(map[string]int)
	for _, Item := range analysisData.RoadRetroReflectivityRange {
		Item.RoadRetroReflectivityID = rrsID

		roundedStruct := helpers.RoundStructPointerFloats(&Item, 3).(*models.RoadRetroReflectivityRange)
		//roundedStruct := roundedStructValue.Interface().(*models.RoadConditionSurvey)
		Item = *roundedStruct

		rcServayId, err := t.Repo.CreateRoadRetroReflectivityRange(tx, Item)
		if err != nil {
			logs.Error(err.Error())
			t.Repo.RollBack(tx)
			return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())
		}
		m100 := fmt.Sprintf("%v-%v", Item.KmStart, Item.KmEnd)
		m100List[m100] = rcServayId

	}

	for m100, rcServayId := range m100List {
		// Check if mItem km range overlaps with kmList[km] km range
		m100Range := strings.Split(m100, "-")
		m100Start, _ := strconv.ParseFloat(m100Range[0], 32)
		m100End, _ := strconv.ParseFloat(m100Range[1], 32)
		subRoadMin := helpers.Min(m100Start, m100End)
		subRoadMax := helpers.Max(m100Start, m100End)
		for _, mItem := range analysisData.RoadRetroReflectivityM {

			switch road.RoadInfo.RefDirectionId {
			case 1:

				if mItem.KmEnd > float64(subRoadMin) && mItem.KmEnd <= float64(subRoadMax) {
					mItem.RoadRetroReflectivityRangeID = rcServayId
				}
			case 2:
				if mItem.KmStart > float64(subRoadMin) && mItem.KmStart <= float64(subRoadMax) {
					mItem.RoadRetroReflectivityRangeID = rcServayId

				}
			}

			if mItem.RoadRetroReflectivityRangeID == 0 {
				continue
			}

			if mItem.RoadRetroReflectivityRangeID != 0 {

				roundedStruct := helpers.RoundStructPointerFloats(&mItem, 3).(*models.RoadRetroReflectivityM)
				//roundedStructValue := helpers.RoundStructFloats(&mItem, 3)
				//roundedStruct := roundedStructValue.Interface().(*models.RoadConditionSurveyM)
				mItem = *roundedStruct

				err := t.Repo.CreateRoadRetroReflectivityM(tx, mItem)
				if err != nil {
					logs.Error(err.Error())
					t.Repo.RollBack(tx)
					return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())
				}
			}

		}

	}
	t.Repo.Commit(tx)

	var resp responses.RoadRetroReflectivityCreate
	resp.ID = rrsID
	return resp, nil
}

func readRoadRetroReflectivityCSVFile(filePath string) ([]models.RoadRetroReflectivityCSV, error) {
	isFirstRow := true
	headerMap := make(map[string]int)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	var RoadRetroReflectivityCSVs []models.RoadRetroReflectivityCSV
	count := 0
	for {
		count++

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		// Handle first row case
		if isFirstRow {
			isFirstRow = false
			// Add mapping: Column/property name --> record index
			for i, v := range record {
				headerMap[v] = i
			}
			// Skip next code
			continue
		}

		RoadRetroReflectivityCSVs = append(RoadRetroReflectivityCSVs, models.RoadRetroReflectivityCSV{
			RoadID:     helpers.StrToInt(record[headerMap["road_id"]]),
			RoadCode:   record[headerMap["road_code"]],
			Name:       record[headerMap["name"]],
			KMStart:    helpers.StrToFloatValidate(record[headerMap["km_start"]]),
			KMEnd:      helpers.StrToFloatValidate(record[headerMap["km_end"]]),
			RetroMin:   helpers.StrToFloatPointerValidate(record[headerMap["retro_min"]]),
			RetroMax:   helpers.StrToFloatPointerValidate(record[headerMap["retro_max"]]),
			RetroAvg:   helpers.StrToFloatPointerValidate(record[headerMap["retro_avg"]]),
			Color:      record[headerMap["color"]],
			StripeType: record[headerMap["stripe_type"]],
		})

	}
	return RoadRetroReflectivityCSVs, nil

}

func RoadRetroReflectivityAnalysisData(csvDatas []models.RoadRetroReflectivityCSV, refStripeType map[string]models.RefStripeType, refStripeColor map[string]models.RefStripeColor, countDivider float64, roadInfo models.RoadInfo) (responses.RoadRetroReflectivityAnalysisData, error) {

	// var kmItem models.RoadConditionSurvey
	// var 100mItem models.RoadConditionSurvey100M

	//เพิ่มการเช็ค color และ stripe type ถ้าไม่ตรงกับ ข้อมูลในฐานข้อมูลใน Return error
	//เพิ่มการเช็ค line ให้ไม่เกิน road lane +1 และไม่เป็น 0

	var currentType string
	var currentCount int
	colorData := make(map[string][]models.RoadRetroReflectivityCSV)

	for _, csvData := range csvDatas {
		if csvData.StripeType+"-"+csvData.Color != currentType {
			currentType = csvData.StripeType + "-" + csvData.Color
			currentCount++
		}

		newColor := fmt.Sprintf("%s-%s-%d", csvData.StripeType, csvData.Color, currentCount)
		colorData[newColor] = append(colorData[newColor], csvData)

	}

	var resps responses.RoadRetroReflectivityAnalysisData

	var roadData models.RoadRetroReflectivityData
	for _, dataSlice := range colorData {

		// parts := strings.Split(stripeType, "-")
		// surveyType := parts[0]

		m100CountDivider := dataSlice[0].KMStart

		for i, data := range dataSlice {

			var mItem models.RoadRetroReflectivityM
			mItem.KmStart = data.KMStart
			mItem.KmEnd = data.KMEnd
			mItem.RetroMin = data.RetroMin
			mItem.RetroMax = data.RetroMax
			mItem.RetroAvg = data.RetroAvg

			subRoadStart := math.Abs(mItem.KmStart-float64(roadInfo.KmStart)) / math.Abs(float64(roadInfo.KmStart)-float64(roadInfo.KmEnd))
			subRoadEnd := math.Abs(mItem.KmEnd-float64(roadInfo.KmStart)) / math.Abs(float64(roadInfo.KmStart)-float64(roadInfo.KmEnd))
			subRoadMin := math.Min(subRoadStart, subRoadEnd)
			subRoadMax := math.Max(subRoadStart, subRoadEnd)

			mItem.TheGeom = fmt.Sprintf("ST_LineSubstring('%s', %f, %f)", roadInfo.TheGeom, subRoadMin, subRoadMax)

			//stripeType := strings.ReplaceAll(data.StripeType, "_", " ")

			mItem.RefStripeTypeID = refStripeType[data.StripeType].ID

			mItem.RefStripeColorID = refStripeColor[data.Color].ID
			resps.RoadRetroReflectivityM = append(resps.RoadRetroReflectivityM, mItem)

			roadData.RetroMinAverage = helpers.CalculateRoadConditionAverage(roadData.RetroMinAverage, mItem.RetroMin, mItem.KmStart, mItem.KmEnd)
			roadData.RetroMaxAverage = helpers.CalculateRoadConditionAverage(roadData.RetroMaxAverage, mItem.RetroMax, mItem.KmStart, mItem.KmEnd)
			roadData.RetroAvgAverage = helpers.CalculateRoadConditionAverage(roadData.RetroAvgAverage, mItem.RetroAvg, mItem.KmStart, mItem.KmEnd)

			roadData.TotalM += math.Abs(mItem.KmStart - mItem.KmEnd)

			roadData.RetroMin, roadData.DividerCountRetroMin = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountRetroMin, roadData.RetroMin, mItem.RetroMin, mItem.KmStart, mItem.KmEnd)
			roadData.RetroMax, roadData.DividerCountRetroMax = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountRetroMax, roadData.RetroMax, mItem.RetroMax, mItem.KmStart, mItem.KmEnd)
			roadData.RetroAvg, roadData.DividerCountRetroAvg = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountRetroAvg, roadData.RetroAvg, mItem.RetroAvg, mItem.KmStart, mItem.KmEnd)

			if int(data.KMEnd)%100 == 0 || i == len(dataSlice)-1 {

				rangeItem := models.RoadRetroReflectivityRange{
					KmStart: m100CountDivider,
					KmEnd:   mItem.KmEnd,
				}

				subRoadStart := math.Abs(rangeItem.KmStart-float64(roadInfo.KmStart)) / math.Abs(float64(roadInfo.KmStart)-float64(roadInfo.KmEnd))
				subRoadEnd := math.Abs(rangeItem.KmEnd-float64(roadInfo.KmStart)) / math.Abs(float64(roadInfo.KmStart)-float64(roadInfo.KmEnd))
				subRoadMin := math.Min(subRoadStart, subRoadEnd)
				subRoadMax := math.Max(subRoadStart, subRoadEnd)
				rangeItem.TheGeom = fmt.Sprintf("ST_LineSubstring('%s', %f, %f)", roadInfo.TheGeom, subRoadMin, subRoadMax)

				if roadData.DividerCountRetroMin != nil {
					rangeItem.RetroMin = helpers.CalculateRoadCondition(roadData.RetroMin, roadData.DividerCountRetroMin)
				}

				if roadData.DividerCountRetroMax != nil {
					rangeItem.RetroMax = helpers.CalculateRoadCondition(roadData.RetroMax, roadData.DividerCountRetroMax)
				}

				if roadData.DividerCountRetroAvg != nil {
					rangeItem.RetroAvg = helpers.CalculateRoadCondition(roadData.RetroAvg, roadData.DividerCountRetroAvg)
				}

				//stripeType := strings.ReplaceAll(data.StripeType, "_", " ")
				rangeItem.RefStripeTypeID = refStripeType[data.StripeType].ID
				rangeItem.RefStripeColorID = refStripeColor[data.Color].ID

				resps.RoadRetroReflectivityRange = append(resps.RoadRetroReflectivityRange, rangeItem)

				//reset value
				*roadData.DividerCountRetroMin = 0
				*roadData.DividerCountRetroMax = 0
				*roadData.DividerCountRetroAvg = 0
				*roadData.RetroMin = 0
				*roadData.RetroMax = 0
				*roadData.RetroAvg = 0

				m100CountDivider = rangeItem.KmEnd

			}
			resps.RoadData = roadData
		}

	}

	return resps, nil
}

func validateData(csvData []models.RoadRetroReflectivityCSV, refStripeType map[string]models.RefStripeType, refStripeColor map[string]models.RefStripeColor, roadInfo models.RoadInfo, directionID int) error {

	//	roadID := csvData[0].RoadId

	kmStart := csvData[0].KMStart
	kmEnd := csvData[len(csvData)-1].KMEnd
	switch directionID {
	case 1:
		if kmStart < float64(roadInfo.KmStart) || kmEnd > float64(roadInfo.KmEnd) {

			return responses.NewAppErr(400, constants.INVALID_ROAD_RETRO_REFLECTIVITY_GEOM_RANGE)

		}
	case 2:
		if kmStart > float64(roadInfo.KmStart) || kmEnd < float64(roadInfo.KmEnd) {

			return responses.NewAppErr(400, constants.INVALID_ROAD_RETRO_REFLECTIVITY_GEOM_RANGE)

		}
	}

	var rtt models.RoadRetroReflectivityCSV
	for i := 0; i < len(csvData); i++ {
		row := i + 1

		if _, exists := refStripeType[csvData[i].StripeType]; !exists {
			return responses.NewAppErr(400, strings.Replace(constants.INVALID_ROAD_RETRO_REFLECTIVITY_STRIPE_TYPE, "_", fmt.Sprintf("%d", row), -1))
		}

		if _, exists := refStripeColor[csvData[i].Color]; !exists {
			return responses.NewAppErr(400, strings.Replace(constants.INVALID_ROAD_RETRO_REFLECTIVITY_COLOR, "_", fmt.Sprintf("%d", row), -1))
		}

		if csvData[i].RoadID == rtt.RoadID ||
			csvData[i].RoadCode == rtt.RoadCode ||
			csvData[i].Name == rtt.Name ||
			csvData[i].RetroMin == nil ||
			csvData[i].RetroMax == nil ||
			csvData[i].RetroAvg == nil {

			return responses.NewAppErr(400, strings.Replace(constants.INVALID_ROAD_CONDITION_VALUE, "_", fmt.Sprintf("%d", row), -1))
		}
		if *csvData[i].RetroMin < 0 ||
			*csvData[i].RetroMax < 0 ||
			*csvData[i].RetroAvg < 0 {

			return responses.NewAppErr(400, strings.Replace(constants.INVALID_ROAD_CONDITION_VALUE, "_", fmt.Sprintf("%d", row), -1))
		}

	}

	return nil
}

func (t *UseCase) UpdateRoadRetroReflectivity(c *gin.Context, uid uint, IDParent int, rrsImport requests.RoadRetroReflectivity, files requests.RoadRetroReflectivityFiles, csvFilenameStatus string) (interface{}, error) {
	CSVDir := os.Getenv("ROAD_RETRO_REFLECTIVITY_CSV_DIR")

	// start transaction

	// no_file
	//    csv [validate], zip[path รูปต้องเป็น path เดิม]
	// delete
	//     csv [validate], zip[ลบpathรูป, ข้อมูล รูปหน้ากล้อง road_condition  ต้องเป็นค่าว่างทั้งหมด]
	// not_edit
	//    csv[ไม่มีการแก้ไขข้อมูล],zip[ไม่มีการแก้ไขข้อมูล]
	// upload
	//    csv[update ข้อมูล road_condition],zip[อัพเดท path img และ ข้อมูลรูปหน้ากล้อง]

	surveyedDate, err := time.Parse("2006-01-02 15:04:00", rrsImport.SurveyedDate)
	if err != nil {
		return responses.RoadConditionUpdate{}, responses.NewAppErr(400, err.Error())
	}

	roadReflective, err := t.Repo.GetAllRoadRetroReflectivityByIdParent(IDParent)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return 0, responses.NewAppErr(500, err.Error())

	}

	if rrsImport.LineNo != roadReflective.LineNo {
		return responses.RoadConditionUpdate{}, responses.NewAppErr(400, constants.INVALID_ROAD_RETRO_REFLECTIVITY_GEOM_RANGE)
	}

	totalLanes, err := t.Repo.GetTotalLanesByRoadID(rrsImport.RoadID)
	if err != nil {
		return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
	}

	//เส้นจราจรจะเท่ากับ lane ทั้งหมดบวก 1 เสมอ
	if int64(rrsImport.LineNo) > totalLanes+1 || rrsImport.LineNo <= 0 {
		return responses.RoadConditionUpdate{}, responses.NewAppErr(400, constants.INVALID_ROAD_RETRO_REFLECTIVITY_LINE)
	}

	latestRevision, err := t.Repo.GetLastRevisionByIdParent(IDParent)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return 0, responses.NewAppErr(500, err.Error())

	}
	rrs := models.RoadRetroReflectivity{
		KmStart:       roadReflective.KmStart,
		KmEnd:         roadReflective.KmEnd,
		RoadID:        roadReflective.RoadID,
		SurveyedDate:  surveyedDate,
		Year:          surveyedDate.Year(),
		LineNo:        roadReflective.LineNo,
		CreatedDate:   roadReflective.CreatedDate,
		UpdatedDate:   time.Now(),
		UpdatedBy:     int(uid),
		CreatedBy:     roadReflective.CreatedBy,
		IDParent:      IDParent,
		Remarks:       rrsImport.Remarks,
		Status:        "A",
		Revision:      latestRevision.Revision + 1,
		RetroMin:      roadReflective.RetroMin,
		RetroMax:      roadReflective.RetroMax,
		RetroAvg:      roadReflective.RetroAvg,
		InputFilePath: roadReflective.InputFilePath,
	}

	var resp responses.RoadConditionUpdate
	var irifilePath string

	tx := t.Repo.StartTransSection()
	if tx.Error != nil {
		logs.Error(err)
		return nil, responses.NewAppErr(500, "Failed to begin transaction")
	}

	switch csvFilenameStatus {
	case "upload":

		irifilePath, err = helpers.SaveFile(c, files.CsvFilename, CSVDir)
		if err != nil {
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
		}
		err := helpers.SetPermissionsTo0775(CSVDir)
		if err != nil {
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
		}

		if roadReflective.Status == "A" {

			err = t.Repo.UpdateStatusIByID(tx, roadReflective.ID, uid)
			if err != nil {
				logs.Error(err.Error())
				t.Repo.RollBack(tx)
				return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
			}
		}
	default:

		// hariphan แก้ เวลาไม่มีไฟล์ csv, zip update ให้ update เฉพาะ surveyed_date, และremarks
		data := models.RoadRetroReflectivity{
			ID:           roadReflective.ID,
			SurveyedDate: surveyedDate,
			Year:         surveyedDate.Year(),
			UpdatedDate:  time.Now(),
			UpdatedBy:    int(uid),
			Remarks:      rrsImport.Remarks,
		}
		_, err = t.Repo.UpdateRoadRetroReflectivityNoIriFile(tx, data)
		if err != nil {
			logs.Error(err.Error())
			t.Repo.RollBack(tx)
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
		}

		resp.ID = roadReflective.ID
		resp.IDParent = IDParent

	}

	if csvFilenameStatus == "upload" {
		CSVData, err := readRoadRetroReflectivityCSVFile(irifilePath)
		if err != nil {
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())

		}
		if len(CSVData) == 0 {
			return responses.RoadConditionUpdate{}, nil
		}
		rrs.RoadID = CSVData[0].RoadID
		rrs.KmStart = CSVData[0].KMStart
		rrs.KmEnd = CSVData[len(CSVData)-1].KMEnd

		if rrs.RoadID != rrsImport.RoadID {
			return responses.RoadConditionUpdate{}, responses.NewAppErr(400, constants.INVALID_ROAD_CONDITION_GEOM_RANGE)
		}

		road, err := t.Repo.GetRoadByID(rrs.RoadID)
		if err != nil {
			return responses.RoadConditionUpdate{}, err

		}

		refStripeTypeTemp := make(map[string]models.RefStripeType)
		refStripeTypes, err := t.Repo.GetRefStripeType()
		if err != nil {
			return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())

		}
		for _, refStripeType := range refStripeTypes {
			refStripeTypeTemp[refStripeType.Name] = refStripeType
		}

		refStripeColorTemp := make(map[string]models.RefStripeColor)
		refStripeColors, err := t.Repo.GetRefStripeColor()
		if err != nil {
			return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())

		}
		for _, refStripeColor := range refStripeColors {
			refStripeColorTemp[refStripeColor.Name] = refStripeColor
		}

		err = validateData(CSVData, refStripeTypeTemp, refStripeColorTemp, road.RoadInfo.RoadInfo, road.RoadInfo.RefDirectionId)
		if err != nil {
			return responses.RoadConditionUpdate{}, err
		}

		kmCountDivider := rrs.KmStart
		analysisData, err := RoadRetroReflectivityAnalysisData(CSVData, refStripeTypeTemp, refStripeColorTemp, kmCountDivider, road.RoadInfo.RoadInfo)
		if err != nil {
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
		}

		if analysisData.RoadData.TotalM >= 0 {
			rrs.RetroMin = helpers.CalculateRoadCondition(analysisData.RoadData.RetroMinAverage, &analysisData.RoadData.TotalM)
			rrs.RetroMax = helpers.CalculateRoadCondition(analysisData.RoadData.RetroMaxAverage, &analysisData.RoadData.TotalM)
			rrs.RetroAvg = helpers.CalculateRoadCondition(analysisData.RoadData.RetroAvgAverage, &analysisData.RoadData.TotalM)

		}

		roundedStruct := helpers.RoundStructPointerFloats(&rrs, 3).(*models.RoadRetroReflectivity)
		rrs = *roundedStruct

		rrs.CreatedBy = roadReflective.CreatedBy
		rrs.CreatedDate = roadReflective.CreatedDate
		rrs.UpdatedBy = int(uid)
		rrs.UpdatedDate = time.Now()
		rrs.IDParent = IDParent

		rrsResp, err := t.Repo.UpdateRoadRetroReflectivity(tx, rrs)
		if err != nil {
			logs.Error(err.Error())
			t.Repo.RollBack(tx)
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())

		}

		filePath := fmt.Sprintf("storages/road/attachments/%d/retro-reflectivity/%d/", rrs.RoadID, rrsResp.ID)
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			logs.Error(err)
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, constants.FAILED_TO_SAVE_FILE)
		}

		iriFileName := filepath.Base(irifilePath)
		// update road damage csv file path
		fullCsvName := filePath + iriFileName

		// move file csv upload
		if err := helpers.MoveFile(irifilePath, fullCsvName); err != nil {
			logs.Error(err)
			return responses.RoadConditionUpdate{}, responses.NewAppErr(400, err.Error())
		}

		err = t.Repo.UpdateRoadRetroReflectivityFilepath(tx, rrsResp.ID, fullCsvName)
		if err != nil {
			logs.Error(err.Error())
			t.Repo.RollBack(tx)
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
		}
		m100List := make(map[string]int)
		for _, Item := range analysisData.RoadRetroReflectivityRange {
			Item.RoadRetroReflectivityID = rrsResp.ID

			roundedStruct := helpers.RoundStructPointerFloats(&Item, 3).(*models.RoadRetroReflectivityRange)
			//roundedStruct := roundedStructValue.Interface().(*models.RoadConditionSurvey)
			Item = *roundedStruct

			rcServayId, err := t.Repo.CreateRoadRetroReflectivityRange(tx, Item)
			if err != nil {
				logs.Error(err.Error())
				t.Repo.RollBack(tx)
				return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())
			}
			m100 := fmt.Sprintf("%v-%v", Item.KmStart, Item.KmEnd)
			m100List[m100] = rcServayId

		}

		for m100, rcServayId := range m100List {
			// Check if mItem km range overlaps with kmList[km] km range
			m100Range := strings.Split(m100, "-")
			m100Start, _ := strconv.ParseFloat(m100Range[0], 32)
			m100End, _ := strconv.ParseFloat(m100Range[1], 32)
			subRoadMin := helpers.Min(m100Start, m100End)
			subRoadMax := helpers.Max(m100Start, m100End)
			for _, mItem := range analysisData.RoadRetroReflectivityM {

				switch road.RoadInfo.RefDirectionId {
				case 1:

					if mItem.KmEnd > float64(subRoadMin) && mItem.KmEnd <= float64(subRoadMax) {
						mItem.RoadRetroReflectivityRangeID = rcServayId
					}
				case 2:
					if mItem.KmStart > float64(subRoadMin) && mItem.KmStart <= float64(subRoadMax) {
						mItem.RoadRetroReflectivityRangeID = rcServayId

					}
				}

				if mItem.RoadRetroReflectivityRangeID == 0 {
					continue
				}

				if mItem.RoadRetroReflectivityRangeID != 0 {

					roundedStruct := helpers.RoundStructPointerFloats(&mItem, 3).(*models.RoadRetroReflectivityM)

					mItem = *roundedStruct

					err := t.Repo.CreateRoadRetroReflectivityM(tx, mItem)
					if err != nil {
						logs.Error(err.Error())
						t.Repo.RollBack(tx)
						return responses.RoadRetroReflectivityCreate{}, responses.NewAppErr(500, err.Error())
					}
				}

			}

		}
		t.Repo.Commit(tx)
		resp.ID = rrsResp.ID
		resp.IDParent = IDParent
	}
	return resp, nil
}

func (t *UseCase) DeleteRoadRetroReflectivity(c *gin.Context, uid uint, idParent int) (bool, error) {

	RoadRetroReflectivity, err := t.Repo.GetAllRoadRetroReflectivityByIdParent(idParent)
	if err != nil {
		return false, err
	}

	switch RoadRetroReflectivity.Status {
	// case "T":
	// 	{
	// 		roadCondition.Status = "D"
	// 		roadCondition.UpdatedDate = time.Now()
	// 		roadCondition.UpdatedBy = int(uid)

	// 		err = t.roadConditionRepo.DeleteRoadCondition(roadCondition)
	// 		if err != nil {
	// 			return false, err
	// 		}
	// 	}
	case "A":
		{

			if RoadRetroReflectivity.Revision > 0 {

				before, err := t.Repo.GetRoadRefactiveStripBeforeLastRevitionByIdParent(idParent, RoadRetroReflectivity.Revision-1)
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						return true, nil
					}
					return false, err
				}
				if before.Status == "I" {

					RoadRetroReflectivity.Status = "D"
					RoadRetroReflectivity.UpdatedDate = time.Now()
					RoadRetroReflectivity.UpdatedBy = int(uid)
					err = t.Repo.DeleteRoadRetroReflectivity(RoadRetroReflectivity.RoadRetroReflectivity)
					if err != nil {
						return false, err
					}

					before.Status = "A"
					before.UpdatedDate = time.Now()
					before.UpdatedBy = int(uid)
					err = t.Repo.DeleteRoadRetroReflectivity(before)
					if err != nil {
						return false, err
					}

				}

			} else {
				RoadRetroReflectivity.Status = "D"
				RoadRetroReflectivity.UpdatedDate = time.Now()
				RoadRetroReflectivity.UpdatedBy = int(uid)
				err = t.Repo.DeleteRoadRetroReflectivity(RoadRetroReflectivity.RoadRetroReflectivity)
				if err != nil {
					return false, err
				}
			}

		}
	}

	return true, nil
}

func (t *UseCase) GetRoadRetroReflectivityTemplate(uid uint, roadID int) (interface{}, error) {

	// roadConditions, err := t.roadConditionRepo.GetRoadConditionByRoadID(roadID)
	// if err != nil {
	// 	return "", err
	// }

	header := []string{"road_id", "road_code", "name", "km_start", "km_end", "stripe_type", "color", "retro_max", "retro_min", "retro_avg"}
	filePath := "storages/template_retro_reflectivity/"

	fileName := ""
	//fileLocations := []string{}

	road, err := t.Repo.GetRoadByID(roadID)
	if err != nil {
		return 0, fmt.Errorf("fullGrom:", err.Error())

	}

	csvData := [][]string{}

	start := fmt.Sprintf("%d", road.Id) + "," + road.RoadSection.RoadGroup.Number + road.RoadSection.Number + "," + road.RoadSection.NameOriginTH + " - " + road.RoadSection.NameDestinationTH + "," + fmt.Sprintf("%d", int(road.RoadInfo.KmStart)) + "," + fmt.Sprintf("%d", int(road.RoadInfo.KmEnd)) + "," + "" + "," + "" + "," + fmt.Sprintf("%d", 0) + "," + fmt.Sprintf("%d", 0) + "," + fmt.Sprintf("%d", 0)
	csvData = append(csvData, strings.Split(start, ","))

	//fileName = "template_reflective_" + road.RoadSection.NameOriginEn + "-" + road.RoadSection.NameDestinationEn + ".csv"
	fileName = "template_reflectivity_" + road.RoadSection.RoadGroup.Number + road.RoadSection.Number + "_" + "road_id_" + fmt.Sprintf("%d", road.Id) + "_" + time.Now().Format("2006-01-02_150405") + ".csv"
	err = csvGenerateFile(header, csvData, filePath, fileName)
	if err != nil {
		return "", responses.NewAppErr(http.StatusInternalServerError, err.Error())
	}
	fileLocation := filePath + fileName

	err = os.Chmod(fileLocation, 0775)
	if err != nil {
		return "", responses.NewAppErr(http.StatusInternalServerError, err.Error())
	}

	//fileLocations = append(fileLocations, fileLocation)

	STORAGE_IP := os.Getenv("STORAGE_IP") + "/"
	// zipName := filePath + "template_reflective_strip_" + time.Now().Format("2006-01-02_150405") + ".zip"
	// err = helpers.Zip(zipName, fileLocations)
	// if err != nil {
	// 	return "", responses.NewAppErr(http.StatusBadRequest, err.Error())
	// }
	// err = os.Chmod(zipName, 0775)
	// if err != nil {
	// 	return "", responses.NewAppErr(http.StatusBadRequest, err.Error())
	// }

	//for _, fileLocation := range fileLocations {
	// err = os.Remove(fileLocation)
	// if err != nil {
	// 	return nil, err
	// }
	//}

	return STORAGE_IP + fileLocation, nil
}

func csvGenerateFile(header []string, data [][]string, filePath, fileName string) error {
	// Create directory
	if err := os.MkdirAll(filePath, 0775); err != nil {
		return err
	}

	// Create a new CSV file
	file, err := os.Create(filePath + fileName)
	if err != nil {
		return errors.New("cannot create file")
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the CSV header
	err = writer.Write(header)
	if err != nil {
		return errors.New("cannot write header")
	}

	for _, row := range data {

		err = writer.Write(row)
		if err != nil {
			return errors.New("cannot write row")
		}

	}
	return nil
}

func (t *UseCase) GetRoadRetroReflectivityList(roadID int) (interface{}, error) {

	result, err := t.Repo.GetRoadRetroReflectivityList(roadID)
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

func (t *UseCase) GetRoadRetroReflectivity(c *gin.Context, uid uint, idParent int) (responses.RoadRetroReflectivity, error) {

	rrs, err := t.Repo.GetRoadRetroReflectivityByIdParent(idParent)
	if err != nil {
		return responses.RoadRetroReflectivity{}, responses.NewAppErr(http.StatusInternalServerError, err.Error())
	}

	road, err := t.Repo.GetRoadByID(rrs.RoadID)
	if err != nil {
		return responses.RoadRetroReflectivity{}, responses.NewAppErr(http.StatusInternalServerError, err.Error())
	}

	STORAGE_IP := os.Getenv("STORAGE_IP") + "/"

	csvUrl := ""
	if rrs.InputFilePath != "" {
		csvUrl = STORAGE_IP + rrs.InputFilePath
	}

	resp := responses.RoadRetroReflectivity{
		ID:           rrs.ID,
		IDParent:     idParent,
		LineNo:       rrs.LineNo,
		SurveyedDate: rrs.SurveyedDate,
		Remarks:      rrs.Remarks,
		CsvFile:      csvUrl,
		Direction: models.RefDirection{
			ID:   road.RoadInfo.RefDirection.ID,
			Name: road.RoadInfo.RefDirection.Name,
		},
	}

	return resp, nil
}

func (t *UseCase) GetRoadRetroReflectivityLineList(roadID int) (interface{}, error) {
	roadLine, err := t.Repo.GetLineListByRoadID(roadID)
	if err != nil {
		return []responses.RoadLaneList{}, responses.NewAppErr(500, err.Error())
	}

	var roadLineLists []responses.RoadRetroReflectivityLineList

	copier.Copy(&roadLineLists, &roadLine)

	if len(roadLineLists) == 0 {
		return []responses.RoadLaneList{}, nil
	}
	return roadLineLists, nil
}

func (t *UseCase) GetRoadRetroReflectivityDetails(uid uint, rangeType string, refStripeTypeIDs string, idParent int) (interface{}, error) {

	var resp responses.RetroReflectivityDetails
	var stripeTypes []string
	if refStripeTypeIDs != "" {
		stripeTypes = strings.Split(refStripeTypeIDs, ",")
		fmt.Println(stripeTypes)
	}

	retroReflectivitys, err := t.Repo.GetRoadRetroReflectivityDetailsByIdParent(idParent, stripeTypes)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil // or return an appropriate empty response
		}

		return nil, responses.NewAppErr(500, err.Error())
	}

	road, err := t.Repo.GetRoadByID(retroReflectivitys.RoadID)
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
					RefStripeTypeID:  retroReflectivityM.RefStripeTypeID,
					RefStripeType:    retroReflectivityM.RefStripeType,
					RefStripeColorID: retroReflectivityM.RefStripeColorID,
					RefStripeColor:   retroReflectivityM.RefStripeColor,
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
	resp.Direction = road.RoadInfo.RefDirection

	userInfo, err := t.Repo.GetUserDepartmentById(retroReflectivitys.UpdatedBy)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
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

	if len(resp.Datas) == 0 {
		return responses.NoData{}, nil
	}

	return resp, nil
}

func (t *UseCase) GetRoadRetroReflectivityCompareLine(roadId int, rcCompare requests.RoadRetroReflectivityCompare) (interface{}, error) {

	years := strings.Split(rcCompare.Years, ",")
	lines := strings.Split(rcCompare.Lines, ",")

	yearInts := []int{}
	for _, year := range years {
		yearInt, err := strconv.Atoi(year)
		if err == nil {
			yearInts = append(yearInts, yearInt)
		}
	}

	lineInts := []int{}
	for _, lane := range lines {
		lineInt, err := strconv.Atoi(lane)
		if err == nil {
			lineInts = append(lineInts, lineInt)
		}
	}

	road, err := t.Repo.GetRoadByID(roadId)
	if err != nil {
		return nil, err
	}
	if !road.IsActive {
		return []responses.RoadConditionLane{}, nil
	}

	roadRetroReflectivitys, err := t.Repo.GetRoadRetroReflectivityCompare(roadId,
		models.RoadRetroReflectivityCompareLine{
			Years: yearInts,
			Lines: lineInts,
		})
	if err != nil {
		return nil, err
	}

	processedYearLanes := make(map[string]bool)
	datas := make(map[int][]responses.RoadRetroReflectivityLineItem)

	for _, rrt := range roadRetroReflectivitys {

		yearLinesKey := fmt.Sprintf("%d-%d", rrt.Year, rrt.LineNo)
		if processedYearLanes[yearLinesKey] {
			continue
		}

		for _, rrtRange := range rrt.RoadRetroReflectivityRanges {

			for _, rrtM := range rrtRange.RoadRetroReflectivityMs {
				var value *float64

				value = rrtM.RetroAvg

				datas[rrt.Year] = append(datas[rrt.Year], responses.RoadRetroReflectivityLineItem{
					LineNo:  rrt.LineNo,
					KmStart: int(rrtM.KmStart),
					KmEnd:   int(rrtM.KmEnd),
					Value:   value,
				})

			}

		}

		processedYearLanes[yearLinesKey] = true

	}

	var resps []responses.RoadRetroReflectivityLine
	for year, items := range datas {
		roadRetroReflectivityLine := responses.RoadRetroReflectivityLine{
			Year:  year,
			Items: items,
		}

		// Sort items by LaneNo and KmStart in ascending order
		sort.Slice(roadRetroReflectivityLine.Items, func(i, j int) bool {
			// First, compare by LaneNo
			if roadRetroReflectivityLine.Items[i].LineNo != roadRetroReflectivityLine.Items[j].LineNo {
				return roadRetroReflectivityLine.Items[i].LineNo < roadRetroReflectivityLine.Items[j].LineNo
			}

			// If LaneNo is the same, compare by KmStart
			if road.RoadInfo.RefDirectionId == 1 {
				return roadRetroReflectivityLine.Items[i].KmStart < roadRetroReflectivityLine.Items[j].KmStart
			} else {
				return roadRetroReflectivityLine.Items[i].KmStart > roadRetroReflectivityLine.Items[j].KmStart
			}

		})

		resps = append(resps, roadRetroReflectivityLine)
	}
	if len(resps) == 0 {
		return []responses.RoadRetroReflectivityLine{}, nil
	}
	return resps, nil
}

func (t *UseCase) GetRoadRetroReflectivityCompareYear(roadId int, rrtCompare requests.RoadRetroReflectivityCompare) (interface{}, error) {

	years := strings.Split(rrtCompare.Years, ",")
	lines := strings.Split(rrtCompare.Lines, ",")

	yearInts := []int{}
	for _, year := range years {
		yearInt, err := strconv.Atoi(year)
		if err == nil {
			yearInts = append(yearInts, yearInt)
		}
	}

	lineInts := []int{}
	for _, line := range lines {
		lineInt, err := strconv.Atoi(line)
		if err == nil {
			lineInts = append(lineInts, lineInt)
		}
	}

	road, err := t.Repo.GetRoadByID(roadId)
	if err != nil {
		return nil, err
	}

	roadInfo, err := t.Repo.GetRoadByID(roadId)
	if err != nil {
		return responses.RoadCondition{}, err
	}

	if !road.IsActive {
		return []responses.RoadConditionYear{}, nil
	}

	roadRetroReflectivitys, err := t.Repo.GetRoadRetroReflectivityCompare(roadId,
		models.RoadRetroReflectivityCompareLine{
			Years: yearInts,
			Lines: lineInts,
		})
	if err != nil {
		return []responses.RoadConditionYear{}, err
	}

	var value float64

	processedLineYears := make(map[string]bool)

	var YearItem responses.RoadRetroReflectivityYearItem
	datas := make(map[int][]responses.RoadRetroReflectivityYearItem)

	for _, rrt := range roadRetroReflectivitys {
		lineYearKey := fmt.Sprintf("%d-%d", rrt.LineNo, rrt.Year)

		if processedLineYears[lineYearKey] {
			continue
		}

		for _, rrtRange := range rrt.RoadRetroReflectivityRanges {

			for _, rrtM := range rrtRange.RoadRetroReflectivityMs {

				value = *rrtM.RetroAvg

				YearItem = responses.RoadRetroReflectivityYearItem{
					Line:    rrt.LineNo,
					Year:    rrt.Year,
					KmStart: int(rrtM.KmStart),
					KmEnd:   int(rrtM.KmEnd),
					Value:   value,
				}

				datas[rrt.LineNo] = append(datas[rrt.LineNo], YearItem)
			}

		}

		processedLineYears[lineYearKey] = true
	}

	var resps []responses.RoadRetroReflectivityYear
	for line, items := range datas {
		year := responses.RoadRetroReflectivityYear{
			Line:  line,
			Items: items,
		}
		// Sort items by LaneNo and KmStart in ascending order
		sort.Slice(year.Items, func(i, j int) bool {
			// First, compare by LaneNo

			if year.Items[i].Year != year.Items[j].Year {
				return year.Items[i].Year > year.Items[j].Year
			}

			// If LaneNo is the same, compare by KmStart
			if roadInfo.RoadInfo.RefDirectionId == 1 {
				return year.Items[i].KmStart < year.Items[j].KmStart
			} else {
				return year.Items[i].KmStart > year.Items[j].KmStart
			}

		})

		resps = append(resps, year)
	}

	sort.SliceStable(resps, func(i, j int) bool {
		return resps[i].Line < resps[j].Line
	})

	if len(resps) == 0 {
		return []responses.RoadConditionYear{}, nil
	}
	return resps, nil

}

func (t *UseCase) GetRoadRetroReflectivityCompareAverage(roadId int, line int) (interface{}, error) {

	road, err := t.Repo.GetRoadByID(roadId)
	if !road.IsActive {
		return responses.NoData{}, nil
	}
	if err != nil {
		return nil, err
	}

	roadConditions, err := t.Repo.GetRoadRetroReflectivityAverage(roadId, line)
	if err != nil {
		return nil, err
	}

	dataMap := make(map[int]models.RoadRetroReflectivityAverage)

	for _, rc := range roadConditions {

		data, ok := dataMap[rc.LineNo]
		if !ok {
			data = models.RoadRetroReflectivityAverage{
				Line:     rc.LineNo,
				Items:    []models.RoadReflectivityAverageItem{},
				Revision: make(map[int]int),
				IDParent: make(map[int]int),
			}
		}

		if _, ok := data.Revision[rc.Year]; !ok {
			data.Revision[rc.Year] = rc.Revision
		}
		if _, ok := data.IDParent[rc.Year]; !ok {
			data.IDParent[rc.Year] = rc.IDParent
		}

		if rc.Revision == data.Revision[rc.Year] && rc.IDParent == data.IDParent[rc.Year] {
			data.Items = append(data.Items, models.RoadReflectivityAverageItem{
				Year:     rc.Year,
				KmStart:  int(rc.KmStart),
				KmEnd:    int(rc.KmEnd),
				RetroAvg: rc.RetroAvg,
			})
		}

		dataMap[rc.LineNo] = data
	}
	dataList := responses.RoadRetroReflectivityAverage{}
	for _, data := range dataMap {
		cleanData := responses.RoadRetroReflectivityAverage{
			Lane:  data.Line,
			Items: data.Items,
		}
		dataList = cleanData
	}
	if dataList.Items == nil {
		return responses.NoData{}, nil
	}

	return dataList, nil
}
