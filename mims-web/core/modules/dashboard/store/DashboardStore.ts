import {
	DashboardService,
	IDashboardAsset,
	IDashboardAssetData,
	IDashboardAssetDetailsRequest,
	IDashboardAssetRequest,
	IDashboardMapAssetRequest,
	IDashboardStateAssetDetail,
	IDashboardStateAssetItems,
	IDashboardRoad,
	IHighChart,
	IMapBound,
	IDashboardSurface,
	IDashboardSurfaceRequest,
	IDashboardSurfaceMapRequest,
	IDashboardSurfaceMap,
	IDataMartCheck,
	IDashboardMaintenanceRequest,
	IDashboardMaintenance,
	IDashboardMaintenanceMapRequest,
	IDashboardMaintenanceMap,
	IHandleCheckboxMode,
	IDashboardConditionMapRequest,
} from "../infrastructure"
import {
	useDashboardConditionStore,
	useDashboardReflectiveStore,
	useDashboardRoadConditionStore,
} from "../condition/store"
import { ITree } from "~/core/shared/types/Tree"
import { IOption } from "~/core/shared/types/Option"
import { IDashboardMapItems } from "~/core/modules/dashboard/condition/infrastructure/DashboardConditionModel"

const localStorageUser = window.localStorage.getItem("init-user")
const initUserStore = localStorageUser ? JSON.parse(localStorageUser) : { accessPermissions: {} }
interface IPavementAge {
	id: number
	name: string
	color: string
}

interface IMarkerData {
	markers: any
	data: { id: number; asset_table_id: number }
}

interface IStateData {
	assets: IDashboardAssetData
	assets_details: IDashboardStateAssetDetail[]
	condition: {}
	surface: IDashboardSurface
	surface_map: IDashboardSurfaceMap[]
	surface_colors: any[]
	surface_cateories: string[]
	surface_array: number[]
	maintenance: IDashboardMaintenance
	maintenance_quantity_index: number[]
	maintenance_budget_index: number[]
	maintenance_map: IDashboardMaintenanceMap[]
}

interface IStateParams {
	road_id: string[]
	depot_code: string[]
	km_start: string
	km_end: string
	year: number | null
	ref_asset_id: number[]
	asset_type: string
	display: number
	limit: number
	page: number
}

type IDashboardMenu = "asset" | "condition" | "surface" | "maintenance"

interface IState {
	menu: IDashboardMenu
	map: any
	longdo: any
	syncing: boolean
	dataMart: IDataMartCheck
	markers: IMarkerData[]
	arrAsset: number[]
	year: number
	display_type: number
	owners: string[]
	roads: string[]
	roadsOptions: ITree[]
	yearOptions: IOption[]
	roadsData: IDashboardRoad
	totalData: number
	dataArray: number[]
	conditions: any[]
	reflect: any[]
	data: IStateData
	params: IStateParams
	toggle: {
		hist: boolean
		whiteLine: boolean
		dataType: string
		ownerID: number
		conditionType: number
	}
	maintenanceTableSearch: Function
	graphSelectFilter: number[]
	dataGraph: IHighChart
	conditionArray: number[]
	conditionList: any[]
	conditionType: number
	surveyRule: number
	lane: number
	labelsArr: string[]
	colors: string[]
	displayParam: number
	loading: boolean
	pavementAge: IPavementAge[]
	delayMapFetch: any
	searchCall: boolean
	mapPage: number
	mapArray: IDashboardMapItems[]
	loadingMap: boolean
	isAssetEventCreated: boolean
	hasPermissionAccressTab: boolean
	isSingleRoad: boolean
	selectCondition: any
}

