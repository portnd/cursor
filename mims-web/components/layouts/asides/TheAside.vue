<script lang="ts" setup>
import TheBrand from "~~/components/layouts/asides/TheBrand.vue"
import TheMenu from "~~/components/layouts/asides/TheMenu.vue"
import TheProfileMenu from "~~/components/layouts/asides/TheProfileMenu.vue"
</script>

<template>
	<div
		id="kt_aside"
		class="aside overflow-visible pb-5 pt-3 pt-lg-0"
		data-kt-drawer="true"
		data-kt-drawer-name="aside"
		data-kt-drawer-activate="{default: true, lg: false}"
		data-kt-drawer-overlay="true"
		data-kt-drawer-width="{default:'80px', '300px': '130px'}"
		data-kt-drawer-direction="start"
		data-kt-drawer-toggle="#kt_aside_mobile_toggle"
	>
		<TheBrand />

		<!--begin::Aside menu-->
		<div class="aside-menu flex-column-fluid">
			<TheMenu />
		</div>
		<!--end::Aside menu-->

		<!--begin::Footer-->
		<div id="kt_aside_footer" class="aside-footer flex-column-auto">
			<!--begin::Menu-->
			<TheProfileMenu />
			<!--end::Menu-->
		</div>
		<!--end::Footer-->
	</div>
</template>

<style scoped></style>
