package usecases

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/constants"
	helpers "gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	requests "gitlab.com/mims-api-service/requests"
	responses "gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadCondition/domains"
	"gorm.io/gorm"
)

type roadConditionUseCase struct {
	roadConditionRepo domains.RoadConditionRepository
}

// init usecase
func NewRoadConditionUseCase(repo domains.RoadConditionRepository) domains.RoadConditionUseCase {
	return &roadConditionUseCase{
		roadConditionRepo: repo,
	}
}

// =========================================================

func (t *roadConditionUseCase) GetMenu(userId uint) ([]models.AccessControl, error) {
	roles, err := t.roadConditionRepo.GetRole(userId)
	if err != nil {
		return nil, err
	}

	var roleIds []int
	for _, item := range roles {
		roleIds = append(roleIds, item.RoleID)
	}

	data, err := t.roadConditionRepo.GetAccessControl(roleIds)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type Data struct {
	GroupId      int    `json:"group_id"`
	GroupName    string `json:"group_name"`
	AssetId      int    `json:"asset_id"`
	AssetName    string `json:"asset_name"`
	GeomType     int    `json:"geom_type"`
	IsInRoad     bool   `json:"is_in_road"`
	IconFilepath string `json:"icon_filepath"`
}

type Groups struct {
	Id    int          `json:"id"`
	Name  string       `json:"name"`
	Items []GroupItems `json:"items"`
}

type GroupItems struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	GeomType     int    `json:"geom_type"`
	IconFilepath string `json:"icon_filepath"`
}

type ItemAsset struct {
	ID    int     `json:"id"`
	Items []items `json:"items"`
}

type items struct {
	ID            int `json:"id"`
	CountTemp     int `json:"count_temp"`
	CountWaiting  int `json:"count_waiting"`
	CountRejected int `json:"count_rejected"`
}

