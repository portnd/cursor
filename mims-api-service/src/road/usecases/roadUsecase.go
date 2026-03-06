package usecases

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	shp "github.com/jonas-p/go-shp"
	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	requests "gitlab.com/mims-api-service/requests"
	responses "gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/road/domains"
	"gorm.io/gorm"

	helperFile "gitlab.com/mims-api-service/helpers/file"
	servicesDB "gitlab.com/mims-api-service/services/database"
)

type roadUseCase struct {
	roadRepo   domains.RoadRepository
	helperFile helperFile.HelpersFileDomain
	servicesDB servicesDB.ServicesDatabaseDomain
}

// init usecase
func NewRoadUseCase(repo domains.RoadRepository, helpFile helperFile.HelpersFileDomain, servicesDB servicesDB.ServicesDatabaseDomain) domains.RoadUseCase {
	return &roadUseCase{
		roadRepo:   repo,
		helperFile: helpFile,
		servicesDB: servicesDB,
	}
}

// =========================================================

func (t *roadUseCase) GetMenu(userId uint) ([]models.AccessControl, error) {
	roles, err := t.roadRepo.GetRole(userId)
	if err != nil {
		return nil, err
	}

	var roleIds []int
	for _, item := range roles {
		roleIds = append(roleIds, item.RoleID)
	}

	data, err := t.roadRepo.GetAccessControl(roleIds)
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

func (t *roadUseCase) GetRoadDetailMenu(userId uint, accessKeys []string, assetType string) (interface{}, error) {
	var assetData []Groups
	// user, err := t.roadRepo.GetUserById(userId)
	// if err != nil {
	// 	return assetData, err
	// }
	// departmentId := user.DepartmentId
	data, err := t.roadRepo.GetRoadDetailMenu(accessKeys, assetType)
	if err != nil {
		return assetData, err
	}
	// helpers.PrintlnJson(accessKeys, departmentId)
	// return data, nil
	switch assetType {
	case "assetin":
		// if helpers.HasPermission([]string{"road_in_asset_manage_data", "road_in_asset_access", "view_all_road_in_asset", "view_department_road_in_asset"}, accessKeys) {
		assetData = t.AssetIn(data)
		if err != nil {
			return assetData, err
		}
		// }
	// case "assetout":
	// 	if helpers.HasPermission([]string{"road_out_asset_manage_data", "road_out_asset_access"}, accessKeys) {
	// 		assetData = t.AssetOut(data)
	// 		if err != nil {
	// 			return assetData, err
	// 		}
	// 	}
	default:
		return []string{}, nil
	}
	if len(assetData) > 0 {
		return assetData, nil
	} else {
		return []string{}, nil
	}

}

func (t *roadUseCase) AssetIn(data []responses.RoadMenuData) []Groups {
	var grops []Groups
	helpers.PrintlnJson(data)
	assetGroupItem := make(map[int]int)
	for _, item := range data {
		if item.IsInRoad {
			var dataResponse2 Groups
			if assetGroupItem[item.GroupId] != item.GroupId {
				fmt.Println(item.GroupName)
				assetGroupItem[item.GroupId] = item.GroupId
				dataResponse2.Id = item.GroupId
				dataResponse2.Name = item.GroupName
				grops = append(grops, dataResponse2)
			}
		}
	}
	var assetinData []Groups
	for _, item := range grops {
		var group Groups
		var groupItems []GroupItems
		for _, item2 := range data {
			if !item2.IsActive {
				continue
			}
			if item.Id == item2.GroupId && item2.IsInRoad {
				var groupItem GroupItems
				groupItem.Id = item2.AssetId
				groupItem.Name = item2.AssetName
				groupItem.GeomType = item2.GeomType
				if item2.IconFilepath == "" {
					groupItem.IconFilepath = ""
				} else {
					groupItem.IconFilepath = os.Getenv("STORAGE_IP") + "/" + item2.IconFilepath //item2.IconFilepath
				}

				groupItems = append(groupItems, groupItem)
			}
		}
		group.Id = item.Id
		group.Name = item.Name
		group.Items = groupItems
		assetinData = append(assetinData, group)
	}
	return assetinData
}

func (t *roadUseCase) AssetOut(data []responses.RoadMenuData) []Groups {
	var grops []Groups
	assetGroupItem := make(map[int]int)
	for _, item := range data {
		if !item.IsInRoad {
			var groups Groups
			if assetGroupItem[item.GroupId] != item.GroupId {
				fmt.Println(item.GroupName)
				assetGroupItem[item.GroupId] = item.GroupId
				groups.Id = item.GroupId
				groups.Name = item.GroupName
				grops = append(grops, groups)
			}
		}
	}
	var assetinData []Groups
	for _, item := range grops {
		var groups Groups
		var groupItems []GroupItems
		for _, item2 := range data {
			if !item2.IsActive {
				continue
			}
			if item.Id == item2.GroupId && !item2.IsInRoad {
				var groupItem GroupItems
				groupItem.Id = item2.AssetId
				groupItem.Name = item2.AssetName
				groupItem.GeomType = item2.GeomType
				if item2.IconFilepath == "" {
					groupItem.IconFilepath = ""
				} else {
					groupItem.IconFilepath = os.Getenv("STORAGE_IP") + "/" + item2.IconFilepath //item2.IconFilepath
				}
				groupItems = append(groupItems, groupItem)
			}
		}
		groups.Id = item.Id
		groups.Name = item.Name
		groups.Items = groupItems
		assetinData = append(assetinData, groups)
	}
	return assetinData
}

func (t *roadUseCase) GetRoadGroupList(userID int, params requests.RoadPrams) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := t.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_road"}, []string{"view_owner_road"})
	//============ end check permission ============

	roadIDs := []int{}
	hasConditionFilter := params.IsIri1000 != nil || params.IsIri100 != nil || params.IsRut100 != nil || params.IsIfi100 != nil || params.IsG7100 != nil

	if hasConditionFilter {
		rfp, err := t.GetRfp(params)
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
		}
		rfpData := rfp.(responses.DataRes)

		if params.IsIri1000 != nil {
			index := 0
			if *params.IsIri1000 {
				index = 1
			}
			roadData := rfpData.Iri1000[index]
			roadIDs = append(roadIDs, roadData...)
		}

		if params.IsIri100 != nil {
			index := 0
			if *params.IsIri100 {
				index = 1
			}
			roadData := rfpData.Iri100[index]
			roadIDs = append(roadIDs, roadData...)
		}

		if params.IsRut100 != nil {
			index := 0
			if *params.IsRut100 {
				index = 1
			}
			roadData := rfpData.Rut100[index]
			roadIDs = append(roadIDs, roadData...)
		}

		if params.IsIfi100 != nil {
			index := 0
			if *params.IsIfi100 {
				index = 1
			}
			roadData := rfpData.Ifi100[index]
			roadIDs = append(roadIDs, roadData...)
		}

		if params.IsG7100 != nil {
			index := 0
			if *params.IsG7100 {
				index = 1
			}
			roadData := rfpData.G7100[index]
			roadIDs = append(roadIDs, roadData...)
		}
	}

	if len(params.RoadId) != 0 {
		roadIDs = append(roadIDs, params.RoadId...)
	}

	data, err := t.roadRepo.GetRoadGroupList(params, roadIDs, isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, responses.NewAppErr(404, err.Error())
	}

	if len(params.RefSurfaceId) != 0 {
		for g, group := range data {
			var sectionData []models.RoadSectionData
			for _, section := range group.Sections {
				var roadData []models.RoadData
				for _, road := range section.Roads {
					childData := make([]models.ChildRoadData, 0)
					for _, child := range road.ChildRoads {
						if helpers.HasDuplicate(child.RefSurface.RefSurfaceId, params.RefSurfaceId) {
							childData = append(childData, child)
						}
					}
					if len(childData) == 0 {
						road.ChildRoads = childData
						if helpers.HasDuplicate(road.RefSurface.RefSurfaceId, params.RefSurfaceId) {
							roadData = append(roadData, road)
						}
					} else {
						roadData = append(roadData, road)
					}
				}
				if len(roadData) != 0 {
					section.Roads = roadData
					sectionData = append(sectionData, section)
				}
			}
			data[g].Sections = sectionData
		}
	}

	for i, res := range data {
		for s, section := range res.Sections {
			for r, road := range section.Roads {
				// check ref road type
				roadOrigin := helpers.CheckOriginDestinationRoad(road.RoadInfo.RefRoadTypeID, section.NameOriginTH, section.NameDestinationTH, "")
				groupNumber := strings.TrimLeft(res.Number, "0")
				data[i].Number = groupNumber
				roadCode := res.Number + section.Number
				responsibleCode := section.RefDistrict.Name + " - " + section.RefDepot.Name
				kmRangeRoad := fmt.Sprintf("%v - %v", road.RoadInfo.KmStart, road.RoadInfo.KmEnd)
				if section.NameDestinationTH != "" {
					data[i].Sections[s].Roads[r].RoadInfo.Name = section.NameOriginTH + " - " + section.NameDestinationTH
				} else {
					data[i].Sections[s].Roads[r].RoadInfo.Name = section.NameOriginTH
				}
				data[i].Sections[s].Roads[r].RoadInfo.OriginToDestination = roadOrigin
				data[i].Sections[s].Roads[r].RoadInfo.RoadCode = roadCode
				data[i].Sections[s].Roads[r].RoadInfo.ResponsibleCode = responsibleCode
				data[i].Sections[s].Roads[r].RoadInfo.KmRange = kmRangeRoad

				for c, child := range road.ChildRoads {
					kmRangeChildRoad := fmt.Sprintf("%v - %v", child.RoadInfo.KmStart, child.RoadInfo.KmEnd)
					data[i].Sections[s].Roads[r].ChildRoads[c].RoadInfo.RoadCode = roadCode
					data[i].Sections[s].Roads[r].ChildRoads[c].RoadInfo.ResponsibleCode = responsibleCode
					data[i].Sections[s].Roads[r].ChildRoads[c].RoadInfo.KmRange = kmRangeChildRoad
				}
			}
		}
	}
	// division แผนก
	// district เขต
	// depot คลังสินค้า
	roadsResponse := make([]models.RoadList, 0)
	if params.Keyword != "" {
		for indexGroup, group := range data {
			statusFindChildRoad := false
			statusFindRoad := false
			sectionWithKeyword := make([]models.RoadSectionData, 0)
			for indexSection, section := range group.Sections {
				roadsWithKeyword := make([]models.RoadData, 0)
				for indexRoad, road := range section.Roads {

					childRoadsWithKeyword := make([]models.ChildRoadData, 0)
					for _, child := range road.ChildRoads {
						// Match by name or road code (รหัสสายทาง e.g. 00070101)
						if strings.Contains(child.RoadInfo.Name, params.Keyword) || strings.Contains(child.RoadInfo.RoadCode, params.Keyword) {
							childRoadsWithKeyword = append(childRoadsWithKeyword, child)
						}
					}
					// Match by name or road code (รหัสสายทาง e.g. 00070101)
					if strings.Contains(road.RoadInfo.Name, params.Keyword) || strings.Contains(road.RoadInfo.RoadCode, params.Keyword) { // if find keyword in level road get all data in child
						roadsWithKeyword = append(roadsWithKeyword, road)
					} else if len(childRoadsWithKeyword) != 0 {
						statusFindChildRoad = true
						data[indexGroup].Sections[indexSection].Roads[indexRoad].ChildRoads = childRoadsWithKeyword
						roadsWithKeyword = append(roadsWithKeyword, data[indexGroup].Sections[indexSection].Roads[indexRoad])
					}

				}
				var sectionFullName string
				if section.NameDestinationTH != "" {
					sectionFullName = section.NameOriginTH + " - " + section.NameDestinationTH
				} else {
					sectionFullName = section.NameOriginTH
				}

				// Section match: name or road code (roads in section share same RoadCode = group.Number+section.Number)
				sectionMatchesKeyword := strings.Contains(sectionFullName, params.Keyword) || (len(section.Roads) > 0 && strings.Contains(section.Roads[0].RoadInfo.RoadCode, params.Keyword))
				if sectionMatchesKeyword {
					sectionWithKeyword = append(sectionWithKeyword, section)
				} else if len(roadsWithKeyword) != 0 {
					statusFindRoad = true
					data[indexGroup].Sections[indexSection].Roads = roadsWithKeyword
					sectionWithKeyword = append(sectionWithKeyword, data[indexGroup].Sections[indexSection])
				}
			}

			if len(sectionWithKeyword) != 0 {
				data[indexGroup].Sections = sectionWithKeyword
			} else if !statusFindRoad && !statusFindChildRoad {
				data[indexGroup].Sections = sectionWithKeyword
			}

			if strings.Contains(group.ShortName, params.Keyword) {
				roadsResponse = append(roadsResponse, group)
			} else if len(data[indexGroup].Sections) != 0 {
				roadsResponse = append(roadsResponse, data[indexGroup])
			}

		}
	} else {
		for _, group := range data {
			if len(group.Sections) != 0 && group.Sections != nil {
				var sectionData []models.RoadSectionData
				for _, section := range group.Sections {
					if len(section.Roads) != 0 && section.Roads != nil {
						sectionData = append(sectionData, section)
					}
				}
				if len(sectionData) != 0 {
					group.Sections = sectionData
					roadsResponse = append(roadsResponse, group)
				}
			}
		}
	}

	for g, group := range roadsResponse {
		for s, section := range group.Sections {
			for r, road := range section.Roads {
				if len(road.MaintenanceRoad) != 0 && len(road.RoadCondition) == 0 {
					roadsResponse[g].Sections[s].Roads[r].SurveyStatus = true
				} else {
					if len(road.MaintenanceRoad) != 0 {
						for _, mainten := range road.MaintenanceRoad {
							checkLane := false
							for _, condition := range road.RoadCondition {
								if mainten.LaneNo == condition.LaneNo {
									checkLane = true
									if mainten.ProjectEndDate.After(condition.SurveyedDate) {
										roadsResponse[g].Sections[s].Roads[r].SurveyStatus = true
										break
									}
								}
							}
							if !checkLane {
								roadsResponse[g].Sections[s].Roads[r].SurveyStatus = true
								break
							}
						}
					} else {
						roadsResponse[g].Sections[s].Roads[r].SurveyStatus = false
					}
				}
				for c, child := range road.ChildRoads {
					if len(child.MaintenanceRoad) != 0 && len(child.RoadCondition) == 0 {
						roadsResponse[g].Sections[s].Roads[r].ChildRoads[c].SurveyStatus = true
					} else {
						if len(road.MaintenanceRoad) != 0 {
							for _, mainten := range road.MaintenanceRoad {
								checkLane := false
								for _, condition := range road.RoadCondition {
									if mainten.LaneNo == condition.LaneNo {
										checkLane = true
										if mainten.ProjectEndDate.After(condition.SurveyedDate) {
											roadsResponse[g].Sections[s].Roads[r].ChildRoads[c].SurveyStatus = true
											break
										}
									}
								}
								if !checkLane {
									roadsResponse[g].Sections[s].Roads[r].ChildRoads[c].SurveyStatus = true
									break
								}
							}
						} else {
							roadsResponse[g].Sections[s].Roads[r].ChildRoads[c].SurveyStatus = false
						}
					}
				}
			}
		}
	}

	return roadsResponse, nil

	// sectionMap := map[string][]models.RoadSectionData{}
	// for _, group := range data {
	// 	for _, section := range group.Sections {
	// 		sectionMap[strconv.Itoa(group.Id)+"_"+section.RefDepot.DepotCode] = append(sectionMap[strconv.Itoa(group.Id)+"_"+section.RefDepot.DepotCode], section)
	// 	}
	// }

	// var roadListNew []models.RoadListNew
	// refDepotDuplicate := map[string]bool{}
	// for _, group := range data {
	// 	var data_temp models.RoadListNew
	// 	data_temp.Id = group.Id
	// 	data_temp.Number = group.Number
	// 	data_temp.Name = group.Name
	// 	data_temp.ShortName = group.ShortName
	// 	data_temp.KmStart = group.KmStart
	// 	data_temp.KmEnd = group.KmEnd
	// 	data_temp.Distance = group.Distance

	// 	for _, section := range group.Sections {
	// 		if !refDepotDuplicate[strconv.Itoa(data_temp.Id)+"_"+section.RefDepot.DepotCode] {
	// 			var refDepot models.RefDepotNew
	// 			refDepot.Id = section.RefDepot.Id
	// 			refDepot.Name = section.RefDepot.Name
	// 			refDepot.DepotCode = section.RefDepot.DepotCode
	// 			refDepot.TheGeom = section.RefDepot.TheGeom

	// 			sectionTemp := sectionMap[strconv.Itoa(data_temp.Id)+"_"+section.RefDepot.DepotCode]

	// 			sort.SliceStable(sectionTemp, func(i, j int) bool {
	// 				i1, err := strconv.Atoi(sectionTemp[i].Number)
	// 				if err != nil {
	// 					i1 = 0
	// 				}

	// 				j1, err := strconv.Atoi(sectionTemp[j].Number)
	// 				if err != nil {
	// 					j1 = 0
	// 				}
	// 				return i1 <= j1
	// 			})

	// 			refDepot.Section = sectionTemp

	// 			data_temp.RefDepot = append(data_temp.RefDepot, refDepot)

	// 			refDepotDuplicate[strconv.Itoa(data_temp.Id)+"_"+section.RefDepot.DepotCode] = true
	// 		}

	// 	}
	// 	roadListNew = append(roadListNew, data_temp)
	// }

	// return roadListNew, nil
}

