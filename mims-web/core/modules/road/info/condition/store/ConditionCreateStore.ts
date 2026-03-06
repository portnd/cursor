import { RoadConditionService } from "../infrastructure/RoadConditionService"
import { IConditionPostParams } from "../infrastructure/RoadConditionRequest"
import { ILane } from "../infrastructure/RoadConditionModel"
import { IFile } from "~/core/shared/types/File"

interface IState {
	loading: boolean
	laneNo: string
	surveyedDate: Date
	remarks: string
	iriFile: IFile
	imageFile: IFile
	laneList: ILane[]
}

export const useConditionCreateStore = defineStore("condition/create", {
	state: (): IState => ({
		loading: false,
		laneNo: "",
		surveyedDate: new Date(),
		remarks: "",
		iriFile: {} as IFile,
		imageFile: {} as IFile,
		laneList: [],
	}),
	actions: {
		async postConditions(roadId: number) {
			const params: IConditionPostParams = {
				lane_no: this.laneNo ? this.laneNo : "",
				surveyed_date: this.surveyedDate ? formatDate(this.surveyedDate) + " 00:00:00" : "",
				remarks: this.remarks,
				iri_filename: this.iriFile.data?.file,
				iri_filename_status: this.iriFile.status,
				image_filename: this.imageFile.data?.file,
				image_filename_status: this.iriFile.status,
			}

			this.loading = true

			const service = new RoadConditionService()
			const res = await service.createConditionData(roadId, params)

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		async getLaneList(roadId: number) {
			const service = new RoadConditionService()
			const res = await service.getLaneList(roadId)

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.laneList = res.data
			}
		},
		resetStore() {
			this.laneNo = ""
			this.surveyedDate = new Date()
			this.remarks = ""

			this.iriFile = {} as IFile
			this.imageFile = {} as IFile
		},
		formatDates(inputDate: Date) {
			const date = new Date(inputDate)
			const year = date.getFullYear()
			const month = String(date.getMonth() + 1).padStart(2, "0")
			const day = String(date.getDate()).padStart(2, "0")
			const hours = String(date.getHours()).padStart(2, "0")
			const minutes = String(date.getMinutes()).padStart(2, "0")
			const seconds = String(date.getSeconds()).padStart(2, "0")

			return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
		},
	},
	getters: {
		getLaneListOptions(state) {
			const { laneList } = state

			const options = laneList.map((lane) => ({ label: lane.lane_no.toString(), value: lane.lane_no }))

			return options ?? []
		},
	},
})
