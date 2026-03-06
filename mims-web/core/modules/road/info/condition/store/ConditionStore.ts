import { IRoadConditionDetails, IRoadConditionList, RoadConditionService } from "../infrastructure"
import { IRoadConditionItem, IHighChart } from "../infrastructure/RoadConditionModel"
import { IOption } from "~/core/shared/types/Option"
import { ConditionList, ISurveyRule } from "~/composables/useSurveyRule"
import { IConditionGrade, IConditionSurfaceType } from "~/core/modules/initData/infrastructure/data"

interface IStateTypeIds {
	type: string
	ref_condition_range_id: number[]
	owner_id: number[]
}

interface IStateGeom {
	geom: string
	color: string
}

interface IStateImage {
	expanded: boolean
	imageID: number
	path: string
	play: any
	speed: number
	pause: Function | null
	playing: boolean
	lastImageID: number | null
}

interface IPlayBtn {
	play: boolean
	speed2x: boolean
	speed4x: boolean
	stop: boolean
	pause: boolean
}

interface IStateChildDataTable {
	id: number
	km_start: string
	km_end: string
	value: number
	image: string
	geom_cl: string
	surface_type: string
	grade_id: number
	expanded: boolean
}

interface IStateDataTable {
	id: string
	km_start: string
	km_end: string
	value: number
	geom_cl: string
	surface_type: string
	child: IStateChildDataTable[]
	grade_id: number
	expanded: boolean
}

interface IStateParams {
	year: number | null
	id_parent: number | null
	old_owner_id: number | null
	owner_id: number | null
	ref_condition_range_id: number | null
	condition_type: string
	graph_type: string
	criteria_id: number[]
}

interface IState {
	loading: boolean
	map: any
	params: IStateParams
	image: IStateImage
	playBtn: IPlayBtn
	refTypeRangeIds: IStateTypeIds[]
	conditionList: IRoadConditionList[]
	conditionGrade: IConditionGrade[]
	conditionInput: ISurveyRule[]
	conditionDataTable: IStateDataTable[]
	refConditionGrade: ISurveyRule[]
	allConditionInput: ISurveyRule[]
	details: IRoadConditionDetails
	mapHeight: string
	criteriaOptions: IOption[]
	lineChart: IHighChart
	geoms: IStateGeom[]
	oldRefOwnerId: number | null
	graphAvg: number
}

