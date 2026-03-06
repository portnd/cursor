package usecases

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	helpers "gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	responses "gitlab.com/mims-api-service/responses"
)

func (t *roadConditionUseCase) GetRoadConditionTemplate(uid uint, roadId int) (interface{}, error) {
	rcTemplate, err := t.roadConditionRepo.GetRoadConditionTemplate(roadId)
	if err != nil {
		return "", err
	}
	// return rcTemplate, nil
	header := []string{"road_id", "road_code", "name", "km_start", "km_end", "img_filepath", "iri", "mpd", "rut", "ifi", "survey_type"}
	filePath := "storages/template_condition/"

	fileName := ""
	fileLocations := []string{}
	haveRoadSurface := false

	var intervalKmRanges = make([]map[string]interface{}, 0)
	// return rcTemplate, nil
	for _, roadGeom := range rcTemplate.RoadGeom {
		if roadGeom.Status != "A" {
			continue
		}
		kms := make(map[float64]bool)
		var keys []float64
		dataKmMap := make(map[string]models.RoadSurfaceData)

		// get road surface
		roadSurfaces, err := t.roadConditionRepo.GetRoadSurfaceByRoadIDByLaneNo(roadId, roadGeom.LaneNo)
		if err != nil {
			return "", err
		}

		var roadSurfacesCleanData []models.RoadSurfaceData
		for _, surface := range roadSurfaces {
			if surface.RoadSurfaceLane.Id != 0 {
				roadSurfacesCleanData = append(roadSurfacesCleanData, surface)
			}
		}
		sort.SliceStable(roadSurfacesCleanData, func(i, j int) bool {
			return roadSurfacesCleanData[i].KmStart < roadSurfacesCleanData[j].KmEnd
		})
		roadSurfaceData := checkConsecutiveRoads(roadSurfacesCleanData)
		roadSurfaceDatas := removeDuplicates(roadSurfaceData)
		var roadSurfaceFinal []models.RoadSurfaceData
		for i, item := range roadSurfaceDatas {
			if i-1 > 0 {
				if item.KmEnd != roadSurfaceDatas[i-1].KmEnd {
					roadSurfaceFinal = append(roadSurfaceFinal, item)
				}
			} else {
				roadSurfaceFinal = append(roadSurfaceFinal, item)
			}
		}
		for _, surface := range roadSurfaceFinal {
			if surface.RoadSurfaceLane.RefSurfaceId == 0 {
				continue
			}
			haveRoadSurface = true
			var kmStart float64
			var kmEnd float64
			if rcTemplate.RoadInfo.RefDirectionId == 1 {
				if (surface.KmStart >= roadGeom.KmStart) && (surface.KmStart < roadGeom.KmEnd) {
					kms[surface.KmStart] = true
					kmStart = surface.KmStart
				} else {
					if (surface.KmStart < roadGeom.KmStart) && (surface.KmStart < roadGeom.KmEnd) {
						kms[roadGeom.KmStart] = true
						kmStart = roadGeom.KmStart
					}
				}

				if (surface.KmEnd <= roadGeom.KmEnd) && (surface.KmEnd > roadGeom.KmStart) {
					kms[surface.KmEnd] = true
					kmEnd = surface.KmEnd
				}
				if kmStart > kmEnd {
					continue
				}
			} else {
				if (surface.KmStart <= roadGeom.KmStart) && (surface.KmStart > roadGeom.KmEnd) {
					kms[surface.KmStart] = true
					kmStart = surface.KmStart
				}

				if (surface.KmEnd >= roadGeom.KmEnd) && (surface.KmEnd < roadGeom.KmStart) {
					kms[surface.KmEnd] = true
					kmEnd = surface.KmEnd
				}
			}

			key := fmt.Sprintf("%v-%v", kmStart, kmEnd)
			helpers.PrintlnJson("key", key)
			dataKmMap[key] = surface

		}
		for key := range kms {
			keys = append(keys, key)
		}
		sort.Float64s(keys)
		kmRanges := make([]map[string]float64, 0)
		/* แบ่งช่วงตามจุดตัด */

		for index := 0; index < len(keys)-1; index++ {
			if rcTemplate.RoadInfo.RefDirectionId == 1 {
				kmRange := make(map[string]float64)
				kmRange["km_start"] = keys[index]
				kmRange["km_end"] = keys[index+1]

				key := fmt.Sprintf("%v-%v", keys[index], keys[index+1])
				refSurface := dataKmMap[key]
				kmRange["ref_surface_id"] = float64(refSurface.RoadSurfaceLane.RefSurfaceId)
				kmRange["lane_no"] = float64(refSurface.RoadSurfaceLane.LaneNo)
				kmRanges = append(kmRanges, kmRange)
			} else {
				kmRange := make(map[string]float64)
				kmRange["km_start"] = keys[index+1]
				kmRange["km_end"] = keys[index]

				key := fmt.Sprintf("%v-%v", keys[index+1], keys[index])
				refSurface := dataKmMap[key]
				kmRange["ref_surface_id"] = float64(refSurface.RoadSurfaceLane.RefSurfaceId)
				kmRange["lane_no"] = float64(refSurface.RoadSurfaceLane.LaneNo)
				kmRanges = append(kmRanges, kmRange)
			}
		}
		// return kmRanges, nil
		/* แบ่งช่วงทุกๆ 25 เมตร */
		interval := 25
		for _, km := range kmRanges {
			kmStart := km["km_start"]
			kmEnd := km["km_end"]
			refSurfaceID := km["ref_surface_id"]
			laneNo := km["lane_no"]

			intvalStart := math.Floor(kmStart / float64(interval)) // หารเอาจำนวนเต็มไม่ปัดเศษ
			intvalEnd := math.Floor(kmEnd / float64(interval))     // หารเอาจำนวนเต็มไม่ปัดเศษ

			// if intvalStart == intvalEnd || math.Abs(kmEnd-kmStart) <= 500 { // กรณีระยะทางรวมไม่เกิน 500 ไม่ต้องแบ่งตาม interval
			// 	intervalKmRanges = append(intervalKmRanges, map[string]interface{}{
			// 		"km_start":       kmStart,
			// 		"km_end":         kmEnd,
			// 		"ref_surface_id": refSurfaceID,
			// 		"lane_no":        laneNo,
			// 	})
			// } else {
			start := kmStart
			// intvalStart := intvalStart
			// intvalEnd := intvalEnd
			var end float64
			if kmStart < kmEnd {
				for start < kmEnd {
					lt := start + float64(interval)
					end = math.Min(lt, kmEnd)

					if intvalStart != intvalEnd {
						end = end - math.Mod(end, float64(interval))
					}

					intervalKmRanges = append(intervalKmRanges, map[string]interface{}{
						"km_start":       start,
						"km_end":         end,
						"ref_surface_id": refSurfaceID,
						"lane_no":        laneNo,
					})

					start = end
					intvalStart = intvalEnd
				}
			} else {
				for start > kmEnd {
					rt := math.Ceil((start-float64(interval))/float64(interval)) * float64(interval)
					end = math.Max(rt, kmEnd)

					if intvalStart != intvalEnd {
						end = end - math.Mod(end, float64(interval))
					}

					intervalKmRanges = append(intervalKmRanges, map[string]interface{}{
						"km_start":       start,
						"km_end":         end,
						"ref_surface_id": refSurfaceID,
						"lane_no":        laneNo,
					})

					start = end
					intvalStart = intvalEnd
				}
			}
			// }
		}
	}
	// return intervalKmRanges, nil
	// return intervalKmRanges, nil
	datas := make(map[int][]models.RoadKmRage25M)
	switch rcTemplate.RoadInfo.RefDirectionId {
	case 1:
		sort.Slice(intervalKmRanges, func(i, j int) bool {
			return intervalKmRanges[i]["km_start"].(float64) < intervalKmRanges[j]["km_end"].(float64)
		})

	case 2:
		sort.Slice(intervalKmRanges, func(i, j int) bool {
			return intervalKmRanges[i]["km_start"].(float64) > intervalKmRanges[j]["km_end"].(float64)
		})
	}

	if !haveRoadSurface {
		return "", responses.NewAppErr(http.StatusBadRequest, "ไม่สามารถสร้าง csv template ได้เนื่องจากสายทางนี้ไม่มีข้อมูลผิวทาง")
	}
	for _, item := range intervalKmRanges {
		var roadKmRage models.RoadKmRage25M
		roadKmRage.KmStart = item["km_start"].(float64)
		roadKmRage.KmEnd = item["km_end"].(float64)
		roadKmRage.RoadID = rcTemplate.RoadInfo.RoadId
		roadKmRage.RoadCode = rcTemplate.RoadSection.RoadGroup.Number + rcTemplate.RoadSection.Number // rcTemplate.RoadCode
		roadKmRage.RoadName = rcTemplate.RoadInfo.Name
		roadKmRage.SurveyType = int(item["ref_surface_id"].(float64))
		datas[int(item["lane_no"].(float64))] = append(datas[int(item["lane_no"].(float64))], roadKmRage)
	}
	refSurfaceData, _ := t.roadConditionRepo.GetRefSurface()
	for i, data := range datas {
		// if data
		csvData := [][]string{}
		for _, value := range data {
			surveyType := refSurfaceData[value.SurveyType]
			start := fmt.Sprintf("%d", value.RoadID) + "," + value.RoadCode + "," + value.RoadName + "," + fmt.Sprintf("%v", value.KmStart) + "," + fmt.Sprintf("%v", value.KmEnd) + "," + "" + "," + fmt.Sprintf("%d", 0) + "," + fmt.Sprintf("%d", 0) + "," + fmt.Sprintf("%d", 0) + "," + fmt.Sprintf("%d", 0) + "," + surveyType
			csvData = append(csvData, strings.Split(start, ","))
		}

		fileName = "template_condition_" + rcTemplate.RoadSection.RoadGroup.Number + rcTemplate.RoadSection.Number + "_road_id_" + fmt.Sprintf("%d", rcTemplate.Road.Id) + "_lane" + fmt.Sprintf("%d", i) + ".csv"
		err = csvGenerateFile(header, csvData, filePath, fileName)
		if err != nil {
			return "", responses.NewAppErr(http.StatusBadRequest, err.Error())
		}
		fileLocation := filePath + fileName

		err = os.Chmod(fileLocation, 0775)
		if err != nil {
			return "", responses.NewAppErr(http.StatusBadRequest, err.Error())
		}

		fileLocations = append(fileLocations, fileLocation)
	}

	STORAGE_IP := os.Getenv("STORAGE_IP") + "/"
	zipName := filePath + "template_condition_" + time.Now().Format("2006-01-02_150405") + ".zip"
	err = helpers.Zip(zipName, fileLocations)
	if err != nil {
		return "", responses.NewAppErr(http.StatusBadRequest, err.Error())
	}
	err = os.Chmod(zipName, 0775)
	if err != nil {
		return "", responses.NewAppErr(http.StatusBadRequest, err.Error())
	}

	for _, fileLocation := range fileLocations {
		err = os.Remove(fileLocation)
		if err != nil {
			return nil, err
		}
	}

	return STORAGE_IP + zipName, nil
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

func checkConsecutiveRoads(roads []models.RoadSurfaceData) []models.RoadSurfaceData {
	temp := make(map[string]models.RoadSurfaceData)
	temp2 := make(map[string]models.RoadSurfaceData)
	var first float64
	data := []models.RoadSurfaceData{}
	isFirst := true
	if len(roads) > 1 {
		for i := 0; i <= len(roads)-1; i++ {
			if i == len(roads)-1 {
				currentRoad := roads[len(roads)-1]
				key := fmt.Sprintf("%v-%v", currentRoad.RoadSurfaceLane.LaneNo, currentRoad.RoadSurfaceLane.RefSurfaceId)
				temp[key] = currentRoad
				data = append(data, temp[key])
				return data
			}
			currentRoad := roads[i]
			nextRoad := roads[i+1]

			if isFirst {
				first = currentRoad.KmStart
			}

			if currentRoad.RoadSurfaceLane.LaneNo == nextRoad.RoadSurfaceLane.LaneNo &&
				currentRoad.RoadSurfaceLane.RefSurfaceId == nextRoad.RoadSurfaceLane.RefSurfaceId &&
				currentRoad.KmEnd == nextRoad.KmStart &&
				currentRoad.Id == nextRoad.Id {
				currentRoad.KmStart = first
				currentRoad.KmEnd = nextRoad.KmEnd
				key := fmt.Sprintf("%v-%v", currentRoad.RoadSurfaceLane.LaneNo, currentRoad.RoadSurfaceLane.RefSurfaceId)
				temp[key] = currentRoad
				isFirst = false
				if i < len(roads)-2 {
					continue
				}

			} else {
				isFirst = true
				key := fmt.Sprintf("%v-%v", currentRoad.RoadSurfaceLane.LaneNo, currentRoad.RoadSurfaceLane.RefSurfaceId)
				temp2[key] = currentRoad
			}
			_, isval := temp[fmt.Sprintf("%v-%v", currentRoad.RoadSurfaceLane.LaneNo, currentRoad.RoadSurfaceLane.RefSurfaceId)]
			if isval {
				data = append(data, temp[fmt.Sprintf("%v-%v", currentRoad.RoadSurfaceLane.LaneNo, currentRoad.RoadSurfaceLane.RefSurfaceId)])
			}

			key2 := fmt.Sprintf("%v-%v", currentRoad.RoadSurfaceLane.LaneNo, currentRoad.RoadSurfaceLane.RefSurfaceId)
			_, isval = temp2[key2]
			if isval {
				data = append(data, temp2[key2])
			}
		}
	} else {
		if len(roads) > 0 {
			currentRoad := roads[0]
			key := fmt.Sprintf("%v-%v", currentRoad.RoadSurfaceLane.LaneNo, currentRoad.RoadSurfaceLane.RefSurfaceId)
			temp[key] = currentRoad
			data = append(data, temp[key])

		}
	}
	return data
}

func removeDuplicates(finalDatas []models.RoadSurfaceData) []models.RoadSurfaceData {
	keys := make(map[string]bool)
	result := []models.RoadSurfaceData{}

	for _, entry := range finalDatas {
		key := generateKey(entry)
		if _, value := keys[key]; !value {
			keys[key] = true
			result = append(result, entry)
		}
	}

	return result
}

func generateKey(data models.RoadSurfaceData) string {
	key, _ := json.Marshal(data)
	return string(key)
}
