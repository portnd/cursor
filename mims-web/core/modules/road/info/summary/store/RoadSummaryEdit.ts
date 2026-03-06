import {
	IRoad,
	IRoadLane,
	IRequestRoad,
	IRoadSummaryItem,
	IRequestRoadSummary,
	RoadSummaryService,
} from "../infrastructure"
import { IFile } from "~/core/shared/types/File"

interface IState {
	id: number
	loading: boolean
	data: IRoadSummaryItem[]
	directionId: number
	roadKmStart: number
	roadKmEnd: number
	km_total: number | string
	params: IRequestRoad
	road: IRoad
	surface_shoulder_left_id: number[]
	surface_shoulder_right_id: number[]
}

export const useRoadSummaryEditStore = defineStore("road/summary/edit", {
	state: (): IState => ({
		id: 0,
		loading: false,
		directionId: 0,
		roadKmStart: 0,
		roadKmEnd: 0,
		data: [
			{
				id: 0,
				km_start: "",
				km_end: "",
				lane: [
					{
						direction: "",
						lane_no: 0,
						geom_cl: "",
						surface: {
							id: 0,
							name: "",
						},
					},
				],
				material_base: {
					id: 0,
					is_initial: false,
					name: "",
				},
				material_subbase: {
					id: 0,
					is_initial: false,
					name: "",
				},
				material_subgrade: {
					id: 0,
					is_initial: false,
					name: "",
				},
				surface_cross_section_code: 0,
				surface_shoulder_left: {
					id: 0,
					name: "",
				},
				surface_shoulder_right: {
					id: 0,
					name: "",
				},
				thickness_concrete_slab: null,
				thickness_base: null,
				thickness_subbase: null,
				thickness_subgrade: null,
				thickness_surface: null,
				width_shoulder_left: null,
				width_shoulder_right: null,
				width_surface: null,
			},
		] as IRoadSummaryItem[],
		road: {} as IRoad,
		surface_shoulder_left_id: [],
		surface_shoulder_right_id: [],
		km_total: 0,
		params: {} as IRequestRoad,
	}),
	actions: {
		setDataRoad(roadDirectionId: number, kmStart: number, kmEnd: number) {
			this.directionId = roadDirectionId
			this.roadKmStart = kmStart
			this.roadKmEnd = kmEnd
		},
		convertStringToKm(data: any) {
			const a = Number(data.split("+")[0])
			const b = Number(data.split("+")[1])
			const c = a * 1000 + b
			return c
		},
		checkKmTotal() {
			let start = 0
			let end = 0
			let total = 0
			let totalRoad = 0
			for (let index = 0; index < this.data.length; index++) {
				start += this.convertStringToKm(this.data[index].km_start)
				end += this.convertStringToKm(this.data[index].km_end)
			}
			if (this.directionId === 1) {
				totalRoad = this.roadKmEnd - this.roadKmStart
				total = end - start
			} else {
				totalRoad = this.roadKmStart - this.roadKmEnd
				total = start - end
			}
			if (total > totalRoad) {
				return {
					status: false,
					Message:
						"โปรดตรวจสอบช่วง กม. เริ่มต้น - กม. สิ้นสุด เนื่องจากผลรวมระยะทางเกิน " +
						toNumber(totalRoad) +
						" ม. ที่กำหนด",
				}
			} else if (total < totalRoad) {
				return { status: false, Message: "โปรดตรวจสอบช่วง กม. เริ่มต้น - กม. สิ้นสุด เนื่องจากผลรวมระยะทางไม่ถูกต้อง" }
			} else {
				return { status: true, Message: "" }
			}
		},
		async getSurface(id: number) {
			this.id = id
			this.loading = true

			const roadSummaryService = new RoadSummaryService()
			const res = await roadSummaryService.getSurface(this.id)
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				// ตรวจสอบค่าใน stage
				const result = ref<any[]>([])
				// กรณีมีข้อมูล
				if (res.data?.length > 0) {
					for (let index = 0; index < res.data?.length; index++) {
						const item: any = res.data[index]?.items
						item.id = res.data[index].id
						item.isNew = false
						item.km_start = convertMeterToKm(Number(res.data[index].items.km_start))
						item.km_end = convertMeterToKm(Number(res.data[index].items.km_end))
						// Base
						item.material_base = { id: res.data[index].items?.material_base?.id ?? "" }
						// Subbase
						item.material_subbase = { id: res.data[index].items?.material_subbase?.id ?? "" }
						// Subgrade
						item.material_subgrade = { id: res.data[index].items?.material_subgrade?.id ?? "" }
						for (let laneIndex = 0; laneIndex < item.lane_count; laneIndex++) {
							if (item.lane.find((el: any) => el.lane_no === laneIndex + 1) === undefined) {
								item.lane.push({
									lane_no: laneIndex + 1,
									surface: { id: -1, value: "ไม่มีผิวทาง", surface_group: "" },
								})
							}
						}
						item.lane = item.lane.sort((a: any, b: any) => {
							return a.lane_no - b.lane_no
						})
						result.value.push(item)
					}
					this.surface_shoulder_left_id = res.data.map((item) => {
						return item.items.surface_shoulder_left?.id
					})
					this.surface_shoulder_right_id = res.data.map((item) => {
						return item.items.surface_shoulder_right?.id
					})
					// 	// this.leftParams = this.initializeLane(res.data, 1)
					// 	// this.rightParams = this.initializeLane(res.data, 2)
				} else {
					result.value = [
						{
							id: 0,
							km_start: "",
							km_end: "",
							lane_count: 0,
							lane: [
								{
									direction: "",
									lane_no: 0,
									geom_cl: "",
									surface: {
										id: 0,
										name: "",
										color_code: "",
										surface_group: "",
									},
								},
							],
							material_base: {
								id: 0,
								is_initial: false,
								name: "",
								layer_coefficient: 0,
								drainage: 0,
								type: "",
							},
							material_subbase: {
								id: 0,
								is_initial: false,
								name: "",
								layer_coefficient: 0,
								drainage: 0,
								type: "",
							},
							material_subgrade: {
								id: 0,
								is_initial: false,
								name: "",
								layer_coefficient: 0,
								drainage: 0,
								type: "",
							},
							surface_cross_section_code: 0,
							surface_shoulder_left: {
								id: 0,
								name: "",
								color_code: "",
								surface_group: "",
							},
							surface_shoulder_right: {
								id: 0,
								name: "",
								color_code: "",
								surface_group: "",
							},
							thickness_concrete_slab: null,
							thickness_base: null,
							thickness_subbase: null,
							thickness_subgrade: null,
							thickness_surface: null,
							width_shoulder_left: null,
							width_shoulder_right: null,
							width_surface: null,
						},
					]
				}
				// Set ค่า stage
				this.data = result.value
			}
		},
		async updateRoad() {
			this.loading = true
			interface MyObject {
				[key: string]: any
			}

			const params: MyObject = {
				id: this.data.map((item) => (item.isNew ? 0 : item.id)),
				km_start: this.data.map((e) => this.convertStringToKm(e.km_start)),
				km_end: this.data.map((e) => this.convertStringToKm(e.km_end)),
				surface_crosssection_code: this.data.map((e) => {
					return Number(e.surface_cross_section_code)
				}),
				width_surface: this.data.map((e) => {
					return Number(e.width_surface)
				}),
				thickness_surface: this.data.map((e) => {
					return Number(e.thickness_surface)
				}),
				thickness_surface_concrete: this.data.map((e) => {
					return Number(e.thickness_surface_concrete)
				}),
				width_shoulder_left: this.data.map((e) => {
					return Number(e.width_shoulder_left)
				}),
				ref_surface_shoulder_id_left: this.data.map((e) => {
					return Number(e.surface_shoulder_left.id)
				}),
				width_shoulder_right: this.data.map((e) => {
					return Number(e.width_shoulder_right)
				}),
				ref_surface_shoulder_id_right: this.data.map((e) => {
					return Number(e.surface_shoulder_right.id)
				}),
				material_base: this.data.map((e) => {
					if (e.material_base.id === null || e.material_base.id === 0 || !e.material_base.id) {
						return null
					} else {
						return Number(e.material_base.id)
					}
				}),
				material_subbase: this.data.map((e) => {
					if (e.material_subbase.id === null || e.material_subbase.id === 0 || !e.material_subbase.id) {
						return null
					} else {
						return Number(e.material_subbase.id)
					}
				}),
				material_subgrade: this.data.map((e) => {
					if (e.material_subgrade.id === null || e.material_subgrade.id === 0 || !e.material_subgrade.id) {
						return null
					} else {
						return Number(e.material_subgrade.id)
					}
				}),
				thickness_base: this.data.map((e) => {
					if (e.thickness_base === null) {
						return null
					} else {
						return Number(e.thickness_base)
					}
				}),
				thickness_subbase: this.data.map((e) => {
					if (e.thickness_subbase === null) {
						return null
					} else {
						return Number(e.thickness_subbase)
					}
				}),
				thickness_subgrade: this.data.map((e) => {
					if (e.thickness_subgrade === null) {
						return null
					} else {
						return Number(e.thickness_subgrade)
					}
				}),
				thickness_concrete_slab: this.data.map((e) => {
					if (e.thickness_concrete_slab === null) {
						return null
					} else {
						return Number(e.thickness_concrete_slab)
					}
				}),
				lane_surface: this.data.map((e) =>
					e.lane.map((item: IRoadLane) => {
						if (item.surface.id === -1) {
							return 0
						} else {
							return Number(item.surface.id)
						}
					})
				),
				road_id: this.id,
			}
			// let count = 0
			// for (let index = 0; index < params.surface_crosssection_code.length; index++) {
			// 	if (params.surface_crosssection_code[index] === 4) {
			// 		count++
			// 	}
			// }
			// if (count === 0) {
			// 	params.thickness_concrete_slab = params.thickness_concrete_slab.map((_e: any) => {
			// 		return (_e = 0)
			// 	})
			// }

			const roadSummaryService = new RoadSummaryService()
			const res = await roadSummaryService.post(params as IRequestRoadSummary)
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		addLane(id: number) {
			const items: IRoadLane = {
				lane_no: this.data[id].lane.length + 1,
				direction: "",
				geom_cl: "",
				surface: {
					color_code: "",
					surface_group: "",
					id: 0,
					name: "",
				},
			}
			this.data[id].lane.push(items)
		},
		addRow() {
			const itemColumn: IRoadSummaryItem = {
				id: this.data.length + 1,
				isNew: true,
				km_start: "",
				km_end: "",
				lane_count: 0,
				lane: [
					{
						direction: "",
						lane_no: 0,
						geom_cl: "",
						surface: {
							id: 0,
							name: "",
							color_code: "",
							surface_group: "",
						},
					},
				],
				material_base: {
					id: 0,
					is_initial: false,
					name: "",
					layer_coefficient: 0,
					drainage: 0,
					type: "",
				},
				material_subbase: {
					id: 0,
					is_initial: false,
					name: "",
					layer_coefficient: 0,
					drainage: 0,
					type: "",
				},
				material_subgrade: {
					id: 0,
					is_initial: false,
					name: "",
					layer_coefficient: 0,
					drainage: 0,
					type: "",
				},
				surface_cross_section_code: 0,
				surface_shoulder_left: {
					id: 0,
					name: "",
					color_code: "",
					surface_group: "",
				},
				surface_shoulder_right: {
					id: 0,
					name: "",
					color_code: "",
					surface_group: "",
				},
				thickness_concrete_slab: null,
				thickness_base: null,
				thickness_subbase: null,
				thickness_subgrade: null,
				thickness_surface: null,
				width_shoulder_left: null,
				width_shoulder_right: null,
				width_surface: null,
			}
			this.data.push(itemColumn)
		},
		deleteLane(id: number, items: IRoadLane) {
			const idToDelete = items.lane_no
			const indexToDelete = this.data[id].lane.findIndex((element) => element.lane_no === idToDelete)
			if (indexToDelete !== -1) {
				this.data[id].lane.splice(indexToDelete, 1)
			}
		},
		deleteColumn(id: number) {
			const idToDeleteColumn = id
			const indexToDeleteColumn = this.data.findIndex((element) => element.id === idToDeleteColumn)
			if (indexToDeleteColumn !== -1) {
				this.data.splice(indexToDeleteColumn, 1)
			}
		},
		clearFile() {
			this.params.center_lane_shape_file = {} as IFile
			this.params.center_line_shape_file = {} as IFile
		},
		checkKM(value: any) {
			const splitText = value?.split("+")
			if (
				splitText?.length !== 2 ||
				splitText[0]?.length === 0 ||
				typeof Number(splitText[0]) !== "number" ||
				isNaN(Number(splitText[0])) ||
				typeof Number(splitText[1]) !== "number" ||
				isNaN(splitText[1]) ||
				splitText[1]?.length !== 3
			) {
				return true
			} else {
				return false
			}
		},
		updateKm(e: string, label: string) {
			const input = ref()
			if (label === "km_start") {
				input.value = document.querySelector(".input-kmStart") as HTMLElement
			} else {
				input.value = document.querySelector(".input-kmEnd") as HTMLElement
			}

			if (e === "" || this.checkKM(e)) {
				if (input) {
					input.value.style.setProperty("margin-bottom", "-15px", "important")
				}
			} else if (e !== "") {
				if (input) {
					input.value.style.setProperty("margin-bottom", "0", "important")
				}
			}
			if (this.road.road_info) {
				switch (label) {
					case "km_end":
						this.road.road_info.km_end = e !== "" && e ? convertStringToKm(e) : 0
						this.calculateDistance()
						break
					case "km_start":
						this.road.road_info.km_start = e !== "" && e ? convertStringToKm(e) : 0
						this.calculateDistance()

						break
					default:
				}
			}
		},
		updateInfo(e: any, label: string) {
			switch (label) {
				case "remark":
					if (this.road.road_info?.remark) {
						this.params.remark = e
					} else {
						this.params.remark = e
					}
					break
				case "color":
					if (this.road.road_info?.road_color_code) {
						this.road.road_info.road_color_code = e
					}
					break
				case "roadType":
					if (this.road.road_info?.ref_road_type_id) {
						this.road.road_info.ref_road_type_id = e
						// console.log(e)
					}

					break
			}
		},
		converKmRef(kmStart: Ref, kmEnd: Ref) {
			kmStart.value = convertMeterToKm(this.road.road_info?.km_start)
			kmEnd.value = convertMeterToKm(this.road.road_info?.km_end)
		},
		calculateDistance() {
			const result = Math.abs(this.road.road_info?.km_end - this.road.road_info?.km_start) / 1000
			if (!isNaN(result) && result !== 0) {
				this.km_total = toNumber(result, 3) ?? "0.000"
			} else {
				this.km_total = "0.000"
			}
		},
		generateParams() {
			const newParams = {} as IRequestRoad
			newParams.road_code = this.road.road_code
			newParams.name = this.road.road_info.name
			newParams.road_group_id = this.road.road_group_id
			newParams.origin = this.road.origin_to_destination.split(" - ")[0]
			newParams.destination = this.road.origin_to_destination.split(" - ")[1]
			newParams.km_start = Number(this.road.road_info.km_start.toString()?.split("+")?.join(""))
			newParams.km_end = Number(this.road.road_info.km_end.toString()?.split("+")?.join(""))
			newParams.road_color_code = this.road.road_info.road_color_code
			newParams.ref_road_type_id = this.road.road_info.ref_road_type_id
			newParams.register_date = formatDate(new Date(), "yyyy-mm-dd hh:mm:ss")
			newParams.center_line_shape_file = this.params.center_line_shape_file
			newParams.center_line_shape_file_status = this.params.center_line_shape_file?.status
			newParams.center_lane_shape_file = this.params.center_lane_shape_file
			newParams.center_lane_shape_file_status = this.params.center_lane_shape_file?.status
			newParams.remark = this.params.remark === undefined ? this.road.road_info.remark : this.params.remark
			newParams.ramp_id = Number(this.road.road_info.ramp_id)
			newParams.year_construction_completed = this.road.road_info.year_construction_completed ?? 0
			return newParams
		},
		async updateRoads(roadId: number) {
			// Loading
			this.loading = true
			const params = this.generateParams()
			const roadSummaryService = new RoadSummaryService()
			const res = await roadSummaryService.updateRoads(roadId, params)
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
	},
	getters: {},
})
