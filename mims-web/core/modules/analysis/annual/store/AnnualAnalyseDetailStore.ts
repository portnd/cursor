import {
	AnnualService,
	IAnnualAnalyzeData,
	IAnnualInterventionCriteria,
	IAnnualRoadGroup,
	IAnnualUpdateModelReq,
} from "../infrastructure"
import { ITree } from "~/core/shared/types/Tree"

interface IStateParams {
	road_id: string[]
	intervention_criteria_id: string[] | undefined
}

interface IState {
	loading: boolean
	submit_loading: boolean
	search_loading: boolean
	roadGroup: IAnnualRoadGroup[]
	intervention_criteria: IAnnualInterventionCriteria[]
	data: IAnnualAnalyzeData
	params: IStateParams
	update_params: IAnnualUpdateModelReq[]
}

export const useAnnualAnalyseDetailStore = defineStore("annual-analyse/detail", {
	state: (): IState => ({
		loading: false,
		search_loading: false,
		submit_loading: false,
		roadGroup: [],
		intervention_criteria: [],
		params: {
			road_id: [],
			intervention_criteria_id: undefined,
		},
		update_params: [],
		data: {} as IAnnualAnalyzeData,
	}),
	actions: {
		async getRoadsTree() {
			const service = new AnnualService()
			const res = await service.getRoadTree()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.roadGroup = res.data
			}
		},
		async getStrategicList() {
			const service = new AnnualService()
			const res = await service.getInterventionCriteria()

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.intervention_criteria = res.data
			}
		},
		async getDetails(id: number) {
			const service = new AnnualService()
			const res = await service.getAnalyzeDefaultDetails(id)

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.data = res.data

				this.params.road_id = res.data.roads.map((road) => String(road.road_id))
				// this.params.intervention_criteria_id = res.data
			}
		},
		async updateModel(id: number) {
			this.submit_loading = true

			const service = new AnnualService()
			const res = await service.updateModel(id, this.update_params)

			this.submit_loading = false
			if (!res.status) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		handleDataTable(data: any) {
			this.update_params = data.map((item: any) => ({
				id: item.id,
				intervention_criteria_id: item.intervention_criteria_id,
			}))
		},
	},
	getters: {
		getRoadGroupOptions(state) {
			const { roadGroup } = state
			const options: ITree[] = roadGroup?.map((road) => ({
				label: road.short_name,
				id: String(road.id),
				children: road.road_sections?.map((section) => ({
					label: `${section.name_origin} - ${section.name_destination}`,
					id: String(section.id),
					children: section.roads?.map((road) => ({ label: road.name, id: String(road.id) })),
				})),
			}))

			return options || []
		},
		getLaneOptions(state) {
			const { roadGroup, params } = state

			if (params.road_id.length === 0) {
				return []
			}

			const roadIds = new Set(params.road_id.map(Number))

			const lanes = roadGroup.flatMap((parent) =>
				parent.road_sections.flatMap(
					(section) => section.roads?.filter((road) => roadIds.has(road.id)).map((road) => road.lane_total) || []
				)
			)

			const validLanes = lanes.filter((lane) => lane != null)
			const maxLanes = validLanes.length > 0 ? Math.max(...validLanes) : 0

			const options = [
				{ label: "ทั้งหมด", value: 0 },
				...Array.from({ length: maxLanes }, (_, index) => ({ label: `${index + 1}`, value: index + 1 })),
			]

			return options
		},
		getInterventionOptions(state) {
			const intervention = state.intervention_criteria

			const options: ITree[] = intervention.map((parent) => ({
				label: parent.label,
				id: `${parent.id}`,
				children: parent.children.map((child) => ({
					label: child.label,
					id: `${child.id}`,
				})),
			}))

			return options || []
		},
	},
})
