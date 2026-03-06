package usecases

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	requests "gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadDamage/domains"
)

type roadDamageUseCase struct {
	repo domains.RoadDamageRepository
}

// init usecase
func NewRoadDamageUseCase(repo domains.RoadDamageRepository) domains.RoadDamageUseCase {
	return &roadDamageUseCase{
		repo: repo,
	}
}

// =========================================================

// =========================================================

func (t *roadDamageUseCase) GetRoadDamageList(roleId int, permissions []string) (interface{}, error) {
	data, err := t.repo.GetRoadDamageList(roleId, permissions)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	var roadDamages []models.RoadDamageListResponse
	year := make(map[int]int)
	dataYear := make(map[int][]models.RoadDamageListResponse)
	temp := make(map[int]bool)
	for _, item := range data {
		_, isVal := temp[item.IdParent]
		if !isVal {
			_, ok := year[item.Year]
			var roadDamage models.RoadDamageListResponse
			var refDirection models.RefDirection
			if !ok {
				year[item.Year] = item.Year
			}
			refDirection.ID = item.DirectionId
			refDirection.Name = item.DirectionName
			roadDamage.Id = item.Id
			roadDamage.IdParent = item.IdParent
			roadDamage.Direction = append(roadDamage.Direction, refDirection)
			roadDamage.LaneNo = item.LaneNo
			roadDamage.SurveyedDate = item.SurveyedDate
			roadDamage.Revision = item.Revision
			roadDamages = append(roadDamages, roadDamage)

			temp[item.IdParent] = true
			dataYear[item.SurveyedDate.Year()] = roadDamages
		}
	}
	var response []responses.RoadDamageList
	for year, items := range dataYear {
		var data responses.RoadDamageList
		data.Year = year

		for _, item := range items {
			if item.SurveyedDate.Year() == year {
				data.Item = append(data.Item, item)
			}
		}
		response = append(response, data)
	}
	return response, nil
}

func (t *roadDamageUseCase) GetDirectionById(roadID int) (models.RefDirection, error) {
	direction, err := t.repo.GetDirectionById(roadID)
	if err != nil {
		logs.Error(err)
		return direction, responses.NewAppErr(400, err.Error())
	}
	return direction, nil
}

func (t *roadDamageUseCase) GetRoadDamageForImport(roadID, IDParent int) (models.RoadDamage, error) {

	data, err := t.repo.GetRoadDamageForImport(roadID, IDParent)
	if err != nil {
		logs.Error(err)
		return data, responses.NewAppErr(400, err.Error())
	}
	return data, nil
}

func (t *roadDamageUseCase) GetRoadDamageByIDParent(parentId int) (models.RoadDamage, error) {
	data, err := t.repo.GetRoadDamageByIDParent(parentId)
	if err != nil {
		logs.Error(err)
		return data, responses.NewAppErr(400, err.Error())
	}
	return data, nil
}

func (t *roadDamageUseCase) GetRoadById(roadId int) (models.Road, error) {
	data, err := t.repo.GetRoadById(roadId)
	if err != nil {
		logs.Error(err)
		return data, responses.NewAppErr(400, err.Error())
	}
	return data, nil
}

func (t *roadDamageUseCase) GetRoadGeom(roadId, laneNo int) (models.RoadGeom, error) {
	data, err := t.repo.GetRoadGeom(roadId, laneNo)
	if err != nil {
		logs.Error(err)
		return data, responses.NewAppErr(400, err.Error())
	}
	return data, nil
}

func (t *roadDamageUseCase) UpdateRoadDamage(id, userId int, req requests.RoadDamageImport) error {
	if err := t.repo.UpdateRoadDamage(id, userId, req); err != nil {
		logs.Error(err)
		return responses.NewAppErr(400, err.Error())
	}
	return nil
}

func (t *roadDamageUseCase) UpdateRoadDamageStatusT(id int, req requests.RoadDamageImport) error {
	if err := t.repo.UpdateRoadDamageStatusT(id, req); err != nil {
		logs.Error(err)
		return responses.NewAppErr(400, err.Error())
	}
	return nil
}

func (t *roadDamageUseCase) UpdateRoadDamageStatusI(id int, req requests.RoadDamageImport) error {
	if err := t.repo.UpdateRoadDamageStatusI(id, req); err != nil {
		logs.Error(err)
		return responses.NewAppErr(400, err.Error())
	}
	return nil
}

type Damage struct {
	Damage models.RoadDamage `json:"damage"`
}

type DamageOld struct {
	Damage models.RoadDamage `json:"damage"`
}

