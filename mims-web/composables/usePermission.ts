import httpStatusCode from "~~/core/shared/http/HttpStatusCode"

enum IUserRolesAccess {
	view_total_dashboard = "view_total_dashboard",
	view_all_asset_dashboard = "view_all_asset_dashboard",
	view_owner_asset_dashboard = "view_owner_asset_dashboard",
	view_all_road_condition_dashboard = "view_all_road_condition_dashboard",
	view_owner_road_condition_dashboard = "view_owner_road_condition_dashboard",
	view_all_surface_dashboard = "view_all_surface_dashboard",
	view_owner_surface_dashboard = "view_owner_surface_dashboard",
	view_all_maint_history_dashboard = "view_all_maint_history_dashboard",
	view_owner_maint_history_dashboard = "view_owner_maint_history_dashboard",

	view_all_road = "view_all_road",
	view_owner_road = "view_owner_road",
	manage_all_road = "manage_all_road",
	manage_owner_road = "manage_owner_road",
	view_road_summary = "view_road_summary",
	manage_road_summary = "manage_road_summary",
	manage_owner_road_summary = "manage_owner_road_summary",
	view_road_in_assets = "view_road_in_assets",
	manage_road_in_assets = "manage_road_in_assets",
	manage_owner_road_in_assets = "manage_owner_road_in_assets",
	view_road_condition = "view_road_condition",
	manage_road_condition = "manage_road_condition",
	manage_owner_road_condition = "manage_owner_road_condition",
	view_road_retro = "view_road_retro",
	manage_road_retro = "manage_road_retro",
	manage_owner_road_retro = "manage_owner_road_retro",
	view_road_damage = "view_road_damage",
	manage_road_damage = "manage_road_damage",
	manage_owner_road_damage = "manage_owner_road_damage",
	view_road_traffic = "view_road_traffic",
	manage_road_traffic = "manage_road_traffic",
	manage_owner_road_traffic = "manage_owner_road_traffic",

	view_all_maint_history = "view_all_maint_history",
	view_owner_maint_history = "view_owner_maint_history",
	manage_all_maint_history = "manage_all_maint_history",
	manage_owner_maint_history = "manage_owner_maint_history",

	view_all_maintenance_analysis = "view_all_maintenance_analysis",
	view_owner_maintenance_analysis = "view_owner_maintenance_analysis",
	manage_myself_maintenance_analysis = "manage_myself_maintenance_analysis",

	download_all_asset_report = "download_all_asset_report",
	download_owner_asset_report = "download_owner_asset_report",
	download_all_map_asset_report = "download_all_map_asset_report",
	download_owner_map_asset_report = "download_owner_map_asset_report",
	download_all_result_asset_report = "download_all_result_asset_report",
	download_owner_result_asset_report = "download_owner_result_asset_report",
	download_all_road_report = "download_all_road_report",
	download_all_surface_report = "download_all_surface_report",
	download_owner_surface_report = "download_owner_surface_report",
	download_all_road_condition_report = "download_all_road_condition_report",
	download_owner_road_condition_report = "download_owner_road_condition_report",
	download_all_road_condition_result_report = "download_all_road_condition_result_report",
	download_owner_road_condition_result_report = "download_owner_road_condition_result_report",
	download_all_road_retro_report = "download_all_road_retro_report",
	download_owner_road_retro_report = "download_owner_road_retro_report",
	download_all_road_retro_result_report = "download_all_road_retro_result_report",
	download_owner_road_retro_result_report = "download_owner_road_retro_result_report",
	download_all_road_damage_report = "download_all_road_damage_report",
	download_owner_road_damage_report = "download_owner_road_damage_report",
	download_all_road_damage_result_report = "download_all_road_damage_result_report",
	download_owner_road_damage_result_report = "download_owner_road_damage_result_report",
	download_all_maint_history_report = "download_all_maint_history_report",
	download_owner_maint_history_report = "download_owner_maint_history_report",
	download_all_maint_kpi_report = "download_all_maint_kpi_report",
	download_owner_maint_kpi_report = "download_owner_maint_kpi_report",
	download_all_aadt_report = "download_all_aadt_report",
	download_owner_aadt_report = "download_owner_aadt_report",

