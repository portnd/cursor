package usecases

import (
	"context"
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

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/xuri/excelize/v2"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	servicesDB "gitlab.com/mims-api-service/services/database"
	"gitlab.com/mims-api-service/src/report/domains"
)

type UseCase struct {
	Repo       domains.Repository
	servicesDB servicesDB.ServicesDatabaseDomain
}

// init usecase
func NewUseCase(repo domains.Repository, servicesDB servicesDB.ServicesDatabaseDomain) domains.UseCase {
	return &UseCase{
		Repo:       repo,
		servicesDB: servicesDB,
	}
}

// ////////////// NEW MIMS ////////////////

func (u *UseCase) CheckReportStatus(id int) (interface{}, error) {
	reportStatus, err := u.Repo.CheckReportStatusById(id)
	if err != nil {
		return nil, err
	}

	var response responses.ReportStatus
	response.IsFinish = reportStatus.IsFinish
	response.Path = reportStatus.Path

	return response, nil
}

func (u *UseCase) ReportStatus() (interface{}, error) {

	reportStatusId, err := u.Repo.CreateReportStatus()
	if err != nil {
		return nil, err
	}

	// ! TODO go func() Start
	path := "test"

	err = u.Repo.UpdateReportStatus(reportStatusId, path)
	if err != nil {
		return nil, err
	}

	// ! TODO go func() End

	var response responses.ReportStatusId
	response.Id = reportStatusId

	return response, nil
}

func (u *UseCase) Report1() (interface{}, error) {
	return nil, nil
}

func (u *UseCase) FilterAssetType1(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"download_all_asset_report"}, []string{"download_owner_asset_report"})
	//============ end check permission ============
	filterAssetRoad, err := u.Repo.FilterAssetRoad(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}

	filterAsset, err := u.Repo.FilterAsset()
	if err != nil {
		return nil, err
	}

	roadGroupDuplicate := map[string]bool{}
	dapotDuplicate := map[int]bool{}
	roadSectionDuplicate := map[string]bool{}
	depots := []responses.FilterDepot{}
	depotToRoadGroups := map[int][]responses.FilterRoadGroup{}
	roadGroupToRoadSections := map[string][]responses.FilterRoadSection{}

	for _, item := range filterAssetRoad {
		if item.RefDepot.Name == "" {
			continue
		}
		if dapotDuplicate[item.RefDepot.Id] == false {
			depot := responses.FilterDepot{}
			if item.RefDepot.Name == "" {
				continue
			}
			depot.Id = item.RefDepot.Id
			depot.Name = item.RefDepot.Name
			depots = append(depots, depot)
			dapotDuplicate[item.RefDepot.Id] = true
		}

		if roadGroupDuplicate[strconv.Itoa(item.RefDepot.Id)+"_"+strconv.Itoa(item.RoadGroupId)] == false {

			number, err := strconv.Atoi(item.RoadGroup.Number)
			if err != nil {
				number = 0
			}

			depotToRoadGroup := responses.FilterRoadGroup{}
			depotToRoadGroup.Id = item.RoadGroup.Id
			depotToRoadGroup.Name = "ทางหลวงหมายเลข " + strconv.Itoa(number) + " " + item.RoadGroup.ShortName
			depotToRoadGroups[item.RefDepot.Id] = append(depotToRoadGroups[item.RefDepot.Id], depotToRoadGroup)
			roadGroupDuplicate[strconv.Itoa(item.RefDepot.Id)+"_"+strconv.Itoa(item.RoadGroupId)] = true
		}

		mapString := fmt.Sprintf(`%d_%d`, item.RefDepot.Id, item.RoadGroup.Id)
		roadGroupToRoadSection := responses.FilterRoadSection{}
		roadGroupToRoadSection.Id = item.Id
		number, err := strconv.Atoi(item.Number)
		if err != nil {
			number = 0
		}

		roadGroupToRoadSection.Number = number
		roadGroupToRoadSection.Name = "ตอนควบคุม " + item.Number + " " + item.NameOriginTH + " - " + item.NameDestinationTH
		roadGroupToRoadSections[mapString] = append(roadGroupToRoadSections[mapString], roadGroupToRoadSection)
		roadSectionDuplicate[strconv.Itoa(item.RoadGroupId)+"_"+strconv.Itoa(item.Id)] = true

	}

	var assetRefAssets []responses.AssetRefAsset
	for _, item1 := range filterAsset {
		var assetRefAsset responses.AssetRefAsset
		assetRefAsset.Id = item1.ID
		assetRefAsset.Name = item1.Name
		var assetRefAssetTables []responses.AssetRefAssetTable
		for _, item2 := range item1.RefAssetTable {
			var assetRefAssetTable responses.AssetRefAssetTable
			assetRefAssetTable.Id = item2.ID
			assetRefAssetTable.Name = item2.TableLabel
			assetRefAssetTables = append(assetRefAssetTables, assetRefAssetTable)
		}

		if len(assetRefAssetTables) != 0 {
			assetRefAsset.Asset = assetRefAssetTables
			assetRefAssets = append(assetRefAssets, assetRefAsset)
		}
	}

	for index1, item1 := range depots {

		if len(depotToRoadGroups[item1.Id]) == 0 {
			depots[index1].RoadGroup = []responses.FilterRoadGroup{}
		} else {
			depots[index1].RoadGroup = depotToRoadGroups[item1.Id]
		}

		for index2, item2 := range depots[index1].RoadGroup {
			mapString := fmt.Sprintf(`%d_%d`, item1.Id, item2.Id)
			if len(roadGroupToRoadSections[mapString]) == 0 {
				depots[index1].RoadGroup[index2].RoadSection = []responses.FilterRoadSection{}
			} else {
				sort.SliceStable(roadGroupToRoadSections[mapString], func(i, j int) bool {
					return roadGroupToRoadSections[mapString][i].Number < roadGroupToRoadSections[mapString][j].Number
				})
				depots[index1].RoadGroup[index2].RoadSection = roadGroupToRoadSections[mapString]
			}
		}

	}

	var asset responses.FilterAssetType1
	asset.FilterRoad = depots
	asset.FilterAsset = assetRefAssets

	return asset, nil
}

func (u *UseCase) FilterAssetType2(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"download_all_map_asset_report"}, []string{"download_owner_map_asset_report"})
	//============ end check permission ============
	filterAssetRoad, err := u.Repo.FilterAssetRoad(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}

	filterAsset, err := u.Repo.FilterAsset()
	if err != nil {
		return nil, err
	}

	roadGroupDuplicate := map[string]bool{}
	dapotDuplicate := map[int]bool{}
	roadSectionDuplicate := map[string]bool{}
	depots := []responses.FilterDepot{}
	depotToRoadGroups := map[int][]responses.FilterRoadGroup{}
	roadGroupToRoadSections := map[string][]responses.FilterRoadSection{}

	for _, item := range filterAssetRoad {
		if item.RefDepot.Name == "" {
			continue
		}
		if !dapotDuplicate[item.RefDepot.Id] {
			depot := responses.FilterDepot{}
			depot.Id = item.RefDepot.Id
			depot.Name = item.RefDepot.Name
			depots = append(depots, depot)
			dapotDuplicate[item.RefDepot.Id] = true
		}

		if !roadGroupDuplicate[strconv.Itoa(item.RefDepot.Id)+"_"+strconv.Itoa(item.RoadGroupId)] {
			number, err := strconv.Atoi(item.RoadGroup.Number)
			if err != nil {
				number = 0
			}
			depotToRoadGroup := responses.FilterRoadGroup{}
			depotToRoadGroup.Id = item.RoadGroup.Id
			depotToRoadGroup.Name = "ทางหลวงหมายเลข " + strconv.Itoa(number) + " " + item.RoadGroup.ShortName
			depotToRoadGroups[item.RefDepot.Id] = append(depotToRoadGroups[item.RefDepot.Id], depotToRoadGroup)
			roadGroupDuplicate[strconv.Itoa(item.RefDepot.Id)+"_"+strconv.Itoa(item.RoadGroupId)] = true
		}

		mapString := fmt.Sprintf(`%d_%d`, item.RefDepot.Id, item.RoadGroup.Id)
		roadGroupToRoadSection := responses.FilterRoadSection{}
		roadGroupToRoadSection.Id = item.Id
		number, err := strconv.Atoi(item.Number)
		if err != nil {
			number = 0
		}

		roadGroupToRoadSection.Number = number
		roadGroupToRoadSection.Name = "ตอนควบคุม " + item.Number + " " + item.NameOriginTH + " - " + item.NameDestinationTH
		roadGroupToRoadSections[mapString] = append(roadGroupToRoadSections[mapString], roadGroupToRoadSection)
		roadSectionDuplicate[strconv.Itoa(item.RoadGroupId)+"_"+strconv.Itoa(item.Id)] = true

	}

	var assetRefAssets []responses.AssetRefAsset
	for _, item1 := range filterAsset {
		var assetRefAsset responses.AssetRefAsset
		assetRefAsset.Id = item1.ID
		assetRefAsset.Name = item1.Name
		var assetRefAssetTables []responses.AssetRefAssetTable
		for _, item2 := range item1.RefAssetTable {
			var assetRefAssetTable responses.AssetRefAssetTable
			assetRefAssetTable.Id = item2.ID
			assetRefAssetTable.Name = item2.TableLabel
			assetRefAssetTables = append(assetRefAssetTables, assetRefAssetTable)
		}

		if len(assetRefAssetTables) != 0 {
			assetRefAsset.Asset = assetRefAssetTables
			assetRefAssets = append(assetRefAssets, assetRefAsset)
		}
	}

	for index1, item1 := range depots {
		if len(depotToRoadGroups[item1.Id]) == 0 {
			depots[index1].RoadGroup = []responses.FilterRoadGroup{}
		} else {
			depots[index1].RoadGroup = depotToRoadGroups[item1.Id]
		}

		for index2, item2 := range depots[index1].RoadGroup {
			mapString := fmt.Sprintf(`%d_%d`, item1.Id, item2.Id)
			if len(roadGroupToRoadSections[mapString]) == 0 {
				depots[index1].RoadGroup[index2].RoadSection = []responses.FilterRoadSection{}
			} else {
				sort.SliceStable(roadGroupToRoadSections[mapString], func(i, j int) bool {
					return roadGroupToRoadSections[mapString][i].Number < roadGroupToRoadSections[mapString][j].Number
				})
				depots[index1].RoadGroup[index2].RoadSection = roadGroupToRoadSections[mapString]
			}
		}

	}

	var asset responses.FilterAssetType2
	asset.FilterRoad = depots
	asset.FilterAsset = assetRefAssets

	return asset, nil
}

func (u *UseCase) FilterAssetType3(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"download_all_result_asset_report"}, []string{"download_owner_result_asset_report"})
	//============ end check permission ============
	filterAssetRoad, err := u.Repo.FilterAssetRoad(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}

	filterAsset, err := u.Repo.FilterAsset()
	if err != nil {
		return nil, err
	}

	roadGroupDuplicate := map[string]bool{}
	dapotDuplicate := map[string]bool{}
	roadSectionDuplicate := map[string]bool{}
	depots := []responses.FilterDepot{}
	depotToRoadGroups := map[int][]responses.FilterRoadGroup{}
	roadGroupToRoadSections := map[string][]responses.FilterRoadSection{}

	for _, item := range filterAssetRoad {
		if item.RefDepot.Name == "" {
			continue
		}
		if !dapotDuplicate[item.RefDepot.DepotCode] {
			depot := responses.FilterDepot{}
			depot.Id = item.RefDepot.Id
			depot.Name = item.RefDepot.Name
			depots = append(depots, depot)
			dapotDuplicate[item.RefDepot.DepotCode] = true
		}
		if !roadGroupDuplicate[item.RefDepot.DepotCode+"_"+strconv.Itoa(item.RoadGroupId)] {
			number, err := strconv.Atoi(item.RoadGroup.Number)
			if err != nil {
				number = 0
			}

			depotToRoadGroup := responses.FilterRoadGroup{}
			depotToRoadGroup.Id = item.RoadGroup.Id
			depotToRoadGroup.Name = "ทางหลวงหมายเลข " + strconv.Itoa(number) + " " + item.RoadGroup.ShortName
			depotToRoadGroups[item.RefDepot.Id] = append(depotToRoadGroups[item.RefDepot.Id], depotToRoadGroup)
			roadGroupDuplicate[item.RefDepot.DepotCode+"_"+strconv.Itoa(item.RoadGroupId)] = true
		}

		mapString := fmt.Sprintf(`%d_%d`, item.RefDepot.Id, item.RoadGroup.Id)
		roadGroupToRoadSection := responses.FilterRoadSection{}
		roadGroupToRoadSection.Id = item.Id
		number, err := strconv.Atoi(item.Number)
		if err != nil {
			number = 0
		}

		roadGroupToRoadSection.Number = number
		roadGroupToRoadSection.Name = "ตอนควบคุม " + item.Number + " " + item.NameOriginTH + " - " + item.NameDestinationTH
		roadGroupToRoadSections[mapString] = append(roadGroupToRoadSections[mapString], roadGroupToRoadSection)
		roadSectionDuplicate[strconv.Itoa(item.RoadGroupId)+"_"+strconv.Itoa(item.Id)] = true
	}

	var assetRefAssets []responses.AssetRefAsset
	for _, item1 := range filterAsset {
		var assetRefAsset responses.AssetRefAsset
		assetRefAsset.Id = item1.ID
		assetRefAsset.Name = item1.Name
		var assetRefAssetTables []responses.AssetRefAssetTable
		for _, item2 := range item1.RefAssetTable {
			var assetRefAssetTable responses.AssetRefAssetTable
			assetRefAssetTable.Id = item2.ID
			assetRefAssetTable.Name = item2.TableLabel
			assetRefAssetTables = append(assetRefAssetTables, assetRefAssetTable)
		}

		if len(assetRefAssetTables) != 0 {
			assetRefAsset.Asset = assetRefAssetTables
			assetRefAssets = append(assetRefAssets, assetRefAsset)
		}
	}

	for index1, item1 := range depots {

		if len(depotToRoadGroups[item1.Id]) == 0 {
			depots[index1].RoadGroup = []responses.FilterRoadGroup{}
		} else {
			depots[index1].RoadGroup = depotToRoadGroups[item1.Id]
		}

		for index2, item2 := range depots[index1].RoadGroup {
			mapString := fmt.Sprintf(`%d_%d`, item1.Id, item2.Id)
			if len(roadGroupToRoadSections[mapString]) == 0 {
				depots[index1].RoadGroup[index2].RoadSection = []responses.FilterRoadSection{}
			} else {
				sort.SliceStable(roadGroupToRoadSections[mapString], func(i, j int) bool {
					return roadGroupToRoadSections[mapString][i].Number < roadGroupToRoadSections[mapString][j].Number
				})
				depots[index1].RoadGroup[index2].RoadSection = roadGroupToRoadSections[mapString]
			}
		}

	}

	var asset responses.FilterAssetType3
	asset.FilterRoad = depots

	return asset, nil
}

func (u *UseCase) FilterRoadType1(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"download_all_road_report"}, []string{"download_all_road_report"})
	//============ end check permission ============
	filterRoadGroup, err := u.Repo.FilterRoadGroup(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}

	var roadGroupNoSections []responses.FilterRoadGroupNoSection
	for _, item := range filterRoadGroup {
		var roadGroupNoSection responses.FilterRoadGroupNoSection

		number, err := strconv.Atoi(item.Number)
		if err != nil {
			number = 0
		}

		roadGroupNoSection.Name = "ทางหลวงหมายเลข " + strconv.Itoa(number) + " " + item.ShortName
		roadGroupNoSection.Id = item.Id

		roadGroupNoSections = append(roadGroupNoSections, roadGroupNoSection)
	}

	var filterRoadType1 responses.FilterRoadType1
	filterRoadType1.FilterRoad = roadGroupNoSections

	return filterRoadType1, nil
}