func (t *roadDamageUseCase) SetRoadDamageFromImport(roadID, idParent, userID int, csvPath, imgPath string, req requests.RoadDamageImport, rcData requests.RcData, damageFilename, imageFilename *multipart.FileHeader) (interface{}, error) {
	// roadDamageId := 0
	// roadDamageRangeId := 0
	filePath := ""
	resRoadDamageID := 0
	// resParentID := 0
	// imgPrefix := time.Now().Format("20060102150405") + "_"
	var roadRanges []requests.RoadRangeItems
	var roadDamages []requests.RoadDamageItem

	// start := 0.0
	// end := 0.0
	//cntContinue := []int{}
	roadDamageRangeKmStart := make(map[int]float64)
	roadDamageRangeKmEnd := make(map[int]float64)
	var DamageOld DamageOld
	if idParent != 0 {
		roaddamageOld, err := t.repo.GetRoadDamageByIDParent(idParent)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		DamageOld.Damage = roaddamageOld
	}

	if req.ImageFilenameStatus == "delete" {
		err := os.Remove(DamageOld.Damage.ImgFilepath)
		if err != nil {
			logs.Error(err)
		}
		err = t.repo.UpdateRoadDamageImgPath(DamageOld.Damage.Id, "")
		if err != nil {
			return responses.RoadConditionUpdate{}, fmt.Errorf("DeleteRoadConditionImgPath:", err.Error())
		}
	}
	// check have csv file
	// if damageFilename != nil {
	var damage Damage
	if req.DamageFilenameStatus == "upload" {
		// create folder
		if err := os.MkdirAll(os.Getenv("ROAD_DAMAGE_CSV_DIR"), 0775); err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(http.StatusUnauthorized, err.Error())
		}

		// get road
		roadInfo, err := t.repo.GetLastRoadInfoByID(roadID)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		direction := roadInfo.RefDirectionId

		csvArray, err := t.RoadDamageReadCSVFile(roadID, csvPath)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}

		switch direction {
		case 1:
			sort.Slice(csvArray, func(i, j int) bool {
				return csvArray[i].KM < csvArray[j].KM
			})
		case 2:
			sort.Slice(csvArray, func(i, j int) bool {
				return csvArray[i].KM > csvArray[j].KM
			})
		}

		rcData.RoadId = csvArray[0].RoadId
		// return csvArray, nil
		// get road info
		// roadInfo, err := t.repo.GetRoadInfo(rcData.RoadId)
		roadGeom, err := t.GetRoadGeom(rcData.RoadId, rcData.LaneNo)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		rcData.KmStart = float64(roadGeom.KmStart)
		rcData.KmEnd = float64(roadGeom.KmEnd)
		rcData.DamageInputFilepath = csvPath
		rcData.DamageImgFilepath = imgPath

		rangeItems := make(map[string]requests.MItem)
		var sumValues requests.SumValues
		//isDown := false

		var currentType string
		var currentCount int
		surveyTypeData := make(map[string][]requests.RoadDamageCsv)
		surfaceTemp := make(map[int]string)
		for index, csvData := range csvArray {
			surfaceTemp[index] = csvData.SurveyType
			if csvData.SurveyType != currentType {
				currentType = csvData.SurveyType
				currentCount++
			}

			newSurveyType := fmt.Sprintf("%s-%d", csvData.SurveyType, currentCount)
			surveyTypeData[newSurveyType] = append(surveyTypeData[newSurveyType], csvData)

		}

		ranges := make(map[string]string)

		dataTemp := make(map[string][]float64)
		kmType := make(map[float64]string)
		for index, item := range csvArray {
			kmType[(item.KM)] = item.SurveyType
			m_start := 0
			m_end := 0
			mIdx := ""
			if direction == 1 {
				m_start = (int(item.KM) / 100) * 100
				m_end = m_start + 100

				if index == 0 && len(csvArray) > 1 {
					if item.SurveyType != csvArray[1].SurveyType {
						m_start = (int(item.KM) / 100) * 100
						m_end = int(item.KM)
					}
				}

				if float64(m_start) == roadGeom.KmEnd {
					m_start = m_start - 100
					m_end = m_end - 100
				}
				mIdx = strconv.Itoa(m_start) + "_" + strconv.Itoa(m_end)
				dataTemp[mIdx] = append(dataTemp[mIdx], item.KM)
			} else if direction == 2 {
				m_end = (int(item.KM) / 100) * 100
				m_start = m_end + 100
				if float64(m_end) == roadGeom.KmStart {
					m_start = m_start - 100
					m_end = m_end - 100
				}
				mIdx = strconv.Itoa(m_start) + "_" + strconv.Itoa(m_end)
				dataTemp[mIdx] = append(dataTemp[mIdx], item.KM)
			}
		}
		// return dataTemp, nil
		/////////
		for index, item := range csvArray {
			m_start := 0
			m_end := 0
			mIdx := ""
			if direction == 1 {
				m_start = (int(item.KM) / 100) * 100
				m_end = m_start + 100

				if index == 0 && len(csvArray) > 1 {
					if item.SurveyType != csvArray[1].SurveyType {
						m_start = (int(item.KM) / 100) * 100
						m_end = int(item.KM)
					}
				}

				if float64(m_start) == roadGeom.KmEnd {
					m_start = m_start - 100
					m_end = m_end - 100
				}
				mIdx = strconv.Itoa(m_start) + "_" + strconv.Itoa(m_end)
			} else if direction == 2 {
				m_end = (int(item.KM) / 100) * 100
				m_start = m_end + 100
				if float64(m_end) == roadGeom.KmStart {
					m_start = m_start - 100
					m_end = m_end - 100
				}
				mIdx = strconv.Itoa(m_start) + "_" + strconv.Itoa(m_end)
			}

			switch direction {
			case 1:
				if item.KM < float64(roadGeom.KmStart) || item.KM > float64(roadGeom.KmEnd) {
					logs.Error(constants.INVALID_GEOM_RANGE)
					return "", responses.NewAppErr(400, constants.INVALID_GEOM_RANGE)

				}
			case 2:
				if item.KM > float64(roadGeom.KmStart) || item.KM < float64(roadGeom.KmEnd) {
					logs.Error(constants.INVALID_GEOM_RANGE)
					return "", responses.NewAppErr(400, constants.INVALID_GEOM_RANGE)
				}
			}
			point := math.Abs(float64(roadGeom.KmStart)-item.KM) / math.Abs(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))
			theGeomPoint := "ST_LineInterpolatePoint('" + roadGeom.TheGeom + "', " + fmt.Sprintf("%f", point) + ")"

			var roadDamage requests.RoadDamageItem
			roadDamage.KM = item.KM
			roadDamage.SurveyType = item.SurveyType
			roadDamage.AcIcrack = item.AcIcrack
			roadDamage.AcUcrack = item.AcUcrack
			roadDamage.AcRavelling = item.AcRavelling
			roadDamage.AcPatching = item.AcPatching
			roadDamage.AcPotholeArea = item.AcPotholeArea
			roadDamage.AcBleeding = item.AcBleeding
			roadDamage.AcPotholeCount = item.AcPotholeCount
			roadDamage.CcTransverseCrack = item.CcTransverseCrack
			roadDamage.CcNonTransverseCrack = item.CcNonTransverseCrack
			roadDamage.CcCornerBreak = item.CcCornerBreak
			roadDamage.CcJointSealDamage = item.CcJointSealDamage
			roadDamage.CcPatching = item.CcPatching
			roadDamage.CcSpalling = item.CcSpalling
			roadDamage.CcScaling = item.CcScaling
			roadDamage.ImgFilepath = item.ImgFilepath
			roadDamage.TheGeomPoint = theGeomPoint

			_, ok := rangeItems[mIdx]
			var mItem requests.MItem
			if !ok {
				m_start_r := 0
				m_end_r := 0
				if direction == 1 {
					m_start_r = (int(item.KM) / 100) * 100
					m_end_r = m_start_r + 100

					if float64(m_start_r) == roadGeom.KmEnd {
						m_start_r = m_start_r - 100
						m_end_r = m_end_r - 100
					}

					surfaceTypeData := make(map[string][]float64)
					temp := ""
					cnt := 0
					for _, i := range dataTemp[mIdx] {
						surfaceType := kmType[i]
						if temp == "" {
							temp = surfaceType
						}

						if temp != surfaceType {
							cnt++
						}
						key := fmt.Sprintf("_%d", cnt)
						surfaceTypeData[surfaceType+key] = append(surfaceTypeData[surfaceType+key], i)
						temp = surfaceType
					}
					if index == 0 && len(csvArray) > 1 {
						if item.SurveyType != csvArray[1].SurveyType {
							helpers.PrintlnJson("item.KM", item.KM)
							m_end_r = int(item.KM)
						}
					}
					if len(surfaceTypeData) > 1 {
						for key, ii := range surfaceTypeData {
							roadRange := FindRange1(ii, key, m_start_r, m_end_r, direction, rcData, mItem, roadGeom, mIdx, item, rangeItems)
							roadRanges = append(roadRanges, roadRange)

						}
					} else {
						isFixData := false
						if index == 0 && len(csvArray) > 1 {
							if item.SurveyType != csvArray[1].SurveyType {
								isFixData = true
							}
						}
						for key, _ := range surfaceTypeData {
							roadRange := FindRange2(isFixData, key, m_start_r, m_end_r, direction, rcData, mItem, roadGeom, mIdx, item, rangeItems)
							roadRanges = append(roadRanges, roadRange)
						}
					}
				} else if direction == 2 {
					m_end_r = (int(item.KM) / 100) * 100
					m_start_r = m_end_r + 100

					if float64(m_end_r) == roadGeom.KmStart {
						m_start_r = m_start_r - 100
						m_end_r = m_end_r - 100
					}

					surfaceTypeData := make(map[string][]float64)
					temp := ""
					cnt := 0
					for _, i := range dataTemp[mIdx] {
						surfaceType := kmType[i]
						if temp == "" {
							temp = surfaceType
						}

						if temp != surfaceType {
							cnt++
						}
						key := fmt.Sprintf("_%d", cnt)
						surfaceTypeData[surfaceType+key] = append(surfaceTypeData[surfaceType+key], i)
						temp = surfaceType
					}
					if len(surfaceTypeData) > 1 {
						for key, ii := range surfaceTypeData {
							roadRange := FindRange1(ii, key, m_start_r, m_end_r, direction, rcData, mItem, roadGeom, mIdx, item, rangeItems)
							roadRanges = append(roadRanges, roadRange)

						}
					} else {
						for key, _ := range surfaceTypeData {
							roadRange := FindRange2(false, key, m_start_r, m_end_r, direction, rcData, mItem, roadGeom, mIdx, item, rangeItems)
							roadRanges = append(roadRanges, roadRange)

						}
					}
				}
			}

			roadDamages = append(roadDamages, roadDamage)
			sumValues.SurveyType = item.SurveyType
			sumValues.AcIcrack += item.AcIcrack
			sumValues.AcUcrack += item.AcUcrack
			sumValues.AcRavelling += item.AcRavelling
			sumValues.AcPatching += item.AcPatching
			sumValues.AcPotholeArea += item.AcPotholeArea
			sumValues.AcBleeding += item.AcBleeding
			sumValues.AcPotholeCount += item.AcPotholeCount
			sumValues.CcTransverseCrack += item.CcTransverseCrack
			sumValues.CcNonTransverseCrack += item.CcNonTransverseCrack
			sumValues.CcCornerBreak += item.CcCornerBreak
			sumValues.CcJointSealDamage += item.CcJointSealDamage
			sumValues.CcPatching += item.CcPatching
			sumValues.CcSpalling += item.CcSpalling
			sumValues.CcScaling += item.CcScaling
		}
		// insert road damage
		roadDamage, err := t.repo.CreateRoadDamage(sumValues, rcData)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		resRoadDamageID = roadDamage.Id
		damage.Damage = roadDamage
		roadDamageId := roadDamage.Id
		if rcData.IDParent == 0 {
			filePath = fmt.Sprintf("storages/road/attachments/%d/damage/%d/", rcData.RoadId, roadDamageId)
			if err := t.repo.UpdateRoadDamageIDParent(roadDamageId); err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, err.Error())
			}
			roadDamage.IdParent = roadDamageId
			damage.Damage = roadDamage
		} else {
			filePath = fmt.Sprintf("storages/road/attachments/%d/damage/%d/", rcData.RoadId, roadDamageId)
			roadDamage.IdParent = idParent
			damage.Damage = roadDamage

		}

		if err := os.MkdirAll(filePath, 0775); err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
		}
		csvName := time.Now().Format("20060102150405") + "_" + damageFilename.Filename
		fullCsvName := filePath + csvName

		if err := t.repo.UpdateRoadDamageCsvFilepath(roadDamageId, fullCsvName); err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}

		// move file csv upload
		if err := helpers.MoveFile(csvPath, fullCsvName); err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}

		if direction == 1 {
			sort.Slice(roadRanges, func(i, j int) bool {
				return roadRanges[i].KmStart < roadRanges[j].KmEnd
			})
		} else {
			sort.Slice(roadRanges, func(i, j int) bool {
				return roadRanges[i].KmStart > roadRanges[j].KmEnd
			})
		}

		// return roadRanges, nil
		roadRangeDatas := RoadRangeDatas(roadRanges)
		// return roadRangeDatas, nil
		for i_r, rag := range roadRangeDatas {
			kmStart := rag.KmStart
			kmEnd := rag.KmEnd

			theGeom := rag.TheGeom
			var roadRangeItem requests.RaodRangeItem
			for _, item := range roadDamages {
				switch direction {
				case 1:
					if (item.KM >= kmStart && item.KM <= kmEnd) || (item.KM == kmEnd && i_r == len(roadDamages)-1) {
						if kmStart < rcData.KmStart {
							roadRangeItem.KMStart = rcData.KmStart
						} else {
							roadRangeItem.KMStart = kmStart
						}
						roadRangeItem.KMEnd = kmEnd
						roadRangeItem.SurveyType = rag.SurveyType //item.SurveyType
						roadRangeItem.AcIcrack += item.AcIcrack
						roadRangeItem.AcUcrack += item.AcUcrack
						roadRangeItem.AcRavelling += item.AcRavelling
						roadRangeItem.AcPatching += item.AcPatching
						roadRangeItem.AcPotholeArea += item.AcPotholeArea
						roadRangeItem.AcBleeding += item.AcBleeding
						roadRangeItem.AcPotholeCount += item.AcPotholeCount
						roadRangeItem.CcTransverseCrack += item.CcTransverseCrack
						roadRangeItem.CcNonTransverseCrack += item.CcNonTransverseCrack
						roadRangeItem.CcCornerBreak += item.CcCornerBreak
						roadRangeItem.CcJointSealDamage += item.CcJointSealDamage
						roadRangeItem.CcPatching += item.CcPatching
						roadRangeItem.CcSpalling += item.CcSpalling
						roadRangeItem.CcScaling += item.CcScaling
						roadRangeItem.TheGeom = theGeom
					}
				case 2:
					if (item.KM <= kmStart && item.KM >= kmEnd) || (item.KM == kmStart && i_r == 0) {
						if kmEnd < rcData.KmEnd {
							roadRangeItem.KMEnd = rcData.KmEnd
						} else {
							roadRangeItem.KMEnd = kmEnd
						}
						roadRangeItem.KMStart = kmStart
						roadRangeItem.SurveyType = item.SurveyType
						roadRangeItem.AcIcrack += item.AcIcrack
						roadRangeItem.AcUcrack += item.AcUcrack
						roadRangeItem.AcRavelling += item.AcRavelling
						roadRangeItem.AcPatching += item.AcPatching
						roadRangeItem.AcPotholeArea += item.AcPotholeArea
						roadRangeItem.AcBleeding += item.AcBleeding
						roadRangeItem.AcPotholeCount += item.AcPotholeCount
						roadRangeItem.CcTransverseCrack += item.CcTransverseCrack
						roadRangeItem.CcNonTransverseCrack += item.CcNonTransverseCrack
						roadRangeItem.CcCornerBreak += item.CcCornerBreak
						roadRangeItem.CcJointSealDamage += item.CcJointSealDamage
						roadRangeItem.CcPatching += item.CcPatching
						roadRangeItem.CcSpalling += item.CcSpalling
						roadRangeItem.CcScaling += item.CcScaling
						roadRangeItem.TheGeom = theGeom
					}
				}
			}
			if roadRangeItem.KMStart == 0 && roadRangeItem.KMEnd == 0 {
				continue
			}
			km := fmt.Sprintf("%d-%d", int(kmStart), int(kmEnd))
			_, ok := ranges[km]
			if ok {
				continue
			} else {
				ranges[km] = km
			}

			if roadRangeItem.KMStart == roadRangeItem.KMEnd {
				// continue
			}

			roadDamageRangeId, _, err := t.repo.CreateRoadDamageRange(roadRangeItem, rcData, roadDamageId)
			if err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, err.Error())
			}
			roadDamageRangeKmStart[roadDamageRangeId] = roadRangeItem.KMStart
			roadDamageRangeKmEnd[roadDamageRangeId] = roadRangeItem.KMEnd
			tempKM := make(map[float64]bool)
			for _, item := range roadDamages {
				switch direction {
				case 1:
					// fmt.Println(item.KM)
					if item.KM > roadRangeItem.KMStart && item.KM <= roadRangeItem.KMEnd {
						tempKM[item.KM] = true
						var roadDamageM requests.RoadDamageMItem
						roadDamageM.RoadDamageRangeID = roadDamageRangeId
						roadDamageM.KM = item.KM
						roadDamageM.TheGeomPoint = item.TheGeomPoint
						if item.ImgFilepath != "" {
							roadDamageM.ImgFilepath = filePath + item.ImgFilepath
						} else {
							roadDamageM.ImgFilepath = ""
						}
						roadDamageM.SurveyType = item.SurveyType
						roadDamageM.AcIcrack = item.AcIcrack
						roadDamageM.AcUcrack = item.AcUcrack
						roadDamageM.AcRavelling = item.AcRavelling
						roadDamageM.AcPatching = item.AcPatching
						roadDamageM.AcPotholeArea = item.AcPotholeArea
						roadDamageM.AcBleeding = item.AcBleeding
						roadDamageM.AcPotholeCount = item.AcPotholeCount
						roadDamageM.CcTransverseCrack = item.CcTransverseCrack
						roadDamageM.CcNonTransverseCrack = item.CcNonTransverseCrack
						roadDamageM.CcCornerBreak = item.CcCornerBreak
						roadDamageM.CcJointSealDamage = item.CcJointSealDamage
						roadDamageM.CcPatching = item.CcPatching
						roadDamageM.CcSpalling = item.CcSpalling
						roadDamageM.CcScaling = item.CcScaling
						err = t.repo.CreateRoadDamageM(roadDamageM, rcData)
						if err != nil {
							logs.Error(err)
							return "", responses.NewAppErr(400, "ไม่สามรถสร้าง Road Damage M")
						}
					} else {
						if (roadRangeItem.KMEnd == rcData.KmEnd && item.KM > roadRangeItem.KMStart) || (roadRangeItem.KMStart == rcData.KmStart && item.KM < roadRangeItem.KMEnd) {
							tempKM[item.KM] = true
							var roadDamageM requests.RoadDamageMItem
							roadDamageM.RoadDamageRangeID = roadDamageRangeId
							roadDamageM.KM = item.KM
							roadDamageM.TheGeomPoint = item.TheGeomPoint
							if item.ImgFilepath != "" {
								roadDamageM.ImgFilepath = filePath + item.ImgFilepath
							} else {
								roadDamageM.ImgFilepath = ""
							}
							roadDamageM.SurveyType = item.SurveyType
							roadDamageM.AcIcrack = item.AcIcrack
							roadDamageM.AcUcrack = item.AcUcrack
							roadDamageM.AcRavelling = item.AcRavelling
							roadDamageM.AcPatching = item.AcPatching
							roadDamageM.AcPotholeArea = item.AcPotholeArea
							roadDamageM.AcBleeding = item.AcBleeding
							roadDamageM.AcPotholeCount = item.AcPotholeCount
							roadDamageM.CcTransverseCrack = item.CcTransverseCrack
							roadDamageM.CcNonTransverseCrack = item.CcNonTransverseCrack
							roadDamageM.CcCornerBreak = item.CcCornerBreak
							roadDamageM.CcJointSealDamage = item.CcJointSealDamage
							roadDamageM.CcPatching = item.CcPatching
							roadDamageM.CcSpalling = item.CcSpalling
							roadDamageM.CcScaling = item.CcScaling
							err = t.repo.CreateRoadDamageM(roadDamageM, rcData)
							if err != nil {
								logs.Error(err)
								return "", responses.NewAppErr(400, "ไม่สามรถสร้าง Road Damage M")
							}
						} else {
							if item.KM == roadRangeItem.KMStart && !tempKM[item.KM] {
								tempKM[item.KM] = true
								var roadDamageM requests.RoadDamageMItem
								roadDamageM.RoadDamageRangeID = roadDamageRangeId
								roadDamageM.KM = item.KM
								roadDamageM.TheGeomPoint = item.TheGeomPoint
								if item.ImgFilepath != "" {
									roadDamageM.ImgFilepath = filePath + item.ImgFilepath
								} else {
									roadDamageM.ImgFilepath = ""
								}
								roadDamageM.SurveyType = item.SurveyType
								roadDamageM.AcIcrack = item.AcIcrack
								roadDamageM.AcUcrack = item.AcUcrack
								roadDamageM.AcRavelling = item.AcRavelling
								roadDamageM.AcPatching = item.AcPatching
								roadDamageM.AcPotholeArea = item.AcPotholeArea
								roadDamageM.AcBleeding = item.AcBleeding
								roadDamageM.AcPotholeCount = item.AcPotholeCount
								roadDamageM.CcTransverseCrack = item.CcTransverseCrack
								roadDamageM.CcNonTransverseCrack = item.CcNonTransverseCrack
								roadDamageM.CcCornerBreak = item.CcCornerBreak
								roadDamageM.CcJointSealDamage = item.CcJointSealDamage
								roadDamageM.CcPatching = item.CcPatching
								roadDamageM.CcSpalling = item.CcSpalling
								roadDamageM.CcScaling = item.CcScaling
								err = t.repo.CreateRoadDamageM(roadDamageM, rcData)
								if err != nil {
									logs.Error(err)
									return "", responses.NewAppErr(400, "ไม่สามรถสร้าง Road Damage M")
								}
							}
						}
					}
				case 2:
					if (item.KM < roadRangeItem.KMStart && item.KM >= roadRangeItem.KMEnd) || (item.KM == roadGeom.KmStart && item.KM == roadRangeItem.KMStart) {
						var roadDamageM requests.RoadDamageMItem
						roadDamageM.RoadDamageRangeID = roadDamageRangeId
						roadDamageM.KM = item.KM
						roadDamageM.TheGeomPoint = item.TheGeomPoint
						if item.ImgFilepath != "" {
							roadDamageM.ImgFilepath = filePath + item.ImgFilepath
						} else {
							roadDamageM.ImgFilepath = ""
						}
						tempKM[item.KM] = true
						roadDamageM.SurveyType = item.SurveyType
						roadDamageM.AcIcrack = item.AcIcrack
						roadDamageM.AcUcrack = item.AcUcrack
						roadDamageM.AcRavelling = item.AcRavelling
						roadDamageM.AcPatching = item.AcPatching
						roadDamageM.AcPotholeArea = item.AcPotholeArea
						roadDamageM.AcBleeding = item.AcBleeding
						roadDamageM.AcPotholeCount = item.AcPotholeCount
						roadDamageM.CcTransverseCrack = item.CcTransverseCrack
						roadDamageM.CcNonTransverseCrack = item.CcNonTransverseCrack
						roadDamageM.CcCornerBreak = item.CcCornerBreak
						roadDamageM.CcJointSealDamage = item.CcJointSealDamage
						roadDamageM.CcPatching = item.CcPatching
						roadDamageM.CcSpalling = item.CcSpalling
						roadDamageM.CcScaling = item.CcScaling
						err = t.repo.CreateRoadDamageM(roadDamageM, rcData)
						if err != nil {
							logs.Error(err)
							return "", responses.NewAppErr(400, "ไม่สามรถสร้าง Road Damage M")
						}
					} else {
						if len(roadRanges)-1 == i_r {
							if item.KM < rcData.KmEnd+100 && item.KM <= rcData.KmEnd {
								var roadDamageM requests.RoadDamageMItem
								roadDamageM.RoadDamageRangeID = roadDamageRangeId
								roadDamageM.KM = item.KM
								roadDamageM.TheGeomPoint = item.TheGeomPoint
								if item.ImgFilepath != "" {
									roadDamageM.ImgFilepath = filePath + item.ImgFilepath
								} else {
									roadDamageM.ImgFilepath = ""
								}
								tempKM[item.KM] = true
								roadDamageM.SurveyType = item.SurveyType
								roadDamageM.AcIcrack = item.AcIcrack
								roadDamageM.AcUcrack = item.AcUcrack
								roadDamageM.AcRavelling = item.AcRavelling
								roadDamageM.AcPatching = item.AcPatching
								roadDamageM.AcPotholeArea = item.AcPotholeArea
								roadDamageM.AcBleeding = item.AcBleeding
								roadDamageM.AcPotholeCount = item.AcPotholeCount
								roadDamageM.CcTransverseCrack = item.CcTransverseCrack
								roadDamageM.CcNonTransverseCrack = item.CcNonTransverseCrack
								roadDamageM.CcCornerBreak = item.CcCornerBreak
								roadDamageM.CcJointSealDamage = item.CcJointSealDamage
								roadDamageM.CcPatching = item.CcPatching
								roadDamageM.CcSpalling = item.CcSpalling
								roadDamageM.CcScaling = item.CcScaling
								err = t.repo.CreateRoadDamageM(roadDamageM, rcData)
								if err != nil {
									logs.Error(err)
									return "", responses.NewAppErr(400, "ไม่สามรถสร้าง Road Damage M")
								}
							}
						} else {
							// else {
							if item.KM == roadRangeItem.KMStart && !tempKM[item.KM] {
								var roadDamageM requests.RoadDamageMItem
								roadDamageM.RoadDamageRangeID = roadDamageRangeId
								roadDamageM.KM = item.KM
								roadDamageM.TheGeomPoint = item.TheGeomPoint
								if item.ImgFilepath != "" {
									roadDamageM.ImgFilepath = filePath + item.ImgFilepath
								} else {
									roadDamageM.ImgFilepath = ""
								}
								tempKM[item.KM] = true
								roadDamageM.SurveyType = item.SurveyType
								roadDamageM.AcIcrack = item.AcIcrack
								roadDamageM.AcUcrack = item.AcUcrack
								roadDamageM.AcRavelling = item.AcRavelling
								roadDamageM.AcPatching = item.AcPatching
								roadDamageM.AcPotholeArea = item.AcPotholeArea
								roadDamageM.AcBleeding = item.AcBleeding
								roadDamageM.AcPotholeCount = item.AcPotholeCount
								roadDamageM.CcTransverseCrack = item.CcTransverseCrack
								roadDamageM.CcNonTransverseCrack = item.CcNonTransverseCrack
								roadDamageM.CcCornerBreak = item.CcCornerBreak
								roadDamageM.CcJointSealDamage = item.CcJointSealDamage
								roadDamageM.CcPatching = item.CcPatching
								roadDamageM.CcSpalling = item.CcSpalling
								roadDamageM.CcScaling = item.CcScaling
								err = t.repo.CreateRoadDamageM(roadDamageM, rcData)
								if err != nil {
									logs.Error(err)
									return "", responses.NewAppErr(400, "ไม่สามรถสร้าง Road Damage M")
								}
							}
							// }
						}
					}
				}
			}
			if err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, "ไม่สามรถสร้าง Road Damage Range")
			}
		}
	}

	damageId := 0
	if idParent == 0 {
		// helpers.PrintlnJson("roadDamageId", damage.Damage.IdParent)
		damage, err := t.GetRoadDamageByIDParent(damage.Damage.IdParent)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		filePath = fmt.Sprintf("storages/road/attachments/%d/damage/%d/", roadID, damage.Id)
		// roadDamageId = damage.Damage.
		damageId = damage.Id
	} else {
		damage, err := t.GetRoadDamageByIDParent(idParent)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		filePath = fmt.Sprintf("storages/road/attachments/%d/damage/%d/", roadID, damage.Id)
		damageId = damage.Id
		// helpers.PrintlnJson(damage)
	}

	// Image Filename
	// if imageFilename != nil {
	if req.ImageFilenameStatus == "upload" {

		// if imageFilename != nil && damageFilename == nil {

		// }
		// helpers.PrintlnJson(filePath)
		if err := os.MkdirAll(filePath, 0775); err != nil {
			// helpers.PrintlnJson(filePath)
			fmt.Println("111")
			logs.Error(constants.FAILED_TO_SAVE_FILE)
			return "", responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
		}
		allowedTypes, _ := helpers.TypeFileAllowed(imageFilename.Filename, "zip|rar")
		if !allowedTypes {
			logs.Error(constants.UNSUPPORTED_FILE_TYPE)
			return "", responses.NewAppErr(400, constants.UNSUPPORTED_FILE_TYPE)
		}
		// helpers.PrintlnJson(imgPath, filePath)
		err := helpers.Unzip(imgPath, filePath)
		if err != nil {
			fmt.Println("222")
			fmt.Println(err)
			logs.Error(constants.FAILED_TO_SAVE_FILE)
			return "", responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
		}
		// move file zip or rar
		err = os.Rename(imgPath, filePath+imageFilename.Filename)
		if err != nil {
			fmt.Println("333")
			logs.Error(err.Error())
			return "", responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
		}

		if err := t.repo.UpdateRoadDamageImgFilepath(damageId, filePath+imageFilename.Filename); err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}

		// helpers.PrintlnJson("test_1", filePath+imageFilename.Filename)
	} else {
		if req.ImageFilenameStatus == "not_edit" {
			// // helpers.PrintlnJson("imgPath", imgPath, filePath+DamageOld.Damage.ImgFilepath, damageId)
			// filePath = fmt.Sprintf("storages/road/attachments/%d/damage/%d/", roadID, damageId)

			// newPaths := helpers.Explode("/", DamageOld.Damage.ImgFilepath)
			// var paths []string
			// fileName := ""
			// for i, item := range newPaths {
			// 	if i < len(newPaths)-1 {
			// 		paths = append(paths, item)
			// 	} else {
			// 		fileName = item
			// 	}

			// }
			// if err := helpers.CopyFile(DamageOld.Damage.ImgFilepath, filePath+fileName); err != nil {
			// 	// Handle the error accordingly.
			// 	panic(err)
			// }
			// err := helpers.Unzip(filePath+fileName, filePath)
			// if err != nil {
			// 	fmt.Println(err)
			// 	logs.Error(constants.FAILED_TO_SAVE_FILE)
			// 	// return "", responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
			// }
			// // err := os.Rename(newPath, filePath)
			// // if err != nil {
			// // 	helpers.PrintlnJson("filePath+newPath, filePath", newPath+"/", filePath)
			// // 	logs.Error(err.Error())
			// // 	return "", responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
			// // }

			// if err := t.repo.UpdateRoadDamageImgFilepath(damageId, DamageOld.Damage.ImgFilepath); err != nil {
			// 	logs.Error(err)
			// 	return "", responses.NewAppErr(400, err.Error())
			// }
		} else {
			if err := t.repo.UpdateRoadDamageImgFilepath(damageId, ""); err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, err.Error())
			}
		}

		// helpers.PrintlnJson("test_2", DamageOld)
		// if err := t.repo.UpdateRoadDamageImgFilepath(roadDamageId, ""); err != nil {
		// 	logs.Error(err)
		// 	return "", responses.NewAppErr(400, err.Error())
		// }
	}

	var data responses.RoadDamageSetData
	data.ID = resRoadDamageID
	if rcData.IDParent != 0 {
		data.IDParent = rcData.IDParent
	} else {
		data.IDParent = resRoadDamageID
	}

	return data, nil
}

