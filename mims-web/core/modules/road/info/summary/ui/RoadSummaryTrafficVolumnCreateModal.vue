<script setup lang="ts">
import { useForm } from "vee-validate"
import { useRoadSummaryStore } from "../store"
import { ITrafficItem, ITrafficModel } from "../infrastructure"
import { IValidate } from "~/core/shared/types/Validate"

// Modal
const { $bootstrap }: any = useNuxtApp()
const store = useRoadSummaryStore()
const modal = ref()

const route = useRoute()
const roadId = Number(route.params.roadId)

const showModal = () => {
	const modalElement = modal.value
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	bootstrapModal.show()
	handleReset()
}

const hideModal = () => {
	const modalElement = modal.value
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}

const validate = computed(() => {
	const validation: IValidate = {}
	validation.surveyed_date = "required"
	validation.veh1 = "required"
	validation.veh2 = "required"
	validation.veh3 = "required"

	return validation
})

const { handleSubmit, handleReset } = useForm({ validationSchema: validate })
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.createTraffic()
	if (res?.status) {
		console.log(res)
		await store.getTrafficRevision(roadId)
		// Dismiss modal
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: async () => {
				await store.getTrafficRevision(roadId)
				const data = store.trafficRevision
					.flatMap((el: ITrafficModel) => el.items)
					.find((child: ITrafficItem) => child.id === res.data.id)
				store.toggle.parentID = data?.id_parent ?? 0
				store.toggle.year = data?.year ?? 0
				store.toggle.aadtId = data?.id ?? 0
				store.getTrafficDetail(store.toggle.aadtId)
				store.createTrafficLine()
				// 		store.$reset()
				// 		await store.getTrafficRevision(groupID).then(() => store.setUpdatedOptions(res.data.id, res.data.id_parent))
				// 		store.createLine()
			},
		})
		hideModal()
		handleReset()
	}
})

defineExpose({
	showModal,
	hideModal,
})
</script>

<template>
	<div id="modal" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0">
					<h1 class="modal-title fw-semibold">ข้อมูลปริมาณจราจร</h1>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<form @submit="onSubmit">
					<div class="modal-body">
						<div class="row">
							<div class="col-12 mb-3">
								<VDatePicker
									v-model="store.trafficParams.surveyedDate"
									name="surveyed_date"
									label="วันที่สำรวจ"
									:required="true"
								/>
							</div>
							<div class="col-12 mb-3">
								<VNumberInput
									v-model="store.trafficParams.veh1"
									name="veh1"
									label="ปริมาณจราจร รถ 4 ล้อ"
									:required="true"
									:precision="0"
									:min="0"
									align="start"
									text-end="คัน/วัน"
								/>
							</div>
							<div class="col-12 mb-3">
								<VNumberInput
									v-model="store.trafficParams.veh2"
									name="veh2"
									label="ปริมาณจราจร รถ 6 ล้อ"
									:required="true"
									:precision="0"
									:min="0"
									align="start"
									text-end="คัน/วัน"
								/>
							</div>
							<div class="col-12 mb-3">
								<VNumberInput
									v-model="store.trafficParams.veh3"
									name="veh3"
									label="ปริมาณจราจร รถมากกว่า 6 ล้อ"
									:required="true"
									:precision="0"
									:min="0"
									align="start"
									text-end="คัน/วัน"
								/>
							</div>
						</div>
					</div>
					<div class="modal-footer justify-content-end">
						<VLoading :loading="store.loading" />
						<div>
							<BtnCancel data-bs-dismiss="modal" />
							<BtnSubmit label="บันทึก" />
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
