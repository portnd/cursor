import { defineStore } from "pinia"
import {
	IRoad,
	ILaneList,
	IRequestRoad,
	IRoadSummaryItem,
	RoadSummaryService,
	IItemConditionCompareAverage,
	IRoadSummaryUpdateBy,
	IMaintenanceProjectsData,
	ISurfaceIcon,
	IRoadLane,
	ITrafficModel,
	ITrafficDetail,
	IRequestTraffic,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IToggle {
	year: number
	parentID: number
	aadtId: number
}

interface IParams {
	surveyedDate: Date
	veh1: number | null
	veh2: number | null
	veh3: number | null
}

interface state {
	loading: boolean
	map: any
	roadId: number
	road: IRoad
	km_total: number
	params: IRequestRoad
	roadSurface: IRoadSummaryItem[]
	surfaceIcon: ISurfaceIcon[]
	// traffic
	trafficRevision: ITrafficModel[]
	trafficDetail: ITrafficDetail
	toggle: IToggle
	trafficParams: IParams
	surveyedDate: Date
	veh1: number
	veh2: number
	veh3: number
	//
	lane_id: number
	status_code: string
	status: string
	reject_reason: string
	conditionLaneId: number
	laneList: ILaneList[]
	conditionCompareAverage: IItemConditionCompareAverage[]
	dataIRIChart: Array<number>
	dataCategories: Array<string>
	dataGNChart: Array<number>
	surfaceSectionCode: number
	dataInPicture: IRoadSummaryItem
	update_by: IRoadSummaryUpdateBy
	update_date: string
	yearList: number[]
	yearParams: number
	maintenanceProjects: IMaintenanceProjectsData[]
}

export const useRoadSummaryStore = defineStore("road-summary", {
	state: (): state => ({
		loading: false,
		map: null,
		roadId: 0,
		road: { ref_depot: { id: 0 } } as IRoad,
		km_total: 0,
		params: {} as IRequestRoad,
		roadSurface: [],
		surfaceIcon: [],
		// traffic
		trafficRevision: [],
		toggle: {
			year: 0,
			parentID: 0,
			aadtId: 0,
		},
		trafficDetail: {} as ITrafficDetail,
		trafficParams: {
			surveyedDate: new Date(),
			veh1: null,
			veh2: null,
			veh3: null,
		},
		surveyedDate: new Date(),
		veh1: 0,
		veh2: 0,
		veh3: 0,
		//
		lane_id: 1,
		status_code: "",
		status: "",
		reject_reason: "",
		conditionLaneId: 1, // ค่าเริ่มต้น เลนที่ 1
		laneList: [],
		conditionCompareAverage: [],
		dataIRIChart: [],
		dataCategories: [],
		dataGNChart: [],
		surfaceSectionCode: 0,
		dataInPicture: {} as IRoadSummaryItem,
		update_by: {} as IRoadSummaryUpdateBy,
		update_date: "",
		yearList: [],
		yearParams: 0,
		maintenanceProjects: [],
	}),
	actions: {
		async getLaneList() {
			// Loading
			this.loading = true

			const roadSummaryService = new RoadSummaryService()
			const res = await roadSummaryService.getLaneList(this.roadId)

			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.laneList = res.data
			}
		},
		async getConditionCompareAverge() {
			// Loading
			this.loading = true
			const roadSummaryService = new RoadSummaryService()
			const res = await roadSummaryService.getConditionCompareAverage(this.roadId, this.conditionLaneId)

			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.conditionCompareAverage = []
				this.dataIRIChart = []
				this.dataCategories = []
				this.dataGNChart = []

				if (Object.keys(res.data).length !== 0) {
					this.conditionCompareAverage = res.data.items
					if (this.conditionCompareAverage.length !== 0) {
						this.conditionCompareAverage.forEach((e: IItemConditionCompareAverage) => {
							const year = e.year + 543
							this.dataCategories.push(year.toString())
							this.dataIRIChart.push(e.iri)
							this.dataGNChart.push(e.gn)
						})
					}
				}
			}
		},
		// begin-surface
		async getRoadSurface(roadId: number) {
			this.roadId = roadId
			// Loading
			this.loading = true
			const roadSummaryService = new RoadSummaryService()
			const res = await roadSummaryService.getSurface(this.roadId)

			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				// สถานะ
				this.status_code = ""
				this.status = ""
				this.reject_reason = ""

				// ข้อมูลผิวทาง
				this.roadSurface = []
				if (res.data.length > 0) {
					for (let index = 0; index < res.data.length; index++) {
						// สถานะ
						this.status_code = res.data[index].status_code
						this.status = res.data[index].status
						this.reject_reason = res.data[index].reject_reason

						// ข้อมูลผิวทาง
						const item: IRoadSummaryItem = res.data[index].items
						item.id = index + 1
						item.km_start = convertMeterToKm(Number(res.data[index].items.km_start))
						item.km_end = convertMeterToKm(Number(res.data[index].items.km_end))
						this.$patch((state) => {
							state.roadSurface.push(res.data[index].items)
						})

						// ข้อมูล update_by, update_date
						this.update_by = res.data[index].update_by
						this.update_date = res.data[index].update_date
						for (let laneIndex = 0; laneIndex < res.data[index]?.items?.lane_count; laneIndex++) {
							if (item.lane.find((el) => el.lane_no === laneIndex + 1) === undefined) {
								item.lane.push({
									lane_no: laneIndex + 1,
									surface: { id: 0, name: "ไม่มีผิวทาง", surface_group: "", color_code: "" },
								} as IRoadLane)
							}
						}
						item.lane = item.lane.sort((a, b) => {
							return a.lane_no - b.lane_no
						})
					}
					this.lane_id = this.roadSurface[0].lane[0].lane_no
					this.roadSurface = this.roadSurface.filter((e: IRoadSummaryItem) => e.id !== 0)
					this.surfaceSectionCode = this.roadSurface[0].surface_cross_section_code
					this.dataInPicture = this.roadSurface[0]
				}
			}
		},
		async getSurfaceIcon() {
			// Loading
			this.loading = true
			const roadSummaryService = new RoadSummaryService()
			const res = await roadSummaryService.getSurfaceIcon(this.roadId)

			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.surfaceIcon = res.data
			}
		},
		setSurfaceLocation(isLane: boolean, item: IRoadSummaryItem, index = 0) {
			const road = ref<string>()
			if (isLane) {
				const indexLane = index + 1

				if (item.lane.find((el) => el.lane_no === indexLane)) {
					road.value = item.lane.filter((el) => el.lane_no === indexLane)[0]?.geom_cl
				}
			} else {
				const data = item.lane
				road.value = data.flatMap((item) => item.geom_cl).filter((el) => el !== undefined)[0]
			}
			// LINESTRING
			let coordinator: Array<string> = []
			if (road.value) {
				if (road.value.includes("LINESTRING")) {
					coordinator = road.value.split(",")[0]?.split("(")[1]?.split(" ")
				} else if (road.value.includes("POINT")) {
					coordinator = road.value.replaceAll(/[()]/g, "").split("POINT")[1]?.split(" ")
				} else {
					// กรณี ไม่ถูก format
					return
				}
			}

			this.map.location({
				lon: coordinator[0],
				lat: coordinator[1],
			})
			this.map.zoom(18)
		},
		// end-surface
		setMap(map: Object) {
			this.$patch((state) => {
				state.map = map
			})
			setTimeout(() => {
				this.createLine()
			}, 1500)
		},
		setLocation(item: IRoadSummaryItem) {
			const geom = item.lane.flatMap((item) => item.geom_cl).filter((geom) => geom !== undefined)[0]

			// LINESTRING
			let coordinator: Array<string> = []
			if (geom.includes("LINESTRING")) {
				coordinator = geom.split(",")[0].split("(")[1].split(" ")
			} else if (geom.includes("POINT")) {
				coordinator = geom.replaceAll(/[()]/g, "").split("POINT")[1].split(" ")
			} else {
				// กรณี ไม่ถูก format
				return
			}

			this.map.location({
				lon: coordinator[0],
				lat: coordinator[1],
			})
			this.map.zoom(18)
		},
		createLine() {
			if (!this.map) {
				return
			}
			const linesGeom = this.roadSurface.flatMap((parent) => {
				return parent.lane?.map((line) => ({ ...line, geom_cl: line.geom_cl, color: line.surface.color_code }))
			})

			if (linesGeom.length > 0) {
				// @ts-ignore
				const longdo = window.longdo
				const lines = linesGeom.flatMap((line) => {
					return longdo.Util.overlayFromWkt(line?.geom_cl ?? "", {
						lineColor: line?.color,
						detail: this.generatePopupDetails(line),
					})
				})
				lines.forEach((line: any) => {
					this.map.Overlays.add(line)
				})
				const latLng = getLatLong(this.road.road_info?.the_geom)
				this.map.location({
					lon: latLng.lon,
					lat: latLng.lat,
				})
			} else if (this.road?.road_info?.the_geom) {
				// ยังไม่มี roadSurface (ยังไม่เปิดแท็บผิวทาง) ใช้ geometry หลักของสายทาง
				// @ts-ignore
				const longdo = window.longdo
				const polyline = longdo.Util.overlayFromWkt(this.road.road_info.the_geom, {
					lineColor: this.road.road_info?.road_color_code,
				})
				if (polyline?.[0]) {
					this.map.Overlays.add(polyline[0])
				}
				const latLng = getLatLong(this.road.road_info.the_geom)
				this.map.location({
					lon: latLng.lon,
					lat: latLng.lat,
				})
			}
			this.map.zoom(18)
		},
		generatePopupDetails(line: any) {
			return `
				</div>
				<div class="row">
					<div class="col-6">ช่องจราจร:</div>
					<div class="col" style="
					color: #5E6278;
					font-size: 12px;
					font-weight: 400;
				">${line?.lane_no}</div>
				</div>

				<div class="row">
					<div class="col-6">ชนิดผิว:</div>
					<div class="col" style="
					color: #5E6278;
					font-size: 12px;
					font-weight: 400;
				">${line?.surface?.name} (${line?.surface?.surface_group})</div>
				</div>
			`
		},
		async getMaintenanceYears(roadId: number) {
			this.loading = true

			const service = new RoadSummaryService()
			const res = await service.getMaintenanceYears(roadId)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
				this.yearList = []
				this.maintenanceProjects = []
			} else {
				this.yearList = Array.isArray(res.data) ? res.data : []

				if (this.yearList.length > 0) {
					const sortYear = [...this.yearList].sort((a, b) => b - a)
					this.yearParams = sortYear[0]
					await this.getMaintenanceProjects(roadId, sortYear[0])
				} else {
					this.yearParams = 0
					this.maintenanceProjects = []
				}
			}
			this.loading = false
		},
		async getMaintenanceProjects(roadId: number, year: number) {
			// this.loading = true

			const service = new RoadSummaryService()
			const res = await service.getMaintenanceProjects(roadId, year)

			// this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.maintenanceProjects = res.data

				this.createMaintenanceHistoryLine()
			}
		},
		toLocation(geom: any) {
			if (geom.type) {
				const latLon = getLatLongObject(geom)
				this.map.location({ lon: latLon.lon, lat: latLon.lat })
				this.map.zoom(18)
			}

			// setTimeout(() => {
			// @ts-ignore
			// const longdo = window.longdo
			// const popup = new longdo.Popup(latLon, {
			// 	detail: this.generatePopupDetails(item),
			// })
			// this.map.Overlays.add(popup)
			// }, 100)
		},
		setDateColor(guranteeDate: Date) {
			const currentDate = new Date()
			const dateToCompare = new Date(guranteeDate)

			const currentDateDay = currentDate.getDate()
			const currentMonth = currentDate.getMonth()
			const currentYear = currentDate.getFullYear()

			const dateToCompareDay = dateToCompare.getDate()
			const compareMonth = dateToCompare.getMonth()
			const compareYear = dateToCompare.getFullYear()

			const diffYear = compareYear - currentYear
			const diffMonth = compareMonth + diffYear * 12 - currentMonth
			const diffDay = dateToCompareDay - currentDateDay

			let diff = diffMonth
			if (diffMonth === 0) {
				diff = diffDay / 30
			}

			if (diff < 0) {
				return "text-gray-600"
			} else if (diff <= 3) {
				return "text-danger"
			} else if (diff <= 6) {
				return "text-warning"
			} else {
				return "text-success"
			}
		},
		createMaintenanceHistoryLine() {
			if (!this.map) {
				return
			}

			this.map.Overlays.clear()

			// @ts-ignore
			const longdo = window.longdo
			const maintenanceProjects = this.maintenanceProjects
			if (Object.keys(maintenanceProjects)?.length > 0) {
				const geoms = maintenanceProjects.flatMap((item) =>
					[...item.roads, ...item.road_histories].flatMap((child) => {
						return {
							geom: child.the_geom,
							color: child.color,
							road: child,
						}
					})
				)

				if (geoms.length) {
					geoms.forEach((item) => {
						if (item.geom && item.geom.type) {
							const strLine = longdo.Util.overlayFromGeoJSON(item.geom, {
								lineColor:
									item.color === "#000000" ? convertHexToRGBA(item.color, 0.5) : convertHexToRGBA(item.color, 1),
								// detail: this.generatePopupDetails(item.road),
							})
							console.log("%c line color", `background: ${convertHexToRGBA(item.color, 0.5)}; color: ${item.color}`)
							this.map.Overlays.add(strLine[0])
						}
					})
				}
			}
		},
		async getRoadDetail(roadId: number) {
			this.roadId = roadId
			// Loading
			this.loading = true
			const roadSummaryService = new RoadSummaryService()
			const res = await roadSummaryService.getRoads(this.roadId)

			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.road = res.data
				this.params.center_line_shape_filepath = this.road.road_info.center_line_shape_file_path
				this.params.center_lane_shape_filepath = this.road.road_info.center_lane_shape_file_path
				// เมื่อมีข้อมูลสายทางแล้ว ถ้า map พร้อมแล้ว ให้วาดเส้นบนแผนที่ทันที (แก้ refresh แล้วแผนที่เป็นสีฟ้า)
				if (this.map) {
					this.createLine()
				}
			}
		},
		// traffic
		createTrafficLine() {
			// @ts-ignore
			const longdo = window.longdo
			const polyline = longdo.Util.overlayFromWkt(this.road.road_info?.the_geom, {
				lineColor: this.road.road_info?.road_color_code,
			})
			this.map?.Overlays.add(polyline[0])
			this.map?.location({
				lon: getLatLong(this.road.road_info?.the_geom).lon,
				lat: getLatLong(this.road.road_info?.the_geom).lat,
			})
			this.map.zoom(18)
		},
		async getTrafficRevision(roadId: number) {
			this.loading = true
			try {
				const roadSummaryService = new RoadSummaryService()
				const res = await roadSummaryService.getTrafficRevision(roadId)
				if (!res.status) {
					useHandlerError(res.code, res.error, { showAlert: true })
					this.trafficRevision = []
				} else {
					this.trafficRevision = Array.isArray(res.data) ? res.data : []
				}
			} finally {
				this.loading = false
			}
		},
	setDefaultOptions() {
		if (!this.trafficRevision || this.trafficRevision.length === 0) {
			this.toggle.year = 0
			this.toggle.aadtId = 0
			return
		}
		const sorted = [...this.trafficRevision].sort((a, b) => b.year - a.year)
		this.toggle.year = sorted[0]?.year ?? 0
		this.toggle.aadtId = sorted[0]?.items?.[0]?.id ?? 0
	},
		setUpdatedOptions(_id: number, _idParent: number) {
			if (!this.trafficRevision || this.trafficRevision.length === 0) {
				this.toggle.year = 0
				this.toggle.aadtId = 0
				return
			}
			const sorted = [...this.trafficRevision].sort((a, b) => b.year - a.year)
			this.toggle.year = sorted[0]?.year ?? 0
			this.toggle.aadtId = sorted[0]?.items?.[0]?.id ?? 0
		},
		async getTrafficDetail(parentID: number) {
			if (!parentID) {
				return
			}
			this.loading = true
			try {
				const roadSummaryService = new RoadSummaryService()
				const res = await roadSummaryService.getTrafficDetail(this.roadId, parentID)
				if (!res.status) {
					useHandlerError(res.code, res.error, { showAlert: true })
				} else {
					this.trafficDetail = res.data
				}
			} finally {
				this.loading = false
			}
		},
		async createTraffic() {
			const params: IRequestTraffic = {
				surveyed_date: this.trafficParams.surveyedDate ? formatDate(this.trafficParams.surveyedDate) : "",
				veh1: Number(this.trafficParams.veh1),
				veh2: Number(this.trafficParams.veh2),
				veh3: Number(this.trafficParams.veh3),
			}
			// Loading
			this.loading = true
			const roadSummaryService = new RoadSummaryService()
			const res = await roadSummaryService.createTraffic(this.roadId, params)

			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		setTrafficParams() {
			if (this.trafficDetail) {
				this.surveyedDate = this.trafficDetail?.surveyed_date
				this.veh1 = this.trafficDetail?.veh1
				this.veh2 = this.trafficDetail?.veh2
				this.veh3 = this.trafficDetail?.veh3
			}
		},
		async updateTraffic(aadtId: number, idParent: number) {
			const params: IRequestTraffic = {
				surveyed_date: this.surveyedDate ? formatDate(this.surveyedDate) : "",
				veh1: Number(this.veh1),
				veh2: Number(this.veh2),
				veh3: Number(this.veh3),
			}
			// Loading
			this.loading = true
			const roadSummaryService = new RoadSummaryService()
			const res = await roadSummaryService.updateTraffic(this.roadId, idParent, aadtId, params)

			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		async deleteTraffic(parentID: number) {
			// Loading
			this.loading = true
			const roadSummaryService = new RoadSummaryService()
			const res = await roadSummaryService.deleteTraffic(this.roadId, parentID)

			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
	},
	getters: {
	getYearListOptions(state) {
		if (!state?.trafficRevision?.length) {
			return []
		}
		return [...state.trafficRevision]
			.sort((a, b) => b.year - a.year)
			.map((item) => ({ label: `${item.year + 543}`, value: item.year }))
	},
	getYearItemsOptions(state) {
		if (!state) {
			return []
		}
		const yearList = state.trafficRevision
		if (!yearList || yearList.length === 0) {
			return []
		}

		const options = yearList.flatMap((item) => {
			if (item.year === state.toggle.year) {
				return (item.items ?? []).map((traffic) => ({
					label: `สำรวจ: ${buddhistFormatDate(traffic.surveyed_date, "dd mmm yyyy")}`,
					value: traffic?.id,
				}))
			}

			return []
		})
		return options
	},
	getUpdateBy(state) {
		if (!state || !state.trafficDetail?.updated_by) {
			return { img: "", name: "", department: "", date: "" }
		}

		const updatedData = state.trafficDetail.updated_by
		const imgProfile = updatedData?.profile_img_path
		const name = `${updatedData?.firstname ?? ""} ${updatedData?.lastname ?? ""}`.trim()
		const department = updatedData?.department?.name
		const dateTimes = state.trafficDetail?.updated_date
			? `เมื่อวันที่ ${buddhistFormatDate(
					state.trafficDetail.updated_date,
					"dd mmm yyyy เวลา HH:ii น."
			  )} น.`
			: ""

		return { img: imgProfile ?? "", name: name ?? "", department: department ?? "", date: dateTimes ?? "" }
	},
	getYearsOptions(state) {
		if (!state || !state.yearList || state.yearList.length === 0) {
			return []
		}

		const options: IOption[] = state.yearList.map((year) => {
			return { label: `${year + 543}`, value: year }
		})

		return options || []
	},
	getSumMaintenanceKm(state) {
		if (!state || !state.maintenanceProjects || state.maintenanceProjects.length === 0) {
			return 0
		}

		const data = state.maintenanceProjects
			const maintenanceRoads = data.flatMap((item) => item.roads)
			const sum = maintenanceRoads
				?.map((item) => {
					const kmEnd = Number(item.km_end)
					const kmStart = Number(item.km_start)

					const result = Math.abs((kmEnd - kmStart) / 1000)

					return result
				})
				.reduce((sum, curr) => sum + curr, 0)

			return sum ?? 0
		},
		getSumMaintenanceHistKm(state) {
			if (state.maintenanceProjects.length === 0) {
				return 0
			}

			const data = state.maintenanceProjects
			const maintenanceRoads = data.flatMap((item) => item.road_histories)
			const sum = maintenanceRoads
				?.map((item) => {
					const kmEnd = Number(item.km_end)
					const kmStart = Number(item.km_start)

					const result = Math.abs((kmEnd - kmStart) / 1000)

					return result
				})
				.reduce((sum, curr) => sum + curr, 0)

			return sum ?? 0
		},
	},
})
