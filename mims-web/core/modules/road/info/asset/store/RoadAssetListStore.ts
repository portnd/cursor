import {
	IRequestAssetDetail,
	IRoadAssetDetailItem,
	IRoadAssetItem,
	IRoadAssetRevision,
	RoadsAssetService,
	location,
} from "../infrastructure"
import { IAssetDetail } from "../infrastructure/RoadAssetModel.d"
import type { THeader } from "~~/core/shared/types/Datatable"
interface IRoadAssetItems {
	[key: string]: any // สามารถใช้ type อื่นแทน any
}
interface ILineString {
	geom_cl: string
	road_color_code: string
}
interface IState {
	id: number
	canEdit: boolean
	refAssetTableId: number
	revisionSelectId: number
	revisionIdParent: number
	loading: boolean
	dataRevision: IRoadAssetRevision[]
	dataAsset: IRoadAssetDetailItem
	headers: THeader[]
	dataGeom: IRoadAssetItem[]
	statusLock: boolean
	dataMap: IRoadAssetItems[]
	map: any
	showVisible: boolean
	locationCurrent: string
	lineString: ILineString
	detail: IRoadAssetDetailItem
	assetType: string
	assetDetail: IAssetDetail
}

export const useRoadAssetListStore = defineStore("roads/asset-in-out/list", {
	state: (): IState => ({
		id: 0,
		canEdit: false,
		refAssetTableId: 0,
		revisionSelectId: 0,
		revisionIdParent: 0,
		loading: false,
		dataRevision: [] as IRoadAssetRevision[],
		dataAsset: {} as IRoadAssetDetailItem,
		headers: [],
		dataGeom: [] as IRoadAssetItem[],
		statusLock: false,
		dataMap: [] as IRoadAssetItems[],
		map: null,
		showVisible: true,
		locationCurrent: "",
		lineString: {} as ILineString,
		detail: {} as IRoadAssetDetailItem,
		assetType: "",
		assetDetail: {} as IAssetDetail,
	}),
	actions: {
		resetStore() {
			this.revisionSelectId = 0
			this.revisionIdParent = 0
			this.dataRevision = []
			this.dataAsset = {} as IRoadAssetDetailItem
			this.headers = []
			this.dataGeom = []
			this.dataMap = []
		},
		async getRevision() {
			// Loading
			this.loading = true
			const assetService = new RoadsAssetService()
			const res = await assetService.getRevision(this.id, this.refAssetTableId)
			const resAssetDetail = await assetService.getAssetDetail(this.refAssetTableId)
			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.dataRevision = res.data as IRoadAssetRevision[]
				if (this.dataRevision.length > 0) {
					this.revisionSelectId = this.dataRevision[0].id
					this.revisionIdParent = this.dataRevision[0].id_parent
				} else {
					this.resetStore()
				}
				this.assetDetail = resAssetDetail.data
				return this.revisionSelectId
			}
		},
		async getDataTable(revisionSelectId: number) {
			// Loading
			if (revisionSelectId === 0) {
				return null
			} else {
				this.loading = true
				this.statusLock = this.dataRevision.some((e) => {
					return revisionSelectId === e.id && e.is_exclusive_lock
				})

				const params: IRequestAssetDetail = {
					ref_asset_table_id: this.refAssetTableId,
					page: "1",
					limit: "9999",
				}

				const assetService = new RoadsAssetService()
				const res = await assetService.getAssetList(this.id, revisionSelectId, params)

				// Loading
				this.loading = false

				if (res.status === false) {
					useHandlerError(res.code, res.error, { showAlert: true })
				} else {
					this.$patch((state) => (state.dataAsset = res.data.items))
					this.statusLock = this.dataAsset.is_exclusive_lock
					this.handledata()
				}
			}
		},
		async putCancel() {
			// Loading
			this.loading = true

			const assetService = new RoadsAssetService()
			const res = await assetService.putAssetCancel(this.revisionIdParent)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				const res = await this.getRevision()
				await this.getDataTable(res ?? 1)
			}
		},
		async putConfirm() {
			// Loading
			this.loading = true

			const assetService = new RoadsAssetService()
			const res = await assetService.putAssetConfirm(this.revisionIdParent)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		setMap(map: Object) {
			this.$patch((state) => (state.map = map))
		},
		setLocation(item: IRoadAssetItem) {
			if (item !== undefined) {
				const data: location = item.geom_2d
				this.map.location({
					lon: data.lon,
					lat: data.lat,
				})
				this.map.zoom(18)
			}
		},
		handledata() {
			if (this.dataGeom.length !== 0) {
				this.dataGeom = []
			}
			if (this.dataMap.length !== 0) {
				this.dataMap = []
			}
			let headers: { text: string; value: string; type: string }[] = []
			this.dataAsset.data_columns.map((e) => {
				if (e.key.includes("_filepath")) {
					headers.push({ text: e.value.component_title, value: e.key, type: e.value.component_type })
				} else if (e.key.includes("_id")) {
					headers.push({ text: e.value.component_title, value: e.key.split("_id")[0], type: e.value.component_type })
				} else {
					headers.push({ text: e.value.component_title, value: e.key, type: e.value.component_type })
				}
				return ""
			})

			headers = headers.filter((e) => !(e.value === "geom_2d" || e.value === "geom_3d" || e.value === "geom_camera"))
			headers.push({ text: "จัดการ", value: "action", type: "action" })

			if (this.dataAsset.status_code === "I" || !this.canEdit) {
				this.headers = headers.filter((item) => item.type !== "action")
			} else {
				this.headers = headers
			}

			// body
			for (let i = 0; i < this.dataAsset.road_assets.length; i++) {
				interface IRoadAssetItems {
					// สามารถใช้ type อื่นแทน any
					[key: string]: any
				}
				const object: IRoadAssetItems = {} as IRoadAssetItems
				const geomObject: IRoadAssetItems = {} as IRoadAssetItems

				for (const [key, value] of Object.entries(this.dataAsset.road_assets[i])) {
					if (typeof value === "object") {
						if (value !== null) {
							try {
								if (value.name !== null) {
									if (key === "ref_asset_sign_image") {
										object.ref_asset_sign_image = value
									} else {
										object[key] = value.name
									}
								} else {
									object[key] = ""
								}
							} catch {
								object[key] = "error"
							}
						} else {
							object[key] = ""
						}
					} else if (key === "geom_2d") {
						const geom = value

						let text = geom.split("(")[1]
						const geomType = geom.split("(")[0]
						text = text.split(")")[0]
						text = text.split(" ")
						const loc = {
							lon: parseFloat(text[0]),
							lat: parseFloat(text[1]),
						}
						object[key] = loc
						geomObject.geom = loc
						geomObject.geom_wkt = value
						geomObject.geom_type = geomType
					} else if (key === "geom_3d") {
						object[key] = value
						const geom3d = value

						let text = geom3d.split("(")[1]
						const geomType = geom3d.split(" (")[0]
						text = text.split(")")[0]
						text = text.split(",")

						const locArr = []
						for (let i = 0; i < text.length; i++) {
							const locWktSplit = text[i].split(" ")
							const loc = {
								lon: parseFloat(locWktSplit[0]),
								lat: parseFloat(locWktSplit[1]),
								alt: parseFloat(locWktSplit[2]),
							}
							locArr.push(loc)
						}

						object[key] = locArr
						geomObject.geom_3d = locArr
						geomObject.geom_3d_wkt = value
						geomObject.geom_3d_type = geomType
					} else if (key === "geom_camera") {
						object[key] = value
						const geomCamera = value
						let text = geomCamera.split("(")[1]

						text = text.split(")")[0]
						text = text.split(" ")
						const loc = {
							lon: parseFloat(text[0]),
							lat: parseFloat(text[1]),
							alt: parseFloat(text[2]),
						}
						object[key] = loc
						geomObject.geom_camera = loc
					} else if (value !== null) {
						if (key.includes("_filepath")) {
							object[key] = value
							// if (value !== "") {
							//  object[key] = value
							// } else {
							//  object[key] = this.dataAsset.thumbnail_icon_filepath
							// }
							geomObject.sign_icon = object[key]
						} else if (key.includes("km")) {
							if (value || value === 0) {
								object[key] = value
							} else {
								object[key] = ""
							}
							geomObject.km = value === "" ? "" : convertMeterToKm(value)
						} else if (key === "km_start") {
							object[key] = value
							geomObject.km_start = convertMeterToKm(value)
						} else if (key === "km_end") {
							object[key] = value
							geomObject.km_end = convertMeterToKm(value)
						} else if (key.includes("_date")) {
							object[key] = buddhistFormatDate(value, "dd mmm yyyy")
							geomObject.date = buddhistFormatDate(value, "dd mmm yyyy")
						} else if (key.includes("ref_asset_sign")) {
							if (!value) {
								object[key] = ""
							} else {
								object[key] = value.id
							}
						} else if (key === "longitude" || key === "latitude") {
							object[key] = value.toFixed(15)
						} else if (key === "altitude") {
							object[key] = (0).toFixed(15)
						} else {
							object[key] = value
						}
					} else {
						object[key] = ""
					}
				}
				if (object.ref_asset_sign_image) {
					if (object.ref_asset_sign_image.filepath) {
						object.sign_filepath = object.ref_asset_sign_image.filepath
						geomObject.sign_icon = object.ref_asset_sign_image.filepath
					} else {
						object.sign_filepath = object.ref_asset_sign_image.thumbnail_filepath
						geomObject.sign_icon = object.ref_asset_sign_image.thumbnail_filepath
					}
				} else {
					geomObject.sign_icon = this.dataAsset.icon_filepath
				}
				if (object.ref_asset_sign) {
					if (object.ref_asset_sign === "ป้ายทั่วไป" || object.ref_asset_sign === "ป้าย LED") {
						object.sign_filepath = object.ref_asset_sign_image.filepath
						object.sign_caption = object.ref_asset_sign_image.name
						geomObject.sign_icon = object.ref_asset_sign_image.filepath
					} else if (object.ref_asset_sign === "ป้าย Overhang" || object.ref_asset_sign === "ป้าย Overhead") {
						geomObject.sign_icon = object.thumbnail_sign_filepath
					}
				} else {
					geomObject.sign_icon = this.dataAsset.icon_filepath
				}
				if (this.dataAsset.line_color_code) {
					geomObject.line_color = this.dataAsset.line_color_code
				} else {
					geomObject.line_color = "red"
				}

				geomObject.asset_id = this.revisionSelectId
				object.action = "action"
				object.raw_data = this.dataAsset.road_assets[i]
				this.dataGeom.push(object as IRoadAssetItem)
				if (Object.keys(geomObject).length > 0) {
					this.dataMap.push(geomObject)
				}
			}
		},

		draw_asset_geom() {
			this.map.Overlays.clear()
			for (let i = 0; i < this.dataMap.length; i++) {
				let kmText = ""
				let latLonText = ""
				let dateText = ""
				let googlemapUrl = ""

				if (this.dataMap[i].km) {
					kmText = "กม. " + this.dataMap[i].km
				}
				if (this.dataMap[i].km_start && this.dataMap[i].km_end) {
					kmText = "ช่วง กม. " + this.dataMap[i].km_start + " - " + this.dataMap[i].km_end
				}
				if (this.dataMap[i].geom) {
					latLonText = "<b>ละติจูด :</b> " + this.dataMap[i].geom.lat + "</br>"
					latLonText += "<b>ลองจิจูด :</b> " + this.dataMap[i].geom.lon + "</br>"

					googlemapUrl =
						'<a href="https://www.google.co.th/maps?cbll=' +
						this.dataMap[i].geom.lat +
						"," +
						this.dataMap[i].geom.lon +
						'&layer=c" target="_blank">Google Maps</a></br>'
				}
				if (this.dataMap[i].date) {
					dateText = "<b>วันที่สำรวจ:</b> " + this.dataMap[i].date + "</br>"
				}

				if (this.dataMap[i].geom_type === "POINT") {
					const geom = this.dataMap[i].geom
					const dataMap = this.dataMap[i]
					const keysIncludeFilepath = this.dataAsset.data_columns
						.filter((item: any) => item.key.includes("filepath"))
						.map((item: any) => item.key)
					const arr = ref([] as any)
					if (keysIncludeFilepath.length !== 0) {
						for (let y = 0; y < keysIncludeFilepath.length; y++) {
							const roadAsset = this.dataAsset?.road_assets[i]
							arr.value.push(roadAsset?.[keysIncludeFilepath[y]])
						}
						arr.value = arr.value.filter((item: any) => item !== "")
						dataMap.sign_icon = arr.value[0]
						if (arr.value.length === 0) {
							if (
								this.dataAsset?.road_assets[i].ref_asset_sign_image?.id === 0 &&
								this.dataAsset?.icon_filepath !== ""
							) {
								dataMap.sign_icon = this.dataAsset?.icon_filepath
							} else if (this.dataAsset?.road_assets[i].ref_asset_sign_image?.id > 0) {
								const image = useInitData()
									.refAssetSignImage()
									?.find((image: any) => image.id === this.dataAsset?.road_assets[i].ref_asset_sign_image.id)
								dataMap.sign_icon = image!.sign_image_filepath
							} else {
								dataMap.sign_icon = "/images/icons/png/location-pin.png"
							}
						}
					} else if (keysIncludeFilepath.length === 0) {
						if (this.dataAsset?.road_assets[i].ref_asset_sign_image) {
							if (
								this.dataAsset?.road_assets[i].ref_asset_sign_image.id === 0 &&
								this.dataAsset?.icon_filepath !== ""
							) {
								dataMap.sign_icon = this.dataAsset.icon_filepath
							} else if (this.dataAsset?.road_assets[i].ref_asset_sign_image.id > 0) {
								const image = useInitData()
									.refAssetSignImage()
									?.find((image: any) => image.id === this.dataAsset?.road_assets[i].ref_asset_sign_image.id)
								dataMap.sign_icon = image!.sign_image_filepath
							} else {
								dataMap.sign_icon = "/images/icons/png/location-pin.png"
							}
						}
					}
					if (this.map) {
						let geomImg = this.dataMap[i].sign_icon
						const image = new Image()
						image.src = geomImg
						image.addEventListener("error", (e) => {
							if (e.type === "error") {
								geomImg = "/images/icons/png/location-pin.png"
								// @ts-ignore
								const longdo = window.longdo
								const marker = new longdo.Marker(geom, {
									icon: {
										html: `<img src="${geomImg}" style="width:25px; height:25px"></img>`,
										offset: { x: 12.5, y: 20.5 },
									},
									visible: true,
									visibleRange: { min: 13, max: 20 },
									popup: {
										html: `<div class="popup"><div class="title">${kmText}</div> ${
											latLonText + dateText + googlemapUrl
										}</div>`,
										closable: false,
									},
								})
								this.map.Overlays.add(marker)
							}
						})
						image.addEventListener("load", () => {
							// @ts-ignore
							const longdo = window.longdo
							const marker = new longdo.Marker(geom, {
								icon: {
									html: `<img src="${geomImg}" style="width:25px; height:25px"></img>`,
									offset: { x: 12.5, y: 20.5 },
								},
								visible: true,
								visibleRange: { min: 13, max: 20 },
								popup: {
									html: `<div class="popup"><div class="title">${kmText}</div> ${
										latLonText + dateText + googlemapUrl
									}</div>`,
									closable: false,
								},
							})

							this.map.Overlays.add(marker)
						})
					}
				} else if (this.dataMap[i].geom_type === "LINESTRING") {
					const geom = this.dataMap[i].geom_wkt
					// @ts-ignore
					const longdo = window.longdo
					const wk = longdo.Util.overlayFromWkt(geom, {
						lineColor: this.dataMap[i].line_color,
						title: kmText,
						detail: latLonText + dateText + googlemapUrl,
					})
					wk.forEach((x: any) => {
						this.map.Overlays.add(x)
					})
				}
			}
			if (this.map) {
				const lines = this.lineString.geom_cl.split("|")
				lines.forEach((e) => {
					// @ts-ignore
					const line = longdo.Util.overlayFromWkt(e, {
						lineColor: this.lineString.road_color_code,
					})
					this.map.Overlays.add(line[0])
				})
				if (this.dataMap.length === 0) {
					const latLng = getLatLong(lines[0])
					this.map.location({
						lon: latLng.lon,
						lat: latLng.lat,
					})
					this.map.zoom(18)
				}
			}
		},
	},
	getters: {},
})
