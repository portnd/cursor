import { useInitDataStore } from "~~/core/modules/initData/store"

const useInitData = () => {
	const initDataStore = useInitDataStore()

	const conditionGrade = () => {
		return initDataStore.data?.condition_grade ? initDataStore.data?.condition_grade : []
	}

	const selectTemplateForm = (key: keyof typeof initDataStore.data) => {
		return initDataStore.data?.[key]
	}

	const refAsset = () => {
		return initDataStore.data?.ref_asset
	}

	const refAssetBarrier = () => {
		return initDataStore.data?.ref_asset_barrier
	}

	const refAssetBuilding = () => {
		return initDataStore.data?.ref_asset_building
	}

	const refAssetCable = () => {
		return initDataStore.data?.ref_asset_cable
	}

	const refAssetCctvPosition = () => {
		return initDataStore.data?.ref_asset_cctv_position
	}

	const refAssetCctvType = () => {
		return initDataStore.data?.ref_asset_cctv_type
	}

	const refAssetCrashcushion = () => {
		return initDataStore.data?.ref_asset_crashcushion
	}

	const refAssetCurve = () => {
		return initDataStore.data?.ref_asset_curve
	}

	const refAssetDrawpit = () => {
		return initDataStore.data?.ref_asset_drawpit
	}

	const refAssetElectricpost = () => {
		return initDataStore.data?.ref_asset_electricpost
	}

	const refAssetEts = () => {
		return initDataStore.data?.ref_asset_ets
	}

	const refAssetExpansionjoint = () => {
		return initDataStore.data?.ref_asset_expansionjoint
	}

	const refAssetFence = () => {
		return initDataStore.data?.ref_asset_fence
	}

	const refAssetFingerjoint = () => {
		return initDataStore.data?.ref_asset_fingerjoint
	}

	const refAssetFlashlight = () => {
		return initDataStore.data?.ref_asset_flashlight
	}

	const refAssetGuardrail = () => {
		return initDataStore.data?.ref_asset_guardrail
	}

	const refAssetKiosk = () => {
		return initDataStore.data?.ref_asset_kiosk
	}

	const refAssetKm = () => {
		return initDataStore.data?.ref_asset_km
	}

	const refAssetLane = () => {
		return initDataStore.data?.ref_asset_lane
	}

	const refAssetLine = () => {
		return initDataStore.data?.ref_asset_line
	}

	const refAssetLineColor = () => {
		return initDataStore.data?.ref_asset_line_color
	}

	const refAssetLocation = () => {
		return initDataStore.data?.ref_asset_location
	}

	const refAssetManholecover = () => {
		return initDataStore.data?.ref_asset_manholecover
	}

	const refAssetNoisebarrier = () => {
		return initDataStore.data?.ref_asset_noisebarrier
	}

	const refAssetPlatejoint = () => {
		return initDataStore.data?.ref_asset_platejoint
	}

	const refAssetPlugjoint = () => {
		return initDataStore.data?.ref_asset_plugjoint
	}

	const refAssetPole = () => {
		return initDataStore.data?.ref_asset_pole
	}

	const refAssetPosition = () => {
		return initDataStore.data?.ref_asset_position
	}

	const refAssetSafetyswitch = () => {
		return initDataStore.data?.ref_asset_safetyswitch
	}

	const refAssetSign = () => {
		return initDataStore.data?.ref_asset_sign
	}

	const refAssetSignImage = () => {
		return initDataStore.data?.ref_asset_sign_image
	}

	const refAssetSignType = () => {
		return initDataStore.data?.ref_asset_sign_type
	}

	const refAssetSupplypillar = () => {
		return initDataStore.data?.ref_asset_supplypillar
	}

	const refAssetTable = () => {
		return initDataStore.data?.ref_asset_table
	}

	const refAssetTableColumns = () => {
		return initDataStore.data?.ref_asset_table_columns
	}

	const refAssetTableStaff = () => {
		return initDataStore.data?.ref_asset_table_staff
	}

	const refAssetTollplaza = () => {
		return initDataStore.data?.ref_asset_tollplaza
	}

	const refAssetTransformer = () => {
		return initDataStore.data?.ref_asset_transformer
	}

	const refAssetTube = () => {
		return initDataStore.data?.ref_asset_tube
	}

	const refColorlist = () => {
		return initDataStore.data?.ref_color_list
	}

	const refDataStatus = () => {
		return initDataStore.data?.ref_data_status
	}

	const refDepartment = () => {
		return initDataStore.data?.ref_department
	}

	const refDirection = () => {
		return initDataStore.data?.ref_direction
	}

	const refGrade = () => {
		return initDataStore.data?.ref_grade
	}

	const refMaterialBase = () => {
		return initDataStore.data?.ref_material_base
	}

	const refMaterialSubbase = () => {
		return initDataStore.data?.ref_material_subbase
	}

	const refMaterialSubgrade = () => {
		return initDataStore.data?.ref_material_subgrade
	}

	const refOwner = () => {
		return initDataStore.data?.ref_owner
	}

	const refRoadType = () => {
		return initDataStore.data?.ref_road_type
	}

	const refRoadTypeLevelFirst = () => {
		return initDataStore.data?.ref_road_type_level_first
	}

	const refRoadTypeLevelSecond = () => {
		return initDataStore.data?.ref_road_type_level_second
	}

	const refStructureSurface = () => {
		return initDataStore.data?.ref_structure_surface
	}

	const refSurface = () => {
		return initDataStore.data?.ref_surface
	}

	const refSurfaceGroup = () => {
		return initDataStore.data?.ref_surface_group
	}

	const refSurfaceType = () => {
		return initDataStore.data?.ref_surface_type
	}

	const refRoadTypeIcon = () => {
		return initDataStore.data?.ref_road_type_icon
	}

	const refTableList = () => {
		return initDataStore.data?.ref_table_list
	}

	const refRoadGroup = () => {
		return initDataStore.data?.road_group
	}

	const refRoadSection = () => {
		return initDataStore.data?.road_section
	}

	const refDistrict = () => {
		return initDataStore.data?.ref_district
	}

	const refConditionRange = () => {
		return initDataStore.data?.ref_condition_range
	}

	const refReflectivityRange = () => {
		return initDataStore.data?.ref_reflectivity_range
	}

	const reflectivityGrade = () => {
		return initDataStore.data?.reflectivity_grade
	}

	const refStripeColor = () => {
		return initDataStore.data?.ref_stripe_color
	}

	const refStripeType = () => {
		return initDataStore.data?.ref_stripe_type
	}

	const refCriteriaType = () => {
		return initDataStore.data?.ref_criteria_type
	}

	const refDivision = () => {
		return initDataStore.data?.ref_division
	}

	const refDivisionDashboardAsset = () => {
		return initDataStore.data?.ref_division_dashboard_asset
	}

	const refDivisionDashboardCondition = () => {
		return initDataStore.data?.ref_division_dashboard_condition
	}

	const refDivisionDashboardMaintenance = () => {
		return initDataStore.data?.ref_division_dashboard_maintenance
	}

	const refDivisionDashboardSurface = () => {
		return initDataStore.data?.ref_division_dashboard_surface
	}

	const refUserOwner = () => {
		return initDataStore.data?.ref_user_owner
	}

	return {
		refAsset,
		refAssetBarrier,
		refAssetBuilding,
		refAssetCable,
		refAssetCctvPosition,
		refAssetCctvType,
		refAssetCrashcushion,
		refAssetCurve,
		refAssetDrawpit,
		refAssetElectricpost,
		refAssetEts,
		refAssetExpansionjoint,
		refAssetFence,
		refAssetFingerjoint,
		refAssetFlashlight,
		refAssetKiosk,
		refAssetKm,
		refAssetGuardrail,
		refAssetLane,
		refAssetLine,
		refAssetLineColor,
		refAssetLocation,
		refAssetManholecover,
		refAssetNoisebarrier,
		refAssetPlatejoint,
		refAssetPlugjoint,
		refAssetPole,
		refAssetPosition,
		refAssetSafetyswitch,
		refAssetSign,
		refAssetSignImage,
		refAssetSignType,
		refAssetSupplypillar,
		refAssetTable,
		refAssetTableColumns,
		refAssetTableStaff,
		refAssetTollplaza,
		refAssetTransformer,
		refAssetTube,
		refColorlist,
		refDataStatus,
		refDepartment,
		refDirection,
		refGrade,
		refMaterialBase,
		refMaterialSubbase,
		refMaterialSubgrade,
		refOwner,
		refRoadType,
		refRoadTypeLevelFirst,
		refRoadTypeLevelSecond,
		refSurface,
		refSurfaceGroup,
		refSurfaceType,
		refRoadTypeIcon,
		refStructureSurface,
		refTableList,
		conditionGrade,
		selectTemplateForm,
		refRoadGroup,
		refRoadSection,
		refDistrict,
		refConditionRange,
		refReflectivityRange,
		reflectivityGrade,
		refStripeColor,
		refStripeType,
		refCriteriaType,
		refDivision,
		refDivisionDashboardAsset,
		refDivisionDashboardCondition,
		refDivisionDashboardMaintenance,
		refDivisionDashboardSurface,
		refUserOwner,
	}
}

export default useInitData
