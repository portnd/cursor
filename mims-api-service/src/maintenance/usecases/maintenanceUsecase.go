package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/jinzhu/copier"
	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	servicesDB "gitlab.com/mims-api-service/services/database"
	"gitlab.com/mims-api-service/src/maintenance/domains"
	"gorm.io/gorm"
)

type Usecase struct {
	Repo       domains.Repository
	servicesDB servicesDB.ServicesDatabaseDomain
}

func NewUsecase(repo domains.Repository, servicesDB servicesDB.ServicesDatabaseDomain) domains.Usecase {
	return &Usecase{Repo: repo, servicesDB: servicesDB}
}

func (u *Usecase) CalculateDistance(request requests.CalDistance) (float64, error) {
	// roadLane, err := u.Repo.GetLaneByRoadID(request.RoadID)
	// if err != nil {
	// 	logs.Error(err)
	// 	return 0, err
	// }
	// if len(roadLane) == 0 {
	// 	return 0, gorm.ErrRecordNotFound
	// }

	// var oldStart float64
	// var oldEnd float64
	// var allNumber []float64
	// for _, v := range roadLane {
	// 	allNumber = append(allNumber, v.KmStart)
	// 	allNumber = append(allNumber, v.KmEnd)
	// }
	// if request.DirectionID == 1 {
	// 	oldStart, oldEnd = helpers.FindMinMaxFloat64(allNumber)
	// } else if request.DirectionID == 2 {
	// 	oldEnd, oldStart = helpers.FindMinMaxFloat64(allNumber)
	// }
	// // check if in road_geom
	// if request.DirectionID == 1 {
	// 	if request.KmStart < oldStart || request.KmEnd > oldEnd {
	// 		err = errors.New(constants.INVALID_GEOM_RANGE)
	// 		logs.Error(err)
	// 		return 0, err
	// 	}
	// } else if request.DirectionID == 2 {
	// 	if request.KmStart > oldStart || request.KmEnd < oldEnd {
	// 		err = errors.New(constants.INVALID_GEOM_RANGE)
	// 		logs.Error(err)
	// 		return 0, err
	// 	}
	// }
	// distance := math.Abs(request.KmEnd - request.KmStart)
	return 11, nil
}

func (u *Usecase) GetRefCriteriaMethod() (interface{}, error) {
	data, err := u.Repo.GetRefCriteriaMethod()
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return data, responses.NewAppErr(400, err.Error())
	}
	// return data, nil
	var interventionCriteriasRes []responses.MaintenanceInterventionCriteria

	for _, data := range data {
		var interventionCriterias responses.MaintenanceInterventionCriteria
		if len(data.Children) != 0 && data.Children != nil {
			copier.Copy(&interventionCriterias, &data)
			interventionCriteriasRes = append(interventionCriteriasRes, interventionCriterias)
		}
	}
	// for _, item := range data {
	// 	itvcData, err := u.Repo.GetInterventionCriteria(item.ID)
	// 	if len(itvcData) == 0 {
	// 		continue
	// 	}
	// 	if err != nil {
	// 		logs.Error(err)
	// 		if err == gorm.ErrRecordNotFound {
	// 			return responses.NoData{}, nil
	// 		}
	// 		return data, responses.NewAppErr(400, err.Error())
	// 	}
	// 	var interventionCriteria responses.MaintenanceInterventionCriteria
	// 	interventionCriteria.ID = item.ID
	// 	interventionCriteria.Label = item.Name
	// 	var childrens []responses.MaintenanceInterventionCriteriaChildren
	// 	for _, item := range itvcData {
	// 		var children responses.MaintenanceInterventionCriteriaChildren
	// 		// var dd interventionCriteria.Children
	// 		children.ID = item.Id
	// 		children.Label = item.MaintenanceStandardName
	// 		childrens = append(childrens, children)
	// 	}
	// 	interventionCriteria.Children = childrens
	// 	interventionCriterias = append(interventionCriterias, interventionCriteria)
	// }
	return interventionCriteriasRes, nil
}

func (u *Usecase) GetMaintenanceBudget() (interface{}, error) {
	data, err := u.Repo.GetMaintenanceBudget()
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return "", responses.NewAppErr(400, err.Error())
	}
	return data, nil
}

func (u *Usecase) GetMaintenanceByRoadID(roadID int, year int) (interface{}, error) {
	responds := []responses.MaintenanceByRoadId{}
	result, err := u.Repo.GetMaintenanceByRoadID(roadID, year)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}

	for _, data := range result {
		var respond responses.MaintenanceByRoadId
		copier.Copy(&respond, &data)
		guaranteeExpirationDate := data.GuaranteeExpirationDate.Unix()
		projectEndDate := data.ProjectEndDate.Unix()
		totalDay := guaranteeExpirationDate - projectEndDate
		leftDay := guaranteeExpirationDate - time.Now().Unix()

		var percent int64
		if totalDay != 0 {
			percent = (leftDay * 100) / totalDay
		} else {
			percent = 0
		}

		var colorCode string
		if percent > 20 {
			colorCode = "#1F70F3"
		} else {
			colorCode = "#F1416C"
		}
		respond.CreatedBy = responses.UserBy{
			Id:         data.CreatedByUser.Id,
			Name:       data.CreatedByUser.Name,
			DepartName: data.CreatedByUser.Department.Name,
		}
		respond.UpdateBy = responses.UserBy{
			Id:         data.UpdateByUser.Id,
			Name:       data.UpdateByUser.Name,
			DepartName: data.UpdateByUser.Department.Name,
		}
		hostPath := os.Getenv("STORAGE_IP") + "/"

		if data.CreatedByUser.ProfileImgPath != "" {
			respond.CreatedBy.ProfileImgPath = hostPath + data.CreatedByUser.ProfileImgPath
		}

		if data.UpdateByUser.ProfileImgPath != "" {
			respond.UpdateBy.ProfileImgPath = hostPath + data.UpdateByUser.ProfileImgPath
		}

		respond.Color = colorCode
		//respond.GuaranteeDay = totalDay / 86400

		if respond.RoadGroupNames != nil {
			for _, data := range data.RoadGroupNames {

				roadGroupName := fmt.Sprintf(`ทางหลวงพิเศษหมายเลข %s`, data)
				respond.RoadGroupNames = append(respond.RoadGroupNames, roadGroupName)
			}

		} else {
			respond.RoadGroupNames = []string{}
		}

		var ownerCode string
		var ownerName string
		if data.RefDivision.DivisionCode != "" {
			ownerCode = fmt.Sprintf(`ref_division_code-%s`, data.RefDivision.DivisionCode)
			ownerName = data.RefDivision.Name
		} else if data.RefDistrict.DistrictCode != "" {
			ownerCode = fmt.Sprintf(`ref_district_code-%s`, data.RefDistrict.DistrictCode)
			ownerName = data.RefDistrict.Name
		} else if data.RefDepot.DepotCode != "" {
			ownerCode = fmt.Sprintf(`ref_depot_code-%s`, data.RefDepot.DepotCode)
			ownerName = data.RefDepot.Name
		}

		respond.OwnerCode = ownerCode
		respond.OwnerName = ownerName

		for i, att := range respond.Attachments {
			respond.Attachments[i].Path = hostPath + att.Path
		}

		if len(data.Roads) != 0 {
			for i, mRoad := range data.Roads {
				respond.Roads[i].Color = colorCode

				var geomJSON responses.GeomJSON
				json.Unmarshal([]byte(mRoad.TheGeomString), &geomJSON)

				respond.Roads[i].TheGeom = geomJSON

				if mRoad.RoadLevel == 1 {
					if mRoad.RoadSecNameDes != "" {
						if mRoad.RefDirectionId == 1 {
							respond.Roads[i].RoadName = fmt.Sprintf(`%s - %s`, mRoad.RoadSecNameOr, mRoad.RoadSecNameDes)
						} else if mRoad.RefDirectionId == 2 {
							respond.Roads[i].RoadName = fmt.Sprintf(`%s - %s`, mRoad.RoadSecNameDes, mRoad.RoadSecNameOr)
						}
					} else {
						respond.Roads[i].RoadName = mRoad.RoadSecNameOr
					}
				}
			}
		} else {
			respond.Roads = []responses.MaintenanceRoadPreload{}
		}

		if len(data.RoadHistories) != 0 {
			for i, mRoadHis := range data.RoadHistories {

				var geomJSON responses.GeomJSON
				json.Unmarshal([]byte(mRoadHis.TheGeom), &geomJSON)

				respond.RoadHistories[i].TheGeom = geomJSON

				respond.RoadHistories[i].Color = colorCode
				if mRoadHis.RoadLevel == 1 {
					if mRoadHis.RoadSecNameDes != "" {
						if mRoadHis.RefDirectionId == 1 {
							respond.Roads[i].RoadName = fmt.Sprintf(`%s - %s`, mRoadHis.RoadSecNameOr, mRoadHis.RoadSecNameDes)
						} else if mRoadHis.RefDirectionId == 2 {
							respond.Roads[i].RoadName = fmt.Sprintf(`%s - %s`, mRoadHis.RoadSecNameDes, mRoadHis.RoadSecNameOr)
						}
					} else {
						respond.Roads[i].RoadName = mRoadHis.RoadSecNameOr
					}
				}
			}
		} else {
			respond.RoadHistories = []responses.MaintenanceRoadPreload{}
		}
		responds = append(responds, respond)
	}

	return responds, nil
}

func (u *Usecase) GetRoadMaintenanceYear(roadID int) (interface{}, error) {
	data, _ := u.Repo.GetRoadMaintenanceYear(roadID)
	year := []int{}
	for _, item := range data {
		year = append(year, item.BudgetYear)
	}
	return year, nil
}

