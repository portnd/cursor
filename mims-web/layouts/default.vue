<script lang="ts" setup>
import TheAside from "~~/components/layouts/asides/TheAside.vue"
import TheHeader from "~~/components/layouts/headers/TheHeader.vue"
import TheFooter from "~~/components/layouts/TheFooter.vue"

useHead({
	bodyAttrs: {
		class: "header-fixed header-tablet-and-mobile-fixed toolbar-enabled toolbar-fixed aside-enabled aside-fixed",
		id: "kt_body",
	},
})

onMounted(() => {
	nextTick(() => {
		reinitializeComponents()
	})
})
</script>

<template>
	<div class="layout-default">
		<div class="d-flex flex-column flex-root">
			<div class="page d-flex flex-row flex-column-fluid">
				<TheAside />

				<div id="kt_wrapper" class="d-flex flex-column flex-row-fluid wrapper">
					<TheHeader />

					<!-- begin:: Content -->
					<div id="kt_content" class="content d-flex flex-column flex-column-fluid mb-5">
						<!-- begin:: Content Body -->
						<div class="post d-flex flex-column-fluid">
							<div id="kt_content_container" class="container-fluid">
								<slot />
							</div>
						</div>
						<!-- end:: Content Body -->
					</div>

					<TheFooter />
				</div>
			</div>
		</div>
		<BtnBackToTop />
	</div>
</template>

<style scoped>
#kt_content {
	min-height: calc(100vh - 85px) !important;
}
</style>
