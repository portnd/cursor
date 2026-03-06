import { IRequestSurveyRule, IRequestConditionSurveyRule, SurveyRuleService } from "../infrastructure"

export const useSurveyRuleCreateStore = defineStore("setting/survey-rules/create", {
	state: () => ({
		name: "",
		conditionRangeId: 1,
		survey: useSurveyRule(1),
		loading: false,
	}),
	actions: {
		async create() {
			// Loading
			this.loading = true

			const result: IRequestConditionSurveyRule[] = []

			this.survey.forEach((item) => {
				for (let i = 0; i < item.ac.conditionList.length; i++) {
					const acCondition = item.ac.conditionList[i]
					const ccCondition = item.cc.conditionList[i]

					result.push({
						grade_id: acCondition.grade_id,
						left_value_ac: Number(acCondition.left_value) ?? 0,
						right_value_ac: Number(acCondition.right_value) ?? 0,
						left_value_cc: Number(ccCondition.left_value) ?? 0,
						right_value_cc: Number(ccCondition.right_value) ?? 0,
						condition_type: item.name,
					})
				}
			})

			const params: IRequestSurveyRule = {
				name: this.name,
				ref_condition_range_id: this.conditionRangeId,
				condition_list: result,
			}

			const surveyRuleService = new SurveyRuleService()
			const res = await surveyRuleService.post(params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		updateRangeId() {
			this.survey = useSurveyRule(this.conditionRangeId)
		},
	},
	getters: {},
})