func (u *Usecase) GetMaintenanceList(userID int, prams requests.MaintenancePrams, limit, offset int64) ([]responses.MaintenanceList, int64, error) {
	//============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_maint_history"}, []string{"view_owner_maint_history"})
	//============ end check permission ============
	var responds []responses.MaintenanceList
	queryFilter := createFilter(prams)
	total, err := u.Repo.GetMaintenanceListCount(queryFilter, prams, isAllData, isOwnerData, depotCode)
	if err != nil {
		logs.Error(err)
		return nil, 0, err
	}
	if total == 0 {
		return []responses.MaintenanceList{}, 0, nil
	}
	result, err := u.Repo.GetMaintenanceList(queryFilter, prams, isAllData, isOwnerData, depotCode, limit, offset)
	if err != nil {
		logs.Error(err)
		return nil, 0, err
	}
	if len(result) == 0 {
		return []responses.MaintenanceList{}, total, nil
	}

	// Batch load geometry for all roads (avoids N+1: was one query per road)
	roadIDs := make([]int, 0, len(result)*4)
	for _, item := range result {
		for _, r := range item.Roads {
			roadIDs = append(roadIDs, r.ID)
		}
	}
	geomByID, err := u.Repo.GetGeomJsonByMaintenanceRoadIDs(roadIDs)
	if err != nil {
		logs.Error(err)
		return nil, 0, err
	}

	for _, item := range result {
		var respond responses.MaintenanceList
		copier.Copy(&respond, &item)
		guaranteeExpirationDate := item.GuaranteeExpirationDate.Unix()
		projectEndDate := item.ProjectEndDate.Unix()
		totalDay := guaranteeExpirationDate - projectEndDate
		leftDay := guaranteeExpirationDate - time.Now().Unix()

		//percent := (leftDay * 100) / totalDay

		var percent int64
		if totalDay != 0 {
			percent = (leftDay * 100) / totalDay
		} else {
			percent = 0
		}

		respond.CreatedBy = responses.UserBy{
			Id:         item.CreatedByUser.Id,
			Name:       item.CreatedByUser.Name,
			DepartName: item.CreatedByUser.Department.Name,
		}
		respond.UpdateBy = responses.UserBy{
			Id:         item.UpdateByUser.Id,
			Name:       item.UpdateByUser.Name,
			DepartName: item.UpdateByUser.Department.Name,
		}
		hostPathImage := os.Getenv("STORAGE_IP") + "/"

		if respond.CreatedBy.ProfileImgPath != "" {
			respond.CreatedBy.ProfileImgPath = hostPathImage + item.CreatedByUser.ProfileImgPath
		}

		if respond.UpdateBy.ProfileImgPath != "" {
			respond.UpdateBy.ProfileImgPath = hostPathImage + item.UpdateByUser.ProfileImgPath
		}

		var colorCode string
		if percent > 20 {
			colorCode = "#1F70F3"
		} else {
			colorCode = "#F1416C"
		}

		respond.Color = colorCode
		nowDateString := time.Now().Format("2006-01-02")

		// Convert nowDateString back to a time.Time object for comparison
		nowDate, err := time.Parse("2006-01-02", nowDateString)
		if err != nil {
			return nil, 0, err
		}

		leftTime := respond.GuaranteeExpirationDate.Sub(nowDate)
		leftTimeDays := int(leftTime.Hours()/24) + 1

		if nowDate.After(respond.GuaranteeExpirationDate) {
			respond.RemainingTime = "หมดการค้ำประกัน"
		} else {
			respond.RemainingTime = "เหลือระยะเวลาค้ำประกัน" + " " + helpers.FormatNumber(leftTimeDays) + " วัน"

		}

		roadGroupNames := []string{}
		if respond.RoadGroupNames != nil {
			for _, data := range respond.RoadGroupNames {

				roadGroupName := fmt.Sprintf(`ทางหลวงพิเศษหมายเลข %s`, data)
				roadGroupNames = append(roadGroupNames, roadGroupName)
			}

		}
		respond.RoadGroupNames = roadGroupNames

		var ownerCode string
		var ownerName string
		if item.RefDivision.DivisionCode != "" {
			ownerCode = fmt.Sprintf(`ref_division_code-%s`, item.RefDivision.DivisionCode)
			ownerName = item.RefDivision.Name
		} else if item.RefDistrict.DistrictCode != "" {
			ownerCode = fmt.Sprintf(`ref_district_code-%s`, item.RefDistrict.DistrictCode)
			ownerName = item.RefDistrict.Name
		} else if item.RefDepot.DepotCode != "" {
			ownerCode = fmt.Sprintf(`ref_depot_code-%s`, item.RefDepot.DepotCode)
			ownerName = item.RefDepot.Name
		}

		respond.OwnerCode = ownerCode
		respond.OwnerName = ownerName
		if len(item.Roads) != 0 {
			for i, mRoad := range item.Roads {

				geomJSON := geomByID[mRoad.ID]
				if len(geomJSON) > 0 {
					theGeom, err := helpers.ConvertThegeomToGeomJSON(geomJSON)
					if err != nil {
						return nil, 0, err
					}
					respond.Roads[i].TheGeom = theGeom
				}

				if mRoad.LaneNo == 0 {
					respond.Roads[i].LaneNo = nil
				} else {
					respond.Roads[i].GridNo = nil
				}

				respond.Roads[i].Color = colorCode
				if mRoad.RoadLevel == 1 {
					if mRoad.RoadSecNameDes != "" {
						if mRoad.RefDirectionId == 1 {
							respond.Roads[i].RoadName = fmt.Sprintf(`%s - %s`, mRoad.RoadSecNameOr, mRoad.RoadSecNameDes)
						} else if mRoad.RefDirectionId == 2 {
							respond.Roads[i].RoadName = fmt.Sprintf(`%s - %s`, mRoad.RoadSecNameDes, mRoad.RoadSecNameOr)
						}
					} else {
						respond.Roads[i].RoadName = mRoad.RoadSecNameOr
					}

				}
			}
		} else {
			respond.Roads = []responses.MaintenanceRoadPreload{}
		}

		// ไม่มี filter ทางสายทาง: แสดงทั้งหมด (รวมรายการที่ไม่มี road_group_names)
		if (len(prams.RoadGroupID) == 0 || prams.RoadGroupID == nil) && (len(prams.RoadGroupIDDashboard) == 0 || prams.RoadGroupIDDashboard == nil) {
			responds = append(responds, respond)
		} else if len(roadGroupNames) > 0 {
			if (len(prams.RoadGroupID) != 0) && (item.Roads != nil && len(item.Roads) != 0) {
				responds = append(responds, respond)
			} else if len(prams.RoadGroupIDDashboard) > 1 && helpers.SlicesAreEqual(prams.RoadGroupIDDashboard, item.RoadGroupID) {
				responds = append(responds, respond)
			} else if len(prams.RoadGroupIDDashboard) == 1 && helpers.IntInSlice(prams.RoadGroupIDDashboard[0], item.RoadGroupID) {
				responds = append(responds, respond)
			}
		}

	}

	return responds, total, nil
}

func (u *Usecase) GetMaintenanceListByID(IDParent int) (interface{}, error) {
	data, err := u.Repo.GetMaintenanceListByID(IDParent)
	if err != nil {
		return data, responses.NewAppErr(404, err.Error())
	}
	// return data, nil
	var respond responses.MaintenanceById
	copier.Copy(&respond, &data)

	budgetMethod, err := u.Repo.GetSettingBudgetMethodByID(data.BudgetMethodID)
	if err != nil {
		logs.Error(err)
		return false, responses.NewAppErr(500, err.Error())
	}

	respond.IsShowMethod = budgetMethod.IsShowMethod

	guaranteeExpirationDate := data.GuaranteeExpirationDate.Unix()
	projectEndDate := data.ProjectEndDate.Unix()
	totalDay := guaranteeExpirationDate - projectEndDate
	leftDay := guaranteeExpirationDate - time.Now().Unix()

	var percent int64
	if totalDay != 0 {
		percent = (leftDay * 100) / totalDay
	} else {
		percent = 0
	}

	var colorCode string
	if percent > 20 {
		colorCode = "#1F70F3"
	} else {
		colorCode = "#F1416C"
	}
	respond.CreatedBy = responses.UserBy{
		Id:         data.CreatedByUser.Id,
		Name:       data.CreatedByUser.Name,
		DepartName: data.CreatedByUser.Department.Name,
	}
	respond.UpdateBy = responses.UserBy{
		Id:         data.UpdateByUser.Id,
		Name:       data.UpdateByUser.Name,
		DepartName: data.UpdateByUser.Department.Name,
	}
	hostPath := os.Getenv("STORAGE_IP") + "/"

	if respond.CreatedBy.ProfileImgPath != "" {
		respond.CreatedBy.ProfileImgPath = hostPath + data.CreatedByUser.ProfileImgPath
	}

	if respond.UpdateBy.ProfileImgPath != "" {
		respond.UpdateBy.ProfileImgPath = hostPath + data.UpdateByUser.ProfileImgPath
	}
	respond.Color = colorCode
	//respond.GuaranteeDay = totalDay / 86400

	roadGroupNames := []string{}
	if respond.RoadGroupNames != nil {
		for _, data := range respond.RoadGroupNames {

			roadGroupName := fmt.Sprintf(`ทางหลวงพิเศษหมายเลข %s`, data)
			roadGroupNames = append(roadGroupNames, roadGroupName)
		}

	}
	respond.RoadGroupNames = roadGroupNames

	nowDateString := time.Now().Format("2006-01-02")

	// Convert nowDateString back to a time.Time object for comparison
	nowDate, err := time.Parse("2006-01-02", nowDateString)
	if err != nil {
		return nil, err
	}

	leftTime := data.GuaranteeExpirationDate.Sub(nowDate)
	leftTimeDays := int(leftTime.Hours()/24) + 1

	if nowDate.After(respond.GuaranteeExpirationDate) {
		respond.RemainingTime = "หมดการค้ำประกัน"
	} else {
		respond.RemainingTime = "เหลือระยะเวลาค้ำประกัน" + " " + helpers.FormatNumber(leftTimeDays) + " วัน"
	}

	var ownerCode string
	var ownerName string
	// if data.RefDivision.DivisionCode != "" {
	// 	ownerCode = fmt.Sprintf(`ref_division_code-%s`, data.RefDivision.DivisionCode)
	// 	ownerName = data.RefDivision.Name
	// } else if data.RefDistrict.DistrictCode != "" {
	// 	ownerCode = fmt.Sprintf(`ref_district_code-%s`, data.RefDistrict.DistrictCode)
	// 	ownerName = data.RefDistrict.Name
	// } else if data.RefDepot.DepotCode != "" {
	// 	ownerCode = fmt.Sprintf(`ref_depot_code-%s`, data.RefDepot.DepotCode)
	// 	ownerName = data.RefDepot.Name
	// }

	ownerCode = fmt.Sprintf(`ref_depot_code-%s`, data.RefDepot.DepotCode)
	ownerName = data.RefDepot.Name

	respond.OwnerCode = ownerCode
	respond.OwnerName = ownerName

	for i, att := range respond.Attachments {
		respond.Attachments[i].Path = hostPath + att.Path
	}

	// Batch load geometry for Roads and RoadHistories (avoids N+1)
	roadIDs := make([]int, 0, len(data.Roads))
	for _, r := range data.Roads {
		roadIDs = append(roadIDs, r.ID)
	}
	roadHistoryIDs := make([]int, 0, len(data.RoadHistories))
	for _, r := range data.RoadHistories {
		roadHistoryIDs = append(roadHistoryIDs, r.ID)
	}
	geomByRoadID, err := u.Repo.GetGeomJsonByMaintenanceRoadIDs(roadIDs)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	geomByRoadHistoryID, err := u.Repo.GetGeomJsonByMaintenanceRoadHistoryIDs(roadHistoryIDs)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	if len(data.Roads) != 0 {
		for i, mRoad := range data.Roads {
			respond.Roads[i].Color = colorCode

			geomJSON := geomByRoadID[mRoad.ID]
			if len(geomJSON) > 0 {
				theGeom, err := helpers.ConvertThegeomToGeomJSON(geomJSON)
				if err != nil {
					return nil, err
				}
				respond.Roads[i].TheGeom = theGeom
			}

			if mRoad.LaneNo == 0 {
				respond.Roads[i].LaneNo = nil
			}

			if mRoad.GridNo == 0 {
				respond.Roads[i].GridNo = nil
			}

			if mRoad.RoadLevel == 1 {
				if mRoad.RoadSecNameDes != "" {
					if mRoad.RefDirectionId == 1 {
						respond.Roads[i].RoadName = fmt.Sprintf(`%s - %s`, mRoad.RoadSecNameOr, mRoad.RoadSecNameDes)
					} else if mRoad.RefDirectionId == 2 {
						respond.Roads[i].RoadName = fmt.Sprintf(`%s - %s`, mRoad.RoadSecNameDes, mRoad.RoadSecNameOr)
					}
				} else {
					respond.Roads[i].RoadName = mRoad.RoadSecNameOr
				}
			}
		}
	} else {
		respond.Roads = []responses.MaintenanceRoadPreload{}
	}

	if len(data.RoadHistories) != 0 {
		for i, mRoadHis := range data.RoadHistories {
			respond.RoadHistories[i].Color = colorCode

			geomJSON := geomByRoadHistoryID[mRoadHis.ID]
			if len(geomJSON) > 0 {
				theGeom, err := helpers.ConvertThegeomToGeomJSON(geomJSON)
				if err != nil {
					return nil, err
				}
				respond.RoadHistories[i].TheGeom = theGeom
			}

			if mRoadHis.LaneNo == 0 {
				respond.RoadHistories[i].LaneNo = nil
			}

			if mRoadHis.GridNo == 0 {
				respond.RoadHistories[i].GridNo = nil
			}

			if mRoadHis.RoadLevel == 1 {
				if mRoadHis.RoadSecNameDes != "" {
					if mRoadHis.RefDirectionId == 1 {
						respond.RoadHistories[i].RoadName = fmt.Sprintf(`%s - %s`, mRoadHis.RoadSecNameOr, mRoadHis.RoadSecNameDes)
					} else if mRoadHis.RefDirectionId == 2 {
						respond.RoadHistories[i].RoadName = fmt.Sprintf(`%s - %s`, mRoadHis.RoadSecNameDes, mRoadHis.RoadSecNameOr)
					}
				} else {
					respond.RoadHistories[i].RoadName = mRoadHis.RoadSecNameOr
				}
			}
		}
	} else {
		respond.RoadHistories = []responses.MaintenanceRoadPreload{}
	}

	return respond, nil
}

