import { IRequestHRIS, HRISService } from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	params: IRequestHRIS
	loading: boolean
	status: number
}

export const useHRISCreateStore = defineStore("setting/hris/create", {
	state: (): IState => ({
		params: {
			office_of_highways_code: null,
			road_number: null,
			section_road_number: null,
			status: null,
		},
		loading: false,
		status: 1,
	}),
	actions: {
		async create() {
			// Loading
			this.loading = true

			this.params.status = this.status !== 0
			const service = new HRISService()
			const res = await service.post(this.params)

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
			this.params = {
				office_of_highways_code: null,
				road_number: null,
				section_road_number: null,
				status: null,
			}
		},
	},
	getters: {
		getStatusOption() {
			return [{ value: 1, label: "Active" } as IOption, { value: 0, label: "Inactive" } as IOption]
		},
	},
})
