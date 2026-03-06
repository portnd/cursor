import {
	IReportParams,
	IStrategicAnalyzeData,
	IStrategicDashboard,
	IStrtegicMapFilter,
	IStrategicMapData,
} from "../infrastructure"
import { IMapDataReq } from "../infrastructure/StrategicRequest"
import { StrategicsService } from "../infrastructure/StrategicsService"

// Display type constants (matches server enum values)
export const DISPLAY_IRI = 1
export const DISPLAY_METHOD = 2

interface IStateExportReport {
	reportType: string
	reportTypeNo: string | null
	reportName: string
	reportPlan: string | null
	reportId: number
}

interface IMapParams {
	plan: number
	year: number
	display: number
	cirteria: number
}

interface IState {
	loading: boolean
	data: IStrategicDashboard
	defaultData: IStrategicAnalyzeData
	plan: string
	planTable: string
	exportReport: IStateExportReport
	mapFilter: IStrtegicMapFilter
	mapData: IStrategicMapData
	map: any
	longdo: any
	map_params: IMapParams
}

export const useStrategicAnalysisDashboardStore = defineStore("analyze/summary/analyze-dashboard", {
	state: (): IState => ({
		loading: false,
		data: {} as IStrategicDashboard,
		defaultData: {} as IStrategicAnalyzeData,
		plan: "",
		planTable: "",
		exportReport: {
			reportType: "",
			reportTypeNo: null,
			reportName: "",
			reportPlan: "",
			reportId: 0,
		},
		mapFilter: {} as IStrtegicMapFilter,
		mapData: {} as IStrategicMapData,
		map: null,
		longdo: null,
		map_params: {
			plan: 0,
			year: 0,
			display: DISPLAY_IRI,
			cirteria: 0,
		},
	}),
	actions: {
		// Fetches analyze details + dashboard data
		async getdata(id: number) {
			this.loading = true
			const service = new StrategicsService()

			const [resDefault, resDashboard] = await Promise.all([
				service.getAnalyzeDefaultDetails(id),
				service.getDashboard(id),
			])

			if (!resDefault.status) {
				useHandlerError(resDefault.code, resDefault.error, { showToast: true })
			} else {
				this.defaultData = resDefault.data
			}

			if (!resDashboard.status) {
				useHandlerError(resDashboard.code, resDashboard.error, { showToast: true })
			} else {
				this.data = resDashboard.data
				// Set initial plan tab for bar2 graph
				const planWithKeyword = this.data.bar2?.datasets?.find(
					(item) => item.plan && item.plan.includes("แผน")
				)
				const firstAvailablePlan = this.data.bar2?.datasets?.find((item) => item.plan)
				this.plan = planWithKeyword?.plan || firstAvailablePlan?.plan || ""
			}

			this.loading = false
		},

		// Fetches map filter options and sets initial map params
		async getMapfilter(id: number) {
			const service = new StrategicsService()
			const res = await service.getMapFilter(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
				return
			}

			this.mapFilter = res.data

			// Set initial map params from filter options
			if (res.data.plan?.length) {
				this.map_params.plan = res.data.plan[0].id
			}
			if (res.data.year?.length) {
				this.map_params.year = res.data.year[0]
			}
			if (res.data.display?.length) {
				this.map_params.display = res.data.display[0].id
			}
			if (res.data.criteria?.length) {
				this.map_params.cirteria = res.data.criteria[0].id
			}
		},

		// Fetches map polyline data with current map_params
		async getMapData(id: number) {
			const service = new StrategicsService()
			const params: IMapDataReq = {
				plan: this.map_params.plan,
				year: this.map_params.year,
				display: this.map_params.display,
			}
			if (this.map_params.display === DISPLAY_IRI && this.map_params.cirteria) {
				params.criteria = this.map_params.cirteria
			}

			const res = await service.getMap(id, params)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
				return
			}

			this.mapData = res.data
			this.drawMapData()
		},

		// Re-fetches map data when filter selection changes
		async onMapSelected(id: number) {
			await this.getMapData(id)
		},

		// Called when Longdo map instance is ready
		setMap(map: any) {
			this.map = map
			if (this.map) {
				// @ts-ignore
				this.longdo = window.longdo
			}
			this.drawMapData()
		},

		// Draws polylines from mapData onto the Longdo map
		drawMapData() {
			if (!this.map || !this.mapData?.items?.length) return

			this.map.Overlays.clear()

			this.mapData.items.forEach((item) => {
				if (!item.the_geom?.coordinates?.length) return

				const wkt = `LINESTRING(${item.the_geom.coordinates.map((c) => `${c[0]} ${c[1]}`).join(",")})`
				const lines = this.longdo?.Util?.overlayFromWkt(wkt, {
					detail: `<b>${item.title || ""}</b>`,
					lineColor: item.color,
				})
				if (lines?.length) {
					this.map.Overlays.add(lines[0])
				}
			})
		},

		togglePlan(planName: string) {
			this.plan = planName
		},

		togglePlanTable(planName: string) {
			this.planTable = planName
		},

		handleReportType(type: string) {
			this.exportReport.reportType = type
		},

		handleReport(id: number, reportNo: string, name: string, plan?: number) {
			this.exportReport.reportId = id
			this.exportReport.reportTypeNo = reportNo
			this.exportReport.reportName = name
			this.exportReport.reportPlan = plan !== undefined ? plan?.toString() : null
		},

		encodeQuery(data: IReportParams) {
			if (!data) return ""
			const keys = Object.keys(data).filter((key) => data[key as keyof typeof data] !== undefined)
			const result: string[] = []
			for (const key of keys) {
				const value = data[key as keyof typeof data]
				if (value !== null && value !== undefined) {
					result.push(`${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
				}
			}
			return result.join("&")
		},
	},
	getters: {
		// --- Road & Filter Info ---
		getRoadsList(state) {
			return (state.data?.road as string[]) || []
		},

		getFilter(state) {
			const filter = state.data?.filter
			if (!filter) return ""
			return `ชนิดผิวทาง : ${filter.surface_type}, ช่องจราจร : ${filter.lane}, จัดกลุ่ม : ${filter.km} กม.`
		},

		getFilterList(state) {
			const filterData = state.data?.filter
			const result = filterData?.filter?.map((condition) => {
				return condition
					.split(" ")
					.map((part, index) => {
						if (index !== 1 && index !== 3) {
							const num = parseFloat(part)
							if (!isNaN(num)) {
								return Number.isInteger(num) ? toNumber(num) : toNumber(num, 2)
							}
						}
						return part
					})
					.join(" ")
			})
			return result || ""
		},

		getDashboardConditions(state) {
			const condition = state.data?.condition
			if (!condition) return ""
			let result = `เงื่อนไข : ${condition.condition}, เป้าหมาย: ${condition.target}`
			if (condition.discount !== null && condition.discount !== undefined) {
				result += `, อัตราคิดลด (Discount Rate) : ${condition.discount}%`
			}
			return result
		},

		getComment(state) {
			const comment = state.data?.comment
			if (!comment) return ["-"]
			const regex = /(\n|<br\s*\/?>)/
			const check = regex.test(comment)
			return check ? comment.split("\n") : [comment]
		},

		// --- Chart Data ---
		getGraphCategories1(state) {
			return state.data?.graph1?.lable?.map((item) => (item + 543).toString()) || []
		},

		getGraph1LegendLabel(state) {
			return state.data?.graph1?.line || []
		},

		getLineChartSeries(state) {
			const graph1 = state.data?.graph1
			if (!graph1?.line?.length) return []
			return graph1.line.map((name, index) => ({
				name,
				data: graph1.value[index] || [],
			}))
		},

		getBar1Series(state) {
			const bar1 = state.data?.bar1
			if (!bar1?.datasets?.length) return []
			return bar1.datasets.map((d) => ({ name: d.lable, data: d.value }))
		},

		getBar1LegendLabel(state) {
			return state.data?.bar1?.datasets?.map((item) => item.lable) || []
		},

		getBar1ColorList(state) {
			return state.data?.bar1?.color || []
		},

		// Returns a function to get bar1 color by index (getter that returns a function)
		getBar1Color: (state) => (key: number) => {
			return state.data?.bar1?.color?.[key] || null
		},

		getBarChartCategories1(state) {
			return state.data?.bar1?.lable?.map((item) => (item + 543).toString()) || []
		},

		getBarChartCategories2(state) {
			return state.data?.bar2?.lable?.map((item) => (item + 543).toString()) || []
		},

		// Bar2 percentage series filtered by selected plan
		getBar2Series(state) {
			if (!Object.keys(state.data).length) return []
			const datasets = state.data?.bar2?.datasets?.filter((item) => item.plan === state.plan)
			const accumulatedData = datasets?.reduce((acc, parent) => {
				;(parent?.data || []).forEach((child) => {
					;(child?.lable || []).forEach((label, index) => {
						const existing = acc.get(label) || []
						existing.push(Number((child?.value[index] || 0).toFixed(2)))
						acc.set(label, existing)
					})
				})
				return acc
			}, new Map<string, number[]>())
			return Array.from(accumulatedData || new Map(), ([name, data]) => ({ name, data }))
		},

		// Bar2 budget data for tooltip
		getBudget(state) {
			if (!Object.keys(state.data).length) return []
			const datasets = state.data?.bar2?.datasets?.filter((item) => item.plan === state.plan)
			const accumulatedData = datasets?.reduce((acc, parent) => {
				;(parent?.data || []).forEach((child) => {
					;(child?.lable || []).forEach((label, index) => {
						const existing = acc.get(label) || []
						existing.push(Number((child?.budget[index] || 0).toFixed(2)))
						acc.set(label, existing)
					})
				})
				return acc
			}, new Map<string, number[]>())
			return Array.from(accumulatedData || new Map(), ([name, data]) => ({ name, data }))
		},

		getLegendBar2(state) {
			const filterData = state.data?.bar2?.datasets?.filter((item) => item?.plan === state.plan)
			return filterData?.flatMap((item) => item?.data)?.flatMap((item) => item?.lable) || []
		},

		// --- Table Data ---
		getSummaryTable(state) {
			const data = state.data?.table?.summary
			return data?.map((item) => {
				const itemData = item?.data
				return {
					name: item.name,
					plans: itemData[0],
					plans_2: itemData?.filter((_, index) => index > 0),
				}
			})
		},

		getSummaryPlanTable(state) {
			let planData = null
			let years: number[] = []

			switch (state.planTable) {
				case "แผนที่ 1":
					planData = state.data?.table?.plan_1
					break
				case "แผนที่ 2":
					planData = state.data?.table?.plan_2
					break
				case "แผนที่ 3":
					planData = state.data?.table?.plan_3
					break
				case "ไม่จำกัดงบประมาณ":
					planData = state.data?.table?.unlimited_plan
					break
			}

			if (planData) {
				years = Array.from(
					new Set(planData.flatMap((item) => item.data.map((child) => child.year + 543)))
				).sort((a, b) => a - b)
			}

			return { plan: planData, years }
		},

		getCalculateColspan(state) {
			return state.data?.table?.summary[0]?.data?.length
		},

		// --- Map Options ---
		getMapPlanOptions(state) {
			return state.mapFilter?.plan?.map((p) => ({ value: p.id, label: p.name })) || []
		},

		getMapYearOptions(state) {
			return state.mapFilter?.year?.map((y) => ({ value: y, label: String(y + 543) })) || []
		},

		getMapDisplayOptions(state) {
			return state.mapFilter?.display?.map((d) => ({ value: d.id, label: d.name })) || []
		},

		getCriteriaOptions(state) {
			return state.mapFilter?.criteria?.map((c) => ({ value: c.id, label: c.name })) || []
		},

		// --- Map Legend ---
		getGradeAC(state) {
			if (!state?.mapFilter) return []
			const selected = state.mapFilter.criteria?.find((c) => c.id === state.map_params?.cirteria)
			return selected?.grade?.map((g) => ({ name: g.name, color: g.color })) || []
		},

		getGradeCC(state) {
			if (!state?.mapFilter) return []
			const selected = state.mapFilter.criteria?.find((c) => c.id === state.map_params?.cirteria)
			return selected?.grade_cc?.map((g) => ({ name: g.name, color: g.color })) || []
		},

		getCriteriaLegend(state) {
			if (!state?.mapData) return []
			return (state.mapData.criteria_method as { name: string; color: string }[]) || []
		},

		// Position helpers for map legend box (unit: px)
		getGradeCCPositionTop(): number {
			if (!this?.getGradeAC) return 0
			const acCount = this.getGradeAC.length
			// 30px content start offset + each AC row ~35px + 30px separator gap
			return 30 + acCount * 35 + 30
		},

		getSymbolBoxHeight(): number {
			if (!this?.map_params) return 0
			if (this.map_params.display === DISPLAY_IRI) {
				const acCount = this.getGradeAC?.length ?? 0
				const ccCount = this.getGradeCC?.length ?? 0
				// 30px header + AC rows + 80px CC section (header+separator) + CC rows + 20px bottom
				return 30 + acCount * 35 + 80 + ccCount * 35 + 20
			}
			const itemCount = this.getCriteriaLegend?.length ?? 0
			const rows = itemCount > 5 ? Math.ceil(itemCount / 2) : itemCount
			return 30 + rows * 35 + 20
		},

		getSymbolBoxPositionTop(): number {
			if (!this?.map_params) return 0
			return this.getSymbolBoxHeight + 55
		},

		getSymbolBoxWidth(): number {
			return (this?.getCriteriaLegend?.length ?? 0) > 5 ? 280 : 230
		},

		getSymbolPositionLeft(): number {
			return (this?.getCriteriaLegend?.length ?? 0) > 5 ? 205 : 155
		},
	},
})
