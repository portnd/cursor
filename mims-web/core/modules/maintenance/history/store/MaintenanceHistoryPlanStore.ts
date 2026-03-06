import {
	IMaintenancePlanListData,
	IMaintenanceHistoryPlanGraph,
	IMaintenanceHistoryDetailData,
	IPlanProgressGraphReportHistTableData,
} from "../infrastructure/MaintenanceHistoryModel"
import { MaintenanceHistoryService } from "../infrastructure/MaintenanceHistoryService"
import { IOption } from "~/core/shared/types/Option"

interface IParamState {
	planId: number[]
}

interface ISumDataTable {
	[key: string]: number
}

interface IState {
	loading: boolean
	projectInfo: IMaintenanceHistoryDetailData
	graphData: IMaintenanceHistoryPlanGraph[]
	tableData: IPlanProgressGraphReportHistTableData
	planList: IMaintenancePlanListData[]
	params: IParamState
	sumDataTable: ISumDataTable[]
	problemList: string[]
	map: any
	longdo: any
}

export const useMaintenanceHistoryPlanStore = defineStore("maintenance/history/plan", {
	state: (): IState => ({
		loading: false,
		projectInfo: {} as IMaintenanceHistoryDetailData,
		graphData: [],
		tableData: {} as IPlanProgressGraphReportHistTableData,
		planList: [],
		params: {
			planId: [],
		},
		sumDataTable: [],
		problemList: [],
		map: null,
		longdo: null,
	}),
	actions: {
		async getPlanList(id: number) {
			this.loading = true

			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceStatusOptions(id)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.planList = res.data
				this.setDefaultStatusCurrent(this.planList)
			}
		},
		async getProjectInfo(id: number) {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceHistoryDetails(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.projectInfo = res.data
			}
		},
		async getPlanGraph(id: number, planId: number[]) {
			this.loading = true

			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceHistoryPlanGraphRerport(id, planId)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.graphData = res.data
				await this.getProjectInfo(id)
				await this.getPlanDataTable(id, planId)
			}

			this.loading = false
		},
		async getPlanDataTable(id: number, planId: number[]) {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceHistoryPlanTableReport(id, planId)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.tableData = res.data
			}
		},
		setMap(map: any) {
			this.map = map
			if (this.map) {
				// @ts-ignore
				this.longdo = window.longdo
			}

			this.createLine()
		},
		getSumDataTable() {
			const data = this.getDataTable

			if (data?.length === 0) {
				return
			}

			const result = data?.map((item) => {
				const values = item.value
				return {
					plan: values.reduce((acc, curr) => acc + curr.plan, 0),
					plan_total: values.reduce((acc, curr) => acc + curr.plan_total, 0),
					progress_plan: values.reduce((acc, curr) => acc + (curr.progress_plan || 0), 0),
					progress_plan_total: values.reduce((acc, curr) => acc + curr.progress_plan_total, 0),
					disburement_plan: values.reduce((acc, curr) => acc + curr.disbursement_plan, 0),
					disburement_plan_total: values.reduce((acc, curr) => acc + curr.disbursement_plan_total, 0),
					disburement_progress: values.reduce((acc, curr) => acc + (curr.disbursement_progress || 0), 0),
					disburement_progress_total: values.reduce((acc, curr) => acc + curr.disbursement_progress_total, 0),
				}
			})

			this.sumDataTable = result
		},
		setDefaultStatusCurrent(statusList: IMaintenancePlanListData[]) {
			if (statusList.length > 0) {
				statusList.forEach((item) => {
					if (item.is_current) {
						this.params.planId.push(item.id)
					}
				})
			}
		},
		async setGraph(id: number, planId: number[], refPlanChart: Ref, refDisburement: Ref) {
			await this.getPlanGraph(id, planId)
			this.getSumDataTable()

			const planColor = this.graphData
				.flatMap((parent) => {
					if (parent.name === "แผนการดำเนินงาน") {
						return parent.data.map((child) => child.color)
					}

					return undefined
				})
				.filter((item) => item !== undefined)

			const planCategories = this.graphData
				.flatMap((parent) => {
					if (parent.name === "แผนการดำเนินงาน") {
						return parent.schedule ? parent.schedule.map((label) => buddhistFormatDate(label, "mmm yy")) : []
					}

					return undefined
				})
				.filter((item) => item !== undefined)

			const disburementColors = this.graphData
				.flatMap((parent) => {
					if (parent.name === "การเบิกจ่ายเงิน") {
						return parent.data.map((child) => {
							return child.color
						})
					}
					return undefined
				})
				.filter((item) => item !== undefined)

			const disburementCategories = this.graphData
				.flatMap((parent) => {
					if (parent.name === "การเบิกจ่ายเงิน") {
						return parent.schedule ? parent.schedule.map((label) => buddhistFormatDate(label, "mmm yy")) : []
					}

					return undefined
				})
				.filter((item) => item !== undefined)

			refPlanChart.value.updateOptions({
				colors: planColor,
				xaxis: {
					categories: planCategories,
					tickAmount: 10,
				},
			})

			refDisburement.value.updateOptions({
				colors: disburementColors,
				xaxis: {
					categories: disburementCategories,
					tickAmount: 10,
				},
			})
		},
		createLine() {
			if (this.map) {
				const geoms = this.projectInfo.maintenance_roads.map((item) => {
					return { geoms: item.the_geom, colors: item.road_info.road_color_code }
				})

				const polylines = geoms.flatMap((item) => {
					return this.longdo.Util.overlayFromWkt(item.geoms, { lineColor: item.colors })
				})

				polylines.forEach((line) => {
					this.map.Overlays.add(line)
				})

				const latlon = getLatLong(geoms[0].geoms)

				this.map.location({
					lon: latlon.lon,
					lat: latlon.lat,
				})

				this.map.zoom(15)
			}
		},
		checkIsInteger(value: number | null) {
			if (value) {
				const isInteger = Number.isInteger(value)
				return toNumber(value, isInteger ? 0 : 2)
			} else {
				return "0"
			}
		},
	},
	getters: {
		getPlanListOptions(state) {
			if (state.planList?.length === 0) {
				return []
			}

			const options: IOption[] = state.planList.map((item) => {
				return { label: item.name, value: item.id }
			})

			return options
		},
		getDataTable(state) {
			if (state.tableData?.maintenance_plan?.length === 0) {
				return []
			}

			const data = state.tableData?.maintenance_plan
			const result = data?.map((item) => {
				return item
			})

			return result
		},

		getProblemDetails(state) {
			if (state.tableData?.problems?.length === 0) {
				return []
			}

			const data = state.tableData.problems
			// const result: string[] = []
			// data.forEach((_, index) => {
			// 	if (data[index]?.plan_name !== data[index + 1]?.plan_name) {
			// 		data.forEach((item) => {
			// 			if (item.plan_name === data[index].plan_name) {
			// 				item.problems.forEach((prob) => {
			// 					result.push(prob)
			// 				})
			// 			}
			// 		})
			// 	}
			// })

			return data
		},
		getSumTable(state) {
			const data = state.tableData.maintenance_plan
			const values = data.map((item) => item.value)

			const plan = values.map((item) => {
				return item
					.map((value) => {
						return value.plan
					})
					.reduce((acc, curr) => acc + curr)
			})

			const progressPlan = values.map((item) => {
				return item
					.map((value) => {
						return value.progress_plan
					})
					.reduce((acc, curr) => acc! + curr!)
			})

			return { plan, progress_plan: progressPlan }
		},
	},
})