func createFilter(params requests.MaintenancePrams) string {
	var queryWhere []string
	if params.BudgetYear != nil {
		queryWhere = append(queryWhere, "maintenance.budget_year = "+fmt.Sprintf("%d", *params.BudgetYear))
	}

	if params.OwnerCode != nil {
		ownerCode := *params.OwnerCode
		ownerCodeSplit := strings.Split(ownerCode, "-")
		if len(ownerCodeSplit) >= 2 {
			queryWhere = append(queryWhere, fmt.Sprintf(`maintenance.%s = '%s'`, ownerCodeSplit[0], ownerCodeSplit[1]))
		}
	}

	if params.Name != nil {
		queryWhere = append(queryWhere, "maintenance.name like '%"+*params.Name+"%'")
	}
	query := strings.Join(queryWhere, " and ")
	return query
}

func (u *Usecase) CreateMaintenance(request requests.MaintenanceReq, userID int, attReqs []requests.MaintenanceAttachmentsReq) (interface{}, error) {
	var maintenance models.Maintenance

	copier.Copy(&maintenance, &request)

	guaranteeExpirationDate, err := time.Parse("2006-01-02", request.GuaranteeExpirationDate)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	projectEndDate, err := time.Parse("2006-01-02", request.ProjectEndDate)
	if err != nil {
		return "", responses.NewAppErr(400, err.Error())
	}

	maintenance.GuaranteeExpirationDate = guaranteeExpirationDate
	maintenance.ProjectEndDate = projectEndDate
	maintenance.IDParent = nil
	maintenance.UpdatedAt = time.Now().UTC()
	maintenance.CreatedAt = time.Now().UTC()
	maintenance.UpdatedBy = userID
	maintenance.CreatedBy = userID
	maintenance.Status = "A"
	maintenance.RefDepotCode = &request.OwnerCode

	ownerCodeSplit := strings.Split(request.OwnerCode, "-")
	if len(ownerCodeSplit) < 2 {
		return "", responses.NewAppErr(400, "owner_code invalid format please check")
	}

	prefix := ownerCodeSplit[0]
	ownerCode := ownerCodeSplit[1]

	if prefix == "ref_division_code" {
		maintenance.RefDivisionCode = &ownerCode
	} else if prefix == "ref_district_code" {
		maintenance.RefDistrictCode = &ownerCode
	} else if prefix == "ref_depot_code" {
		maintenance.RefDepotCode = &ownerCode
	}

	tx := u.Repo.StartTransSection()
	id, err := u.Repo.InsertMaintenance(maintenance, tx)
	if err != nil {
		u.Repo.RollBack(tx)
		logs.Error(err)
		return "", responses.NewAppErr(500, err.Error())
	}

	err = u.Repo.UpdateIDParentMaintenance(id, id, tx)
	if err != nil {
		u.Repo.RollBack(tx)
		logs.Error(err)
		return "", responses.NewAppErr(500, err.Error())
	}

	if len(attReqs) != 0 {
		var maintenanceAttachment []models.MaintenanceAttachment
		for _, att := range attReqs {
			if att.Status == "upload" {
				var maintenanceAtt models.MaintenanceAttachment
				copier.Copy(&maintenanceAtt, &att)
				maintenanceAtt.MaintenanceID = id
				maintenanceAtt.CreatedBy = userID
				maintenanceAtt.UpdatedBy = userID
				maintenanceAtt.CreatedAt = time.Now()
				maintenanceAtt.UpdatedAt = time.Now()
				maintenanceAttachment = append(maintenanceAttachment, maintenanceAtt)
			}
		}

		err = u.Repo.InsertMaintenanceaAttachment(maintenanceAttachment, tx)
		if err != nil {
			u.Repo.RollBack(tx)
			logs.Error(err)
			return "", responses.NewAppErr(500, err.Error())
		}

	}

	u.Repo.Commit(tx)
	return id, nil
}

func (u *Usecase) UpdateMaintenance(idParent int, request requests.MaintenanceReq, userID int, attReqs []requests.MaintenanceAttachmentsReq) (*int, *int, error) {

	var maintenance models.Maintenance
	copier.Copy(&maintenance, &request)

	revision, maintenanceID, err := u.Repo.GetMaxRevisionMaintenanceByIDParent(idParent)
	if err != nil {
		logs.Error(err)
		return nil, nil, responses.NewAppErr(500, err.Error())
	}
	guaranteeExpirationDate, err := time.Parse("2006-01-02", request.GuaranteeExpirationDate)
	if err != nil {
		return nil, nil, responses.NewAppErr(400, err.Error())
	}

	projectEndDate, err := time.Parse("2006-01-02", request.ProjectEndDate)
	if err != nil {
		return nil, nil, responses.NewAppErr(400, err.Error())
	}

	maintenance.GuaranteeExpirationDate = guaranteeExpirationDate
	maintenance.ProjectEndDate = projectEndDate
	maintenance.IDParent = &idParent
	maintenance.UpdatedAt = time.Now().UTC()
	maintenance.CreatedAt = time.Now().UTC()
	maintenance.UpdatedBy = userID
	maintenance.CreatedBy = userID
	maintenance.Status = "A"
	maintenance.Revision = *revision + 1
	maintenance.RefDepotCode = &request.OwnerCode

	ownerCodeSplit := strings.Split(request.OwnerCode, "-")
	if len(ownerCodeSplit) < 2 {
		return nil, nil, responses.NewAppErr(400, "owner_code invalid format please check")
	}

	prefix := ownerCodeSplit[0]
	ownerCode := ownerCodeSplit[1]

	if prefix == "ref_division_code" {
		maintenance.RefDivisionCode = &ownerCode
	} else if prefix == "ref_district_code" {
		maintenance.RefDistrictCode = &ownerCode
	} else if prefix == "ref_depot_code" {
		maintenance.RefDepotCode = &ownerCode
	}

	tx := u.Repo.StartTransSection()
	id, err := u.Repo.InsertMaintenance(maintenance, tx)
	if err != nil {
		u.Repo.RollBack(tx)
		logs.Error(err)
		return nil, nil, responses.NewAppErr(500, err.Error())
	}

	err = u.Repo.UpdateStatusMaintenance(*maintenanceID, tx)
	if err != nil {
		u.Repo.RollBack(tx)
		logs.Error(err)
		return nil, nil, responses.NewAppErr(500, err.Error())
	}

	resultMTR, err := u.Repo.GetMaintenanceRoadByMaintenanceID(*maintenanceID)
	if err != nil {
		u.Repo.RollBack(tx)
		logs.Error(err)
		return nil, nil, responses.NewAppErr(500, err.Error())
	}

	resultMTRH, err := u.Repo.GetMaintenanceRoadHistoryByMaintenanceID(*maintenanceID)
	if err != nil {
		u.Repo.RollBack(tx)
		logs.Error(err)
		return nil, nil, responses.NewAppErr(500, err.Error())
	}

	if len(resultMTR) != 0 {
		var ids []int
		for i, data := range resultMTR {
			resultMTR[i].ID = 0
			resultMTR[i].MaintenanceID = id
			resultMTR[i].Revision = data.Revision + 1
			ids = append(ids, data.ID)
		}
		err = u.Repo.InsertMaintenanceRoad(resultMTR, tx)
		if err != nil {
			u.Repo.RollBack(tx)
			logs.Error(err)
			return nil, nil, responses.NewAppErr(500, err.Error())
		}

		err = u.Repo.UpdateStatusMaintenanceRoad(ids, tx)
		if err != nil {
			u.Repo.RollBack(tx)
			logs.Error(err)
			return nil, nil, responses.NewAppErr(500, err.Error())
		}
	}
	helpers.PrintlnJson(resultMTR)
	var resultMTRHNews []models.MaintenanceRoadHistory

	// copier.Copy(&resultMTRHNew, resultMTRH)
	if len(resultMTRH) != 0 {
		var ids []int
		for _, data := range resultMTR {
			var resultMTRHNew models.MaintenanceRoadHistory
			copier.Copy(&resultMTRHNew, &data)
			resultMTRHNew.ID = 0
			resultMTRHNew.MaintenanceID = id
			resultMTRHNew.Revision = data.Revision + 1
			resultMTRHNews = append(resultMTRHNews, resultMTRHNew)
			// resultMTRH[i].ID = 0
			// resultMTRH[i].MaintenanceID = id
			// resultMTRH[i].Revision = data.Revision + 1
			ids = append(ids, data.ID)
		}
		if len(resultMTRHNews) > 0 {
			err = u.Repo.InsertMaintenanceRoadHistory(resultMTRHNews, tx)
			if err != nil {
				u.Repo.RollBack(tx)
				logs.Error(err)
				return nil, nil, responses.NewAppErr(500, err.Error())
			}

			err = u.Repo.UpdateStatusMaintenanceRoadHistory(ids, tx)
			if err != nil {
				u.Repo.RollBack(tx)
				logs.Error(err)
				return nil, nil, responses.NewAppErr(500, err.Error())
			}
		}

	}

	var maintenanceAttachment []models.MaintenanceAttachment
	var idsAttDelete []int
	for _, att := range attReqs {
		if att.Status == "upload" {
			var maintenanceAtt models.MaintenanceAttachment
			copier.Copy(&maintenanceAtt, &att)
			maintenanceAtt.MaintenanceID = id
			maintenanceAtt.CreatedBy = userID
			maintenanceAtt.UpdatedBy = userID
			maintenanceAtt.CreatedAt = time.Now()
			maintenanceAtt.UpdatedAt = time.Now()
			maintenanceAttachment = append(maintenanceAttachment, maintenanceAtt)
		} else if att.Status == "delete" {
			idsAttDelete = append(idsAttDelete, *att.ID)
		}
	}

	if len(idsAttDelete) != 0 {
		err = u.Repo.DeleteMaintenanceAttachment(idsAttDelete, tx)
		if err != nil {
			u.Repo.RollBack(tx)
			logs.Error(err)
			return nil, nil, responses.NewAppErr(500, err.Error())
		}
	}

	err = u.Repo.UpdateMaintenanceAttachment(*maintenanceID, id, tx)
	if err != nil {
		u.Repo.RollBack(tx)
		logs.Error(err)
		return nil, nil, responses.NewAppErr(500, err.Error())
	}

	if len(maintenanceAttachment) != 0 {
		err = u.Repo.InsertMaintenanceaAttachment(maintenanceAttachment, tx)
		if err != nil {
			u.Repo.RollBack(tx)
			logs.Error(err)
			return nil, nil, responses.NewAppErr(500, err.Error())
		}
	}
	u.Repo.Commit(tx)
	return &id, &idParent, nil
}