func (t *roadUseCase) GetRoadByID(roadID, userID int) (*responses.RoadById, error) {
	//============ start check permission ============
	userInfo, _ := t.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_road"}, []string{"view_owner_road"})
	//============ end check permission ============

	result, err := t.roadRepo.GetRoadByID(roadID)
	if err != nil {
		return nil, responses.NewAppErr(404, err.Error())
	}

	if isOwnerData {
		if result.RoadSection.RefDepotCode != depotCode {
			return nil, responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION)
		}
	}
	if !isAllData && !isOwnerData {
		return nil, responses.NewAppErr(http.StatusForbidden, constants.INVALID_USER_PERMISSION)
	}

	var road responses.RoadById
	road.RefDepot = result.RoadSection.RefDepot
	var province string
	for _, data := range result.RoadSection.Province {
		province += data + ","
	}

	if len(province) > 0 {
		province = province[:len(province)-1]
	}

	var roadOrigin string
	if result.RoadLevel == 1 {
		roadOrigin = helpers.CheckOriginDestinationRoad(result.RoadInfo.RefRoadTypeID, result.RoadSection.NameOriginTH, result.RoadSection.NameDestinationTH, "")
	}

	distance := math.Abs(float64(result.RoadInfo.KmStart)-float64(result.RoadInfo.KmEnd)) / 1000

	copier.Copy(&road, &result)

	road.RoadCode = result.RoadSection.RoadGroup.Number + result.RoadSection.Number
	if result.RoadSection.NameDestinationTH != "" {
		road.RoadSectionNameTH = result.RoadSection.NameOriginTH + " - " + result.RoadSection.NameDestinationTH
		road.RoadSectionNameEN = result.RoadSection.NameOriginEn + " - " + result.RoadSection.NameDestinationEn
	} else {
		road.RoadSectionNameTH = result.RoadSection.NameOriginTH
		road.RoadSectionNameEN = result.RoadSection.NameOriginEn
	}

	road.Province = province
	road.ResponsibleCode = result.RoadSection.RefDistrict.Name + " - " + result.RoadSection.RefDepot.Name
	road.OriginToDestination = roadOrigin
	road.KmRange = fmt.Sprintf("%v - %v", result.RoadInfo.KmStart, result.RoadInfo.KmEnd)
	road.Distance = distance

	centerLineDir := os.Getenv("STORAGE_IP")
	if road.RoadInfo.CenterLaneShapeFilePath != "" && road.RoadInfo.CenterLineShapeFilePath != "" {
		road.RoadInfo.CenterLaneShapeFilePath = fmt.Sprintf("%s/%s", centerLineDir, road.RoadInfo.CenterLaneShapeFilePath)
		road.RoadInfo.CenterLineShapeFilePath = fmt.Sprintf("%s/%s", centerLineDir, road.RoadInfo.CenterLineShapeFilePath)
	}
	return &road, nil
}

