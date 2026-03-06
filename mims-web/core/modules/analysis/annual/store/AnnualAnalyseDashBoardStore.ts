import {
	IAnnualAnalyzeDataDefault,
	IAnnualDashboard,
	IAnnualHandleCheckBox,
	IAnnualMapData,
	IAnnualMapFilter,
	IAnnualMapItem,
} from "../infrastructure"
import { AnnualService } from "../infrastructure/AnnualService"
import { IReportParams } from "../../strategic/infrastructure/StrategicRequest"
import { DISPLAY_METHOD } from "../../strategic/store"

interface IStateExportReport {
	reportType: string
	reportTypeNo: string | null
	reportName: string
	reportPlan: string | null
	reportId: number
}

interface IStateMapParams {
	display: number
	criteria?: number
	year: number
}

interface IState {
	loading: boolean
	defaultData: IAnnualAnalyzeDataDefault
	exportReport: IStateExportReport
	dashboard: IAnnualDashboard
	pie_series_index: number[]
	tree_series_index: number[]
	map_filer_data: IAnnualMapFilter
	map_params: IStateMapParams
	map_data: IAnnualMapData
	map: any
	longdo: any
	isInit: boolean
	retryCreateLine: number
}

export const useAnnualSummaryDashboardStore = defineStore("annual/summary/budget-limit", {
	state: (): IState => ({
		loading: false,
		defaultData: {} as IAnnualAnalyzeDataDefault,
		exportReport: {
			reportType: "",
			reportTypeNo: null,
			reportName: "",
			reportPlan: "",
			reportId: 0,
		},
		dashboard: {} as IAnnualDashboard,
		pie_series_index: [],
		tree_series_index: [],
		map_filer_data: {} as IAnnualMapFilter,
		map_params: {
			display: 1,
			criteria: 1,
			year: 1,
		},
		map_data: {} as IAnnualMapData,
		map: null,
		longdo: null,
		isInit: true,
		retryCreateLine: 0,
	}),
	actions: {
		async getDefaultData(id: number) {
			this.loading = true

			const service = new AnnualService()
			const res = await service.getDefaultData(id)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.defaultData = res.data
			}
		},
		async createFavorite(id: number) {
			const service = new AnnualService()
			const res = await service.createFavorite(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				return res
			}
		},
		async copy(id: number) {
			this.loading = true

			const service = new AnnualService()
			const res = await service.createCopy(id)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				return res
			}
		},
		async getDashboardData(id: number) {
			this.loading = true

			const service = new AnnualService()
			const res = await service.getAnnualDashboard(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.dashboard = res.data
				this.pie_series_index = res.data?.bar1?.data?.map((_, index) => index)
				this.tree_series_index = res.data?.bar2?.data.map((_, index) => index)
				// this.map_params.year = res.data?.analyst_year
			}
			this.loading = false
		},
		async getMapFilter(id: number) {
			const service = new AnnualService()
			const res = await service.getMapFilterOptions(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.map_filer_data = res.data

				if (this.getCriteriaOptions.length) {
					this.map_params.criteria = this.getCriteriaOptions[0]?.value
				}
			}
		},
		async getMapData(id: number) {
			const service = new AnnualService()
			const res = await service.getMapData(id, this.map_params)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.map_data = res.data

				// if (this.isInit) {
				// 	setTimeout(() => {
				// 		this.createLine()
				// 	}, 1500)
				// 	this.isInit = false
				// }

				this.createLine()
			}
		},
		setMap(map: any) {
			this.map = map

			if (this.map) {
				// @ts-ignore
				this.longdo = window.longdo
			}
		},
		// createLine() {
		// 	if (this.map_data && this.map) {
		// 		this.map.Overlays.clear()

		// 		this.map_data?.items?.forEach((item) => {
		// 			const coordinates = item.the_geom?.coordinates?.map((coord) => ({ lon: coord[0], lat: coord[1] }))
		// 			const detail = this.createHtmlPopUp(item)
		// 			const line = new this.longdo.Polyline(coordinates, { lineColor: item.color, detail })

		// 			this.map.Overlays.add(line)

		// 			this.map.location({
		// 				lon: coordinates[0].lon,
		// 				lat: coordinates[0].lat,
		// 			})

		// 			this.map.zoom(18)
		// 		})
		// 	}
		// },
		createLine() {
			if (this.map) {
				this.map.Overlays.clear()
				console.log("createLine 1")
				if (this.map_data?.items?.length > 0) {
					console.log("createLine 2")
					const lines = this.map_data?.items?.map((item) => {
						return this.longdo.Util.overlayFromGeoJSON(item.the_geom, {
							lineColor:
								item.color === "#000000"
									? convertHexToRGBA(item.color ? item.color : "#000000", 0.5)
									: convertHexToRGBA(item.color ? item.color : "#000000", 1),
							detail: this.createHtmlPopUp(item),
						})
					})

					lines.forEach((line) => {
						this.map.Overlays.add(line[0])
					})
					console.log("createLine 3")

					const geom =
						this.map_data?.items[0].the_geom.coordinates.length > 0
							? this.map_data?.items[0].the_geom.coordinates[0]
							: undefined
					// const latLon = getLatLong(geom)
					// @ts-ignore
					// const longdo = window.longdo
					// const popup = new longdo.Popup(latLon, {
					// 	detail: this.generatePopupDetails(this.data, this.getRoadGeom[0].road),
					// })
					// this.map.Overlays.add(popup)
					if (geom) {
						this.map.location({
							lon: geom[0],
							lat: geom[1],
						})
					}

					this.map.zoom(15)
				}
			} else {
				console.log(this.retryCreateLine)
				if (this.retryCreateLine <= 5) {
					setTimeout(() => {
						this.createLine()
					}, 500)
				}

				this.retryCreateLine++
			}
		},
		createHtmlPopUp(item: IAnnualMapItem) {
			let html = ""

			html = `<div class="d-flex gap-2">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <g clip-path="url(#clip0_2983_17320)">
        <path d="M0 9H24V22H16.152C17.586 20.808 18.5 19.011 18.5 17C18.5 14.297 16.849 11.98 14.501 11H9.5C7.152 11.98 5.501 14.297 5.501 17C5.501 19.011 6.415 20.808 7.849 22H0V9ZM14 12.988V15.5C14 16.605 13.105 17.5 12 17.5C10.895 17.5 10 16.605 10 15.5V12.988C8.524 13.726 7.5 15.237 7.5 17C7.5 18.956 8.756 20.605 10.5 21.224V24H13.5V21.224C15.244 20.604 16.5 18.956 16.5 17C16.5 15.237 15.476 13.727 14 12.988ZM24 3V7H0V3C0 1.346 1.346 0 3 0H21C22.654 0 24 1.346 24 3ZM5 3.5C5 2.672 4.328 2 3.5 2C2.672 2 2 2.672 2 3.5C2 4.328 2.672 5 3.5 5C4.328 5 5 4.328 5 3.5ZM9 3.5C9 2.672 8.328 2 7.5 2C6.672 2 6 2.672 6 3.5C6 4.328 6.672 5 7.5 5C8.328 5 9 4.328 9 3.5Z" fill="${
					item.color
				}"/>
        </g>
        <defs>
        <clipPath id="clip0_2983_17320">
        <rect width="24" height="24" fill="white"/>
        </clipPath>
        </defs>
        </svg>
        <h4 class="fw-semibold ps-2" style="color:${item.color}">${item.title}</h4>
      </div>
      <div class="row mt-2">
        <div class="col-4 text-dark mb-2">สายทาง:</div>
        <div class="col-6 text-gray-700 mb-2">${item.road_name}</div>
        <div class="col-4 text-dark mb-2">ปีที่:</div>
        <div class="col-6 text-gray-700 mb-2">${item.year}</div>
        <div class="col-4 text-dark mb-2">IRI ก่อนซ่อม:</div>
        <div class="col-6 text-gray-700 mb-2">${toNumber(item.iri_before, 2)}  ม./กม.</div>
        <div class="col-4 text-dark mb-2">IRI หลังซ่อม:</div>
        <div class="col-6 text-gray-700 mb-2">${toNumber(item.iri_after, 2)} ม./กม.</div>
    </div>`

			return html
		},
		handleReport(id: number, reportNo: string, name: string, plan?: number) {
			this.exportReport.reportId = id
			this.exportReport.reportTypeNo = reportNo
			this.exportReport.reportName = name
			this.exportReport.reportPlan = plan !== undefined ? plan?.toString() : null
		},
		encodeQuery(data: IReportParams) {
			if (data) {
				const result: string[] = []
				const keys = Object.keys(data).filter((key: any) => data[key as keyof typeof data] !== undefined)
				for (const key of keys) {
					const value = data[key as keyof typeof data]
					if (value !== null && value !== undefined) {
						result.push(`${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
					}
				}
				return result.join("&")
			} else {
				return ""
			}
		},
		updateGraph(treeChart: Ref, pieChart: Ref) {
			treeChart.value?.updateOptions({
				dataLabels: {
					formatter: (value: any, opts: any) => {
						return `${value}: ${opts.value}%`
					},
				},
				colors: this.getTreeChartOptions?.colors,
				labels: this.getTreeChartOptions?.label,
				legend: {
					show: false,
					customLegendItems: this.getTreeLegends.label,
				},
			})

			pieChart.value?.updateOptions({
				labels: this.getPieChartOptions?.label,
				colors: this.getPieChartOptions?.colors,
			})
		},

		handleCheckbox(index: number, name: IAnnualHandleCheckBox) {
			const indexMap: Record<IAnnualHandleCheckBox, number[]> = {
				"pie-chart-index": this.pie_series_index,
				"tree-map-chart": this.tree_series_index,
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

			targetIndexArray.sort((a, b) => a - b)
		},
	},
	getters: {
		getRoadsList(state) {
			const { dashboard } = state
			const roadList = dashboard?.road

			return roadList || []
		},
		getDashboardFilter(state) {
			const { dashboard } = state
			const filterInfo = dashboard?.filter

			return `ชนิดผิวทาง : ${filterInfo?.surface_type ?? "-"}, ช่องจราจร : ${filterInfo?.lane ?? "-"}, จัดกลุ่ม : ${
				filterInfo?.km ?? "-"
			} กม.`
		},
		getDashboardConditions(state) {
			const condtions = state.dashboard?.condition

			return `เงื่อนไข : ${condtions?.condition ?? "-"}, เป้าหมาย: ${
				condtions?.target ?? "-"
			}, อัตราคิดลด (Discount Rate) : ${
				condtions?.discount || condtions?.discount === 0 ? condtions?.discount + "%" : "0%"
			} `
		},
		getComment(state) {
			const comment = state.dashboard?.comment
			const regex = /(\n|<br\s*\/?>)/
			const check = regex.test(comment)

			return check ? comment.split("\n") : [comment]
		},
		getPieChartOptions(state) {
			const data = state.dashboard?.bar1
			const indexs = state.pie_series_index
			const series = indexs.map((index) => data?.data[index])
			const label = indexs.map((index) => data.lable[index])
			const area = indexs.map((index) => data?.area[index])
			const colors = indexs.map((index) => data?.color[index])

			return { series, label, area, colors }
		},
		getTreeChartOptions(state) {
			const data = state.dashboard?.bar2
			const indexs = state.tree_series_index
			const series = indexs.map((index) => ({ x: data?.lable[index], y: data?.data[index] }))
			const label = indexs.map((index) => data?.lable[index])
			const colors = indexs.map((index) => data?.color[index])
			const budgets = indexs.map((index) => data?.budget[index])

			return { series: [{ data: series }], label, colors, budgets }
		},
		getPieCategories(state) {
			const data = state.dashboard.bar1
			const color = data?.color
			const label = data?.lable

			return { label, color }
		},
		getTreeLegends(state) {
			const data = state.dashboard.bar2
			const color = data?.color
			const label = data?.lable

			return { label, color }
		},
		getTable1(state) {
			const data = state.dashboard?.table?.table1 || {}

			return data
		},
		getTable2(state) {
			const data = state.dashboard?.table?.table2 || []

			return data
		},
		getBudgets(state) {
			const data = state.dashboard?.bar2

			if (Object.keys(data)?.length === 0) {
				return []
			}

			const budget = state.dashboard?.bar2?.budget

			return budget
		},
		getArea(state) {
			const data = state.dashboard?.bar1

			if (Object.keys(data)?.length === 0) {
				return []
			}

			const area = state.dashboard?.bar1.area

			return area
		},
		getFilterList(state) {
			const data = state.dashboard?.filter

			const result = data?.filter.map((condition) => {
				return condition
					.split(" ")
					.map((part, index) => {
						if (index !== 1 && index !== 3) {
							const num = parseFloat(part)
							if (!isNaN(num)) {
								return toNumber(num, 2)
							}
						}
						return part
					})
					.join(" ")
			})

			return result || ""
		},
		getCriteriaOptions(state) {
			const criteria = state.map_filer_data?.criteria

			const options = criteria.map((item) => ({ label: item.name, value: item.id }))
			return options || []
		},
		getCriteriaLegend(state) {
			const criteriaMethod = state.map_data.criteria_method
			// const criteria = mapData?.items?.map((item) => ({ name: item.title, color: item.color }))
			return criteriaMethod || []
		},
		getCriteriaGradeAC(state) {
			if (!state.map_params.criteria) {
				return state.map_filer_data.criteria.length > 0 ? state.map_filer_data.criteria[0].grade : []
			}
			return (
				state.map_filer_data.criteria?.find((e) => {
					return e.id === state.map_params.criteria
				})?.grade ?? []
			)
		},
		getCriteriaGradeCC(state) {
			if (!state.map_params.criteria) {
				return state.map_filer_data.criteria.length > 0 ? state.map_filer_data.criteria[0].grade_cc : []
			}
			return (
				state.map_filer_data.criteria?.find((e) => {
					return e.id === state.map_params.criteria
				})?.grade_cc ?? []
			)
		},
		getSymbolBoxHeight(state) {
			if (state.map_params.display === DISPLAY_METHOD) {
				return state.map_data?.criteria_method?.length > 5 ? "230" : state.map_data?.criteria_method?.length * 30 + 80
			} else {
				const data = state.map_filer_data.criteria?.find((e) => e.id === state.map_params.criteria)
				if (data) {
					return data.grade?.length * 30 + data.grade_cc?.length * 30 + 155
				} else {
					return "58"
				}
			}
		},
		getSymbolBoxWidth(state) {
			return state.map_data?.criteria_method?.length > 5 ? 280 : 230
		},
		getSymbolPositionLeft(state) {
			return state.map_data?.criteria_method?.length > 5 ? 204 : 154
		},
		getSymbolBoxPositionTop(state) {
			if (state.map_params.display === DISPLAY_METHOD) {
				return (state.map_data?.criteria_method?.length > 5 ? 5 : state.map_data?.criteria_method?.length) * 30 + 38
			} else {
				const data = state.map_filer_data.criteria?.find((e) => e.id === state.map_params.criteria)
				if (data) {
					return data.grade?.length * 30 + data.grade_cc?.length * 30 + 110
				} else {
					return 350
				}
			}
		},
		getGradeCCPositionTop(state) {
			const data = state.map_filer_data.criteria?.find((e) => e.id === state.map_params.criteria)
			if (data) {
				return data.grade?.length * 30 + 70
			} else {
				return 196
			}
		},
	},
})
