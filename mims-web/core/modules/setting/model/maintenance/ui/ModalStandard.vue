<script setup lang="ts">
import { useForm } from "vee-validate"
import { useMaintenanceEditStore, useMaintenanceSequenceStore } from "../store"
import { IMaintenanceSequence } from "../infrastructure"

// Modal
const { $bootstrap } = useNuxtApp()
const modal = ref()

const store = useMaintenanceSequenceStore()
useStoreLifecycle(store)
const MaintenanceEditStore = useMaintenanceEditStore()
const type = [
	{ name: "ลาดยาง (AC)", value: "asphalt" },
	{ name: "คอนกรีต (Concrete)", value: "concrete" },
]
const drag = ref(false)
const dragOption = computed(() => {
	return {
		animation: 200,
		group: "description",
		disabled: false,
		ghostClass: "drag-item",
	}
})
const showModal = async () => {
	await store.get()
	const modalElement = modal.value
	const bootstrapModal = new $bootstrap.Modal(modalElement)
	bootstrapModal.show()
}

const hideModal = () => {
	const modalElement = modal.value
	const bootstrapModal = $bootstrap.Modal.getInstance(modalElement)
	bootstrapModal?.hide()
}
const { handleSubmit } = useForm()
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)

	const res = await handleData()
	if (res.status === true) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: () => {
				MaintenanceEditStore.get()
				hideModal()
			},
		})
	} else {
		useHandlerError(res.code, res.error, { showAlert: true })
	}
})

const handleData = async () => {
	const result: IMaintenanceSequence = {}

	for (let i = 0; i < type.length; i++) {
		if (store.data[type[i].value]) {
			// กรณีเป็นค่าว่าง
			result[type[i].value] = store.data[type[i].value].map((item: any) => item.id)
		}
	}

	const res = await store.post(result)
	return res
}

onMounted(async () => {
	await store.get()
})

defineExpose({
	showModal,
	hideModal,
})

</script>
<template>
	<div id="modal-standard" ref="modal" class="modal fade" data-bs-keyboard="false">
		<div class="modal-dialog modal-lg modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header border-0">
					<h3 class="modal-title fw-semibold">จัดลำดับวิธีการซ่อมบำรุง</h3>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<form @submit="onSubmit">
					<div class="modal-body">
						<div class="row">
							<template v-for="(item, key) in type" :key="key">
								<div class="col-lg-6">
									<table class="table customize-basic-table draggable mb-3 mb-lg-0">
										<thead>
											<tr>
												<th>
													<label class="fw-semibold ms-1 fs-6">{{ item.name }}</label>
												</th>
											</tr>
										</thead>
										<tbody>
											<tr>
												<td class="draggable-content">
													<label v-if="!store.data[item.value]">ไม่พบข้อมูล</label>
													<Draggable
														v-else
														v-model="store.data[item.value]"
														item-key="index"
														tag="el-collapse"
														:component-data="{ tag: 'ul', name: 'flip-list', type: 'transition' }"
														v-bind="dragOption"
														@start="drag = true"
														@end="drag = false"
													>
														<template #item="{ element, index }">
															<table class="table draggable-item cursor-move mb-0">
																<tbody>
																	<tr>
																		<td>{{ index + 1 }}. มาตรฐาน {{ element.name }}</td>
																	</tr>
																</tbody>
															</table>
														</template>
													</Draggable>
												</td>
											</tr>
										</tbody>
									</table>
								</div>
							</template>
						</div>
					</div>
					<div class="modal-footer pt-0 pt-lg-1">
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

<style scoped lang="scss">
.flip-list-move {
	transition: transform 0.5s;
}
.no-move {
	transition: transform 0s;
}
.drag-item {
	opacity: 0.5;
	background: #c8ebfb;
}
.list-group {
	min-height: 20px;
}
.list-group-item {
	cursor: move;
}
.list-group-item i {
	cursor: pointer;
}

.draggable {
	th {
		background-color: var(--kt-gray-300);
		border-top-left-radius: 0.25rem;
		border-top-right-radius: 0.25rem;
		padding: 0.75rem;
	}
}

.draggable-content {
	padding: 10px 12px;
	min-height: 420px;
	display: block;
}
.draggable-item {
	border-radius: 0.75rem;
	margin-top: 8px;
	margin-bottom: 8px !important;
	background-color: var(--kt-input-bg);
}
.draggable-item:first-child {
	margin-top: 2px;
}
.draggable-item:last-child {
	margin-bottom: 2px !important;
}
</style>