func (u *UseCase) FilterRoadType2(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"download_all_surface_report"}, []string{"download_owner_surface_report"})
	//============ end check permission ============
	filterRoadSurface, err := u.Repo.FilterRoadSurface(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}
	// return filterRoadSurface, nil
	years := []responses.FilterRoadYear{}
	yearDuplicate := map[int]bool{}
	depotDuplicate := map[string]bool{}
	roadGroupDuplicate := map[string]bool{}
	roadSectionDuplicate := map[string]bool{}
	yearToDepot := map[int][]responses.FilterDepot{}
	depotToRoadGroups := map[int][]responses.FilterRoadGroup{}
	roadGroupToRoadSections := map[int][]responses.FilterRoadSection{}
	for _, item := range filterRoadSurface {

		if !depotDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode] {
			if item.Road.RoadSection.RefDepot.Name == "" {
				continue
			}
			yearToRoadGroup := responses.FilterDepot{}
			yearToRoadGroup.Id = item.Road.RoadSection.RefDepot.Id
			yearToRoadGroup.Name = item.Road.RoadSection.RefDepot.Name
			depotDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode] = true
			yearToDepot[item.Year] = append(yearToDepot[item.Year], yearToRoadGroup)
		}

		if !yearDuplicate[item.Year] {
			year := responses.FilterRoadYear{}
			year.Year = item.Year
			yearDuplicate[item.Year] = true
			if yearToDepot != nil {
				years = append(years, year)
			}

		}

		if !roadGroupDuplicate[item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] {
			number, err := strconv.Atoi(item.Road.RoadSection.RoadGroup.Number)
			if err != nil {
				number = 0
			}

			depotToRoadGroup := responses.FilterRoadGroup{}
			depotToRoadGroup.Id = item.Road.RoadSection.RoadGroup.Id
			depotToRoadGroup.Name = "ทางหลวงหมายเลข " + strconv.Itoa(number) + " " + item.Road.RoadSection.RoadGroup.ShortName
			depotToRoadGroups[item.Road.RoadSection.RefDepot.Id] = append(depotToRoadGroups[item.Road.RoadSection.RefDepot.Id], depotToRoadGroup)
			roadGroupDuplicate[item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}

		if !roadSectionDuplicate[item.Road.RoadSection.Number+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] {
			roadGroupToRoadSection := responses.FilterRoadSection{}
			roadGroupToRoadSection.Id = item.Road.RoadSection.Id

			number, err := strconv.Atoi(item.Road.RoadSection.Number)
			if err != nil {
				number = 0
			}

			roadGroupToRoadSection.Number = number
			roadGroupToRoadSection.Name = "ตอนควบคุม " + item.Road.RoadSection.Number + " " + item.Road.RoadSection.NameOriginTH + " - " + item.Road.RoadSection.NameDestinationTH
			roadGroupToRoadSections[item.Road.RoadSection.RoadGroup.Id] = append(roadGroupToRoadSections[item.Road.RoadSection.RoadGroup.Id], roadGroupToRoadSection)
			roadSectionDuplicate[item.Road.RoadSection.Number+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}
	}

	for index1, item1 := range years {
		if len(yearToDepot[item1.Year]) != 0 {
			years[index1].Depot = yearToDepot[item1.Year]
		}

		for index2, item2 := range years[index1].Depot {
			if len(depotToRoadGroups[item2.Id]) != 0 {
				years[index1].Depot[index2].RoadGroup = depotToRoadGroups[item2.Id]
			}

			for index3, item3 := range years[index1].Depot[index2].RoadGroup {
				if len(roadGroupToRoadSections[item3.Id]) != 0 {
					sort.SliceStable(roadGroupToRoadSections[item3.Id], func(i, j int) bool {
						return roadGroupToRoadSections[item3.Id][i].Number < roadGroupToRoadSections[item3.Id][j].Number
					})
					years[index1].Depot[index2].RoadGroup[index3].RoadSection = roadGroupToRoadSections[item3.Id]
				}
			}
		}
	}
	var filterRoadType2 responses.FilterRoadType2
	filterRoadType2.FilterRoad = years
	return filterRoadType2, nil
}

func (u *UseCase) FilterRoadType3(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"download_all_road_condition_report"}, []string{"download_owner_road_condition_report"})
	//============ end check permission ============
	filterRoadCondition, err := u.Repo.FilterRoadCondition(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}
	years := []responses.FilterRoadYear{}
	yearDuplicate := map[int]bool{}
	depotDuplicate := map[string]bool{}
	roadGroupDuplicate := map[string]bool{}
	roadSectionDuplicate := map[string]bool{}
	yearToDepot := map[int][]responses.FilterDepot{}
	depotToRoadGroups := map[int][]responses.FilterRoadGroup{}
	roadGroupToRoadSections := map[int][]responses.FilterRoadSection{}

	for _, item := range filterRoadCondition {
		if !depotDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode] {
			if item.Road.RoadSection.RefDepot.Name == "" {
				continue
			}
			yearToRoadGroup := responses.FilterDepot{}
			yearToRoadGroup.Id = item.Road.RoadSection.RefDepot.Id
			yearToRoadGroup.Name = item.Road.RoadSection.RefDepot.Name
			depotDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode] = true
			yearToDepot[item.Year] = append(yearToDepot[item.Year], yearToRoadGroup)
		}

		if !yearDuplicate[item.Year] {
			year := responses.FilterRoadYear{}
			year.Year = item.Year
			yearDuplicate[item.Year] = true
			if yearToDepot != nil {
				years = append(years, year)
			}
		}

		if !roadGroupDuplicate[item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] {
			number, err := strconv.Atoi(item.Road.RoadSection.RoadGroup.Number)
			if err != nil {
				number = 0
			}

			depotToRoadGroup := responses.FilterRoadGroup{}
			depotToRoadGroup.Id = item.Road.RoadSection.RoadGroup.Id
			depotToRoadGroup.Name = "ทางหลวงหมายเลข " + strconv.Itoa(number) + " " + item.Road.RoadSection.RoadGroup.ShortName
			depotToRoadGroups[item.Road.RoadSection.RefDepot.Id] = append(depotToRoadGroups[item.Road.RoadSection.RefDepot.Id], depotToRoadGroup)
			roadGroupDuplicate[item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}

		if roadSectionDuplicate[item.Road.RoadSection.Number+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] == false {
			roadGroupToRoadSection := responses.FilterRoadSection{}
			roadGroupToRoadSection.Id = item.Road.RoadSection.Id

			number, err := strconv.Atoi(item.Road.RoadSection.Number)
			if err != nil {
				number = 0
			}

			roadGroupToRoadSection.Number = number
			roadGroupToRoadSection.Name = "ตอนควบคุม " + item.Road.RoadSection.Number + " " + item.Road.RoadSection.NameOriginTH + " - " + item.Road.RoadSection.NameDestinationTH
			roadGroupToRoadSections[item.Road.RoadSection.RoadGroup.Id] = append(roadGroupToRoadSections[item.Road.RoadSection.RoadGroup.Id], roadGroupToRoadSection)
			roadSectionDuplicate[item.Road.RoadSection.Number+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}
	}

	for index1, item1 := range years {
		if len(yearToDepot[item1.Year]) != 0 {
			years[index1].Depot = yearToDepot[item1.Year]
		}

		for index2, item2 := range years[index1].Depot {
			if len(depotToRoadGroups[item2.Id]) != 0 {
				years[index1].Depot[index2].RoadGroup = depotToRoadGroups[item2.Id]
			}

			for index3, item3 := range years[index1].Depot[index2].RoadGroup {
				if len(roadGroupToRoadSections[item3.Id]) != 0 {
					sort.SliceStable(roadGroupToRoadSections[item3.Id], func(i, j int) bool {
						return roadGroupToRoadSections[item3.Id][i].Number < roadGroupToRoadSections[item3.Id][j].Number
					})
					years[index1].Depot[index2].RoadGroup[index3].RoadSection = roadGroupToRoadSections[item3.Id]
				}
			}
		}
	}

	var filterRoadType3 responses.FilterRoadType3
	filterRoadType3.FilterRoad = years
	filterRoadType3.FilterRange = []int{25, 100, 1000}

	return filterRoadType3, nil
}

func (u *UseCase) FilterRoadType4(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"download_all_road_condition_result_report"}, []string{"download_owner_road_condition_result_report"})
	//============ end check permission ============
	filterRoadCondition, err := u.Repo.FilterRoadCondition(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}

	filterRefOwner, err := u.Repo.FilterRefOwner()
	if err != nil {
		return nil, err
	}

	var filterRefOwners []responses.FilterRefOwner
	for _, item := range filterRefOwner {
		var filterRefOwner responses.FilterRefOwner
		filterRefOwner.Id = item.ID
		filterRefOwner.Name = item.Name
		filterRefOwners = append(filterRefOwners, filterRefOwner)
	}

	years := []responses.FilterRoadYear{}
	yearDuplicate := map[int]bool{}
	depotDuplicate := map[string]bool{}
	roadGroupDuplicate := map[string]bool{}
	roadSectionDuplicate := map[string]bool{}
	yearToDepot := map[int][]responses.FilterDepot{}
	depotToRoadGroups := map[int][]responses.FilterRoadGroup{}
	roadGroupToRoadSections := map[int][]responses.FilterRoadSection{}

	for _, item := range filterRoadCondition {
		if !depotDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode] {
			if item.Road.RoadSection.RefDepot.Name == "" {
				continue
			}
			yearToRoadGroup := responses.FilterDepot{}
			yearToRoadGroup.Id = item.Road.RoadSection.RefDepot.Id
			yearToRoadGroup.Name = item.Road.RoadSection.RefDepot.Name
			depotDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode] = true
			yearToDepot[item.Year] = append(yearToDepot[item.Year], yearToRoadGroup)
		}

		if !yearDuplicate[item.Year] {
			year := responses.FilterRoadYear{}
			year.Year = item.Year
			yearDuplicate[item.Year] = true
			if yearToDepot != nil {
				years = append(years, year)
			}
		}

		if !roadGroupDuplicate[item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] {
			number, err := strconv.Atoi(item.Road.RoadSection.RoadGroup.Number)
			if err != nil {
				number = 0
			}

			depotToRoadGroup := responses.FilterRoadGroup{}
			depotToRoadGroup.Id = item.Road.RoadSection.RoadGroup.Id
			depotToRoadGroup.Name = "ทางหลวงหมายเลข " + strconv.Itoa(number) + " " + item.Road.RoadSection.RoadGroup.ShortName
			depotToRoadGroups[item.Road.RoadSection.RefDepot.Id] = append(depotToRoadGroups[item.Road.RoadSection.RefDepot.Id], depotToRoadGroup)
			roadGroupDuplicate[item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}

		if roadSectionDuplicate[item.Road.RoadSection.Number+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] == false {
			roadGroupToRoadSection := responses.FilterRoadSection{}
			roadGroupToRoadSection.Id = item.Road.RoadSection.Id
			number, err := strconv.Atoi(item.Road.RoadSection.Number)
			if err != nil {
				number = 0
			}

			roadGroupToRoadSection.Number = number
			roadGroupToRoadSection.Name = "ตอนควบคุม " + item.Road.RoadSection.Number + " " + item.Road.RoadSection.NameOriginTH + " - " + item.Road.RoadSection.NameDestinationTH
			roadGroupToRoadSections[item.Road.RoadSection.RoadGroup.Id] = append(roadGroupToRoadSections[item.Road.RoadSection.RoadGroup.Id], roadGroupToRoadSection)
			roadSectionDuplicate[item.Road.RoadSection.Number+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}
	}

	for index1, item1 := range years {
		if len(yearToDepot[item1.Year]) != 0 {
			years[index1].Depot = yearToDepot[item1.Year]
		}

		for index2, item2 := range years[index1].Depot {
			if len(depotToRoadGroups[item2.Id]) != 0 {
				years[index1].Depot[index2].RoadGroup = depotToRoadGroups[item2.Id]
			}

			for index3, item3 := range years[index1].Depot[index2].RoadGroup {
				if len(roadGroupToRoadSections[item3.Id]) != 0 {
					sort.SliceStable(roadGroupToRoadSections[item3.Id], func(i, j int) bool {
						return roadGroupToRoadSections[item3.Id][i].Number < roadGroupToRoadSections[item3.Id][j].Number
					})
					years[index1].Depot[index2].RoadGroup[index3].RoadSection = roadGroupToRoadSections[item3.Id]
				}
			}
		}
	}

	var filterRoadType4 responses.FilterRoadType4
	filterRoadType4.FilterRoad = years
	filterRoadType4.FilterCondition = []string{"IRI", "MPD", "RUT", "IFI"}
	filterRoadType4.FilterCriteria = filterRefOwners

	return filterRoadType4, nil
}

func (u *UseCase) FilterRoadType5(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"download_all_road_retro_report"}, []string{"download_owner_road_retro_report"})
	//============ end check permission ============
	filterRoadCondition, err := u.Repo.FilterRoadRetroReflectivity(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}

	filterRefOwner, err := u.Repo.FilterRefOwnerLine()
	if err != nil {
		return nil, err
	}

	var filterRefOwnerLines []responses.FilterRefOwnerLine
	for _, item := range filterRefOwner {
		var ilterRefOwnerLine responses.FilterRefOwnerLine
		ilterRefOwnerLine.Id = item.ID
		ilterRefOwnerLine.Name = item.Name
		filterRefOwnerLines = append(filterRefOwnerLines, ilterRefOwnerLine)
	}

	years := []responses.FilterRoadYear{}
	yearDuplicate := map[int]bool{}
	depotDuplicate := map[string]bool{}
	roadGroupDuplicate := map[string]bool{}
	roadSectionDuplicate := map[string]bool{}
	yearToDepot := map[int][]responses.FilterDepot{}
	depotToRoadGroups := map[int][]responses.FilterRoadGroup{}
	roadGroupToRoadSections := map[string][]responses.FilterRoadSection{}

	for _, item := range filterRoadCondition {
		if item.Road.RoadSection.RefDepot.DepotCode == "" {
			continue
		}
		if !depotDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode] {
			yearToRoadGroup := responses.FilterDepot{}
			yearToRoadGroup.Id = item.Road.RoadSection.RefDepot.Id
			yearToRoadGroup.Name = item.Road.RoadSection.RefDepot.Name
			depotDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode] = true
			yearToDepot[item.Year] = append(yearToDepot[item.Year], yearToRoadGroup)
		}

		if !yearDuplicate[item.Year] {
			year := responses.FilterRoadYear{}
			year.Year = item.Year
			yearDuplicate[item.Year] = true
			if yearToDepot != nil {
				years = append(years, year)
			}
		}

		if !roadGroupDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] {

			number, err := strconv.Atoi(item.Road.RoadSection.RoadGroup.Number)
			if err != nil {
				number = 0
				strconv.Itoa(number)
			}

			depotToRoadGroup := responses.FilterRoadGroup{}
			depotToRoadGroup.Id = item.Road.RoadSection.RoadGroup.Id
			depotToRoadGroup.Name = "ทางหลวงหมายเลข " + strconv.Itoa(number) + " " + item.Road.RoadSection.RoadGroup.ShortName
			depotToRoadGroups[item.Road.RoadSection.RefDepot.Id] = append(depotToRoadGroups[item.Road.RoadSection.RefDepot.Id], depotToRoadGroup)
			roadGroupDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}

		if !roadSectionDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)+"_"+strconv.Itoa(item.Road.RoadSection.Id)] {
			roadGroupToRoadSection := responses.FilterRoadSection{}
			roadGroupToRoadSection.Id = item.Road.RoadSection.Id
			number, err := strconv.Atoi(item.Road.RoadSection.Number)
			if err != nil {
				number = 0
			}

			roadGroupToRoadSection.Number = number
			roadGroupToRoadSection.Name = "ตอนควบคุม " + item.Road.RoadSection.Number + " " + item.Road.RoadSection.NameOriginTH + " - " + item.Road.RoadSection.NameDestinationTH
			roadGroupToRoadSections[strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)+"_"+strconv.Itoa(item.Year)] = append(roadGroupToRoadSections[strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)+"_"+strconv.Itoa(item.Year)], roadGroupToRoadSection)
			roadSectionDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)+"_"+strconv.Itoa(item.Road.RoadSection.Id)] = true
		}

	}

	for index1, item1 := range years {
		if len(yearToDepot[item1.Year]) != 0 {
			years[index1].Depot = yearToDepot[item1.Year]
		}

		for index2, item2 := range years[index1].Depot {
			if len(depotToRoadGroups[item2.Id]) != 0 {
				years[index1].Depot[index2].RoadGroup = depotToRoadGroups[item2.Id]
			}

			for index3, item3 := range years[index1].Depot[index2].RoadGroup {
				if len(roadGroupToRoadSections[strconv.Itoa(item3.Id)+"_"+strconv.Itoa(item1.Year)]) != 0 {
					sort.SliceStable(roadGroupToRoadSections[strconv.Itoa(item3.Id)+"_"+strconv.Itoa(item1.Year)], func(i, j int) bool {
						return roadGroupToRoadSections[strconv.Itoa(item3.Id)+"_"+strconv.Itoa(item1.Year)][i].Number < roadGroupToRoadSections[strconv.Itoa(item3.Id)+"_"+strconv.Itoa(item1.Year)][j].Number
					})
					years[index1].Depot[index2].RoadGroup[index3].RoadSection = roadGroupToRoadSections[strconv.Itoa(item3.Id)+"_"+strconv.Itoa(item1.Year)]
				}

			}
		}
	}

	var filterRoadType5 responses.FilterRoadType5
	filterRoadType5.FilterRoad = years
	filterRoadType5.FilterCriteria = filterRefOwnerLines

	return filterRoadType5, nil
}

