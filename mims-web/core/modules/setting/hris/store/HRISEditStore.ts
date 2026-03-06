import { IRequestHRIS, HRISService, IHRISItem } from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	id: number
	params: IRequestHRIS
	data: IHRISItem
	loading: boolean
	status: number
}

export const useHRISEditStore = defineStore("setting/HRIS/edit", {
	state: (): IState => ({
		id: 0,
		data: {} as IHRISItem,
		params: {
			office_of_highways_code: null,
			road_number: null,
			section_road_number: null,
			status: null,
		},
		loading: false,
		status: 0,
	}),
	actions: {
		async get(id: number) {
			// ล้างค่า
			this.reset()

			this.id = id

			// Loading
			this.loading = true

			const service = new HRISService()
			const res = await service.get(this.id)

			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.data = res.data
				this.status = this.data.status ? 1 : 0
				return res
			}
		},
		async edit() {
			// Loading
			this.loading = true

			this.params.office_of_highways_code = this.data.office_of_highways_code
			this.params.road_number = this.data.road_number
			this.params.section_road_number = this.data.section_road_number

			this.params.status = this.status !== 0

			const service = new HRISService()
			const res = await service.put(this.id, this.params)

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
