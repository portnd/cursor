import { IRequestSign, ISign, SignService } from "../infrastructure"
import { IFile } from "~~/core/shared/types/File"

interface IState {
	id: number
	data: ISign
	image: IFile
	loading: boolean
}

export const useSignEditStore = defineStore("setting/sign/edit", {
	state: (): IState => ({
		id: 0,
		data: {
			sign_image_filepath: "",
			id: 0,
			name: "",
			abbr: "",
		},
		image: {} as IFile,
		loading: false,
	}),
	actions: {
		async get(id: number) {
			// ล้างค่า
			this.reset()

			this.id = id

			// Loading
			this.loading = true

			const signService = new SignService()
			const res = await signService.get(this.id)

			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.data = res.data
				return res
			}
		},
		async edit() {
			// Loading
			this.loading = true

			const params: IRequestSign = {
				name: this.data.name,
				abbr: this.data.abbr,
				image: this.image.data?.file,
				image_status: this.image.status,
			}

			const signService = new SignService()
			const res = await signService.put(this.id, params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				// ล้างค่า
				this.reset()

				return res
			}
		},
		reset() {
			this.data.id = 0
			this.data.name = ""
			this.data.abbr = ""
			this.data.sign_image_filepath = ""
			this.image = {} as IFile
		},
	},
	getters: {},
})
