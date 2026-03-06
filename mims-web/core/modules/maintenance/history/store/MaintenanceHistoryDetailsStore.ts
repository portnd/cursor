import {
	IMaintenanceHistoryDetailData,
	IMaintenancePlanListData,
	IMaintenanceHistoryPlanStatusData,
	IMaintenanceHistoryPlanStatusSchedule,
	IMaintenanceBudgetCriteria,
	IMaintenanceHistoryDetailRoadInfo,
	IMaintenanceHistoryDetailMaintenanceRoadHistoryGeom,
	IMaintenanceHistoryDetailMaintenanceRoad,
	IMaintenanceHistoryDetailMaintenanceRoadHistory,
} from "../infrastructure/MaintenanceHistoryModel"
import { MaintenanceHistoryService } from "../infrastructure/MaintenanceHistoryService"
import { IOption } from "~/core/shared/types/Option"

interface IRoadGeom {
	geom: IMaintenanceHistoryDetailMaintenanceRoadHistoryGeom
	color: string
	road: any
}

// interface IMaintenanceHistoryDetailMaintenanceRoadHistoryGeom {
// 	type: string;
// 	coordinates: number[][];
// }

interface IParamsState {
	planName: string | null
}

interface IState {
	loading: boolean
	data: IMaintenanceHistoryDetailData
	roadGeoms: IRoadGeom[]
	statusList: IMaintenanceHistoryPlanStatusData[]
	planList: IMaintenancePlanListData[]
	schedule: IMaintenanceHistoryPlanStatusSchedule[]
	budgetCriteriaList: IMaintenanceBudgetCriteria[]
	isShowMethod: boolean
	map: any
	longdo: any
	statusName: string
	params: IParamsState
	fullScreen: boolean
}

