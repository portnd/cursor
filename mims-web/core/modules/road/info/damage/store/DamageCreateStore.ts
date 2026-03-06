import { RoadDamageService } from "../infrastructure"
import { IParams } from "../infrastructure/RoadDamageRequest"
import { IFile } from "~~/core/shared/types/File"

interface IState {
	roadId: number
	date: Date
	lane: string
	loading: boolean
	itemData: IFile
	image: IFile
}

export const useCreateStore = defineStore("road/damage/create", {
	state: (): IState => ({
		roadId: 0,
		date: new Date(),
		lane: "",
		loading: false,
		itemData: {} as IFile,
		image: {} as IFile,
	}),
	actions: {
		async createDamageData(roadId: number) {
			// Loading
			this.loading = true

			const params: IParams = {
				lane_no: this.lane ? this.lane : "",
				surveyed_date: this.date ? formatDate(this.date) + " 00:00:00" : "",
				damage_filename: this.itemData.data?.file,
				damage_filename_status: this.itemData.status,
				image_filename: this.image.data?.file,
				image_filename_status: this.image.status,
			}

			const service = new RoadDamageService()
			const res = await service.createDamageData(roadId, params)

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
