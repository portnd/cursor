import { IDashboardCondition, IHighChart, DashboardConditionService } from "../infrastructure"
import { useDashboardReflectiveStore, useDashboardRoadConditionStore } from "../store"
import { useDashboardStore } from "~/core/modules/dashboard/store/DashboardStore"
interface IState {
	loading: boolean
	menu: string
	data: IDashboardCondition
	// map: any
	// longdo: any
	// arrAsset: number[]
	// year: number
	// owners: string[]
	// roads: string[]
	roads: string[]
	totalData: number
	dataArray: number[]
	// conditions: any[]
	// reflect: any[]
	toggle: {
		hist: boolean
		whiteLine: boolean
		dataType: string
		ownerID: number
		conditionType: number
	}
	graphSelectFilter: number[]
	dataGraph: IHighChart
	conditionList: any[]
	conditionType: number
	conditionTypeString: string
	surveyRule: number
	lane: number
	labelsArr: string[]
	// colors: string[]
	// displayParam: number
	// loading: boolean
	// pavementAge: IPavementAge[]
	// triggerSearch: boolean
	// roadDropdownList: IMaintenanceHistoryRoadDropdownList[]
	reloadData: boolean
	conditionArray: number[]
	conditionColors: Array<any>
	conditionLabel: Array<any>
	colors: Array<any>
	labels: Array<any>
}