type Permission struct {
	CanEdit    bool `json:"can_edit"`
	CanDelete  bool `json:"can_delete"`
	CanApprove bool `json:"can_approve"`
	CanSend    bool `json:"can_send"`
	CanReject  bool `json:"can_reject"`
}

func (t *roadDamageUseCase) GetRoadDamageDetail(roadId, idParent int, permissions []string) (interface{}, error) {
	var result responses.RoadDamageDetailRespond
	directionId, err := t.repo.GetRoadDamageDirection(roadId, idParent, permissions)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	data, _ := t.repo.GetRoadDamageDetail(roadId, idParent, directionId, permissions)
	user, err := t.repo.GetUserById(data.UpdatedBy)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	copier.Copy(&result, &data)
	result.SurveyedDate = (data.SurveyedDate)
	result.CreatedDate = helpers.SetTimeToString(data.CreatedDate)
	result.UpdatedDate = helpers.SetTimeToString(data.UpdatedDate)

	STORAGE_IP := os.Getenv("STORAGE_IP")
	for i := range result.RoadDamageRange {
		if result.ImgFilepath == "" || result.ImgFilepath == STORAGE_IP+"/" {
			for j := range result.RoadDamageRange[i].RoadDamageM {
				result.RoadDamageRange[i].RoadDamageM[j].ImgFilepath = ""
			}
		} else {
			for j := range result.RoadDamageRange[i].RoadDamageM {
				result.RoadDamageRange[i].RoadDamageM[j].ImgFilepath = STORAGE_IP + "/" + result.RoadDamageRange[i].RoadDamageM[j].ImgFilepath
			}
		}
	}

	if directionId == 1 {
		// Sort the RoadDamageRange slice by 'km_start' field in ascending order
		sort.Slice(result.RoadDamageRange, func(i, j int) bool {
			return result.RoadDamageRange[i].KmStart < result.RoadDamageRange[j].KmStart
		})

		//Sort each RoadDamageM slice by 'km' field in ascending order
		for i := range result.RoadDamageRange {

			sort.Slice(result.RoadDamageRange[i].RoadDamageM, func(j, k int) bool {
				return result.RoadDamageRange[i].RoadDamageM[j].Km < result.RoadDamageRange[i].RoadDamageM[k].Km
			})

		}
	} else {
		sort.Slice(result.RoadDamageRange, func(i, j int) bool {
			return result.RoadDamageRange[i].KmStart > result.RoadDamageRange[j].KmStart
		})

		//Sort each RoadDamageM slice by 'km' field in ascending order
		for i := range result.RoadDamageRange {

			sort.Slice(result.RoadDamageRange[i].RoadDamageM, func(j, k int) bool {
				return result.RoadDamageRange[i].RoadDamageM[j].Km > result.RoadDamageRange[i].RoadDamageM[k].Km
			})

		}

	}
	// for i := range data.RoadDamageRange {
	// 	for j, v := range data.RoadDamageRange[i].RoadDamageM {
	// 		result.RoadDamageRange[i].RoadDamageM[j].CreatedDate = helpers.SetTimeToString(v.CreatedDate)
	// 		result.RoadDamageRange[i].RoadDamageM[j].UpdatedDate = helpers.SetTimeToString(v.UpdatedDate)
	// 	}
	// }

	laneNo := data.LaneNo
	roadGeom, err := t.repo.GetRoadGeomByRoadIDByLaneNo(roadId, laneNo)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	roadDamageDetail := responses.RoadDamageListDetail{
		Id:          data.Id,
		IdParent:    data.IdParent,
		UpdatedDate: helpers.SetTimeToString(data.UpdatedDate),
		Status:      data.RoadDamageStatus.Name,
		UpdatedBy:   user,
		RoadDamage:  result,
		TheGeom:     roadGeom.TheGeom,
	}
	return roadDamageDetail, nil
}

