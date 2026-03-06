<script setup lang="ts">
import { useForm } from "vee-validate"

import { useResetPasswordStore } from "../store"

const store = useResetPasswordStore()
const { handleSubmit } = useForm()

const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	await store.resetPassword()
})
</script>

<template>
	<form @submit="onSubmit">
		<template v-if="!store.status">
			<div class="fv-row mb-3">
				<div class="col-12">
					<VTextInput v-model="store.email" name="email" label="ชื่อผู้ใช้" :disabled="true" />
				</div>
			</div>
			<div class="fv-row mb-3">
				<div class="col-12">
					<VPasswordInput v-model="store.newPassword" name="new_password" label="รหัสผ่านใหม่" />
				</div>
			</div>
			<div class="fv-row mb-3">
				<div class="col-12">
					<VPasswordInput v-model="store.confirmNewPassword" name="confirm_new_password" label="ยืนยันรหัสผ่านใหม่" />
				</div>
			</div>

			<div class="d-grid mb-15 pb-3 pt-3">
				<BtnSubmit :loading="store.loading" :disabled="store.loading" label="ยืนยัน" />
			</div>
		</template>
		<template v-else-if="store.status">
			<div class="d-grid mb-15 pb-0 pt-3 text-center">
				<h5 class="fw-normal text-success">ดำเนินการรีเซ็ตรหัสผ่านเรียบร้อยแล้ว</h5>
			</div>
		</template>
	</form>
</template>
