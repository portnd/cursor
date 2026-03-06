<script setup lang="ts">
import { useForm } from "vee-validate"
import { useChangePasswordStore } from "../store"
import { IValidate } from "~/core/shared/types/Validate"

const store = useChangePasswordStore()
useStoreLifecycle(store)

const validate = computed(() => {
	const validations: IValidate = {}

	validations.current_password = "required"
	validations.new_password = "password"
	validations.confirm_new_password = "password"

	return validations
})

// defineRule("current_password", (value: any) => {
// 	const input = value
// 	let message = ""

// 	if (input === store.params.current_password) {
// 		message = "รหัสผ่านใหม่ต้องไม่เหมือนกับรหัสผ่านเดิม"
// 	}

// 	return message
// })

const { handleSubmit, handleReset } = useForm({ validationSchema: validate })

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)

	const res = await store.updatePassword()

	if (res?.status) {
		useHandlerSuccess(res.code, {
			showAlert: true,
			fn: () => {
				handleReset()
				store.$reset()
			},
		})
	}
})

const onCancel = () => {
	navigateTo("/")
}

onUnmounted(() => {
	store.$reset()
})
</script>
<template>
	<div class="row">
		<div class="col-12">
			<div class="row">
				<div class="col-mb-4 col-sm-4 mb-3">
					<VLabel label="ชื่อผู้ใช้งาน" />
					<input type="text" class="form-control" :disabled="true" :value="useInitUser()?.username" />
				</div>
			</div>
		</div>
		<div class="col-12">
			<div class="row">
				<div class="col-mb-4 col-sm-4 mb-3">
					<VPasswordInput
						v-model="store.params.current_password"
						label="รหัสผ่านเดิม"
						name="current_password"
						:required="true"
					/>
				</div>
				<div class="col-mb-4 col-sm-4 mb-3">
					<VLabel label="รหัสผ่านใหม่" :required="true" />
					<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
						<i class="fi fi-sr-interrogation fs-5"></i>
						<template #content>
							<ul class="mt-4">
								<li>ต้องมีจำนวน 8 ตัวอักษรขึ้นไป</li>
								<li>ต้องประกอบด้วยตัวอักษร A-Z ตัวพิมพ์ใหญ่</li>
								<li>ต้องประกอบด้วยตัวอักษร a-z ตัวพิมพ์เล็ก</li>
								<li>ต้องประกอบด้วยตัวเลข 0-9 อย่างน้อย 1 ตัว</li>
							</ul>
						</template>
					</VPopover>
					<VPasswordInput v-model="store.params.new_password" name="new_password" :required="true" />
				</div>
				<div class="col-mb-4 col-sm-4 mb-3">
					<VLabel label="ยืนยันรหัสผ่านใหม่" :required="true" />
					<VPopover class-name="p-1 pb-3 my-0 text-primary fw-normal form-text cursor-pointer">
						<i class="fi fi-sr-interrogation fs-5"></i>
						<template #content>
							<ul class="mt-4">
								<li>ต้องมีจำนวน 8 ตัวอักษรขึ้นไป</li>
								<li>ต้องประกอบด้วยตัวอักษร A-Z ตัวพิมพ์ใหญ่</li>
								<li>ต้องประกอบด้วยตัวอักษร a-z ตัวพิมพ์เล็ก</li>
								<li>ต้องประกอบด้วยตัวเลข 0-9 อย่างน้อย 1 ตัว</li>
							</ul>
						</template>
					</VPopover>
					<VPasswordInput v-model="store.params.confirm_new_password" name="confirm_new_password" :required="true" />
				</div>
			</div>
		</div>
		<div class="d-flex justify-content-end mt-5">
			<BtnCancel @click="onCancel" />
			<BtnSubmit :disabled="store.loading" :loading="store.loading" label="บันทึก" @click="onSubmit" />
		</div>
	</div>
</template>
<style scoped></style>