func (u *UseCase) FilterRoadType6(userID int) (interface{}, error) {
	// ============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"download_all_road_retro_result_report"}, []string{"download_owner_road_retro_result_report"})
	//============ end check permission ============
	filterRoadCondition, err := u.Repo.FilterRoadRetroReflectivity(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}

	filterRefOwner, err := u.Repo.FilterRefOwnerLine()
	if err != nil {
		return nil, err
	}

	var filterRefOwnerLines []responses.FilterRefOwnerLine
	for _, item := range filterRefOwner {
		var ilterRefOwnerLine responses.FilterRefOwnerLine
		ilterRefOwnerLine.Id = item.ID
		ilterRefOwnerLine.Name = item.Name
		filterRefOwnerLines = append(filterRefOwnerLines, ilterRefOwnerLine)
	}

	years := []responses.FilterRoadYear{}
	yearDuplicate := map[int]bool{}
	depotDuplicate := map[string]bool{}
	roadGroupDuplicate := map[string]bool{}
	roadSectionDuplicate := map[string]bool{}
	yearToDepot := map[int][]responses.FilterDepot{}
	depotToRoadGroups := map[int][]responses.FilterRoadGroup{}
	roadGroupToRoadSections := map[int][]responses.FilterRoadSection{}

	for _, item := range filterRoadCondition {
		if item.Road.RoadSection.RefDepot.DepotCode == "" {
			continue
		}
		if !depotDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode] {
			yearToRoadGroup := responses.FilterDepot{}
			yearToRoadGroup.Id = item.Road.RoadSection.RefDepot.Id
			yearToRoadGroup.Name = item.Road.RoadSection.RefDepot.Name
			depotDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode] = true
			yearToDepot[item.Year] = append(yearToDepot[item.Year], yearToRoadGroup)
		}

		if !yearDuplicate[item.Year] {
			year := responses.FilterRoadYear{}
			year.Year = item.Year
			yearDuplicate[item.Year] = true
			if yearToDepot != nil {
				years = append(years, year)
			}
		}

		if !roadGroupDuplicate[item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] {
			number, err := strconv.Atoi(item.Road.RoadSection.RoadGroup.Number)
			if err != nil {
				number = 0
			}
			depotToRoadGroup := responses.FilterRoadGroup{}
			depotToRoadGroup.Id = item.Road.RoadSection.RoadGroup.Id
			depotToRoadGroup.Name = "ทางหลวงหมายเลข " + strconv.Itoa(number) + " " + item.Road.RoadSection.RoadGroup.ShortName
			depotToRoadGroups[item.Road.RoadSection.RefDepot.Id] = append(depotToRoadGroups[item.Road.RoadSection.RefDepot.Id], depotToRoadGroup)
			roadGroupDuplicate[item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}

		if !roadSectionDuplicate[item.Road.RoadSection.Number+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] {
			roadGroupToRoadSection := responses.FilterRoadSection{}
			roadGroupToRoadSection.Id = item.Road.RoadSection.Id
			number, err := strconv.Atoi(item.Road.RoadSection.Number)
			if err != nil {
				number = 0
			}

			roadGroupToRoadSection.Number = number
			roadGroupToRoadSection.Name = "ตอนควบคุม " + item.Road.RoadSection.Number + " " + item.Road.RoadSection.NameOriginTH + " - " + item.Road.RoadSection.NameDestinationTH
			roadGroupToRoadSections[item.Road.RoadSection.RoadGroup.Id] = append(roadGroupToRoadSections[item.Road.RoadSection.RoadGroup.Id], roadGroupToRoadSection)
			roadSectionDuplicate[item.Road.RoadSection.Number+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}
	}

	for index1, item1 := range years {
		if len(yearToDepot[item1.Year]) != 0 {
			years[index1].Depot = yearToDepot[item1.Year]
		}

		for index2, item2 := range years[index1].Depot {
			if len(depotToRoadGroups[item2.Id]) != 0 {
				years[index1].Depot[index2].RoadGroup = depotToRoadGroups[item2.Id]
			}

			for index3, item3 := range years[index1].Depot[index2].RoadGroup {
				if len(roadGroupToRoadSections[item3.Id]) != 0 {
					sort.SliceStable(roadGroupToRoadSections[item3.Id], func(i, j int) bool {
						return roadGroupToRoadSections[item3.Id][i].Number < roadGroupToRoadSections[item3.Id][j].Number
					})
					years[index1].Depot[index2].RoadGroup[index3].RoadSection = roadGroupToRoadSections[item3.Id]
				}
			}
		}
	}

	var filterRoadType6 responses.FilterRoadType6
	filterRoadType6.FilterRoad = years
	filterRoadType6.FilterCriteria = filterRefOwnerLines

	return filterRoadType6, nil
}

func (u *UseCase) FiltertRoadDamageType1(userID int) (interface{}, error) {
	// ============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"download_all_road_damage_report"}, []string{"download_owner_road_damage_report"})
	//============ end check permission ============
	filterRoadDamage, err := u.Repo.FilterRoadDamage(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}

	years := []responses.FilterRoadYear{}
	yearDuplicate := map[int]bool{}
	depotDuplicate := map[string]bool{}
	roadGroupDuplicate := map[string]bool{}
	roadSectionDuplicate := map[string]bool{}
	yearToDepot := map[int][]responses.FilterDepot{}
	depotToRoadGroups := map[int][]responses.FilterRoadGroup{}
	roadGroupToRoadSections := map[int][]responses.FilterRoadSection{}

	for _, item := range filterRoadDamage {
		if !depotDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode] {
			if item.Road.RoadSection.RefDepot.DepotCode == "" {
				continue
			}
			yearToRoadGroup := responses.FilterDepot{}
			yearToRoadGroup.Id = item.Road.RoadSection.RefDepot.Id
			yearToRoadGroup.Name = item.Road.RoadSection.RefDepot.Name
			depotDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode] = true
			yearToDepot[item.Year] = append(yearToDepot[item.Year], yearToRoadGroup)
		}

		if !yearDuplicate[item.Year] {
			year := responses.FilterRoadYear{}
			year.Year = item.Year
			yearDuplicate[item.Year] = true
			if yearToDepot != nil {
				years = append(years, year)
			}
		}

		if !roadGroupDuplicate[item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] {
			number, err := strconv.Atoi(item.Road.RoadSection.RoadGroup.Number)
			if err != nil {
				number = 0

			}

			depotToRoadGroup := responses.FilterRoadGroup{}
			depotToRoadGroup.Id = item.Road.RoadSection.RoadGroup.Id
			depotToRoadGroup.Name = "ทางหลวงหมายเลข " + strconv.Itoa(number) + " " + item.Road.RoadSection.RoadGroup.ShortName
			depotToRoadGroups[item.Road.RoadSection.RefDepot.Id] = append(depotToRoadGroups[item.Road.RoadSection.RefDepot.Id], depotToRoadGroup)
			roadGroupDuplicate[item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}

		if !roadSectionDuplicate[item.Road.RoadSection.Number+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] {
			roadGroupToRoadSection := responses.FilterRoadSection{}
			roadGroupToRoadSection.Id = item.Road.RoadSection.Id
			number, err := strconv.Atoi(item.Road.RoadSection.Number)
			if err != nil {
				number = 0
			}

			roadGroupToRoadSection.Number = number
			roadGroupToRoadSection.Name = "ตอนควบคุม " + item.Road.RoadSection.Number + " " + item.Road.RoadSection.NameOriginTH + " - " + item.Road.RoadSection.NameDestinationTH
			roadGroupToRoadSections[item.Road.RoadSection.RoadGroup.Id] = append(roadGroupToRoadSections[item.Road.RoadSection.RoadGroup.Id], roadGroupToRoadSection)
			roadSectionDuplicate[item.Road.RoadSection.Number+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}
	}

	for index1, item1 := range years {
		if len(yearToDepot[item1.Year]) != 0 {
			years[index1].Depot = yearToDepot[item1.Year]
		}

		for index2, item2 := range years[index1].Depot {
			if len(depotToRoadGroups[item2.Id]) != 0 {
				years[index1].Depot[index2].RoadGroup = depotToRoadGroups[item2.Id]
			}

			for index3, item3 := range years[index1].Depot[index2].RoadGroup {
				if len(roadGroupToRoadSections[item3.Id]) != 0 {
					sort.SliceStable(roadGroupToRoadSections[item3.Id], func(i, j int) bool {
						return roadGroupToRoadSections[item3.Id][i].Number < roadGroupToRoadSections[item3.Id][j].Number
					})
					years[index1].Depot[index2].RoadGroup[index3].RoadSection = roadGroupToRoadSections[item3.Id]
				}
			}
		}
	}

	var filterRoadDamageType1 responses.FilterRoadDamageType1
	filterRoadDamageType1.FilterRoad = years

	return filterRoadDamageType1, nil
}

func (u *UseCase) FiltertRoadDamageType2(userID int) (interface{}, error) {
	// ============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"download_all_road_damage_result_report"}, []string{"download_owner_road_damage_result_report"})
	// ============ end check permission ============
	filterRoadDamage, err := u.Repo.FilterRoadDamage(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}

	years := []responses.FilterRoadYear{}
	yearDuplicate := map[int]bool{}
	depotDuplicate := map[string]bool{}
	roadGroupDuplicate := map[string]bool{}
	roadSectionDuplicate := map[string]bool{}
	yearToDepot := map[int][]responses.FilterDepot{}
	depotToRoadGroups := map[int][]responses.FilterRoadGroup{}
	roadGroupToRoadSections := map[int][]responses.FilterRoadSection{}

	for _, item := range filterRoadDamage {
		if !depotDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode] {
			if item.Road.RoadSection.RefDepot.Name == "" {
				continue
			}
			yearToRoadGroup := responses.FilterDepot{}
			yearToRoadGroup.Id = item.Road.RoadSection.RefDepot.Id
			yearToRoadGroup.Name = item.Road.RoadSection.RefDepot.Name
			depotDuplicate[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RefDepot.DepotCode] = true
			yearToDepot[item.Year] = append(yearToDepot[item.Year], yearToRoadGroup)
		}

		if !yearDuplicate[item.Year] {
			year := responses.FilterRoadYear{}
			year.Year = item.Year
			yearDuplicate[item.Year] = true
			if yearToDepot != nil {
				years = append(years, year)
			}
		}

		if !roadGroupDuplicate[item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] {
			number, err := strconv.Atoi(item.Road.RoadSection.RoadGroup.Number)
			if err != nil {
				number = 0
			}
			depotToRoadGroup := responses.FilterRoadGroup{}
			depotToRoadGroup.Id = item.Road.RoadSection.RoadGroup.Id
			depotToRoadGroup.Name = "ทางหลวงหมายเลข " + strconv.Itoa(number) + " " + item.Road.RoadSection.RoadGroup.ShortName
			depotToRoadGroups[item.Road.RoadSection.RefDepot.Id] = append(depotToRoadGroups[item.Road.RoadSection.RefDepot.Id], depotToRoadGroup)
			roadGroupDuplicate[item.Road.RoadSection.RefDepot.DepotCode+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}

		if !roadSectionDuplicate[item.Road.RoadSection.Number+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] {
			roadGroupToRoadSection := responses.FilterRoadSection{}
			roadGroupToRoadSection.Id = item.Road.RoadSection.Id
			number, err := strconv.Atoi(item.Road.RoadSection.Number)
			if err != nil {
				number = 0
			}

			roadGroupToRoadSection.Number = number
			roadGroupToRoadSection.Name = "ตอนควบคุม " + item.Road.RoadSection.Number + " " + item.Road.RoadSection.NameOriginTH + " - " + item.Road.RoadSection.NameDestinationTH
			roadGroupToRoadSections[item.Road.RoadSection.RoadGroup.Id] = append(roadGroupToRoadSections[item.Road.RoadSection.RoadGroup.Id], roadGroupToRoadSection)
			roadSectionDuplicate[item.Road.RoadSection.Number+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}
	}

	for index1, item1 := range years {
		if len(yearToDepot[item1.Year]) != 0 {
			years[index1].Depot = yearToDepot[item1.Year]
		}

		for index2, item2 := range years[index1].Depot {
			if len(depotToRoadGroups[item2.Id]) != 0 {
				years[index1].Depot[index2].RoadGroup = depotToRoadGroups[item2.Id]
			}

			for index3, item3 := range years[index1].Depot[index2].RoadGroup {
				if len(roadGroupToRoadSections[item3.Id]) != 0 {
					sort.SliceStable(roadGroupToRoadSections[item3.Id], func(i, j int) bool {
						return roadGroupToRoadSections[item3.Id][i].Number < roadGroupToRoadSections[item3.Id][j].Number
					})
					years[index1].Depot[index2].RoadGroup[index3].RoadSection = roadGroupToRoadSections[item3.Id]
				}
			}
		}
	}

	var filterRoadDamageType2 responses.FilterRoadDamageType2
	filterRoadDamageType2.FilterRoad = years

	return filterRoadDamageType2, nil
}

func (u *UseCase) FilterMaintenanceKpiType1(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"download_all_maint_kpi_report"}, []string{"download_owner_maint_kpi_report"})
	//============ end check permission ============

	condition := []string{"IRI", "RUT", "IFI", "G7"}

	filterRoadCondition, err := u.Repo.FilterRoadCondition(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}

	years := []int{}
	yearDuplicate := map[int]bool{}
	yearsToRoadGroups := map[string][]responses.FilterRoadGroup{}
	yearsToRoadGroupDuplicate := map[string]bool{}
	roadGroupToRoadSection := map[string][]responses.FilterRoadSection{}
	roadGroupToRoadSectionDuplicate := map[string]bool{}
	for _, item := range filterRoadCondition {
		if !yearsToRoadGroupDuplicate[strconv.Itoa(item.Year)+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] {
			if item.Road.RoadSection.RoadGroup.ShortName == "" {
				continue
			}
			var filterRoadGroup responses.FilterRoadGroup
			filterRoadGroup.Id = item.Road.RoadSection.RoadGroup.Id
			filterRoadGroup.Name = "ทางหลวงหมายเลข " + item.Road.RoadSection.RoadGroup.Number + " " + item.Road.RoadSection.RoadGroup.ShortName
			yearsToRoadGroups[strconv.Itoa(item.Year)] = append(yearsToRoadGroups[strconv.Itoa(item.Year)], filterRoadGroup)
			yearsToRoadGroupDuplicate[strconv.Itoa(item.Year)+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}

		if !yearDuplicate[item.Year] {
			years = append(years, item.Year)
			yearDuplicate[item.Year] = true
		}

		if !roadGroupToRoadSectionDuplicate[strconv.Itoa(item.Year)+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)+"_"+strconv.Itoa(item.Road.RoadSection.Id)] {
			var filterRoadSection responses.FilterRoadSection
			filterRoadSection.Id = item.Road.RoadSection.Id

			number, err := strconv.Atoi(item.Road.RoadSection.Number)
			if err != nil {
				number = 0
			}
			filterRoadSection.Number = number

			filterRoadSection.Name = "ตอนควบคุม " + item.Road.RoadSection.Number + " " + item.Road.RoadSection.NameOriginTH + " - " + item.Road.RoadSection.NameDestinationTH
			roadGroupToRoadSection[strconv.Itoa(item.Year)+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = append(roadGroupToRoadSection[strconv.Itoa(item.Year)+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)], filterRoadSection)
			roadGroupToRoadSectionDuplicate[strconv.Itoa(item.Year)+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)+"_"+strconv.Itoa(item.Road.RoadSection.Id)] = true
		}

	}

	var filterMaintenanceKpisRoad []responses.FilterMaintenanceYearKpi
	for _, itemYear := range years {
		var filterMaintenanceKpi responses.FilterMaintenanceYearKpi
		filterMaintenanceKpi.Year = itemYear
		filterMaintenanceKpi.RoadGroup = yearsToRoadGroups[strconv.Itoa(itemYear)]

		for indexRoadGroup, itemRoadGroup := range yearsToRoadGroups[strconv.Itoa(itemYear)] {

			data := roadGroupToRoadSection[strconv.Itoa(itemYear)+"_"+strconv.Itoa(itemRoadGroup.Id)]
			sort.SliceStable(data, func(i, j int) bool {
				return data[i].Number < data[j].Number
			})
			filterMaintenanceKpi.RoadGroup[indexRoadGroup].RoadSection = data
		}

		filterMaintenanceKpisRoad = append(filterMaintenanceKpisRoad, filterMaintenanceKpi)
	}

	filterRoadRetroReflectivity, err := u.Repo.FilterRoadRetroReflectivity(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}

	years = []int{}
	yearDuplicate = map[int]bool{}
	yearsToRoadGroups = map[string][]responses.FilterRoadGroup{}
	yearsToRoadGroupDuplicate = map[string]bool{}
	roadGroupToRoadSection = map[string][]responses.FilterRoadSection{}
	roadGroupToRoadSectionDuplicate = map[string]bool{}
	for _, item := range filterRoadRetroReflectivity {

		if !yearDuplicate[item.Year] {
			years = append(years, item.Year)
			yearDuplicate[item.Year] = true
		}

		if !yearsToRoadGroupDuplicate[strconv.Itoa(item.Year)+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] {
			var filterRoadGroup responses.FilterRoadGroup
			filterRoadGroup.Id = item.Road.RoadSection.RoadGroup.Id
			filterRoadGroup.Name = "ทางหลวงหมายเลข " + item.Road.RoadSection.RoadGroup.Number + " " + item.Road.RoadSection.RoadGroup.ShortName
			yearsToRoadGroups[strconv.Itoa(item.Year)] = append(yearsToRoadGroups[strconv.Itoa(item.Year)], filterRoadGroup)
			yearsToRoadGroupDuplicate[strconv.Itoa(item.Year)+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = true
		}

		if !roadGroupToRoadSectionDuplicate[strconv.Itoa(item.Year)+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)+"_"+strconv.Itoa(item.Road.RoadSection.Id)] {
			var filterRoadSection responses.FilterRoadSection
			filterRoadSection.Id = item.Road.RoadSection.Id
			number, err := strconv.Atoi(item.Road.RoadSection.Number)
			if err != nil {
				number = 0
			}

			filterRoadSection.Number = number
			filterRoadSection.Name = "ตอนควบคุม " + item.Road.RoadSection.Number + " " + item.Road.RoadSection.NameOriginTH + " - " + item.Road.RoadSection.NameDestinationTH
			roadGroupToRoadSection[strconv.Itoa(item.Year)+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = append(roadGroupToRoadSection[strconv.Itoa(item.Year)+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)], filterRoadSection)
			roadGroupToRoadSectionDuplicate[strconv.Itoa(item.Year)+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)+"_"+strconv.Itoa(item.Road.RoadSection.Id)] = true
		}

	}

	var filterMaintenanceKpisReflextivity []responses.FilterMaintenanceYearKpi
	for _, itemYear := range years {
		var filterMaintenanceKpi responses.FilterMaintenanceYearKpi
		filterMaintenanceKpi.Year = itemYear
		filterMaintenanceKpi.RoadGroup = yearsToRoadGroups[strconv.Itoa(itemYear)]

		for indexRoadGroup, itemRoadGroup := range yearsToRoadGroups[strconv.Itoa(itemYear)] {
			filterMaintenanceKpi.RoadGroup[indexRoadGroup].RoadSection = roadGroupToRoadSection[strconv.Itoa(itemYear)+"_"+strconv.Itoa(itemRoadGroup.Id)]
		}

		filterMaintenanceKpisReflextivity = append(filterMaintenanceKpisReflextivity, filterMaintenanceKpi)
	}

	var filterMaintenanceConditionKpis []responses.FilterMaintenanceConditionKpi
	for _, item := range condition {
		var filterMaintenanceConditionKpi responses.FilterMaintenanceConditionKpi
		if item == "IRI" || item == "MPD" || item == "RUT" || item == "IFI" {
			filterMaintenanceConditionKpi.Name = item
			filterMaintenanceConditionKpi.Year = filterMaintenanceKpisRoad
		} else {
			filterMaintenanceConditionKpi.Name = item
			filterMaintenanceConditionKpi.Year = filterMaintenanceKpisReflextivity
		}
		filterMaintenanceConditionKpis = append(filterMaintenanceConditionKpis, filterMaintenanceConditionKpi)
	}

	var filterMaintenanceKpiType1 responses.FilterMaintenanceKpiType1
	filterMaintenanceKpiType1.FilterCondition = filterMaintenanceConditionKpis

	return filterMaintenanceKpiType1, nil
}