	setting_asset_place_group_access = "setting_asset_place_group_access",
	setting_in_asset_access = "setting_in_asset_access",
	setting_traffic_sign_access = "setting_traffic_sign_access",
	setting_condition_criteria_access = "setting_condition_criteria_access",
	setting_retro_criteria_access = "setting_retro_criteria_access",
	setting_surface_type_access = "setting_surface_type_access",
	setting_ref_budget_access = "setting_ref_budget_access",
	setting_model_access = "setting_model_access",
	setting_hsms_access = "setting_hsms_access",
	setting_hris_access = "setting_hris_access",

	setting_user_access = "setting_user_acess",
	setting_permission_access = "setting_permission_acess",
}

const usePermission = () => {
	const getAccessControls = (): string[] => {
		if (typeof window !== "undefined") {
			const localStorageAuth = window.localStorage.getItem("init-user")

			if (localStorageAuth) {
				const objectStorageAuth = JSON.parse(localStorageAuth)
				const userAccessKey = objectStorageAuth?.userAccessKey
				// console.log("getAccessControls  ✅", userAccessKey)
				return userAccessKey?.length > 0 ? userAccessKey : []
			}
		}

		return []
	}

	const findAccessKey = (key: string): boolean => {
		console.log("findAccessKey  ✅", key)
		return getAccessControls().includes(key) ?? false
	}

	// ส่ง view Access และ full Access เท่านั้น
	const checkViewAccessMiddleware = (viewAccessKey: IUserRolesAccess, fullAccessKey: IUserRolesAccess) => {
		const canView = findAccessKey(viewAccessKey.toString())
		const canEdit = findAccessKey(fullAccessKey.toString())

		if (!canView && !canEdit) {
			return showError({ statusCode: httpStatusCode.ACCESS_DENIED })
		}
	}

	// ส่ง AccessKey เป็น array, หากพบ AccessKey เพียง 1 จากทั้งหมดถือว่ามีสิทธิ์เข้าใช้
	const checkMultipleAccessMiddleware = (accessKeys: IUserRolesAccess[]) => {
		let canUsePage = false

		for (const index in accessKeys) {
			console.log("key ✅", accessKeys[index])
			canUsePage = findAccessKey(accessKeys[index].toString())
			if (canUsePage) {
				return canUsePage
			}
		}

		if (!canUsePage) {
			return showError({ statusCode: httpStatusCode.ACCESS_DENIED })
		}
	}

	// ส่ง full Access เท่านั้น
	const checkMenuAccessMiddleware = (accessKey: IUserRolesAccess) => {
		// กรณีไม่มีสิทธิ์เข้าถึงหน้าเว็บไซต์ 403

		// ถ้าไม่มี full Access ถือว่าไม่มีสิทธิ์เข้าหน้าแก้ไข
		if (!findAccessKey(accessKey.toString())) {
			return showError({ statusCode: httpStatusCode.ACCESS_DENIED })
		}
	}

	const checkCanAccessPermissionBYKey = (accessKey: IUserRolesAccess) => {
		return findAccessKey(accessKey.toString())
	}

	const hasPermission = (accessKey: string) => {
		console.log("hasPermission ✅", accessKey)
		return true
	}

	const setMiddleware = (accessKey: string) => {
		console.log("setMiddleware ✅", accessKey)
		return true
	}

	return {
		checkMultipleAccessMiddleware,
		checkViewAccessMiddleware,
		checkMenuAccessMiddleware,
		checkCanAccessPermissionBYKey,
		hasPermission,
		setMiddleware,
	}
}

export interface IAccessDetail {
	isTrue: string[]
	isFalse: string[]
}

export interface IAccessMap {
	[key: string]: IAccessDetail
}

