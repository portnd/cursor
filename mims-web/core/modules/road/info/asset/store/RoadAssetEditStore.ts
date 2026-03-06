import { IRequestRoadsAsset, IRoadAssetItem, IRoadAssetTableTemplate, RoadsAssetService } from "../infrastructure"

interface ILocation {
	lat: string
	lon: string
}

interface ILatLng {
	lat: number
	lon: number
}

interface IState {
	id: number
	assetType: string
	location: ILocation
	refAssetTableId: number
	revisionSelectId: number
	idParentAsset: number
	revisionIdParent: number
	loading: boolean
	hideSignImage: boolean
	template: IRoadAssetTableTemplate[]
	imageUrlList: any[]
	actionType: string
	assetObjectId: number
	lineString: string
	geomType: string
	rawData: any[]
	originTextKm: number | null
	roadGeom: ILatLng[]
}

export const useRoadAssetEditStore = defineStore("roads/asset-in-out/edit", {
	state: (): IState => ({
		id: 0,
		assetType: "",
		refAssetTableId: 0,
		revisionSelectId: 0,
		revisionIdParent: 0,
		idParentAsset: 0,
		location: {} as ILocation,
		loading: false,
		hideSignImage: false,
		template: [] as IRoadAssetTableTemplate[],
		imageUrlList: [],
		actionType: "edit", // ใช้สำหรับบอกให้ backend รู้ว่าทำในโหมดอะไร แนบไปกับ api template
		assetObjectId: 0,
		lineString: "",
		geomType: "",
		rawData: [],
		originTextKm: null,
		roadGeom: [],
	}),
	actions: {
		setRoadGeom(stringLine: string) {
			this.roadGeom = parseLineString(stringLine).map((e) => {
				return { lat: e.lat, lon: e.lon } as ILatLng
			})
		},
		resetTemplate() {
			this.template = []
			this.imageUrlList = []
		},
		checkSignImage(item: any) {
			if (item.component_title === "ประเภทป้าย") {
				if (item.value === 1 || item.value === 2) {
					this.hideSignImage = true
				} else {
					this.hideSignImage = false
				}
			}
		},
		pushDataToTemplate(item: IRoadAssetItem) {
			const object = Object.entries(item)
			this.rawData = object
			let latlon: ILocation = {} as ILocation
			console.log("pushDataToTemplate object:", object)
			object.forEach((e) => {
				if (e[0] === "geom_2d") {
					const data = e[1].split(" ")
					const type = data[0].split("(")[0]
					this.geomType = type
					if (type === "POINT") {
						latlon = { lon: data[0].split("(")[1], lat: data[1].split(")")[0] }
					}

					if (type === "LINESTRING") {
						this.lineString = e[1]
					}
				}
			})
			this.location = latlon
			if (this.template) {
				this.template.map((e) => {
					object.map((item) => {
						if (e.column_name === item[0] && e.column_name.includes("_filepath")) {
							this.imageUrlList.push({ name: item[0], path: item[1] })
						}
						return ""
					})
					return ""
				})
			}
		},
		async getTemplateData() {
			// Loading
			this.loading = true
			const assetService = new RoadsAssetService()
			const res = await assetService.getTemplate(this.refAssetTableId, this.actionType, this.assetObjectId)
			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				res.data = res.data?.filter((item) => item.component_type !== "hidden")
				res.data = res.data?.filter((item) => item.column_name !== "the_geom")
				res.data = res.data?.filter((item) => item.column_name !== "the_geom_camera")
				this.template = res.data
				for (let index = 0; index < this.template?.length; index++) {
					if (this.template[index].component_type === "text-km") {
						this.template[index].value =
							this.template[index].value === null ? "" : convertMeterToKm(this.template[index].value)
						this.originTextKm = this.template[index].value
					} else if (this.template[index].column_name === "latitude") {
						this.template[index].value = Number(this.template[index].value).toFixed(15)
					} else if (this.template[index].column_name === "longitude") {
						this.template[index].value = Number(this.template[index].value).toFixed(15)
					} else if (this.template[index].column_name === "altitude") {
						if (this.template[index].value === 0) {
							this.template[index].value = (0).toFixed(15)
						} else {
							this.template[index].value = Number(this.template[index].value)
						}
					} else if (this.template[index].component_type === "text-year") {
						this.template[index].value = Number(this.template[index].value)
					} else if (this.template[index].component_type === "text-number") {
						this.template[index].value = Number(this.template[index].value)
					} else if (this.template[index].column_name.split("_")[1] === "filepath") {
						this.template[index].component_type = "image"
					} else if (this.template[index].component_title === "ประเภทป้าย") {
						this.hideSignImage = true
						this.template[index].value = 1
					}
				}
				return res
			}
		},
		checkDataIsNotShow() {
			interface MyObject {
				[key: string]: any
			}
			const data: MyObject = {}
			for (const b of this.rawData) {
				if (this.template.filter((item) => item.column_name === b[0]).length === 0) {
					if (
						b[0] !== "asset_object_id" &&
						!b[0].includes("geom_") &&
						b[0] !== "id_parent_asset" &&
						!b[0].includes("thumbnail")
					) {
						if (b[0].includes("ref")) {
							const name = b[0] + "_id"
							data[name] = b[1].id
						} else {
							data[b[0]] = b[1]
						}
					}
				}
			}

			return data
		},
		async putAssetRoad(id: number) {
			interface MyObject {
				[key: string]: any
			}
			let data: MyObject = {}
			const hideData = this.checkDataIsNotShow()
			data = Object.assign({}, hideData)
			for (const e of this.template) {
				if (e.component_type === "text-km") {
					data[e.column_name] = e.value ? convertStringToKm(e.value) : undefined
				} else if (e.column_name === "latitude") {
					data[e.column_name] = Number(e.value).toFixed(15)
				} else if (e.column_name === "longitude") {
					data[e.column_name] = Number(e.value).toFixed(15)
				} else if (e.column_name === "altitude") {
					data[e.column_name] = Number(e.value).toFixed(15)
				} else if (e.component_type === "datepicker") {
					data[e.column_name] = e.value ? formatDate(e.value) + " 00:00:00" : undefined
				} else if (e.column_name.includes("_filepath")) {
					if (e.value?.data) {
						data[e.column_name] = e.value.data.base64
					} else {
						data[e.column_name] = ""
					}
				} else if (e.component_type === "text-nember") {
					data[e.column_name] = e.value === 0 ? 0 : Number(e.value)
				} else if (e.component_type === "text-year") {
					data[e.column_name] = e.value === 0 ? 0 : Number(e.value)
				} else if (e.value) {
					data[e.column_name] = e.value
				} else {
					data[e.column_name] = undefined
				}
			}
			if (data.sign_filepath) {
				delete data.ref_asset_sign_image_id
			}
			if (data.ref_asset_sign_image_id) {
				delete data.sign_caption
				delete data.sign_filepath
			}

			data = Object.assign(data, { ref_asset_table_id: this.refAssetTableId })

			let theGeom = ""
			if (this.geomType === "POINT") {
				const currentTextKm = this.template.find((e) => e.component_type === "text-km")
				if (currentTextKm?.value && this.originTextKm && this.originTextKm !== currentTextKm?.value) {
					const point = getLatLngByDistance(this.roadGeom, data[currentTextKm.column_name])
					if (point) {
						theGeom = `POINT(${point?.lon} ${point?.lat})`
					}
				} else {
					theGeom = `POINT(${this.location.lon} ${this.location.lat})`
				}
			}

			if (this.geomType === "LINESTRING") {
				theGeom = this.lineString
			}

			data = Object.assign(data, { the_geom: theGeom })

			if (this.revisionSelectId !== 0) {
				data = Object.assign(data, { id_parent: this.revisionIdParent })
				data = Object.assign(data, { road_asset_id: this.revisionSelectId })
			}
			const params: IRequestRoadsAsset = data

			// Loading
			this.loading = true

			const assetService = new RoadsAssetService()
			const res = await assetService.putAssetRoad(id, this.idParentAsset, params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
	},
	getters: {},
})