func (u *Usecase) DeleteMaintenance(IDParent int) error {
	// mAtt, err := u.Repo.GetMaintenanceaAttachmentByMaintenanceID(maintenanceID)
	// if err != nil {
	// 	logs.Error(err)
	// 	return responses.NewAppErr(400, err.Error())
	// }

	// mrhAtt, err := u.Repo.GetMaintenanceaRoadHisAttachmentByMaintenanceID(maintenanceID)
	// if err != nil {
	// 	logs.Error(err)
	// 	return responses.NewAppErr(400, err.Error())
	// }

	maintenance, err := u.Repo.GetMaintenanceByIDParent(IDParent)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return responses.NewAppErr(500, err.Error())
	}

	err = u.Repo.DeleteMaintenance(maintenance.ID)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return responses.NewAppErr(500, err.Error())
	}

	// for _, att := range mAtt {
	// 	if att.Path != "" {
	// 		helpers.DeleteFile(att.Path)
	// 	}
	// }

	// for _, att := range mrhAtt {
	// 	if att.Path != "" {
	// 		helpers.DeleteFile(att.Path)
	// 	}
	// }

	return nil
}

// +++++++++++++++++++++++++++++++++++++++++//
//
//	Maintenance Road
//
// +++++++++++++++++++++++++++++++++++++++++//
func (u *Usecase) GetMaintenanceRoadByID(IDParent int, mRoadId int) (interface{}, error) {

	maintenance, err := u.Repo.GetMaintenanceByIDParent(IDParent)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return nil, responses.NewAppErr(404, err.Error())
	}

	data, err := u.Repo.GetMaintenanceRoadID(maintenance.ID, mRoadId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return nil, responses.NewAppErr(500, err.Error())
	}
	var respond responses.MaintenanceRoadPreload
	copier.Copy(&respond, &data)
	guaranteeExpirationDate := data.Maintenance.GuaranteeExpirationDate.Unix()
	projectEndDate := data.Maintenance.ProjectEndDate.Unix()
	totalDay := guaranteeExpirationDate - projectEndDate
	leftDay := guaranteeExpirationDate - time.Now().Unix()

	var percent int64
	if totalDay != 0 {
		percent = (leftDay * 100) / totalDay
	} else {
		percent = 0
	}

	var colorCode string
	if percent > 20 {
		colorCode = "#1F70F3"
	} else {
		colorCode = "#F1416C"
	}
	respond.Color = colorCode

	// budgetMethod, err := u.Repo.GetSettingBudgetMethodByID(maintenance.BudgetMethodID)
	// if err != nil {
	// 	logs.Error(err)
	// 	return false, responses.NewAppErr(500, err.Error())
	// }

	// respond.IsShowMethod = budgetMethod.IsShowMethod

	geomJSON, err := u.Repo.GetGeomJsonFromMaintenanceRoadID(mRoadId)
	if err != nil {
		return nil, err
	}

	theGeom, err := helpers.ConvertThegeomToGeomJSON(geomJSON)
	if err != nil {
		return nil, err
	}

	respond.TheGeom = theGeom

	if data.RoadLevel == 1 {

		if data.LaneNo == 0 {
			respond.LaneNo = nil
		}

		if data.GridNo == 0 {
			respond.GridNo = nil
		}

		if data.RoadSecNameDes != "" {
			if data.RefDirectionId == 1 {
				respond.RoadName = fmt.Sprintf(`%s - %s`, data.RoadSecNameOr, data.RoadSecNameDes)
			} else if data.RefDirectionId == 2 {
				respond.RoadName = fmt.Sprintf(`%s - %s`, data.RoadSecNameDes, data.RoadSecNameOr)
			}
		} else {
			respond.RoadName = data.RoadSecNameOr
		}
	}

	return &respond, nil
}

func (u *Usecase) CreateMaintenanceRoad(IDParent int, userID int, req requests.MaintenanceRoadsReq) (interface{}, error) {

	maintenance, err := u.Repo.GetMaintenanceByIDParent(IDParent)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	var maintenanceRoad models.MaintenanceRoad
	copier.Copy(&maintenanceRoad, &req)
	maintenanceRoad.UpdatedAt = time.Now().UTC()
	maintenanceRoad.CreatedAt = time.Now().UTC()
	maintenanceRoad.UpdatedBy = userID
	maintenanceRoad.CreatedBy = userID
	maintenanceRoad.MaintenanceID = maintenance.ID
	maintenanceRoad.Status = "A"

	roadGropId, err := u.Repo.GetRoadGroupByRoadId(*req.RoadID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(500, err.Error())
	}
	maintenanceRoad.RoadGroupId = *roadGropId
	theGeom, err := u.CalTheGeom(maintenanceRoad)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(500, err.Error())
	}

	maintenanceRoad.TheGeom = theGeom
	interventionCriteriaParams, err := u.Repo.GetInterventionCriteriaParamsLatest()
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(500, err.Error())
	}

	maintenanceMethod := 0
	maintenanceSurfaceTypeId := 0
	maintenanceSurfaceTypeIdParams := 0
	if req.InterventionCriteriaID != nil && *req.InterventionCriteriaID != 0 {
		interventionCriteria, err := u.Repo.GetInterventionCriteriaByID(*req.InterventionCriteriaID)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(500, err.Error())
		}
		maintenanceMethod = interventionCriteria.MaintenanceMethod
		maintenanceSurfaceTypeId = interventionCriteria.MaintenanceSurfaceTypeId
		maintenanceSurfaceTypeIdParams = interventionCriteria.MaintenanceSurfaceTypeIdParams
	}

	maintenanceRoad.MaintenanceMethodID = maintenanceMethod
	maintenanceRoad.RefSurfaceID = maintenanceSurfaceTypeId
	maintenanceRoad.RefSurfaceParamsID = maintenanceSurfaceTypeIdParams
	maintenanceRoad.InterventionCriteriaIDParams = interventionCriteriaParams.Id
	tx := u.Repo.StartTransSection()
	id, err := u.Repo.CreateMaintenanceRoad(maintenanceRoad, tx)
	if err != nil {
		u.Repo.RollBack(tx)
		logs.Error(err)
		return "", responses.NewAppErr(500, err.Error())
	}

	err = u.Repo.UpdateIDParentMaintenanceRoad(*id, *id, tx)
	if err != nil {
		u.Repo.RollBack(tx)
		logs.Error(err)
		return "", responses.NewAppErr(500, err.Error())
	}

	u.Repo.Commit(tx)
	return id, nil
}

func (u *Usecase) CheckValidateIsMethod(IDParent int) (bool, error) {

	maintenance, err := u.Repo.GetMaintenanceByIDParent(IDParent)
	if err != nil {
		logs.Error(err)
		return false, responses.NewAppErr(500, err.Error())
	}

	budgetMethod, err := u.Repo.GetSettingBudgetMethodByID(maintenance.BudgetMethodID)
	if err != nil {
		logs.Error(err)
		return false, responses.NewAppErr(500, err.Error())
	}
	if budgetMethod.IsShowMethod {
		return true, nil
	}

	return false, nil
}

func (u *Usecase) UpdateMaintenanceRoad(IDParent int, mRoadId int, userID int, req requests.MaintenanceRoadsReq) (*int, *int, error) {

	maintenance, err := u.Repo.GetMaintenanceByIDParent(IDParent)
	if err != nil {
		logs.Error(err)
		return nil, nil, responses.NewAppErr(500, err.Error())
	}

	revision, idParent, err := u.Repo.GetMaxRevisionMaintenanceRoad(mRoadId)
	if err != nil {
		logs.Error(err)
		return nil, nil, responses.NewAppErr(500, err.Error())
	}
	var maintenanceRoad models.MaintenanceRoad

	copier.Copy(&maintenanceRoad, &req)
	maintenanceRoad.UpdatedAt = time.Now().UTC()
	maintenanceRoad.CreatedAt = time.Now().UTC()
	maintenanceRoad.UpdatedBy = userID
	maintenanceRoad.CreatedBy = userID
	maintenanceRoad.MaintenanceID = maintenance.ID
	maintenanceRoad.Status = "A"
	maintenanceRoad.IDParent = *idParent
	maintenanceRoad.Revision = *revision + 1

	roadGropId, err := u.Repo.GetRoadGroupByRoadId(*req.RoadID)
	if err != nil {
		logs.Error(err)
		return nil, nil, responses.NewAppErr(500, err.Error())
	}
	maintenanceRoad.RoadGroupId = *roadGropId
	theGeom, err := u.CalTheGeom(maintenanceRoad)
	if err != nil {
		logs.Error(err)
		return nil, nil, responses.NewAppErr(500, err.Error())
	}

	maintenanceRoad.TheGeom = theGeom
	interventionCriteriaParams, err := u.Repo.GetInterventionCriteriaParamsLatest()
	if err != nil {
		logs.Error(err)
		return nil, nil, responses.NewAppErr(500, err.Error())
	}

	maintenanceMethodID := 0
	maintenanceSurfaceTypeId := 0
	maintenanceSurfaceTypeIdParams := 0
	if req.InterventionCriteriaID != nil && *req.InterventionCriteriaID != 0 {
		interventionCriteria, err := u.Repo.GetInterventionCriteriaByID(*req.InterventionCriteriaID)
		if err != nil {
			logs.Error(err)
			return nil, nil, responses.NewAppErr(500, err.Error())
		}
		maintenanceMethodID = interventionCriteria.MaintenanceMethod
		maintenanceSurfaceTypeId = interventionCriteria.MaintenanceSurfaceTypeId
		maintenanceSurfaceTypeIdParams = interventionCriteria.MaintenanceSurfaceTypeIdParams
	}

	maintenanceRoad.MaintenanceMethodID = maintenanceMethodID
	maintenanceRoad.RefSurfaceID = maintenanceSurfaceTypeId
	maintenanceRoad.RefSurfaceParamsID = maintenanceSurfaceTypeIdParams
	maintenanceRoad.InterventionCriteriaIDParams = interventionCriteriaParams.Id
	tx := u.Repo.StartTransSection()

	id, err := u.Repo.CreateMaintenanceRoad(maintenanceRoad, tx)
	if err != nil {
		u.Repo.RollBack(tx)
		logs.Error(err)
		return nil, nil, responses.NewAppErr(500, err.Error())
	}

	err = u.Repo.UpdateStatusMaintenanceRoad([]int{mRoadId}, tx)
	if err != nil {
		u.Repo.RollBack(tx)
		logs.Error(err)
		return nil, nil, responses.NewAppErr(500, err.Error())
	}

	u.Repo.Commit(tx)
	return id, idParent, nil
}

