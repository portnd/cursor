package usecases

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/tidwall/gjson"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
)

func (u *useCase) GetDataMart(userID int) (interface{}, error) {
	roadIDs, err := u.repo.GetRoadIDAll()
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	// roadIDs := []int{145}
	// return roadIDs, nil
	// lanes := []int{1}
	lanes := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var roads []responses.Road
	var roadDateBegins []responses.RoadDateBegins
	var roadGeoms []responses.RoadGeoms
	roadSurfaces := make(map[string][]models.RoadSurfacePrepareData)
	maintenances := make(map[string][]responses.MaintenanceData)
	roadTotal := len(roadIDs)
	roadCount := 1
	percent := 0.0
	for _, roadID := range roadIDs {
		for _, lane := range lanes {
			laneNo := lane

			roadGeomData, err := u.GetRoadGeomLaneByID(roadID, laneNo)
			if err != nil {
				continue
			}

			roadInfo, err := u.GetRoadByID(roadID)
			if err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, err.Error())
			}

			roadData, err := u.GetRoadGroupByRoadID(roadID)
			if err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, err.Error())
			}

			roadGroup, err := u.GetRoadGroupByID(roadData.RoadGroupId)
			if err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, err.Error())
			}

			roadSurface, err := u.GetRoadSurfaceByRoadID(roadID, laneNo)
			if err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, err.Error())
			}

			if len(roadSurface) == 0 {
				continue
			}

			if len(roadSurface[0].RoadSurfaceLane) == 0 {
				continue
			}
			index := fmt.Sprintf("%d-%d", roadID, laneNo)
			roadSurfaces[index] = roadSurface
			maintenanceRes, err := u.GetMaintenanceData(roadID, laneNo)
			if err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, err.Error())
			}
			// return maintenanceRes, nil
			var maintenanceDatas []responses.MaintenanceData
			for _, item := range maintenanceRes.([]models.MaintenanceData) {
				for _, road := range item.MaintenanceRoads {
					var maintenanceData responses.MaintenanceData
					maintenanceData.RoadGroupID = road.RoadGroupID
					maintenanceData.RoadID = road.RoadID
					maintenanceData.KmStart = road.KmStart
					maintenanceData.KmEnd = road.KmEnd
					maintenanceData.LaneNo = lane
					maintenanceData.Year = item.BudgetYear
					maintenanceData.ProjectEndDate = item.ProjectEndDate
					criteriaMethod, err := u.GetRefCriteriaMethodByID(road.MaintenanceMethodID)
					if err != nil {
						logs.Error(err)
						return "", responses.NewAppErr(400, err.Error())
					}
					surfaceType := ""
					if criteriaMethod.Surface == "cc" {
						surfaceType = "Concrete"
					} else {
						surfaceType = "Asphalt"
					}

					maintenanceData.IcResult.Type = surfaceType
					maintenanceData.IcResult.Method = criteriaMethod.Name

					params := road.RefSurface.Params
					var data map[string]interface{}
					err = json.Unmarshal([]byte(params), &data)
					if err != nil {
						logs.Error(err)
						return "", responses.NewAppErr(400, err.Error())
					}

					maintenanceData.IcResult.RefSurface.ID = int(data["id"].(float64))
					maintenanceData.IcResult.RefSurface.Name = data["name"].(string)
					maintenanceData.IcResult.RefSurface.Type = data["type"].(string)
					maintenanceData.IcResult.RefSurface.SurfaceGroup = data["surface_group"].(string)

					interventionCriteriaParamsID := road.InterventionCriteriaIDParams
					itvCriteriaParam, err := u.GetInterventionCriteriaParamsById(interventionCriteriaParamsID)
					if err != nil {
						logs.Error(err)
						return "", responses.NewAppErr(400, err.Error())
					}

					value1 := gjson.Get(itvCriteriaParam.Params, strings.ToLower(surfaceType)).String()
					value2 := gjson.Get(value1, criteriaMethod.Name).String()
					var data2 []map[string]interface{}
					err = json.Unmarshal([]byte(value2), &data2)
					if err != nil {
						logs.Error(err)
						// return "", responses.NewAppErr(400, err.Error())
					}

					for _, item := range data2 {
						if int(item["id"].(float64)) == road.InterventionCriteriaID {
							maintenanceData.IcResult.ThicknessRepair = item["maintenance_thickness"].(float64)
							maintenanceData.IcResult.ThicknessScrape = item["maintenance_scraping"].(float64)
						}
					}
					maintenanceDatas = append(maintenanceDatas, maintenanceData)
				}
			}
			key := fmt.Sprintf("%d-%d", roadID, laneNo)
			if maintenanceDatas != nil {
				maintenances[key] = maintenanceDatas
			}
			road := responses.Road{RoadID: roadID, RoadGroupID: roadData.RoadGroupId, RefDirectionID: roadInfo.RefDirectionId, RoadName: roadInfo.Name, RoadGroupName: roadGroup.ShortName, YearConstructionCompleted: roadInfo.YearConstructionCompleted}
			roadGeom := responses.RoadGeoms{RoadID: roadID, KmStart: float64(roadGeomData.KmStart), KmEnd: float64(roadGeomData.KmEnd), LaneNo: roadGeomData.LaneNo}

			roads = append(roads, road)
			// for _, item := range roadDateBeginData {
			roadDateBegin := responses.RoadDateBegins{RoadID: roadID, KmStart: float64(roadInfo.KmStart), KmEnd: float64(roadInfo.KmEnd), Year: roadInfo.YearConstructionCompleted}
			roadDateBegins = append(roadDateBegins, roadDateBegin)
			// }

			roadGeoms = append(roadGeoms, roadGeom)
		}
		percent = float64(roadCount) / float64(roadTotal) * 100
		if percent > 50 {
			percent = 50.00
		}
		err := u.repo.UpdatePercentage(percent, userID)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}

		roadCount++

	}
	data, err := u.PrepareDataFn(roads, roadDateBegins, roadGeoms, roadSurfaces, maintenances)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	err = u.repo.UpdatePercentage(90.00, userID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	if len(data.([]models.DataMart)) > 0 {
		_, err := u.repo.CreateDataMart(data.([]models.DataMart))
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		err = u.repo.UpdatePercentage(100.00, userID)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
	}
	return data, nil
}

