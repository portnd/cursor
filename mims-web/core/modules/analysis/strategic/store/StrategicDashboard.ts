import {
	IMapDataReq,
	IReportParams,
	IStrategicDashboard,
	IStrategicDashboardTableUnlimitedPlan,
	IStrategicMapData,
	IStrategicMapDataCriteriaMethod,
	IStrategicMapDataItem,
	IStrtegicMapFilter,
	StrategicsService,
} from "../infrastructure"

interface IStateExportReport {
	reportType: string
	reportTypeNo: string | null
	reportName: string
	reportPlan: string | null
	reportId: number
}

interface IStateMapParams {
	plan: number
	year: number
	display: number
	cirteria?: number
	method?: number
}

interface IState {
	loading: boolean
	data: IStrategicDashboard
	map_filter: IStrtegicMapFilter
	map_params: IStateMapParams
	map_data: IStrategicMapData
	plan: string
	planTable: string
	map: any
	longdo: any
	map_legend: IStrategicMapDataCriteriaMethod[]
	isInit: boolean
	exportReport: IStateExportReport
	retryCreateLine: number
}

export const DISPLAY_IRI = 1
export const DISPLAY_METHOD = 2

export const useStrategicAnalysisDashboardStore = defineStore("strategic/new-dashboard", {
	state: (): IState => ({
		loading: false,
		data: {} as IStrategicDashboard,
		map_filter: { year: [], criteria: [], method: [], plan: [], display: [] } as IStrtegicMapFilter,
		map_params: {
			plan: 1,
			year: 1,
			display: 1,
			cirteria: undefined,
			method: undefined,
		},
		map_data: {} as IStrategicMapData,
		plan: "",
		planTable: "",
		map: null,
		longdo: null,
		map_legend: [],
		isInit: true,
		exportReport: {
			reportType: "",
			reportTypeNo: null,
			reportName: "",
			reportPlan: "",
			reportId: 0,
		},
		retryCreateLine: 0,
	}),
	actions: {
		async getdata(id: number) {
			const service = new StrategicsService()
			const res = await service.getDashboard(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.data = res.data

				this.setPlan()
			}
		},
		async getMapfilter(id: number) {
			const service = new StrategicsService()
			const res = await service.getMapFilter(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.map_filter = res.data

				switch (true) {
					case res.data.criteria.length > 0:
						this.map_params.cirteria = res.data?.criteria[0]?.id
						break
					case res.data.display.length > 0:
						this.map_params.display = res.data?.display[0]?.id
						break
					case res.data.plan.length > 0:
						this.map_params.plan = res.data?.plan[0].id
						break
					case res.data.year.length > 0:
						this.map_params.year = res.data?.year[0]
						break
					case res.data.method.length > 0:
						this.map_params.method = res.data?.method[0]?.id
						break
				}
			}
		},
		async getMapData(id: number) {
			const params: IMapDataReq = {
				plan: this.map_params.plan,
				year: this.map_params.year,
				display: this.map_params.display,
				criteria: this.map_params.cirteria,
			}

			if (this.map_params.display === DISPLAY_METHOD) {
				params.method = this.map_params.method
			}

			const service = new StrategicsService()
			const res = await service.getMap(id, params)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.map_data = res.data
				// if()
				this.createLine()
				// if (this.isInit) {
				// 	setTimeout(() => {
				// 		res.data.items?.forEach((item) => {
				// 			this.createLine(item, item.color)
				// 		})
				// 	}, 1500)
				// 	this.isInit = false
				// }

				// res.data.items?.forEach((item) => {
				// 	this.createLine(item, item.color)
				// })
			}
		},
		async setFavorite(id: number) {
			const service = new StrategicsService()
			const res = await service.createFavorite(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		updateGraph(lineChart: Ref, barChart1: Ref, barChart2: Ref) {
			if (Object.keys(this.data).length > 0) {
				lineChart?.value?.updateOptions({
					xaxis: {
						categories: this.getGraphCategories1,
						title: {
							text: "ปี พ.ศ.",
						},
						// tickAmount: 10,
					},
					legend: {
						position: "bottom",
						offsetY: 10,
						horizontalAlign: "center",
						showForSingleSeries: true,
						customLegendItems: this.getGraph1LegendLabel,
						itemMargin: {
							horizontal: 15,
						},
						onItemClick: {
							toggleDataSeries: true,
						},
						onItemHover: {
							highlightDataSeries: true,
						},
					},
				})

				barChart1?.value?.updateOptions({
					xaxis: {
						categories: this.getBarChartCategories1,
						title: {
							text: "ปี พ.ศ.",
						},
						// tickAmount: 10,
					},
					legend: {
						position: "bottom",
						offsetY: 10,
						horizontalAlign: "center",
						showForSingleSeries: true,
						customLegendItems: this.getGraph1LegendLabel,
						itemMargin: {
							horizontal: 15,
						},
						onItemClick: {
							toggleDataSeries: true,
						},
						onItemHover: {
							highlightDataSeries: true,
						},
					},
				})

				barChart2?.value?.updateOptions({
					xaxis: {
						categories: this.getBarChartCategories2,
						title: {
							text: "ปี พ.ศ.",
						},
						// tickAmount: 10,
					},
					legend: {
						itemMargin: {
							horizontal: 15,
						},
						showForSingleSeries: true,
						customLegendItems: this.getLegendBar2,
						markers: {
							radius: 12,
							// fillColors: ["#FDB833", "#66abe6"],
						},
						onItemClick: {
							toggleDataSeries: true,
						},
						onItemHover: {
							highlightDataSeries: true,
						},
					},
				})
			}
		},
		togglePlan(planName: string) {
			this.plan = planName
		},
		togglePlanTable(planName: string) {
			this.planTable = planName
		},
		setPlan() {
			// ตั้งค่า init-plan สำหรับ  กราฟตัวสุดท้าย
			const planWithKeyword = this.data.bar2.datasets.find((item) => item.plan && item.plan.includes("แผน"))
			const firstAvailablePlan = this.data.bar2.datasets.find((item) => item.plan)

			this.plan = planWithKeyword?.plan || firstAvailablePlan?.plan || ""
		},
		setMap(map: any) {
			this.map = map
			// @ts-ignore
			this.longdo = window.longdo
		},
		// createLine(item: IStrategicMapDataItem, color: string) {
		// 	if (this.map) {
		// 		const longdo = this.longdo
		// 		const mapData = item?.the_geom?.coordinates?.map((val) => ({ lon: val[0], lat: val[1] }))
		// 		const detail = this.createHtmlPopUp(item)
		// 		const line = new longdo.Polyline(mapData, { lineColor: color, detail })

		// 		this.map.Overlays.add(line)
		// 		this.map.location({ lon: mapData[0]?.lon, lat: mapData[0]?.lat })
		// 		this.map.zoom(14)
		// 	}
		// },
		createLine() {
			if (this.map) {
				this.map.Overlays.clear()
				if (this.map_data.items.length > 0) {
					const lines = this.map_data.items.map((item) => {
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

					const geom =
						this.map_data.items[0].the_geom.coordinates.length > 0
							? this.map_data.items[0].the_geom.coordinates[0]
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
				if (this.retryCreateLine <= 5) {
					setTimeout(() => {
						this.createLine()
					}, 500)
				}

				this.retryCreateLine++
			}
		},
		createHtmlPopUp(item: IStrategicMapDataItem) {
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
      <div class="row">
        <div class="col-4 text-gray-800 mb-2">สายทาง:</div>
        <div class="col-6 text-gray-700 mb-2">${item.road_name}</div>
        <div class="col-4 text-gray-800 mb-2">ปีที่:</div>
        <div class="col-6 text-gray-700 mb-2">${item.year}</div>
        <div class="col-4 text-gray-800 mb-2">IRI ก่อนซ่อม:</div>
        <div class="col-6 text-gray-700 mb-2">${toNumber(item.iri_before, 2)}  ม./กม.</div>
        <div class="col-4 text-gray-800 mb-2">IRI หลังซ่อม:</div>
        <div class="col-6 text-gray-700 mb-2">${toNumber(item.iri_after, 2)} ม./กม.</div>
    </div>`

			return html
		},
		async onMapSelected(id: number) {
			this.map.Overlays.clear()
			await this.getMapData(id)
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
				const keys = Object.keys(data).filter((key) => data[key as keyof typeof data] !== undefined)
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
		getBar1Color(index: number) {
			const data = this.data?.bar1

			return data?.color[index] || ""
		},
	},
	getters: {
		getRoadsList(state) {
			const { data } = state ?? {}
			const roadList = data?.road

			return roadList || []
		},
		getFilter(state) {
			const { data } = state ?? {}
			const filterInfo = data?.filter

			return `ชนิดผิวทาง : ${filterInfo?.surface_type ?? "-"}, ช่องจราจร : ${filterInfo?.lane ?? "-"}, จัดกลุ่ม : ${
				filterInfo?.km ?? "-"
			} กม.`
		},
		getFilterList(state) {
			const filterList = state.data?.filter?.filter
			if (!filterList || !state.data) {
				return ""
			}

			const result = filterList?.map((condition) => {
				return condition
					.split(" ")
					.map((part, index) => {
						if (index === 1 || index === 3) {
							return part
						}

						const num = Number(part)
						return isNaN(num) ? part : num.toFixed(Number.isInteger(num) ? 0 : 2)
					})
					.join(" ")
			})

			return result
		},
		getDashboardConditions(state) {
			const condtions = state.data?.condition

			return `เงื่อนไข : ${condtions?.condition ?? "-"}, เป้าหมาย: ${
				condtions?.target ?? "-"
			}, อัตราคิดลด (Discount Rate) : ${
				condtions?.discount || condtions?.discount === 0 ? condtions?.discount + "%" : "-"
			} `
		},
		getComment(state) {
			const comment = state.data?.comment

			const regex = /(\n|<br\s*\/?>)/
			const check = regex.test(comment)

			return check ? comment.split("\n") : [comment]
		},
		getGraphCategories1(state) {
			const data = state.data?.graph1
			const categories = data?.lable.map((item) => (item + 543).toString())

			return categories || []
		},
		getGraph1LegendLabel(state) {
			const data = state.data?.graph1
			const legend = data?.line

			return legend || []
		},
		getBarChartCategories1(state) {
			const data = state.data?.bar1
			const categories = data?.lable.map((item) => (item + 543).toString())

			return categories || []
		},
		getBarChartCategories2(state) {
			const data = state.data?.bar2
			const categories = data?.lable?.map((item) => (item + 543).toString())

			return categories || []
		},
		getLegendBar2(state) {
			const data = state.data?.bar2
			const filterData = data?.datasets?.filter((item) => item?.plan.includes(String(state.plan)))
			const legend = filterData?.flatMap((item) => item?.data)?.flatMap((item) => item?.lable)

			return legend || []
		},
		getLineChartSeries(state) {
			const data = state.data?.graph1
			const dataset = data?.line.map((parent, index) => {
				return { name: parent, data: data?.value[index] }
			})

			const result = dataset

			return result || []
		},
		getBar1Series(state) {
			const data = state.data?.bar1
			const dataset = data?.datasets.map((item) => {
				return { name: item.lable, data: item.value }
			})

			const result = dataset

			return result || []
			// return []
		},
		getBar1ColorList(state) {
			const data = state.data?.bar1
			const color = data?.color

			return color || []
		},
		getBar1Label(state) {
			const data = state.data?.bar1?.lable

			return data || []
		},
		getBar1LegendLabel(state) {
			const data = state.data?.bar1?.datasets?.map((item) => item.lable)

			return data || []
		},
		getBar2Series(state) {
			const data = state.data?.bar2
			if (!data) {
				return []
			}

			const globalLabels = data?.lable
			const datasets = data?.datasets

			const planData = datasets?.find((item) => item.plan === state.plan)?.data
			if (!planData || !globalLabels) {
				return []
			}

			const acc = new Map()

			planData.forEach((yearData, yearIndex) => {
				yearData?.lable.forEach((label, index) => {
					let existing = acc.get(label)
					if (!existing) {
						existing = new Array(globalLabels.length).fill(null)
						acc.set(label, existing)
					}
					existing[yearIndex] = Number((yearData?.value[index] || 0).toFixed(2))
				})
			})

			const dataset = Array.from(acc, ([name, data]) => ({ name, data }))

			return dataset
		},
		getPlanTable(state) {
			const data = state.data.table

			switch (state.planTable) {
				case "แผนที่ 1":
					return data?.plan_1
				case "แผนที่ 2":
					return data?.plan_2
				case "แผนที่ 3":
					return data?.plan_3
			}
		},
		getSummaryTable(state) {
			const data = state.data?.table?.summary
			const newData = data?.map((item) => {
				const itemData = item?.data
				return {
					name: item.name,
					plans: itemData[0],
					plans_2: itemData?.filter((_, index) => {
						return index > 0
					}),
				}
			})

			return newData
		},
		getSummaryPlanTable(state) {
			let planData: IStrategicDashboardTableUnlimitedPlan[] = []
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
				years = Array.from(new Set(planData.flatMap((item) => item.data.map((child) => child.year + 543)))).sort(
					(a, b) => a - b
				)
			}

			return { plan: planData, years }
		},
		getCalculateColspan(state) {
			return state.data?.table?.summary[0]?.data?.length
		},
		getMapPlanOptions(state) {
			const mapFilter = state.map_filter
			const plans = mapFilter.plan
			const options = plans.map((item) => ({ label: item.name, value: item.id }))

			return options || []
		},
		getMapYearOptions(state) {
			const mapFilter = state.map_filter
			const years = mapFilter.year
			const options = years.map((year) => ({ label: `${year}`, value: year }))

			return options || []
		},
		getMapDisplayOptions(state) {
			const mapFilter = state.map_filter
			const display = mapFilter.display
			const options = display.map((item) => ({ label: item.name, value: item.id }))

			return options || []
		},
		getMapMethodOptions(state) {
			const mapFilter = state.map_filter
			const method = mapFilter.method
			const options = method.map((item) => ({ label: item.name, value: item.id }))

			return options || []
		},
		getCriteriaOptions(state) {
			const mapFilter = state.map_filter
			const criteria = mapFilter.criteria
			const options = criteria.map((item) => ({ label: item.name, value: item.id }))

			return options || []
		},
		getCriteriaLegend(state) {
			const criteriaMethod = state.map_data.criteria_method
			// const criteria = mapData?.items?.map((item) => ({ name: item.title, color: item.color }))
			return criteriaMethod || []
		},
		getMapLegend(state) {
			const mapFilter = state.map_data
			const criteria = mapFilter?.criteria_method

			// const method = mapFilter.method?.map((item) => ({ name: item.name, color: item.color }))

			return criteria || []
		},
		getGradeAC(state) {
			if (!state.map_params.cirteria) {
				return state.map_filter.criteria.length > 0 ? state.map_filter.criteria[0].grade : []
			}
			return (
				state.map_filter.criteria.find((e) => {
					return e.id === state.map_params.cirteria
				})?.grade ?? []
			)
		},
		getGradeCC(state) {
			if (!state.map_params.cirteria) {
				return state.map_filter.criteria.length > 0 ? state.map_filter.criteria[0].grade_cc : []
			}
			return (
				state.map_filter.criteria.find((e) => {
					return e.id === state.map_params.cirteria
				})?.grade_cc ?? []
			)
		},
		getSymbolBoxHeight(state) {
			if (state.map_params.display === DISPLAY_METHOD) {
				return state.map_data?.criteria_method?.length > 0
					? state.map_data?.criteria_method?.length > 5
						? "230"
						: state.map_data?.criteria_method?.length * 30 + 80
					: "80"
			} else {
				// return state.map_filter.criteria[0].grade.length * 30 + state.map_filter.criteria[0].grade_cc.length * 30 + 155

				const data = state.map_filter.criteria.find((e) => e.id === state.map_params.cirteria)
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
			return state.map_data?.criteria_method?.length > 5 ? 205 : 155
		},
		getSymbolBoxPositionTop(state) {
			if (state.map_params.display === DISPLAY_METHOD) {
				return state.map_data?.criteria_method?.length > 0
					? state.map_data?.criteria_method?.length > 5
						? "210"
						: state.map_data?.criteria_method?.length * 30 + 60
					: "55"
			} else {
				const data = state.map_filter.criteria.find((e) => e.id === state.map_params.cirteria)
				if (data) {
					return data.grade?.length * 30 + data.grade_cc?.length * 30 + 110
				} else {
					return 350
				}
			}
		},
		getGradeCCPositionTop(state) {
			const data = state.map_filter.criteria.find((e) => e.id === state.map_params.cirteria)
			if (data) {
				return data.grade?.length * 30 + 70
			} else {
				return 196
			}
		},
	},
})