func (u *Usecase) DeleteMaintenanceRoad(maintenanceID int, mRoadId int) error {
	err := u.Repo.DeleteMaintenanceRoad(maintenanceID, mRoadId)
	if err != nil {
		return responses.NewAppErr(404, err.Error())
	}
	return nil
}

func (u *Usecase) CheckOverlap(data []models.MaintenanceRoad) error {
	dict := make(map[string][]models.MaintenanceRoad)
	for _, v := range data {
		key := fmt.Sprintf("%d_%d", v.RoadID, v.LaneNo)
		dict[key] = append(dict[key], v)
	}
	for _, value := range dict {
		for i, v := range value {
			direction, err := u.Repo.FindRefDirectionIDByRoadID(v.RoadID)
			if err != nil {
				logs.Error(err)
				return responses.NewAppErr(400, err.Error())
			}
			for j := i + 1; j < len(value); j++ {
				if direction == 1 {
					if value[j].KmStart >= value[i].KmStart && value[j].KmStart < value[i].KmEnd {
						err := errors.New(constants.OVERLAPPING_RANGE)
						logs.Error(err)
						return responses.NewAppErr(400, err.Error())
					}
				}
				if direction == 2 {
					if value[j].KmStart <= value[i].KmStart && value[j].KmStart > value[i].KmEnd {
						err := errors.New(constants.OVERLAPPING_RANGE)
						logs.Error(err)
						return responses.NewAppErr(400, err.Error())
					}
				}
			}
		}

	}

	return nil
}