func (t *roadDamageUseCase) GetRoadDamageTemplate(roadID int) (interface{}, error) {
	road, err := t.repo.GetRoadDamageTemplate(roadID)
	if err != nil {
		return "", responses.NewInternalServerError()
	}
	roadCode := road.RoadSection.RoadGroup.Number + road.RoadSection.Number
	header := []string{"road_id", "road_code", "name", "km", "survey_type", "ac_icrack", "ac_ucrack", "ac_ravelling", "ac_patching", "ac_pothole_area", "ac_bleeding", "ac_pothole_count", "cc_transverse_crack", "cc_non_transverse_crack", "cc_corner_break", "cc_joint_seal_damage", "cc_patching", "cc_spalling", "cc_scaling", "img_filepath"}
	filePath := "storages/template_damage/"
	csvFiles := []string{}

	for _, item := range road.RoadGeom {

		surface, err := t.repo.GetRoadSurfaceLane(road.Id, item.LaneNo, item.KmStart)
		if err != nil {
			return "", responses.NewInternalServerError()
		}

		// return surface, nil

		data := [][]string{}
		start := fmt.Sprintf("%d", item.RoadId) + "," + roadCode + "," + road.RoadInfo.Name + "," + fmt.Sprintf("%.0f", item.KmStart) + "," + surface + ",0,0,0,0,0,0,0,0,0,0,0,0,0,0"

		end := fmt.Sprintf("%d", item.RoadId) + "," + roadCode + "," + road.RoadInfo.Name + "," + fmt.Sprintf("%.0f", item.KmEnd) + "," + surface + ",0,0,0,0,0,0,0,0,0,0,0,0,0,0"
		data = append(data, strings.Split(start, ","))
		data = append(data, strings.Split(end, ","))
		fileName := "damage_template_" + roadCode + "_lane" + fmt.Sprintf("%d", item.LaneNo) + ".csv"
		err = CsvGenerateFile(header, data, filePath, fileName)
		if err != nil {
			return "", responses.NewAppErr(http.StatusBadRequest, err.Error())
		}
		csvFiles = append(csvFiles, filePath+fileName)
	}
	zipName := filePath + "template_damage_" + time.Now().Format("2006-01-02_150405") + ".zip"
	err = helpers.Zip(zipName, csvFiles)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(http.StatusBadRequest, err.Error())
	}
	return os.Getenv("STORAGE_IP") + "/" + zipName, nil
}

