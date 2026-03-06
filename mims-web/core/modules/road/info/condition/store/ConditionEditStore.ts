import { RoadConditionService, IConditionPutParams, IConditionData, ILane } from "../infrastructure"
import { IFile } from "~/core/shared/types/File"

interface IState {
	loading: boolean
	laneNo: number
	surveyedDate: Date
	remarks: string
	idParent: number
	iriFile: IFile
	iri_file_path: string
	imageFile: IFile
	image_file_path: string
	conditionData: {
		data: IConditionData
		loading: boolean
	}
	laneList: ILane[]
}

export const useConditionEditStore = defineStore("condition/edit", {
	state: (): IState => ({
		loading: false,
		laneNo: 0,
		surveyedDate: new Date(),
		remarks: "",
		idParent: 0,
		iriFile: {} as IFile,
		iri_file_path: "",
		imageFile: {} as IFile,
		image_file_path: "",
		conditionData: {
			data: {} as IConditionData,
			loading: false,
		},
		laneList: [],
	}),
	actions: {
		async updateCondition(roadId: number) {
			this.loading = true

			const params: IConditionPutParams = {
				lane_no: this.laneNo,
				surveyed_date: this.surveyedDate ? formatDate(this.surveyedDate) + " 00:00:00" : "",
				remarks: this.remarks,
				id_parent: this.idParent,
				iri_filename: this.iriFile.data?.file,
				iri_filename_status: this.iriFile.status,
				image_filename: this.imageFile.data?.file,
				image_filename_status: this.imageFile.status,
			}

			const service = new RoadConditionService()
			const res = await service.updateConditionData(roadId, params)

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		async getDefaultData(idParent: number) {
			this.conditionData.loading = true
			const service = new RoadConditionService()
			const res = await service.getConditionData(idParent)

			this.conditionData.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.conditionData.data = res.data
				this.setDefaultEdit(this.conditionData.data)
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
		setDefaultEdit(data: IConditionData) {
			if (data) {
				this.surveyedDate = new Date(data.surveyed_date)
				this.laneNo = data.lane_no
				this.iri_file_path = data.iri_filename
				this.image_file_path = data.img_filepath
				this.idParent = data.id_parent
				this.remarks = data.remarks
			}
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
