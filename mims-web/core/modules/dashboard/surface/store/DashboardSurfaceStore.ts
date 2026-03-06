import { IDashboardSurface, DashboardSurfaceService } from "../infrastructure"
import { IRoadData, RoadListService } from "~/core/modules/road/roadList/infrastructure"
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
	data: IDashboardSurface
	roadsTree: ITree[]
	roadId: string
	map?: any
	location: Location
	geom: { geom_cl: string; color: string }
	color: string
	roadsOptions: Array<any>
	roads: IRoadData[]
	surfaceArray: Array<number>
	surfaceColors: Array<any>
	surfaceCategories: Array<any>
	colors: Array<any>
	dataArray: Array<any>
	loading: boolean
}

export const useDashboardSurfaceStore = defineStore("dashboard/surface", {
	state: (): IState => ({
		data: {} as IDashboardSurface,
		roadsTree: [] as ITree[],
		roadId: "",
		map: null,
		location: {
			lon: 100.737091,
			lat: 13.730756,
		},
		geom: { geom_cl: "", color: "" },
		color: "",
		roadsOptions: [],
		roads: [],
		surfaceArray: [],
		surfaceColors: [],
		surfaceCategories: [],
		colors: [],
		dataArray: [],
		loading: false,
	}),
	actions: {
		async get() {
			// Loading
			this.loading = true

			const dashboardSurfaceService = new DashboardSurfaceService()

			const params = ref("")
			if (this.roadId !== "") {
				params.value = `?road_id=${this.roadId}`
			}
			const res = await dashboardSurfaceService.get(params.value)
			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.data = res.data
				this.surfaceArray = []
				this.colors = []
				this.surfaceColors = this.data.summary.map((item) => item.surface.color_code)
				this.surfaceCategories =
					this.data.summary?.map((item: any, index) => {
						this.surfaceArray.push(index)
						return item.surface.name
					}) ?? []
				this.dataArray = this.data.summary?.map((item: any) => Number(item.value.toFixed(2)))
				return res
			}
		},
		async getRoadsTree() {
			// Loading
			this.loading = true
			const dashboardSurfaceService = new DashboardSurfaceService()
			const res = await dashboardSurfaceService.getRoadsTree()
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.roadsTree = res.data
				return res
			}
			// Loading
			this.loading = false
		},
		async getRoad() {
			this.loading = true
			const query = {
				keyword: "",
				direction: "",
				road_type: "",
			}
			const service = new RoadListService()
			const res = await service.getRoads(query)
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.roads = res.data
			}
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
			this.createLine()
			// this.defaultLocation()
		},
		defaultLocation() {
			if (this.map) {
				if (this.data?.geom_list) {
					const result = this.data?.geom_list.flatMap((item: any) =>
						item.geom_cl.flatMap((item: any) => item.split("|").map((geom: string) => ({ geom, color: this.color })))
					)

					const location = result[0]?.geom?.split(",")[0].split("(")[1].split(" ")
					this.map.location({
						lon: location[0],
						lat: location[1],
					})
				}
			}
		},
		createLine() {
			if (this.map) {
				this.map.Overlays.clear()
				// @ts-ignore
				const longdo = window.longdo
				const roads = ref<Array<any>>([])
				if (this.data?.geom_list) {
					const data = this.data?.geom_list.filter((item: any) => this.barOptions().labels.includes(item.surface.name))
					data.forEach((item: any) => {
						roads.value = item.geom_cl.flatMap((child: any) =>
							child.split("|").map((geom: string) => ({
								geom,
								color: item.surface.color_code,
							}))
						)
						const lines = roads.value.map(({ geom, color }: GeomWkt) => {
							return longdo.Util.overlayFromWkt(geom, { lineColor: color })
						})

						const firstGeom = getLatLong(roads.value[0].geom)

						lines.forEach((line: Array<Object>) => {
							this.map.Overlays.add(line[0])
						})

						this.map.location({
							lon: firstGeom.lon,
							lat: firstGeom.lat,
						})
					})
				}
			}
		},
		barOptions(): any {
			const categories = ref<Array<string>>([])
			if (this.surfaceCategories) {
				categories.value = this.surfaceArray.map((id: number) => this.surfaceCategories[id]).filter(Boolean)
			}
			const option: any = {
				plotOptions: {
					pie: {
						dataLabels: {
							offset: -5,
						},
					},
				},
				title: {
					text: "ข้อมูลสรุปผิวทาง",
					align: "center",
					style: {
						fontSize: "16px",
					},
				},
				dataLabels: {
					enabled: true,
					formatter: (_: number, { seriesIndex, w }: any) => {
						return w.config.series[seriesIndex] + " กม."
					},
					style: {
						fontSize: "10px",
						colors: ["#fff"],
					},
				},
				colors: this.colors?.length === 0 ? this.surfaceColors : this.colors,
				labels: categories.value,
				tooltip: {
					enabled: true,
					y: {
						show: true,
						formatter: (value: number) => {
							if (value) {
								return value.toLocaleString() + " กม."
							} else {
								return 0 + " กม."
							}
						},
					},
				},
				legend: {
					show: false,
				},
			}
			return option
		},
		barSeries(): any {
			if (this.dataArray) {
				const result = this.surfaceArray
					.map((id: number) => this.dataArray[id])
					.filter((value: any) => value !== null && value !== undefined)
				return result.map((item: any) => Number(item))
			}
			return []
		},
	},
	getters: {},
})
