package usecases

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
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
	"gorm.io/gorm"
)

func (u *UseCase) GetPrepareData(analysisID int, req requests.MaintenanceAnalysis) (interface{}, error) {
	roadIDs := req.Roads
	lanes := []int{}

	if req.LaneTypeId == 0 {
		maxLane, err := u.Repo.GetMaxLane(req.Roads)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		for i := 1; i <= maxLane; i++ {
			lanes = append(lanes, i)
		}
	} else {
		lanes = append(lanes, req.LaneTypeId)
	}

	// laneNo := 1
	// laneNo := helpers.StrToInt(lane_no[0])
	var roads []responses.Road
	var roadDateBegins []responses.RoadDateBegins
	var roadGeoms []responses.RoadGeoms
	roadSurfaces := make(map[string][]models.RoadSurfacePrepareData)
	maintenances := make(map[int][]responses.MaintenanceData)
	growthRates := make(map[int]models.SettingAadtGrowthRate)
	aadtParameters := make(map[int]models.AadtParameterData)
	for _, roadID := range roadIDs {
		for _, lane := range lanes {
			laneNo := lane
			_, aadtParams, _ := u.GetSettingAadtParams(roadID)
			var settingTrafficParameter []models.SettingAADTParameter
			jsonData, err := json.Marshal(aadtParams["aadt_parameter"])
			if err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, err.Error())
			}

			err = json.Unmarshal([]byte(jsonData), &settingTrafficParameter)
			if err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, err.Error())
			}

			trafficParameter := models.SettingAADTParameter{}
			for _, item := range settingTrafficParameter {
				if item.RoadGroupID == int32(roadID) {
					copier.Copy(&trafficParameter, item)
				}
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

			roadGeomData, err := u.GetRoadGeomLaneByID(roadID, laneNo)
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					logs.Error(err)
					return "", responses.NewAppErr(400, err.Error())
				}
				logs.Error(err)
				continue
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
					// return road, nil
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
					maintenanceData.IcResult.RefSurface.LayerCoefficient = data["layer_coefficient"].(float64)
					maintenanceData.IcResult.RefSurface.Drainage = data["drainage"].(float64)
					maintenanceData.IcResult.RefSurface.A = data["a"].(float64)
					maintenanceData.IcResult.RefSurface.B = data["b"].(float64)
					cData := strings.Split(data["c"].(string), "^")
					if len(cData) > 1 {
						cbase, _ := strconv.ParseFloat(cData[0], 64)
						cexp, _ := strconv.ParseFloat(cData[1], 64)
						maintenanceData.IcResult.RefSurface.Cbase = cbase
						maintenanceData.IcResult.RefSurface.Cexp = cexp
					} else {
						maintenanceData.IcResult.RefSurface.Cbase = 0
						maintenanceData.IcResult.RefSurface.Cexp = 0
					}
					crt := 0.0
					if data["crt"] != nil {
						crt = data["crt"].(float64)
					}

					rrf := 0.0
					if data["rrf"] != nil {
						rrf = data["rrf"].(float64)
					}

					raveling_init_a0 := 0.0
					if data["raveling_init_a0"] != nil {
						raveling_init_a0 = data["raveling_init_a0"].(float64)
					}

					raveling_init_a1 := 0.0
					if data["raveling_init_a1"] != nil {
						raveling_init_a1 = data["raveling_init_a1"].(float64)
					}

					raveling_progress_a0 := 0.0
					if data["raveling_progress_a0"] != nil {
						raveling_progress_a0 = data["raveling_progress_a0"].(float64)
					}

					raveling_progress_a1 := 0.0
					if data["raveling_progress_a1"] != nil {
						raveling_progress_a1 = data["raveling_progress_a1"].(float64)
					}

					raveling_progress_a2 := 0.0
					if data["raveling_progress_a2"] != nil {
						raveling_progress_a2 = data["raveling_progress_a2"].(float64)
					}

					asc_init_hsold_equal_0_a0 := 0.0
					if data["asc_init_hsold_equal_0_a0"] != nil {
						asc_init_hsold_equal_0_a0 = data["asc_init_hsold_equal_0_a0"].(float64)
					}

					asc_init_hsold_equal_0_a1 := 0.0
					if data["asc_init_hsold_equal_0_a1"] != nil {
						asc_init_hsold_equal_0_a1 = data["asc_init_hsold_equal_0_a1"].(float64)
					}

					asc_init_hsold_equal_0_a2 := 0.0
					if data["asc_init_hsold_equal_0_a2"] != nil {
						asc_init_hsold_equal_0_a2 = data["asc_init_hsold_equal_0_a2"].(float64)
					}

					asc_init_hsold_equal_0_a3 := 0.0
					if data["asc_init_hsold_equal_0_a3"] != nil {
						asc_init_hsold_equal_0_a3 = data["asc_init_hsold_equal_0_a3"].(float64)
					}

					asc_init_hsold_equal_0_a4 := 0.0
					if data["asc_init_hsold_equal_0_a4"] != nil {
						asc_init_hsold_equal_0_a4 = data["asc_init_hsold_equal_0_a4"].(float64)
					}

					asc_init_hsold_over_0_a0 := 0.0
					if data["asc_init_hsold_over_0_a0"] != nil {
						asc_init_hsold_over_0_a0 = data["asc_init_hsold_over_0_a0"].(float64)
					}

					asc_init_hsold_over_0_a1 := 0.0
					if data["asc_init_hsold_over_0_a1"] != nil {
						asc_init_hsold_over_0_a1 = data["asc_init_hsold_over_0_a1"].(float64)
					}

					asc_init_hsold_over_0_a2 := 0.0
					if data["asc_init_hsold_over_0_a2"] != nil {
						asc_init_hsold_over_0_a2 = data["asc_init_hsold_over_0_a2"].(float64)
					}

					asc_init_hsold_over_0_a3 := 0.0
					if data["asc_init_hsold_over_0_a3"] != nil {
						asc_init_hsold_over_0_a3 = data["asc_init_hsold_over_0_a3"].(float64)
					}

					asc_init_hsold_over_0_a4 := 0.0
					if data["asc_init_hsold_over_0_a4"] != nil {
						asc_init_hsold_over_0_a4 = data["asc_init_hsold_over_0_a4"].(float64)
					}

					asc_progress_hsold_equal_0_a0 := 0.0
					if data["asc_progress_hsold_equal_0_a0"] != nil {
						asc_init_hsold_over_0_a4 = data["asc_progress_hsold_equal_0_a0"].(float64)
					}

					asc_progress_hsold_equal_0_a1 := 0.0
					if data["asc_progress_hsold_equal_0_a1"] != nil {
						asc_progress_hsold_equal_0_a1 = data["asc_progress_hsold_equal_0_a1"].(float64)
					}

					asc_progress_hsold_over_0_a0 := 0.0
					if data["asc_progress_hsold_over_0_a0"] != nil {
						asc_progress_hsold_over_0_a0 = data["asc_progress_hsold_over_0_a0"].(float64)
					}

					asc_progress_hsold_over_0_a1 := 0.0
					if data["asc_progress_hsold_over_0_a1"] != nil {
						asc_progress_hsold_over_0_a1 = data["asc_progress_hsold_over_0_a1"].(float64)
					}

					wsc_init_hsold_equal_0_a0 := 0.0
					if data["wsc_init_hsold_equal_0_a0"] != nil {
						wsc_init_hsold_equal_0_a0 = data["wsc_init_hsold_equal_0_a0"].(float64)
					}

					wsc_init_hsold_equal_0_a1 := 0.0
					if data["wsc_init_hsold_equal_0_a1"] != nil {
						wsc_init_hsold_equal_0_a1 = data["wsc_init_hsold_equal_0_a1"].(float64)
					}

					wsc_init_hsold_equal_0_a2 := 0.0
					if data["wsc_init_hsold_equal_0_a2"] != nil {
						wsc_init_hsold_equal_0_a2 = data["wsc_init_hsold_equal_0_a2"].(float64)
					}

					wsc_init_hsold_over_0_a0 := 0.0
					if data["wsc_init_hsold_over_0_a0"] != nil {
						wsc_init_hsold_over_0_a0 = data["wsc_init_hsold_over_0_a0"].(float64)
					}

					wsc_init_hsold_over_0_a1 := 0.0
					if data["wsc_init_hsold_over_0_a1"] != nil {
						wsc_init_hsold_over_0_a1 = data["wsc_init_hsold_over_0_a1"].(float64)
					}

					wsc_init_hsold_over_0_a2 := 0.0
					if data["wsc_init_hsold_over_0_a2"] != nil {
						wsc_init_hsold_over_0_a2 = data["wsc_init_hsold_over_0_a2"].(float64)
					}

					wsc_progress_hsold_equal_0_a0 := 0.0
					if data["wsc_progress_hsold_equal_0_a0"] != nil {
						wsc_progress_hsold_equal_0_a0 = data["wsc_progress_hsold_equal_0_a0"].(float64)
					}

					wsc_progress_hsold_equal_0_a1 := 0.0
					if data["wsc_progress_hsold_equal_0_a1"] != nil {
						wsc_progress_hsold_equal_0_a1 = data["wsc_progress_hsold_equal_0_a1"].(float64)
					}

					wsc_progress_hsold_over_0_a0 := 0.0
					if data["wsc_progress_hsold_over_0_a0"] != nil {
						wsc_progress_hsold_over_0_a0 = data["wsc_progress_hsold_over_0_a0"].(float64)
					}

					wsc_progress_hsold_over_0_a1 := 0.0
					if data["wsc_progress_hsold_over_0_a1"] != nil {
						wsc_progress_hsold_over_0_a1 = data["wsc_progress_hsold_over_0_a1"].(float64)
					}

					rpd_a0 := 0.0
					if data["rpd_a0"] != nil {
						rpd_a0 = data["rpd_a0"].(float64)
					}

					rpd_a1 := 0.0
					if data["rpd_a1"] != nil {
						rpd_a1 = data["rpd_a1"].(float64)
					}

					rpd_a2 := 0.0
					if data["rpd_a2"] != nil {
						rpd_a2 = data["rpd_a2"].(float64)
					}

					maintenanceData.IcResult.RefSurface.Crt = crt
					maintenanceData.IcResult.RefSurface.Rrf = rrf
					maintenanceData.IcResult.RefSurface.Raveling.Initial.A0 = raveling_init_a0
					maintenanceData.IcResult.RefSurface.Raveling.Initial.A1 = raveling_init_a1
					maintenanceData.IcResult.RefSurface.Raveling.Progression.A0 = raveling_progress_a0
					maintenanceData.IcResult.RefSurface.Raveling.Progression.A1 = raveling_progress_a1
					maintenanceData.IcResult.RefSurface.Raveling.Progression.A2 = raveling_progress_a2

					maintenanceData.IcResult.RefSurface.AllStructuralCrack.Initial.HSOLD.A0 = asc_init_hsold_equal_0_a0
					maintenanceData.IcResult.RefSurface.AllStructuralCrack.Initial.HSOLD.A1 = asc_init_hsold_equal_0_a1
					maintenanceData.IcResult.RefSurface.AllStructuralCrack.Initial.HSOLD.A2 = asc_init_hsold_equal_0_a2
					maintenanceData.IcResult.RefSurface.AllStructuralCrack.Initial.HSOLD.A3 = asc_init_hsold_equal_0_a3
					maintenanceData.IcResult.RefSurface.AllStructuralCrack.Initial.HSOLD.A4 = asc_init_hsold_equal_0_a4

					maintenanceData.IcResult.RefSurface.AllStructuralCrack.Initial.HSOLD_O.A0 = asc_init_hsold_over_0_a0
					maintenanceData.IcResult.RefSurface.AllStructuralCrack.Initial.HSOLD_O.A1 = asc_init_hsold_over_0_a1
					maintenanceData.IcResult.RefSurface.AllStructuralCrack.Initial.HSOLD_O.A2 = asc_init_hsold_over_0_a2
					maintenanceData.IcResult.RefSurface.AllStructuralCrack.Initial.HSOLD_O.A3 = asc_init_hsold_over_0_a3
					maintenanceData.IcResult.RefSurface.AllStructuralCrack.Initial.HSOLD_O.A4 = asc_init_hsold_over_0_a4

					maintenanceData.IcResult.RefSurface.AllStructuralCrack.Progression.HSOLD.A0 = asc_progress_hsold_equal_0_a0
					maintenanceData.IcResult.RefSurface.AllStructuralCrack.Progression.HSOLD.A1 = asc_progress_hsold_equal_0_a1

					maintenanceData.IcResult.RefSurface.AllStructuralCrack.Progression.HSOLD_O.A0 = asc_progress_hsold_over_0_a0
					maintenanceData.IcResult.RefSurface.AllStructuralCrack.Progression.HSOLD_O.A1 = asc_progress_hsold_over_0_a1

					maintenanceData.IcResult.RefSurface.WideStructuralCrack.Initial.HSOLD.A0 = wsc_init_hsold_equal_0_a0
					maintenanceData.IcResult.RefSurface.WideStructuralCrack.Initial.HSOLD.A1 = wsc_init_hsold_equal_0_a1
					maintenanceData.IcResult.RefSurface.WideStructuralCrack.Initial.HSOLD.A2 = wsc_init_hsold_equal_0_a2

					maintenanceData.IcResult.RefSurface.WideStructuralCrack.Initial.HSOLD_O.A0 = wsc_init_hsold_over_0_a0
					maintenanceData.IcResult.RefSurface.WideStructuralCrack.Initial.HSOLD_O.A1 = wsc_init_hsold_over_0_a1
					maintenanceData.IcResult.RefSurface.WideStructuralCrack.Initial.HSOLD_O.A2 = wsc_init_hsold_over_0_a2

					maintenanceData.IcResult.RefSurface.WideStructuralCrack.Progression.HSOLD.A0 = wsc_progress_hsold_equal_0_a0
					maintenanceData.IcResult.RefSurface.WideStructuralCrack.Progression.HSOLD.A1 = wsc_progress_hsold_equal_0_a1

					maintenanceData.IcResult.RefSurface.WideStructuralCrack.Progression.HSOLD_O.A0 = wsc_progress_hsold_over_0_a0
					maintenanceData.IcResult.RefSurface.WideStructuralCrack.Progression.HSOLD_O.A1 = wsc_progress_hsold_over_0_a1

					maintenanceData.IcResult.RefSurface.RuttingPlasticDeformation.A0 = rpd_a0
					maintenanceData.IcResult.RefSurface.RuttingPlasticDeformation.A1 = rpd_a1
					maintenanceData.IcResult.RefSurface.RuttingPlasticDeformation.A2 = rpd_a2

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
			if maintenanceDatas != nil {
				maintenances[roadID] = maintenanceDatas
			}

			road := responses.Road{RoadID: roadID, RoadGroupID: roadData.RoadGroupId, RefDirectionID: roadInfo.RefDirectionId, RoadName: roadInfo.Name, RoadGroupName: roadGroup.ShortName, YearConstructionCompleted: roadInfo.YearConstructionCompleted}

			roadGeom := responses.RoadGeoms{RoadID: roadID, KmStart: float64(roadGeomData.KmStart), KmEnd: float64(roadGeomData.KmEnd), LaneNo: roadGeomData.LaneNo}
			growthRate, _ := u.GetSettingAadtGrowthRate(road.RoadGroupID)
			aadtParameter, err := u.GetSettingAadtParameter(road.RoadGroupID)
			if err != nil {
				aadtParameter = models.AadtParameterData{}
			}
			loadEquivalent6, err := u.GetLoadEquivalent(aadtParameter.Truck6Axle)
			if err != nil {
				loadEquivalent6 = models.RefAadtParameterVehicleType{}
			}
			aadtParameter.Truck6Axle = loadEquivalent6.NumAxle
			aadtParameter.Truck6LoadEquivalent = loadEquivalent6.LoadEquivalent

			loadEquivalent10, err := u.GetLoadEquivalent(aadtParameter.Truck10Axle)
			if err != nil {
				loadEquivalent10 = models.RefAadtParameterVehicleType{}
			}
			aadtParameter.Truck10Axle = loadEquivalent10.NumAxle
			aadtParameter.Truck10LoadEquivalent = loadEquivalent10.LoadEquivalent

			growthRates[roadID] = growthRate
			aadtParameters[roadID] = aadtParameter

			roads = append(roads, road)
			// for _, item := range roadDateBeginData {
			roadDateBegin := responses.RoadDateBegins{RoadID: roadID, KmStart: float64(roadInfo.KmStart), KmEnd: float64(roadInfo.KmEnd), Year: roadInfo.YearConstructionCompleted}
			roadDateBegins = append(roadDateBegins, roadDateBegin)
			// }

			roadGeoms = append(roadGeoms, roadGeom)
		}

	}
	data, err := u.PrepareDataFn(analysisID, roads, roadDateBegins, roadGeoms, roadSurfaces, maintenances, growthRates, aadtParameters, req.GroupKm, req.SurfaceTypeId)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	return data, nil
}

