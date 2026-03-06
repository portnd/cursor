<script setup lang="ts">
defineProps({
	type: {
		type: null,
		required: true,
	},
})
const route = useRoute()
const tab = route.path.split("/")[4]
</script>

<template>
	<div class="mt-0 pb-1">
		<ul class="nav nav-tabs nav-line-tabs mb-5">
			<li class="nav-item" data-bs-toggle="tab" data-bs-target="#detail-surface" role="tab" aria-selected="true">
				<NuxtLink
					to="/settings/models/road-user-cost/acc"
					class="nav-link cursor-pointer"
					:class="{ active: tab === 'acc' }"
				>
					Accident (ACC)
				</NuxtLink>
				<span class="line"></span>
			</li>
			<li class="nav-item" data-bs-toggle="tab" data-bs-target="#detail-surface" role="tab" aria-selected="true">
				<NuxtLink
					to="/settings/models/road-user-cost/ruc"
					class="nav-link cursor-pointer"
					:class="{ active: tab === 'ruc' }"
				>
					Road User Cost (RUC)</NuxtLink
				>
				<span class="line"></span>
			</li>
		</ul>
	</div>
</template>

<style scoped></style>
