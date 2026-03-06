package usecases

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"math"

	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	servicesDB "gitlab.com/mims-api-service/services/database"
	"gitlab.com/mims-api-service/src/maintenanceAnalysis/domains"
	"gorm.io/gorm"
)

type UseCase struct {
	Repo       domains.Repository
	servicesDB servicesDB.ServicesDatabaseDomain
}

func NewUseCase(repo domains.Repository, servicesDB servicesDB.ServicesDatabaseDomain) domains.Usecase {
	return &UseCase{
		Repo:       repo,
		servicesDB: servicesDB,
	}
}

var colors = []string{
	"#008FFB", "#00E396", "#FEB019", "#FF4560", "#775DD0",
	"#3F51B5", "#03A9F4", "#4CAF50", "#F9C1D", "#FF9800",
	"#33B2DF", "#546E7A", "#D4526E", "#13D8AA", "#A5978B",
	"#4ECD4C", "#C7F464", "#81D4FA", "#546E7A", "#FD6A6A",
	"#2B908F", "#F9A3A4", "#90E7E", "#FA4443", "#6D9C27",
	"#449DD1", "#F86E24", "#EA3546", "#662E9B", "#C5D86D",
	"#D7263D", "#1B998B", "#2E294E", "#F46036", "#E2C044",
	"#662E9B", "#F86E24", "#F9C80E", "#EA3546", "#43BCCF",
	"#5C4742", "#A5978B", "#8D5B4C", "#5A2A27", "#C4BBAF",
	"#A300D6", "#7D02EB", "#5653FE", "#2983FF", "#00B1F2",
}

// ////////////////////////////////// LIST ////////////////////////////////////
func (u *UseCase) GetMaintenanceAnalysis(userID int, filter requests.AnalysisFilter, limit, offset int64) ([]responses.AnalysisRes, int64, error) {
	// ============ start check permission ============
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

	total, err := u.Repo.GetMaintenanceAnalysisCount(filter, isAllData, isOwnerData, depotCode)
	if err != nil {
		log.Println(err)
		return nil, 0, responses.NewAppErr(400, err.Error())
	}
	if total == 0 {
		return []responses.AnalysisRes{}, 0, nil
	}

	analysis, err := u.Repo.GetMaintenanceAnalysis(filter, isAllData, isOwnerData, depotCode, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, 0, responses.NewAppErr(400, err.Error())
	}

	var analysisRes []responses.AnalysisRes
	for _, item := range analysis {
		typeAnalysis := ""
		if item.MaintenanceAnalysisTypeId == 1 {
			typeAnalysis = "บำรุงรักษาเชิงกลยุทธ์"
		} else {
			typeAnalysis = "บำรุงรักษาประจำปี"
		}

		status := ""
		if item.Status == "A" {
			status = "สำเร็จ"
		} else if item.Status == "I" {
			status = "กำลังดำเนินการ"
		} else {
			status = "เกิดข้อผิดพลาด"
		}

		analysisRes = append(analysisRes, responses.AnalysisRes{
			ID:                         item.ID,
			TypeAnalysis:               typeAnalysis,
			Name:                       item.Name,
			MaintenanceAnalysisTypeId:  item.MaintenanceAnalysisTypeId,
			MaintenanceConditionTypeId: item.Condition,
			Percentage:                 int(item.Percentage),
			Comment:                    item.Comment,
			Status:                     status,
			AnalysisDate:               item.CreatedAt,
			IsFavorite:                 item.IsFavorite,
		})
	}
	return analysisRes, total, nil
}

// ////////////////////////////////// STEP 1 ////////////////////////////////////
func (u *UseCase) GetMaintenanceAnalysisById(ID int, userID int) (interface{}, error) {
	// ============ start check permission ============
	isAllData := true
	isOwnerData := true
	depotCode := ""
	if userID == 0 {
		isAllData = true
		isOwnerData = true
		depotCode = ""
	} else {
		userInfo, _ := u.servicesDB.UserInfo(userID)
		depotCode = userInfo.RefDepot.DepotCode
		refUserOwnerID := userInfo.RefUserOwnerID
		accessCtrls := userInfo.AccessControl
		accessCtrl := []string{}
		for _, item := range accessCtrls {
			accessCtrl = append(accessCtrl, item.AccessKey)
		}
		isAllData, isOwnerData = helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_maintenance_analysis"}, []string{"view_owner_maintenance_analysis"})
	}

	//============ end check permission ============
	resp, err := u.Repo.GetMaintenanceAnalysisById(ID)
	if err != nil {
		log.Println(err)
		if err.Error() == "record not found" {
			return models.MaintenanceAnalysisPreload{}, responses.NewNotFoundError()
		}
		return models.MaintenanceAnalysisPreload{}, responses.NewAppErr(400, err.Error())
	}

	depot, err := u.Repo.GetRefDepotByRoad()
	if err != nil {
		log.Println(err)
		if err.Error() == "record not found" {
			return "", responses.NewNotFoundError()
		}
		return "", responses.NewAppErr(400, err.Error())
	}

	roads, err := u.Repo.GetMaintenanceAnalysisRoadByID(ID)
	if err != nil {
		logs.Error(err)
		return responses.ResponseMaintenanceAnalysis{}, responses.NewAppErr(400, err.Error())
	}
	isShow := true
	if isOwnerData && !isAllData {
		if len(roads) == 0 {
			isShow = false
		}
		for _, prepare := range roads {
			roadID := prepare.RoadID
			val, isVal := depot[roadID]
			if !isVal {
				isShow = false
			} else {
				isShow = false
				if val == depotCode {
					isShow = true
				}

				goto EXITLOOP
			}
		}
	}
EXITLOOP:

	if !isOwnerData && !isAllData {
		isShow = false
	}

	if !isShow {
		return responses.ResponseMaintenanceAnalysis{}, responses.NewAppErr(403, "")
	}
	// return roads, nil
	var analysisRes responses.AnalysisByIDRes
	copier.Copy(&analysisRes, resp)
	analysisRes.Roads = roads
	return analysisRes, nil
}

// ค้นหา
func (u *UseCase) CreateMaintenanceAnalysis(userID int, req requests.MaintenanceAnalysis) (interface{}, error) {
	dateTime := time.Now()
	var maintenanceAnalysis models.MaintenanceAnalysis
	maintenanceAnalysis.Name = req.Name
	maintenanceAnalysis.MaintenanceAnalysisTypeId = req.MaintenanceAnalysisTypeId
	maintenanceAnalysis.SurfaceTypeId = req.SurfaceTypeId
	maintenanceAnalysis.LaneTypeId = req.LaneTypeId
	maintenanceAnalysis.Iri1 = req.Iri1
	maintenanceAnalysis.Iri2 = req.Iri2
	maintenanceAnalysis.Aadt1 = req.Aadt1
	maintenanceAnalysis.Aadt2 = req.Aadt2
	maintenanceAnalysis.Ifi1 = req.Ifi1
	maintenanceAnalysis.Ifi2 = req.Ifi2
	maintenanceAnalysis.GroupKm = float64(req.GroupKm)
	maintenanceAnalysis.Percentage = 0
	maintenanceAnalysis.Status = ""
	maintenanceAnalysis.IsLatest = true
	maintenanceAnalysis.IsDeleted = false
	maintenanceAnalysis.CreatedBy = userID
	maintenanceAnalysis.CreatedAt = dateTime

	analysis, err := u.Repo.CreateMaintenanceAnalysis(maintenanceAnalysis)
	if err != nil {
		logs.Error(err)
		return responses.ResponseMaintenanceAnalysis{}, responses.NewAppErr(400, err.Error())
	}

	// create road
	err = u.Repo.CreateMaintenanceAnalysisRoad(analysis.ID, userID, req.Roads)
	if err != nil {
		logs.Error(err)
		return responses.ResponseMaintenanceAnalysis{}, responses.NewAppErr(400, err.Error())
	}

	res, err := u.GetMaintenanceAnalysisById(analysis.ID, 0)
	if err != nil {
		logs.Error(err)
		return []string{}, responses.NewAppErr(400, err.Error())
	}
	go u.CreateMaintenanceAnalysisBackground(analysis.ID, req)
	return res, nil
}

func (u *UseCase) CreateMaintenanceAnalysisBackground(analysisID int, req requests.MaintenanceAnalysis) error {

	// create prepare data
	_, err := u.GetPrepareData(analysisID, req)
	if err != nil {
		logs.Error(err)
		return responses.NewAppErr(400, err.Error())
	}

	err = u.Repo.UpdatePrepareDataStatusById(analysisID)
	if err != nil {
		logs.Error(err)
		return responses.NewAppErr(400, err.Error())
	}

	return nil
}