func (t *roadConditionUseCase) GetRoadConditionList(accessKeys []string, roadId int) (interface{}, error) {

	rcStatus := ""
	// if helpers.HasPermission([]string{"road_condition_manage_data"}, accessKeys) {
	rcStatus = "rc.status = 'A'"

	// }

	result, err := t.roadConditionRepo.GetRoadConditionList(rcStatus, roadId)
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

func (t *roadConditionUseCase) GetRoadConditionDetails(uid uint, conditionRangeType string, idParent int) (interface{}, error) {
	if conditionRangeType == "" {
		conditionRangeType = "1"
	}
	var resp responses.RoadConditionDetails

	surveys, err := t.roadConditionRepo.GetAllRoadConditionByIdParent(idParent)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil // or return an appropriate empty response
		}

		return nil, responses.NewAppErr(500, err.Error())
	}

	road, err := t.roadConditionRepo.GetRoadByID(surveys.RoadId)
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

			conditionValue := extractConditionValue(&survey, conditionType)
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

					detailValue := extractConditionValue(&survey100M, conditionType)

					if geom != "LINESTRING EMPTY" {
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

						// fmt.Println("geom : ", surveyM.TheGeom)

						detailValue := extractConditionValue(&surveyM, conditionType)

						img := ""
						if surveyM.ImgFilepath != "" {
							img = os.Getenv("STORAGE_IP") + "/" + surveyM.ImgFilepath
						}
						if geom != "LINESTRING EMPTY" {
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

	userInfo, err := t.roadConditionRepo.GetUserDepartmentById(surveys.UpdatedBy)
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

	return resp, nil
}

func extractConditionValue(survey interface{}, conditionType string) *float64 {
	// Example implementation
	switch s := survey.(type) {
	case *models.RoadConditionSurveyPreload:
		switch conditionType {
		case "IRI":
			return s.IRI
		case "MPD":
			return s.MPD
		case "RUT":
			return s.RUT
		case "IFI":
			return s.IFI
		}
	case *models.RoadConditionSurvey100MPreload:
		// Assuming RoadConditionSurvey100M has similar fields; adjust as necessary
		switch conditionType {
		case "IRI":
			return s.IRI
		case "MPD":
			return s.MPD
		case "RUT":
			return s.RUT
		case "IFI":
			return s.IFI
		}

	case *models.RoadConditionSurveyM:
		// Assuming RoadConditionSurvey100M has similar fields; adjust as necessary
		switch conditionType {
		case "IRI":
			return s.IRI
		case "MPD":
			return s.MPD
		case "RUT":
			return s.RUT
		case "IFI":
			return s.IFI
		}
	}
	return nil
}

// func (t *roadConditionUseCase) GetRoadConditionDetails(uid uint, conditionType string, conditionRangType string, idParent int) (interface{}, error) {

// 	var resp responses.RoadConditionDetails
// 	var rcStatus models.SqlCondition

// 	rcStatus.Where = "rc.status = 'A'"

// 	var sqlCondition models.SqlCondition
// 	conditionType = strings.ToLower(conditionType)
// 	switch conditionType {
// 	case "iri":
// 		sqlCondition.Select = "rc_m.iri AS value,"
// 	case "ifi":
// 		sqlCondition.Select = "rc_m.ifi AS value,"
// 	case "mpd":
// 		sqlCondition.Select = "rc_m.mpd AS value,"
// 	case "rut":
// 		sqlCondition.Select = "rc_m.rut AS value,"
// 	default:
// 		return responses.RoadConditionDetails{}, responses.NewAppErr(http.StatusBadRequest, constants.INVALID_CONDITION_TYPE)
// 	}

// 	roadConditionDetails, err := t.roadConditionRepo.GetRoadConditionDetails(idParent, rcStatus, sqlCondition)
// 	if err != nil {
// 		logs.Error(err)
// 		if err == gorm.ErrRecordNotFound {
// 			return responses.NoData{}, nil
// 		}
// 		return responses.RoadConditionDetails{}, err
// 	}

// 	latestRevision := roadConditionDetails[0].Revision
// 	for _, rc := range roadConditionDetails {
// 		if rc.Revision > latestRevision {
// 			latestRevision = rc.Revision
// 		}
// 	}

// 	latestRoadConditions := []models.RoadConditionDetails{}
// 	for _, rc := range roadConditionDetails {
// 		if rc.Revision == latestRevision {
// 			latestRoadConditions = append(latestRoadConditions, rc)
// 		}
// 	}

// 	type Permissions struct {
// 		CanEdit    bool `json:"can_edit"`
// 		CanDelete  bool `json:"can_delete"`
// 		CanApprove bool `json:"can_approve"`
// 		CanSend    bool `json:"can_send"`
// 		CanReject  bool `json:"can_reject"`
// 	}

// 	updateby := 0

// 	roadCondition, err := t.roadConditionRepo.GetRoadConditionByIdParent(latestRoadConditions[0].IDParent)
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return responses.NoData{}, nil
// 		}
// 		return responses.RoadConditionDetails{}, err
// 	}

// 	for _, rcd := range latestRoadConditions {

// 		direction := models.RefDirection{
// 			ID:   rcd.DirectionID,
// 			Name: rcd.DirectionName,
// 		}

// 		STORAGE_IP := os.Getenv("STORAGE_IP") + "/"
// 		updateby = rcd.UpdatedBy
// 		var items []responses.RoadConditionDetailHeaderItem
// 		itemMap := make(map[string]*responses.RoadConditionDetailHeaderItem)
// 		for _, rc := range latestRoadConditions {
// 			key := fmt.Sprintf("%d-%d", rc.KmStartKM, rc.KmEndKM)
// 			if _, ok := itemMap[key]; !ok {
// 				itemMap[key] = &responses.RoadConditionDetailHeaderItem{
// 					KmStart: rc.KmStartKM,
// 					KmEnd:   rc.KmEndKM,
// 					Items:   []responses.RoadConditionDetailBodyItem{},
// 				}
// 			}

// 			var imgPath string
// 			if roadCondition.ImgFilePath != "" && rc.ImgFilepath != "" {
// 				imgPath = STORAGE_IP + rc.ImgFilepath
// 			}
// 			// For each KmStart and KmEnd group, create a SubItem
// 			subItem := responses.RoadConditionDetailBodyItem{
// 				KmStart:     rc.KmStartM,
// 				KmEnd:       rc.KmEndM,
// 				Value:       rc.Value,
// 				GeomCl:      rc.GeomCL,
// 				ImgFilepath: imgPath,
// 			}

// 			itemMap[key].Items = append(itemMap[key].Items, subItem)

// 			switch resp.Direction.ID {
// 			case 1:

// 				sort.Slice(itemMap[key].Items, func(i, j int) bool {
// 					return itemMap[key].Items[i].KmStart < itemMap[key].Items[j].KmStart
// 				})

// 			case 2:
// 				sort.Slice(itemMap[key].Items, func(i, j int) bool {
// 					return itemMap[key].Items[i].KmStart > itemMap[key].Items[j].KmStart
// 				})
// 			}
// 		}

// 		// Convert map to slice and calculate average value for each item
// 		//var items []responses.RoadConditionDetailHeaderItem
// 		for _, item := range itemMap {
// 			var total float64
// 			for _, subItem := range item.Items {
// 				if subItem.Value != nil {
// 					total += *subItem.Value
// 				}
// 			}
// 			average := total / float64(len(item.Items))
// 			helpers.RoundFloatPointer(&average, 3)
// 			item.Value = &average
// 			items = append(items, *item)
// 		}

// 		switch resp.Direction.ID {
// 		case 1:

// 			sort.Slice(items, func(i, j int) bool {
// 				return items[i].KmStart < items[j].KmStart
// 			})

// 		case 2:
// 			sort.Slice(items, func(i, j int) bool {
// 				return items[i].KmStart > items[j].KmStart
// 			})
// 		}

// 		resp = responses.RoadConditionDetails{
// 			ID:            rcd.RoadID,
// 			IDParent:      rcd.IDParent,
// 			UpdatedDate:   helpers.SetTimeToString(rcd.UpdatedDate),
// 			Status:        rcd.StatusText,
// 			StatusCode:    rcd.Status,
// 			Permissions:   permissions,
// 			Direction:     direction,
// 			RoadTypeID:    rcd.RoadTypeID,
// 			ConditionType: conditionType,
// 			RejectReason:  rcd.RejectReason,
// 			Items:         items,
// 		}

// 	}

// 	userInfo, err := t.roadConditionRepo.GetUserDepartmentById(updateby)
// 	if err != nil {
// 		if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
// 			logs.Error(err)
// 			return responses.RoadConditionDetails{}, err
// 		}
// 	}

// 	resp.UpdatedBy.UID = fmt.Sprintf("%d", userInfo.Id)
// 	resp.UpdatedBy.UserName = userInfo.Email
// 	resp.UpdatedBy.FullName = userInfo.Firstname + " " + userInfo.Lastname
// 	resp.UpdatedBy.Department.ID = userInfo.Department.ID
// 	resp.UpdatedBy.Department.Name = userInfo.Department.Name
// 	resp.UpdatedBy.ProfilePicture = userInfo.ProfileImgPath

// 	return resp, nil
// }

func (t *roadConditionUseCase) CreateRoadCondition(c *gin.Context, uid uint, rcImport requests.RoadCondition, files requests.RoadConditionFiles) (int, error) {

	// start transaction

	tx := t.roadConditionRepo.StartTransSection()

	CSVDir := os.Getenv("ROAD_CONDITION_CSV_DIR")
	ImgDir := os.Getenv("ROAD_CONDITION_ZIP_DIR")

	surveyedDate, err := time.Parse("2006-01-02 15:04:00", rcImport.SurveyedDate)
	if err != nil {
		return 0, responses.NewAppErr(400, err.Error())
	}
	rc := models.RoadCondition{
		SurveyedDate: surveyedDate,
		Year:         surveyedDate.Year(),
		LaneNo:       rcImport.LaneNo,
		CreatedBy:    int(uid),
		UpdatedBy:    int(uid),
		IDParent:     0,
		Remarks:      rcImport.Remarks,
		Status:       "A",
		Revision:     0,
	}

	var irifilePath, imgPath string

	if files.IriFilename.Filename != "" {

		irifilePath, err = helpers.SaveFile(c, files.IriFilename, CSVDir)
		if err != nil {
			return 0, responses.NewAppErr(500, err.Error())
		}
		err := helpers.SetPermissionsTo0775(CSVDir)
		if err != nil {
			return 0, responses.NewAppErr(500, err.Error())
		}
	}

	if files.ImageFilename != nil {

		imgPath, err = helpers.SaveFile(c, files.ImageFilename, ImgDir)
		if err != nil {
			return 0, responses.NewAppErr(500, err.Error())

		}

		err := helpers.SetPermissionsTo0775(ImgDir)
		if err != nil {
			return 0, responses.NewAppErr(500, err.Error())
		}
	}

	CSVData, err := readConfitionCSVFile(irifilePath)
	if err != nil {
		return 0, responses.NewAppErr(500, err.Error())

	}

	if len(CSVData) == 0 {
		return 0, nil
	}
	rc.RoadId = CSVData[0].RoadId
	rc.KmStart = CSVData[0].KMStart
	rc.KmEnd = CSVData[len(CSVData)-1].KMEnd

	if rc.RoadId != rcImport.RoadID {

		return 0, responses.NewAppErr(400, constants.INVALID_ROAD_CONDITION_GEOM_RANGE)
	}

	road, err := t.roadConditionRepo.GetRoadByID(rc.RoadId)
	if err != nil {
		return 0, responses.NewAppErr(500, err.Error())

	}

	fullGeom, err := t.roadConditionRepo.GetFullGeom(rc.RoadId, rc.LaneNo)
	if err != nil {
		return 0, responses.NewAppErr(500, err.Error())

	}

	// var getKmType int
	// switch road.RefDirectionId {
	// case 1:
	// 	getKmType = 2

	// case 2:
	// 	getKmType = 1

	// }

	// kmValues, err := t.roadConditionRepo.GetRoadKmRange(rc.RoadId, int(fullGeom.KmStart), int(fullGeom.KmEnd), getKmType)
	// if err != nil {
	// 	return 0, fmt.Errorf("fullGrom:", err.Error())
	// }

	err = validateData(CSVData, fullGeom, road.RoadInfo.RefDirectionId)
	if err != nil {
		return 0, err
	}

	// var kmItems []models.RoadConditionSurvey
	// var mItems []models.RoadConditionSurveyM
	//var roadData RoadData

	// for i, data := range CSVData {
	// 	var mItem models.RoadConditionSurveyM

	// 	mItem.KmStart = data.KMStart
	// 	mItem.KmEnd = data.KMEnd
	// 	mItem.IRI = data.IRI
	// 	mItem.IFI = data.IFI
	// 	mItem.MPD = data.MPD
	// 	mItem.RUT = data.RUT

	// 	subRoadStart := math.Abs(mItem.KmStart-fullGeom.KmStart) / math.Abs(fullGeom.KmStart-fullGeom.KmEnd)
	// 	subRoadEnd := math.Abs(mItem.KmEnd-fullGeom.KmStart) / math.Abs(fullGeom.KmStart-fullGeom.KmEnd)
	// 	subRoadMin := math.Min(subRoadStart, subRoadEnd)
	// 	subRoadMax := math.Max(subRoadStart, subRoadEnd)

	// 	mItem.TheGeom = fmt.Sprintf("ST_LineSubstring('%s', %f, %f)", fullGeom.Geom, subRoadMin, subRoadMax)

	// 	mItem.ImgFilepath = data.ImgFilepath

	// 	mItems = append(mItems, mItem)

	// 	roadData.IriAverage = helpers.CalculateRoadConditionAverage(roadData.IriAverage, mItem.IRI, mItem.KmStart, mItem.KmEnd)
	// 	roadData.MpdAverage = helpers.CalculateRoadConditionAverage(roadData.MpdAverage, mItem.MPD, mItem.KmStart, mItem.KmEnd)
	// 	roadData.RutAverage = helpers.CalculateRoadConditionAverage(roadData.RutAverage, mItem.RUT, mItem.KmStart, mItem.KmEnd)
	// 	roadData.IfiAverage = helpers.CalculateRoadConditionAverage(roadData.IfiAverage, mItem.IFI, mItem.KmStart, mItem.KmEnd)

	// 	roadData.IriKm, roadData.DividerCountIri = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountIri, roadData.IriKm, mItem.IRI, mItem.KmStart, mItem.KmEnd)
	// 	roadData.IfiKm, roadData.DividerCountIfi = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountIfi, roadData.IfiKm, mItem.IFI, mItem.KmStart, mItem.KmEnd)
	// 	roadData.MpdKm, roadData.DividerCountMpd = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountMpd, roadData.MpdKm, mItem.MPD, mItem.KmStart, mItem.KmEnd)
	// 	roadData.RutKm, roadData.DividerCountRut = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountRut, roadData.RutKm, mItem.RUT, mItem.KmStart, mItem.KmEnd)

	// 	roadData.TotalM += math.Abs(mItem.KmStart - mItem.KmEnd)

	// 	if int(data.KMEnd)%1000 == 0 || i == len(CSVData)-1 {

	// 		kmItem := models.RoadConditionSurvey{
	// 			KmStart: kmCountDivider,
	// 			KmEnd:   mItem.KmEnd,
	// 		}

	// 		if roadData.DividerCountIri != nil {
	// 			kmItem.IRI = helpers.CalculateRoadCondition(roadData.IriKm, roadData.DividerCountIri)
	// 		}

	// 		if roadData.DividerCountIfi != nil {
	// 			kmItem.IFI = helpers.CalculateRoadCondition(roadData.IfiKm, roadData.DividerCountIfi)
	// 		}

	// 		if roadData.DividerCountMpd != nil {
	// 			kmItem.MPD = helpers.CalculateRoadCondition(roadData.MpdKm, roadData.DividerCountMpd)
	// 		}

	// 		if roadData.DividerCountRut != nil {
	// 			kmItem.RUT = helpers.CalculateRoadCondition(roadData.RutKm, roadData.DividerCountRut)
	// 		}

	// 		kmItems = append(kmItems, kmItem)

	// 		kmCountDivider = mItem.KmEnd

	// 	}

	// }

	kmCountDivider := rc.KmStart
	rcAnalysisData, err := RoadConditionAnalysisData(CSVData, kmCountDivider, fullGeom)
	if err != nil {
		return 0, responses.NewAppErr(500, err.Error())
	}

	if rcAnalysisData.RoadData.TotalM >= 0 {
		rc.IRI = helpers.CalculateRoadCondition(rcAnalysisData.RoadData.IriAverage, &rcAnalysisData.RoadData.TotalM)
		rc.IFI = helpers.CalculateRoadCondition(rcAnalysisData.RoadData.IfiAverage, &rcAnalysisData.RoadData.TotalM)
		rc.MPD = helpers.CalculateRoadCondition(rcAnalysisData.RoadData.MpdAverage, &rcAnalysisData.RoadData.TotalM)
		rc.RUT = helpers.CalculateRoadCondition(rcAnalysisData.RoadData.RutAverage, &rcAnalysisData.RoadData.TotalM)
	}

	roundedStruct := helpers.RoundStructPointerFloats(&rc, 3).(*models.RoadCondition)
	rc = *roundedStruct

	rc.CreatedDate = time.Now()
	rc.CreatedBy = int(uid)
	rc.UpdatedDate = time.Now()
	rcId, err := t.roadConditionRepo.CreateRoadCondition(tx, rc)
	if err != nil {
		logs.Error(err.Error())
		t.roadConditionRepo.RollBack(tx)
		return 0, responses.NewAppErr(500, err.Error())
	}

	filePath := fmt.Sprintf("storages/road/attachments/%d/condition/%d/", rc.RoadId, rcId)
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		t.roadConditionRepo.RollBack(tx)
		return 0, responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
	}

	var fullImgPath string
	if imgPath != "" {
		err = helpers.Unzip(imgPath, filePath)
		if err != nil {
			logs.Error(constants.FAILED_TO_SAVE_FILE)
			t.roadConditionRepo.RollBack(tx)
			return 0, responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
		}

		fileName := filepath.Base(imgPath)

		fullImgPath = filePath + fileName
		// move file zip or rar
		err = os.Rename(imgPath, fullImgPath)
		if err != nil {
			logs.Error(err.Error())
			t.roadConditionRepo.RollBack(tx)
			return 0, responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
		}

	}

	iriFileName := filepath.Base(irifilePath)
	// update road damage csv file path
	fullCsvName := filePath + iriFileName
	// move file csv upload
	if err := helpers.MoveFile(irifilePath, fullCsvName); err != nil {
		logs.Error(err)
		t.roadConditionRepo.RollBack(tx)
		return 0, responses.NewAppErr(500, err.Error())
	}

	err = t.roadConditionRepo.UpdateRoadConditionFilepath(tx, rcId, fullCsvName, fullImgPath)
	if err != nil {
		t.roadConditionRepo.RollBack(tx)
		logs.Error(err.Error())
		return 0, responses.NewAppErr(500, err.Error())
	}

	kmList := make(map[string]int)
	for _, kmItem := range rcAnalysisData.RoadConditionSurvey {
		kmItem.RoadConditionId = rcId

		roundedStruct := helpers.RoundStructPointerFloats(&kmItem, 3).(*models.RoadConditionSurvey)
		//roundedStruct := roundedStructValue.Interface().(*models.RoadConditionSurvey)
		kmItem = *roundedStruct

		rcServayId, err := t.roadConditionRepo.CreateRoadConditionSurvey(tx, kmItem)
		if err != nil {
			logs.Error(err.Error())
			t.roadConditionRepo.RollBack(tx)
			return 0, responses.NewAppErr(500, err.Error())
		}
		km := fmt.Sprintf("%v-%v", kmItem.KmStart, kmItem.KmEnd)
		kmList[km] = rcServayId

	}

	m100List := make(map[string]int)
	for km, rcSurveyId := range kmList {

		kmRange := strings.Split(km, "-")
		kmStart, _ := strconv.ParseFloat(kmRange[0], 32)
		kmEnd, _ := strconv.ParseFloat(kmRange[1], 32)
		subRoadMin := helpers.Min(kmStart, kmEnd)
		subRoadMax := helpers.Max(kmStart, kmEnd)

		for _, mItem100 := range rcAnalysisData.RoadConditionSurvey100M {

			switch road.RoadInfo.RefDirectionId {
			case 1:

				if mItem100.KmEnd > float64(subRoadMin) && mItem100.KmEnd <= float64(subRoadMax) {
					mItem100.RoadConditionSurveyId = rcSurveyId
				}
			case 2:
				if mItem100.KmStart > float64(subRoadMin) && mItem100.KmStart <= float64(subRoadMax) {
					mItem100.RoadConditionSurveyId = rcSurveyId

				}
			}
			if mItem100.RoadConditionSurveyId == 0 {
				continue
			}

			roundedStruct := helpers.RoundStructPointerFloats(&mItem100, 3).(*models.RoadConditionSurvey100M)
			//roundedStruct := roundedStructValue.Interface().(*models.RoadConditionSurvey)
			mItem100 = *roundedStruct

			rcSurvey100mId, err := t.roadConditionRepo.CreateRoadConditionSurvey100M(tx, mItem100)
			if err != nil {
				logs.Error(err.Error())
				t.roadConditionRepo.RollBack(tx)
				return 0, responses.NewAppErr(500, err.Error())
			}
			m100 := fmt.Sprintf("%v-%v", mItem100.KmStart, mItem100.KmEnd)
			m100List[m100] = rcSurvey100mId

		}
	}

	for m100, rcServayId := range m100List {
		// Check if mItem km range overlaps with kmList[km] km range
		m100Range := strings.Split(m100, "-")
		m100Start, _ := strconv.ParseFloat(m100Range[0], 32)
		m100End, _ := strconv.ParseFloat(m100Range[1], 32)
		subRoadMin := helpers.Min(m100Start, m100End)
		subRoadMax := helpers.Max(m100Start, m100End)
		for _, mItem := range rcAnalysisData.RoadConditionSurveyM {

			switch road.RoadInfo.RefDirectionId {
			case 1:

				if mItem.KmEnd > float64(subRoadMin) && mItem.KmEnd <= float64(subRoadMax) {
					mItem.RoadConditionSurvay100mID = rcServayId
				}
			case 2:
				if mItem.KmStart > float64(subRoadMin) && mItem.KmStart <= float64(subRoadMax) {
					mItem.RoadConditionSurvay100mID = rcServayId

				}
			}

			if mItem.RoadConditionSurvay100mID == 0 {
				continue
			}

			var imgPath string
			if mItem.ImgFilepath != "" {
				imgPath = filePath + mItem.ImgFilepath
			}
			if mItem.RoadConditionSurvay100mID != 0 {
				mItem.ImgFilepath = imgPath

				roundedStruct := helpers.RoundStructPointerFloats(&mItem, 3).(*models.RoadConditionSurveyM)
				//roundedStructValue := helpers.RoundStructFloats(&mItem, 3)
				//roundedStruct := roundedStructValue.Interface().(*models.RoadConditionSurveyM)
				mItem = *roundedStruct

				err := t.roadConditionRepo.CreateRoadConditionSurveyM(tx, mItem)
				if err != nil {
					logs.Error(err.Error())
					t.roadConditionRepo.RollBack(tx)
					return 0, responses.NewAppErr(500, err.Error())
				}
			}

		}

	}
	t.roadConditionRepo.Commit(tx)

	return rcId, nil
}