func (t *roadUseCase) GetRoadGroup() ([]models.RoadGroup, error) {
	// return []models.RoadGroup, nil
	data, err := t.roadRepo.GetRoadGroup()
	if err != nil {
		return data, err
	}
	return data, nil
}

func (t *roadUseCase) GetRoadTypeIcon() (map[string]int, error) {
	icons := make(map[string]int)
	data, err := t.roadRepo.GetRoadTypeIcon()
	if err != nil {
		return icons, err
	}

	for _, item := range data {
		roadTypeId := strconv.Itoa(item.RoadTypeID)
		icons[roadTypeId] = int(item.ID)
	}
	return icons, nil
}

func (t *roadUseCase) GetRoadDirectionLaneList(roadID int) (interface{}, error) {
	road, err := t.roadRepo.GetRoadDirectionLaneList(roadID)
	if err != nil {
		return []responses.RoadLaneList{}, err
	}

	laneMap := make(map[int]bool)
	var roadLaneLists []responses.RoadLaneList
	for _, item := range road.RoadGeom {
		var roadLaneList responses.RoadLaneList
		roadLaneList.LaneNo = item.LaneNo
		roadLaneList.LaneName = fmt.Sprintf("%s%d", road.RoadInfo.Direction.Name, item.LaneNo)
		if _, exists := laneMap[item.LaneNo]; exists {
			continue
		}
		laneMap[item.LaneNo] = true
		roadLaneLists = append(roadLaneLists, roadLaneList)

	}

	if len(roadLaneLists) == 0 {
		return []responses.RoadLaneList{}, nil
	}
	return roadLaneLists, nil
}

