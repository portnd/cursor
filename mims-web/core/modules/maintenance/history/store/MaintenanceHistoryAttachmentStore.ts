import {
	IMaintenanceHistoryAttrachments,
	IMaintenanceHistoryDetailData,
} from "../infrastructure/MaintenanceHistoryModel"
import { IMaintenanceHistoryFileParams } from "../infrastructure/MaintenanceHistoryRequest"
import { MaintenanceHistoryService } from "../infrastructure/MaintenanceHistoryService"

interface IRoadGeom {
	geom: string
	color: string
}

interface IState {
	loading: boolean
	projectDetais: IMaintenanceHistoryDetailData
	roadGeom: IRoadGeom[]
	attachmentData: IMaintenanceHistoryAttrachments[]
	params: IMaintenanceHistoryFileParams
	map: any
	longdo: any
}

export const useMaintenanceHistoryAttachmentStore = defineStore("maintenance/history/attachment", {
	state: (): IState => ({
		loading: false,
		projectDetais: {} as IMaintenanceHistoryDetailData,
		roadGeom: [],
		attachmentData: [],
		params: {
			file_type: "",
			order: "date",
		},
		map: null,
		longdo: null,
	}),
	actions: {
		async getMaintenanceHistoryFile(id: number) {
			this.loading = true

			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceAttrachment(id, this.params)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.attachmentData = res.data
				await this.getMaintenanceHistoryDetail(id)
			}

			this.loading = false
		},
		async getMaintenanceHistoryDetail(id: number) {
			const service = new MaintenanceHistoryService()
			const res = await service.getMaintenanceHistoryDetails(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.projectDetais = res.data
				this.setRoadGeom()
			}
		},
		setMap(map: any) {
			this.map = map
			if (this.map) {
				// @ts-ignore
				this.longdo = window.longdo

				this.createLine()
			}
		},
		setRoadGeom() {
			if (Object.keys(this.projectDetais).length > 0) {
				const data = this.projectDetais

				const maintenanceHistoryGeom = data.maintenance_road_histories.map((item) => ({
					geom: item.the_geom,
					color: item.road_info.road_color_code,
				}))

				const maintenanceRoadsGeom = data.maintenance_roads.map((item) => ({
					geom: item.the_geom,
					color: item.road_info.road_color_code,
				}))

				this.roadGeom = maintenanceRoadsGeom.concat(maintenanceHistoryGeom)
			}
		},
		createLine() {
			if (this.map || this.roadGeom.length > 0) {
				const polylines = this.roadGeom.flatMap((item) => {
					return this.longdo.Util.overlayFromWkt(item.geom, { lineColor: item.color })
				})

				polylines.forEach((line) => {
					this.map.Overlays.add(line)
				})

				const latLon = getLatLong(this.roadGeom[0]?.geom)
				this.map.location({
					lon: latLon.lon,
					lat: latLon.lat,
				})
				this.map.zoom(13)
			}
		},
	},
	getters: {},
})
