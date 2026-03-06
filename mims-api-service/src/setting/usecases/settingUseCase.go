package usecases

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/setting/domains"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type settingUseCase struct {
	settingRepo domains.SettingRepository
}

func NewSettingUseCase(repo domains.SettingRepository) domains.SettingUseCase {
	return &settingUseCase{settingRepo: repo}
}

func (su *settingUseCase) GetAssetGroups(params requests.QueryParams) (interface{}, error) {
	total, err := su.settingRepo.CountAssetGroups(params.Name)
	if err != nil {
		return responses.Pagination{}, responses.NewInternalServerError()
	}

	limit, offset, page := helpers.GetlimitOffsetPage(params.Limit, params.Page, total)

	assetGroups, err := su.settingRepo.GetAssetGroups(limit, offset, params.Name)
	if err != nil {
		log.Println(err)
		return responses.Pagination{}, responses.NewInternalServerError()
	}

	pagination := helpers.Pagination(assetGroups, limit, page, total)

	return pagination, nil
}

func (su *settingUseCase) GetAssetGroupByID(id string) (models.RefAsset, error) {
	assetGroupId, err := helpers.ConvertStringToInt(id)
	if err != nil {
		return models.RefAsset{}, err
	}

	assetGroup := models.RefAsset{}
	err = su.settingRepo.GetByID(&assetGroup, assetGroupId)
	if err != nil {
		log.Println(err)
		if err.Error() == "record not found" {
			return models.RefAsset{}, responses.NewNotFoundError()
		}

		return models.RefAsset{}, responses.NewInternalServerError()
	}

	return assetGroup, nil
}

func (su *settingUseCase) CreateAssetGroup(assetGroupName string) error {
	assetGroups := []models.RefAsset{}
	if err := su.settingRepo.GetAll(&assetGroups); err != nil {
		return responses.NewInternalServerError()
	}

	for _, assetGroup := range assetGroups {
		if assetGroupName == assetGroup.Name {
			return responses.NewDuplicatedNameError("Name:duplicate")
		}
	}

	assetGroup := models.RefAsset{
		Name:      assetGroupName,
		Status:    1,
		CanDelete: true,
	}

	err := su.settingRepo.CreateAssetGroup(assetGroup)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}

	return nil
}

func (su *settingUseCase) UpdateAssetGroupByID(id, assetGroupName string) error {
	assetGroupId, err := helpers.ConvertStringToInt(id)
	if err != nil {
		return err
	}

	assetGroups := []models.RefAsset{}
	if err := su.settingRepo.GetAll(&assetGroups); err != nil {
		return responses.NewInternalServerError()
	}

	for _, assetGroup := range assetGroups {
		if assetGroupId != assetGroup.ID && assetGroupName == assetGroup.Name {
			return responses.NewDuplicatedNameError("Name:duplicate")
		}
	}

	err = su.settingRepo.UpdateByID("ref_asset", assetGroupName, assetGroupId)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}
	return nil
}

func (su *settingUseCase) DeleteAssetGroupByID(id string) error {
	assetGroupId, err := helpers.ConvertStringToInt(id)
	if err != nil {
		return err
	}

	err = su.settingRepo.DeleteByID("ref_asset", assetGroupId)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}
	return nil
}

func (su *settingUseCase) GetDepartments(params requests.QueryParams) (interface{}, error) {
	total, err := su.settingRepo.CountDepartments(params.Name)
	if err != nil {
		return responses.Pagination{}, responses.NewInternalServerError()
	}

	limit, offset, page := helpers.GetlimitOffsetPage(params.Limit, params.Page, total)

	departments, err := su.settingRepo.GetDepartments(limit, offset, params.Name)
	if err != nil {
		log.Println(err)
		return responses.Pagination{}, responses.NewInternalServerError()
	}

	pagination := helpers.Pagination(departments, limit, page, total)

	return pagination, nil
}

func (su *settingUseCase) GetDepartmentByID(id string) (models.RefDepartment, error) {
	deptId, err := helpers.ConvertStringToInt(id)
	if err != nil {
		return models.RefDepartment{}, err
	}

	department := models.RefDepartment{}
	err = su.settingRepo.GetByID(&department, deptId)
	if err != nil {
		log.Println(err)
		if err.Error() == "record not found" {
			return models.RefDepartment{}, responses.NewNotFoundError()
		}

		return models.RefDepartment{}, responses.NewInternalServerError()
	}

	return department, nil
}

func (su *settingUseCase) CreateDepartment(departmentName string) error {
	departments := []models.RefDepartment{}
	if err := su.settingRepo.GetAll(&departments); err != nil {
		return responses.NewInternalServerError()
	}

	for _, department := range departments {
		if departmentName == department.Name {
			return responses.NewDuplicatedNameError("Name:duplicate")
		}
	}

	department := models.RefDepartment{
		Name:      departmentName,
		Status:    1,
		CanDelete: true,
	}

	err := su.settingRepo.CreateDepartment(department)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}

	return nil
}

func (su *settingUseCase) UpdateDepartmentByID(id, departmentName string) error {
	deptId, err := helpers.ConvertStringToInt(id)
	if err != nil {
		return err
	}

	departments := []models.RefDepartment{}
	if err := su.settingRepo.GetAll(&departments); err != nil {
		return responses.NewInternalServerError()
	}

	for _, department := range departments {
		if deptId != department.ID && departmentName == department.Name {
			return responses.NewDuplicatedNameError("Name:duplicate")
		}
	}

	err = su.settingRepo.UpdateByID("ref_department", departmentName, deptId)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}
	return nil
}

func (su *settingUseCase) DeleteDepartmentByID(id string) error {
	deptId, err := helpers.ConvertStringToInt(id)
	if err != nil {
		return err
	}

	err = su.settingRepo.DeleteByID("ref_department", deptId)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}
	return nil
}

func (su *settingUseCase) GetOwnersRoadLine(params requests.QueryParamsReflectivityRange) (interface{}, error) {
	total, err := su.settingRepo.CountOwnersRoadLine(params)
	if err != nil {
		return responses.Pagination{}, responses.NewInternalServerError()
	}

	limit, offset, page := helpers.GetlimitOffsetPage(params.Limit, params.Page, total)

	owners, err := su.settingRepo.GetOwnersRoadLine(limit, offset, params)
	if err != nil {
		log.Println(err)
		return responses.Pagination{}, responses.NewInternalServerError()
	}

	pagination := helpers.Pagination(owners, limit, page, total)

	return pagination, nil
}

func (su *settingUseCase) GetOwnerRoadLineByID(ownerId int) (models.RefOwnerRoadLine, error) {
	owner, err := su.settingRepo.GetOwnerRoadLineByID(ownerId)
	if err != nil {
		if err.Error() == "record not found" {
			return models.RefOwnerRoadLine{}, responses.NewNotFoundError()
		}
		return models.RefOwnerRoadLine{}, responses.NewInternalServerError()
	}

	return owner, nil
}

func (su *settingUseCase) CreateOwnerRoadLine(request requests.OwnerRoadLineRequest) (int, error) {
	owners, err := su.settingRepo.GetOwnersRoadLineAll()
	if err != nil {
		return 0, responses.NewInternalServerError()
	}

	for _, owner := range owners {
		if request.Name == owner.Name {
			return 0, responses.NewDuplicatedNameError("Name:duplicate")
		}
	}

	owner := models.RefOwnerRoadLine{Name: request.Name, RefReflectivityRangeID: request.RefReflectivityRangeID, IsActive: true}

	if err := su.settingRepo.CreateOwnerRoadLine(&owner); err != nil {
		return 0, responses.NewInternalServerError()
	}

	return owner.ID, nil
}

func (su *settingUseCase) UpdateOwnerRoadLineByID(ownerId int, request requests.OwnerRoadLineRequest) error {
	_, err := su.settingRepo.GetOwnerRoadLineByID(ownerId)
	if err != nil {
		return responses.NewAppErr(400, err.Error())
	}

	owners, err := su.settingRepo.GetOwnersRoadLineAll()
	if err != nil {
		return responses.NewInternalServerError()
	}

	for _, owner := range owners {
		if ownerId != owner.ID && request.Name == owner.Name {
			return responses.NewDuplicatedNameError("Name:duplicate")
		}
	}

	owner := models.RefOwnerRoadLine{Name: request.Name, RefReflectivityRangeID: request.RefReflectivityRangeID}
	err = su.settingRepo.UpdateOwnerRoadLineByID(ownerId, owner)
	if err != nil {
		return responses.NewInternalServerError()
	}
	return nil
}

func (su *settingUseCase) DeleteOwnerRoadLineByID(id string) error {
	ownerId, err := helpers.ConvertStringToInt(id)
	if err != nil {
		return responses.NewInternalServerError()
	}

	err = su.settingRepo.DeleteOwnerRoadLineByID(ownerId)
	if err != nil {
		return responses.NewInternalServerError()
	}

	return nil
}

func (su *settingUseCase) GetOwners(params requests.QueryParams) (interface{}, error) {
	total, err := su.settingRepo.CountOwners(params)
	if err != nil {
		return responses.Pagination{}, responses.NewInternalServerError()
	}

	limit, offset, page := helpers.GetlimitOffsetPage(params.Limit, params.Page, total)

	owners, err := su.settingRepo.GetOwners(limit, offset, params)
	if err != nil {
		log.Println(err)
		return responses.Pagination{}, responses.NewInternalServerError()
	}
	// RefConditionRange := []models.RefConditionRange{}
	// err = su.settingRepo.GetDataList(&RefConditionRange, "")
	// if err != nil {
	// 	log.Println(err)
	// 	return responses.Pagination{}, responses.NewInternalServerError()
	// }
	pagination := helpers.Pagination(owners, limit, page, total)

	return pagination, nil
}

func (su *settingUseCase) GetOwnerByID(ownerId int) (models.RefOwner, error) {
	owner, err := su.settingRepo.GetOwnerByID(ownerId)
	if err != nil {
		if err.Error() == "record not found" {
			return models.RefOwner{}, responses.NewNotFoundError()
		}
		return models.RefOwner{}, responses.NewInternalServerError()
	}

	return owner, nil
}

func (su *settingUseCase) CreateOwner(request requests.OwnerRequest) (int, error) {
	owners, err := su.settingRepo.GetOwnersAll()
	if err != nil {
		return 0, responses.NewInternalServerError()
	}

	for _, owner := range owners {
		if request.Name == owner.Name {
			return 0, responses.NewDuplicatedNameError("Name:duplicate")
		}
	}

	owner := models.RefOwner{Name: request.Name, RefConditionRangeID: request.RefConditionRangeID, IsActive: true}

	if err := su.settingRepo.CreateOwner(&owner); err != nil {
		return 0, responses.NewInternalServerError()
	}

	return owner.ID, nil
}

func (su *settingUseCase) UpdateOwnerByID(ownerId int, request requests.OwnerRequest) error {
	_, err := su.settingRepo.GetOwnerByID(ownerId)
	if err != nil {
		return responses.NewAppErr(400, err.Error())
	}

	owners, err := su.settingRepo.GetOwnersAll()
	if err != nil {
		return responses.NewInternalServerError()
	}

	for _, owner := range owners {
		if ownerId != owner.ID && request.Name == owner.Name {
			return responses.NewDuplicatedNameError("Name:duplicate")
		}
	}

	owner := models.RefOwner{Name: request.Name, RefConditionRangeID: request.RefConditionRangeID}
	err = su.settingRepo.UpdateOwnerByID(ownerId, owner)
	if err != nil {
		return responses.NewInternalServerError()
	}
	return nil
}

func (su *settingUseCase) DeleteOwnerByID(id string) error {
	ownerId, err := helpers.ConvertStringToInt(id)
	if err != nil {
		return responses.NewInternalServerError()
	}

	err = su.settingRepo.DeleteOwnerByID(ownerId)
	if err != nil {
		return responses.NewInternalServerError()
	}

	return nil
}

func (su *settingUseCase) GetConditionList(ownerId int) ([]models.ParamsConditionPreload, error) {
	paramsCondition, err := su.settingRepo.GetParamsCondition(ownerId)
	if err != nil {
		return nil, responses.NewInternalServerError()
	}

	return paramsCondition, nil
}

func (su *settingUseCase) GetRoadLineList(ownerId int) (responses.RoadLineList, error) {
	paramsCondition, err := su.settingRepo.GetParamsRoadLine(ownerId)
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

func (su *settingUseCase) GetOwnerID(params url.Values) int {
	if params.Get("road_id") != "" {
		roadId, _ := helpers.ConvertStringToInt(params.Get("road_id"))
		fmt.Println(roadId)
		roadOwner, err := su.settingRepo.GetOwnerByRoadID(roadId)
		if err != nil {
			return 0
		}

		return roadOwner.RefOwnerID
	}

	ownerId, _ := helpers.ConvertStringToInt(params.Get("owner_id"))

	return ownerId
}

func (su *settingUseCase) CreateConditionList(ownerId int, requests requests.OwnerRequest) (interface{}, error) {
	grades, err := su.settingRepo.GetGrade()
	if err != nil {
		if err.Error() == "record not found" {
			return "", responses.NewNotFoundError()
		}
		return "", responses.NewInternalServerError()
	}

	gradeItems := make(map[int]string)
	for _, item := range grades {
		gradeItems[item.ID] = item.Name
	}

	var paramsConditions []models.ParamsCondition
	for _, v := range requests.ConditionList {
		var paramsCondition models.ParamsCondition
		// if v.LeftValue > v.RightValue {
		// 	if v.GradeID != 5 {
		// 		return "", responses.NewAppErr(400, "ประเภท "+v.ConditionType+" เกณฑ์ "+gradeItems[v.GradeID]+" : Min ต้อง น้อยกว่าหรือเท่ากับ Max")
		// 	}
		// }
		// paramsCondition.RefConditionRangeID = requests.RefConditionRangeID
		paramsCondition.RefOwnerID = ownerId
		paramsCondition.RefGradeID = v.GradeID
		paramsCondition.LeftValueAC = v.LeftValueAC
		paramsCondition.LeftConditionAC = "<="
		paramsCondition.RightValueAC = v.RightValueAC
		paramsCondition.RightConditionAC = "<"
		paramsCondition.LeftValueCC = v.LeftValueCC
		paramsCondition.LeftConditionCC = "<="
		paramsCondition.RightValueCC = v.RightValueCC
		paramsCondition.RightConditionCC = "<"
		paramsCondition.ConditionType = strings.ToUpper(v.ConditionType)
		paramsConditions = append(paramsConditions, paramsCondition)
	}
	if err := su.settingRepo.CreateData(paramsConditions); err != nil {
		return "", responses.NewInternalServerError()
	}
	return paramsConditions, nil
}

func (su *settingUseCase) UpdateConditionList(ownerId int, requests requests.OwnerRequest) error {
	if err := su.settingRepo.DeleteCondition(ownerId); err != nil {
		return err
	}

	grades, err := su.settingRepo.GetGrade()
	if err != nil {
		if err.Error() == "record not found" {
			return responses.NewNotFoundError()
		}
		return responses.NewInternalServerError()
	}

	gradeItems := make(map[int]string)
	for _, item := range grades {
		gradeItems[item.ID] = item.Name
	}

	var paramsConditions []models.ParamsCondition
	for _, v := range requests.ConditionList {
		var paramsCondition models.ParamsCondition
		// if v.LeftValue > v.RightValue {
		// 	if v.GradeID != 5 {
		// 		return responses.NewAppErr(400, "ประเภท "+v.ConditionType+" เกณฑ์ "+gradeItems[v.GradeID]+" : Min ต้อง น้อยกว่าหรือเท่ากับ Max")
		// 	}
		// }
		// paramsCondition.RefConditionRangeID = requests.RefConditionRangeID
		paramsCondition.RefOwnerID = ownerId
		paramsCondition.RefGradeID = v.GradeID
		paramsCondition.LeftValueAC = v.LeftValueAC
		paramsCondition.LeftConditionAC = "<="
		paramsCondition.RightValueAC = v.RightValueAC
		paramsCondition.RightConditionAC = "<"
		paramsCondition.LeftValueCC = v.LeftValueCC
		paramsCondition.LeftConditionCC = "<="
		paramsCondition.RightValueCC = v.RightValueCC
		paramsCondition.RightConditionCC = "<"
		paramsCondition.ConditionType = strings.ToUpper(v.ConditionType)
		paramsConditions = append(paramsConditions, paramsCondition)

	}
	if err := su.settingRepo.CreateData(paramsConditions); err != nil {
		return responses.NewInternalServerError()
	}
	return nil
}

func (su *settingUseCase) CreateRoadLineList(ownerId int, requests requests.OwnerRoadLineRequest) (interface{}, error) {
	grades, err := su.settingRepo.GetGrade()
	if err != nil {
		if err.Error() == "record not found" {
			return "", responses.NewNotFoundError()
		}
		return "", responses.NewInternalServerError()
	}

	gradeItems := make(map[int]string)
	for _, item := range grades {
		gradeItems[item.ID] = item.Name
	}

	var paramsConditions []models.ParamsRoadLine
	for _, v := range requests.ConditionList {
		var paramsCondition models.ParamsRoadLine
		// if v.LeftValue > v.RightValue {
		// 	if v.GradeID != 5 {
		// 		return "", responses.NewAppErr(400, "ประเภท "+v.ConditionType+" เกณฑ์ "+gradeItems[v.GradeID]+" : Min ต้อง น้อยกว่าหรือเท่ากับ Max")
		// 	}
		// }
		// paramsCondition.RefConditionRangeID = requests.RefConditionRangeID
		paramsCondition.RefOwnerRoadLineID = ownerId
		paramsCondition.RefGradeID = v.GradeID
		paramsCondition.LeftValueWhite = v.LeftValueWhite
		paramsCondition.LeftConditionWhite = "<="
		paramsCondition.RightValueWhite = v.RightValueWhite
		paramsCondition.RightConditionWhite = "<"
		paramsCondition.LeftValueYellow = v.LeftValueYellow
		paramsCondition.LeftConditionYellow = "<="
		paramsCondition.RightValueYellow = v.RightValueYellow
		paramsCondition.RightConditionYellow = "<"
		// paramsCondition.ConditionType = strings.ToUpper(v.ConditionType)
		paramsConditions = append(paramsConditions, paramsCondition)
	}
	if err := su.settingRepo.CreateData(paramsConditions); err != nil {
		return "", responses.NewInternalServerError()
	}
	return paramsConditions, nil
}

func (su *settingUseCase) UpdateRoadLineList(ownerId int, requests requests.OwnerRoadLineRequest) error {
	if err := su.settingRepo.DeleteRoadLine(ownerId); err != nil {
		return err
	}

	grades, err := su.settingRepo.GetGrade()
	if err != nil {
		if err.Error() == "record not found" {
			return responses.NewNotFoundError()
		}
		return responses.NewInternalServerError()
	}

	gradeItems := make(map[int]string)
	for _, item := range grades {
		gradeItems[item.ID] = item.Name
	}

	var paramsConditions []models.ParamsRoadLine
	for _, v := range requests.ConditionList {
		var paramsCondition models.ParamsRoadLine
		// if v.LeftValue > v.RightValue {
		// 	if v.GradeID != 5 {
		// 		return responses.NewAppErr(400, "ประเภท "+v.ConditionType+" เกณฑ์ "+gradeItems[v.GradeID]+" : Min ต้อง น้อยกว่าหรือเท่ากับ Max")
		// 	}
		// }
		// paramsCondition.RefConditionRangeID = requests.RefConditionRangeID
		paramsCondition.RefOwnerRoadLineID = ownerId
		paramsCondition.RefGradeID = v.GradeID
		paramsCondition.LeftValueWhite = v.LeftValueWhite
		paramsCondition.LeftConditionWhite = "<="
		paramsCondition.RightValueWhite = v.RightValueWhite
		paramsCondition.RightConditionWhite = "<"
		paramsCondition.LeftValueYellow = v.LeftValueYellow
		paramsCondition.LeftConditionYellow = "<="
		paramsCondition.RightValueYellow = v.RightValueYellow
		paramsCondition.RightConditionYellow = "<"
		// paramsCondition.ConditionType = strings.ToUpper(v.ConditionType)
		paramsConditions = append(paramsConditions, paramsCondition)

	}
	if err := su.settingRepo.CreateData(paramsConditions); err != nil {
		return responses.NewInternalServerError()
	}
	return nil
}

func (su *settingUseCase) GetSigns(c *gin.Context, params requests.QueryParams) (interface{}, error) {
	total, err := su.settingRepo.CountSigns(params.Name)
	if err != nil {
		return responses.Pagination{}, responses.NewInternalServerError()
	}

	limit, offset, page := helpers.GetlimitOffsetPage(params.Limit, params.Page, total)

	signs, err := su.settingRepo.GetSigns(limit, offset, params.Name)
	if err != nil {
		log.Println(err)
		return responses.Pagination{}, responses.NewInternalServerError()
	}

	for i := 0; i < len(signs); i++ {
		signs[i].SignImageFilepath = os.Getenv("STORAGE_IP") + "/" + signs[i].SignImageFilepath
	}

	pagination := helpers.Pagination(signs, limit, page, total)
	return pagination, nil
}

func (su *settingUseCase) GetSignByID(c *gin.Context, id string) (models.RefAssetSignImage, error) {
	signId, err := helpers.ConvertStringToInt(id)
	if err != nil {
		return models.RefAssetSignImage{}, err
	}

	sign := models.RefAssetSignImage{}
	err = su.settingRepo.GetByID(&sign, signId)
	if err != nil {
		log.Println(err)
		if err.Error() == "record not found" {
			return models.RefAssetSignImage{}, responses.NewNotFoundError()
		}

		return models.RefAssetSignImage{}, responses.NewInternalServerError()
	}

	sign.SignImageFilepath = os.Getenv("STORAGE_IP") + "/" + sign.SignImageFilepath
	return sign, nil
}

func (su *settingUseCase) DeleteSignByID(id string) error {
	signId, err := helpers.ConvertStringToInt(id)
	if err != nil {
		return err
	}

	err = su.settingRepo.DeleteByID("ref_asset_sign_image", signId)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}
	return nil
}

func (su *settingUseCase) CreateSign(c *gin.Context, request requests.SignImageRequest) error {
	if !helpers.IsImageTypeJPEGOrPNG(request.Image.Filename) {
		return responses.NewImageTypeError()
	}

	if helpers.IsFileSizeGreaterThanLimit(request.Image.Size, 5) {
		return responses.NewFileSizeExceedLimitError()
	}

	dstPath, err := SaveFile(c, request.Image, os.Getenv("SIGN_DIR"))
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}

	sign := models.RefAssetSignImage{
		Name:              request.Name,
		Abbr:              request.Abbr,
		SignImageFilepath: dstPath,
		StatusCode:        "A",
		Status:            1,
	}

	err = su.settingRepo.CreateSignImage(&sign)
	if err != nil {
		//roll back state when internal server error
		if err := os.Remove(dstPath); err != nil {
			log.Println(err)
		}
		return responses.NewInternalServerError()
	}

	return nil
}

func SaveFile(c *gin.Context, file *multipart.FileHeader, dir string) (string, error) {
	//add time to make image file unique
	dstPath := dir + time.Now().Format("2006-01-02T15:04:05.000") + "_" + file.Filename
	err := c.SaveUploadedFile(file, dstPath)
	if err != nil {
		return "", err
	}

	return dstPath, nil
}

func (su *settingUseCase) UpdateSignByID(c *gin.Context, request requests.SignImageRequest) error {
	signId, err := helpers.ConvertStringToInt(c.Param("id"))
	if err != nil {
		return err
	}

	sign := models.RefAssetSignImage{}

	var dstPath string
	var oldDstPath string
	if request.Image != nil {
		if !helpers.IsImageTypeJPEGOrPNG(request.Image.Filename) {
			return responses.NewImageTypeError()
		}

		if helpers.IsFileSizeGreaterThanLimit(request.Image.Size, 5) {
			return responses.NewFileSizeExceedLimitError()
		}

		err := su.settingRepo.GetByID(&sign, signId)
		if err != nil {
			return responses.NewInternalServerError()
		}

		oldDstPath = sign.SignImageFilepath

		// check whether image file is exist or not if it exist then remove
		if request.ImageStatus != "not_edit" {
			if err := RemoveFileIfExist(sign.SignImageFilepath); err != nil {
				log.Println(err)
				return responses.NewInternalServerError()
			}
		}
		if request.ImageStatus == "not_edit" {
			sign.SignImageFilepath = oldDstPath
		} else {
			storagesFolder := os.Getenv("SIGN_DIR")
			if _, err := os.Stat(storagesFolder); os.IsNotExist(err) {
				err := os.MkdirAll(storagesFolder, 0755)
				if err != nil {
					panic(err)
				}
			}
			dstPath, err = SaveFile(c, request.Image, os.Getenv("SIGN_DIR"))
			if err != nil {
				log.Println(err)
				return responses.NewInternalServerError()
			}
			sign.SignImageFilepath = dstPath
		}

	}

	sign.Name = request.Name
	sign.Abbr = request.Abbr
	err = su.settingRepo.UpdateSignImage(signId, sign)
	if err != nil {
		//roll back state when internal server error
		if err := os.Remove(dstPath); err != nil {
			log.Println(err)
		}
		_, err := os.Create(oldDstPath)
		if err != nil {
			log.Println(err)
		}
		log.Println(err)
		return responses.NewInternalServerError()
	}

	return nil
}

