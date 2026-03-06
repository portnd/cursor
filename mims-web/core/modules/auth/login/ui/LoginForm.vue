<script setup lang="ts">
import { useForm } from "vee-validate"
import { useLoginStore } from "../store"
import { useInitDataStore } from "~~/core/modules/initData/store"
import { useInitMenuStore } from "~~/core/modules/initMenu/store"
import { useInitUserStore } from "~~/core/modules/initUser/store"
import { useAuthStore } from "~~/core/modules/auth/authentication/store"

const store = useLoginStore()
useStoreLifecycle(store)

const { handleSubmit } = useForm()
const onSubmit = handleSubmit(async (_, actions) => {
	useAction(actions)
	const res = await store.login()

	if (res?.status) {
		const authStore = useAuthStore()
		authStore.setToken(res.data.access_token, res.data.refresh_token)

		const initDataStore = useInitDataStore()
		const initMenuStore = useInitMenuStore()
		const initUserStore = useInitUserStore()

		// Only await initUser + initMenu (needed for permission + menu). initData is slow (~30s);
		// run it in background so redirect is fast; dashboard/app will use it when ready.
		await Promise.all([initMenuStore.initMenu(), initUserStore.initUser()])
		void initDataStore.initData()

		// Navigate directly to target (skip "/" and extra redirect)
		return navigateTo(useDefaultLandingUrl())
	}
})

const forgotPassword = () => {
	showAlert({
		title: "ลืมรหัสผ่าน ?",
		html: "โปรดติดต่อ ผู้ดูแลระบบ",
		type: "info",
		// confirmText: "ปิด",
	})
}

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<form class="form w-100 px-0" @submit="onSubmit">
		<div class="text-center mb-8">
			<img src="/images/logos/logo.png" class="mb-4" width="120" alt="logo" />
			<h3 class="text-dark fw-semibold mb-0 mt-5 text-uppercase project-name">{{ useProjectName() }}</h3>
			<h4 class="text-gray-500 mt-1 mb-5 fw-light">{{ useProjectNameTH() }}</h4>
			<h2 class="fw-semibold mb-0">เข้าสู่ระบบ</h2>
		</div>

		<div class="fv-row px-sm-10 mb-5">
			<div class="col-12">
				<VTextInput v-model="store.username" name="username" label="ชื่อผู้ใช้" />
			</div>
		</div>

		<div class="fv-row px-sm-10 mb-3">
			<div class="col-12">
				<VPasswordInput v-model="store.password" name="password" label="รหัสผ่าน" />
			</div>
		</div>

		<div class="d-flex flex-stack flex-wrap gap-3 fs-base fw-semibold px-sm-10">
			<div />
			<NuxtLink class="fw-semibold form-text fs-6 me-4 mt-0 cursor-pointer" @click="forgotPassword">
				ลืมรหัสผ่าน ?
			</NuxtLink>
		</div>

		<div class="d-grid mb-5 pt-5 px-sm-10">
			<BtnSubmit :loading="store.loading" label="เข้าสู่ระบบ" />
		</div>
		<div class="text-gray-500 text-center fw-normal fs-6 mt-8 mb-0">{{ useCopyright() }}</div>
	</form>
</template>

<style scoped>
.project-name {
	line-height: 1.5;
	font-size: 16px;
}

@media screen and (max-width: 576px) {
	h3 {
		font-size: 13px;
	}
}
</style>