// ค้นหา
func (u *UseCase) UpdateMaintenanceAnalysis(ID, userID int, req requests.MaintenanceAnalysis) (interface{}, error) {
	dateTime := time.Now()
	maintenanceAnalysis, err := u.Repo.GetMaintenanceAnalysisById(ID)
	if err != nil {
		if err.Error() == "record not found" {
			return "", responses.NewNotFoundError()
		}
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	maintenanceAnalysis.ID = ID
	maintenanceAnalysis.Name = req.Name
	maintenanceAnalysis.MaintenanceAnalysisTypeId = req.MaintenanceAnalysisTypeId
	maintenanceAnalysis.SurfaceTypeId = req.SurfaceTypeId
	maintenanceAnalysis.LaneTypeId = req.LaneTypeId
	maintenanceAnalysis.Iri1 = req.Iri1
	maintenanceAnalysis.Iri2 = req.Iri2
	maintenanceAnalysis.Aadt1 = req.Aadt1
	maintenanceAnalysis.Aadt2 = req.Aadt2
	maintenanceAnalysis.Ifi1 = req.Ifi1
	maintenanceAnalysis.Ifi2 = req.Ifi2
	maintenanceAnalysis.GroupKm = float64(req.GroupKm)
	maintenanceAnalysis.Percentage = 0
	maintenanceAnalysis.Status = ""
	maintenanceAnalysis.IsLatest = true
	maintenanceAnalysis.IsDeleted = false
	maintenanceAnalysis.CreatedBy = userID
	maintenanceAnalysis.CreatedAt = dateTime
	maintenanceAnalysis.PreviousID = ID

	newMaintenanceAnalysis, err := u.Repo.CreateMaintenanceAnalysis(maintenanceAnalysis)
	if err != nil {
		logs.Error(err)
		return responses.ResponseMaintenanceAnalysis{}, responses.NewAppErr(400, err.Error())
	}

	/// create road
	err = u.Repo.CreateMaintenanceAnalysisRoad(newMaintenanceAnalysis.ID, userID, req.Roads)
	if err != nil {
		logs.Error(err)
		return responses.ResponseMaintenanceAnalysis{}, responses.NewAppErr(400, err.Error())
	}

	res, err := u.GetMaintenanceAnalysisById(newMaintenanceAnalysis.ID, 0)
	if err != nil {
		logs.Error(err)
		return []string{}, responses.NewAppErr(400, err.Error())
	}

	go u.UpdateMaintenanceAnalysisBackground(newMaintenanceAnalysis.ID, req)

	return res, nil
}

func (u *UseCase) UpdateMaintenanceAnalysisBackground(analysisID int, req requests.MaintenanceAnalysis) error {

	// create prepare data
	_, err := u.GetPrepareData(analysisID, req)
	if err != nil {
		logs.Error(err)
		return responses.NewAppErr(400, err.Error())
	}

	err = u.Repo.UpdatePrepareDataStatusById(analysisID)
	if err != nil {
		logs.Error(err)
		return responses.NewAppErr(400, err.Error())
	}

	return nil
}

func (u *UseCase) SelectdPrepareData(ID int, prepareDataID []int) error {
	err := u.Repo.SelectdPrepareData(ID, prepareDataID)
	if err != nil {
		logs.Error(err)
		return responses.NewAppErr(400, err.Error())
	}
	return nil
}

// ////////////////////////////////// STEP 2 ////////////////////////////////////
func (u *UseCase) GetMaintenanceAnalysisConditionById(ID int, prepareDataID []int) (interface{}, error) {
	prepareData, err := u.Repo.GetPrepareDataByAnalysis(ID, prepareDataID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	totalKm := 0.0
	for _, item := range prepareData {
		totalKm += item.Length / 1000
	}
	iriAvg := 0.0
	ifiAvg := 0.0

	iri := 0.0
	ifi := 0.0
	for _, item := range prepareData {
		iri += (item.Iri * (item.Length / 1000))
		ifi += (item.Ifi * (item.Length / 1000))
	}

	iriAvg = iri / totalKm
	if math.IsNaN(iriAvg) {
		iriAvg = 0
	}

	ifiAvg = ifi / totalKm
	if math.IsNaN(ifiAvg) {
		ifiAvg = 0
	}

	var res responses.AnalysisStep2Res

	analysis, err := u.Repo.GetMaintenanceAnalysisById(ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {

		} else {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
	} else {
		res.ConditionID = &analysis.Condition
		res.Discount = analysis.Discount
		res.Year = analysis.Year
		res.Target = analysis.Target
		res.NumberPlan = analysis.NumberPlan
		res.Comment = &analysis.Comment
	}
	plans, err := u.Repo.GetMaintenanceAnalysisPlanByAnalysisID(ID)
	if err != nil {
		log.Println(err)
		if err.Error() == "record not found" {
			logs.Error(err)
		}
		return "", responses.NewAppErr(400, err.Error())
	}
	surface := ""
	if analysis.SurfaceTypeId == 1 {
		surface = "ลาดยาง"
	} else {
		surface = "คอนกรีต"
	}
	res.SurfaceType = surface
	res.IRIAvg = iriAvg
	res.IFIAvg = ifiAvg
	res.TotalKm = totalKm
	res.Budget = analysis.Budget
	res.IRI = analysis.Iri
	res.IFI = analysis.Ifi
	if len(plans) > 0 {
		res.Plans = plans
	} else {
		res.Plans = []string{}
	}

	// analysis.
	return res, nil
}

func (u *UseCase) UpdateMaintenanceAnalysisCondition(ID, userID int, req requests.AnalyzingReq) (interface{}, error) {
	maintenanceAnalysis, err := u.Repo.GetMaintenanceAnalysisById(ID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	maintenanceAnalysis.Name = req.Name
	maintenanceAnalysis.Condition = *req.ConditionID
	maintenanceAnalysis.Discount = req.Discount
	maintenanceAnalysis.Year = req.Year
	maintenanceAnalysis.Target = req.Target
	maintenanceAnalysis.NumberPlan = req.NumberPlan
	maintenanceAnalysis.Comment = *req.Comment
	maintenanceAnalysis.Budget = req.Budget
	maintenanceAnalysis.Iri = req.Iri
	maintenanceAnalysis.Ifi = req.Ifi
	maintenanceAnalysis.Status = "I"

	filter := ""
	if maintenanceAnalysis.MaintenanceAnalysisTypeId == 1 {
		filter = "ผิวลาดยาง"
	} else {
		filter = "คอนกรีต"
	}
	lane := ""
	if maintenanceAnalysis.LaneTypeId == 0 {
		lane += "ทั้งหมด"
	} else {
		lane += fmt.Sprintf("%d", maintenanceAnalysis.LaneTypeId)
	}

	kmGrp := fmt.Sprintf("%d", int(maintenanceAnalysis.GroupKm)) + " กม."
	discount := 0.0
	if maintenanceAnalysis.Discount == nil {
		discount = 0.0
	} else {
		discount = helpers.RoundFloat(*maintenanceAnalysis.Discount, 2)
	}
	dataFillter, err := u.Repo.GetMaintenanceAnalysisRoadDataByID(ID)
	if err != nil {
		logs.Error(err)
		return requests.MaintenanceAnalysisStrategic{}, responses.NewAppErr(400, err.Error())
	}

	roads := []string{}
	roadChkDup := make(map[string]string)
	for _, prepare := range dataFillter.PrepareData {
		_, ok := roadChkDup[prepare.RoadName]
		if !ok {
			roadChkDup[prepare.RoadName] = prepare.RoadName
			roads = append(roads, prepare.RoadName)
		}
	}
	conditionFilters := []string{}
	conditionFilter := ""
	isIri := false
	iri1 := ""
	iri2 := ""
	isAadt := false
	aadt1 := ""
	aadt2 := ""
	isIfi := false
	ifi1 := ""
	ifi2 := ""

	// IRI
	if maintenanceAnalysis.Iri1 != nil {
		isIri = true
		iri1 = fmt.Sprintf("%s < ", helpers.FormatNumber(int(*maintenanceAnalysis.Iri1)))
	}
	if maintenanceAnalysis.Iri2 != nil {
		isIri = true
		iri2 = fmt.Sprintf(" < %s", helpers.FormatNumber(int(*maintenanceAnalysis.Iri2)))
	}
	if isIri {
		conditionFilters = append(conditionFilters, iri1+"IRI"+iri2)
	}

	// AADT
	if maintenanceAnalysis.Aadt1 != nil {
		isAadt = true
		aadt1 = fmt.Sprintf("%s < ", helpers.FormatNumber(int(*maintenanceAnalysis.Aadt1)))
	}
	if maintenanceAnalysis.Aadt2 != nil {
		isAadt = true
		aadt2 = fmt.Sprintf(" < %s", helpers.FormatNumber(int(*maintenanceAnalysis.Aadt2)))
	}
	if isAadt {
		conditionFilters = append(conditionFilters, aadt1+"AADT"+aadt2)
	}

	// Ifi
	if maintenanceAnalysis.Ifi1 != nil {
		isIfi = true
		ifi1 = fmt.Sprintf("%s < ", helpers.FormatNumber(int(*maintenanceAnalysis.Ifi1)))
	}
	if maintenanceAnalysis.Ifi2 != nil {
		isIfi = true
		ifi2 = fmt.Sprintf(" < %s", helpers.FormatNumber(int(*maintenanceAnalysis.Ifi2)))
	}
	if isIfi {
		conditionFilters = append(conditionFilters, ifi1+"IFI"+ifi2)
	}

	conditionFilter = strings.Join(conditionFilters, ", ")
	condition := dataFillter.ConditionData.Name
	target := dataFillter.TargetData.Name
	roadName := helpers.Implode(",", roads)
	filterStr := fmt.Sprintf("สายทาง: %s ตัวกรอง: %s ช่องจราจร: %s จัดกลุ่ม: %s %v ส่วนลด: %v%s เงื่อนไข: %s เป้าหมาย: %s", roadName, filter, lane, kmGrp, conditionFilter, discount, "%", condition, target)
	maintenanceAnalysis.FilterData = filterStr

	// maintenanceAnalysis
	err = u.Repo.UpdateMaintenanceAnalysisStep2(ID, maintenanceAnalysis)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	if len(req.Plans) >= 0 {
		var maintenanceAnalysisPlanList []models.MaintenanceAnalysisPlan
		for _, item := range req.Plans {
			var maintenanceAnalysisPlan models.MaintenanceAnalysisPlan
			maintenanceAnalysisPlan.ID = item.ID
			maintenanceAnalysisPlan.MaintenanceAnalysisID = ID
			maintenanceAnalysisPlan.Plan1 = item.Plan1
			maintenanceAnalysisPlan.Plan2 = item.Plan2
			maintenanceAnalysisPlan.Plan3 = item.Plan3
			maintenanceAnalysisPlan.PlanYear = item.PlanYear
			maintenanceAnalysisPlan.CreatedBy = userID
			maintenanceAnalysisPlan.CreatedAt = time.Now()
			maintenanceAnalysisPlanList = append(maintenanceAnalysisPlanList, maintenanceAnalysisPlan)
		}
		err = u.Repo.CreateMaintenanceAnalysisPlan(ID, maintenanceAnalysisPlanList)
		if err != nil {
			logs.Error(err)
			return requests.MaintenanceAnalysisStrategic{}, responses.NewAppErr(400, err.Error())
		}
	}

	err = u.SelectdPrepareData(ID, req.PrepareDataID)
	if err != nil {
		logs.Error(err)
		return requests.MaintenanceAnalysisStrategic{}, responses.NewAppErr(400, err.Error())
	}

	res, err := u.GetMaintenanceAnalysisConditionById(ID, req.PrepareDataID)
	if err != nil {
		logs.Error(err)
		return requests.MaintenanceAnalysisStrategic{}, responses.NewAppErr(400, err.Error())
	}

	if maintenanceAnalysis.PreviousID != 0 {
		u.Repo.ClearPrepareData(maintenanceAnalysis.PreviousID)
	}

	return res, nil
}

func (u *UseCase) DeleteMaintenanceAnalysis(ID int) error {
	err := u.Repo.DeleteMaintenanceAnalysis(ID)
	if err != nil {
		logs.Error(err)
		return responses.NewAppErr(400, err.Error())
	}
	return nil
}

func (u *UseCase) CopyMaintenanceAnalysis(ID int) (interface{}, error) {
	res, err := u.Repo.CopyMaintenanceAnalysis(ID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	return res, nil
}

func (u *UseCase) FavoriteMaintenanceAnalysis(ID int) (responses.AnalysisIsFavorite, error) {
	var analysisIsFavorite responses.AnalysisIsFavorite
	err := u.Repo.FavoriteMaintenanceAnalysis(ID)
	if err != nil {
		logs.Error(err)
		return analysisIsFavorite, err
	}

	maintenanceAnalysis, err := u.Repo.GetMaintenanceAnalysisById(ID)
	if err != nil {
		logs.Error(err)
		return analysisIsFavorite, err
	}

	analysisIsFavorite.IsFavorite = maintenanceAnalysis.IsFavorite

	return analysisIsFavorite, nil
}

func (u *UseCase) DashboardStrategicMaintenanceAnalysis(id string, userID int) (interface{}, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}
	// ============ start check permission ============
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
	depot, err := u.Repo.GetRefDepotByRoad()
	if err != nil {
		log.Println(err)
		if err.Error() == "record not found" {
			return "", responses.NewNotFoundError()
		}
		return "", responses.NewAppErr(400, err.Error())
	}

	roads, err := u.Repo.GetMaintenanceAnalysisRoadByID(ID)
	if err != nil {
		logs.Error(err)
		return responses.ResponseMaintenanceAnalysis{}, responses.NewAppErr(400, err.Error())
	}
	isShow := true
	if isOwnerData && !isAllData {
		if len(roads) == 0 {
			isShow = false
		}
		for _, prepare := range roads {
			roadID := prepare.RoadID
			val, isVal := depot[roadID]
			if !isVal {
				isShow = false
			} else {
				isShow = false
				if val == depotCode {
					isShow = true
				}

				goto EXITLOOP
			}
		}
	}
EXITLOOP:

	if !isOwnerData && !isAllData {
		isShow = false
	}

	if !isShow {
		return responses.ResponseMaintenanceAnalysis{}, responses.NewAppErr(403, "")
	}

	idInteger, err := strconv.Atoi(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	modelResults, err := u.Repo.GetModelResultDataById(idInteger)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	latestPrepareData, err := u.Repo.GetLatestPrepareDataByMaintenanceAnalysisId(idInteger)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysis, err := u.Repo.GetMaintenanceAnalysisById(idInteger)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysisPlan, err := u.Repo.GetMaintenanceAnalysisPlanByAnalysisID(idInteger)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	// color := []string{"#008FFB", "#00E396", "#FEAF1A", "#FE4560", "#775DD0"}

	interventionCriteria, err := u.Repo.GetInterventionCriteria()
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return models.DashboardStrategicMaintenance{}, err
	}

	colorMap := make(map[string]string)
	for index, item := range interventionCriteria {
		colorMap[item.MaintenanceStandardName] = colors[index]
	}

	colorMap["ซ่อมบำรุงปกติ"] = "#000000"
	// helpers.PrintlnJson("colorMap", colorMap)
	// return colors, nil
	switch maintenanceAnalysis.Condition {
	case 1:
		result, err := u.strategicMaintenanceAnalysisUnlimitBudget(idInteger, modelResults, latestPrepareData, maintenanceAnalysis, maintenanceAnalysisPlan, colors, colorMap)
		if err != nil {
			logs.Error(err)
			return []string{}, err
		}
		return result, nil

	case 2:
		result, err := u.strategicMaintenanceAnalysisLimitBudget(idInteger, modelResults, latestPrepareData, maintenanceAnalysis, maintenanceAnalysisPlan, colors, colorMap)
		if err != nil {
			logs.Error(err)
			return []string{}, err
		}
		return result, nil

	case 3:
		result, err := u.strategicMaintenanceAnalysisTargetIri(idInteger, modelResults, latestPrepareData, maintenanceAnalysis, maintenanceAnalysisPlan, colors, colorMap)
		if err != nil {
			logs.Error(err)
			return []string{}, err
		}
		return result, nil
	}

	return []string{}, nil
}

func (u *UseCase) strategicMaintenanceAnalysisUnlimitBudget(idInteger int, modelResults []models.ModelResult, latestPrepareData models.PrepareData, maintenanceAnalysis models.MaintenanceAnalysis, maintenanceAnalysisPlan []models.MaintenanceAnalysisPlan, color []string, colorMap map[string]string) (interface{}, error) {
	maintenanceAnalysisStrategicBudgetType, err := u.Repo.GetMaintenanceAnalysisStrategicBudgetTypeById(maintenanceAnalysis.Condition)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysisStrategiTarget, err := u.Repo.GetMaintenanceAnalysisStrategiTargetById(*maintenanceAnalysis.Target)
	if err != nil && err.Error() != "record not found" {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysisRoad, err := u.Repo.GetMaintenanceAnalysisRoadNameById(maintenanceAnalysis.ID)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysisRoads := []string{}
	for _, item := range maintenanceAnalysisRoad {
		maintenanceAnalysisRoads = append(maintenanceAnalysisRoads, item.Name)
	}

	lane := ""
	switch maintenanceAnalysis.LaneTypeId {
	case 0:
		lane = "ทั้งหมด"
	case 1:
		lane = "1"
	case 2:
		lane = "2"
	case 3:
		lane = "3"
	}

	surface := ""
	switch maintenanceAnalysis.SurfaceTypeId {
	case 1:
		surface = "ลาดยาง"
	case 2:
		surface = "คอนกรีต"
	}

	comment := ""
	if maintenanceAnalysis.Comment == "" {
		comment = "-"
	} else {
		comment = maintenanceAnalysis.Comment
	}

	var filter = []string{}
	if maintenanceAnalysis.Iri1 != nil || maintenanceAnalysis.Iri2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Iri1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Iri1) + " < "
		}
		textFilter += "IRI"
		if maintenanceAnalysis.Iri2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Iri2)
		}

		filter = append(filter, textFilter)
	}

	if maintenanceAnalysis.Aadt1 != nil || maintenanceAnalysis.Aadt2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Aadt1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Aadt1) + " < "
		}
		textFilter += "AADT"
		if maintenanceAnalysis.Aadt2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Aadt2)
		}

		filter = append(filter, textFilter)
	}

	if maintenanceAnalysis.Ifi1 != nil || maintenanceAnalysis.Ifi2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Ifi1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Ifi1) + " < "
		}
		textFilter += "IFI"
		if maintenanceAnalysis.Ifi2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Ifi2)
		}

		filter = append(filter, textFilter)
	}

	var dashBoard models.DashboardStrategicMaintenance
	dashBoard.Road = maintenanceAnalysisRoads
	dashBoard.Filter.Lane = lane
	dashBoard.Filter.SurfaceType = surface
	dashBoard.Comment = comment
	dashBoard.Filter.Km = maintenanceAnalysis.GroupKm
	dashBoard.Condition.Condition = maintenanceAnalysisStrategicBudgetType.Name
	dashBoard.Condition.Target = maintenanceAnalysisStrategiTarget.Name
	dashBoard.Condition.Discount = maintenanceAnalysis.Discount
	dashBoard.Graph1.Name = "ข้อมูลประจำปี"
	dashBoard.Bar1.Name = "ข้อมูลประจำปี"
	dashBoard.Bar2.Name = "วิธีการซ่อม"
	dashBoard.NumberPlan = *maintenanceAnalysis.NumberPlan
	dashBoard.Filter.Filter = filter

	min := 9999
	for _, item := range modelResults {
		if item.AnalystYear <= min {
			min = item.AnalystYear
		}
	}

	year := *maintenanceAnalysis.Year
	analystYear := min

	var years []int
	for indexYear := 0; indexYear < year; indexYear++ {
		years = append(years, analystYear+indexYear)
	}

	dashBoard.Graph1.Lable = years
	dashBoard.Bar1.Lable = years
	dashBoard.Bar2.Lable = years

	method := []string{"ซ่อมบำรุงปกติ", "ไม่จำกัดงบประมาณ"}
	valueGraph1s := map[string][]float64{}
	valueBar1s := map[string][]float64{}
	for indexPlan := 1; indexPlan < len(method)+1; indexPlan++ {
		var dataset2s models.StrategicDatasetBar2
		key := ""
		if indexPlan == 2 {
			key = "ซ่อมบำรุงปกติ"
		} else {
			key = "ไม่จำกัดงบประมาณ"
		}

		dataset2s.Plan = key
		for indexYear := 1; indexYear <= year; indexYear++ {
			budgetSummaryBar1 := 0.0
			sumGraph1 := 0.0
			sumLengthGraph1 := 0.0
			valueBar2s := map[string][]float64{}
			sumBar2s := map[string]int{}
			keyBar2s := []string{}
			sumBar2 := 0
			for _, item := range modelResults {
				if indexYear == item.Year && indexPlan == item.PlanSequence {
					if item.Data.Repair {
						sumGraph1 = sumGraph1 + item.Data.RweResult.Iri*float64(item.Data.PrepareDataBefore.Length)
						budgetSummaryBar1 = budgetSummaryBar1 + item.Data.UsedBudget

						valueBar2s[item.Data.IcResult.Name] = append(valueBar2s[item.Data.IcResult.Name], item.Data.UsedBudget)
						if !helpers.Contains(keyBar2s, item.Data.IcResult.Name) {
							keyBar2s = append(keyBar2s, item.Data.IcResult.Name)
						}
						sumBar2s[item.Data.IcResult.Name] = sumBar2s[item.Data.IcResult.Name] + int(item.Data.UsedBudget)
						sumBar2 = sumBar2 + int(item.Data.UsedBudget)
					} else {
						sumGraph1 = sumGraph1 + item.Data.DeteriorationResult.Result.Iri*float64(item.Data.PrepareDataBefore.Length)
					}
					sumLengthGraph1 = sumLengthGraph1 + float64(item.Data.PrepareDataBefore.Length)
				}
			}

			result := helpers.RoundFloat(float64(sumGraph1)/float64(sumLengthGraph1), 2)
			if math.IsNaN(result) {
				result = 0
			}

			valueGraph1s[key] = append(valueGraph1s[key], result)
			valueBar1s[key] = append(valueBar1s[key], budgetSummaryBar1)

			datset2Lables := []string{}
			datset2Values := []float64{}
			datset2Budgets := []float64{}
			datset2Colors := []string{}
			for _, itemKey := range keyBar2s {
				datset2Lables = append(datset2Lables, itemKey)
				summary := 0.0
				summaryBudget := 0.0
				for _, itemValue := range valueBar2s[itemKey] {
					summary = summary + ((float64(itemValue)/float64(sumBar2)*100)/100)*100
					summaryBudget = summaryBudget + float64(itemValue)
				}
				if math.IsNaN(summary) {
					summary = 0.0
				}

				color := colorMap[itemKey]
				if color == "" {
					color = "#000000"
				}
				datset2Colors = append(datset2Colors, color)
				datset2Values = append(datset2Values, helpers.RoundFloat(summary, 2))
				datset2Budgets = append(datset2Budgets, summaryBudget)
			}

			var dataBar2 models.StrategicDataBar2
			dataBar2.Lable = datset2Lables
			dataBar2.Value = datset2Values
			dataBar2.Budget = datset2Budgets
			dataBar2.Color = datset2Colors
			dataset2s.Data = append(dataset2s.Data, dataBar2)
		}
		if key != "ซ่อมบำรุงปกติ" {
			dashBoard.Bar2.Datasets = append(dashBoard.Bar2.Datasets, dataset2s)
		}
	}

	colorSort := []string{}
	for index, item := range method {
		dashBoard.Graph1.Value = append(dashBoard.Graph1.Value, valueGraph1s[item])
		colorSort = append(colorSort, color[index])
	}

	dashBoard.Graph1.Line = method

	dashBoard.Graph1.Color = colorSort
	colorSortBar1 := []string{}
	var datasets []models.DatasetBar1
	for index, item := range method {
		if item != "ซ่อมบำรุงปกติ" {
			var dataset models.DatasetBar1
			dataset.Lable = item
			dataset.Value = valueBar1s[item]
			datasets = append(datasets, dataset)
			colorSortBar1 = append(colorSortBar1, color[index])
		}
	}

	dashBoard.Bar1.Datasets = datasets

	dashBoard.Bar1.Color = colorSortBar1

	methodTable := []string{"ไม่จำกัดงบประมาณ", "ซ่อมบำรุงปกติ"}
	var table models.StrategicTable
	titlePlan := []string{"สรุปรวม", "ไม่จำกัดงบประมาณ"}
	for _, item := range titlePlan {
		switch item {
		case "สรุปรวม":
			title := []string{"แผนงบประมาณ (บาท)", "งบประมาณซ่อมบำรุง (บาท)", "IRI (ก่อนวิเคราะห์)", "IRI (หลังวิเคราะห์)"}
			var summaryPlanTables []models.StrategicSummaryPlanTable
			for _, item := range title {
				var summaryPlanTable models.StrategicSummaryPlanTable
				var dataTables []models.StrategicSummaryDataTable
				for indexMethod := len(methodTable); indexMethod > 0; indexMethod-- {
					var dataTable models.StrategicSummaryDataTable
					values := []float64{}
					for indexYear := 1; indexYear <= year; indexYear++ {
						switch item {
						case "แผนงบประมาณ (บาท)":
							values = append(values, 0)
						case "งบประมาณซ่อมบำรุง (บาท)":
							sum := 0.0
							for _, item := range modelResults {
								if indexYear == item.Year && item.PlanSequence == indexMethod {
									if item.Data.Repair {
										sum = sum + item.Data.UsedBudget
									}
								}
							}
							values = append(values, sum)
						case "IRI (ก่อนวิเคราะห์)":
							sumLength := 0.0
							sumIri := 0.0
							for _, item := range modelResults {
								if indexYear == item.Year && item.PlanSequence == indexMethod {
									sumIri = sumIri + (item.Data.PrepareData.RoadCondition.Iri * float64(item.Data.PrepareDataBefore.Length))
									sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
								}
							}

							result := helpers.RoundFloat(sumIri/sumLength, 2)
							if math.IsNaN(result) {
								result = 0
							}

							values = append(values, result)
						case "IRI (หลังวิเคราะห์)":
							sumLength := 0.0
							sumIri := 0.0
							for _, item := range modelResults {
								if indexYear == item.Year && item.PlanSequence == indexMethod {
									if item.Data.Ic {
										sumIri = sumIri + (item.Data.RweResult.Iri * float64(item.Data.PrepareDataBefore.Length))
										sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
									} else {
										sumIri = sumIri + (item.Data.DeteriorationResult.Result.Iri * float64(item.Data.PrepareDataBefore.Length))
										sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
									}
								}
							}

							result := helpers.RoundFloat(sumIri/sumLength, 2)
							if math.IsNaN(result) {
								result = 0
							}

							values = append(values, result)
						}
					}

					dataTable.Value = values
					dataTable.Name = methodTable[indexMethod-1]
					dataTables = append(dataTables, dataTable)
				}

				summaryPlanTable.Name = item
				summaryPlanTable.Data = dataTables
				summaryPlanTables = append(summaryPlanTables, summaryPlanTable)
			}

			table.Summary = summaryPlanTables

		case "ไม่จำกัดงบประมาณ":
			titles := []string{}
			for _, item := range modelResults {
				for indexYear := 1; indexYear <= year; indexYear++ {
					if !helpers.Contains(titles, item.Data.IcResult.Name) && item.Data.IcResult.Name != "" {
						if item.Data.Repair {
							titles = append(titles, item.Data.IcResult.Name)
						}
					}
				}
			}

			var dataTables []models.StrategicDataTable
			var dataTable models.StrategicDataTable
			var DataTableOthers []models.StrategicDataTableOther
			dataTable.MethodName = "งานบำรุงปกติ"
			for indexYear := 1; indexYear <= year; indexYear++ {
				sumKm := 0.0
				for _, itemModel := range modelResults {
					if indexYear == itemModel.Year {
						if !itemModel.Data.Repair {
							sumKm = float64(sumKm) + math.Abs(float64(itemModel.KmEnd)-float64(itemModel.KmStart))
						}
					}
				}
				var dataTableOther models.StrategicDataTableOther
				dataTableOther.Budget = 0
				dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
				dataTableOther.Year = analystYear + indexYear - 1
				DataTableOthers = append(DataTableOthers, dataTableOther)
			}
			dataTable.Data = DataTableOthers
			dataTables = append(dataTables, dataTable)

			for _, itemTitle := range titles {
				var dataTable models.StrategicDataTable
				var DataTableOthers []models.StrategicDataTableOther
				dataTable.MethodName = itemTitle
				for indexYear := 1; indexYear <= year; indexYear++ {
					sumBudget := 0.0
					sumKm := 0
					for _, itemModel := range modelResults {
						if itemModel.Data.IcResult.Name == itemTitle && indexYear == itemModel.Year {
							if itemModel.Data.Repair {
								sumBudget = sumBudget + itemModel.Data.UsedBudget
								sumKm = sumKm + itemModel.Data.PrepareDataBefore.Length
							}
						}
					}
					var dataTableOther models.StrategicDataTableOther
					dataTableOther.Budget = sumBudget
					dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
					dataTableOther.Year = analystYear + indexYear - 1
					DataTableOthers = append(DataTableOthers, dataTableOther)
				}
				dataTable.Data = DataTableOthers
				dataTables = append(dataTables, dataTable)
			}
			table.UnlimitedPlan = dataTables
		}
	}

	if len(table.Plan1) == 0 {
		table.Plan1 = []models.StrategicDataTable{}
	}

	if len(table.Plan2) == 0 {
		table.Plan2 = []models.StrategicDataTable{}
	}

	if len(table.Plan3) == 0 {
		table.Plan3 = []models.StrategicDataTable{}
	}

	if len(table.UnlimitedPlan) == 0 {
		table.UnlimitedPlan = []models.StrategicDataTable{}
	}

	dashBoard.Table = table

	return dashBoard, nil
}