func (u *UseCase) FilterMaintenanceFilterType1(userID int) (interface{}, error) {
	// ============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"download_all_maint_history_report"}, []string{"download_owner_maint_history_report"})
	//============ end check permission ============
	maintenance, err := u.Repo.FilterMaintenance(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}

	yearMin := 99999

	roadSectionDuplicate := map[string]bool{}
	roadSection := map[int][]responses.FilterRoadSection{}
	roadGroupDuplicate := map[int]bool{}
	roadGroup := []responses.FilterRoadGroup{}
	for _, itemMaintenance := range maintenance {

		for _, itemMaintenanceRoad := range itemMaintenance.MaintenanceRoad {

			if !roadSectionDuplicate[strconv.Itoa(itemMaintenanceRoad.Road.RoadSection.RoadGroup.Id)+"_"+strconv.Itoa(itemMaintenanceRoad.Road.RoadSection.Id)] {
				if itemMaintenanceRoad.Road.RoadSection.Number == "" {
					continue
				}
				var filterRoadSection responses.FilterRoadSection
				filterRoadSection.Id = itemMaintenanceRoad.Road.RoadSection.Id
				number, err := strconv.Atoi(itemMaintenanceRoad.Road.RoadSection.Number)
				if err != nil {
					number = 0
				}

				filterRoadSection.Number = number
				filterRoadSection.Name = "ตอนควบคุม " + itemMaintenanceRoad.Road.RoadSection.Number + " " + itemMaintenanceRoad.Road.RoadSection.NameOriginTH + " - " + itemMaintenanceRoad.Road.RoadSection.NameDestinationTH
				roadSection[itemMaintenanceRoad.Road.RoadSection.RoadGroup.Id] = append(roadSection[itemMaintenanceRoad.Road.RoadSection.RoadGroup.Id], filterRoadSection)
				roadSectionDuplicate[strconv.Itoa(itemMaintenanceRoad.Road.RoadSection.RoadGroup.Id)+"_"+strconv.Itoa(itemMaintenanceRoad.Road.RoadSection.Id)] = true
			}

			if !roadGroupDuplicate[itemMaintenanceRoad.Road.RoadSection.RoadGroup.Id] {
				var filterRoadGroup responses.FilterRoadGroup
				filterRoadGroup.Id = itemMaintenanceRoad.Road.RoadSection.RoadGroup.Id
				filterRoadGroup.Name = "ทางหลวงหมายเลข " + itemMaintenanceRoad.Road.RoadSection.RoadGroup.Number + " " + itemMaintenanceRoad.Road.RoadSection.RoadGroup.ShortName
				if roadSection != nil {
					roadGroup = append(roadGroup, filterRoadGroup)
				}

				roadGroupDuplicate[itemMaintenanceRoad.Road.RoadSection.RoadGroup.Id] = true
			}

		}
		if len(roadGroup) > 0 {
			if yearMin > itemMaintenance.BudgetYear {
				yearMin = itemMaintenance.BudgetYear
			}
		}

	}

	yearStart := []int{}
	yearEnd := []int{}
	for i := yearMin; i <= time.Now().Year(); i++ {
		yearStart = append(yearStart, i)
		yearEnd = append(yearEnd, i)
	}

	sort.Slice(yearStart, func(i, j int) bool {
		return yearStart[i] < yearStart[j]
	})

	sort.Slice(yearEnd, func(i, j int) bool {
		return yearEnd[i] > yearEnd[j]
	})

	for index, item := range roadGroup {
		sort.SliceStable(roadSection[item.Id], func(i, j int) bool {
			return roadSection[item.Id][i].Number < roadSection[item.Id][j].Number
		})
		roadGroup[index].RoadSection = roadSection[item.Id]
	}

	var filterMaintenanceYears responses.FilterMaintenanceYear
	filterMaintenanceYears.StartYear = yearStart
	filterMaintenanceYears.EndYear = yearEnd

	var filterMaintenanceType1 responses.FilterMaintenanceType1
	filterMaintenanceType1.FilterRoad = roadGroup
	filterMaintenanceType1.FilterYear = filterMaintenanceYears

	return filterMaintenanceType1, nil
}

func (u *UseCase) FilterAadtType1(userID int) (interface{}, error) {
	// ============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"download_all_aadt_report"}, []string{"download_owner_aadt_report"})
	//============ end check permission ============
	aadt, err := u.Repo.FilterAadt(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, err
	}

	roadSectionMaps := map[string][]responses.FilterRoadSection{}
	roadSectionDuplicateMap := map[string]bool{}
	roadGroupMaps := map[int][]responses.FilterRoadGroup{}
	roadGroupIsDuplicateMap := map[string]bool{}
	years := []int{}
	yearIsduplicateMap := map[int]bool{}
	for _, item := range aadt {
		if item.Road.RoadSection.RoadGroup.Number != "" {

			if !roadSectionDuplicateMap[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RoadGroup.Number+"_"+item.Road.RoadSection.Number] {
				if item.Road.RoadSection.Number == "" {
					continue
				}
				var filterRoadSection responses.FilterRoadSection
				filterRoadSection.Id = item.Road.RoadSection.Id
				number, err := strconv.Atoi(item.Road.RoadSection.Number)
				if err != nil {
					number = 0
				}

				filterRoadSection.Number = number
				filterRoadSection.Name = "ตอนควบคุม " + item.Road.RoadSection.Number + " " + item.Road.RoadSection.NameOriginTH + " - " + item.Road.RoadSection.NameDestinationTH
				roadSectionMaps[strconv.Itoa(item.Year)+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)] = append(roadSectionMaps[strconv.Itoa(item.Year)+"_"+strconv.Itoa(item.Road.RoadSection.RoadGroup.Id)], filterRoadSection)

				roadSectionDuplicateMap[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RoadGroup.Number+"_"+item.Road.RoadSection.Number] = true
			}

			if len(roadSectionMaps) > 0 {
				if !yearIsduplicateMap[item.Year] {
					years = append(years, item.Year)
					yearIsduplicateMap[item.Year] = true
				}
			}

			if !roadGroupIsDuplicateMap[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RoadGroup.Number] {
				var filterRoadGroup responses.FilterRoadGroup
				filterRoadGroup.Id = item.Road.RoadSection.RoadGroup.Id
				filterRoadGroup.Name = "ทางหลวงหมายเลข " + item.Road.RoadSection.RoadGroup.Number + " " + item.Road.RoadSection.RoadGroup.ShortName
				if len(roadSectionMaps) > 0 {
					roadGroupMaps[item.Year] = append(roadGroupMaps[item.Year], filterRoadGroup)
				}
				roadGroupIsDuplicateMap[strconv.Itoa(item.Year)+"_"+item.Road.RoadSection.RoadGroup.Number] = true
			}

		}
	}

	var filterAadts []responses.FilterAadtYear
	for _, item := range years {
		var filterAadt responses.FilterAadtYear
		filterAadt.Year = item

		var filterRoadGroups []responses.FilterRoadGroup
		for _, itemRoadGroup := range roadGroupMaps[item] {
			var filterRoadGroup responses.FilterRoadGroup
			filterRoadGroup = itemRoadGroup

			var filterRoadSections []responses.FilterRoadSection
			for _, itemRoadSection := range roadSectionMaps[strconv.Itoa(item)+"_"+strconv.Itoa(itemRoadGroup.Id)] {
				var filterRoadSection responses.FilterRoadSection
				filterRoadSection = itemRoadSection
				filterRoadSections = append(filterRoadSections, filterRoadSection)
			}

			sort.SliceStable(filterRoadSections, func(i, j int) bool {
				return filterRoadSections[i].Number < filterRoadSections[j].Number
			})

			filterRoadGroup.RoadSection = filterRoadSections
			filterRoadGroups = append(filterRoadGroups, filterRoadGroup)
		}

		filterAadt.RoadGroup = filterRoadGroups
		filterAadts = append(filterAadts, filterAadt)
	}

	var filterAadtType1 responses.FilterAadtType1
	filterAadtType1.FilterRoad = filterAadts

	return filterAadtType1, nil
}

//////////////// NEW MIMS ////////////////