func (u *useCase) PrepareDataFn(roads []responses.Road, roadDateBegins []responses.RoadDateBegins, roadGeom []responses.RoadGeoms, roadSurfaces map[string][]models.RoadSurfacePrepareData, maintenances map[string][]responses.MaintenanceData) (interface{}, error) {
	var intervalKmRanges = make([]map[string]interface{}, 0)
	indexStart := 0
	for road_index, road := range roads {
		kms := make(map[float64]bool)
		var keys []float64
		laneNo := roadGeom[road_index].LaneNo
		indexSurface := fmt.Sprintf("%d-%d", road.RoadID, laneNo)

		if road.RefDirectionID == 1 { // LT
			roadGeom := map[string]float64{
				"km_start": roadGeom[road_index].KmStart,
				"km_end":   roadGeom[road_index].KmEnd,
			}

			kms[roadGeom["km_start"]] = true
			kms[roadGeom["km_end"]] = true

			// inventory ผิวทาง (road surface)
			for _, surface := range roadSurfaces[indexSurface] {
				if (surface.KmStart >= roadGeom["km_start"]) && (surface.KmStart < roadGeom["km_end"]) {
					kms[surface.KmStart] = true
				}

				if (surface.KmEnd <= roadGeom["km_end"]) && (surface.KmEnd > roadGeom["km_start"]) {
					kms[surface.KmEnd] = true
				}
			}
			key := fmt.Sprintf("%d-%d", road.RoadID, laneNo)
			for _, maintenance := range maintenances[key] {
				if (maintenance.KmStart >= roadGeom["km_start"]) && (maintenance.KmStart <= roadGeom["km_end"]) {
					kms[maintenance.KmStart] = true
				}

				if (maintenance.KmEnd <= roadGeom["km_end"]) && (maintenance.KmEnd >= roadGeom["km_start"]) {
					kms[maintenance.KmEnd] = true
				}
			}

			// ประวัติการซ่อมบำรุง
			for key := range kms {
				keys = append(keys, key)
			}
			sort.Float64s(keys)

		} else { // RT
			roadGeom := map[string]float64{
				"km_start": roadGeom[road_index].KmStart,
				"km_end":   roadGeom[road_index].KmEnd,
			}

			kms[roadGeom["km_start"]] = true
			kms[roadGeom["km_end"]] = true

			// inventory ผิวทาง (road surface)
			for _, surface := range roadSurfaces[indexSurface] {
				if (surface.KmStart <= roadGeom["km_start"]) && (surface.KmStart > roadGeom["km_end"]) {
					kms[surface.KmStart] = true
				}

				if (surface.KmEnd >= roadGeom["km_end"]) && (surface.KmEnd < roadGeom["km_start"]) {
					kms[surface.KmEnd] = true
				}
			}
			key := fmt.Sprintf("%d-%d", road.RoadID, laneNo)
			for _, maintenance := range maintenances[key] {
				if (maintenance.KmStart <= roadGeom["km_start"]) && (maintenance.KmStart >= roadGeom["km_end"]) {
					kms[maintenance.KmStart] = true
				}

				if (maintenance.KmEnd >= roadGeom["km_end"]) && (maintenance.KmEnd <= roadGeom["km_start"]) {
					kms[maintenance.KmEnd] = true
				}
			}

			// ประวัติการซ่อมบำรุง
			for key := range kms {
				keys = append(keys, key)
			}
			sort.Float64s(keys)
			// เรียงจากมากไปน้อย
			for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
		kmRanges := make([]map[string]float64, 0)
		/* แบ่งช่วงตามจุดตัด */

		for index := 0; index < len(keys)-1; index++ {
			kmRange := make(map[string]float64)
			kmRange["km_start"] = keys[index]
			kmRange["km_end"] = keys[index+1]
			kmRanges = append(kmRanges, kmRange)
		}

		/* แบ่งช่วง*/
		// interval := 100 * 1000 // 1 kilometer = 1000 meters
		for _, km := range kmRanges {
			kmStart := km["km_start"]
			kmEnd := km["km_end"]
			intervalKmRanges = append(intervalKmRanges, map[string]interface{}{
				"km_start": kmStart,
				"km_end":   kmEnd,
			})
		}

		/* Get ข้อมูลวันที่เริ่มใช้สายทาง, ข้อมูลผิวทาง (road_surface) และประวัติการซ่อมบำรุง (Maintenance) */
		i := 0
		for indexR0 := range intervalKmRanges {
			if indexR0 >= indexStart {
				index := i + indexStart
				intervalKmRanges[index]["road"] = road
				intervalKmRanges[index]["road_geom"] = roadGeom[road_index]
				intervalKmRanges[index]["road_date_begin"] = responses.RoadDateBegins{}
				intervalKmRanges[index]["surface"] = models.RoadSurfacePrepareData{}
				intervalKmRanges[index]["maintenances"] = []interface{}{}
				i++
			}

		}
		j := 0
		// return intervalKmRanges, nil
		for indexR1, kmRange := range intervalKmRanges {
			laneNo := kmRange["road_geom"].(responses.RoadGeoms).LaneNo
			if indexR1 >= indexStart {
				index := j + indexStart
				// Find matching road date begins
				for _, dateBegin := range roadDateBegins {
					if road.RefDirectionID == 1 {
						if (dateBegin.KmStart >= kmRange["km_start"].(float64)) || (kmRange["km_start"].(float64) < dateBegin.KmEnd) { // เงื่อนไข LT
							if dateBegin.KmStart <= kmRange["km_end"].(float64) {
								intervalKmRanges[index]["road_date_begin"] = dateBegin // มีข้อมูลชุดเดียว
							}
						}
					} else {
						if (dateBegin.KmStart <= kmRange["km_start"].(float64)) || (kmRange["km_start"].(float64) > dateBegin.KmEnd) { // เงื่อนไข LT
							if dateBegin.KmStart >= kmRange["km_end"].(float64) {
								intervalKmRanges[index]["road_date_begin"] = dateBegin // มีข้อมูลชุดเดียว
							}
						}
					}

				}
				// return roadSurfaces[road.RoadID], nil
				// Find matching road surfaces
				for _, surface := range roadSurfaces[indexSurface] {
					if road.RefDirectionID == 1 {
						if (surface.KmStart > kmRange["km_start"].(float64)) || (kmRange["km_start"].(float64) < surface.KmEnd) { // เงื่อนไข LT
							if surface.KmStart < kmRange["km_end"].(float64) {
								intervalKmRanges[index]["surface"] = surface // มีข้อมูลชุดเดียว
							}
						}
					} else {
						if (surface.KmStart < kmRange["km_start"].(float64)) || (kmRange["km_start"].(float64) > surface.KmEnd) { // เงื่อนไข LT
							if surface.KmStart >= kmRange["km_end"].(float64) {
								intervalKmRanges[index]["surface"] = surface // มีข้อมูลชุดเดียว
							}
						}
					}

				}

				// // Find matching road maintenances
				var maintenanceArr []responses.MaintenanceData
				key := fmt.Sprintf("%d-%d", road.RoadID, laneNo)
				for _, maintenance := range maintenances[key] {
					if maintenance.LaneNo != laneNo {
						continue
					}
					if road.RefDirectionID == 1 {
						if (maintenance.KmStart >= kmRange["km_start"].(float64)) || (kmRange["km_start"].(float64) < maintenance.KmEnd) { // เงื่อนไข LT
							if maintenance.KmStart < kmRange["km_end"].(float64) {
								maintenanceArr = append(maintenanceArr, maintenance)
							}
						}
					} else {
						if (maintenance.KmStart <= kmRange["km_start"].(float64)) || (kmRange["km_start"].(float64) > maintenance.KmEnd) { // เงื่อนไข LT
							if maintenance.KmStart > kmRange["km_end"].(float64) {
								maintenanceArr = append(maintenanceArr, maintenance)
							}
						}
					}

				}
				if len(maintenanceArr) == 0 {
					intervalKmRanges[index]["maintenances"] = []responses.MaintenanceData{} // มีข้อมูลหลายชุด
				} else {
					intervalKmRanges[index]["maintenances"] = maintenanceArr // มีข้อมูลหลายชุด
				}
				j++
			}
		}
		indexStart = len(intervalKmRanges)
	}

	roadSurfaceData, err := u.repo.GetRoadSurfaceAll()
	if err != nil {
		logs.Error(err)
		return []models.DataMart{}, responses.NewAppErr(400, err.Error())
	}
	var dataMarts []models.DataMart

	for _, road := range intervalKmRanges {
		result := make(map[string]interface{})
		result = road
		copier.Copy(&result, &road)
		/* road length ระยะทางของสายทาง */
		roadGeom := result["road_geom"].(responses.RoadGeoms)
		// roadGeom := result["road_geom"].(RoadGeoms)
		result["lane_length"] = math.Abs(roadGeom.KmStart - roadGeom.KmEnd)
		/* length ระยะทางของช่วงกมที่วิเคราะห์​ */
		result["length"] = math.Abs(result["km_start"].(float64) - result["km_end"].(float64))
		/* Area ของช่วงกมที่วิเคราะห์ ​*/
		laneWidth := 3.5 // ค่าคงที่
		result["lane_width"] = laneWidth
		result["area"] = result["length"].(float64) * laneWidth
		/* สภาพทาง */
		laneNo := roadGeom.LaneNo
		roadID := result["road"].(responses.Road).RoadID
		directionID := result["road"].(responses.Road).RefDirectionID
		roadConditions, surveyedDate, _ := u.GetRoadCondition(roadID, laneNo) // get Road Condition

		roadConditionDataRes := responses.RoadConditionPrepareDataRes{}
		if len(roadConditions) != 0 {
			roadConditionData, _ := u.GetRoadConditionPrepareData(directionID, laneNo, roadConditions, result)
			copier.Copy(&roadConditionDataRes, roadConditionData)
			roadConditionDataRes.SurveyDate = surveyedDate
		} else { //ปรับแก้ case มากกว่า 3 เลน
			// lane 1
			if laneNo == 1 {
				roadConditionLane2, surveyedDate, _ := u.GetRoadCondition(roadID, 2)

				if len(roadConditionLane2) > 0 {
					roadConditionData, _ := u.GetRoadConditionPrepareData(directionID, 2, roadConditionLane2, result)
					copier.Copy(&roadConditionData, roadConditionData)
					roadConditionDataRes.RUT = roadConditionData.RUT
					roadConditionDataRes.IRI = roadConditionData.IRI
					roadConditionDataRes.IFI = roadConditionData.IFI
					roadConditionDataRes.SurveyDate = surveyedDate
				} else {
					roadConditionDataRes.RUT = 0
					roadConditionDataRes.IRI = 0
					roadConditionDataRes.IFI = 0

				}
			} else if laneNo >= 10 {
				// lane 10
				roadConditionLane6, surveyedDate, _ := u.GetRoadCondition(roadID, laneNo-1)

				if len(roadConditionLane6) > 0 {
					roadConditionData, _ := u.GetRoadConditionPrepareData(directionID, laneNo-1, roadConditionLane6, result)
					copier.Copy(&roadConditionData, roadConditionData)
					roadConditionDataRes.RUT = roadConditionData.RUT
					roadConditionDataRes.IRI = roadConditionData.IRI
					roadConditionDataRes.IFI = roadConditionData.IFI
					roadConditionDataRes.SurveyDate = surveyedDate
				} else {
					roadConditionDataRes.RUT = 0
					roadConditionDataRes.IRI = 0
					roadConditionDataRes.IFI = 0
				}
			} else { //เลน 2-9
				// lane n -1
				roadConditionDataRes1 := responses.RoadConditionPrepareDataRes{}
				roadConditionLane1, surveyedDate, _ := u.GetRoadCondition(roadID, laneNo-1)
				if len(roadConditionLane1) > 0 {
					roadConditionData1, _ := u.GetRoadConditionPrepareData(directionID, laneNo-1, roadConditionLane1, result)
					copier.Copy(&roadConditionDataRes1, roadConditionData1)
					roadConditionDataRes1.SurveyDate = surveyedDate
				} else {
					roadConditionDataRes1.RUT = 0
					roadConditionDataRes1.IRI = 0
					roadConditionDataRes1.IFI = 0
					// roadConditionDataRes1.SurveyDate = nil
				}

				// lane n + 1
				roadConditionDataRes3 := responses.RoadConditionPrepareDataRes{}
				roadConditionLane3, surveyedDate, _ := u.GetRoadCondition(roadID, laneNo+1)

				if len(roadConditionLane3) > 0 {
					roadConditionData3, _ := u.GetRoadConditionPrepareData(directionID, laneNo+1, roadConditionLane3, result)
					copier.Copy(&roadConditionDataRes3, roadConditionData3)
					roadConditionDataRes3.SurveyDate = surveyedDate
				} else {
					roadConditionDataRes3.RUT = 0
					roadConditionDataRes3.IRI = 0
					roadConditionDataRes3.IFI = 0
				}

				if roadConditionDataRes3.IRI > roadConditionDataRes1.IRI {
					roadConditionDataRes.RUT = roadConditionDataRes3.RUT
					roadConditionDataRes.IRI = roadConditionDataRes3.IRI
					roadConditionDataRes.IFI = roadConditionDataRes3.IFI
					roadConditionDataRes.SurveyDate = surveyedDate
				} else {
					roadConditionDataRes.RUT = roadConditionDataRes1.RUT
					roadConditionDataRes.IRI = roadConditionDataRes1.IRI
					roadConditionDataRes.IFI = roadConditionDataRes1.IFI
					roadConditionDataRes.SurveyDate = surveyedDate
				}
			}
		}

		result["road_condition"] = roadConditionDataRes

		lastType, curType := u.GetType(road, roadConditionDataRes.SurveyDate)
		result["last_type"] = lastType
		result["type"] = curType
		currentSurfaceID, _ := u.GetCurrentSurfaceID(result, roadConditionDataRes.SurveyDate)
		// currentSurface, _ := u.GetCurrentSurface(result, roadConditionDataRes.SurveyDate)
		// result["current_surface"] = currentSurface
		if roadConditionDataRes.SurveyDate == nil {
			result["analyst_year"] = 0 //roadConditionDataRes.SurveyDate.(time.Time).Year() //year + (t + 1) // ปีที่วิเคราะห์
		} else {
			result["analyst_year"] = roadConditionDataRes.SurveyDate.(time.Time).Year() //year + (t + 1) // ปีที่วิเคราะห์
		}

		if helpers.InArrayInt(result["type"].(int), []int{1, 2, 3, 4}) {
			/* หา Age */
			age, _ := u.GetAgeAC(result, roadConditionDataRes.SurveyDate)
			result["age"] = age
		} else {
			age, _ := u.GetAgeCC(result, roadConditionDataRes.SurveyDate)
			result["age"] = age
			result["hsold_hsnew"] = responses.HsoldHsnewPrepareDataRes{}
		}
		age := road["age"].(responses.AgePrepareDataRes)
		// currentSurfaceRes := road["current_surface"].(responses.RefSurface)
		theGeom, err := u.GenTheGeom(roadID, laneNo, directionID, road["km_start"].(float64), road["km_end"].(float64))
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		roadSurfaceID := 0
		laneCount := 0
		surfaceYear := 0
		for _, item := range roadSurfaceData[roadID] {
			if directionID == 1 {
				if road["km_start"].(float64) >= item.KmStart && road["km_end"].(float64) <= item.KmEnd {
					roadSurfaceID = item.Id
					laneCount = len(item.RoadSurfaceLane)
					surfaceYear = item.Year
					goto EXITLOOP
				}
			} else {
				if road["km_start"].(float64) <= item.KmStart && road["km_end"].(float64) >= item.KmEnd {
					roadSurfaceID = item.Id
					laneCount = len(item.RoadSurfaceLane)
					surfaceYear = item.Year
					goto EXITLOOP
				}
			}

		}
	EXITLOOP:
		var dataMart models.DataMart
		dataMart.TheGeom = theGeom.(string)
		dataMart.RoadSurfaceID = roadSurfaceID
		dataMart.LaneCount = laneCount
		dataMart.RoadID = roadID
		dataMart.LaneNo = laneNo
		dataMart.KmStart = road["km_start"].(float64)
		dataMart.KmEnd = road["km_end"].(float64)
		dataMart.Age = age.Age
		dataMart.SurfaceYear = surfaceYear
		dataMart.RefSurfaceID = currentSurfaceID //currentSurfaceRes.ID
		if len(road["maintenances"].([]responses.MaintenanceData)) > 0 {
			cnt := len(road["maintenances"].([]responses.MaintenanceData))
			dataMart.Year = road["maintenances"].([]responses.MaintenanceData)[cnt-1].BudgetYear
			dataMart.ContractNumber = road["maintenances"].([]responses.MaintenanceData)[cnt-1].ContractNumber
			lastInspectionDate := road["maintenances"].([]responses.MaintenanceData)[cnt-1].ProjectEndDate.Format("2006-01-02")
			dataMart.LastInspectionDate = &lastInspectionDate
			// helpers.PrintlnJson("ddd", road["maintenances"].([]responses.MaintenanceData)[0].LastInspectionDate)
		} else {
			dataMart.LastInspectionDate = nil
		}
		if roadSurfaceID != 0 {
			dataMarts = append(dataMarts, dataMart)

		}
	}
	return dataMarts, nil
}

func (u *useCase) GenTheGeom(roadID, laneNo, directionID int, kmStart, kmEnd float64) (interface{}, error) {
	roadGeom, err := u.repo.GetRoadGeomByRoadIDLaneNo(roadID, laneNo)
	if err != nil {
		return "", err
	}

	var subRoadStart float64
	var subRoadEnd float64
	var theGeom string
	switch directionID {
	case 1:
		if kmStart >= roadGeom.KmStart {
			if kmStart >= (roadGeom.KmEnd) {
				subRoadStart = 0
			} else {
				subRoadStart = (math.Abs((kmStart - (roadGeom.KmStart)))) / (math.Abs(((roadGeom.KmStart) - (roadGeom.KmEnd))))
			}
		} else {
			subRoadStart = 0
		}
		if kmEnd < roadGeom.KmEnd {
			if kmEnd < roadGeom.KmStart {
				subRoadEnd = 0
			} else {
				// subRoadEnd = float64(math.Abs(float64(kmEnd-float64(roadGeom.KmStart)))) / float64(math.Abs(float64(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))))
				subRoadEnd = math.Abs(kmEnd-roadGeom.KmStart) / math.Abs(roadGeom.KmStart-roadGeom.KmEnd)
			}
		} else {
			subRoadEnd = 1
		}

	case 2:
		if kmStart <= roadGeom.KmStart {
			if kmStart < float64(roadGeom.KmEnd) {
				subRoadStart = 0
			} else {
				// subRoadStart = float64(math.Abs(float64(kmStart-float64(roadGeom.KmStart)))) / float64(math.Abs(float64(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))))
				subRoadStart = math.Abs(kmStart-roadGeom.KmStart) / math.Abs(roadGeom.KmStart-roadGeom.KmEnd)
			}
		} else {
			subRoadStart = 1
		}

		if kmEnd >= roadGeom.KmEnd {
			if kmEnd >= roadGeom.KmStart {
				subRoadEnd = 0
			} else {
				// subRoadEnd = float64(math.Abs(float64(kmEnd-float64(roadGeom.KmStart)))) / float64(math.Abs(float64(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))))
				subRoadEnd = math.Abs(kmEnd-roadGeom.KmStart) / math.Abs(roadGeom.KmStart-roadGeom.KmEnd)
			}
		} else {
			subRoadEnd = 0
		}
	}

	subRoadMin := math.Min(subRoadStart, subRoadEnd)
	subRoadMax := math.Max(subRoadStart, subRoadEnd)
	theGeom = fmt.Sprintf("ST_LineSubstring('%s', %f, %f)", roadGeom.TheGeom, subRoadMin, subRoadMax)
	return theGeom, nil
}

type Datasss struct {
	Surface      models.RoadSurfacePrepareData `json:"surface"`
	Maintenances interface{}                   `json:"maintenances"`
}

func (u *useCase) GetData(roads []models.RoadById, roadSurfaces map[int][]models.RoadSurfacePrepareData, maintenances map[int][]responses.MaintenanceData) (interface{}, error) {
	rangeData := make(map[string][]map[string]float64) // []map[string]float64
	indexStart := 0
	for _, road := range roads {
		intervalKmRanges := make([]map[string]interface{}, 0)
		roadID := road.Id
		for _, geom := range road.RoadGeom {
			str := fmt.Sprintf("%d-%d", geom.RoadId, geom.LaneNo)
			kms := make(map[float64]bool)
			var keys []float64
			roadGeom := map[string]float64{
				"km_start": geom.KmStart,
				"km_end":   geom.KmEnd,
			}

			kms[roadGeom["km_start"]] = true
			kms[roadGeom["km_end"]] = true
			// if road.RefDirectionId == 1 { // LT

			// inventory ผิวทาง (road surface)
			for _, surface := range roadSurfaces[roadID] {
				if (surface.KmStart >= roadGeom["km_start"]) && (surface.KmStart < roadGeom["km_end"]) {
					kms[surface.KmStart] = true
				}

				if (surface.KmEnd <= roadGeom["km_end"]) && (surface.KmEnd > roadGeom["km_start"]) {
					kms[surface.KmEnd] = true
				}
			}

			for _, maintenance := range maintenances[roadID] {
				if (maintenance.KmStart >= roadGeom["km_start"]) && (maintenance.KmStart <= roadGeom["km_end"]) {
					kms[maintenance.KmStart] = true
				}

				if (maintenance.KmEnd <= roadGeom["km_end"]) && (maintenance.KmEnd >= roadGeom["km_start"]) {
					kms[maintenance.KmEnd] = true
				}
			}
			// ประวัติการซ่อมบำรุง
			for key := range kms {
				keys = append(keys, key)
			}
			sort.Float64s(keys)
			kmRanges := make([]map[string]float64, 0)
			/* แบ่งช่วงตามจุดตัด */

			for index := 0; index < len(keys)-1; index++ {
				kmRange := make(map[string]float64)
				kmRange["km_start"] = keys[index]
				kmRange["km_end"] = keys[index+1]
				kmRanges = append(kmRanges, kmRange)
			}

			rangeData[str] = kmRanges
			indexStart = len(intervalKmRanges)
			fmt.Println("indexStart", indexStart)
		}
	}

	var Datassssss []Datasss
	for key, ranges := range rangeData {
		var data Datasss
		keyData := strings.Split(key, "-")
		roadID, _ := strconv.Atoi(keyData[0])
		for _, kmRange := range ranges {
			for _, surface := range roadSurfaces[roadID] {
				if (surface.KmStart > kmRange["km_start"]) || (kmRange["km_start"] < surface.KmEnd) { // เงื่อนไข LT
					if surface.KmStart < kmRange["km_end"] {
						data.Surface = surface
					}
				}
			}
		}
		var maintenanceArr []responses.MaintenanceData
		for _, maintenance := range maintenances[roadID] {
			for _, kmRange := range ranges {
				if (maintenance.KmStart >= kmRange["km_start"]) || (kmRange["km_start"] < maintenance.KmEnd) { // เงื่อนไข LT
					if maintenance.KmStart < kmRange["km_end"] {
						maintenanceArr = append(maintenanceArr, maintenance)
					}
				}
			}
		}
		data.Maintenances = maintenanceArr
		Datassssss = append(Datassssss, data)
	}

	return Datassssss, nil
}

func (u *useCase) GetRoadSurfaceByRoadID(roadID, laneNo int) ([]models.RoadSurfacePrepareData, error) {
	surfaceGroup, _ := u.repo.GetRoadSurfaceByGroupRoadID(roadID, laneNo)
	for _, grp := range surfaceGroup {
		data, _ := u.repo.GetRoadSurfaceByRoadID(roadID, laneNo, grp)
		if len(data) > 0 {
			year := data[0].Year
			_, err := u.repo.GetMaintenanceHistoryByRoadID(roadID, year)
			if err == nil {
				var roadSurfaces []models.RoadSurfacePrepareData
				for _, item := range data {
					var roadSurfaceLanes []models.RoadSurfaceLanePrePareData
					for _, item2 := range item.RoadSurfaceLane {
						if item2.RefSurfaceParamsID == 0 {
							continue
						}
						surfaceParams, _ := u.repo.GetRefSurfaceParam(item2.RefSurfaceParamsID)
						params := surfaceParams.Params
						var data models.RefSurface
						err = json.Unmarshal([]byte(params), &data)
						if err != nil {
							continue
						}
						var roadSurfaceLane models.RoadSurfaceLanePrePareData
						copier.Copy(&roadSurfaceLane, item2)
						roadSurfaceLane.RefSurface = data
						roadSurfaceLanes = append(roadSurfaceLanes, roadSurfaceLane)
					}

					var roadSurface models.RoadSurfacePrepareData
					copier.Copy(&roadSurface, item)
					roadSurface.RoadSurfaceLane = roadSurfaceLanes
					roadSurfaces = append(roadSurfaces, roadSurface)
				}
				return roadSurfaces, nil
			}
		}
	}
	surface, err := u.repo.GetRoadSurfaceFirstGrpByRoadID(roadID, laneNo)
	if err != nil {
		logs.Error(err)
		return []models.RoadSurfacePrepareData{}, err
	}
	var roadSurfaces []models.RoadSurfacePrepareData
	for _, item := range surface {
		var roadSurfaceLanes []models.RoadSurfaceLanePrePareData
		for _, item2 := range item.RoadSurfaceLane {
			if item2.RefSurfaceParamsID == 0 {
				continue
			}
			surfaceParams, _ := u.repo.GetRefSurfaceParam(item2.RefSurfaceParamsID)
			params := surfaceParams.Params
			var data models.RefSurface
			err = json.Unmarshal([]byte(params), &data)
			if err != nil {
				continue
			}
			var roadSurfaceLane models.RoadSurfaceLanePrePareData
			copier.Copy(&roadSurfaceLane, item2)
			roadSurfaceLane.RefSurface = data
			roadSurfaceLanes = append(roadSurfaceLanes, roadSurfaceLane)
		}

		var roadSurface models.RoadSurfacePrepareData
		copier.Copy(&roadSurface, item)
		roadSurface.RoadSurfaceLane = roadSurfaceLanes
		roadSurfaces = append(roadSurfaces, roadSurface)
	}
	return roadSurfaces, nil
}

func (u *useCase) GetMaintenanceData(roadId, laneNo int) (interface{}, error) { //([]models.MaintenanceData, error) {
	data, err := u.repo.GetMaintenanceData(roadId, laneNo)
	if err != nil {
		return data, responses.NewAppErr(400, err.Error())
	}
	var data3 []models.MaintenanceData
	var maintenanceInterventionCriteria models.MaintenanceInterventionCriteria
	for _, item := range data {
		if len(item.MaintenanceRoads) <= 0 {
			continue
		}
		var mrs []models.MaintenanceRoadPrepareData
		for _, road := range item.MaintenanceRoads {
			refCriteriaMethod, err := u.repo.GetRefCriteriaMethodByID(road.MaintenanceMethodID)
			if err != nil {
				logs.Error(err)
				return data, responses.NewAppErr(400, err.Error())
			}
			criteriaMethodName := refCriteriaMethod.Name
			surfaceType := ""
			if refCriteriaMethod.Surface == "cc" {
				surfaceType = "concrete"
			} else {
				surfaceType = "asphalt"
			}
			interventionCriteriaParams, err := u.repo.GetInterventionCriteriaParamsById(road.InterventionCriteriaIDParams)
			if err != nil {
				logs.Error(err)
				return data, responses.NewAppErr(400, err.Error())
			}

			params := interventionCriteriaParams.Params
			var data map[string]interface{}
			err = json.Unmarshal([]byte(params), &data)
			if err != nil {
				logs.Error(err)
				return data, responses.NewAppErr(400, err.Error())
			}

			interventionCriteria := data[surfaceType]

			val, ok := interventionCriteria.(map[string]interface{})
			if !ok {
				// return
			}

			var maintenanceItem []MaintenanceItem
			criteria := val[criteriaMethodName]
			jsonBytes, err := json.Marshal(criteria)
			if err != nil {
				logs.Error(err)
				// fmt.Printf("Error marshaling JSON: %s\n", err.Error())
				// return
			}

			err = json.Unmarshal(jsonBytes, &maintenanceItem)
			if err != nil {
				logs.Error(err)
				fmt.Println("Error:", string(jsonBytes))
				// return
			}

			for _, item := range maintenanceItem {
				if item.ID == road.InterventionCriteriaID {
					copier.Copy(&maintenanceInterventionCriteria, item)
				}
			}
			var mr models.MaintenanceRoadPrepareData
			copier.Copy(&mr, road)
			mr.InterventionCriteria = maintenanceInterventionCriteria
			mrs = append(mrs, mr)
		}

		var data2 models.MaintenanceData
		copier.Copy(&data2, item)
		data2.MaintenanceRoads = mrs
		data3 = append(data3, data2)
	}
	return data3, nil
}

func (u *useCase) GetInterventionCriteriaParamsById(ID int) (models.SettingInterventionCriteriaParams, error) {
	data, err := u.repo.GetInterventionCriteriaParamsById(ID)
	if err != nil {
		return models.SettingInterventionCriteriaParams{}, err
	}
	return data, nil
}

func (u *useCase) GetRefCriteriaMethodByID(ID int) (models.RefCriteriaMethod, error) {
	data, err := u.repo.GetRefCriteriaMethodByID(ID)
	if err != nil {
		return models.RefCriteriaMethod{}, err
	}
	return data, nil
}

type MaintenanceItem struct {
	ID                       int                    `json:"id"`
	MaintenanceCondition     []MaintenanceCondition `json:"maintenance_condition"`
	MaintenanceCostPerUnit   int                    `json:"maintenance_cost_per_unit"`
	MaintenanceDescription   string                 `json:"maintenance_description"`
	MaintenanceMethod        string                 `json:"maintenance_method"`
	MaintenanceScraping      float64                `json:"maintenance_scraping"`
	MaintenanceSequence      int                    `json:"maintenance_sequence"`
	MaintenanceStandardName  string                 `json:"maintenance_standard_name"`
	MaintenanceSurfaceTypeID int                    `json:"maintenance_surface_type_id"`
	MaintenanceThickness     float64                `json:"maintenance_thickness"`
	MaintenanceType          string                 `json:"maintenance_type"`
}

type MaintenanceCondition struct {
	ConditionCriterion  string  `json:"condition_criterion"`
	ConditionLink       string  `json:"condition_link"`
	ConditionOperation1 string  `json:"condition_operation_1"`
	ConditionOperation2 string  `json:"condition_operation_2"`
	ConditionValue1     float64 `json:"condition_value_1"`
	ConditionValue2     int     `json:"condition_value_2"`
	ID                  int     `json:"id"`
}

func (u *useCase) GetRoadByID(roadID int) (models.RoadInfo, error) {
	data, err := u.repo.GetRoadInfoByID(roadID)
	if err != nil {
		logs.Error(err)
		return data, err
	}

	return data, nil
}

func (u *useCase) GetRoadDatebegin(roadID int) ([]models.RoadDatebegin, error) {
	data, err := u.repo.GetRoadDatebegin(roadID)
	if err != nil {
		logs.Error(err)
		return data, err
	}

	return data, nil
}

func (u *useCase) GetRoadGroupByRoadID(roadID int) (models.Road, error) {
	data, err := u.repo.GetRoadByRoadID(roadID)
	if err != nil {
		logs.Error(err)
		return data, err
	}

	return data, nil
}

func (u *useCase) GetRoadGroupByID(ID int) (models.RoadGroup, error) {
	data, err := u.repo.GetRoadGroupByID(ID)
	if err != nil {
		logs.Error(err)
		return data, err
	}

	return data, nil
}

func (u *useCase) GetRoadGeomLaneByID(roadID, laneNo int) (models.RoadGeom, error) {
	data, err := u.repo.GetRoadGeomLaneByID(roadID, laneNo)
	if err != nil {
		logs.Error(err)
		return data, err
	}

	return data, nil
}

func (u *useCase) GetCurrentSurface(road map[string]interface{}, surveyDate interface{}) (interface{}, error) {

	currentSurface := road["surface"].(models.RoadSurfacePrepareData).RoadSurfaceLane
	maintenances := road["maintenances"].([]responses.MaintenanceData)
	if len(maintenances) > 0 {
		currentSurfaceRes := responses.RefSurface{}
		for _, item := range maintenances {

			// if item.LastInspectionDate.Before(surveyDate.(time.Time)) || item.LastInspectionDate.Equal(surveyDate.(time.Time)) {
			if road["type"].(int) != 5 {
				current_surface := item.IcResult.RefSurface // ['ic_result']['ref_surface']; // ผิวใหม่
				currentSurfaceRes = current_surface
			}
			if road["type"].(int) == 5 && item.IcResult.Type == "Concrete" { // กรณีผิวคอนกรีต จะอัปเดตผิวเฉพาะซ่อมด้วยคอนกรีตเท่านั้น
				current_surface := item.IcResult.RefSurface //['ic_result']['ref_surface']; // ผิวใหม่
				currentSurfaceRes = current_surface
			}
			// }

		}
		return currentSurfaceRes, nil
	}
	if len(currentSurface) > 0 {

		return ConverstCurrentSurface(currentSurface[0]), nil
	} else {

		return ConverstCurrentSurface(models.RoadSurfaceLanePrePareData{}), nil
	}
}

func (u *useCase) GetCurrentSurfaceID(road map[string]interface{}, surveyDate interface{}) (int, error) {
	currentSurface := road["surface"].(models.RoadSurfacePrepareData).RoadSurfaceLane
	maintenances := road["maintenances"].([]responses.MaintenanceData)
	if len(maintenances) > 0 {
		// currentSurfaceRes := responses.RefSurface{}
		for _, item := range maintenances {

			// if item.LastInspectionDate.Before(surveyDate.(time.Time)) || item.LastInspectionDate.Equal(surveyDate.(time.Time)) {
			if road["type"].(int) != 5 {
				// current_surface := item.IcResult.RefSurface // ['ic_result']['ref_surface']; // ผิวใหม่
				// currentSurfaceRes = current_surface
				return item.IcResult.RefSurface.ID, nil
			}
			if road["type"].(int) == 5 && item.IcResult.Type == "Concrete" { // กรณีผิวคอนกรีต จะอัปเดตผิวเฉพาะซ่อมด้วยคอนกรีตเท่านั้น
				// current_surface := item.IcResult.RefSurface //['ic_result']['ref_surface']; // ผิวใหม่
				// currentSurfaceRes = current_surface
				return item.IcResult.RefSurface.ID, nil
			}

		}

		// return currentSurfaceRes, nil
	}
	if len(currentSurface) > 0 {
		return currentSurface[0].RefSurface.ID, nil
		// return ConverstCurrentSurface(currentSurface[0]), nil
	} else {
		return 0, nil
		// return ConverstCurrentSurface(models.RoadSurfaceLanePrePareData{}), nil
	}
}

func ConverstCurrentSurface(data models.RoadSurfaceLanePrePareData) responses.RefSurface {
	var refSurface responses.RefSurface
	refSurface.ID = data.RefSurface.ID
	refSurface.Name = data.RefSurface.Name
	refSurface.Type = data.RefSurface.Type
	refSurface.SurfaceGroup = data.RefSurface.SurfaceGroup
	refSurface.LayerCoefficient = data.RefSurface.LayerCoefficient
	refSurface.Drainage = data.RefSurface.Drainage
	refSurface.A = data.RefSurface.A
	refSurface.B = data.RefSurface.B
	cData := strings.Split(data.RefSurface.C, "^")
	if len(cData) > 1 {
		cbase, _ := strconv.ParseFloat(cData[0], 64)
		cexp, _ := strconv.ParseFloat(cData[1], 64)
		refSurface.Cbase = cbase
		refSurface.Cexp = cexp
	} else {
		refSurface.Cbase = 0
		refSurface.Cexp = 0
	}
	crt := 0.0
	if data.RefSurface.CRT != nil {
		crt = *data.RefSurface.CRT
	}

	rrf := 0.0
	if data.RefSurface.RRF != nil {
		rrf = *data.RefSurface.RRF
	}

	raveling_init_a0 := 0.0
	if data.RefSurface.RavelingInitA0 != nil {
		raveling_init_a0 = *data.RefSurface.RavelingInitA0
	}

	raveling_init_a1 := 0.0
	if data.RefSurface.RavelingInitA1 != nil {
		raveling_init_a1 = *data.RefSurface.RavelingInitA1
	}

	raveling_progress_a0 := 0.0
	if data.RefSurface.RavelingProgressA0 != nil {
		raveling_progress_a0 = *data.RefSurface.RavelingProgressA0
	}

	raveling_progress_a1 := 0.0
	if data.RefSurface.RavelingProgressA1 != nil {
		raveling_progress_a1 = *data.RefSurface.RavelingProgressA1
	}

	raveling_progress_a2 := 0.0
	if data.RefSurface.RavelingProgressA2 != nil {
		raveling_progress_a2 = *data.RefSurface.RavelingProgressA2
	}

	asc_init_hsold_equal_0_a0 := 0.0
	if data.RefSurface.AscInitHsoldEqual0A0 != nil {
		asc_init_hsold_equal_0_a0 = *data.RefSurface.AscInitHsoldEqual0A0
	}

	asc_init_hsold_equal_0_a1 := 0.0
	if data.RefSurface.AscInitHsoldEqual0A1 != nil {
		asc_init_hsold_equal_0_a1 = *data.RefSurface.AscInitHsoldEqual0A1
	}

	asc_init_hsold_equal_0_a2 := 0.0
	if data.RefSurface.AscInitHsoldEqual0A2 != nil {
		asc_init_hsold_equal_0_a2 = *data.RefSurface.AscInitHsoldEqual0A2
	}

	asc_init_hsold_equal_0_a3 := 0.0
	if data.RefSurface.AscInitHsoldEqual0A3 != nil {
		asc_init_hsold_equal_0_a3 = *data.RefSurface.AscInitHsoldEqual0A3
	}

	asc_init_hsold_equal_0_a4 := 0.0
	if data.RefSurface.AscInitHsoldEqual0A4 != nil {
		asc_init_hsold_equal_0_a4 = *data.RefSurface.AscInitHsoldEqual0A4
	}

	asc_init_hsold_over_0_a0 := 0.0

	if data.RefSurface.AscInitHsoldOver0A0 != nil {
		asc_init_hsold_over_0_a0 = *data.RefSurface.AscInitHsoldOver0A0
	}

	asc_init_hsold_over_0_a1 := 0.0
	if data.RefSurface.AscInitHsoldOver0A1 != nil {
		asc_init_hsold_over_0_a1 = *data.RefSurface.AscInitHsoldOver0A1
	}

	asc_init_hsold_over_0_a2 := 0.0
	if data.RefSurface.AscInitHsoldOver0A2 != nil {
		asc_init_hsold_over_0_a2 = *data.RefSurface.AscInitHsoldOver0A2
	}

	asc_init_hsold_over_0_a3 := 0.0
	if data.RefSurface.AscInitHsoldOver0A3 != nil {
		asc_init_hsold_over_0_a3 = *data.RefSurface.AscInitHsoldOver0A3
	}

	asc_init_hsold_over_0_a4 := 0.0
	if data.RefSurface.AscInitHsoldOver0A4 != nil {
		asc_init_hsold_over_0_a4 = *data.RefSurface.AscInitHsoldOver0A4
	}

	asc_progress_hsold_equal_0_a0 := 0.0
	if data.RefSurface.AscInitHsoldEqual0A0 != nil {
		asc_init_hsold_over_0_a4 = *data.RefSurface.AscInitHsoldEqual0A0
	}

	asc_progress_hsold_equal_0_a1 := 0.0
	if data.RefSurface.AscInitHsoldEqual0A1 != nil {
		asc_progress_hsold_equal_0_a1 = *data.RefSurface.AscInitHsoldEqual0A1
	}

	asc_progress_hsold_over_0_a0 := 0.0
	if data.RefSurface.AscProgressHsoldOver0A0 != nil {
		asc_progress_hsold_over_0_a0 = *data.RefSurface.AscProgressHsoldOver0A0
	}

	asc_progress_hsold_over_0_a1 := 0.0
	if data.RefSurface.AscProgressHsoldOver0A1 != nil {
		asc_progress_hsold_over_0_a1 = *data.RefSurface.AscProgressHsoldOver0A1
	}

	wsc_init_hsold_equal_0_a0 := 0.0
	if data.RefSurface.WscInitHsoldEqual0A0 != nil {
		wsc_init_hsold_equal_0_a0 = *data.RefSurface.WscInitHsoldEqual0A0
	}

	wsc_init_hsold_equal_0_a1 := 0.0
	if data.RefSurface.WscInitHsoldEqual0A1 != nil {
		wsc_init_hsold_equal_0_a1 = *data.RefSurface.WscInitHsoldEqual0A1
	}

	wsc_init_hsold_equal_0_a2 := 0.0
	if data.RefSurface.WscInitHsoldEqual0A2 != nil {
		wsc_init_hsold_equal_0_a2 = *data.RefSurface.WscInitHsoldEqual0A2
	}

	wsc_init_hsold_over_0_a0 := 0.0
	if data.RefSurface.WscInitHsoldOver0A0 != nil {
		wsc_init_hsold_over_0_a0 = *data.RefSurface.WscInitHsoldOver0A0
	}

	wsc_init_hsold_over_0_a1 := 0.0
	if data.RefSurface.WscInitHsoldOver0A1 != nil {
		wsc_init_hsold_over_0_a1 = *data.RefSurface.WscInitHsoldOver0A1
	}

	wsc_init_hsold_over_0_a2 := 0.0
	if data.RefSurface.WscInitHsoldOver0A0 != nil {
		wsc_init_hsold_over_0_a2 = *data.RefSurface.WscInitHsoldOver0A2
	}

	wsc_progress_hsold_equal_0_a0 := 0.0
	if data.RefSurface.WscProgressHsoldEqual0A0 != nil {
		wsc_progress_hsold_equal_0_a0 = *data.RefSurface.WscProgressHsoldEqual0A0
	}

	wsc_progress_hsold_equal_0_a1 := 0.0
	if data.RefSurface.WscProgressHsoldEqual0A1 != nil {
		wsc_progress_hsold_equal_0_a1 = *data.RefSurface.WscProgressHsoldEqual0A1
	}

	wsc_progress_hsold_over_0_a0 := 0.0

	if data.RefSurface.WscProgressHsoldOver0A0 != nil {
		wsc_progress_hsold_over_0_a0 = *data.RefSurface.WscProgressHsoldOver0A0
	}

	wsc_progress_hsold_over_0_a1 := 0.0
	if data.RefSurface.WscProgressHsoldOver0A1 != nil {
		wsc_progress_hsold_over_0_a1 = *data.RefSurface.WscProgressHsoldOver0A1
	}

	rpd_a0 := 0.0
	if data.RefSurface.RpdA0 != nil {
		rpd_a0 = *data.RefSurface.RpdA0
	}

	rpd_a1 := 0.0
	if data.RefSurface.RpdA1 != nil {
		rpd_a1 = *data.RefSurface.RpdA1
	}

	rpd_a2 := 0.0
	if data.RefSurface.RpdA2 != nil {
		rpd_a2 = *data.RefSurface.RpdA2
	}

	refSurface.Crt = crt
	refSurface.Rrf = rrf
	refSurface.Raveling.Initial.A0 = raveling_init_a0
	refSurface.Raveling.Initial.A1 = raveling_init_a1
	refSurface.Raveling.Progression.A0 = raveling_progress_a0
	refSurface.Raveling.Progression.A1 = raveling_progress_a1
	refSurface.Raveling.Progression.A2 = raveling_progress_a2

	refSurface.AllStructuralCrack.Initial.HSOLD.A0 = asc_init_hsold_equal_0_a0
	refSurface.AllStructuralCrack.Initial.HSOLD.A1 = asc_init_hsold_equal_0_a1
	refSurface.AllStructuralCrack.Initial.HSOLD.A2 = asc_init_hsold_equal_0_a2
	refSurface.AllStructuralCrack.Initial.HSOLD.A3 = asc_init_hsold_equal_0_a3
	refSurface.AllStructuralCrack.Initial.HSOLD.A4 = asc_init_hsold_equal_0_a4

	refSurface.AllStructuralCrack.Initial.HSOLD_O.A0 = asc_init_hsold_over_0_a0
	refSurface.AllStructuralCrack.Initial.HSOLD_O.A1 = asc_init_hsold_over_0_a1
	refSurface.AllStructuralCrack.Initial.HSOLD_O.A2 = asc_init_hsold_over_0_a2
	refSurface.AllStructuralCrack.Initial.HSOLD_O.A3 = asc_init_hsold_over_0_a3
	refSurface.AllStructuralCrack.Initial.HSOLD_O.A4 = asc_init_hsold_over_0_a4

	refSurface.AllStructuralCrack.Progression.HSOLD.A0 = asc_progress_hsold_equal_0_a0
	refSurface.AllStructuralCrack.Progression.HSOLD.A1 = asc_progress_hsold_equal_0_a1

	refSurface.AllStructuralCrack.Progression.HSOLD_O.A0 = asc_progress_hsold_over_0_a0
	refSurface.AllStructuralCrack.Progression.HSOLD_O.A1 = asc_progress_hsold_over_0_a1

	refSurface.WideStructuralCrack.Initial.HSOLD.A0 = wsc_init_hsold_equal_0_a0
	refSurface.WideStructuralCrack.Initial.HSOLD.A1 = wsc_init_hsold_equal_0_a1
	refSurface.WideStructuralCrack.Initial.HSOLD.A2 = wsc_init_hsold_equal_0_a2

	refSurface.WideStructuralCrack.Initial.HSOLD_O.A0 = wsc_init_hsold_over_0_a0
	refSurface.WideStructuralCrack.Initial.HSOLD_O.A1 = wsc_init_hsold_over_0_a1
	refSurface.WideStructuralCrack.Initial.HSOLD_O.A2 = wsc_init_hsold_over_0_a2

	refSurface.WideStructuralCrack.Progression.HSOLD.A0 = wsc_progress_hsold_equal_0_a0
	refSurface.WideStructuralCrack.Progression.HSOLD.A1 = wsc_progress_hsold_equal_0_a1

	refSurface.WideStructuralCrack.Progression.HSOLD_O.A0 = wsc_progress_hsold_over_0_a0
	refSurface.WideStructuralCrack.Progression.HSOLD_O.A1 = wsc_progress_hsold_over_0_a1

	refSurface.RuttingPlasticDeformation.A0 = rpd_a0
	refSurface.RuttingPlasticDeformation.A1 = rpd_a1
	refSurface.RuttingPlasticDeformation.A2 = rpd_a2

	return refSurface
}

func (u *useCase) GetRoadCondition(roadID, laneNo int) ([]models.RoadConditionSurveyM, time.Time, error) {
	data, surveyedDate, err := u.repo.GetRoadConditionData(roadID, laneNo)
	if err != nil {
		logs.Error(err)
		return []models.RoadConditionSurveyM{}, surveyedDate, err
	}
	return data, surveyedDate, nil
}

func (u *useCase) GetRoadConditionPrepareData(directionID, laneNo int, roadConditions []models.RoadConditionSurveyM, road map[string]interface{}) (responses.RoadConditionPrepareDataRes, error) {
	// helpers.PrintlnJson("roadConditions", roadConditions)
	km_start := road["km_start"].(float64)
	km_end := road["km_end"].(float64)
	rut := 0.0
	lenght_rut := 0.0

	iri := 0.0
	lenght_iri := 0.0

	gn := 0.0
	lenght_gn := 0.0
	// var newKm NewKm
	/* หาช่วงกมเริ่มต้น, สิ้นสุด ที่ลงกับค่าสำรวจพอดี */
	var newKmStart interface{}
	var newKmEnd interface{}
	if directionID == 1 { // LT
		sort.Slice(roadConditions, func(i, j int) bool {
			return roadConditions[i].KmStart < roadConditions[j].KmEnd
		})
		for _, conditionM := range roadConditions {

			if km_start == conditionM.KmStart {
				newKmStart = conditionM.KmStart
				break
			} else if km_start < conditionM.KmStart {
				diff := math.Abs(km_start - conditionM.KmStart)
				if diff > 12.5 {
					data := conditionM.KmStart - 25
					newKmStart = data
					break
				} else {
					newKmStart = conditionM.KmStart
					break
				}
			}
		}

		for _, conditionM := range roadConditions {
			if km_end == conditionM.KmEnd {
				newKmEnd = conditionM.KmEnd
				break
			} else if km_end < conditionM.KmEnd {
				diff := math.Abs(km_end - conditionM.KmEnd)
				if diff > 12.5 {
					data := conditionM.KmEnd - 25
					newKmEnd = data
					break
				} else {
					newKmEnd = conditionM.KmEnd
					break
				}

			}
		}
	} else {
		sort.Slice(roadConditions, func(i, j int) bool {
			return roadConditions[i].KmStart > roadConditions[j].KmEnd
		})
		for _, conditionM := range roadConditions {
			if km_start == conditionM.KmStart {
				newKmStart = conditionM.KmStart
				break
			} else if km_start > conditionM.KmStart {
				diff := math.Abs(km_start - conditionM.KmStart)
				if diff > 12.5 {
					data := conditionM.KmStart - 25
					newKmStart = data
					break
				} else {
					newKmStart = conditionM.KmStart
					break
				}
			}
		}

		for _, conditionM := range roadConditions {
			if km_end == conditionM.KmEnd {
				newKmEnd = conditionM.KmEnd
				fmt.Println("000", newKmEnd)
				break
			} else if km_end > conditionM.KmEnd {
				diff := math.Abs(km_end - conditionM.KmEnd)
				if diff > 12.5 {
					data := conditionM.KmEnd - 25
					newKmEnd = data
					fmt.Println("111", newKmEnd)
					break
				} else {
					newKmEnd = conditionM.KmEnd
					fmt.Println("2222", newKmEnd)
					break
				}

			}
		}
	}

	// km_start2 := newKmStart
	// km_end2 := newKmEnd
	// helpers.PrintlnJson(road["km_start"].(float64), road["km_end"].(float64), km_start, km_end2)
	// helpers.PrintlnJson(km_start, km_end, newKmStart, newKmEnd)
	// return []models.RoadSurfacePrepareData{}, nil
	// if newKmStart == nil {
	// 	km_start = 0.0
	// } else {
	// 	km_start = *newKmStart
	// }
	// if newKmEnd == nil {
	// 	km_end = 0.0
	// } else {
	// 	km_end = *newKmEnd
	// }
	// helpers.PrintlnJson("km_end, km_end", newKmStart, newKmEnd)
	if newKmStart == nil || newKmEnd == nil {
		// helpers.PrintlnJson("roadConditions", roadConditions)
	} else {
		// helpers.PrintlnJson("newKmStart", newKmStart, newKmEnd)
		if newKmStart.(float64) < newKmEnd.(float64) { // LT

			for _, conditionM := range roadConditions {
				if (conditionM.KmStart >= newKmStart.(float64)) && (conditionM.KmEnd <= newKmEnd.(float64)) {

					lengh := math.Abs(conditionM.KmEnd - conditionM.KmStart)

					// rut
					if conditionM.RUT != nil {
						rut += *conditionM.RUT * lengh
						lenght_rut += lengh
					}

					// iri
					if conditionM.IRI != nil {
						iri += *conditionM.IRI * lengh
						lenght_iri += lengh
					}

					// ifi
					if conditionM.IFI != nil {
						gn += *conditionM.IFI * lengh
						lenght_gn += lengh
					}
					// fmt.Println("*conditionM.IRI", conditionM.KmStart, &conditionM.KmEnd, conditionM.IRI)
				}
			}
		} else { // RT
			for _, conditionM := range roadConditions {
				if (conditionM.KmStart <= newKmStart.(float64)) && (conditionM.KmEnd >= newKmEnd.(float64)) {

					lengh := math.Abs(conditionM.KmEnd - conditionM.KmStart)

					// rut
					if conditionM.RUT != nil {
						rut += *conditionM.RUT * lengh
						lenght_rut += lengh
					}

					// iri
					if conditionM.IRI != nil {
						iri += *conditionM.IRI * lengh
						lenght_iri += lengh
					}

					// ifi
					if conditionM.IFI != nil {
						gn += *conditionM.IFI * lengh
						lenght_gn += lengh
					}
					// fmt.Println("*conditionM.IRI", conditionM.KmStart, &conditionM.KmEnd, conditionM.IRI)
				}
			}
		}

	}

	rut = (rut / lenght_rut)
	iri = (iri / lenght_iri)
	gn = (gn / lenght_gn)
	// helpers.PrintlnJson(rut, iri, gn)
	var roadConditionPrepareDataRes responses.RoadConditionPrepareDataRes
	roadConditionPrepareDataRes.RUT = rut
	roadConditionPrepareDataRes.IRI = iri
	roadConditionPrepareDataRes.IFI = gn
	return roadConditionPrepareDataRes, nil
}

func (u *useCase) GetType(road map[string]interface{}, surveyDate interface{}) (int, int) {
	refStructureSurface := road["surface"].(models.RoadSurfacePrepareData).RefStructureSurface
	structureSurfaceLastType := 0
	structureSurfaceType := 0
	asphaltCountLastType := 0 // ให้ทำการวนลูปเช็คประวัติการซ่อมบำรุงว่าเคยถูกซ่อมด้วย OL-Overlay ไหม
	asphaltCountType := 0
	roadSurface := road["surface"].(models.RoadSurfacePrepareData)
	roadSurfacLane := roadSurface.RoadSurfaceLane
	refStructureSurfaceID := refStructureSurface.ID
	if len(roadSurfacLane) > 0 {
		roadSurfaceId := roadSurfacLane[0].RefSurface
		switch refStructureSurfaceID {
		case 1:
			if roadSurfaceId.SurfaceGroup == "Asphalt" {
				refStructureSurfaceID = refStructureSurface.ID //Asphalt: On Concrete Deck
			} else {
				refStructureSurfaceID = 6 //concrete: on Deck
			}
		case 2:
			if roadSurfaceId.SurfaceGroup == "Asphalt" {
				refStructureSurfaceID = refStructureSurface.ID //Asphalt: On Steel Deck
			} else {
				refStructureSurfaceID = 6 //concrete: on Deck
			}
		case 3:
			if roadSurfaceId.SurfaceGroup == "Asphalt" {
				refStructureSurfaceID = refStructureSurface.ID //Asphalt: On Ground
			} else {
				refStructureSurfaceID = 5 //Concrete: on Ground
			}
		case 4:
			if roadSurfaceId.SurfaceGroup == "Asphalt" {
				refStructureSurfaceID = refStructureSurface.ID //Composite pavement (Asphalt on concrete base)
			} else {
				refStructureSurfaceID = 5 //Concrete: on Ground
			}
		case 5:
			if roadSurfaceId.SurfaceGroup == "Asphalt" {
				refStructureSurfaceID = 3 //Asphalt: On Ground
			} else {
				refStructureSurfaceID = refStructureSurface.ID //Concrete: on Ground
			}
		case 6:
			if roadSurfaceId.SurfaceGroup == "Asphalt" {
				refStructureSurfaceID = 1 //Asphalt: On Concrete Deck
			} else {
				refStructureSurfaceID = refStructureSurface.ID //concrete: on Deck
			}
		default:
			refStructureSurfaceID = refStructureSurface.ID
		}
	} else {
		refStructureSurfaceID = refStructureSurface.ID
	}

	if refStructureSurfaceID != 5 { // กรณีหน้าตัดผิวทางไม่ใช่ Concrete: on Ground
		structureSurfaceLastType = refStructureSurfaceID
		structureSurfaceType = refStructureSurfaceID
		return structureSurfaceLastType, structureSurfaceType
	}
	concreteCountLastType := 0
	concreteCountType := 0
	if refStructureSurfaceID != 5 { // กรณีหน้าตัดผิวทางไม่ใช่ Concrete: on Ground
		if road["maintenances"] != nil {
			maintenances := road["maintenances"].([]responses.MaintenanceData)
			for _, item := range maintenances { //check case RBC เพื่อเปลี่ยนเป็น type = 5,6
				methods := []string{"RBC"}
				maintenanceMethod := item.IcResult.Method
				if helpers.InArray(maintenanceMethod, methods) {
					concreteCountLastType++
					concreteCountType++
				}

			}
			if concreteCountLastType > 0 { // กรณีมีประวัติซ่อมด้วย RBC ให้มองช่วงกม. นี้เป็น concrete pavement
				structureSurfaceLastType = 5
			} else {
				structureSurfaceLastType = refStructureSurfaceID
			}
		} else {
			structureSurfaceLastType = refStructureSurfaceID
			structureSurfaceType = refStructureSurfaceID
		}
		//เพิ่มตรงนี้
	} else { // กรณีหน้าตัดผิวทางเป็นประเภท Concrete: on Ground
		if road["maintenances"] != nil {
			maxYear := 0 // Assuming year cannot be negative. Otherwise, use math.MinInt or similar.
			for _, maintenance := range road["maintenances"].([]responses.MaintenanceData) {
				if maintenance.Year > maxYear {
					maxYear = maintenance.Year
				}
			}
			maintenances := road["maintenances"].([]responses.MaintenanceData)
			for _, item := range maintenances {
				if item.Year == maxYear {
					methods := []string{"OL-Overlay", "M&OL-Mill&Overlay"}
					maintenanceMethod := item.IcResult.Method
					if helpers.InArray(maintenanceMethod, methods) {
						// if surveyDate == nil {
						// 	continue
						// }
						asphaltCountLastType++
						// if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
						asphaltCountType++
						// }
					}
				}
			}

			if asphaltCountLastType > 0 { // กรณีมีประวัติซ่อมด้วย OL-Overlay ให้มองช่วงกม. นี้เป็น Composite pavement (Asphalt on concrete base)
				structureSurfaceLastType = 4
			} else {
				structureSurfaceLastType = refStructureSurfaceID
			}

			if asphaltCountType > 0 { // กรณีมีประวัติซ่อมด้วย OL-Overlay ให้มองช่วงกม. นี้เป็น Composite pavement (Asphalt on concrete base)
				structureSurfaceType = 4
			} else {
				structureSurfaceType = refStructureSurfaceID
			}
		} else {
			structureSurfaceLastType = refStructureSurfaceID
			structureSurfaceType = refStructureSurfaceID
		}
		return structureSurfaceLastType, structureSurfaceType
	}
	return structureSurfaceLastType, structureSurfaceType
}

func (u *useCase) GetAgeAC(road map[string]interface{}, surveyDate interface{}) (interface{}, error) {
	analystYear := road["analyst_year"].(int)                                // ปีที่วิเคราะห์ข้อมูล
	yearRoadBegin := road["road_date_begin"].(responses.RoadDateBegins).Year // ปีที่เปิดใช้งาน

	yearLastOverlays := []int{}
	yearLastSeals := []int{}
	yearLastMolRcls := []int{}
	yearLastReconstructions := []int{}

	maintenances := road["maintenances"].([]responses.MaintenanceData)
	if len(maintenances) > 0 {
		for _, item := range maintenances { // เรียงตามปีที่น้อยไปมาก
			// if surveyDate == nil {
			// 	continue
			// }
			method := item.IcResult.Method
			switch method { // Condition ในการคิด HSOLD
			case "OL-Overlay": // ลาดยาง (AC)
				// if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
				yearLastOverlays = append(yearLastOverlays, item.Year)
				// }
			case "SS-SlurrySeal": // ลาดยาง (AC)
				// if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
				yearLastSeals = append(yearLastSeals, item.Year)
				// }
			case "M&OL-Mill&Overlay": // ลาดยาง (AC)
				// if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
				yearLastMolRcls = append(yearLastMolRcls, item.Year)
				// }
			case "RCL-Recycling": // ลาดยาง (AC)
				// if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
				yearLastMolRcls = append(yearLastMolRcls, item.Year)
				// }
			case "Rc-Reconstruction": // ลาดยาง (AC)
				// if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
				yearLastReconstructions = append(yearLastReconstructions, item.Year)
				// }
			}
		}

	}
	/* หาปีล่าสุดของการซ่อมแต่ละประเภท */
	yearLastOverlay := helpers.FindMaxInt(yearLastOverlays) // max($year_last_overlay);
	yearLastSeal := helpers.FindMaxInt(yearLastSeals)
	yearLastMolRcl := helpers.FindMaxInt(yearLastMolRcls)
	yearLastReconstruction := helpers.FindMaxInt(yearLastReconstructions)
	age := analystYear - helpers.FindMaxInt([]int{yearRoadBegin, yearLastOverlay, yearLastSeal, yearLastMolRcl, yearLastReconstruction})
	if age < 0 {
		age = 0
	}
	var agePrepareDataRes responses.AgePrepareDataRes
	agePrepareDataRes.Age = age
	agePrepareDataRes.YearLastMolRcl = yearLastMolRcl
	agePrepareDataRes.YearLastOverlay = yearLastOverlay
	agePrepareDataRes.YearLastReconstruction = yearLastReconstruction
	agePrepareDataRes.YearLastSeal = yearLastSeal
	return agePrepareDataRes, nil
}

func (u *useCase) GetAgeCC(road map[string]interface{}, surveyDate interface{}) (interface{}, error) {
	analystYear := road["analyst_year"].(int)                                // ปีที่วิเคราะห์ข้อมูล
	yearRoadBegin := road["road_date_begin"].(responses.RoadDateBegins).Year // ปีที่เปิดใช้งาน

	yearLastFdrs := []int{}
	yearLastOvls := []int{}
	yearLastSeals := []int{}
	yearLastMols := []int{}

	maintenances := road["maintenances"].([]responses.MaintenanceData)
	if len(maintenances) > 0 {
		for _, item := range maintenances { // เรียงตามปีที่น้อยไปมาก
			// if surveyDate == nil {
			// 	continue
			// }
			method := item.IcResult.Method
			switch method { // Condition ในการคิด HSOLD
			case "FDR": // คอนกรีต (CC)
				// if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
				yearLastFdrs = append(yearLastFdrs, item.Year)
				// }
			case "OVL": // คอนกรีต (CC)
				// if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
				yearLastOvls = append(yearLastOvls, item.Year)
				// }
			case "M-OL": // คอนกรีต (CC)
				// if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
				yearLastMols = append(yearLastMols, item.Year)
				// }
			case "Seal": // คอนกรีต (CC)
				// if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
				yearLastSeals = append(yearLastSeals, item.Year)
				// }
			}
		}

	}
	/* หาปีล่าสุดของการซ่อมแต่ละประเภท */
	yearLastFdr := helpers.FindMaxInt(yearLastFdrs) // max($year_last_overlay);
	yearLastOvl := helpers.FindMaxInt(yearLastOvls)
	yearLastMol := helpers.FindMaxInt(yearLastMols)
	yearLastSeal := helpers.FindMaxInt(yearLastSeals)
	age := analystYear - helpers.FindMaxInt([]int{yearRoadBegin, yearLastFdr, yearLastOvl, yearLastMol, yearLastSeal})
	if age < 0 {
		age = 0
	}
	var agePrepareDataRes responses.AgePrepareDataRes
	agePrepareDataRes.Age = age
	agePrepareDataRes.YearLastFdr = yearLastFdr
	agePrepareDataRes.YearLastOvl = yearLastOvl
	agePrepareDataRes.YearLastMol = yearLastMol
	agePrepareDataRes.YearLastSeal = yearLastSeal
	return agePrepareDataRes, nil
}

func (u *useCase) GetDataMartCheck() (interface{}, error) {
	data, err := u.repo.GetDataMartCheck()
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	return data, nil
}

func (u *useCase) GetSurfaceDashboard(roadIDs []int, depotCodes []string, filter requests.Asset) (interface{}, error) {
	var result responses.SurfaceRespond
	surfaceInfo, err := u.repo.GetDataMartInfo(roadIDs, depotCodes, filter)
	if err != nil {
		return result, err
	}
	if len(surfaceInfo) == 0 {
		return surfaceInfo, nil
	}
	summary, err := u.repo.GetInitialSurfaceArray()
	if err != nil {
		return result, err
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
	// helpers.SortStructByID(result.GeomList, false)

	return result, nil
}

func (u *useCase) GetSurfaceDashboardMap(roadIDs []int, depotCodes []string, filter requests.Asset, display int) (interface{}, error) {
	var result responses.SurfaceRespond

	surfaceInfo, err := u.repo.GetDataMartInfo(roadIDs, depotCodes, filter)
	if err != nil {
		return result, err
	}
	// return surfaceInfo, nil
	if len(surfaceInfo) == 0 {
		return surfaceInfo, nil
	}

	var geomLists []models.GeomList
	for _, surface := range surfaceInfo {
		var geomList models.GeomList
		// Unmarshal the JSON data into the geom variable
		var geom Geometry
		err := json.Unmarshal([]byte(surface.Geometry), &geom)
		if err != nil {
			fmt.Println(err)
			// return
		}

		age := 0
		years := time.Now().Year()
		if surface.Year != 0 {
			age = int(math.Abs(float64(surface.Year) - float64(years)))
		} else {
			age = surface.Age
		}

		if display == 1 {
			geomList.Title = fmt.Sprintf("อายุผิว %d ปี", age)
			if age > 10 {
				geomList.Color = "#DC3545"
			} else if age >= 6 && age <= 10 {
				geomList.Color = "#FDB833"
			} else if age >= 3 && age <= 5 {
				geomList.Color = "#87C442"
			} else if age >= 0 && age <= 2 {
				geomList.Color = "#50CD89"
			}
			geomList.SurfaceName = surface.SurfaceName
		} else {
			geomList.Title = surface.SurfaceName
			geomList.Color = surface.ColorCode
			if surface.SurfaceGroup == "Asphalt" {
				geomList.SurfaceName = "ยางมะตอย"
			} else {
				geomList.SurfaceName = "คอนกรีต"
			}
			// geomList.SurfaceGroup
		}
		year := fmt.Sprintf("%d", surface.Year)
		if year == "0" {
			year = ""
		}

		geomList.TheGeom = geom
		geomList.KmStart = surface.KmStart
		geomList.KmEnd = surface.KmEnd
		geomList.KmTotal = math.Abs(surface.KmStart-surface.KmEnd) / 1000
		geomList.Year = year
		geomList.Age = age
		geomList.ContractNumber = surface.ContractNumber
		lastInspectionDate := surface.LastInspectionDate.Format("2006-01-02")
		if lastInspectionDate == "0001-01-01" {
			lastInspectionDate = ""
		}
		geomList.LastInspectionDate = &lastInspectionDate
		// geomList.SurfaceName = surface.SurfaceName

		geomList.RoadGroupName = surface.RoadGroupName
		geomList.RoadName = surface.RoadName
		geomLists = append(geomLists, geomList)

	}
	return geomLists, nil
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
			s.ColorCode = s.ColorCode
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
	helpers.PrintlnJson(laneTypeGroup)
	// mapGeomID := make(map[string]models.GeomList)
	helpers.PrintlnJson("len", len(surfaceInfo))
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
			// 	result.Detail.LaneCountList = append(result.Detail.LaneCountList, s.LaneNo)
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
		}
		helpers.PrintlnJson("xxx", mapDetailIdxs)
		for _, detail := range mapDetailIdxs {
			result.Detail.DetailKm = append(result.Detail.DetailKm, detail)
		}
	}
	laneTypeInfo := make(map[string]responses.SurfaceLaneType)
	surfaceKey := make(map[string]int)
	for _, v := range result.Detail.DetailKm {
		helpers.PrintlnJson(v.LaneType, v.Value)
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
	// 	result.GeomList = append(result.GeomList, geom)
	// }
	return result
}

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
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
	helpers.PrintlnJson(count)
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
