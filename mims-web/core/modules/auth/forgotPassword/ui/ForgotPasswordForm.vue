<script setup lang="ts">
import { useForm } from "vee-validate"
import { useForgotPasswordStore } from "../store"

const store = useForgotPasswordStore()
useStoreLifecycle(store)
const status = ref()

const { handleSubmit } = useForm()
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	status.value = await store.forgotPassword()
})

</script>

<template>
	<form class="form w-100 px-11" @submit="onSubmit">
		<div class="text-center mb-10">
			<img src="/images/logos/logo.png" class="mb-5" width="125" alt="logo" />
			<h2 class="text-dark fw-semibold mb-5 mt-5">ลืมรหัสผ่าน ?</h2>
		</div>
		<template v-if="!status">
			<div class="fv-row mb-3">
				<div class="col-12">
					<VTextInput v-model="store.email" name="email" label="อีเมล" :required="true" />
				</div>
			</div>

			<div class="d-grid mb-15 pb-3 pt-3">
				<BtnSubmit :loading="store.loading" :disabled="store.loading" label="รีเซ็ตรหัสผ่าน" />
			</div>
		</template>
		<template v-else-if="status">
			<div class="d-grid mb-15 pb-0 pt-3 text-center">
				<h5 class="fw-normal text-success">ดำเนินการเสร็จสิ้น โปรดตรวจสอบอีเมลของท่าน</h5>
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
