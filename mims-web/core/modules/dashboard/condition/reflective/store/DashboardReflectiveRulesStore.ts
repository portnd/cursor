import {
	ILine,
	IReflectiveListData,
	IReflectivityItem,
	IReflectivityDetails,
	DashboardReflectiveGraphService,
} from "../infrastructure"
import { IReflectivityGrade } from "~/core/modules/initData/infrastructure/data"
import { IOption } from "~/core/shared/types/Option"
import { IHighChart } from "~/core/modules/dashboard/infrastructure"
import { IReflectivityRule, IReflectivityRuleList } from "~/composables/useReflectivityRule"
import { useDashboardStore } from "~/core/modules/dashboard/store/DashboardStore"

interface IStateGeom {
	geom: string
	color: string
}

interface IStateChildDataTable {
	id: number
	km_start: string
	km_end: string
	value: number
	geom_cl: string
	line_type_id: number
	grade_id: number
	expanded: boolean
}

interface IStateDataTable {
	id: string
	km_start: string
	km_end: string
	value: number
	geom_cl: string
	line_type_id: number
	child: IStateChildDataTable[]
	grade_id: number
	expanded: boolean
}

interface IStateParams {
	graph_type: string
	line_color: string
	year: number | null
	id_parent: number | null
	owner_id: number | null
	ref_owner_id: number | null
	owner_name: string
	ref_reflectivity_range_id: number | null
	criteria_id: number[]
	line_type_id: number[]
	old_owner_id: number | null
	toggle_id: number | null
}

interface IState {
	loading: boolean
	params: IStateParams
	map: any
	lineList: ILine[]
	reflectList: IReflectiveListData[]
	details: IReflectivityDetails
	reflectivityGrade: IReflectivityGrade[]
	allReflectivityInput: IReflectivityRule[]
	reflectivityInput: IReflectivityRule[]
	refReflectivityGrade: IReflectivityRule[]
	criteriaOptions: IOption[]
	oldRefOwnerId: number | null
	reflectivityDataTable: IStateDataTable[]
	lineChart: IHighChart
	geoms: IStateGeom[]
}

