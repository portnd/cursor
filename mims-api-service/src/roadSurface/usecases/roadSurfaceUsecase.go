package usecases

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	_ "gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadSurface/domains"
	"gorm.io/gorm"
)

type roadSurfaceUsecase struct {
	roadSurfaceRepo domains.RoadSurfaceRepositories
}

func NewRoadSurfaceUsecase(repo domains.RoadSurfaceRepositories) domains.RoadSurfaceUsecases {
	return &roadSurfaceUsecase{roadSurfaceRepo: repo}
}

func (r *roadSurfaceUsecase) GetRoadSurfaceList(roadId string, access []string, uID int) ([]responses.RoadSurfaceResponds, error) {
	var hasPermission bool
	var respondResult []responses.RoadSurfaceResponds
	result, err := r.roadSurfaceRepo.GetRoadSurfaceData(roadId, hasPermission)
	if len(result) == 0 {
		return nil, nil
	}
	if err != nil {
		logs.Error(err)
		return respondResult, err
	}

	newResult, err := helpers.RemoveEmptySlice(result, "RoadSurfaceLane")
	if err != nil {
		logs.Error(err)
		return respondResult, err
	}
	responds := newResult.([]models.RoadSurfacePreload)

	var refDirection int
	for _, v := range responds {
		var newRespond responses.RoadSurfaceResponds
		newRespond, err = r.getRoadSurface(v)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		newRespond.Permissions = r.permission(v.Status, access)
		// newRespond.StatusCode = v.Status
		newRespond.UpdatedBy = r.updateBy(uID)
		respondResult = append(respondResult, newRespond)
		refDirection = v.RoadInfo.RefDirectionId
	}
	direction := refDirection
	// road_id_int, _ := strconv.Atoi(roadId)
	// direction, err := r.roadSurfaceRepo.FindRefDirectionIDByRoadID(road_id_int)
	// if err != nil {
	// 	logs.Error(err)
	// 	return respondResult, err
	// }

	if direction == 1 {
		sortKmMinToMax(respondResult)
	} else if direction == 2 {
		sortKmMaxToMin(respondResult)
	}
	if err != nil {
		logs.Error(err)
		return respondResult, err
	}
	helpers.SortStructByID(respondResult, false)
	return respondResult, nil
}

// func (r *roadSurfaceUsecase) SumKm(respondResult []responses.RoadSurfaceResponds) float64 {
// 	var totalKm float64
// 	for _, v := range respondResult {
// 		totalKm += math.Abs(v.KmEnd - v.KmStart)
// 	}
// 	roundedTotalKm := math.Round(totalKm*1000) / 1000
// 	return roundedTotalKm
// }

func sortKmMinToMax(respondResult []responses.RoadSurfaceResponds) {
	i := 0
	for i+1 < len(respondResult) {
		head := respondResult[i]
		tail := respondResult[i+1]
		if respondResult[i].Items.KmStart > respondResult[i+1].Items.KmStart {
			respondResult[i] = tail
			respondResult[i+1] = head
			i = 0
			continue
		}
		i++
	}
}

func sortKmMaxToMin(respondResult []responses.RoadSurfaceResponds) {
	i := 0
	for i+1 < len(respondResult) {
		head := respondResult[i]
		tail := respondResult[i+1]
		if respondResult[i].Items.KmStart < respondResult[i+1].Items.KmStart {
			respondResult[i] = tail
			respondResult[i+1] = head
			i = 0
			continue
		}
		i++
	}
}

// func (r *roadSurfaceUsecase) getLanes(data models.RoadSurfacePreload) ([]responses.Lane, error) {
// 	var lanes []responses.Lane
// 	var lane responses.Lane
// 	for _, k := range data.RoadSurfaceLane {
// 		lane.Surface.ID = k.RefSurfaceId
// 		name, surfaceGroup, err := r.roadSurfaceRepo.FindRefSurfaceNameAndSurfaceGroup(k.RefSurfaceId)
// 		if err != nil {
// 			return lanes, err
// 		}
// 		lane.Surface.Name = name
// 		lane.Surface.SurfaceGroup = surfaceGroup
// 		if surfaceGroup == "Asphalt" {
// 			lane.Surface.ColorCode = "#7239ea"
// 		}
// 		if surfaceGroup == "Concrete" {
// 			lane.Surface.ColorCode = "#0e4285"
// 		}

