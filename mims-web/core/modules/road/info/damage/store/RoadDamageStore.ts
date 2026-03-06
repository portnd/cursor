import { RoadDamageService } from "../infrastructure/RoadDamageService"
import { IDatum, IRoadDamageData, IRoadDamageM, IRoadDamageRange, ILane } from "../infrastructure/RoadDamageModel"
import { IOption } from "~/core/shared/types/Option"

export interface IparamsDamage {
	year: number
	parentId: number
	id: number
}

interface IRoadDamageState {
	roadId: number
	loading: boolean
	damageList: IDatum[]
	laneList: ILane[]
	data: IRoadDamageData | null
	map: any
	longdo: any
	params: IparamsDamage
	image: {
		id: number | IOption
		path: string | undefined
	}
	roadColor: string
	roadGeom: string
	isFirstMount: boolean
}

type AccumulatorType = {
	[key in keyof IRoadDamageRange]?: number
}

export const useRoadDamageStore = defineStore("road/damage", {
	state: (): IRoadDamageState => ({
		roadId: 0,
		loading: false,
		damageList: [],
		laneList: [],
		data: null,
		map: null,
		longdo: null,
		params: {
			year: 0,
			parentId: 0,
			id: 0,
		},
		image: {
			id: 1,
			path: "",
		},
		roadColor: "",
		roadGeom: "",
		isFirstMount: true,
	}),
	actions: {
		setMap(map: any) {
			this.map = map
			// @ts-ignore
			this.longdo = window.longdo

			if (this.damageList && this.damageList.length > 0) {
				this.createLine()
			}
		},
		setColor(color: string) {
			this.roadColor = color
		},
		setMainGeom(geom: string) {
			this.roadGeom = geom
		},
		async getLaneList(roadId: number) {
			const service = new RoadDamageService()
			const res = await service.getLaneList(roadId)

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.laneList = res.data
			}
		},
		async getDamageList(id: number) {
			// Loading
			this.loading = true

			const service = new RoadDamageService()
			const res = await service.getRoadDamageList(id)

			// Loading
			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: false })
			} else {
				this.damageList = res.data
			}
		},
		async getRoadDamageDetail(roadId: number) {
			const service = new RoadDamageService()
			const res = await service.getRoadDamageDetails(roadId, this.params.parentId)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: false })
			} else {
				this.data = null

				this.data = res.data
				this.createLine()
			}
			this.setDefaultImage()
		},
		setDefaultparams() {
			if (this.damageList) {
				this.params.year = this.damageList[0].year
				const first = this.damageList[0].items
				this.params.id = first[0].id
				this.params.parentId = this.params.parentId =
					this.damageList
						?.find((item) => item?.year === this?.params?.year)
						?.items.find((item) => item?.id === this.params?.id)?.id_parent || 0
				// this.isFirstMount = false
			}
		},
		setDetails(item: IRoadDamageM) {
			this.image.id = item.id
			this.image.path = item.img_filepath
			this.setLocation(item.the_geom)
		},
		createLine() {
			if (this.map && this.data?.the_geom) {
				this.map.Overlays.clear()

				const geom = this.data?.the_geom

				const line = this.longdo.Util.overlayFromWkt(geom, { lineColor: convertHexToRGBA(this.roadColor, 0.4) })

				this.map.Overlays.add(line[0])

				this.defaultLocation()
				this.map.zoom(13)
				this.createDamagePoints()
			}
		},
		setLocation(geom: string) {
			const latLng = getLatLong(geom)

			this.map.location({
				lon: latLng.lon,
				lat: latLng.lat,
			})
		},
		defaultLocation() {
			if (this.data?.the_geom) {
				const parentGeom = this.data?.the_geom
				const latLng = getLatLong(parentGeom)
				this.map.location({
					lon: latLng.lon,
					lat: latLng.lat,
				})
			}
		},
		setDefaultImage() {
			const roadDamage = this.data?.road_damage
			const roadDamageRange = roadDamage?.road_damage_range

			if (roadDamageRange && roadDamageRange?.length > 0) {
				const roadDamageM = roadDamageRange[0].road_damage_m

				if (roadDamageM?.length > 0) {
					this.image.path = roadDamageM[0].img_filepath || ""
				}
			} else {
				this.image.path = ""
			}
		},
		createOptionsDate() {
			const damageList = this.damageList
			if (!damageList) {
				return []
			}

			const surveyDate = damageList.flatMap((item) => {
				if (item.year === this.params.year && item.items) {
					return item.items.map((child) => ({
						label: `${child?.lane_no} (สำรวจ: ${buddhistFormatDate(child?.surveyed_date, "dd mmm yyyy")})`,
						value: child?.id,
					}))
				}
				return []
			})

			return surveyDate
		},
		createDamagePoints() {
			if (this.getRoadDamageRange.length > 0 && this.map) {
				const damageRange = this.getRoadDamageRange
				const damageM = damageRange.flatMap((item) => item.road_damage_m)
				const geoms = damageM.map((item) => item.the_geom)

				const damagePoints = geoms.flatMap((geom) => {
					return this.longdo.Util.overlayFromWkt(geom, { weight: 10, lineWidth: 7, lineColor: "#CC0000" })
				})

				damagePoints.forEach((point) => {
					this.map.Overlays.add(point)
				})
				this.map.zoom(13)
			}
		},
		async onUpdateYear(roadId: number) {
			const { damageList } = this
			const filterList = damageList.filter((item) => item.year === this.params.year)

			if (filterList.length === 0) {
				return
			}
			this.params.id = filterList[0]?.items[0].id
			this.params.parentId = filterList[0]?.items[0]?.id_parent
			await this.getRoadDamageDetail(roadId)
		},
		async onUpdateIdParent(roadId: number) {
			const { damageList, params } = this
			const findIdParent = damageList
				.find((item) => item.year === params.year)
				?.items?.find((item) => item.id === params.id)

			if (!findIdParent) {
				return
			}
			this.params.parentId = findIdParent?.id_parent
			await this.getRoadDamageDetail(roadId)
		},
	},
	getters: {
		getRoadDamageRange(state) {
			return state.data?.road_damage?.road_damage_range || []
		},
		getParentGeom(state) {
			const data = state.data?.road_damage?.road_damage_range
			if (data?.length) {
				const parentGeom = data.map((item) => item.the_geom)
				return parentGeom
			}

			return []
		},
		getYears(state) {
			if (!state.damageList) {
				return []
			}

			const sortedDate = state.damageList.sort((a, b) => b.year - a.year)

			return (
				sortedDate.map((item) => ({
					label: `${item.year + 543}`,
					value: item.year,
				})) ?? []
			)
		},
		getSumDetails(state) {
			const range = state.data?.road_damage.road_damage_range

			if (!range) {
				return {}
			}

			const result = range.reduce<AccumulatorType>((accumulator, item) => {
				for (const key of Object.keys(item)) {
					if (key.startsWith("ac") || key.startsWith("cc")) {
						const typedKey = key as keyof IRoadDamageRange

						accumulator[typedKey] = (accumulator[typedKey] || 0) + (item[typedKey] as number)
					}
				}
				return accumulator
			}, {})

			const formattedResult = Object.fromEntries(
				Object.entries(result).map(([key, value]) => [key, value === 0 ? 0 : parseFloat(value.toFixed(2))])
			)

			return formattedResult
		},
		createLaneOptions(state) {
			const lanes = state.laneList
			const options = lanes.map((item: ILane) => {
				return { label: item.lane_no.toString(), value: item.lane_no }
			})

			return options
		},
	},
})