func (u *UseCase) strategicMaintenanceAnalysisLimitBudget(idInteger int, modelResults []models.ModelResult, latestPrepareData models.PrepareData, maintenanceAnalysis models.MaintenanceAnalysis, maintenanceAnalysisPlan []models.MaintenanceAnalysisPlan, color []string, colorMap map[string]string) (interface{}, error) {
	maintenanceAnalysisStrategicBudgetType, err := u.Repo.GetMaintenanceAnalysisStrategicBudgetTypeById(maintenanceAnalysis.Condition)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysisStrategiTarget, err := u.Repo.GetMaintenanceAnalysisStrategiTargetById(*maintenanceAnalysis.Target)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysisRoad, err := u.Repo.GetMaintenanceAnalysisRoadNameById(maintenanceAnalysis.ID)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysisRoads := []string{}
	for _, item := range maintenanceAnalysisRoad {
		maintenanceAnalysisRoads = append(maintenanceAnalysisRoads, item.Name)
	}

	min := 9999
	for _, item := range modelResults {
		if item.AnalystYear <= min {
			min = item.AnalystYear
		}
	}

	year := *maintenanceAnalysis.Year
	analystYear := min

	numberPlanDefault := 2
	numberPlanWithDefault := *maintenanceAnalysis.NumberPlan + numberPlanDefault
	numberPlan := *maintenanceAnalysis.NumberPlan
	years := []int{}
	lane := ""

	maintenanceAnalysisPlanMap := make(map[int]map[string]float64)
	for _, item := range maintenanceAnalysisPlan {
		maintenanceAnalysisPlanMapSub := map[string]float64{}
		for indexPlan := 1; indexPlan < *maintenanceAnalysis.NumberPlan+1; indexPlan++ {
			switch indexPlan {
			case 1:
				maintenanceAnalysisPlanMapSub["plan_1"] = *item.Plan1
			case 2:
				maintenanceAnalysisPlanMapSub["plan_2"] = *item.Plan2
			case 3:
				maintenanceAnalysisPlanMapSub["plan_3"] = *item.Plan3
			}
		}
		maintenanceAnalysisPlanMap[int(item.PlanYear)] = maintenanceAnalysisPlanMapSub
	}

	switch maintenanceAnalysis.LaneTypeId {
	case 0:
		lane = "ทั้งหมด"
	case 1:
		lane = "1"
	case 2:
		lane = "2"
	case 3:
		lane = "3"
	}

	surface := ""
	switch maintenanceAnalysis.SurfaceTypeId {
	case 1:
		surface = "ลาดยาง"
	case 2:
		surface = "คอนกรีต"
	}

	comment := ""
	if maintenanceAnalysis.Comment == "" {
		comment = "-"
	} else {
		comment = maintenanceAnalysis.Comment
	}

	var dashBoard models.DashboardStrategicMaintenance
	dashBoard.Road = maintenanceAnalysisRoads
	dashBoard.Filter.Lane = lane
	dashBoard.Filter.SurfaceType = surface
	dashBoard.Comment = comment
	dashBoard.Filter.Km = maintenanceAnalysis.GroupKm
	dashBoard.Condition.Condition = maintenanceAnalysisStrategicBudgetType.Name
	dashBoard.Condition.Target = maintenanceAnalysisStrategiTarget.Name
	dashBoard.Condition.Discount = maintenanceAnalysis.Discount
	dashBoard.Graph1.Name = "ข้อมูลประจำปี"
	dashBoard.Bar1.Name = "ข้อมูลประจำปี"
	dashBoard.Bar2.Name = "วิธีการซ่อม"
	dashBoard.NumberPlan = *maintenanceAnalysis.NumberPlan

	var filter = []string{}
	if maintenanceAnalysis.Iri1 != nil || maintenanceAnalysis.Iri2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Iri1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Iri1) + " < "
		}
		textFilter += "IRI"
		if maintenanceAnalysis.Iri2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Iri2)
		}

		filter = append(filter, textFilter)
	}

	if maintenanceAnalysis.Aadt1 != nil || maintenanceAnalysis.Aadt2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Aadt1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Aadt1) + " < "
		}
		textFilter += "AADT"
		if maintenanceAnalysis.Aadt2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Aadt2)
		}

		filter = append(filter, textFilter)
	}

	if maintenanceAnalysis.Ifi1 != nil || maintenanceAnalysis.Ifi2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Ifi1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Ifi1) + " < "
		}
		textFilter += "IFI"
		if maintenanceAnalysis.Ifi2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Ifi2)
		}

		filter = append(filter, textFilter)
	}

	dashBoard.Filter.Filter = filter

	for indexYear := 0; indexYear < year; indexYear++ {
		years = append(years, analystYear+indexYear)
	}

	dashBoard.Graph1.Lable = years
	dashBoard.Bar1.Lable = years
	dashBoard.Bar2.Lable = years

	valueGraph1s := map[string][]float64{}
	valueBar1s := map[string][]float64{}
	keys := []string{}
	for indexPlan := 1; indexPlan < numberPlanWithDefault+1; indexPlan++ {
		var dataset2s models.StrategicDatasetBar2
		key := ""
		if indexPlan > numberPlan {
			if (indexPlan - numberPlan) == 2 {
				key = "ซ่อมบำรุงปกติ"
			} else {
				key = "ไม่จำกัดงบประมาณ"
			}
		} else {
			key = "แผนที่ " + strconv.Itoa(indexPlan)
		}
		dataset2s.Plan = key
		for indexYear := 1; indexYear < year+1; indexYear++ {
			budgetSummaryBar1 := 0.0
			sumGraph1 := 0.0
			sumLengthGraph1 := 0.0
			valueBar2s := map[string][]float64{}
			sumBar2s := map[string]int{}
			keyBar2s := []string{}
			sumBar2 := 0
			for _, item := range modelResults {
				if (indexYear-1) == item.Year && indexPlan == item.PlanSequence {
					if item.Data.Repair {
						sumGraph1 = sumGraph1 + item.Data.RweResult.Iri*float64(item.Data.PrepareDataBefore.Length)
						budgetSummaryBar1 = budgetSummaryBar1 + item.Data.UsedBudget

						valueBar2s[item.Data.IcResult.Name] = append(valueBar2s[item.Data.IcResult.Name], item.Data.UsedBudget)
						if !helpers.Contains(keyBar2s, item.Data.IcResult.Name) {
							keyBar2s = append(keyBar2s, item.Data.IcResult.Name)
						}
						sumBar2s[item.Data.IcResult.Name] = sumBar2s[item.Data.IcResult.Name] + int(item.Data.UsedBudget)
						sumBar2 = sumBar2 + int(item.Data.UsedBudget)
					} else {
						sumGraph1 = sumGraph1 + item.Data.DeteriorationResult.Result.Iri*float64(item.Data.PrepareDataBefore.Length)
					}
					sumLengthGraph1 = sumLengthGraph1 + float64(item.Data.PrepareDataBefore.Length)
				}

			}

			valueGraph1s[key] = append(valueGraph1s[key], helpers.RoundFloat(float64(sumGraph1)/float64(sumLengthGraph1), 2))
			valueBar1s[key] = append(valueBar1s[key], budgetSummaryBar1)

			datset2Lables := []string{}
			datset2Values := []float64{}
			datset2Budgets := []float64{}
			datset2Colors := []string{}
			for _, itemKey := range keyBar2s {
				datset2Lables = append(datset2Lables, itemKey)
				summary := 0.0
				summaryBudget := 0.0
				for _, itemValue := range valueBar2s[itemKey] {
					summary = summary + ((float64(itemValue)/float64(sumBar2)*100)/100)*100
					summaryBudget = summaryBudget + float64(itemValue)
				}
				if math.IsNaN(summary) {
					summary = 0.0
				}

				color := colorMap[itemKey]
				if color == "" {
					color = "#000000"
				}
				datset2Colors = append(datset2Colors, color)
				datset2Values = append(datset2Values, helpers.RoundFloat(summary, 2))
				datset2Budgets = append(datset2Budgets, summaryBudget)
			}

			var dataBar2 models.StrategicDataBar2
			dataBar2.Lable = datset2Lables
			dataBar2.Value = datset2Values
			dataBar2.Budget = datset2Budgets
			dataBar2.Color = datset2Colors
			dataset2s.Data = append(dataset2s.Data, dataBar2)
		}
		keys = append(keys, key)
		if key != "ซ่อมบำรุงปกติ" {
			dashBoard.Bar2.Datasets = append(dashBoard.Bar2.Datasets, dataset2s)
		}
	}

	var strategicDatasetBar2Sort []models.StrategicDatasetBar2
	strategicDatasetBar2Sort = append(strategicDatasetBar2Sort, dashBoard.Bar2.Datasets[len(dashBoard.Bar2.Datasets)-1])
	for i := 0; i < len(dashBoard.Bar2.Datasets); i++ {
		switch dashBoard.Bar2.Datasets[i].Plan {
		case "แผนที่ 1":
			strategicDatasetBar2Sort = append(strategicDatasetBar2Sort, dashBoard.Bar2.Datasets[i])
		case "แผนที่ 2":
			strategicDatasetBar2Sort = append(strategicDatasetBar2Sort, dashBoard.Bar2.Datasets[i])
		case "แผนที่ 3":
			strategicDatasetBar2Sort = append(strategicDatasetBar2Sort, dashBoard.Bar2.Datasets[i])
		}
	}

	dashBoard.Bar2.Datasets = strategicDatasetBar2Sort

	keySort := []string{}
	colorSort := []string{}
	keySort = append(keySort, keys[len(keys)-1])
	keySort = append(keySort, keys[len(keys)-2])
	for i := 0; i < len(keys); i++ {
		switch keys[i] {
		case "แผนที่ 1":
			keySort = append(keySort, keys[i])
		case "แผนที่ 2":
			keySort = append(keySort, keys[i])
		case "แผนที่ 3":
			keySort = append(keySort, keys[i])
		}
		colorSort = append(colorSort, color[i])
	}

	dashBoard.Graph1.Color = colorSort

	dashBoard.Graph1.Line = keySort

	for _, item := range keySort {
		dashBoard.Graph1.Value = append(dashBoard.Graph1.Value, valueGraph1s[item])
	}

	colorSortBar1 := []string{}
	for index, item := range keySort {
		if item != "ซ่อมบำรุงปกติ" {
			colorSortBar1 = append(colorSortBar1, color[index])
		}
	}

	dashBoard.Bar1.Color = colorSortBar1

	var datasets []models.DatasetBar1
	for _, item := range keySort {
		if item != "ซ่อมบำรุงปกติ" {
			var dataset models.DatasetBar1
			dataset.Lable = item
			dataset.Value = valueBar1s[item]
			datasets = append(datasets, dataset)
		}
	}

	dashBoard.Bar1.Datasets = datasets

	var table models.StrategicTable

	titlePlan := []string{"สรุปรวม"}
	for indexPlan := 1; indexPlan < numberPlanWithDefault; indexPlan++ {
		plan := ""
		if indexPlan > numberPlan {
			if (indexPlan - numberPlan) == 2 {
				plan = "ซ่อมบำรุงปกติ"
			} else {
				plan = "ไม่จำกัดงบประมาณ"
			}
		} else {
			plan = "แผนที่ " + strconv.Itoa(indexPlan)
		}

		titlePlan = append(titlePlan, plan)
	}

	for _, item := range titlePlan {

		switch item {
		case "สรุปรวม":
			title := []string{"แผนงบประมาณ (บาท)", "งบประมาณซ่อมบำรุง (บาท)", "IRI (ก่อนวิเคราะห์)", "IRI (หลังวิเคราะห์)"}
			var summaryPlanTables []models.StrategicSummaryPlanTable
			for _, item := range title {
				var summaryPlanTable models.StrategicSummaryPlanTable
				var dataTables []models.StrategicSummaryDataTable
				for indexPlan := 1; indexPlan < numberPlanWithDefault+1; indexPlan++ {
					var dataTable models.StrategicSummaryDataTable
					plan := ""
					if indexPlan > numberPlan {
						if (indexPlan - numberPlan) == 2 {
							plan = "ซ่อมบำรุงปกติ"
						} else {
							plan = "ไม่จำกัดงบประมาณ"
						}
					} else {
						plan = "แผนที่ " + strconv.Itoa(indexPlan)
					}

					values := []float64{}
					for indexYear := 1; indexYear < year+1; indexYear++ {
						switch item {
						case "แผนงบประมาณ (บาท)":
							if indexPlan > numberPlan {
								values = append(values, 0)
							} else {
								values = append(values, float64(float64(maintenanceAnalysisPlanMap[indexYear]["plan_"+strconv.Itoa(indexPlan)])*1000000))
							}
						case "งบประมาณซ่อมบำรุง (บาท)":
							sum := 0.0
							for _, item := range modelResults {
								if (indexYear-1) == item.Year && indexPlan == item.PlanSequence {
									if item.Data.Repair {
										sum = sum + item.Data.UsedBudget
									}
								}
							}
							values = append(values, sum)
						case "IRI (ก่อนวิเคราะห์)":
							sumLength := 0.0
							sumIri := 0.0
							for _, item := range modelResults {
								if (indexYear-1) == item.Year && indexPlan == item.PlanSequence {
									sumIri = sumIri + (item.Data.PrepareData.RoadCondition.Iri * float64(item.Data.PrepareDataBefore.Length))
									sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
								}
							}
							values = append(values, helpers.RoundFloat(sumIri/sumLength, 2))
						case "IRI (หลังวิเคราะห์)":
							sumLength := 0.0
							sumIri := 0.0
							for _, item := range modelResults {
								if (indexYear-1) == item.Year && indexPlan == item.PlanSequence {
									if item.Data.Repair {
										sumIri = sumIri + (item.Data.RweResult.Iri * float64(item.Data.PrepareDataBefore.Length))
										sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
									} else {
										sumIri = sumIri + (item.Data.DeteriorationResult.Result.Iri * float64(item.Data.PrepareDataBefore.Length))
										sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
									}
								}
							}
							values = append(values, helpers.RoundFloat(sumIri/sumLength, 2))
						}
					}
					dataTable.Value = values
					dataTable.Name = plan
					dataTables = append(dataTables, dataTable)
				}

				var dataTablesSort []models.StrategicSummaryDataTable
				dataTablesSort = append(dataTablesSort, dataTables[len(dataTables)-1])
				dataTablesSort = append(dataTablesSort, dataTables[len(dataTables)-2])
				for i := 0; i < len(dataTables); i++ {
					switch dataTables[i].Name {
					case "แผนที่ 1":
						dataTablesSort = append(dataTablesSort, dataTables[i])
					case "แผนที่ 2":
						dataTablesSort = append(dataTablesSort, dataTables[i])
					case "แผนที่ 3":
						dataTablesSort = append(dataTablesSort, dataTables[i])
					}
				}

				summaryPlanTable.Name = item
				summaryPlanTable.Data = dataTablesSort
				summaryPlanTables = append(summaryPlanTables, summaryPlanTable)
			}

			table.Summary = summaryPlanTables

		case "แผนที่ 1":
			titles := []string{}
			for _, item := range modelResults {
				for indexYear := 1; indexYear < year+1; indexYear++ {
					if item.PlanSequence == 1 {
						if !helpers.Contains(titles, item.Data.IcResult.Name) && item.Data.IcResult.Name != "" {
							if item.Data.Repair {
								titles = append(titles, item.Data.IcResult.Name)
							}
						}
					}
				}
			}

			var dataTables []models.StrategicDataTable
			var dataTable models.StrategicDataTable
			var DataTableOthers []models.StrategicDataTableOther
			dataTable.MethodName = "งานบำรุงปกติ"
			for indexYear := 1; indexYear < year+1; indexYear++ {
				sumKm := 0.0
				for _, itemModel := range modelResults {
					if (indexYear-1) == itemModel.Year && itemModel.PlanSequence == 1 {
						if !itemModel.Data.Repair {
							sumKm = float64(sumKm) + math.Abs(float64(itemModel.KmEnd)-float64(itemModel.KmStart))

						}
					}
				}
				var dataTableOther models.StrategicDataTableOther
				dataTableOther.Budget = 0
				dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
				dataTableOther.Year = analystYear + indexYear - 1
				DataTableOthers = append(DataTableOthers, dataTableOther)
			}
			dataTable.Data = DataTableOthers
			dataTables = append(dataTables, dataTable)

			for _, itemTitle := range titles {
				var dataTable models.StrategicDataTable
				var DataTableOthers []models.StrategicDataTableOther
				dataTable.MethodName = itemTitle
				for indexYear := 1; indexYear < year+1; indexYear++ {
					sumBudget := 0.0
					sumKm := 0
					for _, itemModel := range modelResults {
						if itemModel.Data.IcResult.Name == itemTitle && (indexYear-1) == itemModel.Year && itemModel.PlanSequence == 1 {
							if itemModel.Data.Repair {
								sumBudget = sumBudget + itemModel.Data.UsedBudget
								sumKm = sumKm + itemModel.Data.PrepareDataBefore.Length
							}
						}
					}
					var dataTableOther models.StrategicDataTableOther
					dataTableOther.Budget = sumBudget
					dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
					dataTableOther.Year = analystYear + indexYear - 1
					DataTableOthers = append(DataTableOthers, dataTableOther)
				}
				dataTable.Data = DataTableOthers
				dataTables = append(dataTables, dataTable)
			}

			table.Plan1 = dataTables

		case "แผนที่ 2":
			titles := []string{}
			for _, item := range modelResults {
				for indexYear := 1; indexYear < year+1; indexYear++ {
					if item.PlanSequence == 2 {
						if !helpers.Contains(titles, item.Data.IcResult.Name) && item.Data.IcResult.Name != "" {
							if item.Data.Repair {
								titles = append(titles, item.Data.IcResult.Name)
							}
						}
					}
				}
			}

			var dataTables []models.StrategicDataTable
			var dataTable models.StrategicDataTable
			var DataTableOthers []models.StrategicDataTableOther
			dataTable.MethodName = "งานบำรุงปกติ"
			for indexYear := 1; indexYear < year+1; indexYear++ {
				sumKm := 0.0
				for _, itemModel := range modelResults {
					if (indexYear-1) == itemModel.Year && itemModel.PlanSequence == 2 {
						if !itemModel.Data.Repair {
							sumKm = float64(sumKm) + math.Abs(float64(itemModel.KmEnd)-float64(itemModel.KmStart))
						}
					}
				}
				var dataTableOther models.StrategicDataTableOther
				dataTableOther.Budget = 0
				dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
				dataTableOther.Year = analystYear + indexYear - 1
				DataTableOthers = append(DataTableOthers, dataTableOther)
			}
			dataTable.Data = DataTableOthers
			dataTables = append(dataTables, dataTable)

			for _, itemTitle := range titles {
				var dataTable models.StrategicDataTable
				var DataTableOthers []models.StrategicDataTableOther
				dataTable.MethodName = itemTitle
				for indexYear := 1; indexYear < year+1; indexYear++ {
					sumBudget := 0.0
					sumKm := 0
					for _, itemModel := range modelResults {
						if itemModel.Data.IcResult.Name == itemTitle && (indexYear-1) == itemModel.Year && itemModel.PlanSequence == 2 {
							if itemModel.Data.Repair {
								sumBudget = sumBudget + itemModel.Data.UsedBudget
								sumKm = sumKm + itemModel.Data.PrepareDataBefore.Length
							}
						}
					}
					var dataTableOther models.StrategicDataTableOther
					dataTableOther.Budget = sumBudget
					dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
					dataTableOther.Year = analystYear + indexYear - 1
					DataTableOthers = append(DataTableOthers, dataTableOther)
				}
				dataTable.Data = DataTableOthers
				dataTables = append(dataTables, dataTable)
			}

			table.Plan2 = dataTables

		case "แผนที่ 3":
			titles := []string{}
			for _, item := range modelResults {
				for indexYear := 1; indexYear < year+1; indexYear++ {
					if item.PlanSequence == 3 {
						if !helpers.Contains(titles, item.Data.IcResult.Name) && item.Data.IcResult.Name != "" {
							if item.Data.Repair {
								titles = append(titles, item.Data.IcResult.Name)
							}
						}
					}
				}
			}

			var dataTables []models.StrategicDataTable
			var dataTable models.StrategicDataTable
			var DataTableOthers []models.StrategicDataTableOther
			dataTable.MethodName = "งานบำรุงปกติ"
			for indexYear := 1; indexYear < year+1; indexYear++ {
				sumKm := 0.0
				for _, itemModel := range modelResults {
					if (indexYear-1) == itemModel.Year && itemModel.PlanSequence == 3 {
						if !itemModel.Data.Repair {
							sumKm = float64(sumKm) + math.Abs(float64(itemModel.KmEnd)-float64(itemModel.KmStart))
						}
					}
				}
				var dataTableOther models.StrategicDataTableOther
				dataTableOther.Budget = 0
				dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
				dataTableOther.Year = analystYear + indexYear - 1
				DataTableOthers = append(DataTableOthers, dataTableOther)
			}
			dataTable.Data = DataTableOthers
			dataTables = append(dataTables, dataTable)

			for _, itemTitle := range titles {
				var dataTable models.StrategicDataTable
				var DataTableOthers []models.StrategicDataTableOther
				dataTable.MethodName = itemTitle
				for indexYear := 1; indexYear < year+1; indexYear++ {
					sumBudget := 0.0
					sumKm := 0
					for _, itemModel := range modelResults {
						if itemModel.Data.IcResult.Name == itemTitle && (indexYear-1) == itemModel.Year && itemModel.PlanSequence == 3 {
							if itemModel.Data.Repair {
								sumBudget = sumBudget + itemModel.Data.UsedBudget
								sumKm = sumKm + itemModel.Data.PrepareDataBefore.Length
							}
						}
					}
					var dataTableOther models.StrategicDataTableOther
					dataTableOther.Budget = sumBudget
					dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
					dataTableOther.Year = analystYear + indexYear - 1
					DataTableOthers = append(DataTableOthers, dataTableOther)
				}
				dataTable.Data = DataTableOthers
				dataTables = append(dataTables, dataTable)
			}

			table.Plan3 = dataTables

		case "ไม่จำกัดงบประมาณ":
			titles := []string{}
			for _, item := range modelResults {
				for indexYear := 1; indexYear < year+1; indexYear++ {
					if (item.PlanSequence - numberPlan) == 1 {
						if !helpers.Contains(titles, item.Data.IcResult.Name) && item.Data.IcResult.Name != "" {
							if item.Data.Repair {
								titles = append(titles, item.Data.IcResult.Name)
							}
						}
					}
				}
			}

			var dataTables []models.StrategicDataTable
			var dataTable models.StrategicDataTable
			var DataTableOthers []models.StrategicDataTableOther
			dataTable.MethodName = "งานบำรุงปกติ"
			for indexYear := 1; indexYear < year+1; indexYear++ {
				sumKm := 0.0
				for _, itemModel := range modelResults {
					if (itemModel.PlanSequence - numberPlan) == 1 {
						if (indexYear - 1) == itemModel.Year {
							if !itemModel.Data.Repair {
								sumKm = float64(sumKm) + math.Abs(float64(itemModel.KmEnd)-float64(itemModel.KmStart))
							}
						}
					}
				}
				var dataTableOther models.StrategicDataTableOther
				dataTableOther.Budget = 0
				dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
				dataTableOther.Year = analystYear + indexYear - 1
				DataTableOthers = append(DataTableOthers, dataTableOther)
			}
			dataTable.Data = DataTableOthers
			dataTables = append(dataTables, dataTable)

			for _, itemTitle := range titles {
				var dataTable models.StrategicDataTable
				var DataTableOthers []models.StrategicDataTableOther
				dataTable.MethodName = itemTitle
				for indexYear := 1; indexYear < year+1; indexYear++ {
					sumBudget := 0.0
					sumKm := 0
					for _, itemModel := range modelResults {
						if itemModel.Data.IcResult.Name == itemTitle && (indexYear-1) == itemModel.Year && (itemModel.PlanSequence-numberPlan) == 1 {
							if itemModel.Data.Repair {
								sumBudget = sumBudget + itemModel.Data.UsedBudget
								sumKm = sumKm + itemModel.Data.PrepareDataBefore.Length
							}
						}
					}
					var dataTableOther models.StrategicDataTableOther
					dataTableOther.Budget = sumBudget
					dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
					dataTableOther.Year = analystYear + indexYear - 1
					DataTableOthers = append(DataTableOthers, dataTableOther)
				}
				dataTable.Data = DataTableOthers
				dataTables = append(dataTables, dataTable)
			}
			table.UnlimitedPlan = dataTables
		}
	}

	if len(table.Plan1) == 0 {
		table.Plan1 = []models.StrategicDataTable{}
	}

	if len(table.Plan2) == 0 {
		table.Plan2 = []models.StrategicDataTable{}
	}

	if len(table.Plan3) == 0 {
		table.Plan3 = []models.StrategicDataTable{}
	}

	if len(table.UnlimitedPlan) == 0 {
		table.UnlimitedPlan = []models.StrategicDataTable{}
	}

	dashBoard.Table = table

	return dashBoard, nil
}

