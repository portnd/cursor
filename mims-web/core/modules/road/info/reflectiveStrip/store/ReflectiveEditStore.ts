import { RoadReflectiveService, IReflectiveData, ILine, IReflectiveUpdateParams } from "../infrastructure"
import { IFile } from "~/core/shared/types/File"

interface IState {
	loading: boolean
	lineNo: number
	surveyedDate: Date
	remarks: string
	idParent: number
	csvFile: IFile
	csv_file_path: string
	imageFile: IFile
	image_file_path: string
	conditionData: {
		data: IReflectiveData
		loading: boolean
	}
	lineList: ILine[]
}

export const useReflectiveEditStore = defineStore("reflective/update", {
	state: (): IState => ({
		loading: false,
		lineNo: 0,
		surveyedDate: new Date(),
		remarks: "",
		idParent: 0,
		csvFile: {} as IFile,
		csv_file_path: "",
		imageFile: {} as IFile,
		image_file_path: "",
		conditionData: {
			data: {} as IReflectiveData,
			loading: false,
		},
		lineList: [],
	}),
	actions: {
		async updateReflectData(roadId: number) {
			this.loading = true

			const params: IReflectiveUpdateParams = {
				line_no: this.lineNo ? this.lineNo : null,
				surveyed_date: this.surveyedDate ? formatDate(this.surveyedDate) + " 00:00:00" : "",
				remarks: this.remarks,
				id_parent: this.idParent,
				csv_file: this.csvFile.data?.file,
				csv_filename_status: this.csvFile.status,
			}

			const service = new RoadReflectiveService()
			const res = await service.updateReflectiveData(roadId, params)

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		async getDefaultData(idParent: number) {
			this.conditionData.loading = true
			const service = new RoadReflectiveService()
			const res = await service.getReflectivityData(idParent)

			this.conditionData.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.conditionData.data = res.data
				this.setDefaultEdit(this.conditionData.data)
			}
		},
		setDefaultEdit(data: IReflectiveData) {
			if (data) {
				this.surveyedDate = new Date(data.surveyed_date)
				this.lineNo = data.line_no
				this.csv_file_path = data.csv_file
				this.image_file_path = data.img_filepath
				this.idParent = data.id_parent
				this.remarks = data.remarks
			}
		},
		async getLineList(roadId: number) {
			this.loading = true

			const service = new RoadReflectiveService()
			const res = await service.getLaneList(roadId)

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.lineList = res.data
			}
		},
	},
	getters: {
		getLineOptions(state) {
			const { lineList } = state

			const options = lineList.map((line) => ({ label: line.line_no.toString(), value: line.line_no }))

			return options || []
		},
	},
})