// type ItemAsset struct {
// 	ID    int     `json:"id"`
// 	Items []items `json:"items"`
// }

// type items struct {
// 	ID            int `json:"id"`
// 	CountTemp     int `json:"count_temp"`
// 	CountWaiting  int `json:"count_waiting"`
// 	CountRejected int `json:"count_rejected"`
// }

// func (t *roadUseCase) GetRoadDetailStatus(roadID int, permissions []string) {
// 	// edit_surface_data', 'approve_surface_data', 'view_surface_data'
// 	// permissions :=
// 	if helpers.HasPermission(permissions, []string{"road_summary_manage", "approve_road_surface", "road_summary_access"}) {
// 		// road_surface
// 		t.GetRoadStatusSurface(roadID)
// 	}

// 	// 'edit_asset_in_data', 'approve_asset_in_data', 'view_asset_in_data'
// 	if helpers.HasPermission(permissions, []string{"road_in_asset_manage_data", "approve_road_in_asset_access", "road_in_asset_access"}) {
// 		// road_asset in
// 	}

// 	// 'edit_asset_out_data', 'approve_asset_out_data', 'view_asset_out_data'
// 	if helpers.HasPermission(permissions, []string{"road_out_asset_manage_data", "approve_road_out_asset", "road_out_asset_access"}) {
// 		// road_asset out
// 	}

// 	// 'edit_condition_data', 'approve_condition_data', 'view_condition_data'
// 	if helpers.HasPermission(permissions, []string{"road_condition_manage_data", "approve_road_condition", "road_condition_access"}) {
// 		// condition
// 	}

// 	// ['edit_damage_data', 'approve_damage_data', 'view_damage_data']
// 	if helpers.HasPermission(permissions, []string{"road_damage_manage_data", "approve_road_damage", "road_damage_access"}) {
// 		// damage
// 	}

// }

// func (t *roadUseCase) GetRoadStatusSurface(roadID int) map[int]ItemAsset {
// 	data, _ := t.roadRepo.GetRoadStatusSurface(roadID)
// 	helpers.PrintlnJson(data)
// 	countByObject := make(map[int]ItemAsset)

// 	return countByObject
// }

func (t *roadUseCase) GetRoadTree() (interface{}, error) {
	data, err := t.roadRepo.GetRoadGroup()
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return data, responses.NewAppErr(400, err.Error())
	}
	// return data, nil
	var roads []responses.RoadTree
	for _, item := range data {
		roadData, err := t.roadRepo.GetRoadByRoadGrpID(item.Id)
		if len(roadData) == 0 {
			continue
		}
		if err != nil {
			logs.Error(err)
			if err == gorm.ErrRecordNotFound {
				return responses.NoData{}, nil
			}
			return data, responses.NewAppErr(400, err.Error())
		}
		var road responses.RoadTree
		road.ID = item.Id
		road.Label = item.Name
		var childrens []responses.RoadChildren
		for _, item := range roadData {
			var children responses.RoadChildren
			children.ID = item.RoadId
			children.Label = item.Name
			childrens = append(childrens, children)
		}
		road.Children = childrens
		roads = append(roads, road)
	}

	return roads, nil
}

func (t *roadUseCase) CntVolumeApproved(roadGrpID int) (int, int, error) {
	aadt, err := t.roadRepo.GetVolumeAadtApproved(roadGrpID)
	if err != nil {
		return 0, 0, responses.NewAppErr(400, err.Error())
	}

	helpers.PrintlnJson(aadt)

	// for _, item := range aadt {
	// 	aadtData[item.RoadGroupID] = item
	// }

	accident, err := t.roadRepo.GetVolumeAccidentApproved(roadGrpID)
	if err != nil {
		return 0, 0, responses.NewAppErr(400, err.Error())
	}

	countRejected := aadt.CountRejected + accident.CountRejected
	countWaiting := aadt.CountWaiting + accident.CountWaiting

	return countRejected, countWaiting, nil
}

func (t *roadUseCase) GetRoadSectionByID(RoadSectionId int) (*models.RoadSectionById, error) {
	roadSection, err := t.roadRepo.GetRoadSectionByID(RoadSectionId)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}

	return roadSection, nil
}

func (t *roadUseCase) GetDataById(model interface{}, id int) error {
	err := t.roadRepo.GetDataById(&model, id)
	if err != nil {
		return err
	}
	return nil
}

func (t *roadUseCase) GetRoadInit(id, level, refRoadTypeId int) (*responses.RoadInit, error) {
	var sectionId int
	var roadRes responses.RoadInit
	if level == 1 {
		sectionId = id
	} else {
		var road models.Road
		err := t.roadRepo.GetDataById(&road, id)
		if err != nil {
			return nil, err
		}
		sectionId = road.RoadSectionId
	}

	result, err := t.roadRepo.GetRoadSectionByID(sectionId)
	if err != nil {
		return nil, err
	}

	var province string
	for _, data := range result.Province {
		province += data + ","
	}

	if len(province) > 0 {
		province = province[:len(province)-1]
	}

	roadRes.RoadCode = result.RoadGroup.Number + result.Number
	if result.NameDestinationTH != "" {
		roadRes.RoadSectionNameTH = result.NameOriginTH + " - " + result.NameDestinationTH
		roadRes.RoadSectionNameEN = result.NameOriginEn + " - " + result.NameDestinationEn
	} else {
		roadRes.RoadSectionNameTH = result.NameOriginTH
		roadRes.RoadSectionNameEN = result.NameOriginEn
	}

	roadRes.Province = province
	roadRes.District = result.RefDistrict.Name
	roadRes.Depot = result.RefDepot.Name

	if refRoadTypeId == 0 || result.NameDestinationTH == "" {
		roadRes.Origin = result.NameOriginTH
		roadRes.Destination = result.NameDestinationTH
	} else {
		roadOrigin := strings.Split(helpers.CheckOriginDestinationRoad(refRoadTypeId, result.NameOriginTH, result.NameDestinationTH, ""), " - ")
		if len(roadOrigin) == 2 {
			roadRes.Origin = roadOrigin[0]
			roadRes.Destination = roadOrigin[1]
		}

	}
	return &roadRes, nil
}

