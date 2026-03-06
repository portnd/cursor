<script lang="ts" setup></script>

<template>
	<!--begin::Header-->
	<div id="kt_header" class="header bg-white align-items-stretch d-lg-none">
		<!--begin::Container-->
		<div class="d-flex align-items-stretch justify-content-between">
			<!--begin::Aside mobile toggle-->
			<div class="d-flex align-items-center d-lg-none ms-n3 me-1" title="Show aside menu">
				<div id="kt_aside_mobile_toggle" class="btn btn-icon btn-active-color-primary w-40px h-40px">
					<i class="fi-br-menu-burger text-gray-600 ms-15"></i>
				</div>
			</div>
			<!--end::Aside mobile toggle-->

			<!--begin::Mobile logo-->
			<div class="d-flex align-items-center flex-grow-1 flex-lg-grow-0 ms-8">
				<a href="/" class="d-lg-none">
					<img alt="logo" src="/images/logos/logo.png" class="w-50px" />
				</a>
			</div>
			<!--end::Mobile logo-->
		</div>
		<!--end::Container-->
	</div>
	<!--end::Header-->
</template>

<style scoped>
#kt_aside_mobile_toggle i {
	font-size: 1.8rem;
}
</style>
