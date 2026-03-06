<script lang="ts" setup>
import { useAuthStore } from "~~/core/modules/auth/authentication/store"

definePageMeta({
	layout: "auth",
})

useHead({
	title: "เซสชั่นหมดอายุ",
	bodyAttrs: {
		class: "error-bg",
		id: "kt_body",
	},
})

const authStore = useAuthStore()
useStoreLifecycle(authStore)

onMounted(() => {
	if (authStore.isLoggedIn === false) {
		return navigateTo("/")
	}

	showToast({
		title: "แจ้งเตือน",
		message: "หมดอายุเข้าใช้งานระบบ โปรดเข้าสู่ระบบอีกครั้ง",
		type: "primary",
		callBack: () => {
			return authStore.logout()
		},
	})
})

const handleError = () => {
	authStore.logout()
}

</script>

<template>
	<div class="d-flex h-100">
		<div class="d-flex flex-column flex-root m-auto">
			<div class="d-flex flex-column flex-center flex-column-fluid">
				<div class="d-flex flex-column flex-center text-center p-10">
					<div class="card rounded-5 w-sm-475px m-5 m-md-0">
						<div class="card-body py-16 py-md-18 px-8">
							<h2 class="fw-semibold text-gray-800 mb-8">หมดอายุเข้าใช้งานระบบ โปรดเข้าสู่ระบบอีกครั้ง</h2>
							<div class="mb-16 mt-12">
								<img src="/images/errors/401-timeout.png" class="mw-100 mh-125px theme-light-show" alt="" />
							</div>
							<div class="mb-0">
								<button class="btn btn-primary" @click="handleError">
									<i class="fi fi-rr-arrow-small-left fs-3"></i> ไปยังหน้าเข้าสู่ระบบ
								</button>
							</div>
							<div class="text-gray-500 text-center fw-normal fs-6 mt-12 mb-1">{{ useCopyright() }}</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped></style>