func RoadConditionAnalysisData(csvDatas []models.RoadConditionCSV, countDivider float64, fullGeom models.FullGeom) (responses.RoadConditionAnalysisData, error) {

	// var kmItem models.RoadConditionSurvey
	// var 100mItem models.RoadConditionSurvey100M

	var currentType string
	var currentCount int
	surveyTypeData := make(map[string][]models.RoadConditionCSV)

	for _, csvData := range csvDatas {
		if csvData.SurveyType != currentType {
			currentType = csvData.SurveyType
			currentCount++
		}

		newSurveyType := fmt.Sprintf("%s-%d", csvData.SurveyType, currentCount)
		surveyTypeData[newSurveyType] = append(surveyTypeData[newSurveyType], csvData)

	}

	var resps responses.RoadConditionAnalysisData

	var roadData responses.RoadData
	for surveyType, dataSlice := range surveyTypeData {

		parts := strings.Split(surveyType, "-")
		surveyType := parts[0]

		m100CountDivider := dataSlice[0].KMStart
		kmCountDivider := dataSlice[0].KMStart

		for i, data := range dataSlice {

			var mItem models.RoadConditionSurveyM
			mItem.KmStart = data.KMStart
			mItem.KmEnd = data.KMEnd
			mItem.IRI = data.IRI
			mItem.IFI = data.IFI
			mItem.MPD = data.MPD
			mItem.RUT = data.RUT

			subRoadStart := math.Abs(mItem.KmStart-fullGeom.KmStart) / math.Abs(fullGeom.KmStart-fullGeom.KmEnd)
			subRoadEnd := math.Abs(mItem.KmEnd-fullGeom.KmStart) / math.Abs(fullGeom.KmStart-fullGeom.KmEnd)
			subRoadMin := math.Min(subRoadStart, subRoadEnd)
			subRoadMax := math.Max(subRoadStart, subRoadEnd)

			mItem.TheGeom = fmt.Sprintf("ST_LineSubstring('%s', %f, %f)", fullGeom.Geom, subRoadMin, subRoadMax)

			mItem.ImgFilepath = data.ImgFilepath
			mItem.SurveyType = surveyType
			resps.RoadConditionSurveyM = append(resps.RoadConditionSurveyM, mItem)

			roadData.IriAverage = helpers.CalculateRoadConditionAverage(roadData.IriAverage, mItem.IRI, mItem.KmStart, mItem.KmEnd)
			roadData.MpdAverage = helpers.CalculateRoadConditionAverage(roadData.MpdAverage, mItem.MPD, mItem.KmStart, mItem.KmEnd)
			roadData.RutAverage = helpers.CalculateRoadConditionAverage(roadData.RutAverage, mItem.RUT, mItem.KmStart, mItem.KmEnd)
			roadData.IfiAverage = helpers.CalculateRoadConditionAverage(roadData.IfiAverage, mItem.IFI, mItem.KmStart, mItem.KmEnd)

			roadData.TotalM += math.Abs(mItem.KmStart - mItem.KmEnd)

			roadData.IriKm, roadData.DividerCountIri = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountIri, roadData.IriKm, mItem.IRI, mItem.KmStart, mItem.KmEnd)
			roadData.MpdKm, roadData.DividerCountMpd = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountMpd, roadData.MpdKm, mItem.MPD, mItem.KmStart, mItem.KmEnd)
			roadData.RutKm, roadData.DividerCountRut = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountRut, roadData.RutKm, mItem.RUT, mItem.KmStart, mItem.KmEnd)
			roadData.IfiKm, roadData.DividerCountIfi = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountIfi, roadData.IfiKm, mItem.IFI, mItem.KmStart, mItem.KmEnd)

			roadData.Iri100m, roadData.DividerCount100mIri = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCount100mIri, roadData.Iri100m, mItem.IRI, mItem.KmStart, mItem.KmEnd)
			roadData.Ifi100m, roadData.DividerCount100mIfi = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCount100mIfi, roadData.Ifi100m, mItem.IFI, mItem.KmStart, mItem.KmEnd)
			roadData.Mpd100m, roadData.DividerCount100mMpd = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCount100mMpd, roadData.Mpd100m, mItem.MPD, mItem.KmStart, mItem.KmEnd)
			roadData.Rut100m, roadData.DividerCount100mRut = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCount100mRut, roadData.Rut100m, mItem.RUT, mItem.KmStart, mItem.KmEnd)

			if int(data.KMEnd)%100 == 0 || i == len(dataSlice)-1 {

				mItem100 := models.RoadConditionSurvey100M{
					KmStart: m100CountDivider,
					KmEnd:   mItem.KmEnd,
				}

				subRoadStart := math.Abs(mItem100.KmStart-fullGeom.KmStart) / math.Abs(fullGeom.KmStart-fullGeom.KmEnd)
				subRoadEnd := math.Abs(mItem100.KmEnd-fullGeom.KmStart) / math.Abs(fullGeom.KmStart-fullGeom.KmEnd)
				subRoadMin := math.Min(subRoadStart, subRoadEnd)
				subRoadMax := math.Max(subRoadStart, subRoadEnd)
				mItem100.TheGeom = fmt.Sprintf("ST_LineSubstring('%s', %f, %f)", fullGeom.Geom, subRoadMin, subRoadMax)

				if roadData.DividerCount100mIri != nil {
					mItem100.IRI = helpers.CalculateRoadCondition(roadData.Iri100m, roadData.DividerCount100mIri)
				}

				if roadData.DividerCount100mIfi != nil {
					mItem100.IFI = helpers.CalculateRoadCondition(roadData.Ifi100m, roadData.DividerCount100mIfi)
				}

				if roadData.DividerCount100mMpd != nil {
					mItem100.MPD = helpers.CalculateRoadCondition(roadData.Mpd100m, roadData.DividerCount100mMpd)
				}

				if roadData.DividerCount100mRut != nil {
					mItem100.RUT = helpers.CalculateRoadCondition(roadData.Rut100m, roadData.DividerCount100mRut)
				}
				mItem100.SurveyType = surveyType

				resps.RoadConditionSurvey100M = append(resps.RoadConditionSurvey100M, mItem100)

				//reset value
				*roadData.Iri100m = 0
				*roadData.DividerCount100mIri = 0
				*roadData.Ifi100m = 0
				*roadData.DividerCount100mIfi = 0
				*roadData.Mpd100m = 0
				*roadData.DividerCount100mMpd = 0
				*roadData.Rut100m = 0
				*roadData.DividerCount100mRut = 0

				m100CountDivider = mItem100.KmEnd

			}

			if int(data.KMEnd)%1000 == 0 || i == len(dataSlice)-1 {

				kmItem := models.RoadConditionSurvey{
					KmStart: kmCountDivider,
					KmEnd:   mItem.KmEnd,
				}

				subRoadStart := math.Abs(kmItem.KmStart-fullGeom.KmStart) / math.Abs(fullGeom.KmStart-fullGeom.KmEnd)
				subRoadEnd := math.Abs(kmItem.KmEnd-fullGeom.KmStart) / math.Abs(fullGeom.KmStart-fullGeom.KmEnd)
				subRoadMin := math.Min(subRoadStart, subRoadEnd)
				subRoadMax := math.Max(subRoadStart, subRoadEnd)
				kmItem.TheGeom = fmt.Sprintf("ST_LineSubstring('%s', %f, %f)", fullGeom.Geom, subRoadMin, subRoadMax)

				kmItem.IRI = roadData.IriAverage
				if roadData.DividerCountIri != nil {

					kmItem.IRI = helpers.CalculateRoadCondition(roadData.IriKm, roadData.DividerCountIri)
				}

				if roadData.DividerCountIfi != nil {
					kmItem.IFI = helpers.CalculateRoadCondition(roadData.IfiKm, roadData.DividerCountIfi)
				}

				if roadData.DividerCountMpd != nil {
					kmItem.MPD = helpers.CalculateRoadCondition(roadData.MpdKm, roadData.DividerCountMpd)
				}

				if roadData.DividerCountRut != nil {
					kmItem.RUT = helpers.CalculateRoadCondition(roadData.RutKm, roadData.DividerCountRut)
				}
				kmItem.SurveyType = surveyType

				resps.RoadConditionSurvey = append(resps.RoadConditionSurvey, kmItem)

				kmCountDivider = mItem.KmEnd

				//reset value
				*roadData.IriKm = 0
				*roadData.DividerCountIri = 0
				*roadData.IfiKm = 0
				*roadData.DividerCountIfi = 0
				*roadData.MpdKm = 0
				*roadData.DividerCountMpd = 0
				*roadData.RutKm = 0
				*roadData.DividerCountRut = 0

			}
			resps.RoadData = roadData
		}

	}

	return resps, nil
}

