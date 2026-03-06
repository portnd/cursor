import { IRefStripeColor } from "./data/RefStripeColor"
import { IRefStripeType } from "./data/RefStripeType"
import {
	IRefAsset,
	IRefAssetBarrier,
	IRefAssetBuilding,
	IRefAssetCable,
	IRefAssetCctvPosition,
	IRefAssetCctvType,
	IRefAssetCrashcushion,
	IRefAssetCurve,
	IRefAssetDrawpit,
	IRefAssetElectricpost,
	IRefAssetEts,
	IRefAssetExpansionjoint,
	IRefAssetFence,
	IRefAssetFingerjoint,
	IRefAssetFlashlight,
	IRefAssetGuardrail,
	IRefAssetKiosk,
	IRefAssetKm,
	IRefAssetLane,
	IRefAssetLine,
	IRefAssetLineColor,
	IRefAssetLocation,
	IRefAssetManholecover,
	IRefAssetNoisebarrier,
	IRefAssetPlatejoint,
	IRefAssetPlugjoint,
	IRefAssetPole,
	IRefAssetPosition,
	IRefAssetSafetyswitch,
	IRefAssetSign,
	IRefAssetSignImage,
	IRefAssetSignType,
	IRefAssetSupplypillar,
	IRefAssetTable,
	IRefAssetTableColumns,
	IRefAssetTableStaff,
	IRefAssetTollplaza,
	IRefAssetTransformer,
	IRefAssetTube,
	IRefColorList,
	IRefDataStatus,
	IRefDepartment,
	IRefDirection,
	IRefGrade,
	IRefMaterialBase,
	IRefMaterialSubbase,
	IRefMaterialSubgrade,
	IRefOwner,
	IRefRoadType,
	IRefStructureSurface,
	IRefSurface,
	IRefSurfaceGroup,
	IRefSurfaceType,
	IRefRoadTypeIcon,
	IRefTableList,
	IConditionGroup,
	IRefConditionRange,
	IRefRoadGroup,
	IRefRoadSection,
	IRefDistrict,
	IConditionGrade,
	IReflectivityGrade,
	IRefCriteriaType,
	IRefDivision,
	IRefUserOwner,
} from "./data"

export interface IInitData {
	condition_grade: IConditionGrade[]
	ref_asset: IRefAsset[]
	ref_asset_barrier: IRefAssetBarrier[]
	ref_asset_building: IRefAssetBuilding[]
	ref_asset_cable: IRefAssetCable[]
	ref_asset_cctv_position: IRefAssetCctvPosition[]
	ref_asset_cctv_type: IRefAssetCctvType[]
	ref_asset_crashcushion: IRefAssetCrashcushion[]
	ref_asset_curve: IRefAssetCurve[]
	ref_asset_drawpit: IRefAssetDrawpit[]
	ref_asset_electricpost: IRefAssetElectricpost[]
	ref_asset_ets: IRefAssetEts[]
	ref_asset_expansionjoint: IRefAssetExpansionjoint[]
	ref_asset_fence: IRefAssetFence[]
	ref_asset_fingerjoint: IRefAssetFingerjoint[]
	ref_asset_flashlight: IRefAssetFlashlight[]
	ref_asset_guardrail: IRefAssetGuardrail[]
	ref_asset_kiosk: IRefAssetKiosk[]
	ref_asset_km: IRefAssetKm[]
	ref_asset_lane: IRefAssetLane[]
	ref_asset_line: IRefAssetLine[]
	ref_asset_line_color: IRefAssetLineColor[]
	ref_asset_location: IRefAssetLocation[]
	ref_asset_manholecover: IRefAssetManholecover[]
	ref_asset_noisebarrier: IRefAssetNoisebarrier[]
	ref_asset_platejoint: IRefAssetPlatejoint[]
	ref_asset_plugjoint: IRefAssetPlugjoint[]
	ref_asset_pole: IRefAssetPole[]
	ref_asset_position: IRefAssetPosition[]
	ref_asset_safetyswitch: IRefAssetSafetyswitch[]
	ref_asset_sign: IRefAssetSign[]
	ref_asset_sign_image: IRefAssetSignImage[]
	ref_asset_sign_type: IRefAssetSignType[]
	ref_asset_supplypillar: IRefAssetSupplypillar[]
	ref_asset_table: IRefAssetTable[]
	ref_asset_table_columns: IRefAssetTableColumns[]
	ref_asset_table_staff: IRefAssetTableStaff[]
	ref_asset_tollplaza: IRefAssetTollplaza[]
	ref_asset_transformer: IRefAssetTransformer[]
	ref_asset_tube: IRefAssetTube[]
	ref_color_list: IRefColorList[]
	ref_data_status: IRefDataStatus[]
	ref_department: IRefDepartment[]
	ref_direction: IRefDirection[]
	ref_grade: IRefGrade[]
	ref_material_base: IRefMaterialBase[]
	ref_material_subbase: IRefMaterialSubbase[]
	ref_material_subgrade: IRefMaterialSubgrade[]
	ref_owner: IRefOwner[]
	ref_road_type: IRefRoadType[]
	ref_road_type_level_first: IRefRoadType[]
	ref_road_type_level_second: IRefRoadType[]
	ref_structure_surface: IRefStructureSurface[]
	ref_surface: IRefSurface[]
	ref_surface_group: IRefSurfaceGroup[]
	ref_surface_type: IRefSurfaceType[]
	ref_road_type_icon: IRefRoadTypeIcon[]
	ref_table_list: IRefTableList[]
	ref_condition_range: IRefConditionRange[]
	road_group: IRefRoadGroup[]
	road_section: IRefRoadSection[]
	ref_district: IRefDistrict[]
	ref_reflectivity_range: IRefConditionRange[]
	reflectivity_grade: IReflectivityGrade[]
	ref_stripe_color: IRefStripeColor[]
	ref_stripe_type: IRefStripeType[]
	ref_criteria_type: IRefCriteriaType[]
	ref_division: IRefDivision[]
	ref_division_dashboard_asset: IRefDivision[]
	ref_division_dashboard_condition: IRefDivision[]
	ref_division_dashboard_maintenance: IRefDivision[]
	ref_division_dashboard_surface: IRefDivision[]
	ref_user_owner: IRefUserOwner[]
}
