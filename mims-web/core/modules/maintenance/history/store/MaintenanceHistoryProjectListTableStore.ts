import { LocationQuery } from "vue-router"
import {
	IMaintenanceHistoryRoadGroupListData,
	IMaintenanceBudgetCriteria,
	IMaintenanceHistoryListItem,
	IMaintenanceHistoryDetailMaintenanceRoad,
	IMaintenanceHistoryDevisionList,
	IMaintenanceHistoryRoadDropdownList,
} from "../infrastructure/MaintenanceHistoryModel"
import { MaintenanceHistoryService } from "../infrastructure/MaintenanceHistoryService"

interface IStateParams {
	road_group_id: number[]
	road_group_id_dashboard: number[]
	budget_year: number | null
	budget_type_id: number | null
	budget_method_id: number | null
	budget_maintenance: number | null
	name: string
	owner_code: string | null
}

interface IState {
	loading: boolean
	datas: IMaintenanceHistoryListItem[]
	params: IStateParams
	budgetCriteriaList: IMaintenanceBudgetCriteria[]
	roadGroupList: IMaintenanceHistoryRoadGroupListData[]
	yearList: number[]
	fullScreen: boolean
	map: any
	devisionList: IMaintenanceHistoryDevisionList[]
	roadDropdownList: IMaintenanceHistoryRoadDropdownList[]
}