const getPermissionRelation = (accessKey: string, value: boolean): string[] => {
	const relationObj: IAccessMap = {
		view_total_dashboard: {
			isTrue: [],
			isFalse: [],
		},
		view_all_asset_dashboard: {
			isTrue: [],
			isFalse: [],
		},
		view_owner_asset_dashboard: {
			isTrue: [],
			isFalse: [],
		},
		view_all_road_condition_dashboard: {
			isTrue: [],
			isFalse: [],
		},
		view_owner_road_condition_dashboard: {
			isTrue: [],
			isFalse: [],
		},
		view_all_surface_dashboard: {
			isTrue: [],
			isFalse: [],
		},
		view_owner_surface_dashboard: {
			isTrue: [],
			isFalse: [],
		},
		view_all_maint_history_dashboard: {
			isTrue: [],
			isFalse: [],
		},
		view_owner_maint_history_dashboard: {
			isTrue: [],
			isFalse: [],
		},
		view_all_road: {
			isTrue: ["view_road_summary"],
			isFalse: ["manage_all_road"],
		},
		view_owner_road: {
			isTrue: ["view_road_summary"],
			isFalse: ["manage_owner_road"],
		},
		manage_all_road: {
			isTrue: ["view_all_road"],
			isFalse: [],
		},
		manage_owner_road: {
			isTrue: ["view_owner_road"],
			isFalse: [],
		},
		view_road_summary: {
			isTrue: [],
			isFalse: ["manage_road_summary", "manage_owner_road_summary", "manage_road_traffic", "manage_owner_road_traffic"],
		},
		manage_road_summary: {
			isTrue: ["view_all_road"],
			isFalse: [],
		},
		manage_owner_road_summary: {
			isTrue: ["view_owner_road"],
			isFalse: [],
		},
		view_road_in_assets: {
			isTrue: [],
			isFalse: ["manage_road_in_assets", "manage_owner_road_in_assets"],
		},
		manage_road_in_assets: {
			isTrue: ["view_road_in_assets"],
			isFalse: [],
		},
		manage_owner_road_in_assets: {
			isTrue: ["view_road_in_assets"],
			isFalse: [],
		},
		view_road_condition: {
			isTrue: [],
			isFalse: ["manage_road_condition", "manage_owner_road_condition"],
		},
		manage_road_condition: {
			isTrue: ["view_road_condition"],
			isFalse: [],
		},
		manage_owner_road_condition: {
			isTrue: ["view_road_condition"],
			isFalse: [],
		},
		view_road_retro: {
			isTrue: [],
			isFalse: ["manage_road_retro", "manage_owner_road_retro"],
		},
		manage_road_retro: {
			isTrue: ["view_road_retro"],
			isFalse: [],
		},
		manage_owner_road_retro: {
			isTrue: ["view_road_retro"],
			isFalse: [],
		},
		view_road_damage: {
			isTrue: [],
			isFalse: ["manage_road_damage", "manage_owner_road_damage"],
		},
		manage_road_damage: {
			isTrue: ["view_road_damage"],
			isFalse: [],
		},
		manage_owner_road_damage: {
			isTrue: ["view_road_damage"],
			isFalse: [],
		},
		view_road_traffic: {
			isTrue: [],
			isFalse: ["manage_road_traffic", "manage_owner_road_traffic"],
		},
		manage_road_traffic: {
			isTrue: ["view_road_traffic"],
			isFalse: [],
		},
		manage_owner_road_traffic: {
			isTrue: ["view_road_traffic"],
			isFalse: [],
		},

		view_all_maint_history: {
			isTrue: [],
			isFalse: ["manage_all_maint_history"],
		},
		view_owner_maint_history: {
			isTrue: [],
			isFalse: ["manage_owner_maint_history"],
		},
		manage_all_maint_history: {
			isTrue: ["view_all_maint_history"],
			isFalse: [],
		},
		manage_owner_maint_history: {
			isTrue: ["view_owner_maint_history"],
			isFalse: [],
		},

		view_all_maintenance_analysis: {
			isTrue: [],
			isFalse: [],
		},
		view_owner_maintenance_analysis: {
			isTrue: [],
			isFalse: [],
		},
		manage_myself_maintenance_analysis: {
			isTrue: [],
			isFalse: [],
		},

		download_all_asset_report: {
			isTrue: [],
			isFalse: [],
		},
		download_owner_asset_report: {
			isTrue: [],
			isFalse: [],
		},
		download_all_map_asset_report: {
			isTrue: [],
			isFalse: [],
		},
		download_owner_map_asset_report: {
			isTrue: [],
			isFalse: [],
		},
		download_all_result_asset_report: {
			isTrue: [],
			isFalse: [],
		},
		download_owner_result_asset_report: {
			isTrue: [],
			isFalse: [],
		},
		download_all_road_report: {
			isTrue: [],
			isFalse: [],
		},
		download_all_surface_report: {
			isTrue: [],
			isFalse: [],
		},
		download_owner_surface_report: {
			isTrue: [],
			isFalse: [],
		},
		download_all_road_condition_report: {
			isTrue: [],
			isFalse: [],
		},
		download_owner_road_condition_report: {
			isTrue: [],
			isFalse: [],
		},
		download_all_road_condition_result_report: {
			isTrue: [],
			isFalse: [],
		},
		download_owner_road_condition_result_report: {
			isTrue: [],
			isFalse: [],
		},
		download_all_road_retro_report: {
			isTrue: [],
			isFalse: [],
		},
		download_owner_road_retro_report: {
			isTrue: [],
			isFalse: [],
		},
		download_all_road_retro_result_report: {
			isTrue: [],
			isFalse: [],
		},
		download_owner_road_retro_result_report: {
			isTrue: [],
			isFalse: [],
		},
		download_all_road_damage_report: {
			isTrue: [],
			isFalse: [],
		},
		download_owner_road_damage_report: {
			isTrue: [],
			isFalse: [],
		},
		download_all_road_damage_result_report: {
			isTrue: [],
			isFalse: [],
		},
		download_owner_road_damage_result_report: {
			isTrue: [],
			isFalse: [],
		},
		download_all_maint_history_report: {
			isTrue: [],
			isFalse: [],
		},
		download_owner_maint_history_report: {
			isTrue: [],
			isFalse: [],
		},
		download_all_maint_kpi_report: {
			isTrue: [],
			isFalse: [],
		},
		download_owner_maint_kpi_report: {
			isTrue: [],
			isFalse: [],
		},

		download_all_aadt_report: {
			isTrue: [],
			isFalse: [],
		},
		download_owner_aadt_report: {
			isTrue: [],
			isFalse: [],
		},

		setting_asset_place_group_access: {
			isTrue: [],
			isFalse: [],
		},
		setting_in_asset_access: {
			isTrue: [],
			isFalse: [],
		},
		setting_traffic_sign_access: {
			isTrue: [],
			isFalse: [],
		},
		setting_condition_criteria_access: {
			isTrue: [],
			isFalse: [],
		},
		setting_retro_criteria_access: {
			isTrue: [],
			isFalse: [],
		},
		setting_surface_type_access: {
			isTrue: [],
			isFalse: [],
		},
		setting_ref_budget_access: {
			isTrue: [],
			isFalse: [],
		},
		setting_model_access: {
			isTrue: [],
			isFalse: [],
		},
		setting_hsms_access: {
			isTrue: [],
			isFalse: [],
		},
		setting_hris_access: {
			isTrue: [],
			isFalse: [],
		},

		setting_user_acess: {
			isTrue: [],
			isFalse: [],
		},
		setting_permission_acess: {
			isTrue: [],
			isFalse: [],
		},
	}

	const access = relationObj[accessKey]

	if (!access) {
		return []
	}

	// console.log("accessKey:", accessKey, "value:", value)
	// console.log("isFalse:", access.isFalse)
	// console.log("isTrue:", access.isTrue)

	if (value) {
		return access.isTrue
	} else {
		return access.isFalse
	}
}

export { usePermission, getPermissionRelation, IUserRolesAccess }
