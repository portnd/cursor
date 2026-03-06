package usecases

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/dashboard/domains"
)

type useCaseMaintenance struct {
	repo domains.RepositoryMaintenance
}

func NewUsecaseMaintenance(repo domains.RepositoryMaintenance) domains.UseCaseMaintenance {
	return &useCaseMaintenance{repo: repo}
}

func (u *useCaseMaintenance) GetMaintenanceDashboard(roadIDs []int, depotCodes []string, filter requests.MaintenanceDashboard) (interface{}, error) {

	maintenances, err := u.repo.GetMaintenance(0, roadIDs, depotCodes, filter)
	if err != nil {
		return nil, err
	}

	maintenancesWithLimit, err := u.repo.GetMaintenance(10, roadIDs, depotCodes, filter)
	if err != nil {
		return nil, err
	}

	roadGroups, err := u.repo.GetRoadGroup()
	if err != nil {
		return nil, err
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

	var maintenanceDashboard responses.MaintenanceDashboard
	var numberMaintenanceChart responses.NumberMaintenanceChart
	var maintenanceBudgetChart responses.MaintenanceBudgetChart
	var topTenMaintenanceBudgetChart responses.TopTenMaintenanceBudgetChart

	roadGroupLable := []string{}
	roadGroupMap := map[int]models.RoadGroup{}
	for _, item := range roadGroups {
		if item.Number == "" {
			continue
		}
		number, err := strconv.Atoi(item.Number)
		if err != nil {
			continue
		}

		shortName := "ทางหลวงพิเศษหมายเลข " + strconv.Itoa(number)

		item.ShortName = shortName

		roadGroupLable = append(roadGroupLable, shortName)

		roadGroupMap[item.Id] = item
	}

	roadGroupLable = append(roadGroupLable, "โครงการที่ซ่อมบำรุงมากกว่า 1 สายทาง")

	countRoadGroupBudgetLable := map[string]float64{}
	countRoadGroupLable := map[string]int{}
	mapRoadGroupID := map[string]string{}
	for _, itemMaintenance := range maintenances {
		maintenanceRoadDuplicate := map[int]bool{}
		var countRoadGroupId []int
		roadName := []string{}
		roadSectionName := []string{}
		depotName := []string{}
		for _, itemRoad := range itemMaintenance.MaintenanceRoads {
			if !maintenanceRoadDuplicate[itemRoad.RoadGroupId] {
				roadName = append(roadName, roadGroupMap[itemRoad.RoadGroupId].ShortName)
				roadSectionName = append(roadSectionName, itemRoad.Road.RoadSection.NameOriginTH+" - "+itemRoad.Road.RoadSection.NameDestinationTH)
				depotName = append(depotName, itemRoad.Road.RoadSection.RefDepot.Name)
				countRoadGroupId = append(countRoadGroupId, itemRoad.RoadGroupId)
				maintenanceRoadDuplicate[itemRoad.RoadGroupId] = true
			}
		}

		if len(countRoadGroupId) > 1 {
			countRoadGroupLable["โครงการที่ซ่อมบำรุงมากกว่า 1 สายทาง"] += 1
			countRoadGroupBudgetLable["โครงการที่ซ่อมบำรุงมากกว่า 1 สายทาง"] += itemMaintenance.BudgetMaintenance
			roadGroupIDStr := []string{}
			for _, id := range countRoadGroupId {
				roadGroupIDStr = append(roadGroupIDStr, fmt.Sprintf(`%v`, id))
			}
			mapRoadGroupID["โครงการที่ซ่อมบำรุงมากกว่า 1 สายทาง"] = strings.Join(roadGroupIDStr, ",")
		} else if len(countRoadGroupId) == 1 {
			countRoadGroupLable[roadGroupMap[countRoadGroupId[0]].ShortName] += 1
			countRoadGroupBudgetLable[roadGroupMap[countRoadGroupId[0]].ShortName] += itemMaintenance.BudgetMaintenance
			roadGroupIDStr := []string{}
			for _, id := range countRoadGroupId {
				roadGroupIDStr = append(roadGroupIDStr, fmt.Sprintf(`%v`, id))
			}
			mapRoadGroupID[roadGroupMap[countRoadGroupId[0]].ShortName] = strings.Join(roadGroupIDStr, ",")
		}

	}

	valueNumber := []int{}
	valueBudget := []float64{}
	color := []string{}
	roadGroupID := []string{}
	for index, item := range roadGroupLable {
		valueBudget = append(valueBudget, countRoadGroupBudgetLable[item])
		valueNumber = append(valueNumber, countRoadGroupLable[item])
		color = append(color, Colors[index])
		roadGroupID = append(roadGroupID, mapRoadGroupID[item])
	}

	numberMaintenanceChart.Name = "จำนวนโครงการซ่อมบำรุง"
	numberMaintenanceChart.Lable = roadGroupLable
	numberMaintenanceChart.Data = valueNumber
	numberMaintenanceChart.Color = color
	numberMaintenanceChart.RoadGroupID = roadGroupID

	maintenanceBudgetChart.Name = "งบประมาณการซ่อมบำรุง"
	maintenanceBudgetChart.Lable = roadGroupLable
	maintenanceBudgetChart.Data = valueBudget
	maintenanceBudgetChart.Color = color
	maintenanceBudgetChart.RoadGroupID = roadGroupID

	maintenanceID := []int{}
	maintenanceBudget := []float64{}
	maintenanceName := []string{}
	for _, item := range maintenancesWithLimit {
		maintenanceID = append(maintenanceID, *item.IDParent)
		maintenanceName = append(maintenanceName, item.Name)
		maintenanceBudget = append(maintenanceBudget, item.BudgetMaintenance)
	}

	color = []string{}
	for index, _ := range maintenanceName {
		color = append(color, Colors[index])
	}

	topTenMaintenanceBudgetChart.Name = "จำนวนโครงการซ่อมบำรุง"
	topTenMaintenanceBudgetChart.Lable = maintenanceName
	topTenMaintenanceBudgetChart.Data = maintenanceBudget
	topTenMaintenanceBudgetChart.Color = color
	topTenMaintenanceBudgetChart.MaintenanceID = maintenanceID

	maintenanceDashboard.UpdatedAt = time.Now()
	maintenanceDashboard.NumberMaintenanceChart = numberMaintenanceChart
	maintenanceDashboard.MaintenanceBudgetChart = maintenanceBudgetChart
	maintenanceDashboard.TopTenMaintenanceBudgetChart = topTenMaintenanceBudgetChart

	return maintenanceDashboard, nil
}

func (u *useCaseMaintenance) GetMaintenanceTableDashboard(roadIDs []int, depotCodes []string, filter requests.MaintenanceDashboard) (interface{}, error) {

	maintenances, err := u.repo.GetMaintenance(0, roadIDs, depotCodes, filter)
	if err != nil {
		return nil, err
	}

	roadGroups, err := u.repo.GetRoadGroup()
	if err != nil {
		return nil, err
	}

	var maintenanceTables []responses.MaintenanceTable

	roadGroupLable := []string{}
	roadGroupMap := map[int]models.RoadGroup{}
	for _, item := range roadGroups {
		if item.Number == "" {
			continue
		}
		number, err := strconv.Atoi(item.Number)
		if err != nil {
			continue
		}

		shortName := "ทางหลวงพิเศษหมายเลข " + strconv.Itoa(number)

		item.ShortName = shortName

		roadGroupLable = append(roadGroupLable, shortName)

		roadGroupMap[item.Id] = item
	}

	roadGroupLable = append(roadGroupLable, "โครงการที่ซ่อมบำรุงมากกว่า 1 สายทาง")

	for _, itemMaintenance := range maintenances {
		maintenanceRoadDuplicate := map[int]bool{}
		var countRoadGroupId []int
		roadName := []string{}
		roadSectionName := []string{}
		depotName := []string{}
		for _, itemRoad := range itemMaintenance.MaintenanceRoads {
			if !maintenanceRoadDuplicate[itemRoad.RoadGroupId] {
				roadName = append(roadName, roadGroupMap[itemRoad.RoadGroupId].ShortName)
				roadSectionName = append(roadSectionName, itemRoad.Road.RoadSection.NameOriginTH+" - "+itemRoad.Road.RoadSection.NameDestinationTH)
				depotName = append(depotName, itemRoad.Road.RoadSection.RefDepot.Name)
				countRoadGroupId = append(countRoadGroupId, itemRoad.RoadGroupId)
				maintenanceRoadDuplicate[itemRoad.RoadGroupId] = true
			}
		}

		if len(itemMaintenance.MaintenanceRoads) > 0 {
			var maintenanceTable responses.MaintenanceTable
			maintenanceTable.ID = itemMaintenance.ID
			maintenanceTable.ContractNumber = itemMaintenance.ContractNumber
			maintenanceTable.RoadName = roadName
			maintenanceTable.SectionName = roadSectionName
			maintenanceTable.RefDepotName = depotName
			maintenanceTable.Budget = itemMaintenance.BudgetMaintenance
			start := helpers.TimeThai2DigitYear(itemMaintenance.ProjectEndDate)
			end := helpers.TimeThai2DigitYear(itemMaintenance.GuaranteeExpirationDate)
			maintenanceTable.GuaranteeExpirationDate = start + " - " + end

			dayLeft := int(itemMaintenance.GuaranteeExpirationDate.Sub(time.Now()).Hours() / 24)
			if dayLeft <= 0 {
				dayLeft = 0
			}

			maintenanceTable.RemainDate = dayLeft

			maintenanceTables = append(maintenanceTables, maintenanceTable)
		}

	}

	return maintenanceTables, nil
}

func (u *useCaseMaintenance) GetMaintenanceMapDashboard(roadIDs []int, depotCodes []string, filter requests.MaintenanceDashboard) (interface{}, error) {
	var maintenanceMapDashboards []responses.MaintenanceMapDashboard

	maintenances, err := u.repo.GetMaintenance(0, roadIDs, depotCodes, filter)
	if err != nil {
		return nil, err
	}

	roadGroups, err := u.repo.GetRoadGroup()
	if err != nil {
		return nil, err
	}

	roadGroupLable := []string{}
	roadGroupMap := map[int]models.RoadGroup{}
	for _, item := range roadGroups {
		if item.Number == "" {
			continue
		}
		number, err := strconv.Atoi(item.Number)
		if err != nil {
			continue
		}

		shortName := "ทางหลวงพิเศษหมายเลข " + strconv.Itoa(number)

		item.ShortName = shortName

		roadGroupLable = append(roadGroupLable, shortName)

		roadGroupMap[item.Id] = item
	}

	for _, itemMaintenance := range maintenances {
		var maintenanceMapDashboard responses.MaintenanceMapDashboard
		maintenanceRoadDuplicate := map[int]bool{}
		for _, itemRoad := range itemMaintenance.MaintenanceRoads {
			if !maintenanceRoadDuplicate[itemRoad.RoadGroupId] {

				dayLeftLimit := itemMaintenance.GuaranteeExpirationDate.Sub(itemMaintenance.ProjectEndDate).Hours() / 24

				dayLeft := itemMaintenance.GuaranteeExpirationDate.Sub(time.Now()).Hours() / 24
				if dayLeft <= 0 {
					dayLeft = 0
				}

				dayPercntage := ((dayLeftLimit - dayLeft) / dayLeftLimit) * 100
				if math.IsNaN(dayPercntage) {
					dayPercntage = 0
				}

				title := ""
				color := ""
				if dayPercntage > 20 {
					title = "ระยะเวลาติดค้ำประกัน > 20%"
					color = "#1F70F3"
				} else {
					title = "ระยะเวลาติดค้ำประกัน <= 20%"
					color = "#F1416C"
				}

				var theGeomJson responses.TheGeomJson
				err := json.Unmarshal([]byte(itemRoad.TheGeomJson), &theGeomJson)
				if err != nil {
					var theGeomJsonPoint responses.TheGeomJsonPoint
					err := json.Unmarshal([]byte(itemRoad.TheGeomJson), &theGeomJsonPoint)
					if err != nil {
						return nil, err
					}

					var theGeomJson2 responses.TheGeomJson
					theGeomJson2.Coordinates = append(theGeomJson2.Coordinates, theGeomJsonPoint.Coordinates)
					theGeomJson2.Type = "LineString"

					theGeomJson = theGeomJson2
				}

				maintenanceMapDashboard.IDParent = *itemMaintenance.IDParent
				maintenanceMapDashboard.Title = title
				maintenanceMapDashboard.LaneNo = itemRoad.LaneNo
				maintenanceMapDashboard.RoadName = roadGroupMap[itemRoad.RoadGroupId].ShortName
				maintenanceMapDashboard.SectionName = itemRoad.Road.RoadSection.NameOriginTH + " - " + itemRoad.Road.RoadSection.NameDestinationTH
				maintenanceMapDashboard.ContractNumber = itemMaintenance.ContractNumber
				maintenanceMapDashboard.Name = itemMaintenance.Name
				maintenanceMapDashboard.Color = color
				maintenanceMapDashboard.KmStart = strings.ReplaceAll(u.FloatToString(itemRoad.KmStart/1000), ".", "+")
				maintenanceMapDashboard.KmEnd = strings.ReplaceAll(u.FloatToString(itemRoad.KmEnd/1000), ".", "+")
				maintenanceMapDashboard.KmTotal = math.Abs(itemRoad.KmStart - itemRoad.KmEnd)
				maintenanceMapDashboard.TheGeom = theGeomJson
				maintenanceMapDashboard.RefDepotName = itemRoad.Road.RoadSection.RefDepot.Name

				maintenanceMapDashboards = append(maintenanceMapDashboards, maintenanceMapDashboard)
				maintenanceRoadDuplicate[itemRoad.RoadGroupId] = true
			}
		}

	}

	return maintenanceMapDashboards, nil
}

func (r *useCaseMaintenance) FloatToString(f float64) string {
	return fmt.Sprintf("%.3f", f)
}
