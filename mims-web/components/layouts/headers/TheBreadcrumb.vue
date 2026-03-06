<script setup lang="ts">
interface IBreadcrumbs {
	name: string
	to?: string
}

defineProps({
	breadcrumbs: {
		type: Object as PropType<IBreadcrumbs[]>,
		default: [] as IBreadcrumbs[],
	},
	title: {
		type: String,
		required: true,
	},
})
</script>

<template>
	<div class="row">
		<div class="col-md-12">
			<ol class="breadcrumb text-muted mb-5">
				<template v-for="(breadcrumb, key) in breadcrumbs" :key="key">
					<li v-show="breadcrumb" class="breadcrumb-item pe-2 fw-normal text-gray-500">
						<NuxtLink :to="breadcrumb.to ? breadcrumb.to : ''">
							<span class="pe-2 text-gray-500">{{ breadcrumb.name }}</span>
						</NuxtLink>
					</li>
				</template>
				<li class="breadcrumb-item text-gray-900 fw-semibold">{{ title }}</li>
			</ol>
		</div>
	</div>
</template>

<style scoped>
.breadcrumb-item {
	font-size: 1.1rem;
}

.breadcrumb .breadcrumb-item:after {
	content: "\f12e";
	font-family: uicons-regular-rounded !important;
	font-style: normal;
	font-weight: 500 !important;
}
</style>