func (t *roadUseCase) CreateRoad(c *gin.Context, userID int, req requests.Road) (interface{}, error) {

	centerLineDir := os.Getenv("ROAD_CENTER_LINE_SHAPE_FILE_DIR")
	centerLaneDir := os.Getenv("ROAD_CENTER_LANE_SHAPE_FILE_DIR")

	if centerLineDir == "" || centerLaneDir == "" {
		return nil, responses.NewAppErr(400, constants.FAILED_TO_READ_ENV_FILE)
	}

	var registerDate time.Time
	var year int
	if req.RegisterDate != "" {
		reDate, err := time.Parse("2006-01-02 15:04:05", req.RegisterDate)
		if err != nil {
			return 0, responses.NewAppErr(400, err.Error())
		}
		registerDate = reDate
		year = registerDate.Year()
	}

	err := t.helperFile.EnsureDir(centerLineDir)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewAppErr(500, constants.FAILED_TO_CREATE_DIR)
	}

	err = t.helperFile.EnsureDir(centerLaneDir)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewAppErr(500, constants.FAILED_TO_CREATE_DIR)
	}

	centerLinePath, err := t.helperFile.SaveFile(c, req.CenterLineShapeFile, centerLineDir)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewAppErr(500, constants.FAILED_TO_SAVE_FILE)
	}

	//Save File Center Lane
	centerLanePath, err := t.helperFile.SaveFile(c, req.CenterLaneShapeFile, centerLaneDir)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewAppErr(500, constants.FAILED_TO_SAVE_FILE)

	}

	// Convert Center Line Shape File
	lineString, kmStartStr, kmEndStr, err := t.helperFile.ReadCenterLineShepesFile(centerLinePath)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	kmStart, _ := helpers.KmStringToFloat64(*kmStartStr)
	kmEnd, _ := helpers.KmStringToFloat64(*kmEndStr)

	if kmStart != req.KmStart || kmEnd != req.KmEnd {
		logs.Error(err)
		return 0, responses.NewAppErr(400, constants.INVALID_ROAD_RANGE)
	}

	roadGeoms, err := t.helperFile.ReadCenterLaneShepesFile(centerLanePath)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	roadSectionId := req.RoadSectionID
	roadId := req.RoadID
	var refDirectionId int
	var road models.Road
	road.IsActive = true
	if req.RoadLevel == 1 {
		var roadSection models.RoadSection
		err = t.roadRepo.GetDataById(&roadSection, roadSectionId)
		if err != nil {
			logs.Error(err)
			return 0, responses.NewAppErr(400, err.Error())
		}
		road.ParentRoadId = nil
		if req.RefRoadTypeID == 1 || req.RefRoadTypeID == 3 {
			refDirectionId = 1
		} else if req.RefRoadTypeID == 2 || req.RefRoadTypeID == 4 {
			refDirectionId = 2
		}
		road.RoadGroupId = roadSection.RoadGroupId
	} else {
		var roadData models.Road
		err = t.roadRepo.GetDataById(&roadData, roadId)
		if err != nil {
			logs.Error(err)
			return 0, responses.NewAppErr(400, err.Error())
		}
		roadSectionId = roadData.RoadSectionId
		var roadSection models.RoadSection
		err = t.roadRepo.GetDataById(&roadSection, roadSectionId)
		if err != nil {
			logs.Error(err)
			return 0, responses.NewAppErr(400, err.Error())
		}
		//get ref direction for road level 2 from Parent
		roadInfo, err := t.roadRepo.GetLastRoadInfoByID(roadId)
		if err != nil {
			logs.Error(err)
			return 0, responses.NewAppErr(400, err.Error())
		}
		road.RoadGroupId = roadSection.RoadGroupId
		refDirectionId = roadInfo.RefDirectionId
		road.ParentRoadId = &roadId
	}

	road.RoadSectionId = roadSectionId
	road.RoadLevel = req.RoadLevel
	road.CreatedAt = time.Now()
	road.CreatedBy = userID
	road.IsInit = false

	roadMaxSeq, err := t.roadRepo.GetRoadMaxSeq()
	if err != nil {
		logs.Error(err)
		return nil, responses.NewAppErr(500, err.Error())
	}
	road.Seq = roadMaxSeq + 1

	tx := t.roadRepo.StartTransSection()
	if tx.Error != nil {
		logs.Error(err)
		return nil, responses.NewAppErr(500, "Failed to begin transaction")
	}

	err = t.roadRepo.CreateData(tx, &road)
	if err != nil {
		logs.Error(err)
		t.roadRepo.RollBack(tx)
		return nil, responses.NewAppErr(500, constants.FAILED_TO_CREATE_ROAD)
	}

	var roadInfo models.RoadInfo
	copier.Copy(&roadInfo, &req)
	roadInfo.RoadId = road.Id
	roadInfo.Year = &year
	roadInfo.TheGeom = *lineString
	roadInfo.Status = "A"
	roadInfo.CenterLineShapeFilePath = centerLinePath
	roadInfo.CenterLaneShapeFilePath = centerLanePath
	roadInfo.CreatedAt = time.Now()
	roadInfo.CreatedBy = userID
	roadInfo.UpdatedAt = time.Now()
	roadInfo.UpdatedBy = userID
	roadInfo.RefDirectionId = refDirectionId
	roadInfo.YearConstructionCompleted = req.YearConstructionCompleted
	err = t.roadRepo.CreateRoadInfo(tx, &roadInfo)
	if err != nil {
		logs.Error(err)
		t.roadRepo.RollBack(tx)
		return nil, responses.NewAppErr(500, constants.FAILED_TO_CREATE_ROAD)
	}

	for i := 0; i < len(roadGeoms); i++ {
		roadGeoms[i].RoadId = road.Id
		roadGeoms[i].Revision = 0
		roadGeoms[i].CreatedBy = userID
		roadGeoms[i].CreatedAt = time.Now()
		roadGeoms[i].UpdatedBy = userID
		roadGeoms[i].UpdatedAt = time.Now()
		roadGeoms[i].Status = "A"
	}
	err = t.roadRepo.CreateRoadGeom(tx, roadGeoms)
	if err != nil {
		logs.Error(err)
		t.roadRepo.RollBack(tx)
		return nil, responses.NewAppErr(500, constants.FAILED_TO_CREATE_ROAD)
	}
	t.roadRepo.Commit(tx)
	return road.Id, nil
}

