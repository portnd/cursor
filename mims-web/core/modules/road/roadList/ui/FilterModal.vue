<script setup lang="ts">
import { useRoadListStore } from "../store"

const store = useRoadListStore()

// Modal
const { $bootstrap }: any = useNuxtApp()

const modal = ref()

const showModal = () => {
	const modalElement = modal.value
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	bootstrapModal.show()
}

const hideModal = () => {
	const modalElement = modal.value
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}

defineExpose({
	showModal,
	hideModal,
})
</script>

<template>
	<div id="modal" ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false">
		<div class="modal-dialog modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header pb-5">
					<h1 class="modal-title fw-semibold">สภาพทาง</h1>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<div class="modal-body" @click.stop>
					<div class="row">
						<div class="col-12 mt-5 mb-3">
							<VCheckbox
								v-model="store.condition_all"
								name="checkAll"
								:option="{ label: 'ทั้งหมด' }"
								mode="single"
								@update:model-value="(e:boolean) => store.conditionCheckAll(e)"
							/>
						</div>
						<div class="col-12 mb-3">
							<div>
								<VCheckbox
									v-model="store.condition_parent.iri_1000"
									:option="{ label: 'ค่าดัชนีความขรุขระสากล (IRI) ที่ทดสอบทุกระยะ 1,000 เมตร', value: 0 }"
									mode="single"
									name="check1"
									:inline="false"
									@update:model-value="store.onUpdateCheckBox('iri_1000')"
								/>
							</div>
							<div class="ps-6">
								<VRadio
									v-model="store.condition_temp.is_iri_1000"
									:options="[
										{ label: 'ผ่านเกณฑ์เงื่อนไข RFP', color: '#47BE7D', value: true },
										{ label: 'ไม่ผ่านเกณฑ์เงื่อนไข RFP', color: '#D9214E', value: false },
									]"
									name="radio1"
									:inline="false"
									:disabled="!store.condition_parent.iri_1000"
								/>
							</div>
						</div>

						<div class="col-12 mb-3">
							<div>
								<VCheckbox
									v-model="store.condition_parent.iri_100"
									:option="{ label: 'ค่าดัชนีความขรุขระสากล (IRI) ที่ทดสอบทุกระยะ 100 เมตร ', value: 0 }"
									mode="single"
									name="check2"
									:inline="false"
									@update:model-value="store.onUpdateCheckBox('iri_100')"
								/>
							</div>
							<div class="ps-6">
								<VRadio
									v-model="store.condition_temp.is_iri_100"
									:options="[
										{ label: 'ผ่านเกณฑ์เงื่อนไข RFP', color: '#47BE7D', value: true },
										{ label: 'ไม่ผ่านเกณฑ์เงื่อนไข RFP', color: '#D9214E', value: false },
									]"
									name="radio2"
									:inline="false"
									:disabled="!store.condition_parent.iri_100"
								/>
							</div>
						</div>

						<div class="col-12 mb-3">
							<div>
								<VCheckbox
									v-model="store.condition_parent.rut_100"
									:option="{ label: 'ค่าความลึกร่องล้อ (RUT) ที่ทดสอบทุกระยะ 100 เมตร' }"
									mode="single"
									name="check3"
									:inline="false"
									@update:model-value="store.onUpdateCheckBox('rut_100')"
								/>
							</div>
							<div class="ps-6">
								<VRadio
									v-model="store.condition_temp.is_rut_100"
									:options="[
										{ label: 'ผ่านเกณฑ์เงื่อนไข RFP', color: '#47BE7D', value: true },
										{ label: 'ไม่ผ่านเกณฑ์เงื่อนไข RFP', color: '#D9214E', value: false },
									]"
									name="radio3"
									:inline="false"
									:disabled="!store.condition_parent.rut_100"
								/>
							</div>
						</div>

						<div class="col-12 mb-3">
							<div>
								<VCheckbox
									v-model="store.condition_parent.ifi_100"
									:option="{ label: 'ค่าดัชนีความเสียดทานสากล (IFI) ที่ทดสอบทุกระยะ 100 เมตร ' }"
									mode="single"
									name="check4"
									:inline="false"
									@update:model-value="store.onUpdateCheckBox('ifi_100')"
								/>
							</div>
							<div class="ps-6">
								<VRadio
									v-model="store.condition_temp.is_ifi_100"
									:options="[
										{ label: 'ผ่านเกณฑ์เงื่อนไข RFP', color: '#47BE7D', value: true },
										{ label: 'ไม่ผ่านเกณฑ์เงื่อนไข RFP', color: '#D9214E', value: false },
									]"
									name="radio4"
									:inline="false"
									:disabled="!store.condition_parent.ifi_100"
								/>
							</div>
						</div>

						<div class="col-12 mb-3">
							<div>
								<VCheckbox
									v-model="store.condition_parent.g7_100"
									:option="{ label: 'ค่าสะท้อนแสงของเส้นจราจร (G7) ที่ทดสอบทุกระยะ 100 เมตร ' }"
									mode="single"
									name="check5"
									:inline="false"
									@update:model-value="store.onUpdateCheckBox('g7_100')"
								/>
							</div>
							<div class="ps-6">
								<VRadio
									v-model="store.condition_temp.is_g7_100"
									:options="[
										{ label: 'ผ่านเกณฑ์เงื่อนไข RFP', color: '#47BE7D', value: true },
										{ label: 'ไม่ผ่านเกณฑ์เงื่อนไข RFP', color: '#D9214E', value: false },
									]"
									name="radio5"
									:inline="false"
									:disabled="!store.condition_parent.g7_100"
								/>
							</div>
						</div>
					</div>
				</div>
				<div class="modal-footer d-flex justify-content-end">
					<!-- <VLoading :loading="store.loading" /> -->
					<BtnCancel data-bs-dismiss="modal" @click.stop="store.onCancel" />
					<BtnSubmit
						data-bs-dismiss="modal"
						:disabled="store.loading"
						:loading="store.loading"
						label="บันทึก"
						@click.stop="store.onSubmit"
					/>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