func (u *UseCase) Report9(roadSectionId, filterCriteriaId, year int, typ string) (interface{}, error) {

	retroReflectivity, err := u.Repo.ReportRetroReflectivity(roadSectionId, filterCriteriaId, year)
	if err != nil {
		return nil, err
	}

	paramsRoadLine, err := u.Repo.ParamsRoadLine(filterCriteriaId)
	if err != nil {
		return nil, err
	}

	refStripeColor, err := u.Repo.RefStripeColor()
	if err != nil {
		return nil, err
	}

	colorStripMap := map[int]models.RefStripeColor{}
	for _, item := range refStripeColor {
		colorStripMap[item.ID] = item
	}

	refGrade, err := u.Repo.RefGrade()
	if err != nil {
		return nil, err
	}

	colorRefGrade := map[int]models.RefGrade{}
	for _, item := range refGrade {
		colorRefGrade[item.ID] = item
	}

	var report9 responses.Report9
	report9.RoadGroupNumber = retroReflectivity.RoadGroup.Number
	report9.RoadSectionNumber = retroReflectivity.RoadSection.Number
	report9.RoadSectionName = retroReflectivity.RoadSection.NameOriginTH + " - " + retroReflectivity.RoadSection.NameDestinationTH
	report9.RoadGroupName = retroReflectivity.RoadGroup.ShortName
	report9.KmStart = strings.ReplaceAll(fmt.Sprintf("%.3f", retroReflectivity.RoadSection.KmStart/1000), ".", "+")
	report9.KmEnd = strings.ReplaceAll(fmt.Sprintf("%.3f", retroReflectivity.RoadSection.KmEnd/1000), ".", "+")
	report9.Distance = fmt.Sprintf("%.3f", math.Abs(float64(retroReflectivity.RoadSection.KmStart)-float64(retroReflectivity.RoadSection.KmEnd))/1000)
	report9.Year = strconv.Itoa(year + 543)

	var criterias []responses.Report9Criteria

	for _, itemColor := range refStripeColor {
		var criteria responses.Report9Criteria
		criteria.ColorName = itemColor.NameTH

		var report9CriteriaValues []responses.Report9CriteriaValue
		if itemColor.NameTH == "เส้นสีขาว" {
			for _, item := range paramsRoadLine {
				var report9CriteriaValue responses.Report9CriteriaValue
				report9CriteriaValue.Criteria = item.RefGrade.Name
				report9CriteriaValue.OperationLeft = item.LeftValueWhite
				report9CriteriaValue.OperatorLeft = item.LeftConditionWhite
				report9CriteriaValue.OperationRight = item.RightValueWhite
				report9CriteriaValue.OperatorRight = item.RightConditionWhite

				report9CriteriaValues = append(report9CriteriaValues, report9CriteriaValue)

			}
		} else {
			for _, item := range paramsRoadLine {
				var report9CriteriaValue responses.Report9CriteriaValue
				report9CriteriaValue.Criteria = item.RefGrade.Name
				report9CriteriaValue.OperationLeft = item.LeftValueYellow
				report9CriteriaValue.OperatorLeft = item.LeftConditionYellow
				report9CriteriaValue.OperationRight = item.RightValueYellow
				report9CriteriaValue.OperatorRight = item.RightConditionYellow

				report9CriteriaValues = append(report9CriteriaValues, report9CriteriaValue)

			}
		}

		criteria.CriteriaValue = report9CriteriaValues

		criterias = append(criterias, criteria)
	}

	var tables []responses.Report9Table

	passSummary := 0.0
	allGeom := ""
	passGeom := []string{}
	notPassSummary := 0.0
	notPassGeom := []string{}
	criteriaConditionSummary := 0.0
	g7AvgSummary := 0.0
	for _, itemReflectivity := range retroReflectivity.RoadRetroReflectivity {
		var table responses.Report9Table
		table.Line = itemReflectivity.LineNo
		pass := 0.0
		notPass := 0.0
		g7Avg := 0.0
		g7AvgCount := 0.0
		for _, itemRange := range itemReflectivity.RoadRetroReflectivityRange {
			table.Color = colorStripMap[itemRange.RefStripeColorID].NameTH
			length := math.Abs(itemRange.KmStart - itemRange.KmEnd)
			criteriaConditionSummary += length
			g7Avg += *itemRange.RetroAvg
			g7AvgCount += 1.0
			for _, itemParamsRoadLine := range paramsRoadLine {
				if colorStripMap[itemRange.RefStripeColorID].NameTH == "เส้นสีขาว" {
					if itemParamsRoadLine.LeftValueWhite <= *itemRange.RetroAvg && *itemRange.RetroAvg < itemParamsRoadLine.RightValueWhite {
						pass += length
						passGeom = append(passGeom, itemRange.TheGeom)
					} else {
						notPass += length
						notPassGeom = append(notPassGeom, itemRange.TheGeom)
					}
				} else {
					if itemParamsRoadLine.LeftValueYellow <= *itemRange.RetroAvg && *itemRange.RetroAvg < itemParamsRoadLine.RightValueYellow {
						pass += length
						passGeom = append(passGeom, itemRange.TheGeom)
					} else {
						notPass += length
						notPassGeom = append(notPassGeom, itemRange.TheGeom)
					}
				}
				allGeom = allGeom + strings.ReplaceAll(itemRange.TheGeom, ",", ")) LINESTRING(")
				break
			}

		}
		table.Pass = pass / 1000
		table.NotPass = notPass / 1000
		table.G7Avg = g7Avg / g7AvgCount

		passSummary += table.Pass
		notPassSummary += table.NotPass
		g7AvgSummary += table.G7Avg

		tables = append(tables, table)

	}

	var graph responses.Report9Graph
	graph.Lable = []string{"ผ่าน", "ไม่ผ่าน"}
	graph.Color = []string{"#42d235", "#ff290a"}
	graph.Value = []float64{(passSummary * criteriaConditionSummary) / 100, (notPassSummary * criteriaConditionSummary) / 100}

	var tableSummary responses.Report9TableSummary
	tableSummary.G7Avg = g7AvgSummary
	tableSummary.Pass = passSummary
	tableSummary.NotPass = notPassSummary

	var geoms []responses.Report9Map
	var geom1 responses.Report9Map
	geom1.Color = "#42d235"
	geom1.TheGeom = passGeom
	geoms = append(geoms, geom1)
	var geom2 responses.Report9Map
	geom2.Color = "#ff290a"
	geom2.TheGeom = notPassGeom
	geoms = append(geoms, geom2)

	report9.Criteria = criterias

	report9.Graph = graph

	report9.Table = tables

	report9.TableSummary = tableSummary

	var allmap []string

	allmap = strings.Split(strings.ReplaceAll(allGeom, ")LINESTRING(", ") LINESTRING("), ") ")

	mapCenter := helpers.CalculateMapCenter(allmap)

	report9.MapCenter = mapCenter

	report9.Map = geoms

	if typ == "excel" {
		logo := os.Getenv("LOGO")

		fileLogo, err := os.Open(fmt.Sprint(logo))
		if err != nil {
			fmt.Println("can not open file", err)
		}
		defer fileLogo.Close()

		// imageLogoConfig, _, err := image.DecodeConfig(fileLogo)
		// if err != nil {
		// 	fmt.Println("can not check image size", err)
		// }

		filePath := os.Getenv("GENARAL_EXCEL")
		f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE9_EXCEL"))
		if err != nil {
			return nil, err
		}

		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println(err)
			}
		}()

		pathChart, err := ExportPNGChart(report9, "TEMPLATE_GENARAL_TYPE9_EXCEL_CHART")
		if err != nil {
			return nil, err
		}

		// img size checker
		fileChart, err := os.Open(fmt.Sprint(pathChart))
		if err != nil {
			fmt.Println("can not open file", err)
		}
		defer fileChart.Close()

		// imageChartConfig, _, err := image.DecodeConfig(fileChart)
		// if err != nil {
		// 	fmt.Println("can not check image size", err)
		// }

		defer os.Remove(fmt.Sprint(pathChart))

		pathMap, err := ExportPNGMap(report9, "TEMPLATE_GENARAL_TYPE9_EXCEL_MAP")
		if err != nil {
			return nil, err
		}

		// img size checker
		fileMap, err := os.Open(fmt.Sprint(pathMap))
		if err != nil {
			fmt.Println("can not open file", err)
		}
		defer fileMap.Close()

		imageMapConfig, _, err := image.DecodeConfig(fileMap)
		if err != nil {
			fmt.Println("can not check image size", err)
		}

		//defer os.Remove(fmt.Sprint(pathMap))

		textCenter, err := f.NewStyle(
			&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Font:      &excelize.Font{Bold: true, Color: "00FF0000", Size: 15, Family: "TH SarabunPSK"},
			},
		)
		if err != nil {
			return nil, err
		}

		textCenterWithBoarder, err := f.NewStyle(
			&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Font:      &excelize.Font{Bold: false, Color: "00FF0000", Size: 15, Family: "TH SarabunPSK"},
				Border: []excelize.Border{
					{Type: "left", Color: "00FF0000", Style: 1},
					{Type: "right", Color: "00FF0000", Style: 1},
					{Type: "top", Color: "00FF0000", Style: 1},
					{Type: "bottom", Color: "00FF0000", Style: 1},
				},
			},
		)

		sheetName := "sheet1"
		location := 0
		location = location + 1
		f.AddPicture(sheetName, "A"+fmt.Sprint(location), fmt.Sprint(logo), &excelize.GraphicOptions{
			OffsetX: 45,
			OffsetY: 10,
			ScaleX:  0.05,
			ScaleY:  0.04},
		)
		//f.AddPicture(sheetName, "A"+fmt.Sprint(location), fmt.Sprint(logo), &excelize.GraphicOptions{ScaleX: 160 / float64(imageLogoConfig.Width), ScaleY: 120 / float64(imageLogoConfig.Height)})

		location = location + 1
		f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenter)
		f.MergeCell(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location))
		f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "รายงานสรุปข้อมูลค่าสะท้อนแสงของเส้นจราจร (Retroreflectivity)")

		location = location + 1
		f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenter)
		f.MergeCell(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location))
		f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "หมายเลขทางหลวง : "+fmt.Sprint(report9.RoadGroupNumber)+" ตอนควบคุม : "+fmt.Sprint(report9.RoadSectionNumber)+" ชื่อสายทาง : "+report9.RoadSectionName)

		location = location + 1
		f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenter)
		f.MergeCell(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location))
		f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "กม.เริ่มต้น "+fmt.Sprint(report9.KmStart)+" กม.สิ้นสุด "+fmt.Sprint(report9.KmEnd)+" ระยะทาง "+fmt.Sprint(report9.Distance)+" กม.")

		location = location + 2
		f.SetCellStyle(sheetName, "G"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenter)
		f.MergeCell(sheetName, "G"+fmt.Sprint(location), "J"+fmt.Sprint(location))
		f.SetCellValue(sheetName, "G"+fmt.Sprint(location), "เภณฑ์ที่ใช้ (หน่วย : mcd.x-1.m-2)")

		for _, itemCriteria := range report9.Criteria {
			location = location + 1
			f.SetCellStyle(sheetName, "G"+fmt.Sprint(location), "G"+fmt.Sprint(location), textCenter)
			f.SetCellValue(sheetName, "G"+fmt.Sprint(location), itemCriteria.ColorName)

			for _, itemCriteriaValue := range itemCriteria.CriteriaValue {
				f.SetCellStyle(sheetName, "H"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenter)
				f.SetCellValue(sheetName, "H"+fmt.Sprint(location), itemCriteriaValue.Criteria)
				f.SetCellValue(sheetName, "I"+fmt.Sprint(location), itemCriteriaValue.OperatorLeft)
				f.SetCellValue(sheetName, "J"+fmt.Sprint(location), itemCriteriaValue.OperationLeft)
				location = location + 1

			}
		}

		location = location + 1

		// f.AddPicture(sheetName, "B"+fmt.Sprint(location), fmt.Sprint(pathChart), &excelize.GraphicOptions{ScaleX: 10 / float64(imageChartConfig.Width), ScaleY: 5 / float64(imageChartConfig.Height)})
		f.AddPicture(sheetName, "B"+fmt.Sprint(location), fmt.Sprint(pathChart), &excelize.GraphicOptions{
			OffsetX: 70,
			ScaleX:  1.02,
			ScaleY:  1.08})

		location = location + 23
		f.SetCellStyle(sheetName, "B"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenterWithBoarder)
		f.MergeCell(sheetName, "B"+fmt.Sprint(location), "B"+fmt.Sprint(location+1))
		f.SetCellValue(sheetName, "B"+fmt.Sprint(location), "เส้นจารจร")
		f.MergeCell(sheetName, "C"+fmt.Sprint(location), "C"+fmt.Sprint(location+1))
		f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "สี")
		f.MergeCell(sheetName, "D"+fmt.Sprint(location), "D"+fmt.Sprint(location+1))
		f.SetCellValue(sheetName, "D"+fmt.Sprint(location), "G7 เฉลี่ย (ม./กม.) mcd.x-1.m-2)")
		f.MergeCell(sheetName, "E"+fmt.Sprint(location), "J"+fmt.Sprint(location))
		f.SetCellValue(sheetName, "E"+fmt.Sprint(location), "ระยะทางในแต่ละช่วง (กม.)")
		f.MergeCell(sheetName, "B"+fmt.Sprint(location), "B"+fmt.Sprint(location+1))
		location = location + 1
		f.SetCellStyle(sheetName, "B"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenterWithBoarder)
		f.MergeCell(sheetName, "E"+fmt.Sprint(location), "G"+fmt.Sprint(location))
		f.SetCellValue(sheetName, "E"+fmt.Sprint(location), "ผ่าน")
		f.MergeCell(sheetName, "H"+fmt.Sprint(location), "J"+fmt.Sprint(location))
		f.SetCellValue(sheetName, "H"+fmt.Sprint(location), "ไม่ผ่าน")

		for _, itemTable := range report9.Table {
			location = location + 1

			f.SetCellStyle("Sheet1", "B"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenterWithBoarder)
			f.SetCellValue("Sheet1", "B"+fmt.Sprint(location), itemTable.Line)
			f.SetCellValue("Sheet1", "C"+fmt.Sprint(location), itemTable.Color)
			f.SetCellValue("Sheet1", "D"+fmt.Sprint(location), fmt.Sprintf(`%.2f`, itemTable.G7Avg))
			f.MergeCell(sheetName, "E"+fmt.Sprint(location), "G"+fmt.Sprint(location))
			f.SetCellValue("Sheet1", "E"+fmt.Sprint(location), fmt.Sprintf(`%.2f`, itemTable.Pass))
			f.MergeCell(sheetName, "H"+fmt.Sprint(location), "J"+fmt.Sprint(location))
			f.SetCellValue("Sheet1", "H"+fmt.Sprint(location), fmt.Sprintf(`%.2f`, itemTable.NotPass))
		}
		location = location + 1
		f.SetCellStyle("Sheet1", "B"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenterWithBoarder)
		f.MergeCell(sheetName, "B"+fmt.Sprint(location), "C"+fmt.Sprint(location))
		f.SetCellValue("Sheet1", "B"+fmt.Sprint(location), "รวม")
		f.SetCellValue("Sheet1", "D"+fmt.Sprint(location), fmt.Sprintf(`%.2f`, report9.TableSummary.G7Avg))
		f.MergeCell(sheetName, "E"+fmt.Sprint(location), "G"+fmt.Sprint(location))
		f.SetCellValue("Sheet1", "E"+fmt.Sprint(location), fmt.Sprintf(`%.2f`, report9.TableSummary.Pass))
		f.MergeCell(sheetName, "H"+fmt.Sprint(location), "J"+fmt.Sprint(location))
		f.SetCellValue("Sheet1", "H"+fmt.Sprint(location), fmt.Sprintf(`%.2f`, report9.TableSummary.NotPass))

		if report9.MapCenter != "" {
			location = location + 2
			// f.AddPicture(sheetName, "B"+fmt.Sprint(location), fmt.Sprint(pathMap), &excelize.GraphicOptions{ScaleX: 580 / float64(imageMapConfig.Width), ScaleY: 500 / float64(imageMapConfig.Height)})
			f.AddPicture(sheetName, "B"+fmt.Sprint(location), fmt.Sprint(pathMap), &excelize.GraphicOptions{ScaleX: 580 / float64(imageMapConfig.Width), ScaleY: 500 / float64(imageMapConfig.Height)})
		}

		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			logs.Error(err)
			return nil, err
		}

		helpers.AddFooter(f, sheetName)

		reportName, err := helpers.ReportName("TEMPLATE_GENARAL_TYPE9_EXCEL")
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		code := fmt.Sprintf("%04d", rand.Intn(10000))

		name := code + "_" + reportName

		f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

		return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil

	} else {
		pathResult, err := helpers.RequestExport(report9, "TEMPLATE_GENARAL_TYPE9_HTML", typ)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}

		return pathResult, nil
	}

	return nil, nil
}