func validateData(csvData []models.RoadConditionCSV, fullGeom models.FullGeom, directionID int) error {

	//	roadID := csvData[0].RoadId

	kmStart := csvData[0].KMStart
	kmEnd := csvData[len(csvData)-1].KMEnd
	switch directionID {
	case 1:
		if kmStart < fullGeom.KmStart || kmEnd > fullGeom.KmEnd {

			return responses.NewAppErr(400, constants.INVALID_ROAD_CONDITION_GEOM_RANGE)

		}
	case 2:
		if kmStart > fullGeom.KmStart || kmEnd < fullGeom.KmEnd {

			return responses.NewAppErr(400, constants.INVALID_ROAD_CONDITION_GEOM_RANGE)

		}
	}

	var rc models.RoadConditionCSV
	for i := 0; i < len(csvData); i++ {
		row := i + 1

		if csvData[i].RoadId == rc.RoadId ||
			csvData[i].RoadCode == rc.RoadCode ||
			csvData[i].Name == rc.Name ||
			csvData[i].IFI == nil ||
			csvData[i].IRI == nil ||
			csvData[i].MPD == nil ||
			csvData[i].RUT == nil {

			return responses.NewAppErr(400, strings.Replace(constants.INVALID_ROAD_CONDITION_VALUE, "_", fmt.Sprintf("%d", row), -1))
		}
		if *csvData[i].IFI < 0 ||
			*csvData[i].IRI < 0 ||
			*csvData[i].MPD < 0 ||
			*csvData[i].RUT < 0 {

			return responses.NewAppErr(400, strings.Replace(constants.INVALID_ROAD_CONDITION_VALUE, "_", fmt.Sprintf("%d", row), -1))
		}

	}

	return nil
}