export const useConditionStore = defineStore("condition", {
	state: (): IState => ({
		loading: false,
		map: null,
		params: {
			year: null,
			id_parent: null,
			old_owner_id: null,
			owner_id: null,
			ref_condition_range_id: null,
			criteria_id: [],
			condition_type: "IRI",
			graph_type: "GRAPH",
		},
		image: {
			expanded: false,
			imageID: 0,
			path: "",
			play: null,
			speed: 1000,
			pause: null,
			playing: false,
			lastImageID: null,
		},
		playBtn: {
			play: false,
			speed2x: false,
			speed4x: false,
			stop: false,
			pause: false,
		},
		conditionDataTable: [],
		conditionList: [],
		conditionGrade: useInitData()?.conditionGrade() ? JSON.parse(JSON.stringify(useInitData()?.conditionGrade())) : [],
		conditionInput: [],
		refConditionGrade: [],
		allConditionInput: [],
		details: {} as IRoadConditionDetails,
		mapHeight: "calc(60vh)",
		criteriaOptions: [],
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
		oldRefOwnerId: null,
		refTypeRangeIds: [],
		graphAvg: 0,
	}),
	actions: {
		async getConditionList(roadId: number) {
			this.loading = true

			const service = new RoadConditionService()
			const res = await service.getConditionList(roadId)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.conditionList = res.data
			}
		},
		async getConditionDetails() {
			if (this.params.id_parent && this.params.ref_condition_range_id) {
				this.loading = true

				this.details = {} as IRoadConditionDetails

				const service = new RoadConditionService()
				const res = await service.getConditionDetails(this.params.id_parent, this.params.ref_condition_range_id)

				this.loading = false

				if (!res.status) {
					useHandlerError(res.code, res.error, { showToast: true })
				} else {
					this.details = res.data
				}
			}
		},
		setDefaultParams() {
			const { params, conditionList, criteriaOptions, setConditionInput } = this
			const initData = useInitData()

			params.year = conditionList.length > 0 ? conditionList.sort((a, b) => b.year - a.year)[0]?.year : null
			params.id_parent =
				params.year !== null && conditionList[0]?.items?.length > 0 ? conditionList[0]?.items[0]?.id_parent : null
			this.params.criteria_id = criteriaOptions.length > 0 ? criteriaOptions.map((item) => Number(item.value)) : []
			this.params.owner_id = initData.conditionGrade().length > 0 ? initData.conditionGrade()[0]?.id : null
			this.params.old_owner_id = JSON.parse(JSON.stringify(this.params.owner_id))
			this.params.ref_condition_range_id =
				initData.conditionGrade().length > 0 ? initData.conditionGrade()[0]?.ref_condition_range_id : null

			setConditionInput()
		},
		setMap(map: any) {
			this.map = map
			this.createLine()
			this.defaultLocation()
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
		setConditionInput() {
			const { conditionGrade } = this

			const defaultData: IConditionGrade[] = JSON.parse(JSON.stringify(conditionGrade))

			const allConditionInput = defaultData.flatMap((parent) => {
				const surveyRule = useSurveyRule(parent.ref_condition_range_id)

				surveyRule.forEach((rule) => {
					const matchedData = parent.condition_list.find((data) => data.condition_type === rule.name)?.surface_type

					rule.id = parent.id

					if (rule.id === parent.id) {
						rule.ac.conditionList.forEach((acRule, index) => {
							const ac = matchedData?.ac[index]
							if (ac && acRule.grade_id === ac.grade?.id) {
								acRule.left_value = ac.left_value
								acRule.right_value = ac.right_value
							}
						})

						rule.cc?.conditionList?.forEach((ccRule, index) => {
							const cc = matchedData?.cc?.[index]
							if (cc && ccRule.grade_id === cc.grade?.id) {
								ccRule.left_value = cc.left_value
								ccRule.right_value = cc.right_value
							}
						})
					}
				})

				return surveyRule ?? []
			})

			this.allConditionInput = JSON.parse(JSON.stringify(allConditionInput))
			this.setRefConditionGrade()
		},
		setRefConditionGrade() {
			this.conditionInput = JSON.parse(
				JSON.stringify(
					this.allConditionInput.filter(
						(item) => item.id === this.params.owner_id && item.name === this.params.condition_type
					)
				)
			)
			this.refConditionGrade = JSON.parse(JSON.stringify(this.conditionInput))
			handleFieldCondition(this.conditionInput, this.params.condition_type)
		},
		defaultInput() {
			const { conditionGrade, params } = this

			const defaultData: IConditionGrade[] = JSON.parse(JSON.stringify(conditionGrade))
			const conditionList = defaultData.find((item) => item.id === params.owner_id)?.condition_list

			this.conditionInput.forEach((parent) => {
				if (parent.id === params.owner_id) {
					parent.ac.conditionList.forEach((ac, index) => {
						const matched = conditionList?.find((item) => item.condition_type === parent.name)?.surface_type
						if (ac.grade_id === matched?.ac[index]?.grade?.id) {
							ac.left_value = matched.ac[index]?.left_value
							ac.right_value = matched.ac[index]?.right_value
						}
					})
					parent.cc.conditionList.forEach((cc, index) => {
						const matched = conditionList?.find((item) => item.condition_type === parent.name)?.surface_type
						if (cc.grade_id === matched?.cc[index]?.grade?.id) {
							cc.left_value = matched.cc[index]?.left_value
							cc.right_value = matched.cc[index]?.right_value
						}
					})
				}
			})

			handleFieldCondition(this.conditionInput, this.params.condition_type)
		},
		createCriteriaCheckbox() {
			const { conditionGrade, params } = this
			const filterRange =
				conditionGrade.find((grade) => grade.ref_condition_range_id === params.ref_condition_range_id)
					?.condition_list ?? []
			const filterConditionList =
				filterRange.find((range) => range.condition_type === params.condition_type)?.surface_type ??
				({} as IConditionSurfaceType)
			const criteria =
				filterConditionList.ac?.map((item) => ({
					name: item.grade?.name,
					color: item.grade?.color,
					id: item.grade?.id,
				})) ?? []

			this.params.criteria_id = criteria.length > 0 ? criteria.map((item) => item.id) : []
			const options = criteria.map((item) => ({ label: item.name, value: item.id, color: item.color })) ?? []

			this.criteriaOptions = options ?? []
		},
		cancelRule() {
			this.params.owner_id = JSON.parse(JSON.stringify(this.params.old_owner_id))
			const checkInit = useInitData()
				.conditionGrade()
				.find((item) => item.id === this.params.owner_id)

			if (checkInit) {
				this.params.ref_condition_range_id = checkInit.ref_condition_range_id
			}

			const allInput: ISurveyRule[] = JSON.parse(JSON.stringify(this.allConditionInput))
			this.conditionInput = allInput.filter(
				(item) => item.id === this.params.owner_id && item.name === this.params.condition_type
			)

			this.setRefConditionGrade()
		},
		async submitRule() {
			this.refConditionGrade = JSON.parse(JSON.stringify(this.conditionInput))
			this.allConditionInput.forEach((input) => {
				const conditionGrade = JSON.parse(JSON.stringify(this.refConditionGrade))
				conditionGrade.forEach((grade: ISurveyRule) => {
					if (input.id === grade.id && input.name === grade.name) {
						input.ac = JSON.parse(JSON.stringify(grade.ac))
						input.cc = JSON.parse(JSON.stringify(grade.cc))
					}
				})
			})

			await this.getConditionDetails()
				.then(() => this.createCriteriaCheckbox())
				.then(() => this.createDataTable())
			this.mapHeight = this.params.ref_condition_range_id === 1 ? "calc(60vh)" : "calc(97vh)"
			this.params.criteria_id = this.criteriaOptions.map((item) => Number(item.value))
			this.params.old_owner_id = JSON.parse(JSON.stringify(this.params.owner_id))
			this.oldRefOwnerId = JSON.parse(JSON.stringify(this.params.owner_id))
		},
		toggleData(event: Event) {
			const target = event.target as HTMLElement
			const innerType = target.innerHTML

			if (this.params.condition_type === innerType) {
				return
			}

			this.params.condition_type = innerType

			this.checkConditionType(innerType)
		},
		async checkConditionType(type: string) {
			this.setRefConditionGrade()

			if (!this.params.owner_id) {
				return
			}

			let shouldFetchConditionDetails = false

			const initData = useInitData()?.conditionGrade()
			const currentRefConditionRangeId =
				initData?.find((item) => item.id === this.params.owner_id)?.ref_condition_range_id ?? null

			if (this.oldRefOwnerId && this.oldRefOwnerId !== this.params.owner_id) {
				this.params.owner_id = this.oldRefOwnerId
				this.params.old_owner_id = JSON.parse(JSON.stringify(this.params.owner_id))
				shouldFetchConditionDetails = currentRefConditionRangeId !== 3
			}

			const typeRange = this.refTypeRangeIds.find((item) => item.type === type)
			const isOwnerIdIncluded = typeRange?.owner_id.includes(this.params.owner_id)

			if (typeRange) {
				if (!isOwnerIdIncluded && typeRange?.owner_id?.length > 0) {
					this.params.owner_id = typeRange.owner_id[0]
					this.params.old_owner_id = JSON.parse(JSON.stringify(this.params.owner_id))
					shouldFetchConditionDetails = true
				}

				const newRefConditionRangeId =
					initData?.find((item) => item.id === this.params.owner_id)?.ref_condition_range_id ?? null
				this.params.ref_condition_range_id = newRefConditionRangeId

				if (newRefConditionRangeId === 1 && currentRefConditionRangeId === 1) {
					shouldFetchConditionDetails = false
				}

				if (shouldFetchConditionDetails) {
					await this.getConditionDetails()
				}
				this.mapHeight = this.params.ref_condition_range_id === 1 ? "calc(60vh)" : "calc(97vh)"

				this.setRefConditionGrade()
				this.createCriteriaCheckbox()
				this.createDataTable()
			}
		},
	initTypeRangeId() {
		const dataType = this.getDataType

		if (!dataType || !Array.isArray(dataType) || dataType.length === 0) {
			return
		}

		const conditionTypeToOwnersMap = this.conditionGrade.reduce((acc: { [key: string]: Set<number> }, owner) => {
			if (!owner.condition_list) {
				return acc
			}
			owner.condition_list.forEach((condition) => {
				if (dataType.includes(condition.condition_type)) {
					if (!acc[condition.condition_type]) {
						acc[condition.condition_type] = new Set()
					}
					acc[condition.condition_type].add(owner.ref_condition_range_id)
				}
			})
			return acc
		}, {})
		const conditionTypeOwnerId = this.conditionGrade.reduce((acc: { [key: string]: Set<number> }, owner) => {
			if (!owner.condition_list) {
				return acc
			}
			owner.condition_list.forEach((condition) => {
				if (dataType.includes(condition.condition_type)) {
					if (!acc[condition.condition_type]) {
						acc[condition.condition_type] = new Set()
					}
					acc[condition.condition_type].add(owner.id)
				}
			})
			return acc
		}, {})

			const result = dataType.map((type) => ({
				type,
				ref_condition_range_id: Array.from(conditionTypeToOwnersMap[type] || []),
				owner_id: Array.from(conditionTypeOwnerId[type] || []),
			}))

			this.refTypeRangeIds = result ?? []
		},
		toggleGraph(event: Event) {
			const { params, createDataTable } = this
			const target = event.target as HTMLElement
			const innerType = target.innerHTML

			if (innerType === params.graph_type) {
				return
			}

			params.graph_type = innerType
			nextTick(() => createDataTable())
		},
		checkCriteria(gradeId: number) {
			const { criteriaOptions, conditionGrade, params } = this
			if (!gradeId && criteriaOptions.length === 0) {
				return
			}

			const checkRange =
				conditionGrade.find((item) => item.ref_condition_range_id === params.ref_condition_range_id)?.condition_list ??
				[]
			const checkId = checkRange
				?.find((item) => item.condition_type === params.condition_type)
				?.surface_type?.ac?.find((item) => item.grade.id === gradeId)

			const result = {
				name: checkId?.grade?.name,
				color: checkId?.grade?.color,
			}

			return result
		},
		checkConditionInput(surfaceType: string) {
			const checkDataType = this.conditionInput.find((condition) => condition.name === this.params.condition_type)
			if (surfaceType === "ac") {
				return {
					left_unit: checkDataType?.leftUnit ?? "",
					right_unit: checkDataType?.rightUnit ?? "",
					conditionList: checkDataType?.ac?.conditionList ?? [],
				}
			} else {
				return {
					left_unit: checkDataType?.leftUnit ?? "",
					right_unit: checkDataType?.rightUnit ?? "",
					conditionList: checkDataType?.cc?.conditionList ?? [],
				}
			}
		},
		async createDataTable() {
			this.conditionDataTable = []

			const { refConditionGrade, params, details, processDataTable, createFeature } = this

			const masterRule =
				refConditionGrade.find((condition) => condition?.name === params.condition_type) ?? ({} as ISurveyRule)
			const conditionItems =
				details.condition_datas?.find((condition) => condition.condition_type === params.condition_type)?.items ?? []
			const generatedData = processDataTable(conditionItems, masterRule)
			// @ts-ignore กลับมาแก้ด้วย
			this.conditionDataTable = generatedData.filter((item) => {
				return params.ref_condition_range_id === 3 || (item?.child && item?.child?.length > 0)
			})

			await nextTick()
			createFeature()
		},
		createFeature() {
			const { getGeoms, graph, createLine, histPercent, histSeries, setInitImage } = this
			setInitImage()
			getGeoms()
			graph()
			histPercent()
			histSeries()
			createLine()
		},
		processDataTable(conditionItems: IRoadConditionItem[], masterRule: ISurveyRule) {
			let globalChildId = 0

			const getRules = (surveyType: string) => {
				return surveyType === "AC" ? masterRule.ac?.conditionList : masterRule.cc?.conditionList
			}

			const evaluateExpression = (rules: ConditionList[], value: number) => {
				// if (value === null) {
				// 	return null
				// }
				return (
					rules.find((rule) => {
						if (rule.left_value !== null && rule.right_value !== null) {
							const expression = `${rule.left_value} ${rule.left_symbol} ${value} && ${value} ${rule.right_symbol} ${rule.right_value}`
							return eval(expression) && this.params.criteria_id.includes(rule.grade_id)
						}
						return false
					}) || null
				)
			}

			const processData = (items: IRoadConditionItem[], surveyType: string) => {
				const rules = getRules(surveyType) ?? []
				return items
					.filter((item) => item.survey_type === surveyType)
					.map((parent, index) => {
						const matchingRuleForParent = evaluateExpression(rules, parent.value)

						if (this.params.ref_condition_range_id === 3 && !matchingRuleForParent) {
							return null
						}

						const children =
							parent.items
								?.map((child) => {
									const matchingRuleForChild = evaluateExpression(rules, child.value)
									if (!matchingRuleForChild) {
										return null
									}
									return {
										id: globalChildId++,
										km_start: convertMeterToKm(child.km_start),
										km_end: convertMeterToKm(child.km_end),
										value: child.value === null ? 0 : Number(child.value.toFixed(2)),
										image: child.img_filepath,
										geom_cl: child.geom_cl,
										surface_type: child.survey_type,
										grade_id: matchingRuleForChild.grade_id,
										expanded: false,
									}
								})
								.filter((child) => child !== null) ?? []

						return {
							id: this.params.ref_condition_range_id === 3 ? index : `parent-${surveyType}-${index}`,
							km_start: convertMeterToKm(parent.km_start),
							km_end: convertMeterToKm(parent.km_end),
							value: parent.value === null ? 0 : Number(parent.value.toFixed(2)),
							geom_cl: parent.geom_cl,
							surface_type: parent.survey_type,
							child: children,
							grade_id: matchingRuleForParent ? matchingRuleForParent.grade_id : 0,
							expanded: false,
							km_start_number: parent.km_start,
						}
					})
					.filter((parent) => parent)
			}

			return [...processData(conditionItems, "AC"), ...processData(conditionItems, "CC")].sort((a, b) => {
				const startA = a?.km_start_number ?? 0 // ignore upper and lowercase
				const startB = b?.km_start_number ?? 0 // ignore upper and lowercase

				return startA - startB
			})
		},
		toggleDataTable(item: IStateChildDataTable) {
			const geomCl = item.geom_cl
			const location = getLatLong(geomCl)
			const imagePath = item.image
			this.image.path = imagePath
			this.image.imageID = item.id

			const isGeomNotEmpty = !geomCl.toLowerCase().includes("empty")
			if (isGeomNotEmpty) {
				this.map.location({
					lon: location.lon,
					lat: location.lat,
				})
			}
		},
		toggleParentDataTable(item: IStateDataTable, index: number) {
			if (this.params.ref_condition_range_id === 3) {
				const geomCl = item.geom_cl
				const location = getLatLong(geomCl)

				this.image.imageID = Number(item.id)

				this.map.location({
					lon: location.lon,
					lat: location.lat,
				})
			} else {
				this.conditionDataTable[index].expanded = !this.conditionDataTable[index].expanded
			}
		},
		onUpdateCheckbox() {
			this.createDataTable()
		},
		async onUpdateIdParent() {
			const { getConditionDetails, createDataTable } = this
			await getConditionDetails().then(() => createDataTable())
		},
		onUpdateYear() {
			this.params.id_parent = this.getSurveyLaneOptions.length > 0 ? this.getSurveyLaneOptions[0].value : null
			nextTick(() => this.getConditionDetails().then(() => this.createDataTable()))
		},
		getGeoms() {
			const { conditionDataTable, conditionGrade, params } = this
			this.geoms = []

			if (conditionDataTable.length === 0) {
				return []
			}

			const getColorForGrade = (gradeId: number, surfaceType: string) => {
				const gradeColor = conditionGrade
					.find((item) => item.ref_condition_range_id === params.ref_condition_range_id)
					?.condition_list?.find((condition) => condition.condition_type === params.condition_type)
					?.surface_type[(surfaceType as "ac") || "cc"]?.find((item) => item.grade?.id === gradeId)?.grade?.color
				return gradeColor || ""
			}

			const filterDataTable =
				params.ref_condition_range_id === 3
					? conditionDataTable.filter((item) => params.criteria_id.includes(item.grade_id)) ?? []
					: conditionDataTable.filter((item) =>
							item.child.filter((child) => params.criteria_id.includes(child.grade_id))
					  ) ?? []

			const parentGeoms = filterDataTable.map((data) => ({
				geom: data.geom_cl,
				color: getColorForGrade(data.grade_id, data.surface_type?.toLowerCase()),
			}))

			const childGeoms = filterDataTable.flatMap((data) =>
				(data.child || []).map((child) => ({
					geom: child.geom_cl,
					color: getColorForGrade(child.grade_id, child.surface_type?.toLowerCase()),
				}))
			)

			this.geoms = params.ref_condition_range_id === 3 ? [...parentGeoms, ...childGeoms] : childGeoms ?? []
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
		setGraphAvg(value: number) {
			if (!value || Number.isNaN(value)) {
				this.graphAvg = 0
			}

			this.graphAvg = value
		},
		graph() {
			const graphData = this.conditionDataTable
			const setLocation = this.setLocation
			const type = this.params.condition_type
			const items = graphData
			const setGraphAvg = this.setGraphAvg

			let labelXAxis = items.flatMap((parent) => {
				if (this.params.ref_condition_range_id === 3) {
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
						redraw: function (this: any) {
							const xExtremes = this.xAxis[0].getExtremes()

							const min = xExtremes.min
							const max = xExtremes.max

							this.series.forEach(function (series: any) {
								const filteredData = series.data.filter(function (point: any) {
									return point.x >= min && point.x <= max
								})

								const sum = filteredData.reduce(function (acc: number, point: any) {
									return acc + point.y
								}, 0)

								const average = sum / filteredData.length

								setGraphAvg(average)
							})
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
						text: this.toggleType(type, false),
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
									const geom = graphData
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
			}
			this.lineChart = graph as any
		},
		createSeries() {
			const values = this.extractValues()
			const filterRefGrade = this.refConditionGrade.filter((item) => item.name === this.params.condition_type)
			const zones = filterRefGrade.flatMap((item) => {
				return this.conditionDataTable.flatMap((data) => {
					if (this.params.ref_condition_range_id === 3) {
						return item[(data.surface_type.toLowerCase() as "ac") || "cc"]?.conditionList.map((condition) => {
							return {
								value: condition.right_value,
								color: this.colorRule(condition.grade_id, data.surface_type.toLowerCase()),
							}
						})
					} else {
						return item[(data.surface_type.toLowerCase() as "ac") || "cc"]?.conditionList.flatMap((_) => {
							return data.child?.flatMap((child) => {
								return item[(child.surface_type.toLowerCase() as "ac") || "cc"]?.conditionList.map((value) => {
									return {
										value: value.right_value,
										color: this.colorRule(value.grade_id, child.surface_type.toLowerCase()),
									}
								})
							})
						})
					}
				})
			})

			let uniqueZones = Array.from(new Set(zones.map((zone) => JSON.stringify(zone)))).map((str) => JSON.parse(str))
			uniqueZones = ["MPD", "IFI"].includes(this.params.condition_type) ? uniqueZones.reverse() : uniqueZones

			return [
				{
					name: this.params.condition_type,
					data: values,
					zoneAxis: "y",
					zones: uniqueZones,
					animation: false,
				},
			]
		},
		colorRule(gradeId: number, surfaceType: string) {
			const { conditionGrade, params } = this
			const color =
				conditionGrade
					.find((range) => range.ref_condition_range_id === params.ref_condition_range_id)
					?.condition_list.find((condition) => condition.condition_type === params.condition_type)
					?.surface_type[(surfaceType as "ac") || "cc"]?.find((item) => item.grade?.id === gradeId)?.grade?.color ?? ""

			return color
		},
		extractValues() {
			if (this.params.ref_condition_range_id === 3) {
				return this.conditionDataTable.map((item) => item.value) ?? []
			} else {
				return this.conditionDataTable.flatMap((item) => item.child.map((child) => child.value)) ?? []
			}
		},
		toggleType(name: string, eng = true) {
			const typeName = name.toUpperCase()

			if (eng) {
				if (typeName === "IRI" || typeName === "IFI") {
					return typeName + " (m./km.)"
				} else if (typeName === "RUT" || typeName === "MPD") {
					return typeName + " (mm./mm.)"
				} else {
					return typeName
				}
			} else if (typeName === "IRI" || typeName === "IFI") {
				return typeName + " (ม./กม.)"
			} else if (typeName === "RUT" || typeName === "MPD") {
				return typeName + " (มม.)"
			} else {
				return typeName
			}
		},
		setLocation(latLon: string[], imageID: number) {
			this.image.expanded = true

			if (latLon.length > 0) {
				const location = latLon
				this.map.location({
					lon: location[0],
					lat: location[1],
				})
			}

			this.image.imageID = imageID
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
							text: this.toggleType(params.condition_type, false),
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

		histPercent() {
			const series = this.histSeries()
			const data =
				this.details.condition_datas
					?.filter((item) => item.condition_type === this.params.condition_type)
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
		calculateMax() {
			const { conditionDataTable, params } = this
			let minKm = Infinity
			let maxKm = -Infinity

			conditionDataTable.forEach((item) => {
				if (params.ref_condition_range_id === 3) {
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
		getRules(surveyType: string) {
			const { refConditionGrade, params } = this
			const masterRule =
				refConditionGrade.find((condition) => condition?.name === params.condition_type) ?? ({} as ISurveyRule)
			return surveyType === "AC" ? masterRule.ac?.conditionList : masterRule.cc?.conditionList
		},
		histSeries() {
			const { conditionDataTable, params } = this
			const distanceByGrade: { [key: string]: number } = {}

			conditionDataTable.forEach((item) => {
				let distances
				if (params.ref_condition_range_id === 3) {
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
					name: params.condition_type,
					data: fitlerCriteria.map((data) => (data.distance / 1000).toFixed(2)),
				},
			]

			return result || []
		},
		play(speed = 1, buttonName: keyof typeof this.playBtn) {
			const { image } = this
			this.activeButton(buttonName)

			// Error handling for invalid speed
			if (typeof speed !== "number" || speed <= 0) {
				// console.error("Invalid speed. Speed must be a positive number.")
				return
			}

			if (this.conditionDataTable.length > 0) {
				const imageArray = this.conditionDataTable.flatMap((parent) => parent.child?.map((child) => child.id) ?? [])

				// Error handling for empty imageArray
				if (imageArray.length === 0) {
					console.error("imageArray is empty.")
					return
				}

				// Reset to start if imageID is the last ID in imageArray
				if (image.imageID === imageArray[imageArray.length - 1]) {
					image.imageID = imageArray[0] // Assuming the first ID is at index 0
				}

				image.playing = true

				// Clear previous interval if it exists
				if (image.play) {
					clearInterval(image.play)
				}

				image.play = setInterval(() => {
					if (image.imageID < imageArray[imageArray.length - 1]) {
						image.imageID += 1
					} else {
						clearInterval(image.play)
						image.playing = false
						image.lastImageID = null
					}
				}, image.speed / speed)
			}
		},
		stop(buttonName: keyof typeof this.playBtn) {
			if (this.conditionDataTable.length > 0) {
				this.activeButton(buttonName)
				const geom = this.conditionDataTable?.[0]?.child?.[0]?.geom_cl
				const latLon = geom?.split(",")[0].split("(")[1].split(" ")

				if (latLon) {
					this.map.location({
						lon: latLon[0],
						lat: latLon[1],
					})
				}

				this.resetImage()
			}
		},
		resetImage() {
			const firstImageID = this.conditionDataTable[0]?.child[0]?.id
			clearInterval(this.image.play)
			this.image.play = null
			this.image.playing = false
			this.image.lastImageID = null
			this.image.imageID = firstImageID
			this.image.expanded = false
			this.setInitImage()
		},
		pause(buttonName: keyof typeof this.playBtn) {
			this.activeButton(buttonName)
			const image = this.image
			clearInterval(image.play)
			image.playing = false
			image.lastImageID = image.imageID
		},
		activeButton(buttonName: keyof typeof this.playBtn) {
			switch (buttonName) {
				case "play":
					this.playBtn = { play: true, speed2x: false, speed4x: false, stop: false, pause: false }
					break
				case "speed2x":
					this.playBtn = { play: false, speed2x: true, speed4x: false, stop: false, pause: false }
					break
				case "speed4x":
					this.playBtn = { play: false, speed2x: false, speed4x: true, stop: false, pause: false }
					break
				case "stop":
					this.playBtn = { play: false, speed2x: false, speed4x: false, stop: true, pause: false }
					break
				case "pause":
					this.playBtn = { play: false, speed2x: false, speed4x: false, stop: false, pause: true }
					break
				default:
					console.error(`Unknown button: ${buttonName}`)
			}
		},
		setInitImage() {
			const { conditionDataTable, params } = this
			if (conditionDataTable?.length > 0 && params.ref_condition_range_id === 1) {
				const image = conditionDataTable[0]?.child
				this.image.imageID = image ? image[0].id : 0
				this.image.path = image ? image[0].image : ""
			}
		},
		async callBackUpdateData(roadId: number, method: string) {
			const {
				getConditionList,
				getConditionDetails,
				createCriteriaCheckbox,
				createDataTable,
				defaultLocation,
				setDefaultParams,
			} = this

			this.details = {} as IRoadConditionDetails
			this.conditionDataTable = []

			await getConditionList(roadId)

			if (method === "create" || method === "delete") {
				setDefaultParams()
			} else {
				// ดักกรณี update ข้อมูลแล้วมีการเปลี่ยนปี แล้วให้เรียกตัวเดิม
				const findYear = this.conditionList.find((parent) =>
					parent.items.some((child) => child.id_parent === this.params.id_parent)
				)

				if (findYear) {
					this.params.year = findYear.year

					// const findIdParent = findYear.items.find((item) => item.id === res.data.id)
					// if (findIdParent) {
					// 	this.params.id_parent = findIdParent.id_parent
					// }
				}
			}

			await getConditionDetails()
				.then(() => createCriteriaCheckbox())
				.then(() => createDataTable())
				.then(() => {
					defaultLocation()
				})
		},
		onUpdateOwner() {
			const initData = useInitData()
			this.params.ref_condition_range_id =
				initData.conditionGrade().find((item) => item.id === this.params.owner_id)?.ref_condition_range_id ?? null
			this.setRefConditionGrade()
		},
	},
	getters: {
		getYearOptions(state) {
			const { conditionList } = state

			const options = conditionList?.map((condition) => ({ label: `${condition.year + 543}`, value: condition.year }))

			return options.sort((a, b) => b.value - a.value) || []
		},
		getSurveyLaneOptions(state) {
			const { conditionList, params } = state

			const items = conditionList.find((condition) => condition.year === params.year)?.items ?? []
			const options = items.map((item) => ({
				label: `${item.lane_no} (สำรวจ: ${buddhistFormatDate(item.surveyed_date, "dd mmm yy")})`,
				value: item.id_parent,
			}))

			return options || []
		},
		getRangeOptions(state) {
			const { conditionGrade } = state
			const options = conditionGrade.map((grade) => ({ label: grade.owner_name, value: grade.id }))

			return options || []
		},
	getDataType(state) {
		if (!state || !state.conditionGrade || !state.params) {
			return []
		}
		const { conditionGrade, params } = state
		const types = conditionGrade
			.filter((condition) => condition.ref_condition_range_id === params.ref_condition_range_id)
			.flatMap((condition) => condition.condition_list?.flatMap((item) => item.condition_type) ?? [])

		const uniqueTypes = [...new Set(types)]

		const sortOrder = ["IRI", "MPD", "RUT", "IFI"]

		const sortedResult = uniqueTypes.sort((a, b) => {
			return sortOrder.indexOf(a) - sortOrder.indexOf(b) || Infinity
		})

		return sortedResult
	},
		getUpdateUsers(state) {
			const { details } = state
			const users = details.updated_by

			return {
				username: users?.user_name,
				role: users?.department?.name,
				img: users?.profile_picture,
				full_name: users?.full_name,
				date: buddhistFormatDate(details.updated_date, "dd mmm yy เวลา HH:ii น."),
			}
		},
		getSurveyRangeOptions(state) {
			const { conditionGrade, params } = state
			const check = conditionGrade.filter((parentCondition) =>
				parentCondition.condition_list.some((item) => item.condition_type === params.condition_type)
			)

			const options = check.map((item) => ({ label: item.owner_name, value: item.id }))

			return options ?? []
		},
	},
})
