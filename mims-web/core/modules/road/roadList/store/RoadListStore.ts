import { LocationQuery } from "vue-router"
import { IRequestRoadList, IRoadList, IRoadListRoads, RoadListService } from "../infrastructure"

type KeyInterest = "iri_1000" | "iri_100" | "rut_100" | "ifi_100" | "g7_100"

interface Location {
	lon: number
	lat: number
}

interface geomDetail {
	geom: string
	color: string
	name: string
	section: string
	road: string
	originToDestination: string
	code: string
	km: string
	distance: string
	surface: any[]
	level: number
}

interface IStateParams {
	keyword: string | null
	road_group_id: string[] | null
	road_section_id: number[] | null
	road_id: string[] | null
	km_start: string | null
	km_end: string | null
	depot_code: number[] | null
	ref_surface_id: number[] | null
	is_iri_1000: boolean | null
	is_iri_100: boolean | null
	is_rut_100: boolean | null
	is_ifi_100: boolean | null
	is_g7_100: boolean | null
}

interface IConditionTemp {
	is_iri_1000: boolean | null
	is_iri_100: boolean | null
	is_rut_100: boolean | null
	is_ifi_100: boolean | null
	is_g7_100: boolean | null
}
interface IConditionParentTemp {
	iri_1000: boolean
	iri_100: boolean
	rut_100: boolean
	ifi_100: boolean
	g7_100: boolean
}

// interface GeomWkt {
// 	geom: string
// 	color: string
// }

interface IState {
	loading: boolean
	roads: IRoadList[]
	map?: any
	params: IStateParams
	location: Location
	condition_value: string
	condition_temp: IConditionTemp
	condition_prev: IConditionTemp
	condition_all: boolean
	condition_parent: IConditionParentTemp
	condition_parent_temp: IConditionParentTemp
}

