import { RoadTitleService } from "../infrastructure/index"
import { IRoad, IRoadRefDepot } from "../../../road/info/summary/infrastructure/RoadDetailModel.d"

interface IState {
	data: IRoad
	loading: boolean
	tab: string
}

export const useRoadTitleStore = defineStore("common/road", {
	state: (): IState => ({
		data: { ref_depot: { id: 0 } as IRoadRefDepot } as IRoad,
		loading: false,
		tab: "",
	}),
	actions: {
		async getData(id: number) {
			this.loading = true

			const service = new RoadTitleService()
			const res = await service.getRoad(id)
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.data = res.data
			}
		},
	},
	getters: {
		getSubTitle(state) {
			const item = state?.data
			const kmStart = item.road_info?.km_start
			const kmEnd = item.road_info?.km_end
			return `${item.road_code} | ${item.responsible_code} | กม.ที่ ${convertMeterToKm(kmStart)} - ${convertMeterToKm(
				kmEnd
			)}`
		},
		getSumDistance(state) {
			const item = state?.data
			const kmStart = item.road_info?.km_start
			const kmEnd = item.road_info?.km_end
			const result = Math.abs((kmEnd - kmStart) / 1000).toFixed(3)
			return result
		},
	},
})