func (u *Usecase) CalTheGeom(data models.MaintenanceRoad) (string, error) {
	var subRoadStart float64
	var subRoadEnd float64
	// var subRoadMin float64
	// var subRoadMax float64
	var theGeom string
	direction, err := u.Repo.FindRefDirectionIDByRoadID(data.RoadID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	fmt.Println("direction", direction)
	// fullGeom, err := u.Repo.GetGeomByRoadID(data.RoadID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	roadGeom, err := u.Repo.GetGeomByRoadGeomID(data.RoadID, data.LaneNo)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	switch direction {
	case 1:
		if data.KmStart >= float64(roadGeom.KmStart) {
			if data.KmStart >= float64(roadGeom.KmEnd) {
				subRoadStart = 0
			} else {
				subRoadStart = (math.Abs((data.KmStart - float64(roadGeom.KmStart)))) / (math.Abs((float64(roadGeom.KmStart) - float64(roadGeom.KmEnd))))
			}
		} else {
			subRoadStart = 0
		}
		if data.KmEnd < float64(roadGeom.KmEnd) {
			if data.KmEnd < float64(roadGeom.KmStart) {
				subRoadEnd = 0
			} else {
				// subRoadEnd = float64(math.Abs(float64(data.KmEnd[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
				subRoadEnd = math.Abs(data.KmEnd-float64(roadGeom.KmStart)) / math.Abs(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))
			}
		} else {
			subRoadEnd = 1
		}

	case 2:
		if data.KmStart <= float64(roadGeom.KmStart) {
			if data.KmStart < float64(roadGeom.KmEnd) {
				subRoadStart = 0
			} else {
				// subRoadStart = float64(math.Abs(float64(data.KmStart[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
				subRoadStart = math.Abs(data.KmStart-float64(roadGeom.KmStart)) / math.Abs(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))
			}
		} else {
			subRoadStart = 1
		}

		if data.KmEnd >= float64(roadGeom.KmEnd) {
			if data.KmEnd >= float64(roadGeom.KmStart) {
				subRoadEnd = 0
			} else {
				// subRoadEnd = float64(math.Abs(float64(data.KmEnd[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
				subRoadEnd = math.Abs(data.KmEnd-float64(roadGeom.KmStart)) / math.Abs(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))
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

func (u *Usecase) CalTheGeomHistory(data models.MaintenanceRoadHistory) (string, error) {
	var subRoadStart float64
	var subRoadEnd float64
	// var subRoadMin float64
	// var subRoadMax float64
	var theGeom string
	direction, err := u.Repo.FindRefDirectionIDByRoadID(data.RoadID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	roadGeom, err := u.Repo.GetGeomByRoadGeomID(data.RoadID, data.LaneNo)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	switch direction {
	case 1:
		if data.KmStart >= float64(roadGeom.KmStart) {
			if data.KmStart >= float64(roadGeom.KmEnd) {
				subRoadStart = 0
			} else {
				subRoadStart = (math.Abs((data.KmStart - float64(roadGeom.KmStart)))) / (math.Abs((float64(roadGeom.KmStart) - float64(roadGeom.KmEnd))))
			}
		} else {
			subRoadStart = 0
		}
		if data.KmEnd < float64(roadGeom.KmEnd) {
			if data.KmEnd < float64(roadGeom.KmStart) {
				subRoadEnd = 0
			} else {
				// subRoadEnd = float64(math.Abs(float64(data.KmEnd[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
				subRoadEnd = math.Abs(data.KmEnd-float64(roadGeom.KmStart)) / math.Abs(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))
			}
		} else {
			subRoadEnd = 1
		}

	case 2:
		if data.KmStart <= float64(roadGeom.KmStart) {
			if data.KmStart < float64(roadGeom.KmEnd) {
				subRoadStart = 0
			} else {
				// subRoadStart = float64(math.Abs(float64(data.KmStart[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
				subRoadStart = math.Abs(data.KmStart-float64(roadGeom.KmStart)) / math.Abs(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))
			}
		} else {
			subRoadStart = 1
		}

		if data.KmEnd >= float64(roadGeom.KmEnd) {
			if data.KmEnd >= float64(roadGeom.KmStart) {
				subRoadEnd = 0
			} else {
				// subRoadEnd = float64(math.Abs(float64(data.KmEnd[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
				subRoadEnd = math.Abs(data.KmEnd-float64(roadGeom.KmStart)) / math.Abs(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))
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

func (u *Usecase) CalTheGeomHistoryRoad(data models.MaintenanceRoadData) (string, error) {
	var subRoadStart float64
	var subRoadEnd float64
	// var subRoadMin float64
	// var subRoadMax float64
	var theGeom string
	direction, err := u.Repo.FindRefDirectionIDByRoadID(data.RoadID)
	if err != nil {
		logs.Error(err)
		return "", err
	}
	fmt.Println("direction", direction)
	// fullGeom, err := u.Repo.GetGeomByRoadID(data.RoadID)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	roadGeom, err := u.Repo.GetGeomByRoadGeomID(data.RoadID, data.Lane)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	switch direction {
	case 1:
		if data.KmStart >= float64(roadGeom.KmStart) {
			if data.KmStart >= float64(roadGeom.KmEnd) {
				subRoadStart = 0
			} else {
				subRoadStart = (math.Abs((data.KmStart - float64(roadGeom.KmStart)))) / (math.Abs((float64(roadGeom.KmStart) - float64(roadGeom.KmEnd))))
			}
		} else {
			subRoadStart = 0
		}
		if data.KmEnd < float64(roadGeom.KmEnd) {
			if data.KmEnd < float64(roadGeom.KmStart) {
				subRoadEnd = 0
			} else {
				// subRoadEnd = float64(math.Abs(float64(data.KmEnd[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
				subRoadEnd = math.Abs(data.KmEnd-float64(roadGeom.KmStart)) / math.Abs(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))
			}
		} else {
			subRoadEnd = 1
		}

	case 2:
		if data.KmStart <= float64(roadGeom.KmStart) {
			if data.KmStart < float64(roadGeom.KmEnd) {
				subRoadStart = 0
			} else {
				// subRoadStart = float64(math.Abs(float64(data.KmStart[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
				subRoadStart = math.Abs(data.KmStart-float64(roadGeom.KmStart)) / math.Abs(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))
			}
		} else {
			subRoadStart = 1
		}

		if data.KmEnd >= float64(roadGeom.KmEnd) {
			if data.KmEnd >= float64(roadGeom.KmStart) {
				subRoadEnd = 0
			} else {
				// subRoadEnd = float64(math.Abs(float64(data.KmEnd[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
				subRoadEnd = math.Abs(data.KmEnd-float64(roadGeom.KmStart)) / math.Abs(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))
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

func (u *Usecase) CalTheGeomHistoryRoadHis(data models.MaintenanceRoadHistoryData) (string, error) {
	var subRoadStart float64
	var subRoadEnd float64
	// var subRoadMin float64
	// var subRoadMax float64
	var theGeom string
	direction, err := u.Repo.FindRefDirectionIDByRoadID(data.RoadID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	fmt.Println("direction", direction)
	// fullGeom, err := u.Repo.GetGeomByRoadID(data.RoadID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	roadGeom, err := u.Repo.GetGeomByRoadGeomID(data.RoadID, data.Lane)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	switch direction {
	case 1:
		if data.KmStart >= float64(roadGeom.KmStart) {
			if data.KmStart >= float64(roadGeom.KmEnd) {
				subRoadStart = 0
			} else {
				subRoadStart = (math.Abs((data.KmStart - float64(roadGeom.KmStart)))) / (math.Abs((float64(roadGeom.KmStart) - float64(roadGeom.KmEnd))))
			}
		} else {
			subRoadStart = 0
		}
		if data.KmEnd < float64(roadGeom.KmEnd) {
			if data.KmEnd < float64(roadGeom.KmStart) {
				subRoadEnd = 0
			} else {
				// subRoadEnd = float64(math.Abs(float64(data.KmEnd[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
				subRoadEnd = math.Abs(data.KmEnd-float64(roadGeom.KmStart)) / math.Abs(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))
			}
		} else {
			subRoadEnd = 1
		}

	case 2:
		if data.KmStart <= float64(roadGeom.KmStart) {
			if data.KmStart < float64(roadGeom.KmEnd) {
				subRoadStart = 0
			} else {
				// subRoadStart = float64(math.Abs(float64(data.KmStart[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
				subRoadStart = math.Abs(data.KmStart-float64(roadGeom.KmStart)) / math.Abs(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))
			}
		} else {
			subRoadStart = 1
		}

		if data.KmEnd >= float64(roadGeom.KmEnd) {
			if data.KmEnd >= float64(roadGeom.KmStart) {
				subRoadEnd = 0
			} else {
				// subRoadEnd = float64(math.Abs(float64(data.KmEnd[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
				subRoadEnd = math.Abs(data.KmEnd-float64(roadGeom.KmStart)) / math.Abs(float64(roadGeom.KmStart)-float64(roadGeom.KmEnd))
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

// func (u *Usecase) MaintenanceFinished(maintenanceID int, req requests.MaintenanceFinished) error {
// 	progress, err := u.Repo.MaintenanceCheckProgressComplete(maintenanceID)
// 	if err != nil {
// 		return responses.NewAppErr(400, err.Error())
// 	}
// 	if progress != 100 {
// 		// ปิดชั่วคราว
// 		// return responses.NewAppErr(400, "ไม่สามารถสิ้นสุดโครงการได้ เนื่องจากรายงานความก้าวหน้าไม่ครบ 100%")
// 	}
// 	err = u.Repo.MaintenanceFinished(maintenanceID, req)
// 	if err != nil {
// 		logs.Error(err)
// 		if err == gorm.ErrRecordNotFound {
// 			return nil
// 		}
// 		return responses.NewAppErr(400, err.Error())
// 	}
// 	return nil
// }

func (u *Usecase) CheckMaintenanceDuplicate(idParent int, name string) (bool, error) {

	isDup, err := u.Repo.CheckMaintenanceDuplicate(idParent, name)
	if err != nil {
		logs.Error(err)
		return true, responses.NewAppErr(500, err.Error())
	}

	return isDup, nil
}

// func (u *Usecase) GetMaintenanceStatus(maintenanceID int) (interface{}, error) {
// 	data, err := u.Repo.GetMaintenanceStatus(maintenanceID)
// 	if err != nil {
// 		logs.Error(err)
// 		if err == gorm.ErrRecordNotFound {
// 			return responses.NoData{}, nil
// 		}
// 		return data, responses.NewAppErr(400, err.Error())
// 	}

// 	progressData := make(map[string]models.MaintenancePlanDetailProgress)
// 	progress, err := u.Repo.GetMaintenancePlanProgress(maintenanceID)
// 	if err != nil {
// 		logs.Error(err)
// 		if err == gorm.ErrRecordNotFound {
// 			return responses.NoData{}, nil
// 		}
// 		return data, responses.NewAppErr(400, err.Error())
// 	}

// 	for _, item := range progress {
// 		progressData[item.Schedule.Format("2006-01")] = item
// 		// fmt.Println()
// 	}

// 	// fmt.Println(t.Format("2006-01"))
// 	// return progressData, nil
// 	scheduleArr := make(map[string]string)
// 	var maintenanceStatus []responses.MaintenanceStatus
// 	for _, item := range data {
// 		schedules := item.SchedulePlan
// 		var maintenanceStatusSchedules []responses.MaintenanceStatusSchedules
// 		periodOfWork := 0
// 		for _, schedule := range schedules {
// 			val, isHave := progressData[schedule.Schedule]
// 			if isHave {
// 				status := ""
// 				if schedule.Status.Name == "งวดงาน" {
// 					periodOfWork++
// 					status = fmt.Sprintf("งวดที่ %v", periodOfWork)
// 				} else {
// 					status = schedule.Status.Name
// 				}

// 				maintenanceStatusSchedules = append(maintenanceStatusSchedules, responses.MaintenanceStatusSchedules{IsChecked: val.IsSelect, Schedule: val.Schedule.Format("2006-01"), Status: status, DisbursementDate: val.DisbursementDate})
// 				scheduleArr[schedule.Schedule] = schedule.Schedule
// 			}
// 		}
// 		// for _, item2 := range progress {
// 		// 	_, isHave := scheduleArr[item2.Schedule.Format("2006-01")]
// 		// 	if !isHave {
// 		// 		disbursementDate := ""
// 		// 		if item2.DisbursementDate != nil {
// 		// 			disbursementDate = item2.DisbursementDate.Format("2006-01-02")
// 		// 		} else {
// 		// 			disbursementDate = ""
// 		// 		}

// 		// 		fmt.Println(disbursementDate)
// 		// 		maintenanceStatusSchedules = append(maintenanceStatusSchedules, responses.MaintenanceStatusSchedules{IsChecked: false, Schedule: disbursementDate, Status: ""})
// 		// 	}
// 		// }

// 		maintenanceStatus = append(maintenanceStatus, responses.MaintenanceStatus{IsCurrent: item.IsCurrent, Name: item.Name, Schedules: maintenanceStatusSchedules})
// 		// helpers.PrintlnJson(maintenanceStatusSchedules)
// 	}

// 	// for _, item := range maintenanceStatus {
// 	// 	// if item. == "" {
// 	// 	// 	continue
// 	// 	// }
// 	// }

// 	if len(maintenanceStatus) > 0 {
// 		return maintenanceStatus, nil
// 	} else {
// 		return []string{}, nil
// 	}
// }

// =======================================================================================================
//
//	Maintenance maintenance_history
//
// =======================================================================================================
func (u *Usecase) GetMaintenanceHistoryByID(idParent int, mRoadHisId int) (interface{}, error) {

	maintenance, err := u.Repo.GetMaintenanceByIDParent(idParent)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewAppErr(500, err.Error())
	}

	data, err := u.Repo.GetMaintenanceHistoryByID(maintenance.ID, mRoadHisId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responses.NoData{}, nil
		}
		logs.Error(err)
		return nil, responses.NewAppErr(500, err.Error())
	}

	var respond responses.MaintenanceRoadHistoryPreload
	copier.Copy(&respond, &data)
	guaranteeExpirationDate := data.Maintenance.GuaranteeExpirationDate.Unix()
	projectEndDate := data.Maintenance.ProjectEndDate.Unix()
	totalDay := guaranteeExpirationDate - projectEndDate
	leftDay := guaranteeExpirationDate - time.Now().Unix()

	var percent int64
	if totalDay != 0 {
		percent = (leftDay * 100) / totalDay
	} else {
		percent = 0
	}

	var colorCode string
	if percent > 20 {
		colorCode = "#1F70F3"
	} else {
		colorCode = "#F1416C"
	}
	respond.Color = colorCode

	geomJSON, err := u.Repo.GetGeomJsonFromMaintenanceRoadHistoryID(mRoadHisId)
	if err != nil {
		return nil, err
	}

	theGeom, err := helpers.ConvertThegeomToGeomJSON(geomJSON)
	if err != nil {
		return nil, err
	}

	respond.TheGeom = theGeom

	if data.LaneNo == 0 {
		respond.LaneNo = nil
	}

	if data.GridNo == 0 {
		respond.GridNo = nil
	}

	// budgetMethod, err := u.Repo.GetSettingBudgetMethodByID(maintenance.BudgetMethodID)
	// if err != nil {
	// 	logs.Error(err)
	// 	return false, responses.NewAppErr(500, err.Error())
	// }

	//respond.IsShowMethod = budgetMethod.IsShowMethod

	if data.RoadLevel == 1 {

		if data.RoadSecNameDes != "" {
			if data.RefDirectionId == 1 {
				respond.RoadName = fmt.Sprintf(`%s - %s`, data.RoadSecNameOr, data.RoadSecNameDes)
			} else if data.RefDirectionId == 2 {
				respond.RoadName = fmt.Sprintf(`%s - %s`, data.RoadSecNameDes, data.RoadSecNameOr)
			}
		} else {
			respond.RoadName = data.RoadSecNameOr
		}
	}

	hostPath := os.Getenv("STORAGE_IP") + "/"

	for i, att := range data.Attachments {
		respond.Attachments[i].Path = hostPath + att.Path
	}

	return respond, nil
}

func (u *Usecase) GetMaintenanceHistory(idParent int, maintenanceID int) (interface{}, error) {
	data, err := u.Repo.GetMaintenanceHistory(maintenanceID)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return "", responses.NewAppErr(400, err.Error())
	}
	return data, nil
}

func (u *Usecase) CreateMaintenanceHistory(IDParent int, userID int, req requests.MaintenanceRoadHistoryReq) (interface{}, error) {

	maintenance, err := u.Repo.GetMaintenanceByIDParent(IDParent)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	var maintenanceRoadHistory models.MaintenanceRoadHistory
	copier.Copy(&maintenanceRoadHistory, &req)
	maintenanceRoadHistory.UpdatedAt = time.Now().UTC()
	maintenanceRoadHistory.CreatedAt = time.Now().UTC()
	maintenanceRoadHistory.UpdatedBy = userID
	maintenanceRoadHistory.CreatedBy = userID
	maintenanceRoadHistory.MaintenanceID = maintenance.ID
	maintenanceRoadHistory.Status = "A"

	roadGropId, err := u.Repo.GetRoadGroupByRoadId(*req.RoadID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	maintenanceRoadHistory.RoadGroupId = *roadGropId

	theGeom, err := u.CalTheGeomHistory(maintenanceRoadHistory)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	maintenanceRoadHistory.TheGeom = theGeom
	interventionCriteriaParams, err := u.Repo.GetInterventionCriteriaParamsLatest()
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	maintenanceMethodID := 0
	maintenanceSurfaceTypeId := 0
	maintenanceSurfaceTypeIdParams := 0
	if req.InterventionCriteriaID != nil && *req.InterventionCriteriaID != 0 {
		interventionCriteria, err := u.Repo.GetInterventionCriteriaByID(*req.InterventionCriteriaID)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		maintenanceMethodID = interventionCriteria.MaintenanceMethod
		maintenanceSurfaceTypeId = interventionCriteria.MaintenanceSurfaceTypeId
		maintenanceSurfaceTypeIdParams = interventionCriteria.MaintenanceSurfaceTypeIdParams
	}

	maintenanceRoadHistory.MaintenanceMethodID = maintenanceMethodID
	maintenanceRoadHistory.RefSurfaceID = maintenanceSurfaceTypeId
	maintenanceRoadHistory.RefSurfaceParamsID = maintenanceSurfaceTypeIdParams
	maintenanceRoadHistory.InterventionCriteriaIDParams = interventionCriteriaParams.Id
	tx := u.Repo.StartTransSection()

	id, err := u.Repo.CreateMaintenanceRoadHistory(maintenanceRoadHistory, tx)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	err = u.Repo.UpdateIDParentMaintenanceRoadHistory(*id, *id, tx)
	if err != nil {
		u.Repo.RollBack(tx)
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	// var maintenanceAttachment []models.MaintenanceRoadHistoryAttachment
	// for _, att := range attReqs {
	// 	if att.Status == "upload" {
	// 		var maintenanceRoadHisAtt models.MaintenanceRoadHistoryAttachment
	// 		copier.Copy(&maintenanceRoadHisAtt, &att)
	// 		maintenanceRoadHisAtt.MaintenanceID = maintenanceID
	// 		maintenanceRoadHisAtt.MaintenanceRoadHistoryID = *id
	// 		maintenanceRoadHisAtt.CreatedBy = userID
	// 		maintenanceRoadHisAtt.UpdatedBy = userID
	// 		maintenanceRoadHisAtt.CreatedAt = time.Now()
	// 		maintenanceRoadHisAtt.UpdatedAt = time.Now()
	// 		maintenanceAttachment = append(maintenanceAttachment, maintenanceRoadHisAtt)
	// 	}
	// }

	// err = u.Repo.InsertMaintenanceaRoadHisAttachment(maintenanceAttachment, tx)
	// if err != nil {
	// 	u.Repo.RollBack(tx)
	// 	logs.Error(err)
	// 	return "", responses.NewAppErr(400, err.Error())
	// }

	u.Repo.Commit(tx)
	return id, nil
}

func (u *Usecase) UpdateMaintenanceHistory(IDParent int, historyID int, userID int, req requests.MaintenanceRoadHistoryReq) (*int, *int, error) {

	revision, idParent, err := u.Repo.GetMaxRevisionMaintenanceRoadHistory(historyID)
	if err != nil {
		logs.Error(err)
		return nil, nil, responses.NewAppErr(400, err.Error())
	}

	maintenance, err := u.Repo.GetMaintenanceByIDParent(IDParent)
	if err != nil {
		logs.Error(err)
		return nil, nil, responses.NewAppErr(400, err.Error())
	}

	var maintenanceRoadHistory models.MaintenanceRoadHistory
	copier.Copy(&maintenanceRoadHistory, &req)
	maintenanceRoadHistory.UpdatedAt = time.Now().UTC()
	maintenanceRoadHistory.CreatedAt = time.Now().UTC()
	maintenanceRoadHistory.UpdatedBy = userID
	maintenanceRoadHistory.CreatedBy = userID
	maintenanceRoadHistory.MaintenanceID = maintenance.ID
	maintenanceRoadHistory.Status = "A"
	maintenanceRoadHistory.Revision = *revision + 1
	maintenanceRoadHistory.IDParent = *idParent
	theGeom, err := u.CalTheGeomHistory(maintenanceRoadHistory)
	if err != nil {
		logs.Error(err)
		return nil, nil, responses.NewAppErr(400, err.Error())
	}
	maintenanceRoadHistory.TheGeom = theGeom

	interventionCriteriaParams, err := u.Repo.GetInterventionCriteriaParamsLatest()
	if err != nil {
		logs.Error(err)
		return nil, nil, responses.NewAppErr(400, err.Error())
	}

	roadGropId, err := u.Repo.GetRoadGroupByRoadId(*req.RoadID)
	if err != nil {
		logs.Error(err)
		return nil, nil, responses.NewAppErr(400, err.Error())
	}
	maintenanceRoadHistory.RoadGroupId = *roadGropId

	maintenanceMethodID := 0
	maintenanceSurfaceTypeId := 0
	maintenanceSurfaceTypeIdParams := 0

	if req.InterventionCriteriaID != nil && *req.InterventionCriteriaID != 0 {
		interventionCriteria, err := u.Repo.GetInterventionCriteriaByID(*req.InterventionCriteriaID)
		if err != nil {
			logs.Error(err)
			return nil, nil, responses.NewAppErr(400, err.Error())
		}

		maintenanceRoadHistory.MaintenanceMethodID = maintenanceMethodID
		maintenanceSurfaceTypeId = interventionCriteria.MaintenanceSurfaceTypeId
		maintenanceSurfaceTypeIdParams = interventionCriteria.MaintenanceSurfaceTypeIdParams
	}

	maintenanceRoadHistory.RefSurfaceID = maintenanceSurfaceTypeId
	maintenanceRoadHistory.RefSurfaceParamsID = maintenanceSurfaceTypeIdParams
	maintenanceRoadHistory.InterventionCriteriaIDParams = interventionCriteriaParams.Id

	tx := u.Repo.StartTransSection()

	id, err := u.Repo.CreateMaintenanceRoadHistory(maintenanceRoadHistory, tx)
	if err != nil {
		logs.Error(err)
		return nil, nil, responses.NewAppErr(400, err.Error())
	}

	err = u.Repo.UpdateStatusMaintenanceRoadHistory([]int{historyID}, tx)
	if err != nil {
		u.Repo.RollBack(tx)
		logs.Error(err)
		return nil, nil, responses.NewAppErr(400, err.Error())
	}

	// var maintenanceAttachment []models.MaintenanceRoadHistoryAttachment
	// var idsAttDelete []int
	// for _, att := range attReqs {
	// 	if att.Status == "upload" {
	// 		var maintenanceRoadHisAtt models.MaintenanceRoadHistoryAttachment
	// 		copier.Copy(&maintenanceRoadHisAtt, &att)
	// 		maintenanceRoadHisAtt.MaintenanceID = maintenanceID
	// 		maintenanceRoadHisAtt.MaintenanceRoadHistoryID = *id
	// 		maintenanceRoadHisAtt.CreatedBy = userID
	// 		maintenanceRoadHisAtt.UpdatedBy = userID
	// 		maintenanceRoadHisAtt.CreatedAt = time.Now()
	// 		maintenanceRoadHisAtt.UpdatedAt = time.Now()
	// 		maintenanceAttachment = append(maintenanceAttachment, maintenanceRoadHisAtt)
	// 	} else if att.Status == "delete" {
	// 		idsAttDelete = append(idsAttDelete, *att.ID)
	// 	}
	// }

	// if len(idsAttDelete) != 0 {
	// 	err = u.Repo.DeleteMaintenanceRoadHisAttachment(idsAttDelete, tx)
	// 	if err != nil {
	// 		u.Repo.RollBack(tx)
	// 		logs.Error(err)
	// 		return nil, nil, responses.NewAppErr(400, err.Error())
	// 	}
	// }

	err = u.Repo.UpdateMaintenanceRoadHisAttachment(historyID, *id, tx)
	if err != nil {
		u.Repo.RollBack(tx)
		logs.Error(err)
		return nil, nil, responses.NewAppErr(400, err.Error())
	}

	// if len(maintenanceAttachment) != 0 {
	// 	err = u.Repo.InsertMaintenanceaRoadHisAttachment(maintenanceAttachment, tx)
	// 	if err != nil {
	// 		u.Repo.RollBack(tx)
	// 		logs.Error(err)
	// 		return nil, nil, responses.NewAppErr(400, err.Error())
	// 	}
	// }

	u.Repo.Commit(tx)
	return id, idParent, nil
}

func (u *Usecase) DeleteMaintenanceHistory(maintenanceID int, hisID int) (interface{}, error) {
	mRHAtt, err := u.Repo.GetMaintenanceaRoadHisAttachmentByMRoadHisID(hisID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	err = u.Repo.DeleteMaintenanceHistory(maintenanceID, hisID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	for _, att := range mRHAtt {
		if att.Path != "" {
			helpers.DeleteFile(att.Path)
		}
	}

	return "", nil
}

func (u *Usecase) GetMaintenanceListHistory(prams requests.MaintenancePrams) (interface{}, error) {
	var responds []responses.MaintenanceList
	queryFilter := createFilter(prams)
	result, err := u.Repo.GetMaintenanceListHistory(queryFilter)
	if err != nil {
		logs.Error(err)
		return result, err
	}
	if len(result) == 0 {
		return []responses.MaintenanceList{}, nil
	}

	for _, item := range result {
		var respond responses.MaintenanceList
		copier.Copy(&respond, &item)
		// roadGroup, err := u.Repo.GetRoadGroupInfoByMaintenanceID(item.ID)
		// if err != nil {
		// 	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 		continue
		// 	}
		// 	logs.Error(err)
		// 	return result, err
		// }
		//respond.RoadGroup = roadGroup
		respond.KmTotal, err = u.Repo.GetTotalKmByMaintenanceID(item.ID)
		if err != nil {
			logs.Error(err)
			return result, err
		}
		// if prams.RoadGroupID != nil {
		// 	if roadGroup.Id != *prams.RoadGroupID {
		// 		continue
		// 	}
		// }
		responds = append(responds, respond)
	}

	return responds, nil
}

// func (u *Usecase) GetMaintenanceYearHistory(maintenanceID int) (interface{}, error) {
// 	data, err := u.Repo.GetMaintenanceYearHistory(maintenanceID)
// 	if err != nil {
// 		logs.Error(err)
// 		return "", err
// 	}

// 	year := []int{}
// 	for _, item := range data {
// 		year = append(year, item.BudgetYear)
// 	}
// 	return year, nil
// }

func (u *Usecase) CreateMaintenancePlanReport(dataCharts interface{}, dataTable interface{}, maintenanceID int) (string, error) {
	ctx, cancel := helpers.NewChromedpContext(context.Background(), log.Printf)
	defer cancel()

	dataTheGeom, err := u.Repo.GetMaintenanceRoadTheGeomByID(maintenanceID)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	dataChartByte, err := json.Marshal(dataCharts)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	var maintainChartReport []models.MaintainChartReport
	err = json.Unmarshal(dataChartByte, &maintainChartReport)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	dataTableByte, err := json.Marshal(dataTable)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	var maintainTableReport models.MaintainTableReport
	err = json.Unmarshal(dataTableByte, &maintainTableReport)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	b, err := json.Marshal(maintainChartReport)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	dataMaintenanceListByID, err := u.Repo.GetMaintenanceListByIDWithNotFilterIsComplete(maintenanceID)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	countProblem := 0
	for _, item := range maintainTableReport.Problems {
		countProblem = countProblem + len(item)
	}

	var maintainTableReportlists []models.MaintainTableReportValue
	for index := range maintainTableReport.MaintenancePlan {
		var maintainTableReportlist models.MaintainTableReportValue
		maintainTableReportlist.PlanName = maintainTableReport.MaintenancePlan[len(maintainTableReport.MaintenancePlan)-1-index].PlanName
		maintainTableReportlist.Value = maintainTableReport.MaintenancePlan[len(maintainTableReport.MaintenancePlan)-1-index].Value
		maintainTableReportlists = append(maintainTableReportlists, maintainTableReportlist)
	}

	maintainTableReport.MaintenancePlan = maintainTableReportlists

	var maintainReport models.MaintainReport
	maintainReport.MaintainChartReport = string(b)
	if countProblem > 0 {
		maintainReport.IsProblems = true
	} else {
		maintainReport.IsProblems = false
	}

	maintainReport.MaintainTableReport = maintainTableReport
	maintainReport.MaintainRoadPath = dataTheGeom.TheGeom
	maintainReport.Date = time.Now().Format("02/01/2006")
	maintainReport.Maintain = dataMaintenanceListByID
	maintainReport.ApiKey = os.Getenv("LONGDO_API_KEY")

	filePath := os.Getenv("MAINTENANCE_PROCESS_PDF")

	templateName := os.Getenv("TEMPLATE_MAINTENANCE_PLAN")
	result, err := helpers.InitDataToHtml(templateName, maintainReport, filePath)
	if err != nil {
		logs.Error(err)
		return "", err
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
	if err := chromedp.Run(ctx, helpers.PrintToPDF(string(html), &buf, true)); err != nil {
		logs.Error(err)
		return "", err
	}

	unix := strconv.Itoa(int(time.Now().Unix()))

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return "", err
	}
	fmt.Println()
	fullFilePath := filePath + unix + ".pdf"
	if err := os.WriteFile(fullFilePath, buf, 0777); err != nil {
		logs.Error(err)
		return "", err
	}

	return os.Getenv("STORAGE_IP") + "/" + fullFilePath, nil
}

func (u *Usecase) CreateMaintenanceHistoryPlanReport(dataCharts interface{}, dataTable interface{}, history interface{}, maintenanceID int) (string, error) {
	ctx, cancel := helpers.NewChromedpContext(context.Background(), log.Printf)
	defer cancel()

	dataTheGeom, err := u.Repo.GetMaintenanceRoadTheGeomByID(maintenanceID)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	dataChartByte, err := json.Marshal(dataCharts)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	var maintainChartReport []models.MaintainChartReport
	err = json.Unmarshal(dataChartByte, &maintainChartReport)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	dataTableByte, err := json.Marshal(dataTable)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	var maintainTableReport models.MaintainTableReport
	err = json.Unmarshal(dataTableByte, &maintainTableReport)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	b, err := json.Marshal(maintainChartReport)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	dataMaintenanceListByID, err := u.Repo.GetMaintenanceListByIDWithNotFilterIsComplete(maintenanceID)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	countProblem := 0
	for _, item := range maintainTableReport.Problems {
		countProblem = countProblem + len(item)
	}

	var maintainTableReportlists []models.MaintainTableReportValue

	for index := range maintainTableReport.MaintenancePlan {
		var maintainTableReportlist models.MaintainTableReportValue
		maintainTableReportlist.PlanName = maintainTableReport.MaintenancePlan[len(maintainTableReport.MaintenancePlan)-1-index].PlanName
		maintainTableReportlist.Value = maintainTableReport.MaintenancePlan[len(maintainTableReport.MaintenancePlan)-1-index].Value
		maintainTableReportlists = append(maintainTableReportlists, maintainTableReportlist)
	}

	maintainTableReport.MaintenancePlan = maintainTableReportlists

	var maintainPreload models.MaintainPreload
	dataTableMaintainPreload, err := json.Marshal(history)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	err = json.Unmarshal([]byte(string(dataTableMaintainPreload)), &maintainPreload)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	var maintainReport models.MaintainReport
	maintainReport.MaintainChartReport = string(b)
	if countProblem > 0 {
		maintainReport.IsProblems = true
	} else {
		maintainReport.IsProblems = false
	}

	// if len(maintainPreload.MaintenanceRoadHistories) > 0 {
	// 	maintainReport.IsHistory = true
	// 	maintainReport.MaintainHistoryTableReport = maintainPreload.MaintenanceRoadHistories
	// } else {
	// 	maintainReport.IsHistory = false
	// }

	maintainReport.MaintainTableReport = maintainTableReport
	maintainReport.MaintainRoadPath = dataTheGeom.TheGeom
	maintainReport.Date = time.Now().Format("02/01/2006")
	maintainReport.Maintain = dataMaintenanceListByID
	maintainReport.ApiKey = os.Getenv("LONGDO_API_KEY")

	filePath := os.Getenv("MAINTENANCE_PROCESS_PDF")

	templateName := os.Getenv("TEMPLATE_MAINTENANCE_HISTORY_PLAN")
	result, err := helpers.InitDataToHtml(templateName, maintainReport, filePath)
	if err != nil {
		logs.Error(err)
		return "", err
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
	if err := chromedp.Run(ctx, helpers.PrintToPDF(string(html), &buf, true)); err != nil {
		logs.Error(err)
		return "", err
	}

	unix := strconv.Itoa(int(time.Now().Unix()))

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return "", err
	}
	fullFilePath := filePath + unix + ".pdf"
	if err := os.WriteFile(fullFilePath, buf, 0777); err != nil {
		logs.Error(err)
		return "", err
	}

	return os.Getenv("STORAGE_IP") + "/" + fullFilePath, nil
}

func (u *Usecase) GetMaintenanceYear(roadId int) ([]int, error) {
	maintenances, err := u.Repo.GetMaintenanceYear(roadId)
	if err != nil {
		logs.Error(err)
		return []int{}, responses.NewAppErr(400, err.Error())
	}
	years := []int{}
	for _, item := range maintenances {
		years = append(years, item.BudgetYear)
	}
	return years, nil
}

func (u *Usecase) GetMaintenanceaAttachmentByID(id int) (*models.MaintenanceAttachment, error) {
	att, err := u.Repo.GetMaintenanceaAttachmentByID(id)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewAppErr(400, err.Error())
	}
	return att, nil
}

func (t *Usecase) GetLastRoadInfoByID(roadID int) (*models.RoadInfoGeomData, error) {
	roadInfo, err := t.Repo.GetLastRoadInfoByID(roadID)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	return roadInfo, nil
}

func (u *Usecase) GetDivisionList(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_maint_history"}, []string{"view_owner_maint_history"})
	//============ end check permission ============
	divitionLists, err := u.Repo.GetDivisionList(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, responses.NewAppErr(500, err.Error())
	}

	var divitionList []models.RefDivisionList
	for _, item := range divitionLists {
		var refDistrictInits []models.RefDistrictInit
		for _, item := range item.Districts {
			if len(item.Depots) != 0 {
				refDistrictInits = append(refDistrictInits, item)
			}
			// var refDistrictInits []models.RefDistrictInit
		}
		if len(refDistrictInits) != 0 {
			item.Districts = refDistrictInits
			divitionList = append(divitionList, item)
		}
	}
	return divitionList, nil
}

func (u *Usecase) GetRoadDropdownList(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_maint_history"}, []string{"view_owner_maint_history"})
	//============ end check permission ============
	data, err := u.Repo.GetRoadDropdownList(isAllData, isOwnerData, depotCode)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return []responses.RoadListInitData{}, nil
		}
		return nil, responses.NewAppErr(500, err.Error())
	}
	sectionDup := make(map[int]bool)
	var groupDatas []models.RoadListInit
	for _, group := range data {
		var groupData models.RoadListInit
		var sectionDatas []models.RoadSectionInit
		for _, section := range group.RoadSection {
			if !sectionDup[section.Id] {
				sectionDup[section.Id] = true
				if len(section.Roads) > 0 {
					var sectionData models.RoadSectionInit
					copier.Copy(&sectionData, &section)
					var roadDatas []models.RoadInit
					for _, road := range section.Roads {
						var roadData models.RoadInit
						copier.Copy(&roadData, &road)
						roadData.Name = helpers.CheckOriginDestinationRoad(road.RefRoadTypeId, section.NameOriginTH, section.NameDestinationTH, road.Name)
						roadDatas = append(roadDatas, roadData)
					}
					sectionData.Roads = roadDatas
					if len(roadDatas) > 0 {
						sectionDatas = append(sectionDatas, sectionData)
					}
				}
			}
		}
		copier.Copy(&groupData, &group)
		sort.SliceStable(sectionDatas, func(i, j int) bool {
			return sectionDatas[i].Number < sectionDatas[j].Number
		})
		groupData.RoadSection = sectionDatas
		if len(sectionDatas) > 0 {
			groupDatas = append(groupDatas, groupData)
		}
	}

	if len(groupDatas) == 0 {
		return []string{}, nil
	}
	return groupDatas, nil
}

func (u *Usecase) GetRoadDropdownListAnalyze(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_maintenance_analysis"}, []string{"view_owner_maintenance_analysis"})
	//============ end check permission ============
	data, err := u.Repo.GetRoadDropdownList(isAllData, isOwnerData, depotCode)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return []responses.RoadListInitData{}, nil
		}
		return nil, responses.NewAppErr(500, err.Error())
	}
	sectionDup := make(map[int]bool)
	var groupDatas []models.RoadListInit
	for _, group := range data {
		var groupData models.RoadListInit
		var sectionDatas []models.RoadSectionInit
		for _, section := range group.RoadSection {
			if !sectionDup[section.Id] {
				sectionDup[section.Id] = true
				if len(section.Roads) > 0 {
					var sectionData models.RoadSectionInit
					copier.Copy(&sectionData, &section)
					var roadDatas []models.RoadInit
					for _, road := range section.Roads {
						var roadData models.RoadInit
						copier.Copy(&roadData, &road)
						roadData.Name = helpers.CheckOriginDestinationRoad(road.RefRoadTypeId, section.NameOriginTH, section.NameDestinationTH, road.Name)
						roadDatas = append(roadDatas, roadData)
					}
					sectionData.Roads = roadDatas
					if len(roadDatas) > 0 {
						sectionDatas = append(sectionDatas, sectionData)
					}
				}
			}
		}
		copier.Copy(&groupData, &group)
		sort.SliceStable(sectionDatas, func(i, j int) bool {
			return sectionDatas[i].Number < sectionDatas[j].Number
		})
		groupData.RoadSection = sectionDatas
		if len(sectionDatas) > 0 {
			groupDatas = append(groupDatas, groupData)
		}
	}

	if len(groupDatas) == 0 {
		return []string{}, nil
	}
	return groupDatas, nil
}

func (u *Usecase) GetRoadDropdownListDashboard(userID int, dataType string, ownerCode string) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := u.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	var isAllData, isOwnerData bool
	switch dataType {
	case "asset":
		isAllData, isOwnerData = helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_asset_dashboard"}, []string{"view_owner_asset_dashboard"})
	case "condition":
		isAllData, isOwnerData = helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_road_condition_dashboard"}, []string{"view_owner_road_condition_dashboard"})
	case "surface":
		isAllData, isOwnerData = helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_surface_dashboard"}, []string{"view_owner_surface_dashboard"})
	case "maintenance":
		isAllData, isOwnerData = helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_maint_history_dashboard"}, []string{"view_owner_maint_history_dashboard"})
	default:
		return nil, responses.NewAppErr(400, "กรุณาเลือกประเภท")
	}
	//
	mapOwnerCode := make(map[string][]string)
	if ownerCode != "" {
		ownerCodeSplit := strings.Split(ownerCode, ",")
		for _, data := range ownerCodeSplit {
			dataSplit := strings.Split(data, "-")
			if len(dataSplit) < 2 {
				return "", responses.NewAppErr(400, "owner_code invalid format please check")
			}

			mapOwnerCode[dataSplit[0]] = append(mapOwnerCode[dataSplit[0]], fmt.Sprintf(`'%s'`, dataSplit[1]))
		}
	}
	var query []string
	for key, value := range mapOwnerCode {
		query = append(query, fmt.Sprintf(`road_section.%s in (%s)`, key, strings.Join(value, ", ")))
	}

	queryRelated := strings.Join(query, " or ")

	//============ end check permission ============
	data, err := u.Repo.GetRoadDropdownListDashboard(isAllData, isOwnerData, depotCode, &queryRelated)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return []responses.RoadListInitData{}, nil
		}
		return nil, responses.NewAppErr(500, err.Error())
	}
	sectionDup := make(map[int]bool)
	var groupDatas []models.RoadListInit
	for _, group := range data {
		var groupData models.RoadListInit
		var sectionDatas []models.RoadSectionInit
		for _, section := range group.RoadSection {
			if !sectionDup[section.Id] {
				sectionDup[section.Id] = true
				if len(section.Roads) > 0 {
					var sectionData models.RoadSectionInit
					copier.Copy(&sectionData, &section)
					sectionData.Number = "ตอนควบคุม " + sectionData.Number
					var roadDatas []models.RoadInit
					for _, road := range section.Roads {
						var roadData models.RoadInit
						copier.Copy(&roadData, &road)
						roadData.Name = helpers.CheckOriginDestinationRoad(road.RefRoadTypeId, section.NameOriginTH, section.NameDestinationTH, road.Name)
						roadDatas = append(roadDatas, roadData)
					}
					sectionData.Roads = roadDatas
					if len(roadDatas) > 0 {
						sectionDatas = append(sectionDatas, sectionData)
					}
				}
			}

		}
		copier.Copy(&groupData, &group)
		groupData.RoadSection = sectionDatas
		sort.Slice(sectionDatas, func(i, j int) bool {
			return sectionDatas[i].Number < sectionDatas[j].Number
		})
		if len(sectionDatas) > 0 {
			groupDatas = append(groupDatas, groupData)
		}
	}

	if len(groupDatas) == 0 {
		return []string{}, nil
	}
	return groupDatas, nil
}

// func (u *Usecase) GetRoadDivisionFilter() (interface{}, error) {
// 	divitionList, err := u.Repo.GetRoadDivisionFilter()
// 	if err != nil {
// 		return nil, responses.NewAppErr(500, err.Error())
// 	}
// 	return divitionList, nil
// }