export const useRoadListStore = defineStore("road-list", {
	state: (): IState => ({
		loading: false,
		roads: [],
		map: null,
		params: {
			keyword: null,
			road_group_id: null,
			road_section_id: null,
			road_id: null,
			km_start: null,
			km_end: null,
			depot_code: null,
			ref_surface_id: null,
			is_iri_1000: null,
			is_iri_100: null,
			is_rut_100: null,
			is_ifi_100: null,
			is_g7_100: null,
		},
		location: {
			lon: 100.737091,
			lat: 13.730756,
		},
		condition_value: "",
		condition_temp: {
			is_iri_1000: null,
			is_iri_100: null,
			is_rut_100: null,
			is_ifi_100: null,
			is_g7_100: null,
		},
		condition_prev: {
			is_iri_1000: null,
			is_iri_100: null,
			is_rut_100: null,
			is_ifi_100: null,
			is_g7_100: null,
		},
		condition_parent: {
			iri_1000: false,
			iri_100: false,
			rut_100: false,
			ifi_100: false,
			g7_100: false,
		},
		condition_parent_temp: {
			iri_1000: false,
			iri_100: false,
			rut_100: false,
			ifi_100: false,
			g7_100: false,
		},
		condition_all: false,
	}),
	actions: {
		async getData() {
			this.loading = true

			const params: Partial<IRequestRoadList> = this.checkParams()

			const service = new RoadListService()
			const res = await service.getRoads(params)

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.roads = res.data
			}
			this.loading = false
		},
		setQuriesParams(queries: LocationQuery) {
			Object.keys(queries).forEach((key) => {
				switch (key) {
					case "road_group_id":
						const roadGroupId = String(queries.road_group_id)
						this.params.road_group_id = [roadGroupId]
						break
					case "road_id":
						queries.road_id = queries.road_id as string
						const roadSectionId = queries.road_id?.split(",")

						if (!this.params.road_id) {
							this.params.road_id = []
						}

						roadSectionId.forEach((id) => {
							this.params.road_id?.push(id)
						})

						// ดักเคส duplicate id
						this.params.road_id = [...new Set(this.params.road_id?.filter((id) => id))]
						break
					case "ref_surface_id":
						queries.ref_surface_id = queries.ref_surface_id as string
						const surfaceId = queries.ref_surface_id?.split(",")

						if (!this.params.ref_surface_id) {
							this.params.ref_surface_id = []
						}

						surfaceId.forEach((id) => {
							this.params.ref_surface_id?.push(Number(id))
						})

						// ดักเคส duplicate id
						this.params.ref_surface_id = [...new Set(this.params.ref_surface_id?.filter((id) => id))]
						break
					default:
						break
				}
			})
		},
		checkParams() {
			type KeyInterest = "iri_1000" | "iri_100" | "rut_100" | "ifi_100" | "g7_100"
			type ParamsKeyInterest = "is_iri_1000" | "is_iri_100" | "is_rut_100" | "is_ifi_100" | "is_g7_100"

			const params: Partial<IRequestRoadList> = {
				keyword: this.params.keyword ? this.params.keyword : undefined,
				road_group_id: this.params.road_group_id?.join(",") ? this.params.road_group_id.join(",") : undefined,
				road_section_id: this.params.road_section_id?.join(",") ? this.params.road_section_id.join(",") : undefined,
				road_id: this.params.road_id?.join(",") ? this.params.road_id.join(",") : undefined,
				km_start:
					convertStringToKm(this.params.km_start ?? "").toString() === "NaN"
						? undefined
						: convertStringToKm(this.params.km_start ?? "").toString(),
				km_end:
					convertStringToKm(this.params.km_end ?? "").toString() === "NaN"
						? undefined
						: convertStringToKm(this.params.km_end ?? "").toString(),
				depot_code: this.params.depot_code?.join(",") ? this.params.depot_code.join(",") : undefined,
				ref_surface_id: this.params.ref_surface_id?.join(",") ? this.params.ref_surface_id.join(",") : undefined,
			}

			const keyMapping: Record<KeyInterest, ParamsKeyInterest> = {
				iri_1000: "is_iri_1000",
				iri_100: "is_iri_100",
				rut_100: "is_rut_100",
				ifi_100: "is_ifi_100",
				g7_100: "is_g7_100",
			}

			for (const [key, paramKey] of Object.entries(keyMapping)) {
				const interestKey = key as KeyInterest
				const paramInterestKey = paramKey as ParamsKeyInterest

				if (this.condition_parent[interestKey] && this.params[paramInterestKey] !== null) {
					params[paramInterestKey] = this.params[paramInterestKey] ?? undefined
				} else {
					delete params[paramInterestKey]
				}
			}

			return params
		},
		conditionCheckAll(checked: boolean) {
			const keysInterest: KeyInterest[] = ["iri_1000", "iri_100", "rut_100", "ifi_100", "g7_100"]

			keysInterest.forEach((key) => {
				this.condition_parent[key] = checked

				const tempKey = `is_${key}` as keyof typeof this.condition_temp
				if (checked) {
					if (this.condition_temp[tempKey] === null) {
						this.condition_temp[tempKey] = true
					}
				} else {
					this.condition_temp[tempKey] = null
				}
			})

			this.countCheckBox()
		},
		onUpdateCheckBox(key: KeyInterest) {
			if (!this.condition_parent[key]) {
				this.condition_temp[`is_${key}`] = null
			} else {
				this.condition_temp[`is_${key}`] = true
			}

			const checkedAll = Object.values(this.condition_parent).every((value) => value)
			this.condition_all = checkedAll

			this.countCheckBox()
		},
		countCheckBox() {
			type KeyInterest = "iri_1000" | "iri_100" | "rut_100" | "ifi_100" | "g7_100"
			const keysInterest: KeyInterest[] = ["iri_1000", "iri_100", "rut_100", "ifi_100", "g7_100"]

			const count = keysInterest.reduce((acc, key) => acc + (this.condition_parent[key] ? 1 : 0), 0)

			this.condition_value = count === 0 || count === 4 ? "" : `เลือก ${count} รายการ`
		},
		onSubmit() {
			const conditionTemp = JSON.parse(JSON.stringify(this.condition_temp)) as IConditionTemp
			const conditionParentTemp = JSON.parse(JSON.stringify(this.condition_parent)) as IConditionParentTemp

			this.condition_prev = conditionTemp
			this.condition_parent_temp = conditionParentTemp

			this.params.is_iri_1000 = conditionTemp.is_iri_1000
			this.params.is_iri_100 = conditionTemp.is_iri_100
			this.params.is_rut_100 = conditionTemp.is_rut_100
			this.params.is_ifi_100 = conditionTemp.is_ifi_100
			this.params.is_g7_100 = conditionTemp.is_g7_100

			// this.countCheckBox()
		},
		onCancel() {
			const conditionPrev = JSON.parse(JSON.stringify(this.condition_prev)) as IConditionTemp
			const conditionParentPrev = JSON.parse(JSON.stringify(this.condition_parent_temp)) as IConditionParentTemp

			this.condition_parent = conditionParentPrev
			this.condition_temp = conditionPrev

			this.condition_parent.iri_1000 = conditionParentPrev.iri_1000
			this.condition_parent.iri_100 = conditionParentPrev.iri_100
			this.condition_parent.rut_100 = conditionParentPrev.rut_100
			this.condition_parent.ifi_100 = conditionParentPrev.ifi_100
			this.condition_parent.g7_100 = conditionParentPrev.g7_100

			this.condition_all = Object.values(conditionPrev).every((value) => value)
			this.countCheckBox()
		},
		setMap(map: Object) {
			this.map = map

			// initial Default Location
			// if (this.roads) {
			this.defaultLocation()
			// }

			// Map Functions
			this.createLine()
		},

		generateSurfaceDetail(surface: any[]) {
			const detailHtml = ref("")
			surface.forEach((element: any) => {
				detailHtml.value += `
							<span style="background-color: ${element.color_code}" class="badge badge-primary text-white px-3 py-2 rounded-pill fw-normal">
								${element.name}
							</span>
				`
			})
			return detailHtml.value
		},
		checkRoadLevelDetail(level: number, originToDestination: string) {
			const detailHTML = ref("")
			if (level === 1) {
				detailHTML.value += `<div class="row">
				<div class="col-4">จาก - ถึง</div>
				<div class="col colon">:</div>
				<div class="col">${originToDestination}</div>
			</div>`
			} else {
				detailHTML.value += `<div style="display:none"></div>`
			}
			return detailHTML.value
		},
		createLine() {
			if (this.map) {
				// reset Overlays in the map
				this.map.Overlays.clear()
				// @ts-ignore
				const longdo = window.longdo
				const lines = this.getRoadGeoms.map((item) => {
					return longdo.Util.overlayFromWkt(item.geom, {
						detail: `
              <div class="row">
                <div class="col-4">สายทาง</div>
                <div class="col colon">:</div>
                <div class="col">${item.name}</div>
              </div>
              <div class="row">
                <div class="col-4">ตอนควบคุม</div>
                <div class="col colon">:</div>
                <div class="col">${item.section}</div>
              </div>
              <div class="row">
                <div class="col-4">ถนน</div>
                <div class="col colon">:</div>
                <div class="col">${item.road}</div>
              </div>
							${this.checkRoadLevelDetail(item.level, item.originToDestination)}
              <div class="row">
                <div class="col-4">รหัสสายทาง</div>
                <div class="col colon">:</div>
                <div class="col">${item.code}</div>
              </div>
              <div class="row">
                <div class="col-4">กม.เริ่มต้น - กม.สิ้นสุด</div>
                <div class="col colon">:</div>
                <div class="col">${item.km}</div>
              </div>
              <div class="row">
                <div class="col-4">ระยะทาง</div>
                <div class="col colon">:</div>
                <div class="col">${item.distance} กม.</div>
              </div>
              <div class="row">
                <div class="col-4">ชนิดผิว</div>
                <div class="col colon">:</div>
                <div class="col">
									${this.generateSurfaceDetail(item.surface)}
							</div>
              </div>
            </div>`,
						size: { width: 300 },
						lineColor: item.color,
					})
				})

				lines.forEach((line) => {
					this.map.Overlays.add(line[0])
				})
			}
		},
		createPopup(latLng: any, data: IRoadListRoads) {
			const regex = /\s*\([^)]+\)/
			// @ts-ignore
			const longdo = window.longdo

			const parent = this.roads.filter((item: any) =>
				item.sections.some((sections: any) => sections.id === data.road_section_id)
			)[0]
			const child: any = ref()
			const roadChild: any = ref()
			const html = ref("")
			this.roads.forEach((item) => {
				item.sections.forEach((sectionItem) => {
					if (data.road_level === 1) {
						if (sectionItem.roads.some((road: any) => road.id === data.id)) {
							child.value = sectionItem
							html.value += `${child.value?.number} ${child.value?.name_origin_th} - ${child.value?.name_destination_th}`
						}
					} else {
						sectionItem.roads.forEach((road) => {
							if (road.child_roads && road.child_roads.some((child: any) => child.id === data.id)) {
								roadChild.value = sectionItem
								html.value += `${roadChild.value?.number} ${roadChild.value?.name_origin_th} - ${roadChild.value?.name_destination_th}`
							}
						})
					}
				})
			})

			const popup = new longdo.Popup(
				{
					lon: latLng.lon,
					lat: latLng.lat,
				},
				{
					detail: `
					<div class="row">
						<div class="col-4">สายทาง</div>
						<div class="col colon">:</div>
						<div class="col">ทางหลวงพิเศษหมายเลข ${parent.number} ${parent.short_name.replace(regex, "")}</div>
					</div>
					<div class="row">
						<div class="col-4">ตอนควบคุม</div>
						<div class="col colon">:</div>
						<div class="col">${html.value}</div>
					</div>
					<div class="row">
						<div class="col-4">ถนน</div>
						<div class="col colon">:</div>
						<div class="col">${data.road_info?.name}</div>
					</div>
					${this.checkRoadLevelDetail(data.road_level, data.road_info?.origin_to_destination)}
					<div class="row">
						<div class="col-4">รหัสสายทาง</div>
						<div class="col colon">:</div>
						<div class="col">${data.road_info?.road_code}</div>
					</div>
					<div class="row">
						<div class="col-4">กม.เริ่มต้น - กม.สิ้นสุด</div>
						<div class="col colon">:</div>
						<div class="col">${convertMeterToKm(data.road_info?.km_start)} - ${convertMeterToKm(data.road_info?.km_end)}</div>
					</div>
					<div class="row">
						<div class="col-4">ระยะทาง</div>
						<div class="col colon">:</div>
						<div class="col">${calculateDistance(data.road_info?.km_start, data.road_info?.km_end)} กม.</div>
					</div>
					<div class="row">
						<div class="col-4">ชนิดผิว</div>
						<div class="col colon">:</div>
						<div class="col">
							${this.generateSurfaceDetail(data.road_surface_icon)}
					</div>
					</div>
					`,
					size: { width: 300 },
				}
			)
			return popup
		},
		setLocation(data: IRoadListRoads) {
			if (this.map) {
				if (data.road_info?.the_geom && data.road_info?.the_geom !== "") {
					const latLng = getLatLong(data.road_info?.the_geom)
					const popup = this.createPopup(latLng, data)
					this.map.location({
						lon: latLng.lon,
						lat: latLng.lat,
					})
					this.map.zoom(18)
					setTimeout(() => {
						this.map.Overlays.add(popup)
					}, 250)
				}
			}
		},
		defaultLocation() {
			if (this.getRoadGeoms.length === 0) {
				return
			}
			const latLng = getLatLong(this.getRoadGeoms[0]?.geom)
			this.map.location({
				lon: latLng.lon,
				lat: latLng.lat,
			})
			this.map.zoom(10)
		},
		extractText(name: string, fullname = false) {
			if (fullname === true) {
				const regex = /\s*\([^)]+\)/
				return name.replace(regex, "")
			} else {
				const match = name.match(/\(([^)]+)\)/)
				return match ? match[0] : null
			}
		},
	},
	getters: {
		getRoads(state) {
			return state.roads
		},
		getRoadGeoms(state) {
			const roads = state.roads ?? []

			const sectionsGeom: geomDetail[] = []
			const childrenGeom: geomDetail[] = []
			const regex = /\s*\([^)]+\)/
			roads.forEach((road) => {
				road.sections?.forEach((section) => {
					if (section.roads) {
						section.roads.forEach((roadSection) => {
							sectionsGeom.push({
								geom: roadSection.road_info?.the_geom,
								color: roadSection.road_info?.road_color_code,
								name: "ทางหลวงพิเศษหมายเลข " + road.number + " " + road.short_name.replace(regex, ""),
								section: section.number.toString() + " " + section.name_origin_th + " - " + section.name_destination_th,
								road: roadSection.road_info?.name,
								originToDestination: roadSection.road_info?.origin_to_destination,
								code: roadSection.road_info?.road_code,
								km:
									convertMeterToKm(roadSection.road_info?.km_start) +
									" - " +
									convertMeterToKm(roadSection.road_info?.km_end),
								distance: calculateDistance(roadSection.road_info?.km_start, roadSection.road_info?.km_end),
								surface: roadSection.road_surface_icon,
								level: roadSection.road_level,
							})

							roadSection.child_roads?.forEach((childRoad) => {
								childrenGeom.push({
									geom: childRoad.road_info?.the_geom,
									color: childRoad?.road_info?.road_color_code,
									name: "ทางหลวงพิเศษหมายเลข" + road.number + road.short_name.replace(regex, ""),
									section:
										section.number.toString() + " " + section.name_origin_th + " - " + section.name_destination_th,
									road: childRoad.road_info?.name,
									originToDestination: childRoad.road_info?.origin_to_destination,
									code: childRoad.road_info?.road_code,
									km:
										convertMeterToKm(childRoad.road_info?.km_start) +
										" - " +
										convertMeterToKm(childRoad.road_info?.km_end),
									distance: calculateDistance(childRoad.road_info?.km_start, childRoad.road_info?.km_end),
									surface: childRoad.road_surface_icon,
									level: childRoad.road_level,
								})
							})
						})
					}
				})
			})

			return [...sectionsGeom, ...childrenGeom].filter(Boolean) ?? []
		},
	},
})