func (u *UseCase) strategicMaintenanceAnalysisTargetIri(idInteger int, modelResults []models.ModelResult, latestPrepareData models.PrepareData, maintenanceAnalysis models.MaintenanceAnalysis, maintenanceAnalysisPlan []models.MaintenanceAnalysisPlan, color []string, colorMap map[string]string) (interface{}, error) {

	maintenanceAnalysisStrategicBudgetType, err := u.Repo.GetMaintenanceAnalysisStrategicBudgetTypeById(maintenanceAnalysis.Condition)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysisStrategiTarget, err := u.Repo.GetMaintenanceAnalysisStrategiTargetById(*maintenanceAnalysis.Target)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysisRoad, err := u.Repo.GetMaintenanceAnalysisRoadNameById(maintenanceAnalysis.ID)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysisRoads := []string{}
	for _, item := range maintenanceAnalysisRoad {
		maintenanceAnalysisRoads = append(maintenanceAnalysisRoads, item.Name)
	}

	min := 9999
	for _, item := range modelResults {
		if item.AnalystYear <= min {
			min = item.AnalystYear
		}
	}

	year := *maintenanceAnalysis.Year
	analystYear := min

	numberPlanWithDefault := *maintenanceAnalysis.NumberPlan
	numberPlan := *maintenanceAnalysis.NumberPlan
	years := []int{}
	lane := ""

	maintenanceAnalysisPlanMap := make(map[int]map[string]float64)
	for _, item := range maintenanceAnalysisPlan {
		maintenanceAnalysisPlanMapSub := map[string]float64{}
		for indexPlan := 1; indexPlan < *maintenanceAnalysis.NumberPlan+1; indexPlan++ {
			switch indexPlan {
			case 1:
				maintenanceAnalysisPlanMapSub["plan_1"] = *item.Plan1
			case 2:
				maintenanceAnalysisPlanMapSub["plan_2"] = *item.Plan2
			case 3:
				maintenanceAnalysisPlanMapSub["plan_3"] = *item.Plan3
			}
		}
		maintenanceAnalysisPlanMap[int(item.PlanYear)] = maintenanceAnalysisPlanMapSub
	}

	switch maintenanceAnalysis.LaneTypeId {
	case 0:
		lane = "ทั้งหมด"
	case 1:
		lane = "1"
	case 2:
		lane = "2"
	case 3:
		lane = "3"
	}

	surface := ""
	switch maintenanceAnalysis.SurfaceTypeId {
	case 1:
		surface = "ลาดยาง"
	case 2:
		surface = "คอนกรีต"
	}

	comment := ""
	if maintenanceAnalysis.Comment == "" {
		comment = "-"
	} else {
		comment = maintenanceAnalysis.Comment
	}

	var dashBoard models.DashboardStrategicMaintenance
	dashBoard.Road = maintenanceAnalysisRoads
	dashBoard.Filter.Lane = lane
	dashBoard.Filter.SurfaceType = surface
	dashBoard.Comment = comment
	dashBoard.Filter.Km = maintenanceAnalysis.GroupKm
	dashBoard.Condition.Condition = maintenanceAnalysisStrategicBudgetType.Name
	dashBoard.Condition.Target = maintenanceAnalysisStrategiTarget.Name
	dashBoard.Condition.Discount = maintenanceAnalysis.Discount
	dashBoard.Graph1.Name = "ข้อมูลประจำปี"
	dashBoard.Bar1.Name = "ข้อมูลประจำปี"
	dashBoard.Bar2.Name = "วิธีการซ่อม"
	dashBoard.NumberPlan = *maintenanceAnalysis.NumberPlan

	var filter = []string{}
	if maintenanceAnalysis.Iri1 != nil || maintenanceAnalysis.Iri2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Iri1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Iri1) + " < "
		}
		textFilter += "IRI"
		if maintenanceAnalysis.Iri2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Iri2)
		}

		filter = append(filter, textFilter)
	}

	if maintenanceAnalysis.Aadt1 != nil || maintenanceAnalysis.Aadt2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Aadt1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Aadt1) + " < "
		}
		textFilter += "AADT"
		if maintenanceAnalysis.Aadt2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Aadt2)
		}

		filter = append(filter, textFilter)
	}

	if maintenanceAnalysis.Ifi1 != nil || maintenanceAnalysis.Ifi2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Ifi1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Ifi1) + " < "
		}
		textFilter += "IFI"
		if maintenanceAnalysis.Ifi2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Ifi2)
		}

		filter = append(filter, textFilter)
	}

	dashBoard.Filter.Filter = filter

	for indexYear := 0; indexYear < year; indexYear++ {
		years = append(years, analystYear+indexYear)
	}

	dashBoard.Graph1.Lable = years
	dashBoard.Bar1.Lable = years
	dashBoard.Bar2.Lable = years

	valueGraph1s := map[string][]float64{}
	valueBar1s := map[string][]float64{}
	keys := []string{}
	for indexPlan := 1; indexPlan < numberPlanWithDefault+1; indexPlan++ {
		var dataset2s models.StrategicDatasetBar2
		key := "แผนที่ " + strconv.Itoa(indexPlan)
		dataset2s.Plan = key
		for indexYear := 1; indexYear < year+1; indexYear++ {
			budgetSummaryBar1 := 0.0
			sumGraph1 := 0.0
			sumLengthGraph1 := 0.0
			valueBar2s := map[string][]float64{}
			sumBar2s := map[string]int{}
			keyBar2s := []string{}
			sumBar2 := 0
			for _, item := range modelResults {
				if (indexYear-1) == item.Year && indexPlan == item.PlanSequence {
					if item.Data.Repair {
						sumGraph1 = sumGraph1 + item.Data.RweResult.Iri*float64(item.Data.PrepareDataBefore.Length)
						budgetSummaryBar1 = budgetSummaryBar1 + item.Data.UsedBudget

						valueBar2s[item.Data.IcResult.Name] = append(valueBar2s[item.Data.IcResult.Name], item.Data.UsedBudget)
						if !helpers.Contains(keyBar2s, item.Data.IcResult.Name) {
							keyBar2s = append(keyBar2s, item.Data.IcResult.Name)
						}
						sumBar2s[item.Data.IcResult.Name] = sumBar2s[item.Data.IcResult.Name] + int(item.Data.UsedBudget)
						sumBar2 = sumBar2 + int(item.Data.UsedBudget)
					} else {
						sumGraph1 = sumGraph1 + item.Data.DeteriorationResult.Result.Iri*float64(item.Data.PrepareDataBefore.Length)
					}
					sumLengthGraph1 = sumLengthGraph1 + float64(item.Data.PrepareDataBefore.Length)
				}
			}

			valueGraph1s[key] = append(valueGraph1s[key], helpers.RoundFloat(float64(sumGraph1)/float64(sumLengthGraph1), 2))
			valueBar1s[key] = append(valueBar1s[key], budgetSummaryBar1)

			datset2Lables := []string{}
			datset2Values := []float64{}
			datset2Budgets := []float64{}
			datset2Colors := []string{}
			for _, itemKey := range keyBar2s {
				datset2Lables = append(datset2Lables, itemKey)
				summary := 0.0
				summaryBudget := 0.0
				for _, itemValue := range valueBar2s[itemKey] {
					summary = summary + ((float64(itemValue)/float64(sumBar2)*100)/100)*100
					summaryBudget = summaryBudget + float64(itemValue)
				}
				if math.IsNaN(summary) {
					summary = 0.0
				}

				color := colorMap[itemKey]
				if color == "" {
					color = "#000000"
				}
				datset2Colors = append(datset2Colors, color)
				datset2Values = append(datset2Values, helpers.RoundFloat(summary, 2))
				datset2Budgets = append(datset2Budgets, summaryBudget)
			}

			var dataBar2 models.StrategicDataBar2
			dataBar2.Lable = datset2Lables
			dataBar2.Value = datset2Values
			dataBar2.Budget = datset2Budgets
			dataBar2.Color = datset2Colors
			dataset2s.Data = append(dataset2s.Data, dataBar2)
		}
		keys = append(keys, key)
		dashBoard.Bar2.Datasets = append(dashBoard.Bar2.Datasets, dataset2s)
	}

	var strategicDatasetBar2Sort []models.StrategicDatasetBar2
	for i := 0; i < len(dashBoard.Bar2.Datasets); i++ {
		switch dashBoard.Bar2.Datasets[i].Plan {
		case "แผนที่ 1":
			strategicDatasetBar2Sort = append(strategicDatasetBar2Sort, dashBoard.Bar2.Datasets[i])
		case "แผนที่ 2":
			strategicDatasetBar2Sort = append(strategicDatasetBar2Sort, dashBoard.Bar2.Datasets[i])
		case "แผนที่ 3":
			strategicDatasetBar2Sort = append(strategicDatasetBar2Sort, dashBoard.Bar2.Datasets[i])
		}
	}

	dashBoard.Bar2.Datasets = strategicDatasetBar2Sort

	var table models.StrategicTable

	colorSort := []string{}
	keySort := []string{}
	for i := 0; i < len(keys); i++ {
		switch keys[i] {
		case "แผนที่ 1":
			keySort = append(keySort, keys[i])
		case "แผนที่ 2":
			keySort = append(keySort, keys[i])
		case "แผนที่ 3":
			keySort = append(keySort, keys[i])
		}
		colorSort = append(colorSort, color[i])
	}

	dashBoard.Graph1.Color = colorSort

	dashBoard.Graph1.Line = keySort

	for _, item := range keySort {
		dashBoard.Graph1.Value = append(dashBoard.Graph1.Value, valueGraph1s[item])
	}

	dashBoard.Bar1.Color = colorSort

	var datasets []models.DatasetBar1
	for _, item := range keySort {
		var dataset models.DatasetBar1
		dataset.Lable = item
		dataset.Value = valueBar1s[item]
		datasets = append(datasets, dataset)
	}

	dashBoard.Bar1.Datasets = datasets

	titlePlan := []string{"สรุปรวม"}
	for indexPlan := 1; indexPlan < numberPlanWithDefault+1; indexPlan++ {
		plan := "แผนที่ " + strconv.Itoa(indexPlan)
		titlePlan = append(titlePlan, plan)
	}

	for _, item := range titlePlan {
		switch item {
		case "สรุปรวม":
			title := []string{"IRI เป้าหมาย", "งบประมาณซ่อมบำรุง (บาท)", "IRI (ก่อนวิเคราะห์)", "IRI (หลังวิเคราะห์)"}
			var summaryPlanTables []models.StrategicSummaryPlanTable
			for _, item := range title {
				var summaryPlanTable models.StrategicSummaryPlanTable
				var dataTables []models.StrategicSummaryDataTable
				for indexPlan := 1; indexPlan < numberPlanWithDefault+1; indexPlan++ {
					var dataTable models.StrategicSummaryDataTable
					plan := "แผนที่ " + strconv.Itoa(indexPlan)
					values := []float64{}
					for indexYear := 1; indexYear < year+1; indexYear++ {
						switch item {
						case "IRI เป้าหมาย":
							if indexPlan > numberPlan {
								values = append(values, 0)
							} else {
								values = append(values, float64(float64(maintenanceAnalysisPlanMap[indexYear]["plan_"+strconv.Itoa(indexPlan)])))
							}
						case "งบประมาณซ่อมบำรุง (บาท)":
							sum := 0.0
							for _, item := range modelResults {
								if (indexYear-1) == item.Year && indexPlan == item.PlanSequence {
									if item.Data.Repair {
										sum = sum + item.Data.UsedBudget
									}
								}
							}
							values = append(values, sum)
						case "IRI (ก่อนวิเคราะห์)":
							sumLength := 0.0
							sumIri := 0.0
							for _, item := range modelResults {
								if (indexYear-1) == item.Year && indexPlan == item.PlanSequence {
									sumIri = sumIri + (item.Data.PrepareData.RoadCondition.Iri * float64(item.Data.PrepareDataBefore.Length))
									sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
								}
							}
							values = append(values, helpers.RoundFloat(sumIri/sumLength, 2))
						case "IRI (หลังวิเคราะห์)":
							sumLength := 0.0
							sumIri := 0.0
							for _, item := range modelResults {
								if (indexYear-1) == item.Year && indexPlan == item.PlanSequence {
									if item.Data.Repair {
										sumIri = sumIri + (item.Data.RweResult.Iri * float64(item.Data.PrepareDataBefore.Length))
										sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
									} else {
										sumIri = sumIri + (item.Data.DeteriorationResult.Result.Iri * float64(item.Data.PrepareDataBefore.Length))
										sumLength = sumLength + float64(item.Data.PrepareDataBefore.Length)
									}
								}
							}
							values = append(values, helpers.RoundFloat(sumIri/sumLength, 2))
						}
					}
					dataTable.Value = values
					dataTable.Name = plan
					dataTables = append(dataTables, dataTable)
				}

				var dataTablesSort []models.StrategicSummaryDataTable
				for i := 0; i < len(dataTables); i++ {
					switch dataTables[i].Name {
					case "แผนที่ 1":
						dataTablesSort = append(dataTablesSort, dataTables[i])
					case "แผนที่ 2":
						dataTablesSort = append(dataTablesSort, dataTables[i])
					case "แผนที่ 3":
						dataTablesSort = append(dataTablesSort, dataTables[i])
					}
				}

				summaryPlanTable.Name = item
				summaryPlanTable.Data = dataTablesSort
				summaryPlanTables = append(summaryPlanTables, summaryPlanTable)
			}

			table.Summary = summaryPlanTables
		case "แผนที่ 1":
			titles := []string{}
			for _, item := range modelResults {
				for indexYear := 1; indexYear < year+1; indexYear++ {
					if item.PlanSequence == 1 {
						if !helpers.Contains(titles, item.Data.IcResult.Name) && item.Data.IcResult.Name != "" {
							if item.Data.Repair {
								titles = append(titles, item.Data.IcResult.Name)
							}
						}
					}
				}
			}

			var dataTables []models.StrategicDataTable
			var dataTable models.StrategicDataTable
			var DataTableOthers []models.StrategicDataTableOther
			dataTable.MethodName = "งานบำรุงปกติ"
			for indexYear := 1; indexYear < year+1; indexYear++ {
				sumKm := 0.0
				for _, itemModel := range modelResults {
					if (indexYear-1) == itemModel.Year && itemModel.PlanSequence == 1 {
						if !itemModel.Data.Repair {
							sumKm = float64(sumKm) + math.Abs(float64(itemModel.KmEnd)-float64(itemModel.KmStart))

						}
					}
				}
				var dataTableOther models.StrategicDataTableOther
				dataTableOther.Budget = 0
				dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
				dataTableOther.Year = analystYear + indexYear - 1
				DataTableOthers = append(DataTableOthers, dataTableOther)
			}
			dataTable.Data = DataTableOthers
			dataTables = append(dataTables, dataTable)

			for _, itemTitle := range titles {
				var dataTable models.StrategicDataTable
				var DataTableOthers []models.StrategicDataTableOther
				dataTable.MethodName = itemTitle
				for indexYear := 1; indexYear < year+1; indexYear++ {
					sumBudget := 0.0
					sumKm := 0
					for _, itemModel := range modelResults {
						if itemModel.Data.IcResult.Name == itemTitle && (indexYear-1) == itemModel.Year && itemModel.PlanSequence == 1 {
							if itemModel.Data.Repair {
								sumBudget = sumBudget + itemModel.Data.UsedBudget
								sumKm = sumKm + itemModel.Data.PrepareDataBefore.Length
							}
						}
					}
					var dataTableOther models.StrategicDataTableOther
					dataTableOther.Budget = sumBudget
					dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
					dataTableOther.Year = analystYear + indexYear - 1
					DataTableOthers = append(DataTableOthers, dataTableOther)
				}
				dataTable.Data = DataTableOthers
				dataTables = append(dataTables, dataTable)
			}

			table.Plan1 = dataTables
		case "แผนที่ 2":
			titles := []string{}
			for _, item := range modelResults {
				for indexYear := 1; indexYear < year+1; indexYear++ {
					if item.PlanSequence == 2 {
						if !helpers.Contains(titles, item.Data.IcResult.Name) && item.Data.IcResult.Name != "" {
							if item.Data.Repair {
								titles = append(titles, item.Data.IcResult.Name)
							}
						}
					}
				}
			}

			var dataTables []models.StrategicDataTable
			var dataTable models.StrategicDataTable
			var DataTableOthers []models.StrategicDataTableOther
			dataTable.MethodName = "งานบำรุงปกติ"
			for indexYear := 1; indexYear < year+1; indexYear++ {
				sumKm := 0.0
				for _, itemModel := range modelResults {
					if (indexYear-1) == itemModel.Year && itemModel.PlanSequence == 2 {
						if !itemModel.Data.Repair {
							sumKm = float64(sumKm) + math.Abs(float64(itemModel.KmEnd)-float64(itemModel.KmStart))
						}
					}
				}
				var dataTableOther models.StrategicDataTableOther
				dataTableOther.Budget = 0
				dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
				dataTableOther.Year = analystYear + indexYear - 1
				DataTableOthers = append(DataTableOthers, dataTableOther)
			}
			dataTable.Data = DataTableOthers
			dataTables = append(dataTables, dataTable)

			for _, itemTitle := range titles {
				var dataTable models.StrategicDataTable
				var DataTableOthers []models.StrategicDataTableOther
				dataTable.MethodName = itemTitle
				for indexYear := 1; indexYear < year+1; indexYear++ {
					sumBudget := 0.0
					sumKm := 0
					for _, itemModel := range modelResults {
						if itemModel.Data.IcResult.Name == itemTitle && (indexYear-1) == itemModel.Year && itemModel.PlanSequence == 2 {
							if itemModel.Data.Repair {
								sumBudget = sumBudget + itemModel.Data.UsedBudget
								sumKm = sumKm + itemModel.Data.PrepareDataBefore.Length
							}
						}
					}
					var dataTableOther models.StrategicDataTableOther
					dataTableOther.Budget = sumBudget
					dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
					dataTableOther.Year = analystYear + indexYear - 1
					DataTableOthers = append(DataTableOthers, dataTableOther)
				}
				dataTable.Data = DataTableOthers
				dataTables = append(dataTables, dataTable)
			}

			table.Plan2 = dataTables
		case "แผนที่ 3":
			titles := []string{}
			for _, item := range modelResults {
				for indexYear := 1; indexYear < year+1; indexYear++ {
					if item.PlanSequence == 3 {
						if !helpers.Contains(titles, item.Data.IcResult.Name) && item.Data.IcResult.Name != "" {
							if item.Data.Repair {
								titles = append(titles, item.Data.IcResult.Name)
							}
						}
					}
				}
			}

			var dataTables []models.StrategicDataTable
			var dataTable models.StrategicDataTable
			var DataTableOthers []models.StrategicDataTableOther
			dataTable.MethodName = "งานบำรุงปกติ"
			for indexYear := 1; indexYear < year+1; indexYear++ {
				sumKm := 0.0
				for _, itemModel := range modelResults {
					if (indexYear-1) == itemModel.Year && itemModel.PlanSequence == 3 {
						if !itemModel.Data.Repair {
							sumKm = float64(sumKm) + math.Abs(float64(itemModel.KmEnd)-float64(itemModel.KmStart))
						}
					}
				}
				var dataTableOther models.StrategicDataTableOther
				dataTableOther.Budget = 0
				dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
				dataTableOther.Year = analystYear + indexYear - 1
				DataTableOthers = append(DataTableOthers, dataTableOther)
			}
			dataTable.Data = DataTableOthers
			dataTables = append(dataTables, dataTable)

			for _, itemTitle := range titles {
				var dataTable models.StrategicDataTable
				var DataTableOthers []models.StrategicDataTableOther
				dataTable.MethodName = itemTitle
				for indexYear := 1; indexYear < year+1; indexYear++ {
					sumBudget := 0.0
					sumKm := 0
					for _, itemModel := range modelResults {
						if itemModel.Data.IcResult.Name == itemTitle && (indexYear-1) == itemModel.Year && itemModel.PlanSequence == 3 {
							if itemModel.Data.Repair {
								sumBudget = sumBudget + itemModel.Data.UsedBudget
								sumKm = sumKm + itemModel.Data.PrepareDataBefore.Length
							}
						}
					}
					var dataTableOther models.StrategicDataTableOther
					dataTableOther.Budget = sumBudget
					dataTableOther.Km = math.Abs(float64(sumKm) / float64(1000))
					dataTableOther.Year = analystYear + indexYear - 1
					DataTableOthers = append(DataTableOthers, dataTableOther)
				}
				dataTable.Data = DataTableOthers
				dataTables = append(dataTables, dataTable)
			}

			table.Plan3 = dataTables
		case "ไม่จำกัดงบประมาณ":
		}
	}

	if len(table.Plan1) == 0 {
		table.Plan1 = []models.StrategicDataTable{}
	}

	if len(table.Plan2) == 0 {
		table.Plan2 = []models.StrategicDataTable{}
	}

	if len(table.Plan3) == 0 {
		table.Plan3 = []models.StrategicDataTable{}
	}

	if len(table.UnlimitedPlan) == 0 {
		table.UnlimitedPlan = []models.StrategicDataTable{}
	}

	dashBoard.Table = table

	return dashBoard, nil
}