func (t *roadUseCase) UpdateRoad(c *gin.Context, roadID int, userID int, req requests.RoadUpdate) (interface{}, error) {
	tx := t.roadRepo.StartTransSection()
	if tx.Error != nil {
		return nil, responses.NewAppErr(500, "Failed to begin transaction")
	}

	centerLineDir := os.Getenv("ROAD_CENTER_LINE_SHAPE_FILE_DIR")
	centerLaneDir := os.Getenv("ROAD_CENTER_LANE_SHAPE_FILE_DIR")

	if centerLineDir == "" || centerLaneDir == "" {
		return nil, responses.NewAppErr(400, constants.FAILED_TO_READ_ENV_FILE)
	}

	var registerDate time.Time
	var year int
	if req.RegisterDate != "" {
		reDate, err := time.Parse("2006-01-02 15:04:05", req.RegisterDate)
		if err != nil {
			return 0, responses.NewAppErr(400, err.Error())
		}
		registerDate = reDate
		year = registerDate.Year()
	}

	road, err := t.roadRepo.GetLastRoadInfoByID(roadID)
	if err != nil {
		return nil, responses.NewNotFoundError()
	}

	var checkDirection int
	if req.RefRoadTypeID == 1 || req.RefRoadTypeID == 3 {
		checkDirection = 1
	} else if req.RefRoadTypeID == 2 || req.RefRoadTypeID == 4 {
		checkDirection = 2
	}

	if road.RoadInfo.RefDirectionId != checkDirection {
		countParent, err := t.roadRepo.CountParent(tx, roadID)
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_UPDATE_ROAD_HAVE_CHILD)
		} else if *countParent != 0 {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_UPDATE_ROAD_HAVE_CHILD)
		}
	}

	lineString := ""
	centerLinePath := ""
	centerLanePath := ""
	if req.CenterLaneShapeFileStatus == "upload" && req.CenterLineShapeFileStatus == "upload" {
		//Save File Center Line
		centerLinePath, err := t.helperFile.SaveFile(c, req.CenterLineShapeFile, centerLineDir)
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
		}

		//Save File Center Lane
		centerLanePath, err := t.helperFile.SaveFile(c, req.CenterLaneShapeFile, centerLaneDir)
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
		}

		// Convert Center Line Shape File

		centerLineInZip, err := shp.ShapesInZip(centerLinePath)
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LINE_FILE)

		}
		centerLineInZip = helpers.FilterByPrefix(centerLineInZip, "__MACOSX/")

		if len(centerLineInZip) > 1 {
			return nil, responses.NewAppErr(400, constants.HAS_MANY_SHAPE_FILE_IN_ZIP)
		}
		//Open Center Line Shape File From Zip
		shapeLine, err := shp.OpenShapeFromZip(centerLinePath, centerLineInZip[0])
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LINE_FILE)

		}

		var kmStartStr, kmEndStr string
		// loop through all features in the shapefile
		fields := shapeLine.Fields()
		for shapeLine.Next() {
			_, p := shapeLine.Shape()

			// Check if the geometry is a PolyLine
			polyline, ok := p.(*shp.PolyLine)
			if !ok {
				return nil, responses.NewAppErr(400, constants.INVALID_POLY_LINE)
			}

			//Convert Center Line Polyline To LineString
			lineString = convertPolylineToLineString(polyline)

			for k, f := range fields {
				val := shapeLine.Attribute(k)
				fieldName := strings.TrimSpace(strings.TrimRight(string(f.Name[:]), "\x00"))

				var handled bool
				switch fieldName {
				case "road_code":
					handled = true
				case "name":
					handled = true

				case "km_start":
					handled = true
					kmStartStr = val

				case "km_end":
					handled = true
					kmEndStr = val

				}
				if !handled {
					return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LINE_FILE)
				}
			}

		}

		kmStart, _ := helpers.KmStringToFloat64(kmStartStr)
		kmEnd, _ := helpers.KmStringToFloat64(kmEndStr)
		if kmStart != req.KmStart || kmEnd != req.KmEnd {
			return 0, responses.NewAppErr(400, constants.INVALID_ROAD_RANGE)
		}

		// Convert Center Lane Shape File
		centerLaneInZip, err := shp.ShapesInZip(centerLanePath)
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LANE_FILE)

		}
		centerLaneInZip = helpers.FilterByPrefix(centerLaneInZip, "__MACOSX/")
		if len(centerLaneInZip) > 1 {
			return nil, responses.NewAppErr(400, constants.HAS_MANY_SHAPE_FILE_IN_ZIP)
		}

		shapeLane, err := shp.OpenShapeFromZip(centerLanePath, centerLaneInZip[0])
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LANE_FILE)

		}

		var roadGeoms []models.RoadGeom
		// loop through all features in the shapefile
		for shapeLane.Next() {
			_, p := shapeLane.Shape()

			// Check if the geometry is a PolyLine
			polyline, ok := p.(*shp.PolyLine)
			if !ok {
				return nil, responses.NewAppErr(500, constants.INVALID_POLY_LINE)
			}

			lineString := convertPolylineToLineString(polyline)
			// shapeLaneGeom, err := t.roadRepo.ConvertLineStringToGeom(lineString)
			// if err != nil {
			// 	return 0, responses.NewAppErr(400, constants.INVALID_CONVERT_LINE_STRING_TO_GEOM)
			// }

			var roadGeom models.RoadGeom

			fields := shapeLane.Fields()
			for k, f := range fields {
				val := shapeLane.Attribute(k)
				fieldName := strings.TrimSpace(strings.TrimRight(string(f.Name[:]), "\x00"))
				// Convert [11]byte to string and trim any whitespace

				var handled bool
				switch fieldName {
				case "road_code":
					handled = true

				case "direction":
					handled = true
				case "lane_no":
					laneNo, err := strconv.Atoi(val)
					if err == nil {
						roadGeom.LaneNo = laneNo
					}
					handled = true
				case "km_start":

					val := strings.Replace(val, "+", "", -1)
					kmStart, err := strconv.ParseFloat(val, 64)
					if err == nil {
						roadGeom.KmStart = kmStart
					}
					handled = true
				case "km_end":

					val := strings.Replace(val, "+", "", -1)
					kmEnd, err := strconv.ParseFloat(val, 64)
					if err == nil {
						roadGeom.KmEnd = kmEnd
					}
					handled = true

				}

				if !handled {
					return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LANE_FILE)
				}
			}

			if err != nil {
				return 0, responses.NewAppErr(500, constants.INVALID_CONVERT_LINE_STRING_TO_GEOM)
			}
			roadGeom.TheGeom = lineString

			roadGeoms = append(roadGeoms, roadGeom)

		}
		checkGemo := false
		if len(roadGeoms) == 0 {
			checkGemo = true
			copier.Copy(&roadGeoms, road.RoadGeom)

		}

		for i := 0; i < len(roadGeoms); i++ {
			roadGeoms[i].RoadId = road.RoadId
			roadGeoms[i].Revision = road.Revision + 1
			roadGeoms[i].CreatedBy = userID
			roadGeoms[i].CreatedAt = time.Now()
			roadGeoms[i].UpdatedBy = userID
			roadGeoms[i].UpdatedAt = time.Now()
			roadGeoms[i].Status = "A"
		}

		err = t.roadRepo.DeleteRoadGeom(tx, road.RoadId, userID)
		if err != nil {
			tx.Rollback()
			return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
		}
		if checkGemo {
			err = t.roadRepo.UpdateRoadGeom(tx, roadGeoms)
			if err != nil {
				tx.Rollback()
				return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
			}
		} else {
			err = t.roadRepo.CreateRoadGeom(tx, roadGeoms)
			if err != nil {
				tx.Rollback()
				return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
			}
		}

	}

	err = t.roadRepo.DeleteRoadInfo(tx, road.RoadId, userID)
	if err != nil {
		tx.Rollback()
		return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
	}

	checkGemo := false
	if lineString == "" {
		checkGemo = true
		lineString = road.TheGeom
	}
	if centerLinePath == "" {
		centerLinePath = road.RoadInfo.CenterLineShapeFilePath
	}
	if centerLanePath == "" {
		centerLanePath = road.RoadInfo.CenterLaneShapeFilePath
	}

	var refDirectionId int
	if req.RefRoadTypeID == 1 || req.RefRoadTypeID == 3 {
		refDirectionId = 1
	} else if req.RefRoadTypeID == 2 || req.RefRoadTypeID == 4 {
		refDirectionId = 2
	} else {
		var roadParent models.Road
		err := t.roadRepo.GetDataById(&roadParent, roadID)
		if err != nil {
			return nil, responses.NewNotFoundError()
		}
		roadInfoParent, err := t.roadRepo.GetLastRoadInfoByID(*roadParent.ParentRoadId)
		if err != nil {
			return nil, responses.NewNotFoundError()
		}
		refDirectionId = roadInfoParent.RefDirectionId
	}

	// if refDirectionId != 0 && road.RoadInfo.RefDirectionId != refDirectionId {
	// 	err = t.roadRepo.UpdateDirectionRoad(tx, roadID, refDirectionId)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
	// 	}
	// }

	var roadInfo models.RoadInfo
	copier.Copy(&roadInfo, &req)

	roadInfo.RoadId = road.RoadId
	roadInfo.Year = &year
	roadInfo.TheGeom = lineString
	roadInfo.Status = "A"
	roadInfo.Revision = road.Revision + 1
	roadInfo.CenterLineShapeFilePath = centerLinePath
	roadInfo.CenterLaneShapeFilePath = centerLanePath
	roadInfo.CreatedAt = time.Now()
	roadInfo.CreatedBy = userID
	roadInfo.UpdatedAt = time.Now()
	roadInfo.UpdatedBy = userID
	roadInfo.RefDirectionId = refDirectionId
	roadInfo.YearConstructionCompleted = req.YearConstructionCompleted
	if checkGemo {
		err = t.roadRepo.UpdateRoadInfo(tx, &roadInfo)
		if err != nil {
			tx.Rollback()
			return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
		}
	} else {
		err = t.roadRepo.CreateRoadInfo(tx, &roadInfo)
		if err != nil {
			tx.Rollback()
			return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
		}
	}

	tx.Commit()
	return nil, nil
}