// 		nameDri, err := r.roadSurfaceRepo.FindRefDirectionName(data.RoadId)
// 		if err != nil {
// 			logs.Error(err)
// 			return lanes, err
// 		}
// 		lane.Direction = nameDri
// 		lane.LaneNo = k.LaneNo
// 		geom := ""
// 		if k.TheGeom != "" {
// 			geom = (k.TheGeom)
// 		} else if data.TheGeom != "" {
// 			geom = (data.TheGeom)
// 		}
// 		if err != nil {
// 			logs.Error(err)
// 			return lanes, err
// 		}
// 		lane.GeomCl = geom
// 		lanes = append(lanes, lane)
// 	}

// 	return lanes, nil
// }

func (r *roadSurfaceUsecase) getRoadSurface(data models.RoadSurfacePreload) (responses.RoadSurfaceResponds, error) {
	// var responds []responses.RoadSurfaceResponds
	var respond responses.RoadSurfaceResponds
	var items responses.Items
	copier.Copy(&items, &data)
	respond.Id = data.Id
	respond.UpdatedDate = helpers.SetTimeToString(data.UpdatedDate)

	// status, err := r.roadSurfaceRepo.FindDataStatus(data.Status)
	// if err != nil {
	// 	logs.Error(err)
	// 	return respond, err
	// }
	// respond.Status = status
	// if data.Status == "R" {
	// 	respond.RejectReason = data.RejectReason
	// } else {
	// 	respond.RejectReason = ""
	// }

	items.SurfaceShoulderLeft = data.RefSurfaceIdShoulderLeftData
	items.SurfaceShoulderRight = data.RefSurfaceIdShoulderRightData
	if data.RefMaterialBaseId != 0 {
		items.MaterialBase = &data.MaterialBase
	}
	items.ThicknessSubbase = data.ThicknessSubbase
	if data.RefMaterialSubbaseId != 0 {
		items.MaterialSubbase = &data.MaterialSubbase
	}

	items.ThicknessSubgrade = data.ThicknessSubgrade
	if data.RefMaterialSubgradeId != 0 {
		items.MaterialSubgrade = &data.MaterialSubgrade
	}
	items.SurfaceCrossSectionCode = data.SurfaceCrosssectionCode
	respond.Items = items
	direction := strconv.Itoa(data.RoadInfo.RefDirectionId)
	respond.LaneCount = len(data.RoadSurfaceLane)
	var lanes []responses.Lane
	for _, l := range data.RoadSurfaceLane {
		var lane responses.Lane
		if l.RefSurfaceId != 0 {
			lane.Surface = l.RefSurface
			lane.LaneNo = l.LaneNo
			lane.GeomCl = l.TheGeom
			lane.Direction = direction
			lanes = append(lanes, lane)
		}

	}
	respond.Lanes = lanes
	return respond, nil
}

// func (r *roadSurfaceUsecase) GetMenu(userId uint) ([]models.AccessControl, error) {
// 	roles, err := r.roadSurfaceRepo.GetRole(userId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var roleIds []int
// 	for _, item := range roles {
// 		roleIds = append(roleIds, item.RoleID)
// 	}

// 	data, err := r.roadSurfaceRepo.GetAccessControl(roleIds)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return data, nil
// }

func (r *roadSurfaceUsecase) permission(status string, access []string) responses.Permissions {
	var result responses.Permissions
	if helpers.HasPermission([]string{"road_summary_manage_surface"}, access) && helpers.Contains([]string{"A", "T", "R"}, status) {
		result.CanEdit = true
		result.CanDelete = true
	}
	if helpers.HasPermission([]string{"dashboard_road_surface_access"}, access) && helpers.Contains([]string{"W"}, status) {
		result.CanApprove = true
		result.CanReject = true
	}
	if helpers.HasPermission([]string{"dashboard_road_surface_access"}, access) && helpers.Contains([]string{"T"}, status) {
		result.CanSend = true
	}
	return result
}

func (r *roadSurfaceUsecase) updateBy(userID int) responses.UpdatedBy {
	userInfo := r.roadSurfaceRepo.GetUserInfoByUpdatedBy(userID)
	updatedBy := responses.UpdatedBy{
		Id:       userInfo.Id,
		Email:    userInfo.Email,
		FullName: userInfo.Firstname + " " + userInfo.Lastname,
		// Department: models.Department{
		// 	Id:   userInfo.DepartmentId,
		// 	Name: userInfo.Department.Name,
		// },
		ProfilePicture: userInfo.ProfileImgPath,
	}
	return updatedBy
}

