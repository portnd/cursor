import { IRequestConditionReflectivityRule, IRequestReflectivityRule, ReflectivityRuleService } from "../infrastructure"
import { useInitDataStore } from "~/core/modules/initData/store/InitDataStore"
import { useReflectivityRule } from "~/composables/useReflectivityRule"

const rule = useReflectivityRule()

export const useReflectivityRuleCreateStore = defineStore("setting/ReflectivityRule/create", {
	state: () => ({
		name: "",
		rule,
		reflectivityRangeId: 1,
		loading: false,
	}),
	actions: {
		async create() {
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

			const service = new ReflectivityRuleService()
			const res = await service.post(params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
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
