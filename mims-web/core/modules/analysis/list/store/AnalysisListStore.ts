import {
	AnalysisListService,
	IAnalysisParams,
	IStrategicsList,
	IMaintenanceAnalysis,
	IMaintenanceAnalysisCondition,
} from "../infrastructure"
import { IOption } from "~/core/shared/types/Option"

interface IState {
	loading: boolean
	params: IAnalysisParams
	strategics: IStrategicsList[]
}

export const useAnalysisListStore = defineStore("analysis/list", {
	state: (): IState => ({
		loading: false,
		params: {
			type: null,
			condition: "",
		},
		strategics: [],
	}),
	actions: {
		async getStrategics() {
			this.loading = true
			const service = new AnalysisListService()
			const res = await service.getStrategics()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.strategics = res.data
			}
		},
		toDetailsPage(item: IMaintenanceAnalysis) {
			const analysisMapping: Record<number, string> = {
				1: "strategic",
				2: "annual",
			}

			const conditionMapping: Record<number, string> = {
				1: "no-budget-limit",
				2: "budget-limit",
				3: "iri-target",
			}

			const analysisId = this.strategics.find((e) => e.name === item.type_analysis)?.id
			const basePath = analysisMapping[analysisId as keyof typeof analysisMapping]
			const conditionPath = conditionMapping[item.maintenance_condition_type_id]

			if (item.status === "สำเร็จ" && (analysisId === 1 || analysisId === 2) && conditionPath) {
				navigateTo(`analyses/${basePath}/summary/${item.id}/${conditionPath}`)
			}
		},
		generateCondition(item: IMaintenanceAnalysisCondition) {
			// const roadName = item.road_name.split(",").join(", ")
			const filter = item.filter
			const lane = item.lane
			const group = item.km_group
			const discount = item.discount
			const condition = item.condition
			const target = item.target
			const conditionFilter = item.condition_filter !== "" ? item.condition_filter : ""

			if (conditionFilter === "") {
				return `<span class="fw-semibold">ตัวกรอง:</span> ${filter} , <span class="fw-semibold">ช่องจราจร:</span> ${lane}, <span class="fw-semibold">จัดกลุ่ม:</span> ${group}
      <span class="fw-semibold">อัตราคิดลด (Discount Rate):</span> ${
				discount == null ? "-" : discount + "%"
			}, <span class="fw-semibold">เงื่อนไข:</span> ${condition}, <span class="fw-semibold">เป้าหมาย:</span> ${target}`
			}

			return `<span class="fw-semibold">ตัวกรอง:</span> ${filter} , <span class="fw-semibold">ช่องจราจร:</span> ${lane}, <span class="fw-semibold">จัดกลุ่ม:</span> ${group}
     ${conditionFilter}
      <span class="fw-semibold">อัตราคิดลด (Discount Rate):</span> ${
				discount == null ? "-" : discount + "%"
			}, <span class="fw-semibold">เงื่อนไข:</span> ${condition}, ${
				target && target !== "" ? `<span class="fw-semibold">เป้าหมาย:</span> ${target}` : ""
			}`
		},
		generateComment(comment: string) {
			if (!comment || comment === "-") {
				return "-"
			}

			const maxLength = 250

			const formattedComment = comment.replace(/\n/g, "<br>")

			if (formattedComment.length > maxLength) {
				return formattedComment.substring(0, maxLength) + "..."
			}

			return formattedComment
		},
		generateStatusColors(status: string) {
			if (status === "สำเร็จ") {
				return "badge-light-success text-success" + this.togglerCursor(status)
			} else if (status.includes("กำลังดำเนินการ")) {
				return "badge-light-warning text-warning" + this.togglerCursor(status)
			} else if (status === "เกิดข้อผิดพลาด") {
				return "badge-light-danger text-danger" + this.togglerCursor(status)
			}
		},
		resetParams() {
			this.params.condition = ""
			this.params.type = null
		},
		togglerCursor(status: string) {
			enum IAnalysisStatus {
				Pending = "กำลังดำเนินการ",
				Success = "สำเร็จ",
			}

			if (status === IAnalysisStatus.Success) {
				return "cursor-pointer"
			} else {
				return ""
			}
		},
	},
	getters: {
		getStrategicsOptions(state) {
			if (!state || !state.strategics || state.strategics.length === 0) {
				return []
			}

			const options: IOption[] = state.strategics.map((item) => {
				return { label: item.name, value: item.id }
			})

			return options
		},
	},
})
