<script setup lang="ts">
const route = useRoute()
const tab = route.path.split("/")[3]
</script>

<template>
	<h2 class="fw-semibold fs-5">ประเภทการตั้งค่า</h2>
	<ul class="nav nav-tabs nav-pills border-0 d-flex pb-1">
		<li class="nav-item">
			<NuxtLink to="/settings/models/traffic" class="nav-link" :class="{ active: tab === 'traffic' }">Traffic</NuxtLink>
		</li>
		<li class="nav-item">
			<NuxtLink to="/settings/models/maintenance" class="nav-link" :class="{ active: tab === 'maintenance' }"
				>Intervention Criteria</NuxtLink
			>
		</li>
		<li class="nav-item">
			<NuxtLink to="/settings/models/road-work-effect" class="nav-link" :class="{ active: tab === 'road-work-effect' }"
				>Road Work Effect</NuxtLink
			>
		</li>
		<li class="nav-item">
			<NuxtLink to="/settings/models/road-user-cost" class="nav-link" :class="{ active: tab === 'road-user-cost' }"
				>Road User Cost</NuxtLink
			>
		</li>
		<li class="nav-item">
			<NuxtLink to="/settings/models/optimization" class="nav-link" :class="{ active: tab === 'optimization' }"
				>Optimization</NuxtLink
			>
		</li>
		<li class="nav-item">
			<NuxtLink to="/settings/models/deterioration" class="nav-link" :class="{ active: tab === 'deterioration' }"
				>Deterioration</NuxtLink
			>
		</li>
	</ul>
</template>

<style scoped lang="scss">
.nav-pills .nav-item {
	width: calc(100% / 6);
	margin-right: 0px;
	text-align: -webkit-center;
	margin-top: 5px;
	margin-bottom: 5px;

	@media (max-width: 1360px) {
		width: calc(100% / 4);
	}
	@media (max-width: 1050px) {
		width: calc(100% / 3);
	}
	@media (max-width: 991px) {
		width: calc(100% / 5);
	}
	@media (max-width: 885px) {
		width: calc(100% / 4);
	}
	@media (max-width: 640px) {
		width: calc(100% / 3);
	}
	@media (max-width: 460px) {
		width: calc(100% / 2);
	}
}
.nav-pills .nav-link {
	width: 95%;
	height: 40px;
	line-height: 1.5;
}
</style>
