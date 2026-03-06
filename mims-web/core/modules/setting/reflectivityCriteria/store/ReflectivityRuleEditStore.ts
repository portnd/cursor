import {
	IRequestConditionReflectivityRule,
	IRequestReflectivityRule,
	ReflectivityRuleService,
	IRoadLine,
} from "../infrastructure"
import { IReflectivityRuleList } from "~/composables/useReflectivityRule"
import { useInitDataStore } from "~/core/modules/initData/store/InitDataStore"

const rule = useReflectivityRule()

export const useReflectivityRuleEditStore = defineStore("setting/ReflectivityRule/edit", {
	state: () => ({
		id: 0,
		name: "",
		rule,
		reflectivityRangeId: 0,
		roadLine: {} as IRoadLine,
		loading: false,
	}),
	actions: {
		async get(id: number) {
			this.id = id
			// Loading
			this.loading = true

			const service = new ReflectivityRuleService()
			const res = await service.get(this.id)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.name = res.data.owner_name
				this.reflectivityRangeId = res.data.ref_reflectivity_range_id
				this.roadLine = res.data.road_line
				this.buildDataReflexRule()
				return res
			}
		},
		async edit() {
			// Loading
			this.loading = true

			const result: IRequestConditionReflectivityRule[] = []
			rule.forEach((item) => {
				for (let i = 0; i < item.white.reflectivity_list.length; i++) {
					const whiteCondition = item.white.reflectivity_list[i]
					const yellowCondition = item.yellow.reflectivity_list[i]

					result.push({
						grade_id: whiteCondition.grade_id,
						left_value_yellow: Number(yellowCondition.left_value) ?? 0,
						right_value_yellow: Number(yellowCondition.right_value) ?? 0,
						left_value_white: Number(whiteCondition.left_value) ?? 0,
						right_value_white: Number(whiteCondition.right_value) ?? 0,
					})
				}
			})
			const params: IRequestReflectivityRule = {
				name: this.name,
				ref_reflectivity_range_id: this.reflectivityRangeId,
				road_line_list: result,
			}

			const surveyRuleService = new ReflectivityRuleService()
			const res = await surveyRuleService.put(this.id, params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		// update layout
		buildDataReflexRule() {
			rule.forEach((data) => {
				data.white.reflectivity_list.forEach((cond: IReflectivityRuleList) => {
					const whiteCond = this.roadLine.white.find((c) => c.grade.id === cond.grade_id)
					if (whiteCond) {
						cond.left_value = whiteCond.left_value_white
						cond.right_value = whiteCond.right_value_white
					}
				})

				data.yellow.reflectivity_list.forEach((cond: IReflectivityRuleList) => {
					const yellowCond = this.roadLine.yellow.find((c) => c.grade.id === cond.grade_id)
					if (yellowCond) {
						cond.left_value = yellowCond.left_value_yellow
						cond.right_value = yellowCond.right_value_yellow
					}
				})
			})
		},
	},
	getters: {
		optionGenerator() {
			const initDataStore = useInitDataStore()
			const list = toOptions(initDataStore.data?.ref_reflectivity_range)
			return list
		},
	},
})
