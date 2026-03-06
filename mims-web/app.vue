<script lang="ts" setup>
import { initializeComponents } from "~~/composables/useTheme"
import { useAuthStore } from "~~/core/modules/auth/authentication/store"
import { useInitDataStore } from "~~/core/modules/initData/store"
import { useInitMenuStore } from "~~/core/modules/initMenu/store"
import { useInitUserStore } from "~~/core/modules/initUser/store"

const route = useRoute()

onBeforeMount(() => {
	// Auto refresh JWT Token
	const authStore = useAuthStore()

	// subscribe
	if (authStore.isLoggedIn === true) {
		authStore.subscribeStore()

		authStore.autoRefresh().then(async () => {
			// Init Data
			const initDataStore = useInitDataStore()
			await initDataStore.initData()

			// Init Menu
			const initMenuStore = useInitMenuStore()
			await initMenuStore.initMenu()

			// Init User
			const initUserStore = useInitUserStore()
			await initUserStore.initUser()
		})
	}
})

onMounted(() => {
	nextTick(() => {
		try {
			const bodyStore = useBody()
			bodyStore.addBodyClassname("page-loading")
			useTimeoutFn(() => {
				try {
					bodyStore.removeBodyClassName("page-loading")
					bodyStore.addBodyClassname("page-loaded")
				} catch {
					document.body.classList.remove("page-loading")
					document.body.classList.add("page-loaded")
				}
				document.body.style.backgroundColor = ""
			}, 1000)
			initializeComponents()
		} catch {
			document.body.classList.remove("page-loading")
			document.body.classList.add("page-loaded")
			document.body.style.backgroundColor = ""
		}
	})
})
</script>

<template>
	<div class="page">
		<NuxtLayout>
			<NuxtPage :key="route.fullPath" />
		</NuxtLayout>
	</div>

	<!--begin::Splash screen-->
	<div id="splash-screen" class="splash-screen">
		<img src="/images/splash-screen.png" alt="logo" />
		<div class="d-flex">
			<span>กำลังโหลด</span>
			<div class="dot-flashing" />
		</div>
	</div>
	<!--end::Splash screen-->
</template>

<style lang="scss">
// Icons
@import "@flaticon/flaticon-uicons/css/all/all";

// Themes
@import "~~/assets/themes/sass/dark";
@import "~~/assets/themes/sass/plugins";
@import "~~/assets/themes/sass/style";

// Splash screen
@import "~~/assets/sass/splash-screen";
@import "three-dots";

// Customs
@import "~~/assets/sass/custom";

#nuxt {
	display: contents;
	min-height: 100vh;
}
body {
	padding: 0 !important;
	min-height: 100vh;
	/* ป้องกันหน้าขาว: ให้มีพื้นหลังเสมอ (ไม่ซ่อนทั้งหน้าเมื่อไม่มี class) */
	background-color: var(--kt-gray-100, #f5f8fa);
}
body.page-loading {
	div.page {
		display: none;
	}
}
div.page {
	height: 100%;
	min-height: 100vh;
}
</style>
