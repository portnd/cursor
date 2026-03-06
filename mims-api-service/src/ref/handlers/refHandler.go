package handlers

import (
	"net/http"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/ref/domains"

	"github.com/gin-gonic/gin"
)

type RefHandler struct {
	refUseCase domains.RefUseCase
}

func NewRefHandler(usecase domains.RefUseCase) *RefHandler {
	return &RefHandler{
		refUseCase: usecase,
	}
}

// @summary
// @description
// @tags init data
// @id init_data
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.InitDataResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/initdata [get]
func (rh *RefHandler) InitData(c *gin.Context) {
	userID := helpers.GetUserID(c)
	refAsset, err := rh.refUseCase.GetRefAsset()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetPosition, err := rh.refUseCase.GetRefAssetPosition()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetArea, err := rh.refUseCase.GetRefAssetArea()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetGuardrail, err := rh.refUseCase.GetRefAssetGuardrail()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetReflecType, err := rh.refUseCase.GetRefAssetReflecType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetOwner, err := rh.refUseCase.GetRefAssetOwner()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetLightType, err := rh.refUseCase.GetRefAssetLightType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetLightWatt, err := rh.refUseCase.GetRefAssetLightWatt()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetCleranceType, err := rh.refUseCase.GetRefAssetCleranceType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetCrashcushionType, err := rh.refUseCase.GetRefAssetCrashcushionType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetFenceType, err := rh.refUseCase.GetRefAssetFenceType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetKmstoneType, err := rh.refUseCase.GetRefAssetKmstoneType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetQutterType, err := rh.refUseCase.GetRefAssetQutterType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetTrafficCameraType, err := rh.refUseCase.GetRefAssetTrafficCameraType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetWeightStationType, err := rh.refUseCase.GetRefAssetWeightStationType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetBuildingType, err := rh.refUseCase.GetRefAssetBuildingType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetNoiseBarrier, err := rh.refUseCase.GetRefAssetNoiseBarrier()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetSignImage, err := rh.refUseCase.GetRefAssetSignImage()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetTable, err := rh.refUseCase.GetRefAssetTable()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetTableColumns, err := rh.refUseCase.GetRefAssetTableColumns()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refAssetTableStaff, err := rh.refUseCase.GetRefAssetTableStaff()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	////////////////////////////////////////////////////////////////
	refDataStatus, err := rh.refUseCase.GetRefDataStatus()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refDepartment, err := rh.refUseCase.GetRefDepartment()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refDirection, err := rh.refUseCase.GetRefDirection()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refGrade, err := rh.refUseCase.GetRefGrade()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refMaterialBase, err := rh.refUseCase.GetRefMaterialBase()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refMaterialSubbase, err := rh.refUseCase.GetRefMaterialSubbase()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refMaterialSubgrade, err := rh.refUseCase.GetRefMaterialSubgrade()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refOwner, err := rh.refUseCase.GetRefOwner()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refRoadType, err := rh.refUseCase.GetRefRoadType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refSurface, err := rh.refUseCase.GetRefSurface()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refSurfaceType, err := rh.refUseCase.GetRefSurfaceType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refStructureSurface, err := rh.refUseCase.GetRefStructureSurface()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refSurfaceGroup, err := rh.refUseCase.GetRefSurfaceGroup()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refColorList, err := rh.refUseCase.GetRefColorList()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refRoadTypeIcon, err := rh.refUseCase.GetRefRoadTypeIcon()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refTableList, err := rh.refUseCase.GetRefTableList()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	conditionGrades, err := rh.refUseCase.GetRoadConditionGrades()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	GetRefRoadConditionRange, err := rh.refUseCase.GetRefRoadConditionRange()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	GetRefReflectivityRange, err := rh.refUseCase.GetRefReflectivityRange()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	RoadLineList, err := rh.refUseCase.GetRoadLineConditionGrades()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	roadGroupList, err := rh.refUseCase.GetRoadGroupList()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	roadSectionList, err := rh.refUseCase.GetRoadSectionList(userID)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refDistrictsList, err := rh.refUseCase.GetRefDistrictsInitList(userID)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refRoadTypeOne, err := rh.refUseCase.GetRefRoadTypeLevel("id IN (1, 2, 3, 4)")
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refRoadTypeTwo, err := rh.refUseCase.GetRefRoadTypeLevel("id IN (5, 6, 7)")
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refStripeType, err := rh.refUseCase.GetRefStripeType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refStripeColor, err := rh.refUseCase.GetRefStripeColor()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)

		}
	}

	refDivision, err := rh.refUseCase.GetRefDivisionInitList(userID)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)

		}
	}

	refDivisionDashboardAsset, err := rh.refUseCase.GetRefDivisionInitListDashboardAsset(userID)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
	}

	refDivisionDashboardCondition, err := rh.refUseCase.GetRefDivisionInitListDashboardCondition(userID)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
	}

	refDivisionDashboardSurface, err := rh.refUseCase.GetRefDivisionInitListDashboardSurface(userID)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
	}

	refDivisionDashboardMaintenance, err := rh.refUseCase.GetRefDivisionInitListDashboardMaintenance(userID)
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
	}

	refCriteriaType, err := rh.refUseCase.GetRefCriteriaType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refCriteriaMethod, err := rh.refUseCase.GetRefCriteriaMethod()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	refUserOwner, err := rh.refUseCase.GetRefUserOwner()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	initData := responses.InitData{
		RefAsset:                  refAsset.([]models.RefAsset),
		RefAssetPosition:          refAssetPosition.([]models.RefAssetPosition),
		RefAssetArea:              refAssetArea.([]models.RefAssetArea),
		RefAssetGuardrail:         refAssetGuardrail.([]models.RefAssetGuardrail),
		RefAssetReflecType:        refAssetReflecType.([]models.RefAssetReflecType),
		RefAssetOwner:             refAssetOwner.([]models.RefAssetOwner),
		RefAssetLightType:         refAssetLightType.([]models.RefAssetLightType),
		RefAssetLightWatt:         refAssetLightWatt.([]models.RefAssetLightWatt),
		RefAssetCleranceType:      refAssetCleranceType.([]models.RefAssetCleranceType),
		RefAssetCrashcushionType:  refAssetCrashcushionType.([]models.RefAssetCrashcushionType),
		RefAssetFenceType:         refAssetFenceType.([]models.RefAssetFenceType),
		RefAssetKmstoneType:       refAssetKmstoneType.([]models.RefAssetKmstoneType),
		RefAssetQutterType:        refAssetQutterType.([]models.RefAssetQutterType),
		RefAssetTrafficCameraType: refAssetTrafficCameraType.([]models.RefAssetTrafficCameraType),
		RefAssetWeightStationType: refAssetWeightStationType.([]models.RefAssetWeightStationType),
		RefAssetBuildingType:      refAssetBuildingType.([]models.RefAssetBuildingType),
		RefAssetNoiseBarrier:      refAssetNoiseBarrier.([]models.RefAssetNoiseBarrier),
		RefAssetSignImage:         refAssetSignImage.([]models.RefAssetSignImage),
		RefAssetTable:             refAssetTable.([]models.RefAssetTable),
		RefAssetTableColumns:      refAssetTableColumns.([]models.RefAssetTableColumns),
		RefAssetTableStaff:        refAssetTableStaff.([]models.RefAssetTableStaff),

		RefDataStaus:                    refDataStatus.([]models.RefDataStatus),
		RefDepartment:                   refDepartment.([]models.RefDepartment),
		RefDirection:                    refDirection.([]models.RefDirection),
		RefGrade:                        refGrade.([]models.RefGrade),
		RefMaterialBase:                 refMaterialBase.([]models.RefMaterialBase),
		RefMaterialSubbase:              refMaterialSubbase.([]models.RefMaterialSubbase),
		RefMaterialSubgrade:             refMaterialSubgrade.([]models.RefMaterialSubgrade),
		RefOwner:                        refOwner.([]models.RefOwner),
		RefRoadType:                     refRoadType.([]models.RefRoadType),
		RefSurface:                      refSurface.([]models.RefSurface),
		RefSurfaceType:                  refSurfaceType.([]models.RefSurfaceType),
		RefStructureSurface:             refStructureSurface.([]models.RefStructureSurface),
		RefSurfaceGroup:                 refSurfaceGroup.([]responses.RefSurfaceGroup),
		RefColorList:                    refColorList.([]models.RefColorList),
		RefRoadTypeIcon:                 refRoadTypeIcon.([]models.RefRoadTypeIcon),
		RefTableList:                    refTableList,
		ConditionGrade:                  conditionGrades.([]responses.ConditionRespondInit),
		RoadLineList:                    RoadLineList.([]responses.RoadLineListInit),
		RefConditionRange:               GetRefRoadConditionRange.([]models.RefConditionRange),
		RefReflectivityRange:            GetRefReflectivityRange.([]models.RefReflectivityRange),
		RoadGroup:                       roadGroupList,
		RoadSection:                     roadSectionList,
		RefDistrict:                     refDistrictsList.([]models.RefDistrictInitData),
		RefRoadTypeLevelOne:             refRoadTypeOne,
		RefRoadTypeLevelTwo:             refRoadTypeTwo,
		RefStripeType:                   refStripeType.([]models.RefStripeType),
		RefStripeColor:                  refStripeColor.([]models.RefStripeColor),
		RefDivision:                     refDivision,
		RefDivisionDashboardAsset:       refDivisionDashboardAsset,
		RefDivisionDashboardCondition:   refDivisionDashboardCondition,
		RefDivisionDashboardSurface:     refDivisionDashboardSurface,
		RefDivisionDashboardMaintenance: refDivisionDashboardMaintenance,
		RefCriteriaType:                 refCriteriaType.([]models.RefCriteriaType),
		RefCriteriaMethod:               refCriteriaMethod,
		RefUserOwner:                    refUserOwner,
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(initData, http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_asset
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.RefAssetResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/asset [get]
func (rh *RefHandler) GetRefAsset(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefAsset()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAsset), http.StatusOK))
}

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_barrier
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_barrier [get]
// func (rh *RefHandler) GetRefAssetBarrier(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetBarrier()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetBarrier), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_building
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_building [get]
// func (rh *RefHandler) GetRefAssetBuilding(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetBuilding()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetBuilding), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_cable
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_cable [get]
// func (rh *RefHandler) GetRefAssetCable(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetCable()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetCable), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_cctv_position
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_cctv_position [get]
// func (rh *RefHandler) GetRefAssetCCTVPosition(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetCCTVPosition()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetCctvPosition), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_cctv_type
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_cctv_type [get]
// func (rh *RefHandler) GetRefAssetCCTVType(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetCCTVType()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetCctvType), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_crashcushion
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_crashcushion [get]
// func (rh *RefHandler) GetRefAssetCrashcusion(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetCrashcushion()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetCrashcushion), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_curve
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_curve [get]
// func (rh *RefHandler) GetRefAssetCurve(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetCurve()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetCurve), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_drawpit
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_drawpit [get]
// func (rh *RefHandler) GetRefAssetDrawpit(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetDrawpit()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetDrawpit), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_electricpost
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_electricpost [get]
// func (rh *RefHandler) GetRefAssetElectricpost(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetElectricpost()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetElectricpost), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_ets
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_ets [get]
// func (rh *RefHandler) GetRefAssetEts(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetEts()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetEts), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_expansionjoint
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_expansionjoint [get]
// func (rh *RefHandler) GetRefAssetExpansionjoint(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetExpansionjoint()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetExpansionjoint), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_fence
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_fence [get]
// func (rh *RefHandler) GetRefAssetFence(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetFence()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetFence), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_fingerjoint
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_fingerjoint [get]
// func (rh *RefHandler) GetRefAssetFingerjoint(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetFingerjoint()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetFingerjoint), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_flashlight
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_flashlight [get]
// func (rh *RefHandler) GetRefAssetFlashlight(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetFlashlight()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetFlashlight), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_guardrail
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_guardrail [get]
// func (rh *RefHandler) GetRefAssetGuardrail(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetGuardrail()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetGuardrail), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_kiosk
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_kiosk [get]
// func (rh *RefHandler) GetRefAssetKiosk(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetKiosk()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetKiosk), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_km
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_km [get]
// func (rh *RefHandler) GetRefAssetKm(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetKm()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetKm), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_lane
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_lane [get]
// func (rh *RefHandler) GetRefAssetLane(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetLane()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetLane), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_line
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_line [get]
// func (rh *RefHandler) GetRefAssetLine(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetLine()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetLine), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_line_color
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_line_color [get]
// func (rh *RefHandler) GetRefAssetLineColor(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetLineColor()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetLineColor), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_location
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_location [get]
// func (rh *RefHandler) GetRefAssetLocation(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetLocation()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetLocation), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_manholecover
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_manholecover [get]
// func (rh *RefHandler) GetRefAssetManholecover(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetManholecover()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetManholecover), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_noisebarrier
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_noisebarrier [get]
// func (rh *RefHandler) GetRefAssetNoisebarrier(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetNoisebarrier()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetNoisebarrier), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_platelight
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_platelight [get]
// func (rh *RefHandler) GetRefAssetPlatelight(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetPlatelight()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetPlatelight), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_plugjoint
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_plugjoint [get]
// func (rh *RefHandler) GetRefAssetPlugjoint(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetPlugjoint()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetPlugjoint), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_pole
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_pole [get]
// func (rh *RefHandler) GetRefAssetPole(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetPole()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetPole), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_position
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_position [get]
// func (rh *RefHandler) GetRefAssetPosition(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetPosition()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetPosition), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_safetyswitch
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_safetyswitch [get]
// func (rh *RefHandler) GetRefAssetSafetyswitch(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetSafetyswitch()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetSafetyswitch), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_sign
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_sign [get]
// func (rh *RefHandler) GetRefAssetSign(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetSign()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetSign), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_sign_image
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.RefAsestSignImageReponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_sign_image [get]
// func (rh *RefHandler) GetRefAssetSignImage(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetSignImage()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetSignImage), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_sign_type
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_sign_type [get]
// func (rh *RefHandler) GetRefAssetSignType(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetSignType()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetSignType), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_supplypillar
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_supplypillar [get]
// func (rh *RefHandler) GetRefAssetSupplypillar(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetSupplypillar()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetSupplypillar), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_table
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.RefAssetTableResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_table [get]
// func (rh *RefHandler) GetRefAssetTable(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetTable()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetTable), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_table_columns
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.RefAssetTableColumnsResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_table_columns [get]
// func (rh *RefHandler) GetRefAssetTableColumns(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetTableColumns()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetTableColumns), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_table_staff
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.RefAssetTableStaffResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_table_staff [get]
// func (rh *RefHandler) GetRefAssetTableStaff(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetTableStaff()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetTableStaff), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_tollplaza
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_tollplaza [get]
// func (rh *RefHandler) GetRefAssetTollplaza(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetTollplaza()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetTollplaza), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_transformer
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_transformer [get]
// func (rh *RefHandler) GetRefAssetTransformer(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetTransformer()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetTransformer), http.StatusOK))
// }