func (u *UseCase) Report8(roadSectionId, filterCriteriaId, year int, typ string) (interface{}, error) {

	retroReflectivity, err := u.Repo.ReportRetroReflectivity(roadSectionId, filterCriteriaId, year)
	if err != nil {
		return nil, err
	}

	refStripeColor, err := u.Repo.RefStripeColor()
	if err != nil {
		return nil, err
	}

	colorStripMap := map[int]models.RefStripeColor{}
	for _, item := range refStripeColor {
		colorStripMap[item.ID] = item
	}

	refStripeType, err := u.Repo.RefStripeType()
	if err != nil {
		return nil, err
	}

	typeStripMap := map[int]models.RefStripeType{}
	for _, item := range refStripeType {
		typeStripMap[item.ID] = item
	}

	var report8 responses.Report8

	report8.RoadGroupNumber = retroReflectivity.RoadGroup.Number
	report8.RoadSectionNumber = retroReflectivity.RoadSection.Number
	report8.RoadSectionName = retroReflectivity.RoadSection.NameOriginTH + " - " + retroReflectivity.RoadSection.NameDestinationTH
	report8.RoadGroupName = retroReflectivity.RoadGroup.ShortName
	report8.KmStart = strings.ReplaceAll(fmt.Sprintf("%.3f", retroReflectivity.RoadSection.KmStart/1000), ".", "+")
	report8.KmEnd = strings.ReplaceAll(fmt.Sprintf("%.3f", retroReflectivity.RoadSection.KmEnd/1000), ".", "+")
	report8.Distance = fmt.Sprintf("%.3f", math.Abs(float64(retroReflectivity.RoadSection.KmStart)-float64(retroReflectivity.RoadSection.KmEnd))/1000)

	level2 := map[string][]responses.Report8TableLevel2{}
	for _, item1 := range retroReflectivity.RoadRetroReflectivity {
		for _, item2 := range item1.RoadRetroReflectivityRange {
			var report8TableLevel2 responses.Report8TableLevel2
			report8TableLevel2.KmStart = strings.ReplaceAll(fmt.Sprintf("%.3f", item2.KmStart/1000), ".", "+")
			report8TableLevel2.KmEnd = strings.ReplaceAll(fmt.Sprintf("%.3f", item2.KmEnd/1000), ".", "+")
			report8TableLevel2.RetroReflectivity = fmt.Sprintf("%.3f", *item2.RetroMax)
			report8TableLevel2.Distance = math.Abs(float64(item2.KmStart) - float64(item2.KmEnd))
			level2[strconv.Itoa(item2.RefStripeTypeID)+"_"+strconv.Itoa(item2.RefStripeColorID)] = append(level2[strconv.Itoa(item2.RefStripeTypeID)+"_"+strconv.Itoa(item2.RefStripeColorID)], report8TableLevel2)
		}
	}

	for _, item1 := range retroReflectivity.RoadRetroReflectivity {
		var report8Road responses.Report8Road
		report8Road.Name = item1.RoadInfo.Name
		report8Road.RoadId = item1.RoadID
		level1Isduplicate := map[string]bool{}
		var level1 []responses.Report8TableLevel1
		for _, item2 := range item1.RoadRetroReflectivityRange {
			if level1Isduplicate[strconv.Itoa(item2.RefStripeTypeID)+"_"+strconv.Itoa(item2.RefStripeColorID)] == false {
				var report8TableLevel1 responses.Report8TableLevel1
				report8TableLevel1.Line = item1.LineNo
				report8TableLevel1.Color = colorStripMap[item2.RefStripeColorID].NameTH
				report8TableLevel1.Type = typeStripMap[item2.RefStripeTypeID].NameTH
				report8TableLevel1.Sub = level2[strconv.Itoa(item2.RefStripeTypeID)+"_"+strconv.Itoa(item2.RefStripeColorID)]

				sumDistance := 0.0
				sumRetro := 0.0
				summary := 0.0
				for _, item3 := range level2[strconv.Itoa(item2.RefStripeTypeID)+"_"+strconv.Itoa(item2.RefStripeColorID)] {
					sumDistance += item3.Distance
					retroReflectivity, _ := strconv.ParseFloat(item3.RetroReflectivity, 64)
					sumRetro += retroReflectivity * item3.Distance
				}

				report8TableLevel1.Distance = sumDistance
				summary = sumRetro / sumDistance
				report8TableLevel1.RetroReflectivityAvg = fmt.Sprintf("%.3f", helpers.RoundFloat(summary, 3))

				level1 = append(level1, report8TableLevel1)
				level1Isduplicate[strconv.Itoa(item2.RefStripeTypeID)+"_"+strconv.Itoa(item2.RefStripeColorID)] = true
			}
		}
		report8Road.Table = level1
		report8.Road = append(report8.Road, report8Road)
	}

	if typ == "excel" {

		logo := os.Getenv("LOGO")

		fileLogo, err := os.Open(fmt.Sprint(logo))
		if err != nil {
			fmt.Println("can not open file", err)
		}
		defer fileLogo.Close()

		imageLogoConfig, _, err := image.DecodeConfig(fileLogo)
		if err != nil {
			fmt.Println("can not check image size", err)
		}

		filePath := os.Getenv("GENARAL_EXCEL")
		f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE8_EXCEL"))
		if err != nil {
			return nil, err
		}

		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println(err)
			}
		}()

		textCenter, err := f.NewStyle(
			&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Font:      &excelize.Font{Bold: true, Color: "00FF0000", Size: 15, Family: "TH SarabunPSK"},
			},
		)

		textLeft, err := f.NewStyle(
			&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "left", WrapText: true},
				Font:      &excelize.Font{Bold: false, Color: "00FF0000", Size: 15, Family: "TH SarabunPSK"},
			},
		)
		if err != nil {
			return nil, err
		}

		textCenterWithBoarder, err := f.NewStyle(
			&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Font:      &excelize.Font{Bold: false, Color: "00FF0000", Size: 15, Family: "TH SarabunPSK"},
				Border: []excelize.Border{
					{Type: "left", Color: "00FF0000", Style: 1},
					{Type: "right", Color: "00FF0000", Style: 1},
					{Type: "top", Color: "00FF0000", Style: 1},
					{Type: "bottom", Color: "00FF0000", Style: 1},
				},
			},
		)

		location := 0

		sheetName := ""
		if len(report8.Road) > 0 {
			sheetName = fmt.Sprint(report8.Road[0].RoadId)
		}

		roadName := ""
		if len(report8.Road) > 0 {
			roadName = fmt.Sprint(report8.Road[0].Name)
		}

		index, err := f.NewSheet(sheetName)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		location = location + 1

		f.AddPicture(sheetName, "A"+fmt.Sprint(location), fmt.Sprint(logo), &excelize.GraphicOptions{ScaleX: 160 / float64(imageLogoConfig.Width), ScaleY: 120 / float64(imageLogoConfig.Height)})

		location = location + 1
		f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenter)
		f.MergeCell(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location))
		f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "รายงานค่าสะท้อนแสงของเส้นจราจร (Retroreflectivity)")

		location = location + 1
		f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenter)
		f.MergeCell(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location))
		f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "หมายเลขทางหลวง : "+fmt.Sprint(report8.RoadGroupNumber)+" ตอนควบคุม : "+fmt.Sprint(report8.RoadSectionNumber)+" ชื่อสายทาง : "+report8.RoadSectionName)

		location = location + 1
		f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenter)
		f.MergeCell(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location))
		f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "กม.เริ่มต้น "+fmt.Sprint(report8.KmStart)+" กม.สิ้นสุด "+fmt.Sprint(report8.KmEnd)+" ระยะทาง "+fmt.Sprint(report8.Distance)+" กม.")

		location = location + 3
		f.SetCellStyle(sheetName, "B"+fmt.Sprint(location), "G"+fmt.Sprint(location), textLeft)
		f.MergeCell(sheetName, "B"+fmt.Sprint(location), "G"+fmt.Sprint(location))
		f.SetCellValue(sheetName, "B"+fmt.Sprint(location), fmt.Sprintf(`ชื่อสายทาง: %s`, roadName))
		location = location + 1
		f.SetCellStyle(sheetName, "B"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterWithBoarder)
		f.SetCellValue(sheetName, "B"+fmt.Sprint(location), "เส้นจารจร")
		f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "สีเส้นจารจร")
		f.SetCellValue(sheetName, "D"+fmt.Sprint(location), "ชนิดเส้นจารจร")
		f.SetCellValue(sheetName, "E"+fmt.Sprint(location), "ระยะทาง (ม.)")
		f.SetCellValue(sheetName, "F"+fmt.Sprint(location), "กม.เริ่มต้น")
		f.SetCellValue(sheetName, "G"+fmt.Sprint(location), "กม.เริ่มสิ้นสุด")
		f.MergeCell(sheetName, "H"+fmt.Sprint(location), "I"+fmt.Sprint(location))
		f.SetCellValue(sheetName, "H"+fmt.Sprint(location), "Retroreflectivity\n(mcd.lx-1.m-2)")
		f.MergeCell(sheetName, "J"+fmt.Sprint(location), "K"+fmt.Sprint(location))
		f.SetCellValue(sheetName, "J"+fmt.Sprint(location), "Avg Retroreflectivity\n(mcd.lx-1.m-2)")

		for _, item1 := range report8.Road {
			for _, item2 := range item1.Table {
				location = location + 1
				f.SetCellStyle(sheetName, "B"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterWithBoarder)
				f.SetCellValue(sheetName, "B"+fmt.Sprint(location), item2.Line)
				f.SetCellValue(sheetName, "C"+fmt.Sprint(location), item2.Color)
				f.SetCellValue(sheetName, "D"+fmt.Sprint(location), item2.Type)
				f.SetCellValue(sheetName, "E"+fmt.Sprint(location), item2.Distance)
				f.MergeCell(sheetName, "H"+fmt.Sprint(location), "I"+fmt.Sprint(location))
				f.MergeCell(sheetName, "J"+fmt.Sprint(location), "K"+fmt.Sprint(location))
				f.SetCellValue(sheetName, "J"+fmt.Sprint(location), item2.RetroReflectivityAvg)
				for _, item3 := range item2.Sub {
					location = location + 1
					f.SetCellStyle(sheetName, "B"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterWithBoarder)
					f.SetCellValue(sheetName, "B"+fmt.Sprint(location), "")
					f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "")
					f.SetCellValue(sheetName, "D"+fmt.Sprint(location), "")
					f.SetCellValue(sheetName, "E"+fmt.Sprint(location), "")
					f.SetCellValue(sheetName, "F"+fmt.Sprint(location), item3.KmStart)
					f.SetCellValue(sheetName, "G"+fmt.Sprint(location), item3.KmEnd)
					f.MergeCell(sheetName, "H"+fmt.Sprint(location), "I"+fmt.Sprint(location))
					f.SetCellValue(sheetName, "H"+fmt.Sprint(location), item3.RetroReflectivity)
					f.MergeCell(sheetName, "J"+fmt.Sprint(location), "K"+fmt.Sprint(location))
					f.SetCellValue(sheetName, "J"+fmt.Sprint(location), "")
				}
			}
			helpers.AddFooter(f, sheetName)

			f.SetActiveSheet(index)
		}

		err = f.DeleteSheet("Sheet1")
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			logs.Error(err)
			return nil, err
		}

		reportName, err := helpers.ReportName("TEMPLATE_GENARAL_TYPE8_EXCEL")
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		code := fmt.Sprintf("%04d", rand.Intn(10000))

		name := code + "_" + reportName

		f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

		return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil

	} else {
		pathResult, err := helpers.RequestExport(report8, "TEMPLATE_GENARAL_TYPE8_HTML", typ)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}

		return pathResult, nil
	}

	return report8, nil
}

