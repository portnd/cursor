import { RoadReflectiveService } from "../infrastructure/RoadReflectiveService"
import { IReflectivePostParams } from "../infrastructure/RoadReflectiveRequest"
import { ILine } from "../infrastructure/RoadReflectiveModel"
import { IFile } from "~/core/shared/types/File"

interface IState {
	loading: boolean
	lineNo: number | null
	surveyedDate: Date
	remarks: string
	csvFile: IFile
	imageFile: IFile
	lineList: ILine[]
}

export const useReflectiveCreateStore = defineStore("reflective/create", {
	state: (): IState => ({
		loading: false,
		lineNo: null,
		surveyedDate: new Date(),
		remarks: "",
		csvFile: {} as IFile,
		imageFile: {} as IFile,
		lineList: [],
	}),
	actions: {
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
		async createReflectData(roadId: number) {
			const params: IReflectivePostParams = {
				line_no: this.lineNo ? this.lineNo : null,
				surveyed_date: this.surveyedDate ? formatDate(this.surveyedDate) + " 00:00:00" : "",
				remarks: this.remarks,
				csv_file: this.csvFile.data?.file,
				csv_filename_status: this.csvFile.status,
			}

			this.loading = true

			const service = new RoadReflectiveService()
			const res = await service.createReflectiveData(roadId, params)

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		resetStore() {
			this.lineNo = null
			this.surveyedDate = new Date()
			this.remarks = ""

			this.csvFile = {} as IFile
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
		getLineOptions(state) {
			const { lineList } = state

			const options = lineList.map((line) => ({ label: line.line_no.toString(), value: line.line_no }))

			return options || []
		},
	},
})