func (t *roadConditionUseCase) UpdateRoadCondition(c *gin.Context, uid uint, IDParent int, rcImport requests.RoadConditionUpdate, files requests.RoadConditionFiles, iriFilenameStatus, imageFilenameStatus string) (interface{}, error) {
	CSVDir := os.Getenv("ROAD_CONDITION_CSV_DIR")
	ImgDir := os.Getenv("ROAD_CONDITION_IMG_DIR")

	// start transaction

	tx := t.roadConditionRepo.StartTransSection()
	// no_file
	//    csv [validate], zip[path รูปต้องเป็น path เดิม]
	// delete
	//     csv [validate], zip[ลบpathรูป, ข้อมูล รูปหน้ากล้อง road_condition  ต้องเป็นค่าว่างทั้งหมด]
	// not_edit
	//    csv[ไม่มีการแก้ไขข้อมูล],zip[ไม่มีการแก้ไขข้อมูล]
	// upload
	//    csv[update ข้อมูล road_condition],zip[อัพเดท path img และ ข้อมูลรูปหน้ากล้อง]

	surveyedDate, err := time.Parse("2006-01-02 15:04:00", rcImport.SurveyedDate)
	if err != nil {
		return responses.RoadConditionUpdate{}, responses.NewAppErr(400, err.Error())
	}

	roadCondition, err := t.roadConditionRepo.GetAllRoadConditionByIdParent(IDParent)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return 0, responses.NewAppErr(500, err.Error())

	}

	if rcImport.LaneNo != roadCondition.LaneNo {
		return responses.RoadConditionUpdate{}, responses.NewAppErr(400, constants.INVALID_ROAD_CONDITION_LANE)
	}

	rc := models.RoadCondition{
		KmStart:          roadCondition.KmStart,
		KmEnd:            roadCondition.KmEnd,
		RoadId:           roadCondition.RoadId,
		SurveyedDate:     surveyedDate,
		Year:             surveyedDate.Year(),
		LaneNo:           roadCondition.LaneNo,
		CreatedDate:      roadCondition.CreatedDate,
		UpdatedDate:      time.Now(),
		UpdatedBy:        int(uid),
		CreatedBy:        roadCondition.CreatedBy,
		IDParent:         IDParent,
		Remarks:          rcImport.Remarks,
		Status:           "A",
		Revision:         roadCondition.Revision + 1,
		IRI:              roadCondition.RoadCondition.IRI,
		IFI:              roadCondition.RoadCondition.IFI,
		MPD:              roadCondition.RoadCondition.MPD,
		RUT:              roadCondition.RoadCondition.RUT,
		IRIInputFilePath: roadCondition.IRIInputFilePath,
		ImgFilePath:      roadCondition.ImgFilePath,
	}

	var resp responses.RoadConditionUpdate
	var irifilePath, imgPath string

	switch iriFilenameStatus {
	case "upload":
		irifilePath, err = helpers.SaveFile(c, files.IriFilename, CSVDir)
		if err != nil {
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
		}
		err := helpers.SetPermissionsTo0775(CSVDir)
		if err != nil {
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
		}

		if roadCondition.Status == "A" {

			err = t.roadConditionRepo.UpdateStatusIByID(tx, roadCondition.ID, uid)
			if err != nil {
				logs.Error(err.Error())
				t.roadConditionRepo.RollBack(tx)
				return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
			}
		}
	default:

		// hariphan แก้ เวลาไม่มีไฟล์ csv, zip update ให้ update เฉพาะ surveyed_date, และremarks
		data := models.RoadCondition{
			ID:           roadCondition.ID,
			SurveyedDate: surveyedDate,
			Year:         surveyedDate.Year(),
			UpdatedDate:  time.Now(),
			UpdatedBy:    int(uid),
			Remarks:      rcImport.Remarks,
		}
		_, err = t.roadConditionRepo.UpdateRoadConditionNoIriFile(tx, data)
		if err != nil {
			logs.Error(err.Error())
			t.roadConditionRepo.RollBack(tx)
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
		}

		resp.ID = roadCondition.ID
		resp.IDParent = IDParent

	}

	switch imageFilenameStatus {

	case "delete":

		err = helpers.RemoveFileIfExists(roadCondition.ImgFilePath)
		if err != nil {
			logs.Error(err)
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
		}

		err = t.roadConditionRepo.UpdateRoadConditionImgPath(roadCondition.ID, "")
		if err != nil {
			logs.Error(err.Error())
			t.roadConditionRepo.RollBack(tx)
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
		}

		if iriFilenameStatus == "not_edit" {
			resp.ID = roadCondition.ID
			resp.IDParent = IDParent
			return resp, nil
		}

	case "upload":

		imgPath, err = helpers.SaveFile(c, files.ImageFilename, ImgDir)
		if err != nil {
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
		}

		if iriFilenameStatus == "not_edit" {

			var fullImgPath string

			filePath := fmt.Sprintf("storages/road/attachments/%d/condition/%d/", rc.RoadId, roadCondition.ID)
			if err := os.MkdirAll(filePath, 0775); err != nil {
				logs.Error(err)
				return responses.RoadConditionUpdate{}, responses.NewAppErr(500, constants.FAILED_TO_SAVE_FILE)
			}

			err = helpers.Unzip(imgPath, filePath)
			if err != nil {
				logs.Error(constants.FAILED_TO_SAVE_FILE)
				return responses.RoadConditionUpdate{}, responses.NewAppErr(500, constants.FAILED_TO_SAVE_FILE)
			}

			fileName := filepath.Base(imgPath)

			fullImgPath = filePath + fileName
			// move file zip or rar
			err = os.Rename(imgPath, fullImgPath)
			if err != nil {
				logs.Error(err.Error())
				return responses.RoadConditionUpdate{}, responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
			}

			err := t.roadConditionRepo.UpdateRoadConditionImgPath(roadCondition.ID, fullImgPath)
			if err != nil {
				logs.Error(err.Error())
				t.roadConditionRepo.RollBack(tx)
				return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
			}

		}
	}

	if iriFilenameStatus == "upload" {
		CSVData, err := readConfitionCSVFile(irifilePath)
		if err != nil {
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())

		}
		if len(CSVData) == 0 {
			return responses.RoadConditionUpdate{}, nil
		}
		rc.RoadId = CSVData[0].RoadId
		rc.KmStart = CSVData[0].KMStart
		rc.KmEnd = CSVData[len(CSVData)-1].KMEnd

		if rc.RoadId != rcImport.RoadID {
			return responses.RoadConditionUpdate{}, responses.NewAppErr(400, constants.INVALID_ROAD_CONDITION_GEOM_RANGE)
		}
		road, err := t.roadConditionRepo.GetRoadByID(rc.RoadId)
		if err != nil {
			return responses.RoadConditionUpdate{}, err

		}

		fullGeom, err := t.roadConditionRepo.GetFullGeom(rc.RoadId, rc.LaneNo)
		if err != nil {
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())

		}

		// var getKmType int
		// switch road.RefDirectionId {
		// case 1:
		// 	getKmType = 2

		// case 2:
		// 	getKmType = 1

		// }

		// kmValues, err := t.roadConditionRepo.GetRoadKmRange(rc.RoadId, int(fullGeom.KmStart), int(fullGeom.KmEnd), getKmType)
		// if err != nil {
		// 	return responses.RoadConditionUpdate{}, fmt.Errorf("fullGrom:", err.Error())
		// }

		err = validateData(CSVData, fullGeom, road.RoadInfo.RefDirectionId)
		if err != nil {
			return responses.RoadConditionUpdate{}, err
		}

		// rcGrades, err := t.roadConditionRepo.GetRoadConditionGrades(rc.RoadId)
		// if err != nil {
		// 	return responses.RoadConditionUpdate{}, fmt.Errorf("rcGrades:", err.Error())
		// }

		kmCountDivider := rc.KmStart
		rcAnalysisData, err := RoadConditionAnalysisData(CSVData, kmCountDivider, fullGeom)
		if err != nil {
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
		}

		if rcAnalysisData.RoadData.TotalM >= 0 {
			rc.IRI = helpers.CalculateRoadCondition(rcAnalysisData.RoadData.IriAverage, &rcAnalysisData.RoadData.TotalM)
			rc.IFI = helpers.CalculateRoadCondition(rcAnalysisData.RoadData.IfiAverage, &rcAnalysisData.RoadData.TotalM)
			rc.MPD = helpers.CalculateRoadCondition(rcAnalysisData.RoadData.MpdAverage, &rcAnalysisData.RoadData.TotalM)
			rc.RUT = helpers.CalculateRoadCondition(rcAnalysisData.RoadData.RutAverage, &rcAnalysisData.RoadData.TotalM)
		}

		// kmCountDivider := rc.KmStart
		// var kmItems []models.RoadConditionSurvey
		// var mItems []models.RoadConditionSurveyM
		// var roadData responses.RoadData
		// for i, data := range CSVData {
		// 	var mItem models.RoadConditionSurveyM

		// 	mItem.KmStart = data.KMStart
		// 	mItem.KmEnd = data.KMEnd
		// 	mItem.IFI = data.IFI
		// 	mItem.IRI = data.IRI
		// 	mItem.MPD = data.MPD
		// 	mItem.RUT = data.RUT

		// 	subRoadStart := math.Abs(mItem.KmStart-fullGeom.KmStart) / math.Abs(fullGeom.KmStart-fullGeom.KmEnd)
		// 	subRoadEnd := math.Abs(mItem.KmEnd-fullGeom.KmStart) / math.Abs(fullGeom.KmStart-fullGeom.KmEnd)
		// 	subRoadMin := math.Min(subRoadStart, subRoadEnd)
		// 	subRoadMax := math.Max(subRoadStart, subRoadEnd)

		// 	mItem.TheGeom = fmt.Sprintf("ST_LineSubstring('%s', %f, %f)", fullGeom.Geom, subRoadMin, subRoadMax)

		// 	mItem.ImgFilepath = data.ImgFilepath

		// 	mItems = append(mItems, mItem)

		// 	roadData.IriAverage = helpers.CalculateRoadConditionAverage(roadData.IriAverage, mItem.IRI, mItem.KmStart, mItem.KmEnd)
		// 	roadData.MpdAverage = helpers.CalculateRoadConditionAverage(roadData.MpdAverage, mItem.MPD, mItem.KmStart, mItem.KmEnd)
		// 	roadData.RutAverage = helpers.CalculateRoadConditionAverage(roadData.RutAverage, mItem.RUT, mItem.KmStart, mItem.KmEnd)
		// 	roadData.IfiAverage = helpers.CalculateRoadConditionAverage(roadData.IfiAverage, mItem.IFI, mItem.KmStart, mItem.KmEnd)

		// 	roadData.IriKm, roadData.DividerCountIri = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountIri, roadData.IriKm, mItem.IRI, mItem.KmStart, mItem.KmEnd)
		// 	roadData.IfiKm, roadData.DividerCountIfi = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountIfi, roadData.IfiKm, mItem.IFI, mItem.KmStart, mItem.KmEnd)
		// 	roadData.MpdKm, roadData.DividerCountMpd = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountMpd, roadData.MpdKm, mItem.MPD, mItem.KmStart, mItem.KmEnd)
		// 	roadData.RutKm, roadData.DividerCountRut = helpers.CalculateRoadConditionAverageAndCount(roadData.DividerCountRut, roadData.RutKm, mItem.RUT, mItem.KmStart, mItem.KmEnd)

		// 	roadData.TotalM += math.Abs(mItem.KmStart - mItem.KmEnd)

		// 	if int(data.KMEnd)%1000 == 0 || i == len(CSVData)-1 {

		// 		kmItem := models.RoadConditionSurvey{
		// 			KmStart: kmCountDivider,
		// 			KmEnd:   mItem.KmEnd,
		// 		}

		// 		if roadData.DividerCountIri != nil {
		// 			kmItem.IRI = helpers.CalculateRoadCondition(roadData.IriKm, roadData.DividerCountIri)
		// 		}

		// 		if roadData.DividerCountIfi != nil {
		// 			kmItem.IFI = helpers.CalculateRoadCondition(roadData.IfiKm, roadData.DividerCountIfi)
		// 		}

		// 		if roadData.DividerCountMpd != nil {
		// 			kmItem.MPD = helpers.CalculateRoadCondition(roadData.MpdKm, roadData.DividerCountMpd)
		// 		}

		// 		if roadData.DividerCountRut != nil {
		// 			kmItem.RUT = helpers.CalculateRoadCondition(roadData.RutKm, roadData.DividerCountRut)
		// 		}

		// 		kmItems = append(kmItems, kmItem)
		// 		kmCountDivider = mItem.KmEnd

		// 	}
		// }

		// if roadData.TotalM >= 0 {
		// 	rc.IRI = helpers.CalculateRoadCondition(roadData.IriAverage, &roadData.TotalM)
		// 	rc.IFI = helpers.CalculateRoadCondition(roadData.IfiAverage, &roadData.TotalM)
		// 	rc.MPD = helpers.CalculateRoadCondition(roadData.MpdAverage, &roadData.TotalM)
		// 	rc.RUT = helpers.CalculateRoadCondition(roadData.RutAverage, &roadData.TotalM)
		// }

		roundedStruct := helpers.RoundStructPointerFloats(&rc, 3).(*models.RoadCondition)
		rc = *roundedStruct

		rc.CreatedBy = roadCondition.CreatedBy
		rc.CreatedDate = roadCondition.CreatedDate
		rc.UpdatedBy = int(uid)
		rc.UpdatedDate = time.Now()
		rc.IDParent = IDParent

		if imageFilenameStatus == "delete" {
			rc.ImgFilePath = ""
		}

		rcResp, err := t.roadConditionRepo.UpdateRoadCondition(tx, rc)
		if err != nil {
			logs.Error(err.Error())
			t.roadConditionRepo.RollBack(tx)
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())

		}

		filePath := fmt.Sprintf("storages/road/attachments/%d/condition/%d/", rc.RoadId, rcResp.ID)
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			logs.Error(err)
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, constants.FAILED_TO_SAVE_FILE)
		}

		var fullImgPath string
		if imageFilenameStatus == "upload" {
			err = helpers.Unzip(imgPath, filePath)
			if err != nil {
				logs.Error(constants.FAILED_TO_SAVE_FILE)
				return responses.RoadConditionUpdate{}, responses.NewAppErr(500, constants.FAILED_TO_SAVE_FILE)
			}

			fileName := filepath.Base(imgPath)

			fullImgPath = filePath + fileName
			// move file zip or rar
			err = os.Rename(imgPath, fullImgPath)
			if err != nil {
				logs.Error(err.Error())
				return responses.RoadConditionUpdate{}, responses.NewAppErr(500, constants.FAILED_TO_SAVE_FILE)
			}

		}

		iriFileName := filepath.Base(irifilePath)
		// update road damage csv file path
		fullCsvName := filePath + iriFileName

		// move file csv upload
		if err := helpers.MoveFile(irifilePath, fullCsvName); err != nil {
			logs.Error(err)
			return responses.RoadConditionUpdate{}, responses.NewAppErr(400, err.Error())
		}

		err = t.roadConditionRepo.UpdateRoadConditionFilepath(tx, rcResp.ID, fullCsvName, fullImgPath)
		if err != nil {
			logs.Error(err.Error())
			t.roadConditionRepo.RollBack(tx)
			return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
		}

		kmList := make(map[string]int)
		for _, kmItem := range rcAnalysisData.RoadConditionSurvey {
			kmItem.RoadConditionId = rcResp.ID

			roundedStruct := helpers.RoundStructPointerFloats(&kmItem, 3).(*models.RoadConditionSurvey)
			//roundedStruct := roundedStructValue.Interface().(*models.RoadConditionSurvey)
			kmItem = *roundedStruct

			rcServayId, err := t.roadConditionRepo.CreateRoadConditionSurvey(tx, kmItem)
			if err != nil {
				logs.Error(err.Error())
				t.roadConditionRepo.RollBack(tx)
				return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
			}
			km := fmt.Sprintf("%v-%v", kmItem.KmStart, kmItem.KmEnd)
			kmList[km] = rcServayId

		}

		m100List := make(map[string]int)
		for km, rcSurveyId := range kmList {

			kmRange := strings.Split(km, "-")
			kmStart, _ := strconv.ParseFloat(kmRange[0], 32)
			kmEnd, _ := strconv.ParseFloat(kmRange[1], 32)
			subRoadMin := helpers.Min(kmStart, kmEnd)
			subRoadMax := helpers.Max(kmStart, kmEnd)

			for _, mItem100 := range rcAnalysisData.RoadConditionSurvey100M {

				switch road.RoadInfo.RefDirectionId {
				case 1:

					if mItem100.KmEnd > float64(subRoadMin) && mItem100.KmEnd <= float64(subRoadMax) {
						mItem100.RoadConditionSurveyId = rcSurveyId
					}
				case 2:
					if mItem100.KmStart > float64(subRoadMin) && mItem100.KmStart <= float64(subRoadMax) {
						mItem100.RoadConditionSurveyId = rcSurveyId

					}
				}
				if mItem100.RoadConditionSurveyId == 0 {
					continue
				}

				roundedStruct := helpers.RoundStructPointerFloats(&mItem100, 3).(*models.RoadConditionSurvey100M)
				//roundedStruct := roundedStructValue.Interface().(*models.RoadConditionSurvey)
				mItem100 = *roundedStruct

				rcSurvey100mId, err := t.roadConditionRepo.CreateRoadConditionSurvey100M(tx, mItem100)
				if err != nil {
					logs.Error(err.Error())
					t.roadConditionRepo.RollBack(tx)
					return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
				}
				m100 := fmt.Sprintf("%v-%v", mItem100.KmStart, mItem100.KmEnd)
				m100List[m100] = rcSurvey100mId

			}
		}

		for m100, rcServayId := range m100List {
			// Check if mItem km range overlaps with kmList[km] km range
			m100Range := strings.Split(m100, "-")
			m100Start, _ := strconv.ParseFloat(m100Range[0], 32)
			m100End, _ := strconv.ParseFloat(m100Range[1], 32)
			subRoadMin := helpers.Min(m100Start, m100End)
			subRoadMax := helpers.Max(m100Start, m100End)
			for _, mItem := range rcAnalysisData.RoadConditionSurveyM {

				switch road.RoadInfo.RefDirectionId {
				case 1:

					if mItem.KmEnd > float64(subRoadMin) && mItem.KmEnd <= float64(subRoadMax) {
						mItem.RoadConditionSurvay100mID = rcServayId
					}
				case 2:
					if mItem.KmStart > float64(subRoadMin) && mItem.KmStart <= float64(subRoadMax) {
						mItem.RoadConditionSurvay100mID = rcServayId

					}
				}

				if mItem.RoadConditionSurvay100mID == 0 {
					continue
				}

				var imgPath string
				if mItem.ImgFilepath != "" {
					imgPath = filePath + mItem.ImgFilepath
				}
				if mItem.RoadConditionSurvay100mID != 0 {
					mItem.ImgFilepath = imgPath

					roundedStruct := helpers.RoundStructPointerFloats(&mItem, 3).(*models.RoadConditionSurveyM)
					//roundedStructValue := helpers.RoundStructFloats(&mItem, 3)
					//roundedStruct := roundedStructValue.Interface().(*models.RoadConditionSurveyM)
					mItem = *roundedStruct

					err := t.roadConditionRepo.CreateRoadConditionSurveyM(tx, mItem)
					if err != nil {
						logs.Error(err.Error())
						t.roadConditionRepo.RollBack(tx)
						return responses.RoadConditionUpdate{}, responses.NewAppErr(500, err.Error())
					}
				}

			}

		}
		t.roadConditionRepo.Commit(tx)
		resp.ID = rcResp.ID
		resp.IDParent = IDParent
	}
	return resp, nil
}

