import { IRequestSign, SignService } from "../infrastructure"
import { IFile } from "~~/core/shared/types/File"

interface IState {
	name: string
	abbr: string
	image: IFile
	loading: boolean
}

export const useSignStore = defineStore("setting/sign/create", {
	state: (): IState => ({
		name: "",
		abbr: "",
		image: {} as IFile,
		loading: false,
	}),
	actions: {
		async create() {
			// Loading
			this.loading = true

			const params: IRequestSign = {
				name: this.name,
				abbr: this.abbr,
				image: this.image.data?.file,
				image_status: this.image.status,
			}

			const signService = new SignService()
			const res = await signService.post(params)

			// Loading
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
