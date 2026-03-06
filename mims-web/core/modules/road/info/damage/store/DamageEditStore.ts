import { RoadDamageService } from "../infrastructure/RoadDamageService"
import { IDataImport } from "../infrastructure/RoadDamageModel"
import { IParams } from "../infrastructure/RoadDamageRequest"
import { IFile } from "~~/core/shared/types/File"

interface IState {
	roadId: number
	id: number
	parentId: number

	loading: boolean
	data: IDataImport

	date: Date
	lane_no: string
	csv_data: IFile
	image: IFile
}

export const useEditStore = defineStore("road/damage/edit", {
	state: (): IState => ({
		roadId: 0,
		id: 0,
		parentId: 0,
		loading: false,
		data: {} as IDataImport,
		date: new Date(),
		lane_no: "",
		csv_data: {} as IFile,
		image: {} as IFile,
	}),
	actions: {
		async getDamageDefault(roadId: number, id: number) {
			// ล้างค่า
			this.$reset()

			this.roadId = roadId
			this.id = id

			// Loading
			this.loading = true

			const service = new RoadDamageService()
			const res = await service.getDamageDefault(roadId, this.id)

			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.data = res.data
				this.date = new Date(this.data.surveyed_date)
				this.lane_no = this.data.lane_no.toString()
				this.csv_data = {} as IFile
				this.image = {} as IFile
			}
		},
		async updateDamageData(roadId: number) {
			this.loading = true
			const params: IParams = {
				lane_no: this.lane_no ? this.lane_no.toString() : "",
				surveyed_date: this.date ? formatDate(this.date) + " 00:00:00" : "",
				damage_filename: this.csv_data.data?.file,
				damage_filename_status: this.csv_data.status,
				image_filename: this.image.data?.file,
				image_filename_status: this.image.status,
			}

			const service = new RoadDamageService()
			const res = await service.updateDamageData(roadId, this.parentId, params)

			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				// ล้างค่า
				this.$reset()

				return res
			}
		},
	},
	getters: {},
})