////////////

func (u *UseCase) PrepareDataFn(analysisID int, roads []responses.Road, roadDateBegins []responses.RoadDateBegins, roadGeom []responses.RoadGeoms, roadSurfaces map[string][]models.RoadSurfacePrepareData, maintenances map[int][]responses.MaintenanceData, growthRates map[int]models.SettingAadtGrowthRate, aadtParameters map[int]models.AadtParameterData, groupKm int, surfaceTypeId int) (interface{}, error) {
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

			for _, maintenance := range maintenances[road.RoadID] {
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

			for _, maintenance := range maintenances[road.RoadID] {
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

		/* แบ่งช่วงทุกๆ 1 กิโลเมตร */
		interval := groupKm * 1000 // 1 kilometer = 1000 meters

		for _, km := range kmRanges {
			kmStart := km["km_start"]
			kmEnd := km["km_end"]

			intvalStart := math.Floor(kmStart / float64(interval)) // หารเอาจำนวนเต็มไม่ปัดเศษ
			intvalEnd := math.Floor(kmEnd / float64(interval))     // หารเอาจำนวนเต็มไม่ปัดเศษ

			if intvalStart == intvalEnd || math.Abs(kmEnd-kmStart) <= 500 { // กรณีระยะทางรวมไม่เกิน 500 ไม่ต้องแบ่งตาม interval
				intervalKmRanges = append(intervalKmRanges, map[string]interface{}{
					"km_start": kmStart,
					"km_end":   kmEnd,
				})
			} else {
				start := kmStart
				intvalStart := intvalStart
				intvalEnd := intvalEnd
				var end float64
				if kmStart < kmEnd {
					for start < kmEnd {
						lt := start + float64(interval)
						end = math.Min(lt, kmEnd)

						if intvalStart != intvalEnd {
							end = end - math.Mod(end, float64(interval))
						}

						intervalKmRanges = append(intervalKmRanges, map[string]interface{}{
							"km_start": start,
							"km_end":   end,
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
							"km_start": start,
							"km_end":   end,
						})

						start = end
						intvalStart = intvalEnd
					}
				}
			}
		}
		// helpers.PrintlnJson("roadroadroad", road.RoadID, roadGeom[road_index].RoadID)
		/* Get ข้อมูลวันที่เริ่มใช้สายทาง, ข้อมูลผิวทาง (road_surface) และประวัติการซ่อมบำรุง (Maintenance) */
		i := 0
		for indexR0 := range intervalKmRanges {
			if indexR0 >= indexStart {
				index := i + indexStart
				// helpers.PrintlnJson("roadroadroad", road.RoadID, roadGeom[road_index].RoadID)
				// helpers.PrintlnJson("road", index, len(intervalKmRanges), len(kmRanges), indexR0, indexStart)
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
				for _, maintenance := range maintenances[road.RoadID] {
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
		// fmt.Println("intervalKmRanges", len(intervalKmRanges), indexStart)
	}

	/* Prepare Data */
	dataRes := make([]map[string]interface{}, 0)
	// return intervalKmRanges, nil
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
		// return roadConditions, nil
		roadConditionDataRes := responses.RoadConditionPrepareDataRes{}
		if len(roadConditions) != 0 {
			roadConditionData, _ := u.GetRoadConditionPrepareData(directionID, laneNo, roadConditions, result)
			copier.Copy(&roadConditionDataRes, roadConditionData)
			roadConditionDataRes.SurveyDate = surveyedDate
		} else { //ปรับแก้ case มากกว่า 3 เลน
			// lane 1
			if laneNo == 1 {
				helpers.PrintlnJson("rrrrrrrrrrrrrrrrrr")
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
				helpers.PrintlnJson("bbbbbbbbbb")
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
				helpers.PrintlnJson("oooooooooo")
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

		if surfaceTypeId == 1 {
			if !helpers.InArrayInt(lastType, []int{1, 2, 3, 4}) {
				continue
			}
		} else {
			if !helpers.InArrayInt(lastType, []int{5, 6}) {
				continue
			}
		}

		/* ความเสียหาย */
		roadDamages, surveyedDate, err := u.GetRoadDamage(roadID, laneNo)
		roadDamagesData, _ := u.GetRoadDamagePrepareData(directionID, laneNo, roadDamages, result)
		if err != nil {
			roadDamagesData.SurveyDate = nil
		} else {
			roadDamagesData.SurveyDate = surveyedDate
		}

		result["road_damage"] = roadDamagesData

		roadGroupID := result["road"].(responses.Road).RoadGroupID
		volumeAadt, _ := u.GetVolumeAadt(roadID)
		result["volume_aadt"] = volumeAadt

		volumeRain, _ := u.GetVolumeRain(roadGroupID)
		result["volume_rain"] = volumeRain

		currentSurface, _ := u.GetCurrentSurface(result, roadConditionDataRes.SurveyDate)

		result["current_surface"] = currentSurface

		// year := time.Now().Year() // ปีที่กดปุ่มวิเคราะห์

		t := 0 // ปีที่วิเคราะห์ ปีที่ 1 = 0
		if roadConditionDataRes.SurveyDate == nil {
			result["analyst_year"] = 0 //roadConditionDataRes.SurveyDate.(time.Time).Year() //year + (t + 1) // ปีที่วิเคราะห์
		} else {
			result["analyst_year"] = roadConditionDataRes.SurveyDate.(time.Time).Year() //year + (t + 1) // ปีที่วิเคราะห์
		}

		if helpers.InArrayInt(result["type"].(int), []int{1, 2, 3, 4}) {
			/* หา Age */
			age, _ := u.GetAgeAC(result, roadConditionDataRes.SurveyDate)
			result["age"] = age

			/* หา HSOLD และ HSNEW */
			hsold_hsnew, _ := u.GetHsoldHsnew(result, roadConditionDataRes.SurveyDate)
			result["hsold_hsnew"] = hsold_hsnew
		} else {
			age, _ := u.GetAgeCC(result, roadConditionDataRes.SurveyDate)
			result["age"] = age
			result["hsold_hsnew"] = responses.HsoldHsnewPrepareDataRes{}
		}

		if helpers.InArrayInt(result["type"].(int), []int{3, 4}) {
			/* หา SNP */
			snp, err := u.GetSnp(result, roadConditionDataRes.SurveyDate)
			if err != nil {
				logs.Error(err)
				result["snp"] = responses.SnpPrepareDataRes{}
			} else {
				result["snp"] = snp.(responses.SnpPrepareDataRes)
			}
		} else {
			result["snp"] = responses.SnpPrepareDataRes{}
		}

		/* หา Truck Factor */
		aadt := volumeAadt.Veh1 + volumeAadt.Veh2 + volumeAadt.Veh3
		result["aadt"] = aadt
		growthRate := growthRates[roadID]
		aadtParameter := aadtParameters[roadID]
		truckFactor6 := 0.0
		truckFactor10 := 0.0
		percentTruck := 0.0
		if aadtParameter.TrackFactorOpenInput {
			truckFactor6 = aadtParameter.Truck6TruckFactorInput
			truckFactor10 = aadtParameter.Truck10TruckFactorInput
		} else {
			if float64(volumeAadt.Veh3)+float64(volumeAadt.Veh3) != 0 {
				percentTruckFactor6 := float64((float64(volumeAadt.Veh2) / float64(float64(volumeAadt.Veh2)+float64(volumeAadt.Veh3))) * 100)
				truckFactor6 = float64((percentTruckFactor6 / 100)) * aadtParameter.Truck6LoadEquivalent

				percentTruckFactor10 := float64((float64(volumeAadt.Veh3) / float64(float64(volumeAadt.Veh2)+float64(volumeAadt.Veh3))) * 100)
				truckFactor10 = float64((percentTruckFactor10 / 100)) * aadtParameter.Truck10LoadEquivalent
			} else {
				truckFactor6 = 0.0
				truckFactor10 = 0.0
			}

		}

		truckFactor := truckFactor6 + truckFactor10
		if aadt != 0 {
			percentTruck = float64(((float64(volumeAadt.Veh2) + float64(volumeAadt.Veh3)) / float64(aadt)) * 100)
		}

		result["percent_truck"] = percentTruck
		result["truck_factor"] = truckFactor
		/* หา Esal ปีที่ n */
		result["esal"] = ((float64(percentTruck) / 100) * float64(aadt) * 365 * truckFactor * aadtParameter.LaneDistributionFactor * math.Pow((1+(growthRate.GrowthRate)), float64(t)) * aadtParameter.DirectionalDistributionFactor) / 1000000

		/* YAX */
		car4 := (float64(volumeAadt.Veh1)) * aadtParameter.Car4Axle
		truck6 := float64(volumeAadt.Veh2) * float64(aadtParameter.Truck6Axle)
		truck10 := float64(volumeAadt.Veh3) * float64(aadtParameter.Truck10Axle)

		yax := ((car4 + truck6 + truck10) / (aadtParameter.Elane * 1000000)) * math.Pow((1+(growthRate.GrowthRate)), float64(t))
		result["yax"] = yax

		////////////////////////////////////////////////
		roadGeomRes := road["road_geom"].(responses.RoadGeoms)
		roadDateBegin := road["road_date_begin"].(responses.RoadDateBegins)
		age := road["age"].(responses.AgePrepareDataRes)
		roadCondition := road["road_condition"].(responses.RoadConditionPrepareDataRes)
		roadDamage := road["road_damage"].(responses.RoadDamagePrepareDataRes)
		currentSurfaceRes := road["current_surface"].(responses.RefSurface)
		hsoldHsnewRes := road["hsold_hsnew"].(responses.HsoldHsnewPrepareDataRes)
		snpRes := road["snp"].(responses.SnpPrepareDataRes)
		theGeom, err := u.GenTheGeom(roadID, laneNo, directionID, road["km_start"].(float64), road["km_end"].(float64))
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}

		var prepareData models.PrepareData
		prepareData.TheGeom = theGeom.(string)
		prepareData.MaintenanceAnalysisID = analysisID
		prepareData.IsSelected = false
		prepareData.GroupName = road["road"].(responses.Road).RoadGroupName
		prepareData.RoadName = road["road"].(responses.Road).RoadName
		prepareData.RoadGroupID = roadGroupID
		prepareData.RoadID = roadID
		prepareData.LaneNo = laneNo
		prepareData.LaneKmStart = roadGeomRes.KmStart
		prepareData.LaneKmEnd = roadGeomRes.KmEnd
		prepareData.LaneLength = road["lane_length"].(float64)
		prepareData.LaneWidth = road["lane_width"].(float64)
		prepareData.KmStart = road["km_start"].(float64)
		prepareData.KmEnd = road["km_end"].(float64)
		prepareData.Length = road["length"].(float64)
		prepareData.Area = road["area"].(float64)
		prepareData.LastType = road["last_type"].(int)
		prepareData.Type = road["type"].(int)
		prepareData.AnalystYear = road["analyst_year"].(int)
		prepareData.YearRoadBegin = roadDateBegin.Year
		prepareData.YearLastOverlay = age.YearLastOverlay
		prepareData.YearLastSeal = age.YearLastSeal
		prepareData.YearLastMolRcl = age.YearLastMolRcl
		prepareData.YearLastReconstruction = age.YearLastReconstruction
		prepareData.Age = age.Age
		prepareData.Rut = roadCondition.RUT
		prepareData.Iri = roadCondition.IRI
		prepareData.Ifi = roadCondition.IFI
		prepareData.NumberOfPothole = roadDamage.NumberOfPothole
		prepareData.AreaAcIcrack = roadDamage.AreaAcIcrack
		prepareData.PercentAcIcrack = roadDamage.PercentAcIcrack
		prepareData.AreaAcUcrack = roadDamage.AreaAcUcrack
		prepareData.PercentAcUcrack = roadDamage.PercentAcUcrack
		prepareData.PercentAcRavelling = roadDamage.PercentAcRavelling
		prepareData.CcTransverseCrack = roadDamage.CcTransverseCrack
		prepareData.CcFaulting = roadDamage.CcFaulting
		prepareData.CcSpalling = roadDamage.CcSpalling
		prepareData.CurrentSurfaceID = currentSurfaceRes.ID
		prepareData.CurrentSurfaceName = currentSurfaceRes.Name
		prepareData.CurrentSurfaceType = currentSurfaceRes.Type
		prepareData.CurrentSurfaceSurfaceGroup = currentSurfaceRes.SurfaceGroup
		prepareData.CurrentSurfaceLayerCoefficient = currentSurfaceRes.LayerCoefficient
		prepareData.CurrentSurfaceDrainage = currentSurfaceRes.Drainage
		prepareData.CurrentSurfaceA = currentSurfaceRes.A
		prepareData.CurrentSurfaceB = currentSurfaceRes.B
		prepareData.CurrentSurfaceCBase = currentSurfaceRes.Cbase
		prepareData.CurrentSurfaceCExp = currentSurfaceRes.Cexp
		prepareData.CurrentSurfaceCRT = currentSurfaceRes.Crt
		prepareData.CurrentSurfaceRRF = currentSurfaceRes.Rrf
		prepareData.Hsold = hsoldHsnewRes.Hsold
		prepareData.Hsnew = hsoldHsnewRes.Hsnew
		prepareData.SNPSurface = snpRes.SnpSurface
		prepareData.SNPBase = snpRes.SnpBase
		prepareData.SNPSubbase = snpRes.SnpSubbase
		prepareData.SNP = snpRes.Snp
		prepareData.AADT = float64(road["aadt"].(int))
		prepareData.TruckFactor = road["truck_factor"].(float64)
		prepareData.ESAL = road["esal"].(float64)
		prepareData.YAX = road["yax"].(float64)
		if result == nil {
			continue
		}
		jsonData, err := json.Marshal(result)
		if err != nil {
			continue
		}
		prepareData.Data = string(jsonData)
		prepareData.CreatedBy = 1
		prepareData.UpdatedBy = 1
		prepareData.CreatedAt = time.Now()
		prepareData.UpdatedAt = time.Now()
		_, err = u.Repo.CreateFullPrepareData(prepareData)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		dataRes = append(dataRes, result)
	}
	return dataRes, nil
}
func (u *UseCase) GenTheGeom(roadID, laneNo, directionID int, kmStart, kmEnd float64) (interface{}, error) {
	roadGeom, err := u.Repo.GetRoadGeomByRoadIDLaneNo(roadID, laneNo)
	if err != nil {
		return "", err
	}
	// return roadGeom, nil
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (u *UseCase) GetRoadSurface() ([]models.RoadSurface, error) {
	data, err := u.Repo.GetRoadSurface()
	if err != nil {
		logs.Error(err)
		return data, responses.NewAppErr(400, err.Error())
	}

	return data, nil
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

func (u *UseCase) GetMaintenanceData(roadId, laneNo int) (interface{}, error) { //([]models.MaintenanceData, error) {
	data, err := u.Repo.GetMaintenanceData(roadId, laneNo)
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
			//
			// return road.RefSurfaceParamsID, nil

			refCriteriaMethod, err := u.Repo.GetRefCriteriaMethodByID(road.MaintenanceMethodID)
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
			interventionCriteriaParams, err := u.Repo.GetInterventionCriteriaParamsById(road.InterventionCriteriaIDParams)
			if err != nil {
				// helpers.PrintlnJson("road.InterventionCriteriaIDParams", road.InterventionCriteriaIDParams)
				logs.Error(err)
				return data, responses.NewAppErr(400, err.Error())
			}

			params := interventionCriteriaParams.Params
			var data map[string]interface{}
			err = json.Unmarshal([]byte(params), &data)
			if err != nil {
				logs.Error(err)
				return data, responses.NewAppErr(400, err.Error())

				// return
			}

			interventionCriteria := data[surfaceType]

			val, ok := interventionCriteria.(map[string]interface{})
			if !ok {
				// return
			}

			var maintenanceItem []MaintenanceItem
			// criteria := val[criteriaMethodName]
			// fmt.Println(criteriaMethodName)
			criteria := val[criteriaMethodName]
			// fmt.Println("val", val)
			jsonBytes, err := json.Marshal(criteria)
			if err != nil {
				logs.Error(err)
				fmt.Printf("Error marshaling JSON: %s\n", err.Error())
				// return
			}

			err = json.Unmarshal(jsonBytes, &maintenanceItem)
			if err != nil {
				logs.Error(err)
				fmt.Println("Error:", string(jsonBytes))
				// return
			}

			// helpers.PrintlnJson(maintenanceItem)
			// helpers.PrintlnJson(len(maintenanceItem[0]), criteriaMethodName, maintenanceItem, val)
			for _, item := range maintenanceItem {
				// helpers.PrintlnJson(item.ID, road.InterventionCriteriaID)
				if item.ID == road.InterventionCriteriaID {
					// helpers.PrintlnJson(item)
					// maintenanceInterventionCriteria = item
					copier.Copy(&maintenanceInterventionCriteria, item)
				}
			}
			// sss := val["ol_overlay"]
			// helpers.PrintlnJson(person)
			// MaintenanceProcedure
			//criteriaMethod
			// for _, item := range  {
			var mr models.MaintenanceRoadPrepareData
			copier.Copy(&mr, road)
			mr.InterventionCriteria = maintenanceInterventionCriteria
			mrs = append(mrs, mr)
			// }
		}

		var data2 models.MaintenanceData
		// item.InterventionCriteria = maintenanceInterventionCriteria
		copier.Copy(&data2, item)
		data2.MaintenanceRoads = mrs
		// data2.InterventionCriteria = maintenanceInterventionCriteria
		data3 = append(data3, data2)
	}

	return data3, nil
}

// func (u *UseCase) GetMaintenanceByID(mID int) (models.MaintenanceData, error) {
// 	data, err := u.Repo.GetMaintenanceByID(mID)
// 	if err != nil {
// 		logs.Error(err)
// 		return data, err
// 	}

// 	return data, nil
// }

func (u *UseCase) GetRoadGroupByID(ID int) (models.RoadGroup, error) {
	data, err := u.Repo.GetRoadGroupByID(ID)
	if err != nil {
		logs.Error(err)
		return data, err
	}

	return data, nil
}

func (u *UseCase) GetRoadByID(roadID int) (models.RoadInfo, error) {
	data, err := u.Repo.GetRoadInfoByID(roadID)
	if err != nil {
		logs.Error(err)
		return data, err
	}

	return data, nil
}

func (u *UseCase) GetRoadDatebegin(roadID int) ([]models.RoadDatebegin, error) {
	data, err := u.Repo.GetRoadDatebegin(roadID)
	if err != nil {
		logs.Error(err)
		return data, err
	}

	return data, nil
}

func (u *UseCase) GetRoadGroupByRoadID(roadID int) (models.Road, error) {
	data, err := u.Repo.GetRoadByRoadID(roadID)
	if err != nil {
		logs.Error(err)
		return data, err
	}

	return data, nil
}

func (u *UseCase) GetRoadGeomLaneByID(roadID, laneNo int) (models.RoadGeom, error) {
	data, err := u.Repo.GetRoadGeomLaneByID(roadID, laneNo)
	if err != nil {
		logs.Error(err)
		return data, err
	}

	return data, nil
}

func (u *UseCase) GetRefStructureSurfaceByID(ID int) (models.RefStructureSurface, error) {
	data, err := u.Repo.GetRefStructureSurfaceByID(ID)
	if err != nil {
		logs.Error(err)
		return data, err
	}

	return data, nil
}

func (u *UseCase) GetRoadSurfaceLaneBySurfaceIDByLane(surfaceID, laneNo int) (models.RoadSurfaceLane, error) {
	data, err := u.Repo.GetRoadSurfaceLaneBySurfaceIDByLane(surfaceID, laneNo)
	if err != nil {
		logs.Error(err)
		return data, err
	}
	return data, nil
}

func (u *UseCase) GetRefSurfaceByID(ID int) (models.RefSurface, error) {
	data, err := u.Repo.GetRefSurfaceByID(ID)
	if err != nil {
		logs.Error(err)
		return data, err
	}
	return data, nil
}

func (u *UseCase) GetRefSurfaceParam(ID int) (models.RefSurface, error) {
	data, err := u.Repo.GetRefSurfaceParam(ID)
	if err != nil {
		logs.Error(err)
		return models.RefSurface{}, err
	}

	var refSurface models.RefSurface
	err = json.Unmarshal([]byte(data.Params), &refSurface)
	if err != nil {
		logs.Error(err)
		fmt.Println("error:", err)
	}
	// helpers.PrintlnJson(refSurface)
	return refSurface, nil
}

func (u *UseCase) GetRefMaterialBaseByID(ID int) (models.RefMaterialBase, error) {
	data, err := u.Repo.GetRefMaterialBaseByID(ID)
	if err != nil {
		logs.Error(err)
		return data, err
	}
	return data, nil
}

func (u *UseCase) GetRefMaterialSubgradeByID(ID int) (models.RefMaterialSubgrade, error) {
	data, err := u.Repo.GetRefMaterialSubgradeByID(ID)
	if err != nil {
		logs.Error(err)
		// return data, err
	}
	return data, nil
}

func (u *UseCase) GetRefMaterialSubbaseByID(ID int) (models.RefMaterialSubbase, error) {
	data, err := u.Repo.GetRefMaterialSubbaseByID(ID)
	if err != nil {
		logs.Error(err)
		// return data, err
	}
	return data, nil
}

func (u *UseCase) GetInterventionCriteria() (map[int]models.InterventionCriteria, error) {
	res := make(map[int]models.InterventionCriteria)
	data, err := u.Repo.GetInterventionCriteria()
	if err != nil {
		logs.Error(err)
		// return data, err
	}
	for _, item := range data {
		res[item.Id] = item
	}
	return res, nil
}

func (u *UseCase) GetCriteriaMethod() (map[int]models.RefCriteriaMethod, error) {
	res := make(map[int]models.RefCriteriaMethod)
	data, err := u.Repo.GetCriteriaMethod()
	if err != nil {
		logs.Error(err)
		// return data, err
	}
	for _, item := range data {
		res[item.ID] = item
	}
	return res, nil
}

func (u *UseCase) GetRoadWorkEffect() (models.SettingRoadWorkEffect, error) {
	data, err := u.Repo.GetRoadWorkEffect()
	if err != nil {
		logs.Error(err)
		// return data, err
	}

	return data, nil
}

func (u *UseCase) GetSettingAadtParams(roadGrpID int) (int, map[string]interface{}, error) {
	// var settingAadtParams models.SettingAadtParams
	var data map[string]interface{}
	aadtParams, err := u.Repo.GetSettingAadtParams()
	if err != nil {
		logs.Error(err)
		return 0, data, err
	}

	err = json.Unmarshal([]byte(aadtParams.Params), &data)
	if err != nil {
		logs.Error(err)
		return 0, data, err
	}

	return aadtParams.ID, data, nil
}

func (u *UseCase) GetRoadSurfaceByRoadID(roadID, laneNo int) ([]models.RoadSurfacePrepareData, error) {
	surfaceGroup, _ := u.Repo.GetRoadSurfaceByGroupRoadID(roadID, laneNo)
	for _, grp := range surfaceGroup {
		data, _ := u.Repo.GetRoadSurfaceByRoadID(roadID, laneNo, grp)
		if len(data) > 0 {
			year := data[0].Year
			_, err := u.Repo.GetMaintenanceHistoryByRoadID(roadID, year)
			if err == nil {
				var roadSurfaces []models.RoadSurfacePrepareData
				for _, item := range data {
					var roadSurfaceLanes []models.RoadSurfaceLanePrePareData
					for _, item2 := range item.RoadSurfaceLane {
						if item2.RefSurfaceParamsID == 0 {
							continue
						}
						surfaceParams, _ := u.Repo.GetRefSurfaceParam(item2.RefSurfaceParamsID)
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
	surface, err := u.Repo.GetRoadSurfaceFirstGrpByRoadID(roadID, laneNo)
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
			surfaceParams, _ := u.Repo.GetRefSurfaceParam(item2.RefSurfaceParamsID)
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

func (u *UseCase) GetRefCriteriaMethodByID(ID int) (models.RefCriteriaMethod, error) {
	data, err := u.Repo.GetRefCriteriaMethodByID(ID)
	if err != nil {
		return models.RefCriteriaMethod{}, err
	}
	return data, nil
}

func (u *UseCase) GetInterventionCriteriaParamsById(ID int) (models.SettingInterventionCriteriaParams, error) {
	data, err := u.Repo.GetInterventionCriteriaParamsById(ID)
	if err != nil {
		return models.SettingInterventionCriteriaParams{}, err
	}
	return data, nil
}

func (u *UseCase) GetRoadCondition(roadID, laneNo int) ([]models.RoadConditionSurveyM, time.Time, error) {
	data, surveyedDate, err := u.Repo.GetRoadCondition(roadID, laneNo)
	if err != nil {
		logs.Error(err)
		return []models.RoadConditionSurveyM{}, surveyedDate, err
	}
	return data, surveyedDate, nil
}

type NewKm struct {
	NewKmStart *float64 `json:"new_km_start"`
	NewKmEnd   *float64 `json:"new_km_end"`
}

func (u *UseCase) GetRoadConditionPrepareData(directionID, laneNo int, roadConditions []models.RoadConditionSurveyM, road map[string]interface{}) (responses.RoadConditionPrepareDataRes, error) {
	// helpers.PrintlnJson("roadConditions", roadConditions)
	km_start := road["km_start"].(float64)
	km_end := road["km_end"].(float64)
	rut := 0.0
	lenght_rut := 0.0

	iri := 0.0
	lenght_iri := 0.0

	ifi := 0.0
	lenght_ifi := 0.0
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
	if newKmStart == nil || newKmEnd == nil {
		// helpers.PrintlnJson("roadConditions", roadConditions)
	} else {
		helpers.PrintlnJson("newKmStart", newKmStart, newKmEnd)
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
						ifi += *conditionM.IFI * lengh
						lenght_ifi += lengh
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
						ifi += *conditionM.IFI * lengh
						lenght_ifi += lengh
					}
					// fmt.Println("*conditionM.IRI", conditionM.KmStart, &conditionM.KmEnd, conditionM.IRI)
				}
			}
		}

	}

	rut = (rut / lenght_rut)
	iri = (iri / lenght_iri)
	ifi = (ifi / lenght_ifi)
	// helpers.PrintlnJson(rut, iri, ifi)
	var roadConditionPrepareDataRes responses.RoadConditionPrepareDataRes
	roadConditionPrepareDataRes.RUT = rut
	roadConditionPrepareDataRes.IRI = iri
	roadConditionPrepareDataRes.IFI = ifi
	return roadConditionPrepareDataRes, nil
}

func (u *UseCase) GetRoadDamage(roadID int, laneNo int) ([]models.RoadDamageMPrepareData, time.Time, error) {
	data, surveyedDate, err := u.Repo.GetRoadDamage(roadID, laneNo)
	if err != nil {
		return []models.RoadDamageMPrepareData{}, surveyedDate, err
	}
	return data, surveyedDate, nil
}
func (u *UseCase) GetRoadDamagePrepareData(directionID int, laneNo int, roadDamages []models.RoadDamageMPrepareData, road map[string]interface{}) (responses.RoadDamagePrepareDataRes, error) {
	km_start := road["km_start"].(float64)
	km_end := road["km_end"].(float64)
	/* Asphalt (AC) */
	number_of_pothole := 0.0

	area_ac_icrack := 0.0
	percent_ac_icrack := 0.0

	area_ac_ucrack := 0.0
	percent_ac_ucrack := 0.0

	percent_ac_ravelling := 0.0

	/* Concrete (CC) */
	cc_transverse_crack := 0.0

	cc_faulting := 0.0
	count_cc_faulting := 0.0

	cc_spalling := 0.0

	if km_start < km_end { // LT
		for _, damage := range roadDamages {
			if (damage.Km >= km_start) && (damage.Km <= km_end) {
				/* Asphalt (AC) */
				if damage.AcPothole != nil && *damage.AcPothole != 0 {
					number_of_pothole++
				}

				if damage.AcIcrack != nil {
					area_ac_icrack += *damage.AcIcrack
				}

				if damage.AcUcrack != nil {
					area_ac_ucrack += *damage.AcUcrack
				}

				if damage.AcRavelling != nil {
					percent_ac_ravelling += *damage.AcRavelling
				}

				/* Concrete (CC) */
				if damage.CcTransverseCrack != nil {
					cc_transverse_crack += *damage.CcTransverseCrack
				}

				if damage.CcFaulting != nil {
					cc_faulting += *damage.CcFaulting
					count_cc_faulting++
				}

				if damage.CcSpalling != nil {
					cc_spalling += *damage.CcSpalling
				}
			}
		}
	} else { // RT
		for _, damage := range roadDamages {
			if (damage.Km <= km_start) && (damage.Km >= km_end) {
				/* Asphalt (AC) */
				if damage.AcPothole != nil && *damage.AcPothole != 0 {
					number_of_pothole++
				}

				if damage.AcIcrack != nil {
					area_ac_icrack += *damage.AcIcrack
				}

				if damage.AcUcrack != nil {
					area_ac_ucrack += *damage.AcUcrack
				}

				if damage.AcRavelling != nil {
					percent_ac_ravelling += *damage.AcRavelling
				}

				/* Concrete (CC) */
				if damage.CcTransverseCrack != nil {
					cc_transverse_crack += *damage.CcTransverseCrack
				}

				if damage.CcFaulting != nil {
					cc_faulting += *damage.CcFaulting
					count_cc_faulting++
				}

				if damage.CcSpalling != nil {
					cc_spalling += *damage.CcSpalling
				}
			}
		}
	}

	/* Asphalt (AC) */
	percent_ac_icrack = (area_ac_icrack * 100) / road["area"].(float64)

	area_ac_ucrack = area_ac_ucrack * 0.003
	percent_ac_ucrack = (area_ac_ucrack * 100) / road["area"].(float64)

	percent_ac_ravelling = (percent_ac_ravelling * 100) / road["area"].(float64)

	/* Concrete (CC) */
	cc_transverse_crack = (cc_transverse_crack * 1610) / road["length"].(float64)

	if cc_faulting > 0 {
		cc_faulting = (cc_faulting / count_cc_faulting) / 25.4
	}

	if cc_spalling > 0 {
		cc_spalling = (cc_spalling * 0.2 * 100) / road["area"].(float64)
	}
	var roadDamagePrepareDataRes responses.RoadDamagePrepareDataRes
	roadDamagePrepareDataRes.NumberOfPothole = number_of_pothole
	roadDamagePrepareDataRes.AreaAcIcrack = area_ac_icrack
	roadDamagePrepareDataRes.PercentAcIcrack = percent_ac_icrack
	roadDamagePrepareDataRes.AreaAcUcrack = area_ac_ucrack
	roadDamagePrepareDataRes.PercentAcUcrack = percent_ac_ucrack
	roadDamagePrepareDataRes.PercentAcRavelling = percent_ac_ravelling
	roadDamagePrepareDataRes.CcTransverseCrack = cc_transverse_crack
	roadDamagePrepareDataRes.CcFaulting = cc_faulting
	roadDamagePrepareDataRes.CcSpalling = cc_spalling
	return roadDamagePrepareDataRes, nil
}

func (u *UseCase) GetVolumeAadt(roadId int) (responses.VolumeAadtPrepareDataRes, error) {
	data, err := u.Repo.GetVolumeAadt(roadId)
	if err != nil {
		logs.Error(err)
		return responses.VolumeAadtPrepareDataRes{}, responses.NewAppErr(400, err.Error())
	}
	var volumeAadt responses.VolumeAadtPrepareDataRes
	volumeAadt.Veh1 = data.Veh1
	volumeAadt.Veh2 = data.Veh2
	volumeAadt.Veh3 = data.Veh3
	// volumeAadt.Veh4 = data.Veh4
	return volumeAadt, nil
}

func (u *UseCase) GetVolumeRain(roadId int) (responses.VolumeRainPrepareDataRes, error) {
	data, err := u.Repo.GetVolumeRain(roadId)
	if err != nil {
		logs.Error(err)
		return responses.VolumeRainPrepareDataRes{}, responses.NewAppErr(400, err.Error())
	}
	var volumeRain responses.VolumeRainPrepareDataRes
	volumeRain.RoadGroupID = data.RoadGroupID
	volumeRain.MinRain = data.MinRain
	volumeRain.MaxRain = data.MaxRain
	volumeRain.AvgRain = data.AvgRain
	return volumeRain, nil
}

func (u *UseCase) GetCurrentSurface(road map[string]interface{}, surveyDate interface{}) (interface{}, error) {

	currentSurface := road["surface"].(models.RoadSurfacePrepareData).RoadSurfaceLane
	maintenances := road["maintenances"].([]responses.MaintenanceData)
	if len(maintenances) > 0 {
		for _, item := range maintenances {
			if surveyDate == nil {
				continue
			}
			if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
				if road["type"].(int) != 5 {
					current_surface := item.IcResult.RefSurface // ['ic_result']['ref_surface']; // ผิวใหม่
					return current_surface, nil
				}

				if road["type"].(int) == 5 && item.IcResult.Type == "Concrete" { // กรณีผิวคอนกรีต จะอัปเดตผิวเฉพาะซ่อมด้วยคอนกรีตเท่านั้น
					current_surface := item.IcResult.RefSurface //['ic_result']['ref_surface']; // ผิวใหม่
					return current_surface, nil
				}
			}
		}
	}
	if len(currentSurface) > 0 {
		return ConverstCurrentSurface(currentSurface[0]), nil
	} else {
		return ConverstCurrentSurface(models.RoadSurfaceLanePrePareData{}), nil
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

func (u *UseCase) GetAgeAC(road map[string]interface{}, surveyDate interface{}) (interface{}, error) {
	analystYear := road["analyst_year"].(int)                                // ปีที่วิเคราะห์ข้อมูล
	yearRoadBegin := road["road"].(responses.Road).YearConstructionCompleted // ปีที่เปิดใช้งาน
	yearLastOverlays := []int{}
	yearLastSeals := []int{}
	yearLastMolRcls := []int{}
	yearLastReconstructions := []int{}

	maintenances := road["maintenances"].([]responses.MaintenanceData)
	if len(maintenances) > 0 {
		for _, item := range maintenances { // เรียงตามปีที่น้อยไปมาก
			if surveyDate == nil {
				continue
			}
			method := item.IcResult.Method
			switch method { // Condition ในการคิด HSOLD
			case "OL-Overlay": // ลาดยาง (AC)
				if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
					yearLastOverlays = append(yearLastOverlays, item.Year)
				}
			case "SS-SlurrySeal": // ลาดยาง (AC)
				if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
					yearLastSeals = append(yearLastSeals, item.Year)
				}
			case "M&OL-Mill&Overlay": // ลาดยาง (AC)
				if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
					yearLastMolRcls = append(yearLastMolRcls, item.Year)
				}
			case "RCL-Recycling": // ลาดยาง (AC)
				if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
					yearLastMolRcls = append(yearLastMolRcls, item.Year)
				}
			case "RB-Rehabilitation": // ลาดยาง (AC) แก้จุดที่1 RC -> RB-Rehabilitation
				if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
					yearLastReconstructions = append(yearLastReconstructions, item.Year)
				}
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

func (u *UseCase) GetAgeCC(road map[string]interface{}, surveyDate interface{}) (interface{}, error) {
	analystYear := road["analyst_year"].(int)                                // ปีที่วิเคราะห์ข้อมูล
	yearRoadBegin := road["road"].(responses.Road).YearConstructionCompleted // ปีที่เปิดใช้งาน
	yearLastFdrs := []int{}
	yearLastOvls := []int{}
	yearLastSeals := []int{}
	yearLastMols := []int{}
	yearLastRBCs := []int{} //เพิ่มวิธีซ่อม RBC
	maintenances := road["maintenances"].([]responses.MaintenanceData)
	if len(maintenances) > 0 {
		for _, item := range maintenances { // เรียงตามปีที่น้อยไปมาก
			if surveyDate == nil {
				continue
			}
			method := item.IcResult.Method
			switch method { // Condition ในการคิด HSOLD
			case "FDR": // คอนกรีต (CC)
				if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
					yearLastFdrs = append(yearLastFdrs, item.Year)
				}
			case "BCO": // คอนกรีต (CC)
				if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
					yearLastOvls = append(yearLastOvls, item.Year)
				}
			case "M-OL": // คอนกรีต (CC)
				if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
					yearLastMols = append(yearLastMols, item.Year)
				}
			case "Seal": // คอนกรีต (CC)
				if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
					yearLastSeals = append(yearLastSeals, item.Year)
				}
			case "RBC": // คอนกรีต (CC)
				if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
					yearLastRBCs = append(yearLastRBCs, item.Year)
				}
			}
		}

	}
	/* หาปีล่าสุดของการซ่อมแต่ละประเภท */
	yearLastFdr := helpers.FindMaxInt(yearLastFdrs) // max($year_last_overlay);
	yearLastOvl := helpers.FindMaxInt(yearLastOvls)
	yearLastMol := helpers.FindMaxInt(yearLastMols)
	yearLastSeal := helpers.FindMaxInt(yearLastSeals)
	yearLastRBC := helpers.FindMaxInt(yearLastRBCs) // add
	age := analystYear - helpers.FindMaxInt([]int{yearRoadBegin, yearLastFdr, yearLastOvl, yearLastMol, yearLastSeal, yearLastRBC})
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

func (u *UseCase) GetHsoldHsnew(road map[string]interface{}, surveyDate interface{}) (interface{}, error) {
	maintenances := road["maintenances"].([]responses.MaintenanceData)
	surface := road["surface"].(models.RoadSurfacePrepareData)

	// return surface, nil
	rwe, err := u.Repo.GetSettingRoadWorkEffect()
	if err != nil {
		logs.Error(err)
		// return responses.AadtPercentageVehicleType{}, err
	}
	params := rwe.Params

	var settingRwe models.GetRoadWorkEffectParams
	err = json.Unmarshal([]byte(params), &settingRwe)
	if err != nil {
		logs.Error(err)
		// return responses.AadtPercentageVehicleType{}, err
	}
	var hsoldHsnew responses.HsoldHsnewPrepareDataRes
	if len(maintenances) == 0 { // ไม่มีประวัติซ่อม
		hsoldHsnew.Hsold = 0
		hsoldHsnew.Hsnew = surface.ThicknessSurface * 10 // road_surface
		return hsoldHsnew, nil
	}
	hsold := 0.0
	hsnew := 0.0                                              // ค่า hsnew เริ่มต้นเป็น 0 เสมอ
	if surface.RefStructureSurface.ID == road["type"].(int) { // กรณีประเภทหน้าตัดผิวทาง = type
		hsold = surface.ThicknessSurface // road_surface
	} else { // กรณี concrete เปลี่ยนเป็น composite จะไม่คิด hsold ตั้งต้น
		hsold = 0.0
	}

	for _, item := range maintenances {
		if surveyDate == nil {
			continue
		}
		if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
			method := item.IcResult.Method
			switch method {
			case "OL-Overlay": // ลาดยาง (AC)
				hsold = hsnew + hsold
				hsnew = item.IcResult.ThicknessRepair
				hsoldHsnew.Hsold = hsold * 10
				hsoldHsnew.Hsnew = hsnew * 10
				return hsoldHsnew, nil
			case "SS-SlurrySeal": // ลาดยาง (AC)
				hsold = hsnew + hsold
				hsnew = item.IcResult.ThicknessRepair
				hsoldHsnew.Hsold = hsold * 10
				hsoldHsnew.Hsnew = hsnew * 10
				return hsoldHsnew, nil
			case "M&OL-Mill&Overlay": // ลาดยาง (AC)

				// helpers.PrintlnJson("hsold, hsnew", hsold, hsnew, item.IcResult.ThicknessScrape)
				hsold := (hsold + hsnew) - item.IcResult.ThicknessScrape //$ic_result['thickness_scrape'];
				if hsold < 0 {                                           // ถ้าลบกันแล้วค่าติดลบ จะเซ็ตค่า $hsold  = 0
					hsold = 0
				}
				hsnew = item.IcResult.ThicknessRepair //$ic_result['thickness_repair'];
				hsoldHsnew.Hsold = hsold * 10
				hsoldHsnew.Hsnew = hsnew * 10
				return hsoldHsnew, nil
			case "RCL-Recycling": // ลาดยาง (AC)
				hsold = settingRwe.Asphalt.AsRclDefaultHsOld //$setting_rwe['ac_settings']['RCL-Recycling']['hsold'];
				hsnew = item.IcResult.ThicknessRepair        //$ic_result['thickness_repair'];
				hsoldHsnew.Hsold = hsold * 10
				hsoldHsnew.Hsnew = hsnew * 10
				return hsoldHsnew, nil
			case "RB-Rehabilitation": // ลาดยาง (AC) แก้จุดที่ 2 RC -> RB-Rehabilitation
				hsold = 0
				hsnew = item.IcResult.ThicknessRepair //$ic_result['thickness_repair'];
				hsoldHsnew.Hsold = hsold * 10
				hsoldHsnew.Hsnew = hsnew * 10
				return hsoldHsnew, nil
			default:
				hsoldHsnew.Hsold = hsold * 10
				hsoldHsnew.Hsnew = hsnew * 10
				return hsoldHsnew, nil
			}
		}
	} // ($maintenances as $maintenance) { // เรียงตามปีที่น้อยไปมาก
	hsoldHsnew.Hsold = hsold * 10
	hsoldHsnew.Hsnew = hsnew * 10
	return hsoldHsnew, nil
}

func (u *UseCase) GetSnp(road map[string]interface{}, surveyDate interface{}) (interface{}, error) {
	refMeterialBases, err := u.Repo.GetRefMaterialBaseConcretePavement()
	if err != nil {
		logs.Error(err)
		return models.RefSurface{}, responses.NewAppErr(400, err.Error())
	}
	rwe, _ := u.Repo.GetSettingRoadWorkEffect()
	params := rwe.Params

	//setting_rwe
	var settingRwe models.GetRoadWorkEffectParams
	err = json.Unmarshal([]byte(params), &settingRwe)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	acSettings := settingRwe.Asphalt

	snpSurfaces := []responses.SnpPrepareDataRes{}

	/* หา SNP ของ ผิวทาง */
	surface := road["surface"].(models.RoadSurfacePrepareData)
	surfaceLane := road["surface"].(models.RoadSurfacePrepareData).RoadSurfaceLane
	refSurface := models.RefSurface{}
	if surfaceLane != nil {
		refSurface = surfaceLane[0].RefSurface
	}
	// helpers.PrintlnJson("surfaceLane", surfaceLane)

	if surface.RefStructureSurface.ID == road["type"].(int) { // กรณีประเภทหน้าตัดผิวทาง = type
		snpSurfaces = append(snpSurfaces, responses.SnpPrepareDataRes{Thickness: surface.ThicknessSurface, LayerCoefficient: refSurface.LayerCoefficient, Drainage: refSurface.Drainage})
	} else { // case concrete เป็น composite
		snpSurfaces = append(snpSurfaces, responses.SnpPrepareDataRes{Thickness: 0, LayerCoefficient: 0, Drainage: 0})
	}
	maintenances := road["maintenances"].([]responses.MaintenanceData)

	for _, maintenance := range maintenances {
		if surveyDate == nil {
			continue
		}
		if maintenance.ProjectEndDate.Before(surveyDate.(time.Time)) || maintenance.ProjectEndDate.Equal(surveyDate.(time.Time)) {
			icResult := maintenance.IcResult
			switch icResult.Method {
			case "M&OL-Mill&Overlay":
				/* อัปเดตข้อมูล ความหนา (thickness) ของผิวก่อนหน้า จากปีมากสุดไปน้อยสุด */
				thicknessScrape := icResult.ThicknessScrape // ความหนาของการขูด

				for i := len(snpSurfaces) - 1; i >= 0; i-- {
					updateThickness := snpSurfaces[i].Thickness - thicknessScrape // ความหนาผิวทาง หลังขูดเหลือ
					if updateThickness < 0 {
						updateThickness = 0
					}

					thicknessScrape := thicknessScrape - snpSurfaces[i].Thickness // ความหนาที่ต้องขูดอีก เหลือ
					if thicknessScrape < 0 {
						thicknessScrape = 0
					}
					snpSurfaces[i].Thickness = updateThickness
				}
				/* อัปเดตข้อมูล ความหนา (thickness) ของผิวก่อนหน้า จากปีมากสุดไปน้อยสุด */
				break
			case "RCL-Recycling": // ลาดยาง (AC)
				snpSurfaces = append(snpSurfaces, responses.SnpPrepareDataRes{SnpSurface: 0, SnpBase: 0, SnpSubbase: 0, Snp: acSettings.AsRclSnc})
				return snpSurfaces, nil

				break
			case "Rc-Reconstruction": // ลาดยาง (AC)
				snpSurfaces = append(snpSurfaces, responses.SnpPrepareDataRes{SnpSurface: 0, SnpBase: 0, SnpSubbase: 0, Snp: acSettings.AsRclSnc})
				return snpSurfaces, nil
				break
			}

			/* push ค่าลง snp_surfaces SS-SlurrySeal ไม่คิด SNP */
			switch icResult.Method { //# Condition
			case "OL-Overlay": // ลาดยาง (AC)
				snpSurfaces = append(snpSurfaces, responses.SnpPrepareDataRes{Method: icResult.Method, Thickness: icResult.ThicknessRepair, LayerCoefficient: refSurface.LayerCoefficient, Drainage: icResult.RefSurface.Drainage})
			case "M&OL-Mill&Overlay": // ลาดยาง (AC)
				snpSurfaces = append(snpSurfaces, responses.SnpPrepareDataRes{Method: icResult.Method, Thickness: icResult.ThicknessRepair, LayerCoefficient: refSurface.LayerCoefficient, Drainage: icResult.RefSurface.Drainage})
			}
		}
	}

	/* คำนวณ SNP ของผิว surface ทั้งหมด */
	maxSnpSurfaces := len(snpSurfaces)
	sumSnpSurface := 0.0

	for index, item := range snpSurfaces {
		if (index + 1) == maxSnpSurfaces { // เช็คว่าเป็นผิวชั้นบนสุดไหม ถ้าใช่จะไม่คิด drainage
			sumSnpSurface += item.Thickness * item.LayerCoefficient
		} else {
			sumSnpSurface += (item.Thickness * item.LayerCoefficient * item.Drainage)
		}
	}

	/* หา SNP ของ Base */
	snpBases := []responses.SnpPrepareDataRes{}
	concretePavement := refMeterialBases
	if surface.RefStructureSurface.ID == road["type"].(int) { // กรณีประเภทหน้าตัดผิวทาง = type
		if surface.RefStructureSurface.ID == 4 {
			snpBases = append(snpBases, responses.SnpPrepareDataRes{Thickness: *surface.ThicknessConcreteSlab, LayerCoefficient: concretePavement.LayerCoefficient, Drainage: concretePavement.Drainage})
		}
	} else {
		snpBases = append(snpBases, responses.SnpPrepareDataRes{Thickness: surface.ThicknessSurface, LayerCoefficient: concretePavement.LayerCoefficient, Drainage: concretePavement.Drainage})
		for _, maintenance := range maintenances {
			icResult := maintenance.IcResult
			if icResult.Type == "Concrete" {
				snpBases = append(snpBases, responses.SnpPrepareDataRes{Thickness: icResult.ThicknessRepair, LayerCoefficient: concretePavement.LayerCoefficient, Drainage: concretePavement.Drainage})
			}
		}
	}
	thicknessBase := 0.0
	if surface.ThicknessBase == nil {
		thicknessBase = 0.0
	} else {
		thicknessBase = *surface.ThicknessBase
	}
	snpBases = append(snpBases, responses.SnpPrepareDataRes{Thickness: thicknessBase, LayerCoefficient: surface.RefMaterialBase.LayerCoefficient, Drainage: surface.RefMaterialBase.Drainage})
	sumSnpBase := 0.0
	for _, item := range snpBases {
		sumSnpBase += (item.Thickness * item.LayerCoefficient * item.Drainage)
	}

	/* หา SNP ของ subbase */
	refSubbase := surface.RefMaterialSubbase
	sumSnpSubbase := 0.0
	if surface.ThicknessSubbase != nil {
		sumSnpSubbase = *surface.ThicknessSubbase * refSubbase.LayerCoefficient * refSubbase.Drainage
	} else {
		sumSnpSubbase = refSubbase.LayerCoefficient * refSubbase.Drainage
	}

	/* คำนวณ snp ทั้งหมดของสายทาง */
	snp := sumSnpSurface + sumSnpBase + sumSnpSubbase
	snpRes := responses.SnpPrepareDataRes{}
	snpRes.SnpSurface = sumSnpSurface
	snpRes.SnpBase = sumSnpBase
	snpRes.SnpSubbase = sumSnpSubbase
	snpRes.Snp = snp
	return snpRes, nil

}

func (u *UseCase) GetRrefSurfacrParamsById(ID int) (models.RefSurface, error) {
	surfaceParams, err := u.Repo.GetRrefSurfacrParamsById(ID)
	if err != nil {
		logs.Error(err)
		return models.RefSurface{}, responses.NewAppErr(400, err.Error())
	}

	params := surfaceParams.Params
	var data models.RefSurface
	err = json.Unmarshal([]byte(params), &data)
	if err != nil {
		logs.Error(err)
		fmt.Println("Error decoding JSON:", err)
		// return
	}
	return data, nil
}

func (u *UseCase) GetSettingAadtGrowthRate(RoadGrpID int) (models.SettingAadtGrowthRate, error) {
	data, err := u.Repo.GetSettingAadtGrowthRate(RoadGrpID)
	if err != nil {
		logs.Error(err)
		return models.SettingAadtGrowthRate{}, responses.NewAppErr(400, err.Error())
	}
	return data, nil
}

func (u *UseCase) GetSettingAadtParameter(RoadGrpID int) (models.AadtParameterData, error) {
	data, err := u.Repo.GetSettingAadtParameter(RoadGrpID)
	if err != nil {
		logs.Error(err)
		return models.AadtParameterData{}, responses.NewAppErr(400, err.Error())
	}
	return data, nil
}

func (u *UseCase) GetLoadEquivalent(ID int) (models.RefAadtParameterVehicleType, error) {
	data, err := u.Repo.GetLoadEquivalent(ID)
	if err != nil {
		logs.Error(err)
		return models.RefAadtParameterVehicleType{}, responses.NewAppErr(400, err.Error())
	}
	return data, nil
}

func (u *UseCase) GetType(road map[string]interface{}, surveyDate interface{}) (int, int) {
	refStructureSurface := road["surface"].(models.RoadSurfacePrepareData).RefStructureSurface
	structureSurfaceLastType := 0
	structureSurfaceType := 0
	asphaltCountLastType := 0 // ให้ทำการวนลูปเช็คประวัติการซ่อมบำรุงว่าเคยถูกซ่อมด้วย OL-Overlay ไหม
	asphaltCountType := 0
	roadSurface := road["surface"].(models.RoadSurfacePrepareData)
	roadSurfacLane := roadSurface.RoadSurfaceLane
	refStructureSurfaceID := refStructureSurface.ID
	fmt.Println("ref_structural_id_before", refStructureSurfaceID)
	if len(roadSurfacLane) > 0 {
		roadSurfaceId := roadSurfacLane[0].RefSurface
		fmt.Println("roadSurfaceId", roadSurfaceId)
		switch refStructureSurfaceID {
		case 1:
			fmt.Println("check")
			if roadSurfaceId.SurfaceGroup == "Asphalt" {
				refStructureSurfaceID = refStructureSurface.ID //Asphalt: On Concrete Deck
			} else {
				refStructureSurfaceID = 6 //concrete: on Deck
			}
		case 2:
			fmt.Println("check")
			if roadSurfaceId.SurfaceGroup == "Asphalt" {
				refStructureSurfaceID = refStructureSurface.ID //Asphalt: On Steel Deck
			} else {
				refStructureSurfaceID = 6 //concrete: on Deck
			}
		case 3:
			if roadSurfaceId.SurfaceGroup == "Asphalt" {
				refStructureSurfaceID = refStructureSurface.ID //Asphalt: On Ground
			} else {
				fmt.Println("ref_structural_id_before", refStructureSurfaceID)
				refStructureSurfaceID = 5 //Concrete: on Ground
				fmt.Println("ref_structural_id_check", refStructureSurfaceID)
			}
			fmt.Println("ref_structural_id_after_case3", refStructureSurfaceID)
		case 4:
			fmt.Println("check")
			if roadSurfaceId.SurfaceGroup == "Asphalt" {
				refStructureSurfaceID = refStructureSurface.ID //Composite pavement (Asphalt on concrete base)
			} else {
				refStructureSurfaceID = 5 //Concrete: on Ground
			}
		case 5:
			fmt.Println("check")
			if roadSurfaceId.SurfaceGroup == "Asphalt" {
				refStructureSurfaceID = 3 //Asphalt: On Ground
			} else {
				refStructureSurfaceID = refStructureSurface.ID //Concrete: on Ground
			}
		case 6:
			fmt.Println("check")
			if roadSurfaceId.SurfaceGroup == "Asphalt" {
				refStructureSurfaceID = 1 //Asphalt: On Concrete Deck
			} else {
				refStructureSurfaceID = refStructureSurface.ID //concrete: on Deck
			}
		default:
			refStructureSurfaceID = refStructureSurface.ID
			fmt.Println("ref_structural_id_check1", refStructureSurfaceID)
		}
	} else {
		refStructureSurfaceID = refStructureSurface.ID
		fmt.Println("ref_structural_id_check2", refStructureSurfaceID)
	}

	if refStructureSurfaceID != 5 { // กรณีหน้าตัดผิวทางไม่ใช่ Concrete: on Ground
		structureSurfaceLastType = refStructureSurfaceID
		structureSurfaceType = refStructureSurfaceID
		return structureSurfaceLastType, structureSurfaceType
	}
	concreteCountLastType := 0
	concreteCountType := 0
	fmt.Println("check_case = ", refStructureSurfaceID)
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
						if surveyDate == nil {
							continue
						}
						asphaltCountLastType++
						if item.ProjectEndDate.Before(surveyDate.(time.Time)) || item.ProjectEndDate.Equal(surveyDate.(time.Time)) {
							asphaltCountType++
						}
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
	fmt.Println("ref_structural_last", structureSurfaceLastType, structureSurfaceType)
	return structureSurfaceLastType, structureSurfaceType
}

func (u *UseCase) AnalysisModel(ID int) {
	python := os.Getenv("PYTHON")
	url := fmt.Sprintf(python+"calculate?id=%d", ID)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		logs.Error(err)
		return
	}

	req.Header.Add("app-key", "5ab971d1b81fdcc6bf3b7986a8718a23")
	req.Header.Add("app-secret", "336146ad5ec52021b493c43639a444ac")
	go func() {
		res, err := client.Do(req)
		if err != nil {
			logs.Error(err)
			return
		}
		if res != nil {
			logs.Error("Process the API response here")
		}
	}()
}