func (u *UseCase) Report12(roadSectionId, year int, filterConditionName, typ string) (interface{}, error) {

	if filterConditionName == "IRI" {

		var Report12 responses.Report12Iri
		roadGroup, err := u.Repo.Report12RoadGroup(roadSectionId)
		if err != nil {
			return nil, err
		}

		table, err := u.Repo.Report12Iri(roadSectionId, year)
		if err != nil {
			return nil, err
		}

		// calRol1000 := len(table.ResultIri100) % 19
		// calRol100 := len(table.ResultIri100) % 19
		// Report12.LastRowOnPage1000M = calRol1000
		// Report12.LastRowOnPage100M = calRol100

		Report12.RoadGroupNumber = roadGroup.RoadGroup.Number
		Report12.RoadSectionNumber = roadGroup.RoadSection.Number
		Report12.RoadSectionName = roadGroup.RoadSection.NameOriginTH + " - " + roadGroup.RoadSection.NameDestinationTH
		Report12.RoadGroupName = roadGroup.RoadGroup.ShortName
		Report12.KmStart = strings.ReplaceAll(fmt.Sprintf("%.3f", roadGroup.RoadSection.KmStart/1000), ".", "+")
		Report12.KmEnd = strings.ReplaceAll(fmt.Sprintf("%.3f", roadGroup.RoadSection.KmEnd/1000), ".", "+")
		Report12.Distance = fmt.Sprintf("%.3f", math.Abs(float64(roadGroup.RoadSection.KmStart)-float64(roadGroup.RoadSection.KmEnd))/1000)
		Report12.Table = table

		if typ == "excel" {

			logo := os.Getenv("LOGO")

			fileLogo, err := os.Open(fmt.Sprint(logo))
			if err != nil {
				return nil, err
			}
			defer fileLogo.Close()

			imageLogoConfig, _, err := image.DecodeConfig(fileLogo)
			if err != nil {
				return nil, err
			}

			filePath := os.Getenv("GENARAL_EXCEL")
			f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE12_IRI_EXCEL"))
			if err != nil {

				fmt.Println(err.Error())
				return nil, err
			}

			defer func() {
				if err := f.Close(); err != nil {
				}
			}()

			textCenter, err := f.NewStyle(
				&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Font:      &excelize.Font{Bold: true, Color: "00FF0000", Size: 15, Family: "TH SarabunPSK"},
				},
			)
			if err != nil {
				return nil, err
			}

			textCenterWithBoarder, err := f.NewStyle(
				&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Font:      &excelize.Font{Bold: false, Color: "00FF0000", Size: 15, Family: "TH SarabunPSK"},
					Border: []excelize.Border{
						{Type: "left", Color: "00FF0000", Style: 1},
						{Type: "right", Color: "00FF0000", Style: 1},
						{Type: "top", Color: "00FF0000", Style: 1},
						{Type: "bottom", Color: "00FF0000", Style: 1},
					},
				},
			)

			textCenterRedWithBoarder, err := f.NewStyle(
				&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Font:      &excelize.Font{Bold: false, Color: "FF0000", Size: 15, Family: "TH SarabunPSK"},
					Border: []excelize.Border{
						{Type: "left", Color: "00FF0000", Style: 1},
						{Type: "right", Color: "00FF0000", Style: 1},
						{Type: "top", Color: "00FF0000", Style: 1},
						{Type: "bottom", Color: "00FF0000", Style: 1},
					},
				},
			)

			sheetName := "sheet1"
			location := 0
			location = location + 1
			f.AddPicture(sheetName, "A"+fmt.Sprint(location), fmt.Sprint(logo), &excelize.GraphicOptions{ScaleX: 180 / float64(imageLogoConfig.Width), ScaleY: 120 / float64(imageLogoConfig.Height)})

			location = location + 2
			f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenter)
			f.MergeCell(sheetName, "C"+fmt.Sprint(location), "K"+fmt.Sprint(location))
			f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "รายงานการซ่อมบํารุงตามเกณฑ์ KPI ค่าดัชนีความขรุขระสากล (InternationalRoughness Index : IRI)")

			location = location + 1
			f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenter)
			f.MergeCell(sheetName, "C"+fmt.Sprint(location), "K"+fmt.Sprint(location))
			f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "หมายเลขทางหลวง : "+fmt.Sprint(Report12.RoadGroupNumber)+" ตอนควบคุม : "+fmt.Sprint(Report12.RoadSectionNumber)+" ชื่อตอนควบคุม : "+Report12.RoadSectionName)

			location = location + 1
			f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenter)
			f.MergeCell(sheetName, "C"+fmt.Sprint(location), "K"+fmt.Sprint(location))
			f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "กม.เริ่มต้น "+fmt.Sprint(Report12.KmStart)+" กม.สิ้นสุด "+fmt.Sprint(Report12.KmEnd)+" ระยะทาง "+fmt.Sprint(Report12.Distance)+" กม.")

			location = location + 3
			f.SetCellValue(sheetName, "A"+fmt.Sprint(location), "1,000 เมตร")

			location = location + 1
			f.SetCellStyle(sheetName, "A"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterWithBoarder)
			f.SetCellValue(sheetName, "A"+fmt.Sprint(location), "ลําดับ")
			f.SetCellValue(sheetName, "B"+fmt.Sprint(location), "หมายเลขทางหลวง")
			f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "ตอนควบคุม")
			f.SetCellValue(sheetName, "D"+fmt.Sprint(location), "ชื่อถนน")
			f.SetCellValue(sheetName, "E"+fmt.Sprint(location), "กม.เริ่มต้น")
			f.SetCellValue(sheetName, "F"+fmt.Sprint(location), "กม.เริ่มสิ้นสุด")
			f.SetCellValue(sheetName, "G"+fmt.Sprint(location), "ช่องจราจร")
			f.SetCellValue(sheetName, "H"+fmt.Sprint(location), "ผิวทาง")
			f.SetCellValue(sheetName, "I"+fmt.Sprint(location), "IRI (ม./กม.)")
			f.SetCellValue(sheetName, "J"+fmt.Sprint(location), "กำหนดซ่อมภายใน")
			f.SetCellValue(sheetName, "K"+fmt.Sprint(location), "หมายเหตุ")

			if len(Report12.Table.ResultIri1000) == 0 {
				location = location + 1
				f.MergeCell(sheetName, "A"+fmt.Sprint(location), "K"+fmt.Sprint(location))
				f.SetCellStyle(sheetName, "A"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterWithBoarder)
				f.SetCellValue(sheetName, "A"+fmt.Sprint(location), "ไม่พบข้อมูล")
			} else {
				for index, item := range Report12.Table.ResultIri1000 {
					location = location + 1
					f.SetCellStyle(sheetName, "A"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenterWithBoarder)
					f.SetCellValue(sheetName, "A"+fmt.Sprint(location), index+1)
					f.SetCellValue(sheetName, "B"+fmt.Sprint(location), item.RoadCode)
					f.SetCellValue(sheetName, "C"+fmt.Sprint(location), item.SectionCode)
					f.SetCellValue(sheetName, "D"+fmt.Sprint(location), item.RoadName)
					f.SetCellValue(sheetName, "E"+fmt.Sprint(location), item.KmStart)
					f.SetCellValue(sheetName, "F"+fmt.Sprint(location), item.KmEnd)
					f.SetCellValue(sheetName, "G"+fmt.Sprint(location), item.LaneNo)
					f.SetCellValue(sheetName, "H"+fmt.Sprint(location), item.SurveyType)
					f.SetCellValue(sheetName, "I"+fmt.Sprint(location), item.Iri)
					if item.IsExpire {
						f.SetCellStyle(sheetName, "J"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenterRedWithBoarder)
					} else {
						f.SetCellStyle(sheetName, "J"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenterWithBoarder)
					}
					f.SetCellValue(sheetName, "J"+fmt.Sprint(location), item.MaintExpireDate)

					f.SetCellStyle(sheetName, "K"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterWithBoarder)
					comment := ""
					for idexComment, itemComment := range item.Comment {
						comment = comment + fmt.Sprint(idexComment+1) + ") กม. " + itemComment.KmStart + " - " + itemComment.KmEnd + " \nมีงานซ่อมเมื่อวันที่ " + itemComment.ProjectEndDate + "\n"
					}
					f.SetCellValue(sheetName, "K"+fmt.Sprint(location), comment)
				}
			}

			location = location + 3
			f.SetCellValue(sheetName, "A"+fmt.Sprint(location), "100 เมตร")

			location = location + 1
			f.SetCellStyle(sheetName, "A"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterWithBoarder)
			f.SetCellValue(sheetName, "A"+fmt.Sprint(location), "ลําดับ")
			f.SetCellValue(sheetName, "B"+fmt.Sprint(location), "หมายเลขทางหลวง")
			f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "ตอนควบคุม")
			f.SetCellValue(sheetName, "D"+fmt.Sprint(location), "ชื่อถนน")
			f.SetCellValue(sheetName, "E"+fmt.Sprint(location), "กม.เริ่มต้น")
			f.SetCellValue(sheetName, "F"+fmt.Sprint(location), "กม.เริ่มสิ้นสุด")
			f.SetCellValue(sheetName, "G"+fmt.Sprint(location), "ช่องจราจร")
			f.SetCellValue(sheetName, "H"+fmt.Sprint(location), "ผิวทาง")
			f.SetCellValue(sheetName, "I"+fmt.Sprint(location), "IRI (ม./กม.)")
			f.SetCellValue(sheetName, "J"+fmt.Sprint(location), "กำหนดซ่อมภายใน")
			f.SetCellValue(sheetName, "K"+fmt.Sprint(location), "หมายเหตุ")

			if len(Report12.Table.ResultIri100) == 0 {
				location = location + 1
				f.MergeCell(sheetName, "A"+fmt.Sprint(location), "K"+fmt.Sprint(location))
				f.SetCellStyle(sheetName, "A"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterWithBoarder)
				f.SetCellValue(sheetName, "A"+fmt.Sprint(location), "ไม่พบข้อมูล")
			} else {
				for index, item := range Report12.Table.ResultIri100 {
					location = location + 1
					f.SetCellStyle(sheetName, "A"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterWithBoarder)
					f.SetCellValue(sheetName, "A"+fmt.Sprint(location), index+1)
					f.SetCellValue(sheetName, "B"+fmt.Sprint(location), item.RoadCode)
					f.SetCellValue(sheetName, "C"+fmt.Sprint(location), item.SectionCode)
					f.SetCellValue(sheetName, "D"+fmt.Sprint(location), item.RoadName)
					f.SetCellValue(sheetName, "E"+fmt.Sprint(location), item.KmStart)
					f.SetCellValue(sheetName, "F"+fmt.Sprint(location), item.KmEnd)
					f.SetCellValue(sheetName, "G"+fmt.Sprint(location), item.LaneNo)
					f.SetCellValue(sheetName, "H"+fmt.Sprint(location), item.SurveyType)
					f.SetCellValue(sheetName, "I"+fmt.Sprint(location), item.Iri)
					if item.IsExpire {
						f.SetCellStyle(sheetName, "J"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenterRedWithBoarder)
					} else {
						f.SetCellStyle(sheetName, "J"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenterWithBoarder)
					}
					f.SetCellValue(sheetName, "J"+fmt.Sprint(location), item.MaintExpireDate)

					f.SetCellStyle(sheetName, "K"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterWithBoarder)
					comment := ""
					for idexComment, itemComment := range item.Comment {
						comment = comment + fmt.Sprint(idexComment+1) + ") กม. " + itemComment.KmStart + " - " + itemComment.KmEnd + " \nมีงานซ่อมเมื่อวันที่ " + itemComment.ProjectEndDate + "\n"
					}
					f.SetCellValue(sheetName, "K"+fmt.Sprint(location), comment)
				}
			}

			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				logs.Error(err)
				return nil, err
			}

			reportName, err := helpers.ReportName("TEMPLATE_GENARAL_TYPE12_IRI_EXCEL")
			if err != nil {
				logs.Error(err)
				return nil, err
			}
			code := fmt.Sprintf("%04d", rand.Intn(10000))

			name := code + "_" + reportName

			f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

			return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil

		} else {
			pathResult, err := helpers.RequestExport(Report12, "TEMPLATE_GENARAL_TYPE12_IRI_HTML", typ)
			if err != nil {
				return nil, responses.NewAppErr(400, err.Error())
			}

			return pathResult, nil
		}

	} else if filterConditionName == "IFI" {

		var Report12 responses.Report12Ifi
		roadGroup, err := u.Repo.Report12RoadGroup(roadSectionId)
		if err != nil {
			return nil, err
		}

		table, err := u.Repo.Report12Ifi(roadSectionId, year)
		if err != nil {
			return nil, err
		}

		Report12.RoadGroupNumber = roadGroup.RoadGroup.Number
		Report12.RoadSectionNumber = roadGroup.RoadSection.Number
		Report12.RoadSectionName = roadGroup.RoadSection.NameOriginTH + " - " + roadGroup.RoadSection.NameDestinationTH
		Report12.RoadGroupName = roadGroup.RoadGroup.ShortName
		Report12.KmStart = strings.ReplaceAll(fmt.Sprintf("%.3f", roadGroup.RoadSection.KmStart/1000), ".", "+")
		Report12.KmEnd = strings.ReplaceAll(fmt.Sprintf("%.3f", roadGroup.RoadSection.KmEnd/1000), ".", "+")
		Report12.Distance = fmt.Sprintf("%.3f", math.Abs(float64(roadGroup.RoadSection.KmStart)-float64(roadGroup.RoadSection.KmEnd))/1000)
		Report12.Table = table

		if typ == "excel" {

			logo := os.Getenv("LOGO")

			fileLogo, err := os.Open(fmt.Sprint(logo))
			if err != nil {
				return nil, err
			}
			defer fileLogo.Close()

			imageLogoConfig, _, err := image.DecodeConfig(fileLogo)
			if err != nil {
				return nil, err
			}

			filePath := os.Getenv("GENARAL_EXCEL")
			f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE12_IFI_EXCEL"))
			if err != nil {

				fmt.Println(err.Error())
				return nil, err
			}

			defer func() {
				if err := f.Close(); err != nil {
				}
			}()

			textCenter, err := f.NewStyle(
				&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Font:      &excelize.Font{Bold: true, Color: "00FF0000", Size: 15, Family: "TH SarabunPSK"},
				},
			)
			if err != nil {
				return nil, err
			}

			textCenterWithBoarder, err := f.NewStyle(
				&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Font:      &excelize.Font{Bold: false, Color: "00FF0000", Size: 15, Family: "TH SarabunPSK"},
					Border: []excelize.Border{
						{Type: "left", Color: "00FF0000", Style: 1},
						{Type: "right", Color: "00FF0000", Style: 1},
						{Type: "top", Color: "00FF0000", Style: 1},
						{Type: "bottom", Color: "00FF0000", Style: 1},
					},
				},
			)

			textCenterRedWithBoarder, err := f.NewStyle(
				&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Font:      &excelize.Font{Bold: false, Color: "FF0000", Size: 15, Family: "TH SarabunPSK"},
					Border: []excelize.Border{
						{Type: "left", Color: "00FF0000", Style: 1},
						{Type: "right", Color: "00FF0000", Style: 1},
						{Type: "top", Color: "00FF0000", Style: 1},
						{Type: "bottom", Color: "00FF0000", Style: 1},
					},
				},
			)

			sheetName := "sheet1"
			location := 0
			location = location + 1
			f.AddPicture(sheetName, "A"+fmt.Sprint(location), fmt.Sprint(logo), &excelize.GraphicOptions{ScaleX: 160 / float64(imageLogoConfig.Width), ScaleY: 120 / float64(imageLogoConfig.Height)})

			location = location + 1
			f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenter)
			f.MergeCell(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location))
			f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "รายงานการซ่อมบํารุงตามเกณฑ์ KPI ค่าดัชนีความขรุขระสากล (InternationalFriction Index : IFI)")

			location = location + 1
			f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenter)
			f.MergeCell(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location))
			f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "หมายเลขทางหลวง : "+fmt.Sprint(Report12.RoadGroupNumber)+" ตอนควบคุม : "+fmt.Sprint(Report12.RoadSectionNumber)+" ชื่อตอนควบคุม : "+Report12.RoadSectionName)

			location = location + 1
			f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenter)
			f.MergeCell(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location))
			f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "กม.เริ่มต้น "+fmt.Sprint(Report12.KmStart)+" กม.สิ้นสุด "+fmt.Sprint(Report12.KmEnd)+" ระยะทาง "+fmt.Sprint(Report12.Distance)+" กม.")

			location = location + 3
			f.SetCellValue(sheetName, "B"+fmt.Sprint(location), "100 เมตร")

			location = location + 1
			f.SetCellStyle(sheetName, "B"+fmt.Sprint(location), "L"+fmt.Sprint(location), textCenterWithBoarder)
			f.SetCellValue(sheetName, "B"+fmt.Sprint(location), "ลําดับ")
			f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "หมายเลขทางหลวง")
			f.SetCellValue(sheetName, "D"+fmt.Sprint(location), "ตอนควบคุม")
			f.SetCellValue(sheetName, "E"+fmt.Sprint(location), "ชื่อถนน")
			f.SetCellValue(sheetName, "F"+fmt.Sprint(location), "กม.เริ่มต้น")
			f.SetCellValue(sheetName, "G"+fmt.Sprint(location), "กม.เริ่มสิ้นสุด")
			f.SetCellValue(sheetName, "H"+fmt.Sprint(location), "ช่องจราจร")
			f.SetCellValue(sheetName, "I"+fmt.Sprint(location), "ผิวทาง")
			f.SetCellValue(sheetName, "J"+fmt.Sprint(location), "IFI")
			f.SetCellValue(sheetName, "K"+fmt.Sprint(location), "กำหนดซ่อมภายใน")
			f.SetCellValue(sheetName, "L"+fmt.Sprint(location), "หมายเหตุ")

			if len(Report12.Table) == 0 {
				location = location + 1
				f.MergeCell(sheetName, "B"+fmt.Sprint(location), "L"+fmt.Sprint(location))
				f.SetCellStyle(sheetName, "B"+fmt.Sprint(location), "L"+fmt.Sprint(location), textCenterWithBoarder)
				f.SetCellValue(sheetName, "B"+fmt.Sprint(location), "ไม่พบข้อมูล")
			} else {
				for index, item := range Report12.Table {
					location = location + 1
					f.SetCellStyle(sheetName, "B"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenterWithBoarder)
					f.SetCellValue(sheetName, "B"+fmt.Sprint(location), index+1)
					f.SetCellValue(sheetName, "C"+fmt.Sprint(location), item.RoadCode)
					f.SetCellValue(sheetName, "D"+fmt.Sprint(location), item.SectionCode)
					f.SetCellValue(sheetName, "E"+fmt.Sprint(location), item.RoadName)
					f.SetCellValue(sheetName, "F"+fmt.Sprint(location), item.KmStart)
					f.SetCellValue(sheetName, "G"+fmt.Sprint(location), item.KmEnd)
					f.SetCellValue(sheetName, "H"+fmt.Sprint(location), item.LaneNo)
					f.SetCellValue(sheetName, "I"+fmt.Sprint(location), item.SurveyType)
					f.SetCellValue(sheetName, "J"+fmt.Sprint(location), item.Ifi)
					if item.IsExpire {
						f.SetCellStyle(sheetName, "K"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterRedWithBoarder)
					} else {
						f.SetCellStyle(sheetName, "K"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterWithBoarder)
					}
					f.SetCellValue(sheetName, "K"+fmt.Sprint(location), item.MaintExpireDate)

					f.SetCellStyle(sheetName, "L"+fmt.Sprint(location), "L"+fmt.Sprint(location), textCenterWithBoarder)
					comment := ""
					for idexComment, itemComment := range item.Comment {
						comment = comment + fmt.Sprint(idexComment+1) + ") กม. " + itemComment.KmStart + " - " + itemComment.KmEnd + " \nมีงานซ่อมเมื่อวันที่ " + itemComment.ProjectEndDate + "\n"
					}
					f.SetCellValue(sheetName, "L"+fmt.Sprint(location), comment)
				}
			}

			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				logs.Error(err)
				return nil, err
			}

			reportName, err := helpers.ReportName("TEMPLATE_GENARAL_TYPE12_IFI_EXCEL")
			if err != nil {
				logs.Error(err)
				return nil, err
			}
			code := fmt.Sprintf("%04d", rand.Intn(10000))

			name := code + "_" + reportName

			f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

			return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil

		} else {
			pathResult, err := helpers.RequestExport(Report12, "TEMPLATE_GENARAL_TYPE12_IFI_HTML", typ)
			if err != nil {
				return nil, responses.NewAppErr(400, err.Error())
			}

			return pathResult, nil
		}

	} else if filterConditionName == "RUT" {

		var Report12 responses.Report12Rut
		roadGroup, err := u.Repo.Report12RoadGroup(roadSectionId)
		if err != nil {
			return nil, err
		}

		table, err := u.Repo.Report12Rut(roadSectionId, year)
		if err != nil {
			return nil, err
		}

		Report12.RoadGroupNumber = roadGroup.RoadGroup.Number
		Report12.RoadSectionNumber = roadGroup.RoadSection.Number
		Report12.RoadSectionName = roadGroup.RoadSection.NameOriginTH + " - " + roadGroup.RoadSection.NameDestinationTH
		Report12.RoadGroupName = roadGroup.RoadGroup.ShortName
		Report12.KmStart = strings.ReplaceAll(fmt.Sprintf("%.3f", roadGroup.RoadSection.KmStart/1000), ".", "+")
		Report12.KmEnd = strings.ReplaceAll(fmt.Sprintf("%.3f", roadGroup.RoadSection.KmEnd/1000), ".", "+")
		Report12.Distance = fmt.Sprintf("%.3f", math.Abs(float64(roadGroup.RoadSection.KmStart)-float64(roadGroup.RoadSection.KmEnd))/1000)
		Report12.Table = table

		if typ == "excel" {

			logo := os.Getenv("LOGO")

			fileLogo, err := os.Open(fmt.Sprint(logo))
			if err != nil {
				return nil, err
			}
			defer fileLogo.Close()

			imageLogoConfig, _, err := image.DecodeConfig(fileLogo)
			if err != nil {
				return nil, err
			}

			filePath := os.Getenv("GENARAL_EXCEL")
			f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE12_RUT_EXCEL"))
			if err != nil {

				fmt.Println(err.Error())
				return nil, err
			}

			defer func() {
				if err := f.Close(); err != nil {
				}
			}()

			textCenter, err := f.NewStyle(
				&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Font:      &excelize.Font{Bold: true, Color: "00FF0000", Size: 15, Family: "TH SarabunPSK"},
				},
			)
			if err != nil {
				return nil, err
			}

			textCenterWithBoarder, err := f.NewStyle(
				&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Font:      &excelize.Font{Bold: false, Color: "00FF0000", Size: 15, Family: "TH SarabunPSK"},
					Border: []excelize.Border{
						{Type: "left", Color: "00FF0000", Style: 1},
						{Type: "right", Color: "00FF0000", Style: 1},
						{Type: "top", Color: "00FF0000", Style: 1},
						{Type: "bottom", Color: "00FF0000", Style: 1},
					},
				},
			)

			textCenterRedWithBoarder, err := f.NewStyle(
				&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Font:      &excelize.Font{Bold: false, Color: "FF0000", Size: 15, Family: "TH SarabunPSK"},
					Border: []excelize.Border{
						{Type: "left", Color: "00FF0000", Style: 1},
						{Type: "right", Color: "00FF0000", Style: 1},
						{Type: "top", Color: "00FF0000", Style: 1},
						{Type: "bottom", Color: "00FF0000", Style: 1},
					},
				},
			)

			sheetName := "sheet1"
			location := 0
			location = location + 1
			f.AddPicture(sheetName, "A"+fmt.Sprint(location), fmt.Sprint(logo), &excelize.GraphicOptions{ScaleX: 160 / float64(imageLogoConfig.Width), ScaleY: 120 / float64(imageLogoConfig.Height)})

			location = location + 1
			f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenter)
			f.MergeCell(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location))
			f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "รายงานการซ่อมบํารุงตามเกณฑ์ KPI ค่าดัชนีความขรุขระสากล (RuthDepth : RUT)")

			location = location + 1
			f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenter)
			f.MergeCell(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location))
			f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "หมายเลขทางหลวง : "+fmt.Sprint(Report12.RoadGroupNumber)+" ตอนควบคุม : "+fmt.Sprint(Report12.RoadSectionNumber)+" ชื่อตอนควบคุม : "+Report12.RoadSectionName)

			location = location + 1
			f.SetCellStyle(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenter)
			f.MergeCell(sheetName, "C"+fmt.Sprint(location), "J"+fmt.Sprint(location))
			f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "กม.เริ่มต้น "+fmt.Sprint(Report12.KmStart)+" กม.สิ้นสุด "+fmt.Sprint(Report12.KmEnd)+" ระยะทาง "+fmt.Sprint(Report12.Distance)+" กม.")

			location = location + 3
			f.SetCellValue(sheetName, "B"+fmt.Sprint(location), "100 เมตร")

			location = location + 1
			f.SetCellStyle(sheetName, "B"+fmt.Sprint(location), "L"+fmt.Sprint(location), textCenterWithBoarder)
			f.SetCellValue(sheetName, "B"+fmt.Sprint(location), "ลําดับ")
			f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "หมายเลขทางหลวง")
			f.SetCellValue(sheetName, "D"+fmt.Sprint(location), "ตอนควบคุม")
			f.SetCellValue(sheetName, "E"+fmt.Sprint(location), "ชื่อถนน")
			f.SetCellValue(sheetName, "F"+fmt.Sprint(location), "กม.เริ่มต้น")
			f.SetCellValue(sheetName, "G"+fmt.Sprint(location), "กม.เริ่มสิ้นสุด")
			f.SetCellValue(sheetName, "H"+fmt.Sprint(location), "ช่องจราจร")
			f.SetCellValue(sheetName, "I"+fmt.Sprint(location), "ผิวทาง")
			f.SetCellValue(sheetName, "J"+fmt.Sprint(location), "RUT (มม.)")
			f.SetCellValue(sheetName, "K"+fmt.Sprint(location), "กำหนดซ่อมภายใน")
			f.SetCellValue(sheetName, "L"+fmt.Sprint(location), "หมายเหตุ")

			if len(Report12.Table) == 0 {
				location = location + 1
				f.MergeCell(sheetName, "B"+fmt.Sprint(location), "L"+fmt.Sprint(location))
				f.SetCellStyle(sheetName, "B"+fmt.Sprint(location), "L"+fmt.Sprint(location), textCenterWithBoarder)
				f.SetCellValue(sheetName, "B"+fmt.Sprint(location), "ไม่พบข้อมูล")
			} else {
				for index, item := range Report12.Table {
					location = location + 1
					f.SetCellStyle(sheetName, "B"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenterWithBoarder)
					f.SetCellValue(sheetName, "B"+fmt.Sprint(location), index+1)
					f.SetCellValue(sheetName, "C"+fmt.Sprint(location), item.RoadCode)
					f.SetCellValue(sheetName, "D"+fmt.Sprint(location), item.SectionCode)
					f.SetCellValue(sheetName, "E"+fmt.Sprint(location), item.RoadName)
					f.SetCellValue(sheetName, "F"+fmt.Sprint(location), item.KmStart)
					f.SetCellValue(sheetName, "G"+fmt.Sprint(location), item.KmEnd)
					f.SetCellValue(sheetName, "H"+fmt.Sprint(location), item.LaneNo)
					f.SetCellValue(sheetName, "I"+fmt.Sprint(location), item.SurveyType)
					f.SetCellValue(sheetName, "J"+fmt.Sprint(location), item.Rut)
					if item.IsExpire {
						f.SetCellStyle(sheetName, "K"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterRedWithBoarder)
					} else {
						f.SetCellStyle(sheetName, "K"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterWithBoarder)
					}
					f.SetCellValue(sheetName, "K"+fmt.Sprint(location), item.MaintExpireDate)

					f.SetCellStyle(sheetName, "L"+fmt.Sprint(location), "L"+fmt.Sprint(location), textCenterWithBoarder)
					comment := ""
					for idexComment, itemComment := range item.Comment {
						comment = comment + fmt.Sprint(idexComment+1) + ") กม. " + itemComment.KmStart + " - " + itemComment.KmEnd + " \nมีงานซ่อมเมื่อวันที่ " + itemComment.ProjectEndDate + "\n"
					}
					f.SetCellValue(sheetName, "L"+fmt.Sprint(location), comment)
				}
			}

			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				logs.Error(err)
				return nil, err
			}

			reportName, err := helpers.ReportName("TEMPLATE_GENARAL_TYPE12_RUT_EXCEL")
			if err != nil {
				logs.Error(err)
				return nil, err
			}
			code := fmt.Sprintf("%04d", rand.Intn(10000))

			name := code + "_" + reportName

			f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

			return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil

		} else {
			pathResult, err := helpers.RequestExport(Report12, "TEMPLATE_GENARAL_TYPE12_RUT_HTML", typ)
			if err != nil {
				return nil, responses.NewAppErr(400, err.Error())
			}

			return pathResult, nil
		}
	} else {

		var Report12 responses.Report12G7
		roadGroup, err := u.Repo.Report12RoadGroup(roadSectionId)
		if err != nil {
			return nil, err
		}

		table, err := u.Repo.Report12G7(roadSectionId, year)
		if err != nil {
			return nil, err
		}

		Report12.RoadGroupNumber = roadGroup.RoadGroup.Number
		Report12.RoadSectionNumber = roadGroup.RoadSection.Number
		Report12.RoadSectionName = roadGroup.RoadSection.NameOriginTH + " - " + roadGroup.RoadSection.NameDestinationTH
		Report12.RoadGroupName = roadGroup.RoadGroup.ShortName
		Report12.KmStart = strings.ReplaceAll(fmt.Sprintf("%.3f", roadGroup.RoadSection.KmStart/1000), ".", "+")
		Report12.KmEnd = strings.ReplaceAll(fmt.Sprintf("%.3f", roadGroup.RoadSection.KmEnd/1000), ".", "+")
		Report12.Distance = fmt.Sprintf("%.3f", math.Abs(float64(roadGroup.RoadSection.KmStart)-float64(roadGroup.RoadSection.KmEnd))/1000)
		Report12.Table = table

		if typ == "excel" {

			logo := os.Getenv("LOGO")

			fileLogo, err := os.Open(fmt.Sprint(logo))
			if err != nil {
				return nil, err
			}
			defer fileLogo.Close()

			imageLogoConfig, _, err := image.DecodeConfig(fileLogo)
			if err != nil {
				return nil, err
			}

			filePath := os.Getenv("GENARAL_EXCEL")
			f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE12_G7_EXCEL"))
			if err != nil {

				fmt.Println(err.Error())
				return nil, err
			}

			defer func() {
				if err := f.Close(); err != nil {
				}
			}()

			textCenter, err := f.NewStyle(
				&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Font:      &excelize.Font{Bold: true, Color: "00FF0000", Size: 15, Family: "TH SarabunPSK"},
				},
			)
			if err != nil {
				return nil, err
			}

			textCenterWithBoarder, err := f.NewStyle(
				&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Font:      &excelize.Font{Bold: false, Color: "00FF0000", Size: 12, Family: "TH SarabunPSK"},
					Border: []excelize.Border{
						{Type: "left", Color: "00FF0000", Style: 1},
						{Type: "right", Color: "00FF0000", Style: 1},
						{Type: "top", Color: "00FF0000", Style: 1},
						{Type: "bottom", Color: "00FF0000", Style: 1},
					},
				},
			)

			textCenterRedWithBoarder, err := f.NewStyle(
				&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Font:      &excelize.Font{Bold: false, Color: "FF0000", Size: 12, Family: "TH SarabunPSK"},
					Border: []excelize.Border{
						{Type: "left", Color: "00FF0000", Style: 1},
						{Type: "right", Color: "00FF0000", Style: 1},
						{Type: "top", Color: "00FF0000", Style: 1},
						{Type: "bottom", Color: "00FF0000", Style: 1},
					},
				},
			)

			sheetName := "sheet1"
			location := 0
			location = location + 1
			f.AddPicture(sheetName, "A"+fmt.Sprint(location), fmt.Sprint(logo), &excelize.GraphicOptions{ScaleX: 160 / float64(imageLogoConfig.Width), ScaleY: 120 / float64(imageLogoConfig.Height)})

			location = location + 1
			f.SetCellStyle(sheetName, "D"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenter)
			f.MergeCell(sheetName, "D"+fmt.Sprint(location), "K"+fmt.Sprint(location))
			f.SetCellValue(sheetName, "D"+fmt.Sprint(location), "รายงานการซ่อมบํารุงตามเกณฑ์ KPI ค่าดัชนีความขรุขระสากล (Retroreflectivity)")

			location = location + 1
			f.SetCellStyle(sheetName, "D"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenter)
			f.MergeCell(sheetName, "D"+fmt.Sprint(location), "K"+fmt.Sprint(location))
			f.SetCellValue(sheetName, "D"+fmt.Sprint(location), "หมายเลขทางหลวง : "+fmt.Sprint(Report12.RoadGroupNumber)+" ตอนควบคุม : "+fmt.Sprint(Report12.RoadSectionNumber)+" ชื่อตอนควบคุม : "+Report12.RoadSectionName)

			location = location + 1
			f.SetCellStyle(sheetName, "D"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenter)
			f.MergeCell(sheetName, "D"+fmt.Sprint(location), "K"+fmt.Sprint(location))
			f.SetCellValue(sheetName, "D"+fmt.Sprint(location), "กม.เริ่มต้น "+fmt.Sprint(Report12.KmStart)+" กม.สิ้นสุด "+fmt.Sprint(Report12.KmEnd)+" ระยะทาง "+fmt.Sprint(Report12.Distance)+" กม.")

			location = location + 4
			f.SetCellValue(sheetName, "A"+fmt.Sprint(location), "100 เมตร")

			location = location + 1
			f.SetCellStyle(sheetName, "A"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterWithBoarder)
			f.SetCellValue(sheetName, "A"+fmt.Sprint(location), "ลําดับ")
			f.SetCellValue(sheetName, "B"+fmt.Sprint(location), "หมายเลขทางหลวง")
			f.SetCellValue(sheetName, "C"+fmt.Sprint(location), "ตอนควบคุม")
			f.SetCellValue(sheetName, "D"+fmt.Sprint(location), "ชื่อถนน")
			f.SetCellValue(sheetName, "E"+fmt.Sprint(location), "กม.เริ่มต้น")
			f.SetCellValue(sheetName, "F"+fmt.Sprint(location), "กม.เริ่มสิ้นสุด")
			f.SetCellValue(sheetName, "G"+fmt.Sprint(location), "เส้นจราจร")
			f.SetCellValue(sheetName, "H"+fmt.Sprint(location), "สีจราจร")
			f.SetCellValue(sheetName, "I"+fmt.Sprint(location), "G7 (mcd.lx-1.m-2)")
			f.SetCellValue(sheetName, "J"+fmt.Sprint(location), "กำหนดซ่อมภายใน")
			f.SetCellValue(sheetName, "K"+fmt.Sprint(location), "หมายเหตุ")

			if len(Report12.Table) == 0 {
				location = location + 1
				f.MergeCell(sheetName, "A"+fmt.Sprint(location), "L"+fmt.Sprint(location))
				f.SetCellStyle(sheetName, "A"+fmt.Sprint(location), "L"+fmt.Sprint(location), textCenterWithBoarder)
				f.SetCellValue(sheetName, "A"+fmt.Sprint(location), "ไม่พบข้อมูล")
			} else {
				for index, item := range Report12.Table {
					location = location + 1
					f.SetCellStyle(sheetName, "A"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenterWithBoarder)
					f.SetCellValue(sheetName, "A"+fmt.Sprint(location), index+1)
					f.SetCellValue(sheetName, "B"+fmt.Sprint(location), item.RoadCode)
					f.SetCellValue(sheetName, "C"+fmt.Sprint(location), item.SectionCode)
					f.SetCellValue(sheetName, "D"+fmt.Sprint(location), item.RoadName)
					f.SetCellValue(sheetName, "E"+fmt.Sprint(location), item.KmStart)
					f.SetCellValue(sheetName, "F"+fmt.Sprint(location), item.KmEnd)
					f.SetCellValue(sheetName, "G"+fmt.Sprint(location), item.LineNo)
					f.SetCellValue(sheetName, "H"+fmt.Sprint(location), item.NameStrip)
					f.SetCellValue(sheetName, "I"+fmt.Sprint(location), item.Retro)
					if item.IsExpire {
						f.SetCellStyle(sheetName, "J"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenterRedWithBoarder)
					} else {
						f.SetCellStyle(sheetName, "J"+fmt.Sprint(location), "J"+fmt.Sprint(location), textCenterWithBoarder)
					}
					f.SetCellValue(sheetName, "J"+fmt.Sprint(location), item.MaintExpireDate)

					f.SetCellStyle(sheetName, "K"+fmt.Sprint(location), "K"+fmt.Sprint(location), textCenterWithBoarder)
					comment := ""
					for idexComment, itemComment := range item.Comment {
						comment = comment + fmt.Sprint(idexComment+1) + ") กม. " + itemComment.KmStart + " - " + itemComment.KmEnd + " \nมีงานซ่อมเมื่อวันที่ " + itemComment.ProjectEndDate + "\n"
					}
					f.SetCellValue(sheetName, "K"+fmt.Sprint(location), comment)
				}
			}

			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				logs.Error(err)
				return nil, err
			}

			reportName, err := helpers.ReportName("TEMPLATE_GENARAL_TYPE12_G7_EXCEL")
			if err != nil {
				logs.Error(err)
				return nil, err
			}
			code := fmt.Sprintf("%04d", rand.Intn(10000))

			name := code + "_" + reportName

			f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

			return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil

		} else {
			pathResult, err := helpers.RequestExport(Report12, "TEMPLATE_GENARAL_TYPE12_G7_HTML", typ)
			if err != nil {
				return nil, responses.NewAppErr(400, err.Error())
			}

			return pathResult, nil
		}
	}

	return nil, nil
}

func ExportPNGChart(data interface{}, HTMLTemplate string) (interface{}, error) {
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
	if err := chromedp.Run(ctx, PrintToPNGChart(string(html), &buf, true)); err != nil {
		logs.Error(err)
		return nil, err
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

func PrintToPNGChart(html string, res *[]byte, isDelay bool) chromedp.Tasks {
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

			delay := 15

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
			width, height := 430, 400
			if err := emulation.SetDeviceMetricsOverride(int64(width), int64(height), 0, false).Do(ctx); err != nil {
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

func ExportPNGMap(data interface{}, HTMLTemplate string) (interface{}, error) {
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

	// e := os.Remove(result)
	// if e != nil {
	// 	logs.Error(err)
	// 	return nil, err
	// }

	ctx, cancel := helpers.NewChromedpContext(context.Background(), log.Printf)
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, PrintToPNGMap(string(html), &buf, true)); err != nil {
		logs.Error(err)
		return nil, err
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

func PrintToPNGMap(html string, res *[]byte, isDelay bool) chromedp.Tasks {
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

			delay := 20

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
			width, height := 600, 400
			if err := emulation.SetDeviceMetricsOverride(int64(width), int64(height), 1, false).Do(ctx); err != nil {
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