func (u *UseCase) DashboardAnnualMaintenanceAnalysis(id string, userID int) (interface{}, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}
	// ============ start check permission ============
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
	depot, err := u.Repo.GetRefDepotByRoad()
	if err != nil {
		log.Println(err)
		if err.Error() == "record not found" {
			return "", responses.NewNotFoundError()
		}
		return "", responses.NewAppErr(400, err.Error())
	}

	roads, err := u.Repo.GetMaintenanceAnalysisRoadByID(ID)
	if err != nil {
		logs.Error(err)
		return responses.ResponseMaintenanceAnalysis{}, responses.NewAppErr(400, err.Error())
	}
	isShow := true
	if isOwnerData && !isAllData {
		if len(roads) == 0 {
			isShow = false
		}
		for _, prepare := range roads {
			roadID := prepare.RoadID
			val, isVal := depot[roadID]
			if !isVal {
				isShow = false
			} else {
				isShow = false
				if val == depotCode {
					isShow = true
				}

				goto EXITLOOP
			}
		}
	}
EXITLOOP:

	if !isOwnerData && !isAllData {
		isShow = false
	}

	if !isShow {
		return responses.ResponseMaintenanceAnalysis{}, responses.NewAppErr(403, "")
	}
	idInteger, err := strconv.Atoi(id)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	modelResults, err := u.Repo.GetModelResultDataById(idInteger)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	latestPrepareData, err := u.Repo.GetLatestPrepareDataByMaintenanceAnalysisId(idInteger)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysis, err := u.Repo.GetMaintenanceAnalysisById(idInteger)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	maintenanceAnalysisPlan, err := u.Repo.GetMaintenanceAnalysisPlanByAnalysisID(idInteger)
	if err != nil {
		logs.Error(err)
		return models.DashboardStrategicMaintenance{}, err
	}

	interventionCriteria, err := u.Repo.GetInterventionCriteria()
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return models.DashboardStrategicMaintenance{}, err
	}

	colorMap := make(map[string]string)
	for index, item := range interventionCriteria {
		colorMap[item.MaintenanceStandardName] = colors[index]
	}

	colorMap["ซ่อมบำรุงปกติ"] = "#000000"

	switch maintenanceAnalysis.Condition {
	case 1:
		result, err := u.annualMaintenanceAnalysisUnlimitBudget(idInteger, modelResults, latestPrepareData, maintenanceAnalysis, maintenanceAnalysisPlan, colorMap)
		if err != nil {
			logs.Error(err)
			return []string{}, err
		}
		return result, nil

	case 2:
		result, err := u.annualMaintenanceAnalysisLimitBudget(idInteger, modelResults, latestPrepareData, maintenanceAnalysis, maintenanceAnalysisPlan, colorMap)
		if err != nil {
			logs.Error(err)
			return []string{}, err
		}
		return result, nil

	case 3:
		result, err := u.annualMaintenanceAnalysisTargetIri(idInteger, modelResults, latestPrepareData, maintenanceAnalysis, maintenanceAnalysisPlan, colorMap)
		if err != nil {
			logs.Error(err)
			return []string{}, err
		}
		return result, nil
	}

	return []string{}, nil
}