export const useMaintenanceHistorySearchTableStore = defineStore("maintenance/history/search-table", {
	state: (): IState => ({
		loading: false,
		datas: [],
		params: {
			road_group_id: [],
			road_group_id_dashboard: [],
			budget_year: null,
			budget_type_id: null,
			budget_method_id: null,
			budget_maintenance: null,
			name: "",
			owner_code: null,
		},
		budgetCriteriaList: [],
		roadGroupList: [],
		yearList: [],
		fullScreen: false,
		map: null,
		devisionList: [],
		roadDropdownList: [],
	}),
	actions: {
		// async fetchMaintenanceHistoryList() {
		// 	this.loading = true

		// 	this.setParams()

		// 	const service = new MaintenanceHistoryService()
		// 	const res = await service.getMaintenanceHistoryList(this.params)

		// 	if (!res.status) {
		// 		useHandlerError(res.code, res.error, { showAlert: true })
		// 	} else {
		// 		this.datas = res.data.items
		// 	}
		// 	this.loading = false
		// },
		setDatas(datas: IMaintenanceHistoryListItem[]) {
			this.datas = datas
			// Defer map drawing so table paints first (avoids blocking initial load)
			if (typeof requestIdleCallback !== "undefined") {
				requestIdleCallback(() => this.createLine(), { timeout: 100 })
			} else {
				setTimeout(() => this.createLine(), 0)
			}
		},
		createLine() {
			if (this.map) {
				this.map.Overlays.clear()

				// @ts-ignore
				const longdo = window.longdo

				const items = this.datas.length ? this.datas : []

				items.forEach((item) => {
					item.roads.forEach((subItem) => {
						const strLine = longdo.Util.overlayFromGeoJSON(subItem.the_geom, {
							lineColor: item.color === "#000000" ? convertHexToRGBA(item.color, 0.5) : convertHexToRGBA(item.color, 1),
							detail: this.generatePopupDetails(item),
						})

						this.map.Overlays.add(strLine[0])
					})
				})

				if (items.length) {
					const roads = items.filter((item) => item.roads.length)
					const geom = roads[0].roads[0]?.the_geom?.coordinates?.length
						? roads[0]?.roads[0]?.the_geom?.coordinates[0]
						: []

					if (geom) {
						this.map.location({
							lon: geom[0],
							lat: geom[1],
						})
					}
				}

				this.map.zoom(12)
			}
		},
		resetSearch() {
			const route = useRoute()
			if (Object.keys(route.query).length) {
				navigateTo("/maintenances/history")
			}

			this.params.budget_maintenance = null
			this.params.budget_method_id = null
			this.params.budget_type_id = null
			this.params.budget_year = null
			this.params.name = ""
			this.params.owner_code = null
			this.params.road_group_id = []
		},
		async getYearList() {
			this.loading = true

			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceYearList()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.yearList = res.data
			}

			this.loading = false
		},
		setMap(map: Object) {
			this.map = map
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
		async getBudgetCriteria() {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceBudgetCriteria()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.budgetCriteriaList = res.data
				// โหลด road group พร้อมกัน ไม่ต้องรอ (ไม่ depend กัน)
				this.getRoadGroupList()
			}
		},
		async getDevision() {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceDivision()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.devisionList = res.data
			}
		},
		async getRoadDropdownList() {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceRoadGroup()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.roadDropdownList = res.data
			}
		},
		async getRoadGroupList() {
			const service = new MaintenanceHistoryService()
			const res = await service.getHistoryRoadGroup()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.roadGroupList = res.data
			}
		},
		setDateColor(guaranteeDate: Date) {
			const currentDate = new Date()
			const dateToCompare = new Date(guaranteeDate)

			// Calculate the time difference in milliseconds
			const diffTime = dateToCompare.getTime() - currentDate.getTime()

			// Calculate the time difference in days
			const diffDays = diffTime / (1000 * 3600 * 24)

			// Calculate the time difference in months (approximate)
			const diffMonths = diffDays / 30

			if (diffDays < 0) {
				return ""
			} else if (diffMonths <= 3) {
				return "badge badge-light-danger text-danger"
			} else if (diffMonths <= 6) {
				return "badge badge-light-warning  text-warning"
			} else {
				return "badge badge-light-success  text-success"
			}
		},
		// setParams() {
		// 	for (const key in this.params) {
		// 		const value = this.params[key]
		// 		if (
		// 			value === undefined ||
		// 			(Array.isArray(value) && value.length === 0) ||
		// 			(typeof value === "string" && value.trim() === "")
		// 		) {
		// 			delete this.params[key]
		// 		}
		// 	}
		// },
		setQuriesParams(quries: LocationQuery) {
			quries.road_group_id_dashboard = quries.road_group_id_dashboard as string
			const roadGroupId = quries.road_group_id_dashboard?.split(",")

			roadGroupId.forEach((id) => {
				this.params.road_group_id_dashboard?.push(Number(id))
			})

			// ดักเคส duplicate id
			this.params.road_group_id_dashboard = [...new Set(this.params.road_group_id_dashboard.filter((id) => id))]
		},
		getRoadName(items: IMaintenanceHistoryDetailMaintenanceRoad[]) {
			return items.length > 0 ? `${items[0].road_name}` : "-"
		},
		getKmStartEnd(items: IMaintenanceHistoryDetailMaintenanceRoad[]) {
			return items.length > 0 ? `${convertMeterToKm(items[0].km_start)} - ${convertMeterToKm(items[0].km_end)}` : "-"
		},
		getLane(items: IMaintenanceHistoryDetailMaintenanceRoad[]) {
			return items.length > 0 ? `${items[0].lane_no}` : "-"
		},
		getRoadGroupNames(item: Array<string | string>) {
			return item.length > 0 ? item.join("<br>") : "-"
		},
		// road_group_names
		getProcessHeaderTitle(item: IMaintenanceHistoryListItem) {
			// let icon = ""
			// if (item.percent_progress_plan > item.percent_progress_result) {

			// 	icon = "/images/icons/svg/fi-br-plan-strategy-red.svg"
			// } else if (item.percent_progress_plan === item.percent_progress_result) {

			// 	icon = "/images/icons/svg/fi-br-plan-strategy-yellow.svg"
			// } else if (item.percent_progress_plan < item.percent_progress_result) {

			// 	icon = "/images/icons/svg/fi-br-plan-strategy-green.svg"
			// }

			return `<div class="col-1">

				<i class="fi fi-sr-clock" style="font-size: 23px; color: ${item.color};"></i>
			</div>
			<div class="col"
			style="
			color: ${item.color};
			font-size: 16px;
			font-weight: 500;
			">${item.remaining_time}</div>`
		},
		generatePopupDetails(item: IMaintenanceHistoryListItem) {
			return `
			<div class="row mb-3" >
				${this.getProcessHeaderTitle(item)}
		 	</div>
				<div class="row">
				<div class="col-6">ชื่อโครงการ:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${item.name}</div>
			</div>
				<div class="row">
				<div class="col-6">เลขที่สัญญาโครงการ:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${item.contract_number}</div>
			</div>
			<div class="row">
				<div class="col-6">สายทาง:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${item.road_group_names}</div>
			</div>
			<div class="row mb-2">
				<div class="col-6">จาก - ถึง:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${this.getRoadName(item.roads)}</div>
			</div>
			<div class="row mb-2">
				<div class="col-6">กม. เริ่มต้น - กม. สิ้นสุด:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${this.getKmStartEnd(item.roads)}</div>
			</div>
			<div class="row mb-2">
				<div class="col-6">ช่องจราจรที่:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${this.getLane(item.roads)}</div>
			</div>
			<div class="row mb-2">
				<div class="col-6">ระยะทาง:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${item.km_total} กม.</div>
			</div>
			<div class="row mb-2">
				<div class="col-6">ตรวจรับงานงวดสุดท้าย:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${buddhistFormatDate(item.project_end_date, "dd mmm yy")}</div>
			</div>
			<div class="row mb-2">
				<div class="col-6">หมดค้ำประกัน:</div>
				<div class="col" style="
				color: #5E6278;
				font-size: 12px;
				font-weight: 400;
			">${buddhistFormatDate(item.project_guarantee_expiration_date, "dd mmm yy") || "-"} </div>
			</div>
			<div class="row mt-2 pt-2 maintenance-popup-footer" style="border-top: 1px solid #E4E6EF;">
				<div class="col-12">
					<a href="/maintenances/history/${item.id_parent}/info"
						class="maintenance-popup-view-details d-inline-flex align-items-center gap-1 text-primary text-decoration-none fw-semibold"
						data-maintenance-info-id="${item.id_parent}"
						style="font-size: 13px; cursor: pointer; transition: opacity 0.2s;">
						ดูรายละเอียดโครงการ
						<i class="fi fi-sr-angle-right" style="font-size: 12px;"></i>
					</a>
				</div>
			</div>
      `
		},
		extractText(name: string, fullname = false) {
			if (fullname === true) {
				const regex = /\s*\([^)]+\)/
				return name.replace(regex, "")
			} else {
				const match = name.match(/\(([^)]+)\)/)
				return match ? match[0] : null
			}
		},
		getFileExtension(fileName: String) {
			return fileName.slice(((fileName.lastIndexOf(".") - 1) >>> 0) + 2)
		},
		getDivisionOption() {
			const refDivision = useInitData().refDivision()

			if (refDivision?.length === 0) {
				return []
			}

			const options = useInitData()
				.refDivision()
				?.map((item) => {
					if (item.districts.length !== 0) {
						const districtsArray = item.districts.map((subDistricts) => {
							const depotsArray = subDistricts.depots.map((subDepots) => ({
								label: subDepots.name,
								id: [subDepots.owner_code_key],
							}))
							return {
								label: subDistricts.name,
								id: [subDistricts.owner_code_key],
								children: depotsArray,
							}
						})
						return {
							label: item.name,
							id: [item.owner_code_key],
							children: districtsArray,
						}
					}
					return {
						label: item.name,
						id: [item.owner_code_key],
						children: [],
					}
				})

			return options
		},
	},
	getters: {
		getBudgetCriteriaOptions(state) {
			// if (state.budget.length === 0) {
			if (state.budgetCriteriaList.length === 0) {
				return []
			}

			// const options = state.budgetCriteriaList.map(item => ({
			// 	label: item.name,
			// 	id: item.id,
			// 	...(item.budget_methods.length > 0 && {
			// 		children: item.budget_methods
			// 			.filter(element => !element.is_deleted)
			// 			.map(({ method_name, id }) => ({ label: method_name, id: id }))
			// 	})
			// })).sort((a, b) => a.id - b.id);
			const options = state.budgetCriteriaList
				.map((item) => ({
					label: item.name,
					id: item.id,
					...(item.budget_methods.length > 0 && {
						children: item.budget_methods
							.filter((element) => !element.is_deleted)
							.map(({ method_name: methodName, id }) => ({
								// Updated to camel case
								label: methodName,
								id,
							})),
					}),
				}))
				.sort((a, b) => a.id - b.id)

			return options
		},
		// getDivisionOption(state) {
		// 	// if (state.devisionList.length === 0) {
		// 	// 	return []
		// 	// }

		// 	// const options = state.devisionList.map((item) => {
		// 	// 	if (item.districts.length !== 0) {
		// 	// 		const districtsArray = item.districts.map((subDistricts) => {
		// 	// 			const depotsArray = subDistricts.depots.map((subDepots) => ({
		// 	// 				label: subDepots.name,
		// 	// 				id: [subDepots.owner_code_key],
		// 	// 			}))
		// 	// 			return {
		// 	// 				label: subDistricts.name,
		// 	// 				id: [subDistricts.owner_code_key],
		// 	// 				children: depotsArray,
		// 	// 			}
		// 	// 		})
		// 	// 		return {
		// 	// 			label: item.name,
		// 	// 			id: [item.owner_code_key],
		// 	// 			children: districtsArray,
		// 	// 		}
		// 	// 	}
		// 	// 	return {
		// 	// 		label: item.name,
		// 	// 		id: [item.owner_code_key],
		// 	// 		children: [],
		// 	// 	}
		// 	// })

		// 	const refDivision = useInitData().refDivision()

		// 	console.log("getDivisionOption ---->", refDivision);

		// 	if (refDivision?.length === 0) {
		// 		return []
		// 	}

		// 	const options = useInitData().refDivision()?.map((item) => {
		// 		if (item.districts.length !== 0) {
		// 			const districtsArray = item.districts.map((subDistricts) => {
		// 				const depotsArray = subDistricts.depots.map((subDepots) => ({
		// 					label: subDepots.name,
		// 					id: [subDepots.owner_code_key],
		// 				}))
		// 				return {
		// 					label: subDistricts.name,
		// 					id: [subDistricts.owner_code_key],
		// 					children: depotsArray,
		// 				}
		// 			})
		// 			return {
		// 				label: item.name,
		// 				id: [item.owner_code_key],
		// 				children: districtsArray,
		// 			}
		// 		}
		// 		return {
		// 			label: item.name,
		// 			id: [item.owner_code_key],
		// 			children: [],
		// 		}
		// 	})

		// 	return options
		// },
		getRoadDropdownOption(state) {
			if (state.roadDropdownList.length === 0) {
				return []
			}

			const options = state.roadDropdownList.map((item) => {
				return { label: item.road_number, id: item.id }
			})

			return options
		},
		getRoadGroupOptions(state) {
			if (state.roadGroupList.length === 0) {
				return []
			}

			const options = state.roadGroupList.map((item) => {
				return { label: item.name, value: item.id }
			})

			return options
		},
		getYearsOptions(state) {
			const options = state.yearList.map((item) => {
				return { label: (item + 543).toString(), value: item }
			})

			return options || []
		},
	},
})