func (t *roadUseCase) UpdateRoadInit(c *gin.Context, roadID int, userID int, req requests.RoadUpdateInit) (interface{}, error) {

	tx := t.roadRepo.StartTransSection()
	if tx.Error != nil {
		return nil, responses.NewAppErr(500, "Failed to begin transaction")
	}

	centerLineDir := os.Getenv("ROAD_CENTER_LINE_SHAPE_FILE_DIR")
	centerLaneDir := os.Getenv("ROAD_CENTER_LANE_SHAPE_FILE_DIR")

	if centerLineDir == "" || centerLaneDir == "" {
		return nil, responses.NewAppErr(400, constants.FAILED_TO_READ_ENV_FILE)
	}

	var registerDate time.Time
	var year int
	if req.RegisterDate != "" {
		reDate, err := time.Parse("2006-01-02 15:04:05", req.RegisterDate)
		if err != nil {
			return 0, responses.NewAppErr(400, err.Error())
		}
		registerDate = reDate
		year = registerDate.Year()
	}

	road, err := t.roadRepo.GetLastRoadInfoByID(roadID)
	if err != nil {
		return nil, responses.NewNotFoundError()
	}

	var checkDirection int
	if req.RefRoadTypeID == 1 || req.RefRoadTypeID == 3 {
		checkDirection = 1
	} else if req.RefRoadTypeID == 2 || req.RefRoadTypeID == 4 {
		checkDirection = 2
	}

	if road.RoadInfo.RefDirectionId != checkDirection {
		countParent, err := t.roadRepo.CountParent(tx, roadID)
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_UPDATE_ROAD_HAVE_CHILD)
		} else if *countParent != 0 {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_UPDATE_ROAD_HAVE_CHILD)
		}
	}

	// if (road.RefRoadTypeID == 1 || road.RefRoadTypeID == 3) && (req.RefRoadTypeID != 1 && req.RefRoadTypeID != 3) {
	// 	if err != nil {
	// 		return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
	// 	}
	// } else if (road.RefRoadTypeID == 2 || road.RefRoadTypeID == 4) && (req.RefRoadTypeID != 2 && req.RefRoadTypeID != 4) {
	// 	if err != nil {
	// 		return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
	// 	}
	// }

	lineString := ""
	centerLinePath := ""
	centerLanePath := ""
	if req.CenterLaneShapeFileStatus == "upload" && req.CenterLineShapeFileStatus == "upload" {
		//Save File Center Line
		centerLinePath, err := t.helperFile.SaveFile(c, req.CenterLineShapeFile, centerLineDir)
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
		}

		//Save File Center Lane
		centerLanePath, err := t.helperFile.SaveFile(c, req.CenterLaneShapeFile, centerLaneDir)
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_SAVE_FILE)
		}

		// Convert Center Line Shape File

		centerLineInZip, err := shp.ShapesInZip(centerLinePath)
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LINE_FILE)

		}
		centerLineInZip = helpers.FilterByPrefix(centerLineInZip, "__MACOSX/")

		if len(centerLineInZip) > 1 {
			return nil, responses.NewAppErr(400, constants.HAS_MANY_SHAPE_FILE_IN_ZIP)
		}
		//Open Center Line Shape File From Zip
		shapeLine, err := shp.OpenShapeFromZip(centerLinePath, centerLineInZip[0])
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LINE_FILE)

		}

		var kmStartStr, kmEndStr string
		// loop through all features in the shapefile
		fields := shapeLine.Fields()
		for shapeLine.Next() {
			_, p := shapeLine.Shape()

			// Check if the geometry is a PolyLine
			polyline, ok := p.(*shp.PolyLine)
			if !ok {
				return nil, responses.NewAppErr(400, constants.INVALID_POLY_LINE)
			}

			//Convert Center Line Polyline To LineString
			lineString = convertPolylineToLineString(polyline)

			for k, f := range fields {
				val := shapeLine.Attribute(k)
				fieldName := strings.TrimSpace(strings.TrimRight(string(f.Name[:]), "\x00"))

				var handled bool
				switch fieldName {
				case "road_code":
					handled = true
				case "name":
					handled = true

				case "km_start":
					handled = true
					kmStartStr = val

				case "km_end":
					handled = true
					kmEndStr = val

				}
				if !handled {
					return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LINE_FILE)
				}
			}

		}

		kmStart, _ := helpers.KmStringToFloat64(kmStartStr)
		kmEnd, _ := helpers.KmStringToFloat64(kmEndStr)
		if kmStart != req.KmStart || kmEnd != req.KmEnd {
			return 0, responses.NewAppErr(400, constants.INVALID_ROAD_RANGE)
		}

		// Convert Center Lane Shape File
		centerLaneInZip, err := shp.ShapesInZip(centerLanePath)
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LANE_FILE)

		}
		centerLaneInZip = helpers.FilterByPrefix(centerLaneInZip, "__MACOSX/")
		if len(centerLaneInZip) > 1 {
			return nil, responses.NewAppErr(400, constants.HAS_MANY_SHAPE_FILE_IN_ZIP)
		}

		shapeLane, err := shp.OpenShapeFromZip(centerLanePath, centerLaneInZip[0])
		if err != nil {
			return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LANE_FILE)

		}

		var roadGeoms []models.RoadGeom
		// loop through all features in the shapefile
		for shapeLane.Next() {
			_, p := shapeLane.Shape()

			// Check if the geometry is a PolyLine
			polyline, ok := p.(*shp.PolyLine)
			if !ok {
				return nil, responses.NewAppErr(500, constants.INVALID_POLY_LINE)
			}

			lineString := convertPolylineToLineString(polyline)
			// shapeLaneGeom, err := t.roadRepo.ConvertLineStringToGeom(lineString)
			// if err != nil {
			// 	return 0, responses.NewAppErr(400, constants.INVALID_CONVERT_LINE_STRING_TO_GEOM)
			// }

			var roadGeom models.RoadGeom

			fields := shapeLane.Fields()
			for k, f := range fields {
				val := shapeLane.Attribute(k)
				fieldName := strings.TrimSpace(strings.TrimRight(string(f.Name[:]), "\x00"))
				// Convert [11]byte to string and trim any whitespace

				var handled bool
				switch fieldName {
				case "road_code":
					handled = true

				case "direction":
					handled = true
				case "lane_no":
					laneNo, err := strconv.Atoi(val)
					if err == nil {
						roadGeom.LaneNo = laneNo
					}
					handled = true
				case "km_start":

					val := strings.Replace(val, "+", "", -1)
					kmStart, err := strconv.ParseFloat(val, 64)
					if err == nil {
						roadGeom.KmStart = kmStart
					}
					handled = true
				case "km_end":

					val := strings.Replace(val, "+", "", -1)
					kmEnd, err := strconv.ParseFloat(val, 64)
					if err == nil {
						roadGeom.KmEnd = kmEnd
					}
					handled = true

				}

				if !handled {
					return nil, responses.NewAppErr(400, constants.FAILED_TO_UPLOAD_CENTER_LANE_FILE)
				}
			}

			if err != nil {
				return 0, responses.NewAppErr(500, constants.INVALID_CONVERT_LINE_STRING_TO_GEOM)
			}
			roadGeom.TheGeom = lineString

			roadGeoms = append(roadGeoms, roadGeom)

		}
		checkGemo := false
		if len(roadGeoms) == 0 {
			checkGemo = true
			copier.Copy(&roadGeoms, road.RoadGeom)

		}

		for i := 0; i < len(roadGeoms); i++ {
			roadGeoms[i].RoadId = road.RoadId
			roadGeoms[i].Revision = road.Revision + 1
			roadGeoms[i].CreatedBy = userID
			roadGeoms[i].CreatedAt = time.Now()
			roadGeoms[i].UpdatedBy = userID
			roadGeoms[i].UpdatedAt = time.Now()
			roadGeoms[i].Status = "A"
		}

		err = t.roadRepo.DeleteRoadGeom(tx, road.RoadId, userID)
		if err != nil {
			tx.Rollback()
			return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
		}
		if checkGemo {
			err = t.roadRepo.UpdateRoadGeom(tx, roadGeoms)
			if err != nil {
				tx.Rollback()
				return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
			}
		} else {
			err = t.roadRepo.CreateRoadGeom(tx, roadGeoms)
			if err != nil {
				tx.Rollback()
				return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
			}
		}

	}

	err = t.roadRepo.DeleteRoadInfo(tx, road.RoadId, userID)
	if err != nil {
		tx.Rollback()
		return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
	}
	checkGemo := false
	if lineString == "" {
		checkGemo = true
		lineString = road.TheGeom
	}
	if centerLinePath == "" {
		centerLinePath = road.RoadInfo.CenterLineShapeFilePath
	}
	if centerLanePath == "" {
		centerLanePath = road.RoadInfo.CenterLaneShapeFilePath
	}

	var refDirectionId int
	if req.RefRoadTypeID == 1 || req.RefRoadTypeID == 3 {
		refDirectionId = 1
	} else if req.RefRoadTypeID == 2 || req.RefRoadTypeID == 4 {
		refDirectionId = 2
	} else {
		var roadParent models.Road
		err := t.roadRepo.GetDataById(&roadParent, roadID)
		if err != nil {
			return nil, responses.NewNotFoundError()
		}
		roadInfoParent, err := t.roadRepo.GetLastRoadInfoByID(*roadParent.ParentRoadId)
		if err != nil {
			return nil, responses.NewNotFoundError()
		}
		refDirectionId = roadInfoParent.RefDirectionId
	}

	// if refDirectionId != 0 && road.RoadInfo.RefDirectionId != refDirectionId {
	// 	err = t.roadRepo.UpdateDirectionRoad(tx, roadID, refDirectionId)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
	// 	}
	// }

	var roadInfo models.RoadInfo
	copier.Copy(&roadInfo, &req)

	roadInfo.RoadId = road.RoadId
	roadInfo.Year = &year
	roadInfo.TheGeom = lineString
	roadInfo.Status = "A"
	roadInfo.Revision = road.Revision + 1
	roadInfo.CenterLineShapeFilePath = centerLinePath
	roadInfo.CenterLaneShapeFilePath = centerLanePath
	roadInfo.CreatedAt = time.Now()
	roadInfo.CreatedBy = userID
	roadInfo.UpdatedAt = time.Now()
	roadInfo.UpdatedBy = userID
	roadInfo.RefDirectionId = refDirectionId
	roadInfo.YearConstructionCompleted = req.YearConstructionCompleted
	if checkGemo {
		err = t.roadRepo.UpdateRoadInfo(tx, &roadInfo)
		if err != nil {
			tx.Rollback()
			return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
		}
	} else {
		err = t.roadRepo.CreateRoadInfo(tx, &roadInfo)
		if err != nil {
			tx.Rollback()
			return nil, responses.NewAppErr(500, constants.FAILED_TO_UPDATE_ROAD)
		}
	}

	tx.Commit()
	return nil, nil
}