func (u *UseCase) annualMaintenanceAnalysisUnlimitBudget(idInteger int, modelResults []models.ModelResult, latestPrepareData models.PrepareData, maintenanceAnalysis models.MaintenanceAnalysis, maintenanceAnalysisPlan []models.MaintenanceAnalysisPlan, colorMap map[string]string) (interface{}, error) {

	maintenanceAnalysisStrategicBudgetType, err := u.Repo.GetMaintenanceAnalysisStrategicBudgetTypeById(maintenanceAnalysis.Condition)
	if err != nil {
		logs.Error(err)
		return models.DashboardAnnualMaintenance{}, err
	}

	maintenanceAnalysisStrategiTarget, err := u.Repo.GetMaintenanceAnalysisStrategiTargetById(*maintenanceAnalysis.Target)
	if err != nil {
		logs.Error(err)
		return models.DashboardAnnualMaintenance{}, err
	}

	maintenanceAnalysisRoad, err := u.Repo.GetMaintenanceAnalysisRoadNameById(maintenanceAnalysis.ID)
	if err != nil {
		logs.Error(err)
		return models.DashboardAnnualMaintenance{}, err
	}

	maintenanceAnalysisRoads := []string{}
	for _, item := range maintenanceAnalysisRoad {
		maintenanceAnalysisRoads = append(maintenanceAnalysisRoads, item.Name)
	}
	lane := ""
	switch maintenanceAnalysis.LaneTypeId {
	case 0:
		lane = "ทั้งหมด"
	case 1:
		lane = "1"
	case 2:
		lane = "2"
	case 3:
		lane = "3"
	}

	surface := ""
	switch maintenanceAnalysis.SurfaceTypeId {
	case 1:
		surface = "ลาดยาง"
	case 2:
		surface = "คอนกรีต"
	}

	comment := ""
	if maintenanceAnalysis.Comment == "" {
		comment = "-"
	} else {
		comment = maintenanceAnalysis.Comment
	}

	var dashBoard models.DashboardAnnualMaintenance
	dashBoard.Road = maintenanceAnalysisRoads
	dashBoard.Filter.Lane = lane
	dashBoard.Filter.SurfaceType = surface
	dashBoard.Comment = comment
	dashBoard.Filter.Km = maintenanceAnalysis.GroupKm
	dashBoard.Condition.Condition = maintenanceAnalysisStrategicBudgetType.Name
	dashBoard.Condition.Target = maintenanceAnalysisStrategiTarget.Name
	dashBoard.Condition.Discount = maintenanceAnalysis.Discount

	var filter = []string{}
	if maintenanceAnalysis.Iri1 != nil || maintenanceAnalysis.Iri2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Iri1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Iri1) + " < "
		}
		textFilter += "IRI"
		if maintenanceAnalysis.Iri2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Iri2)
		}

		filter = append(filter, textFilter)
	}

	if maintenanceAnalysis.Aadt1 != nil || maintenanceAnalysis.Aadt2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Aadt1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Aadt1) + " < "
		}
		textFilter += "AADT"
		if maintenanceAnalysis.Aadt2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Aadt2)
		}

		filter = append(filter, textFilter)
	}

	if maintenanceAnalysis.Ifi1 != nil || maintenanceAnalysis.Ifi2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Ifi1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Ifi1) + " < "
		}
		textFilter += "IFI"
		if maintenanceAnalysis.Ifi2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Ifi2)
		}

		filter = append(filter, textFilter)
	}

	dashBoard.Filter.Filter = filter

	values := map[string][]int{}
	valuesAera := map[string][]float64{}
	valuesKm := map[string][]int{}
	countIri := 0
	iriBefore := 0.0
	iriAfter := 0.0
	sums := map[string]int{}
	keys := []string{}
	keys = append(keys, "ซ่อมบำรุงปกติ")
	sum := 0
	sumArea := 0.0
	length := 0.0

	for _, item := range modelResults {
		if item.Data.Repair {
			valuesAera[item.Data.IcResult.Name] = append(valuesAera[item.Data.IcResult.Name], item.Data.PrepareDataBefore.Area)
			values[item.Data.IcResult.Name] = append(values[item.Data.IcResult.Name], int(item.Data.UsedBudget))
			if !helpers.Contains(keys, item.Data.IcResult.Name) {
				keys = append(keys, item.Data.IcResult.Name)
			}

			valuesKm[item.Data.IcResult.Name] = append(valuesKm[item.Data.IcResult.Name], int(math.Abs(float64(item.KmStart)-float64(item.KmEnd))))
			sums[item.Data.IcResult.Name] = sums[item.Data.IcResult.Name] + int(item.Data.UsedBudget)
			sum = sum + int(item.Data.UsedBudget)

			iriAfter += item.Data.RweResult.Iri * float64(item.Data.PrepareDataBefore.Length)
		} else {
			valuesAera["ซ่อมบำรุงปกติ"] = append(valuesAera["ซ่อมบำรุงปกติ"], item.Data.PrepareDataBefore.Area)
			valuesKm["ซ่อมบำรุงปกติ"] = append(valuesKm["ซ่อมบำรุงปกติ"], int(math.Abs(float64(item.KmStart)-float64(item.KmEnd))))

			iriAfter += item.Data.DeteriorationResult.Result.Iri * float64(item.Data.PrepareDataBefore.Length)
		}
		length += float64(item.Data.PrepareDataBefore.Length)

		sumArea = sumArea + item.Data.PrepareDataBefore.Area
		iriBefore += item.Data.PrepareData.RoadCondition.Iri * float64(item.Data.PrepareDataBefore.Length)
		countIri++
	}

	datsetsLable := []string{}
	datsetsValues := []float64{}
	datsetsArea := []float64{}
	datsetsColor := []string{}
	for _, itemKey := range keys {
		datsetsLable = append(datsetsLable, itemKey)
		summary := 0.0
		summaryArea := 0.0
		for _, itemValue := range valuesAera[itemKey] {
			summary = summary + ((float64(itemValue)/float64(sumArea)*100)/100)*100
			summaryArea = summaryArea + float64(itemValue)
		}
		if math.IsNaN(summary) {
			summary = 0.0
		}
		datsetsArea = append(datsetsArea, summaryArea)
		datsetsValues = append(datsetsValues, helpers.RoundFloat(summary, 2))

		color := colorMap[itemKey]
		if color == "" {
			color = "#000000"
		}

		datsetsColor = append(datsetsColor, color)
	}

	var bar models.AnnualBar
	bar.Name = "ปริมาณงาน"
	bar.Lable = datsetsLable
	bar.Data = datsetsValues
	bar.Area = datsetsArea
	bar.Color = datsetsColor
	dashBoard.Bar1 = bar

	var bar2 models.AnnualBar2
	datsets2Lable := []string{}
	datsets2Values := []float64{}
	datsets2Budget := []float64{}
	datsets2Color := []string{}

	for _, itemKey := range keys {
		datsets2Lable = append(datsets2Lable, itemKey)
		summary := 0.0
		summaryBudget := 0.0
		for _, itemValue := range values[itemKey] {
			summary = summary + ((float64(itemValue)/float64(sum)*100)/100)*100
			summaryBudget = summaryBudget + float64(itemValue)
		}
		if math.IsNaN(summary) {
			summary = 0.0
		}
		datsets2Budget = append(datsets2Budget, summaryBudget)
		datsets2Values = append(datsets2Values, helpers.RoundFloat(summary, 2))

		color := colorMap[itemKey]
		if color == "" {
			color = "#000000"
		}

		datsets2Color = append(datsets2Color, color)
	}

	bar2.Name = "ค่าซ่อมบำรุง"
	bar2.Lable = datsets2Lable
	bar2.Data = datsets2Values
	bar2.Budget = datsets2Budget
	bar2.Color = datsets2Color
	dashBoard.Bar2 = bar2

	var annualTableData1 models.AnnualTableData1
	annualTableData1.Budget = 0

	iriAfter = iriAfter / length
	if math.IsNaN(iriAfter) {
		iriAfter = 0
	}

	iriBefore = iriBefore / length
	if math.IsNaN(iriBefore) {
		iriBefore = 0
	}

	annualTableData1.IriAfter = iriAfter
	annualTableData1.IriBefore = iriBefore

	dashBoard.Table.Table1 = annualTableData1

	var annualTableData2s []models.AnnualTableData2

	summaryArea := 0.0
	summaryBudget := 0.0
	summaryKm := 0.0
	for _, item := range modelResults {
		if !item.Data.Repair {
			summaryArea = summaryArea + float64(item.Data.PrepareDataBefore.Area)
			summaryBudget = summaryBudget + item.Data.UsedBudget
			summaryKm = summaryKm + float64(item.Data.PrepareDataBefore.Length)
		}
	}

	for _, itemKey := range keys {
		var annualTableData2 models.AnnualTableData2
		annualTableData2.Name = itemKey
		summaryArea := 0.0
		for _, itemValue := range valuesAera[itemKey] {
			summaryArea = summaryArea + float64(itemValue)
		}
		annualTableData2.Area = summaryArea

		summaryBudget := 0.0
		for _, itemValue := range values[itemKey] {
			summaryBudget = summaryBudget + float64(itemValue)
		}
		annualTableData2.Budget = summaryBudget

		summaryKm := 0.0
		for _, itemValue := range valuesKm[itemKey] {
			summaryKm = summaryKm + float64(itemValue)
		}
		annualTableData2.Range = math.Abs(summaryKm / 1000)

		annualTableData2s = append(annualTableData2s, annualTableData2)
	}

	dashBoard.Table.Table2 = annualTableData2s

	dashBoard.AnalystYear = modelResults[0].AnalystYear

	return dashBoard, nil
}

func (u *UseCase) annualMaintenanceAnalysisLimitBudget(idInteger int, modelResults []models.ModelResult, latestPrepareData models.PrepareData, maintenanceAnalysis models.MaintenanceAnalysis, maintenanceAnalysisPlan []models.MaintenanceAnalysisPlan, colorMap map[string]string) (interface{}, error) {

	maintenanceAnalysisStrategicBudgetType, err := u.Repo.GetMaintenanceAnalysisStrategicBudgetTypeById(maintenanceAnalysis.Condition)
	if err != nil {
		logs.Error(err)
		return models.DashboardAnnualMaintenance{}, err
	}

	maintenanceAnalysisStrategiTarget, err := u.Repo.GetMaintenanceAnalysisStrategiTargetById(*maintenanceAnalysis.Target)
	if err != nil {
		logs.Error(err)
		return models.DashboardAnnualMaintenance{}, err
	}

	maintenanceAnalysisRoad, err := u.Repo.GetMaintenanceAnalysisRoadNameById(maintenanceAnalysis.ID)
	if err != nil {
		logs.Error(err)
		return models.DashboardAnnualMaintenance{}, err
	}

	maintenanceAnalysisRoads := []string{}
	for _, item := range maintenanceAnalysisRoad {
		maintenanceAnalysisRoads = append(maintenanceAnalysisRoads, item.Name)
	}
	lane := ""
	switch maintenanceAnalysis.LaneTypeId {
	case 0:
		lane = "ทั้งหมด"
	case 1:
		lane = "1"
	case 2:
		lane = "2"
	case 3:
		lane = "3"
	}

	surface := ""
	switch maintenanceAnalysis.SurfaceTypeId {
	case 1:
		surface = "ลาดยาง"
	case 2:
		surface = "คอนกรีต"
	}

	comment := ""
	if maintenanceAnalysis.Comment == "" {
		comment = "-"
	} else {
		comment = maintenanceAnalysis.Comment
	}

	var dashBoard models.DashboardAnnualMaintenance
	dashBoard.Road = maintenanceAnalysisRoads
	dashBoard.Filter.Lane = lane
	dashBoard.Filter.SurfaceType = surface
	dashBoard.Comment = comment
	dashBoard.Filter.Km = maintenanceAnalysis.GroupKm
	dashBoard.Condition.Condition = maintenanceAnalysisStrategicBudgetType.Name
	dashBoard.Condition.Target = maintenanceAnalysisStrategiTarget.Name
	dashBoard.Condition.Discount = maintenanceAnalysis.Discount

	var filter = []string{}
	if maintenanceAnalysis.Iri1 != nil || maintenanceAnalysis.Iri2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Iri1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Iri1) + " < "
		}
		textFilter += "IRI"
		if maintenanceAnalysis.Iri2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Iri2)
		}

		filter = append(filter, textFilter)
	}

	if maintenanceAnalysis.Aadt1 != nil || maintenanceAnalysis.Aadt2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Aadt1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Aadt1) + " < "
		}
		textFilter += "AADT"
		if maintenanceAnalysis.Aadt2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Aadt2)
		}

		filter = append(filter, textFilter)
	}

	if maintenanceAnalysis.Ifi1 != nil || maintenanceAnalysis.Ifi2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Ifi1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Ifi1) + " < "
		}
		textFilter += "IFI"
		if maintenanceAnalysis.Ifi2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Ifi2)
		}

		filter = append(filter, textFilter)
	}

	dashBoard.Filter.Filter = filter

	values := map[string][]int{}
	valuesAera := map[string][]float64{}
	valuesKm := map[string][]int{}
	countIri := 0
	sums := map[string]int{}
	keys := []string{}
	keys = append(keys, "ซ่อมบำรุงปกติ")
	sum := 0
	sumArea := 0.0
	length := 0.0
	iriBeforeLength := 0.0
	iriAfterLength := 0.0
	for _, item := range modelResults {
		if item.Data.Repair {
			valuesAera[item.Data.IcResult.Name] = append(valuesAera[item.Data.IcResult.Name], item.Data.PrepareDataBefore.Area)
			values[item.Data.IcResult.Name] = append(values[item.Data.IcResult.Name], int(item.Data.UsedBudget))
			if !helpers.Contains(keys, item.Data.IcResult.Name) {
				keys = append(keys, item.Data.IcResult.Name)
			}

			valuesKm[item.Data.IcResult.Name] = append(valuesKm[item.Data.IcResult.Name], int(math.Abs(float64(item.KmStart)-float64(item.KmEnd))))
			sums[item.Data.IcResult.Name] = sums[item.Data.IcResult.Name] + int(item.Data.UsedBudget)
			sum = sum + int(item.Data.UsedBudget)

			iriAfterLength += float64(item.Data.PrepareDataBefore.Length) * (item.Data.RweResult.Iri)
		} else {
			valuesAera["ซ่อมบำรุงปกติ"] = append(valuesAera["ซ่อมบำรุงปกติ"], item.Data.PrepareDataBefore.Area)
			valuesKm["ซ่อมบำรุงปกติ"] = append(valuesKm["ซ่อมบำรุงปกติ"], int(math.Abs(float64(item.KmStart)-float64(item.KmEnd))))
			iriAfterLength += float64(item.Data.PrepareDataBefore.Length) * (item.Data.DeteriorationResult.Result.Iri)

		}

		iriBeforeLength += float64(item.Data.PrepareDataBefore.Length) * (item.Data.PrepareData.RoadCondition.Iri)

		length += float64(item.Data.PrepareDataBefore.Length)
		sumArea = sumArea + item.Data.PrepareDataBefore.Area
		countIri++
	}

	datsetsLable := []string{}
	datsetsValues := []float64{}
	datsetsArea := []float64{}
	datsetsColor := []string{}
	for _, itemKey := range keys {
		datsetsLable = append(datsetsLable, itemKey)
		summary := 0.0
		summaryArea := 0.0
		for _, itemValue := range valuesAera[itemKey] {
			summary = summary + ((float64(itemValue)/float64(sumArea)*100)/100)*100
			summaryArea = summaryArea + float64(itemValue)
		}
		if math.IsNaN(summary) {
			summary = 0.0
		}
		datsetsArea = append(datsetsArea, summaryArea)
		datsetsValues = append(datsetsValues, helpers.RoundFloat(summary, 2))

		color := colorMap[itemKey]
		if color == "" {
			color = "#000000"
		}

		datsetsColor = append(datsetsColor, color)
	}

	var bar models.AnnualBar
	bar.Name = "ปริมาณงาน"
	bar.Lable = datsetsLable
	bar.Data = datsetsValues
	bar.Area = datsetsArea
	bar.Color = datsetsColor
	dashBoard.Bar1 = bar

	var bar2 models.AnnualBar2
	datsets2Lable := []string{}
	datsets2Values := []float64{}
	datsets2Budget := []float64{}
	datsets2Color := []string{}
	for _, itemKey := range keys {
		datsets2Lable = append(datsets2Lable, itemKey)
		summary := 0.0
		summaryBudget := 0.0
		for _, itemValue := range values[itemKey] {
			summary = summary + ((float64(itemValue)/float64(sum)*100)/100)*100
			summaryBudget = summaryBudget + float64(itemValue)
		}
		if math.IsNaN(summary) {
			summary = 0.0
		}
		datsets2Budget = append(datsets2Budget, summaryBudget)
		datsets2Values = append(datsets2Values, helpers.RoundFloat(summary, 2))

		color := colorMap[itemKey]
		if color == "" {
			color = "#000000"
		}

		datsets2Color = append(datsets2Color, color)
	}

	bar2.Name = "ค่าซ่อมบำรุง"
	bar2.Lable = datsets2Lable
	bar2.Data = datsets2Values
	bar2.Budget = datsets2Budget
	bar2.Color = datsets2Color
	dashBoard.Bar2 = bar2

	var annualTableData1 models.AnnualTableData1

	iriAfter := iriAfterLength / length
	if math.IsNaN(iriAfter) {
		iriAfter = 0
	}

	iriBefore := iriBeforeLength / length
	if math.IsNaN(iriBefore) {
		iriBefore = 0
	}

	budget := 0.0
	if iriAfter > 0 || iriBefore > 0 {
		budget = float64(*maintenanceAnalysis.Budget)
	}
	annualTableData1.Budget = budget * 1000000

	annualTableData1.IriAfter = iriAfter
	annualTableData1.IriBefore = iriBefore

	dashBoard.Table.Table1 = annualTableData1

	var annualTableData2s []models.AnnualTableData2

	summaryArea := 0.0
	summaryBudget := 0.0
	summaryKm := 0.0
	for _, item := range modelResults {
		if !item.Data.Repair {
			summaryArea = summaryArea + float64(item.Data.PrepareDataBefore.Area)
			summaryBudget = summaryBudget + item.Data.UsedBudget
			summaryKm = summaryKm + float64(item.Data.PrepareDataBefore.Length)
		}
	}

	for _, itemKey := range keys {
		var annualTableData2 models.AnnualTableData2
		annualTableData2.Name = itemKey
		summaryArea := 0.0
		for _, itemValue := range valuesAera[itemKey] {
			summaryArea = summaryArea + float64(itemValue)
		}
		annualTableData2.Area = summaryArea

		summaryBudget := 0.0
		for _, itemValue := range values[itemKey] {
			summaryBudget = summaryBudget + float64(itemValue)
		}
		annualTableData2.Budget = summaryBudget

		summaryKm := 0.0
		for _, itemValue := range valuesKm[itemKey] {
			summaryKm = summaryKm + float64(itemValue)
		}
		annualTableData2.Range = math.Abs(summaryKm / 1000)

		annualTableData2s = append(annualTableData2s, annualTableData2)
	}

	dashBoard.Table.Table2 = annualTableData2s

	dashBoard.AnalystYear = modelResults[0].AnalystYear

	return dashBoard, nil
}

