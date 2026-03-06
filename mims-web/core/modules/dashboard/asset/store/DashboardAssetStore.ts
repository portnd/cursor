import { IDashboardAsset, IAssetLocation, DashboardAssetService } from "../infrastructure"
import { ITree } from "~/core/shared/types/Tree"
interface GeomWkt {
	geom: string
	color: string
}
interface Location {
	lon: number
	lat: number
}

interface IState {
	data: IDashboardAsset[]
	asset: IAssetLocation[]
	roadsTree: ITree[]
	roadId: string
	assetId: string
	map?: any
	location: Location
	roadsOptions: Array<any>
	markers: Array<any>
	markercluster: null
	loading: boolean
}

export const useDashboardAssetStore = defineStore("dashboard/asset", {
	state: (): IState => ({
		data: [] as IDashboardAsset[],
		asset: [] as IAssetLocation[],
		roadsTree: [] as ITree[],
		roadId: "",
		assetId: "",
		map: null,
		location: {
			lon: 100.737091,
			lat: 13.730756,
		},
		roadsOptions: [],
		markers: [],
		markercluster: null,
		loading: false,
	}),
	actions: {
		async get(assetType: string) {
			// Loading
			this.loading = true

			const dashboardAssetService = new DashboardAssetService()
			const params = ref("")
			if (this.roadId === "") {
				params.value = `asset_type=${assetType}`
			} else {
				params.value = `road_id=${this.roadId}&asset_type=${assetType}`
			}
			const res = await dashboardAssetService.get(params.value)

			setTimeout(() => {
				// Loading
				this.loading = false
			}, 200)

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.data = res.data
				return res
			}
		},
		async getLocation() {
			// Loading
			// this.loading = true

			const dashboardAssetService = new DashboardAssetService()
			const params = ref("")
			if (this.roadId === "") {
				params.value = `asset_id=${this.assetId}`
			} else {
				params.value = `road_id=${this.roadId}&asset_id=${this.assetId}`
			}
			const res = await dashboardAssetService.getLocation(params.value)

			// Loading
			// this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.asset = res.data
				this.mapCluster()
				return res
			}
		},
		async getRoadsTree() {
			// Loading
			this.loading = true
			const dashboardAssetService = new DashboardAssetService()
			const res = await dashboardAssetService.getRoadsTree()
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.roadsTree = res.data
				return res
			}
			// Loading
			this.loading = false
		},
		createRoadsOptions() {
			const options = this.roadsTree.map((parent: any) => {
				return {
					id: "m" + parent.id,
					label: parent.label,
					children: parent.children.map((child: any) => {
						return { id: child.id, label: child.label }
					}),
				}
			})

			const result = [{ id: "all", label: "ทั้งหมด", children: options }]

			this.roadsOptions = result
		},
		setMap(map: any) {
			this.$patch((state) => {
				state.map = map
			})
			this.map?.zoom(12)
		},
		setLocation() {
			if (this.map) {
				const item = this.asset[0]
				if (item?.wkt) {
					const cleanGeom = getLatLong(item.wkt)
					this.map.location({
						lon: cleanGeom.lon,
						lat: cleanGeom.lat,
					})
					this.map?.zoom(12)
				}
			}
		},
		setIcon(img: string) {
			return {
				html: `<img src="${img}" style="width:25px; height:25px"></img>`,
				offset: { x: 10.5, y: 18.5 },
			}
		},
		clearMarker() {
			if (this.map) {
				if (this.markercluster !== null && this.markercluster !== undefined) {
					// @ts-ignore
					this.markercluster.clearMarkers()
					this.markers = []
					this.markercluster = null
				}
				this.map.Overlays.clear()
			}
		},
		mapCluster() {
			if (this.map) {
				const script = document.createElement("script")
				script.src = "/js/longdomap.markercluster-src.js"
				script.id = "markercluster"
				document.head.appendChild(script)
				this.markers = []
				script.onload = () => {
					// @ts-ignore
					const longdo = window.longdo
					const roads = ref<Array<any>>([])
					if (this.markercluster === null) {
						// @ts-ignore
						this.markercluster = new lmc.MarkerCluster(this.map, {
							maxZoom: 15,
							minZoom: 6,
							minClusterSize: 1,
						})
					}
					for (const index in this.asset) {
						const geom = this.asset[index].wkt as any
						if (geom) {
							if (this.asset[index].type === "POINT") {
								let icon = null
								if (this.asset[index].icon_filepath != null) {
									icon = this.setIcon(this.asset[index].icon_filepath)
								} else if (this.asset[index].thumbnail_icon_filepath != null) {
									icon = this.setIcon(this.asset[index].thumbnail_icon_filepath)
								} else {
									icon = this.setIcon("/images/icons/png/location-pin.png")
								}
								if (icon) {
									const location = getLatLong(geom)
									if (location) {
										this.markers.push(new longdo.Marker({ lon: location.lon, lat: location.lat }, { icon }))
									}
								}
							}
						}
					}
					const result = this.asset.filter((item: any) => item.type === "LINESTRING")
					result.forEach((item: any) => {
						roads.value.push({
							geom: item.wkt,
							color: item.line_color_code,
						})
					})
					const lines = roads.value.map(({ geom, color }: GeomWkt) => {
						return longdo.Util.overlayFromWkt(geom, { lineColor: color })
					})

					lines.forEach((line: Array<Object>) => {
						this.map.Overlays.add(line[0])
					})

					// @ts-ignore
					this.markercluster.addMarkers(this.markers)
					// @ts-ignore
					this.markercluster.render()
				}
			}
		},
	},
	getters: {},
})