func readConfitionCSVFile(filePath string) ([]models.RoadConditionCSV, error) {
	isFirstRow := true
	headerMap := make(map[string]int)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	var roadCondtionCsvs []models.RoadConditionCSV
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

		roadCondtionCsvs = append(roadCondtionCsvs, models.RoadConditionCSV{
			RoadId:      helpers.StrToInt(record[headerMap["road_id"]]),
			RoadCode:    record[headerMap["road_code"]],
			Name:        record[headerMap["name"]],
			KMStart:     helpers.StrToFloatValidate(record[headerMap["km_start"]]),
			KMEnd:       helpers.StrToFloatValidate(record[headerMap["km_end"]]),
			IRI:         helpers.StrToFloatPointerValidate(record[headerMap["iri"]]),
			MPD:         helpers.StrToFloatPointerValidate(record[headerMap["mpd"]]),
			RUT:         helpers.StrToFloatPointerValidate(record[headerMap["rut"]]),
			IFI:         helpers.StrToFloatPointerValidate(record[headerMap["ifi"]]),
			SurveyType:  record[headerMap["survey_type"]],
			ImgFilepath: record[headerMap["img_filepath"]],
		})

	}
	return roadCondtionCsvs, nil

}

// func checkContiion(value float64, rcGrade models.RoadConditionGrade) (int, error) {