func (r *roadSurfaceUsecase) PostRoadSurface(request requests.RoadSurface, uid int) ([]int, error) {

	for index, refSurfaceIds := range request.LaneSurface {
		checkRefId := true
		for _, id := range refSurfaceIds {
			if id != 0 {
				checkRefId = false
				break
			}
		}
		if checkRefId {
			kmStart := helpers.FormatKM(int64(request.KmStart[index]))
			kmEnd := helpers.FormatKM(int64(request.KmEnd[index]))
			errorRes := fmt.Sprintf(`%v %v - %v`, constants.INVALID_SURFACE_ID, kmStart, kmEnd)
			return []int{}, responses.NewAppErr(400, errorRes)
		} else {
			continue
		}
	}

	order := make(map[int]int)
	arrKmStart := make([]float64, len(request.KmStart))
	copy(arrKmStart, request.KmStart)
	direction, err := r.roadSurfaceRepo.FindRefDirectionIDByRoadID(request.RoadId)
	if err != nil {
		logs.Error(err)
		return []int{}, err
	}
	sort.Float64s(arrKmStart)
	if direction == 2 {
		for i, j := 0, len(arrKmStart)-1; i < j; i, j = i+1, j-1 {
			arrKmStart[i], arrKmStart[j] = arrKmStart[j], arrKmStart[i]
		}
	}

	for i, v := range arrKmStart {
		order[i] = FindIndexOfSliceFloat64(v, request.KmStart)
	}
	currentYear := time.Now().Year()

	revision, err := r.roadSurfaceRepo.GetMaxRevisionByRoadID(request.RoadId)
	if err != nil {
		logs.Error(err)
		return []int{}, err
	}

	err = r.checkOverlap(request)
	if err != nil {
		return []int{}, err
	}

	roadInfo, err := r.roadSurfaceRepo.GetLastRoadInfoByID(request.RoadId)
	if err != nil {
		return []int{}, err
	}

	if request.KmStart[0] != float64(roadInfo.KmStart) || request.KmEnd[len(request.KmEnd)-1] != float64(roadInfo.KmEnd) {
		return []int{}, errors.New(constants.INVALID_GEOM_RANGE)
	}

	tx := r.roadSurfaceRepo.StartTransSection()
	// err = r.roadSurfaceRepo.UpdateStatusTToDeleteAll(request.RoadId, tx)
	// if err != nil {
	// 	r.roadSurfaceRepo.RollBack(tx)
	// 	logs.Error(err)
	// 	return []int{}, err
	// }
	if len(request.MaterialBase) == 0 {
		request.MaterialBase = make([]*int, len(request.KmStart))
	}
	if len(request.MaterialSubbase) == 0 {
		request.MaterialSubbase = make([]*int, len(request.KmStart))
	}
	if len(request.MaterialSubgrade) == 0 {
		request.MaterialSubgrade = make([]*int, len(request.KmStart))
	}
	if len(request.ThicknessBase) == 0 {
		request.ThicknessBase = make([]*float64, len(request.KmStart))
	}
	if len(request.ThicknessSubbase) == 0 {
		request.ThicknessSubbase = make([]*float64, len(request.KmStart))
	}
	if len(request.ThicknessSubgrade) == 0 {
		request.ThicknessSubgrade = make([]*float64, len(request.KmStart))
	}
	if len(request.ThicknessConcreteSlab) == 0 {
		request.ThicknessConcreteSlab = make([]*float64, len(request.KmStart))
	}

	if len(request.ThicknessSurfaceConcrete) == 0 {
		request.ThicknessSurfaceConcrete = make([]*float64, len(request.KmStart))
	}

	// find max sruface grpup
	srufaceGrp, err := r.roadSurfaceRepo.GetRoadSurfaceGroupByRoadID(request.RoadId)
	if err != nil {
		r.roadSurfaceRepo.RollBack(tx)

		return []int{}, responses.NewAppErr(400, err.Error())
	}
	maxSrufaceGrp := srufaceGrp.SurfaceGroup + 1

	// get all ref surface param
	refSurfaceParam, err := r.roadSurfaceRepo.GetRefSurfaceParam()
	if err != nil {
		r.roadSurfaceRepo.RollBack(tx)

		return []int{}, responses.NewAppErr(400, err.Error())
	}
	refSurfaceParamData := make(map[int]int)
	for _, item := range refSurfaceParam {
		refSurfaceParamData[item.RefSurfaceID] = item.ID
	}

	for index := range request.KmStart {
		i := order[index]
		var dataInsert models.RoadSurfacePointer
		err = copier.Copy(&dataInsert, &request)
		if err != nil {
			r.roadSurfaceRepo.RollBack(tx)
			logs.Error(err)
			return []int{}, err
		}
		dataInsert.CreatedBy = uid
		dataInsert.Year = currentYear
		dataInsert.Revision = revision + 1
		dataInsert.Status = "A"
		dataInsert.CreatedDate = time.Now()
		dataInsert.KmStart = request.KmStart[i]
		dataInsert.KmEnd = request.KmEnd[i]
		dataInsert.SurfaceCrosssectionCode = request.SurfaceCrosssectionCode[i]
		dataInsert.WidthSurface = (request.WidthSurface[i])
		dataInsert.ThicknessSurface = (request.ThicknessSurface[i])
		dataInsert.ThicknessSurfaceConcrete = (request.ThicknessSurfaceConcrete[i])
		dataInsert.WidthShoulderLeft = request.WidthShoulderLeft[i]
		dataInsert.RefSurfaceIdShoulderLeft = request.RefSurfaceIdShoulderLeft[i]
		dataInsert.WidthShoulderRight = request.WidthShoulderRight[i]
		dataInsert.RefSurfaceIdShoulderRight = request.RefSurfaceIdShoulderRight[i]
		dataInsert.ThicknessBase = request.ThicknessBase[i]
		dataInsert.RefMaterialBaseId = request.MaterialBase[i]
		dataInsert.ThicknessSubbase = request.ThicknessSubbase[i]
		dataInsert.RefMaterialSubbaseId = request.MaterialSubbase[i]
		dataInsert.ThicknessSubgrade = request.ThicknessSubgrade[i]
		dataInsert.RefMaterialSubgradeId = request.MaterialSubgrade[i]
		dataInsert.ThicknessConcreteSlab = request.ThicknessConcreteSlab[i]
		dataInsert.SurfaceGroup = maxSrufaceGrp
		// err := r.roadSurfaceRepo.ClearPreviousData(request.RoadId, uid, request.KmStart[i], request.KmEnd[i], len(request.LaneSurface[i]), tx)
		// if err != nil {
		// 	r.roadSurfaceRepo.RollBack(tx)
		// 	logs.Error(err)
		// 	return []int{}, err
		// }
		err := r.roadSurfaceRepo.ClearPreviousDataStatus(request.RoadId, uid, tx)
		if err != nil {
			r.roadSurfaceRepo.RollBack(tx)
			logs.Error(err)
			return []int{}, err
		}

		id, err := r.roadSurfaceRepo.InsertRoadSurface(dataInsert, tx)
		if err != nil {
			r.roadSurfaceRepo.RollBack(tx)
			logs.Error(err)
			return []int{}, err
		}
		if request.IDs[index] != 0 {
			err = r.roadSurfaceRepo.UpdateIDParent(id, request.IDs[index], tx)
			if err != nil {
				r.roadSurfaceRepo.RollBack(tx)
				logs.Error(err)
				return []int{}, err
			}
		} else {
			err = r.roadSurfaceRepo.UpdateIDParent(id, id, tx)
			if err != nil {
				r.roadSurfaceRepo.RollBack(tx)
				logs.Error(err)
				return []int{}, err
			}
		}

		err = r.PostLaneData(request, id, i, tx, refSurfaceParamData)
		if err != nil {
			return []int{}, err
		}
	}

	if err := r.roadSurfaceRepo.Commit(tx); err != nil {
		logs.Error(err)
		return []int{}, err
	}
	resultID, err := r.roadSurfaceRepo.GetNewIDs(request.RoadId)
	if err != nil {
		logs.Error(err)
		return []int{}, err
	}
	return resultID, nil
}

