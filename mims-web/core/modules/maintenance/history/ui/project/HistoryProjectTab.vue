<script setup lang="ts">
const route = useRoute()
const id = route.params.id
const tab = route.path.split("/")[4]
</script>

<template>
	<div class="mt-0">
		<ul class="nav nav-tabs nav-line-tabs mb-5">
			<li class="nav-item" data-bs-toggle="tab" data-bs-target="" role="tab" aria-selected="true">
				<NuxtLink
					:to="`/maintenances/history/${id}/info`"
					class="nav-link cursor-pointer"
					:class="{ active: tab === 'info' }"
				>
					ข้อมูลโครงการ
				</NuxtLink>
				<span class="line"></span>
			</li>
			<li class="nav-item" data-bs-toggle="tab" data-bs-target="" role="tab" aria-selected="true">
				<NuxtLink
					:to="`/maintenances/history/${id}/plan`"
					class="nav-link cursor-pointer"
					:class="{ active: tab === 'plan' }"
				>
					แผนการดำเนินงาน / งวดงาน
				</NuxtLink>
				<span class="line"></span>
			</li>
			<li class="nav-item" data-bs-toggle="tab" data-bs-target="" role="tab" aria-selected="true">
				<NuxtLink
					:to="`/maintenances/history/${id}/file`"
					class="nav-link cursor-pointer"
					:class="{ active: tab === 'file' }"
				>
					เอกสารและรูปภาพ
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
	padding: 10px 22px;
}
</style>
