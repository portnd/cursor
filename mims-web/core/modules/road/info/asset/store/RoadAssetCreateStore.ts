import { IRequestRoadsAsset, IRoadAssetTableTemplate, Km, RoadsAssetService, location } from "../infrastructure"
import { IResponse } from "~/core/shared/http"

interface IState {
	id: number
	assetType: string
	location: location[]
	refAssetTableId: number
	revisionSelectId: number
	revisionIdParent: number
	loading: boolean
	hideSignImage: boolean
	template: IRoadAssetTableTemplate[]
	pointStart: string
	pointEnd: string
	actionType: string
}

export const useRoadAssetCreateStore = defineStore("roads/asset-in-out/create", {
	state: (): IState => ({
		id: 0,
		assetType: "",
		refAssetTableId: 0,
		revisionSelectId: 0,
		revisionIdParent: 0,
		location: [],
		loading: false,
		hideSignImage: false,
		template: [] as IRoadAssetTableTemplate[],
		pointStart: "",
		pointEnd: "",
		// ใช้สำหรับบอกให้ backend รู้ว่าทำในโหมดอะไร แนบไปกับ api template
		actionType: "add",
	}),
	actions: {
		resetTemplate() {
			this.template = []
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
		async getTemplateData(location: location[], geomType: number) {
			this.location = location
			let paramsKmStart = ""
			let paramsEnd = ""
			let resKmStart = {} as IResponse<Km>
			let resKmEnd = {} as IResponse<Km>
			let kmStart = null
			let kmEnd = null
			if (location.length > 1) {
				paramsKmStart = `POINT(${location[0].lon} ${location[0].lat})`
				paramsEnd = `POINT(${location[1].lon} ${location[1].lat})`
				// `LINESTRING(${location[0].lon} ${location[0].lat},${location[1].lon} ${location[1].lat})`
			} else {
				paramsKmStart = `POINT(${location[0].lon} ${location[0].lat})`
			}

			// Loading
			this.loading = true
			const assetService = new RoadsAssetService()
			if (geomType !== 3) {
				resKmStart = await assetService.getKm(this.id, paramsKmStart)
				if (location.length > 1) {
					resKmEnd = await assetService.getKm(this.id, paramsEnd)
				}
				kmStart = convertMeterToKm(Number(resKmStart.data.km))
				if (location.length > 1) {
					kmEnd = convertMeterToKm(Number(resKmEnd.data.km))
					this.pointEnd = resKmEnd.data.the_geom
				}
				this.pointStart = resKmStart.data.the_geom
			} else {
				this.pointStart = paramsKmStart
			}
			const res = await assetService.getTemplate(this.refAssetTableId, this.actionType)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				res.data = res.data.filter((item) => item.component_type !== "hidden")
				res.data = res.data.filter((item) => item.column_name !== "the_geom")
				res.data = res.data.filter((item) => item.column_name !== "the_geom_camera")
				this.template = res.data
				for (let index = 0; index < this.template.length; index++) {
					if (this.template[index].column_name === "km") {
						this.template[index].value = kmStart
					} else if (this.template[index].column_name === "km_start") {
						this.template[index].value = kmStart
					} else if (this.template[index].column_name === "km_end") {
						this.template[index].value = kmEnd
					} else if (this.template[index].column_name === "latitude") {
						this.template[index].value = Number(this.location[0].lat).toFixed(15)
					} else if (this.template[index].column_name === "longitude") {
						this.template[index].value = Number(this.location[0].lon).toFixed(15)
					} else if (this.template[index].column_name === "altitude") {
						this.template[index].value = (0).toFixed(15)
					} else if (
						this.template[index].component_type === "text-number" ||
						this.template[index].component_type === "text-year"
					) {
						this.template[index].value = ""
					} else if (this.template[index].column_name.split("_")[1] === "filepath") {
						this.template[index].component_type = "image"
					} else if (this.template[index].component_title === "ประเภทป้าย") {
						this.hideSignImage = true
						// this.template[index].value = 1
					}
				}
				return res
			}
		},
		async postAssetRoad() {
			interface MyObject {
				[key: string]: any
			}

			let data: MyObject = {}
			for (const e of this.template) {
				if (e.component_type === "text-km") {
					if (e.value === "") {
						data[e.column_name] = undefined
					} else if (e.value !== "") {
						data[e.column_name] = convertStringToKm(e.value)
					}
				} else if (e.column_name === "km_start") {
					data[e.column_name] = e.value ? convertStringToKm(e.value) : undefined
				} else if (e.column_name === "km_end") {
					data[e.column_name] = e.value ? convertStringToKm(e.value) : undefined
				} else if (e.column_name === "latitude") {
					data[e.column_name] = Number(e.value)
				} else if (e.column_name === "longitude") {
					data[e.column_name] = Number(e.value)
				} else if (e.column_name === "altitude") {
					data[e.column_name] = e.value
				} else if (e.component_type === "datepicker") {
					data[e.column_name] = e.value ? formatDate(e.value) + " 00:00:00" : undefined
				} else if (e.component_type === "text-year") {
					data[e.column_name] = Number(e.value)
				} else if (e.column_name.includes("_filepath")) {
					if (e.value?.data) {
						data[e.column_name] = e.value.data.base64
					} else {
						data[e.column_name] = ""
					}
				} else if (e.component_type === "text-nember") {
					data[e.column_name] = e.value === 0 ? 0 : Number(e.value)
				} else if (e.value) {
					data[e.column_name] = e.value
				} else if (e.component_type === "text") {
					data[e.column_name] = e.value ? e.value : ""
				} else {
					data[e.column_name] = null
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
			if (this.assetType === "in") {
				if (this.location.length > 1) {
					const geomStartlat = this.pointStart.split("(")[1].split(" ")[0]
					const geomStartlon = this.pointStart.split(")")[0].split(" ")[1]
					const geomEndlat = this.pointEnd.split("(")[1].split(" ")[0]
					const geomEndlon = this.pointEnd.split(")")[0].split(" ")[1]
					const theGeom = `LINESTRING(${geomStartlat} ${geomStartlon},${geomEndlat} ${geomEndlon})`
					data = Object.assign(data, { the_geom: theGeom })
				} else {
					data = Object.assign(data, { the_geom: this.pointStart })
				}
			} else {
				const theGeom = `POINT(${data.longitude} ${data.latitude})`

				data = Object.assign(data, { the_geom: theGeom })
			}

			if (this.revisionSelectId !== 0) {
				data = Object.assign(data, { id_parent: this.revisionIdParent })
				data = Object.assign(data, { road_asset_id: this.revisionSelectId })
			}
			const params: IRequestRoadsAsset = data

			// Loading
			this.loading = true
			const assetService = new RoadsAssetService()
			const res = await assetService.postAssetRoad(this.id, params)
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