func (r *roadSurfaceUsecase) checkOverlap(request requests.RoadSurface) error {
	direction, err := r.roadSurfaceRepo.FindRefDirectionIDByRoadID(request.RoadId)
	if err != nil {
		logs.Error(err)
		return err
	}
	// check Km overlap
	for i := 0; i < len(request.KmStart); i++ {
		for j := i + 1; j < len(request.KmStart); j++ {
			if direction == 1 {
				if request.KmStart[j] >= request.KmStart[i] && request.KmStart[j] < request.KmEnd[i] {
					err = errors.New(constants.OVERLAPPING_RANGE)
					logs.Error(err)
					return err
				}
			} else if direction == 2 {
				if request.KmStart[j] <= request.KmStart[i] && request.KmStart[j] > request.KmEnd[i] {
					err = errors.New(constants.OVERLAPPING_RANGE)
					logs.Error(err)
					return err
				}
			}
		}
	}
	return nil
}

func (r *roadSurfaceUsecase) PostLaneData(data requests.RoadSurface, surfaceID int, i int, tx *gorm.DB, refSurfaceParamData map[int]int) error {
	var subRoadStart float64
	var subRoadEnd float64
	// var subRoadMin float64
	// var subRoadMax float64
	var theGeom string
	direction, err := r.roadSurfaceRepo.FindRefDirectionIDByRoadID(data.RoadId)
	if err != nil {
		logs.Error(err)
		return err
	}
	fmt.Println("direction", direction)
	fullGeom, err := r.roadSurfaceRepo.GetGeomByRoadID(data.RoadId)
	if err != nil {
		logs.Error(err)
		return err
	}
	v := data.LaneSurface[i]
	for index, l := range v {
		for _, lg := range fullGeom {
			if index+1 == lg.LaneNo {
				switch direction {
				case 1:
					if data.KmStart[i] >= lg.KmStart {
						if data.KmStart[i] >= (lg.KmEnd) {
							subRoadStart = 0
						} else {
							subRoadStart = (math.Abs((data.KmStart[i] - (lg.KmStart)))) / (math.Abs(((lg.KmStart) - (lg.KmEnd))))
						}
					} else {
						subRoadStart = 0
					}
					if data.KmEnd[i] < lg.KmEnd {
						if data.KmEnd[i] < lg.KmStart {
							subRoadEnd = 0
						} else {
							// subRoadEnd = float64(math.Abs(float64(data.KmEnd[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
							subRoadEnd = math.Abs(data.KmEnd[i]-lg.KmStart) / math.Abs(lg.KmStart-lg.KmEnd)
						}
					} else {
						subRoadEnd = 1
					}

				case 2:
					if data.KmStart[i] <= lg.KmStart {
						if data.KmStart[i] < float64(lg.KmEnd) {
							subRoadStart = 0
						} else {
							// subRoadStart = float64(math.Abs(float64(data.KmStart[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
							subRoadStart = math.Abs(data.KmStart[i]-lg.KmStart) / math.Abs(lg.KmStart-lg.KmEnd)
						}
					} else {
						subRoadStart = 1
					}

					if data.KmEnd[i] >= lg.KmEnd {
						if data.KmEnd[i] >= lg.KmStart {
							subRoadEnd = 0
						} else {
							// subRoadEnd = float64(math.Abs(float64(data.KmEnd[i]-float64(lg.KmStart)))) / float64(math.Abs(float64(float64(lg.KmStart)-float64(lg.KmEnd))))
							subRoadEnd = math.Abs(data.KmEnd[i]-lg.KmStart) / math.Abs(lg.KmStart-lg.KmEnd)
						}
					} else {
						subRoadEnd = 0
					}
				}

				subRoadMin := math.Min(subRoadStart, subRoadEnd)
				subRoadMax := math.Max(subRoadStart, subRoadEnd)
				theGeom = fmt.Sprintf("ST_LineSubstring('%s', %f, %f)", lg.TheGeom, subRoadMin, subRoadMax)
				// if subRoadEnd > subRoadStart {
				// 	subRoadMin = subRoadStart
				// 	subRoadMax = subRoadEnd
				// } else {
				// 	subRoadMin = subRoadEnd
				// 	subRoadMax = subRoadStart
				// }
				// helpers.PrintlnJson(theGeom)
				// theGeom = fmt.Sprintf("ST_LineSubstring('%s', %f, %f)", lg.TheGeom, subRoadMin, subRoadMax)
				// break
			}
		}
		if theGeom == "" {
			err = errors.New(constants.INVALID_GEOM_RANGE)
			logs.Error(err)
			return err
		}
		var dataInsert models.RoadSurfaceLane
		dataInsert.RoadSurfaceId = surfaceID
		dataInsert.LaneNo = index + 1
		dataInsert.RefSurfaceId = l
		dataInsert.TheGeom = theGeom
		dataInsert.RoadId = data.RoadId
		_, isval := refSurfaceParamData[l]
		if isval {
			dataInsert.RefSurfaceParamsID = refSurfaceParamData[l]
		}
		err = r.roadSurfaceRepo.InsertRoadSurfaceLane(dataInsert, tx)
		// if i == 1 {
		// 	err  = errors.New("test rollback")
		// }
		if err != nil {
			r.roadSurfaceRepo.RollBack(tx)
			logs.Error(err)
			return err
		}
	}

	return nil
}

func (r *roadSurfaceUsecase) GetTotalKm(roadID int) (float64, error) {
	result, err := r.roadSurfaceRepo.GetTotalKm(roadID)
	if err != nil {
		logs.Error(err)
		return result, err
	}
	return result, nil
}

func FindIndexOfSliceFloat64(value float64, slice []float64) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

func (r *roadSurfaceUsecase) GetRoadSurfaceIconById(roadId int) ([]models.RoadSurfaceIcon, error) {
	result, err := r.roadSurfaceRepo.GetRoadSurfaceIconById(roadId)
	if err != nil {
		logs.Error(err)
		return result, err
	}
	return result, nil
}