// 	var leftCheck, rightCheck bool
// 	switch rcGrade.LeftCondition {
// 	case "<":
// 		leftCheck = rcGrade.LeftValue < value
// 	case "<=":
// 		leftCheck = rcGrade.LeftValue <= value
// 	case ">":
// 		leftCheck = rcGrade.LeftValue > value
// 	case ">=":
// 		leftCheck = rcGrade.LeftValue >= value
// 	default:
// 		return 0, errors.New("Invalid Left Condition")
// 	}

// 	switch rcGrade.RightCondition {
// 	case "<":
// 		rightCheck = value < rcGrade.RightValue
// 	case "<=":
// 		rightCheck = value <= rcGrade.RightValue
// 	case ">":
// 		rightCheck = value > rcGrade.RightValue
// 	case ">=":
// 		rightCheck = value >= rcGrade.RightValue
// 	default:
// 		return 0, errors.New("Invalid Right Condition")
// 	}

// 	if leftCheck && rightCheck {
// 		return rcGrade.GradeId, nil
// 	}

// 	return 0, nil

// }

func (t *roadConditionUseCase) GetRoadCondition(c *gin.Context, uid uint, idParent int) (responses.RoadCondition, error) {

	roadCondition, err := t.roadConditionRepo.GetRoadConditionByIdParent(idParent)
	if err != nil {
		return responses.RoadCondition{}, err
	}

	road, err := t.roadConditionRepo.GetRoadByID(roadCondition.RoadId)
	if err != nil {
		return responses.RoadCondition{}, err
	}

	STORAGE_IP := os.Getenv("STORAGE_IP") + "/"
	var ImgFilePath string
	if roadCondition.ImgFilePath != "" {
		ImgFilePath = STORAGE_IP + roadCondition.ImgFilePath
	}

	resp := responses.RoadCondition{
		Id:           roadCondition.ID,
		IDParent:     idParent,
		LaneNo:       roadCondition.LaneNo,
		SurveyedDate: roadCondition.SurveyedDate,
		Remarks:      roadCondition.Remarks,
		IriFilename:  STORAGE_IP + roadCondition.IRIInputFilePath,
		ImgFilepath:  ImgFilePath,
		Direction: models.RefDirection{
			ID:   road.RoadInfo.RefDirectionId,
			Name: road.RoadInfo.Direction.Name,
		},
	}

	return resp, nil
}