export const useDashboardStore = defineStore("dashboard", {
	state: (): IState => ({
		menu: "asset" as IDashboardMenu,
		map: null,
		longdo: null,
		syncing: false,
		dataMart: {} as IDataMartCheck,
		markers: [],
		arrAsset: [11, 12],
		year: 2567,
		owners: [],
		roads: [],
		roadsOptions: [],
		yearOptions: [],
		roadsData: {} as IDashboardRoad,
		totalData: 0,
		display_type: 1,
		data: {
			assets: {} as IDashboardAssetData,
			assets_details: [],
			condition: {},
			surface: {} as IDashboardSurface,
			surface_colors: [],
			surface_map: [],
			surface_cateories: [],
			surface_array: [],
			maintenance: {} as IDashboardMaintenance,
			maintenance_quantity_index: [],
			maintenance_budget_index: [],
			maintenance_map: [],
		},
		params: {
			road_id: [],
			depot_code: [],
			km_start: "",
			km_end: "",
			year: null,
			ref_asset_id: [],
			asset_type: "",
			display: 2,
			limit: 10,
			page: 1,
		},
		dataArray: [],
		conditions: [
			{
				id: 1,
				name: "ดีมาก",
				color: "#A4FCA5",
			},
			{
				id: 2,
				name: "ดี",
				color: "#42D235",
			},
			{
				id: 3,
				name: "ปานกลาง",
				color: "#F77A14",
			},
			{
				id: 4,
				name: "แย่",
				color: "#FF290A",
			},
			{
				id: 5,
				name: "แย่มาก",
				color: "#973131",
			},
		],
		reflect: [
			{
				id: 6,
				name: "ผ่าน",
				color: "#42D235",
			},
			{
				id: 7,
				name: "ไม่ผ่าน",
				color: "#FF290A",
			},
		],
		toggle: {
			hist: false,
			whiteLine: false,
			dataType: "IRI",
			ownerID: 1,
			conditionType: 1,
		},
		dataGraph: {
			title: { text: "" },
			chart: {
				width: null,
				height: 300,
				type: "line",
				zoomType: "x",
				events: {},
				stacked: true,
				toolbar: { show: false },
				resetZoomButton: {
					position: {
						y: -5,
						x: 10,
					},
					theme: {
						fill: "#fff",
						stroke: "#fdb833",
						style: {
							color: "#3f4254",
							fontWeight: "500",
							fontSize: "12px",
						},
						states: {
							hover: {
								fill: "#e8edf3",
								stroke: "#fdb833",
								style: {
									color: "#3f4254",
								},
							},
						},
						r: 13,
						padding: 9,
					},
				},
			},
			dataLabels: { enabled: false },
			stroke: { curve: "smooth" },
			xAxis: {
				categories: [
					"7+400",
					"7+375",
					"7+350",
					"7+325",
					"7+300",
					"7+275",
					"7+250",
					"7+225",
					"7+200",
					"7+175",
					"7+150",
					"7+125",
					"7+100",
					"7+075",
					"7+050",
					"7+025",
					"7+000",
					"6+975",
					"6+950",
				],
				tickInterval: 0,
				margin: 10,
				title: { text: "Kilometer (km)", style: { fontSize: "12.5px", color: "#3f4254" } },
				labels: { enable: false, style: { fontSize: "12px", color: "#3f4254" } },
			},
			yAxis: {
				title: { text: "IRI (ม./กม.)", style: { fontSize: "12.5px", color: "#3f4254" } },
				labels: { enable: false, style: { fontSize: "12px", color: "#3f4254" } },
			},
			tooltip: {
				useHTML: true,
				borderRadius: 25,
				backgroundColor: "#fff",
				borderColor: "none",
				padding: 0,
				formatter: function (this: any): string {
					return `<div class="px-4 py-2 fs-7">${this.x}</div>
						<div style="border-top: 1px solid #f3f3f3; width: 100%;" display: block;></div>
						<div class="px-4 py-2 d-flex">
							<span class="rounded me-3 mt-1 fs-7" style="width: 12px; height: 12px; display: block; background: ${this.color};"></span>
							<span>IRI (ม./กม.): <b class="fw-semibold">${this.y}</b></span>
						</div>
						`
				},
			},
			plotOptions: {
				bar: {
					dataLabels: {
						position: "top", // top, center, bottom
					},
				},
				boost: {
					allowForce: true,
					// useGPUTranslations: true,
					enabled: true,
					seriesThreshold: 100,
					useGPUTranslations: true,
					usePreAllocated: true,
				},
				series: {
					bootsThreshold: 100,
					animation: false,
					marker: {
						enabled: true,
					},
					point: {
						events: {},
					},
				},
			},
			grid: {
				borderColor: "#EAEAEA",
				strokeDashArray: 5,
				xaxis: {
					lines: {
						show: false,
					},
				},
				yaxis: {
					lines: {
						show: true,
					},
				},
			},
			legend: {
				enabled: false,
			},
			fill: {
				opacity: 1,
			},
			credits: {
				enabled: false,
			},
			series: [
				{
					name: "IRI",
					data: [
						2.76, 2.66, 2.93, 3.18, 3.02, 3.2, 2.78, 2.87, 2.61, 2.72, 2.73, 2.65, 2.69, 2.31, 2.71, 2.72, 2.23, 2.56,
						2.57,
					],
					zoneAxis: "y",
					zones: [
						{ value: 2.75, color: "#42d235" },
						{ value: 3.5, color: "#a4fca5" },
						{ value: 4.5, color: "#f77a14" },
						{ value: 15, color: "#ff290a" },
					],
					// animation: false,
				},
			],
		} as IHighChart,
		graphSelectFilter: [1, 2, 3, 4, 5],
		conditionArray: [1, 2, 3, 4, 5],
		conditionList: [],
		conditionType: 1,
		surveyRule: 1,
		lane: 1,
		labelsArr: [],
		colors: [],
		displayParam: 1,
		loading: false,
		pavementAge: [
			{ name: "0 - 2 ปี", color: "#50CD89" } as IPavementAge,
			{ name: "3 - 5 ปี", color: "#87C442" } as IPavementAge,
			{ name: "6 - 10 ปี", color: "#FDB833" } as IPavementAge,
			{ name: "มากกว่า 10 ปี", color: "#DC3545" } as IPavementAge,
		],
		delayMapFetch: null,
		maintenanceTableSearch: () => {},
		searchCall: false,
		mapPage: 1,
		mapArray: [],
		loadingMap: false,
		isAssetEventCreated: false,
		hasPermissionAccressTab: true,
		isSingleRoad: false,
		selectCondition: {
			status: false,
			color: "",
		},
	}),
	actions: {
		setMenuByPermission() {
			if (
				initUserStore.accessPermissions[IUserRolesAccess.view_all_asset_dashboard] ||
				initUserStore.accessPermissions[IUserRolesAccess.view_owner_asset_dashboard]
			) {
				this.menu = "asset"
			} else if (
				initUserStore.accessPermissions[IUserRolesAccess.view_all_road_condition_dashboard] ||
				initUserStore.accessPermissions[IUserRolesAccess.view_owner_road_condition_dashboard]
			) {
				this.menu = "condition"
			} else if (
				initUserStore.accessPermissions[IUserRolesAccess.view_all_surface_dashboard] ||
				initUserStore.accessPermissions[IUserRolesAccess.view_owner_surface_dashboard]
			) {
				this.menu = "surface"
			} else if (
				initUserStore.accessPermissions[IUserRolesAccess.view_all_maint_history_dashboard] ||
				initUserStore.accessPermissions[IUserRolesAccess.view_owner_maint_history_dashboard]
			) {
				this.menu = "maintenance"
			}

			this.hasPermissionAccressTab =
				initUserStore.accessPermissions[IUserRolesAccess.view_all_asset_dashboard] ||
				initUserStore.accessPermissions[IUserRolesAccess.view_owner_asset_dashboard] ||
				initUserStore.accessPermissions[IUserRolesAccess.view_all_road_condition_dashboard] ||
				initUserStore.accessPermissions[IUserRolesAccess.view_owner_road_condition_dashboard] ||
				initUserStore.accessPermissions[IUserRolesAccess.view_all_surface_dashboard] ||
				initUserStore.accessPermissions[IUserRolesAccess.view_owner_surface_dashboard] ||
				initUserStore.accessPermissions[IUserRolesAccess.view_all_maint_history_dashboard] ||
				initUserStore.accessPermissions[IUserRolesAccess.view_owner_maint_history_dashboard]
		},
		resetParams() {
			this.params.asset_type = ""
			this.params.ref_asset_id = []
			this.params.km_start = ""
			this.params.km_end = ""
			this.params.year = null
			this.params.depot_code = []
			this.params.road_id = []
		},
		getFirstRoadIdsAndDepotCodesWhenOwnerAccess() {
			let isViewOwner =
				initUserStore.accessPermissions[IUserRolesAccess.view_owner_asset_dashboard] &&
				!initUserStore.accessPermissions[IUserRolesAccess.view_all_asset_dashboard]

			switch (this.menu) {
				case "asset":
					isViewOwner =
						initUserStore.accessPermissions[IUserRolesAccess.view_owner_asset_dashboard] &&
						!initUserStore.accessPermissions[IUserRolesAccess.view_all_asset_dashboard]
					break
				case "condition":
					isViewOwner =
						initUserStore.accessPermissions[IUserRolesAccess.view_owner_road_condition_dashboard] &&
						!initUserStore.accessPermissions[IUserRolesAccess.view_all_road_condition_dashboard]
					break
				case "surface":
					isViewOwner =
						initUserStore.accessPermissions[IUserRolesAccess.view_owner_surface_dashboard] &&
						!initUserStore.accessPermissions[IUserRolesAccess.view_all_surface_dashboard]
					break
				case "maintenance":
					isViewOwner =
						initUserStore.accessPermissions[IUserRolesAccess.view_owner_maint_history_dashboard] &&
						!initUserStore.accessPermissions[IUserRolesAccess.view_all_maint_history_dashboard]
					break

				default:
					break
			}

			const isFirstRoadIdsForViewOwnerAccess = isViewOwner && (!this.params.road_id || this.params.road_id.length === 0)
			const isFirstDepotCodeForViewOwnerAccess =
				isViewOwner && (!this.params.depot_code || this.params.depot_code.length === 0)

			const roadIds = isFirstRoadIdsForViewOwnerAccess
				? this.roadsOptions.flatMap(
						(parent) =>
							parent.children?.flatMap((section) => section.children?.map((child) => child.id ?? "") ?? []) ?? []
				  )
				: this.params.road_id

			const depotCodes = isFirstDepotCodeForViewOwnerAccess
				? this.getOwnerOptions.flatMap(
						(parent) =>
							parent.children?.flatMap((section) => section.children?.map((child) => child.id ?? "") ?? []) ?? []
				  )
				: this.params.depot_code

			return { roadIds, depotCodes }
		},
		async getRoadDropdown() {
			const service = new DashboardService()
			const res = await service.getRoadDropdown(this.menu)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				const options: ITree[] = res.data?.map((item) => ({
					label: `${item.road_number} ${item.short_name}`,
					id: `parent-${item.id}`,
					children: item.road_sections?.map((section) => ({
						label: `${section.number} ${section?.name_origin} - ${section?.name_destination}`,
						id: `section-${section?.id}`,
						children: section.roads?.map((road) => ({
							label: road?.name,
							id: `${road?.id}`,
						})),
					})),
				}))

				this.roadsOptions = options || []
			}
		},
		async getAssets() {
			const { data, params } = this

			const { roadIds, depotCodes } = this.getFirstRoadIdsAndDepotCodesWhenOwnerAccess()

			const req: IDashboardAssetRequest = {
				road_id: roadIds,
				depot_code: depotCodes,
				km_start: Number.isNaN(convertStringToKm(params.km_start)) ? null : convertStringToKm(params.km_start),
				km_end: Number.isNaN(convertStringToKm(params.km_end)) ? null : convertStringToKm(params.km_end),
				year: params.year,
			}

			const service = new DashboardService()
			const res = await service.getDashboardAsset(req)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				data.assets = res.data
			}
		},
		async getAssetDetails() {
			const { params } = this

			const { roadIds, depotCodes } = this.getFirstRoadIdsAndDepotCodesWhenOwnerAccess()

			const req: IDashboardAssetDetailsRequest = {
				road_id: roadIds,
				ref_asset_id: params.ref_asset_id,
				asset_type: params.asset_type,
				depot_code: depotCodes,
				km_start: Number.isNaN(convertStringToKm(params.km_start)) ? null : convertStringToKm(params.km_start),
				km_end: Number.isNaN(convertStringToKm(params.km_end)) ? null : convertStringToKm(params.km_end),
				year: params.year,
			}

			const service = new DashboardService()
			const res = await service.getDashboardAssetDetails(req)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.data.assets_details = res.data?.items?.map((parent) => ({
					id: parent?.asset_group?.id,
					name: parent?.asset_group?.name,
					colors: "#FDB833",
					is_active: true,
					items: parent?.asset_list?.map((item) => ({
						id: item?.asset?.id,
						name: item.asset?.name,
						icon_path: item.asset?.default_icon_url,
						thumbnail_path: item.asset?.thumbnail_icon_url,
						is_active: true,
						value: item?.value,
						color: "#FDB833",
						is_range: item.is_range,
					})),
				}))

				this.params.ref_asset_id = this.data.assets_details?.flatMap((parent) => parent.items.map((item) => item.id))
			}
		},
		setMap(map: any) {
			this.map = map
			if (this.map) {
				// @ts-ignore
				this.longdo = window.longdo
			}

			const self = this

			// first time call
			self.createAssetInMap()
			self.createEventAssetByKey()
			self.createEventAssetByMapToolsbar()

			// Debounce pan/zoom to prevent rapid consecutive API calls (race condition)
			const debouncedCreateAsset = () => {
				if (self.delayMapFetch) clearTimeout(self.delayMapFetch)
				self.delayMapFetch = setTimeout(() => {
					self.createAssetInMap()
				}, 300)
			}

			self.map?.Event.bind("drop", debouncedCreateAsset)
			this.map.Event.bind("zoom", debouncedCreateAsset)
		},

		createAssetInMap() {
			const { map, getAssetsMap, getConditionMap, getMaintenanceMap } = this
			if (map) {
				if (!this.selectCondition.status) {
					// Surface lines are loaded once (not viewport-bound) — do not clear on pan/zoom
					if (this.menu === "surface") return

					map.Overlays.clear()
					if (this.menu === "asset") {
						const bound: IMapBound = map.bound()
						const zoomLv: number = map.zoom()
						getAssetsMap(bound, zoomLv)
					} else if (this.menu === "condition") {
						getConditionMap()
					} else if (this.menu === "maintenance") {
						getMaintenanceMap()
					}
				}
			}
		},
		async getDashboardSurface() {
			const { roadIds, depotCodes } = this.getFirstRoadIdsAndDepotCodesWhenOwnerAccess()
			const params: IDashboardSurfaceRequest = {
				road_id: roadIds,
				depot_code: depotCodes,
				km_start: Number.isNaN(convertStringToKm(this.params.km_start))
					? null
					: convertStringToKm(this.params.km_start),
				km_end: Number.isNaN(convertStringToKm(this.params.km_end)) ? null : convertStringToKm(this.params.km_end),
				year: this.params.year,
			}

			const service = new DashboardService()

			try {
				const res = await service.getDashboardSurface(params)
				if (res.status === false) {
					useHandlerError(res.code, res.error, { showAlert: true })
				} else {
					this.data.surface = res.data
					this.data.surface_array = []
					this.colors = []
					this.data.surface_colors = this.data.surface?.summary.map((item) => item.surface.color_code)
					this.data.surface_cateories =
						this.data.surface?.summary?.map((item: any, index) => {
							this.data.surface_array.push(index)
							return item.surface.name
						}) ?? []
					this.dataArray = this.data.surface?.summary?.map((item: any) => Number(item.value.toFixed(2)))
					return res
				}
			} catch (error) {
				console.error(error)
			}
		},
		async syncDataMart() {
			// Loading
			this.syncing = true

			const service = new DashboardService()
			const check = await service.getDataMartCheck()
			const isNotSyncing = check.data?.stauts
			if (!isNotSyncing) {
				return
			}

			setTimeout(() => {
				this.checkSyncDataMart()
			}, 3000)
			this.$patch((state) => {
				if (state.dataMart) {
					state.dataMart.percent = 0
				}
			})
			const res = await service.getDataMart()

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		async getSurfaceMap() {
			const { roadIds, depotCodes } = this.getFirstRoadIdsAndDepotCodesWhenOwnerAccess()
			const params: IDashboardSurfaceMapRequest = {
				road_id: roadIds,
				depot_code: depotCodes,
				km_start: Number.isNaN(convertStringToKm(this.params.km_start))
					? null
					: convertStringToKm(this.params.km_start),
				km_end: Number.isNaN(convertStringToKm(this.params.km_end)) ? null : convertStringToKm(this.params.km_end),
				year: this.params.year,
				display: this.params.display,
			}

			const service = new DashboardService()
			const res = await service.getSurfaceMap(params)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.data.surface_map = res.data

				const fn = () => {
					if (this.map) {
						this.map.Overlays.clear()
						// Sort so red lines are drawn last and always appear on top
						const sorted = [...(res.data || [])].sort(
							(a, b) => getColorDrawPriority(a.color) - getColorDrawPriority(b.color)
						)
						sorted.forEach((item) => this.createSurfaceLine(item, item.color))
					} else {
						setTimeout(() => {
							fn()
						}, 500)
					}
				}
				fn()
			}
		},
		createSurfaceLine(item: IDashboardSurfaceMap, color = "") {
			const { map, longdo } = this

			if (map) {
				const detail = this.generatePopupDetails(item)

				const geoms = item.the_geom?.coordinates?.map((coord) => ({ lon: coord[0], lat: coord[1] }))
				const line = new longdo.Polyline(geoms, { lineColor: color, detail })

				map.Overlays.add(line)
			}
		},
		async getConditionMap() {
			this.loadingMap = true

			if (this.mapArray.length) {
				this.mapArray = []
			}

			// Guard against race condition: only render results from the latest request
			const requestId = Date.now()
			;(this as any)._conditionRequestId = requestId

			const params = this.setConditionMapParams()
			const dashboardConditionService = new DashboardService()

			try {
				const res = await dashboardConditionService.getCondition_map(params)

				if ((this as any)._conditionRequestId !== requestId) return

				if (!res.status) {
					useHandlerError(res.code, res.error, { showToast: true })
					return
				}

				const totalPages = res.data?.total_pages || 1

				if (totalPages > 1) {
					const promises = Array.from({ length: totalPages }, (_, index) => {
						return dashboardConditionService.getCondition_map({
							...params,
							page: index + 1,
						})
					})

					const responses = await Promise.all(promises)

					if ((this as any)._conditionRequestId !== requestId) return

					this.mapArray = responses.flatMap((response) => response.data.items)
				} else {
					this.mapArray = res.data.items || []
				}

				this.map.Overlays.clear()
				this.createLineCondition()
			} catch (error) {
				console.error("Error fetching condition map data:", error)
			} finally {
				if ((this as any)._conditionRequestId === requestId) {
					this.loadingMap = false
				}
			}
		},

		async getTotalPages(params: IDashboardConditionMapRequest) {
			const dashboardConditionService = new DashboardService()
			const res = await dashboardConditionService.getCondition_map(params)
			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
				return 0
			}
			return res.data.total_pages
		},

		setConditionMapParams() {
			const initData = useInitData()
			const dashboardCondition = useDashboardRoadConditionStore()
			const dashboardReflectiveStore = useDashboardReflectiveStore()
			const conditionStore = useDashboardConditionStore()
			const { roadIds, depotCodes } = this.getFirstRoadIdsAndDepotCodesWhenOwnerAccess()

			// Fallback owner must match the current condition type (IRI/IFI/MPD/RUT).
			// conditionGrade()[0] may not have rules for IFI/MPD/RUT, so we find the
			// first grade owner that actually supports the current condition type string.
			const conditionTypeStr = dashboardCondition.params.condition_type // e.g. "IFI"
			const ownerId =
				initData
					.conditionGrade()
					.find((g) => g.condition_list?.some((c: any) => c.condition_type === conditionTypeStr))?.id ??
				(initData.conditionGrade().length > 0 ? initData.conditionGrade()[0]?.id : null)

			let refConditionRangeId: number | null = 0
			refConditionRangeId =
				conditionStore.conditionType === 5
					? dashboardReflectiveStore.params.owner_id
					: dashboardCondition.params.owner_id

			const kmStart = Number.isNaN(convertStringToKm(this.params.km_start))
				? null
				: convertStringToKm(this.params.km_start)
			const kmEnd = Number.isNaN(convertStringToKm(this.params.km_end)) ? null : convertStringToKm(this.params.km_end)

			const year = Number(this.params.year)

			const bound: IMapBound = this.map.bound()

			const params: IDashboardConditionMapRequest = {
				condition_type: conditionStore.conditionType,
				condition_owner_id: refConditionRangeId || ownerId,
				road_id: roadIds,
				year,
				depot_code: depotCodes,
				km_start: kmStart,
				km_end: kmEnd,
				left: bound.minLon,
				bottom: bound.minLat,
				right: bound.maxLon,
				top: bound.maxLat,
				limit: 10000,
				page: 1,
			}

			return params
		},
		createLineCondition() {
			if (this.map) {
				// @ts-ignore
				const longdo = window.longdo
				const items = this.map === null ? [] : this.mapArray

				// Sort by draw priority so red lines are added last and appear on top
				const sorted = [...items].sort(
					(a, b) => getColorDrawPriority(a.color) - getColorDrawPriority(b.color)
				)

				sorted.forEach((item) => {
					const strLine = longdo.Util.overlayFromGeoJSON(item.the_geom, {
						lineColor: item.color === "#000000" ? convertHexToRGBA(item.color, 0.5) : convertHexToRGBA(item.color, 1),
					})
					this.map.Overlays.add(strLine[0])
				})
			}
		},

		createLineConditionByColor() {
			if (this.map) {
				// @ts-ignore
				const longdo = window.longdo
				let latLng = {
					lon: 0,
					lat: 0,
				}

				const items = this.mapArray.filter((item: any) => item.color === this.selectCondition.color)

				this.map.Overlays.clear()

				items.forEach((item) => {
					const strLine = longdo.Util.overlayFromGeoJSON(item.the_geom, {
						lineColor: item.color === "#000000" ? convertHexToRGBA(item.color, 0.5) : convertHexToRGBA(item.color, 1),
					})
					this.map.Overlays.add(strLine[0])
				})

				if (items.length > 0) {
					latLng = {
						lon: items[0].the_geom.coordinates[0][0],
						lat: items[0].the_geom.coordinates[1][1],
					}

					this.map.location({
						lon: latLng.lon,
						lat: latLng.lat,
					})
				}
			}
		},

		async checkSyncDataMart() {
			const service = new DashboardService()

			const res = await service.getDataMartCheck()

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
				return false
			} else {
				if (res.data.stauts) {
					this.getDashboardSurface()
					// this.getRoad()
					this.getSurfaceMap()
					setTimeout(() => {
						this.syncing = false
					}, 1000)
				} else {
					this.syncing = true
					setTimeout(() => {
						this.checkSyncDataMart()
					}, 5000)
				}
				this.$patch((state) => {
					state.dataMart = res.data
				})

				return !!res.data.stauts
			}
		},

		createLine<T extends IDashboardMaintenanceMap>(item: T, color = "") {
			const { map, longdo } = this
			if (map) {
				let detail = ""

				const geoms = item.the_geom?.coordinates?.map((coord) => ({ lon: coord[0], lat: coord[1] })) || []

				detail = this.createPopupDetail<T>(item)

				const line = new longdo.Polyline(geoms, { lineColor: color, detail })

				map.Overlays.add(line)

				// Add event listener for navigate button after line is added
				// Use longer timeout and multiple attempts to find the button
				// const setupButtonListener = (attempts = 0) => {
				// 	console.log(`Attempt ${attempts + 1} to find button for item:`, item.id_parent)

				// 	let navigateButton = null

				// 	// Try to find button in document (popup might be in document, not in line element)
				// 	navigateButton = document.querySelector(`#navigate-maintenances-${item.id_parent}`)
				// 	console.log("Button found in document:", navigateButton)

				// 	// Also try to find in any iframe or shadow DOM
				// 	if (!navigateButton) {
				// 		// Check all iframes
				// 		const iframes = document.querySelectorAll('iframe')
				// 		for (const iframe of iframes) {
				// 			try {
				// 				const iframeDoc = iframe.contentDocument || iframe.contentWindow?.document
				// 				if (iframeDoc) {
				// 					navigateButton = iframeDoc.querySelector(`#navigate-maintenances-${item.id_parent}`)
				// 					if (navigateButton) {
				// 						console.log("Button found in iframe:", navigateButton)
				// 						break
				// 					}
				// 				}
				// 			} catch (e) {
				// 				// Cross-origin iframe, skip
				// 			}
				// 		}
				// 	}

				// 	if (navigateButton) {
				// 		console.log("Adding click event listener to button")
				// 		navigateButton.addEventListener("click", (e) => {
				// 			e.preventDefault()
				// 			console.log("Button clicked! Navigating to:", `/maintenances/history/${(item as any).id_parent}/info`);

				// 			// Use navigateTo with proper context
				// 			if (typeof navigateTo === 'function') {
				// 				console.log("Using navigateTo function")
				// 				navigateTo(`/maintenances/history/${(item as any).id_parent}/info`)
				// 			} else {
				// 				console.log("navigateTo not available, using window.location")
				// 				// Fallback to window.location if navigateTo is not available
				// 				window.location.href = `/maintenances/history/${(item as any).id_parent}/info`
				// 			}
				// 		})
				// 	} else if (attempts < 10) {
				// 		// Retry after longer delay
				// 		console.log("Button not found, retrying in 500ms...")
				// 		setTimeout(() => setupButtonListener(attempts + 1), 500)
				// 	} else {
				// 		console.log("Button not found after 10 attempts! ID:", `navigate-maintenances-${item.id_parent}`)
				// 	}
				// }

				// // Start looking for the button
				// setTimeout(() => setupButtonListener(), 100)
			}
		},
		generatePopupDetails(item: IDashboardSurfaceMap) {
			return `
					<div class="row mb-3" style="
						width: 220px;" >
						${this.getProcessHeaderTitle(item)}
				 	</div>
					<div class="row">
						<div class="col-4">สายทาง:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${item.road_group_name}</div>
					</div>
					<div class="row mb-2">
						<div class="col-4">เลขที่สัญญา:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${item.contract_number || "-"}</div>
					</div>

					<div class="row mb-2">
						<div class="col-4">ปีงบประมาณ:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${item.year && isNumber(item.year) ? Number(item.year) + 543 : "-"}</div>
					</div>
					<div class="row mb-2">
						<div class="col-4">วันที่ซ่อมล่าสุด:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${buddhistFormatDate(item.last_inspection_date, "dd mmm yyyy") || "-"} </div>
					</div>
					<div class="row mb-2 mt-2">
						<div class="col-4">จาก - ถึง:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${item.road_name} </div>
					</div>
					<div class="row mb-2 mt-2">
						<div class="col-4">กม. เริ่มต้น - กม. สิ้นสุด:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${this.getKmStartEnd(item)} </div>
					</div>
					<div class="row mb-2 mt-2">
						<div class="col-4">ระยะทาง:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${item.km_total.toFixed(2)} กม.</div>
					</div>
					<div class="row mb-2 mt-2">
						<div class="col-4">ผิวทางล่าสุด:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${item.surface_name} </div>
					</div>`
		},
		getProcessHeaderTitle(item: IDashboardSurfaceMap) {
			const icon = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
			<g clip-path="url(#clip0_10409_34188)">
			<path d="M20.2759 2.48C20.1536 1.78514 19.7904 1.15561 19.2499 0.702047C18.7095 0.248483 18.0265 -9.24445e-05 17.3209 4.06091e-06H6.68094C5.97464 -0.0011587 5.29067 0.247414 4.74993 0.701785C4.20918 1.15616 3.84648 1.78707 3.72594 2.483L-0.0390625 24H24.0609L20.2759 2.48ZM12.9999 21H10.9999V17H12.9999V21ZM12.9999 14H10.9999V10H12.9999V14ZM12.9999 7.00001H10.9999V3H12.9999V7.00001Z" fill="${item.color}"/>
			</g>
			<defs>
			<clipPath id="clip0_10409_34188">
			<rect width="24" height="24" fill="white"/>
			</clipPath>
			</defs>
			</svg>`

			return `<div class="col-1">
			${icon}
			</div>
			<div class="col"
			style="
			color: ${item.color};
			font-size: 16px;
			font-weight: 500;
			">${item.title}</div>`
		},
		getKmStartEnd(road: IDashboardSurfaceMap) {
			return `${convertMeterToKm(road.km_start)} - ${convertMeterToKm(road.km_end)}`
		},
		createPopupDetail<T extends IDashboardMaintenanceMap>(data: T) {
			if (data as IDashboardMaintenanceMap) {
				return `
        <div class="row">
          <div class="col-1 align-items-center mb-1">
            <svg width="24" height="24" viewBox="0 0 24 25" fill="none" xmlns="http://www.w3.org/2000/svg">
              <g clip-path="url(#clip0_4337_22762)">
                <path fill="${
									data.color
								}" d="M12 0.5C9.62663 0.5 7.30655 1.20379 5.33316 2.52236C3.35977 3.84094 1.8217 5.71509 0.913451 7.9078C0.00519943 10.1005 -0.232441 12.5133 0.230582 14.8411C0.693605 17.1689 1.83649 19.3071 3.51472 20.9853C5.19295 22.6635 7.33115 23.8064 9.65892 24.2694C11.9867 24.7324 14.3995 24.4948 16.5922 23.5866C18.7849 22.6783 20.6591 21.1402 21.9776 19.1668C23.2962 17.1935 24 14.8734 24 12.5C23.9966 9.31846 22.7312 6.26821 20.4815 4.01852C18.2318 1.76883 15.1815 0.503441 12 0.5ZM13 12.379C13.0001 12.5485 12.9571 12.7153 12.8751 12.8636C12.793 13.0119 12.6746 13.137 12.531 13.227L8.69101 15.627C8.57938 15.6967 8.45512 15.7438 8.3253 15.7655C8.19549 15.7872 8.06268 15.7831 7.93444 15.7535C7.8062 15.7239 7.68505 15.6693 7.57791 15.5929C7.47077 15.5164 7.37974 15.4196 7.31001 15.308C7.24027 15.1964 7.19321 15.0721 7.1715 14.9423C7.1498 14.8125 7.15387 14.6797 7.18349 14.5514C7.21311 14.4232 7.2677 14.302 7.34414 14.1949C7.42059 14.0878 7.51738 13.9967 7.62901 13.927L11 11.825V7.5C11 7.23478 11.1054 6.98043 11.2929 6.79289C11.4804 6.60536 11.7348 6.5 12 6.5C12.2652 6.5 12.5196 6.60536 12.7071 6.79289C12.8946 6.98043 13 7.23478 13 7.5V12.379Z" fill="#1F70F3"/>
              </g>
              <defs>
                <clipPath id="clip0_4337_22762">
                  <rect width="24" height="24" fill="white" transform="translate(0 0.5)"/>
                </clipPath>
              </defs>
            </svg>
          </div>
          <span class="col fs-5 pe-0" style="color: ${data.color}">${data.title}</span>
        </div>
        <div class="row px-0">
          <div class="col-5 pe-0">เลขที่สัญญาโครงการ</div>
          <div class="col colon">: </div>
          <div class="col pe-0">${data.contract_number}</div>
        </div>
        <div class="row px-0">
          <div class="col-5 pe-0">สายทาง</div>
          <div class="col colon">: </div>
          <div class="col pe-0">${data.road_name}</div>
        </div>
        <div class="row px-0">
          <div class="col-5 pe-0">หน่วยงาน</div>
          <div class="col colon">: </div>
          <div class="col pe-0">${data.ref_depot_name}</div>
        </div>
        <div class="row px-0">
          <div class="col-5 pe-0">จาก - ถึง</div>
          <div class="col colon">: </div>
          <div class="col pe-0">${data.section_name}</div>
        </div>
        <div class="row px-0">
          <div class="col-5 pe-0">กม. เริ่มต้น - กม. สิ้นสุด</div>
          <div class="col colon">: </div>
          <div class="col pe-0">${data.km_start} - ${data.km_end}</div>
        </div>
        <div class="row px-0">
          <div class="col-5 pe-0">ช่องจราจร</div>
          <div class="col colon">: </div>
          <div class="col pe-0">${data.lane_no ?? "-"}</div>
        </div>
        <div class="row px-0">
          <div class="col-5 pe-0">ระยะทาง</div>
          <div class="col colon">: </div>
          <div class="col pe-0">${(data.km_total / 1000).toFixed(2)} กม. | ${data.id_parent}</div>
        </div>
		<div class="text-end">
			<button class="btn btn-primary py-0" id="navigate-maintenances-${data.id_parent}">ดูรายละเอียด</button>
		</div>
        `
			}
			return ""
		},
		async getAssetMapDetails(id: number, itemId: number) {
			const service = new DashboardService()
			const res = await service.getAssetsMapDetails(id, itemId)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				return res.data
			}
		},
		async getAssetsMap(bound: IMapBound, zoomLevel: number) {
			const { params } = this

			// reset markers
			this.markers = []

			const { roadIds, depotCodes } = this.getFirstRoadIdsAndDepotCodesWhenOwnerAccess()

			const req: IDashboardMapAssetRequest = {
				depot_code: depotCodes,
				road_id: roadIds,
				km_start: Number.isNaN(convertStringToKm(params.km_start)) ? null : convertStringToKm(params.km_start),
				km_end: Number.isNaN(convertStringToKm(params.km_end)) ? null : convertStringToKm(params.km_end),
				year: params.year,
				ref_asset_id: params.ref_asset_id,
				left: bound.minLon,
				bottom: bound.minLat,
				right: bound.maxLon,
				top: bound.maxLat,
				zoom: zoomLevel,
			}

			const service = new DashboardService()
			const res = await service.getDashboardMapAssets(req)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				const data = res.data

				this.map.Overlays.clear()
				data.forEach((item) => {
					this.generateAssets(item)
				})
			}
		},
		// <a class="btn btn-primary py-0" href="/roads/${id}/in-asset/${asset_table_id}">ดูรายละเอียด</a>

		async addNavigateToDetail(item: IDashboardAsset) {
			const detail = await this.getAssetMapDetails(item.id, item.asset_table_id)

			return `
				${detail} 
				<div class="text-end">
					<button class="btn btn-primary py-0" id="navigate-btn-${item.id}-${item.asset_table_id}">ดูรายละเอียด</button>
				</div>
			`
		},

		generateAssets(item: IDashboardAsset) {
			const addPopup = async (lon: number, lat: number, item: IDashboardAsset) => {
				const detail = await this.addNavigateToDetail(item)

				if (!item.is_cluster) {
					const popup = new this.longdo.Popup(
						{ lon, lat },
						{
							detail,
							closable: true,
						}
					)

					const popupDetailElement = popup.element().querySelector(".ldmap_popup_detail")
					popupDetailElement.classList.add("pe-0", "w-100")

					// Add event listener for navigate button
					const navigateButton = popup.element().querySelector(`#navigate-btn-${item.id}-${item.asset_table_id}`)
					if (navigateButton) {
						navigateButton.addEventListener("click", () => {
							navigateTo(`/roads/${item.id}/in-asset/${item.asset_table_id}`)
						})
					}

					this.map.Overlays.add(popup)
				} else {
					this.map.location({
						lon,
						lat,
					})

					setTimeout(() => {
						if (this.map.zoom() < 11) {
							this.map.zoom(11)
						}
					}, 200)
					setTimeout(
						() => {
							if (this.map.zoom() < 13) {
								this.map.zoom(13)
							}
						},
						this.map.zoom() >= 11 ? 200 : 1000
					)
				}
			}

			const addEventListeners = (element: any, lon: number, lat: number, item: IDashboardAsset) => {
				element?.addEventListener("click", () => addPopup(lon, lat, item))
				element?.classList.add("cursor-pointer")
			}

			if (item.the_geom?.type === "Point") {
				const geoms = item.the_geom
				const lon = geoms.coordinates[0] as number
				const lat = geoms.coordinates[1] as number

				const options = {
					icon: this.setMarkerIcon(item),
				}

				const marker = new this.longdo.Marker({ lon, lat }, options)
				addEventListeners(marker.element(), lon, lat, item)
				this.map.Overlays.add(marker)
			} else {
				const geoms = item.the_geom
				const lon = (geoms.coordinates[0] as number[])[0]
				const lat = (geoms.coordinates[0] as number[])[1]
				const coordinates = item.the_geom?.coordinates.map((coord) => ({
					lon: (coord as number[])[0],
					lat: (coord as number[])[1],
				}))

				const line = new this.longdo.Polyline(coordinates, {
					lineColor: item.line_color_code,
					clickable: true,
				})

				this.map.Event.bind("overlayClick", (overlay: any) => {
					if (overlay === line) {
						if (!item.is_cluster) {
							addPopup(lon, lat, item)
						}
					}
				})

				this.map.Overlays.add(line)
			}
		},

		async getYearOptions() {
			const service = new DashboardService()
			const res = await service.getDashboardYear()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				const options = res.data.map((item) => ({ label: `${item + 543}`, value: item }))

				this.yearOptions = options || []
			}
		},
		async initial() {
			const { getRoadDropdown, getYearOptions } = this
			this.loading = true

			await getRoadDropdown()
			await getYearOptions()

			switch (this.menu) {
				case "asset":
					await this.getAssets()
					await this.getAssetDetails()
					// Refresh map after assets are loaded: fixes race condition where map loads
					// before getAssetDetails() completes (params.ref_asset_id was empty on map load)
					if (this.map) {
						this.createAssetInMap()
					}
					break
				case "surface":
					await this.getDashboardSurface()
					await this.getSurfaceMap()
					this.checkSyncDataMart()
					break
			case "condition": {
				const dashboardConditionInit = useDashboardConditionStore()
				dashboardConditionInit.getCondition()
				await this.getConditionMap()
				break
			}
				case "maintenance":
					await this.getMaintenanceDashboard()
					await this.getMaintenanceMap()
					break
				default:
					await this.getAssets()
					await this.getAssetDetails()
					break
			}

			this.loading = false
		},
		async onSearch() {
			this.loading = true
			const dashboardCondition = useDashboardConditionStore()
			const dashboardRoadCondition = useDashboardRoadConditionStore()
			const storeDashboardReflectiveStore = useDashboardReflectiveStore()

			this.map.Overlays.clear()

			switch (this.menu) {
				case "asset":
					await this.getAssets()
					await this.getAssetDetails()
					this.createAssetInMap()

					break
				case "surface":
					await this.getDashboardSurface()
					await this.getSurfaceMap()
					this.checkSyncDataMart()
					break
				case "condition":
					dashboardRoadCondition.map = this.map
					storeDashboardReflectiveStore.map = this.map

					if (this.params.road_id.length === 1) {
						await dashboardRoadCondition.getRoads()
					} else {
						dashboardCondition.getCondition()
						await this.getConditionMap()
					}

					break
				case "maintenance":
					await this.getMaintenanceDashboard()
					await this.getMaintenanceMap()
					break
				default:
					break
			}
			this.loading = false
		},
		async getRoad() {
			this.loading = true

			const service = new DashboardService()
			const res = await service.getDashboardRoads()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.roadsData = res.data
			}
		},
		async getMaintenanceDashboard() {
			// reset index
			this.data.maintenance_budget_index = []
			this.data.maintenance_quantity_index = []

			const { roadIds, depotCodes } = this.getFirstRoadIdsAndDepotCodesWhenOwnerAccess()

			const params: IDashboardMaintenanceRequest = {
				road_id: roadIds,
				depot_code: depotCodes,
				km_start: Number.isNaN(convertStringToKm(this.params.km_start))
					? null
					: convertStringToKm(this.params.km_start),
				km_end: Number.isNaN(convertStringToKm(this.params.km_end)) ? null : convertStringToKm(this.params.km_end),
				year: this.params.year,
				// limit: this.params.limit,
				// page: this.params.page,
			}

			const service = new DashboardService()
			const res = await service.getMaintenance(params)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.data.maintenance = res.data
				res.data?.number_maintenance_chart?.lable.forEach((_, index) =>
					this.data.maintenance_quantity_index.push(index)
				)

				res.data?.maintenance_budget_chart.lable.forEach((_, index) => this.data.maintenance_budget_index.push(index))
			}
		},
		async getMaintenanceMap() {
			const { roadIds, depotCodes } = this.getFirstRoadIdsAndDepotCodesWhenOwnerAccess()
			const params: IDashboardMaintenanceMapRequest = {
				road_id: roadIds,
				depot_code: depotCodes,
				km_start: Number.isNaN(convertStringToKm(this.params.km_start))
					? null
					: convertStringToKm(this.params.km_start),
				km_end: Number.isNaN(convertStringToKm(this.params.km_end)) ? null : convertStringToKm(this.params.km_end),
				year: this.params.year,
			}

			const service = new DashboardService()
			const res = await service.getMaintenanceMap(params)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.data.maintenance_map = res.data
				res.data?.forEach((item) => this.createLine(item, item.color))
			}
		},
		async handleMenu(menu: IDashboardMenu) {
			const dashboardCondition = useDashboardConditionStore()
			const dashboardRoadCondition = useDashboardRoadConditionStore()

			if (menu === this.menu) {
				return
			}

			this.menu = menu
			this.loading = true
			this.map.Overlays.clear()

			await this.getRoadDropdown()

			switch (menu) {
				case "asset":
					await this.getAssets()
					await this.getAssetDetails()
					this.createAssetInMap()
					if (!this.isAssetEventCreated) {
						this.createEventAssetByKey()
						this.createEventAssetByMapToolsbar()
					}
					break
				case "surface":
					await this.checkSyncDataMart()
					break
				case "condition":
					this.map.zoom(12)
					const store = useDashboardRoadConditionStore()
					const storeDashboardReflectiveStore = useDashboardReflectiveStore()

					await nextTick()

					store.map = this.map
					storeDashboardReflectiveStore.map = this.map

					dashboardRoadCondition.map = this.map
					storeDashboardReflectiveStore.map = this.map

					if (this.params.road_id.length === 1) {
						await dashboardRoadCondition.getRoads()
					} else {
						dashboardCondition.getCondition()
					}
					this.isSingleRoad = this.params.road_id.length === 1

					break
				case "maintenance":
					await this.getMaintenanceDashboard()
					await this.getMaintenanceMap()

					break
				default:
					break
			}
			this.loading = false

			if (this.menu !== "condition") {
				this.defaultLocation()
			}
		},
		setMarkerIcon(item: IDashboardAsset) {
			const checkIconPath = item.icon_filepath !== "" && item.icon_filepath
			const checkThumbnailPath = item.thumbnail_icon_filepath !== "" && item.thumbnail_icon_filepath
			if (item.is_cluster) {
				return {
					html: `
          <div style=" width: 40px;
            height: 40px;
            background-color: ${convertHexToRGBA("#FDB833", 0.5)};
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            position: relative;"
          >
            <span style="
              color: #000;
              font-size: 12px;
              position: relative;
              z-index: 2;
              "
            >
              ${item.cluster}
            </span>
          </div>`,
				}
			} else {
				return checkIconPath
					? this.setIcon(item.icon_filepath)
					: checkThumbnailPath
					? this.setIcon(item.thumbnail_icon_filepath)
					: this.setIcon("/images/icons/png/location-pin.png")
			}
		},
		generateImage(item: IDashboardStateAssetItems) {
			const checkIconPath = item.icon_path !== "" && item.icon_path
			const checkThumbnailPath = item.thumbnail_path !== "" && item.thumbnail_path
			return checkIconPath
				? item.icon_path
				: checkThumbnailPath
				? item.thumbnail_path
				: item.is_active
				? "/images/icons/svg/location-pin.svg"
				: "/images/icons/svg/location-pin-grey.svg"
		},
		createEventAssetByKey() {
			const mapElement = document.querySelector(".longdo-map")

			if (this.menu === "asset") {
				if (mapElement) {
					mapElement.addEventListener("keyup", (event: any) => {
						const arrowKeys = ["ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight", "w", "a", "s", "d"]
						if (arrowKeys.includes(event.key)) {
							this.map.Overlays.clear()
							this.createAssetInMap()
						}
					})

					this.isAssetEventCreated = true
				}
			}
		},
		createEventAssetByMapToolsbar() {
			const toolbars = ["ldmap_goleft", "ldmap_goright", "ldmap_goup", "ldmap_godown"]

			if (this.menu === "asset") {
				toolbars.forEach((tool) => {
					const toolElement = document.querySelector(`.${tool}`)

					if (toolElement) {
						toolElement.addEventListener("click", () => {
							this.map.Overlays.clear()
							this.createAssetInMap()
						})
					}
				})
			}
		},

		setIcon(img: string) {
			const imgElement = document.createElement("img")
			imgElement.src = img
			imgElement.style.width = "25px"
			imgElement.style.height = "25px"
			imgElement.alt = "icon-image"

			return {
				html: imgElement.outerHTML,
				offset: { x: 10.5, y: 18.5 },
			}
		},
		// condition
		updateConditionChart(pieChart: Ref, barChart: Ref) {
			this.totalData = this.dataArray?.reduce((acc: any, item: any) => acc + item)
			pieChart.value.updateOptions({
				series: this.dataArray,
				colors: this.colors,
				labels: this.labelsArr,
			})
			barChart.value.updateOptions({
				colors: this.colors,
				xaxis: {
					categories: this.labelsArr,
				},
			})
		},
		updatehisChart(histChart: Ref): any {
			histChart.value.updateOptions({
				colors: this.colors,
				xaxis: {
					categories: this.labelsArr,
				},
			})
		},
		barSeries(): any {
			if (this.dataArray) {
				const result = this.data.surface_array
					.map((id: number) => this.dataArray[id])
					.filter((value: any) => value !== null && value !== undefined)
				return result.map((item: any) => Number(item))
			}
			return []
		},
		histSeries(): any {
			const series = [{ name: "IRI", data: this.dataArray }]
			return series
		},
		toggleGraph(event: any) {
			const innerType = event.target.innerHTML.toLowerCase()
			if (innerType === "histogram") {
				this.toggle.hist = true
			} else {
				this.toggle.hist = false
			}
		},

		toggleReflect(event: any) {
			const innerType = event.target.innerHTML.toLowerCase()
			if (innerType === "เส้นสีเหลือง") {
				this.toggle.whiteLine = true
			} else {
				this.toggle.whiteLine = false
			}
		},
		getCheckBox() {
			switch (this.conditionType) {
				case 1:
					this.graphSelectFilter = [1, 2, 3, 4]
					this.conditionArray = [1, 2, 3, 4]
					this.dataArray = [44, 55, 13, 32]
					this.conditionList = [
						{
							id: 1,
							label: "เรียบมาก",
							value: 1,
							color: "#42D235",
						},
						{
							id: 2,
							label: "เรียบ",
							value: 2,
							color: "#A4FCA5",
						},
						{
							id: 3,
							label: "ขรุขระ",
							value: 3,
							color: "#F77A14",
						},
						{
							id: 4,
							label: "ขรุขระมาก",
							value: 4,
							color: "#FF290A",
						},
					]
					break
				case 2:
					this.graphSelectFilter = [1, 2, 3]
					this.conditionArray = [1, 2, 3]
					this.dataArray = [44, 55, 13]
					this.conditionList = [
						{
							id: 1,
							label: "หยาบ (ดีมาก)",
							value: 1,
							color: "#A4FCA5",
						},
						{
							id: 2,
							label: "ปานกลาง",
							value: 2,
							color: "#F77A14",
						},
						{
							id: 3,
							label: "ละเอียด (แย่มาก)",
							value: 3,
							color: "#FF290A",
						},
					]
					break
				case 3:
					this.graphSelectFilter = [1, 2, 3, 4]
					this.conditionArray = [1, 2, 3, 4]
					this.dataArray = [44, 55, 13, 32]
					this.conditionList = [
						{
							id: 1,
							label: "ตื้นมาก",
							value: 1,
							color: "#42D235",
						},
						{
							id: 2,
							label: "ตื้น",
							value: 1,
							color: "#A4FCA5",
						},
						{
							id: 3,
							label: "ลึก",
							value: 3,
							color: "#F77A14",
						},
						{
							id: 4,
							label: "ลึกมาก",
							value: 4,
							color: "#FF290A",
						},
					]
					break
				case 4:
				case 5:
					this.graphSelectFilter = [6, 7]
					this.conditionArray = [1, 2]
					this.dataArray = [44, 23]
					this.conditionList = [
						{
							id: 1,
							label: "ผ่าน",
							value: 6,
							color: "#42D235",
						},
						{
							id: 2,
							label: "ไม่ผ่าน",
							value: 7,
							color: "#FF290A",
						},
					]
					break
			}
			this.colors = this.conditionList.map((item) => item.color)
			this.labelsArr = this.conditionList.map((item) => item.label)
		},
		handleSelect(id: number) {
			const index = this.params.ref_asset_id.indexOf(id)

			if (index > -1) {
				this.params.ref_asset_id.splice(index, 1)
			} else {
				this.params.ref_asset_id.push(id)
			}
		},
		onUpdateData(refDonut: Ref) {
			if (!refDonut?.value) return
			refDonut.value.updateOptions({
				chart: {
					events: {
						dataPointSelection: (_: any, __: any, opts: any) => {
							const dataIndex = opts.dataPointIndex
							const roadGroupId = this.roadsData?.road?.road_group_id

							if (roadGroupId.length) {
								navigateTo(`roads?road_group_id=${roadGroupId[dataIndex]}`)
							}
						},
					},
				},
				label: this.getQuantityRoadSeries?.labels,
				colors: this.getRoadsLegend.colors,
			})
		},
		onChecked() {
			const { data, params } = this
			const assetDetails = data?.assets_details

			assetDetails.forEach((parent) => {
				parent.items.forEach((item) => {
					if (params.ref_asset_id.includes(item.id)) {
						item.is_active = true
					} else {
						item.is_active = false
					}
				})

				parent.is_active = parent.items.some((item) => item.is_active)
			})

			this.createAssetInMap()
		},
		setParentActive(item: IDashboardStateAssetDetail) {
			item.is_active = !item.is_active
			item.items.forEach((child) => {
				if (item.is_active) {
					this.params.ref_asset_id.push(child.id)
				} else {
					const index = this.params.ref_asset_id.indexOf(child.id)
					if (index > -1) {
						this.params.ref_asset_id.splice(index, 1)
					}
				}

				child.is_active = item.is_active
			})
			this.params.ref_asset_id = [...new Set(this.params.ref_asset_id)]
			this.createAssetInMap()
		},
		setChildActive(item: IDashboardStateAssetItems) {
			const { data } = this
			item.is_active = !item.is_active

			if (item.is_active) {
				this.params.ref_asset_id.push(item.id)
			} else {
				const index = this.params.ref_asset_id.indexOf(item.id)
				if (index > -1) {
					this.params.ref_asset_id.splice(index, 1)
				}
			}

			this.params.ref_asset_id = [...new Set(this.params.ref_asset_id)]

			const assetDetail = data.assets_details

			if (assetDetail.length) {
				const parentMap = new Map()

				assetDetail.forEach((parent) => {
					parent.items?.forEach((child) => {
						parentMap?.set(child.id, parent)
					})
				})

				const parent = parentMap.get(item.id)
				if (parent) {
					parent.is_active = parent.items?.some((child: IDashboardStateAssetItems) => child.is_active)
				}
			}

			this.createAssetInMap()
		},
		barOptions(): any {
			const categories = ref<Array<string>>([])
			if (this.data.surface_cateories) {
				categories.value = this.data.surface_array.map((id: number) => this.data.surface_cateories[id]).filter(Boolean)
			}

			const option: any = {
				chart: {
					events: {
						dataPointSelection: (_: any, __: any, opts: any) => {
							const dataIndex = opts.dataPointIndex
							const roadId = this.data?.surface?.summary[dataIndex]?.road_id
							const surfaceId = this.data?.surface?.summary[dataIndex]?.surface?.id
							// const roadGroupId = this.roadsData?.road?.road_group_id

							navigateTo(`roads?road_id=${roadId}&ref_surface_id=${surfaceId}`)
						},
					},
				},
				plotOptions: {
					pie: {
						dataLabels: {
							offset: -5,
						},
					},
				},
				title: {
					text: "ข้อมูลสรุปผิวทาง",
					align: "center",
					style: {
						fontSize: "16px",
					},
				},
				dataLabels: {
					enabled: true,
					formatter: (_: number, { seriesIndex, w }: any) => {
						return w.config.series[seriesIndex] + " กม."
					},
					style: {
						fontSize: "10px",
						colors: ["#fff"],
					},
				},
				colors: this.colors?.length === 0 ? this.data.surface_colors : this.colors,
				labels: categories.value,
				tooltip: {
					enabled: true,
					y: {
						show: true,
						formatter: (value: number) => {
							if (value) {
								return value.toLocaleString() + " กม."
							} else {
								return 0 + " กม."
							}
						},
					},
				},
				legend: {
					show: false,
				},
			}
			return option
		},
		defaultLocation() {
			const { map } = this
			if (map) {
				map.location({
					lon: 100.546876,
					lat: 13.739992,
				})
				map.zoom(8)
			}
		},
		handleCheckbox(index: number, name: IHandleCheckboxMode, refChart: Ref) {
			const indexMap: Record<IHandleCheckboxMode, number[]> = {
				"maintenance-quantity": this.data.maintenance_quantity_index,
				"maintenance-budget": this.data.maintenance_budget_index,
			}

			const targetIndexArray = indexMap[name]

			if (targetIndexArray) {
				const idx = targetIndexArray.indexOf(index)
				if (idx !== -1) {
					targetIndexArray.splice(idx, 1)
				} else {
					targetIndexArray.push(index)
				}
			}

			const sumMap: Record<IHandleCheckboxMode, number> = {
				"maintenance-quantity": this.getSumMaintenanceChart?.sum_quantity,
				"maintenance-budget": this.getSumMaintenanceChart?.sum_budget,
			}

			const label = sumMap[name] ? toNumber(sumMap[name]) : "0"

			refChart.value?.updateOptions({
				plotOptions: {
					pie: {
						donut: {
							labels: {
								total: {
									label,
								},
							},
						},
					},
				},
			})
		},
	},

	getters: {
		getSurfaceGrade() {
			const refSurface = useInitData().refSurface()
			if (refSurface) {
				return refSurface
			}

			return []
		},
		getAssetCategories(state) {
			const data = state?.data
			const assets = data?.assets?.items
			const label = assets?.map((item) => item.name)

			return label || []
		},
		getAssetSeries(state) {
			const data = state?.data
			const assets = data?.assets?.items

			const labels = [...new Set(assets?.flatMap((item) => item.label))]

			const series = labels?.flatMap((label, index) => {
				return {
					name: label,
					data: assets?.flatMap((item) => item.data[index] ?? 0),
				}
			})

			return series || []
		},
		getAssetItemsOptions(state) {
			const data = state?.data
			const assets = data?.assets

			const options = assets?.items?.map((item) => ({ label: "", value: item.id }))

			return options || []
		},
		getOwnerOptions(state) {
			let divisions = useInitData()?.refDivisionDashboardAsset()

			switch (state?.menu) {
				case "asset":
					divisions = useInitData()?.refDivisionDashboardAsset()
					break

				case "condition":
					divisions = useInitData()?.refDivisionDashboardCondition()
					break

				case "surface":
					divisions = useInitData()?.refDivisionDashboardSurface()
					break

				case "maintenance":
					divisions = useInitData()?.refDivisionDashboardMaintenance()
					break

				default:
					divisions = useInitData()?.refDivisionDashboardAsset()
					break
			}

			const options: ITree[] =
				divisions?.map((division) => ({
					label: division.name,
					id: division.owner_code_key,
					children: division.districts?.map((district) => ({
						label: district.name,
						id: district.owner_code_key,
						children: district.depots?.map((depot) => ({
							label: depot.name,
							id: depot.depot_code,
						})),
					})),
				})) || []

			return options || []
		},
		getQuantityRoadSeries(state) {
			const roadsData = state?.roadsData
			const roads = roadsData?.road

			const series = roads?.data ?? []
			const labels = roads?.label ?? []
			const options = {
				series,
				labels,
			}

			return options || ({} as keyof typeof options)
		},
		getRoadsLenght(state) {
			const roadsData = state?.roadsData
			const roadsLengths = roadsData?.length_roads

			return roadsLengths || []
		},
		getRoadsAadt(state) {
			const roadsData = state?.roadsData
			const roadsAadt = roadsData?.aadt_roads ?? []

			roadsAadt.forEach((road) => {
				road.year1 = road.year1 + 543
				road.year2 = road.year2 + 543
			})

			return roadsAadt || []
		},
		getMaintnenanceQuantity(state) {
			const data = state?.data
			const maintenance = data?.maintenance
			const dataIndex = data?.maintenance_quantity_index ?? []
			const label = maintenance?.number_maintenance_chart?.lable?.filter((_, index) => dataIndex.includes(index)) ?? []
			const maintenanceData =
				maintenance?.number_maintenance_chart?.data?.filter((_, index) => dataIndex.includes(index)) ?? []
			const colors = maintenance?.number_maintenance_chart?.color?.filter((_, index) => dataIndex.includes(index)) ?? []
			return {
				label,
				series: maintenanceData,
				colors,
			}
		},
		getMaintenanceQuantityCheckbox(state) {
			const data = state?.data
			const maintenance = data?.maintenance
			const labels = maintenance?.number_maintenance_chart?.lable || []
			const colors = maintenance?.number_maintenance_chart?.color || []
			const result = labels.map((label, index) => {
				return {
					name: label,
					color: colors[index],
				}
			})

			return result || []
		},

		getMaintenanceBudget(state) {
			const data = state?.data
			const maintenance = data?.maintenance
			const dataIndex = data?.maintenance_budget_index ?? []
			const label = maintenance?.maintenance_budget_chart?.lable?.filter((_, index) => dataIndex.includes(index)) ?? []
			const maintenanceData =
				maintenance?.maintenance_budget_chart?.data?.filter((_, index) => dataIndex.includes(index)) ?? []
			const colors = maintenance?.maintenance_budget_chart?.color?.filter((_, index) => dataIndex.includes(index)) ?? []
			return {
				label,
				series: maintenanceData,
				colors,
			}
		},
		getMaintenanceBudgetCheckbox(state) {
			const data = state?.data
			const maintenance = data?.maintenance
			const labels = maintenance?.maintenance_budget_chart?.lable || []
			const colors = maintenance?.maintenance_budget_chart?.color || []
			const result = labels.map((label, index) => {
				return {
					name: label,
					color: colors[index],
				}
			})

			return result || []
		},
		getSumMaintenanceChart(state) {
			const data = state?.data
			const dataMap: Record<IHandleCheckboxMode, number[] | undefined> = {
				"maintenance-quantity": data?.maintenance?.number_maintenance_chart?.data,
				"maintenance-budget": data?.maintenance?.maintenance_budget_chart?.data,
			}
			const indexMap: Record<IHandleCheckboxMode, number[]> = {
				"maintenance-quantity": data?.maintenance_quantity_index ?? [],
				"maintenance-budget": data?.maintenance_budget_index ?? [],
			}

			const quantity = dataMap["maintenance-quantity"]?.filter((_, index) =>
				indexMap["maintenance-quantity"]?.includes(index)
			)

			const budget = dataMap["maintenance-budget"]?.filter((_, index) => indexMap["maintenance-budget"]?.includes(index))

			return {
				sum_quantity: quantity?.length ? quantity?.reduce((acc, curr) => acc + curr) : 0,
				sum_budget: budget?.length ? budget?.reduce((acc, curr) => acc + curr) : 0,
			}
		},
		getMaintenanceBarChart(state) {
			const data = state?.data
			const maintenance = data?.maintenance?.top_ten_maintenance_budget_chart
			const categories = maintenance?.lable || []
			const series = maintenance?.data || []
			const colors = maintenance?.color || []

			return { categories, series: [{ name: "งบประมาณ", data: series }], colors }
		},
		getUpdateMaintenanceDate(state) {
			const data = state?.data
			const date = data?.maintenance?.updated_at

			return buddhistFormatDate(date, "dd/mm/yyyy") ?? "-"
		},
		getRoadsLegend(state) {
			const roadsData = state?.roadsData
			const labels = roadsData?.road?.label ?? []
			const colors = roadsData?.road?.color ?? []
			const result = labels.map((label, index) => {
				return {
					name: label,
					color: colors[index],
				}
			})

			return { legend: result, colors }
		},
	},
})
