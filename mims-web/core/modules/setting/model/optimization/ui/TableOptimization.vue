<script setup lang="ts">
import { useForm } from "vee-validate"
import { useOptimizationStore } from "../store"

const { handleSubmit } = useForm({
	validationSchema: {
		bc_ratio_constraint: "required",
		default_design_life: "required",
	},
})

const store = useOptimizationStore()
useStoreLifecycle(store)

onMounted(() => {
	store.get()
})

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.post()
	if (res) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: () => {
				store.get()
			},
		})
	}
})

const onCancel = () => {
	store.get()
}

</script>

<template>
	<VSkeletonLoader :loading="store.loading">
		<form @submit="onSubmit">
			<div class="row">
				<div class="col-12">
					<div class="table-responsive">
						<table class="table customize-basic-table mb-0">
							<thead>
								<tr>
									<th class="text-center fw-semibold">ประเภท</th>
									<th class="text-center fw-semibold">ค่า <span class="required"></span></th>
								</tr>
							</thead>
							<tbody>
								<tr>
									<td class="text-center">B/C RATIO CONSTRAINT</td>
									<td>
										<VNumberInput
											v-model="store.data.bc_ratio_constraint"
											name="bc_ratio_constraint"
											:precision="2"
											:required="true"
										/>
									</td>
								</tr>
								<tr>
									<td class="text-center">DEFAULT DESIGN LIFE</td>
									<td>
										<VNumberInput
											v-model="store.data.default_design_life"
											name="default_design_life"
											:precision="2"
											:required="true"
										/>
									</td>
								</tr>
							</tbody>
						</table>
					</div>
				</div>
			</div>
			<div class="row mt-5">
				<div class="col-12 text-end">
					<VLoading :loading="store.loading" />
					<div>
						<BtnCancel @click="onCancel()" />
						<BtnSubmit label="บันทึก" :disabled="store.loading" />
					</div>
				</div>
			</div>
		</form>
	</VSkeletonLoader>
</template>

<style scoped>
thead tr th {
	padding: 18px 0px;
	border: 0;
}

thead tr th:first-of-type {
	border-top-left-radius: 7px;
}

thead tr th:last-of-type {
	border-top-right-radius: 7px;
}

tbody tr td {
	padding-left: 1.5rem !important;
	padding-right: 1.5rem !important;
	padding-bottom: 0.5rem;
	padding-top: 0.5rem;
	vertical-align: middle;
	width: 50%;
}
/* @media only screen and (max-width: 768px) {
	.customize-basic-table {
		width: max-content;
	}
} */

.customize-basic-table tr td {
	border-right: 1px solid var(--kt-gray-300) !important;
}

.customize-basic-table tr td:last-of-type {
	border: 0 !important;
}

.customize-basic-table tr:last-of-type td:last-of-type,
.customize-basic-table tr:hover:last-of-type td:last-of-type {
	border-radius: 0px !important;
}
</style>