func RemoveFileIfExist(dstPath string) error {
	file, _ := os.Stat(dstPath)
	if file != nil {
		err := os.Remove(dstPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (su *settingUseCase) DeleteInterventionCriteriaById(id string, c *gin.Context) error {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	idInteger, err := helpers.ConvertStringToInt(id)
	if err != nil {
		return err
	}

	getInterventionCriteria, err := su.settingRepo.GetInterventionCriteriaById(idInteger)
	if err != nil {
		log.Println(err)
		return err
	}
	getInterventionCriteria.IsDeleted = true
	getInterventionCriteria.UpdatedAt = dateTime
	getInterventionCriteria.UpdatedBy = userID

	err = su.settingRepo.UpdateInterventionCriteria(&getInterventionCriteria)
	if err != nil {
		log.Println(err)
		return err
	}

	interventionCriteriaConditionList, _ := su.settingRepo.GetInterventionCriteriaConditionListByInterventionCriteriaId(idInteger)

	for range interventionCriteriaConditionList {
		var interventionCriteriaCondition models.InterventionCriteriaCondition
		interventionCriteriaCondition.IsDeleted = true
		interventionCriteriaCondition.UpdatedAt = dateTime
		interventionCriteriaCondition.UpdatedBy = userID
		err := su.settingRepo.UpdateInterventionCriteriaCondition(&interventionCriteriaCondition)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	var asphalt []int
	var concrete []int
	getInterventionCriteriaSequence, err := su.settingRepo.GetInterventionCriteria()
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}
	for _, item := range getInterventionCriteriaSequence {
		criteriaMethod, err := su.settingRepo.GetCriteriaMethodById(item.MaintenanceMethod)
		if err != nil {
			log.Println(err)
			return responses.NewInternalServerError()
		}

		var InterventionCriteriaSequence responses.InterventionCriteriaSequence
		InterventionCriteriaSequence.Id = item.Id
		InterventionCriteriaSequence.Name = criteriaMethod.Name + " | " + item.MaintenanceStandardName

		if criteriaMethod.Surface == "as" {
			asphalt = append(asphalt, item.Id)
		} else if criteriaMethod.Surface == "cc" {
			concrete = append(concrete, item.Id)
		}
	}

	for index, item := range asphalt {
		getInterventionCriteria, err := su.settingRepo.GetInterventionCriteriaById(item)
		if err != nil {
			log.Println(err)
			return responses.NewInternalServerError()
		}
		getInterventionCriteria.MaintenanceSequence = index + 1
		getInterventionCriteria.UpdatedAt = dateTime
		getInterventionCriteria.UpdatedBy = userID
		err = su.settingRepo.UpdateInterventionCriteria(&getInterventionCriteria)
		if err != nil {
			log.Println(err)
			return responses.NewInternalServerError()
		}

	}

	for index, item := range concrete {
		getInterventionCriteria, err := su.settingRepo.GetInterventionCriteriaById(item)
		if err != nil {
			log.Println(err)
			return responses.NewInternalServerError()
		}
		getInterventionCriteria.MaintenanceSequence = index + 1
		getInterventionCriteria.UpdatedAt = dateTime
		getInterventionCriteria.UpdatedBy = userID
		err = su.settingRepo.UpdateInterventionCriteria(&getInterventionCriteria)
		if err != nil {
			log.Println(err)
			return responses.NewInternalServerError()
		}

	}

	err = su.MergeInterventionCriteria(userID)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}

	return nil
}

func (su *settingUseCase) CreateInterventionCriteria(params requests.InterventionCriteria, c *gin.Context) (requests.InterventionCriteria, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	if len(params.Asphalt) > 0 {

		for _, asphalt := range params.Asphalt {

			if asphalt.ID == 0 {
				continue
			}

			for _, interventionCriteriaData := range asphalt.InterventionCriteriaDatas {

				//criteriaMethod := 2
				// criteriaMethod, err := su.settingRepo.GetCriteriaMethodByName(item.MaintenanceMethod)
				// if err != nil {
				// 	log.Println(err)
				// 	return requests.InterventionCriteria{}, responses.NewInternalServerError()
				// }

				latestSurfaceParams, err := su.settingRepo.GetLatestSurfaceParamsById(interventionCriteriaData.MaintenanceSurfaceTypeId)
				if err != nil && err.Error() != "record not found" {
					log.Println(err)
					return requests.InterventionCriteria{}, responses.NewInternalServerError()
				}

				maintenanceSurfaceTypeIdParams := 0
				if (latestSurfaceParams == models.RefSurfaceParam{}) {
					maintenanceSurfaceTypeIdParams = 0
				} else {
					maintenanceSurfaceTypeIdParams = latestSurfaceParams.ID
				}

				getInterventionCriteriaByNotId, err := su.settingRepo.GetInterventionCriteriaByNotId(interventionCriteriaData.Id)
				if err != nil {
					log.Println(err)
				}

				getInterventionCriteriaById, err := su.settingRepo.GetInterventionCriteriaById(interventionCriteriaData.Id)
				if err != nil {
					log.Println(err)
				}

				count, err := su.settingRepo.CountInterventionCriteriaByMaintenanceMethod("as")
				if err != nil {
					log.Println(err)
				}

				for _, item2 := range getInterventionCriteriaByNotId {
					if interventionCriteriaData.MaintenanceStandardName == item2.MaintenanceStandardName {
						return requests.InterventionCriteria{}, responses.NewDuplicatedNameError("standard_name")
					}
				}

				var interventionCriteria models.InterventionCriteria

				interventionCriteria.Id = interventionCriteriaData.Id
				if interventionCriteriaData.IsNew {
					interventionCriteria.Id = 0
				}

				interventionCriteria.MaintenanceMethod = asphalt.ID
				interventionCriteria.MaintenanceCostPerUnit = interventionCriteriaData.MaintenanceCostPerUnit
				interventionCriteria.MaintenanceDescription = interventionCriteriaData.MaintenanceDescription
				interventionCriteria.MaintenanceScraping = interventionCriteriaData.MaintenanceScraping
				if (getInterventionCriteriaById == models.InterventionCriteria{}) {
					interventionCriteria.MaintenanceSequence = count.Count + 1
				} else {
					interventionCriteria.MaintenanceSequence = getInterventionCriteriaById.MaintenanceSequence
				}
				interventionCriteria.MaintenanceStandardName = interventionCriteriaData.MaintenanceStandardName
				interventionCriteria.MaintenanceSurfaceTypeId = interventionCriteriaData.MaintenanceSurfaceTypeId
				interventionCriteria.MaintenanceSurfaceTypeIdParams = maintenanceSurfaceTypeIdParams
				interventionCriteria.MaintenanceThickness = interventionCriteriaData.MaintenanceThickness
				interventionCriteria.IsDeleted = false
				interventionCriteria.IsShow = true
				interventionCriteria.CreatedBy = userID
				interventionCriteria.CreatedAt = dateTime
				interventionCriteria.UpdatedBy = userID
				interventionCriteria.UpdatedAt = dateTime
				err = su.settingRepo.UpdateInterventionCriteria(&interventionCriteria)
				if err != nil {
					log.Println(err)
					return requests.InterventionCriteria{}, responses.NewInternalServerError()
				}

				if interventionCriteriaData.Id != 0 {
					err = su.settingRepo.DeleteInterventionCriteriaConditionById(interventionCriteriaData.Id)
					if err != nil {
						log.Println(err)
						return requests.InterventionCriteria{}, responses.NewInternalServerError()
					}
				}

				for index, item := range interventionCriteriaData.MaintenanceCondition {
					var interventionCriteriaCondition models.InterventionCriteriaCondition

					interventionCriteriaCondition.Id = item.Id
					if item.IsNew {
						interventionCriteriaCondition.Id = 0
					}

					interventionCriteriaCondition.ConditionSequence = index + 1
					interventionCriteriaCondition.InterventionCriteriaId = interventionCriteria.Id
					interventionCriteriaCondition.ConditionCriterion = item.ConditionCriterion
					interventionCriteriaCondition.ConditionLink = item.ConditionLink
					interventionCriteriaCondition.ConditionOperation_1 = item.ConditionOperation1
					interventionCriteriaCondition.ConditionOperation_2 = item.ConditionOperation2
					interventionCriteriaCondition.ConditionValue_1 = item.ConditionValue1
					interventionCriteriaCondition.ConditionValue_2 = item.ConditionValue2
					interventionCriteriaCondition.IsDeleted = false
					interventionCriteriaCondition.UpdatedBy = userID
					interventionCriteriaCondition.UpdatedAt = dateTime
					interventionCriteriaCondition.CreatedBy = userID
					interventionCriteriaCondition.CreatedAt = dateTime
					err = su.settingRepo.UpdateInterventionCriteriaCondition(&interventionCriteriaCondition)
					if err != nil {
						log.Println(err)
						return requests.InterventionCriteria{}, responses.NewInternalServerError()
					}
				}

			}

		}

	}

	if len(params.Concrete) > 0 {

		for _, concrete := range params.Concrete {
			if concrete.ID == 0 {
				continue
			}

			for _, interventionCriteriaData := range concrete.InterventionCriteriaDatas {

				//criteriaMethod := 2
				// criteriaMethod, err := su.settingRepo.GetCriteriaMethodByName(item.MaintenanceMethod)
				// if err != nil {
				// 	log.Println(err)
				// 	return requests.InterventionCriteria{}, responses.NewInternalServerError()
				// }

				latestSurfaceParams, err := su.settingRepo.GetLatestSurfaceParamsById(interventionCriteriaData.MaintenanceSurfaceTypeId)
				if err != nil && err.Error() != "record not found" {
					log.Println(err)
					return requests.InterventionCriteria{}, responses.NewInternalServerError()
				}

				maintenanceSurfaceTypeIdParams := 0
				if (latestSurfaceParams == models.RefSurfaceParam{}) {
					maintenanceSurfaceTypeIdParams = 0
				} else {
					maintenanceSurfaceTypeIdParams = latestSurfaceParams.ID
				}

				getInterventionCriteriaByNotId, err := su.settingRepo.GetInterventionCriteriaByNotId(interventionCriteriaData.Id)
				if err != nil {
					log.Println(err)
				}

				getInterventionCriteriaById, err := su.settingRepo.GetInterventionCriteriaById(interventionCriteriaData.Id)
				if err != nil {
					log.Println(err)
				}

				count, err := su.settingRepo.CountInterventionCriteriaByMaintenanceMethod("cc")
				if err != nil {
					log.Println(err)
				}

				for _, item2 := range getInterventionCriteriaByNotId {
					if interventionCriteriaData.MaintenanceStandardName == item2.MaintenanceStandardName {
						return requests.InterventionCriteria{}, responses.NewDuplicatedNameError("standard_name")
					}
				}

				var interventionCriteria models.InterventionCriteria
				interventionCriteria.Id = interventionCriteriaData.Id
				interventionCriteria.MaintenanceMethod = concrete.ID
				interventionCriteria.MaintenanceCostPerUnit = interventionCriteriaData.MaintenanceCostPerUnit
				interventionCriteria.MaintenanceDescription = interventionCriteriaData.MaintenanceDescription
				interventionCriteria.MaintenanceScraping = interventionCriteriaData.MaintenanceScraping
				if (getInterventionCriteriaById == models.InterventionCriteria{}) {
					interventionCriteria.MaintenanceSequence = count.Count + 1
				} else {
					interventionCriteria.MaintenanceSequence = getInterventionCriteriaById.MaintenanceSequence
				}
				interventionCriteria.MaintenanceStandardName = interventionCriteriaData.MaintenanceStandardName
				interventionCriteria.MaintenanceSurfaceTypeId = interventionCriteriaData.MaintenanceSurfaceTypeId
				interventionCriteria.MaintenanceSurfaceTypeIdParams = maintenanceSurfaceTypeIdParams
				interventionCriteria.MaintenanceThickness = interventionCriteriaData.MaintenanceThickness
				interventionCriteria.IsDeleted = false
				interventionCriteria.IsShow = true
				interventionCriteria.CreatedBy = userID
				interventionCriteria.CreatedAt = dateTime
				interventionCriteria.UpdatedBy = userID
				interventionCriteria.UpdatedAt = dateTime
				err = su.settingRepo.UpdateInterventionCriteria(&interventionCriteria)
				if err != nil {
					log.Println(err)
					return requests.InterventionCriteria{}, responses.NewInternalServerError()
				}

				if interventionCriteriaData.Id != 0 {
					err = su.settingRepo.DeleteInterventionCriteriaConditionById(interventionCriteriaData.Id)
					if err != nil {
						log.Println(err)
						return requests.InterventionCriteria{}, responses.NewInternalServerError()
					}
				}

				for index, item := range interventionCriteriaData.MaintenanceCondition {
					var interventionCriteriaCondition models.InterventionCriteriaCondition

					interventionCriteriaCondition.Id = item.Id
					if item.IsNew {
						interventionCriteriaCondition.Id = 0
					}

					interventionCriteriaCondition.ConditionSequence = index + 1
					interventionCriteriaCondition.InterventionCriteriaId = interventionCriteria.Id
					interventionCriteriaCondition.ConditionCriterion = item.ConditionCriterion
					interventionCriteriaCondition.ConditionLink = item.ConditionLink
					interventionCriteriaCondition.ConditionOperation_1 = item.ConditionOperation1
					interventionCriteriaCondition.ConditionOperation_2 = item.ConditionOperation2
					interventionCriteriaCondition.ConditionValue_1 = item.ConditionValue1
					interventionCriteriaCondition.ConditionValue_2 = item.ConditionValue2
					interventionCriteriaCondition.IsDeleted = false
					interventionCriteriaCondition.UpdatedBy = userID
					interventionCriteriaCondition.UpdatedAt = dateTime
					interventionCriteriaCondition.CreatedBy = userID
					interventionCriteriaCondition.CreatedAt = dateTime
					err = su.settingRepo.UpdateInterventionCriteriaCondition(&interventionCriteriaCondition)
					if err != nil {
						log.Println(err)
						return requests.InterventionCriteria{}, responses.NewInternalServerError()
					}
				}

			}

		}

	}

	err := su.MergeInterventionCriteria(userID)
	if err != nil {
		log.Println(err)
		return requests.InterventionCriteria{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) MergeInterventionCriteria(userID int) error {
	dateTime := time.Now()
	getInterventionCriteria, err := su.settingRepo.GetInterventionCriteria()
	if err != nil {
		log.Println(err)
		return err
	}

	err = su.settingRepo.UpdateInterventionCriteriaParamsByIsLatestIsFalse()
	if err != nil {
		log.Println(err)
		return err
	}

	var interventionCriteriaParams responses.InterventionCriteriaParams

	for _, item := range getInterventionCriteria {

		var interventionCriteriaSurfaceParam responses.InterventionCriteriaSurfaceParams

		criteriaMethod, err := su.settingRepo.GetCriteriaMethodById(item.MaintenanceMethod)
		if err != nil {
			log.Println(err)
			return err
		}

		var interventionCriteriaMaintenance responses.InterventionCriteriaMaintenance
		interventionCriteriaMaintenance.Id = item.Id
		interventionCriteriaMaintenance.MaintenanceCostPerUnit = item.MaintenanceCostPerUnit
		interventionCriteriaMaintenance.MaintenanceDescription = item.MaintenanceDescription
		interventionCriteriaMaintenance.MaintenanceScraping = item.MaintenanceScraping
		interventionCriteriaMaintenance.MaintenanceSequence = item.MaintenanceSequence
		interventionCriteriaMaintenance.MaintenanceStandardName = item.MaintenanceStandardName
		interventionCriteriaMaintenance.MaintenanceSurfaceTypeId = item.MaintenanceSurfaceTypeId
		interventionCriteriaMaintenance.MaintenanceThickness = item.MaintenanceThickness

		var interventionCriteriaCinditionList []responses.InterventionCriteriaCindition
		interventionCriteriaConditionList, _ := su.settingRepo.GetInterventionCriteriaConditionListByInterventionCriteriaId(item.Id)
		for _, itemCondition := range interventionCriteriaConditionList {
			var interventionCriteriaCindition responses.InterventionCriteriaCindition
			interventionCriteriaCindition.Id = itemCondition.Id
			interventionCriteriaCindition.ConditionCriterion = itemCondition.ConditionCriterion
			interventionCriteriaCindition.ConditionLink = itemCondition.ConditionLink
			interventionCriteriaCindition.ConditionOperation1 = itemCondition.ConditionOperation_1
			interventionCriteriaCindition.ConditionOperation2 = itemCondition.ConditionOperation_2
			interventionCriteriaCindition.ConditionValue1 = itemCondition.ConditionValue_1
			interventionCriteriaCindition.ConditionValue2 = itemCondition.ConditionValue_2
			interventionCriteriaCinditionList = append(interventionCriteriaCinditionList, interventionCriteriaCindition)
		}
		interventionCriteriaMaintenance.MaintenanceCondition = interventionCriteriaCinditionList

		interventionCriteriaSurfaceParam.ID = item.Id
		interventionCriteriaSurfaceParam.Name = item.MaintenanceStandardName
		interventionCriteriaSurfaceParam.InterventionCriteriaMaintenances = append(interventionCriteriaSurfaceParam.InterventionCriteriaMaintenances, interventionCriteriaMaintenance)

		interventionCriteriaSurfaceParam.ID = criteriaMethod.ID
		interventionCriteriaSurfaceParam.Name = criteriaMethod.Name

		switch criteriaMethod.Surface {
		case "as":
			interventionCriteriaParams.Asphalt = append(interventionCriteriaParams.Asphalt, interventionCriteriaSurfaceParam)
		case "cc":
			interventionCriteriaParams.Concrete = append(interventionCriteriaParams.Concrete, interventionCriteriaSurfaceParam)

		}

	}

	structByte, err := json.Marshal(interventionCriteriaParams)
	if err != nil {
		log.Println(err)
		return err
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, structByte); err != nil {
		log.Println(err)
		return err
	}

	var settingInterventionCriteria models.SettingInterventionCriteria
	settingInterventionCriteria.Params = strings.ReplaceAll(strings.ReplaceAll(buffer.String(), `\u0026`, "&"), `\u003c`, "<")
	settingInterventionCriteria.IsLatest = true
	settingInterventionCriteria.CreatedAt = dateTime
	settingInterventionCriteria.CreatedBy = userID
	err = su.settingRepo.CreateInterventionCriteriaParams(&settingInterventionCriteria)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (su *settingUseCase) GetInterventionCriteria(c *gin.Context) (responses.InterventionCriteria, error) {

	getInterventionCriteria, err := su.settingRepo.GetInterventionCriteria()
	if err != nil {
		log.Println(err)
		return responses.InterventionCriteria{}, responses.NewInternalServerError()
	}

	var responsesInterventionCriteria responses.InterventionCriteria

	asphaltMap := make(map[int]*responses.InterventionCriteriaSurface)
	concreteMap := make(map[int]*responses.InterventionCriteriaSurface)

	for _, item := range getInterventionCriteria {
		criteriaMethod, err := su.settingRepo.GetCriteriaMethodById(item.MaintenanceMethod)
		if err != nil {
			log.Println(err)
			return responses.InterventionCriteria{}, responses.NewInternalServerError()
		}

		if criteriaMethod.ID == 0 {
			continue
		}

		var interventionCriteriaMaintenance responses.InterventionCriteriaMaintenance
		interventionCriteriaMaintenance.Id = item.Id
		interventionCriteriaMaintenance.MaintenanceCostPerUnit = item.MaintenanceCostPerUnit
		interventionCriteriaMaintenance.MaintenanceDescription = item.MaintenanceDescription
		interventionCriteriaMaintenance.MaintenanceScraping = item.MaintenanceScraping
		interventionCriteriaMaintenance.MaintenanceSequence = item.MaintenanceSequence
		interventionCriteriaMaintenance.MaintenanceStandardName = item.MaintenanceStandardName
		interventionCriteriaMaintenance.MaintenanceSurfaceTypeId = item.MaintenanceSurfaceTypeId
		interventionCriteriaMaintenance.MaintenanceThickness = item.MaintenanceThickness

		var interventionCriteriaConditionList []responses.InterventionCriteriaCindition
		interventionCriteriaConditionList = []responses.InterventionCriteriaCindition{}

		interventionCriteriaConditions, _ := su.settingRepo.GetInterventionCriteriaConditionListByInterventionCriteriaId(item.Id)
		for _, itemCondition := range interventionCriteriaConditions {
			var interventionCriteriaCondition responses.InterventionCriteriaCindition
			interventionCriteriaCondition.Id = itemCondition.Id
			interventionCriteriaCondition.ConditionCriterion = itemCondition.ConditionCriterion
			interventionCriteriaCondition.ConditionLink = itemCondition.ConditionLink
			interventionCriteriaCondition.ConditionOperation1 = itemCondition.ConditionOperation_1
			interventionCriteriaCondition.ConditionOperation2 = itemCondition.ConditionOperation_2
			interventionCriteriaCondition.ConditionValue1 = itemCondition.ConditionValue_1
			interventionCriteriaCondition.ConditionValue2 = itemCondition.ConditionValue_2
			interventionCriteriaConditionList = append(interventionCriteriaConditionList, interventionCriteriaCondition)
		}

		interventionCriteriaMaintenance.MaintenanceCondition = interventionCriteriaConditionList

		var interventionCriteriaSurface *responses.InterventionCriteriaSurface
		var surfaceMap map[int]*responses.InterventionCriteriaSurface

		switch criteriaMethod.Surface {
		case "as":
			surfaceMap = asphaltMap
		case "cc":
			surfaceMap = concreteMap
		}

		if existingSurface, ok := surfaceMap[criteriaMethod.ID]; ok {
			interventionCriteriaSurface = existingSurface
		} else {
			interventionCriteriaSurface = &responses.InterventionCriteriaSurface{
				ID:                               criteriaMethod.ID,
				Name:                             criteriaMethod.Name,
				InterventionCriteriaMaintenances: []responses.InterventionCriteriaMaintenance{},
			}
			surfaceMap[criteriaMethod.ID] = interventionCriteriaSurface
		}

		interventionCriteriaSurface.InterventionCriteriaMaintenances = append(interventionCriteriaSurface.InterventionCriteriaMaintenances, interventionCriteriaMaintenance)
	}

	for _, surface := range asphaltMap {
		responsesInterventionCriteria.Asphalt = append(responsesInterventionCriteria.Asphalt, *surface)
	}

	for _, surface := range concreteMap {
		responsesInterventionCriteria.Concrete = append(responsesInterventionCriteria.Concrete, *surface)
	}

	if len(responsesInterventionCriteria.Asphalt) == 0 {
		responsesInterventionCriteria.Asphalt = []responses.InterventionCriteriaSurface{}
	} else {
		sort.Slice(responsesInterventionCriteria.Asphalt, func(i, j int) bool {
			return responsesInterventionCriteria.Asphalt[i].ID < responsesInterventionCriteria.Asphalt[j].ID
		})
	}

	if len(responsesInterventionCriteria.Concrete) == 0 {
		responsesInterventionCriteria.Concrete = []responses.InterventionCriteriaSurface{}
	} else {
		sort.Slice(responsesInterventionCriteria.Concrete, func(i, j int) bool {
			return responsesInterventionCriteria.Concrete[i].ID < responsesInterventionCriteria.Concrete[j].ID
		})

	}

	return responsesInterventionCriteria, nil
}

func (su *settingUseCase) GetInterventionCriteriaMethod(c *gin.Context) (responses.InterventionCriteriaSequenceCriteriaMethod, error) {
	getRefCriteriaMethod, err := su.settingRepo.GetRefCriteriaMethod()
	if err != nil {
		log.Println(err)
		return responses.InterventionCriteriaSequenceCriteriaMethod{}, responses.NewInternalServerError()
	}

	var InterventionCriteriaSequenceCriteriaMethod responses.InterventionCriteriaSequenceCriteriaMethod
	for _, item := range getRefCriteriaMethod {

		var InterventionCriteriaSequence responses.InterventionCriteriaSequence
		InterventionCriteriaSequence.Id = item.ID
		InterventionCriteriaSequence.Name = item.Name

		if item.Surface == "as" {
			InterventionCriteriaSequenceCriteriaMethod.Asphalt = append(InterventionCriteriaSequenceCriteriaMethod.Asphalt, InterventionCriteriaSequence)
		} else if item.Surface == "cc" {
			InterventionCriteriaSequenceCriteriaMethod.Concrete = append(InterventionCriteriaSequenceCriteriaMethod.Concrete, InterventionCriteriaSequence)
		}
	}

	return InterventionCriteriaSequenceCriteriaMethod, nil
}

func (su *settingUseCase) GetInterventionCriteriaSequence(c *gin.Context) (responses.InterventionCriteriaSequenceCriteriaMethod, error) {
	getInterventionCriteria, err := su.settingRepo.GetInterventionCriteria()
	if err != nil {
		log.Println(err)
		return responses.InterventionCriteriaSequenceCriteriaMethod{}, responses.NewInternalServerError()
	}
	var InterventionCriteriaSequenceCriteriaMethod responses.InterventionCriteriaSequenceCriteriaMethod
	for _, item := range getInterventionCriteria {
		criteriaMethod, err := su.settingRepo.GetCriteriaMethodById(item.MaintenanceMethod)
		if err != nil {
			log.Println(err)
			return responses.InterventionCriteriaSequenceCriteriaMethod{}, responses.NewInternalServerError()
		}

		var InterventionCriteriaSequence responses.InterventionCriteriaSequence
		InterventionCriteriaSequence.Id = item.Id
		InterventionCriteriaSequence.Name = criteriaMethod.Name + " | " + item.MaintenanceStandardName

		if criteriaMethod.Surface == "as" {
			InterventionCriteriaSequenceCriteriaMethod.Asphalt = append(InterventionCriteriaSequenceCriteriaMethod.Asphalt, InterventionCriteriaSequence)
		} else if criteriaMethod.Surface == "cc" {
			InterventionCriteriaSequenceCriteriaMethod.Concrete = append(InterventionCriteriaSequenceCriteriaMethod.Concrete, InterventionCriteriaSequence)
		}
	}

	return InterventionCriteriaSequenceCriteriaMethod, nil
}

func (su *settingUseCase) CreateInterventionCriteriaSequence(params requests.InterventionCriteriaSequenceCriteriaMethod, c *gin.Context) (requests.InterventionCriteriaSequenceCriteriaMethod, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	for index, item := range params.Asphalt {
		getInterventionCriteria, err := su.settingRepo.GetInterventionCriteriaById(item)
		if err != nil {
			log.Println(err)
			return requests.InterventionCriteriaSequenceCriteriaMethod{}, responses.NewInternalServerError()
		}
		getInterventionCriteria.MaintenanceSequence = index + 1
		getInterventionCriteria.CreatedAt = dateTime
		getInterventionCriteria.CreatedBy = userID
		getInterventionCriteria.UpdatedAt = dateTime
		getInterventionCriteria.UpdatedBy = userID
		err = su.settingRepo.UpdateInterventionCriteria(&getInterventionCriteria)
		if err != nil {
			log.Println(err)
			return requests.InterventionCriteriaSequenceCriteriaMethod{}, responses.NewInternalServerError()
		}

	}

	for index, item := range params.Concrete {
		getInterventionCriteria, err := su.settingRepo.GetInterventionCriteriaById(item)
		if err != nil {
			log.Println(err)
			return requests.InterventionCriteriaSequenceCriteriaMethod{}, responses.NewInternalServerError()
		}
		getInterventionCriteria.MaintenanceSequence = index + 1
		getInterventionCriteria.CreatedAt = dateTime
		getInterventionCriteria.CreatedBy = userID
		getInterventionCriteria.UpdatedAt = dateTime
		getInterventionCriteria.UpdatedBy = userID
		err = su.settingRepo.UpdateInterventionCriteria(&getInterventionCriteria)
		if err != nil {
			log.Println(err)
			return requests.InterventionCriteriaSequenceCriteriaMethod{}, responses.NewInternalServerError()
		}

	}

	err := su.MergeInterventionCriteria(userID)
	if err != nil {
		log.Println(err)
		return requests.InterventionCriteriaSequenceCriteriaMethod{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) PostRefSurface(request requests.RefSurface, uID int, isC2Nil bool) error {
	tx := su.settingRepo.StartTransSection()
	var dataToInsertParam models.RefSurfaceParam

	if helpers.CountDecimal(request.A) > 6 {
		s := fmt.Sprintf("%.*f", 6, request.A)
		rounded, _ := strconv.ParseFloat(s, 64)
		request.A = rounded
	} else if helpers.CountDecimal(request.B) > 6 {
		s := fmt.Sprintf("%.*f", 6, request.B)
		rounded, _ := strconv.ParseFloat(s, 64)
		request.B = rounded
	}
	var c string
	c1 := strconv.Itoa(request.C1)
	c2 := strconv.Itoa(request.C2)
	if isC2Nil {
		c = c1
	} else {
		c = fmt.Sprintf("%s^%s", c1, c2)
	}
	dataToInsert := models.NewRefSurface{}
	err := copier.Copy(&dataToInsert, &request)
	if err != nil {
		logs.Error(err)
		su.settingRepo.RollBack(tx)
		return err
	}
	dataToInsert.C = c
	newRefID, err := su.settingRepo.InsertRefSurface(tx, dataToInsert)
	if err != nil {
		logs.Error(err)
		su.settingRepo.RollBack(tx)
		return err
	}
	dataToInsert.ID = newRefID
	dataToInsert.CanDelete = true
	jsonBytes, err := json.Marshal(dataToInsert)
	if err != nil {
		logs.Error(err)
		su.settingRepo.RollBack(tx)
		return err
	}
	jsonString := string(jsonBytes)
	str := strings.ReplaceAll(jsonString, "\\", "")
	dataToInsertParam.CreateDate = time.Now()
	dataToInsertParam.IsLatest = true
	dataToInsertParam.CreateBy = uID
	dataToInsertParam.Params = str
	dataToInsertParam.RefSurfaceID = newRefID
	err = su.settingRepo.InsertRefSurfaceParam(tx, dataToInsertParam)
	if err != nil {
		logs.Error(err)
		su.settingRepo.RollBack(tx)
		return err
	}
	su.settingRepo.Commit(tx)
	return nil
}

func (su *settingUseCase) PutRefSurface(request requests.RefSurface, uID int, id int, isC2Nil bool) error {
	tx := su.settingRepo.StartTransSection()
	var dataToInsertParam models.RefSurfaceParam

	if helpers.CountDecimal(request.A) > 6 {
		s := fmt.Sprintf("%.*f", 6, request.A)
		rounded, _ := strconv.ParseFloat(s, 64)
		request.A = rounded
	} else if helpers.CountDecimal(request.B) > 6 {
		s := fmt.Sprintf("%.*f", 6, request.B)
		rounded, _ := strconv.ParseFloat(s, 64)
		request.B = rounded
	}
	var c string
	c1 := strconv.Itoa(request.C1)
	c2 := strconv.Itoa(request.C2)
	if isC2Nil {
		c = c1
	} else {
		c = fmt.Sprintf("%s^%s", c1, c2)
	}
	dataToInsert := models.NewRefSurface{}
	err := copier.Copy(&dataToInsert, &request)
	if err != nil {
		logs.Error(err)
		su.settingRepo.RollBack(tx)
		return err
	}
	dataToInsert.C = c
	dataToInsert.ID = id
	dataToInsert.CanDelete, err = su.settingRepo.UpdateRefSurface(tx, dataToInsert)
	if err != nil {
		logs.Error(err)
		su.settingRepo.RollBack(tx)
		return err
	}
	refSurface, err := su.settingRepo.GetNewRefSurfaceByID(dataToInsert.ID)
	if err != nil {
		logs.Error(err)
		su.settingRepo.RollBack(tx)
		return err
	}

	jsonBytes, err := json.Marshal(refSurface)
	if err != nil {
		logs.Error(err)
		su.settingRepo.RollBack(tx)
		return err
	}
	jsonString := string(jsonBytes)
	str := strings.ReplaceAll(jsonString, "\\", "")
	dataToInsertParam.CreateDate = time.Now()
	dataToInsertParam.IsLatest = true
	dataToInsertParam.CreateBy = uID
	dataToInsertParam.Params = str
	dataToInsertParam.RefSurfaceID = id
	err = su.settingRepo.InsertRefSurfaceParam(tx, dataToInsertParam)
	if err != nil {
		logs.Error(err)
		su.settingRepo.RollBack(tx)
		return err
	}
	su.settingRepo.Commit(tx)
	return nil
}

func (su *settingUseCase) GetRefSurface(condition string) ([]models.NewRefSurface, error) {
	results, err := su.settingRepo.GetRefSurface(condition)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	err = helpers.SortStructByID(results, false)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return results, nil
}

func (su *settingUseCase) GetRefSurfaceByID(id int) (interface{}, error) {
	result, err := su.settingRepo.GetRefSurfaceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logs.Error(err)
		return nil, err
	}
	return result, nil
}

func (su *settingUseCase) GetParamRefSurface(id int) (interface{}, error) {
	var responds []responses.RefSurfaceParam
	result, err := su.settingRepo.GetParamRefSurfaceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logs.Error(err)
		return nil, err
	}
	for _, v := range result {
		var respond responses.RefSurfaceParam
		var jsonObject map[string]interface{}
		err = json.Unmarshal([]byte(v.Params), &jsonObject)
		if err != nil {
			logs.Error(err)
			return result, err
		}
		err = copier.Copy(&respond, &v)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		respond.JsonParams = jsonObject
		responds = append(responds, respond)
	}

	return responds, nil
}

func (su *settingUseCase) CreateBudget(params requests.Budget, c *gin.Context) (requests.Budget, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	getBudget, err := su.settingRepo.GetBudget()
	if err != nil {
		log.Println(err)
	}

	for _, item := range getBudget {
		if params.Name == item.Name {
			return requests.Budget{}, responses.NewDuplicatedNameError("name")
		}
	}

	var budget models.SettingBudget
	budget.Name = params.Name
	budget.CanDelete = true
	budget.IsDeleted = false
	budget.CreatedBy = userID
	budget.CreatedAt = dateTime

	err = su.settingRepo.UpdateBudget(&budget)
	if err != nil {
		log.Println(err)
		return requests.Budget{}, responses.NewInternalServerError()
	}

	for _, item := range params.Budget {
		var budgetMethod models.SettingBudgetMethod
		budgetMethod.BudgetId = budget.Id
		budgetMethod.MethodName = item.MethodName
		budgetMethod.CostPerUnit = item.CostPerUnit
		budgetMethod.IsDeleted = false
		budgetMethod.CreatedBy = userID
		budgetMethod.CreatedAt = dateTime
		budgetMethod.UpdatedBy = userID
		budgetMethod.UpdatedAt = dateTime
		err := su.settingRepo.UpdateBudgetMethod(&budgetMethod)
		if err != nil {
			log.Println(err)
			return requests.Budget{}, responses.NewInternalServerError()
		}
	}

	return params, nil
}

func (su *settingUseCase) GetBudget(params requests.QueryParams, c *gin.Context) (interface{}, error) {

	total, err := su.settingRepo.CountGetBudgetByName(params.Name)
	if err != nil {
		log.Println(err)
		return []responses.BudgetList{}, responses.NewInternalServerError()
	}

	limit, offset, page := helpers.GetlimitOffsetPage(params.Limit, params.Page, total)

	assetGroups, err := su.settingRepo.GetBudgetByName(limit, offset, params.Name)
	if err != nil {
		log.Println(err)
		return responses.Pagination{}, responses.NewInternalServerError()
	}

	pagination := helpers.Pagination(assetGroups, limit, page, total)

	return pagination, nil
}

func (su *settingUseCase) GetBudgetById(id string, c *gin.Context) (responses.Budget, error) {
	idInt, _ := strconv.Atoi(id)
	getBudget, err := su.settingRepo.GetBudgetById(idInt)
	if err != nil {
		if err.Error() == "record not found" {
			log.Println(err)
			return responses.Budget{}, responses.NewNotFoundError()
		} else {
			log.Println(err)
			return responses.Budget{}, responses.NewInternalServerError()
		}
	}
	getBudgetMothod, err := su.settingRepo.GetBudgetMethodByBudgetId(idInt)
	if err != nil {
		if err.Error() == "record not found" {
			log.Println(err)
		} else {
			log.Println(err)
			return responses.Budget{}, responses.NewInternalServerError()
		}
	}
	var budget responses.Budget
	budget.Id = getBudget.Id
	budget.Name = getBudget.Name
	budget.CanDelete = getBudget.CanDelete
	var budgetMethods []responses.BudgetMethod
	if len(getBudgetMothod) > 0 {
		for _, item := range getBudgetMothod {
			var budgetMethod responses.BudgetMethod
			budgetMethod.Id = item.Id
			budgetMethod.MethodName = item.MethodName
			budgetMethod.CostPerUnit = item.CostPerUnit
			budgetMethod.IsShowMethod = item.IsShowMethod
			budgetMethods = append(budgetMethods, budgetMethod)
		}
	}

	budget.Budget = budgetMethods
	return budget, nil
}

func (su *settingUseCase) UpdateBudget(params requests.UpdateBudget, c *gin.Context) (requests.UpdateBudget, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	getBudgetAll, err := su.settingRepo.GetBudget()
	if err != nil {
		log.Println(err)
	}

	for _, item := range getBudgetAll {
		if params.Name == item.Name && params.Id != item.Id {
			return requests.UpdateBudget{}, responses.NewDuplicatedNameError("Name is duplicate")
		}
	}

	getBudget, err := su.settingRepo.GetBudgetById(params.Id)
	if err != nil {
		log.Println(err)
		return requests.UpdateBudget{}, responses.NewInternalServerError()
	}

	var budget models.SettingBudget
	budget.Id = params.Id
	budget.Name = params.Name
	budget.CanDelete = getBudget.CanDelete
	budget.IsDeleted = false
	budget.CreatedBy = getBudget.CreatedBy
	budget.CreatedAt = getBudget.CreatedAt
	budget.UpdatedBy = userID
	budget.UpdatedAt = dateTime

	err = su.settingRepo.UpdateBudget(&budget)
	if err != nil {
		log.Println(err)
		return requests.UpdateBudget{}, responses.NewInternalServerError()
	}

	var idList = []int{}
	for _, item := range params.Budget {
		if item.Id != 0 {
			idList = append(idList, item.Id)
		}
	}

	getBudgetList, _ := su.settingRepo.GetBudgetMethodListByBudgetId(params.Id)
	for _, item := range getBudgetList {
		if !helpers.MatchInteger(idList, item.Id) {
			getBudget, _ := su.settingRepo.GetBudgetMethodById(item.Id)
			getBudget.IsDeleted = true
			err := su.settingRepo.UpdateBudgetMethod(&getBudget)
			if err != nil {
				log.Println(err)
				return requests.UpdateBudget{}, responses.NewInternalServerError()
			}
		}
	}

	for _, item := range params.Budget {
		var budgetMethod models.SettingBudgetMethod
		getBudget, err := su.settingRepo.GetBudgetMethodById(item.Id)
		if (getBudget != models.SettingBudgetMethod{}) {
			budgetMethod.Id = getBudget.Id
			budgetMethod.BudgetId = budget.Id
			budgetMethod.MethodName = item.MethodName
			budgetMethod.CostPerUnit = item.CostPerUnit
			budgetMethod.IsShowMethod = getBudget.IsShowMethod
			budgetMethod.CreatedBy = getBudget.CreatedBy
			budgetMethod.CreatedAt = getBudget.CreatedAt
			budgetMethod.UpdatedBy = userID
			budgetMethod.UpdatedAt = dateTime
		} else if (getBudget == models.SettingBudgetMethod{}) {
			budgetMethod.BudgetId = budget.Id
			budgetMethod.MethodName = item.MethodName
			budgetMethod.CostPerUnit = item.CostPerUnit
			budgetMethod.IsShowMethod = false
			budgetMethod.IsDeleted = false
			budgetMethod.CreatedBy = userID
			budgetMethod.CreatedAt = dateTime
		} else if err != nil {
			log.Println(err)
			return requests.UpdateBudget{}, responses.NewInternalServerError()
		}

		err = su.settingRepo.UpdateBudgetMethod(&budgetMethod)
		if err != nil {
			log.Println(err)
			return requests.UpdateBudget{}, responses.NewInternalServerError()
		}

	}

	return params, nil
}

func (su *settingUseCase) DeleteBudgetById(id string, c *gin.Context) (interface{}, error) {
	idInt, _ := strconv.Atoi(id)
	getBudget, err := su.settingRepo.GetBudgetById(idInt)
	if err != nil {
		if err.Error() == "record not found" {
			log.Println(err)
			return responses.RoadGroupWithVolumeAadt{}, responses.NewNotFoundError()
		}
		log.Println(err)
		return responses.Budget{}, responses.NewInternalServerError()
	}

	if getBudget.CanDelete {
		getBudget.IsDeleted = true

		err = su.settingRepo.UpdateBudget(&getBudget)
		if err != nil {
			log.Println(err)
			return responses.Budget{}, responses.NewInternalServerError()
		}
	} else {
		return responses.RoadGroupWithVolumeAadt{}, responses.NewNotFoundError()
	}
	return nil, nil
}

func (su *settingUseCase) CreateAadtGrowthRate(params []requests.CreateAadtGrowthRate, c *gin.Context) ([]requests.CreateAadtGrowthRate, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()
	var createsGrowthRate []models.AadtGrowthRate
	for _, item := range params {
		var createGrowthRate models.AadtGrowthRate
		getAadtGrowthRate, err := su.settingRepo.GetAadtGrowthRateByRoadGroupId(item.RoadGroupId)

		if (getAadtGrowthRate == models.AadtGrowthRate{}) {
			createGrowthRate.RoadGroupId = item.RoadGroupId
			createGrowthRate.R = item.R
			createGrowthRate.IsLatest = true
			createGrowthRate.IsDeleted = false
			createGrowthRate.CreatedBy = userID
			createGrowthRate.CreatedAt = dateTime
		} else if (getAadtGrowthRate != models.AadtGrowthRate{}) {
			createGrowthRate.Id = getAadtGrowthRate.Id
			createGrowthRate.RoadGroupId = getAadtGrowthRate.RoadGroupId
			createGrowthRate.R = item.R
			createGrowthRate.IsLatest = true
			createGrowthRate.IsDeleted = false
			createGrowthRate.CreatedBy = userID
			createGrowthRate.CreatedAt = dateTime
			createGrowthRate.UpdatedBy = userID
			createGrowthRate.UpdatedAt = dateTime
		} else if err != nil {
			log.Println(err)
			return []requests.CreateAadtGrowthRate{}, responses.NewInternalServerError()
		}

		err = su.settingRepo.UpdateAadtGrowthRate(&createGrowthRate)
		if err != nil {
			log.Println(err)
			return []requests.CreateAadtGrowthRate{}, responses.NewInternalServerError()
		}
		createsGrowthRate = append(createsGrowthRate, createGrowthRate)
	}

	_, err := su.MergeSettingAadtToParams(userID)
	if err != nil {
		log.Println(err)
		return []requests.CreateAadtGrowthRate{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetAadtGrowthRate(c *gin.Context) ([]responses.AadtGrowthRate, error) {
	getAddtGrowthRate, err := su.settingRepo.GetAadtGrowthRate()
	if err != nil {
		log.Println(err)
		return []responses.AadtGrowthRate{}, responses.NewInternalServerError()
	}

	var aadtGrowthRateList []responses.AadtGrowthRate
	for _, item := range getAddtGrowthRate {
		var aadtGrowthRate responses.AadtGrowthRate

		aadtGrowthRate.RoadGroupId = item.RoadGroupId
		aadtGrowthRate.RoadGroupName = item.RoadGroupName
		aadtGrowthRate.R = item.R
		aadtGrowthRate.Code = item.Number
		aadtGrowthRateList = append(aadtGrowthRateList, aadtGrowthRate)
	}

	return aadtGrowthRateList, nil
}

func (su *settingUseCase) CreateAadtPercentageVehicleType(params requests.AadtPercentageVehicleType, c *gin.Context) (requests.AadtPercentageVehicleType, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.AadtPercentageVehicleType{}, err
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.AadtPercentageVehicleType{}, err
	}

	var aadtPercentageVehicleTypeParams models.AadtPercentageVehicleTypeParams
	aadtPercentageVehicleTypeParams.RoadGroupId = params.RoadGroupId
	aadtPercentageVehicleTypeParams.Params = buffer.String()
	aadtPercentageVehicleTypeParams.IsDeleted = false
	aadtPercentageVehicleTypeParams.IsLatest = true
	aadtPercentageVehicleTypeParams.CreatedBy = userID
	aadtPercentageVehicleTypeParams.CreatedAt = dateTime
	aadtPercentageVehicleTypeParams.UpdatedBy = userID
	aadtPercentageVehicleTypeParams.UpdatedAt = dateTime

	err = su.settingRepo.UpdateAadtPercentageVehicleTypeByIsLatestIsFalseAndRoadGroupId(params.RoadGroupId)
	if err != nil {
		log.Println(err)
		return requests.AadtPercentageVehicleType{}, responses.NewInternalServerError()
	}

	err = su.settingRepo.UpdateAadtPercentageVehicleType(&aadtPercentageVehicleTypeParams)
	if err != nil {
		log.Println(err)
		return requests.AadtPercentageVehicleType{}, responses.NewInternalServerError()
	}

	_, err = su.MergeSettingAadtToParams(userID)
	if err != nil {
		log.Println(err)
		return requests.AadtPercentageVehicleType{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetAadtPercentageVehicleType(c *gin.Context, id string) (responses.AadtPercentageVehicleType, error) {
	idInteger, err := helpers.ConvertStringToInt(id)
	aadtPercentageVehicleTypeParams, err := su.settingRepo.GetAadtPercentageVehicleTypeWithRoadGroupByRoadGroupId(idInteger)
	if err != nil {
		if err.Error() == "record not found" {
			return responses.AadtPercentageVehicleType{}, nil
		}
		log.Println(err)
		return responses.AadtPercentageVehicleType{}, err
	}

	var aadtPercentageVehicleType responses.AadtPercentageVehicleType
	err = json.Unmarshal([]byte(aadtPercentageVehicleTypeParams.Params), &aadtPercentageVehicleType)
	if err != nil {
		log.Println(err)
		return responses.AadtPercentageVehicleType{}, err
	}

	return aadtPercentageVehicleType, nil
}

func (su *settingUseCase) GetAadtParameterRoadGroupWithVolumeAadt(c *gin.Context) ([]responses.RoadGroupWithVolumeAadt, error) {

	getRoadGroup, err := su.settingRepo.GetRoadGroup()
	if err != nil {
		if err.Error() == "record not found" {
			return []responses.RoadGroupWithVolumeAadt{}, responses.NewNotFoundError()
		}
		return []responses.RoadGroupWithVolumeAadt{}, responses.NewInternalServerError()
	}
	var roadGroupWithVolumeAadtList []responses.RoadGroupWithVolumeAadt
	for _, item := range getRoadGroup {
		volumeAadt, err := su.settingRepo.GetVolumeByRoadGroupId(item.Id, "A")
		if err != nil {
			if err.Error() == "record not found" {
				log.Println(err)
			} else {
				log.Println(err)
				return []responses.RoadGroupWithVolumeAadt{}, responses.NewInternalServerError()
			}
		}

		var roadGroupWithVolumeAadt responses.RoadGroupWithVolumeAadt
		roadGroupWithVolumeAadt.RoadGroupId = item.Id
		roadGroupWithVolumeAadt.RoadGroupName = item.Name
		if (volumeAadt != models.VolumeAadt{}) {
			roadGroupWithVolumeAadt.VolumeAadt.Veh1 = volumeAadt.Veh1
			roadGroupWithVolumeAadt.VolumeAadt.Veh2 = volumeAadt.Veh2
			roadGroupWithVolumeAadt.VolumeAadt.Veh3 = volumeAadt.Veh3
			// roadGroupWithVolumeAadt.VolumeAadt.Veh4 = volumeAadt.Veh4
			// roadGroupWithVolumeAadt.VolumeAadt.Calculate.FourWheelTotal = volumeAadt.Veh1 + volumeAadt.Veh2
			// roadGroupWithVolumeAadt.VolumeAadt.Calculate.SixToTenWheelTotal = volumeAadt.Veh3
			// roadGroupWithVolumeAadt.VolumeAadt.Calculate.TenWheelTotal = volumeAadt.Veh4
			// roadGroupWithVolumeAadt.VolumeAadt.Calculate.SixToTenWheelPercentage = (float64(volumeAadt.Veh3) / float64(volumeAadt.Veh3+volumeAadt.Veh4)) * 100
			// roadGroupWithVolumeAadt.VolumeAadt.Calculate.TenWheelPercentage = (float64(volumeAadt.Veh4) / float64(volumeAadt.Veh3+volumeAadt.Veh4)) * 100
		}
		roadGroupWithVolumeAadtList = append(roadGroupWithVolumeAadtList, roadGroupWithVolumeAadt)
	}

	return roadGroupWithVolumeAadtList, nil
}

func (su *settingUseCase) CreateAadtParameter(params requests.CreateAadtParameter, c *gin.Context) (models.AadtParameter, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	var crateAadtParameter models.AadtParameter
	getAadtParameter, err := su.settingRepo.GetAadtParameterByRoadGroupId(params.RoadGroupId)

	if (getAadtParameter == models.AadtParameter{}) {
		crateAadtParameter.RoadGroupId = params.RoadGroupId
		crateAadtParameter.Elane = params.Elane
		crateAadtParameter.FourWheelAxleNumber = params.FourWheelAxleNumber
		crateAadtParameter.FourWheelVehicleVolume = params.FourWheelVehicleVolume
		crateAadtParameter.SixWheelAxleNumberId = params.SixWheelAxleNumberId
		crateAadtParameter.SixWheelVehicleVolume = params.SixWheelVehicleVolume
		crateAadtParameter.SixWheelPercentageTruck = params.SixWheelPercentageTruck
		crateAadtParameter.SixWheelFactorResult = params.SixWheelFactorResult
		crateAadtParameter.TenWheelAxleNumberId = params.TenWheelAxleNumberId
		crateAadtParameter.TenWheelVehicleVolume = params.TenWheelVehicleVolume
		crateAadtParameter.TenWheelPercentageTruck = params.TenWheelPercentageTruck
		crateAadtParameter.TenWheelFactorResult = params.TenWheelFactorResult
		crateAadtParameter.IsTruckFactor = params.IsTruckFactor
		crateAadtParameter.SpeedAverage = params.SpeedAverage
		crateAadtParameter.SpeedHeavyTruck = params.SpeedHeavyTruck
		crateAadtParameter.LaneDistributionFactor = params.LaneDistributionFactor
		crateAadtParameter.DirectionalDistributionFactor = params.DirectionalDistributionFactor
		crateAadtParameter.IsLatest = true
		crateAadtParameter.IsDeleted = false
		crateAadtParameter.CreatedBy = userID
		crateAadtParameter.CreatedAt = dateTime
	} else if (getAadtParameter != models.AadtParameter{}) {
		crateAadtParameter.Id = getAadtParameter.Id
		crateAadtParameter.RoadGroupId = params.RoadGroupId
		crateAadtParameter.Elane = params.Elane
		crateAadtParameter.FourWheelAxleNumber = params.FourWheelAxleNumber
		crateAadtParameter.FourWheelVehicleVolume = params.FourWheelVehicleVolume
		crateAadtParameter.SixWheelAxleNumberId = params.SixWheelAxleNumberId
		crateAadtParameter.SixWheelVehicleVolume = params.SixWheelVehicleVolume
		crateAadtParameter.SixWheelPercentageTruck = params.SixWheelPercentageTruck
		crateAadtParameter.SixWheelFactorResult = params.SixWheelFactorResult
		crateAadtParameter.TenWheelAxleNumberId = params.TenWheelAxleNumberId
		crateAadtParameter.TenWheelVehicleVolume = params.TenWheelVehicleVolume
		crateAadtParameter.TenWheelPercentageTruck = params.TenWheelPercentageTruck
		crateAadtParameter.TenWheelFactorResult = params.TenWheelFactorResult
		crateAadtParameter.IsTruckFactor = params.IsTruckFactor
		crateAadtParameter.SpeedAverage = params.SpeedAverage
		crateAadtParameter.SpeedHeavyTruck = params.SpeedHeavyTruck
		crateAadtParameter.LaneDistributionFactor = params.LaneDistributionFactor
		crateAadtParameter.DirectionalDistributionFactor = params.DirectionalDistributionFactor
		crateAadtParameter.RoadGroupId = getAadtParameter.RoadGroupId
		crateAadtParameter.CreatedBy = userID
		crateAadtParameter.CreatedAt = dateTime
		crateAadtParameter.UpdatedBy = userID
		crateAadtParameter.UpdatedAt = dateTime

	} else if err != nil {
		log.Println(err)
		return models.AadtParameter{}, responses.NewInternalServerError()
	}

	err = su.settingRepo.UpdateAadtParameter(&crateAadtParameter)
	if err != nil {
		log.Println(err)
		return models.AadtParameter{}, responses.NewInternalServerError()
	}

	_, err = su.MergeSettingAadtToParams(userID)
	if err != nil {
		log.Println(err)
		return models.AadtParameter{}, responses.NewInternalServerError()
	}

	return crateAadtParameter, nil
}

func (su *settingUseCase) GetAadtParameter(c *gin.Context, id string) (responses.AadtParameter, error) {
	idInteger, err := helpers.ConvertStringToInt(id)
	getAadtParameter, err := su.settingRepo.GetAadtParameterByRoadGroupId(idInteger)
	if err != nil {
		log.Println(err)
		return responses.AadtParameter{}, responses.NewInternalServerError()
	}

	var aadtParameter responses.AadtParameter
	aadtParameter.RoadGroupId = getAadtParameter.RoadGroupId
	aadtParameter.Elane = getAadtParameter.Elane
	aadtParameter.FourWheelAxleNumber = getAadtParameter.FourWheelAxleNumber
	aadtParameter.FourWheelVehicleVolume = getAadtParameter.FourWheelVehicleVolume
	aadtParameter.SixWheelAxleNumberId = getAadtParameter.SixWheelAxleNumberId
	aadtParameter.SixWheelVehicleVolume = getAadtParameter.SixWheelVehicleVolume
	aadtParameter.SixWheelPercentageTruck = getAadtParameter.SixWheelPercentageTruck
	aadtParameter.SixWheelFactorResult = getAadtParameter.SixWheelFactorResult
	aadtParameter.TenWheelAxleNumberId = getAadtParameter.TenWheelAxleNumberId
	aadtParameter.TenWheelVehicleVolume = getAadtParameter.TenWheelVehicleVolume
	aadtParameter.TenWheelPercentageTruck = getAadtParameter.TenWheelPercentageTruck
	aadtParameter.TenWheelFactorResult = getAadtParameter.TenWheelFactorResult
	aadtParameter.IsTruckFactor = getAadtParameter.IsTruckFactor
	aadtParameter.SpeedAverage = getAadtParameter.SpeedAverage
	aadtParameter.SpeedHeavyTruck = getAadtParameter.SpeedHeavyTruck
	aadtParameter.LaneDistributionFactor = getAadtParameter.LaneDistributionFactor
	aadtParameter.DirectionalDistributionFactor = getAadtParameter.DirectionalDistributionFactor

	return aadtParameter, nil
}

func (su *settingUseCase) CalculateAadtParameterForOverSix(params requests.CalculateAadtParameter, c *gin.Context) (responses.CalculateAadtParameterTruck, error) {
	// volumeAadt, err := su.settingRepo.GetVolumeByRoadGroupId(params.RoadGroupId, "A")
	// if err != nil {
	// 	if err.Error() == "record not found" {
	// 		return responses.CalculateAadtParameterTruck{}, responses.NewNotFoundError()
	// 	}

	// 	log.Println(err)
	// 	return responses.CalculateAadtParameterTruck{}, responses.NewInternalServerError()
	// }

	parameterVehicleType, err := su.settingRepo.GetParameterVehicleTypeById(params.ParameterVehicleTypeId)
	if err != nil {
		if err.Error() == "record not found" {
			return responses.CalculateAadtParameterTruck{}, responses.NewNotFoundError()
		}
		return responses.CalculateAadtParameterTruck{}, responses.NewInternalServerError()
	}

	// truckTotal := float64(volumeAadt.Veh3 + volumeAadt.Veh4)

	var calculateAadtParameterTruck responses.CalculateAadtParameterTruck
	// sixWheel := float64(volumeAadt.Veh3)
	// percentageTruck := (sixWheel / truckTotal) * 100
	// loadEquivalent := parameterVehicleType.LoadEquivalent
	// truckFactor := percentageTruck * loadEquivalent
	// calculateAadtParameterTruck.VehicleVolume = truckTotal
	// calculateAadtParameterTruck.PercentageTruck = percentageTruck
	calculateAadtParameterTruck.LoadEquivalent = parameterVehicleType.LoadEquivalent
	// calculateAadtParameterTruck.TruckFactor = truckFactor

	return calculateAadtParameterTruck, nil
}

func (su *settingUseCase) CalculateAadtParameterForOverTen(params requests.CalculateAadtParameter, c *gin.Context) (responses.CalculateAadtParameterTruck, error) {
	// volumeAadt, err := su.settingRepo.GetVolumeByRoadGroupId(params.RoadGroupId, "A")
	// if err != nil {
	// 	if err.Error() == "record not found" {
	// 		return responses.CalculateAadtParameterTruck{}, responses.NewNotFoundError()
	// 	}

	// 	log.Println(err)
	// 	return responses.CalculateAadtParameterTruck{}, responses.NewInternalServerError()
	// }

	parameterVehicleType, err := su.settingRepo.GetParameterVehicleTypeById(params.ParameterVehicleTypeId)
	if err != nil {
		if err.Error() == "record not found" {
			return responses.CalculateAadtParameterTruck{}, responses.NewNotFoundError()
		}
		return responses.CalculateAadtParameterTruck{}, responses.NewInternalServerError()
	}

	// truckTotal := float64(volumeAadt.Veh3 + volumeAadt.Veh4)
	var calculateAadtParameterTruck responses.CalculateAadtParameterTruck
	// tenWheel := float64(volumeAadt.Veh4)
	// percentageTruck := (tenWheel / truckTotal) * 100
	loadEquivalent := parameterVehicleType.LoadEquivalent
	// truckFactor := percentageTruck * loadEquivalent
	// calculateAadtParameterTruck.VehicleVolume = truckTotal
	// calculateAadtParameterTruck.PercentageTruck = percentageTruck
	calculateAadtParameterTruck.LoadEquivalent = loadEquivalent
	// calculateAadtParameterTruck.TruckFactor = truckFactor

	return calculateAadtParameterTruck, nil
}

func (su *settingUseCase) MergeSettingAadtToParams(userID int) (interface{}, error) {
	growthRate, err := su.settingRepo.GetAllAadtGrowthRate()
	if err != nil {
		log.Println(err)
		return nil, responses.NewInternalServerError()
	}

	percentageVehicleType, err := su.settingRepo.GetAllAadtPercentageVehicleType()
	if err != nil {
		log.Println(err)
		return nil, responses.NewInternalServerError()
	}

	parameter, err := su.settingRepo.GetAllAadtParameter()
	if err != nil {
		log.Println(err)
		return nil, responses.NewInternalServerError()
	}

	err = su.settingRepo.UpdateAadtParamsByIsLatestIsFalse()
	if err != nil {
		log.Println(err)
		return models.AadtParameter{}, responses.NewInternalServerError()
	}
	var percentageVehicleTypes []models.AadtPercentageVehicleType

	for _, item := range percentageVehicleType {
		var percentageVehicleType models.AadtPercentageVehicleType
		err = json.Unmarshal([]byte(item.Params), &percentageVehicleType)
		if err != nil {
			log.Println(err)
			return responses.AadtPercentageVehicleType{}, err
		}
		percentageVehicleTypes = append(percentageVehicleTypes, percentageVehicleType)
	}

	var getAadtParams models.GetAadtParams
	getAadtParams.GetAadtParameter = parameter
	getAadtParams.GetAadtPercentageVehicleType = percentageVehicleTypes
	getAadtParams.GetAadtGrowthRate = growthRate

	structByte, err := json.Marshal(getAadtParams)
	if err != nil {
		log.Println(err)
		return models.AadtParameter{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, structByte); err != nil {
		log.Println(err)
		return []string{}, responses.NewInternalServerError()
	}

	var aadtParams models.AadtParams
	aadtParams.Params = buffer.String()
	aadtParams.IsLatest = true
	aadtParams.CreatedBy = userID
	aadtParams.CreatedAt = time.Now()

	err = su.settingRepo.UpdateAadtParams(aadtParams)
	if err != nil {
		log.Println(err)
		return models.AadtParameter{}, responses.NewInternalServerError()
	}

	return getAadtParams, nil
}

func (su *settingUseCase) CreateRoadWorkEffectAsphalt(params requests.SettingRoadWorkEffectAsphalt, c *gin.Context) (responses.SettingRoadWorkEffectAsphalt, error) {

	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	getRoadWorkEffect, err := su.settingRepo.GetRoadWorkEffect()
	var settingRoadWorkEffect models.SettingRoadWorkEffect
	if (getRoadWorkEffect == models.SettingRoadWorkEffect{}) {
		settingRoadWorkEffect.AsOlOverlayA0 = params.AsOlOverlayA0
		settingRoadWorkEffect.AsOlArVb = params.AsOlArVb
		settingRoadWorkEffect.AsOlPoTb = params.AsOlPoTb
		settingRoadWorkEffect.AsOlAcAb = params.AsOlAcAb
		settingRoadWorkEffect.AsOlRdMb = params.AsOlRdMb
		settingRoadWorkEffect.AsSsRweSsModelA0 = params.AsSsRweSsModelA0
		settingRoadWorkEffect.AsSsDefaultLowerBoundIriAfterSlurrySeal = params.AsSsDefaultLowerBoundIriAfterSlurrySeal
		settingRoadWorkEffect.AsSsArVb = params.AsSsArVb
		settingRoadWorkEffect.AsSsApoTb = params.AsSsApoTb
		settingRoadWorkEffect.AsSsAcAb = params.AsSsAcAb
		settingRoadWorkEffect.AsSsRdMb = params.AsSsRdMb
		settingRoadWorkEffect.AsMolIriAfterMillOverlay = params.AsMolIriAfterMillOverlay
		settingRoadWorkEffect.AsMolArVb = params.AsMolArVb
		settingRoadWorkEffect.AsMolApoTb = params.AsMolApoTb
		settingRoadWorkEffect.AsMolAcAb = params.AsMolAcAb
		settingRoadWorkEffect.AsMolRdMb = params.AsMolRdMb
		settingRoadWorkEffect.AsRclSnc = params.AsRclSnc
		settingRoadWorkEffect.AsRclIriAfterRecycling = params.AsRclIriAfterRecycling
		settingRoadWorkEffect.AsRclArVb = params.AsRclArVb
		settingRoadWorkEffect.AsRclApoTb = params.AsRclApoTb
		settingRoadWorkEffect.AsRclAcAb = params.AsRclAcAb
		settingRoadWorkEffect.AsRclRdMb = params.AsRclRdMb
		settingRoadWorkEffect.AsRclDefaultHsOld = params.AsRclDefaultHsOld
		settingRoadWorkEffect.AsRcSnc = params.AsRcSnc
		settingRoadWorkEffect.AsRcIriAfterReconstruction = params.AsRcIriAfterReconstruction
		settingRoadWorkEffect.AsRcArVb = params.AsRcArVb
		settingRoadWorkEffect.AsRcApoTb = params.AsRcApoTb
		settingRoadWorkEffect.AsRcAcAb = params.AsRcAcAb
		settingRoadWorkEffect.AsRcRdMb = params.AsRcRdMb
		settingRoadWorkEffect.IsDeleted = false
		settingRoadWorkEffect.IsLatest = true
		settingRoadWorkEffect.CreatedBy = userID
		settingRoadWorkEffect.CreatedAt = dateTime
	} else if (getRoadWorkEffect != models.SettingRoadWorkEffect{}) {
		settingRoadWorkEffect.Id = getRoadWorkEffect.Id
		settingRoadWorkEffect.AsOlOverlayA0 = params.AsOlOverlayA0
		settingRoadWorkEffect.AsOlArVb = params.AsOlArVb
		settingRoadWorkEffect.AsOlPoTb = params.AsOlPoTb
		settingRoadWorkEffect.AsOlAcAb = params.AsOlAcAb
		settingRoadWorkEffect.AsOlRdMb = params.AsOlRdMb
		settingRoadWorkEffect.AsSsRweSsModelA0 = params.AsSsRweSsModelA0
		settingRoadWorkEffect.AsSsDefaultLowerBoundIriAfterSlurrySeal = params.AsSsDefaultLowerBoundIriAfterSlurrySeal
		settingRoadWorkEffect.AsSsArVb = params.AsSsArVb
		settingRoadWorkEffect.AsSsApoTb = params.AsSsApoTb
		settingRoadWorkEffect.AsSsAcAb = params.AsSsAcAb
		settingRoadWorkEffect.AsSsRdMb = params.AsSsRdMb
		settingRoadWorkEffect.AsMolIriAfterMillOverlay = params.AsMolIriAfterMillOverlay
		settingRoadWorkEffect.AsMolArVb = params.AsMolArVb
		settingRoadWorkEffect.AsMolApoTb = params.AsMolApoTb
		settingRoadWorkEffect.AsMolAcAb = params.AsMolAcAb
		settingRoadWorkEffect.AsMolRdMb = params.AsMolRdMb
		settingRoadWorkEffect.AsRclSnc = params.AsRclSnc
		settingRoadWorkEffect.AsRclIriAfterRecycling = params.AsRclIriAfterRecycling
		settingRoadWorkEffect.AsRclArVb = params.AsRclArVb
		settingRoadWorkEffect.AsRclApoTb = params.AsRclApoTb
		settingRoadWorkEffect.AsRclAcAb = params.AsRclAcAb
		settingRoadWorkEffect.AsRclRdMb = params.AsRclRdMb
		settingRoadWorkEffect.AsRclDefaultHsOld = params.AsRclDefaultHsOld
		settingRoadWorkEffect.AsRcSnc = params.AsRcSnc
		settingRoadWorkEffect.AsRcIriAfterReconstruction = params.AsRcIriAfterReconstruction
		settingRoadWorkEffect.AsRcArVb = params.AsRcArVb
		settingRoadWorkEffect.AsRcApoTb = params.AsRcApoTb
		settingRoadWorkEffect.AsRcAcAb = params.AsRcAcAb
		settingRoadWorkEffect.AsRcRdMb = params.AsRcRdMb
		settingRoadWorkEffect.CcFdrIriAfterFdr = getRoadWorkEffect.CcFdrIriAfterFdr
		settingRoadWorkEffect.CcFdrFaulting = getRoadWorkEffect.CcFdrFaulting
		settingRoadWorkEffect.CcFdrCracking = getRoadWorkEffect.CcFdrCracking
		settingRoadWorkEffect.CcFdrSpalling = getRoadWorkEffect.CcFdrSpalling
		settingRoadWorkEffect.CcBcoIriAfterBco = getRoadWorkEffect.CcBcoIriAfterBco
		settingRoadWorkEffect.CcBcoFaulting = getRoadWorkEffect.CcBcoFaulting
		settingRoadWorkEffect.CcBcoCracking = getRoadWorkEffect.CcBcoCracking
		settingRoadWorkEffect.CcBcoSpalling = getRoadWorkEffect.CcBcoSpalling
		settingRoadWorkEffect.CcMolIriAfterMol = getRoadWorkEffect.CcMolIriAfterMol
		settingRoadWorkEffect.CcMolFaulting = getRoadWorkEffect.CcMolFaulting
		settingRoadWorkEffect.CcMolCracking = getRoadWorkEffect.CcMolCracking
		settingRoadWorkEffect.CcMolSpalling = getRoadWorkEffect.CcMolSpalling
		settingRoadWorkEffect.CcSealIriAfterSeal = getRoadWorkEffect.CcSealIriAfterSeal
		settingRoadWorkEffect.CcSealFaulting = getRoadWorkEffect.CcSealFaulting
		settingRoadWorkEffect.CcSealCracking = getRoadWorkEffect.CcSealCracking
		settingRoadWorkEffect.CcSealSpalling = getRoadWorkEffect.CcSealSpalling
		settingRoadWorkEffect.IsDeleted = false
		settingRoadWorkEffect.IsLatest = true
		settingRoadWorkEffect.CreatedBy = userID
		settingRoadWorkEffect.CreatedAt = dateTime
		settingRoadWorkEffect.UpdatedBy = userID
		settingRoadWorkEffect.UpdatedAt = dateTime
	} else if err != nil {
		log.Println(err)
		return responses.SettingRoadWorkEffectAsphalt{}, responses.NewInternalServerError()
	}

	var responsesSettingRoadWorkEffectAsphalt responses.SettingRoadWorkEffectAsphalt
	responsesSettingRoadWorkEffectAsphalt.AsOlOverlayA0 = params.AsOlOverlayA0
	responsesSettingRoadWorkEffectAsphalt.AsOlArVb = params.AsOlArVb
	responsesSettingRoadWorkEffectAsphalt.AsOlPoTb = params.AsOlPoTb
	responsesSettingRoadWorkEffectAsphalt.AsOlAcAb = params.AsOlAcAb
	responsesSettingRoadWorkEffectAsphalt.AsOlRdMb = params.AsOlRdMb
	responsesSettingRoadWorkEffectAsphalt.AsSsRweSsModelA0 = params.AsSsRweSsModelA0
	responsesSettingRoadWorkEffectAsphalt.AsSsDefaultLowerBoundIriAfterSlurrySeal = params.AsSsDefaultLowerBoundIriAfterSlurrySeal
	responsesSettingRoadWorkEffectAsphalt.AsSsArVb = params.AsSsArVb
	responsesSettingRoadWorkEffectAsphalt.AsSsApoTb = params.AsSsApoTb
	responsesSettingRoadWorkEffectAsphalt.AsSsAcAb = params.AsSsAcAb
	responsesSettingRoadWorkEffectAsphalt.AsSsRdMb = params.AsSsRdMb
	responsesSettingRoadWorkEffectAsphalt.AsMolIriAfterMillOverlay = params.AsMolIriAfterMillOverlay
	responsesSettingRoadWorkEffectAsphalt.AsMolArVb = params.AsMolArVb
	responsesSettingRoadWorkEffectAsphalt.AsMolApoTb = params.AsMolApoTb
	responsesSettingRoadWorkEffectAsphalt.AsMolAcAb = params.AsMolAcAb
	responsesSettingRoadWorkEffectAsphalt.AsMolRdMb = params.AsMolRdMb
	responsesSettingRoadWorkEffectAsphalt.AsRclSnc = params.AsRclSnc
	responsesSettingRoadWorkEffectAsphalt.AsRclIriAfterRecycling = params.AsRclIriAfterRecycling
	responsesSettingRoadWorkEffectAsphalt.AsRclArVb = params.AsRclArVb
	responsesSettingRoadWorkEffectAsphalt.AsRclApoTb = params.AsRclApoTb
	responsesSettingRoadWorkEffectAsphalt.AsRclAcAb = params.AsRclAcAb
	responsesSettingRoadWorkEffectAsphalt.AsRclRdMb = params.AsRclRdMb
	responsesSettingRoadWorkEffectAsphalt.AsRclDefaultHsOld = params.AsRclDefaultHsOld
	responsesSettingRoadWorkEffectAsphalt.AsRcSnc = params.AsRcSnc
	responsesSettingRoadWorkEffectAsphalt.AsRcIriAfterReconstruction = params.AsRcIriAfterReconstruction
	responsesSettingRoadWorkEffectAsphalt.AsRcArVb = params.AsRcArVb
	responsesSettingRoadWorkEffectAsphalt.AsRcApoTb = params.AsRcApoTb
	responsesSettingRoadWorkEffectAsphalt.AsRcAcAb = params.AsRcAcAb
	responsesSettingRoadWorkEffectAsphalt.AsRcRdMb = params.AsRcRdMb

	err = su.settingRepo.UpdateRoadWorkEffect(&settingRoadWorkEffect)
	if err != nil {
		log.Println(err)
		return responses.SettingRoadWorkEffectAsphalt{}, responses.NewInternalServerError()
	}

	err = su.MergeRoadWorkEffectToParams(userID)
	if err != nil {
		log.Println(err)
		return responses.SettingRoadWorkEffectAsphalt{}, responses.NewInternalServerError()
	}

	return responsesSettingRoadWorkEffectAsphalt, nil
}

func (su *settingUseCase) GetRoadWorkEffectAsphalt(c *gin.Context) (responses.SettingRoadWorkEffectAsphalt, error) {

	getRoadWorkEffect, err := su.settingRepo.GetRoadWorkEffect()
	if err != nil {
		log.Println(err)
		return responses.SettingRoadWorkEffectAsphalt{}, responses.NewInternalServerError()
	}

	var responsesSettingRoadWorkEffectAsphalt responses.SettingRoadWorkEffectAsphalt
	responsesSettingRoadWorkEffectAsphalt.AsOlOverlayA0 = getRoadWorkEffect.AsOlOverlayA0
	responsesSettingRoadWorkEffectAsphalt.AsOlArVb = getRoadWorkEffect.AsOlArVb
	responsesSettingRoadWorkEffectAsphalt.AsOlPoTb = getRoadWorkEffect.AsOlPoTb
	responsesSettingRoadWorkEffectAsphalt.AsOlAcAb = getRoadWorkEffect.AsOlAcAb
	responsesSettingRoadWorkEffectAsphalt.AsOlRdMb = getRoadWorkEffect.AsOlRdMb
	responsesSettingRoadWorkEffectAsphalt.AsSsRweSsModelA0 = getRoadWorkEffect.AsSsRweSsModelA0
	responsesSettingRoadWorkEffectAsphalt.AsSsDefaultLowerBoundIriAfterSlurrySeal = getRoadWorkEffect.AsSsDefaultLowerBoundIriAfterSlurrySeal
	responsesSettingRoadWorkEffectAsphalt.AsSsArVb = getRoadWorkEffect.AsSsArVb
	responsesSettingRoadWorkEffectAsphalt.AsSsApoTb = getRoadWorkEffect.AsSsApoTb
	responsesSettingRoadWorkEffectAsphalt.AsSsAcAb = getRoadWorkEffect.AsSsAcAb
	responsesSettingRoadWorkEffectAsphalt.AsSsRdMb = getRoadWorkEffect.AsSsRdMb
	responsesSettingRoadWorkEffectAsphalt.AsMolIriAfterMillOverlay = getRoadWorkEffect.AsMolIriAfterMillOverlay
	responsesSettingRoadWorkEffectAsphalt.AsMolArVb = getRoadWorkEffect.AsMolArVb
	responsesSettingRoadWorkEffectAsphalt.AsMolApoTb = getRoadWorkEffect.AsMolApoTb
	responsesSettingRoadWorkEffectAsphalt.AsMolAcAb = getRoadWorkEffect.AsMolAcAb
	responsesSettingRoadWorkEffectAsphalt.AsMolRdMb = getRoadWorkEffect.AsMolRdMb
	responsesSettingRoadWorkEffectAsphalt.AsRclSnc = getRoadWorkEffect.AsRclSnc
	responsesSettingRoadWorkEffectAsphalt.AsRclIriAfterRecycling = getRoadWorkEffect.AsRclIriAfterRecycling
	responsesSettingRoadWorkEffectAsphalt.AsRclArVb = getRoadWorkEffect.AsRclArVb
	responsesSettingRoadWorkEffectAsphalt.AsRclApoTb = getRoadWorkEffect.AsRclApoTb
	responsesSettingRoadWorkEffectAsphalt.AsRclAcAb = getRoadWorkEffect.AsRclAcAb
	responsesSettingRoadWorkEffectAsphalt.AsRclRdMb = getRoadWorkEffect.AsRclRdMb
	responsesSettingRoadWorkEffectAsphalt.AsRclDefaultHsOld = getRoadWorkEffect.AsRclDefaultHsOld
	responsesSettingRoadWorkEffectAsphalt.AsRcSnc = getRoadWorkEffect.AsRcSnc
	responsesSettingRoadWorkEffectAsphalt.AsRcIriAfterReconstruction = getRoadWorkEffect.AsRcIriAfterReconstruction
	responsesSettingRoadWorkEffectAsphalt.AsRcArVb = getRoadWorkEffect.AsRcArVb
	responsesSettingRoadWorkEffectAsphalt.AsRcApoTb = getRoadWorkEffect.AsRcApoTb
	responsesSettingRoadWorkEffectAsphalt.AsRcAcAb = getRoadWorkEffect.AsRcAcAb
	responsesSettingRoadWorkEffectAsphalt.AsRcRdMb = getRoadWorkEffect.AsRcRdMb

	return responsesSettingRoadWorkEffectAsphalt, nil
}

func (su *settingUseCase) CreateRoadWorkEffectConcrete(params requests.SettingRoadWorkEffectConcrete, c *gin.Context) (responses.SettingRoadWorkEffectConcrete, error) {

	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	getRoadWorkEffect, err := su.settingRepo.GetRoadWorkEffect()
	var settingRoadWorkEffect models.SettingRoadWorkEffect
	if (getRoadWorkEffect == models.SettingRoadWorkEffect{}) {
		settingRoadWorkEffect.CcFdrIriAfterFdr = params.CcFdrIriAfterFdr
		settingRoadWorkEffect.CcFdrFaulting = params.CcFdrFaulting
		settingRoadWorkEffect.CcFdrCracking = params.CcFdrCracking
		settingRoadWorkEffect.CcFdrSpalling = params.CcFdrSpalling
		settingRoadWorkEffect.CcBcoIriAfterBco = params.CcBcoIriAfterBco
		settingRoadWorkEffect.CcBcoFaulting = params.CcBcoFaulting
		settingRoadWorkEffect.CcBcoCracking = params.CcBcoCracking
		settingRoadWorkEffect.CcBcoSpalling = params.CcBcoSpalling
		settingRoadWorkEffect.CcMolIriAfterMol = params.CcMolIriAfterMol
		settingRoadWorkEffect.CcMolFaulting = params.CcMolFaulting
		settingRoadWorkEffect.CcMolCracking = params.CcMolCracking
		settingRoadWorkEffect.CcMolSpalling = params.CcMolSpalling
		settingRoadWorkEffect.CcSealIriAfterSeal = params.CcSealIriAfterSeal
		settingRoadWorkEffect.CcSealFaulting = params.CcSealFaulting
		settingRoadWorkEffect.CcSealCracking = params.CcSealCracking
		settingRoadWorkEffect.CcSealSpalling = params.CcSealSpalling
		settingRoadWorkEffect.CcRbcIri = params.CcRbcIri
		settingRoadWorkEffect.CcRbcSlabthk = params.CcRbcSlabthk
		settingRoadWorkEffect.CcRbcPercentFaulting = params.CcRbcPercentFaulting
		settingRoadWorkEffect.CcRbcPercentSpalling = params.CcRbcPercentSpalling
		settingRoadWorkEffect.CcRbcPercentCracking = params.CcRbcPercentCracking
		settingRoadWorkEffect.IsDeleted = false
		settingRoadWorkEffect.IsLatest = true
		settingRoadWorkEffect.CreatedBy = userID
		settingRoadWorkEffect.CreatedAt = dateTime
	} else if (getRoadWorkEffect != models.SettingRoadWorkEffect{}) {
		settingRoadWorkEffect.Id = getRoadWorkEffect.Id
		settingRoadWorkEffect.AsOlOverlayA0 = getRoadWorkEffect.AsOlOverlayA0
		settingRoadWorkEffect.AsOlArVb = getRoadWorkEffect.AsOlArVb
		settingRoadWorkEffect.AsOlPoTb = getRoadWorkEffect.AsOlPoTb
		settingRoadWorkEffect.AsOlAcAb = getRoadWorkEffect.AsOlAcAb
		settingRoadWorkEffect.AsOlRdMb = getRoadWorkEffect.AsOlRdMb
		settingRoadWorkEffect.AsSsRweSsModelA0 = getRoadWorkEffect.AsSsRweSsModelA0
		settingRoadWorkEffect.AsSsDefaultLowerBoundIriAfterSlurrySeal = getRoadWorkEffect.AsSsDefaultLowerBoundIriAfterSlurrySeal
		settingRoadWorkEffect.AsSsArVb = getRoadWorkEffect.AsSsArVb
		settingRoadWorkEffect.AsSsApoTb = getRoadWorkEffect.AsSsApoTb
		settingRoadWorkEffect.AsSsAcAb = getRoadWorkEffect.AsSsAcAb
		settingRoadWorkEffect.AsSsRdMb = getRoadWorkEffect.AsSsRdMb
		settingRoadWorkEffect.AsMolIriAfterMillOverlay = getRoadWorkEffect.AsMolIriAfterMillOverlay
		settingRoadWorkEffect.AsMolArVb = getRoadWorkEffect.AsMolArVb
		settingRoadWorkEffect.AsMolApoTb = getRoadWorkEffect.AsMolApoTb
		settingRoadWorkEffect.AsMolAcAb = getRoadWorkEffect.AsMolAcAb
		settingRoadWorkEffect.AsMolRdMb = getRoadWorkEffect.AsMolRdMb
		settingRoadWorkEffect.AsRclSnc = getRoadWorkEffect.AsRclSnc
		settingRoadWorkEffect.AsRclIriAfterRecycling = getRoadWorkEffect.AsRclIriAfterRecycling
		settingRoadWorkEffect.AsRclArVb = getRoadWorkEffect.AsRclArVb
		settingRoadWorkEffect.AsRclApoTb = getRoadWorkEffect.AsRclApoTb
		settingRoadWorkEffect.AsRclAcAb = getRoadWorkEffect.AsRclAcAb
		settingRoadWorkEffect.AsRclRdMb = getRoadWorkEffect.AsRclRdMb
		settingRoadWorkEffect.AsRclDefaultHsOld = getRoadWorkEffect.AsRclDefaultHsOld
		settingRoadWorkEffect.AsRcSnc = getRoadWorkEffect.AsRcSnc
		settingRoadWorkEffect.AsRcIriAfterReconstruction = getRoadWorkEffect.AsRcIriAfterReconstruction
		settingRoadWorkEffect.AsRcArVb = getRoadWorkEffect.AsRcArVb
		settingRoadWorkEffect.AsRcApoTb = getRoadWorkEffect.AsRcApoTb
		settingRoadWorkEffect.AsRcAcAb = getRoadWorkEffect.AsRcAcAb
		settingRoadWorkEffect.AsRcRdMb = getRoadWorkEffect.AsRcRdMb
		settingRoadWorkEffect.CcFdrIriAfterFdr = params.CcFdrIriAfterFdr
		settingRoadWorkEffect.CcFdrFaulting = params.CcFdrFaulting
		settingRoadWorkEffect.CcFdrCracking = params.CcFdrCracking
		settingRoadWorkEffect.CcFdrSpalling = params.CcFdrSpalling
		settingRoadWorkEffect.CcBcoIriAfterBco = params.CcBcoIriAfterBco
		settingRoadWorkEffect.CcBcoFaulting = params.CcBcoFaulting
		settingRoadWorkEffect.CcBcoCracking = params.CcBcoCracking
		settingRoadWorkEffect.CcBcoSpalling = params.CcBcoSpalling
		settingRoadWorkEffect.CcMolIriAfterMol = params.CcMolIriAfterMol
		settingRoadWorkEffect.CcMolFaulting = params.CcMolFaulting
		settingRoadWorkEffect.CcMolCracking = params.CcMolCracking
		settingRoadWorkEffect.CcMolSpalling = params.CcMolSpalling
		settingRoadWorkEffect.CcSealIriAfterSeal = params.CcSealIriAfterSeal
		settingRoadWorkEffect.CcSealFaulting = params.CcSealFaulting
		settingRoadWorkEffect.CcSealCracking = params.CcSealCracking
		settingRoadWorkEffect.CcSealSpalling = params.CcSealSpalling
		settingRoadWorkEffect.CcRbcIri = params.CcRbcIri
		settingRoadWorkEffect.CcRbcSlabthk = params.CcRbcSlabthk
		settingRoadWorkEffect.CcRbcPercentFaulting = params.CcRbcPercentFaulting
		settingRoadWorkEffect.CcRbcPercentSpalling = params.CcRbcPercentSpalling
		settingRoadWorkEffect.CcRbcPercentCracking = params.CcRbcPercentCracking
		settingRoadWorkEffect.IsDeleted = false
		settingRoadWorkEffect.IsLatest = true
		settingRoadWorkEffect.CreatedBy = userID
		settingRoadWorkEffect.CreatedAt = dateTime
		settingRoadWorkEffect.UpdatedBy = userID
		settingRoadWorkEffect.UpdatedAt = dateTime
	} else if err != nil {
		log.Println(err)
		return responses.SettingRoadWorkEffectConcrete{}, responses.NewInternalServerError()
	}

	var responsesSettingRoadWorkEffectConcrete responses.SettingRoadWorkEffectConcrete
	responsesSettingRoadWorkEffectConcrete.CcFdrIriAfterFdr = params.CcFdrIriAfterFdr
	responsesSettingRoadWorkEffectConcrete.CcFdrFaulting = params.CcFdrFaulting
	responsesSettingRoadWorkEffectConcrete.CcFdrCracking = params.CcFdrCracking
	responsesSettingRoadWorkEffectConcrete.CcFdrSpalling = params.CcFdrSpalling
	responsesSettingRoadWorkEffectConcrete.CcBcoIriAfterBco = params.CcBcoIriAfterBco
	responsesSettingRoadWorkEffectConcrete.CcBcoFaulting = params.CcBcoFaulting
	responsesSettingRoadWorkEffectConcrete.CcBcoCracking = params.CcBcoCracking
	responsesSettingRoadWorkEffectConcrete.CcBcoSpalling = params.CcBcoSpalling
	responsesSettingRoadWorkEffectConcrete.CcMolIriAfterMol = params.CcMolIriAfterMol
	responsesSettingRoadWorkEffectConcrete.CcMolFaulting = params.CcMolFaulting
	responsesSettingRoadWorkEffectConcrete.CcMolCracking = params.CcMolCracking
	responsesSettingRoadWorkEffectConcrete.CcMolSpalling = params.CcMolSpalling
	responsesSettingRoadWorkEffectConcrete.CcSealIriAfterSeal = params.CcSealIriAfterSeal
	responsesSettingRoadWorkEffectConcrete.CcSealFaulting = params.CcSealFaulting
	responsesSettingRoadWorkEffectConcrete.CcSealCracking = params.CcSealCracking
	responsesSettingRoadWorkEffectConcrete.CcSealSpalling = params.CcSealSpalling
	responsesSettingRoadWorkEffectConcrete.CcRbcIri = params.CcRbcIri
	responsesSettingRoadWorkEffectConcrete.CcRbcSlabthk = params.CcRbcSlabthk
	responsesSettingRoadWorkEffectConcrete.CcRbcPercentFaulting = params.CcRbcPercentFaulting
	responsesSettingRoadWorkEffectConcrete.CcRbcPercentSpalling = params.CcRbcPercentSpalling
	responsesSettingRoadWorkEffectConcrete.CcRbcPercentCracking = params.CcRbcPercentCracking

	err = su.settingRepo.UpdateRoadWorkEffect(&settingRoadWorkEffect)
	if err != nil {
		log.Println(err)
		return responses.SettingRoadWorkEffectConcrete{}, responses.NewInternalServerError()
	}

	err = su.MergeRoadWorkEffectToParams(userID)
	if err != nil {
		log.Println(err)
		return responses.SettingRoadWorkEffectConcrete{}, responses.NewInternalServerError()
	}

	return responsesSettingRoadWorkEffectConcrete, nil
}

func (su *settingUseCase) GetRoadWorkEffectConcrete(c *gin.Context) (responses.SettingRoadWorkEffectConcrete, error) {

	getRoadWorkEffect, err := su.settingRepo.GetRoadWorkEffect()
	if err != nil {
		log.Println(err)
		return responses.SettingRoadWorkEffectConcrete{}, responses.NewInternalServerError()
	}
	var responsesSettingRoadWorkEffectConcrete responses.SettingRoadWorkEffectConcrete
	responsesSettingRoadWorkEffectConcrete.CcFdrIriAfterFdr = getRoadWorkEffect.CcFdrIriAfterFdr
	responsesSettingRoadWorkEffectConcrete.CcFdrFaulting = getRoadWorkEffect.CcFdrFaulting
	responsesSettingRoadWorkEffectConcrete.CcFdrCracking = getRoadWorkEffect.CcFdrCracking
	responsesSettingRoadWorkEffectConcrete.CcFdrSpalling = getRoadWorkEffect.CcFdrSpalling
	responsesSettingRoadWorkEffectConcrete.CcBcoIriAfterBco = getRoadWorkEffect.CcBcoIriAfterBco
	responsesSettingRoadWorkEffectConcrete.CcBcoFaulting = getRoadWorkEffect.CcBcoFaulting
	responsesSettingRoadWorkEffectConcrete.CcBcoCracking = getRoadWorkEffect.CcBcoCracking
	responsesSettingRoadWorkEffectConcrete.CcBcoSpalling = getRoadWorkEffect.CcBcoSpalling
	responsesSettingRoadWorkEffectConcrete.CcMolIriAfterMol = getRoadWorkEffect.CcMolIriAfterMol
	responsesSettingRoadWorkEffectConcrete.CcMolFaulting = getRoadWorkEffect.CcMolFaulting
	responsesSettingRoadWorkEffectConcrete.CcMolCracking = getRoadWorkEffect.CcMolCracking
	responsesSettingRoadWorkEffectConcrete.CcMolSpalling = getRoadWorkEffect.CcMolSpalling
	responsesSettingRoadWorkEffectConcrete.CcSealIriAfterSeal = getRoadWorkEffect.CcSealIriAfterSeal
	responsesSettingRoadWorkEffectConcrete.CcSealFaulting = getRoadWorkEffect.CcSealFaulting
	responsesSettingRoadWorkEffectConcrete.CcSealCracking = getRoadWorkEffect.CcSealCracking
	responsesSettingRoadWorkEffectConcrete.CcSealSpalling = getRoadWorkEffect.CcSealSpalling
	responsesSettingRoadWorkEffectConcrete.CcRbcIri = getRoadWorkEffect.CcRbcIri
	responsesSettingRoadWorkEffectConcrete.CcRbcSlabthk = getRoadWorkEffect.CcRbcSlabthk
	responsesSettingRoadWorkEffectConcrete.CcRbcPercentFaulting = getRoadWorkEffect.CcRbcPercentFaulting
	responsesSettingRoadWorkEffectConcrete.CcRbcPercentSpalling = getRoadWorkEffect.CcRbcPercentSpalling
	responsesSettingRoadWorkEffectConcrete.CcRbcPercentCracking = getRoadWorkEffect.CcRbcPercentCracking

	return responsesSettingRoadWorkEffectConcrete, nil
}

func (su *settingUseCase) MergeRoadWorkEffectToParams(userID int) error {

	getRoadWorkEffect, err := su.settingRepo.GetRoadWorkEffect()
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}

	var getRoadWorkEffectParams models.GetRoadWorkEffectParams
	getRoadWorkEffectParams.Asphalt.AsOlOverlayA0 = getRoadWorkEffect.AsOlOverlayA0
	getRoadWorkEffectParams.Asphalt.AsOlArVb = getRoadWorkEffect.AsOlArVb
	getRoadWorkEffectParams.Asphalt.AsOlPoTb = getRoadWorkEffect.AsOlPoTb
	getRoadWorkEffectParams.Asphalt.AsOlAcAb = getRoadWorkEffect.AsOlAcAb
	getRoadWorkEffectParams.Asphalt.AsOlRdMb = getRoadWorkEffect.AsOlRdMb
	getRoadWorkEffectParams.Asphalt.AsSsRweSsModelA0 = getRoadWorkEffect.AsSsRweSsModelA0
	getRoadWorkEffectParams.Asphalt.AsSsDefaultLowerBoundIriAfterSlurrySeal = getRoadWorkEffect.AsSsDefaultLowerBoundIriAfterSlurrySeal
	getRoadWorkEffectParams.Asphalt.AsSsArVb = getRoadWorkEffect.AsSsArVb
	getRoadWorkEffectParams.Asphalt.AsSsApoTb = getRoadWorkEffect.AsSsApoTb
	getRoadWorkEffectParams.Asphalt.AsSsAcAb = getRoadWorkEffect.AsSsAcAb
	getRoadWorkEffectParams.Asphalt.AsSsRdMb = getRoadWorkEffect.AsSsRdMb
	getRoadWorkEffectParams.Asphalt.AsMolIriAfterMillOverlay = getRoadWorkEffect.AsMolIriAfterMillOverlay
	getRoadWorkEffectParams.Asphalt.AsMolArVb = getRoadWorkEffect.AsMolArVb
	getRoadWorkEffectParams.Asphalt.AsMolApoTb = getRoadWorkEffect.AsMolApoTb
	getRoadWorkEffectParams.Asphalt.AsMolAcAb = getRoadWorkEffect.AsMolAcAb
	getRoadWorkEffectParams.Asphalt.AsMolRdMb = getRoadWorkEffect.AsMolRdMb
	getRoadWorkEffectParams.Asphalt.AsRclSnc = getRoadWorkEffect.AsRclSnc
	getRoadWorkEffectParams.Asphalt.AsRclIriAfterRecycling = getRoadWorkEffect.AsRclIriAfterRecycling
	getRoadWorkEffectParams.Asphalt.AsRclArVb = getRoadWorkEffect.AsRclArVb
	getRoadWorkEffectParams.Asphalt.AsRclApoTb = getRoadWorkEffect.AsRclApoTb
	getRoadWorkEffectParams.Asphalt.AsRclAcAb = getRoadWorkEffect.AsRclAcAb
	getRoadWorkEffectParams.Asphalt.AsRclRdMb = getRoadWorkEffect.AsRclRdMb
	getRoadWorkEffectParams.Asphalt.AsRclDefaultHsOld = getRoadWorkEffect.AsRclDefaultHsOld
	getRoadWorkEffectParams.Asphalt.AsRcSnc = getRoadWorkEffect.AsRcSnc
	getRoadWorkEffectParams.Asphalt.AsRcIriAfterReconstruction = getRoadWorkEffect.AsRcIriAfterReconstruction
	getRoadWorkEffectParams.Asphalt.AsRcArVb = getRoadWorkEffect.AsRcArVb
	getRoadWorkEffectParams.Asphalt.AsRcApoTb = getRoadWorkEffect.AsRcApoTb
	getRoadWorkEffectParams.Asphalt.AsRcAcAb = getRoadWorkEffect.AsRcAcAb
	getRoadWorkEffectParams.Asphalt.AsRcRdMb = getRoadWorkEffect.AsRcRdMb
	getRoadWorkEffectParams.Concrete.CcFdrIriAfterFdr = getRoadWorkEffect.CcFdrIriAfterFdr
	getRoadWorkEffectParams.Concrete.CcFdrFaulting = getRoadWorkEffect.CcFdrFaulting
	getRoadWorkEffectParams.Concrete.CcFdrCracking = getRoadWorkEffect.CcFdrCracking
	getRoadWorkEffectParams.Concrete.CcFdrSpalling = getRoadWorkEffect.CcFdrSpalling
	getRoadWorkEffectParams.Concrete.CcBcoIriAfterBco = getRoadWorkEffect.CcBcoIriAfterBco
	getRoadWorkEffectParams.Concrete.CcBcoFaulting = getRoadWorkEffect.CcBcoFaulting
	getRoadWorkEffectParams.Concrete.CcBcoCracking = getRoadWorkEffect.CcBcoCracking
	getRoadWorkEffectParams.Concrete.CcBcoSpalling = getRoadWorkEffect.CcBcoSpalling
	getRoadWorkEffectParams.Concrete.CcMolIriAfterMol = getRoadWorkEffect.CcMolIriAfterMol
	getRoadWorkEffectParams.Concrete.CcMolFaulting = getRoadWorkEffect.CcMolFaulting
	getRoadWorkEffectParams.Concrete.CcMolCracking = getRoadWorkEffect.CcMolCracking
	getRoadWorkEffectParams.Concrete.CcMolSpalling = getRoadWorkEffect.CcMolSpalling
	getRoadWorkEffectParams.Concrete.CcSealIriAfterSeal = getRoadWorkEffect.CcSealIriAfterSeal
	getRoadWorkEffectParams.Concrete.CcSealFaulting = getRoadWorkEffect.CcSealFaulting
	getRoadWorkEffectParams.Concrete.CcSealCracking = getRoadWorkEffect.CcSealCracking
	getRoadWorkEffectParams.Concrete.CcSealSpalling = getRoadWorkEffect.CcSealSpalling
	getRoadWorkEffectParams.Concrete.CcRbcIri = getRoadWorkEffect.CcRbcIri
	getRoadWorkEffectParams.Concrete.CcRbcSlabthk = getRoadWorkEffect.CcRbcSlabthk
	getRoadWorkEffectParams.Concrete.CcRbcPercentFaulting = getRoadWorkEffect.CcRbcPercentFaulting
	getRoadWorkEffectParams.Concrete.CcRbcPercentSpalling = getRoadWorkEffect.CcRbcPercentSpalling
	getRoadWorkEffectParams.Concrete.CcRbcPercentCracking = getRoadWorkEffect.CcRbcPercentCracking

	err = su.settingRepo.UpdateRoadWorkEffectParamsByIsLatestIsFalse()
	if err != nil {
		log.Println(err)
		return err
	}

	b, err := json.Marshal(getRoadWorkEffectParams)
	if err != nil {
		fmt.Println(err)
		return err
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return err
	}

	var SettingRoadWorkEffectParams models.SettingRoadWorkEffectParams
	SettingRoadWorkEffectParams.Params = buffer.String()
	SettingRoadWorkEffectParams.IsLatest = true
	SettingRoadWorkEffectParams.CreatedBy = userID
	SettingRoadWorkEffectParams.CreatedAt = time.Now()

	err = su.settingRepo.CreateRoadWorkEffectParams(&SettingRoadWorkEffectParams)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (m *settingUseCase) DeleteSettingRefSurfaceByID(ID int) error {
	err := m.settingRepo.DeleteSettingRefSurfaceByID(ID)
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (su *settingUseCase) CreateRoadUserCostAccLossValue(params requests.RoadUserCostAccLossValue, c *gin.Context) (responses.RoadUserCostAccLossValue, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	err := su.settingRepo.UpdateRoadUserCostLossValueParamsByIsLatestIsFalse()
	if err != nil {
		log.Println(err)
		return responses.RoadUserCostAccLossValue{}, err
	}

	helpers.RoundFloatPointer(params.ValueOfFatalAccidents, 2)
	helpers.RoundFloatPointer(params.ValueOfAccidentsWithSeriousInjuries, 2)
	helpers.RoundFloatPointer(params.ValueOfAccidentsWithMinorInjuries, 2)
	helpers.RoundFloatPointer(params.ValueOfAccidentsWithPropertyDamaged, 2)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return responses.RoadUserCostAccLossValue{}, err
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return responses.RoadUserCostAccLossValue{}, err
	}

	var settingRoadUserCostLossValue models.SettingRoadUserCostLossValue
	settingRoadUserCostLossValue.Params = buffer.String()
	settingRoadUserCostLossValue.IsDeleted = false
	settingRoadUserCostLossValue.IsLatest = true
	settingRoadUserCostLossValue.CreatedBy = userID
	settingRoadUserCostLossValue.CreatedAt = dateTime
	settingRoadUserCostLossValue.UpdatedBy = userID
	settingRoadUserCostLossValue.UpdatedAt = dateTime

	err = su.settingRepo.CreateRoadUserCostLossValueParams(&settingRoadUserCostLossValue)
	if err != nil {
		log.Println(err)
		return responses.RoadUserCostAccLossValue{}, responses.NewInternalServerError()
	}

	err = su.MergeRoadUserCostParams(userID)
	if err != nil {
		log.Println(err)
		return responses.RoadUserCostAccLossValue{}, responses.NewInternalServerError()
	}

	var roadUserCostAccLossValue responses.RoadUserCostAccLossValue
	roadUserCostAccLossValue.ValueOfAccidentsWithMinorInjuries = params.ValueOfAccidentsWithMinorInjuries
	roadUserCostAccLossValue.ValueOfAccidentsWithPropertyDamaged = params.ValueOfAccidentsWithPropertyDamaged
	roadUserCostAccLossValue.ValueOfAccidentsWithSeriousInjuries = params.ValueOfAccidentsWithSeriousInjuries
	roadUserCostAccLossValue.ValueOfFatalAccidents = params.ValueOfFatalAccidents

	return roadUserCostAccLossValue, nil
}

func (su *settingUseCase) GetRoadUserCostAccLossValue(c *gin.Context) (responses.RoadUserCostAccLossValue, error) {

	costLossValue, err := su.settingRepo.GetRoadUserCostLossValueParamsByLatest()
	if err != nil {
		log.Println(err)
		return responses.RoadUserCostAccLossValue{}, err
	}
	var roadUserCostAccLossValue responses.RoadUserCostAccLossValue

	err = json.Unmarshal([]byte(costLossValue.Params), &roadUserCostAccLossValue)
	if err != nil {
		fmt.Println(err)
		return responses.RoadUserCostAccLossValue{}, err
	}

	return roadUserCostAccLossValue, nil
}

func (su *settingUseCase) CreateRoadUserCostAccChanceOfAccident(params requests.RoadUserCostAccChanceOfAccident, c *gin.Context) (requests.RoadUserCostAccChanceOfAccident, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	helpers.RoundFloatPointer(params.NumberOfFatalAccidents, 2)
	helpers.RoundFloatPointer(params.NumberOfAccidentsWithSeriousInjuries, 2)
	helpers.RoundFloatPointer(params.NumberOfAccidentsWithMinorInjuries, 2)
	helpers.RoundFloatPointer(params.NumberOfAccidentsWithPropertyDamage, 2)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.RoadUserCostAccChanceOfAccident{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.RoadUserCostAccChanceOfAccident{}, responses.NewInternalServerError()
	}

	var roadUserCostChanceOfAccident models.SettingRoadUserCostChanceOfAccident
	roadUserCostChanceOfAccident.RoadGroupId = params.RoadGroupId
	roadUserCostChanceOfAccident.Params = buffer.String()
	roadUserCostChanceOfAccident.IsLatest = true
	roadUserCostChanceOfAccident.IsDeleted = false
	roadUserCostChanceOfAccident.CreatedAt = dateTime
	roadUserCostChanceOfAccident.CreatedBy = userID
	roadUserCostChanceOfAccident.UpdatedAt = dateTime
	roadUserCostChanceOfAccident.UpdatedBy = userID

	err = su.settingRepo.UpdateRoadUserCostLossValueParamsByIsLatestIsFalseAndRoadGroupId(params.RoadGroupId)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostAccChanceOfAccident{}, responses.NewInternalServerError()
	}

	err = su.settingRepo.CreateChanceOfAccidentParams(&roadUserCostChanceOfAccident)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostAccChanceOfAccident{}, responses.NewInternalServerError()
	}

	err = su.MergeRoadUserCostParams(userID)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostAccChanceOfAccident{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetRoadUserCostAccChanceOfAccident(id string, c *gin.Context) (responses.RoadUserCostAccChanceOfAccident, error) {
	idInteger, err := helpers.ConvertStringToInt(id)
	chanceOfAccident, err := su.settingRepo.GetRoadUserCostAccChanceOfAccidentParamsByLatestAndRoadGroupId(idInteger)
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return responses.RoadUserCostAccChanceOfAccident{}, responses.NewInternalServerError()
	}

	var responsesChanceOfAccident responses.RoadUserCostAccChanceOfAccident

	err = json.Unmarshal([]byte(chanceOfAccident.Params), &responsesChanceOfAccident)
	if err != nil {
		fmt.Println(err)
		return responses.RoadUserCostAccChanceOfAccident{}, err
	}

	return responsesChanceOfAccident, nil
}

func (su *settingUseCase) CreateRoadUserCostRucDefaultData(params requests.RoadUserCostRusDefaultData, c *gin.Context) (requests.RoadUserCostRusDefaultData, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.FuelUCost, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.OilUCost, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.TypeUCost, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.VehUCost, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.MUpper, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.MLower, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Wheels, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.NumPass, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.TUCost, 4)

	helpers.RoundFloatPointer(params.CarOverThanSeven.FuelUCost, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.OilUCost, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.TypeUCost, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.VehUCost, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.MUpper, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.MLower, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Wheels, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.NumPass, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.TUCost, 4)

	helpers.RoundFloatPointer(params.LightBus.FuelUCost, 4)
	helpers.RoundFloatPointer(params.LightBus.OilUCost, 4)
	helpers.RoundFloatPointer(params.LightBus.TypeUCost, 4)
	helpers.RoundFloatPointer(params.LightBus.VehUCost, 4)
	helpers.RoundFloatPointer(params.LightBus.MUpper, 4)
	helpers.RoundFloatPointer(params.LightBus.MLower, 4)
	helpers.RoundFloatPointer(params.LightBus.Wheels, 4)
	helpers.RoundFloatPointer(params.LightBus.NumPass, 4)
	helpers.RoundFloatPointer(params.LightBus.TUCost, 4)

	helpers.RoundFloatPointer(params.MediumBus.FuelUCost, 4)
	helpers.RoundFloatPointer(params.MediumBus.OilUCost, 4)
	helpers.RoundFloatPointer(params.MediumBus.TypeUCost, 4)
	helpers.RoundFloatPointer(params.MediumBus.VehUCost, 4)
	helpers.RoundFloatPointer(params.MediumBus.MUpper, 4)
	helpers.RoundFloatPointer(params.MediumBus.MLower, 4)
	helpers.RoundFloatPointer(params.MediumBus.Wheels, 4)
	helpers.RoundFloatPointer(params.MediumBus.NumPass, 4)
	helpers.RoundFloatPointer(params.MediumBus.TUCost, 4)

	helpers.RoundFloatPointer(params.HeavyBus.FuelUCost, 4)
	helpers.RoundFloatPointer(params.HeavyBus.OilUCost, 4)
	helpers.RoundFloatPointer(params.HeavyBus.TypeUCost, 4)
	helpers.RoundFloatPointer(params.HeavyBus.VehUCost, 4)
	helpers.RoundFloatPointer(params.HeavyBus.MUpper, 4)
	helpers.RoundFloatPointer(params.HeavyBus.MLower, 4)
	helpers.RoundFloatPointer(params.HeavyBus.Wheels, 4)
	helpers.RoundFloatPointer(params.HeavyBus.NumPass, 4)
	helpers.RoundFloatPointer(params.HeavyBus.TUCost, 4)

	helpers.RoundFloatPointer(params.LightTruck.FuelUCost, 4)
	helpers.RoundFloatPointer(params.LightTruck.OilUCost, 4)
	helpers.RoundFloatPointer(params.LightTruck.TypeUCost, 4)
	helpers.RoundFloatPointer(params.LightTruck.VehUCost, 4)
	helpers.RoundFloatPointer(params.LightTruck.MUpper, 4)
	helpers.RoundFloatPointer(params.LightTruck.MLower, 4)
	helpers.RoundFloatPointer(params.LightTruck.Wheels, 4)
	helpers.RoundFloatPointer(params.LightTruck.NumPass, 4)
	helpers.RoundFloatPointer(params.LightTruck.TUCost, 4)

	helpers.RoundFloatPointer(params.MediumTruck.FuelUCost, 4)
	helpers.RoundFloatPointer(params.MediumTruck.OilUCost, 4)
	helpers.RoundFloatPointer(params.MediumTruck.TypeUCost, 4)
	helpers.RoundFloatPointer(params.MediumTruck.VehUCost, 4)
	helpers.RoundFloatPointer(params.MediumTruck.MUpper, 4)
	helpers.RoundFloatPointer(params.MediumTruck.MLower, 4)
	helpers.RoundFloatPointer(params.MediumTruck.Wheels, 4)
	helpers.RoundFloatPointer(params.MediumTruck.NumPass, 4)
	helpers.RoundFloatPointer(params.MediumTruck.TUCost, 4)

	helpers.RoundFloatPointer(params.HeavyTruck.FuelUCost, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.OilUCost, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.TypeUCost, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.VehUCost, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.MUpper, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.MLower, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.Wheels, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.NumPass, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.TUCost, 4)

	helpers.RoundFloatPointer(params.FullTrailor.FuelUCost, 4)
	helpers.RoundFloatPointer(params.FullTrailor.OilUCost, 4)
	helpers.RoundFloatPointer(params.FullTrailor.TypeUCost, 4)
	helpers.RoundFloatPointer(params.FullTrailor.VehUCost, 4)
	helpers.RoundFloatPointer(params.FullTrailor.MUpper, 4)
	helpers.RoundFloatPointer(params.FullTrailor.MLower, 4)
	helpers.RoundFloatPointer(params.FullTrailor.Wheels, 4)
	helpers.RoundFloatPointer(params.FullTrailor.NumPass, 4)
	helpers.RoundFloatPointer(params.FullTrailor.TUCost, 4)

	helpers.RoundFloatPointer(params.SemiTrailor.FuelUCost, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.OilUCost, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.TypeUCost, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.VehUCost, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.MUpper, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.MLower, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.Wheels, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.NumPass, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.TUCost, 4)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.RoadUserCostRusDefaultData{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.RoadUserCostRusDefaultData{}, responses.NewInternalServerError()
	}

	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return requests.RoadUserCostRusDefaultData{}, responses.NewInternalServerError()
	}
	getRoadUserCostRuc.DefaultData = buffer.String()
	getRoadUserCostRuc.IsLatest = true
	getRoadUserCostRuc.IsDeleted = false
	getRoadUserCostRuc.CreatedAt = dateTime
	getRoadUserCostRuc.CreatedBy = userID
	getRoadUserCostRuc.UpdatedAt = dateTime
	getRoadUserCostRuc.UpdatedBy = userID

	err = su.settingRepo.CreateRoadUserCostRuc(getRoadUserCostRuc)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusDefaultData{}, responses.NewInternalServerError()
	}

	err = su.MergeRoadUserCostParams(userID)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusDefaultData{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetRoadUserCostRucDefaultData(c *gin.Context) (responses.RoadUserCostRusDefaultData, error) {
	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil {
		log.Println(err)
		return responses.RoadUserCostRusDefaultData{}, responses.NewInternalServerError()
	}
	var responseRusDefaultData responses.RoadUserCostRusDefaultData

	err = json.Unmarshal([]byte(getRoadUserCostRuc.DefaultData), &responseRusDefaultData)
	if err != nil {
		fmt.Println(err)
		return responses.RoadUserCostRusDefaultData{}, err
	}
	return responseRusDefaultData, nil
}

func (su *settingUseCase) CreateRoadUserCostRucDriving(params requests.RoadUserCostRusDriving, c *gin.Context) (requests.RoadUserCostRusDriving, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.CrA1, 5)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.CrA2, 5)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.P, 5)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Cd, 5)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Cduml, 5)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Ad, 5)

	helpers.RoundFloatPointer(params.CarOverThanSeven.CrA1, 5)
	helpers.RoundFloatPointer(params.CarOverThanSeven.CrA2, 5)
	helpers.RoundFloatPointer(params.CarOverThanSeven.P, 5)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Cd, 5)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Cduml, 5)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Ad, 5)

	helpers.RoundFloatPointer(params.LightBus.CrA1, 5)
	helpers.RoundFloatPointer(params.LightBus.CrA2, 5)
	helpers.RoundFloatPointer(params.LightBus.P, 5)
	helpers.RoundFloatPointer(params.LightBus.Cd, 5)
	helpers.RoundFloatPointer(params.LightBus.Cduml, 5)
	helpers.RoundFloatPointer(params.LightBus.Ad, 5)

	helpers.RoundFloatPointer(params.MediumBus.CrA1, 5)
	helpers.RoundFloatPointer(params.MediumBus.CrA2, 5)
	helpers.RoundFloatPointer(params.MediumBus.P, 5)
	helpers.RoundFloatPointer(params.MediumBus.Cd, 5)
	helpers.RoundFloatPointer(params.MediumBus.Cduml, 5)
	helpers.RoundFloatPointer(params.MediumBus.Ad, 5)

	helpers.RoundFloatPointer(params.HeavyBus.CrA1, 5)
	helpers.RoundFloatPointer(params.HeavyBus.CrA2, 5)
	helpers.RoundFloatPointer(params.HeavyBus.P, 5)
	helpers.RoundFloatPointer(params.HeavyBus.Cd, 5)
	helpers.RoundFloatPointer(params.HeavyBus.Cduml, 5)
	helpers.RoundFloatPointer(params.HeavyBus.Ad, 5)

	helpers.RoundFloatPointer(params.LightTruck.CrA1, 5)
	helpers.RoundFloatPointer(params.LightTruck.CrA2, 5)
	helpers.RoundFloatPointer(params.LightTruck.P, 5)
	helpers.RoundFloatPointer(params.LightTruck.Cd, 5)
	helpers.RoundFloatPointer(params.LightTruck.Cduml, 5)
	helpers.RoundFloatPointer(params.LightTruck.Ad, 5)

	helpers.RoundFloatPointer(params.MediumTruck.CrA1, 5)
	helpers.RoundFloatPointer(params.MediumTruck.CrA2, 5)
	helpers.RoundFloatPointer(params.MediumTruck.P, 5)
	helpers.RoundFloatPointer(params.MediumTruck.Cd, 5)
	helpers.RoundFloatPointer(params.MediumTruck.Cduml, 5)
	helpers.RoundFloatPointer(params.MediumTruck.Ad, 5)

	helpers.RoundFloatPointer(params.HeavyTruck.CrA1, 5)
	helpers.RoundFloatPointer(params.HeavyTruck.CrA2, 5)
	helpers.RoundFloatPointer(params.HeavyTruck.P, 5)
	helpers.RoundFloatPointer(params.HeavyTruck.Cd, 5)
	helpers.RoundFloatPointer(params.HeavyTruck.Cduml, 5)
	helpers.RoundFloatPointer(params.HeavyTruck.Ad, 5)

	helpers.RoundFloatPointer(params.FullTrailor.CrA1, 5)
	helpers.RoundFloatPointer(params.FullTrailor.CrA2, 5)
	helpers.RoundFloatPointer(params.FullTrailor.P, 5)
	helpers.RoundFloatPointer(params.FullTrailor.Cd, 5)
	helpers.RoundFloatPointer(params.FullTrailor.Cduml, 5)
	helpers.RoundFloatPointer(params.FullTrailor.Ad, 5)
	helpers.RoundFloatPointer(params.FullTrailor.Ad, 5)

	helpers.RoundFloatPointer(params.SemiTrailor.CrA1, 5)
	helpers.RoundFloatPointer(params.SemiTrailor.CrA2, 5)
	helpers.RoundFloatPointer(params.SemiTrailor.P, 5)
	helpers.RoundFloatPointer(params.SemiTrailor.Cd, 5)
	helpers.RoundFloatPointer(params.SemiTrailor.Cduml, 5)
	helpers.RoundFloatPointer(params.SemiTrailor.Ad, 5)
	helpers.RoundFloatPointer(params.SemiTrailor.Ad, 5)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.RoadUserCostRusDriving{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.RoadUserCostRusDriving{}, responses.NewInternalServerError()
	}

	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return requests.RoadUserCostRusDriving{}, responses.NewInternalServerError()
	}
	getRoadUserCostRuc.Driving = buffer.String()
	getRoadUserCostRuc.IsLatest = true
	getRoadUserCostRuc.IsDeleted = false
	getRoadUserCostRuc.CreatedAt = dateTime
	getRoadUserCostRuc.CreatedBy = userID
	getRoadUserCostRuc.UpdatedAt = dateTime
	getRoadUserCostRuc.UpdatedBy = userID

	err = su.settingRepo.CreateRoadUserCostRuc(getRoadUserCostRuc)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusDriving{}, responses.NewInternalServerError()
	}

	err = su.MergeRoadUserCostParams(userID)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusDriving{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetRoadUserCostRucDriving(c *gin.Context) (responses.RoadUserCostRusDriving, error) {
	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil {
		log.Println(err)
		return responses.RoadUserCostRusDriving{}, responses.NewInternalServerError()
	}
	var responseRusDriving responses.RoadUserCostRusDriving

	err = json.Unmarshal([]byte(getRoadUserCostRuc.Driving), &responseRusDriving)
	if err != nil {
		fmt.Println(err)
		return responses.RoadUserCostRusDriving{}, err
	}
	return responseRusDriving, nil
}

func (su *settingUseCase) CreateRoadUserCostRucEngineSpeed(params requests.RoadUserCostRusEngineSpeed, c *gin.Context) (requests.RoadUserCostRusEngineSpeed, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.RpmA0, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.RpmA1, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.RpmA2, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.RpmIdle, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Rpm100, 4)

	helpers.RoundFloatPointer(params.CarOverThanSeven.RpmA0, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.RpmA1, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.RpmA2, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.RpmIdle, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Rpm100, 4)

	helpers.RoundFloatPointer(params.LightBus.RpmA0, 4)
	helpers.RoundFloatPointer(params.LightBus.RpmA1, 4)
	helpers.RoundFloatPointer(params.LightBus.RpmA2, 4)
	helpers.RoundFloatPointer(params.LightBus.RpmIdle, 4)
	helpers.RoundFloatPointer(params.LightBus.Rpm100, 4)

	helpers.RoundFloatPointer(params.MediumBus.RpmA0, 4)
	helpers.RoundFloatPointer(params.MediumBus.RpmA1, 4)
	helpers.RoundFloatPointer(params.MediumBus.RpmA2, 4)
	helpers.RoundFloatPointer(params.MediumBus.RpmIdle, 4)
	helpers.RoundFloatPointer(params.MediumBus.Rpm100, 4)

	helpers.RoundFloatPointer(params.HeavyBus.RpmA0, 4)
	helpers.RoundFloatPointer(params.HeavyBus.RpmA1, 4)
	helpers.RoundFloatPointer(params.HeavyBus.RpmA2, 4)
	helpers.RoundFloatPointer(params.HeavyBus.RpmIdle, 4)
	helpers.RoundFloatPointer(params.HeavyBus.Rpm100, 4)

	helpers.RoundFloatPointer(params.LightTruck.RpmA0, 4)
	helpers.RoundFloatPointer(params.LightTruck.RpmA1, 4)
	helpers.RoundFloatPointer(params.LightTruck.RpmA2, 4)
	helpers.RoundFloatPointer(params.LightTruck.RpmIdle, 4)
	helpers.RoundFloatPointer(params.LightTruck.Rpm100, 4)

	helpers.RoundFloatPointer(params.MediumTruck.RpmA0, 4)
	helpers.RoundFloatPointer(params.MediumTruck.RpmA1, 4)
	helpers.RoundFloatPointer(params.MediumTruck.RpmA2, 4)
	helpers.RoundFloatPointer(params.MediumTruck.RpmIdle, 4)
	helpers.RoundFloatPointer(params.MediumTruck.Rpm100, 4)

	helpers.RoundFloatPointer(params.HeavyTruck.RpmA0, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.RpmA1, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.RpmA2, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.RpmIdle, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.Rpm100, 4)

	helpers.RoundFloatPointer(params.FullTrailor.RpmA0, 4)
	helpers.RoundFloatPointer(params.FullTrailor.RpmA1, 4)
	helpers.RoundFloatPointer(params.FullTrailor.RpmA2, 4)
	helpers.RoundFloatPointer(params.FullTrailor.RpmIdle, 4)
	helpers.RoundFloatPointer(params.FullTrailor.Rpm100, 4)

	helpers.RoundFloatPointer(params.SemiTrailor.RpmA0, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.RpmA1, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.RpmA2, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.RpmIdle, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.Rpm100, 4)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.RoadUserCostRusEngineSpeed{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.RoadUserCostRusEngineSpeed{}, responses.NewInternalServerError()
	}

	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return requests.RoadUserCostRusEngineSpeed{}, responses.NewInternalServerError()
	}
	getRoadUserCostRuc.EngineSpeed = buffer.String()
	getRoadUserCostRuc.IsLatest = true
	getRoadUserCostRuc.IsDeleted = false
	getRoadUserCostRuc.CreatedAt = dateTime
	getRoadUserCostRuc.CreatedBy = userID
	getRoadUserCostRuc.UpdatedAt = dateTime
	getRoadUserCostRuc.UpdatedBy = userID

	err = su.settingRepo.CreateRoadUserCostRuc(getRoadUserCostRuc)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusEngineSpeed{}, responses.NewInternalServerError()
	}

	err = su.MergeRoadUserCostParams(userID)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusEngineSpeed{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetRoadUserCostRucEngineSpeed(c *gin.Context) (responses.RoadUserCostRusEngineSpeed, error) {
	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil {
		log.Println(err)
		return responses.RoadUserCostRusEngineSpeed{}, responses.NewInternalServerError()
	}
	var responseRusEngineSpeed responses.RoadUserCostRusEngineSpeed

	err = json.Unmarshal([]byte(getRoadUserCostRuc.EngineSpeed), &responseRusEngineSpeed)
	if err != nil {
		fmt.Println(err)
		return responses.RoadUserCostRusEngineSpeed{}, err
	}
	return responseRusEngineSpeed, nil
}

func (su *settingUseCase) CreateRoadUserCostRucFuelConsumption(params requests.RoadUserCostRusFuelConsumption, c *gin.Context) (requests.RoadUserCostRusFuelConsumption, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.IdleFuel, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.DfFuel, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.ZeTab, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Ehp, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Edt, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Prat, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Kpea, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.PaccsA0, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.PctPeng, 4)

	helpers.RoundFloatPointer(params.CarOverThanSeven.IdleFuel, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.DfFuel, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.ZeTab, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Ehp, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Edt, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Prat, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Kpea, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.PaccsA0, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.PctPeng, 4)

	helpers.RoundFloatPointer(params.LightBus.IdleFuel, 4)
	helpers.RoundFloatPointer(params.LightBus.DfFuel, 4)
	helpers.RoundFloatPointer(params.LightBus.ZeTab, 4)
	helpers.RoundFloatPointer(params.LightBus.Ehp, 4)
	helpers.RoundFloatPointer(params.LightBus.Edt, 4)
	helpers.RoundFloatPointer(params.LightBus.Prat, 4)
	helpers.RoundFloatPointer(params.LightBus.Kpea, 4)
	helpers.RoundFloatPointer(params.LightBus.PaccsA0, 4)
	helpers.RoundFloatPointer(params.LightBus.PctPeng, 4)

	helpers.RoundFloatPointer(params.MediumBus.IdleFuel, 4)
	helpers.RoundFloatPointer(params.MediumBus.DfFuel, 4)
	helpers.RoundFloatPointer(params.MediumBus.ZeTab, 4)
	helpers.RoundFloatPointer(params.MediumBus.Ehp, 4)
	helpers.RoundFloatPointer(params.MediumBus.Edt, 4)
	helpers.RoundFloatPointer(params.MediumBus.Prat, 4)
	helpers.RoundFloatPointer(params.MediumBus.Kpea, 4)
	helpers.RoundFloatPointer(params.MediumBus.PaccsA0, 4)
	helpers.RoundFloatPointer(params.MediumBus.PctPeng, 4)

	helpers.RoundFloatPointer(params.HeavyBus.IdleFuel, 4)
	helpers.RoundFloatPointer(params.HeavyBus.DfFuel, 4)
	helpers.RoundFloatPointer(params.HeavyBus.ZeTab, 4)
	helpers.RoundFloatPointer(params.HeavyBus.Ehp, 4)
	helpers.RoundFloatPointer(params.HeavyBus.Edt, 4)
	helpers.RoundFloatPointer(params.HeavyBus.Prat, 4)
	helpers.RoundFloatPointer(params.HeavyBus.Kpea, 4)
	helpers.RoundFloatPointer(params.HeavyBus.PaccsA0, 4)
	helpers.RoundFloatPointer(params.HeavyBus.PctPeng, 4)

	helpers.RoundFloatPointer(params.LightTruck.IdleFuel, 4)
	helpers.RoundFloatPointer(params.LightTruck.DfFuel, 4)
	helpers.RoundFloatPointer(params.LightTruck.ZeTab, 4)
	helpers.RoundFloatPointer(params.LightTruck.Ehp, 4)
	helpers.RoundFloatPointer(params.LightTruck.Edt, 4)
	helpers.RoundFloatPointer(params.LightTruck.Prat, 4)
	helpers.RoundFloatPointer(params.LightTruck.Kpea, 4)
	helpers.RoundFloatPointer(params.LightTruck.PaccsA0, 4)
	helpers.RoundFloatPointer(params.LightTruck.PctPeng, 4)

	helpers.RoundFloatPointer(params.MediumTruck.IdleFuel, 4)
	helpers.RoundFloatPointer(params.MediumTruck.DfFuel, 4)
	helpers.RoundFloatPointer(params.MediumTruck.ZeTab, 4)
	helpers.RoundFloatPointer(params.MediumTruck.Ehp, 4)
	helpers.RoundFloatPointer(params.MediumTruck.Edt, 4)
	helpers.RoundFloatPointer(params.MediumTruck.Prat, 4)
	helpers.RoundFloatPointer(params.MediumTruck.Kpea, 4)
	helpers.RoundFloatPointer(params.MediumTruck.PaccsA0, 4)
	helpers.RoundFloatPointer(params.MediumTruck.PctPeng, 4)

	helpers.RoundFloatPointer(params.HeavyTruck.IdleFuel, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.DfFuel, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.ZeTab, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.Ehp, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.Edt, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.Prat, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.Kpea, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.PaccsA0, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.PctPeng, 4)

	helpers.RoundFloatPointer(params.FullTrailor.IdleFuel, 4)
	helpers.RoundFloatPointer(params.FullTrailor.DfFuel, 4)
	helpers.RoundFloatPointer(params.FullTrailor.ZeTab, 4)
	helpers.RoundFloatPointer(params.FullTrailor.Ehp, 4)
	helpers.RoundFloatPointer(params.FullTrailor.Edt, 4)
	helpers.RoundFloatPointer(params.FullTrailor.Prat, 4)
	helpers.RoundFloatPointer(params.FullTrailor.Kpea, 4)
	helpers.RoundFloatPointer(params.FullTrailor.PaccsA0, 4)
	helpers.RoundFloatPointer(params.FullTrailor.PctPeng, 4)

	helpers.RoundFloatPointer(params.SemiTrailor.IdleFuel, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.DfFuel, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.ZeTab, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.Ehp, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.Edt, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.Prat, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.Kpea, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.PaccsA0, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.PctPeng, 4)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.RoadUserCostRusFuelConsumption{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.RoadUserCostRusFuelConsumption{}, responses.NewInternalServerError()
	}

	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return requests.RoadUserCostRusFuelConsumption{}, responses.NewInternalServerError()
	}
	getRoadUserCostRuc.FuelConsumption = buffer.String()
	getRoadUserCostRuc.IsLatest = true
	getRoadUserCostRuc.IsDeleted = false
	getRoadUserCostRuc.CreatedAt = dateTime
	getRoadUserCostRuc.CreatedBy = userID
	getRoadUserCostRuc.UpdatedAt = dateTime
	getRoadUserCostRuc.UpdatedBy = userID

	err = su.settingRepo.CreateRoadUserCostRuc(getRoadUserCostRuc)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusFuelConsumption{}, responses.NewInternalServerError()
	}

	err = su.MergeRoadUserCostParams(userID)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusFuelConsumption{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetRoadUserCostRucFuelConsumption(c *gin.Context) (responses.RoadUserCostRusFuelConsumption, error) {
	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil {
		log.Println(err)
		return responses.RoadUserCostRusFuelConsumption{}, responses.NewInternalServerError()
	}
	var responseRusFuelConsumption responses.RoadUserCostRusFuelConsumption

	err = json.Unmarshal([]byte(getRoadUserCostRuc.FuelConsumption), &responseRusFuelConsumption)
	if err != nil {
		fmt.Println(err)
		return responses.RoadUserCostRusFuelConsumption{}, err
	}
	return responseRusFuelConsumption, nil
}

func (su *settingUseCase) CreateRoadUserCostRucLubricantConsumption(params requests.RoadUserCostRusLubricantConsumption, c *gin.Context) (requests.RoadUserCostRusLubricantConsumption, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.OilCont, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.OilOper, 4)

	helpers.RoundFloatPointer(params.CarOverThanSeven.OilCont, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.OilOper, 4)

	helpers.RoundFloatPointer(params.LightBus.OilCont, 4)
	helpers.RoundFloatPointer(params.LightBus.OilOper, 4)

	helpers.RoundFloatPointer(params.MediumBus.OilCont, 4)
	helpers.RoundFloatPointer(params.MediumBus.OilOper, 4)

	helpers.RoundFloatPointer(params.HeavyBus.OilCont, 4)
	helpers.RoundFloatPointer(params.HeavyBus.OilOper, 4)

	helpers.RoundFloatPointer(params.LightTruck.OilCont, 4)
	helpers.RoundFloatPointer(params.LightTruck.OilOper, 4)

	helpers.RoundFloatPointer(params.MediumTruck.OilCont, 4)
	helpers.RoundFloatPointer(params.MediumTruck.OilOper, 4)

	helpers.RoundFloatPointer(params.HeavyTruck.OilCont, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.OilOper, 4)

	helpers.RoundFloatPointer(params.FullTrailor.OilCont, 4)
	helpers.RoundFloatPointer(params.FullTrailor.OilOper, 4)

	helpers.RoundFloatPointer(params.SemiTrailor.OilCont, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.OilOper, 4)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.RoadUserCostRusLubricantConsumption{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.RoadUserCostRusLubricantConsumption{}, responses.NewInternalServerError()
	}

	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return requests.RoadUserCostRusLubricantConsumption{}, responses.NewInternalServerError()
	}
	getRoadUserCostRuc.LubricantConsumption = buffer.String()
	getRoadUserCostRuc.IsLatest = true
	getRoadUserCostRuc.IsDeleted = false
	getRoadUserCostRuc.CreatedAt = dateTime
	getRoadUserCostRuc.CreatedBy = userID
	getRoadUserCostRuc.UpdatedAt = dateTime
	getRoadUserCostRuc.UpdatedBy = userID

	err = su.settingRepo.CreateRoadUserCostRuc(getRoadUserCostRuc)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusLubricantConsumption{}, responses.NewInternalServerError()
	}

	err = su.MergeRoadUserCostParams(userID)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusLubricantConsumption{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetRoadUserCostRucLubricantConsumption(c *gin.Context) (responses.RoadUserCostRusLubricantConsumption, error) {
	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil {
		log.Println(err)
		return responses.RoadUserCostRusLubricantConsumption{}, responses.NewInternalServerError()
	}
	var responseRusLubricantConsumption responses.RoadUserCostRusLubricantConsumption

	err = json.Unmarshal([]byte(getRoadUserCostRuc.LubricantConsumption), &responseRusLubricantConsumption)
	if err != nil {
		fmt.Println(err)
		return responses.RoadUserCostRusLubricantConsumption{}, err
	}
	return responseRusLubricantConsumption, nil
}

func (su *settingUseCase) CreateRoadUserCostRucWasteOfConsumption(params requests.RoadUserCostRusWasteOfConsumption, c *gin.Context) (requests.RoadUserCostRusWasteOfConsumption, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.C0Tc, 5)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.CtCte, 5)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Vol, 5)

	helpers.RoundFloatPointer(params.CarOverThanSeven.C0Tc, 5)
	helpers.RoundFloatPointer(params.CarOverThanSeven.CtCte, 5)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Vol, 5)

	helpers.RoundFloatPointer(params.LightBus.C0Tc, 5)
	helpers.RoundFloatPointer(params.LightBus.CtCte, 5)
	helpers.RoundFloatPointer(params.LightBus.Vol, 5)

	helpers.RoundFloatPointer(params.MediumBus.C0Tc, 5)
	helpers.RoundFloatPointer(params.MediumBus.CtCte, 5)
	helpers.RoundFloatPointer(params.MediumBus.Vol, 5)

	helpers.RoundFloatPointer(params.HeavyBus.C0Tc, 5)
	helpers.RoundFloatPointer(params.HeavyBus.CtCte, 5)
	helpers.RoundFloatPointer(params.HeavyBus.Vol, 5)

	helpers.RoundFloatPointer(params.LightTruck.C0Tc, 5)
	helpers.RoundFloatPointer(params.LightTruck.CtCte, 5)
	helpers.RoundFloatPointer(params.LightTruck.Vol, 5)

	helpers.RoundFloatPointer(params.MediumTruck.C0Tc, 5)
	helpers.RoundFloatPointer(params.MediumTruck.CtCte, 5)
	helpers.RoundFloatPointer(params.MediumTruck.Vol, 5)

	helpers.RoundFloatPointer(params.HeavyTruck.C0Tc, 5)
	helpers.RoundFloatPointer(params.HeavyTruck.CtCte, 5)
	helpers.RoundFloatPointer(params.HeavyTruck.Vol, 5)

	helpers.RoundFloatPointer(params.FullTrailor.C0Tc, 5)
	helpers.RoundFloatPointer(params.FullTrailor.CtCte, 5)
	helpers.RoundFloatPointer(params.FullTrailor.Vol, 5)

	helpers.RoundFloatPointer(params.SemiTrailor.C0Tc, 5)
	helpers.RoundFloatPointer(params.SemiTrailor.CtCte, 5)
	helpers.RoundFloatPointer(params.SemiTrailor.Vol, 5)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.RoadUserCostRusWasteOfConsumption{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.RoadUserCostRusWasteOfConsumption{}, responses.NewInternalServerError()
	}

	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return requests.RoadUserCostRusWasteOfConsumption{}, responses.NewInternalServerError()
	}
	getRoadUserCostRuc.WasteOfConsumption = buffer.String()
	getRoadUserCostRuc.IsLatest = true
	getRoadUserCostRuc.IsDeleted = false
	getRoadUserCostRuc.CreatedAt = dateTime
	getRoadUserCostRuc.CreatedBy = userID
	getRoadUserCostRuc.UpdatedAt = dateTime
	getRoadUserCostRuc.UpdatedBy = userID

	err = su.settingRepo.CreateRoadUserCostRuc(getRoadUserCostRuc)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusWasteOfConsumption{}, responses.NewInternalServerError()
	}

	err = su.MergeRoadUserCostParams(userID)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusWasteOfConsumption{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetRoadUserCostRucWasteOfConsumption(c *gin.Context) (responses.RoadUserCostRusWasteOfConsumption, error) {
	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil {
		log.Println(err)
		return responses.RoadUserCostRusWasteOfConsumption{}, responses.NewInternalServerError()
	}
	var responseRusWasteOfConsumption responses.RoadUserCostRusWasteOfConsumption

	err = json.Unmarshal([]byte(getRoadUserCostRuc.WasteOfConsumption), &responseRusWasteOfConsumption)
	if err != nil {
		fmt.Println(err)
		return responses.RoadUserCostRusWasteOfConsumption{}, err
	}
	return responseRusWasteOfConsumption, nil
}

func (su *settingUseCase) CreateRoadUserCostRucMaintenance(params requests.RoadUserCostRusMaintenance, c *gin.Context) (requests.RoadUserCostRusMaintenance, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Kpc, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Akmo, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Life0, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Kp, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.A0, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.A1, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.CpCon, 4)

	helpers.RoundFloatPointer(params.CarOverThanSeven.Kpc, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Akmo, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Life0, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Kp, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.A0, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.A1, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.CpCon, 4)

	helpers.RoundFloatPointer(params.LightBus.Kpc, 4)
	helpers.RoundFloatPointer(params.LightBus.Akmo, 4)
	helpers.RoundFloatPointer(params.LightBus.Life0, 4)
	helpers.RoundFloatPointer(params.LightBus.Kp, 4)
	helpers.RoundFloatPointer(params.LightBus.A0, 4)
	helpers.RoundFloatPointer(params.LightBus.A1, 4)
	helpers.RoundFloatPointer(params.LightBus.CpCon, 4)

	helpers.RoundFloatPointer(params.MediumBus.Kpc, 4)
	helpers.RoundFloatPointer(params.MediumBus.Akmo, 4)
	helpers.RoundFloatPointer(params.MediumBus.Life0, 4)
	helpers.RoundFloatPointer(params.MediumBus.Kp, 4)
	helpers.RoundFloatPointer(params.MediumBus.A0, 4)
	helpers.RoundFloatPointer(params.MediumBus.A1, 4)
	helpers.RoundFloatPointer(params.MediumBus.CpCon, 4)

	helpers.RoundFloatPointer(params.HeavyBus.Kpc, 4)
	helpers.RoundFloatPointer(params.HeavyBus.Akmo, 4)
	helpers.RoundFloatPointer(params.HeavyBus.Life0, 4)
	helpers.RoundFloatPointer(params.HeavyBus.Kp, 4)
	helpers.RoundFloatPointer(params.HeavyBus.A0, 4)
	helpers.RoundFloatPointer(params.HeavyBus.A1, 4)
	helpers.RoundFloatPointer(params.HeavyBus.CpCon, 4)

	helpers.RoundFloatPointer(params.LightTruck.Kpc, 4)
	helpers.RoundFloatPointer(params.LightTruck.Akmo, 4)
	helpers.RoundFloatPointer(params.LightTruck.Life0, 4)
	helpers.RoundFloatPointer(params.LightTruck.Kp, 4)
	helpers.RoundFloatPointer(params.LightTruck.A0, 4)
	helpers.RoundFloatPointer(params.LightTruck.A1, 4)
	helpers.RoundFloatPointer(params.LightTruck.CpCon, 4)

	helpers.RoundFloatPointer(params.MediumTruck.Kpc, 4)
	helpers.RoundFloatPointer(params.MediumTruck.Akmo, 4)
	helpers.RoundFloatPointer(params.MediumTruck.Life0, 4)
	helpers.RoundFloatPointer(params.MediumTruck.Kp, 4)
	helpers.RoundFloatPointer(params.MediumTruck.A0, 4)
	helpers.RoundFloatPointer(params.MediumTruck.A1, 4)
	helpers.RoundFloatPointer(params.MediumTruck.CpCon, 4)

	helpers.RoundFloatPointer(params.HeavyTruck.Kpc, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.Akmo, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.Life0, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.Kp, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.A0, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.A1, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.CpCon, 4)

	helpers.RoundFloatPointer(params.FullTrailor.Kpc, 4)
	helpers.RoundFloatPointer(params.FullTrailor.Akmo, 4)
	helpers.RoundFloatPointer(params.FullTrailor.Life0, 4)
	helpers.RoundFloatPointer(params.FullTrailor.Kp, 4)
	helpers.RoundFloatPointer(params.FullTrailor.A0, 4)
	helpers.RoundFloatPointer(params.FullTrailor.A1, 4)
	helpers.RoundFloatPointer(params.FullTrailor.CpCon, 4)

	helpers.RoundFloatPointer(params.SemiTrailor.Kpc, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.Akmo, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.Life0, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.Kp, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.A0, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.A1, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.CpCon, 4)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.RoadUserCostRusMaintenance{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.RoadUserCostRusMaintenance{}, responses.NewInternalServerError()
	}

	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return requests.RoadUserCostRusMaintenance{}, responses.NewInternalServerError()
	}
	getRoadUserCostRuc.Maintenance = buffer.String()
	getRoadUserCostRuc.IsLatest = true
	getRoadUserCostRuc.IsDeleted = false
	getRoadUserCostRuc.CreatedAt = dateTime
	getRoadUserCostRuc.CreatedBy = userID
	getRoadUserCostRuc.UpdatedAt = dateTime
	getRoadUserCostRuc.UpdatedBy = userID

	err = su.settingRepo.CreateRoadUserCostRuc(getRoadUserCostRuc)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusMaintenance{}, responses.NewInternalServerError()
	}

	err = su.MergeRoadUserCostParams(userID)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusMaintenance{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetRoadUserCostRucMaintenance(c *gin.Context) (responses.RoadUserCostRusMaintenance, error) {
	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil {
		log.Println(err)
		return responses.RoadUserCostRusMaintenance{}, responses.NewInternalServerError()
	}
	var responseRusMaintenance responses.RoadUserCostRusMaintenance

	err = json.Unmarshal([]byte(getRoadUserCostRuc.Maintenance), &responseRusMaintenance)
	if err != nil {
		fmt.Println(err)
		return responses.RoadUserCostRusMaintenance{}, err
	}
	return responseRusMaintenance, nil
}

func (su *settingUseCase) CreateRoadUserCostRucTravelTime(params requests.RoadUserCostRusTravelTime, c *gin.Context) (requests.RoadUserCostRusTravelTime, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.PcTwk, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.PcTwk, 4)
	helpers.RoundFloatPointer(params.LightBus.PcTwk, 4)
	helpers.RoundFloatPointer(params.MediumBus.PcTwk, 4)
	helpers.RoundFloatPointer(params.HeavyBus.PcTwk, 4)
	helpers.RoundFloatPointer(params.LightTruck.PcTwk, 4)
	helpers.RoundFloatPointer(params.MediumTruck.PcTwk, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.PcTwk, 4)
	helpers.RoundFloatPointer(params.FullTrailor.PcTwk, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.PcTwk, 4)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.RoadUserCostRusTravelTime{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.RoadUserCostRusTravelTime{}, responses.NewInternalServerError()
	}

	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return requests.RoadUserCostRusTravelTime{}, responses.NewInternalServerError()
	}
	getRoadUserCostRuc.TravelTime = buffer.String()
	getRoadUserCostRuc.IsLatest = true
	getRoadUserCostRuc.IsDeleted = false
	getRoadUserCostRuc.CreatedAt = dateTime
	getRoadUserCostRuc.CreatedBy = userID
	getRoadUserCostRuc.UpdatedAt = dateTime
	getRoadUserCostRuc.UpdatedBy = userID

	err = su.settingRepo.CreateRoadUserCostRuc(getRoadUserCostRuc)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusTravelTime{}, responses.NewInternalServerError()
	}

	err = su.MergeRoadUserCostParams(userID)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusTravelTime{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetRoadUserCostRucTravelTime(c *gin.Context) (responses.RoadUserCostRusTravelTime, error) {
	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil {
		log.Println(err)
		return responses.RoadUserCostRusTravelTime{}, responses.NewInternalServerError()
	}
	var responseRusTravelTime responses.RoadUserCostRusTravelTime

	err = json.Unmarshal([]byte(getRoadUserCostRuc.TravelTime), &responseRusTravelTime)
	if err != nil {
		fmt.Println(err)
		return responses.RoadUserCostRusTravelTime{}, err
	}
	return responseRusTravelTime, nil
}

func (su *settingUseCase) CreateRoadUserCostRucVehicleSpeedCalculation(params requests.RoadUserCostRusVehicleSpeedCalculation, c *gin.Context) (requests.RoadUserCostRusVehicleSpeedCalculation, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Cw1, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Cw2, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Cw3, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.A2, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.A3, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Pd, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.Pb, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.ArvMax, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.AUpper0, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.ALower0, 4)
	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.A1, 4)

	helpers.RoundFloatPointer(params.CarOverThanSeven.Cw1, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Cw2, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Cw3, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.A2, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.A3, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Pd, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.Pb, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.ArvMax, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.AUpper0, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.ALower0, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.A1, 4)

	helpers.RoundFloatPointer(params.LightBus.Cw1, 4)
	helpers.RoundFloatPointer(params.LightBus.Cw2, 4)
	helpers.RoundFloatPointer(params.LightBus.Cw3, 4)
	helpers.RoundFloatPointer(params.LightBus.A2, 4)
	helpers.RoundFloatPointer(params.LightBus.A3, 4)
	helpers.RoundFloatPointer(params.LightBus.Pd, 4)
	helpers.RoundFloatPointer(params.LightBus.Pb, 4)
	helpers.RoundFloatPointer(params.LightBus.ArvMax, 4)
	helpers.RoundFloatPointer(params.LightBus.AUpper0, 4)
	helpers.RoundFloatPointer(params.LightBus.ALower0, 4)
	helpers.RoundFloatPointer(params.LightBus.A1, 4)

	helpers.RoundFloatPointer(params.MediumBus.Cw1, 4)
	helpers.RoundFloatPointer(params.MediumBus.Cw2, 4)
	helpers.RoundFloatPointer(params.MediumBus.Cw3, 4)
	helpers.RoundFloatPointer(params.MediumBus.A2, 4)
	helpers.RoundFloatPointer(params.MediumBus.A3, 4)
	helpers.RoundFloatPointer(params.MediumBus.Pd, 4)
	helpers.RoundFloatPointer(params.MediumBus.Pb, 4)
	helpers.RoundFloatPointer(params.MediumBus.ArvMax, 4)
	helpers.RoundFloatPointer(params.MediumBus.AUpper0, 4)
	helpers.RoundFloatPointer(params.MediumBus.ALower0, 4)
	helpers.RoundFloatPointer(params.MediumBus.A1, 4)

	helpers.RoundFloatPointer(params.HeavyBus.Cw1, 4)
	helpers.RoundFloatPointer(params.HeavyBus.Cw2, 4)
	helpers.RoundFloatPointer(params.HeavyBus.Cw3, 4)
	helpers.RoundFloatPointer(params.HeavyBus.A2, 4)
	helpers.RoundFloatPointer(params.HeavyBus.A3, 4)
	helpers.RoundFloatPointer(params.HeavyBus.Pd, 4)
	helpers.RoundFloatPointer(params.HeavyBus.Pb, 4)
	helpers.RoundFloatPointer(params.HeavyBus.ArvMax, 4)
	helpers.RoundFloatPointer(params.HeavyBus.AUpper0, 4)
	helpers.RoundFloatPointer(params.HeavyBus.ALower0, 4)
	helpers.RoundFloatPointer(params.HeavyBus.A1, 4)

	helpers.RoundFloatPointer(params.LightTruck.Cw1, 4)
	helpers.RoundFloatPointer(params.LightTruck.Cw2, 4)
	helpers.RoundFloatPointer(params.LightTruck.Cw3, 4)
	helpers.RoundFloatPointer(params.LightTruck.A2, 4)
	helpers.RoundFloatPointer(params.LightTruck.A3, 4)
	helpers.RoundFloatPointer(params.LightTruck.Pd, 4)
	helpers.RoundFloatPointer(params.LightTruck.Pb, 4)
	helpers.RoundFloatPointer(params.LightTruck.ArvMax, 4)
	helpers.RoundFloatPointer(params.LightTruck.AUpper0, 4)
	helpers.RoundFloatPointer(params.LightTruck.ALower0, 4)
	helpers.RoundFloatPointer(params.LightTruck.A1, 4)

	helpers.RoundFloatPointer(params.MediumTruck.Cw1, 4)
	helpers.RoundFloatPointer(params.MediumTruck.Cw2, 4)
	helpers.RoundFloatPointer(params.MediumTruck.Cw3, 4)
	helpers.RoundFloatPointer(params.MediumTruck.A2, 4)
	helpers.RoundFloatPointer(params.MediumTruck.A3, 4)
	helpers.RoundFloatPointer(params.MediumTruck.Pd, 4)
	helpers.RoundFloatPointer(params.MediumTruck.Pb, 4)
	helpers.RoundFloatPointer(params.MediumTruck.ArvMax, 4)
	helpers.RoundFloatPointer(params.MediumTruck.AUpper0, 4)
	helpers.RoundFloatPointer(params.MediumTruck.ALower0, 4)
	helpers.RoundFloatPointer(params.MediumTruck.A1, 4)

	helpers.RoundFloatPointer(params.HeavyTruck.Cw1, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.Cw2, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.Cw3, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.A2, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.A3, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.Pd, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.Pb, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.ArvMax, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.AUpper0, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.ALower0, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.A1, 4)

	helpers.RoundFloatPointer(params.FullTrailor.Cw1, 4)
	helpers.RoundFloatPointer(params.FullTrailor.Cw2, 4)
	helpers.RoundFloatPointer(params.FullTrailor.Cw3, 4)
	helpers.RoundFloatPointer(params.FullTrailor.A2, 4)
	helpers.RoundFloatPointer(params.FullTrailor.A3, 4)
	helpers.RoundFloatPointer(params.FullTrailor.Pd, 4)
	helpers.RoundFloatPointer(params.FullTrailor.Pb, 4)
	helpers.RoundFloatPointer(params.FullTrailor.ArvMax, 4)
	helpers.RoundFloatPointer(params.FullTrailor.AUpper0, 4)
	helpers.RoundFloatPointer(params.FullTrailor.ALower0, 4)
	helpers.RoundFloatPointer(params.FullTrailor.A1, 4)

	helpers.RoundFloatPointer(params.SemiTrailor.Cw1, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.Cw2, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.Cw3, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.A2, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.A3, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.Pd, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.Pb, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.ArvMax, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.AUpper0, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.ALower0, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.A1, 4)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.RoadUserCostRusVehicleSpeedCalculation{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.RoadUserCostRusVehicleSpeedCalculation{}, responses.NewInternalServerError()
	}

	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return requests.RoadUserCostRusVehicleSpeedCalculation{}, responses.NewInternalServerError()
	}
	getRoadUserCostRuc.VehicleSpeedCalculation = buffer.String()
	getRoadUserCostRuc.IsLatest = true
	getRoadUserCostRuc.IsDeleted = false
	getRoadUserCostRuc.CreatedAt = dateTime
	getRoadUserCostRuc.CreatedBy = userID
	getRoadUserCostRuc.UpdatedAt = dateTime
	getRoadUserCostRuc.UpdatedBy = userID

	err = su.settingRepo.CreateRoadUserCostRuc(getRoadUserCostRuc)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusVehicleSpeedCalculation{}, responses.NewInternalServerError()
	}

	err = su.MergeRoadUserCostParams(userID)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusVehicleSpeedCalculation{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetRoadUserCostRucVehicleSpeedCalculation(c *gin.Context) (responses.RoadUserCostRusVehicleSpeedCalculation, error) {
	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil {
		log.Println(err)
		return responses.RoadUserCostRusVehicleSpeedCalculation{}, responses.NewInternalServerError()
	}
	var responseRusVehicleSpeedCalculation responses.RoadUserCostRusVehicleSpeedCalculation

	err = json.Unmarshal([]byte(getRoadUserCostRuc.VehicleSpeedCalculation), &responseRusVehicleSpeedCalculation)
	if err != nil {
		fmt.Println(err)
		return responses.RoadUserCostRusVehicleSpeedCalculation{}, err
	}
	return responseRusVehicleSpeedCalculation, nil
}

func (su *settingUseCase) CreateRoadUserCostRucTrafficData(params requests.RoadUserCostRusTrafficData, c *gin.Context) (requests.RoadUserCostRusTrafficData, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	helpers.RoundFloatPointer(params.CarLessThanEqualSeven.PcuEquivalent, 4)
	helpers.RoundFloatPointer(params.CarOverThanSeven.PcuEquivalent, 4)
	helpers.RoundFloatPointer(params.LightBus.PcuEquivalent, 4)
	helpers.RoundFloatPointer(params.MediumBus.PcuEquivalent, 4)
	helpers.RoundFloatPointer(params.HeavyBus.PcuEquivalent, 4)
	helpers.RoundFloatPointer(params.LightTruck.PcuEquivalent, 4)
	helpers.RoundFloatPointer(params.MediumTruck.PcuEquivalent, 4)
	helpers.RoundFloatPointer(params.HeavyTruck.PcuEquivalent, 4)
	helpers.RoundFloatPointer(params.FullTrailor.PcuEquivalent, 4)
	helpers.RoundFloatPointer(params.SemiTrailor.PcuEquivalent, 4)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.RoadUserCostRusTrafficData{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.RoadUserCostRusTrafficData{}, responses.NewInternalServerError()
	}

	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return requests.RoadUserCostRusTrafficData{}, responses.NewInternalServerError()
	}
	getRoadUserCostRuc.TrafficData = buffer.String()
	getRoadUserCostRuc.IsLatest = true
	getRoadUserCostRuc.IsDeleted = false
	getRoadUserCostRuc.CreatedAt = dateTime
	getRoadUserCostRuc.CreatedBy = userID
	getRoadUserCostRuc.UpdatedAt = dateTime
	getRoadUserCostRuc.UpdatedBy = userID

	err = su.settingRepo.CreateRoadUserCostRuc(getRoadUserCostRuc)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusTrafficData{}, responses.NewInternalServerError()
	}

	err = su.MergeRoadUserCostParams(userID)
	if err != nil {
		log.Println(err)
		return requests.RoadUserCostRusTrafficData{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetRoadUserCostRucTrafficData(c *gin.Context) (responses.RoadUserCostRusTrafficData, error) {
	getRoadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil {
		log.Println(err)
		return responses.RoadUserCostRusTrafficData{}, responses.NewInternalServerError()
	}
	var responseRusTrafficData responses.RoadUserCostRusTrafficData

	err = json.Unmarshal([]byte(getRoadUserCostRuc.TrafficData), &responseRusTrafficData)
	if err != nil {
		fmt.Println(err)
		return responses.RoadUserCostRusTrafficData{}, err
	}
	return responseRusTrafficData, nil
}

func (su *settingUseCase) MergeRoadUserCostParams(userID int) error {
	dateTime := time.Now()

	roadUserCostRuc, err := su.settingRepo.GetRoadUserCostRuc()
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return err
	}

	costLossValue, err := su.settingRepo.GetRoadUserCostLossValueParamsByLatest()
	if err != nil {
		log.Println(err)
		return err
	}

	chanceOfAccident, err := su.settingRepo.GetRoadUserCostAccChanceOfAccidentParamsByLatest()
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return err
	}

	err = su.settingRepo.UpdateRoadUserCostParamsByIsLatestIsFalse()
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return err
	}

	var mergeRoadUserCost models.MergeRoadUserCost

	var defaultData models.RoadUserCostRusDefaultData
	var driving models.RoadUserCostRusDriving
	var engineSpeed models.RoadUserCostRusEngineSpeed
	var fuelConsumption models.RoadUserCostRusFuelConsumption
	var lubricantConsumption models.RoadUserCostRusLubricantConsumption
	var wasteOfConsumption models.RoadUserCostRusWasteOfConsumption
	var maintenance models.RoadUserCostRusMaintenance
	var travelTime models.RoadUserCostRusTravelTime
	var vehicleSpeedCalculation models.RoadUserCostRusVehicleSpeedCalculation
	var trafficData models.RoadUserCostRusTrafficData

	json.Unmarshal([]byte(roadUserCostRuc.DefaultData), &defaultData)
	mergeRoadUserCost.Ruc.DefaultData = defaultData
	json.Unmarshal([]byte(roadUserCostRuc.Driving), &driving)
	mergeRoadUserCost.Ruc.Driving = driving
	json.Unmarshal([]byte(roadUserCostRuc.EngineSpeed), &engineSpeed)
	mergeRoadUserCost.Ruc.EngineSpeed = engineSpeed
	json.Unmarshal([]byte(roadUserCostRuc.FuelConsumption), &fuelConsumption)
	mergeRoadUserCost.Ruc.FuelConsumption = fuelConsumption
	json.Unmarshal([]byte(roadUserCostRuc.LubricantConsumption), &lubricantConsumption)
	mergeRoadUserCost.Ruc.LubricantConsumption = lubricantConsumption
	json.Unmarshal([]byte(roadUserCostRuc.WasteOfConsumption), &wasteOfConsumption)
	mergeRoadUserCost.Ruc.WasteOfConsumption = wasteOfConsumption
	json.Unmarshal([]byte(roadUserCostRuc.Maintenance), &maintenance)
	mergeRoadUserCost.Ruc.Maintenance = maintenance
	json.Unmarshal([]byte(roadUserCostRuc.TravelTime), &travelTime)
	mergeRoadUserCost.Ruc.TravelTime = travelTime
	json.Unmarshal([]byte(roadUserCostRuc.VehicleSpeedCalculation), &vehicleSpeedCalculation)
	mergeRoadUserCost.Ruc.VehicleSpeedCalculation = vehicleSpeedCalculation
	json.Unmarshal([]byte(roadUserCostRuc.TrafficData), &trafficData)
	mergeRoadUserCost.Ruc.TrafficData = trafficData

	var roadUserCostAccLossValue models.RoadUserCostAccLossValue

	err = json.Unmarshal([]byte(costLossValue.Params), &roadUserCostAccLossValue)
	if err != nil {
		fmt.Println(err)
		return err
	}

	mergeRoadUserCost.Acc.LossValue = roadUserCostAccLossValue
	for _, item := range chanceOfAccident {
		var roadUserCostAccChanceOfAccident models.RoadUserCostAccChanceOfAccident
		err = json.Unmarshal([]byte(item.Params), &roadUserCostAccChanceOfAccident)
		if err != nil {
			fmt.Println(err)
			return err
		}

		mergeRoadUserCost.Acc.ChanceOfAccident = append(mergeRoadUserCost.Acc.ChanceOfAccident, roadUserCostAccChanceOfAccident)
	}

	b, err := json.Marshal(mergeRoadUserCost)
	if err != nil {
		fmt.Println(err)
		return err
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return err
	}

	var roadUserCostParams models.SettingRoadUserCostParams
	roadUserCostParams.IsDeleted = false
	roadUserCostParams.IsLatest = true
	roadUserCostParams.Params = buffer.String()
	roadUserCostParams.UpdatedAt = dateTime
	roadUserCostParams.UpdatedBy = userID
	roadUserCostParams.CreatedAt = dateTime
	roadUserCostParams.CreatedBy = userID

	err = su.settingRepo.CreateRoadUserCostParams(&roadUserCostParams)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (su *settingUseCase) CreateOptimization(params requests.Optimization, c *gin.Context) (requests.Optimization, error) {
	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	helpers.RoundFloatPointer(params.BcRatioConstraint, 2)
	helpers.RoundFloatPointer(params.DefaultDesignLife, 2)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.Optimization{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.Optimization{}, responses.NewInternalServerError()
	}

	err = su.settingRepo.UpdateOptimizationByIsLatestIsFalse()
	if err != nil {
		log.Println(err)
		return requests.Optimization{}, responses.NewInternalServerError()
	}

	var settingOptimization models.SettingOptimization
	settingOptimization.IsDeleted = false
	settingOptimization.IsLatest = true
	settingOptimization.Params = buffer.String()
	settingOptimization.UpdatedAt = dateTime
	settingOptimization.UpdatedBy = userID
	settingOptimization.CreatedAt = dateTime
	settingOptimization.CreatedBy = userID

	err = su.settingRepo.CreateOptimizationParams(&settingOptimization)
	if err != nil {
		log.Println(err)
		return requests.Optimization{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetOptimization(c *gin.Context) (responses.Optimization, error) {

	getSettingOptimization, err := su.settingRepo.GetOptimizationParams()
	if err != nil {
		if err.Error() == "record not found" {
			return responses.Optimization{}, responses.NewNotFoundError()
		}
		return responses.Optimization{}, responses.NewInternalServerError()
	}

	var optimization models.Optimization
	err = json.Unmarshal([]byte(getSettingOptimization.Params), &optimization)
	if err != nil {
		fmt.Println(err)
		return responses.Optimization{}, responses.NewInternalServerError()
	}

	var responsesOptimization responses.Optimization
	responsesOptimization.BcRatioConstraint = optimization.BcRatioConstraint
	responsesOptimization.DefaultDesignLife = optimization.DefaultDesignLife

	return responsesOptimization, nil
}

func (su *settingUseCase) CreateDeteriorationAsphalt(params requests.DeteriorationAsphalt, c *gin.Context) (requests.DeteriorationAsphalt, error) {

	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	getDeterioration, err := su.settingRepo.GetDeteriorationByRoadGroupId(params.RoadGroupId)
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return requests.DeteriorationAsphalt{}, responses.NewInternalServerError()
	}

	helpers.RoundFloat(params.Tlf, 2)
	helpers.RoundFloat(params.Cdb, 2)
	helpers.RoundFloat(params.Cds, 2)
	helpers.RoundFloat(params.Comp, 2)
	helpers.RoundFloat(params.Kvi, 2)
	helpers.RoundFloat(params.Kvp, 2)
	helpers.RoundFloat(params.Kpi, 2)
	helpers.RoundFloat(params.Kpp, 2)
	helpers.RoundFloat(params.Krid, 2)
	helpers.RoundFloat(params.Krst, 2)
	helpers.RoundFloat(params.Krpd, 2)
	helpers.RoundFloat(params.Cmod, 2)
	helpers.RoundFloat(params.Kgm, 2)
	helpers.RoundFloat(params.Kgp, 2)
	helpers.RoundFloat(params.Kcia, 2)
	helpers.RoundFloat(params.Kciw, 2)
	helpers.RoundFloat(params.Kcpa, 2)
	helpers.RoundFloat(params.Kcpw, 2)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.DeteriorationAsphalt{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.DeteriorationAsphalt{}, responses.NewInternalServerError()
	}

	getDeterioration.RoadGroupId = params.RoadGroupId
	getDeterioration.UpdatedAt = dateTime
	getDeterioration.UpdatedBy = userID
	getDeterioration.CreatedAt = dateTime
	getDeterioration.CreatedBy = userID
	getDeterioration.IsDeleted = false
	getDeterioration.IsLatest = true
	getDeterioration.ParamsAsphalt = buffer.String()

	err = su.settingRepo.UpdateDeterioration(&getDeterioration)
	if err != nil {
		log.Println(err)
		return requests.DeteriorationAsphalt{}, responses.NewInternalServerError()
	}

	err = su.MergeDeteriorationParams(userID)
	if err != nil {
		log.Println(err)
		return requests.DeteriorationAsphalt{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetDeteriorationAsphalt(roadGroupId string, c *gin.Context) (responses.DeteriorationAsphalt, error) {
	idInteger, err := helpers.ConvertStringToInt(roadGroupId)
	var responsesDeteriorationAsphalt responses.DeteriorationAsphalt
	responsesDeteriorationAsphalt.RoadGroupId = idInteger
	getDeterioration, err := su.settingRepo.GetDeteriorationByRoadGroupId(idInteger)
	if err != nil {
		if err.Error() == "record not found" {
			return responsesDeteriorationAsphalt, nil
		}
		return responses.DeteriorationAsphalt{}, responses.NewInternalServerError()
	}

	if getDeterioration.ParamsAsphalt == "" {
		return responsesDeteriorationAsphalt, nil
	}

	var deteriorationAsphalt models.DeteriorationAsphalt
	err = json.Unmarshal([]byte(getDeterioration.ParamsAsphalt), &deteriorationAsphalt)
	if err != nil {
		fmt.Println(err)
		return responses.DeteriorationAsphalt{}, responses.NewInternalServerError()
	}

	responsesDeteriorationAsphalt.Tlf = deteriorationAsphalt.Tlf
	responsesDeteriorationAsphalt.Cdb = deteriorationAsphalt.Cdb
	responsesDeteriorationAsphalt.Cds = deteriorationAsphalt.Cds
	responsesDeteriorationAsphalt.Comp = deteriorationAsphalt.Comp
	responsesDeteriorationAsphalt.Kvi = deteriorationAsphalt.Kvi
	responsesDeteriorationAsphalt.Kvp = deteriorationAsphalt.Kvp
	responsesDeteriorationAsphalt.Kpi = deteriorationAsphalt.Kpi
	responsesDeteriorationAsphalt.Kpp = deteriorationAsphalt.Kpp
	responsesDeteriorationAsphalt.Krid = deteriorationAsphalt.Krid
	responsesDeteriorationAsphalt.Krst = deteriorationAsphalt.Krst
	responsesDeteriorationAsphalt.Krpd = deteriorationAsphalt.Krpd
	responsesDeteriorationAsphalt.Kgm = deteriorationAsphalt.Kgm
	responsesDeteriorationAsphalt.Kgp = deteriorationAsphalt.Kgp
	responsesDeteriorationAsphalt.Kcia = deteriorationAsphalt.Kcia
	responsesDeteriorationAsphalt.Cmod = deteriorationAsphalt.Cmod
	responsesDeteriorationAsphalt.Kciw = deteriorationAsphalt.Kciw
	responsesDeteriorationAsphalt.Kcpa = deteriorationAsphalt.Kcpa
	responsesDeteriorationAsphalt.Kcpw = deteriorationAsphalt.Kcpw

	return responsesDeteriorationAsphalt, nil
}

func (su *settingUseCase) CreateDeteriorationConcrete(params requests.DeteriorationConcrete, c *gin.Context) (requests.DeteriorationConcrete, error) {

	user_id, _ := c.Get("userID")
	userID := int(user_id.(float64))
	dateTime := time.Now()

	getDeterioration, err := su.settingRepo.GetDeteriorationByRoadGroupId(params.RoadGroupId)
	if err != nil && err.Error() != "record not found" {
		log.Println(err)
		return requests.DeteriorationConcrete{}, responses.NewInternalServerError()
	}

	helpers.RoundFloat(params.PSteel, 2)
	helpers.RoundFloat(params.Ec, 2)
	helpers.RoundFloat(params.Mi, 2)
	helpers.RoundFloat(params.Fi, 2)
	helpers.RoundFloat(params.Kjrc, 2)
	helpers.RoundFloat(params.BStress, 2)
	helpers.RoundFloat(params.JtSpace, 2)
	helpers.RoundFloat(params.Kjrf, 2)
	helpers.RoundFloat(params.Widened, 2)
	helpers.RoundFloat(params.PredSeal, 2)
	helpers.RoundFloat(params.DwlCor, 2)
	helpers.RoundFloat(params.Kjrs, 2)
	helpers.RoundFloat(params.Kjrr, 2)

	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return requests.DeteriorationConcrete{}, responses.NewInternalServerError()
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return requests.DeteriorationConcrete{}, responses.NewInternalServerError()
	}

	getDeterioration.RoadGroupId = params.RoadGroupId
	getDeterioration.UpdatedAt = dateTime
	getDeterioration.UpdatedBy = userID
	getDeterioration.CreatedAt = dateTime
	getDeterioration.CreatedBy = userID
	getDeterioration.IsDeleted = false
	getDeterioration.IsLatest = true
	getDeterioration.ParamsConcrete = buffer.String()

	err = su.settingRepo.UpdateDeterioration(&getDeterioration)
	if err != nil {
		log.Println(err)
		return requests.DeteriorationConcrete{}, responses.NewInternalServerError()
	}

	err = su.MergeDeteriorationParams(userID)
	if err != nil {
		log.Println(err)
		return requests.DeteriorationConcrete{}, responses.NewInternalServerError()
	}

	return params, nil
}

func (su *settingUseCase) GetDeteriorationConcrete(roadGroupId string, c *gin.Context) (responses.DeteriorationConcrete, error) {
	idInteger, err := helpers.ConvertStringToInt(roadGroupId)
	var responsesDeteriorationConcrete responses.DeteriorationConcrete
	responsesDeteriorationConcrete.RoadGroupId = idInteger
	getDeterioration, err := su.settingRepo.GetDeteriorationByRoadGroupId(idInteger)
	if err != nil {
		if err.Error() == "record not found" {

			return responsesDeteriorationConcrete, nil
		}
		return responses.DeteriorationConcrete{}, responses.NewInternalServerError()
	}

	if getDeterioration.ParamsConcrete == "" {
		return responsesDeteriorationConcrete, nil
	}

	var deteriorationConcrete models.DeteriorationConcrete
	err = json.Unmarshal([]byte(getDeterioration.ParamsConcrete), &deteriorationConcrete)
	if err != nil {
		fmt.Println(err)
		return responses.DeteriorationConcrete{}, responses.NewInternalServerError()
	}

	responsesDeteriorationConcrete.PSteel = deteriorationConcrete.PSteel
	responsesDeteriorationConcrete.Ec = deteriorationConcrete.Ec
	responsesDeteriorationConcrete.Mi = deteriorationConcrete.Mi
	responsesDeteriorationConcrete.Fi = deteriorationConcrete.Fi
	responsesDeteriorationConcrete.Kjrc = deteriorationConcrete.Kjrc
	responsesDeteriorationConcrete.BStress = deteriorationConcrete.BStress
	responsesDeteriorationConcrete.JtSpace = deteriorationConcrete.JtSpace
	responsesDeteriorationConcrete.Kjrf = deteriorationConcrete.Kjrf
	responsesDeteriorationConcrete.Widened = deteriorationConcrete.Widened
	responsesDeteriorationConcrete.PredSeal = deteriorationConcrete.PredSeal
	responsesDeteriorationConcrete.DwlCor = deteriorationConcrete.DwlCor
	responsesDeteriorationConcrete.Kjrs = deteriorationConcrete.Kjrs
	responsesDeteriorationConcrete.Kjrr = deteriorationConcrete.Kjrr

	return responsesDeteriorationConcrete, nil
}

func (su *settingUseCase) MergeDeteriorationParams(userID int) error {

	getDeterioration, err := su.settingRepo.GetDeteriorationList()
	if err != nil {
		if err.Error() == "record not found" {
			return err
		}
		return err
	}

	var deteriorationParams models.DeteriorationParams
	var deteriorationAsphaltList []models.DeteriorationAsphalt
	var deteriorationConcreteList []models.DeteriorationConcrete

	for _, item := range getDeterioration {
		if item.ParamsAsphalt != "" {
			var deteriorationAsphalt models.DeteriorationAsphalt
			json.Unmarshal([]byte(item.ParamsAsphalt), &deteriorationAsphalt)
			deteriorationAsphaltList = append(deteriorationAsphaltList, deteriorationAsphalt)
		}

		if item.ParamsConcrete != "" {
			var deteriorationConcrete models.DeteriorationConcrete
			json.Unmarshal([]byte(item.ParamsConcrete), &deteriorationConcrete)
			deteriorationConcreteList = append(deteriorationConcreteList, deteriorationConcrete)
		}
	}

	deteriorationParams.Asphalt = deteriorationAsphaltList

	deteriorationParams.Concrete = deteriorationConcreteList

	b, err := json.Marshal(deteriorationParams)
	if err != nil {
		fmt.Println(err)
		return err
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		log.Println(err)
		return err
	}
	err = su.settingRepo.UpdateDeteriorationParamsByIsLatestIsFalse()
	if err != nil {
		log.Println(err)
		return err
	}

	var settingDeteriorationParams models.SettingDeteriorationParams
	settingDeteriorationParams.Params = buffer.String()
	settingDeteriorationParams.IsLatest = true
	settingDeteriorationParams.UpdatedBy = userID
	settingDeteriorationParams.UpdatedAt = time.Now()
	settingDeteriorationParams.CreatedBy = userID
	settingDeteriorationParams.CreatedAt = time.Now()

	err = su.settingRepo.CreateDeteriorationParams(&settingDeteriorationParams)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (u *settingUseCase) GetHris() (interface{}, error) {
	hris, err := u.settingRepo.GetHris()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var hrisResponses []responses.RefHris
	for _, item := range hris {
		var hrisResponse responses.RefHris
		hrisResponse.Id = item.Id
		hrisResponse.RoadNumber = item.RoadNumber
		hrisResponse.OfficeOfHighwaysCode = item.OfficeOfHighwaysCode
		hrisResponse.SectionRoadNumber = item.SectionRoadNumber
		hrisResponse.Status = item.Status

		hrisResponses = append(hrisResponses, hrisResponse)
	}

	return hrisResponses, nil
}

func (u *settingUseCase) GetHrisById(id string) (interface{}, error) {

	integerId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	hris, err := u.settingRepo.GetHrisById(integerId)
	if err != nil {
		return nil, err
	}

	if (hris == models.RefHris{}) {
		return nil, responses.NewNotFoundError()
	}

	var hrisResponse responses.RefHris
	hrisResponse.Id = hris.Id
	hrisResponse.RoadNumber = hris.RoadNumber
	hrisResponse.OfficeOfHighwaysCode = hris.OfficeOfHighwaysCode
	hrisResponse.SectionRoadNumber = hris.SectionRoadNumber
	hrisResponse.Status = hris.Status

	return hrisResponse, nil
}

func (u *settingUseCase) GetHrisPreview() (interface{}, error) {
	refHirs, err := u.settingRepo.GetHrisByStatus()
	if err != nil {
		return nil, err
	}

	sectionGeoms, err := u.SectionGeomLogic(refHirs)
	if err != nil {
		return nil, err
	}

	roadLatest, err := u.RoadLatestLogic(refHirs)
	if err != nil {
		return nil, err
	}

	var hrisSectionGeoms []responses.HrisSectionGeom
	for _, item := range sectionGeoms {
		var hrisSectionGeom responses.HrisSectionGeom
		hrisSectionGeom.RoadGroupNumber = item.RoadNumber
		hrisSectionGeom.SectionRoadNumber = item.SectionRoadNumber
		hrisSectionGeom.SectionRoadThName = item.SectionRoadThName
		hrisSectionGeom.SectionRoadEngName = item.SectionRoadEngName
		hrisSectionGeom.KmStart = item.KmBegin
		hrisSectionGeom.KmEnd = item.KmEnd

		hrisSectionGeoms = append(hrisSectionGeoms, hrisSectionGeom)
	}

	var hrisRoadLists []responses.HrisRoadList
	for _, item := range roadLatest {
		var hrisRoadList responses.HrisRoadList
		hrisRoadList.Name = item.RoadName
		hrisRoadList.Number = item.RoadCode

		hrisRoadLists = append(hrisRoadLists, hrisRoadList)
	}

	var refHrisPreview responses.RefHrisPreview
	refHrisPreview.RoadGroup = hrisRoadLists
	refHrisPreview.RoadSection = hrisSectionGeoms

	return refHrisPreview, nil
}

func (u *settingUseCase) SectionGeomLogic(refHirs []models.RefHris) ([]models.Item, error) {
	var level3 bson.A
	for _, item := range refHirs {
		var level1 bson.A

		if item.OfficeOfHighwaysCode != "" {
			level1 = append(level1, bson.M{"office_of_highways_code": item.OfficeOfHighwaysCode})
		}

		if item.RoadNumber != "" {
			level1 = append(level1, bson.M{"road_number": item.RoadNumber})
		}

		if item.SectionRoadNumber != "" {
			level1 = append(level1, bson.M{"section_road_number": item.SectionRoadNumber})
		}

		if len(level1) > 0 {
			level2 := bson.D{
				{"$and", level1}}

			level3 = append(level3, level2)
		}

	}

	filter := bson.D{
		{"$or", level3},
	}

	sectionGeom, err := u.settingRepo.GetSectionGeomWithFilter(filter)
	if err != nil {
		return nil, err
	}

	return sectionGeom, nil
}

func (u *settingUseCase) RoadLatestLogic(refHirs []models.RefHris) ([]models.RoadLatest, error) {

	var level3 bson.A
	for _, item := range refHirs {
		var level1 bson.A

		if item.RoadNumber != "" {
			level1 = append(level1, bson.M{"road_code": item.RoadNumber})

			level2 := bson.D{
				{"$and", level1}}

			level3 = append(level3, level2)

		}

	}

	filter := bson.D{
		{"$or", level3},
	}

	roadLatest, err := u.settingRepo.GetRoadLatest(filter)
	if err != nil {
		return nil, err
	}

	return roadLatest, nil
}

func (u *settingUseCase) CreateHris(request requests.CreateRefHris, userId int) (interface{}, error) {
	var refHris models.RefHris
	refHris.RoadNumber = request.RoadNumber
	refHris.OfficeOfHighwaysCode = request.OfficeOfHighwaysCode
	refHris.SectionRoadNumber = request.SectionRoadNumber
	refHris.Status = request.Status
	refHris.IsDeleted = false
	refHris.CreatedBy = userId
	refHris.CreatedAt = time.Now()

	err := u.settingRepo.CreateHris(refHris)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	hris, err := u.settingRepo.GetHrisByStatusBylatest()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var hrisResponse responses.RefHris
	hrisResponse.Id = hris.Id
	hrisResponse.RoadNumber = hris.RoadNumber
	hrisResponse.OfficeOfHighwaysCode = hris.OfficeOfHighwaysCode
	hrisResponse.SectionRoadNumber = hris.SectionRoadNumber
	hrisResponse.Status = hris.Status

	return hrisResponse, nil
}

func (u *settingUseCase) UpdateHris(request requests.UpdateRefHris, id string, userId int) (interface{}, error) {

	integerId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	hris, err := u.settingRepo.GetHrisById(integerId)
	if err != nil {
		return nil, err
	}

	if (hris == models.RefHris{}) {
		return nil, responses.NewNotFoundError()
	}

	hris.RoadNumber = request.RoadNumber
	hris.OfficeOfHighwaysCode = request.OfficeOfHighwaysCode
	hris.SectionRoadNumber = request.SectionRoadNumber
	hris.Status = request.Status
	hris.IsDeleted = false
	hris.UpdatedBy = userId
	hris.UpdatedAt = time.Now()

	err = u.settingRepo.UpdateHris(hris)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var hrisResponse responses.RefHris
	hrisResponse.Id = hris.Id
	hrisResponse.RoadNumber = hris.RoadNumber
	hrisResponse.OfficeOfHighwaysCode = hris.OfficeOfHighwaysCode
	hrisResponse.SectionRoadNumber = hris.SectionRoadNumber
	hrisResponse.Status = hris.Status

	return hrisResponse, nil
}

func (u *settingUseCase) DeleteHris(id string, userId int) (interface{}, error) {

	integerId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	hris, err := u.settingRepo.GetHrisById(integerId)
	if err != nil {
		return nil, err
	}

	if (hris == models.RefHris{}) {
		return nil, responses.NewNotFoundError()
	}

	hris.IsDeleted = true
	hris.UpdatedBy = userId
	hris.UpdatedAt = time.Now()

	err = u.settingRepo.UpdateHris(hris)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return hris, nil
}

func (u *settingUseCase) ImportHris() (interface{}, error) {

	tx := u.settingRepo.StartTransSection()

	err := u.GetRoadLatest()
	if err != nil {
		u.settingRepo.RollBack(tx)
		return nil, err
	}

	time.Sleep(3 * time.Second)

	err = u.GetSectionGeom()
	if err != nil {
		u.settingRepo.RollBack(tx)
		return nil, err
	}

	refHirs, err := u.settingRepo.GetRefHirsWithXt(tx)
	if err != nil {
		u.settingRepo.RollBack(tx)
		return nil, err
	}

	sectionGeoms, err := u.SectionGeomLogic(refHirs)
	if err != nil {
		u.settingRepo.RollBack(tx)
		return nil, err
	}

	roadLatest, err := u.RoadLatestLogic(refHirs)
	if err != nil {
		u.settingRepo.RollBack(tx)
		return nil, err
	}

	roadGroups, err := u.settingRepo.GetRoadGroupWithXt(tx)
	if err != nil {
		u.settingRepo.RollBack(tx)
		return nil, err
	}

	roadSections, err := u.settingRepo.GetRoadSectionWithXt(tx)
	if err != nil {
		u.settingRepo.RollBack(tx)
		return nil, err
	}

	err = u.MatchDataModule1(sectionGeoms, roadGroups, roadLatest, tx)
	if err != nil {
		u.settingRepo.RollBack(tx)
		return nil, err
	}

	roadGroups, err = u.settingRepo.GetRoadGroupWithXt(tx)
	if err != nil {
		u.settingRepo.RollBack(tx)
		return nil, err
	}

	err = u.MatchDataModule2(sectionGeoms, roadGroups, roadSections, tx)
	if err != nil {
		u.settingRepo.RollBack(tx)
		return nil, err
	}

	u.settingRepo.Commit(tx)

	var hrisSectionGeoms []responses.HrisSectionGeom
	for _, item := range sectionGeoms {
		var hrisSectionGeom responses.HrisSectionGeom
		hrisSectionGeom.RoadGroupNumber = item.RoadNumber
		hrisSectionGeom.SectionRoadNumber = item.SectionRoadNumber
		hrisSectionGeom.SectionRoadThName = item.SectionRoadThName
		hrisSectionGeom.SectionRoadEngName = item.SectionRoadEngName
		hrisSectionGeom.KmStart = item.KmBegin
		hrisSectionGeom.KmEnd = item.KmEnd

		hrisSectionGeoms = append(hrisSectionGeoms, hrisSectionGeom)
	}

	var hrisRoadLists []responses.HrisRoadList
	for _, item := range roadLatest {
		var hrisRoadList responses.HrisRoadList
		hrisRoadList.Name = item.RoadName
		hrisRoadList.Number = item.RoadCode

		hrisRoadLists = append(hrisRoadLists, hrisRoadList)
	}

	var refHrisPreview responses.RefHrisPreview
	refHrisPreview.RoadGroup = hrisRoadLists
	refHrisPreview.RoadSection = hrisSectionGeoms

	return refHrisPreview, nil
}

func (u *settingUseCase) MatchDataModule1(sectionGeoms []models.Item, roadGroups []models.RoadGroup, roadLatest []models.RoadLatest, xt *gorm.DB) error {

	isDuplicate := map[string]bool{}
	sectionGeomRoadCodes := []string{}
	for _, item := range sectionGeoms {
		if !isDuplicate[item.RoadNumber] {
			sectionGeomRoadCodes = append(sectionGeomRoadCodes, item.RoadNumber)
			isDuplicate[item.RoadNumber] = true
		}
	}

	roadGroupMap := map[string]models.RoadGroup{}
	for _, item := range roadGroups {
		roadGroupMap[item.Number] = item
	}

	roadLatestMap := map[string]models.RoadLatest{}
	for _, item := range roadLatest {
		roadLatestMap[item.RoadCode] = item
	}

	for _, item := range sectionGeomRoadCodes {

		roadName := ""
		roadNames := strings.Split(roadLatestMap[item].RoadName, " ")
		if len(roadNames) >= 3 {
			roadName = roadNames[0] + " " + roadNames[1] + " " + roadNames[2]
		} else {
			roadName = roadNames[0]
		}

		kmStart, err := u.KmToM(roadLatestMap[item].KmStart)
		if err != nil {
			kmStart = 0
		}

		kmEnd, err := u.KmToM(roadLatestMap[item].KmEnd)
		if err != nil {
			kmEnd = 0
		}

		length, err := strconv.ParseFloat(roadLatestMap[item].Length, 32)
		if err != nil {
			length = 0
		}

		if (roadGroupMap[item] == models.RoadGroup{}) {
			var insert models.InsertRoadGroup
			insert.Number = roadLatestMap[item].RoadCode
			insert.Name = roadLatestMap[item].RoadName
			insert.ShortName = roadName
			insert.KmStart = float32(kmStart)
			insert.KmEnd = float32(kmEnd)
			insert.Distance = float32(length)
			err = u.settingRepo.InsertRoadGroup(insert, xt)
			if err != nil {
				return err
			}

		} else {
			if roadGroupMap[item].Name != roadLatestMap[item].RoadName {
				err := u.settingRepo.UpdateRoadGroupNameByNumberWithXt(roadLatestMap[item].RoadCode, roadLatestMap[item].RoadName, roadName, xt)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (u *settingUseCase) MatchDataModule2(sectionGeoms []models.Item, roadGroups []models.RoadGroup, roadSections []models.RoadSection, xt *gorm.DB) error {

	roadSectionsMap := map[string]models.RoadSection{}
	for _, item := range roadSections {
		roadSectionsMap[strconv.Itoa(item.RoadGroupId)+"-"+item.Number] = item
	}

	roadGroupNumberMap := map[string]int{}
	for _, item := range roadGroups {
		roadGroupNumberMap[item.Number] = item.Id
	}

	for _, item := range sectionGeoms {
		nameOriginTH := ""
		nameDestinationTH := ""
		nameOriginEn := ""
		nameDestinationEn := ""

		nameTh := strings.Split(item.SectionRoadThName, " - ")
		if len(nameTh) == 1 {
			nameOriginTH = item.SectionRoadThName
			nameDestinationTH = item.SectionRoadThName
		} else {
			nameOriginTH = nameTh[0]
			nameDestinationTH = nameTh[1]
		}

		nameEn := strings.Split(item.SectionRoadEngName, " - ")
		if len(nameEn) == 1 {
			nameOriginEn = item.SectionRoadEngName
			nameDestinationEn = item.SectionRoadEngName
		} else {
			nameOriginEn = nameEn[0]
			nameDestinationEn = nameEn[1]
		}

		provinceCord, err := strconv.Atoi(item.ProvinceCord)
		if err != nil {
			provinceCord = 0
		}

		kmEnd, err := u.KmToM(item.KmEnd)
		if err != nil {
			kmEnd = 0
		}

		kmStart, err := u.KmToM(item.KmBegin)
		if err != nil {
			kmStart = 0
		}

		length, err := strconv.ParseFloat(item.Distance, 32)
		if err != nil {
			length = 0
		}

		if roadSectionsMap[strconv.Itoa(roadGroupNumberMap[item.RoadNumber])+"-"+item.SectionRoadNumber].Id == 0 {
			var insert models.InsertRoadSection

			sectionPartID, err := strconv.Atoi(item.SectionPartID)
			if err != nil {
				return err
			}

			insert.Id = sectionPartID
			insert.NameOriginTH = nameOriginTH
			insert.NameDestinationTH = nameDestinationTH
			insert.NameOriginEn = nameOriginEn
			insert.NameDestinationEn = nameDestinationEn
			insert.KmStart = float32(kmStart)
			insert.KmEnd = float32(kmEnd)
			insert.RoadGroupId = roadGroupNumberMap[item.RoadNumber]
			insert.Number = item.SectionRoadNumber
			insert.Distance = float32(length)
			insert.RefDivisionCode = item.OfficeOfHighwaysCode
			insert.RefDistrictCode = item.HighwayDistrictCode
			insert.RefDepotCode = item.DepotCode
			insert.ProvinceCode = append(insert.ProvinceCode, int64(provinceCord))

			err = u.settingRepo.InsertRoadSectionWithXt(insert, xt)
			if err != nil {
				return err
			}

		} else {
			var update models.InsertRoadSection

			sectionPartID, err := strconv.Atoi(item.SectionPartID)
			if err != nil {
				return err
			}

			update.Id = sectionPartID
			update.NameOriginTH = nameOriginTH
			update.NameDestinationTH = nameDestinationTH
			update.NameOriginEn = nameOriginEn
			update.NameDestinationEn = nameDestinationEn
			update.RoadGroupId = roadGroupNumberMap[item.RoadNumber]
			update.Number = item.SectionRoadNumber
			update.KmStart = float32(kmStart)
			update.KmEnd = float32(kmEnd)
			update.Distance = float32(length)
			update.RefDivisionCode = item.OfficeOfHighwaysCode
			update.RefDistrictCode = item.HighwayDistrictCode
			update.RefDepotCode = item.DepotCode
			update.ProvinceCode = append(update.ProvinceCode, int64(provinceCord))

			err = u.settingRepo.InsertRoadSectionWithXt(update, xt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *settingUseCase) KmToM(value string) (int, error) {

	stringKm := strings.ReplaceAll(strings.ReplaceAll(value, ".", ""), "+", "")

	stringInt, err := strconv.Atoi(stringKm)
	if err != nil {
		return 0, err
	}

	return stringInt, nil
}

func (u *settingUseCase) GetSectionGeom() error {

	run := 1
	limit := 10
	var err error
	for run <= limit {
		errLogic := u.LogicSectionGeom()
		if err == nil {
			return nil
		} else {
			err = errLogic
		}

		run++
		time.Sleep(3 * time.Second)
	}

	return err
}

func (u *settingUseCase) LogicSectionGeom() error {
	var GetSectionGeom models.SectionGeom

	url := os.Getenv("HRIS2_SECTION_GEOM_URL")
	method := "GET"

	response, err := u.Request(url, method)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(response, &GetSectionGeom)
	if err != nil {
		return err
	}

	var insertData []interface{}
	for _, item := range GetSectionGeom.Item {
		item.IsLatested = true
		insertData = append(insertData, item)
	}

	err = u.settingRepo.InsertSectionGeom(insertData)
	if err != nil {
		return err
	}

	return nil
}

func (u *settingUseCase) GetRoadLatest() error {

	run := 1
	limit := 10
	var err error
	for run <= limit {
		errLogic := u.LogicRoadLatest()
		if err == nil {
			return nil
		} else {
			err = errLogic
		}

		run++
		time.Sleep(3 * time.Second)
	}

	return err
}

func (u *settingUseCase) LogicRoadLatest() error {
	var roadLatest []models.RoadLatest

	url := os.Getenv("HRIS2_ROAD_LATEST_URL") + "?date=2014-08-31"
	method := "GET"

	response, err := u.Request(url, method)
	if err != nil {
		return err
	}

	err = json.Unmarshal(response, &roadLatest)
	if err != nil {
		return err
	}

	var insertData []interface{}
	for _, item := range roadLatest {
		item.IsLatested = true
		if item.Status == "A" {
			insertData = append(insertData, item)
		}
	}

	err = u.settingRepo.InsertRoadLatest(insertData)
	if err != nil {
		return err
	}

	return nil
}

func (u *settingUseCase) Request(url, method string) ([]byte, error) {

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (u *settingUseCase) GetAllHsms(req requests.FilterHsms) ([]responses.HsmsAll, error) {

	assetName := req.AssetName
	isAssetName := false
	if assetName == "" {
		isAssetName = true
	}

	var hsmsAll []responses.HsmsAll

	roadGroup, err := u.settingRepo.GetRoadGroup()
	if err != nil {
		return hsmsAll, err
	}

	roadGroupNumbeMap := map[string]models.RoadGroup{}
	for _, item := range roadGroup {
		roadGroupNumbeMap[item.Number] = item
	}

	roadSection, err := u.settingRepo.GetRoadSection()
	if err != nil {
		return hsmsAll, err
	}

	roadSectionNumbeMap := map[string]models.RoadSection{}
	for _, item := range roadSection {
		roadSectionNumbeMap[item.Number] = item
	}

	refAssetTable, err := u.settingRepo.GetRefAssetTable()
	if err != nil {
		return hsmsAll, err
	}

	refAssetTableTableNameMap := map[string]models.RefAssetTable{}
	for _, item := range refAssetTable {
		refAssetTableTableNameMap[item.TableNameColumn] = item
	}
	/////////////////////////////////////////////////////////////

	bridgeTable := u.settingRepo.GetTableFromStruct(models.HsmsMotorwayFootbridge{})

	bridgeTableType := u.settingRepo.GetTableFromStruct(models.Hsms01Bridge{})

	guardTable := u.settingRepo.GetTableFromStruct(models.HsmsMotorwayGuardrail{})

	guardTableType := u.settingRepo.GetTableFromStruct(models.Hsms01Guard{})

	interchangeTable := u.settingRepo.GetTableFromStruct(models.HsmsMotorwayInterchange{})

	interchangeTableType := u.settingRepo.GetTableFromStruct(models.Hsms01Interchange{})

	intersectionTable := u.settingRepo.GetTableFromStruct(models.HsmsMotorwayIntersection{})

	intersectionTableType := u.settingRepo.GetTableFromStruct(models.Hsms01Intersection{})

	lightTable := u.settingRepo.GetTableFromStruct(models.HsmsMotorwayStreetlight{})

	lightTableType := u.settingRepo.GetTableFromStruct(models.Hsms01Light{})

	railwaycrossingTable := u.settingRepo.GetTableFromStruct(models.HsmsMotorwayRailwaycrossing{})

	railwaycrossingTableType := u.settingRepo.GetTableFromStruct(models.Hsms01Railwaycrossing{})

	trafficLightTable := u.settingRepo.GetTableFromStruct(models.HsmsMotorwayTrafficlight{})

	trafficLightTableType := u.settingRepo.GetTableFromStruct(models.Hsms01Signal{})

	uturnBridgeTable := u.settingRepo.GetTableFromStruct(models.HsmsMotorwayUTurn{})

	uturnBridgeTableType := u.settingRepo.GetTableFromStruct(models.Hsms01Uturnbridge{})

	/////////////////////////////////////////////////////////////

	bridge, err := u.settingRepo.GetHsmsBridge()
	if err != nil {
		return hsmsAll, err
	}

	guard, err := u.settingRepo.GetHsmsGuard()
	if err != nil {
		return hsmsAll, err
	}

	interchange, err := u.settingRepo.GetHsmsInterchange()
	if err != nil {
		return hsmsAll, err
	}

	intersection, err := u.settingRepo.GetHsmsIntersection()
	if err != nil {
		return hsmsAll, err
	}

	light, err := u.settingRepo.GetHsmsStreetlight()
	if err != nil {
		return hsmsAll, err
	}

	railwaycrossing, err := u.settingRepo.GetHsmsRailwaycrossing()
	if err != nil {
		return hsmsAll, err
	}

	trafficLight, err := u.settingRepo.GetHsmsTrafficlight()
	if err != nil {
		return hsmsAll, err
	}

	uturnBridge, err := u.settingRepo.GetHsmsUturnbridge()
	if err != nil {
		return hsmsAll, err
	}

	/////////////////////////////////////////////////////////////

	if isAssetName || strings.Contains(refAssetTableTableNameMap[bridgeTable].TableLabel, assetName) {
		for _, item := range bridge {
			var hsms responses.HsmsAll
			hsms.Id = item.ID
			hsms.Type = bridgeTableType
			hsms.AssetName = u.IsStringEmpty(refAssetTableTableNameMap[bridgeTable].TableLabel)
			hsms.RoadGroupName = u.IsStringEmpty(roadGroupNumbeMap[item.RoadCode].ShortName)
			hsms.RoadName = u.IsStringEmpty(roadSectionNumbeMap[item.SectionCode].NameOriginTH + " - " + roadSectionNumbeMap[item.SectionCode].NameDestinationTH)
			hsms.Km = u.IsStringEmpty(item.KM)
			hsms.KmRange = u.IsStringEmpty("")
			hsms.LocationName = u.IsStringEmpty(item.Location)
			hsms.LocationTypeName = u.IsStringEmpty("")
			hsms.DepotName = u.IsStringEmpty(item.DepotName)

			hsmsAll = append(hsmsAll, hsms)
		}
	}

	if isAssetName || strings.Contains(refAssetTableTableNameMap[guardTable].TableLabel, assetName) {
		for _, item := range guard {
			var hsms responses.HsmsAll
			hsms.Id = item.Id
			hsms.Type = guardTableType
			hsms.AssetName = u.IsStringEmpty(refAssetTableTableNameMap[guardTable].TableLabel)
			hsms.RoadGroupName = u.IsStringEmpty(roadGroupNumbeMap[item.RoadCode].ShortName)
			hsms.RoadName = u.IsStringEmpty(roadSectionNumbeMap[item.SectionCode].NameOriginTH + " - " + roadSectionNumbeMap[item.SectionCode].NameDestinationTH)
			hsms.Km = u.IsStringEmpty("")
			hsms.KmRange = u.IsStringEmpty(item.KmStart + " - " + item.KmEnd)
			hsms.LocationName = u.IsStringEmpty(item.Location)
			hsms.LocationTypeName = u.IsStringEmpty(item.LocationTypeText)
			hsms.DepotName = u.IsStringEmpty(item.DepotName)

			hsmsAll = append(hsmsAll, hsms)
		}
	}

	if isAssetName || strings.Contains(refAssetTableTableNameMap[interchangeTable].TableLabel, assetName) {
		for _, item := range interchange {
			var hsms responses.HsmsAll
			hsms.Id = item.Id
			hsms.Type = interchangeTableType
			hsms.AssetName = u.IsStringEmpty(refAssetTableTableNameMap[interchangeTable].TableLabel)
			hsms.RoadGroupName = u.IsStringEmpty(roadGroupNumbeMap[item.RoadCode].ShortName)
			hsms.RoadName = u.IsStringEmpty(roadSectionNumbeMap[item.SectionCode].NameOriginTH + " - " + roadSectionNumbeMap[item.SectionCode].NameDestinationTH)
			hsms.Km = u.IsStringEmpty(item.Km)
			hsms.KmRange = u.IsStringEmpty("")
			hsms.LocationName = u.IsStringEmpty(item.Location)
			hsms.LocationTypeName = u.IsStringEmpty("-")
			hsms.DepotName = u.IsStringEmpty(item.DepotName)

			hsmsAll = append(hsmsAll, hsms)
		}
	}

	if isAssetName || strings.Contains(refAssetTableTableNameMap[intersectionTable].TableLabel, assetName) {
		for _, item := range intersection {
			var hsms responses.HsmsAll
			hsms.Id = item.Id
			hsms.Type = intersectionTableType
			hsms.AssetName = u.IsStringEmpty(refAssetTableTableNameMap[intersectionTable].TableLabel)
			hsms.RoadGroupName = u.IsStringEmpty(roadGroupNumbeMap[item.RoadCode].ShortName)
			hsms.RoadName = u.IsStringEmpty(roadSectionNumbeMap[item.SectionCode].NameOriginTH + " - " + roadSectionNumbeMap[item.SectionCode].NameDestinationTH)
			hsms.Km = u.IsStringEmpty(item.Km)
			hsms.KmRange = u.IsStringEmpty("")
			hsms.LocationName = u.IsStringEmpty("")
			hsms.LocationTypeName = u.IsStringEmpty("")
			hsms.DepotName = u.IsStringEmpty(item.DepotName)

			hsmsAll = append(hsmsAll, hsms)
		}
	}

	if isAssetName || strings.Contains(refAssetTableTableNameMap[lightTable].TableLabel, assetName) {
		for _, item := range light {
			var hsms responses.HsmsAll
			hsms.Id = item.Id
			hsms.Type = lightTableType
			hsms.AssetName = u.IsStringEmpty(refAssetTableTableNameMap[lightTable].TableLabel)
			hsms.RoadGroupName = u.IsStringEmpty(roadGroupNumbeMap[item.RoadCode].ShortName)
			hsms.RoadName = u.IsStringEmpty(roadSectionNumbeMap[item.SectionCode].NameOriginTH + " - " + roadSectionNumbeMap[item.SectionCode].NameDestinationTH)
			hsms.Km = u.IsStringEmpty("")
			hsms.KmRange = u.IsStringEmpty(item.KmStart + " - " + item.KmEnd)
			hsms.LocationName = u.IsStringEmpty(item.Location)
			hsms.LocationTypeName = u.IsStringEmpty(item.LocationTypeText)
			hsms.DepotName = u.IsStringEmpty(item.DepotName)

			hsmsAll = append(hsmsAll, hsms)
		}
	}

	if isAssetName || strings.Contains(refAssetTableTableNameMap[railwaycrossingTable].TableLabel, assetName) {
		for _, item := range railwaycrossing {
			var hsms responses.HsmsAll
			hsms.Id = item.Id
			hsms.Type = railwaycrossingTableType
			hsms.AssetName = u.IsStringEmpty(refAssetTableTableNameMap[railwaycrossingTable].TableLabel)
			hsms.RoadGroupName = u.IsStringEmpty(roadGroupNumbeMap[item.RoadCode].ShortName)
			hsms.RoadName = u.IsStringEmpty(roadSectionNumbeMap[item.SectionCode].NameOriginTH + " - " + roadSectionNumbeMap[item.SectionCode].NameDestinationTH)
			hsms.Km = u.IsStringEmpty(item.Km)
			hsms.KmRange = u.IsStringEmpty("")
			hsms.LocationName = u.IsStringEmpty(item.Location)
			hsms.LocationTypeName = u.IsStringEmpty("")
			hsms.DepotName = u.IsStringEmpty(item.DepotName)

			hsmsAll = append(hsmsAll, hsms)
		}
	}

	if isAssetName || strings.Contains(refAssetTableTableNameMap[trafficLightTable].TableLabel, assetName) {
		for _, item := range trafficLight {
			var hsms responses.HsmsAll
			hsms.Id = item.Id
			hsms.Type = trafficLightTableType
			hsms.AssetName = u.IsStringEmpty(refAssetTableTableNameMap[trafficLightTable].TableLabel)
			hsms.RoadGroupName = u.IsStringEmpty(roadGroupNumbeMap[item.RoadCode].ShortName)
			hsms.RoadName = u.IsStringEmpty(roadSectionNumbeMap[item.SectionCode].NameOriginTH + " - " + roadSectionNumbeMap[item.SectionCode].NameDestinationTH)
			hsms.Km = u.IsStringEmpty(item.Km)
			hsms.KmRange = u.IsStringEmpty("")
			hsms.LocationName = u.IsStringEmpty(item.Location)
			hsms.LocationTypeName = u.IsStringEmpty(item.LocationTypeText)
			hsms.DepotName = u.IsStringEmpty(item.DepotName)

			hsmsAll = append(hsmsAll, hsms)
		}
	}

	if isAssetName || strings.Contains(refAssetTableTableNameMap[uturnBridgeTable].TableLabel, assetName) {
		for _, item := range uturnBridge {
			var hsms responses.HsmsAll
			hsms.Id = item.Id
			hsms.Type = uturnBridgeTableType
			hsms.AssetName = u.IsStringEmpty(refAssetTableTableNameMap[uturnBridgeTable].TableLabel)
			hsms.RoadGroupName = u.IsStringEmpty(roadGroupNumbeMap[item.RoadCode].ShortName)
			hsms.RoadName = u.IsStringEmpty(roadSectionNumbeMap[item.SectionCode].NameOriginTH + " - " + roadSectionNumbeMap[item.SectionCode].NameDestinationTH)
			hsms.Km = u.IsStringEmpty(item.Km)
			hsms.KmRange = u.IsStringEmpty("")
			hsms.LocationName = u.IsStringEmpty("")
			hsms.LocationTypeName = u.IsStringEmpty("")
			hsms.DepotName = u.IsStringEmpty(item.DepotName)

			hsmsAll = append(hsmsAll, hsms)
		}
	}
	return hsmsAll, nil
}

func (r *settingUseCase) FloatToString(f float64) string {
	return fmt.Sprintf("%.3f", f)
}

func (r *settingUseCase) IsStringEmpty(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func (r *settingUseCase) DeleteHsmsByTypeAndId(typeData, id string) error {
	err := r.settingRepo.DeleteHsmsByTypeAndId(typeData, id)
	if err != nil {
		return err
	}
	return nil
}
