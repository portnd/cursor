<script setup lang="ts">
const route = useRoute()
const tab = route.path.split("/")[4]
</script>

<template>
	<div class="mt-0 pb-1">
		<ul class="nav nav-tabs nav-line-tabs mb-5">
			<li class="nav-item" data-bs-toggle="tab" data-bs-target="#detail-surface" role="tab" aria-selected="true">
				<NuxtLink
					to="/settings/models/traffic/growth-rate"
					class="nav-link cursor-pointer"
					:class="{ active: tab === 'growth-rate' }"
				>
					Growth Rate</NuxtLink
				>
				<span class="line"></span>
			</li>
			<li class="nav-item" data-bs-toggle="tab" data-bs-target="#detail-surface" role="tab" aria-selected="true">
				<NuxtLink
					to="/settings/models/traffic/vehicle-type"
					class="nav-link cursor-pointer"
					:class="{ active: tab === 'vehicle-type' }"
				>
					% Vehicle Type</NuxtLink
				>
				<span class="line"></span>
			</li>
			<li class="nav-item" data-bs-toggle="tab" data-bs-target="#detail-surface" role="tab" aria-selected="true">
				<NuxtLink
					to="/settings/models/traffic/parameter"
					class="nav-link cursor-pointer"
					:class="{ active: tab === 'parameter' }"
				>
					Parameter</NuxtLink
				>
				<span class="line"></span>
			</li>
		</ul>
	</div>
</template>

<style scoped></style>
