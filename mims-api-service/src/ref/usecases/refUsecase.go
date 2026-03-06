package usecases

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/jinzhu/copier"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	servicesDB "gitlab.com/mims-api-service/services/database"
	"gitlab.com/mims-api-service/src/ref/domains"
)

type refUseCase struct {
	refRepo    domains.RefRepository
	servicesDB servicesDB.ServicesDatabaseDomain
}

func NewRefUseCase(repo domains.RefRepository, servicesDB servicesDB.ServicesDatabaseDomain) domains.RefUseCase {
	return &refUseCase{refRepo: repo, servicesDB: servicesDB}
}

func (ru *refUseCase) GetRefAsset() (interface{}, error) {
	model := []models.RefAsset{}

	err := ru.refRepo.GetRefStatus(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetPosition() (interface{}, error) {
	model := []models.RefAssetPosition{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetArea() (interface{}, error) {
	model := []models.RefAssetArea{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetGuardrail() (interface{}, error) {
	model := []models.RefAssetGuardrail{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetReflecType() (interface{}, error) {
	model := []models.RefAssetReflecType{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetOwner() (interface{}, error) {
	model := []models.RefAssetOwner{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetLightType() (interface{}, error) {
	model := []models.RefAssetLightType{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetLightWatt() (interface{}, error) {
	model := []models.RefAssetLightWatt{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetCleranceType() (interface{}, error) {
	model := []models.RefAssetCleranceType{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetCrashcushionType() (interface{}, error) {
	model := []models.RefAssetCrashcushionType{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetFenceType() (interface{}, error) {
	model := []models.RefAssetFenceType{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetKmstoneType() (interface{}, error) {
	model := []models.RefAssetKmstoneType{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetQutterType() (interface{}, error) {
	model := []models.RefAssetQutterType{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetTrafficCameraType() (interface{}, error) {
	model := []models.RefAssetTrafficCameraType{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetWeightStationType() (interface{}, error) {
	model := []models.RefAssetWeightStationType{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetBuildingType() (interface{}, error) {
	model := []models.RefAssetBuildingType{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetNoiseBarrier() (interface{}, error) {
	model := []models.RefAssetNoiseBarrier{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetSign() (interface{}, error) {
	model := []models.RefAssetSign{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetSignImage() (interface{}, error) {
	model := []models.RefAssetSignImage{}

	err := ru.refRepo.GetRefAssetSignImage(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}
	datas := []models.RefAssetSignImage{}
	for _, item := range model {
		data := models.RefAssetSignImage{}
		copier.Copy(&data, &item)
		data.SignImageFilepath = os.Getenv("STORAGE_IP") + "/" + item.SignImageFilepath
		datas = append(datas, data)
	}

	return datas, nil
}

func (ru *refUseCase) GetRefAssetSignType() (interface{}, error) {
	model := []models.RefAssetSignType{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetTable() (interface{}, error) {
	model := []models.RefAssetTable{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetTableColumns() (interface{}, error) {
	model := []models.RefAssetTableColumns{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefAssetTableStaff() (interface{}, error) {
	model := []models.RefAssetTableStaff{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefDataStatus() (interface{}, error) {
	model := []models.RefDataStatus{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefDepartment() (interface{}, error) {
	model := []models.RefDepartment{}

	err := ru.refRepo.GetRefStatus(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefDirection() (interface{}, error) {
	model := []models.RefDirection{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefGrade() (interface{}, error) {
	model := []models.RefGrade{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefMaterialBase() (interface{}, error) {
	model := []models.RefMaterialBase{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefMaterialSubbase() (interface{}, error) {
	model := []models.RefMaterialSubbase{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefMaterialSubgrade() (interface{}, error) {
	model := []models.RefMaterialSubgrade{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefOwner() (interface{}, error) {
	model := []models.RefOwner{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	resp := []models.RefOwner{}
	for _, item := range model {
		if item.IsActive {
			resp = append(resp, item)
		}
	}

	return resp, nil
}

func (ru *refUseCase) GetRefRoadType() (interface{}, error) {
	model := []models.RefRoadType{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefSurface() (interface{}, error) {

	surfaces := []models.RefSurface{}
	err := ru.refRepo.GetRef(&surfaces)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return surfaces, nil
}

func (ru *refUseCase) GetRefSurfaceType() (interface{}, error) {

	surfacesType := []models.RefSurfaceType{}
	err := ru.refRepo.GetRef(&surfacesType)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return surfacesType, nil
}

func (ru *refUseCase) GetRefSurfaceGroup() (interface{}, error) {

	var resps []responses.RefSurfaceGroup

	surfaces := []models.RefSurface{}
	err := ru.refRepo.GetRef(&surfaces)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	uniqueGroups := make(map[string]bool)
	for _, surface := range surfaces {
		var resp responses.RefSurfaceGroup
		resp.Name = surface.SurfaceGroup
		if surface.SurfaceGroup == "Asphalt" {
			resp.ColorCode = "#7239ea"
		} else if surface.SurfaceGroup == "Concrete" {
			resp.ColorCode = "#0e4285"
		}
		if !uniqueGroups[surface.SurfaceGroup] {
			uniqueGroups[surface.SurfaceGroup] = true
			resps = append(resps, resp)
		}
	}

	return resps, nil
}

func (ru *refUseCase) GetRefStructureSurface() (interface{}, error) {

	model := []models.RefStructureSurface{}

	err := ru.refRepo.GetRef(&model)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return model, nil
}

func (ru *refUseCase) GetRefTableList() ([]models.RefTableList, error) {
	tableLists, err := ru.refRepo.GetRefTableList()
	if err != nil {
		return []models.RefTableList{}, responses.NewInternalServerError()
	}

	return tableLists, nil
}

func (ru *refUseCase) GetRefColorList() (interface{}, error) {
	tableLists, err := ru.refRepo.GetRefColorList()
	if err != nil {
		return []models.RefTableList{}, responses.NewInternalServerError()
	}

	return tableLists, nil
}

func (ru *refUseCase) GetRefRoadTypeIcon() (interface{}, error) {
	icons, err := ru.refRepo.GetRefRoadTypeIcon()
	if err != nil {
		return []models.RefTableList{}, responses.NewInternalServerError()
	}

	return icons, nil
}

func (ru *refUseCase) GetRoadConditionGrades() (interface{}, error) {
	var owners []models.RefOwner
	err := ru.refRepo.GetDataList(&owners, "is_active = true")
	if err != nil {
		return "", responses.NewInternalServerError()
	}
	var resps []responses.ConditionRespondInit
	for _, v := range owners {
		resp, err := ru.GetConditionList(v.ID)
		if err != nil {
			return "", responses.NewInternalServerError()
		}
		resps = append(resps, resp)
	}
	// grads, err := ru.refRepo.GetParamsCondition()
	// if err != nil {
	// 	return resps, responses.NewInternalServerError()
	// }

	// groupedResponses := make(map[int]responses.ConditionGradeResponse)
	// conditionTypes := make(map[int]map[string][]responses.ConditionListNew)
	// for _, grad := range grads {
	// 	if !grad.RefOwner.IsActive {
	// 		continue
	// 	} else {
	// 		resp, exists := groupedResponses[grad.RefOwner.ID]
	// 		if !exists {
	// 			resp = responses.ConditionGradeResponse{
	// 				ID:   grad.RefOwner.ID,
	// 				Name: grad.RefOwner.Name,
	// 			}
	// 			groupedResponses[grad.RefOwner.ID] = resp
	// 		}

	// 		conditionGroup, exists := conditionTypes[grad.RefOwner.ID]
	// 		if !exists {
	// 			conditionGroup = make(map[string][]responses.ConditionListNew)
	// 			conditionTypes[grad.RefOwner.ID] = conditionGroup
	// 		}

	// 		condition := responses.ConditionGrade{
	// 			Grade:            grad.RefGrade,
	// 			LeftValueAC:      grad.LeftValueAC,
	// 			LeftConditionAC:  grad.LeftConditionAC,
	// 			RightValueAC:     grad.RightValueAC,
	// 			RightConditionAC: grad.RightConditionAC,
	// 			LeftValueCC:      grad.LeftValueCC,
	// 			LeftConditionCC:  grad.LeftConditionCC,
	// 			RightValueCC:     grad.RightValueCC,
	// 			RightConditionCC: grad.RightConditionCC,
	// 		}

	// 		conditionGroup[grad.ConditionType] = append(conditionGroup[grad.ConditionType], condition)

	// 	}

	// }

	// for id, resp := range groupedResponses {
	// 	for conditionType, conditions := range conditionTypes[id] {
	// 		resp.ConditionGroups = append(resp.ConditionGroups, responses.ConditionGroup{
	// 			ConditionType: conditionType,
	// 			Conditions:    conditions,
	// 		})
	// 	}
	// 	resps = append(resps, resp)
	// }

	sort.SliceStable(resps, func(i, j int) bool {
		return resps[i].ID < resps[j].ID
	})

	return resps, nil
}
func GetConditionType(condition models.ParamsConditionPreload) responses.ConditionType {
	return responses.ConditionType{
		Grade:            condition.RefGrade,
		LeftValueAC:      condition.LeftValueAC,
		LeftConditionAC:  condition.LeftConditionAC,
		RightValueAC:     condition.RightValueAC,
		RightConditionAC: condition.RightConditionAC,
		LeftValueCC:      condition.LeftValueCC,
		LeftConditionCC:  condition.LeftConditionCC,
		RightValueCC:     condition.RightValueCC,
		RightConditionCC: condition.RightConditionCC,
	}
}
func GetCondition(paramsCondition []models.ParamsConditionPreload) []responses.ConditionListNew {
	ConditionListNews := []responses.ConditionListNew{}
	if len(paramsCondition) == 0 {
		return ConditionListNews
	}
	condition := responses.Condition{}
	for _, v := range paramsCondition {

		if !v.RefOwner.IsActive {
			continue
		}

		switch v.ConditionType {
		case "IFI":
			ifi := GetConditionType(v)
			condition.IFI = append(condition.IFI, ifi)
		case "IRI":
			iri := GetConditionType(v)
			condition.IRI = append(condition.IRI, iri)
		case "MPD":
			mpd := GetConditionType(v)
			condition.MPD = append(condition.MPD, mpd)
		case "RUT":
			rut := GetConditionType(v)
			condition.RUT = append(condition.RUT, rut)
		default:
			log.Println("there is no condition type")
		}
	}
	if len(condition.IFI) != 0 {
		ConditionListNew := GenRespond(condition, "IFI")
		ConditionListNews = append(ConditionListNews, ConditionListNew)
	}
	if len(condition.IRI) != 0 {
		ConditionListNew := GenRespond(condition, "IRI")
		ConditionListNews = append(ConditionListNews, ConditionListNew)
	}
	if len(condition.MPD) != 0 {
		ConditionListNew := GenRespond(condition, "MPD")
		ConditionListNews = append(ConditionListNews, ConditionListNew)
	}
	if len(condition.RUT) != 0 {
		ConditionListNew := GenRespond(condition, "RUT")
		ConditionListNews = append(ConditionListNews, ConditionListNew)
	}
	return ConditionListNews
}

func GenRespond(condition responses.Condition, conditionType string) responses.ConditionListNew {
	ConditionListNew := responses.ConditionListNew{ConditionType: conditionType}
	SurfaceTypeCondition := responses.SurfaceTypeCondition{}
	AC := []responses.ConditionTypeAC{}
	CC := []responses.ConditionTypeCC{}
	switch conditionType {
	case "IFI":
		for _, v := range condition.IFI {
			ac := responses.ConditionTypeAC{Grade: v.Grade,
				LeftValueAC:      v.LeftValueAC,
				LeftConditionAC:  v.LeftConditionAC,
				RightValueAC:     v.RightValueAC,
				RightConditionAC: v.RightConditionAC}
			cc := responses.ConditionTypeCC{Grade: v.Grade,
				LeftValueCC:      v.LeftValueCC,
				LeftConditionCC:  v.LeftConditionCC,
				RightValueCC:     v.RightValueCC,
				RightConditionCC: v.RightConditionCC}
			AC = append(AC, ac)
			CC = append(CC, cc)

		}
	case "IRI":
		for _, v := range condition.IRI {
			ac := responses.ConditionTypeAC{Grade: v.Grade,
				LeftValueAC:      v.LeftValueAC,
				LeftConditionAC:  v.LeftConditionAC,
				RightValueAC:     v.RightValueAC,
				RightConditionAC: v.RightConditionAC}
			cc := responses.ConditionTypeCC{Grade: v.Grade,
				LeftValueCC:      v.LeftValueCC,
				LeftConditionCC:  v.LeftConditionCC,
				RightValueCC:     v.RightValueCC,
				RightConditionCC: v.RightConditionCC}
			AC = append(AC, ac)
			CC = append(CC, cc)

		}
	case "MPD":
		for _, v := range condition.MPD {
			ac := responses.ConditionTypeAC{Grade: v.Grade,
				LeftValueAC:      v.LeftValueAC,
				LeftConditionAC:  v.LeftConditionAC,
				RightValueAC:     v.RightValueAC,
				RightConditionAC: v.RightConditionAC}
			cc := responses.ConditionTypeCC{Grade: v.Grade,
				LeftValueCC:      v.LeftValueCC,
				LeftConditionCC:  v.LeftConditionCC,
				RightValueCC:     v.RightValueCC,
				RightConditionCC: v.RightConditionCC}
			AC = append(AC, ac)
			CC = append(CC, cc)

		}
	case "RUT":
		for _, v := range condition.RUT {
			ac := responses.ConditionTypeAC{Grade: v.Grade,
				LeftValueAC:      v.LeftValueAC,
				LeftConditionAC:  v.LeftConditionAC,
				RightValueAC:     v.RightValueAC,
				RightConditionAC: v.RightConditionAC}
			cc := responses.ConditionTypeCC{Grade: v.Grade,
				LeftValueCC:      v.LeftValueCC,
				LeftConditionCC:  v.LeftConditionCC,
				RightValueCC:     v.RightValueCC,
				RightConditionCC: v.RightConditionCC}
			AC = append(AC, ac)
			CC = append(CC, cc)

		}
	}

	SurfaceTypeCondition.AC = AC
	SurfaceTypeCondition.CC = CC
	ConditionListNew.SurfaceType = SurfaceTypeCondition
	return ConditionListNew
}

func GenRespondInit(condition responses.Condition, conditionType string) responses.ConditionListNewInit {
	ConditionListNew := responses.ConditionListNewInit{ConditionType: conditionType}
	SurfaceTypeCondition := responses.SurfaceTypeConditionInit{}
	AC := []responses.ConditionTypeInit{}
	CC := []responses.ConditionTypeInit{}
	switch conditionType {
	case "IFI":
		for _, v := range condition.IFI {
			ac := responses.ConditionTypeInit{Grade: v.Grade,
				LeftValue:      v.LeftValueAC,
				LeftCondition:  v.LeftConditionAC,
				RightValue:     v.RightValueAC,
				RightCondition: v.RightConditionAC}
			cc := responses.ConditionTypeInit{Grade: v.Grade,
				LeftValue:      v.LeftValueCC,
				LeftCondition:  v.LeftConditionCC,
				RightValue:     v.RightValueCC,
				RightCondition: v.RightConditionCC}
			AC = append(AC, ac)
			CC = append(CC, cc)

		}
	case "IRI":
		for _, v := range condition.IRI {
			ac := responses.ConditionTypeInit{Grade: v.Grade,
				LeftValue:      v.LeftValueAC,
				LeftCondition:  v.LeftConditionAC,
				RightValue:     v.RightValueAC,
				RightCondition: v.RightConditionAC}
			cc := responses.ConditionTypeInit{Grade: v.Grade,
				LeftValue:      v.LeftValueCC,
				LeftCondition:  v.LeftConditionCC,
				RightValue:     v.RightValueCC,
				RightCondition: v.RightConditionCC}
			AC = append(AC, ac)
			CC = append(CC, cc)

		}
	case "MPD":
		for _, v := range condition.MPD {
			ac := responses.ConditionTypeInit{Grade: v.Grade,
				LeftValue:      v.LeftValueAC,
				LeftCondition:  v.LeftConditionAC,
				RightValue:     v.RightValueAC,
				RightCondition: v.RightConditionAC}
			cc := responses.ConditionTypeInit{Grade: v.Grade,
				LeftValue:      v.LeftValueCC,
				LeftCondition:  v.LeftConditionCC,
				RightValue:     v.RightValueCC,
				RightCondition: v.RightConditionCC}
			AC = append(AC, ac)
			CC = append(CC, cc)

		}
	case "RUT":
		for _, v := range condition.RUT {
			ac := responses.ConditionTypeInit{Grade: v.Grade,
				LeftValue:      v.LeftValueAC,
				LeftCondition:  v.LeftConditionAC,
				RightValue:     v.RightValueAC,
				RightCondition: v.RightConditionAC}
			cc := responses.ConditionTypeInit{Grade: v.Grade,
				LeftValue:      v.LeftValueCC,
				LeftCondition:  v.LeftConditionCC,
				RightValue:     v.RightValueCC,
				RightCondition: v.RightConditionCC}
			AC = append(AC, ac)
			CC = append(CC, cc)

		}
	}

	SurfaceTypeCondition.AC = AC
	SurfaceTypeCondition.CC = CC
	ConditionListNew.SurfaceType = SurfaceTypeCondition
	return ConditionListNew
}

func GetConditionInit(paramsCondition []models.ParamsConditionPreload) []responses.ConditionListNewInit {
	ConditionListNews := []responses.ConditionListNewInit{}
	if len(paramsCondition) == 0 {
		return ConditionListNews
	}
	condition := responses.Condition{}
	for _, v := range paramsCondition {

		if !v.RefOwner.IsActive {
			continue
		}

		switch v.ConditionType {
		case "IFI":
			ifi := GetConditionType(v)
			condition.IFI = append(condition.IFI, ifi)
		case "IRI":
			iri := GetConditionType(v)
			condition.IRI = append(condition.IRI, iri)
		case "MPD":
			mpd := GetConditionType(v)
			condition.MPD = append(condition.MPD, mpd)
		case "RUT":
			rut := GetConditionType(v)
			condition.RUT = append(condition.RUT, rut)
		default:
			log.Println("there is no condition type")
		}
	}
	if len(condition.IFI) != 0 {
		ConditionListNew := GenRespondInit(condition, "IFI")
		ConditionListNews = append(ConditionListNews, ConditionListNew)
	}
	if len(condition.IRI) != 0 {
		ConditionListNew := GenRespondInit(condition, "IRI")
		ConditionListNews = append(ConditionListNews, ConditionListNew)
	}
	if len(condition.MPD) != 0 {
		ConditionListNew := GenRespondInit(condition, "MPD")
		ConditionListNews = append(ConditionListNews, ConditionListNew)
	}
	if len(condition.RUT) != 0 {
		ConditionListNew := GenRespondInit(condition, "RUT")
		ConditionListNews = append(ConditionListNews, ConditionListNew)
	}
	return ConditionListNews
}

func (su *refUseCase) GetConditionList(ownerId int) (responses.ConditionRespondInit, error) {
	paramsCondition, err := su.refRepo.GetParamsCondition(ownerId)
	if err != nil {
		return responses.ConditionRespondInit{}, responses.NewInternalServerError()
	}
	if len(paramsCondition) == 0 {
		return responses.ConditionRespondInit{}, nil
	}

	data := responses.ConditionRespondInit{
		ID:                  paramsCondition[0].RefOwnerID,
		RefConditionRangeID: paramsCondition[0].RefOwner.RefConditionRangeID,
		OwnerName:           paramsCondition[0].RefOwner.Name,
		ConditionList:       GetConditionInit(paramsCondition),
	}
	return data, nil
}

func (su *refUseCase) GetRoadLineList(ownerId int) (responses.RoadLineListInit, error) {
	paramsCondition, err := su.refRepo.GetParamsRoadLine(ownerId)
	if err != nil {
		return responses.RoadLineListInit{}, responses.NewInternalServerError()
	}
	if len(paramsCondition) == 0 {
		return responses.RoadLineListInit{}, nil
	}
	responds := responses.RoadLineListInit{
		ID:                     paramsCondition[0].RefOwnerRoadLineID,
		RefReflectivityRangeID: paramsCondition[0].RefOwnerRoadLine.RefReflectivityRangeID,
		OwnerName:              paramsCondition[0].RefOwnerRoadLine.Name,
	}
	SurfaceTypeRoadLine := responses.SurfaceTypeRoadLineInit{}
	yellow := []responses.ConditionTypeInit{}
	white := []responses.ConditionTypeInit{}
	for _, condition := range paramsCondition {
		resYellow := responses.ConditionTypeInit{
			Grade:          condition.RefGrade,
			LeftValue:      condition.LeftValueYellow,
			LeftCondition:  condition.LeftConditionYellow,
			RightValue:     condition.RightValueYellow,
			RightCondition: condition.RightConditionYellow,
		}
		resWhite := responses.ConditionTypeInit{
			Grade:          condition.RefGrade,
			LeftValue:      condition.LeftValueWhite,
			LeftCondition:  condition.LeftConditionWhite,
			RightValue:     condition.RightValueWhite,
			RightCondition: condition.RightConditionWhite,
		}
		yellow = append(yellow, resYellow)
		white = append(white, resWhite)
	}
	SurfaceTypeRoadLine.White = white
	SurfaceTypeRoadLine.Yellow = yellow
	responds.RoadLine = SurfaceTypeRoadLine
	return responds, nil
}

func (su *refUseCase) GetRoadLineListInit(ownerId int) (responses.RoadLineList, error) {
	paramsCondition, err := su.refRepo.GetParamsRoadLine(ownerId)
	if err != nil {
		return responses.RoadLineList{}, responses.NewInternalServerError()
	}
	if len(paramsCondition) == 0 {
		return responses.RoadLineList{}, nil
	}
	responds := responses.RoadLineList{
		RefReflectivityRangeID: paramsCondition[0].RefOwnerRoadLine.RefReflectivityRangeID,
		OwnerName:              paramsCondition[0].RefOwnerRoadLine.Name,
	}
	SurfaceTypeRoadLine := responses.SurfaceTypeRoadLine{}
	yellow := []responses.ConditionTypeYellow{}
	white := []responses.ConditionTypeWhite{}
	for _, condition := range paramsCondition {
		resYellow := responses.ConditionTypeYellow{
			Grade:                condition.RefGrade,
			LeftValueYellow:      condition.LeftValueYellow,
			LeftConditionYellow:  condition.LeftConditionYellow,
			RightValueYellow:     condition.RightValueYellow,
			RightConditionYellow: condition.RightConditionYellow,
		}
		resWhite := responses.ConditionTypeWhite{
			Grade:               condition.RefGrade,
			LeftValueWhite:      condition.LeftValueWhite,
			LeftConditionWhite:  condition.LeftConditionWhite,
			RightValueWhite:     condition.RightValueWhite,
			RightConditionWhite: condition.RightConditionWhite,
		}
		yellow = append(yellow, resYellow)
		white = append(white, resWhite)
	}
	SurfaceTypeRoadLine.White = white
	SurfaceTypeRoadLine.Yellow = yellow
	responds.RoadLine = SurfaceTypeRoadLine
	return responds, nil
}

func (ru *refUseCase) GetRoadLineConditionGrades() (interface{}, error) {
	var owners []models.RefOwnerRoadLineInit
	err := ru.refRepo.GetDataList(&owners, "is_active = true")
	if err != nil {
		return "", responses.NewInternalServerError()
	}
	var resps []responses.RoadLineListInit
	for _, v := range owners {
		resp, err := ru.GetRoadLineList(v.ID)
		if err != nil {
			return "", responses.NewInternalServerError()
		}
		if resp.ID == 0 {
			continue
		}
		resps = append(resps, resp)
	}
	if resps == nil {
		return []responses.RoadLineListInit{}, nil
	}
	return resps, nil
}

func (ru *refUseCase) GetRefCriteriaType() (interface{}, error) {
	criteriaType, err := ru.refRepo.GetRefCriteriaType()
	if err != nil {
		return []models.RefCriteriaType{}, responses.NewInternalServerError()
	}

	return criteriaType, nil
}

func (ru *refUseCase) GetRefRoadConditionRange() (interface{}, error) {
	result := []models.RefConditionRange{}
	err := ru.refRepo.GetDataList(&result, "")
	if err != nil {
		return []models.RefConditionRange{}, responses.NewInternalServerError()
	}

	return result, nil
}

func (ru *refUseCase) GetRefReflectivityRange() (interface{}, error) {
	result := []models.RefReflectivityRange{}
	err := ru.refRepo.GetDataList(&result, "")
	if err != nil {
		return []models.RefReflectivityRange{}, responses.NewInternalServerError()
	}

	return result, nil
}

func (ru *refUseCase) GetParameterVehicleTypeList(road_group_id string) ([]responses.RefAadtParameterVehicleType, error) {
	idInteger, err := helpers.ConvertStringToInt(road_group_id)
	tableLists, err := ru.refRepo.GetParameterVehicleTypeListByRoadGroupId(idInteger)
	if err != nil {
		return []responses.RefAadtParameterVehicleType{}, responses.NewInternalServerError()
	}
	var refAadtParameterVehicleTypeList []responses.RefAadtParameterVehicleType
	for _, item := range tableLists {
		var refAadtParameterVehicleType responses.RefAadtParameterVehicleType
		refAadtParameterVehicleType.Id = item.Id
		refAadtParameterVehicleType.Name = item.Name
		refAadtParameterVehicleType.NumAxle = item.NumAxle
		refAadtParameterVehicleType.NumWheel = item.NumWheel
		refAadtParameterVehicleType.LoadEquivalent = item.LoadEquivalent
		refAadtParameterVehicleType.ImagePath = os.Getenv("STORAGE_IP") + "/" + item.ImagePath
		refAadtParameterVehicleTypeList = append(refAadtParameterVehicleTypeList, refAadtParameterVehicleType)
	}

	return refAadtParameterVehicleTypeList, nil
}

func (ru *refUseCase) GetRoadUserCostAcc() (interface{}, error) {
	getRoadUserCostAcc, err := ru.refRepo.GetRoadUserCostAcc()
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return getRoadUserCostAcc, nil
}

func (ru *refUseCase) GetRoadUserCostRuc() (interface{}, error) {
	getRoadUserCostRuc, err := ru.refRepo.GetRoadUserCostRuc()
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return getRoadUserCostRuc, nil
}

func (ru *refUseCase) GetMaintenanceAnalysisStrategicBudgetType() (interface{}, error) {
	getRoadUserCostRuc, err := ru.refRepo.GetMaintenanceAnalysisStrategicBudgetType()
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return getRoadUserCostRuc, nil
}

func (ru *refUseCase) GetMaintenanceAnalysisStrategicTargetType() (interface{}, error) {
	getRoadUserCostRuc, err := ru.refRepo.GetMaintenanceAnalysisStrategicTargetType()
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return getRoadUserCostRuc, nil
}

func (ru *refUseCase) GetMaintenanceAnalysisStrategic() (interface{}, error) {
	getRoadUserCostRuc, err := ru.refRepo.GetMaintenanceAnalysisStrategic()
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	isDuplicate := map[string]bool{}
	budgetMapId := map[string]int{}
	for _, item1 := range getRoadUserCostRuc {
		for _, item2 := range item1.Budget {
			if isDuplicate[item2.Name] == false {
				budgetMapId[item2.Name] = item2.Id
				isDuplicate[item2.Name] = true
			}
		}
	}

	for index1, _ := range getRoadUserCostRuc {
		for index2, _ := range getRoadUserCostRuc[index1].Budget {
			if budgetMapId[getRoadUserCostRuc[index1].Budget[index2].Name] != 0 {
				getRoadUserCostRuc[index1].Budget[index2].Id = budgetMapId[getRoadUserCostRuc[index1].Budget[index2].Name]
			}
		}
	}

	return getRoadUserCostRuc, nil
}

func (ru *refUseCase) GetRefDistrictsList() (interface{}, error) {
	resp, err := ru.refRepo.GetRefDistrictsList()

	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return resp, nil
}

func (ru *refUseCase) GetRoadGroupList() ([]responses.RoadGroupInitData, error) {
	var roadGroupRes []responses.RoadGroupInitData
	result, err := ru.refRepo.GetRoadGroupList()
	if err != nil {
		return nil, responses.NewInternalServerError()
	}
	for g, group := range result {
		groupNumber := strings.TrimLeft(group.Number, "0")
		result[g].Number = groupNumber
	}
	copier.Copy(&roadGroupRes, &result)
	return roadGroupRes, nil
}

func (ru *refUseCase) GetRoadSectionList(userID int) ([]models.RoadSectionInitData, error) {
	//============ start check permission ============
	userInfo, _ := ru.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_road"}, []string{"view_owner_road"})
	//============ end check permission ============
	result, err := ru.refRepo.GetRoadSectionList(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return result, nil
}

func (ru *refUseCase) GetRefDistrictsInitList(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := ru.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_road"}, []string{"view_owner_road"})
	//============ end check permission ============
	resp, err := ru.refRepo.GetRefDistrictsInitList(isAllData, isOwnerData, depotCode)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	var res []models.RefDistrictInitData
	for _, item := range resp.([]models.RefDistrictInitData) {
		var data models.RefDistrictInitData
		if len(item.Depots) > 0 {
			err := copier.Copy(&data, item)
			if err != nil {
				logs.Error(err)
				return nil, responses.NewInternalServerError()
			}
			res = append(res, data)
		}
	}

	return res, nil
}

func (ru *refUseCase) GetRefRoadTypeLevel(where string) ([]models.RefRoadTypeInit, error) {
	result := []models.RefRoadTypeInit{}
	err := ru.refRepo.GetDataList(&result, where)
	if err != nil {
		return []models.RefRoadTypeInit{}, responses.NewInternalServerError()
	}

	return result, nil
}

func (ru *refUseCase) GetRefDivisionInitList(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := ru.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_road"}, []string{"view_owner_road"})
	//============ end check permission ============
	resp, err := ru.refRepo.GetRefDivisionInitList(isAllData, isOwnerData, depotCode)
	var refDivisionLists []models.RefDivisionList
	for _, item := range resp {
		var refDivisionList models.RefDivisionList
		var refDistricts []models.RefDistrictInit
		if len(item.Districts) > 0 {
			refDivisionList.Id = item.Id
			refDivisionList.Name = item.Name
			refDivisionList.DivisionCode = item.DivisionCode
			refDivisionList.OwnerCodeKey = item.OwnerCodeKey
			for _, district := range item.Districts {
				var refDistrict models.RefDistrictInit
				if len(district.Depots) > 0 {
					var refDepots []models.RefDepotInit
					for _, depot := range district.Depots {
						if depot.Name == "" {
							continue
						}
						var refDepot models.RefDepotInit
						copier.Copy(&refDepot, &depot)
						refDepots = append(refDepots, refDepot)
					}
					copier.Copy(&refDistrict, &district)
					refDistrict.Depots = refDepots
					if len(refDepots) > 0 {
						refDistricts = append(refDistricts, refDistrict)
					}
				}
			}
			refDivisionList.Districts = refDistricts
			if len(refDistricts) > 0 {
				refDivisionLists = append(refDivisionLists, refDivisionList)
			}
		}
	}

	if err != nil {
		return nil, responses.NewInternalServerError()
	}
	if len(refDivisionLists) == 0 {
		return []string{}, nil
	}
	return refDivisionLists, nil
}

func (ru *refUseCase) GetRefDivisionInitListDashboardAsset(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := ru.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_asset_dashboard"}, []string{"view_owner_asset_dashboard"})
	//============ end check permission ============
	resp, err := ru.refRepo.GetRefDivisionInitList(isAllData, isOwnerData, depotCode)
	var refDivisionLists []models.RefDivisionList
	for _, item := range resp {
		var refDivisionList models.RefDivisionList
		var refDistricts []models.RefDistrictInit
		if len(item.Districts) > 0 {
			refDivisionList.Id = item.Id
			refDivisionList.Name = item.Name
			refDivisionList.DivisionCode = item.DivisionCode
			refDivisionList.OwnerCodeKey = item.OwnerCodeKey
			for _, district := range item.Districts {
				var refDistrict models.RefDistrictInit
				if len(district.Depots) > 0 {
					var refDepots []models.RefDepotInit
					for _, depot := range district.Depots {
						if depot.Name == "" {
							continue
						}
						var refDepot models.RefDepotInit
						copier.Copy(&refDepot, &depot)
						refDepots = append(refDepots, refDepot)
					}
					copier.Copy(&refDistrict, &district)
					refDistrict.Depots = refDepots
					if len(refDepots) > 0 {
						refDistricts = append(refDistricts, refDistrict)
					}
				}
			}
			refDivisionList.Districts = refDistricts
			if len(refDistricts) > 0 {
				refDivisionLists = append(refDivisionLists, refDivisionList)
			}
		}
	}

	if err != nil {
		return nil, responses.NewInternalServerError()
	}
	if len(refDivisionLists) == 0 {
		return []string{}, nil
	}
	return refDivisionLists, nil
}

func (ru *refUseCase) GetRefDivisionInitListDashboardCondition(userID int) (interface{}, error) {
	//============ start check permission ============
	userInfo, _ := ru.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_road_condition_dashboard"}, []string{"view_owner_road_condition_dashboard"})
	//============ end check permission ============
	resp, err := ru.refRepo.GetRefDivisionInitList(isAllData, isOwnerData, depotCode)
	var refDivisionLists []models.RefDivisionList
	for _, item := range resp {
		var refDivisionList models.RefDivisionList
		var refDistricts []models.RefDistrictInit
		if len(item.Districts) > 0 {
			refDivisionList.Id = item.Id
			refDivisionList.Name = item.Name
			refDivisionList.DivisionCode = item.DivisionCode
			refDivisionList.OwnerCodeKey = item.OwnerCodeKey
			for _, district := range item.Districts {
				var refDistrict models.RefDistrictInit
				if len(district.Depots) > 0 {
					var refDepots []models.RefDepotInit
					for _, depot := range district.Depots {
						if depot.Name == "" {
							continue
						}
						var refDepot models.RefDepotInit
						copier.Copy(&refDepot, &depot)
						refDepots = append(refDepots, refDepot)
					}
					copier.Copy(&refDistrict, &district)
					refDistrict.Depots = refDepots
					if len(refDepots) > 0 {
						refDistricts = append(refDistricts, refDistrict)
					}
				}
			}
			refDivisionList.Districts = refDistricts
			if len(refDistricts) > 0 {
				refDivisionLists = append(refDivisionLists, refDivisionList)
			}
		}
	}

	if err != nil {
		return nil, responses.NewInternalServerError()
	}
	if len(refDivisionLists) == 0 {
		return []string{}, nil
	}
	return refDivisionLists, nil
}

func (ru *refUseCase) GetRefDivisionInitListDashboardSurface(userID int) (interface{}, error) {
	// ============ start check permission ============
	userInfo, _ := ru.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_surface_dashboard"}, []string{"view_owner_surface_dashboard"})
	// ============ end check permission ============
	resp, err := ru.refRepo.GetRefDivisionInitList(isAllData, isOwnerData, depotCode)
	var refDivisionLists []models.RefDivisionList
	for _, item := range resp {
		var refDivisionList models.RefDivisionList
		var refDistricts []models.RefDistrictInit
		if len(item.Districts) > 0 {
			refDivisionList.Id = item.Id
			refDivisionList.Name = item.Name
			refDivisionList.DivisionCode = item.DivisionCode
			refDivisionList.OwnerCodeKey = item.OwnerCodeKey
			for _, district := range item.Districts {
				var refDistrict models.RefDistrictInit
				if len(district.Depots) > 0 {
					var refDepots []models.RefDepotInit
					for _, depot := range district.Depots {
						if depot.Name == "" {
							continue
						}
						var refDepot models.RefDepotInit
						copier.Copy(&refDepot, &depot)
						refDepots = append(refDepots, refDepot)
					}
					copier.Copy(&refDistrict, &district)
					refDistrict.Depots = refDepots
					if len(refDepots) > 0 {
						refDistricts = append(refDistricts, refDistrict)
					}
				}
			}
			refDivisionList.Districts = refDistricts
			if len(refDistricts) > 0 {
				refDivisionLists = append(refDivisionLists, refDivisionList)
			}
		}
	}

	if err != nil {
		return nil, responses.NewInternalServerError()
	}
	if len(refDivisionLists) == 0 {
		return []string{}, nil
	}
	return refDivisionLists, nil
}

func (ru *refUseCase) GetRefDivisionInitListDashboardMaintenance(userID int) (interface{}, error) {
	// ============ start check permission ============
	userInfo, _ := ru.servicesDB.UserInfo(userID)
	depotCode := userInfo.RefDepot.DepotCode
	refUserOwnerID := userInfo.RefUserOwnerID
	accessCtrls := userInfo.AccessControl
	accessCtrl := []string{}
	for _, item := range accessCtrls {
		accessCtrl = append(accessCtrl, item.AccessKey)
	}
	isAllData, isOwnerData := helpers.CheckPermission(refUserOwnerID, accessCtrl, []string{"view_all_maint_history_dashboard"}, []string{"view_owner_maint_history_dashboard"})
	// ============ end check permission ============
	resp, err := ru.refRepo.GetRefDivisionInitList(isAllData, isOwnerData, depotCode)
	var refDivisionLists []models.RefDivisionList
	for _, item := range resp {
		var refDivisionList models.RefDivisionList
		var refDistricts []models.RefDistrictInit
		if len(item.Districts) > 0 {
			refDivisionList.Id = item.Id
			refDivisionList.Name = item.Name
			refDivisionList.DivisionCode = item.DivisionCode
			refDivisionList.OwnerCodeKey = item.OwnerCodeKey
			for _, district := range item.Districts {
				var refDistrict models.RefDistrictInit
				if len(district.Depots) > 0 {
					var refDepots []models.RefDepotInit
					for _, depot := range district.Depots {
						if depot.Name == "" {
							continue
						}
						var refDepot models.RefDepotInit
						copier.Copy(&refDepot, &depot)
						refDepots = append(refDepots, refDepot)
					}
					copier.Copy(&refDistrict, &district)
					refDistrict.Depots = refDepots
					if len(refDepots) > 0 {
						refDistricts = append(refDistricts, refDistrict)
					}
				}
			}
			refDivisionList.Districts = refDistricts
			if len(refDistricts) > 0 {
				refDivisionLists = append(refDivisionLists, refDivisionList)
			}
		}
	}

	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	if len(refDivisionLists) == 0 {
		return []string{}, nil
	}
	return refDivisionLists, nil
}

func (ru *refUseCase) GetRefStripeColor() (interface{}, error) {

	refStripeColor := []models.RefStripeColor{}
	err := ru.refRepo.GetRef(&refStripeColor)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return refStripeColor, nil
}

func (ru *refUseCase) GetRefStripeType() (interface{}, error) {

	refStripeType := []models.RefStripeType{}
	err := ru.refRepo.GetRef(&refStripeType)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return refStripeType, nil

}

func (ru *refUseCase) GetRefCriteriaMethod() ([]models.RefCriteriaMethod, error) {
	data, err := ru.refRepo.GetRefCriteriaMethod()
	if err != nil {
		return data, responses.NewInternalServerError()
	}

	return data, nil
}

func (ru *refUseCase) GetRefUserOwner() ([]responses.RefUserOwner, error) {
	data, err := ru.refRepo.GetRefUserOwner()
	if err != nil {
		return data, responses.NewInternalServerError()
	}
	return data, nil
}
