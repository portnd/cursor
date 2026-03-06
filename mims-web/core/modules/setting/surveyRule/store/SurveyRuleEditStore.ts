import { SurveyRuleService, IRequestSurveyRule, IRequestConditionSurveyRule, IResSurveyRule } from "../infrastructure"
import { ISurveyRule } from "~/composables/useSurveyRule"

interface IState {
	id: number
	name: string
	data: IResSurveyRule
	conditionRangeId: number
	survey: ISurveyRule[]
	loading: boolean
	originRangeId: number
}

export const useSurveyRuleEditStore = defineStore("setting/survey-rules/edit", {
	state: (): IState => ({
		id: -1,
		name: "",
		data: {} as IResSurveyRule,
		conditionRangeId: 1,
		survey: [],
		loading: false,
		originRangeId: 1,
	}),
	actions: {
		async get(id: number) {
			// console.log("get")
			if (this.id !== -1) {
				// this.id !== -1 คือให้ดึง api ครั้งแรกครั้งเดียว เพราะการเลือกประเภทจะทำให้ทั้งหน้า re-render ใหม่ทั้งหทด onMounted จะถูกเรียกทุกครั้ง
				return
			}
			this.id = id
			// Loading
			this.loading = true

			const surveyRuleService = new SurveyRuleService()
			const res = await surveyRuleService.get(this.id)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.name = res.data.owner_name
				this.conditionRangeId = res.data.ref_condition_range_id
				this.originRangeId = res.data.ref_condition_range_id
				this.data = res.data
				this.updateSurveyRul()

				return res
			}
		},
		async edit() {
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
			const res = await surveyRuleService.put(this.id, params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		buildDataSurveyRule() {
			if (this.data.condition_list) {
				this.data.condition_list.forEach((condition) => {
					// หา index ของข้อมูลใน initdata ที่ตรงกับ condition_type
					const index = this.survey.findIndex((data) => data.name === condition.condition_type)

					if (index !== -1) {
						// Loop ผ่านแต่ละ condition ใน ac และ cc เพื่ออัพเดทข้อมูล
						condition.surface_type.ac.forEach((acCondition) => {
							const acIndex = this.survey[index].ac.conditionList.findIndex(
								(item) => item.grade_id === acCondition.grade.id
							)
							if (acIndex !== -1) {
								this.survey[index].ac.conditionList[acIndex].left_value = acCondition.left_value_ac
								this.survey[index].ac.conditionList[acIndex].right_value = acCondition.right_value_ac
							}
						})

						condition.surface_type.cc.forEach((ccCondition) => {
							const ccIndex = this.survey[index].cc.conditionList.findIndex(
								(item) => item.grade_id === ccCondition.grade.id
							)
							if (ccIndex !== -1) {
								this.survey[index].cc.conditionList[ccIndex].left_value = ccCondition.left_value_cc
								this.survey[index].cc.conditionList[ccIndex].right_value = ccCondition.right_value_cc
							}
						})
					}
				})
			}
		},
		updateSurveyRul() {
			this.survey = useSurveyRule(this.conditionRangeId)
			if (this.conditionRangeId === this.originRangeId) {
				this.buildDataSurveyRule()
			}
		},
	},
	getters: {
		generativeDataSurveyRule(state) {
			switch (state.conditionRangeId) {
				case 1:
					return state.survey
				case 2:
					return state.survey.filter((e) => e.name !== "MPD")

				case 3:
					return state.survey.filter((e) => e.name === "IRI")

				default:
					break
			}
		},
		getSurveyConditionType(state) {
			const { survey } = state
			return survey.map((item) => item.name) ?? []
		},
	},
})