// // @summary
// // @description
// // @tags master data ref
// // @id ref_asset_tube
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// // @response 200 {object} responses.FirstGroupResponse "OK"
// // @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// // @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// // @router /api/v1/ref/asset_tube [get]
// func (rh *RefHandler) GetRefAssetTube(c *gin.Context) {
// 	resp, err := rh.refUseCase.GetRefAssetTube()
// 	if err != nil {
// 		appErr, ok := err.(*responses.AppErr)
// 		if ok {
// 			errResponse := responses.FailRespone(appErr)
// 			c.JSON(appErr.StatusCode, errResponse)
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefAssetTube), http.StatusOK))
// }

// @summary
// @description
// @tags master data ref
// @id ref_data_status
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.RefDataStatusResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/data_status [get]
func (rh *RefHandler) GetRefDataStatus(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefDataStatus()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefDataStatus), http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_department
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.FirstGroupResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/department [get]
func (rh *RefHandler) GetRefDepartment(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefDepartment()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefDepartment), http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_direction
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.FirstGroupResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/direction [get]
func (rh *RefHandler) GetRefDirection(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefDirection()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefDirection), http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_grade
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.FirstGroupResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/grade [get]
func (rh *RefHandler) GetRefGrade(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefGrade()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefGrade), http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_material_base
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.MaterialResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/material_base [get]
func (rh *RefHandler) GetRefMaterialBase(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefMaterialBase()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefMaterialBase), http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_material_subbase
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.MaterialResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/material_subbase [get]
func (rh *RefHandler) GetRefMaterialSubbase(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefMaterialSubbase()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefMaterialSubbase), http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_material_subgrade
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.MaterialResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/material_subgrade [get]
func (rh *RefHandler) GetRefMaterialSubgrade(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefMaterialSubgrade()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefMaterialSubgrade), http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_owner
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.RefOwnerResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/owner [get]
func (rh *RefHandler) GetRefOwner(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefOwner()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefOwner), http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_road_type
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.FirstGroupResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/road_type [get]
func (rh *RefHandler) GetRefRoadType(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefRoadType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefRoadType), http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_surface
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.FirstGroupResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/surface [get]
func (rh *RefHandler) GetRefSurface(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefSurface()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefSurface), http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_structure_surface
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.FirstGroupResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/structure_surface [get]
func (rh *RefHandler) GetRefStructureSurface(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefSurface()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp.([]models.RefSurface), http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_table_list
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.RefTableListResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/table_list [get]
func (rh *RefHandler) GetRefTableList(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefTableList()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_color_list
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.RefTableListResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/color_list [get]
func (rh *RefHandler) GetRefColorList(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefColorList()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id criteria_type
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.RefTableListResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/criteria_type [get]
func (rh *RefHandler) GetRefCriteriaType(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefCriteriaType()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_parameter_vehicle_type
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @param road_group_id path string false "Insert your road group id"
// @response 200 {object} []responses.RefAadtParameterVehicleType "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/aadt_parameter_vehicle_type/{road_group_id} [get]
func (rh *RefHandler) GetParameterVehicleType(c *gin.Context) {
	resp, err := rh.refUseCase.GetParameterVehicleTypeList(c.Param("road_group_id"))
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_road_user_cost_acc
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.RefTableListResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/road_user_cost/acc [get]
func (rh *RefHandler) GetRoadUserCostAcc(c *gin.Context) {
	resp, err := rh.refUseCase.GetRoadUserCostAcc()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_road_user_cost_ruc
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.RefTableListResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/road_user_cost/ruc [get]
func (rh *RefHandler) GetRoadUserCostRuc(c *gin.Context) {
	resp, err := rh.refUseCase.GetRoadUserCostRuc()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary ข้อมูล master data
// @description
// @tags Analysis
// @id ref_maintenance_analysis_strategic
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.RefTableListResponse "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/maintenance_analysis_strategic [get]
func (rh *RefHandler) GetMaintenanceAnalysisStrategic(c *gin.Context) {
	resp, err := rh.refUseCase.GetMaintenanceAnalysisStrategic()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}

// @summary
// @description
// @tags master data ref
// @id ref_districts
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @response 200 {object} responses.Success{data=[]models.RefDistrictData} "OK"
// @response 401 {object} responses.UnauthorizedErrorResponse "Unauthorized"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/v1/ref/districts [get]
func (rh *RefHandler) GetRefDistrictsList(c *gin.Context) {
	resp, err := rh.refUseCase.GetRefDistrictsList()
	if err != nil {
		appErr, ok := err.(*responses.AppErr)
		if ok {
			errResponse := responses.FailRespone(appErr)
			c.JSON(appErr.StatusCode, errResponse)
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(resp, http.StatusOK))
}
