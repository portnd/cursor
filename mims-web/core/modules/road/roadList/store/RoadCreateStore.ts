import { ICreateRoad, IRequestRoadInit, RoadListService } from "../infrastructure"

interface IInitParams {
	roadId: string
	level: string
	roadTypeId: string
}

interface IState {
	loading: boolean
	roadInitParams: IInitParams
	roadInit: IRequestRoadInit
	params: ICreateRoad
	roadId: number
}

export const useRoadCreateStore = defineStore("road-create", {
	state: (): IState => ({
		loading: false,
		roadInitParams: {} as IInitParams,
		roadInit: {} as IRequestRoadInit,
		params: {} as ICreateRoad,
		roadId: 0,
	}),
	actions: {
		async getRoadInit(id: string, level: string, roadType: string) {
			this.loading = true
			const service = new RoadListService()
			const res = await service.getRoadinit(id, level, roadType)
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.roadInit = res.data
			}
			this.loading = false
		},
		async createRoad() {
			this.loading = true

			const params = this.generateParams()
			const service = new RoadListService()
			const res = await service.createRoad(params)
			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.roadId = res.data.id
				return res
			}
		},
		updateDistance(distance: Ref) {
			const kmStart = convertStringToKm(this.params?.km_start?.toString() ?? "")
			const kmEnd = convertStringToKm(this.params?.km_end?.toString() ?? "")
			const result = Math.abs(kmEnd - kmStart) / 1000
			if (!isNaN(result) && result !== 0) {
				distance.value = toNumber(result, 3) ?? "0.000"
			} else {
				distance.value = "0.000"
			}
		},
		setRoadType(roadType: number) {
			this.roadInitParams.roadTypeId = roadType?.toString()
			this.getRoadInit(this.roadInitParams.roadId, this.roadInitParams.level, this.roadInitParams.roadTypeId)
		},
		clearFile() {
			this.params.center_lane_shape_file = {} as File
			this.params.center_line_shape_file = {} as File
		},
		generateParams() {
			const newParams = {} as ICreateRoad
			newParams.name = this.params.name
			newParams.road_section_id = Number(this.roadInitParams.roadId)
			newParams.road_id = Number(this.roadInitParams.roadId)
			newParams.road_level = Number(this.roadInitParams.level)
			newParams.ramp_id = Number(this.params.ramp_id)
			newParams.km_start = Number(this.params.km_start.toString()?.split("+")?.join(""))
			newParams.km_end = Number(this.params.km_end.toString()?.split("+")?.join(""))
			newParams.road_color_code = this.params.road_color_code
			newParams.ref_road_type_id = Number(this.roadInitParams.roadTypeId)
			newParams.register_date = formatDate(new Date(), "yyyy-mm-dd hh:mm:ss")
			newParams.center_line_shape_file = this.params.center_line_shape_file
			newParams.center_lane_shape_file = this.params.center_lane_shape_file
			newParams.remark = this.params.remark ?? ""
			newParams.year_construction_completed = this.params.year_construction_completed ?? 0

			return newParams
		},
	},
	getters: {},
})