func (u *UseCase) annualMaintenanceAnalysisTargetIri(idInteger int, modelResults []models.ModelResult, latestPrepareData models.PrepareData, maintenanceAnalysis models.MaintenanceAnalysis, maintenanceAnalysisPlan []models.MaintenanceAnalysisPlan, colorMap map[string]string) (interface{}, error) {

	maintenanceAnalysisStrategicBudgetType, err := u.Repo.GetMaintenanceAnalysisStrategicBudgetTypeById(maintenanceAnalysis.Condition)
	if err != nil {
		logs.Error(err)
		return models.DashboardAnnualMaintenance{}, err
	}

	maintenanceAnalysisStrategiTarget, err := u.Repo.GetMaintenanceAnalysisStrategiTargetById(*maintenanceAnalysis.Target)
	if err != nil {
		logs.Error(err)
		return models.DashboardAnnualMaintenance{}, err
	}

	maintenanceAnalysisRoad, err := u.Repo.GetMaintenanceAnalysisRoadNameById(maintenanceAnalysis.ID)
	if err != nil {
		logs.Error(err)
		return models.DashboardAnnualMaintenance{}, err
	}

	maintenanceAnalysisRoads := []string{}
	for _, item := range maintenanceAnalysisRoad {
		maintenanceAnalysisRoads = append(maintenanceAnalysisRoads, item.Name)
	}
	lane := ""
	switch maintenanceAnalysis.LaneTypeId {
	case 0:
		lane = "ทั้งหมด"
	case 1:
		lane = "1"
	case 2:
		lane = "2"
	case 3:
		lane = "3"
	}

	surface := ""
	switch maintenanceAnalysis.SurfaceTypeId {
	case 1:
		surface = "ลาดยาง"
	case 2:
		surface = "คอนกรีต"
	}

	comment := ""
	if maintenanceAnalysis.Comment == "" {
		comment = "-"
	} else {
		comment = maintenanceAnalysis.Comment
	}

	var dashBoard models.DashboardAnnualMaintenance
	dashBoard.Road = maintenanceAnalysisRoads
	dashBoard.Filter.Lane = lane
	dashBoard.Filter.SurfaceType = surface
	dashBoard.Comment = comment
	dashBoard.Filter.Km = maintenanceAnalysis.GroupKm
	dashBoard.Condition.Condition = maintenanceAnalysisStrategicBudgetType.Name
	dashBoard.Condition.Target = maintenanceAnalysisStrategiTarget.Name
	dashBoard.Condition.Discount = maintenanceAnalysis.Discount

	var filter = []string{}
	if maintenanceAnalysis.Iri1 != nil || maintenanceAnalysis.Iri2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Iri1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Iri1) + " < "
		}
		textFilter += "IRI"
		if maintenanceAnalysis.Iri2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Iri2)
		}

		filter = append(filter, textFilter)
	}

	if maintenanceAnalysis.Aadt1 != nil || maintenanceAnalysis.Aadt2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Aadt1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Aadt1) + " < "
		}
		textFilter += "AADT"
		if maintenanceAnalysis.Aadt2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Aadt2)
		}

		filter = append(filter, textFilter)
	}

	if maintenanceAnalysis.Ifi1 != nil || maintenanceAnalysis.Ifi2 != nil {
		textFilter := ""
		if maintenanceAnalysis.Ifi1 != nil {
			textFilter += fmt.Sprintf("%f", *maintenanceAnalysis.Ifi1) + " < "
		}
		textFilter += "IFI"
		if maintenanceAnalysis.Ifi2 != nil {
			textFilter += " < " + fmt.Sprintf("%f", *maintenanceAnalysis.Ifi2)
		}

		filter = append(filter, textFilter)
	}

	dashBoard.Filter.Filter = filter

	values := map[string][]int{}
	valuesAera := map[string][]float64{}
	valuesKm := map[string][]int{}
	countIri := 0
	iriBefore := 0.0
	iriAfter := 0.0
	length := 0.0
	sums := map[string]int{}
	keys := []string{}
	keys = append(keys, "ซ่อมบำรุงปกติ")
	sum := 0
	sumArea := 0.0

	for _, item := range modelResults {
		if item.Data.Repair {
			valuesAera[item.Data.IcResult.Name] = append(valuesAera[item.Data.IcResult.Name], item.Data.PrepareDataBefore.Area)
			values[item.Data.IcResult.Name] = append(values[item.Data.IcResult.Name], int(item.Data.UsedBudget))
			if !helpers.Contains(keys, item.Data.IcResult.Name) {
				keys = append(keys, item.Data.IcResult.Name)
			}

			valuesKm[item.Data.IcResult.Name] = append(valuesKm[item.Data.IcResult.Name], int(math.Abs(float64(item.KmStart)-float64(item.KmEnd))))
			sums[item.Data.IcResult.Name] = sums[item.Data.IcResult.Name] + int(item.Data.UsedBudget)
			sum = sum + int(item.Data.UsedBudget)
			iriAfter += item.Data.RweResult.Iri * float64(item.Data.PrepareDataBefore.Length)
		} else {
			valuesAera["ซ่อมบำรุงปกติ"] = append(valuesAera["ซ่อมบำรุงปกติ"], item.Data.PrepareDataBefore.Area)
			valuesKm["ซ่อมบำรุงปกติ"] = append(valuesKm["ซ่อมบำรุงปกติ"], int(math.Abs(float64(item.KmStart)-float64(item.KmEnd))))
			iriAfter += item.Data.DeteriorationResult.Result.Iri * float64(item.Data.PrepareDataBefore.Length)
		}

		length += float64(item.Data.PrepareDataBefore.Length)

		sumArea = sumArea + item.Data.PrepareDataBefore.Area
		iriBefore += item.Data.PrepareData.RoadCondition.Iri * float64(item.Data.PrepareDataBefore.Length)
		countIri++
	}

	datsetsLable := []string{}
	datsetsValues := []float64{}
	datsetsArea := []float64{}
	datsetsColor := []string{}
	for _, itemKey := range keys {
		datsetsLable = append(datsetsLable, itemKey)
		summary := 0.0
		summaryArea := 0.0
		for _, itemValue := range valuesAera[itemKey] {
			summary = summary + ((float64(itemValue)/float64(sumArea)*100)/100)*100
			summaryArea = summaryArea + float64(itemValue)
		}
		if math.IsNaN(summary) {
			summary = 0.0
		}
		datsetsArea = append(datsetsArea, summaryArea)
		datsetsValues = append(datsetsValues, helpers.RoundFloat(summary, 2))

		color := colorMap[itemKey]
		if color == "" {
			color = "#000000"
		}

		datsetsColor = append(datsetsColor, color)
	}

	var bar models.AnnualBar
	bar.Name = "ปริมาณงาน"
	bar.Lable = datsetsLable
	bar.Data = datsetsValues
	bar.Area = datsetsArea
	bar.Color = datsetsColor
	dashBoard.Bar1 = bar

	var bar2 models.AnnualBar2
	datsets2Lable := []string{}
	datsets2Values := []float64{}
	datsets2Budget := []float64{}
	datsets2Color := []string{}
	for _, itemKey := range keys {
		datsets2Lable = append(datsets2Lable, itemKey)
		summary := 0.0
		summaryBudget := 0.0
		for _, itemValue := range values[itemKey] {
			summary = summary + ((float64(itemValue)/float64(sum)*100)/100)*100
			summaryBudget = summaryBudget + float64(itemValue)
		}
		if math.IsNaN(summary) {
			summary = 0.0
		}
		datsets2Budget = append(datsets2Budget, summaryBudget)
		datsets2Values = append(datsets2Values, helpers.RoundFloat(summary, 2))

		color := colorMap[itemKey]
		if color == "" {
			color = "#000000"
		}

		datsets2Color = append(datsets2Color, color)
	}

	bar2.Name = "ค่าซ่อมบำรุง"
	bar2.Lable = datsets2Lable
	bar2.Data = datsets2Values
	bar2.Budget = datsets2Budget
	bar2.Color = datsets2Color
	dashBoard.Bar2 = bar2

	var annualTableData1 models.AnnualTableData1

	budget := 0.0
	if maintenanceAnalysis.Iri != nil {
		budget = float64(*maintenanceAnalysis.Iri)
	}

	iriAfter = iriAfter / length
	if math.IsNaN(iriAfter) {
		iriAfter = 0
	}

	iriBefore = iriBefore / length
	if math.IsNaN(iriBefore) {
		iriBefore = 0
	}

	annualTableData1.Budget = budget
	annualTableData1.IriAfter = iriAfter
	annualTableData1.IriBefore = iriBefore

	dashBoard.Table.Table1 = annualTableData1

	var annualTableData2s []models.AnnualTableData2

	summaryArea := 0.0
	summaryBudget := 0.0
	summaryKm := 0.0
	for _, item := range modelResults {
		if !item.Data.Repair {
			summaryArea = summaryArea + float64(item.Data.PrepareDataBefore.Area)
			summaryBudget = summaryBudget + item.Data.UsedBudget
			summaryKm = summaryKm + float64(item.Data.PrepareDataBefore.Length)
		}
	}

	for _, itemKey := range keys {
		var annualTableData2 models.AnnualTableData2
		annualTableData2.Name = itemKey
		summaryArea := 0.0
		for _, itemValue := range valuesAera[itemKey] {
			summaryArea = summaryArea + float64(itemValue)
		}
		annualTableData2.Area = summaryArea

		summaryBudget := 0.0
		for _, itemValue := range values[itemKey] {
			summaryBudget = summaryBudget + float64(itemValue)
		}
		annualTableData2.Budget = summaryBudget

		summaryKm := 0.0
		for _, itemValue := range valuesKm[itemKey] {
			summaryKm = summaryKm + float64(itemValue)
		}
		annualTableData2.Range = math.Abs(summaryKm / 1000)

		annualTableData2s = append(annualTableData2s, annualTableData2)
	}

	dashBoard.Table.Table2 = annualTableData2s

	dashBoard.AnalystYear = modelResults[0].AnalystYear

	return dashBoard, nil
}

func (u *UseCase) ExportData(ID int) (interface{}, error) {
	folderPath := os.Getenv("ANALYSES_EXPORT")
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		return "", err
	}

	// Create a new CSV file for writing
	fileName := folderPath + uuid.New().String() + ".csv"
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Define the data to be written to the CSV file
	data := [][]string{
		{"Name", "Age", "City"},
		{"Alice", "25", "New York"},
		{"Bob", "30", "Los Angeles"},
		{"Charlie", "22", "Chicago"},
	}

	// Write the data to the CSV file
	for _, record := range data {
		err := writer.Write(record)
		if err != nil {
			panic(err)
		}
	}
	return os.Getenv("STORAGE_IP") + "/" + fileName, nil
}

func (u *UseCase) GetMaintenanceAnalysisPreviousID(ID int) (int, error) {
	newID, err := u.Repo.GetMaintenanceAnalysisPreviousID(ID)
	if err != nil {
		if err.Error() != "record not found" {
			logs.Error(err)
			return 0, responses.NewAppErr(400, err.Error())
		}
		return ID, nil
	}
	return newID, nil
}

