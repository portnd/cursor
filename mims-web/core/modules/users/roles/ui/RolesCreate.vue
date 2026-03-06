<script setup lang="ts">
import { useForm } from "vee-validate"
import { useRolesCreateStore } from "../store"
import {
	IRolesAccessGroup,
	IRolesAccessControl,
	IRolesAccessDetail,
	IRequestRolesAccessDetail,
} from "../infrastructure"
import { IValidate } from "~/core/shared/types/Validate"
import { useInitDataStore } from "~/core/modules/initData/store"

const store = useRolesCreateStore()
useStoreLifecycle(store)

const handleValidate = computed(() => {
	const validations: IValidate = {
		name: "required",
	}
	return validations
})

const result = ref<Array<IRequestRolesAccessDetail>>([])
const { handleSubmit, isSubmitting, errors } = useForm({ validationSchema: handleValidate })

const onSubmit = handleSubmit(async (_, actions) => {
	result.value = store.control
		.flatMap((el: IRolesAccessGroup) => el.menu)
		.flatMap((el: IRolesAccessControl) => el.access_control)
		.filter((el: IRolesAccessDetail) => el.is_check === true)
		.map((el: IRolesAccessDetail) => ({ access_control_id: el.id }))

	if (result.value.length < 1) {
		showAlert({
			title: "แจ้งเตือน",
			message: "โปรดระบุสิทธิ์การใช้งานอย่างน้อย 1 ตัวเลือก",
			type: "warning",
		})
	} else {
		useAction(actions)
		const res = await store.create()
		if (res?.status) {
			useHandlerSuccess(res.code, {
				showAlert: true,
				fn: function () {
					const initDataStore = useInitDataStore()
					initDataStore.initData()
					store.clearInput()
					navigateTo("/users/roles")
				},
			})
		}
	}
})

watch(
	() => store.control,
	() => {
		console.log("watch:")
		console.log(store.currentItemChecked)
		if (store.currentItemChecked !== undefined) {
			store.updateCheckedByRelation(store.currentItemChecked.is_check)
		}
	},
	{
		deep: true,
	}
)

// เลื่อนไปที่ฟิลด์ Error
watch(isSubmitting, () => {
	if (Object.keys(errors.value).length > 0) {
		scrollIntoInvalidField()
	}
})

const onCancel = () => {
	return navigateTo("/users/roles")
}

onMounted(() => {
	store.get()
})