export const useDashboardConditionStore = defineStore("dashboard/condition", {
	state: (): IState => ({
		loading: false,
		menu: "asset",
		data: {} as IDashboardCondition,
		// map: null,
		// longdo: null,
		// arrAsset: [11, 12],
		// year: 2567,
		// owners: [],
		roads: [],
		totalData: 0,
		// dataArray: [44, 55, 13, 43, 22],
		dataArray: [],
		// conditions: [
		// 	{
		// 		id: 1,
		// 		name: "ดีมาก",
		// 		color: "#A4FCA5",
		// 	},
		// 	{
		// 		id: 2,
		// 		name: "ดี",
		// 		color: "#42D235",
		// 	},
		// 	{
		// 		id: 3,
		// 		name: "ปานกลาง",
		// 		color: "#F77A14",
		// 	},
		// 	{
		// 		id: 4,
		// 		name: "แย่",
		// 		color: "#FF290A",
		// 	},
		// 	{
		// 		id: 5,
		// 		name: "แย่มาก",
		// 		color: "#973131",
		// 	},
		// ],
		// reflect: [
		// 	{
		// 		id: 6,
		// 		name: "ผ่าน",
		// 		color: "#42D235",
		// 	},
		// 	{
		// 		id: 7,
		// 		name: "ไม่ผ่าน",
		// 		color: "#FF290A",
		// 	},
		// ],
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
		conditionList: [],
		conditionType: 1,
		conditionTypeString: "",
		surveyRule: 1,
		lane: 1,
		labelsArr: [],
		// displayParam: 1,
		// loading: false,
		// pavementAge: [
		// 	{ name: "0 - 2 ปี", color: "#50CD89" } as IPavementAge,
		// 	{ name: "3 - 5 ปี", color: "#87C442" } as IPavementAge,
		// 	{ name: "6 - 10 ปี", color: "#FDB833" } as IPavementAge,
		// 	{ name: "มากกว่า 10 ปี", color: "#DC3545" } as IPavementAge,
		// ],
		// triggerSearch: false,
		// roadDropdownList: [],
		reloadData: false,
		conditionArray: [1, 2, 3, 4, 5],
		conditionColors: [],
		conditionLabel: [],
		colors: [],
		labels: [],
	}),
	actions: {
		getRoads() {
			this.loading = true
			const dashboardStore = useDashboardStore()

			this.roads = dashboardStore.params.road_id

			this.loading = false
		},
		async getCondition() {
			const initData = useInitData()
			const ownerId = initData.conditionGrade().length > 0 ? initData.conditionGrade()[0]?.id : null
			const dashboardCondition = useDashboardRoadConditionStore()

			const dashboardStore = useDashboardStore()

			const { roadIds, depotCodes } = dashboardStore.getFirstRoadIdsAndDepotCodesWhenOwnerAccess()

			this.roads = roadIds

			// แปลงค่า km จาก "00+000" เป็นตัวเลข (เมตร) เหมือนกับ tab อื่น
			const kmStart = Number.isNaN(convertStringToKm(dashboardStore.params.km_start))
				? null
				: convertStringToKm(dashboardStore.params.km_start)
			const kmEnd = Number.isNaN(convertStringToKm(dashboardStore.params.km_end))
				? null
				: convertStringToKm(dashboardStore.params.km_end)

			// ใช้ roadIds จาก access helper แทน raw params เพื่อรองรับ owner access
			const year = Number(dashboardStore.params.year)

			const dashboardReflectiveStore = useDashboardReflectiveStore()
			let refConditionRangeId: number | null = 0
			refConditionRangeId =
				this.conditionType === 5 ? dashboardReflectiveStore.params.owner_id : dashboardCondition.params.owner_id

			const params = ref("")
			params.value = `condition_type=${this.conditionType}`
			if (roadIds.length) {
				params.value += `&road_id=${roadIds}`
			}
			if (year) {
				params.value += `&year=${year}`
			}
			if (depotCodes.length) {
				params.value += `&depot_code=${depotCodes}`
			}
			if (kmStart !== null) {
				params.value += `&km_start=${kmStart}`
			}
			if (kmEnd !== null) {
				params.value += `&km_end=${kmEnd}`
			}
			if (refConditionRangeId) {
				params.value += `&condition_owner_id=${refConditionRangeId}`
			} else {
				params.value += `&condition_owner_id=${ownerId}`
				dashboardCondition.params.owner_id = ownerId
			}

			const dashboardConditionService = new DashboardConditionService()
			const res = await dashboardConditionService.getCondition(params.value)

			if (!res.status) {
				// useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.data = res.data
				this.conditionArray = this.data.chart?.color?.map((_: string, index: number) => index)
				this.conditionColors = this.data.chart?.color?.map((item: any) => item)
				this.conditionLabel = this.data.chart?.lable?.map((item: any) => item)
				this.dataArray = this.data.chart?.data?.map((item: any) => item)
				// this.totalData = this.dataArray?.data((acc: any, item: any) => acc + item)
			}

			// await this.getCondition_map();
		},
		// getCheckBox() {
		// 	switch (this.conditionType) {
		// 		case 1:
		// 			this.graphSelectFilter = [1, 2, 3, 4]
		// 			this.conditionArray = [1, 2, 3, 4]
		// 			this.dataArray = [44, 55, 13, 32]
		// 			this.conditionList = [
		// 				{
		// 					id: 1,
		// 					label: "เรียบมาก",
		// 					value: 1,
		// 					color: "#42D235",
		// 				},
		// 				{
		// 					id: 2,
		// 					label: "เรียบ",
		// 					value: 2,
		// 					color: "#A4FCA5",
		// 				},
		// 				{
		// 					id: 3,
		// 					label: "ขรุขระ",
		// 					value: 3,
		// 					color: "#F77A14",
		// 				},
		// 				{
		// 					id: 4,
		// 					label: "ขรุขระมาก",
		// 					value: 4,
		// 					color: "#FF290A",
		// 				},
		// 			]
		// 			break
		// 		case 2:
		// 			this.graphSelectFilter = [1, 2, 3]
		// 			this.conditionArray = [1, 2, 3]
		// 			this.dataArray = [44, 55, 13]
		// 			this.conditionList = [
		// 				{
		// 					id: 1,
		// 					label: "หยาบ (ดีมาก)",
		// 					value: 1,
		// 					color: "#A4FCA5",
		// 				},
		// 				{
		// 					id: 2,
		// 					label: "ปานกลาง",
		// 					value: 2,
		// 					color: "#F77A14",
		// 				},
		// 				{
		// 					id: 3,
		// 					label: "ละเอียด (แย่มาก)",
		// 					value: 3,
		// 					color: "#FF290A",
		// 				},
		// 			]
		// 			break
		// 		case 3:
		// 			this.graphSelectFilter = [1, 2, 3, 4]
		// 			this.conditionArray = [1, 2, 3, 4]
		// 			this.dataArray = [44, 55, 13, 32]
		// 			this.conditionList = [
		// 				{
		// 					id: 1,
		// 					label: "ตื้นมาก",
		// 					value: 1,
		// 					color: "#42D235",
		// 				},
		// 				{
		// 					id: 2,
		// 					label: "ตื้น",
		// 					value: 1,
		// 					color: "#A4FCA5",
		// 				},
		// 				{
		// 					id: 3,
		// 					label: "ลึก",
		// 					value: 3,
		// 					color: "#F77A14",
		// 				},
		// 				{
		// 					id: 4,
		// 					label: "ลึกมาก",
		// 					value: 4,
		// 					color: "#FF290A",
		// 				},
		// 			]
		// 			break
		// 		case 4:
		// 		case 5:
		// 			this.graphSelectFilter = [6, 7]
		// 			this.conditionArray = [1, 2]
		// 			this.dataArray = [44, 23]
		// 			this.conditionList = [
		// 				{
		// 					id: 1,
		// 					label: "ผ่าน",
		// 					value: 6,
		// 					color: "#42D235",
		// 				},
		// 				{
		// 					id: 2,
		// 					label: "ไม่ผ่าน",
		// 					value: 7,
		// 					color: "#FF290A",
		// 				},
		// 			]
		// 			break
		// 	}
		// 	this.colors = this.conditionList.map((item) => item.color)
		// 	this.labelsArr = this.conditionList.map((item) => item.label)
		// },
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
		histSeries(): any {
			const series = [{ name: "IRI", data: this.dataArray }]
			return series
		},
		pieOptions(): any {
			if (this.data.chart) {
				this.totalData = this.data.chart.data?.reduce(
					(accumulator: any, currentValue: any) => accumulator + currentValue,
					0
				)
			}

			const option: any = {
				chart: {
					events: {
						dataPointSelection: (_: any, __: any, opts: any) => {
							// const dataIndex = opts.dataPointIndex
							// const roadId = this.data?.chart?.road_id[dataIndex]
							// // const roadGroupId = this.roadsData?.road?.road_group_id

							// navigateTo(`roads?road_id=${roadId}`)

							this.filterMap(_, __, opts)
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
					text: "ข้อมูลสภาพทาง",
					align: "center",
					style: {
						fontSize: "16px",
					},
				},
				dataLabels: {
					style: {
						fontSize: "10px",
						colors: ["#fff"],
					},
					formatter: (_: any, opts: any) =>
						((opts.w.config.series[opts.seriesIndex] / this.totalData) * 100).toFixed(1) + " %",
				},
				colors: [],
				labels: [],
				tooltip: {
					enabled: true,
					fillSeriesColor: false,
					theme: "light",
					y: {
						show: true,
						formatter: (value: number) => {
							return Number(value.toFixed(2)).toLocaleString() + " กม."
						},
					},
				},
				legend: {
					show: false,
				},
			}

			if (this.colors.length > 0) {
				option.colors = this.colors
			} else {
				option.colors = this.conditionColors ?? []
			}
			if (this.labels.length > 0) {
				option.labels = this.labels
			} else {
				option.labels = this.conditionLabel ?? []
			}
			return option
		},
		pieSeries(): any {
			const result = this.conditionArray
				?.map((id: number) => {
					const value = this.dataArray[id]
					return value !== undefined ? Number(value.toFixed(2)) : undefined
				})
				.filter((item) => item !== undefined)

			return result ?? []
		},
		barSeries(): any {
			const final = [
				{
					name: "ระยะทาง",
					data: <Array<any>>[],
				},
			]
			if (this.conditionArray) {
				const result = this.conditionArray
					.map((id: number) => {
						const value = this.dataArray[id]
						return value !== undefined ? Number(value.toFixed(2)) : undefined
					})
					.filter((item) => item !== undefined)

				final[0].data = result ?? []
			}

			return final ?? []
		},
		barOption(): any {
			const option: any = {
				chart: {
					toolbar: {
						show: false,
					},
					events: {
						dataPointSelection: (_: any, __: any, opts: any) => {
							// const dataIndex = opts.dataPointIndex
							// const roadId = this.data?.chart?.road_id[dataIndex]
							// // const roadGroupId = this.roadsData?.road?.road_group_id

							// navigateTo(`roads?road_id=${roadId}`)
							this.filterMap(_, __, opts)
						},
					},
				},
				dataLabels: {
					enabled: true,
					offsetY: -20,
					style: {
						fontSize: "10px",
						colors: ["#304758"],
					},
				},
				tooltip: {
					enabled: true,
					y: {
						show: true,
						formatter: (value: number) => {
							const total = this.dataArray.reduce((acc, item) => acc + item)
							return ((value / total) * 100).toFixed(1) + " %"
						},
					},
				},
				stroke: {
					show: true,
					width: 2,
					colors: ["transparent"],
					curve: "smooth",
				},
				legend: {
					show: false,
				},
				title: {
					text: "ข้อมูลสรุปสภาพทาง",
					align: "center",
					style: {
						fontSize: "16px",
					},
				},
				colors: [],
				plotOptions: {
					bar: {
						horizontal: false,
						columnWidth: "50",
						distributed: true,
						endingShape: "rounded",
						borderRadius: 4,
						dataLabels: {
							position: "top",
						},
					},
				},
				xaxis: {
					categories: [],
					title: {
						text: "เกณฑ์การประเมิน",
					},
					labels: {
						style: {
							colors: "#304758",
							fontSize: "12px",
						},
					},
				},
				yaxis: [
					{
						title: {
							text: "ระยะทาง (กม.)",
						},
						labels: {
							formatter: (value: number) => {
								if (value) {
									return Number(value.toFixed(2)).toLocaleString()
								} else {
									return 0
								}
							},
						},
					},
				],
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
				fill: {
					opacity: 1,
				},
			}
			if (this.colors.length > 0) {
				option.colors = this.colors
			} else {
				option.colors = this.conditionColors ?? []
			}
			if (this.labels.length > 0) {
				option.xaxis.categories = this.labels
			} else {
				option.xaxis.categories = this.conditionLabel ?? []
			}

			return option
		},
		filterMap(_: any, __: any, opts: any) {
			const dashboardStore = useDashboardStore()

			const dataIndex = opts.dataPointIndex
			// const roadId = this.data?.chart?.road_id[dataIndex]
			// const roadGroupId = this.roadsData?.road?.road_group_id

			const isSelected = opts.selectedDataPoints?.[0]?.[0] !== undefined

			if (isSelected) {
				// dashboardStore.map.Overlays.clear()
				dashboardStore.selectCondition = {
					status: true,
					color: this.data?.chart.color[dataIndex],
				}
				dashboardStore.selectCondition.status = true
				dashboardStore.selectCondition.color = this.data?.chart.color[dataIndex]
				dashboardStore.createLineConditionByColor()
			} else {
				dashboardStore.selectCondition.status = false

				dashboardStore.selectCondition = {
					status: false,
					color: "",
				}

				dashboardStore.createAssetInMap()
			}
		},
	},
	getters: {
		getYearOptions(state) {
			const { conditionList } = state
			const options = conditionList?.map((condition) => ({ label: `${condition.year + 543}`, value: condition.year }))

			return options.sort((a, b) => b.value - a.value) || []
		},
	},
})
