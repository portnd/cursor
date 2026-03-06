<script setup lang="ts">
import { useResetPasswordStore } from "../store"
import { ResetPasswordForm } from "../ui"

const route = useRoute()
const store = useResetPasswordStore()
useStoreLifecycle(store)

onMounted(async () => {
	store.setResetPasswordToken(route.params.token as string)
	await store.checkToken()
})

</script>

<template>
	<form class="form w-100 px-11">
		<div class="text-center mb-10">
			<img src="/images/logos/logo.png" class="mb-5" width="125" alt="logo" />
			<h2 class="text-dark fw-semibold mb-5 mt-5">รีเซ็ตรหัสผ่าน</h2>
		</div>
		<template v-if="store.statusResetPasswordToken">
			<div>
				<ResetPasswordForm />
			</div>
		</template>
		<template v-else>
			<div class="d-grid mb-15 pb-0 pt-3 text-center">
				<h5 class="fw-normal text-danger">หมดเวลาดำเนินการ โปรดดำเนินการใหม่อีกครั้ง</h5>
			</div>
		</template>

		<div class="fv-row mb-5">
			<div class="col-12">
				<div class="text-center">
					<NuxtLink to="/auth/login" class="text-gray-700 mt-2 fs-6"
						><i class="fi fi-rr-arrow-small-left"></i> กลับไปหน้าเข้าสู่ระบบ
					</NuxtLink>
				</div>
			</div>
		</div>

		<div class="text-gray-500 text-center fw-normal fs-6 mt-12 mb-1">{{ useCopyright() }}</div>
	</form>
</template>

<style scoped></style>