func (u *UseCase) GetAnalysisModel(ID int, req requests.MaintenanceAnalysisModel) (interface{}, error) {
	data, err := u.Repo.GetAnalysisModel(ID, req)
	if err != nil {
		if err.Error() != "record not found" {
			logs.Error(err)
			return 0, responses.NewAppErr(400, err.Error())
		}
		return data, nil
	}

	roadInfoData, err := u.Repo.GetRoadInfo()
	if err != nil {
		return data, nil
	}

	roadGroupData, err := u.Repo.GetRoadGroup()
	if err != nil {
		return data, nil
	}

	roadData, err := u.Repo.GetRoads()
	if err != nil {
		return data, nil
	}

	var res []responses.AnalysesModel
	// refCriteriaMethod, _ := u.GetRefCriteriaMethod()
	roadIDData := make(map[int]bool)
	roadIDs := strings.Split(req.RoadIDs, ",")
	for _, roadID := range roadIDs {
		num, _ := strconv.Atoi(roadID)
		roadIDData[num] = true
	}

	for _, item := range data {
		isData := false
		iri := item.Data.IriAfter
		aadt := item.Data.PrepareDataBefore.Aadt
		ifi := item.Data.IfiAfter
		laneNo := item.Data.PrepareDataBefore.RoadGeom.LaneNo
		interventionCriteriaChangeID := 0
		if item.Data.Repair {
			interventionCriteriaChangeID = item.Data.InterventionCriteriaChangeID
		}

		if len(req.RoadIDs) > 0 {
			_, isVal := roadIDData[item.RoadID]
			if !isVal {
				continue
			}
		} else {
			isData = true
		}

		if req.LaneNo != nil {
			if laneNo == *req.LaneNo {
				isData = true
			}
		} else {
			isData = true
		}

		if req.InterventionCriteriaID != nil {
			if interventionCriteriaChangeID == *req.InterventionCriteriaID {
				isData = true
			}
		} else {
			isData = true
		}

		if req.Iri1 != nil && req.Iri2 != nil {
			if *req.Iri1 <= iri && *req.Iri2 >= iri {
				isData = true
			} else {
				isData = false
			}
		} else if req.Iri1 != nil && req.Iri2 == nil {
			if *req.Iri1 <= iri {
				isData = true
			} else {
				isData = false
			}
		} else if req.Iri1 == nil && req.Iri2 != nil {
			if *req.Iri2 >= iri {
				isData = true
			} else {
				isData = false
			}
		}

		if req.Aadt1 != nil && req.Aadt2 != nil {
			if *req.Aadt1 <= aadt && *req.Aadt2 >= aadt {
				isData = true
			} else {
				isData = false
			}
		} else if req.Aadt1 != nil && req.Aadt2 == nil {
			if *req.Aadt1 <= aadt {
				isData = true
			} else {
				isData = false
			}
		} else if req.Aadt1 == nil && req.Aadt2 != nil {
			if *req.Aadt2 >= aadt {
				isData = true
			} else {
				isData = false
			}
		}

		if req.Ifi1 != nil && req.Ifi2 != nil {
			if *req.Ifi1 <= ifi && *req.Ifi2 >= ifi {
				isData = true
			} else {
				isData = false
			}
		} else if req.Ifi1 != nil && req.Ifi2 == nil {
			if *req.Ifi1 <= ifi {
				isData = true
			} else {
				isData = false
			}
		} else if req.Ifi1 == nil && req.Ifi2 != nil {
			if *req.Ifi2 >= ifi {
				isData = true
			} else {
				isData = false
			}
		}

		if !isData {
			continue
		}

		var r responses.AnalysesModel
		roadinfo := roadInfoData[item.RoadID]
		road := roadData[item.RoadID]
		roadGrp := roadGroupData[road.RoadGroupId]
		if item.MaintenanceAnalysisTypeID == 2 {
			r.Plan = "ซ่อมบำรุงปกติ"
		} else {
			if item.RepairBudgetType == "No Budget" {
				r.Plan = "ซ่อมบำรุงปกติ"
			} else if item.RepairBudgetType == "Unlimited Budget" {
				r.Plan = "ไม่จำกัดงบประมาณ"
			} else {
				r.Plan = fmt.Sprintf("แผน %d", item.Plan)
			}
		}

		r.ID = item.ID
		r.Year = item.AnalystYear
		r.RoadGroupName = roadGrp.Name
		r.RoadName = roadinfo.Name
		r.KmStart = float64(item.KmStart)
		r.KmEnd = float64(item.KmEnd)
		r.Distance = math.Abs(r.KmStart-r.KmEnd) / 1000
		r.Lane = item.Data.PrepareDataBefore.RoadGeom.LaneNo
		r.InterventionCriteria = interventionCriteriaChangeID
		r.Cost = item.Data.UsedBudget
		r.Area = item.Data.PrepareDataBefore.Area
		r.VolumeAadt = int(item.Data.PrepareDataBefore.Aadt)
		r.BC = item.Data.BcAfter
		r.IriBefore = item.Data.PrepareData.RoadCondition.Iri
		r.IriAfter = item.Data.IriAfter
		res = append(res, r)
	}

	return res, nil
}

func (u *UseCase) GetRefCriteriaMethod() (interface{}, error) {
	data, err := u.Repo.GetRefCriteriaMethod()
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return data, responses.NewAppErr(400, err.Error())
	}

	// return data, nil
	var interventionCriterias []responses.MaintenanceInterventionCriteriaModel
	var childrenNotRepaires []responses.MaintenanceInterventionCriteriaChildrenModel
	childrenNotRepaire := responses.MaintenanceInterventionCriteriaChildrenModel{ID: 0, Label: "ซ่อมบำรุงปกติ", Selected: false}
	childrenNotRepaires = append(childrenNotRepaires, childrenNotRepaire)
	interventionCriterias = append(interventionCriterias, responses.MaintenanceInterventionCriteriaModel{ID: 0, Label: "ซ่อมบำรุงปกติ", Children: childrenNotRepaires})
	for _, item := range data {
		itvcData, err := u.Repo.GetInterventionCriteriaByID(item.ID)
		if len(itvcData) == 0 {
			continue
		}
		if err != nil {
			logs.Error(err)
			if err == gorm.ErrRecordNotFound {
				return responses.NoData{}, nil
			}
			return data, responses.NewAppErr(400, err.Error())
		}
		var interventionCriteria responses.MaintenanceInterventionCriteriaModel
		interventionCriteria.ID = item.ID
		interventionCriteria.Label = item.Name
		var childrens []responses.MaintenanceInterventionCriteriaChildrenModel
		for _, item := range itvcData {
			var children responses.MaintenanceInterventionCriteriaChildrenModel
			// var dd interventionCriteria.Children
			children.ID = item.Id
			children.Label = item.MaintenanceStandardName
			children.Selected = false
			childrens = append(childrens, children)
		}
		interventionCriteria.Children = childrens
		interventionCriterias = append(interventionCriterias, interventionCriteria)

	}

	return interventionCriterias, nil
}

func (u *UseCase) UpdateAnalysisModel(ID int, reqs []requests.MaintenanceAnalysisModelReq) (interface{}, error) {
	for _, item := range reqs {
		dataID := item.ID
		interventionCriteriaID := item.InterventionCriteriaID
		u.Repo.UpdateAnalysisModel(ID, dataID, interventionCriteriaID)
	}
	return "", nil
}

func (u *UseCase) DashboardMapFilter(ID int) (interface{}, error) {
	data, err := u.Repo.GetMaintenanceAnalysisById(ID)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return data, responses.NewAppErr(400, err.Error())
	}

	refOwners, err := u.Repo.GetRefOwner()
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return data, responses.NewAppErr(400, err.Error())
	}

	refCriteriaMethods, err := u.Repo.GetCriteriaMethod()
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return data, responses.NewAppErr(400, err.Error())
	}

	var criterias []responses.Criterias
	for _, item := range refOwners {
		gradeData, _ := u.Repo.GetRefOwnerCondition(item.ID)
		grade := gradeData["AC"]
		gradeCC := gradeData["CC"]
		criterias = append(criterias, responses.Criterias{ID: item.ID, Name: item.Name, Grade: grade, GradeCC: gradeCC})
	}

	var methods []responses.Methods
	for _, item := range refCriteriaMethods {
		methods = append(methods, responses.Methods{ID: item.ID, Name: item.Name, Color: item.Color})
	}
	var filter responses.Filter
	// return data, nil
	if data.MaintenanceAnalysisTypeId == 1 {
		switch data.Condition {
		case 1:
			years := []int{}
			for i := 0; i < *data.Year; i++ {
				years = append(years, i+1)
			}

			filter.Year = years
			var plans []responses.Plan
			plans = append(plans, responses.Plan{ID: 1, Name: "ซ่อมบำรุงปกติ"})
			plans = append(plans, responses.Plan{ID: 0, Name: "ไม่จำกัดงบประมาณ"})

			filter.Year = years
			filter.Plan = plans
			filter.Display = append(filter.Display, responses.Displays{ID: 1, Name: "IRI"})
			filter.Display = append(filter.Display, responses.Displays{ID: 2, Name: "วิธีการซ่อม"})
			filter.Criteria = criterias
			filter.Method = methods
			return filter, nil
		case 2:
			years := []int{}
			for i := 0; i < *data.Year; i++ {
				years = append(years, i+1)
			}
			filter.Year = years
			var plans []responses.Plan
			plans = append(plans, responses.Plan{ID: *data.NumberPlan + 2, Name: "ซ่อมบำรุงปกติ"})
			plans = append(plans, responses.Plan{ID: 0, Name: "ไม่จำกัดงบประมาณ"})

			for i := 0; i < *data.NumberPlan; i++ {
				plans = append(plans, responses.Plan{ID: i + 1, Name: fmt.Sprintf("แผนที่ %d", i+1)})
			}
			filter.Year = years
			filter.Plan = plans
			filter.Display = append(filter.Display, responses.Displays{ID: 1, Name: "IRI"})
			filter.Display = append(filter.Display, responses.Displays{ID: 2, Name: "วิธีการซ่อม"})
			filter.Criteria = criterias
			filter.Method = methods
			return filter, nil
		case 3:
			var filter responses.Filter
			years := []int{}
			for i := 0; i < *data.Year; i++ {
				years = append(years, i+1)
			}
			filter.Year = years
			var plans []responses.Plan
			for i := 0; i < *data.NumberPlan; i++ {
				plans = append(plans, responses.Plan{ID: i + 1, Name: fmt.Sprintf("แผนที่ %d", i+1)})
			}
			filter.Year = years
			filter.Plan = plans
			filter.Display = append(filter.Display, responses.Displays{ID: 1, Name: "IRI"})
			filter.Display = append(filter.Display, responses.Displays{ID: 2, Name: "วิธีการซ่อม"})
			filter.Criteria = criterias
			filter.Method = methods
			return filter, nil
		}

	} else {

		filter.Year = []int{}
		filter.Plan = []responses.Plan{}
		filter.Display = append(filter.Display, responses.Displays{ID: 1, Name: "IRI"})
		filter.Display = append(filter.Display, responses.Displays{ID: 2, Name: "วิธีการซ่อม"})
		filter.Criteria = criterias
		filter.Method = methods
		return filter, nil
	}
	return filter, nil
}

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type DashboardMap struct {
	Title     string   `json:"title"`
	RoadName  string   `json:"road_name"`
	TheGeom   Geometry `json:"the_geom"`
	IriBefore float64  `json:"iri_before"`
	IriAfter  float64  `json:"iri_after"`
	Year      int      `json:"year"`
	Color     string   `json:"color"`
}

func (u *UseCase) DashboardMap(ID int, req requests.MapFilter) (interface{}, error) {
	analysis, err := u.Repo.GetMaintenanceAnalysisById(ID)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return "", responses.NewAppErr(400, err.Error())
	}

	year := *req.Year
	check := helpers.IntInSlice(analysis.Condition, []int{2, 3, 6, 7})
	if check {
		year = year - 1
	}

	var refCriteriaMethods []models.RefCriteriaMethod
	colorCriteria := make(map[string]string)
	helpers.PrintlnJson(ID, analysis.MaintenanceAnalysisTypeId, req.Plan)
	dataAllYear, err := u.Repo.DashboardMapAllYear(ID, analysis.MaintenanceAnalysisTypeId, req.Plan)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return dataAllYear, responses.NewAppErr(400, err.Error())
	}
	colorIndex := 0
	for _, item := range dataAllYear {
		if req.Display == 2 {
			if item.Data.Repair {
				//set color\
				if colorCriteria[item.Data.IcResult.Name] == "" {
					colorCriteria[item.Data.IcResult.Name] = colors[colorIndex]
					refCriteriaMethods = append(refCriteriaMethods, models.RefCriteriaMethod{Name: item.Data.IcResult.Name, Color: colors[colorIndex]})
					colorIndex++
				}
			}
		}
	}
	// return refCriteriaMethods, nil
	refCriteriaMethods = append(refCriteriaMethods, models.RefCriteriaMethod{Name: "ซ่อมบำรุงปกติ", Color: "#000000"})

	data, err := u.Repo.DashboardMap(ID, analysis.MaintenanceAnalysisTypeId, year, req.Plan)
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return responses.NoData{}, nil
		}
		return data, responses.NewAppErr(400, err.Error())
	}

	var dashboardMapDatas []responses.DashboardMapData
	grades, _ := u.Repo.GetRefOwnerCondition(*req.Criteria)
	for _, item := range data {
		var dashboardMapData responses.DashboardMapData
		dashboardMapData.KmStart = item.KmStart
		dashboardMapData.KmEnd = item.KmEnd
		dashboardMapData.Display = req.Display
		dashboardMapData.TheGeom = Geometry(item.Data.Geometry)
		dashboardMapData.IriBefore = item.Data.DeteriorationResult.Result.Iri
		if item.Data.Repair {
			dashboardMapData.IriAfter = item.Data.RweResult.Iri
		} else {
			dashboardMapData.IriAfter = item.Data.DeteriorationResult.Result.Iri
		}

		dashboardMapData.Year = req.Year
		dashboardMapData.RoadName = item.Data.PrepareData.Road.RoadName
		if req.Display == 1 {
			iri := dashboardMapData.IriAfter
			if item.Data.Repair {
				dashboardMapData.Title = fmt.Sprintf("IRI %.2f ม./กม.", iri)
				helpers.PrintlnJson("ddddddd", item.Data.IcResult.RefSurface.SurfaceGroup)
				if item.Data.IcResult.RefSurface.SurfaceGroup == "Asphalt" {
					for _, grade := range grades["AC"] {

						if grade.LeftValue <= iri && iri < grade.RightValue {
							dashboardMapData.Color = grade.Color
							dashboardMapDatas = append(dashboardMapDatas, dashboardMapData)
							goto EXITLOOP
						}
					}
				} else if item.Data.IcResult.RefSurface.SurfaceGroup == "Concrete" {
					for _, grade := range grades["CC"] {
						if grade.LeftValue <= iri && iri < grade.RightValue {
							dashboardMapData.Color = grade.Color
							dashboardMapDatas = append(dashboardMapDatas, dashboardMapData)
							goto EXITLOOP
						}
					}
				}
			} else {
				if item.Data.DeteriorationResult.IsSurfaceAc {
					for _, grade := range grades["AC"] {

						if grade.LeftValue <= iri && iri < grade.RightValue {
							dashboardMapData.Color = grade.Color
							dashboardMapDatas = append(dashboardMapDatas, dashboardMapData)
							goto EXITLOOP
						}
					}
				} else {
					for _, grade := range grades["CC"] {
						if grade.LeftValue <= iri && iri < grade.RightValue {
							dashboardMapData.Color = grade.Color
							dashboardMapDatas = append(dashboardMapDatas, dashboardMapData)
							goto EXITLOOP
						}
					}
				}
			}

		} else {
			if item.Data.Repair {
				helpers.PrintlnJson("vvvvvv", item.Data.IcResult.Name)
				//======
				dashboardMapData.Title = item.Data.IcResult.Name
				dashboardMapData.Color = colorCriteria[item.Data.IcResult.Name]

				// colors[colorIndex] = colors[colorIndex]

			} else {
				dashboardMapData.Title = "ซ่อมบำรุงปกติ"
				dashboardMapData.Color = "#000000"
			}
		}

	EXITLOOP:
		dashboardMapDatas = append(dashboardMapDatas, dashboardMapData)
	}
	// return colors, nil
	// refCriteriaMethod, err := u.Repo.GetRefCriteriaMethod()
	// if err != nil {
	// 	logs.Error(err)
	// 	if err == gorm.ErrRecordNotFound {
	// 		return responses.NoData{}, nil
	// 	}
	// 	return data, responses.NewAppErr(400, err.Error())
	// }
	var dashboardMapCriteriaMethod []responses.DashboardMapCriteriaMethod
	checkDup := make(map[string]bool)
	for _, item := range refCriteriaMethods {
		if !checkDup[item.Name] {
			dashboardMapCriteriaMethod = append(dashboardMapCriteriaMethod, responses.DashboardMapCriteriaMethod{Name: item.Name, Color: item.Color})
			checkDup[item.Name] = true
		}

	}
	var dashboardMap responses.DashboardMap
	dashboardMap.CriteriaMethod = dashboardMapCriteriaMethod
	dashboardMap.Items = dashboardMapDatas
	return dashboardMap, nil
}

func (u *UseCase) CheckPrepareDataById(ID int) (interface{}, error) {

	maintenanceAnalysis, err := u.Repo.GetMaintenanceAnalysisById(ID)
	if err != nil {
		log.Println(err)
		if err.Error() == "record not found" {
			return models.MaintenanceAnalysis{}, responses.NewNotFoundError()
		}
		return models.MaintenanceAnalysis{}, responses.NewAppErr(400, err.Error())
	}

	var checkPrepareDataStatus responses.CheckPrepareDataStatus
	checkPrepareDataStatus.Status = maintenanceAnalysis.PrepareDataStatus

	return checkPrepareDataStatus, nil
}

func (u *UseCase) GetPrepareDataIdById(ID int) (interface{}, error) {

	prepareDataIds, err := u.Repo.GetPrepareDataIdByMaintenanceAnalysisId(ID)
	if err != nil {
		log.Println(err)
		if err.Error() == "record not found" {
			return models.MaintenanceAnalysis{}, responses.NewNotFoundError()
		}
		return models.MaintenanceAnalysis{}, responses.NewAppErr(400, err.Error())
	}

	return prepareDataIds, nil
}

func (u *UseCase) GetPrepareDataById(ID int) (interface{}, error) {
	prepareData, err := u.Repo.GetPrepareDataById(ID)
	if err != nil {
		log.Println(err)
		if err.Error() == "record not found" {
			return models.MaintenanceAnalysis{}, responses.NewNotFoundError()
		}
		return models.MaintenanceAnalysis{}, responses.NewAppErr(400, err.Error())
	}

	return prepareData, nil
}

func (u *UseCase) GetPrepareDataAllByAnalysisSelected(ID int) (interface{}, error) {
	data, err := u.Repo.GetPrepareDataAllByAnalysisSelected(ID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	ids := []int{}
	for _, item := range data {
		ids = append(ids, item.ID)
	}
	return ids, nil
}
