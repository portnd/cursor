import { IRoadsAssetItem, RoadsAssetService } from "../infrastructure"

interface IState {
	assetType: string
	loading: boolean
	menu: IRoadsAssetItem[]
	map: any
	geom: string
	color: string
}

export const useRoadAssetItemStore = defineStore("roads/asset-in-out/item", {
	state: (): IState => ({
		assetType: "",
		loading: false,
		menu: [] as IRoadsAssetItem[],
		map: null,
		geom: "",
		color: "",
	}),
	actions: {
		async getAsset() {
			// Loading
			this.loading = true
			const assetService = new RoadsAssetService()
			const res = await assetService.getMenu(this.assetType)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.menu = res.data as IRoadsAssetItem[]
			}
		},
		setMap(map: Object) {
			this.$patch((state) => {
				state.map = map
			})
			// วาดเส้นและเซ็ตตำแหน่งเมื่อมี geom แล้ว (ป้องกันแผนที่สีฟ้าเมื่อ map พร้อมก่อนข้อมูล)
			if (this.geom) {
				this.defaultLocation()
				this.createLine()
			}
		},
		defaultLocation() {
			if (!this.map || !this.geom) {
				return
			}
			const cleanLeftGeom = getLatLong(this.geom)
			this.map.location({
				lon: cleanLeftGeom.lon,
				lat: cleanLeftGeom.lat,
			})
			this.map?.zoom(18)
		},
		createLine() {
			if (!this.map || !this.geom) {
				return
			}
			// @ts-ignore
			const longdo = window.longdo
			const line = longdo.Util.overlayFromWkt(this.geom, { lineColor: this.color })
			if (line?.[0]) {
				this.map.Overlays.add(line[0])
			}
		},
	},
	getters: {},
})