export const useDashboardReflectiveRulesStore = defineStore("dashboard/reflectiveStore", {
	state: (): IState => ({
		loading: false,
		params: {
			graph_type: "graph",
			line_color: "",
			year: null,
			id_parent: null,
			owner_id: null,
			owner_name: "",
			ref_reflectivity_range_id: null,
			ref_owner_id: null,
			criteria_id: [],
			line_type_id: [],
			old_owner_id: null,
			toggle_id: null,
		},
		map: null,
		lineList: [],
		reflectList: [],
		details: {} as IReflectivityDetails,
		reflectivityGrade: JSON.parse(JSON.stringify(useInitData()?.reflectivityGrade())),
		allReflectivityInput: [],
		reflectivityInput: [],
		refReflectivityGrade: [],
		criteriaOptions: [],
		oldRefOwnerId: null,
		reflectivityDataTable: [],
		lineChart: {
			title: {
				text: "",
			},
			chart: {
				width: 650,
				height: 300,
				type: "line",
				zoomType: "x",
				stacked: true,
				toolbar: {
					show: false,
				},
				events: {},
			},
			dataLabels: {
				enabled: false,
			},
			stroke: {
				curve: "smooth",
			},
			xAxis: {
				categories: [],
				tickInterval: 0,
				margin: 10,
				title: {
					text: null,
					style: {
						fontSize: "12.5px",
						color: "#3f4254",
					},
				},
				labels: {
					enable: false,
					style: {
						fontSize: "12px",
						color: "#3f4254",
					},
				},
			},
			yAxis: {
				title: {
					text: " (m/km)",
					style: {
						fontSize: "12.5px",
						color: "#3f4254",
					},
				},
			},
			tooltip: {
				useHTML: true,
				borderRadius: 10,
				backgroundColor: "#fff",
				borderColor: "none",
				padding: 8,
				formatter: function (this: any): string {
					return (
						// "กม. ที่ " +
						this.x +
						"<div style='border-bottom: 1px solid #dedede; width: 100%;' class='mt-1 mb-4' ></div>" +
						" <span class='rounded' style='width:10px; height: 10px;  display:block; background: " +
						this.color +
						";'></span>" +
						"<span>" +
						": " +
						this.y +
						"</span>"
					)
				},
			},
			plotOptions: {
				series: {
					turboThreshold: 10000,
					marker: {
						enabled: false,
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

			series: [],
		} as IHighChart,
		geoms: [],
	}),
	actions: {
		setMap(map: any) {
			this.map = map

			this.createLine()
			this.defaultLocation()
		},
		async getLineList() {
			this.loading = true
			const dashboard = useDashboardStore()
			const roadId = Number(dashboard.params.road_id[0])

			const service = new DashboardReflectiveGraphService()
			const res = await service.getLaneList(roadId)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.lineList = res.data
			}
		},
		async getReflectivityList() {
			this.loading = true
			const dashboard = useDashboardStore()
			const roadId = Number(dashboard.params.road_id[0])

			const service = new DashboardReflectiveGraphService()
			const res = await service.getReflectivityList(roadId)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.reflectList = res.data
			}
		},
		setDefault() {
			const initData = useInitData()
			const refReflectivityGrade = initData?.reflectivityGrade() ?? []
			this.params.year = this.getYearOptions.length > 0 ? this.getYearOptions[0]?.value : null
			this.params.id_parent = this.getLineOptions.length > 0 ? this.getLineOptions[0]?.value : null
			this.params.criteria_id =
				this.criteriaOptions.length > 0 ? this.criteriaOptions.map((item) => Number(item.value)) : []
			this.params.line_type_id = initData?.refStripeType()?.map((item) => item.id) ?? []
			this.params.owner_id = refReflectivityGrade?.length > 0 ? refReflectivityGrade[0]?.id : null
			this.params.owner_name = refReflectivityGrade?.length > 0 ? refReflectivityGrade[0]?.owner_name : ""
			this.params.ref_owner_id = JSON.parse(JSON.stringify(this.params.owner_id))
			this.params.ref_reflectivity_range_id =
				refReflectivityGrade.length > 0 ? refReflectivityGrade[0]?.ref_reflectivity_range_id : null
			this.params.old_owner_id = JSON.parse(JSON.stringify(this.params.owner_id))

			console.log("this.params.owner_id ----->", this.params.owner_id)

			nextTick()
			this.setReflectInput()
		},
		async getReflectivityDetails() {
			if (this.params.id_parent && this.params.ref_reflectivity_range_id) {
				this.loading = true

				const service = new DashboardReflectiveGraphService()
				const res = await service.getReflectivityDetails(this.params.id_parent, this.params.ref_reflectivity_range_id)

				this.loading = false

				if (!res.status) {
					useHandlerError(res.code, res.error, { showToast: true })
				} else {
					this.details = res.data
					this.params.line_color = this.details.has_white_line ? "white" : this.details.has_yellow_line ? "yellow" : ""
				}
			}
		},
		toggleGraph(event: any) {
			const innerType = event.target.innerHTML.toLowerCase()
			this.params.graph_type = innerType

			this.createDataTable()
		},
		async onUpdateIdParent() {
			const { getReflectivityDetails, createDataTable } = this
			await getReflectivityDetails().then(() => createDataTable())
		},
		setReflectInput() {
			const { reflectivityGrade } = this

			const defaultData: IReflectivityGrade[] = JSON.parse(JSON.stringify(reflectivityGrade))

			const allReflectivityInput = defaultData.flatMap((parent) => {
				const reflectRule = useReflectivityRule()

				reflectRule.forEach((rule) => {
					rule.id = parent.id

					if (rule.id === parent.id) {
						rule.name = parent.owner_name
						rule.white.reflectivity_list.forEach((whiteRule, index) => {
							const white = parent.road_line?.white[index]
							if (whiteRule.grade_id === white.grade?.id) {
								whiteRule.left_value = white.left_value
								whiteRule.right_value = white.right_value
							}
						})

						rule.yellow?.reflectivity_list?.forEach((yellowRule, index) => {
							const yellow = parent.road_line?.yellow[index]
							if (yellowRule.grade_id === yellow?.grade?.id) {
								yellowRule.left_value = yellow?.left_value
								yellowRule.right_value = yellow?.right_value
							}
						})
					}
				})

				return reflectRule ?? []
			})

			this.allReflectivityInput = JSON.parse(JSON.stringify(allReflectivityInput))
			this.setRefReflectivityGrade()
		},
		setRefReflectivityGrade() {
			this.reflectivityInput = JSON.parse(JSON.stringify(this.allReflectivityInput))

			this.refReflectivityGrade = JSON.parse(JSON.stringify(this.reflectivityInput))
			handleFieldReflectivityCondition(this.reflectivityInput, this.params.owner_name)
		},
		checkReflectivitiyInput(surfaceType: string) {
			const checkDataType = this.reflectivityInput.find((item) => item.id === this.params.owner_id)

			if (surfaceType === "white") {
				return {
					left_unit: checkDataType?.leftUnit ?? "",
					right_unit: checkDataType?.rightUnit ?? "",
					reflectivity_list: checkDataType?.white?.reflectivity_list ?? [],
				}
			} else {
				return {
					left_unit: checkDataType?.leftUnit ?? "",
					right_unit: checkDataType?.rightUnit ?? "",
					reflectivity_list: checkDataType?.yellow?.reflectivity_list ?? [],
				}
			}
		},
		onUpdateOwner() {
			console.log("onUpdateOwner ----->")

			const initData = useInitData()
			this.params.ref_reflectivity_range_id =
				initData?.reflectivityGrade()?.find((item) => item.id === this.params.owner_id)?.ref_reflectivity_range_id ??
				null

			console.log("this.params.ref_reflectivity_range_id ----->", this.params.ref_reflectivity_range_id)

			this.params.owner_name =
				initData?.reflectivityGrade()?.find((item) => item.id === this.params.owner_id)?.owner_name ?? ""
			console.log("this.params.owner_name ----->", this.params.owner_name)
			this.setRefReflectivityGrade()
		},
		checkCriteria(gradeId: number) {
			const { criteriaOptions, reflectivityGrade, params } = this
			if (!gradeId && criteriaOptions.length === 0) {
				return
			}

			const checkRange =
				reflectivityGrade.find((item) => item.ref_reflectivity_range_id === params.ref_reflectivity_range_id)
					?.road_line[(this.params.line_color as "white") || "yellow"] ?? []
			const checkId = checkRange?.find((item) => item.grade.id === gradeId)

			const result = {
				name: checkId?.grade?.name,
				color: checkId?.grade?.color,
			}

			return result
		},
		createCriteriaCheckbox() {
			const { reflectivityGrade, params } = this
			const filterRange =
				reflectivityGrade.find((grade) => grade.ref_reflectivity_range_id === params.ref_reflectivity_range_id)
					?.road_line[(this.params.line_color as "white") || "yellow"] ?? []

			const criteria =
				filterRange.map((item) => ({
					name: item.grade?.name,
					color: item.grade?.color,
					id: item.grade?.id,
				})) ?? []

			this.params.criteria_id = criteria.length > 0 ? criteria.map((item) => item.id) : []
			const options = criteria.map((item) => ({ label: item.name, value: item.id, color: item.color })) ?? []

			this.criteriaOptions = options ?? []
		},
		defaultInput() {
			const { reflectivityGrade, params } = this

			const defaultData: IReflectivityGrade[] = JSON.parse(JSON.stringify(reflectivityGrade))
			const reflectivityList = defaultData.find((item) => item.id === params.ref_owner_id)?.road_line

			this.reflectivityInput.forEach((item) => {
				if (item.id === params.ref_owner_id) {
					item.white.reflectivity_list.forEach((white, index) => {
						if (white.grade_id === reflectivityList?.white[index]?.grade?.id) {
							white.left_value = reflectivityList?.white[index]?.left_value
							white.right_value = reflectivityList?.white[index]?.right_value
						}
					})

					item.yellow.reflectivity_list.forEach((yellow, index) => {
						if (yellow.grade_id === reflectivityList?.yellow[index]?.grade?.id) {
							yellow.left_value = reflectivityList?.yellow[index]?.left_value
							yellow.right_value = reflectivityList?.yellow[index]?.right_value
						}
					})
				}
			})
		},
		cancelRule() {
			this.params.owner_id = JSON.parse(JSON.stringify(this.params.old_owner_id))
			this.params.ref_owner_id = JSON.parse(JSON.stringify(this.params.owner_id))
			const checkInit = useInitData()
				?.reflectivityGrade()
				?.find((item) => item.id === this.params.owner_id)

			if (checkInit) {
				this.params.ref_reflectivity_range_id = checkInit.ref_reflectivity_range_id
			}

			const allInput: IReflectivityRule[] = JSON.parse(JSON.stringify(this.allReflectivityInput))
			const filterAllInput = allInput?.filter((item) => item.id === this.params.ref_owner_id)
			this.reflectivityInput = filterAllInput
			this.setRefReflectivityGrade()
		},
		async submitRule() {
			this.refReflectivityGrade = JSON.parse(JSON.stringify(this.reflectivityInput))
			this.allReflectivityInput.forEach((input) => {
				const conditionGrade = JSON.parse(JSON.stringify(this.refReflectivityGrade))
				conditionGrade.forEach((grade: IReflectivityRule) => {
					if (input.id === grade.id && input.name === grade.name) {
						input.white = JSON.parse(JSON.stringify(grade.white))
						input.yellow = JSON.parse(JSON.stringify(grade.yellow))
					}
				})
			})

			this.params.criteria_id = this.criteriaOptions.map((item) => Number(item.value))
			this.params.old_owner_id = JSON.parse(JSON.stringify(this.params.owner_id))
			this.params.ref_owner_id = JSON.parse(JSON.stringify(this.params.owner_id))
			this.params.ref_reflectivity_range_id =
				useInitData()
					?.reflectivityGrade()
					?.find((item) => item.id === this.params.ref_owner_id)?.ref_reflectivity_range_id ?? null

			await this.getReflectivityDetails()
				.then(() => this.createCriteriaCheckbox())
				.then(() => this.createDataTable())
			this.oldRefOwnerId = JSON.parse(JSON.stringify(this.params.owner_id))
		},
		onUpdateCheckbox() {
			this.createDataTable()
		},
		async createDataTable() {
			this.reflectivityDataTable = []

			const { refReflectivityGrade, params, details, processDataTable, createFeature } = this

			const masterRule =
				refReflectivityGrade.find((reflect) => reflect?.id === params.ref_owner_id) ?? ({} as IReflectivityRule)
			const reflectivityItems =
				details.datas?.find((data) => data.color.toLowerCase() === params.line_color)?.items ?? []
			const generatedData = processDataTable(reflectivityItems, masterRule)
			// @ts-ignore กลับมาแก้ด้วย
			this.reflectivityDataTable = generatedData.filter((item) => {
				return params.ref_reflectivity_range_id === 2 || (item?.child && item?.child?.length > 0)
			})

			await nextTick()
			createFeature()
		},
		processDataTable(reflectivityItems: IReflectivityItem[], masterRule: IReflectivityRule) {
			let globalChildId = 0

			const getRules = (lineColor: string) => {
				return lineColor === "white" ? masterRule.white?.reflectivity_list : masterRule.yellow?.reflectivity_list
			}

			const evaluateExpression = (rules: IReflectivityRuleList[], value: number) => {
				// if (value === null) {
				// 	return null
				// }
				return (
					rules.find((rule) => {
						if (
							rule.left_value !== null &&
							rule.right_value !== null &&
							this.params.criteria_id.includes(rule.grade_id)
						) {
							const expression = `${rule.left_value} ${rule.left_symbol} ${value} && ${value} ${rule.right_symbol} ${rule.right_value}`

							return eval(expression) && this.params.criteria_id.includes(rule.grade_id)
						}
						return false
					}) || null
				)
			}

			const processData = (items: IReflectivityItem[], lineColor: string) => {
				const rules = getRules(lineColor) ?? []
				return items
					.map((parent, index) => {
						const matchingRuleForParent = evaluateExpression(rules, parent.retro_avg)

						if (this.params.ref_reflectivity_range_id === 2 && !matchingRuleForParent) {
							return null
						}

						const children =
							parent.items
								?.map((child) => {
									const matchingRuleForChild = evaluateExpression(rules, child.retro_avg)
									if (!matchingRuleForChild) {
										return null
									}

									return {
										id: globalChildId++,
										km_start: convertMeterToKm(child.km_start),
										km_end: convertMeterToKm(child.km_end),
										value: child.retro_avg === null ? 0 : Number(child.retro_avg.toFixed(2)),
										geom_cl: child.geom_cl,
										line_type_id: child.ref_stripe_type_id,
										grade_id: matchingRuleForChild.grade_id,
										expanded: false,
									}
								})
								.filter((child) => child !== null && this.params.line_type_id.includes(child.line_type_id)) ?? []

						return {
							id: this.params.ref_reflectivity_range_id === 2 ? index : `parent${index}`,
							km_start: convertMeterToKm(parent.km_start),
							km_end: convertMeterToKm(parent.km_end),
							value: parent.retro_avg === null ? 0 : Number(parent.retro_avg.toFixed(2)),
							geom_cl: parent.geom_cl,
							line_type_id: parent.ref_stripe_type_id,
							child: children,
							grade_id: matchingRuleForParent ? matchingRuleForParent.grade_id : 0,
							expanded: false,
						}
					})
					.filter((parent) => parent && this.params.line_type_id.includes(parent.line_type_id))
			}

			return processData(reflectivityItems, this.params.line_color)
		},
		toggleData(lineColor: string) {
			if (lineColor !== this.params.line_color) {
				this.params.line_color = lineColor
				this.createDataTable()
			}
		},
		setLocation(latLon: string[], id: number) {
			// this.image.expanded = true

			if (latLon.length > 0) {
				const location = latLon
				this.map.location({
					lon: location[0],
					lat: location[1],
				})
			}

			this.params.toggle_id = id
			// this.image.imageID = imageID
		},
		graph() {
			const graphData = this.reflectivityDataTable
			const setLocation = this.setLocation
			const type = this.params.line_color === "white" ? "เส้นสีขาว" : "เส้นสีเหลือง"
			const items = graphData
			const refRangeId = this.params.ref_reflectivity_range_id

			let labelXAxis = items.flatMap((parent) => {
				if (this.params.ref_reflectivity_range_id === 2) {
					return parent.km_start
				} else {
					return parent.child.map((child) => child.km_start)
				}
			})

			// Deduplicate
			labelXAxis = [...new Set(labelXAxis)]

			if (labelXAxis.length < 2) {
				labelXAxis = items.map((parent) => parent.km_start)
			}

			// Calculate tickInterval once
			const tickInterval = labelXAxis.length < 105 ? 0 : 35

			// const type = this.toggle.dataType

			const graph = {
				title: {
					text: "",
				},
				chart: {
					height: 300,
					type: "line",
					zoomType: "x",
					events: {
						selection: function (this: any) {
							const geom = graphData
								.flatMap((parent) => parent.child.map((child) => (this.x === child.id ? child.geom_cl : undefined)))
								.filter((item) => item !== undefined)
							const latlon = geom[0]?.split(",")[0].split("(")[1].split(" ")
							setLocation(latlon ?? [], this.x)
						},
					},
					stacked: true,
					toolbar: {
						show: false,
					},
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
				dataLabels: {
					enabled: false,
				},
				stroke: {
					curve: "smooth",
				},
				xAxis: {
					categories: labelXAxis,
					tickInterval,
					margin: 10,
					title: {
						text: null,
						style: {
							fontSize: "12.5px",
							color: "#3f4254",
						},
					},
					labels: {
						enable: false,
						style: {
							fontSize: "12px",
							color: "#3f4254",
						},
					},
				},
				yAxis: {
					title: {
						text: "Retro Reflectivity Average",
						style: {
							fontSize: "12.5px",
							color: "#3f4254",
						},
					},
					labels: {
						style: {
							fontSize: "12px",
						},
					},
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
                <span>${type}: <b class="fw-semibold">${this.y}</b></span>
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
							events: {
								click: function (this: any) {
									const geom =
										refRangeId === 2
											? graphData.filter((item) => item.id === this.x).map((item) => item.geom_cl)
											: graphData
													.flatMap((parent) => {
														return parent.child.map((child) => {
															if (this.x === child.id) {
																return child.geom_cl
															}

															return undefined
														})
													})
													.filter((item) => item !== undefined)

									const latlon = geom[0]?.split(",")[0].split("(")[1].split(" ")
									setLocation(latlon ?? [], this.x)
								},
							},
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
				series: this.createSeries() ?? [],
				// series: [],
			}

			this.lineChart = graph as any
		},
		createSeries() {
			const values = this.extractValues()
			const filterRefGrade = this.refReflectivityGrade.filter((item) => item.id === this.params.ref_owner_id)
			const zones = filterRefGrade.flatMap((item) => {
				return this.reflectivityDataTable.flatMap((data) => {
					if (this.params.ref_reflectivity_range_id === 2) {
						return item[(this.params.line_color as "white") || "yellow"]?.reflectivity_list.map((reflect) => {
							return {
								value: reflect.right_value,
								color: this.colorRule(reflect.grade_id, this.params.line_color.toLowerCase()),
							}
						})
					} else {
						return item[(this.params.line_color as "white") || "yellow"]?.reflectivity_list.flatMap((_) => {
							return data.child?.flatMap((_) => {
								return item[(this.params.line_color as "white") || "yellow"]?.reflectivity_list.map((reflect) => {
									return {
										value: reflect.right_value,
										color: this.colorRule(reflect.grade_id, this.params.line_color.toLowerCase()),
									}
								})
							})
						})
					}
				})
			})

			let uniqueZones = Array.from(new Set(zones.map((zone) => JSON.stringify(zone)))).map((str) => JSON.parse(str))
			uniqueZones = uniqueZones.reverse()

			return [
				{
					name: this.params.line_color,
					data: values,
					zoneAxis: "y",
					zones: uniqueZones,
					animation: false,
				},
			]
		},
		extractValues() {
			if (this.params.ref_reflectivity_range_id === 2) {
				return this.reflectivityDataTable.map((item) => item.value) ?? []
			} else {
				return this.reflectivityDataTable.flatMap((item) => item.child.map((child) => child.value)) ?? []
			}
		},
		colorRule(gradeId: number, lineColor: string) {
			const { reflectivityGrade, params } = this
			const color = reflectivityGrade
				.find((range) => range.id === params.ref_owner_id)
				?.road_line[(lineColor as "white") || "yellow"]?.find((item) => item.grade?.id === gradeId)?.grade?.color

			return color ?? ""
		},
		calculateMax() {
			const { reflectivityDataTable, params } = this
			let minKm = Infinity
			let maxKm = -Infinity

			reflectivityDataTable.forEach((item) => {
				if (params.ref_reflectivity_range_id === 2) {
					minKm = Math.min(minKm, convertStringToKm(item.km_start), convertStringToKm(item.km_end))
					maxKm = Math.max(maxKm, convertStringToKm(item.km_start), convertStringToKm(item.km_end))
				} else {
					item.child.forEach((child) => {
						minKm = Math.min(minKm, convertStringToKm(child.km_start), convertStringToKm(child.km_end))
						maxKm = Math.max(maxKm, convertStringToKm(child.km_start), convertStringToKm(child.km_end))
					})
				}
			})

			const result = Math.abs(maxKm - minKm) / 1000

			return result
		},
		histSeries() {
			const { reflectivityDataTable, params } = this
			const distanceByGrade: { [key: string]: number } = {}

			reflectivityDataTable.forEach((item) => {
				let distances
				if (params.ref_reflectivity_range_id === 2) {
					distances = [{ km_start: item.km_start, km_end: item.km_end, grade_id: item.grade_id }]
				} else {
					distances = item.child.map((child) => ({
						km_start: child.km_start,
						km_end: child.km_end,
						grade_id: child.grade_id,
					}))
				}

				distances.forEach((e) => {
					const startKm = convertStringToKm(e.km_start)
					const endKm = convertStringToKm(e.km_end)
					const distance = Math.abs(startKm - endKm)
					distanceByGrade[e.grade_id] = (distanceByGrade[e.grade_id] || 0) + distance
				})
			})

			const data = this.criteriaOptions.map((key) => ({
				grade_id: Number(key.value),
				distance: distanceByGrade[Number(key.value)] ?? 0,
			}))

			const fitlerCriteria = data.filter((data) => params.criteria_id.includes(data.grade_id)) || []

			const result = [
				{
					name: params.line_color === "white" ? "เส้นสีขาว" : "เส้นสีเหลือง",
					data: fitlerCriteria.map((data) => (data.distance / 1000).toFixed(2)),
				},
			]

			return result || []
		},
		histPercent() {
			const series = this.histSeries()
			const data =
				this.details.datas
					?.filter((item) => item.color.toLowerCase() === this.params.line_color)
					?.flatMap((item) => item.items) ?? []

			if (data.length === 0) {
				return []
			}

			const totalDistance = Math.abs(data[0]?.km_start - data[data.length - 1]?.km_end) / 1000

			const result = series.flatMap((parent) => {
				return parent.data.map((value) => {
					return (Number(value) / totalDistance) * 100
				})
			})

			return result ?? []
		},
		histGraph() {
			const { params, criteriaOptions } = this
			const criteriaArray = criteriaOptions.reduce((accumulator, current) => {
				const id = Number(current.value)

				if (this.params.criteria_id.includes(id)) {
					accumulator.push(current.label)
				}

				return accumulator
			}, [] as string[])

			const distance = this.calculateMax()
			if (criteriaArray.length > 0) {
				const colors = criteriaOptions.reduce((acc: string[], criteria) => {
					if (params.criteria_id.includes(Number(criteria.value))) {
						acc.push(criteria.color ?? "")
					}

					return acc
				}, [])

				const percent = this.histPercent() ?? []

				const graph = {
					chart: {
						zoom: {
							enabled: false,
						},
						stacked: true,
						toolbar: {
							show: false,
						},
					},
					dataLabels: {
						enabled: true,
						distributed: true,
						offsetY: -30,
						formatter: function (value: any, opt: any) {
							return `${value} กม. (${percent[opt.dataPointIndex]?.toFixed(2)}%)`
						},
						style: {
							position: "absolute",
							fontSize: "12px",
							colors: ["#000000"],
							fontWeight: "regular",
						},
					},
					stroke: {
						curve: "smooth",
					},
					legend: {
						show: false,
					},
					xaxis: {
						type: "string",
						categories: criteriaArray,
						tickPlacement: "on",
						title: {
							offsetY: -10,
							text: "Retro Reflectivity Average",
							style: {
								color: "#3F4254",
								fontSize: "12px",
								fontWeight: 400,
							},
						},
					},
					yaxis: {
						min: 0,
						max: distance,
						tickAmount: 4,
						title: {
							text: "ระยะทาง (กม.)",
						},
						labels: {
							formatter: function (val: number) {
								const km = val.toFixed(1)
								return `${km}`
							},
						},
					},
					colors,
					plotOptions: {
						bar: {
							horizontal: false,
							columnWidth: `80px`,
							endingShape: "rounded",
							distributed: true,
							borderRadius: 4,
							dataLabels: {
								position: "top",
								hideOverflowingLabels: false,
							},
						},
					},
					markers: {
						size: 5,
						hover: {
							size: 9,
						},
					},
					tooltip: {
						Html: true,
						enabled: true,
						custom: function ({ dataPointIndex }: any) {
							let html = `<div class="apexcharts-result">`

							html += `<div><span class="dot" style="background-color: ${colors[dataPointIndex]};"></span> ${
								criteriaArray[dataPointIndex]
							}: ${percent[dataPointIndex]?.toFixed(2)}%</div>`

							html += "</div>"

							return html
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
					fill: {
						opacity: 1,
					},
				}

				return graph
			}
		},
		getGeoms() {
			const { reflectivityDataTable, reflectivityGrade, params } = this
			this.geoms = []

			if (reflectivityDataTable.length === 0) {
				return []
			}

			const getColorForGrade = (gradeId: number) => {
				const gradeColor = reflectivityGrade
					.find((item) => item.ref_reflectivity_range_id === params.ref_reflectivity_range_id)
					?.road_line[(params.line_color as "white") || "yellow"]?.find((item) => item.grade?.id === gradeId)
					?.grade?.color
				return gradeColor || ""
			}

			const filterDataTable =
				params.ref_reflectivity_range_id === 2
					? reflectivityDataTable.filter(
							(item) => params.criteria_id.includes(item.grade_id) && params.line_type_id.includes(item.line_type_id)
					  ) ?? []
					: reflectivityDataTable.filter((item) =>
							item.child.filter(
								(child) =>
									params.criteria_id.includes(child.grade_id) && params.line_type_id.includes(child.line_type_id)
							)
					  ) ?? []

			const parentGeoms = filterDataTable.map((data) => ({
				geom: data.geom_cl,
				color: getColorForGrade(data.grade_id),
			}))

			const childGeoms = filterDataTable.flatMap((data) =>
				(data.child || []).map((child) => ({
					geom: child.geom_cl,
					color: getColorForGrade(child.grade_id),
				}))
			)

			this.geoms = params.ref_reflectivity_range_id === 2 ? [...parentGeoms, ...childGeoms] : childGeoms ?? []
		},
		createLine() {
			const { geoms } = this
			this.map?.Overlays.clear()

			if (geoms.length === 0) {
				return
			}

			// @ts-ignore
			const longdo = window.longdo

			if (longdo) {
				const polylines = geoms.flatMap((geom) => {
					return longdo?.Util.overlayFromWkt(geom.geom, { lineColor: geom.color })
				})

				if (polylines.length > 0) {
					polylines.forEach((line) => {
						this.map?.Overlays.add(line)
					})
				}
			}
		},
		defaultLocation() {
			if (this.geoms.length > 0) {
				// default location
				const location = getLatLong(this.geoms[0]?.geom)
				this.map?.location({
					lon: location.lon,
					lat: location.lat,
				})
			}

			this.map?.zoom(15)
		},
		createFeature() {
			const { getGeoms, graph, histPercent, histSeries, createLine, setToggleId } = this

			setToggleId()
			getGeoms()
			graph()
			histPercent()
			histSeries()
			createLine()
		},
		toggleDataTable(item: IStateChildDataTable) {
			const geomCl = item.geom_cl
			const location = getLatLong(geomCl)
			this.params.toggle_id = item.id

			this.map.location({
				lon: location.lon,
				lat: location.lat,
			})
		},
		toggleParentDataTable(item: IStateDataTable, index: number) {
			if (this.params.ref_reflectivity_range_id === 2) {
				const geomCl = item.geom_cl
				const location = getLatLong(geomCl)

				this.map.location({
					lon: location.lon,
					lat: location.lat,
				})

				this.params.toggle_id = Number(item.id)
			} else {
				this.reflectivityDataTable[index].expanded = !this.reflectivityDataTable[index].expanded
			}
		},
		async callBackUpdateData(method: string) {
			const {
				getReflectivityList,
				getReflectivityDetails,
				createCriteriaCheckbox,
				createDataTable,
				defaultLocation,
				setDefault,
			} = this

			this.details = {} as IReflectivityDetails
			this.reflectivityDataTable = []

			await getReflectivityList()

			if (method === "create" || method === "delete") {
				setDefault()
			} else {
				// ดักกรณี update ข้อมูลแล้วมีการเปลี่ยนปี แล้วให้เรียกตัวเดิม
				const findYear = this.reflectList.find((parent) =>
					parent.items.some((child) => child.id_parent === this.params.id_parent)
				)

				if (findYear) {
					this.params.year = findYear.year
				}
			}

			await getReflectivityDetails()
				.then(() => createCriteriaCheckbox())
				.then(() => createDataTable())
				.then(() => {
					defaultLocation()
				})
		},
		setToggleId() {
			this.params.toggle_id =
				this.reflectivityDataTable.length > 0 && this.params.ref_reflectivity_range_id === 2
					? Number(this.reflectivityDataTable[0]?.id)
					: this.reflectivityDataTable[0]?.child[0]?.id
		},
		onUpdateYear() {
			this.params.id_parent = this.getLineOptions.length > 0 ? this.getLineOptions[0].value : null
			nextTick(() => this.getReflectivityDetails().then(() => this.createDataTable()))
		},
	},
	getters: {
		getYearOptions(state) {
			const { reflectList } = state
			const options = reflectList.map((item) => ({ label: `${item.year + 543}`, value: item.year }))

			return options || []
		},
		getLineOptions(state) {
			const { reflectList, params } = state
			const items = reflectList.find((item) => item.year === params.year)?.items ?? []
			const options = items.map((item) => ({
				label: `${item.line_no} (สำรวจ: ${buddhistFormatDate(item.surveyed_date, "dd mmm yy")})`,
				value: item.id_parent,
			}))

			return options || []
		},
		getUserUpdated(state) {
			const { details } = state
			const user = details.updated_by

			return {
				username: user?.user_name ?? "",
				fullname: user?.full_name ?? "",
				department: user?.department?.name ?? "",
				img_path: user?.profile_picture ?? "",
				date: buddhistFormatDate(details.updated_date, "dd mmm yy เวลา HH:ii น.") ?? "",
			}
		},
		getOwnerOptions(state) {
			const { reflectivityGrade } = state
			const options = reflectivityGrade.map((grade) => ({ label: grade.owner_name, value: grade.id }))

			return options || []
		},
		getLineTypeOptions(_) {
			const data = useInitData()?.refStripeType()
			const options = data?.map((item) => ({ label: item.name_th, value: item.id }))

			return options || []
		},
		getLineListOptions(state) {
			const { lineList } = state
			const options = lineList.map((item) => ({ label: item.line_no.toString(), value: item.line_no }))
			return options || []
		},
	},
})
