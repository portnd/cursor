import { IDashboardSurface, DashboardService, ISurfaceMap, IDataMartCheck } from "../infrastructure"
import { IRoadData, RoadListService } from "~/core/modules/road/roadList/infrastructure"
import { ITree } from "~/core/shared/types/Tree"

interface Location {
	lon: number
	lat: number
}

interface IPavementAge {
	id: number
	name: string
	color: string
}

interface IState {
	data: IDashboardSurface
	mapData: ISurfaceMap[]
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
	displayParam: number
	pavementAge: IPavementAge[]
	retryCreateLine: number
	syncing: boolean
	dataMart?: IDataMartCheck
}

export const useDashboardSurfaceStore = defineStore("dashboard/surface", {
	state: (): IState => ({
		data: {} as IDashboardSurface,
		mapData: [],
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
		displayParam: 1,
		pavementAge: [
			{ name: "0 - 2 ปี", color: "#50CD89" } as IPavementAge,
			{ name: "3 - 5 ปี", color: "#87C442" } as IPavementAge,
			{ name: "6 - 10 ปี", color: "#FDB833" } as IPavementAge,
			{ name: "มากกว่า 10 ปี", color: "#DC3545" } as IPavementAge,
		],
		retryCreateLine: 0,
		syncing: false,
		dataMart: undefined,
	}),
	actions: {
		async get() {
			// Loading
			this.loading = true
			console.log("get -----> loading")

			const service = new DashboardService()

			const params = ref("")
			if (this.roadId !== "") {
				params.value = `?road_id=${this.roadId}`
			}
			const res = await service.get(params.value)
			// Loading
			console.log("DashboardService -----> loading", res.data)
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				console.log("DashboardService ----->", res.data)

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

		async getMap() {
			// Loading
			this.loading = true

			const service = new DashboardService()

			const res = await service.getSurfaceMap(this.displayParam, this.roadId)
			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.mapData = res.data
				this.createLine()
				return res
			}
		},
		async syncDataMart() {
			// Loading
			this.syncing = true

			const service = new DashboardService()
			const check = await service.getDataMartCheck()
			const isNotSyncing = check.data?.stauts
			if (!isNotSyncing) {
				return
			}

			setTimeout(() => {
				this.checkSyncDataMart()
			}, 3000)
			this.$patch((state) => {
				if (state.dataMart) {
					state.dataMart.percent = 0
				}
			})
			const res = await service.getDataMart()

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		async checkSyncDataMart() {
			const service = new DashboardService()

			const res = await service.getDataMartCheck()

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
				return false
			} else {
				if (res.data.stauts) {
					this.get()
					this.getRoad()
					this.getMap()
					setTimeout(() => {
						this.syncing = false
					}, 1000)
				} else {
					this.syncing = true
					setTimeout(() => {
						this.checkSyncDataMart()
					}, 5000)
				}
				this.$patch((state) => {
					state.dataMart = res.data
				})

				return !!res.data.stauts
			}
		},
		createLine() {
			if (this.map) {
				this.map.Overlays.clear()
				// @ts-ignore
				const longdo = window.longdo

				if (this.mapData.length > 0) {
					console.log("this.retryCreateLine 0 =", this.retryCreateLine)

					this.mapData.forEach((item) => {
						const strLine = longdo.Util.overlayFromGeoJSON(item.the_geom, {
							lineColor: item.color === "#000000" ? convertHexToRGBA(item.color, 0.5) : convertHexToRGBA(item.color, 1),
							detail: this.generatePopupDetails(item),
						})
						console.log("%c line color", `background: ${convertHexToRGBA(item.color, 0.5)}; color: ${item.color}`)
						this.map.Overlays.add(strLine[0])
					})

					const geom =
						this.mapData[0].the_geom.coordinates.length > 0 ? this.mapData[0].the_geom.coordinates[0] : undefined

					if (geom) {
						this.map.location({
							lon: geom[0],
							lat: geom[1],
						})
					}

					this.map.zoom(12)
					this.retryCreateLine = 0
				}
			} else {
				console.log("this.retryCreateLine =", this.retryCreateLine)
				if (this.retryCreateLine <= 5) {
					setTimeout(() => {
						this.createLine()
					}, 500)
				}

				this.retryCreateLine++
			}
		},
		getKmStartEnd(road: ISurfaceMap) {
			return `${convertMeterToKm(road.km_start)} - ${convertMeterToKm(road.km_end)}`
		},
		getKmTotal(road: ISurfaceMap) {
			const result = Math.abs(road.km_end - road.km_start)
			return result / 1000
		},
		getProcessHeaderTitle(item: ISurfaceMap) {
			const icon = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
			<g clip-path="url(#clip0_10409_34188)">
			<path d="M20.2759 2.48C20.1536 1.78514 19.7904 1.15561 19.2499 0.702047C18.7095 0.248483 18.0265 -9.24445e-05 17.3209 4.06091e-06H6.68094C5.97464 -0.0011587 5.29067 0.247414 4.74993 0.701785C4.20918 1.15616 3.84648 1.78707 3.72594 2.483L-0.0390625 24H24.0609L20.2759 2.48ZM12.9999 21H10.9999V17H12.9999V21ZM12.9999 14H10.9999V10H12.9999V14ZM12.9999 7.00001H10.9999V3H12.9999V7.00001Z" fill="${item.color}"/>
			</g>
			<defs>
			<clipPath id="clip0_10409_34188">
			<rect width="24" height="24" fill="white"/>
			</clipPath>
			</defs>
			</svg>`

			return `<div class="col-1">
			${icon}
			</div>
			<div class="col"
			style="
			color: ${item.color};
			font-size: 16px;
			font-weight: 500;
			">${item.title}</div>`
		},
		generatePopupDetails(item: ISurfaceMap) {
			return `
					<div class="row mb-3" style="
						width: 220px;" >
						${this.getProcessHeaderTitle(item)}
				 	</div>
					<div class="row">
						<div class="col-6">สายทาง:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${item.road_group_name}</div>
					</div>
					<div class="row mb-2">
						<div class="col-6">เลขที่สัญญา:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${item.contract_number}</div>
					</div>

					<div class="row mb-2">
						<div class="col-6">ปีงบประมาณ:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${item.year && isNumber(item.year) ? Number(item.year) + 543 : ""}</div>
					</div>
					<div class="row mb-2">
						<div class="col-6">วันที่ซ่อมล่าสุด:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${buddhistFormatDate(item.last_inspection_date, "dd mmm yyyy")} </div>
					</div>
					<div class="row mb-2 mt-2">
						<div class="col-6">จาก - ถึง:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${item.road_name} </div>
					</div>
					<div class="row mb-2 mt-2">
						<div class="col-6">กม. เริ่มต้น - กม. สิ้นสุด:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${this.getKmStartEnd(item)} </div>
					</div>
					<div class="row mb-2 mt-2">
						<div class="col-6">ระยะทาง:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${this.getKmTotal(item)} กม.</div>
					</div>
					<div class="row mb-2 mt-2">
						<div class="col-6">ผิวทางล่าสุด:</div>
						<div class="col" style="
						color: #5E6278;
						font-size: 12px;
						font-weight: 400;
						">${item.surface_name} </div>
					</div>`
		},
	},
	getters: {},
})
