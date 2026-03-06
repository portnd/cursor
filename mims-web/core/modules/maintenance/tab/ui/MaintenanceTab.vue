<script setup lang="ts">
const route = useRoute()
const tab = route.path.split("/")[2]
</script>

<template>
	<div class="mt-0">
		<ul class="nav nav-tabs nav-line-tabs mb-5">
			<li class="nav-item" data-bs-toggle="tab" data-bs-target="" role="tab" aria-selected="true">
				<NuxtLink to="/maintenances/follow-up" class="nav-link cursor-pointer" :class="{ active: tab === 'follow-up' }">
					ติดตามการซ่อมบำรุง
				</NuxtLink>
				<span class="line"></span>
			</li>
			<li class="nav-item" data-bs-toggle="tab" data-bs-target="" role="tab" aria-selected="true">
				<NuxtLink to="/maintenances/history" class="nav-link cursor-pointer" :class="{ active: tab === 'history' }">
					ประวัติการซ่อมบำรุง
				</NuxtLink>
				<span class="line"></span>
			</li>
		</ul>
	</div>
</template>

<style scoped lang="scss">
#nav-tab .active {
	background-color: #fff0d9 !important;
	color: #fdb833 !important;
}

#nav-tab .nav-link {
	background: var(--kt-gray-100);
}

.nav-line-tabs .nav-item .nav-link {
	@media (max-width: 767px) {
		padding: 10px 20px;
	}
}
</style>