export const useMaintenanceHistoryDetailsStore = defineStore("maintenance/history/details", {
	state: (): IState => ({
		loading: false,
		data: { ref_depot: { id: 0 } } as IMaintenanceHistoryDetailData,
		roadGeoms: [],
		statusList: [],
		planList: [],
		schedule: [],
		budgetCriteriaList: [],
		isShowMethod: false,
		map: null,
		longdo: null,
		statusName: "",
		params: {
			planName: null,
		},
		fullScreen: false,
	}),
	actions: {
		async getMaintenanceHistoryDetail(id: number) {
			this.loading = true

			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceHistoryDetails(id)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.data = res.data
				this.setRoadGeom()
				// Defer map drawing so content paints first
				if (typeof requestIdleCallback !== "undefined") {
					requestIdleCallback(() => this.createLine(), { timeout: 100 })
				} else {
					setTimeout(() => this.createLine(), 0)
				}
				// await this.getStatusListOptions(id)
				// await this.getPlanStatus(id)
			}
		},
		async getStatusListOptions(id: number) {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceStatusOptions(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.planList = res.data
				this.checkIsCurrent()
			}
		},
		async getPlanStatus(id: number) {
			this.loading = true

			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenancePlanStatusList(id)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.statusList = res.data
			}
		},
		async getMaintenanceBudgets() {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceBudgetCriteria()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.budgetCriteriaList = res.data
			}
		},
		checkIsShowMethod(methodId: number) {
			this.budgetCriteriaList.forEach((item) => {
				item.budget_methods.forEach((child) => {
					if (methodId === child.id) {
						this.isShowMethod = child.is_show_method
					}
				})
			})
		},
		setMap(map: any) {
			this.map = map
			// @ts-ignore
			this.longdo = window.longdo
			this.map.Event.bind("fullscreen", this.fullscreen)
			this.createLine()
		},
		fullscreen() {
			const element = document.querySelector("#map")
			const map = element?.querySelector(".longdo-map")
			const fullScreenElement = map?.querySelector(".ldmap_placeholder_fullscreen") as HTMLElement
			if (fullScreenElement !== null) {
				this.fullScreen = true
			} else {
				this.fullScreen = false
			}
		},
		createLine() {
			if (this.map || this.roadGeoms.length > 0) {
				const items = this.roadGeoms === null ? [] : this.roadGeoms

				this.map.Overlays.clear()

				// @ts-ignore
				const longdo = window.longdo

				items.forEach((item) => {
					const strLine = longdo.Util.overlayFromGeoJSON(item.geom, {
						lineColor: item.color === "#000000" ? convertHexToRGBA(item.color, 0.5) : convertHexToRGBA(item.color, 1),
						detail: this.generatePopupDetails(item.road),
					})
					this.map.Overlays.add(strLine[0])
				})

				const geom = items[0]?.road?.the_geom.coordinates.length > 0 ? items[0].road.the_geom.coordinates[0] : undefined

				if (geom) {
					this.map.location({
						lon: geom[0],
						lat: geom[1],
					})
					this.map.zoom(12)

					setTimeout(() => {
						const latLon = getLatLongObject(items[0].road.the_geom)
						const popup = new longdo.Popup(latLon, {
							detail: this.generatePopupDetails(items[0].road),
						})
						this.map.Overlays.add(popup)
					}, 100)
				}
			}
		},
		toLocation(item: any) {
			if (item) {
				const geom = item.the_geom

				const latLon = getLatLongObject(geom)
				this.map.location({ lon: latLon.lon, lat: latLon.lat })
				this.map.zoom(12)

				setTimeout(() => {
					// @ts-ignore
					const longdo = window.longdo
					const popup = new longdo.Popup(latLon, {
						detail: this.generatePopupDetails(item),
					})
					this.map.Overlays.add(popup)
				}, 100)
			}
		},
		setSchedulePlan(planName: string) {
			this.statusList.forEach((item) => {
				if (item.name === planName) {
					this.schedule = item.schedules
				}
			})
		},
		checkIsCurrent() {
			if (this.planList.length > 0) {
				this.planList.forEach((item) => {
					if (item.is_current) {
						this.params.planName = item.name
						this.setSchedulePlan(this.params.planName)
					}
				})
			}
		},
		setRoadGeom() {
			if (Object.keys(this.data).length > 0) {
				const data = this.data

				const maintenanceHistoryGeom = data.road_histories.map((item) => ({
					geom: item.the_geom,
					color: item.color,
					road: item,
				}))

				const maintenanceRoadsGeom = data.roads.map((item) => ({
					geom: item.the_geom,
					color: item.color,
					road: item,
				}))

				this.roadGeoms = maintenanceRoadsGeom.concat(maintenanceHistoryGeom)
			}
		},
		setDateColor(guaranteeDate: Date) {
			const currentDate = new Date()
			const dateToCompare = new Date(guaranteeDate)

			// คำนวนความต่างของเวลา
			const diffTime = dateToCompare.getTime() - currentDate.getTime()

			// คำนวนความต่างของวัน
			const diffDays = diffTime / (1000 * 3600 * 24)

			// คำนวนความต่างของเดือน
			const diffMonths = diffDays / 30

			if (diffDays < 0) {
				return ""
			} else if (diffMonths <= 3) {
				return "text-danger"
			} else if (diffMonths <= 6) {
				return "text-warning"
			} else {
				return "text-success"
			}
		},
		generatePopupDetails(
			item: IMaintenanceHistoryDetailRoadInfo | IMaintenanceHistoryDetailMaintenanceRoad | IMaintenanceHistoryDetailMaintenanceRoadHistory
		) {
			return `
			<div class="row mb-3" >
				${this.getProcessHeaderTitle()}
		 	</div>
			<div class="row">
				<div class="col-6">ชื่อโครงการ:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${this.data.name}</div>
			</div>
			<div class="row">
				<div class="col-6">เลขที่สัญญาโครงการ:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${this.data.contract_number}</div>
			</div>
			<div class="row mb-2">
				<div class="col-6">สายทาง:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${this.data.road_group_names}</div>
			</div>
			<div class="row mb-2">
				<div class="col-6">จาก - ถึง:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${item.road_name}</div>
			</div>
			<div class="row mb-2">
				<div class="col-6">กม. เริ่มต้น - กม. สิ้นสุด:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${this.getKmStartEnd(item)}</div>
			</div>
			<div class="row mb-2">
				<div class="col-6">ช่องจราจรที่:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${item.lane_no}</div>
			</div>
			<div class="row mb-2">
				<div class="col-6">วิธีการซ่อมบำรุง:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${this.getInterventionCriteriaName(item)}</div>
			</div>
			<div class="row mb-2">
				<div class="col-6">ตรวจรับงานงวดสุดท้าย:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${buddhistFormatDate(this.data.project_end_date, "dd mmm yyyy")}</div>
			</div>
			<div class="row mb-2">
				<div class="col-6">หมดการค้ำประกัน:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${buddhistFormatDate(this.data.guarantee_expiration_date, "dd mmm yyyy")}</div>
			</div>
			<div class="row mt-2 pt-2 maintenance-popup-footer" style="border-top: 1px solid #E4E6EF;">
				<div class="col-12">
					<a href="/maintenances/history"
						class="maintenance-popup-view-details d-inline-flex align-items-center gap-1 text-primary text-decoration-none fw-semibold"
						data-maintenance-back="true"
						style="font-size: 13px; cursor: pointer; transition: opacity 0.2s;">
						กลับหน้าประวัติการซ่อมบำรุง
						<i class="fi fi-sr-angle-left" style="font-size: 12px;"></i>
					</a>
				</div>
			</div>
      `
		},
		getKmStartEnd(items: IMaintenanceHistoryDetailRoadInfo) {
			return items ? `${convertMeterToKm(items.km_start)} - ${convertMeterToKm(items.km_end)}` : "-"
		},
		/**
		 * ดึงชื่อวิธีการซ่อมบำรุงจาก road/road_history (มี intervention_criteria) หรือ road_info (ไม่มี)
		 * item จาก roadGeoms คือ full object จาก roads/road_histories ที่มี intervention_criteria
		 */
		getInterventionCriteriaName(
			item: IMaintenanceHistoryDetailRoadInfo | IMaintenanceHistoryDetailMaintenanceRoad | IMaintenanceHistoryDetailMaintenanceRoadHistory
		): string {
			const withCriteria = item as IMaintenanceHistoryDetailMaintenanceRoad | IMaintenanceHistoryDetailMaintenanceRoadHistory
			return withCriteria?.intervention_criteria?.maintenance_standard_name ?? "-"
		},
		getProcessHeaderTitle() {
			return `<div class="col-1">

				<i class="fi fi-sr-clock" style="font-size: 23px; color: ${this.data.color};"></i>
			</div>
			<div class="col"
			style="
			color: ${this.data.color};
			font-size: 16px;
			font-weight: 500;
			">${this.data.remaining_time}</div>`
		},
	},
	getters: {
		getRoadName(state) {
			if (Object.keys(state.data).length === 0) {
				return ""
			}

			let name = ""

			state.data.maintenance_roads.forEach((item) => {
				name = item.road_group.name
			})

			return name
		},
		getBudgetName(state) {
			if (Object.keys(state.data).length === 0) {
				return ""
			}

			const data = state.data
			const budgetData = data.budget ?? undefined
			const budgetName = budgetData?.name

			return budgetName ?? ""
		},
		getBudgetMethodName(state) {
			if (Object.keys(state.data).length === 0) {
				return ""
			}

			const data = state.data
			const budgetData = data.budget_method ?? undefined
			const budgetMethodName = budgetData?.method_name

			return budgetMethodName ?? ""
		},
		getMaintenanceHistoryRoadsTable(state) {
			// if (state.data.maintenance_roads?.length === 0 || Object.keys(state.data).length === 0) {
			if (state.data.roads?.length === 0 || Object.keys(state.data).length === 0) {
				return []
			}

			const data = state.data
			// const maintenanceRoads = data.maintenance_roads
			const maintenanceRoads = data.roads ?? []

			return maintenanceRoads
		},
		getSumMaintenanceRoadsTable(state) {
			// if (state.data.maintenance_roads?.length === 0 || Object.keys(state.data).length === 0) {
			if (state.data.roads?.length === 0 || Object.keys(state.data).length === 0) {
				return 0
			}

			const data = state.data
			// const maintenanceRoads = data.maintenance_roads
			const maintenanceRoads = data.roads
			const distances = maintenanceRoads.map((item) => Math.abs(item.km_end - item.km_start))
			const sumDistance = distances.reduce((acc, curr) => acc + curr) / 1000

			return +sumDistance.toFixed(2)
		},
		getMaintenanceHistoryDataTable(state) {
			if (state.data.road_histories?.length === 0 || Object.keys(state.data).length === 0) {
				return []
			}

			const data = state.data
			const maintenanceHistory = data.road_histories ?? []

			return maintenanceHistory ?? []
		},
		getSumMaintenanceHistoryDataTable(state) {
			if (state.data.road_histories?.length === 0 || Object.keys(state.data).length === 0) {
				return 0
			}

			const data = state.data
			const maintenanceHistory = data.road_histories
			const distances = maintenanceHistory.map((item) => Math.abs(item.km_end - item.km_start))
			const sumDistance = distances.reduce((acc, curr) => acc + curr) / 1000

			return +sumDistance.toFixed(2)
		},
		getRoadGroupId(state) {
			if (Object.keys(state?.data)?.length === 0) {
				return 0
			}

			const data = state.data.maintenance_roads
			const roadId = data[0].road_group_id

			return roadId
		},
		getInterventionParentId(state) {
			if (Object.keys(state?.data).length === 0) {
				return 0
			}
			const data = state.data.maintenance_roads
			const roadId = data[0].intervention_criteria_id

			return roadId
		},
		getPlanOptions(state) {
			if (state.planList.length === 0) {
				return []
			}

			const options: IOption[] = state.planList.map((item) => {
				return { label: item.name, value: item.name }
			})

			return options
		},
		getScheduleList(state) {
			if (state.statusList.length === 0) {
				return []
			}

			let scheduleList: IMaintenanceHistoryPlanStatusSchedule[] = []

			state.statusList.forEach((item) => {
				if (item.name === state.params.planName) {
					scheduleList = item.schedules
				}
			})

			return scheduleList
		},
		getUpdater(state) {
			if (Object.keys(state.data).length === 0) {
				return {}
			}
			const data = state.data?.updated_by

			return {
				name: data?.firstname && data?.lastname ? `${data?.firstname} ${data?.lastname}` : "",
				date: state.data?.updated_at
					? buddhistFormatDate(state.data?.updated_at, "เมื่อวันที่ dd mmm yyyy เวลา HH:ii น.")
					: "",
				department: data?.department?.name ? data.department?.name : "",
			}
		},
		getBudgetWithOutVat(state) {
			if (Object.keys(state.data).length === 0) {
				return {}
			}

			const data = state.data
			// const budget_procurement = data.budget_procurement
			const budgetWithOutVat = (data.budget_procurement * 100) / 107

			return budgetWithOutVat.toFixed(2)
		},
	},
})