func (t *roadUseCase) DeleteRoad(roadID int, userID int) (interface{}, error) {

	tx := t.roadRepo.StartTransSection()

	err := t.roadRepo.DeleteRoad(tx, roadID, userID)
	if err != nil {
		tx.Rollback()
		return nil, responses.NewAppErr(500, constants.FAILED_TO_DELETE_ROAD)
	}

	err = t.roadRepo.DeleteRoadInfo(tx, roadID, userID)
	if err != nil {
		tx.Rollback()
		return nil, responses.NewAppErr(500, constants.FAILED_TO_DELETE_ROAD)
	}

	err = t.roadRepo.DeleteRoadGeom(tx, roadID, userID)
	if err != nil {
		tx.Rollback()
		return nil, responses.NewAppErr(500, constants.FAILED_TO_DELETE_ROAD)
	}
	tx.Commit()

	return nil, nil
}

func (t *roadUseCase) GetLastRoadInfoByID(roadID int) (*models.RoadInfoGeomData, error) {
	roadInfo, err := t.roadRepo.GetLastRoadInfoByID(roadID)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	return &roadInfo, nil
}

func (t *roadUseCase) GetRoadLanes(roadID int) ([]responses.RoadLanes, error) {
	roadLanse, err := t.roadRepo.GetRoadLanes(roadID)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}

	var roadLanseRes []responses.RoadLanes
	copier.Copy(&roadLanseRes, &roadLanse)

	return roadLanseRes, nil
}

func convertPolylineToLineString(polyline *shp.PolyLine) string {
	var points []string
	for _, point := range polyline.Points {
		points = append(points, fmt.Sprintf("%f %f", point.X, point.Y))
	}
	return fmt.Sprintf(`LINESTRING(%s)`, strings.Join(points, ", "))
}