func (t *roadDamageUseCase) DeleteRoadDamageForImport(idParent, userID int) (interface{}, error) {
	err := t.repo.DeleteRoadDamageForImport(idParent, userID)
	if err != nil {
		logs.Error(err)
		return "", errors.New(err.Error())
	}
	return "", nil
}

// ////////////////////// Private function ////////////////////////
func (t *roadDamageUseCase) RoadDamageReadCSVFile(roadID int, filePath string) ([]requests.RoadDamageCsv, error) {
	isFirstRow := true
	headerMap := make(map[string]int)
	f, _ := os.Open(filePath)
	r := csv.NewReader(f)
	var roadDamageCsvs []requests.RoadDamageCsv
	line := 2
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		// helpers.PrintlnJson("record", len(record))
		if len(record) < 19 {
			return []requests.RoadDamageCsv{}, errors.New("ไม่สามารถอัพโหลดไฟล์ได้")
		}

		// Handle first row case
		if isFirstRow {
			isFirstRow = false
			for i, v := range record {
				headerMap[v] = i
			}
			// Skip next code
			continue
		}

		if helpers.StrToInt(record[headerMap["road_id"]]) <= 0 || !helpers.CheckInt(record[headerMap["road_id"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล road_id บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if roadID != helpers.StrToInt(record[headerMap["road_id"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล road_id บรรทัดที่ %v ไม่ตรงกับสายทางที่เลือก", line)
		}

		if record[headerMap["road_code"]] == "" {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล road_code บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if record[headerMap["name"]] == "" {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล name บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if record[headerMap["survey_type"]] == "" {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล survey_type บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToInt(record[headerMap["km"]]) < 0 || !helpers.CheckInt(record[headerMap["km"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล km บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToFloat(record[headerMap["ac_icrack"]]) < 0 || !helpers.CheckFloat(record[headerMap["ac_icrack"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล ac_icrack บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToFloat(record[headerMap["ac_ucrack"]]) < 0 || !helpers.CheckFloat(record[headerMap["ac_ucrack"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล ac_ucrack บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToFloat(record[headerMap["ac_ravelling"]]) < 0 || !helpers.CheckFloat(record[headerMap["ac_ravelling"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล ac_ravelling บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToFloat(record[headerMap["ac_patching"]]) < 0 || !helpers.CheckFloat(record[headerMap["ac_patching"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล ac_patching บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToFloat(record[headerMap["ac_pothole_area"]]) < 0 || !helpers.CheckFloat(record[headerMap["ac_pothole_area"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล ac_pothole_area บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToFloat(record[headerMap["ac_bleeding"]]) < 0 || !helpers.CheckFloat(record[headerMap["ac_bleeding"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล ac_bleeding บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToFloat(record[headerMap["ac_pothole_count"]]) < 0 || !helpers.CheckFloat(record[headerMap["ac_pothole_count"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล ac_pothole_count บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToFloat(record[headerMap["cc_transverse_crack"]]) < 0 || !helpers.CheckFloat(record[headerMap["cc_transverse_crack"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล cc_transverse_crack บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToFloat(record[headerMap["cc_non_transverse_crack"]]) < 0 || !helpers.CheckFloat(record[headerMap["cc_non_transverse_crack"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล cc_non_transverse_crack บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToFloat(record[headerMap["cc_corner_break"]]) < 0 || !helpers.CheckFloat(record[headerMap["cc_corner_break"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล cc_corner_break บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToFloat(record[headerMap["cc_joint_seal_damage"]]) < 0 || !helpers.CheckFloat(record[headerMap["cc_joint_seal_damage"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล cc_joint_seal_damage บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToFloat(record[headerMap["cc_patching"]]) < 0 || !helpers.CheckFloat(record[headerMap["cc_patching"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล cc_patching บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToFloat(record[headerMap["cc_spalling"]]) < 0 || !helpers.CheckFloat(record[headerMap["cc_spalling"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล cc_spalling บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}

		if helpers.StrToFloat(record[headerMap["cc_scaling"]]) < 0 || !helpers.CheckFloat(record[headerMap["cc_scaling"]]) {
			return []requests.RoadDamageCsv{}, fmt.Errorf("ข้อมูล cc_scaling บรรทัดที่ %v ไม่ถูกต้องโปรดตรวจสอบข้อมูล", line)
		}
		line++
		imgFilepath := ""
		if len(record) == 19 {
			imgFilepath = ""
		} else {
			imgFilepath = record[headerMap["img_filepath"]]
		}
		roadDamageCsvs = append(roadDamageCsvs, requests.RoadDamageCsv{

			RoadId:               helpers.StrToInt(record[headerMap["road_id"]]),
			RoadCode:             record[headerMap["road_code"]],
			Name:                 record[headerMap["name"]],
			KM:                   helpers.StrToFloat(record[headerMap["km"]]),
			SurveyType:           record[headerMap["survey_type"]],
			AcIcrack:             helpers.StrToFloat(record[headerMap["ac_icrack"]]),
			AcUcrack:             helpers.StrToFloat(record[headerMap["ac_ucrack"]]),
			AcRavelling:          helpers.StrToFloat(record[headerMap["ac_ravelling"]]),
			AcPatching:           helpers.StrToFloat(record[headerMap["ac_patching"]]),
			AcPotholeArea:        helpers.StrToFloat(record[headerMap["ac_pothole_area"]]),
			AcBleeding:           helpers.StrToFloat(record[headerMap["ac_bleeding"]]),
			AcPotholeCount:       helpers.StrToFloat(record[headerMap["ac_pothole_count"]]),
			CcTransverseCrack:    helpers.StrToFloat(record[headerMap["cc_transverse_crack"]]),
			CcNonTransverseCrack: helpers.StrToFloat(record[headerMap["cc_non_transverse_crack"]]),
			CcCornerBreak:        helpers.StrToFloat(record[headerMap["cc_corner_break"]]),
			CcJointSealDamage:    helpers.StrToFloat(record[headerMap["cc_joint_seal_damage"]]),
			CcPatching:           helpers.StrToFloat(record[headerMap["cc_patching"]]),
			CcSpalling:           helpers.StrToFloat(record[headerMap["cc_spalling"]]),
			CcScaling:            helpers.StrToFloat(record[headerMap["cc_scaling"]]),
			ImgFilepath:          imgFilepath,
		})
	}
	return roadDamageCsvs, nil
}

func CsvGenerateFile(header []string, data [][]string, filePath, fileName string) error {
	// Create directory
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return err
	}

	// Create a new CSV file
	file, err := os.Create(filePath + fileName)
	if err != nil {
		logs.Error(err)
		return errors.New("cannot create file")
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Set the CSV header
	err = writer.Write(header)
	if err != nil {
		logs.Error(err)
		return errors.New("cannot write header")
	}

	for _, row := range data {
		err = writer.Write(row)
		if err != nil {
			logs.Error(err)
			return errors.New("cannot write row")
		}
	}
	return nil
}