func (t *roadConditionUseCase) DeleteRoadCondition(c *gin.Context, uid uint, idParent int) (bool, error) {

	roadCondition, err := t.roadConditionRepo.GetRoadConditionByIdParent(idParent)
	if err != nil {
		return false, err
	}

	switch roadCondition.Status {
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

			if roadCondition.Revision > 0 {

				before, err := t.roadConditionRepo.GetRoadConditionBeforeLastRevitionByIdParent(idParent, roadCondition.Revision-1)
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						return true, nil
					}
					return false, err
				}
				if before.Status == "I" {

					roadCondition.Status = "D"
					roadCondition.UpdatedDate = time.Now()
					roadCondition.UpdatedBy = int(uid)
					err = t.roadConditionRepo.DeleteRoadCondition(roadCondition)
					if err != nil {
						return false, err
					}

					before.Status = "A"
					before.UpdatedDate = time.Now()
					before.UpdatedBy = int(uid)
					err = t.roadConditionRepo.DeleteRoadCondition(before)
					if err != nil {
						return false, err
					}

				}

			} else {
				roadCondition.Status = "D"
				roadCondition.UpdatedDate = time.Now()
				roadCondition.UpdatedBy = int(uid)
				err = t.roadConditionRepo.DeleteRoadCondition(roadCondition)
				if err != nil {
					return false, err
				}
			}

		}
	}

	return true, nil
}

func (t *roadConditionUseCase) GetRoadConditionCompareLane(roadId int, rcCompare requests.RoadConditionCompare) (interface{}, error) {

	years := strings.Split(rcCompare.Years, ",")
	lanes := strings.Split(rcCompare.Lanes, ",")
	conditionType := strings.ToLower(rcCompare.ConditionType)
	yearInts := []int{}
	for _, year := range years {
		yearInt, err := strconv.Atoi(year)
		if err == nil {
			yearInts = append(yearInts, yearInt)
		}
	}

	laneInts := []int{}
	for _, lane := range lanes {
		laneInt, err := strconv.Atoi(lane)
		if err == nil {
			laneInts = append(laneInts, laneInt)
		}
	}

	road, err := t.roadConditionRepo.GetRoadByID(roadId)
	if err != nil {
		return nil, err
	}
	if !road.IsActive {
		return []responses.RoadConditionLane{}, nil
	}

	roadConditions, err := t.roadConditionRepo.GetRoadConditionCompare(roadId,
		models.RoadConditionCompareLane{
			Years: yearInts,
			Lanes: laneInts,
		})
	if err != nil {
		return nil, err
	}

	processedYearLanes := make(map[string]bool)
	datas := make(map[int][]responses.RoadConditionLaneItem)

	for _, rc := range roadConditions {

		yearLanesKey := fmt.Sprintf("%d-%d", rc.Year, rc.LaneNo)
		if processedYearLanes[yearLanesKey] {
			continue
		}

		for _, rcKm := range rc.RoadConditionSurveys {
			for _, rc100m := range rcKm.RoadConditionSurvey100Ms {
				for _, rcM := range rc100m.RoadConditionSurveyMs {
					var value *float64

					switch conditionType {
					case "rut":
						value = rcM.RUT
					case "iri":
						value = rcM.IRI
					case "mpd":
						value = rcM.MPD
					case "ifi":
						value = rcM.IFI
					}

					datas[rc.Year] = append(datas[rc.Year], responses.RoadConditionLaneItem{
						LaneNo:  rc.LaneNo,
						KmStart: int(rcM.KmStart),
						KmEnd:   int(rcM.KmEnd),
						Value:   value,
					})

				}
			}
		}

		processedYearLanes[yearLanesKey] = true

	}

	var resps []responses.RoadConditionLane
	for year, items := range datas {
		roadConditionLane := responses.RoadConditionLane{
			Year:  year,
			Items: items,
		}

		// Sort items by LaneNo and KmStart in ascending order
		sort.Slice(roadConditionLane.Items, func(i, j int) bool {
			// First, compare by LaneNo
			if roadConditionLane.Items[i].LaneNo != roadConditionLane.Items[j].LaneNo {
				return roadConditionLane.Items[i].LaneNo < roadConditionLane.Items[j].LaneNo
			}

			// If LaneNo is the same, compare by KmStart
			if road.RoadInfo.RefDirectionId == 1 {
				return roadConditionLane.Items[i].KmStart < roadConditionLane.Items[j].KmStart
			} else {
				return roadConditionLane.Items[i].KmStart > roadConditionLane.Items[j].KmStart
			}

		})

		resps = append(resps, roadConditionLane)
	}
	if len(resps) == 0 {
		return []responses.RoadConditionLane{}, nil
	}
	return resps, nil
}

func (t *roadConditionUseCase) GetRoadConditionCompareYear(roadId int, rcCompare requests.RoadConditionCompare) (interface{}, error) {

	years := strings.Split(rcCompare.Years, ",")
	lanes := strings.Split(rcCompare.Lanes, ",")
	conditionType := strings.ToLower(rcCompare.ConditionType)
	yearInts := []int{}
	for _, year := range years {
		yearInt, err := strconv.Atoi(year)
		if err == nil {
			yearInts = append(yearInts, yearInt)
		}
	}

	laneInts := []int{}
	for _, lane := range lanes {
		laneInt, err := strconv.Atoi(lane)
		if err == nil {
			laneInts = append(laneInts, laneInt)
		}
	}

	road, err := t.roadConditionRepo.GetRoadById(roadId)
	if err != nil {
		return nil, err
	}

	roadInfo, err := t.roadConditionRepo.GetRoadByID(roadId)
	if err != nil {
		return responses.RoadCondition{}, err
	}

	if !road.IsActive {
		return []responses.RoadConditionYear{}, nil
	}

	roadConditions, err := t.roadConditionRepo.GetRoadConditionCompare(roadId,
		models.RoadConditionCompareLane{
			Years: yearInts,
			Lanes: laneInts,
		})
	if err != nil {
		return []responses.RoadConditionYear{}, err
	}

	var value *float64

	processedLaneYears := make(map[string]bool)

	var YearItem responses.RoadConditionYearItem
	datas := make(map[int][]responses.RoadConditionYearItem)

	for _, rc := range roadConditions {
		laneYearKey := fmt.Sprintf("%d-%d", rc.LaneNo, rc.Year)

		if processedLaneYears[laneYearKey] {
			continue
		}

		for _, rcKm := range rc.RoadConditionSurveys {

			for _, rc100m := range rcKm.RoadConditionSurvey100Ms {

				for _, rcM := range rc100m.RoadConditionSurveyMs {
					switch conditionType {
					case "rut":
						value = rcM.RUT
					case "iri":
						value = rcM.IRI
					case "mpd":
						value = rcM.MPD
					case "ifi":
						value = rcM.IFI
					}
					YearItem = responses.RoadConditionYearItem{
						Lane:    rc.LaneNo,
						Year:    rc.Year,
						KmStart: int(rcM.KmStart),
						KmEnd:   int(rcM.KmEnd),
						Value:   value,
					}

					datas[rc.LaneNo] = append(datas[rc.LaneNo], YearItem)
				}

			}

		}

		processedLaneYears[laneYearKey] = true
	}

	var resps []responses.RoadConditionYear
	for lane, items := range datas {
		year := responses.RoadConditionYear{
			Lane:  lane,
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
		return resps[i].Lane < resps[j].Lane
	})

	if len(resps) == 0 {
		return []responses.RoadConditionYear{}, nil
	}
	return resps, nil

}

func (t *roadConditionUseCase) GetRoadConditionCompareAverage(roadId int, lane int) (interface{}, error) {

	road, err := t.roadConditionRepo.GetRoadById(roadId)
	if !road.IsActive {
		return responses.NoData{}, nil
	}
	if err != nil {
		return nil, err
	}

	roadConditions, err := t.roadConditionRepo.GetRoadConditionCompareAverage(roadId, lane)
	if err != nil {
		return nil, err
	}

	dataMap := make(map[int]models.RoadConditionAverage)

	for _, rc := range roadConditions {

		data, ok := dataMap[rc.LaneNo]
		if !ok {
			data = models.RoadConditionAverage{
				Lane:     rc.LaneNo,
				Items:    []models.RoadConditionAverageItem{},
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
			data.Items = append(data.Items, models.RoadConditionAverageItem{
				Year:    rc.Year,
				KmStart: int(rc.KmStart),
				KmEnd:   int(rc.KmEnd),
				IRI:     rc.IRI,
				MPD:     rc.MPD,
				RUT:     rc.RUT,
				IFI:     rc.IFI,
			})
		}

		dataMap[rc.LaneNo] = data
	}
	dataList := responses.RoadConditionAverage{}
	for _, data := range dataMap {
		cleanData := responses.RoadConditionAverage{
			Lane:  data.Lane,
			Items: data.Items,
		}
		dataList = cleanData
	}
	if dataList.Items == nil {
		return responses.NoData{}, nil
	}

	return dataList, nil
}