</script>
<template>
	<div class="row">
		<div class="col-xl-12">
			<div class="card p-5 mb-5">
				<form @submit.prevent="onSubmit">
					<div class="row mb-5">
						<div class="col-lg-6">
							<VTextInput v-model="store.name" label="ชื่อกลุ่มสิทธิ์การใช้งาน" name="name" :required="true" />
						</div>
					</div>
					<div class="row">
						<div class="col-xl-12">
							<table class="table customize-basic-table mb-0">
								<thead>
									<tr class="text-center">
										<th>ลำดับ</th>
										<th>เมนู</th>
										<th>สิทธิ์การใช้งาน</th>
									</tr>
								</thead>
								<tbody>
									<template v-for="(level1, index1) in store.control" :key="index1">
										<tr v-if="level1.menu?.length > 1">
											<td class="text-center">{{ index1 + 1 }}</td>
											<td colspan="2">
												<span class="fw-semibold">{{ level1.name }}</span>
											</td>
										</tr>
										<tr v-if="level1.menu === null">
											<td class="text-center">{{ index1 + 1 }}</td>
											<td colspan="2">
												<span class="fw-semibold">{{ level1.name }}</span>
											</td>
										</tr>

										<template v-for="(level2, index2) in level1.menu" :key="index2">
											<tr class="level2">
												<td
													v-if="level1.name === level2.name"
													class="text-center"
													:class="`${level1.menu.length === 1 ? 'showBorder' : ''}`"
												>
													{{ index1 + 1 }}
												</td>
												<td v-else :class="`${index2 + 1 === level1.menu.length ? 'showBorder' : ''}`"></td>
												<td v-if="index2 + 1 === level1.menu.length" class="showBorder">
													<span :class="`${level1.menu.length === 1 ? 'fw-semibold' : ''}`">{{ level2.name }}</span>
												</td>
												<td v-else>
													<span>{{ level2.name }}</span>
												</td>
												<td v-if="index2 + 1 === level1.menu.length" class="showBorder">
													<div
														v-for="(level3, index3) in level2.access_control"
														:key="index3"
														:class="`${level2.access_control.length > 1 ? 'level2-list' : ''}`"
														class="showBorder"
													>
														<VCheckbox
															v-model="store.control[index1].menu[index2].access_control[index3].is_check"
															:option="{ label: level3.name }"
															:name="level3.access_key"
															:inline="false"
															mode="single"
															@update:model-value="(isCheck:boolean) => {
																console.log('isCheck', isCheck)
																console.log('access_key', store.control[index1].menu[index2].access_control[index3].access_key)

																store.updateCurrentItemChecked(store.control[index1].menu[index2].access_control[index3])



															}"
														/>
													</div>
												</td>
												<td v-else>
													<div
														v-for="(level3, index3) in level2.access_control"
														:key="index3"
														:class="`${level2.access_control.length > 1 ? 'level2-list' : ''}`"
														class="showBorder"
													>
														<VCheckbox
															v-model="store.control[index1].menu[index2].access_control[index3].is_check"
															:option="{ label: level3.name }"
															:name="level3.access_key"
															:inline="false"
															mode="single"
														/>
													</div>
												</td>
											</tr>
										</template>
									</template>
								</tbody>
								<!-- <tbody>
									<template v-for="(level1, index1) in store.control" :key="index1">
										<tr v-if="level1.menu?.length > 1">
											<td class="text-center">{{ index1 + 1 }}</td>
											<td colspan="2">
												<span class="fw-semibold">{{ level1.name }}</span>
											</td>
										</tr>
										<tr v-if="level1.menu === null">
											<td class="text-center">{{ index1 + 1 }}</td>
											<td colspan="2">
												<span class="fw-semibold">{{ level1.name }}</span>
											</td>
										</tr>

										<template v-for="(level2, index2) in level1.menu" :key="index2">
											<tr class="level2">
												<td
													v-if="level1.name === level2.name"
													class="text-center"
													:class="`${level1.menu.length === 1 ? 'showBorder' : ''}`"
												>
													{{ index1 + 1 }}
												</td>
												<td v-else :class="`${index2 + 1 === level1.menu.length ? 'showBorder' : ''}`"></td>
												<td v-if="index2 + 1 === level1.menu.length" class="showBorder">
													<span :class="`${level1.menu.length === 1 ? 'fw-semibold' : ''}`">{{ level2.name }}</span>
												</td>
												<td v-else>
													<span>{{ level2.name }}</span>
												</td>
												<td v-if="index2 + 1 === level1.menu.length" class="showBorder">
													<div
														v-for="(level3, index3) in level2.access_control"
														:key="index3"
														:class="`${level2.access_control.length > 1 ? 'level2-list' : ''}`"
														class="showBorder"
													>
														<VCheckbox
															:model-value="store.control[index1].menu[index2].access_control[index3].is_check"
															:option="{ label: level3.name }"
															:name="level3.access_key"
															:inline="false"
															mode="single"
															@click="
																() =>
																	store.updateRole(
																		store.control[index1].menu[index2].access_control[index3],
																		store.control[index1].menu[index2].access_control
																	)
															"
														/>
													</div>
												</td>
												<td v-else>
													<div
														v-for="(level3, index3) in level2.access_control"
														:key="index3"
														:class="`${level2.access_control.length > 1 ? 'level2-list' : ''}`"
														class="showBorder"
													>
														<VCheckbox
															:model-value="store.control[index1].menu[index2].access_control[index3].is_check"
															:option="{ label: level3.name }"
															:name="level3.access_key"
															:inline="false"
															mode="single"
															@click="
																() =>
																	store.updateRole(
																		store.control[index1].menu[index2].access_control[index3],
																		store.control[index1].menu[index2].access_control
																	)
															"
														/>
													</div>
												</td>
											</tr>
										</template>
									</template>
								</tbody> -->
							</table>
						</div>
					</div>
					<div class="d-flex justify-content-end mt-5">
						<BtnCancel @click="onCancel" />
						<BtnSubmit label="บันทึก" />
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style lang="scss" scoped>
tr.level2 {
	&:nth-child(even) {
		background-color: rgba(var(--kt-gray-300), 0.3);
	}
}
td.showBorder {
	border-bottom: solid 1px #ddd;
}
.level2-list {
	margin-bottom: 1.25em;
	&:nth-last-child(1) {
		margin-bottom: 0;
	}
}
</style>
